//go:build windows

package ball

// Ticket-65 rework, integration leg (A1594): the WebView2 renderer window for
// the upstream liquid-glass orb. The owner's order for this leg is "完全跟人家
// 的一致" - ship LerSent001/orb's own export page (orbasset.go embed) unmodified
// instead of re-authoring a D2D shader, which is exactly what this host does:
// it puts that page on screen and pushes state/audio at it. Nothing here draws.
//
// Shape of the window (the frozen ball set, SPEC-08 §2, minus LAYERED):
// WS_POPUP + WS_EX_TRANSPARENT|TOPMOST|TOOLWINDOW|NOACTIVATE, fixed edge =
// the SAME WindowEdgePx edge the ball window uses, transparency via
// ICoreWebView2Controller2::DefaultBackgroundColor(a=0) on top of the page's
// own premultiplied-alpha canvas, with the window side of the chain made
// per-pixel-alpha-capable by DwmEnableBlurBehindWindow's empty-region form.
// The WebView2 machinery is the jchv/go-webview2 module the panel host
// already ships with; the one controller extension used,
// ICoreWebView2Controller2, is already bound by the module
// (ICoreWebView2Controller.GetICoreWebView2Controller2) - no new COM binding
// is written in this repo.
//
// Transparency chain (ticket 65, orchestrator-measured black corners on the
// first cut, 2026-11-07): four layers each had to be right. The page layer
// was already correct (premultiplied canvas, clear a=0, injected bootstrap
// clears the body fill); the WebView2 layer was already correct
// (DefaultBackgroundColor a=0, documented as supported in windowed hosting
// since 1.0.774.44). The WINDOW layer was the broken one: the first cut
// carried WS_EX_LAYERED but never called SetLayeredWindowAttributes or
// UpdateLayeredWindow, so the top-level's own redirection surface stayed an
// uninitialized (black) layer UNDER the WebView2 child - the page's alpha-0
// pixels composited over that black square. WS_EX_LAYERED buys nothing here
// (this window never presents a ULW bitmap - that is the ball window's D2D
// shape), so it is GONE, and the mainstream Tauri/tao+wry recipe for a
// transparent WebView2 host takes its place: DwmEnableBlurBehindWindow with
// an empty region (CreateRectRgn(0,0,-1,-1), DWM_BB_ENABLE|DWM_BB_BLURREGION)
// makes the window's redirection surface alpha-aware without any blur, which
// is exactly how tao builds transparent windows for wry's windowed WebView2
// (tao window.rs on_create: "making the window transparent"). The DWM layer
// of the chain is therefore satisfied by the same single call; no
// DwmExtendFrameIntoClientArea and no DWMWA_SYSTEMBACKDROP_TYPE are needed
// (those route through the same alpha-aware-surface requirement this call
// already meets).
//
// Threading (D38a/D38b): the orb window is created on, owned by, and torn down
// on the ONE ui-sta thread - no new goroutine exists in this leg. Everything
// from other threads arrives through Ball's existing PostTask/uiRun paths.
//
// Serving (D29): no socket, no port, and no response-construction seam. The
// entry document is delivered under the synthetic host https://wisp.orb via
// ICoreWebView2_3.SetVirtualHostNameToFolderMapping onto a per-user site
// folder written from the embed (orbasset.go), with CORS denied; the CSP
// rides INSIDE the document as an injected <meta http-equiv> (a response-
// header CSP is not available to either delivery shape). A mapping failure
// degrades to NavigateToString, which loses only the WebGPU origin.
//
// The serving history this shape replaced, measured black on 2026-11-07: the
// first cut served the page out of the embed through a WebResourceRequested
// filter, and the jchv wrapper for ICoreWebView2Environment.
// CreateWebResourceResponse reports success by reading the RAW thread
// GetLastError after the COM call - on the orb's STA thread that value reads
// ERROR_NO_TOKEN (1008) persistently, while a direct-vtbl probe read
// HRESULT = S_OK and a live response object in the out param (evidence doc).
// The wrapper discards both and returns nil, the filter served nothing, the
// page never loaded, the window stayed black. Neither current delivery path
// contains a GetLastError success check.
//
// Visibility: the module's Embed never puts IsVisible(true) and the current
// runtimes hand the controller to the app with the composition children
// HIDDEN - measured 2026-11-07: without the explicit browser.Show() below
// the orb window composited nothing (the desktop showed straight through
// it), with it the page's pixels appear. The panel host gets away without
// the call because its window is activated+focused at creation; the orb's
// never-activate show path needs the explicit put.
//
// Fallback (A1594 ③): the D2D renderer is untouched. If WebView2 cannot come
// up here, this file logs one line and returns an error; the caller leaves
// b.orb nil and the D2D path keeps rendering exactly as before.
//
// Clickability + single-ball architecture (ticket 65 rework, 2026-11-07): the
// orb window IS the ball's visible surface - it overlays the ball window's
// rect (follow() center-aligns the two) and carries the whole picture, while
// the BALL window stays the interaction layer: it keeps every frozen
// hit-test/drag/dock/click semantic and stops painting its D2D body when the
// orb is attached (ball_windows.go renderFrame's overlay mode), drawing only
// what the page cannot (approval badge, progress ring, text overrides) plus a
// near-invisible hit-test silhouette. That is the mainstream overlay shape -
// one visible layer (this window) over one interaction layer (the ball
// window) - and it kills the "two balls side by side" state the side-by-side
// follow geometry produced (orchestrator-measured 2026-11-07).
//
// THIS window is and stays pure passthrough: WS_EX_TRANSPARENT hands every
// mouse hit to whatever sits beneath, and the wndproc answers HTTRANSPARENT
// to WM_NCHITTEST on top of that, so the orb can never eat a click meant for
// the ball (or any other window) whether or not the system's WindowFromPoint
// honors WS_EX_TRANSPARENT on a non-layered window - all interaction stays on
// the ball window's frozen semantics.

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
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
)

