package memory

// Ticket 196 AC#1 (the single authority) and AC#2 (the CHECK + its migration
// path). Written against the three doors this ticket measured BEFORE any
// change, all three of which were open: StartTaskLog, FinishTaskLog with an
// empty error class, and raw SQL straight past the DAO. Every case below rings
// on the pre-196 tree and only goes quiet once both halves of the shape are in:
// the roster in taskstate.go (甲) and the CHECK in ddlV4 (乙).

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// theFourLegacyWords are the four the ticket names as "really written into the
// database today". They are listed here literally, not derived, so this leg
// still means something if the roster is ever edited: it is the ticket's
// reading, not a mirror of the code under test.
var theFourLegacyWords = []string{"done", "cancelled", "running", "succeeded"}

const illegal196 = "definitely-not-a-task-state"

// TestTaskState196RosterIsExplicit pins the roster's shape: unique, non-empty,
// and total over both vocabularies. A roster that silently lost a word, or
// gained a duplicate, fails here rather than at a user's database.
func TestTaskState196RosterIsExplicit(t *testing.T) {
	roster := TaskStateRoster()
	if len(roster) == 0 {
		t.Fatal("roster is empty: the CHECK would accept nothing")
	}
	seen := map[string]bool{}
	for _, w := range roster {
		if w == "" {
			t.Error("roster holds an empty word")
		}
		if seen[w] {
			t.Errorf("roster holds %q twice", w)
		}
		seen[w] = true
	}
	// The four the ticket names must all be accepted, and so must the D43 names
	// a task record may reach (the ticket's first hard boundary: the CHECK may
	// lock neither "only the four" nor "only D43's twenty").
	for _, w := range theFourLegacyWords {
		if !ValidTaskState(w) {
			t.Errorf("legacy word %q is not accepted: the CHECK would reject rows existing databases already hold", w)
		}
	}
	for _, d43 := range []statemachine.State{
		statemachine.StateThinking, statemachine.StateSettling,
		statemachine.StateMuted, statemachine.StateStuck, statemachine.StateError,
	} {
		if !ValidTaskState(string(d43)) {
			t.Errorf("D43 name %q is not accepted: a future migration would need a schema bump", string(d43))
		}
	}

	// The roster must NOT be the whole ball alphabet: Sleeping in a task row is
	// not a thing the ticket allows the CHECK to permit.
	for _, ball := range []statemachine.State{
		statemachine.StateSleeping, statemachine.StateListening, statemachine.StateArmed,
	} {
		if ValidTaskState(string(ball)) {
			t.Errorf("ball state %q is on the task roster: the CHECK now permits a state no task row can mean", string(ball))
		}
	}

	// A caller must not be able to widen the vocabulary through the accessor.
	got := TaskStateRoster()
	got[0] = "injected"
	if ValidTaskState("injected") {
		t.Error("TaskStateRoster handed out the live slice: a caller widened the vocabulary")
	}
}

// TestTaskState196MappingBothDirections is AC#1's read half: both vocabularies
// are retained and every translation goes through this one table.
func TestTaskState196MappingBothDirections(t *testing.T) {
	for _, w := range TaskStateRoster() {
		d43, ok := TaskStateFor(w)
		if !ok {
			t.Errorf("TaskStateFor(%q) has no answer although %q is on the roster", w, w)
			continue
		}
		if !statemachine.Valid(d43) {
			t.Errorf("TaskStateFor(%q) = %q, which is not one of D43's 20 names", w, string(d43))
		}
		back, ok := TaskStateFromD43(d43)
		if !ok {
			t.Errorf("TaskStateFromD43(%q) has no answer although it came from the roster", string(d43))
			continue
		}
		if !ValidTaskState(string(back)) {
			t.Errorf("TaskStateFromD43(%q) = %q, which is off the roster", string(d43), string(back))
		}
	}
	// The four legacy words and their D43 readings, spelled out: this is the
	// table a reader is supposed to consult instead of writing its own switch.
	for w, want := range map[string]statemachine.State{
		"running":   statemachine.StateThinking,
		"done":      statemachine.StateSettling,
		"succeeded": statemachine.StateSettling,
		"cancelled": statemachine.StateMuted,
	} {
		got, ok := TaskStateFor(w)
		if !ok || got != want {
			t.Errorf("TaskStateFor(%q) = (%q, %v), want (%q, true)", w, string(got), ok, string(want))
		}
	}
	// A word off the roster is never translated, and never silently maps to a
	// plausible-looking D43 name.
	if got, ok := TaskStateFor(illegal196); ok || got != "" {
		t.Errorf("TaskStateFor(%q) = (%q, %v), want (\"\", false)", illegal196, string(got), ok)
	}
	if _, ok := TaskStateFromD43(statemachine.StateSleeping); ok {
		t.Error("TaskStateFromD43(Sleeping) answered: a ball state is not a task state")
	}

	// The read-side door on the row itself.
	tl := TaskLog{State: "succeeded"}
	if got, ok := tl.D43State(); !ok || got != statemachine.StateSettling {
		t.Errorf("TaskLog.D43State() = (%q, %v), want (Settling, true)", string(got), ok)
	}
	if got, ok := (TaskLog{State: illegal196}).D43State(); ok || got != "" {
		t.Errorf("TaskLog.D43State() on garbage = (%q, %v), want (\"\", false)", string(got), ok)
	}
}

