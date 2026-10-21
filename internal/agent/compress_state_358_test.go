package agent

// Ticket 358 AC#2: the tri-state outcome and the provider-reported trigger
// basis, pinned with teeth.
//
// The shape this ticket is about - "该压但压不动" - had NO reading on any
// surface before this: Need() is a bool, Ran is a bool, and the compressor's
// only trace fired on the success edge. A run whose history was over budget
// with nothing foldable left answered exactly like a run whose history never
// came near the budget. Both of those are compression passes that folded
// nothing, and they are different facts.
//
// What each case below bites, and the mutation that turns it red:
//
//   - TestCompressionStateAttemptedIsNotSkipped: delete the `entered` arm of
//     the state switch (or make the switch `if Ran {compacted} else
//     {skipped}`) and an over-budget/no-fold pass reports skipped - red.
//   - TestCompressionStateSkippedStaysSilent: delete the `default` arm and a
//     pass under budget reports attempted - red.
//   - TestCompressionAttemptedLeavesItsOwnTrace: drop the
//     `else if rep.State == CompressionAttempted` block and the shape leaves
//     no record at all - red, and the record's own attribute set is asserted
//     by name so a record saying the wrong thing is red too.
//   - TestCompressionTriggerPrefersProviderUsage: make Trigger ignore
//     `observed` (always estimate) and the provider-usage case reads
//     "estimate" - red.
//   - TestCompressionTriggerFallsBackToEstimateOnFirstTurn: make Trigger
//     compare `observed` unconditionally and the first-turn case compares 0 -
//     red (nothing would ever compress on a fresh task).
//   - TestCompressionStateRanIsTheMirrorOfCompacted: assign Ran from the loop
//     instead of from State and the two disagree - red.
//
// Every threshold is read back out of Budgets; no reference value is written.

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/llm"
)

const traceSkipMsg = "agent: history compress skipped"

// stuckAtFloorHistory is the shape ticket 358 names: over the fold budget, with
// exactly KeepRawRounds raw rounds left, so the fold loop breaks on its first
// iteration. Built from the budget, never from written numbers.
func stuckAtFloorHistory(b Budgets) []llm.Message {
	return buildRoundHistory(b.KeepRawRounds, 400)
}

func TestCompressionStateAttemptedIsNotSkipped(t *testing.T) {
	b := BudgetsFor(4096)
	c := NewCompressor(b, nil)
	hist := stuckAtFloorHistory(b)

	// Positive control: the fixture really is over budget and really has
	// nothing foldable. Without both halves the case does not test the shape.
	if !c.Need(hist) {
		t.Fatalf("setup: %d tokens does not clear the threshold %d, so no attempt would be made",
			c.TotalTokens(hist), b.HistoryCompressTokens)
	}
	if raw := len(rawRoundIndexes(groupRounds(hist))); raw != b.KeepRawRounds {
		t.Fatalf("setup: %d raw rounds, want exactly KeepRawRounds=%d or the fold could take one",
			raw, b.KeepRawRounds)
	}

	out, rep, err := c.Compress(context.Background(), hist)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	if rep.Ran {
		t.Fatalf("setup: the pass folded, so this is not the shape under test")
	}
	if len(out) != len(hist) {
		t.Fatalf("the pass did not fold but changed the message count: %d -> %d", len(hist), len(out))
	}

	if rep.State != CompressionAttempted {
		t.Errorf("State = %q, want %q: the history is over the %d-token budget and nothing was foldable, "+
			"which is an ATTEMPT that folded nothing - not a pass that never triggered",
			rep.State, CompressionAttempted, b.HistoryCompressTokens)
	}
	if !rep.State.Valid() {
		t.Errorf("State %q is not one of the three enum values", rep.State)
	}
	// The two "folded nothing" states must be distinguishable, which is the
	// whole point of AC#1: this is the assertion a companion bool could not
	// express, because it could only ever say "not skipped" or "not run".
	if rep.State == CompressionSkipped {
		t.Error("State collapsed onto skipped, so '该压但压不动' is again indistinguishable from '没触发'")
	}
}

