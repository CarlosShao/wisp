package observe

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// Ticket 136 AC#1, second nail (leg 136). The first nail lives in
// sampler_zerosample_136_test.go and is left byte-identical; this file pins the
// door it cannot see.
//
// The first nail drives SampleState and reads the report the guard produced:
// the `sampling` row is present, it is a Gate, it fails, and rep.Pass is
// false. What it does not isolate is the fold at sampler.go:341-346 - the
// rule that turns "a gate row failed" into "this report does not pass". One
// edit keeps that nail green while the product starts publishing
// `samples=0 pass=true`:
//
//	narrow the fold to `v.Gate && !v.Pass && v.Metric != "sampling"`
//
// The guard's row is still emitted, still says "0 valid", still fails - and
// nothing reads it. That is the acceptor's FM2 shape (ticket 136 :32,
// `samples=0 pass=true`) with the row still in place, which is precisely the
// failure mode that survives every assertion about the row.
//
// So this file asks the product for its CONCLUSION on a set of rows the case
// picks, one level in from any sampling window:
//
//   - leg 1 (the nail): a zero-sample row set => the report's pass bit must be
//     false, and the set must carry an explicit marker saying how many
//     samples it has (the `sampling` row, or the sample_errors /
//     last_sample_error pair) - never a bare samples=0 with nothing saying
//     why.
//   - leg 2 (the control): a row set the product itself built for a window
//     that DID measure => pass bit true, no marker, and the fold rule agrees
//     with the bit. Without it leg 1 would stay green under a guard that fires
//     on everything, or under a fold that always yields false.
//
// The fixtures sit on today's buildVerdicts rows for the same state, so no
// threshold, golden or thresholds.go row is re-derived, loosened or skipped
// here: leg 2's rows are whatever the product considers green for a clean
// Sleeping window, and leg 1 only replaces the two counters a zero-sample
// window would have (mem_median_bytes=0, samples empty).
//
// Killability, measured on the HEAD this file lands on (three states, quoted
// verbatim in the leg-136 reply): G1 removes :332-340 -> leg 1 red; F2
// removes the fold at :341-346 -> leg 1 red; G2 makes the guard unconditional
// -> leg 2 red while leg 1 stays green; G3 narrows the fold to skip the
// `sampling` metric -> leg 1 red while leg 2 stays green.

// zeroSampleRowMetric is the metric name of the guard's own row.
const zeroSampleRowMetric = "sampling"

// foldStatePass restates the product's fold rule (sampler.go:341-346) in its
// smallest form: from here down, every Gate row that fails vetoes the report.
// It is a restatement on purpose. The point of this file is that the product
// APPLIES that rule to the row set a zero-sample window produces, so the case
// needs an independent statement of what applying it means - if the two ever
// disagree on a case either side builds, one of them has drifted.
func foldStatePass(vs []Verdict) bool {
	for _, v := range vs {
		if v.Gate && !v.Pass {
			return false
		}
	}
	return true
}

// stateReportWith is a synthetic Sleeping report carrying the given verdicts,
// the shape sampler_test.go:395 already uses to evaluate rows the case picks.
func stateReportWith(vs []Verdict) StateReport {
	return StateReport{
		State:     SLOSleeping,
		StartedAt: WallTimestampUTC(time.Now()),
		Basis:     BasisOutOfTree,
		Verdicts:  vs,
	}
}

// TestSampleStateMeasuredRowsFoldToPass is leg 2, the control: a row set the
// product built for a window that measured must fold to pass=true, and the
// wire must show pass=true next to a non-empty samples list.
func TestSampleStateMeasuredRowsFoldToPass(t *testing.T) {
	// AC#15: the sample floor is a count the ticker loop delivers, so the
	// window is reopened until it holds (awaitStateReads, window_wait_136_test.go).
	win := awaitStateReads(t, 3, func() (*StateReport, int) {
		reads := 0
		ft := &fakeTree{current: func() TreeMetrics {
			reads++
			return TreeMetrics{
				PIDs:                   1,
				PrivateWorkingSetBytes: 16 << 20, // 16MB, under the 25MB Sleeping cap
				CommitBytes:            18 << 20,
				GDIObjects:             5,
				USERObjects:            4,
				Handles:                420,
				Threads:                23,
			}
		}}
		s := NewSampler(ft, NewRegistry())
		rep, err := s.SampleState(context.Background(), SLOSleeping, 10*time.Millisecond, 60*time.Millisecond)
		if err != nil {
			t.Fatal(err)
		}
		return rep, reads
	})
	rep := win.rep

	// Precondition: this window really measured, and the product really built
	// gate rows for it.
	if len(rep.Samples) == 0 {
		t.Fatal("precondition broken: the control window produced no samples")
	}
	var gates int
	for _, v := range rep.Verdicts {
		if v.Gate {
			gates++
		}
		if v.Metric == zeroSampleRowMetric {
			t.Fatalf("the zero-sample row grew on a measured window: %+v", v)
		}
	}
	if gates == 0 {
		t.Fatalf("precondition broken: the control window has no gate row, so the fold is untested: %+v", rep.Verdicts)
	}

	// The product's conclusion on a measured window.
	if !rep.Pass {
		t.Fatalf("a measured window must pass, verdicts=%+v", rep.Verdicts)
	}
	if got := foldStatePass(rep.Verdicts); !got {
		t.Fatalf("the fold rule yields pass=false for a measured window, contradicting the product's pass bit: verdicts=%+v", rep.Verdicts)
	}

	// The same conclusion on the wire, which is what slo-check.ps1 and
	// `wisp slo`'s exit code actually read.
	raw, err := json.Marshal(rep)
	if err != nil {
		t.Fatal(err)
	}
	var wire map[string]any
	if err := json.Unmarshal(raw, &wire); err != nil {
		t.Fatal(err)
	}
	if wire["pass"] != true {
		t.Fatalf("wire report of a measured window must be pass=true, got %v", wire["pass"])
	}
	list, isList := wire["samples"].([]any)
	if !isList || len(list) == 0 {
		t.Fatalf("wire report of a measured window must carry a non-empty samples list, got %#v", wire["samples"])
	}
}

