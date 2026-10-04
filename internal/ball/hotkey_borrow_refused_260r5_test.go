//go:build windows

package ball

// Ticket 260 AC#2, leg 260-r5: a ruler that can HEAR the borrow being refused.
//
// What the cell asks for, from the ticket: plant "the borrow request went out but
// Win32 never lent the key" and make the step that says so go red - in a
// non-winelive case, with no t.Skip, and NOT in the shape the two existing rulers
// use, where the test itself pastes the failure line back into the report with
// withCancel(cancelFailedLine(...)) and the wiring is never walked.
//
// Why nobody had one (260-v1's measurement, re-run by this leg before writing):
// every case that calls Ball.TakeEscForCancel sits behind //go:build windows &&
// winlive (logs/ruler-takeesc-testcallers.txt and logs/ruler-build-tags.txt under
// .scratch/wisp/probes/260/r5/), so deleting the failure-receipt write at
// ball_windows.go:892-899 out of production leaves this whole package green.
// That silence is AC#2's defect.
//
// Reachability, answered before writing anything (full reading: instrument.md
// section 1 in that same directory):
//
//   - The only non-nil error inside the closure comes from the real Win32 call.
//     takeEscBorrow (hotkey_windows.go:631-633) hands takeEscWithAcc the concrete
//     hotkeyRegisterer, which is a function and not a variable: the registerFn
//     seam that registerAllWith takes was never parameterised on the borrow path.
//     Making it one is a production change this leg is not authorised to make.
//   - The candidate shape "an empty or invalid hwnd makes RegisterHotKey fail on
//     its own" is HALF true, measured here instead of assumed: hWnd is an optional
//     parameter, so hwnd=NULL SUCCEEDS (rc=1) and would really take the combination
//     away from the desktop; only an invalid NON-ZERO handle refuses (rc=0,
//     errno=1400 ERROR_INVALID_WINDOW_HANDLE), and it refuses before anything is
//     registered, so this case lends nothing.
//   - Getting INTO the closure needs no window: Ball.uiRun runs the closure inline
//     when the calling thread is the ui-sta thread (ball_windows.go:741-745), and
//     staThread.start records its thread id BEFORE the create step and returns
//     without ever pumping when create fails (sta_windows.go:59-90). That is the
//     same door ticket 33's sta_release_windows_test.go measures at non-winelive.
//
// So the rig is a ball with no window, no renderer and no tray, whose ui-sta thread
// IS the test goroutine, whose report starts as the healthy idle one (three live
// slots plus ticket 245's standby cancel line, built through the deterministic fake
// registry - which is never called for the cancel slot at all), and whose b.hwnd
// names a handle this process does not own. Then the shipped TakeEscForCancel runs
// end to end and the refusal it meets is the real one.
//
// Deliberate restraint: the planted cancel binding is a modifier combination and
// never the bare Esc. Handle validation happens before registration, so the attempt
// below cannot succeed - but a ruler whose worst-case failure mode is "took Esc from
// every window on the machine" is not a ruler, and that is ticket 245's whole point,
// so the default spelling is never walked through real Win32 here.

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"testing"

	"golang.org/x/sys/windows"
)

const (
	// t260r5PlantHWND is a USER handle value this process does not own. Not 0
	// (that half registers), not 0xFFFFFFFF (HWND_MESSAGE, a real special window),
	// and the cases below ask user32 itself whether it is a window.
	t260r5PlantHWND = uintptr(0x001D0A5C)

	// t260r5CancelBinding is the planted [hotkey] cancel value: a combination, so
	// even a misread Win32 could not take a bare key off the desktop.
	t260r5CancelBinding = "Ctrl+Alt+F9"

	// t260r5WerrInvalidWindowHandle is raw Win32 1400, spelled as a number rather
	// than through the x/sys name for the same reason accelDefault260 exists in
	// ticket 260-r1's file: a ruler that reads its expectation out of the same
	// constant the code under test uses cannot see that constant moving.
	t260r5WerrInvalidWindowHandle = 1400

	// t260r5RefusedLineWord and t260r5RefusedLogWord are the two sentences the
	// wiring owes the user, quoted from production instead of restated: Problems()'
	// attempted-but-refused branch (hotkey_windows.go:399-402) and takeEscWithAcc's
	// own complaint (hotkey_windows.go:612-613). The Win32 text itself is NOT
	// spelled out here - syscall.Errno.Error() runs FormatMessage, which is
	// localized per machine, so the cases compare against line.Err.Error().
	t260r5RefusedLineWord = "was not registered"
	t260r5RefusedLogWord  = "cancel takeover refused by Win32"
)

