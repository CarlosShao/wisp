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
//     truncating count below is therefore REPORTED, not gated - 322-a1 measured the same
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

// whyGuardRefusal322r1 is the verbatim reason winsec_windows.go's preOpenLeafGuard
// carries. Pinned rather than paraphrased: "it returned an error" is also true when a
// DIFFERENT leg answered, which is the false green 120-r3 measured.
const (
	whyGuardRefusal322r1 = "the name carries a reparse point, so the open that follows would go through that " +
		"link and truncate whatever it targets"
	winRaceAttempts322r1 = 400
	winRacePayload322r1  = "raced private payload 322r1 windows"
)

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
// measured on this machine rather than inferred from the POSIX row. The name is flipped
// by rename between a real file and a link at a foreign witness; the hard judgement is
// that the NEW leg answered at least once through PrivateFile, because that is the only
// way this file can say the guard is on the route and not merely on the shelf. The
// truncating count is printed as a reading: the guard cannot lock its own answer, so a
// non-zero residual is the shape the adjudication named, not a failure of this leg.
func TestAC3WindowsPrivateFileRouteHasTheLeafGuardSpeaking322r1(t *testing.T) {
	dir := SealableTempDirForTest124(t)
	name := filepath.Join(dir, "caller-name-322r1.bin")
	witness := filepath.Join(dir, "foreign-witness-322r1.bin")
	stageReal := filepath.Join(dir, "staged-real-322r1.bin")
	stageLink := filepath.Join(dir, "staged-leaf-322r1.lnk")
	const wbody = "foreign bytes are not yours to truncate"
	if err := os.WriteFile(witness, []byte(wbody), 0o666); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(name, []byte("the object the caller owns"), 0o600); err != nil {
		t.Fatalf("give the caller a name to stand on: %v", err)
	}

	var flips, guardSpoke, otherRefused, nilReturned int32
	var flipErrs int64
	var flipMu sync.Mutex
	var firstFlip, firstGuard string
	stop, done := make(chan struct{}), make(chan struct{})
	var stopOnce sync.Once
	stopRacer := func() { stopOnce.Do(func() { close(stop) }); <-done }
	go func() {
		defer close(done)
		defer func() {
			if p := recover(); p != nil {
				atomic.AddInt64(&flipErrs, 1)
			}
		}()
		for {
			select {
			case <-stop:
				return
			default:
			}
			if err := os.WriteFile(stageReal, []byte("the object the caller owns"), 0o600); err != nil {
				flipMu.Lock()
				if firstFlip == "" {
					firstFlip = err.Error()
				}
				flipMu.Unlock()
				continue
			}
			_ = os.Rename(stageReal, name)
			_ = os.Remove(stageLink)
			if err := os.Symlink(filepath.Base(witness), stageLink); err != nil {
				flipMu.Lock()
				if firstFlip == "" {
					firstFlip = err.Error()
				}
				flipMu.Unlock()
				continue
			}
			_ = os.Rename(stageLink, name)
			atomic.AddInt32(&flips, 1)
		}
	}()
	defer stopRacer()

	for i := 0; i < winRaceAttempts322r1; i++ {
		err := PrivateFile(name, []byte(winRacePayload322r1), 0o600)
		switch {
		case err == nil:
			atomic.AddInt32(&nilReturned, 1)
		case strings.Contains(err.Error(), whyGuardRefusal322r1):
			atomic.AddInt32(&guardSpoke, 1)
			flipMu.Lock()
			if firstGuard == "" {
				firstGuard = err.Error()
			}
			flipMu.Unlock()
		default:
			atomic.AddInt32(&otherRefused, 1)
		}
	}
	stopRacer()

	truncated := false
	if got, rErr := os.ReadFile(witness); rErr != nil || string(got) != wbody {
		truncated = true
	}
	info, _ := os.Lstat(witness)
	// Judged: the new leg answered on the production route, and the racer really flipped.
	if guardSpoke == 0 || flips == 0 {
		t.Errorf("AC#3 RED instrument: %d raced attempts produced guard-leg refusals=%d and racer flips=%d "+
			"(flipErrs=%d, firstFlipErr=%q), so nothing below proves the guard is on this route",
			winRaceAttempts322r1, guardSpoke, flips, flipErrs, firstFlip)
	}
	size, mode := int64(-1), "unreachable"
	if info != nil {
		size, mode = info.Size(), fmt.Sprint(info.Mode())
	}
	// Reported, because the probe cannot lock its own answer: a residual open that still
	// followed a swapped leaf is the named half-open window, not a new hole.
	t.Logf("AC#3 local paired reading (winsec_windows.go leaf guard): attempts=%d flips=%d guardRefusals=%d "+
		"otherRefusals=%d returnedNil=%d witnessTruncated=%v witnessSize=%d witnessMode=%v firstGuardVerbatim=%q",
		winRaceAttempts322r1, flips, guardSpoke, otherRefused, nilReturned, truncated, size, mode, firstGuard)
}
