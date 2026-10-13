//go:build windows

package main

// The panel page host drives WebView2's exported pkg/edge layer itself (ticket
// 299 AC#1, shape ruled 2026-10-13): this file owns the window class, the
// message pump, the RPC pipe and the resource filter that go-webview2's
// high-level wrapper used to keep to itself.
//
// WHY NOT THE WRAPPER ANY MORE. The one registration D29 names -
// AddWebResourceRequestedFilter fed from the embedded bytes - needs the
// *edge.Chromium the wrapper creates internally, and the wrapper is exactly the
// layer that will not hand it out: it stores the object in an unexported field
// (webview.go:49-59 `browser browser`, assigned at :107) behind an unexported
// type, and its 12-method WebView interface (common.go:26-84) has no getter.
// pkg/edge itself is fully exported (chromium.go:47 NewChromium, :42 the
// WebResourceRequestedCallback arm, :292 AddWebResourceRequestedFilter, :72
// Embed which pumps this thread until environment AND controller are built), so
// the host drives that layer directly while still satisfying webview2.WebView -
// which is what keeps panel_resident_windows.go's post()/Run() contract, and
// every ruler 33/35/255 already wrote against it, unchanged.
//
// WHAT IS PORTED, AND FROM WHERE. Every semantic below is the library's, read
// at module github.com/jchv/go-webview2 v0.0.0-20260205173254-56598839c808:
//   - the window class + wndproc + create             webview.go:222-345
//   - Run()'s WM_APP drain / WM_QUIT exit             webview.go:351-379
//   - Dispatch / Terminate / Destroy                  webview.go:347-383
//   - Window / SetTitle / SetSize / Navigate / SetHtml
//     Init / Eval                                     webview.go:385-441
//   - the Bind stub + msgcb + callbinding RPC pipe    webview.go:131-220, :450-482
//
// ONE NAMED DEVIATION, and it is forced by this repo's own vet gate. The
// library's wndproc writes the min/max track sizes through
// `(*w32.MinMaxInfo)(unsafe.Pointer(lParam))` (webview.go:244-251), and
// `go vet -tags winlive ./cmd/wisp/` - the CI step that type-checks this package
// on Windows (.github/workflows/ci.yml:672) - flags exactly that uintptr-to-
// Pointer conversion as "possible misuse of unsafe.Pointer" (measured, reading
// in .scratch/wisp/probes/299/r1/logs/vet-unsafepointer.txt). The case below
// therefore returns 0 with no write, which is byte-for-byte what the library did
// in THIS product: its two fields are only ever set by SetSize with HintMin /
// HintMax, and the only SetSize caller in this tree (requestGeometryOnReshow)
// passes webview2.HintNone. SetSize still records the pair (its body is the
// library's), but nothing reads the two fields now.
//
// WHY A SYNTHETIC ORIGIN. The page's entry references its bundle relatively
// (frontend/vite.config.ts:23 `base: "./"`, so `./assets/index-*.js` and
// `./assets/*.css`), and a document's relative URLs resolve against its BASE
// URI. SetHtml produces an about:blank document (common.go:62), which no
// relative URL can be resolved against at all - measured 2026-10-13: the page
// never asks, so a filter alone intercepts nothing and the window stays white.
// So the entry is NAVIGATED TO at panelEntryURL, an origin THIS process owns,
// and every request under it is answered from the embed by the filter
// installAssetFilter registers. This is not a listener: no socket, no port, no
// byte leaves the process. D29's forbidden shape is the localhost HTTP server;
// this is the shape D29 names instead.

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unsafe"

	"github.com/CarlosShao/wisp/internal/panel"
	webview2 "github.com/jchv/go-webview2"
	"github.com/jchv/go-webview2/pkg/edge"
	"golang.org/x/sys/windows"
)

// The synthetic origin, spelled in ONE place. `.panel` is a name no DNS
// resolver can answer, so a request that ever escaped the filter would fail
// offline instead of reaching a real host. panelEntryURL is what serveEntry
// navigates to AND what the headless rulers ask the filter's lookup about
// (panel_pageover_33r10), so the navigation and the assertions cannot drift.
const (
	panelOrigin             = "https://wisp.panel"
	panelAssetFilterPattern = panelOrigin + "/*"
	panelEntryURL           = panelOrigin + "/" + panel.EntryFile
)