func TestCompressionStateSkippedStaysSilent(t *testing.T) {
	recs := newTraceCapture()
	recs.assertRulerLive(t)

	b := BudgetsFor(4096)
	c := NewCompressor(b, nil, WithLogger(recs.logger()))
	hist := buildRoundHistory(2, 20) // far under the budget
	if c.Need(hist) {
		t.Fatalf("setup: %d tokens clears %d, want a pass that never triggers",
			c.TotalTokens(hist), b.HistoryCompressTokens)
	}

	_, rep, err := c.Compress(context.Background(), hist)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	if rep.State != CompressionSkipped {
		t.Errorf("State = %q, want %q for a history under the budget", rep.State, CompressionSkipped)
	}
	// Both record channels stay empty: a pass that was never an attempt has
	// nothing to report, and a record here would be the false positive ticket
	// 139's fang guards.
	if got := recs.with(traceMsg); len(got) != 0 {
		t.Errorf("%q records = %d, want none: %v", traceMsg, len(got), flattenRecordMsgs(got))
	}
	if got := recs.with(traceSkipMsg); len(got) != 0 {
		t.Errorf("%q records = %d, want none for a pass that never triggered: %v",
			traceSkipMsg, len(got), flattenRecordMsgs(got))
	}
}

// The record the over-budget/no-fold pass leaves, asserted by name. This is the
// trace AC#1's "留得下痕" actually means: the shape has to be readable off a
// disk a day later, without re-deriving it from other counters.
func TestCompressionAttemptedLeavesItsOwnTrace(t *testing.T) {
	recs := newTraceCapture()
	recs.assertRulerLive(t)

	b := BudgetsFor(4096)
	c := NewCompressor(b, nil, WithLogger(recs.logger()))
	hist := stuckAtFloorHistory(b)
	const sentinel = "PRIVACY-SENTINEL-attempted-358"
	for i := range hist {
		for j, p := range hist[i].Content {
			if tp, ok := p.(llm.TextPart); ok {
				hist[i].Content[j] = llm.TextPart{Text: tp.Text + " " + sentinel}
			}
		}
	}
	if !c.Need(hist) {
		t.Fatalf("setup: %d tokens does not clear %d", c.TotalTokens(hist), b.HistoryCompressTokens)
	}

	_, rep, err := c.Compress(context.Background(), hist)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	if rep.Ran {
		t.Fatalf("setup: the pass folded")
	}

	hits := recs.with(traceSkipMsg)
	if len(hits) != 1 {
		t.Fatalf("%q records = %d, want exactly 1 for the over-budget/no-fold shape (all: %v)",
			traceSkipMsg, len(hits), flattenRecordMsgs(recs.all()))
	}
	rec := hits[0]
	// The record must not pretend the pass compressed anything.
	if got := recs.with(traceMsg); len(got) != 0 {
		t.Errorf("the same pass also booked %q %d time(s): %v", traceMsg, len(got), flattenRecordMsgs(got))
	}

	for _, key := range []string{
		"state", "reason", "tokens_before", "tokens_after", "threshold",
		"msgs_before", "msgs_after", "compressed_msgs", "kept_raw_rounds",
		"history_changed", "trigger_source", "trigger_tokens",
	} {
		if _, ok := rec.attr[key]; !ok {
			t.Fatalf("attempted record lacks %q; present: %v", key, rec.attr)
		}
	}
	if got, ok := rec.attr["state"].(string); !ok || got != string(CompressionAttempted) {
		t.Errorf("record state = %v (%T), want %q", rec.attr["state"], rec.attr["state"], CompressionAttempted)
	}
	if got, ok := rec.attr["reason"].(string); !ok || got != skipReasonKeepRawFloor {
		t.Errorf("record reason = %v (%T), want %q: the fold loop stopped at the KeepRawRounds floor",
			rec.attr["reason"], rec.attr["reason"], skipReasonKeepRawFloor)
	}
	if got := attrInt(t, rec, "threshold"); got != int64(b.HistoryCompressTokens) {
		t.Errorf("record threshold = %d, want the budget-derived %d", got, b.HistoryCompressTokens)
	}
	if got := attrInt(t, rec, "tokens_before"); got != int64(rep.TokensBefore) {
		t.Errorf("record tokens_before = %d, want the report's %d", got, rep.TokensBefore)
	}
	if got := attrInt(t, rec, "tokens_after"); got != int64(rep.TokensAfter) {
		t.Errorf("record tokens_after = %d, want the report's %d", got, rep.TokensAfter)
	}
	// The "before" number is what tells a reader the history WAS over budget.
	if got := attrInt(t, rec, "tokens_before"); got <= int64(b.HistoryCompressTokens) {
		t.Errorf("record tokens_before = %d, which does not read as over the %d threshold",
			got, b.HistoryCompressTokens)
	}
	if got := attrInt(t, rec, "kept_raw_rounds"); got != int64(b.KeepRawRounds) {
		t.Errorf("record kept_raw_rounds = %d, want the floor %d", got, b.KeepRawRounds)
	}
	if got := attrInt(t, rec, "compressed_msgs"); got != 0 {
		t.Errorf("record compressed_msgs = %d, want 0 for a pass that folded nothing", got)
	}
	if attrBool(t, rec, "history_changed") {
		t.Error("record history_changed = true for a pass that did not move the history")
	}
	// AC#5: the numerator behind the attempt is named, and this pass had no
	// provider measurement, so it must say estimate rather than stay blank.
	if got, ok := rec.attr["trigger_source"].(string); !ok || got != string(TriggerEstimate) {
		t.Errorf("record trigger_source = %v, want %q for a call with no provider measurement",
			rec.attr["trigger_source"], TriggerEstimate)
	}
	// Privacy: counts only, same hard rule as the success record.
	if strings.Contains(rec.flat(), sentinel) {
		t.Errorf("attempted record leaked history content: %q in %q", sentinel, rec.flat())
	}
}

