//go:build windows

package proc

import (
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sys/windows"
)

// Per-session single instance (D42#7, D38): the mutex lives in the `Local\`
// namespace, so it is scoped to the login session - a second user (or RDP
// session) can run their own Wisp, but within one session only one instance
// runs. A second launch signals the first instance's activation event
// (bringing the ball to the front is the ball's job, ticket 07) and exits.

// ErrAlreadyRunning is returned by AcquireSingleInstance when another Wisp
// process of this session owns the mutex.
var ErrAlreadyRunning = errors.New("proc: another wisp instance is running in this session")

// SingleInstance owns the mutex and the activation event of the running
// instance. Create with AcquireSingleInstance.
type SingleInstance struct {
	MutexName string
	EventName string

	mu     sync.Mutex
	mutex  windows.Handle
	event  windows.Handle
	closed bool
}

// AcquireSingleInstance creates and owns the named mutex, and creates the
// activation event (reusing a stale one from a crashed instance is safe: the
// first act of the owner is a Reset). It returns ErrAlreadyRunning when the
// mutex already exists in this session. Cleanup happens at process death by
// OS fiat; Release exists for tests and clean hand-back.
func AcquireSingleInstance(mutexName, eventName string) (*SingleInstance, error) {
	if mutexName == "" || eventName == "" {
		return nil, fmt.Errorf("proc: single-instance: empty mutex/event name")
	}
	mname, err := windows.UTF16PtrFromString(mutexName)
	if err != nil {
		return nil, fmt.Errorf("proc: single-instance mutex name: %w", err)
	}
	// initialOwner=true removes the ownership race window.
	h, err := windows.CreateMutex(nil, true, mname)
	switch {
	case h != 0 && err != nil:
		// x/sys returns (valid handle, ERROR_ALREADY_EXISTS) when it exists.
		_ = windows.CloseHandle(h)
		return nil, fmt.Errorf("proc: mutex %q: %w", mutexName, ErrAlreadyRunning)
	case err != nil:
		return nil, fmt.Errorf("proc: create mutex %q: %w", mutexName, err)
	}

	ename, err := windows.UTF16PtrFromString(eventName)
	if err != nil {
		_ = windows.CloseHandle(h)
		return nil, fmt.Errorf("proc: single-instance event name: %w", err)
	}
	// Manual-reset, initially nonsignaled.
	ev, err := windows.CreateEvent(nil, 1, 0, ename)
	if ev != 0 && err != nil {
		// Event existed (stale from a crashed instance): reuse it.
		err = nil
	}
	if err != nil {
		_ = windows.CloseHandle(h)
		return nil, fmt.Errorf("proc: create activation event %q: %w", eventName, err)
	}

	return &SingleInstance{
		MutexName: mutexName,
		EventName: eventName,
		mutex:     h,
		event:     ev,
	}, nil
}

// ActivateEvent returns the activation event handle for the owner's wait
// loop (ball, ticket 07). Handle with care: do not close it.
func (s *SingleInstance) ActivateEvent() windows.Handle {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.event
}

// ResetActivation clears the activation signal (manual-reset event) after
// the owner handled it, so the next second-launch can signal again.
func (s *SingleInstance) ResetActivation() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return fmt.Errorf("proc: single-instance already released")
	}
	if err := windows.ResetEvent(s.event); err != nil {
		return fmt.Errorf("proc: reset activation event: %w", err)
	}
	return nil
}

// Release closes both handles. The OS releases the mutex anyway when the
// owning process exits.
func (s *SingleInstance) Release() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	if s.event != 0 {
		_ = windows.CloseHandle(s.event)
		s.event = 0
	}
	if s.mutex != 0 {
		_ = windows.CloseHandle(s.mutex)
		s.mutex = 0
	}
	return nil
}

// SignalExistingInstance signals the activation event of the running
// instance. Used by the second launch right before exiting.
func SignalExistingInstance(eventName string) error {
	ename, err := windows.UTF16PtrFromString(eventName)
	if err != nil {
		return fmt.Errorf("proc: activation event name: %w", err)
	}
	h, err := windows.OpenEvent(windows.EVENT_MODIFY_STATE, false, ename)
	if err != nil {
		return fmt.Errorf("proc: open activation event %q: %w", eventName, err)
	}
	defer windows.CloseHandle(h)
	if err := windows.SetEvent(h); err != nil {
		return fmt.Errorf("proc: signal activation event %q: %w", eventName, err)
	}
	return nil
}
