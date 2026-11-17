package main

// Readface wave 2 (2026-11-17), ticket 38 AC#4's task-card sources: the
// composition root's live per-task ledger that the roster reader
// (panel_pump.go's taskRosterState) reads the three new TaskRowView fields
// from. It is a HOLDER, not a judge: every value is written by an event the
// console sink already receives (internal/agent/sink.go's EvToolStart /
// EvToolEnd / EvUsage / EvDone), carried verbatim, and read back verbatim.
// Nothing is derived, defaulted or invented here; a dimension the events do
// not carry stays at its zero and the wire keeps it absent (omitempty).
//
// WHY A LEDGER AND NOT A QUERY. The three facts the ticket names exist as
// events on the sink face and as rows on the store face:
//
//   - startedAt: stamped at admitTask/noteTask - the moment the loop-minted
//     task id first exists in this process, the same hook that makes the row
//     nameable. The roster is process-local by design (164 定案②), so the
//     admission stamp is its matching scope; the DB row's started_at is the
//     persisted twin, written by the loop's journal from its own start.
//   - costToDate: the loop's Cost accumulator publishes its cumulative
//     task-so-far totals on EvUsage (CostMicros and the currency pair travel
//     together or not at all, census 145 §2.13); the ledger keeps the latest
//     per task. A subagent's spend is booked on its own finalize and its
//     EvUsage never reaches this sink (internal/tools' subagentTextSink
//     carries text deltas and EvDone only) - so a subagent row reads
//     costToDate absent, which is the honest "this host's sink published no
//     spend event for this row", never "cost zero".
//   - currentStep: the tool call the loop has started and not yet ended
//     (EvToolStart / EvToolEnd), retired on end and on the task's EvDone -
//     the same retirement rule RunFeed's running rows follow. A subagent's
//     tool events likewise do not reach this sink, so its rows read
//     currentStep absent.
//
// The mutex covers the two goroutines that genuinely meet here: the task
// goroutine publishing events (consoleSink.Publish, the loop's own call) and
// the pump's readers (approval UI's and sink's publishes).

import (
	"sync"

	"github.com/CarlosShao/wisp/internal/panel"
)

// taskLiveLedger is the per-task live half of the roster's three card fields.
type taskLiveLedger struct {
	mu     sync.Mutex
	starts map[string]int64
	usage  map[string]*panel.RunUsageRow
	step   map[string]string
}

func newTaskLiveLedger() *taskLiveLedger {
	return &taskLiveLedger{
		starts: map[string]int64{},
		usage:  map[string]*panel.RunUsageRow{},
		step:   map[string]string{},
	}
}

// startStamped records the admission moment of one task id. Duplicate ids are
// not re-stamped: the FIRST admission is the start, and admitTask's idempotent
// registration must not move a row's clock.
func (l *taskLiveLedger) startStamped(taskID string, unix int64) {
	if l == nil || taskID == "" || unix == 0 {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.starts[taskID]; ok {
		return
	}
	l.starts[taskID] = unix
}

// usageRecorded keeps the latest cumulative spend an EvUsage carried for one
// task. The micros/currency pair travels as it arrived - this method derives
// nothing and never splits it.
func (l *taskLiveLedger) usageRecorded(taskID string, in, out, cached int, micros int64, currency string) {
	if l == nil || taskID == "" {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.usage[taskID] = &panel.RunUsageRow{
		TokensIn: in, TokensOut: out, CachedTokens: cached,
		Micros: micros, Currency: currency,
	}
}

// stepStarted records one in-flight call; stepEnded retires it - the same
// call cannot be both. taskDone retires whatever step a task was in when its
// terminal event arrived: a task that died between a start and an end must
// not leave a step "in flight" forever (RunFeed's rule, same holder family).
func (l *taskLiveLedger) stepStarted(taskID, tool string) {
	if l == nil || taskID == "" || tool == "" {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.step[taskID] = tool
}

func (l *taskLiveLedger) stepEnded(taskID string) {
	if l == nil || taskID == "" {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.step, taskID)
}

func (l *taskLiveLedger) taskDone(taskID string) {
	l.stepEnded(taskID)
}

// cardFields reads the three values for one row. A task the ledger never saw
// reads zero/nil/empty - the wire keeps those absent, and the pump's reader
// copies them without filling.
func (l *taskLiveLedger) cardFields(taskID string) (startedAt int64, cost *panel.RunUsageRow, step string) {
	if l == nil || taskID == "" {
		return 0, nil, ""
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.starts[taskID], l.usage[taskID], l.step[taskID]
}
