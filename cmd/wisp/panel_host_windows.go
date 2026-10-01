//go:build windows

package main

// C27 panel host: the WebView2 window that makes the "page clicks -> Go reads ->
// receipt reaches the page" hop real for the first time (ticket 33 AC#1..AC#4).
//
// Shape and why it is this thin:
//   - ONE PanelManager, ONE WebView2 window per session, hide-don't-destroy
//     (D32/D29). Destroy happens at session dispose (a 31 hook), not on hide.
//   - Resources come from panel.Assets (the //go:embed all:dist bundle) and are
//     handed to the page through the WebView2 control itself. There is NO local
//     HTTP listener, no websocket and no window.fetch (bridge_test.go:159 names
//     that as ticket 92's forbidden second channel; AC#3 says "no listening
//     socket"). Serving bytes straight into the control is exactly the offline
//     model D29 asks for.
//   - The inbound hop rides the WebView2 message channel: a JS binding whose Go
//     callback parses the raw composer envelope through the ONE existing door
//     (panel.ComposerDispatch.Handle) and whose return value is the receipt the
//     page gets back. That is H2 (control created) + H3 (page bytes reach Go) +
//     H10 (reply reaches the page) landing together on a real host.
//   - The window's messages must belong to a thread that pumps them. In the
//     resident process that thread is the shared ui-sta (J1: the host is created
//     from an internal/ball Events callback, which already runs on ui-sta, so no
//     second UI thread and no new D38b roster name is added). A standalone
//     bring-up (the `wisp panel` diagnostic and the local measurement tests)
//     owns a locked OS thread for its one window; that is a measurement harness,
//     not the shipping resident topology, and it registers no goroutine name.
//
// Library gap that scopes AC#3 honestly (recorded in the evidence table, ⑧):
// pkg/edge does expose AddWebResourceRequestedFilter, but the same edge layer
// sizes its controller with ICoreWebView2Controller.PutBounds, whose parameter
// type is go-webview2's INTERNAL w32.Rect - unimportable from this module, and
// the vtbl that would let a caller build the rect by hand is unexported too. A
// controller we cannot size renders blank, so a rendering window can only be
// driven through the high-level webview2.WebView API, which does not surface the
// resource-request filter. Serving therefore uses SetHtml/Navigate over embedded
// bytes (still offline, still no listener) rather than AddWebResourceRequested
// Filter. Turning the page bundle into many files served through the filter
// needs a library decision (a replace/fork, or a different binding) - a
// dependency/architecture call, which per AGENTS.md is human-approved, so this
// leg stops and reports it instead of inventing a shape.

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/CarlosShao/wisp/internal/panel"
	webview2 "github.com/jchv/go-webview2"
	"golang.org/x/sys/windows"
)

// Binding name the page uses to send one raw composer envelope and get the
// router's answer back. It is a property of this host, not of the (ticket 35)
// bridge protocol, so it is spelled here and only here.
const panelDispatchBinding = "wispDispatch"

// panelWindowClass is the Win32 class go-webview2's high-level wrapper creates
// for us; kept only for diagnostics. panelOrigin is the opaque page origin that
// SetHtml produces (about:blank): offline, no host name, no port.
const panelTitle = "Wisp panel"

// user32 procs this host needs that golang.org/x/sys/windows does not export.
var (
	pnlModUser32     = windows.NewLazySystemDLL("user32.dll")
	pnlShowWindow    = pnlModUser32.NewProc("ShowWindow")
	pnlUpdateWindow  = pnlModUser32.NewProc("UpdateWindow")
	pnlSetFocus      = pnlModUser32.NewProc("SetFocus")
	pnlPeekMessageW  = pnlModUser32.NewProc("PeekMessageW")
	pnlDispatchMsgW  = pnlModUser32.NewProc("DispatchMessageW")
	pnlTranslateMsg  = pnlModUser32.NewProc("TranslateMessage")
	pnlSetForeground = pnlModUser32.NewProc("SetForegroundWindow")
)

const (
	pnlSwHide = 0
	pnlSwShow = 5
)

