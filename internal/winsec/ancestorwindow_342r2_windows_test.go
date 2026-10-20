//go:build windows

// Tickets 342 AC#1 + 348 AC#2, leg 342-r2: the DETERMINISTIC Windows window instruments. The DAMAGE
// column and the REFUSAL column are printed and gated as two numerators that never merge.
//
// WHAT THE ADJUDICATIONS ASKED FOR, in their own words. 342-v2: the landing side still leaks 0~2/400
// `truncated`, and the fix it named is a deterministic instrument for "an ancestor that flips only after
// the open", not a racer. 348-v1: "修前分子>0 在 Windows 取不到" = 缺世界, because the steady state is
// answered by the installed C26 resolver and by preOpenLeafGuard before the route ever reaches its
// removal.
//
// THIS LEG DID NOT INVENT A SEAM TO GET THE WORLD IT WANTED. It measured first, outside the repository
// (bench probes 342r2 probe1..probe6 plus probe8; verbatim rows archived under
// .scratch/wisp/probes/342/r2/logs/bench-probes-verbatim.txt - probe7 never compiled and is
// void-registered in that same file, and nothing in these instruments leans on it),
// which name an attacker can move at each instant the route owns, and it then drives only the worlds the
// platform actually offers, through the seam this package ALREADY declares:
//
//	var applyDescriptor = applyDescriptorWindows   // winsec_windows.go:35-38: "Tests replace it to
//	                                               // inject a failure (AC#5)"
//
// private_fail_test.go:14 replaces it, and sealhandle_window_120_other_test.go:163-178
// (inIdentityWindow120r3, the POSIX counterpart the dispatch names) already uses the same var as a
// WINDOW rather than a substitution: run the real leg, keep its verdict, then make the move an attacker
// would have made. The plant here runs the REAL seal through the REAL platform implementation first, so
// the route's state at the plant is what production has at that hop. No production name is added, no
// function shape changes, and winsec.go:430's wrapPath is untouched.
//
// THE MEASUREMENT, in the route's own handle states. The sharing mask is the mechanism:
//
//	BASELINE open = os.OpenFile  -> syscall.Open      -> share READ|WRITE          (no FILE_SHARE_DELETE)
//	LANDING  open = Root.OpenFile -> internal Openat  -> share READ|WRITE|DELETE
//
//	S1  BASELINE, at its removal (leaf closed, nothing pinned): rename an ANCESTOR    = ok
//	S2  LANDING,  at its removal (parent pinned, leaf closed):   rename an ANCESTOR    = Access is denied
//	S2b LANDING,  at its removal:                                rename the PINNED dir = in use, another process
//	S3  BASELINE, at the seal hop (leaf open by name):           vacate the leaf name  = in use; rename into it = Access is denied
//	S4  LANDING,  at the seal hop (leaf open through the pin):   vacate the leaf name  = POSSIBLE; rename into it = POSSIBLE
//	S5  positive control, nothing open inside the tree:          every one of those moves lands.
//
// So dispatch ①'s sentence - "父钉已拿到、祖先在 open 那一跳之后才翻" - splits into measured answers, and
// this file states both instead of choosing the convenient one:
//
//   - an ANCESTOR cannot be flipped after the open on the landing at all: rows S2/S2b say the platform
//     refuses the move while the pin's handle stands inside it, at both the seal hop and the removal hop.
//     Making that flip land would need a seam INSIDE the route, between pinnedParent's descent and the
//     acts that follow it - a new production seam, which this dispatch forbids me to add. It is reported
//     in case THREE's rows instead: that is the second human approval, and it is the same wall 348-v1
//     hit from the steady-state side;
//   - what the pin DOES open, measurably and deterministically, is the NAME: the landing's handle carries
//     FILE_SHARE_DELETE, so the caller's leaf name can be vacated and re-filled with an object the caller
//     never named while the route holds it (S4), which the baseline's own handle forbids (S3). The route's
//     refusal-path removal is relative to the pinned DIRECTORY, and a directory-relative removal of a
//     name that now carries a foreign object deletes that object's only remaining name. That is a residual
//     of THIS landing, found by this leg, and case ONE measures it as a number on both trees.
//
// CASE TWO is the steady-state ancestor junction - deterministic on both trees with no seam at all - kept
// in two columns: 348-v1 read its refusal and reported "witness intact" as an adjective. This case
// attaches numbers, and it attaches a POSITIVE CONTROL proving the damage column is a ruler that can see
// damage at all: a blade that bites nothing may not report green, which is this repository's own booked
// failure mode.
//
// WHAT NONE OF THE THREE CASES CLAIMS: none of them makes openPinnedLeaf's identity gate speak (its
// window is between the guard's verdict and the open, and there is no seam there); none reproduces
// ticket 342's `truncated` numerator (while the landing's handle is open, the truncation has already gone
// to the object the HANDLE names, not to whatever the name later carries - probe rows: a write through
// the handle lands in the pinned object, not in the swapped-in one); and nothing here may be read as "the
// window is closed" on either platform. winsec.go:118-140's statement is not this leg's to edit: both
// platforms measured is the gate for that sentence, and the update right belongs after the adjudication.
package winsec

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	// winDetShots342r2 is NOT a race denominator. Every case re-arms its own world and drives one
	// production call, so the instrument tells the truth only when the same columns come back shot after
	// shot. Three is this family's quiet-replicate count.
	winDetShots342r2 = 3

	winDetLeaf342r2       = "leaf-342r2w.bin"
	winDetSwapIn342r2     = "the-unnamed-object-342r2.bin"
	winDetOwnBody342r2    = "the caller's own object under its own name"
	winDetSwapBody342r2   = "an object the caller never named, standing in the same directory"
	winDetPayload342r2    = "deterministic payload 342r2 - must never land outside the object the caller named"
	winDetInjected342r2   = "injected through this package's declared AC#5 seam: the descriptor could not be applied"
	winDetAncData342r2    = "data"
	winDetAncSub342r2     = "sub"
	winDetAncLeaf342r2    = "leaf-342r2q.bin"
	winDetAncPark342r2    = "parked-real-data-342r2q"
	winDetAncForeign342r2 = "someone-elses-tree-342r2q"
	winDetAncControl342r2 = "the-control-342r2q.bin"
	winDetStageLink342r2  = "staged-anc-junction-342r2q"
	// Bodies differ in LENGTH from the caller's own on purpose: an O_TRUNC that lands and writes nothing
	// leaves a short body, and equal-length bodies let damage read as untouched.
	winDetWitness342r2     = "foreign bytes are not yours to touch, and this body is longer on purpose"
	winDetControlBody342r2 = "the control witness exists so the damage column can prove it can see damage at all"
)

