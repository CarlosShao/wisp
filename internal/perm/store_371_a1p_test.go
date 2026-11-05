package perm

// Leg 371-a1p: the READ-ONLY evidence ruler for ticket 371's AC#0 ②③ and
// AC#0' ②④. Nothing here is a fix; every case measures what Set does TODAY, at
// HEAD 786a13b22.
//
// THE RULER. One limb type carries two counters, because "confirm_calls" in the
// ticket has to mean "a card was put in front of an operator and that operator
// answered allow", not "a Go function was entered":
//
//	invoked  the Store called this ConfirmFunc (the gate "reached the channel")
//	asked    the limb actually raised a card and got an allow back (an ask)
//
// a1pLimb.answer is the real shape of the production leg: cmd/wisp/run.go's
// confirmModeSwitchVia returns nil only after approval.Gate.PendingApproval
// answered tools.AnswerAllow (run.go:1126-1132), i.e. a nil return is supposed
// to mean "someone answered allow". NewTestStore's in-package stub for that
// shape is run_mode101_test.go's confirm closure ("stands for the native side
// clicked allow once", :225-228) and panel_mode_confirm_114_windows_test.go's
// countingPermConfirm. The fake limb below is the new shape AC#0'② names: it is
// non-nil, it returns nil, and it never raises a card at all.
//
// EXPECTATION FRESHNESS, stated the way store_367_test.go states it: the A and B
// readings are stable (R20/M4). The D and F readings pinned a LIVE HOLE, and the
// note they carried said what would happen when ticket 371's AC#2' landed - "a
// card that was never answered must not apply" - so the D/F assertions below were
// FLIPPED by leg 371-r1p to applied=0 / refused with the unanswered-card reason.
// That flip is the fix landing, not a regression here; the readings themselves
// (invoked=1 / asked=0 for the fake limb) are unchanged, because the limb-side
// facts are what any fix has to keep true.
//
// On the 372 blob A read refused-confirm instead of applied (367's direction rule
// was refusing the middle tier); on the 367 blob B read confirm_calls=0 because
// the middle tier was free; since R9 both middle-tier rows read one card, because
// every loosening costs one.

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// a1pLimb is one confirmation channel with the two readings the ticket needs.
type a1pLimb struct {
	mu      sync.Mutex
	invoked int
	asked   int
	answer  error
	corr    string
}

// real is a limb that actually asks: it raises the card (asked++) and returns an
// answered-allow credential only when that card came back allow - the production
// limb's contract (ticket 371: the credential, not the nil error, is what says so).
func (l *a1pLimb) real(context.Context, Switch) (Card, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.invoked++
	l.asked++
	if l.answer != nil {
		return Card{}, l.answer
	}
	return AnsweredCard(l.corr), nil
}

// fake is NON-NIL and NEVER ASKS: it is entered, records that it was entered, and
// returns the ZERO credential without a card ever existing. AC#0'②'s new shape -
// the one ticket 371's gate has to refuse.
func (l *a1pLimb) fake(context.Context, Switch) (Card, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.invoked++
	return Card{}, nil
}

func (l *a1pLimb) counts() (invoked, asked int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.invoked, l.asked
}

// a1pReading is one measurement, logged in one pipe-delimited line so the
// evidence file can quote the command's output verbatim.
func a1pReading(t *testing.T, tag string, from, to risk.Mode, invoked, asked int, setErr error, last Switch, ok bool, persisted string, saves int) {
	t.Helper()
	result := "<no-record>"
	if ok {
		result = string(last.Result)
	}
	t.Logf("371A1P|%s|%s->%s|invoked=%d|asked=%d|set_err=%v|last_result=%s|persisted=%s|saves=%d",
		tag, from, to, invoked, asked, setErr, result, persisted, saves)
}

