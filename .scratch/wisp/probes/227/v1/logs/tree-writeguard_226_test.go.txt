package config

// Ticket 226 - a programmatic config write must not overwrite the operator's own
// hand edits, and must not claim their content as the program's own write.
//
// This file is the FIRST direct coverage of the [fs] write path: before this
// ticket, allowdirs.go had no test of its own (its behaviour was reached only
// sideways through cmd/wisp's reply-listener cases), which is how "the whole
// snapshot goes to SaveFile" survived a landing. Every case below therefore
// asserts the FILE, read back through LoadFile, not a field on a struct.
//
// The pair that carries the load is AC#1's: TestAC1... above says the hand edit
// survives a guarded write, and TestAC1Control... below says the shape this path
// had before (snapshot -> SaveFile) still reverts it. Without the second case the
// first could be green because the test never changed a second key.
//
// The two safety behaviours ticket 201's ruling pinned are pinned here too, so a
// later edit cannot soften them by accident: an entry with a relative hop is
// refused, and a failed write rolls the in-memory value back.

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/risk"
)

// handEdited is matrixBase with one key the operator changed while the process
// ran - a [ball] size, which has nothing whatever to do with [fs].
var handEdited = strings.Replace(matrixBase, "size = 56", "size = 60", 1)

const testDir = `D:\data`

// readLoaded reads the file back the way the next start would.
func readLoaded(t *testing.T, path string) *Config {
	t.Helper()
	c, _, err := LoadFile(path, nil)
	if err != nil {
		t.Fatalf("the file on disk does not load after the write: %v", err)
	}
	return c
}

// captureSlog sends the process default logger into a buffer for one test. Same
// instrument manager_test.go's logging case uses.
func captureSlog(t *testing.T) *bytes.Buffer {
	t.Helper()
	buf := &bytes.Buffer{}
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	t.Cleanup(func() { slog.SetDefault(old) })
	return buf
}

// AC#1: 手改不许被吞.
func TestAC1AllowedDirsWriteKeepsAHandEditedKey(t *testing.T) {
	m, path, mutate := newTestManager(t, matrixBase)
	mutate(handEdited)

	if err := m.AddAllowedDir(testDir); err != nil {
		t.Fatalf("AddAllowedDir: %v", err)
	}

	// The hand edit is still in the file...
	got := readLoaded(t, path)
	if got.Ball.Size != 60 {
		t.Errorf("the operator's [ball] size edit was reverted by the allowed_dirs write: got %d, want 60",
			got.Ball.Size)
	}
	if !strings.Contains(readText(t, path), "size = 60") {
		t.Errorf("config.toml no longer says size = 60 on its own lines:\n%s", readText(t, path))
	}
	// ...and the widening this answer asked for is in there too, so this case
	// cannot be green because nothing was written.
	if !reflect.DeepEqual(got.FS.AllowedDirs, []string{testDir}) {
		t.Errorf("fs.allowed_dirs on disk = %v, want [%s]", got.FS.AllowedDirs, testDir)
	}
	// Merging the file in must not smuggle its values into memory: the process
	// still holds what it loaded at startup (reading them back is ticket 223's
	// reload path, and a locked loosening there still needs ConfirmLocked).
	if m.Config().Ball.Size != 56 {
		t.Errorf("the guarded write applied the hand edit to memory too; this run now runs a "+
			"setting nobody confirmed: got %d, want the 56 this process loaded", m.Config().Ball.Size)
	}
	if !reflect.DeepEqual(m.Config().FS.AllowedDirs, []string{testDir}) {
		t.Errorf("memory fs.allowed_dirs = %v, want [%s]", m.Config().FS.AllowedDirs, testDir)
	}
}

