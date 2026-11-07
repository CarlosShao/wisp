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
// above true. Since ticket 148 it is the copy primitive of BOTH outbound doors:
// Queue.replay uses it for the whole Decision it re-opens (the sibling outlet
// ticket 148 AC#1 closed), and panelVerdict below uses it for the two slots the
// whitelist still carries (Args, RulesHit).
// The copy is ONE level deep: what a caller cannot now do is write a slot the
// queue reports; what it can still do, and what no Go type here can prevent, is
// reach THROUGH Params' values into the containers a JSON decode left inside
// them (a []any or nested map[string]any is still shared). Treat those values
// as read-only too. ticket146_liveapprovals_backing_test.go is the instrument
// for the eight named slots on cloneDecision, and it walks tools.Decision by
// reflection so a ninth reference-typed field cannot arrive quietly. Ticket
// 148 Ⓑ removed the reachable path through LiveApprovals: the panel payload
// does not carry Params at all, so the second-level containers cannot leave the
// queue through that door any more.

// LiveApproval is one approval still waiting for an answer, as an in-process
// observer sees it. Decision is the verdict as it was admitted, and the two
// sentences about it are not the same sentence: the frozen-contract note in
// internal/tools/gate.go:12-13 (ticket 17) names exactly TWO fields the card
// renders VERBATIM - RulesHit and Reason - while the four the L2 card actually
// carries (Level/RulesHit/Reason/SessionOverrideBlocked) come from the panel's
// own mapping note at cmd/wisp/panel_pump.go:44-46, which calls that mapping
// "a copy, not a judgement" and is not a frozen contract. panel.CardViewFromDecision
// (internal/panel/approval.go:72-88) is where all four are read onto the view.
// Neither note is enforced by tools/d22scan (its pathresolver-bypass ban is
// about filepath.Clean/Abs outside risk.PathResolver, and its phantom-citation
// ban checks only that a cited path exists), so both are conventions this
// header states rather than invariants a gate checks. Decision is projected per
// LiveApprovals through panelVerdict, so the row carries only the fields the
// panel reader names and shares no backing with the queue for the reference
// slots it does carry; Position is the 1-based place in the FIFO, i.e. the
// depth badge C18 puts on the card.
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

// panelVerdict is the OUTBOUND FIELD WHITELIST of LiveApprovals (ticket 148
// Ⓑ, owner direction: name the narrowing as "the fields the panel needs"). The
// set is not a judgement, it is a census of the one production reader:
// cmd/wisp/panel_pump.go's liveVerdicts reads exactly six fields off the row's
// Decision - Tool, Args (through panelArgs), Level, RulesHit, Reason and
// SessionOverrideBlocked - plus CorrelationID and Position, which travel on the
// LiveApproval wrapper above and are therefore not Decision fields at all.
// Everything else a tools.Decision carries stays in the queue:
//
//   - Params is NOT on the list, and that is the ticket-146 second-level hole's
//     disposal, not an omission: the containers a JSON decode left inside
//     Params' values ([]any, nested map[string]any) were still shared after the
//     one-level clone, and no Go copy here can fix that. A field that never
//     leaves the queue cannot have its containers reached from outside, so the
//     hole closes by exclusion. cloneDecision keeps its one-level Params copy
//     for Queue.replay, whose returned Decision is consumed in-process by the
//     approval flow, not rendered by the panel.
//   - Paths, Capabilities, Blacklist, Refusal, Mode*, the id strings and the
//     timeouts are likewise not panel payload; the untrusted card projection
//     (PanelItem, queue.go's viewLocked) builds its own bounded view and never
//     reads them off this row.
//
// The two reference slots the whitelist DOES carry (Args, RulesHit) are copied
// through the same cloneBacking primitive cloneDecision uses, so an in-place
// write by a reader still cannot reach the queue's record. The instrument that
// keeps this list honest is ticket148_panel_whitelist_test.go: it walks
// tools.Decision by reflection, fails closed on a field nobody classified, and
// goes red when the whitelist and the census drift apart in either direction.
func panelVerdict(d tools.Decision) tools.Decision {
	return tools.Decision{
		Tool:                   d.Tool,
		Args:                   cloneBacking(d.Args),
		Level:                  d.Level,
		RulesHit:               cloneBacking(d.RulesHit),
		Reason:                 d.Reason,
		SessionOverrideBlocked: d.SessionOverrideBlocked,
	}
}

// LiveApprovals lists the pending items, head first.
//
// It reports the queue's state at the instant it is called, which for a
// concurrent host means "some answer may have landed a moment later" - that is
// the same honesty the panel's generatedAt stamp already carries, and it is why
// the pump takes a reader instead of a value: the next state change asks again.
//
// COST, stated not assumed: one panelVerdict per reported row, so two fresh
// backings (Args, RulesHit) per pending item per call - down from
// cloneDecision's map-plus-seven, and down further from the pre-ticket-146 bare
// value copy that shared every backing. The bound is the queue's own
// (DefaultMaxPending, queue.go), and the pump does not run on streamed text
// deltas (cmd/wisp/panel_pump.go). The band-level readings, including any that
// ended up unmeasured, are recorded in docs/evidence/s1/146-liveapprovals-r1.md
// §4 and re-taken in docs/evidence/s1/148-r3.md §AC#2.
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
			Decision:      panelVerdict(it.Dec),
			Position:      q.position(it),
		})
	}
	return out
}
