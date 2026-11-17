package panel

// Readface wave 1 (2026-11-17) behavior teeth, panel half:
//
//   - ticket 37 AC#5: the aggregate-card fields ride a NativeVerdict's Batch
//     onto the card and the wire, and a non-batch card's JSON stays free of
//     them;
//   - ticket 39 AC#1/AC#4: the config.get payload carries the full section
//     tree and the line-numbered load errors a leg hands it, and omits both
//     keys when the leg answered neither;
//   - ticket 39 AC#1's reload event: an accepted config.set receipt carries
//     the machine event line; a refused write emits none;
//   - ticket 40 AC#3/AC#4: the cost payload carries the aggregate block when
//     the reader has the aggregate face, and says "known=false" when it does
//     not - never a fabricated number.

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// ---- ticket 37 AC#5: the aggregate-card fields ----

func TestBatchVerdictCardCarriesAggregateFields(t *testing.T) {
	v := NativeVerdict{
		CorrelationID: "corr-batch-37",
		Tool:          "fs.write",
		Level:         risk.L2,
		RulesHit:      []risk.RuleID{risk.R7},
		Reason:        "batch",
		Batch: &NativeBatchView{
			Total:        10,
			AffectedDirs: []string{"C:/data/x", "C:/data/y"},
			First:        []string{"C:/data/x/a.txt", "C:/data/x/b.txt"},
			Deferred:     8,
			Summary:      "本次调用包含 10 项同类 L2 操作",
		},
	}
	card := v.CardView()
	if card.BatchTotal != 10 || card.BatchDeferred != 8 || card.BatchSummary == "" {
		t.Errorf("card batch fields = %d/%d/%q, want the verdict's own numbers through",
			card.BatchTotal, card.BatchDeferred, card.BatchSummary)
	}
	if len(card.BatchDirs) != 2 || len(card.BatchFirst) != 2 {
		t.Errorf("card batch lists = %v / %v, want the verdict's own copies", card.BatchDirs, card.BatchFirst)
	}
	if !card.HasBatch() {
		t.Error("HasBatch() = false for a card that carries a batch")
	}
	// The wire carries the five keys (they are ledger-registered judged gaps
	// until the page declares them) and the numbers on it are the card's.
	data, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal card: %v", err)
	}
	var wire map[string]any
	if err := json.Unmarshal(data, &wire); err != nil {
		t.Fatalf("unmarshal card: %v", err)
	}
	for _, k := range []string{"batchTotal", "batchDirs", "batchFirst", "batchDeferred", "batchSummary"} {
		if _, ok := wire[k]; !ok {
			t.Errorf("wire carries no %q key although this card is about a batch", k)
		}
	}
	if got, _ := wire["batchTotal"].(float64); int(got) != 10 {
		t.Errorf("wire batchTotal = %v, want 10", wire["batchTotal"])
	}
}

func TestPlainCardWireCarriesNoBatchKeys(t *testing.T) {
	card := CardViewFromDecision(ApprovalSubject{CorrelationID: "c-plain", Tool: "fs.write"},
		risk.Decision{Level: risk.L2, RulesHit: []risk.RuleID{risk.R2}, Reason: "plain"})
	data, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, k := range []string{"batchTotal", "batchDirs", "batchFirst", "batchDeferred", "batchSummary"} {
		if strings.Contains(string(data), `"`+k+`"`) {
			t.Errorf("a plain card's wire carries %q - a non-batch card must not wear batch fields", k)
		}
	}
	if card.HasBatch() {
		t.Error("HasBatch() = true for a card with no batch")
	}
}

// ---- ticket 39 AC#1/AC#4: sections + errors on the config.get payload ----

type readfaceStore456 struct {
	view   SettingsView
	apply  SettingWriteResult
	applyE error
}

func (s *readfaceStore456) ReadSettings(context.Context) (SettingsView, error) { return s.view, nil }
func (s *readfaceStore456) ApplySetting(context.Context, SettingWrite) (SettingWriteResult, error) {
	return s.apply, s.applyE
}
func (s *readfaceStore456) StoreCredential(context.Context, string) (SettingWriteResult, error) {
	return SettingWriteResult{}, nil
}

func readfaceRequest456(method, field, value string) ComposerRequest {
	return ComposerRequest{
		Method: method, RequestID: "rid-456", Source: ComposerRequestSource,
		ConfigField: field, ConfigValue: value,
	}
}

