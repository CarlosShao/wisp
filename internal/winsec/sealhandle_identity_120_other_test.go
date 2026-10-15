//go:build !windows

// Ticket 120, residual gap 1: sealHandle's OWN two inline re-checks
// (winsec_other.go:269 the link face, :272 the SameFile face) had no case that
// reaches them. The acceptance leg 120-v1 measured that by taking the teeth out:
// both conditions rewritten to `false && ...` on a snapshot, whole package run,
// RUN=55 PASS=55 FAIL=0 SKIP=0, not one reading changed colour, so this half of
// the gate had a shape and no evidence (docs/evidence/s1/120-v1.md section 2).
//
// Why nothing reached it is read from the code here and re-measured by
// TestAC1POSIXSealHandleNotReachedByAPrePlantedLink120r3 below: every existing
// case plants the link on the name BEFORE the call, so sealHandle's first
// statement - applyDescriptor, which walks and pins by name - already refuses at
// sameObjectPinned (winsec_other.go:188). sealHandle returns at :258 and its own
// :269 never runs. The two legs are two legs: 120-r2's control neutered
// sameObjectPinned and turned the ticket's own shape 3 red, and that reading says
// nothing at all about this one.
//
// The only world that reaches :269 and :272 is the one sealHandle's doc comment
// describes and no fixture had built: the CALLER already holds an open descriptor
// on object A, and the name changes after the pinned seal has run and before the
// identity re-check reads it. Both shapes are built here, separately:
//
//   - the SameFile face (:272): rename(2) a second REAL object onto the name
//     before the call. Nothing is a link at any instant, so the pinned walk is
//     satisfied by construction and only sealHandle's own comparison can refuse.
//     Fully deterministic, no fixture trick; the blind racer in
//     sealhandle_privateroute_120_other_test.go hits this same state on its own.
//   - the link face (:269): the name has to carry a link at :265 while the walk
//     that ran a handful of syscalls earlier saw a plain file. A link planted in
//     advance cannot produce that world - it is refused by the walk, which is the
//     gap's cause above - so the plant has to land IN the window. It lands through
//     the package's own declared seam: applyDescriptor is a var on purpose
//     (winsec_other.go:14-23) and private_fail_test.go already replaces it for
//     AC#5. This fixture does NOT replace it: it calls applyDescriptorPOSIX,
//     keeps that verdict, and only then performs the move an attacker would have
//     performed. The leg under test still sees the real pinned seal's answer; the
//     fixture decides WHEN the swap happens, which is the one thing a blind racer
//     cannot be promised. Its blind distribution is reported, not asserted,
//     because 120-v1 section 2's third ruling says a count that cannot be
//     promised must not become a gate.
//
// What the two cases assert beyond the refusal, and why: each one reads BOTH
// objects afterwards. The pinned walk narrows whatever the name carried when it
// walked, and the caller's descriptor still stands on the object the write would
// go to, so "the object that was narrowed is not the object that would be written"
// is a fact on disk, not a restatement of the error string. That is also what makes
// the control honest: with :272 neutered the same read-back shows the caller's own
// object being narrowed by a seal that was answering about a different object, and
// with :269 neutered the call returns nil while the name carries a link.
//
// Control plan (the pairing this file exists for): neuter :269 alone and the only
// readings that change colour are the link face here and its twin on the
// PrivateFile route; neuter :272 alone and only the SameFile face and its route
// twin change colour. Both neuters are physical rewrites of winsec_other.go plus a
// sha256 restore, on the same container table as the paired green run
// (.scratch/wisp/probes/120/r3/). -overlay is not used and cannot be: the guard
// under test reads the disk at run time.
//
// The container runs as uid 0, so the permission faces of the pinned walk are not
// asserted anywhere here (ticket 113's standing note, ticket 120 r2 measured them
// at uid 1000 instead). Mode bits are stored on the container's overlay /tmp, which
// is what makes the read-backs below mean something.
package winsec

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

// The four refusal reasons of winsec_other.go, transcribed as they are formatted
// by placementRefusal so a reading can name WHICH leg spoke. Two belong to the
// pinned walk (sameObjectPinned), two to sealHandle's own re-check; the pairs are
// deliberately near-identical in wording, which is exactly why matching on a
// sentinel alone cannot tell the legs apart and why these constants are spelled out.
const (
	whyWalkLink120r3   = "the name carries a link, so the pinned object is not what the spelling reaches"
	whyWalkSwap120r3   = "the name carries a different object than the one that was pinned"
	whyHandleLink120r3 = "the name carries a link, so this descriptor is not the object the caller named"
	whyHandleSwap120r3 = "this descriptor is not the object the name carries now, so the bytes would not " +
		"land in what was just sealed"
)

