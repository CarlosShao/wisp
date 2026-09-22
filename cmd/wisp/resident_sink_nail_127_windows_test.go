//go:build windows

package main

// Ticket 127 AC#1: the resident GUI leg's log listener gets a nail of its own.
//
// The reading this file exists to answer is `acceptor-ticket117` M4: deleting
// the whole install block of cmd/wisp/resident_windows.go (the 20 lines from
// the "Ticket 117" comment through `defer sink.close()`) left `go build` rc=0,
// `go vet` rc=0 and `go test ./cmd/wisp/` green at 146/146, while the same
// deletion on the `run` leg (M1) took three cases red. Nothing in this package
// proved that the leg owner actually double-clicks has a listener at all -
// the same family booked seven times now: capability finished, nobody listening.
//
// Why these cases are subprocess cases while the run leg's (logsink_windows_test.go)
// are in-process ones, and why that is the SAME judgment and not a new one:
// the run leg's production entry point is a function that returns
// (`runTextTask`, exit code in hand), so a test may call it and then read the
// disk. The resident leg's entry point is `runResident()` (main.go, no-args
// branch), which never returns until the process is told to stop, and which
// writes to os.Stdout/os.Stderr and may os.Exit(1)/os.Exit(2). Calling it in
// process would mean a test reimplementing its surroundings, which is exactly
// the shape R-117-A counted as zero. So the driver here is the shipped
// wisp.exe started with NO arguments - the same command line a double-click
// produces - and the evidence is what that process leaves on disk.
//
// Three properties these cases are built to keep, from the ticket:
//
//   - No window. Ticket 07 has not landed, so this leg boots, prints and parks
//     in the empty event loop (MainWindowHandle = 0). Nothing here opens,
//     waits for or reads a window; that is why the leg can be nailed at all.
//   - No t.Skip, no relaxed assertion. If the install is gone the cases go red
//     by name, and the reason they give is that the shipped resident process
//     wrote nothing where its own listener says it writes.
//   - Under WISP_ENV=test: the test env registers no single-instance mutex
//     (internal/proc/envfork.go), so this child can boot next to an owner's
//     running Wisp without activating it, and WISP_TEST_DATA_DIR keeps every
//     byte it writes inside the harness' own directory.
//
// What is deliberately NOT claimed here: this leg has no sealing site today
// (`proc.Boot` imports no winsec), so the `winsec: seal cleared principals...`
// WARN cannot be produced by it - see ticket 117 R-117-A and the ticket's own
// §二. The record nailed below is the install record and the D38(e) shutdown
// trail, which are the two things this leg really does log.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
	"golang.org/x/sys/windows"
)

// The three sentences the resident leg's 20 lines are responsible for, copied
// as literals on purpose (the same rule as sealNoticeMsg): if any of them is
// reworded, a case in this file goes red and somebody has to say out loud where
// the record that proves the listener was there now lives.
const (
	// residentInstallMsg is installLogSink's own booking record (logsink.go),
	// the one line that says "this process has a persistent listener from this
	// moment, into this directory".
	residentInstallMsg = "wisp: persistent log sink installed"
	// residentReachedLoop is printed by runResident only after the install
	// attempt and the shutdown defer are both behind it.
	residentReachedLoop = "empty event loop running"
	// residentRefusalPrefix is the loud fallback that must appear on the
	// console when the sink cannot be installed, and residentRefusalPromise is
	// the half of that sentence telling the operator what is lost.
	residentRefusalPrefix  = "wisp: 持久日志未启用"
	residentRefusalPromise = "不会落盘"
	// residentShutdownLine is runResident's own exit report, printed by the
	// defer that the sink-close defer is registered behind (D38(e) first).
	residentShutdownLine = "exited through the D38(e) shutdown order"
	// residentShutdownRecord is what internal/proc logs per step that has
	// nothing attached; on this leg most steps are of that kind, and their
	// landing in the same file is the ordering proof.
	residentShutdownRecord = "shutdown step"
)

// sinkInstallRecord is one JSONL record with the fields this ticket's claims are
// about. sinkRecord (logsink_windows_test.go) carries the winsec notice's fields;
// the install record's own attributes (which data root, which level) are what the
// resident leg's claims need, so they are decoded here and nowhere else.
type sinkInstallRecord struct {
	Time     string `json:"time"`
	Level    string `json:"level"`
	Msg      string `json:"msg"`
	Dir      string `json:"dir"`
	MinLevel string `json:"min_level"`
	Step     int    `json:"step"`
	Name     string `json:"name"`
}

