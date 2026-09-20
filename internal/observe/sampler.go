package observe

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"slices"
	"sort"
	"sync/atomic"
	"time"
)

// SLO sampler (ticket 08; D32 16.3.2, D42#10, SPEC-10 §3).
//
// One sample = the seven D32 metrics over the whole PROCESS TREE (C30
// JobScope, not single-process RSS - the WebView2 children are the classic
// blind spot):
//
//	1. tree private bytes   2. CPU (all-core mean)   3. GDI objects
//	4. USER objects         5. handles               6. goroutines
//	7. threads
//
// plus the two Sleeping hard-constraint inputs (periodic disk writes, live
// TCP connections). Thresholds and gate/target classification live in
// thresholds.go and are FROZEN: nothing here may loosen or skip them
// (D22 run-away mode 6).
//
// The sampler is platform-neutral; the process-tree reader is injected
// (WindowsTreeReader on Windows over proc.JobScope). Goroutine counts come
// from the D38b registry (the product source of truth), with
// runtime.NumGoroutine recorded alongside for drift diagnosis.

// SLOState is one of the six D32 16.3.2 sampling states. Sleeping/Armed/
// Warm/Conversation are D43 machine states (names pinned by test to
// internal/statemachine); PanelOpen and WorkPeak are SLO postures over the
// machine, not machine states themselves.
type SLOState string

const (
	SLOSleeping     SLOState = "Sleeping"
	SLOArmed        SLOState = "Armed"
	SLOWarm         SLOState = "Warm"
	SLOConversation SLOState = "Conversation"
	SLOPanelOpen    SLOState = "PanelOpen"
	SLOWorkPeak     SLOState = "WorkPeak"
)

// SLOStates lists the six frozen sampling states in D32 table order.
var SLOStates = []SLOState{SLOSleeping, SLOArmed, SLOWarm, SLOConversation, SLOPanelOpen, SLOWorkPeak}

// Valid reports whether s is one of the six states.
func (s SLOState) Valid() bool { return slices.Contains(SLOStates, s) }

// TreeMetrics is one raw read of the process tree. CPU and write counters
// are cumulative since process start; the sampler derives window deltas.
//
// Memory units (docs/SLO.md §7): the frozen thresholds are private working
// set numbers (Task Manager "Memory (private working set)"), so the gate
// metric is PrivateWorkingSetBytes; CommitBytes (C30 PrivateUsage, always
// >= private WS) is recorded alongside for diagnostics.
type TreeMetrics struct {
	// PIDs is the number of processes in the tree.
	PIDs int
	// PrivateWorkingSetBytes is the tree private working set - the D32 gate
	// metric in the units docs/SLO.md was measured in.
	PrivateWorkingSetBytes int64
	// CommitBytes is the tree private commit charge (C30 PrivateUsage sum).
	CommitBytes int64
	// CPUTotalNanos is the summed user+kernel CPU time of the tree since
	// process start (100ns units merged into nanos by the reader).
	CPUTotalNanos int64
	// GDIObjects / USERObjects are the summed GDI/USER object counts.
	GDIObjects, USERObjects int
	// Handles is the summed process handle count.
	Handles int
	// Threads is the summed thread count.
	Threads int
	// WriteOps is the summed write-operation count of the tree (excluding
	// the sampling process itself - the observer must not fail its own
	// measurement; the self counter is reported alongside).
	WriteOps, SelfWriteOps int64
	// TCPConnections counts TCP connections owned by tree processes.
	TCPConnections int
}

// TreeReader reads the process tree (C30 JobScope on Windows).
type TreeReader interface {
	ReadTree() (TreeMetrics, error)
}

// ErrTreeUnsupported is returned by the placeholder reader on platforms
// without the Job Object reader; sampling logic stays testable there via
// injected fakes.
var ErrTreeUnsupported = errors.New("observe: process-tree metrics unavailable on this platform")

