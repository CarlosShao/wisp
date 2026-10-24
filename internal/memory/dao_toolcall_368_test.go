package memory

// Ticket 368: InsertToolCall fills the aging timestamps a row proves have
// already happened (stampToolCallInsert). Retention ages tool_call by
// COALESCE(ended_at, started_at, decided_at), so before this fallback a row
// inserted with all three NULL could never age: the child outlived its task_log
// parent for ever and every later pass skipped it again. The cases below are
// the two writers that reach the DAO (tools.Bridge.book, whose row already
// carries a decision and an outcome, and agent.taskJournal.startCall, whose row
// is pending), the two rulers that must keep working (normal aging, and three
// caller-supplied timestamps that must not be rewritten), plus the column order
// the COALESCE expression depends on.
//
// What pins them to today: cases 1, 2, 3 and 5 assert timestamps that the
// unfixed DAO binds as SQL NULL, so they fail on the pre-368 bytes and pass
// only with the fallback. Case 4 is the opposite direction - it fails if the
// fallback ever starts rewriting or inventing a value the caller supplied - and
// case 6 is the ruler control that has to pass on both sides.

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

const t368TTL = 30 * 24 * time.Hour

// t368TS renders one nullable unix second so a failure message names the value
// instead of a pointer.
func t368TS(p *int64) string {
	if p == nil {
		return "<NULL>"
	}
	return strconv.FormatInt(*p, 10)
}

// t368Booked is the row tools.Bridge.book builds: verdict and outcome set, no
// timestamps at all (bridge.go, memory.ToolCall literal).
func t368Booked(taskID string, seq int64) ToolCall {
	return ToolCall{
		TaskID:        taskID,
		Seq:           seq,
		Tool:          "fs.read",
		ArgsJSON:      `{}`,
		RiskLevel:     RiskL0,
		Decision:      "allow",
		Outcome:       "success",
		CorrelationID: "corr-" + taskID,
	}
}

// t368Pending is the row agent.taskJournal.startCall builds: neither the
// decision nor the outcome exists yet, and it sets no timestamps either.
func t368Pending(taskID string, seq int64) ToolCall {
	return ToolCall{
		TaskID:        taskID,
		Seq:           seq,
		Tool:          "fs.read",
		ArgsJSON:      `{}`,
		RiskLevel:     RiskL0,
		CorrelationID: "corr-" + taskID,
	}
}

// t368Retention runs one retention pass against a clock 31 days ahead, so a row
// inserted now is strictly older than the 30-day window. The parent is seeded
// at the same pinned "now", which is what makes it age in the same pass.
func t368Retention(t *testing.T, s *Store, now time.Time) RetentionResult {
	t.Helper()
	pinned := now.Add(31 * 24 * time.Hour)
	return s.RunRetentionOnce(context.Background(),
		RetentionConfig{Now: func() time.Time { return pinned }})
}

// Case 1 (tooth): a booked row carries the two instants it proves. Both land on
// the bookkeeping instant, started_at stays NULL because that writer never saw
// the call begin.
func TestTicket368BookedRowCarriesTheInstantsItProves(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	id, err := s.InsertToolCall(ctx, t368Booked("task-368-booked", 1))
	if err != nil {
		t.Fatalf("InsertToolCall: %v", err)
	}
	row, err := s.ToolCallByID(ctx, id)
	if err != nil {
		t.Fatalf("ToolCallByID: %v", err)
	}
	if row.DecidedAt == nil {
		t.Error("decided_at is NULL on a row inserted with decision=allow: " +
			"the booking instant is the decision instant, so the fallback must fill it")
	}
	if row.EndedAt == nil {
		t.Error("ended_at is NULL on a row inserted with outcome=success: " +
			"the call had already finished when it was booked, so the fallback must fill it")
	}
	if row.DecidedAt != nil && row.EndedAt != nil && *row.DecidedAt != *row.EndedAt {
		t.Errorf("decided_at=%d ended_at=%d, want one instant: a booked row records "+
			"both events at the moment it is written", *row.DecidedAt, *row.EndedAt)
	}
	if row.StartedAt != nil {
		t.Errorf("started_at=%d, want NULL: the booking writer did not observe the call begin, "+
			"so inventing a start here would be the column lying about what it measured", *row.StartedAt)
	}
	if row.Decision != "allow" || row.Outcome != "success" {
		t.Errorf("enum columns rewritten: decision=%q outcome=%q", row.Decision, row.Outcome)
	}
}

