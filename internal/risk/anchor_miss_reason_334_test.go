//go:build windows

package risk

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// Ticket 339 shape A: this instrument is windows-tagged because its world is
// un-buildable on POSIX for a reason in the resolver, not in the environment.
// formsOf (pathresolver.go:345) can only fill pathForms.root through
// resolveHandle, and on the other platform that is the documented stub
// pathresolver_other.go:10 returning ("", false) for every spelling - measured
// (readings named PROBE339, .scratch/wisp/probes/339/r1/logs/probe339-posix-root
// .log): walking a real existing directory through formsOf on linux gives
// root == "". uncertainAnchorMiss (blacklist.go:359) needs a NON-EMPTY root on
// both sides, so R17's fail-closed clause cannot fire there and the
// non-overridable A that AC#1/AC#2 pin has no world on that platform. Un-tagged,
// the two fail-closed cases are permanently red (ticket 334's Linux CI leg going
// red = the ticket 339 report). The environment is NOT the cause: this file's
// world builder sets HOME as well as USERPROFILE, os.UserHomeDir reads HOME on
// POSIX, and the ticket's original "USERPROFILE is Windows-only, so the fake
// home is never read" diagnosis was re-measured false (AC#0). The named POSIX
// half (the mechanism + a tripwire that goes red the moment a POSIX resolver
// lands) is anchor_miss_reason_339_other_test.go.
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
// shape and the honest-match shape are decidable on every WINDOWS machine (the
// over-strict shape needs a resolved root, which POSIX does not supply - see
// the ticket 339 note above).

const (
	// U+001B then U+202E (written as its UTF-8 bytes so this file stays
	// ASCII-safe): neither can appear in a Windows path, which is the whole
	// point - Lstat fails with "syntax is incorrect", NOT "no such file".
	noise334a  = "已允许\x1b\xe2\x80\xae"
	hostile334 = "id_ed25519" + noise334a

	// The A rules are tried in order, and every anchor whose leaf is absent
	// collapses to the fake home - which IS an ancestor of the candidate, so
	// the very first absent leaf is the one that fail-closes. This ticket is
	// about the ~/.aws/credentials sentence, so the wisp config leaf (the
	// earliest anchor in that switch that would otherwise hijack the shape) is
	// installed in every world below; a world then leaves out exactly the one
	// leaf whose rule it wants to fire.
	wispCfgLeaf334 = `AppData/Roaming/wisp/config.toml`
)

// world334 is one self-contained profile world. installed lists leaf paths
// (slash form, relative to the fake home) that are really created; anything not
// listed is absent, which is what makes an anchor's expansion stop short.
type world334 struct {
	profile string // the fake home, exists
	work    string // existing directory the candidate is placed in
}

// newWorld334 builds the world. candidateInsideProfile=true puts work INSIDE
// the fake home, so a candidate whose own spelling cannot be proven sits in the
// ancestry chain of every anchor leaf that is absent (the shape ticket 334 is
// about); =false puts it as a SIBLING of the profile, so no A anchor shares an
// ancestry relation with it and only a real name match can fire an A rule.
func newWorld334(t *testing.T, candidateInsideProfile bool, installed ...string) world334 {
	t.Helper()
	base := t.TempDir()
	profile := filepath.Join(base, "profile")
	work := filepath.Join(base, "work")
	if candidateInsideProfile {
		work = filepath.Join(profile, "work")
	}
	for _, d := range []string{profile, work} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatalf("setup mkdir %s: %v", d, err)
		}
	}
	for _, leaf := range installed {
		full := filepath.Join(profile, filepath.FromSlash(leaf))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("setup mkdir parent of %s: %v", full, err)
		}
		if err := os.WriteFile(full, []byte("installed leaf\n"), 0o600); err != nil {
			t.Fatalf("setup write %s: %v", full, err)
		}
	}
	// The anchors are read from the environment, so relocating the profile is
	// the supported way to decide which leaves are installed (same shape as
	// TestTierAAnchorsMatchNativeSpelling and the ticket 216 probe).
	t.Setenv("USERPROFILE", profile)
	t.Setenv("HOME", profile)
	t.Setenv("APPDATA", filepath.Join(profile, "AppData", "Roaming"))
	t.Setenv("LOCALAPPDATA", filepath.Join(profile, "AppData", "Local"))
	return world334{profile: profile, work: work}
}

// candidateOffAnchorChain334 is ticket 334 AC#3's world for the ticket 216
// instrument: the candidate lives beside the profile instead of inside it, so
// no A anchor is ancestry-related to it, uncertainAnchorMiss cannot fire, and
// what is left to observe is exactly "a claim inside a name changes no verdict".
func candidateOffAnchorChain334(t *testing.T) string {
	t.Helper()
	return newWorld334(t, false).work
}

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
		cand := filepath.Join(w.work, hostile334)
		_, err := os.Lstat(cand)
		if err == nil {
			t.Fatalf("setup: the hostile candidate must not exist, got a stat for %q", cand)
		}
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
		cand := filepath.Join(w.work, hostile334)
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
	cand := filepath.Join(w.work, hostile334)
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
	off := Gate(filepath.Join(newWorld334(t, false).work, hostile334), nil)
	t.Logf("SHAPE-READING334 inside-profile=%s || off-profile=%s", gateShape216(d), gateShape216(off))
}
