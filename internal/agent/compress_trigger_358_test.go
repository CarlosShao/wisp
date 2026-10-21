package agent

// Ticket 358 AC#5's wiring teeth (the half a compressor unit test cannot reach):
// does the input-token count the PROVIDER reported actually become the trigger's
// numerator inside a real Loop.run, and does the fallback show up as estimate?
//
// tool-then-text.sse reports prompt_tokens=150 on its first response and
// prompt_tokens=180 on its second. With the window knocked down to 1024 the
// scaled threshold is 96, so 150 clears it while the loop's own estimate of the
// same seeded history (tens of tokens) does not. That gap is the reading:
// a trigger saying provider_usage/150 can only have come from turn.Usage.
//
// Mutations that turn these red:
//   - delete `prevUsage = turn.Usage` in loop.go: the second iteration's trigger
//     falls back to estimate and the case reads "estimate" - red.
//   - pass 0 instead of prevUsage.InputTokens into Trigger: same red.
//   - make the loop keep calling Need(hist) instead of Trigger: same red.

import (
	"testing"
)

// loopTriggerHarness builds a harness whose scaled threshold is low enough that
// the fixture's own prompt_tokens clears it while the seeded history's estimate
// does not. Every number is derived from the budget, never written by hand.
func loopTriggerHarness(t *testing.T, fixture string) (*harness, *traceCapture, int) {
	t.Helper()
	// Find the smallest window (in the 1024..4096 decade) whose scaled
	// threshold stays under the fixture's 150-token first response, so the
	// provider count is the numerator that can fire.
	for _, window := range []int{1024, 2048, 3072, 4096} {
		if b := BudgetsFor(window); b.HistoryCompressTokens < 150 {
			recs := newTraceCapture()
			recs.assertRulerLive(t)
			h := newHarness(t, fixture,
				withConfig(func(c *Config) { c.ContextWindow = window }),
				withLogger(recs.logger()))
			return h, recs, b.HistoryCompressTokens
		}
	}
	t.Fatalf("no window in [1024,4096] derives a threshold under the fixture's 150-token report")
	return nil, nil, 0
}

// The provider's number is the numerator, and it survives onto the report.
func TestLoopTriggerUsesProviderReportedInputTokens(t *testing.T) {
	h, recs, threshold := loopTriggerHarness(t, "tool-then-text")

	// Seed a small history: the point is that the ESTIMATE stays under the
	// threshold while the provider's count does not.
	for _, m := range buildRoundHistory(2, 20) {
		h.loop.append(m)
	}
	est := NewCompressor(h.loop.Budgets(), nil).TotalTokens(h.loop.History())
	if est >= threshold {
		t.Fatalf("setup: seeded estimate %d already clears %d, so this case could not tell the two bases apart",
			est, threshold)
	}

	res := h.run("再确认一下")
	if res.Status != StatusCompleted {
		t.Fatalf("run status = %s (%s), want completed", res.Status, res.Message)
	}
	if res.Rounds < 2 {
		t.Fatalf("run rounds = %d, want at least 2 so a second turn had a previous usage to read", res.Rounds)
	}
	if res.Usage.InputTokens <= 0 {
		t.Fatalf("run reports no input tokens at all, so the fixture no longer exercises AC#5")
	}

	rep := res.Compression
	if rep.TriggerSource != TriggerProviderUsage {
		t.Errorf("Compression.TriggerSource = %q, want %q on a run whose provider reported %d input tokens",
			rep.TriggerSource, TriggerProviderUsage, res.Usage.InputTokens)
	}
	// The number is the provider's, not the estimate. 150 is the fixture's own
	// first-response prompt_tokens; it is read from the report, and the
	// assertion is that it is NOT the local estimate (which is tens).
	if rep.TriggerTokens <= 0 {
		t.Fatalf("Compression.TriggerTokens = %d, want the provider's count", rep.TriggerTokens)
	}
	if rep.TriggerTokens == est {
		t.Errorf("Compression.TriggerTokens = %d equals the local estimate: the provider's count was dropped", rep.TriggerTokens)
	}
	if rep.TriggerTokens <= threshold {
		t.Errorf("Compression.TriggerTokens = %d does not clear the %d threshold, yet a trigger fired",
			rep.TriggerTokens, threshold)
	}
	// The attempt is visible on the outwire, which is the AC#1 half at loop level.
	if rep.State != CompressionAttempted {
		t.Errorf("Compression.State = %q, want %q: the provider's %d cleared the budget over a history the "+
			"estimate calls small, so the pass tried and had nothing foldable",
			rep.State, CompressionAttempted, rep.TriggerTokens)
	}
	if rep.Ran {
		t.Error("Compression.Ran = true, but nothing was foldable")
	}

	// ...and the trace on the sink names the same basis and the owning task, so
	// a reader off a disk sees which numerator fired and who paid for it.
	hits := recs.with(traceSkipMsg)
	if len(hits) != 1 {
		t.Fatalf("%q records = %d, want exactly 1 (all: %v)", traceSkipMsg, len(hits), flattenRecordMsgs(recs.all()))
	}
	rec := hits[0]
	if got, ok := rec.attr["trigger_source"].(string); !ok || got != string(TriggerProviderUsage) {
		t.Errorf("record trigger_source = %v, want %q", rec.attr["trigger_source"], TriggerProviderUsage)
	}
	if got := attrInt(t, rec, "trigger_tokens"); got != int64(rep.TriggerTokens) {
		t.Errorf("record trigger_tokens = %d, but the report says %d", got, rep.TriggerTokens)
	}
	if got := attrInt(t, rec, "threshold"); got != int64(threshold) {
		t.Errorf("record threshold = %d, want the window-derived %d", got, threshold)
	}
	if v, ok := rec.attr["task"].(string); !ok || v != res.TaskID {
		t.Errorf("record task = %v, want the Result's %q - the attempt must be attributable like the success record",
			rec.attr["task"], res.TaskID)
	}
}

