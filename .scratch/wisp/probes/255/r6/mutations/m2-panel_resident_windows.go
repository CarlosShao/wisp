//go:build windows

package main

// The resident panel thread: the ONE place in this process that owns the panel
// window's OS thread (ticket 33, orchestrator ruling P1 at 10-01 13:12).
//
// Why a dedicated thread and not the ball's ui-sta. Ruling J1 in this ticket first
// picked "post it to the existing ui-sta", and 33-r1 stopped because it measured
// something bad happening there. The measurement has now been repeated with a
// counter-control (probe 33-p1, R34/R35 plus positive control R36): calling
// bring-up from inside a window callback of a thread that is already pumping does
// NOT panic - it returns in about half a second, and while it is inside the
// library's nested pump the OUTER pump's iteration counter freezes, so up to five
// tasks that were already queued on that thread get executed early by the nested
// pump. That is order reversal on the ball's own gesture queue, which is worse to
// ship than a crash. The branch that does panic is reachable but needs a forced
// WM_QUIT during create (R37) or an MTA apartment (R32/R39). So: a thread of its
// own, and the ball's pump is never re-entered.
//
// Why the library's Run(). webview.Dispatch only appends a closure to the
// library's private queue and posts a thread message, and Run() is that queue's
// ONLY reader, so a host that pumps by hand can never deliver a Go -> page reply:
// 33-p1 §A measured the page's awaited binding replies staying unresolved with the
// hand pump and arriving in ~107ms with Run(). Run() occupies the thread, so this
// file also owns the way out (Terminate -> WM_QUIT -> Run returns -> Destroy on
// the same thread -> the handle's goroutine ends), which is what makes the panel
// thread's clean exit an observable fact instead of a hope.
//
// Why COM is initialised here explicitly, and why a failure is fatal FOR THIS
// THREAD ONLY. 33-p1 §B-4 measured three shapes on one matrix: STA (0x2) creates
// fine, an un-initialised thread also creates fine, and a thread initialised as
// MTA (0x0) panics inside the library (pkg/edge/chromium.go:171 does not catch the
// failed HRESULT, :175 dereferences the nil environment). "It worked without
// initialising" is the shape this repo's host and its test harness shipped until
// now, so this thread asks for 0x2 by name and refuses to create a window if
// Win32 answers anything other than S_OK / S_FALSE. The cause of the MTA failure is
// still undetermined (the HRESULT needs an unexported handler to read); the
// phenomenon is what is encoded here.
//
// What this file does NOT do. It never calls Shutdown or os.Exit (A493), it starts
// no second window, it grants nothing: the only door it exposes to the page is the
// one cmd/wisp/panel_host_windows.go already binds (panel.ComposerDispatch.Handle
// behind window.wispDispatch), whose widening branches are fail-closed, and it adds
// no approval channel of its own. Its teardown is a defer in runResident, the same
// shape ticket 228 chose for the ball, so the frozen D38(e) ten-step order and its
// closed hook roster are untouched.
//
// What it DOES do since 33-r9, and why that is not a new shape: a show request that
// bringUp refuses by name (the thread owes somebody an undispatched WM_CLOSE) now
// ends this thread instead of leaving it up, taking requests and refusing every one
// of them forever (ticket 33, 33-v2 问①: "拒绝没有出口"). It rides the fatal path this
// file already had for the CoInitializeEx guard, so the exit adds no thread, no pump
// and no second window, and it removes nothing from anybody's message queue - which
// is the thing 33-r9's fifth-shape reading says no recovery here may do: the close it
// refuses on can belong to a window that is still alive and still owned (readings in
// cmd/wisp/panel_host_windows.go's bringUp comment).

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/panel"
	webview2 "github.com/jchv/go-webview2"
	wv2 "github.com/jchv/go-webview2"
	"golang.org/x/sys/windows"
)

