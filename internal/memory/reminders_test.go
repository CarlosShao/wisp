package memory

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Ticket 31 (leg 31-r2) acceptance for the memory layer:
//
//	(a) the 5->6 migration creates `reminders` and keeps every older row (旧数据不丢);
//	(b) the step is idempotent (re-running applyV6 inside a tx neither errors nor restates rows);
//	(c) the status CHECK follows reminderStatusRoster (the ticket-196 rule, new enum);
//	(d) the DAO round-trips and the one-way armed -> fired | cancelled transitions hold.

// seedV5Database builds a REAL v5 database through the production chain (not a
// hand-copied DDL), seeds one cost_daily row and one profile row, and returns
// its directory so a later Open can migrate it 5->6.
func seedV5Database(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	s, err := Open(dir, withSchemaTarget(5))
	if err != nil {
		t.Fatalf("open v5 fixture: %v", err)
	}
	ctx := context.Background()
	if err := s.BumpCostDay(ctx, "2026-11-02", 100, 50, 7000); err != nil {
		t.Fatalf("seed cost_daily: %v", err)
	}
	if err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO profile(slot, content, source, updated_at) VALUES('pref.language','zh-CN','manual',?)`,
			time.Now().UTC().Unix())
		return err
	}); err != nil {
		t.Fatalf("seed profile: %v", err)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("close v5 fixture: %v", err)
	}
	return dir
}

func scalarAt(t *testing.T, path, q string, args ...any) int {
	t.Helper()
	db, err := sql.Open(DriverName, dsn(path, true))
	if err != nil {
		t.Fatalf("open reader: %v", err)
	}
	defer db.Close()
	var n int
	if err := db.QueryRow(q, args...).Scan(&n); err != nil {
		t.Fatalf("query %q: %v", q, err)
	}
	return n
}

// TestMigration31V5ToV6KeepsOldData is (a): the production Open runs 5->6, the
// reminders table exists, and the v5-era rows are byte-identical.
func TestMigration31V5ToV6KeepsOldData(t *testing.T) {
	dir := seedV5Database(t)
	// The fixture is the file the production Open just wrote, named the way this
	// platform names it (filepath.Join yields dir\wisp.db on Windows, byte for
	// byte the spelling this line used to build by hand, and dir/wisp.db
	// everywhere else - where the hand-built "dir\\wisp.db" opened a second,
	// empty database and read 0 reminders).
	dbPath := filepath.Join(dir, dbFileName)

	s, err := Open(dir)
	if err != nil {
		t.Fatalf("open (migrate v5->v6): %v", err)
	}
	ctx := context.Background()

	// The watermark landed on 6 through the production step.
	var v string
	if err := s.reader.QueryRowContext(ctx,
		`SELECT value FROM schema_meta WHERE key='schema_version'`).Scan(&v); err != nil {
		t.Fatal(err)
	}
	if v != "6" {
		t.Errorf("schema_version = %q, want 6", v)
	}

	// reminders exists and starts empty.
	if got := scalarAt(t, dbPath,
		`SELECT count(*) FROM sqlite_master WHERE type='table' AND name='reminders'`); got != 1 {
		t.Fatalf("reminders table after the 5->6 step = %d, want 1", got)
	}

	// The v5-era rows are untouched (旧数据不丢).
	cost, err := s.CostDay(ctx, "2026-11-02")
	if err != nil {
		t.Fatalf("read seeded cost_daily after migration: %v", err)
	}
	if cost.TokensIn != 100 || cost.TokensOut != 50 || cost.CostMicros != 7000 || cost.Tasks != 1 {
		t.Errorf("cost_daily after migration = %+v, want the seeded 100/50/7000/1", cost)
	}
	if got := scalarAt(t, dbPath, `SELECT count(*) FROM profile`); got != 1 {
		t.Errorf("profile rows after migration = %d, want 1", got)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}

	// Reopen is a no-op: no second 5->6 backup, no data movement.
	s2, err := Open(dir)
	if err != nil {
		t.Fatalf("reopen at v6: %v", err)
	}
	defer s2.Close()
	if got := scalarAt(t, dbPath, `SELECT count(*) FROM profile`); got != 1 {
		t.Errorf("profile rows after reopen = %d, want 1", got)
	}
}

// TestMigration31StepIsIdempotent is (b): applyV6 runs twice inside one
// transaction without error and without restating a row. This is the shape a
// retried 5->6 transaction takes.
func TestMigration31StepIsIdempotent(t *testing.T) {
	dir := t.TempDir()
	s, err := Open(dir)
	if err != nil {
		t.Fatalf("open (production chain to v6): %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	ctx := context.Background()
	if err := s.InsertReminder(ctx, Reminder{
		ID: "rem_idempotent", FireAt: time.Now().UTC().Add(time.Hour).Unix(),
		Note: "幂等复跑不重述这一行", Status: ReminderArmed,
	}); err != nil {
		t.Fatalf("seed reminder: %v", err)
	}
	tx, err := s.writer.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // the tx is only a probe
	for i := 0; i < 2; i++ {
		if err := applyV6(tx); err != nil {
			t.Fatalf("applyV6 re-run #%d failed (the step is not idempotent): %v", i+1, err)
		}
	}
	var n int
	if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM reminders`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Errorf("reminder rows after two applyV6 runs = %d, want 1 (a step must never restate rows)", n)
	}
}

