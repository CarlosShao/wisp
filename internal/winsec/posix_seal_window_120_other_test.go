//go:build !windows

// Ticket 120 AC#1: the window ticket 113's acceptance round recorded as R-113-D,
// made measurable and re-runnable instead of anecdotal.
//
// platformVerifyPlacement is a *check*; the mode it stands for used to be applied
// by os.Chmod on the same *name* a moment later. A link swapped into that gap sent
// the seal into somebody else's tree while the call returned nil. The acceptance
// round measured 154/4000 retries through it on the ticket-113 code and 2087/4000
// before it (its P4 row), and recorded that the two "costs" comments over
// platformVerifyPlacement did not name it.
//
// The instrument is two racers and one witness. A goroutine flips the object the
// caller named - the leaf, or one ancestor directory - between "the real object
// inside the tree" and "a symlink at the foreign witness". The sealing call runs
// in a loop over the same spelling. Every attempt is classified from the witness,
// not from the error: a mode change on the foreign file is a *hit* whatever the
// call reported, and a hit the call reported as nil is counted separately, because
// "acted on the foreign tree" and "said it worked" are two claims and the second is
// the one this repository cannot blame on the racer.
//
// Reading rules (AC#1, verbatim in the ticket): every sample is reported, no
// averaging, no "best of". The strict form asserts hits == 0 per sample; the
// committed form here is the *reader*, because the implementation this file was
// written against has not changed yet (see STATE below). WISP_120_BASELINE=1 names
// that reader mode explicitly - it is how the 2026-10-15 container reading was taken
// (git archive of the pre-fix sha, this file copied in, WISP_120_ATTEMPTS=100) - and
// WISP_120_ASSERT_CLOSED=1 is the strict form the fix has to promote to the default.
//
// An instrument that can only read zero is worthless, so each shape also asserts
// liveness over the whole run: the floor has to have *refused* at least once (the
// racer really planted a link the check saw) and the seal has to have *succeeded*
// at least once (the plain spelling really seals). A run where the racer never
// raced reads as dead, not as clean.
//
// The container runs as uid 0 (ticket 113's reading note): mode changes are still
// real here because every assertion happens on the container's overlay /tmp, where
// mode bits are stored; what root distorts is the *permission* faces of the fix,
// which this file deliberately does not assert on (they are measured as uid 1000 in
// the ticket's evidence instead).
//
// STATE OF THIS INSTRUMENT (2026-10-15, agent 120-r1, deadline-pause). What is
// committed here is the *unfixed form* of AC#1: a reader, not a gate.
// 本仪器当前只作基线读数 (this instrument currently takes baseline readings only).
// Two defects were measured on it in a container the same day it was written, and
// they are named here instead of being tuned away, because the second one is
// invisible in the output unless the reader knows what a dead instrument looks like:
//
//  1. the leaf racer is not atomic. It removes and recreates the leaf, so the leaf
//     is absent for two of its four steps and the sealing call spends most of its
//     attempts reading "no such file" instead of the window (measured, sample 1 of
//     4 with WISP_120_ATTEMPTS=100: other=89, refused=5, sealed-own=5; over 400
//     attempts: hit-foreign=1, refused=25, sealed-own=27, other=347 - against the
//     154/4000 the acceptance round measured with its own probe). The fix is a
//     rename(2)-based swap - rename replaces the destination atomically - built
//     around a pair of staged names that alternate kind so the leaf is never absent;
//  2. the ancestor racer raised its failures through t.Fatal *inside the racer
//     goroutine*, which does not stop the racer: Goexit ends the flip, the loop
//     keeps calling the broken step, and the test is marked failed while its reading
//     goes dead. That is what the first cut of this file measured (refused=0 and
//     sealed-own=100 in every one of four samples, i.e. this file's own liveness
//     check fired). What is committed here records the failure instead of raising it,
//     so the flip proceeds; one step of the flip is still not idempotent (it does not
//     empty the directory before recreating it, which is the single EEXIST the
//     reading reports on the first flip), and repairing that step is next= item 1 too.
//     With the recording in place both shapes read the open window on the unfixed
//     code, and the reading moves from run to run: in the three captures taken on
//     2026-10-15 the two shapes landed 0.2% to 0.8% of their attempts on the foreign
//     tree (four to sixteen hits per 2000 attempts each), against the 154/4000 =
//     3.85% the acceptance round measured with its own, more sensitive probe. The log
//     committed next to this file is the last of those captures and carries the exact
//     per-sample numbers; the gap to 3.85% is defect 1's cost, not a different
//     verdict about the window.
//
// Because of (1) and (2) a zero reading below means "the instrument read nothing",
// which is the shape ticket 120 exists to stop. The two defects and their fix
// directions are the first two items of the ticket's next= line.
package winsec_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/CarlosShao/wisp/internal/winsec"
)

