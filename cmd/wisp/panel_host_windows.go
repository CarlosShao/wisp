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
//   - The window's messages belong to the thread that pumps them, so bringUp runs
//     on the resident panel's OWN dedicated STA thread
//     (cmd/wisp/panel_resident_windows.go, orchestrator ruling P1 at 10-01 13:12),
//     never on the ball's ui-sta. Why not ui-sta, now measured rather than asserted:
//     calling bring-up re-entrantly from a window callback of a thread that is
//     already pumping returns fine, but the outer pump's iteration freezes for the
//     whole nested pump and up to five tasks already queued on that thread get run
//     EARLY by the nested one (about 0.55s per open) - docs/evidence/s1/... the 33-p1
//     probe, R34/R35 with positive control R36. 33-p1 also corrected this file's
//     earlier claim that re-entry "panics on this box": the natural shape does not
//     panic; the nil dereference needs a FORCED WM_QUIT during create (R37) or an
//     MTA-initialised thread (R32/R39). A test harness still owns its own locked
//     thread for the local measurements; that harness is not the shipping topology.
//   - Because the resident thread hands the pump over to the library's own Run(),
//     a Go->page reply is actually delivered: webview.Dispatch only appends to the
//     library's private queue and posts a thread message, and the queue's only
//     reader is Run(). A host that pumps by hand instead never drains it, so the
//     page's awaited binding reply does not arrive (33-p1 §A, R25 vs R26).

// What is NOT a library gap (this paragraph was wrong before and is corrected
// here, ticket 33 AC#13's stop-and-report): AddWebResourceRequestedFilter is
// exported by pkg/edge, and the claim that a caller cannot size a controller
// because ICoreWebView2Controller.PutBounds takes go-webview2's INTERNAL w32.Rect
// does not block anything - (*edge.Chromium).Resize() (pkg/edge/
// chromium_amd64.go:12) is exported, sizes the controller from its own
// GetClientRect, and is what the high-level wrapper already calls after a
// successful Embed (webview.go:343). So "many files served through a resource
// request filter" is a product-shape decision, not a dependency boundary.
// What is still NOT proven, and is not claimed here: whether that path gives the
// panel the 420x260 size this host asks WindowOptions for. Today serving goes
// through SetHtml over embedded bytes (still offline, still no listener).

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
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
	pnlGetAncestor   = pnlModUser32.NewProc("GetAncestor")
)

