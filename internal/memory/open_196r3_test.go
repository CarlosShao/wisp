package memory

// Ticket 196 leg r3: the migration connection declares its own foreign-key
// switch instead of inheriting one from the shared DSN.
//
// Why this file exists: r2 made the v4 rebuild refuse to run with enforcement
// ON, and left the residual that Open() then cannot open such a database at all.
// The fix is SQLite's own 12-step step 1 - turn the switch off on the connection
// that migrates, before the transaction starts. The three tests below pin the
// three things that make it a fix rather than a hope:
//
//  1. the shared DSN still answers the foreign-key question with silence, so no
//     business connection is switched off (or on) by this package;
//  2. the migration pragma list carries exactly ONE foreign_keys entry and it
//     says OFF - not two, see TestMigrationPragmas196r3ExactlyOneForeignKeysOff
//     for why the count is load-bearing rather than stylistic;
//  3. on the wire, a migration connection really reads 0 while the same path
//     opened from a list that enables enforcement reads 1.

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDsn196r3BusinessPoolsKeepTheirOwnForeignKeys is hard boundary 1, made
// executable: the migration fix must not be a global switch.
//
// The first half freezes the shared DSN byte for byte, because the pragma set is
// contract-level (D35 rule 3) and because "unchanged" is a comparison, not a
// claim. The second half is the actual boundary: a caller who turns enforcement
// on gets it, on a connection this package built.
func TestDsn196r3BusinessPoolsKeepTheirOwnForeignKeys(t *testing.T) {
	const p = `C:\Users\x\wisp\wisp.db`
	for _, forWriter := range []bool{false, true} {
		want := "file:C:/Users/x/wisp/wisp.db?_pragma=busy_timeout(3000)" +
			"&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)" +
			"&_pragma=wal_autocheckpoint(1000)"
		if forWriter {
			want += "&_txlock=immediate"
		}
		if got := dsn(p, forWriter); got != want {
			t.Errorf("dsn(forWriter=%v) changed shape:\n got  %s\n want %s", forWriter, got, want)
		}
		if strings.Contains(dsn(p, forWriter), "foreign_keys") {
			t.Errorf("the shared DSN now answers the foreign-key question; business pools must keep "+
				"the caller's value:\n%s", dsn(p, forWriter))
		}
	}

	dir := filepath.Join(sealableTempDir124(t), "data")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	// The caller's own switch, added on top of this package's DSN exactly the
	// way a future dsn() change would add it.
	db, err := sql.Open(DriverName, dsn(filepath.Join(dir, dbFileName), false)+"&_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()
	var fk int
	if err := db.QueryRowContext(context.Background(), pragmaForeignKeysQuery).Scan(&fk); err != nil {
		t.Fatalf("read pragma: %v", err)
	}
	if fk != 1 {
		t.Errorf("a business connection built from dsn() reads foreign_keys = %d after the caller asked "+
			"for 1: the migration fix leaked into the shared DSN", fk)
	}
}

// TestMigrationPragmas196r3ExactlyOneForeignKeysOff is the tooth for the one
// decision this leg adds.
//
// The count of foreign_keys entries is load-bearing, not stylistic: modernc
// sorts the _pragma values and applies them in sorted order (sqlite.go in
// v1.59.0), so a duplicate key is resolved by the sort rather than by DSN order.
// "Append an override after dsn()" is therefore a lottery played on the
// spelling, measured over eight spellings of the shared entry and the one
// override this leg adds: foreign_keys(1) and foreign_keys = 1 happened to end
// up OFF, while foreign_keys(ON), (on), (true) and (yes) all ended up ON, and
// the numeric pair the dispatch suggested - (1) then (0) - ends up ON too,
// 200 of 200 fresh connections. Contradicting the list instead (drop, then
// append exactly one) is the only shape that does not depend on a sort.
func TestMigrationPragmas196r3ExactlyOneForeignKeysOff(t *testing.T) {
	base := dsnPragmas()
	cases := []struct {
		name   string
		shared []string
	}{
		{"shared list as shipped", base},
		{"shared list enforcing", append(append([]string{}, base...), "foreign_keys(1)")},
		{"enforcing first", append([]string{"foreign_keys(ON)"}, base...)},
		{"spaced spelling", append(append([]string{}, base...), "foreign_keys = ON")},
		{"only the switch", []string{"foreign_keys(1)"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := migrationPragmas(tc.shared)
			joined := strings.Join(got, "&")
			if n := strings.Count(joined, "foreign_keys"); n != 1 {
				t.Fatalf("migration pragmas carry %d foreign_keys entries, want exactly 1 (%v): "+
					"two entries resolve to ON on this driver", n, got)
			}
			last := got[len(got)-1]
			if last != pragmaForeignKeysOff {
				t.Errorf("last migration pragma is %q, want %q: the switch must be declared, not left "+
					"to SQLite's default", last, pragmaForeignKeysOff)
			}
			var kept, want []string
			for _, p := range got {
				if !isForeignKeysPragma(p) {
					kept = append(kept, p)
				}
			}
			for _, p := range tc.shared {
				if !isForeignKeysPragma(p) {
					want = append(want, p)
				}
			}
			if strings.Join(kept, "&") != strings.Join(want, "&") {
				t.Errorf("shared pragmas were not carried over unchanged and in order:\n got  %v\n want %v",
					kept, want)
			}
		})
	}
}

// TestMigrationDSN196r3SwitchesForeignKeysOffOnTheWire closes the loop the
// string tests cannot: the driver really applies the list, and the two lists
// really differ when a connection is opened.
//
// The twin connection is the future case built by hand - the same path, the same
// shared pragmas, plus enforcement - so a migration pragma transform that fails
// to remove it shows up as foreign_keys = 1 here rather than as a green string
// comparison.
func TestMigrationDSN196r3SwitchesForeignKeysOffOnTheWire(t *testing.T) {
	ctx := context.Background()
	dir := filepath.Join(sealableTempDir124(t), "data")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	dbPath := filepath.Join(dir, dbFileName)

	read := func(label, dsnStr string) int {
		t.Helper()
		db, err := sql.Open(DriverName, dsnStr)
		if err != nil {
			t.Fatalf("%s: open: %v", label, err)
		}
		defer db.Close()
		var fk int
		if err := db.QueryRowContext(ctx, pragmaForeignKeysQuery).Scan(&fk); err != nil {
			t.Fatalf("%s: read pragma: %v", label, err)
		}
		return fk
	}

	// Control: the same pragmas the migration connection starts from, with
	// enforcement on. If this does not read 1 the subject half measures nothing.
	twin := buildDSN(dbPath, true, append(append([]string{}, dsnPragmas()...), "foreign_keys(1)"))
	if fk := read("twin", twin); fk != 1 {
		t.Fatalf("control: the enforcing twin connection reads foreign_keys = %d, want 1", fk)
	}
	if fk := read("migration", migrationDSN(dbPath)); fk != 0 {
		t.Errorf("the migration connection reads foreign_keys = %d, want 0: the v4 rebuild's DROP TABLE "+
			"would abort on tool_call's REFERENCES clause", fk)
	}

	// Every spelling a shared list could carry, on the wire. This is where the
	// "append an override" shape loses: it keeps the shared entry, and the
	// driver's sort then decides - for ON/on/true/yes it decides ON.
	for _, enforcing := range []string{
		"foreign_keys(1)", "foreign_keys(0)", "foreign_keys(ON)", "foreign_keys(on)",
		"foreign_keys(true)", "foreign_keys(yes)", "foreign_keys = 1", "foreign_keys(OFF)",
	} {
		shared := append(append([]string{}, dsnPragmas()...), enforcing)
		got := buildDSN(dbPath, true, migrationPragmas(shared))
		if fk := read("migration/shared "+enforcing, got); fk != 0 {
			t.Errorf("with %q in the shared list the migration connection reads foreign_keys = %d, "+
				"want 0: the transform let the shared entry survive (%s)", enforcing, fk, got)
		}
	}
}