func TestConfigGetPayloadCarriesSectionsAndErrors(t *testing.T) {
	h := &ConfigWriteHandler{Store: &readfaceStore456{view: SettingsView{
		Readable: true,
		Sections: []ConfigSectionRow{
			{Key: "ball.size", Value: "90", Tier: "hot"},
			{Key: "risk.permission_mode", Value: "ask_every_step", Tier: "locked"},
		},
		LoadErrors: []ConfigLoadError{{Line: 5, Message: `unknown key "bogus" at line 5`}},
	}}}
	reply, err := h.HandleConfigRequest(context.Background(), readfaceRequest456(MethodConfigGet, "", ""), "")
	if err != nil {
		t.Fatalf("config.get refused: %v", err)
	}
	var payloadLine string
	for _, line := range strings.Split(reply, "\n") {
		if strings.HasPrefix(line, SettingsPayloadMarker) {
			payloadLine = strings.TrimPrefix(line, SettingsPayloadMarker)
		}
	}
	if payloadLine == "" {
		t.Fatalf("config.get reply carries no %s line: %q", SettingsPayloadMarker, reply)
	}
	var payload struct {
		Sections []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
			Tier  string `json:"tier"`
		} `json:"sections"`
		Errors []struct {
			Line    int    `json:"line"`
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.Unmarshal([]byte(payloadLine), &payload); err != nil {
		t.Fatalf("payload line is not JSON: %v\n%s", err, payloadLine)
	}
	if len(payload.Sections) != 2 {
		t.Fatalf("payload sections = %d rows, want 2", len(payload.Sections))
	}
	if payload.Sections[0].Key != "ball.size" || payload.Sections[0].Value != "90" || payload.Sections[0].Tier != "hot" {
		t.Errorf("sections row 0 = %+v, want the leg's own row", payload.Sections[0])
	}
	if payload.Sections[1].Tier != "locked" {
		t.Errorf("locked section tier = %q, want the registry's locked marker", payload.Sections[1].Tier)
	}
	if len(payload.Errors) != 1 || payload.Errors[0].Line != 5 || !strings.Contains(payload.Errors[0].Message, "bogus") {
		t.Errorf("payload errors = %+v, want the line-numbered row through", payload.Errors)
	}
}

func TestConfigGetPayloadOmitsUnansweredFaces(t *testing.T) {
	h := &ConfigWriteHandler{Store: &readfaceStore456{view: SettingsView{Readable: true}}}
	reply, err := h.HandleConfigRequest(context.Background(), readfaceRequest456(MethodConfigGet, "", ""), "")
	if err != nil {
		t.Fatalf("config.get refused: %v", err)
	}
	var payloadLine string
	for _, line := range strings.Split(reply, "\n") {
		if strings.HasPrefix(line, SettingsPayloadMarker) {
			payloadLine = strings.TrimPrefix(line, SettingsPayloadMarker)
		}
	}
	if strings.Contains(payloadLine, `"sections"`) || strings.Contains(payloadLine, `"errors"`) {
		t.Errorf("a leg that answered neither face still put the keys on the wire: %s", payloadLine)
	}
}

// ---- ticket 39 AC#1: the reload event on an accepted write ----

func TestAcceptedWriteCarriesReloadEvent(t *testing.T) {
	h := &ConfigWriteHandler{Store: &readfaceStore456{
		apply: SettingWriteResult{Field: "ball_size", Written: []string{"ball.size", "ball.opacity_idle"}, Tier: EffectiveRestart},
	}}
	reply, err := h.HandleConfigRequest(context.Background(),
		readfaceRequest456(MethodConfigSet, "ball_size", "90"), "")
	if err != nil {
		t.Fatalf("config.set refused: %v", err)
	}
	var eventLine string
	for _, line := range strings.Split(reply, "\n") {
		if strings.HasPrefix(line, ConfigReloadEventMarker) {
			eventLine = strings.TrimPrefix(line, ConfigReloadEventMarker)
		}
	}
	if eventLine == "" {
		t.Fatalf("accepted write carries no %s event line: %q", ConfigReloadEventMarker, reply)
	}
	var ev struct {
		Event    string   `json:"event"`
		Tier     string   `json:"tier"`
		Written  []string `json:"written"`
		Sections []string `json:"sections"`
	}
	if err := json.Unmarshal([]byte(eventLine), &ev); err != nil {
		t.Fatalf("event line is not JSON: %v\n%s", err, eventLine)
	}
	if ev.Event != "config.reload" {
		t.Errorf("event = %q, want config.reload", ev.Event)
	}
	if ev.Tier != string(EffectiveRestart) {
		t.Errorf("tier = %q, want %q", ev.Tier, EffectiveRestart)
	}
	if len(ev.Written) != 2 {
		t.Errorf("written = %v, want the write result's own key paths", ev.Written)
	}
	if len(ev.Sections) != 1 || ev.Sections[0] != "ball" {
		t.Errorf("sections = %v, want [ball] (the distinct top-level segments)", ev.Sections)
	}
}

func TestRefusedWriteCarriesNoReloadEvent(t *testing.T) {
	h := &ConfigWriteHandler{Store: &readfaceStore456{applyE: errors.New("no leg")}}
	// A refusal before the store: unknown field.
	reply, err := h.HandleConfigRequest(context.Background(),
		readfaceRequest456(MethodConfigSet, "not_a_field", "x"), "")
	if err == nil {
		t.Fatal("unknown field must refuse")
	}
	if strings.Contains(reply, ConfigReloadEventMarker) {
		t.Errorf("a refused write emitted a reload event: %q (nothing changed, nothing reloads)", reply)
	}
}

// ---- ticket 40 AC#3/AC#4: the cost aggregate face ----

type aggReader456 struct{ err error }

func (r *aggReader456) CostDays(context.Context) ([]CostDayRow, error) {
	return []CostDayRow{{Day: "2026-11-17", TokensIn: 10, TokensOut: 2, CostMicros: 300, Tasks: 1}}, nil
}
func (r *aggReader456) CostAggregate(context.Context) (CostAggregateRow, error) {
	if r.err != nil {
		return CostAggregateRow{}, r.err
	}
	return CostAggregateRow{
		MonthCostMicros:     900,
		DailyBudgetMicros:   1000,
		MonthlyBudgetMicros: 2000,
		BudgetState:         "warn",
		BudgetAllowNewTask:  true,
		PricingVersion:      3,
		PricingUpdated:      "2026-11-01",
	}, nil
}

type daysOnlyReader456 struct{}

func (r *daysOnlyReader456) CostDays(context.Context) ([]CostDayRow, error) {
	return []CostDayRow{{Day: "2026-11-17", CostMicros: 300, Tasks: 1}}, nil
}

func costPayloadMap456(t *testing.T, reply string) map[string]any {
	t.Helper()
	for _, line := range strings.Split(reply, "\n") {
		if strings.HasPrefix(line, CostSummaryPayloadMarker) {
			var m map[string]any
			if err := json.Unmarshal([]byte(strings.TrimPrefix(line, CostSummaryPayloadMarker)), &m); err != nil {
				t.Fatalf("cost payload is not JSON: %v", err)
			}
			return m
		}
	}
	t.Fatalf("cost reply carries no %s line: %q", CostSummaryPayloadMarker, reply)
	return nil
}

func costRequest456() ComposerRequest {
	return ComposerRequest{Method: MethodCostSummary, RequestID: "rid-cost-456", Source: ComposerRequestSource}
}

func TestCostSummaryCarriesAggregateWhenReaderHasTheFace(t *testing.T) {
	h := &CostReadHandler{Reader: &aggReader456{}}
	reply, err := h.HandleCostSummary(context.Background(), costRequest456())
	if err != nil {
		t.Fatalf("cost.summary refused: %v", err)
	}
	m := costPayloadMap456(t, reply)
	if known, _ := m["aggregateKnown"].(bool); !known {
		t.Errorf("aggregateKnown = %v, want true", m["aggregateKnown"])
	}
	agg, _ := m["aggregate"].(map[string]any)
	if agg == nil {
		t.Fatal("aggregate block missing from the payload")
	}
	if got, _ := agg["monthCostMicros"].(float64); int(got) != 900 {
		t.Errorf("aggregate monthCostMicros = %v, want 900", agg["monthCostMicros"])
	}
	if got, _ := agg["budgetState"].(string); got != "warn" {
		t.Errorf("aggregate budgetState = %q, want the leg's own word", got)
	}
	if got, _ := agg["pricingVersion"].(float64); int(got) != 3 {
		t.Errorf("aggregate pricingVersion = %v, want 3", agg["pricingVersion"])
	}
	if known, _ := m["cachedKnown"].(bool); known {
		t.Error("cachedKnown = true: this build's storage aggregates no cached tokens and must not claim it")
	}
}

func TestCostSummarySaysKnownFalseWithoutTheFace(t *testing.T) {
	h := &CostReadHandler{Reader: &daysOnlyReader456{}}
	reply, err := h.HandleCostSummary(context.Background(), costRequest456())
	if err != nil {
		t.Fatalf("cost.summary refused: %v (a days-only reader still answers)", err)
	}
	m := costPayloadMap456(t, reply)
	if known, _ := m["aggregateKnown"].(bool); known {
		t.Errorf("aggregateKnown = true for a days-only reader: the payload would claim a read nobody did")
	}
	if _, ok := m["aggregate"]; ok {
		t.Error("payload carries an aggregate block although the reader answered none")
	}
}

func TestCostSummaryAggregateFailureIsSaidNotSilent(t *testing.T) {
	h := &CostReadHandler{Reader: &aggReader456{err: errors.New("store closed")}}
	reply, err := h.HandleCostSummary(context.Background(), costRequest456())
	if err != nil {
		t.Fatalf("cost.summary refused: %v (the days half still answered)", err)
	}
	m := costPayloadMap456(t, reply)
	if known, _ := m["aggregateKnown"].(bool); known {
		t.Error("aggregateKnown = true although the aggregate read failed")
	}
	if !strings.Contains(reply, "未取到") {
		t.Errorf("the prose does not say the aggregate was not read: %q", reply)
	}
}
