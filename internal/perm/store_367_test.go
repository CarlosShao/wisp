package perm

// Ticket 367: Set's gate judges the DIRECTION, not the target value.
//
// The hole this file nails shut, measured on 4ee57d318 by the 367-a1 census and
// re-measured here: `Set` asked "is the target auto_approve?" (`store.go`'s
// `if to == risk.ModeAutoApprove`), so ask_every_step -> ask_high_risk - a
// widening - was applied with zero confirmation and zero refusal whenever the
// caller was not the panel handler. The panel handler had the direction right
// and the store had it wrong: two transition points, two different judges.
//
// Every expectation below is the POST-fix one, which is what makes these teeth:
// on the pre-fix blob the first table's ask_every_step -> ask_high_risk row
// reads applied=1 and fails.

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

// TestSetGatesEveryWideningAndNoNarrowing is ticket 367's acceptance reading as
// a table: with NO confirmation channel wired, all three widenings are refused
// (including the non-auto_approve one) and all three narrowings still apply.
//
// The narrowing half is not decoration. It is the tooth for the boundary that
// says "only widening is gated": a predicate that treated any change as a
// widening would keep the first half green and turn the second half red, and a
// gate that refused everything would do the same. One table, both directions.
func TestSetGatesEveryWideningAndNoNarrowing(t *testing.T) {
	ctx := context.Background()
	pairs := sixPairs()
	if len(pairs) != 6 {
		t.Fatalf("pair table holds %d switches, want 6 (3 modes, no self-switch)", len(pairs))
	}
	widening, narrowing := 0, 0
	for _, p := range pairs {
		s, mc, rec := newTestStore(t, p.from.String(), nil)
		before, _ := mc.saved()
		err := s.Set(ctx, p.to, "t367", "owner")
		last, ok := s.Last()
		if !ok {
			t.Fatalf("%s -> %s: the switch left no audit record at all", p.from, p.to)
		}
		switch {
		case p.widens:
			widening++
			if err == nil {
				t.Errorf("%s -> %s was APPLIED with no confirmation channel; a widening must be refused (ticket 367)",
					p.from, p.to)
			}
			if last.Result != ResultRefused {
				t.Errorf("%s -> %s audit result = %q, want %q", p.from, p.to, last.Result, ResultRefused)
			}
			if got := s.PermissionMode(); got != p.from {
				t.Errorf("%s -> %s moved the mode to %s despite the refusal", p.from, p.to, got)
			}
			if mode, saves := mc.saved(); mode != before || saves != 0 {
				t.Errorf("%s -> %s persisted anyway: mode=%q saves=%d, want %q/0",
					p.from, p.to, mode, saves, before)
			}
			// The refusal has to name both ends, or an operator reading the log
			// cannot tell which widening was attempted.
			if err != nil && (!strings.Contains(err.Error(), p.from.String()) ||
				!strings.Contains(err.Error(), p.to.String())) {
				t.Errorf("%s -> %s refusal does not name both modes: %v", p.from, p.to, err)
			}
			if rec.count(string(ResultRefused)) != 1 {
				t.Errorf("%s -> %s: %d refusal lines in the audit, want 1: %v",
					p.from, p.to, rec.count(string(ResultRefused)), rec.all())
			}
		case !p.widens:
			narrowing++
			if err != nil {
				t.Errorf("%s -> %s was REFUSED (%v); getting stricter must never be blocked (R20)", p.from, p.to, err)
			}
			if last.Result != ResultApplied {
				t.Errorf("%s -> %s audit result = %q, want %q", p.from, p.to, last.Result, ResultApplied)
			}
			if got := s.PermissionMode(); got != p.to {
				t.Errorf("%s -> %s did not apply: mode = %s", p.from, p.to, got)
			}
			if mode, saves := mc.saved(); mode != p.to.String() || saves != 1 {
				t.Errorf("%s -> %s persisted mode=%q saves=%d, want %q/1", p.from, p.to, mode, saves, p.to)
			}
			if rec.count(string(ResultRefused)) != 0 {
				t.Errorf("%s -> %s wrote a refusal line for a narrowing: %v", p.from, p.to, rec.all())
			}
		}
	}
	if widening != 3 || narrowing != 3 {
		t.Fatalf("table classified %d widening / %d narrowing, want 3/3", widening, narrowing)
	}
}

// TestSetWideningPassesOnceAChannelIsAttached is the other half of the gate: it
// refuses for want of a channel, not because widening is forbidden. A composition
// that did wire one may widen - and, per R20/M4's asymmetry, the non-auto step
// still costs no ask at all.
func TestSetWideningPassesOnceAChannelIsAttached(t *testing.T) {
	ctx := context.Background()
	spy := &confirmSpy{}
	s, mc, rec := newTestStore(t, risk.ModeAskEveryStepName, spy.Confirm)
	if err := s.Set(ctx, risk.ModeAskHighRisk, "t367", "owner"); err != nil {
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
	if err := s.Set(ctx, risk.ModeAutoApprove, "t367", "owner"); err != nil {
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

// TestSetAsksAutoApproveEvenAfterAnUnreadableStoredMode pins the invariant the
// gate's structure leans on: PermissionMode() fail-closes to ask_every_step, so
// "to == auto_approve" is always a widening and Set never reaches the ask with a
// nil channel. Without that, dropping the old `if s.confirm == nil` from inside
// the auto_approve branch would be a nil dereference instead of a refusal.
func TestSetAsksAutoApproveEvenAfterAnUnreadableStoredMode(t *testing.T) {
	s, mc, _ := newTestStore(t, "not_a_mode", nil)
	if got := s.PermissionMode(); got != risk.ModeAskEveryStep {
		t.Fatalf("an unreadable stored mode reads as %s, want the strict default ask_every_step", got)
	}
	err := s.Set(context.Background(), risk.ModeAutoApprove, "t367", "owner")
	if err == nil {
		t.Fatal("auto_approve was reached from an unreadable stored mode with no channel")
	}
	if _, saves := mc.saved(); saves != 0 {
		t.Errorf("state moved under a refusal: saves=%d", saves)
	}
}

// TestIsWideningOnTheDocumentedLadder is this package's end of the shared
// predicate: risk.Mode's constants are ordered strict -> loose ("Do not reorder
// these constants"), and every adjacent step up is a widening while every step
// down is not. internal/panel drives the SAME function over the SAME names, so a
// reorder breaks both ends rather than silently flipping one.
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

// TestSetConfirmRefusalOnAWideningIsClassifiedAsRefused keeps the two widening
// failure modes distinguishable in the audit: "the operator said no" is not
// "there was no way to ask". The nil-channel case above is ResultRefused; this
// is the answered-no case on the same non-auto widening.
func TestSetConfirmRefusalOnAWideningIsClassifiedAsRefused(t *testing.T) {
	spy := &confirmSpy{answer: errors.New("用户拒绝")}
	s, mc, rec := newTestStore(t, risk.ModeAskEveryStepName, spy.Confirm)
	err := s.Set(context.Background(), risk.ModeAutoApprove, "t367", "owner")
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
