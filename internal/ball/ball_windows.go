//go:build windows

package ball

// The ball surface: a WS_EX_LAYERED|NOACTIVATE|TOOLWINDOW|TOPMOST window
// rendered with UpdateLayeredWindow + Direct2D/DirectWrite (SPEC-08 §2).
//
// Interaction contract: the window NEVER activates (WS_EX_NOACTIVATE +
// MA_NOACTIVATE); the orb circle is clickable (drag moves it, click fires
// OnClickBall), the transparent corners fall through via WM_NCHITTEST
// returning HTTRANSPARENT. Animation timers exist ONLY where
// AnimationPolicy grants them - Sleeping has none (zero-timer discipline).

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// EventKind enumerates the user gestures the ball surfaces to its consumer
// (which owns the state machine and decides what each gesture means - e.g.
// a ball click is summon in Sleeping/Warm but veto in Listening/Confirming).
type EventKind uint8

const (
	EvClickBall EventKind = iota
	EvSummonHotkey
	EvMuteHotkey
	EvCancelHotkey // the cancel binding (Esc while taken over, per B1)
	EvPanelHotkey
	EvTrayPanel
	EvTrayMute
	EvTrayPauseWake
	EvTrayExit
	EvDragEnd // position was persisted
)

// Events carries the consumer callbacks. Nil entries are skipped.
type Events struct {
	OnClickBall     func()
	OnSummonHotkey  func()
	OnMuteHotkey    func()
	OnCancelHotkey  func()
	OnPanelHotkey   func()
	OnTrayPanel     func()
	OnTrayMute      func()
	OnTrayPauseWake func()
	OnTrayExit      func()
	OnDragEnd       func()
}

// Options configures New.
type Options struct {
	SizePx   int // configured orb size 44..72 (0 = default 56)
	Theme    Theme
	Initial  statemachine.State
	Hotkeys  HotkeyConfig   // zero value = DefaultHotkeys()
	Store    *PositionStore // nil = no persistence (position defaults)
	Events   Events
	Registry *observe.Registry // nil = observe.Default
	// StartHidden skips ShowWindow (tests that only need the COM stack).
	StartHidden bool
	// WindowTitle overrides the Win32 window title (default "Wisp"). The
	// debug harness gives its child process a private title so evidence
	// scripts can FindWindow the exact ball window instead of a foreign
	// Wisp window left on the desktop by another run.
	WindowTitle string
}

// Ball is the floating ball surface. All mutating methods are goroutine-safe
// (they post to the ui-sta thread). One Ball per process.
type Ball struct {
	opts Options
	sta  *staThread
	hwnd windows.HWND
	rend *renderer
	tray *tray

	// UI-thread-only state below (accessed on the STA thread unless atomic).
	mu              sync.Mutex // guards the fields the STA thread reads on entry
	curState        statemachine.State
	curVisual       Visual
	badge           int
	progress        float32
	badgeText       string
	dragging        bool
	dragMoved       bool
	dragStart       point
	animTimerActive bool
	fadeStart       time.Time
	stateEntered    time.Time
	animPhase       float64
	constAlpha      uint8 // ULW SourceConstantAlpha (Warm breathing / Settling fade)

	// Audio-driven liquid (ticket 62 step 3). liq is the pure model from
	// liquid.go (UI-thread-owned, no locks); liqRaw keeps the last delivered
	// envelope so the bounded burst timer can advance a transition between
	// samples. Only a scalar level ever arrives here (SetAudioLevel) - no
	// audio samples, no transcript text (C25 contamination surface).
	liq               liquidMotion
	liqRaw            float32
	lastMotion        time.Time
	liquidTimerActive bool

	// Edge dock (ticket 62 item 2). UI-thread-only, like liq, and deliberately
	// WITHOUT a timer of its own: lastDock is the monotonic anchor of its frame
	// budget (dockFrameDt), and every step arrives on a mouse message or a
	// window-position event. See dock.go/dock_windows.go.
	dock     dockMotion
	lastDock time.Time

	registeredHotkeys map[uint32]Accelerator
	hotkeyReport      HotkeyReport // outcome of the last registration pass (ticket 64 A1b)
	escTakenOver      bool
	cancelBinding     string // configured cancel binding, for B1 release

	trayMuted     bool
	trayPauseWake bool

	startedErr error
	closed     atomic.Bool
}

