package observe

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

// Ticket 136 AC#12 - "this window lost readings" has to be readable from the
// settle report.
//
// The two shapes below are the acceptor's probes from
// docs/evidence/s1/136-ac8-ac9-r1-acceptance.md §3.2 (R-136-7, severity 中):
// a settle window where only one read out of five was trustworthy answered
// `samples=1 back=60 pass=true`, and one where half the reads failed answered
// `samples=2 pass=true`. Both were true statements about the reads that were
// kept; neither report said anything about the reads that were dropped, because
// SettleReport had no such field (sample_errors / last_sample_error /
// dropped_reads were all ABSENT on the wire) while the neighbouring
// StateReport had carried the first two since ticket 66 (sampler.go:164-168).
//
// What these legs nail is the disclosure, not the verdict: a partially covered
// window still passes today (changing that verdict is not this cell's job), but
// it can no longer pass silently. The counting branches in CheckSettle
// (sampler.go:488-500) are what these legs stand on; removing either counter
// turns its own leg red (evidence §2.4, mutations MD/ME/MF).
//
// Thresholds, goldens, thresholds.go and the SampleState side are untouched by
// this file, and every fixture sits under the frozen Sleeping cap (25MB) on
// purpose. The existing settle legs - sampler_test.go:243/:276 and the AC#9
// nails - are neither edited nor skipped.

// errSettleProbeRead is the distinctive failure text these legs look for in
// last_sample_error, so a report cannot satisfy them with some other reason.
var errSettleProbeRead = errors.New("settle probe: transient tree read failure")

// coverageStep is one scripted answer of the settle seam.
type coverageStep struct {
	m   TreeMetrics
	err error
}

// coverageScriptTree answers ReadTree from a script and keeps answering with
// the last step once the script runs out, so how many ticker ticks fired stays
// a timing question while WHAT the window saw does not.
type coverageScriptTree struct {
	steps []coverageStep
	i     int
	reads int
}

func (t *coverageScriptTree) ReadTree() (TreeMetrics, error) {
	t.reads++
	if len(t.steps) == 0 {
		// Dogfooding AC#13: an unscripted fake is a broken fixture, not a
		// measurement. Refuse it here instead of indexing into an empty
		// script and taking the test binary down with this case.
		return TreeMetrics{}, errors.New("observe test seam: coverageScriptTree has no scripted reads")
	}
	step := t.steps[t.i]
	if t.i < len(t.steps)-1 {
		t.i++
	}
	return step.m, step.err
}

// trustworthyRead is a live-tree reading the sampler may record: 4MB, well
// under the frozen Sleeping cap, so nothing but coverage can be at issue.
func trustworthyRead() coverageStep {
	return coverageStep{m: TreeMetrics{PIDs: 1, PrivateWorkingSetBytes: 4 << 20, CommitBytes: 5 << 20, Handles: 10}}
}

// zeroFootprint is a live tree reporting no private working set: untrustworthy
// (see SampleState) and dropped.
func zeroFootprint() coverageStep {
	return coverageStep{m: TreeMetrics{PIDs: 1, PrivateWorkingSetBytes: 0}}
}

// failedRead is a tree read that simply did not produce metrics.
func failedRead() coverageStep {
	return coverageStep{m: TreeMetrics{}, err: errSettleProbeRead}
}

// settleSUT runs one CheckSettle over the given script with the release leg
// satisfied, and returns the report plus how many reads the seam actually took.
func settleSUT(t *testing.T, steps []coverageStep) (*SettleReport, int) {
	t.Helper()
	ft := &coverageScriptTree{steps: steps}
	s := NewSampler(ft, NewRegistry())
	// The FreeOSMemory counter and released=true must stay satisfied, so the
	// only thing these legs can ever be red about is coverage disclosure.
	prev := debugFreeOSMemory
	debugFreeOSMemory = func() {}
	t.Cleanup(func() { debugFreeOSMemory = prev })
	ReleaseMemory()
	rep, err := s.CheckSettle(context.Background(), SLOSleeping, 100*time.Millisecond, 10*time.Millisecond, true, TreeMetrics{}, 100<<20)
	if err != nil {
		t.Fatal(err)
	}
	return rep, ft.reads
}

// settleWire marshals the report the way `wisp slo -settle` writes it, and
// fails the calling case if that ever stops parsing.
func settleWire(t *testing.T, rep *SettleReport) map[string]any {
	t.Helper()
	raw, err := json.Marshal(rep)
	if err != nil {
		t.Fatal(err)
	}
	var wire map[string]any
	if err := json.Unmarshal(raw, &wire); err != nil {
		t.Fatal(err)
	}
	return wire
}