// winDetLegs342r2 is THIS instrument's own attribution roster, declared in order so two runs diff line
// by line. Every `why` string is copied from the production text it names, never paraphrased: an
// attribution that does not match the bytes a leg formats is a leg that never spoke. It is deliberately
// local - the r1 racer's roster belongs to that instrument, and editing another leg's roster so these
// readings classify would be the same act as loosening an assertion.
var winDetLegs342r2 = []struct{ why, name string }{
	{"path escapes from parent", "os.Root containment, inside the descent (act time)"},
	{"the name carries a link, so the pinned object", "pinnedParent sameObjectPinnedWindows, link face (BEFORE the open)"},
	{"the name carries a different object than the one that was pinned", "pinnedParent sameObjectPinnedWindows, SameFile face (BEFORE the open)"},
	{"the name carries a link, so this handle", "openPinnedLeaf re-check, link face (AFTER the open)"},
	{"this handle is not the object the name carries now", "openPinnedLeaf re-check, SameFile face (AFTER the open)"},
	{"the name carries a reparse point", "preOpenLeafGuard, leaf refusal (BEFORE the act, ticket 322)"},
	{"the reparse-point probe could not read the name", "preOpenLeafGuard, unreadable-name refusal"},
	{"reaches it through the link at", "RemoveUnlinked's firstLinkAncestor (the by-name reclaim guard)"},
	{"traverses a reparse point at", "the Windows placement floor's own reparse walk (placement_windows.go)"},
	{"path traverses a reparse point", "the installed C26 resolver's steady-state refusal (internal/risk)"},
	{"cannot apply a private security descriptor", "seal refusal on this route (case one injects it through applyDescriptor and names it as injected; case two never reaches this leg)"},
	{"not provably resolved", "resolveString refused the spelling"},
}

const (
	winDetOtherLeg342r2 = "some other failure"
	winDetNilLeg342r2   = "none: the call returned nil"
)

func spokeWinDet342r2(err error) string {
	if err == nil {
		return winDetNilLeg342r2
	}
	msg := err.Error()
	for _, leg := range winDetLegs342r2 {
		if strings.Contains(msg, leg.why) {
			return leg.name
		}
	}
	return winDetOtherLeg342r2
}

// ---- read-back: always by absolute spelling, never through a name that may have been re-pointed --------

func readBack342r2(path string) (kind, detail string) {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "absent", err.Error()
		}
		return "unreadable", err.Error()
	}
	if info.Mode()&fs.ModeSymlink != 0 || reparseAt(path) {
		return "link", "mode " + info.Mode().String()
	}
	body, rErr := os.ReadFile(path)
	if rErr != nil {
		return "unreadable", rErr.Error()
	}
	return "present", fmt.Sprintf("size %d body [%s]", info.Size(), body)
}

