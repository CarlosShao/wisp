package agent

// Ticket 179: Loop.decideRisk only knew two declared levels (L1, L2) and let
// everything else fall into the unclassified branch, so a tool whose host
// bridge declares the lowest tier (memory.RiskL0, produced by
// tools.levelString's default arm) was refused with the "unclassified"
// sentence whenever the host left Config.PassThroughUnclassifiedRisk off -
// which is the real composition root (cmd/wisp sets it nowhere). D4 says an
// L0 read-only call runs directly and never interrupts, and riskLabel in this
// same file already asks "is it RiskUnclassified" instead of "is it not
// L1/L2". These criteria are the teeth that sentence never had: zero tests in
// the repository pinned the refusal text before this file.
//
// Three shapes, deliberately paired:
//
//   - AC#2 positive: pass-through OFF + a DECLARED L0 call must actually run.
//     Red before the fix, green after.
//   - AC#3 negative: pass-through OFF + a genuinely unclassified call (the
//     directory has no entry for that name, so byName yields the zero
//     ToolInfo) must STILL be refused with the same verbatim sentence. Green
//     before and after - it is the "do not wash out the gate" guard, not this
//     ticket's new tooth.
//   - AC#4 guards: the declared L1/L2 routing is unchanged word for word, both
//     with and without an approval layer wired.

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/observe"
)

// unclassifiedRefusal is the verbatim fail-closed sentence at loop.go. It is
// copied here on purpose: a fix that rewords it must turn AC#3 red, and AC#3
// is the one that proves the gate still exists.
const unclassifiedRefusal = "风险未分级且直通开关关闭，已拒绝执行"

// l0Directory is the dev provider wired with the DECLARED L0 tier, i.e. the
// shape tools.Bridge.Tools produces for fs.read / fs.list / task.output
// (internal/tools/fs.go:300, fs.go:313, task.go:285 -> levelString -> memory.RiskL0).
func l0Directory() []ToolInfo {
	return []ToolInfo{
		{
			Name: "echo", Description: "Echo back the given text.", Resident: true,
			RiskLevel: RiskL0,
			Parameters: json.RawMessage(
				`{"type":"object","properties":{"text":{"type":"string"}},"required":["text"]}`),
		},
		{
			Name: "sleep", Description: "Sleep for the given number of milliseconds.", Resident: true,
			RiskLevel: RiskL0,
			Parameters: json.RawMessage(
				`{"type":"object","properties":{"ms":{"type":"integer"}},"required":["ms"]}`),
		},
	}
}

func openStore179(t *testing.T, dir string) *memory.Store {
	t.Helper()
	store, err := memory.Open(dir,
		memory.WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))))
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

