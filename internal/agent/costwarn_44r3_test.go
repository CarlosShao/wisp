package agent

import (
	"context"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Ticket 44 AC#3, the half 44-v1 named as missing (§3.2 ② / §7 ②): the 80% line
// produced a BudgetVerdict{Level:"warn", Notice:"…"} that was true, tested, and
// never published — the ruler was `git grep -n "publish(Event" -- internal/agent`,
// which found exactly one cost-path event (loop.go's EvError at the 100% refusal)
// and nothing at all on a warn verdict, because AllowNewTask=true made
// preDispatchBudget return nil before any bus write.
//
// These cases pin the event exit itself, not the Notice: a warn verdict must put
// an EvWarn on the sink while the task still runs, the 100% refusal must keep
// riding EvError (and must not double-fire a warn), and an under-threshold run
// must publish NOTHING (so a producer that warned on every task would not pass).

// TestAC3WarnBudgetPublishesWarnEvent puts a task at exactly the 80% line: the
// day already spent res.CostMicros, the cap is 1.25x that, so the basis points
// land on 8000 = warn, in pause mode where warn does not gate.
func TestAC3WarnBudgetPublishesWarnEvent(t *testing.T) {
	store := newStore44(t)
	first := newHarness(t, "tool-then-text", withJournal(store),
		withConfig(func(c *Config) { c.Currency = "CNY"; c.Price = price44 }))
	spend := first.run("先跑一个")
	if spend.CostMicros == 0 {
		t.Fatalf("setup produced no spend: %+v", spend)
	}

	cap1 := spend.CostMicros + spend.CostMicros/4 // 1.25x: the day sits at 80%
	params, err := observe.BudgetFrom(cap1, 0, 0.8, "pause")
	if err != nil {
		t.Fatal(err)
	}
	second := newHarness(t, "tool-then-text", withJournal(store),
		withConfig(func(c *Config) { c.Currency = "CNY"; c.Price = price44 }))
	second.loop.SetCostBudget(params)

	// The verdict really is at warn before the task starts, or this case tests
	// an ok run and the event assertion below would be worthless.
	pre, err := second.loop.CostBudget(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if pre.Verdict.Level != observe.LevelWarn || !pre.Verdict.AllowNewTask {
		t.Fatalf("setup: verdict = %+v, want level=warn with new tasks still allowed", pre.Verdict)
	}

	out := second.run("80% 预警下照常开工")
	if out.Status != StatusCompleted {
		t.Fatalf("status = %s (%s): a warn-level budget must not refuse the task", out.Status, out.Message)
	}
	if second.requests() == 0 {
		t.Error("the warn-level task never reached the provider")
	}

	ev := second.sink.LastOf(EvWarn)
	if ev == nil {
		t.Fatal("no EvWarn on the bus at the 80% line: this is the missing event exit 44-v1 §3.2 ② named")
	}
	if ev.TaskID != out.TaskID {
		t.Errorf("the warn event names task %q, want the warned task %q", ev.TaskID, out.TaskID)
	}
	if !strings.Contains(ev.Text, "预算") || !strings.Contains(ev.Text, "80%") {
		t.Errorf("the warn event text is not the budget notice the user needs: %q", ev.Text)
	}
	if ev.Text != pre.Verdict.Notice {
		t.Errorf("the event carries %q, want the verdict Notice verbatim %q", ev.Text, pre.Verdict.Notice)
	}
	// And it is a warning, not a failure: the refusal kind stays unused here.
	if second.sink.LastOf(EvError) != nil {
		t.Error("a warn-level run also published EvError: the two kinds must not be conflated")
	}
}

// TestAC3WarnEventFiresPastTheLimitInWarnMode covers the second unannounced
// shape: over_budget=warn at 100% keeps starting tasks, and until this publish
// point that crossing was likewise silent on the bus.
func TestAC3WarnEventFiresPastTheLimitInWarnMode(t *testing.T) {
	store := newStore44(t)
	first := newHarness(t, "tool-then-text", withJournal(store),
		withConfig(func(c *Config) { c.Currency = "CNY"; c.Price = price44 }))
	spend := first.run("先跑一个")

	params, err := observe.BudgetFrom(spend.CostMicros, 0, 0.8, "warn")
	if err != nil {
		t.Fatal(err)
	}
	second := newHarness(t, "tool-then-text", withJournal(store),
		withConfig(func(c *Config) { c.Currency = "CNY"; c.Price = price44 }))
	second.loop.SetCostBudget(params)

	out := second.run("warn 模式越过 100%")
	if out.Status != StatusCompleted {
		t.Fatalf("over_budget=warn refused the task: %s %s", out.Status, out.Message)
	}
	ev := second.sink.LastOf(EvWarn)
	if ev == nil {
		t.Fatal("crossing 100% in warn mode published no event at all")
	}
	if !strings.Contains(ev.Text, "上限") {
		t.Errorf("the event does not say the limit was reached: %q", ev.Text)
	}
}

// TestAC3UnderTheAlertLinePublishesNothing is the reverse half of the ruler: an
// always-on warn producer would satisfy the two cases above and still spam the
// user, so ok-level runs must be silent.
func TestAC3UnderTheAlertLinePublishesNothing(t *testing.T) {
	store := newStore44(t)
	first := newHarness(t, "tool-then-text", withJournal(store),
		withConfig(func(c *Config) { c.Currency = "CNY"; c.Price = price44 }))
	spend := first.run("先跑一个")

	// A cap 100x the day's spend: 1% used, far under the 80% line.
	params, err := observe.BudgetFrom(spend.CostMicros*100, 0, 0.8, "pause")
	if err != nil {
		t.Fatal(err)
	}
	second := newHarness(t, "tool-then-text", withJournal(store),
		withConfig(func(c *Config) { c.Currency = "CNY"; c.Price = price44 }))
	second.loop.SetCostBudget(params)

	if rep, err := second.loop.CostBudget(context.Background()); err != nil || rep.Verdict.Level != observe.LevelOK {
		t.Fatalf("setup: verdict = %+v err = %v, want level=ok", rep.Verdict, err)
	}
	if out := second.run("预算充裕的任务"); out.Status != StatusCompleted {
		t.Fatalf("status = %s (%s)", out.Status, out.Message)
	}
	if ev := second.sink.LastOf(EvWarn); ev != nil {
		t.Errorf("published a warning at 1%% of the budget: %q", ev.Text)
	}
}

// TestAC3RefusalStaysAnError pins that the 100% pause-mode refusal was not
// moved off EvError by this landing (the refusal is a failure of the request,
// the warn is a notice about money).
func TestAC3RefusalStaysAnError(t *testing.T) {
	store := newStore44(t)
	first := newHarness(t, "tool-then-text", withJournal(store),
		withConfig(func(c *Config) { c.Currency = "CNY"; c.Price = price44 }))
	spend := first.run("先跑一个")

	params, err := observe.BudgetFrom(spend.CostMicros, 0, 0.8, "pause")
	if err != nil {
		t.Fatal(err)
	}
	second := newHarness(t, "tool-then-text", withJournal(store),
		withConfig(func(c *Config) { c.Currency = "CNY"; c.Price = price44 }))
	second.loop.SetCostBudget(params)

	refused := second.run("该拒的任务")
	if refused.Status != StatusFailed {
		t.Fatalf("status = %s, want failed", refused.Status)
	}
	if second.sink.LastOf(EvError) == nil {
		t.Error("the refusal published no EvError")
	}
	if second.sink.LastOf(EvWarn) != nil {
		t.Error("the refusal published a warn event too: a refused task is announced once, as an error")
	}
}
