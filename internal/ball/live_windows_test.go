//go:build windows && winlive

package ball

// Live-window tests: they create the REAL ball window on the user's desktop
// (hidden=false, topmost) and MUST tear everything down. Run explicitly:
//
//	go test -tags winlive ./internal/ball/ -run TestBallLive -v
//
// They are excluded from the default suite so `go test ./...` never pops
// windows. The ticket 07 acceptance items exercised here: window stack
// handle gate (<600), zero animation timers in Sleeping, B1 Esc
// takeover/release, per-monitor position persistence, full teardown.

import (
	"context"
	"path/filepath"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// handlesOfProcess is the PROCESS-LOCAL handle count (GetProcessHandleCount):
// no other process can raise it, which is why registry A6's growth question can
// be answered here even though registry A5's window lookup could not.
func handlesOfProcess() uint32 {
	var n uint32
	procGetProcessHandleCount.Call(uintptr(windows.CurrentProcess()), uintptr(unsafe.Pointer(&n)))
	return n
}

var procGetProcessHandleCount = modKernel32.NewProc("GetProcessHandleCount")

// TestBallLiveLifecycle walks all 20 states on a real window and tears down.
// Ticket 64 rewrote its window assertions to go by the test's OWN hwnd and
// added the quiet-desktop precondition: a stray ball from another process used
// to make this test red for the wrong reason (registry A5), and its handle
// numbers are now logged every iteration so a leak can be classified from data
// instead of from a guess (registry A6).
func TestBallLiveLifecycle(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "ballpos.json")
	store, err := OpenPositionStore(storePath)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	requireQuietBallDesktop(t)
	baseHandles := handlesOfProcess()
	baseUser := uiUserObjects()

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
	hwnd := b.DebugHWND()
	if !windowAlive(hwnd) {
		t.Fatalf("ball window %v is not a window after New", hwnd)
	}
	if got := b.HotkeyReport(); !got.AllLive() || len(got.Live()) != 4 {
		t.Fatalf("the live ball did not register its four hotkeys: %+v", got.Bindings())
	}
	afterNew := handlesOfProcess()

	t.Cleanup(func() {
		b.Close()
		if windowAlive(hwnd) {
			t.Errorf("ball window %v still alive after Close", hwnd)
		}
		afterClose := handlesOfProcess()
		t.Logf("HANDLES base=%d afterNew=%d afterClose=%d delta=%d | USER objs base=%d afterClose=%d",
			baseHandles, afterNew, afterClose, int64(afterClose)-int64(baseHandles), baseUser, uiUserObjects())
		requireQuietBallDesktop(t)
	})

	// All 20 states render without hanging; animated states hold a timer,
	// static states none.
	for _, s := range []statemachine.State{
		statemachine.StateFirstRun, statemachine.StateSleeping, statemachine.StateArmed,
		statemachine.StateMuted, statemachine.StateListening, statemachine.StateThinking,
		statemachine.StateActing, statemachine.StateSpeaking, statemachine.StateWarm,
		statemachine.StateConversation, statemachine.StateConfirming,
		statemachine.StateAwaitingApproval, statemachine.StateSettling,
		statemachine.StateDownloading, statemachine.StateUnconfigured,
		statemachine.StateNoNetwork, statemachine.StateError, statemachine.StateQueued,
		statemachine.StateStuck, statemachine.StateWatchdogAlert,
	} {
		b.SetState(s)
		time.Sleep(80 * time.Millisecond)
		if want := Animated(s) && s != statemachine.StateActing || s == statemachine.StateWarm || s == statemachine.StateSettling; want != b.TimersAlive() {
			t.Errorf("%s: TimersAlive=%v, want %v", s, b.TimersAlive(), want)
		}
	}

	// Zero-timer discipline in Sleeping (ticket acceptance).
	b.SetState(statemachine.StateSleeping)
	time.Sleep(150 * time.Millisecond)
	if b.TimersAlive() {
		t.Fatal("animation timer alive in Sleeping (zero-timer discipline broken)")
	}

	// B1: Esc takeover during Confirming, released at session end.
	b.SetState(statemachine.StateConfirming)
	b.TakeEscForCancel()
	time.Sleep(100 * time.Millisecond)
	if !b.EscTakenOver() {
		t.Log("Esc registration busy (another app holds it); takeover not assertable this run")
	} else {
		b.SetState(statemachine.StateSleeping)
		b.ReleaseEscAfterSession()
		time.Sleep(100 * time.Millisecond)
		if b.EscTakenOver() {
			t.Fatal("Esc not returned after session end (B1 violation)")
		}
	}

	// Handle gate (orchestrator ruling, docs/SLO.md): window-bearing stack
	// stays below 600.
	if h := handlesOfProcess(); h >= 600 {
		t.Errorf("handle gate exceeded: %d >= 600 (base %d)", h, baseHandles)
	} else {
		t.Logf("handles: base=%d peak-suite=%d (gate <600)", baseHandles, h)
	}
}

