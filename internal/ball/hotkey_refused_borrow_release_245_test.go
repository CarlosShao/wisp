//go:build windows

package ball

// Ticket 245 AC#7: a borrow Win32 REFUSED must not leave the report describing a
// claim that never existed.
//
// The defect this file judges (ticket 245 AC#7, V-F1): TakeEscForCancel's failure
// branch writes the refusal onto the report but never sets escTakenOver, so
// ReleaseEscAfterSession - whose first line returns when that flag is false -
// never runs its own confirmation step. The ball therefore ends the card still
// reading cancel=HotkeyError, AllLive() forever false and Problems() saying the
// key "was not registered", while the desktop is in fact holding nothing at all:
// the refusal is exactly the case where NOTHING was bound. Two symptoms, one
// stale bookkeeping line.
//
// The ruling keeps the refusal VISIBLE for as long as a card is waiting (ticket
// 260 AC#2's receipt is not being taken away here); what this ruler demands is
// that the return step puts the slot back to the idle truth once the session
// ends, the same way it does after a borrow that worked.
//
// Reachability is the same windowless door 260-r5 measured (and this file
// re-measures rather than assumes): a Ball whose ui-sta thread is this test
// goroutine, b.hwnd naming a handle this process does not own, so RegisterHotKey
// refuses with ERROR_INVALID_WINDOW_HANDLE before anything is lent. No t.Skip
// anywhere, no desktop taken: the planted cancel binding is a modifier
// combination, never the bare Esc.

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"syscall"
	"testing"

	"golang.org/x/sys/windows"
)

const (
	// t245PlantHWND is a USER handle value this process does not own. Not 0
	// (RegisterHotKey with a NULL owner really does register desktop-wide), not
	// 0xFFFFFFFF (HWND_MESSAGE is a real special window).
	t245PlantHWND = uintptr(0x001D0A5E)

	// t245CancelBinding is a modifier combination: even if the plant were
	// somehow accepted, no bare key could leave the desktop.
	t245CancelBinding = "Ctrl+Alt+F8"

	// t245WerrInvalidWindowHandle is raw Win32 1400, spelled as a number so this
	// file's expectation does not move with the constant the code under test uses.
	t245WerrInvalidWindowHandle = 1400
)

// t245ErrNoWindowDoor stops staThread.start before the pump: the create step
// returns an error, the thread records the error and leaves (sta_windows.go), so
// no window and no message loop exist for the length of this rig.
var t245ErrNoWindowDoor = errors.New("245 rig: this door creates no window and runs no pump")

// t245IsWindow is this file's own handle to user32's IsWindow (pIsWindow belongs
// to the winlive package layer and is not compiled in the default suite).
var t245IsWindow = modUser32.NewProc("IsWindow")

// t245Mode says which shape of session this rig stands up.
type t245Mode int

const (
	// t245RefusedNoRelease: the borrow is refused and the card is still waiting.
	// The control leg: the refusal has to stay named here.
	t245RefusedNoRelease t245Mode = iota
	// t245RefusedThenRelease: the borrow is refused, then the session ends. This
	// is AC#7's ruler.
	t245RefusedThenRelease
	// t245HeldThenRelease: the ball already holds a borrow that worked, then the
	// session ends. Control for the return path that existed before AC#7's fix,
	// so this fix cannot buy the new behaviour by breaking the old one.
	t245HeldThenRelease
)

// t245Reading is everything measured on the ui-sta thread, carried out for the
// test goroutine to judge (the ball lives on that thread; nothing is read from
// here concurrently).
type t245Reading struct {
	plantIsWin     uintptr
	inline         bool
	entered        bool
	seed           HotkeyReport
	seedProbs      []string
	seedLive       map[uint32]Accelerator
	cancelSeed     HotkeyBinding
	attemptLine    HotkeyBinding // cancel line right after the borrow step
	attemptProbs   []string
	attemptAllLive bool
	attemptTaken   bool
	attemptLive    map[uint32]Accelerator
	final          HotkeyReport
	finalProbs     []string
	finalLive      map[uint32]Accelerator
	finalTaken     bool
	secondRelease  HotkeyReport // a second return step must be a no-op
	doorErr        error
	staHWND        uintptr
}

