//go:build windows && winlive

package ball

// The four SPEC-08 §2 interaction clauses, tested as BEHAVIOUR on a real
// window (ticket 64 A2). Before this file only HitTest()'s pure geometric
// partition had tests; the clauses themselves - "click summons", "Esc/click
// cancels from Confirming", "never takes focus", "transparent pixels fall
// through" - had none.
//
// Each test drives the real window procedure through posted/sent Win32
// messages or injected input, and asserts on what the machine and the window
// ended up as. Where a clause needs a human hand that no harness can supply
// (a second monitor, the felt experience of focus) the test says so in this
// file rather than faking a proof: see the notes at each clause.

import (
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

const (
	gwlExStyle     = ^uintptr(19) // GWL_EXSTYLE (-20)
	htTransparent  = ^uintptr(0)  // LRESULT(-1)
	maNoActivate   = uintptr(3)
	htClientResult = uintptr(HTClient)
	mouseLeftDown  = 0x0002
	mouseLeftUp    = 0x0004
)

var procMouseEvent = modUser32.NewProc("mouse_event") // MOUSEEVENTF_LEFT* below

// pointParam packs screen coordinates the way WM_NCHITTEST carries them in
// lParam: x in the low WORD, y in the high WORD.
func pointParam(x, y int32) uintptr {
	return uintptr(uint32(uint16(x)) | uint32(uint16(y))<<16)
}

// pointReg packs a POINT for an API that takes the STRUCT by value. On x64 an
// 8-byte POINT goes in one register: x in the low DWORD, y in the HIGH DWORD -
// which is NOT the lParam packing above (mixing the two makes WindowFromPoint
// look up a garbage point and return NULL).
func pointReg(x, y int32) uintptr {
	return uintptr(uint32(x)) | uintptr(uint32(y))<<32
}

// TestLiveClickSummonsAndDragDoesNot covers clause 1 (单击球 Sleeping ->
// Listening) and the drag/click discrimination that clause rests on.
func TestLiveClickSummonsAndDragDoesNot(t *testing.T) {
	m := statemachine.New(statemachine.Options{Initial: statemachine.StateSleeping})
	var b *Ball
	var mu sync.Mutex
	var clicks, drags int
	b = newLiveBall(t, Options{
		Initial: statemachine.StateSleeping,
		Hotkeys: liveHotkeys(),
		Events: Events{
			OnClickBall: func() {
				mu.Lock()
				clicks++
				mu.Unlock()
				if _, err := m.Dispatch(statemachine.EvSummon, nil); err != nil {
					t.Errorf("dispatch EvSummon: %v", err)
					return
				}
				b.SetState(m.State())
			},
			OnDragEnd: func() {
				mu.Lock()
				drags++
				mu.Unlock()
			},
		},
	})
	requireFreeOfDock(t, b)

	hwnd := b.DebugHWND()
	var wr rect
	pGetWindowRect.Call(uintptr(hwnd), unsafePtr(&wr))
	client := lParamPoint(wr.width()/2, wr.height()/2)

	// A drag past the 4px threshold moves the window and fires OnDragEnd, and
	// must NOT summon. Done first so the ball starts from Sleeping.
	before := windowRect(t, b)
	sendMouse(b, wmLButtonDown, 0, client)
	sendMouse(b, wmMouseMove, 0, lParamPoint(wr.width()/2+40, wr.height()/2+30))
	sendMouse(b, wmLButtonUp, 0, lParamPoint(wr.width()/2+40, wr.height()/2+30))
	mu.Lock()
	c0, d0 := clicks, drags
	mu.Unlock()
	if c0 != 0 {
		t.Fatalf("a drag fired OnClickBall (%d times), which would summon the ball on every release", c0)
	}
	if d0 != 1 {
		t.Fatalf("a drag fired OnDragEnd %d times, want 1 (position must persist)", d0)
	}
	if after := windowRect(t, b); after.L == before.L && after.T == before.T {
		t.Fatalf("the drag did not move the window: %+v -> %+v", before, after)
	}
	if got := readState(t, b); got != statemachine.StateSleeping {
		t.Fatalf("the drag left the ball in %s, want the untouched Sleeping", got)
	}
	if m.State() != statemachine.StateSleeping {
		t.Fatalf("the drag moved the machine to %s, want Sleeping", m.State())
	}

	// Clause 1: a single click on the orb wakes it - Sleeping -> Listening
	// (D43 #4), on the machine AND on the rendered ball.
	pGetWindowRect.Call(uintptr(hwnd), unsafePtr(&wr))
	client = lParamPoint(wr.width()/2, wr.height()/2)
	sendMouse(b, wmLButtonDown, 0, client)
	sendMouse(b, wmLButtonUp, 0, client)
	if got := readState(t, b); got != statemachine.StateListening {
		t.Fatalf("a single click on the orb left the ball rendering %s, want Listening (D43 #4)", got)
	}
	if m.State() != statemachine.StateListening {
		t.Fatalf("machine after the click = %s, want Listening", m.State())
	}
	mu.Lock()
	c1, d1 := clicks, drags
	mu.Unlock()
	if c1 != 1 || d1 != 1 {
		t.Fatalf("click=%d drag=%d after one click, want 1/1 (a click must not also read as a drag)", c1, d1)
	}
}

// TestLiveConfirmingCancelAndEscReturned covers clause 2 (B1: cancel from
// Confirming, and the Esc takeover handed back at session end).
func TestLiveConfirmingCancelAndEscReturned(t *testing.T) {
	m := statemachine.New(statemachine.Options{Initial: statemachine.StateSleeping})
	var b *Ball
	var mu sync.Mutex
	var cancels int
	syncBall := func() {
		b.SetState(m.State())
		if m.State() == statemachine.StateConfirming {
			b.TakeEscForCancel() // B1: Esc is the cancel key while Confirming
		} else {
			b.ReleaseEscAfterSession() // ...and is handed back the moment it is not
		}
	}
	b = newLiveBall(t, Options{
		Initial: statemachine.StateSleeping,
		Hotkeys: liveHotkeys(),
		Events: Events{
			OnCancelHotkey: func() {
				mu.Lock()
				cancels++
				mu.Unlock()
				if _, err := m.Dispatch(statemachine.EvVeto, nil); err != nil {
					t.Errorf("dispatch EvVeto: %v", err)
					return
				}
				syncBall()
			},
			OnClickBall: func() {
				if m.State() != statemachine.StateListening && m.State() != statemachine.StateConfirming {
					return
				}
				mu.Lock()
				cancels++
				mu.Unlock()
				if _, err := m.Dispatch(statemachine.EvVeto, nil); err != nil {
					t.Errorf("click veto: %v", err)
					return
				}
				syncBall()
			},
		},
	})

	// Walk D43 to Confirming: #4 summon, #11 vad-stop, #15 first token with a
	// tool call, #17 approval needed at L1.
	steps := []struct {
		ev    statemachine.Event
		facts *statemachine.Facts
		want  statemachine.State
	}{
		{statemachine.EvSummon, nil, statemachine.StateListening},
		{statemachine.EvVadStop, &statemachine.Facts{SpeechMS: 900}, statemachine.StateThinking},
		{statemachine.EvFirstToken, &statemachine.Facts{HasToolCall: true}, statemachine.StateActing},
		{statemachine.EvApprovalNeeded, &statemachine.Facts{ApprovalLevel: 1}, statemachine.StateConfirming},
	}
	for _, s := range steps {
		if _, err := m.Dispatch(s.ev, s.facts); err != nil {
			t.Fatalf("dispatch %s: %v", s.ev, err)
		}
		if m.State() != s.want {
			t.Fatalf("after %s the machine is %s, want %s", s.ev, m.State(), s.want)
		}
	}
	syncBall()
	if !waitFor(2*time.Second, func() bool { return b.EscTakenOver() }) {
		t.Fatal("entering Confirming did not take Esc over (B1)")
	}
	if got := readState(t, b); got != statemachine.StateConfirming {
		t.Fatalf("ball renders %s, want Confirming", got)
	}

	// The cancel binding (Esc while taken over) vetoes: D43 #22 -> Acting.
	sendMessage(b, wmHotkey, hkCancel, 0)
	if !waitFor(2*time.Second, func() bool { return m.State() == statemachine.StateActing }) {
		t.Fatalf("the cancel hotkey in Confirming left the machine in %s, want Acting (D43 #22)", m.State())
	}
	if b.EscTakenOver() {
		t.Fatal("Esc still taken over outside Confirming: B1 return violated")
	}
	// Esc itself must be ours to give back: while Confirming, an injected bare
	// Esc is the user's cancel.
	if _, err := m.Dispatch(statemachine.EvApprovalNeeded, &statemachine.Facts{ApprovalLevel: 1}); err != nil {
		t.Fatalf("back into Confirming: %v", err)
	}
	syncBall()
	if !b.EscTakenOver() {
		t.Fatal("re-entering Confirming did not take Esc over")
	}
	if !injectBinding(t, "Esc") {
		t.Skipf("SKIP-LOUD: injected Esc never reached the input stream, so this run did NOT prove the " +
			"physical bare-Esc cancel during Confirming. The takeover itself, the WM_HOTKEY routing of the " +
			"cancel id and the D43 #22 veto were asserted above; the B1 return is asserted below by the real " +
			"registration of the configured binding.")
	}
	if !waitFor(2*time.Second, func() bool { return m.State() == statemachine.StateActing }) {
		t.Fatalf("the taken-over Esc did not cancel Confirming (machine %s)", m.State())
	}
	syncBall()
	if b.EscTakenOver() {
		t.Fatal("Esc not returned after the session")
	}
	if rep := b.HotkeyReport(); !rep.IsLive(hkCancel) {
		t.Errorf("after the B1 return the configured cancel binding is not live: %+v", rep.Bindings())
	}

	// A single click cancels too (the same veto event, the gesture table says).
	if _, err := m.Dispatch(statemachine.EvApprovalNeeded, &statemachine.Facts{ApprovalLevel: 1}); err != nil {
		t.Fatalf("into Confirming again: %v", err)
	}
	syncBall()
	requireFreeOfDock(t, b)
	var wr rect
	pGetWindowRect.Call(uintptr(b.DebugHWND()), unsafePtr(&wr))
	client := lParamPoint(wr.width()/2, wr.height()/2)
	sendMouse(b, wmLButtonDown, 0, client)
	sendMouse(b, wmLButtonUp, 0, client)
	if !waitFor(2*time.Second, func() bool { return m.State() == statemachine.StateActing }) {
		t.Fatalf("a click in Confirming left the machine in %s, want the veto to Acting", m.State())
	}
	syncBall()
	if b.EscTakenOver() {
		t.Error("click-cancel left Esc taken over")
	}
}

// TestLiveNeverStealsFocus covers clause 3. What is asserted here is the
// measurable half of "does not take focus": the window carries
// WS_EX_NOACTIVATE, its WM_MOUSEACTIVATE answer is MA_NOACTIVATE, the ball
// thread never owns focus, and the foreground window is the same handle after
// a real injected click as before it. The half no harness can assert - that a
// human does not FEEL the focus move - is the owner sign-off, not this test.
func TestLiveNeverStealsFocus(t *testing.T) {
	b := newLiveBall(t, Options{Initial: statemachine.StateListening, Hotkeys: liveHotkeys()})
	requireFreeOfDock(t, b)
	hwnd := b.DebugHWND()
	if !windowAlive(hwnd) {
		t.Fatal("the ball window is not alive")
	}

	// 1. The style bit the OS enforces.
	ex, _, _ := pGetWindowLongPtrW.Call(uintptr(hwnd), gwlExStyle)
	if ex&uintptr(wsExNoActivate) == 0 {
		t.Fatalf("the live window does not carry WS_EX_NOACTIVATE: exStyle=0x%X", ex)
	}
	if ex&uintptr(wsExLayered) == 0 {
		t.Errorf("the live window lost WS_EX_LAYERED: exStyle=0x%X", ex)
	}

	// 2. The message answer, for the mouse button down that would activate:
	//    lParam is MAKELONG(HTCLIENT, WM_LBUTTONDOWN).
	lParam := uintptr(uint32(HTClient) | uint32(wmLButtonDown)<<16)
	if got := sendMessage(b, wmMouseActivate, 0, lParam); got != maNoActivate {
		t.Fatalf("WM_MOUSEACTIVATE returned %d, want MA_NOACTIVATE(3)", got)
	}

	// 3. Focus belongs to somebody else before and after a full click.
	fgBefore, _, _ := pGetForegroundWindow.Call()
	sendMouse(b, wmLButtonDown, 0, lParamPoint(20, 20))
	sendMouse(b, wmLButtonUp, 0, lParamPoint(20, 20))
	if fgBefore != 0 {
		if fgNow, _, _ := pGetForegroundWindow.Call(); fgNow != fgBefore {
			t.Fatalf("the foreground window changed across a ball click: %x -> %x", fgBefore, fgNow)
		}
	}
	if f, _, _ := pGetForegroundWindow.Call(); f == uintptr(hwnd) {
		t.Fatal("the ball became the foreground window")
	}

	// 4. The real thing: move the physical pointer onto the orb and click it.
	var wr rect
	pGetWindowRect.Call(uintptr(hwnd), unsafePtr(&wr))
	cx, cy := int32(wr.l+wr.width()/2), int32(wr.t+wr.height()/2)
	procSetCursorPos.Call(uintptr(cx), uintptr(cy))
	var cp point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&cp)))
	if cp.x != cx || cp.y != cy {
		t.Skipf("SKIP-LOUD: SetCursorPos was refused, so this run did NOT click the orb with the physical "+
			"pointer (asked (%d,%d), pointer at (%d,%d)). The NOACTIVATE style, the MA_NOACTIVATE answer and "+
			"the unchanged foreground across the posted click were all asserted above.", cx, cy, cp.x, cp.y)
	}
	procMouseEvent.Call(mouseLeftDown, 0, 0, 0)
	procMouseEvent.Call(mouseLeftUp, 0, 0, 0)
	time.Sleep(150 * time.Millisecond)
	if f, _, _ := pGetForegroundWindow.Call(); f == uintptr(hwnd) {
		t.Fatalf("a physical click on the orb activated it: foreground is now the ball (%x)", hwnd)
	}
	if fgBefore != 0 {
		if fgNow, _, _ := pGetForegroundWindow.Call(); fgNow != fgBefore {
			t.Fatalf("a physical click moved the foreground: %x -> %x", fgBefore, fgNow)
		}
	}
	if h, _, _ := pWindowFromPoint.Call(pointReg(cx, cy)); h != uintptr(hwnd) {
		t.Fatalf("the pixel under the pointer is not the ball after moving there: hwnd=%x want %x", h, uintptr(hwnd))
	}
}

