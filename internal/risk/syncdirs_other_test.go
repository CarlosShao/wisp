//go:build !windows

package risk

import (
	"path/filepath"
	"testing"
)

// Ticket 82, AC#2 — the POSIX tier of the sync-root grading family.
//
// STATE AFTER TICKET 347. This file was written against a world where
// resolveHandle (pathresolver_other.go) was the DEFERRED(macOS/Linux) stub: no
// root could reach C26-canonical, so no grade name could be "confirmed",
// SyncDetectionComplete() was structurally false, and the under-profile suspect
// net was permanently armed. 347 landed realpath + lstat-reject there, and the
// three consequences were handled rather than re-tagged:
//
//   - the tripwire that asserted the absence (TestSyncNoGradeIsConfirmedOnPosix)
//     is GONE, and its inverted claim is now asserted on BOTH platforms by the
//     untagged TestSyncConfirmedGradesDisarmFallback in syncdirs_test.go. That
//     file's own header prescribed exactly this ("delete
//     TestSyncNoGradeIsConfirmedOnPosix") at the moment the resolver landed;
//   - the "weak grades may never disarm anything" half did not move: it lives in
//     TestSyncFallbackNotDisarmableByWeakRoot (untagged), which after 347 is a
//     real tooth - the injected root is canonical there too, so the grade is the
//     only thing keeping the net armed;
//   - TestSyncMembershipDecidesOffProfilePosix keeps its subject (root
//     MEMBERSHIP decides, not the blanket net) and states the flipped premise
//     out loud instead of silently inheriting it.
//
// What is still owed here, and by whom: syncdirs_other.go's registryProbe is
// still `return nil` under a DEFERRED(P12-macos) marker, so a real POSIX machine
// has no confirmed-grade location evidence (env/client-config only) and the
// macOS port's CloudStorage/~/Library probe is ticket 55's. The branch parked in
// provenance_syncdirs_other_test.go behind SyncDetectionComplete() flipped live
// by itself with this ticket and needs no edit here, which is why it was written
// as a branch and not as a Skip.

// TestSyncMembershipDecidesOffProfilePosix is the half that ticket 75's N-9 fix
// bought and that must not be lost again: an incomplete detector must not mean
// "every path is suspect", and membership is decided before the fallback. Both
// verdicts below come from root membership, which match() evaluates before the
// fallback, so they are assertable without the net.
//
// The third block used to assert that an under-profile write fail-closes to
// sync-suspect, which was true only while the net could never be disarmed. With
// the resolver in place the injected registry-grade root IS canonical, detection
// IS complete, and the same input must NOT be suspect any more. The polarity
// flipped because the measured truth flipped; the armed direction is not lost -
// TestSyncFallbackNotDisarmableByWeakRoot sweeps it on both platforms.
func TestSyncMembershipDecidesOffProfilePosix(t *testing.T) {
	e := membershipEngine(t, "registry")
	if st := e.p.IsSyncPath(e.syncTarget); !st.Sync || st.Root.Source != "registry" {
		t.Errorf("a write under the injected root must be sync BY MEMBERSHIP, got %+v", st)
	}
	if st := e.p.IsSyncPath(e.plainTarget); st.Sync {
		t.Errorf("ESCALATING FALSE POSITIVE: a plain out-of-profile write is sync-suspect on POSIX: %+v", st)
	}
	if !e.p.SyncDetectionComplete() {
		t.Fatal("premise of this file since ticket 347: the injected registry-grade root is a real directory, so detection must be complete here (it was structurally false before the POSIX resolver landed)")
	}
	// The profile itself: with a confirmed root the blanket net is lifted, so an
	// under-profile write is decided by membership like any other - and it is not
	// under the injected root.
	if st := e.p.IsSyncPath(filepath.Join(e.home, "Documents", "x.txt")); st.Sync {
		t.Errorf("with detection complete an under-profile write outside every root must not be suspect, got %+v", st)
	}
}

// TestSyncRegistryProbeIsAStubHere is the machine-checkable form of the reason
// this family needed tiering at all: syncdirs_other.go's registryProbe is a
// `return nil` carrying a DEFERRED(P12-macos) marker, so POSIX has no
// confirmed-grade location evidence to work with. It asserts the stub, so the
// day ticket 55 implements the probe this test goes red and points straight at
// the comment that promised it.
//
// Ticket 347 does NOT touch this: the resolver is not the probe. That is why the
// two halves are separate tests and why only the resolver half moved.
func TestSyncRegistryProbeIsAStubHere(t *testing.T) {
	home, _, _ := sandbox(t)
	if got := registryProbe(probeEnv{Home: home}); got != nil {
		t.Fatalf("registryProbe is the documented POSIX stub and must return nil until ticket 55 replaces it, got %+v", got)
	}
}