// t245Run stands the windowless ui-sta door up around one ball and walks the
// shipped TakeEscForCancel / ReleaseEscAfterSession in the requested shape.
func t245Run(t *testing.T, mode t245Mode) t245Reading {
	t.Helper()

	cfg := ApplyHotkeyDefaults(HotkeyConfig{Cancel: t245CancelBinding})

	// The seed is the healthy idle report production builds at boot: three live
	// slots plus ticket 245's standby cancel line, through the deterministic fake
	// registry. Nothing in it is a failure line.
	seed := registerAllWith(0, cfg, (&fakeRegistry{}).register)
	seedCancel, ok := seed.Binding(hkCancel)
	if !ok || seedCancel.Status != HotkeyStandby || seedCancel.Err != nil {
		t.Fatalf("the rig's idle seed cancel line = %+v (present=%t), want standby with no error", seedCancel, ok)
	}
	if p := seed.Problems(); len(p) != 0 {
		t.Fatalf("the rig's idle seed already complains (%q): the readings below would judge the fixture", p)
	}

	b := &Ball{}
	s := newSTAThread(nil)
	b.sta = s
	b.boundCfg = cfg
	b.cancelBinding = cfg.Cancel
	b.registeredHotkeys = seed.Live()
	b.hotkeyReport = seed
	if mode == t245HeldThenRelease {
		// The shape a borrow that WORKED leaves behind, and the flag the return
		// step looks for. A precondition, not an assertion target for that leg.
		b.escTakenOver = true
		b.hotkeyReport = seed.withCancel(cancelBorrowedLineFor(resolveCancelBorrow(cfg.Cancel)))
		b.registeredHotkeys = b.hotkeyReport.Live()
	}

	var o t245Reading
	o.seed = seed
	o.seedProbs = seed.Problems()
	o.seedLive = seed.Live()
	o.cancelSeed = seedCancel

	s.start(func(st *staThread) error {
		o.plantIsWin, _, _ = t245IsWindow.Call(t245PlantHWND)
		o.inline = windows.GetCurrentThreadId() == st.threadID()
		b.hwnd = windows.HWND(t245PlantHWND)
		o.entered = true

		if mode != t245HeldThenRelease {
			b.TakeEscForCancel()
		}
		o.attemptLine, _ = b.hotkeyReport.Binding(hkCancel)
		o.attemptProbs = b.hotkeyReport.Problems()
		o.attemptAllLive = b.hotkeyReport.AllLive()
		o.attemptTaken = b.escTakenOver
		o.attemptLive = b.registeredHotkeys

		if mode != t245RefusedNoRelease {
			b.ReleaseEscAfterSession()
			o.final = b.hotkeyReport
			o.finalProbs = b.hotkeyReport.Problems()
			o.finalLive = b.registeredHotkeys
			o.finalTaken = b.escTakenOver
			// Idempotency: the host may sync state twice. A second return step
			// must not move anything - in particular it must not turn a healthy
			// standby line back into a failure or re-register a key.
			b.ReleaseEscAfterSession()
			o.secondRelease = b.hotkeyReport
		}
		return t245ErrNoWindowDoor
	})
	o.doorErr = s.waitStarted()
	s.mu.Lock()
	o.staHWND = uintptr(s.hwnd)
	s.mu.Unlock()
	return o
}

