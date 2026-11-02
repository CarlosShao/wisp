package observe

// C23 price cards and the price table (ticket 44).
//
// A price is declared in INTEGER micro-currency units per 1,000,000 tokens —
// the unit task_log.cost_amount_micro and cost_daily.cost_micros already store
// (D35 / SPEC-02 §3). Money math therefore never touches a float: one card and
// a token count produce int64 micros (PriceMicros below), and every rollup
// (memory's cost_daily upsert, this package's day/month sums) adds int64.
//
// The table below is the starter set the ticket asks for ("in-repo JSON,
// versioned; editable via config override"). It deliberately ships NO invented
// numbers: 44-costmeter-c23.md:19-20 says an unknown price must read
// tokens-only with the 价格未知 marker and "never fabricate amounts", so a row
// is only added when its source is the provider's own published card. Rows are
// keyed "provider/model" — the same spelling config's llm.text_chain elements
// and config.checkChainElement already use — and a host adds prices through
// PriceTable.WithOverrides (config override), never by editing this file.
//
// Currency: a card states its own currency and the meter carries that label
// verbatim. config.Price is documented as micro-USD per 1M tokens while
// task_log.currency defaults to CNY (agent/cost.go's unit note); this package
// never converts between them, because an exchange rate would be an invented
// number.

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"math"
	"math/bits"
	"sort"
	"strings"
)

// PriceUnknownMarker is the user-visible marker for a model with no price card.
// It is 44's own wording (票面:19 逐字 "价格未知").
const PriceUnknownMarker = "价格未知"

// million is the denominator every card is declared against.
const million = 1_000_000

// Price is one model's three published rates, in micro-currency units per 1M
// tokens. Cached == 0 means the provider declares no separate cache rate, and
// cached tokens are then billed at In (the conservative reading agent/cost.go
// already documents).
type Price struct {
	In     int64
	Out    int64
	Cached int64
}

// Known reports whether the card carries at least one rate. An all-zero card
// is indistinguishable from "no card at all", so the meter treats it as
// unknown rather than as a free model: reporting 0 micros for a model nobody
// priced would read as "this cost nothing", which is the fabrication the
// ticket forbids. A genuinely free model must say so with a card marked
// Billing "free" (see entrySource).
func (p Price) Known() bool { return p.In != 0 || p.Out != 0 || p.Cached != 0 }

// Tokens is the C6 usage triple this package accumulates (InputTokens /
// OutputTokens / CachedTokens as llm.Usage reports them; observe cannot import
// llm, so the host copies the three counters).
type Tokens struct {
	In     int64
	Out    int64
	Cached int64
}

// PriceMicros prices one usage block against one card, in integer micro
// units. Three properties are pinned by tests (ticket 44 AC#1/#4):
//
//   - integer math only: the product tokens*rate is built with bits.Mul64 and
//     divided once, with a half-up remainder rule identical to what
//     math.Round produced for the pre-44 float path;
//   - unknown card (or negative token counts) => 0, never a guess;
//   - cached tokens never exceed the input tokens they are carved out of.
//
// A product that cannot fit int64 saturates at MaxInt64 instead of wrapping:
// saturating is visible (the aggregate reads absurdly large and the budget
// gate stops), wrapping would silently read as a tiny cost.
func PriceMicros(p Price, t Tokens) int64 {
	if !p.Known() || t.In < 0 || t.Out < 0 || t.Cached < 0 {
		return 0
	}
	cached := t.Cached
	if cached > t.In {
		cached = t.In
	}
	cachedRate := p.Cached
	if cachedRate == 0 {
		cachedRate = p.In
	}
	uncached := t.In - cached

	var num uint64
	for _, term := range []struct {
		tokens, rate int64
	}{
		{uncached, p.In},
		{cached, cachedRate},
		{t.Out, p.Out},
	} {
		if term.tokens == 0 || term.rate <= 0 {
			continue
		}
		hi, lo := bits.Mul64(uint64(term.tokens), uint64(term.rate))
		if hi != 0 {
			return math.MaxInt64
		}
		var carry uint64
		num, carry = bits.Add64(num, lo, 0)
		if carry != 0 {
			return math.MaxInt64
		}
	}
	if num > math.MaxInt64-million {
		return math.MaxInt64
	}
	q := (num + million/2) / million
	if q > math.MaxInt64 {
		return math.MaxInt64
	}
	return int64(q)
}