func dirListing342r2(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "readdir-error: " + err.Error()
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return "dir entries [" + strings.Join(names, " | ") + "]"
}

func sameFileViaSpelling342r2(viaSpelling, ownSpelling string) bool {
	a, errA := os.Stat(viaSpelling)
	b, errB := os.Stat(ownSpelling)
	if errA != nil || errB != nil {
		return false
	}
	return os.SameFile(a, b)
}

func noteOr342r2(s string) string {
	if s == "" {
		return "<none>"
	}
	return s
}

// winExec342r2 runs a helper and hands its output back as text plus an error, so a construction failure
// becomes a named reading. mustExec ends the run, and execNoCheck belongs to the r1 racer's file; this
// leg's helper is its own so a tree that does not carry that file still compiles this one.
func winExec342r2(name string, args ...string) (string, error) {
	var out bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout, cmd.Stderr = &out, &out
	err := cmd.Run()
	return out.String(), err
}

// mkJunction342r2 is this platform's no-privilege construction (`mklink /J`). It returns an error rather
// than skipping, and callers treat a failure as a bench fault: a junction needs neither elevation nor
// developer mode here, so a fixture that quietly skipped is how a zero damage numerator would come to be
// read as a green.
func mkJunction342r2(link, target string) error {
	out, err := winExec342r2("cmd", "/c", "mklink", "/J", link, target)
	if err != nil {
		return fmt.Errorf("mklink /J %s: %v (%s)", filepath.Base(link), err, strings.TrimSpace(out))
	}
	if !reparseAt(link) {
		return fmt.Errorf("%s is not a reparse point after mklink, so this fixture built nothing", link)
	}
	return nil
}

// winClearUnlinked342r2 removes an entry WITHOUT following it, tolerating absence, through this package's
// own promise (RemoveUnlinked). A51② is the case where the naive call resolves the target and then cites
// its contents as the reason not to delete; leaving such an entry standing is what makes a temp root
// untraversable afterwards.
func winClearUnlinked342r2(path string) error {
	if _, err := os.Lstat(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return RemoveUnlinked(path)
}

// ---- case one: the caller's NAME swapped after the seal's verdict -------------------------------------

// winDetSwapWorld342r2 is the caller's directory with two objects: the name the caller spells, and a bare
// sibling in the SAME directory that the caller never named - dispatch ③'s answerable form: nothing
// escapes, no containment rule can answer, so only this package's own acts can. Rows S3/S4 decide what
// each tree can do with that shape: the landing's handle shares DELETE, so the name can be vacated and
// re-filled while the route holds it; the baseline's handle does not, so the platform freezes the name
// and the world is not enterable there at all.
type winDetSwapWorld342r2 struct {
	t        *testing.T
	dir      string
	name     string // the caller's spelling; object A before the plant, object B's only name after it
	swapIn   string // object B's own name, consumed by the plant's rename
	sealRun  bool   // the REAL seal ran first, inside the seam
	entered  bool   // the plant landed inside the route's window
	plantErr string // verbatim refusal when the platform would not let the move happen
}

func newWinDetSwapWorld342r2(t *testing.T) *winDetSwapWorld342r2 {
	t.Helper()
	w := &winDetSwapWorld342r2{t: t, dir: SealableTempDirForTest124(t)}
	w.name = filepath.Join(w.dir, winDetLeaf342r2)
	w.swapIn = filepath.Join(w.dir, winDetSwapIn342r2)
	w.restage()
	return w
}

// restage puts both objects back with their own bodies. The plant consumes one of them by rename, so the
// world has to be re-armed before every shot exactly as the r1 family's stageFlip re-arms its staged link
// after every cycle.
func (w *winDetSwapWorld342r2) restage() {
	for p, body := range map[string]string{w.name: winDetOwnBody342r2, w.swapIn: winDetSwapBody342r2} {
		if info, err := os.Lstat(p); err == nil && (info.Mode()&fs.ModeSymlink != 0 || reparseAt(p)) {
			if rErr := winClearUnlinked342r2(p); rErr != nil {
				w.t.Fatalf("clear %s before restaging: %v", filepath.Base(p), rErr)
			}
		}
		if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
			w.t.Fatalf("stage %s: %v", filepath.Base(p), err)
		}
	}
}

