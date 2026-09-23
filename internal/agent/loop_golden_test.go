package agent

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Acceptance criterion 1: golden-driven loop tests. Every case here replays
// recorded SSE bytes through the real openai-chat adapter (the C5 seam), so
// what is asserted is what the loop sees in production.

func TestGoldenTextReply(t *testing.T) {
	h := newHarness(t, "text-reply")
	res := h.run("提醒我五分钟后关火")

	if res.Status != StatusCompleted {
		t.Fatalf("status = %s (%s), want completed", res.Status, res.Message)
	}
	if res.Text != "好的，已经把提醒设好了。" {
		t.Errorf("text = %q", res.Text)
	}
	if res.Rounds != 1 {
		t.Errorf("rounds = %d, want 1", res.Rounds)
	}
	if h.requests() != 1 {
		t.Errorf("provider requests = %d, want 1", h.requests())
	}
	if res.ToolCalls != 0 {
		t.Errorf("tool calls = %d, want 0", res.ToolCalls)
	}
	// C23 data: usage lands on the result from the Usage event.
	if res.Usage.InputTokens != 140 || res.Usage.OutputTokens != 9 {
		t.Errorf("usage = %+v", res.Usage)
	}
	if res.Usage.CachedTokens != 96 {
		t.Errorf("cached = %d, want 96 (prompt cache hit is visible)", res.Usage.CachedTokens)
	}
	// The streamed text reached the sink as deltas.
	if h.sink.Of(EvTextDelta) == 0 {
		t.Error("no text deltas published")
	}
	if got := h.sink.LastOf(EvDone); got == nil {
		t.Error("no done event published")
	}
}

func TestGoldenSingleToolCall(t *testing.T) {
	h := newHarness(t, "tool-then-text")
	res := h.run("现在天气怎么样")

	if res.Status != StatusCompleted {
		t.Fatalf("status = %s (%s), want completed", res.Status, res.Message)
	}
	if n := h.echo().CallCount(); n != 1 {
		t.Fatalf("tool executions = %d, want 1", n)
	}
	call := h.echo().Calls()[0]
	if call.Req.Name != "echo" || call.Req.CallID != "call_e1" {
		t.Errorf("call = %+v", call.Req)
	}
	if call.Req.TaskID != res.TaskID || call.Req.CorrelationID != res.TaskID {
		t.Errorf("call identity = %+v, want task %s", call.Req, res.TaskID)
	}
	if res.Text != "现在是 22 摄氏度，晴。" {
		t.Errorf("final text = %q", res.Text)
	}
	// The assistant turn and its tool result both stay in the history, keyed by
	// the same call id (C25/D21 depend on the chain being intact).
	if !historyHasCallID(h.loop.History(), "call_e1") {
		t.Errorf("history lost call id: %s", dumpHistory(h.loop.History()))
	}
}

func TestGoldenParallelToolCalls(t *testing.T) {
	h := newHarness(t, "parallel-tools")
	res := h.run("两件事一起办")

	if res.Status != StatusCompleted {
		t.Fatalf("status = %s (%s)", res.Status, res.Message)
	}
	if n := h.echo().CallCount(); n != 2 {
		t.Fatalf("tool executions = %d, want 2", n)
	}
	// The loop books results in wire order regardless of execution order.
	if len(res.ToolLog) != 2 || res.ToolLog[0].CallID != "call_p1" ||
		res.ToolLog[1].CallID != "call_p2" {
		t.Errorf("tool log order = %+v, want call_p1,call_p2", res.ToolLog)
	}
	if res.ToolLog[0].Outcome != OutcomeSuccess || res.ToolLog[1].Outcome != OutcomeSuccess {
		t.Errorf("tool outcomes = %+v", res.ToolLog)
	}
	// Two calls of one message -> two role:tool messages, in call order.
	var results []string
	for _, m := range h.loop.History() {
		if m.Role != llm.RoleTool {
			continue
		}
		tr, ok := m.Content[0].(llm.ToolResultPart)
		if !ok {
			t.Fatalf("tool message carries %T", m.Content[0])
		}
		results = append(results, tr.ID)
	}
	if strings.Join(results, ",") != "call_p1,call_p2" {
		t.Errorf("tool result ids = %v", results)
	}
}