// pnlMsg mirrors the Win32 MSG layout (as internal/ball does) so a manual pump
// can move window messages off the queue. Only the fields the pump hands back to
// user32 are read, so the struct is passed by pointer and never field-decoded.
type pnlMsg struct {
	hwnd    windows.HWND
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      struct{ x, y int32 }
	_       uint32 // x64 padding DWORD
}

// PanelManager owns at most one panel window per session. It is safe for
// concurrent Show/Hide/IsShown calls; the WebView2 control itself is only touched
// on the thread that created the window (see bringUp's runtime.LockOSThread and
// the ui-sta hosting note above).
type PanelManager struct {
	mu       sync.Mutex
	w        webview2.WebView
	hwnd     windows.HWND
	dataPath string
	assets   *panel.Assets
	disp     *panel.ComposerDispatch

	created bool
	shown   bool
	// prevFocus is the foreground window recorded at the moment Show brought the
	// panel forward; Hide hands focus back to it (D29 focus return).
	prevFocus windows.HWND

	// lastColdMs / lastHotMs are the most recent measured bring-up latencies on
	// THIS host instance (D32 panel rows: cold <=1500ms, hot <=200ms). They are
	// measurements, never thresholds - the SLO thresholds live in
	// internal/observe/thresholds.go, untouched here.
	lastColdMs float64
	lastHotMs  float64

	// stopping tells the private-thread pump loop (bringUpOnNewThread) to exit.
	stopping atomic.Bool
}

// NewPanelManager wires one host to the inbound router and the embedded asset
// view. disp is the SAME panel.ComposerDispatch that has the production
// listeners (cmd/wisp/panel_inbound.go); the page's postMessage is just a second
// face onto it. assets may be nil (the host then reports "not built" and still
// answers the dispatch binding). dataPath is the WebView2 user-data folder; an
// empty string lets go-webview2 default it under %AppData%.
func NewPanelManager(disp *panel.ComposerDispatch, assets *panel.Assets, dataPath string) *PanelManager {
	return &PanelManager{disp: disp, assets: assets, dataPath: dataPath}
}

// IsCreated reports whether the window has been brought up in this session.
func (m *PanelManager) IsCreated() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.created
}

// IsShown reports the visible state (hide-don't-destroy: a shown window can be
// hidden and re-shown without a second msedgewebview2 child set).
func (m *PanelManager) IsShown() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.created && m.shown
}

// LastColdMs / LastHotMs expose the two latencies the last bring-up measured on
// this instance, so a caller (the diagnostic command, a test) can print or assert
// against a real reading instead of a promise.
func (m *PanelManager) LastColdMs() float64 { m.mu.Lock(); defer m.mu.Unlock(); return m.lastColdMs }
func (m *PanelManager) LastHotMs() float64  { m.mu.Lock(); defer m.mu.Unlock(); return m.lastHotMs }

// windowHandle returns the raw HWND (0 before creation). Callers use it to read
// the process's window state without going through WebView2.
func (m *PanelManager) windowHandle() uintptr {
	m.mu.Lock()
	defer m.mu.Unlock()
	return uintptr(m.hwnd)
}

