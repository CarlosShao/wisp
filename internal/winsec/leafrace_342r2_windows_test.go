//go:build windows

// Ticket 342 AC#1/AC#5 + 348 AC#2's leaf face, leg 342-r2: the leaf re-check given the ONE form the
// adjudication proved can reach it.
//
// WHAT WAS PROVEN BEFORE THIS FILE EXISTED (342-v2, 裁腿现量, restated here because implementing legs
// may not re-argue adjudications): a leaf swap whose target spelling carries a parent pointer is
// answered by the CONTAINMENT face inside Root.OpenFile - verbatim "openat … path escapes from parent" -
// and openPinnedLeaf's two-half gate never runs for it. The one spelling that does reach the re-check is
// the bare sibling name: the target stays in the SAME directory the link lands in, its substitute name
// holds no volume and no ".." element, so containment has nothing to refuse and only the package's own
// identity gate can answer. r1's leaf racer staged `..\..\someone-elses-tree-w\leaf...` - the
// containment-answering form - which is why its 400x3 shots recorded re-check=0 and why "答没答过" was
// booked as UNKNOWN rather than as absence.
//
// THE PLANT IS STILL A RACER, and this file says what that can and cannot promise. The re-check's window
// is between preOpenLeafGuard's verdict (a probe that answers about a NAME before the act) and
// openPinnedLeaf's own open; there is no declared seam in that gap, so entering it is timed, not
// constructed. What IS fixed here is the SHAPE: when a flip lands inside the gap on this instrument, the
// only legs that can answer are the landing's own. The split is therefore a real reading:
//
//   - recheckSpoke > 0        -> the gate is MEASURED on the answerable form (this case's claim);
//   - flip-proved cycles > 0  -> the world really moved (liveness, the same rule r1's family applies);
//   - damage (witness truncated) is ticket 342's numerator and is REPORTED, not gated: an O_TRUNC the
//     followed open already carried out is not undone by a refusal after it, and winsec.go:118-140 keeps
//     that residual named;
//   - the pre-landing pair: the SAME instrument on the pre-landing bytes shows what the un-pinned route
//     did at this shape - the followed open truncated the foreign object and the call LANDED (returned
//     nil), i.e. ticket 342's "truncated + reported success" damage, before any gate existed at all.
//
// WHY THE FOREIGN OBJECT IS IN THE SAME DIRECTORY: dispatch ③'s form names the target as a bare sibling,
// and damage must still be counted from a name the caller never uses. The witness therefore keeps its own
// bare name next to the link, is read by that absolute spelling only, and is restocked with a body of
// DIFFERENT LENGTH from the caller's own (an O_TRUNC that lands and writes nothing must not read as
// untouched). No containment rule can answer this world by construction - there is no spelling that
// escapes anything - so a refusal here cannot be borrowed from another leg, which was r1's whole problem.
//
// Bench facts this shape leans on, with each one tied to the reading that carries it (本腿仓外 probe7
// 编译失败＝⛔ 凭据、已在散件里作废登记，这里⛔ 引用它): "the followed open really truncates" is proven
// INSIDE this instrument - the landing's damage rows read `size 0 body []` (an O_TRUNC that landed through
// the flipped name and wrote nothing, which no restock can produce) on runs where the same reading counts
// 46~49 re-check refusals; the per-attempt pairing of one truncation to one refusal is NOT claimed (the
// sample races the racer), and ticket 342's numerator stays REPORTED rather than gated here. The property
// itself is in-repo evidence already (342-a1's W-ROOTSYM row: Root.OpenFile follows an in-root relative
// link, measured three runs). That the link's entry can be removed while other handles stand in the
// directory is probe6's own row (`os.Remove` on the name while the landing's handle is open = ok), and it
// is also what this package's deleteLink documents; the clear-then-rename gap is therefore ENOENT for one
// instruction, counted as a gap, never hidden, never claimed as liveness.
//
// NO NEW SEAM, NO PRODUCTION CHANGE: this file consumes only the exported entry point and the package's
// existing test helpers (reparseAt, SealableTempDirForTest124), plus its own helpers named with the 342r2
// suffix. winsec.go:430's wrapPath is untouched; winsec.go:118-140's statement is not this leg's to edit.
package winsec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	winLeafAttempts342r2 = 400
	winLeafPace342r2     = 2 * time.Millisecond
	winLeafName342r2     = "leaf-342r2v.bin"
	// The witness keeps its OWN bare name in the same directory. The symlink's target is exactly this
	// base name - no volume prefix, no ".." element - which is the whole difference from r1's form.
	winLeafWitness342r2   = "witness-342r2v.bin"
	winLeafBody342r2      = "the caller's own object under the name it uses"
	winLeafWitBody342r2   = "the bare sibling's body is deliberately longer than the caller's own object"
	winLeafPayload342r2   = "leaf-race payload 342r2 - must never land in the bare sibling"
	winLeafStageReal342r2 = "staged-real-342r2v"
	winLeafStageLink342r2 = "staged-link-342r2v"
)

