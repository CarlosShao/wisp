//go:build !windows

// Ticket 120 residual gap 1, the four cases that call sealHandle directly. The world
// each one builds, why the two faces need two different constructions, and what the
// read-backs are for are stated once in sealhandle_window_120_other_test.go; this
// file only asks the questions.
//
// Cases 1 and 2 are the two control targets: neuter the link face at
// winsec_other.go:269 alone and only case 2 (plus its twin on the PrivateFile route)
// changes colour; neuter the SameFile face at :272 alone and only case 1 (plus its
// route twin) does. Case 3 is the gap's cause, mechanically: the same link planted
// BEFORE the call is answered by the walk, never by :269, so it must stay green under
// both controls. Case 4 is the direction a leg that refused everything would fake:
// the honest path still has to seal.
package winsec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAC1POSIXSealHandleNameCarriesAnotherRealObject120r3 is the SameFile face: the
// caller's fd stands on object A, a second real object is renamed onto A's spelling
// before the call, so the pinned walk finds exactly what the name carries and is
// satisfied, and the only thing left that can notice the write would go somewhere
// else is sealHandle's own comparison. No link exists at any instant in this world.
func TestAC1POSIXSealHandleNameCarriesAnotherRealObject120r3(t *testing.T) {
	dir, name, _ := sealHandleFixture120r3(t)
	f := openHandleOn120r3(t, name)

	swapIn := filepath.Join(dir, "swap-in-120r3.bin")
	writeWide120r3(t, swapIn, "the object swapped into the name")
	if err := os.Rename(swapIn, name); err != nil {
		t.Fatalf("swap a second real object onto the name: %v", err)
	}
	// Fixture liveness, before the call, so a world that was never built cannot be
	// read as a leg that refuses everything.
	world := nameCarriesFdObject120r3(f, name)
	if world != "two different objects" {
		t.Fatalf("fixture did not build the world: the name and the fd hold the %s", world)
	}
	t.Logf("AC#1' world: fd and name = %s; name kind=%s mode=%s",
		world, kindOfName120r3(name), modeOfName120r3(name))

	err := sealHandle(f)
	assertLegSpoke120r3(t, err, whyHandleSwap120r3, "name-carries-another-real-object")

	// Both objects, read back: the walk narrowed the object the NAME denotes, which
	// is not the object the pending write stands on. That is why refusing is the only
	// correct answer, and it is the reading a neutered :272 changes.
	t.Logf("AC#1' read-back: name's object=%s, caller's object=%s",
		modeOfName120r3(name), modeOfFd120r3(f))
	if got := modeOfName120r3(name); got != "-rw-------" {
		t.Errorf("AC#1 RED: the pinned walk left the object the name carries at %s, want -rw------- "+
			"(so this run never reached the identity re-check at all)", got)
	}
	if got := modeOfFd120r3(f); got != "-rw-rw-rw-" {
		t.Errorf("AC#1 RED: the caller's own object was narrowed to %s while the refusal says the bytes "+
			"would not land in what was just sealed - the face at winsec_other.go:272 is not speaking", got)
	}
}