const (
	// panelSTAName is this thread's registry name. It is deliberately NOT one of
	// observe.ResidentNames' six frozen resident names: adding a name there is a
	// contract change (D38b), so a leg that needs one asks instead of renaming. An
	// unlisted name is classified CategoryUnknown, which logs and does not fail.
	panelSTAName  = "panel-sta"
	panelSTAOwner = "panel host (ticket 33)"

	// coinitApartmentThreaded is COINIT_APARTMENTTHREADED. internal/ball spells the
	// same value (win32_windows.go) for its own ui-sta; that constant is package
	// private, so this file spells it too rather than reaching into internal/ball.
	coinitApartmentThreaded = 0x2
	// sFalse is what CoInitializeEx returns when the thread already has an
	// apartment; S_OK (0) and S_FALSE (1) both mean "this thread is STA now".
	sFalse = 1

	// panelThreadExitBudget bounds how long stop() waits for the thread to be
	// seen gone. It is a monotonic deadline (time.After), not a wall-clock
	// difference, and exceeding it is reported, never retried.
	panelThreadExitBudget = 5 * time.Second
)

var (
	pnlModOle32       = windows.NewLazySystemDLL("ole32.dll")
	pnlCoInitializeEx = pnlModOle32.NewProc("CoInitializeEx")
)

// residentPanel is the assembly root's handle on the panel thread. Every method is
// safe to call from any goroutine, including the ball's ui-sta callbacks, which is
// the contract those callbacks are written to ("quick and non-blocking"): show,
// hide and dispose hand a closure over and return.
type residentPanel struct {
	mgr *PanelManager
	reg *observe.Registry

	// mu guards wRef and the "should still accept the channel route" decision
	// together, so a request can never be left in a channel nobody reads and can
	// never be Dispatched into a queue nobody drains.
	mu      sync.Mutex
	wRef    webview2.WebView
	tasks   chan func()
	startUp error

	entered  chan struct{} // closed just before the library pump takes over
	finished chan struct{} // closed when the thread's goroutine returns
	stopNow  chan struct{} // closed by stop() while the pump is not live yet

	thread     *observe.Handle
	closeOnce  sync.Once
	toggles    atomic.Int64
	shows      atomic.Int64
	hides      atomic.Int64
	disposals  atomic.Int64
	failedPost atomic.Int64

	// stopRequested is what makes stop() effective no matter where in the thread's
	// life it lands: the pump cannot be entered and then ignored, because a thread
	// that reaches Run() after a stop was already asked would otherwise pump
	// forever with nobody left to send WM_QUIT.
	stopRequested atomic.Bool
}

// startResidentPanel builds the thread but NOT the window: nothing is created
// until the user asks (a ball gesture, a tray item, a test). That is what keeps a
// machine with no desktop, or a machine with no WebView2 runtime, booting exactly
// as it does today - the bring-up verdict lands on this thread and is logged, and
// runResident's ball leg is unaffected either way.
func startResidentPanel(reg *observe.Registry, mgr *PanelManager) *residentPanel {
	rp := &residentPanel{
		mgr:      mgr,
		reg:      reg,
		tasks:    make(chan func(), 16),
		entered:  make(chan struct{}),
		finished: make(chan struct{}),
		stopNow:  make(chan struct{}),
	}
	root := observe.NewRoot("panel-host")
	rp.thread = reg.Spawn(panelSTAName, panelSTAOwner, root, rp.loop)
	return rp
}

// residentPanelActor names whoever drives the resident panel's writes, so a mode
// switch asked for by the page can never be read as the CLI seam's (the audit line
// carries this string verbatim).
const residentPanelActor = "resident-panel"

const residentPanelHotReloadNote = "panel host (resident): this leg does not tick config.toml either; " +
	"the reload tick lives in `wisp run` (config_reload.go), and this leg has no approval card, so a loosening " +
	"it could read would have nowhere to be confirmed"

