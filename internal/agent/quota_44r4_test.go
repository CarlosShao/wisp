package agent

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Ticket 44 AC#5 wiring tests (leg 44-r4). The disk mouth is asserted against
// the real memory DAO at compile time: if provider_health ever loses
// RecordProviderError, the pre-dispatch signal wiring stops compiling here
// rather than silently going inert.
var _ QuotaSignalWriter = (*memory.Store)(nil)

// countingProvider records whether a chain element was actually dispatched.
type countingProvider struct {
	name  string
	calls int
	err   error
}

func (p *countingProvider) Stream(ctx context.Context, req *llm.Request, emit func(llm.StreamEvent) error) error {
	p.calls++
	return p.err
}
func (p *countingProvider) Info() llm.ProviderInfo { return llm.ProviderInfo{Model: p.name} }

type capturedSignal struct {
	provider, model, class, detail, state string
}

type fakeSignalWriter struct{ got []capturedSignal }

func (w *fakeSignalWriter) RecordProviderError(ctx context.Context, provider, model, class, detail, quotaState string) error {
	w.got = append(w.got, capturedSignal{provider, model, class, detail, quotaState})
	return nil
}

func newTestGate(t *testing.T, now func() time.Time, params observe.QuotaParams,
	cards map[string]observe.QuotaModel, providers map[string]observe.QuotaProvider,
) (*QuotaGate, *[]Event) {
	t.Helper()
	events := &[]Event{}
	g := NewQuotaGate(params, cards, providers, observe.NewQuotaLedger(now),
		func(e Event) { *events = append(*events, e) })
	return g, events
}

// AC#5 headline fixture: remaining 0 → the over-quota model is NOT dispatched
// (zero bytes out), the chain keeps advancing to the next entry, and the skip
// is a visible event, not silence (票 44:53-55, SPEC-03 §3.1 不得静默).
func TestAC5RemainingZeroSkipsDispatchChainAdvances(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC) }
	cards := map[string]observe.QuotaModel{
		"pa/m1": {Provider: "pa", Model: "m1", Billing: observe.BillingPayPerToken, DailyMicros: 1000},
		"pb/m2": {Provider: "pb", Model: "m2", Billing: observe.BillingPayPerToken, DailyMicros: 1000},
	}
	g, events := newTestGate(t, now, observe.QuotaParams{AlertBasisPoints: 8000}, cards, nil)
	if err := g.Book("pa/m1", 1000); err != nil {
		t.Fatalf("book exhausted spend: %v", err)
	}
	first := &countingProvider{name: "pa/m1"}
	second := &countingProvider{name: "pb/m2"}
	// The refusal happens BEFORE the element's Stream: the classified budget
	// error is exactly what ticket 11's ChainRunner steps over (its own test
	// feeds a ClassBudget element into the chain and moves to the next entry).
	err := g.GateProvider(first, "pa/m1").Stream(context.Background(), &llm.Request{}, func(llm.StreamEvent) error { return nil })
	if err == nil {
		t.Fatal("remaining 0 must refuse the dispatch (anti-negative pre-dispatch check)")
	}
	class, _ := observe.ClassOf(err)
	if class != observe.ClassBudget {
		t.Errorf("refusal class = %v, want budget (the class the chain steps on)", class)
	}
	if first.calls != 0 {
		t.Errorf("over-quota model was dispatched %d times, want 0 (never into debt)", first.calls)
	}
	// Chain advances: next entry dispatches normally.
	if err := g.GateProvider(second, "pb/m2").Stream(context.Background(), &llm.Request{}, func(llm.StreamEvent) error { return nil }); err != nil {
		t.Fatalf("fallback element must still dispatch: %v", err)
	}
	if second.calls != 1 {
		t.Errorf("fallback calls = %d, want 1 (chain continued)", second.calls)
	}
	chosen, skips := g.Chain([]string{"pa/m1", "pb/m2"})
	if chosen != "pb/m2" || len(skips) != 1 || skips[0].Model != "pa/m1" {
		t.Errorf("chain reading chosen=%q skips=%v, want pb/m2 with pa/m1 skipped", chosen, skips)
	}
	// Visible event with the user-facing reason, not silence.
	if len(*events) == 0 || (*events)[0].Kind != EvWarn || !strings.Contains((*events)[0].Text, "跳过") {
		t.Errorf("skip produced no visible event, got %v", *events)
	}
}

