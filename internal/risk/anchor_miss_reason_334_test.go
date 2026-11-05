package risk

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// Ticket 339 shape B, completed by ticket 347: this instrument is NOT
// windows-tagged. It was, for one shape (339-r1: //go:build windows plus a
// POSIX tripwire), because its world could not be built here - and the
// condition was in the resolver, not in the environment. formsOf
// (pathresolver.go) fills pathForms.root only through resolveHandle, and
// pathresolver_other.go used to return ("", false) for every spelling, so
// uncertainAnchorMiss (blacklist.go) could not fire on POSIX at all (it needs a
// non-empty root on both sides) and the non-overridable A these cases pin had
// no world on that platform. Ticket 347 landed realpath + lstat there, the
// tripwire went red exactly as designed, and the ticket-339 split collapses
// here: the same three cases below now run on every platform. The environment
// was never the cause - the world builder sets HOME as well as USERPROFILE, and
// os.UserHomeDir reads HOME on POSIX (re-measured for the ticket's AC#0).
//
// What still differs per platform is HOW a spelling becomes unproven
// (pathForms.partial), because the OS decides that: a Windows name carrying
// U+001B / U+202E cannot be parsed at all, while both bytes are ordinary
// filename bytes on POSIX, so there the same "present but unprovable" pair is a
// dangling symlink. That single step is partialCandidate334, supplied per
// platform (anchor_miss_reason_334_windows_test.go, and the POSIX half
// delegating to the maker ticket 347 landed, anchor_miss_reason_339_other_test.go).
// Everything else - the profile world and its installed/absent leaf lists -
// comes from the shared builder in anchor_world_334_test.go, so the two
// platforms cannot drift into two copies of one expectation.
//
// Ticket 334. Two things multiply into one dishonest rejection:
//
//   - pathExists (pathresolver.go) is deliberately fail-toward-exists: only a
//     clean "no such file" reads as absent, so a name Windows refuses to even
//     parse ("The filename, directory name, or volume label syntax is
//     incorrect") leaves the candidate's spelling UNPROVEN (partial=true);
//   - an A-anchor leaf that is not installed on this machine (say ~/.aws) stops
//     ITS expansion at the profile home, so anchor.root is the home.
//
// uncertainAnchorMiss (R17's fail-closed clause) then sees two ancestry-related
// roots with one unproven side and says "treat it as protected" - a DENY, in the
// non-overridable tier. That verdict stays exactly as it is (AC#2). What ticket
// 334 fixes is the sentence printed beside it: it used to name the credential
// file the candidate was never shown to be.
//
// The worlds below are built, not inherited from the host: every A-anchor leaf
// is either installed by name or deliberately absent, so both the over-strict
// shape and the honest-match shape are decidable on every machine.

// quotedFold334 puts a path into the shape aRuleReason prints it in: %q quoting
// (backslashes doubled, control characters escaped) and lowercased, because
// classifyForms speaks in normPath shapes (comparison-folded, so lowercased on
// Windows) while t.TempDir() hands back a mixed-case path.
func quotedFold334(s string) string { return strings.ToLower(strconv.Quote(s)) }

func containFolded334(t *testing.T, haystack, needle, what string) {
	t.Helper()
	if !strings.Contains(strings.ToLower(haystack), quotedFold334(needle)) {
		t.Errorf("reason must name %s (%q), got:\n\t%s", what, needle, haystack)
	}
}

