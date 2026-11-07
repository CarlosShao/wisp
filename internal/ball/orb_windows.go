//go:build windows

package ball

// Ticket-65 rework, integration leg (A1594): the WebView2 renderer window for
// the upstream liquid-glass orb. The owner's order for this leg is "完全跟人家
// 的一致" - ship LerSent001/orb's own export page (orbasset.go embed) unmodified
// instead of re-authoring a D2D shader, which is exactly what this host does:
// it puts that page on screen and pushes state/audio at it. Nothing here draws.
//
// Shape of the window (the frozen ball set, SPEC-08 §2): WS_POPUP +
// WS_EX_LAYERED|TOPMOST|TOOLWINDOW|NOACTIVATE, fixed edge = the configured orb
// edge, transparency via ICoreWebView2Controller2::DefaultBackgroundColor(a=0)
// on top of the page's own premultiplied-alpha canvas. The WebView2 machinery
// is the jchv/go-webview2 module the panel host already ships with; the one
// controller extension used, ICoreWebView2Controller2, is already bound by the
// module (ICoreWebView2Controller.GetICoreWebView2Controller2) - no new COM
// binding is written in this repo.
//
// Threading (D38a/D38b): the orb window is created on, owned by, and torn down
// on the ONE ui-sta thread - no new goroutine exists in this leg. Everything
// from other threads arrives through Ball's existing PostTask/uiRun paths.
//
// Serving (D29): no socket, no port. The page is answered out of the embed by
// a WebResourceRequested filter under a synthetic origin that no resolver can
// answer, and every served response carries the CSP header with connect-src
// 'none' - the same no-network shape the panel host uses.
//
// Fallback (A1594 ③): the D2D renderer is untouched. If WebView2 cannot come
// up here, this file logs one line and returns an error; the caller leaves
// b.orb nil and the D2D path keeps rendering exactly as before. The render
// guard is three early returns (renderFrame/present), and the D2D code below
// them is dead while the orb is up, not deleted.

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"

	orbedge "github.com/jchv/go-webview2/pkg/edge"
	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

const (
	// orbWindowClassName is this host's class, deliberately NOT the ball's
	// "WispBallWindow": a desktop census can tell the two windows apart by
	// name, and the live tests' quiet-desktop rule is written against the
	// ball class only.
	orbWindowClassName = "WispOrbWindow"

	// orbWindowEdge is the fixed square edge of the orb host window (96-dpi
	// px, scaled by the destination monitor's DPI like the ball edge). The
	// page's shader fits the orb to min(w,h) with a soft margin, so a little
	// extra room around the configured size is what keeps the glow band and
	// the contour deform inside the window. 1.5x, rounded.
	orbWindowEdgeBase = BallSizeMaxPx*96/64 + 24
)

// orbWindowExStyle: the ball's frozen set minus nothing - layered (the page
// composites per-pixel through ULW-style transparency), topmost, tool window
// (no taskbar button), never activates.
const orbWindowExStyle = uintptr(wsExLayered | wsExTopmost | wsExToolWindow | wsExNoActivate)

// orbShowCommand: same never-activate show the ball uses.
const orbShowCommand = swShownoactivate

// WM_WINDOWPOSCHANGED: the hook the ball's window forwards so the orb follows
// a dragged/docked/DPI-moved ball. Not in win32_windows.go's set, so it is
// named here (this leg's only message constant).
const wmWindowPosChanged = 0x0047

// HWND_TOPMOST as a uintptr literal ((HWND)-1).
const orbHwndTopmost = ^uintptr(0)

// orbWndProcPtr / activeOrb: the class callback plus the package-level
// pointer that routes its messages to the one live orb window - the same
// dispatcher shape the ball window uses (activeBall/wndProcPtr).
var (
	orbWndProcPtr   = windows.NewCallback(orbWndProc)
	activeOrb       atomic.Pointer[orbWindow]
	orbClassOnce    sync.Once
	orbClassErr     error
	orbModKernel    = windows.NewLazySystemDLL("kernel32.dll")
	orbSetLastError = orbModKernel.NewProc("SetLastError")
)

// registerOrbClass registers the orb window class once per process.
func registerOrbClass() error {
	orbClassOnce.Do(func() {
		clsName := utf16(orbWindowClassName)
		wc := wndClassEx{
			size:      uint32(unsafe.Sizeof(wndClassEx{})),
			wndProc:   orbWndProcPtr,
			instance:  windows.Handle(moduleHandle()),
			cursor:    windows.Handle(loadArrowCursor()),
			className: clsName,
		}
		if r, _, _ := pRegisterClassExW.Call(unsafePtr(&wc)); r == 0 {
			orbClassErr = fmt.Errorf("ball: RegisterClassExW(%s) failed", orbWindowClassName)
		}
	})
	return orbClassErr
}

