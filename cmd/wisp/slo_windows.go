//go:build windows

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/proc"
)

// wisp slo - the SLO sampling driver (ticket 08, SPEC-10 §3 test hook;
// measurement basis reworked by ticket 66).
//
// This subcommand is the seam scripts/slo-check.ps1 (and CI) drive: it boots
// the REAL runtime skeleton (Job Object, registry, log pipeline - telemetry
// is real, not simulated), enters one SLO posture, samples the seven D32
// metrics over the process tree for N seconds and emits a JSON verdict.
// Exit code 0 = pass, 1 = fail (verdict red), 2 = usage/environment error.
//
// WHO MEASURES WHOM (ticket 66 / docs/SLO.md registry A15, orchestrator
// ruling 2). The acceptance numbers come from OUTSIDE the measured tree: this
// process (the observer) starts a second skeleton process (the SUBJECT) in the
// Job, enters the posture, and is sampled by proc.ExternalSampler through the
// kernel's own per-pid counters. The subject simultaneously samples ITSELF the
// way tickets 08/12 did - that in-tree report is embedded under "observer"
// and its CPU row is marked observer_cost (recorded, never a gate), because a
// sampler inside the tree it measures is billing its own ReadTree burn (~1.3ms
// per read at a 250ms cadence = 0.52% all-core, i.e. above the <=0.5% gate on
// its own). The D32 limit itself is untouched: cpu_percent_all_core still gates
// at <=0.5% for Sleeping, it is just read from where the product is not also
// its own observer.
//
// Measurement protocol (docs/SLO.md §7, the spike's own methodology): after
// state entry the driver runs runtime.GC()x2 + the instrumented
// observe.ReleaseMemory, waits -settle-ms for the footprint to plateau, and
// only then opens the sampling window (CPU totals are window deltas, so the
// baseline excludes boot work). This is measurement methodology, not gate
// tuning: thresholds and gate/target classification are FROZEN in
// observe/thresholds.go.
//
// Modes:
//
//	wisp slo -state Sleeping [-seconds 5] [-leak]        acceptance run (out-of-tree gate)
//	wisp slo -subject -state Sleeping ...                measurement subject (self-sample, in-tree)
//	wisp slo -settle [-seconds 10]                       settle-row check (in-process, no CPU row)
//
// The leak fixture allocates and holds 100MB (page-touched) in the SUBJECT
// before sampling, which must flip the memory gate red - a green run proves
// the sampler is broken. The settle fixture allocates a session-shaped buffer,
// releases it through observe.ReleaseMemory (the instrumented C11 step-4 call)
// and requires the tree memory back under the Sleeping cap within the window
// AND a FreeOSMemory counter > 0.

const (
	leakFixtureBytes = 100 << 20
	settlePeakBytes  = 120 << 20

	// subjectGrace is how long the instrumented subject stays alive after ITS
	// sampling window closes, so that the report file it writes lands strictly
	// after the observer's last read. A subject write inside the window would
	// show up in the D32 disk_write_ops gate and put the instrument's own cost
	// back into the measured number - the exact class of bug ticket 66 removes.
	subjectGrace = 3 * time.Second
	// subjectIdleSlack is how much longer than the window the measured subject
	// is willing to sit still before declaring its observer gone.
	subjectIdleSlack = 2 * time.Minute
	// subjectReadyBudget / subjectReportBudget are monotonic waits (D42#9);
	// subjectPollInterval is how often the files are re-checked.
	subjectReadyBudget  = 60 * time.Second
	subjectReportBudget = 30 * time.Second
	subjectPollInterval = 20 * time.Millisecond
)

// leakHold keeps the forced leak alive until process exit.
var leakHold []byte

// settleHold keeps the settle fixture alive until ReleaseMemory runs.
var settleHold []byte

type sloRun struct {
	Mode      string  `json:"mode"` // "state" | "subject" | "subject-in-tree" | "settle"
	State     string  `json:"state"`
	Posture   string  `json:"posture"`
	StartedAt string  `json:"started_at"`
	Seconds   float64 `json:"seconds"`
	// SubjectPID is the MEASURED child (product posture, no sampler inside it)
	// when this run samples out-of-tree.
	SubjectPID uint32 `json:"subject_pid,omitempty"`
	// ObserverSubjectPID is the INSTRUMENTED child: the same posture with a
	// sampler running inside it, whose self-report lands in Observer.
	ObserverSubjectPID uint32 `json:"observer_subject_pid,omitempty"`
	// ObserverSubjectPrivateWS is ONE out-of-tree read of the instrumented
	// subject, taken after the window closes. It exists so a reader can tell
	// whether the in-tree report's footprint is that process's real memory or
	// the sampler's own scratch: an in-process ReadTree keeps a multi-megabyte
	// SystemProcessInformation buffer on the measured tree's heap, which
	// docs/SLO.md §7 already flagged for the spike.
	ObserverSubjectPrivateWS int64  `json:"observer_subject_private_working_set_bytes,omitempty"`
	ObserverSubjectErr       string `json:"observer_subject_read_error,omitempty"`
	// Report holds the ACCEPTANCE numbers: measured out-of-tree in "state"
	// mode (proc.ExternalSampler over the subject pid), in-process otherwise.
	Report *observe.StateReport `json:"report,omitempty"`
	// Observer holds the subject's own in-tree report of the same posture and
	// the same window length (ticket 66): recorded evidence, whose CPU row is
	// marked observer_cost and never gates.
	Observer *sloRun               `json:"observer,omitempty"`
	Settle   *observe.SettleReport `json:"settle,omitempty"`
	Pass     bool                  `json:"pass"`
}

