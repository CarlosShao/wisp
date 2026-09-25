package observe

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

// Ticket 136 AC#9 - the settle side of the zero-sample family.
//
// CheckSettle has the same fail-closed comparison that SampleState has: a
// zero-footprint read of a live tree is not a measurement, so sampler.go:477
// drops it instead of recording it. AC#1 (sampler_zerosample_136_test.go)
// nailed that branch for SampleState; nothing nailed it here, and ticket 136's
// acceptor measured why that matters (docs/evidence/s1/
// 136-ac1-adversarial-acceptance.md §7 R-136-1): relaxing `> 0` to `>= 0` lets
// a settle window that never took a trustworthy reading hand back
// samples=6 pass=true back_within_cap_ms=10 final_bytes=0, while the whole
// package stays green.
//
// What is nailed, in the shape AC#1 used for the neighbouring function:
//
//	zero trustworthy samples => the settle check may never pass, and the
//	report it returns says so itself (empty samples list, the -1
//	"never within cap" sentinel, pass=false on the wire).
//
// Attribution is isolated on purpose: the release counter and the memory
// comparison are both left satisfied by the fixture, so the only thing that
// can make this report red is that no trustworthy read was ever recorded.
//
// Nothing here touches a threshold, a golden or thresholds.go; the fixtures sit
// under the frozen caps by design. The two existing TestCheckSettle* cases
// (sampler_test.go:243, :276) are untouched and not skipped.

// TestCheckSettleZeroTrustworthySamplesFailsClosed is leg A, the nail.
func TestCheckSettleZeroTrustworthySamplesFailsClosed(t *testing.T) {
	// AC#15: unlike SampleState, CheckSettle has no read before the first tick
	// (sampler.go:493-501), so the "the tree was never read" floor below IS a
	// scheduling outcome and is waited for. The window is reopened whole -
	// fresh tree, fresh sampler, fresh report - so the count the wait judged and
	// the count the assertions below read are the same window's.
	win := awaitSettleReads(t, 1, func() (*SettleReport, int) {
		reads := 0
		ft := &fakeTree{current: func() TreeMetrics {
			reads++
			// A live tree (PIDs 1) reporting no private working set: every read
			// is untrustworthy and must be dropped, never recorded as progress.
			return TreeMetrics{PIDs: 1, PrivateWorkingSetBytes: 0}
		}}
		s := NewSampler(ft, NewRegistry())

		// Keep the release leg satisfied so Pass=false cannot be attributed to
		// it: the gate's own FreeOSMemory counter must be >0 and released=true.
		prev := debugFreeOSMemory
		debugFreeOSMemory = func() {}
		t.Cleanup(func() { debugFreeOSMemory = prev })
		ReleaseMemory()

		rep, err := s.CheckSettle(context.Background(), SLOSleeping, 100*time.Millisecond, 20*time.Millisecond, true, TreeMetrics{}, 100<<20)
		if err != nil {
			t.Fatal(err)
		}
		return rep, reads
	})
	rep, reads := win.rep, win.reads

	// Precondition legs: this really is a zero-trustworthy-sample window.
	// (Relaxing sampler.go:477 so zero-footprint reads are recorded turns
	// THIS red first - the report starts claiming samples it never measured.)
	if reads == 0 {
		t.Fatal("precondition broken: the tree was never read")
	}
	if len(rep.Samples) != 0 {
		t.Fatalf("precondition broken: a settle window of zero-footprint reads recorded %d samples, report=%+v", len(rep.Samples), rep)
	}

	// Attribution: both other factors of rep.Pass are satisfied here, so the
	// only one that can fail is "came back within cap", which can only be
	// reached by recording a trustworthy read.
	if !rep.FreeOSMemoryRequested {
		t.Fatal("attribution broken: the report does not carry released=true")
	}
	if rep.FreeOSMemoryCount == 0 {
		t.Fatal("attribution broken: the release counter is 0, so Pass=false would prove nothing about the sample guard")
	}
	if rep.FinalBytes > rep.CapBytes {
		t.Fatalf("attribution broken: the memory comparison alone is red (final=%d cap=%d)", rep.FinalBytes, rep.CapBytes)
	}

	// The nail: an unmeasured settle window is never a pass, and the report
	// says it measured nothing instead of inventing a timestamp.
	if rep.Pass {
		t.Fatalf("settle window with 0 trustworthy samples must never pass, got pass=true report=%+v", rep)
	}
	if rep.BackWithinCapMS >= 0 {
		t.Fatalf("report claims it settled within cap at %dms after recording 0 trustworthy reads: %+v", rep.BackWithinCapMS, rep)
	}

	// The other half of the wording: the wire form carries its own "not
	// measured" evidence to whoever reads the JSON.
	raw, err := json.Marshal(rep)
	if err != nil {
		t.Fatal(err)
	}
	var wire map[string]any
	if err := json.Unmarshal(raw, &wire); err != nil {
		t.Fatal(err)
	}
	if _, present := wire["samples"]; !present {
		t.Fatal("the wire report dropped the samples field entirely")
	}
	if list, isList := wire["samples"].([]any); isList && len(list) != 0 {
		t.Fatalf("wire report must carry samples=0 for an unmeasured window, got %v", wire["samples"])
	} else if !isList && wire["samples"] != nil {
		t.Fatalf("wire samples field is neither empty nor a list: %#v", wire["samples"])
	}
	if wire["pass"] != false {
		t.Fatalf("wire report of an unmeasured settle window must be judged unusable (pass=false), got %v", wire["pass"])
	}
	if v, ok := wire["back_within_cap_ms"]; !ok || v != float64(-1) {
		t.Fatalf("an unmeasured settle window must keep the -1 'never within cap' sentinel, got %v (present=%v)", v, ok)
	}
}