// orbWindowExStyle: the ball's frozen set MINUS WS_EX_LAYERED, plus
// WS_EX_TRANSPARENT - transparent (every mouse hit falls through to whatever
// sits beneath: the orb is pure passthrough, it must never eat a click),
// topmost, tool window (no taskbar button), never activates. LAYERED is
// deliberately absent: this window never presents a ULW bitmap and never
// calls SetLayeredWindowAttributes, so the bit only served to keep an
// uninitialized black redirection surface under the WebView2 child (the
// transparency-chain note in the file header).
const orbWindowExStyle = uintptr(wsExTransparent | wsExTopmost | wsExToolWindow | wsExNoActivate)

// orbShowCommand: same never-activate show the ball uses.
const orbShowCommand = swShownoactivate

// WM_WINDOWPOSCHANGED / WS_EX_TRANSPARENT: not in win32_windows.go's set, so
// they are named here (this file's own message/style constants).
const (
	wmWindowPosChanged = 0x0047
	wsExTransparent    = 0x00000020
)

// dwmapi plumbing for the window side of the transparency chain: the
// empty-region blur-behind call that makes the top-level's redirection
// surface per-pixel-alpha-aware (Tauri/tao's exact transparent-window recipe,
// file header). The DWM_BLURBEHIND layout is x64-ABI: four 32-bit fields with
// the HRGN aligned inside.
var (
	modDwmapi                  = windows.NewLazySystemDLL("dwmapi.dll")
	pDwmEnableBlurBehindWindow = modDwmapi.NewProc("DwmEnableBlurBehindWindow")
	pCreateRectRgn             = modGDI32.NewProc("CreateRectRgn")
)

const (
	dwmBBEnable     = 0x00000001
	dwmBBBlurRegion = 0x00000002
)

// dwmBlurBehind is DWM_BLURBEHIND.
type dwmBlurBehind struct {
	dwFlags                uint32
	fEnable                uint32 // BOOL
	hRgnBlur               windows.Handle
	fTransitionOnMaximized uint32 // BOOL
}