// Ran must stay the mirror of the compacted state, or every existing reader
// (loop.go's replaceHistory gate, cmd/wisp's line, seven ticket-139/355 nails)
// silently means something else.
func TestCompressionStateRanIsTheMirrorOfCompacted(t *testing.T) {
	b := BudgetsFor(4096)
	c := NewCompressor(b, nil)

	cases := []struct {
		name string
		hist []llm.Message
		want CompressionState
	}{
		{"under-budget", buildRoundHistory(2, 20), CompressionSkipped},
		{"stuck-at-floor", buildRoundHistory(b.KeepRawRounds, 400), CompressionAttempted},
		{"real-fold", buildRoundHistory(6, 400), CompressionCompacted},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, rep, err := c.Compress(context.Background(), tc.hist)
			if err != nil {
				t.Fatalf("Compress: %v", err)
			}
			if rep.State != tc.want {
				t.Errorf("State = %q, want %q (est %d tokens, threshold %d)",
					rep.State, tc.want, c.TotalTokens(tc.hist), b.HistoryCompressTokens)
			}
			if want := rep.State == CompressionCompacted; rep.Ran != want {
				t.Errorf("Ran = %v but State = %q: Ran must be exactly (State == compacted), "+
					"or the old bool and the new enum answer differently", rep.Ran, rep.State)
			}
		})
	}
}

