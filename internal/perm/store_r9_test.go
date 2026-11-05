package perm

// R9 (owner 2026-11-05, ledger A1541 item 4): "every loosening shows a
// confirmation card".
//
// WHAT WAS TRUE BEFORE THIS FILE. The runtime gate judged the TARGET: only the
// switch INTO auto_approve cost an answered L2 card, so ask_every_step ->
// ask_high_risk - a loosening that silences every L1 pre-execution window - landed
// with zero cards. The config-file channel (D36's rule 1) has carded every
// loosening since ticket 90; the runtime channel did not, so the SAME request read
// differently depending on which door it came through. R9 closes that gap by
// judging the DIRECTION at the runtime gate as well.
//
// WHAT THESE CASES PIN, in pairs, so a gate that simply refuses everything cannot
// pass either:
//
//   - every loosening, with a real asking limb: applied once, with exactly one
//     answered card, and the audit row says asked=1 and names the card;
//   - every loosening, with a limb that never asks: refused as
//     card-not-answered (ticket 371's rule now covers the middle tier too);
//   - every TIGHTENING: still free - zero cards - and still lands. Getting
//     stricter is never guarded (PLAN.md:132), and a rule that carded it would be
//     a gate that protects nothing.
//
// The no-channel half of the same table is store_367_test.go's (R9 rewrote it):
// a loosening this machine cannot card is refused outright, fail-closed.

import (
	"context"
	"errors"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// TestR9EveryLooseningCostsOneAnsweredCard is the ticket's headline reading: with
// a confirmation channel attached, each of the three loosening pairs raises
// exactly one card and applies; each of the three tightening pairs raises none.
func TestR9EveryLooseningCostsOneAnsweredCard(t *testing.T) {
	ctx := context.Background()
	pairs := sixPairs()
	if len(pairs) != 6 {
		t.Fatalf("pair table holds %d switches, want 6 (3 modes, no self-switch)", len(pairs))
	}
	loosened, tightened := 0, 0
	for _, p := range pairs {
		spy := &confirmSpy{}
		s, mc, _ := newTestStore(t, p.from.String(), spy.Confirm)
		err := s.Set(ctx, p.to, "r9", "owner")
		last, ok := s.Last()
		if !ok {
			t.Fatalf("%s -> %s: the switch left no audit record at all", p.from, p.to)
		}
		if p.widens {
			loosened++
			if err != nil {
				t.Errorf("%s -> %s is a loosening and a channel is attached; it was refused: %v (R9)",
					p.from, p.to, err)
			}
			if got := spy.n(); got != 1 {
				t.Errorf("%s -> %s cost %d L2 card(s), want exactly 1 (R9: every loosening shows a card)",
					p.from, p.to, got)
			}
			if last.Result != ResultApplied || !last.Asked || last.Card == "" {
				t.Errorf("%s -> %s: result=%q asked=%v card=%q, want applied/true/<id> - the card that was "+
					"answered has to be readable off the row", p.from, p.to, last.Result, last.Asked, last.Card)
			}
			if len(spy.calls) == 1 && (spy.calls[0].From != p.from || spy.calls[0].To != p.to) {
				t.Errorf("the card was raised about %s -> %s, want %s -> %s",
					spy.calls[0].From, spy.calls[0].To, p.from, p.to)
			}
			if mode, saves := mc.saved(); mode != p.to.String() || saves != 1 {
				t.Errorf("%s -> %s persisted mode=%q saves=%d, want %q/1",
					p.from, p.to, mode, saves, p.to)
			}
		} else {
			tightened++
			if err != nil {
				t.Errorf("%s -> %s is a tightening; it was refused: %v (getting stricter is never guarded)",
					p.from, p.to, err)
			}
			if got := spy.n(); got != 0 {
				t.Errorf("%s -> %s cost %d L2 card(s); a tightening must cost none (PLAN.md:132)",
					p.from, p.to, got)
			}
			if last.Asked || last.Card != "" {
				t.Errorf("%s -> %s: asked=%v card=%q, want false/\"\" - nobody was asked",
					p.from, p.to, last.Asked, last.Card)
			}
			if last.Result != ResultApplied {
				t.Errorf("%s -> %s: result=%q, want %q", p.from, p.to, last.Result, ResultApplied)
			}
			if mode, saves := mc.saved(); mode != p.to.String() || saves != 1 {
				t.Errorf("%s -> %s persisted mode=%q saves=%d, want %q/1",
					p.from, p.to, mode, saves, p.to)
			}
		}
	}
	if loosened != 3 || tightened != 3 {
		t.Fatalf("table classified %d loosening / %d tightening, want 3/3", loosened, tightened)
	}
}

// TestR9ANeverAskingLimbCannotLoosenTheMiddleTier is ticket 371's answered-card
// rule applied to the middle tier. Before R9 the never-asking shape only mattered
// on auto_approve; a loosening to ask_high_risk never reached a limb at all. Now
// the two tiers are the same rule, so the same fake limb has to be refused there
// too - and the TIGHTENING next to it must still land without touching the limb,
// which is what keeps this from being "refuse everything".
func TestR9ANeverAskingLimbCannotLoosenTheMiddleTier(t *testing.T) {
	ctx := context.Background()
	limb := &leg371{}
	s, mc, rec := newTestStore(t, risk.ModeAskEveryStepName, neverAsks371(limb))

	err := s.Set(ctx, risk.ModeAskHighRisk, "r9", "owner")
	if err == nil {
		t.Fatal("ask_every_step -> ask_high_risk applied although no card was ever answered: R9's " +
			"middle tier inherited ticket 371's hole")
	}
	if !errors.Is(err, ErrCardNotAnswered) {
		t.Errorf("refusal = %v, want it to wrap ErrCardNotAnswered (the limb was entered and never asked, "+
			"which is not the no-channel story)", err)
	}
	if limb.invoked != 1 {
		t.Errorf("the limb was entered %d time(s), want 1 - the gate must ASK for a loosening", limb.invoked)
	}
	if last, _ := s.Last(); last.Result != ResultCardNotAnswered {
		t.Errorf("audit result = %q, want %q", last.Result, ResultCardNotAnswered)
	}
	if mode, saves := mc.saved(); mode != risk.ModeAskEveryStepName || saves != 0 {
		t.Errorf("persisted mode=%q saves=%d, want ask_every_step/0 (nothing may land)", mode, saves)
	}
	if got := rec.count(string(ResultCardNotAnswered)); got != 1 {
		t.Errorf("unanswered-card rows = %d, want 1: %v", got, rec.all())
	}

	// The other direction, a separate store seeded on the looser档: a tightening
	// is not a request for a card, so it lands and the limb stays untouched.
	tight := &leg371{}
	s2, mc2, _ := newTestStore(t, risk.ModeAskHighRiskName, neverAsks371(tight))
	if err := s2.Set(ctx, risk.ModeAskEveryStep, "r9", "owner"); err != nil {
		t.Fatalf("the tightening was refused: %v", err)
	}
	if tight.invoked != 0 {
		t.Errorf("the tightening entered the confirmation limb %d time(s); it must cost none", tight.invoked)
	}
	if mode, saves := mc2.saved(); mode != risk.ModeAskEveryStepName || saves != 1 {
		t.Errorf("persisted mode=%q saves=%d, want ask_every_step/1", mode, saves)
	}
}
