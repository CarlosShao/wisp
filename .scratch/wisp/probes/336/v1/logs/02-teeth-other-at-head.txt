//go:build !windows

package risk

import (
	"os"
	"path/filepath"
	"testing"
)

// Ticket 336 AC#2, the other platform's half of the same roster, named so that the
// windows-tagged teeth cannot quietly "run less" (the ticket 339 shape: an
// instrument whose world only exists on one OS, discovered later by a different
// leg on a container).
//
// Why there is no tooth here, measured rather than assumed: uncertainAnchorMiss
// requires a NON-EMPTY resolved root on both sides, and this platform's resolver is
// the documented stub - pathresolver_other.go's resolveHandle returns ("", false) -
// so formsOf never fills root for any candidate or any anchor, the branch can never
// fire, and no world on this platform can make "remove this anchor's branch"
// observable. A tooth here would be permanently red, which AC#2 forbids, and a fake
// green would be worse: nothing below asserts the A tier.
//
// This is the tripwire instead. It pins the mechanism, and it goes RED the moment a
// POSIX handle resolver lands (root stops being empty), which is exactly when
// anchor_fallback_teeth_336_windows_test.go must be ported to this platform.

func teeth336OtherAnchors(homeS, appDataS, localAppDataS string) []string {
	return []string{
		anchorPath(homeS, ".git-credentials"),
		anchorPath(appDataS, "wisp"),
		anchorPath(appDataS, "wisp", "config.toml"),
		anchorPath(homeS, ".ssh"),
		anchorPath(homeS, ".aws", "credentials"),
		anchorPath(homeS, ".kube", "config"),
		anchorPath(appDataS, "microsoft", "protect"),
		anchorPath(localAppDataS, "microsoft", "credentials"),
	}
}

// TestFallbackBranchHasNoWorldOnThisPlatform336 is the named POSIX half of ticket
// 336 AC#2: every A-tier anchor and every candidate has an empty resolved root here,
// so the fail-closed branch cannot be load-bearing on any input.
func TestFallbackBranchHasNoWorldOnThisPlatform336(t *testing.T) {
	base := t.TempDir()
	profile := filepath.Join(base, "profile")
	work := filepath.Join(profile, "work")
	for _, d := range []string{profile, work, filepath.Join(base, "work")} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatalf("setup mkdir %s: %v", d, err)
		}
	}
	// Same relocation the windows teeth use. os.UserHomeDir reads HOME here and
	// USERPROFILE on Windows, so the fake home does take effect on this platform -
	// the missing ingredient is the resolved root, not the environment.
	t.Setenv("HOME", profile)
	t.Setenv("USERPROFILE", profile)
	t.Setenv("APPDATA", filepath.Join(profile, "AppData", "Roaming"))
	t.Setenv("LOCALAPPDATA", filepath.Join(profile, "AppData", "Local"))

	homeS := normDir(userHomeDir())
	appDataS := normDir(envOr("APPDATA", filepath.Join(userHomeDir(), "AppData", "Roaming")))
	localAppDataS := normDir(envOr("LOCALAPPDATA", filepath.Join(userHomeDir(), "AppData", "Local")))

	for _, relCand := range []string{"work", "AppData/Roaming", "AppData/Roaming/wisp/logs", "AppData/Local"} {
		dir := work
		if relCand != "work" {
			dir = filepath.Join(profile, filepath.FromSlash(relCand))
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatalf("setup mkdir %s: %v", dir, err)
			}
		}
		cand := filepath.Join(dir, hostile334)
		f := formsOf(cand)
		if f.root != "" {
			t.Errorf("tripwire: this platform now resolves a root (%q) - pathresolver_other.go has grown a real resolver, so the per-anchor fail-closed teeth in anchor_fallback_teeth_336_windows_test.go must be ported here (ticket 336 AC#2) or a POSIX-only downgrade will go unobserved", f.root)
			continue
		}
		for _, spelling := range teeth336OtherAnchors(homeS, appDataS, localAppDataS) {
			a := anchorForms(spelling)
			if a.root != "" {
				t.Errorf("tripwire: anchor %q now resolves a root (%q) on this platform; port the teeth", spelling, a.root)
				continue
			}
			if uncertainAnchorMiss(f, a) {
				t.Errorf("the branch fired with empty roots for anchor %q: uncertainAnchorMiss changed shape, re-read this file", spelling)
			}
		}
		t.Logf("TEETH336-OTHER relCand=%s cand_root=EMPTY verdict=%s (the branch cannot fire here; see pathresolver_other.go)",
			relCand, Classify(cand))
	}
}