// t245RequireRefusal is the shared precondition: this rig only judges a refusal
// if a refusal is what actually happened.
func t245RequireRefusal(t *testing.T, o t245Reading) {
	t.Helper()
	if o.plantIsWin != 0 {
		t.Fatalf("IsWindow(0x%X) says the planted handle is a live window, so the borrow could have succeeded", t245PlantHWND)
	}
	if !o.entered {
		t.Fatal("the uiRun closure never ran, so the borrow wiring was never reached")
	}
	if !o.inline {
		t.Fatal("Ball.uiRun did not take its inline branch; with no window the task would be dropped (sta_windows.go)")
	}
	if o.staHWND != 0 {
		t.Errorf("the rig's staThread recorded a window handle 0x%X - this case promised to create no window", o.staHWND)
	}
	if !errors.Is(o.doorErr, t245ErrNoWindowDoor) {
		t.Errorf("the create-error door was not taken (err=%v): a pump running here would make the readings unreadable", o.doorErr)
	}
	if o.attemptTaken {
		t.Fatal("escTakenOver reads true after a refused borrow: the ball claims a key Win32 never lent")
	}
	if o.attemptLine.Status != HotkeyError || o.attemptLine.Err == nil {
		t.Fatalf("the borrow was not refused as this rig requires: cancel line = %+v", o.attemptLine)
	}
	var errno syscall.Errno
	if !errors.As(o.attemptLine.Err, &errno) || errno != syscall.Errno(t245WerrInvalidWindowHandle) {
		t.Fatalf("the refusal carries %v (%T), want Win32 %d ERROR_INVALID_WINDOW_HANDLE",
			o.attemptLine.Err, o.attemptLine.Err, t245WerrInvalidWindowHandle)
	}
	if _, live := o.attemptLive[hkCancel]; live {
		t.Fatalf("the registration set claims the cancel key after a refusal: %v", o.attemptLive)
	}
}

// TestRefusedBorrowReleaseReturnsTheIdleRoster245 is AC#7's ruler. After a borrow
// Win32 refused and the session ended, the report must describe the truth of the
// desktop - three idle slots, cancel on standby (not bound), no complaint - and
// not the stale refusal line that used to stay there for the life of the process.
//
// RED before the fix: ReleaseEscAfterSession returns on its first line because
// escTakenOver is false, so the HotkeyError line written by the refusal stays,
// AllLive() stays false and Problems() keeps saying "was not registered".
func TestRefusedBorrowReleaseReturnsTheIdleRoster245(t *testing.T) {
	o := t245Run(t, t245RefusedThenRelease)
	t245RequireRefusal(t, o)

	// AC#7's judgement, in the order the user meets it.
	cancel, ok := o.final.Binding(hkCancel)
	if !ok {
		t.Fatal("the report has no cancel line at all after the session ended")
	}
	if cancel.Status != HotkeyStandby {
		t.Errorf("after the session ended the cancel line reads %q (err=%v), want %q: nothing was ever bound, so the "+
			"report must say the slot is on standby, not that a registration failed",
			cancel.Status, cancel.Err, HotkeyStandby)
	}
	if cancel.Err != nil {
		t.Errorf("the idle cancel line carries a stale error after the session ended: %v", cancel.Err)
	}
	if cancel.Binding != t245CancelBinding {
		t.Errorf("the idle cancel line names %q, want the configured %q: the slot still knows what the user asked for",
			cancel.Binding, t245CancelBinding)
	}
	if got := o.finalProbs; len(got) != 0 {
		t.Errorf("Problems() after the session ended = %q, want none: the refusal was real while the card waited, "+
			"but it is not a standing fact once the session is over", got)
	}
	if !o.final.AllLive() {
		t.Errorf("AllLive() = false after a refused borrow was returned (lines: %s): the desktop holds the idle set, "+
			"so the report has to say so", describe245(o.final.Bindings()))
	}
	// The registration set itself: exactly the three idle slots, and the cancel
	// slot is NOT among them (ticket 245's steady-state rule is not undone by a
	// refused borrow).
	if !reflect.DeepEqual(o.finalLive, o.seedLive) {
		t.Errorf("the registration set after the session ended = %v, want the idle seed %v", o.finalLive, o.seedLive)
	}
	if len(o.finalLive) != 3 {
		t.Errorf("the registration set after the session ended holds %d keys (%v), want the three idle slots",
			len(o.finalLive), o.finalLive)
	}
	if _, live := o.finalLive[hkCancel]; live {
		t.Errorf("the cancel key is in the registration set after a refused borrow was returned: %v", o.finalLive)
	}
	if o.finalTaken {
		t.Error("escTakenOver reads true after the session ended: the ball claims a key it never got")
	}
	// A second return step (the host may sync state twice) changes nothing.
	if !reflect.DeepEqual(o.secondRelease.Bindings(), o.final.Bindings()) {
		t.Errorf("a second ReleaseEscAfterSession moved the report (%s -> %s): the return step is not idempotent",
			describe245(o.final.Bindings()), describe245(o.secondRelease.Bindings()))
	}
	t.Logf("245 AC#7: refused line before release = %+v / %q; after release = %+v / problems=%q live=%v",
		o.attemptLine, o.attemptProbs, cancel, o.finalProbs, o.finalLive)
}

