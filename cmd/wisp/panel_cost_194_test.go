package main

// Ticket 194's cost door teeth (ledger A1681): the door is wired to the REAL
// reader - memory.Store.ListCostDaily - and the tests below prove the rows
// actually travel, not that "the handler was called". None of them stops at a
// nil check alone:
//
//	Test194CostDoorEndToEndReturnsRealCostDailyRows 正控: bumped rows come back
//	                                                through the chain, values intact.
//	Test194CostDoorEmptyTableAnswersZeroRows        an empty machine answers a true
//	                                                zero, not an error.
//	Test194CostDoorWiringNailInProductionChain      the wiring nail: the assembled
//	                                                production chain carries the reader,
//	                                                so unplugging it (AC 反向) reddens.

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/panel"
)

func auditQuiet194(string, ...any) {}

func costEnvelope194(requestID string) string {
	return `{"method":"` + panel.MethodCostSummary + `","requestId":"` + requestID +
		`","source":"` + panel.ComposerRequestSource + `"}`
}

// costPayload194 is the page's view of the marker line (cost_handlers.go's
// costPayload), declared here so the test parses the WIRE shape, not the
// handler's internal struct.
type costPayload194 struct {
	Days []struct {
		Day        string `json:"day"`
		TokensIn   int64  `json:"tokensIn"`
		TokensOut  int64  `json:"tokensOut"`
		CostMicros int64  `json:"costMicros"`
		Tasks      int64  `json:"tasks"`
	} `json:"days"`
	Count int `json:"count"`
}

// parseCostReply194 splits one Handle reply into the prose sentence and the
// marker payload; the marker line is the only line the page reads.
func parseCostReply194(t *testing.T, reply string) (string, costPayload194) {
	t.Helper()
	var payload costPayload194
	prose, found := "", false
	for _, line := range strings.Split(reply, "\n") {
		if !strings.HasPrefix(line, panel.CostSummaryPayloadMarker) {
			continue
		}
		raw := strings.TrimPrefix(line, panel.CostSummaryPayloadMarker)
		if err := json.Unmarshal([]byte(raw), &payload); err != nil {
			t.Fatalf("cost payload line is not the page's JSON: %v\n%s", err, raw)
		}
		prose, found = strings.TrimSuffix(reply, "\n"+line), true
	}
	if !found {
		t.Fatalf("reply carries no %s line: %q", panel.CostSummaryPayloadMarker, reply)
	}
	return prose, payload
}

// Test194CostDoorEndToEndReturnsRealCostDailyRows is the AC's 正控: two real
// BumpCostDay rows go in one side of the C23 table and come out the cost door
// with their values intact, newest day first. If the adapter stops reading
// ListCostDaily - or the chain stops injecting it - this case reddens.
func Test194CostDoorEndToEndReturnsRealCostDailyRows(t *testing.T) {
	dir := newInboundDataDir33(t, "ask_every_step")
	// Seed the machine's real spending, then close this connection so the
	// chain's own lazy open is the only one the request uses.
	st, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("memory.Open(%s): %v", dir, err)
	}
	ctx := context.Background()
	if err := st.BumpCostDay(ctx, "2026-11-09", 100, 200, 5000); err != nil {
		t.Fatalf("BumpCostDay 11-09: %v", err)
	}
	if err := st.BumpCostDay(ctx, "2026-11-10", 7, 9, 111); err != nil {
		t.Fatalf("BumpCostDay 11-10: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("seed store Close: %v", err)
	}

	disp, _, err := newComposerDispatchChain(dir, auditQuiet194, panelInboundActor, nil, nil, nil)
	if err != nil {
		t.Fatalf("the production chain did not assemble on %s: %v", dir, err)
	}
	reply, err := disp.Handle(ctx, costEnvelope194("rid-194-e2e"))
	if err != nil {
		t.Fatalf("cost.summary refused on a wired chain: %v\nreply=%q", err, reply)
	}

	prose, payload := parseCostReply194(t, reply)
	if payload.Count != 2 {
		t.Fatalf("payload count = %d, want the 2 bumped rows", payload.Count)
	}
	if !strings.Contains(prose, "共 2 行") {
		t.Fatalf("prose does not name the row count: %q", prose)
	}
	// ListCostDaily's own contract: newest day first. The adapter reorders
	// nothing; these two assertions prove the rows came out in that order.
	if len(payload.Days) != 2 || payload.Days[0].Day != "2026-11-10" || payload.Days[1].Day != "2026-11-09" {
		t.Fatalf("day order = %v, want newest first", payload.Days)
	}
	if payload.Days[0].TokensIn != 7 || payload.Days[0].TokensOut != 9 ||
		payload.Days[0].CostMicros != 111 || payload.Days[0].Tasks != 1 {
		t.Fatalf("newest row values drifted: %+v", payload.Days[0])
	}
	if payload.Days[1].TokensIn != 100 || payload.Days[1].TokensOut != 200 ||
		payload.Days[1].CostMicros != 5000 || payload.Days[1].Tasks != 1 {
		t.Fatalf("older row values drifted: %+v", payload.Days[1])
	}
}