// sloSubject is one running measurement subject: the child process plus the
// two handshake files it exchanges with the observer.
type sloSubject struct {
	cmd       *exec.Cmd
	pid       uint32
	dir       string
	readyPath string
	outPath   string
	// readReportFile is the report loop's ONLY touch of the subject's report
	// file. Production leaves it nil and the loop calls os.ReadFile; a test
	// hands in a scripted sequence of readings instead. The seam exists for one
	// claim that cannot be made any other way on this host: "a report that is
	// genuinely corrupt fails closed ON THE SPOT" has to be shown by counting
	// reads, because the machine ticket 144 runs on is not allowed to be
	// trusted for elapsed-time readings (three other processes may be sampling).
	readReportFile func(path string) ([]byte, error)
}

func cmdSLOUsage() {
	fmt.Print(`wisp slo - SLO sampling driver (ticket 08, measurement basis reworked by ticket 66)

Usage:
  wisp slo -state <Sleeping|Armed|Warm|Conversation|PanelOpen|WorkPeak> [flags]
  wisp slo -settle [flags]
  wisp slo -subject [-self-sample] -state <state> -ready-file f [flags]  (child)

Flags:
  -state string      SLO state to sample (required unless -settle)
  -settle            run the settle-row check (peak buffer, ReleaseMemory,
                     tree memory back under the Sleeping cap within 10s)
  -subject           act as a measurement child of another wisp slo process:
                     boot, settle, write the readiness marker; never spawns a
                     child of its own
  -self-sample      subject mode: also run the LEGACY in-tree self-sampling
                     (sampler inside the measured tree) and write the report
                     -out asks for; this is the observer-cost record
  -ready-file string subject mode: readiness marker path (required there)
  -seconds float     sampling window in seconds (default 5)
  -interval-ms int   sampling interval in ms (default 250)
  -settle-ms int     settle delay before the sampling window (default 2000)
  -leak              force the leak fixture: allocate + hold 100MB in the
                     subject (must flip the gate red; used by slo-check.ps1
                     self-test)
  -out string        write the JSON report to this path (default stdout)

Acceptance numbers are read OUT OF TREE (ticket 66): a state run starts a
subject process in the Job and measures it from the parent; the subject's own
in-tree self-sample rides along under "observer" as a record, and its CPU row
carries observer_cost=true because the sampler there is inside what it measures.

Environment: runs under WISP_ENV=test by default (no single-instance mutex,
hermetic data dir via WISP_TEST_DATA_DIR).
`)
}

