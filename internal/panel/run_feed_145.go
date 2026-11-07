package panel

import (
	"fmt"
	"sync"
)

// Ticket 145's working-state carrier: the snapshot's top-level "run" key - the
// dimensions the fourteen-row working-state table (docs/PLAN.md's UI-state
// table) asks the panel to render while a task is running.
//
// READER-GATED, THE 213/187/65 PATTERN, NO SECOND PARADIGM. A pump assembled
// without the reader sends no "run" key at all - the shape the four-key byte
// nails pin (pump_test.go:123, :291, subagent_roster_197_test.go:214) - and a
// pump that HAS one always sends the key, with an empty feed rendered as
// "tools":[] rather than a vanished key. Nothing here opens a door and no C17
// method arrives with it.
//
// THE VALUES ARE THE LOOP'S OWN, NOT A COPY. RunFeed is the holder the
// assembly root's event sink feeds; every field below is written by a value
// an agent.Event already carries (internal/agent/loop.go's publish sites),
// carried through verbatim. Nothing is derived, defaulted or invented here:
// a dimension the events do not carry is ABSENT from this section, and the
// absent ones are named in the field comments plus
// docs/evidence/s1/145-snapshot-fields.md so the next leg does not mistake
// an omission for an oversight.
//
// Verified event sources this carrier rides (readings 2026-11-07, anchors in
// the evidence file):
//   - tool rows: EvToolEnd only (loop.go's executeCalls / truncated /
//     failOpenCalls publish sites). EvToolStart has a kind, a console handler
//     and NO publisher in production code - the "running" chip state has no
//     writer, so a row appears here only when its call has ENDED, carrying the
//     outcome it actually ended in. The outcome vocabulary is internal/agent's
//     own journal vocabulary (success | error | cancelled | truncated),
//     carried through untouched like the D43 names in the tasks section; this
//     package spells none of the words.
//   - reasoning: EvReasoningDelta text deltas, accumulated under a per-rune
//     bound with an honest elided count (same honesty rule as StreamLog's
//     markers: the bound is counted, never silent).
//   - error: EvError events that carry the loop's classified error - the
//     class and the log-safe detail travel verbatim. The human sentence for
//     an error state has no source in the tree and is NOT invented here.
//   - stuck: EvStuck - the loop's own visible brake sentence plus the C22
//     ladder rung that fired. The three buttons the table draws next to it
//     are inbound C17 surface (Q-76) and arrive with nothing here.
//
// Registered as NOT landable today (no event source, so landing any of them
// would be a constant dressed as data - the exact shape ticket 145's AC#1
// forbids): the "running" tool state (EvToolStart is never published), tool
// args digest / risk badge / elapsed (the events carry none of them), usage
// and cost (the loop-local accumulator dies inside Execute; EvDone's final
// in/out echo has no cached/money half and no status word), the terminal
// status word and cancel reason (EvDone.Text is res.Text || the terminal
// sentence, so a cancelled run's event text is its partial REPLY, not a
// reason - carrying it as one would fabricate the row-13 state), phase /
// elapsed timers, the human error sentence, and the L2 window's remainingMs.
// pending[].position and context.occupancy are ticket 167's territory (the
// two legs must not each add their own pending[] keys); composer.git's
// branches half is ticket 186's; the "which screen" view dimension is held
// by AC#3's verdict (the Go side does not know the screen names, and the
// panel-to-Go half reopens the C17 method whitelist).

// ToolRunRow is one ENDED tool call. Rows are written by EvToolEnd only, in
// end order; the fields are the event's own CallID, ToolName, Outcome and
// result text, verbatim.
type ToolRunRow struct {
	CallID     string `json:"callId"`
	Name       string `json:"name"`
	Outcome    string `json:"outcome"`
	ResultText string `json:"resultText,omitempty"`
}

// RunErrorRow is the last classified error the loop published. Class and
// Detail are observe.Error's own fields as the event carried them; Text is
// the event's message string. A notice-shaped EvError (no classified error
// attached) does not write this row - inventing a class for it would be the
// constant-as-data shape again.
type RunErrorRow struct {
	Class  string `json:"class"`
	Detail string `json:"detail,omitempty"`
	Text   string `json:"text,omitempty"`
}

// RunStuckRow is the loop's own stuck sentence and the C22 rung that fired.
// The gradient-brake buttons are C17 inbound surface and are not this key's.
type RunStuckRow struct {
	Text        string `json:"text"`
	RepeatLevel int    `json:"repeatLevel"`
}

