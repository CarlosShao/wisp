package approval

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// Ticket 285 AC#1 (leg 285-r1). Placement is form A: the whole case lives inside
// this package, no production byte moved, cmd/wisp untouched.
//
// The gap this file closes was named by the non-implementer leg 242-v2
// (docs/evidence/s1/242-grant-binding-v2.md:102-110). Two cards that are BOTH
// live at the moment of the answer, the second card's correlation id named by the
// call, the first card's UNSPENT grant handed to it, and the audit face read back
// to see WHICH denial that booked - nobody asserted that anywhere. What existed:
//   - the only assertion in this package putting a foreign correlation id on
//     Native().Allow was TestGrantIsSingleUseAndBoundToItsItem in queue_test.go,
//     and the id it aims at ("corr-B") is never pushed. That call returns
//     ErrUnknownCorrelation and books zero GRANT-DENY lines, so it never read a
//     denial name; the "not a route that is not running" dimension was carried by
//     a missing card, not by the store.
//   - the genuine two-live-cards shot lived only in
//     cmd/wisp/subagent_selfapproval_197_test.go, which asserts the error value
//     and never captures the gate's Logf, so the denial name was dropped on the
//     floor the moment it was written.
//
// Scope, held deliberately:
//   - the outward answer stays the one merged ErrBadGrant (ticket 259 form a,
//     ledger A562 / A619). Nothing below lifts a denial token onto an error value
//     or a card face, and TestTicket259R1OutwardAnswerStaysMergedAcrossDenials
//     remains the ruler for that merge.
//   - this case does NOT claim the route can produce denialMisbound, and it does
//     not claim it cannot. It claims the route NAMES which denial it did produce,
//     and that the name is the "the presented proof is not a live row in the store
//     of the card the answer named" family - which is where a cross-card grant
//     lands today, because grantStore is one per qitem (queue.go push) and
//     grantStore.spend scans membership before it compares a binding.
//
// Fixture reuse rather than a new one: r1Audit, r1Card and r1CardIn from
// ticket259_denial_rulers_test.go. r1CardIn exists precisely so one case can hold
// two live cards in one funnel behind one audit face; rebuilding that here would
// have put a second copy of the capture in the package.

const (
	// xcardDenialToken is the audit token this case demands for a proof that is
	// not a row in the named card's own store. Written as a literal and NOT as
	// denialSpentNonce.label(), so an edit that moves the production label and a
	// derived expectation together in one stroke cannot pass this silently.
	xcardDenialToken = "denial=spent-or-never-live-nonce"

	// xcardOtherNames are the three remaining denial names the audit can print.
	// The cross-card shape is supposed to land in exactly one bucket; if the
	// branch it travels moves to another, this is the case that says so.
	xcardOtherNames = "denial=missing-nonce|denial=misbound|denial=card-not-pending"
)

// xcardLive asserts one card is still pending AND still indexed under the
// correlation id this case handed it. That is what turns the refusal asserted
// below into a grant denial rather than "this route is not running": with both
// guards open, the only thing left that can refuse the answer is the store.
func xcardLive(t *testing.T, q *Queue, it *qitem, who string) {
	t.Helper()
	q.mu.Lock()
	defer q.mu.Unlock()
	if it.state != statePending {
		t.Fatalf("285 fixture: %s (%s) left the pending set before the shot; the denial below would then be booked by the state guard, not by the grant store", who, it.Corr)
	}
	if got := q.lookupForAllowLocked(it.Corr); got != it {
		t.Fatalf("285 fixture: %s (%s) is not the item the allow side resolves by its own correlation id; the answer below would not reach this card", who, it.Corr)
	}
}

