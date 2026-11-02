package agent

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Ticket 44 AC#3 measured on the loop itself: reaching 100% rejects a NEW task
// before the provider is touched, says why in the user's language, stays in the
// product's own record, leaves a running task alone, and does none of that in
// warn mode.

// flipSink sets the loop's money budget the first time the running task emits a
// text delta - i.e. mid-flight, after the entry gate has already been passed.
// It exists so "running tasks are unaffected" is MEASURED rather than asserted:
// the budget changes underneath a task that is already going.
type flipSink struct {
	mu     sync.Mutex
	inner  *RecordSink
	h      *harness
	params observe.BudgetParams
	flips  int
}

func (s *flipSink) Publish(e Event) {
	if e.Kind == EvTextDelta {
		s.mu.Lock()
		if s.flips == 0 {
			s.flips++
			s.h.loop.SetCostBudget(s.params)
		}
		s.mu.Unlock()
	}
	s.inner.Publish(e)
}

// flipped reports whether the mid-flight budget change actually happened - the
// case is worthless without it, because a task that never published a delta
// would "prove" the rule by never testing it.
func (s *flipSink) flipped() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.flips
}

func withMidFlightBudget(p observe.BudgetParams, out **flipSink) harnessOpt {
	return func(h *harness, o *Options) {
		fs := &flipSink{inner: h.sink, h: h, params: p}
		*out = fs
		o.Sink = fs
	}
}

