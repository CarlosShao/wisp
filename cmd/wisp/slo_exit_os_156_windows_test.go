//go:build windows

package main

// Ticket 156, cases 15-19. Q-57's ruling 乙, landed: `wisp slo`'s predicate for
// "is my subject gone" asks the OPERATING SYSTEM, not "have we happened to reap
// it". The three pre-existing callers of that predicate (waitReady's give-up, the
// report loop's give-up, stop's kill guard) are the whole blast radius, and two of
// them are the sentences an operator reads when the instrument goes red.
//
// What this file refuses to do is prove the fix with the shape that made the bug
// invisible. Ticket 149's case 13 reaches the loop's exited branch by handing the
// loop a child it had ALREADY reaped (`dead.Run()` fills ProcessState before the
// loop starts) - ticket 152's acceptance run named that wiring as production-
// nonexistent, because startSubject only ever STARTS a subject and the first
// cmd.Wait() in the file is inside stop(), deferred past the loop. So every cell
// below keeps the child UNREAPED across the assertion, says so in a visible way
// (ProcessState == nil is asserted in the same breath as the verdict), and asks
// the OS the same question over a DIFFERENT door (a handle opened fresh from the
// pid) before asking our own predicate. A future "fix" that reaps inside exited()
// would keep every verdict here true and still fail these assertions, which is the
// point: the claim of this ticket is about the authority, not about the answer.
//
// Nothing here changes a budget or a threshold: subjectReportBudget, subjectGrace,
// subjectReadyBudget and subjectPollInterval are not written to, and no SLO number
// is asserted. Case 17 hands the loop the production pair BY VALUE (read, not
// moved) so that the sentence it compares is the one an operator would get; cases
// 4-14 of slo_report_144_windows_test.go hand it literal parameters for the same
// reason. Case 18 calls waitReady(), whose budget IS the frozen constant - it
// returns on the first poll once the predicate asks the OS, and it is the shape
// ticket 152's second outlet (slo_windows.go:520, now :522) never produced.
//
// The role children are this test binary re-entering itself by name (same idiom as
// .scratch/wisp/probes/156/zz156probe_windows_test.go), so no cmd.exe or ping.exe
// is depended on and TerminateProcess lands on exactly one pid. They enter through
// case 19 and not through a graded cell - the reason is measured, in case 19.

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/windows"
)

const (
	// slo156EnvRole names the child role; unset in the parent. It is deliberately
	// spelled differently from the 156 PROBE's variable (WISP156_CHILD_ROLE), which
	// lives outside this package's shipped bytes, so the two rulers can be
	// compiled into one test binary and still not read each other's children.
	slo156EnvRole = "WISP_SLO156_CHILD_ROLE"
	// slo156RoleExit dies with slo156ChildExit immediately: the report file is
	// never created, which is ticket 152's "nofile" shape.
	slo156RoleExit = "exit"
	// slo156RoleSleep stays alive until killed, which is the negative control: a
	// predicate that asked the OS wrongly ("always gone") answers case 15 fine and
	// fails this one.
	slo156RoleSleep = "sleep"
	// slo156ChildExit is the code the child dies with; two sentences name it.
	slo156ChildExit = 7
	// slo156StillActive is Win32's answer for a running process, and the reason
	// the witness in this file and the shipped arm judge liveness by different
	// doors (see slo156OsStatus).
	slo156StillActive = 259
	// slo156DeathWait bounds how long a cell waits for its fixture child to
	// really die. It is a fixture guard: exceeding it fails the test rather than
	// grading anything about the code under test.
	slo156DeathWait = 20 * time.Second
	slo156SleepFor  = 30 * time.Second
)