func (r sinkInstallRecord) String() string {
	return fmt.Sprintf("{time=%s level=%s msg=%q dir=%s step=%d}", r.Time, r.Level, r.Msg, r.Dir, r.Step)
}

// lockedBuf collects a child's output while it is still alive, so a case can
// wait for a printed line instead of sleeping for a guessed duration.
type lockedBuf struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (l *lockedBuf) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.Write(p)
}

func (l *lockedBuf) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.String()
}

func (l *lockedBuf) has(sub string) bool {
	return strings.Contains(l.String(), sub)
}

// residentLeg is one shipped wisp.exe standing in the no-args branch, with its
// console captured and its process hand still available.
type residentLeg struct {
	dataDir string
	cmd     *exec.Cmd
	stdout  *lockedBuf
	stderr  *lockedBuf
	// done receives the result of cmd.Wait exactly once, from the reaper
	// goroutine started at boot. cmd.Wait also flushes the two copy goroutines,
	// so no case may read stdout/stderr before it has reported.
	done chan error
	// waitErr is what the child ended with; waited says whether it ended.
	waitErr error
	waited  bool
}

// bootResidentLeg starts the real binary with no arguments. The command line is
// the whole point: anything other than "no args" would be a different leg.
//
// CREATE_NEW_PROCESS_GROUP keeps the console break used by the ordering case
// addressed to this child alone; it does not disable CTRL_BREAK (only CTRL_C),
// which is why that case can ask for a clean exit at all.
func bootResidentLeg(t *testing.T, exe, dataDir string) *residentLeg {
	t.Helper()
	cmd := exec.Command(exe)
	cmd.Dir = filepath.Dir(exe)
	cmd.Env = append(os.Environ(), "WISP_ENV=test", procTestDataDirEnv+"="+dataDir)
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: windows.CREATE_NEW_PROCESS_GROUP}
	leg := &residentLeg{
		dataDir: dataDir,
		cmd:     cmd,
		stdout:  &lockedBuf{},
		stderr:  &lockedBuf{},
	}
	cmd.Stdout, cmd.Stderr, cmd.Stdin = leg.stdout, leg.stderr, nil
	if err := cmd.Start(); err != nil {
		t.Fatalf("start the resident leg (%s, no args): %v", exe, err)
	}
	leg.done = make(chan error, 1)
	go func() { leg.done <- cmd.Wait() }()
	// Whatever a case does afterwards, no child may outlive it: an orphaned
	// wisp.exe would hold its temp dir (and, in the dev/prod envs, the session
	// mutex) open past the test.
	t.Cleanup(leg.stop)
	return leg
}

const procTestDataDirEnv = "WISP_TEST_DATA_DIR"

// pid of the running child (its own process-group id, which is what a
// CREATE_NEW_PROCESS_GROUP child is addressed by).
func (l *residentLeg) pid() uint32 { return uint32(l.cmd.Process.Pid) }

// breakToLoop asks the child to leave the event loop the way Ctrl+C does: a
// CTRL_BREAK_EVENT addressed to its own process group. Go's console handler
// turns that into os.Interrupt, proc.RunEventLoop returns "signal", and the
// deferred chain (D38(e) sequence, then sink.close) runs.
func (l *residentLeg) breakToLoop() error {
	return windows.GenerateConsoleCtrlEvent(windows.CTRL_BREAK_EVENT, l.pid())
}

// exitedWithin gives the child a bounded number of 50ms ticks to finish leaving
// through its own shutdown path (no wall-clock deadline arithmetic, D42#9).
func (l *residentLeg) exitedWithin(attempts int) bool {
	for i := 0; i < attempts; i++ {
		if l.reap(50 * time.Millisecond) {
			return true
		}
	}
	return l.reap(50 * time.Millisecond)
}

// reap takes the Wait result if it has arrived within the tick budget.
func (l *residentLeg) reap(tick time.Duration) bool {
	if l.waited {
		return true
	}
	select {
	case err := <-l.done:
		l.waited = true
		l.waitErr = err
		return true
	case <-time.After(tick):
		return false
	}
}

