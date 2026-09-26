//go:build windows

package main

// Ticket 156 PROBE. This file is NOT part of the shipped package: it lives in
// .scratch/wisp/probes/156/ (which the go tool ignores, dot-directory) and is fed
// to `go test -overlay` as cmd/wisp/zz156probe_windows_test.go, so no probe byte
// ever lands in the tree - same shape as ticket 152's zz152probe.
//
// The one question it exists to answer, in AC#1's own words: reproduce THE SHOT
// as this process's own reading (not as a citation of ticket 152's numbers), and
// size the window in which the OS already knows the subject is dead and our code
// does not - with a named read count and milliseconds.
//
// Production wiring, replayed and not approximated:
//   * startSubject() only ever starts the child (Job.StartInJob -> cmd.Start), so
//     cmd.ProcessState is nil for the whole life of the observer's loop;
//   * the first cmd.Wait() in the file is inside stop(), which runOutOfTree defers
//     past collectReportWithin (dispatch premise 2, re-measured at this anchor);
//   * so the loop's own predicate s.exited() can only answer "false" about a child
//     the kernel has already finished, and the give-up the operator reads is the
//     budget sentence - the neighbour's sentence - instead of "exited (code N)".
//
// Every cell spawns a REAL child (exec.Command(os.Args[0], ...) re-entering this
// test binary in a child role: a real pid, a real partial file, a real exit code
// 7) and drives the REAL collectReportWithin over the REAL os.ReadFile, wrapped by
// a counting seam so "how many times did it look" is a reading too.
//
// The external witness is deliberately a DIFFERENT API from the one the fix may
// use: OpenProcess(pid) + GetExitCodeProcess on a freshly opened handle, asked
// from the pid, not the exec package's own handle. The fix's "OS says dead" is
// therefore never corroborated by the same call it makes.
//
// Readings print as `P156|<cell>|<key>=<value>` lines so they can be quoted
// verbatim from the raw log.

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"golang.org/x/sys/windows"
)

const (
	p156EnvRole  = "WISP156_CHILD_ROLE"
	p156EnvOut   = "WISP156_CHILD_OUT"
	p156EnvBytes = "WISP156_CHILD_BYTES"
	// p156ChildExit is the code the child dies with; the sentence names it.
	p156ChildExit = 7
	// p156StillActive is Win32 STILL_ACTIVE: what GetExitCodeProcess answers while
	// the process object is still running. The witness below is the ONLY place in
	// this probe allowed to read it, and it is not the mechanism the fix uses.
	p156StillActive = 259
	p156Budget      = 2 * time.Second
	p156Poll        = 10 * time.Millisecond
	// p156WitnessWait caps how long the probe waits for the child to really die.
	// It bounds a fixture, not a shipped timeout, and it is a fixture the test
	// FAILS on rather than grades on.
	p156WitnessWait = 20 * time.Second
)

// p156Say prints one reading as a flat, greppable line. The value is always
// formatted by the CALLER with a constant format string, so this file never hands
// a non-constant format to fmt (go vet's printf check runs inside `go test`).
func p156Say(cell, key, value string) {
	fmt.Printf("P156|%s|%s=%s\n", cell, key, value)
}

// p156ChildBody is the child half: a real process that stops mid-report and then
// dies. The bytes it leaves are a proper prefix of the payload writeSLO writes
// (same serializer as the shipped fixture), which is exactly the state a subject
// killed inside its single os.WriteFile leaves behind - ticket 144's case 1
// measures that every one of those prefixes classifies unwritten.
func p156ChildBody(t *testing.T) {
	role := os.Getenv(p156EnvRole)
	out := os.Getenv(p156EnvOut)
	n, err := strconv.Atoi(os.Getenv(p156EnvBytes))
	switch role {
	case "sleep":
		time.Sleep(60 * time.Second) // negative control: a live child
		os.Exit(0)
	case "nofile":
		os.Exit(p156ChildExit) // died before the report file existed at all
	}
	doc := slo144Report(t)
	switch {
	case role != "prefix":
		fmt.Printf("P156|child|unknown_role=%q\n", role)
		os.Exit(6)
	case err != nil || n <= 0 || n >= len(doc):
		fmt.Printf("P156|child|bad_bytes=%d err=%v\n", n, err)
		os.Exit(5)
	}
	if err := os.WriteFile(out, append([]byte(nil), doc[:n]...), 0o644); err != nil {
		fmt.Printf("P156|child|write_failed=%v\n", err)
		os.Exit(4)
	}
	os.Exit(p156ChildExit)
}

// p156OsStatus asks the kernel, over a handle opened fresh from the pid, whether
// the pid is still running and what code it carries. Returns alive / code / how.
func p156OsStatus(pid uint32) (alive bool, code int64, how string) {
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		if errors.Is(err, syscall.Errno(windows.ERROR_INVALID_PARAMETER)) {
			// No process object at all for this pid - the OS knows it is gone.
			return false, -1, "OpenProcess: ERROR_INVALID_PARAMETER (no process object)"
		}
		return true, -1, fmt.Sprintf("OpenProcess: %v (not judged dead)", err)
	}
	defer func() { _ = windows.CloseHandle(h) }()
	var c uint32
	if err := windows.GetExitCodeProcess(h, &c); err != nil {
		return true, -1, fmt.Sprintf("GetExitCodeProcess: %v (not judged dead)", err)
	}
	if c == p156StillActive {
		return true, int64(c), "GetExitCodeProcess: STILL_ACTIVE"
	}
	return false, int64(c), "GetExitCodeProcess: terminated"
}

