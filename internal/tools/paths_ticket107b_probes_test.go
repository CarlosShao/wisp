package tools

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// Ticket 107, round 2 (agent-ticket107b). The three probes the adversarial
// acceptance ran inside a /tmp snapshot are re-runnable cases here, because a
// probe that cannot reach CI does not exist (acceptor's own next= line).
//
// What they pin down is the direction rule, not one platform's behaviour: the
// release side may authorize a path only when the tree it judges, the tree the
// tool opens and the tree the operator named are the SAME tree. Every assertion
// below expects "not authorized", is meant to hold on Windows and on POSIX, and
// has no skip path: a link that cannot be built, or whose presence cannot be
// observed, makes the probe fail loudly instead of passing vacuously.

// makeDirLink107b builds a directory link with a shape this platform can build
// without elevation: a symlink on POSIX, an NTFS junction on Windows (same
// command internal/winsec/seam_guard_windows_test.go uses). For the judgment
// under test the two are one threat object: a path component whose real tree
// lives somewhere else.
func makeDirLink107b(t *testing.T, link, target string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		out, err := exec.Command("cmd", "/c", "mklink", "/J", link, target).CombinedOutput()
		t.Logf("mklink /J %s -> %s: %s", link, target, strings.TrimSpace(string(out)))
		if err != nil {
			t.Fatalf("mklink /J failed: %v (%s): without the link this probe measures nothing", err, out)
		}
		return
	}
	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("os.Symlink(%s, %s): %v: without the link this probe measures nothing", link, target, err)
	}
}

// dirLinkEvidence107b proves the probe's premise on the platform it runs on and
// reports the tree the link routes to. It tries every signal this repository
// could reasonably consult and logs all of them, because the strength of the
// fix's own confirmation depends on which ones answer:
//
//	lstat mode bit, os.Readlink, filepath.EvalSymlinks on the link, and
//	filepath.EvalSymlinks on a file reached THROUGH the link.
//
// No signal seeing a switch is a Fatalf: the probe then says it measured
// nothing rather than going green for the wrong reason.
func dirLinkEvidence107b(t *testing.T, link, target string) string {
	t.Helper()
	marker := filepath.Join(target, "marker.txt")
	if err := os.WriteFile(marker, []byte("TOPSECRET"), 0o600); err != nil {
		t.Fatalf("write marker: %v", err)
	}
	fi, lerr := os.Lstat(link)
	modeLink := lerr == nil && fi.Mode()&os.ModeSymlink != 0
	rdst, rerr := os.Readlink(link)
	lres, lerr2 := filepath.EvalSymlinks(link)
	vres, verr := filepath.EvalSymlinks(filepath.Join(link, "marker.txt"))
	t.Logf("premise signals for %q -> %q: lstat(err=%v, modeLink=%v) readlink(%q, err=%v) eval(link)=%q,%v eval(through)=%q,%v",
		link, target, lerr, modeLink, rdst, rerr, lres, lerr2, vres, verr)
	seen := modeLink || rerr == nil || (lerr2 == nil && foldPath(lres) != foldPath(link)) ||
		(verr == nil && foldPath(vres) != foldPath(filepath.Join(link, "marker.txt")))
	if !seen {
		t.Fatalf("no signal reports %q as a link onto another tree, so this probe measures nothing on this platform", link)
	}
	if verr == nil {
		return vres
	}
	if lerr2 == nil {
		return filepath.Join(lres, "marker.txt")
	}
	return marker
}

// judge107b hands the canonicalizer what a tool would actually carry: its
// Canonicalize output. Where C26 refuses the path outright (Windows denies
// reparse traversal with risk.ErrReparseDenied) the refusal is itself the
// fail-closed outcome, but the judgment leg must still be exercised on the
// spelling a POSIX run feeds it, so the lexical join is used instead.
func judge107b(t *testing.T, pc *PathCanonicalizer, path string) string {
	t.Helper()
	c, err := pc.Canonicalize(path)
	if err != nil {
		t.Logf("Canonicalize(%q) refused: %v (judging the lexical spelling instead)", path, err)
		return path
	}
	return c
}

// TestTicket107bProbeASymlinkedRewrittenRootAuthorizesNothing is acceptance
// probe A: the ROOT itself is a link. allowed_dirs spells %VAR%<sep>proj, proj
// routes onto base/outside, and outside holds a file the operator never named.
// Letting that through is a cross-tree fail-open, and the three conditions the
// acceptance listed hold together: the release side says yes, the pre-fix code
// said no, the tree that opens is not the tree that was spelled.
func TestTicket107bProbeASymlinkedRewrittenRootAuthorizesNothing(t *testing.T) {
	base := t.TempDir()
	outside := filepath.Join(base, "outside")
	if err := os.MkdirAll(outside, 0o700); err != nil {
		t.Fatalf("mkdir outside: %v", err)
	}
	proj := filepath.Join(base, "proj")
	makeDirLink107b(t, proj, outside)
	real := dirLinkEvidence107b(t, proj, outside)

	via := filepath.Join(proj, "marker.txt")
	if body, err := os.ReadFile(via); err != nil || string(body) != "TOPSECRET" {
		t.Fatalf("reading %q gave %q, %v: the tool would not actually open the other tree", via, body, err)
	}

	const name = "WISP107B_A_ROOT"
	t.Setenv(name, base)
	spelled := "%" + name + "%" + pathSep + "proj"
	pc := NewPathCanonicalizer([]string{spelled}, nil)
	cand := judge107b(t, pc, via)
	t.Logf("roots=%v rewritten=%v unusable=%v judged=%q opened=%q",
		pc.Roots(), pc.RewrittenRoots(), pc.UnusableRoots(), cand, real)
	// The root leg has to have its own teeth: a name that is not the tree it
	// points at must never enter the authorization book, so that the audit line
	// and the deny leg's reason name a tree that exists rather than a link.
	if r := pc.Roots(); len(r) != 0 {
		t.Errorf("AC#3 RED: Roots() = %v: the named root %q is a link onto %q, so it authorizes nothing and must not be in the book",
			r, spelled, outside)
	}
	if u := pc.UnusableRoots(); len(u) == 0 {
		t.Errorf("UnusableRoots() is empty: a dropped root must leave the operator a reason, got roots=%v", pc.Roots())
	}
	if pc.InAllowlist(cand) {
		t.Errorf("AC#3 RED: InAllowlist(%q) = true although the named root %q expands onto %q, a link onto %q: the tree that opens (%q) is one the operator never named",
			cand, spelled, proj, outside, real)
	}
}

