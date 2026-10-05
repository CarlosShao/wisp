package approval

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/tools"
)

// Ticket 259 AC#2 (leg 259-r1): "spend" used to answer with a bool, so the four
// refusal causes - no proof presented, proof not live on this card, proof live
// but bound elsewhere, card already out of pending - folded into one reply and
// the audit line could only list them instead of naming which one happened.
//
// Scope, exactly as the orchestrator ruled it: the DISTINCTION lives in the
// internal value and in the audit record; the OUTWARD answer stays merged (one
// error value, one sentence, ui.go's anti-probe comment untouched). Nothing
// below lifts a cause onto a user-visible face, and TestTicket259R1OutwardAnswer
// is the ruler that goes red if anyone does.
//
// Ruling side note, so the next leg reads a fact and not an inference: these
// cases do NOT decide ticket 259 AC#1 (form a "name the gap" vs form b "make
// the binding check bite"). Cause 3 (misbound) is asserted where that branch
// actually is reachable - the store - and no case below claims the route can or
// cannot produce it. That question, and the wording change it implies, is
// AC#1's and belongs to the orchestrator.

// r1Audit captures what the queue booked, so a denial is asserted at the face
// ticket 259 AC#2 is about: the audit record, not a return value a caller could
// branch on.
type r1Audit struct {
	mu    sync.Mutex
	lines []string
}

func (a *r1Audit) write(format string, args ...any) {
	line := fmt.Sprintf(format, args...)
	a.mu.Lock()
	defer a.mu.Unlock()
	a.lines = append(a.lines, line)
}

func (a *r1Audit) has(want string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, l := range a.lines {
		if strings.Contains(l, want) {
			return true
		}
	}
	return false
}

func (a *r1Audit) dump() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return strings.Join(a.lines, "\n")
}

// r1Card stands up one pending L2 item with one live native grant, without
// going through a display: push + grantNonce are the same two production calls
// PendingApproval makes, so the item is a real card with a real nonce, and every
// refusal below travels through the exported answer surface (NativeAPI /
// DecideFrom*). No goroutine and no clock are involved, which is what makes the
// four readings below deterministic.
func r1Card(t *testing.T, taskID, corr, tool string) (*Gate, *Queue, *qitem, string, *r1Audit) {
	t.Helper()
	a := &r1Audit{}
	g := New(Options{Logf: a.write})
	g.AdmitTextTask(taskID)
	it, nonce := r1CardIn(t, g.Queue(), taskID, corr, tool)
	return g, g.Queue(), it, nonce, a
}

// r1CardIn stands up one more card on an EXISTING queue, so a case can hold two
// live cards in the one funnel at once (the laundering case needs a second
// subject, and two separate gates would not be that).
func r1CardIn(t *testing.T, q *Queue, taskID, corr, tool string) (*qitem, string) {
	t.Helper()
	it, err := q.push(tools.Decision{
		Tool: tool, TaskID: taskID, CorrelationID: corr, Level: 2,
		Args: json.RawMessage(`{"cmd":"259"}`), Paths: []string{"C:/tmp/259-" + corr},
		Reason: "259 fixture card",
	})
	if err != nil {
		t.Fatalf("259 fixture: push: %v", err)
	}
	nonce, err := q.grantNonce(it)
	if err != nil {
		t.Fatalf("259 fixture: grantNonce: %v", err)
	}
	t.Cleanup(func() {
		q.mu.Lock()
		settled := it.state != statePending
		q.mu.Unlock()
		if !settled {
			_ = q.reject(it.Corr, "259 fixture cleanup")
		}
		q.mu.Lock()
		stillPending := it.state == statePending
		q.mu.Unlock()
		if stillPending {
			t.Errorf("259 fixture: card %s was left pending; every case settles its own item", it.Corr)
		}
	})
	return it, nonce
}

// AC#2 cause 1 of 4: no proof was presented at all.
func TestTicket259R1DenialNamesMissingNonce(t *testing.T) {
	g, _, it, _, a := r1Card(t, "task-259-missing", "corr-259-missing", "shell.run")

	err := g.Native().Allow(context.Background(), it.Corr, "")
	if !errors.Is(err, ErrBadGrant) {
		t.Fatalf("AC#2 RED: the missing-proof route returned %v, want the merged ErrBadGrant (the outward face must not split)", err)
	}
	want := "approval: GRANT-DENY corr=" + it.Corr + " tool=shell.run denial=missing-nonce"
	if !a.has(want) {
		t.Fatalf("AC#2 RED: the audit does not name the missing-proof cause (%s)\naudit:\n%s", want, a.dump())
	}
	if !a.has("approval: FORGED-OR-STALE allow rejected corr=" + it.Corr) {
		t.Fatalf("AC#2 RED: the merged FORGED-OR-STALE sentence this line has always written is gone\naudit:\n%s", a.dump())
	}
}

