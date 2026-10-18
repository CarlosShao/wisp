package config

// Ticket 227 - two hand-edit shapes the guarded write reaches that ticket 226's
// AC#1/AC#2 cases never planted, plus the boundary of the "already says it,
// nothing written" promise. This file is tests only; it changes no behaviour.
//
// WHY THESE CASES DID NOT EXIST YET. writeguard_226_test.go plants exactly one
// hand-edit shape: it replaces "size = 56" with "size = 60" (a key whose VALUE
// changed). Nothing there deletes a line, deletes a whole section, or drops a map
// entry, and the "nothing written" case's fixture carries a current schema_version
// - so the migration that readConfigFile runs was out of its reach. Ticket 226's
// AC#1/AC#2 are therefore only true along the "changed a value" branch; this file
// pins the other branches as permanent, red-on-regression cases.
//
// WHAT EACH CASE PROVES (the ticket's own three-for-the-delete-shape list):
//   1. a hand-deleted NON-default scalar line comes back, and its value equals the
//      schema default (56), NOT the value this process held in memory (60);
//   2. kept_in_file_not_in_memory DOES name that key (memory and file diverge);
//   3. the operator's comment on that line is gone, and the log now SAYS the
//      delete-shape falls back to the default rather than hiding it in a comment.

import (
	"reflect"
	"strings"
	"testing"
)

// AC#1①②③: hand-delete a non-default scalar line, run one guarded write.
func TestHandDeletedScalarLineComesBackAtSchemaDefault(t *testing.T) {
	logs := captureSlog(t)
	// Startup load: a current-schema file whose [ball] size is 60 (non-default,
	// so memory holds 60) with the operator's comment on that very line.
	start := strings.Replace(matrixBase,
		"size = 56", "size = 60  # my size", 1)
	m, path, mutate := newTestManager(t, start)
	if m.Config().Ball.Size != 60 {
		t.Fatalf("setup: startup did not load the non-default size 60, got %d", m.Config().Ball.Size)
	}

	// The operator deletes the whole line (and its comment) by hand, leaving the
	// current schema_version in place so this is a pure delete-shape, not a
	// migration.
	deleted := strings.Replace(start, "size = 60  # my size\n", "", 1)
	if strings.Contains(deleted, "size = 60") || strings.Contains(deleted, "my size") {
		t.Fatal("setup: the size line or its comment was not actually deleted")
	}
	mutate(deleted)

	if err := m.AddAllowedDir(testDir); err != nil {
		t.Fatalf("AddAllowedDir: %v", err)
	}

	// ① The line is back, at the schema default (56), not at memory's 60.
	got := readLoaded(t, path)
	if got.Ball.Size != 56 {
		t.Errorf("the hand-deleted size line came back as %d; ticket 227 pins it at the schema "+
			"default 56, NOT the value this process held in memory", got.Ball.Size)
	}
	if !strings.Contains(readText(t, path), "size = 56") {
		t.Errorf("the deleted size line was not restored to the file:\n%s", readText(t, path))
	}
	// And the running process still holds what it loaded at startup (the guarded
	// write merges toward the FILE, it does not apply the fallback to memory).
	if m.Config().Ball.Size != 60 {
		t.Errorf("the guarded write pushed the default back into memory too: got %d, want 60",
			m.Config().Ball.Size)
	}

	// ② kept_in_file_not_in_memory names ball.size (file says 56, memory 60).
	text := logs.String()
	if !strings.Contains(text, "kept_in_file_not_in_memory=[ball.size]") {
		t.Errorf("the fallback-to-default divergence was not reported; log:\n%s", text)
	}

	// ③ the operator's comment is gone, and the write says the delete-shape in
	//    words a user can see (ticket 227 AC#1③: not only in a source comment).
	if strings.Contains(readText(t, path), "my size") {
		t.Error("the operator's comment survived a canonical re-lay; it must be reported gone, not silently kept")
	}
	if !strings.Contains(text, "a key line the operator deleted by hand comes back at its schema default value") {
		t.Errorf("the log must state the delete-shape plainly; log:\n%s", text)
	}
}

// AC#1's map-entry reverse (ticket 226 b2 P4): deleting a provider sub-table is
// NOT filled back in - the file loses that key rather than reviving it.
func TestHandDeletedMapEntryIsNotRestored(t *testing.T) {
	// A current-schema file that declares one provider so memory holds it.
	start := matrixBase + `
[llm]
[llm.providers.deepseek]
base_url = "https://api.deepseek.com/v1"
`
	m, path, mutate := newTestManager(t, start)
	if _, ok := m.Config().LLM.Providers["deepseek"]; !ok {
		t.Fatal("setup: startup did not load the deepseek provider")
	}

	// The operator deletes the [llm.providers.deepseek] table by hand.
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
		t.Errorf("a hand-deleted map entry came BACK; ticket 226 b2 P4 pins this shape as NOT restored")
	}
	if !reflect.DeepEqual(got.FS.AllowedDirs, []string{testDir}) {
		t.Errorf("the owned key did not land: %v", got.FS.AllowedDirs)
	}
}

