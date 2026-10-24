package memory

// Ticket 373: deleting a task_log parent takes its tool_call children with it,
// and the children go FIRST in the same write transaction. Three product paths
// are covered here, one case each, plus the two readings the ticket's acceptance
// is phrased in: the ruler reads dirty before and clean after (section 5 of the
// census this ticket rests on), and a failure of the parent statement rolls the
// child statement back with it.
//
// What makes these cases bite on the pre-373 bytes: the child rows exist before
// the parent delete (seeded through the production DAO or RAW SQL for the legacy
// shapes), so "row count of tool_call" and PRAGMA foreign_key_check both keep
// the old reading there. The old readings are 3/2/1 orphans for the three paths
// (census 373-a1 section 2, P1/P2/P4); every case below asserts 0.

import (
	"context"
	"database/sql"
	"errors"
	"sort"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// t373SeedTwoTasks builds the census fixture: task-A with two children and
// task-B with one, through the production DAO.
func t373SeedTwoTasks(t *testing.T, s *Store) {
	t.Helper()
	ctx := context.Background()
	for _, id := range []string{"task-A", "task-B"} {
		if err := s.StartTaskLog(ctx, TaskLog{ID: id, State: "running", QueryText: "q-" + id}); err != nil {
			t.Fatalf("seed task_log %s: %v", id, err)
		}
	}
	for _, tc := range []ToolCall{
		{TaskID: "task-A", Seq: 1, Tool: "fs.read", ArgsJSON: `{"p":"a"}`, RiskLevel: RiskL0, CorrelationID: "c-a1"},
		{TaskID: "task-A", Seq: 2, Tool: "fs.write", ArgsJSON: `{"p":"b"}`, RiskLevel: RiskL1, CorrelationID: "c-a2"},
		{TaskID: "task-B", Seq: 1, Tool: "fs.read", ArgsJSON: `{"p":"c"}`, RiskLevel: RiskL0, CorrelationID: "c-b1"},
	} {
		if _, err := s.InsertToolCall(ctx, tc); err != nil {
			t.Fatalf("seed tool_call %s#%d: %v", tc.TaskID, tc.Seq, err)
		}
	}
}

// t373ChildTaskIDs lists tool_call.task_id sorted, so a case can name exactly
// whose forensic rows it removed - "the count went down" is not enough
// (ticket 373: every cascade must say what it deleted).
func t373ChildTaskIDs(t *testing.T, s *Store) []string {
	t.Helper()
	rows, err := s.ListToolCalls(context.Background())
	if err != nil {
		t.Fatalf("list tool_call: %v", err)
	}
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.TaskID)
	}
	sort.Strings(out)
	return out
}

// t373Removed answers the multiset difference before-after, both sorted.
func t373Removed(before, after []string) []string {
	var out []string
	j := 0
	for _, id := range before {
		if j < len(after) && after[j] == id {
			j++
			continue
		}
		out = append(out, id)
	}
	return out
}

// t373FkCount is the acceptance ruler: PRAGMA foreign_key_check row count on a
// business-shaped connection (the independent second ruler of ticket 366's
// helper, not the Store's own reader).
func t373FkCount(t *testing.T, s *Store) int {
	t.Helper()
	return len(t366FkRows(t, s.dbPath))
}

// TestTicket373PurgeTaskLogsTakesTheChildren is P1 (census: task_log 2 ->
// tool_call 3 -> foreign_key_check 3). After the fix the purge removes both
// tables, and it names which children it removed.
func TestTicket373PurgeTaskLogsTakesTheChildren(t *testing.T) {
	s := openTestStore(t)
	t373SeedTwoTasks(t, s)
	if n := t373FkCount(t, s); n != 0 {
		t.Fatalf("foreign_key_check before the purge = %d, want 0: the fixture is not the shape", n)
	}
	childrenBefore := t373ChildTaskIDs(t, s)

	n, err := s.PurgePrivacy(context.Background(), PrivacyTaskLog)
	if err != nil {
		t.Fatalf("PurgePrivacy(PrivacyTaskLog): %v", err)
	}
	if n != 2 {
		t.Errorf("purge reported %d task rows, want 2 (the count the UI shows is task rows)", n)
	}
	if got := countRows(t, s, "task_log"); got != 0 {
		t.Errorf("task_log rows = %d after the purge, want 0", got)
	}
	if got := countRows(t, s, "tool_call"); got != 0 {
		t.Errorf("tool_call rows = %d after the purge, want 0: the children carry tool names "+
			"and arguments, so a cleared task log with them left behind is only half cleared", got)
	}
	if removed := t373Removed(childrenBefore, t373ChildTaskIDs(t, s)); len(removed) != 3 ||
		removed[0] != "task-A" || removed[1] != "task-A" || removed[2] != "task-B" {
		t.Errorf("children removed = %v, want [task-A task-A task-B]: every child of a deleted "+
			"parent goes, and nothing else is touched", removed)
	}
	if n := t373FkCount(t, s); n != 0 {
		t.Errorf("foreign_key_check after the purge = %d, want 0", n)
	}
}

