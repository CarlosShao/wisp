//go:build windows && winlive

package main

// 319-v2 out-of-repo probe (injected with -overlay as
// cmd/wisp/probe_319v2_shapes_winlive_test.go; never written into the tracked
// tree). Ticket 319 AC#2's other two shapes, through the same production
// second-launch hop the visible probe drives:
//
//   - TestProbe319v2HiddenPanelComesForward (AC#2-c): the panel window EXISTS but
//     is hidden (opened through the gesture road, then hidden through the
//     production toggle request - the tray/hotkey semantics, untouched). The
//     second launch must show it again AND leave it as the foreground window.
//   - TestProbe319v2DeadThreadRebuildsAndRaises (AC#2-d): the close chain is
//     planted; the second launch must rebuild a NEW window, shown AND in front.
//
// Shape c's prior window is built on the panel thread (the thread is alive);
// shape d's is built on this goroutine, because the panel thread's own windows
// die with it. In d the probe reports where the foreground actually was at
// signal time instead of failing on the prior install: the assertion is about
// the REBUILT window taking the front. Layouts are per-run; no second
// `wisp --env dev` is started.

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

// startSecondLaunchHop319v2 boots a runtime with this run's own layout, wires
// the production activation hook to the holder, and returns the receipts channel.
func startSecondLaunchHop319v2(t *testing.T, h *panelHolder, tag string) (chan string, func()) {
	t.Helper()
	reg := observe.NewRegistry()
	layout := proc.Layout{
		Env:               buildinfo.EnvProd,
		DataDir:           t.TempDir(),
		MutexName:         fmt.Sprintf(`Local\wisp-test-319v2-%s-%d-%d`, tag, os.Getpid(), time.Now().UnixNano()),
		ActivateEventName: fmt.Sprintf(`Local\wisp-test-319v2-%s-%d-%d.activate`, tag, os.Getpid(), time.Now().UnixNano()),
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
	return receipts, func() {
		if err := proc.SignalExistingInstance(rt.Layout.ActivateEventName); err != nil {
			t.Fatalf("SignalExistingInstance: %v", err)
		}
	}
}

func TestProbe319v2HiddenPanelComesForward(t *testing.T) {
	h, mgr := startPanelHolder319(t)

	if !h.toggle("319v2-probe-c-open") {
		t.Fatalf("the holder refused the open gesture")
	}
	waitPanelTrue(t, "the panel window to be shown", func() bool { return mgr.IsShown() })
	hwnd := mgr.windowHandle()
	if hwnd == 0 {
		t.Fatalf("the open reported a shown window with HWND 0")
	}
	// Hide it through the production toggle request (an open panel gets hidden).
	if !h.toggle("319v2-probe-c-hide") {
		t.Fatalf("the holder refused the hide gesture")
	}
	waitPanelTrue(t, "the panel to hide", func() bool { return !mgr.IsShown() })
	aliveWhileHidden := isLiveWindow(hwnd)

	receipts, signal := startSecondLaunchHop319v2(t, h, "hidden")

	prior := editorOnPanelThread(t, h.rp)
	if fg := uintptr(windows.GetForegroundWindow()); fg != uintptr(prior) {
		t.Fatalf("the prior window 0x%X did not take the foreground (the foreground is 0x%X); failed measurement", uintptr(prior), fg)
	}
	t.Logf("BEFORE: hwnd=0x%X mgr.IsShown=false alive-while-hidden=%v; prior=0x%X foreground=0x%X", hwnd, aliveWhileHidden, uintptr(prior), uintptr(windows.GetForegroundWindow()))

	signal()
	var receipt string
	select {
	case receipt = <-receipts:
	case <-time.After(panelThreadWait):
		t.Fatalf("the activation hook never answered within %v", panelThreadWait)
	}
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
	t.Logf("AFTER: receipt=%q; mgr.IsShown=%v hwnd=0x%X (same=%v) alive=%v; foreground=0x%X (panel=%v, prior=0x%X); SetForegroundWindow answered %d", receipt, shownAfter, afterHwnd, afterHwnd == hwnd, isLiveWindow(afterHwnd), fgAfter, fgAfter == afterHwnd, uintptr(prior), raiseRet)

	if !shownAfter {
		t.Errorf("AC#2-c RED: the hidden panel did not come back on screen (mgr.IsShown=false after the second launch)")
	}
	if !isLiveWindow(afterHwnd) {
		t.Errorf("AC#2-c RED: the panel window handle 0x%X is not a live window after the second launch", afterHwnd)
	}
	if fgAfter != afterHwnd {
		t.Errorf("AC#2-c RED: the panel is shown (hwnd=0x%X) but the foreground window is 0x%X, not the panel (prior=0x%X, SetForegroundWindow answered %d): the second launch must bring the existing instance's UI to the front", afterHwnd, fgAfter, uintptr(prior), raiseRet)
	}
}

func TestProbe319v2DeadThreadRebuildsAndRaises(t *testing.T) {
	h, mgr := startPanelHolder319(t)

	if !h.toggle("319v2-probe-d-open") {
		t.Fatalf("the holder refused the first open gesture")
	}
	waitPanelTrue(t, "the first panel window to be shown", func() bool { return mgr.IsShown() })
	old := mgr.windowHandle()
	if old == 0 {
		t.Fatalf("the first show reported a shown window with HWND 0")
	}
	dead := h.rp

	// A prior that survives the panel thread: built on THIS goroutine, not on the
	// panel thread (a locked thread's windows die with it), so it can still be the
	// foreground after the death.
	prior, err := createEditorWindow("wisp 319v2 d prior")
	if err != nil {
		t.Skipf("the ruler could not build its own prior window (failed measurement, not a product verdict): %v", err)
	}
	if prior == 0 {
		t.Fatalf("CreateWindowExW returned 0 with no error code")
	}
	t.Cleanup(func() { destroyEditorWindow(prior) })
	pnlShowWindow.Call(uintptr(prior), pnlSwShow)
	pnlUpdateWindow.Call(uintptr(prior))

	// The 15:41:45 chain on the real window (the pair the system sends for the X).
	r, _, callErr := t33r9NailPostMessage.Call(old, pan319WmSysCommand, pan319ScClose, 0)
	if r == 0 {
		t.Fatalf("PostMessageW(hwnd 0x%X, WM_SYSCOMMAND, SC_CLOSE) was refused: %v", old, callErr)
	}
	waitPanelTrue(t, "the panel thread to end after its window was closed", dead.isFinished)

	// Put the surviving prior in front where possible, and REPORT where the
	// foreground really is - the assertion below is about the rebuilt window.
	pnlSetForeground.Call(uintptr(prior))
	fgBefore := uintptr(windows.GetForegroundWindow())
	t.Logf("BEFORE: old hwnd=0x%X closed (alive after close=%v); mgr.IsCreated=%v mgr hwnd=0x%X; prior=0x%X foreground-before-signal=0x%X (prior-in-front=%v)", old, isLiveWindow(old), mgr.IsCreated(), mgr.windowHandle(), uintptr(prior), fgBefore, fgBefore == uintptr(prior))

	receipts, signal := startSecondLaunchHop319v2(t, h, "dead")
	signal()

	waitPanelTrue(t, "the activation to rebuild a fresh panel window", func() bool {
		return mgr.IsCreated() && mgr.IsShown() && mgr.windowHandle() != 0 && mgr.windowHandle() != old
	})
	neu := mgr.windowHandle()
	fgAfter := uintptr(windows.GetForegroundWindow())
	raiseRet := mgr.LastRaiseSetForeground()
	var receipt string
	select {
	case receipt = <-receipts:
	case <-time.After(panelThreadWait):
		t.Fatalf("the activation hook never answered within %v", panelThreadWait)
	}
	t.Logf("AFTER: receipt=%q; rebuilt hwnd=0x%X (new=%v) mgr.IsCreated=%v mgr.IsShown=%v alive=%v; old window alive=%v; foreground=0x%X (rebuilt-in-front=%v); SetForegroundWindow answered %d", receipt, neu, neu != old, mgr.IsCreated(), mgr.IsShown(), isLiveWindow(neu), isLiveWindow(old), fgAfter, fgAfter == neu, raiseRet)

	if neu == old {
		t.Errorf("AC#2-d RED: the rebuild reports the OLD window handle 0x%X; the close chain destroyed that window, so a rebuild reusing the dead handle is not a rebuild", old)
	}
	if !isLiveWindow(neu) {
		t.Errorf("AC#2-d RED: the rebuilt window handle 0x%X is not a live window", neu)
	}
	if isLiveWindow(old) {
		t.Errorf("AC#2-d RED: the closed window 0x%X is alive again; the close chain should have destroyed it", old)
	}
	if fgAfter != neu {
		t.Errorf("AC#2-d RED: the rebuilt panel (hwnd=0x%X, shown=%v) is not the foreground window (foreground=0x%X, prior=0x%X, SetForegroundWindow answered %d): a second launch after a dead thread must rebuild the window AND bring it to the front", neu, mgr.IsShown(), fgAfter, uintptr(prior), raiseRet)
	}
}
