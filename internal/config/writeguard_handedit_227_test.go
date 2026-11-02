package config

// Ticket 227 - hand-edit shapes the guarded write reaches that ticket 226 never
// planted (delete a line / drop a map entry) plus the boundary of the "already
// says it" promise. Tests only; no behaviour changes (AC#4). 226 planted only a
// changed VALUE ("size = 56" -> "size = 60") on a fixture whose schema_version was
// current, so its AC#1/AC#2 held along that branch only: nothing here deletes a
// line. Ticket 345 (2026-11-02) re-nailed ONE case of this file - the schema-
// boundary case below - after migration left the read pipeline; the delete-shape
// cases above are byte-for-byte the promises 227 pinned and stay green.

import (
	"os"
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

// AC#2's boundary, RE-NAILED by ticket 345 (2026-11-02, shape A, owner-approved;
// this replaces TestNothingWrittenPromiseOnlyHoldsWhenSchemaIsCurrent - no
// assertion was deleted, the pinned TRUTH is the opposite half, named below):
//
//   - OLD promise: "this write added nothing" held ONLY for a current-schema
//     file, because a behind-version file had already been re-laid in canonical
//     form by the migration INSIDE the re-read - bytes moved, mtime moved, the
//     operator's comment died, and the file said schema_version = 2, all while
//     the program logged that it had added nothing.
//   - NEW promise: the nothing-written branch writes NOTHING, whatever the
//     file's schema_version declares. The re-read migrates in memory only
//     (ticket 345), so the bytes stay put, the mtime stays put, the operator's
//     comment survives, no config.toml.bak-1 appears (no SaveFile ran), and the
//     file still declares no schema_version because nothing rewrote it.
//   - AC#4's written-key half is untouched: writeguard_226_test.go:126-131
//     still pins wrote=[fs.allowed_dirs] and forbids wrote=[ball.size....
func TestNothingWrittenStaysUntouchedEvenWhenSchemaIsStale(t *testing.T) {
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
	if !strings.Contains(before, "hand set") {
		t.Fatal("setup: the stale fixture had no comment to keep")
	}

	if err := m.AddAllowedDir(testDir); err != nil {
		t.Fatalf("AddAllowedDir: %v", err)
	}

	text := logs.String()
	if !strings.Contains(text, "already carries this key") {
		t.Fatalf("expected the nothing-written branch to fire; log:\n%s", text)
	}
	// The sentence may no longer state the schema-current condition: under
	// ticket 345 the promise is unconditional.
	if strings.Contains(text, "only because its schema_version was already current") {
		t.Errorf("the nothing-written sentence still claims the old schema-current condition; log:\n%s", text)
	}
	if got := readText(t, path); got != before {
		t.Errorf("ticket 345: a nothing-written branch on a stale-schema file must leave the bytes alone; got:\n%s", got)
	}
	if !statMod(t, path).Equal(mtimeBefore) {
		t.Error("mtime moved although this branch writes nothing")
	}
	if !strings.Contains(readText(t, path), "hand set") {
		t.Error("the operator's comment did not survive the re-read - the whole point of ticket 345")
	}
	if strings.Contains(readText(t, path), "schema_version") {
		t.Error("the stale file was re-laid after all - migration-on-read is back")
	}
	if _, err := os.Stat(path + ".bak-1"); !os.IsNotExist(err) {
		t.Errorf("no backup may appear when no write ran (stat err=%v)", err)
	}
}

// The paired control under the NEW promise: nothing-written leaves the file
// untouched for a stale-schema file too - and for this current-schema control
// it must stay untouched exactly as before, which keeps the case above from
// proving a blanket "nothing is ever written" rather than the schema-
// unconditional boundary of ticket 345.
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