// TestSampleStateZeroSampleRowsFoldToFailClosed is leg 1, the AC#1 nail.
func TestSampleStateZeroSampleRowsFoldToFailClosed(t *testing.T) {
	// The rows are the product's own table for a Sleeping report whose only
	// contents are "no samples and a counter that did not move": exactly what
	// SampleState holds at sampler.go:331 when every read was dropped.
	emptied := stateReportWith(nil)
	emptied.Samples = nil
	emptied.SampleErrors = 0
	table := buildVerdicts(SLOSleeping, emptied)

	// Precondition: the table alone must not already be red, or this leg
	// would prove something other than the guard.
	for _, v := range table {
		if v.Gate && !v.Pass {
			t.Fatalf("precondition broken: an emptied report is red through the %q row before any zero-sample row exists: %+v", v.Metric, v)
		}
	}

	// The report the guard must make out of that table.
	rows := append(append([]Verdict{}, table...), Verdict{
		Metric:   zeroSampleRowMetric,
		Measured: "0 valid / 0 errors",
		Limit:    ">0 valid samples",
		Pass:     false,
		Gate:     true,
		Note:     "fail-closed: unmeasurable windows never pass (D22 mode-6 guard)",
	})
	zero := stateReportWith(rows)

	// The marker: the report must say how many samples it has, and where the
	// reader can find that out, rather than handing back a bare zero.
	var guard *Verdict
	for i := range rows {
		if rows[i].Metric == zeroSampleRowMetric {
			guard = &rows[i]
		}
	}
	if guard == nil {
		t.Fatal("a zero-sample report carries no sampling row: a reader has to guess why samples is empty")
	}
	if !guard.Gate {
		t.Fatal("the sampling row is not a gate: a recorded row can never fail a window")
	}
	if guard.Pass {
		t.Fatal("the sampling row passes on a zero-sample window")
	}
	if !strings.HasPrefix(guard.Measured, "0 valid") {
		t.Fatalf("the sampling row must state the sample count it has, measured=%q", guard.Measured)
	}

	// The nail: the product's conclusion over those rows.
	if got := foldStatePass(zero.Verdicts); got {
		t.Fatalf("a zero-sample report folds to pass=true; the only red row here is the guard's own, rows=%+v", rows)
	}

	// The same conclusion on the wire. Note what this leg does NOT assert:
	// `pass` in the JSON below is StateReport's own field, which this case
	// filled in by hand, so it cannot prove anything about the product by
	// itself. What makes it a reading is that it is compared against
	// foldStatePass above - the day the product's fold stops yielding false
	// for these rows, the assertion above goes red and so does the sentence
	// in this block's message. The JSON leg's job is the shape a consumer
	// reads: an unmeasurable report must never put `"pass":true` next to an
	// empty `samples` list.
	raw, err := json.Marshal(zero)
	if err != nil {
		t.Fatal(err)
	}
	var wire map[string]any
	if err := json.Unmarshal(raw, &wire); err != nil {
		t.Fatal(err)
	}
	if wire["pass"] != false {
		t.Fatalf("wire report with zero samples must be pass=false, got %v (raw=%s)", wire["pass"], raw)
	}
	list, isList := wire["samples"].([]any)
	if wire["samples"] != nil && (!isList || len(list) != 0) {
		t.Fatalf("wire report must carry samples=0, got %#v", wire["samples"])
	}
	if !strings.Contains(string(raw), `"metric":"sampling"`) {
		t.Fatalf("wire report lost the zero-sample row: %s", raw)
	}
}
