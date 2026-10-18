package config

// Ticket 227 - hand-edit shapes the guarded write reaches that ticket 226 never
// planted (delete a line / drop a map entry) plus the boundary of the "already
// says it" promise. Tests only; no behaviour changes (AC#4). 226 planted only a
// changed VALUE ("size = 56" -> "size = 60") on a fixture whose schema_version was
// current, so its AC#1/AC#2 held along that branch only: nothing here deletes a
// line, and nothing there reached the migration inside readConfigFile. These cases
// pin the delete shapes and that migration, permanently, red on regression.

import (
	"reflect"
	"strings"
	"testing"
)

// AC#1, the ticket's three assertions in one case: a hand-deleted NON-default
// scalar line comes back at the schema default (not at memory's value),
// kept_in_file_not_in_memory names it as kept rather than swallowed, the
// operator's comment on that line is gone, and the log says the delete-shape out
// loud - ticket 227 AC#1③ requires the sentence a user can see, not only a source
// comment.
func TestHandDeletedScalarLineComesBackAtSchemaDefault(t *testing.T) {
	logs := captureSlog(t)
	start := strings.Replace(matrixBase, "size = 56", "size = 60  # my size", 1)
	m, path, mutate := newTestManager(t, start)
	if m.Config().Ball.Size != 60 {
		t.Fatalf("setup: startup did not load the non-default size 60, got %d", m.Config().Ball.Size)
	}

	deleted := strings.Replace(start, "size = 60  # my size\n", "", 1)
	if strings.Contains(deleted, "size = 60") || strings.Contains(deleted, "my size") {
		t.Fatal("setup: the size line or its comment was not actually deleted")
	}
	mutate(deleted)
	if err := m.AddAllowedDir(testDir); err != nil {
		t.Fatalf("AddAllowedDir: %v", err)
	}

	// AC#1 (1): back at the schema default, not at memory's value.
	got := readLoaded(t, path)
	if got.Ball.Size != 56 {
		t.Errorf("the hand-deleted size line came back as %d, want the schema default 56: a deleted "+
			"line falls back to the default, NOT to the value held in memory (ticket 227)", got.Ball.Size)
	}
	if !strings.Contains(readText(t, path), "size = 56") {
		t.Errorf("the deleted size line was not restored to the file:\n%s", readText(t, path))
	}
	if m.Config().Ball.Size != 60 {
		t.Errorf("the guarded write pushed the fallback into memory too: got %d, want 60", m.Config().Ball.Size)
	}

	// (2) the divergence is reported, not swallowed.
	text := logs.String()
	if !strings.Contains(text, "kept_in_file_not_in_memory=[ball.size]") {
		t.Errorf("the fallback-to-default divergence was not reported; log:\n%s", text)
	}

	// (3) the comment is gone AND the log names the delete-shape.
	if strings.Contains(readText(t, path), "my size") {
		t.Error("the operator's comment survived a canonical re-lay; it must be reported gone, not kept silently")
	}
	if !strings.Contains(text, "a key line the operator deleted by hand comes back at its schema default value") {
		t.Errorf("the log must state the delete-shape plainly; log:\n%s", text)
	}
}

// AC#1's map-entry half (226 b2 P4, reversed): a deleted provider sub-table is
// NOT filled back in.
func TestHandDeletedMapEntryIsNotRestored(t *testing.T) {
	start := matrixBase + `
[llm]
[llm.providers.deepseek]
base_url = "https://api.deepseek.com/v1"
`
	m, path, mutate := newTestManager(t, start)
	if _, ok := m.Config().LLM.Providers["deepseek"]; !ok {
		t.Fatal("setup: startup did not load the deepseek provider")
	}

	noTable := strings.Replace(start,
		"\n[llm]\n[llm.providers.deepseek]\nbase_url = \"https://api.deepseek.com/v1\"\n", "\n", 1)
	if strings.Contains(noTable, "llm.providers.deepseek") {
		t.Fatal("setup: the provider table was not actually deleted")
	}
	mutate(noTable)
	if err := m.AddAllowedDir(testDir); err != nil {
		t.Fatalf("AddAllowedDir: %v", err)
	}

	got := readLoaded(t, path)
	if _, ok := got.LLM.Providers["deepseek"]; ok {
		t.Error("a hand-deleted map entry came BACK; this shape is pinned as NOT restored (b2 P4)")
	}
	if !reflect.DeepEqual(got.FS.AllowedDirs, []string{testDir}) {
		t.Errorf("the owned key did not land: %v", got.FS.AllowedDirs)
	}
}