// t260r5IsWindow is this file's own handle to user32's IsWindow. It is deliberately
// not named pIsWindow: live_guard_windows_test.go owns that name and is compiled
// next to this file under -tags winlive.
var t260r5IsWindow = modUser32.NewProc("IsWindow")

// t260r5ErrNoWindowDoor is what the rig's create step returns so that
// staThread.start never reaches the message pump: no pump, no window, and the
// thread goes back through the create-error door.
var t260r5ErrNoWindowDoor = errors.New("260-r5 rig: this door creates no window and runs no pump")

// t260r5Errno is Win32 1400 as an error value, in one place, so the two
// comparisons below read as the same number this file claims to expect.
func t260r5Errno() syscall.Errno { return syscall.Errno(t260r5WerrInvalidWindowHandle) }

// t260r5LockedWriter is a slog writer a handler may write to from any goroutine;
// the borrow runs on the ui-sta thread while the test goroutine reads the tail.
type t260r5LockedWriter struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (w *t260r5LockedWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.Write(p)
}

func (w *t260r5LockedWriter) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.String()
}

// t260r5Reading is everything the rig measured on the ui-sta thread, carried out to
// the test goroutine to be judged.
type t260r5Reading struct {
	seed        HotkeyReport // the healthy idle report the ball was handed
	seedProbs   []string     // what Problems() said BEFORE the borrow
	plantIsWin  uintptr      // IsWindow(plant): must be 0 or this run measured nothing
	inline      bool         // did Ball.uiRun really take its inline branch
	entered     bool         // did the closure actually run
	report      HotkeyReport // what the ball's report reads after the borrow
	takenOver   bool         // b.escTakenOver after the borrow, as read by the rig
	lentKey     bool         // the attempt leg came back holding the key (the plant was accepted)
	liveAfter   map[uint32]Accelerator
	logged      string  // error-level output the production borrow said out loud
	staHWND     uintptr // staThread.hwnd at the door: this rig never sets it
	doorErr     error
	unregNullRc uintptr // release attempt for the unexpected-success door
}

// t260r5RunBorrow stands the windowless ui-sta door up around one ball and runs the
// shipped Ball.TakeEscForCancel on it. cfg.Cancel is the planted binding.
// preBorrowed leaves the ball in the state a SUCCESSFUL borrow already got to, which
// is the leg that keeps the refused-borrow assertions from being a standing
// announcement.
func t260r5RunBorrow(t *testing.T, cfg HotkeyConfig, preBorrowed bool) t260r5Reading {
	t.Helper()

	// The seed: the idle report production builds at boot (createOnSTA,
	// ball_windows.go:250), through the deterministic fake registry ticket 64/245
	// built. Nothing in it is a failure line - the only thing that can make the
	// cancel slot stop being standby is the code under test.
	reg := &fakeRegistry{}
	seed := registerAllWith(0, cfg, reg.register)
	if n := reg.attemptsOf(hkCancel); n != 0 {
		t.Fatalf("the rig's own seed registered the cancel slot (%d attempts): ticket 245's idle rule is already "+
			"broken inside the fixture, so nothing below can tell a fixture from a wiring failure", n)
	}
	if c, ok := seed.Binding(hkCancel); !ok || c.Status != HotkeyStandby || c.Err != nil {
		t.Fatalf("the seed cancel line = %+v (present=%t), want a standby line with no error - the failure this "+
			"case looks for must be written by ball_windows.go:895, not already be there before the borrow", c, ok)
	}
	seedProbs := seed.Problems()

	b := &Ball{}
	s := newSTAThread(nil)
	b.sta = s
	b.boundCfg = cfg
	b.cancelBinding = cfg.Cancel
	b.registeredHotkeys = seed.Live()
	b.escTakenOver = preBorrowed
	b.hotkeyReport = seed
	held := HotkeyReport{}
	if preBorrowed {
		// The shape a borrow that WORKED leaves behind: the cancel slot reads live
		// and carries the configured key. A pre-condition, not an assertion target -
		// the leg that uses it demands that a borrow which was never attempted
		// changes nothing about it.
		held = seed.withCancel(cancelBorrowedLineFor(resolveCancelBorrow(cfg.Cancel)))
		b.hotkeyReport = held
	}

	var o t260r5Reading
	o.seed = seed
	o.seedProbs = seedProbs
	w := &t260r5LockedWriter{}
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelError})))

	// start() runs on THIS goroutine: it locks the OS thread, records its id as the
	// ui-sta thread and runs create() on it, so b.uiRun's first branch (same thread
	// id) executes the closure inline. create() returns an error, so the pump never
	// runs and staThread.hwnd is never set - no window is created anywhere in here.
	s.start(func(st *staThread) error {
		o.plantIsWin, _, _ = t260r5IsWindow.Call(t260r5PlantHWND)
		o.inline = windows.GetCurrentThreadId() == st.threadID()
		b.hwnd = windows.HWND(t260r5PlantHWND)
		o.entered = true
		b.TakeEscForCancel()
		o.report = b.hotkeyReport
		o.takenOver = b.escTakenOver
		o.liveAfter = b.registeredHotkeys
		o.lentKey = b.escTakenOver && !preBorrowed
		if o.lentKey {
			// The plant was accepted after all, so a real combination is now held by
			// this thread: release it on the thread that took it, both doors, before
			// the case below complains about it. This is the rig's own hygiene, and
			// only the attempt leg can need it.
			b.ReleaseEscAfterSession()
			r, _, _ := pUnregisterHotKey.Call(0, uintptr(hkCancel))
			o.unregNullRc = r
		}
		return t260r5ErrNoWindowDoor
	})

	slog.SetDefault(previous)
	o.logged = w.String()
	o.doorErr = s.waitStarted()
	s.mu.Lock()
	o.staHWND = uintptr(s.hwnd)
	s.mu.Unlock()
	if preBorrowed {
		o.seed = held // the re-entry leg compares against what it actually started from
	}
	return o
}