// TestSLO156FixtureChildrenDoWhatTheirNamesSay - case 19, the fixture's own nail,
// and the entry point every cell in this file re-enters as a child.
//
// It exists because of a mistake this file made first and fixed second. Its
// earlier draft put the role dispatch inside case 15 and spawned children with
// `-test.run=^<case 15's name>$`; the mutation ruler then "neutralised case 15" by
// renaming that function - which quietly broke the CHILDREN of cases 17 and 18,
// so two cells went red for a reason that had nothing to do with the code under
// test (measured: probes/156/mut-156/case15-off.log, red=17,18 with the shipped
// bytes untouched). A fixture whose wires are one of the graded cases cannot be
// taken apart to test anything.
//
// So the role lives in its own test, and this test also STATES the two facts every
// cell leans on: the "exit" child really dies with slo156ChildExit and reaches the
// unreaped-with-a-code state, and the "sleep" child really is alive to the OS until
// something kills it. Neutralising it is never a legitimate mutation cell, and the
// ruler says so out loud (my156.py has no cell that renames this test).
//
// The measurement behind this paragraph, not a story about it:
// .scratch/wisp/probes/156/mut-156-first-run-summaries.txt, line `case15-off` -
// rc=1, FAIL=2, red = cases 17 and 18, with slo_windows.go at the SHIPPED bytes.
// (That run's raw per-cell logs are gone - this程 rotated the directory and broke
// "临时件只建不删" doing it; the surviving reading is the ruler's own stdout, and
// the evidence file's §9 carries the account.)
func TestSLO156FixtureChildrenDoWhatTheirNamesSay(t *testing.T) {
	if role := os.Getenv(slo156EnvRole); role != "" {
		slo156PlayRole(role) // never returns
		return
	}

	// 19a - the dying child: a real code, and the state every cell asks about.
	cmd := slo156Spawn(t, slo156RoleExit)
	pid := uint32(cmd.Process.Pid)
	code, how := slo156OsDead(t, pid)
	if code != slo156ChildExit {
		t.Fatalf("the OS reported exit code %d for the %q child, want %d (%s) - every give-up sentence in this file names that number",
			code, slo156RoleExit, slo156ChildExit, how)
	}
	if cmd.ProcessState != nil {
		t.Fatalf("the %q child already has a ProcessState (%v) before anyone waited on it - the witness is not the only thing looking",
			slo156RoleExit, cmd.ProcessState)
	}
	if err := cmd.Wait(); err == nil {
		t.Errorf("Wait on a child that exits %d returned nil error; every cell's exit status assumption is then wrong", slo156ChildExit)
	}
	if cmd.ProcessState == nil || cmd.ProcessState.ExitCode() != slo156ChildExit {
		t.Errorf("after Wait, ProcessState says %v, want exit code %d", cmd.ProcessState, slo156ChildExit)
	}

	// 19b - the living child: alive until killed, and the kill is visible too.
	live := slo156Spawn(t, slo156RoleSleep)
	lpid := uint32(live.Process.Pid)
	if !slo156OsAlive(t, lpid) {
		t.Fatalf("the %q child is not alive to the OS - it cannot serve as the negative control", slo156RoleSleep)
	}
	if err := live.Process.Kill(); err != nil {
		t.Fatalf("kill the %q child: %v", slo156RoleSleep, err)
	}
	if c, khow := slo156OsDead(t, lpid); c == slo156ChildExit {
		t.Errorf("killed child reports exit code %d, the same as the %q role (%s) - the two roles are no longer tellable apart",
			c, slo156RoleExit, khow)
	}
}

// TestSLO156ExitedAsksTheOSForAChildNobodyReaped - case 15, the mechanism itself.
func TestSLO156ExitedAsksTheOSForAChildNobodyReaped(t *testing.T) {
	cmd := slo156Spawn(t, slo156RoleExit)
	s := &sloSubject{cmd: cmd, pid: uint32(cmd.Process.Pid), outPath: "unused-in-this-case.json"}

	code, how := slo156OsDead(t, uint32(cmd.Process.Pid))
	if code != slo156ChildExit {
		t.Fatalf("the OS itself reported exit code %d, fixture wanted %d (%s) - the assertions below name that number", code, slo156ChildExit, how)
	}
	if s.cmd.ProcessState != nil {
		t.Fatalf("fixture drifted: ProcessState is already filled (%v), so this cell would be measuring our own record and not the OS", s.cmd.ProcessState)
	}

	// THE NAIL: the OS says dead, so we say dead - on the first ask, while nobody
	// has reaped anything. Before ticket 156 this is false/-1, because ProcessState
	// is only filled by Wait and the only Wait in this type is stop()'s.
	exited := s.exited()
	got := s.exitCode()
	if !exited {
		t.Errorf("exited() is false for a subject the OS reports terminated with code %d (%s) while ProcessState is still nil: that is the shape Q-57 ruled on, and it is why the loop prints the budget sentence over a corpse", code, how)
	}
	if got != slo156ChildExit {
		t.Errorf("exitCode() is %d, want %d (the code the OS reports): the give-up sentence prints this number", got, slo156ChildExit)
	}
	if s.cmd.ProcessState != nil {
		t.Errorf("exited()/exitCode() filled ProcessState (%v): asking the OS must not turn into reaping the child here - stop() owns the reap, and this ticket does not move it", s.cmd.ProcessState)
	}
	// Asking twice must not consume anything: the loop asks once per poll, for the
	// whole budget.
	if !s.exited() {
		t.Error("the second ask of exited() answered false for a child the OS reports dead")
	}
}

