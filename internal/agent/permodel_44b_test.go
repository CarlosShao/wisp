package agent

import (
	"context"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Ticket 44b agent-layer acceptance (AC#3): the pre-dispatch quota gate reads
// the same per-model numbers after a "restart" as it did before one, rebuilt
// from the cost_daily_model disk side. The fake below is the durable half:
// both gates share it, so closing gate A and opening gate B is the restart.

type t44bStore struct {
	rows []observe.CostModelDayTotal
}

func (s *t44bStore) BumpCostDayModel(_ context.Context, day, provider, model string, tokensIn, tokensOut, micros int64) error {
	if tokensIn < 0 || tokensOut < 0 || micros < 0 {
		panic("fake store must not receive negative bookings (observe guards upstream)")
	}
	for i := range s.rows {
		r := &s.rows[i]
		if r.Day == day && r.Provider == provider && r.Model == model {
			r.TokensIn += tokensIn
			r.TokensOut += tokensOut
			r.Micros += micros
			r.Tasks++
			return nil
		}
	}
	s.rows = append(s.rows, observe.CostModelDayTotal{
		Day: day, Provider: provider, Model: model,
		TokensIn: tokensIn, TokensOut: tokensOut, Micros: micros, Tasks: 1,
	})
	return nil
}

func (s *t44bStore) CostModelDaysBetween(_ context.Context, fromDay, toDay string) ([]observe.CostModelDayTotal, error) {
	var out []observe.CostModelDayTotal
	for _, r := range s.rows {
		if r.Day >= fromDay && r.Day <= toDay {
			out = append(out, r)
		}
	}
	return out, nil
}

func t44bCards() map[string]observe.QuotaModel {
	return map[string]observe.QuotaModel{
		"pa/m1": {Provider: "pa", Model: "m1", Billing: observe.BillingPayPerToken, DailyMicros: 1000, MonthlyMicros: 5000},
		"pl/m9": {Provider: "pl", Model: "m9", Billing: observe.BillingPlan},
	}
}

func newT44bGate(clock func() time.Time, store *t44bStore) *QuotaGate {
	g := NewQuotaGate(observe.QuotaParams{AlertBasisPoints: 8000}, t44bCards(),
		map[string]observe.QuotaProvider{"pl": {Provider: "pl", PlanCreditTotalMicro: 100000}},
		observe.NewQuotaLedger(clock), nil)
	g.SetModelStore(store, store, clock)
	return g
}

// TestBookModelPersistsAndGates is AC#2+AC#5 wiring shape: one BookModel call
// lands one same-batch row set on the disk side and moves the session ledger.
func TestBookModelPersistsAndGates(t *testing.T) {
	now := time.Date(2026, 11, 3, 9, 0, 0, 0, time.UTC)
	store := &t44bStore{}
	g := newT44bGate(func() time.Time { return now }, store)

	if err := g.BookModel(context.Background(), "pa/m1", 100, 40, 300); err != nil {
		t.Fatalf("BookModel: %v", err)
	}
	if len(store.rows) != 1 {
		t.Fatalf("store rows = %d, want 1", len(store.rows))
	}
	r := store.rows[0]
	if r.Day != "2026-11-03" || r.Provider != "pa" || r.Model != "m1" ||
		r.TokensIn != 100 || r.TokensOut != 40 || r.Micros != 300 || r.Tasks != 1 {
		t.Errorf("stored row %+v does not match the task", r)
	}
	sp := g.Ledger().Spend()
	if sp.ModelDay["pa/m1"] != 300 {
		t.Errorf("session ledger not moved by BookModel: %+v", sp.ModelDay)
	}

	// The gate now says over at the cap using the booked number. The second
	// booking merges into the same (day,provider,model) row — that is the
	// upsert dimension, not a second row.
	if err := g.BookModel(context.Background(), "pa/m1", 10, 5, 700); err != nil {
		t.Fatal(err)
	}
	if len(store.rows) != 1 {
		t.Fatalf("store rows = %d, want 1 (same day/provider/model accumulates into one row)", len(store.rows))
	}
	if store.rows[0].Micros != 1000 || store.rows[0].Tasks != 2 ||
		store.rows[0].TokensIn != 110 || store.rows[0].TokensOut != 45 {
		t.Errorf("merged row = %+v, want 1000/2/110/45", store.rows[0])
	}
	if v := g.Decide("pa/m1"); v.Level != observe.LevelOver || v.AllowDispatch {
		t.Errorf("verdict after hitting the 1000-micro daily cap = %+v, want over/refused", v)
	}

	// Bad shapes never touch the disk.
	if err := g.BookModel(context.Background(), "nomodelkey", 1, 1, 1); err == nil {
		t.Error("key without a slash accepted")
	}
	if len(store.rows) != 1 {
		t.Errorf("refused call wrote rows: %d, want the single merged row", len(store.rows))
	}
}

// TestSeedBaselineSurvivesRestart is AC#3's fixture: gate A books, "restarts"
// (fresh gate, same store), gate B seeds, and both gates' readings agree —
// including the plan-credit bucket and the daily cap verdict.
func TestSeedBaselineSurvivesRestart(t *testing.T) {
	now := time.Date(2026, 11, 3, 9, 0, 0, 0, time.UTC)
	store := &t44bStore{rows: []observe.CostModelDayTotal{
		// Yesterday's rows are inside the current month: month folds include
		// them, the day bucket must not.
		{Day: "2026-11-01", Provider: "pa", Model: "m1", TokensIn: 10, TokensOut: 2, Micros: 200, Tasks: 1},
		{Day: "2026-11-03", Provider: "pa", Model: "m1", TokensIn: 20, TokensOut: 4, Micros: 300, Tasks: 1},
		{Day: "2026-11-03", Provider: "pl", Model: "m9", TokensIn: 1, TokensOut: 1, Micros: 50, Tasks: 1},
	}}

	live := newT44bGate(func() time.Time { return now }, store)
	seeded := newT44bGate(func() time.Time { return now }, store)
	n, err := seeded.SeedBaseline(context.Background())
	if err != nil {
		t.Fatalf("SeedBaseline: %v", err)
	}
	if n != 3 {
		t.Errorf("seeded %d rows, want 3 (the month's rows)", n)
	}
	// Drive both to the same state: the live gate books nothing yet, so mirror
	// the disk truth by feeding it the same rows through its own ledger.
	for _, r := range store.rows {
		if err := live.Ledger().SeedFrom([]observe.CostModelDayTotal{r}, now.Format("2006-01-02"), t44bCards()); err != nil {
			t.Fatal(err)
		}
	}

	a, b := live.Ledger().Spend(), seeded.Ledger().Spend()
	for _, key := range []string{"pa/m1", "pl/m9"} {
		if a.ModelDay[key] != b.ModelDay[key] || a.ModelMonth[key] != b.ModelMonth[key] {
			t.Errorf("restart drift on %s: live day/month %d/%d vs seeded %d/%d",
				key, a.ModelDay[key], a.ModelMonth[key], b.ModelDay[key], b.ModelMonth[key])
		}
	}
	if a.ProviderPlan["pl"] != 50 || b.ProviderPlan["pl"] != 50 {
		t.Errorf("plan credit not rebuilt equally: live %d seeded %d, want 50",
			a.ProviderPlan["pl"], b.ProviderPlan["pl"])
	}
	if a.ModelDay["pa/m1"] != 300 || a.ModelMonth["pa/m1"] != 500 {
		t.Errorf("day/month bucketing wrong after seed: day=%d month=%d, want 300/500",
			a.ModelDay["pa/m1"], a.ModelMonth["pa/m1"])
	}

	// A seeded gate refuses at its cap exactly like a live one.
	if v := seeded.Decide("pa/m1"); !v.AllowDispatch || v.Level != observe.LevelOK {
		t.Errorf("fresh-day seeded verdict = %+v, want OK/allow (500 of 1000 monthly, 300 daily)", v)
	}
	store.rows = append(store.rows, observe.CostModelDayTotal{Day: "2026-11-03", Provider: "pa", Model: "m1", Micros: 900})
	g2 := newT44bGate(func() time.Time { return now }, store)
	if _, err := g2.SeedBaseline(context.Background()); err != nil {
		t.Fatal(err)
	}
	if v := g2.Decide("pa/m1"); v.AllowDispatch || v.Level != observe.LevelOver {
		t.Errorf("after restart the gate that should be over quota dispatches anyway: %+v", v)
	}

	// Rollover after seeding is model-correct: a later day empties the day
	// bucket while the month keeps November's total until December.
	clock := time.Date(2026, 11, 4, 1, 0, 0, 0, time.UTC)
	g3 := newT44bGate(func() time.Time { return clock }, store)
	if _, err := g3.SeedBaseline(context.Background()); err != nil {
		t.Fatal(err)
	}
	sp := g3.Ledger().Spend()
	if sp.ModelDay["pa/m1"] != 0 {
		t.Errorf("after crossing midnight the seeded day bucket still holds %d, want 0", sp.ModelDay["pa/m1"])
	}
	if sp.ModelMonth["pa/m1"] == 0 {
		t.Error("after a day rollover the month bucket emptied; it must keep the month-to-date total")
	}
}

// TestBookModelWithoutLedgerNamesTheGap keeps the un-wired shape honest (the
// 44-r4 session-only behaviour must not silently pretend to persist).
func TestBookModelWithoutLedgerNamesTheGap(t *testing.T) {
	g := NewQuotaGate(observe.QuotaParams{}, nil, nil, nil, nil)
	if err := g.BookModel(context.Background(), "pa/m1", 1, 1, 1); err == nil {
		t.Error("BookModel succeeded without a model ledger")
	}
	if _, err := g.SeedBaseline(context.Background()); err == nil {
		t.Error("SeedBaseline succeeded without a model reader")
	}
}
