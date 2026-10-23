package memory

// Ticket 366 leg r1: the 3->4 rebuild now runs SQLite's 12-step procedure step 10
// (PRAGMA foreign_key_check) inside its transaction and RECORDS the result.
//
// Why the record is not a gate, in one sentence: refusing on a non-zero count
// rolls the rebuild back, so schema_version stays at 3, and every later start takes
// the same path and refuses again - the database is locked out forever, not once
// (probes/orch/20261022 census 366-a1 section 5 measured exactly that). The cases
// below pin the record, the four reachable orphan shapes 366-a1 found, and the
// property that separates this shape from the hard gate: opening twice.
//
// The counts 3/2/1/1 are not copied from the census; they are re-measured here on
// fixtures this file builds through the production chain (Open with a v3 target),
// and each case also checks them against a second, independent ruler - a plain
// PRAGMA foreign_key_check on a business-shaped connection after the migration.

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// t366FkRows reads PRAGMA foreign_key_check on a business-shaped connection
// (dsn(), foreign_keys left alone). The check does not need enforcement on.
func t366FkRows(t *testing.T, dbPath string) []foreignKeyCheckRow {
	t.Helper()
	db, err := sql.Open(DriverName, dsn(dbPath, false))
	if err != nil {
		t.Fatalf("open for the independent foreign_key_check: %v", err)
	}
	defer db.Close()
	rows, err := db.QueryContext(context.Background(), pragmaForeignKeyCheckQuery)
	if err != nil {
		t.Fatalf("%s on %s: %v", pragmaForeignKeyCheckQuery, dbPath, err)
	}
	defer rows.Close()
	var out []foreignKeyCheckRow
	for rows.Next() {
		var r foreignKeyCheckRow
		var rowID sql.NullInt64
		if err := rows.Scan(&r.Child, &rowID, &r.Parent, &r.FKID); err != nil {
			t.Fatalf("scan independent foreign_key_check: %v", err)
		}
		r.RowID = rowID.Int64
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("drain independent foreign_key_check: %v", err)
	}
	return out
}

// t366Scalar reads one integer out of dbPath (counts and the watermark).
func t366Scalar(t *testing.T, dbPath, query string) int {
	t.Helper()
	db, err := sql.Open(DriverName, dsn(dbPath, false))
	if err != nil {
		t.Fatalf("open for %q: %v", query, err)
	}
	defer db.Close()
	var n int
	if err := db.QueryRowContext(context.Background(), query).Scan(&n); err != nil {
		t.Fatalf("%q: %v", query, err)
	}
	return n
}

// t366SeedV3 builds a real v3 database through the production migration chain and
// seeds the shape 366-a1 measured: task-A with two tool_call rows and task-B with
// one. The store is returned open at v3 so the caller can drive a production
// orphan path and then close it; the migration under test is the next Open.
func t366SeedV3(t *testing.T) (string, *Store) {
	t.Helper()
	dir := filepath.Join(sealableTempDir124(t), "data")
	s, err := Open(dir, withSchemaTarget(3))
	if err != nil {
		t.Fatalf("seed: open a v3 database: %v", err)
	}
	ctx := context.Background()
	now := time.Now().Unix()
	for _, id := range []string{"task-A", "task-B"} {
		if err := s.StartTaskLog(ctx, TaskLog{ID: id, StartedAt: now - 60, State: "done", QueryText: "q-" + id}); err != nil {
			t.Fatalf("seed task_log %s: %v", id, err)
		}
	}
	for _, tc := range []ToolCall{
		{TaskID: "task-A", Seq: 1, Tool: "read_file", ArgsJSON: `{}`, RiskLevel: RiskL0, CorrelationID: "c1"},
		{TaskID: "task-A", Seq: 2, Tool: "read_file", ArgsJSON: `{}`, RiskLevel: RiskL0, CorrelationID: "c2"},
		{TaskID: "task-B", Seq: 1, Tool: "read_file", ArgsJSON: `{}`, RiskLevel: RiskL0, CorrelationID: "c3"},
	} {
		if _, err := s.InsertToolCall(ctx, tc); err != nil {
			t.Fatalf("seed tool_call for %s: %v", tc.TaskID, err)
		}
	}
	return dir, s
}

