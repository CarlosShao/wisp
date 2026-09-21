//go:build !windows

package risk

import (
	"os"
	"path/filepath"
	"testing"
)

// Ticket 82, AC#2 — the POSIX tier of the sync-root grading family.
//
// AC#2 forbids this file from being empty or from being a wall of t.Skip, so
// every case below asserts what POSIX ACTUALLY does today, and none of it is
// conditional:
//
//   - no root ever reaches C26-canonical here, because resolveHandle
//     (pathresolver_other.go) is the DEFERRED(macOS/Linux) stub; therefore no
//     probe grade, however strong its name, is "confirmed";
//   - so SyncDetectionComplete() is false, the under-profile suspect net stays
//     ARMED, and every write inside the profile fail-closes to sync-suspect —
//     the safe side of P12, not a bug to route around;
//   - and yet root MEMBERSHIP is still decided properly (syncdirs_test.go's
//     membershipEngine, which all six portable cases now run through), i.e.
//     POSIX is not "everything is dirty", it is "nothing may be EXCUSED".
//
// HOW THIS FILE TURNS LIVE WHEN TICKET 55 LANDS (A51⑧, spelled out so nobody
// has to re-derive it): the moment a darwin/linux probe produces
// confirmed-grade evidence that Resolve can canonicalize — a CloudStorage /
// ~/Library/CloudServices directory probe on macOS, realpath+lstat in
// resolveHandle — `TestSyncNoGradeIsConfirmedOnPosix` must FAIL. That is its
// purpose: it is a tripwire, not a comfort. Then, in this order:
//  1. move TestSyncConfirmedGradesDisarmFallbackWindows and
//     TestWriteGatePlainTargetInsideProfileWindows from syncdirs_windows_test.go
//     into syncdirs_test.go verbatim (do NOT copy: two tiers with the same
//     expectation drift), and delete syncdirs_windows_test.go once it is empty;
//  2. delete TestSyncNoGradeIsConfirmedOnPosix and
//     TestSyncRegistryProbeIsAStubHere from this file;
//  3. keep TestSyncMembershipDecidesOffProfilePosix — the out-of-profile
//     shape stays a legitimate second configuration — and un-gate the dormant
//     negative half of TestExfilSyncWritePosixSpellingInvariant, which
//     provenance_syncdirs_other_test.go already parks behind
//     SyncDetectionComplete(); that branch flips by itself and needs no edit
//     here, which is why it was written as a branch and not as a Skip.
//
// Ticket 55 owns steps 1-3; nothing in this file is a substitute for that work.

// TestSyncNoGradeIsConfirmedOnPosix pins the fail-closed reality behind
// TestSyncFallbackNotDisarmableByWeakRoot's mirror half: on POSIX NO grade name
// — not even the confirmed-grade ones — disarms the under-profile net, because
// the second half of the definition (a C26-canonical root) can never be
// satisfied by this platform's resolver yet.
func TestSyncNoGradeIsConfirmedOnPosix(t *testing.T) {
	home, _, _ := sandbox(t)
	root := filepath.Join(home, "OneDrive")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	// The mechanism, asserted rather than assumed: an EXISTING directory is
	// still not handle-resolvable here, which is the single reason finalize()
	// never sets complete. If this starts passing, everything below it is the
	// first thing ticket 55 has to revisit.
	if res, err := Resolve(root, nil); err != nil || res.Resolved {
		t.Fatalf("premise of this whole tier: Resolve(%q) must report Resolved=false on POSIX, got %+v err=%v", root, res, err)
	}
	for _, src := range []string{"default", "fixture", "options", "", "registry", "config", "env"} {
		p := NewProvenance(ProvOptions{
			NoProbe: true, HomeDir: home,
			SyncRoots: []SyncRoot{{Provider: "OneDrive", Path: root, Source: src}},
		})
		if p.SyncDetectionComplete() {
			t.Errorf("source %q: POSIX detection cannot be complete while resolveHandle is the DEFERRED stub (ticket 55)", src)
		}
		st := p.IsSyncPath(filepath.Join(home, "Documents", "exfil.md"))
		if !st.Sync || st.Root.Source != "suspect-fallback" {
			t.Errorf("source %q: the under-profile net must stay armed and be NAMED as the reason, got %+v", src, st)
		}
	}
}

// TestSyncMembershipDecidesOffProfilePosix is the half that ticket 75's N-9 fix
// bought and that must not be lost again: incomplete detection means writes
// inside the profile are suspect, NOT that every path is suspect. Both verdicts
// below come from root membership, which match() evaluates before the fallback,
// so they are assertable without anything ticket 55 owes us.
func TestSyncMembershipDecidesOffProfilePosix(t *testing.T) {
	e := membershipEngine(t, "registry")
	if st := e.p.IsSyncPath(e.syncTarget); !st.Sync || st.Root.Source != "registry" {
		t.Errorf("a write under the injected root must be sync BY MEMBERSHIP, got %+v", st)
	}
	if st := e.p.IsSyncPath(e.plainTarget); st.Sync {
		t.Errorf("ESCALATING FALSE POSITIVE: a plain out-of-profile write is sync-suspect on POSIX: %+v", st)
	}
	if e.p.SyncDetectionComplete() {
		t.Error("premise of this file: see TestSyncNoGradeIsConfirmedOnPosix")
	}
	// And the profile itself is exactly where the net does apply — the
	// asymmetry the two assertions above depend on.
	if st := e.p.IsSyncPath(filepath.Join(e.home, "Documents", "x.txt")); !st.Sync || st.Root.Source != "suspect-fallback" {
		t.Errorf("an under-profile write must fail closed as sync-suspect, got %+v", st)
	}
}

// TestSyncRegistryProbeIsAStubHere is the machine-checkable form of the reason
// this family needed tiering at all: syncdirs_other.go's registryProbe is a
// `return nil` carrying a DEFERRED(P12-macos) marker, so POSIX has no
// confirmed-grade location evidence to work with. It asserts the stub, so the
// day ticket 55 implements the probe this test goes red and points straight at
// the comment that promised it (step 2 of this file's header).
func TestSyncRegistryProbeIsAStubHere(t *testing.T) {
	home, _, _ := sandbox(t)
	if got := registryProbe(probeEnv{Home: home}); got != nil {
		t.Fatalf("registryProbe is the documented POSIX stub and must return nil until ticket 55 replaces it, got %+v", got)
	}
}
