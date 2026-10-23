package memory

// Task-state vocabulary: the ONE authority for task_log.state (ticket 196).
//
// Two vocabularies coexist in this repository and this file is where they meet,
// on purpose and in one place:
//
//	(A) the memory vocabulary - what the L3 task row has actually been written
//	    with since the table existed: running, done, succeeded, cancelled,
//	    interrupted, stuck, error. Provenance, read off the call sites rather
//	    than off prose: internal/agent/journal.go:131 inserts "running";
//	    internal/agent/loop.go's taskLogState maps the loop's terminal status
//	    onto done/stuck/cancelled/interrupted/error; cmd/wisp/run.go's own
//	    taskLogState maps the CLI's onto done/cancelled/error. "succeeded" has
//	    no producer left in this tree - it survives in the package's own tests
//	    (dao_test.go, privacy_test.go) and, in the field, in rows written by
//	    older builds, which is exactly why it must stay accepted.
//
//	(B) the D43 vocabulary - the 20 frozen machine names in
//	    internal/statemachine/states.go:11-31 (D43; C12 frozen). Not one of
//	    (A)'s words is among them, which is the defect ticket 196 was filed
//	    for.
//
// Ticket 196 deliberately does NOT unify them: the ticket's forbidden zone
// says in so many words that task_log.state must not be renamed to D43 names,
// because that would wash away the four words existing producers write. So both
// sets are RETAINED, and the reconciliation lives here instead:
//
//   - writes go through ValidTaskState (the DAO refuses anything off the
//     roster, and the column carries the same roster as a SQL CHECK since
//     schema v4 - see ddlV4 in schema.go);
//   - reads go through TaskStateFor / TaskLog.D43State, so a reader that wants
//     a D43 name asks this table instead of inventing a switch of its own.
//
// The D43 half of the roster is deliberately NOT all twenty names. Only the
// five a task record can legitimately hold are accepted, and they are the ones
// this repository already chose for a task's own state in
// internal/tools/subagent_197.go:107-110 (Thinking while it runs, Settling when
// it ends with an answer, Muted when the host stopped it, Error when it ended
// without one), plus Stuck, which (A) already writes as "stuck". Accepting all
// twenty would put Sleeping and Listening - ball states - within reach of a
// task row and would make the CHECK assert something nobody means.

import (
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// TaskState is one value of the task_log.state column. It is a string type so
// the existing DAO signatures and the DDL keep working unchanged; the type is
// what gives the vocabulary a name to point at.
type TaskState string

// The memory vocabulary (A). Every one of these has been written by a shipped
// producer or by this package's own tests; see the file comment for the census.
const (
	TaskStateRunning     TaskState = "running"     // row opened, no terminal state yet
	TaskStateDone        TaskState = "done"        // loop/CLI completed
	TaskStateSucceeded   TaskState = "succeeded"   // legacy spelling of "done"; no producer left, kept for old rows
	TaskStateCancelled   TaskState = "cancelled"   // user or host cancelled
	TaskStateInterrupted TaskState = "interrupted" // truncated by the round/token ceiling
	TaskStateStuck       TaskState = "stuck"       // loop hit its repeated-call ceiling
	TaskStateError       TaskState = "error"       // ended without an answer
)

// taskStateRoster is the accepted set, in a fixed order so the generated SQL
// CHECK text is stable and readable. Adding a word here is the ONLY way to
// widen the column: the CHECK is built from this slice (schema.go's buildDDLV4)
// and the DAO validates against it, so the two cannot disagree.
var taskStateRoster = []TaskState{
	TaskStateRunning,
	TaskStateDone,
	TaskStateSucceeded,
	TaskStateCancelled,
	TaskStateInterrupted,
	TaskStateStuck,
	TaskStateError,

	// The D43 names a task record may hold (see the file comment).
	TaskState(statemachine.StateThinking),
	TaskState(statemachine.StateSettling),
	TaskState(statemachine.StateMuted),
	TaskState(statemachine.StateStuck),
	TaskState(statemachine.StateError),
}

// taskStateToD43 maps the memory vocabulary onto D43 names. It is total over
// (A): every word a producer can write has an answer, so no reader has to
// invent a fallback. The targets are the four names
// internal/tools/subagent_197.go already uses for a task, and the mapping is
// many-to-one on purpose (both "done" and "succeeded" mean an answer was
// produced; "cancelled" and "interrupted" both mean the host stopped it).
var taskStateToD43 = map[TaskState]statemachine.State{
	TaskStateRunning:     statemachine.StateThinking,
	TaskStateDone:        statemachine.StateSettling,
	TaskStateSucceeded:   statemachine.StateSettling,
	TaskStateCancelled:   statemachine.StateMuted,
	TaskStateInterrupted: statemachine.StateMuted,
	TaskStateStuck:       statemachine.StateStuck,
	TaskStateError:       statemachine.StateError,
}

// taskStateFromD43 is the canonical inverse of taskStateToD43: when several
// memory words share one D43 name, the one listed FIRST in taskStateRoster
// wins, so the round trip is deterministic instead of map-ordered.
var taskStateFromD43 = func() map[statemachine.State]TaskState {
	out := make(map[statemachine.State]TaskState, len(taskStateToD43))
	for _, ts := range taskStateRoster {
		d43, ok := taskStateToD43[ts]
		if !ok {
			// A D43 word already in (A)'s spelling answers for itself.
			d43 = statemachine.State(ts)
		}
		if _, taken := out[d43]; !taken {
			out[d43] = ts
		}
	}
	return out
}()

// TaskStateRoster returns the accepted task_log.state values in roster order.
// The returned slice is a copy: a caller cannot widen the vocabulary by
// appending to it.
func TaskStateRoster() []string {
	out := make([]string, len(taskStateRoster))
	for i, ts := range taskStateRoster {
		out[i] = string(ts)
	}
	return out
}

// ValidTaskState reports whether s is on the roster. This is the judge the DAO
// write path uses, and it is the same set the column's CHECK declares.
func ValidTaskState(s string) bool {
	for _, ts := range taskStateRoster {
		if TaskState(s) == ts {
			return true
		}
	}
	return false
}

// TaskStateFor answers the D43 name for a task_log.state word. It accepts both
// vocabularies: a memory word is translated, a D43 word already on the roster is
// returned as itself. ok=false for anything off the roster, and a non-roster
// word is never silently translated.
func TaskStateFor(s string) (statemachine.State, bool) {
	if !ValidTaskState(s) {
		return "", false
	}
	if d43, ok := taskStateToD43[TaskState(s)]; ok {
		return d43, true
	}
	return statemachine.State(s), true
}

// TaskStateFromD43 answers the canonical memory word for a D43 name. It is the
// read direction for a surface that only ever learned the D43 spelling.
func TaskStateFromD43(s statemachine.State) (TaskState, bool) {
	ts, ok := taskStateFromD43[s]
	return ts, ok
}

// D43State is the read-side door: the task row's state expressed as a D43 name.
// A reader asks this instead of switching on State's spelling itself, so there
// is one translation in the program rather than one per surface.
func (tl TaskLog) D43State() (statemachine.State, bool) {
	return TaskStateFor(tl.State)
}

// invalidTaskStateError is the one sentence a producer sees when it writes a
// word off the roster. The roster is spelled out in the message so the fix is
// readable from the failure alone.
func invalidTaskStateError(s string) error {
	return fmt.Errorf("memory: invalid task_log.state %q (want one of: %s)",
		s, strings.Join(TaskStateRoster(), "|"))
}