// t366OrphanFixture builds a v3 database carrying orphan tool_call rows through
// the production path that reaches them, and returns the directory with the store
// closed - ready for the migration under test.
func t366OrphanFixture(t *testing.T, path string) string {
	t.Helper()
	dir, s := t366SeedV3(t)
	ctx := context.Background()
	switch path {
	case "purge-tasklogs": // privacy one-click clear: parents go, children stay
		n, err := s.PurgeTaskLogs(ctx)
		if err != nil || n != 2 {
			t.Fatalf("purge: rows=%d err=%v, want 2 and nil", n, err)
		}
	case "delete-tasklog": // one task removed, its children stay
		if err := s.DeleteTaskLog(ctx, "task-A"); err != nil {
			t.Fatalf("delete task-A: %v", err)
		}
	case "manual-insert": // a writer that bypasses the DAO
		db, err := sql.Open(DriverName, dsn(filepath.Join(dir, dbFileName), false))
		if err != nil {
			t.Fatalf("manual insert: open: %v", err)
		}
		now := time.Now().Unix()
		_, ierr := db.ExecContext(ctx,
			`INSERT INTO tool_call(task_id, seq, tool, args_json, risk_level, decision, decided_at,
			   started_at, ended_at, outcome, error_class, correlation_id, grant_id)
			 VALUES('ghost-missing-task', 9, 'read_file', '{}', 'L0', 'allow', ?, ?, ?, 'success', NULL, 'c9', NULL)`,
			now, now, now)
		db.Close()
		if ierr != nil {
			t.Fatalf("manual insert: %v", ierr)
		}
	case "retention": // a NULL-timestamp child outlives its parent (journal.startCall)
		now := time.Now()
		old := now.Add(-31 * 24 * time.Hour).Unix()
		if err := s.StartTaskLog(ctx, TaskLog{ID: "task-old", StartedAt: old, State: "done", QueryText: "q-old"}); err != nil {
			t.Fatalf("retention: seed task-old: %v", err)
		}
		if _, err := s.InsertToolCall(ctx, ToolCall{
			TaskID: "task-old", Seq: 1, Tool: "read_file",
			ArgsJSON: `{}`, RiskLevel: RiskL0, CorrelationID: "c-old-null",
		}); err != nil {
			t.Fatalf("retention: seed the NULL-timestamp child: %v", err)
		}
		oldTS := old
		if _, err := s.InsertToolCall(ctx, ToolCall{
			TaskID: "task-old", Seq: 2, Tool: "read_file",
			ArgsJSON: `{}`, RiskLevel: RiskL0, CorrelationID: "c-old-ts",
			StartedAt: &oldTS, EndedAt: &oldTS,
		}); err != nil {
			t.Fatalf("retention: seed the old-timestamp child: %v", err)
		}
		res := s.RunRetentionOnce(ctx, RetentionConfig{Now: func() time.Time { return now }})
		if res.TaskLogDeleted != 1 || res.ToolCallDeleted != 1 {
			t.Fatalf("retention: deleted task_log=%d tool_call=%d, want 1 and 1",
				res.TaskLogDeleted, res.ToolCallDeleted)
		}
	default:
		t.Fatalf("unknown orphan path %q", path)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("close the v3 fixture: %v", err)
	}
	return dir
}

// t366ChildTaskIDs lists tool_call.task_id in row order: the fixture's forensic
// rows are what must survive the migration untouched.
func t366ChildTaskIDs(t *testing.T, dbPath string) []string {
	t.Helper()
	db, err := sql.Open(DriverName, dsn(dbPath, false))
	if err != nil {
		t.Fatalf("open for child ids: %v", err)
	}
	defer db.Close()
	rows, err := db.QueryContext(context.Background(), `SELECT task_id FROM tool_call ORDER BY id`)
	if err != nil {
		t.Fatalf("read child ids: %v", err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			t.Fatalf("scan child id: %v", err)
		}
		out = append(out, id)
	}
	return out
}

