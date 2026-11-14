//go:build windows && winlive

package main

// Ticket 255 AC#4, leg 255-r8 - the ON-SCREEN reading of the re-show hop, the one
// 255-r1's own header names as missing: "the SAME configured number produces
// DIFFERENT on-screen widths depending on which hop carried it, by a frame delta
// that is 仅本机可量 - only a winlive reading can name that many pixels, and no case
// here pretends to have one" (panel_reshow_255r1_windows_test.go).
//
// WHY THE EXISTING WINLIVE CASE CANNOT SAY THIS. panel_geometry_255_winlive_test.go
// measures the CREATE hop (create, Destroy, create again): its numbers come out of
// CreateWindowExW, which reads the pair as the OUTER FRAME. 255-r1's hop is a
// different one - Show's re-show branch posts SetSize(w,h,HintNone) to a window that
// was never destroyed, and panelControl.SetSize runs AdjustWindowRect first
// (cmd/wisp/panel_control_windows.go, the library's semantics at webview.go:405), so
// THAT pair is the CLIENT area. This file drives the re-show branch on a live window
// and reports both rects on both sides of the edit.
//
// WHAT IT MEASURES, END TO END - the production chain, not a stand-in: a real
// config.toml on disk; the assembly root's own closure (panelGeometrySource via
// withGeometrySource, the same two functions cmd/wisp/panel_resident_windows.go:273
// hands the shipping host); the production panel thread (startResidentPanel: bringUp
// then the library pump Run(), which is the ONLY drainer of the Dispatch queue); and
// the production request path (RequestToggle to hide, RequestShow to show again -
// what a hot key or a tray item posts). GetClientRect / GetWindowRect are read before
// and after, and the three facts the hop owes are asserted:
//
//  1. the same HWND is still there (no recreate; C27's hide-don't-destroy);
//  2. the CLIENT area lands on the number [panel] width carries after the edit
//     (SetSize's pair is a client rect);
//  3. the OUTER frame follows by the SAME frame delta the first sample showed - that
//     delta is MEASURED here, never assumed. Because the two hops have the two
//     different semantics above, the same config number is NOT the same on-screen
//     width across them; this file reports both numbers rather than claiming they
//     are equal.
//
// HONEST LIMITS, named. The display scale is taken from the first sample (outer ==
// asked on a 100% display) and the second assertion is stated through it; the DPI
// logical-to-physical layer itself does not exist in this repository (255-a2
// measured that, ledger A700 section 2) and this file claims nothing about it.
// Teardown rides the production stop() (WM_QUIT to the panel thread, then a bounded
// wait for it to be seen gone) and asserts isFinished() before returning.
//
// PATH warning (ticket 98): this package's test binary links sherpa-onnx and dies at
// load (0xc0000135) unless third_party/sherpa-onnx is on PATH.

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/panel"
	"golang.org/x/sys/windows"
)

// reshowPanelWait255r8 bounds the wait for the SetSize hop to land. A monotonic
// deadline (time.Now().Add + After), never a wall-clock difference - the same
// shape settleThreadWindows255 uses for the same reason.
const reshowPanelWait255r8 = 10 * time.Second

var (
	reshow255r8User32        = windows.NewLazySystemDLL("user32.dll")
	reshow255r8GetClientRect = reshow255r8User32.NewProc("GetClientRect")
)

func clientRect255r8(hwnd uintptr) (rect255, error) {
	var r rect255
	rc, _, callErr := reshow255r8GetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&r)))
	if rc == 0 {
		return r, fmt.Errorf("GetClientRect(0x%X) returned 0: %w", hwnd, callErr)
	}
	return r, nil
}

// panelReshowSeed255r8 is a minimal valid config.toml carrying exactly the numbers
// this reading needs. schema_version 2 is internal/config's current one, and the two
// [panel] keys are written explicitly so the reading leans on no default.
func panelReshowSeed255r8(width, height int) string {
	return fmt.Sprintf("schema_version = 2\n\n[panel]\nwidth = %d\nheight = %d\n", width, height)
}