// installWindow is the use this package already makes of its own seam: the real applyDescriptorWindows
// runs first and keeps its verdict, the move happens only after that verdict, and the seam then refuses,
// which puts the route on its removal path with the world already changed. The plant's own failure is
// captured and reported instead of ending the run, because on the pre-landing bytes the platform refusing
// the move IS the reading (row S3).
func (w *winDetSwapWorld342r2) installWindow(t *testing.T) {
	t.Helper()
	orig := applyDescriptor
	applyDescriptor = func(path string, dir bool) error {
		err := orig(path, dir)
		w.sealRun = true
		if err != nil {
			return err // the seal refused on its own: no world change, keep that answer
		}
		if rErr := os.Remove(w.name); rErr != nil {
			w.plantErr = "vacate the caller's name: " + rErr.Error()
			return err
		}
		if rErr := os.Rename(w.swapIn, w.name); rErr != nil {
			w.plantErr = "rename the unnamed object into the caller's name: " + rErr.Error()
			return err
		}
		w.entered = true
		return errors.New(winDetInjected342r2)
	}
	t.Cleanup(func() { applyDescriptor = orig })
}

// TestAC1WindowsNameSwappedAfterTheSealVerdictReportsDamageAndRefusalAsTwoNumbers342r2 is case one.
//
// The REFUSAL column answers: which leg answered, and is it the one this instrument injected? The DAMAGE
// column answers a different question with a number: after the refusal, is the object the caller never
// named still reachable under any name? They are two rows gated on separate conditions, because merging
// them is exactly how one ticket would borrow the other's green (ticket 348 AC#3's rule).
//
// What is ASSERTED, and both are contract statements this package already makes:
//   - a refused seal answers with ErrNotSealable (AC#5; sealError is what wraps it);
//   - "a refused write leaves no half-made file": the name is gone after the refusal. ★Ticket 353 form 乙
//     narrowed the second sentence to what it always meant, and says so here instead of loosening it: no
//     half-made file THIS ROUTE MADE. Once the name has been vacated and re-filled with an object the
//     caller never named, removing that name would destroy somebody else's only name - which is ticket
//     353's damage - so the route re-checks the name against the object its own handle stands on (the
//     landing's removal, in unlinkPinnedLeaf) and refuses the removal when they differ. In THIS case's
//     world the name is therefore expected to SURVIVE, and it survives carrying the swapped-in body; the
//     un-topped-up world (this file's own positive control lives in nametopup_353_windows_test.go) is
//     where the name must still be gone.
//
// What is REPORTED and not asserted: that removing that name destroyed an object the caller never named.
// Before ticket 353's landing that is what happened on every deterministic shot, and the rows below still
// print it as a number, because a case whose red branch can no longer fire has no way to see the damage
// come back. ★The residual ticket 342-r2 booked here is now GUARDED for every top-up that arrived before
// the re-check answered, and ★NOT guarded for one arriving after it and before Deleteat binds the name -
// that window is inside the platform (NtOpenFile), this package has no seam there, and no row in this
// file may be read as covering it (ticket 322 AC#1's "report the window, do not paper over it").
func TestAC1WindowsNameSwappedAfterTheSealVerdictReportsDamageAndRefusalAsTwoNumbers342r2(t *testing.T) {
	w := newWinDetSwapWorld342r2(t)
	entered, frozen, damageDestroyed, guardHeld := 0, 0, 0, 0
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
				// The plant moved B's only name onto the caller's spelling, so a removed name means the
				// route's removal destroyed the object the caller never named.
				damageDestroyed++
			}
		} else {
			frozen++
		}
		t.Logf("AC#1-one shot %d DAMAGE : the caller's name=%s (%s) | the unnamed object's own name=%s (%s) | "+
			"the swapped-in object left with no name=%v | %s", shot, kName, dName, kSwap, dSwap,
			w.entered && kName == "absent", dirListing342r2(w.dir))
		t.Logf("AC#1-one shot %d REFUSAL: leg=%s | realSealRanFirst=%v windowEntered=%v plantErr=%s | "+
			"isNotSealable=%v | verbatim=%v", shot, leg, w.sealRun, w.entered, noteOr342r2(w.plantErr),
			errors.Is(err, ErrNotSealable), err)
		if !w.entered {
			t.Logf("AC#1-one shot %d NOT AN ENTERED WORLD: the platform refused the move while the route held "+
				"its by-name handle (%s). That is bench row S3: this tree's state, reported, and NOT counted "+
				"as a green and NOT as a red.", shot, w.plantErr)
			continue
		}
		if !w.sealRun {
			t.Errorf("shot %d RED instrument: the real seal never ran before the plant, so this case did not "+
				"enter the window through the production route", shot)
		}
		if err == nil {
			t.Errorf("shot %d RED: the route returned nil after the caller's name had been moved out from under "+
				"it and replaced by an object nobody named, so it neither refused nor reported what it did", shot)
			continue
		}
		if !errors.Is(err, ErrNotSealable) {
			t.Errorf("shot %d RED: the seal-refusal answer dropped the sentinel AC#5 promises: %v", shot, err)
		}
		if leg != winDetLegs342r2[10].name {
			t.Errorf("shot %d RED: the refusal came from %s, want the seal leg this instrument injected "+
				"through the declared seam (verbatim: %v)", shot, leg, err)
		}
		if kName == "present" {
			// ★Ticket 353 form 乙 restated this branch, and it is a NARROWER promise here, not a relaxed
			// one. A surviving name is no longer read on its own: a name surviving with the swapped-in
			// body inside it is the route's own identity re-check refusing to remove a name that had come
			// to carry an object the caller never named (that is ticket 353 being closed, and it is the
			// reading this case now produces on the landed bytes). What stays RED is a name surviving with
			// anything ELSE inside it: that is the removal failing to run for a reason this instrument
			// never authorised.
			if !strings.Contains(dName, winDetSwapBody342r2) {
				t.Errorf("shot %d RED: AC#5's half-made file survived the refusal carrying a body that is "+
					"neither the caller's nor the one the plant moved in (%s) - so the refusal-path removal "+
					"did not run, and nothing on this route explains why", shot, dName)
			}
			guardHeld++
		}
		if strings.Contains(dName, winDetPayload342r2) || strings.Contains(dSwap, winDetPayload342r2) {
			t.Errorf("shot %d RED: the refused write's payload landed in the tree (name=%s sibling=%s)", shot, dName, dSwap)
		}
	}
	t.Logf("AC#1-one over %d shots: worlds entered=%d, worlds the platform refused=%d, "+
		"unnamed objects left with no name=%d, tops-up refused by the identity re-check (name kept)=%d",
		winDetShots342r2, entered, frozen, damageDestroyed, guardHeld)
	if entered == winDetShots342r2 && damageDestroyed == winDetShots342r2 {
		// ★This branch is the PRE-乙 reading, and it is kept on purpose rather than deleted: a case whose
		// red branch can no longer fire cannot see the damage come back. On these bytes it is a REGRESSION
		// reading, not a residual this landing accepts - ticket 353 booked exactly this number, form 乙 was
		// chosen to make it 0, and nametopup_353_windows_test.go is the instrument that now owns the claim.
		t.Errorf("AC#1-one RED: %d of %d entered worlds left the swapped-in object with NO name - the "+
			"refusal-path removal ran on a name the plant had re-filled with an object the caller never "+
			"named (bench row S4: the landing's handle shares FILE_SHARE_DELETE and the baseline's does "+
			"not). That is ticket 353's damage; the identity re-check in unlinkPinnedLeaf is what refuses "+
			"it, and it refuses every top-up that arrived BEFORE it answered - NOT one arriving after it "+
			"and before Deleteat binds the name, which no instrument here can reach.",
			damageDestroyed, winDetShots342r2)
	}
	if entered == winDetShots342r2 && damageDestroyed == 0 && guardHeld == winDetShots342r2 {
		t.Logf("AC#1-one GUARD HELD, stated with its limit: on these bytes the caller's name survived all %d "+
			"entered shots and it survived CARRYING the swapped-in body, i.e. the route refused to remove a "+
			"name whose object it could not prove it owned - ticket 353's damage numerator is 0. The price, "+
			"ratified with form 乙: in the attacked state a name now stays standing that used to be removed, "+
			"and it is the attacker's object standing there. NOT covered and NOT claimed: a top-up landing "+
			"after the re-check answered and before Deleteat's NtOpenFile binds the name.", guardHeld)
	}
	if entered == 0 {
		t.Logf("AC#1-one: the world was not enterable on this tree (%d/%d refused by the platform), so the pair "+
			"for this case is the other tree's reading and these rows state only their own side",
			frozen, winDetShots342r2)
	}
}

