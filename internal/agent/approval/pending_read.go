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

// LiveApproval is one approval still waiting for an answer, as an in-process
// observer sees it. Decision is the verdict as it was admitted (RulesHit,
// Reason, Level and SessionOverrideBlocked are the fields the L2 card shows
// verbatim, ticket 17's frozen-contract note in tools/gate.go); Position is the
// 1-based place in the FIFO, i.e. the depth badge C18 puts on the card.
type LiveApproval struct {
	CorrelationID string
	Decision      tools.Decision
	Position      int
}

// LiveApprovals lists the pending items, head first.
//
// It reports the queue's state at the instant it is called, which for a
// concurrent host means "some answer may have landed a moment later" - that is
// the same honesty the panel's generatedAt stamp already carries, and it is why
// the pump takes a reader instead of a value: the next state change asks again.
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
			Decision:      it.Dec,
			Position:      q.position(it),
		})
	}
	return out
}
