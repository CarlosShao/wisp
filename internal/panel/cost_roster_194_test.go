package panel

// Ticket 194's resident tests for the one roster name this batch landed
// (cost.summary) plus the boundary records for the three spellings it measured
// and could NOT land (see bridge.go's roster note: approval.current /
// approval.queue by the frozen grant-boundary instrument's vocabulary,
// tasks.list / task.detail by internal/tools' ticket-188 task-axis scan, whose
// package is outside this ticket's write surface).
//
// What each block is FOR, so a later reader can re-derive rather than trust:
//
//   - the whitelist/routing block is AC#4's clause 1 (an off-roster name is
//     refused AND audited, never silently dropped) applied to the name this
//     batch added, plus the honest states of the doors themselves: a nil socket
//     refuses BY NAME (the router's own rule), and the three blocked spec
//     spellings stay refused here so nobody "completes" them under an alias
//     without a fresh owner ruling.
//   - the handler block is AC#4's clause 3 (attribution is read, never
//     rewritten - the audit line carries the envelope's own requestId verbatim)
//     and clause 2's negative space (the door carries no decision path at all:
//     a verdict-shaped field on the envelope is dropped by the shared parse and
//     the reader interface has no write to reach).
//
// The reader here is a fake of the interface cost_handlers.go declares - the
// same seam the next cmd/wisp wave fills with the real memory adapter, so these
// tests stay green on both sides of that wiring.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
)

// auditSink194 collects the lines one audit sink was handed, with the sink's
// nil-safety free of charge (a nil sink silences the line, never the refusal -
// asserted below by simply not attaching one).
type auditSink194 struct {
	lines []string
}

func (s *auditSink194) sink(format string, args ...any) {
	s.lines = append(s.lines, fmt.Sprintf(format, args...))
}

func (s *auditSink194) joined() string { return strings.Join(s.lines, "\n") }

// fakeCostReader194 is the fake cost leg: one read method, no third - the
// interface shape IS the no-write statement (AC#4 clause 2's negative space).
type fakeCostReader194 struct {
	rows []CostDayRow
	err  error

	daysSeen int
}

func (f *fakeCostReader194) CostDays(context.Context) ([]CostDayRow, error) {
	f.daysSeen++
	if f.err != nil {
		return nil, f.err
	}
	return f.rows, nil
}

func costDispatch194(costs CostSummaryHandler, audit AuditFunc) *ComposerDispatch {
	return &ComposerDispatch{Costs: costs, Audit: audit}
}

// TestRoster194CostDoorIsWhitelistedAndRouted is the landed name's end-to-end
// behaviour: the guard answers it, parse accepts it with its identity fields, a
// nil socket refuses BY NAME (never silently), and an attached socket routes to
// the handler that owns the door.
func TestRoster194CostDoorIsWhitelistedAndRouted(t *testing.T) {
	if !knownComposerMethod(MethodCostSummary) {
		t.Fatal("premise: the running guard refuses cost.summary - this test would be measuring a door that does not exist")
	}
	req, err := ParseComposerRequest(okEnvelope(t, MethodCostSummary, nil))
	if err != nil {
		t.Fatalf("cost.summary refused at parse: %v", err)
	}
	if req.Source != ComposerRequestSource || req.RequestID == "" {
		t.Fatalf("cost.summary parsed without its identity fields: %+v", req)
	}

	t.Run("a nil socket refuses the listed name, audited, with nothing answered", func(t *testing.T) {
		sink := &auditSink194{}
		d := costDispatch194(nil, sink.sink)
		raw := okEnvelope(t, MethodCostSummary, nil)
		reply, err := d.Handle(context.Background(), raw)
		if err == nil || !errors.Is(err, ErrNoHandlerAttached) {
			t.Fatalf("cost.summary on a nil socket = (%q, %v), want ErrNoHandlerAttached", reply, err)
		}
		if reply == "" || !strings.Contains(reply, "cost.summary") {
			t.Fatalf("refusal %q does not name the door", reply)
		}
		parsed, _ := ParseComposerRequest(raw)
		if !strings.Contains(sink.joined(), parsed.RequestID) || !strings.Contains(sink.joined(), "cost.summary") {
			t.Fatalf("refusal audit line %q does not carry the request id and method (AC#4 clause 1)", sink.joined())
		}
	})

	t.Run("an attached socket routes to the cost handler", func(t *testing.T) {
		reader := &fakeCostReader194{rows: []CostDayRow{{Day: "2026-11-10", Tasks: 1}}}
		sink := &auditSink194{}
		d := costDispatch194(&CostReadHandler{Reader: reader, Audit: sink.sink, Actor: "test"}, sink.sink)
		if _, err := d.Handle(context.Background(), okEnvelope(t, MethodCostSummary, nil)); err != nil {
			t.Fatalf("cost.summary on an attached socket errored: %v", err)
		}
		if reader.daysSeen != 1 {
			t.Fatalf("reader asked %d times, want 1", reader.daysSeen)
		}
	})

	t.Run("the three blocked spec spellings stay refused - no alias lands them", func(t *testing.T) {
		// Two frozen instruments refuse these spellings as answered routes, and
		// neither instrument is in this ticket's write surface:
		//   - approval.current / approval.queue: l2_grant_boundary_test.go's
		//     grantRouteWords flags every answered route carrying an approval
		//     word (D33/F2, ban #6) - a spec-vs-instrument conflict, owner item.
		//   - tasks.list / task.detail: internal/tools' ticket-188 task-axis
		//     scan (task_state_188_test.go) flags every Method* constant whose
		//     value contains "task"/"state" as a panel-side write leg.
		// These refusals are the standing record of that boundary.
		for _, m := range []string{"approval.current", "approval.queue", "tasks.list", "task.detail"} {
			raw := okEnvelope(t, m, nil)
			if _, err := ParseComposerRequest(raw); err == nil || !errors.Is(err, ErrComposerRequest) {
				t.Fatalf("%s parsed - a route the frozen instruments ban must refuse at the door, got %v", m, err)
			}
		}
	})
}