// residentPanelGeometryNote is the third sentence this leg owes about config.toml,
// and it is not the same claim as residentPanelHotReloadNote above. 票 255 AC#4
// wired [panel] width/height into the window, and the way it could be wired in
// THIS process is not through a Manager: the panel host is built at
// cmd/wisp/resident_windows.go:142, long before the run.go pipeline that owns
// config.NewManager exists, so what the assembly root hands the host is a closure
// that re-reads config.toml at every window creation (panelGeometrySource below).
// What still does not move, and this line says so out loud instead of letting the
// operator infer it: the tick, and with it every OTHER hot section of this
// process, keeps needing a restart.
const residentPanelGeometryNote = "config: PANEL-GEOMETRY state=per-create-and-per-reshow reads=[panel] width/height " +
	"detail=\"面板宿主每次建窗现读一次 [panel] width/height：关窗再开即跟上新值；已建好的窗口在下一次重新显示（Show）时" +
	"也会把此刻解析出的那一对数发给它一次（票 255-r1，走库的 SetSize，客户区语义，与建窗那份外框语义不是同一个宽度）。本腿不轮询 config.toml，也没有任何东西在文件被保存那一刻去按这一下，所以面板尺寸要等下一次显示请求才跟上，其余热加载段仍要重启进程才生效。\""

// panelGeometrySource is [panel] width/height, re-read from disk every time the
// panel host is about to create a window. Built by the assembly root, handed down
// as a function value: cmd/wisp/panel_host_windows.go imports nothing that parses
// config.toml (票 255 AC#4's ⛔ no panel->config edge, ⛔ no host reading the disk
// itself), and the host's own constants stay what they were for a host nobody
// sized.
//
// WHY a fresh config.LoadFile and not a Manager: the only Manager-shaped reader in
// this repo is internal/config's Manager, and this leg owns none (its header
// comment above, and `grep -n 'config\.' cmd/wisp/resident_windows.go` answers zero
// production hits). Standing one up here would make residentPanelHotReloadNote's
// 「this leg does not tick config.toml」 false, which is another leg's ruling. The
// one-shot load is this repository's existing shape for "read what is on disk in
// this process" - cmd/wisp/models.go:184 does exactly that for `wisp models`.
//
// The failure branch is loud and stays at the constants: a host whose config could
// not be parsed says so rather than silently answering someone else's number.
func panelGeometrySource(dataDir string) func() (width, height int) {
	cfgPath := filepath.Join(dataDir, configFileName)
	return func() (width, height int) {
		cfg, _, err := config.LoadFile(cfgPath, nil)
		if err != nil {
			slog.Warn("panel host: [panel] geometry source unreadable, sizing at the host's own default",
				"path", cfgPath, "err", err, "default", fmt.Sprintf("%dx%d", panelWidthPx, panelHeightPx))
			return 0, 0
		}
		return cfg.Panel.Width, cfg.Panel.Height
	}
}

func newResidentComposerDispatch(dataDir string, auditf panel.AuditFunc) (*panel.ComposerDispatch, error) {
	// Said once per start, same posture the CLI seam takes: name the tier that was
	// not taken rather than leave the operator to infer it.
	auditf("%s", residentPanelHotReloadNote)
	return newComposerDispatchChain(dataDir, auditf, residentPanelActor)
}

// newResidentPanelManager is the assembly the resident leg uses. The chain is the
// production one (the same five constructors cmd/wisp/panel_inbound.go assembles
// for the CLI seam and run.go:648 builds for the run seam): config Manager ->
// perm.Store -> panel.ModeWriteHandler -> panel.ComposerDispatch. Confirm stays
// nil on purpose: this leg has no card of its own, so a widening request is
// refused before Store.Set is ever reached, and the panel side never produces an
// "allow" (internal/agent/approval/ui.go: PanelAPI has no Allow method).
//
// The actor is named for what actually drives it, so the audit line cannot be read
// as the CLI seam's.
func newResidentPanelManager(dataDir string) (*PanelManager, error) {
	auditf := func(format string, args ...any) {
		line := fmt.Sprintf(format, args...)
		slog.Info("panel host: " + line)
		fmt.Printf("wisp: %s\n", line)
	}
	disp, err := newResidentComposerDispatch(dataDir, auditf)
	if err != nil {
		return nil, err
	}
	assets, aerr := panel.BuiltinAssets()
	if aerr != nil {
		// Not fatal and not hidden: nil assets make the host report "bundle not
		// built" and the dispatch door still answers, which is what
		// NewPanelManager's own contract says.
		slog.Warn("panel host: embedded bundle unavailable, serving nothing", "err", aerr)
		assets = nil
	}
	// The WebView2 user-data folder lives under the data root the boot already
	// resolved; this leg parses no data root of its own (ticket 128's leg table).
	dataPath := filepath.Join(dataDir, "panel-webview2")
	// 票 255 AC#4: the one thing this leg now reads out of config.toml, and it reads
	// it per window creation rather than holding it. Said once per start, in the
	// same posture as the hot-reload note above.
	auditf("%s", residentPanelGeometryNote)
	return NewPanelManager(disp, assets, dataPath, withGeometrySource(panelGeometrySource(dataDir))), nil
}

