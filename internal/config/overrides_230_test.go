package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// Ticket 230 AC#4 (owner ruling 2026-11-05, ledger A1549): risk.blacklist_overrides
// is a REAL setting now - the persisted, user-customizable B-tier exemption list.
// The three readings this file pins:
//
//  1. a valid B-tier entry survives the load pipeline as an exemption list, in the
//     C26-canonical spelling the gate compares against (not the raw string);
//  2. an A-tier entry is REFUSED at load. A-tier is absolutely denied and no
//     exemption, grant or config line can unlock it (SPEC-06 4.1); accepting the
//     line would be the lying option unwired.go used to reject for the opposite
//     reason ("nothing reads this key");
//  3. an entry no B-tier rule can ever honour is refused too, so a typo cannot
//     sit in the file reading as a standing exemption.

// pinOverrideAnchors230 points every environment anchor the A/B tables read at
// one fake profile, so the AC's examples can be spelled the way an operator
// spells them ("~/.ssh/id_rsa", "%APPDATA%\wisp\config.toml") without touching
// the machine's real home. Both spellings of home are set because
// os.UserHomeDir answers from $HOME on POSIX and %USERPROFILE% on Windows.
func pinOverrideAnchors230(t *testing.T, root string) {
	t.Helper()
	t.Setenv("HOME", root)
	t.Setenv("USERPROFILE", root)
	t.Setenv("APPDATA", filepath.Join(root, "AppData", "Roaming"))
	t.Setenv("LOCALAPPDATA", filepath.Join(root, "AppData", "Local"))
}

// tomlLiteral230 quotes one path as a TOML literal string: Windows paths are
// full of backslashes and a basic string would read them as escapes.
func tomlLiteral230(p string) string { return "'" + p + "'" }

// overrideConfig230 builds a minimal config.toml carrying exactly the given
// override entries.
func overrideConfig230(entries ...string) string {
	quoted := make([]string, len(entries))
	for i, e := range entries {
		quoted[i] = tomlLiteral230(e)
	}
	return "schema_version = 2\n\n[risk]\nblacklist_overrides = [" + strings.Join(quoted, ", ") + "]\n"
}