// cmdSLO is `wisp slo`, the SLO sampling driver (ticket 08).
//
// WISP-LEG-SINK-RULING: this leg resolves the JSONL pipeline itself
// (observe.InitLog, slo_windows.go) rather than going through installLogSink,
// because the pipeline is its OUTPUT and not its listener: the command's job is
// to sample a subject and print a JSON verdict, and the log files it writes
// under <data root>\logs are the sweep target its retention rule counts. Giving
// it a second, process-wide fan-out handler on top of a sink it already opens by
// hand would duplicate records into the tree it is measuring, which is a
// different claim from the one ticket 117's install makes on the other legs. The
// enumeration gate in leg_sink_gate_131_test.go reads this function's own doc
// comment for the sentence above and its body for the calls below, so the row it
// books for `wisp slo` in that gate's ledger is install=false with records=true
// and ruled=true, and the state word on that line is "ruled" rather than
// "nailed". This sentence is the reason that is a decision and not an
// omission.
//
// WISP-LEG-COVERAGE-RULING: slo is dispatched by main and driven by no case in
// this package: its driver is scripts/slo-check.ps1 and the CI SLO jobs, which
// start a real subject process (ticket 08's protocol), and a second reading of it
// inside this test binary would sample a process that is not the one the gate
// measures. The ruling above this one is why it also has no listener nail.
func cmdSLO(args []string) int {
	fs := flag.NewFlagSet("slo", flag.ContinueOnError)
	fs.Usage = cmdSLOUsage
	state := fs.String("state", "", "SLO state to sample")
	settle := fs.Bool("settle", false, "run the settle-row check")
	subject := fs.Bool("subject", false, "act as a measurement child of another wisp slo process (boot, settle, report ready)")
	selfSample := fs.Bool("self-sample", false, "subject mode: also run the legacy in-tree self-sampling and write the report the observer collects")
	readyFile := fs.String("ready-file", "", "subject mode: readiness marker path")
	seconds := fs.Float64("seconds", 5, "sampling window seconds")
	intervalMS := fs.Int("interval-ms", 250, "sampling interval ms")
	settleMS := fs.Int("settle-ms", 2000, "settle delay before the sampling window (SLO.md §7 protocol)")
	leak := fs.Bool("leak", false, "allocate + hold 100MB (leak fixture)")
	out := fs.String("out", "", "output JSON path")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "wisp slo: unexpected argument %q\n", fs.Arg(0))
		return 2
	}
	if *seconds <= 0 || *intervalMS <= 0 {
		fmt.Fprintln(os.Stderr, "wisp slo: -seconds and -interval-ms must be positive")
		return 2
	}
	if *settle && *subject {
		fmt.Fprintln(os.Stderr, "wisp slo: -settle and -subject are mutually exclusive")
		return 2
	}
	if *subject && *readyFile == "" {
		fmt.Fprintln(os.Stderr, "wisp slo: -subject needs -ready-file")
		return 2
	}
	if !*settle && *state == "" {
		fmt.Fprintln(os.Stderr, "wisp slo: -state is required (or use -settle)")
		return 2
	}
	if *settle && *leak {
		fmt.Fprintln(os.Stderr, "wisp slo: -settle and -leak are mutually exclusive")
		return 2
	}
	if *state != "" && !observe.SLOState(*state).Valid() {
		fmt.Fprintf(os.Stderr, "wisp slo: unknown state %q (want one of %s)\n", *state, sloStateList())
		return 2
	}

	// Sampling runs hermetic by default: test env registers no mutex and
	// honors WISP_TEST_DATA_DIR for the data dir.
	if os.Getenv("WISP_ENV") == "" {
		_ = os.Setenv("WISP_ENV", "test")
	}
	attachParentConsole()

	env, err := buildinfo.ResolveEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "wisp slo: %v\n", err)
		return 2
	}
	rt, err := proc.Boot(env)
	if err == proc.ErrAlreadyRunning {
		fmt.Fprintln(os.Stderr, "wisp slo: another instance is running (set WISP_ENV=test)")
		return 2
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "wisp slo: boot failed: %v\n", err)
		return 2
	}

	// Real log pipeline (redaction + rolling + log-flusher goroutine).
	pipeline, err := observe.InitLog(observe.LogConfig{
		Dir:   filepath.Join(rt.Layout.DataDir, "logs"),
		Level: "info",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "wisp slo: log pipeline: %v\n", err)
		rt.Shutdown(false)
		return 2
	}
	pipeline.InstallAsDefault()

	sampler := observe.NewSampler(proc.NewTreeSampler(rt.Job), rt.Registry)
	run := &sloRun{
		State:     *state,
		StartedAt: observe.WallTimestampUTC(observe.NowWallUTC()),
		Seconds:   *seconds,
	}
	ctx := context.Background()
	exitCode := 0

	switch {
	case *settle:
		run.Mode = "settle"
		run.State = string(observe.SLOSleeping)
		run.Posture = "skeleton+settle-fixture"
		// In-process by construction and there is no CPU row in the settle
		// verdict, so the ticket 66 basis question does not arise here.
		run.Settle = runSettle(ctx, sampler, proc.NewTreeSampler(rt.Job))
		run.Pass = run.Settle.Pass
	case *subject:
		code, err := runSubject(ctx, sampler, run, *state, *seconds,
			time.Duration(*intervalMS)*time.Millisecond, *settleMS, *readyFile, *leak, *selfSample)
		if err != nil {
			fmt.Fprintf(os.Stderr, "wisp slo: %v\n", err)
			_ = pipeline.Close()
			rt.Shutdown(false)
			return code
		}
	default:
		code, err := runOutOfTree(ctx, rt, run, *state, *seconds,
			time.Duration(*intervalMS)*time.Millisecond, *settleMS, *leak)
		if err != nil {
			fmt.Fprintf(os.Stderr, "wisp slo: %v\n", err)
			_ = pipeline.Close()
			rt.Shutdown(false)
			return code
		}
	}

	if err := writeSLO(run, *out); err != nil {
		fmt.Fprintf(os.Stderr, "wisp slo: %v\n", err)
		rt.Shutdown(false)
		return 2
	}

	// Close the log pipeline before the D38(e) sequence (workers stopped
	// already; the pipeline must be flushed, not lost).
	_ = pipeline.Close()
	records := rt.Shutdown(false)
	for _, rec := range records {
		if rec.Err != nil {
			fmt.Fprintf(os.Stderr, "wisp slo: shutdown step %d (%s) failed: %v\n", rec.Step, rec.Name, rec.Err)
		}
	}
	if !run.Pass {
		exitCode = 1
	}
	return exitCode
}

