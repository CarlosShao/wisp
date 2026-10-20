//go:build windows

// Ticket 353 AC#2, leg 353-r1: the DETERMINISTIC instrument for the refusal-path removal's
// NAME-TOP-UP damage. The DAMAGE column and the ENTRY column are printed and gated as two numerators
// that never merge and never subtract (ticket 348 AC#3's rule, and this ticket's own AC#2 wording).
//
// WHAT THIS MEASURES. `privateFile`'s failure path removes the leaf NAME, relative to the pinned
// parent directory (`Root.Remove(base)`). Whoever can vacate that name and rename a second object
// into it while the route holds the leaf makes that removal delete an object the caller never named:
// the swapped-in object's only name is the caller's spelling, so after the removal it has no name at
// all. `damage` counts exactly that, per shot, read back through the caller's own absolute spelling.
//
// THE WORLD IS NOT INVENTED HERE. It is the world ticket 342-r2 already built through the seam this
// package ALREADY declares (`applyDescriptor`, winsec_windows.go:35-38: "Tests replace it to inject a
// failure (AC#5)"): installWindow runs the REAL seal first, keeps its verdict, and only then vacates
// the caller's name and renames the unnamed object into it, and finally returns the injected refusal -
// so when the route reaches its removal the world has already changed. This case reuses that
// instrument (`newWinDetSwapWorld342r2` / `installWindow`) verbatim: no new production seam, no new
// exported name, no racer, no timing. Ticket 342-r2's bench rows S3/S4 are why `entered` is a number
// and not an assumption: the landing's leaf handle carries FILE_SHARE_DELETE, so the move lands
// (`entered=3/3`); the pre-pinning by-name handle freezes the name (`entered=0/3`).
//
// NO ANCESTOR SWAP IS A PREMISE (AC#2's fourth criterion). Nothing here renames a directory, plants a
// junction or touches containment: only a NAME in the caller's own directory moves. That is the whole
// reason this is a new ticket and not a row of the 342/348 families - those families' numerators
// (`truncated` / `deleted`) read 0 in this world, and borrowing either would be borrowing the other
// ticket's green.
//
// ★WHAT THIS INSTRUMENT CANNOT COVER, stated because a green must not be read wider than the case:
// form 乙 re-checks the name against the object the handle stands on IMMEDIATELY BEFORE the removal,
// so it refuses deterministically every top-up that landed BEFORE that re-check. A top-up that lands
// in the window AFTER the re-check has answered and BEFORE `Deleteat`'s `NtOpenFile` binds the name
// is NOT refused and NOT measurable here: producing that world needs a hook INSIDE the platform's
// `NtOpenFile` (`internal/syscall/windows/at_windows.go:247`), and this package has no seam there and
// must not gain one. The damage case below is therefore a claim about "a top-up that arrived before
// the re-check" and about nothing else, and its green says so out loud. Same shape as ticket 322
// AC#1's "report the window, do not paper over it": a blade that covers one world may not report
// green for another.
package winsec

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	// nameTopUpControl353r1 is the positive control's own name. The control is the world where nothing
	// tops up: the seal is refused with the caller's name still carrying the object the route opened.
	// There the removal MUST still run, or AC#5's "a refused write leaves no half-made file" has been
	// traded away for the top-up guard - and the damage column would be a ruler that cannot see an
	// absent name at all, which is this repository's booked failure mode ("a blade that bites nothing
	// may not report green").
	nameTopUpControl353r1 = "no-top-up-control-353r1.bin"
	// nameTopUpControlBody353r1 is written before the control call so a successful O_TRUNC plus a
	// removal is observable as an ABSENT name rather than as an untouched file.
	nameTopUpControlBody353r1 = "the caller's own object, standing under its own name, before a refused write"
)

