package main

// Ticket 378 件二: the reliable acceptance tooth for "a task.spawn child's
// tool_call rows must have a parent task_log row" (the second orphan source
// ticket 373 cannot reach).
//
// WHAT THIS PINS, and why it is not a mock (AGENTS 1.3): the fix lives in
// cmd/wisp/run.go:1573-1585 - when the store exists, the loop options carry
// taskLifecycleJournal (the task-row half) and rt.loopOpt is republished so
// subagent_197.go:320's copy hands the SAME journal to every child loop. The
// observable consequence is a persistence fact: a child that opens its
// task_log row through that journal half makes the bridge-written tool_call
// rows name an existing parent, so PRAGMA foreign_key_check reads 0. This test
// drives taskLifecycleJournal (the exact type run.go installs) against a REAL
// memory.Store and asserts that fact with the same independent ruler ticket 373
// uses (a fresh connection running PRAGMA foreign_key_check). It is a component
// test at the seam run.go occupies, not a mock standing in for it; the wiring
// that hands the child this journal is verified by reading run.go:1573-1585 and
// subagent_197.go:320 (see docs/evidence/s1/378-ac-verdict.md).

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
)

// TestTicket378ChildTaskRowParentsItsToolCalls is P1: with the child's
// task_log row written the way the child loop's journal half writes it, the
// child's tool_call is parented and foreign_key_check reads 0.
func TestTicket378ChildTaskRowParentsItsToolCalls(t *testing.T) {
	dir := t.TempDir()
	st, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("memory.Open(%s): %v", dir, err)
	}
	t.Cleanup(func() { _ = st.Close() })
	ctx := context.Background()

	// Exactly what run.go hands the child loop: the task-row half only.
	j := &taskLifecycleJournal{store: st}
	if err := j.StartTaskLog(ctx, memory.TaskLog{
		ID: "child-378", StartedAt: time.Now().Unix(), State: string(memory.TaskStateRunning), QueryText: "子任务读笔记",
	}); err != nil {
		t.Fatalf("taskLifecycleJournal.StartTaskLog: %v", err)
	}

	// The row is really there: one row per child task, which is the whole point
	// of the payload (run.go comment "ONE row per child task and nothing else").
	rows, err := st.ListTaskLogs(ctx)
	if err != nil {
		t.Fatalf("ListTaskLogs: %v", err)
	}
	found := false
	for _, r := range rows {
		if r.ID == "child-378" {
			found = true
		}
	}
	if !found {
		t.Fatalf("task_log has no row for child-378 after StartTaskLog (rows=%d): the child loop's journal half did not land its row, which is exactly ticket 378's bug", len(rows))
	}

	// The child's tool_call (bridge-written). The journal's own InsertToolCall
	// must stay inert so rows-per-call does not double (378 hard reminder #1).
	if n, err := j.InsertToolCall(ctx, memory.ToolCall{TaskID: "child-378"}); err != nil || n != 0 {
		t.Errorf("taskLifecycleJournal.InsertToolCall = (%d,%v), want (0,nil): the task-row half must not book tool_call rows", n, err)
	}
	if _, err := st.InsertToolCall(ctx, memory.ToolCall{
		TaskID: "child-378", Seq: 1, Tool: "task.spawn", ArgsJSON: "{}", RiskLevel: "L0", CorrelationID: "corr-378",
	}); err != nil {
		t.Fatalf("store.InsertToolCall(child-378): %v", err)
	}

	if n := fkCount378(t, dir); n != 0 {
		t.Errorf("PRAGMA foreign_key_check after the child row + its tool_call = %d, want 0", n)
	}
}

// TestTicket378OrphanToolCallIsWhatTheRulerSees is the "0 means something"
// guard: a tool_call whose task row was never written (the pre-fix state 378
// fixes) IS counted by the ruler. A 0 in P1 would prove nothing if this were
// also 0.
func TestTicket378OrphanToolCallIsWhatTheRulerSees(t *testing.T) {
	dir := t.TempDir()
	st, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("memory.Open(%s): %v", dir, err)
	}
	t.Cleanup(func() { _ = st.Close() })
	ctx := context.Background()

	if _, err := st.InsertToolCall(ctx, memory.ToolCall{
		TaskID: "ghost-378", Seq: 1, Tool: "x", ArgsJSON: "{}", RiskLevel: "L0", CorrelationID: "corr-ghost",
	}); err != nil {
		t.Fatalf("store.InsertToolCall(ghost-378): %v", err)
	}
	if n := fkCount378(t, dir); n != 1 {
		t.Fatalf("PRAGMA foreign_key_check on an orphan tool_call = %d, want 1: the ruler cannot see the pre-fix shape, so P1's 0 would prove nothing", n)
	}
}

// fkCount378 counts PRAGMA foreign_key_check rows on a fresh connection to the
// store's wisp.db - ticket 373's acceptance ruler, replicated here because its
// helpers (t373FkCount/t366FkRows) are private to internal/memory.
func fkCount378(t *testing.T, dir string) int {
	t.Helper()
	db, err := sql.Open(memory.DriverName, filepath.Join(dir, "wisp.db"))
	if err != nil {
		t.Fatalf("sql.Open for the independent foreign_key_check: %v", err)
	}
	defer func() { _ = db.Close() }()
	rows, err := db.QueryContext(context.Background(), "PRAGMA foreign_key_check")
	if err != nil {
		t.Fatalf("PRAGMA foreign_key_check: %v", err)
	}
	defer func() { _ = rows.Close() }()
	n := 0
	for rows.Next() {
		n++
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("drain foreign_key_check: %v", err)
	}
	return n
}