// Sample is one derived sample (JSON shape of the slo output).
type Sample struct {
	At                string  `json:"at"`                 // wall RFC3339 UTC - persisted record
	TreePrivateBytes  int64   `json:"tree_private_bytes"` // private working set (SLO.md §7 gate units)
	TreeCommitBytes   int64   `json:"tree_commit_bytes"`  // C30 commit charge, recorded
	CPUPercent        float64 `json:"cpu_percent"`        // all-core mean over the interval
	GDIObjects        int     `json:"gdi_objects"`
	USERObjects       int     `json:"user_objects"`
	Handles           int     `json:"handles"`
	Goroutines        int     `json:"goroutines"`         // D38b registry count
	RuntimeGoroutines int     `json:"runtime_goroutines"` // runtime.NumGoroutine
	Threads           int     `json:"threads"`
	TreeProcesses     int     `json:"tree_processes"`
	WriteOpsDelta     int64   `json:"write_ops_delta"`      // tree minus self, interval delta
	SelfWriteOpsDelta int64   `json:"self_write_ops_delta"` // observer itself, interval delta
	TCPConnections    int     `json:"tcp_connections"`
}

// Transition is one recorded state transition (wall timestamp, persisted).
type Transition struct {
	From string `json:"from"`
	To   string `json:"to"`
	At   string `json:"at"` // wall RFC3339 UTC
}

// Verdict is one metric's gate evaluation for a state window.
type Verdict struct {
	Metric   string `json:"metric"`
	Measured string `json:"measured"`
	Limit    string `json:"limit"`
	Pass     bool   `json:"pass"`
	// Gate false = the row is a recorded target/observation, not an
	// acceptance gate (e.g. the 700MB work peak before S3/S5 backfill,
	// D32: "target, not acceptance"). A green non-gate row never counts as
	// an SLO pass; a red one never blocks.
	Gate bool   `json:"gate"`
	Note string `json:"note,omitempty"`
	// ObserverCost true = the row was measured from inside the measured
	// process tree, so the number contains the sampler's own cost. It stays
	// in the report and stays compared against the frozen limit, but it is
	// never an acceptance gate (ticket 66 ruling 2 / registry A15).
	ObserverCost bool `json:"observer_cost,omitempty"`
}

// Measurement bases (ticket 66, registry A15). The basis says where the
// observer sits relative to what it measures - it never changes a threshold,
// only which number is allowed to gate.
const (
	// BasisOutOfTree: the reader measures another process (parent reads the
	// subject pid). This is the D32 product-side basis for the CPU row.
	BasisOutOfTree = "out-of-tree"
	// BasisInTree: the reader runs inside the measured tree, so its own
	// ReadTree burn lands in the measured CPU (docs/SLO.md appendix B.2:
	// ~1.3ms of CPU per read, i.e. 0.52% all-core at a 250ms cadence - above
	// the <=0.5% gate on its own). The CPU row stays in the report and is
	// marked observer_cost; it must not gate (orchestrator ruling 2).
	BasisInTree = "in-tree"
)

// StateReport aggregates one sampling window over one state.
type StateReport struct {
	State        SLOState `json:"state"`
	StartedAt    string   `json:"started_at"`
	DurationSec  float64  `json:"duration_sec"`
	IntervalSec  float64  `json:"interval_sec"`
	Samples      []Sample `json:"samples"`
	SampleErrors int      `json:"sample_errors"`
	// LastSampleError keeps WHY the most recent read was dropped. A bare count
	// let an instrument lose samples without saying what it lost (ticket 66:
	// the whole point is that this instrument stops hiding things).
	LastSampleError string `json:"last_sample_error,omitempty"`

	// Basis is BasisOutOfTree or BasisInTree (ticket 66).
	Basis string `json:"measurement_basis"`
	// ObserverCost true means this whole report was produced from inside the
	// measured process tree: the CPU row is the observer's own cost and is
	// recorded, never gated.
	ObserverCost bool `json:"observer_cost,omitempty"`

	MemMedianBytes       int64   `json:"mem_median_bytes"` // private working set median (gate units)
	MemMedianCommitBytes int64   `json:"mem_median_commit_bytes"`
	CPUMeanPercent       float64 `json:"cpu_mean_percent"` // exact window mean from totals
	GDIMax               int     `json:"gdi_max"`
	USERMax              int     `json:"user_max"`
	HandlesMax           int     `json:"handles_max"`
	GoroutinesMax        int     `json:"goroutines_max"`
	ThreadsMax           int     `json:"threads_max"`
	WriteOpsTotal        int64   `json:"write_ops_total"`
	SelfWriteTotal       int64   `json:"self_write_ops_total"`
	TCPMax               int     `json:"tcp_max"`

	Verdicts []Verdict `json:"verdicts"`
	Pass     bool      `json:"pass"` // all Gate verdicts pass
}

