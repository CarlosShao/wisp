//go:build !windows

// Ticket 120, residual gap 1, second half: the two legs of sealHandle's inline
// re-check (winsec_other.go:269 and :272) reached through the ENTRY POINT THAT
// ACTUALLY HAS A CALLER, privateFile -> sealHandle, plus the reading that says the
// window is reached by an unsynchronised racer too and not only by the fixture's
// synchronised plant in sealhandle_identity_120_other_test.go.
//
// Why the route matters: sealHandle exists for one caller shape - the process
// already holds a descriptor and is about to write through it - so a fixture that
// only ever hands it a hand-made descriptor proves the leg's arithmetic, not that
// anything in this package can reach it. Here the bytes, the refusal and the
// unlink-on-refusal (winsec.go:117-121) are the production ones: the payload must
// not be reachable under any name afterwards, and on the link shape the witness the
// planted link points at must keep both its mode and its bytes.
//
// Both route cases plant the swap in the window through the package's declared
// seam var applyDescriptor, calling the real applyDescriptorPOSIX first and keeping
// its verdict - see the long note in sealhandle_identity_120_other_test.go for why
// a pre-planted link can never reach :269 (the walk refuses it, which is the gap's
// cause) and why a blind racer cannot be PROMISED a landing.
//
// The blind racer is the third case: it flips the leaf between a real file and a
// relative sibling link by rename(2) only (so the name is never absent, which is
// 120-r2's fix for its own dead instrument), drives PrivateFile over that spelling
// N times, and reports which leg answered each attempt. Its judgement is the safety
// half only - the witness never narrows, whatever the distribution - because 120-v1
// section 2's ruling 3 holds that a probabilistic landing count must not become a
// hard gate. What it buys is the statement "sealHandle's own legs answer on their
// own, with no seam help", measured on the same table as the control pair.
//
// uid 0 in the container again: no permission face is asserted here either.
package winsec

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// windowProbeAttempts120r3 is the blind racer's denominator, reported in every row.
const windowProbeAttempts120r3 = 2000

// assertRefusedWriteLeavesNothing120r3 is AC#5's claim measured on THIS world: after
// a refused seal the payload must not be readable under any name in the tree, and
// the spelling the caller asked for must not be holding it either (reading it
// follows whatever the name carries now, which is the whole point of the leg).
func assertRefusedWriteLeavesNothing120r3(t *testing.T, dir, name, payload string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read the tree back: %v", err)
	}
	reachable := 0
	for _, e := range entries {
		if !e.Type().IsRegular() && e.Type()&os.ModeSymlink == 0 {
			continue
		}
		p := filepath.Join(dir, e.Name())
		body, rErr := os.ReadFile(p)
		if rErr != nil {
			t.Logf("AC#1' %s is not readable (%v), so no payload claim is made about it", e.Name(), rErr)
			continue
		}
		reachable++
		if strings.Contains(string(body), payload) {
			t.Errorf("AC#1 RED: the refused write left the payload in %s (mode %s, kind %s)",
				p, modeOfName120r3(p), kindOfName120r3(p))
		}
	}
	if body, rErr := os.ReadFile(name); rErr == nil && strings.Contains(string(body), payload) {
		t.Errorf("AC#1 RED: the spelling the caller named reads back the payload it refused to write")
	}
	t.Logf("AC#1' tree after the refusal: %d readable entries, none carrying the payload", reachable)
}

// TestAC1POSIXPrivateFileRefusesWhenNameCarriesAnotherRealObject120r3 is the SameFile
// face on the production route: privateFile opens the caller's object, the pinned
// seal narrows what the name carries at that instant, and only after that is a
// second REAL object renamed onto the name. sealHandle must refuse the write, and
// privateFile must unlink the name it refused (the swap-in object goes with it).
func TestAC1POSIXPrivateFileRefusesWhenNameCarriesAnotherRealObject120r3(t *testing.T) {
	dir, name, _ := sealHandleFixture120r3(t)
	const payload = "artifact bytes that must not land in a file nobody sealed"

	swapIn := filepath.Join(dir, "swap-in-route-120r3.bin")
	writeWide120r3(t, swapIn, "the object swapped into the name")
	inIdentityWindow120r3(t, func() {
		if err := os.Rename(swapIn, name); err != nil {
			t.Errorf("fixture: swap a second real object onto the name: %v", err)
		}
	})

	err := PrivateFile(name, []byte(payload), 0o600)
	assertLegSpoke120r3(t, err, whyHandleSwap120r3, "private-file-route-name-carries-another-real-object")
	if got := kindOfName120r3(name); got != "absent" {
		t.Errorf("AC#1 RED: the refused write left the name carrying a %s, want it unlinked (%s)", got, err)
	}
	assertRefusedWriteLeavesNothing120r3(t, dir, name, payload)
}

