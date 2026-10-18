package config

// Ticket 341 AC#1 - the four hand-edit shapes ticket 227's matrix reported,
// honestly, as still unnailed: P2 (delete a scalar line whose value was ALREADY
// the schema default), P3 (delete a WHOLE SECTION, the [ball] header included),
// P5 (delete the schema_version line while the WRITE branch is the one that runs)
// and P6 (delete the base_url line of one provider). 227 pinned P1 (:23), P4 (:70)
// and P5' (:106) in writeguard_handedit_227_test.go; these four had one-off bench
// readings or nothing at all. This file adds them permanently and changes no
// behaviour: not one non-comment byte of writeguard.go, loader.go or migrate.go
// moved (ticket 341 AC#5), and the migration-on-read behaviour change that P5
// exposes is ticket 345's business (the "yi" candidate), not this ticket's.
//
// Denominator = 226-v1's bench piece (not re-invented):
// .scratch/wisp/probes/226/v1/b2-probe-deleted-shapes.txt - P2 :12-18, P3 :20-27,
// P5 :38-44, P6 :46-52. 227-v1's independent census
// (.scratch/wisp/probes/227/v1/logs/06-ac3-independent-census.txt) reached the
// same four gaps. The census this leg ran for itself (every shape, its own
// content anchors, hit counts and carrier lines at anchor 98dbf20a) is
// .scratch/wisp/probes/341/r1/logs/01-start-gauges.txt, gauge 3.
//
// Each case names the guard it pins and the ONE content-anchored line whose
// removal makes it red. The four de-gum landings are four different files or four
// different lines (AC#1 requires per-shape distinct landings); the shots, their
// before/after anchor hit counts, the verbatim red lines and the restore hashes
// are in .scratch/wisp/probes/341/r1/logs/05-mu-p2.txt, 05-mu-p3.txt,
// 05-mu-p5.txt and 05-mu-p6.txt.

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// P2 - delete a scalar line whose value was ALREADY the schema default.
//
// Why this is not P1 wearing a hat (ticket 341 AC#2 forbids merging them): the
// shape looks identical - one line gone - and the RESULT is not. P1's deleted line
// held 60 while memory held 60, so the file came back at 56 and the divergence is
// visible: the WARN branch prints kept_in_file_not_in_memory=[ball.size]
// (writeguard.go:184-190). Here the deleted line held 56, memory holds 56 and the
// restore is 56 too, so nothing in the report can point at the event: b2 P2
// (:12-18) measured the INFO branch (writeguard.go:172-176), where the operator's
// comment is destroyed with no divergence list to show for it.
//
// "You cannot see the difference" is the fact this nail pins - and it pins it AS
// "the sentence still has to say what was lost". The restore itself is silent by
// construction here, so the only non-silent part of this shape is the delete-shape
// clause in the INFO branch. That clause is what its de-gum shot removes
// (writeguard.go:175), which is why this case is the one that goes red when a
// later edit shortens the sentence back to "only the keys listed changed value".
//
// Guards pinned, in order: the re-read overlays the file onto NewDefaults
// (loader.go:131) so a deleted line returns at the schema default; MarshalCanonical
// emits every static field (parse.go:160-163) so the line is physically back; the
// INFO branch names the delete-shape out loud (writeguard.go:172-176).
func TestHandDeletedDefaultValuedLineIsRestoredInvisibly(t *testing.T) {
	logs := captureSlog(t)
	const p2Default = 56
	if NewDefaults().Ball.Size != p2Default {
		t.Fatalf("setup: this shape means the deleted line held the schema default; ball.size default is now %d, so the fixture and this case need re-cutting (report it, do not relax it)",
			NewDefaults().Ball.Size)
	}
	start := strings.Replace(matrixBase, "size = 56", "size = 56  # default to me", 1)
	m, path, mutate := newTestManager(t, start)
	if m.Config().Ball.Size != p2Default {
		t.Fatalf("setup: memory holds %d and the schema default is %d - P2 requires them equal",
			m.Config().Ball.Size, p2Default)
	}

	deleted := strings.Replace(start, "size = 56  # default to me\n", "", 1)
	if strings.Contains(deleted, "size") || strings.Contains(deleted, "default to me") {
		t.Fatal("setup: the size line or its comment was not actually deleted")
	}
	mutate(deleted)
	if err := m.AddAllowedDir(testDir); err != nil {
		t.Fatalf("AddAllowedDir: %v", err)
	}

	// (1) the line is physically back, at the default.
	text := readText(t, path)
	if !strings.Contains(text, "size = 56") {
		t.Errorf("the hand-deleted default-valued line did not come back at the schema default:\n%s", text)
	}
	// (2) the comment went with it: a guarded write re-lays canonical form.
	if strings.Contains(text, "default to me") {
		t.Error("the operator's comment survived a canonical re-lay; the cost must be reported, not kept silently")
	}

	// (3) the report cannot see this shape, and that is the reading: memory and
	// file agree, so mergeWrite took the INFO branch and no divergence list was
	// printed. If a later build starts printing one here, that means the two
	// branches converged and P2/P1 need re-cutting - so this stays a hard
	// assertion, and this case is never skipped.
	logText := logs.String()
	if !strings.Contains(logText, "wrote merged change into config.toml") {
		t.Errorf("expected a real write (the file held no allowed_dirs entry); log:\n%s", logText)
	}
	if strings.Contains(logText, "kept hand edits this process does not hold") ||
		strings.Contains(logText, "kept_in_file_not_in_memory") {
		t.Errorf("a deleted default-valued line is invisible to the divergence report - that is P2's whole point; log:\n%s", logText)
	}
	if strings.Contains(logText, "wrote=[ball.size") {
		t.Errorf("the write claimed a key it does not own; log:\n%s", logText)
	}

	// (4) and still the INFO sentence names the delete-shape out loud, because the
	// numbers alone hide the event: this is the assertion its de-gum shot kills.
	if !strings.Contains(logText, "a key line the operator deleted by hand comes back at its schema default value") {
		t.Errorf("the INFO branch must state the delete-shape even when the values hide it; log:\n%s", logText)
	}
	// (5) file and memory end up equal: indistinguishable from a write that
	// restored nothing.
	if got := readLoaded(t, path); got.Ball.Size != p2Default || m.Config().Ball.Size != p2Default {
		t.Errorf("file says %d and memory says %d, want %d and %d", got.Ball.Size, m.Config().Ball.Size, p2Default, p2Default)
	}
}