// wndproc routing: windows.NewCallback needs stable function pointers, so a
// single package-level proc dispatches to the active Ball.
var (
	activeBall atomic.Pointer[Ball]
	wndProcPtr = windows.NewCallback(ballWndProc)
)

// New boots the ui-sta thread and creates the ball window on it.
func New(opts Options) (*Ball, error) {
	if opts.SizePx == 0 {
		opts.SizePx = BallSizeDefaultPx
	}
	if opts.SizePx < BallSizeMinPx {
		opts.SizePx = BallSizeMinPx
	}
	if opts.SizePx > BallSizeMaxPx {
		opts.SizePx = BallSizeMaxPx
	}
	if opts.Hotkeys == (HotkeyConfig{}) {
		opts.Hotkeys = DefaultHotkeys()
	}
	if opts.Initial == "" {
		opts.Initial = statemachine.StateSleeping
	}
	if opts.Registry == nil {
		opts.Registry = observe.Default
	}

	enablePerMonitorV2DPI()

	b := &Ball{opts: opts, registeredHotkeys: map[uint32]Accelerator{}}
	b.sta = newSTAThread(opts.Registry)

	// ui-sta: the frozen resident roster name (D38b). The D2D factories and
	// every window live and die on this thread.
	b.sta.handle = opts.Registry.Spawn("ui-sta", "ball", nil, func(_ context.Context) {
		b.sta.start(func(s *staThread) error { return b.createOnSTA(s) })
	})

	if err := b.sta.waitStarted(); err != nil {
		return nil, err
	}
	activeBall.Store(b)
	return b, nil
}

// createOnSTA runs on the ui-sta thread: window + renderer + tray + hotkeys.
func (b *Ball) createOnSTA(s *staThread) error {
	if err := ensureFactories(); err != nil {
		return fmt.Errorf("ball: %w", err)
	}

	if err := registerBallClass(); err != nil {
		return err
	}
	clsName := utf16("WispBallWindow")
	title := b.opts.WindowTitle
	if title == "" {
		title = "Wisp"
	}

	exStyle := uintptr(wsExLayered | wsExTopmost | wsExToolWindow | wsExNoActivate)
	hwnd, _, err := pCreateWindowExW.Call(
		exStyle,
		unsafePtr(clsName),
		unsafePtr(utf16(title)),
		0,          // WS_OVERLAPPED; sized/shown below
		0, 0, 1, 1, // placed by restorePosition
		0, 0, uintptr(moduleHandle()), 0)
	if hwnd == 0 {
		return fmt.Errorf("ball: CreateWindowExW: %v", err)
	}
	s.mu.Lock()
	s.hwnd = windows.HWND(hwnd)
	s.mu.Unlock()
	b.hwnd = windows.HWND(hwnd)

	// Window DPI + initial size/position (per-monitor restore). The window is
	// MOVED while hidden, so no WM_DPICHANGED fires - the destination
	// monitor DPI is read explicitly after the move (SPEC-08 §2).
	// The window edge is FIXED for the process lifetime (see applyStateLocked):
	// resolve and place at that size, reading the destination monitor DPI
	// explicitly (hidden windows get no WM_DPICHANGED).
	edge96 := int32(WindowEdgePx(b.opts.SizePx, 96))
	x, y := b.resolveInitial(edge96)
	moveWindow(uintptr(b.hwnd), x, y, 1, 1)
	dpi := getDpiForWindow(b.hwnd) // true DPI of the destination monitor
	edge := int32(WindowEdgePx(b.opts.SizePx, dpi))
	moveWindow(uintptr(b.hwnd), x, y, edge, edge)

	var errR error
	b.rend, errR = newRenderer(b.hwnd, edge, edge, dpi)
	if errR != nil {
		return errR
	}

	// Tray + hotkeys live on the same window.
	t, errT := addTrayIcon(b.hwnd, "Wisp")
	if errT != nil {
		return errT
	}
	b.tray = t
	b.cancelBinding = b.opts.Hotkeys.Cancel
	b.hotkeyReport = registerAll(b.hwnd, b.opts.Hotkeys)
	b.registeredHotkeys = b.hotkeyReport.Live()

	if !b.opts.StartHidden {
		pShowWindow.Call(hwnd, swShownoactivate)
	}

	// Initial state render (no animation timers in static states).
	b.applyStateLocked(b.opts.Initial)
	// A position restored at a work-area edge comes back as the tab the user
	// left it in: the same evaluation a drag ending there runs, and it starts
	// no timer (ticket 62 item 2).
	b.dockCommit()
	return nil
}

