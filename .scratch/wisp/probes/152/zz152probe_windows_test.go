//go:build windows

package main

// Ticket 152 PROBE. This file is NOT part of the shipped package: it lives
// outside the repository and is fed to `go test -overlay` as
// cmd/wisp/zz152probe_windows_test.go, so no probe byte ever lands in the tree.
//
// The one question it exists to answer, in the shape tickets 144/147/149 all
// declined to answer it in:
//
//	case 13 of slo_report_144_windows_test.go reaches the loop's s.exited()
//	branch with (a) a child reaped BEFORE the loop starts - `dead.Run()` fills
//	ProcessState - and (b) the report bytes handed in by scriptedReader. Neither
//	half is how startSubject + collectReport actually wire a subject: startSubject
//	only ever calls cmd.Start (via Job.StartInJob) and the first Wait in the
//	production path is stop(), which is deferred past collectReport.
//
// So what does the REAL loop print when a REAL subject process dies for real,
// with a REAL file on disk read by the REAL os.ReadFile? Every cell below
// spawns a real child process (exec.Command(os.Args[0], ...), a real pid, a real
// exit status 7) and drives the real collectReportWithin. The only things varied
// between cells are the two halves case 13 fixed by hand: whether the child was
// reaped (Wait) before the loop, and whether the reader is the production one.
//
// Readings are printed as `P152|<cell>|<key>=<value>` lines so they can be
// quoted verbatim from the raw log.

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
	p152EnvRole  = "WISP152_CHILD_ROLE"
	p152EnvOut   = "WISP152_CHILD_OUT"
	p152EnvBytes = "WISP152_CHILD_BYTES"
	// p152ChildExit is the exit code the child dies with; the sentence names it.
	p152ChildExit   = 7
	p152StillActive = 259 // Win32 STILL_ACTIVE, what GetExitCodeProcess says while alive
	p152Poll        = 20 * time.Millisecond
	p152Budget      = 4 * time.Second
	p152ChildWait   = 20 * time.Second
)

// p152Say prints one reading as a flat, greppable line. The value is always
// formatted by the CALLER with a constant format string, so this function never
// hands a non-constant format to fmt (go vet's printf check runs inside
// `go test` and would otherwise refuse the build).
func p152Say(cell, key, value string) {
	fmt.Printf("P152|%s|%s=%s\n", cell, key, value)
}

// p152ChildBody is the child half: a real process that stops mid-report. The
// bytes it hands the file are a proper prefix of the payload writeSLO writes
// (same serializer, slo144Report), which is exactly the state a subject that
// died inside its single os.WriteFile leaves behind - case 1 of the shipped file
// measures that every one of those prefixes classifies unwritten. The child then
// exits with a code, so the writer is gone for good: no tail ever arrives.
func p152ChildBody(t *testing.T) {
	role := os.Getenv(p152EnvRole)
	out := os.Getenv(p152EnvOut)
	n, err := strconv.Atoi(os.Getenv(p152EnvBytes))
	doc := slo144Report(t)
	switch {
	case err != nil || n <= 0 || n >= len(doc):
		fmt.Printf("P152|child|bad_bytes=%d err=%v\n", n, err)
		os.Exit(5)
	case role == "nofile":
		os.Exit(p152ChildExit) // died before the report file existed at all
	case role != "prefix":
		fmt.Printf("P152|child|unknown_role=%q\n", role)
		os.Exit(6)
	}
	if err := os.WriteFile(out, append([]byte(nil), doc[:n]...), 0o644); err != nil {
		fmt.Printf("P152|child|write_failed=%v\n", err)
		os.Exit(4)
	}
	os.Exit(p152ChildExit)
}

// p152OsStatus asks the OS - not Go's exec package - whether the pid is still
// running and what exit code it carries. This is the external witness that a
// cell's child is really dead while the loop's own predicate says otherwise.
func p152OsStatus(pid uint32) (alive bool, code int64, how string) {
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		if errors.Is(err, syscall.Errno(windows.ERROR_INVALID_PARAMETER)) {
			return false, -1, "OpenProcess: ERROR_INVALID_PARAMETER (no process object)"
		}
		return true, -1, fmt.Sprintf("OpenProcess: %v (not judged dead)", err)
	}
	defer func() { _ = windows.CloseHandle(h) }()
	var c uint32
	if err := windows.GetExitCodeProcess(h, &c); err != nil {
		return true, -1, fmt.Sprintf("GetExitCodeProcess: %v (not judged dead)", err)
	}
	if c == p152StillActive {
		return true, int64(c), "GetExitCodeProcess: STILL_ACTIVE"
	}
	return false, int64(c), "GetExitCodeProcess: terminated"
}

// p152Spawn starts a real child. Stdout/Stderr are nil exactly as startSubject
// leaves them, so nothing the child prints can be mistaken for the observer's
// console.
func p152Spawn(t *testing.T, path string, prefixBytes int, role string) *exec.Cmd {
	t.Helper()
	cmd := exec.Command(os.Args[0],
		"-test.run=^TestP152SubjectDeathProbe$", "-test.timeout=120s")
	cmd.Env = append(os.Environ(),
		p152EnvRole+"="+role,
		p152EnvOut+"="+path,
		p152EnvBytes+"="+strconv.Itoa(prefixBytes))
	cmd.Stdout, cmd.Stderr = nil, nil
	if err := cmd.Start(); err != nil {
		t.Fatalf("P152 start subject child: %v", err)
	}
	return cmd
}

