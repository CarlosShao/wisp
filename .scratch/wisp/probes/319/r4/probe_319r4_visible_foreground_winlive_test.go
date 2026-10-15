//go:build windows && winlive

package main

// 319-r4 out-of-repo probe (injected with -overlay as
// cmd/wisp/probe_319r4_visible_foreground_winlive_test.go; never committed to the
// tree). The question it answers, on the PRE-CHANGE tree and on the fixed tree
// alike: with the panel ALREADY visible and the prior window in front, where does
// the production second-launch hop (SignalExistingInstance -> RunEventLoop ->
// secondLaunchActivation319 -> panelHolder.show) leave the foreground - does the
// committed pin's assertion `GetForegroundWindow() == the panel's HWND` hold?
//
// It must NOT hold on the pre-change tree: that is this leg's "unfixed tree, does
// the ruler ring" reading for ticket 319 AC#2's raise half. Run against the fixed
// tree (the same probe, no pre-change mapping) the same assertion must hold. The
// probe deliberately references only symbols that exist on BOTH sides of the
// change (startPanelHolder319, panelHolder.toggle, the production hook,
// editorOnPanelThread, isLiveWindow), so one file can be the ruler for both
// trees; the two runs and their readings are in .scratch/wisp/probes/319/r4/.

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

func TestProbe319r4VisibleForeground(t *testing.T) {
	h, mgr := startPanelHolder319(t)

	// The operator's panel is already on screen, opened through the ordinary
	// gesture road (toggle's show half).
	if !h.toggle("319r4-probe-setup") {
		t.Fatalf("the holder refused the setup gesture")
	}
	waitPanelTrue(t, "the setup panel window to be shown", func() bool { return mgr.IsShown() })
	hwnd := mgr.windowHandle()
	if hwnd == 0 {
		t.Fatalf("shown window with HWND 0")
	}

	// The runtime + event loop + production hook, the in-process shape the
	// committed cases build (no second dev process: a second `wisp --env dev`
	// would signal the owner's instance).
	reg := observe.NewRegistry()
	layout := proc.Layout{
		Env:               buildinfo.EnvProd,
		DataDir:           t.TempDir(),
		MutexName:         fmt.Sprintf(`Local\wisp-test-319r4probe-%d-%d`, os.Getpid(), time.Now().UnixNano()),
		ActivateEventName: fmt.Sprintf(`Local\wisp-test-319r4probe-%d-%d.activate`, os.Getpid(), time.Now().UnixNano()),
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

	// Where the operator came from: a window of this process in front of the
	// visible panel, so "was the panel raised" is decidable.
	prior := editorOnPanelThread(t, h.rp)
	if fg := uintptr(windows.GetForegroundWindow()); fg != uintptr(prior) {
		t.Fatalf("the prior window 0x%X did not take the foreground (fg=0x%X); measurement impossible", uintptr(prior), fg)
	}

	// The second launch: the identical call the ErrAlreadyRunning branch makes.
	if err := proc.SignalExistingInstance(rt.Layout.ActivateEventName); err != nil {
		t.Fatalf("SignalExistingInstance: %v", err)
	}
	select {
	case <-receipts:
	case <-time.After(panelThreadWait):
		t.Fatalf("no activation receipt within %v", panelThreadWait)
	}

	// Deterministic "the activation's requests were served": a sentinel posted
	// after them through the same route is FIFO behind them.
	served := make(chan struct{})
	if !h.rp.post(func() { close(served) }) {
		t.Fatalf("the sentinel post was refused by the panel thread")
	}
	select {
	case <-served:
	case <-time.After(panelThreadWait):
		t.Fatalf("the sentinel was not served within %v", panelThreadWait)
	}

	fg := uintptr(windows.GetForegroundWindow())
	t.Logf("PROBE 319-r4 (visible panel, production second-launch hop): hwnd=0x%X shown=%v alive=%v; foreground=0x%X panel-was-fg=%v prior=0x%X",
		hwnd, mgr.IsShown(), isLiveWindow(hwnd), fg, fg == hwnd, uintptr(prior))
	if fg != hwnd {
		t.Errorf("the panel was left in the background: the panel window 0x%X is shown (alive=%v) but the foreground window is 0x%X, not the panel (prior window 0x%X) - the committed pin's assertion `GetForegroundWindow() == the panel's HWND` is RED on this tree", hwnd, isLiveWindow(hwnd), fg, uintptr(prior))
	}
}
