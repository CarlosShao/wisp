package observe

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

// Ticket 136 AC#14 - the settle report has to say, in its own rows, which gate
// it checked and whether it held.
//
// Before this file the settle side reported two numbers about coverage
// (sample_errors / last_sample_error, added by AC#12) and one total boolean
// (pass). Nothing in the report stated that "1 read out of 20" and "20 reads
// out of 20" are different verdicts, so a reader who did not already know to
// look at sample_errors could not tell them apart - the same "one step exists
// but never produces a conclusion" shape ticket 136 was filed on. buildSettleVerdicts
// (sampler.go) is the row that ends that; these legs are its nails.
//
// Scope discipline, since two other cells share this package:
//   - the StateReport side is FROZEN by the ticket (issue 136 :278-279), so
//     nothing here edits or re-derives buildVerdicts' table rows; leg
//     TestStateReportVerdictBuilderStaysSinglePurpose pins that the two
//     builders stayed two functions instead of one two-purpose one;
//   - AC#15 owns the flake account of sampler_settle_coverage_136_test.go, so
//     this file carries its OWN fixtures (gateStep/gateScriptTree) and never
//     touches theirs;
//   - no threshold, no golden, nothing under thresholds.go or cmd/wisp.
//
// Every window below is deliberately generous (200ms at 10ms, about 20 reads)
// and every precondition is stated as "at least", never as an exact read
// count, so these legs cannot join the "the window only took N reads" family
// AC#15 is chasing.

// errGateProbeRead is this file's own distinctive failure text, so no other
// leg's reason can satisfy these assertions by accident.
var errGateProbeRead = errors.New("settle gate probe: scripted tree read failure")

// gateStep is one scripted answer of gateScriptTree.
type gateStep struct {
	m   TreeMetrics
	err error
}

// gateTrustworthy is a live-tree read under the frozen Sleeping cap (25MB), so
// nothing but coverage can be at issue in these legs.
func gateTrustworthy() gateStep {
	return gateStep{m: TreeMetrics{PIDs: 1, PrivateWorkingSetBytes: 4 << 20, CommitBytes: 5 << 20, Handles: 10}}
}

// gateZeroFootprint is a live tree reporting no private working set: a read
// that is not a measurement, dropped by the fail-closed comparison AC#9 nailed.
func gateZeroFootprint() gateStep {
	return gateStep{m: TreeMetrics{PIDs: 1, PrivateWorkingSetBytes: 0}}
}

// gateFailed is a tree read that produced no metrics at all.
func gateFailed() gateStep {
	return gateStep{m: TreeMetrics{}, err: errGateProbeRead}
}

// gateScriptTree answers ReadTree from a script and keeps repeating the last
// step once the script runs out, so how many ticks fired stays a timing
// question while what the window saw does not. An unscripted tree is a broken
// fixture, not a measurement, and is refused through the error channel (the
// AC#13 lesson: a panic in an instrument eats the rest of the roster).
type gateScriptTree struct {
	steps []gateStep
	i     int
	reads int
}

func (t *gateScriptTree) ReadTree() (TreeMetrics, error) {
	t.reads++
	if len(t.steps) == 0 {
		return TreeMetrics{}, errors.New("observe test seam: gateScriptTree has no scripted reads")
	}
	step := t.steps[t.i]
	if t.i < len(t.steps)-1 {
		t.i++
	}
	return step.m, step.err
}

// gateAlternatingTree fails every other read, so about half the window is
// missing no matter how many ticks fired.
type gateAlternatingTree struct{ reads int }

func (t *gateAlternatingTree) ReadTree() (TreeMetrics, error) {
	t.reads++
	if t.reads%2 == 0 {
		return gateTrustworthy().m, nil
	}
	return TreeMetrics{}, errGateProbeRead
}