// AC#1's positive control, spelled out as a permanent case rather than a
// mutation someone has to remember to run: handing the in-memory snapshot to
// SaveFile is what this write path did before ticket 226, and it does revert the
// hand edit. If this case ever goes green the wrong way (a snapshot write stops
// clobbering), TestAC1... above stops proving anything.
func TestAC1ControlSnapshotWriteIsWhatRevertsTheHandEdit(t *testing.T) {
	m, path, mutate := newTestManager(t, matrixBase)
	mutate(handEdited)

	// Exactly the old body: m.cur (the startup snapshot) -> SaveFile.
	if err := SaveFile(path, m.Config()); err != nil {
		t.Fatalf("SaveFile: %v", err)
	}
	if got := readLoaded(t, path); got.Ball.Size == 60 {
		t.Fatal("the control has lost its teeth: a snapshot write no longer reverts the hand " +
			"edit, so AC#1's green cannot be attributed to the merge")
	}
}

// AC#2: the merge must say which keys it wrote and which hand edits it kept.
func TestAC2GuardedWriteReportsTheKeysItWroteAndTheOnesItKept(t *testing.T) {
	logs := captureSlog(t)
	m, path, mutate := newTestManager(t, matrixBase)
	mutate(handEdited)

	if err := m.AddAllowedDir(testDir); err != nil {
		t.Fatalf("AddAllowedDir: %v", err)
	}
	text := logs.String()
	if !strings.Contains(text, "wrote=[fs.allowed_dirs]") {
		t.Errorf("the write did not report the keys it actually changed; log:\n%s", text)
	}
	if strings.Contains(text, "wrote=[ball.size") {
		t.Errorf("the write claimed a key it did not change; log:\n%s", text)
	}
	if !strings.Contains(text, "kept_in_file_not_in_memory=[ball.size]") {
		t.Errorf("the preserved hand edit was not named; log:\n%s", text)
	}
	if !strings.Contains(text, "level=WARN") {
		t.Errorf("diverging from the file must be louder than Info; log:\n%s", text)
	}
	// And it is a report, not a refusal: the widening landed.
	if got := readLoaded(t, path); !reflect.DeepEqual(got.FS.AllowedDirs, []string{testDir}) {
		t.Errorf("fs.allowed_dirs on disk = %v, want [%s]", got.FS.AllowedDirs, testDir)
	}
}

// AC#2's clean half: with nothing diverging, the report is one key at Info and
// says nothing was kept, because nothing was.
func TestAC2CleanWriteReportsOneKeyAtInfoLevel(t *testing.T) {
	logs := captureSlog(t)
	m, _, _ := newTestManager(t, matrixBase)
	if err := m.AddAllowedDir(testDir); err != nil {
		t.Fatalf("AddAllowedDir: %v", err)
	}
	text := logs.String()
	if !strings.Contains(text, "wrote=[fs.allowed_dirs]") || !strings.Contains(text, "level=INFO") {
		t.Errorf("a clean write must report the one key it wrote at Info; log:\n%s", text)
	}
	if strings.Contains(text, "kept_in_file_not_in_memory") {
		t.Errorf("a clean write reported kept hand edits; log:\n%s", text)
	}
}

// AC#2's other half: an unreadable file is refused, loudly, and nothing is
// written - the one shape in which "merge" cannot be computed, and the shape
// where overwriting would destroy the operator's edit outright.
func TestAC2UnreadableFileIsRefusedAndNothingIsWritten(t *testing.T) {
	m, path, mutate := newTestManager(t, matrixBase)
	mutate("this is not toml [[[")

	err := m.AddAllowedDir(testDir)
	if err == nil {
		t.Fatal("a guarded write accepted an unreadable config.toml and overwrote it")
	}
	if !strings.Contains(err.Error(), "cannot be read back") {
		t.Errorf("refusal must say why, got: %v", err)
	}
	if !strings.Contains(err.Error(), "keeping the previous allowlist in memory") {
		t.Errorf("refusal must say memory was left alone, got: %v", err)
	}
	if got := readText(t, path); got != "this is not toml [[[" {
		t.Errorf("the refused write still changed the file:\n%s", got)
	}
	if len(m.Config().FS.AllowedDirs) != 0 {
		t.Errorf("a refused write must roll the in-memory list back, got %v", m.Config().FS.AllowedDirs)
	}
}

