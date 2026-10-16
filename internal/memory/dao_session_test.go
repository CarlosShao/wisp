package memory

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"
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

func ptrTime(v time.Time) *time.Time { return &v }

func sessionIDList(ss []Session) []string {
	out := make([]string, len(ss))
	for i, s := range ss {
		out[i] = s.ID
	}
	return out
}

// --- ticket 326 AC#1: appending a turn is a "the conversation was touched" ---

// TestAppendSessionTurnSurfacesTheConversation plants the shape ticket 326
// reports: session A created first, session B created second, then ONE more
// turn inside A. The list sorts by `updated_at DESC`, so a conversation somebody
// is still talking in must come back first — and "the order looked right" is not
// the finish line: the same read asserts A.updated_at equals the stamp of the
// turn it just gained. That second assertion is what keeps the case honest
// against the two opposite shapes that also put A on top by accident: refreshing
// to a constant, and refreshing to A's own created_at (a no-op here, because
// InsertSession seeds updated_at = created_at).
//
// Stamps are one minute apart on purpose. Two rows sharing a second tie on
// updated_at and the ORDER BY falls through to "id DESC", where id is random hex
// — ticket 166 AC#6's coin, which ticket 326 AC#2 explicitly keeps OUT of this
// scope, so the case pins its own clock instead of leaning on the wall clock.
func TestAppendSessionTurnSurfacesTheConversation(t *testing.T) {
	ctx := context.Background()
	s := openTestStore(t)

	clock := ptrTime(time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC))
	stamp := func() int64 {
		*clock = clock.Add(time.Minute)
		return clock.Unix()
	}

	aAt, bAt := stamp(), stamp()
	for _, sess := range []Session{
		{ID: "sess-a", Title: "the one I keep talking in", CreatedAt: aAt, UpdatedAt: aAt},
		{ID: "sess-b", Title: "created a minute later", CreatedAt: bAt, UpdatedAt: bAt},
	} {
		if err := s.InsertSession(ctx, sess); err != nil {
			t.Fatalf("InsertSession(%s): %v", sess.ID, err)
		}
	}
	before, err := s.ListSessions(ctx, SessionListOptions{})
	if err != nil {
		t.Fatalf("ListSessions: %v", err)
	}
	if len(before) != 2 || before[0].ID != "sess-b" {
		t.Fatalf("premise broken: list = %v, want sess-b first before the append", sessionIDList(before))
	}

	turnAt := stamp()
	if _, err := s.AppendSessionTurn(ctx, SessionTurn{
		SessionID: "sess-a", Role: "user", Content: "still here", CreatedAt: turnAt,
	}); err != nil {
		t.Fatalf("AppendSessionTurn: %v", err)
	}

	row, err := s.SessionByID(ctx, "sess-a")
	if err != nil {
		t.Fatalf("SessionByID(sess-a): %v", err)
	}
	after, err := s.ListSessions(ctx, SessionListOptions{})
	if err != nil {
		t.Fatalf("ListSessions after append: %v", err)
	}
	t.Logf("326 AC#1 same-row evidence: sess-a.updated_at %d -> %d (delta=%d) for turn.created_at=%d | list = %v",
		aAt, row.UpdatedAt, row.UpdatedAt-aAt, turnAt, sessionIDList(after))

	if row.UpdatedAt != turnAt {
		t.Errorf("appending a turn did not reach session.updated_at: sess-a.updated_at = %d, want the appended turn's stamp %d (before=%d, delta=%d)",
			row.UpdatedAt, turnAt, aAt, row.UpdatedAt-aAt)
	}
	if len(after) != 2 || after[0].ID != "sess-a" {
		t.Errorf("continuing a chat does not surface it: list = %v, want sess-a first after one more turn in it",
			sessionIDList(after))
	}
	if b, err := s.SessionByID(ctx, "sess-b"); err != nil || b.UpdatedAt != bAt {
		t.Errorf("the append disturbed the OTHER row: sess-b.updated_at = %d, want the untouched %d (%v)",
			b.UpdatedAt, bAt, err)
	}
}

// TestAppendSessionTurnRefreshUsesStoredClockFamily: the caller may hand the turn
// no timestamp at all (CreatedAt == 0), which orNow fills from
// observe.NowWallUTC().Unix(). The refresh must ride THAT same value — one clock
// per write, no third one — so the assertions are "updated_at == the created_at
// the row itself got" and "it moved forward", read back from the database.
func TestAppendSessionTurnRefreshUsesStoredClockFamily(t *testing.T) {
	ctx := context.Background()
	s := openTestStore(t)
	at := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC).Unix()
	if err := s.InsertSession(ctx, Session{ID: "s", Title: "no timestamp on the turn", CreatedAt: at, UpdatedAt: at}); err != nil {
		t.Fatalf("InsertSession: %v", err)
	}
	if _, err := s.AppendSessionTurn(ctx, SessionTurn{SessionID: "s", Role: "user", Content: "hello"}); err != nil {
		t.Fatalf("AppendSessionTurn: %v", err)
	}
	turns, err := s.ListSessionTurns(ctx, "s")
	if err != nil || len(turns) != 1 {
		t.Fatalf("ListSessionTurns = %+v (%v), want 1 turn", turns, err)
	}
	row, err := s.SessionByID(ctx, "s")
	if err != nil {
		t.Fatalf("SessionByID: %v", err)
	}
	t.Logf("326 AC#1 clock family: turn.created_at = %d, sess.updated_at %d -> %d (delta=%d)",
		turns[0].CreatedAt, at, row.UpdatedAt, row.UpdatedAt-at)
	if row.UpdatedAt != turns[0].CreatedAt {
		t.Errorf("refresh used a different clock than the turn it wrote: updated_at = %d, want turn.created_at %d",
			row.UpdatedAt, turns[0].CreatedAt)
	}
	if row.UpdatedAt <= at {
		t.Errorf("updated_at did not move forward: %d -> %d (appending is a touch)", at, row.UpdatedAt)
	}
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
