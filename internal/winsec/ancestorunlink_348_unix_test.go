//go:build unix

// Ticket 348 AC#2 + AC#3, leg 342-r1: the DELETION leg - the route's own refusal-path unlink,
// measured through the production entry point. Shared world and liveness rules: see
// ancestorworld_342_348_unix_test.go. This file judges ONE numerator: deleted.
//
// THE DAMAGE THIS GATE OWNS. privateFile unlinks on two failures - the seal refusal and the write
// failure - and until this landing both did it BY NAME (`_ = os.Remove(resolved)`; there were TWO of
// them: ticket 348's own table named one, 342-v1's verdict corrected it to two, and the dispatch
// names exactly that as the trap a leg falls into by changing only one). A by-name removal
// re-resolves the whole spelling, so an ancestor carrying a link sends the unlink into somebody
// else's tree and removes a name nobody named. 342-a1's A3 measured it DETERMINISTICALLY - unlink
// err=<nil> and the foreign witness gone (probes/342/a1/logs/24-container-A3-unlink-through-ancestor
// .log). Not a race to be won; a route to be walked.
//
// WHY THE PARENT PIN IS THE FIX, AND THE COUPLING THAT MAKES IT ONE. unlinkat(parentFd, base, 0)
// resolves only `base`, inside the object the descriptor names, so the ancestor is not in the path at
// all. That holds only while the pin is STILL OPEN when the refusal fires: releasing it as soon as the
// leaf was opened would leave this unlink nothing to be relative to. That is why winsec.go's deferred
// close sits after both removals instead of inside the open step, and it is the coupling 342-v1's
// answer (v) named and the dispatch re-states as the shape a leg gets wrong by copying the probe.
//
// THE PAIR. The deterministic case below has the teeth this ticket needs: on the pre-landing bytes it
// deletes the foreign witness every time it runs, and on the landing it cannot. The raced case after
// it holds ticket 348's own liveness口径 over the same world ticket 342 drives, and reads a different
// numerator, so neither ticket borrows the other's green (AC#3 verbatim).
package winsec

import (
	"os"
	"path/filepath"
	"testing"
)

