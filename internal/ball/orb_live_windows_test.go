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
	// WS_EX_LAYERED is asserted ABSENT (ticket 65 transparency fix): the orb
	// never presents a ULW bitmap, so the layered bit only kept an
	// uninitialized black redirection surface under the WebView2 child - the
	// measured "black square" defect. The per-pixel alpha path is
	// DwmEnableBlurBehindWindow's empty region instead (Tauri/tao's recipe).
	if r, _, _ := pIsWindow.Call(uintptr(orbHWND)); r == 0 {
		t.Fatalf("DebugOrbHWND is not a live window: %v", orbHWND)
	}
	exStyle, _, _ := orbTestGetWindowLongPtrW.Call(uintptr(orbHWND), uintptr(gwlExStyleTest))
	const wantEx = wsExTransparent | wsExTopmost | wsExToolWindow | wsExNoActivate
	if exStyle&wantEx != wantEx {
		t.Errorf("orb ex-style %#x missing bits of %#x", exStyle, wantEx)
	}
	if exStyle&wsExLayered != 0 {
		t.Errorf("orb ex-style %#x carries WS_EX_LAYERED: the bit is the black-square defect, not transparency", exStyle)
	}
	// Overlay contract (ticket 65 rework): the orb is the ball's visible face,
	// so its rect must cover the ball window's rect (equal WindowEdgePx edges
	// make the unclamped alignment the identical rect).
	ballHWND := uintptr(b.DebugHWND())
	var bwr, owr rect
	if r, _, _ := pGetWindowRect.Call(ballHWND, unsafePtr(&bwr)); r == 0 {
		t.Fatal("GetWindowRect(ball) failed")
	}
	if r, _, _ := pGetWindowRect.Call(uintptr(orbHWND), unsafePtr(&owr)); r == 0 {
		t.Fatal("GetWindowRect(orb) failed")
	}
	if bwr.width() != owr.width() || bwr.height() != owr.height() {
		t.Errorf("orb window edge %dx%d differs from the ball window edge %dx%d: the overlay is mis-sized",
			owr.width(), owr.height(), bwr.width(), bwr.height())
	}
	work, workOk := monitorRectsForWindow(windows.HWND(ballHWND))
	wantX, wantY := orbOverlayXY(bwr, owr.width(), work.Work, workOk)
	if owr.l != wantX || owr.t != wantY {
		t.Errorf("orb rect (%d,%d) is not the overlay position follow() computes (%d,%d)",
			owr.l, owr.t, wantX, wantY)
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

	// The ball window itself is still alive and still the interaction layer:
	// the orb window is its visible FACE (overlaying it), not a replacement
	// for it - the whole frozen hit-test/drag/dock/click set lives here, and
	// the D2D fallback reclaims the picture if the orb ever fails to attach.
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

// TestBallLiveOrbFollowPlacementDecision pins the pure placement math behind
// follow() against real work-area shapes (ticket 65 rework): the orb overlays
// the ball (center-aligned), clamped fully inside the work area when monitor
// info exists - including the flush-docked shapes where the raw alignment
// would push the orb's glow past the screen edge - and raw when it does not.
func TestBallLiveOrbFollowPlacementDecision(t *testing.T) {
	const edge = 132
	work := Rect{L: 0, T: 0, R: 3440, B: 1320}
	cases := []struct {
		name   string
		ball   rect
		work   Rect
		workOk bool
		wantX  int32
		wantY  int32
	}{
		// 100x100 ball at (3000,100): center (3050,150), orb centered on it.
		{"free ball centers the orb on it", rect{l: 3000, t: 100, r: 3100, b: 200}, work, true, 3050 - edge/2, 150 - edge/2},
		// Ball flush at the right work edge: raw alignment x=3390-66=3324
		// overhangs (3324+132=3456 > 3440), so the clamp pulls it back.
		{"docked right edge clamps into the work area", rect{l: 3340, t: 100, r: 3440, b: 200}, work, true, 3440 - edge, 150 - edge/2},
		// Ball flush at the bottom: raw y=1290-66=1224 overhangs (1224+132=1356
		// > 1320), so the clamp pulls it back; x needs no clamp.
		{"docked bottom clamps the vertical axis", rect{l: 1000, t: 1260, r: 1100, b: 1320}, work, true, 1050 - edge/2, 1320 - edge},
		// Ball half above the work top: raw y=(0)-66=-66 < 0, clamped to 0.
		{"ball above the work top clamps down", rect{l: 1000, t: -50, r: 1100, b: 50}, work, true, 1050 - edge/2, 0},
		// Narrow work area narrower than the orb: both x-clamp rules hit; the
		// left-edge rule wins last.
		{"narrow screen clamps instead of overhanging", rect{l: 0, t: 100, r: 200, b: 200}, Rect{L: 0, T: 0, R: 100, B: 1320}, true, 0, 150 - edge/2},
		// No monitor info: raw center alignment, never a guessed boundary.
		{"no monitor info keeps the raw alignment", rect{l: 3340, t: 100, r: 3440, b: 200}, work, false, 3390 - edge/2, 150 - edge/2},
	}
	for _, tc := range cases {
		gotX, gotY := orbOverlayXY(tc.ball, edge, tc.work, tc.workOk)
		if gotX != tc.wantX || gotY != tc.wantY {
			t.Errorf("%s: orbOverlayXY() = (%d,%d), want (%d,%d)", tc.name, gotX, gotY, tc.wantX, tc.wantY)
		}
	}
}

// TestBallLiveOrbFollowStaysInsideWorkArea is the off-screen regression
// (ticket 65, orchestrator-measured), re-pinned for the rework: with the ball
// parked flush against the bottom-right corner of its monitor's WORK area,
// follow() must land the orb EXACTLY over the ball (equal edges make the
// center alignment the identical rect - both rects end flush at the work
// edges) and never outside the work area. The pre-rework side-by-side shape
// and any raw overhang both fail here: the assertions read the actual window
// rects, never "follow was called".
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
	// Exact overlay coordinates, not just containment: the orb and the ball
	// share one WindowEdgePx edge, so the flush-docked ball's rect IS the
	// aligned overlay position - no gap, no overhang, no side flip.
	if orb.l != br.l || orb.t != br.t || orb.r != br.r || orb.b != br.b {
		t.Errorf("orb did not overlay the ball exactly: ball=(%d,%d)-(%d,%d) orb=(%d,%d)-(%d,%d)",
			br.l, br.t, br.r, br.b, orb.l, orb.t, orb.r, orb.b)
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

// orbTestReapLeakedWebView2 terminates every msedgewebview2.exe process whose
// PARENT is this test process and reports how many it killed. Ticket 65
// (measured 2026-11-07): the jchv wrapper's ICoreWebView2Controller has no
// Close method, so production destroy() leaks each orb's browser processes;
// across a count=10 suite those leaks degrade rendering round over round
// (40.6% -> 27.2% -> 7.7% -> three all-black attempts) until the page cannot
// paint. Reaping OUR OWN children - matched by image name AND parent PID, so
// no other application's WebView2 can ever match - removes that accumulated
// pressure without touching any other process on the desktop.
func orbTestReapLeakedWebView2(t *testing.T) {
	t.Helper()
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		t.Logf("orb flake-fix: process snapshot failed (leak reaping skipped): %v", err)
		return
	}
	defer windows.CloseHandle(snap)
	mine := windows.GetCurrentProcessId()
	killed := 0
	var pe windows.ProcessEntry32
	pe.Size = uint32(unsafe.Sizeof(pe))
	err = windows.Process32First(snap, &pe)
	for err == nil {
		if pe.ParentProcessID == mine &&
			windows.UTF16ToString(pe.ExeFile[:]) == "msedgewebview2.exe" {
			if h, oerr := windows.OpenProcess(windows.PROCESS_TERMINATE, false, pe.ProcessID); oerr == nil {
				if terr := windows.TerminateProcess(h, 1); terr == nil {
					killed++
				}
				windows.CloseHandle(h)
			}
		}
		err = windows.Process32Next(snap, &pe)
	}
	if killed > 0 {
		t.Logf("orb flake-fix: reaped %d leaked WebView2 browser process(es) parented by this test process", killed)
	}
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
//
// Flake treatment (ticket 65, 2026-11-07): the sampling already waits on the
// real signal - painted screen pixels, polled - so the wait shape was never
// the problem. The measured red mechanism is ACCUMULATION: every WebView2
// instance this process ever embedded leaks its browser processes (the jchv
// wrapper's ICoreWebView2Controller has no Close method, so production
// destroy() cannot shut the runtime down gracefully), and the full-suite
// count=10 run showed the render degrading round over round - non-dark ratio
// 40.6% x7 rounds, then 27.2%, 7.7%, then three consecutive all-black
// attempts - i.e. leaked browser processes pile up GPU pressure until a new
// instance cannot paint at all. Two bounded countermeasures, both test-side,
// both zero assertion changes: reap the leaked msedgewebview2.exe children of
// this process (at test start and before every rebuild - they are OUR
// orphans, parented by us, so no other app's WebView2 can ever match), and
// keep the fresh-instance rebuild so a single wedged boot still gets a
// second chance.
func TestBallLiveOrbPageRendersPixels(t *testing.T) {
	store, err := OpenPositionStore(filepath.Join(t.TempDir(), "ballpos.json"))
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	requireQuietBallDesktop(t)
	orbTestReapLeakedWebView2(t)

	const attempts = 3
	var lastBest float64
	var lastSamples int
	for attempt := 1; attempt <= attempts; attempt++ {
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
		lit, best, samples := orbPixelAttempt(t, b, attempt)
		if lit {
			return
		}
		lastBest, lastSamples = best, samples
		// A dark boot keeps its WebView2 instance; release it and reap every
		// leaked browser process this test process has accumulated (the
		// measured red mechanism) before the rebuild is genuinely fresh
		// (Close is idempotent, the t.Cleanup below stays harmless).
		b.Close()
		orbTestReapLeakedWebView2(t)
		if attempt < attempts {
			t.Logf("attempt %d/%d stayed dark (best %.3f%% over %d sampled px); rebuilding ball and WebView2 instance",
				attempt, attempts, best*100, samples)
		}
	}
	t.Fatalf("orb window content stayed black after %d attempts: best non-dark ratio %.3f%% over %d sampled px (page did not render; delivery or WebGPU broken)",
		attempts, lastBest*100, lastSamples)
}

// orbPixelAttempt runs one create-and-sample round against the LIVE desktop:
// it waits for the orb window to attach, parks it centered over a known-black
// backdrop, and polls the composited screen pixels until they light up or the
// attempt budget expires. It returns (lit, best non-dark ratio, sample count)
// and reports true the moment the page has painted; the caller owns the
// retry decision. Fatal paths inside are hard regressions (attach loss,
// window-API failure), not the cold-boot flake - only the dark-expiry return
// is retryable.
func orbPixelAttempt(t *testing.T, b *Ball, attempt int) (bool, float64, int) {
	t.Cleanup(b.Close)

	deadline := time.Now().Add(15 * time.Second)
	var orbHWND windows.HWND
	for {
		orbHWND = b.DebugOrbHWND()
		if orbHWND != 0 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("attempt %d: orb window did not attach within 15s (WebView2 runtime missing or Embed refused)", attempt)
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
	// topmost band before any sampling. It is destroyed here - before the
	// caller can Close the ball - because its pump is this ball's STA thread.
	backdrop := make(chan uintptr, 1)
	b.sta.PostTask(func() { backdrop <- orbTestCreateBackdrop(wr) })
	bh := <-backdrop
	if bh == 0 {
		t.Fatal("black backdrop window failed to create")
	}
	pSetWindowPos.Call(uintptr(orbHWND), orbHwndTopmost, 0, 0, 0, 0,
		uintptr(swpNoSize|swpNoMove|swpNoActivate))
	defer func() {
		done := make(chan struct{})
		b.sta.PostTask(func() {
			pDestroyWindow.Call(bh)
			close(done)
		})
		<-done
	}()

	// WebGPU boot is not instant; poll until the orb lights up or the
	// deadline proves this attempt never will.
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
			t.Logf("attempt %d: orb page rendered: %d/%d sampled px non-dark (%.2f%% > 1%% floor)",
				attempt, int(ratio*float64(n)), n, ratio*100)
			return true, best, samples
		}
		if time.Now().After(deadline) {
			return false, best, samples
		}
		time.Sleep(300 * time.Millisecond)
	}
}

// orbTestCreateSolidBrush: this file's handle to the GDI entry the colored
// transparency floor below needs (its KNOWN floor is not black).
var orbTestCreateSolidBrush = modGDI32.NewProc("CreateSolidBrush")

// orbTestFloorColor is the known floor of the transparency regression, as a
// COLORREF (0x00BBGGRR): the desktop gray the orchestrator measured (51,51,50).
const orbTestFloorColor = 0x00323333

// orbTestColoredFloorWndProc: the colored floor never handles anything itself.
func orbTestColoredFloorWndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	r, _, _ := pDefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return r
}

