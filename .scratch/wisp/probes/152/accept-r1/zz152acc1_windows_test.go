// Probe file for the R1 acceptance run of ticket 152. NOT part of the repo:
// it is mapped into cmd/wisp by -overlay from outside the working tree, so
// `git status` stays clean and no shipped byte is touched.
//
// Question it answers with a reading, not a claim: with a REAL child process
// that writes a real prefix of a subject report and then really exits 7, does
// the `if s.exited()` leg of collectReportWithin ever fire on the production
// wiring? Which sentence does the operator get? And what changes the answer?
package main

import (
	"errors"
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
	accEnv      = "WISP152ACC_CHILD"
	accEnvOut   = "WISP152ACC_OUT"
	accExitCode = 7
)

// TestAcc152Helper is the child role. It is only ever selected by the helper
// command line; a normal `-run` selection that excludes it never sees it.
func TestAcc152Helper(t *testing.T) {
	if os.Getenv(accEnv) != "1" {
		t.Fatal("helper role: reached without " + accEnv + "=1")
	}
	out := os.Getenv(accEnvOut)
	body := `{"mode":"subject-in-tree","report":[` + strings.Repeat("0,", 400)
	if err := os.WriteFile(out, []byte(body), 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "ACC-CHILD|write-failed|%v\n", err)
	}
	fmt.Fprintf(os.Stderr, "ACC-CHILD|wrote|%d|then exit %d\n", len(body), accExitCode)
	os.Exit(accExitCode)
}

// acc152OsDead is the external witness: what the OS itself says about this pid,
// asked without reaping it (no Wait, no ProcessState involvement).
func acc152OsDead(t *testing.T, pid uint32) (bool, uint32) {
	t.Helper()
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		// ERROR_INVALID_PARAMETER (87) / no such process means the pid is gone.
		t.Logf("ACC|os|OpenProcess pid=%d err=%v", pid, err)
		return false, 0xFFFFFFFF
	}
	defer windows.CloseHandle(h)
	var code uint32
	if err := windows.GetExitCodeProcess(h, &code); err != nil {
		t.Fatalf("ACC|os|GetExitCodeProcess: %v", err)
	}
	const stillActive = 259 // STILL_ACTIVE
	return code != stillActive, code
}

