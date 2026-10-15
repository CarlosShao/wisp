//go:build !windows

// Ticket 120 AC#1: the window ticket 113's acceptance round recorded as R-113-D,
// made measurable in r1 and judged here.
//
// platformVerifyPlacement is a *check*; the mode it stands for used to be applied
// by os.Chmod on the same *name* a moment later. A link swapped into that gap sent
// the seal into somebody else's tree while the call returned nil. The acceptance
// round measured 154/4000 retries through it on the ticket-113 code and 2087/4000
// before it (its P4 row), and recorded that the two "costs" comments over
// platformVerifyPlacement did not name it.
//
// The instrument is racers and one witness. A goroutine flips the object the caller
// named - the leaf, or one ancestor directory - between "the real object inside the
// tree" and "a link at the foreign witness". The sealing call runs in a loop over
// the same spelling. Every attempt is classified from the witness, not from the
// error: a mode change on the foreign file is a *hit* whatever the call reported,
// and a hit the call reported as nil is counted separately, because "acted on the
// foreign tree" and "said it worked" are two claims and the second is the one this
// repository cannot blame on the racer.
//
// READING RULES (AC#1, verbatim in the ticket): every sample is reported, no
// averaging, no "best of". hits == 0 per sample is the requirement, and it is the
// DEFAULT form of this file: a bare run is a gate. WISP_120_BASELINE=1 downgrades
// exactly the same rows to readings, which is how the pre-fix half of a paired
// reading is taken on code that has not been fixed under it - it is opt-in, it is
// named in the output, and it is not what CI runs.
//
// An instrument that can only read zero is worthless, so liveness is asserted WITH
// the hits and never in place of them: the floor has to have refused at least once,
// the seal has to have succeeded on its own object at least once, and the racer's
// own flip must have recorded no errors. A run where the racer never raced reads as
// dead, and dead is RED here - r1's first cut produced refused=0 and sealed-own=100
// in every sample, which is what a dead instrument looks like and is the reason
// these three assertions cannot be traded against each other.
//
// refusals are split by which leg produced them, because that is the difference
// between "the check before the act caught it" and "the pinned act refused to be
// moved": refused-at-check carries only ErrUnresolvedPath (resolve.go's floor),
// refused-at-act carries ErrUnresolvedPath and ErrNotSealable together
// (winsec_other.go's identity gate). A run with refused>0 and refused-at-act==0 has
// not measured the gate at all.
//
// The container runs as uid 0 (ticket 113's reading note): mode changes are still
// real here because every assertion happens on the container's overlay /tmp, where
// mode bits are stored; what root distorts is the *permission* faces of the fix,
// which this file deliberately does not assert on (they are measured as uid 1000 in
// the ticket's evidence instead).
//
// STATE OF THIS INSTRUMENT (2026-10-15, agent 120-r2). r1 committed this file with
// two defects named in its header and left their repair as next= item 1; both are
// repaired here, and the repair is measured rather than claimed:
//
//  1. the leaf racer is now a rename(2) swap between two staged names that
//     alternate kind. rename replaces the destination atomically and does not follow
//     a symlink standing at the destination (both measured in r2's probe), so the
//     leaf is present in every state at every instant and no attempt can land in an
//     absence. r1's remove-and-recreate left the leaf missing for two of its four
//     steps: on the pre-fix code 1720 of 2000 leaf-swap attempts read "no such file
//     or directory" and the hit rate was a third of what the acceptance round's own
//     probe got. The swap-back is also why the strict form is trustworthy: a racer
//     that spends its time deleting the object is not racing anything.
//  2. the ancestor flip is idempotent. It clears the ancestor by Lstat first: a real
//     directory is emptied (blob, then rmdir) and a name that currently carries a
//     link is unlinked directly. Removing through a name that carries a link would
//     delete the witness - which is what r1's version did to the foreign tree it was
//     supposed to be reading - and mkdir over a directory that still holds its blob
//     is the EEXIST r1's reading reported as "recorded flip errors: 1". Now the flip
//     alternates real and link states with 0 recorded errors, so hits==0 cannot be
//     explained by a racer that stopped flipping.
//  3. a third shape was added, and it is the one the fix stands or falls on. r1
//     measured that os.Root *follows* a symbolic link whose target stays inside the
//     root, and refuses only absolute targets and targets that climb out with "..".
//     Shapes 1 and 2 plant absolute-target links, so a zero there could still be
//     read as "the stdlib said no" rather than "the identity gate works". Shape 3
//     plants an ancestor link whose target is a *sibling directory inside the same
//     parent* - "sub -> other", no "..", reachable by Root - and the witness is
//     other/blob.bin, which is not the tree the call names (R-108-2's rule). Only
//     sameObjectPinned can refuse it. With that gate neutered this file is RED on
//     shape 3 (r2's control reading); that is what makes it a gate and not a
//     decoration.
package winsec_test

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/winsec"
)

