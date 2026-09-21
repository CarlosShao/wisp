package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
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

// collectReport waits for the subject to write its own in-tree report and
// parses it. The subject only writes after subjectGrace past the end of its
// window, so this never races the observer's last sample.
func (s *sloSubject) collectReport() (*sloRun, error) {
	budget := observe.NewTimeout(subjectReportBudget + subjectGrace)
	for {
		if data, err := os.ReadFile(s.outPath); err == nil {
			var rep sloRun
			if err := json.Unmarshal(data, &rep); err != nil {
				return nil, fmt.Errorf("wisp slo: subject report: %w", err)
			}
			if rep.Report == nil {
				return nil, fmt.Errorf("wisp slo: subject report carries no state report")
			}
			return &rep, nil
		}
		if s.exited() {
			return nil, fmt.Errorf("wisp slo: subject %d exited (code %d) without writing its report",
				s.pid, s.exitCode())
		}
		if budget.Expired() {
			return nil, fmt.Errorf("wisp slo: subject %d did not write its report within %s",
				s.pid, budget.Budget())
		}
		time.Sleep(subjectPollInterval)
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
