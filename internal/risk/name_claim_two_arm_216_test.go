package risk

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Ticket 216 AC#6 arm B (the two-arm differential). AC#3's half-sentence "no
// verdict may change because a name spells 已允许 inside it" was never
// attacked by any case, so each carrier below runs the same call twice - a
// plain name and the same name carrying the claim plus U+001B plus one Cf -
// and the verdict must come out identical: same tier/level, same rules, same
// reason.
//
// The four carriers are the only spellings that reach this package at all. An
// attachment's display name is NOT among them (the word "attachment" appears
// in no production file of risk/perm/tools/agent), and pinning it here would
// be a permanently green decoration - AC#6 names that shape and forbids it.
//   - a file base name into the A/B blacklist (Classify, Gate);
//   - argv[0] into the shell allowlist (Assess, R6);
//   - a host into the domain allowlist (Assess, R5);
//   - an override key into Gate's single-file confirmation map.
//
// Non-tautology proof: the off-tree mutation readings recorded in
// docs/evidence/s1/216-allow-substring-pin-r2.md.

const (
	claim216      = "已允许"
	claimNoise216 = claim216 + "\x1b\u202e" // plus U+001B and a Cf (U+202E)
)

func withClaim216(name string) string { return name + claimNoise216 }

func assessShape216(d Decision) string {
	return fmt.Sprintf("level=%s rules=%v reason=%q", d.Level, d.RulesHit, d.Reason)
}

func gateShape216(d PathDecision) string {
	return fmt.Sprintf("class=%s allow=%v needL2=%v reason=%q", d.Class, d.Allow, d.NeedL2, d.Reason)
}

// assertSameAssess216 is the two-arm claim in one place: every field equal.
func assertSameAssess216(t *testing.T, label string, plain, hostile Decision) {
	t.Helper()
	if plain.Level != hostile.Level || plain.Reason != hostile.Reason ||
		plain.SessionOverrideBlocked != hostile.SessionOverrideBlocked ||
		!equalRules216(plain.RulesHit, hostile.RulesHit) {
		t.Errorf("%s: the claim moved the verdict\n\tplain   %s\n\thostile %s",
			label, assessShape216(plain), assessShape216(hostile))
	}
}

func equalRules216(a, b []RuleID) bool {
	return strings.Join(toStrings216(a), ",") == strings.Join(toStrings216(b), ",")
}

func toStrings216(in []RuleID) []string {
	out := make([]string, len(in))
	for i, r := range in {
		out[i] = string(r)
	}
	return out
}

// TestFileBaseNameClaimChangesNoClassOrReason is carrier 1, both tiers.
func TestFileBaseNameClaimChangesNoClassOrReason(t *testing.T) {
	// Ticket 334 AC#3 added this world; NO assertion below was touched, and the
	// tier each pair must agree on is still the one the table says.
	//
	// Why it was needed: on the CI surface (GOTMPDIR built from %LOCALAPPDATA%,
	// scripts/portable-tests.sh:152-159, TMP/TEMP untouched) t.TempDir() roots
	// the candidate INSIDE the profile. The hostile arm carries U+001B, which
	// Lstat refuses to parse - "syntax is incorrect", not "no such file" - so
	// pathExists leaves its spelling unproven (partial), and any A anchor whose
	// leaf is not installed here (say ~/.aws) has its own expansion stopped at
	// the home. uncertainAnchorMiss then sees two ancestry-related roots with an
	// unproven side and fail-closes to a non-overridable A. That is a reading of
	// ANCHORS, not of the claim in the name: the instrument was observing
	// something other than what it claims, and going red on CI for it.
	//
	// candidateOffAnchorChain334 relocates the profile beside the candidate, so
	// no A anchor shares an ancestry chain with it and the claim differential is
	// the only thing left observable. The inside-profile shape is not dropped:
	// ticket 334 keeps it as a named reading plus a pinned level (see
	// anchor_miss_reason_334_test.go).
	dir := candidateOffAnchorChain334(t)
	cases := []struct {
		name           string
		plain, hostile string
		want           Class
	}{
		{
			"A tier .git/config", filepath.Join(dir, "repo", ".git", "config"),
			filepath.Join(dir, "repo"+claimNoise216, ".git", "config"), ClassA,
		},
		{
			"B tier id_*", filepath.Join(dir, "id_ed25519"),
			filepath.Join(dir, withClaim216("id_ed25519")), ClassB,
		},
		{
			"B tier .env*", filepath.Join(dir, ".env.local"),
			filepath.Join(dir, withClaim216(".env.local")), ClassB,
		},
		{
			"B tier secrets.*", filepath.Join(dir, "secrets.toml"),
			filepath.Join(dir, withClaim216("secrets.toml")), ClassB,
		},
	}
	// Non-vacuity guard: an ordinary file in the same directory must be none,
	// so the pairs cannot be agreeing because everything reads alike.
	if got := Classify(filepath.Join(dir, "plain-note.txt")); got != ClassNone {
		t.Fatalf("fixture guard: an ordinary file here classified %s, want none", got)
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if got := Classify(tc.plain); got != tc.want {
				t.Fatalf("fixture guard: plain %q classified %s, want %s", tc.plain, got, tc.want)
			}
			if got := Classify(tc.hostile); got != Classify(tc.plain) {
				t.Errorf("class changed by a claim in the name: plain %q = %s, hostile %q = %s",
					tc.plain, Classify(tc.plain), tc.hostile, got)
			}
			plain, hostile := Gate(tc.plain, nil), Gate(tc.hostile, nil)
			if plain != hostile {
				t.Errorf("Gate verdict changed by a claim in the name:\n\tplain   %s\n\thostile %s",
					gateShape216(plain), gateShape216(hostile))
			}
		})
	}
}