// A fully exhausted chain returns chosen == "" with every entry named: the
// gate says so and steps aside — ending the task is ticket 11's runner.
func TestAC5WholeChainExhaustedStillAdvancesNotBreaks(t *testing.T) {
	cards := map[string]observe.QuotaModel{
		"pa/m1": {Provider: "pa", Model: "m1", DailyMicros: 500},
		"pb/m2": {Provider: "pb", Model: "m2", DailyMicros: 500},
	}
	g, _ := newTestGate(t, nil, observe.QuotaParams{AlertBasisPoints: 8000}, cards, nil)
	for _, k := range []string{"pa/m1", "pb/m2"} {
		if err := g.Book(k, 500); err != nil {
			t.Fatalf("book %s: %v", k, err)
		}
	}
	chosen, skips := g.Chain([]string{"pa/m1", "pb/m2"})
	if chosen != "" || len(skips) != 2 {
		t.Fatalf("exhausted chain = %q / %v, want empty choice with both skipped", chosen, skips)
	}
}

// plan_credit_total_micro 分账: manual total minus booked consumption, the
// remaining carries the 估算 marker, and a plan draw never darkens a
// pay-per-token model's number (AC#5 "tracked separately per model").
func TestAC5PlanCreditBalanceEstimatedAndSeparateFromPay(t *testing.T) {
	cards := map[string]observe.QuotaModel{
		"pl/m1": {Provider: "pl", Model: "m1", Billing: observe.BillingPlan},
		"pl/m2": {Provider: "pl", Model: "m2", Billing: observe.BillingPayPerToken},
	}
	provs := map[string]observe.QuotaProvider{
		"pl": {Provider: "pl", PlanCreditTotalMicro: 1000},
	}
	g, events := newTestGate(t, nil, observe.QuotaParams{AlertBasisPoints: 8000}, cards, provs)
	if err := g.Book("pl/m1", 900); err != nil {
		t.Fatalf("book plan spend: %v", err)
	}
	v := g.Decide("pl/m1")
	if v.RemainingMicros != 100 {
		t.Errorf("plan remaining = %d, want 100 (manual total 1000 - booked 900)", v.RemainingMicros)
	}
	if !v.Estimated || v.Binding != "plan-credit" {
		t.Errorf("plan verdict = %+v, want estimated plan-credit binding", v)
	}
	if v.Level != observe.LevelWarn {
		t.Errorf("90%% of plan credit level = %q, want warn", v.Level)
	}
	if len(*events) != 0 {
		t.Errorf("decide must not publish; got %v", *events)
	}
	// The pay-per-token sibling is untouched by the plan draw: tracked
	// separately by construction.
	pv := g.Decide("pl/m2")
	if pv.Level != observe.LevelOK || pv.RemainingMicros != 0 || pv.Binding != "" {
		t.Errorf("pay-per-token sibling verdict = %+v, want uncapped OK untouched by plan spend", pv)
	}
	// Crossing the plan line hard-stops ONLY the plan model.
	if err := g.Book("pl/m1", 100); err != nil {
		t.Fatalf("book rest of plan: %v", err)
	}
	plan := &countingProvider{name: "pl/m1"}
	err := g.GateProvider(plan, "pl/m1").Stream(context.Background(), &llm.Request{}, func(llm.StreamEvent) error { return nil })
	if err == nil || plan.calls != 0 {
		t.Fatalf("exhausted plan credit must skip dispatch (err=%v calls=%d)", err, plan.calls)
	}
	if !strings.Contains(err.Error(), observe.PlanCreditEstimatedMarker) {
		t.Errorf("stop notice must keep the 估算 marker, got %q", err.Error())
	}
	pay := &countingProvider{name: "pl/m2"}
	if err := g.GateProvider(pay, "pl/m2").Stream(context.Background(), &llm.Request{}, func(llm.StreamEvent) error { return nil }); err != nil || pay.calls != 1 {
		t.Errorf("pay-per-token sibling must still dispatch: err=%v calls=%d", err, pay.calls)
	}
}