// TestLiveTransparentCornerFallsThrough covers clause 4 on a real window: the
// window procedure answers HTTRANSPARENT for the corner pixels, and the OS
// agrees that a point there does not belong to the ball, while the centre of
// the orb does.
func TestLiveTransparentCornerFallsThrough(t *testing.T) {
	b := newLiveBall(t, Options{Initial: statemachine.StateListening, Hotkeys: liveHotkeys()})
	requireFreeOfDock(t, b)
	hwnd := b.DebugHWND()
	var wr rect
	pGetWindowRect.Call(uintptr(hwnd), unsafePtr(&wr))
	edge := wr.width()
	cx, cy := wr.l+edge/2, wr.t+edge/2
	nx, ny := wr.l+1, wr.t+1

	if r := sendMessage(b, wmNCHitTest, 0, pointParam(cx, cy)); r != htClientResult {
		t.Fatalf("WM_NCHITTEST at the orb centre (%d,%d) = %d, want HTCLIENT(1)", cx, cy, int32(r))
	}
	if r := sendMessage(b, wmNCHitTest, 0, pointParam(nx, ny)); r != htTransparent {
		t.Fatalf("WM_NCHITTEST at the transparent corner = %d, want HTTRANSPARENT(-1)", r)
	}
	// WindowFromPoint is the OS's own answer to "who gets this click", and for
	// a layered window it also honours the alpha bitmap: the orb centre must
	// resolve to the ball, the transparent corner must not.
	if h, _, _ := pWindowFromPoint.Call(pointReg(cx, cy)); h != uintptr(hwnd) {
		t.Fatalf("WindowFromPoint(centre) = %x, want the ball %x - the orb must be clickable", h, uintptr(hwnd))
	}
	if h, _, _ := pWindowFromPoint.Call(pointReg(nx, ny)); h == uintptr(hwnd) {
		t.Fatalf("WindowFromPoint(corner) still returns the ball: a click on a fully transparent pixel "+
			"would be eaten by it (hwnd=%x)", h)
	}
	// A click on the corner must not change the window's own answer either: the
	// corner keeps belonging to the desktop no matter what was posted at it.
	before := sendMessage(b, wmNCHitTest, 0, pointParam(nx, ny))
	sendMouse(b, wmLButtonDown, 0, lParamPoint(1, 1))
	sendMouse(b, wmLButtonUp, 0, lParamPoint(1, 1))
	if after := sendMessage(b, wmNCHitTest, 0, pointParam(nx, ny)); after != before {
		t.Errorf("the corner's hit-test answer changed after a click: %d -> %d", before, after)
	}
}