// gateSUT runs one CheckSettle over a script with the release leg satisfied.
func gateSUT(t *testing.T, tree TreeReader) (*SettleReport, int) {
	t.Helper()
	s := NewSampler(tree, NewRegistry())
	prev := debugFreeOSMemory
	debugFreeOSMemory = func() {}
	t.Cleanup(func() { debugFreeOSMemory = prev })
	ReleaseMemory()
	rep, err := s.CheckSettle(context.Background(), SLOSleeping, 200*time.Millisecond, 10*time.Millisecond, true, TreeMetrics{}, 100<<20)
	if err != nil {
		t.Fatal(err)
	}
	switch x := tree.(type) {
	case *gateScriptTree:
		return rep, x.reads
	case *gateAlternatingTree:
		return rep, x.reads
	}
	return rep, 0
}

// gateCoverageRow returns the settle report's own sampling row, or fails the
// calling case. Its absence is the AC#14 finding, so "no row" is never a pass.
func gateCoverageRow(t *testing.T, rep *SettleReport) Verdict {
	t.Helper()
	for _, v := range rep.Verdicts {
		if v.Metric == settleCoverageMetric {
			return v
		}
	}
	t.Fatalf("settle report carries no %q verdict row (AC#14): verdicts=%+v", settleCoverageMetric, rep.Verdicts)
	return Verdict{}
}

// gateRowWire marshals one verdict the way `wisp slo -settle` writes it.
func gateRowWire(t *testing.T, v Verdict) map[string]any {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var wire map[string]any
	if err := json.Unmarshal(raw, &wire); err != nil {
		t.Fatal(err)
	}
	return wire
}

// TestSettleCoverageRowExistsAndPassesWhenFullyMeasured is the positive control
// and the shape nail: a window whose reads all land carries exactly one
// coverage row, it says it passed, its measured text is formatted exactly like
// the StateReport row it is meant to be 同形 with, and the wire form shows the
// row's own gate/limit/pass fields rather than only the report's total boolean.
func TestSettleCoverageRowExistsAndPassesWhenFullyMeasured(t *testing.T) {
	rep, reads := gateSUT(t, &gateScriptTree{steps: []gateStep{gateTrustworthy()}})
	if reads < 1 || len(rep.Samples) != reads {
		t.Fatalf("precondition broken: this leg needs a window that recorded every read, got samples=%d reads=%d", len(rep.Samples), reads)
	}
	if rep.SampleErrors != 0 {
		t.Fatalf("precondition broken: nothing should have been dropped, sample_errors=%d", rep.SampleErrors)
	}
	row := gateCoverageRow(t, rep)
	if !row.Pass {
		t.Fatalf("a fully measured window must not be told it measured too little: %+v", row)
	}
	if want := fmt.Sprintf("%d valid / %d errors", len(rep.Samples), rep.SampleErrors); row.Measured != want {
		t.Fatalf("coverage row measured text must carry both counts, want %q got %q", want, row.Measured)
	}
	if !strings.Contains(row.Limit, "sample_errors") {
		t.Fatalf("the row must name the field it gates on so a red is traceable, got limit=%q", row.Limit)
	}
	if !strings.Contains(row.Note, "sample_errors=0") {
		t.Fatalf("row note must state the coverage it saw, got %q", row.Note)
	}
	if len(rep.Verdicts) != 1 {
		t.Fatalf("the settle side is entitled to exactly this one self-describing row, got %d: %+v", len(rep.Verdicts), rep.Verdicts)
	}
	if !rep.Pass {
		t.Fatalf("a measurable, settled, released window must pass, report=%+v", rep)
	}

	wire := gateRowWire(t, row)
	for _, key := range []string{"metric", "measured", "limit", "pass", "gate"} {
		if _, ok := wire[key]; !ok {
			t.Fatalf(`coverage row is missing the wire key %q that makes it self-describing: %#v`, key, wire)
		}
	}
	if wire["metric"] != settleCoverageMetric {
		t.Fatalf("coverage row metric=%v", wire["metric"])
	}
}

// TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed is AC#14 judgement 2's
// nail: half the reads fail, and the ROW - not the bare pass bit - is what has
// to say this window did not measure what it is reporting on. The red text
// names sample_errors, because a red that cannot be attributed is a red that
// gets widened.
func TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed(t *testing.T) {
	rep, reads := gateSUT(t, &gateAlternatingTree{})
	kept, lost := len(rep.Samples), rep.SampleErrors
	if kept < 1 || lost < 1 {
		t.Fatalf("precondition broken: an alternating tree must both keep and lose reads within the window, kept=%d lost=%d reads=%d", kept, lost, reads)
	}
	if kept+lost != reads {
		t.Fatalf("kept %d + dropped %d != %d reads taken: the window is still hiding readings", kept, lost, reads)
	}
	row := gateCoverageRow(t, rep)
	if row.Pass {
		t.Fatalf("a window that dropped %d of %d reads must be judged by its own row as not fully measured: %+v", lost, reads, row)
	}
	if want := fmt.Sprintf("%d errors", lost); !strings.Contains(row.Measured, want) {
		t.Fatalf("the failing row must print the dropped-read count it failed on, want %q inside measured=%q", want, row.Measured)
	}
	if !strings.Contains(row.Note, fmt.Sprintf("sample_errors=%d", lost)) {
		t.Fatalf("the failing row must name sample_errors by field name, got note=%q", row.Note)
	}
	if !strings.Contains(row.Note, "settle gate probe: scripted tree read failure") {
		t.Fatalf("the failing row must carry the reason the last read was dropped, got note=%q", row.Note)
	}
	// The verdict consequence, stated as an invariant rather than a hard-coded
	// gate flag: a report may only pass while every gate row of its own passes.
	// settleCoverageRowGates is true at HEAD and this row does not pass, so the
	// invariant is live here: rep.Pass has to be false. Recorded rather than
	// gated, the fold has nothing to veto; no leg here had to change for it.
	if row.Gate && rep.Pass {
		t.Fatalf("a failing gate row and pass=true at once: the fold at foldSettlePass is not being applied, row=%+v report=%+v", row, rep)
	}
}

// TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured pins the reason the
// row exists: the two windows must be distinguishable from the row itself, not
// from the report's total boolean. It also pins the third shape (a report with
// no samples and no errors) at the row builder, because no CheckSettle window
// can produce it - the loop always reads at least once - and the synthetic
// report in cmd/wisp/slo_windows.go:623 is the one that does (AC#10's ground).
func TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured(t *testing.T) {
	full, _ := gateSUT(t, &gateScriptTree{steps: []gateStep{gateTrustworthy()}})
	none, reads := gateSUT(t, &gateScriptTree{steps: []gateStep{gateZeroFootprint()}})
	if reads < 1 || len(none.Samples) != 0 {
		t.Fatalf("precondition broken: this leg needs a window that recorded nothing, samples=%d reads=%d", len(none.Samples), reads)
	}
	if none.SampleErrors < 1 {
		t.Fatalf("precondition broken: every dropped read must be counted, sample_errors=%d", none.SampleErrors)
	}
	fullRow, noneRow := gateCoverageRow(t, full), gateCoverageRow(t, none)
	if !fullRow.Pass {
		t.Fatalf("the fully measured control must read as measured: %+v", fullRow)
	}
	if noneRow.Pass {
		t.Fatalf("a window that recorded 0 trustworthy reads must be judged by its own row: %+v", noneRow)
	}
	if fullRow.Measured == noneRow.Measured {
		t.Fatalf(`"measured everything" and "measured nothing" print the same measured text %q: the row says nothing`, fullRow.Measured)
	}
	if !strings.HasPrefix(noneRow.Measured, "0 valid / ") {
		t.Fatalf("the unmeasured window's row must show 0 valid reads, got %q", noneRow.Measured)
	}

	// The third shape, at the builder: 0 valid / 0 errors is "never measured",
	// and must not read the same as a window that measured nothing but lost
	// reads too. It says so in its own note rather than in a bare pass=false.
	never := buildSettleVerdicts(SettleReport{})[0]
	if never.Pass {
		t.Fatalf("a report with no samples and no errors must not claim to have measured, got %+v", never)
	}
	if never.Measured != "0 valid / 0 errors" {
		t.Fatalf("never-measured shape, got measured=%q", never.Measured)
	}
	if !strings.Contains(never.Note, "0 valid samples") {
		t.Fatalf(`the never-measured note must say "measured nothing" rather than blame sample_errors, got %q`, never.Note)
	}
	if never.Measured == noneRow.Measured && never.Note == noneRow.Note {
		t.Fatalf(`"never measured" and "every read failed" are indistinguishable on the row: %+v`, never)
	}
}

