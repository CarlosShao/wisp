package agent

import (
	"context"
	"log/slog"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Ticket 44 (C23): the loop is the one place a task's final numbers exist, so
// the loop is where the daily aggregate gets booked.
//
// Why the loop and not the composition root: task_log is closed from two sites
// (the loop's own journal, and cmd/wisp's after-the-fact write for the shapes
// that never reached the loop, run.go's !opened branch), and only the loop sees
// the per-task Cost it just accumulated. Booking here means every task that
// ever produced a usage event reaches cost_daily, which is what closes the
// "cost_daily has zero production writers" row in the ledger (GAP-19 / ticket
// 300's census). cmd/wisp needs no change to get the aggregate.
//
// The Journal interface is NOT widened: ticket 10 froze it, and test doubles
// across this package implement it. The aggregate mouths are discovered by
// assertion, so a Journal that is only a Journal (a fake, or a store built
// without the cost DAO) simply books nothing, and the loop says so in the log
// rather than inventing an aggregate.

// costMeterFor assembles the C23 meter over whatever the journal can do. A nil
// result means "this journal has no cost surface", which the loop then reports
// once per task instead of pretending the day row exists.
func costMeterFor(j Journal, log *slog.Logger) *observe.CostMeter {
	if j == nil {
		return nil
	}
	reader, ok := j.(observe.DayReader)
	if !ok {
		return nil
	}
	opts := []observe.Option{observe.WithMeterLogger(log)}
	if ledger, ok := j.(observe.DayLedger); ok {
		opts = append(opts, observe.WithLedger(ledger))
	}
	if writer, ok := j.(observe.DayWriter); ok {
		opts = append(opts, observe.WithWriter(writer))
	}
	if tasks, ok := j.(observe.TaskReader); ok {
		opts = append(opts, observe.WithTasks(tasks))
	}
	m, err := observe.NewCostMeter(reader, opts...)
	if err != nil {
		return nil
	}
	return m
}

// bookCostDay posts the finished task's rollup into cost_daily. It runs on the
// same detached context as the task_log close (a cancelled task still owes its
// forensics, and its cost is part of that record), and its failure is a warning
// with the day named, never a change to the task's result: the task ended, and
// an aggregate that missed one row is repaired by CostMeter.RefreshDay rather
// than by failing the task that found out too late.
func (l *Loop) bookCostDay(ctx context.Context, cost Cost, taskID string) {
	if l.costMeter == nil {
		if l.opt.Journal != nil {
			l.log().Warn("agent: cost_daily not booked: the wired journal has no cost surface",
				"task", taskID)
		}
		return
	}
	day, err := l.costMeter.PostTask(ctx,
		int64(cost.Usage.InputTokens), int64(cost.Usage.OutputTokens), cost.Micros)
	if err != nil {
		l.log().Warn("agent: cost_daily book failed", "err", err, "task", taskID, "day", day)
		return
	}
	l.log().Debug("agent: cost_daily booked",
		"task", taskID, "day", day,
		"tokens_in", cost.Usage.InputTokens, "tokens_out", cost.Usage.OutputTokens,
		"micros", cost.Micros, "priced", cost.Priced)
}

// CostReport is what the host reads before it starts a new task: the spend the
// numbers came from, plus the verdict computed from them.
type CostReport struct {
	// Day is the aggregate for the calendar day the meter is in; Month is the
	// month-to-date aggregate. Both are zero when the meter has no store.
	Day   observe.DayTotal
	Month observe.DayTotal
	// Verdict is the budget reading over those two numbers.
	Verdict observe.BudgetVerdict
	// Metered says whether a cost surface was wired at all. false means the
	// three fields above are the "nothing was measured" shape, NOT the "cost
	// is zero" shape — a caller that reads AllowNewTask without checking this
	// gets the right answer for the wrong reason.
	Metered bool
}

// CostBudget reports the money budget the loop's meter judges new tasks
// against, plus the day/month spend it was judged on. The host calls this
// before starting a task (ticket 44 AC#3: 80% warns, 100% pauses NEW tasks); a
// task already running is never read here, because reaching a budget does not
// stop the work in flight — that is the loop's own token/round brake (C22), a
// different layer with a different number.
func (l *Loop) CostBudget(ctx context.Context) (CostReport, error) {
	if l.costMeter == nil {
		return CostReport{Verdict: observe.BudgetVerdict{AllowNewTask: true, Level: observe.LevelOK}}, nil
	}
	day, err := l.costMeter.Today(ctx)
	if err != nil {
		return CostReport{Metered: true}, err
	}
	month, err := l.costMeter.MonthToDate(ctx)
	if err != nil {
		return CostReport{Day: day, Metered: true}, err
	}
	return CostReport{
		Day: day, Month: month,
		Verdict: l.costMeter.Judge(day, month),
		Metered: true,
	}, nil
}

// preDispatchBudget is the C23 gate in front of a NEW task (ticket 44 AC#3).
// It returns the verdict when the task must be refused, and nil when it may
// start — including the three shapes where refusing would be a lie: no cost
// surface wired, no cap configured, or the aggregate could not be read.
//
// The read failure is the deliberate one: it logs and lets the task start.
// A store hiccup that blocks every task turns an accounting problem into a
// product outage, and the outage is the more expensive of the two. It is loud
// (Warn, with the error), not silent, so a budget that never actually gates is
// still discoverable from the log.
func (l *Loop) preDispatchBudget(ctx context.Context, taskID string) *observe.BudgetVerdict {
	if l.costMeter == nil {
		return nil
	}
	rep, err := l.CostBudget(ctx)
	if err != nil {
		l.log().Warn("agent: cost budget unreadable, the new task starts anyway", "err", err)
		return nil
	}
	// The 80% half of AC#3: a task that starts anyway still announces the
	// budget line it crossed. Before this publish point the warn Notice was only
	// ever returned (BudgetVerdict.Notice) and never produced onto the bus, so the
	// ticket"s "80% warn event" had no event exit — 44-v1 §3.2 ②.
	if rep.Metered && rep.Verdict.AllowNewTask && rep.Verdict.Level != observe.LevelOK {
		l.publish(Event{Kind: EvWarn, TaskID: taskID, Text: rep.Verdict.Notice})
	}
	if !rep.Metered || rep.Verdict.AllowNewTask {
		return nil
	}
	return &rep.Verdict
}

// SetCostBudget hands the loop's meter the [cost] caps. Until the host calls
// this the meter judges against no cap (uncapped), which is the honest reading
// of a budget nobody configured — see the ticket report's wiring note.
func (l *Loop) SetCostBudget(b observe.BudgetParams) {
	if l.costMeter != nil {
		l.costMeter.SetBudget(b)
	}
}

// CostWiring adapts the host's two C23 knobs into the shapes the loop consumes
// (ticket 44 生产咬合: until a host calls this the price card is a zero value, every
// turn costs 0 micros, and the meter judges against no cap at all — the [cost] section
// and llm.providers.<p>.models.<m>.price were loadable and dead). The price card is
// looked up, never guessed: a model config.toml does not carry a card for returns the
// zero Price, which is PriceKnown() == false and the 价格未知 reading (AC#4), not a
// fabricated number. The caps go through observe.BudgetFrom so a bad [cost] leaf is
// refused there and leaves the loop uncapped with the host naming the reason.
func CostWiring(cfg *config.Config, provider, model string) (config.Price, observe.BudgetParams) {
	var card config.Price
	var budget observe.BudgetParams
	if cfg == nil {
		return card, budget
	}
	if p, ok := cfg.LLM.Providers[provider]; ok {
		if m, ok := p.Models[model]; ok {
			card = m.Price
		}
	}
	b, err := observe.BudgetFrom(cfg.Cost.DailyBudget, cfg.Cost.MonthlyBudget,
		cfg.Cost.AlertThreshold, cfg.Cost.OverBudget)
	if err != nil {
		return card, observe.BudgetParams{}
	}
	return card, b
}
