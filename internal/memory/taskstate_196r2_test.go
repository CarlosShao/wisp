package memory

// Ticket 196 leg r2: the v4 rebuild drops task_log while tool_call.task_id
// references it (ddlV1). Today that is safe only because PRAGMA foreign_keys is
// off - nothing in this repository ever sets it (dsn() in open.go lists four
// pragmas and not that one). This file makes that fact load-bearing in two
// directions:
//
//  1. TestTaskState196LegacyV3WithToolCallRowsStillOpens is the tripwire. It
//     runs the production Open() on a v3 database that holds a child row. It is
//     green today, and it goes red - with SQLite's opaque "FOREIGN KEY
//     constraint failed" - the day anyone adds the pragma to the migration DSN.
//     That is the "measurable failure" the r1 verdict asked for in residual 4.
//
//  2. TestTaskState196RefusesRebuildWithForeignKeysOn pins what happens on that
//     day: applyV4 refuses at its first statement and says why, instead of
//     dying three statements in with a message that names neither table nor row.
//
// Both fixtures are built from the production DDL constants (seedLegacyV3), not
// from hand-written SQL, so they cannot drift away from the shape they represent.

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"testing"
)

// addLegacyToolCallRow inserts one tool_call row whose task_id is read OUT of
// task_log, so the row genuinely references a task the rebuild will copy and
// then drop the parent of. A hand-picked id would make both tests below pass for
// the wrong reason - an orphan child does not make DROP TABLE violate anything.
func addLegacyToolCallRow(t *testing.T, dir string) string {
	t.Helper()
	db, err := sql.Open(DriverName, dsn(filepath.Join(dir, dbFileName), true))
	if err != nil {
		t.Fatalf("open raw: %v", err)
	}
	defer db.Close()
	var parent string
	if err := db.QueryRow(`SELECT id FROM task_log ORDER BY id LIMIT 1`).Scan(&parent); err != nil {
		t.Fatalf("read a task_log id to hang the child row on: %v", err)
	}
	if _, err := db.Exec(
		`INSERT INTO tool_call(task_id, seq, tool, args_json, risk_level, correlation_id)
		 VALUES(?,?,?,?,?,?)`,
		parent, 1, "fs.read", `{"path":"x"}`, "L0", "corr-196r2"); err != nil {
		t.Fatalf("seed tool_call: %v", err)
	}
	return parent
}

// TestTaskState196LegacyV3WithToolCallRowsStillOpens is residual 4 turned into a
// test: a v3 database that already holds a tool_call row must still open and
// reach v4, with the child row alive and still pointing at the same task.
//
// Why it rings on the day the hazard becomes real: Open() migrates on a
// connection whose foreign-key switch comes from dsn(). With the switch off (as
// today) DROP TABLE does not look at the child rows; with it on, the same call
// aborts inside the 3->4 step with "FOREIGN KEY constraint failed (787)". So the
// pragma turning on in open.go turns this test red - no other code change needed.
func TestTaskState196LegacyV3WithToolCallRowsStillOpens(t *testing.T) {
	dir := filepath.Join(sealableTempDir124(t), "data")
	seedLegacyV3(t, dir, theFourLegacyWords)
	parent := addLegacyToolCallRow(t, dir)

	s, err := Open(dir)
	if err != nil {
		t.Fatalf("Open on a v3 database holding a tool_call row: %v\n"+
			"if this says \"FOREIGN KEY constraint failed\", the migration DSN now enables "+
			"PRAGMA foreign_keys and the v4 rebuild can no longer drop task_log - see "+
			"requireForeignKeysOff in schema.go", err)
	}
	defer s.Close()
	ctx := context.Background()

	var ver string
	if err := s.reader.QueryRowContext(ctx,
		`SELECT value FROM schema_meta WHERE key='schema_version'`).Scan(&ver); err != nil {
		t.Fatal(err)
	}
	if ver != "4" {
		t.Errorf("schema_version = %q, want 4: the 3->4 step did not land", ver)
	}

	var rows, matched int
	if err := s.reader.QueryRowContext(ctx, `SELECT count(*) FROM tool_call`).Scan(&rows); err != nil {
		t.Fatal(err)
	}
	if err := s.reader.QueryRowContext(ctx,
		`SELECT count(*) FROM tool_call WHERE task_id=?`, parent).Scan(&matched); err != nil {
		t.Fatal(err)
	}
	if rows != 1 || matched != 1 {
		t.Errorf("tool_call rows = %d (referencing %q: %d), want 1 and 1: the rebuild lost or "+
			"repointed the child row", rows, parent, matched)
	}

	var debris int
	if err := s.reader.QueryRowContext(ctx,
		`SELECT count(*) FROM sqlite_master WHERE name LIKE 'task_log_v4%'`).Scan(&debris); err != nil {
		t.Fatal(err)
	}
	if debris != 0 {
		t.Errorf("leftover task_log_v4 objects = %d, want 0", debris)
	}
}