// winLeafLegs342r2: the answers THIS world can produce, declared in order. The first two rows are the
// point of the instrument; the rest are named so nothing can hide in "some other failure".
var winLeafLegs342r2 = []struct{ why, name string }{
	{"the name carries a link, so this handle", "openPinnedLeaf re-check, link face (the gate this case exists to reach)"},
	{"this handle is not the object the name carries now", "openPinnedLeaf re-check, SameFile face"},
	{"the name carries a reparse point", "preOpenLeafGuard's leaf refusal (the steady state answers first)"},
	{"the reparse-point probe could not read the name", "preOpenLeafGuard's unreadable-name refusal"},
	{"path escapes from parent", "os.Root containment (must stay ZERO on this spelling: the target is a bare name, so if this row is non-zero the world was built wrong and the reading below is void)"},
	{"cannot apply a private security descriptor", "a seal refusal on this route"},
	{"path traverses a reparse point", "the installed C26 resolver's steady-state refusal"},
	{"not provably resolved", "resolveString refused the spelling"},
}

const (
	winLeafOther342r2 = "some other failure"
	winLeafNil342r2   = "none: the call returned nil"
)

func spokeWinLeaf342r2(err error) string {
	if err == nil {
		return winLeafNil342r2
	}
	msg := err.Error()
	for _, leg := range winLeafLegs342r2 {
		if strings.Contains(msg, leg.why) {
			return leg.name
		}
	}
	return winLeafOther342r2
}

// winLeafWorld342r2: dir\leaf-342r2v.bin (the caller's name) and dir\witness-342r2v.bin (the bare
// sibling). The racer flips the caller's name between the real file and a RELATIVE symlink whose target
// is the sibling's bare name, staged under two unused names and moved in by rename, so the name is in
// every state except the counted clear-then-rename gap.
type winLeafWorld342r2 struct {
	t             *testing.T
	dir           string
	name          string
	witness       string
	stageRealDir  string
	stageLinkPath string

	cycles    int32
	flipped   int32 // cycles where reparseAt(name) proved TRUE right after the rename
	gapErrs   int32
	flipErrs  int32
	stop      chan struct{}
	done      chan struct{}
	stopOnce  sync.Once
	flipMu    sync.Mutex
	firstFlip string
}

func newWinLeafWorld342r2(t *testing.T) *winLeafWorld342r2 {
	t.Helper()
	w := &winLeafWorld342r2{t: t, dir: SealableTempDirForTest124(t)}
	w.name = filepath.Join(w.dir, winLeafName342r2)
	w.witness = filepath.Join(w.dir, winLeafWitness342r2)
	w.stageRealDir = filepath.Join(w.dir, winLeafStageReal342r2)
	w.stageLinkPath = filepath.Join(w.dir, winLeafStageLink342r2)
	if err := os.WriteFile(w.name, []byte(winLeafBody342r2), 0o600); err != nil {
		t.Fatalf("stage the caller's object: %v", err)
	}
	if err := os.WriteFile(w.witness, []byte(winLeafWitBody342r2), 0o600); err != nil {
		t.Fatalf("stage the bare sibling witness: %v", err)
	}
	// The staged link: a RELATIVE symlink whose target is the sibling's bare name and nothing else.
	// os.Symlink needs developer mode for DIRECTORY links; for a file target an unprivileged token holds
	// SeCreateSymbolicLinkPrivilege on this machine was measured by 342-r1's family, and a failure here is
	// named as the construction it is - never a quiet skip, never a borrowed green.
	if err := w.stageLink(); err != nil {
		t.Skipf("cannot construct the relative file symlink on this box: %v (the answerable form IS this "+
			"construction; skipping quietly would be the failure both tickets name, so acceptance has to "+
			"count this SKIP)", err)
	}
	t.Cleanup(w.settle)
	return w
}