// TestAC3RunningTaskFinishesPastTheLimit crosses the line while the task runs.
func TestAC3RunningTaskFinishesPastTheLimit(t *testing.T) {
	store := newStore44(t)
	var fs *flipSink
	// One micro of headroom: the task's own spend crosses it on the way.
	cap1 := observe.BudgetParams{DailyMicros: 1, AlertBasisPoints: 8000, PauseAtLimit: true}
	h := newHarness(t, "tool-then-text", withJournal(store),
		withConfig(func(c *Config) { c.Currency = "CNY"; c.Price = price44 }),
		withMidFlightBudget(cap1, &fs))
	res := h.run("天气怎么样")

	if fs.flipped() == 0 {
		t.Fatal("the mid-flight budget change never happened: this case measured nothing")
	}
	if res.Status != StatusCompleted {
		t.Fatalf("status = %s (%s): the money budget stopped a task that was already running",
			res.Status, res.Message)
	}
	if res.CostMicros <= 1 {
		t.Fatalf("setup: the task cost %d micros, which does not cross the 1-micro cap", res.CostMicros)
	}
	rep, err := h.loop.CostBudget(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if rep.Verdict.AllowNewTask {
		t.Error("after the crossing a NEW task should be gated; finishing this one is only half the rule")
	}
}

// TestAC3NewTaskRefusedBeforeTheProvider is the refusal half: no byte leaves the
// machine, the message is explicit, the row is auditable, the two cost surfaces
// still agree.
func TestAC3NewTaskRefusedBeforeTheProvider(t *testing.T) {
	store := newStore44(t)
	first := newHarness(t, "tool-then-text", withJournal(store),
		withConfig(func(c *Config) { c.Currency = "CNY"; c.Price = price44 }))
	res := first.run("先跑一个")
	if res.CostMicros == 0 {
		t.Fatalf("setup produced no spend: %+v", res)
	}

	params, err := observe.BudgetFrom(res.CostMicros, 0, 0.8, "pause")
	if err != nil {
		t.Fatal(err)
	}
	second := newHarness(t, "tool-then-text", withJournal(store),
		withConfig(func(c *Config) { c.Currency = "CNY"; c.Price = price44 }))
	second.loop.SetCostBudget(params)

	refused := second.run("再跑一个")
	if refused.Status != StatusFailed {
		t.Fatalf("status = %s (%s), want failed: 100%% of a pause-mode cap must reject a new task",
			refused.Status, refused.Message)
	}
	if refused.Err == nil || refused.Err.Class != observe.ClassBudget {
		t.Errorf("refusal error = %+v, want class %q so a reader can tell a budget refusal from a provider failure",
			refused.Err, observe.ClassBudget)
	}
	if !strings.Contains(refused.Message, "暂停新任务") {
		t.Errorf("the rejection message is not explicit: %q", refused.Message)
	}
	if second.requests() != 0 {
		t.Errorf("the refused task made %d provider requests, want 0: a budget that gates after the call is a budget that already paid for the task",
			second.requests())
	}

	// The refusal is in the product's record, not only in a log line.
	row, err := store.TaskLogByID(context.Background(), refused.TaskID)
	if err != nil {
		t.Fatalf("the refused task left no task_log row: %v", err)
	}
	if row.ErrorClass != string(observe.ClassBudget) {
		t.Errorf("task_log.error_class = %q, want %q", row.ErrorClass, observe.ClassBudget)
	}
	if !strings.Contains(strings.Join(deref(row.SummaryText), ""), "暂停新任务") {
		t.Errorf("the refusal reason never reached the row: %+v", row)
	}

	// And the two cost surfaces still agree: the refusal booked a task that cost
	// nothing, in both of them.
	m, err := observe.NewCostMeter(store, observe.WithTasks(store), observe.WithLocation(time.Local))
	if err != nil {
		t.Fatal(err)
	}
	rec, err := m.Reconcile(context.Background(), localDay(t))
	if err != nil {
		t.Fatal(err)
	}
	if !rec.Agree {
		t.Errorf("the refusal broke the reconciliation invariant: agg=%+v base=%+v", rec.Aggregate, rec.Base)
	}
	if rec.Aggregate.Tasks != 2 {
		t.Errorf("the day aggregate counts %d tasks, want 2 (the paid one plus the refused one)", rec.Aggregate.Tasks)
	}
	if rec.Aggregate.Micros != res.CostMicros {
		t.Errorf("the refusal added money: %d, want the first task's %d", rec.Aggregate.Micros, res.CostMicros)
	}
}

// TestAC3WarnModeStillRunsTasks pins the alternative the ticket names
// (over_budget=warn): past 100% tasks keep running, and the notice says why.
func TestAC3WarnModeStillRunsTasks(t *testing.T) {
	store := newStore44(t)
	first := newHarness(t, "tool-then-text", withJournal(store),
		withConfig(func(c *Config) { c.Currency = "CNY"; c.Price = price44 }))
	res := first.run("先跑一个")

	params, err := observe.BudgetFrom(res.CostMicros, 0, 0.8, "warn")
	if err != nil {
		t.Fatal(err)
	}
	second := newHarness(t, "tool-then-text", withJournal(store),
		withConfig(func(c *Config) { c.Currency = "CNY"; c.Price = price44 }))
	second.loop.SetCostBudget(params)
	out := second.run("warn 模式继续跑")
	if out.Status != StatusCompleted {
		t.Fatalf("over_budget=warn refused the task (%s): %s", out.Status, out.Message)
	}
	if second.requests() == 0 {
		t.Error("warn mode must still reach the provider")
	}
	if out.CostMicros <= 0 {
		t.Errorf("the warn-mode task booked no cost: %+v", out)
	}
	// Warn-mode still produces the visible event a toast/badge would ride on.
	rep, err := second.loop.CostBudget(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if rep.Verdict.Level != observe.LevelOver || rep.Verdict.Notice == "" {
		t.Errorf("over-cap warn mode produced %+v, want level=over with a notice", rep.Verdict)
	}
	if !strings.Contains(rep.Verdict.Notice, "不拦") && !strings.Contains(rep.Verdict.Notice, "warn") {
		t.Errorf("the notice does not say the mode it is in: %q", rep.Verdict.Notice)
	}
}

func deref(p *string) []string {
	if p == nil {
		return nil
	}
	return []string{*p}
}