// ---- case two: the steady-state ancestor junction, damage and refusal, with a sensitivity control ------

// winDetForeign342r2 names the objects inside somebody else's tree. Both are read ONLY through these
// absolute spellings: after the flip the caller's spelling reaches them too, and a damage count taken
// through a name that can be re-pointed is not a damage count.
type winDetForeign342r2 struct {
	dir     string // <foreign>\sub: where the caller's flipped spelling lands
	witness string // what the caller's spelling resolves to after the flip, read by this own spelling
	control string // the positive control's own name, never handed to the route
}

// winDetAncWorld342r2 is <base>\data\sub\<leaf> with the foreign tree a sibling of `data`, so `data` can
// become a junction at it and the caller's spelling then resolves inside that other tree. No seam, no
// injection, no timing: the state every steady-state refusal is about. The real `data` is PARKED, not
// deleted - rename keeps its contents, and the parked spelling is how the caller's own object is read
// back through a spelling the flip cannot re-point.
type winDetAncWorld342r2 struct {
	t       *testing.T
	base    string
	data    string
	park    string
	name    string
	own     string // the caller's object, read through the parked spelling
	foreign string
	f       winDetForeign342r2
}

func newWinDetAncWorld342r2(t *testing.T) *winDetAncWorld342r2 {
	t.Helper()
	w := &winDetAncWorld342r2{t: t, base: SealableTempDirForTest124(t)}
	w.data = filepath.Join(w.base, winDetAncData342r2)
	w.park = filepath.Join(w.base, winDetAncPark342r2)
	w.name = filepath.Join(w.data, winDetAncSub342r2, winDetAncLeaf342r2)
	w.own = filepath.Join(w.park, winDetAncSub342r2, winDetAncLeaf342r2)
	w.foreign = filepath.Join(w.base, winDetAncForeign342r2)
	w.f.dir = filepath.Join(w.foreign, winDetAncSub342r2)
	w.f.witness = filepath.Join(w.f.dir, winDetAncLeaf342r2)
	w.f.control = filepath.Join(w.f.dir, winDetAncControl342r2)
	for _, d := range []string{filepath.Dir(w.name), w.f.dir} {
		if err := os.MkdirAll(d, 0o700); err != nil {
			t.Fatalf("stage %s: %v", d, err)
		}
	}
	t.Cleanup(w.settle)
	return w
}