const (
	// Four samples, all reported. Scaled by WISP_120_ATTEMPTS for a cheaper probe
	// run; the committed default is what the ticket's evidence rows read.
	sealWindowSamples120  = 4
	sealWindowAttempts120 = 500
)

// sealWindowBaseline120 names the reader form explicitly: same rows, no failures.
// It exists for the pre-fix half of the paired reading and for nobody else - the
// default of this file is the gate (sealWindowStrict120).
func sealWindowBaseline120() bool { return os.Getenv("WISP_120_BASELINE") == "1" }

// sealWindowAssertClosed120 is the strict form: per-sample hits == 0 plus liveness.
// Promoted to the default by next=4 of r1's hand-off, which is what "the instrument
// goes from reading to answering" means. WISP_120_BASELINE=1 is the only way out of
// it, and it says so in every row it prints.
func sealWindowAssertClosed120() bool { return !sealWindowBaseline120() }

func sealWindowAttemptsFor120() int {
	if v := os.Getenv("WISP_120_ATTEMPTS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return sealWindowAttempts120
}

// sealWindowReading120 is one sample's tally, printed whole.
type sealWindowReading120 struct {
	shape        string
	attempts     int
	hitForeign   int // the witness's mode changed: the seal acted on the foreign tree
	hitAndReport int // ... and the call returned nil, i.e. reported success
	refused      int // a refusal: the pre-act floor or the pinned act's gate
	refusedAct   int // ... of which the pinned act's identity gate (both sentinels)
	sealedOwn    int // nil, witness untouched: the plain object in the tree was sealed
	other        int // absent/ENOTDIR/ErrNotSealable alone: neither refusal nor hit
	firstOther   string
}

func (r sealWindowReading120) String() string {
	return fmt.Sprintf("attempts=%d hit-foreign=%d of-which-reported-nil=%d refused-at-check=%d of-which-refused-at-act=%d sealed-own-object=%d other=%d",
		r.attempts, r.hitForeign, r.hitAndReport, r.refused, r.refusedAct, r.sealedOwn, r.other)
}

func (r *sealWindowReading120) note(other string) {
	r.other++
	if r.firstOther == "" {
		r.firstOther = other
	}
}

func (r *sealWindowReading120) add(o sealWindowReading120) {
	r.attempts += o.attempts
	r.hitForeign += o.hitForeign
	r.hitAndReport += o.hitAndReport
	r.refused += o.refused
	r.refusedAct += o.refusedAct
	r.sealedOwn += o.sealedOwn
	r.other += o.other
	if r.firstOther == "" {
		r.firstOther = o.firstOther
	}
}

// sealWindowRun120 drives one shape: `call` seals the spelling, `reset` puts the
// witness back to 0666 before every attempt, `witnessMode` reads the witness.
func sealWindowRun120(t *testing.T, shape string, attempts int, call func() error, witnessMode func() os.FileMode, reset func()) sealWindowReading120 {
	t.Helper()
	r := sealWindowReading120{shape: shape}
	for i := 0; i < attempts; i++ {
		reset()
		err := call()
		r.attempts++
		if got := witnessMode(); got == 0o600 {
			r.hitForeign++
			if err == nil {
				r.hitAndReport++
			}
			continue
		}
		switch {
		case err == nil:
			r.sealedOwn++
		case errors.Is(err, winsec.ErrUnresolvedPath):
			r.refused++
			if errors.Is(err, winsec.ErrNotSealable) {
				r.refusedAct++
			}
		default:
			r.note(err.Error())
		}
	}
	return r
}

// judgeSample120 is AC#1's per-sample rule. In the reader form the same rows are
// logged; nothing else about the reading changes, which is what makes the pair
// comparable.
func judgeSample120(t *testing.T, sample int, r sealWindowReading120) {
	t.Helper()
	t.Logf("AC#1 sample %d/%d shape=%s: %s", sample, sealWindowSamples120, r.shape, r)
	if r.firstOther != "" {
		t.Logf("AC#1 sample %d/%d shape=%s other-shape example: %s", sample, sealWindowSamples120, r.shape, r.firstOther)
	}
	if r.hitForeign == 0 {
		return
	}
	if sealWindowAssertClosed120() {
		t.Errorf("AC#1 RED sample %d/%d shape=%s: the seal landed on the foreign tree %d time(s), %d of them while returning nil (window open on THIS sample; every sample is judgement, no averaging)",
			sample, sealWindowSamples120, r.shape, r.hitForeign, r.hitAndReport)
		return
	}
	t.Logf("AC#1 WINDOW OPEN sample %d/%d shape=%s: hit-foreign=%d of-which-reported-nil=%d (reader mode WISP_120_BASELINE=1: reads the window, does not fail on it)",
		sample, sealWindowSamples120, r.shape, r.hitForeign, r.hitAndReport)
}

// judgeRun120 is the liveness half of the same judgement: a racer that cannot flip,
// or a run in which the floor never refused and the seal never sealed its own
// object, proves nothing and is RED in the default form.
func judgeRun120(t *testing.T, totals sealWindowReading120, flipErrs int, firstFlipErr string) {
	t.Helper()
	t.Logf("AC#1 total shape=%s: %s", totals.shape, totals)
	if totals.firstOther != "" {
		t.Logf("AC#1 shape=%s last-seen other-shape example: %s", totals.shape, totals.firstOther)
	}
	if flipErrs > 0 {
		if sealWindowAssertClosed120() {
			t.Errorf("AC#1 RED instrument shape=%s: the racer's own flip recorded %d error(s) (first: %s) - a racer that cannot flip is not evidence about the window",
				totals.shape, flipErrs, firstFlipErr)
		} else {
			t.Logf("AC#1 FLIP BROKEN shape=%s: %d recorded flip error(s) (first: %s)", totals.shape, flipErrs, firstFlipErr)
		}
	}
	if totals.refused == 0 || totals.sealedOwn == 0 {
		if sealWindowAssertClosed120() {
			t.Errorf("AC#1 RED instrument is dead shape=%s: refused=%d, sealed-own=%d - a zero window reading here proves nothing",
				totals.shape, totals.refused, totals.sealedOwn)
		} else {
			t.Logf("AC#1 INSTRUMENT DEAD reading shape=%s: refused=%d, sealed-own=%d", totals.shape, totals.refused, totals.sealedOwn)
		}
	}
	if totals.refused > 0 && totals.refusedAct == 0 {
		t.Logf("AC#1 shape=%s: every refusal came from the check before the act (refused-at-act=0), so this run does NOT exercise the pinned identity gate", totals.shape)
	}
	if !sealWindowAssertClosed120() {
		t.Logf("AC#1 reader mode (WISP_120_BASELINE=1): the rows above are readings, the gate is the default form of this file")
	}
}

// racerSwap120 keeps flipping what the spelling resolves to until stopped. Flip
// failures are recorded, never raised: t.Fatal from inside this goroutine marks the
// test failed, ends the flip with a Goexit and leaves the loop calling a broken
// step - which is how r1's first cut read as a dead instrument (its defect 2).
type racerSwap120 struct {
	stop  chan struct{}
	done  chan struct{}
	once  sync.Once
	errs  int
	first string
}

func startRacer120(flip func(note func(error))) *racerSwap120 {
	r := &racerSwap120{stop: make(chan struct{}), done: make(chan struct{})}
	note := func(err error) {
		if err == nil {
			return
		}
		r.errs++
		if r.first == "" {
			r.first = err.Error()
		}
	}
	go func() {
		defer close(r.done)
		defer func() {
			if p := recover(); p != nil {
				note(fmt.Errorf("racer panic: %v", p))
			}
		}()
		for {
			select {
			case <-r.stop:
				return
			default:
			}
			flip(note)
		}
	}()
	return r
}

// stopNow is idempotent and joins the goroutine: the shapes call it once to settle
// the counters before they are read and again through defer, so a t.Fatal in the
// middle cannot leave a racer flipping inside a finished test.
func (r *racerSwap120) stopNow() {
	r.once.Do(func() { close(r.stop) })
	<-r.done
}

func (r *racerSwap120) flips() (int, string) { return r.errs, r.first }

// witnessPair is the fixture every shape shares: a place to put bytes that are not
// this tree's, and the two closures the loop needs.
func witnessPair120(t *testing.T, witness string) (func() os.FileMode, func()) {
	t.Helper()
	mode := func() os.FileMode {
		info, err := os.Stat(witness)
		if err != nil {
			t.Fatalf("witness: %v", err)
		}
		return info.Mode().Perm()
	}
	reset := func() {
		if err := os.Chmod(witness, 0o666); err != nil {
			t.Fatal(err)
		}
	}
	return mode, reset
}

// TestPOSIX120SealWindowLeafSwapIsClosed is the shape the ticket names first: the
// leaf the caller named is swapped between a real file and a link at the witness
// while SealFile runs over the same spelling.
func TestPOSIX120SealWindowLeafSwapIsClosed(t *testing.T) {
	base := winsec.SealableTempDirForTest124(t)
	foreign := filepath.Join(base, "foreign")
	if err := os.Mkdir(foreign, 0o700); err != nil {
		t.Fatal(err)
	}
	witness := filepath.Join(foreign, "witness.bin")
	if err := os.WriteFile(witness, []byte("not this tree's data"), 0o600); err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(base, "root")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatal(err)
	}
	leaf := filepath.Join(root, "blob.bin")
	if err := os.WriteFile(leaf, []byte("the tree's own bytes"), 0o666); err != nil {
		t.Fatal(err)
	}
	witnessMode, reset := witnessPair120(t, witness)
	reset()

	// One flip = two atomic replaces. Both staged names are re-staged before they
	// are used, so the leaf itself is never the thing that is missing.
	stageLink := filepath.Join(root, "stage-link.bin")
	stageReal := filepath.Join(root, "stage-real.bin")
	racer := startRacer120(func(note func(error)) {
		if err := os.Symlink(witness, stageLink); err != nil && !errors.Is(err, fs.ErrExist) {
			note(err)
			return
		}
		note(os.Rename(stageLink, leaf))
		if err := os.WriteFile(stageReal, []byte("the tree's own bytes"), 0o666); err != nil {
			note(err)
			return
		}
		note(os.Rename(stageReal, leaf))
	})
	defer racer.stopNow() // idempotent: it also covers a t.Fatal in the loop below

	attempts := sealWindowAttemptsFor120()
	var totals sealWindowReading120
	totals.shape = "leaf-swap"
	for sample := 1; sample <= sealWindowSamples120; sample++ {
		r := sealWindowRun120(t, totals.shape, attempts, func() error { return winsec.SealFile(leaf) }, witnessMode, reset)
		judgeSample120(t, sample, r)
		totals.add(r)
	}
	racer.stopNow() // settled before the counters below are read, so this is not a race on the reading
	flips, firstFlip := racer.flips()
	judgeRun120(t, totals, flips, firstFlip)
}