// ---------------------------------------------------------------- helpers

// requireFreeOfDock parks the ball in the middle of the work area with no
// dock, so click tests are about clicks and not about the tab. It runs SYNCHRONOUSLY
// on the UI thread: a caller that reads the window rect right after this must
// see the parked position, not the pre-park one (a PostTask here made the
// hit-test assertions race the move).
func requireFreeOfDock(t *testing.T, b *Ball) {
	t.Helper()
	b.uiRun(func() {
		work, edgePx, _, _, _, _, _, ok := b.dockGeometry()
		if !ok {
			t.Errorf("no dock geometry for the live window")
			return
		}
		x := work.L + (work.R-work.L-int(edgePx))/2
		y := work.T + (work.B-work.T-int(edgePx))/2
		pSetWindowPos.Call(uintptr(b.hwnd), 0, uintptr(x), uintptr(y), 0, 0,
			uintptr(swpNoSize|swpNoZOrder|swpNoActivate))
		b.dock.edge, b.dock.p, b.dock.owns = EdgeNone, 0, false
		b.dockCommit()
	})
	if e, _, owns := b.Docked(); owns || e != EdgeNone {
		t.Fatalf("could not free the ball from the dock (edge=%d owns=%v); click tests need a floating orb", e, owns)
	}
}

// windowRect reads the live window rectangle.
func windowRect(t *testing.T, b *Ball) Rect {
	t.Helper()
	var wr rect
	pGetWindowRect.Call(uintptr(b.DebugHWND()), unsafePtr(&wr))
	return Rect{int(wr.l), int(wr.t), int(wr.r), int(wr.b)}
}
