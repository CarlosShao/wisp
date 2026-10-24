package perm

// Tickets 367 / 372: what Set's gate judges, and what stayed from 367.
//
// 367's rule - every widening refused without a confirmation channel - is no
// longer this package's rule. Ticket 372 reverted the gate to the TARGET because
// the mainstream survey found no harness that gates on "wider"
// (docs/evidence/s1/m4-survey.md:96) and because R20/M4 cards only the switch INTO
// auto_approve. What the revert KEPT from 367 is the discipline: one predicate
// (perm.RequiresCard, direction.go) called by both transition points, instead of
// the two private judges 367-a1 measured.
//
// The expectations below are the post-372 ones, which is what makes them teeth:
// on the 367 blob the first table's ask_every_step -> ask_high_risk row reads
// refused-confirm and fails; on the pre-367 blob both auto_approve rows read
// applied and fail there.

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// sixPairs is the complete set of switches between three distinct modes: three
// widening and three narrowing. The total is asserted, not assumed, so a fourth
// mode added later cannot silently make this table partial.
func sixPairs() []struct {
	from, to risk.Mode
	widens   bool
} {
	ladder := []risk.Mode{risk.ModeAskEveryStep, risk.ModeAskHighRisk, risk.ModeAutoApprove}
	var out []struct {
		from, to risk.Mode
		widens   bool
	}
	for i := range ladder {
		for j := range ladder {
			if i == j {
				continue
			}
			out = append(out, struct {
				from, to risk.Mode
				widens   bool
			}{ladder[i], ladder[j], j > i})
		}
	}
	return out
}

// TestSetGatesOnlyTheAutoApproveTargetAndEveryNarrowingApplies is ticket 372's
// acceptance reading as a table: with NO confirmation channel wired, exactly the
// two switches whose TARGET is auto_approve are refused, and the other four -
// the middle-tier widening included - apply.
//
// The direction columns are not decoration. The table asserts 3 widening / 3
// narrowing, so a predicate that went back to "refuse every widening" would still
// refuse 2 of the 3 widenings and the count alone would not catch it; the
// ask_every_step -> ask_high_risk row is the line that does, and the three
// narrowing rows keep a gate that refuses everything from hiding behind it.
func TestSetGatesOnlyTheAutoApproveTargetAndEveryNarrowingApplies(t *testing.T) {
	ctx := context.Background()
	pairs := sixPairs()
	if len(pairs) != 6 {
		t.Fatalf("pair table holds %d switches, want 6 (3 modes, no self-switch)", len(pairs))
	}
	widening, narrowing, refused, applied := 0, 0, 0, 0
	for _, p := range pairs {
		s, mc, rec := newTestStore(t, p.from.String(), nil)
		before, _ := mc.saved()
		err := s.Set(ctx, p.to, "t372", "owner")
		last, ok := s.Last()
		if !ok {
			t.Fatalf("%s -> %s: the switch left no audit record at all", p.from, p.to)
		}
		if p.widens {
			widening++
		} else {
			narrowing++
		}

		// Spelled literally on purpose: a test that asked RequiresCard would agree
		// with a broken predicate instead of judging it.
		if p.to == risk.ModeAutoApprove {
			refused++
			if err == nil {
				t.Errorf("%s -> %s was APPLIED with no confirmation channel; the auto_approve target must be refused (R20/M4)",
					p.from, p.to)
			}
			if last.Result != ResultRefused {
				t.Errorf("%s -> %s audit result = %q, want %q", p.from, p.to, last.Result, ResultRefused)
			}
			if mode, saves := mc.saved(); mode != before || saves != 0 {
				t.Errorf("%s -> %s persisted anyway: mode=%q saves=%d, want %q/0",
					p.from, p.to, mode, saves, before)
			}
			// The refusal has to name both ends, or an operator reading the log
			// cannot tell which switch was attempted.
			if err != nil && (!strings.Contains(err.Error(), p.from.String()) ||
				!strings.Contains(err.Error(), p.to.String())) {
				t.Errorf("%s -> %s refusal does not name both modes: %v", p.from, p.to, err)
			}
			if rec.count(string(ResultRefused)) != 1 {
				t.Errorf("%s -> %s: %d refusal lines in the audit, want 1: %v",
					p.from, p.to, rec.count(string(ResultRefused)), rec.all())
			}
		} else {
			applied++
			if err != nil {
				t.Errorf("%s -> %s was REFUSED (%v); R20/M4 cards only the auto_approve target", p.from, p.to, err)
			}
			if last.Result != ResultApplied {
				t.Errorf("%s -> %s audit result = %q, want %q", p.from, p.to, last.Result, ResultApplied)
			}
			if mode, saves := mc.saved(); mode != p.to.String() || saves != 1 {
				t.Errorf("%s -> %s persisted mode=%q saves=%d, want %q/1",
					p.from, p.to, mode, saves, p.to)
			}
			if rec.count(string(ResultRefused)) != 0 {
				t.Errorf("%s -> %s wrote a refusal line for a target M4 does not card: %v", p.from, p.to, rec.all())
			}
		}
		// Both halves read the same effective档: with no channel the read side is
		// downgraded to ask_every_step whatever the file now says (372-r2).
		if got := s.PermissionMode(); got != risk.ModeAskEveryStep {
			t.Errorf("%s -> %s: the effective read is %s, want the downgrade to ask_every_step (no channel)", p.from, p.to, got)
		}
	}
	if widening != 3 || narrowing != 3 {
		t.Fatalf("table classified %d widening / %d narrowing, want 3/3", widening, narrowing)
	}
	if refused != 2 || applied != 4 {
		t.Fatalf("table produced %d refused / %d applied, want 2/4 (only the two auto_approve targets are gated)", refused, applied)
	}
}

