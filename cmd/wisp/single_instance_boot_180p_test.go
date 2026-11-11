package main

// The deferred-disposition 20261110 #1 tests: [app] single_instance reaches the
// boot layout through the proc.WithLayout seam.
//
// THE REVERSE CONTRACT (摘喂线 -> 红): every case below asserts the READ half,
// not just the answer half. If the config read is gutted, the false case's
// cleared-mutex-name assertions go red (the layout keeps proc's prod names);
// if the clear-the-names half is gutted, the name assertions red while the
// flag still passes; if the fail-closed direction is inverted (missing /
// unreadable config dissolving the mutex), the two fail-closed cases red.

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/proc"
)

// singleInstance180PProdLayout is a mutex-carrying layout that belongs to a
// throwaway root: proc.LayoutFor's prod fork verbatim, with the data dir
// pointed at the test's own tree (the helper reads config from the path the
// caller hands it, so no real user data dir is ever touched).
func singleInstance180PProdLayout(t *testing.T, root string) proc.Layout {
	t.Helper()
	l, err := proc.LayoutFor(buildinfo.EnvProd, root)
	if err != nil {
		t.Fatalf("LayoutFor(prod): %v", err)
	}
	if !l.MutexEnabled || l.MutexName == "" || l.ActivateEventName == "" {
		t.Fatalf("the premise moved: the prod layout no longer carries a mutex (%+v)", l)
	}
	l.DataDir = root
	return l
}

func singleInstance180PWriteConfig(t *testing.T, root, body string) string {
	t.Helper()
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatalf("mkdir data root: %v", err)
	}
	p := filepath.Join(root, configFileName)
	if body != "" {
		if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
			t.Fatalf("write config: %v", err)
		}
	}
	return p
}

// Test180PSingleInstanceFalseClearsTheMutex is the forward control: the key
// says false, and the boot layout arrives with the mutex flag AND both names
// cleared (proc's own invariant: MutexName/ActivateEventName are "" when
// MutexEnabled is false). The name assertions are the read half - a helper
// that only flips the flag would leave a mutex NAME somebody could still
// acquire with, and these two lines are what catch it.
func Test180PSingleInstanceFalseClearsTheMutex(t *testing.T) {
	root := t.TempDir()
	cfgPath := singleInstance180PWriteConfig(t, root, "schema_version = 2\n\n[app]\nsingle_instance = false\n")
	l, off, source := singleInstanceBootLayoutRead180P(singleInstance180PProdLayout(t, root), cfgPath)
	if !off {
		t.Fatal("the key said false and the helper did not turn the mutex off: the feed does not reach the seam")
	}
	if l.MutexEnabled {
		t.Error("MutexEnabled stayed true off a false file - proc.Boot would take the mutex anyway")
	}
	if l.MutexName != "" || l.ActivateEventName != "" {
		t.Errorf("the mutex names survived the flip (name=%q event=%q): a layout that says no mutex must not still carry a mutex name",
			l.MutexName, l.ActivateEventName)
	}
	if source != cfgPath {
		t.Errorf("source = %q, want the file the decision was read from (%q)", source, cfgPath)
	}
	if l.DataDir != root {
		t.Errorf("DataDir moved: %q, want %q (the key must not relocate the data dir)", l.DataDir, root)
	}
}

// Test180PSingleInstanceTrueAndMissingKeepTheMutex pins the default direction:
// the key true (the schema default) and a MISSING config.toml both leave the
// layout exactly as the env fork built it - same flag, same names - so a
// no-config first run keeps the single-instance guarantee it never got to read.
func Test180PSingleInstanceTrueAndMissingKeepTheMutex(t *testing.T) {
	root := t.TempDir()
	base := singleInstance180PProdLayout(t, root)
	for _, tc := range []struct {
		name, body string
		write      bool
	}{
		{"key-true", "schema_version = 2\n\n[app]\nsingle_instance = true\n", true},
		{"no-key", "schema_version = 2\n", true},
		{"file-missing", "", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfgPath := singleInstance180PWriteConfig(t, root, tc.body)
			if !tc.write {
				_ = os.Remove(cfgPath)
			}
			l, off, source := singleInstanceBootLayoutRead180P(base, cfgPath)
			if off {
				t.Fatalf("%s: the mutex was turned off without the key saying false", tc.name)
			}
			if source != "" {
				t.Fatalf("%s: the boot line would print with no config decision behind it (source=%q)", tc.name, source)
			}
			if !l.MutexEnabled || l.MutexName != base.MutexName || l.ActivateEventName != base.ActivateEventName {
				t.Fatalf("%s: the layout moved off the env fork (mutex=%v name=%q event=%q)",
					tc.name, l.MutexEnabled, l.MutexName, l.ActivateEventName)
			}
		})
	}
}

// Test180PSingleInstanceUnreadableConfigFailsClosed is the fail-closed
// direction: a config file that exists but cannot be parsed must not dissolve
// the mutex it cannot even name. The refusal belongs to the boot (proc.Boot
// reads nothing here); this helper's only job is to keep the lock ON.
func Test180PSingleInstanceUnreadableConfigFailsClosed(t *testing.T) {
	root := t.TempDir()
	base := singleInstance180PProdLayout(t, root)
	cfgPath := singleInstance180PWriteConfig(t, root, "schema_version = 2\n\n[[[\nnot a config\n")
	l, off, _ := singleInstanceBootLayoutRead180P(base, cfgPath)
	if off || l.MutexEnabled != base.MutexEnabled || l.MutexName != base.MutexName {
		t.Fatalf("an unreadable config moved the mutex shape (off=%v mutex=%v name=%q): fail-closed means keep the env fork verbatim",
			off, l.MutexEnabled, l.MutexName)
	}
	// The premise, said out loud: the same bytes ARE a parse failure for the
	// loader, so this case is not silently testing an empty file.
	if _, _, err := config.LoadFile(cfgPath, nil); err == nil {
		t.Fatal("the premise moved: the planted bytes now load cleanly, so this case pins nothing")
	}
}
