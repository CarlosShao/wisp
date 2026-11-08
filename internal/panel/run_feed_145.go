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
// Verified event sources this carrier rides (readings 2026-11-07/08, anchors
// in the evidence file):
//   - tool rows: EvToolEnd only (loop.go's executeCalls / truncated /
//     failOpenCalls publish sites). The outcome vocabulary is internal/agent's
//     own journal vocabulary (success | error | cancelled | truncated),
//     carried through untouched like the D43 names in the tasks section; this
//     package spells none of the words. The row's argsJson is NOT an event
//     field: it is looked up from the call's booked forensics row
//     (memory.ToolCall.ArgsJSON, the bridge's sanitized booking) by the
//     C18 correlation id the event now carries - the census's named source
//     for the expanded-args dimension. A call the bridge never booked
//     (rejected in-loop, list_tools answered locally) has no row and
//     travels with argsJson absent - "never executed with args" is the
//     truth, not a gap to fill.
//   - running rows: EvToolStart - forwarded by loop.stream from the
//     adapters' own llm.EvToolCallStart (all three adapters emit it; until
//     2026-11-08 the agent kind had zero producers, the registered 33-r14
//     AC#8b zero-emitter). A running row retires when its call's EvToolEnd
//     arrives, or when the task itself ends (EvDone) - a task that died
//     between start and end must not leave a chip running forever.
//   - reasoning: EvReasoningDelta text deltas, accumulated under a per-rune
//     bound with an honest elided count (same honesty rule as StreamLog's
//     markers: the bound is counted, never silent).
//   - usage: EvUsage - the Cost accumulator's cumulative totals at the
//     moment a turn's usage landed (publishUsage), the same numbers Result
//     books. Micros and Currency travel as the pair census §2.13's unit rule
//     demands; neither is rendered by this package.
//   - error: EvError events that carry the loop's classified error - the
//     class and the log-safe detail travel verbatim, plus the class policy's
//     own Retryable answer (observe.Error.Retryable). The human sentence
//     for an error state has no source in the tree and is NOT invented here.
//   - stuck: EvStuck - the loop's own visible brake sentence plus the C22
//     ladder rung that fired. The three buttons the table draws next to it
//     are inbound C17 surface (Q-76) and arrive with nothing here.
//
// Registered as NOT landable today (no event source, so landing any of them
// would be a constant dressed as data - the exact shape ticket 145's AC#1
// forbids): the terminal status word and cancel reason (census §2.14's
// orientation - Event.Stop vs Result.Status - is still unruded, and
// EvDone.Text is res.Text || the terminal sentence, so a cancelled run's
// event text is its partial REPLY, not a reason), tool duration (the
// forensics row's started_at has no production writer, so ended-start is
// not computable - census §2.6), phase / elapsed timers, the human error
// sentence, and the L2 window's remainingMs.
// pending[].position and context.occupancy are ticket 167's territory (the
// two legs must not each add their own pending[] keys); composer.git's
// branches half is ticket 186's; the "which screen" view dimension is held
// by AC#3's verdict (the Go side does not know the screen names, and the
// panel-to-Go half reopens the C17 method whitelist).

// ToolRunRow is one ENDED tool call. Rows are written by EvToolEnd only, in
// end order; the fields are the event's own CallID, ToolName, Outcome and
// result text, verbatim. ArgsJSON is the call's booked forensics row's
// sanitized args (memory.ToolCall.ArgsJSON), joined by the event's C18
// correlation id at the feed line - a call with no booked row travels with
// it absent (omitempty), which is the truth "nothing executed these args".
type ToolRunRow struct {
	CallID     string `json:"callId"`
	Name       string `json:"name"`
	Outcome    string `json:"outcome"`
	ResultText string `json:"resultText,omitempty"`
	ArgsJSON   string `json:"argsJson,omitempty"`
}

// ToolRunningRow is one call IN FLIGHT: the adapters' tool_call_start as the
// loop's EvToolStart forwarded it. No outcome exists yet - an outcome here
// would be invented - and the row retires when its call ends or the task
// does.
type ToolRunningRow struct {
	CallID string `json:"callId"`
	Name   string `json:"name"`
}

