//go:build windows && winlive

package ball

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// Ticket 65 hard-judgment probe (docs/evidence/s1/65-jchv-controller-leak.md
// §3): after the fix, production destroy() must actually CLOSE the WebView2
// controller, and the measurable consequence is that no msedgewebview2.exe
// process parented by this test process survives a run of create/destroy
// rounds. This file is deliberately self-contained - it shares no helpers
// with orb_live_windows_test.go - so the leak judgment never depends on
// another leg's in-flight churn in that file.

// orbCloseProbeWebView2Children counts every msedgewebview2.exe process whose
// PARENT is this test process. The parent-PID plus image-name double match is
// the same discipline as the suite's reaper: no other application's WebView2
// processes can ever match, this only ever sees children this process spawned.
func orbCloseProbeWebView2Children(t *testing.T) int {
	t.Helper()
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		t.Fatalf("probe: process snapshot failed: %v", err)
	}
	defer windows.CloseHandle(snap)
	mine := windows.GetCurrentProcessId()
	n := 0
	var pe windows.ProcessEntry32
	pe.Size = uint32(unsafe.Sizeof(pe))
	for err := windows.Process32First(snap, &pe); err == nil; err = windows.Process32Next(snap, &pe) {
		if pe.ParentProcessID == mine && windows.UTF16ToString(pe.ExeFile[:]) == "msedgewebview2.exe" {
			n++
		}
	}
	return n
}

// orbCloseProbeReapOwnWebView2 terminates every own-parented
// msedgewebview2.exe (leftovers from an earlier crashed run would otherwise
// sit in the baseline) and reports how many it killed.
func orbCloseProbeReapOwnWebView2(t *testing.T) int {
	t.Helper()
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		t.Fatalf("probe: process snapshot failed: %v", err)
	}
	defer windows.CloseHandle(snap)
	mine := windows.GetCurrentProcessId()
	killed := 0
	var pe windows.ProcessEntry32
	pe.Size = uint32(unsafe.Sizeof(pe))
	for err := windows.Process32First(snap, &pe); err == nil; err = windows.Process32Next(snap, &pe) {
		if pe.ParentProcessID != mine || windows.UTF16ToString(pe.ExeFile[:]) != "msedgewebview2.exe" {
			continue
		}
		if h, oerr := windows.OpenProcess(windows.PROCESS_TERMINATE, false, pe.ProcessID); oerr == nil {
			if terr := windows.TerminateProcess(h, 1); terr == nil {
				killed++
			}
			windows.CloseHandle(h)
		}
	}
	return killed
}

// TestOrbControllerCloseStopsBrowserProcessLeak is the re-runnable measurement
// behind the fix judgment: five full ball create/destroy rounds, then the
// process table must show ZERO net growth of own-parented msedgewebview2
// processes. Before the fork's controller.Close this exact shape leaked one
// browser process set per round (the measured 40.6%->0 pixel collapse across
// count=10); after it the count returns to baseline.
func TestOrbControllerCloseStopsBrowserProcessLeak(t *testing.T) {
	if killed := orbCloseProbeReapOwnWebView2(t); killed > 0 {
		t.Logf("probe: reaped %d leftover own-parented msedgewebview2 process(es) from earlier runs", killed)
		time.Sleep(2 * time.Second) // terminated entries need a beat to leave the table
	}
	baseline := orbCloseProbeWebView2Children(t)

	store, err := OpenPositionStore(filepath.Join(t.TempDir(), "ballpos.json"))
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	title := fmt.Sprintf("Wisp Leak Probe %d", windows.GetCurrentProcessId())

	const rounds = 5
	for i := 1; i <= rounds; i++ {
		b, err := New(Options{
			SizePx:      BallSizeDefaultPx,
			Initial:     statemachine.StateSleeping,
			Store:       store,
			WindowTitle: title,
		})
		if err != nil {
			t.Fatalf("round %d: New: %v", i, err)
		}
		deadline := time.Now().Add(15 * time.Second)
		for b.DebugOrbHWND() == 0 {
			if time.Now().After(deadline) {
				b.Close()
				t.Fatalf("round %d: orb did not attach within 15s (WebView2 runtime missing or Embed refused)", i)
			}
			time.Sleep(100 * time.Millisecond)
		}
		// Let the browser process tree finish spawning so the round is a real
		// leak opportunity, then tear the whole ball down through the
		// production Close path (the only path under judgment).
		time.Sleep(300 * time.Millisecond)
		b.Close()
	}

	// controller.Close retires the browser process tree; give it a bounded
	// grace window, then require the table to be back at baseline. Anything
	// still parented by us past the budget is a leak this probe exists to
	// catch.
	var after int
	deadline := time.Now().Add(15 * time.Second)
	for {
		after = orbCloseProbeWebView2Children(t)
		if after <= baseline || time.Now().After(deadline) {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if after > baseline {
		t.Fatalf("WebView2 controller leak: %d create/destroy rounds left %d msedgewebview2 child process(es) parented by this test process (baseline %d) - destroy() is not closing the controller",
			rounds, after, baseline)
	}
	t.Logf("probe: %d create/destroy rounds, own-parented msedgewebview2 baseline=%d after=%d (no growth - controller Close verified)",
		rounds, baseline, after)
}