// The other half of AC#5: on a first turn nothing has been observed yet, so the
// numerator must be the local estimate - named as such, never silently blended.
//
// This case drives the same fixture at a window whose threshold the ESTIMATE
// clears, which is what makes the estimate the reading.
func TestLoopTriggerFallsBackToEstimateWhenNothingObserved(t *testing.T) {
	recs := newTraceCapture()
	recs.assertRulerLive(t)
	// 128000 keeps the reference threshold; the seeded history is sized to
	// clear it, so the trigger fires with observed == 0 (first turn).
	h := newHarness(t, "text-reply",
		withConfig(func(c *Config) { c.ContextWindow = 128000 }),
		withLogger(recs.logger()))
	threshold := h.loop.Budgets().HistoryCompressTokens

	for _, m := range buildRoundHistory(15, 4000) {
		h.loop.append(m)
	}
	seeded := NewCompressor(h.loop.Budgets(), nil).TotalTokens(h.loop.History())
	if seeded <= threshold {
		t.Fatalf("setup: seeded %d tokens does not clear %d", seeded, threshold)
	}

	res := h.run("收个尾")
	if res.Status != StatusCompleted {
		t.Fatalf("run status = %s (%s), want completed", res.Status, res.Message)
	}

	rep := res.Compression
	if rep.TriggerSource != TriggerEstimate {
		t.Errorf("Compression.TriggerSource = %q, want %q: nothing had been observed when this trigger fired",
			rep.TriggerSource, TriggerEstimate)
	}
	if rep.State != CompressionCompacted {
		t.Fatalf("Compression.State = %q, want %q for a history with oldest rounds to fold",
			rep.State, CompressionCompacted)
	}
	// The estimate is the report's own pre-pass total, by construction: the
	// trigger runs on the history Compress then measures.
	if rep.TriggerTokens != rep.TokensBefore {
		t.Errorf("Compression.TriggerTokens = %d, want the estimate %d that decided the trigger",
			rep.TriggerTokens, rep.TokensBefore)
	}
	if !rep.Ran {
		t.Error("Compression.Ran = false while the pass folded")
	}
	hits := recs.with(traceMsg)
	if len(hits) != 1 {
		t.Fatalf("%q records = %d, want 1", traceMsg, len(hits))
	}
	if got, ok := hits[0].attr["trigger_source"].(string); !ok || got != string(TriggerEstimate) {
		t.Errorf("record trigger_source = %v, want %q", hits[0].attr["trigger_source"], TriggerEstimate)
	}
}

// The fallback must not be a permanent state: once a task HAS a measurement, the
// next trigger must use it rather than drifting back to the estimate. This is
// the regression a "compute once outside the loop" wiring would cause.
func TestLoopTriggerBasisCanChangeBetweenTurns(t *testing.T) {
	h, _, threshold := loopTriggerHarness(t, "tool-then-text")

	// Record what each turn's trigger said, by driving the compressor the loop
	// itself owns with the two numerators the run will produce.
	comp := h.loop.comp
	hist := buildRoundHistory(2, 20)
	first := comp.Trigger(hist, 0)
	if first.Source != TriggerEstimate {
		t.Fatalf("first turn basis = %q, want estimate", first.Source)
	}
	if first.Need {
		t.Fatalf("first turn: estimate %d clears %d, so this case cannot show the basis changing",
			first.Tokens, threshold)
	}
	second := comp.Trigger(hist, 150)
	if second.Source != TriggerProviderUsage {
		t.Errorf("second turn basis = %q, want provider_usage once a count exists", second.Source)
	}
	if !second.Need {
		t.Errorf("second turn: provider's 150 does not clear %d", threshold)
	}
	// The two readings describe the same history and still disagree, which is
	// exactly why AC#5 forbids mixing them silently.
	if first.Tokens == second.Tokens {
		t.Errorf("both bases reported %d tokens; the fixture no longer separates them", first.Tokens)
	}
}