// runOutOfTree is the acceptance path (ticket 66 AC#3). It starts TWO children
// in the Job over the same window:
//
//   - the MEASURED subject: the product posture and nothing else (boot the real
//     skeleton, apply the SLO.md §7 settle protocol, then sit still). The
//     observer samples it from outside through proc.ExternalSampler and those
//     rows are the D32 gates.
//   - the INSTRUMENTED subject: the same posture plus a sampler running inside
//     it, i.e. exactly what `wisp slo` measured up to ticket 66. Its own
//     in-tree report rides along as the observer-cost record.
//
// The pair is what proves the difference is the OBSERVER and not the location
// of the read: both children are the same binary in the same posture sampled at
// the same cadence in the same window, and the only difference between them is
// whether a ReadTree happens inside what it measures.
//
// Fail-closed everywhere: if a subject cannot be started, measured, or cannot
// hand back its own report, this exits 2 - silently degrading to the old
// in-tree basis is exactly what the ticket exists to prevent.
func runOutOfTree(ctx context.Context, rt *proc.Runtime, run *sloRun,
	state string, seconds float64, interval time.Duration, settleMS int, leak bool,
) (int, error) {
	run.Mode = "state"
	run.Posture = "skeleton"
	settleBeforeSample(settleMS)

	base := subjectSpec{state: state, seconds: seconds, interval: interval, settleMS: settleMS, leak: leak}
	subject, err := startSubject(rt, base)
	if err != nil {
		return 2, err
	}
	defer subject.stop()
	run.SubjectPID = subject.pid

	instrumented, err := startSubject(rt, subjectSpec{
		state: state, seconds: seconds, interval: interval, settleMS: settleMS, selfSample: true,
	})
	if err != nil {
		return 2, err
	}
	defer instrumented.stop()
	run.ObserverSubjectPID = instrumented.pid

	sampler := observe.NewSampler(proc.NewExternalSampler(subject.pid), rt.Registry)
	sampler.MarkTransition("boot", state)
	rep, err := sampler.SampleState(ctx, observe.SLOState(state), interval,
		time.Duration(seconds*float64(time.Second)))
	if err != nil {
		return 2, fmt.Errorf("out-of-tree sampling failed: %w", err)
	}
	run.Report = rep

	observer, err := instrumented.collectReport()
	if err != nil {
		return 2, fmt.Errorf("in-tree record unavailable (fail-closed, no silent downgrade to the tree basis): %w", err)
	}
	run.Observer = observer
	run.Pass = rep.Pass && observer.Pass
	return 0, nil
}

// runSubject is a child started by runOutOfTree. In -self-sample mode it enters
// the posture the way `wisp slo` always has, marks itself ready AFTER the settle
// protocol (so the observer's window opens on a plateau and the marker write
// itself stays outside the window), then samples ITSELF in-tree; that report is
// the observer-cost record. Without -self-sample it is the measured subject:
// the same boot, the same settle, then nothing at all, because a product
// posture that is busy being an instrument measures the instrument.
func runSubject(ctx context.Context, sampler *observe.Sampler, run *sloRun, state string,
	seconds float64, interval time.Duration, settleMS int, readyPath string, leak, selfSample bool,
) (int, error) {
	run.Mode = "subject"
	run.Posture = "skeleton"
	run.SubjectPID = uint32(os.Getpid())
	settleBeforeSample(settleMS)
	if leak {
		run.Posture = "skeleton+leak-fixture"
		holdLeakFixture()
	}
	if err := writeSubjectReady(readyPath, run.SubjectPID); err != nil {
		return 2, err
	}
	if !selfSample {
		// The measured subject does nothing until the observer is done with it.
		// The bound is a config sum (window + slack), never a wall-clock
		// difference; if it ever fires the observer fails closed instead of
		// measuring a corpse.
		time.Sleep(time.Duration(seconds*float64(time.Second)) + subjectIdleSlack)
		return 2, fmt.Errorf("wisp slo: subject outlived its observer's window")
	}
	// Inside the tree: demote the CPU row to a record (observe/thresholds.go
	// keeps the frozen limit; only the gate flag moves).
	sampler.MarkObserverInsideTree()
	sampler.MarkTransition("boot", state)
	rep, err := sampler.SampleState(ctx, observe.SLOState(state), interval,
		time.Duration(seconds*float64(time.Second)))
	if err != nil {
		return 2, fmt.Errorf("subject sampling failed: %w", err)
	}
	run.Mode = "subject-in-tree"
	run.Report = rep
	run.Pass = rep.Pass
	// Stay parked so the observer's window closes first; then the report write
	// is outside the measured interval (see subjectGrace).
	time.Sleep(subjectGrace)
	return 0, nil
}

