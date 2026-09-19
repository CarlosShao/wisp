// Package openaichat is the OpenAI Chat Completions adapter (C5 implementor,
// ticket 09): /v1/chat/completions with tool_calls, SSE streaming, usage
// blocks, and the D37 error classification at the seam.
//
// Wire mapping notes:
//   - normalized ToolUse/ToolResult parts map to assistant tool_calls and
//     role:"tool" messages; argument fragments reassemble by the wire
//     "index" field (gaps/out-of-order tolerated only in compat.loose);
//   - finish_reason mapping: stop->end_turn, length->max_tokens,
//     tool_calls|function_call->tool_use, content_filter->content_filter;
//     missing finish_reason / missing usage are tolerated exactly when the
//     provider's compat switches say so (SPEC-03 sec 3.1: "missing fields are
//     tolerated per the switch, or error out");
//   - Request.ThinkingIntensity is intentionally NOT mapped: the chat
//     completions wire has no portable field, and sending provider-private
//     params to strict endpoints causes 400s. Adapters that own a thinking
//     parameter (anthropic, ticket 11) consume it.
//   - streaming always requests stream_options.include_usage unless
//     compat.loose (several OpenAI-compatible servers reject it).
package openaichat

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

func init() {
	llm.RegisterProtocol("openai-chat", func(ep llm.EndpointOptions) (llm.LlmProvider, error) {
		return New(ep), nil
	})
}

// Adapter implements llm.LlmProvider for OpenAI Chat Completions.
type Adapter struct {
	opts llm.EndpointOptions
	http *http.Client
}

// New builds an adapter. A nil HTTPClient gets a proxy-aware default client.
func New(opts llm.EndpointOptions) *Adapter {
	c := opts.HTTPClient
	if c == nil {
		c = llm.NewDefaultHTTPClient("", "")
	}
	return &Adapter{opts: opts, http: c}
}

// Info implements llm.LlmProvider.
func (a *Adapter) Info() llm.ProviderInfo {
	return llm.ProviderInfo{
		Protocol:         "openai-chat",
		Provider:         a.opts.Provider,
		Model:            a.opts.Model,
		Cache:            llm.CacheSupport{Mode: llm.CacheImplicit, Breakpoints: false},
		MaxContextWindow: a.opts.ContextWindow,
	}
}

// Stream implements llm.LlmProvider. See package llm (events.go) for the
// event contract.
func (a *Adapter) Stream(ctx context.Context, req *llm.Request, emit func(llm.StreamEvent) error) error {
	if err := req.Validate(); err != nil {
		return a.fail(emit, observe.Wrap(observe.ClassInternal, err, "invalid normalized request"))
	}

	payload, err := a.buildWire(req)
	if err != nil {
		return a.fail(emit, err)
	}

	url := strings.TrimRight(a.opts.BaseURL, "/") + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return a.fail(emit, observe.Wrap(observe.ClassInternal, err, "build http request"))
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")
	if a.opts.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+a.opts.APIKey)
	}
	for k, v := range a.opts.Compat.ExtraHeaders {
		httpReq.Header.Set(k, v)
	}

	resp, err := a.http.Do(httpReq)
	if err != nil {
		// ctx cancellation surfaces as a transport error: classify by cause.
		if ctx.Err() != nil {
			return a.cancelled(emit)
		}
		return a.fail(emit, llm.ClassifyTransportError(err))
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
		detail := llm.HTTPErrorDetail{Status: resp.StatusCode}
		if e := parseWireError(body); e != nil {
			detail.Code, detail.Type, detail.Message = e.codeString(), e.Type, e.Message
		} else {
			detail.Message = strings.TrimSpace(string(body))
		}
		return a.fail(emit, llm.NewHTTPError(detail,
			llm.ParseRetryAfter(resp.Header.Get("Retry-After"), time.Now())))
	}
	defer resp.Body.Close()

	return a.consume(ctx, resp.Body, emit)
}

