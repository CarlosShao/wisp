package observe

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// Ticket 136 AC#1 — the zero-sample fail-closed guard at the end of
// SampleState is the load-bearing wall under "a report exists, therefore
// there are numbers in it". Before this file it had no case of its own:
// deleting the guard left the whole package green (54 PASS, rc=0), which is
// the reason this ticket was filed (docs/evidence/s1/134-adversarial-acceptance.md
// §8.2, R-134-6).
//
// What is nailed here, in the direction the guard moves:
//
//	zero trustworthy samples => the report may never claim a pass, and if a
//	report is emitted anyway it must carry its own samples=0 marker and a
//	failing gate row saying why.
//
// Nothing in this file touches a threshold (D32's Sleeping CPU <=0.5% /
// RSS <=25MB are untouched, and every fixture sits under them on purpose, so
// the only thing that can make these reports red is the guard).
//
// Killability (the sanity legs are not mute - evidence §3):
//   - the guard's `if len(rep.Samples) == 0` block removed  -> leg A red;
//   - the zero-footprint discard branch made never-discard  -> leg A red on
//     its own precondition (the fixture stops being a zero-sample window);
//   - the same branch made always-discard                   -> leg B red (the
//     trustworthy window stops producing samples and stops passing).

// zeroFootprintTree answers every read successfully with a zero private
// working set: a live tree with no footprint is not a measurement, so each
// read is dropped and the window ends with zero samples.
func zeroFootprintTree(reads *int) *fakeTree {
	return &fakeTree{current: func() TreeMetrics {
		if reads != nil {
			*reads++
		}
		return TreeMetrics{PIDs: 1, PrivateWorkingSetBytes: 0}
	}}
}

// TestSampleStateZeroSampleWindowFailsClosed is the AC#1 nail: an unmeasurable
// window must fail closed, never pass on zeros.
func TestSampleStateZeroSampleWindowFailsClosed(t *testing.T) {
	var reads int
	s := NewSampler(zeroFootprintTree(&reads), NewRegistry())
	rep, err := s.SampleState(context.Background(), SLOSleeping, 10*time.Millisecond, 50*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	// Precondition leg: the window really holds zero trustworthy samples, and
	// the sampler says so instead of hiding it. (Mutating the discard branch
	// so zero-footprint reads are accepted turns this red - evidence §3.)
	if len(rep.Samples) != 0 {
		t.Fatalf("precondition broken: unmeasurable window produced %d samples", len(rep.Samples))
	}
	// AC#15 census note: this one floor is NOT behind awaitStateReads, because
	// it cannot come out of scheduling. SampleState reads the tree before the
	// ticker loop starts (sampler.go:269) and again to close the window
	// (:306), so any window that returns at all has taken at least 2 reads;
	// the tick-dependent floor of this leg is the sample count above, which
	// leg B (TestSampleStateTrustworthyWindowNotMarkedUnmeasurable) is now
	// waited for. Wrapping this guard in a wait would be a wait that can never
	// be the one that saves you.
	if reads == 0 {
		t.Fatal("precondition broken: the tree was never read")
	}
	if rep.SampleErrors == 0 {
		t.Fatal("every dropped read must be counted in sample_errors")
	}
	if !strings.Contains(rep.LastSampleError, "zero private working set") {
		t.Fatalf("the report must say WHY it has no samples, got %q", rep.LastSampleError)
	}

	// The nail: the unmeasurable window is judged unusable.
	if rep.Pass {
		t.Fatalf("zero-sample window must never pass, got pass=true verdicts=%+v", rep.Verdicts)
	}
	var guard *Verdict
	for i := range rep.Verdicts {
		if rep.Verdicts[i].Metric == "sampling" {
			guard = &rep.Verdicts[i]
		}
	}
	if guard == nil {
		t.Fatal("zero-sample report carries no `sampling` verdict row: the window is unmeasurable and nothing says it")
	}
	if !guard.Gate {
		t.Fatal("the sampling row must be a gate, not a recorded note")
	}
	if guard.Pass {
		t.Fatal("the sampling row must fail on a zero-sample window")
	}
	if !strings.HasPrefix(guard.Measured, "0 valid") {
		t.Fatalf("the sampling row must state the sample count, measured=%q", guard.Measured)
	}

	// Attribution: the sampling row is the ONLY thing failing here. Every
	// fixture number sits under the frozen caps, so if this loop finds a red
	// gate other than `sampling`, leg A stopped proving the guard.
	for _, v := range rep.Verdicts {
		if v.Metric == "sampling" {
			continue
		}
		if v.Gate && !v.Pass {
			t.Fatalf("a non-sampling gate is red (%+v): the zero-sample verdict is no longer isolated", v)
		}
	}

	// The other half of the AC wording: if a report is emitted anyway, it
	// carries its samples=0 marker to whoever reads the JSON.
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
	if raw0, isList := wire["samples"].([]any); isList && len(raw0) != 0 {
		t.Fatalf("wire report must carry samples=0, got %v", wire["samples"])
	} else if !isList && wire["samples"] != nil {
		t.Fatalf("wire samples field is neither empty nor a list: %#v", wire["samples"])
	}
	if wire["pass"] != false {
		t.Fatalf("wire report with zero samples must be judged unusable (pass=false), got %v", wire["pass"])
	}
	if !strings.Contains(string(raw), `"metric":"sampling"`) || !strings.Contains(string(raw), `"gate":true`) {
		t.Fatalf("wire report lost the failing gate row: %s", raw)
	}
}

// TestSampleStateTrustworthyWindowNotMarkedUnmeasurable is leg B, the positive
// control leg A is measured against: the same sampler on the same state, on
// reads the sampler can trust, must produce samples, must not grow a sampling
// row, and must pass. Without it, leg A would stay green for a guard that
// fires on every window (a guard that cannot tell a good report from a bad
// one is not a gate). Mutating the discard branch so that every read is
// dropped turns this red - evidence §3.
func TestSampleStateTrustworthyWindowNotMarkedUnmeasurable(t *testing.T) {
	// AC#15: samples only ever come from the ticker loop, so this floor is
	// waited for (3 reads = the 2 SampleState takes outside the loop + 1 tick).
	win := awaitStateReads(t, 3, func() (*StateReport, int) {
		reads := 0
		ft := &fakeTree{current: func() TreeMetrics {
			reads++
			return TreeMetrics{PIDs: 1, PrivateWorkingSetBytes: 16 << 20, GDIObjects: 5, Handles: 420, Threads: 23}
		}}
		s := NewSampler(ft, NewRegistry())
		rep, err := s.SampleState(context.Background(), SLOSleeping, 10*time.Millisecond, 60*time.Millisecond)
		if err != nil {
			t.Fatal(err)
		}
		return rep, reads
	})
	rep := win.rep
	if len(rep.Samples) == 0 {
		t.Fatal("trustworthy reads must be sampled")
	}
	if rep.SampleErrors != 0 {
		t.Fatalf("trustworthy reads must not be dropped: %d errors (%s)", rep.SampleErrors, rep.LastSampleError)
	}
	for _, v := range rep.Verdicts {
		if v.Metric == "sampling" {
			t.Fatalf("the sampling guard fired on a measurable window: %+v", v)
		}
	}
	if !rep.Pass {
		t.Fatalf("a clean measurable window must pass, verdicts=%+v", rep.Verdicts)
	}
}
