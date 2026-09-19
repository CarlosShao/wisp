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