func moduleHandle() uintptr {
	h, _, _ := procGetModuleHandleW.Call(0)
	return h
}

func loadArrowCursor() uintptr {
	c, _, _ := pLoadCursorW.Call(0, idcArrow)
	return c
}

// resolveInitial returns the restored position (SPEC-08 §2: per-monitor
// persistence; off-screen / detached-monitor falls back to the primary).
func (b *Ball) resolveInitial(edge int32) (int32, int32) {
	monitors := enumMonitors()
	if len(monitors) == 0 {
		return 100, 100
	}
	var (
		saved      PosEntry
		savedKnown bool
		savedDev   string
	)
	if b.opts.Store != nil {
		// The ball has ONE position: the first saved monitor that still
		// exists wins (positions are keyed by device name).
		for _, m := range monitors {
			if p, ok := b.opts.Store.Get(m.Device); ok {
				saved, savedKnown, savedDev = p, true, m.Device
				break
			}
		}
	}
	_, pos := ResolvePosition(monitors, savedDev, saved, savedKnown, int(edge))
	return int32(pos.X), int32(pos.Y)
}

// ---------------------------------------------------------------- state API

// SetState switches the rendered state (the consumer drives this after its
// statemachine dispatch). Safe from any goroutine.
func (b *Ball) SetState(s statemachine.State) {
	b.sta.PostTask(func() { b.applyStateLocked(s) })
}

// SetBadge updates the approval-queue badge depth (AwaitingApproval).
func (b *Ball) SetBadge(n int) {
	b.sta.PostTask(func() {
		b.badge = n
		b.repaintCurrent()
	})
}

// SetProgress updates the Downloading progress ring (0..1).
func (b *Ball) SetProgress(p float32) {
	b.sta.PostTask(func() {
		b.progress = p
		b.repaintCurrent()
	})
}

// SetBadgeText sets the center text override (L1 countdown digits / percent).
func (b *Ball) SetBadgeText(s string) {
	b.sta.PostTask(func() {
		b.badgeText = s
		b.repaintCurrent()
	})
}

// applyStateLocked is the STA-thread state switch: timer policy, window
// size, one static frame. (The name is historical: the state fields are
// STA-thread-owned, not mutex-shared.)
func (b *Ball) applyStateLocked(s statemachine.State) {
	prev := b.curState
	b.curState = s
	b.motionStateChanged(prev, s) // liquid session start/stop on state edges
	b.curVisual = VisualFor(pal, b.opts.SizePx, s, 0)
	b.badge = 0
	b.progress = -1
	b.badgeText = ""
	b.animPhase = 0
	b.stateEntered = time.Now()

	policy := AnimationPolicy(s)
	b.stopAnimTimer()
	b.animTimerActive = policy.PeriodMs > 0
	if b.animTimerActive {
		pSetTimer.Call(uintptr(b.hwnd), timerAnimID, uintptr(policy.PeriodMs), 0)
	}
	b.syncLiquidTimer() // arms only where liquidDriven+Animated already grant motion

	// Window size: FIXED for the process lifetime (the configured orb + ring
	// margins). States render within it - the Sleeping micro dot is a 12px
	// visual centered in the same window, and its extra margins are
	// click-through (WM_NCHITTEST). Resizing would force a DC-render-target
	// rebind, which empirically yields transparent frames - avoided entirely.

	// Constant alpha: 255 for everything except Settling (fade target) -
	// Warm breathing and the Settling fade ramp drive this value per tick.
	b.constAlpha = 255

	b.rend.setAssets(b.curVisual)
	b.renderFrame()
}