// Test371A1PAC0RealLimbOnTheMiddleTier measures AC#0②: ask_every_step ->
// ask_high_risk with a REAL, asking limb attached. On 786a13b22 this reading was
// invoked=0 / asked=0 because R20/M4 carded only the auto_approve target. R9
// (owner 2026-11-05, ledger A1541 item 4) cards EVERY loosening, so the same row
// now reads invoked=1 / asked=1 and still applies - the middle tier moved from
// "no card" to "one card", not to "refused".
func Test371A1PAC0RealLimbOnTheMiddleTier(t *testing.T) {
	ctx := context.Background()
	limb := &a1pLimb{}
	s, mc, _ := newTestStore(t, risk.ModeAskEveryStep.String(), limb.real)

	setErr := s.Set(ctx, risk.ModeAskHighRisk, "t371a1p", "owner")
	invoked, asked := limb.counts()
	last, ok := s.Last()
	persisted, saves := mc.saved()
	a1pReading(t, "AC0-2 AE->AH real limb", risk.ModeAskEveryStep, risk.ModeAskHighRisk,
		invoked, asked, setErr, last, ok, persisted, saves)

	if setErr != nil {
		t.Fatalf("ask_every_step -> ask_high_risk was refused with an answering limb attached: %v", setErr)
	}
	if invoked != 1 || asked != 1 {
		t.Errorf("the middle tier cost %d invocations / %d asks, want 1/1 (R9: every loosening shows one card)",
			invoked, asked)
	}
	if !ok || last.Result != ResultApplied {
		t.Errorf("audit result = %q, want %q - the switch applied", last.Result, ResultApplied)
	}
	if persisted != risk.ModeAskHighRisk.String() || saves != 1 {
		t.Errorf("persisted = %q / saves=%d, want %q / 1 (applied)", persisted, saves, risk.ModeAskHighRisk)
	}
}

// Test371A1PAC0PositiveControlOneCardOnAutoApprove is AC#0③: the POSITIVE
// CONTROL. Same limb, same store, only the target changes to auto_approve. If
// this reads asked=0 the ruler cannot see an ask at all and AC#0②'s zero is
// worthless; it must read exactly one.
func Test371A1PAC0PositiveControlOneCardOnAutoApprove(t *testing.T) {
	ctx := context.Background()
	limb := &a1pLimb{}
	s, mc, _ := newTestStore(t, risk.ModeAskEveryStep.String(), limb.real)

	setErr := s.Set(ctx, risk.ModeAutoApprove, "t371a1p", "owner")
	invoked, asked := limb.counts()
	last, ok := s.Last()
	persisted, saves := mc.saved()
	a1pReading(t, "AC0-3 AE->AA real limb (positive control)", risk.ModeAskEveryStep, risk.ModeAutoApprove,
		invoked, asked, setErr, last, ok, persisted, saves)

	if setErr != nil {
		t.Fatalf("ask_every_step -> auto_approve with an answering limb was refused: %v", setErr)
	}
	if invoked != 1 || asked != 1 {
		t.Fatalf("the auto_approve target read invocations=%d asks=%d, want 1/1: this is the "+
			"positive control, and a ruler that reads 0 here cannot certify any other zero", invoked, asked)
	}
	if !ok || last.Result != ResultApplied {
		t.Errorf("audit result = %q, want %q", last.Result, ResultApplied)
	}
	if persisted != risk.ModeAutoApprove.String() || saves != 1 {
		t.Errorf("persisted = %q / saves=%d, want %q / 1", persisted, saves, risk.ModeAutoApprove)
	}
}

// Test371A1PAC0PrimeFakeLimbIsRefusedNotSilentlyApplied is AC#0'② and ④. It was
// written as Test371A1PAC0PrimeFakeLimbIsSilentlyApplied on the 786a13b22 blob,
// where this reading WAS the hole: auto_approve, a non-nil limb that never asks,
// and the gate could not tell it apart from the answering limb above - the switch
// applied with asked=0.
//
// Leg 371-r1p flipped it, which is the fix landing: the same fake limb is now
// entered once, never asks, and the switch is REFUSED with a reason that names the
// unanswered card. This is not ticket 375's hole and not 372's: those are the
// "cannot ask" side (they fire on nil). This limb CAN ask; it simply is not made
// to, and the gate no longer takes that for an answer.
func Test371A1PAC0PrimeFakeLimbIsRefusedNotSilentlyApplied(t *testing.T) {
	ctx := context.Background()
	limb := &a1pLimb{}
	s, mc, _ := newTestStore(t, risk.ModeAskEveryStep.String(), limb.fake)

	setErr := s.Set(ctx, risk.ModeAutoApprove, "t371a1p", "owner")
	invoked, asked := limb.counts()
	last, ok := s.Last()
	persisted, saves := mc.saved()
	a1pReading(t, "AC0p-2 AE->AA fake limb (the hole)", risk.ModeAskEveryStep, risk.ModeAutoApprove,
		invoked, asked, setErr, last, ok, persisted, saves)

	// Limb-side facts, stable under any fix: the fake was entered once and
	// never raised a card.
	if invoked != 1 {
		t.Fatalf("the fake limb was entered %d times, want 1 - Set does call the leg for this target", invoked)
	}
	if asked != 0 {
		t.Fatalf("the fake limb reported %d asks; by construction it never asks", asked)
	}
	// Gate-side reading SINCE 371: entered, never asked, refused - and the reason
	// is its own, not the nil-channel one.
	if setErr == nil {
		t.Fatalf("the fake limb was silently accepted; ticket 371's gate is gone - "+
			"the never-answering limb must be refused (%v)", setErr)
	}
	if !errors.Is(setErr, ErrCardNotAnswered) {
		t.Errorf("refusal = %v, want it to wrap ErrCardNotAnswered (a card was never answered, "+
			"which is not the nil-channel refusal)", setErr)
	}
	if !strings.Contains(setErr.Error(), "card was not answered") {
		t.Errorf("refusal %q does not name the unanswered card", setErr.Error())
	}
	if last.Result != ResultCardNotAnswered {
		t.Errorf("audit result = %q, want %q (not %q, which is the nil-channel verdict)",
			last.Result, ResultCardNotAnswered, ResultRefused)
	}
	if last.Asked || last.Card != "" {
		t.Errorf("audit asked=%v card=%q, want false/\"\": nothing was answered", last.Asked, last.Card)
	}
	if persisted != risk.ModeAskEveryStep.String() || saves != 0 {
		t.Errorf("persisted = %q / saves=%d, want %q / 0 (the loosest tier must not land)",
			persisted, saves, risk.ModeAskEveryStep)
	}
}

