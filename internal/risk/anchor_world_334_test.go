package risk

import (
	"os"
	"path/filepath"
	"testing"
)

// The A-anchor instrument world, shared by the per-platform teeth.
//
// Ticket 347 collapsed the two byte-symmetric copies of this block (one in
// anchor_miss_reason_334_test.go behind //go:build windows, one in
// anchor_miss_reason_339_other_test.go behind //go:build !windows) into one
// untagged file, which is what both of those files said would happen "when the
// teeth move here" (339_other_test.go:43 before this ticket). There is now one
// world builder, so the two platforms cannot drift into "the same expectation"
// being maintained twice.
//
// What still differs per platform is HOW a spelling becomes unproven
// (pathForms.partial), and that is the only thing the platform files keep:
//
//   - Windows: hostile334 carries U+001B / U+202E, which the OS refuses to even
//     parse - Lstat fails with "syntax is incorrect", not "no such file", so
//     pathExists reads the component as present-but-unproven;
//   - POSIX: both characters are legal filename bytes, so the same string is an
//     ordinary missing name and the walk proves the miss (measured: ticket 347
//     landed a real resolver and the ticket-334 shape stayed class B with
//     clause=false). The POSIX partial maker is therefore a DANGLING SYMLINK:
//     Lstat succeeds (the entry exists) while Stat/EvalSymlinks cannot open it
//     (the target does not), which is the same present-but-unproven pair
//     without needing the OS to refuse anything.
//
// hostile334 stays here, untagged, because the windows-tagged teeth are not the
// only reader of the constant and because the comment above is where the
// contrast belongs.

const (
	// U+001B then U+202E (written as its UTF-8 bytes so this file stays
	// ASCII-safe): neither can appear in a Windows path, which is the whole
	// point - Lstat fails with "syntax is incorrect", NOT "no such file".
	noise334a  = "已允许\x1b\xe2\x80\xae"
	hostile334 = "id_ed25519" + noise334a

	// hostilePOSIX347 is the POSIX candidate name: an ordinary B-tier name
	// (id_*), because on POSIX nothing about the NAME can be refused. The
	// unproven-spelling half comes from the dangling symlink built around it.
	hostilePOSIX347 = "id_ed25519"

	// danglingTarget347 is the missing target a POSIX partial-maker link
	// points at. Relative on purpose: it resolves against the link's own
	// directory, so the link is dangling wherever the test's temp root is.
	danglingTarget347 = "wisp347-absent-target"

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
