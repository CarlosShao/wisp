//go:build windows

package risk

import (
	"os"
	"path/filepath"
	"testing"
)

// Ticket 82, AC#1/AC#2 — the WINDOWS tier of the sync-root grading family.
//
// Everything in this file asserts a consequence of ONE capability: a sync root
// can become C26-CANONICAL, because resolveHandle asks the OS for the final path
// of an existing directory. syncSet.add sets canonical=true only on that answer,
// and syncSet.finalize counts a root as "confirmed" only when it is BOTH
// canonical and confirmed-GRADE (env/registry/config).
//
// TICKET 347 CHANGED WHAT IS LEFT HERE. Until it landed, resolveHandle
// (pathresolver_other.go) was the DEFERRED(macOS/Linux) stub, so canonical - and
// with it SyncDetectionComplete() - was structurally unreachable on POSIX and
// every case below was Windows-only for that reason. POSIX now resolves paths
// (realpath + lstat-reject, PLAN.md:2385's macOS half), so the grade consequence
// itself is no longer platform-bound: TestSyncConfirmedGradesDisarmFallback was
// MOVED (not copied) into syncdirs_test.go, per the instruction this file used to
// carry, and the POSIX tripwire that asserted the absence went with it.
//
// What stays: the case whose subject really is this platform's API
// (TestWriteGatePlainTargetInsideProfileWindows is Windows-only by name and by
// the write-gate fixture it uses), and the live registry probe below. Neither may
// be moved by a ticket that does not own the write gate, so ticket 55's step 1
// ("move each case below") is now half done and this note says which half.
//
// syncdirs_other.go's registryProbe still returns nil, which is why the live
// probe below remains the only evidence this platform can produce.

// TestWriteGatePlainTargetInsideProfileWindows is the assertion the four
// TestWriteGate* cases used to carry as their own precondition, and the one
// thing about the write gate that is genuinely platform-bound: a plain target
// that sits INSIDE the user profile is only judged non-sync when a confirmed
// root has disarmed the blanket net (membershipEngine had to move its targets
// out of the profile to become portable — see AC#1's table in the ticket).
//
// Both directions are here on purpose. The negative half is what an incomplete
// detector must never grant; the positive half (same engine, same profile, a
// target inside the root) is what proves the case is not being carried by
// "nothing flags anywhere", the exact failure mode ticket 75's report named.
func TestWriteGatePlainTargetInsideProfileWindows(t *testing.T) {
	base := t.TempDir()
	home := filepath.Join(base, "profile")
	root := filepath.Join(home, "OneDrive")
	work := filepath.Join(home, "work")
	for _, d := range []string{root, work} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	p := NewProvenance(ProvOptions{NoProbe: true, HomeDir: home, SyncRoots: []SyncRoot{
		{Provider: "OneDrive", Path: root, Source: "registry"},
	}})
	if !p.SyncDetectionComplete() {
		t.Fatal("precondition: the injected registry-grade root must confirm detection, so the suspect net is NOT what decides these cases")
	}
	inside := filepath.Join(work, "brand-new.md") // under-profile, not under any root
	if st := p.IsSyncPath(inside); st.Sync {
		t.Fatalf("a confirmed root must have disarmed the under-profile net for %s: %+v", inside, st)
	}
	p.OpenScope("task-1")
	if !p.Mark("task-1", SrcWebFetch, "https://x", "quote: "+marker) {
		t.Fatal("mark rejected")
	}
	if _, ok := p.Inspect("task-1", "fs.write", map[string]any{
		"path": inside, "content": "local bytes " + marker,
	}); ok {
		t.Error("false positive: an in-profile write of a disarmed engine became an exfil channel")
	}
	// Positive half: the root in the same profile still decides.
	if _, ok := p.Inspect("task-1", "fs.write", map[string]any{
		"path": filepath.Join(root, "Notes", "out.md"), "content": "leak " + marker,
	}); !ok {
		t.Error("ESCAPIABLE: with the net disarmed, a write into the confirmed root stopped flagging")
	}
}

// TestSyncRegistryProbeLive is P12 live registry evidence (skipped where the
// hive has nothing to say): it asserts the probe's SHAPE, never machine state.
//
// Ticket 93 AC#2 moved it here verbatim from syncdirs_test.go. Reason: the
// subject is `registryProbe`, whose POSIX implementation (syncdirs_other.go)
// has no registry to read, so on that platform the loop body can never execute
// and `len(roots) == 0` is a property of the OS, not of the code. Under a bare
// `go test` that reads as `ok` on both platforms, which is why this case had
// never produced a verdict in CI. It is now evaluated only where the API
// exists; the file-level //go:build windows above is the constraint, and the
// package itself stays portable-tagged, so `go vet ./...` still type-checks it
// on both platforms (see the go-list numbers on the ticket).
func TestSyncRegistryProbeLive(t *testing.T) {
	roots := registryProbe(probeEnv{Home: userHomeDir()})
	for _, r := range roots {
		if r.Source != "registry" || r.Provider == "" || r.Path == "" {
			t.Fatalf("live registry probe returned a malformed root %+v", r)
		}
		t.Logf("P12 evidence: registry-confirmed %s root %s", r.Provider, r.Path)
	}
	if len(roots) == 0 {
		t.Skip("no registry-grade sync record on this machine (HKCU Accounts without UserFolder is the documented reality here)")
	}
}
