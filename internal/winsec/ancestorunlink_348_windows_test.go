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
// 票 369 (leg 369-r1) REPLACED (b): it used to be "guardSpoke+recheckSpoke > 0 in 400 raced attempts",
// which is a count of how often the racer landed inside a two-syscall gap, and it read guard=0 here
// and re-check=0 there on unchanged bytes (0/400 and 1/400 on runner head 23b663c2). (b) is now the
// landing's own leg DRIVEN DETERMINISTICALLY on a bare-sibling leaf link - the one spelling containment
// cannot answer (342-v2's adjudication, restated in leafrace_342r2_windows_test.go:5-14) - so it is red
// when the re-check is broken and green when it is not, on every machine. The raced world stays, and
// keeps its own numerator: `deleted`, which is ticket 348's and is the half that must never move.
//
// Judged per ticket 348 AC#3: `deleted` is this ticket's numerator, gated here on both forms.
// `truncated` is ticket 342's and is gated in its own file on the junction form.
package winsec

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAC2WindowsPinRefusalOnAnAbsentAncestorIsNotFiledAsASealRefusal348r1 is the case for the one
// roster row this leg added, and it is deterministic rather than raced.
//
// WHY THE ROW NEEDS A CASE. A pin refusal and a seal refusal both carry ErrNotSealable, so the roster's
// coarse "cannot apply a private security descriptor" row filed BOTH under "a seal refusal on this
// route". That is not a cosmetic mislabel: this file's liveness gate counts attempts the ROUTE answered
// about a flipped ancestor, and an ancestor that is ABSENT is not a flipped ancestor - so a 400-attempt
// run that spent 395 attempts on a parked (absent) ancestor reported 395 "seal refusals", proved three
// flips, and passed. Measured before the row existed (348-r1, this machine).
//
// WHY THE WORLD IS BUILT THIS WAY. The ancestor is deleted outright and nothing is ever flipped: no
// reparse point is involved anywhere, so no reparse leg can answer, and the only possible answer is the
// pinned walk failing to pin `sub`. That makes the attribution question - not the damage question - the
// whole subject, and it makes the case independent of the racer's timing.
func TestAC2WindowsPinRefusalOnAnAbsentAncestorIsNotFiledAsASealRefusal348r1(t *testing.T) {
	base := SealableTempDirForTest124(t)
	data := filepath.Join(base, "data")
	sub := filepath.Join(data, "sub")
	name := filepath.Join(sub, winAncLeaf342r1)
	if err := os.MkdirAll(sub, 0o700); err != nil {
		t.Fatalf("stage the ancestor: %v", err)
	}
	if err := os.WriteFile(name, []byte(winAncOwnBody342r1), 0o600); err != nil {
		t.Fatal(err)
	}
	// Park the ancestor: the caller's leaf goes with it, exactly as a half-finished flip leaves it.
	if err := os.Remove(name); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(sub); err != nil {
		t.Fatal(err)
	}
	err := PrivateFile(name, []byte(winAncPayload342r1), 0o600)
	if err == nil {
		t.Fatalf("AC#2 RED instrument: a private write under an absent ancestor returned nil, so this "+
			"case never reached the walk it is about (name %s)", name)
	}
	leg := spokeWinAnc342r1(err)
	want := winAncPinWalkLeg348r1
	if leg != want {
		t.Fatalf("AC#2 RED: an ABSENT ancestor was filed as %q; want the pin-walk row %q. Both the pin "+
			"and the seal wrap ErrNotSealable, so without this split an absent world is counted as the "+
			"seal leg - and the seal leg is what the liveness gate reads as 'the route met a flipped "+
			"ancestor'. verbatim: %v", leg, want, err)
	}
	if winAncIsFormLeg342r1(leg) {
		t.Fatalf("AC#2 RED: an absent ancestor was counted as one of form 乙's reparse legs (%q), so "+
			"the liveness gate would read a world that never reached the link state as having reached it",
			leg)
	}
	t.Logf("AC#2' attribution: the absent ancestor is answered by the pin-walk row, not by the seal "+
		"row (verbatim: %v)", err)
}

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
	// (A) THE LANDING'S OWN LEG, DRIVEN DETERMINISTICALLY, and this is 票 369's replacement for the old
	// `guardSpoke+recheckSpoke > 0`. The spelling is the bare sibling (§342-v2: containment has nothing to
	// refuse, so the package's own identity gate is the only leg that CAN answer), the leaf link is
	// constructed and not raced, and the leg is called the way the route calls it. Red when the re-check
	// is broken, green when it is not, on every machine.
	dir := SealableTempDirForTest124(t)
	sub := filepath.Join(dir, "sub")
	if err := os.MkdirAll(sub, 0o700); err != nil {
		t.Fatal(err)
	}
	bystander := filepath.Join(sub, "bystander-348r1.bin")
	const bystanderBody = "an object the caller never named, standing in the same directory, longer on purpose"
	if err := os.WriteFile(bystander, []byte(bystanderBody), 0o600); err != nil {
		t.Fatal(err)
	}
	linked := filepath.Join(sub, winAncLeaf342r1)
	if err := os.Symlink(filepath.Base(bystander), linked); err != nil {
		t.Skipf("cannot construct the relative file symlink this leg is about: %v (needs developer mode or "+
			"SeCreateSymbolicLinkPrivilege; acceptance has to count this SKIP, never a quiet pass)", err)
	}
	if !reparseAt(linked) {
		t.Fatalf("%s is not a reparse point, so this case measured nothing", filepath.Base(linked))
	}
	pin, err := pinnedParent(linked)
	if err != nil {
		// A pin failure here would be the world, not the leg: the leaf form keeps `sub` real.
		t.Fatalf("AC#2 RED instrument: pinning the parent of a real directory failed: %v", err)
	}
	handle, _, leafErr := openPinnedLeaf(pin, filepath.Base(linked), os.O_WRONLY|os.O_CREATE|os.O_TRUNC|oNoFollowPlatform, 0o600)
	_ = pin.close()
	if handle != nil {
		_ = handle.Close()
	}
	const wantLinkFace = "the name carries a link, so this handle is not the object the caller named"
	if leafErr == nil {
		t.Fatalf("AC#2 RED: openPinnedLeaf returned a usable handle for a name that carries a link, so this "+
			"route's bytes land in an object the caller never named (name %s) - the re-check is the leg that "+
			"stops it and it did not answer", linked)
	}
	if !strings.Contains(leafErr.Error(), wantLinkFace) {
		t.Fatalf("AC#2 RED: the refusal did not carry the re-check's link face %q, so some other leg answered; "+
			"verbatim: %v", wantLinkFace, leafErr)
	}
	if !errors.Is(leafErr, ErrUnresolvedPath) || !errors.Is(leafErr, ErrNotSealable) {
		t.Fatalf("AC#2 RED: the re-check's refusal lost a sentinel (ErrUnresolvedPath/ErrNotSealable): %v", leafErr)
	}
	t.Logf("AC#2' the landing's own leg answered, verbatim: %v", leafErr)
	// The residual, measured rather than raced: the followed open carried out the truncation BEFORE the
	// re-check could refuse it. Reported and NOT gated - a future landing that closes it must make this
	// line read "intact", and gating "truncated" would forbid exactly that improvement.
	after, rErr := os.ReadFile(bystander)
	t.Logf("AC#2' RESIDUAL, deterministic on this spelling: the followed open left the bystander as %q "+
		"(read err %v) - the bytes were refused afterwards and the NAME was never removed, which is the half "+
		"this platform may claim", after, rErr)

	// (B) LIVENESS, deterministically: the flipped state is entered by construction and PROVED, instead of
	// being inferred from a racer's count. Keep that apart from (C): a liveness proof is about the world,
	// never about the leg.
	w := newWinAncWorld342r1(t, winAncFormLeafSym)
	if err := w.toFlip(); err != nil {
		t.Fatalf("AC#2 RED instrument: the world cannot enter the flipped leaf state at all: %v", err)
	}
	if !reparseAt(w.name) {
		t.Fatalf("AC#2 RED instrument: %s is not a reparse point after the synchronous flip, so the liveness "+
			"this case used to read off the racer was never proved", w.name)
	}
	t.Logf("AC#2' liveness: the flipped leaf state is entered synchronously and proved by reparseAt(%s)=true",
		filepath.Base(w.name))
	if err := w.toReal(); err != nil {
		t.Fatalf("settle the world back to the ordinary state: %v", err)
	}

	// (C) THE RACED WORLD, judged only on ticket 348's numerator.
	st := w.run(winAncAttempts342r1, func() error {
		return PrivateFile(w.name, []byte(winAncPayload342r1), 0o600)
	})
	st.print(t, "AC#2 (ticket 348 + the leaf residual, Windows leaf-symlink form)")
	if st.deleted > 0 {
		t.Errorf("AC#2 RED: %d of %d attempts found the foreign witness GONE after a leaf name carrying a "+
			"link, so the refusal-path removal still resolved into somebody else's tree (first: %s)",
			st.deleted, st.attempts, st.firstProof)
	}
	// The split is a READING about the clock: which gap a racer lands in is not a property of the landing,
	// and this file stopped gating it in 票 369 (see the file header). Printed so the residual story stays
	// readable, never demanded.
	guardSpoke := st.legs["preOpenLeafGuard's leaf refusal (ticket 322's Windows leaf leg)"]
	recheckSpoke := st.legs["leaf identity re-check inside openPinnedLeaf (os.Root pins ancestors, NOT the leaf)"] +
		st.legs["leaf identity re-check, SameFile face"]
	t.Logf("AC#2' leaf-leg split, REPORTED and not gated: preOpenLeafGuard (check before the open)=%d, "+
		"openPinnedLeaf re-check (after the open)=%d, provedReparse=%d - that split IS the residual story of "+
		"ticket 322's Windows leaf leg, and it is a reading about the racer's timing, not about the landing",
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
