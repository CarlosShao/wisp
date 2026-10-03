//go:build windows && winlive

package main

// Ticket 255 AC#4 - the real-window half, measured on this desktop session.
//
// WHAT panel_geometry_255_test.go cannot say. Its four rulers prove that the
// number resolved from [panel] is the number handed to webview2.NewWithOptions,
// which is the whole of what a CI box can decide. They do not say that Win32 then
// sizes a window with it, and 票 255 AC#4 names that second reading as belonging
// to the 只有本机可量 family. This file is that reading, and it is tagged winlive so
// it is a deliberate local measurement: the default `go test ./cmd/wisp/` does not
// compile it, so no other leg's package run can be poisoned by a window it opens.
//
// WHY IT NAILS ITS OWN TEARDOWN, IN THE SAME BREATH AS ITS ASSERTION. The血账 this
// ticket carries is a single real window poisoning the NEXT case that lands on the
// same M: Destroy only posts, and UnlockOSThread would hand back a thread still
// owing those messages. Three things here defuse it, and the third one is asserted
// rather than trusted:
//   - every host call (create, destroy, re-create) runs on the harness thread
//     through hostThreadHarness.call, never from this goroutine;
//   - hostThreadHarness.runOnThread never calls runtime.UnlockOSThread, by name and
//     with the reason written at its own site (panel_host_windows_test.go:90-102),
//     so the thread and whatever it still owes die with the goroutine;
//   - after the final Destroy this case counts the windows still on that thread by
//     handle, and a non-zero count is a leak it reports as its OWN failure.
//
// DPI, stated so the number is read correctly: the assertion compares the two
// SAMPLES against the two REQUESTED values through the scale factor the first
// sample reveals. On a 100% display that is the requested pixel width exactly; on
// a scaled display it is the same requested width times the same factor, and what
// would be a false claim - "420 became 900" - is still the claim being made.

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	g255User32                   = windows.NewLazySystemDLL("user32.dll")
	g255GetWindowRect            = g255User32.NewProc("GetWindowRect")
	g255GetWindowThreadProcessID = g255User32.NewProc("GetWindowThreadProcessId")
)

// rect255 mirrors RECT.
type rect255 struct{ left, top, right, bottom int32 }

func windowRect255(hwnd uintptr) (rect255, error) {
	var r rect255
	rc, _, callErr := g255GetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&r)))
	if rc == 0 {
		return r, fmt.Errorf("GetWindowRect(0x%X) returned 0: %w", hwnd, callErr)
	}
	return r, nil
}

func windowThread255(hwnd uintptr) (uint32, error) {
	var pid uint32
	tid, _, callErr := g255GetWindowThreadProcessID.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	if tid == 0 {
		return 0, fmt.Errorf("GetWindowThreadProcessId(0x%X) returned 0: %w", hwnd, callErr)
	}
	return uint32(tid), nil
}

