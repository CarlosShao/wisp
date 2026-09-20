package tools_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/llm/golden"
	_ "github.com/CarlosShao/wisp/internal/llm/openaichat" // protocol registration
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/tools"
)

// THE LOOP IS NOW WIRED TO THE GATE (ticket 12, ruling A13).
//
// internal/agent/loop.go used to refuse every declared L1/L2 call before
// Execute, because there was nothing behind it to ask - the stopguard was the
// only approval the loop had. That guard made fs.write unreachable (it died
// one line short of the bridge) and made the approval layer's 19 tests a
// description of code nothing called.
//
// The two cases below are the whole point of the swap, and they are a pair on
// purpose:
//
//	TestLoopPassesDeclaredL1WriteThroughTheGate      gate wired  -> window opens,
//	                                                                 file lands
//	TestLoopStillRefusesL1WhenNoGateIsRegistered     gate absent -> refuses
//
// If the early-reject had simply been deleted, the second case would fail: a
// declared write would execute in a host that forgot its approval layer. That
// is the failure mode the dispatch called worse than shipping nothing, so it
// is pinned here rather than promised in a comment.
//
// The provider is a real HTTP server replaying recorded SSE through the real
// openai-chat adapter (SPEC-05 §9: no function mocks at the C5 seam). The
// fixture is generated per test only because the target path lives in
// t.TempDir(); the bytes are the golden format, unmodified.

// composeLoop wires loop + bridge + gate + one real fs.write target.
type loopFixture struct {
	loop   *agent.Loop
	bridge *tools.Bridge
	gate   *approval.Gate
	ui     *uiSpy
	mem    *memory.Store
	target string
}

func composeLoop(t *testing.T, dir string, admit bool) *loopFixture {
	t.Helper()
	target := filepath.Join(dir, "loop-written.txt")
	sse := writeThenTextSSE(t, target)

	res, err := golden.Parse(sse)
	if err != nil {
		t.Fatalf("parse sse: %v", err)
	}
	rep := golden.NewReplayer(res)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(body))
		rep.ServeHTTP(w, r)
	}))
	t.Cleanup(srv.Close)
	prov, err := llm.NewProvider(llm.EndpointOptions{
		Provider: "mock", Model: "probe-model", Protocol: "openai-chat",
		BaseURL: srv.URL, ContextWindow: 128000,
		HTTPClient: &http.Client{Transport: &http.Transport{Proxy: nil}},
	})
	if err != nil {
		t.Fatalf("provider: %v", err)
	}

	mem, err := memory.Open(filepath.Join(dir, "db"))
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	t.Cleanup(func() { _ = mem.Close() })

	ui := &uiSpy{}
	gate := approval.New(approval.Options{UI: ui, Channels: approval.NewChannels()})
	deps := tools.FSDeps{Paths: tools.NewPathCanonicalizer(
		[]string{mustCanon(t, dir)}, nil)}
	reg := tools.NewRegistry()
	for _, e := range tools.BuiltinFSEntries(deps) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	b := tools.New(tools.Options{
		Registry:   reg,
		Paths:      deps.Paths,
		Gate:       gate,
		Cancel:     gate.ToolsCancelBus(),
		Journal:    mem,
		Provenance: nil,
		Logf:       func(string, ...any) {},
	})
	opt := agent.Options{
		Provider: prov,
		Tools:    b,
		Sink:     &agent.RecordSink{},
		Registry: observe.NewRegistry(),
		// The bridge owns the tool_call rows in this composition (the loop's
		// Journal is nil), so exactly one row per call exists and it carries
		// the ASSESSED level plus the gate's decision.
		Journal: nil,
		Config: agent.Config{
			Model: "probe-model", ContextWindow: 128000,
			ArtifactsDir:   filepath.Join(dir, "artifacts"),
			PerToolTimeout: 10 * time.Second,
			// The fs tools classify themselves, so nothing here is
			// unclassified: the pass-through flag stays off, which keeps the
			// fail-closed default honest.
			PassThroughUnclassifiedRisk: false,
		},
	}
	if admit {
		// D47: the loop registers its own task id with the gate, because the
		// id is minted inside run().
		opt.AdmitTask = gate.AdmitTextTask
	} else {
		// The journal follows the DECISION, not the code path: a host that
		// registers no approval layer keeps the pre-execution refusals in the
		// loop, so the loop must own the rows. Exactly one writer per call
		// either way - see the row assertions in both cases.
		opt.Journal = mem
	}
	loop, err := agent.New(opt)
	if err != nil {
		t.Fatalf("agent.New: %v", err)
	}
	return &loopFixture{loop: loop, bridge: b, gate: gate, ui: ui, mem: mem, target: target}
}

// rows reads back what the composition booked.
func (f *loopFixture) rows(t *testing.T, taskID string) []memory.ToolCall {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := f.mem.ListToolCallsByTask(ctx, taskID)
	if err != nil {
		t.Fatalf("ListToolCallsByTask: %v", err)
	}
	return rows
}

