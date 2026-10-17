//go:build windows

package proc

import (
	"sync"
)

// The stop-request bit of ticket 228 AC#11 (form 乙-1, orchestrator ruling A493).
//
// WHY THIS EXISTS. Until now the only thing that could end RunEventLoop was a
// console signal: signal.NotifyContext(os.Interrupt, SIGTERM) inside this
// function (boot_windows.go). cmd/wisp's tray "Exit" item therefore had a record
// and no executor (resident_ball_windows.go's recordTrayExit printed
// outcome=ignored), and ticket 244's -H=windowsgui cut would have left a
// double-clicked resident process with no way out that keeps the D38(e) order.
// The bit is what AC#11 adds: one request, recorded, and the loop ends itself.
//
// WHY A BIT AND NOT A WAITER. The tray callback runs ON the ball's ui-sta thread
// (internal/ball/ball_windows.go's wndproc -> Ball.fire, contract "quick and
// non-blocking"), and every teardown step in this process waits on something
// (Ball.Close is PostTask + <-done), so a callback that tried to close the
// process would wait on the thread it is standing on. Recording a request and
// returning is the only shape that cannot self-deadlock; it is the same shape
// cmd/wisp's vetoByEsc already uses for the Esc key.
//
// WHAT IT DOES NOT TOUCH. No step of D38(e) is added, removed or reordered: the
// sequence still walks 1..10 in RunShutdownSequence (shutdown.go), the ten-step
// chain still hangs off runResident returning normally, and hookableRoster stays
// the closed eight slots it is (shutdown_hooks.go). A new TRIGGER is not a new
// STEP - that is the reading A493 ruled, and the doc below states the trigger
// roster instead of leaving it stale.
type stopRequest struct {
	mu  sync.Mutex
	why string
	// wake is closed by the first request. It is a latch, not an event: a request
	// that arrives while nobody is selecting yet is still there when the select
	// happens, which is exactly what the boot window needs.
	wake chan struct{}
}

// newStopRequest builds the bit. Boot calls it, so every runtime this package
// hands out can be asked to stop from the moment it exists.
func newStopRequest() *stopRequest {
	return &stopRequest{wake: make(chan struct{})}
}

// request records the first request and ignores every later one: the loop can
// only end once, and the reason the operator gave the first time is the one the
// audit trail should carry.
func (s *stopRequest) request(why string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.wake == nil {
		// A Runtime assembled without Boot (a hand-built literal in a test) still
		// must not panic on the thread that owns a click.
		s.wake = make(chan struct{})
	}
	if s.why != "" {
		return false
	}
	s.why = why
	close(s.wake)
	return true
}

// signal returns the channel that is closed once a request is on record. The
// zero value answers nil, and a nil channel simply never readies a select.
func (s *stopRequest) signal() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.wake
}

// requested reports whether a request is on record.
func (s *stopRequest) requested() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.why != ""
}

// reason answers why this runtime was asked to stop, or the empty string when
// nobody asked.
func (s *stopRequest) reason() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.why
}

// RequestStop records a stop request against this runtime and returns
// immediately. It waits on nothing and touches no Win32, which is what makes it
// legal to call from a ball callback running on the ui-sta thread; the loop
// itself ends the process, through the same normal return that every ordered
// exit has always used.
//
// It returns true when this call was the one the loop will act on, and false
// when a request was already on record (or when rt is nil, which is legal for
// the injected-callback shape this seam exists for).
func (rt *Runtime) RequestStop(why string) bool {
	if rt == nil {
		return false
	}
	if rt.stopReq == nil {
		rt.stopReq = newStopRequest()
	}
	return rt.stopReq.request(why)
}

// StopSignal is the closed-on-request channel of this runtime's stop bit, and
// the thing a waiter selects on. A runtime with no bit answers nil, and a nil
// channel readies no select case, so "nothing can ask this process to stop" is
// the behaviour instead of a panic.
func (rt *Runtime) StopSignal() <-chan struct{} {
	if rt == nil || rt.stopReq == nil {
		return nil
	}
	return rt.stopReq.signal()
}

// StopRequested reports whether a stop request is already on record. It reads a
// latch, so it answers the same question in the boot window as inside the loop.
func (rt *Runtime) StopRequested() bool {
	if rt == nil || rt.stopReq == nil {
		return false
	}
	return rt.stopReq.requested()
}

// StopReason is the reason the first request carried, or "" when nobody asked.
// The event loop puts it into the string it returns, which is the sentence
// cmd/wisp prints before it runs the D38(e) order.
func (rt *Runtime) StopReason() string {
	if rt == nil || rt.stopReq == nil {
		return ""
	}
	return rt.stopReq.reason()
}