// subjectSpec describes one child runOutOfTree starts.
type subjectSpec struct {
	state      string
	seconds    float64
	interval   time.Duration
	settleMS   int
	leak       bool // hold the 100MB leak fixture in this child
	selfSample bool // run the legacy in-tree self-sampling and write a report
}

// startSubject spawns this executable in -subject mode inside the Job (so
// KILL_ON_JOB_CLOSE can never leave an orphan behind to pollute the next
// person's measurement), in its own hermetic data dir so the pipelines never
// share a file, and waits for the readiness marker.
func startSubject(rt *proc.Runtime, spec subjectSpec) (*sloSubject, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("wisp slo: os.Executable: %w", err)
	}
	dir, err := os.MkdirTemp("", "wisp-slo-subject-")
	if err != nil {
		return nil, fmt.Errorf("wisp slo: subject dir: %w", err)
	}
	s := &sloSubject{
		dir:       dir,
		readyPath: filepath.Join(dir, "ready"),
		outPath:   filepath.Join(dir, "subject-report.json"),
	}
	args := []string{
		"slo", "-subject", "-state", spec.state,
		"-seconds", fmt.Sprintf("%g", spec.seconds),
		"-interval-ms", fmt.Sprint(int(spec.interval / time.Millisecond)),
		"-settle-ms", fmt.Sprint(spec.settleMS),
		"-ready-file", s.readyPath,
		"-out", s.outPath,
	}
	if spec.selfSample {
		args = append(args, "-self-sample")
	}
	if spec.leak {
		args = append(args, "-leak")
	}
	cmd := exec.Command(exe, args...)
	cmd.Env = append(os.Environ(), proc.TestDataDirEnv+"="+filepath.Join(dir, "data"))
	// Nothing is streamed: a subject that wrote to the shared console would
	// put its I/O counters back into the very window being measured.
	cmd.Stdout, cmd.Stderr = nil, nil
	// StartInJob starts AND assigns (and kills the child itself if the
	// assignment fails), so a subject can never escape the Job and outlive the
	// observer's cleanup.
	p, err := rt.Job.StartInJob(cmd)
	if err != nil {
		_ = os.RemoveAll(dir)
		return nil, fmt.Errorf("wisp slo: start subject in job: %w", err)
	}
	s.cmd = cmd
	s.pid = uint32(p.Pid)
	if err := s.waitReady(); err != nil {
		s.stop()
		return nil, err
	}
	return s, nil
}

// waitReady blocks until the subject reports steady state, or fails closed if
// it dies or never gets there (budgets are monotonic, D42#9).
func (s *sloSubject) waitReady() error {
	budget := observe.NewTimeout(subjectReadyBudget)
	for {
		if fileHas(s.readyPath, "ready=1") {
			return nil
		}
		if s.exited() {
			return fmt.Errorf("wisp slo: subject %d exited early (code %d) before reporting ready",
				s.pid, s.exitCode())
		}
		if budget.Expired() {
			return fmt.Errorf("wisp slo: subject %d never reported ready within %s", s.pid, subjectReadyBudget)
		}
		time.Sleep(subjectPollInterval)
	}
}

// subjectReportState classifies ONE look at the subject's report file. The two
// states this ticket separates are reportUnwritten and reportCorrupt: before
// ticket 144 the loop had no word for the first one, so a report whose tail was
// still in flight was reported with the same sentence as a report that is
// actually broken.
type subjectReportState int

const (
	// reportUnwritten - what is on disk is a PREFIX of the document the subject
	// is still writing: zero bytes, padding only, or a value the decoder ran
	// out of input inside. Nothing in these bytes contradicts a subject report,
	// so the only honest reading is "look again", and the budget is what
	// decides whether looking again was worth it.
	reportUnwritten subjectReportState = iota
	// reportCorrupt - what is on disk CONTRADICTS the document the subject
	// writes: a head that cannot begin it, a syntax error that is not at the end
	// of the input, a second value after the first one closed, or a document
	// that parsed cleanly and carries no state report. No amount of waiting
	// changes any of those, so they fail closed on the spot (ticket 128 / 136:
	// never a silent downgrade to the tree basis).
	reportCorrupt
	// reportComplete - exactly one sloRun and nothing after it.
	reportComplete
)

// String keeps a failure in a test reading like a sentence about the file, not
// like a number the reader has to look up.
func (s subjectReportState) String() string {
	switch s {
	case reportUnwritten:
		return "unwritten"
	case reportCorrupt:
		return "corrupt"
	default:
		return "complete"
	}
}