// flip parks the real `data` and lands a junction at its name. rename(2) cannot replace a non-empty
// directory, which is why the real ancestor is MOVED OUT first - the same construction the POSIX
// deterministic case uses (ancestorunlink_348_unix_test.go:83's park), for the same reason.
func (w *winDetAncWorld342r2) flip() error {
	stage := filepath.Join(w.base, winDetStageLink342r2)
	if err := winClearUnlinked342r2(stage); err != nil {
		return err
	}
	if err := mkJunction342r2(stage, w.foreign); err != nil {
		return err
	}
	if _, err := os.Lstat(w.park); err == nil {
		return errors.New("the parked ancestor is still standing, so the world was not settled")
	}
	if err := os.Rename(w.data, w.park); err != nil {
		return fmt.Errorf("move the real ancestor out: %v", err)
	}
	if err := os.Rename(stage, w.data); err != nil {
		return fmt.Errorf("land the junction at the ancestor's name: %v", err)
	}
	if !reparseAt(w.data) {
		return errors.New("the ancestor is not a reparse point after the flip")
	}
	return nil
}

// settle restores the real ancestor and re-stocks every object with its own body, so each shot starts
// from a known state and nothing is left standing that a temp-root cleanup would have to traverse.
func (w *winDetAncWorld342r2) settle() {
	if reparseAt(w.data) {
		if err := winClearUnlinked342r2(w.data); err != nil {
			w.t.Logf("342r2 cleanup: clear the flipped ancestor: %v", err)
		}
	}
	if _, err := os.Lstat(w.park); err == nil {
		if _, dErr := os.Lstat(w.data); dErr == nil {
			w.t.Logf("342r2 cleanup: the ancestor name is occupied by something else, leaving the parked tree "+
				"where it is: %s", dirListing342r2(w.base))
		} else if rErr := os.Rename(w.park, w.data); rErr != nil {
			w.t.Logf("342r2 cleanup: rename the real ancestor back: %v", rErr)
		}
	}
	if _, err := os.Lstat(filepath.Dir(w.name)); err != nil {
		if err := os.MkdirAll(filepath.Dir(w.name), 0o700); err != nil {
			w.t.Logf("342r2 cleanup: re-create the chain: %v", err)
		}
	}
	for p, body := range map[string]string{
		w.name:      winDetOwnBody342r2,
		w.f.witness: winDetWitness342r2,
		w.f.control: winDetControlBody342r2,
	} {
		if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
			w.t.Logf("342r2 cleanup: restock %s: %v", filepath.Base(p), err)
		}
	}
}

