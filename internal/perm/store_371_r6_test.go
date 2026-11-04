package perm

// Ticket 371, R6 (the ticket's F4, item 〔己〕): the GateAsk branch must REFUSE
// a nil confirmation channel instead of panicking.
//
// THE OLD INVARIANT. Set judged the target and the channel once
// (JudgeCard(to, s.confirm != nil)) and then, in the GateAsk branch, called
// s.confirm directly. That was safe only as long as the predicate kept refusing
// every auto_approve switch whose channel was nil; the comment said so and the
// panic shot in docs/evidence/s1/372-r2-impl.md used the panic as the witness.
// R6 keeps the refusal but not the crash: a broken predicate is an accident in
// the gate, and the worst user-visible outcome of an accident while someone is
// changing a security posture is a refusal, never a dead process.
//
// HOW THIS TEST REACHES THE BRANCH. It enters GateAsk with a REAL limb attached
// (the predicate answers "ask"), then drops the channel between the judgement
// and the call. Set stamps the switch (clock call 1) before the gate and starts
// the confirmation timer (clock call 2) inside GateAsk after the gate has
// decided; dropping the channel on call 2 is therefore a direct unit that hands
// the GateAsk path a nil confirm while JudgeCard still saw a non-nil one. On the
// pre-R6 code that call dereferences a nil func and the test fails with a panic;
// on R6 it must end in ResultRefused with the inner-guard detail. The same guard
// is what refuses a predicate that lets a nil channel through outright - see the
// overlay readings in docs/evidence/s1/371-r7-boundary.md.
//
// The test also refuses to pass vacuously: it asserts the hook was reached
// exactly twice, so a future reordering of Set that skips the GateAsk timer
// fails loudly instead of silently testing nothing.

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/risk"
)

func Test371R6TheGateAskPathRefusesANilChannelInsteadOfPanicking(t *testing.T) {
	ctx := context.Background()
	spy := &confirmSpy{}
	s, mc, rec := newTestStore(t, risk.ModeAskEveryStep.String(), spy.Confirm)

	base := s.now
	calls := 0
	s.now = func() time.Time {
		calls++
		if calls == 2 {
			// The gate has already said GateAsk with a non-nil channel; this is
			// the instant before the limb would be called.
			s.confirm = nil
		}
		return base()
	}

	err := s.Set(ctx, risk.ModeAutoApprove, "t371-r6", "owner")

	if calls != 2 {
		t.Fatalf("the clock hook saw %d calls, want 2: the test never reached the GateAsk timer, "+
			"so it proves nothing about the branch it names", calls)
	}
	if err == nil {
		t.Fatal("auto_approve was applied although the confirmation channel was nil at the call - " +
			"the inner nil guard is gone and a broken predicate would apply the loosest tier")
	}
	last, ok := s.Last()
	if !ok {
		t.Fatal("the refused switch left no audit record")
	}
	if last.Result != ResultRefused {
		t.Errorf("audit result = %q, want %q (the inner nil guard's own refusal, not a panic "+
			"and not %q, which is the answered-card verdict)", last.Result, ResultRefused, ResultCardNotAnswered)
	}
	if !strings.Contains(last.Detail, "内层 nil 兜底") {
		t.Errorf("refusal detail = %q, want it to name the inner nil guard: an operator must be "+
			"able to tell a broken permission predicate from an ordinary no-channel assembly", last.Detail)
	}
	// The ordinary GateNoChannel wording belongs to the other branch. If this
	// refusal reads like that one, the test never reached the GateAsk path.
	if strings.Contains(last.Detail, "没有接入确认通道") {
		t.Errorf("refusal detail = %q is the GateNoChannel wording; the test did not reach GateAsk", last.Detail)
	}
	if spy.n() != 0 {
		t.Errorf("the confirmation limb was called %d time(s) after the channel was dropped, want 0", spy.n())
	}
	if mode, saves := mc.saved(); mode != risk.ModeAskEveryStep.String() || saves != 0 {
		t.Errorf("persisted = %q / saves=%d, want %q / 0: a refused inner guard must not move the file",
			mode, saves, risk.ModeAskEveryStep)
	}
	if rec.count(string(ResultRefused)) != 1 {
		t.Errorf("the inner-guard refusal was not audited exactly once: %v", rec.all())
	}
}