// AC#2 cause 2 of 4: the presented value is not live on this card, which covers
// both "already spent" (asserted at the store, where the spent row is knowable)
// and "never issued" (asserted on the route).
func TestTicket259R1DenialNamesSpentOrNeverLiveNonce(t *testing.T) {
	g, _, it, _, a := r1Card(t, "task-259-spent", "corr-259-spent", "shell.run")

	err := g.Native().Allow(context.Background(), it.Corr, "grant_0000_never_issued")
	if !errors.Is(err, ErrBadGrant) {
		t.Fatalf("AC#2 RED: the never-issued-proof route returned %v, want ErrBadGrant", err)
	}
	want := "approval: GRANT-DENY corr=" + it.Corr + " tool=shell.run denial=spent-or-never-live-nonce"
	if !a.has(want) {
		t.Fatalf("AC#2 RED: the audit does not name the not-live-proof cause (%s)\naudit:\n%s", want, a.dump())
	}

	// The "already spent" half, at the store: the first spend lands, the second
	// attempt at the SAME value must be classified as not-live rather than
	// silently read as a fresh grant.
	s := newGrantStore()
	s.issue("nonce-259-single-use", "digest-259")
	if d := s.spend("nonce-259-single-use", "digest-259"); d != denialNone {
		t.Fatalf("AC#2 RED: the first spend of a live grant was denied (%s); the round trip itself is broken", d.label())
	}
	if d := s.spend("nonce-259-single-use", "digest-259"); d != denialSpentNonce {
		t.Fatalf("AC#2 RED: a spent nonce was classified %s, want the spent-or-never-live cause named apart", d.label())
	}
}

// AC#2 cause 3 of 4: the proof is live on this card but the binding digest it is
// presented against is not the one it was issued under. Asserted at the store,
// because that is the only face where this branch is reachable today (ticket 259
// present-count 1-3: the route passes the same digest it stored). This case
// takes no position on AC#1: it does not claim the route can produce this
// denial, and it does not claim it cannot.
func TestTicket259R1DenialNamesMisbound(t *testing.T) {
	s := newGrantStore()
	minted := bindDigest("corr-259", "task-259", "shell.run", "L2", 7, []byte(`{"cmd":"259"}`))
	s.issue("nonce-259-misbound", minted)

	d := s.spend("nonce-259-misbound", minted+"-tampered")
	if d != denialMisbound {
		t.Fatalf("AC#2 RED: a spend against a mismatched binding was classified %s, want the binding cause named apart", d.label())
	}
	// The retry killer stays the retry killer under the new shape: a refused
	// spend still consumes the nonce.
	if s.live() != 0 {
		t.Fatalf("AC#2 RED: the misbound spend left %d live nonces; refusing must still burn the grant", s.live())
	}
	if d2 := s.spend("nonce-259-misbound", minted); d2 == denialNone {
		t.Fatal("AC#2 RED: the nonce survived a misbound refusal and spent on the retry")
	}
}

// AC#2 cause 4 of 4: the card exists but already left pending, so no proof of
// any quality could be spent on it. The fixture writes the state by hand
// (indexed, not pending) because both production writes of it.state sit in the
// same critical section as the drop that de-indexes the item - which is ticket
// 259 present-count 3's finding, recorded here as a reading rather than used to
// delete a guard. The sibling site the route CAN reach (the losing half of a
// concurrent answer, where deliver reports the item as settled) books the same
// token, and that one is asserted below by calling deliver directly.
func TestTicket259R1DenialNamesCardNotPending(t *testing.T) {
	g, q, it, nonce, a := r1Card(t, "task-259-settled", "corr-259-settled", "shell.run")

	q.mu.Lock()
	it.state = stateAnswered
	q.mu.Unlock()

	err := g.Native().Allow(context.Background(), it.Corr, nonce)
	if !errors.Is(err, ErrNotPending) {
		t.Fatalf("AC#2 RED: an answer for a settled card returned %v, want ErrNotPending", err)
	}
	want := "approval: GRANT-DENY corr=" + it.Corr + " tool=shell.run denial=card-not-pending"
	if !a.has(want) {
		t.Fatalf("AC#2 RED: the audit does not name the settled-card cause (%s)\naudit:\n%s", want, a.dump())
	}
	if it.grants.live() == 0 {
		t.Fatal("AC#2 RED: the settled-card refusal burned the live grant; that branch never reaches the store")
	}
	// The route-reachable sibling: deliver refuses to settle an item that is
	// already settled, and allowScoped books that as the same cause.
	if q.deliver(it, answer{a: tools.AnswerAllow, why: "259: pre-settle for the deliver branch"}) {
		t.Fatal("AC#2 RED: deliver settled an item that had already left pending; the sibling site cannot be reached")
	}
	q.mu.Lock()
	q.dropLocked(it, false)
	q.mu.Unlock()
	if n := q.Depth(); n != 0 {
		t.Fatalf("AC#2 fixture: %d pending item(s) after cleanup", n)
	}
}

