package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Ticket 75, consumer side. PathCanonicalizer hands two things to the rest of
// the bridge: a string the fs tools open, and the comparison form that decides
// R2 (in-allowlist or not). Both used to be built on a Windows assumption, so
// on Linux Canonicalize returned "<cwd>/\tmp\u001\a.txt" - a string no OS call
// can open, and one whose separators had been folded into a character that is
// legal inside a POSIX file name.
//
// No build tag: "platform-shaped" is assertable on every platform, and the
// Windows-only guard is precisely the blind spot the defect lived in.

// foreignSepT is the separator that is not this platform's.
func foreignSepT() string {
	if filepath.Separator == '\\' {
		return "/"
	}
	return `\`
}

func TestCanonicalizeReturnsAPathTheOSCanOpen(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "payload.txt")
	if err := os.WriteFile(target, []byte("ticket 75"), 0o600); err != nil {
		t.Fatal(err)
	}

	pc := NewPathCanonicalizer([]string{dir}, nil)
	if bad := pc.UnusableRoots(); len(bad) != 0 {
		t.Fatalf("the allowlist root was rejected: %v", bad)
	}
	canon, err := pc.Canonicalize(target)
	if err != nil {
		t.Fatalf("Canonicalize(%q): %v", target, err)
	}
	if strings.Contains(canon, foreignSepT()) {
		t.Errorf("canonical %q carries the foreign separator %q", canon, foreignSepT())
	}
	if _, err := os.Lstat(canon); err != nil {
		t.Fatalf("fs.read/fs.write open what Canonicalize returned, and that path is not there: %v (%q)",
			err, canon)
	}
	if !pc.InAllowlist(canon) {
		t.Fatalf("InAllowlist(%q) = false for a file inside the only allowed root %q: R2 would send "+
			"every plain write to L2", canon, dir)
	}
}

// TestCanonicalizeAgreesWithTheOSName is the cheap version of "the verdict and
// the bytes describe the same file": after one round trip the canonical must
// still name the file the caller asked about.
func TestCanonicalizeAgreesWithTheOSName(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "round-trip.txt")
	if err := os.WriteFile(target, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	canon, err := NewPathCanonicalizer(nil, nil).Canonicalize(target)
	if err != nil {
		t.Fatalf("Canonicalize: %v", err)
	}
	body, err := os.ReadFile(canon)
	if err != nil || string(body) != "x" {
		t.Fatalf("reading back the canonical %q gave %v %q, want the bytes of %q",
			canon, err, body, target)
	}
}

func TestFoldPathKeepsPosixBackslashesDistinct(t *testing.T) {
	// Guard, not t.Skip: the property is about a character that is only a
	// filename byte on POSIX, and the package's SKIP ledger must stay at zero
	// so no red can hide behind a skip (ticket 75 AC#5).
	if filepath.Separator == '\\' {
		return
	}
	// Folding '/' into '\' used to make these two POSIX paths one comparison
	// key. foldPath's output authorizes roots, so that merge is a fail-open.
	a := foldPath("/home/u/we/ird")
	b := foldPath(`/home/u/we\ird`)
	if a == b {
		t.Fatalf("two different POSIX paths fold to the same allowlist key %q", a)
	}
	if foldPath("/Home/U/x") != foldPath("/home/u/x") {
		t.Error("the case fold is part of the contract and must stay")
	}
}

func TestDirOfBaseOfCutOnThePlatformSeparator(t *testing.T) {
	dir := t.TempDir()
	child := filepath.Join(dir, "notes.txt")
	if got := dirOf(child); got != dir {
		t.Errorf("dirOf(%q) = %q, want %q", child, got, dir)
	}
	if got := baseOf(child); got != "notes.txt" {
		t.Errorf("baseOf(%q) = %q, want notes.txt", child, got)
	}
	if strings.Contains(dirOf(child), foreignSepT()) || strings.Contains(baseOf(child), foreignSepT()) {
		t.Errorf("dirOf/baseOf leaked the foreign separator into the audit text: %q / %q",
			dirOf(child), baseOf(child))
	}
	if filepath.Separator != '\\' {
		// A POSIX file name may contain a backslash; cutting must ignore it.
		p := filepath.Join(dir, `we\ird.txt`)
		if got := baseOf(p); got != `we\ird.txt` {
			t.Errorf("baseOf(%q) = %q, want the whole name", p, got)
		}
		if got := dirOf(p); got != dir {
			t.Errorf("dirOf(%q) = %q, want %q", p, got, dir)
		}
	}
}
