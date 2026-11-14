package observe

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Concurrency instrument for the QuotaLedger (ticket 44 family; the red this
// reproduces is ticket 70's CI step 10, `Test197Subagent…`: a recovered
// goroutine panic inside QuotaLedger.roll). Evidence:
// docs/evidence/s1/44-quota-ledger-race.md.
//
// The product shape being modelled is 44-r4's QuotaGate: Decide calls
// ledger.Spend() on every dispatch (quota_44r4.go), Book calls ledger.Post
// after the turn, and several tasks (a root agent plus its subagents) run at
// the same time — so Spend and Post are concurrently live, and Spend/Post each
// call roll(), which writes the bucket stamps and clears the maps whenever the
// calendar day or month moved. Both tests fail on a ledger with no
// synchronization and pass on one that has it; neither asserts a latency or a
// schedule, so they cannot flake into a pass.

// quotaRaceCard is the card both tests spend against: caps far above the
// amounts booked so that no dispatch decision is involved — this instrument
// is about memory safety, not about the 80/100 boundaries.
var quotaRaceCard = QuotaModel{
	Provider:      "raceprov",
	Model:         "racemodel",
	Billing:       BillingPlan,
	DailyMicros:   1 << 40,
	MonthlyMicros: 1 << 44,
}

// TestQuotaLedgerConcurrentSpendAcrossBucketRolls drives N goroutines through
// Spend() while the injected clock advances one day per read, so every roll()
// has a bucket to change: the exact writer/reader pair ticket 70 step 10 read
// out of the CI log (a roll clearing the maps while another goroutine is
// ranging over them).
func TestQuotaLedgerConcurrentSpendAcrossBucketRolls(t *testing.T) {
	const goroutines = 16
	const rounds = 200

	base := time.Date(2026, 3, 1, 9, 0, 0, 0, time.Local)
	var tick atomic.Int64
	clock := func() time.Time {
		// -1 so the first read lands exactly on base.
		return base.AddDate(0, 0, int(tick.Add(1)-1))
	}
	l := NewQuotaLedger(clock)

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < rounds; j++ {
				sp := l.Spend()
				if sp.ModelDay == nil || sp.ModelMonth == nil || sp.ProviderPlan == nil {
					// Report without t.Fatalf from a non-test goroutine: a
					// nil map would mean the snapshot lost its buckets.
					panic("QuotaLedger.Spend returned a nil bucket map")
				}
			}
		}()
	}
	wg.Wait()
}

// TestQuotaLedgerConcurrentSpendAndPostKeepEveryMicro runs readers and writers
// at the same time on a stable calendar day (so roll() takes its quiet branch
// and the maps are the only shared state): 8 goroutines snapshotting with
// Spend(), 8 booking with Post(). Every goroutine's booking must be visible in
// the end state — a lost update is a safety bug even when the race detector is
// quiet.
func TestQuotaLedgerConcurrentSpendAndPostKeepEveryMicro(t *testing.T) {
	const readers = 8
	const writers = 8
	const rounds = 200
	const micros = 5

	now := time.Date(2026, 3, 2, 10, 0, 0, 0, time.Local)
	l := NewQuotaLedger(func() time.Time { return now })
	key := quotaRaceCard.Key()

	var wg sync.WaitGroup
	wg.Add(readers + writers)
	for i := 0; i < readers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < rounds; j++ {
				_ = l.Spend()
			}
		}()
	}
	for i := 0; i < writers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < rounds; j++ {
				if err := l.Post(quotaRaceCard, micros); err != nil {
					panic("QuotaLedger.Post refused a legal positive spend: " + err.Error())
				}
			}
		}()
	}
	wg.Wait()

	want := int64(writers * rounds * micros)
	sp := l.Spend()
	if got := sp.ModelDay[key]; got != want {
		t.Fatalf("ModelDay[%s] = %d, want %d (a booked micro was lost)", key, got, want)
	}
	if got := sp.ModelMonth[key]; got != want {
		t.Fatalf("ModelMonth[%s] = %d, want %d (a booked micro was lost)", key, got, want)
	}
	if got := sp.ProviderPlan[quotaRaceCard.Provider]; got != want {
		t.Fatalf("ProviderPlan[%s] = %d, want %d (plan credit lost a booked micro)",
			quotaRaceCard.Provider, got, want)
	}
}

