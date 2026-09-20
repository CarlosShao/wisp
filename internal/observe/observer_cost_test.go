package observe

import (
	"context"
	"testing"
	"time"
)

// Ticket 66 / registry A15: the CPU row's measurement BASIS decides whether it
// gates, and nothing else. The frozen D32 limit itself is untouched - these
// tests pin both halves of that sentence, because "move the gate out of the
// tree" and "delete the gate" look identical in a report that only carries a
// boolean.

func verdictFor(t *testing.T, vs []Verdict, metric string) Verdict {
	t.Helper()
	for _, v := range vs {
		if v.Metric == metric {
			return v
		}
	}
	t.Fatalf("no %q verdict in %d rows", metric, len(vs))
	return Verdict{}
}

func TestInTreeCPURowRecordedNeverGated(t *testing.T) {
	rep := StateReport{
		State: SLOSleeping, CPUMeanPercent: 0.66, ObserverCost: true, Basis: BasisInTree,
		MemMedianBytes: 9 << 20, GDIMax: 4, USERMax: 8, HandlesMax: 360,
		GoroutinesMax: 1, ThreadsMax: 20, Samples: []Sample{{TreePrivateBytes: 9 << 20}},
	}
	cpu := verdictFor(t, buildVerdicts(SLOSleeping, rep), "cpu_percent_all_core")
	if cpu.Gate {
		t.Error("in-tree CPU row must not be a gate (ruling 2: a sampler inside its own tree cannot gate)")
	}
	if !cpu.ObserverCost {
		t.Error("in-tree CPU row must carry observer_cost so the report says why it does not gate")
	}
	if cpu.Pass {
		t.Errorf("0.66%% must still read FAIL against <=0.5%%: %+v", cpu)
	}
	if cpu.Limit != "<=0.5%" {
		t.Errorf("in-tree CPU limit drifted from the frozen D32 row: %q", cpu.Limit)
	}
}

func TestOutOfTreeCPURowStillGatesAtSameLimit(t *testing.T) {
	in := StateReport{
		State: SLOSleeping, CPUMeanPercent: 0.66, ObserverCost: true, Basis: BasisInTree,
		MemMedianBytes: 9 << 20, Samples: []Sample{{TreePrivateBytes: 9 << 20}},
	}
	out := in
	out.ObserverCost = false
	out.Basis = BasisOutOfTree

	inCPU := verdictFor(t, buildVerdicts(SLOSleeping, in), "cpu_percent_all_core")
	outCPU := verdictFor(t, buildVerdicts(SLOSleeping, out), "cpu_percent_all_core")
	if !outCPU.Gate {
		t.Error("out-of-tree CPU row is the D32 acceptance gate and must gate")
	}
	if outCPU.ObserverCost {
		t.Error("out-of-tree row must not be marked observer_cost")
	}
	if inCPU.Limit != outCPU.Limit {
		t.Errorf("the same limit must appear on both rows: %q vs %q", inCPU.Limit, outCPU.Limit)
	}
	if inCPU.Pass != outCPU.Pass {
		t.Errorf("both rows must judge the same number the same way: %v vs %v", inCPU.Pass, outCPU.Pass)
	}
	// A red out-of-tree CPU row fails the window: the demotion must not have
	// moved the gate, only the location.
	rep, err := NewSampler(&fakeTree{current: func() TreeMetrics {
		return TreeMetrics{PIDs: 1, PrivateWorkingSetBytes: 9 << 20, Handles: 10, CPUTotalNanos: 0}
	}}, NewRegistry()).MarkObserverInsideTree().
		SampleState(context.Background(), SLOSleeping, 20*time.Millisecond, 60*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if rep.Basis != BasisInTree || !rep.ObserverCost {
		t.Errorf("sampler marked in-tree still reports basis %q observer_cost=%v", rep.Basis, rep.ObserverCost)
	}
	if v := verdictFor(t, rep.Verdicts, "cpu_percent_all_core"); v.Gate {
		t.Error("SampleState must propagate the in-tree basis into the report's CPU verdict")
	}
}

func TestInTreeDemotionTouchesOnlyTheCPURow(t *testing.T) {
	base := StateReport{
		State: SLOSleeping, CPUMeanPercent: 0.66, MemMedianBytes: 9 << 20,
		GDIMax: 4, USERMax: 8, HandlesMax: 360, GoroutinesMax: 1, ThreadsMax: 20,
		Samples: []Sample{{TreePrivateBytes: 9 << 20}}, Basis: BasisOutOfTree,
	}
	inTree := base
	inTree.ObserverCost = true
	inTree.Basis = BasisInTree

	vsIn := buildVerdicts(SLOSleeping, inTree)
	vsOut := buildVerdicts(SLOSleeping, base)
	if len(vsIn) != len(vsOut) {
		t.Fatalf("row count differs in-tree/out-of-tree: %d vs %d", len(vsIn), len(vsOut))
	}
	for i := range vsIn {
		if vsIn[i].Metric != vsOut[i].Metric {
			t.Fatalf("row %d is %q out-of-tree but %q in-tree", i, vsOut[i].Metric, vsIn[i].Metric)
		}
		if vsIn[i].Metric == "cpu_percent_all_core" {
			continue
		}
		if vsIn[i].Gate != vsOut[i].Gate {
			t.Errorf("%q gate changed with the basis (%v -> %v): only the CPU row may move",
				vsIn[i].Metric, vsOut[i].Gate, vsIn[i].Gate)
		}
	}
	// An in-tree report whose write-op gate fails must still fail: demoting
	// the CPU row must not quietly un-gate the window.
	dirty := inTree
	dirty.WriteOpsTotal = 7
	if v := verdictFor(t, buildVerdicts(SLOSleeping, dirty), "disk_write_ops"); !v.Gate || v.Pass {
		t.Errorf("disk_write_ops gate must stay a gate in-tree: %+v", v)
	}
}