// TestMigration366AuditRecordsCleanDatabase is the positive control: the ruler is
// not stuck above zero. Both paths that reach a clean v4 must report 0 - a brand
// new file (0->1->2->3->4) and an existing clean v3 database (3->4).
func TestMigration366AuditRecordsCleanDatabase(t *testing.T) {
	t.Run("fresh database", func(t *testing.T) {
		dir := filepath.Join(sealableTempDir124(t), "data")
		s, err := Open(dir)
		if err != nil {
			t.Fatalf("Open on a fresh directory: %v", err)
		}
		defer s.Close()
		audit := s.ForeignKeyAudit()
		if !audit.Checked {
			t.Fatal("a fresh database ran 0->4 and the 3->4 step is audited, so Checked must be true: " +
				"the check never ran, or the step lost its auditForeignKeys flag")
		}
		if audit.SchemaVersion != SchemaVersionTarget || audit.OrphanRows != 0 || len(audit.ChildTables) != 0 {
			t.Errorf("fresh audit = {version %d, orphans %d, tables %v}, want {%d, 0, []}",
				audit.SchemaVersion, audit.OrphanRows, audit.ChildTables, SchemaVersionTarget)
		}
		dbPath := filepath.Join(dir, dbFileName)
		if n := len(t366FkRows(t, dbPath)); n != 0 {
			t.Errorf("independent foreign_key_check on the fresh database = %d rows, want 0", n)
		}
	})

	t.Run("clean v3 database", func(t *testing.T) {
		dir, s := t366SeedV3(t)
		if err := s.Close(); err != nil {
			t.Fatalf("close the v3 fixture: %v", err)
		}
		s2, err := Open(dir)
		if err != nil {
			t.Fatalf("Open on a clean v3 database: %v", err)
		}
		defer s2.Close()
		audit := s2.ForeignKeyAudit()
		if !audit.Checked || audit.OrphanRows != 0 || audit.SchemaVersion != SchemaVersionTarget {
			t.Errorf("clean audit = {checked %v, version %d, orphans %d}, want {true, %d, 0}: "+
				"a database with no dangling references must read 0, not a constant",
				audit.Checked, audit.SchemaVersion, audit.OrphanRows, SchemaVersionTarget)
		}
		dbPath := filepath.Join(dir, dbFileName)
		if n := len(t366FkRows(t, dbPath)); n != 0 {
			t.Errorf("independent foreign_key_check = %d rows on the clean fixture, want 0", n)
		}
		if got := t366Scalar(t, dbPath, `SELECT count(*) FROM task_log`); got != 2 {
			t.Errorf("task_log = %d after the migration, want 2 (the rebuild copies rows, it does not drop them)", got)
		}
		if got := strings.Join(t366ChildTaskIDs(t, dbPath), ","); got != "task-A,task-A,task-B" {
			t.Errorf("tool_call.task_id = %q after the migration, want task-A,task-A,task-B", got)
		}
	})
}

