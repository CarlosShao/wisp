package panel

// Ticket 194 batch 1 - the cost door's own handler (cost.summary), the third
// read-only roster name this ticket lands (SPEC-08's table, owner ruling A381).
//
// Shape and ruling are config_handlers.go's (ticket 248): the leg is
// an interface this package names (no internal/memory import edge), the row
// view carries display facts only, integer micro units only (no float math
// crosses this seam), and the door is read-only by construction - the
// interface declares one read method and nothing else. The reply is the
// settings-route shape: one prose sentence plus one payload line the page
// finds by its marker.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// ErrNoCostReader is the refusal for a cost request on a machine with no cost
// leg attached: named so the page's ask is never a silence.
var ErrNoCostReader = errors.New("panel: 成本读取腿未接入")

// CostDayRow is the display shape of one cost_daily row (C23). Integers only:
// micro units per the schema, currency left to the row's own field, and no
// total is computed here that the schema does not already store - a summary
// that derives a number is a second source of truth.
type CostDayRow struct {
	Day        string `json:"day"`
	TokensIn   int64  `json:"tokensIn"`
	TokensOut  int64  `json:"tokensOut"`
	CostMicros int64  `json:"costMicros"`
	Tasks      int64  `json:"tasks"`
	// CachedTokens is the per-day cached-token count. Readface wave 1
	// (2026-11-17): the cost page's IN/OUT/CACHED column set needs the cached
	// half, and this build's storage does not aggregate it - cost_daily has
	// no cached column (D35's table, which this leg may not widen), so the
	// adapter hands 0 here. The payload's CachedKnown flag is the gate: while
	// it is false the page MUST NOT render this field as a reading - it is
	// "cached tokens are not aggregated by this build's storage", not "zero
	// cached tokens". Inventing a number here would be the fabricated-data
	// failure the 40 page leg's own note refuses on its side too.
	CachedTokens int64 `json:"cachedTokens"`
}

// CostAggregateRow is the aggregate face of the cost read (ticket 40 AC#3's
// Go half, readface wave 1): the month-to-date totals, the [cost] budgets and
// their verdict, and the pricing version the meter priced against. Integers
// only; the budget comparison itself happens in internal/observe (basis
// points, no floats touch money) and arrives here as a word plus the two
// integer budgets.
type CostAggregateRow struct {
	// DayCostMicros is today's own cost_daily row (0 = no tasks today).
	DayCostMicros int64 `json:"dayCostMicros"`
	// Month* are the month-to-date aggregates over cost_daily, summed from
	// the store's own rows (never re-derived per model or guessed from days
	// the reader did not fetch).
	MonthTokensIn   int64 `json:"monthTokensIn"`
	MonthTokensOut  int64 `json:"monthTokensOut"`
	MonthCostMicros int64 `json:"monthCostMicros"`
	// DailyBudgetMicros / MonthlyBudgetMicros are the [cost] caps as
	// configured; 0 = uncapped (the schema's own "0 = no budget" word).
	DailyBudgetMicros   int64 `json:"dailyBudgetMicros"`
	MonthlyBudgetMicros int64 `json:"monthlyBudgetMicros"`
	// BudgetState is the three-state word the page renders the budget bar
	// from (ticket 40 AC#4): "normal" (<80%), "warn" (>= the alert
	// threshold, 80% by default), "paused" (>=100% AND over_budget=pause,
	// i.e. new tasks are actually being refused). The mapping from
	// observe.BudgetVerdict lives in the adapter (cmd/wisp), which owns the
	// config read; the handler carries the word through untouched.
	BudgetState string `json:"budgetState"`
	// BudgetAllowNewTask mirrors observe.BudgetVerdict.AllowNewTask: the
	// machine-checkable half of the state (false only at paused).
	BudgetAllowNewTask bool `json:"budgetAllowNewTask"`
	// PricingVersion / PricingUpdated identify the price table the meter
	// priced against (observe's embedded pricing.json version + date): the
	// cost page's pricing-version footnote. 0/"" = the leg could not read the
	// table, which the page renders as "unknown", never as a version.
	PricingVersion int    `json:"pricingVersion"`
	PricingUpdated string `json:"pricingUpdated"`
}

// CostAggregateReader is the OPTIONAL second face of the cost leg: an adapter
// that can also answer the aggregate read. It is a separate interface (not a
// second method on CostReader) so the existing CostReader implementations and
// test doubles keep compiling untouched - the same shape the models route's
// interface took - and the handler asks for it with a type assertion: a reader
// that answers days only still answers cost.summary, with the aggregate marked
// known=false.
type CostAggregateReader interface {
	CostAggregate(ctx context.Context) (CostAggregateRow, error)
}

// CostReader is the injected cost leg: one read method, no third. nil reader =
// every cost request refused with a name.
type CostReader interface {
	// CostDays answers "what did this machine spend": per-day rows, newest
	// first, in the leg's own order (the adapter owns the ordering contract).
	CostDays(ctx context.Context) ([]CostDayRow, error)
}

// CostReadHandler is the native handler for the cost door.
type CostReadHandler struct {
	// Reader is the single leg. nil = every cost request refused.
	Reader CostReader
	// Audit is the project's existing "[audit]" sink, as in every other handler.
	Audit AuditFunc
	// Actor names who triggered it, for the audit line. Empty stays "unknown".
	Actor string
}