// TestAnchorMissReasonNamesCandidateAndMissingAnchor is ticket 334 AC#1: the
// fail-closed rejection has to say which candidate was rejected and which
// anchor was not installed, instead of accusing the candidate of being the
// credential file that rule guards.
func TestAnchorMissReasonNamesCandidateAndMissingAnchor(t *testing.T) {
	t.Run("absent-aws-leaf-accuses-aws-and-says-so", func(t *testing.T) {
		// .git-credentials and .ssh ARE installed, so those two rules miss
		// honestly and the first fail-closed hit is the ~/.aws one - the exact
		// reason string ticket 334 was written about.
		w := newWorld334(t, true, ".git-credentials", ".ssh/id_installed", wispCfgLeaf334)
		cand := partialCandidate334(t, w.work)
		d := Gate(cand, nil)
		if d.Class != ClassA {
			t.Fatalf("fixture guard: the fail-closed shape must still be A here, got %s", gateShape216(d))
		}
		if d.Reason == "A-list sensitive path (non-overridable): ~/.aws/credentials" {
			t.Errorf("the reason is back to naming only the anchor, which is the false accusation ticket 334 is about: %q", d.Reason)
		}
		containFolded334(t, d.Reason, normPath(cand), "the REAL candidate being rejected")
		containFolded334(t, d.Reason, filepath.Join(w.profile, ".aws", "credentials"), "the REAL absent anchor")
		containFolded334(t, d.Reason, filepath.Join(w.profile, "work"), "where the candidate's own walk stopped")
		// The collapse that makes the two trees comparable at all: the anchor
		// that is not installed has its root stop at the profile home.
		containFolded334(t, d.Reason, w.profile, "the home the un-installed anchor collapsed to")
		for _, want := range []string{"fail-closed", "NOT a match"} {
			if !strings.Contains(d.Reason, want) {
				t.Errorf("reason must carry the word %q so a reader cannot take it for a match, got:\n\t%s", want, d.Reason)
			}
		}
		t.Logf("honest reject: %s", d.Reason)
	})
	t.Run("installing-aws-moves-the-named-anchor-to-kube", func(t *testing.T) {
		// Ticket 334 现量 row 3, built instead of inherited: the reason follows
		// WHICH anchor is missing, never the candidate's identity.
		w := newWorld334(t, true, ".git-credentials", ".ssh/id_installed", ".aws/not_credentials", wispCfgLeaf334)
		cand := partialCandidate334(t, w.work)
		d := Gate(cand, nil)
		if d.Class != ClassA {
			t.Fatalf("fixture guard: want A from the fail-closed clause, got %s", gateShape216(d))
		}
		containFolded334(t, d.Reason, filepath.Join(w.profile, ".kube", "config"), "the now-missing kube anchor")
		containFolded334(t, d.Reason, normPath(cand), "the REAL candidate")
		if strings.Contains(strings.ToLower(d.Reason), quotedFold334(filepath.Join(w.profile, ".aws", "credentials"))) {
			t.Errorf("the aws anchor is installed here, it must not be the one the reason names: %s", d.Reason)
		}
	})
	t.Run("proven-match-keeps-the-bare-literal", func(t *testing.T) {
		// The honest half of AC#1: a candidate that really IS the credential
		// file gets the rule's own literal, unchanged, with no fail-closed
		// sentence bolted on.
		w := newWorld334(t, true, ".git-credentials", ".ssh/id_installed", ".aws/credentials", wispCfgLeaf334)
		cand := filepath.Join(w.profile, ".aws", "credentials")
		if _, err := os.Lstat(cand); err != nil {
			t.Fatalf("setup: %v", err)
		}
		d := Gate(cand, nil)
		if d.Class != ClassA {
			t.Fatalf("fixture guard: the real credential file must still be A, got %s", gateShape216(d))
		}
		if want := "A-list sensitive path (non-overridable): ~/.aws/credentials"; d.Reason != want {
			t.Errorf("an honest match must keep the bare rule literal verbatim\n\tgot  %q\n\twant %q", d.Reason, want)
		}
	})
}

// TestAnchorMissKeepsNonOverridableClassAndNeedL2 is ticket 334 AC#2's gate.
// Relaxing this shape is a contract change and is NOT in this ticket's reach
// (see the ticket's "未定案那一支"), so the three values below are the reading
// the reject had before this commit and must keep after it - and the positive
// control (mutation M-C in .scratch/wisp/probes/334/r1/) proves this test goes
// red the moment anyone downgrades it, even by one tier.
func TestAnchorMissKeepsNonOverridableClassAndNeedL2(t *testing.T) {
	w := newWorld334(t, true, ".git-credentials", ".ssh/id_installed", wispCfgLeaf334)
	cand := partialCandidate334(t, w.work)
	if got := Classify(cand); got != ClassA {
		t.Errorf("class: got %s, want A (the fail-closed tier is unchanged by AC#1)", got)
	}
	d := Gate(cand, nil)
	if d.Class != ClassA || d.NeedL2 || d.Allow {
		t.Errorf("Gate: got %s, want class=A needL2=false allow=false (non-overridable)", gateShape216(d))
	}
	// Non-overridable means the B-tier single-file confirmation buys nothing.
	if u := Gate(cand, map[string]bool{normPath(cand): true}); u.Allow || u.Class != ClassA {
		t.Errorf("an A deny from the fail-closed clause was unlocked by an override key: %s", gateShape216(u))
	}
	// AC#3's named reading, kept as a reading and NOT as a permanently red
	// assertion: the same input beside the profile is a plain B deny. The
	// over-strictness inside the profile chain is real, it is today's behavior,
	// and its denominator is still unmeasured (the ticket's 未定案那一支).
	off := Gate(partialCandidate334(t, newWorld334(t, false).work), nil)
	t.Logf("SHAPE-READING334 inside-profile=%s || off-profile=%s", gateShape216(d), gateShape216(off))
}
