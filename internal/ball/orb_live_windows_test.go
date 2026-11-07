//go:build windows && winlive

package ball

// Ticket-65 integration leg (A1594): the LIVE window test for the WebView2
// orb. It creates the real ball, waits for the real WebView2 runtime to attach
// the orb beside it, and requires the whole takeover to have happened; it also
// verifies the attach contract itself (one fixed edge, topmost, tool window,
// never-activate show) and tears everything down. Run explicitly:
//
//	go test -tags winlive ./internal/ball/ -run TestBallLiveOrb -v
//
// Excluded from the default suite like every live test, so `go test ./...`
// never pops windows. The fallback half of the contract is pinned in the
// default suite: a ball created under -tags winlive WITHOUT a WebView2 runtime
// fails THIS test loudly rather than silently rendering D2D - the fallback is
// a decision the harness reads (OrbDebugLine), never a guess.

import (
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// This file's own handle to user32's GetWindowLongPtrW: pIsWindow belongs to
// live_guard_windows_test.go (in-binary single declaration), and a negative
// GWL_* index cannot be a const uintptr, so it lives in a variable like the
// panel host's ctlGwlStyle does.
var orbTestGetWindowLongPtrW = modUser32.NewProc("GetWindowLongPtrW")

var gwlExStyleTest = -20 // GWL_EXSTYLE

// TestBallLiveOrbAttachStateAndTeardown: real ball + real orb window on the
// real desktop. Assertions ride the test's OWN hwnds (never FindWindow by
// class), the quiet-desktop precondition holds the desktop, and every handle
// this test creates is destroyed before it returns.
func TestBallLiveOrbAttachStateAndTeardown(t *testing.T) {
	store, err := OpenPositionStore(filepath.Join(t.TempDir(), "ballpos.json"))
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	requireQuietBallDesktop(t)
	baseHandles := handlesOfProcess()

	b, err := New(Options{
		SizePx:      BallSizeDefaultPx,
		Initial:     statemachine.StateSleeping,
		Store:       store,
		WindowTitle: liveBallTitle(t),
		Hotkeys:     liveHotkeys(),
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() {
		hwnd := b.DebugHWND()
		b.Close()
		if r, _, _ := pIsWindow.Call(uintptr(hwnd)); r != 0 {
			t.Fatalf("ball window survived Close: hwnd=%v", hwnd)
		}
	})

	// The orb attach is asynchronous from New's point of view only in the
	// WebView2 runtime cold start INSIDE attachOrb; the PostTask it rides
	// lands before any later task, and SetState posts after it - so by the
	// time the deadline below expires without a window, the runtime really
	// did not come up.
	deadline := time.Now().Add(10 * time.Second)
	var orbHWND windows.HWND
	for {
		orbHWND = b.DebugOrbHWND()
		if orbHWND != 0 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("orb window did not attach within 10s (WebView2 runtime missing or Embed refused); "+
				"debug face says %q", OrbDebugLine(false, statemachine.StateSleeping))
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Attach contract: the orb window is a real window, of the orb class's
	// process, sized as ONE fixed square, shown without activation, topmost.
	if r, _, _ := pIsWindow.Call(uintptr(orbHWND)); r == 0 {
		t.Fatalf("DebugOrbHWND is not a live window: %v", orbHWND)
	}
	exStyle, _, _ := orbTestGetWindowLongPtrW.Call(uintptr(orbHWND), uintptr(gwlExStyleTest))
	const wantEx = wsExLayered | wsExTopmost | wsExToolWindow | wsExNoActivate
	if exStyle&wantEx != wantEx {
		t.Errorf("orb ex-style %#x missing bits of %#x", exStyle, wantEx)
	}
	var wr rect
	if r, _, _ := pGetWindowRect.Call(uintptr(orbHWND), unsafePtr(&wr)); r == 0 {
		t.Fatal("GetWindowRect(orb) failed")
	}
	if w, h := wr.width(), wr.height(); w != h || w <= 0 {
		t.Errorf("orb window is not a square: %dx%d", w, h)
	}

	// The takeover contract, on the REAL surface: the debug face names the
	// webview2 orb, and a SetState through the ball reaches the orb's page.
	if got := OrbDebugLine(b.DebugOrbAttached(), b.DebugOrbState()); got != "orb-surface=webview2-orb orb-state=idle" {
		t.Errorf("debug face after attach = %q", got)
	}
	b.SetState(statemachine.StateThinking)
	deadline = time.Now().Add(5 * time.Second)
	for {
		if got := OrbDebugLine(b.DebugOrbAttached(), b.DebugOrbState()); got == "orb-surface=webview2-orb orb-state=thinking" {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("orb state push never surfaced; face stuck at %q", OrbDebugLine(b.DebugOrbAttached(), b.DebugOrbState()))
		}
		time.Sleep(50 * time.Millisecond)
	}

	// LEAK sensitivity: repeated pushes must be free. Handles ate ~300 on the
	// attach itself (measured twice on this box, and of the same shape the
	// vendored module imposes on the panel host - it never releases the
	// environment). That one-time-per-process cost is bounded by construction
	// (attachOrb is CAS-guarded to ONE orb per process), so the number that
	// can still grow and hide a leak is the per-push delta: drive the same
	// burst the audio capture would (50 state pushes + 50 audio pushes) and
	// require the handle count to come back.
	afterAttach := handlesOfProcess()
	for i := 0; i < 50; i++ {
		b.SetState(statemachine.StateListening)
		b.SetAudioLevel(0.5)
		b.SetState(statemachine.StateThinking)
	}
	deadline = time.Now().Add(5 * time.Second)
	for {
		// The last push lands on the STA thread; give the queue a moment to
		// drain before reading the counter.
		time.Sleep(100 * time.Millisecond)
		if got := handlesOfProcess(); got <= afterAttach+50 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("handle growth across 100 orb pushes: before=%d after=%d",
				afterAttach, handlesOfProcess())
		}
	}

	// The ball window itself is still alive and still the ball: the orb is a
	// SECOND window, not a replacement.
	if r, _, _ := pIsWindow.Call(uintptr(b.DebugHWND())); r == 0 {
		t.Fatal("ball window died while the orb was attached")
	}

	// Teardown: Close must destroy BOTH windows and stay inside the one-time
	// attach budget (documented above; the window-death checks are the hard
	// assertions, this one only catches unbounded growth).
	b.Close()
	if r, _, _ := pIsWindow.Call(uintptr(orbHWND)); r != 0 {
		t.Fatalf("orb window survived Close: hwnd=%v", orbHWND)
	}
	if got := handlesOfProcess(); got > baseHandles+600 {
		t.Errorf("handle growth beyond the one-time attach budget: base=%d end=%d", baseHandles, got)
	}
}

// TestBallLiveOrbFollowPlacementDecision pins the pure side-choice math behind
// follow() against real work-area shapes, including the off-screen shapes the
// orchestrator measured (ticket 65): a ball docked flush at the right work
// edge must flip the orb to the LEFT, and a ball too centered for either side
// must clamp the orb fully inside the work area.
func TestBallLiveOrbFollowPlacementDecision(t *testing.T) {
	const edge, gap = 132, 12
	work := Rect{L: 0, T: 0, R: 3440, B: 1320}
	cases := []struct {
		name   string
		ball   rect
		work   Rect
		workOk bool
		wantX  int32
		wantY  int32
	}{
		{"room on the right keeps the legacy side", rect{l: 3000, t: 100, r: 3100, b: 200}, work, true, 3112, 100},
		{"docked right edge flips to the left", rect{l: 3307, t: 100, r: 3440, b: 200}, work, true, 3307 - gap - edge, 100},
		{"docked bottom clamps the vertical axis", rect{l: 3307, t: 1220, r: 3440, b: 1320}, work, true, 3307 - gap - edge, 1320 - edge},
		{"ball above the work top clamps down", rect{l: 1000, t: -50, r: 1100, b: 50}, work, true, 1112, 0},
		{"narrow screen clamps instead of flipping", rect{l: 60, t: 100, r: 200, b: 200}, Rect{L: 0, T: 0, R: 260, B: 1320}, true, 260 - edge, 100},
		{"no monitor info keeps the legacy shape", rect{l: 3000, t: 100, r: 3100, b: 200}, work, false, 3112, 100},
	}
	for _, tc := range cases {
		gotX, gotY := orbFollowXY(tc.ball, edge, gap, tc.work, tc.workOk)
		if gotX != tc.wantX || gotY != tc.wantY {
			t.Errorf("%s: orbFollowXY() = (%d,%d), want (%d,%d)", tc.name, gotX, gotY, tc.wantX, tc.wantY)
		}
	}
}

// TestBallLiveOrbFollowStaysInsideWorkArea is the off-screen regression
// (ticket 65, orchestrator-measured): with the ball parked flush against the
// bottom-right corner of its monitor's WORK area, follow() must land the orb
// FULLY inside that work area at real coordinates. The pre-fix behavior
// (unconditional wr.r+12) starts the orb 12 px outside the work area here, so
// this test is red on the old code and green on the fix - the assertions read
// the actual window rects, never "follow was called".
func TestBallLiveOrbFollowStaysInsideWorkArea(t *testing.T) {
	store, err := OpenPositionStore(filepath.Join(t.TempDir(), "ballpos.json"))
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	requireQuietBallDesktop(t)

	b, err := New(Options{
		SizePx:      BallSizeDefaultPx,
		Initial:     statemachine.StateSleeping,
		Store:       store,
		WindowTitle: liveBallTitle(t),
		Hotkeys:     liveHotkeys(),
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(b.Close)

	deadline := time.Now().Add(10 * time.Second)
	for b.DebugOrbHWND() == 0 {
		if time.Now().After(deadline) {
			t.Fatalf("orb window did not attach within 10s (WebView2 runtime missing or Embed refused)")
		}
		time.Sleep(100 * time.Millisecond)
	}

	ballHWND := uintptr(b.DebugHWND())
	mon, ok := monitorRectsForWindow(windows.HWND(ballHWND))
	if !ok {
		t.Fatal("monitorRectsForWindow(ball) returned no monitor")
	}
	work := mon.Work

	// Park the ball flush against the bottom-right corner of the work area:
	// ball.r == work.R and ball.b == work.B, so BOTH legacy axes point
	// outside and the fix has to flip the side AND clamp vertically.
	var br rect
	if r, _, _ := pGetWindowRect.Call(ballHWND, unsafePtr(&br)); r == 0 {
		t.Fatal("GetWindowRect(ball) failed")
	}
	pSetWindowPos.Call(ballHWND, 0,
		uintptr(int32(work.R)-br.width()), uintptr(int32(work.B)-br.height()), 0, 0,
		uintptr(swpNoSize|swpNoZOrder|swpNoActivate))
	// Re-read where the ball really settled (the wndproc arms inside the
	// SetWindowPos above run synchronously on the STA thread), so the exact
	// assertions below measure against the settled rect.
	if r, _, _ := pGetWindowRect.Call(ballHWND, unsafePtr(&br)); r == 0 {
		t.Fatal("GetWindowRect(ball) after the move failed")
	}

	b.followOrb()

	var orb rect
	if r, _, _ := pGetWindowRect.Call(uintptr(b.DebugOrbHWND()), unsafePtr(&orb)); r == 0 {
		t.Fatal("GetWindowRect(orb) failed")
	}
	t.Logf("work=%+v ball=(%d,%d)-(%d,%d) orb=(%d,%d)-(%d,%d)",
		work, br.l, br.t, br.r, br.b, orb.l, orb.t, orb.r, orb.b)

	if orb.l < int32(work.L) || orb.r > int32(work.R) {
		t.Fatalf("orb horizontally outside the work area: work.R=%d orb=(%d,%d)-(%d,%d)",
			work.R, orb.l, orb.t, orb.r, orb.b)
	}
	if orb.t < int32(work.T) || orb.b > int32(work.B) {
		t.Fatalf("orb vertically outside the work area: work.B=%d orb=(%d,%d)-(%d,%d)",
			work.B, orb.l, orb.t, orb.r, orb.b)
	}
	// Exact coordinates, not just containment: a flush-right ball with room on
	// its left must flip the orb to the LEFT (orb.r == ball.l - gap) and clamp
	// its bottom to the work bottom.
	edge := orb.width()
	if left := br.l - orbFollowGap - edge; left >= int32(work.L) {
		if orb.r != br.l-orbFollowGap {
			t.Errorf("orb did not flip to the left of the docked-right ball: want orb.r=%d, got orb=(%d,%d)-(%d,%d)",
				br.l-orbFollowGap, orb.l, orb.t, orb.r, orb.b)
		}
	}
	if br.b >= int32(work.B) && orb.b != int32(work.B) {
		t.Errorf("orb bottom was not clamped to the work bottom: want orb.b=%d, got %d",
			int32(work.B), orb.b)
	}
}

// orbTestGetPixel / orbTestGetStockObject: this file's own handles to the two
// GDI entries the pixel regression below needs (screen-DC sampling and the
// black class brush).
var (
	orbTestGetPixel       = modGDI32.NewProc("GetPixel")
	orbTestGetStockObject = modGDI32.NewProc("GetStockObject")
)

// GetStockObject(BLACK_BRUSH).
const orbTestBlackBrush = 4

// orbTestBackdropWndProc: the backdrop window never handles anything itself.
func orbTestBackdropWndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	r, _, _ := pDefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return r
}

var (
	orbTestBackdropProcPtr = windows.NewCallback(orbTestBackdropWndProc)
	orbTestBackdropOnce    sync.Once
	orbTestBackdropErr     error
)

// orbTestCreateBackdrop registers (once) and shows a plain black WS_POPUP
// window covering the orb's rect, so the screen pixels the regression samples
// have a KNOWN floor behind the orb: through the orb's transparent pixels the
// sample reads black, not whatever wallpaper happens to sit on this desktop.
// It is created TOPMOST (below the orb, which the test re-asserts afterwards)
// so ordinary application windows can never composite above the floor and
// contaminate the reading. It MUST be called on the ui-sta thread (via
// Ball.sta.PostTask) - a window on a thread nobody pumps never paints, and
// the whole point is painted black pixels.
func orbTestCreateBackdrop(wr rect) uintptr {
	orbTestBackdropOnce.Do(func() {
		brush, _, _ := orbTestGetStockObject.Call(orbTestBlackBrush)
		wc := wndClassEx{
			size:      uint32(unsafe.Sizeof(wndClassEx{})),
			wndProc:   orbTestBackdropProcPtr,
			instance:  windows.Handle(moduleHandle()),
			cursor:    windows.Handle(loadArrowCursor()),
			bkgnd:     windows.Handle(brush),
			className: utf16("WispOrbPxBackdrop"),
		}
		if r, _, _ := pRegisterClassExW.Call(unsafePtr(&wc)); r == 0 {
			orbTestBackdropErr = fmt.Errorf("RegisterClassExW(backdrop) failed")
		}
	})
	if orbTestBackdropErr != nil {
		return 0
	}
	hwnd, _, _ := pCreateWindowExW.Call(0x8, // WS_EX_TOPMOST: above every normal window, below the re-asserted orb
		unsafePtr(utf16("WispOrbPxBackdrop")), unsafePtr(utf16("Wisp Orb Px Backdrop")),
		uintptr(wsPopup),
		uintptr(wr.l), uintptr(wr.t), uintptr(wr.width()), uintptr(wr.height()),
		0, 0, uintptr(moduleHandle()), 0)
	if hwnd == 0 {
		return 0
	}
	pShowWindow.Call(hwnd, swShownoactivate)
	return hwnd
}

// orbTestNonDarkRatio samples the SCREEN DC (the composited desktop, i.e.
// exactly what a screenshot of the orb region shows) across the orb rect and
// returns (fraction of sampled pixels with any channel > 30, sample count).
// It is the pixel-truth reading the regression's teeth are made of.
func orbTestNonDarkRatio(l, t, w, h int32) (float64, int) {
	hdc, _, _ := pGetDC.Call(0)
	if hdc == 0 {
		return 0, 0
	}
	defer pReleaseDC.Call(0, hdc)
	const step = 3 // px between samples
	nonDark, total := 0, 0
	for y := t; y < t+h; y += step {
		for x := l; x < l+w; x += step {
			c, _, _ := orbTestGetPixel.Call(hdc, uintptr(x), uintptr(y))
			total++
			if c&0xff > 30 || (c>>8)&0xff > 30 || (c>>16)&0xff > 30 {
				nonDark++
			}
		}
	}
	if total == 0 {
		return 0, 0
	}
	return float64(nonDark) / float64(total), total
}

// TestBallLiveOrbPageRendersPixels is the black-orb regression (ticket 65,
// orchestrator-measured 2026-11-07): the first serving seam answered the
// entry document through a CreateWebResourceResponse the jchv wrapper
// misreported as failed, the page NEVER loaded, and the window stayed pure
// black. This test puts the orb on-screen over a KNOWN black backdrop and
// requires the composited screen pixels to actually light up - at least 1% of
// the sampled orb-area pixels carrying a channel above 30. It fails loudly on
// the old shape (all-dark window) and on any future regression that breaks
// page delivery or WebGPU boot: "no error logged" is not a passing shape
// here, painted pixels are.
func TestBallLiveOrbPageRendersPixels(t *testing.T) {
	store, err := OpenPositionStore(filepath.Join(t.TempDir(), "ballpos.json"))
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	requireQuietBallDesktop(t)

	b, err := New(Options{
		SizePx:      BallSizeDefaultPx,
		Initial:     statemachine.StateSleeping,
		Store:       store,
		WindowTitle: liveBallTitle(t),
		Hotkeys:     liveHotkeys(),
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(b.Close)

	deadline := time.Now().Add(15 * time.Second)
	var orbHWND windows.HWND
	for {
		orbHWND = b.DebugOrbHWND()
		if orbHWND != 0 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("orb window did not attach within 15s (WebView2 runtime missing or Embed refused)")
		}
		time.Sleep(100 * time.Millisecond)
	}

	mon, ok := monitorRectsForWindow(windows.HWND(b.DebugHWND()))
	if !ok {
		t.Fatal("monitorRectsForWindow(ball) returned no monitor")
	}
	work := mon.Work

	// Park the ball at the work-area corner so nothing of ours overlaps the
	// sampled region, and center the orb in the work area so every sampled
	// pixel is on-screen.
	ballHWND := uintptr(b.DebugHWND())
	pSetWindowPos.Call(ballHWND, 0,
		uintptr(int32(work.L)+8), uintptr(int32(work.T)+8), 0, 0,
		uintptr(swpNoSize|swpNoZOrder|swpNoActivate))
	var wr rect
	if r, _, _ := pGetWindowRect.Call(uintptr(orbHWND), unsafePtr(&wr)); r == 0 {
		t.Fatal("GetWindowRect(orb) failed")
	}
	edge := wr.width()
	ox := (int32(work.L) + int32(work.R)) / 2
	oy := (int32(work.T) + int32(work.B)) / 2
	pSetWindowPos.Call(uintptr(orbHWND), orbHwndTopmost,
		uintptr(ox-edge/2), uintptr(oy-edge/2), 0, 0,
		uintptr(swpNoSize|swpNoActivate))
	if r, _, _ := pGetWindowRect.Call(uintptr(orbHWND), unsafePtr(&wr)); r == 0 {
		t.Fatal("GetWindowRect(orb) after the move failed")
	}

	// Known-black floor behind the orb, created and painted on the ui-sta
	// thread (its pump is what makes the window actually paint). The floor
	// went in topmost, so the orb must be re-asserted ABOVE it in the
	// topmost band before any sampling.
	backdrop := make(chan uintptr, 1)
	b.sta.PostTask(func() { backdrop <- orbTestCreateBackdrop(wr) })
	bh := <-backdrop
	if bh == 0 {
		t.Fatal("black backdrop window failed to create")
	}
	pSetWindowPos.Call(uintptr(orbHWND), orbHwndTopmost, 0, 0, 0, 0,
		uintptr(swpNoSize|swpNoMove|swpNoActivate))
	t.Cleanup(func() {
		done := make(chan struct{})
		b.sta.PostTask(func() {
			pDestroyWindow.Call(bh)
			close(done)
		})
		<-done
	})

	// WebGPU boot is not instant; poll until the orb lights up or the
	// deadline proves it never will.
	const budget = 25 * time.Second
	deadline = time.Now().Add(budget)
	best := 0.0
	var samples int
	for {
		ratio, n := orbTestNonDarkRatio(wr.l+4, wr.t+4, wr.width()-8, wr.height()-8)
		samples = n
		if ratio > best {
			best = ratio
		}
		if ratio > 0.01 {
			t.Logf("orb page rendered: %d/%d sampled px non-dark (%.2f%% > 1%% floor)",
				int(ratio*float64(n)), n, ratio*100)
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("orb window content stayed black after %v: best non-dark ratio %.3f%% over %d sampled px (page did not render; delivery or WebGPU broken)",
				budget, best*100, samples)
		}
		time.Sleep(300 * time.Millisecond)
	}
}
