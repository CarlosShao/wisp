//go:build windows

// Ticket 342 AC#1 + AC#3, leg 342-r1: the ANCESTOR leg of the open inside privateFile, driven
// through the production entry point ON THIS PLATFORM. World, three forms and the liveness rule: see
// ancestorworld_342_348_windows_test.go. This file judges the damage the open does, on both ancestor
// spellings, because they are answered by two different mechanisms and a green on one is not a green
// on the other.
//
// FORM 1, the junction. `mklink /J` stores an ABSOLUTE substitute name, so the pinned chain meets a
// reparse point with OBJ_DONT_REPARSE set and the kernel answers inside the descent - act time, with
// no window left to go stale. This is the row 342-a1 measured for form 乙 (WYI damage=0, root-refuse
// 1312/1500) and the shape Windows users actually have.
//
// FORM 2, a directory symlink with a RELATIVE target. Containment does NOT refuse it (the target stays
// inside the volume root the chain descends from), so Go's doInRoot follows it and the descent would
// land on somebody else's directory. The answer here can only be the chain's own identity gate -
// sameObjectPinnedWindows, which Lstats the name THROUGH THE PARENT HANDLE and refuses a name carrying
// a reparse point. This is the form that separates "pinned the ancestor" from "trusted the ancestor",
// and it is the Windows face of the reason the package's POSIX pinPlacement verifies rather than
// trusts. Needs developer mode: an unavailable construction is a t.Skipf that acceptance counts, never
// a quiet pass.
//
// JUDGED on each form:
//  1. truncated == 0 - no foreign object lost its body to a route that resolved through a flipped
//     ancestor, and returnedNil is printed because on the pre-landing bytes a following open SUCCEEDS
//     (that row is evidence about the window, not a pass);
//  2. liveness in the same reading: provedReparse > 0 (reparseAt confirmed the ancestor name really was
//     a reparse point after each flip - the key ticket 342 AC#1 names for this platform) AND
//     formRefusals > 0 (the refusal came from one of the landing's own legs, attributed per attempt);
//  3. deleted is printed and NOT gated - that numerator is ticket 348's (AC#3: its pin may not reuse
//     this one's green).
package winsec

import (
	"testing"
)

func TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1(t *testing.T) {
	for _, tc := range []struct{ form, why string }{
		{winAncFormJunction, "the act-time face: an absolute reparse target refused inside the descent"},
		{winAncFormSymDir, "the identity face: a relative in-root target refused by the chain's own re-check"},
	} {
		t.Run(tc.form, func(t *testing.T) {
			w := newWinAncWorld342r1(t, tc.form)
			st := w.run(winAncAttempts342r1, func() error {
				return PrivateFile(w.name, []byte(winAncPayload342r1), 0o600)
			})
			st.print(t, "AC#1 (ticket 342, Windows ancestor form "+tc.form+": "+tc.why+")")
			if st.truncated > 0 {
				t.Errorf("AC#1 RED: %d of %d attempts found the foreign witness's body replaced, so the open "+
					"inside privateFile resolved through an ancestor swapped into a link and its truncate "+
					"landed in somebody else's tree (first: %s)",
					st.truncated, st.attempts, st.firstProof)
			}
			if st.provedReparse == 0 || st.formRefusals == 0 {
				t.Errorf("AC#1 RED instrument: reparseAt proved %d flipped cycles and the landing's own legs "+
					"answered %d attempts across %d (flipErrs=%d, first %q) - so the zeroes above are not "+
					"evidence about a window that was ever reached",
					st.provedReparse, st.formRefusals, st.attempts, st.flipErrs, st.firstFlip)
			}
		})
	}
}
