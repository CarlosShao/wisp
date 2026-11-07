package agent

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Ticket 165 AC#6 (form A of the census's three shapes): the FIRST standing
// instrument for "the model really does knock on the brake".
//
// The gap this closes (census docs/evidence/s1/165-doorbell-census.md, §问 A
// 断口①): internal/agent/loop.go:637-640 is the sink - every tool call the
// provider stream assembles is read there as c.Name / c.Args and booked as one
// tool_call row (the INSERT is internal/memory/dao_toolcall.go:24). The chain
// existed and the sink existed, but no case in any package had ever sent a
// plan.present call through those two lines. Ticket 165's five accepted ACs all
// measure clauses, wiring and the roster; none of them measures the hop. This
// case measures the hop: replay plan-present.sse, then read the real store back.
//
// What this case deliberately does NOT do - reach the real host bridge.
// internal/tools/bridge.go:13 already imports internal/agent and package agent
// has no agent_test external test package, so importing internal/tools from
// here is an import cycle. The directory below therefore MIRRORS the frozen row
// rather than reaching through the bridge, and the tool body is never run (the
// assertion on CallCount pins that). Nothing is mocked in place of the real
// thing to claim the feature works: the one and only judgement is what the
// store read back says about the two sink lines. The bridge half is cmd/wisp
// territory and belongs to another leg.

// The brake's own name, spelled rather than imported: docs/PLAN.md:2567 is the
// frozen D34 row (plan.present, L2, S3) and internal/tools/plan_present.go:72
// is the same string on the tool side.
const planPresentName165r3 = "plan.present"

// The call id, summary and steps are the bytes of response 1 of
// internal/agent/testdata/golden/plan-present.sse, restated here on purpose: a
// change on either side then has to be made twice.
const (
	planPresentCallID165r3  = "call_pp1"
	planPresentFixture165r3 = "plan-present"
	planPresentSummary165r3 = "先把项目收拾成能跑的样子，这一轮只出清单不动手"
)

var planPresentSteps165r3 = []string{
	"列出根目录里散落的构建产物",
	"把三处过期脚本改名并同步 README",
}

type planPresentArgs165r3 struct {
	Summary string   `json:"summary"`
	Steps   []string `json:"steps"`
}

// planPresentDirectory165r3 mirrors the frozen row: Resident true, Declared L2
// (internal/tools/plan_present.go:123-132, PlanPresentDecl; :129 is the
// Resident true). The description is the first clause of the frozen one
// (internal/tools/plan_present.go:77-81). Only the tier and the residency
// matter to the loop here: they decide which arm of decideRisk the call takes
// and which risk_level lands in the row.
func planPresentDirectory165r3() []ToolInfo {
	return []ToolInfo{{
		Name:        planPresentName165r3,
		Description: "先交一份「我要干什么」的清单并停下等点头",
		Resident:    true,
		RiskLevel:   RiskL2,
		Parameters: json.RawMessage(`{"type":"object","properties":{` +
			`"summary":{"type":"string"},` +
			`"steps":{"type":"array","minItems":1,"items":{"type":"string"}}` +
			`},"required":["steps"],"additionalProperties":false}`),
	}}
}

// sinkCount165r3 IS THE READING of the sink at internal/agent/loop.go:637-640:
// how many tool_call rows that loop booked for one tool name. This is the line
// the vacuity check of ticket 165 AC#6 replaces with a constant non-zero value;
// rowOfTool165r3 re-reads the same rows independently, so faking this count
// does not fake the case.
func sinkCount165r3(rows []memory.ToolCall, tool string) int {
	n := 0
	for _, r := range rows {
		if r.Tool == tool {
			n++
		}
	}
	return n
}

// rowOfTool165r3 returns the first row booked for the brake, or fails the test
// on the spot: "no row" is precisely the absence this case exists to catch.
func rowOfTool165r3(t *testing.T, rows []memory.ToolCall, tool string) memory.ToolCall {
	t.Helper()
	for _, r := range rows {
		if r.Tool == tool {
			return r
		}
	}
	t.Fatalf("no tool_call row was booked for %q (rows read back = %d): the call never reached internal/agent/loop.go:637-640",
		tool, len(rows))
	return memory.ToolCall{}
}

