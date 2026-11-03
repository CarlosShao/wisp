package observe

import (
	"context"
	"testing"
	"time"
)

// Ticket 44b observe-layer tests: the same-batch model booking, the seed into
// the quota ledger, and the model-dimension rollover.

type t44bLedger struct {
	bumped []t44bRow
	err    error
}

type t44bRow struct {
	day, provider, model   string
	in, out, micros, tasks int64
}

func (l *t44bLedger) BumpCostDayModel(_ context.Context, day, provider, model string, tokensIn, tokensOut, costMicros int64) error {
	if l.err != nil {
		return l.err
	}
	l.bumped = append(l.bumped, t44bRow{day, provider, model, tokensIn, tokensOut, costMicros, 1})
	return nil
}

func newT44bMeter(t *testing.T, now time.Time, ledger ModelLedger) *CostMeter {
	t.Helper()
	m, err := NewCostMeter(&fakeStore{days: map[string]DayTotal{}},
		WithClock(func() time.Time { return now }),
		WithLocation(time.UTC),
		WithModelLedger(ledger))
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func TestPostModelTaskBooksOneBatchUnderTheMeterDay(t *testing.T) {
	l := &t44bLedger{}
	now := time.Date(2026, 11, 3, 9, 0, 0, 0, time.UTC)
	m := newT44bMeter(t, now, l)
	day, err := m.PostModelTask(context.Background(), "pa", "m1", 100, 40, 3000)
	if err != nil {
		t.Fatal(err)
	}
	if day != "2026-11-03" {
		t.Errorf("day label = %q, want the meter calendar's day", day)
	}
	if len(l.bumped) != 1 {
		t.Fatalf("bumped %d rows, want exactly 1 (one call = one batch)", len(l.bumped))
	}
	r := l.bumped[0]
	if r.day != "2026-11-03" || r.provider != "pa" || r.model != "m1" ||
		r.in != 100 || r.out != 40 || r.micros != 3000 {
		t.Errorf("bumped row %+v does not match the task", r)
	}
}

func TestPostModelTaskRefusesBadShapes(t *testing.T) {
	l := &t44bLedger{}
	m := newT44bMeter(t, time.Date(2026, 11, 3, 9, 0, 0, 0, time.UTC), l)
	ctx := context.Background()
	if _, err := m.PostModelTask(ctx, "", "m1", 1, 1, 1); err == nil {
		t.Error("empty provider accepted")
	}
	if _, err := m.PostModelTask(ctx, "pa", "", 1, 1, 1); err == nil {
		t.Error("empty model accepted")
	}
	if _, err := m.PostModelTask(ctx, "pa", "m1", -1, 0, 0); err == nil {
		t.Error("negative tokens accepted")
	}
	if _, err := m.PostModelTask(ctx, "pa", "m1", 0, 0, -5); err == nil {
		t.Error("negative micros accepted")
	}
	if len(l.bumped) != 0 {
		t.Errorf("refused books reached the ledger: %v", l.bumped)
	}
	// Unknown price: tokens-only with micros 0 is a legal booking (44's marker
	// lives at read time, not as a fabricated number at book time).
	if _, err := m.PostModelTask(ctx, "px", "unknown", 55, 11, 0); err != nil {
		t.Errorf("tokens-only booking refused: %v", err)
	}
	// No ledger wired: the method says so instead of silently booking nothing.
	bare, err := NewCostMeter(&fakeStore{days: map[string]DayTotal{}}, WithClock(func() time.Time { return time.Unix(0, 0) }))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := bare.PostModelTask(ctx, "pa", "m1", 1, 1, 1); err == nil {
		t.Error("PostModelTask without WithModelLedger succeeded")
	}
}

func TestMonthFirst(t *testing.T) {
	if got := MonthFirst("2026-11-03"); got != "2026-11-01" {
		t.Errorf("MonthFirst = %q, want 2026-11-01", got)
	}
}

func TestQuotaLedgerSeedFromBucketsDayMonthPlan(t *testing.T) {
	now := time.Date(2026, 11, 3, 9, 0, 0, 0, time.UTC)
	l := NewQuotaLedger(func() time.Time { return now })
	cards := map[string]QuotaModel{
		"pa/m1": {Provider: "pa", Model: "m1", Billing: BillingPayPerToken, DailyMicros: 10000, MonthlyMicros: 100000},
		"pl/m9": {Provider: "pl", Model: "m9", Billing: BillingPlan},
	}
	rows := []CostModelDayTotal{
		{Day: "2026-11-01", Provider: "pa", Model: "m1", Micros: 500},
		{Day: "2026-11-03", Provider: "pa", Model: "m1", Micros: 300},
		{Day: "2026-11-02", Provider: "pl", Model: "m9", Micros: 70},
	}
	if err := l.SeedFrom(rows, "2026-11-03", cards); err != nil {
		t.Fatal(err)
	}
	sp := l.Spend()
	if sp.ModelDay["pa/m1"] != 300 {
		t.Errorf("seeded ModelDay[pa/m1] = %d, want 300 (only today's rows)", sp.ModelDay["pa/m1"])
	}
	if sp.ModelMonth["pa/m1"] != 800 {
		t.Errorf("seeded ModelMonth[pa/m1] = %d, want 800 (whole month folds)", sp.ModelMonth["pa/m1"])
	}
	if sp.ProviderPlan["pl"] != 70 {
		t.Errorf("seeded ProviderPlan[pl] = %d, want 70 (plan billing draws credit)", sp.ProviderPlan["pl"])
	}
	if sp.ProviderPlan["pa"] != 0 {
		t.Errorf("pay-per-token spend darkened the plan bucket: %d", sp.ProviderPlan["pa"])
	}
	// Negative stored rows are refused, not credited.
	if err := l.SeedFrom([]CostModelDayTotal{{Day: "2026-11-03", Provider: "pa", Model: "m1", Micros: -1}}, "2026-11-03", cards); err == nil {
		t.Error("negative seeded spend accepted")
	}
}

// TestQuotaLedgerSeedRolloverIsModelCorrect is AC#4: after seeding, crossing
// midnight empties the day buckets (per model), and crossing the month empties
// the month buckets too — the seeded baseline rolls exactly like live spend.
func TestQuotaLedgerSeedRolloverIsModelCorrect(t *testing.T) {
	clock := time.Date(2026, 11, 30, 23, 0, 0, 0, time.UTC)
	l := NewQuotaLedger(func() time.Time { return clock })
	cards := map[string]QuotaModel{"pa/m1": {Provider: "pa", Model: "m1"}}
	rows := []CostModelDayTotal{
		{Day: "2026-11-30", Provider: "pa", Model: "m1", Micros: 400},
	}
	if err := l.SeedFrom(rows, DayKey(clock), cards); err != nil {
		t.Fatal(err)
	}
	if sp := l.Spend(); sp.ModelDay["pa/m1"] != 400 || sp.ModelMonth["pa/m1"] != 400 {
		t.Fatalf("pre-rollover spend = %+v", sp)
	}
	// Midnight: the day bucket empties, the month still holds November's total.
	clock = clock.Add(2 * time.Hour) // 2026-12-01 01:00 — month rolls too.
	sp := l.Spend()
	if sp.ModelDay["pa/m1"] != 0 || sp.ModelMonth["pa/m1"] != 0 {
		t.Errorf("after crossing into December day/month still hold %d/%d, want 0/0 (month rolled)",
			sp.ModelDay["pa/m1"], sp.ModelMonth["pa/m1"])
	}
	// And a same-month day crossing keeps the month number: seed with the
	// ledger's real today (2026-11-03); a row from the 30th folds into the
	// month only, never into today's bucket.
	l2 := NewQuotaLedger(func() time.Time { return time.Date(2026, 11, 3, 9, 0, 0, 0, time.UTC) })
	if err := l2.SeedFrom(rows, "2026-11-03", cards); err != nil {
		t.Fatal(err)
	}
	// row day 2026-11-30 != today 2026-11-03: it folds into the month only.
	sp2 := l2.Spend()
	if sp2.ModelDay["pa/m1"] != 0 || sp2.ModelMonth["pa/m1"] != 400 {
		t.Errorf("month-only row bucketed wrong: day=%d month=%d, want 0/400",
			sp2.ModelDay["pa/m1"], sp2.ModelMonth["pa/m1"])
	}
}
