//go:build windows

package ball

// Ticket 33 close-out, leg 33-r8: the thread-owner invariant on the ball's
// ui-sta (D38a). staThread.start pins its goroutine to an OS thread
// (runtime.LockOSThread) and hands that SAME thread back to the Go pool when it
// returns (defer runtime.UnlockOSThread, sta_windows.go:52-53). Whatever is on
// the thread at that instant - a window it created, a message still queued, an
// initialised COM apartment - is inherited by whoever the scheduler puts on it
// next.
//
// Both cases below run start() on a thread this test owns and locks, so the
// verdict is read ON the thread that is about to be released. That is not a
// stylistic choice: PeekMessage can only ever examine the calling thread's
// queue, so nobody but the owner can measure it at that moment.
// EnumThreadWindows IS legal cross-thread, so the window tally is reported from
// whichever goroutine asks.
//
// Nothing here plants another owner's poison inside the product or drains
// somebody else's queue: the two cases are the two exits staThread.start
// actually has, and every assertion is about what the thread looks like when its
// owner lets go of it.

import (
	"errors"
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Test-local Win32 entry points. The t33r8 prefix is load-bearing: under
// -tags winlive this file compiles next to live_guard_windows_test.go, which
// already owns pIsWindow / pEnumThreadWindows / pGetClassNameW and friends.
var (
	t33r8PeekMessageW      = modUser32.NewProc("PeekMessageW")
	t33r8EnumThreadWindows = modUser32.NewProc("EnumThreadWindows")
	t33r8GetClassNameW     = modUser32.NewProc("GetClassNameW")
	t33r8IsWindow          = modUser32.NewProc("IsWindow")
	t33r8PostThreadMessage = modUser32.NewProc("PostThreadMessageW")
)

const (
	// t33r8PmRemove is PM_REMOVE: take the message off the queue and dispatch it.
	t33r8PmRemove = 1
	// t33r8PmNoRemove is PM_NOREMOVE: ask what the head is without changing the
	// answer, so reading a queue is never a side effect.
	t33r8PmNoRemove = 0
	// t33r8WmQuit is WM_QUIT, named so a reading can say it out loud.
	t33r8WmQuit = 0x0012
	// t33r8BallClass is the ball's own window class (registerBallClass).
	t33r8BallClass = "WispBallWindow"
	// t33r8WmLeftover is a WM_APP slot nothing in this package owns. The plant
	// posts it because it must be a message the queue REPORTS: the first draft of
	// this file planted WM_NULL, and a release reading of "empty" came back even
	// though PostMessageW had answered rc=1 - WM_NULL is not something these
	// instruments can see, so it would have pinned nothing.
	t33r8WmLeftover = 0x82F8
	// t33r8PumpCap bounds a release pump. Hitting the cap is a loud reading,
	// never a silent green.
	t33r8PumpCap = 4096
)

// The EnumThreadWindows callback carries no usable context through LPARAM, so
// the tallies sit at package scope and are reset before each enumeration (same
// shape the panel family's tally uses).
var (
	t33r8BallWins atomic.Int64
	t33r8AllWins  atomic.Int64
	// t33r8Classes is the class list of the enumeration in flight. Plain slice,
	// no lock: EnumThreadWindows calls back synchronously on the calling thread,
	// and every reader in this file runs on the one plant goroutine.
	t33r8Classes []string
	t33r8Tally   = windows.NewCallback(func(hwnd, _ uintptr) uintptr {
		buf := make([]uint16, 64)
		t33r8GetClassNameW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
		name := windows.UTF16ToString(buf)
		t33r8AllWins.Add(1)
		t33r8Classes = append(t33r8Classes, name)
		if name == t33r8BallClass {
			t33r8BallWins.Add(1)
		}
		return 1
	})
)

// threadWindowTally asks user32, not Go, how many top-level windows a thread
// owns, which classes they are, and how many of them are the ball's own class.
// The class list is what makes the ball-count honest: a Go thread is not raw
// count 0 (see live_guard_windows_test.go's note about the runtime's own helper
// window), so the invariant can only ever be about the windows THIS owner made.
func threadWindowTally(tid uint32) (ball, all int, classes string) {
	t33r8BallWins.Store(0)
	t33r8AllWins.Store(0)
	t33r8Classes = nil
	t33r8EnumThreadWindows.Call(uintptr(tid), uintptr(t33r8Tally), 0)
	return int(t33r8BallWins.Load()), int(t33r8AllWins.Load()), fmt.Sprint(t33r8Classes)
}

// queueHead names the first message still queued on the CALLING thread, or
// "empty". An unfiltered PM_NOREMOVE peek: 33-r7 measured that a FILTERED peek
// cannot see a latched quit (cmd/wisp/panel_host_windows.go's
// drainStaleQuitBeforeCreate returned removed=0 while the quit that killed the
// next create was still there), so a filtered head would be a lying instrument.
func queueHead() string {
	var m msg
	r, _, _ := t33r8PeekMessageW.Call(unsafePtr(&m), 0, 0, 0, uintptr(t33r8PmNoRemove))
	if r == 0 {
		return "empty"
	}
	if m.message == t33r8WmQuit {
		return "WM_QUIT (latched quit)"
	}
	return fmt.Sprintf("msg 0x%X on hwnd 0x%X", m.message, uintptr(m.hwnd))
}

// releaseReading is what a thread looks like at the moment its owner could hand
// it back.
type releaseReading struct {
	tid      uint32
	ballWins int
	allWins  int
	classes  string
	head     string
}

// readThisThread takes the releaseReading for the CALLING thread.
func readThisThread() releaseReading {
	tid := uint32(windows.GetCurrentThreadId())
	ball, all, classes := threadWindowTally(tid)
	return releaseReading{tid: tid, ballWins: ball, allWins: all, classes: classes, head: queueHead()}
}

// t33r8OwnWindow makes one real (hidden) window of the ball's own class on the
// calling thread and records it the way the production owner does, through
// staThread.hwnd (ball_windows.go:218-220). It does not pump, so whatever
// creation queued is left queued.
func t33r8OwnWindow(s *staThread) (windows.HWND, error) {
	if err := registerBallClass(); err != nil {
		return 0, err
	}
	cls := utf16(t33r8BallClass)
	title := utf16("wisp-33r8-thread-owner-nail")
	h, _, err := pCreateWindowExW.Call(
		uintptr(wsExToolWindow|wsExNoActivate), unsafePtr(cls), unsafePtr(title),
		uintptr(wsPopup), 0, 0, 40, 40, 0, 0, uintptr(moduleHandle()), 0)
	if h == 0 {
		return 0, fmt.Errorf("plant: CreateWindowExW: %w", err)
	}
	s.mu.Lock()
	s.hwnd = windows.HWND(h)
	s.mu.Unlock()
	return windows.HWND(h), nil
}

// t33r8Teardown is this file's own hygiene, the same duty 33-r7 wrote into the
// panel family: a case that reads a dirty thread must not then BE the thing that
// poisons the pool, or a red run here breaks whichever later case the scheduler
// puts on the thread. Every assertion is taken BEFORE it runs; what it touches
// is only ever a window this file made on this thread.
func t33r8Teardown(hwnd windows.HWND) {
	if r, _, _ := t33r8IsWindow.Call(uintptr(hwnd)); r != 0 {
		pDestroyWindow.Call(uintptr(hwnd))
	}
	pumpThisThreadToQuiet(t33r8PumpCap)
}

// pumpThisThreadToQuiet dispatches (never drops) the calling thread's queued
// messages until user32 reports the queue empty, bounded by limit. It returns
// how many it moved.
func pumpThisThreadToQuiet(limit int) int {
	pumped := 0
	for pumped < limit {
		var m msg
		r, _, _ := t33r8PeekMessageW.Call(unsafePtr(&m), 0, 0, 0, uintptr(t33r8PmRemove))
		if r == 0 {
			return pumped
		}
		pTranslateMessage.Call(unsafePtr(&m))
		pDispatchMessageW.Call(unsafePtr(&m))
		pumped++
	}
	return pumped
}

// TestSTAReleaseAfterFailedCreateHandsBackNoWindow is the create-error exit door
// of staThread.start, and it is a shape production can reach today: createOnSTA
// records staThread.hwnd at ball_windows.go:219 and can still fail AFTER that,
// at newRenderer (:239) or addTrayIcon (:247). start then closes s.started and
// returns at sta_windows.go:70; its deferred UnlockOSThread hands the thread to
// the pool with the window still on it. No pump ever runs, so nothing unwinds -
// and nothing drains either, which is why this is the door where a message the
// create step left queued also goes back to the pool (planted below, with the
// rc and the head reported as read).
func TestSTAReleaseAfterFailedCreateHandsBackNoWindow(t *testing.T) {
	type outcome struct {
		panicked       any
		startErr       error
		plantErr       error
		hwnd           windows.HWND
		queuedRc       uintptr
		queuedThreadRc uintptr
		afterPlant     releaseReading
		afterQueued    releaseReading
		atRelease      releaseReading
		afterOwn       releaseReading
	}
	res := make(chan outcome, 1)
	go func() {
		// owner: this test's plant goroutine. It locks its own thread, runs the
		// production start(), reads the thread before it can be handed back, and
		// recovers below so a panic reaches the report instead of the binary.
		var o outcome
		defer func() {
			o.panicked = recover()
			res <- o
		}()
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		s := newSTAThread(nil)
		s.start(func(st *staThread) error {
			h, err := t33r8OwnWindow(st)
			if err != nil {
				o.plantErr = err
			} else {
				o.hwnd = h
			}
			o.afterPlant = readThisThread()
			// Two plants, because they are cleaned by two different things and the
			// assertion has to be able to tell them apart. Both stand for what a
			// real create step can leave queued on a door that never pumps: a tray
			// callback (Shell_NotifyIcon answers asynchronously, and addTrayIcon
			// is one of the two failures this door stands for) or a hotkey already
			// pressed while the window was being built, both window-addressed, and
			// anything at all that is addressed to the THREAD rather than to a
			// window. Measured (mutation M2,
			// .scratch/wisp/probes/33/r8/logs/mutation-M2.txt): with only the
			// window-addressed plant this case stayed GREEN after the release pump
			// was dropped, because DestroyWindow takes the messages addressed to
			// the window it destroys out of the queue with it. The thread-addressed
			// one is what DestroyWindow cannot purge, so that is the plant that
			// makes the pump's half of the duty bite. It is also exactly what a
			// later owner of this thread is handed as "mine".
			o.queuedRc, _, _ = pPostMessageW.Call(uintptr(st.hwnd), t33r8WmLeftover, 0, 0)
			o.queuedThreadRc, _, _ = t33r8PostThreadMessage.Call(
				uintptr(windows.GetCurrentThreadId()), t33r8WmLeftover+1, 0, 0)
			o.afterQueued = readThisThread()
			// This is the door: newRenderer / addTrayIcon failing after the
			// window exists. The error value is the product's own return path.
			return errors.New("plant: renderer init failed")
		})
		o.startErr = s.waitStarted()
		o.atRelease = readThisThread()
		t33r8Teardown(o.hwnd)
		o.afterOwn = readThisThread()
	}()
	o := <-res

	t.Logf("33-r8 create-error door: tid=%d planted hwnd=0x%X | after plant: ball windows=%d all windows=%d classes=%s queue head=%s",
		o.atRelease.tid, uintptr(o.hwnd), o.afterPlant.ballWins, o.afterPlant.allWins, o.afterPlant.classes, o.afterPlant.head)
	t.Logf("33-r8 create-error door: after the plants (window-addressed PostMessageW rc=%d, thread-addressed PostThreadMessageW rc=%d): ball windows=%d all=%d classes=%s queue head=%s",
		o.queuedRc, o.queuedThreadRc, o.afterQueued.ballWins, o.afterQueued.allWins, o.afterQueued.classes, o.afterQueued.head)
	t.Logf("33-r8 create-error door: the instant start() returned: ball windows=%d all windows=%d classes=%s queue head=%s (start err=%v plant err=%v)",
		o.atRelease.ballWins, o.atRelease.allWins, o.atRelease.classes, o.atRelease.head, o.startErr, o.plantErr)

	if o.panicked != nil {
		t.Fatalf("the plant goroutine panicked: %v", o.panicked)
	}
	if o.plantErr != nil {
		t.Fatalf("the plant made no window (%v), so this run measured nothing about the thread's release state", o.plantErr)
	}
	if o.startErr == nil {
		t.Fatalf("start() returned nil although the create step failed: the create-error door was never taken, so this run measured nothing")
	}
	if o.afterPlant.ballWins != 1 {
		t.Fatalf("the plant took no effect on this thread's window tally (ball windows=%d, want 1), so the assertions below would be measuring nothing", o.afterPlant.ballWins)
	}
	if o.queuedRc == 0 || o.queuedThreadRc == 0 {
		t.Fatalf("the plants did not both take (window-addressed rc=%d, thread-addressed rc=%d), so this run cannot tell a drained queue from an unmeasured one", o.queuedRc, o.queuedThreadRc)
	}
	if o.afterQueued.head == "empty" {
		t.Fatalf("the plant is not visible to this instrument (head=empty right after a successful PostMessageW on a door that never pumps), so the queue assertion below would be measuring nothing. WM_NULL is exactly such an invisible plant - measured, see the case comment")
	}

	// The invariant: a thread that still owns a window, or still carries a
	// queued message, is not the Go pool's to reuse.
	if o.atRelease.ballWins != 0 {
		t.Errorf("staThread.start is about to hand OS thread %d back to the Go pool while it still owns %d live %s window(s) (planted hwnd 0x%X). The create step failed, so nothing unwound it. A thread that owns a window its Go manager has already let go of is what makes the NEXT create on it fail (33-r7 shape ①), and no pre-create check on the receiving side can see it: the panel's staleCloseQueued looks at the QUEUE, not at the handle table. Fix: the owner destroys the windows it started before it releases the thread",
			o.atRelease.tid, o.atRelease.ballWins, t33r8BallClass, uintptr(o.hwnd))
	}
	if o.atRelease.head != "empty" {
		t.Errorf("staThread.start is about to hand OS thread %d back with its message queue non-empty (head=%s). A thread released while it still carries messages is the order-dependent disease 33-r7 named in the panel family: whoever lands here next pumps messages that were not addressed to them. Fix: the owner dispatches its own queue to empty before it releases the thread",
			o.atRelease.tid, o.atRelease.head)
	}
	// This case's own teardown: a red reading must not become pool poison for a
	// later case.
	if o.afterOwn.ballWins != 0 || o.afterOwn.head != "empty" {
		t.Errorf("this case released OS thread %d with ball windows=%d and queue head=%s AFTER its own teardown: a red run here would now poison whichever later test the scheduler puts on this thread",
			o.afterOwn.tid, o.afterOwn.ballWins, o.afterOwn.head)
	}
}

// TestSTAReleaseAfterPumpExitDispatchesItsQueue is the other exit door: the pump
// ran and ended. The door is taken twice on ONE locked thread because the two
// answers are different questions:
//
//	kept-window: the pump ended for a reason that was NOT the ball's own Close.
//	    Production reaches this today - the GetMessageW error break
//	    (sta_windows.go:83-90), and on any thread that inherits a latched quit,
//	    the pump's very first iteration. The window stays live, and a message
//	    posted after the quit stays queued.
//	close-shaped: DestroyWindow, then staThread.quit - the order Ball.Close uses
//	    (ball_windows.go:946-962). This one is the honest control: it is what the
//	    shipping path leaves behind, and whatever it reads is reported as read,
//	    green or not.
//
// The queued leftover in the kept-window door is a plant, and the plant has a
// production name: Close never stops the task port (staThread.hwnd is not
// cleared by Ball.Close, which zeroes only Ball.hwnd), so any goroutine calling
// SetState / SetBadge / PostTask between quit() and the thread's exit lands its
// message BEHIND the quit, where no pump will ever read it. Same queue shape,
// different msg id.
func TestSTAReleaseAfterPumpExitDispatchesItsQueue(t *testing.T) {
	type run struct {
		door        string
		hwnd        windows.HWND
		createErr   error
		afterCreate releaseReading
		leftoverRc  uintptr
		staHwnd     uintptr
		pumpedAtRel int
		atRelease   releaseReading
		afterOwn    releaseReading
	}
	type outcome struct {
		panicked any
		runs     []run
	}
	res := make(chan outcome, 1)
	go func() {
		// owner: this test's plant goroutine; locks its thread, runs the real
		// pump twice, reads the thread at each release, recovers below.
		var o outcome
		defer func() {
			o.panicked = recover()
			res <- o
		}()
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		for _, realClose := range []bool{false, true} {
			r := run{door: "quit-only"}
			if realClose {
				r.door = "real-Close"
			}
			s := newSTAThread(nil)

			// The task port runs through the ONE package-level window procedure,
			// which routes to the active Ball by handle (ball_windows.go:540-547).
			// Without an active Ball a posted task is DefWindowProc'd and never
			// runs (measured the hard way: this plant hung on its first draft),
			// so the plant stands up the minimum Ball the routing needs. It is
			// not a real ball: no renderer, no tray, no timers, no hotkeys.
			b := &Ball{}
			b.sta = s
			restored := activeBall.Swap(b)

			go func() {
				// owner: the exit driver for one pump run. It waits for the
				// thread to be up through the product's own door and then takes
				// the exit door from OUTSIDE the pump, which is where Ball.Close
				// is called from in production (a shutdown hook, not the pump).
				if err := s.waitStarted(); err != nil {
					return
				}
				if realClose {
					b.Close() // the shipping door, whole: ball_windows.go:935-966
					return
				}
				done := make(chan struct{})
				s.PostTask(func() {
					defer close(done)
					s.quit() // WM_QUIT only (sta_windows.go:156-162): the window stays up
					// The plant: a message that loses the race against the
					// exit and lands behind the quit. In production this is a
					// SetState / SetBadge / PostTask that arrives after
					// Ball.Close has already posted WM_QUIT.
					r.leftoverRc, _, _ = pPostMessageW.Call(uintptr(s.hwnd), t33r8WmLeftover, 0, 0)
				})
				<-done
			}()

			s.start(func(st *staThread) error {
				h, err := t33r8OwnWindow(st)
				if err != nil {
					return err
				}
				b.hwnd = h // ballWndProc's routing key, same thread
				r.hwnd = h
				r.afterCreate = readThisThread()
				return nil // production door: the pump now runs
			})
			r.createErr = s.waitStarted()
			r.atRelease = readThisThread()
			s.mu.Lock()
			r.staHwnd = uintptr(s.hwnd) // the owner's own record, at the release door
			s.mu.Unlock()
			r.pumpedAtRel = pumpThisThreadToQuiet(t33r8PumpCap)
			t33r8Teardown(r.hwnd)
			r.afterOwn = readThisThread()
			activeBall.CompareAndSwap(b, restored)
			o.runs = append(o.runs, r)
		}
	}()
	o := <-res

	if o.panicked != nil {
		t.Fatalf("the plant goroutine panicked: %v", o.panicked)
	}
	if len(o.runs) != 2 {
		t.Fatalf("got %d runs, want 2 (kept-window door and close-shaped door)", len(o.runs))
	}
	for _, r := range o.runs {
		t.Logf("33-r8 pump-exit door=%s: tid=%d hwnd=0x%X createErr=%v | after create: ball windows=%d all=%d classes=%s head=%s | AT RELEASE: ball windows=%d all=%d classes=%s head=%s staThread.hwnd=0x%X | plant PostMessage rc=%d | this file's own pump moved %d | after own teardown: ball windows=%d head=%s",
			r.door, r.atRelease.tid, uintptr(r.hwnd), r.createErr,
			r.afterCreate.ballWins, r.afterCreate.allWins, r.afterCreate.classes, r.afterCreate.head,
			r.atRelease.ballWins, r.atRelease.allWins, r.atRelease.classes, r.atRelease.head, r.staHwnd,
			r.leftoverRc, r.pumpedAtRel, r.afterOwn.ballWins, r.afterOwn.head)
	}

	for _, r := range o.runs {
		if r.createErr != nil {
			t.Fatalf("door=%s: the create step failed (%v), so this run never took the pump-exit door", r.door, r.createErr)
		}
		if r.afterCreate.ballWins != 1 {
			t.Fatalf("door=%s: the window the create step made is not on this thread's tally (ball windows=%d), so this run measured nothing", r.door, r.afterCreate.ballWins)
		}
		if r.atRelease.ballWins != 0 {
			t.Errorf("door=%s: the pump ended and staThread.start is about to hand OS thread %d back to the Go pool still owning %d live %s window(s) (hwnd 0x%X). Nobody unwound it, so the next owner of that thread creates beside a window no manager owns, and this process keeps a USER handle it has lost the record of. Fix: the owner destroys the windows it started before it releases the thread",
				r.door, r.atRelease.tid, r.atRelease.ballWins, t33r8BallClass, uintptr(r.hwnd))
		}
		if r.atRelease.head != "empty" {
			t.Errorf("door=%s: the pump ended on a thread whose queue is NOT empty (head=%s). GetMessageW left it there when it broke on the quit, and UnlockOSThread now gives the whole queue to whoever the scheduler puts on thread %d, which cannot tell those messages from its own. Fix: the owner dispatches its own queue to empty before it releases the thread",
				r.door, r.atRelease.head, r.atRelease.tid)
		}
		if r.staHwnd != 0 {
			t.Errorf("door=%s: the thread is about to go back to the Go pool with staThread.hwnd still holding 0x%X. That field is the owner's list of what it started: left set, the release step reaches for a handle this process no longer owns (Win32 may hand it to another thread's window in the meantime), and PostTask keeps posting messages to a window that has no pump any more. Fix: whoever destroys the window forgets it in the same step",
				r.door, r.staHwnd)
		}
		if r.afterOwn.ballWins != 0 || r.afterOwn.head != "empty" {
			t.Errorf("this case released OS thread %d with ball windows=%d and queue head=%s AFTER its own teardown (door=%s): a red run here would now poison whichever later test the scheduler puts on this thread",
				r.afterOwn.tid, r.afterOwn.ballWins, r.afterOwn.head, r.door)
		}
	}
}
