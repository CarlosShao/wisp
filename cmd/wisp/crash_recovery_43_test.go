package main

// Ticket 43, crash recovery: the marker round trip, the interrupted sweep and
// the no-rerun proof, against the real SQLite store.
//
// These are cross-platform on purpose: the sweep runs from assembleRuntime
// (cmd/wisp/run.go), which the CLI walks on every OS. The kill -9 itself is a
// live-tier reading (lifecycle_live_43_windows_test.go); this file proves the
// semantics the sweep is made of: a live marker means the previous process died
// without its own shutdown, still-'running' rows become interrupted exactly
// once, the row's own numbers survive the transition, and nothing re-executes
// anything.

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/CarlosShao/wisp/internal/memory"
)

func Test43RunMarkerRoundTrip(t *testing.T) {
	dir := t.TempDir()
	if hasRunMarker(dir) {
		t.Fatal("a fresh data root already has a run marker")
	}
	if err := writeRunMarker(dir); err != nil {
		t.Fatalf("writeRunMarker: %v", err)
	}
	if !hasRunMarker(dir) {
		t.Fatal("the marker is gone right after it was written")
	}
	m, err := readRunMarker(dir)
	if err != nil {
		t.Fatalf("readRunMarker: %v", err)
	}
	if m.PID == 0 {
		t.Fatalf("marker pid = %d, want this process's pid", m.PID)
	}
	if err := clearRunMarker(dir); err != nil {
		t.Fatalf("clearRunMarker: %v", err)
	}
	if hasRunMarker(dir) {
		t.Fatal("the marker survived its own clear: a clean exit would read as a crash")
	}
	// A double clear is the boot-with-no-marker shape; it must stay silent.
	if err := clearRunMarker(dir); err != nil {
		t.Fatalf("second clearRunMarker: %v", err)
	}
}

// seedRunningTask43 plants one task_log row in the state a kill -9 leaves
// behind: opened (state 'running', ended_at NULL) and never finished.
func seedRunningTask43(t *testing.T, dir, id string) memory.TaskLog {
	t.Helper()
	st, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("memory.Open(%s): %v", dir, err)
	}
	t.Cleanup(func() { _ = st.Close() })
	tl := memory.TaskLog{
		ID: id, State: string(memory.TaskStateRunning), QueryText: "上次没跑完的任务",
		CostTokensIn: 11, CostTokensOut: 7, CostAmountMicro: 30, SummaryText: strPtr43("死前已写的摘要"),
	}
	if err := st.StartTaskLog(context.Background(), tl); err != nil {
		t.Fatalf("StartTaskLog(%s): %v", id, err)
	}
	// One finished row alongside: the sweep must not touch terminal rows.
	if err := st.StartTaskLog(context.Background(), memory.TaskLog{ID: id + "-done", State: string(memory.TaskStateRunning), QueryText: "已收口任务"}); err != nil {
		t.Fatalf("StartTaskLog(%s-done): %v", id, err)
	}
	if err := st.FinishTaskLog(context.Background(), id+"-done", string(memory.TaskStateDone), "", 1, 1, 1, ""); err != nil {
		t.Fatalf("FinishTaskLog(%s-done): %v", id, err)
	}
	return tl
}

func strPtr43(s string) *string { return &s }

func taskByID43(t *testing.T, dir, id string) memory.TaskLog {
	t.Helper()
	st, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("memory.Open(%s): %v", dir, err)
	}
	t.Cleanup(func() { _ = st.Close() })
	tl, err := st.TaskLogByID(context.Background(), id)
	if err != nil {
		t.Fatalf("TaskLogByID(%s): %v", id, err)
	}
	return tl
}

func Test43CrashSweepMarksRunningTasksInterruptedOnce(t *testing.T) {
	dir := t.TempDir()
	seedRunningTask43(t, dir, "task-43-crash")
	if err := writeRunMarker(dir); err != nil {
		t.Fatalf("writeRunMarker: %v", err)
	}

	st, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	n, err := recoverCrashedTasks(context.Background(), st, dir)
	if err != nil {
		t.Fatalf("recoverCrashedTasks: %v", err)
	}
	if n != 1 {
		t.Fatalf("sweep marked %d tasks, want exactly the one still-running row", n)
	}

	tl := taskByID43(t, dir, "task-43-crash")
	if tl.State != string(memory.TaskStateInterrupted) {
		t.Fatalf("task state = %q, want %q (the listing's 中断 word)", tl.State, memory.TaskStateInterrupted)
	}
	if tl.EndedAt == nil {
		t.Fatal("task ended_at is NULL: the row never reached a terminal state")
	}
	if tl.SummaryText == nil || *tl.SummaryText != "死前已写的摘要" {
		t.Fatalf("summary = %v, want the dead process's own summary preserved", tl.SummaryText)
	}
	if tl.CostTokensIn != 11 || tl.CostTokensOut != 7 || tl.CostAmountMicro != 30 {
		t.Fatalf("cost columns moved: %+v, want 11/7/30 unchanged", tl)
	}

	// The marker stays: the NEXT boot is the one that acts. A sweep under the
	// same marker is idempotent (the rows are terminal now), and the count
	// says so instead of re-marking.
	n2, err := recoverCrashedTasks(context.Background(), st, dir)
	if err != nil {
		t.Fatalf("second sweep: %v", err)
	}
	if n2 != 0 {
		t.Fatalf("second sweep marked %d tasks, want 0: terminal rows are not re-marked", n2)
	}

	// NO RERUN: the sweep's only writes are the terminal transitions above.
	// The row count before and after is the structural proof.
	rows, err := st.ListTaskLogs(context.Background())
	if err != nil {
		t.Fatalf("ListTaskLogs: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("task rows = %d, want 2 (the interrupted one and the finished one): a rerun would have opened a new row", len(rows))
	}
}

func Test43CrashSweepWithoutMarkerIsANoOp(t *testing.T) {
	dir := t.TempDir()
	seedRunningTask43(t, dir, "task-43-alive")
	st, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	n, err := recoverCrashedTasks(context.Background(), st, dir)
	if err != nil {
		t.Fatalf("recoverCrashedTasks: %v", err)
	}
	if n != 0 {
		t.Fatalf("sweep without a marker marked %d tasks, want 0", n)
	}
	tl := taskByID43(t, dir, "task-43-alive")
	if tl.State != string(memory.TaskStateRunning) {
		t.Fatalf("task state = %q, want running: a still-alive process's row must not be touched", tl.State)
	}
}

func Test43UpdatePendingHandoffBooksOnlyWhenPresent(t *testing.T) {
	dir := t.TempDir()
	if markUpdatePendingHandoff(dir) {
		t.Fatal("a data root with no update-pending.json booked an apply-handoff event")
	}
	p := filepath.Join(dir, updatePendingFile)
	if err := os.WriteFile(p, []byte("{}"), 0o600); err != nil {
		t.Fatalf("write cue: %v", err)
	}
	if !markUpdatePendingHandoff(dir) {
		t.Fatal("the cue existed and no apply-handoff event was booked")
	}
}