// TestTaskState196DAORejectsIllegalOnStart is door 1. Before 196 this returned
// nil and the row landed with the garbage in it.
func TestTaskState196DAORejectsIllegalOnStart(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	err := s.StartTaskLog(ctx, TaskLog{ID: "d1", State: illegal196, QueryText: "q"})
	if err == nil {
		t.Fatalf("StartTaskLog(state=%q) returned nil: the DAO still accepts any word", illegal196)
	}
	if !strings.Contains(err.Error(), illegal196) {
		t.Errorf("refusal %q does not quote the offending value", err.Error())
	}
	if _, err := s.TaskLogByID(ctx, "d1"); !errors.Is(err, ErrNotFound) {
		t.Errorf("a row was written despite the refusal: TaskLogByID err = %v", err)
	}
}

// TestTaskState196DAORejectsIllegalOnFinish is door 2, the widest one: the
// finish path used to validate ONLY when the caller passed an error class, so
// every clean exit wrote its state unchecked.
func TestTaskState196DAORejectsIllegalOnFinish(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	if err := s.StartTaskLog(ctx, TaskLog{ID: "d2", State: "running", QueryText: "q"}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	err := s.FinishTaskLog(ctx, "d2", illegal196, "s", 0, 0, 0, "")
	if err == nil {
		t.Fatalf("FinishTaskLog(state=%q, errorClass=\"\") returned nil: the clean-exit branch is still unguarded", illegal196)
	}
	if !strings.Contains(err.Error(), illegal196) {
		t.Errorf("refusal %q does not quote the offending value", err.Error())
	}
	got, err := s.TaskLogByID(ctx, "d2")
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if got.State != "running" {
		t.Errorf("state moved to %q despite the refusal, want it still \"running\"", got.State)
	}
	if got.EndedAt != nil {
		t.Error("ended_at was stamped despite the refusal")
	}
}

// TestTaskState196CheckRejectsIllegalRawSQL is door 3 and the point of 乙: a
// writer that never touches the DAO is refused by the database itself.
func TestTaskState196CheckRejectsIllegalRawSQL(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, e := tx.ExecContext(ctx,
			`INSERT INTO task_log(id, started_at, state, query_text) VALUES(?,?,?,?)`,
			"raw1", 1, illegal196, "q")
		return e
	})
	if err == nil {
		t.Fatalf("raw INSERT of state=%q succeeded: the column still has no CHECK", illegal196)
	}
	if !strings.Contains(strings.ToUpper(err.Error()), "CHECK") {
		t.Errorf("refusal %q does not mention the CHECK constraint", err.Error())
	}
	// The refusal text is recorded here so the evidence file can quote the
	// database's own words rather than paraphrase them.
	t.Logf("CHECK refusal (verbatim): %v", err)

	// And the same on UPDATE, so the hole is not merely closed on INSERT.
	if err := s.StartTaskLog(ctx, TaskLog{ID: "raw2", State: "running", QueryText: "q"}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	uerr := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, e := tx.ExecContext(ctx, `UPDATE task_log SET state=? WHERE id=?`, illegal196, "raw2")
		return e
	})
	if uerr == nil {
		t.Fatalf("raw UPDATE to state=%q succeeded: the CHECK holds on INSERT only", illegal196)
	}
}

// TestTaskState196CheckAcceptsTheFourLegacyWords is the CHECK's other
// direction: today's four words must still be writable, or the ticket would
// have "washed away" exactly what its forbidden zone protects.
func TestTaskState196CheckAcceptsTheFourLegacyWords(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	for i, w := range theFourLegacyWords {
		if err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
			_, e := tx.ExecContext(ctx,
				`INSERT INTO task_log(id, started_at, state, query_text) VALUES(?,?,?,?)`,
				"four-"+w, int64(i+1), w, "q")
			return e
		}); err != nil {
			t.Errorf("the CHECK rejected legacy word %q: %v", w, err)
		} else {
			t.Logf("legacy word %q accepted by the CHECK", w)
		}
	}
	var n int
	if err := s.reader.QueryRowContext(ctx, `SELECT count(*) FROM task_log`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != len(theFourLegacyWords) {
		t.Errorf("rows = %d, want %d (one per legacy word)", n, len(theFourLegacyWords))
	}
}