// entry points this file needs. Reused from panel_host_windows.go (same
// package): pnlShowWindow, pnlUpdateWindow, pnlSetFocus, pnlGetAncestor,
// pnlTranslateMsg, pnlDispatchMsgW.
var (
	pnlModKernel32 = windows.NewLazySystemDLL("kernel32.dll")
	pnlGetModule   = pnlModKernel32.NewProc("GetModuleHandleW")
	// ctlSetLastError clears this thread's last error before the one library
	// call whose success check reads it back - see respondWithEmbeddedAsset.
	ctlSetLastError = pnlModKernel32.NewProc("SetLastError")
)

// Constants of the ported shapes. The numeric values are the ones the library's
// own w32 package carries (internal/w32/w32.go:55-120); the pair spelled here so
// this file reads without that module open.
const (
	// ctlWindowClassName is this host's class. Deliberately NOT the library's
	// "webview": a later reader can tell which create made a window.
	ctlWindowClassName = "wisp-panel-host"

	ctlWsOverlappedWindow = 0x00CF0000
	ctlCwUseDefault       = 0x80000000
	ctlSmCxIcon           = 11
	ctlSmCyIcon           = 12
	ctlIdiApplication     = 32512

	ctlWmDestroy       = 0x0002
	ctlWmMove          = 0x0003
	ctlWmSize          = 0x0005
	ctlWmActivate      = 0x0006
	ctlWmClose         = 0x0010
	ctlWmQuit          = 0x0012
	ctlWmGetMinMaxInfo = 0x0024
	ctlWmNcLButtonDown = 0x00A1
	ctlWmMoving        = 0x0216
	ctlWmApp           = 0x8000

	// ctlGaRoot is GA_ROOT (2), the ancestor the library's Run() resolves before
	// IsDialogMessage (webview.go:371). GA_ROOTOWNER (3) is a different value
	// panel_host_windows.go's pnlGaRoot carries for its own, different call.
	ctlGaRoot = 2
	// ctlWaInactive is WA_INACTIVE (winuser.h), the activate state the library's
	// wndproc ignores.
	ctlWaInactive    = 0
	ctlWsThickFrame  = 0x00040000
	ctlWsMaximizeBox = 0x00010000
	// the SWP_ flags the library's SetSize posts with (w32.go:80-83).
	ctlSwpNoZOrder     = 0x0004
	ctlSwpNoActivate   = 0x0010
	ctlSwpNoMove       = 0x0002
	ctlSwpFrameChanged = 0x0020
)

// ctlGwlStyle is GWL_STYLE (-16). A variable, not a constant: the two
// GetWindowLongPtrW / SetWindowLongPtrW calls below pass it through uintptr, and
// a negative constant cannot be converted to uintptr at compile time. The
// library's own w32_64bit.go:7-14 converts an int variable the same way, so the
// value the call sees is the same sign-extended -16.
var ctlGwlStyle = -16

// ctlWndClassExW mirrors WNDCLASSEXW field for field (x64), shape-for-shape
// from the library's internal/w32.WndClassExW - which a package outside that
// module cannot import.
type ctlWndClassExW struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	ClsExtra      int32
	WndExtra      int32
	HInstance     windows.Handle
	HIcon         windows.Handle
	HCursor       windows.Handle
	HbrBackground windows.Handle
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       windows.Handle
}

// ctlPoint / ctlRect mirror POINT / RECT.
type (
	ctlPoint struct{ X, Y int32 }
	ctlRect  struct{ Left, Top, Right, Bottom int32 }
)

// user32 entry points, looked up once.
var (
	ctlUser32 = windows.NewLazySystemDLL("user32.dll")

	ctlRegisterClassExW   = ctlUser32.NewProc("RegisterClassExW")
	ctlCreateWindowExW    = ctlUser32.NewProc("CreateWindowExW")
	ctlDefWindowProcW     = ctlUser32.NewProc("DefWindowProcW")
	ctlDestroyWindow      = ctlUser32.NewProc("DestroyWindow")
	ctlPostMessageW       = ctlUser32.NewProc("PostMessageW")
	ctlPostThreadMessageW = ctlUser32.NewProc("PostThreadMessageW")
	ctlPostQuitMessage    = ctlUser32.NewProc("PostQuitMessage")
	ctlGetMessageW        = ctlUser32.NewProc("GetMessageW")
	ctlIsDialogMessage    = ctlUser32.NewProc("IsDialogMessage")
	ctlGetWindowLongPtrW  = ctlUser32.NewProc("GetWindowLongPtrW")
	ctlSetWindowLongPtrW  = ctlUser32.NewProc("SetWindowLongPtrW")
	ctlAdjustWindowRect   = ctlUser32.NewProc("AdjustWindowRect")
	ctlSetWindowPos       = ctlUser32.NewProc("SetWindowPos")
	ctlSetWindowTextW     = ctlUser32.NewProc("SetWindowTextW")
	ctlLoadImageW         = ctlUser32.NewProc("LoadImageW")
	ctlGetSystemMetrics   = ctlUser32.NewProc("GetSystemMetrics")
)

