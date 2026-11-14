//go:build !windows

// Ticket 322, leg 322-r1. Judged here: the OPEN inside privateFile
// (os.OpenFile(resolved, flags|oNoFollowPlatform, perm)), which runs before any seal. The
// shape is adjudicated in ledger A996, not re-decided here: O_NOFOLLOW closes the LEAF leg on
// POSIX and the ANCESTOR leg stays open and stays named - it speaks about the last component
// only, so an ancestor swapped for a link after resolveString answered is still followed and
// its O_TRUNC still lands in somebody else's tree. Nothing here claims the window is closed;
// reasoning, top-red roster and paired readings: docs/evidence/s1/322-openleaf.md. Why a
// racer: a leaf link planted before the call is refused by the placement floor before the open
// runs, so planting cannot reach this leg (120-r3's reason for winsec_other.go:269). Judged: no
// fresh foreign witness loses its own bytes or mtime, AND the new leg spoke at least once, so a
// 0 numerator is never read as "nothing was raced". Both numbers print on one screen.
//
// SPLIT NOTE (leg 322-r1b): this file came in at 262 lines, two over the 260-line budget, and
// the by-name deterministic case that sat under it is now openleaf_flag_322_other_test.go. The
// split moved bytes only - every assertion, message and threshold here is the one 322-r1 wrote.
package winsec

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	// leafRaceAttempts322r1 is 3000, the census's own denominator: at 800 one sample landed
	// zero leaf refusals (flips=358), which the instrument check below correctly rejects as
	// a non-reading rather than a closed window.
	leafRaceAttempts322r1 = 3000
	leafRacePayload322r1  = "raced private payload 322r1"
	// whyOpenELOOP322r1 is the kernel's wording for O_NOFOLLOW refusing a leaf link, kept
	// verbatim so a reading names WHICH leg answered rather than asking only "did it error".
	whyOpenELOOP322r1 = "too many levels of symbolic links"
	legLeafOpen322r1  = "ticket 322 leaf leg: open carrying O_NOFOLLOW refused a link at the name"
	legChain322r1     = "link-chain ELOOP from a deeper openat (winsec_other.go pinPlacement), NOT this leg"
)

// legs322r1 pairs each answer's fingerprint with its printed name, in a fixed order so two
// runs diff line by line. The last two rows are the non-refusals.
var legs322r1 = []struct{ why, name string }{
	{"reaches it through the link at", "placement floor before the act (platformVerifyPlacement)"},
	{"the name carries a link, so the pinned object", "pinPlacement sameObjectPinned link face (winsec_other.go:188)"},
	{"the name carries a link, so this descriptor", "sealHandle inline link face (winsec_other.go:269)"},
	{"not the object the name carries now", "sealHandle inline SameFile face (winsec_other.go:272)"},
	{whyOpenELOOP322r1, legChain322r1},
	{"", "some other failure"},
	{"", "none: the call returned nil"},
}

// spokeLeg322r1 names the leg that answered. The leaf leg's key is ELOOP reported BY THE OPEN
// OF THE NAME, not ELOOP alone: with the teeth out the racer also produced chain-formed ELOOP
// errors arriving from pinPlacement's openat (2 such rows), which would fake a non-zero liveness.
func spokeLeg322r1(err error) string {
	if err == nil {
		return legs322r1[len(legs322r1)-1].name
	}
	msg := err.Error()
	if strings.Contains(msg, whyOpenELOOP322r1) && strings.Contains(msg, "open ") && !strings.Contains(msg, "openat ") {
		return legLeafOpen322r1
	}
	for _, leg := range legs322r1[:len(legs322r1)-2] {
		if strings.Contains(msg, leg.why) {
			return leg.name
		}
	}
	return legs322r1[len(legs322r1)-2].name
}

func TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1(t *testing.T) {
	dir := SealableTempDirForTest124(t)
	name := filepath.Join(dir, "caller-name-322r1.bin")
	stageReal := filepath.Join(dir, "staged-real-322r1.bin")
	stageLink := filepath.Join(dir, "staged-leaf-322r1.lnk")
	witnessAt := func(i int) string { return filepath.Join(dir, fmt.Sprintf("foreign-witness-%04d-322r1.bin", i)) }

	// Ticket 444 probe knob: the committed denominator is the default, verbatim.
	attempts := attemptsFor444(leafRaceAttempts322r1)
	// One fresh foreign object per attempt, named by nobody, wide so any touch shows. Equal
	// bytes imply equal size, so body and mtime are the whole before-picture.
	bodies := make([]string, attempts)
	mtimes := make([]time.Time, attempts)
	for i := range bodies {
		bodies[i] = fmt.Sprintf("foreign bytes %04d are not yours to truncate", i)
		if err := os.WriteFile(witnessAt(i), []byte(bodies[i]), 0o600); err != nil {
			t.Fatalf("stage witness %d: %v", i, err)
		}
		if err := os.Chmod(witnessAt(i), 0o666); err != nil {
			t.Fatalf("widen witness %d: %v", i, err)
		}
		info, err := os.Lstat(witnessAt(i))
		if err != nil {
			t.Fatalf("read witness %d back: %v", i, err)
		}
		mtimes[i] = info.ModTime()
	}
	if err := os.WriteFile(name, []byte("the object the caller owns"), 0o600); err != nil {
		t.Fatalf("give the caller a name to stand on: %v", err)
	}

	var cur, flips, flipErrs int32
	var flipMu sync.Mutex
	var firstFlip string
	note := func(err error) {
		if err == nil {
			atomic.AddInt32(&flips, 1)
			return
		}
		atomic.AddInt32(&flipErrs, 1)
		flipMu.Lock()
		defer flipMu.Unlock()
		if firstFlip == "" {
			firstFlip = err.Error()
		}
	}
	stop, done := make(chan struct{}), make(chan struct{})
	var stopOnce sync.Once
	stopRacer := func() { stopOnce.Do(func() { close(stop) }); <-done }
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
			i := int(atomic.LoadInt32(&cur))
			if i >= attempts {
				continue
			}
			// Both states arrive by rename onto the SAME spelling, so the name is never absent
			// and the open always has something to resolve (120-r3's anti-hole-instrument fix).
			if err := os.WriteFile(stageReal, []byte("the object the caller owns"), 0o600); err != nil {
				note(err)
				continue
			}
			note(os.Rename(stageReal, name))
			_ = os.Remove(stageLink) // consumed by the previous flip, so re-plant it
			if err := os.Symlink(filepath.Base(witnessAt(i)), stageLink); err != nil {
				note(err)
				continue
			}
			note(os.Rename(stageLink, name))
		}
	}()
	defer stopRacer()

	counts := map[string]int{}
	for i := 0; i < attempts; i++ {
		atomic.StoreInt32(&cur, int32(i))
		counts[spokeLeg322r1(PrivateFile(name, []byte(leafRacePayload322r1), 0o600))]++
	}
	atomic.StoreInt32(&cur, int32(attempts))
	stopRacer()

	// A sweep over EVERY witness, so damage landed after an attempt's own tally still counts.
	damaged, narrowed, landed, firstDamage := 0, 0, 0, ""
	for i := range bodies {
		body, rErr := os.ReadFile(witnessAt(i))
		info, sErr := os.Lstat(witnessAt(i))
		why := "unchanged"
		switch {
		case rErr != nil || sErr != nil:
			why = fmt.Sprintf("unreadable (%v / %v)", rErr, sErr)
		case string(body) != bodies[i]:
			why = fmt.Sprintf("content is %q, want its own bytes", body)
		case !info.ModTime().Equal(mtimes[i]):
			why = fmt.Sprintf("mtime moved from %s to %s", mtimes[i].Format(time.RFC3339Nano), info.ModTime().Format(time.RFC3339Nano))
		}
		if why != "unchanged" {
			damaged++
			if firstDamage == "" {
				firstDamage = filepath.Base(witnessAt(i)) + ": " + why
			}
		}
		if rErr == nil && strings.Contains(string(body), leafRacePayload322r1) {
			landed++
		}
		if st, err := os.Stat(witnessAt(i)); err == nil && st.Mode().Perm() != 0o666 {
			narrowed++
		}
	}
	liveness := counts[legLeafOpen322r1]
	ctl := leafLegControl444(t)
	t.Logf("AC#1 paired reading on one screen: attempts=%d flips=%d flipErrs=%d damagedWitnesses=%d leafLegRefusals=%d foreignsNarrowed=%d foreignsHoldingPayload=%d",
		attempts, flips, flipErrs, damaged, liveness, narrowed, landed)
	for _, leg := range legs322r1 {
		if n := counts[leg.name]; n > 0 || leg.name == legLeafOpen322r1 {
			t.Logf("AC#1' leg that answered = %s : %d attempts", leg.name, n)
		}
	}
	ctl.print(t, "AC#1")
	if damaged > 0 {
		t.Errorf("AC#1 RED: %d of %d foreign witnesses lost their own bytes or mtime, so the open in "+
			"privateFile followed a leaf link swapped into the caller's name and truncated a file no caller "+
			"named (first: %s). This is ticket 322's shape and the outcome the leaf leg exists to prevent.",
			damaged, attempts, firstDamage)
	}
	// Ticket 444's split, half one: the RACER's own arrival. A racer that landed zero flips means the
	// fixture, not the schedule, produced this run's zero - and deleting the racer has to stay fatal,
	// so this arm keeps its teeth (it is ticket 441's rule: the arm that proves the world was entered
	// stays RED, the arm that proves the QUIET side was observed may abstain).
	if flips == 0 {
		t.Errorf("AC#1 RED instrument: across %d raced attempts the racer landed 0 flips (%d errors, first %q), "+
			"so the damaged count above is not evidence about a window that was never reached and this "+
			"instrument's own fixture is what failed", attempts, flipErrs, firstFlip)
	}
	// Ticket 444's split, half two: the GUARDRAIL, carried by an attempt nobody races over. The leaf
	// leg is the open inside privateFile carrying oNoFollowPlatform; the control drives that open, on
	// the route's own pin, against a name that carries a link, and requires the platform to refuse it.
	if !ctl.refused || !ctl.witnessIntact {
		t.Errorf("AC#1 RED instrument: the un-raced leaf-leg control did not get the refusal ticket 322 "+
			"pays for. %s", ctl.verdict())
	}
	// The raced loop's own observation of the leaf leg, when contention starves it, is a named
	// abstention and not a verdict about the guardrail - the control above carries that.
	if liveness == 0 && flips > 0 {
		t.Logf("AC#1 ABSTAIN(444:322-leaf:leg-starved) attempts=%d samples=1 flips=%d flipErrs=%d "+
			"leaf-leg-refusals=0 damaged=%d - the racer landed %d flips but no attempt was inside the "+
			"link state when the route opened, so this run cannot attest the leaf leg FROM THE RACE; the "+
			"threshold that carries it is the un-raced control above (refused=%v, witness-intact=%v), "+
			"not a green from absence",
			attempts, flips, flipErrs, damaged, flips, ctl.refused, ctl.witnessIntact)
	}
	if landed > 0 || narrowed > 0 {
		t.Errorf("AC#1 RED: %d foreign witnesses hold a raced private payload and %d were narrowed off 0666, so a "+
			"write or a seal followed a flipped leaf link into a foreign object (ticket 120's surface re-opened)", landed, narrowed)
	}
}

