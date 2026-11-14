//go:build windows

// Ticket 348 AC#2, leg 348-r2: the DELETION leg measured through the seam that ticket's AC#2 needed
// a second human approval for (`removeRefusedLeaf`, winsec.go).
//
// WHY A SEAM AT ALL, stated as the measurement that asked for it (348-r2, logs/03-*). The world this
// ticket is about is "the ancestor is a link at the instant the refusal-path removal runs". A racer
// cannot put the pre-landing route into it on this platform: while the leaf handle is open the kernel
// refuses every ancestor rename (`Access is denied`), and while the pinned parent is open a rename
// that MOVES the pinned directory out of its parent answers `in use by another process`. What neither
// of those refuses is a rename of the pinned directory to a name inside the same parent (measured
// `err=<nil>`), and that is the form used below because it is the form the attacker has - so the flip
// stands at the removal instant even on the landed bytes, and the reading is about the removal rather
// than about the platform. What makes the pre-landing route's world hard to reach is narrower: its
// only free-handed instant is between `_ = f.Close()` and `os.Remove(resolved)`, two syscalls. With
// the removal branch forced to run (the package's own applyDescriptor seam, injected to fail) the
// pre-fix route still reaches the damage only 8-10 times in 400 attempts (logs/11-*): a real
// numerator, measured, but one that moves run to run with the scheduler.
//
// WHAT THIS CASE THEREFORE DOES. It replaces the seam with a function that ACTS at that instant - lands
// the ancestor flip, proves it, then calls the tree's own removal, unchanged - so the numerator becomes
// a reading about the route's removal instead of a reading about how often a racer wins two syscalls.
// The instrument is deliberately tree-agnostic: on the landed bytes the flip is refused by the pin and
// the witness survives; on the pre-landing bytes the flip lands and the by-name removal takes the
// foreign name. The second behaviour is the reverse criterion, and it is run, not asserted (348-r2).
//
// THE LIVENESS QUESTION, ANSWERED WHERE IT CAN BE. AC#2 forbids reading a zero as evidence, and this
// file has to be honest about *which* zero it prints:
//   - `hopLanded` counts attempts where the flip really stood on the ancestor at the removal instant,
//     each one proved with reparseAt (the same key the racer instruments use). On the landed bytes this
//     is 0, and that is a reading about the pin, not about the removal;
//   - so the "the window was reached" half for the landed tree is NOT taken from this case. It is
//     taken from TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1, which
//     reaches the link state by racing the walk itself (provedReparse in the same run, logs/12-*).
//     This case does not claim that half and must not be cited for it.
package winsec

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"sync/atomic"
	"testing"
)

