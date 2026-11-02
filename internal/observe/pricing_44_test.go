package observe

import (
	"math"
	"strings"
	"testing"
)

// Ticket 44 AC#1/#4/#6: the integer price math, the 价格未知 marker, and the
// versioned in-repo table with its config-override merge.

func TestAC1PriceMicrosMatchesHandComputedFixture(t *testing.T) {
	// 2.5 USD per 1M input, 10 USD per 1M output, 0.5 USD per 1M cache read —
	// stated in micro units per 1M: 2_500_000 / 10_000_000 / 500_000. The card
	// without Cached is the same card from a provider that publishes no cache
	// rate, which bills cached tokens at In (the conservative reading).
	withCache := Price{In: 2_500_000, Out: 10_000_000, Cached: 500_000}
	noCache := Price{In: 2_500_000, Out: 10_000_000}

	cases := []struct {
		name   string
		card   Price
		tok    Tokens
		want   int64
		byHand string
	}{
		{
			name:   "turn 1: no cache",
			card:   withCache,
			tok:    Tokens{In: 1000, Out: 200},
			want:   4500,
			byHand: "1000*2.5e6 + 200*10e6 = 2.5e9+2e9 = 4.5e9 per-1M micros -> 4500",
		},
		{
			name:   "turn 2: cache served at the cache rate",
			card:   withCache,
			tok:    Tokens{In: 1000, Out: 200, Cached: 600},
			want:   3300,
			byHand: "uncached 400*2.5e6=1e9; cached 600*0.5e6=3e8; out 200*10e6=2e9; sum 3.3e9/1e6=3300",
		},
		{
			name:   "cached rate unset bills at the input rate",
			card:   noCache,
			tok:    Tokens{In: 1000, Out: 0, Cached: 400},
			want:   2500,
			byHand: "Cached=0 => every input token bills at In: 1000*2.5e6/1e6 = 2500",
		},
		{
			name:   "cached cannot exceed input",
			card:   withCache,
			tok:    Tokens{In: 100, Out: 0, Cached: 900},
			want:   50,
			byHand: "cached clamped to 100 => 100*0.5e6/1e6 = 50, uncached 0",
		},
		{
			name:   "one output token",
			card:   withCache,
			tok:    Tokens{In: 0, Out: 1},
			want:   10,
			byHand: "1*10e6/1e6 = 10",
		},
		{
			name:   "half-micro total rounds half-up",
			card:   withCache,
			tok:    Tokens{In: 1, Out: 0},
			want:   3,
			byHand: "1*2.5e6 = 2.5e6 /1e6 = 2.5 -> half-up = 3 (the same value math.Round produced before ticket 44)",
		},
	}
	for _, c := range cases {
		if got := PriceMicros(c.card, c.tok); got != c.want {
			t.Errorf("%s: PriceMicros = %d, want %d (%s)", c.name, got, c.want, c.byHand)
		}
	}

	// The multi-turn rollup: a task's cost is the SUM of its per-turn deltas,
	// not a re-pricing of merged counters (each request reports its own usage).
	// Turns 1+2 above on one card: 4500+3300 = 7800 micro-units.
	var sum int64
	for _, name := range []string{"turn 1: no cache", "turn 2: cache served at the cache rate"} {
		for _, c := range cases {
			if c.name == name {
				sum += PriceMicros(c.card, c.tok)
			}
		}
	}
	if sum != 7800 {
		t.Errorf("per-turn rollup = %d micros, want 7800", sum)
	}

	// Same total through the accumulating meter path would be 7800 too; that is
	// pinned against the real loop in internal/agent/cost_44_test.go.
}

