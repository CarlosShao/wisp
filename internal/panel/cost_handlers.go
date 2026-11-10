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
	h.record(req, nil, "answered")
	var b strings.Builder
	fmt.Fprintf(&b, "成本按天共 %d 行（新的一天在前）。", len(days))
	b.WriteString("\n" + CostSummaryPayloadMarker + marshalCostPayload(costPayload{
		Days:  days,
		Count: len(days),
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
type costPayload struct {
	Days  []CostDayRow `json:"days"`
	Count int          `json:"count"`
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