var (
	orbTestColoredFloorProcPtr = windows.NewCallback(orbTestColoredFloorWndProc)
	orbTestColoredFloorOnce    sync.Once
	orbTestColoredFloorErr     error
)

// orbTestCreateColoredFloor registers (once) and shows a (51,51,50) WS_POPUP
// window covering the orb's rect - the transparency regression's KNOWN floor.
// A BLACK floor cannot serve: transparent-over-black and the black-square
// defect both read (0,0,0), so the reading could not tell the fix from the
// defect. Same creation law as orbTestCreateBackdrop: topmost (the orb is
// re-asserted above it), created and painted on the ui-sta thread via
// Ball.sta.PostTask - a window on a thread nobody pumps never paints.
func orbTestCreateColoredFloor(wr rect) uintptr {
	orbTestColoredFloorOnce.Do(func() {
		brush, _, _ := orbTestCreateSolidBrush.Call(orbTestFloorColor)
		wc := wndClassEx{
			size:      uint32(unsafe.Sizeof(wndClassEx{})),
			wndProc:   orbTestColoredFloorProcPtr,
			instance:  windows.Handle(moduleHandle()),
			cursor:    windows.Handle(loadArrowCursor()),
			bkgnd:     windows.Handle(brush),
			className: utf16("WispOrbPxClearFloor"),
		}
		if r, _, _ := pRegisterClassExW.Call(unsafePtr(&wc)); r == 0 {
			orbTestColoredFloorErr = fmt.Errorf("RegisterClassExW(clear-floor) failed")
		}
	})
	if orbTestColoredFloorErr != nil {
		return 0
	}
	hwnd, _, _ := pCreateWindowExW.Call(0x8, // WS_EX_TOPMOST, below the re-asserted orb
		unsafePtr(utf16("WispOrbPxClearFloor")), unsafePtr(utf16("Wisp Orb Clear Floor")),
		uintptr(wsPopup),
		uintptr(wr.l), uintptr(wr.t), uintptr(wr.width()), uintptr(wr.height()),
		0, 0, uintptr(moduleHandle()), 0)
	if hwnd == 0 {
		return 0
	}
	pShowWindow.Call(hwnd, swShownoactivate)
	return hwnd
}