// loop is the panel thread. Registry.Spawn owns the recover and the panic sink,
// so this is not a bare goroutine (ban #1): it has an owner, it is booked in the
// registry, and it locks its OS thread for its whole life.
func (rp *residentPanel) loop(ctx context.Context) {
	defer close(rp.finished)

	runtime.LockOSThread()
	if r, _, _ := pnlCoInitializeEx.Call(0, coinitApartmentThreaded); r != 0 && r != uintptr(sFalse) {
		rp.setStartUp(fmt.Errorf("panel thread: CoInitializeEx(COINIT_APARTMENTTHREADED) returned %#x - refusing to create a WebView2 window on an apartment that is not STA", r))
		slog.Error("panel thread: refusing to bring up the panel window", "err", rp.startUpErr())
		return
	}

	for {
		rp.drainTasks()
		if rp.isCreated() || rp.startUpErr() != nil {
			break
		}
		select {
		case <-ctx.Done():
			rp.teardown("context cancelled before the panel window existed")
			return
		case <-rp.stopNow:
			rp.teardown("stop requested before the panel window existed")
			return
		case fn := <-rp.tasks:
			fn()
		}
	}

	w := rp.handOverPump()
	if w == nil {
		rp.teardown("no window to pump")
		return
	}
	if rp.stopRequested.Load() {
		// A stop that arrived while the thread was still inside bring-up: the pump
		// is about to start, so ask it to end before it begins. Run() then returns
		// on the first GetMessageW, and the teardown below is the same one.
		w.Terminate()
	}
	// The library pump. It returns on WM_QUIT, which Terminate posts - the exit
	// this thread is designed around (see the header).
	w.Run()
	rp.teardown("the library pump returned")
}

// drainTasks runs everything queued on the channel. Only the thread itself calls
// it, so the closures run on the window's thread - which is the entire point of
// the post() routing below.
func (rp *residentPanel) drainTasks() {
	for {
		select {
		case fn := <-rp.tasks:
			fn()
		default:
			return
		}
	}
}

// handOverPump publishes the control (so later requests route through
// webview.Dispatch, which Run() drains) and takes one final channel drain under
// the same lock. That ordering is what makes the routing airtight: post() either
// sees wRef empty and lands in the channel this drain reads, or sees wRef set and
// lands in the queue the pump reads.
func (rp *residentPanel) handOverPump() webview2.WebView {
	rp.mu.Lock()
	w := rp.mgr.currentWindow()
	rp.wRef = w
	var pending []func()
drain:
	for {
		select {
		case fn := <-rp.tasks:
			pending = append(pending, fn)
		default:
			break drain
		}
	}
	rp.mu.Unlock()
	for _, fn := range pending {
		fn()
	}
	if w != nil {
		close(rp.entered)
	}
	return w
}

// post hands fn to the panel thread without waiting for it. False means "this
// thread will not run your request" (its queue is full, or it is already going
// away), which callers report instead of pretending the window moved.
func (rp *residentPanel) post(fn func()) bool {
	if rp == nil {
		return false
	}
	if rp.startUpErr() != nil || rp.isFinished() {
		return false
	}
	rp.mu.Lock()
	if rp.wRef != nil {
		w := rp.wRef
		rp.mu.Unlock()
		w.Dispatch(fn)
		return true
	}
	select {
	case rp.tasks <- fn:
		rp.mu.Unlock()
		return true
	default:
		rp.mu.Unlock()
		rp.failedPost.Add(1)
		return false
	}
}