// AC#2 - the positive tooth. With the host-level pass-through switch OFF (the
// real composition root's shape), one DECLARED L0 call must be executed, not
// refused: the provider ran it, its produced body came back to the model, and
// the journal row books an allow plus a success.
func TestDeclaredL0RunsWhenPassThroughIsOff179(t *testing.T) {
	dir := sealableTempDir124(t)
	store := openStore179(t, dir)

	ep := NewEchoProvider(l0Directory()...)
	h := newHarness(t, "tool-then-text", withJournal(store), withTools(ep),
		withConfig(func(c *Config) {
			c.ArtifactsDir = filepath.Join(dir, "artifacts")
			c.PassThroughUnclassifiedRisk = false
		}))
	res := h.run("现在天气怎么样")

	rows, err := store.ListToolCallsByTask(context.Background(), res.TaskID)
	if err != nil {
		t.Fatalf("tool_call lookup: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("tool_call rows = %d, want 1", len(rows))
	}
	row := rows[0]

	if n := h.echo().CallCount(); n != 1 {
		t.Errorf("executions = %d, want 1: a declared L0 call must reach the provider", n)
	}
	if row.Decision != DecisionAllow {
		t.Errorf("decision = %q, want %q", row.Decision, DecisionAllow)
	}
	if row.Outcome != OutcomeSuccess {
		t.Errorf("outcome = %q, want %s", row.Outcome, OutcomeSuccess)
	}
	if row.ErrorClass != "" {
		t.Errorf("error_class = %q, want none", row.ErrorClass)
	}
	if row.RiskLevel != memory.RiskL0 {
		t.Errorf("risk_level = %q, want %s", row.RiskLevel, memory.RiskL0)
	}

	// "It ran" is pinned on the payload, not on the absence of an error: the
	// tool's own output text is what goes back into the next model request.
	calls := h.echo().Calls()
	if len(calls) != 1 {
		t.Fatalf("recorded calls = %d, want 1", len(calls))
	}
	produced := calls[0].Out.Text
	if strings.TrimSpace(produced) == "" {
		t.Fatal("the provider produced empty text; nothing to pin the payload on")
	}
	bodies := h.requestBodies()
	if len(bodies) < 2 {
		t.Fatalf("provider requests = %d, want >= 2 (the tool result must go back)", len(bodies))
	}
	if !strings.Contains(string(bodies[1]), produced) {
		t.Errorf("second request does not carry the tool's produced body %q", produced)
	}
	if strings.Contains(string(bodies[1]), unclassifiedRefusal) {
		t.Errorf("second request carries the refusal sentence %q for a declared L0 call",
			unclassifiedRefusal)
	}
}

// AC#3 - the paired negative guard. Same switch state, but the call is
// genuinely unclassified (no directory entry -> zero ToolInfo -> RiskLevel ==
// ""), and it must still fail closed with the same sentence. Any fix that
// turns the default arm into a pass-through silences this one.
func TestUnclassifiedCallStillRefusedWhenPassThroughIsOff179(t *testing.T) {
	dir := sealableTempDir124(t)
	store := openStore179(t, dir)

	// Default directory: echo/sleep are offered but declare NO risk level,
	// which is exactly RiskUnclassified (tools.go:37, memory never sees an L
	// tier for them). The golden fixture calls "echo".
	ep := NewEchoProvider(ToolInfo{
		Name: "sleep", Description: "Sleep.", Resident: true, RiskLevel: RiskL0,
		Parameters: json.RawMessage(`{"type":"object"}`),
	})
	h := newHarness(t, "tool-then-text", withJournal(store), withTools(ep),
		withConfig(func(c *Config) {
			c.ArtifactsDir = filepath.Join(dir, "artifacts")
			c.PassThroughUnclassifiedRisk = false
		}))
	res := h.run("现在天气怎么样")

	rows, err := store.ListToolCallsByTask(context.Background(), res.TaskID)
	if err != nil {
		t.Fatalf("tool_call lookup: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("tool_call rows = %d, want 1", len(rows))
	}
	row := rows[0]

	if n := h.echo().CallCount(); n != 0 {
		t.Errorf("executions = %d, want 0: an unclassified call must never run", n)
	}
	if row.Decision != DecisionReject {
		t.Errorf("decision = %q, want %q", row.Decision, DecisionReject)
	}
	if row.Outcome != OutcomeError {
		t.Errorf("outcome = %q, want %s", row.Outcome, OutcomeError)
	}
	if row.ErrorClass != string(observe.ClassPermissionDenied) {
		t.Errorf("error_class = %q, want %s", row.ErrorClass, observe.ClassPermissionDenied)
	}

	bodies := h.requestBodies()
	if len(bodies) < 2 {
		t.Fatalf("provider requests = %d, want >= 2", len(bodies))
	}
	if !strings.Contains(string(bodies[1]), unclassifiedRefusal) {
		t.Errorf("second request lost the verbatim fail-closed sentence %q", unclassifiedRefusal)
	}
}

// AC#4 - the declared L1/L2 route is untouched. Without an approval layer the
// old refusal stands word for word; with one wired the loop passes the call
// but books NO decision, because the gate owns that column (ruling A13).
func TestDeclaredL1L2RoutingUnchanged179(t *testing.T) {
	ctx := context.Background()
	dir := sealableTempDir124(t)
	store := openStore179(t, dir)

	for _, tier := range []string{RiskL1, RiskL2} {
		tier := tier
		t.Run(tier+"/no-approval-layer", func(t *testing.T) {
			l := &Loop{} // zero Options: AdmitTask == nil, pass-through off
			j := newTaskJournal(store, "task-179-nil-"+tier, "task-179-nil-"+tier)
			rowID := j.startCall(ctx, "call-1", "fs.write", json.RawMessage(`{}`), riskColumn(tier))

			ok, unbooked, why := l.decideRisk(ctx, j, rowID, ToolInfo{Name: "fs.write", RiskLevel: tier})
			if ok || unbooked {
				t.Errorf("ok/unbooked = %v/%v, want false/false", ok, unbooked)
			}
			want := "该操作属于 " + tier + " 级，审批通道尚未接入（ticket 21），已拒绝执行"
			if why != want {
				t.Errorf("why = %q, want %q", why, want)
			}
			rows, err := store.ListToolCallsByTask(ctx, j.taskID)
			if err != nil {
				t.Fatalf("tool_call lookup: %v", err)
			}
			if len(rows) != 1 || rows[0].Decision != DecisionReject {
				t.Errorf("booked decision = %+v, want one row with %q", rows, DecisionReject)
			}
		})

		t.Run(tier+"/with-approval-layer", func(t *testing.T) {
			l := &Loop{}
			l.opt.AdmitTask = func(string) func() { return nil }
			j := newTaskJournal(store, "task-179-admit-"+tier, "task-179-admit-"+tier)
			rowID := j.startCall(ctx, "call-1", "fs.write", json.RawMessage(`{}`), riskColumn(tier))

			ok, unbooked, why := l.decideRisk(ctx, j, rowID, ToolInfo{Name: "fs.write", RiskLevel: tier})
			if !ok || !unbooked || why != "" {
				t.Errorf("ok/unbooked/why = %v/%v/%q, want true/true/empty", ok, unbooked, why)
			}
			rows, err := store.ListToolCallsByTask(ctx, j.taskID)
			if err != nil {
				t.Fatalf("tool_call lookup: %v", err)
			}
			if len(rows) != 1 || rows[0].Decision != "" {
				t.Errorf("decision = %+v, want the row left pending so the gate owns it", rows)
			}
		})
	}

	// A declared L0 must not be swept into the gate-owned routing either: it
	// books its own allow and reports "booked something".
	l := &Loop{}
	l.opt.AdmitTask = func(string) func() { return nil }
	j := newTaskJournal(store, "task-179-l0-direct", "task-179-l0-direct")
	rowID := j.startCall(ctx, "call-1", "fs.read", json.RawMessage(`{}`), riskColumn(RiskL0))
	ok, unbooked, why := l.decideRisk(ctx, j, rowID, ToolInfo{Name: "fs.read", RiskLevel: RiskL0})
	if !ok || unbooked || why != "" {
		t.Errorf("declared L0 with a gate wired: ok/unbooked/why = %v/%v/%q, want true/false/empty",
			ok, unbooked, why)
	}
}
