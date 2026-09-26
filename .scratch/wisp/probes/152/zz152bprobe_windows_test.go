//go:build windows

package main

// Ticket 152 PROBE #2 - the SECOND outlet of the same predicate.
//
// waitReady() (slo_windows.go:514-529) asks the very same s.exited() before
// anything has ever called Wait on the subject, so its "exited early (code N)
// before reporting ready" sentence is claimed unreachable by the same root
// cause. The code census says so; this probe measures it: a real child that
// exits 7 WITHOUT writing the readiness marker, then the real waitReady() -
// which waits the full subjectReadyBudget (60s, a constant this probe must not
// and does not change).
//
// Same overlay mechanics as probe #1; nothing lands in the repository.

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

const (
	p152bEnvKillWithoutReady = "WISP152B_CHILD"
	p152bBudget              = 75 * time.Second // > subjectReadyBudget (60s)
)

// p152bSay keeps the format string constant at every call site (go vet's printf
// check runs inside `go test`, so a non-constant format would refuse the build).
func p152bSay(key, value string) {
	fmt.Printf("P152B|%s|%s\n", key, value)
}

func TestP152BWaitReadyProbe(t *testing.T) {
	if os.Getenv(p152bEnvKillWithoutReady) == "1" {
		// child: never writes the ready marker, dies with a code
		os.Exit(7)
		return
	}

	dir := t.TempDir()
	ready := filepath.Join(dir, "ready")
	out := filepath.Join(dir, "subject-report.json")

	cmd := exec.Command(os.Args[0], "-test.run=^TestP152BWaitReadyProbe$", "-test.timeout=180s")
	cmd.Env = append(os.Environ(), p152bEnvKillWithoutReady+"=1")
	cmd.Stdout, cmd.Stderr = nil, nil
	if err := cmd.Start(); err != nil {
		t.Fatalf("P152B start: %v", err)
	}
	defer func() { _ = cmd.Wait() }()
	pid := uint32(cmd.Process.Pid)

	alive, code, how := p152OsStatus(pid)
	for tries := 0; alive && tries < 100; tries++ {
		time.Sleep(20 * time.Millisecond)
		alive, code, how = p152OsStatus(pid)
	}
	p152bSay("child", fmt.Sprintf("os_alive=%v os_code=%d how=%s pid=%d", alive, code, how, pid))

	s := &sloSubject{cmd: cmd, pid: pid, dir: dir, readyPath: ready, outPath: out}
	p152bSay("predicate", fmt.Sprintf("exited_before_waitReady=%v exitCode=%d ready_file_exists=%v",
		s.exited(), s.exitCode(), fileHas(ready, "ready=1")))

	start := time.Now()
	err := s.waitReady()
	p152bSay("waitReady", fmt.Sprintf("elapsed=%s err=%v", time.Since(start).Round(10*time.Millisecond), err))
	p152bSay("after", fmt.Sprintf("exited=%v exitCode=%d", s.exited(), s.exitCode()))
	_ = cmd.Wait()
	p152bSay("after_reap", fmt.Sprintf("exited=%v exitCode=%d", s.exited(), s.exitCode()))
}