// enableDWMAlpha turns the window's redirection surface alpha-aware via
// DwmEnableBlurBehindWindow's empty-region form (CreateRectRgn(0,0,-1,-1) is
// the empty region tao uses - no blur is actually requested, only the alpha
// path). A failure is logged, not fatal: without it the window composites
// black where the page is transparent (the measured first-cut defect), but
// the ball still works.
func enableDWMAlpha(hwnd windows.HWND) {
	region, _, _ := pCreateRectRgn.Call(0, 0, ^uintptr(0), ^uintptr(0)) // (0,0,-1,-1): empty
	if region == 0 {
		slog.Warn("ball: orb CreateRectRgn(empty) failed")
		return
	}
	defer pDeleteObject.Call(region)
	bb := dwmBlurBehind{
		dwFlags:  dwmBBEnable | dwmBBBlurRegion,
		fEnable:  1,
		hRgnBlur: windows.Handle(region),
	}
	if hr, _, _ := pDwmEnableBlurBehindWindow.Call(uintptr(hwnd), unsafePtr(&bb)); hr != 0 {
		// HRESULT: 0 = S_OK.
		slog.Warn("ball: orb DwmEnableBlurBehindWindow failed; window may composite black outside the page", "hr", hr)
	}
}

// HWND_TOPMOST as a uintptr literal ((HWND)-1).
const orbHwndTopmost = ^uintptr(0)

// orbWndProcPtr / activeOrb: the class callback plus the package-level
// pointer that routes its messages to the one live orb window - the same
// dispatcher shape the ball window uses (activeBall/wndProcPtr).
var (
	orbWndProcPtr = windows.NewCallback(orbWndProc)
	activeOrb     atomic.Pointer[orbWindow]
	orbClassOnce  sync.Once
	orbClassErr   error
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

	// style is the visual family this window is currently SERVING (UI-thread
	// owned, set at attach and updated by setStyle): it selects the pushed
	// page's seam (window.liquidOrb for siri, window.jarvisOrb for jarvis)
	// and whether audio bands exist at all (jarvis has no audio input).
	style Style

	// servedByMapping records that the virtual-host mapping came up and the
	// page was delivered as https://wisp.orb/<entry> (vs the NavigateToString
	// fallback) - the delivery shape setStyle must reproduce on a swap.
	servedByMapping bool

	// state is the last ball state this window was told about (UI-thread
	// owned). pushAudio reads it to decide whether the mic drives the bands,
	// and setStyle re-pushes it after the page swap.
	state statemachine.State
}

