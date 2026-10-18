//go:build !windows

// Ticket 120 residual gap 1: the same two sealHandle faces reached through the entry
// point that actually has callers - privateFile -> sealHandle (winsec.go:117) - plus the
// reading that says the window is not a fixture invention.
// A hand-made descriptor (sealhandle_identity_120_other_test.go) proves the leg's
// arithmetic; this file proves that a production call refuses, that its bytes do not
// land, and that the unlink-on-refusal removes the planted name and not what that link
// points at. The world, and the honesty of the seam-windowed plant, are stated once in
// sealhandle_window_120_other_test.go. The third case is the blind racer: same fixture,
// no seam help, the flip driven by its own goroutine. Its classification rows are
// READINGS, because 120-v1 section 2's third ruling holds that a landing count nobody
// can promise must not become a gate; what is judged is the racer's own health and the
// invariants these legs do own.
package winsec

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
)

// windowProbeAttempts120r3 is the blind racer's denominator, printed with every row.
const windowProbeAttempts120r3 = 2000

// assertRefusedWriteLeavesNothing120r3 is AC#5's promise measured on THIS world: after a
// refused seal the payload must not be readable under any name in the tree - and reading
// the spelling the caller asked for follows whatever the name carries now, which is the
// whole subject of the leg.
func assertRefusedWriteLeavesNothing120r3(t *testing.T, dir, name, payload string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read the tree back: %v", err)
	}
	readable := 0
	for _, e := range entries {
		p := filepath.Join(dir, e.Name())
		if !e.Type().IsRegular() && e.Type()&os.ModeSymlink == 0 {
			continue
		}
		body, rErr := os.ReadFile(p)
		if rErr != nil {
			t.Logf("AC#1' %s is not readable (%v), so no payload claim is made about it", e.Name(), rErr)
			continue
		}
		readable++
		if strings.Contains(string(body), payload) {
			t.Errorf("AC#1 RED: the refused write left the payload in %s (kind %s, mode %s)",
				p, kindOfName120r3(p), modeOfName120r3(p))
		}
	}
	if body, rErr := os.ReadFile(name); rErr == nil && strings.Contains(string(body), payload) {
		t.Errorf("AC#1 RED: the spelling the caller named reads back the payload its refusal would not write")
	}
	t.Logf("AC#1' tree after the refusal: %d readable entries, none carrying the payload", readable)
}

// TestAC1POSIXPrivateFileRefusesWhenNameCarriesAnotherRealObject120r3 is the SameFile
// face on the production route: privateFile opens the caller's object, the pinned seal
// narrows what the name carries at that instant, and only after that is a second REAL
// object renamed onto the name. The write must be refused by :272, and privateFile
// must unlink the name it refused.
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
		t.Errorf("AC#1 RED: the refused write left the name carrying a %s, want it unlinked; err=%v", got, err)
	}
	assertRefusedWriteLeavesNothing120r3(t, dir, name, payload)
}

// TestAC1POSIXPrivateFileRefusesWhenNameBecomesLinkAfterPinnedSeal120r3 is the link
// face on the production route: after the pinned seal has run for real on the caller's
// object the name is replaced by a relative sibling link at the witness. The refusal
// has to name :269, the witness has to keep its mode AND its bytes, and the unlink on
// the refusal path has to remove the link rather than what it points at.
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
		t.Errorf("AC#1 RED: the refused write left the name carrying a %s, want the link unlinked; err=%v", got, err)
	}
	assertRefusedWriteLeavesNothing120r3(t, dir, name, payload)
}

// TestAC1POSIXPrivateFileWindowIsRacedWithoutTheSeam120r3 is the reading that the
// window is not a fixture invention. A racer flips the caller's name between a real
// file and a relative sibling link using rename(2) only, so the name is present in
// every state (120-r2's fix for its own dead instrument, reused in shape here), and
// PrivateFile runs over that spelling windowProbeAttempts120r3 times. Every attempt is
// classified by the leg that answered it and printed. The distribution is a READING;
// what is judged is the racer's health and the two invariants that hold whatever it says.
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
	var stopOnce sync.Once
	stopRacer := func() {
		stopOnce.Do(func() { close(stop) })
		<-done
	}
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
			// destination atomically and does not follow a link standing at it, so
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
	defer stopRacer()

	counts := map[string]int{}
	var firstOther string
	for i := 0; i < windowProbeAttempts120r3; i++ {
		err := PrivateFile(name, []byte("window probe payload"), 0o600)
		leg := spokeLeg120r3(err)
		switch {
		case err == nil:
			leg = "returned nil: the flip missed every check"
		case leg == "some other failure":
			leg = "other failure"
			if firstOther == "" {
				firstOther = err.Error()
			}
		}
		counts[leg]++
	}
	stopRacer() // settled before the counters below are read, so this is not a race on the reading

	t.Logf("AC#1' natural-window reading, denominator=%d, attempts classified by the leg that answered:",
		windowProbeAttempts120r3)
	for _, leg := range sortedLegs120r3(counts) {
		t.Logf("AC#1'   %s = %d", leg, counts[leg])
	}
	if firstOther != "" {
		t.Logf("AC#1' natural-window reading, first other-failure verbatim: %s", firstOther)
	}
	if flipErrs > 0 {
		t.Errorf("AC#1 RED instrument: the racer's own flip recorded %d error(s) (first: %s), so nothing "+
			"was raced and the rows above are not a reading", flipErrs, firstFlip)
	}
	// What the race owes, judged: whatever missed the window, no seal may have
	// narrowed the object a flipped link pointed at, and none of the raced payloads
	// may have reached it.
	if got := modeOfName120r3(witness); got != "-rw-rw-rw-" {
		t.Errorf("AC#1 RED: after %d raced attempts the witness stands at %s, so a seal followed a flipped "+
			"link into a foreign object", windowProbeAttempts120r3, got)
	}
	body := string(mustReadFile120r3(t, witness))
	switch {
	case body == "not this write's object":
		t.Logf("AC#1' witness holds its own bytes after %d raced attempts: no open, seal or write reached "+
			"it through a flipped link", windowProbeAttempts120r3)
	case strings.Contains(body, "window probe payload"):
		t.Errorf("AC#1 RED: a raced payload landed IN the foreign object the flipped link pointed at: %q", body)
	default:
		// Counted and named, NOT judged: an attempt whose flip landed before the open
		// gets O_TRUNC through that link. Since ticket 322 the POSIX LEAF leg of that
		// open is closed by oNoFollowPlatform (winsec_other.go:103-113), so a leaf link
		// swapped into the name answers ELOOP instead of truncating. What stays open on
		// this route is the ANCESTOR leg - O_NOFOLLOW speaks only the last name - and
		// that is nobody's gate either. Nobody may read any of this as the window closed.
		t.Logf("AC#1' R-322 SHAPE OBSERVED, counted and not judged: the witness now holds %q, so at least "+
			"one open truncated it through a flipped link (ticket 322 closed only the POSIX leaf leg; "+
			"the ancestor leg stays open). The sealHandle legs answered the write: see the rows above.", body)
	}
}

// sortedLegs120r3 prints the classification roster in a stable order so two runs of
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