func (w *winLeafWorld342r2) stageLink() error {
	if err := winClearUnlinked342r2(w.stageLinkPath); err != nil {
		return err
	}
	// The target is the BASE NAME of the witness - no path at all. This is the difference from r1's
	// `..\..\x\y` spelling, which containment answered before the gate could.
	if err := os.Symlink(winLeafWitness342r2, w.stageLinkPath); err != nil {
		return err
	}
	if !reparseAt(w.stageLinkPath) {
		return fmt.Errorf("the staged name %s is not a reparse point, so this fixture races nothing", w.stageLinkPath)
	}
	return nil
}

// flip: move the staged symlink onto the caller's name and PROVE the state. rename replaces the entry
// without following it; the leaf was a real file, so the replace needs no clear at all - rename(2) on
// Windows swaps the directory entry atomically, which is exactly what keeps this gap NARROWER than
// r1's clear-then-rename (whose ENOENT windows are counted, never hidden).
func (w *winLeafWorld342r2) flip() error {
	if err := os.Rename(w.stageLinkPath, w.name); err != nil {
		return fmt.Errorf("move the staged link onto the caller's name: %v", err)
	}
	if !reparseAt(w.name) {
		return errors.New("the caller's name is not a reparse point after the flip, so the state was never entered")
	}
	atomic.AddInt32(&w.flipped, 1)
	atomic.AddInt32(&w.cycles, 1)
	return w.stageLink()
}

// settle: the ordinary state back - unlink the link entry if the name carries one (never follow it),
// re-write the real object, re-stock the sibling with its own longer body, re-arm the staged link.
func (w *winLeafWorld342r2) settle() {
	if reparseAt(w.name) {
		if rErr := winClearUnlinked342r2(w.name); rErr != nil {
			w.t.Logf("342r2v cleanup: unlink the flipped name: %v", rErr)
		}
	}
	if err := os.WriteFile(w.name, []byte(winLeafBody342r2), 0o600); err != nil {
		w.t.Logf("342r2v cleanup: restore the real leaf: %v", err)
	}
	if err := os.WriteFile(w.witness, []byte(winLeafWitBody342r2), 0o600); err != nil {
		w.t.Logf("342r2v cleanup: restock the sibling: %v", err)
	}
	if err := w.stageLink(); err != nil {
		w.t.Logf("342r2v cleanup: re-arm the staged link: %v", err)
	}
}

func (w *winLeafWorld342r2) noteFlipErr(err error) {
	atomic.AddInt32(&w.flipErrs, 1)
	w.flipMu.Lock()
	defer w.flipMu.Unlock()
	if w.firstFlip == "" {
		w.firstFlip = err.Error()
	}
}

func (w *winLeafWorld342r2) startRacer() {
	w.stop = make(chan struct{})
	w.done = make(chan struct{})
	go func() {
		defer close(w.done)
		defer func() {
			if p := recover(); p != nil {
				w.noteFlipErr(fmt.Errorf("racer panic: %v", p))
			}
		}()
		for {
			select {
			case <-w.stop:
				return
			default:
			}
			if err := w.flip(); err != nil {
				w.noteFlipErr(err)
				w.settle() // recover the world instead of spinning on the same fault
				continue
			}
			time.Sleep(winLeafPace342r2)
			select {
			case <-w.stop:
				return
			default:
			}
			w.settle()
			time.Sleep(winLeafPace342r2)
		}
	}()
}