// --- the window-context registry ---------------------------------------------
//
// Win32 delivers messages to ONE class-wide callback, so the per-window
// behaviour lives in a map keyed by HWND. The library ships the same registry
// (webview.go:21-36); the one difference is that this one DROPS an entry on
// WM_DESTROY instead of keeping it, because a destroyed HWND's value can be
// reused by the next CreateWindowEx and a stale entry would then hijack a
// different window's messages. Post-destroy messages reach DefWindowProc either
// way - that is the library's own fallback for an unknown context too.

var (
	ctlWindowsMu sync.RWMutex
	ctlWindows   = map[uintptr]*panelControl{}
)

func ctlPutWindow(hwnd windows.HWND, c *panelControl) {
	ctlWindowsMu.Lock()
	ctlWindows[uintptr(hwnd)] = c
	ctlWindowsMu.Unlock()
}

func ctlLookupWindow(hwnd uintptr) (*panelControl, bool) {
	ctlWindowsMu.RLock()
	c, ok := ctlWindows[hwnd]
	ctlWindowsMu.RUnlock()
	return c, ok
}

func ctlDropWindow(hwnd windows.HWND) {
	ctlWindowsMu.Lock()
	delete(ctlWindows, uintptr(hwnd))
	ctlWindowsMu.Unlock()
}

// ctlWndProcHandle is the stable callback Windows holds for the class.
var ctlWndProcHandle = windows.NewCallback(ctlWndProc)

// ctlWndProc is the library's wndproc (webview.go:222-261), message for
// message, on this package's registry. Every arm returns 0 except the
// WM_NCLBUTTONDOWN and default arms, which hand back DefWindowProcW's answer -
// the library's own split.
func ctlWndProc(hwnd, msg, wparam, lparam uintptr) uintptr {
	c, ok := ctlLookupWindow(hwnd)
	if !ok {
		r, _, _ := ctlDefWindowProcW.Call(hwnd, msg, wparam, lparam)
		return r
	}
	log.Printf("MARK299 wndproc hwnd=%#x msg=%#x known=%v", hwnd, msg, ok)
	switch msg {
	case ctlWmMove, ctlWmMoving:
		// The WebView2 child tracks its parent's moves through this call.
		_ = c.browser.NotifyParentWindowPositionChanged()
	case ctlWmNcLButtonDown:
		pnlSetFocus.Call(hwnd)
		r, _, _ := ctlDefWindowProcW.Call(hwnd, msg, wparam, lparam)
		return r
	case ctlWmSize:
		c.browser.Resize()
	case ctlWmActivate:
		if wparam == ctlWaInactive {
			break
		}
		// Autofocus is false in every construction of this host (the wrapper
		// passed AutoFocus:false); the branch is the library's, kept whole.
		if c.autofocus {
			c.browser.Focus()
		}
	case ctlWmClose:
		// Destroy() only POSTS this message; the owning thread dispatches it
		// into here and the window goes away for real.
		_, _, _ = ctlDestroyWindow.Call(hwnd)
	case ctlWmDestroy:
		ctlDropWindow(windows.HWND(hwnd))
		// This is the library's Terminate, invoked where it lands: WM_DESTROY
		// is dispatched on the owning thread, and PostQuitMessage ends THAT
		// thread's GetMessageW loop (webview.go:242-243 -> :381-383).
		ctlPostQuitMessage.Call(0)
	case ctlWmGetMinMaxInfo:
		// See the file header: the library's write path is a vet-flagged
		// uintptr-to-Pointer conversion, and the two fields it would write are
		// only ever non-zero under HintMin/HintMax, which no caller in this
		// tree passes. Returning 0 without writing is what the library did in
		// this product.
		return 0
	default:
		r, _, _ := ctlDefWindowProcW.Call(hwnd, msg, wparam, lparam)
		return r
	}
	return 0
}