// TestPlanPresentCallBooksTheToolCallRow165r3 replays recorded SSE in which the
// model asks for a plan before touching anything, through the genuine
// openai-chat adapter (C5), the genuine Loop and the genuine memory.Store.
//
// The tier matters: with Declared L2 and no Options.AdmitTask wired, this
// harness is an ungated host, so decideRisk refuses the call
// (internal/agent/loop.go:820-823) even though the pass-through switch is on -
// an unapproved write must not run. The refusal does not erase the row: the
// sink books it BEFORE the switch, which is exactly why this table is the
// forensic answer to "did the doorbell ring".
func TestPlanPresentCallBooksTheToolCallRow165r3(t *testing.T) {
	dir := sealableTempDir124(t)
	store, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	defer store.Close()

	h := newHarness(t, planPresentFixture165r3,
		withTools(NewEchoProvider(planPresentDirectory165r3()...)),
		withJournal(store),
		withConfig(func(c *Config) { c.ArtifactsDir = filepath.Join(dir, "artifacts") }))
	res := h.run("帮我把这个项目收拾一下")

	// REDUNDANT NON-ZERO GATE, not THE JUDGMENT: 165-v2 measured that feeding
	// sinkCount165r3 a constant keeps this green, because the tooth of this case
	// is the independent row walk at line 145 below (rowOfTool165r3, :95-105).
	rows, err := store.ListToolCallsByTask(context.Background(), res.TaskID)
	if err != nil {
		t.Fatalf("ListToolCallsByTask: %v", err)
	}
	n := sinkCount165r3(rows, planPresentName165r3)
	if n == 0 {
		t.Fatalf("doorbell sink reading = 0 rows for %q over %d tool_call rows: the replay never reached internal/agent/loop.go:637-640",
			planPresentName165r3, len(rows))
	}

	// One call in the fixture, so one row for the brake and one row total.
	row := rowOfTool165r3(t, rows, planPresentName165r3)
	if n != 1 {
		t.Errorf("sink reading for %q = %d, want 1 (the fixture carries exactly one call)", planPresentName165r3, n)
	}
	if len(rows) != 1 {
		t.Errorf("tool_call rows = %d, want 1: %+v", len(rows), rows)
	}

	// The row is the loop's own, not a fixture: seq, the declared tier that
	// byName[c.Name] resolved, and the args bytes c.Args carried.
	if row.Seq != 1 {
		t.Errorf("row.Seq = %d, want 1", row.Seq)
	}
	if row.RiskLevel != memory.RiskL2 {
		t.Errorf("row.risk_level = %q, want %q (the frozen tier must survive the trip to the row)",
			row.RiskLevel, memory.RiskL2)
	}
	if row.CorrelationID == "" {
		t.Errorf("row.correlation_id is empty (C18): %+v", row)
	}
	if !json.Valid([]byte(row.ArgsJSON)) {
		t.Errorf("row.args_json is not valid JSON: %q", row.ArgsJSON)
	}
	var args planPresentArgs165r3
	if err := json.Unmarshal([]byte(row.ArgsJSON), &args); err != nil {
		t.Fatalf("row.args_json does not parse as the plan.present arguments: %v (%q)", err, row.ArgsJSON)
	}
	if args.Summary != planPresentSummary165r3 {
		t.Errorf("args.summary = %q, want %q", args.Summary, planPresentSummary165r3)
	}
	if len(args.Steps) != len(planPresentSteps165r3) {
		t.Fatalf("args.steps = %v, want %v", args.Steps, planPresentSteps165r3)
	}
	for i, want := range planPresentSteps165r3 {
		if args.Steps[i] != want {
			t.Errorf("args.steps[%d] = %q, want %q", i, args.Steps[i], want)
		}
	}

	// The refusal is booked on the row too: this host wires no approval channel.
	if row.Decision != DecisionReject {
		t.Errorf("row.decision = %q, want %q (an ungated host must not run an L2 write)",
			row.Decision, DecisionReject)
	}
	if row.Outcome != OutcomeError {
		t.Errorf("row.outcome = %q, want %q", row.Outcome, OutcomeError)
	}
	if row.ErrorClass != string(observe.ClassPermissionDenied) {
		t.Errorf("row.error_class = %q, want %q", row.ErrorClass, string(observe.ClassPermissionDenied))
	}

	// The loop read the name and the call id off the wire.
	if res.ToolCalls != 1 || len(res.ToolLog) != 1 {
		t.Fatalf("result tool calls = %d, tool log = %+v, want one call", res.ToolCalls, res.ToolLog)
	}
	if got := res.ToolLog[0]; got.Name != planPresentName165r3 || got.CallID != planPresentCallID165r3 || !got.Rejected {
		t.Errorf("tool log entry = %+v, want a rejected %s/%s", got, planPresentName165r3, planPresentCallID165r3)
	}
	if res.Status != StatusCompleted {
		t.Errorf("status = %s (%s), want completed", res.Status, res.Message)
	}

	// No tool body ran: this case instruments the hop, it does not stage a
	// substitute for the real plan.present implementation.
	if got := h.echo().CallCount(); got != 0 {
		t.Errorf("executed tool calls = %d, want 0 (an ungated L2 call must never execute)", got)
	}

	// The other half of "the model was offered this door and knocked": the same
	// run put the brake on the wire in round 1, and the refusal answer fed
	// round 2. Byte assertions on the real request bodies, the shape already
	// used by TestGoldenToolResultFeedsNextTurn.
	bodies := h.requestBodies()
	if len(bodies) != 2 {
		t.Fatalf("provider requests = %d, want 2 (the call must open a second round)", len(bodies))
	}
	mustContain(t, "round 1 body", string(bodies[0]), `"name":"`+planPresentName165r3+`"`)
	mustContain(t, "round 2 body", string(bodies[1]), `"tool_call_id":"`+planPresentCallID165r3+`"`)
}
