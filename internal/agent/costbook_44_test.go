package agent

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Ticket 44 AC#2/AC#3 at the loop's own seam: the task that ran books its day
// row, and the day row is what the budget gate reads. Everything here goes
// through the golden C5 seam (real SSE bytes, real adapter) and a real SQLite
// store — the two surfaces the ticket names, not a stand-in for either.

func localDay(t *testing.T) string {
	t.Helper()
	return time.Now().In(time.Local).Format("2006-01-02")
}

func newStore44(t *testing.T) *memory.Store {
	t.Helper()
	store, err := memory.Open(sealableTempDir124(t),
		memory.WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))))
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	return store
}

var price44 = config.Price{In: 2_000_000, Out: 10_000_000, Cached: 1_000_000}

// TestLoopBooksFinishedTaskIntoCostDay is the production-writer half of AC#2:
// before ticket 44 the loop closed task_log and cost_daily stayed at zero
// writers (the census in docs/reports/pending-and-issues.md names it GAP-19).
func TestLoopBooksFinishedTaskIntoCostDay(t *testing.T) {
	store := newStore44(t)
	h := newHarness(t, "tool-then-text", withJournal(store), withConfig(func(c *Config) {
		c.Currency = "CNY"
		c.Price = price44
	}))
	res := h.run("天气怎么样")
	if res.CostMicros == 0 {
		t.Fatalf("the run priced nothing: usage=%+v", res.Usage)
	}

	day := localDay(t)
	row, err := store.CostDay(context.Background(), day)
	if err != nil {
		t.Fatal(err)
	}
	if row.Tasks != 1 {
		t.Fatalf("cost_daily[%s] = %+v, want exactly one booked task (the loop never posted it)", day, row)
	}
	if row.CostMicros != res.CostMicros || row.TokensIn != int64(res.Usage.InputTokens) ||
		row.TokensOut != int64(res.Usage.OutputTokens) {
		t.Errorf("cost_daily row %+v does not match the task that produced it {%d/%d/micro=%d}",
			row, res.Usage.InputTokens, res.Usage.OutputTokens, res.CostMicros)
	}

	// A second task (its own golden replay, same store) adds to the same day
	// row instead of replacing it.
	h2 := newHarness(t, "tool-then-text", withJournal(store), withConfig(func(c *Config) {
		c.Currency = "CNY"
		c.Price = price44
	}))
	res2 := h2.run("再来一次")
	row, err = store.CostDay(context.Background(), day)
	if err != nil {
		t.Fatal(err)
	}
	if row.Tasks != 2 {
		t.Errorf("after two tasks cost_daily.tasks = %d, want 2", row.Tasks)
	}
	if row.CostMicros != res.CostMicros+res2.CostMicros {
		t.Errorf("aggregate micros = %d, want the two tasks summed (%d + %d)",
			row.CostMicros, res.CostMicros, res2.CostMicros)
	}
}

// TestAC2ReconcilesAgainstTaskLog is the "reconciliation with task_log exact"
// clause, taken through the same two surfaces the product uses: the loop wrote
// the rows, the meter compares them.
func TestAC2ReconcilesAgainstTaskLog(t *testing.T) {
	store := newStore44(t)
	for _, q := range []string{"一", "二"} {
		nh := newHarness(t, "tool-then-text", withJournal(store), withConfig(func(c *Config) {
			c.Currency = "CNY"
			c.Price = price44
		}))
		nh.run(q)
	}

	m, err := observe.NewCostMeter(store,
		observe.WithTasks(store), observe.WithLocation(time.Local))
	if err != nil {
		t.Fatal(err)
	}
	rec, err := m.Reconcile(context.Background(), localDay(t))
	if err != nil {
		t.Fatal(err)
	}
	if !rec.Agree {
		t.Fatalf("aggregate and base table disagree after two clean runs: agg=%+v base=%+v",
			rec.Aggregate, rec.Base)
	}
	if rec.Base.Tasks != 2 || rec.Aggregate.Tasks != 2 {
		t.Errorf("tasks agg=%d base=%d, want 2 and 2", rec.Aggregate.Tasks, rec.Base.Tasks)
	}
}

