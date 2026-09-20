package agent

import (
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// The loop's outward event surface. Everything a user can see about a running
// task goes through Sink: the ball/panel/CLI (tickets 12, 28, 33+) forward
// these into the D43 state machine, and the C22 rule "触发必须显式告知用户"
// (never stop silently) is enforced by requiring every brake to publish a
// visible message.

// EventKind enumerates what the loop reports.
type EventKind string

const (
	// EvTextDelta streams assistant text.
	EvTextDelta EventKind = "text_delta"
	// EvReasoningDelta streams provider reasoning text (provider-gated).
	EvReasoningDelta EventKind = "reasoning_delta"
	// EvToolStart announces an executed call.
	EvToolStart EventKind = "tool_start"
	// EvToolEnd announces a finished call with its outcome.
	EvToolEnd EventKind = "tool_end"
	// EvReminder announces an injected C22 gradient reminder.
	EvReminder EventKind = "reminder"
	// EvStuck announces a C22 brake: the task stopped, and why, visibly.
	EvStuck EventKind = "stuck"
	// EvControl announces a D11(1) control word handled without the LLM.
	EvControl EventKind = "control"
	// EvError announces a classified failure (SPEC-05 §3.4: failures are
	// visible and distinguishable, never a silent downgrade).
	EvError EventKind = "error"
	// EvDone closes the task.
	EvDone EventKind = "done"
)

// Event is one loop publication. Fields are meaningful per Kind; Text is
// always safe to show to a user (it never carries tool args beyond a
// truncated preview, and never carries secrets - observe.Error is log-safe).
type Event struct {
	Kind      EventKind
	TaskID    string
	Text      string
	ToolName  string
	CallID    string
	Outcome   string
	Verb      ControlVerb
	Stop      llm.StopReason
	TokensIn  int
	TokensOut int
	// RepeatLevel carries the C22 ladder rung that fired (D43 row #19 guard
	// input Facts.RepeatLevel).
	RepeatLevel int
	Err         *observe.Error
}

// Sink receives loop events. Implementations must not block: the loop calls
// Publish synchronously on the task goroutine.
type Sink interface {
	Publish(Event)
}

// NopSink discards events (used when no UI is wired).
type NopSink struct{}

// Publish implements Sink.
func (NopSink) Publish(Event) {}

// collector is the test/host helper that records events in order.
type collector struct {
	events []Event
}

// RecordSink is an in-memory Sink; safe for the single task goroutine that
// publishes to it (tests use it directly).
type RecordSink struct{ collector }

// Publish implements Sink.
func (s *RecordSink) Publish(e Event) { s.events = append(s.events, e) }

// Events returns the recorded events.
func (s *RecordSink) Events() []Event { return s.events }

// LastOf returns the last event of a kind, or nil.
func (s *RecordSink) LastOf(k EventKind) *Event {
	for i := len(s.events) - 1; i >= 0; i-- {
		if s.events[i].Kind == k {
			return &s.events[i]
		}
	}
	return nil
}

// Of counts events of a kind.
func (s *RecordSink) Of(k EventKind) int {
	n := 0
	for _, e := range s.events {
		if e.Kind == k {
			n++
		}
	}
	return n
}