// RequestShow asks for the panel to appear. Safe from a ui-sta callback: it does
// not create a window on the caller's thread and does not wait for the pump. A nil
// residentPanel is legal and answers "no", which is what the boot does when the
// assembly failed and the ball host still has the gesture wired.
func (rp *residentPanel) RequestShow(via string) bool {
	if rp == nil {
		return false
	}
	rp.shows.Add(1)
	return rp.post(func() { rp.showOnThread(via) })
}

// showOnThread runs one show request ON the panel thread and, when bringUp refuses
// the thread by name, makes that refusal the thread's state instead of a per-request
// log line.
//
// Why this exists (ticket 33, 33-v2 问①, leg 33-r9). Before this, a refusal left the
// thread exactly as it found it: not created, not failed, still waiting in loop's
// select for the next task. Every later show then walked the same road - post() said
// yes, bringUp said no by name, nothing recorded it - so the process could never open
// a panel and never said so outside a log line, and a caller with a bounded wait read
// 15 seconds of silence instead of a refusal. 33-r7's four-shape table and 33-r9's
// fifth shape (cmd/wisp/panel_host_windows.go's bringUp comment, with readings) both
// say bringUp may NOT clean the queue itself: the message it refuses on can belong to
// a window that is still alive and still owned. So the exit is at THIS level: the
// thread stops pretending it can serve.
//
// The mechanism is the one this file already has for its other fatal verdict (the
// CoInitializeEx guard at the top of loop): store the error as startUp. That makes
// post() answer false on the spot (it already checks startUpErr before routing),
// makes loop retire through its existing teardown instead of waiting forever for a
// task that can only fail, and puts the reason in statusLine(), which is the sentence
// the boot report prints. No new thread, no new pump, no second window, and nothing
// removed or dispatched from anybody else's queue.
func (rp *residentPanel) showOnThread(via string) {
	err := rp.mgr.Show(context.Background())
	if err == nil {
		fmt.Printf("wisp: panel window is up (%s, cold %.1f ms, hot path %.1f ms)\n",
			via, rp.mgr.LastColdMs(), rp.mgr.LastHotMs())
		return
	}
	slog.Error("panel host: show failed on the panel thread", "via", via, "err", err)
	fmt.Printf("wisp: panel could not open (%v): %s\n", err, via)
	if errors.Is(err, errPanelRefusedThread) {
		rp.setStartUp(err)
		slog.Error("panel thread: retiring without a panel window after a named refusal",
			"via", via, "err", err, "shows", rp.shows.Load())
		fmt.Printf("wisp: panel thread will take no further requests: %v\n", err)
	}
}

// RequestToggle is what the panel hot key and the tray item drive: an open panel
// gets hidden (and hands focus back), a hidden or absent one gets shown.
func (rp *residentPanel) RequestToggle(via string) bool {
	if rp == nil {
		return false
	}
	rp.toggles.Add(1)
	if rp.mgr.IsShown() {
		rp.hides.Add(1)
		return rp.post(func() { rp.mgr.Hide() })
	}
	return rp.RequestShow(via)
}

// RequestDispose tears the window down for good on the thread that owns it (the
// recreate path). It keeps the pump running, so a later RequestShow builds a
// fresh window; stop() is what ends the thread.
func (rp *residentPanel) RequestDispose() bool {
	if rp == nil {
		return false
	}
	rp.disposals.Add(1)
	return rp.post(func() { rp.mgr.Destroy() })
}