// Test371A1PAC0PrimeTheAutoApproveGateSplitsTheTwoLimbs is AC#0'④. It was written
// as ...TheAutoApproveGateHasTheSameHole: it asserted that the gate's DECISION was
// identical for the answering and the never-asking limb, and that equality WAS the
// hole. Leg 371-r1p flipped the judgement - the gate must now split them, and the
// only reading that may stay equal is the limb-side one.
//
// It runs both limbs against the same target and asserts (a) the asks differ (1
// vs 0, unchanged), (b) the answered limb applies, and (c) the never-asking limb
// is refused with the unanswered-card result. If a later change lets the fake
// through again, (c) reddens with both outcomes printed.
func Test371A1PAC0PrimeTheAutoApproveGateSplitsTheTwoLimbs(t *testing.T) {
	ctx := context.Background()

	type outcome struct {
		refused   bool
		result    Result
		persisted string
		saves     int
		asked     int
	}
	run := func(t *testing.T, limb ConfirmFunc) outcome {
		t.Helper()
		s, mc, _ := newTestStore(t, risk.ModeAskEveryStep.String(), limb)
		err := s.Set(ctx, risk.ModeAutoApprove, "t371a1p", "owner")
		last, _ := s.Last()
		persisted, saves := mc.saved()
		return outcome{refused: err != nil, result: last.Result, persisted: persisted, saves: saves}
	}

	answering := &a1pLimb{}
	neverAsks := &a1pLimb{}
	realOut := run(t, answering.real)
	fakeOut := run(t, neverAsks.fake)
	_, realAsked := answering.counts()
	_, fakeAsked := neverAsks.counts()
	realOut.asked, fakeOut.asked = realAsked, fakeAsked

	t.Logf("371A1P|AC0p-4 answered_limb|refused=%v|result=%s|persisted=%s|saves=%d|asked=%d",
		realOut.refused, realOut.result, realOut.persisted, realOut.saves, realOut.asked)
	t.Logf("371A1P|AC0p-4 never-asking_limb|refused=%v|result=%s|persisted=%s|saves=%d|asked=%d",
		fakeOut.refused, fakeOut.result, fakeOut.persisted, fakeOut.saves, fakeOut.asked)

	if realOut.asked != 1 || fakeOut.asked != 0 {
		t.Fatalf("the two limbs did not differ where they must: answered=%d never-asking=%d asks",
			realOut.asked, fakeOut.asked)
	}
	if realOut.refused {
		t.Errorf("the answered limb was refused; AA with an answering card must apply (got %+v)", realOut)
	}
	if realOut.result != ResultApplied {
		t.Errorf("the answered limb's result = %q, want %q", realOut.result, ResultApplied)
	}
	if !fakeOut.refused || fakeOut.result != ResultCardNotAnswered {
		t.Errorf("the NEVER-ANSWERING limb was not refused as unanswered: got %+v, want refused with %q. "+
			"If the two limbs read the same verdict again, ticket 371's gate is gone", fakeOut,
			ResultCardNotAnswered)
	}
	if realOut.persisted == fakeOut.persisted && realOut.saves == fakeOut.saves {
		t.Errorf("both limbs left the same persisted state (%q/%d): the gate did not split them",
			realOut.persisted, realOut.saves)
	}
}
