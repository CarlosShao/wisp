package config

// Ticket 39 AC#2's named roster of security-loosening keys (write leg
// 2026-11-18, dispatch .scratch/wisp/probes/orch/20261118/txt/39-security-l2-dispatch.txt).
//
// WHAT THIS LIST IS. The one machine-readable answer to "which config keys can
// a write LOOSEN". Every entry is copied from the D36 direction auditors
// themselves - the four *Direction functions in manager.go are the only
// authority on what a loosening is (sectionwrite.go's lockedDirectionOf and
// the reload path's planLocked both dispatch to them) - so a key is on this
// roster because a named auditor line emits it, never because this file
// invented a category. The Source field carries that auditor's file:line, and
// Test39LoosenRosterSourcesReReadTheAuditors re-reads manager.go at test time
// and fails if the line the roster names no longer carries the key - the
// roster cannot rot away from its sources without a red test.
//
// The [plugins] section is dynamic (the section is a map keyed by plugin id),
// so its entries are FAMILIES, spelled exactly the way pluginsDirection
// spells them: the auditor emits "plugins.<id>" for a new entry and
// "plugins.<id>.<leaf>" for the four leaves it walks. IsSecurityLoosenKey
// matches the same families and nothing wider.
//
// WHAT THIS LIST IS FOR. The section-write gate (sectionwrite.go) consults it
// before queueing a loosening: a key the direction auditors call a loosening
// but this roster does not know is a registrar drift, and the gate refuses it
// fail-closed rather than queueing a card for something the roster cannot
// name. It is a completeness pin in the direction ticket 255's TierRegistry
// already established for tiers, pointed at the security half.
//
// WHAT THIS LIST IS NOT. It is not a second verdict: loosen-vs-tighten is
// still decided only by the auditors. It is not an allowlist for callers -
// knowing a key is loosen-capable grants nobody the right to loosen it; the
// native L2 card still stands between the attempt and the write.

import "strings"

// LoosenKey is one roster row: the dotted key path (or the plugins family
// pattern), the locked section it lives in, and the file:line of the direction
// auditor that emits it as a loosening.
type LoosenKey struct {
	// Key is the dotted key path exactly as the auditor emits it. For the
	// dynamic [plugins] families the id part is the literal `<id>` the
	// auditor's own spelling uses; IsSecurityLoosenKey does the matching.
	Key string
	// Section is the locked section (risk|fs|net|plugins).
	Section string
	// Source is the direction auditor line that emits this key as a
	// loosening, as file:line relative to this package.
	Source string
	// Pin is the literal that must appear on the Source line. For static keys
	// it is the key path as the auditor spells it in its append(); for the
	// dynamic families it is the auditor's own concatenation fragment, so the
	// pin pins the CODE that builds the name, not a copied string.
	Pin string
}

// SecurityLoosenRoster is the named roster. Static keys first, then the
// dynamic [plugins] families; order is the auditor's own order.
var SecurityLoosenRoster = []LoosenKey{
	// ---- [risk] - riskDirection (manager.go:465) ----
	{Key: "risk.shell_enabled", Section: "risk", Source: "manager.go:467", Pin: `"risk.shell_enabled"`},
	{Key: "risk.allow_shell_string", Section: "risk", Source: "manager.go:472", Pin: `"risk.allow_shell_string"`},
	{Key: "risk.l1_window_sec", Section: "risk", Source: "manager.go:477", Pin: `"risk.l1_window_sec"`},
	{Key: "risk.permission_mode", Section: "risk", Source: "manager.go:491", Pin: `"risk.permission_mode"`},
	{Key: "risk.shell_allowlist", Section: "risk", Source: "manager.go:495", Pin: `"risk.shell_allowlist"`},
	{Key: "risk.blacklist_overrides", Section: "risk", Source: "manager.go:496", Pin: `"risk.blacklist_overrides"`},
	// ---- [fs] - fsDirection (manager.go:501) ----
	{Key: "fs.allowed_dirs", Section: "fs", Source: "manager.go:502", Pin: `"fs.allowed_dirs"`},
	{Key: "fs.reparse_point_exceptions", Section: "fs", Source: "manager.go:503", Pin: `"fs.reparse_point_exceptions"`},
	{Key: "fs.delete_enabled", Section: "fs", Source: "manager.go:505", Pin: `"fs.delete_enabled"`},
	// ---- [net] - netDirection (manager.go:514) ----
	{Key: "net.allowlist", Section: "net", Source: "manager.go:515", Pin: `"net.allowlist"`},
	{Key: "net.block_private_ranges", Section: "net", Source: "manager.go:517", Pin: `"net.block_private_ranges"`},
	{Key: "net.proxy.mode", Section: "net", Source: "manager.go:525", Pin: `"net.proxy.mode"`},
	// ---- [plugins] - pluginsDirection (manager.go:534) ----
	{Key: "plugins.tier2_enabled", Section: "plugins", Source: "manager.go:536", Pin: `"plugins.tier2_enabled"`},
	// The dynamic families: the auditor builds these names per declared id,
	// so the pin is the concatenation fragment the auditor's own append
	// carries - pinning the code, not a copy of the result.
	{Key: "plugins.<id>", Section: "plugins", Source: "manager.go:542", Pin: `"plugins."+id`},
	{Key: "plugins.<id>.enabled", Section: "plugins", Source: "manager.go:547", Pin: `"plugins."+id+".enabled"`},
	{Key: "plugins.<id>.capabilities", Section: "plugins", Source: "manager.go:551", Pin: `"plugins."+id+".capabilities"`},
	{Key: "plugins.<id>.net_allowlist", Section: "plugins", Source: "manager.go:552", Pin: `"plugins."+id+".net_allowlist"`},
	{Key: "plugins.<id>.host_api", Section: "plugins", Source: "manager.go:553", Pin: `"plugins."+id+".host_api"`},
}

// staticLoosenKeys is the static half of the roster, for exact-match lookups.
var staticLoosenKeys = func() map[string]bool {
	m := make(map[string]bool, len(SecurityLoosenRoster))
	for _, k := range SecurityLoosenRoster {
		if !strings.Contains(k.Key, "<id>") {
			m[k.Key] = true
		}
	}
	return m
}()

// pluginsLoosenLeaves are the per-id leaves the [plugins] auditor walks, as
// suffixes under "plugins.<id>.".
var pluginsLoosenLeaves = []string{".enabled", ".capabilities", ".net_allowlist", ".host_api"}

// IsSecurityLoosenKey reports whether the dotted key path is one the D36
// direction auditors can emit as a LOOSENING: a static roster key, a new
// [plugins.<id>] entry, or one of the four per-id plugin leaves. Unknown paths
// answer false - the gate treats an unknown loosen key as registrar drift and
// refuses, never queues a card for a name the roster cannot carry.
func IsSecurityLoosenKey(path string) bool {
	if staticLoosenKeys[path] {
		return true
	}
	const prefix = "plugins."
	if !strings.HasPrefix(path, prefix) {
		return false
	}
	rest := path[len(prefix):]
	if rest == "" || strings.Contains(rest, "<") {
		return false // the family spellings themselves are not key paths
	}
	if !strings.Contains(rest, ".") {
		// "plugins.<id>": a NEW plugin entry is itself a loosening
		// (manager.go:531-535).
		return true
	}
	for _, leaf := range pluginsLoosenLeaves {
		if strings.HasSuffix(rest, leaf) {
			return true
		}
	}
	return false
}
