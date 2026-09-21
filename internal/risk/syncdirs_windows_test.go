//go:build windows

package risk

import (
	"os"
	"path/filepath"
	"testing"
)

// Ticket 82, AC#1/AC#2 — the WINDOWS tier of the sync-root grading family.
//
// Everything in this file asserts a consequence of ONE capability that exists
// only here: a sync root can become C26-CANONICAL, because resolveHandle
// (pathresolver_windows.go) asks the OS for the final path of an existing
// directory. syncSet.add sets canonical=true only on that answer, and
// syncSet.finalize counts a root as "confirmed" only when it is BOTH canonical
// and confirmed-GRADE (env/registry/config). So `SyncDetectionComplete()` is
// structurally always false on POSIX (pathresolver_other.go's resolveHandle is
// the DEFERRED(macOS/Linux) stub, syncdirs_other.go's registryProbe returns nil)
// — those were the 8 reds, and they were red on a precondition, not on a
// verdict.
//
// The portable subjects (root membership, the payload-key-blind write gate, the
// N-10 `..` hardening, the weak grades that may never disarm anything) stayed in
// syncdirs_test.go / provenance_test.go and now run on both platforms via
// membershipEngine. What is left here is only what needs a disarmed blanket net,
// i.e. a verdict that POSIX cannot produce yet. syncdirs_other_test.go asserts
// that absence instead of inheriting these expectations as a red.
//
// WHEN TICKET 55 LANDS (the macOS/Linux probe), this file is the donor: move
// each case below into syncdirs_test.go verbatim and delete this file with its
// //go:build windows line. Do NOT copy them — a duplicated expectation is how
// the two tiers drift apart.

// TestSyncConfirmedGradesDisarmFallbackWindows is the mirror half of
// TestSyncFallbackNotDisarmableByWeakRoot: which grades are strong enough to
// switch the under-profile suspect net off. Windows-only today for the reason
// stated in this file's header.
func TestSyncConfirmedGradesDisarmFallbackWindows(t *testing.T) {
	home, _, _ := sandbox(t)
	root := filepath.Join(home, "OneDrive")
	if err := os.MkdirAll(root, 0o755); err != nil { // confirmed roots are real dirs
		t.Fatal(err)
	}
	for _, src := range []string{"registry", "config", "env"} {
		p := NewProvenance(ProvOptions{
			NoProbe: true, HomeDir: home,
			SyncRoots: []SyncRoot{{Provider: "OneDrive", Path: root, Source: src}},
		})
		if !p.SyncDetectionComplete() {
			t.Errorf("source %q must count as confirmed", src)
		}
		if p.IsSyncPath(filepath.Join(home, "Documents", "exfil.md")).Sync {
			t.Errorf("source %q: confirmed root must lift the blanket suspect net", src)
		}
		// ...while the confirmed root itself keeps flagging.
		if !p.IsSyncPath(filepath.Join(root, "notes.md")).Sync {
			t.Errorf("source %q: writes under the confirmed root must be sync", src)
		}
	}
}

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
