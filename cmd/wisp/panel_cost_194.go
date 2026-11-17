package main

// Ticket 194's cost door gets its real leg here (ledger A1681's AC#5 boundary:
// the internal/panel batch landed the door with a nil socket that refused by
// name; this file is the cmd/wisp wave that attaches the reader).
//
// WHAT IT IS. panel.CostReader is the injected cost leg (cost_handlers.go: one
// read method, nothing else, so the door is read-only by construction). This
// adapter answers it from memory.Store.ListCostDaily (dao_misc.go) - the same
// C23 cost_daily table the task loop's BumpCostDay writes - so cost.summary
// returns the machine's real per-day rows instead of the named refusal. The
// mapping is one struct field per column, integers only; the adapter derives
// nothing, adds nothing, and reorders nothing: ListCostDaily already returns
// newest-day-first and panel's handler repeats that contract in its prose.
//
// WHY OPEN PER READ. The same shape every read outlet in this package already
// uses (sessions.go's list/search/outlets: memory.Open, one query, defer
// Close). The cost door is a user-triggered, rare read - not a hot path - and
// a handle held for the process's whole life would pin wisp.db open for a
// session that may never ask again, while a failed open would either fail the
// whole chain's assembly for a door the request never used or have to be
// cached as a permanent refusal. Per-read keeps every failure on the request
// that actually hit it: a database that will not open refuses THE COST REQUEST
// - loudly, through the handler's read-failed audit branch - and the next
// request starts clean.

import (
	"context"
	"fmt"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/panel"
)

// costReader194 adapts memory.Store's cost DAOs into panel.CostReader (and,
// since readface wave 1, panel.CostAggregateReader - ticket 40 AC#3's Go
// half). It holds the data root and the config Manager the [cost] section
// comes from; no database handle lives between reads.
type costReader194 struct {
	dir string
	mgr *config.Manager
}

// newCostReader194 builds the adapter over the data root and the Manager the
// chain was given. No database is opened here - see the file comment. The
// Manager may be nil (a chain with no config): the aggregate face then
// reports uncapped budgets from a nil config read, honestly, instead of
// guessing the [cost] section.
func newCostReader194(dir string, mgr *config.Manager) *costReader194 {
	return &costReader194{dir: dir, mgr: mgr}
}

// CostDays answers panel.CostReader: the real cost_daily rows, newest day
// first, one view row per stored row. An empty table answers an empty slice -
// zero rows is a true answer for a machine that has not run a task, not an
// error and not a stub.
func (r *costReader194) CostDays(ctx context.Context) ([]panel.CostDayRow, error) {
	st, err := memory.Open(r.dir)
	if err != nil {
		return nil, fmt.Errorf("成本存储不可用（%s）：%w", r.dir, err)
	}
	defer func() { _ = st.Close() }()
	days, err := st.ListCostDaily(ctx)
	if err != nil {
		return nil, fmt.Errorf("成本读取失败（%s）：%w", r.dir, err)
	}
	out := make([]panel.CostDayRow, 0, len(days))
	for _, d := range days {
		out = append(out, panel.CostDayRow{
			Day:        d.Day,
			TokensIn:   d.TokensIn,
			TokensOut:  d.TokensOut,
			CostMicros: d.CostMicros,
			Tasks:      d.Tasks,
			// CachedTokens stays 0 by construction: cost_daily aggregates no
			// cached-token column and this leg may not widen the D35 table.
			// The payload's CachedKnown=false gate (cost_handlers.go) is what
			// tells the page not to read the zero as a number.
			CachedTokens: 0,
		})
	}
	return out, nil
}

// CostAggregate answers panel.CostAggregateReader (readface wave 1): the
// month-to-date totals summed from the store's own cost_daily rows, the
// [cost] budgets read off the live config, the three-state verdict computed
// by observe.BudgetParams.Judge (the meter's own judge - basis points, no
// float touches money here), and the pricing version the embedded price
// table carries.
func (r *costReader194) CostAggregate(ctx context.Context) (panel.CostAggregateRow, error) {
	row := panel.CostAggregateRow{BudgetState: "normal", BudgetAllowNewTask: true}
	now := time.Now()
	today := observe.DayKey(now)
	st, err := memory.Open(r.dir)
	if err != nil {
		return row, fmt.Errorf("成本存储不可用（%s）：%w", r.dir, err)
	}
	defer func() { _ = st.Close() }()

	// Today's own row (0 when nothing ran today - the store's own answer).
	if day, err := st.CostDay(ctx, today); err == nil {
		row.DayCostMicros = day.CostMicros
	}

	// Month-to-date: one range read over the store, summed from the rows it
	// returns. The month window is the LOCAL calendar month (a monthly budget
	// is a human concept, the same reading the meter's day bucketing uses);
	// observe.DayKey is the same label the booking side computes.
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	fromDay := observe.DayKey(firstOfMonth)
	days, err := st.CostDaysBetween(ctx, fromDay, today)
	if err != nil {
		return row, fmt.Errorf("成本月累计读取失败（%s）：%w", r.dir, err)
	}
	month := observe.DayTotal{}
	for _, d := range days {
		month = month.Add(d)
	}
	row.MonthTokensIn = month.TokensIn
	row.MonthTokensOut = month.TokensOut
	row.MonthCostMicros = month.Micros

	// The budgets come off the live [cost] section; the verdict is the
	// meter's own judge over the two aggregates just read.
	var daily, monthly int64
	var alertRatio float64
	overBudget := "pause"
	if r.mgr != nil {
		if cfg := r.mgr.Config(); cfg != nil {
			daily, monthly = cfg.Cost.DailyBudget, cfg.Cost.MonthlyBudget
			alertRatio, overBudget = cfg.Cost.AlertThreshold, cfg.Cost.OverBudget
		}
	}
	row.DailyBudgetMicros = daily
	row.MonthlyBudgetMicros = monthly
	if b, err := observe.BudgetFrom(daily, monthly, alertRatio, overBudget); err == nil {
		v := b.Judge(observe.DayTotal{Micros: row.DayCostMicros}, observe.DayTotal{Micros: row.MonthCostMicros})
		row.BudgetState = budgetStateWord(v)
		row.BudgetAllowNewTask = v.AllowNewTask
	}

	// The pricing version footnote: the embedded table's own version + date.
	if tbl, err := observe.DefaultPriceTable(); err == nil {
		row.PricingVersion = tbl.Version
		row.PricingUpdated = tbl.Updated
	}
	return row, nil
}

// budgetStateWord maps observe's BudgetVerdict onto the page's three-state
// word (ticket 40 AC#4): normal below the alert line, warn at/above it,
// paused when the limit is reached AND over_budget=pause - i.e. when new
// tasks are actually being refused. The word list is these three; anything
// else would be a second vocabulary.
func budgetStateWord(v observe.BudgetVerdict) string {
	switch {
	case v.Level == observe.LevelOver && !v.AllowNewTask:
		return "paused"
	case v.Level == observe.LevelOver || v.Level == observe.LevelWarn:
		return "warn"
	default:
		return "normal"
	}
}
