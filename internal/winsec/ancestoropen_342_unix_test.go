//go:build unix

// Ticket 342 AC#1 + AC#3, leg 342-r1: the ANCESTOR leg of the open inside privateFile, driven
// through the production entry point on this platform. World, racer and liveness proof: see
// ancestorworld_342_348_unix_test.go. This file judges ONE numerator: truncated.
//
// THE DAMAGE THIS GATE OWNS. An ancestor renamed into a link after resolveString and the placement
// floor had answered used to be followed by the route's open, and that open carries O_TRUNC: 342-a1
// measured the production route at damage=106 on POSIX, and 342-v1 re-measured it (VPROD,
// truncated>0 under its pace sweeps). Form 乙 closes it by pinning the parent before the act: the
// leaf is opened relative to a descriptor no rename of an ancestor can re-point, so a truncation can
// only ever land on a child of the object the caller's spelling named before the gap.
//
// THE FORM THIS RUNS. A RELATIVE link target (`../someone-elses-tree`). That is the only spelling
// where the answer can only be the pin: os.Root's containment rule refuses an ABSOLUTE target for
// escaping the "/" root, but resolves a target that stays inside it - following the link - so with a
// relative target the walk's own identity gate (sameObjectPinned) is the sole leg that can refuse,
// and the OLD code walks through with nothing to refuse. Reading this row on the old bytes is what
// gives the numerator its teeth: on the pre-landing tree the open, the seal and the write all land in
// the foreign witness.
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

func TestAC1POSIXPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1(t *testing.T) {
	w := newAncWorld342r1(t, ancTargetRel)
	st := w.run(ancAttempts342r1, func() error {
		return PrivateFile(w.name, []byte(ancPayload342r1), 0o600)
	})
	st.print(t, "AC#1 (ticket 342, damaged-by-the-open gate, in-root ancestor link)")
	if st.truncated > 0 || st.deleted > 0 {
		t.Errorf("AC#1 RED: %d attempts found the foreign witness's body replaced and %d found it gone, so this "+
			"route's act reached THROUGH an ancestor swapped into a link after resolveString had answered "+
			"(first: %s). Both numerators count here on purpose: on this route an O_TRUNC that lands is very "+
			"often followed by the refusal-path unlink taking the name away, so a gate reading only the "+
			"surviving body would go GREEN on a damaged tree - the shape ticket 342 AC#3 refuses. Ticket 348's "+
			"own gate is the removal alone and does not borrow this reading.",
			st.truncated, st.deleted, st.firstProof)
	}
	if st.provedLinkCycles == 0 || st.pinnedParentLeg == 0 {
		t.Errorf("AC#1 RED instrument: the racer proved %d link-state cycles and the parent pin refused %d "+
			"attempts BEFORE any open (walk-after-open answers: %d, attempts: %d, flipErrs=%d, first %q) - so "+
			"the zeroes above are not evidence about a window that was ever reached, and a refusal that came "+
			"after the open proves nothing about it",
			st.provedLinkCycles, st.pinnedParentLeg, st.sealSideLeg, st.attempts, st.flipErrs, st.firstFlip)
	}
}
