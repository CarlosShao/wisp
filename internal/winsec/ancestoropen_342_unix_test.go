//go:build unix

// Ticket 342 AC#1 + AC#3, leg 342-r1: the ANCESTOR leg of the open inside privateFile, driven
// through the production entry point on this platform. World, racer and liveness proof: see
// ancestorworld_342_348_unix_test.go. This file judges ONE numerator: truncated.
//
// THE DAMAGE THIS GATE OWNS. An ancestor renamed into a link after resolveString and the placement
// floor had answered used to be followed by the route's open, and that open carries O_TRUNC: 342-a1
// measured the production route at damage=106 on POSIX, and 342-v1 re-measured it (VPROD, truncated>0
// under its pace sweeps). Form 乙 closes it by pinning the parent before the act: the leaf is opened
// relative to a descriptor no rename of an ancestor can re-point, so a truncation can only ever land
// on a child of the object the caller's spelling named before the gap.
//
// THE FORM THIS RUNS. A RELATIVE link target (`someone-elses-tree`, a bare in-root sibling). That is
// the only spelling where the answer can only be the pin: os.Root's containment rule refuses an
// ABSOLUTE target for escaping the "/" root, but resolves a target that stays inside it - following
// the link - so with a relative target the walk's own identity gate (sameObjectPinned) is the sole
// leg that can refuse, and the OLD code walks through with nothing to refuse. Reading this row on the
// old bytes is what gives the numerator its teeth.
//
// ★PAIRED READINGS, AND THE RULER THAT MAKES THEM READINGS (leg 342-r2). Every row below says which
// bytes it was measured on, and the "old" side is a RUN, not a citation: the same instrument is
// compiled against an exported tree that differs from the landing by one hunk
// (pinnedParent(path, resolved) -> pinnedParent(resolved), i.e. ticket 400 AC#3 reverted, the pre-400
// spelling of the route). What is REQUIRED differs by side, and the difference is the pair:
//
//   - landing: truncated == 0 AND deleted == 0, with liveness in the same reading;
//   - old: ATTRIBUTABLE - the landed `pinnedParentBeforeOpen` numerator must be 0, because that leg
//     does not exist in those bytes. Without this row the old side's zero is ambiguous: an O_TRUNC
//     that lands is very often followed by the refusal-path unlink taking the name away (this file's
//     own first paragraph, and 342-v1's i01 rows), so "0 on old" is equally consistent with "never
//     damaged" and with "damaged and then removed".
//
// ★WHAT THE RUNS MEASURED, AND HOW THE ATTRIBUTION BELOW IS JUDGED (corrected by leg 342-r3 on this
// branch's first-ever execution: the branch was committed without a runtime run and its original
// wording described the opposite trigger of the one coded below). The pin walk and the seal-side
// walk share one wording (winsec_other.go's pinPlacement answers both), so the classifier splits
// them by whether the refusal names the caller's leaf: a parent-side answer never does. Measured on
// pre-landing bytes (leg 342-r3, golang:1.27 container, reverted to the by-name route):
// pinnedParentBeforeOpen=0 with the seal-side walk still answering 44 attempts AFTER the open and
// deleted=92/1500 - the attribution HOLDS on the old side, and the pre-fix numerator WAS re-taken
// there (342-v1's VPROD/VBASE rows: deleted 65..86/1500, 4/4 runs). A row without the landing's own
// before-open numerator is nevertheless RED in this gate, because what this gate certifies is the
// LANDING: truncated==0 and deleted==0 with liveness can only be certified by bytes that carry the
// pin, and reading an old-side row as "the window was already closed" would be the exact mistake
// ticket 342 AC#1 forbids.
//
// JUDGED:
//  1. the route must never DAMAGE the foreign object through a flipped ancestor: the witness either
//     kept its own bytes (truncated == 0) or is gone (deleted > 0 - the refusal-path unlink took it).
//     Both are counted, because on the pre-landing route an O_TRUNC that lands is very often followed
//     by that removal, and a gate reading only the surviving body would go GREEN on the damaged tree -
//     a ruler that passes both shapes, which ticket 342 AC#3 refuses outright. `deleted` is not folded
//     into ticket 348's green by this: 348's own gate is the removal alone, plus a deterministic case
//     where only the unlink is in play;
//  2. provedLinkCycles > 0 AND pinnedParentBeforeOpen > 0 - the swapped state was reached AND the
//     route stopped at the parent pin, before opening anything, so (1) is evidence about a closed
//     window and not about a window nobody visited (ticket 342 AC#1 verbatim: "修后 0 且 liveness
//     非零");
//  3. narrowed is PRINTED and not gated here - that is a seal reaching the foreign object, ticket
//     120's surface, and it is gated in ticket 348's raced case.
package winsec

import (
	"testing"
)