// GetAncestor flags (winuser.h). GA_ROOT walks to the top-level window a handle
// belongs to, which is how Show tells "the window the user came from" apart from
// one of the panel's own child windows.
const (
	pnlGaRoot = 3
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

// PanelManager owns at most one panel window per session. The mutex makes its
// state readable from any goroutine; it does NOT make the WebView2 control
// thread-safe, and this file deliberately contains no runtime.LockOSThread of its
// own - the rule is that every method touching the control (bringUp/Show/Hide/
// Destroy) runs on the thread that created the window. In the resident process
// that thread is cmd/wisp/panel_resident_windows.go's dedicated panel thread,
// which is the only place that locks the OS thread, initialises COM as STA and
// hands the pump to the library's Run(). Earlier wording here pointed at a
// "bringUp's runtime.LockOSThread" that does not exist in this file; the lock is
// real now, and it lives in that one caller (33-r5, ticket 33 dispatch item 6).
type PanelManager struct {
	mu       sync.Mutex
	w        webview2.WebView
	hwnd     windows.HWND
	dataPath string
	assets   *panel.Assets
	disp     *panel.ComposerDispatch

	created bool
	shown   bool
	// prevFocus is the foreground window recorded BEFORE anything in Show can
	// move the foreground, and never the panel's own window (see
	// setPriorFocusLocked). Hide hands focus back to it and then clears it, so a
	// second Hide cannot restore a stale handle (D29 focus return, ticket 33 AC#4).
	prevFocus windows.HWND

	// lastRestoreTo / lastRestoreSetForeground / lastRestoreSetFocus are what the
	// most recent Hide tried and what Win32 answered (0 means "no restore was
	// attempted"). They exist so AC#4's readings are per-run facts instead of
	// inference; see Hide.
	lastRestoreTo            windows.HWND
	lastRestoreSetForeground uintptr
	lastRestoreSetFocus      uintptr

	// lastColdMs / lastHotMs are the most recent measured bring-up latencies on
	// THIS host instance (D32 panel rows: cold <=1500ms, hot <=200ms). They are
	// measurements, never thresholds - the SLO thresholds live in
	// internal/observe/thresholds.go, untouched here.
	lastColdMs float64
	lastHotMs  float64
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

// currentWindow returns the live control handle (nil before creation and after
// Destroy). The resident panel thread publishes exactly this value to decide where
// a request goes: while it is nil the thread's own channel is the route, once it is
// set the library's Dispatch queue (which only Run() drains) is. Reading it under
// the manager lock is what keeps "the window exists" and "the pump can be reached"
// from being two different instants.
func (m *PanelManager) currentWindow() webview2.WebView {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.w
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

	// The creating thread must not carry a WM_QUIT left by a prior window owner,
	// or the library's Embed pump (pkg/edge chromium.go:96-111) dequeues that quit
	// and breaks its GetMessageW loop with inited still 0 and e.webview nil, and
	// Init (chromium.go:130-136, no guard on the nil) panics - the panic behind
	// ticket 33's order-dependent AC#13 red. go-webview2 leaves the quit for free:
	// its wndproc calls w.Terminate() on WM_DESTROY (webview.go:242-243) and
	// Terminate is a bare PostQuitMessage (webview.go:381-383) onto whichever thread
	// pumped the close. A thread that has hosted a now-destroyed panel window can
	// therefore carry a pending quit into the next bring-up. The resident thread
	// never reuses one (loop locks for life and the thread is destroyed on exit, so
	// its quit dies with it) - but a test harness that returns a window thread to
	// the Go pool does, and that is the trigger ticket 33's full-package red showed.
	// This drain is the product-side half of the invariant "a thread that has pumped
	// a panel window is quit-clean before it creates another"; it is deliberately
	// NARROW (only WM_QUIT, only removed, never dispatched) - dispatching a stale
	// close here would tear down a prior owner's window and re-post a fresh quit,
	// and 33-r6 measured that purging the whole queue without dispatching leaves
	// zombie WebView2 controllers that hang the NEXT create. The companion fix - the
	// test harness keeping its window thread locked so it is destroyed, not pooled -
	// is what stops WM_CLOSE-shaped poison from ever reaching the pool.
	drainStaleQuitBeforeCreate()

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

	// AC#13's order, and it is the whole fix: prove the message channel first,
	// hand the page over LAST. firstRoundTripLocked shows a document of its own,
	// so running it after serveEntry meant every cold start finished on the probe
	// page instead of the panel (ticket 33 AC#13, 33-v1 §AC#3 "供给那半"). The
	// probe stays - it is where cold "usable" is decided - it just no longer gets
	// the last word about what the user sees.
	rtMs := m.firstRoundTripLocked(ctx, t0)

	if err := m.serveEntry(); err != nil {
		m.serveNotBuiltNoticeLocked()
	}

	m.mu.Lock()
	m.lastColdMs = rtMs
	m.mu.Unlock()
	return nil
}

// serveNotBuiltNoticeLocked replaces the page with an explicit offline notice when
// the embed carries no bundle. It is a SetHtml of a document this host writes, so
// the user is never left looking at the round-trip probe page, and it still opens
// no socket and loads nothing over the network.
func (m *PanelManager) serveNotBuiltNoticeLocked() {
	m.mu.Lock()
	w := m.w
	m.mu.Unlock()
	if w == nil {
		return
	}
	w.SetHtml(`<!doctype html><html><head><meta charset="utf-8"><title>panel assets unavailable</title>` +
		`</head><body><p>panel assets unavailable: the embedded bundle is not built.</p></body></html>`)
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
//
// AC#4's whole hop depends on the ORDER of two statements here: the prior
// foreground window is sampled BEFORE bringUp runs, because bringing the control
// up already takes the foreground (measured by 33-r4: in every run the value read
// after creation was the panel's own HWND, so Hide handed focus back to the window
// it had just hidden). The sampled value goes through setPriorFocusLocked, which
// refuses to record the panel itself or one of its own child windows.
func (m *PanelManager) Show(ctx context.Context) error {
	m.mu.Lock()
	created := m.created
	m.mu.Unlock()

	prior := windows.GetForegroundWindow()
	if !created {
		if err := m.bringUp(ctx); err != nil {
			return err
		}
	}

	m.mu.Lock()
	hwnd := m.hwnd
	shown := m.shown
	m.setPriorFocusLocked(prior)
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

// setPriorFocusLocked records the caller's foreground sample for Hide to hand
// focus back to, and refuses samples that cannot be a real "where the user came
// from": the panel's own top-level window, or any window whose root IS the panel.
// Refusing keeps prevFocus at the last honest value instead of overwriting it with
// the panel, which is what made AC#4 red (33-v1 §A#27, 33-r4 §① 格 2).
//
// The caller must hold m.mu.
func (m *PanelManager) setPriorFocusLocked(prior windows.HWND) {
	if prior == 0 {
		return
	}
	if m.hwnd != 0 && sameRootWindow(prior, m.hwnd) {
		return
	}
	m.prevFocus = prior
}

// sameRootWindow reports whether a and b are the same top-level window or one is a
// child of the other, by comparing what GetAncestor(GA_ROOT) resolves them to.
// An unreadable handle resolves to 0 and therefore to "not the same window": a
// failed lookup must not silently drop a legitimate focus-restore target.
func sameRootWindow(a, b windows.HWND) bool {
	if a == b {
		return true
	}
	ra, _, _ := pnlGetAncestor.Call(uintptr(a), pnlGaRoot)
	rb, _, _ := pnlGetAncestor.Call(uintptr(b), pnlGaRoot)
	if ra == 0 || rb == 0 {
		return false
	}
	return ra == rb
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
// the window Show recorded as the user's prior foreground (D29).
//
// The recorded window is NOT cleared afterwards, and that is a measured decision,
// not an oversight: with a clear, a Show whose only foreground sample is the panel
// itself (which setPriorFocusLocked correctly refuses) leaves the next Hide with no
// target at all, and the reading on this box is "focus never went back". Keeping it
// means the honest prior survives one refused sample, and a stale or already-closed
// handle is a SetForegroundWindow that returns 0 - recorded below, never retried.
//
// The two Win32 return codes are stored, not acted on: a denial is Windows'
// foreground lock (a process that owns no foreground rights cannot move the
// foreground), and this host does not sleep-and-retry it, because that would turn an
// unstable product path into a green test (ticket 33 ruling, 10-01 12:12, item 2).
// lastRestoreTo / lastRestoreSetForeground / lastRestoreSetFocus carry the per-run
// handles instead.
func (m *PanelManager) Hide() {
	m.mu.Lock()
	hwnd := m.hwnd
	prev := m.prevFocus
	m.shown = false
	m.lastRestoreTo = 0
	m.lastRestoreSetForeground = 0
	m.lastRestoreSetFocus = 0
	m.mu.Unlock()
	if hwnd == 0 {
		return
	}
	pnlShowWindow.Call(uintptr(hwnd), pnlSwHide)
	pnlUpdateWindow.Call(uintptr(hwnd))
	if prev == 0 {
		return
	}
	rf, _, _ := pnlSetForeground.Call(uintptr(prev))
	rc, _, _ := pnlSetFocus.Call(uintptr(prev))
	m.mu.Lock()
	m.lastRestoreTo = prev
	m.lastRestoreSetForeground = rf
	m.lastRestoreSetFocus = rc
	m.mu.Unlock()
}

// terminateOnThisThread asks the window's message loop to end. It must be called
// ON the thread that owns the window: the underlying library implements Terminate
// as PostQuitMessage, which posts WM_QUIT to the CALLING thread's queue. Calling it
// from anywhere else ends the wrong loop and leaves the panel pumping forever.
func (m *PanelManager) terminateOnThisThread() {
	m.mu.Lock()
	w := m.w
	m.mu.Unlock()
	if w != nil {
		w.Terminate()
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

// firstRoundTripLocked drives one JS -> Go cycle and returns its duration, pumping
// the thread's queue until the Go side sees the probe binding or a monotonic
// deadline passes (a bounded wait, not a wall-clock timeout).
//
// What this proves and what it does NOT (orchestrator ruling P2, 10-01 13:12):
// its range is "the page reached Go" (H3). The `done` channel closes inside the Go
// body of the binding, which happens before any reply has to travel back, so this
// function is NOT evidence for AC#14 - 33-p1 §A measured a run where ECHO_CALLS=3
// (page to Go, arrived) and the awaited replies all stayed unresolved in the same
// process. AC#14's receipt hop has its own ruler, whose assertion is made of what
// the PAGE reports back, and a second one for Go-side Eval push; the two are
// separate dimensions and stay separate tests.
//
// The probe page it shows is transient by design and, since AC#13 was fixed, is no
// longer the last word: bringUp runs this BEFORE serveEntry, so the document the
// user ends on is the embedded entry.
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

const (
	// pnlWmQuit is WM_QUIT. The PeekMessage range filter is inclusive, so
	// [pnlWmQuit, pnlWmQuit] matches only the quit message - a stale close or any
	// other leftover is left untouched (dispatching them here is the zombie hazard
	// 33-r6 measured).
	pnlWmQuit = 0x12
	// pnlPmRemove is PeekMessage's PM_REMOVE flag; without it a message would be
	// reported again on the next call and the drain would loop forever.
	pnlPmRemove = 1
	// pnlQuitDrainCap bounds the loop. More than this many WM_QUIT on one thread
	// means several prior owners piled up; stop rather than trust the count, and say
	// so out loud so a failed drain is a loud reading, not a silent green.
	pnlQuitDrainCap = 64
)

// drainStaleQuitBeforeCreate removes every WM_QUIT still queued on the CALLING
// thread so the library's Embed pump cannot dequeue it before the control exists.
// It returns how many it removed. It is narrow by construction: only WM_QUIT (the
// PeekMessage range filter is [pnlWmQuit, pnlWmQuit]), only this thread, only
// before a window is created, and it REMOVES without DISPATCHING (so a stale close
// cannot be pumped into a fresh quit here - the harness keeps its window thread
// locked and destroyed so a close never reaches a pooled thread at all). See the
// comment at the top of bringUp for why the quit it prevents is a real hazard.
func drainStaleQuitBeforeCreate() int {
	removed := 0
	for removed < pnlQuitDrainCap {
		var m pnlMsg
		// filter [WM_QUIT..WM_QUIT] matches only the quit; PM_REMOVE drops it.
		r, _, _ := pnlPeekMessageW.Call(
			uintptr(unsafe.Pointer(&m)), 0, pnlWmQuit, pnlWmQuit, pnlPmRemove)
		if r == 0 {
			return removed
		}
		removed++
	}
	slog.Warn("panel host: pre-create quit-drain hit its cap before the queue was clean",
		"cap", pnlQuitDrainCap, "removed", removed)
	return removed
}