// TestRoster194OffRosterNameIsRefusedAndAudited is AC#4 clause 1 in its own
// right: a name outside the whitelist is refused at parse, the refusal names
// the request, and the audit line exists even though nothing was routed.
func TestRoster194OffRosterNameIsRefusedAndAudited(t *testing.T) {
	sink := &auditSink194{}
	d := costDispatch194(nil, sink.sink)
	raw := okEnvelope(t, "cost.histogram", nil)
	reply, err := d.Handle(context.Background(), raw)
	if err == nil || !errors.Is(err, ErrComposerRequest) {
		t.Fatalf("cost.histogram = (%q, %v), want a parse refusal", reply, err)
	}
	parsed, _ := ParseComposerRequest(raw)
	if !strings.Contains(reply, parsed.RequestID) || !strings.Contains(reply, "cost.histogram") {
		t.Fatalf("refusal %q does not name what was refused", reply)
	}
	if !strings.Contains(sink.joined(), parsed.RequestID) || !strings.Contains(sink.joined(), "cost.histogram") {
		t.Fatalf("off-roster refusal not audited: %q", sink.joined())
	}
}

// TestCostDoorHandlerGatesAndAnswers covers the cost door's own gates - no leg,
// foreign selector, leg error: each a named refusal with an audit line - and
// the payload contract (prose plus ONE marker line whose JSON renders an empty
// history as [], never null).
func TestCostDoorHandlerGatesAndAnswers(t *testing.T) {
	ctx := context.Background()

	t.Run("no leg is a named refusal, not a silence", func(t *testing.T) {
		sink := &auditSink194{}
		h := &CostReadHandler{Reader: nil, Audit: sink.sink, Actor: "test"}
		_, err := h.HandleCostSummary(ctx, mustParseCost194(t))
		if err == nil || !errors.Is(err, ErrNoCostReader) {
			t.Fatalf("no-leg cost request = %v, want ErrNoCostReader", err)
		}
		if !strings.Contains(sink.joined(), "no-leg") {
			t.Fatalf("no-leg refusal not audited with its outcome: %q", sink.joined())
		}
	})

	t.Run("a filter the door never offered is refused, not ignored", func(t *testing.T) {
		raw := okEnvelope(t, MethodCostSummary, map[string]any{"configField": "x"})
		req, err := ParseComposerRequest(raw)
		if err != nil {
			t.Fatalf("envelope refused: %v", err)
		}
		h := &CostReadHandler{Reader: &fakeCostReader194{}, Audit: nil, Actor: "test"}
		_, err = h.HandleCostSummary(ctx, req)
		if err == nil || !errors.Is(err, ErrSelectorRefused) {
			t.Fatalf("cost.summary with a filter = %v, want ErrSelectorRefused", err)
		}
	})

	t.Run("a leg error is a state the user is told about, not an empty list", func(t *testing.T) {
		sink := &auditSink194{}
		h := &CostReadHandler{Reader: &fakeCostReader194{err: errors.New("cost_daily: unreadable")}, Audit: sink.sink, Actor: "test"}
		_, err := h.HandleCostSummary(ctx, mustParseCost194(t))
		if err == nil {
			t.Fatal("an unreadable cost store answered as an empty summary")
		}
		if !strings.Contains(err.Error(), "unreadable") {
			t.Fatalf("refusal %v does not name the leg's own error", err)
		}
	})

	t.Run("the payload is one marker line whose JSON renders empty as []", func(t *testing.T) {
		h := &CostReadHandler{Reader: &fakeCostReader194{}, Audit: nil, Actor: "test"}
		reply, err := h.HandleCostSummary(ctx, mustParseCost194(t))
		if err != nil {
			t.Fatalf("an empty history answered as an error: %v", err)
		}
		payload := markerPayload194(t, reply, CostSummaryPayloadMarker)
		days, ok := payload["days"].([]any)
		if !ok || days == nil {
			t.Fatalf("payload days = %v, want a [] (an empty list is a state, never null)", payload["days"])
		}
		if len(days) != 0 {
			t.Fatalf("payload days = %v, want empty", days)
		}
	})

	t.Run("a filled history renders its rows", func(t *testing.T) {
		h := &CostReadHandler{Reader: &fakeCostReader194{rows: []CostDayRow{
			{Day: "2026-11-09", TokensIn: 10, TokensOut: 5, CostMicros: 100, Tasks: 2},
		}}, Audit: nil, Actor: "test"}
		reply, err := h.HandleCostSummary(ctx, mustParseCost194(t))
		if err != nil {
			t.Fatalf("cost summary errored: %v", err)
		}
		payload := markerPayload194(t, reply, CostSummaryPayloadMarker)
		if payload["count"] != float64(1) {
			t.Fatalf("payload count = %v, want 1", payload["count"])
		}
		days := payload["days"].([]any)
		if len(days) != 1 || days[0].(map[string]any)["day"] != "2026-11-09" {
			t.Fatalf("payload days = %v, want the one row the leg returned", days)
		}
	})
}

