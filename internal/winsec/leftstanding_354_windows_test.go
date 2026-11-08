//go:build windows

// Ticket 354 AC#2, leg 354-r1: the DETERMINISTIC instrument for the residual ticket 353's form 乙
// ratified - after the refusal-path removal's identity re-check refuses, the caller's NAME STAYS, and
// in the attacked state it stays carrying the object the caller never named. This leg names that
// reading as its own numerator, `leftStandingForeignName`, and GATES it. Ticket 353's instrument
// (nametopup_353_windows_test.go) only prints the shape (its DAMAGE row) and gates the pre-乙 damage
// branch; nothing on these bytes gates the residual itself, so a landing that silently changed what
// stands under the caller's name would turn nothing red. The paired 去牙 of AC#3 (this leg's evidence
// file) proves the numerator bites: pull the re-check tooth and the numerator collapses to 0 with a
// red; pull the removal tooth and the CONTROL below goes red instead. The two teeth are pulled at two
// different sites and neither pull is committed - the landed bytes keep both.
//
// DETERMINISM, stated as the construction rather than claimed: the SAME fixed pre-state on every shot
// (fixed per-shot names, fixed bodies), the move made synchronously inside the seam this package
// already declares (applyDescriptor, winsec_windows.go:35-38: "Tests replace it to inject a failure
// (AC#5)"), and the world is winDetSwapWorld342r2 reused VERBATIM (newWinDetSwapWorld342r2 /
// installWindow) - no racer, no timing, no sleeps, no new production seam, no new exported name. The
// plant runs the REAL seal first, keeps its verdict, vacates the caller's name, renames the unnamed
// object into it, and only then injects the refusal, so the route's state at the plant is what
// production has at that hop. What is new here is the ruler's shape, two things 353r1's DAMAGE row
// does not do:
//   - the numerator is ATTRIBUTABLE per shot: each shot's two objects carry the shot's own sequence
//     number in their names (the 353-v1 lesson: same-name fixtures across shots made whole rows of a
//     raw log unreconstructable), and every reading line carries that same number;
//   - the readings are asserted BYTE-EXACT, not counted loosely: the foreign body is a fixed 64-byte
//     constant whose length the case asserts up front, and the surviving name must read back as
//     `size 64 body [<that constant>]` verbatim while the swapped-in object's own name reads absent -
//     i.e. the caller's name is the foreign object's ONLY name. Same input, same bytes out, every run.
//
// WHAT THIS INSTRUMENT CANNOT COVER, stated because a green must not be read wider than the case
// (ticket 353's own GREEN SCOPE, carried over unchanged): the re-check refuses every top-up that
// landed BEFORE it answered. A top-up landing AFTER the re-check answered and BEFORE Deleteat's
// NtOpenFile binds the name is inside the platform, this package has no seam there and must not gain
// one, and no row below claims that window.
package winsec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	// leftStandingShots354r1 is NOT a race denominator. Each shot re-arms a fixed world and drives one
	// production call, so the instrument tells the truth only when the same columns come back shot
	// after shot. Three is this family's quiet-replicate count (winDetShots342r2).
	leftStandingShots354r1 = 3

	// leftStandingForeignBody354r1 is the fixed body of the object the caller never named. The LENGTH
	// is part of the instrument: exactly 64 bytes, asserted at the top of the case, so the surviving
	// name's byte-exact readback (`size 64 body [...]`) cannot drift when the constant is edited.
	leftStandingForeignBody354r1 = "an object the caller never named, standing in the same directory"

	// leftStandingOwnBody354r1 is what the caller's name carries BEFORE the plant. Deliberately a
	// different length from the foreign body, so a surviving name whose body reads the own length can
	// never be confused with the foreign standing.
	leftStandingOwnBody354r1 = "the caller's own object, standing under its own name before the refusal"

	// leftStandingPayload354r1 is the refused write's payload; landing anywhere in the tree is red.
	leftStandingPayload354r1 = "deterministic payload 354r1 - must never land outside the object the caller named"

	// leftStandingInjected354r1 is the injected refusal this leg names as its own, through the same
	// declared seam the 342r2 window uses.
	leftStandingInjected354r1 = "354r1 injected through this package's declared AC#5 seam: the descriptor could not be applied"

	// leftStandingControl354r1 is the positive control's own name: the world where nothing tops up, so
	// the name still carries the object the route opened and the removal MUST still run. It is the leg
	// that sees the OTHER tooth: with the removal pulled the name stays, and this control - not the
	// numerator above - is what goes red.
	leftStandingControl354r1     = "no-top-up-control-354r1.bin"
	leftStandingControlBody354r1 = "the caller's own object, standing under its own name, before a refused write with nothing topped up"
)

