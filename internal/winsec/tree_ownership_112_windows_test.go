//go:build windows

package winsec_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/winsec"
	"golang.org/x/sys/windows"
)

// This file is ticket 112's instrument for the first two of the three reds that
// internal/winsec's new CI step caught on its first real run (run 35595651898,
// job 106319703680):
//
//	resolve_windows_test.go:30: no C26 pipeline installed in winsec: internal/risk/winsec_c26.go's init did not run
//	resolve_windows_test.go:90: C26 is not installed, so this leg measured the fallback instead
//	resolve_windows_test.go:97: C26 is not installed, so this leg measured the fallback instead
//
// Both AC3 legs and the wiring test share one cause, and the CI log names it in
// the ERROR record the guard itself wrote:
//
//	winsec: refusing to install a path resolver into the sealing seam resolver=risk.c26Pipeline
//	reason="it answered "C:\\Users\\RUNNER~1\\AppData\\Local\\Temp\\wisp-108-tree-ownership-probe"
//	with "C:\\Users\\RUNNER~1\\...", which is not inside the tree
//	"C:\\Users\\runneradmin\\AppData\\Local\\Temp" it answered for that path's own parent"
//
// The seam guard asked about an existing directory (which the pipeline answers
// with the filesystem's real path, 8.3 expanded) and a missing child (which it
// answers lexically, still RUNNER~1), compared the two spellings, and read an
// honest resolver as a tree mover. That only happens on a machine whose profile
// directory is longer than eight characters - GitHub's runner logs in as
// "runneradmin", this box does not - so the reproduction has to be planted
// rather than waited for (AC#3, ticket 106's shape): make a long directory name
// under a temporary root, ask the OS for its short form, and drive the guard's
// own probe at that pair.
//
// The //go:build above is not a dodge: an 8.3 short name is a Windows filesystem
// property and there is no POSIX counterpart to plant. The POSIX leg of the same
// seam guard lives in ancestor_separator_108_other_test.go (ticket 113's file).

// c26Pipeline112 adapts the real C26 pipeline to the seam shape. It is the same
// three legs as the production wiring in internal/risk/winsec_c26.go, which
// TestC26PipelineIsWiredIntoWinsec pins as installed; here the point is to put
// that pipeline in front of the guard's probe with a pair the production init
// cannot choose for itself.
type c26Pipeline112 struct{}

func (c26Pipeline112) Resolve(input string) (string, error) {
	res, err := risk.Resolve(input, nil)
	if err != nil {
		return "", err
	}
	return res.Actable()
}

func (c26Pipeline112) ResolveAccounted(input string) (string, bool, error) {
	res, err := risk.Resolve(input, nil)
	if err != nil {
		return "", false, err
	}
	path, err := res.Actable()
	if err != nil {
		return "", res.Rewritten, err
	}
	return path, res.Rewritten, nil
}

// rerootChildOnly112 names a foreign tree for the child while answering honestly
// for the child's own parent directory: exactly the shape the second witness is
// there to keep refusing after ticket 112 stopped treating a spelling difference
// as a tree move.
type rerootChildOnly112 struct{ foreign string }

func (r rerootChildOnly112) Resolve(input string) (string, error) {
	p, _, err := r.ResolveAccounted(input)
	return p, err
}

func (r rerootChildOnly112) ResolveAccounted(input string) (string, bool, error) {
	if filepath.Dir(input) == r.foreign || input == r.foreign {
		return input, false, nil
	}
	return filepath.Join(r.foreign, filepath.Base(input)), false, nil
}

// constant112 answers every input with one spelling, which is ticket 103's probe
// P1b shape and is caught before containment by either witness.
type constant112 struct{ fixed string }

func (c constant112) Resolve(string) (string, error) { return c.fixed, nil }

func (c constant112) ResolveAccounted(string) (string, bool, error) {
	return c.fixed, false, nil
}