func TestGoldenToolResultFeedsNextTurn(t *testing.T) {
	h := newHarness(t, "tool-then-text")
	h.run("现在天气怎么样")

	bodies := h.requestBodies()
	if len(bodies) != 2 {
		t.Fatalf("requests = %d, want 2 (tool result must open a second round)", len(bodies))
	}
	// Round 2's wire body must carry the assistant tool_calls plus the tool
	// result content, addressed by call id (role:"tool" + tool_call_id).
	second := string(bodies[1])
	mustContain(t, "round 2 body", second, `"tool_call_id":"call_e1"`)
	mustContain(t, "round 2 body", second, "22 摄氏度，晴")
	mustContain(t, "round 2 body", second, `"name":"echo"`)
	// Round 1 must NOT contain any tool result yet.
	if strings.Contains(string(bodies[0]), "tool_call_id") {
		t.Errorf("round 1 body already carries a tool result:\n%s", trunc(string(bodies[0]), 600))
	}
}

func TestGoldenConcurrencyCeiling(t *testing.T) {
	h := newHarness(t, "six-tools")
	res := h.run("六个任务并行")
	if res.Status != StatusCompleted {
		t.Fatalf("status = %s (%s)", res.Status, res.Message)
	}
	if n := h.echo().CallCount(); n != 6 {
		t.Fatalf("executions = %d, want 6", n)
	}
	if got := h.echo().MaxConcurrent(); got > MaxToolConcurrency {
		t.Errorf("max concurrent tools = %d, want <= %d (D38d ceiling)",
			got, MaxToolConcurrency)
	}
	if got := h.echo().MaxConcurrent(); got < 2 {
		t.Logf("note: observed concurrency %d (scheduler did not overlap; ceiling still honored)", got)
	}
	// NOTE: with an instant-returning tool the check above is satisfied by a
	// fully serial pool, so it only proves the ceiling. TestToolExecutionRunsFourAcross
	// below holds its calls open and proves the parallelism.
}

// MINOR-6: D38d's ceiling is only half the claim - the other half is that tool
// calls really do run four-abreast. EchoProvider returns so fast that
// MaxConcurrent() stays 1 no matter how the pool is built, so the ceiling test
// above cannot tell "limit 4" from "strictly serial". Here the provider holds
// every call until the test says so, which makes four overlapping executions a
// precondition rather than an accident: a serial implementation parks at one
// and fails the wait below.
func TestToolExecutionRunsFourAcross(t *testing.T) {
	tools := newBlockingProvider()
	t.Cleanup(tools.release)
	h := newHarness(t, "six-tools", withTools(tools), withConfig(func(c *Config) {
		c.PerToolTimeout = 30 * time.Second // the release ends the hold, not the brake
	}))

	task := h.loop.RunAsync(context.Background(), "六个任务并行")
	if !h.waitForFirstRequest(5 * time.Second) {
		task.Cancel()
		task.Wait()
		t.Fatal("provider never received a request")
	}
	if !tools.waitForInflight(MaxToolConcurrency, 5*time.Second) {
		tools.release()
		task.Wait()
		t.Fatalf("only %d calls ran at once, want %d: the pool is serial, so the "+
			"ceiling assertion in TestGoldenConcurrencyCeiling is vacuous",
			tools.MaxConcurrent(), MaxToolConcurrency)
	}
	// Four of six are parked right now: the ceiling is a gate, not a hint, so
	// not even a fifth call may have entered the provider.
	if got := tools.MaxConcurrent(); got > MaxToolConcurrency {
		t.Errorf("max concurrent = %d, want <= %d (D38d ceiling)", got, MaxToolConcurrency)
	}
	if got := tools.Started(); got > MaxToolConcurrency {
		t.Errorf("calls entered while four were held = %d, want <= %d", got, MaxToolConcurrency)
	}

	tools.release()
	res := task.Wait()
	if res.Status != StatusCompleted {
		t.Fatalf("status = %s (%s), want completed", res.Status, res.Message)
	}
	if got := tools.CallCount(); got != 6 {
		t.Errorf("executions = %d, want 6", got)
	}
	if got := tools.MaxConcurrent(); got != MaxToolConcurrency {
		t.Errorf("observed concurrency = %d, want exactly %d: the ceiling must be "+
			"reached, not merely respected", got, MaxToolConcurrency)
	}
}

