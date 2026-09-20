// Package openairesponses is the OpenAI Responses API adapter (C5
// implementor, ticket 11): /v1/responses with reasoning items, function-call
// items and the D37 error classification at the seam.
//
// # Wire mapping
//
// Request (C7 -> wire):
//   - Request.System flattens into the top-level `instructions` string (the
//     protocol has no system message role in `input`).
//   - Messages become `input` items: user/assistant turns are
//     {type:"message",role,content:[{type:"input_text"|"output_text",text}]},
//     an assistant ToolUsePart is a standalone {type:"function_call",
//     call_id,name,arguments} item, and a tool result is
//     {type:"function_call_output",call_id,output}. Images are
//     {type:"input_image",image_url:"data:..."} inside a user message.
//   - Request.ThinkingIntensity maps to reasoning.effort (low|medium|high)
//     with summary="auto" so reasoning summaries stream; "off"/"" omits the
//     reasoning object entirely (provider default).
//   - Request.MaxOutputTokens -> max_output_tokens; StopSequences has NO
//     equivalent in this protocol, so a request that sets them fails with a
//     ClassInternal error unless compat.loose (which drops them).
//   - Caching: the Responses API caches prompt prefixes implicitly, so
//     Request.CacheBreakpoints is reported as having no effect
//     (CacheSupport{Mode: CacheImplicit, Breakpoints: false}) and is ignored
//     rather than smuggled into a provider-private field.
//
// Response (wire -> C6), by SSE `type`:
//   - response.created / response.in_progress -> nothing
//   - response.output_item.added(function_call) -> ToolCallStart (call_id is
//     the tool id the caller will echo back; `id` is the item id deltas use)
//   - response.output_text.delta / response.refusal.delta -> TextDelta
//   - response.reasoning_summary_text.delta /
//     response.reasoning_text.delta -> ReasoningDelta (AC#3)
//   - response.function_call_arguments.delta -> ToolCallArgsDelta
//   - response.output_item.done -> ToolCallEnd (after back-filling arguments
//     or text a non-incremental server never streamed); a reasoning item that
//     ends emits nothing
//   - response.completed / response.incomplete / response.failed ->
//     terminal: usage merged, status + reason mapped (table below); a bare
//     {"type":"response.usage"} frame from a gateway merges usage only
//   - [DONE] -> the terminal sentinel (chat's analogue)
//   - {"type":"error"} -> classified failure
//
// # Stop-reason mapping table (documented per ticket 11)
//
//	wire status + detail                 C6 StopReason      notes
//	completed, function_call emitted     tool_use           this protocol has no stop_reason field
//	completed                            end_turn           -
//	completed, refusal item              content_filter     -
//	incomplete, max_output_tokens        max_tokens         drives failToolCallsFromTruncatedMessage
//	incomplete, content_filter           content_filter     -
//	incomplete, other reason             max_tokens         token/limit truncation is the only other cause
//	failed                               (failure path)     Error{class} + Stop{error}
//	"" (no terminal status, [DONE] seen) strict: provider error; loose/allow_missing_usage: end_turn
package openairesponses

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Protocol is the config protocol enum this adapter registers.
const Protocol = "openai-responses"

func init() {
	llm.RegisterProtocol(Protocol, func(ep llm.EndpointOptions) (llm.LlmProvider, error) {
		return New(ep), nil
	})
}

// Adapter implements llm.LlmProvider for the OpenAI Responses API.
type Adapter struct {
	opts llm.EndpointOptions
	http *http.Client
}

// New builds an adapter; a nil HTTPClient gets the proxy-aware default.
func New(opts llm.EndpointOptions) *Adapter {
	c := opts.HTTPClient
	if c == nil {
		c = llm.NewDefaultHTTPClient("", "")
	}
	return &Adapter{opts: opts, http: c}
}

// Info reports implicit prefix caching (no explicit breakpoints in this
// protocol) - see the package doc.
func (a *Adapter) Info() llm.ProviderInfo {
	return llm.ProviderInfo{
		Protocol:         Protocol,
		Provider:         a.opts.Provider,
		Model:            a.opts.Model,
		Cache:            llm.CacheSupport{Mode: llm.CacheImplicit, Breakpoints: false},
		MaxContextWindow: a.opts.ContextWindow,
	}
}