// AC#3: statOwnWrite may claim only the content this write produced.
func TestAC3AdoptionClaimsOnlyWhatThisWriteProduced(t *testing.T) {
	// (a) Nothing diverged: the write is ours, so the poll sees no change...
	m, path, mutate := newTestManager(t, matrixBase)
	if err := m.AddAllowedDir(testDir); err != nil {
		t.Fatalf("AddAllowedDir: %v", err)
	}
	rep, err := m.CheckAndReload()
	if err != nil {
		t.Fatalf("CheckAndReload after a clean write: %v", err)
	}
	if rep != nil {
		t.Errorf("a write this process holds byte-for-byte was still read back as a change: %+v", rep)
	}
	// ...and a hand edit AFTER the write is still seen, which is AC#3's own test.
	time1 := readText(t, path)
	mutate(strings.Replace(time1, "size = 56", "size = 60", 1))
	rep, err = m.CheckAndReload()
	if err != nil {
		t.Fatalf("CheckAndReload after a later hand edit: %v", err)
	}
	if rep == nil || !slices.Contains(rep.Hot, "ball") {
		t.Fatalf("a hand edit made after our write was not seen (the claim never narrowed): %+v", rep)
	}

	// (b) Something diverged: the write merged foreign content into the file, so
	// it must NOT be claimed - the next poll has to read it back (and take the
	// D36 verdict on it, denied for a loosening).
	m2, _, mutate2 := newTestManager(t, matrixBase)
	mutate2(handEdited)
	if err := m2.AddAllowedDir(testDir); err != nil {
		t.Fatalf("AddAllowedDir (diverged): %v", err)
	}
	rep2, err := m2.CheckAndReload()
	if err != nil {
		t.Fatalf("CheckAndReload after a merged write: %v", err)
	}
	if rep2 == nil || !slices.Contains(rep2.Hot, "ball") {
		t.Fatalf("the merged write adopted a stat for content this process never read; "+
			"the hand edit is now invisible forever: %+v", rep2)
	}
	if rep2 != nil && len(rep2.Locked) != 0 {
		t.Errorf("a [ball] edit is hot-tier; a locked verdict was reported: %+v", rep2.Locked)
	}
}

// AC#3's foreign-loosening half: preserving a hand-added [fs] entry must not put
// it into memory behind ConfirmLocked's back.
func TestAC3HandAddedEntryIsPreservedButNotAppliedToMemory(t *testing.T) {
	m, path, mutate := newTestManager(t, matrixBase)
	mutate(strings.Replace(matrixBase, "allowed_dirs = []", `allowed_dirs = ["E:\\hand"]`, 1))

	if err := m.AddAllowedDir(testDir); err != nil {
		t.Fatalf("AddAllowedDir: %v", err)
	}
	// Both lines are in the file: nobody's entry was overwritten...
	got := readLoaded(t, path)
	if !reflect.DeepEqual(got.FS.AllowedDirs, []string{`E:\hand`, testDir}) {
		t.Errorf("fs.allowed_dirs on disk = %v, want the hand entry plus the confirmed one", got.FS.AllowedDirs)
	}
	// ...and only the confirmed one is in memory.
	if !reflect.DeepEqual(m.Config().FS.AllowedDirs, []string{testDir}) {
		t.Errorf("memory fs.allowed_dirs = %v; the guarded write must not adopt a key it never "+
			"confirmed (D36 rule 1 belongs to the reload path)", m.Config().FS.AllowedDirs)
	}
	// The file is not claimed as our own write, so the reload path still runs its
	// verdict on the hand entry - and denies it while nothing is wired.
	m.ConfirmLocked = func(string, []string) bool { return false }
	rep, err := m.CheckAndReload()
	if err != nil {
		t.Fatalf("CheckAndReload: %v", err)
	}
	if rep == nil || len(rep.Locked) == 0 {
		t.Fatal("the hand-added [fs] entry was never judged; our write hid it (AC#3)")
	}
	if len(m.Config().FS.AllowedDirs) != 1 || m.Config().FS.AllowedDirs[0] != testDir {
		t.Errorf("a denied loosening must keep memory at the confirmed entry, got %v",
			m.Config().FS.AllowedDirs)
	}
}

