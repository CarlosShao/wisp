//go:build windows

package risk

import (
	"os"
	"path/filepath"
	"testing"
)

// The Windows half of ticket 336 AC#2's world: the roster and every assertion
// live in anchor_fallback_teeth_336_roster_test.go (ticket 347), and what is
// left here is only what this platform answers differently - how a spelling is
// made present-but-unprovable.
//
// On Windows that is a name the OS refuses to PARSE: hostile334 carries U+001B
// and U+202E, Lstat fails with "syntax is incorrect" rather than "no such file",
// pathExists reads that as present, and pathForms.partial becomes true. The
// POSIX twin cannot use it (both characters are legal there), so it builds a
// dangling symlink instead; see anchor_world_334_test.go for the contrast.

// teeth336Candidate places the hostile candidate in the roster row's directory,
// after checking the two setup facts the whole shape rests on.
func teeth336Candidate(t *testing.T, w world334, tc teeth336Case) string {
	t.Helper()
	dir := w.work
	if tc.relCand != "work" {
		dir = filepath.Join(w.profile, filepath.FromSlash(tc.relCand))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("setup mkdir %s: %v", dir, err)
		}
	}
	if _, err := os.Lstat(filepath.Join(w.profile, filepath.FromSlash(tc.omit))); err == nil {
		t.Fatalf("setup: the leaf under test must be absent: %s", tc.omit)
	}
	cand := filepath.Join(dir, hostile334)
	if _, err := os.Lstat(cand); err == nil {
		t.Fatalf("setup: the hostile candidate must not exist at %q", cand)
	}
	return cand
}

// teeth336OffChainCandidate is the off-chain arm's candidate: the same hostile
// name beside the profile instead of inside it.
func teeth336OffChainCandidate(t *testing.T, installed []string) string {
	t.Helper()
	return filepath.Join(newWorld334(t, false, installed...).work, hostile334)
}