// --- the control ---------------------------------------------------------------

// panelControl is this package's own WebView2 host: the window, the browser
// control and the RPC pipe in one object. It satisfies webview2.WebView so
// panel_resident_windows.go's routing (post() -> Dispatch, the pump handover
// -> Run, stop() -> Terminate) and every existing ruler keep working against
// the same contract.
type panelControl struct {
	hwnd       windows.HWND
	mainthread uint32
	browser    *edge.Chromium
	autofocus  bool
	// maxsz / minsz are the library's fields for HintMax / HintMin. SetSize
	// records them (its body is the library's) but the WM_GETMINMAXINFO write
	// that read them is the one ported arm this repo's vet gate forbids - see
	// the file header.
	maxsz ctlPoint
	minsz ctlPoint

	m         sync.Mutex
	bindings  map[string]interface{}
	dispatchq []func()
}

var _ webview2.WebView = (*panelControl)(nil)

// newPanelControl creates the ONE panel window and embeds a WebView2 control in
// it. It must run on the thread that will pump the window's messages - the
// resident panel thread, or a locked thread in a test harness - which is the
// same rule the wrapper shipped with. Embed pumps this thread's queue until the
// environment AND the controller exist and installs the
// window.external.invoke shim the Bind stub rides (pkg/edge/chromium.go:72-114),
// so the control is live when this returns.
func newPanelControl(options webview2.WindowOptions, dataPath string) (*panelControl, error) {
	hinstance, _, _ := pnlGetModule.Call(0)

	// The control OBJECT (and with it the browser) exists before the window does,
	// which is the library's own order (NewWithOptions assigns w.browser at
	// webview.go:107, before its CreateWithOptions). It has to be: ShowWindow
	// below dispatches WM_SIZE/WM_MOVE synchronously, and this host's wndproc
	// answers those through the browser - a window shown with a nil browser
	// dereferences it on the spot (measured 2026-10-13: "panel host pump
	// panicked: runtime error: invalid memory address or nil pointer
	// dereference" 0.03s into bringUp, probe logs/run-ac2-live.txt).
	c := &panelControl{
		mainthread: windows.GetCurrentThreadId(),
		bindings:   map[string]interface{}{},
	}
	browser := edge.NewChromium()
	browser.MessageCallback = c.msgcb
	browser.DataPath = dataPath
	// The wrapper's one permission grant, kept (webview.go:105).
	browser.SetPermission(edge.CoreWebView2PermissionKindClipboardRead, edge.CoreWebView2PermissionStateAllow)
	c.browser = browser

	// Default application icon, the library's own loading shape (webview.go:273-282).
	icow, _, _ := ctlGetSystemMetrics.Call(ctlSmCxIcon)
	icoh, _, _ := ctlGetSystemMetrics.Call(ctlSmCyIcon)
	icon, _, _ := ctlLoadImageW.Call(hinstance, ctlIdiApplication, icow, icoh, 0)

	className, err := windows.UTF16PtrFromString(ctlWindowClassName)
	if err != nil {
		return nil, fmt.Errorf("panel host: window class name: %w", err)
	}
	wc := ctlWndClassExW{
		CbSize:        uint32(unsafe.Sizeof(ctlWndClassExW{})),
		HInstance:     windows.Handle(hinstance),
		LpszClassName: className,
		HIcon:         windows.Handle(icon),
		HIconSm:       windows.Handle(icon),
		LpfnWndProc:   ctlWndProcHandle,
	}
	// A recreate re-registers the class; Win32 answers "class already exists"
	// and the first registration (same callback, same name) keeps serving. The
	// library ignores this return value the same way (webview.go:293).
	_, _, _ = ctlRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))

	windowName, err := windows.UTF16PtrFromString(options.Title)
	if err != nil {
		return nil, fmt.Errorf("panel host: window title: %w", err)
	}

	// The library's 0 -> 640x480 fallback (webview.go:297-304), kept so a
	// caller that hands no geometry behaves exactly as it did through the
	// wrapper. The shipping host never takes it: windowOptions answers the
	// host's own constants when nobody sized it.
	windowWidth := options.Width
	if windowWidth == 0 {
		windowWidth = 640
	}
	windowHeight := options.Height
	if windowHeight == 0 {
		windowHeight = 480
	}

	log.Printf("MARK299 before CreateWindowExW")
	hwnd, _, callErr := ctlCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)),
		ctlWsOverlappedWindow,
		ctlCwUseDefault,
		ctlCwUseDefault,
		uintptr(windowWidth),
		uintptr(windowHeight),
		0,
		0,
		hinstance,
		0,
	)
	if hwnd == 0 {
		return nil, fmt.Errorf("panel host: CreateWindowExW returned 0 (class %s, %dx%d): %v",
			ctlWindowClassName, windowWidth, windowHeight, callErr)
	}
	c.hwnd = windows.HWND(hwnd)
	log.Printf("MARK299 after CreateWindowExW hwnd=%#x", hwnd)
	ctlPutWindow(c.hwnd, c)

	// The library shows the window as part of the create (webview.go:336-338);
	// PanelManager.Show then shows it again, which is idempotent.
	pnlShowWindow.Call(hwnd, pnlSwShow)
	log.Printf("MARK299 after ShowWindow")
	pnlUpdateWindow.Call(hwnd)
	log.Printf("MARK299 after UpdateWindow")
	pnlSetFocus.Call(hwnd)
	log.Printf("MARK299 after SetFocus")

	log.Printf("MARK299 before Embed")
	if !browser.Embed(hwnd) {
		ctlDropWindow(c.hwnd)
		_, _, _ = ctlDestroyWindow.Call(hwnd)
		return nil, errors.New("panel host: WebView2 Embed failed (runtime missing or blocked)")
	}
	log.Printf("MARK299 after Embed")
	browser.Resize()

	// Context menu + devtools follow the wrapper's Debug=false shape
	// (webview.go:113-126); a failure is fatal there and stays fatal here.
	settings, serr := browser.GetSettings()
	if serr != nil {
		log.Fatal(serr)
	}
	if err := settings.PutAreDefaultContextMenusEnabled(false); err != nil {
		log.Fatal(err)
	}
	if err := settings.PutAreDevToolsEnabled(false); err != nil {
		log.Fatal(err)
	}
	return c, nil
}

