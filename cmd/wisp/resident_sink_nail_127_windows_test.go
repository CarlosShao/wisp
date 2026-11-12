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
// §二. The records nailed below are the install record and the D38(e) shutdown
// trail, which are the two things this leg really does log, plus the one record
// ticket 130 replays into this file from before the listener existed.
//
// WHY RECORD 0 CHANGED, and what it used to protect (ticket 130, A110 row 3
// authorization; the intent is kept and the claim gets stronger, not weaker):
//
// These cases used to assert `recs[0].Msg == "wisp: persistent log sink
// installed"`. That sentence was doing two jobs. One was the real claim: the
// listener's own booking is the landmark in the file, and nothing the process
// meant to audit precedes the thing that audits it. The other was incidental to
// the pre-ticket-130 world: since a record made during a package init() could
// not reach the file at all, the install booking was necessarily the first line
// - so "index 0" silently also encoded "no boot-time record exists", which is
// the very defect ticket 130 was opened for. Yesterday's assertion passed BECAUSE
// the defect was present.
//
// The claim after the fix keeps the landmark job and adds to it: record 0 must
// be the replayed early record, the install booking must come next, every record
// before the install booking must be one of the replayed ones, and its timestamp
// must be older than the install booking's. Delete the replay in
// cmd/wisp/logsink.go and these cases go red again, in the same direction as
// before - the early record is simply not in the file - so the authorization to
// touch this nail did not buy a weaker instrument.

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
	//
	// RE-ANCHORED by ticket 246 AC#7 (leg 246-r2), and the reason is this file's
	// own rule rather than a convenience: the sentence it pointed at said "empty
	// event loop running (no task is taken by this loop yet)", and AC#7 is the
	// increment that makes that clause false - the resident leg now has a task
	// entry, so the landmark had to move onto the part of the sentence that is
	// still doing the same job. What the claim protects did not move and is not
	// weaker:
	//
	//   - duty 1 (the sink landmark): this line still only arrives from the branch
	//     that ran AFTER installLogSink and the shutdown defer, so a child whose
	//     console never reaches it is still read as "the disk record proves
	//     little";
	//   - duty 2 (no earlier start record): unchanged, and it is the record index
	//     assertions above that carry it, not this string;
	//   - delete the replay in cmd/wisp/logsink.go and these cases still go red on
	//     record 0, exactly as before.
	//
	// The new token is true in BOTH branches of the print (console entry wired or
	// not) and false in the during-boot-exit branch, which is the same shape the
	// old one had. The task-source posture now rides on the same line, and a case
	// in resident_task_source_246_windows_test.go reads that half.
	residentReachedLoop = "resident event loop running"
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
	// residentReplayRecord is the prefix of the ten-step replay ticket 43 added
	// to resident_windows.go's shutdown defer, AFTER the sequence itself has run.
	// Ticket 425 re-anchored this file's ordering claim onto it: the trail is
	// read as "booked before the replay", which is the promise the two defers
	// make, instead of "last record in the file", which the replay made
	// unsatisfiable. Copied as a literal on the same rule as the two above.
	residentReplayRecord = "wisp: D38(e) step"
	// residentEarlyResolverMsg is internal/winsec's sealing-seam verdict, emitted
	// from internal/risk's package init() - before installLogSink can exist - and
	// replayed into this leg's file by ticket 130's buffer. Copied as a literal on
	// the same rule as the two above: if the sentence moves, this file has to say
	// where the boot-time record of the sealing seam went instead.
	residentEarlyResolverMsg = "winsec: sealing path resolver installed"
)

// installRecordIndex127 finds the listener's own booking record, the landmark
// every ordering claim in this file is read against. -1 means the resident leg
// never booked its sink at all, which is the ticket 117 failure this file was
// written for.
func installRecordIndex127(recs []sinkInstallRecord) int {
	for i, r := range recs {
		if r.Msg == residentInstallMsg {
			return i
		}
	}
	return -1
}

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
	Resolver string `json:"resolver"`
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