// TestCheckSettleTrustworthyReadsAreRecorded is leg B, the positive control
// leg A is measured against: the same sampler on reads it can trust must
// record them, must retire the -1 sentinel and must pass. Without it, leg A
// would stay green for a CheckSettle that records nothing ever (a guard that
// cannot tell a settled tree from an unmeasured one is not a gate). This is
// not the account of the two existing settle cases: sampler_test.go:243 pins
// the release-counter contract and :276 pins the never-reaches-cap shape,
// neither asserts that a trustworthy read lands in Samples.
func TestCheckSettleTrustworthyReadsAreRecorded(t *testing.T) {
	// AC#15: CheckSettle's samples are the tick reads that passed the
	// fail-closed comparison, so a window that took no tick has no sample to
	// check and the floor below is a scheduling outcome. Waited for; the
	// criterion is the seam's own read count, so a recording branch that stops
	// recording can never be the thing that satisfies the wait.
	win := awaitSettleReads(t, 1, func() (*SettleReport, int) {
		reads := 0
		ft := &fakeTree{current: func() TreeMetrics {
			reads++
			// Under the frozen Sleeping cap on purpose: the only reason this
			// window may not pass is a broken recording branch.
			return TreeMetrics{PIDs: 1, PrivateWorkingSetBytes: 4 << 20, Handles: 10}
		}}
		s := NewSampler(ft, NewRegistry())
		prev := debugFreeOSMemory
		debugFreeOSMemory = func() {}
		t.Cleanup(func() { debugFreeOSMemory = prev })
		ReleaseMemory()

		rep, err := s.CheckSettle(context.Background(), SLOSleeping, 100*time.Millisecond, 20*time.Millisecond, true, TreeMetrics{}, 100<<20)
		if err != nil {
			t.Fatal(err)
		}
		return rep, reads
	})
	rep := win.rep
	if len(rep.Samples) == 0 {
		t.Fatal("trustworthy settle reads must be recorded in samples")
	}
	if rep.BackWithinCapMS < 0 {
		t.Fatalf("a settled tree must report when it came back within cap, got %d with %+v", rep.BackWithinCapMS, rep)
	}
	if rep.FinalBytes <= 0 {
		t.Fatalf("a recorded read must carry its footprint, got %d", rep.FinalBytes)
	}
	if !rep.Pass {
		t.Fatalf("a measurable, settled, released window must pass, report=%+v", rep)
	}
}