// repaintCurrent re-renders the current state (badge/progress/text updates).
func (b *Ball) repaintCurrent() {
	v := VisualFor(pal, b.opts.SizePx, b.curState, 0)
	v.BadgeCount = b.badge
	if b.progress >= 0 {
		v.Progress = b.progress
	}
	if b.badgeText != "" {
		v.BadgeText = b.badgeText
	}
	b.curVisual = v
	b.rend.setAssets(v)
	b.renderFrame()
}

// renderFrame draws and presents via UpdateLayeredWindow.
func (b *Ball) renderFrame() {
	if b.rend == nil || b.rend.dcRT == nil {
		return
	}
	t0 := time.Now()
	b.rend.drawFrame(b.frameVisual(), b.animPhase)
	if d := time.Since(t0); d > 30*time.Millisecond {
		slog.Warn("ball: slow drawFrame", "dur", d.String(), "state", string(b.curState))
	}

	var wr rect
	pGetWindowRect.Call(uintptr(b.hwnd), unsafePtr(&wr))
	dst := point{wr.l, wr.t}
	src := point{0, 0}
	sz := size{b.rend.w, b.rend.h}
	blend := blendFunction{
		blendOp:             acSrcOver,
		blendFlags:          0,
		sourceConstantAlpha: b.constAlpha,
		alphaFormat:         acSrcAlpha,
	}
	r, _, err := pUpdateLayeredWindow.Call(
		uintptr(b.hwnd), 0, // hdcDst
		unsafePtr(&dst),
		unsafePtr(&sz),
		uintptr(b.rend.hdcMem),
		unsafePtr(&src),
		0,
		unsafePtr(&blend),
		ulwAlpha)
	if d := time.Since(t0); d > 30*time.Millisecond {
		slog.Warn("ball: slow ULW frame", "dur", d.String(), "state", string(b.curState))
	}
	if r == 0 {
		slog.Warn("ball: UpdateLayeredWindow failed", "err", err)
	}
}

// stopAnimTimer kills the animation timer (zero-timer discipline for static
// states, especially Sleeping).
func (b *Ball) stopAnimTimer() {
	if b.animTimerActive {
		pKillTimer.Call(uintptr(b.hwnd), timerAnimID)
		b.animTimerActive = false
	}
}

// TimersAlive reports whether an animation timer (state loop, breathing, fade
// or the bounded liquid burst) is currently armed. The consumer can assert
// TimersAlive()==false in Sleeping (ticket acceptance).
func (b *Ball) TimersAlive() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.animTimerActive || b.liquidTimerActive
}

// DebugTimersAlive reads the LIVE timer flags on the UI thread (the value the
// evidence harness prints). Unlike TimersAlive it is synchronised with the
// STA thread, so it cannot report a stale value right after SetState.
func (b *Ball) DebugTimersAlive() bool {
	done := make(chan bool, 1)
	b.sta.PostTask(func() { done <- b.animTimerActive || b.liquidTimerActive })
	return <-done
}

// onAnimTick runs on the STA thread from WM_TIMER.
func (b *Ball) onAnimTick() {
	state := b.curState
	policy := AnimationPolicy(state)
	if policy.Kind == AnimNone || policy.PeriodMs == 0 {
		b.stopAnimTimer()
		return
	}
	elapsed := time.Since(b.stateEntered)

	switch policy.Kind {
	case AnimAlpha: // Warm breathing: constant alpha only, no re-render.
		t := (elapsed % (WarmBreathPeriodMs * time.Millisecond)).Seconds() / float64(WarmBreathPeriodMs)
		w := 0.5 - 0.5*math.Cos(2*math.Pi*t) // ease-in-out 0..1
		a := WarmOpacityLow + (WarmOpacityHigh-WarmOpacityLow)*w
		b.constAlpha = uint8(a * 255)
		b.renderFrame()
	case AnimFade: // Settling: 260ms one-shot ramp 1 -> 0.35 (SPEC-08 §2.1).
		p := elapsed.Seconds() / (float64(SettlingFadeMs) / 1000)
		if p >= 1 {
			b.constAlpha = u8(float32(SleepOpacity) * 255.0)
			b.stopAnimTimer()
		} else {
			b.constAlpha = u8(float32((1 - p*(1-SleepOpacity)) * 255.0))
		}
		b.renderFrame()
	default: // AnimFrame: phase-advanced full re-render at <= 30fps.
		cycleMs := animCycleMs(state)
		b.animPhase = elapsed.Seconds() * 1000 / float64(cycleMs)
		b.repaintAnimated()
	}
}

