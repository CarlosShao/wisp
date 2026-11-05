package config

// risk.blacklist_overrides: the persisted B-tier exemption list (ticket 230
// AC#4, owner ruling 2026-11-05, ledger A1549).
//
// WHY THIS FILE REPLACES A GUARD ROW. Until this ticket the key was registered
// in unwired.go: the config layer asserted it had a security consequence, no
// code read it, and a non-empty value was refused with "this key is written but
// does nothing". The owner's ruling made the opposite call - the B tier is
// default-deny, and the operator may persist per-file exemptions for it, the
// same way they persist the other policy inputs in [risk]. The A tier is not
// part of that decision and never can be: SPEC-06 4.1 pins it as absolutely
// denied, non-overridable, "no grant, session grant or config relaxation can
// unlock it" (internal/risk/blacklist.go's own header).
//
// So the key is now a real setting, and this file is what the config layer owes
// it at load. Three duties, all fail-loud, because the failure mode this repo
// keeps paying for is a line that reads as a policy change and is not one:
//
//  1. CANONICALIZE. Each entry goes through the C26 pipeline (risk.Resolve) and
//     the Config carries the resulting canonical spelling. This is not
//     cosmetic: the gate compares override keys against the handle-resolved
//     comparison form of the path it is asked about, so an entry that keeps the
//     operator's raw spelling is an exemption that may silently never fire
//     (8.3 short names, junction targets, mixed separators).
//  2. REFUSE A TIER. An entry whose canonical spelling classifies ClassA is a
//     load error. Accepting it would recreate, one layer down, exactly the
//     lying option unwired.go existed to prevent: the file would carry a
//     standing "exemption" for a file no exemption can ever unlock.
//  3. REFUSE AN ENTRY THAT CANNOT FIRE. B membership is decided by the file
//     name (internal/risk's B rules), so an entry matching no B rule can never
//     become true later. A typo must not sit in config.toml reading as a
//     standing exemption.
//
// The expansion-before-classification order in (1)/(2) is load-bearing: without
// it "~/.ssh/id_rsa" would be classified on the literal string, match the
// B-tier "id_*" name rule, and be accepted - the exact fail-open the tier check
// exists to stop.
//
// STILL UNDONE (reported as the seam, not hand-waved here): nothing consumes
// the list yet. tools.Options.Confirmations is the documented seam for Gate's
// bOverrides map, and keying it needs risk's comparison fold, which is
// unexported in a frozen package. See docs/evidence/s1/230-ac4-config-half.md.

import (
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/risk"
)

// validateBlacklistOverrides canonicalizes c.Risk.BlacklistOverrides in place
// and rejects every entry an exemption could never honour. An empty (or absent)
// list is the default and stays loadable: the guard must not punish a file that
// never asked for an exemption.
func validateBlacklistOverrides(c *Config) error {
	if c == nil || len(c.Risk.BlacklistOverrides) == 0 {
		return nil
	}
	canonical := make([]string, 0, len(c.Risk.BlacklistOverrides))
	seen := make(map[string]bool, len(c.Risk.BlacklistOverrides))
	for i, entry := range c.Risk.BlacklistOverrides {
		key := fmt.Sprintf("risk.blacklist_overrides[%d]", i)
		if strings.TrimSpace(entry) == "" {
			return observe.New(observe.ClassConfig, fmt.Sprintf(
				"config.toml: %s is empty: an exemption must name the file it exempts, "+
					"and an empty entry would sit in the file exempting nothing", key))
		}
		// The C26 pipeline, with the same [fs] exceptions the tool path uses,
		// so "canonical here" and "canonical there" are one spelling.
		res, err := risk.Resolve(entry, c.FS.ReparsePointExceptions)
		if err != nil {
			return observe.New(observe.ClassConfig, fmt.Sprintf(
				"config.toml: %s %q cannot be canonicalized: %v. An exemption must be the real "+
					"path of the file it exempts (SPEC-06 4.1: the blacklist takes the C26 resolver's answer)", key, entry, err))
		}
		path := res.Canonical
		switch class := risk.Classify(path); class {
		case risk.ClassA:
			return observe.New(observe.ClassConfig, fmt.Sprintf(
				"config.toml: %s %q resolves to %q, which is an A-tier sensitive path: the A tier is "+
					"absolutely denied and is never overridable (SPEC-06 4.1), so no exemption - standing, "+
					"per-session or in this file - can unlock it. Remove this entry", key, entry, path))
		case risk.ClassB:
			// The only tier an exemption can honour.
		default:
			return observe.New(observe.ClassConfig, fmt.Sprintf(
				"config.toml: %s %q resolves to %q, which matches no B-tier rule: the tier is decided by "+
					"the file name (internal/risk's B rules), so this exemption could never fire. "+
					"Check the spelling or remove the entry", key, entry, path))
		}
		// Two spellings of one file are one exemption; keeping both would only
		// make the gate's set do the same work twice.
		if seen[path] {
			continue
		}
		seen[path] = true
		canonical = append(canonical, path)
	}
	c.Risk.BlacklistOverrides = canonical
	return nil
}