// TestSetPassesTheMiddleTargetAndCardsOnlyTheAutoTarget is the other half of the
// rule: a channel is what makes the auto_approve target reachable, and the
// non-auto step costs no ask at all (R20/M4). The first switch would also pass
// without a channel since 372 - it is here to keep the asymmetry reading: zero
// cards for ask_high_risk, exactly one for auto_approve.
func TestSetPassesTheMiddleTargetAndCardsOnlyTheAutoTarget(t *testing.T) {
	ctx := context.Background()
	spy := &confirmSpy{}
	s, mc, rec := newTestStore(t, risk.ModeAskEveryStepName, spy.Confirm)
	if err := s.Set(ctx, risk.ModeAskHighRisk, "t372", "owner"); err != nil {
		t.Fatalf("ask_every_step -> ask_high_risk refused although a confirmation channel is attached: %v", err)
	}
	if got := s.PermissionMode(); got != risk.ModeAskHighRisk {
		t.Errorf("mode = %s after an allowed widening, want ask_high_risk", got)
	}
	if mode, saves := mc.saved(); mode != risk.ModeAskHighRiskName || saves != 1 {
		t.Errorf("persisted mode=%q saves=%d, want ask_high_risk/1", mode, saves)
	}
	if spy.n() != 0 {
		t.Errorf("the non-auto widening cost %d confirmations; R20/M4 says only auto_approve does", spy.n())
	}
	if rec.count(string(ResultApplied)) != 1 {
		t.Errorf("the allowed widening is not in the audit as applied: %v", rec.all())
	}

	// And the own-target case still costs exactly one answered card, which is
	// what "the gate did not replace R20/M4" means.
	if err := s.Set(ctx, risk.ModeAutoApprove, "t372", "owner"); err != nil {
		t.Fatalf("ask_high_risk -> auto_approve refused with a channel attached: %v", err)
	}
	if spy.n() != 1 {
		t.Errorf("auto_approve cost %d confirmations, want exactly 1 (R20/M4)", spy.n())
	}
	if len(spy.calls) == 1 && (spy.calls[0].From != risk.ModeAskHighRisk || spy.calls[0].To != risk.ModeAutoApprove) {
		t.Errorf("the card was raised about %s -> %s, want ask_high_risk -> auto_approve",
			spy.calls[0].From, spy.calls[0].To)
	}
}