// Budget exhaustion is the C22 token brake: Stuck state plus an explicit,
// user-visible message that names the consumption. No silent stop.
func TestGoldenBudgetExhaustionStuck(t *testing.T) {
	var cfgTokenBudget int = 300 // two rounds of 170 tokens
	h := newHarness(t, "budget-loop", withConfig(func(c *Config) {
		c.TokenBudget = cfgTokenBudget
	}))
	res := h.run("分三步把这件事做完")

	if res.Status != StatusStuck {
		t.Fatalf("status = %s, want stuck (message: %s)", res.Status, res.Message)
	}
	if res.Brake != BrakeToken {
		t.Errorf("brake = %s, want %s", res.Brake, BrakeToken)
	}
	if !strings.Contains(res.Message, "token 预算") {
		t.Errorf("message must name the brake: %q", res.Message)
	}
	if !strings.Contains(res.Message, "本次任务已用 340 token") {
		t.Errorf("message must state consumption: %q", res.Message)
	}
	st := h.sink.LastOf(EvStuck)
	if st == nil {
		t.Fatal("no stuck event published (C22 forbids stopping silently)")
	}
	if st.Text != res.Message {
		t.Errorf("published stuck text = %q, result message = %q", st.Text, res.Message)
	}
	// The third round must never have been requested.
	if h.requests() != 2 {
		t.Errorf("requests = %d, want 2", h.requests())
	}
	if res.RootPending != 0 {
		t.Errorf("root pending = %d, want 0", res.RootPending)
	}
}

// Root ctx cancellation must reach stream and tool execution end-to-end
// (D21#5: AbortSignal through loop -> streamFn -> tool.execute).
func TestGoldenCancellationEndToEnd(t *testing.T) {
	h := newHarness(t, "slow-tool", withConfig(func(c *Config) {
		c.PerToolTimeout = 10 * time.Second // cancellation, not the timeout, must end it
	}))
	task := h.loop.RunAsync(context.Background(), "慢一点也行")
	if !h.waitForFirstRequest(3 * time.Second) {
		t.Fatal("provider never received a request")
	}
	task.Cancel()

	res := task.Wait() // the loop must return on its own after cancellation
	if res.Status != StatusCancelled {
		t.Fatalf("status = %s (%s), want cancelled", res.Status, res.Message)
	}
	if res.RootPending != 0 {
		t.Errorf("task root still has %d goroutines pending", res.RootPending)
	}
	if got := task.Pending(); got != 0 {
		t.Errorf("RunAsync root pending = %d, want 0", got)
	}
}

// D38e: task completion means the root's join counter drained, and the
// registry roster is back to zero for this task.
func TestGoroutineBudgetDrains(t *testing.T) {
	reg := observe.NewRegistry()
	h := newHarness(t, "parallel-tools", withRegistry(reg))

	task := h.loop.RunAsync(context.Background(), "两件事一起办")
	res := task.Wait()
	if res.Status != StatusCompleted {
		t.Fatalf("status = %s (%s)", res.Status, res.Message)
	}
	if got := task.Pending(); got != 0 {
		t.Errorf("root pending = %d, want 0 (WaitGroup must drain)", got)
	}
	if got := reg.Count(); got != 0 {
		t.Errorf("registry live goroutines = %d, want 0", got)
	}
	rep := reg.RosterReport()
	if len(rep.Unknown) != 0 {
		t.Errorf("goroutines outside the D38 roster: %v", rep.Unknown)
	}
	if rep.PerTask != 0 {
		t.Errorf("per-task goroutines still live: %+v", reg.Snapshot())
	}
	if rep.ResidentOverBaseline {
		t.Errorf("resident baseline exceeded: %+v", rep)
	}
}

