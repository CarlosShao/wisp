// Package anthropic is the Anthropic Messages adapter (C5 implementor,
// ticket 11): /v1/messages with native tool use, extended thinking, explicit
// prompt-caching breakpoints and the D37 error classification at the seam.
//
// # Wire mapping
//
// Request (C7 -> wire):
//   - Request.System becomes the top-level `system` text-block array. A
//     CacheBreakpoints entry of -1 puts `cache_control:{"type":"ephemeral"}`
//     on the LAST system block (prefix 1/5/6/2 of SPEC-05 sec 4.1); entries
//     >= 0 put it on the last text block of Messages[i]. Nothing else ever
//     gets a breakpoint, so the volatile suffix is never cached. Anthropic
//     accepts at most 4 breakpoints: more than 4 declared is a ClassInternal
//     error in strict mode (our bug) and the first 4 in compat.loose.
//   - assistant ToolUsePart -> content block {type:"tool_use",id,name,input};
//     tool messages -> a USER message carrying {type:"tool_result",
//     tool_use_id,content,is_error} (Anthropic has no tool role).
//   - Request.ThinkingIntensity maps to the `thinking` parameter:
//     off/"" -> omitted (provider default), low -> budget_tokens 1024,
//     medium -> 4096, high -> 8192, each clamped to max_tokens-1 (Anthropic
//     requires max_tokens > budget). Because Anthropic rejects
//     temperature != 1 while thinking is enabled, Temperature is dropped in
//     that case. These budgets are adapter-local numbers, not contract.
//   - Request.MaxOutputTokens: `max_tokens` is REQUIRED by the protocol, so a
//     request that omits it gets defaultMaxTokens (4096).
//
// Response (wire -> C6), by SSE `type`:
//   - message_start      -> Usage{input, cached} event (cached from
//     cache_read_input_tokens)
//   - content_block_start(text|tool_use|thinking) -> ToolCallStart for
//     tool_use blocks (a non-empty `input` at start is replayed as one args
//     delta); other block kinds emit nothing
//   - content_block_delta: text_delta -> TextDelta, thinking_delta ->
//     ReasoningDelta, signature_delta -> dropped (no C6 event exists),
//     input_json_delta -> ToolCallArgsDelta for the block's tool id
//   - content_block_stop -> ToolCallEnd (tool_use blocks only; a tool_use
//     block that never stops stays open, which is what the agent core's
//     failToolCallsFromTruncatedMessage path needs)
//   - message_delta      -> records stop_reason; its cumulative usage emits a
//     second Usage event (the collector max-merges)
//   - message_stop       -> the terminal sentinel (the [DONE] analogue: a
//     stream ending without it is a mid-stream disconnect)
//   - error              -> classified failure (both the `error` SSE event and
//     an {"type":"error"} data payload)
//   - ping               -> ignored
//
// # Stop-reason mapping table (documented per ticket 11)
//
//	wire stop_reason            C6 StopReason      notes
//	end_turn                    end_turn           -
//	max_tokens                  max_tokens         drives failToolCallsFromTruncatedMessage
//	stop_sequence               stop_sequence      -
//	tool_use                    tool_use           -
//	refusal                     content_filter     Anthropic refuses instead of filtering
//	pause_turn                  end_turn           server tool-pause; turn is complete
//	(max)                       -
//	"" (absent)                 strict: provider error; loose/allow_missing_usage: end_turn
//	unknown value               strict: provider error; loose: end_turn
//
// # Usage mapping
//
// Anthropic's `input_tokens` counts NON-cached prompt tokens only (cache
// reads and cache writes are reported separately), unlike OpenAI's
// `prompt_tokens`, which includes cached tokens. The seam maps verbatim
// (InputTokens = input_tokens, CachedTokens = cache_read_input_tokens) and
// does NOT fold cache reads into InputTokens: inventing a sum would change
// what C6 means. Consumers that price a turn (C23, ticket 44) must add
// CachedTokens to InputTokens for this protocol; that asymmetry is recorded
// as an open question in ticket 11 rather than resolved by guessing.
package anthropic

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
const Protocol = "anthropic"

// AnthropicVersion is the required anthropic-version header value.
const AnthropicVersion = "2023-06-01"

// defaultMaxTokens applies when the caller sets none (the protocol requires
// the field).
const defaultMaxTokens = 4096

// maxCacheControl is how many cache_control breakpoints one request may carry.
const maxCacheControl = 4

func init() {
	llm.RegisterProtocol(Protocol, func(ep llm.EndpointOptions) (llm.LlmProvider, error) {
		return New(ep), nil
	})
}

// Adapter implements llm.LlmProvider for the Anthropic Messages API.
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

// Info reports the explicit prompt-caching capability (C5 cache-breakpoint
// capability, SPEC-05 sec 4.1) - this is the one protocol in the seam that
// honors Request.CacheBreakpoints.
func (a *Adapter) Info() llm.ProviderInfo {
	return llm.ProviderInfo{
		Protocol: Protocol,
		Provider: a.opts.Provider,
		Model:    a.opts.Model,
		Cache:    llm.CacheSupport{Mode: llm.CacheExplicit, Breakpoints: true},
		// The 1-hour cache tier and per-tier TTLs are C5's business later;
		// the seam only exposes "explicit breakpoints honored".
		MaxContextWindow: a.opts.ContextWindow,
	}
}