// TestTicket107bProbeCLinkInsideAllowedRootStaysOutside is acceptance probe C:
// the root is a real tree, but the target reaches another tree through a link
// INSIDE it. The first draft of the fix let this through too, which is why the
// resolved form has to be consulted on both legs - handling only the root half
// fixes one place out of two in the same invariant.
func TestTicket107bProbeCLinkInsideAllowedRootStaysOutside(t *testing.T) {
	base := t.TempDir()
	proj := filepath.Join(base, "proj")
	outside := filepath.Join(base, "outside")
	for _, d := range []string{proj, outside} {
		if err := os.MkdirAll(d, 0o700); err != nil {
			t.Fatalf("mkdir %q: %v", d, err)
		}
	}
	legit := filepath.Join(proj, "a.txt")
	if err := os.WriteFile(legit, []byte("wisp"), 0o600); err != nil {
		t.Fatalf("write legit: %v", err)
	}
	esc := filepath.Join(proj, "esc")
	makeDirLink107b(t, esc, outside)
	real := dirLinkEvidence107b(t, esc, outside)

	const name = "WISP107B_C_ROOT"
	t.Setenv(name, base)
	spelled := "%" + name + "%" + pathSep + "proj"
	pc := NewPathCanonicalizer([]string{spelled}, nil)
	t.Logf("roots=%v rewritten=%v unusable=%v", pc.Roots(), pc.RewrittenRoots(), pc.UnusableRoots())

	// The named tree keeps working: this is the leg ticket 102 went red on, so
	// tightening the link case may not pay for it by breaking the plain one.
	if c := judge107b(t, pc, legit); !pc.InAllowlist(c) {
		t.Errorf("InAllowlist(%q) = false, want true: %q is a real tree the operator named; roots=%v unusable=%v",
			c, proj, pc.Roots(), pc.UnusableRoots())
	}

	via := filepath.Join(esc, "marker.txt")
	if body, err := os.ReadFile(via); err != nil || string(body) != "TOPSECRET" {
		t.Fatalf("reading %q gave %q, %v", via, body, err)
	}
	cand := judge107b(t, pc, via)
	if pc.InAllowlist(cand) {
		t.Errorf("AC#3 RED: InAllowlist(%q) = true: it leaves the allowed root %q through the link %q and lands on %q",
			cand, proj, esc, real)
	}
}

// TestTicket107bProbeBJudgedTreeIsTheOpenedTree measures the shape the
// acceptance filed as the POSIX legacy account: the root is a link but spelled
// as a plain absolute path, so ticket 102's rewritten-root account never ran on
// it. The assertion is the one that is true on both platforms: if the release
// side says
// yes, then the tree that opens must be owned by some root in the authorization
// book. Windows satisfies it because the handle-based resolver puts the real
// tree in the book; POSIX at this commit does not, and that is the old account
// this ticket now also has to speak about honestly.
func TestTicket107bProbeBJudgedTreeIsTheOpenedTree(t *testing.T) {
	base := t.TempDir()
	outside := filepath.Join(base, "outside")
	if err := os.MkdirAll(outside, 0o700); err != nil {
		t.Fatalf("mkdir outside: %v", err)
	}
	proj := filepath.Join(base, "proj")
	makeDirLink107b(t, proj, outside)
	real := dirLinkEvidence107b(t, proj, outside)

	pc := NewPathCanonicalizer([]string{proj}, nil)
	via := filepath.Join(proj, "marker.txt")
	if body, err := os.ReadFile(via); err != nil || string(body) != "TOPSECRET" {
		t.Fatalf("reading %q gave %q, %v", via, body, err)
	}
	cand := judge107b(t, pc, via)
	got := pc.InAllowlist(cand)
	t.Logf("unexpanded link root: roots=%v unusable=%v judged=%q opened=%q InAllowlist=%v",
		pc.Roots(), pc.UnusableRoots(), cand, real, got)
	if !got {
		return // fail-closed is always an acceptable answer for the release side
	}
	owned := false
	for _, r := range pc.Roots() {
		if f := foldPath(real); f == r || strings.HasPrefix(f, r+pathSep) {
			owned = true
		}
	}
	if !owned {
		t.Errorf("AC#3 RED: InAllowlist(%q) = true but the tree that opens (%q) is named by no root in the book (roots=%v)",
			cand, real, pc.Roots())
	}
}
