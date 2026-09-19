package llm

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/observe"
)

// C6 normalized stream events (SPEC-05 sec 3.2, contract-level). The event
// set is EXACT - adapters may not invent additional kinds, and the agent core
// must not see protocol shapes. Ordering contract of a Stream call:
//
//	data events (TextDelta / ReasoningDelta / ToolCall* / Usage) ...
//	  then exactly one terminal pair:
//	    success:  Stop{reason != error, != cancelled}
//	    failure:  Error{...} then Stop{reason == error}
//	    cancel:   Stop{reason == cancelled}          (no Error event)
//	  then Done (always last, exactly once).
//
// Stream() returns nil on success and cancellation, and the same classified
// *observe.Error it emitted on failure (callers may rely on either signal).
// Hard requirement from SPEC-05 sec 3.2: Stop.reason == max_tokens must make
// the agent core fail open tool calls (failToolCallsFromTruncatedMessage,
// ticket 10); the seam guarantees the reason is emitted verbatim.

// StreamEventType discriminates the C6 union.
type StreamEventType uint8

const (
	EvTextDelta StreamEventType = iota + 1
	EvReasoningDelta
	EvToolCallStart
	EvToolCallArgsDelta
	EvToolCallEnd
	EvUsage
	EvStop
	EvError
	EvDone
)

func (t StreamEventType) String() string {
	switch t {
	case EvTextDelta:
		return "text_delta"
	case EvReasoningDelta:
		return "reasoning_delta"
	case EvToolCallStart:
		return "tool_call_start"
	case EvToolCallArgsDelta:
		return "tool_call_args_delta"
	case EvToolCallEnd:
		return "tool_call_end"
	case EvUsage:
		return "usage"
	case EvStop:
		return "stop"
	case EvError:
		return "error"
	case EvDone:
		return "done"
	default:
		return fmt.Sprintf("stream_event(%d)", uint8(t))
	}
}

// StopReason is the normalized terminal reason (C6). The set is exact.
type StopReason string

const (
	StopEndTurn       StopReason = "end_turn"
	StopMaxTokens     StopReason = "max_tokens"
	StopToolUse       StopReason = "tool_use"
	StopStopSequence  StopReason = "stop_sequence"
	StopContentFilter StopReason = "content_filter"
	StopCancelled     StopReason = "cancelled"
	StopError         StopReason = "error"
)

func (r StopReason) Valid() bool {
	switch r {
	case StopEndTurn, StopMaxTokens, StopToolUse, StopStopSequence,
		StopContentFilter, StopCancelled, StopError:
		return true
	}
	return false
}

// Usage is the token accounting event (C23 inputs). CachedTokens counts
// prompt tokens served from the provider cache.
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	CachedTokens int `json:"cached_tokens"`
}

// IsZero reports whether no field carries a value.
func (u Usage) IsZero() bool { return u == Usage{} }

// Add merges o into u. All three counters are monotonic within one response
// (OpenAI sends one final block; Anthropic-style streams report cumulative
// output in every delta), so the max-merge is the safe aggregation for both
// and collapses duplicate reports to their freshest value.
func (u Usage) Add(o Usage) Usage {
	if o.InputTokens > u.InputTokens {
		u.InputTokens = o.InputTokens
	}
	if o.OutputTokens > u.OutputTokens {
		u.OutputTokens = o.OutputTokens
	}
	if o.CachedTokens > u.CachedTokens {
		u.CachedTokens = o.CachedTokens
	}
	return u
}

// StreamEvent is one C6 event. Fields not relevant to Type are zero. The
// struct (rather than an interface union) keeps event streams comparable with
// reflect.DeepEqual, which the golden tests rely on.
type StreamEvent struct {
	Type StreamEventType

	Text       string // TextDelta, ReasoningDelta
	ToolCallID string // ToolCallStart, ToolCallArgsDelta, ToolCallEnd
	ToolName   string // ToolCallStart
	ArgsDelta  string // ToolCallArgsDelta (raw JSON fragment)
	Usage      Usage  // EvUsage
	Stop       StopReason
	Err        *observe.Error // EvError (class is a D37 enum, see observe.ErrorClass)
}

// IsData reports whether the event carries turn payload (as opposed to the
// terminal events). Retry wrappers treat data events as "the attempt was
// committed to the consumer": after the first data event a failing stream is
// never retried (mid-stream reconnect would duplicate content).
func (e StreamEvent) IsData() bool {
	switch e.Type {
	case EvTextDelta, EvReasoningDelta, EvToolCallStart,
		EvToolCallArgsDelta, EvToolCallEnd, EvUsage:
		return true
	}
	return false
}

