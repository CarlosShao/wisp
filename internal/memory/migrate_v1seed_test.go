package memory

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// A REAL v1 database (schema_version=1, with seed data written through the
// normal write path) must migrate to v2 on reopen: provider_health appears,
// the seed data survives, and the pre-migration backup exists. This is the
// upgrade path every existing Wisp installation takes when the binary moves
// from the ticket-04 schema to the ticket-09 schema.

func TestMigrateSeededV1DatabaseToV2(t *testing.T) {
	dir := filepath.Join(sealableTempDir124(t), "data")

	// 1. Create a genuine v1 database via the production chain pinned to
	// target 1 (what the ticket-04 binary produced), then seed rows through
	// the regular DAO path.
	old, err := Open(dir, withSchemaTarget(1))
	if err != nil {
		t.Fatalf("open v1: %v", err)
	}
	ctx := context.Background()
	err = old.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO profile(slot, content, source, updated_at) VALUES('pref.language','zh-CN','manual',?)`,
			time.Now().UTC().Unix()); err != nil {
			return err
		}
		_, err := tx.ExecContext(ctx,
			`INSERT INTO memory(kind, content, keywords, created_at) VALUES('explicit','记得带伞','带 伞',?)`,
			time.Now().UTC().Unix())
		return err
	})
	if err != nil {
		t.Fatalf("seed v1: %v", err)
	}
	// A v1 database has no provider_health table yet.
	var hasPH int
	if err := old.reader.QueryRowContext(ctx,
		`SELECT count(*) FROM sqlite_master WHERE type='table' AND name='provider_health'`).Scan(&hasPH); err != nil {
		t.Fatal(err)
	}
	if hasPH != 0 {
		t.Fatal("v1 target unexpectedly created provider_health (test premise broken)")
	}
	if err := old.Close(); err != nil {
		t.Fatal(err)
	}

	// 2. Reopen at the production target (2): the 1->2 step must run.
	fresh, err := Open(dir)
	if err != nil {
		t.Fatalf("reopen (migrate v1->v2): %v", err)
	}
	defer fresh.Close()

	var v string
	if err := fresh.reader.QueryRowContext(ctx,
		`SELECT value FROM schema_meta WHERE key='schema_version'`).Scan(&v); err != nil {
		t.Fatal(err)
	}
	if v != "2" {
		t.Errorf("schema_version = %q, want 2", v)
	}

	// 3. provider_health now exists with the v2 shape and is usable.
	if err := fresh.UpsertProviderProbe(ctx, "openai", "gpt-4o", ProbeFlags{}, true, 500, time.Now().UTC()); err != nil {
		t.Fatalf("provider_health write after migration: %v", err)
	}
	row, err := fresh.ProviderHealthRow(ctx, "openai", "gpt-4o")
	if err != nil {
		t.Fatalf("provider_health read after migration: %v", err)
	}
	if row.Provider != "openai" || row.Model != "gpt-4o" {
		t.Errorf("row = %+v", row)
	}

	// 4. Seed data survived the migration.
	var profiles int
	if err := fresh.reader.QueryRowContext(ctx, `SELECT count(*) FROM profile`).Scan(&profiles); err != nil {
		t.Fatal(err)
	}
	if profiles != 1 {
		t.Errorf("profile rows after migration = %d, want 1", profiles)
	}
	var memories int
	if err := fresh.reader.QueryRowContext(ctx, `SELECT count(*) FROM memory`).Scan(&memories); err != nil {
		t.Fatal(err)
	}
	if memories != 1 {
		t.Errorf("memory rows after migration = %d, want 1", memories)
	}

	// 5. The pre-migration backup for the 1->2 step exists (SPEC-02 §7);
	// the dir also holds the 0->1 backup from CREATING the v1 database, and
	// nothing else.
	if _, err := os.Stat(filepath.Join(dir, "backup", "wisp.db.bak-1-2")); err != nil {
		t.Errorf("pre-migration backup missing: %v", err)
	}
	entries, err := os.ReadDir(filepath.Join(dir, "backup"))
	if err != nil {
		t.Fatal(err)
	}
	wantBackups := []string{"wisp.db.bak-0-1", "wisp.db.bak-1-2"}
	if len(entries) != len(wantBackups) {
		t.Errorf("backup dir = %v, want %v", namesOf(entries), wantBackups)
	}
	for i, w := range wantBackups {
		if i < len(entries) && entries[i].Name() != w {
			t.Errorf("backup[%d] = %s, want %s", i, entries[i].Name(), w)
		}
	}
}