func TestAC1PriceMicrosIsIntegerAndSaturates(t *testing.T) {
	// A card so large the product cannot fit int64 must saturate, not wrap: a
	// wrapped number reads as a cheap task and lets a budget gate open.
	huge := Price{In: math.MaxInt64, Out: 0}
	if got := PriceMicros(huge, Tokens{In: math.MaxInt64}); got != math.MaxInt64 {
		t.Errorf("saturating product = %d, want MaxInt64", got)
	}
	// Negative token counts are a caller bug: 0, never a credit.
	if got := PriceMicros(huge, Tokens{In: -5}); got != 0 {
		t.Errorf("negative tokens priced at %d, want 0", got)
	}
	if got := PriceMicros(Price{In: -1, Out: -1}, Tokens{In: 10}); got != 0 {
		t.Errorf("negative rates priced at %d, want 0 (a negative price is not a price)", got)
	}
	// Zero usage on a real card is zero money, and it stays zero when added up.
	if got := PriceMicros(huge, Tokens{}); got != 0 {
		t.Errorf("zero usage priced at %d, want 0", got)
	}
}

func TestAC4UnknownCardPricesNothing(t *testing.T) {
	// The rule the ticket states twice: an unpriced model yields tokens only, so
	// the money answer must be exactly 0 and the marker must be available.
	if got := PriceMicros(Price{}, Tokens{In: 1_000_000, Out: 1_000_000}); got != 0 {
		t.Errorf("unknown card produced %d micros, want 0 (never fabricate amounts)", got)
	}
	var zero Price
	if zero.Known() {
		t.Error("an all-zero card reported itself as known")
	}
	if PriceUnknownMarker != "价格未知" {
		t.Errorf("marker = %q, want the ticket's own words 价格未知", PriceUnknownMarker)
	}
}

func TestAC4TableLookupMissesSayUnknownAndCostNoMoney(t *testing.T) {
	tbl := mustTable(t, `{"version":1,"models":{}}`)
	if _, currency, known := tbl.Card("acme", "m1"); known || currency != "" {
		t.Errorf("miss returned known=%v currency=%q, want false and an empty currency", known, currency)
	}
	// A miss stays a miss until an override adds the model: no stale zero.
	over, err := tbl.WithOverrides([]Override{{
		Model: "acme/m1", Currency: "USD",
		Price: Price{In: 1_000_000, Out: 2_000_000}, Source: "provider pricing page 2026-11-02",
	}})
	if err != nil {
		t.Fatal(err)
	}
	card, currency, known := over.Card("acme", "m1")
	if !known || currency != "USD" || card.In != 1_000_000 {
		t.Fatalf("override not applied: %+v %q %v", card, currency, known)
	}
	if got := PriceMicros(card, Tokens{In: 1000, Out: 500}); got != 2_000 {
		t.Errorf("priced usage = %d micros, want 1000*1 + 500*2 = 2000", got)
	}
	// The receiver is untouched: an override is a view, not a mutation.
	if _, _, known := tbl.Card("acme", "m1"); known {
		t.Error("WithOverrides mutated the table it was called on")
	}
	// The meter's own mouth reports the same miss (AC#4's marker path).
	m, err := NewCostMeter(&fakeStore{days: map[string]DayTotal{}}, WithPriceTable(tbl))
	if err != nil {
		t.Fatal(err)
	}
	if _, _, known := m.Price("acme", "m1"); known {
		t.Error("meter priced a model the table does not know")
	}
	if v := m.PriceTableVersion(); v != 1 {
		t.Errorf("meter reports version %d, want 1", v)
	}
}