// bootFactCeiling127 is the CEILING on a wait for a boot fact. It is not the
// wait and it is not the measurement: a wait built on it ends the instant the
// fact becomes observable, and the number below only decides how long a boot
// that never produces the fact takes to report red instead of hanging the
// package.
//
// WHY THIS SHAPE AND NOT A BIGGER ATTEMPT COUNT (ticket 431). The two cases in
// this file used to hand `pollUntil127` a fixed attempt count, which is a
// duration written as a multiplication: 200 x 50ms = 10s, 20 x 50ms = 1s. Both
// numbers are guesses about how long this machine takes to boot, and both were
// wrong in a way no amount of tuning fixes, because the quantity being waited
// for is not a duration - it is a fact ("this process reached its event loop"),
// and a duration either outlives the fact (and then hides a real hang behind a
// long silence) or does not (and then reports red on a slow machine for the
// crime of being slow). Ticket 431 asked which of the two shapes is the
// mainstream one; the answer is the fact, with a monotonic ceiling as the only
// duration in sight, because the ceiling is spent ONLY when the fact never
// arrives.
//
// MEASURED, NOT ASSUMED (ticket 431 AC#0, this host, 2026-11-12). The
// boot-complete line "resident event loop running" first becomes observable at
// 2.80s / 3.31s / 3.49s / 3.28s / 3.32s / 3.18s / 3.86s / 3.71s / 4.11s /
// 3.35s over ten boots, and 2.77-4.32s over six more with the log directory
// blocked; the VAD construction that dominates it (speech.NewDetector,
// resident_voice_windows.go:713) is 2.4-2.6s of that, and a boot with no voice
// chain finishes at 1.75-1.84s. So 60s is ~14x the worst measured boot: a
// machine ten times slower than this one still passes, and a boot that is
// actually wedged still reports inside the package's own -timeout.
//
// The ceiling is NOT what a real fault costs, and that is the point of
// waitForBootFact127's second exit below: a child that dies during boot is
// reported the moment it is reaped, so the ceiling is spent only by the one
// failure mode it is actually for - a process that is alive and never arrives.
const bootFactCeiling127 = 60 * time.Second

// bootFact127 is the outcome of one wait for a boot fact.
type bootFact127 struct {
	// ok says the fact became observable before either exit below fired.
	ok bool
	// elapsed is how long the wait took, whichever way it ended.
	elapsed time.Duration
	// exitedEarly says the child was reaped without ever producing the fact:
	// the boot died rather than merely being slow. This is a DIFFERENT finding
	// from a spent ceiling and the failure messages say which one happened.
	exitedEarly bool
	// exitStatus is the child's ending, when it ended.
	exitStatus string
}