// TestMigration366AuditRecordsOrphansAndStillOpens is the subject half: all four
// orphan shapes 366-a1 found reach v4, the store opens, the count is the truth
// (3/2/1/1), and no row is deleted to make the number smaller.
func TestMigration366AuditRecordsOrphansAndStillOpens(t *testing.T) {
	cases := []struct {
		path        string
		orphans     int
		taskLogs    int
		toolCalls   int
		childTaskID string
	}{
		{"purge-tasklogs", 3, 0, 3, "task-A,task-A,task-B"},
		{"delete-tasklog", 2, 1, 3, "task-A,task-A,task-B"},
		{"manual-insert", 1, 2, 4, "task-A,task-A,task-B,ghost-missing-task"},
		{"retention", 1, 2, 4, "task-A,task-A,task-B,task-old"},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			dir := t366OrphanFixture(t, tc.path)
			dbPath := filepath.Join(dir, dbFileName)
			s, err := Open(dir)
			if err != nil {
				t.Fatalf("Open refused a database with %d orphan row(s): %v - recording is not a gate "+
					"(refusing here leaves schema_version at 3 and refuses again on every start)",
					tc.orphans, err)
			}
			defer s.Close()
			audit := s.ForeignKeyAudit()
			if !audit.Checked {
				t.Fatal("the migration ran the 3->4 step but the audit says Checked=false: " +
					"the foreign_key_check hop is missing from applyMigrationStep")
			}
			if audit.SchemaVersion != SchemaVersionTarget {
				t.Errorf("audit.SchemaVersion = %d, want %d", audit.SchemaVersion, SchemaVersionTarget)
			}
			if audit.OrphanRows != tc.orphans {
				t.Errorf("audit.OrphanRows = %d, want %d (366-a1 measured %d for this path)",
					audit.OrphanRows, tc.orphans, tc.orphans)
			}
			if want := []string{"tool_call"}; strings.Join(audit.ChildTables, ",") != strings.Join(want, ",") {
				t.Errorf("audit.ChildTables = %v, want %v: the dangling rows live in the child table "+
					"whose parent was removed", audit.ChildTables, want)
			}
			// Independent ruler: the same question asked again on a fresh connection
			// after the migration. If these two disagree, the audit is not the truth.
			if n := len(t366FkRows(t, dbPath)); n != tc.orphans {
				t.Errorf("independent foreign_key_check = %d rows, want %d (audit said %d)",
					n, tc.orphans, audit.OrphanRows)
			}
			if got := t366Scalar(t, dbPath, `SELECT count(*) FROM task_log`); got != tc.taskLogs {
				t.Errorf("task_log = %d, want %d: the migration must neither create nor delete parents", got, tc.taskLogs)
			}
			if got := t366Scalar(t, dbPath, `SELECT count(*) FROM tool_call`); got != tc.toolCalls {
				t.Errorf("tool_call = %d, want %d: forensic rows are never deleted to satisfy the check", got, tc.toolCalls)
			}
			if got := strings.Join(t366ChildTaskIDs(t, dbPath), ","); got != tc.childTaskID {
				t.Errorf("tool_call.task_id = %q, want %q: the forensic rows must be untouched", got, tc.childTaskID)
			}
			if v := t366Scalar(t, dbPath, `SELECT value FROM schema_meta WHERE key='schema_version'`); v != 4 {
				t.Errorf("schema_version = %d, want 4", v)
			}
		})
	}
}