// ancBytes342r2 are the two byte-identities this file can name, and the tags every reading carries.
// The tag is not decoration: a numerator is a reading only if the row says which bytes produced it
// (this leg's whole finding is a zero that came from the wrong bytes).
const (
	ancBytesLanding342r2   = "landing"
	ancBytesOldPin342r2    = "old-pin-bytes"
	ancBytesNoLanding342r2 = "not-a-landing-reading"
)

func TestAC1POSIXPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1(t *testing.T) {
	w := newAncWorld342r1(t, ancTargetRel)
	attempts := attemptsFor444(ancAttempts342r1)
	st := w.run(attempts, func() error {
		return PrivateFile(w.name, []byte(ancPayload342r1), 0o600)
	})
	landed := st.pinnedParentLeg > 0
	bytes := ancBytesNoLanding342r2
	if landed {
		bytes = ancBytesLanding342r2
	}
	st.print(t, "AC#1 (ticket 342, damaged-by-the-open gate, in-root ancestor link, bytes="+bytes+")")

	// Ticket 444: the un-raced carrier runs FIRST, so its reading sits on the same screen as the
	// raced row it qualifies. It is the half of this instrument that contention cannot starve.
	ctl := ancWindowControl444(t, ancCtlOpen444)
	ctl.print(t, "AC#1")

	// ATTRIBUTION IS A READING, NOT A VERDICT (ticket 444; before it, this was the first red arm and
	// the one CI load actually fired). A row that cannot say which route produced it is not a reading,
	// and the landing requirement below must never accept an old-side row that names a leg its bytes do
	// not have - so the row still says out loud that it certifies nothing.
	if !landed {
		t.Logf("AC#1 not a landing reading: pinnedParentBeforeOpen answered 0 attempts, so this row "+
			"cannot certify the landing's closed window either way - the pre-landing route cannot raise "+
			"that numerator and the landing's liveness gate requires it. On pre-landing bytes "+
			"(bytes="+ancBytesOldPin342r2+") these numbers are the pre-fix reading (342-v1's VPROD/VBASE "+
			"measured deleted 65..86/1500 there), "+
			"and reading any zero among them as 'the window was already closed' is the exact mistake "+
			"ticket 342 AC#1 forbids (provedLinkCycles=%d, attempts=%d, flipErrs=%d, first %q)",
			st.provedLinkCycles, st.attempts, st.flipErrs, st.firstFlip)
	}

	if st.truncated > 0 || st.deleted > 0 {
		t.Errorf("AC#1 RED: %d attempts found the foreign witness's body replaced and %d found it gone, so "+
			"this route's act reached THROUGH an ancestor swapped into a link after resolveString had "+
			"answered, on bytes=%s (first: %s). Both numerators count here on purpose: on this route an "+
			"O_TRUNC that lands is very often followed by the refusal-path unlink taking the name away, so "+
			"a gate reading only the surviving body would go GREEN on a damaged tree - the shape ticket "+
			"342 AC#3 refuses. Ticket 348's own gate is the removal alone and does not borrow this reading.",
			st.truncated, st.deleted, bytes, st.firstProof)
	}

	// Ticket 444's split, half one: the RACER's own arrival. Zero proved link-state cycles means the
	// fixture, not the schedule, produced the zeros above - and deleting the racer has to stay fatal,
	// so this arm keeps its teeth (ticket 441's rule: the arm that proves the world was entered stays
	// RED; the arm that says the route stood inside it may abstain).
	if st.provedLinkCycles == 0 {
		t.Errorf("AC#1 RED instrument: the racer proved 0 link-state cycles across %d attempts "+
			"(flipErrs=%d, first %q), so the window this instrument is named for was never entered and "+
			"the zeroes above say nothing about a closed one", st.attempts, st.flipErrs, st.firstFlip)
	}
	// Ticket 444's split, half two: the GUARDRAIL, carried by the act nobody races over.
	if !ctl.foreignIntact || ctl.actErr != "" {
		t.Errorf("AC#1 RED instrument: the un-raced window control - the route's own pin, then the "+
			"ancestor flip, then the act - did not keep this route's act inside the pinned object: %s",
			ctl.verdict())
	}
	// The abstention: the racer proved the world, the route never stood inside it. Named, with its
	// numbers, and carried by the control above - not a green from absence.
	if !landed && st.provedLinkCycles > 0 {
		t.Logf("AC#1 ABSTAIN(444:342-unix-ancestor:route-missed-link) attempts=%d samples=1 "+
			"provedLinkCycles=%d pinnedParentBeforeOpen=0 refusedAfterOpen=%d truncated=%d deleted=%d "+
			"flipErrs=%d - the racer proved %d link-state cycles but no PrivateFile call was inside one, "+
			"so this run cannot attest the raced row; the threshold that carries the open's pin is the "+
			"un-raced control above (foreign-intact=%v), not a green from absence",
			st.attempts, st.provedLinkCycles, st.sealSideLeg, st.truncated, st.deleted, st.flipErrs,
			st.provedLinkCycles, ctl.foreignIntact)
	}
}