// animCycleMs returns the base cycle of the state's phase animation.
func animCycleMs(s statemachine.State) int {
	switch s {
	case statemachine.StateListening:
		return ListeningWaveMs
	case statemachine.StateThinking:
		return ThinkingSweepMs
	case statemachine.StateConfirming:
		return ConfirmingPulseMs
	case statemachine.StateSpeaking:
		return SpeakingBreathMs
	}
	return 1000
}

// repaintAnimated re-renders with the current overlay values (no asset
// rebuild - the cached brushes are phase-independent).
func (b *Ball) repaintAnimated() {
	v := b.frameVisual()
	v.BadgeCount = b.badge
	if b.progress >= 0 {
		v.Progress = b.progress
	}
	if b.badgeText != "" {
		v.BadgeText = b.badgeText
	}
	b.rend.drawFrame(v, b.animPhase)
	b.present()
}

// present is renderFrame without the draw (used after manual draws).
func (b *Ball) present() {
	var wr rect
	pGetWindowRect.Call(uintptr(b.hwnd), unsafePtr(&wr))
	dst := point{wr.l, wr.t}
	src := point{0, 0}
	sz := size{b.rend.w, b.rend.h}
	blend := blendFunction{acSrcOver, 0, b.constAlpha, acSrcAlpha}
	r, _, err := pUpdateLayeredWindow.Call(
		uintptr(b.hwnd), 0, unsafePtr(&dst), unsafePtr(&sz),
		uintptr(b.rend.hdcMem), unsafePtr(&src), 0, unsafePtr(&blend), ulwAlpha)
	if r == 0 {
		slog.Warn("ball: UpdateLayeredWindow failed", "err", err)
	}
}

// renderFrame = draw + present.
// (kept as the single entry: drawFrame then present)

// ------------------------------------------------------------------ wndproc

func ballWndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	b := activeBall.Load()
	if b == nil || b.hwnd == 0 || hwnd != uintptr(b.hwnd) {
		r, _, _ := pDefWindowProcW.Call(hwnd, msg, wParam, lParam)
		return r
	}
	return b.wndProc(hwnd, msg, wParam, lParam)
}

