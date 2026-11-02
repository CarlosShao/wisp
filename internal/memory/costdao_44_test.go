package memory

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Ticket 44 (C23): the read/repair surface over the frozen D35 tables. These
// tests pin the DAO against a real SQLite file; the boundary arithmetic that
// feeds them (day bounds, basis points, the 价格未知 path) is pinned in
// internal/observe/costmeter_44_test.go.

func TestAC2CostDaysBetweenFiltersAndOrders(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	for _, row := range []struct {
		day                                string
		tokensIn, tokensOut, micros, tasks int64
	}{
		{"2026-10-31", 1, 1, 1, 1}, // one day BEFORE the range
		{"2026-11-01", 100, 20, 500, 2},
		{"2026-11-02", 50, 5, 250, 1},
		{"2026-12-01", 7, 7, 7, 1}, // one day AFTER the range
	} {
		if err := s.BumpCostDay(ctx, row.day, row.tokensIn, row.tokensOut, row.micros); err != nil {
			t.Fatal(err)
		}
		for i := int64(1); i < row.tasks; i++ {
			if err := s.BumpCostDay(ctx, row.day, 0, 0, 0); err != nil {
				t.Fatal(err)
			}
		}
	}

	got, err := s.CostDaysBetween(ctx, "2026-11-01", "2026-11-02")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("range returned %d rows (%+v), want exactly the two days inside it", len(got), got)
	}
	if got[0].Day != "2026-11-01" || got[1].Day != "2026-11-02" {
		t.Errorf("rows are not oldest-first: %s then %s", got[0].Day, got[1].Day)
	}
	if got[0].Micros != 500 || got[0].Tasks != 2 {
		t.Errorf("2026-11-01 = %+v, want 500 micros over 2 tasks", got[0])
	}

	// An empty range is a fact, not an error: "no cost today" must render.
	empty, err := s.CostDaysBetween(ctx, "2026-11-15", "2026-11-16")
	if err != nil {
		t.Fatalf("an empty range must not error: %v", err)
	}
	if len(empty) != 0 {
		t.Errorf("empty range returned %d rows", len(empty))
	}

	// The day label is the only ordering key the SQL trusts, so a malformed
	// endpoint is refused rather than compared lexically into nonsense.
	if _, err := s.CostDaysBetween(ctx, "2026-11-1", "2026-11-02"); err == nil {
		t.Error("a one-digit day was accepted")
	}
	if _, err := s.CostDaysBetween(ctx, "2026-11-02", "2026-11-01"); err == nil {
		t.Error("an inverted range was accepted")
	} else if !strings.Contains(err.Error(), "inverted") {
		t.Errorf("inverted range error %q must say inverted", err)
	}
}

func TestAC2SumTasksEndedBetweenCountsOnlyFinishedRows(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	dayStart := time.Date(2026, 11, 2, 0, 0, 0, 0, time.UTC)
	seed := []struct {
		id             string
		ended          *time.Time // nil = still running
		in, out, micro int64
	}{
		{"t-in-first-second", ptr(dayStart), 10, 1, 100},
		{"t-in-last-second", ptr(dayStart.Add(24*time.Hour - time.Second)), 20, 2, 200},
		{"t-yesterday", ptr(dayStart.Add(-time.Second)), 99, 9, 999},
		{"t-tomorrow", ptr(dayStart.Add(24 * time.Hour)), 99, 9, 999},
		{"t-still-running", nil, 55, 5, 555},
	}
	for _, r := range seed {
		if err := s.StartTaskLog(ctx, TaskLog{
			ID: r.id, StartedAt: dayStart.Add(-time.Hour).Unix(),
			State: string(TaskStateRunning), QueryText: "seed",
		}); err != nil {
			t.Fatalf("seed %s: %v", r.id, err)
		}
		if r.ended != nil {
			if err := setEndedAt(ctx, s, r.id, r.ended.Unix(), r.in, r.out, r.micro); err != nil {
				t.Fatalf("seed finish %s: %v", r.id, err)
			}
		}
	}

	// [00:00, 24:00) of 2026-11-02 holds exactly the two in-range rows: the
	// first second of the day is IN, the last second is IN, the next day's
	// first second is OUT, and a running row is OUT no matter its numbers.
	sum, err := s.SumTasksEndedBetween(ctx, dayStart.Unix(), dayStart.Add(24*time.Hour).Unix())
	if err != nil {
		t.Fatal(err)
	}
	if sum.Tasks != 2 || sum.TokensIn != 30 || sum.TokensOut != 3 || sum.Micros != 300 {
		t.Errorf("day sum = %+v, want {2 tasks, 30 in, 3 out, 300 micros}", sum)
	}

	// A range nobody ended in reads zero with no error (COUNT(*) over no rows).
	none, err := s.SumTasksEndedBetween(ctx, dayStart.Add(72*time.Hour).Unix(), dayStart.Add(96*time.Hour).Unix())
	if err != nil {
		t.Fatal(err)
	}
	if none.Tasks != 0 || none.Micros != 0 {
		t.Errorf("empty range read %+v, want the zero row", none)
	}
	// An empty interval is a caller bug: it would silently read as "no cost".
	if _, err := s.SumTasksEndedBetween(ctx, dayStart.Unix(), dayStart.Unix()); err == nil {
		t.Error("an empty range (from==to) was accepted")
	}
}

