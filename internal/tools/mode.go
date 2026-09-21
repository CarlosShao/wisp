package tools

import (
	"github.com/CarlosShao/wisp/internal/risk"
)

// Permission mode plumbing (ticket 90, ruling R20).
//
// AC#1 is the reason this file exists as a seam rather than a variable:
// "模式作为显式参数进入决策链，不许做成包级全局变量 —— 全局变量等于任何代码路径
// 都能改权限". So:
//
//   - the value is read once per tool call from an injected ModeSource
//     (nil = the strictest mode, so an unwired composition cannot accidentally
//     run relaxed);
//   - it is then carried as a PARAMETER into route(), which is the one place a
//     level is turned into "ask the user" or "do not ask";
//   - nothing in this package stores a mode it could be talked into changing.
//
// The setter side (a user switching档) is not here on purpose: internal/perm
// owns it, because a switch needs a confirmation and an audit record, and the
// bridge is the component that must not be able to grant itself either.

// ModeSource is the read side of the permission mode, satisfied directly by
// *config.Config (so the bridge never imports the config package) and by
// internal/perm's Store.
type ModeSource interface {
	// PermissionMode returns the mode in effect. It must not block, and it must
	// not be able to change anything.
	PermissionMode() risk.Mode
}

// permissionMode resolves the mode for one call. Fail-closed by construction:
// no source, or a source that panics, means the strictest档.
func (b *Bridge) permissionMode() (m risk.Mode) {
	if b.modes == nil {
		return risk.DefaultMode()
	}
	defer func() {
		if rec := recover(); rec != nil {
			b.log("tools: MODE-READ-FAILED recovering=%v; falling back to %s",
				rec, risk.DefaultMode())
			m = risk.DefaultMode()
		}
	}()
	return b.modes.PermissionMode()
}

// BlacklistNote is risk.Gate's own reading of one call's paths, kept beside the
// assessor's verdict rather than merged into it.
//
// Why a second opinion at the routing layer, when the assessor already ran R3:
// AC#5 asked for risk.Gate to stop being a function with zero production call
// sites. It is the only code in this repository that knows the B-tier
// single-file-override rule (SPEC-06 §4.1), and the confirmation card needs
// exactly that fact - "this file can be unlocked by one L2 confirm" - which the
// frozen SensitiveClassifier interface cannot express (it returns a Tier and
// nothing else). Consulting it here also means a mode can never be the only
// witness: the A/B evidence reaches the card and the audit line on its own.
//
// It is NOT an authority. Reading this note grants nothing: an
// AlreadyUnlocked entry requires a confirmation record that only the approval
// queue (ticket 21) can produce, and until something populates
// Options.Confirmations that map stays empty, so Gate answers every B-tier path
// with "ask first" - the fail-closed side of the pair.
type BlacklistNote struct {
	// Class is the strictest tier seen across the call's paths.
	Class risk.Class
	// Absolute names the A-tier paths: non-overridable, no mode, no grant.
	Absolute []string
	// Unlockable names the B-tier paths that ONE L2 confirmation can unlock.
	Unlockable []string
	// AlreadyUnlocked names the B-tier paths whose single-file confirmation is
	// already on record.
	AlreadyUnlocked []string
	// Reason is the rule text Gate produced for the strictest finding.
	Reason string
	// Unresolved counts the paths this call could not canonicalize. They are
	// reported, never skipped: Gate only answers for a path C26 could resolve.
	Unresolved int
}

// readBlacklist runs the frozen A/B gate over the call's raw paths, each
// through the same C26 canonicalizer the judge used (SPEC-06 §4: whitelist and
// blacklist membership both go through the same resolver).
func (b *Bridge) readBlacklist(raw []string) BlacklistNote {
	var n BlacklistNote
	if b.paths == nil || len(raw) == 0 {
		return n
	}
	overrides := b.confirmations()
	for _, p := range raw {
		c, err := b.paths.Canonicalize(p)
		if err != nil {
			n.Unresolved++
			continue
		}
		d := risk.Gate(c, overrides)
		switch d.Class {
		case risk.ClassA:
			n.Absolute = append(n.Absolute, c)
		case risk.ClassB:
			if d.Allow {
				n.AlreadyUnlocked = append(n.AlreadyUnlocked, c)
			} else {
				n.Unlockable = append(n.Unlockable, c)
			}
		}
		if d.Class > n.Class {
			n.Class = d.Class
			n.Reason = d.Reason
		}
	}
	return n
}

// blacklistedSensitive reports whether the note found anything worth putting on
// a card: A-tier, B-tier, or an unresolved path on a call the judge already
// flagged as sensitive.
func (n BlacklistNote) blacklistedSensitive() bool {
	return n.Class != risk.ClassNone
}