// TestSetRefusesAutoApproveOnAnUnreadableStoredMode pins the invariant the gate
// leans on: PermissionMode() fail-closes to ask_every_step, so an unreadable
// config is still a *different* mode from auto_approve and the target gate refuses
// the switch. Without the gate's nil check this path would dereference a nil
// channel instead of refusing (the panic shot in 372-r2's evidence).
func TestSetRefusesAutoApproveOnAnUnreadableStoredMode(t *testing.T) {
	s, mc, _ := newTestStore(t, "not_a_mode", nil)
	if got := s.PermissionMode(); got != risk.ModeAskEveryStep {
		t.Fatalf("an unreadable stored mode reads as %s, want the strict default ask_every_step", got)
	}
	err := s.Set(context.Background(), risk.ModeAutoApprove, "t372", "owner")
	if err == nil {
		t.Fatal("auto_approve was reached from an unreadable stored mode with no channel")
	}
	if _, saves := mc.saved(); saves != 0 {
		t.Errorf("state moved under a refusal: saves=%d", saves)
	}
}

// TestIsWideningOnTheDocumentedLadder keeps the ladder-order witness alive: since
// ticket 372 this function is NOT a gate, but the order it reads is still the
// vocabulary the audit words "widened"/"narrowed" rest on. internal/panel drives
// the SAME function over the SAME names, so a reorder breaks both ends rather
// than silently flipping one.
func TestIsWideningOnTheDocumentedLadder(t *testing.T) {
	ladder := []risk.Mode{risk.ModeAskEveryStep, risk.ModeAskHighRisk, risk.ModeAutoApprove}
	names := risk.ModeNames()
	if len(names) != len(ladder) {
		t.Fatalf("risk.ModeNames() holds %d names, this predicate's ladder holds %d", len(names), len(ladder))
	}
	for i := 1; i < len(ladder); i++ {
		if !IsWidening(ladder[i-1], ladder[i]) {
			t.Errorf("%s -> %s must read as widening", ladder[i-1], ladder[i])
		}
		if IsWidening(ladder[i], ladder[i-1]) {
			t.Errorf("%s -> %s must read as narrowing", ladder[i], ladder[i-1])
		}
	}
	for _, m := range ladder {
		if IsWidening(m, m) {
			t.Errorf("%s -> %s must not read as widening (a no-op is not a loosening)", m, m)
		}
	}
}

// TestSetConfirmRefusalOnTheAutoTargetIsClassifiedAsRefused keeps the two failure
// modes distinguishable in the audit: "the operator said no" is not "there was no
// way to ask". The nil-channel case above is ResultRefused at the gate; this is
// the answered-no case on the one target that asks.
func TestSetConfirmRefusalOnTheAutoTargetIsClassifiedAsRefused(t *testing.T) {
	spy := &confirmSpy{answer: errors.New("用户拒绝")}
	s, mc, rec := newTestStore(t, risk.ModeAskEveryStepName, spy.Confirm)
	err := s.Set(context.Background(), risk.ModeAutoApprove, "t372", "owner")
	if err == nil {
		t.Fatal("Set succeeded although the confirmation was refused")
	}
	if last, _ := s.Last(); last.Result != ResultRefused {
		t.Errorf("audit result = %q, want %q", last.Result, ResultRefused)
	}
	if _, saves := mc.saved(); saves != 0 {
		t.Errorf("a refused confirmation still persisted (%d writes)", saves)
	}
	if rec.count(string(ResultRefused)) != 1 {
		t.Errorf("refusal not audited: %v", rec.all())
	}
}