// TestRefusedBorrowStaysNamedWhileTheCardWaits245 is the control on the ruler
// above, and the guard on ticket 260 AC#2's receipt: the refusal must still be
// named loudly while the card is actually waiting. Without this leg, "the report
// reads standby afterwards" could also be produced by never writing the refusal.
func TestRefusedBorrowStaysNamedWhileTheCardWaits245(t *testing.T) {
	o := t245Run(t, t245RefusedNoRelease)
	t245RequireRefusal(t, o)

	probs := o.attemptProbs
	if len(probs) != 1 {
		t.Fatalf("problem lines while a card waits on a refused borrow = %d (%q), want exactly 1 naming it", len(probs), probs)
	}
	for _, want := range []string{"cancel", t245CancelBinding, "was not registered"} {
		if !strings.Contains(probs[0], want) {
			t.Errorf("the refusal line while the card waits does not name %q; it reads: %s", want, probs[0])
		}
	}
	if o.attemptLine.Status != HotkeyError {
		t.Errorf("the cancel line while a card waits reads %q, want %q", o.attemptLine.Status, HotkeyError)
	}
	if o.attemptAllLive {
		t.Error("AllLive() reports true while a card waits with no cancel key: that is the zero-symptom shape")
	}
	t.Logf("245 control: refusal while the card waits = %+v / %q", o.attemptLine, probs)
}

// TestHeldBorrowReleaseStillReturnsStandby245 keeps the pre-existing return path
// pinned: a ball that DOES hold the borrow ends on the same idle line. AC#7 must
// not buy the refused case by weakening the working one.
func TestHeldBorrowReleaseStillReturnsStandby245(t *testing.T) {
	o := t245Run(t, t245HeldThenRelease)

	if !o.entered || !o.inline {
		t.Fatal("the uiRun closure never ran inline, so the return step judged nothing")
	}
	if o.attemptLine.Status != HotkeyLive {
		t.Fatalf("the held-borrow precondition line = %+v, want the live cancel line", o.attemptLine)
	}
	cancel, ok := o.final.Binding(hkCancel)
	if !ok || cancel.Status != HotkeyStandby || cancel.Err != nil {
		t.Errorf("a released borrow ended reading %+v (present=%t), want standby with no error", cancel, ok)
	}
	if len(o.finalLive) != 3 {
		t.Errorf("after releasing a held borrow the registration set holds %d keys (%v), want the three idle slots",
			len(o.finalLive), o.finalLive)
	}
	if o.finalTaken {
		t.Error("escTakenOver reads true after the borrow was returned")
	}
	if got := o.finalProbs; len(got) != 0 {
		t.Errorf("Problems() after returning a held borrow = %q, want none", got)
	}
}

// describe245 renders a report without leaning on %v of an unexported slice.
func describe245(bindings []HotkeyBinding) string {
	parts := make([]string, 0, len(bindings))
	for _, b := range bindings {
		parts = append(parts, fmt.Sprintf("%s=%q", b.Name, b.Status))
	}
	return strings.Join(parts, " ")
}