// TestCheckSettleSingleTrustworthyReadReportsItsLoss is probe ① ("整窗只 1 枚
// 可信"): one trustworthy read followed by nothing but failures.
func TestCheckSettleSingleTrustworthyReadReportsItsLoss(t *testing.T) {
	rep, reads := settleSUT(t, []coverageStep{trustworthyRead(), failedRead(), failedRead()})
	if reads < 3 {
		t.Fatalf("precondition broken: the window only took %d reads, it must lose some", reads)
	}
	// The shape the acceptor measured: one sample out of the whole window.
	if len(rep.Samples) != 1 {
		t.Fatalf("precondition broken: want exactly 1 trustworthy sample, got %d (reads=%d report=%+v)", len(rep.Samples), reads, rep)
	}
	// The disclosure this cell is about.
	if rep.SampleErrors != reads-1 {
		t.Fatalf("report counted %d dropped reads but the seam took %d reads and kept %d: %d unaccounted",
			rep.SampleErrors, reads, len(rep.Samples), reads-len(rep.Samples)-rep.SampleErrors)
	}
	if rep.SampleErrors < 2 {
		t.Fatalf("this window lost %d of %d reads and must say so, sample_errors=%d", reads-1, reads, rep.SampleErrors)
	}
	if !strings.Contains(rep.LastSampleError, "settle probe: transient tree read failure") {
		t.Fatalf("the report must carry the reason it dropped reads, got %q", rep.LastSampleError)
	}
	if !rep.Pass {
		t.Fatalf("this leg pins disclosure, not the verdict; a covered-enough window must still pass: %+v", rep)
	}

	wire := settleWire(t, rep)
	samples, isList := wire["samples"].([]any)
	if !isList {
		t.Fatalf("wire samples must be a list, got %#v", wire["samples"])
	}
	if len(samples) != 1 {
		t.Fatalf("wire samples must carry the one trustworthy read, got %d", len(samples))
	}
	if got, ok := wire["sample_errors"]; !ok {
		t.Fatal(`wire report has no "sample_errors": a partially covered window is silent again`)
	} else if got != float64(rep.SampleErrors) {
		t.Fatalf("wire sample_errors=%v but the report says %d", got, rep.SampleErrors)
	}
	if got, ok := wire["last_sample_error"]; !ok {
		t.Fatal(`wire report has no "last_sample_error" although reads were dropped`)
	} else if !strings.Contains(got.(string), "settle probe: transient tree read failure") {
		t.Fatalf("wire last_sample_error=%v", got)
	}
	// Reads that happened are either kept or counted: nothing disappears.
	if lost, ok := wire["sample_errors"].(float64); !ok {
		t.Fatalf(`wire "sample_errors" is not a number: %#v`, wire["sample_errors"])
	} else if len(samples)+int(lost) != reads {
		t.Fatalf("wire keeps %d and counts %d losses but the seam took %d reads", len(samples), int(lost), reads)
	}
}

// TestCheckSettleHalfTheReadsFailedReportsItsLoss is probe ② ("一半读数报错"):
// every other read fails, so no matter how many ticks fired, about half the
// window is missing and the report has to say so.
func TestCheckSettleHalfTheReadsFailedReportsItsLoss(t *testing.T) {
	tree := &alternatingTree{}
	s := NewSampler(tree, NewRegistry())
	prev := debugFreeOSMemory
	debugFreeOSMemory = func() {}
	t.Cleanup(func() { debugFreeOSMemory = prev })
	ReleaseMemory()
	rep, err := s.CheckSettle(context.Background(), SLOSleeping, 100*time.Millisecond, 10*time.Millisecond, true, TreeMetrics{}, 100<<20)
	if err != nil {
		t.Fatal(err)
	}
	kept, lost := len(rep.Samples), rep.SampleErrors
	if tree.reads < 4 {
		t.Fatalf("precondition broken: only %d reads taken, half-and-half needs a window to lose in", tree.reads)
	}
	if kept < 2 {
		t.Fatalf("precondition broken: the fixture kept only %d of %d reads, this leg needs trustworthy reads to compare against", kept, tree.reads)
	}
	// The property: about half the window failed, and the report has to carry
	// that instead of just the kept half.
	if lost < 2 {
		t.Fatalf("the seam lost %d of %d reads but the report says sample_errors=%d: a half-covered window must report its losses", tree.reads-kept, tree.reads, lost)
	}
	if kept+lost != tree.reads {
		t.Fatalf("kept %d + dropped %d != %d reads taken: the window is still hiding readings", kept, lost, tree.reads)
	}
	if kept > lost+1 || lost > kept+1 {
		t.Fatalf("this leg is about a half-covered window, kept=%d lost=%d reads=%d", kept, lost, tree.reads)
	}
	if !strings.Contains(rep.LastSampleError, "settle probe: transient tree read failure") {
		t.Fatalf("a dropped read must leave its reason behind, got %q", rep.LastSampleError)
	}
	if !rep.Pass {
		t.Fatalf("disclosure leg, not a verdict leg: this window still passes, report=%+v", rep)
	}

	wire := settleWire(t, rep)
	if got := wire["sample_errors"]; got != float64(lost) {
		t.Fatalf("wire must expose how many reads this window lost, got sample_errors=%v (report says %d)", got, lost)
	}
	if samples, ok := wire["samples"].([]any); !ok || len(samples) != kept {
		t.Fatalf("wire samples must show the %d kept reads, got %#v", kept, wire["samples"])
	}
}