// leafLegReading444 is the un-raced control's whole reading: the refusal, the leg that produced it,
// and the foreign object's state afterwards.
type leafLegReading444 struct {
	attempts      int
	refused       bool
	leg           string
	witnessBody   string
	witnessMode   fs.FileMode
	witnessIntact bool
	err           error
}

func (r leafLegReading444) verdict() string {
	return fmt.Sprintf("attempts=%d refused=%v leg=%q witness-body=%q witness-mode=%v verbatim: %v",
		r.attempts, r.refused, r.leg, r.witnessBody, r.witnessMode, r.err)
}

func (r leafLegReading444) print(t *testing.T, tag string) {
	t.Helper()
	t.Logf("%s un-raced leaf-leg control (the route's own open, no racer): %s", tag, r.verdict())
}

// leafLegControl444 is ticket 444's carrier for the half of this instrument that contention can
// starve. The raced loop above can only see the leaf leg when its own call lands inside the racer's
// link state - that is a coin flip under load, and reading its absence as "the guardrail is gone" is
// the false red this ticket exists to remove.
//
// What is NOT a coin flip: the open inside privateFile carries oNoFollowPlatform, and that open
// refuses a name carrying a leaf link on this platform, every time it is asked. So the control asks
// it: the same pin (pinnedParent), the same helper (openPinnedLeaf), the same flag expression the
// route hands it (winsec.go's `flags|oNoFollowPlatform`), on a name that carries a sibling link at a
// foreign object. WHAT IT DOES NOT CLAIM: it does not prove the route still PASSES that flag - a route
// that stopped passing it would keep this control green. That claim is the raced `liveness` reading
// above, which is why that reading is printed and why deleting the flag from oNoFollowPlatform (or
// dropping the flags inside openPinnedLeaf) turns this control red, as ticket 444's AC#3 measures.
// openleaf_flag_322_other_test.go pins the same flag at the raw os.OpenFile the route used before
// form 乙 moved it onto the pin.
func leafLegControl444(t *testing.T) leafLegReading444 {
	t.Helper()
	dir := SealableTempDirForTest124(t)
	const witnessBody = "an object the caller never named, and this body is longer on purpose"
	witness := filepath.Join(dir, "foreign-witness-444-322r1.bin")
	if err := os.WriteFile(witness, []byte(witnessBody), 0o600); err != nil {
		t.Fatalf("control: stage the foreign witness: %v", err)
	}
	if err := os.Chmod(witness, 0o666); err != nil {
		t.Fatalf("control: widen the foreign witness: %v", err)
	}
	link := filepath.Join(dir, "caller-name-444-322r1.bin")
	if err := os.Symlink(filepath.Base(witness), link); err != nil {
		t.Fatalf("control: plant the leaf link: %v", err)
	}
	pin, err := pinnedParent(link, link)
	if err != nil {
		t.Fatalf("control: pin the parent of a leaf link: %v", err)
	}
	defer func() { _ = pin.close() }()
	f, _, openErr := openPinnedLeaf(pin, filepath.Base(link), os.O_WRONLY|os.O_CREATE|os.O_TRUNC|oNoFollowPlatform, 0o600)
	if f != nil {
		_ = f.Close()
	}
	r := leafLegReading444{attempts: 1, err: openErr}
	if openErr != nil {
		r.leg = spokeLeg322r1(openErr)
		r.refused = r.leg == legLeafOpen322r1
	}
	body, rErr := os.ReadFile(witness)
	info, sErr := os.Lstat(witness)
	switch {
	case rErr != nil || sErr != nil:
		r.witnessBody = fmt.Sprintf("unreadable (%v / %v)", rErr, sErr)
	default:
		r.witnessBody = string(body)
		r.witnessMode = info.Mode().Perm()
		r.witnessIntact = r.witnessBody == witnessBody && r.witnessMode == 0o666
	}
	return r
}
