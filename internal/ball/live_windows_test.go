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

var (
	procFindWindowW = modUser32.NewProc("FindWindowW")
)

func findBallWindow() uintptr {
	cls := utf16("WispBallWindow")
	h, _, _ := procFindWindowW.Call(unsafePtr(cls), 0)
	return h
}

func handlesOfProcess() uint32 {
	var n uint32
	procGetProcessHandleCount.Call(uintptr(windows.CurrentProcess()), uintptr(unsafe.Pointer(&n)))
	return n
}

var procGetProcessHandleCount = modKernel32.NewProc("GetProcessHandleCount")

// TestBallLiveLifecycle walks all 20 states on a real window and tears down.
func TestBallLiveLifecycle(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "ballpos.json")
	store, err := OpenPositionStore(storePath)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	baseHandles := handlesOfProcess()

	b, err := New(Options{
		SizePx:  BallSizeDefaultPx,
		Initial: statemachine.StateSleeping,
		Store:   store,
		Hotkeys: HotkeyConfig{Summon: "Ctrl+Alt+F9", Mute: "Ctrl+Alt+F10", Cancel: "Ctrl+Alt+F11", Panel: "Ctrl+Alt+F12"},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer func() {
		b.Close()
		if h := findBallWindow(); h != 0 {
			t.Errorf("ball window still alive after Close")
		}
	}()

	if findBallWindow() == 0 {
		t.Fatal("ball window not found after New")
	}

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

	b, err := New(Options{Initial: statemachine.StateSleeping, Store: store})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
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
	time.Sleep(100 * time.Millisecond)

	got, ok := store.Get(dev)
	if !ok || got != want {
		t.Fatalf("persisted position: got %+v ok=%v, want %+v (device %s)", got, ok, want, dev)
	}

	// Re-open: the restored top-left must equal the saved position.
	b2, err := New(Options{Initial: statemachine.StateSleeping, Store: store})
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer b2.Close()
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
	for _, s := range []statemachine.State{statemachine.StateArmed, statemachine.StateMuted,
		statemachine.StateConversation, statemachine.StateError} {
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

// TestBallLiveAudioLiquidGate (ticket 62 wiring): the render-side audio seam
// must animate session states and must NEVER touch Sleeping's zero-timer
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