// newOrbWindow creates the orb host window and the WebView2 inside it. It
// MUST run on the ui-sta thread (window + WebView2 + Embed's internal pump
// are all thread-bound). ballHWND is the ball window the orb slots beneath
// (the z-law in follow()'s doc). Every failure mode returns an error and
// nothing else: the D2D fallback decision is the caller's, made from
// b.orb == nil.
func newOrbWindow(s *staThread, sizePx int, dpi uint32, ballHWND windows.HWND) (*orbWindow, error) {
	if err := registerOrbClass(); err != nil {
		return nil, err
	}

	// Same fixed square edge as the ball window, at the destination monitor's
	// DPI: the orb is the ball's visible face, so its host overlays the ball
	// window's rect exactly (follow() center-aligns the two; with equal edges
	// the unclamped alignment is the identical rect).
	edge := int32(WindowEdgePx(sizePx, dpi))

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
	// Window side of the transparency chain: alpha-aware redirection surface
	// (the removed WS_EX_LAYERED never provided one). Before the WebView2
	// child composites, an alpha-aware surface with nothing painted in it is
	// FULLY transparent - no black flash between window creation and attach.
	enableDWMAlpha(windows.HWND(hwnd))

	o := &orbWindow{sta: s, hwnd: windows.HWND(hwnd), style: ActiveStyle()}

	// Transparency, WebView2 layer, from BROWSER BIRTH: the documented cure
	// for the API put's own known issue - "setting the color by API can
	// still leave the app with a white flicker before the
	// DefaultBackgroundColor takes effect" (WebView2 docs,
	// ICoreWebView2Controller2.put_DefaultBackgroundColor). Measured
	// 2026-11-07 on real desktop boots: some instances composited the page
	// over an opaque WHITE default for their whole lifetime despite the
	// successful API put below. The environment variable applies the
	// transparent default when the browser process is born; the API put
	// stays as the belt (get returns the env value; later updates go
	// through the API). Process-wide by nature: the panel host's WebView2
	// shares this process env, and its own page paints an opaque CSS
	// background, so an undefined-background window is not a shape it ever
	// shows.
	if err := os.Setenv("WEBVIEW2_DEFAULT_BACKGROUND_COLOR", "0x00FFFFFF"); err != nil {
		slog.Warn("ball: orb default-background env not set; white boot possible", "err", err)
	}

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

	// The bootstrap runs BEFORE the page's module script (document-created
	// injection): it clears the export page's own dark tab background so the
	// pixels outside the orb are the window's transparency, not that fill.
	browser.Init(orbPageBootstrap)

	// Serving (D29), second shape - a REAL https origin, because the SIRI
	// page's shader is WebGPU and the runtime only exposes navigator.gpu to a
	// secure context (measured 2026-11-07: under NavigateToString the page ran
	// fine - liquidOrb=object, ready=complete - but navigator.gpu was
	// undefined and the page itself printed "WebGPU is not supported in this
	// environment"). The active style's page is materialized byte-for-byte
	// (plus the injected CSP meta, orbasset.go) into a per-user folder - BOTH
	// style entries are written so a runtime switch is a pure navigation -
	// and the synthetic host https://wisp.orb is mapped onto it with CORS
	// denied - the offline no-socket shape survives: the mapping exists only
	// inside THIS browser instance, the host name is under a TLD no resolver
	// answers, and the meta CSP still cuts every network verb. The module's
	// SetVirtualHostNameToFolderMapping wrapper reports failure off the raw
	// thread GetLastError (the same defect that broke the response seam), so
	// a reported failure is logged and answered with the NavigateToString
	// fallback rather than trusted blindly. (The style itself was read into
	// o.style at window creation, above.)
	siteDir := filepath.Join(orbDataPath(), "site")
	served := false
	if err := orbWriteSite(siteDir); err != nil {
		slog.Warn("ball: orb site materialization failed; WebGPU origin unavailable", "err", err)
	} else if w3 := browser.GetICoreWebView2_3(); w3 == nil {
		slog.Warn("ball: orb ICoreWebView2_3 unavailable; WebGPU origin unavailable")
	} else if err := w3.SetVirtualHostNameToFolderMapping(
		orbSiteHost, siteDir, orbedge.COREWEBVIEW2_HOST_RESOURCE_ACCESS_KIND_DENY_CORS); err != nil {
		slog.Warn("ball: orb virtual host mapping reported failure (the wrapper's raw-GetLastError check may lie); falling back to NavigateToString", "err", err)
	} else {
		_, _, entryURL := orbPageForStyle(o.style)
		browser.Navigate(entryURL)
		served = true
		o.servedByMapping = true
	}
	if !served {
		// Degraded delivery: the page still loads (interface, states all work)
		// but the SIRI page's runtime hides navigator.gpu off this origin, so
		// that page shows its own "WebGPU is not supported" face. The jarvis
		// page is WebGL and renders off this origin too.
		doc, _, _ := orbPageForStyle(o.style)
		browser.NavigateToString(string(doc))
	}

	// Controller visibility: the module's Embed NEVER puts IsVisible(true)
	// and on the current runtimes the controller comes up with the children
	// hidden - measured 2026-11-07: without this call the orb window
	// composited NOTHING (screen showed straight through it), with it the
	// page's own pixels appear. The panel host gets away without it because
	// its window is activated+focused at creation; the orb's never-activate
	// show path needs the explicit put.
	if err := browser.Show(); err != nil {
		slog.Warn("ball: orb controller Show failed", "err", err)
	}

	// Never-activate show, then one z-order slot: DIRECTLY BELOW the ball
	// window (the z-law in follow()'s doc - the ball is the top layer, the
	// orb its face beneath it). The first follow() re-slots after every move.
	pShowWindow.Call(hwnd, orbShowCommand)
	pSetWindowPos.Call(hwnd, uintptr(ballHWND), 0, 0, 0, 0,
		uintptr(swpNoSize|swpNoMove|swpNoActivate))
	return o, nil
}