// Case 2 (tooth, the ticket's acceptance): the booked row ages out with its
// parent, and the orphan count returns to zero. Two passes, because the census
// measured that the second pass used to leave the same orphan behind.
func TestTicket368BookedRowIsNotImmortal(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	now := observe.NowWallUTC()

	seedTaskLogAt(t, s, "task-368-booked", now)
	if _, err := s.InsertToolCall(ctx, t368Booked("task-368-booked", 1)); err != nil {
		t.Fatalf("InsertToolCall: %v", err)
	}
	if n := len(t366FkRows(t, s.dbPath)); n != 0 {
		t.Fatalf("foreign_key_check before retention = %d rows, want 0: the fixture is not the shape", n)
	}

	first := t368Retention(t, s, now)
	if first.TaskLogDeleted != 1 || first.ToolCallDeleted != 1 {
		t.Errorf("pass 1 deleted task_log=%d tool_call=%d, want 1 and 1: the child must go with its parent",
			first.TaskLogDeleted, first.ToolCallDeleted)
	}
	if got := countRows(t, s, "tool_call"); got != 0 {
		t.Errorf("tool_call rows after pass 1 = %d, want 0 (the booked row aged out)", got)
	}
	if rows := t366FkRows(t, s.dbPath); len(rows) != 0 {
		t.Errorf("foreign_key_check after pass 1 = %v, want none: an aging row must not leave an orphan", rows)
	}

	second := t368Retention(t, s, now)
	if second.TaskLogDeleted != 0 || second.ToolCallDeleted != 0 {
		t.Errorf("pass 2 deleted task_log=%d tool_call=%d, want 0 and 0 (nothing left to age)",
			second.TaskLogDeleted, second.ToolCallDeleted)
	}
	if rows := t366FkRows(t, s.dbPath); len(rows) != 0 {
		t.Errorf("foreign_key_check after pass 2 = %v, want none", rows)
	}
}

// Case 3 (tooth): a pending row - startCall's shape, inserted before the call
// runs - ages from its insert instant instead of living for ever.
func TestTicket368PendingRowAgesFromItsInsert(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	now := observe.NowWallUTC()

	seedTaskLogAt(t, s, "task-368-pending", now)
	id, err := s.InsertToolCall(ctx, t368Pending("task-368-pending", 1))
	if err != nil {
		t.Fatalf("InsertToolCall: %v", err)
	}
	row, err := s.ToolCallByID(ctx, id)
	if err != nil {
		t.Fatalf("ToolCallByID: %v", err)
	}
	if row.StartedAt == nil {
		t.Fatal("started_at is NULL on a pending row: with decision and outcome both pending, " +
			"the insert instant is the only start the row has, and without it the " +
			"COALESCE in retention sees NULL and the row can never age")
	}
	if row.DecidedAt != nil || row.EndedAt != nil {
		t.Errorf("pending row got decided_at=%s ended_at=%s, want both NULL: "+
			"nothing has been decided or finished yet", t368TS(row.DecidedAt), t368TS(row.EndedAt))
	}

	res := t368Retention(t, s, now)
	if res.TaskLogDeleted != 1 || res.ToolCallDeleted != 1 {
		t.Errorf("deleted task_log=%d tool_call=%d, want 1 and 1: a pending child must age with its parent",
			res.TaskLogDeleted, res.ToolCallDeleted)
	}
	if rows := t366FkRows(t, s.dbPath); len(rows) != 0 {
		t.Errorf("foreign_key_check = %v, want none: the child was left behind as an orphan", rows)
	}
}

// Case 4 (reverse tooth): caller-supplied timestamps survive untouched. This is
// the case that fails if the fallback ever turns into a rewrite - on a booked
// row and on a pending row alike.
func TestTicket368ExplicitTimestampsAreNotOverwritten(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	now := observe.NowWallUTC()

	decided := now.Add(-400 * 24 * time.Hour).Unix()
	started := now.Add(-401 * 24 * time.Hour).Unix()
	ended := now.Add(-402 * 24 * time.Hour).Unix()
	booked := t368Booked("task-368-explicit", 1)
	booked.DecidedAt, booked.StartedAt, booked.EndedAt = &decided, &started, &ended
	id, err := s.InsertToolCall(ctx, booked)
	if err != nil {
		t.Fatalf("InsertToolCall(booked): %v", err)
	}
	row, err := s.ToolCallByID(ctx, id)
	if err != nil {
		t.Fatalf("ToolCallByID: %v", err)
	}
	if row.DecidedAt == nil || *row.DecidedAt != decided {
		t.Errorf("decided_at = %s, want the supplied %d (the fallback only fills NULLs)", t368TS(row.DecidedAt), decided)
	}
	if row.StartedAt == nil || *row.StartedAt != started {
		t.Errorf("started_at = %s, want the supplied %d", t368TS(row.StartedAt), started)
	}
	if row.EndedAt == nil || *row.EndedAt != ended {
		t.Errorf("ended_at = %s, want the supplied %d", t368TS(row.EndedAt), ended)
	}

	pending := t368Pending("task-368-explicit-pending", 2)
	pending.StartedAt = &started
	id, err = s.InsertToolCall(ctx, pending)
	if err != nil {
		t.Fatalf("InsertToolCall(pending): %v", err)
	}
	row, err = s.ToolCallByID(ctx, id)
	if err != nil {
		t.Fatalf("ToolCallByID: %v", err)
	}
	if row.StartedAt == nil || *row.StartedAt != started {
		t.Errorf("pending started_at = %s, want the supplied %d", t368TS(row.StartedAt), started)
	}

	// A row booked after the fact without an outcome: the insert instant would
	// postdate the end it already carries, so the start must stay unknown.
	endedOnly := t368Pending("task-368-explicit-endonly", 3)
	endedOnly.EndedAt = &ended
	id, err = s.InsertToolCall(ctx, endedOnly)
	if err != nil {
		t.Fatalf("InsertToolCall(ended only): %v", err)
	}
	row, err = s.ToolCallByID(ctx, id)
	if err != nil {
		t.Fatalf("ToolCallByID: %v", err)
	}
	if row.EndedAt == nil || *row.EndedAt != ended {
		t.Errorf("ended_at = %s, want the supplied %d", t368TS(row.EndedAt), ended)
	}
	if row.StartedAt != nil {
		t.Errorf("started_at = %s, want NULL: the supplied ended_at is %d and stamping the insert "+
			"instant would make the row start after it finishes", t368TS(row.StartedAt), ended)
	}
}

