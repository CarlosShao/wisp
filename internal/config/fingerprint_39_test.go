package config

// Ticket 39 AC#4's teeth, config layer (2026-11-18):
//
//   - FileSHA256 answers the disk bytes' real hash, and "no file" is
//     "no fingerprint", never a hash of an empty read;
//   - the no-op save (the write's value is the value the file already
//     carries) leaves config.toml BYTE-IDENTICAL: same sha256, same mtime,
//     comments and blank lines and hand-chosen order all still there. The
//     tooth bites in both directions: a real change of the same key moves
//     the bytes (the harness would catch a no-op that went through a real
//     rewrite - the sha would differ and this test would go red).
//
// The judgement is on the FILE, never on a return value: writeguard.go's
// mergeWrite answers "already carries this key" by not writing at all, and
// these cases would catch the day that early-return is lost to a rewrite.

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strings"
	"testing"
)

// noopBase39 is a config with hand-written comments, blank lines and a
// non-canonical key order - exactly the formatting a no-op save must never
// re-lay (mergeWrite's canonical rewrite only happens on a real change).
const noopBase39 = `# 手写注释：这份文件是人改过的，注释必须活得过任何 no-op 保存。
schema_version = 2

[ball]
# 球的大小，手写的注释行
size = 90

[app]
theme = "dark"
`

func Test39FileSHA256AnswersTheDiskBytesHash(t *testing.T) {
	m, path, _ := newTestManager(t, noopBase39)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read the fixture back: %v", err)
	}
	sum := sha256.Sum256(raw)
	want := hex.EncodeToString(sum[:])
	got, ferr := m.FileSHA256()
	if ferr != nil {
		t.Fatalf("FileSHA256: %v", ferr)
	}
	if got != want {
		t.Errorf("FileSHA256 = %s, want the disk bytes' own sha256 %s", got, want)
	}

	// The missing-file shape: no file is no fingerprint (""), not an error and
	// not a hash of an empty read.
	if err := os.Remove(path); err != nil {
		t.Fatalf("remove the fixture: %v", err)
	}
	got2, ferr2 := m.FileSHA256()
	if ferr2 != nil {
		t.Errorf("FileSHA256 over a missing file returned %v; want (\"\", nil)", ferr2)
	}
	if got2 != "" {
		t.Errorf("FileSHA256 over a missing file = %q, want the empty no-fingerprint answer", got2)
	}
}

func Test39NoopSectionWriteLeavesBytesAndMtimeUntouched(t *testing.T) {
	m, path, _ := newTestManager(t, noopBase39)

	beforeBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read before: %v", err)
	}
	beforeStat, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat before: %v", err)
	}
	beforeSHA, err := m.FileSHA256()
	if err != nil {
		t.Fatalf("hash before: %v", err)
	}

	// The no-op save: ball.size is already 90 on disk, and the write carries
	// 90 again. Written must come back empty - the receipt says "nothing
	// changed" - and the file must not move at all.
	written, werr := m.SetBallSize(90)
	if werr != nil {
		t.Fatalf("the no-op save was refused: %v", werr)
	}
	if len(written) != 0 {
		t.Errorf("the no-op save's receipt named key paths %v, want an empty list (nothing changed)", written)
	}

	afterStat, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat after: %v", err)
	}
	afterBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read after: %v", err)
	}
	afterSHA, err := m.FileSHA256()
	if err != nil {
		t.Fatalf("hash after: %v", err)
	}

	if !beforeStat.ModTime().Equal(afterStat.ModTime()) {
		t.Errorf("mtime moved across a no-op save (%v -> %v): the file was rewritten", beforeStat.ModTime(), afterStat.ModTime())
	}
	if afterSHA != beforeSHA {
		t.Errorf("sha256 moved across a no-op save (%s -> %s): the file's bytes changed", beforeSHA, afterSHA)
	}
	if string(afterBytes) != string(beforeBytes) {
		t.Errorf("a no-op save re-laid the file:\nbefore:\n%s\nafter:\n%s", beforeBytes, afterBytes)
	}
	for _, want := range []string{"# 手写注释：这份文件是人改过的", "# 球的大小，手写的注释行"} {
		if !strings.Contains(string(afterBytes), want) {
			t.Errorf("a no-op save dropped the hand-written comment %q:\n%s", want, afterBytes)
		}
	}

	// The positive control on the SAME key: a real change must move the bytes.
	// Without this arm the test above could stay green on a harness that never
	// writes anything at all.
	if _, err := m.SetBallSize(120); err != nil {
		t.Fatalf("the real change was refused: %v", err)
	}
	changedSHA, err := m.FileSHA256()
	if err != nil {
		t.Fatalf("hash after the real change: %v", err)
	}
	if changedSHA == beforeSHA {
		t.Error("a real value change left the sha256 untouched: the harness cannot tell a write from a no-op")
	}
}
