//go:build windows

package risk

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Ticket 336 AC#2, shape A: one tooth per A-tier branch operand (aEq / aUnder =
// "proven match OR uncertainAnchorMiss"). Ticket 334 pinned that branch as a whole;
// its final judge measured that removing ONE case's branch leaves every shipped
// instrument green while that one shape drops from non-overridable A to overridable
// B (docs/evidence/s1/334-v1-verdict.md, AC#2 half-established). The per-operand
// mutation sweep showing each tooth red is .scratch/wisp/probes/336/r1/logs/.
//
// Per operand: install EVERY OTHER branch-capable leaf (an installed leaf resolves
// to itself, so its root is not ancestry-related to the candidate's and it misses
// honestly instead of hijacking the ordered switch), leave the leaf under test
// absent, and pin both the verdict (A + unoverridable + the reason naming THIS
// literal) and both halves of that case's condition - anchorMatches false AND
// uncertainAnchorMiss true, exactly what swapping aEq/aUnder for anchorMatches
// deletes. The off-chain arm pins the same candidate beside the profile as B, so no
// tooth can be satisfied by a clause that reads the whole machine as A.
//
// The windows tag is a measured limit named in the report: uncertainAnchorMiss needs
// a non-empty root on both sides and pathresolver_other.go's resolveHandle returns
// ("", false), so there the branch cannot fire and no world makes "remove it"
// observable - untagged, these teeth would be permanently red (the ticket 339 shape).

type teeth336Case struct {
	name    string
	anchor  string // anchor spelling, built the way classifyForms builds it
	rel     string // "eq" | "under" - the relation that operand asks for
	literal string // the rule string its case returns, verbatim from blacklist.go
	omit    string // the leaf under test: NOT created
	relCand string // profile-relative dir the candidate is placed in
}

// teeth336Roster comes from `git grep -n 'aEq\|aUnder' <anchor> --
// internal/risk/blacklist.go`: eight branch operands over seven A anchors. The
// %APPDATA%\wisp case ANDs aUnder with aEq, so it contributes two rows, each in the
// world where ITS OWN half hangs on the branch (run teeth-00: from <profile>\work the
// aUnder operand is a CERTAIN miss, <APPDATA>\wisp's root is not ancestry-related to
// it, so the case never fires - that half reads from <APPDATA>\Roaming; the aEq half
// reads inside <APPDATA>\wisp\logs, where aUnder's proven match already holds).
func teeth336Roster() []teeth336Case {
	return []teeth336Case{
		{"git-credentials", "home:.git-credentials", "eq", "~/.git-credentials", ".git-credentials", "work"},
		{"wisp-config-under", "appdata:wisp", "under", "%APPDATA%\\wisp\\config.toml", wispCfgLeaf334, `AppData/Roaming`},
		{"wisp-config-eq", "appdata:wisp:config.toml", "eq", "%APPDATA%\\wisp\\config.toml", wispCfgLeaf334, `AppData/Roaming/wisp/logs`},
		{"ssh", "home:.ssh", "under", "~/.ssh/**", ".ssh/id_installed", "work"},
		{"aws", "home:.aws:credentials", "eq", "~/.aws/credentials", ".aws/credentials", "work"},
		{"kube", "home:.kube:config", "eq", "~/.kube/config", ".kube/config", "work"},
		{"dpapi-protect", "appdata:microsoft:protect", "under", "%APPDATA%\\Microsoft\\Protect\\** (DPAPI master keys)", `AppData/Roaming/Microsoft/Protect/marker.txt`, `AppData/Roaming/Microsoft/Windows`},
		{"local-credentials", "local:microsoft:credentials", "under", "%LOCALAPPDATA%\\Microsoft\\Credentials\\**", `AppData/Local/Microsoft/Credentials/marker.txt`, `AppData/Local`},
	}
}

