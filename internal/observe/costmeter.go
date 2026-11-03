package observe

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// CostMeter (C23, ticket 44).
//
// What lives here: the reading side of cost accounting — the daily/monthly
// aggregates over cost_daily, the budget verdict (80% warn / 100% pause,
// over_budget=warn alternative), the price lookup with the 价格未知 marker, and
// the reconciliation check against task_log. What does NOT live here: the
// per-task counters (internal/agent's Cost, one per running task) and the
// storage primitive (internal/memory's cost_daily upsert).
//
// Why interfaces instead of *memory.Store: memory imports observe (its DAO
// validates D37 error classes through this package), so observe must not import
// memory back. The three interfaces below are the seams memory.Store satisfies
// — the same "narrow interface, concrete store satisfies it" shape
// internal/agent/journal.go already uses for the 04-tables writer.
//
// Day arithmetic: a "day" is one calendar day in the meter's Location (default
// time.Local, because a daily budget is a human concept: "today, in the user's
// timezone"). Persisted timestamps stay wall UTC unix seconds (D42#9); only the
// bucket label is local. A fixture clock is injected with WithClock, so
// rollover tests never wait for a real midnight.

// DayTotal is one cost_daily row (or the sum of several).
type DayTotal struct {
	Day       string
	TokensIn  int64
	TokensOut int64
	Micros    int64
	Tasks     int64
}

// Add folds another row's counters into t. Day is not touched: a sum over a
// range has no single day, and callers who print one use Days() to see which
// days went in.
func (t DayTotal) Add(o DayTotal) DayTotal {
	t.TokensIn += o.TokensIn
	t.TokensOut += o.TokensOut
	t.Micros += o.Micros
	t.Tasks += o.Tasks
	return t
}

// DayLedger is the cost_daily writer (internal/memory.Store.BumpCostDay).
type DayLedger interface {
	BumpCostDay(ctx context.Context, day string, tokensIn, tokensOut, costMicros int64) error
}

// DayReader is the aggregate read over a closed day range (inclusive both
// ends). internal/memory.Store.CostDaysBetween implements it with
// `day BETWEEN ? AND ?` — day strings are fixed-width 'YYYY-MM-DD', so the
// lexical order is the chronological order.
type DayReader interface {
	CostDaysBetween(ctx context.Context, fromDay, toDay string) ([]DayTotal, error)
}

// DayWriter overwrites one day's aggregate from a recomputation. It exists for
// the repair path (RefreshDay), not for the hot path: PostTask adds one task at
// a time through DayLedger.
type DayWriter interface {
	SetCostDay(ctx context.Context, day string, tokensIn, tokensOut, costMicros, tasks int64) error
}

// TaskReader is the base table behind the aggregate: the finished task_log rows
// whose cost columns must add up to the same numbers (ticket 44 AC#2's
// reconciliation).
type TaskReader interface {
	SumTasksEndedBetween(ctx context.Context, fromUnix, toUnix int64) (DayTotal, error)
}

// Clock is the meter's injectable time source. It returns wall time for
// bucketing persisted records only (D42#9); no interval math goes through it.
type Clock func() time.Time

// Option configures a CostMeter.
type Option func(*CostMeter)

// WithLedger attaches the cost_daily writer PostTask needs. Nil = the meter
// reads and judges but cannot book (a read-only view, e.g. a panel query).
func WithLedger(l DayLedger) Option {
	return func(m *CostMeter) { m.ledger = l }
}

// WithTasks attaches the task_log reader Reconcile needs.
func WithTasks(r TaskReader) Option {
	return func(m *CostMeter) { m.tasks = r }
}

// WithLocation sets the calendar a day belongs to. nil is ignored (the default
// stays time.Local) rather than becoming UTC by surprise: a budget that
// silently resets at 08:00 is a defect someone has to choose on purpose.
func WithLocation(loc *time.Location) Option {
	return func(m *CostMeter) {
		if loc != nil {
			m.loc = loc
		}
	}
}

// WithClock injects a fixture clock for rollover tests.
func WithClock(c Clock) Option {
	return func(m *CostMeter) {
		if c != nil {
			m.now = c
		}
	}
}