// TestBallLivePositionPersistence drags the ball programmatically (moves the
// window), closes, re-opens and expects the position restored.
func TestBallLivePositionPersistence(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "ballpos.json")
	store, err := OpenPositionStore(storePath)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	requireQuietBallDesktop(t)

	b, err := New(Options{
		Initial: statemachine.StateSleeping, Store: store,
		WindowTitle: liveBallTitle(t), Hotkeys: liveHotkeys(),
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	hwnd := b.DebugHWND()
	time.Sleep(100 * time.Millisecond)

	// Simulate a drag: move the window, then persist through the same path
	// the wndproc uses.
	var wr rect
	pGetWindowRect.Call(uintptr(b.hwnd), unsafePtr(&wr))
	want := PosEntry{X: int(wr.l) - 100, Y: int(wr.t) - 60}
	moveWindow(uintptr(b.hwnd), int32(want.X), int32(want.Y), wr.width(), wr.height())
	b.persistPosition()
	dev := monitorFromWindow(b.hwnd)
	b.Close()
	if windowAlive(hwnd) {
		t.Fatal("the ball window survived Close")
	}
	time.Sleep(100 * time.Millisecond)

	got, ok := store.Get(dev)
	if !ok || got != want {
		t.Fatalf("persisted position: got %+v ok=%v, want %+v (device %s)", got, ok, want, dev)
	}

	// Re-open: the restored top-left must equal the saved position.
	b2, err := New(Options{
		Initial: statemachine.StateSleeping, Store: store,
		WindowTitle: liveBallTitle(t) + "-2", Hotkeys: liveHotkeys(),
	})
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer func() {
		h2 := b2.DebugHWND()
		b2.Close()
		if windowAlive(h2) {
			t.Errorf("the reopened ball window %v survived Close", h2)
		}
	}()
	time.Sleep(150 * time.Millisecond)
	var wr2 rect
	pGetWindowRect.Call(uintptr(b2.hwnd), unsafePtr(&wr2))
	if int(wr2.l) != want.X || int(wr2.t) != want.Y {
		t.Fatalf("restored rect (%d,%d), want (%d,%d)", wr2.l, wr2.t, want.X, want.Y)
	}
}

// TestBallLiveIdleBorderTransition (ticket 62 item 1): the border must arrive
// through a transition that STATE changes drive, with no audio in the picture,
// and the bounded burst that carries it must retire itself and must never
// appear in a state the frozen policy keeps static.
func TestBallLiveIdleBorderTransition(t *testing.T) {
	restored := PrototypeVisualsEnabled()
	t.Cleanup(func() { EnablePrototypeVisuals(restored) })
	EnablePrototypeVisuals(true)

	b, err := New(Options{Initial: statemachine.StateSleeping, StartHidden: true})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer b.Close()

	probe := make(chan Visual, 1)
	readVisual := func() Visual {
		b.sta.PostTask(func() { probe <- b.frameVisual() })
		return <-probe
	}
	ask := func(f func() bool) bool {
		done := make(chan bool, 1)
		b.sta.PostTask(func() { done <- f() })
		return <-done
	}
	timerOn := func() bool { return ask(func() bool { return b.liquidTimerActive }) }

	// Sleeping: the resting frame carries no border and no motion authority.
	b.SetState(statemachine.StateSleeping)
	if v := readVisual(); v.BorderAlpha != 0 {
		t.Fatalf("the Sleeping frame must carry no border: %v", v.BorderAlpha)
	}
	if timerOn() {
		t.Fatal("Sleeping armed the transition timer")
	}

	// A session, then the idle states that follow it: every one of them that
	// the frozen policy grants a timer to must GLIDE the border in, which means
	// at least one painted frame strictly between 0 and 1.
	for _, s := range []statemachine.State{
		statemachine.StateListening, statemachine.StateThinking,
		statemachine.StateActing, statemachine.StateConfirming,
		statemachine.StateWarm, statemachine.StateSettling,
	} {
		b.SetState(statemachine.StateSpeaking) // a fresh session: border away
		if v := readVisual(); v.BorderAlpha != 0 {
			t.Fatalf("%s from Speaking: border should be out, got %v", s, v.BorderAlpha)
		}
		b.SetState(s)
		var frames []float32
		settled := false
		for i := 0; i < 300; i++ { // <=2.4s of polls; the burst is bounded at ~1s
			frames = append(frames, readVisual().BorderAlpha)
			if !timerOn() {
				settled = true
				break
			}
			time.Sleep(8 * time.Millisecond)
		}
		if !settled {
			t.Fatalf("%s: the border burst never retired (polled %d times)", s, len(frames))
		}
		// The tick that lands the border also retires the timer, and it can
		// fire between an iteration's read and its timer probe: read once
		// more, which is then the authoritative end value.
		frames = append(frames, readVisual().BorderAlpha)
		last := frames[len(frames)-1]
		if last != borderAtRest(s) {
			t.Fatalf("%s: the border must land at its resting level, frames=%v", s, frames)
		}
		mid := false
		for _, f := range frames {
			if f > 0 && f < 1 {
				mid = true
				break
			}
		}
		if !mid {
			t.Fatalf("%s: the border arrived in one frame, no transition frames: %v", s, frames)
		}
	}

	// Static frames keep the frozen mapping and no timer, even after a border
	// has been travelling.
	for _, s := range []statemachine.State{
		statemachine.StateArmed, statemachine.StateMuted,
		statemachine.StateConversation, statemachine.StateError,
	} {
		b.SetState(s)
		if timerOn() {
			t.Fatalf("%s armed the transition timer (frozen policy grants it none)", s)
		}
	}

	// Back to Sleeping: nothing armed, and the frame is the committed one.
	b.SetState(statemachine.StateSleeping)
	if b.DebugTimersAlive() {
		t.Fatal("a timer survived the return to Sleeping")
	}
	if v := readVisual(); v.BorderAlpha != 0 || !v.Glass {
		t.Fatalf("Sleeping frame after a session: %+v", v)
	}
}

// TestBallLiveEdgeDock (ticket 62 item 2, owner: "球在靠近电脑侧边的时候，应该会
// 自动收缩…鼠标悬停或单击→丝滑弹回完整球"): on a real window, docking each of the
// four edges must move the window so the squashed orb is tangent to that
// monitor's WORK area (never under the taskbar), leave a docked Sleeping orb at
// zero timers, and let the hover path walk it back out.
func TestBallLiveEdgeDock(t *testing.T) {
	restored := PrototypeVisualsEnabled()
	t.Cleanup(func() { EnablePrototypeVisuals(restored) })
	EnablePrototypeVisuals(true)

	b, err := New(Options{Initial: statemachine.StateSleeping, StartHidden: true})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer b.Close()

	liveRect := func() (Rect, MonitorRect) {
		var wr rect
		pGetWindowRect.Call(uintptr(b.hwnd), unsafePtr(&wr))
		m, ok := monitorRectsForWindow(b.hwnd)
		if !ok {
			t.Fatal("no monitor for the ball window")
		}
		return Rect{int(wr.l), int(wr.t), int(wr.r), int(wr.b)}, m
	}

	for _, e := range []Edge{EdgeLeft, EdgeRight, EdgeTop, EdgeBottom} {
		if !b.DebugDock(e) {
			t.Fatalf("edge %d: DebugDock did not commit the dock", e)
		}
		edge, p, owns := b.Docked()
		if !owns || edge != e || p != 1 {
			t.Fatalf("edge %d: docked reports edge=%d p=%v owns=%v", e, edge, p, owns)
		}
		// D32: a docked, sleeping orb is still a zero-timer orb.
		if b.DebugTimersAlive() {
			t.Fatalf("edge %d: docking armed a timer in Sleeping", e)
		}
		// The law holds on the live window: the drawn orb touches the work
		// boundary along the dock axis and nothing else moved it.
		wr, m := liveRect()
		edgePx := wr.R - wr.L
		marginPx := int32(RingMarginPx * float32(getDpiForWindow(b.hwnd)) / 96)
		R := int32(edgePx/2) - marginPx
		drawn := int(float32(R) * DockSquash(p))
		half := edgePx / 2
		var got, want int
		switch e {
		case EdgeLeft:
			got, want = wr.L+half-drawn, m.Work.L
		case EdgeRight:
			got, want = wr.L+half+drawn, m.Work.R
		case EdgeTop:
			got, want = wr.T+half-drawn, m.Work.T
		case EdgeBottom:
			got, want = wr.T+half+drawn, m.Work.B
		}
		if got < want-2 || got > want+2 {
			t.Fatalf("edge %d: docked orb boundary %d, want the work edge %d", e, got, want)
		}
		// Only TRANSPARENT pixels may ride past the work boundary: the orb is
		// tangent to it, so no part of the ball hides under the taskbar
		// (D42#3). `got` is the orb's boundary, from the check just above.
		if e == EdgeBottom && got > m.Work.B+2 {
			t.Fatalf("a bottom tab must keep the orb above the taskbar: orb bottom %d, work %d",
				got, m.Work.B)
		}
		// The pop-out ramp (a hidden window cannot be hovered; the real-pointer
		// hover is TestBallLiveEdgeDockHover below). These are the two calls a
		// WM_MOUSEMOVE makes, with the dt the message would have carried.
		done := make(chan float32, 1)
		step := func(target float32) float32 {
			b.sta.PostTask(func() {
				b.dock.target = target
				b.dockStep(40 * time.Millisecond) // the 30fps budget of one move
				done <- b.dock.p
			})
			return <-done
		}
		var seen []float32
		for i := 0; ; i++ {
			p = step(0)
			seen = append(seen, p)
			if p == 0 || i > 20 {
				break
			}
		}
		if p != 0 {
			t.Fatalf("edge %d: the pop-out ramp never landed: %v", e, seen)
		}
		mid := 0
		for _, f := range seen {
			if f > 0 && f < 1 {
				mid++
			}
		}
		if mid < 2 {
			t.Fatalf("edge %d: the ball must slide out through a ramp, not one frame: %v", e, seen)
		}
		if b.DebugTimersAlive() {
			t.Fatalf("edge %d: the dock ramp armed a timer", e)
		}
		// Fully out means fully visible: the whole orb inside the work area.
		wr, m = liveRect()
		x0, x1 := wr.L+half-int(R), wr.L+half+int(R)
		y0, y1 := wr.T+half-int(R), wr.T+half+int(R)
		if x0 < m.Work.L-2 || x1 > m.Work.R+2 || y0 < m.Work.T-2 || y1 > m.Work.B+2 {
			t.Fatalf("edge %d: the popped-out orb is not fully visible: rect %+v work %+v", e, wr, m.Work)
		}
		// And the pointer leaving walks it back into a tab.
		for i := 0; ; i++ {
			p = step(1)
			if p == 1 || i > 20 {
				break
			}
		}
		if p != 1 {
			t.Fatalf("edge %d: the retreat never re-docked: %v", e, p)
		}
		if b.DebugTimersAlive() {
			t.Fatalf("edge %d: the dock owns a timer", e)
		}
	}

	// Dragged back to the middle of the work area, the commit must free the
	// ball: no edge is nearby, so no dock and no position authority. (Releasing
	// AT an edge re-docks by design, which is what the loop above proved.)
	b.sta.PostTask(func() {
		work, edgePx, _, _, _, _, _, ok := b.dockGeometry()
		if !ok {
			t.Error("no dock geometry for the live window")
			return
		}
		x := work.L + (work.R-work.L-int(edgePx))/2
		y := work.T + (work.B-work.T-int(edgePx))/2
		pSetWindowPos.Call(uintptr(b.hwnd), 0, uintptr(x), uintptr(y), 0, 0,
			uintptr(swpNoSize|swpNoZOrder|swpNoActivate))
		b.dock.edge, b.dock.p, b.dock.owns = EdgeNone, 0, false
		b.dockCommit()
	})
	if e, _, owns := b.Docked(); e != EdgeNone || owns {
		t.Fatalf("after a commit with no edge nearby: edge=%d owns=%v", e, owns)
	}
	b.SetState(statemachine.StateSleeping)
	if b.DebugTimersAlive() {
		t.Fatal("a docked-then-freed Sleeping orb holds a timer")
	}
}

var (
	procSetCursorPos = modUser32.NewProc("SetCursorPos")
	procGetCursorPos = modUser32.NewProc("GetCursorPos")
)

// TestBallLiveEdgeDockHover is the end-to-end half of "悬停 → 弹回完整球": a real
// pointer walks onto a docked tab and the ball slides out, then walks off and
// the tab closes again. Every frame here came from a mouse message - the test
// asserts no timer was ever armed while the orb sat in Sleeping the whole time.
// If some other window owns the pixel the pointer is sent to, the hover cannot
// be delivered and the test says so instead of failing on an occlusion - LOUDLY,
// and only about the wiring: the pop-back behavior itself is pinned by the
// non-skipping TestDockHoverPopBackWalksTheRampHome in dock_test.go, so a skip
// here never leaves the direction unproven.
func TestBallLiveEdgeDockHover(t *testing.T) {
	restored := PrototypeVisualsEnabled()
	t.Cleanup(func() { EnablePrototypeVisuals(restored) })
	EnablePrototypeVisuals(true)

	b, err := New(Options{Initial: statemachine.StateSleeping})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer b.Close()

	var prev point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&prev)))
	t.Cleanup(func() { procSetCursorPos.Call(uintptr(prev.x), uintptr(prev.y)) })
	// A background process can be refused cursor input; that is an environment
	// limit, not a product defect, so say so instead of failing.
	moveCursor := func(x, y int32) {
		procSetCursorPos.Call(uintptr(x), uintptr(y))
		var cp point
		procGetCursorPos.Call(uintptr(unsafe.Pointer(&cp)))
		if cp.x != x || cp.y != y {
			t.Skipf("SKIP-LOUD: SetCursorPos refused (the pointer is owned elsewhere): asked (%d,%d), at (%d,%d). "+
				"This run therefore did NOT exercise the real WM_MOUSEMOVE/WM_MOUSELEAVE path. "+
				"The pop-back BEHAVIOR is proven by the non-skipping deterministic test "+
				"TestDockHoverPopBackWalksTheRampHome in dock_test.go, which drives the same "+
				"dockRampFrame step law the hover and leave handlers call; this test only proves the "+
				"WIRING delivers the event, and it proved nothing this run.",
				x, y, cp.x, cp.y)
		}
	}

	if !b.DebugDock(EdgeRight) {
		t.Fatalf("the ball would not dock on the right edge")
	}
	// MoveWindow places with HWND_TOP, which can drop a TOPMOST window below
	// another app's always-on-top surface: for a hover test the tab has to be
	// the thing under the pointer, so put the topmost bit back.
	const hwndTopMost = uintptr(0xFFFFFFFFFFFFFFFE) // (HWND)(LONG_PTR)-2
	pSetWindowPos.Call(uintptr(b.hwnd), hwndTopMost, 0, 0, 0, 0,
		uintptr(swpNoMove|swpNoSize|swpNoActivate))
	var wr rect
	pGetWindowRect.Call(uintptr(b.hwnd), unsafePtr(&wr))
	m, ok := monitorRectsForWindow(b.hwnd)
	if !ok {
		t.Fatal("no monitor")
	}
	half := int(wr.width() / 2)
	// The docked orb's centre is on-screen (that is the point of the position
	// law); parking the pointer there is what a hover on the tab looks like.
	cx, cy := int(wr.l)+half, int(wr.t)+half
	moveCursor(int32(m.Work.L+40), int32(m.Work.T+40)) // away first
	time.Sleep(150 * time.Millisecond)

	settled := func() float32 {
		done := make(chan float32, 1)
		b.sta.PostTask(func() { done <- b.dock.p })
		return <-done
	}
	hovered := func() bool {
		done := make(chan bool, 1)
		b.sta.PostTask(func() { done <- b.dock.hover })
		return <-done
	}
	waitFor := func(cond func() bool, budget time.Duration) bool {
		ctx, cancel := context.WithTimeout(context.Background(), budget)
		defer cancel()
		for {
			if cond() {
				return true
			}
			if ctx.Err() != nil {
				return false
			}
			time.Sleep(30 * time.Millisecond)
		}
	}
	// postMsg runs the wndproc route for one mouse message with the pointer in
	// the tab's client area - the fallback when the session cannot own the
	// physical cursor.
	postMsg := func(msg uint32) {
		done := make(chan struct{})
		b.sta.PostTask(func() {
			lParam := uintptr(uint32(uint16(36)) | uint32(uint16(36))<<16)
			pPostMessageW.Call(uintptr(b.hwnd), uintptr(msg), 0, lParam)
			close(done)
		})
		<-done
		time.Sleep(60 * time.Millisecond)
	}

	// One WM_MOUSEMOVE on the tab must START the pop-out: the hover flag flips
	// and the ramp leaves its docked end. (A human keeps moving the mouse, so
	// the rest of the ramp arrives with the rest of the messages; with the 33ms
	// frame budget a single move carries a single step.)
	moveCursor(int32(cx), int32(cy))
	started := waitFor(func() bool { return hovered() && settled() < 1 }, 2*time.Second)
	via := "cursor"
	if !started && !hovered() {
		via = "posted WM_MOUSEMOVE"
		postMsg(wmMouseMove)
		started = waitFor(func() bool { return hovered() && settled() < 1 }, 2*time.Second)
	}
	if !started {
		t.Skipf("SKIP-LOUD: hover undeliverable this run (neither the physical cursor nor a posted WM_MOUSEMOVE "+
			"reached the tab): centre (%d,%d), hovered=%v p=%v. This run therefore did NOT exercise the real "+
			"WM_MOUSEMOVE/WM_MOUSELEAVE path. The pop-back BEHAVIOR is still proven, by the non-skipping "+
			"deterministic test TestDockHoverPopBackWalksTheRampHome in dock_test.go, which drives the same "+
			"dockRampFrame step law these handlers call; what this test could not show today is the wiring.",
			cx, cy, hovered(), settled())
	}
	if b.DebugTimersAlive() {
		t.Fatal("the hover pop-out armed a timer in Sleeping")
	}
	done := make(chan float32, 1)
	step := func(target float32) float32 {
		b.sta.PostTask(func() {
			b.dock.target = target
			b.dockStep(40 * time.Millisecond)
			done <- b.dock.p
		})
		return <-done
	}
	for i := 0; step(0) != 0 && i < 20; i++ {
	}
	if p := settled(); p != 0 {
		t.Fatalf("the %s pop-out never landed (p=%v)", via, p)
	}
	t.Logf("pop-out driven by %s", via)

	// Pointer away: WM_MOUSELEAVE closes the tab again.
	moveCursor(int32(m.Work.L+40), int32(m.Work.T+40))
	closed := waitFor(func() bool { return settled() == 1 }, 2*time.Second)
	if !closed {
		postMsg(wmMouseLeave)
		closed = waitFor(func() bool { return settled() == 1 }, 2*time.Second)
	}
	if !closed {
		t.Fatalf("the pointer left but the tab never closed (p=%v)", settled())
	}
	if b.DebugTimersAlive() {
		t.Fatal("a docked Sleeping orb holds a timer")
	}
}