// orbTestPixelAt samples one composited screen pixel, as (r,g,b).
func orbTestPixelAt(hdc uintptr, x, y int32) (int32, int32, int32) {
	c, _, _ := orbTestGetPixel.Call(hdc, uintptr(x), uintptr(y))
	return int32(c & 0xff), int32((c >> 8) & 0xff), int32((c >> 16) & 0xff)
}

// The probe constants the corners regression samples with (wave68 moved them
// to file scope so the bounded-retry wrapper and its final Fatal share them).
const (
	orbCornerInset  = int32(3)  // the orchestrator's probe: 3 px inside the rect
	orbCornerTol    = int32(12) // per-channel distance from the floor color
	orbPageDiffGate = 25        // a page pixel differs from the floor by more than this
)

// TestBallLiveOrbCornersReadTheDesktop is the transparency regression (ticket
// 65 rework, orchestrator-measured black corners on the first cut: left-up
// (30,24,29) / right-up (38,33,40) / left-down (26,21,32) / right-down
// (38,31,38) against a (51,51,50) desktop). The orb is parked over a KNOWN
// (51,51,50) floor, and the four corners of the orb window - sampled 3 px
// inside, exactly the orchestrator's probe shape - must read the floor within
// 12 per channel. The page must ALSO actually be alive (some sampled pixel
// differing from the floor by more than 25), so a dead-page window cannot
// pass by being invisible. Both readings poll to a deadline: WebGPU boot is
// not instant.
//
// Flake treatment (wave68, the file's established shape): a dark WebView2
// boot fails the life half with page lit ratio 0.000% while the corners still
// read the floor - the page simply never painted inside the budget (measured
// on this box 2026-11-07, the A1602 dark-boot family). orbPixelAttempt's
// bounded rebuild is the in-file cure, test-side, zero assertion changes: up
// to three fresh ball+WebView2 instances before the hard fail below.
func TestBallLiveOrbCornersReadTheDesktop(t *testing.T) {
	store, err := OpenPositionStore(filepath.Join(t.TempDir(), "ballpos.json"))
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	requireQuietBallDesktop(t)
	orbTestReapLeakedWebView2(t)

	const attempts = 3
	var lastReadings [4][3]int32
	lastLit := 0.0
	for attempt := 1; attempt <= attempts; attempt++ {
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
		ok, readings, lit := orbCornerAttempt(t, b)
		if ok {
			return
		}
		lastReadings, lastLit = readings, lit
		// A dark boot keeps its WebView2 instance; release it and reap every
		// leaked browser process (the measured accumulation mechanism) so the
		// next attempt starts genuinely fresh.
		b.Close()
		orbTestReapLeakedWebView2(t)
		if attempt < attempts {
			t.Logf("attempt %d/%d did not clear both halves (page lit ratio %.3f%%); rebuilding ball and WebView2 instance",
				attempt, attempts, lastLit*100)
		}
	}
	t.Fatalf("orb transparency regression after %d attempts: corners read %v (want each channel within +/-%d of (51,51,50)); page lit ratio %.3f%% (want > 1%%)",
		attempts, lastReadings, orbCornerTol, lastLit*100)
}