// spokeLeg120r3 names the leg an error came from, for the log lines and for the
// "and not some other leg" half of the judgement.
func spokeLeg120r3(err error) string {
	if err == nil {
		return "none, the call returned nil"
	}
	msg := err.Error()
	for _, leg := range []struct{ why, name string }{
		{whyWalkLink120r3, "pinPlacement sameObjectPinned, link face (winsec_other.go:188)"},
		{whyWalkSwap120r3, "pinPlacement sameObjectPinned, SameFile face (winsec_other.go:191)"},
		{whyHandleLink120r3, "sealHandle inline re-check, link face (winsec_other.go:269)"},
		{whyHandleSwap120r3, "sealHandle inline re-check, SameFile face (winsec_other.go:272)"},
	} {
		if strings.Contains(msg, leg.why) {
			return leg.name
		}
	}
	return "some other failure"
}

// assertLegSpoke120r3 is the judgement every case below shares: a refusal, carrying
// both sentinels placementRefusal promises, coming from exactly ONE named leg.
func assertLegSpoke120r3(t *testing.T, err error, want, tag string) {
	t.Helper()
	t.Logf("AC#1' %s: leg that spoke = %s", tag, spokeLeg120r3(err))
	if err == nil {
		t.Errorf("AC#1 RED %s: the call returned nil, so %s never spoke; want the refusal %q "+
			"(this is the shape of taking that leg's teeth out)", tag, spokeLeg120r3(err), want)
		return
	}
	if !errors.Is(err, ErrNotSealable) {
		t.Errorf("AC#1 RED %s: refusal does not carry ErrNotSealable: %v", tag, err)
	}
	if !errors.Is(err, ErrUnresolvedPath) {
		t.Errorf("AC#1 RED %s: refusal does not carry ErrUnresolvedPath: %v", tag, err)
	}
	if !strings.Contains(err.Error(), want) {
		t.Errorf("AC#1 RED %s: refusal came from %s, want the leg carrying %q; verbatim: %v",
			tag, spokeLeg120r3(err), want, err)
	}
}

// sealHandleFixture120r3 builds the tree the caller names: one regular WIDE file at
// the spelling, and a sibling regular WIDE file standing where a planted link would
// point (the witness). Both are 0666 on purpose - "still wide" is what makes a
// seal that landed somewhere visible as a mode change, ticket 113's instrument.
func sealHandleFixture120r3(t *testing.T) (dir, name, witness string) {
	t.Helper()
	dir = SealableTempDirForTest124(t)
	name = filepath.Join(dir, "caller-object-120r3.bin")
	witness = filepath.Join(dir, "witness-sibling-120r3.bin")
	writeWide120r3(t, name, "the object the caller opened")
	writeWide120r3(t, witness, "not this write's object")
	return dir, name, witness
}

func writeWide120r3(t *testing.T, path, body string) {
	t.Helper()
	// Created at 0600 then widened: the widening has to be the stored answer of the
	// filesystem, so a later read-back of 0600 can only mean a seal changed it.
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("stage %s: %v", filepath.Base(path), err)
	}
	if err := os.Chmod(path, 0o666); err != nil {
		t.Fatalf("widen %s: %v", filepath.Base(path), err)
	}
}

// openHandleOn120r3 hands back the fd the caller holds. privateFile opens with
// O_WRONLY|O_CREATE|O_TRUNC; for the identity question both descriptors stand on the
// same object, and O_RDONLY|O_NONBLOCK is the open pinPlacement itself uses
// (winsec_other.go:158), so nothing here borrows a privilege the leg does not have.
func openHandleOn120r3(t *testing.T, path string) *os.File {
	t.Helper()
	f, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		t.Fatalf("open the object the caller named: %v", err)
	}
	t.Cleanup(func() { _ = f.Close() })
	return f
}

// inIdentityWindow120r3 is the fixture's one trick, and it is a window rather than a
// substitution: the REAL pinned seal runs first and keeps its verdict; the swap only
// happens after that verdict is nil, i.e. in the gap sealHandle's own re-check was
// written to close.
func inIdentityWindow120r3(t *testing.T, plant func()) {
	t.Helper()
	orig := applyDescriptor
	applyDescriptor = func(path string, dir bool) error {
		err := applyDescriptorPOSIX(path, dir)
		if err == nil {
			plant()
		}
		return err
	}
	t.Cleanup(func() { applyDescriptor = orig })
}

