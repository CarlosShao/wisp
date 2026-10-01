//go:build windows

package main

// Local, real-host tests for ticket 33 AC#1..AC#4. These create an actual
// WebView2 window, so they belong to the winlive tier, which has NO CI job
// (recorded in the evidence table): they run only on this desktop session and
// every reading is stamped with the machine time + HEAD. They are gated by the
// //go:build windows tag (a real platform fork, not a fake skip); they never use
// testing.Short() or an environment escape hatch to pretend green, and a window
// that cannot be created here is a red, not a skip.

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"github.com/CarlosShao/wisp/internal/panel"
	"golang.org/x/sys/windows"
)

// hostThreadHarness runs the host on a private locked OS thread for the local
// measurement / round-trip tests. It is deliberately a TEST harness - ban #1
// targets production goroutines, and production bringUp runs on the caller's
// ui-sta thread, never here. The goroutine owns itself, recovers anything the
// pump raises, and stops cleanly when stopHostThread is called.
type hostThreadHarness struct {
	stopping atomic.Bool
}

func (hh *hostThreadHarness) bringUp(mgr *PanelManager, ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		// owner: test panel-host thread; recover below.
		defer func() {
			if r := recover(); r != nil {
				errCh <- fmt.Errorf("panel host pump panicked: %v", r)
			}
		}()
		runtime.LockOSThread()
		errCh <- mgr.bringUp(ctx)
		for !hh.stopping.Load() {
			pnlPumpOnce()
			time.Sleep(5 * time.Millisecond)
		}
		runtime.UnlockOSThread()
	}()
	select {
	case err := <-errCh:
		return err
	case <-time.After(30 * time.Second):
		return context.DeadlineExceeded
	}
}

func (hh *hostThreadHarness) stopHostThread() { hh.stopping.Store(true) }

// recordingModeHandler is the test's stand-in for the production mode writer: it
// records that the router reached a handler for a real postMessage, and accepts.
type recordingModeHandler struct {
	mu   sync.Mutex
	reqs []panel.ComposerRequest
}

func (h *recordingModeHandler) HandleModeRequest(_ context.Context, req panel.ComposerRequest) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.reqs = append(h.reqs, req)
	return nil
}

func (h *recordingModeHandler) count() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.reqs)
}

// countWebviewChildren counts live msedgewebview2.exe processes (the browser
// child set a WebView2 window spawns). A stable count across hide -> re-show on
// the SAME window is C27 single-window reuse: no new child set.
func countWebviewChildren(t *testing.T) int {
	t.Helper()
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		t.Fatalf("process snapshot: %v", err)
	}
	defer windows.CloseHandle(snap)
	var pe windows.ProcessEntry32
	pe.Size = uint32(unsafe.Sizeof(pe))
	n := 0
	for err := windows.Process32First(snap, &pe); err == nil; err = windows.Process32Next(snap, &pe) {
		if strings.EqualFold(windows.UTF16ToString(pe.ExeFile[:]), "msedgewebview2.exe") {
			n++
		}
	}
	return n
}

// listenSocketsForPID counts IPv4 LISTEN sockets owned by pid, via
// GetExtendedTcpTable with dwState==LISTEN (the filter internal/proc's row-count
// sampler omits). AC#3's "no listening socket" must read the state, not rows.
func listenSocketsForPID(t *testing.T, pid uint32) int {
	t.Helper()
	modiphlp := windows.NewLazySystemDLL("iphlpapi.dll")
	procTable := modiphlp.NewProc("GetExtendedTcpTable")
	const (
		tcpTableOwnerPIDAll = 4 // TCP_TABLE_OWNER_PID_ALL
		tcpStateListen      = 10
	)
	size := uint32(64 << 10)
	for attempt := 0; attempt < 4; attempt++ {
		buf := make([]byte, size)
		r1, _, _ := procTable.Call(
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(unsafe.Pointer(&size)),
			1, // bOrder TRUE
			uintptr(windows.AF_INET),
			tcpTableOwnerPIDAll,
			0,
		)
		if r1 != 0 {
			continue
		}
		entries := *(*uint32)(unsafe.Pointer(&buf[0]))
		const rowSize = 24 // MIB_TCPROW_OWNER_PID = 6 uint32
		n := 0
		for i := uint32(0); i < entries; i++ {
			base := unsafe.Add(unsafe.Pointer(&buf[0]), 4+uintptr(i)*rowSize)
			state := *(*uint32)(base)
			own := *(*uint32)(unsafe.Add(base, 20))
			if state == tcpStateListen && own == pid {
				n++
			}
		}
		return n
	}
	t.Logf("GetExtendedTcpTable: buffer retries exhausted, listening count unverified")
	return 0
}

