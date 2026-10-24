//go:build windows

// Ticket 322, leg 322-r1, AC#3 - the Windows half, measured on this machine.
//
// Shape is adjudicated in ledger A996 and is not re-decided here: on Windows
// O_NOFOLLOW has no spelling (oNoFollowPlatform is 0), so the leaf leg is answered by
// preOpenLeafGuard opening the name with FILE_FLAG_OPEN_REPARSE_POINT, asking that
// handle one question - is this name a link? - and refusing if it is. The handle is NOT
// a write path: measured by 322-a1, a Write or Truncate against a reparse handle
// "succeeds" without ever reaching the target, and the CREATE_ALWAYS variant of the same
// probe was measured replacing the LINK ENTRY with an empty regular file, which is a
// second damage rather than a fix. So this half closes a judgment, and two residuals
// stay open and stay named:
//
//   - the probe cannot hold its answer. Between the verdict and the os.OpenFile that
//     follows, the leaf can still be swapped and that open still follows it. The
//     truncating count is therefore REPORTED, not gated - 322-a1 measured the same
//     race at damaged=571/759/780 on this machine with no guard at all, so a drop with a
//     non-zero guard-leg count is what this shape can honestly claim;
//   - ancestors are not probed (the placement floor covers the steady state, not the gap
//     after it has answered).
//
// Also out of this ticket and named so nobody reads these rows as full Windows cover:
// sealHandle here still has no same-object re-check (winsec_windows.go), so "private
// bytes land in a foreign object" is ticket 120's open Windows leg.
//
// 〔仅本机可量〕 CI and the container cannot reach this file: the container is linux and
// the CI runner's symlink-privilege tier was never measured (322-a1 未测格 3). Where the
// construction itself needs a privilege this machine has, the case says so out loud with
// t.Skipf rather than passing quietly.
//
// 票 369 (leg 369-r1) REPLACED THIS FILE'S ROUTE CRITERION, and the reason is a
// mechanism rather than a flake. The old case ran 400 raced attempts and gated
// `guardSpoke > 0`. That count is not a property of the landing: privateFile calls
// resolveString first (winsec.go:163), and BOTH steady-state refusals it owns walk the
// LEAF as one of their prefixes - the installed resolver answers "path traverses a
// reparse point" (internal/risk, reached through resolve.go:605-611) and the floor's own
// walk answers "traverses a reparse point at" (placement_windows.go:59-68, because
// pathPieces' last piece IS the full path, winsec.go:408-420). So in the steady state a
// leaf link is refused BEFORE winsec.go:167 is reached, and preOpenLeafGuard can only
// answer when the swap lands inside the two syscalls between that walk and its own
// CreateFile. The pinned chain's window (winsec.go:179-187, i.e. pinnedParent plus the
// descent inside openPinnedLeaf) strictly CONTAINS the guard's, so the chain wins almost
// every gap hit: one 400-attempt run on this box counted "path escapes from parent"=75
// and preOpenLeafGuard=0. Measured counts of the same question: 0 and 1 on the same
// runner head 23b663c2 (runs 39884431726 / 40034438521), 0..15 on this box. A gate that
// is 0 or 1 is a gate on the machine's timing.
//
// What replaces it, and what each replacement can and cannot catch:
//   - TestAC3WindowsTheRouteCallsTheLeafGuardBeforeThePin322r1: the WIRING, read off the
//     call graph of winsec.go - the guard is called, and it is called between
//     resolveString and the pin. Red if the call is deleted or moved after the pin;
//   - TestAC3WindowsLeafGuardRefusesLinkNamesAndAcceptsRealOnes322r1: the LEG, driven
//     directly. Red if the guard stops refusing link names;
//   - TestAC3WindowsPrivateFileRouteHasTheLeafGuardSpeaking322r1: the ROUTE, with the
//     leaf-link state entered synchronously instead of raced. It gates "the route never
//     writes through a leaf link, and the leg that answers is the steady-state one",
//     plus a positive control that the same call lands once the link is gone. The gap
//     count itself stays a REPORTED reading (348's leaf world prints it), because it is
//     a reading about the clock.
package winsec

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// whyGuardRefusal322r1 is the verbatim reason winsec_windows.go's preOpenLeafGuard
// carries. Pinned rather than paraphrased: "it returned an error" is also true when a
// DIFFERENT leg answered, which is the false green 120-r3 measured.
const (
	whyGuardRefusal322r1 = "the name carries a reparse point, so the open that follows would go through that " +
		"link and truncate whatever it targets"
	winRacePayload322r1 = "raced private payload 322r1 windows"
	winLeafShots322r1   = 3 // quiet replicates: every shot re-enters the same state and must read the same
)