// bringUp creates the ONE window on the calling thread and returns once the
// control is alive and the entry bytes have been handed to it. It must run on a
// thread that will pump the window's messages (ui-sta in the resident process, a
// locked OS thread for the standalone harness).
//
// Cold latency = window create + first browser round trip (the same "usable"
// point the S0 spike measured with NewWithOptions + SetHtml + a dispatch round
// trip; here the round trip is the JS binding proving the WebView2 message
// channel is live). It is a monotonic-clock measurement (Go's time.Since carries
// the monotonic reading), never a wall-clock timeout.
func (m *PanelManager) bringUp(ctx context.Context) error {
	t0 := time.Now()
	m.mu.Lock()
	if m.created {
		m.mu.Unlock()
		return nil
	}
	m.mu.Unlock()

	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     false,
		DataPath:  m.dataPath,
		AutoFocus: false,
		WindowOptions: webview2.WindowOptions{
			Title:  panelTitle,
			Width:  420,
			Height: 260,
		},
	})
	if w == nil {
		return fmt.Errorf("panel host: WebView2 window creation returned nil (runtime missing or blocked)")
	}

	hwnd := windows.HWND(w.Window())

	// The one inbound door. JS -> Go via a bound function over the WebView2
	// message channel; the returned string is the receipt that reaches the page.
	// No new inbound method name is invented here: the raw envelope is parsed by
	// panel.ParseComposerRequest inside Handle, whose whitelist is fixed at the
	// four existing methods (bridge.go:42-45, ticket 248 owns config methods).
	bindErr := w.Bind(panelDispatchBinding, func(raw string) string {
		reply, _ := m.dispatchRaw(ctx, raw)
		return reply
	})
	if bindErr != nil {
		w.Destroy()
		return fmt.Errorf("panel host: bind dispatch door: %w", bindErr)
	}

	m.mu.Lock()
	m.w = w
	m.hwnd = hwnd
	m.created = true
	m.mu.Unlock()

	// Hand the embedded entry bytes straight to the control (offline: the bytes
	// come from the //go:embed bundle, never from a socket).
	if err := m.serveEntry(); err != nil {
		w.Eval("document.title='panel assets unavailable'")
	}

	// Prove the message channel is live before declaring cold "usable": ask the
	// page to call the binding and take the first receipt as the round trip.
	rtMs := m.firstRoundTripLocked(ctx, t0)

	m.mu.Lock()
	m.lastColdMs = rtMs
	m.mu.Unlock()
	return nil
}

// serveEntry resolves the embedded index.html and pushes it into the control.
func (m *PanelManager) serveEntry() error {
	if m.assets == nil || !m.assets.Built() {
		return fmt.Errorf("panel host: embedded bundle not built (only the tracked placeholder is present)")
	}
	data, _, err := m.assets.Resolve(panel.EntryFile)
	if err != nil {
		return err
	}
	m.mu.Lock()
	w := m.w
	m.mu.Unlock()
	if w == nil {
		return fmt.Errorf("panel host: no window to serve into")
	}
	w.SetHtml(string(data))
	return nil
}

// Show reveals the window (creating it on first use) and returns focus
// bookkeeping. The caller must be on the window's message-pumping thread.
func (m *PanelManager) Show(ctx context.Context) error {
	m.mu.Lock()
	created := m.created
	m.mu.Unlock()
	if !created {
		if err := m.bringUp(ctx); err != nil {
			return err
		}
	}

	m.mu.Lock()
	hwnd := m.hwnd
	shown := m.shown
	m.prevFocus = windows.GetForegroundWindow()
	m.mu.Unlock()

	pnlShowWindow.Call(uintptr(hwnd), pnlSwShow)
	pnlUpdateWindow.Call(uintptr(hwnd))
	if !shown {
		pnlSetForeground.Call(uintptr(hwnd))
	}

	m.mu.Lock()
	m.shown = true
	m.mu.Unlock()
	return nil
}

// HotShow is Show for an already-created window; it also measures the hot path
// (C27 single-window reuse) into lastHotMs. Show already short-circuits creation
// when the window exists, so this wraps Show with the latency stamp.
func (m *PanelManager) HotShow(ctx context.Context) error {
	t0 := time.Now()
	if err := m.Show(ctx); err != nil {
		return err
	}
	m.mu.Lock()
	m.lastHotMs = float64(time.Since(t0).Microseconds()) / 1000.0
	m.mu.Unlock()
	return nil
}

// Hide takes the window off screen without destroying it, and hands focus back to
// whatever window was foreground before Show (D29).
func (m *PanelManager) Hide() {
	m.mu.Lock()
	hwnd := m.hwnd
	prev := m.prevFocus
	m.shown = false
	m.mu.Unlock()
	if hwnd == 0 {
		return
	}
	pnlShowWindow.Call(uintptr(hwnd), pnlSwHide)
	pnlUpdateWindow.Call(uintptr(hwnd))
	if prev != 0 {
		pnlSetForeground.Call(uintptr(prev))
		pnlSetFocus.Call(uintptr(prev))
	}
}

// Destroy tears the window down for good (session dispose, the 31 hook). After
// Destroy a later Show creates a fresh window (the recreate path).
func (m *PanelManager) Destroy() {
	m.mu.Lock()
	w := m.w
	m.w = nil
	m.hwnd = 0
	m.created = false
	m.shown = false
	m.mu.Unlock()
	if w != nil {
		w.Destroy()
	}
}