// TestSLO156LiveSubjectStaysAliveUntilTheOSDisagrees - case 16, the other
// direction, and both arms of exitStatus in one timeline.
//
// 16a: a running child is NOT dead. Without this, "exited() := true" would answer
// every other case in this file and keep the loop's sentence short-circuited.
// 16b: a child WE killed, still unreaped, is dead - the same arm as case 15, from
// a different cause of death.
// 16c: after stop()-style reaping, the answer comes from ProcessState, and the
// handle arm is not even consulted (os/exec will not lend it after Wait).
func TestSLO156LiveSubjectStaysAliveUntilTheOSDisagrees(t *testing.T) {
	cmd := slo156Spawn(t, slo156RoleSleep)
	s := &sloSubject{cmd: cmd, pid: uint32(cmd.Process.Pid)}

	if !slo156OsAlive(t, uint32(cmd.Process.Pid)) {
		t.Fatal("fixture child is not alive to the OS before the assertions ran")
	}
	if s.exited() {
		t.Errorf("exited() is true for a subject the OS reports running: give-up sentences would fire on a live child and stop() would skip its Kill")
	}
	if got := s.exitCode(); got != exitCodeUnknown {
		t.Errorf("exitCode() is %d for a live subject, want %d (exitCodeUnknown)", got, exitCodeUnknown)
	}

	if err := cmd.Process.Kill(); err != nil {
		t.Fatalf("kill the fixture child: %v", err)
	}
	code, how := slo156OsDead(t, uint32(cmd.Process.Pid))
	if !s.exited() {
		t.Errorf("exited() is false after the OS reported our own Kill with code %d (%s), ProcessState nil=%v", code, how, s.cmd.ProcessState == nil)
	}

	if err := cmd.Wait(); err != nil {
		t.Logf("Wait on a killed child reports the kill as an error (%v); that is the fixture, not a failure", err)
	}
	if s.cmd.ProcessState == nil {
		t.Fatal("Wait returned without a ProcessState - 16c cannot tell the two arms apart")
	}
	if !s.exited() {
		t.Errorf("exited() is false for a reaped child: the ProcessState arm went to sleep")
	}
}