// TestAC2WindowsAfterTheRefusalTheCallersNameStandsCarryingTheForeignObject354r1 is ticket 354 AC#2.
//
// The numerator, one shot at a time: after the route answers the injected seal refusal through the
// planted window, read the caller's name back by its own absolute spelling and count the shot ONLY
// when all three hold - the name is PRESENT, its body is the byte-exact foreign pre-state, and the
// swapped-in object's own name is ABSENT (the plant moved that object's only name onto the caller's
// spelling). The gate is exact, not loose: leftStandingForeignName must equal the entered count, and
// the entered count must equal the shot count. Anything else is red with the reading named.
func TestAC2WindowsAfterTheRefusalTheCallersNameStandsCarryingTheForeignObject354r1(t *testing.T) {
	if got := len(leftStandingForeignBody354r1); got != 64 {
		t.Fatalf("354r1 AC#2 RED instrument: the fixed foreign body is %d bytes, want the fixed 64-byte "+
			"pre-state this instrument's byte-exact readback is built on", got)
	}
	if len(leftStandingOwnBody354r1) == len(leftStandingForeignBody354r1) {
		t.Fatalf("354r1 AC#2 RED instrument: the own body and the foreign body are the same length, so a " +
			"surviving name's body could not tell the two objects apart")
	}
	// The world is 353r1's, reused verbatim; only the fixtures are re-named per shot, with the shot's
	// own sequence number in each name so every raw-log line is attributable to its shot.
	w := newWinDetSwapWorld342r2(t)
	// wantStanding is the byte-exact readback the surviving name must produce: readBack342r2 formats
	// present names as `size <n> body [<body>]`, and both halves are fixed constants here.
	wantStanding := fmt.Sprintf("size %d body [%s]",
		len(leftStandingForeignBody354r1), leftStandingForeignBody354r1)
	leftStandingForeignName, entered := 0, 0
	for shot := 1; shot <= leftStandingShots354r1; shot++ {
		w.name = filepath.Join(w.dir, fmt.Sprintf("left-standing-354r1-shot%d.bin", shot))
		w.swapIn = filepath.Join(w.dir, fmt.Sprintf("the-unnamed-object-354r1-shot%d.bin", shot))
		w.restage()
		w.sealRun, w.entered, w.plantErr = false, false, ""
		w.installWindow(t)
		err := PrivateFile(w.name, []byte(leftStandingPayload354r1), 0o600)
		leg := spokeWinDet342r2(err)
		kName, dName := readBack342r2(w.name)
		kSwap, dSwap := readBack342r2(w.swapIn)
		if w.entered {
			entered++
			if kName == "present" && dName == wantStanding && kSwap == "absent" {
				leftStandingForeignName++
			}
		}
		t.Logf("354r1 AC#2 shot %d LEFT-STANDING : the caller's name=%s (%s) | the swapped-in object's own "+
			"name=%s (%s) | %s", shot, kName, dName, kSwap, dSwap, dirListing342r2(w.dir))
		t.Logf("354r1 AC#2 shot %d LIVENESS+REFUSAL: leg=%s | realSealRanFirst=%v windowEntered=%v "+
			"plantErr=%s | isNotSealable=%v | verbatim=%v", shot, leg, w.sealRun, w.entered,
			noteOr342r2(w.plantErr), errors.Is(err, ErrNotSealable), err)
		if !w.entered {
			t.Errorf("354r1 AC#2 RED: shot %d did not ENTER the window (%s), so neither the numerator nor "+
				"the re-check was reachable on this shot - a zero leftStandingForeignName from a world that "+
				"never entered is the baseline's handle shape, NOT a green", shot, noteOr342r2(w.plantErr))
			continue
		}
		if !w.sealRun {
			t.Errorf("354r1 AC#2 RED instrument: shot %d entered the window without the REAL seal having "+
				"run first, so the plant did not happen inside the production route's own window", shot)
		}
		if kSwap != "absent" {
			t.Errorf("354r1 AC#2 RED instrument: shot %d entered the window yet the swapped-in object's own "+
				"name still reads %s (%s), so the plant did not move that object's only name onto the "+
				"caller's spelling and the numerator below is counting nothing", shot, kSwap, dSwap)
		}
		if err == nil {
			t.Errorf("354r1 AC#2 RED: shot %d returned nil after the caller's name had been vacated and "+
				"refilled by an object nobody named, so the route neither refused nor reported what it did", shot)
		} else if !errors.Is(err, ErrNotSealable) {
			t.Errorf("354r1 AC#2 RED: shot %d dropped the sentinel AC#5 promises: %v", shot, err)
		}
		if leg != winDetLegs342r2[10].name {
			t.Errorf("354r1 AC#2 RED: shot %d refused from %s, want the seal leg this instrument injected "+
				"through the declared seam (verbatim: %v)", shot, leg, err)
		}
		if strings.Contains(dName, leftStandingPayload354r1) || strings.Contains(dSwap, leftStandingPayload354r1) {
			t.Errorf("354r1 AC#2 RED: shot %d landed the refused write's payload in the tree (name=%s "+
				"swap-in=%s)", shot, dName, dSwap)
		}
		if kName != "present" {
			t.Errorf("354r1 AC#2 RED: shot %d the caller's name did NOT survive the refusal (reads %s: %s) - "+
				"the identity re-check in unlinkPinnedLeaf is the leg that leaves it standing, and on these "+
				"bytes it no longer does; leftStandingForeignName cannot be counted from this shot", shot, kName, dName)
			continue
		}
		if dName != wantStanding {
			t.Errorf("354r1 AC#2 RED: shot %d the surviving name carries [%s], want the byte-exact foreign "+
				"pre-state [%s] - what stands under the caller's name is not the object the plant moved in",
				shot, dName, wantStanding)
		}
	}
	t.Logf("354r1 AC#2 over %d shots: entered=%d/%d, leftStandingForeignName=%d (the caller's name standing "+
		"after the refusal, carrying the foreign object as its ONLY name)",
		leftStandingShots354r1, entered, leftStandingShots354r1, leftStandingForeignName)
	if entered == leftStandingShots354r1 && leftStandingForeignName != entered {
		t.Errorf("354r1 AC#2 RED: leftStandingForeignName=%d of %d entered worlds - the ratified residual "+
			"moved. Form 乙's price is that the caller's name STAYS under the refused write and stays "+
			"carrying the object the caller never named; a landing that changed either half changed what "+
			"stands under a caller's name without saying so", leftStandingForeignName, entered)
	}
	if entered == leftStandingShots354r1 && leftStandingForeignName == entered {
		t.Logf("354r1 AC#2 GREEN SCOPE, stated so this row is not read wider than the case: the numerator " +
			"says the re-check LEFT the caller's name standing, carrying the foreign object as its only " +
			"name, on every entered shot - the residual ticket 353's owner accepted with form 乙. It says " +
			"NOTHING about a top-up that lands AFTER the re-check answered and BEFORE Deleteat's NtOpenFile " +
			"binds the name: that window is inside the platform, this package has no hook there, and this " +
			"instrument does not cover it and may not be read as covering it.")
	}
}

