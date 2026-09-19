//go:build windows

package proc

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/windows"
)

// The single-instance tests need two REAL processes of the same session:
// the parent test process and a helper subprocess (WISP_HELPER=single-instance-holder
// handled in TestHelperProcess).

// startHolder starts the helper that acquires the given mutex and waits for
// the activation event, then returns its stdout pipe.
func startHolder(t *testing.T, mutex, event string) (*exec.Cmd, io.ReadCloser) {
	t.Helper()
	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable: %v", err)
	}
	cmd := exec.Command(exe, "-test.run=^TestHelperProcess$", "-test.count=1")
	cmd.Env = append(os.Environ(),
		"WISP_HELPER=single-instance-holder",
		"WISP_HELPER_MUTEX="+mutex,
		"WISP_HELPER_EVENT="+event,
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("StdoutPipe: %v", err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatalf("start holder: %v", err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill() })
	return cmd, stdout
}

// TestSingleInstanceSecondExitsAfterSignallingFirst is the D42#7 acceptance:
// within one session, a second launch detects the running instance, signals
// it via the activation event, and exits; the mutex frees when the first
// instance goes away.
func TestSingleInstanceSecondExitsAfterSignallingFirst(t *testing.T) {
	mutex := fmt.Sprintf(`Local\wisp-test-single-instance-%d-%d`, os.Getpid(), time.Now().UnixNano())
	event := mutex + ".activate"

	cmd, stdout := startHolder(t, mutex, event)
	waitForLine(t, stdout, "READY", 15*time.Second)

	// Second instance in the SAME session (this test process) must detect
	// the running one.
	si, err := AcquireSingleInstance(mutex, event)
	if !errors.Is(err, ErrAlreadyRunning) {
		if si != nil {
			_ = si.Release()
		}
		t.Fatalf("second AcquireSingleInstance err = %v, want ErrAlreadyRunning", err)
	}

	// The second launch signals the first and exits.
	if err := SignalExistingInstance(event); err != nil {
		t.Fatalf("SignalExistingInstance: %v", err)
	}
	waitForLine(t, stdout, "ACTIVATED", 10*time.Second)

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case <-done:
		// Helper exited after signalling handling.
	case <-time.After(10 * time.Second):
		t.Fatal("holder did not exit after activation")
	}

	// With the first instance gone, the mutex must be acquirable again.
	si, err = AcquireSingleInstance(mutex, event)
	if err != nil {
		t.Fatalf("acquire after holder exit: %v", err)
	}
	// Reset + signal round-trip within one instance.
	if err := si.ResetActivation(); err != nil {
		t.Fatalf("ResetActivation: %v", err)
	}
	if err := SignalExistingInstance(event); err != nil {
		t.Fatalf("self-signal: %v", err)
	}
	if code, err := windows.WaitForSingleObject(si.ActivateEvent(), 2000); err != nil || code != windows.WAIT_OBJECT_0 {
		t.Fatalf("activation event not signaled after self-signal: code=%d err=%v", code, err)
	}
	if err := si.Release(); err != nil {
		t.Fatalf("Release: %v", err)
	}
	if err := si.Release(); err != nil {
		t.Fatalf("second Release = %v, want nil (idempotent)", err)
	}
}

// TestAcquireRejectsEmptyNames is a guard against silent global-defaulting.
func TestAcquireRejectsEmptyNames(t *testing.T) {
	if _, err := AcquireSingleInstance("", "x"); err == nil {
		t.Fatal("empty mutex name must fail")
	}
	if _, err := AcquireSingleInstance(`Local\m`, ""); err == nil {
		t.Fatal("empty event name must fail")
	}
}

// waitForLine reads the pipe until the given line shows up or timeout.
func waitForLine(t *testing.T, r io.Reader, want string, within time.Duration) {
	t.Helper()
	found := make(chan bool, 1)
	go func() {
		buf := make([]byte, 1)
		line := make([]byte, 0, 64)
		for {
			n, err := r.Read(buf)
			if n > 0 {
				line = append(line, buf[:n]...)
				if len(line) > 4096 {
					line = line[len(line)-4096:]
				}
				if strings.Contains(string(line), want) {
					found <- true
					return
				}
			}
			if err != nil {
				found <- false
				return
			}
		}
	}()
	select {
	case ok := <-found:
		if !ok {
			t.Fatalf("stream ended before %q appeared", want)
		}
	case <-time.After(within):
		t.Fatalf("%q did not appear within %v", want, within)
	}
}
