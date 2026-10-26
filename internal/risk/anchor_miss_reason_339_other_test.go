//go:build !windows

package risk

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Ticket 339 AC#1 shape A, the POSIX half - and since ticket 347 a real tooth
// instead of a tripwire.
//
// What stood here before: the ticket 334 instrument is tagged //go:build windows
// because on this platform its world was un-buildable, and the condition was in
// the resolver, not in the environment. resolveHandle (pathresolver_other.go)
// was the DEFERRED(macOS/Linux) stub returning ("", false) for EVERY spelling,
// so formsOf never filled pathForms.root, uncertainAnchorMiss (blacklist.go)
// needs a non-empty root on both sides, and R17's fail-closed clause could not
// fire on any input. The file therefore asserted the MECHANISM and went red the
// moment a POSIX resolver landed ("port the teeth named above"), because a
// permanently red tooth is the shape ticket 336 AC#2 forbids and a green that
// asserts nothing is worse.
//
// Ticket 347 landed the resolver, that tripwire went red exactly as designed,
// and the two teeth it asked for are below. The worlds still come from the
// shared builder (anchor_world_334_test.go) so this half cannot drift from the
// Windows half; the one thing this platform has to supply itself is the
// present-but-unprovable spelling, and it is a dangling symlink rather than
// hostile334 - U+001B and U+202E are legal filename bytes here, so the Windows
// maker would leave partial=false and the clause correctly inert (measured: the
// ticket-334 shape reads class B / clause=false on POSIX even with the resolver
// in place, docs/evidence/s1/347-r1-impl.md).
//
// The fixture names ticket 216 reads (candidateOffAnchorChain334 and friends)
// now live in anchor_world_334_test.go; the name is kept because
// name_claim_two_arm_216_test.go:89 names it.

// posixPartialCandidate347 creates the present-but-unprovable candidate: a
// symlink whose target does not exist. Lstat succeeds (the entry is there, which
// is what pathForms.partial reports), Stat and EvalSymlinks cannot open it (the
// target is not), and that pair is exactly what pathExists/tailExistsBelow read
// as "this component exists but its true spelling is unknown".
func posixPartialCandidate347(t *testing.T, dir, name string) string {
	t.Helper()
	cand := filepath.Join(dir, name)
	if _, err := os.Lstat(cand); err == nil {
		t.Fatalf("setup: %q must not exist before the maker runs", cand)
	}
	if err := os.Symlink(danglingTarget347, cand); err != nil {
		t.Fatalf("setup symlink %s: %v", cand, err)
	}
	if _, err := os.Stat(cand); err == nil {
		t.Fatalf("setup: %q follows to something; it must dangle", cand)
	}
	return cand
}