// orbCornerAttempt runs one create-and-sample round of the transparency
// regression against the LIVE desktop: it waits for the orb window to attach,
// centers the overlay over a known (51,51,50) floor, and polls until BOTH
// halves clear - the four corners reading the floor within tolerance AND the
// page compositing somewhere - or the attempt budget expires. Hard failures
// (attach loss, window-API breakage) are Fatal regressions, not retryable;
// only the dark-expiry return drives the caller's rebuild. It returns
// (cleared, last corner readings, last lit ratio).
func orbCornerAttempt(t *testing.T, b *Ball) (bool, [4][3]int32, float64) {
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

	// Center the overlay in the work area: the ball window is the anchor the
	// orb follows, so move the BALL to the center and let follow() align the
	// orb over it.
	ballHWND := uintptr(b.DebugHWND())
	var bwr rect
	if r, _, _ := pGetWindowRect.Call(ballHWND, unsafePtr(&bwr)); r == 0 {
		t.Fatal("GetWindowRect(ball) failed")
	}
	cx := (int32(work.L) + int32(work.R)) / 2
	cy := (int32(work.T) + int32(work.B)) / 2
	pSetWindowPos.Call(ballHWND, 0,
		uintptr(cx-bwr.width()/2), uintptr(cy-bwr.height()/2), 0, 0,
		uintptr(swpNoSize|swpNoZOrder|swpNoActivate))
	b.followOrb()

	var wr rect
	if r, _, _ := pGetWindowRect.Call(uintptr(orbHWND), unsafePtr(&wr)); r == 0 {
		t.Fatal("GetWindowRect(orb) failed")
	}

	// Known floor below the orb, created and painted on the ui-sta thread.
	// Z-order for the sampling: ball > orb > floor - the floor goes in
	// topmost (above everything), then slots DIRECTLY BELOW the orb
	// (hWndInsertAfter = orb), which leaves the ball > orb > floor stack the
	// z-law in orb_windows.go follow() holds on the live desktop.
	floor := make(chan uintptr, 1)
	b.sta.PostTask(func() { floor <- orbTestCreateColoredFloor(wr) })
	fh := <-floor
	if fh == 0 {
		t.Fatal("colored floor window failed to create")
	}
	pSetWindowPos.Call(fh, uintptr(orbHWND), 0, 0, 0, 0,
		uintptr(swpNoSize|swpNoMove|swpNoActivate))
	defer func() {
		done := make(chan struct{})
		b.sta.PostTask(func() {
			pDestroyWindow.Call(fh)
			close(done)
		})
		<-done
	}()

	corners := [4][2]int32{
		{wr.l + orbCornerInset, wr.t + orbCornerInset},
		{wr.r - orbCornerInset - 1, wr.t + orbCornerInset},
		{wr.l + orbCornerInset, wr.b - orbCornerInset - 1},
		{wr.r - orbCornerInset - 1, wr.b - orbCornerInset - 1},
	}
	fr := int32(orbTestFloorColor & 0xff)
	fg := int32((orbTestFloorColor >> 8) & 0xff)
	fb := int32((orbTestFloorColor >> 16) & 0xff)
	deadline = time.Now().Add(25 * time.Second)
	lit := 0.0
	var readings [4][3]int32
	for {
		hdc, _, _ := pGetDC.Call(0)
		if hdc == 0 {
			t.Fatal("GetDC(NULL) failed")
		}
		// Transparency half: all four corners read the floor color.
		cornersOK := true
		for i, pt := range corners {
			r, g, bl := orbTestPixelAt(hdc, pt[0], pt[1])
			readings[i] = [3]int32{r, g, bl}
			if absi(int(r-fr)) > int(orbCornerTol) || absi(int(g-fg)) > int(orbCornerTol) || absi(int(bl-fb)) > int(orbCornerTol) {
				cornersOK = false
			}
		}
		// Life half: the page must be compositing something somewhere.
		lit = 0
		total := 0
		for y := wr.t + 4; y < wr.b-4; y += 3 {
			for x := wr.l + 4; x < wr.r-4; x += 3 {
				r, g, bl := orbTestPixelAt(hdc, x, y)
				total++
				if absi(int(r-fr)) > orbPageDiffGate || absi(int(g-fg)) > orbPageDiffGate || absi(int(bl-fb)) > orbPageDiffGate {
					lit++
				}
			}
		}
		pReleaseDC.Call(0, hdc)
		if total > 0 {
			lit /= float64(total)
		}
		if cornersOK && lit > 0.01 {
			t.Logf("orb transparent over a known floor: corners %v within +/-%d of (51,51,50); page lit ratio %.2f%%",
				readings, orbCornerTol, lit*100)
			return true, readings, lit
		}
		if time.Now().After(deadline) {
			return false, readings, lit
		}
		time.Sleep(300 * time.Millisecond)
	}
}