func acc152Start(t *testing.T) (*exec.Cmd, string, string, int) {
	t.Helper()
	dir := t.TempDir()
	out := filepath.Join(dir, "subject-report.json")
	cmd := exec.Command(os.Args[0], "-test.run=^TestAcc152Helper$", "-test.timeout=60s")
	cmd.Env = append(os.Environ(), accEnv+"=1", accEnvOut+"="+out)
	if err := cmd.Start(); err != nil {
		t.Fatalf("ACC|start: %v", err)
	}
	pid := uint32(cmd.Process.Pid)
	// Wait for the real bytes to land on the real disk, then for the OS to say
	// the pid is dead. Neither of these touches cmd.Wait().
	// os.WriteFile CREATES the file and then fills it, so a size taken the
	// instant the name exists reads 0 (measured: the first cut of this probe
	// asserted against a want of 0 and went red on a correct sentence). The
	// byte count is therefore sampled after the OS reports the pid terminated.
	deadline := time.Now().Add(20 * time.Second)
	dead := false
	var code uint32
	for time.Now().Before(deadline) {
		dead, code = acc152OsDead(t, pid)
		if dead {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	if !dead {
		t.Fatalf("ACC|OS still reports the pid alive after 20s")
	}
	want := -1
	for time.Now().Before(deadline.Add(20 * time.Second)) {
		st, serr := os.Stat(out)
		if serr == nil && st.Size() > 0 {
			want = int(st.Size())
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	if want < 0 {
		t.Fatalf("ACC|the child never wrote a single byte")
	}
	t.Logf("ACC|setup|pid=%d file_bytes=%d os_says=terminated exit_code=%d cmd.ProcessState==nil?%v",
		pid, want, code, cmd.ProcessState == nil)
	return cmd, out, dir, want
}

// Cell A - the production wiring: nobody reaps. Which leg speaks?
func TestAcc152ProbeUnreapedTakesTheTimeoutLeg(t *testing.T) {
	cmd, out, dir, want := acc152Start(t)
	s := &sloSubject{cmd: cmd, pid: uint32(cmd.Process.Pid), dir: dir, outPath: out}
	// readReportFile stays nil: the loop calls the real os.ReadFile.
	exitedBefore := s.exited()
	codeBefore := s.exitCode()
	start := time.Now()
	_, err := s.collectReportWithin(2*time.Second, 5*time.Millisecond)
	spent := time.Since(start)
	if err == nil {
		t.Fatalf("ACC|A|a dead subject with an unwritten report passed the loop")
	}
	sent := err.Error()
	t.Logf("ACC|A|sawReported=false|exited_before=%v exitCode_before=%d elapsed=%s", exitedBefore, codeBefore, spent)
	t.Logf("ACC|A|sentence=%s", sent)
	t.Logf("ACC|A|after_loop|exited=%v exitCode=%d", s.exited(), s.exitCode())

	if exitedBefore {
		t.Errorf("ACC|A|s.exited() was TRUE before the loop although nobody called Wait(): the finding is refuted")
	}
	if strings.Contains(sent, "exited (code") {
		t.Errorf("ACC|A|the exited leg fired on the production wiring - the ticket's root cause is wrong: %q", sent)
	}
	if !strings.Contains(sent, "never wrote a complete report within") {
		t.Errorf("ACC|A|expected the budget leg to speak, got %q", sent)
	}
	if !strings.Contains(sent, fmt.Sprintf("%d bytes read", want)) {
		t.Errorf("ACC|A|the sentence does not name the bytes it really read (%d): %q", want, sent)
	}
	if spent < 2*time.Second {
		t.Errorf("ACC|A|the loop gave up in %s although the child was already dead: budget claim wrong", spent)
	}
	// Does the classifier really see a prefix and not an empty envelope?
	data, rerr := os.ReadFile(out)
	if rerr != nil {
		t.Fatalf("ACC|A|reread: %v", rerr)
	}
	obs := readSubjectReport(data)
	t.Logf("ACC|A|last_reading|state=%s offset=%d bytes=%d summary=%q", obs.state, obs.offset, obs.bytes, obs.summary())
	if obs.state != reportUnwritten {
		t.Errorf("ACC|A|a real prefix classified as %s, not unwritten", obs.state)
	}
	if obs.summary() == "" {
		t.Errorf("ACC|A|the reading carried into the give-up sentence is an empty string")
	}
	_ = errors.Is
}

// Cell B - the same child, one extra line: someone reaps. This is what case 13
// does and what production does not.
func TestAcc152ProbeReapedTakesTheExitedLeg(t *testing.T) {
	cmd, out, dir, want := acc152Start(t)
	if err := cmd.Wait(); err != nil {
		t.Logf("ACC|B|Wait returned %v (expected: exit status 7)", err)
	}
	s := &sloSubject{cmd: cmd, pid: uint32(cmd.Process.Pid), dir: dir, outPath: out}
	t.Logf("ACC|B|after_reap|exited=%v exitCode=%d", s.exited(), s.exitCode())
	if !s.exited() {
		t.Fatalf("ACC|B|even after Wait() exited() is false: the mechanism reading is wrong")
	}
	start := time.Now()
	_, err := s.collectReportWithin(2*time.Second, 5*time.Millisecond)
	spent := time.Since(start)
	if err == nil {
		t.Fatalf("ACC|B|a reaped dead subject passed the loop")
	}
	sent := err.Error()
	t.Logf("ACC|B|elapsed=%s sentence=%s", spent, sent)
	if !strings.Contains(sent, "exited (code 7) without writing its report") {
		t.Errorf("ACC|B|reaping did not produce the exited sentence: %q", sent)
	}
	if strings.Contains(sent, "never wrote a complete report within") {
		t.Errorf("ACC|B|the budget leg spoke where the exited leg should: %q", sent)
	}
	if !strings.Contains(sent, fmt.Sprintf("%d bytes read", want)) {
		t.Errorf("ACC|B|the exited sentence does not carry the last real reading (%d bytes): %q", want, sent)
	}
	if spent >= 2*time.Second {
		t.Errorf("ACC|B|gave up on the spot expected, took %s", spent)
	}
}

// Cell C - the claim as stated in the orchestrator's ledger: "`last` is nil /
// an empty envelope, so the exited leg prints nothing". subjectReportRead is a
// VALUE type here, and the loop assigns `last = obs` before either give-up
// check. This cell pins what a zero value would even render as, so the wording
// can be checked instead of repeated.
func TestAcc152ProbeLastIsAValueNotANilEnvelope(t *testing.T) {
	var zero subjectReportRead
	t.Logf("ACC|C|zero_value_summary=%q state=%q offset=%d", zero.summary(), zero.state, zero.offset)
	if zero.summary() == "" {
		t.Errorf("ACC|C|a zero subjectReportRead renders an empty string")
	}
	// The loop cannot reach either give-up check without having assigned `last`:
	// `last = obs` sits above both of them in the same iteration.
	src, err := os.ReadFile(filepath.Join("..", "..", "cmd", "wisp", "slo_windows.go"))
	if err != nil {
		t.Logf("ACC|C|source read skipped: %v", err)
		return
	}
	body := string(src)
	// Scope the ordering check to collectReportWithin: `if s.exited()` also
	// appears in waitReady further up the file, and my first cut of this cell
	// compared the two different call sites (recorded as a self-refutation).
	start := strings.Index(body, "func (s *sloSubject) collectReportWithin")
	if start < 0 {
		t.Fatalf("ACC|C|cannot find collectReportWithin in the shipped bytes")
	}
	fn := body[start:]
	i := strings.Index(fn, "last = obs")
	e := strings.Index(fn, "if s.exited()")
	if i < 0 || e < 0 || i > e {
		t.Errorf("ACC|C|`last = obs` no longer precedes the exited check inside collectReportWithin (last=%d exited=%d)", i, e)
	}
	t.Logf("ACC|C|inside_collectReportWithin|last_assign_at=%d exited_check_at=%d", i, e)
}
