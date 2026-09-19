package memory

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// openTestStore opens a Store in a fresh t.TempDir subdirectory.
func openTestStore(t *testing.T, opts ...Option) *Store {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "data")
	s, err := Open(dir, opts...)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

// -----------------------------------------------------------------------------
// Contract: sqlite_master introspection vs SPEC-02 §3 (D35, eight tables).
// -----------------------------------------------------------------------------

type colSpec struct {
	name     string
	dataType string
	notNull  bool
	pk       bool
}

type tableSpec struct {
	cols    []colSpec
	indexes []string
}

// contractSpec mirrors SPEC-02 §3 exactly: eight tables, their columns in
// order, and the indexes the spec declares. Any drift between schema.go and
// the spec fails here.
//
// notNull mirrors what PRAGMA table_info reports for THIS DDL: PRIMARY KEY
// columns report notnull=0 in SQLite (rowid-alias and TEXT PKs alike) even
// though the key can never be NULL in practice — the pk flag is the binding
// part of the contract for those columns.
var contractSpec = map[string]tableSpec{
	"schema_meta": {
		cols: []colSpec{
			{"key", "TEXT", false, true},
			{"value", "TEXT", true, false},
		},
	},
	"profile": {
		cols: []colSpec{
			{"id", "INTEGER", false, true},
			{"slot", "TEXT", true, false},
			{"content", "TEXT", true, false},
			{"source", "TEXT", true, false},
			{"updated_at", "INTEGER", true, false},
		},
		indexes: []string{"idx_profile_updated"},
	},
	"memory": {
		cols: []colSpec{
			{"id", "INTEGER", false, true},
			{"kind", "TEXT", true, false},
			{"content", "TEXT", true, false},
			{"keywords", "TEXT", true, false},
			{"created_at", "INTEGER", true, false},
			{"last_hit_at", "INTEGER", true, false},
			{"hit_count", "INTEGER", true, false},
		},
		indexes: []string{"idx_memory_keywords", "idx_memory_last_hit"},
	},
	"task_log": {
		cols: []colSpec{
			{"id", "TEXT", false, true},
			{"started_at", "INTEGER", true, false},
			{"ended_at", "INTEGER", false, false},
			{"state", "TEXT", true, false},
			{"query_text", "TEXT", true, false},
			{"summary_text", "TEXT", false, false},
			{"cost_tokens_in", "INTEGER", true, false},
			{"cost_tokens_out", "INTEGER", true, false},
			{"cost_amount_micro", "INTEGER", true, false},
			{"currency", "TEXT", true, false},
			{"error_class", "TEXT", false, false},
		},
		indexes: []string{"idx_task_log_started"},
	},
	"tool_call": {
		cols: []colSpec{
			{"id", "INTEGER", false, true},
			{"task_id", "TEXT", true, false},
			{"seq", "INTEGER", true, false},
			{"tool", "TEXT", true, false},
			{"args_json", "TEXT", true, false},
			{"risk_level", "TEXT", true, false},
			{"decision", "TEXT", false, false},
			{"decided_at", "INTEGER", false, false},
			{"started_at", "INTEGER", false, false},
			{"ended_at", "INTEGER", false, false},
			{"outcome", "TEXT", false, false},
			{"error_class", "TEXT", false, false},
			{"correlation_id", "TEXT", true, false},
			{"grant_id", "INTEGER", false, false},
		},
		indexes: []string{"idx_tool_call_task", "idx_tool_call_corr"},
	},
	"approval_grant": {
		cols: []colSpec{
			{"id", "INTEGER", false, true},
			{"scope", "TEXT", true, false},
			{"tool", "TEXT", true, false},
			{"pattern", "TEXT", true, false},
			{"session_id", "TEXT", true, false},
			{"created_at", "INTEGER", true, false},
			{"expires_at", "INTEGER", true, false},
			{"revoked_at", "INTEGER", false, false},
		},
		indexes: []string{"idx_grant_session"},
	},
	"cost_daily": {
		cols: []colSpec{
			{"day", "TEXT", false, true},
			{"tokens_in", "INTEGER", true, false},
			{"tokens_out", "INTEGER", true, false},
			{"cost_micros", "INTEGER", true, false},
			{"tasks", "INTEGER", true, false},
		},
	},
	"plugin_state": {
		cols: []colSpec{
			{"id", "TEXT", false, true},
			{"version", "TEXT", true, false},
			{"enabled", "INTEGER", true, false},
			{"capabilities_json", "TEXT", true, false},
			{"net_allowlist_json", "TEXT", false, false},
			{"installed_at", "INTEGER", true, false},
			{"hash", "TEXT", true, false},
			{"exe_hash", "TEXT", false, false},
		},
	},
}