// TestTicket373DeleteTaskLogTakesOnlyThatTasksChildren is P2 (census 2
// orphans). The other task's child must survive untouched.
func TestTicket373DeleteTaskLogTakesOnlyThatTasksChildren(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	t373SeedTwoTasks(t, s)
	childrenBefore := t373ChildTaskIDs(t, s)

	if err := s.DeletePrivacyItem(ctx, PrivacyTaskLog, "task-A"); err != nil {
		t.Fatalf("DeletePrivacyItem(PrivacyTaskLog, task-A): %v", err)
	}
	if _, err := s.TaskLogByID(ctx, "task-A"); !errors.Is(err, ErrNotFound) {
		t.Errorf("task-A parent still readable: %v", err)
	}
	if _, err := s.TaskLogByID(ctx, "task-B"); err != nil {
		t.Errorf("task-B parent must survive: %v", err)
	}
	after := t373ChildTaskIDs(t, s)
	if len(after) != 1 || after[0] != "task-B" {
		t.Errorf("tool_call.task_id after deleting task-A = %v, want [task-B] only", after)
	}
	if removed := t373Removed(childrenBefore, after); len(removed) != 2 ||
		removed[0] != "task-A" || removed[1] != "task-A" {
		t.Errorf("children removed = %v, want [task-A task-A]: deleting one task takes its own "+
			"children and not another task's", removed)
	}
	if n := t373FkCount(t, s); n != 0 {
		t.Errorf("foreign_key_check after deleting task-A = %d, want 0", n)
	}
	// A missing parent is still an error, and its children (if any) are not
	// silently half-deleted: the statement order is child-then-parent inside one
	// transaction, so the not-found rollback restores the child delete too.
	orphanSeed := ToolCall{
		TaskID: "task-never-booked", Seq: 1, Tool: "fs.read",
		ArgsJSON: `{}`, RiskLevel: RiskL0, CorrelationID: "c-orphan",
	}
	if _, err := s.InsertToolCall(ctx, orphanSeed); err != nil {
		t.Fatalf("seed the orphan child: %v", err)
	}
	if err := s.DeleteTaskLog(ctx, "task-never-booked"); !errors.Is(err, ErrNotFound) {
		t.Errorf("DeleteTaskLog on a missing parent = %v, want ErrNotFound", err)
	}
	if got := t373ChildTaskIDs(t, s); len(got) != 2 || got[0] != "task-B" || got[1] != "task-never-booked" {
		t.Errorf("tool_call.task_id after the not-found delete = %v, want [task-B task-never-booked]: "+
			"the child delete must roll back with the parent's not-found error", got)
	}
}

