//go:build !windows

// Ticket 120, residual gap 1, shared apparatus. Leg 120-r3, test surface only:
// nothing in here changes a product line, and winsec_other.go is byte-identical
// before and after (its sha256 is on the record in .scratch/wisp/probes/120/r3/).
//
// THE GAP. sealHandle (winsec_other.go:256-276) consults applyDescriptor first and
// then adds the question a seal BY NAME cannot ask: is the descriptor the caller is
// about to write through the object that was just narrowed? Its two inline
// re-checks - the link face at :269 and the SameFile face at :272 - had no case
// reaching them. The acceptance leg measured that by taking the teeth out: both
// conditions rewritten to `false && ...` physically on a snapshot, whole package
// run, RUN=55 PASS=55 FAIL=0 SKIP=0, not one reading changed colour
// (docs/evidence/s1/120-v1.md section 2). A shape with no way to be made to speak
// is not evidence, so this file builds the world.
//
// WHY NOTHING REACHED IT, re-measured rather than quoted (see
// TestAC1POSIXSealHandleNotReachedByAPrePlantedLink120r3): every existing case
// plants the link on the name BEFORE the call, so sealHandle's own first statement -
// applyDescriptor, which walks from the volume root and pins every component -
// refuses at sameObjectPinned (:188) and the call returns at :258. That is a
// different leg: 120-r2's control neutered sameObjectPinned and turned the ticket's
// own shape 3 red, which says nothing about :269 or :272.
//
// THE ONLY WORLD THAT SPEAKS. The caller already holds an open descriptor on object
// A, and the name changes AFTER the pinned seal has run and BEFORE the identity
// re-check reads it. Two shapes, and they need two different constructions:
//
//   - SameFile face (:272): rename(2) a second REAL object onto the name. Nothing is
//     a link at any instant, so the walk is satisfied by construction and only
//     sealHandle's own comparison can refuse. Deterministic, no fixture trick; the
//     blind racer in sealhandle_privateroute_120_other_test.go hits it on its own.
//   - link face (:269): the name has to carry a link at :265 while the walk a handful
//     of syscalls earlier saw a plain file. Pre-planting it produces the world above
//     (the walk refuses), so the plant must land IN the window. It does, through the
//     package's own declared seam var: applyDescriptor is a var ON PURPOSE
//     (winsec_other.go:14-23) and private_fail_test.go already replaces it for AC#5.
//     inIdentityWindow120r3 does NOT replace it - it calls applyDescriptorPOSIX,
//     keeps that verdict, and only then performs the move an attacker would have
//     made. The leg under test still gets the real pinned seal's answer; the fixture
//     decides WHEN the swap happens, which is the one thing a blind racer cannot be
//     promised. The blind distribution is therefore REPORTED and not gated, per
//     120-v1 section 2's ruling that a landing count nobody can promise must not
//     become a hard assertion.
//
// Each case also reads BOTH objects back after the call. The walk narrows whatever
// the name carried when it walked, while the caller's descriptor still stands on the
// object the write goes to, so "the narrowed object is not the written object" is a
// fact on disk rather than a restatement of an error string - and it is the half that
// changes colour when a tooth is removed, independent of the refusal itself.
//
// The container runs as uid 0, so no permission face is asserted anywhere in this
// family (ticket 113's standing note; ticket 120 r2 measured those at uid 1000).
// Mode bits are stored on the container's overlay /tmp, which is what makes a
// read-back of 0600 mean a seal happened.
package winsec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

// The four refusal reasons winsec_other.go formats through placementRefusal,
// transcribed so a reading can name WHICH leg spoke. The two pairs are near-identical
// in wording on purpose in the product, which is exactly why matching a sentinel
// alone cannot tell the legs apart.
const (
	whyWalkLink120r3   = "the name carries a link, so the pinned object is not what the spelling reaches"
	whyWalkSwap120r3   = "the name carries a different object than the one that was pinned"
	whyHandleLink120r3 = "the name carries a link, so this descriptor is not the object the caller named"
	whyHandleSwap120r3 = "this descriptor is not the object the name carries now, so the bytes would not " +
		"land in what was just sealed"
)