// TestAC2WindowsSteadyStateAncestorJunctionReportsDamageAndRefusalAsTwoNumbers342r2 is case two. The
// REFUSAL column is what 348-v1 measured (verbatim, identical on both trees: the installed C26 resolver
// answers before any act). The DAMAGE column is what it left as an adjective, and it carries a positive
// control - the same shape of act through a different name, which MUST be able to destroy a witness;
// without that row, "no damage" is a ruler pointing at nothing.
//
// It also states its own limit in the same reading, because that limit is the ticket's finding: this is
// the PRE-ACT face. It is not ticket 348's pre-landing numerator, which needs the ancestor to move AFTER
// the pin is held, and bench rows S2/S2b measure the platform refusing exactly that move while the
// landing's handle stands inside the tree.
func TestAC2WindowsSteadyStateAncestorJunctionReportsDamageAndRefusalAsTwoNumbers342r2(t *testing.T) {
	w := newWinDetAncWorld342r2(t)
	refusalLegs := map[string]int{}
	for shot := 1; shot <= winDetShots342r2; shot++ {
		w.settle()
		if err := w.flip(); err != nil {
			t.Fatalf("shot %d: the steady-state world is a construction and not a race, so a failure here is "+
				"a bench fault: %v", shot, err)
		}

		// ---- the damage column's positive control, on the control's OWN name, before the route runs. ----
		if err := os.WriteFile(w.f.control, []byte(winDetControlBody342r2), 0o600); err != nil {
			t.Fatalf("shot %d stage the control: %v", shot, err)
		}
		viaSpelling := filepath.Join(w.data, winDetAncSub342r2, winDetAncControl342r2)
		same := sameFileViaSpelling342r2(viaSpelling, w.f.control)
		remErr := os.Remove(viaSpelling)
		kCtrl, dCtrl := readBack342r2(w.f.control)
		t.Logf("AC#2-two shot %d CONTROL : the flipped spelling and the foreign spelling name one object=%v | "+
			"by-name remove err=%v | the foreign control name is now=%s (%s)", shot, same, remErr, kCtrl, dCtrl)
		if !same || kCtrl != "absent" {
			t.Errorf("shot %d RED INSTRUMENT: the damage column's positive control could NOT destroy its "+
				"witness through the flipped ancestor (same object=%v, remove err=%v, control now %s), so the "+
				"no-damage row below would be a blade that bites nothing reporting green", shot, same, remErr, kCtrl)
		}

		// ---- the route, on the same flipped world. ----
		err := PrivateFile(w.name, []byte(winDetPayload342r2), 0o600)
		leg := spokeWinDet342r2(err)
		refusalLegs[leg]++
		kWit, dWit := readBack342r2(w.f.witness)
		kOwn, dOwn := readBack342r2(w.own)
		t.Logf("AC#2-two shot %d DAMAGE  : the foreign object the caller's spelling now reaches=%s (%s) | the "+
			"caller's own object, read through the PARKED spelling=%s (%s) | %s", shot, kWit, dWit, kOwn, dOwn,
			dirListing342r2(w.f.dir))
		t.Logf("AC#2-two shot %d REFUSAL: leg=%s | refused=%v | isUnresolved=%v | isNotSealable=%v | verbatim=%v",
			shot, leg, err != nil, errors.Is(err, ErrUnresolvedPath), errors.Is(err, ErrNotSealable), err)
		if err == nil {
			t.Errorf("shot %d RED: the route reported success through an ancestor that is a junction, so its "+
				"act landed in a tree the caller never named", shot)
		}
		if kWit != "present" {
			t.Errorf("shot %d RED (DAMAGE column): the foreign object reads %s after the route ran, so an act "+
				"reached it whatever the refusal said", shot, kWit)
		}
		if strings.Contains(dWit, winDetPayload342r2) {
			t.Errorf("shot %d RED (DAMAGE column): the refused write's payload is readable in the foreign "+
				"object: %s", shot, dWit)
		}
		if kOwn == "absent" {
			t.Errorf("shot %d RED (DAMAGE column): the caller's own object is gone from the parked spelling "+
				"(%s), so a refused act still damaged the tree it stood on", shot, dOwn)
		}
		w.settle()
	}
	for _, l := range winDetLegs342r2 {
		if n := refusalLegs[l.name]; n > 0 {
			t.Logf("AC#2-two': leg that answered = %s : %d shots", l.name, n)
		}
	}
	if n := refusalLegs[winDetOtherLeg342r2]; n > 0 {
		t.Errorf("AC#2-two RED: %d shots answered with a leg this roster does not know, so the refusal column "+
			"is not attributable", n)
	}
	t.Logf("AC#2-two limit, stated instead of glossed: both columns read zero damage with the refusal coming " +
		"from a leg that answers BEFORE any act (the resolver / the floor). That reproduces 348-v1's " +
		"缺世界 deterministically, and it holds on the pre-landing shape too: no damage numerator can be " +
		"produced for a steady-state ancestor link on this platform, and this case is NOT ticket 348's " +
		"pre-landing numerator.")
}

// ---- case three: the platform property the ancestor leg depends on, measured in both handle states -----

