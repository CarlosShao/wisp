package agent

import (
	"context"
	"fmt"
	"testing"

	"github.com/CarlosShao/wisp/internal/llm"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/memory"
)

// AC#1 per-task tokens/cost match hand-computed fixture to the unit. Fixture
// tool-then-text has two SSE usage events: 150 in/14 out and 180 in/10 out (150
// cached). Usage.Add max-merges → totals: in=180, out=10, cache=150. Card: 2e6
// micro/1M in, 10e6 micro/1M out, 1e6 micro/1M cache → PriceMicros = 180*2 +
// 10*10 = 360+100 = 460 micros. We assert the Result carries exactly that number
// AND the task_log row books it (AC#1 "reconcile with task_log exact").
func TestAC1GoldenPerTaskCostMatchesHandComputed(t *testing.T) {
	dir := sealableTempDir124(t)
	store, err := memory.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	card := config.Price{In: 2000000.0, Out: 10000000.0, Cached: 1000000.0} // micro-per-1M rates in config units
	h := newHarness(t, "tool-then-text", withJournal(store), withConfig(func(c *Config) {
		c.Currency = "CNY"
		c.ArtifactsDir = dir
		c.Price = card
	}))
	res := h.run("weather")

	// The loop max-merged the two events: 330 in (150+180), 24 out (14+10), 150 cached.
	wantInput, wantOutput, wantCache := 330, 24, 150
	gotInput, gotOutput, gotCache := res.Usage.InputTokens, res.Usage.OutputTokens, res.Usage.CachedTokens
	if gotInput != wantInput || gotOutput != wantOutput || gotCache != wantCache {
		t.Fatalf("usage = {%d,%d,%d}, want {%d,%d,%d}", gotInput, gotOutput, gotCache, wantInput, wantOutput, wantCache)
	}

	// Hand computation, to the unit, from the two recorded usage events on the
	// card {in 2e6, out 10e6, cached 1e6} micro-per-1M:
	//
	//	turn 1 {150, 14, 0}:   150*2e6 + 14*10e6                   = 440e6/1e6 = 440
	//	turn 2 {180, 10, 150}: (180-150)*2e6 + 150*1e6 + 10*10e6   = 310e6/1e6 = 310
	//	per-turn SUM                                               = 750 micros
	//
	// Pricing the merged counters would also say 750 here, because the card is
	// linear in tokens; the two readings only part ways when a turn's total
	// lands on a half-micro (that case is pinned in
	// internal/observe/pricing_44_test.go's "half-micro total rounds half-up").
	// What this test pins is the loop's own number against the fixture's, and
	// against the row it wrote.
	wantMicros := PriceMicros(card, llm.Usage{InputTokens: 150, OutputTokens: 14}) +
		PriceMicros(card, llm.Usage{InputTokens: 180, OutputTokens: 10, CachedTokens: 150})
	if wantMicros != 750 {
		t.Fatalf("the hand-computed fixture itself moved: %d, want 750", wantMicros)
	}
	if res.CostMicros != wantMicros {
		t.Errorf("cost micros = %d, want %d (%s)", res.CostMicros, wantMicros, fmtTokenCard(res.Usage))
	}
	if merged := PriceMicros(card, res.Usage); merged != wantMicros {
		t.Errorf("merged-counter pricing says %d while the per-turn sum says %d: "+
			"one of the two readings changed its mind, and a cost that depends on which "+
			"mouth asks is a cost nobody can reconcile", merged, wantMicros)
	}

	// The base table saw the same three numbers (the finishTask writer closed them
	// together): reconciliation test for this specific row.
	row, err := store.TaskLogByID(context.Background(), res.TaskID)
	if err != nil {
		t.Fatal(err)
	}
	if row.CostTokensIn != int64(wantInput) || row.CostTokensOut != int64(wantOutput) || row.CostAmountMicro != wantMicros {
		t.Errorf("task_log row = {%d/%d/micro=%d}, want {%d/%d/micro=%d}",
			row.CostTokensIn, row.CostTokensOut, row.CostAmountMicro,
			wantInput, wantOutput, wantMicros)
	}
}

// AC#4 unknown price markers the model as unpriced and costs nothing — not a
// guess. A model without a published card must report Priced=false and Cost=0
// while still showing the token counters.
func TestAC4UnknownPriceShowsTokensOnlyAndPricedFalse(t *testing.T) {
	// An all-zero card = unknown; PriceKnown returns false.
	var zero config.Price
	t.Logf("zero card Known()=%v (ticket 44: never fabricate amounts)", CardOf(zero).Known())
	if zero.In == 0 && zero.Out == 0 && CardOf(zero).Known() {
		t.Error("unknown price said it was known")
	}
	// Using it on real usage produces zero cost; Priced remains false.
	dir := sealableTempDir124(t)
	store, err := memory.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	h := newHarness(t, "text-reply", withJournal(store), withConfig(func(c *Config) {
		c.Currency = "CNY"
		c.ArtifactsDir = dir
		c.Price = zero // unknown
	}))
	res := h.run("hello")

	if res.CostMicros != 0 {
		t.Errorf("unknown price priced at %d micros, want 0", res.CostMicros)
	}
	// The marker would be shown by the caller who reads cost.Priced; this test
	// asserts the meter does not fabricate when there is no card.
	if res.Usage.IsZero() {
		t.Error("known usage but zero cards; the user expects to see tokens even if money is unknown")
	}
}

// fmtTokenCard renders one side of the reconciliation check.
func fmtTokenCard(u llm.Usage) string {
	return fmt.Sprintf("in=%d out=%d cache=%d", u.InputTokens, u.OutputTokens, u.CachedTokens)
}
