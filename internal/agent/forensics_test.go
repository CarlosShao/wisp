package agent

import (
	"context"
	"io"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/observe"
)

// MAJOR-1 (adversarial review of ticket 10): cancellation is one of the four
// named terminal states, and the contract requires the cancelled task's
// task_log row plus its tool_call rows to carry error_class/decision. Because
// loop.go rebinds the task ctx to the root's cancelable ctx (the fix that let
// the D11(1) control layer stop a running task), the terminal writes used to be
// handed an ALREADY-CANCELLED ctx, memory's writer refused them at BeginTx, and
// the rows were stranded: state="running", ended_at NULL, and zero tool_call
// rows for the calls that were in flight. This test pins the three components
// the review demanded.

func TestCancelledTaskPersistsTerminalRows(t *testing.T) {
	dir := t.TempDir()
	store, err := memory.Open(dir,
		memory.WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))))
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	defer store.Close()

	tools := newBlockingProvider()
	h := newHarness(t, "six-tools", withJournal(store), withTools(tools),
		withConfig(func(c *Config) {
			// Cancellation, never the C22 timeout, must be what ends this task.
			c.PerToolTimeout = 10 * time.Second
			c.ArtifactsDir = filepath.Join(dir, "artifacts")
		}))

	task := h.loop.RunAsync(context.Background(), "六个任务并行")
	if !h.waitForFirstRequest(3 * time.Second) {
		tools.release()
		t.Fatal("provider never received a request")
	}
	// Wait until four calls are genuinely in flight: their rows are open, and
	// the cancellation now lands mid-execution.
	if !tools.waitForInflight(MaxToolConcurrency, 3*time.Second) {
		tools.release()
		task.Wait()
		t.Fatalf("only %d/%d tool calls started; the cancelled path is untested",
			tools.MaxConcurrent(), MaxToolConcurrency)
	}
	task.Cancel()
	tools.release()

	res := task.Wait()
	if res.Status != StatusCancelled {
		t.Fatalf("status = %s (%s), want cancelled", res.Status, res.Message)
	}

	// Read back on a fresh, LIVE ctx: the point of the defect is that the rows
	// were written (or not) under the dead one.
	ctx := context.Background()
	tl, err := store.TaskLogByID(ctx, res.TaskID)
	if err != nil {
		t.Fatalf("task_log lookup: %v", err)
	}
	if tl.State == "running" {
		t.Errorf("task_log.state is still %q: the terminal write was abandoned by "+
			"the cancelled ctx (MAJOR-1 regression)", tl.State)
	}
	if tl.State != "cancelled" {
		t.Errorf("task_log.state = %q, want cancelled", tl.State)
	}
	if tl.EndedAt == nil {
		t.Error("task_log.ended_at is NULL: the row never reached a terminal state")
	}

	rows, err := store.ListToolCallsByTask(ctx, res.TaskID)
	if err != nil {
		t.Fatalf("tool_call lookup: %v", err)
	}
	if len(rows) != 6 {
		t.Fatalf("tool_call rows = %d, want 6: the in-flight calls of a cancelled "+
			"task must be booked, not lost (res.ToolLog=%+v)", len(rows), res.ToolLog)
	}
	for _, r := range rows {
		if r.Outcome != OutcomeCancelled {
			t.Errorf("row %s/%s outcome = %q, want %s", r.Tool, r.CorrelationID,
				r.Outcome, OutcomeCancelled)
		}
		if r.EndedAt == nil {
			t.Errorf("row %s ended_at is NULL (still pending)", r.Tool)
		}
		if r.ErrorClass != string(observe.ClassCancelled) {
			t.Errorf("row %s error_class = %q, want %s", r.Tool, r.ErrorClass,
				observe.ClassCancelled)
		}
		// The gate column is part of the contract's row requirement and is
		// written before execution, so it must be present on this path too.
		if r.Decision != DecisionAllow {
			t.Errorf("row %s decision = %q, want %s", r.Tool, r.Decision, DecisionAllow)
		}
	}
	// The in-memory result surface agrees with the persisted rows.
	if len(res.ToolLog) != 6 {
		t.Errorf("res.ToolLog = %d entries, want 6", len(res.ToolLog))
	}
	for _, l := range res.ToolLog {
		if l.Outcome != OutcomeCancelled {
			t.Errorf("result tool log %s outcome = %q, want %s", l.CallID, l.Outcome,
				OutcomeCancelled)
		}
	}
}