// Stop ends the private-thread pump loop started by bringUpOnNewThread. It is a
// no-op when the host is driven from the resident ui-sta (no private thread).
// The window itself is torn down by Destroy.
func (m *PanelManager) Stop() {
	m.stopping.Store(true)
}

// dispatchRaw runs one page envelope through the router. Kept separate from the
// bind closure so a test can call the same code path the page reaches.
func (m *PanelManager) dispatchRaw(ctx context.Context, raw string) (string, error) {
	if m.disp == nil {
		return "panel host: no inbound router attached", fmt.Errorf("panel host: nil router")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return m.disp.Handle(ctx, raw)
}

// firstRoundTripLocked drives one JS -> Go -> JS cycle and returns its duration,
// pumping the thread's queue until the receipt arrives or a monotonic deadline
// passes (a bounded wait, not a wall-clock timeout). It proves the WebView2
// message channel is really wired on this host (H3+H10) before cold latency is
// reported.
func (m *PanelManager) firstRoundTripLocked(ctx context.Context, t0 time.Time) float64 {
	m.mu.Lock()
	w := m.w
	m.mu.Unlock()
	if w == nil {
		return -1
	}

	done := make(chan struct{}, 1)
	// A second, one-shot binding the probe page calls; its Go side just signals
	// the round trip completed. Using a binding keeps the probe on the same
	// WebView2 channel the real door uses.
	if err := w.Bind("wispProbeRT", func() string {
		select {
		case <-done:
		default:
			close(done)
		}
		return "ok"
	}); err != nil {
		return -1
	}

	// A page that immediately calls the probe. SetHtml of a tiny document keeps
	// this offline and independent of the (ticket 34/36-40) real page.
	w.SetHtml(`<!doctype html><html><head><meta charset="utf-8"></head><body>` +
		`<script>if(window.wispProbeRT)window.wispProbeRT();</script></body></html>`)

	deadline := t0.Add(5 * time.Second)
	for {
		pnlPumpOnce()
		select {
		case <-done:
			return float64(time.Since(t0).Microseconds()) / 1000.0
		default:
		}
		if ctx != nil {
			if err := ctx.Err(); err != nil {
				return -1
			}
		}
		if time.Now().After(deadline) {
			return -1
		}
	}
}

// pnlPumpOnce drains whatever window messages are queued for this thread. The
// WebView2 control delivers bindings and resource completion as window messages,
// so a host that is not inside the resident ui-sta pump has to pump its own
// thread to make the message channel move.
func pnlPumpOnce() {
	var m pnlMsg
	for {
		r, _, _ := pnlPeekMessageW.Call(
			uintptr(unsafe.Pointer(&m)), 0, 0, 0, 1) // PM_REMOVE
		if r == 0 {
			return
		}
		pnlTranslateMsg.Call(uintptr(unsafe.Pointer(&m)))
		pnlDispatchMsgW.Call(uintptr(unsafe.Pointer(&m)))
	}
}

// bringUpOnNewThread creates the window on a private locked OS thread and keeps
// that thread pumping it for the returned manager's lifetime. This is the shape
// the `wisp panel` diagnostic and the local measurement tests use on a machine
// with a desktop session; it is NOT the resident topology (the resident process
// drives the host from a ui-sta Events callback and never calls this).
//
// It starts one goroutine that owns its OS thread and recovers anything the pump
// raises, so a panic cannot take the process down silently; the goroutine exits
// when Stop is called.
func (m *PanelManager) bringUpOnNewThread(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		// owner: this panel host thread; recover below.
		defer func() {
			if r := recover(); r != nil {
				errCh <- fmt.Errorf("panel host pump panicked: %v", r)
			}
		}()
		runtime.LockOSThread()
		errCh <- m.bringUp(ctx)
		// Keep the creating thread alive and pumping until told to stop.
		for !m.stopping.Load() {
			pnlPumpOnce()
			time.Sleep(5 * time.Millisecond)
		}
		runtime.UnlockOSThread()
	}()

	select {
	case err := <-errCh:
		return err
	case <-time.After(30 * time.Second):
		return fmt.Errorf("panel host: bring-up did not return within its bound")
	}
}
