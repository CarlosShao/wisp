//go:build windows && winlive

package main

// 319-v2 out-of-repo probe (injected with -overlay as
// cmd/wisp/probe_319v2_visible_winlive_test.go; never written into the tracked
// tree). Ticket 319 AC#2-a + AC#2-b, measured by THIS leg through the production
// second-launch hop: proc.SignalExistingInstance -> RunEventLoop
// (Runtime.handleActivation) -> SetActivationHook -> secondLaunchActivation319 ->
// panelHolder.show -> residentPanel.RequestShow/RequestRaise -> mgr.Show/Raise.
//
// Shape: the panel is ALREADY visible (opened through the ordinary gesture road)
// and the foreground is a prior window of this process (the operator's "where
// they came from"), then the second launch arrives. The probe asserts the two
// fact groups separately so a mutant can redden either half alone:
//   - AC#2-a: mgr.IsShown() stays true, the HWND is unchanged, the window is
//     alive. Route the hook back to h.toggle (the D1 regression, mutant A) and
//     this group reddens while the raise group does too.
//   - AC#2-b: GetForegroundWindow() == the panel's HWND, with the prior handle,
//     the foreground handle and SetForegroundWindow's own answer printed as
//     readings. Drop the RequestRaise hop (mutant B) and ONLY this group
//     reddens - IsShown must not be used as a stand-in for "in front".
//
// The prior window is this process's own editor window on the panel thread; the
// probe fails the MEASUREMENT (not the product) if that window cannot hold the
// foreground first. Layout names are per-run: a second `wisp --env dev` would
// signal the owner's instance, so the second launch is this process signalling a
// name it just created.

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/proc"
	"golang.org/x/sys/windows"
)

func TestProbe319v2VisibleStaysShownAndRaised(t *testing.T) {
	h, mgr := startPanelHolder319(t)

	// The operator's panel is already on screen, opened through the ordinary
	// gesture road (toggle's show half) - how an operator gets there before a
	// later double click.
	if !h.toggle("319v2-probe-visible-setup") {
		t.Fatalf("the holder refused the setup gesture")
	}
	waitPanelTrue(t, "the setup panel window to be shown", func() bool { return mgr.IsShown() })
	hwndBefore := mgr.windowHandle()
	if hwndBefore == 0 {
		t.Fatalf("the setup show reported a shown window with HWND 0")
	}

	// Runtime + event loop + production hook, per-run names.
	reg := observe.NewRegistry()
	layout := proc.Layout{
		Env:               buildinfo.EnvProd,
		DataDir:           t.TempDir(),
		MutexName:         fmt.Sprintf(`Local\wisp-test-319v2-visible-%d-%d`, os.Getpid(), time.Now().UnixNano()),
		ActivateEventName: fmt.Sprintf(`Local\wisp-test-319v2-visible-%d-%d.activate`, os.Getpid(), time.Now().UnixNano()),
		MutexEnabled:      true,
	}
	rt, err := proc.Boot(buildinfo.EnvProd, proc.WithLayout(layout), proc.WithRegistry(reg))
	if err != nil {
		t.Fatalf("proc.Boot: %v", err)
	}
	t.Cleanup(func() { rt.Shutdown(true) })

	hook := secondLaunchActivation319(h, nil)
	receipts := make(chan string, 4)
	rt.SetActivationHook(func() string {
		s := hook()
		select {
		case receipts <- s:
		default:
		}
		return s
	})
	go func() { rt.RunEventLoop() }()

	// Where the operator's foreground is when the second launch arrives: a window
	// of this process in front of the visible panel. Same-foreground failure is a
	// failed measurement, not a product verdict.
	prior := editorOnPanelThread(t, h.rp)
	if fg := uintptr(windows.GetForegroundWindow()); fg != uintptr(prior) {
		t.Fatalf("the prior window 0x%X did not take the foreground (the foreground is 0x%X); the reading below would be about an unknown desktop", uintptr(prior), fg)
	}
	raisesBefore := h.rp.raises.Load()
	t.Logf("BEFORE: panel hwnd=0x%X mgr.IsShown=%v alive=%v; prior=0x%X foreground=0x%X; resident raises=%d",
		hwndBefore, mgr.IsShown(), isLiveWindow(hwndBefore), uintptr(prior), uintptr(windows.GetForegroundWindow()), raisesBefore)

	// The second launch: the identical call the ErrAlreadyRunning branch makes.
	if err := proc.SignalExistingInstance(rt.Layout.ActivateEventName); err != nil {
		t.Fatalf("SignalExistingInstance: %v", err)
	}
	var receipt string
	select {
	case receipt = <-receipts:
	case <-time.After(panelThreadWait):
		t.Fatalf("the activation hook never answered within %v", panelThreadWait)
	}

	// The activation's requests were handed to the panel thread when the receipt
	// was written; this sentinel is posted through the SAME route afterwards, so
	// its closure runs after them (FIFO) and makes the state read deterministic
	// for the fixed road and for the mutant alike.
	served := make(chan struct{})
	if !h.rp.post(func() { close(served) }) {
		t.Fatalf("the sentinel post was refused by the panel thread")
	}
	select {
	case <-served:
	case <-time.After(panelThreadWait):
		t.Fatalf("the sentinel was not served within %v - failed measurement", panelThreadWait)
	}

	shownAfter := mgr.IsShown()
	afterHwnd := mgr.windowHandle()
	fgAfter := uintptr(windows.GetForegroundWindow())
	raiseRet := mgr.LastRaiseSetForeground()
	raisesAfter := h.rp.raises.Load()
	t.Logf("AFTER: receipt=%q; mgr.IsShown=%v hwnd=0x%X (same=%v) alive=%v; foreground=0x%X (panel=%v, prior=0x%X); resident raises %d->%d; SetForegroundWindow answered %d",
		receipt, shownAfter, afterHwnd, afterHwnd == hwndBefore, isLiveWindow(afterHwnd), fgAfter, fgAfter == hwndBefore, uintptr(prior), raisesBefore, raisesAfter, raiseRet)

	// AC#2-a: an already-visible panel must not be hidden.
	if !shownAfter {
		t.Errorf("AC#2-a RED (the visible panel was HIDDEN by the second launch): the window 0x%X was shown before the activation and mgr.IsShown() is false after it; AC#2 says a second launch brings the running instance's UI to the front, and hiding it is the D1 defect", hwndBefore)
	}
	if afterHwnd != hwndBefore {
		t.Errorf("AC#2-a RED: the second launch changed the window handle (0x%X -> 0x%X); an already-visible panel must stay the window the operator was looking at", hwndBefore, afterHwnd)
	}
	if !isLiveWindow(hwndBefore) {
		t.Errorf("AC#2-a RED: the panel's window 0x%X is no longer a live window after the activation", hwndBefore)
	}
	// AC#2-b: it must be the foreground window afterwards.
	if fgAfter != hwndBefore {
		t.Errorf("AC#2-b RED (the second launch left the panel in the background): the panel window 0x%X is shown=%v (alive=%v) but the foreground window is 0x%X, not the panel (the prior window was 0x%X, SetForegroundWindow answered %d); AC#2's letter is 'bring the existing instance's UI to the front', so this red says the raise hop did not run", hwndBefore, shownAfter, isLiveWindow(hwndBefore), fgAfter, uintptr(prior), raiseRet)
	}
}