// RunSection is the snapshot's top-level "run" key.
type RunSection struct {
	// Tools lists the tool calls that have ENDED in feed order, oldest kept
	// first. Empty (never null) when nothing has ended yet.
	Tools []ToolRunRow `json:"tools"`
	// ToolsDropped counts the oldest rows the cap evicted. A bound without
	// its count would turn "the oldest rows fell off" into silence - the
	// same honesty rule StreamLog's dropped ledger states.
	ToolsDropped int `json:"toolsDropped,omitempty"`
	// Reasoning is the reasoning text streamed so far for the live turn,
	// under the bound below. Empty when none streamed yet; a bounded feed
	// carries its elided count inside the marker, never silently.
	Reasoning string `json:"reasoning,omitempty"`
	// Error is the last classified error event, or nil while none arrived.
	Error *RunErrorRow `json:"error,omitempty"`
	// Stuck is the loop's last stuck sentence, or nil while the task is not
	// stuck.
	Stuck *RunStuckRow `json:"stuck,omitempty"`
}

const (
	// runToolRowsKept is how many ended tool rows the feed retains at full
	// fidelity. A long task can end hundreds of calls; the panel renders the
	// recent ones and the count above keeps the eviction honest (64 rows of
	// ordinary ids and outcomes is a bounded price, like StreamLog's caps).
	runToolRowsKept = 64

	// runReasoningKeepRunes is the tail window reasoning accumulates under.
	// Runes, not bytes - a byte cut would split the Chinese this project
	// streams and hand the panel mojibake. 4096 matches the per-row volume
	// bound the results log pays, so neither region outbounds the other for
	// the same flood.
	runReasoningKeepRunes = 4096
)

// RunFeed is the working-state holder the assembly root feeds and the pump's
// reader reads. It is a StreamLog-shaped dependency of the snapshot pump: the
// console sink (cmd/wisp's event switch, which already receives every one of
// these events) feeds it, run.go's assembly root hands Section to PumpSources,
// and the mutex covers the two goroutines that genuinely meet here - the task
// goroutine publishing events and whatever drives the pump.
//
// The feed methods take plain values, not agent.Event, on purpose: this
// package does not import internal/agent, and the narrow seams keep the
// carrier honest about exactly what it was given.
type RunFeed struct {
	mu            sync.Mutex
	tools         []ToolRunRow
	toolsDropped  int
	reasoningTail []rune
	reasoningSeen int
	lastError     *RunErrorRow
	stuck         *RunStuckRow
}

// NewRunFeed builds an empty feed.
func NewRunFeed() *RunFeed { return &RunFeed{} }

// ToolEnded records one ended tool call. Oldest rows past the cap are evicted
// and counted, never folded into a neighbour.
func (f *RunFeed) ToolEnded(callID, name, outcome, resultText string) {
	if f == nil {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.tools = append(f.tools, ToolRunRow{
		CallID: callID, Name: name, Outcome: outcome, ResultText: resultText,
	})
	if len(f.tools) > runToolRowsKept {
		f.tools = f.tools[len(f.tools)-runToolRowsKept:]
		f.toolsDropped++
	}
}

// ReasoningDelta accumulates one reasoning text delta under the tail window.
func (f *RunFeed) ReasoningDelta(text string) {
	if f == nil || text == "" {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.reasoningTail = append(f.reasoningTail, []rune(text)...)
	f.reasoningSeen += len([]rune(text))
	if over := len(f.reasoningTail) - runReasoningKeepRunes; over > 0 {
		f.reasoningTail = f.reasoningTail[over:]
	}
}

// Errored records the last classified error event. An event that carries no
// class and no detail and no text writes nothing - there is nothing true to
// write.
func (f *RunFeed) Errored(class, detail, text string) {
	if f == nil || (class == "" && detail == "" && text == "") {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.lastError = &RunErrorRow{Class: class, Detail: detail, Text: text}
}

// Stuck records the loop's last stuck sentence and its C22 rung.
func (f *RunFeed) Stuck(text string, repeatLevel int) {
	if f == nil {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.stuck = &RunStuckRow{Text: text, RepeatLevel: repeatLevel}
}

// Section builds the wire shape. A feed that exists always yields a non-nil
// section - the key's presence is the reader's existence, and an empty feed
// renders "tools":[] rather than null, the same rule the commands and
// modelSwitch sections state.
func (f *RunFeed) Section() *RunSection {
	if f == nil {
		return nil
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	tools := make([]ToolRunRow, len(f.tools))
	copy(tools, f.tools)
	sec := &RunSection{
		Tools:        tools,
		ToolsDropped: f.toolsDropped,
		Reasoning:    f.reasoningLocked(),
		Error:        f.lastError,
		Stuck:        f.stuck,
	}
	return sec
}

// reasoningLocked renders the retained reasoning under the bound. Below the
// window the text travels whole; past it the elided count is inside the
// marker, so a bound that ate runes is never invisible.
func (f *RunFeed) reasoningLocked() string {
	elided := f.reasoningSeen - len(f.reasoningTail)
	if elided <= 0 {
		return string(f.reasoningTail)
	}
	return fmt.Sprintf("...[truncated: %d runes elided]...%s", elided, string(f.reasoningTail))
}
