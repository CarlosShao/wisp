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
// it at load. Four duties, all fail-loud, because the failure mode this repo
// keeps paying for is a line that reads as a policy change and is not one:
//
//  1. CONSUME THE C26 REWRITE ACCOUNT. This leg does not merely classify the
//     resolver's answer - it PUBLISHES that answer as the key the gate will
//     match - so it has to ask what the pipeline substituted to get there
//     (Result.Actable, ticket 102). An entry whose %VAR%/$VAR/leading-~ was
//     expanded is a load error: keying the gate on the expanded tree would
//     exempt a file this line never named, while the line still reads as the
//     spelling the operator typed. Expansion stays in the pipeline on purpose
//     (SPEC-06 4 step 1); what is refused is letting its answer travel into the
//     exemption set unaccounted for.
//  2. CANONICALIZE. Each accepted entry goes through the same C26 pipeline
//     (risk.Resolve) and the Config carries the resulting canonical spelling.
//     This is not cosmetic: the gate compares override keys against the
//     handle-resolved comparison form of the path it is asked about, so an
//     entry that keeps the operator's raw spelling is an exemption that may
//     silently never fire (8.3 short names, junction targets, mixed
//     separators). Those are handle-resolution steps, not rewrites, so they do
//     not trip (1).
//  3. REFUSE A TIER. An entry whose expanded spelling classifies ClassA is a
//     load error. Accepting it would recreate, one layer down, exactly the
//     lying option unwired.go existed to prevent: the file would carry a
//     standing "exemption" for a file no exemption can ever unlock.
//  4. REFUSE AN ENTRY THAT CANNOT FIRE. B membership is decided by the file
//     name (internal/risk's B rules), so an entry matching no B rule can never
//     become true later. A typo must not sit in config.toml reading as a
//     standing exemption.
//
// The tier check still runs on the EXPANDED form, on purpose: "~/.ssh/id_rsa"
// must reach ClassA instead of matching the literal B-tier "id_*" name rule.
// The rewrite account binds where the answer becomes an exemption key, i.e. the
// B tier; an A-tier refusal keeps its tier reason and additionally reports the
// expansion (see the ClassA branch below).
//
// THE REMAINING WIRING LEG (reported as the seam, not hand-waved here): no
// committed code in cmd/wisp arms the list yet. The consumers exist -
// tools.Options.PersistedExemptions is the provider risk.GateStanding takes as
// its third argument (the second is the runtime L2 confirmations) - so arming
// it is a pure fold with risk.OverrideKey, not a new seam. See
// docs/evidence/s1/230-ac4-config-rewrite-account.md.

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
		// The tier is decided on the expanded form on purpose (SPEC-06 4:
		// expand, then classify): that is what makes "~/.ssh/id_rsa" reach
		// ClassA instead of matching the literal B-tier "id_*" name rule.
		switch class := risk.Classify(res.Canonical); class {
		case risk.ClassA:
			reason := fmt.Sprintf(
				"config.toml: %s %q resolves to %q, which is an A-tier sensitive path: the A tier is "+
					"absolutely denied and is never overridable (SPEC-06 4.1), so no exemption - standing, "+
					"per-session or in this file - can unlock it. Remove this entry", key, entry, res.Canonical)
			if res.Rewritten {
				reason += fmt.Sprintf(" (C26 also expanded this spelling: %s; the refusal stands on the "+
					"expanded form, so neither form may be exempted)", strings.Join(res.Rewrites, "+"))
			}
			return observe.New(observe.ClassConfig, reason)
		case risk.ClassB:
			// The only tier an exemption can honour - and the only branch whose
			// answer becomes a published key, so this is where the rewrite
			// account binds (ticket 102). An expanded spelling is not the file
			// the operator named, so it never reaches c.Risk.BlacklistOverrides.
			path, err := res.Actable()
			if err != nil {
				return observe.New(observe.ClassConfig, fmt.Sprintf(
					"config.toml: %s %q cannot be a standing exemption: %v. "+
						"A spelling C26 has to substitute into does not name the file it ends up at; "+
						"publishing the substituted tree as the exemption key would exempt a path this "+
						"entry never spelled. Write the entry as the literal path, or remove it", key, entry, err))
			}
			// Two spellings of one file are one exemption; keeping both would only
			// make the gate's set do the same work twice.
			if seen[path] {
				continue
			}
			seen[path] = true
			canonical = append(canonical, path)
		default:
			reason := fmt.Sprintf(
				"config.toml: %s %q resolves to %q, which matches no B-tier rule: the tier is decided by "+
					"the file name (internal/risk's B rules), so this exemption could never fire. "+
					"Check the spelling or remove the entry", key, entry, res.Canonical)
			if res.Rewritten {
				reason += fmt.Sprintf(" (C26 also expanded this spelling: %s; an expanded spelling is not "+
					"the path this entry names)", strings.Join(res.Rewrites, "+"))
			}
			return observe.New(observe.ClassConfig, reason)
		}
	}
	c.Risk.BlacklistOverrides = canonical
	return nil
}