// AC#5(1): when the provider reported a count, THAT is the numerator - even
// when the local estimate disagrees, in both directions.
func TestCompressionTriggerPrefersProviderUsage(t *testing.T) {
	b := BudgetsFor(4096)
	c := NewCompressor(b, nil)
	hist := buildRoundHistory(2, 20) // estimate far under the budget

	est := c.TotalTokens(hist)
	if est >= b.HistoryCompressTokens {
		t.Fatalf("setup: estimate %d already clears %d, want it under so the two bases disagree",
			est, b.HistoryCompressTokens)
	}

	// The provider says the request really cost more than the budget.
	over := b.HistoryCompressTokens * 2
	r := c.Trigger(hist, over)
	if r.Source != TriggerProviderUsage {
		t.Errorf("Source = %q, want %q when a provider count is available", r.Source, TriggerProviderUsage)
	}
	if r.Tokens != over {
		t.Errorf("Tokens = %d, want the provider's %d (NOT the local estimate %d)", r.Tokens, over, est)
	}
	if !r.Need {
		t.Errorf("Need = false, but the provider reported %d over the %d threshold", over, b.HistoryCompressTokens)
	}
	if r.Threshold != b.HistoryCompressTokens {
		t.Errorf("Threshold = %d, want the frozen scaled %d", r.Threshold, b.HistoryCompressTokens)
	}
	if r.Observed != over {
		t.Errorf("Observed = %d, want %d", r.Observed, over)
	}

	// ...and the provider saying it fits must not be overruled by an estimate
	// that says it does not. Use a history the estimate calls over-budget.
	big := buildRoundHistory(6, 400)
	bigEst := c.TotalTokens(big)
	if bigEst <= b.HistoryCompressTokens {
		t.Fatalf("setup: estimate %d does not clear %d", bigEst, b.HistoryCompressTokens)
	}
	under := b.HistoryCompressTokens / 2
	r2 := c.Trigger(big, under)
	if r2.Source != TriggerProviderUsage || r2.Tokens != under {
		t.Errorf("Source/Tokens = %q/%d, want %q/%d: a provider count always wins over the estimate",
			r2.Source, r2.Tokens, TriggerProviderUsage, under)
	}
	if r2.Need {
		t.Errorf("Need = true on the provider's %d against %d, want false: the estimate's %d must not override it",
			under, b.HistoryCompressTokens, bigEst)
	}
}

// AC#5(2): no provider count (first turn, or a provider that reports nothing)
// falls back to exactly the local estimate - and says so.
func TestCompressionTriggerFallsBackToEstimateOnFirstTurn(t *testing.T) {
	b := BudgetsFor(4096)
	c := NewCompressor(b, nil)

	small := buildRoundHistory(2, 20)
	big := buildRoundHistory(6, 400)

	for _, tc := range []struct {
		name string
		hist []llm.Message
		want bool
	}{
		{"under", small, false},
		{"over", big, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			est := c.TotalTokens(tc.hist)
			r := c.Trigger(tc.hist, 0)
			if r.Source != TriggerEstimate {
				t.Errorf("Source = %q, want %q when the provider reported nothing", r.Source, TriggerEstimate)
			}
			if r.Tokens != est {
				t.Errorf("Tokens = %d, want the local estimate %d", r.Tokens, est)
			}
			if r.Need != tc.want {
				t.Errorf("Need = %v for %d tokens against %d, want %v", r.Need, est, b.HistoryCompressTokens, tc.want)
			}
			if r.Observed != 0 {
				t.Errorf("Observed = %d, want 0 (nothing was observed)", r.Observed)
			}
			// The fallback must agree with Need(), which is the estimate-only
			// reading the frozen budget contract is written against. If these
			// two ever disagree, a caller that kept Need() compresses on a
			// different trigger than the loop does.
			if got := c.Need(tc.hist); got != r.Need {
				t.Errorf("Trigger(hist, 0).Need = %v but Need(hist) = %v: the fallback is not the estimate path",
					r.Need, got)
			}
		})
	}
}