// consume parses the SSE body into C6 events.
func (a *Adapter) consume(ctx context.Context, body io.Reader, emit func(llm.StreamEvent) error) error {
	asm := newToolAssembler(a.loose())
	sawFinish := ""
	sawUsage := false
	dataDone := false // [DONE] sentinel seen

	sc := bufio.NewScanner(body)
	sc.Buffer(make([]byte, 0, 64<<10), 4<<20)
	for sc.Scan() {
		if ctx.Err() != nil {
			return a.cancelled(emit)
		}
		line := sc.Text()
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}
		if !strings.HasPrefix(line, "data:") {
			continue // event:/id:/retry: lines carry nothing for us
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			dataDone = true
			break
		}
		var chunk wireChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			return a.fail(emit, observe.Wrap(observe.ClassProvider, err,
				"openai-chat: stream chunk is not valid JSON"))
		}
		if chunk.Error != nil {
			return a.fail(emit, chunk.Error.classify())
		}

		// Tool call fragments first (they carry payload).
		for _, choice := range chunk.Choices {
			for _, frag := range choice.Delta.ToolCalls {
				events, err := asm.fragment(frag)
				if err != nil {
					return a.fail(emit, err)
				}
				if err := emitAll(emit, events); err != nil {
					return err
				}
			}
			if t := choice.Delta.text(); t != "" {
				if err := emit(emitText(t)); err != nil {
					return err
				}
			}
			if r := choice.Delta.reasoning(); r != "" {
				if err := emit(emitReasoning(r)); err != nil {
					return err
				}
			}
			if choice.FinishReason != "" && choice.FinishReason != "null" {
				sawFinish = choice.FinishReason
				events, err := asm.finish() // ToolCallEnd events, index order
				if err != nil {
					return a.fail(emit, err)
				}
				if err := emitAll(emit, events); err != nil {
					return err
				}
			}
		}

		if u := chunk.Usage; u != nil && !u.isZero() {
			sawUsage = true
			if err := emit(llm.StreamEvent{Type: llm.EvUsage, Usage: u.normalized()}); err != nil {
				return err
			}
		}
	}
	if err := sc.Err(); err != nil {
		// Network read failed mid-body (the disconnect case). Cancellation
		// masquerades as a read error - check the ctx first.
		if ctx.Err() != nil {
			return a.cancelled(emit)
		}
		return a.fail(emit, llm.ClassifyTransportError(err))
	}
	if ctx.Err() != nil {
		return a.cancelled(emit)
	}

	if !dataDone {
		// Stream ended without [DONE]: a genuine disconnect unless the
		// provider (compat) routinely omits the sentinel after finish_reason.
		if sawFinish != "" && a.loose() {
			// tolerate: finish_reason already seen, treat as complete
		} else {
			return a.fail(emit, llm.WithCode(
				observe.New(observe.ClassNetwork, "openai-chat: stream ended without [DONE]"),
				llm.CodeStreamDisconnect))
		}
	}

	// Terminal mapping with the compat matrix.
	stop, err := a.mapFinish(sawFinish, sawUsage)
	if err != nil {
		return a.fail(emit, err)
	}
	if !sawUsage {
		if !a.loose() && !a.opts.Compat.AllowMissingUsage {
			return a.fail(emit, observe.New(observe.ClassProvider,
				"openai-chat: stream finished without a usage block (set llm.providers.<p>.compat.allow_missing_usage or compat.loose)"))
		}
	}

	if err := emit(llm.StreamEvent{Type: llm.EvStop, Stop: stop}); err != nil {
		return err
	}
	return emit(llm.StreamEvent{Type: llm.EvDone})
}

// mapFinish maps the wire finish_reason to the C6 StopReason under the
// compat rules. Missing usage does not block the stop mapping (the usage
// check happens separately).
func (a *Adapter) mapFinish(finish string, sawUsage bool) (llm.StopReason, error) {
	switch finish {
	case "stop":
		return llm.StopEndTurn, nil
	case "length":
		return llm.StopMaxTokens, nil // must drive failToolCallsFromTruncatedMessage (ticket 10)
	case "tool_calls", "function_call":
		return llm.StopToolUse, nil
	case "content_filter":
		return llm.StopContentFilter, nil
	case "stop_sequence":
		// some OpenAI-compatible servers emit the anthropic vocabulary
		return llm.StopStopSequence, nil
	case "":
		if a.loose() || a.opts.Compat.AllowMissingUsage {
			return llm.StopEndTurn, nil
		}
		return "", observe.New(observe.ClassProvider,
			"openai-chat: stream finished without finish_reason (set compat.loose to tolerate)")
	default:
		if a.loose() {
			return llm.StopEndTurn, nil
		}
		return "", observe.New(observe.ClassProvider,
			fmt.Sprintf("openai-chat: unknown finish_reason %q", finish))
	}
}

func (a *Adapter) loose() bool { return a.opts.Compat.Loose }

// fail emits the terminal failure triple (Error, Stop{error}, Done) and
// returns the same classified error. err may be an error or *observe.Error;
// unclassified errors collapse to ClassInternal (observe.ClassOfOrInternal
// semantics).
func (a *Adapter) fail(emit func(llm.StreamEvent) error, err error) error {
	e, ok := err.(*observe.Error)
	if !ok {
		e = observe.Wrap(observe.ClassOfOrInternal(err), err, "openai-chat failure")
	}
	if emit != nil {
		_ = emit(llm.StreamEvent{Type: llm.EvError, Err: e})
		_ = emit(llm.StreamEvent{Type: llm.EvStop, Stop: llm.StopError})
		_ = emit(llm.StreamEvent{Type: llm.EvDone})
	}
	return e
}

// cancelled emits Stop{cancelled} + Done and returns nil (a cancellation is
// not a provider failure).
func (a *Adapter) cancelled(emit func(llm.StreamEvent) error) error {
	if emit != nil {
		_ = emit(llm.StreamEvent{Type: llm.EvStop, Stop: llm.StopCancelled})
		_ = emit(llm.StreamEvent{Type: llm.EvDone})
	}
	return nil
}

func emitText(t string) llm.StreamEvent {
	return llm.StreamEvent{Type: llm.EvTextDelta, Text: t}
}

func emitReasoning(r string) llm.StreamEvent {
	return llm.StreamEvent{Type: llm.EvReasoningDelta, Text: r}
}

func emitAll(emit func(llm.StreamEvent) error, events []llm.StreamEvent) error {
	for _, ev := range events {
		if err := emit(ev); err != nil {
			return err
		}
	}
	return nil
}