// TestAC2WindowsRemovalStillRunsWhenNothingTopsUpTheName354r1 is this instrument's positive control,
// and it is load-bearing for the numerator above in BOTH directions: with nothing topped up the
// caller's name still carries the object the route opened, so the re-check matches and the removal
// must still run - the name reads ABSENT afterwards. That proves the ruler can see an absent name at
// all (a numerator that only ever reads "standing" is a blade that bites nothing), and it is the leg
// that answers AC#3's second tooth: pull the removal tooth and THIS control goes red, while the
// numerator above (which measures what the re-check leaves standing) is untouched by that pull.
func TestAC2WindowsRemovalStillRunsWhenNothingTopsUpTheName354r1(t *testing.T) {
	dir := SealableTempDirForTest124(t)
	name := filepath.Join(dir, leftStandingControl354r1)
	if err := os.WriteFile(name, []byte(leftStandingControlBody354r1), 0o600); err != nil {
		t.Fatalf("stage the control's object: %v", err)
	}
	injected := errors.New("354r1 control injected through the same declared seam: the descriptor could not be applied")
	withInjectedSealFailure(t, injected)
	err := PrivateFile(name, []byte(leftStandingPayload354r1), 0o600)
	kind, detail := readBack342r2(name)
	t.Logf("354r1 AC#2 CONTROL (no top-up, the name still carries the object the route opened): "+
		"isNotSealable=%v | leg=%s | name=%s (%s)", errors.Is(err, ErrNotSealable),
		spokeWinDet342r2(err), kind, detail)
	if err == nil || !errors.Is(err, ErrNotSealable) {
		t.Errorf("354r1 AC#2 RED instrument: the control call did not answer with the seal refusal AC#5 "+
			"promises (%v), so this control is not measuring the failure path at all", err)
	}
	if kind != "absent" {
		t.Errorf("354r1 AC#2 RED: with NOTHING topped up the refused write still left the name standing "+
			"(%s: %s), so the refusal-path removal no longer runs - the numerator above would then be "+
			"counting a removal that cannot happen, and AC#5's refused write now leaves a half-made file",
			kind, detail)
	}
}