// TestLoopPassesDeclaredL1WriteThroughTheGate is the headline of the swap: a
// model that asks for fs.write reaches the L1 block window through the real
// loop, and the write happens because the window ran out unvetoed - not
// because anything approved it silently.
func TestLoopPassesDeclaredL1WriteThroughTheGate(t *testing.T) {
	dir := t.TempDir()
	f := composeLoop(t, dir, true)

	start := time.Now()
	res := f.loop.Run(t.Context(), "把这句话写进文件")
	if res.Status != agent.StatusCompleted {
		t.Fatalf("status = %s (%s): the gate-owned route must complete, not fail\n%s",
			res.Status, res.Message, dumpToolLog(res))
	}
	if res.ToolCalls != 1 {
		t.Fatalf("ToolCalls = %d, want the one fs.write the model asked for:\n%s",
			res.ToolCalls, dumpToolLog(res))
	}
	for _, l := range res.ToolLog {
		if l.Rejected {
			t.Fatalf("the loop still rejected a declared L1 call it should have routed: %+v", l)
		}
	}
	cards := f.ui.cards()
	if len(cards) != 1 {
		t.Fatalf("the gate displayed %d cards, want exactly the one L1 window", len(cards))
	}
	if cards[0].Level != "L1" || cards[0].Tool != "fs.write" {
		t.Errorf("card = %s %s, want L1 fs.write", cards[0].Level, cards[0].Tool)
	}
	if !f.ui.started() {
		t.Error("the gate never reported the handoff to execution")
	}
	body, err := os.ReadFile(f.target)
	if err != nil || string(body) != "written by the loop" {
		t.Fatalf("the write did not land through the window: %v %q", err, body)
	}
	if time.Since(start) < 2*time.Second {
		t.Errorf("the task returned in %v: an unvetoed L1 window must BLOCK", time.Since(start))
	}

	rows := f.rows(t, res.TaskID)
	if len(rows) != 1 {
		t.Fatalf("tool_call rows = %d, want exactly one per call: %+v", len(rows), rows)
	}
	r := rows[0]
	if r.Tool != "fs.write" || r.RiskLevel != memory.RiskL1 || r.Decision != agent.DecisionAllow {
		t.Errorf("row = %+v, want fs.write/L1/allow (the gate's own verdict)", r)
	}
	if r.Outcome != agent.OutcomeSuccess {
		t.Errorf("outcome = %q, want success", r.Outcome)
	}
	if r.CorrelationID == "" || r.CorrelationID != res.TaskID {
		t.Errorf("correlation_id = %q, want the task id %q (C18)", r.CorrelationID, res.TaskID)
	}
	// The loop must not have booked a second row of its own: one call, one
	// authoritative record.
	if r.Seq != 1 {
		t.Errorf("seq = %d, want 1 (a second writer would have bumped it)", r.Seq)
	}
}

// TestLoopStillRefusesL1WhenNoGateIsRegistered is the half that keeps the
// early-reject's deletion safe: a host that assembles a loop without the D47
// registration gets the old refusal back, and the file never appears.
func TestLoopStillRefusesL1WhenNoGateIsRegistered(t *testing.T) {
	dir := t.TempDir()
	f := composeLoop(t, dir, false)

	res := f.loop.Run(t.Context(), "把这句话写进文件")
	if len(f.ui.cards()) != 0 {
		t.Fatalf("a loop with no admitted task must not reach a window: %+v", f.ui.cards())
	}
	if len(res.ToolLog) != 1 || !res.ToolLog[0].Rejected {
		t.Fatalf("the call must be refused before execution:\n%s", dumpToolLog(res))
	}
	if res.ToolLog[0].ErrorClass != string(observe.ClassPermissionDenied) {
		t.Errorf("error_class = %q, want permission_denied", res.ToolLog[0].ErrorClass)
	}
	if !strings.Contains(res.ToolLog[0].Text, "审批通道尚未接入") {
		t.Errorf("the refusal must name the missing approval channel, got %q", res.ToolLog[0].Text)
	}
	if _, err := os.Stat(f.target); !os.IsNotExist(err) {
		t.Fatalf("an ungated host wrote to disk anyway: %v", err)
	}
	if rows := f.rows(t, res.TaskID); len(rows) != 1 ||
		rows[0].Decision != agent.DecisionReject || rows[0].RiskLevel != memory.RiskL1 {
		t.Errorf("the refusal must still be booked, as L1/reject: %+v", rows)
	}
}

func dumpToolLog(res agent.Result) string {
	var b strings.Builder
	for _, l := range res.ToolLog {
		fmt.Fprintf(&b, "  call=%s tool=%s outcome=%s class=%s rejected=%v text=%q\n",
			l.CallID, l.Name, l.Outcome, l.ErrorClass, l.Rejected, l.Text)
	}
	return b.String()
}

// writeThenTextSSE renders the two turns in the golden SSE format: an fs.write
// tool call, then a text answer.
func writeThenTextSSE(t *testing.T, target string) []byte {
	t.Helper()
	args, err := json.Marshal(map[string]string{
		"path": filepath.ToSlash(target), "content": "written by the loop"})
	if err != nil {
		t.Fatal(err)
	}
	argStr, err := json.Marshal(string(args))
	if err != nil {
		t.Fatal(err)
	}
	return []byte(fmt.Sprintf(`# wisp golden sse v1
# @scenario ticket 12 loop+gate: response 1 = one L1 fs.write call, response 2 = final text
# @response 200
data: {"choices":[{"index":0,"delta":{"role":"assistant","content":""},"finish_reason":null}]}

data: {"choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_w1","type":"function","function":{"name":"fs.write","arguments":%s}}]},"finish_reason":null}]}

data: {"choices":[{"index":0,"delta":{},"finish_reason":"tool_calls"}]}

data: {"choices":[],"usage":{"prompt_tokens":30,"completion_tokens":8,"total_tokens":38}}

data: [DONE]
# @response 200
data: {"choices":[{"index":0,"delta":{"content":"已写入。"},"finish_reason":null}]}

data: {"choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}

data: {"choices":[],"usage":{"prompt_tokens":50,"completion_tokens":4,"total_tokens":54}}

data: [DONE]
`, argStr))
}