// P3 - delete a WHOLE SECTION: the [ball] header plus every line under it.
//
// b2 P3 (:20-27) measured this once and the census above hits zero carriers for a
// section-level delete, so it had no regression protection at all. A section delete
// is not a key delete: the re-read finds no table at all, so every field of
// BallSection returns from its own tag default - including keys the operator never
// wrote (opacity_idle, click_through, hide_on_fullscreen) - while the divergence
// report names exactly ONE leaf, ball.size, the only default that differs from
// memory. A section-level event is therefore reported at key-level granularity, and
// the resurrected keys are invisible to that report: that asymmetry is what
// "the cost of deleting a section differs from the cost of deleting a key" means
// in practice, and it is why this shape needs a nail of its own next to P1.
//
// Guard pinned: the merge base comes from the FILE, not from memory
// (writeguard.go:140-142, the whole point of candidate jia). De-gum landing: that
// deepCopyConfig(disk) replaced by deepCopyConfig(m.cur) - the deleted section then
// comes back at memory's 60 instead of the schema default 56, which is precisely
// the silent revert ticket 226 exists to stop.
func TestHandDeletedWholeSectionComesBackFromDefaultsNotMemory(t *testing.T) {
	logs := captureSlog(t)
	start := strings.Replace(matrixBase, "size = 56", "size = 60  # my size", 1)
	m, path, mutate := newTestManager(t, start)
	if m.Config().Ball.Size != 60 {
		t.Fatalf("setup: memory did not load the non-default size, got %d", m.Config().Ball.Size)
	}

	deleted := strings.Replace(start, "[ball]\nsize = 60  # my size\n\n", "", 1)
	if strings.Contains(deleted, "[ball]") || strings.Contains(deleted, "size =") || strings.Contains(deleted, "my size") {
		t.Fatal("setup: the [ball] section was not deleted whole (header plus body)")
	}
	mutate(deleted)
	if err := m.AddAllowedDir(testDir); err != nil {
		t.Fatalf("AddAllowedDir: %v", err)
	}

	// (1) the section header itself is back.
	text := readText(t, path)
	if !strings.Contains(text, "[ball]") {
		t.Errorf("the deleted section header did not come back:\n%s", text)
	}
	// (2) and it came back at the SCHEMA DEFAULT, not at memory's value.
	if !strings.Contains(text, "size = 56") || strings.Contains(text, "size = 60") {
		t.Errorf("the resurrected section carried memory's value instead of the schema default:\n%s", text)
	}
	// (3) keys nobody ever wrote come back with the section: the section-level cost.
	if !strings.Contains(text, "opacity_idle = 0.35") || !strings.Contains(text, "click_through = true") {
		t.Errorf("the re-laid section did not carry the tag defaults nobody wrote:\n%s", text)
	}
	if strings.Contains(text, "my size") {
		t.Error("the operator's comment survived a canonical re-lay")
	}
	// (4) one leaf is reported for a section-level event; the phantoms are not.
	logText := logs.String()
	if !strings.Contains(logText, "kept_in_file_not_in_memory=[ball.size]") {
		t.Errorf("the size divergence was not reported; log:\n%s", logText)
	}
	if strings.Contains(logText, "ball.opacity_idle") || strings.Contains(logText, "ball.click_through") {
		t.Errorf("the report named a resurrected key the operator never had - a claim this write did not produce; log:\n%s", logText)
	}
	if got := readLoaded(t, path); got.Ball.Size != 56 || !got.Ball.ClickThrough {
		t.Errorf("a reload reads size=%d click_through=%v, want 56 and true", got.Ball.Size, got.Ball.ClickThrough)
	}
	if m.Config().Ball.Size != 60 {
		t.Errorf("the guarded write pushed the fallback into memory too: %d, want 60", m.Config().Ball.Size)
	}
}