// AC#1, the missing reading: two live cards, A's unspent grant aimed at B, and
// the audit face naming which refusal that booked.
func TestTicket285R1CrossCardAllowNamesTheDenialInTheAudit(t *testing.T) {
	g, q, cardA, nonceA, a := r1Card(t, "task-285-xcard-a", "corr-285-xcard-a", "shell.run")
	cardB, _ := r1CardIn(t, q, "task-285-xcard-b", "corr-285-xcard-b", "shell.run")

	// Two cards, two ids, both queued at once. The queue silently rewrites a
	// colliding correlation id (queue.go push), so assert the ids this case means
	// are the ids the queue issued.
	if cardA.Corr != "corr-285-xcard-a" || cardB.Corr != "corr-285-xcard-b" {
		t.Fatalf("285 fixture: a correlation id was rewritten (%s / %s); the shape asserted below is not the shape on the queue", cardA.Corr, cardB.Corr)
	}
	if cardA.Corr == cardB.Corr {
		t.Fatalf("285 fixture: both cards carry the same correlation id (%s); there is no cross-card shot to make", cardA.Corr)
	}
	xcardLive(t, q, cardA, "card A")
	xcardLive(t, q, cardB, "card B")
	if n := q.Depth(); n != 2 {
		t.Fatalf("285 fixture: depth=%d, want both cards live at once. 242-v2's complaint is precisely that the second card was never pushed; a one-card run would reproduce ErrUnknownCorrelation instead of a grant denial", n)
	}

	// A's grant is unspent, and it is A's: one live row in A's store, one live row
	// in B's. Without this the word "unspent" in the criterion is decorative.
	if n := cardA.grants.live(); n != 1 {
		t.Fatalf("285 fixture: card A holds %d live grant(s), want 1; the value handed to B below is not an unspent proof", n)
	}
	if n := cardB.grants.live(); n != 1 {
		t.Fatalf("285 fixture: card B holds %d live grant(s), want 1; B is not a card a native answer could have opened, which changes what the refusal means", n)
	}

	err := g.Native().Allow(context.Background(), cardB.Corr, nonceA)

	// Outward: one merged error value, unchanged (ticket 259 form a).
	if !errors.Is(err, ErrBadGrant) {
		t.Fatalf("AC#1 RED: the two-live-card cross-card answer returned %v, want the merged ErrBadGrant", err)
	}
	// The two names this shape must NOT answer with, because they are exactly how
	// a green-looking run could hide the hole 242-v2 named: refusing a card that
	// is not there, or a card that already left.
	if errors.Is(err, ErrUnknownCorrelation) {
		t.Fatalf("AC#1 RED: the answer aimed at a live card (%s) came back ErrUnknownCorrelation; the card stopped being reachable before the store was consulted, so the denial named below proves nothing about grants", cardB.Corr)
	}
	if errors.Is(err, ErrNotPending) {
		t.Fatalf("AC#1 RED: the answer aimed at a pending card (%s) came back ErrNotPending; the state guard took this shot, not the grant store", cardB.Corr)
	}

	// The reading this ticket exists for: the audit face names WHICH denial the
	// route booked, and the name is the not-live-in-B's-store family.
	want := "approval: GRANT-DENY corr=" + cardB.Corr + " tool=shell.run " + xcardDenialToken
	if !a.has(want) {
		t.Fatalf("AC#1 RED: the assertion face cannot read which denial the two-live-card cross-card route booked (want %q)\naudit:\n%s", want, a.dump())
	}

	// Attributed to the card the ANSWER pointed at, not to the card that minted
	// the proof: an auditor following this line has to be sent to B.
	if a.has("approval: GRANT-DENY corr=" + cardA.Corr) {
		t.Errorf("AC#1 RED: the cross-card denial was booked against card A (%s) although the answer named card B (%s); the audit then describes the wrong refusal\naudit:\n%s", cardA.Corr, cardB.Corr, a.dump())
	}
	dump := a.dump()
	if n := strings.Count(dump, "approval: GRANT-DENY"); n != 1 {
		t.Errorf("AC#1 RED: one refused answer booked %d GRANT-DENY lines, want exactly 1; a second line means another guard answered this call too\naudit:\n%s", n, dump)
	}
	for _, other := range strings.Split(xcardOtherNames, "|") {
		if strings.Contains(dump, other) {
			t.Errorf("AC#1 RED: the cross-card route booked %s alongside or instead of %s; this case names the not-live-in-this-store family, so the branch this answer travels has moved\naudit:\n%s", other, xcardDenialToken, dump)
		}
	}
	// The merged outward-facing audit sentence this funnel has always written is
	// still written next to the named line (queue.go keeps them as two lines).
	if !a.has("approval: FORGED-OR-STALE allow rejected corr=" + cardB.Corr) {
		t.Errorf("AC#1 RED: the merged FORGED-OR-STALE sentence for the answered card is gone; the named line and the merged line are two lines on purpose\naudit:\n%s", dump)
	}

	// Positive control, and it is the load-bearing one: the refused attempt burned
	// nothing on A, and A still opens with the very same value. That is what makes
	// the refusal above a statement about B's store rather than about the quality
	// of the proof. Run last, because spending settles A.
	if n := cardA.grants.live(); n != 1 {
		t.Fatalf("AC#1 RED: A's grant did not survive an answer aimed at B (%d live); the foreign attempt reached A's store, which is a different and worse bug than the one this case reads", n)
	}
	if err := g.Native().Allow(context.Background(), cardA.Corr, nonceA); err != nil {
		t.Fatalf("AC#1 RED: card A refused its own live grant after a foreign answer aimed at B (%v); the proof the cross-card call carried was live, so the control says the denial above was not about the value", err)
	}
	if !a.has("approval: ANSWER-ALLOW corr=" + cardA.Corr) {
		t.Errorf("AC#1 RED: the control allow booked no ANSWER-ALLOW line for card A\naudit:\n%s", a.dump())
	}
}

// The other half of the same gap, and the reason the case above is not vacuous:
// the shape queue_test.go actually asserts - a foreign correlation id that was
// never pushed - refuses for a reason the grant store never sees, and books no
// GRANT-DENY line at all. Pinning that here is what stops a future edit from
// reading "a denial was refused" off the unknown-correlation case and calling
// AC#1 covered.
func TestTicket285R1UnknownCorrelationBooksNoDenialLine(t *testing.T) {
	g, q, cardA, nonceA, a := r1Card(t, "task-285-unknown-a", "corr-285-unknown-a", "shell.run")
	const ghost = "corr-285-never-pushed"

	xcardLive(t, q, cardA, "card A")
	if n := cardA.grants.live(); n != 1 {
		t.Fatalf("285 fixture: card A holds %d live grant(s), want 1", n)
	}

	err := g.Native().Allow(context.Background(), ghost, nonceA)
	if !errors.Is(err, ErrUnknownCorrelation) {
		t.Fatalf("AC#1 RED: an answer aimed at a correlation id that was never pushed returned %v, want ErrUnknownCorrelation", err)
	}
	if n := strings.Count(a.dump(), "approval: GRANT-DENY"); n != 0 {
		t.Fatalf("AC#1 RED: the unknown-correlation route booked %d GRANT-DENY line(s); once it names a denial too, reading a denial name off that call is no longer proof of the two-live-card shape\naudit:\n%s", n, a.dump())
	}
}