// TestAC2WindowsNameTopUpBeforeTheReCheckLeavesTheTopUpObjectNamed353r1 is ticket 353 AC#2 case one.
//
// Two numerators, gated separately:
//   - ENTRY (liveness): `entered=3/3` with `realSealRanFirst=true` in the SAME shot. Zero entered shots
//     is a RED instrument, never a green: it means the platform froze the name and the damage question
//     was never reachable on this tree;
//   - DAMAGE: the number of entered shots in which the swapped-in object was left with NO name. Form
//     乙's promise is that this is 0, because the route re-checks the name against the object its own
//     handle stands on before it removes, and a name carrying somebody else's object fails that check.
//
// The refusal column (AC#5's sentinel and the leg that answered) is asserted too, because a removal
// that was skipped for the wrong reason - a route that never reached its failure path, or a seal that
// refused without the world changing - must not be reported as "the guard held".
//
// The ratified residual, written into the judgement instead of smoothed over: when the re-check
// refuses, the caller's NAME STAYS. In the attacked state that is an observable behaviour change -
// the name survives and now carries the object the caller never named, and `privateFile` still returns
// only the seal failure (reporting "I did not delete this" would be a change to C26's entry contract,
// which is not this landing's). That is the price the machine owner accepted with form 乙, and it is
// what the RESIDUAL line below prints on every shot that the guard refused.
func TestAC2WindowsNameTopUpBeforeTheReCheckLeavesTheTopUpObjectNamed353r1(t *testing.T) {
	w := newWinDetSwapWorld342r2(t)
	entered, frozen, damage, refused := 0, 0, 0, 0
	for shot := 1; shot <= winDetShots342r2; shot++ {
		w.restage()
		w.sealRun, w.entered, w.plantErr = false, false, ""
		w.installWindow(t)
		err := PrivateFile(w.name, []byte(winDetPayload342r2), 0o600)
		leg := spokeWinDet342r2(err)
		kName, dName := readBack342r2(w.name)
		kSwap, dSwap := readBack342r2(w.swapIn)
		if w.entered {
			entered++
			if kName == "absent" {
				// The plant moved the sibling's ONLY name onto the caller's spelling, so an absent name
				// after the refusal means the route's removal destroyed that object.
				damage++
			}
			if kName == "present" {
				refused++
			}
		} else {
			frozen++
		}
		t.Logf("353r1 AC#2 shot %d DAMAGE : the caller's name=%s (%s) | the swapped-in object's own name=%s (%s) | "+
			"swapped-in object left with no name=%v | %s", shot, kName, dName, kSwap, dSwap,
			w.entered && kName == "absent", dirListing342r2(w.dir))
		t.Logf("353r1 AC#2 shot %d LIVENESS+REFUSAL: leg=%s | realSealRanFirst=%v windowEntered=%v plantErr=%s | "+
			"isNotSealable=%v | verbatim=%v", shot, leg, w.sealRun, w.entered, noteOr342r2(w.plantErr),
			errors.Is(err, ErrNotSealable), err)
		if !w.entered {
			// Reported verbatim and then failed on purpose: this is the BASELINE side of the ENTRY pair
			// (bench row S3, the by-name handle freezing the name), and it is the reading 去牙发 B has to
			// produce. Naming it here is what stops a zero damage numerator from being read as a guard.
			t.Errorf("353r1 AC#2 RED: world %d/%d did not ENTER - the platform refused to vacate the caller's "+
				"name while the route held its leaf (%s), so neither the damage question nor form 乙's re-check "+
				"was reachable on this shot. This is the baseline's handle shape, NOT a green.",
				shot, winDetShots342r2, w.plantErr)
			continue
		}
		if !w.sealRun {
			t.Errorf("353r1 AC#2 RED instrument: shot %d entered the window without the REAL seal having run "+
				"first, so the move did not happen inside the production route's own window", shot)
		}
		if kSwap != "absent" {
			t.Errorf("353r1 AC#2 RED instrument: shot %d entered the window yet the swapped-in object's own name "+
				"still reads %s (%s), so the plant did not move its only name onto the caller's spelling and "+
				"the damage column below is measuring nothing", shot, kSwap, dSwap)
		}
		if err == nil {
			t.Errorf("353r1 AC#2 RED: shot %d returned nil after the caller's name had been vacated and refilled "+
				"by an object nobody named, so the route neither refused nor reported what it removed", shot)
		} else if !errors.Is(err, ErrNotSealable) {
			t.Errorf("353r1 AC#2 RED: shot %d dropped the sentinel AC#5 promises: %v", shot, err)
		}
		if leg != winDetLegs342r2[10].name {
			t.Errorf("353r1 AC#2 RED: shot %d refused from %s, want the seal leg this instrument injected "+
				"through the declared seam (verbatim: %v)", shot, leg, err)
		}
		if strings.Contains(dName, winDetPayload342r2) || strings.Contains(dSwap, winDetPayload342r2) {
			t.Errorf("353r1 AC#2 RED: shot %d landed the refused write's payload in the tree (name=%s sibling=%s)",
				shot, dName, dSwap)
		}
	}
	t.Logf("353r1 AC#2 over %d shots: entered=%d/%d, worlds the platform froze=%d, "+
		"unnamed objects left with no name (DAMAGE)=%d, refusals that left the caller's name standing=%d",
		winDetShots342r2, entered, winDetShots342r2, frozen, damage, refused)
	if entered == winDetShots342r2 && damage > 0 {
		t.Errorf("353r1 AC#2 RED: %d of %d entered worlds left the swapped-in object with NO name, so the "+
			"refusal-path removal ran on a name that carried an object the caller never named - ticket 353's "+
			"damage, on the platform where the route's own handle makes that name vacate-able (bench row S4)",
			damage, winDetShots342r2)
	}
	if entered == winDetShots342r2 && damage == 0 {
		t.Logf("353r1 AC#2 GREEN SCOPE, stated so this row is not read wider than the case: the re-check ran " +
			"BEFORE the removal, so it refused every top-up that had already landed - and it left the " +
			"caller's name standing (an observable behaviour change in the attacked state, ratified with " +
			"form 乙). It says NOTHING about a top-up that lands AFTER the re-check answered and BEFORE " +
			"Deleteat's NtOpenFile binds the name: that window is inside the platform, this package has no " +
			"hook there, and this instrument does not cover it and may not be read as covering it.")
	}
}