// orbDataPath picks a PER-PROCESS user-data folder for THIS browser instance.
// History: it used to hand EVERY process the same folder
// (%LocalAppData%\wisp\orb), so `go test ./cmd/wisp` - several ball window
// tests, each booting WebView2, across several processes - self-collided on
// one user-data folder: a contended CreateEnvironment completed with S_OK and
// a NULL environment out-param, the panic died on the recover boundary, the
// STA thread was gone and ball.New waited forever (whole-package timeout,
// ticket 408). WebView2's guidance for multiple instances is one user-data
// folder per process, so the folder is now <base>\wisp\orb\<pid>-<8 hex>.
// Decision record (ticket 408): pid alone - candidate 甲 - still collides
// when Windows reuses a pid and hands the new process a dead predecessor's
// folder with its stale locks; the random suffix buys that class away. The
// exe name - candidate 丙 - adds no uniqueness among live processes and
// lengthens the path, so it was dropped. The path MUST be stable for the
// process lifetime: one attach reads it twice (browser.DataPath above and
// the site folder below), so the value is computed exactly once.
func orbDataPath() string {
	orbDataPathOnce.Do(func() {
		orbDataDir, orbDataLock = acquireOrbDataDir()
	})
	return orbDataDir
}

var (
	orbDataPathOnce sync.Once
	orbDataDir      string
	// orbDataLock is held OPEN from before the folder exists until process
	// exit - that open handle IS the sweep's liveness signal.
	orbDataLock windows.Handle
)

// orbDataBase is the folder root: LOCALAPPDATA, falling back to the temp dir
// when unset (the old fallback, kept byte-for-byte).
func orbDataBase() string {
	base := os.Getenv("LOCALAPPDATA")
	if base == "" {
		base = os.TempDir()
	}
	return base
}

// orbDataDirFor is the pure form of the path: the leaf this process would use,
// rooted at base. Tests drive it directly so the leaf shape and the base
// fallback have teeth without re-exec.
func orbDataDirFor(base string) string {
	return filepath.Join(base, "wisp", "orb", fmt.Sprintf("%d-%s", os.Getpid(), orbInstanceSuffix()))
}

func acquireOrbDataDir() string {
	base := orbDataBase()
	root := filepath.Join(base, "wisp", "orb")
	orbSweepStaleDirs(root)
	if err := os.MkdirAll(root, 0o755); err != nil {
		slog.Warn("ball: orb user-data root mkdir failed", "err", err)
	}
	dir := orbDataDirFor(base)
	// Lock FIRST, folder second: the sweep's only liveness signal is an
	// exclusively held lock file, so a folder must never be observable
	// without its lock (a sweep would otherwise race the creator). The
	// handle stays open until process exit. A failure here logs and
	// continues - uniqueness still holds; only this dir's own sweep
	// protection is lost.
	h, err := windows.CreateFile(dir+".lock", windows.GENERIC_READ|windows.GENERIC_WRITE,
		0 /* no sharing: any second opener reads the owner as alive */, nil,
		windows.OPEN_ALWAYS, windows.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		slog.Warn("ball: orb user-data lock not held; this dir is invisible to the stale sweep", "err", err)
	} else {
		orbDataLock = h
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		slog.Warn("ball: orb user-data dir mkdir failed", "err", err)
	}
	return dir
}

// orbSweepStaleDirs bounds the garbage the per-process folder leaves behind:
// an entry whose lock file can be opened EXCLUSIVELY has a dead owner - the
// owner holds that lock from before the folder exists until process exit -
// so it is removed, lock and all. A busy lock (an alive owner, including
// this process) and a folder with no lock at all are skipped. The pid in the
// name is never consulted, so pid reuse cannot strand an entry forever.
// Best-effort by design: any failure leaves the entry for the next attach.
func orbSweepStaleDirs(root string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return
	}
	for _, entry := range entries {
		name := entry.Name()
		if !orbDataLeafRe.MatchString(name) {
			continue
		}
		h, err := windows.CreateFile(filepath.Join(root, name+".lock"), windows.GENERIC_READ,
			0, nil, windows.OPEN_EXISTING, windows.FILE_ATTRIBUTE_NORMAL, 0)
		if err != nil {
			continue // alive owner or unreadable: conservative skip
		}
		windows.CloseHandle(h)
		if entry.IsDir() {
			os.RemoveAll(filepath.Join(root, name))
		}
		os.Remove(filepath.Join(root, name+".lock"))
	}
}