var contractTableOrder = []string{
	"schema_meta", "profile", "memory", "task_log", "tool_call",
	"approval_grant", "cost_daily", "plugin_state",
}

func TestSchemaContractIntrospection(t *testing.T) {
	s := openTestStore(t)
	db := s.reader

	// 1. Exactly the eight tables exist (plus none other).
	rows, err := db.QueryContext(context.Background(),
		`SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name`)
	if err != nil {
		t.Fatalf("query tables: %v", err)
	}
	defer rows.Close()
	var gotTables []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			t.Fatal(err)
		}
		gotTables = append(gotTables, n)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if len(gotTables) != len(contractTableOrder) {
		t.Errorf("table count = %d, want %d: %v", len(gotTables), len(contractTableOrder), gotTables)
	}
	for _, want := range contractTableOrder {
		found := false
		for _, got := range gotTables {
			if got == want {
				found = true
			}
		}
		if !found {
			t.Errorf("missing table %q (got %v)", want, gotTables)
		}
	}

	// 2. Per-table columns: name, type, notnull, pk — order and values.
	for table, spec := range contractSpec {
		t.Run("columns/"+table, func(t *testing.T) { checkColumns(t, db, table, spec) })
		t.Run("indexes/"+table, func(t *testing.T) { checkIndexes(t, db, table, spec) })
	}

	// 3. schema_version = 1 seeded.
	var v string
	if err := db.QueryRowContext(context.Background(),
		`SELECT value FROM schema_meta WHERE key='schema_version'`).Scan(&v); err != nil {
		t.Fatalf("read schema_version: %v", err)
	}
	if v != "1" {
		t.Errorf("schema_version = %q, want %q", v, "1")
	}

	// 4. The DDL script round-trips through SQLite (statement splitting is
	// sound: no semicolons hidden in literals/comments).
	// 8 CREATE TABLE + 7 CREATE INDEX = exactly 15 statements in the spec DDL.
	if n := len(splitSQLStatements(ddlV1)); n != 15 {
		t.Errorf("splitSQLStatements produced %d statements, want 15 (8 tables + 7 indexes)", n)
	}
}

func checkColumns(t *testing.T, db *sql.DB, table string, spec tableSpec) {
	t.Helper()
	rows, err := db.QueryContext(context.Background(), "PRAGMA table_info("+table+")")
	if err != nil {
		t.Fatalf("table_info(%s): %v", table, err)
	}
	defer rows.Close()
	var got []colSpec
	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull, pk int
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &dataType, &notNull, &dflt, &pk); err != nil {
			t.Fatal(err)
		}
		got = append(got, colSpec{name: name, dataType: dataType, notNull: notNull == 1, pk: pk > 0})
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if len(got) != len(spec.cols) {
		t.Fatalf("column count = %d, want %d\ngot:  %v\nwant: %v", len(got), len(spec.cols), colNames(got), colNames(spec.cols))
	}
	for i, w := range spec.cols {
		g := got[i]
		if g.name != w.name {
			t.Errorf("col[%d] name = %q, want %q", i, g.name, w.name)
		}
		if g.dataType != w.dataType {
			t.Errorf("%s.%s type = %q, want %q", table, g.name, g.dataType, w.dataType)
		}
		if g.notNull != w.notNull {
			t.Errorf("%s.%s notnull = %v, want %v", table, g.name, g.notNull, w.notNull)
		}
		if g.pk != w.pk {
			t.Errorf("%s.%s pk = %v, want %v", table, g.name, g.pk, w.pk)
		}
	}
}

func checkIndexes(t *testing.T, db *sql.DB, table string, spec tableSpec) {
	t.Helper()
	rows, err := db.QueryContext(context.Background(),
		`SELECT name FROM sqlite_master WHERE type='index' AND tbl_name=? AND name NOT LIKE 'sqlite_%'`, table)
	if err != nil {
		t.Fatalf("indexes(%s): %v", table, err)
	}
	defer rows.Close()
	var got []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			t.Fatal(err)
		}
		got = append(got, n)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if len(got) != len(spec.indexes) {
		t.Fatalf("%s index count = %d (%v), want %d (%v)", table, len(got), got, len(spec.indexes), spec.indexes)
	}
	gotSet := map[string]bool{}
	for _, g := range got {
		gotSet[g] = true
	}
	for _, w := range spec.indexes {
		if !gotSet[w] {
			t.Errorf("%s missing index %q", table, w)
		}
	}
}

func colNames(cols []colSpec) []string {
	out := make([]string, len(cols))
	for i, c := range cols {
		out[i] = c.name
	}
	return out
}

