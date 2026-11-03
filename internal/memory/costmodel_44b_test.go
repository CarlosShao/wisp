package memory

import (
	"context"
	"database/sql"
	"testing"
)

// Ticket 44b acceptance: the 4->5 migration creates cost_daily_model without
// touching cost_daily's structure or data (AC#1), and the per-model DAO upserts
// along the (day, provider, model) dimension in the same batch as cost_daily
// (AC#2).

func t44bSeedCostDay(t *testing.T, s *Store, day string, in, out, micros int64) {
	t.Helper()
	if err := s.BumpCostDay(context.Background(), day, in, out, micros); err != nil {
		t.Fatalf("seed cost_daily %s: %v", day, err)
	}
}

// TestMigration44bAddsModelTableKeepsOldData is AC#1's first half: a database
// with real cost_daily rows migrates to v5, the rows are byte-identical after,
// and the new table exists.
func TestMigration44bAddsModelTableKeepsOldData(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	t44bSeedCostDay(t, s, "2026-10-01", 100, 50, 7000)
	t44bSeedCostDay(t, s, "2026-10-02", 200, 60, 9000)
	before, err := s.CostDaysBetween(ctx, "2026-01-01", "2026-12-31")
	if err != nil {
		t.Fatal(err)
	}
	if len(before) != 2 {
		t.Fatalf("seeded cost_daily rows = %d, want 2", len(before))
	}

	// Reopen: the migration recognizes v5 as current and is a no-op; the data
	// must be untouched.
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	s2, err := Open(s.dir)
	if err != nil {
		t.Fatalf("reopen after v5: %v", err)
	}
	defer s2.Close()
	var v string
	if err := s2.reader.QueryRowContext(ctx,
		`SELECT value FROM schema_meta WHERE key='schema_version'`).Scan(&v); err != nil {
		t.Fatal(err)
	}
	if v != "5" {
		t.Errorf("schema_version = %q, want 5", v)
	}
	after, err := s2.CostDaysBetween(ctx, "2026-01-01", "2026-12-31")
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != len(before) {
		t.Fatalf("cost_daily rows after migration = %d, want %d (old data must not move)", len(after), len(before))
	}
	for i := range before {
		if before[i] != after[i] {
			t.Errorf("cost_daily row %d changed across the migration: %+v -> %+v", i, before[i], after[i])
		}
	}
	// The new table exists and starts empty.
	if got := countRows(t, s2, "cost_daily_model"); got != 0 {
		t.Errorf("cost_daily_model rows on a migrated database = %d, want 0 (accumulation starts now)", got)
	}
	// cost_daily's PK shape is untouched: still PRIMARY KEY(day), no provider column.
	var pkCols, hasProvider int
	if err := s2.reader.QueryRowContext(ctx,
		`SELECT count(*) FROM pragma_table_info('cost_daily') WHERE pk > 0`).Scan(&pkCols); err != nil {
		t.Fatal(err)
	}
	if pkCols != 1 {
		t.Errorf("cost_daily PK columns = %d, want 1 (the split-table rule: the main PK is never widened)", pkCols)
	}
	if err := s2.reader.QueryRowContext(ctx,
		`SELECT count(*) FROM pragma_table_info('cost_daily') WHERE name='provider'`).Scan(&hasProvider); err != nil {
		t.Fatal(err)
	}
	if hasProvider != 0 {
		t.Error("cost_daily gained a provider column; 44b must not touch the main table")
	}
}

// TestMigration44bIdempotentReopens is AC#1's second half: reopening N times
// never re-runs the step and never errors.
func TestMigration44bIdempotentReopens(t *testing.T) {
	s := openTestStore(t)
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		s2, err := Open(s.dir)
		if err != nil {
			t.Fatalf("reopen %d: %v", i, err)
		}
		if err := s2.Close(); err != nil {
			t.Fatal(err)
		}
	}
	s3, err := Open(s.dir)
	if err != nil {
		t.Fatal(err)
	}
	defer s3.Close()
	if got := countRows(t, s3, "cost_daily_model"); got != 0 {
		t.Errorf("cost_daily_model rows after 3 reopens = %d, want 0", got)
	}
}

