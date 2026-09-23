package tools

import (
	"os"
	"path/filepath"
	"testing"
)

// Ticket 107. Two portable guards over the same seam, both with NO skip path:
// they must run on Windows and on POSIX, because the defect this ticket fixes
// was invisible to every local (Windows) gate.

// TestTicket107AllowlistJudgmentTwoShapes feeds ONE input in two shapes to the
// judgment function (PathCanonicalizer.InAllowlist):
//
//	shape R (resolved real-tree form): the [fs] allowed_dirs entry is spelled as
//	           the absolute path C26 itself produces for the tree on disk.
//	shape E (expanded form):           the same entry spelled with the %VAR%
//	           construct a config uses, so it travels C26's expansion step.
//
// The target string handed to InAllowlist is byte-identical in both arms. If
// both judge true the comparison fold is innocent and whatever misjudges does
// so upstream of it; if only one arm judges wrong, that arm names the guilty
// leg. The cross feed pins this down further: it hands the shape-E canonical to
// the canonicalizer whose root is already the resolved tree, so the comparison
// leg is exercised on the exact string the failing platform rejected.
func TestTicket107AllowlistJudgmentTwoShapes(t *testing.T) {
	existing := sealableTempDir124(t)
	sub := filepath.Join(existing, "proj")
	if err := os.MkdirAll(sub, 0o700); err != nil {
		t.Fatalf("mkdir root: %v", err)
	}
	target := filepath.Join(sub, "a.txt")
	if err := os.WriteFile(target, []byte("wisp"), 0o600); err != nil {
		t.Fatalf("write target: %v", err)
	}

	const name = "WISP107_ROOT"
	t.Setenv(name, existing)
	spelled := "%" + name + "%" + pathSep + "proj"

	cr := NewPathCanonicalizer([]string{sub}, nil)
	tr, err := cr.Canonicalize(target)
	if err != nil {
		t.Fatalf("shape R Canonicalize: %v", err)
	}
	gotR := cr.InAllowlist(tr)

	ce := NewPathCanonicalizer([]string{spelled}, nil)
	te, err := ce.Canonicalize(target)
	if err != nil {
		t.Fatalf("shape E Canonicalize: %v", err)
	}
	gotE := ce.InAllowlist(te)
	gotCross := cr.InAllowlist(te)

	t.Logf("shape R root=%q target=%q roots=%v InAllowlist=%v", sub, tr, cr.Roots(), gotR)
	t.Logf("shape E root=%q target=%q roots=%v rewritten=%v unusable=%v InAllowlist=%v",
		spelled, te, ce.Roots(), ce.RewrittenRoots(), ce.UnusableRoots(), gotE)
	t.Logf("cross feed (root from shape R, canonical from shape E) InAllowlist=%v", gotCross)

	if !gotR {
		t.Errorf("shape R (root already resolved) InAllowlist(%q) = false, want true", tr)
	}
	if !gotE {
		t.Errorf("shape E (root through C26 expansion) InAllowlist(%q) = false, want true: "+
			"the expanded root is a real tree, so dropping it turns ticket 102 into option (A); roots=%v unusable=%v",
			te, ce.Roots(), ce.UnusableRoots())
	}
	if !gotCross {
		t.Errorf("comparison leg rejects a canonical the resolved root owns: %q", te)
	}
}

// TestTicket107AllowlistBoundaryIsComponentWise guards the OTHER direction: the
// release side must recognize only the one resolved form, so a sibling tree
// whose name merely starts with an allowed root's characters (or the root's own
// parent) stays outside the allowlist. Weakening the fold to "any prefix" makes
// this red.
func TestTicket107AllowlistBoundaryIsComponentWise(t *testing.T) {
	existing := sealableTempDir124(t)
	root := filepath.Join(existing, "proj")
	sibling := root + "-evil"
	for _, d := range []string{root, sibling} {
		if err := os.MkdirAll(d, 0o700); err != nil {
			t.Fatalf("mkdir %q: %v", d, err)
		}
	}
	pc := NewPathCanonicalizer([]string{root}, nil)
	if len(pc.Roots()) != 1 {
		t.Fatalf("Roots() = %v, want the one real root kept", pc.Roots())
	}

	inside, err := pc.Canonicalize(filepath.Join(root, "a.txt"))
	if err != nil {
		t.Fatalf("Canonicalize(inside): %v", err)
	}
	if !pc.InAllowlist(inside) {
		t.Errorf("InAllowlist(%q) = false, want true: a file under the allowed root must pass", inside)
	}

	outside, err := pc.Canonicalize(filepath.Join(sibling, "b.txt"))
	if err != nil {
		t.Fatalf("Canonicalize(sibling): %v", err)
	}
	if pc.InAllowlist(outside) {
		t.Errorf("InAllowlist(%q) = true: root %q does not own a sibling that only shares its prefix",
			outside, root)
	}

	parent, err := pc.Canonicalize(existing)
	if err != nil {
		t.Fatalf("Canonicalize(parent): %v", err)
	}
	if pc.InAllowlist(parent) {
		t.Errorf("InAllowlist(%q) = true: a root authorizes its subtree, not its parent", parent)
	}
}