// P5 - delete the schema_version line while the WRITE branch is the one that runs.
//
// 227 pinned the other half only (writeguard_handedit_227_test.go:106, the
// nothing-written branch); b2 P5 (:38-44) measured this one as a bench reading.
// They are two different paths through mergeWrite, and running both is what makes
// "which branch did a schema-less file take" a reading instead of an assumption:
// here the file still disagrees with memory on the owned key, so the write really
// happens - and the migration inside the re-read has already rewritten the file
// before that write (readConfigFile sees version 0, treats it as v1,
// applyMigrations replaces the bytes and drops config.toml.bak-1).
//
// What this case pins is what the file promises afterwards: the deleted version
// line comes back FORCED to the current number, and the report must not say two
// things about it - schema_version is not a key this write wrote, and it is not a
// hand edit preserved either. diffKeyPaths skips SchemaVersion on purpose
// (writeguard.go:225-229) precisely because SaveFile stamps it, so a claim naming
// it would be a lie about a key nobody asked this write to touch.
//
// Guard pinned: SaveFile stamps the current schema version (loader.go:257). De-gum
// landing: that stamp moved to 1 - the version the operator deleted then comes
// back WRONG and nothing in the report says so, which is the failure this nail
// exists for.
func TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed(t *testing.T) {
	logs := captureSlog(t)
	m, path, mutate := newTestManager(t, matrixBase)
	hand := strings.Replace(matrixBase, "size = 56", "size = 56  # hand note", 1)
	noVersion := strings.Replace(hand, "schema_version = 2\n", "", 1)
	if strings.Contains(noVersion, "schema_version") {
		t.Fatal("setup: the schema_version line was not removed")
	}
	mutate(noVersion)
	before := readText(t, path)
	mtimeBefore := statMod(t, path)
	if strings.Contains(before, "schema_version") {
		t.Fatal("setup: the file on disk still declares a version")
	}

	if err := m.AddAllowedDir(testDir); err != nil {
		t.Fatalf("AddAllowedDir: %v", err)
	}

	// (1) the version line is back at the CURRENT number: forced, not merged.
	text := readText(t, path)
	if !strings.Contains(text, "schema_version = 2") {
		t.Fatalf("the write branch did not re-stamp the deleted schema_version line at the current version:\n%s", text)
	}
	// (2) this is a write: bytes and mtime moved (the migration had already moved
	// them once inside the re-read).
	if text == before {
		t.Error("a schema-less file was left byte-identical although this write owned a change")
	}
	if statMod(t, path).Equal(mtimeBefore) {
		t.Error("mtime did not move on the write branch")
	}
	// (3) the operator's comment died in the canonical re-lay, and the migration
	// kept the pre-migration bytes verbatim as the backup: that file is the one
	// place the destroyed comment still exists, so it is pinned here too.
	if strings.Contains(text, "hand note") {
		t.Error("the operator's comment survived the migration-on-read plus write")
	}
	bak, err := os.ReadFile(path + ".bak-1")
	if err != nil {
		t.Fatalf("the migration backup config.toml.bak-1 is missing: %v", err)
	}
	if !bytes.Contains(bak, []byte("hand note")) || bytes.Contains(bak, []byte("schema_version")) {
		t.Errorf("the backup is not the pre-migration bytes verbatim:\n%s", bak)
	}
	// (4) the sentence says WRITE, never borrows the nothing-written promise, and
	// never claims schema_version among the keys it wrote or kept.
	logText := logs.String()
	if !strings.Contains(logText, "wrote merged change into config.toml") {
		t.Errorf("expected the write branch sentence; log:\n%s", logText)
	}
	if strings.Contains(logText, "already carries this key") {
		t.Errorf("the write branch must not say the nothing-written sentence; log:\n%s", logText)
	}
	if strings.Contains(logText, "wrote=[schema_version") ||
		strings.Contains(logText, "kept_in_file_not_in_memory=[schema_version") {
		t.Errorf("schema_version was claimed as a key this write touched; log:\n%s", logText)
	}
	if !strings.Contains(logText, "wrote=[fs.allowed_dirs]") {
		t.Errorf("the owned key was not the key written; log:\n%s", logText)
	}
}

