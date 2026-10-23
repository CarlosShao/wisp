package panel

// Ticket 367: the panel gate and the store gate must reach the SAME verdict.
//
// The disease the ticket names is "two transition points, two different judges":
// this package has always judged the DIRECTION (modeIsWidening, now
// perm.IsWidening) while perm.Store judged the TARGET ("is it auto_approve?").
// A panel request that widened to ask_high_risk was refused here and would have
// been applied unchanged by Set had anything else called it.
//
// The rows below are the invariant, not either implementation: for every ordered
// pair of the three modes, with no confirmation channel wired anywhere, the two
// layers must refuse the same switches. Fixing one layer alone turns this red,
// which is why it is a tooth for the class rather than for the pair.

import (
	"context"
	"fmt"
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
	pairs, widening, narrowing := 0, 0, 0
	for fromIdx, from := range ladder {
		for toIdx, to := range ladder {
			if fromIdx == toIdx {
				continue
			}
			pairs++
			if toIdx > fromIdx {
				widening++
			} else {
				narrowing++
			}

			// Layer 1: the real perm.Store, no confirmation channel. This is the
			// composition both production entries build (panel_inbound.go:254).
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

			if (storeErr == nil) != (handlerErr == nil) {
				t.Errorf("%s -> %s: the two gates disagree - perm.Store refused=%v (%v), panel handler refused=%v (%v)",
					from, to, storeErr != nil, storeErr, handlerErr != nil, handlerErr)
			}
			// A refusal on either layer must leave the mode where it was.
			if storeErr != nil && st.PermissionMode() != from {
				t.Errorf("%s -> %s: the store refused but moved to %s", from, to, st.PermissionMode())
			}
			if handlerErr != nil && w.setCalls != 0 {
				t.Errorf("%s -> %s: the handler refused but called the write leg %d time(s)", from, to, w.setCalls)
			}
			// And every widening must be the refused direction on BOTH layers:
			// two judges that happen to agree on "refuse everything" would pass
			// the comparison above and fail this line.
			if toIdx > fromIdx {
				if storeErr == nil || handlerErr == nil {
					t.Errorf("%s -> %s is a widening and both layers must refuse it without a channel: store=%v handler=%v",
						from, to, storeErr, handlerErr)
				}
			}
			if toIdx < fromIdx {
				if storeErr != nil || handlerErr != nil {
					t.Errorf("%s -> %s is a narrowing and neither layer may refuse it: store=%v handler=%v",
						from, to, storeErr, handlerErr)
				}
			}
		}
	}
	if pairs != 6 || widening != 3 || narrowing != 3 {
		t.Fatalf("matrix covered %d pairs (%d widening / %d narrowing), want 6 (3/3)", pairs, widening, narrowing)
	}
}