// orbWindow is the WebView2 host. All fields below are owned by the ui-sta
// thread; other threads reach it only through Ball's posted closures.
type orbWindow struct {
	sta     *staThread
	hwnd    windows.HWND
	browser *orbedge.Chromium

	// state is the last ball state this window was told about (UI-thread
	// owned). pushAudio reads it to decide whether the mic drives the bands.
	state statemachine.State
}

// newOrbWindow creates the orb host window and the WebView2 inside it. It
// MUST run on the ui-sta thread (window + WebView2 + Embed's internal pump
// are all thread-bound). Every failure mode returns an error and nothing
// else: the D2D fallback decision is the caller's, made from b.orb == nil.
func newOrbWindow(s *staThread, sizePx int, dpi uint32) (*orbWindow, error) {
	if err := registerOrbClass(); err != nil {
		return nil, err
	}

	// Fixed square edge at the destination monitor's DPI (same scaling rule
	// as WindowEdgePx for the ball; the orb window has no ring margins of its
	// own, so the base is this file's constant scaled linearly).
	edge := int32(int(orbWindowEdgeBase) * int(dpi) / 96)
	_ = sizePx // the window is one fixed square; the page fits the orb itself

	clsName := utf16(orbWindowClassName)
	title := utf16("Wisp Orb")
	hwnd, _, callErr := pCreateWindowExW.Call(
		orbWindowExStyle,
		unsafePtr(clsName),
		unsafePtr(title),
		uintptr(wsPopup),
		0, 0, uintptr(edge), uintptr(edge),
		0, 0, uintptr(moduleHandle()), 0)
	if hwnd == 0 {
		return nil, fmt.Errorf("ball: orb CreateWindowExW: %v", callErr)
	}

	o := &orbWindow{sta: s, hwnd: windows.HWND(hwnd)}

	browser := orbedge.NewChromium()
	// MessageCallback stays nil on purpose: the upstream page never posts to
	// the host, and the library's MessageReceived arm would bounce every
	// message back into the page. Nothing arrives on this channel; if a
	// future leg needs page->host messages, that is a NEW seam, not this one.
	browser.DataPath = orbDataPath()
	o.browser = browser

	if !browser.Embed(uintptr(o.hwnd)) {
		o.destroy()
		return nil, fmt.Errorf("ball: orb WebView2 Embed failed (runtime missing or blocked)")
	}
	browser.Resize()

	// Context menu + devtools off, the panel host's shape. A failure here is
	// a settings-only downgrade, not an attach failure: log and continue.
	if settings, err := browser.GetSettings(); err != nil {
		slog.Warn("ball: orb settings read failed", "err", err)
	} else {
		if err := settings.PutAreDefaultContextMenusEnabled(false); err != nil {
			slog.Warn("ball: orb context-menu disable failed", "err", err)
		}
		if err := settings.PutAreDevToolsEnabled(false); err != nil {
			slog.Warn("ball: orb devtools disable failed", "err", err)
		}
	}

	// Transparency: controller2 (IID {c979903e-d4ca-4228-92eb-47ee3fa96eab})
	// puts an alpha-0 default background behind the page's premultiplied
	// canvas. Without this the window is a black square - the page itself is
	// correct (alpha 0 outside the orb); only the HOST background matters, and
	// the injected bootstrap (orbasset.go) clears the export page's own dark
	// body fill so that host transparency actually shows through.
	if c2 := browser.GetController().GetICoreWebView2Controller2(); c2 == nil {
		o.destroy()
		return nil, fmt.Errorf("ball: orb controller2 QI failed (transparency unavailable)")
	} else if err := c2.PutDefaultBackgroundColor(orbedge.COREWEBVIEW2_COLOR{A: 0, R: 0, G: 0, B: 0}); err != nil {
		o.destroy()
		return nil, fmt.Errorf("ball: orb PutDefaultBackgroundColor: %v", err)
	}

	// Serving: the embed answers everything under the synthetic origin; a miss
	// is an empty 404 with the CSP header on it.
	browser.WebResourceRequestedCallback = o.serve
	browser.AddWebResourceRequestedFilter(orbAssetFilterPattern, orbedge.COREWEBVIEW2_WEB_RESOURCE_CONTEXT_ALL)

	// The bootstrap runs BEFORE the page's module script (document-created
	// injection): it clears the export page's own dark tab background so the
	// pixels outside the orb are the window's transparency, not that fill.
	browser.Init(orbPageBootstrap)
	browser.Navigate(orbEntryURL)

	// Never-activate show, then one explicit topmost re-assert so the orb
	// lands in the same topmost band as the ball that owns it.
	pShowWindow.Call(hwnd, orbShowCommand)
	pSetWindowPos.Call(hwnd, orbHwndTopmost, 0, 0, 0, 0,
		uintptr(swpNoSize|swpNoMove|swpNoActivate))
	return o, nil
}