// TestMigration366AuditNeverLocksTheDatabase is the boundary that separates this
// shape from the hard gate 366-a1 measured: the same orphan database opens twice,
// both times with a nil error.
//
// Under ">0 => refuse" the first Open fails inside the rebuild's transaction, the
// watermark stays at 3, and the second Open takes the same path and fails again -
// forever. Here the first Open migrates and records, the second Open takes
// open.go's `version == s.target` shortcut and records nothing (Checked=false),
// and the orphan rows are still there for whoever wants them.
func TestMigration366AuditNeverLocksTheDatabase(t *testing.T) {
	dir := t366OrphanFixture(t, "purge-tasklogs")
	dbPath := filepath.Join(dir, dbFileName)

	first, err := Open(dir)
	if err != nil {
		t.Fatalf("first Open refused: %v", err)
	}
	firstAudit := first.ForeignKeyAudit()
	if !firstAudit.Checked || firstAudit.OrphanRows != 3 {
		t.Fatalf("first Open audit = {checked %v, orphans %d}, want {true, 3}",
			firstAudit.Checked, firstAudit.OrphanRows)
	}
	if err := first.Close(); err != nil {
		t.Fatalf("close after the first Open: %v", err)
	}

	second, err := Open(dir)
	if err != nil {
		t.Fatalf("second Open refused a database the first Open had already migrated: %v - "+
			"that is the permanent-lock shape (watermark at 3, refused on every start)", err)
	}
	defer second.Close()
	secondAudit := second.ForeignKeyAudit()
	if secondAudit.Checked {
		t.Errorf("second Open audit says Checked=true: no step ran (schema_version was already %d), "+
			"so the baseline must report itself as absent rather than as a fresh zero", SchemaVersionTarget)
	}
	if secondAudit.SchemaVersion != SchemaVersionTarget || secondAudit.OrphanRows != 0 {
		t.Errorf("second Open audit = {version %d, orphans %d}, want {%d, 0}",
			secondAudit.SchemaVersion, secondAudit.OrphanRows, SchemaVersionTarget)
	}
	if v := t366Scalar(t, dbPath, `SELECT value FROM schema_meta WHERE key='schema_version'`); v != 4 {
		t.Errorf("schema_version = %d after two Opens, want 4", v)
	}
	if got := strings.Join(t366ChildTaskIDs(t, dbPath), ","); got != "task-A,task-A,task-B" {
		t.Errorf("tool_call.task_id = %q after two Opens, want task-A,task-A,task-B: nothing may be deleted", got)
	}
	if n := len(t366FkRows(t, dbPath)); n != 3 {
		t.Errorf("independent foreign_key_check = %d after two Opens, want 3 (the rows are kept)", n)
	}
}

// TestMigration366AuditIsLogged pins the operator's half of the record: the count
// is in the log, in both shapes, with the child table and the statement that lists
// the rows. The Store accessor is the test's half; this is the one a person reads.
func TestMigration366AuditIsLogged(t *testing.T) {
	t.Run("clean", func(t *testing.T) {
		logger, buf := bufLogger(t)
		dir := filepath.Join(sealableTempDir124(t), "data")
		s, err := Open(dir, WithLogger(logger))
		if err != nil {
			t.Fatalf("Open: %v", err)
		}
		defer s.Close()
		logs := buf.String()
		if !strings.Contains(logs, "PRAGMA foreign_key_check found no dangling references") {
			t.Errorf("the clean baseline is not in the log; the operator sees no evidence the check ran:\n%s", logs)
		}
		if !strings.Contains(logs, "orphan_rows=0") {
			t.Errorf("the clean baseline does not carry orphan_rows=0:\n%s", logs)
		}
	})

	t.Run("orphans", func(t *testing.T) {
		logger, buf := bufLogger(t)
		dir := t366OrphanFixture(t, "delete-tasklog")
		s, err := Open(dir, WithLogger(logger))
		if err != nil {
			t.Fatalf("Open: %v", err)
		}
		defer s.Close()
		logs := buf.String()
		for _, want := range []string{
			"PRAGMA foreign_key_check found dangling references",
			"orphan_rows=2",
			"child_tables=tool_call",
			"opening anyway",
		} {
			if !strings.Contains(logs, want) {
				t.Errorf("migration log does not contain %q:\n%s", want, logs)
			}
		}
		if strings.Contains(logs, "level=ERROR") {
			t.Errorf("the baseline logged an error for a database that opened fine:\n%s", logs)
		}
	})
}

// TestMigration366AuditFlagIsOnTheRebuildStep keeps the opt-in honest: the audit
// belongs to the step that rebuilds a table, and the chain must not quietly grow a
// second audited step (the Store keeps one baseline) or lose this one.
func TestMigration366AuditFlagIsOnTheRebuildStep(t *testing.T) {
	var audited []string
	for _, step := range migrationChain {
		if step.auditForeignKeys {
			audited = append(audited, fmt.Sprintf("%d->%d", step.from, step.to))
		}
	}
	if got := strings.Join(audited, ","); got != "3->4" {
		t.Errorf("audited steps = %q, want exactly 3->4 (the v4 rebuild): the Store keeps one baseline, "+
			"and a step that does not rebuild a table has nothing for step 10 to check", got)
	}
}