// TestAC1POSIXPrivateFileRefusesWhenNameBecomesLinkAfterPinnedSeal120r3 is the link
// face on the production route: after the pinned seal has run for real on the
// caller's object, the name is replaced by a relative sibling link at the witness.
// The refusal has to name sealHandle's own leg, the witness has to keep its mode and
// its bytes, and the unlink must remove the link rather than what it points at.
func TestAC1POSIXPrivateFileRefusesWhenNameBecomesLinkAfterPinnedSeal120r3(t *testing.T) {
	dir, name, witness := sealHandleFixture120r3(t)
	const payload = "artifact bytes that must not be written through a swapped name"

	inIdentityWindow120r3(t, func() {
		if err := os.Remove(name); err != nil {
			t.Errorf("fixture: take the name away: %v", err)
			return
		}
		if err := os.Symlink(filepath.Base(witness), name); err != nil {
			t.Errorf("fixture: plant the sibling link: %v", err)
		}
	})

	err := PrivateFile(name, []byte(payload), 0o600)
	assertLegSpoke120r3(t, err, whyHandleLink120r3, "private-file-route-name-becomes-a-sibling-link")

	if got := modeOfName120r3(witness); got != "-rw-rw-rw-" {
		t.Errorf("AC#1 RED: the witness the link pointed at is at %s, so a seal followed it", got)
	}
	if body, rErr := os.ReadFile(witness); rErr != nil {
		t.Errorf("AC#1 RED: the witness is no longer readable, so the refused write reached through it: %v", rErr)
	} else if string(body) != "not this write's object" {
		t.Errorf("AC#1 RED: the witness holds %q, so the refused write landed in it", body)
	}
	if got := kindOfName120r3(name); got != "absent" {
		t.Errorf("AC#1 RED: the refused write left the name carrying a %s, want the link unlinked (%v)", got, err)
	}
	assertRefusedWriteLeavesNothing120r3(t, dir, name, payload)
}

// TestAC1POSIXPrivateFileWindowIsRacedWithoutTheSeam120r3 is the reading that the
// window is not a fixture invention. A racer flips the caller's name between a real
// file and a relative sibling link, exclusively by rename(2) so the name is present
// in every state, and PrivateFile runs over that spelling windowProbeAttempts120r3
// times. Every attempt is classified by the leg that answered it and printed; the
// distribution is a READING, not a gate, and the one hard assertion is the safety
// half that has to hold whatever the distribution is: the witness is never narrowed.
func TestAC1POSIXPrivateFileWindowIsRacedWithoutTheSeam120r3(t *testing.T) {
	dir, name, witness := sealHandleFixture120r3(t)
	stageReal := filepath.Join(dir, "stage-real-120r3.bin")
	stageLink := filepath.Join(dir, "stage-link-120r3.bin")

	var flipErrs int
	var firstFlip string
	note := func(err error) {
		if err == nil {
			return
		}
		flipErrs++
		if firstFlip == "" {
			firstFlip = err.Error()
		}
	}
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer func() {
			if p := recover(); p != nil {
				note(fmt.Errorf("racer panic: %v", p))
			}
		}()
		for {
			select {
			case <-stop:
				return
			default:
			}
			// Two staged names, re-staged before each use: rename replaces the
			// destination atomically and does not follow a link standing there, so
			// the caller's name is never the thing that is missing.
			if err := os.WriteFile(stageReal, []byte("the tree's own bytes"), 0o600); err != nil {
				note(err)
				continue
			}
			note(os.Rename(stageReal, name))
			if err := os.Symlink(filepath.Base(witness), stageLink); err != nil && !errors.Is(err, fs.ErrExist) {
				note(err)
				continue
			}
			note(os.Rename(stageLink, name))
		}
	}()
	defer func() {
		close(stop)
		<-done
	}()

	counts := map[string]int{}
	var firstOther string
	for i := 0; i < windowProbeAttempts120r3; i++ {
		err := PrivateFile(name, []byte("window probe payload"), 0o600)
		leg := spokeLeg120r3(err)
		switch {
		case err == nil:
			leg = "returned nil (the flip missed the window)"
		case leg == "some other failure":
			if firstOther == "" {
				firstOther = err.Error()
			}
		}
		counts[leg]++
	}
	close(stop)
	<-done

	t.Logf("AC#1' natural-window reading, denominator=%d: attempts classified by the leg that answered",
		windowProbeAttempts120r3)
	for _, leg := range sortedLegs120r3(counts) {
		t.Logf("AC#1'   %-72s = %d", leg, counts[leg])
	}
	if firstOther != "" {
		t.Logf("AC#1' natural-window reading, first other-failure verbatim: %s", firstOther)
	}
	if flipErrs > 0 {
		t.Errorf("AC#1 RED instrument: the racer's own flip recorded %d error(s) (first: %s), "+
			"so nothing was raced and the rows above are not a reading", flipErrs, firstFlip)
	}
	if got := modeOfName120r3(witness); got != "-rw-rw-rw-" {
		t.Errorf("AC#1 RED: after %d raced writes the witness stands at %s, so a seal or a write followed "+
			"the flipped link into a foreign object", windowProbeAttempts120r3, got)
	}
	if string(mustReadFile120r3(t, witness)) != "not this write's object" {
		t.Errorf("AC#1 RED: the witness no longer holds its own bytes")
	}
}

// sortedLegs120r3 prints the classification roster in a stable order, so two runs of
// this reading can be compared line by line.
func sortedLegs120r3(counts map[string]int) []string {
	legs := make([]string, 0, len(counts))
	for leg := range counts {
		legs = append(legs, leg)
	}
	sort.Strings(legs)
	return legs
}

func mustReadFile120r3(t *testing.T, path string) []byte {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return body
}