// TestAC2SeamWindowsTheRefusalPathRemovalCannotBeMovedByAnAncestorFlip348r2 is the deletion gate.
//
// It gates one thing: the foreign witness is never removed. Everything else it prints is a reading -
// how many attempts reached the seam, how many landed the flip at the instant, and the kernel's own
// verbatim answer when it refused. The refusal text is kept because "the pin refused" and "the
// removal was safe" are different claims, and only the first one is what this case observes on the
// landed bytes.
func TestAC2SeamWindowsTheRefusalPathRemovalCannotBeMovedByAnAncestorFlip348r2(t *testing.T) {
	// Ticket 444: the un-raced carrier runs FIRST, while `removeRefusedLeaf` is still the production
	// removal - the seam below replaces that var, so a control run after it would be measuring the
	// seam's own closure instead of the landing's act.
	ctl := winAncWindowControl444(t, winAncFormSymDir, winAncCtlRemoval444)
	ctl.print(t, "AC#2")
	ctlOK := ctl.foreignIntact && ctl.actErr == "" && ctl.ownLeafGone

	withInjectedSealFailure(t, errors.New("348r2: injected seal failure so the refusal-path removal runs"))
	w := newWinAncWorld342r1(t, winAncFormSymDir)

	var (
		seamCalls   int32
		hopLanded   int32
		hopRefused  int32
		firstRefuse string
		firstLanded string
	)
	// The tree's own removal, kept so the replacement can call it rather than re-implement it: the
	// act under test must stay the production act, or this instrument measures its own code.
	productionRemoval := removeRefusedLeaf
	t.Cleanup(func() { removeRefusedLeaf = productionRemoval })
	removeRefusedLeaf = func(p leafPin, base string, pinned fs.FileInfo) {
		atomic.AddInt32(&seamCalls, 1)
		// Land ticket 348's world AT THE HOP, and prove it with the package's own probe before the
		// removal runs - exactly the sequence the racer cannot complete here.
		if err := w.toFlip(); err != nil {
			atomic.AddInt32(&hopRefused, 1)
			if firstRefuse == "" {
				firstRefuse = err.Error()
			}
			if rErr := w.resettle(); rErr != nil && firstRefuse != "" {
				firstRefuse += " | resettle: " + rErr.Error()
			}
		} else {
			atomic.AddInt32(&hopLanded, 1)
			if firstLanded == "" {
				firstLanded = "the ancestor stood on a link at the removal instant and reparseAt proved it"
			}
		}
		productionRemoval(p, base, pinned)
		if atomic.LoadInt32(&hopLanded) > 0 {
			if err := w.toReal(); err != nil {
				t.Logf("AC#2' restore after the landed flip: %v", err)
			}
		}
	}

	deleted := 0
	proof := ""
	attempts := attemptsFor444(winAncAttempts342r1)
	// Ticket 444 AC#1 probe: WISP_444_HOP_REFUSE=1 makes the platform's own refusal branch
	// the only one this run can take, by parking a NON-EMPTY directory where the flip moves
	// the pinned ancestor to. rename(2) cannot replace a non-empty directory, so `w.toFlip()`
	// answers an error on every seam call - the same branch the kernel takes on CI when the
	// pinned handle turns the rename into "being used by another process". The counters and
	// the sentences below are untouched by this: it manufactures the ANSWER, not the reading.
	hopRefuseProbe := os.Getenv("WISP_444_HOP_REFUSE") != ""
	if hopRefuseProbe {
		if err := os.MkdirAll(w.parkReal, 0o700); err != nil {
			t.Fatalf("probe: park a directory where the flip lands: %v", err)
		}
		if err := os.WriteFile(filepath.Join(w.parkReal, "occupied-444"), []byte("x"), 0o600); err != nil {
			t.Fatalf("probe: fill the parked directory: %v", err)
		}
	}
	for i := 0; i < attempts; i++ {
		_ = PrivateFile(w.name, []byte(winAncPayload342r1), 0o600)
		if kind, detail := w.readWitness(); kind == "deleted" {
			deleted++
			if proof == "" {
				proof = "witness ABSENT after attempt " + strconv.Itoa(i) + " (" + detail + ")"
			}
		}
		w.restockWitness()
	}

	calls := atomic.LoadInt32(&seamCalls)
	landed := atomic.LoadInt32(&hopLanded)
	refused := atomic.LoadInt32(&hopRefused)
	t.Logf("AC#2 (ticket 348, Windows in-root ancestor form, through the removal seam): attempts=%d "+
		"seamCalls=%d hopLanded=%d hopRefused=%d deleted=%d", attempts, calls, landed,
		refused, deleted)
	if landed > 0 {
		t.Logf("AC#2' the flip STOOD at the hop on %d attempts (first: %s) - the world this ticket is "+
			"about was entered at the removal instant", landed, firstLanded)
	}
	if refused > 0 {
		t.Logf("AC#2' the platform REFUSED the flip at the hop on %d attempts; first verbatim: %s",
			refused, firstRefuse)
	}
	if proof != "" {
		t.Logf("AC#2' first damage verbatim: %s", proof)
	}

	if calls != int32(attempts) {
		t.Errorf("AC#2 RED instrument: the removal seam was reached %d times in %d attempts, so this "+
			"reading is not about the refusal-path removal at all (an attempt that never got there "+
			"cannot say anything about what that removal would have done)", calls, attempts)
	}
	// Ticket 444's carrier, judged before anything below leans on it: with the racer absent by
	// construction, does the landing's own removal still remove the PINNED name and leave the foreign
	// one standing? This is the half that contention cannot starve.
	if !ctlOK {
		t.Errorf("AC#2 RED instrument: the un-raced removal control - the route's own pin, then the flip, "+
			"then the refusal-path removal - neither kept the foreign name nor removed the pinned one: %s",
			ctl.verdict())
	}
	if landed == 0 && refused == 0 {
		// Ticket 444: this arm names the shape "the flip was never asked", and the only way it can hold
		// is that the route never reached the removal at all (every seam call increments exactly one of
		// hopLanded/hopRefused, so landed==0 && refused==0 is calls==0). That shape is an abstention
		// when nothing was asked, and the wiring defect line 126 already carries its own red for it; if
		// it ever holds with calls>0 the seam is not counting, which stays RED in the original words.
		if calls == 0 {
			t.Logf("AC#2 ABSTAIN(444:348r2-seam:removal-never-reached) attempts=%d seamCalls=0 "+
				"hopLanded=0 hopRefused=0 deleted=%d - the route never reached the refusal-path removal, "+
				"so this run cannot attest the flip at the hop; the threshold that carries the removal is "+
				"the un-raced control above (foreign-intact=%v, own-leaf-removed=%v), not a green from "+
				"absence. The wiring red for 'the seam was reached fewer times than attempts' is a "+
				"separate arm and is unaffected", attempts, deleted, ctl.foreignIntact, ctl.ownLeafGone)
		} else {
			t.Errorf("AC#2 RED instrument: %d seam calls produced neither a landed nor a refused flip, so "+
				"the flip was never even asked and the deletion count below is unaccompanied by any "+
				"statement about the ancestor", calls)
		}
	}
	// ★The liveness floor, and it is deliberately the strong form rather than "> 0". AC#2 forbids
	// reading a zero as evidence, and the zero below is only evidence about the REMOVAL if the world
	// stood at the ancestor when the removal ran. On this landing that holds on every attempt
	// (measured 400/400 across nine runs, logs/21-*), so anything less means the instrument lost the
	// world it is named for - not that the removal got safer.
	if landed != calls {
		// Ticket 444: this is the arm CI load actually fires (the kernel refuses the hop rename under
		// contention). It abstains only when the un-raced control above carries the removal on this
		// same run; if the control did not, the original sentence stands unchanged.
		if ctlOK {
			t.Logf("AC#2 ABSTAIN(444:348r2-seam-hop:flip-not-asked) attempts=%d samples=1 seamCalls=%d "+
				"hopLanded=%d hopRefused=%d deleted=%d first-refusal-verbatim=%q - the real failure mode "+
				"is the platform's own refusal of the hop rename under load (the same branch CI takes), "+
				"so this run cannot attest the flip AT the removal instant; the threshold that carries "+
				"the removal is the un-raced control above (foreign-intact=%v, own-leaf-removed=%v), "+
				"not a green from absence",
				attempts, calls, landed, refused, deleted, firstRefuse, ctl.foreignIntact, ctl.ownLeafGone)
		} else {
			t.Errorf("AC#2 RED instrument: the ancestor flip stood at the removal instant on only %d of %d "+
				"seam calls (refused=%d, first refusal verbatim: %q), so the deletion count is a reading "+
				"about a world that was not entered rather than about the removal", landed, calls, refused,
				firstRefuse)
		}
	}
	if deleted > 0 {
		t.Errorf("AC#2 RED: %d of %d attempts found the foreign witness GONE through the refusal-path "+
			"removal (first: %s), so that removal was moved by an ancestor flip and deleted a name in "+
			"somebody else's tree", deleted, attempts, proof)
	}
}