func (b *Ball) wndProc(hwnd, m, wParam, lParam uintptr) uintptr {
	switch m {
	case wmNCHitTest:
		// Click-through outside the orb circle (SPEC-08 §2). The clickable
		// circle follows the CURRENT visual size (the Sleeping micro dot is
		// small; the configured orb is full-size).
		var pt point
		pt.x, pt.y = loSigned(lParam), hiSigned(lParam) // screen coords
		pScreenToClient.Call(hwnd, uintptr(unsafe.Pointer(&pt)))
		var wr rect
		pGetWindowRect.Call(hwnd, unsafePtr(&wr))
		code := HitTest(pt.x, pt.y, wr.width(), int(b.curVisual.SizePx), getDpiForWindow(windows.HWND(hwnd)))
		if code == HTTransparent {
			return htTransparentResult()
		}
		return uintptr(HTClient)

	case wmMouseActivate:
		// Never steal focus (SPEC-08 §2).
		return 3 // MA_NOACTIVATE

	case wmLButtonDown:
		b.mu.Lock()
		b.dragging = true
		b.dragMoved = false
		b.dragStart = point{loSigned(lParam), hiSigned(lParam)}
		b.mu.Unlock()
		b.dock.supped = false // a new drag re-arms the edge dock
		b.lastDock = time.Now()
		pSetCapture.Call(hwnd)
		return 0

	case wmMouseMove:
		b.mu.Lock()
		dragging := b.dragging
		b.mu.Unlock()
		if dragging {
			var pt point
			pt.x, pt.y = loSigned(lParam), hiSigned(lParam)
			b.mu.Lock()
			dx := pt.x - b.dragStart.x
			dy := pt.y - b.dragStart.y
			moved := b.dragMoved
			if dx*dx+dy*dy > 16 { // 4px drag threshold
				b.dragMoved = true
				moved = true
			}
			b.mu.Unlock()
			if moved {
				var wr rect
				pGetWindowRect.Call(hwnd, unsafePtr(&wr))
				pSetWindowPos.Call(hwnd, 0,
					uintptr(wr.l+dx), uintptr(wr.t+dy), 0, 0,
					uintptr(swpNoSize|swpNoZOrder|swpNoActivate))
				// The push toward an edge IS the shrink animation: this message
				// carries the frame, so the dock needs no timer (D32).
				b.dockProximity()
			}
			return 0
		}
		// Not dragging, and the message reached us at all only because
		// WM_NCHITTEST called this pixel part of the ball: the pointer is on a
		// docked tab, so the orb slides back out (ticket 62 item 2).
		b.dockHoverMove()
		return 0

	case wmLButtonUp:
		pReleaseCapture.Call(hwnd)
		b.mu.Lock()
		wasDrag := b.dragMoved
		b.dragging = false
		b.dragMoved = false
		b.mu.Unlock()
		if wasDrag {
			// Decide the dock BEFORE the position is written down, so what
			// survives a restart is what the user ended with.
			b.dockCommit()
			b.persistPosition()
			// The callback runs on the STA thread: it must be quick and
			// non-blocking (dispatch the machine, PostTask the rest) - the
			// message pump may not stall.
			b.fire(b.opts.Events.OnDragEnd)
		} else if !b.dockClickOut() {
			b.fire(b.opts.Events.OnClickBall)
		}
		return 0

	case wmMouseLeave:
		b.dockLeave()
		return 0

	case wmCaptureChanged:
		b.mu.Lock()
		b.dragging = false
		b.mu.Unlock()
		return 0

	case wmTimer:
		if wParam == timerAnimID {
			b.onAnimTick()
			return 0
		}
		if wParam == timerLiquidID {
			b.onLiquidTick()
			return 0
		}

	case wmHotkey:
		switch wParam {
		case hkSummon:
			b.fire(b.opts.Events.OnSummonHotkey)
		case hkMute:
			b.fire(b.opts.Events.OnMuteHotkey)
		case hkCancel:
			b.fire(b.opts.Events.OnCancelHotkey)
		case hkPanel:
			b.fire(b.opts.Events.OnPanelHotkey)
		}
		return 0

	case wmAppTray:
		switch lParam & 0xFFFF {
		case wmLButtonUp: // left click = open panel (SPEC-08 §7)
			b.fire(b.opts.Events.OnTrayPanel)
		case wmRButtonUp:
			sel := showMenu(b.hwnd, b.trayMuted, b.trayPauseWake)
			switch sel {
			case menuOpenPanel:
				b.fire(b.opts.Events.OnTrayPanel)
			case menuMute:
				b.fire(b.opts.Events.OnTrayMute)
			case menuPauseWake:
				b.fire(b.opts.Events.OnTrayPauseWake)
			case menuExit:
				b.fire(b.opts.Events.OnTrayExit)
			}
		}
		return 0

	case wmAppTask:
		b.sta.runTask(uint64(wParam))
		return 0

	case wmDPICHanged:
		// Per-Monitor V2: full renderer rebuild at the new DPI (S1: the
		// ticket's sanctioned simplification - DC render targets are
		// recreated, not rebound).
		newDpi := uint32(wParam & 0xFFFF)
		sug := rectFromUintptr(lParam)
		edge := int32(WindowEdgePx(b.opts.SizePx, newDpi))
		moveWindow(hwnd, sug.l, sug.t, edge, edge)
		b.rend.release()
		var errR error
		b.rend, errR = newRenderer(windows.HWND(hwnd), edge, edge, newDpi)
		if errR != nil {
			slog.Error("ball: DPI renderer rebuild failed", "err", errR)
			return 0
		}
		b.rend.setAssets(b.curVisual)
		b.dockCommit() // tab geometry is physical px: re-derive it at the new DPI
		b.renderFrame()
		return 0

	case wmDisplayChange:
		// Topology change: re-clamp into the visible area, then re-decide the
		// dock against the new work area (a detached monitor must not leave a
		// tab stranded off-screen).
		b.reclampPosition()
		b.dockCommit()
		return 0
	}
	r, _, _ := pDefWindowProcW.Call(hwnd, m, wParam, lParam)
	return r
}

