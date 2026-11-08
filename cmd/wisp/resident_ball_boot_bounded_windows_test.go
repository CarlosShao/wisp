//go:build windows

package main

// Ticket 408 r3, the test-layer half: every test that boots the resident ball
// waits for an external event (WebView2 environment + controller creation on
// the ui-sta thread) with NO deadline of its own - ball.New's waitStarted has
// none, and the orb's Embed pumps until a completion callback fires. When the
// whole-package run degrades (ticket 408's measured dead points: 17m35s at
// Test258BridgeRebindsLiveKeysFromConfigEdit, 10 minutes at
// TestAC294MuteHotkeyRegistersLiveOnARealBallWindow, both inside
// startResidentBall), that wait never resolved and the package burned its
// entire -timeout budget on ONE silent hang.
//
// The mainstream harness rule this file installs: a test's wait for an
// external event is BOUNDED, and expiry FAILS the test loudly with
// diagnostics instead of hanging. The budget is generous against every
// healthy reading measured on this tree (single boots 1.3-1.5s; the worst
// loaded boot 7.9s), so an expiry means the boot really wedged - the same
// judgement go test's own -timeout makes, one level down, with the full
// goroutine dump the wedged state owes the evidence.
//
// What expiry does NOT do: it does not unwind the wedged boot. The goroutine
// inside startResidentBall stays wedged (its STA thread, window and hotkeys
// with it), so a run that already hit one expiry is already red and may read
// further reds from the wedged desktop state - that is the honest cost of a
// wedged boot, made visible, not hidden.

import (
	"runtime"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/observe"
)

// residentBallBootBudget bounds one startResidentBall call. 2 minutes is
// ~15x the worst healthy boot measured on this tree (7.9s) and caps a whole
// run's worst-case burn at 8 sites x 2m = 16m inside the 40m package budget.
const residentBallBootBudget = 2 * time.Minute

// startResidentBallBounded is the bounded form of startResidentBall every
// test in this package uses. Same construction, same verdict, same nil-able
// result - the only difference is that a wedged boot turns into a loud
// t.Fatalf carrying the goroutine dump instead of a silent hang.
func startResidentBallBounded(t *testing.T, reg *observe.Registry, onCancelEsc escVetoFunc, hotCfg func() ball.HotkeyConfig, hotReload func() ball.HotkeyConfig, hooks ...ballHostHook) *residentBall {
	t.Helper()
	type outcome struct {
		rb *residentBall
	}
	done := make(chan outcome, 1)
	go func() {
		done <- outcome{startResidentBall(reg, onCancelEsc, hotCfg, hotReload, hooks...)}
	}()
	timer := time.NewTimer(residentBallBootBudget)
	defer timer.Stop()
	select {
	case r := <-done:
		return r.rb
	case <-timer.C:
		buf := make([]byte, 1<<20)
		n := runtime.Stack(buf, true)
		t.Fatalf("the resident ball boot did not come back within %s (healthy boots on this tree measure 1.3-7.9s; past this budget the boot is wedged - ticket 408's dead-point shape). Full goroutine dump:\n%s",
			residentBallBootBudget, buf[:n])
		return nil
	}
}