// stop kills and reaps the child if a case has not let it exit by itself.
func (l *residentLeg) stop() {
	if l.reap(0) {
		return
	}
	_ = l.cmd.Process.Kill()
	// Bounded: a child wedged inside its own shutdown is a finding the caller
	// hears about through the case that asked for a clean exit, not a hang.
	if !l.exitedWithin(200) {
		return
	}
	l.waited = true
}

// exitStatus reports how the child ended, for the failure messages.
func (l *residentLeg) exitStatus() string {
	if !l.waited {
		return "still running"
	}
	if l.waitErr != nil {
		var ee *exec.ExitError
		if errors.As(l.waitErr, &ee) {
			return fmt.Sprintf("exit code %d (%v)", ee.ExitCode(), l.waitErr)
		}
		return l.waitErr.Error()
	}
	return "exited 0"
}

func (l *residentLeg) console() string {
	return "stdout:\n" + l.stdout.String() + "\nstderr:\n" + l.stderr.String()
}

// pollUntil127 runs done on a fixed attempt budget (the shape ticket 63's
// pollCommandLine uses) and reports whether it ever said yes.
func pollUntil127(attempts int, done func() bool) bool {
	for i := 0; i < attempts; i++ {
		if done() {
			return true
		}
		time.Sleep(50 * time.Millisecond)
	}
	return done()
}

// readResidentSink returns the records the resident process wrote, refusing a
// directory that holds no pipeline file. The file naming rule comes from
// observe itself (CountLogFiles), so a renamed pipeline is a red case here
// rather than a silently empty read.
func readResidentSink(t *testing.T, dir string) []sinkInstallRecord {
	t.Helper()
	n, err := observe.CountLogFiles(dir)
	if err != nil {
		t.Fatalf("reading the resident leg's sink dir %s: %v", dir, err)
	}
	if n == 0 {
		t.Fatalf("no wisp-<day>-<seq>.jsonl in %s: the resident process wrote nothing there", dir)
	}
	if n != 1 {
		t.Errorf("pipeline files in %s = %d, want 1 (one resident boot, one file)", dir, n)
	}
	var out []sinkInstallRecord
	for _, name := range mustGlob117(t, filepath.Join(dir, "wisp-*.jsonl")) {
		b, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		for _, line := range strings.Split(string(b), "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			var rec sinkInstallRecord
			if err := json.Unmarshal([]byte(line), &rec); err != nil {
				t.Fatalf("line %q in %s is not a JSON object: %v", line, name, err)
			}
			out = append(out, rec)
		}
	}
	t.Logf("resident sink %s: %d file(s), %d record(s)", dir, n, len(out))
	return out
}

// jsonlFilesUnder lists every *.jsonl in a tree, so a case can say "this byte
// landed nowhere else" and mean it.
func jsonlFilesUnder(t *testing.T, root string) []string {
	t.Helper()
	var found []string
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(p, ".jsonl") {
			found = append(found, p)
		}
		return nil
	})
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("walking %s: %v", root, err)
	}
	return found
}

// sinkHasARecord reports whether the pipeline directory holds a file with at
// least one record in it. The rolling writer opens its file as soon as the
// pipeline starts, so "the file exists" is not yet "the listener booked
// itself" - and a case that read on the file's existence alone would race its
// own 500ms flush tick.
func sinkHasARecord(dir string) bool {
	names, err := filepath.Glob(filepath.Join(dir, "wisp-*.jsonl"))
	if err != nil {
		return false
	}
	for _, n := range names {
		b, err := os.ReadFile(n)
		if err != nil {
			return false
		}
		for _, line := range strings.Split(string(b), "\n") {
			if strings.TrimSpace(line) != "" {
				return true
			}
		}
	}
	return false
}