// WithBudget sets the [cost] caps this meter judges against.
func WithBudget(b BudgetParams) Option {
	return func(m *CostMeter) { m.budget = b }
}

// WithPriceTable sets the price table (version + cards) used for lookups and
// recorded in the startup line (AC#6).
func WithPriceTable(t PriceTable) Option {
	return func(m *CostMeter) { m.table = t }
}

// WithMeterLogger sends the meter's own records somewhere (nil = slog.Default).
func WithMeterLogger(l *slog.Logger) Option {
	return func(m *CostMeter) {
		if l != nil {
			m.log = l
		}
	}
}

// WithWriter attaches the aggregate repair mouth (SetCostDay) RefreshDay needs.
func WithWriter(w DayWriter) Option {
	return func(m *CostMeter) { m.writer = w }
}

// CostMeter is the C23 read/decide surface. The budgets and the price table can
// be replaced at runtime (config's hot tier, D36), so they sit behind mu; the
// store handles are frozen at construction and read without locking.
type CostMeter struct {
	days   DayReader
	ledger DayLedger
	writer DayWriter
	tasks  TaskReader
	// modelLedger books cost_daily (+ cost_daily_model in one storage batch);
	// models reads the per-day model rows (ticket 44b). Both are optional:
	// a meter built without them keeps 44's global-only behaviour.
	modelLedger ModelLedger
	models      ModelReader
	loc         *time.Location
	now         Clock
	log         *slog.Logger

	mu     sync.RWMutex
	budget BudgetParams
	table  PriceTable
}

// Budget returns the caps the meter judges against.
func (m *CostMeter) Budget() BudgetParams {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.budget
}

// SetBudget replaces the caps (a config reload hands the new [cost] section to
// the meter that is already running, instead of rebuilding the meter and
// losing its clock/calendar injections).
func (m *CostMeter) SetBudget(b BudgetParams) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.budget = b
}

// Judge evaluates the meter's caps against the day and month-to-date spends the
// caller already read (see BudgetParams.Judge for the rule).
func (m *CostMeter) Judge(day, month DayTotal) BudgetVerdict {
	return m.Budget().Judge(day, month)
}

// Rebuild is the reload path for the price table: the caps and the calendar
// stay, the cards change.
func (m *CostMeter) Rebuild(table PriceTable) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.table = table
}

// NewCostMeter builds a meter over the aggregate reader. The reader is required
// (a meter that can read nothing can only say "no data", which is not worth an
// object), the writer and the task reader are options.
func NewCostMeter(days DayReader, opts ...Option) (*CostMeter, error) {
	if days == nil {
		return nil, errors.New("observe: NewCostMeter needs a DayReader (cost_daily has no other mouth)")
	}
	m := &CostMeter{
		days: days,
		loc:  time.Local,
		now:  func() time.Time { return time.Now() },
		log:  slog.Default(),
	}
	for _, o := range opts {
		o(m)
	}
	return m, nil
}

// DayKey formats the calendar day an instant belongs to in this meter's
// calendar. Exported so the writer that books a row and the reader that asks
// for it compute the SAME label from the same instant: the day a task lands in
// is decided once, not re-derived by each side's own arithmetic.
func (m *CostMeter) DayKey(t time.Time) string {
	return DayKey(t.In(m.loc))
}

// DayKey is the canonical 'YYYY-MM-DD' label (the shape memory's validDay
// enforces on cost_daily.day).
func DayKey(t time.Time) string {
	return t.Format("2006-01-02")
}

// Today returns the aggregate for the calendar day the meter's clock is in.
func (m *CostMeter) Today(ctx context.Context) (DayTotal, error) {
	day := m.DayKey(m.now())
	rows, err := m.days.CostDaysBetween(ctx, day, day)
	if err != nil {
		return DayTotal{}, fmt.Errorf("observe: cost meter today %s: %w", day, err)
	}
	var sum DayTotal
	for _, r := range rows {
		sum = sum.Add(r)
	}
	if sum.Day == "" {
		sum.Day = day
	}
	return sum, nil
}

