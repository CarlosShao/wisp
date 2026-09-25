package approval

import "github.com/CarlosShao/wisp/internal/tools"

// The read side of the live approval state (ticket 35's snapshot pump).
//
// WHY A NEW FUNCTION WHEN PanelAPI ALREADY EXISTS. approval.PanelAPI (ui.go) is
// the projection the UNTRUSTED side is allowed to hold, and its header says so:
// PanelItem deliberately carries no Params and no grant. The snapshot pump is
// not that side - it runs in-process, on the native side of the boundary, and
// the packet it builds for the panel is the one the panel renders as an L2 card,
// which has to name the rules that fired (ApprovalCardView.rulesHit) and whether
// R4 blocked a scope override (sessionOverrideBlocked). Both of those live on
// qitem.Dec (a tools.Decision), which PanelItem drops. Re-adding them to
// PanelItem would widen the surface a page can ask for; reading them here widens
// nothing, because this function is not on PanelAPI and cannot answer a request
// that arrives from the page.
//
// WHAT THIS IS NOT: an answer, a grant, or an authority. It returns copies of
// verdicts the assessor already reached and the queue already admitted. It
// cannot allow anything, cannot be called by the renderer, and moves no item out
// of pending. The three ways to settle an item stay exactly where they were:
// Queue.allow (native, needs a grant), Queue.reject (needs nothing), and expiry.
//
// "COPIES" IS LOAD-BEARING, AND TICKET 146 MADE IT SO. tools.Decision is copied
// by value, and it carries reference-typed slots - Params plus seven more
// slices (Args, RulesHit, Paths, Capabilities, and the three inside Blacklist).
// A bare value copy therefore shared all eight backings with the queue's stored
// record, so one in-place write by a caller rewrote an admitted verdict.
// cloneDecision copies each of those slots, which is what makes the sentence
// above true.
// The copy is ONE level deep: what a caller cannot now do is write a slot the
// queue reports; what it can still do, and what no Go type here can prevent, is
// reach THROUGH Params' values into the containers a JSON decode left inside
// them (a []any or nested map[string]any is still shared). Treat those values
// as read-only too. ticket146_liveapprovals_backing_test.go is the instrument
// for the eight named slots, and it walks tools.Decision by reflection so a
// ninth reference-typed field cannot arrive quietly.

// LiveApproval is one approval still waiting for an answer, as an in-process
// observer sees it. Decision is the verdict as it was admitted (RulesHit,
// Reason, Level and SessionOverrideBlocked are the fields the L2 card shows
// verbatim, ticket 17's frozen-contract note in tools/gate.go), copied per
// LiveApprovals so writing to it cannot reach the queue; Position is the
// 1-based place in the FIFO, i.e. the depth badge C18 puts on the card.
type LiveApproval struct {
	CorrelationID string
	Decision      tools.Decision
	Position      int
}

// cloneDecision copies one verdict out of the queue's own record. Every field
// of tools.Decision whose type is a map or a slice gets its own backing; the
// scalars and strings travel by the struct copy. A nil slot stays nil, because a
// caller cannot be asked to tell "no rules hit" from "the queue handed me an
// empty slice it allocated".
func cloneDecision(d tools.Decision) tools.Decision {
	out := d
	out.Params = cloneParamsMap(d.Params)
	out.Args = cloneBacking(d.Args)
	out.RulesHit = cloneBacking(d.RulesHit)
	out.Paths = cloneBacking(d.Paths)
	out.Capabilities = cloneBacking(d.Capabilities)
	out.Blacklist.Absolute = cloneBacking(d.Blacklist.Absolute)
	out.Blacklist.Unlockable = cloneBacking(d.Blacklist.Unlockable)
	out.Blacklist.AlreadyUnlocked = cloneBacking(d.Blacklist.AlreadyUnlocked)
	return out
}

// cloneBacking copies one slice's elements into fresh memory.
func cloneBacking[S ~[]E, E any](in S) S {
	if in == nil {
		return nil
	}
	return append(make(S, 0, len(in)), in...)
}

// cloneParamsMap copies the parameter map one level (see the header note on
// what that does and does not protect).
func cloneParamsMap(in map[string]any) map[string]any {
	if in == nil {
		return nil
	}
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

// LiveApprovals lists the pending items, head first.
//
// It reports the queue's state at the instant it is called, which for a
// concurrent host means "some answer may have landed a moment later" - that is
// the same honesty the panel's generatedAt stamp already carries, and it is why
// the pump takes a reader instead of a value: the next state change asks again.
//
// COST, stated not assumed: one cloneDecision per reported row, so one map plus
// up to seven fresh backings per pending item per call. The bound is the queue's
// own (DefaultMaxPending, queue.go), and the pump does not run on streamed text
// deltas (cmd/wisp/panel_pump.go). Whether that moves a resource band is ticket
// 146 AC#3's question; its readings - including any band that ended up
// unmeasured - are recorded in docs/evidence/s1/146-liveapprovals-r1.md §4.
func (q *Queue) LiveApprovals() []LiveApproval {
	if q == nil {
		return nil
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]LiveApproval, 0, len(q.pending))
	for _, it := range q.pending {
		if it == nil || it.state != statePending {
			continue
		}
		out = append(out, LiveApproval{
			CorrelationID: it.Corr,
			Decision:      cloneDecision(it.Dec),
			Position:      q.position(it),
		})
	}
	return out
}