// AC#4: the older twin of the same shape - SetPermissionMode - keeps hand edits
// too, and still adopts its own write when nothing diverged.
func TestAC4PermissionModeWriteKeepsAHandEditedKey(t *testing.T) {
	m, path, mutate := newTestManager(t, matrixBase)
	mutate(handEdited)

	if err := m.SetPermissionMode(risk.ModeAutoApprove); err != nil {
		t.Fatalf("SetPermissionMode: %v", err)
	}
	got := readLoaded(t, path)
	if got.Ball.Size != 60 {
		t.Errorf("the mode switch reverted the operator's [ball] size edit: got %d, want 60", got.Ball.Size)
	}
	if got.Risk.PermissionMode != risk.ModeAutoApproveName {
		t.Errorf("risk.permission_mode on disk = %q, want %q",
			got.Risk.PermissionMode, risk.ModeAutoApproveName)
	}
	if m.Config().Ball.Size != 56 {
		t.Errorf("the mode switch must not apply a foreign hand edit to memory, got %d",
			m.Config().Ball.Size)
	}

	// The clean half: with nothing diverging the switch is still claimed as our
	// own write (that is what R20/M3's persistence means for CheckAndReload).
	m2, _, _ := newTestManager(t, matrixBase)
	if err := m2.SetPermissionMode(risk.ModeAutoApprove); err != nil {
		t.Fatalf("SetPermissionMode (clean): %v", err)
	}
	rep, err := m2.CheckAndReload()
	if err != nil {
		t.Fatalf("CheckAndReload: %v", err)
	}
	if rep != nil {
		t.Errorf("a clean mode switch was read back as somebody else's edit: %+v", rep)
	}
}

// AC#5's落库 half: revoking a long-lived rule is a write of the same single key,
// so it persists and it persists WITHOUT touching anything else in the file.
// (The other half of AC#5 - a reply grammar that triggers this - is not in this
// ticket's reach and is reported as open, see the evidence file.)
func TestAC5SetAllowedDirsPersistsARevocationWithoutAClobber(t *testing.T) {
	start := strings.Replace(matrixBase, "allowed_dirs = []", `allowed_dirs = ["D:\\data", "D:\\notes"]`, 1)
	m, path, mutate := newTestManager(t, start)
	mutate(strings.Replace(handEdited, "allowed_dirs = []", `allowed_dirs = ["D:\\data", "D:\\notes"]`, 1))

	if err := m.SetAllowedDirs([]string{`D:\notes`}); err != nil {
		t.Fatalf("SetAllowedDirs: %v", err)
	}
	got := readLoaded(t, path)
	if !reflect.DeepEqual(got.FS.AllowedDirs, []string{`D:\notes`}) {
		t.Errorf("the revoked entry is still in the file: %v", got.FS.AllowedDirs)
	}
	if !reflect.DeepEqual(m.Config().FS.AllowedDirs, []string{`D:\notes`}) {
		t.Errorf("memory and file diverged after a successful revoke: %v", m.Config().FS.AllowedDirs)
	}
	if got.Ball.Size != 60 {
		t.Errorf("the revocation reverted the hand edit: got %d, want 60", got.Ball.Size)
	}
}

// One pinned safety behaviour, untouched by this ticket: a relative hop is not a
// root and is refused before anything is read or written.
func TestAddAllowedDirStillRefusesARelativeHop(t *testing.T) {
	m, path, mutate := newTestManager(t, matrixBase)
	mutate(handEdited)
	before := readText(t, path)

	for _, bad := range []string{`D:\data\..\escape`, `..\who`} {
		err := m.AddAllowedDir(bad)
		if err == nil {
			t.Fatalf("an allowed_dirs entry with a relative hop was accepted: %q", bad)
		}
		if !strings.Contains(err.Error(), "..") {
			t.Errorf("refusal of %q must name the hop, got: %v", bad, err)
		}
	}
	if got := readText(t, path); got != before {
		t.Errorf("a refused entry changed the file:\n%s", got)
	}
	if len(m.Config().FS.AllowedDirs) != 0 {
		t.Errorf("a refused entry changed memory: %v", m.Config().FS.AllowedDirs)
	}
	if err := m.SetAllowedDirs([]string{`D:\data`, `..\escape`}); err == nil {
		t.Error("SetAllowedDirs accepted a relative hop")
	}
	if len(m.Config().FS.AllowedDirs) != 0 {
		t.Errorf("SetAllowedDirs changed memory before refusing: %v", m.Config().FS.AllowedDirs)
	}
}