// AC#2's range boundary: the "already carries this key; nothing written" promise
// holds ONLY when the file already declares the current schema_version. When the
// operator stripped schema_version, readConfigFile's migration re-lays the file
// during the very re-read that ends in "nothing written" - so the file is touched
// (bytes move, mtime moves, a .bak appears) while the log says it was not. This
// case pins that truth so the writeguard_226_test.go "does not rewrite" case can
// no longer be read as unconditional.
func TestNothingWrittenPromiseOnlyHoldsWhenSchemaIsCurrent(t *testing.T) {
	logs := captureSlog(t)
	// Build a schema-less fixture whose [fs] line already carries the entry, with
	// a comment on it so we can see the migration erase it.
	hand := strings.Replace(matrixBase,
		"allowed_dirs = []", `allowed_dirs = ["D:\\data"]  # hand set`, 1)
	noVersion := strings.Replace(hand, "schema_version = 2\n", "", 1)
	if strings.Contains(noVersion, "schema_version") {
		t.Fatal("setup: schema_version line was not removed")
	}

	// Startup load through matrixBase (a current-schema file), then hand-plant the
	// schema-less form: memory is a clean v2 snapshot with an empty allowlist, the
	// file already says D:\data but no longer declares its version.
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
	// The sentence must be the narrowed, conditional one. If it ever reverts to
	// the flat "nothing written", this case goes red rather than letting the lie
	// back in quietly (ticket 227 AC#2 甲: the log says what actually happens).
	if !strings.Contains(text, "only because its schema_version was already current") {
		t.Errorf("the nothing-written sentence must state the schema-current condition it depends on; log:\n%s", text)
	}
	// ...yet the file WAS re-laid by the migration during the re-read.
	if readText(t, path) == before {
		t.Error("the schema-less file was NOT re-laid even though the log said nothing written - "+
			"either the migration moved or this case stopped proving the boundary")
	}
	if statMod(t, path).Equal(mtimeBefore) {
		t.Error("mtime did not move on a write that logged 'nothing written' with a stale schema - "+
			"the migration-on-read boundary is no longer being exercised")
	}
	if !strings.Contains(readText(t, path), "schema_version = 2") {
		t.Errorf("the migrated file does not declare the current schema:\n%s", readText(t, path))
	}
	if strings.Contains(readText(t, path), "hand set") {
		t.Error("the operator's comment survived the migration-on-read; the boundary must report it gone")
	}
}

// The paired POSITIVE CONTROL the ticket demands: the same "already says it,
// nothing written" shape, but with schema_version already current (the exact
// fixture writeguard_226_test.go:393 uses). There the promise is true - the file
// is byte-for-byte untouched. This is what keeps the boundary case above from
// being a case that is red only because its setup always mutates the file.
func TestNothingWrittenLeavesFileUntouchedWhenSchemaIsCurrent(t *testing.T) {
	m, path, mutate := newTestManager(t, matrixBase)
	hand := strings.Replace(matrixBase,
		"allowed_dirs = []", `allowed_dirs = ["D:\\data"]  # hand set`, 1)
	mutate(hand)
	before := readText(t, path)
	mtimeBefore := statMod(t, path)

	if err := m.AddAllowedDir(testDir); err != nil {
		t.Fatalf("AddAllowedDir: %v", err)
	}
	if readText(t, path) != before {
		t.Errorf("a current-schema file that already carried the key was rewritten:\n%s", readText(t, path))
	}
	if !statMod(t, path).Equal(mtimeBefore) {
		t.Error("mtime moved on a true nothing-written write (schema current) - the control lost its teeth")
	}
	if !strings.Contains(before, "hand set") {
		t.Fatal("setup: control fixture had no comment to keep")
	}
	if !strings.Contains(readText(t, path), "hand set") {
		t.Error("a current-schema nothing-written write destroyed the operator's comment")
	}
}

// AC#4's nameless guard: the write must claim exactly the owned key and no other.
// A deleted-line fixture (the shape this file's AC#1 case uses) is run here to
// prove that even when the re-read fills a key back in, the report never lists a
// key other than fs.allowed_dirs under wrote= - this leg widens coverage, not the
// merge's reach.
func TestHandEditCasesWriteOnlyTheOwnedKey(t *testing.T) {
	logs := captureSlog(t)
	m, _, mutate := newTestManager(t, matrixBase)
	deleted := strings.Replace(matrixBase, "size = 56\n", "", 1)
	mutate(deleted)
	if err := m.AddAllowedDir(testDir); err != nil {
		t.Fatalf("AddAllowedDir: %v", err)
	}
	text := logs.String()
	if !strings.Contains(text, "wrote=[fs.allowed_dirs]") {
		t.Errorf("the owned key was not the one written; log:\n%s", text)
	}
	if strings.Contains(text, "wrote=[ball.size") {
		t.Errorf("the write claimed a key it did not own (AC#4 must not widen); log:\n%s", text)
	}
}