// The anti-folding ruler: four causes, four DIFFERENT names. If a later edit
// collapses two labels into one token - or routes two causes through the same
// branch - this goes red without any other case having to change.
func TestTicket259R1FourDenialLabelsArePairwiseDistinct(t *testing.T) {
	cases := []grantDenial{denialMissingNonce, denialSpentNonce, denialMisbound, denialNotPending}
	seen := map[string]grantDenial{}
	for _, d := range cases {
		lab := d.label()
		if lab == "" || lab == denialNone.label() {
			t.Fatalf("AC#2 RED: denial %d has no name of its own (label=%q)", d, lab)
		}
		if other, dup := seen[lab]; dup {
			t.Errorf("AC#2 RED: denials %d and %d share the label %q; four causes must not fold into fewer names", other, d, lab)
		}
		seen[lab] = d
	}
	if len(seen) != len(cases) {
		t.Errorf("AC#2 RED: %d causes produced %d distinct labels; the four sentences are the deliverable", len(cases), len(seen))
	}
}

// The merge the orchestrator ruled STAYS merged, pinned from the inside: every
// grant denial leaves this package as one error value and one sentence, and no
// denial token is reachable from the outward text. If a leg lifts the four
// causes onto a UI face, or reuses the permission-denied wording for them, this
// case is the one that says so.
func TestTicket259R1OutwardAnswerStaysMergedAcrossDenials(t *testing.T) {
	const merged = "approval: 原生令牌无效（缺失/已用/与本次请求不绑定）"
	if got := ErrBadGrant.Error(); got != merged {
		t.Fatalf("AC#2 RED: the merged outward sentence moved (%q); ticket 259 keeps it merged, the split is internal-only", got)
	}
	g, _, it, _, a := r1Card(t, "task-259-merge", "corr-259-merge", "shell.run")
	if err := g.Native().Allow(context.Background(), it.Corr, ""); !errors.Is(err, ErrBadGrant) {
		t.Fatalf("AC#2 RED: missing proof returned %v, want ErrBadGrant", err)
	}
	if err := g.Native().Allow(context.Background(), it.Corr, "grant_nope"); !errors.Is(err, ErrBadGrant) {
		t.Fatalf("AC#2 RED: never-live proof returned %v, want ErrBadGrant", err)
	}
	for _, lab := range []string{denialMissingNonce.label(), denialSpentNonce.label(), denialMisbound.label(), denialNotPending.label()} {
		if strings.Contains(merged, lab) {
			t.Errorf("AC#2 RED: the outward sentence now carries the internal cause %q; the merge is the anti-probe property", lab)
		}
	}
	// The internal attribution exists, and it exists ONLY in the audit face.
	if !a.has("denial=missing-nonce") || !a.has("denial=spent-or-never-live-nonce") {
		t.Fatalf("AC#2 RED: the audit stopped naming the two causes it merged outward\naudit:\n%s", a.dump())
	}
	// Two things, two sentences: a grant denial must not borrow the wording of
	// the refusals that already travel this package (panel-route, D47 admission,
	// unloaded channel). Reusing one of those for a fourth cause would fold the
	// distinction back into an existing bucket instead of naming it.
	borrowed := []string{
		ErrPanelAllow.Error(),
		ErrNotAdmitted.Error(),
		ErrChannelUnavailable.Error(),
		ErrNotPending.Error(),
		ErrUnknownCorrelation.Error(),
	}
	for _, d := range []grantDenial{denialMissingNonce, denialSpentNonce, denialMisbound, denialNotPending} {
		for _, sentence := range borrowed {
			if strings.Contains(sentence, d.label()) {
				t.Errorf("AC#2 RED: denial %q reuses an existing refusal sentence (%q); each cause carries its own words", d.label(), sentence)
			}
		}
	}
}