// TestAC2POSIXRefusalPathWithASwappedAncestorNeverDeletesAForeignFile342r1 is the deterministic face:
// one production call, the ANCESTOR swapped inside the window the seal's own walk owns, then the
// foreign tree read back.
//
// WHY THE SWAP GOES IN THROUGH THE PACKAGE'S OWN SEAM. The placement floor refuses a steady-state
// ancestor link before the act, so pre-planting it means the route never reaches the unlink this
// ticket is about; the removal only runs on the route's own refusal, which means the swap has to land
// AFTER the open and before the identity re-check. inIdentityWindow120r3 is exactly that window and it
// is already the package's declared fixture: it runs the REAL pinned seal first and keeps its verdict
// (sealhandle_privateroute_120_other_test.go:99 is its prior use for the leaf; this file reuses the
// seam, not the shape).
//
// WHY `park` AND NOT `remove`: rename(2) cannot replace a non-empty directory and a directory with a
// symlink cannot simply be written over, so the real ancestor is MOVED OUT - one rename, which keeps
// its contents, including the caller's leaf - and only then is the link planted at the vacated name.
// The moved-out directory is what the pinned descriptor still names on the landing side, and that is
// the whole point of this case: the landing's unlink removes a leaf inside an object the caller owned,
// while the pre-landing unlink goes by name, re-resolves through the fresh link, and takes the foreign
// witness with it.
func TestAC2POSIXRefusalPathWithASwappedAncestorNeverDeletesAForeignFile342r1(t *testing.T) {
	base := SealableTempDirForTest124(t)
	data := filepath.Join(base, "data")
	sub := filepath.Join(data, "sub")
	park := filepath.Join(base, "parked-real-sub")
	name := filepath.Join(sub, ancLeaf342r1)
	foreign := filepath.Join(data, ancOutsideDir342r1)
	for _, d := range []string{data, sub, foreign} {
		if err := os.MkdirAll(d, 0o700); err != nil {
			t.Fatalf("stage %s: %v", d, err)
		}
	}
	witness := filepath.Join(foreign, ancLeaf342r1)
	if err := os.WriteFile(witness, []byte(ancWitnessBody342r1), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(witness, 0o666); err != nil {
		t.Fatal(err)
	}
	before, err := os.Lstat(witness)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(name, []byte(ancOwnBody342r1), 0o600); err != nil {
		t.Fatal(err)
	}

	planted := false
	var plantErr string
	inIdentityWindow120r3(t, func() {
		if err := os.Rename(sub, park); err != nil {
			plantErr = "rename the real ancestor out: " + err.Error()
			return
		}
		// A BARE sibling name, exactly like the raced instrument's relative form: no ".." element, so
		// nothing here is refused by a containment rule and the only leg that can answer is this
		// package's own identity gate (see the world file's two-forms note).
		if err := os.Symlink(ancOutsideDir342r1, sub); err != nil {
			plantErr = "plant the ancestor link: " + err.Error()
			return
		}
		info, lErr := os.Lstat(sub)
		if lErr != nil || info.Mode()&os.ModeSymlink == 0 {
			plantErr = "the ancestor is not a symlink after the swap"
			return
		}
		planted = true
	})

	err = PrivateFile(name, []byte(ancPayload342r1), 0o600)
	if plantErr != "" {
		t.Fatalf("AC#2 RED instrument: %s, so the route was never driven through a swapped ancestor", plantErr)
	}
	if !planted {
		t.Fatal("AC#2 RED instrument: the plant never landed")
	}
	if err == nil {
		t.Fatal("AC#2 RED instrument: the swapped ancestor produced no refusal, so the route's unlink never " +
			"ran and the reads below would prove nothing about it")
	}
	t.Logf("AC#2' refusal verbatim (the leg that answered): %v", err)

	after, statErr := os.Lstat(witness)
	if statErr != nil {
		t.Errorf("AC#2 RED: the foreign witness is GONE after the refusal-path unlink (%v), so the removal "+
			"re-resolved the caller's spelling through the swapped ancestor and deleted a name in somebody "+
			"else's tree - ticket 348's shape, the second damage of the pair", statErr)
		return
	}
	if body, rErr := os.ReadFile(witness); rErr != nil || string(body) != ancWitnessBody342r1 {
		t.Errorf("AC#2 RED: the foreign witness holds %q (err %v), want its own bytes intact", body, rErr)
	}
	// The mode of the survivor is REPORTED and not gated here: a seal-side walk can narrow what a
	// relative link leads it to, and that leg belongs to ticket 120, not to this landing. Saying so out
	// loud is what stops a green deletion row from being read as a closed family.
	t.Logf("AC#2' survivor read-back: size %d (was %d), mode %v (was %v), mtime moved=%v - a narrowed mode on "+
		"this path is ticket 120's surface and is NOT claimed closed by tickets 342/348",
		after.Size(), before.Size(), after.Mode(), before.Mode(), !after.ModTime().Equal(before.ModTime()))
}

// TestAC2POSIXPrivateFileRaceThroughASwappedAncestorNeverDeletesAForeignFile342r1 is the same route
// under the racer with ticket 348's liveness口径 attached: the deleted numerator may not be read as a
// green unless the link state was PROVED and the pinned walk answered it.
//
// The spelling is the RELATIVE one (see the world file's two-forms note): it is the form on which the
// pre-landing bytes actually reach their unlink - the open lands in the foreign object and the seal's
// own walk then refuses by name - so `deleted` has a way to be non-zero here, and a zero is a reading.
func TestAC2POSIXPrivateFileRaceThroughASwappedAncestorNeverDeletesAForeignFile342r1(t *testing.T) {
	w := newAncWorld342r1(t, ancTargetRel)
	st := w.run(ancAttempts342r1, func() error {
		return PrivateFile(w.name, []byte(ancPayload342r1), 0o600)
	})
	st.print(t, "AC#2 (ticket 348, deleted gate, relative ancestor link under the racer)")
	if st.deleted > 0 {
		t.Errorf("AC#2 RED: %d of %d attempts found the foreign witness GONE, so the refusal-path unlink in "+
			"privateFile re-resolved the spelling through an ancestor swapped into a link and removed a "+
			"name in somebody else's tree (first: %s)", st.deleted, st.attempts, st.firstProof)
	}
	if st.provedLinkCycles == 0 || st.liveness == 0 {
		t.Errorf("AC#2 RED instrument: the racer proved %d link-state cycles and the pinned walk refused %d "+
			"attempts across %d (flipErrs=%d, first %q), so the zero above is not evidence about a removal "+
			"that was ever reachable", st.provedLinkCycles, st.liveness, st.attempts, st.flipErrs, st.firstFlip)
	}
	if st.narrowed > 0 {
		t.Errorf("AC#2 RED: %d attempts left the foreign witness narrowed off 0666, so a seal followed the "+
			"flipped ancestor too - ticket 120's surface, re-opened on this route", st.narrowed)
	}
}