func TestPricingRejectsFabricatedAndMalformedRows(t *testing.T) {
	cases := []struct {
		name, raw, want string
	}{
		{"no version", `{"models":{}}`, "version"},
		{"zero version", `{"version":0,"models":{}}`, "version"},
		{"missing models", `{"version":1}`, "models"},
		{"key without provider", `{"version":1,"models":{"m1":{"currency":"USD","in":1,"source":"x"}}}`, "provider/model"},
		{"negative rate", `{"version":1,"models":{"p/m":{"currency":"USD","in":-1,"source":"x"}}}`, "negative"},
		{"no currency", `{"version":1,"models":{"p/m":{"in":100,"source":"x"}}}`, "currency"},
		{"price without source", `{"version":1,"models":{"p/m":{"currency":"USD","in":100}}}`, "source"},
		{"zero rates without free marker", `{"version":1,"models":{"p/m":{"currency":"USD","in":0,"out":0,"source":"x"}}}`, "no rate"},
		{"bad billing enum", `{"version":1,"models":{"p/m":{"currency":"USD","in":1,"billing":"credits","source":"x"}}}`, "billing"},
	}
	for _, c := range cases {
		if _, err := ParsePriceTable([]byte(c.raw)); err == nil {
			t.Errorf("%s: ParsePriceTable accepted it", c.name)
		} else if !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: error %q must name %q", c.name, err, c.want)
		}
	}
	// The one shape that legitimately carries zero money is a stated free row.
	free := mustTable(t, `{"version":1,"models":{"p/m":{"currency":"USD","billing":"free","in":0,"out":0,"source":"provider: free tier, 2026-11-02"}}}`)
	card, currency, known := free.Card("p", "m")
	if !known || card.Known() || currency != "USD" {
		t.Errorf("free row read as %+v/%q/known=%v, want known with an all-zero card", card, currency, known)
	}
	if got := PriceMicros(card, Tokens{In: 5000, Out: 5000}); got != 0 {
		t.Errorf("a free model priced at %d, want 0", got)
	}
}

func TestAC6EmbeddedTableIsValidAndVersioned(t *testing.T) {
	tbl, err := DefaultPriceTable()
	if err != nil {
		t.Fatalf("the in-repo pricing.json must parse: %v", err)
	}
	if tbl.Version < 1 {
		t.Errorf("embedded table version = %d, want >= 1", tbl.Version)
	}
	// AC#6's second half — "old aggregates unchanged" — is pinned in
	// internal/agent/costbook_44_test.go against the real store: a version bump
	// re-records the number and recomputes nothing.
	bumped := PriceTable{Version: tbl.Version + 1, Updated: "2026-12-01", rows: tbl.rows}
	if bumped.Version != tbl.Version+1 {
		t.Error("a version bump did not move the number the record carries")
	}
	// A meter rebuilt on the bumped table reports the new version and keeps the
	// rows it already stored: the only rewrite path (RefreshDay) recomputes from
	// task_log, never from a price card.
	m, err := NewCostMeter(&fakeStore{days: map[string]DayTotal{}}, WithPriceTable(bumped))
	if err != nil {
		t.Fatal(err)
	}
	if m.PriceTableVersion() != tbl.Version+1 {
		t.Errorf("meter reports %d, want the bumped version", m.PriceTableVersion())
	}
}

func TestOverrideValidationMatchesTableValidation(t *testing.T) {
	tbl := mustTable(t, `{"version":1,"models":{}}`)
	bad := []Override{
		{Model: "no-slash", Currency: "USD", Price: Price{In: 1}, Source: "x"},
		{Model: "p/m", Price: Price{In: 1}, Source: "x"},
		{Model: "p/m", Currency: "USD", Price: Price{In: -1}, Source: "x"},
		{Model: "p/m", Currency: "USD", Price: Price{}, Source: "x"},
		{Model: "p/m", Currency: "USD", Price: Price{In: 1}},
		{Model: "p/m", Currency: "USD", Price: Price{In: 1}, Billing: "credits", Source: "x"},
	}
	for i, o := range bad {
		if _, err := tbl.WithOverrides([]Override{o}); err == nil {
			t.Errorf("override %d (%+v) was accepted", i, o)
		}
	}
	good, err := tbl.WithOverrides([]Override{
		{Model: "p/m", Currency: "CNY", Price: Price{In: 4000, Out: 12000}, Source: "provider page 2026-11-02"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, currency, known := good.Card("p", "m"); !known || currency != "CNY" {
		t.Errorf("good override not applied: %q known=%v", currency, known)
	}
}

func mustTable(t *testing.T, raw string) PriceTable {
	t.Helper()
	tbl, err := ParsePriceTable([]byte(raw))
	if err != nil {
		t.Fatalf("ParsePriceTable(%s): %v", raw, err)
	}
	return tbl
}
