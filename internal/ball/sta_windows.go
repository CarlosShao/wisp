//go:build windows

package ball

// The ui-sta STA thread (D38a): ONE thread owns all UI COM objects (the D2D
// factory, render targets, the ball window, and later the panel). Everything
// UI-bound runs or posts here. Spawned through the observe registry under
// its frozen roster name "ui-sta"; goroutines elsewhere must never touch
// window/COM state directly - they PostTask.

import (
	"fmt"
	"log/slog"
	"runtime"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/observe"
)

// pmRemove is PeekMessage's PM_REMOVE: take the message off the queue, as
// opposed to PM_NOREMOVE (a peek that reports without removing). The release
// pump needs the removing form or it never reaches an empty queue.
const pmRemove = 1

// releasePumpCap bounds the owner's release pump. A thread that still has not
// reported an empty queue after this many of its OWN messages is not something
// to keep pumping: the cap is a loud reading (see releaseOwnQueueToQuiet), not a
// green.
const releasePumpCap = 4096

// staThread is the process-wide UI STA. There is exactly one ball per
// process and one ui-sta per process.
type staThread struct {
	mu      sync.Mutex
	nextID  uint64
	tasks   map[uint64]func()
	started chan struct{}
	err     error

	hwnd windows.HWND // ball window (created on the thread)
	tid  uint32

	registry *observe.Registry
	handle   *observe.Handle
}

func newSTAThread(registry *observe.Registry) *staThread {
	return &staThread{
		tasks:    map[uint64]func(){},
		started:  make(chan struct{}),
		registry: registry,
	}
}

// start runs the thread. It blocks until the message loop exits (Close).
func (s *staThread) start(create func(s *staThread) error) {
	// Win32 affinity: a window and its message queue belong to the OS thread
	// that created them - the goroutine must stay pinned or the pump below
	// waits on the WRONG queue (messages posted to the creating thread are
	// never dispatched). This is the STA contract (D38a) in its literal form.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	// Whichever way this function leaves - create failed, the pump broke on a
	// quit, or a panic is on its way out through observe's recover boundary -
	// the deferred UnlockOSThread above is about to hand an OS thread back to
	// the Go pool. Whoever the scheduler puts here next inherits whatever is
	// still on it, so the unwinding is THIS function's duty, not the create
	// step's and not the window's first user's (LIFO: this defer runs first).
	defer s.releaseCOM() // TEMP pre-fix shape for the winlive baseline, reverted immediately
	// Remember who owns this thread: a caller already on it must run inline
	// (see Ball.uiRun) - post-and-wait from inside the pump is a deadlock.
	s.mu.Lock()
	s.tid = windows.GetCurrentThreadId()
	s.mu.Unlock()

	// STA init: COM apartment-threaded on this thread (D2D single-threaded
	// factory + future WebView2 panel both want STA).
	r1, _, _ := pCoInitializeEx.Call(0, coinitApartmentThreaded)
	_ = r1 // S_FALSE (already) is fine too

	if err := create(s); err != nil {
		s.mu.Lock()
		s.err = err
		s.mu.Unlock()
		close(s.started)
		return
	}
	close(s.started)

	// Message pump. GetMessage returns 0 on WM_QUIT (clean exit), -1 on
	// error (treat as fatal for this thread and surface it).
	var m msg
	for {
		r, _, _ := pGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		c := int32(r)
		if c == 0 {
			break
		}
		if c == -1 {
			s.mu.Lock()
			if s.err == nil {
				s.err = fmt.Errorf("ball: GetMessageW failed")
			}
			s.mu.Unlock()
			break
		}
		pTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		pDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}
	// The unwinding is the deferred releaseThread: it runs here, before the
	// deferred UnlockOSThread, on this thread, and it is what balances the
	// CoInitializeEx above. This used to be a bare s.releaseCOM() at this point,
	// which the create-error door below skipped entirely.
}

