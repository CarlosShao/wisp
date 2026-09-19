//go:build windows

package proc

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/windows"
)

// testHelper is a re-executed instance of this test binary running
// TestHelperProcess in a given mode.
type testHelper struct {
	cmd   *exec.Cmd
	ready io.ReadCloser
}

// startHelper spawns the helper via job (nil = plain start) and waits for
// its READY handshake, proving it reached real Go code (not just process
// creation).
func startHelper(t *testing.T, job *JobScope, mode string) *testHelper {
	t.Helper()
	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable: %v", err)
	}
	cmd := exec.Command(exe, "-test.run=^TestHelperProcess$", "-test.count=1")
	cmd.Env = append(os.Environ(), "WISP_HELPER="+mode)
	ready, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("StdoutPipe: %v", err)
	}
	h := &testHelper{cmd: cmd, ready: ready}
	if job != nil {
		if _, err := job.StartInJob(cmd); err != nil {
			t.Fatalf("StartInJob(%s): %v", mode, err)
		}
	} else {
		if err := cmd.Start(); err != nil {
			t.Fatalf("helper start(%s): %v", mode, err)
		}
	}
	buf := make([]byte, 8)
	n, err := ready.Read(buf)
	if err != nil || !strings.HasPrefix(string(buf[:n]), "READY") {
		t.Fatalf("helper %s did not signal READY (n=%d err=%v)", mode, n, err)
	}
	return h
}

// TestHelperProcess is not a real test in the parent process: the tests below
// re-execute the test binary with WISP_HELPER set (go test idiom).
func TestHelperProcess(t *testing.T) {
	switch os.Getenv("WISP_HELPER") {
	case "job-child":
		// Announce liveness, then stay alive for up to 30s; the JobScope
		// tests kill us earlier via the Job. Sleep in slices so a kill
		// lands promptly.
		fmt.Println("READY")
		for i := 0; i < 600; i++ {
			time.Sleep(50 * time.Millisecond)
		}
	case "single-instance-holder":
		// Acquire the single-instance mutex named by the parent, print
		// READY, wait for the activation event (max 30s), print ACTIVATED
		// and exit - the ticket 03 second-instance acceptance.
		si, err := AcquireSingleInstance(
			os.Getenv("WISP_HELPER_MUTEX"), os.Getenv("WISP_HELPER_EVENT"))
		if err != nil {
			fmt.Printf("ACQUIRE-FAILED: %v\n", err)
			os.Exit(3)
		}
		fmt.Println("READY")
		if code, err := windows.WaitForSingleObject(si.ActivateEvent(), 30_000); err != nil || code != windows.WAIT_OBJECT_0 {
			fmt.Printf("WAIT-FAILED: code=%d err=%v\n", code, err)
			os.Exit(4)
		}
		fmt.Println("ACTIVATED")
		_ = si.Release()
	default:
		t.Skip("helper process mode not set")
	}
}

// waitGone waits for the child to exit; fails if it survives the deadline.
func (h *testHelper) waitGone(t *testing.T, within time.Duration) {
	t.Helper()
	done := make(chan error, 1)
	go func() { done <- h.cmd.Wait() }()
	select {
	case <-done:
	case <-time.After(within):
		t.Fatalf("helper survived for %v after Job close", within)
	}
}

// TestJobScopeKillsChildOnClose is the C30 acceptance: a child spawned into
// the Job dies when the Job closes, and TreePrivateBytes returns sane numbers
// (a live Go child has at least a few MB of private bytes).
func TestJobScopeKillsChildOnClose(t *testing.T) {
	job, err := OpenJobScope()
	if err != nil {
		t.Fatalf("OpenJobScope: %v", err)
	}
	t.Cleanup(func() { _ = job.Close() })

	h := startHelper(t, job, "job-child")

	// Sane private bytes for a live Go child: at least 1MB (Go runtime maps
	// several MB), below 8GB (plausibility ceiling vs Task Manager).
	deadline := time.Now().Add(5 * time.Second)
	var pb int64
	for {
		pb, err = job.TreePrivateBytes()
		if err == nil && pb >= 1<<20 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("TreePrivateBytes = %d (err %v), want >= 1MB for a live Go child", pb, err)
		}
		time.Sleep(50 * time.Millisecond)
	}
	if pb > 8<<30 {
		t.Fatalf("TreePrivateBytes = %d, implausibly large", pb)
	}
	if n, err := job.TreeProcessCount(); err != nil || n < 1 {
		t.Fatalf("TreeProcessCount = %d (err %v), want >= 1", n, err)
	}

	// Closing the Job must kill the child (KILL_ON_JOB_CLOSE).
	if err := job.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	h.waitGone(t, 5*time.Second)

	// Post-close behavior: closed flag set, queries fail, Close idempotent.
	if !job.Closed() {
		t.Fatal("Closed() = false after Close")
	}
	if _, err := job.TreePrivateBytes(); err == nil {
		t.Fatal("TreePrivateBytes after Close must fail")
	}
	if err := job.Close(); err != nil {
		t.Fatalf("second Close = %v, want nil (idempotent)", err)
	}
}

// TestJobScopeTreeAccounting covers the multi-child tree: both children are
// accounted; an empty job reports zero.
func TestJobScopeTreeAccounting(t *testing.T) {
	job, err := OpenJobScope()
	if err != nil {
		t.Fatalf("OpenJobScope: %v", err)
	}
	t.Cleanup(func() { _ = job.Close() })

	if n, err := job.TreeProcessCount(); err != nil || n != 0 {
		t.Fatalf("empty job: TreeProcessCount = %d (err %v), want 0", n, err)
	}
	if pb, err := job.TreePrivateBytes(); err != nil || pb != 0 {
		t.Fatalf("empty job: TreePrivateBytes = %d (err %v), want 0", pb, err)
	}

	children := []*testHelper{
		startHelper(t, job, "job-child"),
		startHelper(t, job, "job-child"),
	}

	deadline := time.Now().Add(5 * time.Second)
	var n int
	for {
		n, err = job.TreeProcessCount()
		if err == nil && n >= 2 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("TreeProcessCount = %d (err %v), want 2", n, err)
		}
		time.Sleep(50 * time.Millisecond)
	}

	if err := job.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	for i, h := range children {
		h.waitGone(t, 5*time.Second)
		_ = i
	}
}

// TestStartInJobRejectsBadCommand verifies the escape-proof contract: a
// command that fails to start does not leave a half-registered child.
func TestStartInJobRejectsBadCommand(t *testing.T) {
	job, err := OpenJobScope()
	if err != nil {
		t.Fatalf("OpenJobScope: %v", err)
	}
	t.Cleanup(func() { _ = job.Close() })

	cmd := exec.Command("wisp-this-binary-does-not-exist.exe")
	if _, err := job.StartInJob(cmd); err == nil {
		t.Fatal("StartInJob with a nonexistent binary must fail")
	}
	if n, err := job.TreeProcessCount(); err != nil || n != 0 {
		t.Fatalf("job must stay empty after failed start: %d (err %v)", n, err)
	}
}