// -----------------------------------------------------------------------------
// Migration chain (SPEC-02 §7).
// -----------------------------------------------------------------------------

func TestMigrateFreshDatabaseCreatesV1(t *testing.T) {
	s := openTestStore(t)
	// Contract test above covers content; here: reopen is a no-op (no second
	// backup, no error) and version stays 1.
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	s2, err := Open(s.dir)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer s2.Close()
	var v string
	if err := s2.reader.QueryRowContext(context.Background(),
		`SELECT value FROM schema_meta WHERE key='schema_version'`).Scan(&v); err != nil {
		t.Fatal(err)
	}
	if v != "1" {
		t.Errorf("schema_version after reopen = %q", v)
	}
	// Reopen must not have written another backup: the bak-0-1 file exists
	// exactly once from the first open.
	entries, err := os.ReadDir(s.backupDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "wisp.db.bak-0-1" {
		t.Errorf("backup dir after reopen = %v, want [wisp.db.bak-0-1]", namesOf(entries))
	}
}

func namesOf(entries []os.DirEntry) []string {
	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = e.Name()
	}
	return out
}

// withTestMigration registers a test-only v1->v2 step so the chain mechanics
// (backup name, single transaction, watermark update) are exercised without
// polluting the production chain. It mutates the package chain and restores
// it at cleanup.
func withTestMigration(t *testing.T, next func(tx txExec) error) Option {
	t.Helper()
	saved := migrationChain
	t.Cleanup(func() { migrationChain = saved })
	return func(o *storeOptions) {
		migrationChain = append(saved, migrationStep{from: 1, to: 2, apply: next})
	}
}