// Sampler samples the SLO metric set and records state transitions.
type Sampler struct {
	tree   TreeReader
	reg    *Registry
	now    func() time.Time // wall clock, injectable for tests
	trans  []Transition
	volume int // max transitions retained (bounded)
	// inTreeObserver true = this sampler runs inside the tree it measures, so
	// its CPU row is recorded with observer_cost instead of gating (ticket 66).
	inTreeObserver bool
}

// NewSampler builds a sampler over the given tree reader and registry. The
// default basis is out-of-tree (numbers gate). A sampler that measures the
// process it lives in MUST call MarkObserverInsideTree.
func NewSampler(tree TreeReader, reg *Registry) *Sampler {
	if reg == nil {
		reg = Default
	}
	return &Sampler{
		tree:   tree,
		reg:    reg,
		now:    time.Now,
		volume: 1024,
	}
}

// MarkObserverInsideTree declares that this sampler's reader measures the
// process tree the sampler itself runs in. The CPU row of every report it
// produces is then marked observer_cost and is not a gate (D32 thresholds are
// untouched: what moves is where the product number is read from).
func (s *Sampler) MarkObserverInsideTree() *Sampler {
	s.inTreeObserver = true
	return s
}

// MarkTransition records a state transition with a wall timestamp (D32:
// state-transition timestamps are part of the sampling evidence).
func (s *Sampler) MarkTransition(from, to string) Transition {
	tr := Transition{From: from, To: to, At: WallTimestampUTC(s.now())}
	s.trans = append(s.trans, tr)
	if n := len(s.trans); n > s.volume {
		s.trans = slices.Clone(s.trans[n-s.volume:])
	}
	return tr
}

// Transitions returns the retained transition history (copy).
func (s *Sampler) Transitions() []Transition {
	out := make([]Transition, len(s.trans))
	copy(out, s.trans)
	return out
}

// SampleState samples the tree every interval for duration and evaluates
// the frozen thresholds for st. interval <= 0 defaults to 500ms. The first
// read only sets the CPU/write baseline (no verdicts against a zero window);
// a sampling window shorter than one interval still yields real samples at
// the first tick.
func (s *Sampler) SampleState(ctx context.Context, st SLOState, interval, duration time.Duration) (*StateReport, error) {
	if !st.Valid() {
		return nil, New(ClassConfig, "observe: unknown SLO state "+string(st))
	}
	if interval <= 0 {
		interval = 500 * time.Millisecond
	}
	rep := &StateReport{
		State:        st,
		StartedAt:    WallTimestampUTC(s.now()),
		Basis:        BasisOutOfTree,
		ObserverCost: s.inTreeObserver,
	}
	if s.inTreeObserver {
		rep.Basis = BasisInTree
	}
	start := time.Now() // monotonic
	windowStart, err := s.tree.ReadTree()
	if err != nil {
		return nil, Wrap(ClassResource, err, "observe: baseline tree read")
	}
	prev := windowStart
	prevAt := start

	t := time.NewTicker(interval)
	defer t.Stop()
	deadline := start.Add(duration)
	for {
		select {
		case <-ctx.Done():
			return nil, Wrap(ClassCancelled, ctx.Err(), "observe: sampling cancelled")
		case <-t.C:
		}
		m, err := s.tree.ReadTree()
		nowAt := time.Now()
		if err != nil {
			rep.SampleErrors++
			rep.LastSampleError = "read: " + err.Error()
		} else if m.PrivateWorkingSetBytes <= 0 {
			// A zero-footprint read of a live tree is an untrustworthy
			// measurement, never a pass (fail-closed).
			rep.SampleErrors++
			rep.LastSampleError = "read returned a zero private working set for a live tree"
		} else {
			sample := s.derive(prev, prevAt, m, nowAt)
			rep.Samples = append(rep.Samples, sample)
			prev, prevAt = m, nowAt
		}
		if !nowAt.Before(deadline) {
			break
		}
	}

	// Window aggregates need a closing read for the exact CPU mean.
	windowEnd, err := s.tree.ReadTree()
	if err != nil {
		windowEnd = prev
	}
	elapsed := time.Since(start).Seconds()
	if elapsed <= 0 {
		elapsed = float64(interval) / float64(time.Second)
	}
	cores := float64(runtime.NumCPU())
	rep.DurationSec = elapsed
	rep.IntervalSec = interval.Seconds()
	rep.CPUMeanPercent = cpuPercent(windowStart.CPUTotalNanos, windowEnd.CPUTotalNanos, elapsed, cores)

	rep.MemMedianBytes = medianInt64(samplesField(rep.Samples, func(x Sample) int64 { return x.TreePrivateBytes }))
	rep.MemMedianCommitBytes = medianInt64(samplesField(rep.Samples, func(x Sample) int64 { return x.TreeCommitBytes }))
	for _, sm := range rep.Samples {
		rep.GDIMax = maxInt(rep.GDIMax, sm.GDIObjects)
		rep.USERMax = maxInt(rep.USERMax, sm.USERObjects)
		rep.HandlesMax = maxInt(rep.HandlesMax, sm.Handles)
		rep.GoroutinesMax = maxInt(rep.GoroutinesMax, sm.Goroutines)
		rep.ThreadsMax = maxInt(rep.ThreadsMax, sm.Threads)
		rep.TCPMax = maxInt(rep.TCPMax, sm.TCPConnections)
		rep.WriteOpsTotal += sm.WriteOpsDelta
		rep.SelfWriteTotal += sm.SelfWriteOpsDelta
	}
	rep.Verdicts = buildVerdicts(st, *rep)
	if len(rep.Samples) == 0 {
		// No trustworthy sample in the whole window: the report fails with
		// its own gate verdict rather than silently passing on zeros.
		rep.Verdicts = append(rep.Verdicts, Verdict{
			Metric: "sampling", Measured: fmt.Sprintf("%d valid / %d errors", len(rep.Samples), rep.SampleErrors),
			Limit: ">0 valid samples", Pass: false, Gate: true,
			Note: "fail-closed: unmeasurable windows never pass (D22 mode-6 guard)",
		})
	}
	rep.Pass = true
	for _, v := range rep.Verdicts {
		if v.Gate && !v.Pass {
			rep.Pass = false
		}
	}
	return rep, nil
}