// Case 5 (tooth): the column ORDER the retention expression relies on still
// holds once the fallback is in: a call starts at its insert, is decided at the
// gate, and ends when the tool returns.
func TestTicket368ColumnOrderSurvivesTheWrites(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	id, err := s.InsertToolCall(ctx, t368Pending("task-368-order", 1))
	if err != nil {
		t.Fatalf("InsertToolCall: %v", err)
	}
	if err := s.DecideToolCall(ctx, id, "allow", nil); err != nil {
		t.Fatalf("DecideToolCall: %v", err)
	}
	if err := s.FinishToolCall(ctx, id, "success", ""); err != nil {
		t.Fatalf("FinishToolCall: %v", err)
	}
	row, err := s.ToolCallByID(ctx, id)
	if err != nil {
		t.Fatalf("ToolCallByID: %v", err)
	}
	if row.StartedAt == nil || row.DecidedAt == nil || row.EndedAt == nil {
		t.Fatalf("timestamps after the full lifecycle = {started %v, decided %v, ended %v}, want all three set: "+
			"the COALESCE order ended_at > started_at > decided_at is only meaningful when every "+
			"event that happened has its own instant", t368TS(row.StartedAt), t368TS(row.DecidedAt), t368TS(row.EndedAt))
	}
	if *row.EndedAt < *row.DecidedAt || *row.DecidedAt < *row.StartedAt {
		t.Errorf("order started=%d decided=%d ended=%d, want started <= decided <= ended",
			*row.StartedAt, *row.DecidedAt, *row.EndedAt)
	}

	res := t368Retention(t, s, observe.NowWallUTC())
	if res.ToolCallDeleted != 1 {
		t.Errorf("completed call deleted = %d, want 1 (it ages by ended_at like any finished row)", res.ToolCallDeleted)
	}
}

// Case 6 (ruler control): rows that already carried a timestamp keep aging
// normally. This one passes before and after the fallback - it is here so a red
// above cannot be blamed on a stuck retention ruler.
func TestTicket368NormalAgingStillWorks(t *testing.T) {
	s := openTestStore(t)
	now := observe.NowWallUTC()
	old := now.Add(-31 * 24 * time.Hour)

	seedTaskLogAt(t, s, "task-368-old", old)
	// A row written by a pre-368 build: only ended_at is set, by raw SQL, so the
	// fallback cannot be what makes it age.
	seedToolCallFull(t, s, "task-368-old", 1, old, old, true)

	res := s.RunRetentionOnce(context.Background(),
		RetentionConfig{Now: func() time.Time { return now }})
	if res.TaskLogDeleted != 1 || res.ToolCallDeleted != 1 {
		t.Errorf("deleted task_log=%d tool_call=%d, want 1 and 1: an old child must go with its parent",
			res.TaskLogDeleted, res.ToolCallDeleted)
	}
	if rows := t366FkRows(t, s.dbPath); len(rows) != 0 {
		t.Errorf("foreign_key_check = %v, want none", rows)
	}
}

// A guard for the constant above: the package TTL and this file's window must
// stay the same number, or case 2/3 would silently measure a different window.
func TestTicket368WindowMatchesTheContract(t *testing.T) {
	if t368TTL != TaskLogTTL {
		t.Fatalf("t368TTL = %s, TaskLogTTL = %s: the cases pin 30 days", t368TTL, TaskLogTTL)
	}
}