// TestTicket255RealWindowWidthFollowsTheConfig is AC#4's on-screen reading: it
// creates the panel window for real, destroys it, changes what the geometry source
// answers, creates it again, and compares what Win32 reported for the two windows.
// The second create is the path the ticket's ruling names - dispose then show, the
// only "read it again" this host has.
func TestTicket255RealWindowWidthFollowsTheConfig(t *testing.T) {
	const (
		// asked1 is 640, internal/config/schema.go's own [panel] width default;
		// asked2 is 900, a number no constant in this repository can produce, so
		// only the config path can put it on screen.
		asked1 = 640
		asked2 = 900
		askedH = 400
	)
	var (
		mu      sync.Mutex
		width   = asked1
		height  = askedH
		reads   int
		dataDir = t.TempDir()
	)
	src := func() (int, int) {
		mu.Lock()
		defer mu.Unlock()
		reads++
		return width, height
	}

	mgr := NewPanelManager(panelDispatchForGeometry255(), nil, filepath.Join(dataDir, "webview2"),
		withGeometrySource(src))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hh := &hostThreadHarness{}
	if err := hh.bringUp(mgr, ctx); err != nil {
		t.Fatalf("the first real create: %v", err)
	}
	defer hh.stopHostThread()
	if !mgr.IsCreated() {
		t.Fatal("bringUp reported success and IsCreated() is still false")
	}

	hwnd1 := mgr.windowHandle()
	if hwnd1 == 0 {
		t.Fatal("the created window has no handle to measure - this is a failed measurement, not a green")
	}
	tid, err := windowThread255(hwnd1)
	if err != nil {
		t.Fatalf("which thread owns the panel window: %v", err)
	}
	r1, err := windowRect255(hwnd1)
	if err != nil {
		t.Fatalf("the first window's rect: %v", err)
	}
	got1 := int(r1.right - r1.left)
	h1 := int(r1.bottom - r1.top)

	// Dispose, then ask for the panel again - the recreate path, on the thread that
	// owns the window.
	hh.call(func() { mgr.Destroy() })
	if mgr.IsCreated() {
		t.Fatal("Destroy left IsCreated() true, so the second Show would not create anything and the reading below would be about the first window")
	}
	// Destroy only POSTS WM_CLOSE, so a count taken on the same instruction is
	// always non-zero and would be a false red (measured: 3 windows, the panel plus
	// its WebView2 children, on the first run of this case). The nail is that the
	// pump dispatches it and the count reaches zero - and that is also the
	// precondition bringUp itself refuses without (staleCloseQueued), so a thread
	// still owing the close here would make the recreate below fail for the reason
	// the host documents instead of silently.
	if n := settleThreadWindows255(t, hh, tid, 3*time.Second); n != 0 {
		t.Errorf("after Destroy the panel thread (tid %d) still owns %d window(s) three seconds later: that is the leak shape this case exists to nail", tid, n)
	}

	mu.Lock()
	width, height = asked2, askedH
	mu.Unlock()

	var showErr error
	hh.call(func() { showErr = mgr.Show(ctx) })
	if showErr != nil {
		t.Fatalf("the recreate after dispose: %v", showErr)
	}
	hwnd2 := mgr.windowHandle()
	if hwnd2 == 0 {
		t.Fatal("the recreated window has no handle to measure")
	}
	r2, err := windowRect255(hwnd2)
	if err != nil {
		t.Fatalf("the recreated window's rect: %v", err)
	}
	got2 := int(r2.right - r2.left)
	h2 := int(r2.bottom - r2.top)

	// The scale factor the FIRST sample reveals, applied to the SECOND request. At
	// 100% that is 1.0 and the assertion is "the pixel width is the configured
	// width"; at 150% it is "the same configured ratio", which is the claim being
	// made either way.
	scale := float64(got1) / float64(asked1)
	want2 := float64(asked2) * scale
	t.Logf("REAL WINDOW 255: asked %dx%d -> GetWindowRect %dx%d (hwnd 0x%X); after dispose+show asked %dx%d -> %dx%d (hwnd 0x%X); implied scale %.4f; geometry source read %d time(s) for 2 creates",
		asked1, askedH, got1, h1, hwnd1, asked2, askedH, got2, h2, hwnd2, scale, reads)

	if got2 <= got1 {
		t.Errorf("[panel] width went from %d to %d across a dispose and a fresh create and the real window did not follow: %d px then %d px (scale %.4f)",
			asked1, asked2, got1, got2, scale)
	}
	if d := float64(got2) - want2; d > 2 || d < -2 {
		t.Errorf("the recreated window is %d px wide, want %.0f (= %d at the scale %.4f the first window showed) +/- 2: the number reached the create call but the window is not sized off it",
			got2, want2, asked2, scale)
	}
	if reads != 2 {
		t.Errorf("the geometry source was consulted %d times for two creates, want 2 (one per create, and the recreate must consult it again)", reads)
	}
	if h1 != h2 {
		t.Errorf("height was asked the same (%d) across the recreate and the window answers %d then %d", askedH, h1, h2)
	}

	// Self-nail, and it is an assertion rather than a defer: a window this case
	// started and did not finish is exactly what poisons the next tenant of the
	// thread.
	hh.call(func() { mgr.Destroy() })
	if mgr.IsCreated() {
		t.Error("the final Destroy did not retire the manager's window")
	}
	left := settleThreadWindows255(t, hh, tid, 3*time.Second)
	if left != 0 {
		t.Errorf("TEARDOWN NAIL: thread %d still owns %d window(s) at the end of this case, and it is never unlocked, so the leak stays in this process: %s",
			tid, left, "see panel_host_windows_test.go:90 for why this thread is not returned to the pool")
	}
	var dirty bool
	var stale uintptr
	hh.call(func() { dirty, stale = staleCloseQueued() })
	// The library's Terminate posts WM_QUIT on WM_DESTROY (webview.go:242-243 +
	// 381-383), so a quit here is the expected shape of a thread that hosted a
	// window and is now dying; an undispatched WM_CLOSE with no window to close is
	// the shape bringUp refuses to create over. Only the second one is reported.
	if dirty {
		t.Errorf("TEARDOWN NAIL: the panel thread still owes an undispatched WM_CLOSE for hwnd 0x%X with %d window(s) of its own left; that is the poison shape ticket 255's dispatch item names",
			stale, left)
	}
	t.Logf("TEARDOWN NAIL: thread %d windows=%d, staleCloseQueued=%v (a latched WM_QUIT is expected and is why the harness never unlocks this thread)",
		tid, left, dirty)
}

// settleThreadWindows255 waits, on a bounded monotonic deadline, for the harness
// thread to dispatch the close a Destroy posted and for its window tally to fall to
// zero. It returns the tally as last seen, so a caller that wants to report a leak
// reports the number it measured rather than a guess about the settle time.
//
// Why a wait and not an immediate read: measured on this box, the tally taken on
// the same instruction as Destroy was 3 (the panel plus its WebView2 children),
// because Destroy only posts WM_CLOSE and the owning thread is the one that gets to
// dispatch it. An immediate assert there is a false red; an unbounded wait is a
// hang. Zero within the bound is the shape this case's own nail requires, and it is
// the precondition bringUp refuses to create over.
func settleThreadWindows255(t *testing.T, hh *hostThreadHarness, tid uint32, limit time.Duration) int {
	t.Helper()
	deadline := time.Now().Add(limit)
	var left int
	for {
		left = threadWindowCount(tid)
		if left == 0 || time.Now().After(deadline) {
			return left
		}
		hh.call(func() {}) // keep the pump servicing while the tally settles
		time.Sleep(50 * time.Millisecond)
	}
}