// subjectReportRead is one observation of the file: the bytes, where inside
// them the document stopped, and what that means. It is a value, not a side
// effect, so the loop can carry its LAST reading into the error that gives up -
// which is the difference between "unexpected end of JSON input" and "waited
// 33s, the file never grew past 1 431 bytes".
type subjectReportRead struct {
	state  subjectReportState
	report *sloRun
	bytes  int
	// offset is where inside these bytes the document stopped. It means exactly
	// two things (ticket 147 pins both): -1 when no decoder ran at all - nothing
	// read yet, no file yet, a file that would not open - and N >= 0 when one
	// did. On reportComplete N is the end of the document; on reportUnwritten
	// the document ran out of input, so N is the last byte that arrived and the
	// missing part is what follows it. On reportCorrupt N is the byte position the
	// decoder objects AT, counted the way the decoder counts it - the byte it
	// objects to is included, so a bad byte at index K is named K+1 (see
	// contradictionOffset). Of the 20 shapes .scratch/wisp/probes/149 hands in, 18
	// classify as corrupt and none of those 18 produced a name smaller than 1; the
	// two corrupt legs with no decoder error keep the end of what DID close (a
	// second value after the first, or a document that parsed cleanly and carries
	// no state report), and only the first of those prints its position.
	// What this comment claimed until ticket 149 - that on reportCorrupt the
	// number "is 0 whenever it could not begin reading a value at all" - was
	// measured false the same week it was written: 8 of those 18 printed 0, and two
	// of them had a perfectly good object head behind the 0 (a document whose only
	// fault is 0xff landing at index 40, and one whose first 500 bytes parse and
	// whose byte 500 is '@'). It is recorded here rather than deleted quietly,
	// because the next reader meeting "at offset 0" would otherwise re-derive the
	// same wrong rule. Case 11
	// (TestSLO149CorruptSentenceNamesThePositionTheDecoderObjectedAt) is the
	// instrument that makes the rule above a checkable fact, and case 12 the one
	// that keeps the fill off the two legs that have no decoder error.
	// Only the -1 leg is never printed: a note-bearing observation renders its
	// note instead (see summary).
	offset int64
	// note is what the loop adds for observations the classifier never saw
	// (no file yet, file unreadable this instant).
	note string
	// err is non-nil only for reportCorrupt.
	err error
}

// summary renders an observation for an error message.
func (o subjectReportRead) summary() string {
	if o.note != "" {
		return fmt.Sprintf("%d bytes read, %s", o.bytes, o.note)
	}
	switch o.state {
	case reportUnwritten:
		return fmt.Sprintf("%d bytes read, document still open at offset %d: the tail had not arrived", o.bytes, o.offset)
	case reportCorrupt:
		return fmt.Sprintf("%d bytes read, report corrupt: %v", o.bytes, o.err)
	default:
		return fmt.Sprintf("%d bytes read, complete document at offset %d", o.bytes, o.offset)
	}
}

// contradictionOffset is the position INSIDE THE BYTES that a decoder error
// names, for a document that contradicts a subject report.
//
// Why not dec.InputOffset(), which is what the corrupt leg used to print:
// measured over the 20 corrupt shapes of .scratch/wisp/probes/149 (see
// probe-corrupt-census.log there), InputOffset() reports 0 for 8 of them - among
// them a document whose first 500 bytes parse as a perfectly good object head
// ('@' written at index 500), one whose only fault is a second comma at index 26,
// and one whose whole head is the subject's own document until 0xff lands at index
// 1341. On this leg that 0 is an artefact of how the decoder buffers, not a fact
// about the file, and it is precisely the number a reader of "at offset 0" takes
// to mean "nothing arrived" - the same misdirection ticket 147 removed from the
// unwritten leg.
//
// What the decoder's own error names IS a position, and it counts the byte it
// objected to as consumed: a bad byte at index K is named as K+1 (measured over
// the same 20 shapes; the smallest name ever produced was 1, for a file whose
// very first byte is not JSON, and no shape produced 0). *json.SyntaxError covers
// a head or a byte that cannot be part of the document, *json.UnmarshalTypeError
// covers a value that closes where a field of another type should be; both name
// their position. When an error names neither - no shape in that probe reached it,
// so this fallback leg carries no witness today and ticket 149's acceptance note
// says so - the position the decoder stopped at is all there is, so it is kept.
func contradictionOffset(err error, inputOffset int64) int64 {
	var syn *json.SyntaxError
	if errors.As(err, &syn) {
		return syn.Offset
	}
	var typ *json.UnmarshalTypeError
	if errors.As(err, &typ) {
		return typ.Offset
	}
	return inputOffset
}