// must animate session states and must NEVER touch Sleeping's zero-timer
// discipline. Posts to the UI thread are FIFO, so every assertion below is
// ordered against the levels it just pushed, not racing them.
// TestBallLiveAudioLiquidGate (ticket 62 wiring): the render-side audio seam must
// animate session states and must NEVER touch Sleeping's zero-timer
// discipline. Posts to the UI thread are FIFO, so every assertion below is
// ordered against the levels it just pushed, not racing them.
func TestBallLiveAudioLiquidGate(t *testing.T) {
	restored := PrototypeVisualsEnabled()
	t.Cleanup(func() { EnablePrototypeVisuals(restored) })

	b, err := New(Options{Initial: statemachine.StateSleeping, StartHidden: true})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer b.Close()

	// Frozen mode first: the seam is inert, levels cannot arm anything.
	for i := 0; i < 4; i++ {
		b.SetAudioLevel(0.9)
	}
	if b.DebugTimersAlive() {
		t.Fatal("frozen mode: a level armed a timer")
	}

	EnablePrototypeVisuals(true)

	// Sleeping + a hot mic: levels are dropped on the floor (D32).
	b.SetState(statemachine.StateSleeping)
	for i := 0; i < 6; i++ {
		b.SetAudioLevel(0.9)
	}
	if b.DebugTimersAlive() {
		t.Fatal("audio level armed a timer in Sleeping")
	}

	// Listening: the summon burst must drive the bounded liquid timer...
	b.SetState(statemachine.StateListening)
	b.SetAudioLevel(0.9)
	liquid := make(chan bool, 1)
	b.sta.PostTask(func() { liquid <- b.liquidTimerActive })
	if !<-liquid {
		t.Fatal("a summoned orb in Listening must run the liquid burst timer")
	}

	// ...and it must retire by itself: the model's burst ends, and with the
	// raw feed above the silence gate the border never needs it after.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			t.Fatal("liquid timer outlived its bounded burst (must retire)")
		default:
		}
		b.sta.PostTask(func() { liquid <- b.liquidTimerActive })
		if !<-liquid {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Back to Sleeping: motion parked, no timer survives the state edge.
	b.SetState(statemachine.StateSleeping)
	if b.DebugTimersAlive() {
		t.Fatal("liquid timer survived the return to Sleeping")
	}
	b.SetAudioLevel(0.9)
	if b.DebugTimersAlive() {
		t.Fatal("post-Sleeping level feed re-armed a timer")
	}
}
