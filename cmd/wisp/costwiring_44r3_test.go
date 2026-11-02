package main

import (
	"os"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Ticket 44 生产咬合 (44-v1 §3.2 ③ / §7 ①): `observe.BudgetFrom` and
// `Loop.SetCostBudget` had ZERO callers outside tests, and `agent.Config.Price`
// was never filled by the composition root — so in a real `wisp run` every turn
// priced at 0 micros and the meter judged against no cap, however much money
// config.toml said to spend. The code path was right and unreachable.
//
// These two cases close it from both ends: the adapter hands the loop the real
// [cost] caps and the model's real price card (and says 价格未知 rather than
// inventing a number when the config carries no card), and run.go provably
// applies them on every task.

func cfg44r3(cost config.CostSection) *config.Config {
	c := &config.Config{}
	c.LLM.Providers = map[string]config.Provider{"p44r3": {
		Models: map[string]config.ModelSpec{"m44r3": {
			Price: config.Price{In: 2_000_000, Out: 10_000_000, Cached: 1_000_000},
		}},
	}}
	c.Cost = cost
	return c
}

// Test44R3CostWiringCarriesPriceAndCaps is the positive reading.
func Test44R3CostWiringCarriesPriceAndCaps(t *testing.T) {
	c := cfg44r3(config.CostSection{
		DailyBudget: 5_000_000, MonthlyBudget: 90_000_000,
		AlertThreshold: 0.8, OverBudget: "pause",
	})
	card, budget := agent.CostWiring(c, "p44r3", "m44r3")

	if card.In != 2_000_000 || card.Out != 10_000_000 || card.Cached != 1_000_000 {
		t.Fatalf("price card = %+v, want the model's own [llm.providers.p44r3.models.m44r3.price] card", card)
	}
	if got := observe.PriceMicros(agent.CardOf(card), observe.Tokens{In: 1_000, Out: 500}); got != 7_000 {
		t.Errorf("the card the host hands the loop prices 1000/500 at %d micros, want 2*1000+10*500 = 7000", got)
	}
	if budget.DailyMicros != 5_000_000 || budget.MonthlyMicros != 90_000_000 {
		t.Errorf("budget caps = %+v, want the [cost] leaves carried across", budget)
	}
	if budget.AlertBasisPoints != 8000 || !budget.PauseAtLimit {
		t.Errorf("budget knobs = %+v, want 0.8 to 8000 bp and over_budget=pause to mean PauseAtLimit", budget)
	}
}

// Test44R3CostWiringSaysUnknownPriceAndBadBudget is the reverse half: an
// uncarded model yields a zero card (成本 0 可见, and Priced=false downstream —
// AC#4's rule, not a guessed price), and a bad [cost] leaf yields uncapped
// params instead of a silently clamped policy.
func Test44R3CostWiringSaysUnknownPriceAndBadBudget(t *testing.T) {
	c := cfg44r3(config.CostSection{DailyBudget: 1, AlertThreshold: 0.8, OverBudget: "pause"})

	card, budget := agent.CostWiring(c, "p44r3", "model-with-no-card")
	if agent.CardOf(card).Known() {
		t.Errorf("a model the config has no card for came back priced: %+v", card)
	}
	if got := observe.PriceMicros(agent.CardOf(card), observe.Tokens{In: 1_000_000, Out: 1_000_000}); got != 0 {
		t.Errorf("the uncarded model priced at %d micros, want 0 with the 价格未知 marker", got)
	}
	if budget.DailyMicros != 1 {
		t.Errorf("a missing card must not drop the caps: %+v", budget)
	}

	// alert_threshold outside [0,1] is a config bug: BudgetFrom refuses it and
	// the loop keeps judging against no cap rather than against a clamped one.
	bad := cfg44r3(config.CostSection{DailyBudget: 5_000_000, AlertThreshold: 1.5, OverBudget: "pause"})
	if _, b := agent.CostWiring(bad, "p44r3", "m44r3"); b != (observe.BudgetParams{}) {
		t.Errorf("bad [cost] produced %+v, want the zero BudgetParams (uncapped, refused upstream)", b)
	}
	// Unknown over_budget enum likewise.
	worse := cfg44r3(config.CostSection{DailyBudget: 5_000_000, AlertThreshold: 0.8, OverBudget: "stop"})
	if _, b := agent.CostWiring(worse, "p44r3", "m44r3"); b.DailyMicros != 0 {
		t.Errorf("over_budget=stop produced %+v, want uncapped", b)
	}
}

// Test44R3RunGoAppliesTheWiring is the 咬合 pin in the composition root: the
// adapter alone would be another dead function, so the file that assembles the
// loop has to name the price field, hand over the caps after the loop exists,
// and read the warn event the loop now publishes (44-v1's own ruler was a
// `git grep` over this file, so this pins the same thing from inside the test).
func Test44R3RunGoAppliesTheWiring(t *testing.T) {
	raw, err := os.ReadFile("run.go")
	if err != nil {
		t.Fatal(err)
	}
	src := string(raw)
	for _, want := range []string{
		"card, budget := agent.CostWiring(cfg, ep.Provider, ep.Model)",
		"Price: card,",
		"loop.SetCostBudget(budget)",
		"case agent.EvWarn:",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("run.go does not contain %q: the C23 knobs are wired in the adapter but not in the composition root", want)
		}
	}
	// Order matters: the caps go to the meter after the loop is assembled, and
	// the price card travels in the Options the loop is built from.
	iOpts := strings.Index(src, "card, budget := agent.CostWiring")
	iLoop := strings.Index(src, "loop, err := agent.New(loopOpts)")
	iSet := strings.Index(src, "loop.SetCostBudget(budget)")
	if !(iOpts >= 0 && iOpts < iLoop && iLoop < iSet) {
		t.Errorf("wiring order wrong: CostWiring=%d New=%d SetCostBudget=%d", iOpts, iLoop, iSet)
	}
	// The live-config half: the call sits after the per-task re-read, which is
	// what makes ticket 223's reload tick reach [cost] on the next task.
	iCfg := strings.Index(src, "cfg = rt.mgr.Config()")
	if !(iCfg >= 0 && iCfg < iOpts) {
		t.Errorf("CostWiring is not reading the live per-task config: mgr.Config=%d CostWiring=%d", iCfg, iOpts)
	}
}