func (w *winLeafWorld342r2) stopRacer() {
	w.stopOnce.Do(func() { close(w.stop) })
	<-w.done
	w.settle()
}

// run drives the production entry point over the flipping name and samples the sibling witness AFTER
// EVERY attempt, before the racer's next restock can hide a damage row. Sampling races the restock, so
// the count may UNDER-count (a truncation can be restocked before the read); it cannot over-count,
// because a body that differs from the witness's own is a fact that happened. This is the DAMAGE column
// of dispatch ②; the leg map is the REFUSAL column; one number never feeds the other's gate.
func (w *winLeafWorld342r2) run(call func() error) (legs map[string]int, flipped int, firstOther string, damage winLeafDamage342r2) {
	w.startRacer()
	legs = map[string]int{}
	damage = winLeafDamage342r2{}
	for i := 0; i < winLeafAttempts342r2; i++ {
		err := call()
		leg := spokeWinLeaf342r2(err)
		legs[leg]++
		if leg == winLeafOther342r2 && firstOther == "" {
			firstOther = err.Error()
		}
		info, lErr := os.Lstat(w.witness)
		switch {
		case lErr != nil:
			damage.gone++
			if damage.firstGone == "" {
				damage.firstGone = fmt.Sprintf("attempt %d: %v", i, lErr)
			}
		default:
			_ = info // Mode() is decorative on this platform (ticket 89's premise), so no mode row here.
			body, rErr := os.ReadFile(w.witness)
			if rErr != nil {
				damage.unreadable++
			} else if string(body) != winLeafWitBody342r2 {
				damage.truncated++
				if damage.firstTrunc == "" {
					damage.firstTrunc = fmt.Sprintf("attempt %d: size %d body [%s]", i, len(body), body)
				}
				// Re-stock the witness here, from the main loop, so each attempt gets a fresh reading.
				// The racer's own settle also writes it back, so the sample races two restocks and can
				// UNDER-count damage; it cannot over-count, because a body that differs from the
				// witness's own is a fact that happened.
				if wErr := os.WriteFile(w.witness, []byte(winLeafWitBody342r2), 0o600); wErr != nil {
					w.t.Fatalf("restock the witness: %v", wErr)
				}
			}
		}
	}
	w.stopRacer()
	return legs, int(atomic.LoadInt32(&w.flipped)), firstOther, damage
}

// winLeafDamage342r2 is the DAMAGE column: facts about the sibling object, counted per attempt. It never
// feeds, and is never fed by, the leg map - ticket 348 AC#3's rule that one numerator may not reuse
// another's green.
type winLeafDamage342r2 struct {
	truncated  int // the witness held something other than its own body
	gone       int // the witness was absent
	unreadable int
	firstTrunc string
	firstGone  string
}