// TestAC1WindowsAncestorMoveIsRefusedWhileAParentHandleIsHeld342r2 is dispatch ①'s residual statement
// made as a measurement instead of an argument. The world the two tickets want ("the parent pin is already
// held and an ancestor flips after the open hop") is not enterable on this platform through any seam this
// package declares, and the reason is a property of the handles, not of the route: Windows refuses to
// rename a directory while a handle stands inside it. Two states an instrument CAN build deterministically
// with nothing but the platform's own primitives - so both trees reach them, and neither one is a
// simulation of the other:
//
//	B = the pre-landing route's state at its removal (the leaf was opened by NAME, nothing is pinned):
//	    renaming an ancestor LANDS (row S1), which is what lets a by-name removal re-resolve into the flip;
//	L = the landing's state at the same hop (the parent pinned by handle, the leaf closed):
//	    the same rename is REFUSED (row S2), so there is no window for the removal to be re-pointed into.
//
// The asymmetry is the mechanism form 乙 claims on this platform, and it is measured here rather than
// asserted. It is NOT a claim that the production route was driven: no entry point runs in this case, and
// nothing here may be read as ticket 348's pre-landing numerator on the real route - producing that would
// need a plant point between the leaf handle's close and the removal, and this package has no seam there.
// That seam (a hook inside privateFile's act sequence, or a function-shape change) is a SECOND human
// approval; it is named in this leg's report for the orchestrator to put to the machine owner, and this
// leg does not add it.
func TestAC1WindowsAncestorMoveIsRefusedWhileAParentHandleIsHeld342r2(t *testing.T) {
	base := SealableTempDirForTest124(t)
	data := filepath.Join(base, winDetAncData342r2)
	sub := filepath.Join(data, winDetAncSub342r2)
	leaf := filepath.Join(sub, winDetAncLeaf342r2)
	park := filepath.Join(base, winDetAncPark342r2)
	if err := os.MkdirAll(sub, 0o700); err != nil {
		t.Fatalf("stage the chain: %v", err)
	}
	if err := os.WriteFile(leaf, []byte(winDetOwnBody342r2), 0o600); err != nil {
		t.Fatalf("stage the leaf: %v", err)
	}

	// State B: the pre-landing route's state at the removal hop - the leaf was opened by name and its
	// handle is closed by the time the refusal-path removal runs, and nothing stands on the ancestors.
	f, err := os.OpenFile(leaf, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		t.Fatalf("open the leaf the baseline way: %v", err)
	}
	_ = f.Close()
	moveB := os.Rename(data, park)
	t.Logf("AC#1-three B (the baseline's handle state at its removal hop): rename(an ancestor) err=%v", moveB)
	if moveB != nil {
		t.Errorf("B RED INSTRUMENT: the baseline's own handle state was measured (bench row S1) to ALLOW the "+
			"ancestor move and it did not (%v), so the pair below proves nothing - re-measure the bench", moveB)
	} else {
		if rErr := os.Rename(park, data); rErr != nil {
			t.Fatalf("put the chain back after B: %v", rErr)
		}
	}

	// State L: the landing's state at the same hop. Open the leaf through a pinned parent handle exactly
	// as pinnedParent + openPinnedLeaf do, keep the pin open (the route's deferred close has not run),
	// close the leaf - and try the same move.
	pin, err := os.OpenRoot(sub)
	if err != nil {
		t.Fatalf("pin the parent the landing way (os.OpenRoot on the leaf's directory): %v", err)
	}
	t.Cleanup(func() { _ = pin.Close() })
	lf, err := pin.OpenFile(filepath.Base(leaf), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		t.Fatalf("open the leaf through the pinned handle: %v", err)
	}
	_ = lf.Close()
	moveL := os.Rename(data, park)
	kLeaf, dLeaf := readBack342r2(leaf)
	t.Logf("AC#1-three L (the landing's handle state): rename(an ancestor) err=%v | the chain still intact: "+
		"%s (%s) | %s", moveL, kLeaf, dLeaf, dirListing342r2(data))
	if moveL == nil {
		t.Errorf("L RED: the ancestor move LANDED while a handle was pinned on a directory inside it, so the "+
			"platform does not refuse what form 乙 relies on it refusing (%v) - the window ticket 348 booked "+
			"is open on this platform after all, and the pin is not the lock the landing's comment claims", moveL)
		if _, err := os.Lstat(park); err == nil {
			if rErr := os.Rename(park, data); rErr != nil {
				t.Errorf("restore the chain after L: %v", rErr)
			}
		}
		return
	}
	if _, err := os.Lstat(park); err == nil {
		t.Errorf("L RED instrument: the move reported an error yet the parked spelling exists, so the reading "+
			"above is not trustworthy: %v", err)
	}
	t.Logf("AC#1-three MEASURED, and this is the sentence the交件 carries: on this machine the platform "+
		"refuses an ancestor move while the landing's parent handle is held (%v). Form 乙's ancestor leg on "+
		"Windows is therefore a platform lock and not a route behaviour - so dispatch ①'s deterministic pair "+
		"CANNOT be produced on the real route through any existing seam: the pin is the refusal. Reporting a "+
		"red for the landing's damage numerator would take a simulation, and a mock standing in for the real "+
		"thing is what this repository forbids.", moveL)
}