// TestWindowlessSTAReachesTheBorrowWiring260r5 is the premise case: it pins that the
// rig really is what AC#2's judge needs - a non-winelive run that reaches
// ball_windows.go:892-899 through the shipped TakeEscForCancel, with no window, no
// renderer, no tray and no hot key left registered. If any of this stops being true
// the case goes red (it never skips), which is also the earliest warning that the
// two cases below stopped measuring anything.
func TestWindowlessSTAReachesTheBorrowWiring260r5(t *testing.T) {
	if windows.ERROR_INVALID_WINDOW_HANDLE != t260r5WerrInvalidWindowHandle {
		t.Fatalf("raw Win32 1400 is not what x/sys calls ERROR_INVALID_WINDOW_HANDLE any more (it is %d): this "+
			"file's expectation went stale and the refusal below would go unrecognised",
			windows.ERROR_INVALID_WINDOW_HANDLE)
	}
	cfg := t260r5Config()
	// The plant must be a combination and never the bare Esc, or this file would be
	// one Win32 change away from taking Esc from the whole desktop.
	acc, err := ParseAccelerator(t260r5CancelBinding)
	if err != nil || acc.VK == vkEscape || acc.Mods&modAlt == 0 || acc.Mods&modControl == 0 {
		t.Fatalf("the planted cancel binding %q is not a modifier combination (%+v, %v) - refuse to run",
			t260r5CancelBinding, acc, err)
	}

	o := t260r5RunBorrow(t, cfg, false)

	if o.plantIsWin != 0 {
		t.Fatalf("IsWindow(0x%X) says the planted handle is a live window, so it is not a refused handle and this "+
			"run measured a borrow that could have succeeded", t260r5PlantHWND)
	}
	if !o.entered {
		t.Fatal("the uiRun closure never ran, so ball_windows.go:892-899 was never reached")
	}
	if !o.inline {
		t.Fatal("Ball.uiRun did not take its inline branch: the borrow went through PostTask, which with no window " +
			"drops the task (sta_windows.go:240-247), so this rig cannot claim anything below was executed")
	}
	if o.staHWND != 0 {
		t.Errorf("the rig's staThread recorded a window handle 0x%X - this case promised to create no window at all", o.staHWND)
	}
	if !errors.Is(o.doorErr, t260r5ErrNoWindowDoor) {
		t.Errorf("the create-error door was not taken (err=%v): if the pump ran instead, this rig is holding a real "+
			"message loop on a test goroutine and nothing below is readable", o.doorErr)
	}
	if o.lentKey {
		t.Fatalf("the ball ends the borrow claiming the cancel key is its own (escTakenOver=true) on a handle IsWindow "+
			"calls invalid: either Win32 really lent it (release rc=%d on the NULL door) or the wiring reports a key it "+
			"never got, which is the fall-through door at ball_windows.go:897. Neither leaves this case judging a refusal",
			o.unregNullRc)
	}
	line, ok := o.report.Binding(hkCancel)
	if !ok {
		t.Fatal("the report has no cancel line at all after the borrow")
	}
	if line.Err == nil {
		t.Fatalf("the cancel line carries no error (%+v): the rig never met a refusal", line)
	}
	var errno syscall.Errno
	if !errors.As(line.Err, &errno) || errno != t260r5Errno() {
		t.Errorf("the cancel line carries %v (%T), want Win32 %d ERROR_INVALID_WINDOW_HANDLE - the error the "+
			"receipt prints has to be the one user32 set", line.Err, line.Err, t260r5WerrInvalidWindowHandle)
	}
	if _, live := o.liveAfter[hkCancel]; live {
		t.Errorf("the registration set still claims the cancel key after a refusal: %v", o.liveAfter)
	}
	t.Logf("260-r5 rig: IsWindow(0x%X)=%d uiRun-inline=%t closure-entered=%t sta.hwnd=0x%X door=%q refusal=%q",
		t260r5PlantHWND, o.plantIsWin, o.inline, o.entered, o.staHWND, o.doorErr, line.Err)
}