// MonthToDate returns the aggregate over the calendar month the meter's clock
// is in, from its first day to the current day. The range ends at "today", not
// at the last day of the month: cost_daily only ever holds days that already
// happened, so both give the same rows, and "to today" is the honest name for
// what a month budget measures (the ticket calls it 月累计, not 月预测).
func (m *CostMeter) MonthToDate(ctx context.Context) (DayTotal, error) {
	now := m.now().In(m.loc)
	day := DayKey(now)
	first := DayKey(time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, m.loc))
	rows, err := m.days.CostDaysBetween(ctx, first, day)
	if err != nil {
		return DayTotal{}, fmt.Errorf("observe: cost meter month %s..%s: %w", first, day, err)
	}
	var sum DayTotal
	for _, r := range rows {
		sum = sum.Add(r)
	}
	sum.Day = first + ".." + day
	return sum, nil
}

// Range returns the aggregates for a closed day range as stored (one row per
// day that has one), newest last.
func (m *CostMeter) Range(ctx context.Context, fromDay, toDay string) ([]DayTotal, error) {
	if fromDay > toDay {
		return nil, fmt.Errorf("observe: cost meter range %q..%q is inverted", fromDay, toDay)
	}
	return m.days.CostDaysBetween(ctx, fromDay, toDay)
}

// PostTask books one finished task into the daily aggregate (ticket 44's
// "per-task rollup -> cost_daily upsert"). It is the only place the meter
// writes, and it refuses a negative on the way in rather than letting an
// aggregate drift below zero: a negative cost is a bug in the caller's
// arithmetic, and storing it would hide that bug behind a smaller total.
//
// The day label comes from this meter's calendar and clock, and is returned so
// the caller's log line names the same day the row landed in.
func (m *CostMeter) PostTask(ctx context.Context, tokensIn, tokensOut, micros int64) (string, error) {
	if m.ledger == nil {
		return "", errors.New("observe: cost meter has no ledger (WithLedger was never given one)")
	}
	if tokensIn < 0 || tokensOut < 0 {
		return "", fmt.Errorf("observe: cannot book negative tokens (%d in / %d out)", tokensIn, tokensOut)
	}
	if micros < 0 {
		return "", fmt.Errorf("observe: cannot book a negative cost (%d micros)", micros)
	}
	day := m.DayKey(m.now())
	if err := m.ledger.BumpCostDay(ctx, day, tokensIn, tokensOut, micros); err != nil {
		return day, fmt.Errorf("observe: book cost_daily %s: %w", day, err)
	}
	return day, nil
}

// RefreshDay recomputes one day's aggregate from the task_log base table and
// writes it back (the repair path behind Reconcile). The day is re-derived from
// the meter's calendar, so a row booked under a different calendar than the one
// being read gets healed rather than staying split in two — that split is what
// an aggregate cache and its base table drift into when the writer and the
// reader each keep their own idea of "today".
//
// It refuses to run without both mouths (reader of the base table and writer of
// the aggregate): recomputing from nothing and then overwriting a real row with
// zeros would destroy the number it was supposed to protect.
func (m *CostMeter) RefreshDay(ctx context.Context, day string) (DayTotal, error) {
	if m.tasks == nil || m.writer == nil {
		return DayTotal{}, fmt.Errorf(
			"observe: RefreshDay(%s) needs both a task reader (WithTasks) and an aggregate writer (WithWriter)", day)
	}
	start, end, err := m.dayBounds(day)
	if err != nil {
		return DayTotal{}, err
	}
	base, err := m.tasks.SumTasksEndedBetween(ctx, start.Unix(), end.Unix())
	if err != nil {
		return DayTotal{}, fmt.Errorf("observe: refresh %s from task_log: %w", day, err)
	}
	if err := m.writer.SetCostDay(ctx, day, base.TokensIn, base.TokensOut, base.Micros, base.Tasks); err != nil {
		return DayTotal{}, fmt.Errorf("observe: refresh %s write: %w", day, err)
	}
	base.Day = day
	m.log.Info("cost: day aggregate refreshed", "day", day, "tasks", base.Tasks, "micros", base.Micros)
	return base, nil
}

