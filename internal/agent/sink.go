package agent

import (
	"sync"

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
	// EvWarn announces a non-fatal warning the user has to see: today the C23 money
	// budget crossing [cost] alert_threshold (ticket 44 AC#3, "80% warn event"), where
	// the task does start and nothing else announces it. It is its own kind rather than
	// a reused EvError because an error says a task failed and a warning says one went
	// through against a budget the user should know is close (票 40 renders the two
	// apart).
	EvWarn EventKind = "warn"
	// EvUsage reports the task's cumulative spend at the moment one turn's usage
	// landed in the accumulator (publishUsage). The values are the Cost
	// accumulator's own totals - the same numbers Result books - not a
	// per-response echo: llm.EvUsage is per-response, this is the task-so-far
	// truth the panel's working-state rows render (ticket 145 rows 2 and 12).
	EvUsage EventKind = "usage"
	// EvDone closes the task.
	EvDone EventKind = "done"
)

// Event is one loop publication. Fields are meaningful per Kind; Text is
// always safe to show to a user (it never carries tool args beyond a
// truncated preview, and never carries secrets - observe.Error is log-safe).
type Event struct {
	Kind     EventKind
	TaskID   string
	Text     string
	ToolName string
	CallID   string
	Outcome  string
	Verb     ControlVerb
	// Stop is the terminal word the loop publishes on its terminal events:
	// EvDone carries it on every ending (ticket 145 row 13, census §2.14's
	// OPTION A ruling), and the terminal EvError publications carry it where
	// the provider's own stop reason exists. Empty means the ending was the
	// loop's own (a C22 brake, a refusal) with no provider word - absent is
	// the truth, never filled from Status.
	Stop      llm.StopReason
	TokensIn  int
	TokensOut int
	// RepeatLevel carries the C22 ladder rung that fired (D43 row #19 guard
	// input Facts.RepeatLevel).
	RepeatLevel int
	Err         *observe.Error
	// Usage is the Cost accumulator's cumulative token counters (EvUsage):
	// in, out and cached - llm.Usage's own three, verbatim. EvDone's
	// TokensIn/TokensOut stay the terminal echo they always were; this field
	// is the live one.
	Usage llm.Usage
	// CostMicros and Currency travel as a PAIR or not at all (census
	// 145-snapshot-field-census-r1.md §2.13's unit rule: the money number is
	// micro-denominated while Currency is the config's own label, so one
	// without the other manufactures a fake unit). Carried on EvUsage.
	CostMicros int64
	Currency   string
	// CorrelationID is the C18 routing key the call's forensics row is booked
	// under, minted by the same callCorr the dispatch used (EvToolEnd). It is
	// the join key between the event face (call id) and the store face (row).
	CorrelationID string
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

// RecordSink is an in-memory Sink that records events in order.
//
// It is safe for concurrent publishers, which the loop can genuinely produce:
// one task's goroutine and a second (control-word) task running on the caller's
// goroutine both publish to the sink the harness shares, so an unsynchronised
// append here is a data race under -race, not merely a theoretical one.
type RecordSink struct {
	mu     sync.Mutex
	events []Event
}

// Publish implements Sink.
func (s *RecordSink) Publish(e Event) {
	s.mu.Lock()
	s.events = append(s.events, e)
	s.mu.Unlock()
}

// Events returns a snapshot copy of the recorded events.
func (s *RecordSink) Events() []Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Event, len(s.events))
	copy(out, s.events)
	return out
}

// LastOf returns the last event of a kind, or nil.
func (s *RecordSink) LastOf(k EventKind) *Event {
	evs := s.Events()
	for i := len(evs) - 1; i >= 0; i-- {
		if evs[i].Kind == k {
			return &evs[i]
		}
	}
	return nil
}

// Of counts events of a kind.
func (s *RecordSink) Of(k EventKind) int {
	n := 0
	for _, e := range s.Events() {
		if e.Kind == k {
			n++
		}
	}
	return n
}