// entry is one row of the embedded price table.
type entry struct {
	// Currency labels the three rates (e.g. "USD", "CNY"). Required: a rate
	// without a currency is not a price anyone can act on.
	Currency string `json:"currency"`
	// Billing is "pay-per-token" or "free". Only "free" lets a row carry three
	// zero rates, and it says so out loud instead of leaving the meter to
	// guess whether 0 means free or unpriced.
	Billing string `json:"billing"`
	// In/Out/Cached: micro-currency units per 1M tokens.
	In     int64 `json:"in"`
	Out    int64 `json:"out"`
	Cached int64 `json:"cached"`
	// Source names where the card came from (provider pricing page + access
	// date). A row without a source is a fabricated number and fails Load.
	Source string `json:"source"`
}

// priceFile is the JSON shape of pricing.json.
type priceFile struct {
	Version int              `json:"version"`
	Updated string           `json:"updated"`
	Notes   string           `json:"notes"`
	Models  map[string]entry `json:"models"`
}

// PriceTable is a parsed, validated price table.
type PriceTable struct {
	Version int
	Updated string
	// rows is the validated set keyed "provider/model".
	rows map[string]entry
}

// pricingJSON is the in-repo starter table (ticket 44: "in-repo JSON,
// versioned"). Editing it means a version bump; user edits belong in the config
// override instead, so a preset refresh never clobbers them.
//
//go:embed pricing.json
var pricingJSON []byte

// DefaultPriceTable parses and validates the embedded starter table. It errors
// rather than degrading to an empty table: a corrupt price file must stop the
// caller that asked for it, and "no prices" is a claim about the world.
func DefaultPriceTable() (PriceTable, error) { return ParsePriceTable(pricingJSON) }

// ParsePriceTable validates raw JSON (exported so a test or an override file
// can go through the same rules as the embedded table).
func ParsePriceTable(raw []byte) (PriceTable, error) {
	var f priceFile
	if err := json.Unmarshal(raw, &f); err != nil {
		return PriceTable{}, fmt.Errorf("observe: cost price table: %w", err)
	}
	if f.Version < 1 {
		return PriceTable{}, fmt.Errorf("observe: cost price table: version %d must be >= 1", f.Version)
	}
	if f.Models == nil {
		return PriceTable{}, fmt.Errorf("observe: cost price table: missing models object")
	}
	rows := make(map[string]entry, len(f.Models))
	for key, e := range f.Models {
		if err := validModelKey(key); err != nil {
			return PriceTable{}, fmt.Errorf("observe: cost price table row %q: %w", key, err)
		}
		if e.Currency == "" {
			return PriceTable{}, fmt.Errorf("observe: cost price table row %q: currency is required", key)
		}
		if e.In < 0 || e.Out < 0 || e.Cached < 0 {
			return PriceTable{}, fmt.Errorf("observe: cost price table row %q: a negative price is not a price", key)
		}
		switch e.Billing {
		case "free":
			// A free model is stated, not implied: zero rates with no marker
			// would be indistinguishable from a missing card.
			if e.In == 0 && e.Out == 0 {
				if e.Source == "" {
					return PriceTable{}, fmt.Errorf("observe: cost price table row %q: billing=free needs a source", key)
				}
			}
		case "pay-per-token", "":
			card := Price{In: e.In, Out: e.Out, Cached: e.Cached}
			if !card.Known() {
				return PriceTable{}, fmt.Errorf("observe: cost price table row %q: billing=%q but no rate to price with", key, orDefault(e.Billing, "pay-per-token"))
			}
			if e.Source == "" {
				return PriceTable{}, fmt.Errorf("observe: cost price table row %q: a price needs the source that published it", key)
			}
		default:
			return PriceTable{}, fmt.Errorf("observe: cost price table row %q: billing %q must be pay-per-token or free", key, e.Billing)
		}
		rows[key] = e
	}
	return PriceTable{Version: f.Version, Updated: f.Updated, rows: rows}, nil
}

