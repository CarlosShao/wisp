package observe

import (
	"context"
	"errors"
	"fmt"
)

// Per-model cost dimension (ticket 44b, D35 addition approved as ledger A1445).
//
// What lives here: the reading/writing seams for the new cost_daily_model table
// (day × provider × model), the same-batch booking method the loop's journal
// path uses, the model-side reconciliation, and the QuotaLedger seed that makes
// the quota gate's per-model baseline survive a restart. What does NOT live
// here: the SQL (internal/memory owns both tables' DAOs and the 4->5 migration
// step), the price cards (pricing.go), the chain stepping (ticket 11).
//
// Why interfaces again (same rule as costmeter.go's header): memory imports
// observe, so observe must not import memory back.

// CostModelDayTotal is one cost_daily_model row (or the sum of several).
type CostModelDayTotal struct {
	Day       string
	Provider  string
	Model     string
	TokensIn  int64
	TokensOut int64
	Micros    int64
	Tasks     int64
}

// Add folds another row's counters into t. Day/Provider/Model keep the first
// row's values: a folded row is a sum and printing it needs its own Day().
func (t CostModelDayTotal) Add(o CostModelDayTotal) CostModelDayTotal {
	t.TokensIn += o.TokensIn
	t.TokensOut += o.TokensOut
	t.Micros += o.Micros
	t.Tasks += o.Tasks
	return t
}

// Key is the "provider/model" spelling the quota ledger and the config
// text_chain already use.
func (t CostModelDayTotal) Key() string { return t.Provider + "/" + t.Model }

// ModelReader is the per-model aggregate read over a closed day range (memory's
// CostModelDaysBetween). Nil means "no model surface wired": reads report so
// rather than inventing zeros.
type ModelReader interface {
	CostModelDaysBetween(ctx context.Context, fromDay, toDay string) ([]CostModelDayTotal, error)
}

// ModelLedger is the SAME-BATCH writer: one call books the task into
// cost_daily and cost_daily_model in a single storage transaction (memory
// implements both INSERTs inside one write tx — a crash between the two must
// not leave a day whose model rows disagree with its global row).
type ModelLedger interface {
	BumpCostDayModel(ctx context.Context, day, provider, model string, tokensIn, tokensOut, costMicros int64) error
}

// WithModelLedger attaches the same-batch writer PostModelTask needs.
func WithModelLedger(l ModelLedger) Option {
	return func(m *CostMeter) { m.modelLedger = l }
}

// WithModelReader attaches the per-model reader SeedQuota/ReconcileModels need.
func WithModelReader(r ModelReader) Option {
	return func(m *CostMeter) { m.models = r }
}

// PostModelTask books one finished task into cost_daily AND cost_daily_model in
// the same batch (ticket 44b: "与 cost_daily 同批"). The day label comes from
// this meter's calendar and clock — the same single reading PostTask uses, so
// the two rows can never land in two different days.
//
// Negative counters are refused on the way in (same anti-negative rule as
// PostTask). Unknown price is booked as tokens-only with micros == 0: the 价格
// 未知 marker is the READER's rendering (PriceUnknownMarker at display time),
// never a second number invented at book time.
func (m *CostMeter) PostModelTask(ctx context.Context, provider, model string, tokensIn, tokensOut, micros int64) (string, error) {
	if m.modelLedger == nil {
		return "", errors.New("observe: cost meter has no model ledger (WithModelLedger was never given one)")
	}
	if provider == "" || model == "" {
		return "", fmt.Errorf("observe: cost_daily_model needs both provider and model (got %q/%q)", provider, model)
	}
	if tokensIn < 0 || tokensOut < 0 {
		return "", fmt.Errorf("observe: cannot book negative tokens (%d in / %d out)", tokensIn, tokensOut)
	}
	if micros < 0 {
		return "", fmt.Errorf("observe: cannot book a negative cost (%d micros)", micros)
	}
	day := m.DayKey(m.now())
	if err := m.modelLedger.BumpCostDayModel(ctx, day, provider, model, tokensIn, tokensOut, micros); err != nil {
		return day, fmt.Errorf("observe: book cost_daily(+model) %s %s/%s: %w", day, provider, model, err)
	}
	return day, nil
}

// ModelDay returns one day's per-model rows (empty slice = no model rows yet,
// which is also the honest shape for a day booked entirely before v5: the
// ticket's named backfill policy is "accumulate from now on", see 44br1-impl.md).
func (m *CostMeter) ModelDay(ctx context.Context, day string) ([]CostModelDayTotal, error) {
	if m.models == nil {
		return nil, errors.New("observe: cost meter has no model reader (WithModelReader was never given one)")
	}
	return m.models.CostModelDaysBetween(ctx, day, day)
}