// TestTicket373RetentionTakesChildrenOutsideTheirWindow is P4, the second
// mechanism the ticket is for: parent started 31 days ago, its children 5 and 2
// days ago. The old pass deleted the parent and left both children (1 more
// orphan every pass, for ever); the fix deletes the children first.
func TestTicket373RetentionTakesChildrenOutsideTheirWindow(t *testing.T) {
	s := openTestStore(t)
	now := observe.NowWallUTC()
	day := 24 * time.Hour

	seedTaskLogAt(t, s, "task-old", now.Add(-31*day))
	seedToolCallFull(t, s, "task-old", 1, now.Add(-5*day), now.Add(-5*day), true)
	seedToolCallFull(t, s, "task-old", 2, now.Add(-2*day), now.Add(-1*day), true)
	// Control: a live parent keeps its live child, and a live child that is past
	// its own window still ages out on its own (the statement after the cascade).
	seedTaskLogAt(t, s, "task-live", now.Add(-time.Hour))
	seedToolCallFull(t, s, "task-live", 1, now.Add(-time.Hour), now, true)

	childrenBefore := t373ChildTaskIDs(t, s)
	if n := t373FkCount(t, s); n != 0 {
		t.Fatalf("foreign_key_check before the pass = %d, want 0", n)
	}

	res := s.RunRetentionOnce(context.Background(),
		RetentionConfig{Now: func() time.Time { return now }})
	if res.TaskLogDeleted != 1 {
		t.Errorf("task_log deleted = %d, want 1 (task-old's parent)", res.TaskLogDeleted)
	}
	if res.ToolCallDeleted != 2 {
		t.Errorf("tool_call deleted = %d, want 2 (both of task-old's children, including the "+
			"one still inside its own window)", res.ToolCallDeleted)
	}
	after := t373ChildTaskIDs(t, s)
	if len(after) != 1 || after[0] != "task-live" {
		t.Errorf("tool_call.task_id after the pass = %v, want [task-live]: the live parent's "+
			"child must not be touched by the cascade", after)
	}
	if removed := t373Removed(childrenBefore, after); len(removed) != 2 ||
		removed[0] != "task-old" || removed[1] != "task-old" {
		t.Errorf("children removed = %v, want [task-old task-old]", removed)
	}
	if n := t373FkCount(t, s); n != 0 {
		t.Errorf("foreign_key_check after the pass = %d, want 0: a parent taken while its "+
			"children stay is exactly the orphan this ticket exists for", n)
	}
	// Second pass must find nothing: the orphan was not merely postponed.
	second := s.RunRetentionOnce(context.Background(),
		RetentionConfig{Now: func() time.Time { return now }})
	if second.TaskLogDeleted != 0 || second.ToolCallDeleted != 0 {
		t.Errorf("pass 2 deleted task_log=%d tool_call=%d, want 0 and 0",
			second.TaskLogDeleted, second.ToolCallDeleted)
	}
	if n := t373FkCount(t, s); n != 0 {
		t.Errorf("foreign_key_check after pass 2 = %d, want 0", n)
	}
}

// TestTicket373PositiveControlTheRulerReadsDirtyThenClean is the control the
// census used before trusting the 3/2/1 readings: the ruler can say 3, and the
// same rows read 0 once the parent/child pair is deleted together. The dirty
// state is the shape an older build left on disk (parent removed by hand, as
// the pre-373 statement removed it), so the purge here has no parent to count
// and only the inherited orphan rows to take.
func TestTicket373PositiveControlTheRulerReadsDirtyThenClean(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	t373SeedTwoTasks(t, s)

	if err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `DELETE FROM task_log`)
		return err
	}); err != nil {
		t.Fatalf("seed the pre-373 state: %v", err)
	}
	if n := t373FkCount(t, s); n != 3 {
		t.Fatalf("foreign_key_check on the dirty fixture = %d, want 3: the ruler cannot read the "+
			"state this case is about, so a later 0 would prove nothing", n)
	}

	n, err := s.PurgeTaskLogs(ctx)
	if err != nil || n != 0 {
		t.Fatalf("purge = %d, %v; want 0 task rows (none were left) and nil", n, err)
	}
	if got := countRows(t, s, "tool_call"); got != 0 {
		t.Errorf("tool_call rows = %d after the purge, want 0: an inherited orphan is a child of "+
			"a parent that is already gone, and the one-click clear must take it too", got)
	}
	if got := countRows(t, s, "task_log"); got != 0 {
		t.Errorf("task_log rows = %d after the purge, want 0", got)
	}
	if n := t373FkCount(t, s); n != 0 {
		t.Errorf("foreign_key_check after the purge = %d, want 0", n)
	}
}

// t373BlockParentDeletes plants a trigger that refuses every task_log delete, so
// the parent statement fails AFTER the child statement has already run inside the
// same transaction. That is what makes the rollback evidence about the
// transaction rather than about statement order alone.
func t373BlockParentDeletes(t *testing.T, s *Store) {
	t.Helper()
	err := s.write(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			`CREATE TRIGGER t373_block_task_log_delete BEFORE DELETE ON task_log
			 BEGIN SELECT RAISE(ABORT, 't373: parent delete blocked'); END`)
		return err
	})
	if err != nil {
		t.Fatalf("plant the blocking trigger: %v", err)
	}
}

