package observe

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// fakeTree feeds scripted metrics to the sampler (seam injection: the
// sampling logic is platform-neutral; the real Job Object read is exercised
// by the windows-only live test and by slo-check.ps1 end to end).
type fakeTree struct {
	mu      []TreeMetrics // sequence returned per read
	i       int
	err     error
	current func() TreeMetrics
}

func (f *fakeTree) ReadTree() (TreeMetrics, error) {
	if f.err != nil {
		return TreeMetrics{}, f.err
	}
	if f.current != nil {
		return f.current(), nil
	}
	if f.i >= len(f.mu) {
		return f.mu[len(f.mu)-1], nil
	}
	m := f.mu[f.i]
	f.i++
	return m, nil
}

func TestSLOStateNamesPinnedToMachineStates(t *testing.T) {
	// Sleeping/Armed/Warm/Conversation are D43 machine states; the SLO
	// vocabulary must not drift from the 20-state table (ticket 07).
	pairs := map[SLOState]statemachine.State{
		SLOSleeping:     statemachine.StateSleeping,
		SLOArmed:        statemachine.StateArmed,
		SLOWarm:         statemachine.StateWarm,
		SLOConversation: statemachine.StateConversation,
	}
	for slo, st := range pairs {
		if string(slo) != string(st) {
			t.Errorf("SLO state %q drifts from D43 state %q", slo, st)
		}
	}
	for _, s := range SLOStates {
		if !s.Valid() {
			t.Errorf("SLOStates contains invalid member %q", s)
		}
	}
	if len(SLOStates) != 6 {
		t.Fatalf("want 6 D32 states, got %d", len(SLOStates))
	}
}

func TestSampleStateAllMetricsAndVerdicts(t *testing.T) {
	// A clean skeleton tree: everything under the frozen caps.
	ft := &fakeTree{current: func() TreeMetrics {
		return TreeMetrics{
			PIDs: 1, PrivateBytes: 16 << 20,
			CPUTotalNanos: 5_000_000, GDIObjects: 5, USERObjects: 14,
			Handles: 420, Threads: 23,
		}
	}}
	reg := NewRegistry()
	s := NewSampler(ft, reg)
	rep, err := s.SampleState(context.Background(), SLOSleeping, 20*time.Millisecond, 120*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Samples) < 3 {
		t.Fatalf("expected several samples, got %d", len(rep.Samples))
	}
	if !rep.Pass {
		t.Fatalf("clean tree must pass, verdicts: %+v", rep.Verdicts)
	}
	// All seven D32 metric families present in the verdict set.
	want := map[string]bool{
		"tree_private_bytes": false, "cpu_percent_all_core": false, "gdi_objects": false,
		"user_objects": false, "handles": false, "goroutines": false, "threads": false,
		"disk_write_ops": false, "tcp_connections": false,
	}
	for _, v := range rep.Verdicts {
		want[v.Metric] = true
	}
	for m, seen := range want {
		if !seen {
			t.Errorf("metric %q missing from verdicts", m)
		}
	}
	if rep.GoroutinesMax != 0 && rep.GoroutinesMax > goroutineLimitSleeping {
		t.Fatalf("goroutine gate wrong: %d", rep.GoroutinesMax)
	}
	if rep.MemMedianBytes != 16<<20 {
		t.Fatalf("mem median wrong: %d", rep.MemMedianBytes)
	}
}

