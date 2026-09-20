//go:build windows

package ball

// The ui-sta STA thread (D38a): ONE thread owns all UI COM objects (the D2D
// factory, render targets, the ball window, and later the panel). Everything
// UI-bound runs or posts here. Spawned through the observe registry under
// its frozen roster name "ui-sta"; goroutines elsewhere must never touch
// window/COM state directly - they PostTask.

import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/observe"
)

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

	s.releaseCOM()
}

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