// ModelMonthToDate returns the current calendar month's per-model rows,
// month-first-day through today (same range rule as MonthToDate).
func (m *CostMeter) ModelMonthToDate(ctx context.Context) ([]CostModelDayTotal, error) {
	if m.models == nil {
		return nil, errors.New("observe: cost meter has no model reader (WithModelReader was never given one)")
	}
	now := m.now().In(m.loc)
	first := MonthFirst(DayKey(now))
	today := DayKey(now)
	return m.models.CostModelDaysBetween(ctx, first, today)
}

// MonthFirst turns a 'YYYY-MM-DD' label into the 'YYYY-MM-01' label of the same
// calendar month. It is string surgery on purpose: the day shape is enforced
// narrow on every write (validDay), and a budget's month is a label bucket, not
// a computed date.
func MonthFirst(day string) string {
	if len(day) != 10 {
		return day
	}
	return day[:len("YYYY-MM")] + "-01"
}

// ReconcileModelDay folds one day's model rows and checks them against BOTH
// cost_daily and the task_log base table (ticket 44b AC#4: reconciliation exact
// at the model dimension). task_log carries no provider/model column (D35), so
// the model dimension can only ever be verified through its sums: the fold of
// every model row must equal the global day row, and the global day row must
// equal the task_log sum (the AC#2 rule costmeter.go's Reconcile already pins).
type ReconcileModelDay struct {
	Day      string
	ModelSum DayTotal // Σ over the day's cost_daily_model rows
	Agg      DayTotal // the cost_daily row
	Base     DayTotal // the task_log sum
	ModelsOk bool     // ModelSum == Agg on all four counters
	BaseOk   bool     // Agg == Base on all four counters
	Agree    bool     // both halves
}

// ReconcileModels runs that comparison. It needs the model reader, the day
// reader and the task reader; missing any is an error, never a vacuous "agree".
func (m *CostMeter) ReconcileModels(ctx context.Context, day string) (ReconcileModelDay, error) {
	if m.models == nil {
		return ReconcileModelDay{}, errors.New("observe: ReconcileModels needs a model reader (WithModelReader was never given one)")
	}
	if m.tasks == nil {
		return ReconcileModelDay{}, errors.New("observe: ReconcileModels needs a task reader (WithTasks was never given one)")
	}
	rows, err := m.models.CostModelDaysBetween(ctx, day, day)
	if err != nil {
		return ReconcileModelDay{}, fmt.Errorf("observe: reconcile models %s: %w", day, err)
	}
	var sum DayTotal
	for _, r := range rows {
		sum = sum.Add(DayTotal{TokensIn: r.TokensIn, TokensOut: r.TokensOut, Micros: r.Micros, Tasks: r.Tasks})
	}
	rec, err := m.Reconcile(ctx, day)
	if err != nil {
		return ReconcileModelDay{}, err
	}
	modelsOk := sum.TokensIn == rec.Aggregate.TokensIn && sum.TokensOut == rec.Aggregate.TokensOut &&
		sum.Micros == rec.Aggregate.Micros && sum.Tasks == rec.Aggregate.Tasks
	return ReconcileModelDay{
		Day: day, ModelSum: sum, Agg: rec.Aggregate, Base: rec.Base,
		ModelsOk: modelsOk,
		BaseOk:   rec.Agree,
		Agree:    modelsOk && rec.Agree,
	}, nil
}

// SeedFrom rebuilds the ledger's per-model buckets from stored rows
// (ticket 44b AC#3: the quota gate reads the same numbers after a restart as
// before it). Rows are the current month's cost_daily_model rows; currentDay is
// today's label; cards decides which model spends also draw the provider plan
// credit bucket (billing=plan, same rule as Post). Negative rows are refused
// rather than credited — the anti-negative rule Post already states.
//
// The whole body runs under l.mu and rolls through rollLocked, so the seed is
// one atomic step against Spend/Post like any other writer: a boot seed that
// overlapped the first dispatch's Decide used to write the three maps with no
// lock at all (docs/evidence/s1/44-quota-seedfrom-lock.md). Taking the lock
// here and then calling roll() would re-enter the same mutex and deadlock, so
// this path must use the caller-holds-the-lock spelling throughout.
func (l *QuotaLedger) SeedFrom(rows []CostModelDayTotal, currentDay string, cards map[string]QuotaModel) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.rollLocked()
	for _, r := range rows {
		if r.Micros < 0 {
			return fmt.Errorf("observe: quota ledger refuses negative stored spend (%d micros for %s)", r.Micros, r.Key())
		}
		key := r.Key()
		l.modelMon[key] += r.Micros
		if r.Day == currentDay {
			l.modelDay[key] += r.Micros
		}
		if c, ok := cards[key]; ok && c.Billing == BillingPlan {
			l.planMon[r.Provider] += r.Micros
		}
	}
	return nil
}