// Stream implements llm.LlmProvider (see package llm events.go for the
// contract).
func (a *Adapter) Stream(ctx context.Context, req *llm.Request, emit func(llm.StreamEvent) error) error {
	if err := req.Validate(); err != nil {
		return a.fail(emit, observe.Wrap(observe.ClassInternal, err, "invalid normalized request"))
	}
	payload, err := a.buildWire(req)
	if err != nil {
		return a.fail(emit, err)
	}

	url := strings.TrimRight(a.opts.BaseURL, "/") + "/messages"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return a.fail(emit, observe.Wrap(observe.ClassInternal, err, "build http request"))
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")
	httpReq.Header.Set("anthropic-version", AnthropicVersion)
	if a.opts.APIKey != "" {
		// Anthropic uses x-api-key, not a Bearer scheme.
		httpReq.Header.Set("x-api-key", a.opts.APIKey)
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
			detail.Code, detail.Type, detail.Message = e.Type, e.Type, e.Message
		} else {
			detail.Message = strings.TrimSpace(string(body))
		}
		return a.fail(emit, llm.NewHTTPError(detail,
			llm.ParseRetryAfter(resp.Header.Get("Retry-After"), time.Now())))
	}
	defer resp.Body.Close()

	return a.consume(ctx, resp.Body, emit)
}

// consume parses Anthropic SSE (event:/data: pairs) into C6 events.
func (a *Adapter) consume(ctx context.Context, body io.Reader, emit func(llm.StreamEvent) error) error {
	st := newStreamState()
	sc := bufio.NewScanner(body)
	sc.Buffer(make([]byte, 0, 64<<10), 4<<20)

	var eventName string
	var data []byte

	for sc.Scan() {
		if ctx.Err() != nil {
			return a.cancelled(emit)
		}
		line := strings.TrimRight(sc.Text(), "\r")
		switch {
		case line == "":
			// Blank line dispatches one SSE event.
			if data == nil {
				eventName = ""
				continue
			}
			stop, err := st.handle(ctx, eventName, data, emit)
			data = nil
			eventName = ""
			if err != nil {
				if ferr, ok := err.(*observe.Error); ok && ferr.Class == observe.ClassCancelled {
					return a.cancelled(emit)
				}
				return a.fail(emit, err)
			}
			if stop {
				// message_stop with no dispatch left: terminal handled below.
			}
		case strings.HasPrefix(line, ":"):
			// SSE comment (keep-alive).
		case strings.HasPrefix(line, "event:"):
			eventName = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
		case strings.HasPrefix(line, "data:"):
			chunk := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if chunk == "" {
				continue
			}
			if data == nil {
				data = []byte(chunk)
			} else {
				data = append(data, []byte("\n"+chunk)...)
			}
		default:
			// id:/retry: and unknown fields carry nothing for this dialect.
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

	// Terminal sentinel: a Messages stream that ends without message_stop is
	// a mid-stream disconnect (the [DONE] analogue), unless compat.loose and
	// the terminal delta already arrived.
	if !st.sawMessageStop {
		if st.stopReason != "" && a.loose() {
			// tolerated: stop_reason seen, sentinel omitted by this provider
		} else {
			return a.fail(emit, llm.WithCode(
				observe.New(observe.ClassNetwork, "anthropic: stream ended without message_stop"),
				llm.CodeStreamDisconnect))
		}
	}

	stop, err := a.mapStop(st.stopReason)
	if err != nil {
		return a.fail(emit, err)
	}
	if !st.sawUsage {
		if !a.loose() && !a.opts.Compat.AllowMissingUsage {
			return a.fail(emit, observe.New(observe.ClassProvider,
				"anthropic: stream carried no usage block (set compat.allow_missing_usage or compat.loose)"))
		}
	}
	if st.openBlocks() {
		// An unclosed content block is a truncated turn: the agent core
		// fails the open tool calls, the seam never invents ToolCallEnd.
		// Nothing to emit here; the collector reports it via HasOpenToolCalls.
		_ = st
	}

	if err := emit(llm.StreamEvent{Type: llm.EvStop, Stop: stop}); err != nil {
		return err
	}
	return emit(llm.StreamEvent{Type: llm.EvDone})
}

// mapStop translates the wire stop_reason (table in the package doc).
func (a *Adapter) mapStop(reason string) (llm.StopReason, error) {
	switch reason {
	case "end_turn":
		return llm.StopEndTurn, nil
	case "max_tokens":
		return llm.StopMaxTokens, nil
	case "stop_sequence":
		return llm.StopStopSequence, nil
	case "tool_use":
		return llm.StopToolUse, nil
	case "refusal":
		return llm.StopContentFilter, nil
	case "pause_turn":
		// Server-side tool pause: the visible turn ended, the loop re-sends.
		return llm.StopEndTurn, nil
	case "":
		if a.loose() || a.opts.Compat.AllowMissingUsage {
			return llm.StopEndTurn, nil
		}
		return "", observe.New(observe.ClassProvider,
			"anthropic: stream finished without a stop_reason (set compat.loose to tolerate)")
	default:
		if a.loose() {
			return llm.StopEndTurn, nil
		}
		return "", observe.New(observe.ClassProvider,
			fmt.Sprintf("anthropic: unknown stop_reason %q", reason))
	}
}

func (a *Adapter) loose() bool { return a.opts.Compat.Loose }

// fail emits the terminal failure triple and returns the classified error.
func (a *Adapter) fail(emit func(llm.StreamEvent) error, err error) error {
	e, ok := err.(*observe.Error)
	if !ok {
		e = observe.Wrap(observe.ClassOfOrInternal(err), err, "anthropic failure")
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
