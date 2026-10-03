package approval

import "sort"

// The read side of the L1 confirmation windows (ticket 220 AC#2, landing point
// 甲 as ruled by the orchestrator and booked in the ledger as A562).
//
// WHY A NEW NAME WHEN LiveApprovals ALREADY EXISTS. LiveApprovals (pending_read.go)
// is the L2 queue's view, and an L1 window is deliberately not a queue item - see
// gate.go's own note beside the ANSWER-VETO line: the queue's ANSWER-ALLOW /
// ANSWER-REJECT bookkeeping can only ever describe an L2 card. So the one
// dimension "which tasks are blocked in the 2-3s block right now" had NO producer
// anywhere in the tree: g.windows is written by openWindow / closeWindow
// (gate.go:346-360) and read by name by Veto (gate.go:415ff), which needs the id
// it is looking for already. A panel that wants to say "this row is waiting" can
// ask by name only when it already knows the answer.
//
// WHAT THIS IS NOT: an answer, a veto, or a second door. It returns three id
// strings per live window and no method that could act on one. It cannot allow
// (there is no route to an allow in this package that does not go through
// Queue.allow with a native grant, Veto, or the window's own expiry), it cannot
// veto (it never hands out the window's channel or context), and it moves nothing
// out of g.windows. The three ways to settle an L1 window stay exactly where they
// were. TestReadingL1WindowsMovesNoGateState is the instrument for that sentence,
// and it is a white-box case on purpose: "the read moved no state" is only
// measurable against the unexported maps.
//
// IT IS NOT ON C17. Nothing a page can call reaches this: internal/panel receives
// a copy of these ids through a READER the composition root supplies
// (PumpSources.L1Windows), and the inbound roster ruler
// (internal/panel/inbound_roster_253_test.go) keeps the closed set of page-callable
// method names at six. Widening that set is a contract change and is not this
// file's to make.

// L1Window is one confirmation window an in-process observer sees as waiting for
// a veto right now. CorrelationID is the id the window is keyed by in g.windows
// (which is the task id whenever the decision carried no correlation of its own,
// per orDefaultText at gate.go:275); TaskID is the task the call belongs to, so a
// reader that pairs rows by task id can join without guessing at the correlation.
// Tool names what the window is holding back. All three are copies of strings the
// gate already filed - a read, never a re-assessment - and the type carries no
// behaviour at all.
type L1Window struct {
	CorrelationID string
	TaskID        string
	Tool          string
}

// LiveL1Windows lists the confirmation windows open at the instant it is called,
// ordered by correlation id so a caller that puts the result on a hash-logged
// packet gets reproducible bytes.
//
// Like LiveApprovals, it reports state as of the call: a veto or an expiry a
// moment later is a different instant's truth, which is why the pump takes a
// reader rather than a value. COST, stated not assumed: one small slice and one
// sort per call, bounded by the number of live windows; it copies no decision and
// no parameter map, so it is strictly cheaper than LiveApprovals per row.
//
// A window past its deadline is NOT listed: markStarted has already moved that
// correlation into g.running, and "still executing" is not "waiting for a veto".
// Lighting a blocked cell for a call that is writing files would be the mirror
// image of the lie this ticket is fixing.
func (g *Gate) LiveL1Windows() []L1Window {
	if g == nil {
		return nil
	}
	g.mu.Lock()
	out := make([]L1Window, 0, len(g.windows))
	for _, w := range g.windows {
		if w == nil {
			continue
		}
		out = append(out, L1Window{CorrelationID: w.corr, TaskID: w.taskID, Tool: w.tool})
	}
	g.mu.Unlock()
	sort.Slice(out, func(i, j int) bool { return out[i].CorrelationID < out[j].CorrelationID })
	return out
}