// p156Spawn starts a real child. Stdout/Stderr are nil exactly as startSubject
// leaves them, so nothing the child prints is mistaken for the observer's console.
func p156Spawn(t *testing.T, path string, prefixBytes int, role string) *exec.Cmd {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestP156SubjectDeathProbe$", "-test.timeout=180s")
	cmd.Env = append(os.Environ(),
		p156EnvRole+"="+role,
		p156EnvOut+"="+path,
		p156EnvBytes+"="+strconv.Itoa(prefixBytes))
	cmd.Stdout, cmd.Stderr = nil, nil
	if err := cmd.Start(); err != nil {
		t.Fatalf("P156 start subject child: %v", err)
	}
	return cmd
}

func TestP156SubjectDeathProbe(t *testing.T) {
	if role := os.Getenv(p156EnvRole); role != "" {
		p156ChildBody(t) // never returns
		return
	}

	doc := slo144Report(t)
	half := len(doc) / 2
	p156Say("fixture", "bytes", fmt.Sprintf("doc=%d child_writes=%d child_exit=%d budget=%s poll=%s",
		len(doc), half, p156ChildExit, p156Budget, p156Poll))

	cells := []struct {
		name   string
		role   string
		reaped bool // Wait() BEFORE the loop: ticket 149's case 13 shape, the one production never uses
	}{
		{name: "prefix-unreaped", role: "prefix"},
		{name: "nofile-unreaped", role: "nofile"},
		{name: "prefix-REAPED-case13shape", role: "prefix", reaped: true},
		{name: "live-child", role: "sleep"},
	}

	for _, c := range cells {
		t.Run(c.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "subject-report.json")
			spawn := time.Now()
			cmd := p156Spawn(t, path, half, c.role)
			pid := uint32(cmd.Process.Pid)
			reaped := false
			defer func() {
				if !reaped {
					_ = cmd.Wait() // never leave a child behind
				}
			}()

			s := &sloSubject{cmd: cmd, pid: pid, dir: dir, outPath: path}
			reads := 0
			s.readReportFile = func(p string) ([]byte, error) {
				reads++
				return os.ReadFile(p)
			}

			// --- the external witness: when does the OS itself know? ---
			osDeadAt := time.Time{}
			witnessPolls := 0
			alive, code, how := true, int64(-1), "not asked"
			for {
				witnessPolls++
				alive, code, how = p156OsStatus(pid)
				if !alive || c.role == "sleep" {
					break
				}
				if time.Since(spawn) > p156WitnessWait {
					t.Fatalf("P156 child %d never reported dead to the OS within %s (%s)",
						pid, p156WitnessWait, how)
				}
				time.Sleep(time.Millisecond)
			}
			if c.role == "sleep" {
				p156Say(c.name, "witness", fmt.Sprintf("alive=%v code=%d polls=%d (%s)",
					alive, code, witnessPolls, "child deliberately still running"))
			} else {
				osDeadAt = time.Now()
				p156Say(c.name, "witness", fmt.Sprintf("os_knows_dead_ms_since_spawn=%d polls=%d code=%d (%s)",
					int64(osDeadAt.Sub(spawn)/time.Millisecond), witnessPolls, code, how))
			}

			if c.reaped {
				werr := cmd.Wait()
				reaped = true
				p156Say(c.name, "wait_before_loop", fmt.Sprintf("%v", werr))
			}

			// --- the loop's OWN predicate, at the instant the OS already knows ---
			p156Say(c.name, "record", fmt.Sprintf("ProcessState_nil=%v exited=%v exit_code=%d",
				s.cmd.ProcessState == nil, s.exited(), s.exitCode()))

			// How long until OUR record agrees with the OS, without anybody reaping?
			// p156AgreePolls counts the asks and p156AgreeMS the elapsed ms; a cell
			// that never agrees reports agree_ms=-1 agree_polls=<n>.
			agreePolls := 0
			var agreeMS int64 = -1
			if !osDeadAt.IsZero() {
				deadline := time.Now().Add(200 * time.Millisecond)
				for {
					agreePolls++
					if s.exited() {
						agreeMS = int64(time.Since(osDeadAt) / time.Millisecond)
						break
					}
					if time.Now().After(deadline) {
						break
					}
					time.Sleep(time.Millisecond)
				}
				p156Say(c.name, "agree", fmt.Sprintf("os_dead_to_code_dead_ms=%d asks=%d", agreeMS, agreePolls))
			}

			start := time.Now()
			_, err := s.collectReportWithin(p156Budget, p156Poll)
			elapsed := time.Since(start)
			sentence := "<NO ERROR - the loop came back green>"
			if err != nil {
				sentence = strings.ReplaceAll(err.Error(), "\n", " ")
			}
			kind := "other"
			switch {
			case strings.Contains(sentence, "exited (code"):
				kind = "DEAD-CHILD sentence (the one that names the cause)"
			case strings.Contains(sentence, "never wrote a complete report within"):
				kind = "BUDGET sentence (the neighbour: names a spent budget, not a corpse)"
			case err == nil:
				kind = "no error"
			}
			p156Say(c.name, "loop", fmt.Sprintf("elapsed_ms=%d reads=%d budget=%s verdict_kind=%s",
				int64(elapsed/time.Millisecond), reads, p156Budget, kind))
			p156Say(c.name, "sentence", sentence)
			p156Say(c.name, "record_after_loop", fmt.Sprintf("exited=%v exit_code=%d ProcessState_nil=%v",
				s.exited(), s.exitCode(), s.cmd.ProcessState == nil))
			if !osDeadAt.IsZero() {
				p156Say(c.name, "window", fmt.Sprintf("os_knew_at_ms=%d_code_still_said_alive_after_loop_ms=%d",
					int64(osDeadAt.Sub(spawn)/time.Millisecond), int64(time.Since(osDeadAt)/time.Millisecond)))
			}
		})
	}
}