// TestBlacklistOverridesParseAsAnExemptionList is the positive half. The load
// must succeed (this key no longer lies), and what the Config carries afterwards
// must be the C26-canonical spelling of each entry - the same form the gate's
// comparison keys are built from - not the string the operator typed.
func TestBlacklistOverridesParseAsAnExemptionList(t *testing.T) {
	root := t.TempDir()
	pinOverrideAnchors230(t, root)

	// A real B-tier file (".env" is the canonical example of the tier) plus a
	// home-relative spelling that only exists after expansion.
	project := filepath.Join(root, "project")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatal(err)
	}
	envFile := filepath.Join(project, ".env")
	if err := os.WriteFile(envFile, []byte("TOKEN=fixture\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	homeRelative := filepath.Join(root, "later", ".env")

	path := writeConfigFile(t, overrideConfig230(envFile, "~/later/.env"))
	c, _, err := LoadFile(path, nil)
	if err != nil {
		t.Fatalf("a valid B-tier exemption must load: %v", err)
	}
	if len(c.Risk.BlacklistOverrides) != 2 {
		t.Fatalf("exemption list lost entries: %v", c.Risk.BlacklistOverrides)
	}

	// The expectation is derived through the same resolver the production
	// pipeline uses, so the assertion is "the value is C26-canonical", not
	// "the value is whatever this test happens to think a path looks like".
	want := make([]string, 0, 2)
	for _, spelling := range []string{envFile, homeRelative} {
		res, err := risk.Resolve(spelling, nil)
		if err != nil {
			t.Fatalf("Resolve(%q): %v", spelling, err)
		}
		want = append(want, res.Canonical)
	}
	if !equalStrings(c.Risk.BlacklistOverrides, want) {
		t.Fatalf("exemption list is not the canonical spelling:\n got=%v\nwant=%v",
			c.Risk.BlacklistOverrides, want)
	}
	// Every entry that survives the load must name a path the gate treats as
	// B tier; otherwise the "exemption" could never fire and the setting would
	// be lying again, one layer down.
	for _, entry := range c.Risk.BlacklistOverrides {
		if got := risk.Classify(entry); got != risk.ClassB {
			t.Fatalf("accepted entry %q classifies %s, want B", entry, got)
		}
	}
}

// TestBlacklistOverridesRefuseATier is the negative half and the one the AC
// turns on: every example the ticket lists as "never overridable" must be
// refused by name, INCLUDING when it is spelled with ~ or %APPDATA% (the
// expansion has to happen before the tier is decided, or ".ssh/id_rsa" would
// read as the B-tier "id_*" name pattern and slip through).
func TestBlacklistOverridesRefuseATier(t *testing.T) {
	root := t.TempDir()
	pinOverrideAnchors230(t, root)

	cases := []struct {
		name  string
		entry string
	}{
		{"ssh key spelled with ~", "~/.ssh/id_rsa"},
		{"ssh key spelled absolute", filepath.Join(root, ".ssh", "id_rsa")},
		{"git credentials", filepath.Join(root, ".git-credentials")},
		{"aws credentials", filepath.Join(root, ".aws", "credentials")},
		{"kube config", filepath.Join(root, ".kube", "config")},
		{"browser credential store", filepath.Join(root, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Login Data")},
		{"dpapi master keys", filepath.Join(root, "AppData", "Roaming", "Microsoft", "Protect", "S-1-5-21-0", "master")},
		{"windows credential vault", filepath.Join(root, "AppData", "Local", "Microsoft", "Credentials", "vault")},
		// Spelled with the %APPDATA% token and the PLATFORM separator: on POSIX a
		// backslash is an ordinary file-name character, so a literal
		// "%APPDATA%\wisp\config.toml" spelling there names one long file name
		// and is legitimately NOT the config file.
		{"wisp's own config", "%APPDATA%" + string(filepath.Separator) + filepath.Join("wisp", "config.toml")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := writeConfigFile(t, overrideConfig230(tc.entry))
			_, _, err := LoadFile(path, nil)
			if err == nil {
				t.Fatalf("A-tier entry %q was accepted as a B-tier exemption", tc.entry)
			}
			msg := err.Error()
			if !strings.Contains(msg, "risk.blacklist_overrides") {
				t.Errorf("error must name the key, got %q", msg)
			}
			if !strings.Contains(msg, "A-tier") {
				t.Errorf("error must name the tier that refused it, got %q", msg)
			}
			if !strings.Contains(msg, "never") {
				t.Errorf("error must say A-tier is never overridable, got %q", msg)
			}
			// The rejection must be the A-tier verdict, not the old
			// "this key does nothing" guard: that sentence would mean the
			// registration at unwired.go was left in place.
			if strings.Contains(msg, "does nothing") {
				t.Errorf("the unwired loud-failure still owns this key, got %q", msg)
			}
		})
	}
}

// TestBlacklistOverridesRefuseAnInertEntry: an entry that names no B-tier shape
// can never fire (B membership is decided by the file name, so it cannot become
// true later). Accepting it would put a standing "exemption" in the file that
// exempts nothing.
func TestBlacklistOverridesRefuseAnInertEntry(t *testing.T) {
	root := t.TempDir()
	pinOverrideAnchors230(t, root)
	path := writeConfigFile(t, overrideConfig230(filepath.Join(root, "notes.txt")))
	_, _, err := LoadFile(path, nil)
	if err == nil {
		t.Fatal("an entry no B-tier rule matches must not load")
	}
	if msg := err.Error(); !strings.Contains(msg, "risk.blacklist_overrides") ||
		!strings.Contains(msg, "no B-tier rule") || strings.Contains(msg, "does nothing") {
		t.Fatalf("error must be the tier verdict naming the key (not the old unwired guard), got %q", msg)
	}
}

// TestATierStaysDeniedEvenWhenItReachesTheOverrideSet is the belt under the
// braces: the load-time refusal above is the fast path, but "A-tier is never
// overridable" must also hold for a set that was built some other way. This
// keys the override map with the strongest candidate spellings (the canonical
// form AND its lowercased form) and asserts the gate still answers ClassA with
// no allow and no L2 route, and still emits its A-tier security line.
//
// The comparison fold itself is not reproducible from outside internal/risk,
// which is exactly the seam this wave reports (docs/evidence/s1/230-ac4-*.md):
// the wiring wave must not re-implement it.
func TestATierStaysDeniedEvenWhenItReachesTheOverrideSet(t *testing.T) {
	root := t.TempDir()
	pinOverrideAnchors230(t, root)

	ssh := filepath.Join(root, ".ssh", "id_rsa")
	res, err := risk.Resolve(ssh, nil)
	if err != nil {
		t.Fatalf("Resolve(%q): %v", ssh, err)
	}
	if got := risk.Classify(res.Canonical); got != risk.ClassA {
		t.Fatalf("fixture %q classifies %s, want A", res.Canonical, got)
	}

	var lines []string
	prev := risk.Logf
	risk.Logf = func(format string, args ...any) { lines = append(lines, format) }
	t.Cleanup(func() { risk.Logf = prev })

	d := risk.Gate(res.Canonical, map[string]bool{
		res.Canonical:                  true,
		strings.ToLower(res.Canonical): true,
	})
	if d.Class != risk.ClassA || d.Allow || d.NeedL2 {
		t.Fatalf("A-tier path got %+v; A must stay denied, non-allowed and not L2-unlockable", d)
	}
	if len(lines) == 0 {
		t.Fatal("the A-tier deny emitted no security line")
	}
}