// TestQuotaLedgerSeedFromOverlapsSpend is the positive control for the one
// writer the mutex skipped: SeedFrom rebuilds the three buckets from stored
// rows (QuotaGate.SeedBaseline on boot), and on the unfixed ledger it did that
// with no lock while Spend was already making its copy under one. The product
// shape is a boot seed that loses the race with the first dispatch's Decide:
// one goroutine ranges a bucket map while the seeder is writing it, and the
// pair has no happens-before edge to hide behind.
//
// The assertions after the barrier are the safety half, and they are what keeps
// this test honest once the detector is quiet: every seeded micro must land in
// the snapshot exactly once, so a fix that serialises the reads but drops a
// seeded row still fails.
func TestQuotaLedgerSeedFromOverlapsSpend(t *testing.T) {
	const readers = 4
	const seedRounds = 40
	const readerRounds = 200
	const dayMicros = 7
	const monthMicros = 11

	currentDay := "2026-11-03"
	now := time.Date(2026, 11, 3, 9, 0, 0, 0, time.UTC)

	rows := []CostModelDayTotal{
		// Today's plan-billed row: day bucket, month bucket, provider credit.
		{Day: currentDay, Provider: quotaRaceCard.Provider, Model: quotaRaceCard.Model, Micros: dayMicros},
		// A row from yesterday of a model with no card: month bucket only.
		{Day: "2026-11-02", Provider: "seedprov", Model: "seedmodel", Micros: monthMicros},
	}
	cards := map[string]QuotaModel{quotaRaceCard.Key(): quotaRaceCard}

	l := NewQuotaLedger(func() time.Time { return now })

	var wg sync.WaitGroup
	wg.Add(readers + 1)
	for i := 0; i < readers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < readerRounds; j++ {
				sp := l.Spend()
				if sp.ModelDay == nil || sp.ModelMonth == nil || sp.ProviderPlan == nil {
					// Report without t.Fatalf from a non-test goroutine.
					t.Errorf("QuotaLedger.Spend returned a nil bucket map while a seed was in flight")
					return
				}
			}
		}()
	}
	go func() {
		defer wg.Done()
		for j := 0; j < seedRounds; j++ {
			if err := l.SeedFrom(rows, currentDay, cards); err != nil {
				t.Errorf("QuotaLedger.SeedFrom refused legal seed rows: %v", err)
				return
			}
		}
	}()
	wg.Wait()

	wantDay := int64(seedRounds * dayMicros)
	wantMonth := int64(seedRounds * (dayMicros + monthMicros))
	sp := l.Spend()
	if got := sp.ModelDay[quotaRaceCard.Key()]; got != wantDay {
		t.Fatalf("ModelDay[%s] = %d, want %d (a seeded micro was lost)", quotaRaceCard.Key(), got, wantDay)
	}
	if got := sp.ModelMonth[quotaRaceCard.Key()]; got != wantDay {
		t.Fatalf("ModelMonth[%s] = %d, want %d (a seeded micro was lost)", quotaRaceCard.Key(), got, wantDay)
	}
	if got := sp.ProviderPlan[quotaRaceCard.Provider]; got != wantDay {
		t.Fatalf("ProviderPlan[%s] = %d, want %d (plan credit lost a seeded micro)",
			quotaRaceCard.Provider, got, wantDay)
	}
	if got := sp.ModelMonth["seedprov/seedmodel"]; got != wantMonth-wantDay {
		t.Fatalf("ModelMonth[seedprov/seedmodel] = %d, want %d (a month-only row was lost)",
			got, wantMonth-wantDay)
	}
	if got := sp.ModelDay["seedprov/seedmodel"]; got != 0 {
		t.Fatalf("ModelDay[seedprov/seedmodel] = %d, want 0 (a non-today row leaked into the day bucket)", got)
	}
}