func TestMigrationChainWithBackup(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "data")

	// 1. Create a v1 database with data in it.
	s1, err := Open(dir)
	if err != nil {
		t.Fatalf("open v1: %v", err)
	}
	ctx := context.Background()
	err = s1.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO profile(slot, content, source, updated_at) VALUES('pref.language','zh-CN','manual',?)`,
			time.Now().UTC().Unix())
		return err
	})
	if err != nil {
		t.Fatalf("seed v1: %v", err)
	}
	if err := s1.Close(); err != nil {
		t.Fatal(err)
	}

	// 2. Reopen targeting v2 via a test-only migration that adds a table in
	// its own transaction.
	mkV2 := withTestMigration(t, func(tx txExec) error {
		_, err := tx.ExecContext(context.Background(),
			`CREATE TABLE schema_probe_v2(x INTEGER NOT NULL)`)
		return err
	})
	s2, err := Open(dir, mkV2, withSchemaTarget(2))
	if err != nil {
		t.Fatalf("open v2: %v", err)
	}
	defer s2.Close()

	// 3. Backup written before the step, named per SPEC-02 §7, and it still
	// holds the pre-migration content (v1 backup = full db copy).
	bakPath := filepath.Join(dir, "backup", "wisp.db.bak-1-2")
	fi, err := os.Stat(bakPath)
	if err != nil {
		t.Fatalf("pre-migration backup missing: %v", err)
	}
	if fi.Size() == 0 {
		t.Error("pre-migration backup is empty")
	}

	// 4. Watermark advanced inside the same transaction chain.
	var v int
	if err := s2.reader.QueryRowContext(ctx, `SELECT CAST(value AS INTEGER) FROM schema_meta WHERE key='schema_version'`).Scan(&v); err != nil {
		t.Fatal(err)
	}
	if v != 2 {
		t.Errorf("schema_version = %d, want 2", v)
	}
	// Pre-existing data survived the migration.
	var n int
	if err := s2.reader.QueryRowContext(ctx, `SELECT count(*) FROM profile`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Errorf("profile rows after migration = %d, want 1", n)
	}
	// The step's own DDL landed.
	var probe int
	if err := s2.reader.QueryRowContext(ctx,
		`SELECT count(*) FROM sqlite_master WHERE type='table' AND name='schema_probe_v2'`).Scan(&probe); err != nil {
		t.Fatal(err)
	}
	if probe != 1 {
		t.Error("test-only v2 table missing after migration")
	}
}

// TestMigrationFailedStepIsAtomic proves each step is a single transaction
// (SPEC-02 §7): a step that fails mid-way must leave the watermark and the
// schema untouched, with the pre-backup in place for manual recovery.
func TestMigrationFailedStepIsAtomic(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "data")
	s1, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := s1.Close(); err != nil {
		t.Fatal(err)
	}

	failStep := withTestMigration(t, func(tx txExec) error {
		if _, err := tx.ExecContext(context.Background(),
			`CREATE TABLE schema_probe_v2(x INTEGER NOT NULL)`); err != nil {
			return err
		}
		// Violate the NOT NULL inside the same transaction: the CREATE above
		// must roll back together with the watermark update.
		_, err := tx.ExecContext(context.Background(),
			`INSERT INTO schema_probe_v2(x) VALUES(NULL)`)
		return err
	})
	_, err = Open(dir, failStep, withSchemaTarget(2))
	if err == nil {
		t.Fatal("Open succeeded despite failing migration step")
	}

	// Database is still at v1, probe table rolled back.
	db, err := sql.Open(DriverName, dsn(filepath.Join(dir, dbFileName), true))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var v int
	if err := db.QueryRow(`SELECT CAST(value AS INTEGER) FROM schema_meta WHERE key='schema_version'`).Scan(&v); err != nil {
		t.Fatal(err)
	}
	if v != 1 {
		t.Errorf("schema_version after failed step = %d, want 1 (atomic rollback)", v)
	}
	var probe int
	if err := db.QueryRow(`SELECT count(*) FROM sqlite_master WHERE name='schema_probe_v2'`).Scan(&probe); err != nil {
		t.Fatal(err)
	}
	if probe != 0 {
		t.Error("probe table survived a failed migration step (not atomic)")
	}
	// Pre-backup exists.
	if _, err := os.Stat(filepath.Join(dir, "backup", "wisp.db.bak-1-2")); err != nil {
		t.Errorf("pre-migration backup missing: %v", err)
	}
}

func TestMigrationNewerSchemaIsUnmigratable(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "data")
	s, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}

	// Forge a "future" database written by a newer binary.
	db, err := sql.Open(DriverName, dsn(s.dbPath, true))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`UPDATE schema_meta SET value='99' WHERE key='schema_version'`); err != nil {
		t.Fatal(err)
	}
	// Checkpoint so the forged value lives in the main db file (the tamper
	// must survive even if the -wal sidecar were removed).
	if _, err := db.Exec(startupCheckpointQuery); err != nil {
		t.Fatal(err)
	}
	db.Close()
	_ = os.Remove(s.dbPath + "-wal")
	_ = os.Remove(s.dbPath + "-shm")
	before, err := os.ReadFile(s.dbPath)
	if err != nil {
		t.Fatal(err)
	}

	s2, err := Open(dir)
	if err == nil {
		_ = s2.Close()
		t.Fatal("Open succeeded on a newer schema; want ErrSchemaUnmigratable")
	}
	if !IsSchemaUnmigratable(err) {
		t.Fatalf("error = %v, want errors.Is ErrSchemaUnmigratable", err)
	}
	var se *SchemaError
	if !errors.As(err, &se) || se.CurrentVersion != 99 {
		t.Errorf("SchemaError detail = %+v, want CurrentVersion=99", se)
	}

	// The database must be untouched (never silently overwritten).
	after, err := os.ReadFile(s.dbPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Error("database file changed after failed migration (must never be modified)")
	}
}

func TestMigrationForeignDatabaseIsUnmigratable(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "data")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open(DriverName, dsn(filepath.Join(dir, dbFileName), true))
	if err != nil {
		t.Fatal(err)
	}
	// A database with tables but no schema_meta watermark.
	if _, err := db.Exec(`CREATE TABLE something_else(a TEXT NOT NULL)`); err != nil {
		t.Fatal(err)
	}
	db.Close()

	_, err = Open(dir)
	if !IsSchemaUnmigratable(err) {
		t.Fatalf("error = %v, want ErrSchemaUnmigratable", err)
	}
	if !strings.Contains(err.Error(), "no schema_meta") {
		t.Errorf("error should name the missing watermark: %v", err)
	}
}

// TestParseVersionStrict is the MINOR-5 regression: the watermark parser must
// reject anything that is not a plain decimal (lenient prefix parsing would
// read "1abc" as 1 and migrate a foreign database).
func TestParseVersionStrict(t *testing.T) {
	for _, ok := range []string{"0", "1", "42", " 3 "} {
		v, err := parseVersion(ok)
		if err != nil {
			t.Errorf("parseVersion(%q) = %v, want accepted", ok, err)
		} else if v != 3 && ok == " 3 " {
			t.Errorf("parseVersion(%q) = %d", ok, v)
		}
	}
	// "1\n" is accepted BY DESIGN: the value is whitespace-trimmed before the
	// strict digit check (a stray newline from manual editing is tolerated;
	// embedded garbage like "1abc" is not).
	for _, bad := range []string{"", "   ", "1abc", "abc", "1.5", "-1", "+2", "1e2", "0x1"} {
		if _, err := parseVersion(bad); err == nil {
			t.Errorf("parseVersion(%q) accepted, want rejection", bad)
		}
	}
}
