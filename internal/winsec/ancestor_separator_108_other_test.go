//go:build !windows

// Ticket 108 AC#2's platform-correctness half, on the only platform where a
// backslash is an ordinary character inside a file name.
//
// The tempting "fix" for the Windows separator bypass (P2, R-103-2) is to treat
// both / and \ as separators everywhere. On POSIX that is wrong in both
// directions and these four cases are the pins:
//
//   - folding \ into a separator splits one real name into two, so a prefix that
//     is not an ancestor of the input gets Lstat'ed: a symlink called "a" then
//     refuses a spelling whose actual ancestor is the directory "a\b" (a false
//     refusal, i.e. the reclaim route broken for no reason);
//   - rejoining the pieces on filepath.Join collapses the two names into one, so
//     the ancestor that IS a link, "x\y", is never looked at under its real name
//     and the unlink lands in somebody else's tree - the fail-open this
//     repository already booked once as A74(3).
//
// Run for real (not compile-only) per ledger A79(1): see the ticket's AC#5 row
// for the container command and its rc.
package winsec_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/CarlosShao/wisp/internal/winsec"
)

// foreignTree108 is somebody else's directory with one file in it, plus the
// marker of what "nothing happened" has to mean on this platform: the file is
// still there, and its mode is still the wide one it was created with.
func foreignTree108(t *testing.T, name string) (dir, victim string) {
	t.Helper()
	dir = filepath.Join(t.TempDir(), name)
	if err := os.Mkdir(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	victim = filepath.Join(dir, "keep-me.txt")
	if err := os.WriteFile(victim, []byte("not this tree's data"), 0o666&^fs.FileMode(0o077)); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(victim, 0o666); err != nil {
		t.Fatal(err)
	}
	return dir, victim
}

func assertStillThere108(t *testing.T, path string, why string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Errorf("%s: %s is gone (%v)", why, path, err)
		return
	}
	if perm := info.Mode().Perm(); perm != 0o666 {
		t.Errorf("%s: %s mode changed from 0666 to %o, so the seal/removal acted on the wrong tree", why, path, perm)
	}
}

// TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink is the positive leg on
// this platform: an ancestor that is a symlink must refuse the unlink even though
// the spelling is entirely forward slashes.
func TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink(t *testing.T) {
	_, victim := foreignTree108(t, "foreign")
	root := filepath.Join(t.TempDir(), "root")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "link")
	if err := os.Symlink(filepath.Dir(victim), link); err != nil {
		t.Skipf("no symlink privilege on this host: %v", err)
	}
	spelled := filepath.ToSlash(filepath.Join(link, "keep-me.txt"))
	err := winsec.RemoveUnlinked(spelled)
	t.Logf("RemoveUnlinked(%q) -> err=%v", spelled, err)
	if err == nil {
		t.Errorf("AC#2 RED: a slash-spelled path through a symlink returned nil")
	} else if !errors.Is(err, winsec.ErrIsReparsePoint) {
		t.Errorf("AC#2: refusal must name ErrIsReparsePoint, got %v", err)
	}
	assertStillThere108(t, victim, "AC#2 RED the foreign file was deleted through the symlink")
}