// AC#5's exit: a trigger that fired on the provider's number while this pass's
// own estimate is under the fold budget still reports an ATTEMPT with that
// numerator - the second shape that reaches the attempted state, and the one a
// bool could never express.
func TestCompressionAttemptedWhenProviderSaysOverButFoldHasNothing(t *testing.T) {
	recs := newTraceCapture()
	recs.assertRulerLive(t)

	b := BudgetsFor(4096)
	c := NewCompressor(b, nil, WithLogger(recs.logger()))
	hist := buildRoundHistory(2, 20)
	observed := b.HistoryCompressTokens * 10

	trig := c.Trigger(hist, observed)
	if !trig.Need || trig.Source != TriggerProviderUsage {
		t.Fatalf("setup: trigger = %+v, want a provider-usage trigger that fires", trig)
	}

	_, rep, err := c.Compress(context.Background(), hist, WithTrigger(trig))
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	if rep.Ran {
		t.Fatalf("setup: the pass folded, so this is not the shape")
	}
	if rep.State != CompressionAttempted {
		t.Errorf("State = %q, want %q: the trigger fired, so this pass WAS an attempt",
			rep.State, CompressionAttempted)
	}
	if rep.TriggerSource != TriggerProviderUsage {
		t.Errorf("report TriggerSource = %q, want %q - the basis must survive onto the report",
			rep.TriggerSource, TriggerProviderUsage)
	}
	if rep.TriggerTokens != observed {
		t.Errorf("report TriggerTokens = %d, want the provider's %d", rep.TriggerTokens, observed)
	}
	// The trigger did not get silently downgraded to the estimate.
	if rep.TriggerTokens == c.TotalTokens(hist) {
		t.Errorf("report TriggerTokens = %d equals the local estimate: the provider's %d was dropped",
			rep.TriggerTokens, observed)
	}

	hits := recs.with(traceSkipMsg)
	if len(hits) != 1 {
		t.Fatalf("%q records = %d, want 1 (all: %v)", traceSkipMsg, len(hits), flattenRecordMsgs(recs.all()))
	}
	rec := hits[0]
	if got, ok := rec.attr["reason"].(string); !ok || got != skipReasonUnderFoldBudget {
		t.Errorf("record reason = %v, want %q: the estimate was under the fold budget, so the loop never took a round",
			rec.attr["reason"], skipReasonUnderFoldBudget)
	}
	if got, ok := rec.attr["trigger_source"].(string); !ok || got != string(TriggerProviderUsage) {
		t.Errorf("record trigger_source = %v, want %q", rec.attr["trigger_source"], TriggerProviderUsage)
	}
	if got := attrInt(t, rec, "trigger_tokens"); got != int64(observed) {
		t.Errorf("record trigger_tokens = %d, want %d", got, observed)
	}
}

// A bare Compress call (no caller trigger) over the budget is still an attempt:
// the reading must not depend on a caller remembering to compute one.
func TestCompressionAttemptedWithoutCallerTrigger(t *testing.T) {
	b := BudgetsFor(4096)
	c := NewCompressor(b, nil)
	hist := stuckAtFloorHistory(b)

	_, rep, err := c.Compress(context.Background(), hist)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	if rep.State != CompressionAttempted {
		t.Errorf("State = %q, want %q: Compress's own estimate clears the budget, so this pass tried",
			rep.State, CompressionAttempted)
	}
	// With no caller reading, the report names its own estimate rather than
	// leaving the basis blank - a blank basis is the silent mixing AC#5 forbids.
	if rep.TriggerSource != TriggerEstimate {
		t.Errorf("TriggerSource = %q, want %q", rep.TriggerSource, TriggerEstimate)
	}
	if rep.TriggerTokens != rep.TokensBefore {
		t.Errorf("TriggerTokens = %d, want the pass's own estimate %d", rep.TriggerTokens, rep.TokensBefore)
	}
}

// Format contract: the enum round-trips by name, and an unset state does not
// masquerade as one of the three.
func TestCompressionStateStrings(t *testing.T) {
	for _, s := range []CompressionState{CompressionSkipped, CompressionAttempted, CompressionCompacted} {
		if !s.Valid() {
			t.Errorf("%q reported invalid", s)
		}
		if got := fmt.Sprintf("%v", s); got != string(s) {
			t.Errorf("%q formatted as %q", s, got)
		}
	}
	if CompressionState("").Valid() {
		t.Error("the empty state reported valid: an unfilled report must not read as a decision")
	}
	if got := CompressionState("bogus").String(); !strings.Contains(got, "bogus") {
		t.Errorf("String() of an unknown state = %q, want it to name the value", got)
	}
}