// TestRefusedBorrowIsNamedOnTheReport260r5 is AC#2's ruler: the borrow request went
// out, Win32 refused it, and the step that says so is the one production runs.
// Deleting that write (ball_windows.go:895) makes this case red - the receipt keeps
// the standby line it started with, Problems() stays empty, and the card sits there
// with no cancel key and nothing to show for it. That is exactly the zero-symptom
// shape the ticket names, and this case is the instrument that hears it.
func TestRefusedBorrowIsNamedOnTheReport260r5(t *testing.T) {
	cfg := t260r5Config()
	o := t260r5RunBorrow(t, cfg, false)

	if o.plantIsWin != 0 || !o.entered || o.lentKey {
		t.Fatalf("the rig did not reach a refused borrow (IsWindow=%d entered=%t the-ball-claims-the-key=%t): either the "+
			"plant stopped being a refusal or the wiring claims a key Win32 refused - in both readings the receipt below "+
			"would be judging a fixture and not ball_windows.go:895", o.plantIsWin, o.entered, o.lentKey)
	}
	line, ok := o.report.Binding(hkCancel)
	if !ok {
		t.Fatal("the refused borrow left no cancel line in the report: the receipt write never ran")
	}
	// The user-visible symptom is checked FIRST, because that is what AC#2 is about:
	// zero problem lines before the borrow, exactly the one that names it afterwards.
	// Deleting the write at ball_windows.go:895 lands here.
	if len(o.seedProbs) != 0 {
		t.Fatalf("the seed report already complained (%q) before the borrow was attempted, so counting problem "+
			"lines below measures the fixture and not the wiring", o.seedProbs)
	}
	probs := o.report.Problems()
	if len(probs) != 1 {
		t.Fatalf("problem lines after a borrow Win32 refused = %d (%q), want the 1 that says so: the seed said %q, "+
			"and the write at ball_windows.go:895 is the only thing between the two", len(probs), probs, o.seedProbs)
	}
	said := probs[0]
	errText := "<nil>"
	if line.Err != nil {
		errText = line.Err.Error()
	}
	for _, want := range []string{"cancel", t260r5CancelBinding, t260r5RefusedLineWord, errText} {
		if !strings.Contains(said, want) {
			t.Errorf("the refused-borrow line does not say %q; it reads: %s", want, said)
		}
	}
	if line.Status != HotkeyError || !line.Status.Attempted() {
		t.Errorf("cancel line after a refused borrow = %q (attempted=%t), want HotkeyError and attempted: the "+
			"borrow WAS attempted and Win32 said no, which is the pair Problems() renders", line.Status, line.Status.Attempted())
	}
	if line.Err == nil {
		t.Fatal("the cancel line carries no error after a real Win32 refusal: the failure was swallowed")
	}
	var errno syscall.Errno
	if !errors.As(line.Err, &errno) || errno != t260r5Errno() {
		t.Errorf("the cancel line names %v, want the errno user32 set (%d): the receipt must print the real refusal "+
			"and not a value this file could have written itself", line.Err, t260r5WerrInvalidWindowHandle)
	}
	// Condition 2 of the ruling, at the wiring layer: the key named in the failure
	// line is the key that was really handed to RegisterHotKey.
	wantAcc, err := ParseAccelerator(t260r5CancelBinding)
	if err != nil {
		t.Fatalf("ParseAccelerator(%q) fails, so the plant itself is broken: %v", t260r5CancelBinding, err)
	}
	if line.Binding != t260r5CancelBinding || line.Acc != wantAcc {
		t.Errorf("the refused borrow reports %q / %+v, want %q / %+v - the line names a key that was not the one "+
			"attempted", line.Binding, line.Acc, t260r5CancelBinding, wantAcc)
	}

	// A card with no cancel key must not be reported as holding one.
	if o.takenOver {
		t.Errorf("escTakenOver reads true after a refused borrow: the ball claims a key Win32 never lent, and " +
			"ReleaseEscAfterSession would later drop whatever it thinks it took")
	}
	// withCancel replaces the cancel line and nothing else: the three slots an idle
	// ball does hold must still read the way the seed left them.
	for _, other := range []struct {
		id   uint32
		name string
	}{{hkSummon, "summon"}, {hkMute, "mute"}, {hkPanel, "panel"}} {
		before, _ := o.seed.Binding(other.id)
		after, ok := o.report.Binding(other.id)
		if !ok || after.Status != before.Status || after.Binding != before.Binding {
			t.Errorf("the refused borrow moved the %s line from %+v to %+v (present=%t): a receipt write that "+
				"clobbers the other slots reports a registration set Win32 does not hold", other.name, before, after, ok)
		}
	}
	// The other symptom channel, said by production (takeEscWithAcc) and not by this
	// file. It is the second half of "零症状": a refusal that logs nothing is only
	// half fixed by a receipt line.
	if !strings.Contains(o.logged, t260r5RefusedLogWord) {
		t.Errorf("the refused borrow logged nothing naming %q; error-level output captured was %q",
			t260r5RefusedLogWord, o.logged)
	}
	t.Logf("260-r5 refused borrow: line=%+v problems=%q logged=%q", line, probs, o.logged)
}

