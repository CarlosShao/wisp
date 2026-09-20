package agent

import (
	"math"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/llm"
)

// C23 per-task cost data. Only the numbers a task owns live here: token
// counters, tool-call count, elapsed time and the money estimate from the
// model's price card. The daily/monthly aggregation (cost_daily) and the
// budget-pause policy belong to the CostMeter ticket (44), which is the only
// consumer of this besides the task_log row.
//
// Unit note (open question registered in the ticket report): config.Price is
// documented as "micro-USD per 1M tokens" while task_log.currency defaults to
// CNY. The loop therefore computes micro-units of whatever currency the price
// card is declared in and labels the visible message with Config.Currency
// verbatim; it never invents an exchange rate.

// Cost tracks one task's C23 counters.
type Cost struct {
	Usage     llm.Usage
	ToolCalls int
	// Micros is the accumulated 1e-6-currency-unit estimate.
	Micros int64
	// Currency labels Micros (config-provided; "" = micro units only).
	Currency string
}

// AddUsage prices one Usage block against a price card and accumulates it.
// Counters SUM across turns (each request reports its own usage), unlike
// llm.Usage.Add's max-merge which collapses duplicate reports inside one
// stream. Cached prompt tokens are billed at Price.Cached when the card sets
// it (0 = treat cached like the input rate, the conservative reading of a card
// that declares no cache price).
func (c *Cost) AddUsage(p config.Price, u llm.Usage) int64 {
	delta := PriceMicros(p, u)
	c.Usage = llm.Usage{
		InputTokens:  c.Usage.InputTokens + u.InputTokens,
		OutputTokens: c.Usage.OutputTokens + u.OutputTokens,
		CachedTokens: c.Usage.CachedTokens + u.CachedTokens,
	}
	c.Micros += delta
	return delta
}

// PriceMicros computes the 1e-6-unit cost of one usage block.
func PriceMicros(p config.Price, u llm.Usage) int64 {
	in := int64(u.InputTokens)
	cached := int64(u.CachedTokens)
	if cached > in {
		cached = in
	}
	cachedRate := p.Cached
	if cachedRate == 0 {
		cachedRate = p.In
	}
	uncached := float64(in-cached)*p.In + float64(cached)*cachedRate
	total := uncached + float64(u.OutputTokens)*p.Out
	return int64(math.Round(total / 1_000_000))
}