// derive computes one sample from consecutive raw reads.
func (s *Sampler) derive(prev TreeMetrics, prevAt time.Time, m TreeMetrics, nowAt time.Time) Sample {
	elapsed := nowAt.Sub(prevAt).Seconds()
	if elapsed <= 0 {
		elapsed = float64(time.Millisecond) / float64(time.Second)
	}
	cores := float64(runtime.NumCPU())
	return Sample{
		At:                WallTimestampUTC(nowAt),
		TreePrivateBytes:  m.PrivateWorkingSetBytes,
		TreeCommitBytes:   m.CommitBytes,
		CPUPercent:        cpuPercent(prev.CPUTotalNanos, m.CPUTotalNanos, elapsed, cores),
		GDIObjects:        m.GDIObjects,
		USERObjects:       m.USERObjects,
		Handles:           m.Handles,
		Goroutines:        s.reg.Count(),
		RuntimeGoroutines: runtime.NumGoroutine(),
		Threads:           m.Threads,
		TreeProcesses:     m.PIDs,
		WriteOpsDelta:     m.WriteOps - prev.WriteOps,
		SelfWriteOpsDelta: m.SelfWriteOps - prev.SelfWriteOps,
		TCPConnections:    m.TCPConnections,
	}
}

func cpuPercent(prevTotal, total int64, elapsedSec, cores float64) float64 {
	if elapsedSec <= 0 || cores <= 0 || total <= prevTotal {
		return 0
	}
	return (float64(total-prevTotal) / 1e9) / elapsedSec / cores * 100
}

func samplesField(ss []Sample, f func(Sample) int64) []int64 {
	out := make([]int64, 0, len(ss))
	for _, s := range ss {
		out = append(out, f(s))
	}
	return out
}

func medianInt64(v []int64) int64 {
	if len(v) == 0 {
		return 0
	}
	c := slices.Clone(v)
	sort.Slice(c, func(i, j int) bool { return c[i] < c[j] })
	return c[len(c)/2]
}

func maxInt(a, b int) int {
	if b > a {
		return b
	}
	return a
}