// waitForBootFact127 waits for one boot fact on a monotonic ceiling, and gives
// up EARLY when the child it is watching has already exited: a boot that died
// cannot produce the fact however long the ceiling is, so making that case sit
// out the full 60s would be the exact cost ticket 431 refused to accept from a
// bigger fixed number. Two exits, one duration.
//
// The budget is an observe.Timeout (D42#9: monotonic clock only, never a
// wall-clock delta), so a system clock jump cannot shorten or extend it. The
// elapsed time is returned rather than only logged so a caller can put the
// measurement in its own failure message: "waited 60s for X" and "died after
// 3s without X" are different findings, and a reader should not have to guess
// which one a red is reporting.
func waitForBootFact127(leg *residentLeg, budget time.Duration, fact func() bool) bootFact127 {
	tm := observe.NewTimeout(budget)
	for {
		if fact() {
			return bootFact127{ok: true, elapsed: tm.Elapsed()}
		}
		// reap(0) is non-blocking and idempotent, so this neither steals the
		// Wait result from leg.stop nor races the reaper goroutine. cmd.Wait
		// also flushes both copy goroutines before the reaper reports, so once
		// this is true the two buffers hold every byte the child wrote - which
		// is why the fact is re-read here rather than assumed absent: a child
		// that printed the fact and exited inside the same tick is a pass, not
		// an early death.
		if leg.reap(0) {
			if fact() {
				return bootFact127{ok: true, elapsed: tm.Elapsed()}
			}
			return bootFact127{
				elapsed:     tm.Elapsed(),
				exitedEarly: true,
				exitStatus:  leg.exitStatus(),
			}
		}
		if tm.Expired() {
			return bootFact127{elapsed: tm.Elapsed()}
		}
		time.Sleep(50 * time.Millisecond)
	}
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

	// Wait for the FACT that the listener booked itself on disk, not for a
	// duration. The old shape handed this a fixed 200 x 50ms = 10s and the
	// comment below claimed that was "~20 flush ticks"; ticket 431 AC#0 measured
	// the record arriving at 611ms, so the budget never bound at all and the
	// number was doing no work. What the case needs is the fact, and the ceiling
	// is only what a boot that never produces it spends before reporting red.
	sawRecord := waitForBootFact127(leg, bootFactCeiling127, func() bool { return sinkHasARecord(sinkDir) })
	if !sawRecord.ok {
		files := jsonlFilesUnder(t, dataDir)
		leg.stop()
		if sawRecord.exitedEarly {
			t.Fatalf("the resident process exited (%s) after %s without writing a log file under %s: "+
				"the boot died before its listener booked itself.\n"+
				"data root tree holds %d *.jsonl: %v\n%s\n"+
				"that is the claim this case exists for: the leg owner actually runs has no listener on disk.",
				sawRecord.exitStatus, sawRecord.elapsed, sinkDir, len(files), files, leg.console())
		}
		t.Fatalf("the resident process wrote no log file under %s within %s.\n"+
			"data root tree holds %d *.jsonl: %v\n%s\n"+
			"that is the claim this case exists for: the leg owner actually runs has no listener on disk.",
			sinkDir, bootFactCeiling127, len(files), files, leg.console())
	}
	t.Logf("resident listener booked itself on disk after %s", sawRecord.elapsed)

	recs := readResidentSink(t, sinkDir)
	if len(recs) == 0 {
		t.Fatalf("%s holds a pipeline file with no record in it", sinkDir)
	}
	// Ticket 130 moved what record 0 is. Read the header of this file for why the
	// claim below is the same job with a stronger fact behind it, and for what
	// used to make this assertion pass.
	installIdx := installRecordIndex127(recs)
	if installIdx < 0 {
		t.Fatalf("no %q record in the file the resident process wrote (records: %v)",
			residentInstallMsg, msgsOf127(recs))
	}
	first := recs[0]
	if first.Msg != residentEarlyResolverMsg {
		t.Errorf("record 0 = %q, want %q: the boot-time sealing verdict must be replayed into this leg's "+
			"file ahead of the listener's own booking. If the booking is record 0 again, the replay in "+
			"cmd/wisp/logsink.go is gone and ticket 130's gap is open (records: %v)",
			first.Msg, residentEarlyResolverMsg, recs)
	}
	if first.Resolver == "" {
		t.Errorf("record 0 carries no resolver attribute, so it is not the sealing-seam verdict this leg booted with: %v", first)
	}
	if installIdx != 1 {
		t.Errorf("install booking is record %d, want 1: the replayed early records come first and nothing else "+
			"may sit before the booking (records: %v)", installIdx, msgsOf127(recs))
	}
	booking := recs[installIdx]
	earlyTime, earlyErr := time.Parse(time.RFC3339Nano, first.Time)
	bookingTime, bookingErr := time.Parse(time.RFC3339Nano, booking.Time)
	switch {
	case earlyErr != nil || bookingErr != nil:
		t.Errorf("the stamps are not the RFC3339Nano form the pipeline writes: early %q (%v), booking %q (%v)",
			first.Time, earlyErr, booking.Time, bookingErr)
	case !earlyTime.Before(bookingTime):
		t.Errorf("the replayed record's stamp %s is not older than the booking's %s: a replay that re-dates the "+
			"record it rescued says the listener decided to speak, not that the boot did", earlyTime, bookingTime)
	}
	if booking.Level != "INFO" {
		t.Errorf("install record level = %q, want INFO", booking.Level)
	}
	// WISP_TEST_DATA_DIR is an identity contract (internal/proc/envfork.go): it
	// comes back verbatim, so this comparison is exact and not a tree match.
	if booking.Dir != sinkDir {
		t.Errorf("install record names dir = %q, want %q: the resident listener must book the data root it was booted with",
			booking.Dir, sinkDir)
	}
	if booking.MinLevel != "info" {
		t.Errorf("install record min_level = %q, want the schema default %q", booking.MinLevel, "info")
	}

	// Ticket 431 AC#0's finding, and the reason this wait is here rather than a
	// bigger number on the old one: this case used to stop the child as soon as
	// the on-disk record appeared (measured: 611ms), and THEN assert the console
	// carried the boot-complete line - which this host prints at 2.8-4.3s,
	// because the voice chain's VAD construction (speech.NewDetector,
	// resident_voice_windows.go:713) alone costs 2.4-2.6s. So the child was
	// killed roughly 2.7s BEFORE the fact the next assertion is about could
	// exist, and the case reported "the child never reached the event loop"
	// about a process that had simply not been allowed to get there. That is not
	// a budget being too small; it is the assertion's own fact being unobservable
	// at the moment it was read.
	//
	// The fix is to wait for the fact before stopping the child. The fact is
	// unchanged - still the shipped binary's own "resident event loop running"
	// line, still asserted below - and the ceiling only decides how long a boot
	// that never prints it takes to report red.
	sawLoop := waitForBootFact127(leg, bootFactCeiling127, func() bool { return leg.stdout.has(residentReachedLoop) })
	if sawLoop.ok {
		t.Logf("resident leg reached its event loop after %s", sawLoop.elapsed)
	}

	leg.stop()
	// The console half of the fan-out, read from the child's own captured
	// stderr: this is what ties the FILE to THIS process, so the pass cannot be
	// bought by another writer that happened to use the same directory.
	if !leg.stderr.has(residentInstallMsg) {
		t.Errorf("the install record never reached the child's console; the fan-out went one way only.\n%s", leg.console())
	}
	// Ticket 130's other half, on this leg: the copy that got a second home in
	// the file keeps its first one. Zero means the mirror was switched off to
	// make the file cheap; two means the replay double-printed.
	if n := strings.Count(leg.stderr.String(), residentEarlyResolverMsg); n != 1 {
		t.Errorf("the boot-time sealing verdict reached the child's console %d times, want exactly 1 "+
			"(a-with-mirror: never switched off, never doubled).\n%s", n, leg.console())
	}
	if !sawLoop.ok {
		if sawLoop.exitedEarly {
			t.Errorf("the child exited (%s) after %s without ever reaching the event loop, so the disk read above proves little.\n%s",
				sawLoop.exitStatus, sawLoop.elapsed, leg.console())
		} else {
			t.Errorf("the child never reached the event loop within %s, so the disk read above proves little.\n%s",
				bootFactCeiling127, leg.console())
		}
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

	// Both waits below are on a FACT, on the same monotonic ceiling, and each
	// reports what it spent. The old shape gave the refusal 200 x 50ms = 10s and
	// then the boot 20 x 50ms = 1s, which is the ticket-431 defect in its purest
	// form: two durations picked by hand for two facts that arrive at completely
	// different times. The refusal is printed at ~0.1s (it is the sink's first
	// act); the boot-complete line cannot arrive before the voice chain finishes,
	// which this host measures at 2.8-4.3s. So the second number was not a
	// tolerance, it was a coin flip against a construction that takes three to
	// four times as long - and on this host it lost 5 runs out of 5.
	//
	// The ceiling is shared rather than doubled on purpose: it bounds the WHOLE
	// boot, and each fact gets the remainder of that one budget, so a boot that
	// never produces either fact cannot burn 2 x 60s.
	bootCeiling := observe.NewTimeout(bootFactCeiling127)
	sawRefusal := waitForBootFact127(leg, bootCeiling.Remaining(), func() bool {
		return leg.stderr.has(residentRefusalPrefix)
	})
	sawLoop := waitForBootFact127(leg, bootCeiling.Remaining(), func() bool {
		return leg.stdout.has(residentReachedLoop)
	})
	leg.stop()

	if sawRefusal.ok {
		t.Logf("the refusal reached the console after %s", sawRefusal.elapsed)
	}
	if sawLoop.ok {
		t.Logf("the resident leg reached its event loop after %s", sawLoop.elapsed)
	}

	if !sawRefusal.ok {
		t.Fatalf("a log directory that cannot be opened produced no refusal on the console within %s.\n%s\n"+
			"the operator has to be told the security notices are not being kept; that is the other half of the 20 lines.",
			bootFactCeiling127, leg.console())
	}
	if !sawLoop.ok {
		if sawLoop.exitedEarly {
			t.Fatalf("the refusal stopped the app: the child exited (%s) after %s without ever reaching the event loop.\n%s",
				sawLoop.exitStatus, sawLoop.elapsed, leg.console())
		}
		t.Fatalf("the refusal stopped the app: the child never reached the event loop within %s.\n%s",
			bootFactCeiling127, leg.console())
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
// handler raises as os.Interrupt), and the file then has to carry all three
// halves in this order: the boot-time record ticket 130 replays, the install
// booking, the shutdown trail last.
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
	// Ticket 131 adds this guard, and the reason it is here is a measurement:
	// running this package's full suite on a loaded host, the child took the
	// CTRL_BREAK and died at 0xc000013a before its own close flushed anything,
	// so readResidentSink returned a file with zero records and the index below
	// panicked. A panic in a Go test binary eats the red name of every case that
	// had not run yet and truncates the package output, which is how a real
	// failure (this case's own failure to hold its exit budget) reads as "there
	// were no failures". R-127-4 named this shape; the same guard already exists
	// at line 370 in the sibling case above.
	if len(recs) == 0 {
		t.Fatalf("%s holds a pipeline file with no record in it, so the ordering has nothing to be read out of; the leg exited %s",
			sinkDir, leg.exitStatus())
	}
	if recs[0].Msg != residentEarlyResolverMsg {
		t.Fatalf("record 0 = %q, want %q: the trail below cannot be read as \"the listener caught the shutdown\" "+
			"unless the file starts where the process started - and record 0 being the booking again means the "+
			"replay of the boot-time sealing verdict is gone (records: %v)",
			recs[0].Msg, residentEarlyResolverMsg, recs)
	}
	installIdx := installRecordIndex127(recs)
	if installIdx < 0 {
		t.Fatalf("no %q record in the file, so the shutdown trail has no landmark to be read against (records: %v)",
			residentInstallMsg, msgsOf127(recs))
	}
	if installIdx != 1 {
		t.Fatalf("install booking is record %d, want 1: only the records ticket 130 replays may precede the "+
			"listener's own booking (records: %v)", installIdx, msgsOf127(recs))
	}
	// Ticket 425: this is the ordering claim, re-anchored on the promise the two
	// defers actually make instead of on "the last record in the file".
	//
	// WHY THE OLD SHAPE HAD TO GO. It asserted `recs[len(recs)-1]` contained
	// residentShutdownRecord. Ticket 43 (5056a365c) then put the ten-step replay
	// loop at resident_windows.go:246 - inside the same defer as the sequence,
	// but AFTER it - so the file's last record is necessarily
	// "wisp: D38(e) step 10 (exit) ..." and a "shutdown step ..." record can
	// never be last again. That assertion was red for a reason that has nothing
	// to do with the product, and the non-empty guard below it was already green
	// (this case's own log line reports 5 shutdown records, steps [1 2 5 6 7]).
	//
	// WHY THIS SHAPE IS NOT THE WEAKENING THE TICKET FORBIDS. "The trail exists"
	// alone would be no assertion at all: that reading already passed before this
	// change. What is asserted here is the thing the two defers exist for, read
	// off the sink in its own append order (and corroborated by the stamps), as
	// three conjuncts:
	//
	//  1. the trail is not empty - the guard is kept, and it is the conjunct that
	//     fires when the per-step booking is deleted (delete the slog.Info at
	//     internal/proc/shutdown.go:130 and this case goes red: measured through
	//     a go test -overlay, not asserted in a comment);
	//  2. the steps are booked in D38(e) order - strictly increasing step
	//     numbers, which is what "in order" means for an audit trail, and which
	//     a reordered sequence (or a reordered run() roster) breaks while every
	//     record is still present;
	//  3. EVERY shutdown-step record precedes the FIRST "wisp: D38(e) step"
	//     record. This is the ordering itself: the sequence books each step while
	//     the sink is open, the replay that follows is the same still-open sink,
	//     so the trail must land entirely ahead of the replay. Move the replay
	//     ahead of the sequence, or book a step from inside the replay, and this
	//     conjunct goes red while conjuncts 1 and 2 stay green.
	//
	// The old assertion could not distinguish any of those three; it could only
	// say "some shutdown record is last", which after ticket 43 is unsatisfiable
	// and therefore measures nothing.
	var trail []sinkInstallRecord
	var trailAt []int
	for i := installIdx + 1; i < len(recs); i++ {
		if strings.Contains(recs[i].Msg, residentShutdownRecord) {
			trail = append(trail, recs[i])
			trailAt = append(trailAt, i)
		}
	}
	if len(trail) == 0 {
		t.Fatalf("the D38(e) trail is not in the file the same process closed: %d record(s) after the install, all msgs %v",
			len(recs)-installIdx-1, msgsOf127(recs))
	}
	for i := 1; i < len(trail); i++ {
		if trail[i].Step <= trail[i-1].Step {
			t.Errorf("the D38(e) trail books step %d (record %d) after step %d (record %d), steps %v: the sequence's "+
				"contract is that the steps are booked in D38(e) order, so a trail that runs backwards is not an audit of "+
				"that order even though every record is present",
				trail[i].Step, trailAt[i], trail[i-1].Step, trailAt[i-1], stepsOf127(trail))
		}
	}
	replayAt := -1
	for i := installIdx + 1; i < len(recs); i++ {
		if strings.HasPrefix(recs[i].Msg, residentReplayRecord) {
			replayAt = i
			break
		}
	}
	if replayAt < 0 {
		t.Errorf("no %q record after the install booking, so this case cannot read the ordering it was re-anchored on "+
			"(ticket 425): ticket 43's replay loop is what the trail is read against. If the replay was deliberately "+
			"removed or reworded, this case has to be re-anchored onto wherever the exit path now books itself - do not "+
			"delete the claim (records: %v)", residentReplayRecord, msgsOf127(recs))
	} else {
		lastTrailAt := trailAt[len(trailAt)-1]
		if lastTrailAt > replayAt {
			t.Errorf("the last shutdown-step record is record %d of %d, but the first %q record is record %d: the "+
				"D38(e) sequence must book its steps BEFORE the replay runs, in the sink the same process closes "+
				"(records: %v)", lastTrailAt, len(recs), residentReplayRecord, replayAt, msgsOf127(recs))
		}
		// The stamps corroborate the append order rather than replace it: the file
		// order above is the sink's own contract, and a record dated after the
		// replay began would say the booking happened later than the position says.
		// Equal stamps are allowed - a tie is not evidence of a regression.
		trailTime, trailErr := time.Parse(time.RFC3339Nano, trail[len(trail)-1].Time)
		replayTime, replayErr := time.Parse(time.RFC3339Nano, recs[replayAt].Time)
		switch {
		case trailErr != nil || replayErr != nil:
			t.Errorf("the stamps are not the RFC3339Nano form the pipeline writes: trail %q (%v), replay %q (%v)",
				trail[len(trail)-1].Time, trailErr, recs[replayAt].Time, replayErr)
		case trailTime.After(replayTime):
			t.Errorf("the last shutdown-step record is stamped %s, after the first replay record's %s: the booking "+
				"must happen before the replay, not during or after it", trailTime, replayTime)
		}
		t.Logf("RESIDENT LEG ORDERING: %d shutdown record(s) at records %v, first replay record %d of %d, last trail stamp %s, first replay stamp %s",
			len(trail), trailAt, replayAt, len(recs), trail[len(trail)-1].Time, recs[replayAt].Time)
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
