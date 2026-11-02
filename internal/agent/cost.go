package agent

import (
	"math"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
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
//
// Ticket 44 moved the arithmetic itself into observe (integer micro-units, no
// float): this file keeps the per-task accumulator and the config-card adapter,
// and PriceMicros below is a one-line delegation, so there is exactly ONE
// money calculator in the tree. A second one would drift the moment a rounding
// rule changes, and task_log, cost_daily and the visible line are three mouths
// that must say the same number.

// Cost tracks one task's C23 counters.
type Cost struct {
	Usage llm.Usage
	// ToolCalls counts this task's tool calls (C23's "工具次数" column).
	ToolCalls int
	// Micros is the accumulated 1e-6-currency-unit estimate.
	Micros int64
	// Currency labels Micros (config-provided; "" = micro units only).
	Currency string
	// Priced records whether the card that produced Micros was known. false
	// means the token counters are real and the money figure is deliberately
	// absent: the visible line then shows tokens plus observe's 价格未知 marker
	// instead of a number nobody priced (ticket 44 AC#4).
	Priced bool
}

// AddUsage prices one Usage block against a price card and accumulates it.
// Counters SUM across turns (each request reports its own usage), unlike
// llm.Usage.Add's max-merge which collapses duplicate reports inside one
// stream. Cached prompt tokens are billed at Price.Cached when the card sets
// it (0 = treat cached like the input rate, the conservative reading of a card
// that declares no cache price).
func (c *Cost) AddUsage(p config.Price, u llm.Usage) int64 {
	card := CardOf(p)
	delta := observe.PriceMicros(card, TokensOf(u))
	c.Usage = llm.Usage{
		InputTokens:  c.Usage.InputTokens + u.InputTokens,
		OutputTokens: c.Usage.OutputTokens + u.OutputTokens,
		CachedTokens: c.Usage.CachedTokens + u.CachedTokens,
	}
	c.Micros += delta
	if card.Known() {
		c.Priced = true
	}
	return delta
}

// CardOf adapts config's price card to the meter's integer card. The config
// field is a float64 because D36 froze it that way (internal/config/schema.go's
// Price); converting the declared RATE is the only place a float is read, and
// it is rounded to the nearest whole micro-per-1M rather than truncated — a
// half-micro rate would otherwise silently vanish, and a rate that vanishes is
// a task that books zero. Money AMOUNTS are int64 from here on.
func CardOf(p config.Price) observe.Price {
	return observe.Price{
		In:     wholeMicros(p.In),
		Out:    wholeMicros(p.Out),
		Cached: wholeMicros(p.Cached),
	}
}

// TokensOf copies C6's usage triple into the meter's counter shape (observe
// cannot import llm, so the copy happens here).
func TokensOf(u llm.Usage) observe.Tokens {
	return observe.Tokens{
		In:     int64(u.InputTokens),
		Out:    int64(u.OutputTokens),
		Cached: int64(u.CachedTokens),
	}
}

func wholeMicros(v float64) int64 {
	if v <= 0 {
		return 0
	}
	return int64(math.Round(v))
}

// PriceMicros computes the 1e-6-unit cost of one usage block. It is observe's
// integer calculator, kept here under the name ticket 10 shipped so the loop
// and its callers have one mouth (see the package comment above).
func PriceMicros(p config.Price, u llm.Usage) int64 {
	return observe.PriceMicros(CardOf(p), TokensOf(u))
}