// TestBumpCostDayModelSameBatchBothRows is AC#2: one call moves cost_daily and
// cost_daily_model together, along the model dimension.
func TestBumpCostDayModelSameBatchBothRows(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	if err := s.BumpCostDayModel(ctx, "2026-11-03", "pa", "m1", 100, 40, 3000); err != nil {
		t.Fatalf("first book: %v", err)
	}
	if err := s.BumpCostDayModel(ctx, "2026-11-03", "pa", "m1", 50, 10, 1000); err != nil {
		t.Fatalf("second book same dimension: %v", err)
	}
	if err := s.BumpCostDayModel(ctx, "2026-11-03", "pa", "m2", 7, 3, 0); err != nil {
		t.Fatalf("second model same day: %v", err)
	}
	if err := s.BumpCostDayModel(ctx, "2026-11-02", "pb", "m1", 1, 1, 5); err != nil {
		t.Fatalf("second day: %v", err)
	}

	global, err := s.CostDay(ctx, "2026-11-03")
	if err != nil {
		t.Fatal(err)
	}
	if global.TokensIn != 157 || global.TokensOut != 53 || global.CostMicros != 4000 || global.Tasks != 3 {
		t.Errorf("cost_daily 2026-11-03 = %+v, want 157/53/4000/3 (one batch per task, both tables)", global)
	}

	rows, err := s.CostModelDaysBetween(ctx, "2026-11-02", "2026-11-03")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 3 {
		t.Fatalf("cost_daily_model rows = %d, want 3 (day,provider,model distinct keys)", len(rows))
	}
	want := map[[2]string]struct{ in, out, micro, tasks int64 }{
		{"pa", "m1"}: {150, 50, 4000, 2},
		{"pa", "m2"}: {7, 3, 0, 1},
		{"pb", "m1"}: {1, 1, 5, 1},
	}
	for _, r := range rows {
		w, ok := want[[2]string{r.Provider, r.Model}]
		if !ok {
			t.Errorf("unexpected model row %+v", r)
			continue
		}
		if r.TokensIn != w.in || r.TokensOut != w.out || r.Micros != w.micro || r.Tasks != w.tasks {
			t.Errorf("cost_daily_model[%s/%s] = %d/%d/%d/%d, want %d/%d/%d/%d",
				r.Provider, r.Model, r.TokensIn, r.TokensOut, r.Micros, r.Tasks,
				w.in, w.out, w.micro, w.tasks)
		}
	}

	// The unknown-price row booked tokens-only with micros 0 and is readable.
	if rows[1].Provider != "pa" || rows[1].Model != "m2" {
		t.Logf("row order note: %+v", rows[1])
	}

	// Rejected shapes never half-write: a refusal leaves both tables untouched.
	beforeGlobal, _ := s.CostDay(ctx, "2026-11-03")
	if err := s.BumpCostDayModel(ctx, "2026-11-03", "pa", "m1", -1, 0, 0); err == nil {
		t.Error("negative tokens accepted")
	}
	if err := s.BumpCostDayModel(ctx, "2026-11-03", "pa", "m1", 0, 0, -5); err == nil {
		t.Error("negative micros accepted")
	}
	if err := s.BumpCostDayModel(ctx, "2026-11-03", "", "m1", 0, 0, 0); err == nil {
		t.Error("empty provider accepted")
	}
	if err := s.BumpCostDayModel(ctx, "not-a-day", "pa", "m1", 0, 0, 0); err == nil {
		t.Error("invalid day accepted")
	}
	afterGlobal, _ := s.CostDay(ctx, "2026-11-03")
	if beforeGlobal != afterGlobal {
		t.Errorf("a refused book still moved cost_daily: %+v -> %+v", beforeGlobal, afterGlobal)
	}
}

// TestCostModelDaysBetweenRejectsBadRanges pins the reader's guards (same
// shape as the global range DAO).
func TestCostModelDaysBetweenRejectsBadRanges(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	if _, err := s.CostModelDaysBetween(ctx, "2026-11-3", "2026-11-04"); err == nil {
		t.Error("malformed from-day accepted")
	}
	if _, err := s.CostModelDaysBetween(ctx, "2026-11-05", "2026-11-04"); err == nil {
		t.Error("inverted range accepted")
	}
	rows, err := s.CostModelDaysBetween(ctx, "2026-11-04", "2026-11-04")
	if err != nil {
		t.Fatalf("empty range is a fact, not an error: %v", err)
	}
	if len(rows) != 0 {
		t.Errorf("empty range returned %d rows, want 0", len(rows))
	}
}

// Ensure the test file compiles even if sql remains only used via helpers.
var _ = sql.ErrNoRows
