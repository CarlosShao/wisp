//go:build windows

package tools

import (
	"testing"
)

// Ticket 252 R2 leg (252-r2), W1: the VALUE witness for the lexical leg.
//
// The ban it measures is AC#2's, verbatim: "绝不许把两次包含改成一次". This is
// the half of it that had no instrument at all - the acceptance table
// (docs/evidence/s1/252-allowlist-sameform-v1.md) deleted paths.go's first
// containment and this package answered 257 PASS / 0 FAIL, a reading this leg
// re-ran itself before writing anything (see .scratch/wisp/probes/252/r2/).
//
// Why a value test for that leg has to look strange: ticket 252's fix moved the
// same-form step into Canonicalize, so every string production hands
// InAllowlist is already one form on both sides and the two legs cannot
// disagree. The only inputs on which they can are the ones that skipped that
// step - so W1 calls InAllowlist directly, with no Canonicalize in front of it,
// on a spelling that has not been aligned. That is not a shape a shipped tool
// reaches today, and it must never become the reason the leg is deleted: the
// leg is what makes InAllowlist safe for callers that are not this file's
// Canonicalize (risk.PathCanonicalizer is an interface precisely so other
// feeds exist), and it is the reason "two containments" is a pair rather than
// one check written twice.
//
// The carrier reuses the 252-p1 meter (measure252 / legs252 / shortName252 /
// repoRoot252, same package, same build tag) and the same rule it states: short
// names come from the OS, never typed by hand, and where this machine gives no
// usable alias the case FAILS LOUDLY. Zero t.Skip calls in this file.
//
// Each cell also carries its own positive control - the same raw ask taken
// through the full production round trip must be authorized - so a future
// tightening can never let this ruler pass by refusing everything.

// TestTicket252R2LexicalLegRefusesWhatOnlyItRefuses is W1 proper: two inputs of
// shape LEG1=false / LEG2=true. Delete the lexical containment and both flip to
// authorized on the spot; delete the resolvedForm containment and this test
// stays green, which is the point - it is the mirror image of W2, and the pair
// is what "each leg reachable on its own" means as a reading.
func TestTicket252R2LexicalLegRefusesWhatOnlyItRefuses(t *testing.T) {
	long := repoRoot252(t)
	short := shortName252(t, long)
	const leaf = "q252r2-lexical-leg-only.md"

	cases := []struct {
		name  string
		asked string
		why   string
	}{
		// The ticket 252 shape itself: an 8.3 alias of an existing tree plus a
		// leaf that does not exist yet. C26's handle query answers for the tree
		// and not for the leaf, so this is the spelling that used to arrive.
		{
			name: "B1 short-of-existing-tree + missing-leaf", asked: short + pathSep + leaf,
			why: "the 8.3 spelling of the allowed root, new file under it",
		},
		// Same disagreement without the missing-leaf half: the whole ask exists,
		// so this cell shows the disagreement is about WHO expands the alias,
		// not about whether the path is there.
		{
			name: "B2 short-of-existing-tree itself", asked: short,
			why: "the allowed root spelled by its own alias",
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			pc := NewPathCanonicalizer([]string{long}, nil)
			if u := pc.UnusableRoots(); len(u) != 0 {
				t.Fatalf("carrier broken: the configured root was dropped: %q", u)
			}
			l := measure252(pc, c.asked)
			l.log(t, "  ")

			// The premise, asserted rather than assumed: this input must be one
			// the legs disagree on, or the case measures nothing.
			if l.leg1 {
				t.Fatalf("carrier moved: the lexical leg already answers true for %q (folded=%q, roots=%q), "+
					"so this cell no longer isolates that leg", c.asked, l.folded, l.roots)
			}
			if !l.leg2 {
				t.Fatalf("carrier moved: the resolvedForm leg answers false for %q (resolved=%q ok=%v), so the "+
					"refusal below cannot be attributed to the lexical leg", c.asked, l.rf, l.rfOK)
			}
			if !l.final {
				t.Errorf("TICKET 252 BAN RED (lexical containment missing / two containments merged into one): "+
					"InAllowlist(%q) = true (%s) although the lexical leg refuses it: foldPath(canonical)=%q sits "+
					"under no root in %q, and only resolvedForm=%q (ok=%v) would. AC#2's "+
					"\"绝不许把两次包含改成一次\" is what keeps that leg in the code; ticket 107 round 2's own note "+
					"(paths.go:199-207, \"Requiring both is strictly narrower\") is the reason it may not be called "+
					"redundant and removed.",
					c.asked, c.why, l.folded, l.roots, l.rf, l.rfOK)
			}

			// Positive control: the same raw ask through Canonicalize - the
			// production round trip - must be authorized. This is 252-r1's
			// verdict, re-measured here so W1 cannot be green merely because
			// the judge got stricter somewhere else.
			canonical, err := pc.Canonicalize(c.asked)
			if err != nil {
				t.Fatalf("carrier broken: Canonicalize(%q): %v", c.asked, err)
			}
			cl := measure252(pc, canonical)
			cl.log(t, "  aligned ")
			if !cl.final {
				t.Errorf("POSITIVE CONTROL RED for W1: %q through the production round trip is still refused "+
					"(canonical=%q leg1=%v leg2=%v roots=%q), so this ruler is being satisfied by a judge that "+
					"refuses both spellings and its negative reading means nothing",
					c.asked, canonical, cl.leg1, cl.leg2, cl.roots)
			}
			if cl.leg1 != cl.leg2 {
				t.Errorf("AC#1's split is live again on an aligned input: leg1=%v leg2=%v (canonical=%q)",
					cl.leg1, cl.leg2, canonical)
			}
		})
	}
}