// overlayArcTealPixels audits the ULW source DIB for the Downloading progress
// arc (renderer_windows.go step 9): the pal.Accent teal stroke is the ONLY
// green/teal-dominant element the overlay branch may legally paint - the badge
// disc is Danger (red-dominant), the queue dot Info (blue-dominant), the ring
// background BorderSoft (near-white at alpha ~26), the hit-test silhouette
// (alpha 1/255) - so a hue-relative predicate pins the arc's presence on the
// presented frame. The predicate compares PREMULTIPLIED channel/alpha ratios,
// which stay invariant under opacity scaling and antialiasing: Accent must
// have r<=0.6a, g>=0.4a and b<g, which admits both theme rows of --accent
// (#86C2B9 dark / #3E837A light) and rejects every other overlay element.
// Runs on the STA thread (serialized with drawFrame, the same shape the glyph
// audit uses). Returns -1 when there is no DIB to audit.
func overlayArcTealPixels(b *Ball) int {
	done := make(chan int, 1)
	b.sta.PostTask(func() {
		r := b.rend
		if r.bits == nil || r.w <= 0 || r.h <= 0 {
			done <- -1
			return
		}
		buf := unsafe.Slice((*uint32)(r.bits), int(r.w)*int(r.h))
		n := 0
		for _, px := range buf {
			a := px >> 24 // BGRA little-endian: alpha is the top byte
			if a < 100 {
				continue
			}
			bl := int(px & 0xff)
			g := int((px >> 8) & 0xff)
			rr := int((px >> 16) & 0xff)
			if rr*10 <= int(a)*6 && g*10 >= int(a)*4 && bl < g {
				n++
			}
		}
		done <- n
	})
	select {
	case n := <-done:
		return n
	case <-time.After(10 * time.Second):
		return -1
	}
}

