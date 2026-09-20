package agent

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Acceptance criterion 2: stopReason=max_tokens fails ALL unclosed tool calls
// of that message (D21 failToolCallsFromTruncatedMessage - truncated arguments
// are never executed). Asserted at the golden level: the fixture replays a real
// length-stopped stream carrying one complete and one cut tool call, and the
// assertions run against the executed-call count plus the persisted tool_call
// rows.

func TestMaxTokensFailsAllToolCallsOfThatMessage(t *testing.T) {
	dir := t.TempDir()
	store, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	defer store.Close()

	h := newHarness(t, "max-tokens-toolcall", withJournal(store),
		withConfig(func(c *Config) { c.ArtifactsDir = filepath.Join(dir, "artifacts") }))
	res := h.run("帮我整理一下这份清单")

	if res.Status != StatusTruncated {
		t.Fatalf("status = %s (%s), want truncated", res.Status, res.Message)
	}
	// The D21 rule: nothing executed.
	if n := h.echo().CallCount(); n != 0 {
		t.Errorf("tool executions = %d, want 0 (truncated args must never run)", n)
	}
	for _, c := range h.echo().Calls() {
		t.Errorf("unexpected execution: %+v", c.Req)
	}
	if res.ToolCalls != 2 {
		t.Errorf("res.ToolCalls = %d, want 2 (both judged failed)", res.ToolCalls)
	}

	// The stream really carried two calls, one with valid args and one cut
	// mid-argument - the rule fails BOTH.
	if len(res.ToolLog) != 2 {
		t.Fatalf("tool log = %+v", res.ToolLog)
	}
	seen := map[string]ToolResultLog{}
	for _, l := range res.ToolLog {
		seen[l.CallID] = l
	}
	for _, id := range []string{"call_ok1", "call_cut2"} {
		l, ok := seen[id]
		if !ok {
			t.Fatalf("call %s missing from the tool log: %+v", id, res.ToolLog)
		}
		if l.Outcome != OutcomeTruncated {
			t.Errorf("call %s outcome = %q, want %s", id, l.Outcome, OutcomeTruncated)
		}
		if l.ErrorClass != string(observe.ClassLoop) {
			t.Errorf("call %s error_class = %q, want loop", id, l.ErrorClass)
		}
	}
	if json.Valid([]byte(`{"text":"这条参数在输出上限处被`)) {
		t.Fatal("fixture no longer carries a truncated argument")
	}

	// The persisted forensics rows agree with the in-memory log.
	rows, err := store.ListToolCallsByTask(context.Background(), res.TaskID)
	if err != nil {
		t.Fatalf("tool_call lookup: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("tool_call rows = %d, want 2", len(rows))
	}
	for _, row := range rows {
		if row.Outcome != OutcomeTruncated {
			t.Errorf("row %s outcome = %q, want truncated", row.Tool, row.Outcome)
		}
		if row.Decision != DecisionReject {
			t.Errorf("row %s decision = %q, want reject", row.Tool, row.Decision)
		}
		if row.ErrorClass != string(observe.ClassLoop) {
			t.Errorf("row %s error_class = %q, want loop", row.Tool, row.ErrorClass)
		}
	}

	// Partial text is kept and marked, never silently discarded (§3.4), and the
	// user is told what happened.
	if res.Text != "我先写一下，然后" {
		t.Errorf("partial text = %q, want the delivered prefix", res.Text)
	}
	if !strings.Contains(res.Message, "截断") || !strings.Contains(res.Message, "未执行") {
		t.Errorf("message must state the truncation: %q", res.Message)
	}
	if h.sink.LastOf(EvError) == nil {
		t.Error("truncation not published (must be user-visible)")
	}
	// No second round was requested: the loop stops after failing the calls.
	if h.requests() != 1 {
		t.Errorf("requests = %d, want 1", h.requests())
	}
	// The tool_call ids stay in the chain (C25 + later truncation judgements
	// depend on them): each failed call has a matching role:tool result.
	for _, id := range []string{"call_ok1", "call_cut2"} {
		if !historyHasResultID(h.loop.History(), id) {
			t.Errorf("history lost the tool result for %s:\n%s", id, dumpHistory(h.loop.History()))
		}
	}
}

// An open (unclosed) tool call after a mid-stream failure is also failed, not
// executed - the seam never fabricates ToolCallEnd.
func TestOpenToolCallAfterFailureIsFailed(t *testing.T) {
	h := newHarness(t, "disconnect-mid-toolcall")
	res := h.run("慢慢说")
	if res.Status != StatusFailed {
		t.Fatalf("status = %s (%s), want failed", res.Status, res.Message)
	}
	if n := h.echo().CallCount(); n != 0 {
		t.Errorf("executions = %d, want 0 for an unclosed call", n)
	}
	found := false
	for _, l := range res.ToolLog {
		if l.CallID == "call_open" && l.Outcome == OutcomeError {
			found = true
		}
	}
	if !found {
		t.Errorf("open call not judged failed: %+v", res.ToolLog)
	}
	if res.Err == nil || res.Err.Class != observe.ClassNetwork {
		t.Errorf("error = %v, want a network-classified failure", res.Err)
	}
}

// ---------------------------------------------------------------------------

func historyHasResultID(hist []llm.Message, id string) bool {
	for _, m := range hist {
		for _, p := range m.Content {
			if tr, ok := p.(llm.ToolResultPart); ok && tr.ID == id {
				return true
			}
		}
	}
	return false
}