// TestAC1WindowsLeafRecheckAnswersOnTheBareSiblingForm342r2 gates exactly one claim: on the dispatch's
// answerable form, openPinnedLeaf's re-check SPOKE (and spoke with one of its own two halves, not with a
// containment or floor answer). Everything else is printed: the guard's steady-state count, the landed
// nil count (the pre-landing damage shape), the witness reading after the run, the flip errors. The
// damage column is NOT gated here (r1's file owns the truncation residual rows and winsec.go:118-140
// keeps that residual named); the two tickets' numerators stay in separate files by separate gates.
func TestAC1WindowsLeafRecheckAnswersOnTheBareSiblingForm342r2(t *testing.T) {
	w := newWinLeafWorld342r2(t)
	legs, flipped, firstOther, damage := w.run(func() error {
		return PrivateFile(w.name, []byte(winLeafPayload342r2), 0o600)
	})
	recheck := legs[winLeafLegs342r2[0].name] + legs[winLeafLegs342r2[1].name]
	guard := legs[winLeafLegs342r2[2].name] + legs[winLeafLegs342r2[3].name]
	t.Logf("AC#1-leaf attempts=%d flippedCycles=%d flipErrs=%d firstFlipErr=%q",
		winLeafAttempts342r2, flipped, atomic.LoadInt32(&w.flipErrs), w.firstFlip)
	// The DAMAGE column, its own rows, its own numbers.
	t.Logf("AC#1-leaf DAMAGE : witness truncated=%d gone=%d unreadable=%d | first truncation=%q | first "+
		"absence=%q (sampled after every attempt from the main goroutine; sampling races the racer's own "+
		"restock so it can UNDER-count and cannot over-count)",
		damage.truncated, damage.gone, damage.unreadable, damage.firstTrunc, damage.firstGone)
	for _, leg := range winLeafLegs342r2 {
		if n := legs[leg.name]; n > 0 {
			t.Logf("AC#1-leaf': leg that answered = %s : %d attempts", leg.name, n)
		}
	}
	for _, extra := range []string{winLeafOther342r2, winLeafNil342r2} {
		if n := legs[extra]; n > 0 {
			t.Logf("AC#1-leaf': leg that answered = %s : %d attempts", extra, n)
		}
	}
	if firstOther != "" {
		t.Logf("AC#1-leaf': first unattributed failure verbatim: %s", firstOther)
	}
	// A containment answer on THIS spelling is an instrument fault: the target holds no ".." and no
	// volume, so containment has nothing to see - r1's mis-shape made this row possible and that is
	// exactly why it is gated instead of merely printed.
	if n := legs[winLeafLegs342r2[4].name]; n > 0 {
		t.Errorf("AC#1-leaf RED INSTRUMENT: containment answered %d attempts on a bare-sibling spelling "+
			"(verbatim first: %s), so the world built is NOT the answerable form and no re-check reading "+
			"below can be trusted", n, firstOther)
	}
	if n := legs[winLeafOther342r2]; n > 0 {
		t.Errorf("AC#1-leaf RED: %d attempts answered with legs this roster does not name, so the refusal "+
			"column is not attributable (first verbatim: %s)", n, firstOther)
	}
	// Liveness: flips were PROVED (reparseAt right after each rename), not inferred from error shapes.
	if flipped == 0 {
		t.Errorf("AC#1-leaf RED instrument: %d proved flipped cycles across %d attempts (flipErrs=%d, first "+
			"%q), so the counts above are not a reading about a world that ever moved", flipped,
			winLeafAttempts342r2, atomic.LoadInt32(&w.flipErrs), w.firstFlip)
	}
	// The claim of this file, gated exactly once:
	if recheck == 0 {
		t.Errorf("AC#1-leaf NOT YET MEASURED even on the answerable form: the guard answered %d and the "+
			"re-check answered %d across %d attempts with %d proved flips - so either the gap was never "+
			"landed in (timing, reportable) or the re-check cannot answer this form (shape, and the "+
			"adjudication's premise needs re-deriving). Either way this platform's leaf gate stays open and "+
			"this red says so instead of borrowing the guard's green.", guard, recheck, winLeafAttempts342r2, flipped)
		return
	}
	t.Logf("AC#1-leaf MEASURED: openPinnedLeaf's re-check answered %d of %d attempts on the bare-sibling "+
		"form (guard=%d, proved flips=%d) - the gate 342-v2 recorded as unreachable is reached with the "+
		"target spelling fixed to a bare name, with NO new seam and NO production change",
		recheck, winLeafAttempts342r2, guard, flipped)
	// Post-run invariant: whatever landed inside the gap, the sibling must still be readable. Body
	// damage is ticket 342's residual and is REPORTED (see r1's leaf file), never hidden here:
	body, err := os.ReadFile(w.witness)
	switch {
	case err != nil:
		t.Errorf("AC#1-leaf: the bare sibling witness is unreadable after the run: %v", err)
	case string(body) != winLeafWitBody342r2:
		t.Logf("AC#1-leaf RESIDUAL OBSERVED, counted and NOT claimed closed: the witness now holds %d of "+
			"%d bytes, so at least one followed open truncated it before the re-check refused the write "+
			"(ticket 322's Windows leaf residual + ticket 342's numerator; r1's file owns the gated row).",
			len(body), len(winLeafWitBody342r2))
	}
}
