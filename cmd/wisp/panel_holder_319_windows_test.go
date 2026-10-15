//go:build windows

package main

// Ticket 319 AC#1/AC#3 (N1), the HEADLESS half of "a panel thread that died is
// not the end of this process's panel".
//
// The defect (ticket 319, root cause fixed by leg 319-a1: .scratch/wisp/probes/
// 319/a1/census.md section 1): the panel window being closed - the owner's X, a
// system SC_CLOSE, or host RequestDispose, which is the same road - runs
// WM_CLOSE -> DestroyWindow -> WM_DESTROY -> PostQuitMessage(0), the library
// pump Run() returns, the thread's loop() returns, rp.finished closes, and every
// later gesture is refused by post() forever. The 2026-10-14 live log has both
// halves: "panel thread ending / why: the library pump returned" at 15:41:45,
// then "panel gesture could not be handed over" for every tray/ball gesture
// afterwards.
//
// This file pins BOTH sides of the fix. The raw half - a gesture handed to the
// dead thread - must stay refused (that is today's only path, verbatim, and the
// refusal itself is not the bug). The holder half - the same gesture through the
// assembly root's panelHolder - must rebuild: a fresh thread on the same
// manager. The rebuild is where the nail has teeth: delete the isFinished()
// branch in panelHolder.toggle and the holder half goes red on the spot, both on
// "toggle returned true" and on "the spawn was asked a second time".

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/panel"
)

// panel319DataPath is the WebView2 user-data folder this file's fixtures point
// at. Nothing here creates a window (the headless half never asks bringUp), so
// the folder is a name and not a profile - kept unique-ish per go test run only
// where it matters, and shared here on purpose: two fixtures in one run would
// otherwise race over one WebView2 profile if a later edit ever gave them a
// window.
func panel319DataPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(os.TempDir(), "wisp-319r1-holder-profile")
}

// newPanelFixture319 builds a production manager and the production panel
// thread, with NO window ever requested (this is the shape every resident boot
// starts in).
func newPanelFixture319(t *testing.T) (*residentPanel, *PanelManager) {
	t.Helper()
	disp := &panel.ComposerDispatch{Mode: &recordingModeHandler{}}
	mgr := NewPanelManager(disp, nil, panel319DataPath(t))
	rp := startResidentPanel(observe.NewRegistry(), mgr)
	t.Cleanup(rp.stop)
	return rp, mgr
}

// newDeadPanelThread319 ends the fixture's thread through its own stop() and
// returns it only once isFinished() is true - the observable state a closed
// window leaves behind. stop() is used and not a real SC_CLOSE because this
// file is the headless tier (no window may exist here); the live replay of the
// close chain itself lives in panel_death_replay_319_winlive_test.go, and what
// post() sees is the same fact either way: the thread's goroutine has returned.
func newDeadPanelThread319(t *testing.T) *residentPanel {
	t.Helper()
	rp, _ := newPanelFixture319(t)
	rp.stop()
	if !rp.isFinished() {
		t.Fatalf("the fixture's thread did not finish after stop(); every assertion below would be about a live thread")
	}
	return rp
}

// TestDeadPanelThreadRefusesEveryGesture319 is the "today, unfixed" reading at
// the headless level, kept as a case so the refusal that IS correct (a dead
// thread takes no request) can never be quietly replaced by a delivery that
// pretends to have happened. It goes red if post() ever answers true for a
// thread whose goroutine has returned.
func TestDeadPanelThreadRefusesEveryGesture319(t *testing.T) {
	rp := newDeadPanelThread319(t)
	toggle := rp.RequestToggle("319-dead-thread-toggle")
	show := rp.RequestShow("319-dead-thread-show")
	if toggle || show {
		t.Errorf("a gesture handed to a dead panel thread was accepted (toggle=%v show=%v) - post() must refuse a thread whose goroutine has returned, or the caller reports a window move that never happened", toggle, show)
	}
	if rp.isFinished() != true {
		t.Errorf("isFinished() = false on a thread whose goroutine returned")
	}
	t.Logf("319 dead-thread raw reading: RequestToggle=%v RequestShow=%v isFinished=%v", toggle, show, rp.isFinished())
}