// Stream implements llm.LlmProvider (see package llm events.go).
func (a *Adapter) Stream(ctx context.Context, req *llm.Request, emit func(llm.StreamEvent) error) error {
	if err := req.Validate(); err != nil {
		return a.fail(emit, observe.Wrap(observe.ClassInternal, err, "invalid normalized request"))
	}
	payload, err := a.buildWire(req)
	if err != nil {
		return a.fail(emit, err)
	}

	url := strings.TrimRight(a.opts.BaseURL, "/") + "/responses"
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

// consume parses the Responses SSE dialect into C6 events.
func (a *Adapter) consume(ctx context.Context, body io.Reader, emit func(llm.StreamEvent) error) error {
	st := newStreamState()
	sc := bufio.NewScanner(body)
	sc.Buffer(make([]byte, 0, 64<<10), 4<<20)

	for sc.Scan() {
		if ctx.Err() != nil {
			return a.cancelled(emit)
		}
		line := strings.TrimRight(sc.Text(), "\r")
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}
		if !strings.HasPrefix(line, "data:") {
			continue // event:/id:/retry: carry nothing this protocol needs
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			st.sawDone = true
			break
		}
		stop, err := st.handle(data, emit)
		if err != nil {
			return a.fail(emit, err)
		}
		if stop {
			break
		}
	}
	if err := sc.Err(); err != nil {
		if ctx.Err() != nil {
			return a.cancelled(emit)
		}
		return a.fail(emit, llm.ClassifyTransportError(err))
	}
	if ctx.Err() != nil {
		return a.cancelled(emit)
	}

	if !st.sawDone && !st.sawTerminal {
		// No sentinel and no terminal status frame: mid-stream disconnect.
		if a.loose() && st.status != "" {
			// tolerated: a terminal status arrived without the sentinel
		} else {
			return a.fail(emit, llm.WithCode(
				observe.New(observe.ClassNetwork, "openai-responses: stream ended without [DONE] or a terminal status"),
				llm.CodeStreamDisconnect))
		}
	}

	stop, err := a.mapTerminal(st)
	if err != nil {
		return a.fail(emit, err)
	}
	if !st.sawUsage {
		if !a.loose() && !a.opts.Compat.AllowMissingUsage {
			return a.fail(emit, observe.New(observe.ClassProvider,
				"openai-responses: stream carried no usage block (set compat.allow_missing_usage or compat.loose)"))
		}
	}
	if err := emit(llm.StreamEvent{Type: llm.EvStop, Stop: stop}); err != nil {
		return err
	}
	return emit(llm.StreamEvent{Type: llm.EvDone})
}

// mapTerminal derives the C6 stop reason (table in the package doc): this
// protocol reports a status plus per-item types instead of a finish_reason,
// so tool_use is inferred from emitted function_call items.
func (a *Adapter) mapTerminal(st *streamState) (llm.StopReason, error) {
	switch st.status {
	case "completed":
		if st.sawRefusal {
			return llm.StopContentFilter, nil
		}
		if st.calledFunction {
			return llm.StopToolUse, nil
		}
		return llm.StopEndTurn, nil
	case "incomplete":
		switch st.incompleteReason {
		case "content_filter":
			return llm.StopContentFilter, nil
		case "max_output_tokens", "":
			return llm.StopMaxTokens, nil
		default:
			// Any other incomplete reason is a limit truncation of the turn.
			return llm.StopMaxTokens, nil
		}
	case "failed", "cancelled":
		// A failed terminal is the failure path; reach here only when the
		// error payload was empty, so report it as a provider failure.
		return "", observe.New(observe.ClassProvider,
			"openai-responses: response terminal status "+st.status)
	case "":
		if a.loose() || a.opts.Compat.AllowMissingUsage {
			if st.calledFunction {
				return llm.StopToolUse, nil
			}
			return llm.StopEndTurn, nil
		}
		return "", observe.New(observe.ClassProvider,
			"openai-responses: stream finished without a terminal status (set compat.loose to tolerate)")
	default:
		if a.loose() {
			if st.calledFunction {
				return llm.StopToolUse, nil
			}
			return llm.StopEndTurn, nil
		}
		return "", observe.New(observe.ClassProvider,
			fmt.Sprintf("openai-responses: unknown response status %q", st.status))
	}
}

func (a *Adapter) loose() bool { return a.opts.Compat.Loose }

// fail emits the terminal failure triple and returns the classified error.
func (a *Adapter) fail(emit func(llm.StreamEvent) error, err error) error {
	e, ok := err.(*observe.Error)
	if !ok {
		e = observe.Wrap(observe.ClassOfOrInternal(err), err, "openai-responses failure")
	}
	if emit != nil {
		_ = emit(llm.StreamEvent{Type: llm.EvError, Err: e})
		_ = emit(llm.StreamEvent{Type: llm.EvStop, Stop: llm.StopError})
		_ = emit(llm.StreamEvent{Type: llm.EvDone})
	}
	return e
}

// cancelled emits Stop{cancelled} + Done and returns nil.
func (a *Adapter) cancelled(emit func(llm.StreamEvent) error) error {
	if emit != nil {
		_ = emit(llm.StreamEvent{Type: llm.EvStop, Stop: llm.StopCancelled})
		_ = emit(llm.StreamEvent{Type: llm.EvDone})
	}
	return nil
}