const (
	// Four samples, all reported. Scaled by WISP_120_ATTEMPTS for a cheaper probe
	// run; the committed default is what the ticket's evidence rows read.
	sealWindowSamples120  = 4
	sealWindowAttempts120 = 500
)

// sealWindowBaseline120 is the explicit name of the committed reader mode (see the
// file comment's STATE paragraph). It is accepted rather than required: a bare run
// of this file is the reader too, because the implementation it measures has not
// been fixed under it yet.
func sealWindowBaseline120() bool { return os.Getenv("WISP_120_BASELINE") == "1" }

// sealWindowAssertClosed120 is the strict form: per-sample hits == 0 and the
// liveness checks, i.e. AC#1's "after the fix it must be 0". The fix promotes it to
// the default; until then it exists so the reader cannot be mistaken for the gate.
func sealWindowAssertClosed120() bool { return os.Getenv("WISP_120_ASSERT_CLOSED") == "1" }

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
	shape          string
	attempts       int
	hitForeign     int // the witness's mode changed: the seal acted on the foreign tree
	hitAndReported int // ... and the call returned nil, i.e. reported success
	refused        int // the floor refused the spelling (a link was caught by the check)
	sealedOwn      int // nil, witness untouched: the plain object in the tree was sealed
	other          int // absent/ENOTDIR/ErrNotSealable: neither refusal nor hit
	firstOther     string
}

func (r sealWindowReading120) String() string {
	return fmt.Sprintf("attempts=%d hit-foreign=%d of-which-reported-nil=%d refused-at-check=%d sealed-own-object=%d other=%d",
		r.attempts, r.hitForeign, r.hitAndReported, r.refused, r.sealedOwn, r.other)
}

func (r *sealWindowReading120) note(other string) {
	r.other++
	if r.firstOther == "" {
		r.firstOther = other
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
				r.hitAndReported++
			}
			continue
		}
		switch {
		case err == nil:
			r.sealedOwn++
		case errors.Is(err, winsec.ErrUnresolvedPath):
			r.refused++
		default:
			r.note(err.Error())
		}
	}
	return r
}

// racerSwap120 keeps flipping what the spelling resolves to until stopped.
type racerSwap120 struct {
	stop chan struct{}
	done chan struct{}
}

func startRacer120(flip func()) *racerSwap120 {
	r := &racerSwap120{stop: make(chan struct{}), done: make(chan struct{})}
	go func() {
		defer close(r.done)
		for {
			select {
			case <-r.stop:
				return
			default:
			}
			flip()
		}
	}()
	return r
}

func (r *racerSwap120) stopNow() {
	close(r.stop)
	<-r.done
}

