//go:build windows

// Ticket 348 AC#2 + AC#3, leg 342-r1: the DELETION leg on Windows, plus the leg this platform had NO
// answer for before the landing. Two files' worth of judgement in one place, because ticket 348's own
// AC#1 named both: `os.Root`'s chain is "act-time refusal, but that family does not exist here today",
// and the leaf half of it does not exist at all - `os.Root` pins ancestors and follows an in-root
// relative link AT the leaf (342-a1's W-ROOTSYM row, followed=true, three runs).
//
// CASE 1 - the removal, on the in-root form. The ancestor flips to a RELATIVE directory symlink, which
// containment does NOT refuse, so the chain has to reach its identity gate to answer at all. When the
// route then refuses, its unlink is relative to the pinned handle, so the foreign name cannot be
// reached by it. On the pre-landing bytes that same spelling is followed and the refusal-path removal
// goes by name - which is ticket 348's damage, on the platform that had no pinning family at all.
// Needs developer mode: unavailable construction is a t.Skipf, never a quiet pass.
//
// CASE 2 - the leaf, stated as a measurement rather than as a claim. The ancestor stays real; only the
// LEAF name flips between a file and a relative symlink at the foreign object. `openPinnedLeaf`
// re-checks the name against the handle it got, so the bytes are refused and the removal is
// fd-relative. What this platform cannot promise, and what this case therefore reports instead of
// gating: if the followed open truncated the foreign object before the re-check ran, a refusal after
// the fact does not restore those bytes. So the gate here is (a) the foreign object is never REMOVED,
// and (b) the landing's leaf leg actually spoke; the truncation count is printed for the record, and
// winsec.go's statement keeps that residual named rather than writing "the window is closed".
//
// Judged per ticket 348 AC#3: `deleted` is this ticket's numerator, gated here on both forms.
// `truncated` is ticket 342's and is gated in its own file on the junction form.
package winsec

import (
	"testing"
)

func TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1(t *testing.T) {
	w := newWinAncWorld342r1(t, winAncFormSymDir)
	st := w.run(winAncAttempts342r1, func() error {
		return PrivateFile(w.name, []byte(winAncPayload342r1), 0o600)
	})
	st.print(t, "AC#2 (ticket 348, Windows in-root ancestor form: the removal face)")
	if st.deleted > 0 {
		t.Errorf("AC#2 RED: %d of %d attempts found the foreign witness GONE, so the refusal-path unlink in "+
			"privateFile re-resolved the caller's spelling through a swapped ancestor and removed a name in "+
			"somebody else's tree (first: %s)", st.deleted, st.attempts, st.firstProof)
	}
	if st.provedReparse == 0 || st.formRefusals == 0 {
		t.Errorf("AC#2 RED instrument: reparseAt proved %d flipped cycles and the landing's own legs answered "+
			"%d attempts across %d (flipErrs=%d, first %q) - so the zero above is not evidence about a removal "+
			"that was ever reachable", st.provedReparse, st.formRefusals, st.attempts, st.flipErrs, st.firstFlip)
	}
}

func TestAC2WindowsLeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1(t *testing.T) {
	w := newWinAncWorld342r1(t, winAncFormLeafSym)
	st := w.run(winAncAttempts342r1, func() error {
		return PrivateFile(w.name, []byte(winAncPayload342r1), 0o600)
	})
	st.print(t, "AC#2 (ticket 348 + the leaf residual, Windows leaf-symlink form)")
	// (a) the removal may not reach the foreign name, on this platform either.
	if st.deleted > 0 {
		t.Errorf("AC#2 RED: %d of %d attempts found the foreign witness GONE after a leaf name carrying a "+
			"link, so the refusal-path removal still resolved into somebody else's tree (first: %s)",
			st.deleted, st.attempts, st.firstProof)
	}
	// (b) SOME leaf leg has to have answered, or this reading says nothing about the leaf. Two legs can
	// answer a flipped leaf on this platform and they are NOT interchangeable: preOpenLeafGuard probes
	// the NAME before the open (ticket 322's Windows leg, and it answers first in the steady state),
	// while openPinnedLeaf's re-check answers only in the gap between that verdict and the open - which
	// is the residual 322-r1 named and the very reason the re-check exists. Both counts are printed and
	// the gate requires the pair to be non-zero; the landing's own leg is reported separately rather than
	// demanded, because demanding it would make this instrument a claim about how often a racer lands
	// inside a two-syscall window instead of a claim about the removal.
	guardSpoke := st.legs["preOpenLeafGuard's leaf refusal (ticket 322's Windows leaf leg)"]
	recheckSpoke := st.legs["leaf identity re-check inside openPinnedLeaf (os.Root pins ancestors, NOT the leaf)"] +
		st.legs["leaf identity re-check, SameFile face"]
	if guardSpoke+recheckSpoke == 0 || st.provedReparse == 0 {
		t.Errorf("AC#2 RED instrument: no leaf leg answered (guard=%d, landing's re-check=%d) while reparseAt "+
			"proved %d flipped leaf cycles across %d (flipErrs=%d, first %q) - so nothing here shows a leaf "+
			"leg is on the route rather than on the shelf",
			guardSpoke, recheckSpoke, st.provedReparse, st.attempts, st.flipErrs, st.firstFlip)
	}
	t.Logf("AC#2' leaf-leg split: preOpenLeafGuard (check before the open)=%d, openPinnedLeaf re-check (after "+
		"the open)=%d, provedReparse=%d - that split IS the residual story of ticket 322's Windows leaf leg",
		guardSpoke, recheckSpoke, st.provedReparse)
	// (c) the residual, REPORTED and not gated: a truncation the followed open already carried out is
	// not undone by a refusal after it. Naming the number is what keeps winsec.go's statement honest,
	// and it is why this platform's window may not be written as closed.
	if st.truncated > 0 {
		t.Logf("AC#2' RESIDUAL OBSERVED, counted and NOT claimed closed: %d of %d attempts left the foreign "+
			"witness with its body replaced, so that open still followed the flipped leaf name and truncated "+
			"it before the re-check ran (first: %s). The bytes were refused and the name was not removed; the "+
			"truncation itself is this platform's named residual (ticket 322's Windows leaf leg + ticket 342).",
			st.truncated, st.attempts, st.firstProof)
	} else {
		t.Logf("AC#2' residual reading on this run: truncated=0 across %d attempts with %d proved flipped-leaf "+
			"cycles - a zero on one machine is a reading, not a promise; it is NOT written as the window closing",
			st.attempts, st.provedReparse)
	}
}