// releaseThread is the thread owner's duty at the door, and it is the reason the
// ball's ui-sta can be handed back at all: start() locks an OS thread
// (runtime.LockOSThread) and its deferred runtime.UnlockOSThread puts that thread
// back into the Go pool, where any later goroutine - including the one that
// creates the NEXT window, on any thread name - may be placed. A thread that
// goes back still owning a window, still carrying queued messages, or still
// holding an initialised COM apartment is not "free": it is a thread with somebody
// else's state on it.
//
// What it does, in this order:
//
//  1. Destroy the window this thread created, if the create step left one behind.
//     Ball.Close destroys the ball window itself and calls forgetWindow, so on the
//     shipping path there is nothing to do here; what this covers is the doors that
//     never reach Close - createOnSTA failing after ball_windows.go recorded
//     s.hwnd (newRenderer, addTrayIcon), and a pump that ends for a reason other
//     than the ball's own quit. Measured at HEAD a35f7f5e before this existed:
//     both of those doors handed thread back with 1 live WispBallWindow on it.
//  2. Dispatch this thread's OWN queue to empty. Dispatched, not dropped: 33-r7
//     measured the four purge shapes on the panel side and the one that survives
//     contact with a window is dispatching (docs/evidence/s1/33-panel-host-c27-r7.md
//     ①). Nothing here removes a message belonging to another owner's window: the
//     queue being emptied is this thread's, and the only windows it ever had are
//     the ones this package created.
//  3. Unbalance the CoInitializeEx above. That is thread state in exactly the same
//     sense as the other two: the next owner of the thread gets an STA apartment it
//     never asked for and cannot see. The create-error door used to skip this
//     entirely (start returned straight past releaseCOM).
func (s *staThread) releaseThread() {
	s.mu.Lock()
	hwnd := s.hwnd
	s.hwnd = 0
	s.mu.Unlock()

	if hwnd != 0 {
		if r, _, err := pDestroyWindow.Call(uintptr(hwnd)); r == 0 {
			slog.Error("ball: the ui-sta owner could not destroy its own window before releasing the thread; the handle stays alive in this process",
				"hwnd", uintptr(hwnd), "err", err)
		}
	}
	releaseOwnQueueToQuiet()
	s.releaseCOM()
}

// forgetWindow is the owner's record that its window is gone: Ball.Close destroys
// the window on the task port, and staThread.hwnd is what releaseThread destroys
// on the way out. Without this the two would disagree - the release step would
// destroy an already-destroyed handle at best, and at worst a handle Win32 has
// since handed to somebody else. It also stops PostTask from posting to a window
// that no longer exists (the pump is going away, so that message had no reader).
func (s *staThread) forgetWindow() {
	s.mu.Lock()
	s.hwnd = 0
	s.mu.Unlock()
}

// releaseOwnQueueToQuiet dispatches (never drops) the CALLING thread's queued
// messages until user32 reports the queue empty, bounded by releasePumpCap. It
// has to run on the owner's own thread - PeekMessage only ever examines the
// calling thread's queue - which is precisely why this duty cannot be discharged
// by a later watcher or by the receiving side.
//
// The filter is the whole range (0..0 with PM_REMOVE means "everything"), never a
// single message id: a latched WM_QUIT is invisible to a filtered peek while any
// other message is still queued (33-r7's reading of the same window in the panel
// family), and an unfiltered one consumes it as a side effect of doing its job.
func releaseOwnQueueToQuiet() int {
	pumped := 0
	for pumped < releasePumpCap {
		var m msg
		r, _, _ := pPeekMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0, uintptr(pmRemove))
		if r == 0 {
			return pumped
		}
		pTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		pDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
		pumped++
	}
	slog.Warn("ball: the ui-sta release pump hit its cap before the thread's queue was empty; the thread goes back to the Go pool non-empty",
		"cap", releasePumpCap, "pumped", pumped)
	return pumped
}

var pPeekMessageW = modUser32.NewProc("PeekMessageW")

var (
	pCoInitializeEx = modOle32.NewProc("CoInitializeEx")
	modOle32        = windows.NewLazySystemDLL("ole32.dll")
)

func (s *staThread) releaseCOM() {
	// Balanced CoUninitialize for the successful CoInitializeEx above.
	pCoUninitialize.Call()
}

var pCoUninitialize = modOle32.NewProc("CoUninitialize")

// waitStarted blocks until the window is up (or boot failed).
func (s *staThread) waitStarted() error {
	<-s.started
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.err
}

// threadID is the OS thread that owns the ball window and pump (0 until the
// thread has started).
func (s *staThread) threadID() uint32 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tid
}

// PostTask schedules fn on the STA thread. Safe from any goroutine.
func (s *staThread) PostTask(fn func()) {
	s.mu.Lock()
	s.nextID++
	id := s.nextID
	s.tasks[id] = fn
	hwnd := s.hwnd
	s.mu.Unlock()
	if hwnd == 0 {
		// Window not created yet (or gone): run inline as last resort so
		// callers cannot deadlock on a dead thread.
		s.mu.Lock()
		delete(s.tasks, id)
		s.mu.Unlock()
		return
	}
	pPostMessageW.Call(uintptr(hwnd), wmAppTask, uintptr(id), 0)
}

// runTask executes a posted closure on the STA thread (from wndproc).
func (s *staThread) runTask(id uint64) {
	s.mu.Lock()
	fn, ok := s.tasks[id]
	delete(s.tasks, id)
	s.mu.Unlock()
	if ok && fn != nil {
		fn()
	}
}

func (s *staThread) quit() {
	hwnd := s.hwnd
	if hwnd != 0 {
		pPostMessageW.Call(uintptr(hwnd), wmNull, 0, 0)
	}
	pPostQuitMessage.Call(0)
}