// readSubjectReport classifies one read of the subject's report file. It waits
// for nothing, touches nothing and keeps no state: the same bytes always give
// the same answer, which is what makes the retry loop below a decision instead
// of a hope.
//
// The credential for telling the two apart is the shape of writeSLO's single
// os.WriteFile (slo_windows.go, writeSLO): the whole JSON document is produced
// in memory first, then written by ONE open+write from offset 0, with no
// temporary file and no rename. A reader therefore sees either nothing or a
// PREFIX of the final document - a prefix can be missing its tail, but it can
// never have a head that the finished document does not also have. So:
//
//   - the decoder ran out of input (io.EOF on zero bytes / padding,
//     io.ErrUnexpectedEOF inside a value) => reportUnwritten, keep polling;
//   - anything else the decoder objects to, or a whole document whose content
//     is wrong => reportCorrupt, fail closed now.
//
// The one shape this cannot tell apart from a prefix is a subject that DIED
// part-way through its write; that waits out the budget and still reports red,
// now with the budget and the last byte count in the sentence.
func readSubjectReport(data []byte) subjectReportRead {
	obs := subjectReportRead{bytes: len(data), offset: -1}
	dec := json.NewDecoder(bytes.NewReader(data))
	var rep sloRun
	if err := dec.Decode(&rep); err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			// The document ran out of input: every byte that arrived was
			// consumed, so the position it stopped at IS the last byte on disk,
			// and the missing part is what would have followed it (ticket 147).
			// dec.InputOffset() is deliberately not used here - measured over all
			// 1978 prefixes of ticket 144's fixture, it reports 0 for every
			// truncated document (it names the end of the last COMPLETED token,
			// and a truncated document has none), so on this branch it carries no
			// information and printing it pointed a reader at "nothing was
			// written" when the truth was "1949 bytes are here, the tail is not".
			obs.state = reportUnwritten
			obs.offset = int64(obs.bytes)
			return obs
		}
		obs.state = reportCorrupt
		// The position this sentence prints is the one the decoder objected at,
		// not the one it buffered up to - see contradictionOffset.
		obs.offset = contradictionOffset(err, dec.InputOffset())
		obs.err = fmt.Errorf("%d bytes contradict a subject report at offset %d: %w", obs.bytes, obs.offset, err)
		return obs
	}
	obs.offset = dec.InputOffset()
	if tok, terr := dec.Token(); !errors.Is(terr, io.EOF) {
		obs.state = reportCorrupt
		obs.err = fmt.Errorf("%d bytes hold more than one document: a value closed at offset %d and %d bytes follow (next token %v, %v)",
			obs.bytes, obs.offset, int64(obs.bytes)-obs.offset, tok, terr)
		return obs
	}
	if rep.Report == nil {
		// Parsed cleanly, complete, and still not a subject report: this is the
		// leg AC#2 of ticket 144 refuses to let the retry leg swallow.
		obs.state = reportCorrupt
		obs.err = fmt.Errorf("%d bytes parsed as a subject run but carry no state report", obs.bytes)
		return obs
	}
	obs.state = reportComplete
	obs.report = &rep
	return obs
}

// collectReport is the production entry point: the budget and the cadence are
// the ones this file has always used, named exactly once (ticket 144 moves no
// threshold and no budget).
func (s *sloSubject) collectReport() (*sloRun, error) {
	return s.collectReportWithin(subjectReportBudget+subjectGrace, subjectPollInterval)
}

// collectReportWithin waits for the subject to write its own in-tree report and
// parses it, and it does NOT treat "the file exists" as "the report is there".
//
// What the three lines above this function used to claim - that the write
// lands after subjectGrace past the end of the subject's window, "so this never
// races the observer's last sample" - was argued about the wrong pair of facts.
// It says something about WHEN the write starts relative to the observer's last
// read; it says nothing about what a reader sees while that single write is in
// flight, and "the file is there" and "the bytes are all there" are two
// different states (os.WriteFile creates the file, then fills it). So the race
// this loop has to handle is not observer-vs-observer, it is reader-vs-writer,
// and the two states stay separable for the reason spelled out on
// readSubjectReport: the writer produces the whole document in memory and
// writes it once from offset 0, so whatever is visible is a prefix of it.
//
// reportUnwritten stays inside the budget and is re-read. reportCorrupt returns
// an error on the spot - fail-closed is NOT relaxed here, it is only given a
// reason: if this gives up, the sentence names the budget it spent, the byte
// count of its last read and where inside those bytes the document stopped.
func (s *sloSubject) collectReportWithin(budget, poll time.Duration) (*sloRun, error) {
	read := s.readReportFile
	if read == nil {
		read = os.ReadFile
	}
	timeout := observe.NewTimeout(budget)
	last := subjectReportRead{state: reportUnwritten, offset: -1, note: "nothing read yet"}
	for {
		data, err := read(s.outPath)
		var obs subjectReportRead
		switch {
		case errors.Is(err, os.ErrNotExist):
			obs = subjectReportRead{state: reportUnwritten, offset: -1, note: "no report file yet"}
		case err != nil:
			// An unreadable file is transient here exactly as it was before this
			// ticket (the writer has it open); the loop polls and says so at the
			// end instead of staying silent about what it last saw.
			obs = subjectReportRead{state: reportUnwritten, offset: -1, note: fmt.Sprintf("report file unreadable: %v", err)}
		default:
			obs = readSubjectReport(data)
		}
		last = obs
		switch obs.state {
		case reportComplete:
			return obs.report, nil
		case reportCorrupt:
			return nil, fmt.Errorf("wisp slo: subject report is corrupt (failing closed on the first read, not a partial write): %w", obs.err)
		}
		if s.exited() {
			return nil, fmt.Errorf("wisp slo: subject %d exited (code %d) without writing its report (%s)",
				s.pid, s.exitCode(), last.summary())
		}
		if timeout.Expired() {
			return nil, fmt.Errorf("wisp slo: subject %d never wrote a complete report within %s (last read: %s)",
				s.pid, timeout.Budget(), last.summary())
		}
		time.Sleep(poll)
	}
}