// --- the resource filter (ticket 299 AC#1) --------------------------------------

// installAssetFilter is the registration hop: after the control exists (the
// filter lives on the CoreWebView2 the controller-completion callback creates)
// and before any navigation that can request a resource. The library attaches
// the WebResourceRequested EVENT handler itself (chromium.go:186-231, the
// AddWebResourceRequested call at :211); what an application registers is the
// filter pattern (:292) plus the callback arm (:42, checked at :286).
func (c *panelControl) installAssetFilter(assets *panel.Assets) {
	c.browser.WebResourceRequestedCallback = func(request *edge.ICoreWebView2WebResourceRequest, args *edge.ICoreWebView2WebResourceRequestedEventArgs) {
		c.respondWithEmbeddedAsset(assets, request, args)
	}
	c.browser.AddWebResourceRequestedFilter(panelAssetFilterPattern, edge.COREWEBVIEW2_WEB_RESOURCE_CONTEXT_ALL)
}

// respondWithEmbeddedAsset answers one intercepted request out of the embed.
// Every request under panelOrigin is answered here; a miss is a 404 with an
// EMPTY body, never a fallback document - a silent index.html for every miss is
// the blank-panel-that-still-looks-alive shape assets.go:69-74 refuses.
func (c *panelControl) respondWithEmbeddedAsset(assets *panel.Assets, request *edge.ICoreWebView2WebResourceRequest, args *edge.ICoreWebView2WebResourceRequestedEventArgs) {
	uri, err := request.GetUri()
	if err != nil {
		log.Printf("panel host: resource request without a readable URI: %v", err)
		return
	}
	if !strings.HasPrefix(uri, panelOrigin) {
		// Unreachable through the filter pattern; left alone rather than
		// answered by a guess.
		return
	}
	data, contentType, status, reason := panelAssetLookup(assets, uri)
	env := c.browser.Environment()
	if env == nil {
		log.Printf("panel host: %s: no WebView2 environment to build a response", uri)
		return
	}
	headers := ""
	if status == 200 {
		// The content type comes from Assets.Resolve's own second return value
		// (assets.go:75); nothing here invents one.
		headers = "Content-Type: " + contentType
	}
	response, err := panelBuildResponse(env, data, status, reason, headers)
	if err != nil {
		log.Printf("panel host: %s: CreateWebResourceResponse (status %d, %d byte(s), %q): %v",
			uri, status, len(data), reason, err)
		return
	}
	if err := args.PutResponse(response); err != nil {
		log.Printf("panel host: %s: PutResponse: %v", uri, err)
	}
}