func writePanelReshowSeed255r8(t *testing.T, dir string, width, height int) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(panelReshowSeed255r8(width, height)), 0o600); err != nil {
		t.Fatalf("plant config.toml: %v", err)
	}
}

// TestTicket255r8RealWindowFollowsTheConfigAcrossAReshow is AC#4's re-show reading
// on a real desktop session: the window is born off [panel] width, the config is
// rewritten on disk, the window is hidden and shown again without being destroyed,
// and the two rects before and after are compared.
func TestTicket255r8RealWindowFollowsTheConfigAcrossAReshow(t *testing.T) {
	const (
		asked1W, asked1H = 640, 400 // the schema's [panel] width default, and a height nothing auto-sizes
		asked2W          = 900      // no constant in this repository can produce this number
	)
	dir := t.TempDir()
	writePanelReshowSeed255r8(t, dir, asked1W, asked1H)

	// The production closure, wrapped only to count consultations. Everything
	// inside it (the fresh config.LoadFile, the loud unreadable-config branch) is
	// the shipping code the assembly root hands the host.
	var reads atomic.Int64
	base := panelGeometrySource(dir)
	src := func() (int, int) {
		reads.Add(1)
		return base()
	}

	disp := &panel.ComposerDispatch{Mode: &recordingModeHandler{}}
	// The WebView2 user-data folder deliberately lives under the SYSTEM temp dir,
	// not under t.TempDir(): TempDir's RemoveAll runs while Edge's profile files can
	// still be flushing, and that race reddens the cleanup rather than the reading
	// (measured on this leg: "EBWebView\\Default: The directory is not empty" after
	// the panel thread had already exited cleanly). startPanelForTest uses the same
	// shape for the same reason.
	mgr := NewPanelManager(disp, builtinAssetsOrTestNil(t),
		filepath.Join(os.TempDir(), "wisp-255r8-panel-profile"), withGeometrySource(src))
	rp := startResidentPanel(observe.NewRegistry(), mgr)
	t.Cleanup(rp.stop)

	// First show: the window is born off [panel] (the create hop, OUTER semantics).
	showAndWait(t, rp)
	hwnd1 := uintptr(mgr.windowHandle())
	if hwnd1 == 0 {
		t.Fatal("the first show created no window to measure - a failed measurement, not a pass")
	}
	outer1, err := windowRect255(hwnd1)
	if err != nil {
		t.Fatalf("the first window's rect: %v", err)
	}
	client1, err := clientRect255r8(hwnd1)
	if err != nil {
		t.Fatalf("the first window's client rect: %v", err)
	}
	outer1W := int(outer1.right - outer1.left)
	client1W := int(client1.right - client1.left)
	frameW1 := outer1W - client1W

	// The edit this ticket is about: a new [panel] width on disk. No restart, and no
	// destroy below - the window under measurement stays alive throughout.
	writePanelReshowSeed255r8(t, dir, asked2W, asked1H)

	// 关窗再开，不销毁: hide through the production toggle request, then the show
	// request that takes Show's re-show branch.
	if !rp.RequestToggle("255r8-hide") {
		t.Fatalf("the panel thread refused the hide request (startUp error: %v)", rp.startUpErr())
	}
	waitPanelTrue(t, "the panel to hide", func() bool { return !rp.mgr.IsShown() })
	if !rp.RequestShow("255r8-reshow") {
		t.Fatalf("the panel thread refused the re-show request (startUp error: %v)", rp.startUpErr())
	}
	waitPanelTrue(t, "the panel to show again", func() bool { return rp.mgr.IsShown() })

	// The SetSize hop arrives through the Dispatch queue, which the pump drains
	// after the show request itself; wait for the client area to land. Silence
	// within the bound is a failed measurement and is reported as one.
	scale := float64(outer1W) / float64(asked1W)
	wantClient2W := float64(asked2W) * scale
	deadline := time.Now().Add(reshowPanelWait255r8)
	for {
		if uintptr(mgr.windowHandle()) != hwnd1 {
			break // recreated; the assertions below name it with both handles
		}
		cur, cerr := clientRect255r8(hwnd1)
		if cerr == nil {
			d := float64(int(cur.right-cur.left)) - wantClient2W
			if d <= 2 && d >= -2 {
				break
			}
		}
		if time.Now().After(deadline) {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	hwnd2 := uintptr(mgr.windowHandle())
	outer2, err := windowRect255(hwnd2)
	if err != nil {
		t.Fatalf("the re-shown window's rect: %v", err)
	}
	client2, err := clientRect255r8(hwnd2)
	if err != nil {
		t.Fatalf("the re-shown window's client rect: %v", err)
	}
	outer2W := int(outer2.right - outer2.left)
	client2W := int(client2.right - client2.left)
	frameW2 := outer2W - client2W

	t.Logf("REAL WINDOW 255r8: config %dx%d -> create hop: outer %dx%d / client %dx%d (frame delta %d px) hwnd=0x%X; config rewritten to %dx%d, hide+show with no destroy -> outer %dx%d / client %dx%d (frame delta %d px) hwnd=0x%X (same hwnd: %t); implied scale %.4f; geometry source read %d time(s)",
		asked1W, asked1H, outer1W, int(outer1.bottom-outer1.top), client1W, int(client1.bottom-client1.top), frameW1, hwnd1,
		asked2W, asked1H, outer2W, int(outer2.bottom-outer2.top), client2W, int(client2.bottom-client2.top), frameW2, hwnd2, hwnd2 == hwnd1,
		scale, reads.Load())

	if hwnd2 != hwnd1 {
		t.Errorf("the re-show replaced the window (hwnd 0x%X -> 0x%X): 票 255 AC#4's 甲形 resizes the window that exists, and C27 keeps hide-don't-destroy", hwnd1, hwnd2)
	}
	if d := float64(client2W) - wantClient2W; d > 2 || d < -2 {
		t.Errorf("after the re-show the CLIENT width is %d px, want %.0f (= %d at the scale %.4f the first window showed) +/- 2: the [panel] number did not reach the live window",
			client2W, wantClient2W, asked2W, scale)
	}
	// Fact 3, stated the way the two semantics actually compose: the create reads its
	// pair as the OUTER frame (so the first CLIENT was asked1W minus that frame),
	// while the re-show reads its pair as the CLIENT area - so comparing
	// client2W-client1W against the bare configured delta would be the wrong ruler by
	// exactly frameW1 px. The right statement is the outer as client plus frame.
	if d := float64(outer2W) - (wantClient2W + float64(frameW1)); d > 2 || d < -2 {
		t.Errorf("the OUTER width is %d px, want %.0f (= the CLIENT width %.0f plus the %d px frame delta the first sample measured) +/- 2: the outer frame did not follow the client by the same frame",
			outer2W, wantClient2W+float64(frameW1), wantClient2W, frameW1)
	}
	if frameW2 != frameW1 {
		t.Errorf("the frame delta moved from %d to %d px across the resize: the outer read and the client read are no longer about the same window frame", frameW1, frameW2)
	}
	if outer2W <= outer1W {
		t.Errorf("the OUTER width did not grow across the edit: %d px then %d px - the client reading above and this one cannot both be about a real window", outer1W, outer2W)
	}
	if got := reads.Load(); got != 2 {
		t.Errorf("the geometry source was consulted %d time(s), want 2 (one per create, one per re-show): a different count means one of the two hops did not read the config this reading is about", got)
	}

	// Teardown, asserted rather than trusted: stop() posts WM_QUIT to the panel
	// thread (the only way the library pump returns) and waits, bounded, for the
	// thread to be seen gone.
	rp.stop()
	if !rp.isFinished() {
		t.Errorf("TEARDOWN NAIL: the panel thread is still running after stop() returned - its window may still be open")
	}
}