// modeOfName120r3 follows the spelling the way a reader would; modeOfFd120r3 reads
// the object the descriptor stands on, name or no name.
func modeOfName120r3(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return "unreachable"
	}
	return info.Mode().String()
}

func modeOfFd120r3(f *os.File) string {
	info, err := f.Stat()
	if err != nil {
		return "unreadable"
	}
	return info.Mode().String()
}

// kindOfName120r3 is what the name carries right now, without following it.
func kindOfName120r3(path string) string {
	info, err := os.Lstat(path)
	if err != nil {
		return "absent"
	}
	switch {
	case info.Mode()&os.ModeSymlink != 0:
		return "link"
	case info.IsDir():
		return "dir"
	default:
		return "file"
	}
}

func identity120r3(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Sprintf("unreachable(%v)", err)
	}
	if st, ok := info.Sys().(*syscall.Stat_t); ok {
		return fmt.Sprintf("dev=%d ino=%d", st.Dev, st.Ino)
	}
	return "no platform identity"
}

func fdIdentity120r3(f *os.File) string {
	info, err := f.Stat()
	if err != nil {
		return fmt.Sprintf("unreadable(%v)", err)
	}
	if st, ok := info.Sys().(*syscall.Stat_t); ok {
		return fmt.Sprintf("dev=%d ino=%d", st.Dev, st.Ino)
	}
	return "no platform identity"
}

// TestAC1POSIXSealHandleNameCarriesAnotherRealObject120r3 is the SameFile face,
// built with no link standing anywhere: the caller's fd is on object A, a second
// real object is renamed onto A's spelling before the call, so the pinned walk
// finds exactly what the name carries and is satisfied, and the only thing left to
// notice that the write would go somewhere else is sealHandle's own comparison.
func TestAC1POSIXSealHandleNameCarriesAnotherRealObject120r3(t *testing.T) {
	dir, name, _ := sealHandleFixture120r3(t)
	f := openHandleOn120r3(t, name)
	fdBefore := fdIdentity120r3(f)

	swapIn := filepath.Join(dir, "swap-in-120r3.bin")
	writeWide120r3(t, swapIn, "the object swapped into the name")
	if err := os.Rename(swapIn, name); err != nil {
		t.Fatalf("swap a second real object onto the name: %v", err)
	}

	// Fixture liveness, stated before the call so a world that was not built cannot
	// be read as a leg that refused everything: the name and the fd must already be
	// two objects here.
	if os.SameFile(mustStat120r3(t, name), mustFdStat120r3(t, f)) {
		t.Fatalf("fixture did not build the world: the name still carries the fd's object (%s)", fdBefore)
	}
	t.Logf("AC#1' world: fd stands on %s, the name %s now carries %s, no link at any instant",
		fdBefore, name, identity120r3(name))

	err := sealHandle(f)
	assertLegSpoke120r3(t, err, whyHandleSwap120r3, "name-carries-another-real-object")

	// The two read-backs are the reason this leg exists, measured rather than
	// quoted: the walk narrowed the object the NAME denotes, which is not the object
	// the pending write stands on.
	t.Logf("AC#1' read-back: object the name carries = %s, caller's object = %s",
		modeOfName120r3(name), modeOfFd120r3(f))
	if got := modeOfName120r3(name); got != "-rw-------" {
		t.Errorf("AC#1 RED: the pinned walk left the object the name carries at %s, want -rw------- "+
			"(so this run never reached the identity re-check at all)", got)
	}
	if got := modeOfFd120r3(f); got != "-rw-rw-rw-" {
		t.Errorf("AC#1 RED: the caller's own object was narrowed to %s while the refusal says the write "+
			"would not land in what was just sealed - the leg at winsec_other.go:272 is not speaking", got)
	}
}