// TestTaskState196CheckAcceptsTheWholeRoster walks the roster through raw SQL,
// so a word the DAO accepts but the CHECK rejects cannot hide.
func TestTaskState196CheckAcceptsTheWholeRoster(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	for i, w := range TaskStateRoster() {
		if err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
			_, e := tx.ExecContext(ctx,
				`INSERT INTO task_log(id, started_at, state, query_text) VALUES(?,?,?,?)`,
				"roster-"+w, int64(i+1), w, "q")
			return e
		}); err != nil {
			t.Errorf("the CHECK rejected roster word %q, which the DAO accepts: the two mouths disagree: %v", w, err)
		}
	}
}

// TestTaskState196CheckIsDeclaredOnTheColumn proves the CHECK is really in the
// column's DDL - not merely in the Go validator - and that the value list it
// carries is the roster, word for word.
func TestTaskState196CheckIsDeclaredOnTheColumn(t *testing.T) {
	s := openTestStore(t)
	var ddl string
	if err := s.reader.QueryRowContext(context.Background(),
		`SELECT sql FROM sqlite_master WHERE type='table' AND name='task_log'`).Scan(&ddl); err != nil {
		t.Fatalf("read task_log DDL: %v", err)
	}
	if !strings.Contains(ddl, "CHECK") {
		t.Fatalf("task_log DDL carries no CHECK:\n%s", ddl)
	}
	for _, w := range TaskStateRoster() {
		if !strings.Contains(ddl, "'"+w+"'") {
			t.Errorf("roster word %q is missing from the CHECK the column declares", w)
		}
	}
	if strings.Contains(ddl, "'"+illegal196+"'") {
		t.Error("the garbage word appears in the CHECK: the roster was widened by accident")
	}
	// The column must still be NOT NULL: the rebuild had to change the CHECK
	// and nothing else.
	if !strings.Contains(ddl, "state          TEXT NOT NULL CHECK") {
		t.Errorf("state column is no longer TEXT NOT NULL CHECK:\n%s", ddl)
	}
}

// seedLegacyV3 builds a REAL schema-v3 database: the ddlV1/ddlV2/ddlV3 text a
// shipped v3 binary would have written, the watermark at 3, and task rows whose
// state holds the given words. Built from the production DDL constants on
// purpose - a hand-written fixture would drift away from the shape it is meant
// to represent.
func seedLegacyV3(t *testing.T, dir string, states []string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	db, err := sql.Open(DriverName, dsn(filepath.Join(dir, dbFileName), true))
	if err != nil {
		t.Fatalf("open raw: %v", err)
	}
	defer db.Close()
	for _, script := range []string{ddlV1, ddlV2, ddlV3} {
		for _, stmt := range splitSQLStatements(script) {
			if _, err := db.Exec(stmt); err != nil {
				t.Fatalf("seed v3 DDL: %v", err)
			}
		}
	}
	if _, err := db.Exec(`INSERT INTO schema_meta(key, value) VALUES('schema_version','3')`); err != nil {
		t.Fatalf("seed watermark: %v", err)
	}
	for i, w := range states {
		if _, err := db.Exec(
			`INSERT INTO task_log(id, started_at, state, query_text) VALUES(?,?,?,?)`,
			"legacy-"+w+"-"+string(rune('a'+i)), int64(i+1), w, "q"); err != nil {
			t.Fatalf("seed row %q: %v", w, err)
		}
	}
}

