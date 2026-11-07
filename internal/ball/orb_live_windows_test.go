//go:build windows && winlive

package ball

// Ticket-65 integration leg (A1594): the LIVE window test for the WebView2
// orb. It creates the real ball, waits for the real WebView2 runtime to attach
// the orb beside it, and requires the whole takeover to have happened; it also
// verifies the attach contract itself (one fixed edge, topmost, tool window,
// never-activate show) and tears everything down. Run explicitly:
//
//	go test -tags winlive ./internal/ball/ -run TestBallLiveOrb -v
//
// Excluded from the default suite like every live test, so `go test ./...`
// never pops windows. The fallback half of the contract is pinned in the
// default suite: a ball created under -tags winlive WITHOUT a WebView2 runtime
// fails THIS test loudly rather than silently rendering D2D - the fallback is
// a decision the harness reads (OrbDebugLine), never a guess.

import (
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// This file's own handle to user32's GetWindowLongPtrW: pIsWindow belongs to
// live_guard_windows_test.go (in-binary single declaration), and a negative
// GWL_* index cannot be a const uintptr, so it lives in a variable like the
// panel host's ctlGwlStyle does.
var orbTestGetWindowLongPtrW = modUser32.NewProc("GetWindowLongPtrW")

var gwlExStyleTest = -20 // GWL_EXSTYLE

// TestBallLiveOrbAttachStateAndTeardown: real ball + real orb window on the
// real desktop. Assertions ride the test's OWN hwnds (never FindWindow by
// class), the quiet-desktop precondition holds the desktop, and every handle
// this test creates is destroyed before it returns.
func TestBallLiveOrbAttachStateAndTeardown(t *testing.T) {
	store, err := OpenPositionStore(filepath.Join(t.TempDir(), "ballpos.json"))
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	requireQuietBallDesktop(t)
	baseHandles := handlesOfProcess()

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
	t.Cleanup(func() {
		hwnd := b.DebugHWND()
		b.Close()
		if r, _, _ := pIsWindow.Call(uintptr(hwnd)); r != 0 {
			t.Fatalf("ball window survived Close: hwnd=%v", hwnd)
		}
	})

	// The orb attach is asynchronous from New's point of view only in the
	// WebView2 runtime cold start INSIDE attachOrb; the PostTask it rides
	// lands before any later task, and SetState posts after it - so by the
	// time the deadline below expires without a window, the runtime really
	// did not come up.
	deadline := time.Now().Add(10 * time.Second)
	var orbHWND windows.HWND
	for {
		orbHWND = b.DebugOrbHWND()
		if orbHWND != 0 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("orb window did not attach within 10s (WebView2 runtime missing or Embed refused); "+
				"debug face says %q", OrbDebugLine(false, statemachine.StateSleeping))
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Attach contract: the orb window is a real window, of the orb class's
	// process, sized as ONE fixed square, shown without activation, topmost.
	if r, _, _ := pIsWindow.Call(uintptr(orbHWND)); r == 0 {
		t.Fatalf("DebugOrbHWND is not a live window: %v", orbHWND)
	}
	exStyle, _, _ := orbTestGetWindowLongPtrW.Call(uintptr(orbHWND), uintptr(gwlExStyleTest))
	const wantEx = wsExLayered | wsExTopmost | wsExToolWindow | wsExNoActivate
	if exStyle&wantEx != wantEx {
		t.Errorf("orb ex-style %#x missing bits of %#x", exStyle, wantEx)
	}
	var wr rect
	if r, _, _ := pGetWindowRect.Call(uintptr(orbHWND), unsafePtr(&wr)); r == 0 {
		t.Fatal("GetWindowRect(orb) failed")
	}
	if w, h := wr.width(), wr.height(); w != h || w <= 0 {
		t.Errorf("orb window is not a square: %dx%d", w, h)
	}

	// The takeover contract, on the REAL surface: the debug face names the
	// webview2 orb, and a SetState through the ball reaches the orb's page.
	if got := OrbDebugLine(b.DebugOrbAttached(), b.DebugOrbState()); got != "orb-surface=webview2-orb orb-state=idle" {
		t.Errorf("debug face after attach = %q", got)
	}
	b.SetState(statemachine.StateThinking)
	deadline = time.Now().Add(5 * time.Second)
	for {
		if got := OrbDebugLine(b.DebugOrbAttached(), b.DebugOrbState()); got == "orb-surface=webview2-orb orb-state=thinking" {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("orb state push never surfaced; face stuck at %q", OrbDebugLine(b.DebugOrbAttached(), b.DebugOrbState()))
		}
		time.Sleep(50 * time.Millisecond)
	}

	// LEAK sensitivity: repeated pushes must be free. Handles ate ~300 on the
	// attach itself (measured twice on this box, and of the same shape the
	// vendored module imposes on the panel host - it never releases the
	// environment). That one-time-per-process cost is bounded by construction
	// (attachOrb is CAS-guarded to ONE orb per process), so the number that
	// can still grow and hide a leak is the per-push delta: drive the same
	// burst the audio capture would (50 state pushes + 50 audio pushes) and
	// require the handle count to come back.
	afterAttach := handlesOfProcess()
	for i := 0; i < 50; i++ {
		b.SetState(statemachine.StateListening)
		b.SetAudioLevel(0.5)
		b.SetState(statemachine.StateThinking)
	}
	deadline = time.Now().Add(5 * time.Second)
	for {
		// The last push lands on the STA thread; give the queue a moment to
		// drain before reading the counter.
		time.Sleep(100 * time.Millisecond)
		if got := handlesOfProcess(); got <= afterAttach+50 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("handle growth across 100 orb pushes: before=%d after=%d",
				afterAttach, handlesOfProcess())
		}
	}

	// The ball window itself is still alive and still the ball: the orb is a
	// SECOND window, not a replacement.
	if r, _, _ := pIsWindow.Call(uintptr(b.DebugHWND())); r == 0 {
		t.Fatal("ball window died while the orb was attached")
	}

	// Teardown: Close must destroy BOTH windows and stay inside the one-time
	// attach budget (documented above; the window-death checks are the hard
	// assertions, this one only catches unbounded growth).
	b.Close()
	if r, _, _ := pIsWindow.Call(uintptr(orbHWND)); r != 0 {
		t.Fatalf("orb window survived Close: hwnd=%v", orbHWND)
	}
	if got := handlesOfProcess(); got > baseHandles+600 {
		t.Errorf("handle growth beyond the one-time attach budget: base=%d end=%d", baseHandles, got)
	}
}