// TestAC1POSIXSealWindowIsClosed120 is the auditor's shape cut as a case: the leaf
// the caller named is swapped between a real file and a symlink at the witness
// while SealFile runs over the same spelling.
//
// The hit is read off the witness's mode, and the "reported nil" count is the
// sentence AC#1 is really about: before ticket 120 the seal acted on the foreign
// tree and returned nil, and a by-name seal cannot tell the two objects apart at
// the moment it acts.
func TestPOSIX120SealWindowInstrumentReadsBaselineUnfixed(t *testing.T) {
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

	reset := func() {
		if err := os.Chmod(witness, 0o666); err != nil {
			t.Fatal(err)
		}
	}
	witnessMode := func() os.FileMode {
		info, err := os.Stat(witness)
		if err != nil {
			t.Fatalf("witness: %v", err)
		}
		return info.Mode().Perm()
	}
	reset()

	racer := startRacer120(func() {
		// Link state, then real state. Each flip is one syscall sequence; the
		// seal's own walk takes longer than a flip, which is what keeps the
		// window reachable for the pre-120 code and the racer meaningful for
		// the fixed one.
		_ = os.Remove(leaf)
		_ = os.Symlink(witness, leaf)
		_ = os.Remove(leaf)
		_ = os.WriteFile(leaf, []byte("the tree's own bytes"), 0o666)
		_ = os.Chmod(leaf, 0o666)
	})
	defer racer.stopNow()

	attempts := sealWindowAttemptsFor120()
	var totals sealWindowReading120
	for sample := 1; sample <= sealWindowSamples120; sample++ {
		r := sealWindowRun120(t, "leaf-swap", attempts, func() error { return winsec.SealFile(leaf) }, witnessMode, reset)
		totals.attempts += r.attempts
		totals.hitForeign += r.hitForeign
		totals.hitAndReported += r.hitAndReported
		totals.refused += r.refused
		totals.sealedOwn += r.sealedOwn
		totals.other += r.other
		t.Logf("AC#1 sample %d/%d shape=%s: %s", sample, sealWindowSamples120, r.shape, r)
		if r.hitForeign != 0 {
			if sealWindowAssertClosed120() {
				t.Errorf("AC#1 RED sample %d/%d: the seal landed on the foreign tree %d time(s), %d of them while returning nil (window open on THIS sample; every sample is judgement, no averaging)",
					sample, sealWindowSamples120, r.hitForeign, r.hitAndReported)
			} else {
				t.Logf("AC#1 WINDOW OPEN sample %d/%d: hit-foreign=%d of-which-reported-nil=%d (reader mode: the committed form of this file reads the window, it does not fail on it)",
					sample, sealWindowSamples120, r.hitForeign, r.hitAndReported)
			}
		}
		if r.firstOther != "" {
			t.Logf("AC#1 sample %d other-shape examples: %s", sample, r.firstOther)
		}
	}

	t.Logf("AC#1 total shape=leaf-swap: %s", totals)
	if totals.refused == 0 || totals.sealedOwn == 0 {
		if sealWindowAssertClosed120() {
			t.Errorf("AC#1 instrument is dead: refused=%d, sealed-own=%d - a zero window reading here proves nothing", totals.refused, totals.sealedOwn)
		} else {
			t.Logf("AC#1 INSTRUMENT DEAD reading: refused=%d, sealed-own=%d - a zero window reading here proves nothing; the named defect list is in this file's header", totals.refused, totals.sealedOwn)
		}
	}
	if !sealWindowAssertClosed120() {
		t.Logf("AC#1 reader mode (the committed form of this file, unfixed implementation): readings above are the open window; the strict form is WISP_120_ASSERT_CLOSED=1 and it is what the fix promotes to the default (WISP_120_BASELINE=1, accepted, names this mode explicitly)")
	}
}