func (s *sloSubject) exited() bool {
	return s.cmd != nil && s.cmd.ProcessState != nil && s.cmd.ProcessState.Exited()
}

func (s *sloSubject) exitCode() int {
	if !s.exited() {
		return -1
	}
	return s.cmd.ProcessState.ExitCode()
}

// stop guarantees the subject this run created is gone, then reaps it. A
// subject is a pure measurement stand-in (no window, no tray, no session), so
// termination needs no dispose sequence; whatever remains is the Job's own
// KILL_ON_JOB_CLOSE accounting, which rt.Shutdown drives.
func (s *sloSubject) stop() {
	if s == nil {
		return
	}
	if s.cmd != nil && s.cmd.Process != nil && !s.exited() {
		_ = s.cmd.Process.Kill()
	}
	if s.cmd != nil {
		_ = s.cmd.Wait()
	}
	_ = os.RemoveAll(s.dir)
}

// writeSubjectReady drops the readiness marker the observer polls for. Written
// before the subject's sampling window opens, so its own write counters cannot
// land inside the measured interval.
func writeSubjectReady(path string, pid uint32) error {
	body := fmt.Sprintf("ready=1\npid=%d\n", pid)
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		return fmt.Errorf("wisp slo: subject ready marker: %w", err)
	}
	return nil
}

func fileHas(path, needle string) bool {
	b, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.Contains(string(b), needle)
}

// runSettle executes the settle-row fixture: allocate a session-shaped
// buffer (page-touched), release it through the instrumented
// observe.ReleaseMemory (C11 step 4), then require the tree memory back
// under the Sleeping cap within the window.
func runSettle(ctx context.Context, sampler *observe.Sampler, reader *proc.TreeSampler) *observe.SettleReport {
	sampler.MarkTransition("boot", string(observe.SLOWarm))
	// Trim boot garbage first so the check measures the fixture delta,
	// not startup residue (same SLO.md §7 protocol).
	settleBeforeSample(2000)
	settleHold = make([]byte, settlePeakBytes)
	for i := 0; i < len(settleHold); i += 4096 {
		settleHold[i] = 1
	}
	peak := int64(settlePeakBytes)
	if m, err := reader.ReadTree(); err == nil {
		peak = m.PrivateWorkingSetBytes
	}
	// C11 dispose order: the buffer must be DEAD before the release call -
	// debug.FreeOSMemory's internal GC can only return memory nothing
	// references. (Getting this order wrong is exactly the kind of bug the
	// settle gate exists to catch: counter > 0, memory still parked.)
	settleHold = nil
	runtime.GC()
	observe.ReleaseMemory() // C11 step 4 through the instrumented wrapper
	sampler.MarkTransition(string(observe.SLOWarm), string(observe.SLOSleeping))
	rep, err := sampler.CheckSettle(ctx, observe.SLOSleeping, 10*time.Second, 250*time.Millisecond,
		true, observe.TreeMetrics{PrivateWorkingSetBytes: peak}, peak)
	if err != nil {
		return &observe.SettleReport{TargetState: observe.SLOSleeping, Pass: false}
	}
	return rep
}

func writeSLO(run *sloRun, out string) error {
	data, err := json.MarshalIndent(run, "", "  ")
	if err != nil {
		return fmt.Errorf("json: %w", err)
	}
	if out == "" {
		fmt.Println(string(data))
		return nil
	}
	return os.WriteFile(out, data, 0o644)
}

func sloStateList() string {
	out := make([]string, 0, len(observe.SLOStates))
	for _, s := range observe.SLOStates {
		out = append(out, string(s))
	}
	return strings.Join(out, "|")
}

// holdLeakFixture allocates and page-touches 100MB, held until exit.
func holdLeakFixture() {
	leakHold = make([]byte, leakFixtureBytes)
	for i := 0; i < len(leakHold); i += 4096 {
		leakHold[i] = 1
	}
}

// settleBeforeSample applies the SLO.md §7 measurement protocol: GC twice,
// one instrumented release, then wait for the footprint to plateau so the
// sampling window measures steady state rather than boot noise.
func settleBeforeSample(ms int) {
	runtime.GC()
	runtime.GC()
	observe.ReleaseMemory()
	if ms > 0 {
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
}