// HandleCostSummary answers one parsed cost request and returns the sentence
// the page is shown. The request's identity fields are read and never
// assigned; the door takes no selector (cost.summary in the spec table carries
// none), so anything selector-shaped on the envelope is refused rather than
// silently ignored.
func (h *CostReadHandler) HandleCostSummary(ctx context.Context, req ComposerRequest) (string, error) {
	if h == nil || h.Reader == nil {
		err := fmt.Errorf("%w: 成本请求 %q（requestId=%q）没有读取腿可用，未发生任何变更",
			ErrNoCostReader, req.Method, req.RequestID)
		h.record(req, err, "no-leg")
		return "", err
	}
	if violations := foreignSelectors(req); violations != "" {
		err := fmt.Errorf("%w: cost.summary 不接受选择子（带了 %s）（requestId=%q）",
			ErrSelectorRefused, violations, req.RequestID)
		h.record(req, err, "foreign-selector")
		return "", err
	}
	days, err := h.Reader.CostDays(ctx)
	if err != nil {
		refused := fmt.Errorf("panel: 成本读取失败（requestId=%q）：%w", req.RequestID, err)
		h.record(req, refused, "read-failed")
		return "", refused
	}
	if days == nil {
		days = []CostDayRow{}
	}
	// The aggregate face is optional (CostAggregateReader): a reader that
	// answers it extends the payload; one that answers days only still gets
	// its answer through, with the aggregate marked known=false - a partial
	// answer that SAYS it is partial, never a silence.
	var agg *CostAggregateRow
	aggregateKnown := false
	if aggReader, ok := h.Reader.(CostAggregateReader); ok {
		row, aggErr := aggReader.CostAggregate(ctx)
		if aggErr == nil {
			agg = &row
			aggregateKnown = true
		} else {
			h.record(req, aggErr, "aggregate-failed")
		}
	}
	h.record(req, nil, "answered")
	var b strings.Builder
	fmt.Fprintf(&b, "成本按天共 %d 行（新的一天在前）。", len(days))
	if !aggregateKnown {
		b.WriteString(" 本月聚合/预算读数未取到（读取腿没有接聚合面，或读取失败）。")
	} else {
		fmt.Fprintf(&b, " 本月累计 %d micro；预算状态 %s。", agg.MonthCostMicros, agg.BudgetState)
	}
	b.WriteString("\n" + CostSummaryPayloadMarker + marshalCostPayload(costPayload{
		Days:           days,
		Count:          len(days),
		CachedKnown:    false,
		AggregateKnown: aggregateKnown,
		Aggregate:      agg,
	}))
	return b.String(), nil
}

// ErrSelectorRefused is the refusal for a filter a no-selector door was handed.
// It is shared vocabulary on purpose: the read-only doors this ticket measured
// (see bridge.go's roster note for the ones not landed today) would refuse with
// the same one sentence, "this door takes no filter".
var ErrSelectorRefused = errors.New("panel: 请求的选择子不成立")

// foreignSelectors names which envelope selector fields are filled on a
// no-selector door. Any filled selector is a filter the door never offered,
// and the returned string is the refusal's evidence list.
func foreignSelectors(req ComposerRequest) string {
	var filled []string
	add := func(name, v string) {
		if strings.TrimSpace(v) != "" {
			filled = append(filled, name)
		}
	}
	add("to", req.To)
	add("path", req.Path)
	add("text", req.Text)
	add("configField", req.ConfigField)
	add("configValue", req.ConfigValue)
	return strings.Join(filled, ", ")
}

// CostSummaryPayloadMarker starts the cost read payload. Same two properties
// as every marker in this family: no `panel.` prefix and no dot (it is neither
// a method name nor a route).
const CostSummaryPayloadMarker = "wisp-cost-summary:"

// costPayload is the cost door's wire shape, separate from the view type so a
// field added to the view does not silently change the page's JSON contract.
//
// Readface wave 1 (2026-11-17, ticket 40 AC#3/AC#4's Go half): the payload
// grows the aggregate/budget/pricing block (optional, gated by AggregateKnown)
// and the cached-token gate. CachedKnown is FALSE on this build by
// construction: cost_daily carries no cached-token column and this leg may
// not widen the D35 table, so the gate travels with the rows and the page
// renders the CACHED column as "未聚合" while it is false - never as the
// rows' zero. Aggregate is nil unless the leg answered; omitempty keeps the
// old payload shape byte-identical for a days-only leg.
type costPayload struct {
	Days  []CostDayRow `json:"days"`
	Count int          `json:"count"`
	// CachedKnown gates the per-row CachedTokens field: false = this build's
	// storage does not aggregate cached tokens (see CostDayRow).
	CachedKnown bool `json:"cachedKnown"`
	// AggregateKnown gates Aggregate: false = the leg has no aggregate face
	// or its read failed (the prose says so too); true = Aggregate carries
	// the month totals, budgets, verdict and pricing version.
	AggregateKnown bool              `json:"aggregateKnown"`
	Aggregate      *CostAggregateRow `json:"aggregate,omitempty"`
}

// marshalCostPayload renders the payload line; the payloadError branch is the
// same never-silent shape the settings payload carries.
func marshalCostPayload(p costPayload) string {
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Sprintf("%s{\"payloadError\":%q}", CostSummaryPayloadMarker, err.Error())
	}
	return string(data)
}

func (h *CostReadHandler) actor() string {
	if h == nil || h.Actor == "" {
		return "unknown"
	}
	return h.Actor
}

// record writes one audit line per outcome; a nil sink silences the line,
// never the refusal.
func (h *CostReadHandler) record(req ComposerRequest, err error, outcome string) {
	if h == nil || h.Audit == nil {
		return
	}
	line := "ok"
	if err != nil {
		line = err.Error()
	}
	h.Audit("panel: COST-ROUTE request=%q method=%q source=%q actor=%q outcome=%q err=%s",
		req.RequestID, req.Method, req.Source, h.actor(), outcome, line)
}