// Card returns one model's price card. known=false means the table has no
// card for it, which the meter renders as tokens-only + PriceUnknownMarker.
func (t PriceTable) Card(provider, model string) (Price, string, bool) {
	e, ok := t.rows[provider+"/"+model]
	if !ok {
		return Price{}, "", false
	}
	if e.Billing == "free" {
		// A stated-free row is priced: the answer is genuinely 0, not "we do
		// not know". Currency still rides along for the label.
		return Price{}, e.Currency, true
	}
	return Price{In: e.In, Out: e.Out, Cached: e.Cached}, e.Currency, true
}

// Models lists the keys the table carries, sorted (a stable surface for logs
// and for AC#6's version record).
func (t PriceTable) Models() []string {
	out := make([]string, 0, len(t.rows))
	for k := range t.rows {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Override is one model's price supplied by the host (config override, ticket
// 44: "editable via config override"). Source names where it came from so a
// user-added card is distinguishable from the in-repo one in a report.
type Override struct {
	Model    string // "provider/model"
	Currency string
	Price    Price
	Billing  string // "" = pay-per-token
	Source   string
}

// WithOverrides returns a table with the given rows applied on top. The
// receiver is never mutated: an override is a view, so a hot-reloaded config
// cannot rewrite the table another meter is already pricing against.
// Overrides are validated by the same rules as the embedded rows (currency and
// source required, no negative rates), and the version is reported as the
// embedded one plus nothing: callers who change rates change pricing.json's
// version, and AC#6's record is that version number.
func (t PriceTable) WithOverrides(overs []Override) (PriceTable, error) {
	rows := make(map[string]entry, len(t.rows)+len(overs))
	for k, v := range t.rows {
		rows[k] = v
	}
	for _, o := range overs {
		if err := validModelKey(o.Model); err != nil {
			return PriceTable{}, fmt.Errorf("observe: cost price override %q: %w", o.Model, err)
		}
		if o.Currency == "" {
			return PriceTable{}, fmt.Errorf("observe: cost price override %q: currency is required", o.Model)
		}
		if o.Price.In < 0 || o.Price.Out < 0 || o.Price.Cached < 0 {
			return PriceTable{}, fmt.Errorf("observe: cost price override %q: a negative price is not a price", o.Model)
		}
		billing := o.Billing
		if billing == "" {
			billing = "pay-per-token"
		}
		if billing != "free" && billing != "pay-per-token" {
			return PriceTable{}, fmt.Errorf("observe: cost price override %q: billing %q must be pay-per-token or free", o.Model, billing)
		}
		if billing == "pay-per-token" && !o.Price.Known() {
			return PriceTable{}, fmt.Errorf("observe: cost price override %q: billing %q but no rate to price with", o.Model, billing)
		}
		if o.Source == "" {
			return PriceTable{}, fmt.Errorf("observe: cost price override %q: a price needs the source that published it", o.Model)
		}
		rows[o.Model] = entry{
			Currency: o.Currency, Billing: billing,
			In: o.Price.In, Out: o.Price.Out, Cached: o.Price.Cached,
			Source: o.Source,
		}
	}
	return PriceTable{Version: t.Version, Updated: t.Updated, rows: rows}, nil
}

// validModelKey checks the "provider/model" spelling the catalog already uses.
func validModelKey(key string) error {
	if strings.TrimSpace(key) != key || key == "" {
		return fmt.Errorf("model key must not be empty or padded")
	}
	if strings.Count(key, "/") != 1 {
		return fmt.Errorf("model key %q must be \"provider/model\"", key)
	}
	p, m, _ := strings.Cut(key, "/")
	if p == "" || m == "" {
		return fmt.Errorf("model key %q must name both a provider and a model", key)
	}
	if strings.ContainsAny(key, " \t\n") {
		return fmt.Errorf("model key %q must not contain whitespace", key)
	}
	return nil
}

func orDefault(v, def string) string {
	if v == "" {
		return def
	}
	return v
}
