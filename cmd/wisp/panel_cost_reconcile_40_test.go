package main

// Readface wave 1 (2026-11-17), ticket 40 AC#3/AC#4's teeth (cmd/wisp half):
//
//	Test40CostReconcileAgreesToTheMicroUnit        正控: real task_log rows and
//	                                               the booked cost_daily row
//	                                               agree to the micro; a one-micro
//	                                               drift is named by Reconcile.
//	Test40BudgetThreeStatesFromTheAdapter          AC#4: <80 normal / >=80 warn /
//	                                               >=100 paused (and the
//	                                               over_budget=warn polarity),
//	                                               measured through the real
//	                                               adapter over a real store and
//	                                               config, never off a table.
//	Test40CostSummaryPayloadCarriesAggregateOnChain 正控 through the production
//	                                               chain: the payload line carries
//	                                               the aggregate block and the
//	                                               cached gate says known=false.

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/panel"
)

// Test40CostReconcileAgreesToTheMicroUnit is AC#3's 对账 tooth: two real
// task_log rows (StartTaskLog + FinishTaskLog, the loop's own close path) and
// the cost_daily row booked through the same ledger the meter uses must
// reconcile to the micro-unit - and a one-micro drift must be named, not
// averaged away.
func Test40CostReconcileAgreesToTheMicroUnit(t *testing.T) {
	dir := newInboundDataDir33(t, "ask_every_step")
	st, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	ctx := context.Background()
	if err := st.StartTaskLog(ctx, memory.TaskLog{ID: "task-rec-1", State: "running", QueryText: "rec"}); err != nil {
		t.Fatalf("StartTaskLog 1: %v", err)
	}
	if err := st.StartTaskLog(ctx, memory.TaskLog{ID: "task-rec-2", State: "running", QueryText: "rec"}); err != nil {
		t.Fatalf("StartTaskLog 2: %v", err)
	}
	if err := st.FinishTaskLog(ctx, "task-rec-1", "succeeded", "", 10, 2, 300, ""); err != nil {
		t.Fatalf("FinishTaskLog 1: %v", err)
	}
	if err := st.FinishTaskLog(ctx, "task-rec-2", "succeeded", "", 30, 4, 700, ""); err != nil {
		t.Fatalf("FinishTaskLog 2: %v", err)
	}
	// Book the two tasks the way the loop does - ONE BumpCostDay per task, so
	// the aggregate's task count is the base table's count too.
	today := observe.DayKey(time.Now())
	if err := st.BumpCostDay(ctx, today, 10, 2, 300); err != nil {
		t.Fatalf("BumpCostDay task 1: %v", err)
	}
	if err := st.BumpCostDay(ctx, today, 30, 4, 700); err != nil {
		t.Fatalf("BumpCostDay task 2: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	// Fresh read-only meter over the same store: the exact shape the panel
	// query could run. Agree=true means every micro booked into cost_daily is
	// accounted for in task_log and vice versa.
	st2, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer func() { _ = st2.Close() }()
	meter, err := observe.NewCostMeter(st2, observe.WithTasks(st2))
	if err != nil {
		t.Fatalf("NewCostMeter: %v", err)
	}
	rec, err := meter.Reconcile(ctx, today)
	if err != nil {
		t.Fatalf("Reconcile: %v", err)
	}
	if !rec.Agree {
		t.Fatalf("cost_daily vs task_log disagreed: aggregate=%+v base=%+v (AC#3 对账到微单位)", rec.Aggregate, rec.Base)
	}
	if rec.Aggregate.Micros != 1000 || rec.Base.Micros != 1000 {
		t.Errorf("reconciled micros = %d/%d, want 1000/1000 (the sum of the two task rows)", rec.Aggregate.Micros, rec.Base.Micros)
	}

	// The negative half: one extra micro booked into cost_daily that task_log
	// does not know about must flip the verdict, with both numbers named.
	st3, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("reopen for drift: %v", err)
	}
	if err := st3.BumpCostDay(ctx, today, 0, 0, 1); err != nil {
		t.Fatalf("BumpCostDay drift: %v", err)
	}
	_ = st3.Close()
	st4, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("reopen after drift: %v", err)
	}
	defer func() { _ = st4.Close() }()
	meter2, err := observe.NewCostMeter(st4, observe.WithTasks(st4))
	if err != nil {
		t.Fatalf("NewCostMeter after drift: %v", err)
	}
	rec2, err := meter2.Reconcile(ctx, today)
	if err != nil {
		t.Fatalf("Reconcile after drift: %v", err)
	}
	if rec2.Agree {
		t.Fatal("a one-micro drift still reconciled - the tooth cannot see the drift it was built to catch")
	}
	if rec2.Aggregate.Micros != 1001 || rec2.Base.Micros != 1000 {
		t.Errorf("drift readings = aggregate %d / base %d, want 1001/1000", rec2.Aggregate.Micros, rec2.Base.Micros)
	}
}

// writeCostConfig40 writes a minimal config.toml with the [cost] leaves the
// adapter reads and builds a Manager over it.
func writeCostConfig40(t *testing.T, dir string, overBudget string) *config.Manager {
	t.Helper()
	body := "schema_version = 2\n\n[cost]\nmonthly_budget = 1000\ndaily_budget = 0\nalert_threshold = 0.8\nover_budget = \"" + overBudget + "\"\n"
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	mgr, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	return mgr
}