// The two steady-state rows a leaf link can be answered by on the route, verbatim from
// the legs that format them. They are kept apart on purpose: which one answers is the
// reading 票 369 is about, so a single "it was refused" assertion would erase it.
const (
	winLeafRowResolver322r1    = "path traverses a reparse point"
	winLeafRowFloorWalk322r1   = "traverses a reparse point at"
	winLeafRowGuardLink322r1   = whyGuardRefusal322r1
	winLeafRowGuardProbe322r1  = "the reparse-point probe could not read the name"
	winLeafRowRecheckLink322r1 = "the name carries a link, so this handle"
)

// leakRow322r1 answers WHICH leg refused, or "" when the text names none of them. It is
// deliberately a whitelist: an unrecognised refusal is not silently accepted as "some
// leg answered", because that is how a rename of a reason string would turn this case
// green without anything having changed.
func leakRow322r1(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	for _, why := range []string{
		winLeafRowGuardLink322r1,
		winLeafRowGuardProbe322r1,
		winLeafRowRecheckLink322r1,
		winLeafRowFloorWalk322r1,
		winLeafRowResolver322r1,
	} {
		if strings.Contains(msg, why) {
			return why
		}
	}
	return ""
}

// mkFileSymlink is mkDirSymlink's twin for a FILE target, with the same honesty rule:
// the privilege is conditional, so an unavailable construction is stated as a SKIP that
// acceptance has to count, never as an ok.
func mkFileSymlink(t *testing.T, link, target string) {
	t.Helper()
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("cannot construct a file symlink here: os.Symlink(%s, %s) = %v (needs developer mode or "+
			"SeCreateSymbolicLinkPrivilege)", target, link, err)
	}
	if !reparseAt(link) {
		t.Fatalf("%s is not a reparse point, so this test measured nothing", filepath.Base(link))
	}
}

func assertGuardRefusal322r1(t *testing.T, err error, tag string) {
	t.Helper()
	if err == nil {
		t.Errorf("AC#3 RED %s: preOpenLeafGuard accepted a name carrying a reparse point, so the open that "+
			"follows it can still truncate whatever that link targets", tag)
		return
	}
	if !strings.Contains(err.Error(), whyGuardRefusal322r1) {
		t.Errorf("AC#3 RED %s: the refusal did not carry the guard's own reason %q, so some other leg answered;"+
			" verbatim: %v", tag, whyGuardRefusal322r1, err)
	}
	if !errors.Is(err, ErrUnresolvedPath) || !errors.Is(err, ErrNotSealable) {
		t.Errorf("AC#3 RED %s: the refusal lost a sentinel (ErrUnresolvedPath/ErrNotSealable): %v", tag, err)
	}
	t.Logf("AC#3' %s: leg that spoke = preOpenLeafGuard, verbatim: %v", tag, err)
}