// Reconcile compares one day's stored aggregate with the sum of the task_log
// rows that ended in it (AC#2: "reconciliation with task_log exact"). A nil
// TaskReader means the meter has no base table to check against, which is
// reported as an error rather than as "they agree".
func (m *CostMeter) Reconcile(ctx context.Context, day string) (Reconcile, error) {
	if m.tasks == nil {
		return Reconcile{}, errors.New("observe: cost meter has no task reader (WithTasks was never given one)")
	}
	start, end, err := m.dayBounds(day)
	if err != nil {
		return Reconcile{}, err
	}
	rows, err := m.days.CostDaysBetween(ctx, day, day)
	if err != nil {
		return Reconcile{}, fmt.Errorf("observe: reconcile %s aggregate: %w", day, err)
	}
	var agg DayTotal
	for _, r := range rows {
		agg = agg.Add(r)
	}
	base, err := m.tasks.SumTasksEndedBetween(ctx, start.Unix(), end.Unix())
	if err != nil {
		return Reconcile{}, fmt.Errorf("observe: reconcile %s task_log: %w", day, err)
	}
	return Reconcile{
		Day: day, Aggregate: agg, Base: base,
		Agree: agg.TokensIn == base.TokensIn && agg.TokensOut == base.TokensOut &&
			agg.Micros == base.Micros && agg.Tasks == base.Tasks,
	}, nil
}

// Reconcile is the outcome of comparing cost_daily with task_log for one day.
// Base.Day keeps the day label so a report says which day it compared.
type Reconcile struct {
	Day       string
	Aggregate DayTotal
	Base      DayTotal
	Agree     bool
}

// dayBounds returns [start, end) in wall time for one calendar day in the
// meter's calendar. A day that does not exist in that calendar (a DST spring
// forward can delete one) is refused by name: silently widening or shrinking
// the bucket is how an aggregate starts disagreeing with its base table.
func (m *CostMeter) dayBounds(day string) (time.Time, time.Time, error) {
	if len(day) != 10 || day[4] != '-' || day[7] != '-' {
		return time.Time{}, time.Time{}, fmt.Errorf("observe: cost day %q is not YYYY-MM-DD", day)
	}
	start, err := time.ParseInLocation("2006-01-02", day, m.loc)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("observe: cost day %q: %w", day, err)
	}
	end := start.AddDate(0, 0, 1)
	if end.Equal(start) {
		return time.Time{}, time.Time{}, fmt.Errorf("observe: cost day %q has no bounds in calendar %s", day, m.loc.String())
	}
	return start, end, nil
}

// Price looks a model's card up in the table the meter carries. known=false is
// the 价格未知 path: the caller shows tokens only, and PriceMicros for that
// card returns 0 rather than a guess.
func (m *CostMeter) Price(provider, model string) (Price, string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.table.Card(provider, model)
}

// PriceTableVersion reports the version of the table the meter prices with
// (AC#6's recorded number).
func (m *CostMeter) PriceTableVersion() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.table.Version
}

// RecordPriceTableVersion writes the startup record: which price table this
// meter prices with, and how many models it carries. AC#6 asks that a version
// bump be "recorded" — this is the durable record (slog JSONL is the diagnostic
// surface, SPEC-01 §observe), and the meter never recomputes an old aggregate:
// nothing in this file rewrites cost_daily except RefreshDay's repair, which
// recomputes from task_log, not from a price card.
func (m *CostMeter) RecordPriceTableVersion() {
	m.mu.RLock()
	version, updated, models := m.table.Version, m.table.Updated, len(m.table.Models())
	m.mu.RUnlock()
	m.log.Info("cost: price table",
		"version", version,
		"updated", updated,
		"models", models,
		"recomputes_aggregates", false)
}

// Now exposes the meter's clock so a caller's log line and the meter's buckets
// agree on one reading. The clock is frozen at construction (WithClock): a meter
// is rebuilt for a new calendar, never re-timed underneath a running one.
func (m *CostMeter) Now() time.Time { return m.now() }

// Calendar names the location the meter buckets by (for a log line that says
// which "today" it means).
func (m *CostMeter) Calendar() *time.Location { return m.loc }
