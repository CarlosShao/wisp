//go:build windows && winlive

package main

// AC#1's timing clause, moved out of the default tier by 33-r5 (orchestrator ruling
// 10-01 12:22, form B).
//
// What was moved and why: "the WebView2 children THIS process started exit within
// 2s of Destroy" is the one clause in this family that cannot be decided from a
// single deterministic observation - 33-r4 measured it red in 1 of 6 whole-package
// runs and green in 5 of 5 isolated runs (its §② second final run and §⑤ 25), and
// both denominators moved in the same direction in the red run, so it is a load
// effect and not a denominator artifact.
//
// The bound is UNCHANGED at 2 seconds. Nothing was widened to 5s or 10s, no retry
// was added around the assertion, and the clause was not deleted. What changed is
// which tier it lives in, and the cost of that is stated here and in
// docs/evidence/s1/33-panel-host-c27-r5.md: winlive has NO CI job, so this clause
// is now measurable only on a desktop box and is never seen by the pipeline. Any
// table quoting it has to carry that sentence.
//
// The default-tier lifecycle test keeps the dimensions that ARE decidable from one
// sample: the same HWND across hide -> re-show, a single window, our tree holding
// browser children while the window is up, and Destroy leaving no window behind.

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/panel"
)

func TestPanelHostWebViewChildrenExitWithinTwoSeconds_WinLive(t *testing.T) {
	dataPath := filepath.Join(os.TempDir(), "wisp-33r1-panel-profile")
	h := &recordingModeHandler{}
	mgr := NewPanelManager(&panel.ComposerDispatch{Mode: h}, nil, dataPath)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	self := uint32(os.Getpid())
	baseline := readWebviewTree(t, self)

	hh := &hostThreadHarness{}
	if err := hh.bringUp(mgr, ctx); err != nil {
		t.Fatalf("bringUp for the winlive exit clause: %v", err)
	}
	defer hh.stopHostThread()

	up := readWebviewTree(t, self)
	if up.TreeWebview < 1 {
		t.Fatalf("the window is up but our tree holds %d msedgewebview2 process(es) - there is nothing for the 2s exit clause to be about (tree pids %d, machine-wide %d)",
			up.TreeWebview, len(up.TreePIDs), up.MachineNamed)
	}
	mgr.Destroy()

	// Same bounded, monotonic wait and the same 2s bound 33-r4 shipped. The
	// denominator is this process tree, not the machine (33-v1 §A#26).
	deadline := time.Now().Add(2 * time.Second)
	final := readWebviewTree(t, self)
	for final.TreeWebview > baseline.TreeWebview && time.Now().Before(deadline) {
		time.Sleep(50 * time.Millisecond)
		final = readWebviewTree(t, self)
	}
	t.Logf("winlive exit clause: our tree webview baseline %d -> peak %d -> now %d (tree pids %d) | machine-wide %d -> %d",
		baseline.TreeWebview, up.TreeWebview, final.TreeWebview, len(final.TreePIDs), baseline.MachineNamed, final.MachineNamed)
	if final.TreeWebview > baseline.TreeWebview {
		t.Errorf("after Destroy the WebView2 children this process started did not exit within 2s: our tree went baseline %d -> now %d (tree pids %d). The machine-wide count (%d -> %d) is reported only and is NOT the denominator",
			baseline.TreeWebview, final.TreeWebview, len(final.TreePIDs), baseline.MachineNamed, final.MachineNamed)
	}
}
