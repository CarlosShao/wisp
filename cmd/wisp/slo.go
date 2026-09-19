package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/proc"
)

// wisp slo - the SLO sampling driver (ticket 08, SPEC-10 §3 test hook).
//
// This subcommand is the seam scripts/slo-check.ps1 (and CI) drive: it boots
// the REAL runtime skeleton (Job Object, registry, log pipeline - telemetry
// is real, not simulated), enters one SLO posture, samples the seven D32
// metrics over the process tree for N seconds and emits a JSON verdict.
// Exit code 0 = pass, 1 = fail (verdict red), 2 = usage/environment error.
//
// Measurement protocol (docs/SLO.md §7, the spike's own methodology): after
// state entry the driver runs runtime.GC()x2 + the instrumented
// observe.ReleaseMemory, waits -settle-ms for the footprint to plateau, and
// only then opens the sampling window (CPU totals are window deltas, so the
// baseline excludes boot work). This is measurement methodology, not gate
// tuning: thresholds and sample counts are frozen in observe/thresholds.go.
//
// Postures, honestly scoped to what exists at ticket 08 (the JSON marks
// this in "posture"): the model-bearing states (Armed/Warm/Conversation)
// and the panel/peak postures sample the SKELETON footprint until tickets
// 15/26/33 load real engines. Sampling is never skipped and thresholds are
// never loosened; docs/SLO.md stays the authority for the measured backfill.
//
// Modes:
//
//	wisp slo -state Sleeping [-seconds 5] [-leak]     state sampling
//	wisp slo -settle [-seconds 10]                    settle-row check
//
// The leak fixture allocates and holds 100MB (page-touched) before sampling,
// which must flip the memory gate red - a green run proves the sampler is
// broken. The settle fixture allocates a session-shaped buffer, releases it
// through observe.ReleaseMemory (the instrumented C11 step-4 call) and
// requires the tree memory back under the Sleeping cap within the window AND
// a FreeOSMemory counter > 0.

const (
	leakFixtureBytes = 100 << 20
	settlePeakBytes  = 120 << 20
)

// leakHold keeps the forced leak alive until process exit.
var leakHold []byte

// settleHold keeps the settle fixture alive until ReleaseMemory runs.
var settleHold []byte

type sloRun struct {
	Mode      string                `json:"mode"` // "state" | "settle"
	State     string                `json:"state"`
	Posture   string                `json:"posture"`
	StartedAt string                `json:"started_at"`
	Seconds   float64               `json:"seconds"`
	Report    *observe.StateReport  `json:"report,omitempty"`
	Settle    *observe.SettleReport `json:"settle,omitempty"`
	Pass      bool                  `json:"pass"`
}

func cmdSLOUsage() {
	fmt.Print(`wisp slo - SLO sampling driver (ticket 08)

Usage:
  wisp slo -state <Sleeping|Armed|Warm|Conversation|PanelOpen|WorkPeak> [flags]
  wisp slo -settle [flags]

Flags:
  -state string      SLO state to sample (required unless -settle)
  -settle            run the settle-row check (peak buffer, ReleaseMemory,
                     tree memory back under the Sleeping cap within 10s)
  -seconds float     sampling window in seconds (default 5)
  -interval-ms int   sampling interval in ms (default 250)
  -leak              force the leak fixture: allocate + hold 100MB (must
                     flip the gate red; used by slo-check.ps1 self-test)
  -out string        write the JSON report to this path (default stdout)

Environment: runs under WISP_ENV=test by default (no single-instance mutex,
hermetic data dir via WISP_TEST_DATA_DIR).
`)
}

func cmdSLO(args []string) int {
	fs := flag.NewFlagSet("slo", flag.ContinueOnError)
	fs.Usage = cmdSLOUsage
	state := fs.String("state", "", "SLO state to sample")
	settle := fs.Bool("settle", false, "run the settle-row check")
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
		run.Settle = runSettle(ctx, sampler, proc.NewTreeSampler(rt.Job))
		run.Pass = run.Settle.Pass
	default:
		run.Mode = "state"
		run.Posture = "skeleton"
		settleBeforeSample(*settleMS)
		if *leak {
			run.Posture = "skeleton+leak-fixture"
			holdLeakFixture()
		}
		sampler.MarkTransition("boot", run.State)
		rep, err := sampler.SampleState(ctx, observe.SLOState(run.State),
			time.Duration(*intervalMS)*time.Millisecond, time.Duration(*seconds*float64(time.Second)))
		if err != nil {
			fmt.Fprintf(os.Stderr, "wisp slo: sampling failed: %v\n", err)
			rt.Shutdown(false)
			return 2
		}
		run.Report = rep
		run.Pass = rep.Pass
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