func seedTodayCost40(t *testing.T, dir string, micros int64) {
	t.Helper()
	st, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	defer func() { _ = st.Close() }()
	if err := st.BumpCostDay(context.Background(), observe.DayKey(time.Now()), 1, 1, micros); err != nil {
		t.Fatalf("BumpCostDay: %v", err)
	}
}

// Test40BudgetThreeStatesFromTheAdapter walks the three states through the
// real adapter (real store rows, real config file), plus the
// over_budget=warn polarity that keeps 100% a warning instead of a pause.
func Test40BudgetThreeStatesFromTheAdapter(t *testing.T) {
	cases := []struct {
		name       string
		overBudget string
		seedMicros int64
		want       string
		wantAllow  bool
	}{
		{"below 80 percent is normal", "pause", 700, "normal", true},
		{"at 80 percent is warn", "pause", 800, "warn", true},
		{"at 100 percent with pause is paused", "pause", 1000, "paused", false},
		{"at 100 percent with warn stays warn", "warn", 1000, "warn", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := newInboundDataDir33(t, "ask_every_step")
			seedTodayCost40(t, dir, c.seedMicros)
			mgr := writeCostConfig40(t, dir, c.overBudget)
			r := newCostReader194(dir, mgr)
			row, err := r.CostAggregate(context.Background())
			if err != nil {
				t.Fatalf("CostAggregate: %v", err)
			}
			if row.BudgetState != c.want {
				t.Errorf("BudgetState = %q, want %q (month=%d cap=%d)", row.BudgetState, c.want, row.MonthCostMicros, row.MonthlyBudgetMicros)
			}
			if row.BudgetAllowNewTask != c.wantAllow {
				t.Errorf("BudgetAllowNewTask = %v, want %v", row.BudgetAllowNewTask, c.wantAllow)
			}
			if row.MonthCostMicros != c.seedMicros {
				t.Errorf("MonthCostMicros = %d, want the seeded %d (month window includes today)", row.MonthCostMicros, c.seedMicros)
			}
			if row.MonthlyBudgetMicros != 1000 {
				t.Errorf("MonthlyBudgetMicros = %d, want the configured 1000", row.MonthlyBudgetMicros)
			}
			if row.PricingVersion < 1 {
				t.Errorf("PricingVersion = %d, want the embedded table's own version (>=1)", row.PricingVersion)
			}
			if row.DayCostMicros != c.seedMicros {
				t.Errorf("DayCostMicros = %d, want the seeded %d", row.DayCostMicros, c.seedMicros)
			}
		})
	}
}

// Test40CostSummaryPayloadCarriesAggregateOnChain reads the aggregate block
// off the production chain's cost.summary payload: the block is on the wire
// with the store's own numbers, and the cached gate says known=false so the
// page never renders the zero as a reading.
func Test40CostSummaryPayloadCarriesAggregateOnChain(t *testing.T) {
	dir := newInboundDataDir33(t, "ask_every_step")
	seedTodayCost40(t, dir, 1234)
	disp, _, err := newComposerDispatchChain(dir, auditQuiet194, panelInboundActor, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("the production chain did not assemble on %s: %v", dir, err)
	}
	reply, err := disp.Handle(context.Background(), costEnvelope194("rid-40-agg"))
	if err != nil {
		t.Fatalf("cost.summary refused on a wired chain: %v\nreply=%q", err, reply)
	}
	var payload struct {
		Count          int  `json:"count"`
		CachedKnown    bool `json:"cachedKnown"`
		AggregateKnown bool `json:"aggregateKnown"`
		Aggregate      *struct {
			DayCostMicros      int64  `json:"dayCostMicros"`
			MonthCostMicros    int64  `json:"monthCostMicros"`
			BudgetState        string `json:"budgetState"`
			BudgetAllowNewTask bool   `json:"budgetAllowNewTask"`
			PricingVersion     int    `json:"pricingVersion"`
		} `json:"aggregate"`
	}
	for _, line := range strings.Split(reply, "\n") {
		if strings.HasPrefix(line, panel.CostSummaryPayloadMarker) {
			if err := json.Unmarshal([]byte(strings.TrimPrefix(line, panel.CostSummaryPayloadMarker)), &payload); err != nil {
				t.Fatalf("payload is not JSON: %v", err)
			}
		}
	}
	if payload.Aggregate == nil || !payload.AggregateKnown {
		t.Fatalf("payload carries no aggregate block (known=%v): %+v", payload.AggregateKnown, payload)
	}
	if payload.Aggregate.DayCostMicros != 1234 || payload.Aggregate.MonthCostMicros != 1234 {
		t.Errorf("aggregate micros = day %d month %d, want the seeded 1234/1234",
			payload.Aggregate.DayCostMicros, payload.Aggregate.MonthCostMicros)
	}
	if payload.Aggregate.BudgetState == "" {
		t.Error("aggregate carries no budget state word")
	}
	if payload.Aggregate.PricingVersion < 1 {
		t.Errorf("aggregate carries pricing version %d, want the embedded table's own", payload.Aggregate.PricingVersion)
	}
	if payload.CachedKnown {
		t.Error("payload claims cachedKnown=true: cost_daily aggregates no cached column and the gate must say so")
	}
	if payload.Count != 1 {
		t.Errorf("days count = %d, want the seeded 1 row", payload.Count)
	}
}