// alternatingTree answers every other read with a failure and every other one
// with a trustworthy live-tree footprint.
type alternatingTree struct{ reads int }

func (t *alternatingTree) ReadTree() (TreeMetrics, error) {
	t.reads++
	if t.reads%2 == 0 {
		return TreeMetrics{PIDs: 1, PrivateWorkingSetBytes: 4 << 20, CommitBytes: 5 << 20, Handles: 10}, nil
	}
	return TreeMetrics{}, errSettleProbeRead
}

// TestCheckSettleZeroFootprintDropsAreCountedToo covers the OTHER drop branch
// (the fail-closed comparison AC#9 nailed): a read that returns a live tree with
// no footprint is dropped as untrustworthy and must be counted with the same
// sentence StateReport uses, or a window could still lose reads silently while
// sample_errors stays 0.
func TestCheckSettleZeroFootprintDropsAreCountedToo(t *testing.T) {
	rep, reads := settleSUT(t, []coverageStep{trustworthyRead(), zeroFootprint(), zeroFootprint()})
	if reads < 3 {
		t.Fatalf("precondition broken: the window only took %d reads", reads)
	}
	if len(rep.Samples) != 1 {
		t.Fatalf("precondition broken: want 1 recorded sample, got %d (report=%+v)", len(rep.Samples), rep)
	}
	if rep.SampleErrors != reads-1 || rep.SampleErrors < 2 {
		t.Fatalf("zero-footprint drops must be counted: sample_errors=%d reads=%d kept=1", rep.SampleErrors, reads)
	}
	if rep.LastSampleError != "read returned a zero private working set for a live tree" {
		t.Fatalf("the settle side must use the reason the sampling side uses, got %q", rep.LastSampleError)
	}
	if !rep.Pass {
		t.Fatalf("report=%+v", rep)
	}
	wire := settleWire(t, rep)
	if got, ok := wire["sample_errors"]; !ok || got != float64(reads-1) {
		t.Fatalf("wire sample_errors must count the zero-footprint drops, got %v (present=%v)", got, ok)
	}
}

// TestCheckSettleFullyMeasuredWindowReportsNoLoss is the positive control that
// keeps the three legs above from being satisfied by a report that always
// claims loss: a window whose reads all land must show sample_errors=0 and no
// last_sample_error at all.
func TestCheckSettleFullyMeasuredWindowReportsNoLoss(t *testing.T) {
	rep, reads := settleSUT(t, []coverageStep{trustworthyRead()})
	if reads < 3 {
		t.Fatalf("precondition broken: only %d reads taken", reads)
	}
	if len(rep.Samples) != reads {
		t.Fatalf("a fully trustworthy window must record every read: samples=%d reads=%d", len(rep.Samples), reads)
	}
	if rep.SampleErrors != 0 {
		t.Fatalf("nothing was dropped, but the report claims %d losses (last=%q)", rep.SampleErrors, rep.LastSampleError)
	}
	if rep.LastSampleError != "" {
		t.Fatalf("last_sample_error must stay empty when nothing failed, got %q", rep.LastSampleError)
	}
	if !rep.Pass {
		t.Fatalf("a measurable, settled, released window must pass, report=%+v", rep)
	}
	wire := settleWire(t, rep)
	if got, ok := wire["sample_errors"]; !ok || got != float64(0) {
		t.Fatalf("wire sample_errors must be present and 0 for a fully measured window, got %v (present=%v)", got, ok)
	}
	if _, present := wire["last_sample_error"]; present {
		t.Fatalf(`wire must not carry last_sample_error when nothing failed (omitempty), got %#v`, wire["last_sample_error"])
	}
}
