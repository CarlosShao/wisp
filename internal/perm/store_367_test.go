package perm

// Tickets 367 / 372 / R9: what Set's gate judges.
//
// THREE RULINGS IN A ROW, and this file is the seam between them:
//
//   - ticket 367 made the DIRECTION the gate (every widening refused without a
//     confirmation channel), which also made ask_every_step -> ask_high_risk
//     unreachable on a legless host;
//   - ticket 372 reverted the judgement to the TARGET ("only auto_approve"),
//     because the mainstream survey found no harness that REFUSES a switch for
//     being wider (docs/evidence/s1/m4-survey.md:96), and kept 367's discipline:
//     one predicate (perm.RequiresCard, direction.go) called by every transition
//     point;
//   - R9 (owner 2026-11-05, ledger A1541 item 4: "every loosening shows a
//     confirmation card") put the DIRECTION back as the COST rule. It does not
//     contradict the survey: the survey refuted refusing a widening, and R9 makes
//     a widening ask for the card the config-file channel has asked for since D36.
//     A host with no card surface therefore cannot loosen at all, which is the
//     ruling's accepted cost (docs/evidence/s1/R9-impl.md).
//
// The expectations below are R9's. On the 367 blob they differ in the middle-tier
// row only for the widening arm (367 refused it, R9 asks a card for it); on the
// 372 blob the same row applies with zero cards and fails here.

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