// fire runs a consumer callback on the STA thread. CONTRACT: callbacks must
// be quick and non-blocking (statemachine dispatch + PostTask); a blocking
// callback stalls the message pump. No goroutine is spawned per gesture
// (D38b roster discipline).
func (b *Ball) fire(fn func()) {
	if fn != nil {
		fn()
	}
}

var pScreenToClient = modUser32.NewProc("ScreenToClient")

// persistPosition saves the window position keyed by its current monitor.
func (b *Ball) persistPosition() {
	if b.opts.Store == nil {
		return
	}
	var wr rect
	pGetWindowRect.Call(uintptr(b.hwnd), unsafePtr(&wr))
	dev := monitorFromWindow(b.hwnd)
	if dev == "" {
		return
	}
	if err := b.opts.Store.Put(dev, PosEntry{X: int(wr.l), Y: int(wr.t)}); err != nil {
		slog.Warn("ball: position persist failed", "err", err)
	}
}

// reclampPosition re-places the window when it would be off-screen (monitor
// detach / resolution change): falls back to the primary's default spot.
func (b *Ball) reclampPosition() {
	monitors := enumMonitors()
	if len(monitors) == 0 {
		return
	}
	var wr rect
	pGetWindowRect.Call(uintptr(b.hwnd), unsafePtr(&wr))
	dev := monitorFromWindow(b.hwnd)
	cur := Rect{int(wr.l), int(wr.t), int(wr.r), int(wr.b)}
	visible := false
	for _, mo := range monitors {
		if mo.Device == dev && mo.Screen.intersects(cur) {
			visible = true
			break
		}
	}
	if visible {
		return
	}
	m, pos := ResolvePosition(monitors, dev, PosEntry{int(wr.l), int(wr.t)}, true, int(wr.width()))
	_ = m
	moveWindow(uintptr(b.hwnd), int32(pos.X), int32(pos.Y), wr.width(), wr.height())
}

// ------------------------------------------------------------------ hotkeys

// RebindHotkeys re-registers all hotkeys after a config change (hot reload)
// and returns what happened. Callers must NOT drop the report on the floor:
// a binding that is occupied by another app is a user-visible fact, and the
// HotkeyReloader bridge logs it (ticket 64 A1). Rebind is idempotent - same
// bindings in, same live set out - so a host may call it on every config
// poll; HotkeyReloader still diffs so the Win32 churn happens once per change.
func (b *Ball) RebindHotkeys(cfg HotkeyConfig) HotkeyReport {
	done := make(chan struct{})
	var rep HotkeyReport
	b.sta.PostTask(func() {
		defer close(done)
		unregisterAll(b.hwnd)
		rep = registerAll(b.hwnd, cfg)
		b.hotkeyReport = rep
		b.registeredHotkeys = rep.Live()
		b.cancelBinding = cfg.Cancel
		b.escTakenOver = false
	})
	<-done
	return rep
}

// HotkeyReport returns the outcome of the last registration pass (boot or
// rebind). Synchronised with the UI thread, so it is authoritative the moment
// RebindHotkeys/New return.
func (b *Ball) HotkeyReport() HotkeyReport {
	done := make(chan HotkeyReport, 1)
	b.sta.PostTask(func() { done <- b.hotkeyReport })
	return <-done
}

// ConfiguredHotkeys returns the binding set the ball was told to hold (the
// input of the last registration pass, not its outcome - see HotkeyReport for
// what Win32 actually accepted). opts is immutable after New.
func (b *Ball) ConfiguredHotkeys() HotkeyConfig { return b.opts.Hotkeys }

// RegisteredHotkeys returns the id -> accelerator set Win32 actually holds
// for us right now (the registration set the acceptance asks tests to
// assert, not a mirror of the config).
func (b *Ball) RegisteredHotkeys() map[uint32]Accelerator {
	done := make(chan map[uint32]Accelerator, 1)
	b.sta.PostTask(func() { done <- cloneAccel(b.registeredHotkeys) })
	return <-done
}