// TestTaskState196RefusesRebuildWithForeignKeysOn drives the real 3->4 step
// (applyMigrationStep, the same function migrate() calls) twice over two
// identical fixtures. The only difference between them is the foreign-key
// switch, so the two outcomes can only be caused by that switch.
func TestTaskState196RefusesRebuildWithForeignKeysOn(t *testing.T) {
	step, ok := nextStep(3)
	if !ok {
		t.Fatal("the migration chain has no 3->4 step")
	}
	logger, _ := bufLogger(t)
	ctx := context.Background()

	// Control: production DSN, no FK pragma. The step must apply - otherwise the
	// subject half below would be measuring a broken fixture rather than the
	// pragma.
	dirOff := filepath.Join(sealableTempDir124(t), "data")
	seedLegacyV3(t, dirOff, theFourLegacyWords)
	addLegacyToolCallRow(t, dirOff)
	dbOff, err := sql.Open(DriverName, dsn(filepath.Join(dirOff, dbFileName), true))
	if err != nil {
		t.Fatalf("open control: %v", err)
	}
	defer dbOff.Close()
	if _, err := applyMigrationStep(dbOff, step, logger); err != nil {
		t.Fatalf("control: the same fixture failed to migrate with foreign_keys OFF: %v", err)
	}
	var offVersion, offChildren int
	if err := dbOff.QueryRowContext(ctx,
		`SELECT value FROM schema_meta WHERE key='schema_version'`).Scan(&offVersion); err != nil {
		t.Fatalf("control watermark: %v", err)
	}
	if err := dbOff.QueryRowContext(ctx, `SELECT count(*) FROM tool_call`).Scan(&offChildren); err != nil {
		t.Fatalf("control children: %v", err)
	}
	if offVersion != 4 || offChildren != 1 {
		t.Fatalf("control: schema_version = %d, tool_call = %d, want 4 and 1", offVersion, offChildren)
	}

	// Subject: the same fixture, the only change being the pragma.
	dirOn := filepath.Join(sealableTempDir124(t), "data")
	seedLegacyV3(t, dirOn, theFourLegacyWords)
	addLegacyToolCallRow(t, dirOn)
	dbOn, err := sql.Open(DriverName, dsn(filepath.Join(dirOn, dbFileName), true)+"&_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatalf("open subject: %v", err)
	}
	defer dbOn.Close()
	var fk int
	if err := dbOn.QueryRowContext(ctx, pragmaForeignKeysQuery).Scan(&fk); err != nil {
		t.Fatalf("read subject pragma: %v", err)
	}
	if fk != 1 {
		t.Fatalf("subject connection has foreign_keys = %d, want 1: this half would measure nothing", fk)
	}

	_, err = applyMigrationStep(dbOn, step, logger)
	if err == nil {
		t.Fatal("the 3->4 step ran to completion with foreign_keys ON and a tool_call row present: " +
			"either the pragma stopped taking effect or the rebuild stopped dropping task_log - " +
			"re-measure before trusting the guard")
	}
	msg := err.Error()
	// The readable form names the switch, the child clause and the fix. The raw
	// SQLite failure ("constraint failed: FOREIGN KEY constraint failed (787)")
	// names none of the three, so these three assertions are what separate a
	// guard from an accident.
	for _, want := range []string{
		"foreign_keys",
		"tool_call.task_id REFERENCES task_log(id)",
		"Nothing was modified.",
	} {
		if !strings.Contains(msg, want) {
			t.Errorf("refusal does not mention %q; it is not the guard's message:\n%v", want, msg)
		}
	}
	// The refusal text is recorded so the evidence file can quote the database's
	// own words rather than paraphrase them.
	t.Logf("refusal (verbatim): %v", msg)

	// A refusal must leave the database exactly as it was: watermark still 3,
	// both tables intact, no half-built table behind.
	if err := dbOn.QueryRowContext(ctx,
		`SELECT value FROM schema_meta WHERE key='schema_version'`).Scan(&offVersion); err != nil {
		t.Fatalf("read watermark after refusal: %v", err)
	}
	if offVersion != 3 {
		t.Errorf("schema_version = %d after a refused step, want 3", offVersion)
	}
	var tasks, children, debris int
	if err := dbOn.QueryRowContext(ctx, `SELECT count(*) FROM task_log`).Scan(&tasks); err != nil {
		t.Fatalf("count task_log: %v", err)
	}
	if err := dbOn.QueryRowContext(ctx, `SELECT count(*) FROM tool_call`).Scan(&children); err != nil {
		t.Fatalf("count tool_call: %v", err)
	}
	if err := dbOn.QueryRowContext(ctx,
		`SELECT count(*) FROM sqlite_master WHERE name LIKE 'task_log_v4%'`).Scan(&debris); err != nil {
		t.Fatalf("count debris: %v", err)
	}
	if tasks != len(theFourLegacyWords) || children != 1 || debris != 0 {
		t.Errorf("after a refused step: task_log = %d, tool_call = %d, task_log_v4 objects = %d; "+
			"want %d, 1 and 0", tasks, children, debris, len(theFourLegacyWords))
	}
}
