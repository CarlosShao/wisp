package memory

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
)

// Session-storage cases for the v3 (ticket 166) schema.

// TestSessionSearchEscapesWildcards: a query carrying LIKE wildcards must be
// treated as literal text, not as a pattern (models.go escapeLike).
func TestSessionSearchEscapesWildcards(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	for _, sess := range []Session{
		{ID: "s1", Title: "100% done", Workspace: "w"},
		{ID: "s2", Title: "100 done", Workspace: "w"},
	} {
		if err := s.InsertSession(ctx, sess); err != nil {
			t.Fatalf("InsertSession: %v", err)
		}
	}
	got, err := s.SearchSessions(ctx, "%", SessionListOptions{})
	if err != nil {
		t.Fatalf("SearchSessions: %v", err)
	}
	if len(got) != 1 || got[0].ID != "s1" {
		t.Errorf("searching a literal %% returned %d rows, want only s1 (wildcards must be escaped)", len(got))
	}
}

func hasTable(t *testing.T, db *sql.DB, name string) bool {
	t.Helper()
	var n int
	if err := db.QueryRow(
		`SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?`, name).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n > 0
}

// TestSessionMigrationCrashAtomicAndReplay is ticket 166 AC#3's evidence. The
// v3 step is ONE transaction: aborting it in the middle leaves no half-created
// table and the watermark unmoved, and a later clean open replays it whole.
//
// New and old databases both reach v3 through the same 2->3 step (ddlV1/ddlV2
// are never edited), so neither side can "lack a table the other has" — the
// divergence AC#3's second half warns about needs a CREATE for fresh databases
// and an ALTER for upgraded ones, which this shape simply does not have.
func TestSessionMigrationCrashAtomicAndReplay(t *testing.T) {
	dir := filepath.Join(sealableTempDir124(t), "data")

	// 1. A genuine v2 database has no session tables yet.
	old, err := Open(dir, withSchemaTarget(2))
	if err != nil {
		t.Fatalf("open v2: %v", err)
	}
	if hasTable(t, old.reader, "session") {
		t.Fatal("v2 target unexpectedly created the session table (premise broken)")
	}
	if err := old.Close(); err != nil {
		t.Fatal(err)
	}

	// 2. Fault-inject the REAL 2->3: run ddlV3, then abort before commit.
	saved := migrationChain
	t.Cleanup(func() { migrationChain = saved })
	migrationChain = []migrationStep{
		{from: 0, to: 1, apply: applyV1},
		{from: 1, to: 2, apply: applyV2},
		{from: 2, to: 3, apply: func(tx txExec) error {
			if err := applyV3(tx); err != nil {
				return err
			}
			// NOT NULL on session.title: this statement must take the whole
			// step (tables + watermark) down with it.
			_, err := tx.ExecContext(context.Background(),
				`INSERT INTO session(id, title, created_at, updated_at) VALUES('s', NULL, 0, 0)`)
			return err
		}},
	}
	if s, err := Open(dir, withSchemaTarget(3)); err == nil {
		_ = s.Close()
		t.Fatal("Open succeeded despite a failing v3 step")
	}

	// 3. All-or-nothing: watermark still 2, no session table, pre-backup kept.
	db, err := sql.Open(DriverName, dsn(filepath.Join(dir, dbFileName), true))
	if err != nil {
		t.Fatal(err)
	}
	var v int
	if err := db.QueryRow(`SELECT CAST(value AS INTEGER) FROM schema_meta WHERE key='schema_version'`).Scan(&v); err != nil {
		t.Fatal(err)
	}
	if v != 2 {
		t.Errorf("schema_version after aborted v3 = %d, want 2 (no half migration)", v)
	}
	if hasTable(t, db, "session") || hasTable(t, db, "session_turn") {
		t.Error("a session table survived an aborted v3 step")
	}
	db.Close()
	if _, err := os.Stat(filepath.Join(dir, "backup", "wisp.db.bak-2-3")); err != nil {
		t.Errorf("pre-migration backup missing: %v", err)
	}

	// 4. Reopen with the production chain: the step replays whole.
	migrationChain = saved
	replayed, err := Open(dir, withSchemaTarget(3))
	if err != nil {
		t.Fatalf("replay v3: %v", err)
	}
	defer replayed.Close()
	if !hasTable(t, replayed.reader, "session") || !hasTable(t, replayed.reader, "session_turn") {
		t.Error("session tables missing after replaying the v3 step")
	}
	if err := replayed.InsertSession(context.Background(),
		Session{ID: "s", Title: "after replay"}); err != nil {
		t.Errorf("v3 tables not usable after replay: %v", err)
	}
}