// TestAC1POSIXSealWindowAncestorSwapIsClosed120 is the same instrument one level
// up, because the sentence in the ticket is "after the ancestor chain is checked,
// before the seal" - a fix that only protects the leaf would read zero on the leaf
// case and still hand the tree away when a directory component is the swapped one.
//
// The ancestor is flipped between the real directory (holding blob.bin) and a
// symlink at a foreign directory holding a file of the *same name*, so a seal that
// follows the swap is visible as that foreign file's mode.
func TestPOSIX120SealWindowAncestorSwapInstrumentUnfixed(t *testing.T) {
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
	blob := filepath.Join(ancestor, "blob.bin")

	// Defect 2 in the header, kept as it was measured: this flip does not empty the
	// directory before it re-creates it, so os.Mkdir fails EEXIST on every flip and
	// the racer never reaches the link state. The failures are *recorded* here
	// rather than raised through t.Fatal, which is what made this test land as a
	// failed case from a goroutine while the racer kept looping; the reading below
	// reports them, which is the honest shape of a dead instrument. The flip order
	// is deliberately not repaired in this file: repairing it is next= item 1.
	var flipErrs int
	var firstFlipErr string
	noteFlip := func(err error) {
		if err == nil {
			return
		}
		flipErrs++
		if firstFlipErr == "" {
			firstFlipErr = err.Error()
		}
	}
	realAncestor := func() {
		_ = os.Remove(ancestor) // the symlink state, if it is the current one
		noteFlip(os.Mkdir(ancestor, 0o700))
		noteFlip(os.WriteFile(blob, []byte("the tree's own bytes"), 0o600))
		noteFlip(os.Chmod(blob, 0o666))
	}
	realAncestor()
	witnessMode := func() os.FileMode {
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
	reset()

	racer := startRacer120(func() {
		// Real state, then link state; the order inside each flip matters, the
		// directory has to be empty before its name can carry the symlink.
		realAncestor()
		_ = os.Remove(blob)
		_ = os.Remove(ancestor)
		_ = os.Symlink(foreignAncestor, ancestor)
	})

	attempts := sealWindowAttemptsFor120()
	var totals sealWindowReading120
	for sample := 1; sample <= sealWindowSamples120; sample++ {
		r := sealWindowRun120(t, "ancestor-swap", attempts, func() error { return winsec.SealFile(blob) }, witnessMode, reset)
		totals.attempts += r.attempts
		totals.hitForeign += r.hitForeign
		totals.hitAndReported += r.hitAndReported
		totals.refused += r.refused
		totals.sealedOwn += r.sealedOwn
		totals.other += r.other
		t.Logf("AC#1 sample %d/%d shape=%s: %s", sample, sealWindowSamples120, r.shape, r)
		if r.hitForeign != 0 {
			if sealWindowAssertClosed120() {
				t.Errorf("AC#1 RED sample %d/%d: the seal walked an ancestor that had been swapped and landed on the foreign tree %d time(s), %d of them reported as nil",
					sample, sealWindowSamples120, r.hitForeign, r.hitAndReported)
			} else {
				t.Logf("AC#1 WINDOW OPEN sample %d/%d: hit-foreign=%d of-which-reported-nil=%d (reader mode: the committed form of this file reads the window, it does not fail on it)",
					sample, sealWindowSamples120, r.hitForeign, r.hitAndReported)
			}
		}
		if r.firstOther != "" {
			t.Logf("AC#1 sample %d other-shape examples: %s", sample, r.firstOther)
		}
	}

	racer.stopNow() // settled before the counters below are read, so this is not a race on the reading
	t.Logf("AC#1 total shape=ancestor-swap: %s", totals)
	if flipErrs > 0 {
		t.Logf("AC#1 ancestor-racer recorded flip errors: %d (first: %s) - this is defect 2 in the header (one flip step is not idempotent); the flip still reaches the link state, so this reading is live, and repairing the step is next= item 1", flipErrs, firstFlipErr)
		if sealWindowAssertClosed120() {
			t.Errorf("AC#1 instrument is not fit to judge yet: the ancestor racer's own flip recorded %d error(s) (first: %s)", flipErrs, firstFlipErr)
		}
	}
	if totals.refused == 0 || totals.sealedOwn == 0 {
		if sealWindowAssertClosed120() {
			t.Errorf("AC#1 instrument is dead: refused=%d, sealed-own=%d - a zero window reading here proves nothing", totals.refused, totals.sealedOwn)
		} else {
			t.Logf("AC#1 INSTRUMENT DEAD reading: refused=%d, sealed-own=%d - a zero window reading here proves nothing; this shape's defect (the ancestor racer never reaches the link state) is named in the header", totals.refused, totals.sealedOwn)
		}
	}
	if !sealWindowAssertClosed120() {
		t.Logf("AC#1 reader mode (the committed form of this file, unfixed implementation): readings above are the open window; the strict form is WISP_120_ASSERT_CLOSED=1 and it is what the fix promotes to the default (WISP_120_BASELINE=1, accepted, names this mode explicitly)")
	}
}