// TestTreeOwnershipProbeAcceptsAn83ShortSpellingOfItsOwnParent is AC#1's third
// reading and AC#3's instrument: the honest pipeline must survive the guard on a
// pair spelled the way the CI runner spells it, and the guard must still bite on
// the fakes it bit before.
func TestTreeOwnershipProbeAcceptsAn83ShortSpellingOfItsOwnParent(t *testing.T) {
	t.Run("on a volume that generates 8.3 names", func(t *testing.T) {
		base := t.TempDir()
		long := filepath.Join(base, "wisp11283parentdir")
		if err := os.Mkdir(long, 0o700); err != nil {
			t.Fatal(err)
		}
		short := shortPathOf(t, long)
		if strings.EqualFold(short, long) {
			// Capability gate, and the only honest one available: with 8.3 name
			// generation disabled (the NtfsDisable8dot3NameCreation default on
			// some volumes) this machine cannot spell a directory two ways, so the
			// runner's asymmetry cannot be planted here. Kept in a subtest because
			// a top-level SKIP is fatal to the CI gate (ticket 71 AC#3), which is
			// the right rule: the leg below is what the gate is for, and it is
			// reported by name as skipped rather than quietly absent.
			t.Skipf("this volume gives %q no short name (GetShortPathName returned %q), so the runner's RUNNER~1 shape cannot be planted", long, short)
		}
		// The child must not exist: that is half of the asymmetry, because the
		// pipeline answers a missing leaf lexically and an existing directory with
		// the filesystem's own real path.
		child := filepath.Join(short, "wisp-112-missing-child")

		pipeline := c26Pipeline112{}
		parentAns, _, err := pipeline.ResolveAccounted(short)
		if err != nil {
			t.Fatalf("the pipeline refused the planted short spelling %q: %v", short, err)
		}
		childAns, _, err := pipeline.ResolveAccounted(child)
		if err != nil {
			t.Fatalf("the pipeline refused the missing child %q: %v", child, err)
		}
		// Read the instrument, not the verdict: if the two answers are already in
		// the same spelling family then nothing was planted and this test would
		// pass for the reason the CI run did not fail for.
		if strings.EqualFold(parentAns, short) {
			t.Fatalf("the instrument measured nothing: the pipeline answered the existing %q with itself, so no 8.3 expansion happened on this box", short)
		}
		if !strings.EqualFold(childAns, child) {
			t.Fatalf("the instrument measured nothing: the pipeline expanded the missing child to %q, so this box's pipeline does not have the asymmetry that refused the install on the runner", childAns)
		}

		if reason := winsec.TreeOwnershipProbeForTest(pipeline, short, child); reason != "" {
			t.Fatalf("AC#1 RED: the guard still refuses the honest pipeline on the runner's own shape: %s\nanswers: parent %q -> %q, child %q -> %q",
				reason, short, parentAns, child, childAns)
		}
		// And the same probe, with the short spelling replaced by the long one, has
		// to keep the verdict it always had: the fix makes the guard spelling-blind,
		// it does not retire the leg.
		if reason := winsec.TreeOwnershipProbeForTest(pipeline, parentAns, filepath.Join(parentAns, "wisp-112-missing-child")); reason != "" {
			t.Fatalf("the guard refuses the honest pipeline on the long spelling too: %s", reason)
		}

		foreign := filepath.Join(t.TempDir(), "someone-elses-tree")
		if err := os.Mkdir(foreign, 0o700); err != nil {
			t.Fatal(err)
		}
		mover := rerootChildOnly112{foreign: foreign}
		if reason := winsec.TreeOwnershipProbeForTest(mover, short, child); reason == "" {
			t.Errorf("AC#2 REGRESSION: the resolver that answers %q with something under %q was accepted once the guard stopped comparing spellings; the refusal it earned: %q", child, foreign, reason)
		} else {
			t.Logf("the tree-moving candidate is still refused: %s", reason)
		}
		fixed := filepath.Join(foreign, "fixed")
		if reason := winsec.TreeOwnershipProbeForTest(constant112{fixed: fixed}, short, child); reason == "" {
			t.Errorf("AC#2 REGRESSION: the constant fake %q was accepted by the tree-ownership leg", fixed)
		}
		// A refusal is still a narrow answer, on every machine.
		if reason := winsec.TreeOwnershipProbeForTest(unresolving112{}, short, child); reason != "" {
			t.Errorf("a candidate that refuses the probe was refused by the guard: %q", reason)
		}
	})
}

type unresolving112 struct{}

func (unresolving112) Resolve(string) (string, error) {
	return "", fmt.Errorf("%w: instrument refuses everything", winsec.ErrUnresolvedPath)
}

func (unresolving112) ResolveAccounted(in string) (string, bool, error) {
	return "", false, fmt.Errorf("instrument refuses %s", in)
}

// shortPathOf asks the OS for the 8.3 rendering of an existing object. No
// hand-built short name: the whole defect is about a spelling this package does
// not control, so the fixture takes the one the filesystem chose.
func shortPathOf(t *testing.T, path string) string {
	t.Helper()
	p, err := windows.UTF16PtrFromString(path)
	if err != nil {
		t.Fatal(err)
	}
	n, err := windows.GetShortPathName(p, nil, 0)
	if err != nil && n == 0 {
		t.Fatalf("GetShortPathName(%s) size probe: %v", path, err)
	}
	buf := make([]uint16, n)
	if _, err := windows.GetShortPathName(p, &buf[0], n); err != nil {
		t.Fatalf("GetShortPathName(%s): %v", path, err)
	}
	return windows.UTF16ToString(buf)
}