func cloneAccel(m map[uint32]Accelerator) map[uint32]Accelerator {
	out := make(map[uint32]Accelerator, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// TakeEscForCancel temporarily binds Esc as the cancel key during Confirming
// (B1). Idempotent.
func (b *Ball) TakeEscForCancel() {
	done := make(chan struct{})
	b.sta.PostTask(func() {
		defer close(done)
		if b.escTakenOver {
			return
		}
		if takeEsc(b.hwnd) {
			b.escTakenOver = true
		}
	})
	<-done
}

// ReleaseEscAfterSession hands Esc back to the configured cancel binding
// (B1: "会话结束必须归还"). Idempotent.
func (b *Ball) ReleaseEscAfterSession() {
	done := make(chan struct{})
	b.sta.PostTask(func() {
		defer close(done)
		if !b.escTakenOver {
			return
		}
		releaseEsc(b.hwnd, b.cancelBinding)
		b.escTakenOver = false
	})
	<-done
}

// EscTakenOver reports the B1 takeover state (test seam).
func (b *Ball) EscTakenOver() bool {
	done := make(chan bool)
	b.sta.PostTask(func() { done <- b.escTakenOver })
	return <-done
}

// --------------------------------------------------------------------- tray

// SetTrayChecks updates the mute / pause-wake checkmarks.
func (b *Ball) SetTrayChecks(muted, pausedWake bool) {
	b.sta.PostTask(func() {
		b.trayMuted = muted
		b.trayPauseWake = pausedWake
	})
}

// SetTrayTip updates the tray tooltip (state surfacing).
func (b *Ball) SetTrayTip(tip string) {
	b.sta.PostTask(func() { b.tray.setTip(tip) })
}

// ------------------------------------------------------------------ closing

// Close destroys every window/resource on the STA thread and stops it.
// After Close returns, the process holds no ball-side USER/GDI/D2D objects.
func (b *Ball) Close() {
	if b.closed.Swap(true) {
		return
	}
	done := make(chan struct{})
	b.sta.PostTask(func() {
		defer close(done)
		b.stopAnimTimer()
		b.stopLiquidTimer()
		unregisterAll(b.hwnd)
		if b.tray != nil {
			b.tray.remove()
		}
		if b.rend != nil {
			b.rend.release()
		}
		if b.hwnd != 0 {
			// Ticket 64 A6: the return value used to be dropped, which made a
			// failed destroy (and the USER/GDI handles it leaves behind)
			// invisible. A dead window is now loud, not silent.
			if r, _, err := pDestroyWindow.Call(uintptr(b.hwnd)); r == 0 {
				slog.Error("ball: DestroyWindow failed; window handles stay alive in this process",
					"hwnd", uintptr(b.hwnd), "err", err)
			}
			b.hwnd = 0
		}
		activeBall.CompareAndSwap(b, nil)
		b.sta.quit()
	})
	<-done
	if b.sta.handle != nil {
		<-b.sta.handle.Done() // join the ui-sta thread
	}
}

// recenterAt re-sizes/re-positions the window for the destination monitor's
// DPI after the hidden move (keeps the restored top-left corner).
func (b *Ball) recenterAt(dpi uint32, edge int32) {
	var wr rect
	pGetWindowRect.Call(uintptr(b.hwnd), unsafePtr(&wr))
	moveWindow(uintptr(b.hwnd), wr.l, wr.t, edge, edge)
}

// registerBallClass registers the ball window class once per process
// (re-registration is an error; a second Ball in-process reuses it).
var (
	classOnce sync.Once
	classErr  error
)

func registerBallClass() error {
	classOnce.Do(func() {
		clsName := utf16("WispBallWindow")
		wc := wndClassEx{
			size:      uint32(unsafe.Sizeof(wndClassEx{})),
			wndProc:   wndProcPtr,
			instance:  windows.Handle(moduleHandle()),
			cursor:    windows.Handle(loadArrowCursor()),
			className: clsName,
		}
		if r, _, _ := pRegisterClassExW.Call(unsafePtr(&wc)); r == 0 {
			classErr = fmt.Errorf("ball: RegisterClassExW failed")
		}
	})
	return classErr
}

// DebugHWND exposes the ball window handle (debug harness/evidence only).
func (b *Ball) DebugHWND() windows.HWND { return b.hwnd }