// P6 - delete the base_url line of ONE provider, inside a table that survives.
//
// Adjacent to 227's P4 (which pins the map ENTRY being dropped and NOT restored)
// but not the same shape: b2 P6 (:46-52) is a leaf delete inside a live table. The
// provider name here is a built-in preset, so the re-read refills base_url from the
// preset table (applyPresets, loader.go:202-218) - the line comes back at NEITHER
// what the operator wrote NOR what memory holds, but at the endpoint this build
// ships. It shows in the divergence report only because memory still holds the hand
// value; the file now carries a URL the operator never typed. That is the same harm
// P1 pins (a deleted line silently replaced by a built-in default) reached through a
// different seam: the preset table instead of the struct tags, which is why it
// needs its own nail.
//
// Guard pinned: the preset fill (loader.go:213, p.BaseURL = pre.BaseURL). De-gum
// landing: that fill made a no-op. The paired control below - a provider name no
// preset knows - must stay GREEN while this one goes RED; that is what makes the
// shot a reading about the preset path rather than just another broken build.
func TestHandDeletedProviderBaseUrlComesBackAsThePresetValue(t *testing.T) {
	logs := captureSlog(t)
	const handURL = "https://my-gateway.example.invalid/v1"
	const presetURL = "https://api.deepseek.com/v1"
	start := matrixBase + `
[llm]
[llm.providers.deepseek]
base_url = "https://my-gateway.example.invalid/v1"  # mine
`
	m, path, mutate := newTestManager(t, start)
	if got := m.Config().LLM.Providers["deepseek"].BaseURL; got != handURL {
		t.Fatalf("setup: memory did not load the hand base_url, got %q", got)
	}

	deleted := strings.Replace(start, "base_url = \"https://my-gateway.example.invalid/v1\"  # mine\n", "", 1)
	if strings.Contains(deleted, "base_url") || strings.Contains(deleted, "mine") {
		t.Fatal("setup: the base_url line or its comment was not actually deleted")
	}
	mutate(deleted)
	if err := m.AddAllowedDir(testDir); err != nil {
		t.Fatalf("AddAllowedDir: %v", err)
	}

	text := readText(t, path)
	// MarshalCanonical emits TOML literal strings, so every quoted value in the
	// re-laid file lands between single quotes - measured here, and the house
	// convention (b2's own dump says base_url = ''; settings_313_test.go:270
	// asserts model = 'deepseek-reasoner'). An earlier draft of this case matched
	// double quotes and was red on that alone: a ruler defect, fixed here, ⛔ a
	// finding about the behaviour.
	if !strings.Contains(text, "base_url = 'https://api.deepseek.com/v1'") {
		t.Errorf("the deleted base_url line came back as something other than the built-in preset:\n%s", text)
	}
	if strings.Contains(text, "my-gateway") || strings.Contains(text, "mine") {
		t.Error("the operator's URL or comment survived the re-lay")
	}
	// memory is not pushed onto the preset value: the substitution lives in the file.
	if got := m.Config().LLM.Providers["deepseek"].BaseURL; got != handURL {
		t.Errorf("the guarded write pushed the preset endpoint into memory: %q", got)
	}
	if got := readLoaded(t, path); got.LLM.Providers["deepseek"].BaseURL != presetURL {
		t.Errorf("a reload reads %q, want the preset endpoint the file now carries",
			got.LLM.Providers["deepseek"].BaseURL)
	}
	logText := logs.String()
	if !strings.Contains(logText, "kept_in_file_not_in_memory=[llm.providers.deepseek.base_url]") {
		t.Errorf("the preset substitution was not reported as file-holds-what-memory-does-not; log:\n%s", logText)
	}
	if strings.Contains(logText, "wrote=[llm.providers.deepseek.base_url") {
		t.Errorf("the write claimed a key it does not own; log:\n%s", logText)
	}
}

