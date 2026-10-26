//go:build !windows

package risk

import (
	"os"
	"path/filepath"
	"testing"
)

// The POSIX half of ticket 336 AC#2's world, which ticket 347 is what finally
// made buildable. The roster and every assertion live in
// anchor_fallback_teeth_336_roster_test.go, unchanged from the Windows file
// ticket 336 wrote; only the way a spelling is made present-but-unprovable
// differs here, and it is MEASURED rather than assumed:
//
//   - before ticket 347 this file could hold no tooth at all. resolveHandle was
//     the DEFERRED(macOS/Linux) stub, pathForms.root stayed "" on both sides,
//     uncertainAnchorMiss needs a non-empty root on both, so the branch was
//     structurally unreachable and a tooth would have been permanently red.
//     What stood here was a tripwire that said exactly that;
//   - when the resolver landed, the tripwire went red and the ticket-334 shape
//     still read class B with clause=false, because U+001B and U+202E are LEGAL
//     filename bytes on POSIX: the name is merely missing, Lstat says "no such
//     file", pathExists proves the miss and partial stays false. Copying the
//     Windows maker here would have produced a tooth that passes for the wrong
//     reason (or an off-chain arm that is accidentally right);
//   - so the maker is a DANGLING SYMLINK. The entry exists (Lstat succeeds, which
//     is what partial reports) and cannot be opened (Stat/EvalSymlinks follow it
//     onto a target that is not there), which is the same pair the Windows name
//     produces - exactly the "dangling reparse point" case pathForms.partial
//     names in its own comment.

// teeth336Candidate places the hostile candidate in the roster row's directory,
// after checking the setup facts the whole shape rests on. The first two are
// ticket 336's verbatim (the leaf under test is absent, the candidate was not
// there before); the third is this platform's partial maker, asserted rather
// than assumed.
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
	cand := filepath.Join(dir, hostilePOSIX347)
	if _, err := os.Lstat(cand); err == nil {
		t.Fatalf("setup: the hostile candidate must not exist at %q", cand)
	}
	if err := os.Symlink(danglingTarget347, cand); err != nil {
		t.Fatalf("setup symlink %s: %v", cand, err)
	}
	if _, err := os.Stat(cand); err == nil {
		t.Fatalf("setup: the partial maker must not be openable, got a following stat for %q", cand)
	}
	if _, err := os.Lstat(cand); err != nil {
		t.Fatalf("setup: the partial maker must still be a directory entry (Lstat is what makes the spelling unproven), got %v", err)
	}
	return cand
}

// teeth336OffChainCandidate is the off-chain arm's candidate: the same dangling
// link beside the profile instead of inside it, so no A anchor is
// ancestry-related to it and only a real name match can reach the B tier.
func teeth336OffChainCandidate(t *testing.T, installed []string) string {
	t.Helper()
	cand := filepath.Join(newWorld334(t, false, installed...).work, hostilePOSIX347)
	if err := os.Symlink(danglingTarget347, cand); err != nil {
		t.Fatalf("setup symlink %s: %v", cand, err)
	}
	return cand
}