// orbDataPath picks a per-user user-data folder for THIS browser instance.
// Sharing the panel's folder would tie two surfaces to one runtime state, and
// the library's %AppData%\<exe> default would collide with the panel's. An
// empty LOCALAPPDATA falls back to the temp dir rather than a relative path.
func orbDataPath() string {
	base := os.Getenv("LOCALAPPDATA")
	if base == "" {
		base = os.TempDir()
	}
	return filepath.Join(base, "wisp", "orb")
}

// serve answers one intercepted request. Same contract as the panel host:
// the entry document out of the embed, an EMPTY 404 for every other path
// (never a fallback document), the CSP header on every response.
func (o *orbWindow) serve(request *orbedge.ICoreWebView2WebResourceRequest, args *orbedge.ICoreWebView2WebResourceRequestedEventArgs) {
	uri, err := request.GetUri()
	if err != nil {
		slog.Warn("ball: orb resource request without a readable URI", "err", err)
		return
	}
	if !strings.HasPrefix(uri, orbOrigin) {
		// Unreachable through the filter pattern; left alone rather than
		// answered by a guess.
		return
	}
	status, reason, contentType := 404, "Not Found", ""
	data := []byte(nil)
	if strings.HasPrefix(uri, orbEntryURL) {
		data, status, reason, contentType = orbPageHTML, 200, "OK", "text/html; charset=utf-8"
	}
	env := o.browser.Environment()
	if env == nil {
		slog.Warn("ball: orb request with no environment to build a response", "uri", uri)
		return
	}
	response, err := orbBuildResponse(env, data, status, reason, orbResponseHeaders(status, contentType))
	if err != nil {
		slog.Warn("ball: orb CreateWebResourceResponse failed", "uri", uri, "err", err)
		return
	}
	if err := args.PutResponse(response); err != nil {
		slog.Warn("ball: orb PutResponse failed", "uri", uri, "err", err)
	}
}

// orbResponseHeaders is the header LINE one served response carries. 200 gets
// the Content-Type Resolve decided; every response gets the CSP.
func orbResponseHeaders(status int, contentType string) string {
	if status == 200 {
		return "Content-Type: " + contentType + "\r\n" + orbCSPHeader
	}
	return orbCSPHeader
}

// orbBuildResponse wraps env.CreateWebResourceResponse with the bounded retry
// the panel host measured (2026-10-13): that helper reports failure off a
// stale GetLastError and a SUCCESSFUL call can read errno 1008 on its first
// invocation per cold start. Retry, loudly; the stale value is not stable.
func orbBuildResponse(env *orbedge.ICoreWebView2Environment, data []byte, status int, reason, headers string) (*orbedge.ICoreWebView2WebResourceResponse, error) {
	const attempts = 3
	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		orbSetLastError.Call(0)
		response, err := env.CreateWebResourceResponse(data, status, reason, headers)
		if err == nil {
			return response, nil
		}
		lastErr = err
		slog.Warn("ball: orb response create retry", "attempt", attempt, "err", err)
	}
	return nil, fmt.Errorf("ball: orb CreateWebResourceResponse after %d attempts: %w", attempts, lastErr)
}

// pushState maps and pushes one ball state to the page (STA thread).
func (o *orbWindow) pushState(s statemachine.State) {
	o.state = s
	if o.browser == nil {
		return
	}
	o.browser.Eval(orbStateScript(orbStateFor(s)))
}

// pushAudio feeds one mic envelope to the page's bands (STA thread). The
// caller (Ball.applyLevel) already gated this on the audio-driven family;
// the double check here keeps the seam honest on its own.
func (o *orbWindow) pushAudio(level float32) {
	if o.browser == nil || !orbAudioFor(o.state) {
		return
	}
	l := float64(level)
	o.browser.Eval(orbAudioScript(l, l, l, l))
}