// TestArgvZeroClaimChangesNoShellVerdict is carrier 2. Both shapes are pinned:
// a claim inside the program's directory (the program is allowlisted, so the
// verdict must not move) and a claim glued onto an allowlisted program name
// (a different program, so the claim must buy it nothing).
func TestArgvZeroClaimChangesNoShellVerdict(t *testing.T) {
	a := NewRiskAssessor()
	allow := []string{"git.exe", "wisp"}
	t.Run("allowlisted-program-with-claim-in-its-directory", func(t *testing.T) {
		plain := shell216(a, allow, `C:\bin\git.exe`)
		hostile := shell216(a, allow, `C:\bin`+claimNoise216+`\git.exe`)
		if plain.Level != L1 {
			t.Fatalf("fixture guard: plain argv[0] = %s, want the allowlisted L1", assessShape216(plain))
		}
		assertSameAssess216(t, "argv[0] with a claim inside a directory segment", plain, hostile)
	})
	t.Run("claim-glued-onto-an-allowlisted-program-buys-nothing", func(t *testing.T) {
		hostile := shell216(a, allow, withClaim216("git.exe"))
		if hostile.Level != L2 || hostile.Reason != `R6: 程序不在 shell 白名单: `+withClaim216("git.exe") {
			t.Errorf("a program name asserting 已允许 was not left on the not-allowlisted branch: %s",
				assessShape216(hostile))
		}
	})
	t.Run("two-unlisted-programs-one-claiming", func(t *testing.T) {
		plain, hostile := shell216(a, allow, "helper"), shell216(a, allow, withClaim216("helper"))
		if plain.Level != L2 || hostile.Level != plain.Level ||
			!equalRules216(plain.RulesHit, hostile.RulesHit) {
			t.Errorf("the level moved:\n\tplain   %s\n\thostile %s",
				assessShape216(plain), assessShape216(hostile))
		}
		// These two arms name two DIFFERENT programs, so their reasons cannot
		// be the same string - the frozen reason shape is "the program is not
		// in the shell allowlist: <argv[0]>", which quotes the name verbatim.
		// What is pinned is that the only difference is that echo: the claim
		// buys no extra sentence, no downgrade and no upgrade.
		if want := `R6: 程序不在 shell 白名单: helper`; plain.Reason != want {
			t.Errorf("plain reason = %q, want %q", plain.Reason, want)
		}
		if want := `R6: 程序不在 shell 白名单: ` + withClaim216("helper"); hostile.Reason != want {
			t.Errorf("hostile reason = %q, want the verbatim echo of the name only: %q", hostile.Reason, want)
		}
	})
}

// TestHostClaimChangesNoNetworkVerdict is carrier 3.
func TestHostClaimChangesNoNetworkVerdict(t *testing.T) {
	a := NewRiskAssessor()
	allow := []string{"example.com"}
	t.Run("allowlisted-subdomain-with-claim-inside-a-label", func(t *testing.T) {
		plain := net216(a, allow, "api.example.com")
		hostile := net216(a, allow, "api"+claimNoise216+".example.com")
		if plain.Level != L0 {
			t.Fatalf("fixture guard: plain host = %s, want the allowlisted clean verdict", assessShape216(plain))
		}
		assertSameAssess216(t, "host with a claim inside a label", plain, hostile)
	})
	t.Run("claim-glued-onto-an-allowlisted-domain-buys-nothing", func(t *testing.T) {
		hostile := net216(a, allow, withClaim216("example.com"))
		if hostile.Level != L2 || hostile.Reason != "R5: 域名在白名单之外: "+withClaim216("example.com") {
			t.Errorf("a host asserting 已允许 was not left outside the allowlist: %s", assessShape216(hostile))
		}
	})
}

// TestOverrideKeySpelledWithClaimDoesNotUnlock is carrier 4, the allow side of
// the blacklist: Gate's confirmation map, keyed on one exact spelling.
func TestOverrideKeySpelledWithClaimDoesNotUnlock(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "id_ed25519")
	if err := os.WriteFile(file, []byte("private key material"), 0o600); err != nil {
		t.Fatal(err)
	}
	res, err := Resolve(file, nil)
	if err != nil {
		t.Fatalf("Resolve(%s): %v", file, err)
	}
	base := Gate(res.Canonical, nil)
	if base.Class != ClassB || !base.NeedL2 || base.Allow {
		t.Fatalf("fixture guard: the file must sit in the tier that has an override seam, got %s",
			gateShape216(base))
	}
	if d := Gate(res.Canonical, map[string]bool{normPath(res.Canonical): true}); !d.Allow {
		t.Fatalf("fixture guard: the plain override key must unlock its own file, got %s", gateShape216(d))
	}
	claimed := map[string]bool{normPath(withClaim216(res.Canonical)): true}
	if d := Gate(res.Canonical, claimed); d != base {
		t.Errorf("an override recorded under a claiming spelling unlocked the file:\n\twithout %s\n\twith    %s",
			gateShape216(base), gateShape216(d))
	}
}

func shell216(a *RiskAssessor, allow []string, argv0 string) Decision {
	return a.Assess("shell.exec", nil, Facts{ShellArgv: []string{argv0, "status"}, ShellAllowlist: allow})
}

func net216(a *RiskAssessor, allow []string, host string) Decision {
	return a.Assess("web.open", nil, Facts{Network: &NetTarget{Protocol: "https", Host: host, Allowlist: allow}})
}