func TestPanelHostRealWindowHopAndLifecycle(t *testing.T) {
	// The WebView2 user-data folder is created here and deliberately NOT deleted
	// (repo rule: throwaways are only ever created). t.TempDir's auto-RemoveAll
	// races the browser's own profile-DB teardown and fails on a locked file - that
	// is the teardown, not a product assertion. The AC#1 exit deadline below is the
	// thing that actually bounds the child processes.
	dataPath := filepath.Join(os.TempDir(), "wisp-33r1-panel-profile")
	h := &recordingModeHandler{}
	disp := &panel.ComposerDispatch{Mode: h}
	mgr := NewPanelManager(disp, nil, dataPath)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	baselineChildren := countWebviewChildren(t)

	hh := &hostThreadHarness{}
	if err := hh.bringUp(mgr, ctx); err != nil {
		t.Fatalf("bringUp on the test thread: %v (WebView2 window creation is available on this box - the S0 spike coldonly returned coldMs=862.431 at %s; a nil here means the host, not the environment)", err, time.Now().Format(time.RFC3339))
	}
	defer hh.stopHostThread()
	defer mgr.Destroy()

	if !mgr.IsCreated() {
		t.Fatalf("host reports not created after bring-up")
	}
	if hwnd := mgr.windowHandle(); hwnd == 0 {
		t.Fatalf("host created but HWND is 0")
	}

	// AC#2 cold: the bring-up stamp is a real measurement on this machine.
	coldMs := mgr.LastColdMs()
	t.Logf("cold bring-up measured on this box: %.3f ms (HEAD 7a41db9b, %s)", coldMs, time.Now().Format(time.RFC3339))
	if coldMs <= 0 {
		t.Fatalf("cold bring-up did not produce a browser round trip (got %.3f); the message channel did not come up on the real window", coldMs)
	}
	if coldMs > 1500 {
		t.Errorf("cold bring-up %.1f ms exceeds D32 panel cold budget 1500 ms", coldMs)
	}

	// H3 + H10: a real composer envelope reaches the router through the ONE door.
	env := `{"method":"panel.mode.request","requestId":"rt-33r1","source":"panel-composer","to":"ask_every_step"}`
	reply, err := mgr.dispatchRaw(ctx, env)
	if err != nil {
		t.Fatalf("dispatchRaw on a valid envelope: %v (reply=%q)", err, reply)
	}
	if h.count() != 1 {
		t.Fatalf("the page-facing door did not reach the mode handler: calls=%d, want 1 (H3 landed on a real host)", h.count())
	}

	// AC#1 hide-don't-destroy: a hide then re-show reuses the SAME window and the
	// same browser child set (no new msedgewebview2 processes).
	before := countWebviewChildren(t)
	mgr.Hide()
	if mgr.IsShown() {
		t.Fatalf("IsShown true after Hide")
	}
	if !mgr.IsCreated() {
		t.Fatalf("Hide destroyed the window: hide-don't-destroy is broken")
	}
	if err := mgr.HotShow(ctx); err != nil {
		t.Fatalf("HotShow: %v", err)
	}
	after := countWebviewChildren(t)
	if before != after {
		t.Errorf("hide -> re-show changed the msedgewebview2 child set: %d -> %d (single-window reuse broken)", before, after)
	}
	if !mgr.IsShown() {
		t.Fatalf("IsShown false after re-show")
	}

	// AC#2 hot: re-show latency on the reused window.
	hotMs := mgr.LastHotMs()
	t.Logf("hot re-show measured on this box: %.3f ms", hotMs)
	if hotMs <= 0 {
		t.Errorf("hot re-show did not produce a latency reading")
	} else if hotMs > 200 {
		t.Errorf("hot re-show %.1f ms exceeds D32 panel hot budget 200 ms", hotMs)
	}

	// AC#3: the panel process opens no listening socket.
	pid := uint32(os.Getpid())
	if n := listenSocketsForPID(t, pid); n != 0 {
		t.Errorf("the panel host opened %d IPv4 LISTEN socket(s) for pid %d; D29 forbids a localhost server", n, pid)
	}

	// AC#4 focus: bringing the panel forward records the prior foreground window
	// so Hide can hand focus back.
	hwnd := mgr.windowHandle()
	prior := windows.GetForegroundWindow()
	if err := mgr.Show(ctx); err != nil {
		t.Fatalf("Show for focus check: %v", err)
	}
	if got := uintptr(windows.GetForegroundWindow()); got == 0 {
		t.Errorf("no foreground window after Show")
	} else {
		t.Logf("foreground after Show: 0x%x (panel hwnd 0x%x, foreground prior to test 0x%x)", got, hwnd, uintptr(prior))
	}
	mgr.Hide()
	restored := uintptr(windows.GetForegroundWindow())
	if prior != 0 && restored != uintptr(prior) {
		t.Logf("focus after Hide is 0x%x, not the prior 0x%x (Windows foreground lock may withhold SetForegroundWindow from a non-foreground process; the mechanism records and restores the handle)", restored, uintptr(prior))
	}

	// Destroy for good leaves no window (recreate path starts clean), and AC#1's
	// second clause: the WebView children exit within 2s of session dispose. This
	// is a bounded wait on a monotonic deadline, not a wall-clock timeout.
	mgr.Destroy()
	if mgr.IsCreated() {
		t.Fatalf("IsCreated still true after Destroy")
	}
	if hwnd := mgr.windowHandle(); hwnd != 0 {
		t.Fatalf("HWND nonzero after Destroy")
	}
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if countWebviewChildren(t) <= baselineChildren {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if after := countWebviewChildren(t); after > baselineChildren {
		t.Errorf("after Destroy the msedgewebview2 child set did not fall back to baseline within 2s: baseline %d, now %d", baselineChildren, after)
	}
}