// TestAC3BudgetGateReadsTheBookedSpend is AC#3's gate: the spend that trips the
// threshold is the one cost_daily carries, not a number the caller brings.
func TestAC3BudgetGateReadsTheBookedSpend(t *testing.T) {
	store := newStore44(t)
	h := newHarness(t, "tool-then-text", withJournal(store), withConfig(func(c *Config) {
		c.Currency = "CNY"
		c.Price = price44
	}))
	res := h.run("天气")

	ctx := context.Background()
	// Cap the day at exactly the cost of one task: 100% reached, pause mode.
	params, err := observe.BudgetFrom(res.CostMicros, 0, 0.8, "pause")
	if err != nil {
		t.Fatal(err)
	}
	rep, err := h.costBudgetWith(ctx, params)
	if err != nil {
		t.Fatal(err)
	}
	if !rep.Metered {
		t.Fatal("the loop reported an unmetered budget while a store journal was wired")
	}
	if rep.Verdict.AllowNewTask {
		t.Errorf("at exactly 100%% of the daily cap new tasks were still allowed: %+v", rep.Verdict)
	}
	if rep.Verdict.Level != observe.LevelOver {
		t.Errorf("level = %q, want over", rep.Verdict.Level)
	}
	if rep.Day.Micros != res.CostMicros {
		t.Errorf("the gate read day spend %d, want the booked %d", rep.Day.Micros, res.CostMicros)
	}

	// One micro under the cap is a warn at 80%, and it still allows the task.
	params, err = observe.BudgetFrom(res.CostMicros*5, 0, 0.8, "pause")
	if err != nil {
		t.Fatal(err)
	}
	rep, err = h.costBudgetWith(ctx, params)
	if err != nil {
		t.Fatal(err)
	}
	if !rep.Verdict.AllowNewTask {
		t.Errorf("at 20%% of a 5x cap new tasks were blocked: %+v", rep.Verdict)
	}
	if rep.Verdict.Level != observe.LevelOK {
		t.Errorf("at 20%% the level should be ok, got %q", rep.Verdict.Level)
	}

	// warn mode: over the cap, still allowed, and the notice says why.
	params, err = observe.BudgetFrom(res.CostMicros, 0, 0.8, "warn")
	if err != nil {
		t.Fatal(err)
	}
	rep, err = h.costBudgetWith(ctx, params)
	if err != nil {
		t.Fatal(err)
	}
	if !rep.Verdict.AllowNewTask {
		t.Errorf("over_budget=warn blocked a task: %+v", rep.Verdict)
	}
	if rep.Verdict.Level != observe.LevelOver || !strings.Contains(rep.Verdict.Notice, "warn") {
		t.Errorf("warn-mode verdict = %+v", rep.Verdict)
	}
}

// TestAC6PriceTableBumpLeavesAggregatesAlone pins the second half of AC#6: the
// version in use is recorded, and stored aggregates are not re-derived.
func TestAC6PriceTableBumpLeavesAggregatesAlone(t *testing.T) {
	store := newStore44(t)
	h := newHarness(t, "tool-then-text", withJournal(store), withConfig(func(c *Config) {
		c.Currency = "CNY"
		c.Price = price44
	}))
	res := h.run("天气")
	day := localDay(t)
	before, err := store.CostDay(context.Background(), day)
	if err != nil {
		t.Fatal(err)
	}
	if before.Tasks != 1 || before.CostMicros != res.CostMicros {
		t.Fatalf("setup: %+v", before)
	}

	tbl, err := observe.DefaultPriceTable()
	if err != nil {
		t.Fatal(err)
	}
	m, err := observe.NewCostMeter(store,
		observe.WithPriceTable(tbl), observe.WithTasks(store), observe.WithWriter(store),
		observe.WithLocation(time.Local))
	if err != nil {
		t.Fatal(err)
	}
	m.RecordPriceTableVersion()

	// A version bump is a new table with a higher number; re-reading the day
	// after the bump returns exactly the row that was stored before it.
	bumped := tbl
	bumped.Version = tbl.Version + 1
	m.Rebuild(bumped)
	if m.PriceTableVersion() != tbl.Version+1 {
		t.Fatalf("meter did not take the bumped version: %d", m.PriceTableVersion())
	}
	after, err := store.CostDay(context.Background(), day)
	if err != nil {
		t.Fatal(err)
	}
	if after != before {
		t.Errorf("a price-table version bump changed a stored aggregate: before=%+v after=%+v", before, after)
	}

	// And the next task prices against the NEW table's version while the old
	// row stays where it was: the record says which version a number came from.
	if _, err := m.RefreshDay(context.Background(), day); err != nil {
		t.Fatal(err)
	}
	healed, err := store.CostDay(context.Background(), day)
	if err != nil {
		t.Fatal(err)
	}
	if healed.CostMicros != res.CostMicros {
		t.Errorf("RefreshDay re-priced history: %d, want the stored %d", healed.CostMicros, res.CostMicros)
	}
}

// TestAC4UnknownPriceKeepsTokensPins the loop-side half of AC#4: with no card
// the money is 0 and the marker is what a caller renders, while the token
// counters stay real.
func TestAC4UnknownPriceKeepsTokens(t *testing.T) {
	store := newStore44(t)
	h := newHarness(t, "tool-then-text", withJournal(store), withConfig(func(c *Config) {
		c.Currency = "CNY"
		c.Price = config.Price{} // nothing published
	}))
	res := h.run("天气")
	if res.CostMicros != 0 {
		t.Errorf("an unknown card produced %d micros", res.CostMicros)
	}
	if res.Usage.InputTokens == 0 || res.Usage.OutputTokens == 0 {
		t.Fatalf("tokens vanished with the money: %+v", res.Usage)
	}
	if res.Currency != "CNY" {
		t.Errorf("currency label = %q, want the config's own word (no conversion invented)", res.Currency)
	}
	// The marker is available and says the ticket's words, not a number.
	if observe.PriceUnknownMarker != "价格未知" {
		t.Errorf("marker = %q", observe.PriceUnknownMarker)
	}
	// The day row still books the task (a task that cost nothing is still a
	// task), with zero money.
	day := localDay(t)
	row, err := store.CostDay(context.Background(), day)
	if err != nil {
		t.Fatal(err)
	}
	if row.Tasks != 1 || row.CostMicros != 0 {
		t.Errorf("unpriced task booked %+v, want 1 task / 0 micros", row)
	}
}

// costBudgetWith installs caps on the loop's meter and reads the report back,
// so the gate is tested against the spend the store actually holds.
func (h *harness) costBudgetWith(ctx context.Context, p observe.BudgetParams) (CostReport, error) {
	h.t.Helper()
	h.loop.SetCostBudget(p)
	return h.loop.CostBudget(ctx)
}