// TestSetGatesEveryLooseningAndLetsEveryNarrowingApply is R9's acceptance reading
// as a table: with NO confirmation channel wired, exactly the three switches that
// REMOVE questions are refused, and the three that add them apply.
//
// The direction columns are not decoration. The table asserts 3 widening / 3
// narrowing, so a predicate that refused everything (or nothing) cannot pass on a
// count alone, and the three narrowing rows keep the gate from turning "no
// confirmation channel" into "no mode change ever". The middle widening - the row
// tickets 367 and 372 disagreed about - is the one R9 decides: it costs a card, so
// without a channel it is refused.
func TestSetGatesEveryLooseningAndLetsEveryNarrowingApply(t *testing.T) {
	ctx := context.Background()
	pairs := sixPairs()
	if len(pairs) != 6 {
		t.Fatalf("pair table holds %d switches, want 6 (3 modes, no self-switch)", len(pairs))
	}
	widening, narrowing, refused, applied := 0, 0, 0, 0
	for _, p := range pairs {
		s, mc, rec := newTestStore(t, p.from.String(), nil)
		before, _ := mc.saved()
		err := s.Set(ctx, p.to, "r9", "owner")
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
		if p.widens {
			refused++
			if err == nil {
				t.Errorf("%s -> %s was APPLIED with no confirmation channel; every loosening must be "+
					"refused without one (R9)", p.from, p.to)
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
				t.Errorf("%s -> %s was REFUSED (%v); getting stricter is never guarded (PLAN.md:132)",
					p.from, p.to, err)
			}
			if last.Result != ResultApplied {
				t.Errorf("%s -> %s audit result = %q, want %q", p.from, p.to, last.Result, ResultApplied)
			}
			if mode, saves := mc.saved(); mode != p.to.String() || saves != 1 {
				t.Errorf("%s -> %s persisted mode=%q saves=%d, want %q/1",
					p.from, p.to, mode, saves, p.to)
			}
			if rec.count(string(ResultRefused)) != 0 {
				t.Errorf("%s -> %s wrote a refusal line for a tightening: %v", p.from, p.to, rec.all())
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
	if refused != 3 || applied != 3 {
		t.Fatalf("table produced %d refused / %d applied, want 3/3 (every loosening gated, every tightening free)",
			refused, applied)
	}
}

// TestSetCardsTheMiddleLooseningAndNeverACardForATightening is the other half of
// the rule: a channel is what makes a loosening reachable, and it costs exactly
// one card each time - including the middle tier, which is the cell R9 moved. The
// tightening at the end is the control: it must cost none.
func TestSetCardsTheMiddleLooseningAndNeverACardForATightening(t *testing.T) {
	ctx := context.Background()
	spy := &confirmSpy{}
	s, mc, rec := newTestStore(t, risk.ModeAskEveryStepName, spy.Confirm)
	if err := s.Set(ctx, risk.ModeAskHighRisk, "r9", "owner"); err != nil {
		t.Fatalf("ask_every_step -> ask_high_risk refused although a confirmation channel is attached: %v", err)
	}
	if got := s.PermissionMode(); got != risk.ModeAskHighRisk {
		t.Errorf("mode = %s after the allowed loosening, want ask_high_risk", got)
	}
	if mode, saves := mc.saved(); mode != risk.ModeAskHighRiskName || saves != 1 {
		t.Errorf("persisted mode=%q saves=%d, want ask_high_risk/1", mode, saves)
	}
	if spy.n() != 1 {
		t.Errorf("the middle-tier loosening cost %d confirmations; R9 cards every loosening, want 1", spy.n())
	}
	if len(spy.calls) == 1 && (spy.calls[0].From != risk.ModeAskEveryStep || spy.calls[0].To != risk.ModeAskHighRisk) {
		t.Errorf("the card was raised about %s -> %s, want ask_every_step -> ask_high_risk",
			spy.calls[0].From, spy.calls[0].To)
	}
	if rec.count(string(ResultApplied)) != 1 {
		t.Errorf("the allowed loosening is not in the audit as applied: %v", rec.all())
	}

	// The next loosening costs exactly one more card, and it is the auto_approve
	// one - "one card per loosening", not "one card per process".
	if err := s.Set(ctx, risk.ModeAutoApprove, "r9", "owner"); err != nil {
		t.Fatalf("ask_high_risk -> auto_approve refused with a channel attached: %v", err)
	}
	if spy.n() != 2 {
		t.Errorf("two loosenings cost %d confirmations in total, want exactly 2", spy.n())
	}
	if len(spy.calls) == 2 && (spy.calls[1].From != risk.ModeAskHighRisk || spy.calls[1].To != risk.ModeAutoApprove) {
		t.Errorf("the second card was raised about %s -> %s, want ask_high_risk -> auto_approve",
			spy.calls[1].From, spy.calls[1].To)
	}

	// The control: a tightening is not a request for a card.
	if err := s.Set(ctx, risk.ModeAskHighRisk, "r9", "owner"); err != nil {
		t.Fatalf("auto_approve -> ask_high_risk (a tightening) was refused: %v", err)
	}
	if spy.n() != 2 {
		t.Errorf("the tightening cost a confirmation (%d total); getting stricter must cost none", spy.n())
	}
	if mode, saves := mc.saved(); mode != risk.ModeAskHighRiskName || saves != 3 {
		t.Errorf("persisted mode=%q saves=%d, want ask_high_risk/3", mode, saves)
	}
}

// TestSetRefusesAutoApproveOnAnUnreadableStoredMode pins the invariant the gate
// leans on: PermissionMode() fail-closes to ask_every_step, so an unreadable
// config is still a *different* mode from auto_approve and the direction gate
// refuses the switch (ask_every_step -> auto_approve is a loosening). Without the
// gate's nil check this path would dereference a nil channel instead of refusing
// (the panic shot in 372-r2's evidence).
func TestSetRefusesAutoApproveOnAnUnreadableStoredMode(t *testing.T) {
	s, mc, _ := newTestStore(t, "not_a_mode", nil)
	if got := s.PermissionMode(); got != risk.ModeAskEveryStep {
		t.Fatalf("an unreadable stored mode reads as %s, want the strict default ask_every_step", got)
	}
	err := s.Set(context.Background(), risk.ModeAutoApprove, "r9", "owner")
	if err == nil {
		t.Fatal("auto_approve was reached from an unreadable stored mode with no channel")
	}
	if _, saves := mc.saved(); saves != 0 {
		t.Errorf("state moved under a refusal: saves=%d", saves)
	}
}

// TestIsWideningOnTheDocumentedLadder keeps the ladder-order witness alive: since
// R9 this function IS the gate's first term (RequiresCard), and it is also the
// vocabulary the audit words "widened"/"narrowed" rest on. internal/panel and
// internal/command drive the SAME order over the SAME names, so a reorder breaks
// all three ends rather than silently flipping one.
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
// the answered-no case on a loosening that asks.
func TestSetConfirmRefusalOnTheAutoTargetIsClassifiedAsRefused(t *testing.T) {
	spy := &confirmSpy{answer: errors.New("用户拒绝")}
	s, mc, rec := newTestStore(t, risk.ModeAskEveryStepName, spy.Confirm)
	err := s.Set(context.Background(), risk.ModeAutoApprove, "r9", "owner")
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