// p152WaitDeath blocks until the OS reports the child terminated, so no cell can
// be accused of asking exited() about a child that was still alive.
func p152WaitDeath(t *testing.T, pid uint32) (int64, string) {
	t.Helper()
	deadline := time.Now().Add(p152ChildWait)
	for {
		alive, code, how := p152OsStatus(pid)
		if !alive {
			return code, how
		}
		if time.Now().After(deadline) {
			t.Fatalf("P152 child %d never reported dead to the OS within %s (%s)",
				pid, p152ChildWait, how)
		}
		time.Sleep(p152Poll)
	}
}

// p152Cell is one wiring of the loop: a real child, a real file (unless the
// reader is scripted), and the two halves case 13 set by hand.
type p152Cell struct {
	name   string
	role   string // "prefix" | "nofile"
	reaped bool   // Wait() BEFORE the loop, as case 13's dead.Run() does
	reader string // "production" (nil seam) | "counting" | "scripted"
}

func TestP152SubjectDeathProbe(t *testing.T) {
	if role := os.Getenv(p152EnvRole); role != "" {
		p152ChildBody(t) // never returns
		return
	}

	doc := slo144Report(t)
	half := len(doc) / 2
	p152Say("fixture", "bytes", fmt.Sprintf("%d child_writes=%d child_exit=%d", len(doc), half, p152ChildExit))

	cells := []p152Cell{
		// PRODUCTION SHAPE: nothing reaps the child before the loop, and the
		// loop reads the real file through os.ReadFile.
		{name: "prod-prefix-unreaped", role: "prefix", reader: "production"},
		{name: "prod-nofile-unreaped", role: "nofile", reader: "production"},
		// The same, with a counting wrapper AROUND os.ReadFile (no scripted
		// bytes) so the number of real readings is a reading too.
		{name: "prod-prefix-unreaped-count", role: "prefix", reader: "counting"},
		// The two halves case 13 fixes by hand, one at a time.
		{name: "prod-prefix-REAPED", role: "prefix", reaped: true, reader: "production"},
		{name: "prod-nofile-REAPED", role: "nofile", reaped: true, reader: "production"},
		{name: "case13-shape-scripted-REAPED", role: "prefix", reaped: true, reader: "scripted"},
		{name: "unreaped-scripted", role: "prefix", reader: "scripted"},
	}

	for _, c := range cells {
		t.Run(c.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "subject-report.json")
			cmd := p152Spawn(t, path, half, c.role)
			pid := uint32(cmd.Process.Pid)
			reaped := false
			defer func() {
				if !reaped {
					_ = cmd.Wait() // never leave a child behind
				}
			}()

			s := &sloSubject{cmd: cmd, pid: pid, dir: dir, outPath: path}

			var reads int
			var scripted *scriptedReader
			switch c.reader {
			case "counting":
				s.readReportFile = func(p string) ([]byte, error) {
					reads++
					return os.ReadFile(p)
				}
			case "scripted":
				scripted = &scriptedReader{t: t, repeatLast: true,
					got: []func() ([]byte, error){scriptBytes(doc[:half])}}
				s.readReportFile = scripted.read
			}

			if c.reaped {
				werr := cmd.Wait()
				reaped = true
				p152Say(c.name, "wait_before_loop", fmt.Sprintf("%v", werr))
			} else {
				code, how := p152WaitDeath(t, pid)
				p152Say(c.name, "os_says", fmt.Sprintf("alive=false code=%d (%s)", code, how))
			}

			// The loop's OWN predicate, at the instant the OS has already
			// reported the child terminated.
			p152Say(c.name, "exited_before_loop", fmt.Sprintf("%v", s.exited()))
			p152Say(c.name, "exitcode_before_loop", fmt.Sprintf("%d", s.exitCode()))

			start := time.Now()
			_, err := s.collectReportWithin(p152Budget, p152Poll)
			elapsed := time.Since(start)
			if err == nil {
				p152Say(c.name, "sentence", "<NO ERROR - the loop came back green>")
			} else {
				p152Say(c.name, "sentence", strings.ReplaceAll(err.Error(), "\n", " "))
			}
			p152Say(c.name, "elapsed", fmt.Sprintf("%s budget=%s", elapsed.Round(p152Poll), p152Budget))
			p152Say(c.name, "reads", fmt.Sprintf("counting_seam=%d scripted_seam=%s",
				reads, func() string {
					if scripted == nil {
						return "n/a"
					}
					return strconv.Itoa(scripted.calls)
				}()))
			p152Say(c.name, "exited_after_loop", fmt.Sprintf("%v", s.exited()))
			p152Say(c.name, "exited_after_wait", func() string {
				if !reaped {
					_ = cmd.Wait()
					reaped = true
				}
				return fmt.Sprintf("%v code=%d", s.exited(), s.exitCode())
			}())

			onDisk, rerr := os.ReadFile(path)
			obs := readSubjectReport(onDisk)
			p152Say(c.name, "file_on_disk", fmt.Sprintf("bytes=%d read_err=%v classified=%s (%s)",
				len(onDisk), rerr, obs.state, obs.summary()))
		})
	}
}