// Test194CostDoorEmptyTableAnswersZeroRows: a machine that has run no cost
// task yet answers a true zero - the door reads the table, it does not
// invent rows and it does not error.
func Test194CostDoorEmptyTableAnswersZeroRows(t *testing.T) {
	dir := newInboundDataDir33(t, "ask_every_step")
	disp, _, err := newComposerDispatchChain(dir, auditQuiet194, panelInboundActor, nil, nil, nil)
	if err != nil {
		t.Fatalf("the production chain did not assemble on %s: %v", dir, err)
	}
	reply, err := disp.Handle(context.Background(), costEnvelope194("rid-194-empty"))
	if err != nil {
		t.Fatalf("cost.summary refused on an empty machine: %v\nreply=%q", err, reply)
	}
	_, payload := parseCostReply194(t, reply)
	if payload.Count != 0 || len(payload.Days) != 0 {
		t.Fatalf("empty table answered %+v, want zero rows", payload)
	}
}

// Test194CostDoorWiringNailInProductionChain is the 反向 nail: it reads the
// assembled production chain off the disk (the same literal the CLI seam and
// the resident panel build) and requires the cost socket to carry the real
// reader. Unplug the wiring in newComposerDispatchChain and this case reddens
// before any envelope is sent - which is exactly the AC's reverse test
// (摘注入 ⇒ 红).
func Test194CostDoorWiringNailInProductionChain(t *testing.T) {
	dir := newInboundDataDir33(t, "ask_every_step")
	disp, _, err := newComposerDispatchChain(dir, auditQuiet194, panelInboundActor, nil, nil, nil)
	if err != nil {
		t.Fatalf("the production chain did not assemble on %s: %v", dir, err)
	}
	if disp.Costs == nil {
		t.Fatal("the production chain carries no cost leg: the door is back on the nil socket that refuses by name")
	}
	h, ok := disp.Costs.(*panel.CostReadHandler)
	if !ok {
		t.Fatalf("cost leg is %T, want *panel.CostReadHandler", disp.Costs)
	}
	if h.Reader == nil {
		t.Fatal("CostReadHandler carries no Reader: the door answers nothing but refusals")
	}
	if _, isReal := h.Reader.(*costReader194); !isReal {
		t.Fatalf("Reader is %T, want the ListCostDaily adapter costReader194 (a stub here would be the non-真 answer the AC forbids)", h.Reader)
	}
	// And the refusal's mirror shape stays true: a nil socket must still be a
	// named refusal, never a silence - this is the same property the panel
	// package's nil-socket test pins, read here through the production assembly.
	if !strings.Contains(panel.ErrNoCostReader.Error(), "成本读取腿未接入") {
		t.Fatalf("ErrNoCostReader lost its name: %q (want the refusal that names the missing leg)", panel.ErrNoCostReader)
	}
}
