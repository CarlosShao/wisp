package main

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Ticket 44-r5 装配腿: the pre-dispatch quota gate was delivered at the
// component layer (internal/agent/quota_44r4.go) but 44-v2-verdict §4①(a)
// registered that cmd/wisp had ZERO callers of QuotaWiring/NewQuotaGate/
// GateProvider - so in a real `wisp run` an over-quota model was still
// dispatched into debt. These cases pin the assembly from three sides:
// the behaviour (over-quota element skipped, chain advanced, event visible),
// the reverse judge (no gate installed == the debt dispatch), and the
// composition root (the wiring is actually in run.go, in the right order).

// stub44r5 is a chain element that counts dispatches. A call here means bytes
// left the machine for that model, which is the thing AC#5 forbids at zero
// remaining.
type stub44r5 struct {
	name    string
	calls   int
	serr    error
	info    llm.ProviderInfo
	wroteIt bool
}

func (s *stub44r5) Info() llm.ProviderInfo { return s.info }

func (s *stub44r5) Stream(_ context.Context, _ *llm.Request, _ func(llm.StreamEvent) error) error {
	s.calls++
	return s.serr
}

func cfg44r5() *config.Config {
	c := &config.Config{}
	c.LLM.Providers = map[string]config.Provider{
		"pa": {Billing: observe.BillingPayPerToken, Models: map[string]config.ModelSpec{
			// m1 is capped hard at 1000 micros/day, m2 at 5000.
			"m1": {QuotaDailyMicro: 1000},
			"m2": {QuotaDailyMicro: 5000},
		}},
	}
	c.Cost = config.CostSection{
		DailyBudget: 5_000_000, MonthlyBudget: 90_000_000,
		AlertThreshold: 0.8, OverBudget: "pause",
	}
	return c
}

// Test44R5OverQuotaElementIsSkippedAndChainAdvances is the ticket's whole
// sentence in the host shape: pa/m1 has hit its daily cap, so it is NEVER
// dispatched, the chain steps to pa/m2 and the task still gets an answer, and
// both the skip and the step come out as visible warn events.
func Test44R5OverQuotaElementIsSkippedAndChainAdvances(t *testing.T) {
	first := &stub44r5{name: "pa/m1", info: llm.ProviderInfo{Model: "m1"}}
	second := &stub44r5{name: "pa/m2", info: llm.ProviderInfo{Model: "m2"}}
	var events []agent.Event
	notify := func(e agent.Event) { events = append(events, e) }

	prov, gate, _, servedKey := quotaDispatchChain(cfg44r5(),
		[]llm.LlmProvider{first, second}, []string{"pa/m1", "pa/m2"}, nil, notify)
	if gate == nil {
		t.Fatal("no quota gate assembled even though an event mouth was handed in")
	}
	if err := gate.Book("pa/m1", 1000); err != nil {
		t.Fatalf("booking the cap: %v", err)
	}

	if err := prov.Stream(context.Background(), &llm.Request{}, func(llm.StreamEvent) error { return nil }); err != nil {
		t.Fatalf("the chain should have advanced to pa/m2 and succeeded, got %v", err)
	}
	if first.calls != 0 {
		t.Errorf("over-quota pa/m1 was dispatched %d time(s): 票 44:34-35 says it is SKIPPED, never dispatched into debt", first.calls)
	}
	if second.calls != 1 {
		t.Errorf("pa/m2 dispatched %d time(s), want 1: the chain has to advance, not break", second.calls)
	}
	if got := servedKey(); got != "pa/m2" {
		t.Errorf("servedKey() = %q, want the element that actually served (the spend must be booked to the right bucket)", got)
	}
	var skip, step bool
	for _, e := range events {
		if e.Kind != agent.EvWarn {
			continue
		}
		if strings.Contains(e.Text, "跳过") && strings.Contains(e.Text, "pa/m1") {
			skip = true
		}
		if strings.Contains(e.Text, "派发链推进") && strings.Contains(e.Text, "pa/m2") {
			step = true
		}
	}
	if !skip {
		t.Errorf("no visible skip event, events = %+v (SPEC-03 §3.1 不得静默)", events)
	}
	if !step {
		t.Errorf("no visible chain-advance event, events = %+v", events)
	}
}

// Test44R5NoEventMouthMeansNoGateIsTheDebtShape is the reverse judge as a
// behaviour: an assembly that installs no gate hands back the raw primary
// element, which dispatches an over-quota model into debt. This is the shape
// run.go was in before 44-r5; if the wiring is removed the pin below goes red.
func Test44R5NoEventMouthMeansNoGateIsTheDebtShape(t *testing.T) {
	first := &stub44r5{name: "pa/m1", info: llm.ProviderInfo{Model: "m1"}}
	prov, gate, _, servedKey := quotaDispatchChain(cfg44r5(),
		[]llm.LlmProvider{first}, []string{"pa/m1"}, nil, nil)
	if gate != nil {
		t.Fatal("a gate with no event mouth would announce nothing - it must not be installed")
	}
	if err := prov.Stream(context.Background(), &llm.Request{}, func(llm.StreamEvent) error { return nil }); err != nil {
		t.Fatalf("ungated shape must dispatch straight through, got %v", err)
	}
	if first.calls != 1 {
		t.Errorf("calls = %d, want 1: this is 派进债务, and it is only ever reachable when nothing is wired", first.calls)
	}
	if servedKey() != "pa/m1" {
		t.Errorf("servedKey() = %q, want the primary key", servedKey())
	}
}