// TestCostDoorNeverRewritesAttribution is AC#4 clause 3: the request id and
// source travel through the hop untouched, and the audit line carries the
// envelope's own id - the observable field is RequestID, per A380.
func TestCostDoorNeverRewritesAttribution(t *testing.T) {
	sink := &auditSink194{}
	reader := &fakeCostReader194{rows: []CostDayRow{{Day: "2026-11-10"}}}
	h := &CostReadHandler{Reader: reader, Audit: sink.sink, Actor: "test"}
	req, err := ParseComposerRequest(okEnvelope(t, MethodCostSummary, nil))
	if err != nil {
		t.Fatalf("cost.summary refused: %v", err)
	}
	if _, err := h.HandleCostSummary(context.Background(), req); err != nil {
		t.Fatalf("cost.summary errored: %v", err)
	}
	if !strings.Contains(sink.joined(), req.RequestID) {
		t.Fatalf("audit %q does not carry the envelope's own requestId (attribution rewritten or dropped)", sink.joined())
	}
	if !strings.Contains(sink.joined(), ComposerRequestSource) {
		t.Fatalf("audit %q does not carry the source claim", sink.joined())
	}
}

// TestCostDoorCarriesNoDecisionPath is AC#4 clause 2's negative space: a
// verdict-shaped field on the envelope is dropped by the shared parse (the
// ComposerRequest has no field for it), and the read-only interface gives the
// handler nothing a decision could ride through - measured here by the reader
// being asked for exactly its one read and nothing else.
func TestCostDoorCarriesNoDecisionPath(t *testing.T) {
	reader := &fakeCostReader194{rows: []CostDayRow{{Day: "2026-11-10"}}}
	h := &CostReadHandler{Reader: reader, Audit: nil, Actor: "test"}
	// outcome/decision-shaped extras on the envelope: the parse drops them (no
	// field binds them), the door answers its read, and no write exists to land
	// them in.
	req, err := ParseComposerRequest(okEnvelope(t, MethodCostSummary, map[string]any{
		"outcome": "grant", "decision": "allow",
	}))
	if err != nil {
		t.Fatalf("cost.summary with verdict-shaped extras refused at parse: %v", err)
	}
	if req.To != "" || req.ConfigField != "" {
		t.Fatalf("verdict-shaped extras bound to selector fields: %+v", req)
	}
	if _, err := h.HandleCostSummary(context.Background(), req); err != nil {
		t.Fatalf("the read door errored on dropped extras: %v", err)
	}
	if reader.daysSeen != 1 {
		t.Fatalf("reader asked %d times, want exactly the one read", reader.daysSeen)
	}
}

// mustParseCost194 parses one cost envelope.
func mustParseCost194(t *testing.T) ComposerRequest {
	t.Helper()
	req, err := ParseComposerRequest(okEnvelope(t, MethodCostSummary, nil))
	if err != nil {
		t.Fatalf("cost.summary refused: %v", err)
	}
	return req
}

// markerPayload194 splits a receipt at its marker line and parses what follows.
// The payload is ONE line beginning with the marker; the prose never contains
// it (the renderer appends it last).
func markerPayload194(t *testing.T, reply, marker string) map[string]any {
	t.Helper()
	idx := strings.Index(reply, marker)
	if idx < 0 {
		t.Fatalf("reply %q carries no %s payload line", reply, marker)
	}
	raw := reply[idx+len(marker):]
	if i := strings.IndexByte(raw, '\n'); i >= 0 {
		raw = raw[:i]
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("payload %q does not parse: %v", raw, err)
	}
	return out
}