// TestFoldSettlePassOnlyGateRowsVeto nails the folding rule itself, so the
// shape is covered independently of which way the shipped gate constant is set
// today (recorded rows must not veto; gate rows must). Deleting the loop body
// turns this red on its own, with no other leg to lean on.
func TestFoldSettlePassOnlyGateRowsVeto(t *testing.T) {
	cases := []struct {
		name string
		base bool
		rows []Verdict
		want bool
	}{
		{"no rows keeps the base verdict", false, nil, false},
		{"no rows keeps a true base", true, nil, true},
		{"a passing gate cannot veto", true, []Verdict{{Metric: settleCoverageMetric, Pass: true, Gate: true}}, true},
		{"a failing gate vetoes", true, []Verdict{{Metric: settleCoverageMetric, Pass: false, Gate: true}}, false},
		{"a failing recorded row cannot veto", true, []Verdict{{Metric: settleCoverageMetric, Pass: false, Gate: false}}, true},
		{"one failing gate among passing rows vetoes", true, []Verdict{
			{Metric: "other", Pass: true, Gate: true},
			{Metric: settleCoverageMetric, Pass: false, Gate: true},
			{Metric: "third", Pass: true, Gate: false},
		}, false},
		{"a base that already failed stays failed", false, []Verdict{{Metric: settleCoverageMetric, Pass: true, Gate: true}}, false},
	}
	for _, tc := range cases {
		got := foldSettlePass(tc.base, tc.rows)
		if got != tc.want {
			t.Fatalf("foldSettlePass(%v, %+v) = %v, want %v (%s)", tc.base, tc.rows, got, tc.want, tc.name)
		}
	}
}

// TestSettleReportPassNeverContradictsItsGateRows is the cross-consistency leg
// over real windows: whatever the shipped gate flag is, no report may hand back
// pass=true while carrying a failing gate row. It is the invariant that makes
// flipping settleCoverageRowGates a one-line move rather than a shape change.
func TestSettleReportPassNeverContradictsItsGateRows(t *testing.T) {
	scripts := []struct {
		name string
		tree TreeReader
	}{
		{"fully measured", &gateScriptTree{steps: []gateStep{gateTrustworthy()}}},
		{"half failed", &gateAlternatingTree{}},
		{"all reads untrustworthy", &gateScriptTree{steps: []gateStep{gateZeroFootprint()}}},
	}
	for _, sc := range scripts {
		rep, _ := gateSUT(t, sc.tree)
		for _, v := range rep.Verdicts {
			if v.Gate && !v.Pass && rep.Pass {
				t.Fatalf("%s: pass=true alongside a failing gate row %q, report=%+v", sc.name, v.Metric, rep)
			}
		}
	}
}

// TestStateReportVerdictBuilderStaysSinglePurpose guards the ticket's frozen
// boundary (issue 136 :278-279): the settle rows are a second builder, and
// buildVerdicts must keep answering only the state-table question for a
// StateReport. Merging the two into one two-purpose function breaks this at
// compile time (the argument types) and, if someone widened it instead, at the
// assertion below (a settle coverage row would appear among the table rows).
func TestStateReportVerdictBuilderStaysSinglePurpose(t *testing.T) {
	var rep StateReport
	rep.SampleErrors = 3
	rows := buildVerdicts(SLOSleeping, rep)
	if len(rows) == 0 {
		t.Fatal("precondition broken: buildVerdicts stopped producing the frozen table rows")
	}
	for _, v := range rows {
		if v.Metric == settleCoverageMetric || strings.Contains(v.Limit, "sample_errors") {
			t.Fatalf("the settle coverage row leaked into buildVerdicts, which is frozen state-report shape: %+v", v)
		}
	}
}
