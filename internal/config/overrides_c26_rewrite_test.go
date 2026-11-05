package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestBlacklistOverridesRefuseARewrittenEntry is the ticket 102 rewrite account
// applied to ticket 230 AC#4's exemption list: the list is PUBLISHED as the key
// risk.GateStanding matches, so it may only carry a spelling the operator
// actually named. C26's expansion step (%VAR%, $VAR, leading ~) stays in the
// pipeline - classification needs it - but an entry it had to substitute into
// is refused at load instead of being canonicalized: keying the gate on the
// expanded tree would exempt a file the config line never spelled, while the
// line still reads as that spelling.
//
// Both halves matter: a silent rejection would be its own kind of lie (the
// operator writes a line, gets an error that does not say why), so this also
// pins that the message names the rewrite C26 performed.
func TestBlacklistOverridesRefuseARewrittenEntry(t *testing.T) {
	root := t.TempDir()
	pinOverrideAnchors230(t, root)
	const envName = "WISP230_REWRITE_ROOT"
	t.Setenv(envName, root)

	project := filepath.Join(root, "project")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, ".env"), []byte("TOKEN=fixture\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	sep := string(filepath.Separator)
	cases := []struct {
		name  string
		entry string
		kind  string
	}{
		{"env variable", "%" + envName + "%" + sep + filepath.Join("project", ".env"), "env"},
		{"home tilde", "~" + sep + filepath.Join("later", ".env"), "home"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := writeConfigFile(t, overrideConfig230(tc.entry))
			_, _, err := LoadFile(path, nil)
			if err == nil {
				t.Fatalf("a rewritten exemption %q loaded; its key would name the expanded tree, not the spelling", tc.entry)
			}
			msg := err.Error()
			if !strings.Contains(msg, "risk.blacklist_overrides") {
				t.Errorf("error must name the key, got %q", msg)
			}
			if !strings.Contains(msg, "cannot be a standing exemption") {
				t.Errorf("the refusal must be the rewrite account, not a tier verdict, got %q", msg)
			}
			if !strings.Contains(msg, "expands to") || !strings.Contains(msg, tc.kind) {
				t.Errorf("error must report the rewrite (%q) it refused, got %q", tc.kind, msg)
			}
			if !strings.Contains(msg, "literal path") {
				t.Errorf("error must say how to fix it (spell the literal path), got %q", msg)
			}
			if strings.Contains(msg, "does nothing") {
				t.Errorf("the unwired loud-failure must not own this key any more, got %q", msg)
			}
		})
	}
}