// TestAC3WindowsLeafGuardRefusesLinkNamesAndAcceptsRealOnes322r1 drives the guard
// directly over the names the open would have to resolve: a file symlink at a foreign
// file, a directory symlink, a junction (the shape that reports no ModeSymlink to
// os.Lstat on this platform), a real file (must stay writable - over-refusal is the other
// way to break the four production callers) and an absent name (the create case).
func TestAC3WindowsLeafGuardRefusesLinkNamesAndAcceptsRealOnes322r1(t *testing.T) {
	dir := filepath.Join(SealableTempDirForTest124(t), "tree")
	outside := filepath.Join(SealableTempDirForTest124(t), "outside-tree")
	// Both roots are built through filepath.Join rather than used raw: t.TempDir()
	// inherits the separator style of GOTMPDIR, and a forward-slash spelling handed to
	// cmd's mklink was measured here failing as `Invalid switch - "Users"` - a fixture
	// artefact that must never be read as a product refusal.
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outside, 0o700); err != nil {
		t.Fatal(err)
	}
	witness := filepath.Join(outside, "someone-elses-file-322r1.bin")
	const body = "the bytes a truncating open would have taken away"
	if err := os.WriteFile(witness, []byte(body), 0o666); err != nil {
		t.Fatal(err)
	}
	before, err := os.Lstat(witness)
	if err != nil {
		t.Fatal(err)
	}

	real := filepath.Join(dir, "plain-322r1.bin")
	if err := os.WriteFile(real, []byte("the tree's own object"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := preOpenLeafGuard(real); err != nil {
		t.Errorf("AC#3 RED: the guard refused a real file (%v), which would break every PrivateFile caller", err)
	}
	absent := filepath.Join(dir, "not-created-yet-322r1.bin")
	if err := preOpenLeafGuard(absent); err != nil {
		t.Errorf("AC#3 RED: the guard refused a name that is not there yet (%v), which would break every create", err)
	}

	mkFileSymlink(t, filepath.Join(dir, "file-sym-322r1.lnk"), witness)
	assertGuardRefusal322r1(t, preOpenLeafGuard(filepath.Join(dir, "file-sym-322r1.lnk")), "file-symlink")
	mkDirSymlink(t, filepath.Join(dir, "dir-sym-322r1.lnk"), outside)
	assertGuardRefusal322r1(t, preOpenLeafGuard(filepath.Join(dir, "dir-sym-322r1.lnk")), "directory-symlink")
	mkJunction(t, filepath.Join(dir, "junction-322r1.lnk"), outside)
	assertGuardRefusal322r1(t, preOpenLeafGuard(filepath.Join(dir, "junction-322r1.lnk")), "junction")

	after, err := os.Lstat(witness)
	if err != nil {
		t.Fatalf("the witness became unreachable: %v", err)
	}
	if got, rErr := os.ReadFile(witness); rErr != nil || string(got) != body {
		t.Errorf("AC#3 RED: the witness holds %q (err %v), want its own bytes", got, rErr)
	}
	if after.Size() != before.Size() || !after.ModTime().Equal(before.ModTime()) {
		t.Errorf("AC#3 RED: the witness moved from size %d mtime %s to size %d mtime %s", before.Size(),
			before.ModTime().Format(time.RFC3339Nano), after.Size(), after.ModTime().Format(time.RFC3339Nano))
	}
	for _, l := range []string{"file-sym-322r1.lnk", "junction-322r1.lnk"} {
		if _, err := os.Lstat(filepath.Join(dir, l)); err != nil {
			t.Errorf("AC#3 RED: the guard's refusal removed the link %s (%v)", l, err)
		}
	}
}

// TestAC3WindowsPrivateFileRouteHasTheLeafGuardSpeaking322r1 is the production wiring,
// entered deterministically instead of raced (see the header for why the count could not
// stay a gate).
//
// The world: the caller's name carries a RELATIVE symlink whose target is a bare sibling
// name, so no containment rule can answer it and the steady-state refusals are the only
// legs in front of the guard. The link is moved onto the name by rename and PROVED with
// reparseAt, so the state this case judges is entered by construction, not by timing.
//
// Judged: (a) the route refuses the caller's name while it carries a link - a nil here
// would mean privateFile wrote through a link, which is exactly the damage ticket 322
// exists for; (b) the refusal is attributable to a NAMED leg, and specifically to the
// steady state, which is the mechanism the header states; (c) neither the caller's name
// nor the foreign witness is removed, and the payload never lands in the witness; (d) the
// positive control: with the link replaced by the caller's own object, the same call
// lands and writes the payload - so (a) is a refusal about the link, not about a world
// where everything refuses.
//
// Reported, not gated: which steady-state refusal answered, and the truncation residual
// (a zero here is one box's reading; 322-a1's 571/759/780 is the same window with no
// guard at all, and the raced world in ancestorunlink_348_windows_test.go still counts
// it).
func TestAC3WindowsPrivateFileRouteHasTheLeafGuardSpeaking322r1(t *testing.T) {
	dir := SealableTempDirForTest124(t)
	name := filepath.Join(dir, "caller-name-322r1.bin")
	witness := filepath.Join(dir, "foreign-witness-322r1.bin")
	stageLink := filepath.Join(dir, "staged-leaf-322r1.lnk")
	const wbody = "foreign bytes are not yours to truncate, and this body is longer on purpose"
	if err := os.WriteFile(witness, []byte(wbody), 0o666); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(name, []byte("the object the caller owns"), 0o600); err != nil {
		t.Fatalf("give the caller a name to stand on: %v", err)
	}
	var answered []string
	truncatedShots := 0
	for shot := 0; shot < winLeafShots322r1; shot++ {
		// The link is staged inside the loop, not once above it: the rename below CONSUMES the
		// staged name, and a fixture that only works on its first shot is the same "measured once"
		// shape this whole change exists to remove.
		mkFileSymlink(t, stageLink, filepath.Base(witness))
		if err := os.Rename(stageLink, name); err != nil {
			t.Fatalf("shot %d: move the link onto the caller's name: %v", shot, err)
		}
		if !reparseAt(name) {
			t.Fatalf("shot %d: the caller's name is not a reparse point after the flip, so the state this "+
				"case judges was never entered and nothing below is a reading about a leaf link", shot)
		}
		err := PrivateFile(name, []byte(winRacePayload322r1), 0o600)
		leg := leakRow322r1(err)
		if err == nil {
			t.Fatalf("AC#3 RED shot %d: PrivateFile returned nil while the caller's name carried a link at a "+
				"foreign object - the route wrote THROUGH the leaf, which is ticket 322's damage in its "+
				"strongest form (name %s)", shot, name)
		}
		if leg == "" {
			t.Fatalf("AC#3 RED shot %d: the refusal names no leg this case declares, so it cannot be "+
				"attributed: %v", shot, err)
		}
		answered = append(answered, leg)
		if got, rErr := os.ReadFile(witness); rErr != nil {
			t.Fatalf("AC#3 shot %d: the witness became unreadable (%v), which is a removal rather than a "+
				"truncation and is ticket 348's numerator", shot, rErr)
		} else if string(got) != wbody {
			truncatedShots++
			if err := os.WriteFile(witness, []byte(wbody), 0o666); err != nil {
				t.Fatalf("restock the witness: %v", err)
			}
		}
		if strings.Contains(mustRead322r1(t, name), winRacePayload322r1) {
			t.Fatalf("AC#3 RED shot %d: the raced payload landed in the caller's name while it carried a link, "+
				"so the bytes went through it", shot)
		}
		if err := os.Remove(name); err != nil {
			t.Fatalf("shot %d: unlink the leaf link: %v", shot, err)
		}
		if err := os.WriteFile(name, []byte("the object the caller owns"), 0o600); err != nil {
			t.Fatalf("shot %d: restore the caller's own object: %v", shot, err)
		}
	}
	for i, leg := range answered {
		t.Logf("AC#3' shot %d: leg that answered while the leaf carried a link = %q", i, leg)
		// The gate is the FAMILY, not one member: either steady-state leg answering is the same
		// claim (the route refuses the link before any act), and which member answers is a fact
		// about what is linked into this binary (internal/risk's pipeline here, the package's own
		// floor when nothing is), so pinning one name would make the case red for a reason that is
		// not a defect. What may NOT answer is a leaf leg: that would mean the swap reached the
		// guard or the open, i.e. the steady state had stopped covering the leaf.
		if leg != winLeafRowResolver322r1 && leg != winLeafRowFloorWalk322r1 {
			t.Errorf("AC#3 RED shot %d: the leaf-link state was answered by %q, want one of the two "+
				"steady-state rows (%q / %q). Both of them walk the LEAF (resolve.go:605-611 then "+
				"placement_windows.go:59-68, whose last prefix IS the full path), which is the mechanism "+
				"this ticket's reading rests on; a leaf leg answering in the steady state means that order "+
				"moved and the gap counts have to be re-derived", i, leg,
				winLeafRowResolver322r1, winLeafRowFloorWalk322r1)
		}
	}
	if truncatedShots > 0 {
		t.Logf("AC#3' RESIDUAL OBSERVED on %d of %d synchronously entered shots: the foreign witness lost "+
			"its body. Reported, not claimed closed, and not gated: the guard's answer is about a NAME and "+
			"nothing in this route can lock it.", truncatedShots, winLeafShots322r1)
	}

	// Positive control: the same call, same directory, no link on the name.
	if err := PrivateFile(name, []byte(winRacePayload322r1), 0o600); err != nil {
		t.Errorf("AC#3 RED positive control: PrivateFile refused a plain, caller-owned, link-free name (%v), "+
			"so the refusals above say nothing about the link", err)
	} else if got := mustRead322r1(t, name); got != winRacePayload322r1 {
		t.Errorf("AC#3 RED positive control: the caller's name holds %q, want the payload", got)
	}
}

func mustRead322r1(t *testing.T, path string) string {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(body)
}
