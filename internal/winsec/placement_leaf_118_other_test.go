//go:build !windows

// Ticket 118 AC#6, from acceptor-ticket113's R-113-C: the POSIX link leg had no
// case for the leaf direction.
//
// What shipped for ticket 113 AC#1 plants every link in an *ancestor* position
// (`<root>/link/keep-me.txt`), and the one case that puts a link where the caller
// named the object is TestAC1POSIXSealDirThroughASymlinkRefuses. That left the
// acceptance side's MUT-B - `pieces = pieces[:len(pieces)-1]`, an implementation
// that walks the ancestors and stops before the leaf - reddening exactly one case.
// It then wrote two cases of its own, SealFile(link) and PrivateFile(link), that
// are green on the shipped code and red on MUT-B: the implementation is right and
// the coverage was missing two. Those two are here.
//
// They are the direction that matters most on this platform, not a bonus: the
// comment over platformVerifyPlacement says the leg "checks the leaf as well as
// the ancestors", because os.Chmod follows the leaf - so a symlink standing where
// a seal was asked for is the one position where the wrong object is guaranteed to
// be the object the mode change lands on.
//
// Run for real, not compile-only, in a Linux container: see the ticket's AC#6 row
// for the mount, the rc and the MUT-B reading.
package winsec_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/CarlosShao/wisp/internal/winsec"
)

// leafLinkTo118 puts a symlink AT the position a caller would name an object,
// pointing straight at the foreign victim, and returns the spelling. This is the
// shape the 113 file never plants: the link is not an ancestor of the argument, it
// IS the argument.
func leafLinkTo118(t *testing.T, root, linkName, target string) string {
	t.Helper()
	link := filepath.Join(root, linkName)
	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("cannot build the leaf symlink this case measures: %v", err)
	}
	info, err := os.Lstat(link)
	if err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("Lstat(%s) after os.Symlink: info=%v err=%v", link, info, err)
	}
	// The spelling has to reach the victim through the link and nowhere else: if
	// the target were also reachable as a plain path inside root, a refusal could
	// be credited to the wrong leg.
	if same, err := filepath.EvalSymlinks(link); err != nil || same != target {
		t.Fatalf("the leaf link does not answer with its own target: %q vs %q (err=%v)", same, target, err)
	}
	return link
}

// ownTree118 is the tree the caller legitimately holds: one private directory with
// nothing foreign in it, used as the anchor the leaf link hangs off.
func ownTree118(t *testing.T) string {
	t.Helper()
	root := filepath.Join(winsec.SealableTempDirForTest124(t), "root")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatal(err)
	}
	return root
}

// TestAC118POSIXSealFileRefusesALinkStandingWhereTheFileWasNamed is the first of
// the two missing directions: SealFile asked to seal a link. The mode change this
// call performs follows the link, so the only correct answer is the refusal - and
// the victim's own mode is the read that says which object got the chmod.
//
// MUT-B (acceptance's half-fix, the walk stopping one component short of the leaf)
// reddens this case twice over: assertRefused113 sees a nil, and untouched sees the
// victim drop from 0666 to 0600.
func TestAC118POSIXSealFileRefusesALinkStandingWhereTheFileWasNamed(t *testing.T) {
	f := newForeign113(t, "foreign")
	root := ownTree118(t)
	link := leafLinkTo118(t, root, "artifact.txt", f.victim)

	victimBefore := statFact113(f.victim)
	assertRefused113(t, "SealFile", link, winsec.SealFile(link))

	if got := statFact113(f.victim); got != victimBefore {
		t.Errorf("AC#6 RED: SealFile(%q) left the victim at %s, it was %s - the mode change followed the leaf", link, got, victimBefore)
	}
	f.untouched(t, "AC#6 SealFile on a leaf link")
	if info, err := os.Lstat(link); err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Errorf("AC#6: a refusal must not touch anything, but the link now reads %v (%v)", info, err)
	}
}

// TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasNamed is the second
// missing direction, on the entry point that writes bytes: PrivateFile over a leaf
// link opens with O_WRONLY|O_CREATE|O_TRUNC, which follows the link, so the failure
// this case measures is not a wide file but somebody else's content replaced.
func TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasNamed(t *testing.T) {
	f := newForeign113(t, "foreign")
	root := ownTree118(t)
	link := leafLinkTo118(t, root, "artifact.txt", f.victim)

	contentBefore, err := os.ReadFile(f.victim)
	if err != nil {
		t.Fatal(err)
	}
	factBefore := statFact113(f.victim)
	err = winsec.PrivateFile(link, []byte("top secret"), 0o600)
	assertRefused113(t, "PrivateFile", link, err)

	contentAfter, err := os.ReadFile(f.victim)
	if err != nil {
		t.Fatalf("AC#6: the victim is not even readable any more: %v", err)
	}
	if string(contentAfter) != string(contentBefore) {
		t.Errorf("AC#6 RED: PrivateFile(%q) replaced the victim's bytes (%q -> %q); the open followed the leaf",
			link, contentBefore, contentAfter)
	}
	if got := statFact113(f.victim); got != factBefore {
		t.Errorf("AC#6 RED: PrivateFile(%q) left the victim at %s, it was %s", link, got, factBefore)
	}
	f.untouched(t, "AC#6 PrivateFile on a leaf link")
	if _, err := os.Lstat(link); err != nil {
		t.Errorf("AC#6: the refusal removed the link it was handed: %v", err)
	}
}