// TestPOSIX120SealWindowAncestorAbsoluteLinkIsClosed is r1's ancestor shape: a
// directory component is swapped between the real directory holding blob.bin and a
// link at a foreign directory holding a file of the same name. The link target is
// absolute, so os.Root refuses to follow it and this shape can be closed by the
// standard library as well as by the gate.
func TestPOSIX120SealWindowAncestorAbsoluteLinkIsClosed(t *testing.T) {
	base := winsec.SealableTempDirForTest124(t)
	foreignAncestor := filepath.Join(base, "foreign", "sub")
	if err := os.MkdirAll(foreignAncestor, 0o700); err != nil {
		t.Fatal(err)
	}
	witness := filepath.Join(foreignAncestor, "blob.bin")
	if err := os.WriteFile(witness, []byte("not this tree's data"), 0o600); err != nil {
		t.Fatal(err)
	}
	data := filepath.Join(base, "root", "data")
	if err := os.MkdirAll(data, 0o700); err != nil {
		t.Fatal(err)
	}
	ancestor := filepath.Join(data, "sub")
	judgeAncestorSwap120(t, "ancestor-swap-absolute-link", ancestor,
		filepath.Join(ancestor, "blob.bin"), witness, foreignAncestor)
}

// TestPOSIX120SealWindowAncestorRelativeLinkIsClosed is the same swap with the one
// spelling difference that matters for this fix: the link target is a sibling
// directory inside the same parent ("sub -> other", no ".."), which os.Root follows
// (measured in r2's probe: OpenFile on a sibling-relative link handed back the
// TARGET's descriptor). Root's own refusal cannot close this shape, so a zero here
// is the identity gate's reading and nothing else's - and neutering that gate turns
// this case red, which is the control the ticket asks for.
func TestPOSIX120SealWindowAncestorRelativeLinkIsClosed(t *testing.T) {
	base := winsec.SealableTempDirForTest124(t)
	data := filepath.Join(base, "root", "data")
	if err := os.MkdirAll(data, 0o700); err != nil {
		t.Fatal(err)
	}
	other := filepath.Join(data, "other")
	if err := os.MkdirAll(other, 0o700); err != nil {
		t.Fatal(err)
	}
	witness := filepath.Join(other, "blob.bin")
	if err := os.WriteFile(witness, []byte("not this tree's data"), 0o600); err != nil {
		t.Fatal(err)
	}
	ancestor := filepath.Join(data, "sub")
	judgeAncestorSwap120(t, "ancestor-swap-relative-sibling-link", ancestor,
		filepath.Join(ancestor, "blob.bin"), witness, "other")
}

