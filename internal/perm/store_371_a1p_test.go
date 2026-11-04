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
// readings are stable (R20/M4). The D and F readings pin a LIVE HOLE, so the day
// ticket 371's AC#2' lands - "a card that was never answered must not apply" -
// the D/F assertions below flip to applied=0 with a refusal naming the unasked
// card. Flipping that assertion is part of that fix, not a regression here.
//
// On the 372 blob (before the gate moved back to the target) A read
// refused-confirm instead of applied; on the 367 blob B read confirm_calls=0
// because the direction rule let the middle tier past without a card.

import (
	"context"
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
}

// real is a limb that actually asks: it raises the card (asked++) and returns
// nil only when that card came back allow - the production limb's contract.
func (l *a1pLimb) real(context.Context, Switch) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.invoked++
	l.asked++
	return l.answer
}

// fake is NON-NIL and NEVER ASKS: it is entered, records that it was entered,
// and returns nil without a card ever existing. AC#0'②'s new shape.
func (l *a1pLimb) fake(context.Context, Switch) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.invoked++
	return nil
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
// ask_high_risk with a REAL, asking limb attached. R20/M4 says this target
// costs no card, so the reading must be invoked=0 / asked=0 / applied.
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
		t.Fatalf("ask_every_step -> ask_high_risk was refused (%v); R20/M4 cards only auto_approve", setErr)
	}
	if invoked != 0 || asked != 0 {
		t.Errorf("the middle tier cost %d invocations / %d asks, want 0/0 (R20/M4)", invoked, asked)
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

// Test371A1PAC0PrimeFakeLimbIsSilentlyApplied is AC#0'② and ④: the LIVE HOLE.
// auto_approve, a non-nil limb that never asks, and the gate cannot tell it
// apart from the answering limb above - the switch applies with asked=0.
//
// This is not ticket 375's hole and not 372's: those are the "cannot ask" side
// (they fire on nil). This limb CAN ask; it simply is not made to, and nothing
// in Set asks whether it did.
func Test371A1PAC0PrimeFakeLimbIsSilentlyApplied(t *testing.T) {
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
	// Gate-side reading TODAY: entered, never asked, applied anyway.
	if setErr != nil {
		t.Fatalf("the fake limb was not silently accepted (%v); the hole AC#0'② records is gone - "+
			"flip this case and the evidence's D row instead of deleting it", setErr)
	}
	if !ok || last.Result != ResultApplied {
		t.Errorf("audit result = %q, want %q (this is the hole: applied without an ask)", last.Result, ResultApplied)
	}
	if persisted != risk.ModeAutoApprove.String() || saves != 1 {
		t.Errorf("persisted = %q / saves=%d, want %q / 1 (the loosest tier landed on disk)", persisted, saves, risk.ModeAutoApprove)
	}
}

// Test371A1PAC0PrimeTheAutoApproveGateHasTheSameHole is AC#0'④: the judgement,
// not a second measurement. It runs the answering limb and the never-asking
// limb against the same target and asserts that the gate's DECISION - the pair
// (error is nil, result, persisted) - is identical while the ask counts differ.
// That equality IS the hole: the gate's verdict is a function of non-nil-ness,
// which both limbs satisfy, so "was answered" contributes nothing to it.
func Test371A1PAC0PrimeTheAutoApproveGateHasTheSameHole(t *testing.T) {
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
	if realOut.refused != fakeOut.refused || realOut.result != fakeOut.result ||
		realOut.persisted != fakeOut.persisted || realOut.saves != fakeOut.saves {
		t.Errorf("the gate split the two limbs: answered=%+v never-asking=%+v. If that is ticket "+
			"371's fix landing, make this the positive judgement instead of deleting it", realOut, fakeOut)
	}
	if realOut.refused {
		t.Errorf("the answered limb was refused; AA with an answering card must apply")
	}
}