// follow places the orb to the right of the ball's current window rect. Runs
// on the ui-sta thread (from the ball's WM_WINDOWPOSCHANGED arm), which is
// why no lock guards the window geometry.
func (o *orbWindow) follow(ballHWND windows.HWND) {
	if o.browser == nil || o.hwnd == 0 {
		return
	}
	var wr rect
	r, _, _ := pGetWindowRect.Call(uintptr(ballHWND), unsafePtr(&wr))
	if r == 0 {
		return
	}
	gap := int32(12)
	pSetWindowPos.Call(uintptr(o.hwnd), orbHwndTopmost,
		uintptr(wr.r+gap), uintptr(wr.t), 0, 0,
		uintptr(swpNoSize|swpNoActivate))
}

// destroy tears the browser and the window down (STA thread). Order matters,
// measured in the live test (ticket 65): pDestroyWindow delivers
// WM_DESTROY/WM_WINDOWPOSCHANGED SYNCHRONOUSLY on this thread, so routing
// must stop BEFORE the destroy call - a message that lands mid-destroy would
// otherwise reach the library's NotifyParentWindowPositionChanged with a
// controller whose object is already going away (nil deref inside the module,
// recovered by the observe boundary but leaving the window alive).
func (o *orbWindow) destroy() {
	if activeOrb.Load() == o {
		activeOrb.Store(nil)
	}
	o.browser = nil
	if o.hwnd != 0 {
		if r, _, err := pDestroyWindow.Call(uintptr(o.hwnd)); r == 0 {
			slog.Warn("ball: orb DestroyWindow failed", "hwnd", uintptr(o.hwnd), "err", err)
		}
		o.hwnd = 0
	}
}

// orbWndProc: the class callback. Windows delivers messages for every window
// of the class here; anything that is not the live orb window falls through
// to DefWindowProc, exactly like the ball's dispatcher.
func orbWndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	o := activeOrb.Load()
	if o == nil || o.hwnd == 0 || hwnd != uintptr(o.hwnd) {
		r, _, _ := pDefWindowProcW.Call(hwnd, msg, wParam, lParam)
		return r
	}
	switch msg {
	case wmDestroy:
		// The owner (Ball.Close) drives teardown; this arm only stops the
		// library's per-thread quit from being posted twice - the window is
		// destroyed by pDestroyWindow from destroy(), which lands here.
		return 0
	case wmWindowPosChanged:
		// The controller must track its parent window's moves. The browser
		// nil-check is teardown: pDestroyWindow itself raises this message,
		// and destroy() has already dropped the browser by then (the
		// routing-level guard above usually answers first).
		if o.browser != nil {
			_ = o.browser.NotifyParentWindowPositionChanged()
		}
		r, _, _ := pDefWindowProcW.Call(hwnd, msg, wParam, lParam)
		return r
	}
	r, _, _ := pDefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return r
}

// attachOrb is the New-side attach attempt: create the WebView2 orb on the
// ui-sta thread. NEVER fails the ball - a failure logs one line and leaves
// b.orb nil, which keeps the D2D fallback rendering (A1594 ③).
func (b *Ball) attachOrb(s *staThread) {
	o, err := newOrbWindow(s, b.opts.SizePx, getDpiForWindow(b.hwnd))
	if err != nil {
		slog.Warn("ball: webview orb not attached; D2D fallback stays active", "err", err)
		return
	}
	if !activeOrb.CompareAndSwap(nil, o) {
		// One orb per process; a second attach attempt loses and cleans up.
		o.destroy()
		return
	}
	b.orb.Store(o)
	// Seed the page with the state the ball already carries, and the audio
	// level it already has, so the orb never opens on the wrong face.
	o.pushState(b.curState)
	if orbAudioFor(b.curState) {
		o.pushAudio(b.liqRaw)
	}
}

// detachOrb tears the orb down on the STA thread (Ball.Close's arm).
func (b *Ball) detachOrb() {
	if o := b.orb.Swap(nil); o != nil {
		o.destroy()
	}
}

// followOrb is the ball-window move hook: keep the orb next to the ball.
func (b *Ball) followOrb() {
	if o := b.orb.Load(); o != nil {
		o.follow(b.hwnd)
	}
}

// orbLive is the render guard: while the orb owns the pixels, the D2D
// renderer must not paint into the layered window behind the transparent
// WebView2 child.
func (b *Ball) orbLive() bool { return b.orb.Load() != nil }

// DebugOrbHWND exposes the orb window handle (debug harness/evidence only;
// zero when the D2D fallback is active).
func (b *Ball) DebugOrbHWND() windows.HWND {
	if o := b.orb.Load(); o != nil {
		return o.hwnd
	}
	return 0
}