// spokeLeg120r3 names the leg an error came from, for the log rows and for the
// "and not some other leg" half of every judgement below.
func spokeLeg120r3(err error) string {
	if err == nil {
		return "none: the call returned nil"
	}
	msg := err.Error()
	for _, leg := range []struct{ why, name string }{
		{whyWalkLink120r3, "pinPlacement sameObjectPinned link face (winsec_other.go:188)"},
		{whyWalkSwap120r3, "pinPlacement sameObjectPinned SameFile face (winsec_other.go:191)"},
		{whyHandleLink120r3, "sealHandle inline link face (winsec_other.go:269)"},
		{whyHandleSwap120r3, "sealHandle inline SameFile face (winsec_other.go:272)"},
		{"reaches it through the link at", "floor before the act (winsec_other.go platformVerifyPlacement)"},
	} {
		if strings.Contains(msg, leg.why) {
			return leg.name
		}
	}
	return "some other failure"
}

// assertLegSpoke120r3 is the judgement every case shares: a refusal, carrying both
// sentinels placementRefusal promises, coming from exactly the named leg.
func assertLegSpoke120r3(t *testing.T, err error, want, tag string) {
	t.Helper()
	t.Logf("AC#1' %s: leg that spoke = %s", tag, spokeLeg120r3(err))
	if err == nil {
		t.Errorf("AC#1 RED %s: the call returned nil, so the leg carrying %q never spoke - "+
			"this is the shape of that leg with its teeth taken out", tag, want)
		return
	}
	if !errors.Is(err, ErrNotSealable) {
		t.Errorf("AC#1 RED %s: the refusal does not carry ErrNotSealable: %v", tag, err)
	}
	if !errors.Is(err, ErrUnresolvedPath) {
		t.Errorf("AC#1 RED %s: the refusal does not carry ErrUnresolvedPath: %v", tag, err)
	}
	if !strings.Contains(err.Error(), want) {
		t.Errorf("AC#1 RED %s: refusal came from %s, want the leg carrying %q; verbatim: %v",
			tag, spokeLeg120r3(err), want, err)
	}
}

// sealHandleFixture120r3 builds the tree the caller names: a regular WIDE file at the
// spelling, and a sibling regular WIDE file standing where a planted link points (the
// witness). Both 0666 on purpose - "still wide" is what makes a seal that wandered
// visible as a mode change (ticket 113's instrument).
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
	// Created then widened: the wide mode has to be the filesystem's stored answer,
	// so a later read-back of 0600 can only mean something changed it.
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("stage %s: %v", filepath.Base(path), err)
	}
	if err := os.Chmod(path, 0o666); err != nil {
		t.Fatalf("widen %s: %v", filepath.Base(path), err)
	}
}

// openHandleOn120r3 hands back the descriptor the caller holds. privateFile opens
// with O_WRONLY|O_CREATE|O_TRUNC; for the identity question both stand on the same
// object, and O_RDONLY|O_NONBLOCK is the open pinPlacement itself uses (:158), so
// nothing here borrows a privilege the leg under test does not have.
func openHandleOn120r3(t *testing.T, path string) *os.File {
	t.Helper()
	f, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		t.Fatalf("open the object the caller named: %v", err)
	}
	t.Cleanup(func() { _ = f.Close() })
	return f
}

// inIdentityWindow120r3 is the fixture's one trick, and it is a window, not a
// substitution: the REAL pinned seal runs first and keeps its verdict, and the swap
// happens only after that verdict came back nil - i.e. in the gap sealHandle's own
// re-check exists to close.
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

// The three questions a reading can ask of a name: what it carries without following
// it, what the object behind it is at, and whether it is the object a descriptor
// stands on. Kept separate because the legs disagree about exactly these.
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

func nameCarriesFdObject120r3(f *os.File, path string) string {
	fdInfo, err := f.Stat()
	if err != nil {
		return fmt.Sprintf("fd unreadable (%v)", err)
	}
	nameInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Sprintf("name unreachable (%v)", err)
	}
	if os.SameFile(fdInfo, nameInfo) {
		return "same object"
	}
	return "two different objects"
}