// RunUsageRow is the Cost accumulator's cumulative spend for the live task,
// as the EvUsage events carried it. Micros and Currency are the census
// §2.13 pair: they travel together or the row does not travel, so no
// consumer can hang a currency symbol on a bare number.
type RunUsageRow struct {
	TokensIn     int    `json:"tokensIn"`
	TokensOut    int    `json:"tokensOut"`
	CachedTokens int    `json:"cachedTokens,omitempty"`
	Micros       int64  `json:"micros"`
	Currency     string `json:"currency"`
}

// RunErrorRow is the last classified error the loop published. Class and
// Detail are observe.Error's own fields as the event carried them; Text is
// the event's message string. Retryable is the class policy's own answer
// (observe.Error.Retryable), carried verbatim - a false here means "this
// class never retries", which is information, not absence. A notice-shaped
// EvError (no classified error attached) does not write this row - inventing
// a class for it would be the constant-as-data shape again.
type RunErrorRow struct {
	Class     string `json:"class"`
	Detail    string `json:"detail,omitempty"`
	Text      string `json:"text,omitempty"`
	Retryable bool   `json:"retryable"`
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
	// Running lists the calls in flight right now, in start order. Empty
	// (omitted) while nothing runs - which is also what a finished task
	// leaves behind, because EvDone retires every running row.
	Running []ToolRunningRow `json:"running,omitempty"`
	// Usage is the task's cumulative spend as the EvUsage events carried it,
	// or nil until the first usage event arrived. Absent is honest: a task
	// whose provider reported no usage has no numbers to show.
	Usage *RunUsageRow `json:"usage,omitempty"`
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
	running       []ToolRunningRow
	usage         *RunUsageRow
	reasoningTail []rune
	reasoningSeen int
	lastError     *RunErrorRow
	stuck         *RunStuckRow
}

// NewRunFeed builds an empty feed.
func NewRunFeed() *RunFeed { return &RunFeed{} }

// ToolStarted records one call in flight, in start order. A call id already
// running is not re-added: a duplicate start cannot make two chips out of
// one call.
func (f *RunFeed) ToolStarted(callID, name string) {
	if f == nil || callID == "" {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, r := range f.running {
		if r.CallID == callID {
			return
		}
	}
	f.running = append(f.running, ToolRunningRow{CallID: callID, Name: name})
}

// ToolEnded records one ended tool call and retires its running row - the
// same call cannot be both. Oldest rows past the cap are evicted and
// counted, never folded into a neighbour. argsJSON is the booked forensics
// row's sanitized args as the feed line looked it up; empty stays absent.
func (f *RunFeed) ToolEnded(callID, name, outcome, resultText, argsJSON string) {
	if f == nil {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	kept := f.running[:0:0]
	for _, r := range f.running {
		if r.CallID != callID {
			kept = append(kept, r)
		}
	}
	f.running = kept
	f.tools = append(f.tools, ToolRunRow{
		CallID: callID, Name: name, Outcome: outcome, ResultText: resultText,
		ArgsJSON: argsJSON,
	})
	if len(f.tools) > runToolRowsKept {
		f.tools = f.tools[len(f.tools)-runToolRowsKept:]
		f.toolsDropped++
	}
}

// Usage records the Cost accumulator's cumulative totals as one EvUsage
// event carried them. Micros and Currency are stored as the pair they
// arrived in; this method derives nothing.
func (f *RunFeed) Usage(in, out, cached int, micros int64, currency string) {
	if f == nil {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.usage = &RunUsageRow{
		TokensIn: in, TokensOut: out, CachedTokens: cached,
		Micros: micros, Currency: currency,
	}
}

// TaskEnded retires every running row: a task that died between a call's
// start and its end must not leave a chip running forever. The ended rows
// and every other dimension stay - they are the run's record.
func (f *RunFeed) TaskEnded() {
	if f == nil {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.running = nil
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

// Errored records the last classified error event: the class policy's
// retryable answer travels with it. An event that carries no class and no
// detail and no text writes nothing - there is nothing true to write.
func (f *RunFeed) Errored(class, detail, text string, retryable bool) {
	if f == nil || (class == "" && detail == "" && text == "") {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.lastError = &RunErrorRow{Class: class, Detail: detail, Text: text, Retryable: retryable}
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
		Usage:        f.usage,
	}
	if len(f.running) > 0 {
		running := make([]ToolRunningRow, len(f.running))
		copy(running, f.running)
		sec.Running = running
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