// TestTaskState196LegacyV3DatabaseOpensAndReads is the ticket's second hard
// boundary, measured rather than asserted: a database that already holds the
// four words must OPEN, migrate to v4, and read back unchanged.
func TestTaskState196LegacyV3DatabaseOpensAndReads(t *testing.T) {
	dir := filepath.Join(sealableTempDir124(t), "data")
	seedLegacyV3(t, dir, theFourLegacyWords)

	s, err := Open(dir)
	if err != nil {
		t.Fatalf("Open on a v3 database holding the four legacy words: %v", err)
	}
	defer s.Close()
	ctx := context.Background()

	var ver string
	if err := s.reader.QueryRowContext(ctx,
		`SELECT value FROM schema_meta WHERE key='schema_version'`).Scan(&ver); err != nil {
		t.Fatal(err)
	}
	if ver != "4" {
		t.Errorf("schema_version = %q, want 4 (the 3->4 step did not land)", ver)
	}

	rows, err := s.ListTaskLogs(ctx)
	if err != nil {
		t.Fatalf("ListTaskLogs: %v", err)
	}
	if len(rows) != len(theFourLegacyWords) {
		t.Fatalf("rows = %d, want %d: the rebuild lost or duplicated rows", len(rows), len(theFourLegacyWords))
	}
	got := map[string]bool{}
	for _, r := range rows {
		got[r.State] = true
	}
	for _, w := range theFourLegacyWords {
		if !got[w] {
			t.Errorf("legacy word %q did not survive the rebuild (got %v)", w, rows)
		}
	}

	// The CHECK must be in force on the migrated database, not only on a fresh
	// one - that is the difference between "migrated" and "rebuilt".
	if err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, e := tx.ExecContext(ctx,
			`INSERT INTO task_log(id, started_at, state, query_text) VALUES(?,?,?,?)`,
			"post-v4", 99, illegal196, "q")
		return e
	}); err == nil {
		t.Error("the migrated database still accepts an illegal state: the v4 CHECK was not applied")
	}

	// And the read-side authority agrees about the rows that came from the old
	// binary.
	for _, w := range theFourLegacyWords {
		if _, ok := TaskStateFor(w); !ok {
			t.Errorf("word %q survived the migration but has no D43 reading", w)
		}
	}
}

// TestTaskState196RefusesToMigrateAnOutOfRosterDatabase is the other side of
// compatibility: a database that does NOT satisfy the roster must fail loudly,
// naming the values, and leave the watermark where it was.
func TestTaskState196RefusesToMigrateAnOutOfRosterDatabase(t *testing.T) {
	dir := filepath.Join(sealableTempDir124(t), "data")
	seedLegacyV3(t, dir, []string{"running", illegal196})

	_, err := Open(dir)
	if err == nil {
		t.Fatal("Open on a v3 database holding an out-of-roster state succeeded")
	}
	if !strings.Contains(err.Error(), illegal196) {
		t.Errorf("migration refusal %q does not name the offending value", err.Error())
	}

	// The step must have rolled back: the watermark is still 3 and the garbage
	// row is still there, untouched. A half-rebuilt table would show up here.
	db, oerr := sql.Open(DriverName, dsn(filepath.Join(dir, dbFileName), false))
	if oerr != nil {
		t.Fatalf("reopen raw: %v", oerr)
	}
	defer db.Close()
	var ver string
	if err := db.QueryRow(`SELECT value FROM schema_meta WHERE key='schema_version'`).Scan(&ver); err != nil {
		t.Fatalf("read watermark: %v", err)
	}
	if ver != "3" {
		t.Errorf("schema_version = %q after a refused step, want 3 (the transaction must roll back)", ver)
	}
	var n int
	if err := db.QueryRow(`SELECT count(*) FROM task_log`).Scan(&n); err != nil {
		t.Fatalf("count rows: %v", err)
	}
	if n != 2 {
		t.Errorf("rows = %d after a refused step, want 2 (nothing copied, nothing dropped)", n)
	}
}

// TestTaskState196DDLV4KeepsTheShape is the guard on the rebuild itself: the
// generated script must still declare the same columns and the same index, so
// the introspection contract cannot be satisfied by a table that lost a column.
func TestTaskState196DDLV4KeepsTheShape(t *testing.T) {
	stmts := splitSQLStatements(ddlV4)
	if len(stmts) != 5 {
		t.Fatalf("ddlV4 produced %d statements, want 5 (create, copy, drop, rename, reindex): %q", len(stmts), stmts)
	}
	for _, want := range []string{
		"CREATE TABLE task_log_v4", "INSERT INTO task_log_v4",
		"DROP TABLE task_log", "ALTER TABLE task_log_v4 RENAME TO task_log",
		"CREATE INDEX idx_task_log_started ON task_log(started_at)",
	} {
		found := false
		for _, s := range stmts {
			if strings.Contains(s, want) {
				found = true
			}
		}
		if !found {
			t.Errorf("ddlV4 statement %q is missing; the rebuild is not the documented 12-step shape", want)
		}
	}
	// ddlV1 must be untouched: this ticket adds a step, it does not edit the
	// frozen contract text.
	if n := len(splitSQLStatements(ddlV1)); n != 15 {
		t.Errorf("splitSQLStatements(ddlV1) = %d, want 15: ddlV1 must not be edited by ticket 196", n)
	}
}
