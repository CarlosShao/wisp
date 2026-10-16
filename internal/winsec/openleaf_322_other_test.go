//go:build !windows

// Ticket 322, leg 322-r1. Judged here: the OPEN inside privateFile
// (os.OpenFile(resolved, flags|oNoFollowPlatform, perm)), which runs before any seal. The
// shape is adjudicated in ledger A996, not re-decided here: O_NOFOLLOW closes the LEAF leg on
// POSIX and the ANCESTOR leg stays open and stays named - it speaks about the last component
// only, so an ancestor swapped for a link after resolveString answered is still followed and
// its O_TRUNC still lands in somebody else's tree. Nothing here claims the window is closed;
// reasoning, top-red roster and paired readings: docs/evidence/s1/322-open-nofollow.md. Why a
// racer: a leaf link planted before the call is refused by the placement floor before the open
// runs, so planting cannot reach this leg (120-r3's reason for winsec_other.go:269). Judged: no
// fresh foreign witness loses its own bytes or mtime, AND the new leg spoke at least once, so a
// 0 numerator is never read as "nothing was raced". Both numbers print on one screen.
package winsec

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
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

	// One fresh foreign object per attempt, named by nobody, wide so any touch shows. Equal
	// bytes imply equal size, so body and mtime are the whole before-picture.
	bodies := make([]string, leafRaceAttempts322r1)
	mtimes := make([]time.Time, leafRaceAttempts322r1)
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
			if i >= leafRaceAttempts322r1 {
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
	for i := 0; i < leafRaceAttempts322r1; i++ {
		atomic.StoreInt32(&cur, int32(i))
		counts[spokeLeg322r1(PrivateFile(name, []byte(leafRacePayload322r1), 0o600))]++
	}
	atomic.StoreInt32(&cur, leafRaceAttempts322r1)
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
	t.Logf("AC#1 paired reading on one screen: attempts=%d flips=%d flipErrs=%d damagedWitnesses=%d leafLegRefusals=%d foreignsNarrowed=%d foreignsHoldingPayload=%d",
		leafRaceAttempts322r1, flips, flipErrs, damaged, liveness, narrowed, landed)
	for _, leg := range legs322r1 {
		if n := counts[leg.name]; n > 0 || leg.name == legLeafOpen322r1 {
			t.Logf("AC#1' leg that answered = %s : %d attempts", leg.name, n)
		}
	}
	if damaged > 0 {
		t.Errorf("AC#1 RED: %d of %d foreign witnesses lost their own bytes or mtime, so the open in "+
			"privateFile followed a leaf link swapped into the caller's name and truncated a file no caller "+
			"named (first: %s). This is ticket 322's shape and the outcome the leaf leg exists to prevent.",
			damaged, leafRaceAttempts322r1, firstDamage)
	}
	if liveness == 0 || flips == 0 {
		t.Errorf("AC#1 RED instrument: across %d raced attempts the leaf leg (%q) answered %d and the racer landed "+
			"%d flips (%d errors, first %q), so the damaged count above is not evidence about a window that was "+
			"never reached", leafRaceAttempts322r1, whyOpenELOOP322r1, liveness, flips, flipErrs, firstFlip)
	}
	if landed > 0 || narrowed > 0 {
		t.Errorf("AC#1 RED: %d foreign witnesses hold a raced private payload and %d were narrowed off 0666, so a "+
			"write or a seal followed a flipped leaf link into a foreign object (ticket 120's surface re-opened)", landed, narrowed)
	}
}

// TestAC322LeafFlagRefusesALeafLinkWithoutTouchingTheTarget322r1 pins the constant's kernel
// effect deterministically, including two shapes only this can show: a refusal must not
// delete the link it was told about (322-a1's nearest miss, the unlink on privateFile's
// seal-failure path), and a DANGLING leaf must stop creating an object outside the tree.
func TestAC322LeafFlagRefusesALeafLinkWithoutTouchingTheTarget322r1(t *testing.T) {
	dir, outside := SealableTempDirForTest124(t), SealableTempDirForTest124(t)
	witness := filepath.Join(outside, "someone-elses-file-322r1.bin")
	const body = "the bytes a truncating open would have taken away"
	if err := os.WriteFile(witness, []byte(body), 0o666); err != nil {
		t.Fatal(err)
	}
	before, err := os.Lstat(witness)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ label, link, target string }{
		{"relative sibling", filepath.Join(dir, "sib-322r1.lnk"), filepath.Base(witness)},
		{"dangling, target never created", filepath.Join(dir, "dangling-322r1.lnk"), "never-created-322r1.bin"},
	} {
		if err := os.Symlink(tc.target, tc.link); err != nil {
			t.Fatalf("%s: stage the leaf: %v", tc.label, err)
		}
		f, err := os.OpenFile(tc.link, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|oNoFollowPlatform, 0o600)
		if err == nil {
			_ = f.Close()
			t.Errorf("AC#1 RED %s: the open carrying oNoFollowPlatform took a leaf link, want ELOOP", tc.label)
		} else if !errors.Is(err, syscall.ELOOP) || !strings.Contains(err.Error(), whyOpenELOOP322r1) {
			t.Errorf("AC#1 RED %s: the open answered %v, want the leaf leg carrying %q", tc.label, err, whyOpenELOOP322r1)
		}
		if _, sErr := os.Lstat(tc.link); sErr != nil {
			t.Errorf("AC#1 RED %s: the refusal removed the link itself (%v), so a refusal deleted the entry it was told about", tc.label, sErr)
		}
	}
	if _, err := os.Lstat(filepath.Join(dir, "never-created-322r1.bin")); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("AC#1 RED: the dangling-leaf spelling created an object outside the tree")
	}
	after, err := os.Lstat(witness)
	if err != nil {
		t.Fatalf("the witness became unreadable: %v", err)
	}
	if got, rErr := os.ReadFile(witness); rErr != nil || string(got) != body {
		t.Errorf("AC#1 RED: the witness holds %q, want its own bytes untouched", got)
	}
	if after.Size() != before.Size() || !after.ModTime().Equal(before.ModTime()) {
		t.Errorf("AC#1 RED: the witness moved from size %d mtime %s to size %d mtime %s", before.Size(),
			before.ModTime().Format(time.RFC3339Nano), after.Size(), after.ModTime().Format(time.RFC3339Nano))
	}
}