// P6's paired control - the same leaf delete under a provider name the preset table
// does not know (its protocol therefore has to be spelled out in the file,
// validate.go:208-213, and it stays in the file when only base_url is deleted).
// Nothing can refill base_url here, so the file comes back with an EXPLICIT empty
// line and that empty line is what the next start reads. This case exists for the
// shot above: it must stay green when the preset fill is removed, which localises
// that mutation to the preset path.
func TestHandDeletedBaseUrlOfAnUnknownProviderStaysEmptyInTheFile(t *testing.T) {
	logs := captureSlog(t)
	start := matrixBase + `
[llm]
[llm.providers.mycustom]
protocol = "openai-chat"
base_url = "https://my-gateway.example.invalid/v1"
`
	m, path, mutate := newTestManager(t, start)
	const handURL = "https://my-gateway.example.invalid/v1"
	if got := m.Config().LLM.Providers["mycustom"].BaseURL; got != handURL {
		t.Fatalf("setup: memory did not load the hand base_url, got %q", got)
	}

	deleted := strings.Replace(start, "base_url = \"https://my-gateway.example.invalid/v1\"\n", "", 1)
	if strings.Contains(deleted, "base_url") {
		t.Fatal("setup: the base_url line was not actually deleted")
	}
	mutate(deleted)
	if err := m.AddAllowedDir(testDir); err != nil {
		t.Fatalf("AddAllowedDir: %v", err)
	}

	text := readText(t, path)
	if !strings.Contains(text, "base_url = ''") {
		t.Errorf("the deleted line of an unpresetted provider did not come back as an explicit empty string:\n%s", text)
	}
	if strings.Contains(text, "my-gateway") {
		t.Error("an unpresetted provider kept the operator's URL; the file is re-laid, not patched")
	}
	if !strings.Contains(text, "protocol = 'openai-chat'") {
		t.Errorf("the untouched key of the same table did not survive the re-lay:\n%s", text)
	}
	if got := readLoaded(t, path); got.LLM.Providers["mycustom"].BaseURL != "" {
		t.Errorf("a reload reads %q, want the empty string nothing refills", got.LLM.Providers["mycustom"].BaseURL)
	}
	logText := logs.String()
	if !strings.Contains(logText, "kept_in_file_not_in_memory=[llm.providers.mycustom.base_url]") {
		t.Errorf("the empty substitution was not reported; log:\n%s", logText)
	}
}