// TestBallLiveOrbOverlayFrameHasNoGlyphPixels is the ticket-65 defect the
// owner measured on the real desktop (2026-11-07): with the orb attached, the
// ball window's overlay frame still drew the D2D center text override, so the
// white "69%" download percent rode ON TOP of the glass web orb (evidence
// screenshot D:\tmp\wisp\20261107\orblab\orb-FINAL.png). The overlay branch's
// no-text law is asserted on the pixels the layered window actually presents:
// the real ball runs the exact defect shape - Downloading at 69% plus the
// same SetBadgeText call cmd/wisp's presenter uses - and then the DIB audit
// runs on the STA thread (serialized with drawFrame, so the read is the
// presented frame). It must find ZERO opaque near-white pixels: OnSolid is
// the color every drawText call uses and the ONLY near-white element the
// overlay branch may legally paint (the silhouette is alpha 1/255, the
// progress ring is BorderSoft/Accent, the badge dot is BGInset/Danger, the
// queue dot Info - none near-white), so any nonzero count is text or digit
// pixels stamped on the ball's face.
//
// Wave68 adds the OTHER half the first cut missed (the owner's real-desktop
// follow-up: the outer cyan progress arc of orb-FINAL.png must survive the
// no-text law): the same Downloading frame must still carry the step-9
// progress arc, because step 9 is unguarded in overlay mode ON PURPOSE - the
// arc is the progress signal the page's frozen API cannot draw. A future
// guard that "cleans up" the overlay branch by suppressing step 9 would pass
// every glyph audit while silently removing the arc; the teal audit below
// turns that shape red.
func TestBallLiveOrbOverlayFrameHasNoGlyphPixels(t *testing.T) {
	store, err := OpenPositionStore(filepath.Join(t.TempDir(), "ballpos.json"))
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	requireQuietBallDesktop(t)
	orbTestReapLeakedWebView2(t)

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
	for {
		if b.DebugOrbHWND() != 0 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("orb window did not attach within 15s (WebView2 runtime missing or Embed refused); face %q",
				OrbDebugLine(b.DebugOrbAttached(), b.DebugOrbState()))
		}
		time.Sleep(100 * time.Millisecond)
	}

	// audit posts the DIB walk to the STA thread: the FIFO queue puts it
	// behind every repaint the setters posted, so the count reads the settled
	// overlay frame without polling or races.
	audit := func() int {
		done := make(chan int, 1)
		b.sta.PostTask(func() { done <- b.rend.overlayGlyphPixels() })
		select {
		case n := <-done:
			return n
		case <-time.After(10 * time.Second):
			t.Fatal("overlay glyph audit never ran on the STA thread")
			return -1
		}
	}

	// Phase 1: the measured defect shape. applyStateLocked resets the
	// overlay fields, so state first, then the three host-only values the
	// old branch painted over the face.
	b.SetState(statemachine.StateDownloading)
	b.SetProgress(0.69)
	b.SetBadgeText("69%")
	b.SetBadge(3)
	if n := audit(); n != 0 {
		t.Fatalf("overlay frame carries %d near-white glyph pixels ON the glass face: text/numbers still ride the orb window (ticket-65 \"69%%\" regression)", n)
	}
	// Phase 1b: the SAME frame carries the Downloading progress arc (step 9,
	// the outer cyan stroke - the only teal element the overlay branch may
	// paint). No-text means no GLYPHS, not no progress: the arc is the signal
	// the page cannot draw and its suppression is the silent regression the
	// owner's orb-FINAL vs orb-noglyph2 comparison hunted (wave68).
	if n := overlayArcTealPixels(b); n < 50 {
		t.Fatalf("overlay Downloading frame carries only %d progress-arc pixels (want >=50): the outer cyan arc the page cannot draw is suppressed (ticket-65 arc regression)", n)
	} else {
		t.Logf("overlay Downloading frame keeps the progress arc: %d teal pixels (floor 50)", n)
	}

	// Phase 2: the resting frame is silhouette-only by law.
	b.SetState(statemachine.StateSleeping)
	if n := audit(); n != 0 {
		t.Fatalf("sleeping overlay frame carries %d near-white glyph pixels", n)
	}
}