// TestAnchorMissReasonNamesCandidateAndMissingAnchorPosix339 is ticket 334 AC#1
// on the platform ticket 347 opened: the fail-closed rejection has to say which
// candidate was rejected and which anchor was not installed, instead of
// accusing the candidate of being the credential file that rule guards.
func TestAnchorMissReasonNamesCandidateAndMissingAnchorPosix339(t *testing.T) {
	// .git-credentials and .ssh ARE installed, so those two rules miss honestly
	// and the first fail-closed hit is the ~/.aws one - the exact sentence
	// ticket 334 was written about, now reachable here.
	w := newWorld334(t, true, ".git-credentials", ".ssh/id_installed", wispCfgLeaf334)
	cand := posixPartialCandidate347(t, w.work, hostilePOSIX347)

	// The reading below is only about the resolver if the environment still does
	// its part: the fake home must be the home the package actually sees.
	homeS := normDir(userHomeDir())
	if homeS != normPath(w.profile) {
		t.Errorf("fixture: the fake home stopped taking effect on this platform (home %q, profile %q); the readings would then measure the environment, not the resolver", homeS, w.profile)
	}

	// The mechanism, asserted rather than assumed (this is what the tripwire used
	// to say was impossible here). The anchor spelling is the RAW home, the way
	// classifyForms builds it since ticket 347.
	f := formsOf(cand)
	a := anchorForms(anchorPath(userHomeDir(), ".aws", "credentials"))
	if f.root == "" || a.root == "" {
		t.Fatalf("premise: both sides must have a resolved root on this platform since ticket 347, got cand=%+v anchor=%+v", f, a)
	}
	if anchorMatches(f, a, eqPath) {
		t.Fatalf("fixture: the candidate must NOT be a proven match for the aws anchor; a tooth may not be carried by a real name hit")
	}
	if !uncertainAnchorMiss(f, a) {
		t.Fatalf("fixture: the fail-closed clause must be what carries this verdict; uncertainAnchorMiss=false for cand=%+v anchor=%+v", f, a)
	}

	d := Gate(cand, nil)
	if d.Class != ClassA {
		t.Fatalf("fixture guard: the fail-closed shape must be A here, got %s", gateShape216(d))
	}
	if want := "A-list sensitive path (non-overridable): ~/.aws/credentials"; d.Reason == want {
		t.Errorf("the reason is back to naming only the anchor, which is the false accusation ticket 334 is about: %q", d.Reason)
	}
	for _, want := range []struct{ what, needle string }{
		{"the REAL candidate being rejected", normPath(cand)},
		{"the REAL absent anchor", filepath.Join(w.profile, ".aws", "credentials")},
		{"where the candidate's own walk stopped", w.work},
		{"the home the un-installed anchor collapsed to", w.profile},
	} {
		// Case-insensitive on purpose, and for the same reason the Windows twin
		// folds: the reason prints the spellings it compared, and the anchor side
		// is now the environment's own spelling while the candidate side is
		// normPath's fold (t.TempDir() carries the test name, i.e. uppercase).
		if !strings.Contains(strings.ToLower(d.Reason), strings.ToLower(want.needle)) {
			t.Errorf("reason must name %s (%q), got:\n\t%s", want.what, want.needle, d.Reason)
		}
	}
	for _, want := range []string{"fail-closed", "NOT a match"} {
		if !strings.Contains(d.Reason, want) {
			t.Errorf("reason must carry the word %q so a reader cannot take it for a match, got:\n\t%s", want, d.Reason)
		}
	}
	t.Logf("honest reject: %s", d.Reason)
}

// TestAnchorMissKeepsNonOverridableClassPosix339 is ticket 334 AC#2's level on
// this platform: the tier is unchanged by ticket 347, and the off-chain control
// keeps it from being satisfied by "this machine reads everything as A".
func TestAnchorMissKeepsNonOverridableClassPosix339(t *testing.T) {
	w := newWorld334(t, true, ".git-credentials", ".ssh/id_installed", wispCfgLeaf334)
	cand := posixPartialCandidate347(t, w.work, hostilePOSIX347)
	if got := Classify(cand); got != ClassA {
		t.Errorf("class: got %s, want A (the fail-closed tier is the one the windows-tagged 334 teeth pin)", got)
	}
	d := Gate(cand, nil)
	if d.Class != ClassA || d.NeedL2 || d.Allow {
		t.Errorf("Gate: got %s, want class=A needL2=false allow=false (non-overridable)", gateShape216(d))
	}
	// Non-overridable means the B-tier single-file confirmation buys nothing.
	if u := Gate(cand, map[string]bool{normPath(cand): true}); u.Allow || u.Class != ClassA {
		t.Errorf("an A deny from the fail-closed clause was unlocked by an override key: %s", gateShape216(u))
	}
	// AC#3's named reading: the same candidate beside the profile is a plain B
	// deny, so no tooth above is satisfied by a clause that reads the whole
	// machine as A.
	off := filepath.Join(candidateOffAnchorChain334(t), hostilePOSIX347)
	if err := os.Symlink(danglingTarget347, off); err != nil {
		t.Fatalf("setup symlink %s: %v", off, err)
	}
	if got := Classify(off); got != ClassB {
		t.Errorf("off-chain arm: got %s, want B (the branch protects only the candidate's own ancestry chain): %s", got, gateShape216(Gate(off, nil)))
	}
	t.Logf("ANCHORMISS339-OTHER-347 cand=%s class=%s clause=%v offchain=%s", gateShape216(d), d.Class, true, ClassB)
}