// --- memory release instrumentation (D32 settle row, C11 step 4) ---
//
// The settle gate must VERIFY that debug.FreeOSMemory was actually invoked
// (SPEC-10 §3: "without it this necessarily fails"). Product dispose paths
// therefore call ReleaseMemory, never debug.FreeOSMemory directly, and the
// counter is the settle evidence.

var freeOSMemoryCalls atomic.Uint64

// ReleaseMemory invokes debug.FreeOSMemory through the instrumented wrapper.
func ReleaseMemory() {
	freeOSMemoryCalls.Add(1)
	debugFreeOSMemory()
}

// FreeOSMemoryCount returns how many instrumented releases happened since
// process start.
func FreeOSMemoryCount() uint64 { return freeOSMemoryCalls.Load() }

// debugFreeOSMemory is a variable for test injection.
var debugFreeOSMemory = func() { debug.FreeOSMemory() }

// SettleReport is the D32 settle-row verdict: after Settling the tree memory
// must be back under the target state's cap within 10s and FreeOSMemory must
// have been invoked.
type SettleReport struct {
	TargetState           SLOState `json:"target_state"`
	StartedAt             string   `json:"started_at"`
	PeakBytes             int64    `json:"peak_bytes"`
	CapBytes              int64    `json:"cap_bytes"`
	FreeOSMemoryCount     uint64   `json:"free_os_memory_count"`
	FreeOSMemoryRequested bool     `json:"free_os_memory_requested"`
	BackWithinCapMS       int64    `json:"back_within_cap_ms"` // -1 = never within window
	ElapsedMS             int64    `json:"elapsed_ms"`
	FinalBytes            int64    `json:"final_bytes"`
	Samples               []Sample `json:"samples"`
	Pass                  bool     `json:"pass"`
}

// CheckSettle samples until the tree memory falls to the target state's cap
// (or the within budget expires). released=true asserts the caller invoked
// ReleaseMemory; with released=false the check fails by contract - the gate
// exists to prove the release path ran.
func (s *Sampler) CheckSettle(ctx context.Context, target SLOState, within, interval time.Duration, released bool, start TreeMetrics, peak int64) (*SettleReport, error) {
	if !target.Valid() {
		return nil, New(ClassConfig, "observe: unknown SLO state "+string(target))
	}
	if interval <= 0 {
		interval = 250 * time.Millisecond
	}
	capBytes := stateMemCap(target)
	rep := &SettleReport{
		TargetState:           target,
		StartedAt:             WallTimestampUTC(s.now()),
		PeakBytes:             peak,
		CapBytes:              capBytes,
		FreeOSMemoryCount:     FreeOSMemoryCount(),
		FreeOSMemoryRequested: released,
		BackWithinCapMS:       -1,
	}
	startAt := time.Now()
	deadline := startAt.Add(within)
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, Wrap(ClassCancelled, ctx.Err(), "observe: settle check cancelled")
		case <-t.C:
		}
		m, err := s.tree.ReadTree()
		if err == nil && m.PrivateWorkingSetBytes > 0 {
			// Zero-footprint reads of a live tree are untrustworthy (see
			// SampleState) and are dropped, never recorded as progress.
			sm := Sample{
				At:               WallTimestampUTC(time.Now()),
				TreePrivateBytes: m.PrivateWorkingSetBytes,
				TreeCommitBytes:  m.CommitBytes,
				Goroutines:       s.reg.Count(),
			}
			rep.Samples = append(rep.Samples, sm)
			rep.FinalBytes = m.PrivateWorkingSetBytes
			if m.PrivateWorkingSetBytes <= capBytes && rep.BackWithinCapMS < 0 {
				rep.BackWithinCapMS = time.Since(startAt).Milliseconds()
			}
		}
		if time.Now().After(deadline) {
			break
		}
	}
	rep.ElapsedMS = time.Since(startAt).Milliseconds()
	rep.FreeOSMemoryCount = FreeOSMemoryCount()
	backInTime := rep.BackWithinCapMS >= 0 && rep.BackWithinCapMS <= within.Milliseconds()
	memOK := rep.FinalBytes <= capBytes
	releaseOK := released && rep.FreeOSMemoryCount > 0
	rep.Pass = memOK && backInTime && releaseOK
	return rep, nil
}