// orbDataLeafRe matches the per-process leaf `<pid>-<8 hex>` and its `.lock`
// file spelling; anything else under the root (the pre-fix shared layout's
// own EBWebView/site children, say) is foreign and never touched.
var orbDataLeafRe = regexp.MustCompile(`^\d{1,10}-[0-9a-f]{8}(\.lock)?$`)

func orbInstanceSuffix() string {
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		// Unreachable in practice; pid-level uniqueness among live processes
		// still holds, and the lock-based sweep keeps pid reuse safe.
		return "00000000"
	}
	return hex.EncodeToString(b[:])
}

// pushState maps and pushes one ball state to the SERVED page (STA thread).
// The seam follows the window's style: the siri page exposes window.liquidOrb
// (two names, throws on unknown - orbstate.go orbStateFor), the jarvis page
// window.jarvisOrb (four names, ignores unknown - jarvisStateFor).
func (o *orbWindow) pushState(s statemachine.State) {
	o.state = s
	if o.browser == nil {
		return
	}
	if o.style == StyleJarvis {
		o.browser.Eval(orbJarvisStateScript(jarvisStateFor(s)))
		return
	}
	o.browser.Eval(orbStateScript(orbStateFor(s)))
}

// pushAudio feeds one mic envelope to the page's bands (STA thread). The
// caller (Ball.applyLevel) already gated this on the audio-driven family;
// the double check here keeps the seam honest on its own. The jarvis page is
// skipped outright: the upstream component has no audio-band input (verified
// states.ts/renderer.ts/JarvisOrb.tsx) and none is invented.
func (o *orbWindow) pushAudio(level float32) {
	if o.browser == nil || o.style == StyleJarvis || !orbAudioFor(o.state) {
		return
	}
	l := float64(level)
	o.browser.Eval(orbAudioScript(l, l, l, l))
}

// setStyle swaps the served page in place on a runtime style switch (STA
// thread): a pure navigation between the two site entries (orbWriteSite wrote
// both at attach), then a re-push of the window's own last state so the new
// page opens on the face the ball is in. The re-push rides the same guarded
// seam as attach-time seeding: if it lands before the page's module script
// has defined the interface object the guard no-ops and the next SetState
// corrects the face - the identical race the siri attach accepted and
// measured benign. Served=false (the NavigateToString fallback) re-runs the
// fallback with the other style's document.
func (o *orbWindow) setStyle(style Style) {
	if o.browser == nil || o.style == style {
		return
	}
	o.style = style
	doc, _, entryURL := orbPageForStyle(style)
	if entryURL != "" && o.servedByMapping {
		o.browser.Navigate(entryURL)
	} else {
		o.browser.NavigateToString(string(doc))
	}
	o.pushState(o.state)
}

// orbFollowGap is gone: the orb no longer floats BESIDE the ball (that
// geometry drew a second ball on screen - orchestrator-measured, ticket 65).
// The orb overlays the ball window's rect; see orbOverlayXY.

// follow places the orb over the ball's current window rect (center-aligned,
// work-area clamped) and slots it DIRECTLY BELOW the ball window in z-order.
// Runs on the ui-sta thread (ball drag move/end, DPI change, attach), which is
// why no lock guards the window geometry.
//
// Z-law (ticket 65 rework, measured): the ball window is the TOP layer and
// the orb sits immediately beneath it. The pointer therefore always lands on
// the ball first (WindowFromPoint at the ball centre must answer the ball -
// the frozen interaction contract), the ball's own transparent pixels let the
// orb show through, and the orb's WS_EX_TRANSPARENT + HTTRANSPARENT wndproc
// remain as passthrough belts for real input routing. The reverse order
// measured broken: a non-layered WS_EX_TRANSPARENT window is NOT skipped by
// WindowFromPoint, so an orb above the ball answered WindowFromPoint(centre)
// the moment follow() aligned the two.
func (o *orbWindow) follow(ballHWND windows.HWND) {
	if o.browser == nil || o.hwnd == 0 {
		return
	}
	var wr rect
	r, _, _ := pGetWindowRect.Call(uintptr(ballHWND), unsafePtr(&wr))
	if r == 0 {
		return
	}
	var or rect
	if r, _, _ := pGetWindowRect.Call(uintptr(o.hwnd), unsafePtr(&or)); r == 0 {
		return
	}
	work, workOk := monitorRectsForWindow(ballHWND)
	x, y := orbOverlayXY(wr, or.width(), work.Work, workOk)
	// hWndInsertAfter = the ball window: the orb slots DIRECTLY BELOW it
	// (the z-law above), never at the topmost band's head.
	pSetWindowPos.Call(uintptr(o.hwnd), uintptr(ballHWND),
		uintptr(x), uintptr(y), 0, 0,
		uintptr(swpNoSize|swpNoActivate))
}