// String renders the event for logs. Error details are included (they are
// already log-safe: observe.Error never carries secrets), usage and deltas
// are truncated.
func (e StreamEvent) String() string {
	var b strings.Builder
	b.WriteString(e.Type.String())
	switch e.Type {
	case EvTextDelta, EvReasoningDelta:
		fmt.Fprintf(&b, " %q", truncateRunes(e.Text, 80))
	case EvToolCallStart:
		fmt.Fprintf(&b, " id=%s name=%s", e.ToolCallID, e.ToolName)
	case EvToolCallArgsDelta:
		fmt.Fprintf(&b, " id=%s delta=%q", e.ToolCallID, truncateRunes(e.ArgsDelta, 80))
	case EvToolCallEnd:
		fmt.Fprintf(&b, " id=%s", e.ToolCallID)
	case EvUsage:
		fmt.Fprintf(&b, " %+v", e.Usage)
	case EvStop:
		b.WriteString(" " + string(e.Stop))
	case EvError:
		if e.Err != nil {
			fmt.Fprintf(&b, " %s", e.Err.Error())
		}
	}
	return b.String()
}

// ToolCall is one assembled tool invocation (collector output).
type ToolCall struct {
	ID       string
	Name     string
	Args     json.RawMessage // assembled argument JSON; nil when incomplete
	Complete bool            // true when ToolCallEnd was seen
}

// TurnResult is the assembled outcome of one streamed turn (C7 turn shape).
// Incomplete is true when the stream ended in error or cancellation after
// payload was already delivered: the partial content is retained and marked,
// never silently discarded (SPEC-05 sec 3.4). Tool calls left open by such a
// stream are reported with Complete=false - the agent core (ticket 10) turns
// them into failed tool results; the seam never fabricates ToolCallEnd.
type TurnResult struct {
	Text       string
	Reasoning  string
	ToolCalls  []ToolCall
	Usage      Usage
	Stop       StopReason
	Err        *observe.Error
	Incomplete bool
}

// TurnCollector assembles a C6 event stream into a TurnResult. It is shared
// infrastructure for the agent loop (ticket 10) and the golden tests; the
// event contract above is what it relies on, nothing else.
type TurnCollector struct {
	res        TurnResult
	text       strings.Builder
	reasoning  strings.Builder
	byID       map[string]int // tool call id -> index in res.ToolCalls
	sawPayload bool
	finished   bool
}

// NewTurnCollector creates a collector.
func NewTurnCollector() *TurnCollector {
	return &TurnCollector{byID: map[string]int{}}
}

// Observe feeds one event. Events after Done are ignored (idempotent close).
func (c *TurnCollector) Observe(ev StreamEvent) {
	if c.finished {
		return
	}
	switch ev.Type {
	case EvTextDelta:
		c.text.WriteString(ev.Text)
		c.sawPayload = true
	case EvReasoningDelta:
		c.reasoning.WriteString(ev.Text)
		c.sawPayload = true
	case EvToolCallStart:
		c.res.ToolCalls = append(c.res.ToolCalls, ToolCall{ID: ev.ToolCallID, Name: ev.ToolName})
		c.byID[ev.ToolCallID] = len(c.res.ToolCalls) - 1
		c.sawPayload = true
	case EvToolCallArgsDelta:
		i, ok := c.byID[ev.ToolCallID]
		if !ok {
			// Adapter contract violation (ArgsDelta before Start): record as an
			// orphan fragment so it is visible instead of silently dropped.
			c.res.ToolCalls = append(c.res.ToolCalls, ToolCall{ID: ev.ToolCallID})
			c.byID[ev.ToolCallID] = len(c.res.ToolCalls) - 1
			i = len(c.res.ToolCalls) - 1
		}
		c.res.ToolCalls[i].Args = append(c.res.ToolCalls[i].Args, ev.ArgsDelta...)
		c.sawPayload = true
	case EvToolCallEnd:
		if i, ok := c.byID[ev.ToolCallID]; ok {
			c.res.ToolCalls[i].Complete = true
		}
	case EvUsage:
		c.res.Usage = c.res.Usage.Add(ev.Usage)
	case EvStop:
		c.res.Stop = ev.Stop
		if ev.Stop == StopError || ev.Stop == StopCancelled {
			c.res.Incomplete = c.sawPayload
		}
	case EvError:
		c.res.Err = ev.Err
		c.res.Incomplete = c.res.Incomplete || c.sawPayload
	case EvDone:
		c.finished = true
		c.res.Text = c.text.String()
		c.res.Reasoning = c.reasoning.String()
	}
}

// SawPayload reports whether any data event has been observed.
func (c *TurnCollector) SawPayload() bool { return c.sawPayload }

// Result returns the assembled turn (valid, and frozen, after Done).
func (c *TurnCollector) Result() TurnResult { return c.res }

// HasOpenToolCalls reports whether any observed tool call never ended. The
// agent core maps these to failed tool results (SPEC-05 sec 3.4).
func (c *TurnCollector) HasOpenToolCalls() bool {
	for _, tc := range c.res.ToolCalls {
		if !tc.Complete {
			return true
		}
	}
	return false
}
