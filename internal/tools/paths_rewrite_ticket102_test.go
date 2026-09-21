package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Ticket 102 (fix B), authorization leg: an [fs] allowed_dirs entry whose
// spelling C26's expansion step substituted is (a) recorded and operator-visible
// through RewrittenRoots, and (b) dropped - authorizing nothing - when the tree
// it moved onto cannot be confirmed on disk. Expansion itself stays in the
// pipeline, because a root spelled ~/... or %VAR%\... is contract behaviour
// (SPEC-06 §4 step 1), not a defect.

func TestPathCanonicalizerAccountsForRewrittenRoots(t *testing.T) {
	existing := t.TempDir()
	sub := filepath.Join(existing, "proj")
	if err := os.MkdirAll(sub, 0o700); err != nil {
		t.Fatalf("mkdir root: %v", err)
	}

	const name = "WISP102TOOLS_ROOT"
	t.Setenv(name, existing)
	spelled := "%" + name + "%" + pathSep + "proj"

	pc := NewPathCanonicalizer([]string{spelled}, nil)
	rewrote := pc.RewrittenRoots()
	if len(rewrote) != 1 {
		t.Fatalf("RewrittenRoots() = %v, want the one %q root recorded", rewrote, spelled)
	}
	if !strings.Contains(rewrote[0], "%"+name+"%") || !strings.Contains(rewrote[0], sub) || !strings.Contains(rewrote[0], "env") {
		t.Errorf("RewrittenRoots()[0] = %q, want spelling + expanded tree + the construct that moved it", rewrote[0])
	}

	// Half the point of (B): the mark is only true when something expanded.
	plain := NewPathCanonicalizer([]string{sub}, nil)
	if got := plain.RewrittenRoots(); len(got) != 0 {
		t.Errorf("RewrittenRoots() = %v for a root spelled as an absolute path, want empty", got)
	}

	// A root C26 moved onto a tree that is not there authorizes nothing.
	const ghost = "WISP102TOOLS_GHOST"
	t.Setenv(ghost, filepath.Join(existing, "no-such-tree"))
	gp := NewPathCanonicalizer([]string{"%" + ghost + "%" + pathSep + "proj"}, nil)
	if len(gp.Roots()) != 0 {
		t.Errorf("Roots() = %v, want empty: an expanded root that cannot be confirmed must not authorize", gp.Roots())
	}
	if len(gp.UnusableRoots()) != 1 || !strings.Contains(gp.UnusableRoots()[0], "not confirmed on disk") {
		t.Errorf("UnusableRoots() = %v, want the reason recorded for the operator", gp.UnusableRoots())
	}

	// Platform capability probe: handle-based resolution only exists on Windows
	// (see internal/risk/pathresolver_other.go), so the "a confirmed rewritten
	// root still authorizes its tree" half is only asserted where it can be
	// falsified.
	if len(NewPathCanonicalizer([]string{sub}, nil).Roots()) == 1 {
		c, err := pc.Canonicalize(filepath.Join(sub, "a.txt"))
		if err != nil {
			t.Fatalf("Canonicalize(target): %v", err)
		}
		if !pc.InAllowlist(c) {
			t.Errorf("InAllowlist(%q) = false although the expanded root is a real tree; the fix must not turn into option (A)", c)
		}
	} else {
		t.Logf("platform offers no handle-resolved form for an existing directory; authorization half not asserted here")
	}
}
