//go:build !windows

// Ticket 113 AC#1/AC#2/AC#3, on the platform whose floor used to answer every
// placement question with "yes, unchanged": winsec_other.go's
// platformVerifyPlacement was `return path, nil`, so the sealing entry points
// had no link leg at all here. PROBE P3 (ticket 103) and its Windows re-cut
// (ticket 108) measured the same outcome one platform away: SealFile on a
// spelling that reaches its target through a symlink returned nil while the
// chmod landed on somebody else's file.
//
// Every refusal case below therefore measures two things, and the second one is
// the judgement: an error is only worth anything if the foreign tree is standing
// exactly where it was, mode and owner both unchanged. The plain-path cases are
// the reverse half (AC#3) - a guard that refused everything would pass the
// first group and be worthless, so "the tree's own file still gets sealed" is
// asserted with the same instrument.
//
// The separator rule under test is pathPieces', which on POSIX cuts on '/'
// alone. Two of these cases exist only to keep that rule honest from the seal
// side: a backslash is an ordinary character in a file name here, so folding it
// in would either refuse the real directory named `a\b` (false refusal) or check
// a rebuilt `x/y` while the link is `x\y` (fail-open, booked as A74(3)).
//
// Run for real, not compile-only, in a Linux container: see the ticket's AC#1
// row for the mount and the rc.
package winsec_test

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/CarlosShao/wisp/internal/winsec"
)

// foreign113 is somebody else's tree: one directory with one file in it, both
// deliberately wide (0777 / 0666), because "wide" is what makes a seal that
// wandered in visible as a mode change.
type foreign113 struct {
	dir        string
	victim     string
	dirFact    string
	victimFact string
}

func newForeign113(t *testing.T, name string) *foreign113 {
	t.Helper()
	f := &foreign113{}
	f.dir = filepath.Join(winsec.SealableTempDirForTest124(t), name)
	if err := os.Mkdir(f.dir, 0o700); err != nil {
		t.Fatal(err)
	}
	f.victim = filepath.Join(f.dir, "keep-me.txt")
	if err := os.WriteFile(f.victim, []byte("not this tree's data"), 0o600); err != nil {
		t.Fatal(err)
	}
	// Chmod rather than a wide umask: the point is a mode the kernel has agreed
	// to store, so a later read-back of 0600 can only mean something changed it.
	if err := os.Chmod(f.victim, 0o666); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(f.dir, 0o777); err != nil {
		t.Fatal(err)
	}
	f.dirFact, f.victimFact = statFact113(f.dir), statFact113(f.victim)
	t.Logf("AC#1 foreign tree before: dir %s = %s, victim %s = %s", f.dir, f.dirFact, f.victim, f.victimFact)
	return f
}

// untouched asserts that nothing in the foreign tree changed: still there, same
// mode, same owner. Owner is measured because a seal that reached the wrong tree
// could also be a chown, and because "I only chmod'ed it" is not a claim this
// file should take on faith.
func (f *foreign113) untouched(t *testing.T, why string) {
	t.Helper()
	for _, e := range []struct{ path, want string }{{f.dir, f.dirFact}, {f.victim, f.victimFact}} {
		got := statFact113(e.path)
		t.Logf("%s: %s before=%s after=%s", why, e.path, e.want, got)
		switch {
		case got == "missing":
			t.Errorf("AC#1 RED %s: %s is gone", why, e.path)
		case got != e.want:
			t.Errorf("AC#1 RED %s: %s was %s and is now %s, so the call acted on the foreign tree",
				why, e.path, e.want, got)
		}
	}
}

// statFact113 is mode + owner in one comparable string, read with os.Stat (the
// following kind of stat, which is the whole reason a link in the ancestor chain
// matters: os.Chmod goes through it).
func statFact113(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return "missing"
	}
	uid, gid := -1, -1
	if st, ok := info.Sys().(*syscall.Stat_t); ok {
		uid, gid = int(st.Uid), int(st.Gid)
	}
	return fmt.Sprintf("mode=%o uid=%d gid=%d", info.Mode().Perm(), uid, gid)
}

// linkTo113 puts a symlink at the named place inside root, pointing at target,
// and returns the spelling the caller will hand to a sealing entry point.
func linkTo113(t *testing.T, root, linkName, target string) string {
	t.Helper()
	link := filepath.Join(root, linkName)
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("cannot build the symlink this case measures (%v): without it the leg is untested", err)
	}
	info, err := os.Lstat(link)
	if err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("Lstat(%s) after os.Symlink: info=%v err=%v", link, info, err)
	}
	return link
}

// assertRefused113 is the "either refuse, or only touch the link itself" half of
// AC#1. This floor can only refuse (it never rewrites), so refusal is the
// outcome this ticket's fix produces and the outcome every case requires; a nil
// here means the seal walked through the link.
func assertRefused113(t *testing.T, what, spelled string, err error) {
	t.Helper()
	t.Logf("AC#1 %s(%q) -> err=%v", what, spelled, err)
	if err == nil {
		t.Errorf("AC#1 RED: %s(%q) returned nil, i.e. it sealed through a symlink and reported success. AC#1 requires either a refusal or an action confined to the link itself.",
			what, spelled)
		return
	}
	if !errors.Is(err, winsec.ErrUnresolvedPath) {
		t.Errorf("AC#1: %s(%q) refused with %v, which does not name ErrUnresolvedPath: a refusal this call did not make cannot be attributed to the floor",
			what, spelled, err)
	}
}

// TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone is the
// outcome PROBE P3 reported, cut on this platform: one symlink standing in the
// middle of the spelling.
func TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone(t *testing.T) {
	f := newForeign113(t, "foreign")
	root := filepath.Join(t.TempDir(), "root")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatal(err)
	}
	link := linkTo113(t, root, "link", f.dir)
	spelled := filepath.Join(link, "keep-me.txt")

	assertRefused113(t, "SealFile", spelled, winsec.SealFile(spelled))

	f.untouched(t, "AC#1 the foreign tree behind the link")
	if info, err := os.Lstat(link); err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Errorf("AC#1: a refusal must not touch anything, but the link now reads %v (%v)", info, err)
	}
}

// TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth is the mutation bait for
// AC#4's "half-fix" shape: an implementation that inspects one ancestor (or only
// the parent) passes depth 1 and fails here, because the link sits two, three
// and four components below the root the caller named.
func TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth(t *testing.T) {
	for _, depth := range []int{1, 2, 3, 4} {
		t.Run(fmt.Sprintf("link-at-depth-%d", depth), func(t *testing.T) {
			f := newForeign113(t, "foreign")
			root := filepath.Join(t.TempDir(), "root")
			if err := os.Mkdir(root, 0o700); err != nil {
				t.Fatal(err)
			}
			// Walk down depth-1 real directories, then the link.
			place := root
			for i := 1; i < depth; i++ {
				next := filepath.Join(place, fmt.Sprintf("dir%d", i))
				if err := os.Mkdir(next, 0o700); err != nil {
					t.Fatal(err)
				}
				place = next
			}
			link := linkTo113(t, place, "link", f.dir)
			spelled := filepath.Join(link, "keep-me.txt")
			t.Logf("AC#1 depth %d: spelled=%s (pieces below root: %d)", depth, spelled, depth+1)

			assertRefused113(t, "SealFile", spelled, winsec.SealFile(spelled))
			f.untouched(t, fmt.Sprintf("AC#1 depth %d", depth))
		})
	}
}

// TestAC1POSIXSealDirThroughASymlinkRefusesIsTheSameLegOnTheDirectoryEntry:
// SealDir is the other exported way in, and on POSIX it chmods a directory
// through a link exactly the same way.
func TestAC1POSIXSealDirThroughASymlinkRefuses(t *testing.T) {
	f := newForeign113(t, "foreign")
	root := filepath.Join(t.TempDir(), "root")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatal(err)
	}
	link := linkTo113(t, root, "link", f.dir)

	assertRefused113(t, "SealDir", link, winsec.SealDir(link))
	f.untouched(t, "AC#1 the foreign directory behind the link")
}

// TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing covers the entry
// point that creates bytes: refusing after the open would already have written
// into the foreign tree, so this measures that the placement leg runs before any
// content exists on the other side.
func TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing(t *testing.T) {
	f := newForeign113(t, "foreign")
	root := filepath.Join(t.TempDir(), "root")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatal(err)
	}
	link := linkTo113(t, root, "link", f.dir)
	spelled := filepath.Join(link, "secret.txt")

	assertRefused113(t, "PrivateFile", spelled, winsec.PrivateFile(spelled, []byte("top secret"), 0o600))

	f.untouched(t, "AC#1 the foreign tree behind the link")
	if _, err := os.Lstat(filepath.Join(f.dir, "secret.txt")); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("AC#1 RED: PrivateFile put bytes in the foreign tree (%v)", err)
	}
	if _, err := os.Lstat(spelled); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("AC#1: the refusal created %q anyway (%v)", spelled, err)
	}
}

// TestAC1POSIXSealFileThroughABackslashNamedLink is AC#2's fail-open pin seen
// from the seal side: the ancestor's own name contains a backslash, which POSIX
// reads as an ordinary character. An implementation that cuts on the backslash
// too (or rebuilds prefixes with filepath.Join) ends up Lstat'ing `root/x/y`,
// finding nothing, concluding "not a link" and sealing through `root/x\y`.
func TestAC1POSIXSealFileThroughABackslashNamedLink(t *testing.T) {
	f := newForeign113(t, "foreign")
	root := filepath.Join(t.TempDir(), "root")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatal(err)
	}
	linkTo113(t, root, `x\y`, f.dir)
	spelled := filepath.Join(root, `x\y`, "keep-me.txt")

	assertRefused113(t, "SealFile", spelled, winsec.SealFile(spelled))
	f.untouched(t, "AC#2 the foreign tree behind a backslash-named link")
}

