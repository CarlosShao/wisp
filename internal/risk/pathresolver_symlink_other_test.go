//go:build !windows

package risk

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Ticket 347 AC#1, the POSITIVE direction: the A tier's inputs on POSIX, now
// that C26 has the two primitives PLAN.md:2385 promised for this side.
//
// The negative/half-world direction is carried by the moved teeth
// (anchor_fallback_teeth_336_roster_test.go, anchor_miss_reason_339_other_test.go):
// they assert that a spelling whose own walk stopped short now fail-closes into
// the non-overridable A tier here. This file asserts the other half - that a
// REAL symlink is refused by default and, when explicitly exempted, is resolved
// onto its target (realpath) rather than passed through.

// mustMkdir347 is the local mkdir-or-fail helper.
func mustMkdir347(t *testing.T, dirs ...string) {
	t.Helper()
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatalf("setup mkdir %s: %v", d, err)
		}
	}
}

// TestResolveDeniesSymlinkTraversalByDefault347 is SPEC-06 §4's default-deny
// rule on this platform: a link whose target sits outside the tree the caller
// named must not be followed into it, and the refusal has to be C26's own error
// (not a bare stat failure) so the caller can report the reason.
func TestResolveDeniesSymlinkTraversalByDefault347(t *testing.T) {
	base := t.TempDir()
	inside := filepath.Join(base, "profile")
	outside := filepath.Join(base, "outside")
	mustMkdir347(t, inside, outside)
	target := filepath.Join(outside, "secret.txt")
	if err := os.WriteFile(target, []byte("outside the named tree\n"), 0o600); err != nil {
		t.Fatalf("setup write %s: %v", target, err)
	}
	link := filepath.Join(inside, "shortcut")
	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("setup symlink %s: %v", link, err)
	}

	res, err := Resolve(link, nil)
	if !errors.Is(err, ErrReparseDenied) {
		t.Fatalf("Resolve(%q) = (%+v, %v), want ErrReparseDenied: an unexempted symlink traversal is SPEC-06 §4's default deny", link, res, err)
	}
	if res.Resolved {
		t.Errorf("a refused traversal must not report a resolved canonical: %+v", res)
	}
	// The walk reports the component AS SPELLED, which is what an exception
	// entry and the audit line have to carry.
	comps := reparseComponents(lexCanonical(link))
	if len(comps) != 1 || comps[0] != lexCanonical(link) {
		t.Errorf("reparseComponents(%q) = %v, want exactly [%q]", link, comps, lexCanonical(link))
	}
	// An exception must name the COMPONENT, not the destination: listing the
	// target does not authorize traversing the link (the table holds explicit
	// paths, never prefixes or "the same tree").
	if _, err := Resolve(link, []string{target}); !errors.Is(err, ErrReparseDenied) {
		t.Errorf("Resolve(%q, exception=%q) did not deny: the exception table names components, not destinations, got %v", link, target, err)
	}
	t.Logf("POSIX default-deny: %v", err)
}

// TestResolveFollowsExemptedSymlinkOntoRealPath347 is the positive control for
// the same primitive: with the exact component listed in
// [fs] reparse_point_exceptions the link is traversed, and what C26 hands back
// is realpath's answer - the TARGET - not the link's own spelling. That is the
// difference between "resolved" and "the same string back", and it is what makes
// the A tier reachable here by resolution instead of by spelling.
func TestResolveFollowsExemptedSymlinkOntoRealPath347(t *testing.T) {
	base := t.TempDir()
	home := filepath.Join(base, "profile")
	work := filepath.Join(home, "work")
	ssh := filepath.Join(home, ".ssh")
	mustMkdir347(t, home, work, ssh)
	// The A anchors are read from the environment; HOME is what os.UserHomeDir
	// reads here (USERPROFILE is set too so the world looks like the other files').
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	target := filepath.Join(ssh, "id_ed25519")
	if err := os.WriteFile(target, []byte("PRIVATE KEY\n"), 0o600); err != nil {
		t.Fatalf("setup write %s: %v", target, err)
	}
	link := filepath.Join(work, "innocent.txt")
	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("setup symlink %s: %v", link, err)
	}

	// Default deny applies even when the destination is inside the tree: the
	// contract refuses TRAVERSAL, it does not try to judge the destination.
	if _, err := Resolve(link, nil); !errors.Is(err, ErrReparseDenied) {
		t.Fatalf("an in-tree link must still be refused without an exception, got %v", err)
	}

	res, err := Resolve(link, []string{link})
	if err != nil {
		t.Fatalf("Resolve(%q, exception=%q): %v", link, link, err)
	}
	if !res.Resolved {
		t.Fatalf("with the component exempted the link must resolve, got %+v", res)
	}
	if !res.Reparse {
		t.Errorf("Resolved with an exempted traversal must report Reparse=true (the exception was consumed): %+v", res)
	}
	if res.Canonical != target {
		t.Errorf("canonical = %q, want the realpath answer %q (a link passed through unchanged is not resolution)", res.Canonical, target)
	}
	if got := Classify(res.Canonical); got != ClassA {
		t.Errorf("classify(%q) = %s, want A: the resolved form is ~/.ssh/** and must reach the A tier", res.Canonical, got)
	}
	t.Logf("POSIX resolved through an exempted symlink: %q -> %q class=%s", link, res.Canonical, Classify(res.Canonical))
}

// TestNormalizeLocalUNCLeavesPosixSpellingsAlone347 pins the compile-time early
// exit ticket 347 added to normalizeLocalUNC (pathresolver.go).
//
// The wake-up it closes is real but only reachable off the pipeline: a POSIX
// path may start with two slashes (an ordinary absolute path; POSIX leaves "//"
// implementation-defined and Go's Clean folds it to "/"), and before the early
// exit that spelling was rewritten wholesale into backslashes - a string no
// syscall can open and no comparison should see. Through Resolve it was dead
// only because lexCanonical (step 2) runs first, which is a property of the call
// ORDER, not of the function; the guard makes it a property of the platform.
func TestNormalizeLocalUNCLeavesPosixSpellingsAlone347(t *testing.T) {
	for _, in := range []string{"//localhost/c$/x", "//server/share/y", `\\server\share\z`} {
		if got := normalizeLocalUNC(in); got != in {
			t.Errorf("normalizeLocalUNC(%q) = %q: there is no UNC spelling to normalize on POSIX", in, got)
		}
	}
	// The same reading through the whole pipeline, kept because it is the one a
	// caller can actually reach: TestDoubleSlashSpellingIsNotUNC pins the
	// canonical of this input, this pins the primitive that would have broken it.
	res, err := Resolve("//localhost/c$/x", nil)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if strings.Contains(res.Canonical, `\`) {
		t.Errorf("canonical %q carries a backslash: a POSIX double-slash spelling was rewritten into UNC shape", res.Canonical)
	}
}