// TestRefusedBorrowIsNotAStandingAnnouncement260r5 is the sensitivity control on the
// same ruler: the success-shaped state. The ball already holds the borrow, so
// TakeEscForCancel's idempotency branch (ball_windows.go:883-885) returns before any
// Win32 call and nothing may change - no new line, no new complaint, no rewrite of
// the receipt. A ruler that also went red here would be a ruler that always answers
// "the borrow failed", which is the other way AC#2 could be faked.
func TestRefusedBorrowIsNotAStandingAnnouncement260r5(t *testing.T) {
	cfg := t260r5Config()
	o := t260r5RunBorrow(t, cfg, true)

	if !o.entered {
		t.Fatal("the uiRun closure never ran, so this control leg judged nothing")
	}
	if o.lentKey {
		t.Fatal("the rig's own bookkeeping is wrong: the re-entry leg cannot lend a key it already held")
	}
	held, _ := o.seed.Binding(hkCancel)
	after, ok := o.report.Binding(hkCancel)
	if !ok {
		t.Fatal("the live cancel line disappeared from the report")
	}
	if after.Status != HotkeyLive || after.Err != nil {
		t.Errorf("a ball that already holds the borrow ended up reading %q (err=%v), want it unchanged from %+v: the "+
			"receipt is being rewritten on a path that never attempted a registration", after.Status, after.Err, held.Status)
	}
	if p := o.report.Problems(); len(p) != 0 {
		t.Errorf("a held cancel key reports problems (%q); Problems() then cannot tell a refused borrow from a "+
			"standing announcement, and the ruler above is insensitive", p)
	}
	if o.logged != "" {
		t.Errorf("the idempotent re-entry logged something at error level (%q); it must not reach Win32 at all", o.logged)
	}
	if !reflect.DeepEqual(o.report.Bindings(), o.seed.Bindings()) {
		t.Errorf("the report changed on the re-entry leg (%s vs %s): a re-entry that does nothing must leave the "+
			"whole receipt alone, not just the line this case reads",
			describe260r5(o.report.Bindings()), describe260r5(o.seed.Bindings()))
	}
	if !o.takenOver {
		t.Error("escTakenOver went false on a re-entry that does nothing: the ball now reports no cancel key while a " +
			"card is waiting, which is the same zero-symptom shape seen from the other side")
	}
}

// t260r5Config is the planted [hotkey] section: cancel moved off the shipped default
// (ticket 260 形 a is what makes the borrow read the config at all), the other three
// slots left where ApplyHotkeyDefaults puts them so the seed report looks like a
// normal boot.
func t260r5Config() HotkeyConfig {
	return ApplyHotkeyDefaults(HotkeyConfig{Cancel: t260r5CancelBinding})
}

// describe260r5 renders a report for an error message without leaning on %v of an
// unexported slice.
func describe260r5(bindings []HotkeyBinding) string {
	parts := make([]string, 0, len(bindings))
	for _, b := range bindings {
		parts = append(parts, fmt.Sprintf("%s=%q", b.Name, b.Status))
	}
	return strings.Join(parts, " ")
}