// TestAC1POSIXSealHandleNameBecomesLinkAfterPinnedSeal120r3 is the link face, and
// the only construction that reaches it: the pinned seal runs for real on the plain
// file, the name is then taken away and replaced by a RELATIVE SIBLING link at the
// witness, and sealHandle's own Lstat is the first thing to see it.
func TestAC1POSIXSealHandleNameBecomesLinkAfterPinnedSeal120r3(t *testing.T) {
	_, name, witness := sealHandleFixture120r3(t)
	f := openHandleOn120r3(t, name)

	inIdentityWindow120r3(t, func() {
		// The attacker's move, landed in the window instead of hoped for. Relative
		// sibling target: no "..", no absolute spelling, so this is the one link
		// shape the walk itself would follow if it ever saw it.
		if err := os.Remove(name); err != nil {
			t.Errorf("fixture: take the name away: %v", err)
			return
		}
		if err := os.Symlink(filepath.Base(witness), name); err != nil {
			t.Errorf("fixture: plant the sibling link: %v", err)
			return
		}
	})

	err := sealHandle(f)
	assertLegSpoke120r3(t, err, whyHandleLink120r3, "name-becomes-a-sibling-link-after-the-pinned-seal")

	// Liveness of BOTH halves: the plant really happened, and the pinned seal really
	// ran for the caller's object (it is the object that came back narrowed).
	t.Logf("AC#1' read-back: the name carries a %s at %s, object the fd holds = %s, witness = %s",
		kindOfName120r3(name), identity120r3(witness), fdIdentity120r3(f), modeOfName120r3(witness))
	if got := kindOfName120r3(name); got != "link" {
		t.Errorf("AC#1 RED: the name carries a %s, not the planted link, so this run did not build the world", got)
	}
	if got := modeOfFd120r3(f); got != "-rw-------" {
		t.Errorf("AC#1 RED: the caller's object is still %s, so applyDescriptor did not run for real "+
			"before the plant - the seam was substituted, not windowed", got)
	}
	if got := modeOfName120r3(witness); got != "-rw-rw-rw-" {
		t.Errorf("AC#1 RED: the witness the link points at was narrowed to %s, so a seal followed it", got)
	}
}

// TestAC1POSIXSealHandleNotReachedByAPrePlantedLink120r3 reproduces the gap's cause
// rather than asserting it: the same link, planted before the call the way every
// existing PrivateFile case does it, is refused by the pinned WALK, and the reading
// names which leg spoke. This is the case that shows :269 has no world of its own
// unless the swap lands in the window - and it is deliberately not the case the
// :269 control turns red, because it never reaches that line.
func TestAC1POSIXSealHandleNotReachedByAPrePlantedLink120r3(t *testing.T) {
	_, name, witness := sealHandleFixture120r3(t)
	if err := os.Remove(name); err != nil {
		t.Fatalf("fixture: %v", err)
	}
	if err := os.Symlink(filepath.Base(witness), name); err != nil {
		t.Fatalf("fixture: plant the link before the call: %v", err)
	}
	f := openHandleOn120r3(t, name)
	t.Logf("AC#1' world: the name carries a %s pointing at a relative sibling; the fd stands on %s",
		kindOfName120r3(name), fdIdentity120r3(f))

	err := sealHandle(f)
	assertLegSpoke120r3(t, err, whyWalkLink120r3, "pre-planted-link-is-refused-by-the-walk")
	if strings.Contains(err.Error(), whyHandleLink120r3) {
		t.Errorf("AC#1 RED: the reading names sealHandle's own link face, so this case is not the gap's cause")
	}
	if got := modeOfName120r3(witness); got != "-rw-rw-rw-" {
		t.Errorf("AC#1 RED: the witness was narrowed to %s through the planted link", got)
	}
}

// TestAC1POSIXSealHandlePlainObjectStillSeals120r3 is the other direction: a leg
// that refused everything would pass the two cases above and be worthless. With the
// name carrying the caller's own object, sealHandle must return nil and land the
// mode. It is also the case that must NOT change colour under either control.
func TestAC1POSIXSealHandlePlainObjectStillSeals120r3(t *testing.T) {
	_, name, witness := sealHandleFixture120r3(t)
	f := openHandleOn120r3(t, name)
	if err := sealHandle(f); err != nil {
		t.Errorf("AC#1 RED: the honest path was refused: %v (leg that spoke = %s)", err, spokeLeg120r3(err))
	}
	if got := modeOfName120r3(name); got != "-rw-------" {
		t.Errorf("AC#1 RED: the object the caller named is at %s, want -rw-------", got)
	}
	if got := modeOfName120r3(witness); got != "-rw-rw-rw-" {
		t.Errorf("AC#1 RED: the witness was narrowed to %s on a path with no link in it", got)
	}
}

func mustStat120r3(t *testing.T, path string) fs.FileInfo {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	return info
}

func mustFdStat120r3(t *testing.T, f *os.File) fs.FileInfo {
	t.Helper()
	info, err := f.Stat()
	if err != nil {
		t.Fatalf("fstat: %v", err)
	}
	return info
}