// judgeAncestorSwap120 is the ancestor body for both link spellings: flip the
// directory between "real, holding the tree's own blob" and "link at the other
// directory", seal the same spelling over and over, and judge every sample.
func judgeAncestorSwap120(t *testing.T, shape, ancestor, blob, witness, target string) {
	t.Helper()
	witnessMode, reset := witnessPair120(t, witness)
	// clearByKind is the idempotent half of the flip: whatever stands at the name is
	// taken away by its own kind. A directory is emptied first (rmdir on a non-empty
	// one is EEXIST-then-ENOENT noise, which is r1's defect 2), and a name carrying a
	// link is unlinked directly - removing blob.bin while the name carries a link
	// would go THROUGH the link and delete the witness this reading is made of.
	clearByKind := func(note func(error)) {
		info, err := os.Lstat(ancestor)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				note(err)
			}
			return
		}
		if info.Mode()&os.ModeSymlink != 0 {
			note(os.Remove(ancestor))
			return
		}
		if err := os.Remove(blob); err != nil && !errors.Is(err, fs.ErrNotExist) {
			note(err)
		}
		note(os.Remove(ancestor))
	}
	// Seed the real state before the racer starts, so nothing here races the test's
	// own setup: the first attempt then has something to seal and every later
	// absence is the racer's doing, which is what the reading is about.
	if err := os.Mkdir(ancestor, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(blob, []byte("the tree's own bytes"), 0o666); err != nil {
		t.Fatal(err)
	}
	reset()

	toLink := false
	racer := startRacer120(func(note func(error)) {
		clearByKind(note)
		if toLink {
			note(os.Symlink(target, ancestor))
		} else {
			note(os.Mkdir(ancestor, 0o700))
			note(os.WriteFile(blob, []byte("the tree's own bytes"), 0o666))
		}
		toLink = !toLink
	})
	defer racer.stopNow() // idempotent: it also covers a t.Fatal in the loop below

	attempts := sealWindowAttemptsFor120()
	var totals sealWindowReading120
	totals.shape = shape
	for sample := 1; sample <= sealWindowSamples120; sample++ {
		r := sealWindowRun120(t, shape, attempts, func() error { return winsec.SealFile(blob) }, witnessMode, reset)
		judgeSample120(t, sample, r)
		totals.add(r)
	}
	racer.stopNow() // settled before the counters below are read, so this is not a race on the reading
	flips, firstFlip := racer.flips()
	judgeRun120(t, totals, flips, firstFlip)
}