// The other pinned safety behaviour: a failed write cannot leave memory holding
// what the file does not.
func TestGuardedWriteFailureRollsBackMemory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(matrixBase), 0o600); err != nil {
		t.Fatal(err)
	}
	m, err := NewManager(path, nil)
	if err != nil {
		t.Fatal(err)
	}
	// Nowhere to write: the parent directory does not exist, so atomicWrite's
	// temp file cannot be created. This is the real write-failure path, not a
	// stand-in for it.
	m.path = filepath.Join(dir, "gone", "config.toml")

	if err := m.AddAllowedDir(testDir); err == nil {
		t.Fatal("a write into a missing directory reported success")
	} else if !strings.Contains(err.Error(), "keeping the previous allowlist in memory") {
		t.Errorf("the failure must say memory was rolled back, got: %v", err)
	}
	if len(m.Config().FS.AllowedDirs) != 0 {
		t.Errorf("a failed write left the widening in memory: %v", m.Config().FS.AllowedDirs)
	}
	if err := m.SetPermissionMode(risk.ModeAutoApprove); err == nil {
		t.Fatal("a mode write into a missing directory reported success")
	}
	if m.Config().Risk.PermissionMode == risk.ModeAutoApproveName {
		t.Error("a failed mode write left the new mode in memory")
	}
}

// A file that already says what this write would put in it is not rewritten at
// all: no clobber-by-canonicalization, no lost comments, for nothing.
func TestGuardedWriteDoesNotRewriteAFileThatAlreadySaysIt(t *testing.T) {
	m, path, mutate := newTestManager(t, matrixBase)
	hand := strings.Replace(matrixBase, "allowed_dirs = []", `allowed_dirs = ["D:\\data"]`, 1)
	mutate(hand)
	before := readText(t, path)
	mtimeBefore := statMod(t, path)

	if err := m.AddAllowedDir(testDir); err != nil {
		t.Fatalf("AddAllowedDir: %v", err)
	}
	if got := readText(t, path); got != before {
		t.Errorf("the file already carried the entry and was rewritten anyway:\n%s", got)
	}
	if !statMod(t, path).Equal(mtimeBefore) {
		t.Error("the file's mtime moved on a write that wrote nothing")
	}
	if !reflect.DeepEqual(m.Config().FS.AllowedDirs, []string{testDir}) {
		t.Errorf("memory must still hold the confirmed entry, got %v", m.Config().FS.AllowedDirs)
	}
}

// The differ AC#2 reports through: a reported path must be a real key in the
// file, per map element, and must not report what did not change.
func TestDiffKeyPathsNamesRealKeyPaths(t *testing.T) {
	a := NewDefaults()
	b := deepCopyConfig(a)
	if got := diffKeyPaths(a, b); len(got) != 0 {
		t.Fatalf("identical configs differ at %v", got)
	}
	b.Ball.Size = a.Ball.Size + 4
	b.LLM.Providers = map[string]Provider{"acme": {BaseURL: "https://example.invalid/v1"}}
	b.Plugins.Entries = map[string]PluginEntry{"p1": {Enabled: true}}
	b.App.Portable = !a.App.Portable // `toml:"-"`: the file never holds it

	got := diffKeyPaths(a, b)
	want := []string{"ball.size", "llm.providers.acme", "plugins.p1"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("diffKeyPaths = %v, want %v", got, want)
	}
	// schema_version is forced by SaveFile and is not an operator's key.
	c := deepCopyConfig(a)
	c.SchemaVersion = 1
	if got := diffKeyPaths(a, c); len(got) != 0 {
		t.Errorf("schema_version was reported as a changed key: %v", got)
	}
}

func readText(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func statMod(t *testing.T, path string) time.Time {
	t.Helper()
	st, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	return st.ModTime()
}