// AC#2's boundary: the nothing-written promise holds only when the file's
// schema_version is already current. Strip it and the migration inside
// readConfigFile re-lays the file DURING the re-read that ends in that log line:
// bytes move, mtime moves, the comment is gone, while the program says it wrote
// nothing. AC#4's written-key half needs no new case: writeguard_226_test.go:126-131
// already pins wrote=[fs.allowed_dirs] and forbids wrote=[ball.size....
func TestNothingWrittenPromiseOnlyHoldsWhenSchemaIsCurrent(t *testing.T) {
	logs := captureSlog(t)
	hand := strings.Replace(matrixBase,
		"allowed_dirs = []", `allowed_dirs = ["D:\\data"]  # hand set`, 1)
	noVersion := strings.Replace(hand, "schema_version = 2\n", "", 1)
	if strings.Contains(noVersion, "schema_version") {
		t.Fatal("setup: the schema_version line was not removed")
	}

	m, path, mutate := newTestManager(t, matrixBase)
	mutate(noVersion)
	before := readText(t, path)
	mtimeBefore := statMod(t, path)

	if err := m.AddAllowedDir(testDir); err != nil {
		t.Fatalf("AddAllowedDir: %v", err)
	}

	text := logs.String()
	if !strings.Contains(text, "already carries this key") {
		t.Fatalf("expected the nothing-written branch to fire; log:\n%s", text)
	}
	if !strings.Contains(text, "only because its schema_version was already current") {
		t.Errorf("the nothing-written sentence must state the schema-current condition; log:\n%s", text)
	}
	if readText(t, path) == before {
		t.Error("the schema-less file was NOT re-laid although the log said this write added nothing - " +
			"either the migration stopped moving the file or this case stopped proving the boundary")
	}
	if statMod(t, path).Equal(mtimeBefore) {
		t.Error("mtime did not move on a stale-schema file that logged 'added nothing' - the " +
			"migration-on-read boundary is no longer being exercised")
	}
	if !strings.Contains(readText(t, path), "schema_version = 2") {
		t.Errorf("the migrated file does not declare the current schema:\n%s", readText(t, path))
	}
	if strings.Contains(readText(t, path), "hand set") {
		t.Error("the operator's comment survived the migration-on-read; the boundary must report it gone")
	}
}

// The paired control that keeps the case above from being red because its setup
// always mutates the file: same shape, schema already current, promise is true -
// byte-for-byte untouched, comment still there.
func TestNothingWrittenLeavesFileUntouchedWhenSchemaIsCurrent(t *testing.T) {
	m, path, mutate := newTestManager(t, matrixBase)
	hand := strings.Replace(matrixBase,
		"allowed_dirs = []", `allowed_dirs = ["D:\\data"]  # hand set`, 1)
	mutate(hand)
	before := readText(t, path)
	mtimeBefore := statMod(t, path)
	if !strings.Contains(before, "hand set") {
		t.Fatal("setup: the control fixture had no comment to keep")
	}

	if err := m.AddAllowedDir(testDir); err != nil {
		t.Fatalf("AddAllowedDir: %v", err)
	}
	if readText(t, path) != before {
		t.Errorf("a current-schema file that already carried the key was rewritten:\n%s", readText(t, path))
	}
	if !statMod(t, path).Equal(mtimeBefore) {
		t.Error("mtime moved on a true nothing-written write - the control lost its teeth")
	}
	if !strings.Contains(readText(t, path), "hand set") {
		t.Error("a current-schema nothing-written write destroyed the operator's comment")
	}
}

// AC#4's guard, in the NEW shape: writeguard_226_test.go:126-131 already pins
// wrote=[fs.allowed_dirs] and forbids wrote=[ball.size... along the changed-value
// branch. Along the delete branch that is not the same proof: here the file was
// missing ball.size and the write put the default back, so a report that claimed
// ball.size as a key this write changed would be a lie about a key nobody asked it
// to touch. This leg adds coverage, it does not widen the merge.
func TestDeletedLineRestorationIsNotClaimedAsAWrittenKey(t *testing.T) {
	logs := captureSlog(t)
	m, _, mutate := newTestManager(t, matrixBase)
	mutate(strings.Replace(matrixBase, "size = 56\n", "", 1))
	if err := m.AddAllowedDir(testDir); err != nil {
		t.Fatalf("AddAllowedDir: %v", err)
	}
	text := logs.String()
	if !strings.Contains(text, "wrote=[fs.allowed_dirs]") {
		t.Errorf("the owned key was not the one written; log:\n%s", text)
	}
	if strings.Contains(text, "wrote=[ball.size") {
		t.Errorf("the write claimed a key it does not own (AC#4 forbids widening); log:\n%s", text)
	}
}