// stop ends the thread: Terminate posts WM_QUIT so Run() returns and the teardown
// below runs on the panel thread itself. It waits a bounded, monotonic interval for
// the thread to be gone and says out loud whether it made it - the reading ticket
// 33 AC#1's "children exit" clause and ticket 228 AC#11's stop paths both want.
func (rp *residentPanel) stop() {
	if rp == nil {
		return
	}
	rp.closeOnce.Do(func() {
		rp.stopRequested.Store(true)
		rp.mu.Lock()
		w := rp.wRef
		rp.mu.Unlock()
		if w == nil {
			// The thread may already hold a control it has not published yet (it is
			// still inside bring-up). Take that one, or the stop is a request nobody
			// will ever service.
			w = rp.mgr.currentWindow()
		}
		if w == nil {
			close(rp.stopNow)
			return
		}
		// The exit has to be ASKED FOR ON THE PANEL THREAD. go-webview2's
		// Terminate is a bare PostQuitMessage (webview.go:381-383), and that posts
		// WM_QUIT to the queue of the thread that CALLS it, not of the thread that
		// owns the window. Called from here - another goroutine, another thread - the
		// panel thread never sees a quit, Run() keeps pumping, and the process hangs
		// on a window nobody will close. That is how the first version of this
		// function failed on this box, and 33-p1 R26 only worked because its
		// Terminate ran inside a closure on the webview's own thread.
		if !rp.post(func() { rp.mgr.terminateOnThisThread() }) {
			close(rp.stopNow)
		}
	})
	select {
	case <-rp.finished:
		slog.Info("panel thread exited cleanly", "shows", rp.shows.Load(), "toggles", rp.toggles.Load(), "disposals", rp.disposals.Load())
	case <-time.After(panelThreadExitBudget):
		slog.Warn("panel thread did NOT exit within its budget", "budget", panelThreadExitBudget.String())
		fmt.Printf("wisp: panel thread still running after %v - its window may still be open\n", panelThreadExitBudget)
	}
}

// teardown runs the window's last act on the owning thread: Destroy, then clear
// the published handle so nothing can Dispatch into a dead pump.
func (rp *residentPanel) teardown(why string) {
	rp.mu.Lock()
	w := rp.wRef
	rp.mu.Unlock()
	if w != nil {
		rp.mgr.Destroy()
	}
	rp.mu.Lock()
	rp.wRef = nil
	rp.mu.Unlock()
	slog.Info("panel thread ending", "why", why, "window_opened", w != nil)
}

// --- small state readers, so the routing above never touches mgr.mu directly ---

func (rp *residentPanel) isCreated() bool { return rp.mgr.IsCreated() }

func (rp *residentPanel) setStartUp(err error) {
	rp.mu.Lock()
	rp.startUp = err
	rp.mu.Unlock()
}

func (rp *residentPanel) startUpErr() error {
	rp.mu.Lock()
	defer rp.mu.Unlock()
	return rp.startUp
}

// isFinished reports the thread's own handle being done, i.e. its goroutine has
// returned. Requests arriving after that are refused rather than queued forever.
func (rp *residentPanel) isFinished() bool {
	if rp.thread == nil {
		return false
	}
	select {
	case <-rp.thread.Done():
		return true
	default:
		return false
	}
}

// statusLine is the honest one-line answer to "does this process have a panel".
// It reports the thread, not a window: the window is created on first use, so a
// booted process with no gesture yet has a thread and no HWND, and that is the
// truth the boot line has to carry.
func (rp *residentPanel) statusLine() string {
	if rp == nil {
		return "this process has NO panel thread (the panel host never started)"
	}
	if err := rp.startUpErr(); err != nil {
		return "panel thread could not start: " + err.Error()
	}
	created := rp.mgr.IsCreated()
	return fmt.Sprintf("panel thread up (STA, dedicated pump; window %s, %d show request(s), %d toggle(s))",
		map[bool]string{true: "created", false: "not created yet"}[created], rp.shows.Load(), rp.toggles.Load())
}

// probe255r6AliasedCreateSite is the r6 M2 plant: the same create reached
// through an ALIASED import, proving the census matches the selector name, not
// the receiver literal. Restored to the HEAD blob right after the run.
func probe255r6AliasedCreateSite() {
	w := wv2.NewWithOptions(wv2.WebViewOptions{})
	_ = w
}