// TestReminderStatusCheckFollowsRoster is (c): the DDL CHECK and
// reminderStatusRoster cannot fork (ticket 196's rule, first use on a new
// enum), and a value outside the roster is refused by SQLite itself.
func TestReminderStatusCheckFollowsRoster(t *testing.T) {
	ddl := string(ddlV6)
	for _, st := range reminderStatusRoster {
		if !strings.Contains(ddl, "'"+st+"'") {
			t.Errorf("the v6 CHECK does not carry roster value %q: %s", st, ddl)
		}
	}
	// Strip -- comments before the recurrence guard: the DDL deliberately
	// documents "无 recurrence/cron 列" in comments, which must not trip a
	// substring check meant for actual SQL keywords/shapes.
	var code strings.Builder
	for _, ln := range strings.Split(ddl, "\n") {
		if i := strings.Index(ln, "--"); i >= 0 {
			ln = ln[:i]
		}
		code.WriteString(ln)
		code.WriteByte('\n')
	}
	if strings.Contains(code.String(), "'recurring'") || strings.Contains(code.String(), "cron") {
		t.Errorf("the v6 DDL grew a recurrence shape, which ticket 31 rejects:\n%s", ddl)
	}

	s := openTestStore(t)
	ctx := context.Background()
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO reminders(id, fire_at_wallclock_utc, note, status, created_utc)
			 VALUES('rem_check', ?, 'x', 'repeating', ?)`,
			time.Now().UTC().Add(time.Hour).Unix(), time.Now().UTC().Unix())
		return err
	})
	if err == nil {
		t.Error("CHECK accepted a status outside the roster")
	} else if !strings.Contains(err.Error(), "CHECK constraint failed") {
		t.Errorf("rejection did not come from the CHECK: %v", err)
	}
}

// TestReminderDAORoundTripAndOneWayTransitions is (d).
func TestReminderDAORoundTripAndOneWayTransitions(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	now := time.Now().UTC()

	future := Reminder{ID: "rem_future", FireAt: now.Add(2 * time.Hour).Unix(), Note: "两小时后", Status: ReminderArmed}
	if err := s.InsertReminder(ctx, future); err != nil {
		t.Fatalf("insert: %v", err)
	}
	overdue := Reminder{ID: "rem_overdue", FireAt: now.Add(-time.Minute).Unix(), Note: "已过期"}
	if err := s.InsertReminder(ctx, overdue); err != nil {
		t.Fatalf("insert overdue (default status): %v", err)
	}

	got, err := s.ReminderByID(ctx, "rem_future")
	if err != nil {
		t.Fatalf("by id: %v", err)
	}
	if got.Status != ReminderArmed || got.FireAt != future.FireAt || got.Note != future.Note ||
		got.CreatedAt == 0 || got.FiredAt != nil {
		t.Errorf("round-trip = %+v, want the inserted armed row with fired_utc NULL", got)
	}

	list, err := s.ListReminders(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 2 || list[0].ID != "rem_overdue" {
		t.Errorf("list = %+v, want 2 rows in fire_at order (overdue first)", list)
	}

	// CountMissedArmed counts only armed rows whose wall-clock instant passed.
	missed, err := s.CountMissedArmed(ctx, now.Unix())
	if err != nil {
		t.Fatalf("count missed: %v", err)
	}
	if missed != 1 {
		t.Errorf("CountMissedArmed = %d, want 1", missed)
	}

	// armed -> fired is one-way: the second MarkFired reports "not delivered".
	applied, err := s.MarkReminderFired(ctx, "rem_overdue")
	if err != nil || !applied {
		t.Fatalf("mark fired = (%v, %v), want (true, nil)", applied, err)
	}
	again, err := s.MarkReminderFired(ctx, "rem_overdue")
	if err != nil || again {
		t.Errorf("second mark fired = (%v, %v), want (false, nil): fired is terminal", again, err)
	}
	// a fired row cannot be cancelled
	_, err = s.CancelReminder(ctx, "rem_overdue")
	if err != nil {
		t.Fatalf("cancel on a fired row is a no-op, not an error: %v", err)
	}
	// the fired row carries fired_utc
	fired, err := s.ReminderByID(ctx, "rem_overdue")
	if err != nil {
		t.Fatalf("read fired row: %v", err)
	}
	if fired.Status != ReminderFired || fired.FiredAt == nil {
		t.Errorf("fired row = %+v, want status fired with fired_utc stamped", fired)
	}

	// armed -> cancelled
	applied, err = s.CancelReminder(ctx, "rem_future")
	if err != nil || !applied {
		t.Fatalf("cancel = (%v, %v), want (true, nil)", applied, err)
	}
	again, err = s.CancelReminder(ctx, "rem_future")
	if err != nil || again {
		t.Errorf("second cancel = (%v, %v), want (false, nil): cancelled is terminal", again, err)
	}
	// a cancelled reminder is no longer missed
	missed, err = s.CountMissedArmed(ctx, now.Add(time.Hour).Unix())
	if err != nil {
		t.Fatal(err)
	}
	if missed != 0 {
		t.Errorf("CountMissedArmed after cancel = %d, want 0: only armed rows count", missed)
	}

	// Refused shapes never half-write.
	if err := s.InsertReminder(ctx, Reminder{ID: "", FireAt: now.Unix(), Note: "x"}); err == nil {
		t.Error("empty id accepted")
	}
	if err := s.InsertReminder(ctx, Reminder{ID: "rem_bad", FireAt: 0, Note: "x"}); err == nil {
		t.Error("non-positive fire_at accepted")
	}
	if err := s.InsertReminder(ctx, Reminder{ID: "rem_bad", FireAt: now.Unix(), Note: ""}); err == nil {
		t.Error("empty note accepted")
	}
	if err := s.InsertReminder(ctx, Reminder{ID: "rem_bad", FireAt: now.Unix(), Note: "x", Status: "repeating"}); err == nil {
		t.Error("status outside the roster accepted")
	}
	if _, err := s.ReminderByID(ctx, "rem_bad"); err == nil {
		t.Error("a refused insert left a row behind")
	}
}