// TestAC3POSIXSealFileStillNarrowsAPlainFileInsideTheNamedTree is AC#3's reverse
// half: the target really is inside the named tree and no ancestor is a link, so
// the call has to keep working. Without this case a "fix" that refuses
// everything would read as green.
func TestAC3POSIXSealFileStillNarrowsAPlainFileInsideTheNamedTree(t *testing.T) {
	root := filepath.Join(winsec.SealableTempDirForTest124(t), "data")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatal(err)
	}
	plain := filepath.Join(root, "blob.bin")
	if err := os.WriteFile(plain, []byte("ours"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(plain, 0o666); err != nil {
		t.Fatal(err)
	}
	before := statFact113(plain)
	t.Logf("AC#3 before: %s = %s", plain, before)

	if err := winsec.SealFile(plain); err != nil {
		t.Errorf("AC#3 RED: a plain file inside the named tree was refused: %v", err)
	}
	after := statFact113(plain)
	t.Logf("AC#3 after: %s = %s", plain, after)
	info, err := os.Stat(plain)
	if err != nil {
		t.Fatalf("AC#3: %s is gone: %v", plain, err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Errorf("AC#3 RED: the seal did not narrow anything, mode is %o (before %s)", got, before)
	}
}

// TestAC3POSIXSealStillWorksNextToAndThroughRealDirectoriesAndLinks covers the
// two ways a link-leg can over-refuse: a sibling link that is not an ancestor,
// and a link below the file being sealed. It also walks the other entry points
// (PrivateDirAll, PrivateFile, SealDir) so the leg cannot be bolted onto
// SealFile alone.
func TestAC3POSIXSealStillWorksNextToAndThroughRealDirectoriesAndLinks(t *testing.T) {
	f := newForeign113(t, "foreign")
	root := filepath.Join(winsec.SealableTempDirForTest124(t), "data")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatal(err)
	}
	// A link that is NOT on the ancestor chain of anything sealed below.
	linkTo113(t, root, "not-an-ancestor", f.dir)

	deep := filepath.Join(root, "a", "b", "c")
	if err := winsec.PrivateDirAll(deep, 0o700); err != nil {
		t.Fatalf("AC#3 RED: PrivateDirAll refused its own tree: %v", err)
	}
	for _, p := range []string{filepath.Join(root, "a"), filepath.Join(root, "a", "b"), deep} {
		info, err := os.Stat(p)
		if err != nil {
			t.Fatalf("AC#3: %s was not created: %v", p, err)
		}
		if got := info.Mode().Perm(); got != 0o700 {
			t.Errorf("AC#3 RED: %s landed %o, not 0700", p, got)
		}
	}
	file := filepath.Join(deep, "artifact.txt")
	if err := winsec.PrivateFile(file, []byte("tool output"), 0o644); err != nil {
		t.Fatalf("AC#3 RED: PrivateFile refused a plain path inside the named tree: %v", err)
	}
	if got, err := os.Stat(file); err != nil || got.Mode().Perm() != 0o600 {
		t.Errorf("AC#3: %s is %v (%v), want -rw-------", file, got.Mode(), err)
	}
	if err := winsec.SealDir(deep); err != nil {
		t.Errorf("AC#3 RED: SealDir refused its own directory: %v", err)
	}
	// A symlink sitting below a file that is being sealed is not an ancestor of
	// it, so sealing the tree's own file still works.
	if err := os.Symlink(filepath.Join(root, "elsewhere"), filepath.Join(deep, "dangling")); err != nil {
		t.Skipf("no symlink here: %v", err)
	}
	if err := winsec.SealDir(deep); err != nil {
		t.Errorf("AC#3 RED: SealDir refused a directory that merely has a link among its children: %v", err)
	}
	f.untouched(t, "AC#3 the foreign tree pointed at by a sibling link")
}

// TestAC3POSIXSealDoesNotFoldABackslashIntoASeparator is AC#3's other reverse
// half, and the pair to the fail-open pin above: there is a symlink at root/a and
// a real directory named `a\b` beside it. The input names the real directory, so
// the seal has to happen - an implementation that treats the backslash as a
// separator sees root/a as an ancestor and refuses its own tree.
func TestAC3POSIXSealDoesNotFoldABackslashIntoASeparator(t *testing.T) {
	f := newForeign113(t, "foreign")
	root := filepath.Join(winsec.SealableTempDirForTest124(t), "root")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatal(err)
	}
	linkTo113(t, root, "a", f.dir)
	realDir := filepath.Join(root, `a\b`)
	if err := os.Mkdir(realDir, 0o700); err != nil {
		t.Fatal(err)
	}
	mine := filepath.Join(realDir, "blob.bin")
	if err := os.WriteFile(mine, []byte("this tree's own data"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(mine, 0o666); err != nil {
		t.Fatal(err)
	}

	if err := winsec.SealFile(mine); err != nil {
		t.Errorf("AC#3 RED: POSIX folded the backslash in %q into a separator and refused the tree's own file: %v", mine, err)
	} else if info, err := os.Stat(mine); err != nil || info.Mode().Perm() != 0o600 {
		t.Errorf("AC#3: %s is %v (%v), want -rw-------", mine, info.Mode(), err)
	}
	f.untouched(t, "AC#3 the symlink target that shares a prefix once a backslash is folded")
}