// TestSLO156ReportLoopNamesTheDeadSubjectItWasWaitingOn - case 17, the shot
// itself at the production seam: real unreaped corpse, the loop's own give-up.
//
// This is case 13's promise with the one hand-tied half untied: the child is dead
// and NOT reaped, which is the only wiring startSubject actually produces. The
// report bytes stay scripted for the reason case 5 and case 13 give - counting
// readings is how "on the spot" is shown without a wall-clock claim - and the
// fixture child writes no file at all, which is 152's `nofile` shape and the one
// where the old sentence named nothing but a spent budget.
//
// Against the pre-ticket口径 the same cell burns the budget it is handed and comes
// back with the neighbour's sentence: mutation m1 (the OS arm removed, nothing else
// changed) measures 1469 real readings in 30.07s here, and the same removal in
// waitReady costs a further 60.06s - both in
// .scratch/wisp/probes/156/mut-156/m1-os-arm-removed.log (an earlier run of the
// same cell read 1468; the count is this host's, the sentence is the point).
func TestSLO156ReportLoopNamesTheDeadSubjectItWasWaitingOn(t *testing.T) {
	cmd := slo156Spawn(t, slo156RoleExit)
	slo156OsDead(t, uint32(cmd.Process.Pid))

	head := []byte(`{"mode":`) // 8 bytes whose tail never arrives
	r := &scriptedReader{t: t, repeatLast: true, got: []func() ([]byte, error){scriptBytes(head)}}
	s := &sloSubject{cmd: cmd, pid: uint32(cmd.Process.Pid), outPath: "scripted.json", readReportFile: r.read}
	_, err := s.collectReportWithin(subjectReportBudget, subjectPollInterval)
	if err == nil {
		t.Fatal("a subject the OS reports dead passed the report loop")
	}
	msg := strings.ReplaceAll(err.Error(), "\n", " ")
	if s.cmd.ProcessState != nil {
		t.Errorf("fixture drifted: something reaped the child before the give-up (ProcessState %v), so this cell is case 13 again", s.cmd.ProcessState)
	}
	if !strings.Contains(msg, "exited (code 7) without writing its report") {
		t.Errorf("give-up %q does not name the dead child and its code", msg)
	}
	if strings.Contains(msg, "never wrote a complete report within") {
		t.Errorf("an unreaped dead subject was reported as a spent budget, not as a corpse - the exact misattribution Q-57 was raised for: %q", msg)
	}
	if !strings.Contains(msg, "8 bytes read, document still open at offset 8") {
		t.Errorf("give-up %q lost the last reading case 13 pinned at this same call site", msg)
	}
	if r.calls != 1 {
		t.Errorf("dead subject was read %d times, want exactly 1 (the loop asks the OS and gives up on the spot)", r.calls)
	}
}

// TestSLO156WaitReadyNamesTheDeadSubjectToo - case 18, the second outlet.
//
// slo_windows.go:520 is the other caller of the same predicate, and ticket 152's
// second格 measured that it could never produce its own sentence either: a subject
// that dies before writing its readiness marker printed "never reported ready
// within 1m0s", naming a budget the operator did not spend. Unlike case 17 this
// one runs on the FROZEN constant budget, so the pre-ticket shape costs 60 real
// seconds before it goes red - the red is the point, the 60s is what the loop
// would have burned on a machine that had nothing left to write.
func TestSLO156WaitReadyNamesTheDeadSubjectToo(t *testing.T) {
	cmd := slo156Spawn(t, slo156RoleExit)
	slo156OsDead(t, uint32(cmd.Process.Pid))

	dir := t.TempDir()
	s := &sloSubject{
		cmd: cmd, pid: uint32(cmd.Process.Pid), dir: dir,
		readyPath: filepath.Join(dir, "ready"),
	}
	err := s.waitReady()
	if err == nil {
		t.Fatal("a subject that died before its readiness marker passed waitReady")
	}
	msg := strings.ReplaceAll(err.Error(), "\n", " ")
	if !strings.Contains(msg, "exited early (code 7) before reporting ready") {
		t.Errorf("waitReady give-up %q does not name the dead subject and its code", msg)
	}
	if strings.Contains(msg, "never reported ready within") {
		t.Errorf("a dead subject was reported as a spent readiness budget: %q", msg)
	}
}

// ---------------------------------------------------------------------------
// helpers. Each one states a FACT about the fixture (a real child, the OS's own
// answer); none of them contains a verdict about the code under test, so no
// assertion below can be satisfied by the thing that set it up.
// ---------------------------------------------------------------------------

// slo156PlayRole is the child half: a real process with a real pid that either
// dies with a named code or stays alive until killed. It never returns.
func slo156PlayRole(role string) {
	switch role {
	case slo156RoleExit:
		os.Exit(slo156ChildExit)
	case slo156RoleSleep:
		time.Sleep(slo156SleepFor)
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "unknown %s=%q\n", slo156EnvRole, role)
		os.Exit(6)
	}
}