// TestAC2WindowsRemovalStillRunsWhenNobodyTopsUpTheName353r1 is the positive control, and it is
// load-bearing for the case above: a damage numerator that can only ever read 0 is not evidence.
//
// Same route, same injected seal refusal, no plant at all: the caller's name still carries the object
// the route opened, so form 乙's re-check matches and the removal still runs. That is AC#5's "a refused
// write leaves no half-made file", which this landing must keep on both platforms - and it is proof
// that `readBack342r2` returns "absent" when an object really was removed under this instrument's nose.
func TestAC2WindowsRemovalStillRunsWhenNobodyTopsUpTheName353r1(t *testing.T) {
	dir := SealableTempDirForTest124(t)
	name := filepath.Join(dir, nameTopUpControl353r1)
	if err := os.WriteFile(name, []byte(nameTopUpControlBody353r1), 0o600); err != nil {
		t.Fatalf("stage the control's object: %v", err)
	}
	injected := errors.New("353r1 control injected through the same declared seam: the descriptor could not be applied")
	withInjectedSealFailure(t, injected)
	err := PrivateFile(name, []byte(winDetPayload342r2), 0o600)
	kind, detail := readBack342r2(name)
	t.Logf("353r1 AC#2 CONTROL (no top-up, name still carries the object the route opened): "+
		"isNotSealable=%v | leg=%s | name=%s (%s) | %s", errors.Is(err, ErrNotSealable),
		spokeWinDet342r2(err), kind, detail, dirListing342r2(dir))
	if err == nil || !errors.Is(err, ErrNotSealable) {
		t.Errorf("353r1 AC#2 RED instrument: the control call did not answer with the seal refusal AC#5 "+
			"promises (%v), so this pair is not measuring the failure path at all", err)
	}
	if kind != "absent" {
		t.Errorf("353r1 AC#2 RED instrument: with NOTHING topped up the refused write still left the name "+
			"standing (%s: %s), so the refusal-path removal no longer runs - and the case above's zero "+
			"damage numerator would be a blade that bites nothing, not a guard", kind, detail)
	}
}