// orbOverlayXY is the pure placement decision follow() applies, split out so
// the suite can pin it against real work-area shapes without windows. The orb
// is centered on the ball window's rect (equal edges make that the identical
// rect), then clamped fully inside the work area when monitor info exists.
// A workOk=false result (no monitor info available) keeps the raw alignment
// rather than guessing a boundary that was never read.
func orbOverlayXY(ball rect, orbEdge int32, work Rect, workOk bool) (int32, int32) {
	x := (ball.l+ball.r)/2 - orbEdge/2
	y := (ball.t+ball.b)/2 - orbEdge/2
	if !workOk {
		return x, y
	}
	if x+orbEdge > int32(work.R) {
		x = int32(work.R) - orbEdge
	}
	if x < int32(work.L) {
		x = int32(work.L)
	}
	if y+orbEdge > int32(work.B) {
		y = int32(work.B) - orbEdge
	}
	if y < int32(work.T) {
		y = int32(work.T)
	}
	return x, y
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
	// Ticket 65 (docs/evidence/s1/65-jchv-controller-leak.md): the WebView2
	// controller must be CLOSED - not just unreferenced - before the host
	// window dies. Without controller.Close the browser process tree
	// (msedgewebview2) outlives the orb, and every create/destroy round leaks
	// another set of processes; across a count=10 live suite that accumulation
	// degraded rendering round over round until pages could not paint at all.
	// Close (fork wrapper, third_party/go-webview2) unregisters handlers and
	// lets the browser processes exit; Release drops the controller's COM
	// reference the wrapper's AddRef-side of Embed took. Both run on this STA
	// thread, before o.browser is dropped, so no other call can reach the
	// dying controller in between.
	if o.browser != nil {
		if c := o.browser.GetController(); c != nil {
			if err := c.Close(); err != nil {
				slog.Warn("ball: orb controller Close failed", "err", err)
			}
			c.Release()
		}
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
	case wmNCHitTest:
		// Pure passthrough, second belt: the system's hit-test walk asks this
		// window first (WindowFromPoint consults WM_NCHITTEST top-down) and
		// HTTRANSPARENT sends it to whatever sits beneath - the same answer
		// the ball window gives outside its clickable circle. WS_EX_TRANSPARENT
		// already asks for this; on a NON-layered window the system's honoring
		// of that bit for top-level hit tests is not something to bet on, and
		// the answer here makes the passthrough a property of THIS wndproc
		// instead of a style-interpretation detail.
		return htTransparentResult()
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
	o, err := newOrbWindow(s, b.opts.SizePx, getDpiForWindow(b.hwnd), b.hwnd)
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

// ApplyStyle switches the style family AT RUNTIME: SetStyle (the process-wide
// family the D2D renderer repaints from) plus the served page swap that makes
// the WebView2 face follow it. Safe from any goroutine; returns false for a
// name outside the vocabulary (styleGlassName|styleJarvisName) with nothing
// changed. With the orb not attached (D2D fallback), only SetStyle runs - the
// fallback face is the renderer's jarvis/glass branch, untouched here.
func (b *Ball) ApplyStyle(name string) bool {
	if !SetStyle(name) {
		return false
	}
	b.sta.PostTask(func() {
		if o := b.orb.Load(); o != nil {
			o.setStyle(ActiveStyle())
		}
	})
	return true
}

// followOrb is the ball-window move hook: keep the orb next to the ball.
func (b *Ball) followOrb() {
	if o := b.orb.Load(); o != nil {
		o.follow(b.hwnd)
	}
}

// DebugOrbHWND exposes the orb window handle (debug harness/evidence only;
// zero when the D2D fallback is active).
func (b *Ball) DebugOrbHWND() windows.HWND {
	if o := b.orb.Load(); o != nil {
		return o.hwnd
	}
	return 0
}