// TestAC2POSIXDoesNotFoldABackslashIntoASeparator is AC#2's required reverse
// case: on POSIX the name "a\b" is one component, not two. There is a symlink at
// root/a and a real directory at root/"a\b"; the input names the real directory,
// so the unlink must succeed and root/a's target must be untouched. An
// implementation that treats backslash as a separator refuses (or worse, rejoins
// and deletes) here.
func TestAC2POSIXDoesNotFoldABackslashIntoASeparator(t *testing.T) {
	foreignDir, victim := foreignTree108(t, "foreign")
	root := filepath.Join(t.TempDir(), "root")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(foreignDir, filepath.Join(root, "a")); err != nil {
		t.Skipf("no symlink privilege on this host: %v", err)
	}
	realDirName := `a\b`
	if err := os.Mkdir(filepath.Join(root, realDirName), 0o700); err != nil {
		t.Fatal(err)
	}
	stray := filepath.Join(root, realDirName, "keep-me.txt")
	if err := os.WriteFile(stray, []byte("this tree's own stray"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := winsec.RemoveUnlinked(stray); err != nil {
		t.Errorf("AC#2 RED: POSIX folded the backslash in %q into a separator and refused the tree's own file: %v", stray, err)
	}
	if _, err := os.Stat(stray); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("the stray inside root/%s should have been removed: %v", realDirName, err)
	}
	assertStillThere108(t, victim, "AC#2 RED the symlink target was reached through a folded name")
}

// TestAC2POSIXABackslashInALinkNameIsStillALinkAncestor is the fail-open pin: the
// ancestor whose *own name* contains a backslash is a symlink, so an
// implementation that splits on backslash and rebuilds prefixes with Join looks
// at root/x/y (which does not exist), concludes "not a link", and unlinks inside
// the foreign tree.
func TestAC2POSIXABackslashInALinkNameIsStillALinkAncestor(t *testing.T) {
	foreignDir, victim := foreignTree108(t, "foreign")
	root := filepath.Join(t.TempDir(), "root")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatal(err)
	}
	linkName := `x\y`
	if err := os.Symlink(foreignDir, filepath.Join(root, linkName)); err != nil {
		t.Skipf("no symlink privilege on this host: %v", err)
	}
	spelled := filepath.Join(root, linkName, "keep-me.txt")
	err := winsec.RemoveUnlinked(spelled)
	t.Logf("RemoveUnlinked(%q) -> err=%v", spelled, err)
	if err == nil {
		t.Errorf("AC#2 RED: the link ancestor %q was not checked because its name was split on the backslash", linkName)
	} else if !errors.Is(err, winsec.ErrIsReparsePoint) {
		t.Errorf("AC#2: refusal must name ErrIsReparsePoint, got %v", err)
	}
	assertStillThere108(t, victim, "AC#2 RED the foreign file was deleted through a backslash-named link")
}

// TestAC4POSIXFloorAnswersInsideTheNamedTree is AC#4's non-Windows half. The
// platform resolver here is the stub in internal/risk/pathresolver_other.go,
// whose Resolved flag is false on every answer, so nothing about this leg may
// depend on Windows running: with no resolver installed the floor has to keep
// the answer inside the tree the caller named, and has to refuse a spelling that
// does not.
func TestAC4POSIXFloorAnswersInsideTheNamedTree(t *testing.T) {
	if prev := winsec.PathResolverInstalled(); prev != nil {
		t.Skipf("this platform's binary links a real C26 pipeline (seam holds %T), so the floor leg measures nothing", prev)
	}
	root := filepath.Join(t.TempDir(), "data")
	dir := filepath.Join(root, "sub")
	if err := winsec.PrivateDirAll(dir, 0o700); err != nil {
		t.Fatalf("PrivateDirAll(%s): %v", dir, err)
	}
	got, err := winsec.ResolvePath(dir)
	if err != nil {
		t.Fatalf("ResolvePath(%s): %v", dir, err)
	}
	if got.String() != dir {
		t.Errorf("AC#4 RED: the floor answered %q for %q, i.e. it rewrote the spelling it was asked about", got.String(), dir)
	}
	if filepath.Dir(got.String()) != root {
		t.Errorf("AC#4 RED: the answer %q is not inside the named tree %q", got.String(), root)
	}
	info, statErr := os.Stat(dir)
	if statErr != nil {
		t.Errorf("AC#4: nothing was created where the call named it: %v", statErr)
	} else if perm := info.Mode().Perm(); perm != 0o700 {
		t.Errorf("AC#4: %s landed %o, not 0700", dir, perm)
	}
	if _, err := winsec.ResolvePath(root + "/../" + "elsewhere"); err == nil {
		t.Errorf("AC#4 RED: the floor accepted a spelling that traverses a parent pointer")
	} else {
		t.Logf("parent-pointer spelling refused as: %v", err)
	}
}
