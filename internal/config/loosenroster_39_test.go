package config

// Ticket 39 AC#2's roster teeth (write leg 2026-11-18).
//
// The claim: the named roster (loosenroster_39.go) and the D36 direction
// auditors (manager.go's four *Direction functions) name the SAME set of
// loosening keys - every key an auditor can emit as a loosening is covered by
// the roster, every static roster key is emittable by its auditor, every pin
// re-reads its source line, and every roster row names a locked section.
// Deleting a roster row, retyping a Source line, or teaching an auditor a new
// loosening without registering it all redden one of these.

import (
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// rosterSourceLine re-reads the source file the roster names and returns the
// claimed line. The working directory of an in-package test is the package
// directory, so the relative Sources resolve to the auditors themselves.
func rosterSourceLine(t *testing.T, source string) string {
	t.Helper()
	file, lineNo, ok := strings.Cut(source, ":")
	if !ok {
		t.Fatalf("roster source %q is not file:line", source)
	}
	n, err := strconv.Atoi(lineNo)
	if err != nil {
		t.Fatalf("roster source %q has no line number", source)
	}
	raw, err := os.ReadFile(filepath.FromSlash(file))
	if err != nil {
		t.Fatalf("re-reading the roster source %s: %v", source, err)
	}
	lines := strings.Split(string(raw), "\n")
	if n < 1 || n > len(lines) {
		t.Fatalf("roster source %s points at line %d, but %s has %d lines",
			source, n, file, len(lines))
	}
	return lines[n-1]
}

// The pin teeth: every roster row's Source line must still carry the key
// literal (or the auditor's own concatenation fragment for the dynamic
// families). A roster that has rotted away from the auditors - a key moved, a
// line inserted, a spelling changed - fails here with the exact row named.
func Test39LoosenRosterSourcesReReadTheAuditors(t *testing.T) {
	if len(SecurityLoosenRoster) == 0 {
		t.Fatal("the security-loosen roster is empty")
	}
	for _, k := range SecurityLoosenRoster {
		line := rosterSourceLine(t, k.Source)
		if !strings.Contains(line, k.Pin) {
			t.Errorf("roster row %s pins %q at %s, but that line is:\n\t%s\nthe roster has rotted away from its auditor",
				k.Key, k.Pin, k.Source, strings.TrimSpace(line))
		}
	}
}

// The coverage tooth, auditor side: drive each direction auditor with
// all-loosening flips and require every emitted loosening key to pass
// IsSecurityLoosenKey. An auditor that learns a new loosening without the
// roster learning it fails here (and the write gate refuses it fail-closed).
func Test39LoosenRosterCoversEveryAuditorEmittedLoosening(t *testing.T) {
	riskLoosen, _ := riskDirection(&RiskSection{}, &RiskSection{
		ShellEnabled:       true,
		AllowShellString:   true,
		L1WindowSec:        1,
		PermissionMode:     "auto_approve",
		ShellAllowlist:     []string{"x"},
		BlacklistOverrides: []string{"y"},
	})
	fsLoosen, _ := fsDirection(&FSSection{}, &FSSection{
		AllowedDirs:            []string{`C:\x`},
		ReparsePointExceptions: []string{`C:\y`},
		DeleteEnabled:          true,
	})
	netLoosen, _ := netDirection(
		&NetSection{BlockPrivateRanges: true, Proxy: ProxyConfig{Mode: ProxyNone}},
		&NetSection{Allowlist: []string{"example.com"}, Proxy: ProxyConfig{Mode: "manual"}})
	pluginsLoosen, _ := pluginsDirection(
		&PluginsSection{Entries: map[string]PluginEntry{"dyn": {}}},
		&PluginsSection{
			Tier2Enabled: true,
			Entries: map[string]PluginEntry{
				"dyn": {Enabled: true, Capabilities: []string{"fs.read"},
					NetAllowlist: []string{"a"}, HostAPI: []string{"b"}},
			},
		})

	all := slices.Concat(riskLoosen, fsLoosen, netLoosen, pluginsLoosen)
	if len(all) == 0 {
		t.Fatal("no auditor emitted any loosening; the driver is broken")
	}
	for _, want := range []string{"risk.shell_enabled", "fs.allowed_dirs", "net.allowlist",
		"net.proxy.mode", "plugins.tier2_enabled", "plugins.dyn.enabled"} {
		if !slices.Contains(all, want) {
			t.Fatalf("the drivers did not emit %q: %v", want, all)
		}
	}
	for _, k := range all {
		if !IsSecurityLoosenKey(k) {
			t.Errorf("auditor-emitted loosening key %q is not covered by the roster; the roster and the auditors disagree", k)
		}
	}
}

// The coverage tooth, roster side: every STATIC roster key is emittable by its
// own auditor in the loosening direction, and every roster row names a section
// the tier registry marks locked. A row naming a key no auditor emits (or a
// section that is not locked) is an invented category - exactly what the
// dispatch forbids.
func Test39LoosenRosterKeysAreAuditorEmittedInLockedSections(t *testing.T) {
	riskLoosen, _ := riskDirection(&RiskSection{}, &RiskSection{
		ShellEnabled:       true,
		AllowShellString:   true,
		L1WindowSec:        1,
		PermissionMode:     "auto_approve",
		ShellAllowlist:     []string{"x"},
		BlacklistOverrides: []string{"y"},
	})
	fsLoosen, _ := fsDirection(&FSSection{}, &FSSection{
		AllowedDirs:            []string{"x"},
		ReparsePointExceptions: []string{"y"},
		DeleteEnabled:          true,
	})
	netLoosen, _ := netDirection(
		&NetSection{BlockPrivateRanges: true, Proxy: ProxyConfig{Mode: ProxyNone}},
		&NetSection{Allowlist: []string{"a"}, Proxy: ProxyConfig{Mode: "manual"}})
	pluginsLoosen, _ := pluginsDirection(&PluginsSection{}, &PluginsSection{
		Tier2Enabled: true,
		Entries: map[string]PluginEntry{"dyn": {Enabled: true, Capabilities: []string{"c"},
			NetAllowlist: []string{"n"}, HostAPI: []string{"h"}}}})
	emitted := slices.Concat(riskLoosen, fsLoosen, netLoosen, pluginsLoosen)

	for _, k := range SecurityLoosenRoster {
		if tier, ok := TierRegistry[k.Section]; !ok || tier != TierLocked {
			t.Errorf("roster row %s names section %q, which the tier registry does not mark locked (tier=%q ok=%v)",
				k.Key, k.Section, tier, ok)
		}
		if strings.Contains(k.Key, "<id>") {
			continue // dynamic family: matched by the width test below
		}
		if !slices.Contains(emitted, k.Key) {
			t.Errorf("static roster key %q is not emittable by its auditor as a loosening; the row is an invention", k.Key)
		}
	}
}

// The matcher's width: static keys and plugin-family spellings match, the
// family templates themselves do not, and an unknown path answers false (the
// write gate refuses an unknown loosen key fail-closed instead of queueing a
// card for a name the roster cannot carry).
func Test39IsSecurityLoosenKeyWidth(t *testing.T) {
	l1 := "risk." + "l1_" + "window_" + "sec"
	tier2 := "plugins." + "tier" + "2_" + "enabled"
	yes := []string{
		"risk.shell_enabled", l1, "risk.permission_mode",
		"risk.shell_allowlist", "risk.blacklist_overrides",
		"fs.allowed_dirs", "fs.delete_enabled", "fs.reparse_point_exceptions",
		"net.allowlist", "net.proxy.mode", "net.block_private_ranges",
		tier2,
		"plugins.dyn", "plugins.dyn.enabled", "plugins.dyn.capabilities",
		"plugins.dyn.net_allowlist", "plugins.dyn.host_api",
	}
	for _, k := range yes {
		if !IsSecurityLoosenKey(k) {
			t.Errorf("IsSecurityLoosenKey(%q) = false, want true", k)
		}
	}
	no := []string{
		"", "ball.size", "voice.asr", "app.theme",
		"plugins.", "plugins.<id>", "plugins.<id>.enabled",
		"plugins.dyn.unknown_leaf", "fs.something_else", "risk.new_key",
	}
	for _, k := range no {
		if IsSecurityLoosenKey(k) {
			t.Errorf("IsSecurityLoosenKey(%q) = true, want false", k)
		}
	}
	// The roster itself must carry the same names the width test pins, so the
	// constants cannot drift from the matcher.
	for _, want := range []string{"risk.shell_enabled", l1, tier2, "net.proxy.mode", "fs.delete_enabled"} {
		found := false
		for _, k := range SecurityLoosenRoster {
			if k.Key == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("the roster does not carry %q, but the matcher (and the auditors) do", want)
		}
	}
}