// t373DropParentDeleteBlock removes the trigger again, so the same fixture can
// take the unblocked path: without that second half a "the rows are still here"
// reading would also be produced by a cascade that never ran.
func t373DropParentDeleteBlock(t *testing.T, s *Store) {
	t.Helper()
	err := s.write(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `DROP TRIGGER t373_block_task_log_delete`)
		return err
	})
	if err != nil {
		t.Fatalf("drop the blocking trigger: %v", err)
	}
}

// TestTicket373ChildDeleteRollsBackWithTheParent proves the two statements are
// one transaction on both single-row paths: with the parent delete refused, the
// child delete that already ran is undone, and the database is left exactly as
// it was (foreign_key_check still 0, three tool_call rows). The unblocking half
// at the end is the control: the same fixture, same call, without the trigger,
// does remove the children - so "still three" above was a rollback and not a
// cascade that never ran.
func TestTicket373ChildDeleteRollsBackWithTheParent(t *testing.T) {
	t.Run("DeleteTaskLog", func(t *testing.T) {
		s := openTestStore(t)
		ctx := context.Background()
		t373SeedTwoTasks(t, s)
		t373BlockParentDeletes(t, s)

		if err := s.DeleteTaskLog(ctx, "task-A"); err == nil {
			t.Fatal("DeleteTaskLog returned nil although the parent delete was refused")
		}
		if got := countRows(t, s, "tool_call"); got != 3 {
			t.Errorf("tool_call rows after the refused delete = %d, want 3: the child statement "+
				"ran first, so anything but a rollback means the pair is not one transaction", got)
		}
		if got := countRows(t, s, "task_log"); got != 2 {
			t.Errorf("task_log rows after the refused delete = %d, want 2", got)
		}
		if n := t373FkCount(t, s); n != 0 {
			t.Errorf("foreign_key_check = %d, want 0", n)
		}

		t373DropParentDeleteBlock(t, s)
		if err := s.DeleteTaskLog(ctx, "task-A"); err != nil {
			t.Fatalf("DeleteTaskLog after unblocking: %v", err)
		}
		if got := countRows(t, s, "tool_call"); got != 1 {
			t.Errorf("tool_call rows after the unblocked delete = %d, want 1: the trigger, not "+
				"the code, is what kept them above", got)
		}
	})

	t.Run("RunRetentionOnce", func(t *testing.T) {
		s := openTestStore(t)
		now := observe.NowWallUTC()
		day := 24 * time.Hour
		seedTaskLogAt(t, s, "task-old", now.Add(-31*day))
		seedToolCallFull(t, s, "task-old", 1, now.Add(-5*day), now.Add(-5*day), true)
		t373BlockParentDeletes(t, s)

		s.RunRetentionOnce(context.Background(),
			RetentionConfig{Now: func() time.Time { return now }})
		// The pass reports the failure on the log and returns; the database is
		// the authority here: the cascade's child delete must have been undone
		// with the parent statement that failed after it.
		if got := countRows(t, s, "tool_call"); got != 1 {
			t.Errorf("tool_call rows after the refused pass = %d, want 1: the cascade ran before "+
				"the parent delete, so a kept row is the proof both are in one transaction", got)
		}
		if got := countRows(t, s, "task_log"); got != 1 {
			t.Errorf("task_log rows after the refused pass = %d, want 1", got)
		}
		if n := t373FkCount(t, s); n != 0 {
			t.Errorf("foreign_key_check = %d, want 0", n)
		}

		t373DropParentDeleteBlock(t, s)
		s.RunRetentionOnce(context.Background(),
			RetentionConfig{Now: func() time.Time { return now }})
		if got := countRows(t, s, "tool_call"); got != 0 {
			t.Errorf("tool_call rows after the unblocked pass = %d, want 0", got)
		}
		if got := countRows(t, s, "task_log"); got != 0 {
			t.Errorf("task_log rows after the unblocked pass = %d, want 0", got)
		}
	})
}