// panelBuildResponse creates one web-resource response, with a bounded retry of
// the library helper because that helper's success check is unreliable.
//
// The defect, measured rather than assumed (2026-10-13, this box, every
// AC#2 run): `ICoreWebView2Environment.CreateWebResourceResponse`
// (pkg/edge/corewebview2.go:169) reports success by comparing the raw
// GetLastError it reads after the COM call against ERROR_SUCCESS - and a
// SUCCESSFUL call can leave a stale nonzero value there. Every cold start's
// FIRST entry-document invocation read back errno 1008 (ERROR_NO_TOKEN, a Win32
// code no WebView2 HRESULT could be), the created-and-leaked response was
// dropped, the navigation fell through to the network, and Chromium re-issued
// it ~1.2 s later, where a later invocation of the same call served it. So the
// response object exists in that failure; only the wrapper's return value lies.
// Retrying is the one repair available from outside the module: a retry runs
// the same call again, and the stale value is not stable across attempts
// (measured: attempt 2 returns clean in every failing run so far). The bounded
// count and the loud line on a retry keep this from turning a real failure into
// a silent one.
func panelBuildResponse(env *edge.ICoreWebView2Environment, data []byte, status int, reason, headers string) (*edge.ICoreWebView2WebResourceResponse, error) {
	const attempts = 3
	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		// Cleared immediately before the call for the same reason: what the
		// helper reads back must be what THIS call left.
		ctlSetLastError.Call(0)
		response, err := env.CreateWebResourceResponse(data, status, reason, headers)
		if err == nil {
			if attempt > 1 {
				log.Printf("panel host: CreateWebResourceResponse needed attempt %d of %d (the stale-GetLastError defect named above)", attempt, attempts)
			}
			return response, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

// panelAssetLookup maps one intercepted URL under panelOrigin to the embedded
// bytes: (bytes, contentType, status, reason). It is a package function, not a
// closure, so a headless ruler can ask the filter's exact question - which is
// what keeps the window-free AC#13 ruler (panel_pageover_33r10) able to say
// "the entry URL serves the entry's own bytes" after the entry handover became
// a navigation.
func panelAssetLookup(assets *panel.Assets, uri string) ([]byte, string, int, string) {
	path := strings.TrimPrefix(uri, panelOrigin)
	if cut := strings.IndexAny(path, "?#"); cut >= 0 {
		path = path[:cut]
	}
	if assets == nil {
		return nil, "", 404, "Not Found"
	}
	data, contentType, err := assets.Resolve(path)
	if err != nil {
		return nil, "", 404, "Not Found"
	}
	return data, contentType, 200, "OK"
}

// --- webview2.WebView -----------------------------------------------------------

// Run is the library's pump (webview.go:351-379): WM_APP drains the dispatch
// queue, WM_QUIT returns, everything else goes through the dialog check and
// DispatchMessageW. It is the ONLY reader of that queue - which is why
// panel_resident_windows.go hands its thread over to this call instead of
// pumping by hand (a hand pump can never deliver a Go->page reply).
func (c *panelControl) Run() {
	var msg pnlMsg
	for {
		_, _, _ = ctlGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if msg.message == ctlWmApp {
			c.m.Lock()
			q := append([]func(){}, c.dispatchq...)
			c.dispatchq = []func(){}
			c.m.Unlock()
			for _, v := range q {
				v()
			}
		} else if msg.message == ctlWmQuit {
			return
		}
		r, _, _ := pnlGetAncestor.Call(uintptr(msg.hwnd), ctlGaRoot)
		r, _, _ = ctlIsDialogMessage.Call(r, uintptr(unsafe.Pointer(&msg)))
		if r != 0 {
			continue
		}
		pnlTranslateMsg.Call(uintptr(unsafe.Pointer(&msg)))
		pnlDispatchMsgW.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

// Terminate posts WM_QUIT to the CALLING thread's queue (webview.go:381-383):
// it ends Run() only when called on the thread that owns the window, which is
// why the resident leg routes it through the panel thread.
func (c *panelControl) Terminate() {
	ctlPostQuitMessage.Call(0)
}

// Dispatch appends f to the queue Run() drains and posts the wake-up. Same
// shape as the library's (webview.go:443-448): the append and the post are what
// makes "post() landed in the queue" and "the pump will read it" one event.
func (c *panelControl) Dispatch(f func()) {
	c.m.Lock()
	c.dispatchq = append(c.dispatchq, f)
	c.m.Unlock()
	if c.mainthread != 0 {
		_, _, _ = ctlPostThreadMessageW.Call(uintptr(c.mainthread), ctlWmApp, 0, 0)
	}
}

// Destroy posts WM_CLOSE; the owning thread's pump dispatches it and the window
// is really destroyed there (webview.go:347-349).
func (c *panelControl) Destroy() {
	_, _, _ = ctlPostMessageW.Call(uintptr(c.hwnd), ctlWmClose, 0, 0)
}

// Window is part of webview2.WebView and has NO faithful implementation here.
// Returning the handle as a pointer is a uintptr-to-unsafe.Pointer conversion,
// and `go vet -tags winlive ./cmd/wisp/` - the CI step that type-checks this
// package on Windows - flags exactly that shape (measured:
// probes/299/r1/logs/vet-unsafepointer.txt); the library only gets away with it
// because a dependency's body is never vetted. The host never goes through this
// door: it brings the window up, sizes it and tears it down through
// panelControl.hwnd, and before ticket 299-r1 the only caller of
// (*webview2.WebView).Window() in this tree was that same bringUp line this
// ticket replaced. A caller that still reaches the method gets nil - "0, no
// window" - a named, safe answer rather than a punned pointer.
func (c *panelControl) Window() unsafe.Pointer { return nil }

// SetTitle / SetSize / Navigate / SetHtml / Init / Eval are the library's
// methods, forwarded to the edge layer it forwards them to (webview.go:385-441).
func (c *panelControl) SetTitle(title string) {
	t, err := windows.UTF16FromString(title)
	if err != nil {
		t, _ = windows.UTF16FromString("")
	}
	_, _, _ = ctlSetWindowTextW.Call(uintptr(c.hwnd), uintptr(unsafe.Pointer(&t[0])))
}

// SetSize keeps the library's four-hint semantics (webview.go:405-433):
// HintFixed clears the thick-frame / maximize bits, HintMin/HintMax only record
// their pair (the arm that consumed them is the one ported-vet deviation above),
// and everything else converts the pair from CLIENT area to a window rect and
// asks the controller to resize.
func (c *panelControl) SetSize(width int, height int, hints webview2.Hint) {
	style, _, _ := ctlGetWindowLongPtrW.Call(uintptr(c.hwnd), uintptr(ctlGwlStyle))
	if hints == webview2.HintFixed {
		style &^= (ctlWsThickFrame | ctlWsMaximizeBox)
	} else {
		style |= (ctlWsThickFrame | ctlWsMaximizeBox)
	}
	_, _, _ = ctlSetWindowLongPtrW.Call(uintptr(c.hwnd), uintptr(ctlGwlStyle), style)

	if hints == webview2.HintMax {
		c.maxsz.X = int32(width)
		c.maxsz.Y = int32(height)
	} else if hints == webview2.HintMin {
		c.minsz.X = int32(width)
		c.minsz.Y = int32(height)
	} else {
		r := ctlRect{Left: 0, Top: 0, Right: int32(width), Bottom: int32(height)}
		_, _, _ = ctlAdjustWindowRect.Call(uintptr(unsafe.Pointer(&r)), ctlWsOverlappedWindow, 0)
		_, _, _ = ctlSetWindowPos.Call(
			uintptr(c.hwnd), 0,
			uintptr(r.Left), uintptr(r.Top),
			uintptr(r.Right-r.Left), uintptr(r.Bottom-r.Top),
			ctlSwpNoZOrder|ctlSwpNoActivate|ctlSwpNoMove|ctlSwpFrameChanged)
		c.browser.Resize()
	}
}

func (c *panelControl) Navigate(url string) { c.browser.Navigate(url) }
func (c *panelControl) SetHtml(html string) { c.browser.NavigateToString(html) }
func (c *panelControl) Init(js string)      { c.browser.Init(js) }
func (c *panelControl) Eval(js string)      { c.browser.Eval(js) }

// --- the RPC pipe, ported from the library --------------------------------------
//
// Bind / msgcb / callbinding are webview.go:131-220 and :450-482 with the
// receiver renamed: the page calls window.<name>(...), the stub turns the call
// into an RPC frame through window.external.invoke, MessageCallback hands the
// frame to msgcb, and the binding's answer travels back as an Eval of
// window._rpc[<id>].resolve(...) posted through Dispatch. panel_transport_35r2
// models this pipe call for call, so the two descriptions can be read against
// each other.

type ctlRpcMessage struct {
	ID     int               `json:"id"`
	Method string            `json:"method"`
	Params []json.RawMessage `json:"params"`
}

func ctlJsString(v interface{}) string { b, _ := json.Marshal(v); return string(b) }

func (c *panelControl) msgcb(msg string) {
	d := ctlRpcMessage{}
	if err := json.Unmarshal([]byte(msg), &d); err != nil {
		log.Printf("invalid RPC message: %v", err)
		return
	}

	id := strconv.Itoa(d.ID)
	if res, err := c.callbinding(d); err != nil {
		c.Dispatch(func() {
			c.Eval("window._rpc[" + id + "].reject(" + ctlJsString(err.Error()) + "); window._rpc[" + id + "] = undefined")
		})
	} else if b, err := json.Marshal(res); err != nil {
		c.Dispatch(func() {
			c.Eval("window._rpc[" + id + "].reject(" + ctlJsString(err.Error()) + "); window._rpc[" + id + "] = undefined")
		})
	} else {
		c.Dispatch(func() {
			c.Eval("window._rpc[" + id + "].resolve(" + string(b) + "); window._rpc[" + id + "] = undefined")
		})
	}
}

func (c *panelControl) callbinding(d ctlRpcMessage) (interface{}, error) {
	c.m.Lock()
	f, ok := c.bindings[d.Method]
	c.m.Unlock()
	if !ok {
		return nil, nil
	}

	v := reflect.ValueOf(f)
	isVariadic := v.Type().IsVariadic()
	numIn := v.Type().NumIn()
	if (isVariadic && len(d.Params) < numIn-1) || (!isVariadic && len(d.Params) != numIn) {
		return nil, errors.New("function arguments mismatch")
	}
	args := []reflect.Value{}
	for i := range d.Params {
		var arg reflect.Value
		if isVariadic && i >= numIn-1 {
			arg = reflect.New(v.Type().In(numIn - 1).Elem())
		} else {
			arg = reflect.New(v.Type().In(i))
		}
		if err := json.Unmarshal(d.Params[i], arg.Interface()); err != nil {
			return nil, err
		}
		args = append(args, arg.Elem())
	}

	errorType := reflect.TypeOf((*error)(nil)).Elem()
	res := v.Call(args)
	switch len(res) {
	case 0:
		return nil, nil

	case 1:
		if res[0].Type().Implements(errorType) {
			if res[0].Interface() != nil {
				return nil, res[0].Interface().(error)
			}
			return nil, nil
		}
		return res[0].Interface(), nil

	case 2:
		if !res[1].Type().Implements(errorType) {
			return nil, errors.New("second return value must be an error")
		}
		if res[1].Interface() == nil {
			return res[0].Interface(), nil
		}
		return res[0].Interface(), res[1].Interface().(error)

	default:
		return nil, errors.New("unexpected number of return values")
	}
}

func (c *panelControl) Bind(name string, f interface{}) error {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		return errors.New("only functions can be bound")
	}
	if n := v.Type().NumOut(); n > 2 {
		return errors.New("function may only return a value or a value+error")
	}
	c.m.Lock()
	c.bindings[name] = f
	c.m.Unlock()

	c.Init("(function() { var name = " + ctlJsString(name) + ";" + `
		var RPC = window._rpc = (window._rpc || {nextSeq: 1});
		window[name] = function() {
		  var seq = RPC.nextSeq++;
		  var promise = new Promise(function(resolve, reject) {
			RPC[seq] = {
			  resolve: resolve,
			  reject: reject,
			};
		  });
		  window.external.invoke(JSON.stringify({
			id: seq,
			method: name,
			params: Array.prototype.slice.call(arguments),
		  }));
		  return promise;
		}
	})()`)

	return nil
}
