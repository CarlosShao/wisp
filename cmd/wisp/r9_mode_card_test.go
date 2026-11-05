package main

// R9 (owner 2026-11-05, ledger A1541 item 4): "every loosening shows a
// confirmation card" - read from the two production doors that can move the档.
//
// THE GAP THIS FILE MEASURES. The runtime gate used to judge the TARGET, so
// `wisp run`'s seam and the panel pipe's seam applied ask_every_step ->
// ask_high_risk with zero cards while the config-file channel (D36) carded the
// same loosening. R9 makes the direction the judgement at every door, and these
// cases are that sentence as two readings:
//
//   - the run seam (run.go's rt.modes, the SAME perm.Store the bridge reads):
//     every loosening costs exactly one answered card, every tightening costs
//     none, and the file ends where the last switch put it;
//   - the panel pipe (panel_inbound.go's chain, which has no card surface): a
//     loosening is REFUSED with panel.ErrNoL2Confirm and nothing lands, while a
//     tightening still lands. That refusal is R9's accepted cost, named in
//     docs/evidence/s1/R9-impl.md: a pipe that cannot show a card cannot loosen.
//
// The pairing is the point: a gate that refused everything would fail the
// tightening rows, and a gate that judged the target would fail the middle-tier
// loosening row.

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/perm"
	"github.com/CarlosShao/wisp/internal/risk"
)

// r9Limb is the C18 leg in the shape both seams consume: it records the target of
// every switch it was asked about and answers allow, i.e. it is the "the operator
// clicked" arm rather than the "nobody was asked" one (that arm is internal/perm's
// store_r9_test.go).
type r9Limb struct {
	mu sync.Mutex
	to []risk.Mode
}

func (l *r9Limb) fn() perm.ConfirmFunc {
	return func(_ context.Context, sw perm.Switch) (perm.Card, error) {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.to = append(l.to, sw.To)
		return perm.AnsweredCard("r9-native"), nil
	}
}

func (l *r9Limb) asked() []risk.Mode {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]risk.Mode(nil), l.to...)
}

// r9SwitchLine returns the MODE-SWITCH row for one from -> to pair, or "".
func r9SwitchLine(log, from, to string) string {
	for _, line := range strings.Split(log, "\n") {
		if strings.Contains(line, "MODE-SWITCH") &&
			strings.Contains(line, "from="+from) && strings.Contains(line, "to="+to) {
			return line
		}
	}
	return ""
}

// TestR9TheRunSeamCardsEveryLooseningAndNoTightening drives the run seam's own
// Store: two loosenings and two tightenings in one boot, then reads the leg, the
// audit rows and the file the next process would start from.
func TestR9TheRunSeamCardsEveryLooseningAndNoTightening(t *testing.T) {
	h := t101boot(t, "") // fresh install: ask_every_step
	leg := &r9Limb{}
	var setErrs []string
	hook := func(rt *agentRuntime) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// loosening, loosening, tightening, tightening.
		for _, to := range []risk.Mode{
			risk.ModeAskHighRisk, risk.ModeAutoApprove,
			risk.ModeAskHighRisk, risk.ModeAskEveryStep,
		} {
			if err := rt.modes.Set(ctx, to, "cli", "r9"); err != nil {
				setErrs = append(setErrs, to.String()+": "+err.Error())
			}
		}
	}
	code, log := h.start(t, leg.fn(), hook)
	if code != 0 {
		t.Fatalf("exit %d\n%s", code, log)
	}
	if len(setErrs) != 0 {
		t.Fatalf("a switch was refused although the leg answers allow: %v", setErrs)
	}

	want := []risk.Mode{risk.ModeAskHighRisk, risk.ModeAutoApprove}
	if got := leg.asked(); len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("L2 cards raised = %v, want %v (one per loosening, none for a tightening): "+
			"the runtime gate is not judging the direction (R9)", got, want)
	}

	for _, row := range []struct {
		from, to string
		asked    string
	}{
		{risk.ModeAskEveryStepName, risk.ModeAskHighRiskName, "asked=1"},
		{risk.ModeAskHighRiskName, risk.ModeAutoApproveName, "asked=1"},
		{risk.ModeAutoApproveName, risk.ModeAskHighRiskName, "asked=0"},
		{risk.ModeAskHighRiskName, risk.ModeAskEveryStepName, "asked=0"},
	} {
		line := r9SwitchLine(log, row.from, row.to)
		if line == "" {
			t.Errorf("no MODE-SWITCH row for %s -> %s:\n%s", row.from, row.to, log)
			continue
		}
		if !strings.Contains(line, row.asked) || !strings.Contains(line, "result=applied") {
			t.Errorf("%s -> %s row does not read %s/applied: %s", row.from, row.to, row.asked, line)
		}
	}

	// The last switch was the tightening, so the next process starts strict: the
	// card requirement did not turn a tightening into a refusal that left the looser
	// value on disk.
	if got := h.readMode(t); got != risk.ModeAskEveryStep {
		t.Errorf("config.toml reads %v after the final tightening, want ask_every_step", got)
	}
}

// TestR9TheCLIPipeRefusesALooseningItCannotCard is the other door. The pipe has
// nowhere to show an L2 card, so under R9 it can no longer loosen at all - the
// accepted cost of the ruling, refused before the write leg sees it - while the
// direction that makes the machine safer still works with no card surface at all.
func TestR9TheCLIPipeRefusesALooseningItCannotCard(t *testing.T) {
	ctx := context.Background()

	t.Run("a loosening is refused and nothing lands", func(t *testing.T) {
		dir := newInboundDataDir33(t, risk.ModeAskEveryStepName)
		disp, err := newPanelInboundDispatch(dir, func(string, ...any) {})
		if err != nil {
			t.Fatalf("the CLI chain did not assemble: %v", err)
		}
		_, err = disp.Handle(ctx, modeRequest114r2("r9-cli-ae-ah", risk.ModeAskHighRiskName))
		if !errors.Is(err, panel.ErrNoL2Confirm) {
			t.Errorf("ask_every_step -> ask_high_risk must be refused with ErrNoL2Confirm on a pipe "+
				"with no card surface (R9, fail-closed); got %v", err)
		}
		if got := readInboundCfg33(t, dir); got != risk.ModeAskEveryStepName {
			t.Errorf("the refused loosening still moved the file: reads %q, want ask_every_step", got)
		}
	})

	t.Run("a tightening still lands with no card surface", func(t *testing.T) {
		dir := newInboundDataDir33(t, risk.ModeAskHighRiskName)
		disp, err := newPanelInboundDispatch(dir, func(string, ...any) {})
		if err != nil {
			t.Fatalf("the CLI chain did not assemble: %v", err)
		}
		if _, err := disp.Handle(ctx, modeRequest114r2("r9-cli-ah-ae", risk.ModeAskEveryStepName)); err != nil {
			t.Errorf("the tightening was refused on a legless pipe: %v (getting stricter is never guarded)", err)
		}
		if got := readInboundCfg33(t, dir); got != risk.ModeAskEveryStepName {
			t.Errorf("the tightening reached no write leg: reads %q, want ask_every_step", got)
		}
	})
}
