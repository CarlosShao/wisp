package config

import (
	"fmt"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Unwired security keys (ticket 83; ruling A53② = ticket 80's option (C)).
//
// A key is in this file when all three hold:
//
//  1. it sits in a locked (🔒) section and the direction auditor in manager.go
//     names it as a loosen/tighten key, i.e. the config layer already asserts
//     it has a security consequence;
//  2. nothing in this repository reads its value (ticket 80 AC#1's read/write
//     inventory, plus the indirect-consumption sweep recorded there and in
//     ticket 83 AC#1: no reflection-by-name lookup, no name-keyed map, no
//     TOML passthrough feeds any of these);
//  3. its non-default value asks for a *relaxation* of a gate that the running
//     code does not have.
//
// Such a key lies: the user writes one line and the program swallows it. The
// repo's existing answer to a lying key is the hard-coded read-only trio
// (`models.verify_signature`, `audio.half_duplex`, `privacy.keep_*`) -
// SPEC-03 sec 3 pins the shape ("写 false 报错而非生效"), validate.go implements
// it. These keys are the same shape with one difference: they are not
// hard-coded to a value, the capability behind them simply is not built yet.
// So the error says exactly that and names who builds it.
//
// Deliberately NOT here: zero-consumer keys of unlocked sections (ball, panel,
// cost, voice, memory, ...). Those are the ordinary not-built-yet debt of the
// subsystem that owns them (see ticket 83 AC#1's table); none of them is
// referenced by a security verdict, so writing them cannot read as "the
// security posture changed". The [plugins] keys are the one open question and
// are handed to the orchestrator rather than decided here.
//
// The guard fires only on a NON-DEFAULT value, so a file that never mentions
// these keys - and one that SaveFile writes back out while they sit at their
// defaults - keeps loading exactly as before.

// unwiredKey is one lying security key plus the sentence that explains it.
type unwiredKey struct {
	// path is the config.toml key path, spelled exactly as the direction
	// auditor spells it in manager.go, so the error and the audit log agree.
	path string
	// missing states which consumer is absent, by symbol, so the claim is
	// checkable rather than an assertion of faith.
	missing string
	// lands names who turns the key into a real switch.
	lands string
	// fires reports whether this config asks for the unimplemented capability.
	fires func(*Config) bool
}

// unwiredKeys is the guard table. Deleting a row is part of landing the
// capability it names: a load-time error must not outlive the consumer that
// justifies it.
var unwiredKeys = []unwiredKey{
	{
		path:    "risk.shell_enabled",
		missing: "no shell.exec tool is registered in internal/tools, so nothing reads this flag",
		lands:   "it lands with the shell.exec tool itself (SPEC-07 sec 3, S3, default-disabled per D14)",
		fires:   func(c *Config) bool { return c.Risk.ShellEnabled },
	},
	{
		path:    "risk.allow_shell_string",
		missing: "risk.Facts.ShellString has no production writer, so R6's string-mode verdict never runs",
		lands:   "it lands with the shell.exec tool (SPEC-06 sec 10 forced-argv rule)",
		fires:   func(c *Config) bool { return c.Risk.AllowShellString },
	},
	{
		path:    "risk.shell_allowlist",
		missing: "risk.Facts.ShellAllowlist has no production writer, so R6's allowlist verdict never runs",
		lands:   "it lands with the shell.exec tool (SPEC-06 sec 10)",
		fires:   func(c *Config) bool { return len(c.Risk.ShellAllowlist) > 0 },
	},
	{
		path:    "risk.blacklist_overrides",
		missing: "risk.Gate has zero production call sites (the live chain calls risk.Classify, whose signature takes no override set) and nothing populates bOverrides",
		lands:   "the contract defines a B-tier exemption as a runtime event - one L2 confirmation plus a log line (SPEC-06 sec 4.1), which lands with the approval queue, ticket 21; wiring a static pre-authorization instead would need D22 (ticket 80 AC#2)",
		fires:   func(c *Config) bool { return len(c.Risk.BlacklistOverrides) > 0 },
	},
	{
		path:    "net.allowlist",
		missing: "risk.NetTarget.Allowlist has no production writer, and no web/open tool exists to build one, so R5's domain gate never runs",
		lands:   "it lands with the web tool family and its D30 egress layers, ticket 22",
		fires:   func(c *Config) bool { return len(c.Net.Allowlist) > 0 },
	},
	{
		path:    "net.block_private_ranges",
		missing: "R5's private-range verdict is unconditional in internal/risk/rules_network.go, so there is no switch to turn off",
		lands:   "switching it off would mean a new relaxation seam in the frozen risk layer: SPEC-06 sec 6.3 pins the SSRF check, so this needs a contract decision (D22), not a consumer",
		fires:   func(c *Config) bool { return !c.Net.BlockPrivateRanges },
	},
}

// validateUnwired rejects a config that asks for a capability nothing reads.
func validateUnwired(c *Config) error {
	for _, k := range unwiredKeys {
		if !k.fires(c) {
			continue
		}
		return observe.New(observe.ClassConfig, fmt.Sprintf(
			"config.toml: %s is written but does nothing: %s. %s. "+
				"Remove the key or set it back to its default - accepting it silently would be a lying config option.",
			k.path, k.missing, k.lands))
	}
	return nil
}

// lockedKeyDisposition records, for the completeness pin in unwired_test.go,
// what makes each locked-section key honest. A value of "unwired:<path>" must
// match a row of unwiredKeys; anything else is a key that is allowed to load
// silently today, with the reason written next to it. Adding a locked key
// without adding it here fails TestEveryLockedSectionKeyIsAccountedFor, which
// is what keeps this list from going stale the way the four [risk] keys did.
var lockedKeyDisposition = map[string]string{
	// [risk]
	"risk.confirm_timeout_sec": "consumed: cmd/wisp/run.go builds the approval timeout from it",
	"risk.l1_window_sec":       "consumed: cmd/wisp/run.go builds the L1 auto-approve window from it",
	"risk.shell_enabled":       "unwired:risk.shell_enabled",
	"risk.allow_shell_string":  "unwired:risk.allow_shell_string",
	"risk.shell_allowlist":     "unwired:risk.shell_allowlist",
	"risk.blacklist_overrides": "unwired:risk.blacklist_overrides",
	// [fs]
	"fs.allowed_dirs":             "consumed: cmd/wisp/run.go feeds the C26 canonicalizer",
	"fs.reparse_point_exceptions": "consumed: cmd/wisp/run.go feeds the C26 canonicalizer",
	"fs.delete_enabled":           "consumed: cmd/wisp/run.go gates the delete-capable tools",
	// [net]
	"net.allowlist":            "unwired:net.allowlist",
	"net.block_private_ranges": "unwired:net.block_private_ranges",
	"net.proxy":                "table header: the leaves are net.proxy.mode / net.proxy.url",
	"net.proxy.mode":           "consumed: chainOptions -> llm transport",
	"net.proxy.url":            "consumed: chainOptions -> llm transport",
	// [plugins] - same defect class, deliberately NOT guarded by this ticket:
	// the whole plugin engine (ticket 50/51) is unbuilt, so there is no half
	// soldered wire here, just an absent subsystem. Handed to the orchestrator.
	"plugins.tier2_enabled":      "not built: plugin engine is ticket 50/51 (orchestrator call on guarding it)",
	"plugins.<id>":               "not built: one table per plugin id, ticket 50/51",
	"plugins.<id>.enabled":       "not built: plugin engine is ticket 50/51",
	"plugins.<id>.capabilities":  "not built: plugin engine is ticket 50/51",
	"plugins.<id>.net_allowlist": "not built: plugin engine is ticket 50/51",
	"plugins.<id>.host_api":      "not built: plugin engine is ticket 50/51",
}
