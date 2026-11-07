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
	ancBytesLanding342r2 = "landing"
	ancBytesOldPin342r2  = "old-pin-bytes"
)

func TestAC1POSIXPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1(t *testing.T) {
	w := newAncWorld342r1(t, ancTargetRel)
	st := w.run(ancAttempts342r1, func() error {
		return PrivateFile(w.name, []byte(ancPayload342r1), 0o600)
	})
	landed := st.pinnedParentLeg > 0
	bytes := ancBytesOldPin342r2
	if landed {
		bytes = ancBytesLanding342r2
	}
	st.print(t, "AC#1 (ticket 342, damaged-by-the-open gate, in-root ancestor link, bytes="+bytes+")")

	// ATTRIBUTION FIRST: a row that cannot say which route produced it is not a reading, and the
	// landing requirement below would silently accept an old-side row that names a leg its bytes do
	// not have. pinnedParentBeforeOpen is the landing's own numerator: the pre-landing route cannot
	// raise it, and the landing's liveness gate requires it, so a row carrying 0 there is not a
	// landing-side reading either way - on pre-landing bytes it is the pre-fix reading, and in this
	// gate that is red, because what this gate certifies is the landing.
	if !landed {
		t.Errorf("AC#1 RED (not a landing reading): pinnedParentBeforeOpen answered %d attempts, and that "+
			"numerator is the landing's own - the pre-landing route cannot raise it and the landing's "+
			"liveness gate requires it - so this row cannot certify the landing's closed window either way. "+
			"On pre-landing bytes these numbers are the pre-fix reading (342-v1's VPROD/VBASE measured "+
			"deleted 65..86/1500 there; leg 342-r3 re-took deleted=%d of %d on the by-name route), and "+
			"reading any zero among them as 'the window was already closed' is the exact mistake ticket "+
			"342 AC#1 forbids (provedLinkCycles=%d, attempts=%d, flipErrs=%d, first %q)",
			st.pinnedParentLeg, st.deleted, st.attempts, st.provedLinkCycles, st.attempts,
			st.flipErrs, st.firstFlip)
		return
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
	if st.provedLinkCycles == 0 || st.pinnedParentLeg == 0 {
		t.Errorf("AC#1 RED instrument: the racer proved %d link-state cycles and the parent pin refused %d "+
			"attempts BEFORE any open (walk-after-open answers: %d, attempts: %d, flipErrs=%d, first %q) - "+
			"so the zeroes above are not evidence about a window that was ever reached, and a refusal that "+
			"came after the open proves nothing about it",
			st.provedLinkCycles, st.pinnedParentLeg, st.sealSideLeg, st.attempts, st.flipErrs, st.firstFlip)
	}
}