// TestAC1POSIXSealHandleNameBecomesLinkAfterPinnedSeal120r3 is the link face, and the
// only construction that reaches it: the pinned seal runs for real on the caller's
// plain file, and only then is the name taken away and replaced by a RELATIVE SIBLING
// link at the witness (the one link spelling the walk itself would follow, which is
// why the walk never sees it here).
func TestAC1POSIXSealHandleNameBecomesLinkAfterPinnedSeal120r3(t *testing.T) {
	_, name, witness := sealHandleFixture120r3(t)
	f := openHandleOn120r3(t, name)

	inIdentityWindow120r3(t, func() {
		// The attacker's move, landed in the window instead of hoped for.
		if err := os.Remove(name); err != nil {
			t.Errorf("fixture: take the name away: %v", err)
			return
		}
		if err := os.Symlink(filepath.Base(witness), name); err != nil {
			t.Errorf("fixture: plant the sibling link: %v", err)
		}
	})

	err := sealHandle(f)
	assertLegSpoke120r3(t, err, whyHandleLink120r3, "name-becomes-a-sibling-link-after-the-pinned-seal")

	// Liveness of BOTH halves: the plant really stands at the name, and the pinned
	// seal really ran for the caller's object (it is the object that came back
	// narrowed) - so the seam was windowed, not substituted.
	t.Logf("AC#1' read-back: name kind=%s, caller's object=%s, witness=%s",
		kindOfName120r3(name), modeOfFd120r3(f), modeOfName120r3(witness))
	if got := kindOfName120r3(name); got != "link" {
		t.Errorf("AC#1 RED: the name carries a %s, not the planted link, so this run did not build "+
			"the world for :269", got)
	}
	if got := modeOfFd120r3(f); got != "-rw-------" {
		t.Errorf("AC#1 RED: the caller's object is still %s, so applyDescriptor did not run for real "+
			"before the plant - the seam was substituted, not windowed", got)
	}
	if got := modeOfName120r3(witness); got != "-rw-rw-rw-" {
		t.Errorf("AC#1 RED: the witness the link points at was narrowed to %s, so a seal followed it", got)
	}
}

// TestAC1POSIXSealHandleNotReachedByAPrePlantedLink120r3 reproduces the gap's cause
// instead of asserting it: the same link, planted before the call the way every
// existing PrivateFile case does it, is answered by the pinned WALK, and the reading
// names which leg spoke. This is the case that shows :269 had no world of its own -
// and it is deliberately NOT a control target, because it never reaches that line.
func TestAC1POSIXSealHandleNotReachedByAPrePlantedLink120r3(t *testing.T) {
	_, name, witness := sealHandleFixture120r3(t)
	if err := os.Remove(name); err != nil {
		t.Fatalf("fixture: %v", err)
	}
	if err := os.Symlink(filepath.Base(witness), name); err != nil {
		t.Fatalf("fixture: plant the link before the call: %v", err)
	}
	f := openHandleOn120r3(t, name)
	t.Logf("AC#1' world: the name carries a %s at a relative sibling; the fd and the name hold the %s",
		kindOfName120r3(name), nameCarriesFdObject120r3(f, name))

	err := sealHandle(f)
	assertLegSpoke120r3(t, err, whyWalkLink120r3, "pre-planted-link-is-answered-by-the-walk")
	if err == nil || strings.Contains(err.Error(), whyHandleLink120r3) {
		t.Errorf("AC#1 RED: this reading names sealHandle's own link face, so it is not the gap's cause "+
			"(leg=%s)", spokeLeg120r3(err))
	}
	if got := modeOfName120r3(witness); got != "-rw-rw-rw-" {
		t.Errorf("AC#1 RED: the witness was narrowed to %s through the pre-planted link", got)
	}
}

// TestAC1POSIXSealHandlePlainObjectStillSeals120r3 is the other direction: a leg that
// refused everything would pass the cases above and be worthless. With the name
// carrying the caller's own object, sealHandle must return nil and land the mode -
// and this case must NOT change colour under either control.
func TestAC1POSIXSealHandlePlainObjectStillSeals120r3(t *testing.T) {
	_, name, witness := sealHandleFixture120r3(t)
	f := openHandleOn120r3(t, name)
	if err := sealHandle(f); err != nil {
		t.Errorf("AC#1 RED: the honest path was refused: %v (leg that spoke = %s)", err, spokeLeg120r3(err))
	}
	if got := modeOfName120r3(name); got != "-rw-------" {
		t.Errorf("AC#1 RED: the object the caller named is at %s, want -rw-------", got)
	}
	if got := modeOfName120r3(witness); got != "-rw-rw-rw-" {
		t.Errorf("AC#1 RED: the witness was narrowed to %s on a spelling with no link in it", got)
	}
}
