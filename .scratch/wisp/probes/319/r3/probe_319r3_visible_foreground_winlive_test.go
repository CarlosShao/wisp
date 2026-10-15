//go:build windows && winlive

package main

// 319-r3 out-of-repo probe (injected with -overlay as
// cmd/wisp/probe_319r3_visible_foreground_winlive_test.go; never committed to the
// tree). Question it answers: on the fixed tree, with the panel ALREADY visible
// and the prior window in front, where does the foreground land after the
// second-launch activation, and would an explicit SetForegroundWindow or
// ShowWindow(SW_SHOW) on the same HWND raise it?
//
// This is a measurement, not a verdict: it opens a real window pair, reads
// GetForegroundWindow after each step, and prints. Production is untouched.

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

func TestProbe319r3VisibleForeground(t *testing.T) {
	h, mgr := startPanelHolder319(t)

	if !h.toggle("319r3-probe-setup") {
		t.Fatalf("the holder refused the setup gesture")
	}
	waitPanelTrue(t, "the setup panel window to be shown", func() bool { return mgr.IsShown() })
	hwnd := mgr.windowHandle()
	if hwnd == 0 {
		t.Fatalf("shown window with HWND 0")
	}

	reg := observe.NewRegistry()
	layout := proc.Layout{
		Env:               buildinfo.EnvProd,
		DataDir:           t.TempDir(),
		MutexName:         fmt.Sprintf(`Local\wisp-test-319r3probe-%d-%d`, os.Getpid(), time.Now().UnixNano()),
		ActivateEventName: fmt.Sprintf(`Local\wisp-test-319r3probe-%d-%d.activate`, os.Getpid(), time.Now().UnixNano()),
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

	prior := editorOnPanelThread(t, h.rp)
	if fg := uintptr(windows.GetForegroundWindow()); fg != uintptr(prior) {
		t.Fatalf("prior window 0x%X did not take the foreground (fg=0x%X); measurement impossible", uintptr(prior), fg)
	}

	if err := proc.SignalExistingInstance(rt.Layout.ActivateEventName); err != nil {
		t.Fatalf("SignalExistingInstance: %v", err)
	}
	select {
	case <-receipts:
	case <-time.After(panelThreadWait):
		t.Fatalf("no activation receipt within %v", panelThreadWait)
	}

	// Deterministic "the activation's request was served": a sentinel posted after
	// it through the same route (post -> library Dispatch -> Run), so it runs
	// after the activation closure whatever that closure did.
	served := make(chan struct{})
	if !h.rp.post(func() { close(served) }) {
		t.Fatalf("the sentinel post was refused by the panel thread")
	}
	select {
	case <-served:
	case <-time.After(panelThreadWait):
		t.Fatalf("the sentinel was not served within %v", panelThreadWait)
	}

	afterShown := mgr.IsShown()
	afterHwnd := mgr.windowHandle()
	fgAfterShow := uintptr(windows.GetForegroundWindow())
	t.Logf("PROBE A (fixed tree, visible panel): shown=%v hwnd=0x%X same-hwnd=%v alive=%v fg-after-Show=0x%X panel-was-fg=%v prior=0x%X",
		afterShown, afterHwnd, afterHwnd == hwnd, isLiveWindow(afterHwnd), fgAfterShow, fgAfterShow == afterHwnd, uintptr(prior))
	time.Sleep(1 * time.Second)
	fgAfterSettle := uintptr(windows.GetForegroundWindow())
	t.Logf("PROBE A2 (one second later, nothing else called): fg=0x%X panel-was-fg=%v", fgAfterSettle, fgAfterSettle == afterHwnd)

	// Put the prior window back in front, then measure ONE candidate call at a
	// time from a clean state - otherwise each probe pollutes the next reading.
	putPriorBack := func() {
		t.Helper()
		pnlSetForeground.Call(uintptr(prior))
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			if uintptr(windows.GetForegroundWindow()) == uintptr(prior) {
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
		t.Fatalf("the prior window 0x%X would not take the foreground back; the probes below cannot be isolated", uintptr(prior))
	}

	// C: ShowWindow(SW_SHOW) alone - the exact call mgr.Show's re-show branch makes.
	putPriorBack()
	rC, _, errC := pnlShowWindow.Call(afterHwnd, pnlSwShow)
	time.Sleep(300 * time.Millisecond)
	fgAfterSW := uintptr(windows.GetForegroundWindow())
	t.Logf("PROBE C (ShowWindow(SW_SHOW) alone, visible HWND, from a clean prior-in-front state): ret=%d err=%v fg-after-SW_SHOW=0x%X panel-was-fg=%v",
		rC, errC, fgAfterSW, fgAfterSW == afterHwnd)

	// B: SetForegroundWindow alone - what mgr.Show skips when already shown.
	putPriorBack()
	rB, _, errB := pnlSetForeground.Call(afterHwnd)
	time.Sleep(300 * time.Millisecond)
	fgAfterSet := uintptr(windows.GetForegroundWindow())
	t.Logf("PROBE B (SetForegroundWindow alone, from a clean prior-in-front state): ret=%d err=%v fg-after-SetForeground=0x%X panel-was-fg=%v",
		rB, errB, fgAfterSet, fgAfterSet == afterHwnd)
}