// task_log + tool_call rows (the 04 tables) through the real store.
func TestTaskLogAndToolCallRows(t *testing.T) {
	dir := sealableTempDir124(t)
	store, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	defer store.Close()

	h := newHarness(t, "tool-then-text", withJournal(store),
		withConfig(func(c *Config) { c.ArtifactsDir = filepath.Join(dir, "artifacts") }))
	res := h.run("现在天气怎么样")
	if res.Status != StatusCompleted {
		t.Fatalf("status = %s (%s)", res.Status, res.Message)
	}

	ctx := context.Background()
	tl, err := store.TaskLogByID(ctx, res.TaskID)
	if err != nil {
		t.Fatalf("task_log lookup: %v", err)
	}
	if tl.State != "done" {
		t.Errorf("task_log.state = %q, want done", tl.State)
	}
	if tl.QueryText != "现在天气怎么样" {
		t.Errorf("task_log.query_text = %q", tl.QueryText)
	}
	if tl.CostTokensIn != 330 || tl.CostTokensOut != 24 {
		t.Errorf("task_log cost = %d/%d, want 330/24 (both turns)",
			tl.CostTokensIn, tl.CostTokensOut)
	}
	if tl.SummaryText == nil || !strings.Contains(*tl.SummaryText, "22 摄氏度") {
		t.Errorf("task_log.summary_text = %v", tl.SummaryText)
	}

	rows, err := store.ListToolCallsByTask(ctx, res.TaskID)
	if err != nil {
		t.Fatalf("tool_call lookup: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("tool_call rows = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Tool != "echo" || row.CorrelationID != res.TaskID {
		t.Errorf("row = %+v", row)
	}
	if row.Decision != DecisionAllow {
		t.Errorf("decision = %q, want %q (pre-gate pass-through)", row.Decision, DecisionAllow)
	}
	if row.Outcome != OutcomeSuccess {
		t.Errorf("outcome = %q, want %s", row.Outcome, OutcomeSuccess)
	}
	if row.RiskLevel != memory.RiskL0 {
		t.Errorf("risk_level = %q, want L0 placeholder until ticket 21", row.RiskLevel)
	}
	if !json.Valid([]byte(row.ArgsJSON)) {
		t.Errorf("args_json is not valid JSON: %q", row.ArgsJSON)
	}
}

// A failing provider (golden 5xx-with-no-retry path) must surface as a
// classified, visible failure - never a silent downgrade (SPEC-05 §3.4).
func TestGoldenProviderFailureIsVisible(t *testing.T) {
	h := newHarness(t, "provider-500")
	res := h.run("办件事")
	if res.Status != StatusFailed {
		t.Fatalf("status = %s, want failed (%s)", res.Status, res.Message)
	}
	if res.Err == nil {
		t.Fatal("no classified error")
	}
	if res.Err.Class != observe.ClassProvider && res.Err.Class != observe.ClassNetwork {
		t.Errorf("error class = %s, want provider/network", res.Err.Class)
	}
	if h.sink.LastOf(EvError) == nil {
		t.Error("failure was not published (must be user-visible)")
	}
}

// ---------------------------------------------------------------------------

func historyHasCallID(hist []llm.Message, id string) bool {
	for _, m := range hist {
		for _, p := range m.Content {
			switch c := p.(type) {
			case llm.ToolUsePart:
				if c.ID == id {
					return true
				}
			case llm.ToolResultPart:
				if c.ID == id {
					return true
				}
			}
		}
	}
	return false
}

func dumpHistory(hist []llm.Message) string {
	var b strings.Builder
	for _, m := range hist {
		b.WriteString(string(m.Role) + ": ")
		for _, p := range m.Content {
			switch c := p.(type) {
			case llm.TextPart:
				b.WriteString(trunc(c.Text, 40) + " ")
			case llm.ToolUsePart:
				b.WriteString("<use " + c.ID + "/" + c.Name + "> ")
			case llm.ToolResultPart:
				b.WriteString("<result " + c.ID + "> ")
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}