// slo156Spawn starts a role child WITHOUT reaping it: cmd.Start only, which is the
// same call startSubject makes through Job.StartInJob, and the reason ProcessState
// stays nil in production. The child re-enters THIS binary through
// TestSLO156FixtureChildrenDoWhatTheirNamesSay (case 19), never through a graded
// cell, for the reason spelled out there. Stdout/Stderr are nil exactly as
// startSubject leaves them, so nothing a child prints can be mistaken for the
// observer's console.
func slo156Spawn(t *testing.T, role string) *exec.Cmd {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestSLO156FixtureChildrenDoWhatTheirNamesSay$",
		"-test.timeout=300s")
	cmd.Env = append(os.Environ(), slo156EnvRole+"="+role)
	cmd.Stdout, cmd.Stderr = nil, nil
	if err := cmd.Start(); err != nil {
		t.Fatalf("spawn %s child: %v", role, err)
	}
	t.Cleanup(func() { _ = cmd.Wait() }) // never leave a child behind
	return cmd
}

// slo156Witness opens the child fresh FROM ITS PID with the right to wait on it
// and to read its code. It is not the handle the implementation reads (os/exec's
// own, handed over by CreateProcess), and opening one without SYNCHRONIZE is not
// enough to wait at all - measured in this file's first draft, a handle opened with
// PROCESS_QUERY_LIMITED_INFORMATION alone answers WaitForSingleObject with
// WAIT_FAILED / "Access is denied". Same process object, two callers of it.
func slo156Witness(pid uint32) (windows.Handle, error) {
	return windows.OpenProcess(windows.SYNCHRONIZE|windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
}

// slo156OsDead blocks until the kernel has FINISHED terminating the child - the
// process object is signaled - and hands back the code it terminated with, so no
// cell can be accused of asking exited() about a child that was still running.
//
// "Finished" is the state the assertions are graded against, and it is a later
// state than "the exit code is already readable": this anchor measured
// GetExitCodeProcess on a pid-opened handle answering the child's code about 6ms
// before that child's process object was signaled
// (.scratch/wisp/probes/156/probe-post/probe-156.log, cell prefix-unreaped - code
// door at 50ms, the shipped predicate agreeing at 56ms). Grading at the signaled
// instant is the sharp claim available ("the OS has finished saying dead and you
// answer dead on the same ask"); the few ms before it are the shipped arm answering
// "not yet", which is the conservative direction and costs one 20ms poll.
//
// The wait is bounded by a fixture guard that FAILS the test rather than grading
// anything: a child that never dies is a broken cell, not a verdict.
func slo156OsDead(t *testing.T, pid uint32) (int, string) {
	t.Helper()
	h, err := slo156Witness(pid)
	if err != nil {
		t.Fatalf("witness OpenProcess(%d): %v", pid, err)
	}
	defer func() { _ = windows.CloseHandle(h) }()
	if waited, werr := windows.WaitForSingleObject(h, uint32(slo156DeathWait/time.Millisecond)); waited != windows.WAIT_OBJECT_0 {
		t.Fatalf("the OS never finished terminating child %d within %s (wait=%d err=%v) - fixture, not a verdict",
			pid, slo156DeathWait, waited, werr)
	}
	var c uint32
	if err := windows.GetExitCodeProcess(h, &c); err != nil {
		t.Fatalf("witness GetExitCodeProcess(%d): %v", pid, err)
	}
	if c == slo156StillActive {
		t.Fatalf("child %d's process object is signaled while GetExitCodeProcess still says STILL_ACTIVE", pid)
	}
	return int(c), "process object signaled, exit code read back"
}

// slo156OsAlive asks the same object whether it is still running, with a
// zero-timeout wait: WAIT_TIMEOUT is the only answer that means running, so a
// WAIT_FAILED here reads as "not alive" and case 16's own guard fails the cell
// rather than letting a broken fixture pass as a verdict.
func slo156OsAlive(t *testing.T, pid uint32) bool {
	t.Helper()
	h, err := slo156Witness(pid)
	if err != nil {
		t.Fatalf("witness OpenProcess(%d): %v", pid, err)
	}
	defer func() { _ = windows.CloseHandle(h) }()
	waited, _ := windows.WaitForSingleObject(h, 0)
	// The conversion is not decoration: x/sys/types declares WAIT_OBJECT_0 as an
	// untyped 0 and WAIT_TIMEOUT as a typed syscall.Errno, which is the same
	// asymmetry slo_windows.go's queryProcess side-steps by only ever comparing
	// against WAIT_OBJECT_0.
	return waited == uint32(windows.WAIT_TIMEOUT)
}
