//go:build !windows

package risk

import (
	"os"
	"path/filepath"
	"testing"
)

// Ticket 339 AC#1 shape A, the POSIX half. The ticket 334 instrument
// (anchor_miss_reason_334_test.go) is now tagged //go:build windows, because
// on this platform its world is un-buildable and the condition is in the
// resolver, not in the environment:
//
//   - resolveHandle (pathresolver_other.go:10) is the documented stub, it
//     returns ("", false) for EVERY spelling - the leg probe walked a real,
//     existing directory through formsOf and got root == "" (reading
//     ANCHORMISS339-PROBE, .scratch/wisp/probes/339/r1/logs/linux-ac0-*),
//   - uncertainAnchorMiss (blacklist.go:359) requires a non-empty root on
//     BOTH sides, so R17's fail-closed clause cannot fire on any input here,
//   - so the non-overridable A verdict ticket 334 AC#1/AC#2 pin has no world
//     on this platform, and t.Setenv cannot manufacture one: the fake home
//     DOES take effect here (HOME is set by the world builder and os.UserHomeDir
//     reads it - the ticket's USERPROFILE diagnosis was re-measured false at
//     AC#0).
//
// A permanent red is the shape ticket 336 AC#2 forbids, and a green that
// asserts nothing is worse, so this file asserts the MECHANISM instead: it
// goes red the moment this platform grows a real resolver (a root appears),
// which is exactly when the windows-tagged teeth must be ported here and the
// twin fixtures below deleted. It also pins today's honest reading - the
// fail-closed shape classifies B here - so the claim cannot rot quietly into
// a false green either way.
//
// The fixture block is a copy, not a refactor: the untagged
// name_claim_two_arm_216_test.go and the !windows
// anchor_fallback_teeth_336_other_test.go read hostile334, wispCfgLeaf334 and
// candidateOffAnchorChain334 on this platform, and tagging the whole 334 file
// without re-declaring them breaks the linux build (measured with go vet:
// undefined: hostile334, undefined: candidateOffAnchorChain334; log
// jia-naive-buildbreak). Keep this block byte-symmetric with the twin file.
// When the teeth move here, both copies collapse back into one untagged file.

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

// TestAnchorMissHasNoWorldOnThisPlatform339 is the named POSIX half of ticket
// 339 AC#1 shape A: every spelling resolves root == "" here, so
// uncertainAnchorMiss is inert, and the A verdict ticket 334 pins cannot be
// produced by any test-side world. This is the tripwire, not a tooth: it goes
// red when a POSIX handle resolver lands (root stops being empty), which is
// exactly when anchor_miss_reason_334_test.go's teeth must be ported here.
func TestAnchorMissHasNoWorldOnThisPlatform339(t *testing.T) {
	w := newWorld334(t, true, ".git-credentials", ".ssh/id_installed", wispCfgLeaf334)
	cand := filepath.Join(w.work, hostile334)
	if _, err := os.Lstat(cand); err == nil {
		t.Fatalf("setup: the hostile candidate must not exist, got a stat for %q", cand)
	}
	// The reading below is only about the resolver if the environment still
	// does its part: the fake home must be the home the package actually sees
	// (os.UserHomeDir reads HOME on this platform, the ticket's own diagnosis
	// of "HOME is not read here" was re-measured false at anchor_miss_reason
	// _334_test.go:88 / AC#0).
	homeS := normDir(userHomeDir())
	if homeS != normPath(w.profile) {
		t.Errorf("fixture: relocating the fake home stopped taking effect on this platform (home %q, profile %q); the readings below would then measure the environment, not the resolver", homeS, w.profile)
	}
	f := formsOf(cand)
	a := anchorForms(anchorPath(homeS, ".aws", "credentials"))
	if f.root != "" {
		t.Errorf("tripwire: this platform now resolves a candidate root (%q) - pathresolver_other.go has grown a real resolver, so the windows-tagged fail-closed teeth in anchor_miss_reason_334_test.go must be ported here (ticket 339 AC#1 shape A's own condition)", f.root)
	}
	if a.root != "" {
		t.Errorf("tripwire: this platform now resolves an anchor root (%q); port the teeth named above", a.root)
	}
	if uncertainAnchorMiss(f, a) {
		t.Errorf("the fail-closed clause fired with these roots on this platform; uncertainAnchorMiss changed shape, re-read blacklist.go and port the teeth")
	}
	if got := Classify(cand); got != ClassB {
		t.Errorf("class: got %s, want B (with an inert fail-closed clause the hostile leaf only reaches the B tier here; if it reaches A the premise of the windows tag broke - port the teeth and delete this file)", got)
	}
	d := Gate(cand, nil)
	t.Logf("ANCHORMISS339-OTHER cand.root=EMPTY anchor.root=EMPTY clause=false class=%s shape=%s (ticket 334's non-overridable A has no world on this platform; see pathresolver_other.go:10)", d.Class, gateShape216(d))
}