// teeth336Leaves is every branch-capable A leaf; each world installs all but one.
var teeth336Leaves = []string{
	".git-credentials", wispCfgLeaf334, ".ssh/id_installed", ".aws/credentials",
	".kube/config", `AppData/Roaming/Microsoft/Protect/marker.txt`,
	`AppData/Local/Microsoft/Credentials/marker.txt`,
}

// teeth336Installed is every branch-capable leaf except the one under test.
func teeth336Installed(omit string) []string {
	var out []string
	for _, leaf := range teeth336Leaves {
		if leaf != omit {
			out = append(out, leaf)
		}
	}
	return out
}

// teeth336Anchor resolves a roster spec against the same three env values
// classifyForms reads, joined by the same anchorPath it uses.
func teeth336Anchor(spec, homeS, appDataS, localAppDataS string) string {
	root, rest, _ := strings.Cut(spec, ":")
	parts := strings.Split(rest, ":")
	switch root {
	case "home":
		return anchorPath(homeS, parts...)
	case "appdata":
		return anchorPath(appDataS, parts...)
	}
	return anchorPath(localAppDataS, parts...)
}

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

// TestAnchorFallbackTeethPerAnchor336 is ticket 336 AC#2: removing the fail-closed
// branch of any single load-bearing A-tier operand must go this red, while the same
// worlds stay green on unmutated bytes (the paired control AC#2 requires).
func TestAnchorFallbackTeethPerAnchor336(t *testing.T) {
	for _, tc := range teeth336Roster() {
		installed := teeth336Installed(tc.omit)
		rel := eqPath
		if tc.rel == "under" {
			rel = isUnder
		}
		t.Run(tc.name, func(t *testing.T) {
			cand := teeth336Candidate(t, newWorld334(t, true, installed...), tc)
			d := Gate(cand, nil)
			if d.Class != ClassA {
				t.Errorf("class: got %s, want A: the fail-closed branch of %s carries this verdict", gateShape216(d), tc.literal)
			}
			if d.Allow || d.NeedL2 {
				t.Errorf("this shape stopped being non-overridable: %s", gateShape216(d))
			}
			if u := Gate(cand, map[string]bool{normPath(cand): true}); u.Allow || u.Class != ClassA {
				t.Errorf("an A deny carried by %s was unlocked by an override key: %s", tc.literal, gateShape216(u))
			}
			if !strings.Contains(d.Reason, "A-list sensitive path (non-overridable): "+tc.literal) {
				t.Errorf("the reason must name THIS anchor %q (the switch is ordered: naming another anchor means this case stopped firing): %s", tc.literal, d.Reason)
			}
			if !strings.Contains(d.Reason, "fail-closed, NOT a match") {
				t.Errorf("the reason must admit the case fired on the branch: %s", d.Reason)
			}
			f := formsOf(cand)
			a := anchorForms(teeth336Anchor(tc.anchor, normDir(userHomeDir()),
				normDir(envOr("APPDATA", filepath.Join(userHomeDir(), "AppData", "Roaming"))),
				normDir(envOr("LOCALAPPDATA", filepath.Join(userHomeDir(), "AppData", "Local")))))
			matched, clause := anchorMatches(f, a, rel), uncertainAnchorMiss(f, a)
			t.Logf("TEETH336 operand=%s class=%s matched=%v clause=%v", tc.name, d.Class, matched, clause)
			if matched || !clause {
				t.Errorf("this row must observe the branch alone: anchorMatches=%v uncertainAnchorMiss=%v for %s (a row whose verdict does not hang on the branch is not a tooth)", matched, clause, tc.literal)
			}
		})
		t.Run(tc.name+"/offchain", func(t *testing.T) {
			off := filepath.Join(newWorld334(t, false, installed...).work, hostile334)
			if got := Classify(off); got != ClassB {
				t.Errorf("off-chain arm: got %s, want B (the branch must protect only the candidate's own ancestry chain); shape %s", got, gateShape216(Gate(off, nil)))
			}
		})
	}
}