// Test44R5UnderCapChainDispatchesThePrimaryNormally is the positive control:
// with an empty ledger nothing is skipped and no warning is produced, so the
// gate cannot be a machine that blocks real work.
func Test44R5UnderCapChainDispatchesThePrimaryNormally(t *testing.T) {
	first := &stub44r5{name: "pa/m1", info: llm.ProviderInfo{Model: "m1"}}
	second := &stub44r5{name: "pa/m2", info: llm.ProviderInfo{Model: "m2"}}
	var events []agent.Event
	prov, _, _, servedKey := quotaDispatchChain(cfg44r5(),
		[]llm.LlmProvider{first, second}, []string{"pa/m1", "pa/m2"}, nil,
		func(e agent.Event) { events = append(events, e) })

	if err := prov.Stream(context.Background(), &llm.Request{}, func(llm.StreamEvent) error { return nil }); err != nil {
		t.Fatalf("an under-cap chain must not fail: %v", err)
	}
	if first.calls != 1 || second.calls != 0 {
		t.Errorf("dispatch split = %d/%d, want 1/0 (no card tripped)", first.calls, second.calls)
	}
	if len(events) != 0 {
		t.Errorf("events = %+v, want none below the alert line", events)
	}
	if servedKey() != "pa/m1" {
		t.Errorf("servedKey() = %q, want pa/m1", servedKey())
	}
}

// Test44R5UncappedModelIsNeverBlocked covers 44-r4's no-card rule through the
// host wiring: a config that sets no quota anywhere dispatches as before.
func Test44R5UncappedModelIsNeverBlocked(t *testing.T) {
	c := cfg44r5()
	c.LLM.Providers["pa"] = config.Provider{Models: map[string]config.ModelSpec{"m1": {}}}
	first := &stub44r5{name: "pa/m1", info: llm.ProviderInfo{Model: "m1"}}
	prov, gate, _, _ := quotaDispatchChain(c, []llm.LlmProvider{first}, []string{"pa/m1"}, nil,
		func(agent.Event) {})
	if err := prov.Stream(context.Background(), &llm.Request{}, func(llm.StreamEvent) error { return nil }); err != nil || first.calls != 1 {
		t.Errorf("a model with no quota card was blocked (calls=%d err=%v)", first.calls, err)
	}
	if gate == nil {
		t.Error("gate handle missing: run.go needs it to book spend and to write quota_state")
	}
}

// Test44R5RunGoInstallsTheGate is the 咬合 pin in the composition root, the
// same ruler 44-v2 used (a git grep over non-test cmd/wisp) but inside a test:
// remove the assembly and this goes red.
func Test44R5RunGoInstallsTheGate(t *testing.T) {
	raw, err := os.ReadFile("run.go")
	if err != nil {
		t.Fatal(err)
	}
	src := string(raw)
	for _, want := range []string{
		"prov, quotaGate, ledger, servedKey := quotaDispatchChain(cfg, rt.provs, quotaKeys, rt.quotaLedger, sink.Publish)",
		"rt.quotaLedger = ledger",
		"quotaGate.SetSignalWriter(rt.store, rt.auditf)",
		"Provider: prov,",
		"Sink:     sink,",
		"quotaGate.Book(served, res.CostMicros)",
		"agent.QuotaWiring(cfg)",
		"agent.NewQuotaGate(params, cards, providers, ledger, notify)",
		"gate.GateProvider(p, k)",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("run.go does not contain %q: the quota gate is a component again, not a dispatch path", want)
		}
	}
	// Order: the live per-task config is read first, then the gate is assembled,
	// then the loop is built from it, and the booking happens after the task ran.
	iCfg := strings.Index(src, "cfg = rt.mgr.Config()")
	iGate := strings.Index(src, ":= quotaDispatchChain(")
	iNew := strings.Index(src, "loop, err := agent.New(loopOpts)")
	iWait := strings.Index(src, "res := bg.Wait()")
	iBook := strings.Index(src, "quotaGate.Book(served, res.CostMicros)")
	if !(iCfg >= 0 && iCfg < iGate && iGate < iNew) {
		t.Errorf("assembly order wrong: live cfg=%d gate=%d agent.New=%d", iCfg, iGate, iNew)
	}
	if !(iWait >= 0 && iWait < iBook) {
		t.Errorf("booking is not on the finished task: wait=%d book=%d", iWait, iBook)
	}
	// The card key and the gate key are the same endpoint pair, so a price and a
	// quota decision cannot name two different models.
	if !strings.Contains(src, "card, budget := agent.CostWiring(cfg, ep.Provider, ep.Model)") ||
		!strings.Contains(src, `quotaKeys[0] = ep.Provider + "/" + ep.Model`) {
		t.Error("run.go keys the price card and the quota gate from different places")
	}
}

// Test44R5ChainRunnerIsNotReimplemented keeps this leg inside its own writing
// face: the host assembles the existing chain surface and does NOT grow its own
// stepping logic into internal/llm or internal/agent.
func Test44R5ChainRunnerIsNotReimplemented(t *testing.T) {
	raw, err := os.ReadFile("run.go")
	if err != nil {
		t.Fatal(err)
	}
	src := string(raw)
	helper := src[strings.Index(src, "func quotaDispatchChain("):]
	helper = helper[:strings.Index(helper, "\n// gatedChainDispatch")]
	for _, forbidden := range []string{"for i, p := range r.Elements", "flushTerminals", "ClassBudget", "delivered"} {
		if strings.Contains(helper, forbidden) {
			t.Errorf("the assembly re-implements chain stepping (%q) instead of using llm.ChainRunner", forbidden)
		}
	}
	if !strings.Contains(helper, "llm.ChainRunner{") {
		t.Error("the assembly does not dispatch through ticket 11's llm.ChainRunner")
	}
}