func TestAC2ReconcileAgainstRealStoreAndRepair(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	now := time.Now()

	m, err := observe.NewCostMeter(s,
		observe.WithLedger(s), observe.WithWriter(s), observe.WithTasks(s),
		observe.WithLocation(time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	day := m.DayKey(now)

	// Two tasks that went through the loop: their task_log rows carry the same
	// numbers the meter booked into cost_daily.
	for _, task := range []struct{ in, out, micros int64 }{{120, 30, 640}, {80, 12, 260}} {
		id := "reconcile-" + strconv.FormatInt(task.in, 10)
		if err := s.StartTaskLog(ctx, TaskLog{
			ID: id, StartedAt: now.Add(-time.Minute).Unix(),
			State: string(TaskStateRunning), QueryText: "q",
		}); err != nil {
			t.Fatal(err)
		}
		if err := s.FinishTaskLog(ctx, id, string(TaskStateDone), "s", task.in, task.out, task.micros, ""); err != nil {
			t.Fatal(err)
		}
		if _, err := m.PostTask(ctx, task.in, task.out, task.micros); err != nil {
			t.Fatal(err)
		}
	}

	rec, err := m.Reconcile(ctx, day)
	if err != nil {
		t.Fatal(err)
	}
	if !rec.Agree {
		t.Fatalf("reconcile = %+v, want agreement", rec)
	}
	if rec.Aggregate.Micros != 900 || rec.Aggregate.Tasks != 2 {
		t.Errorf("aggregate = %+v, want 900 micros over 2 tasks", rec.Aggregate)
	}

	// Now the drift the aggregate cache can actually develop: a row written by
	// someone who never went through the meter (cmd/wisp closes task_log for
	// the shapes that never reached the loop). The instrument must see it.
	stray := "stray-row"
	if err := s.StartTaskLog(ctx, TaskLog{
		ID: stray, StartedAt: now.Add(-time.Minute).Unix(),
		State: string(TaskStateRunning), QueryText: "q",
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.FinishTaskLog(ctx, stray, string(TaskStateDone), "s", 5, 5, 50, ""); err != nil {
		t.Fatal(err)
	}
	rec, err = m.Reconcile(ctx, day)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Agree {
		t.Fatal("a base-only row went unnoticed: the reconciliation cell would be decoration")
	}
	if rec.Base.Tasks != 3 || rec.Aggregate.Tasks != 2 {
		t.Errorf("drift read wrong: base=%+v agg=%+v", rec.Base, rec.Aggregate)
	}

	fixed, err := m.RefreshDay(ctx, day)
	if err != nil {
		t.Fatal(err)
	}
	if fixed.Tasks != 3 || fixed.Micros != 950 {
		t.Errorf("refresh = %+v, want the base numbers (3 tasks / 950 micros)", fixed)
	}
	if rec, err = m.Reconcile(ctx, day); err != nil || !rec.Agree {
		t.Errorf("after repair: agree=%v err=%v", rec.Agree, err)
	}
}

func TestSetCostDayOverwritesAndRefusesNegatives(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	if err := s.BumpCostDay(ctx, "2026-11-02", 100, 10, 1000); err != nil {
		t.Fatal(err)
	}
	// Additive first (the hot path), then overwrite (the repair path): the row
	// ends where the base table says it should, not doubled.
	if err := s.SetCostDay(ctx, "2026-11-02", 30, 4, 300, 1); err != nil {
		t.Fatal(err)
	}
	got, err := s.CostDay(ctx, "2026-11-02")
	if err != nil {
		t.Fatal(err)
	}
	if got.TokensIn != 30 || got.TokensOut != 4 || got.CostMicros != 300 || got.Tasks != 1 {
		t.Errorf("SetCostDay did not replace the row: %+v", got)
	}
	for _, bad := range []struct{ in, out, micros, tasks int64 }{
		{-1, 0, 0, 0}, {0, -1, 0, 0}, {0, 0, -1, 0}, {0, 0, 0, -1},
	} {
		if err := s.SetCostDay(ctx, "2026-11-02", bad.in, bad.out, bad.micros, bad.tasks); err == nil {
			t.Errorf("negative aggregate %+v was written", bad)
		}
	}
	if err := s.SetCostDay(ctx, "2026-11-2", 1, 1, 1, 1); err == nil {
		t.Error("a malformed day label was written")
	}
	// A zero row is a stated fact, not a deletion.
	if err := s.SetCostDay(ctx, "2026-11-03", 0, 0, 0, 0); err != nil {
		t.Fatal(err)
	}
	rows, err := s.CostDaysBetween(ctx, "2026-11-03", "2026-11-03")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Tasks != 0 {
		t.Errorf("zero row not stored: %+v", rows)
	}
}

// setEndedAt rewrites a seeded task_log row's finish instant. FinishTaskLog
// stamps wall-clock "now" and the boundary cases need a fixed one (yesterday,
// the first second, the last second), so the test sets the column directly on
// the same table the DAO owns.
func setEndedAt(ctx context.Context, s *Store, id string, ended int64, in, out, micros int64) error {
	return s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			`UPDATE task_log SET ended_at=?, state=?, cost_tokens_in=?, cost_tokens_out=?, cost_amount_micro=? WHERE id=?`,
			ended, string(TaskStateDone), in, out, micros, id)
		return err
	})
}

func ptr(t time.Time) *time.Time { return &t }