// Anti-negative bookkeeping: a ledger that credited negatives is one
// subtraction from fabricating remaining.
func TestAC5NegativeSpendRefusedLedgerUnchanged(t *testing.T) {
	cards := map[string]observe.QuotaModel{
		"pa/m1": {Provider: "pa", Model: "m1", DailyMicros: 1000},
	}
	g, _ := newTestGate(t, nil, observe.QuotaParams{AlertBasisPoints: 8000}, cards, nil)
	if err := g.Book("pa/m1", 400); err != nil {
		t.Fatalf("book: %v", err)
	}
	if err := g.Book("pa/m1", -400); err == nil {
		t.Fatal("negative spend accepted into the quota ledger")
	}
	if v := g.Decide("pa/m1"); v.RemainingMicros != 600 {
		t.Errorf("remaining after refused negative = %d, want 600 (no fabricated credit)", v.RemainingMicros)
	}
}

// The 80% half of the quota layers has an event exit too (same discipline as
// r3's [cost] warn fix): warn dispatches anyway, but announces first.
func TestAC5QuotaWarnPublishesVisibleEventAndStillDispatches(t *testing.T) {
	cards := map[string]observe.QuotaModel{
		"pa/m1": {Provider: "pa", Model: "m1", DailyMicros: 1000},
	}
	g, events := newTestGate(t, nil, observe.QuotaParams{AlertBasisPoints: 8000}, cards, nil)
	if err := g.Book("pa/m1", 800); err != nil {
		t.Fatalf("book: %v", err)
	}
	p := &countingProvider{name: "pa/m1"}
	if err := g.GateProvider(p, "pa/m1").Stream(context.Background(), &llm.Request{}, func(llm.StreamEvent) error { return nil }); err != nil || p.calls != 1 {
		t.Fatalf("warn must still dispatch: err=%v calls=%d", err, p.calls)
	}
	if len(*events) != 1 || (*events)[0].Kind != EvWarn || !strings.Contains((*events)[0].Text, "80%") {
		t.Errorf("warn event missing or wrong: %v", *events)
	}
}

// Disk-side bookkeeping signal: exhaustion lands in provider_health via the
// existing DAO, and a signal failure never blocks the refusal.
func TestAC5ExhaustedSignalWritesProviderHealthQuotaState(t *testing.T) {
	cards := map[string]observe.QuotaModel{
		"pa/m1": {Provider: "pa", Model: "m1", DailyMicros: 10},
	}
	g, _ := newTestGate(t, nil, observe.QuotaParams{AlertBasisPoints: 8000}, cards, nil)
	w := &fakeSignalWriter{}
	g.SetSignalWriter(w, func(string, ...any) {})
	if err := g.Book("pa/m1", 10); err != nil {
		t.Fatalf("book: %v", err)
	}
	p := &countingProvider{}
	if err := g.GateProvider(p, "pa/m1").Stream(context.Background(), &llm.Request{}, func(llm.StreamEvent) error { return nil }); err == nil {
		t.Fatal("want refusal at remaining 0")
	}
	if len(w.got) != 1 || w.got[0].provider != "pa" || w.got[0].model != "m1" ||
		w.got[0].class != "budget" || w.got[0].state != string(memory.QuotaExhausted) {
		t.Errorf("provider_health signal = %+v, want pa/m1 budget exhausted", w.got)
	}
}

// QuotaWiring on an empty config says "no quota anywhere" — it invents no
// caps and no alert drift (the CostWiring reading discipline).
func TestAC5QuotaWiringEmptyConfigInventsNothing(t *testing.T) {
	cards, provs, params := QuotaWiring(&config.Config{})
	if len(cards) != 0 || len(provs) != 0 {
		t.Errorf("empty config produced quota cards: %v / %v", cards, provs)
	}
	if params.AlertBasisPoints != 0 {
		t.Errorf("empty config alert = %d, want 0 (no invented line)", params.AlertBasisPoints)
	}
	cards, provs, params = QuotaWiring(nil)
	if len(cards) != 0 || len(provs) != 0 || params.AlertBasisPoints != 0 {
		t.Errorf("nil config must produce the zero shape, got %v/%v/%v", cards, provs, params)
	}
}