// TestAC1ResidentLegInstallsItsLogListenerOnDisk is the closing nail: it drives
// the shipped binary the way owner launches it (no arguments, no terminal,
// stderr going nowhere) and requires the listener's own record to be on disk
// under the data root, written by that process and by no other writer.
//
// Delete the install block in resident_windows.go and this case is the one that
// goes red, with a name that says which leg lost its listener.
func TestAC1ResidentLegInstallsItsLogListenerOnDisk(t *testing.T) {
	exe := buildWispForTest(t)
	dataDir := t.TempDir()
	leg := bootResidentLeg(t, exe, dataDir)
	sinkDir := logSinkDir(dataDir)

	// The pipeline's flusher runs every 500ms (internal/observe/logging.go), so
	// a 10s budget is ~20 flush ticks: enough for a slow disk, not enough to
	// hide a missing install.
	if !pollUntil127(200, func() bool { return sinkHasARecord(sinkDir) }) {
		files := jsonlFilesUnder(t, dataDir)
		leg.stop()
		t.Fatalf("the resident process wrote no log file under %s within the budget.\n"+
			"data root tree holds %d *.jsonl: %v\n%s\n"+
			"that is the claim this case exists for: the leg owner actually runs has no listener on disk.",
			sinkDir, len(files), files, leg.console())
	}

	recs := readResidentSink(t, sinkDir)
	if len(recs) == 0 {
		t.Fatalf("%s holds a pipeline file with no record in it", sinkDir)
	}
	first := recs[0]
	if first.Msg != residentInstallMsg {
		t.Errorf("record 0 = %q, want %q: the install record must be the first thing in the file, "+
			"or the listener missed its own booking (records: %v)", first.Msg, residentInstallMsg, recs)
	}
	if first.Level != "INFO" {
		t.Errorf("install record level = %q, want INFO", first.Level)
	}
	// WISP_TEST_DATA_DIR is an identity contract (internal/proc/envfork.go): it
	// comes back verbatim, so this comparison is exact and not a tree match.
	if first.Dir != sinkDir {
		t.Errorf("install record names dir = %q, want %q: the resident listener must book the data root it was booted with",
			first.Dir, sinkDir)
	}
	if first.MinLevel != "info" {
		t.Errorf("install record min_level = %q, want the schema default %q", first.MinLevel, "info")
	}

	leg.stop()
	// The console half of the fan-out, read from the child's own captured
	// stderr: this is what ties the FILE to THIS process, so the pass cannot be
	// bought by another writer that happened to use the same directory.
	if !leg.stderr.has(residentInstallMsg) {
		t.Errorf("the install record never reached the child's console; the fan-out went one way only.\n%s", leg.console())
	}
	if !leg.stdout.has(residentReachedLoop) {
		t.Errorf("the child never reached the event loop, so the disk read above proves little.\n%s", leg.console())
	}
	// Placement on this leg, same rule as the run leg's AC#3: the bytes stay
	// inside <data root>\logs.
	for _, p := range jsonlFilesUnder(t, dataDir) {
		if rel, err := filepath.Rel(sinkDir, p); err != nil || strings.HasPrefix(rel, "..") {
			t.Errorf("the resident listener wrote %q outside %s", p, sinkDir)
		}
	}
}

// TestAC1ResidentLegOutlivesItsOwnLogFailure is the fallback branch of the same
// 20 lines: a data root whose log directory cannot be opened must produce a
// named, loud refusal AND still boot, because "the log directory is broken" is
// not a reason to keep Wisp from starting. The notices then fall back to the
// console, which is the state before ticket 117 - and this case is what keeps
// that sentence from being a comment.
//
// The shape is real rather than hypothetical: the sink's first act is
// os.MkdirAll(<data root>\logs), so anything already standing there that is not
// a directory (an operator's plain file, a leftover symlink target) is exactly
// the condition that makes the install fail.
func TestAC1ResidentLegOutlivesItsOwnLogFailure(t *testing.T) {
	exe := buildWispForTest(t)
	dataDir := t.TempDir()
	blocker := logSinkDir(dataDir)
	if err := os.WriteFile(blocker, []byte("not a directory"), 0o600); err != nil {
		t.Fatal(err)
	}
	leg := bootResidentLeg(t, exe, dataDir)

	sawRefusal := pollUntil127(200, func() bool { return leg.stderr.has(residentRefusalPrefix) })
	sawLoop := pollUntil127(20, func() bool { return leg.stdout.has(residentReachedLoop) })
	leg.stop()

	if !sawRefusal {
		t.Fatalf("a log directory that cannot be opened produced no refusal on the console.\n%s\n"+
			"the operator has to be told the security notices are not being kept; that is the other half of the 20 lines.",
			leg.console())
	}
	if !sawLoop {
		t.Fatalf("the refusal stopped the app: the child never reached the event loop.\n%s", leg.console())
	}
	if !leg.stderr.has(residentRefusalPromise) {
		t.Errorf("the refusal does not say what is lost (%q); captured stderr:\n%s", residentRefusalPromise, leg.stderr.String())
	}
	// A refusal that silently kept a half-open sink would be worse than no
	// sink: nothing may have been written.
	if files := jsonlFilesUnder(t, dataDir); len(files) != 0 {
		t.Errorf("the refused install produced %d jsonl file(s) anyway: %v", len(files), files)
	}
	if st, err := os.Stat(blocker); err != nil || st.IsDir() {
		t.Errorf("the blocker at %s is gone or turned into a directory (stat err %v)", blocker, err)
	}
}

// TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink is the ordering half of
// these 20 lines, and it is the reading ticket 117 registered as missing
// (§一 row 7 and §四: "GUI 腿；但这半格的读数没采到"). The two defers in
// resident_windows.go are registered in a specific order so that the D38(e)
// sequence runs while the sink is still open and the close flushes its records.
//
// So the child is asked to leave the loop the way an operator does (Ctrl+C,
// here a CTRL_BREAK_EVENT addressed to its own process group, which Go's console
// handler raises as os.Interrupt), and the file then has to carry both halves:
// the install record first, the shutdown trail last.
func TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink(t *testing.T) {
	exe := buildWispForTest(t)
	dataDir := t.TempDir()
	sinkDir := logSinkDir(dataDir)
	leg := bootResidentLeg(t, exe, dataDir)
	t.Cleanup(func() {
		t.Logf("resident leg exit: %s", leg.exitStatus())
	})

	if !pollUntil127(200, func() bool {
		n, err := observe.CountLogFiles(sinkDir)
		return err == nil && n > 0
	}) {
		leg.stop()
		t.Fatalf("no sink file appeared, so the ordering has nothing to be read out of.\n%s", leg.console())
	}
	if err := leg.breakToLoop(); err != nil {
		leg.stop()
		t.Fatalf("GenerateConsoleCtrlEvent(CTRL_BREAK_EVENT, %d) on the leg's own process group: %v\n%s",
			leg.pid(), err, leg.console())
	}
	if !leg.exitedWithin(400) {
		leg.stop()
		t.Fatalf("the child took the break and did not exit through its own shutdown path.\n%s", leg.console())
	}
	if leg.waitErr != nil {
		t.Errorf("the child exited through the shutdown path with %v, want 0\n%s", leg.waitErr, leg.console())
	}

	recs := readResidentSink(t, sinkDir)
	if recs[0].Msg != residentInstallMsg {
		t.Fatalf("record 0 = %q, want the install record: the trail below cannot be read as \"the listener caught the shutdown\" unless it starts with the listener (records: %v)", recs[0].Msg, recs)
	}
	var trail []sinkInstallRecord
	for _, r := range recs[1:] {
		if strings.Contains(r.Msg, residentShutdownRecord) {
			trail = append(trail, r)
		}
	}
	if len(trail) == 0 {
		t.Fatalf("the D38(e) trail is not in the file the same process closed: %d record(s) after the install, all msgs %v",
			len(recs)-1, msgsOf127(recs))
	}
	last := recs[len(recs)-1]
	if !strings.Contains(last.Msg, residentShutdownRecord) {
		t.Errorf("the last record in the sink is %q, want a shutdown-step record: the shutdown must be logged before the sink closes (records: %v)", last.Msg, recs)
	}
	if !leg.stdout.has(residentShutdownLine) {
		t.Errorf("the child printed no D38(e) exit report; captured:\n%s", leg.console())
	}
	t.Logf("RESIDENT LEG RECORDED ON DISK: %d shutdown record(s), steps %v, first=%q last=%q",
		len(trail), stepsOf127(trail), trail[0].Msg, trail[len(trail)-1].Msg)
}

func msgsOf127(recs []sinkInstallRecord) []string {
	out := make([]string, 0, len(recs))
	for _, r := range recs {
		out = append(out, r.Msg)
	}
	return out
}

func stepsOf127(recs []sinkInstallRecord) []int {
	out := make([]int, 0, len(recs))
	for _, r := range recs {
		out = append(out, r.Step)
	}
	return out
}
