package panel

// Ticket 367: the panel gate and the store gate must reach the SAME verdict.
//
// The disease the ticket names is "two transition points, two different judges":
// this package has always judged the DIRECTION (modeIsWidening, now
// perm.IsWidening) while perm.Store judged the TARGET ("is it auto_approve?").
// A panel request that widened to ask_high_risk was refused here and would have
// been applied unchanged by Set had anything else called it. Ticket 367 made both
// layers judge the direction; ticket 372 found that judgement refuted by the
// mainstream survey (docs/evidence/s1/m4-survey.md:96: no surveyed harness gates
// on "wider"), reverted it to the target and kept ONE definition of the target
// rule (perm.RequiresCard), called here and in Set.
//
// The rows below are the invariant, not either implementation: for every ordered
// pair of the three modes, with no confirmation channel wired anywhere, the two
// layers must refuse the same switches. Since 372 the refused set is exactly "the
// target is auto_approve" - the one档 R20/M4 cards - so the auto_approve rows are
// the teeth, and ask_every_step -> ask_high_risk is the row the revert buys.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/perm"
	"github.com/CarlosShao/wisp/internal/risk"
)

// t367Manager is the smallest perm.ConfigManager a real Store can be handed: the
// test's stand-in for the config.Manager the composition root passes in.
type t367Manager struct{ mode risk.Mode }

func (m *t367Manager) Config() *config.Config {
	c := config.NewDefaults()
	c.Risk.PermissionMode = m.mode.String()
	return c
}

func (m *t367Manager) SetPermissionMode(mode risk.Mode) error {
	m.mode = mode
	return nil
}

// TestPanelGateAndStoreGateAgreeOnEveryPair367 is the cross-layer agreement.
func TestPanelGateAndStoreGateAgreeOnEveryPair367(t *testing.T) {
	ladder := modeLadder(t)
	names := risk.ModeNames()
	ctx := context.Background()
	pairs, refused := 0, 0
	sawMiddleWidening := false
	for fromIdx, from := range ladder {
		for toIdx, to := range ladder {
			if fromIdx == toIdx {
				continue
			}
			pairs++
			wantRefused := perm.RequiresCard(to)
			if wantRefused {
				refused++
			}
			if from == risk.ModeAskEveryStep && to == risk.ModeAskHighRisk {
				sawMiddleWidening = true
			}

			// Layer 1: the real perm.Store, no confirmation channel. This is the
			// composition both production entries build (panel_inbound.go:272).
			mgr := &t367Manager{mode: from}
			st, err := perm.New(perm.Options{Manager: mgr})
			if err != nil {
				t.Fatalf("perm.New: %v", err)
			}
			storeErr := st.Set(ctx, to, "t367", "owner")

			// Layer 2: the panel handler in front of a writer that records calls,
			// with no confirm leg either.
			w := &fakeModeWriter{current: from, currentOK: true}
			h := &ModeWriteHandler{Modes: w}
			id := fmt.Sprintf("agree-%s-%s", names[fromIdx], names[toIdx])
			handlerErr := h.HandleModeRequest(ctx, modeReq(names[toIdx], id))

			if (storeErr != nil) != wantRefused {
				t.Errorf("%s -> %s: perm.Store refused=%v (%v), want refused=%v (target is auto_approve: %v)",
					from, to, storeErr != nil, storeErr, wantRefused, perm.RequiresCard(to))
			}
			if (handlerErr != nil) != wantRefused {
				t.Errorf("%s -> %s: panel handler refused=%v (%v), want refused=%v (target is auto_approve: %v)",
					from, to, handlerErr != nil, handlerErr, wantRefused, perm.RequiresCard(to))
			}
			if (storeErr == nil) != (handlerErr == nil) {
				t.Errorf("%s -> %s: the two gates disagree - perm.Store refused=%v (%v), panel handler refused=%v (%v)",
					from, to, storeErr != nil, storeErr, handlerErr != nil, handlerErr)
			}
			if wantRefused {
				// A refusal on either layer must leave the PERSISTED档 where it was.
				// The read side is deliberately not asked: without a leg it answers
				// ask_every_step even when the file says otherwise (372-r2).
				if mgr.mode != from {
					t.Errorf("%s -> %s: the refusal still moved the persisted档 to %s", from, to, mgr.mode)
				}
				if w.setCalls != 0 {
					t.Errorf("%s -> %s: the handler refused but called the write leg %d time(s)", from, to, w.setCalls)
				}
				if !errors.Is(handlerErr, ErrNoL2Confirm) {
					t.Errorf("%s -> %s: the refusal is not ErrNoL2Confirm: %v", from, to, handlerErr)
				}
				// Both ends must be named, or an operator reading the log cannot
				// tell which switch was attempted.
				if storeErr == nil || !strings.Contains(storeErr.Error(), from.String()) ||
					!strings.Contains(storeErr.Error(), to.String()) {
					t.Errorf("%s -> %s: the store refusal does not name both modes: %v", from, to, storeErr)
				}
				continue
			}
			// A target M4 does not card passes BOTH layers: nothing here is a
			// widening rule in disguise, and the middle档 is reachable.
			if mgr.mode != to {
				t.Errorf("%s -> %s: the allowed switch did not persist the target (file says %s)", from, to, mgr.mode)
			}
			if w.setCalls != 1 || w.setMode != to {
				t.Errorf("%s -> %s: the handler should call the write leg once with %s; called %d time(s) with %s",
					from, to, to, w.setCalls, w.setMode)
			}
		}
	}
	if pairs != 6 || refused != 2 {
		t.Fatalf("matrix covered %d pairs with %d refused, want 6 pairs and 2 refused (auto_approve from the two stricter档)", pairs, refused)
	}
	if !sawMiddleWidening {
		t.Fatal("the ask_every_step -> ask_high_risk row is missing: the revert's own case went untested")
	}
}