func TestSampleStateSleepingDiskWriteGateFails(t *testing.T) {
	writes := int64(0)
	ft := &fakeTree{current: func() TreeMetrics {
		writes += 7 // periodic writes between reads
		return TreeMetrics{PIDs: 1, PrivateBytes: 16 << 20, Handles: 10, WriteOps: writes}
	}}
	s := NewSampler(ft, NewRegistry())
	rep, err := s.SampleState(context.Background(), SLOSleeping, 10*time.Millisecond, 60*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if rep.Pass {
		t.Fatal("periodic disk writes must fail the Sleeping gate")
	}
	for _, v := range rep.Verdicts {
		if v.Metric == "disk_write_ops" && (v.Pass || !v.Gate) {
			t.Fatalf("disk_write_ops verdict wrong: %+v", v)
		}
	}
}

func TestSampleStateSleepingTCPGateFails(t *testing.T) {
	ft := &fakeTree{current: func() TreeMetrics {
		return TreeMetrics{PIDs: 1, PrivateBytes: 16 << 20, Handles: 10, TCPConnections: 2}
	}}
	s := NewSampler(ft, NewRegistry())
	rep, err := s.SampleState(context.Background(), SLOSleeping, 10*time.Millisecond, 40*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if rep.Pass {
		t.Fatal("long-lived TCP connections must fail the Sleeping gate")
	}
}

func TestSampleStateHandleGateUsesRulingLimit(t *testing.T) {
	// SLO.md ruling 1: <600 across all states; 640 handles = red.
	ft := &fakeTree{current: func() TreeMetrics {
		return TreeMetrics{PIDs: 1, PrivateBytes: 16 << 20, Handles: 640, GDIObjects: 5}
	}}
	s := NewSampler(ft, NewRegistry())
	rep, err := s.SampleState(context.Background(), SLOWarm, 10*time.Millisecond, 40*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if rep.Pass {
		t.Fatal("640 handles must fail the handle guardrail")
	}
}

func TestSampleStateWorkPeakMemoryIsTargetNotGate(t *testing.T) {
	ft := &fakeTree{current: func() TreeMetrics {
		return TreeMetrics{PIDs: 1, PrivateBytes: 690 << 20, Handles: 100, GDIObjects: 5}
	}}
	s := NewSampler(ft, NewRegistry())
	rep, err := s.SampleState(context.Background(), SLOWorkPeak, 10*time.Millisecond, 40*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range rep.Verdicts {
		if v.Metric == "tree_private_bytes" && v.Gate {
			t.Fatal("work-peak memory row must stay a non-gate target (D32/SLO.md)")
		}
	}
	if !rep.Pass {
		t.Fatalf("non-gate rows must not fail the report: %+v", rep.Verdicts)
	}
}

func TestSampleStateLeakFixtureFlipsRed(t *testing.T) {
	// Forced-leak fixture contract (ticket acceptance): a 100MB allocation
	// must flip the Sleeping memory gate to fail.
	ft := &fakeTree{current: func() TreeMetrics {
		return TreeMetrics{PIDs: 1, PrivateBytes: (16 + 100) << 20, Handles: 10, GDIObjects: 5}
	}}
	s := NewSampler(ft, NewRegistry())
	rep, err := s.SampleState(context.Background(), SLOSleeping, 10*time.Millisecond, 40*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if rep.Pass {
		t.Fatal("leak fixture must flip the memory gate")
	}
	for _, v := range rep.Verdicts {
		if v.Metric == "tree_private_bytes" && v.Pass {
			t.Fatalf("memory verdict must be red: %+v", v)
		}
	}
}

func TestSampleStateCPUTotalDrivenMean(t *testing.T) {
	// 0.5 core-second per second on a 12-core box = 4.166% all-core.
	cpu := int64(0)
	ft := &fakeTree{current: func() TreeMetrics {
		cpu += int64(0.5 * float64(20*time.Millisecond) / float64(time.Nanosecond))
		return TreeMetrics{PIDs: 1, PrivateBytes: 16 << 20, Handles: 10, CPUTotalNanos: cpu}
	}}
	s := NewSampler(ft, NewRegistry())
	rep, err := s.SampleState(context.Background(), SLOWarm, 20*time.Millisecond, 100*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if rep.CPUMeanPercent <= 0 || rep.CPUMeanPercent > 10 {
		t.Fatalf("cpu mean implausible: %.3f%%", rep.CPUMeanPercent)
	}
}

func TestSampleStateUnknownStateAndReaderError(t *testing.T) {
	s := NewSampler(&fakeTree{}, NewRegistry())
	if _, err := s.SampleState(context.Background(), SLOState("Nope"), time.Millisecond, time.Millisecond); err == nil {
		t.Fatal("unknown state must error")
	}
	s2 := NewSampler(&fakeTree{err: errors.New("boom")}, NewRegistry())
	if _, err := s2.SampleState(context.Background(), SLOSleeping, time.Millisecond, time.Millisecond); err == nil {
		t.Fatal("reader error must propagate")
	}
}

func TestMarkTransitionTimestamps(t *testing.T) {
	s := NewSampler(&fakeTree{}, NewRegistry())
	tr := s.MarkTransition("boot", string(SLOSleeping))
	if tr.From != "boot" || tr.To != "Sleeping" || tr.At == "" {
		t.Fatalf("transition record wrong: %+v", tr)
	}
	if got := s.Transitions(); len(got) != 1 || got[0] != tr {
		t.Fatalf("transitions history wrong: %+v", got)
	}
}

func TestCheckSettleVerifiesReleaseCounter(t *testing.T) {
	ft := &fakeTree{current: func() TreeMetrics {
		return TreeMetrics{PIDs: 1, PrivateBytes: 16 << 20, Handles: 10}
	}}
	s := NewSampler(ft, NewRegistry())

	// released=false: the gate must fail BY CONTRACT even if memory is low -
	// the row exists to prove the release path ran (SPEC-10 §3).
	rep, err := s.CheckSettle(context.Background(), SLOSleeping, 100*time.Millisecond, 20*time.Millisecond, false, TreeMetrics{}, 100<<20)
	if err != nil {
		t.Fatal(err)
	}
	if rep.Pass {
		t.Fatal("settle without a FreeOSMemory invocation must fail")
	}

	// released=true + memory back under cap: pass (inject a release call).
	prev := debugFreeOSMemory
	debugFreeOSMemory = func() {} // keep the test off the real GC release
	t.Cleanup(func() { debugFreeOSMemory = prev })
	ReleaseMemory()
	rep, err = s.CheckSettle(context.Background(), SLOSleeping, 100*time.Millisecond, 20*time.Millisecond, true, TreeMetrics{}, 100<<20)
	if err != nil {
		t.Fatal(err)
	}
	if !rep.Pass {
		t.Fatalf("settle should pass: %+v", rep)
	}
	if rep.FreeOSMemoryCount == 0 {
		t.Fatal("counter must be >0 after ReleaseMemory")
	}
}

func TestCheckSettleNeverReachesCap(t *testing.T) {
	ft := &fakeTree{current: func() TreeMetrics {
		return TreeMetrics{PIDs: 1, PrivateBytes: 120 << 20, Handles: 10}
	}}
	s := NewSampler(ft, NewRegistry())
	prev := debugFreeOSMemory
	debugFreeOSMemory = func() {}
	t.Cleanup(func() { debugFreeOSMemory = prev })
	ReleaseMemory()
	rep, err := s.CheckSettle(context.Background(), SLOSleeping, 80*time.Millisecond, 20*time.Millisecond, true, TreeMetrics{}, 200<<20)
	if err != nil {
		t.Fatal(err)
	}
	if rep.Pass || rep.BackWithinCapMS >= 0 {
		t.Fatalf("settle must fail when the cap is never reached: %+v", rep)
	}
}

// TestSamplerGoroutineAccountingFollowsRegistry pins that the sampler's
// goroutine metric is the D38b registry count (the gate's source of truth),
// not runtime.NumGoroutine (runtime workers would poison the baseline).
func TestSamplerGoroutineAccountingFollowsRegistry(t *testing.T) {
	reg := NewRegistry()
	ft := &fakeTree{current: func() TreeMetrics { return TreeMetrics{PIDs: 1, PrivateBytes: 1 << 20} }}
	s := NewSampler(ft, reg)
	rep, err := s.SampleState(context.Background(), SLOSleeping, 10*time.Millisecond, 30*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if rep.GoroutinesMax != 0 {
		t.Fatalf("empty registry must report 0, got %d", rep.GoroutinesMax)
	}
	if rep.Samples[0].RuntimeGoroutines < 1 {
		t.Fatal("runtime goroutine companion field must be recorded")
	}
}

func TestLiveRegistryBaselineWithinSleepingGate(t *testing.T) {
	// The D38b resident baseline is 6; the process registry (log-flusher
	// included once wired) must stay within the Sleeping gate even with the
	// pipeline running. This is the guard the SLO driver exercises for real.
	reg := Default
	if reg.Count() > goroutineLimitSleeping {
		t.Fatalf("live process registry exceeds the Sleeping gate: %d (runtime: %d)",
			reg.Count(), runtime.NumGoroutine())
	}
}

func TestThresholdTableCoversAllStates(t *testing.T) {
	for _, st := range SLOStates {
		capBytes := stateMemCap(st)
		if capBytes <= 0 {
			t.Errorf("state %s has no memory cap", st)
		}
		if stateCPULimit(st) <= 0 {
			t.Errorf("state %s has no cpu limit", st)
		}
		rep := StateReport{State: st}
		rep.Verdicts = buildVerdicts(st, rep)
		if len(rep.Verdicts) == 0 {
			t.Errorf("state %s produced no verdicts", st)
		}
	}
}
