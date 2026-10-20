package models

// Ticket 109 rework leg 109-r1, cell A (AC#2 landing + wording), shape (1) from
// docs/evidence/s1/109-adversarial-acceptance.md:43-44: "把守卫挪到 Ensure 的返回路径上".
//
// WHY SHAPE (1) AND NOT SHAPE (2). Shape (2) is "预置守卫，待接线启用" plus a red
// nail saying nobody calls it. That was the honest answer on 2026-09-21, when
// `go list -deps ./cmd/wisp` did not contain internal/models and
// DownloadingBridge.Run had zero production callers (R-109-4, the acceptance
// table's 攻2). Ticket 121 then landed AC#2's wiring and AC#3's gate, and both
// are green on the tree this leg measured:
//
//   - internal/models/assembly_reachability_121_test.go:126
//     TestAC3CapabilityPackagesReachTheCompositionRoot - PASS on GOOS=windows,
//     GOOS=linux, GOOS=darwin in this leg's baseline run (see the ticket face).
//   - the assembly chain that instrument now judges is a real call graph, not a
//     promise: cmd/wisp/models.go:305 calls models.WireDownloading, whose
//     DownloadingBridge.Run (internal/models/bridge.go:39) calls Ensure at
//     internal/models/bridge.go:44 and the hand-off guard at internal/models/bridge.go:58;
//     cmd/wisp/models.go:249 calls Manager.VerifyInstalled directly for
//     `wisp models verify`.
//
// With a production caller present, shape (2) would be dressing up a hole that no
// longer exists, and the acceptance table's own remaining complaint (R-109-4) is
// about the OTHER half: the guard sat on the bridge, so any caller that uses
// Manager.Ensure directly - the panel, an engine loader, a CLI sub-command - still
// received an unguarded hand-off. Shape (1) is the only one of the two that closes
// that, and it is what this file nails: Ensure's own return paths are the guard.
//
// THE SHAPE THAT IS RED BEFORE THE FIX. Ensure has three success paths, and only
// two of them ever looked at bytes after the last mutation point:
//
//	downloader.go:163 local_override  -> VerifyDir, then progress, then return
//	downloader.go:173 cache hit       -> VerifyDir, then progress, then return
//	downloader.go:203 fresh install   -> staging hash + move, then progress, then return
//
// The third one never re-reads the install directory after the move lands: the
// files are hashed inside staging (downloader.go:300) and then moved
// (downloader.go:710), and what the caller receives is "those bytes, wherever they
// are now". A principal that can write the install directory - the grant AC#1
// measured at SID level - owns everything between the move and the return. The
// cases below swap a file at exactly that moment, which is the same construction
// ticket 109's own hook uses: inside the Progress(PhaseDone) callback, after the
// last check that ran and before Ensure returns.
//
// AC#2's reverse leg is here too: a clean install and a warm cache hit must still
// hand back the same directory with the same number of server hits, so moving the
// guard cannot silently turn into "Ensure now refuses everything".

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// swapAtPhaseDone returns a Progress hook that flips one byte of one installed
// file the first time Ensure announces completion, plus the swap count and a func
// putting the byte back. The hook runs inside Ensure, after its last verification
// and before its return, so the swap is the post-verification write AC#1 measured
// - not a hypothesis about one. The restore matters: legs below share one store,
// and a byte left swapped would make the next reading red for the wrong reason.
func swapAtPhaseDone(t *testing.T, dir, id, path string) (hook func(ProgressEvent), count func() int, restoreAll func()) {
	t.Helper()
	swaps := 0
	var restore func()
	hook = func(ev ProgressEvent) {
		if ev.Phase != PhaseDone || ev.ModelID != id || swaps > 0 {
			return
		}
		swaps++
		restore = swapOneInstalledFile(t, filepath.Join(dir, id), path)
	}
	count = func() int { return swaps }
	restoreAll = func() {
		if restore != nil {
			restore()
		}
	}
	return hook, count, restoreAll
}

// TestAC21FreshInstallIsGuardedOnEnsuresOwnReturnPath is cell A's red-before-fix
// case: the fresh-install path (downloader.go:203-208), reached with NO bridge, so
// the guard ticket 109 put on the hand-off point cannot be what saves it.
func TestAC21FreshInstallIsGuardedOnEnsuresOwnReturnPath(t *testing.T) {
	srv := newCountingServer(t, map[string][]byte{
		"/kws-fixture/tiny-model.tar.bz2": mustRead(t, filepath.Join("testdata", "tiny-model.tar.bz2")),
	})
	m := manifestWithURLs(t, srv.URL())
	entry, err := m.FindModel("kws-fixture")
	if err != nil {
		t.Fatal(err)
	}
	store := filepath.Join(t.TempDir(), "models")
	if err := os.MkdirAll(store, 0o755); err != nil {
		t.Fatal(err)
	}
	hook, count, restoreSwap := swapAtPhaseDone(t, store, "kws-fixture", entry.InstalledFiles()[0].Path)
	mgr, err := NewManager(Options{
		DataDir: store, Manifest: m, VerifySignature: true, Attempts: 2, BackoffBase: 0, Progress: hook,
	})
	if err != nil {
		t.Fatal(err)
	}

	dir, err := mgr.Ensure(context.Background(), "kws-fixture")
	restoreSwap()
	if err == nil {
		t.Fatalf("AC#2/R-109-4 (shape 1): Ensure returned a model directory as available although 1 byte of %s "+
			"was swapped after the last check that ran on its return path. dir=%s swaps=%d - "+
			"a caller that uses Manager.Ensure directly, without DownloadingBridge, gets the hand-off with no "+
			"guard on it at all, which is exactly what R-109-4 says the bridge-only placement left open.",
			entry.InstalledFiles()[0].Path, dir, count())
	}
	if count() == 0 {
		t.Fatalf("the swap never fired, so this case proved nothing about the window (err=%v)", err)
	}
	if !strings.Contains(err.Error(), "sha256 mismatch") {
		t.Errorf("the refusal must come from the re-verification and name the hash, got: %v", err)
	}
	t.Logf("Ensure's own return path refused the post-verification swap: %v", err)
}

// TestAC22CacheHitAndOverrideReturnPathsAreGuarded covers the other two success
// paths by the same route: no bridge, so only a guard Ensure itself runs can see
// the write. The cache hit verifies before its PhaseDone tick (downloader.go:173),
// so the swap lands inside the span that remains.
func TestAC22CacheHitAndOverrideReturnPathsAreGuarded(t *testing.T) {
	srv := newCountingServer(t, map[string][]byte{
		"/kws-fixture/tiny-model.tar.bz2": mustRead(t, filepath.Join("testdata", "tiny-model.tar.bz2")),
	})
	m := manifestWithURLs(t, srv.URL())
	entry, err := m.FindModel("kws-fixture")
	if err != nil {
		t.Fatal(err)
	}
	store := filepath.Join(t.TempDir(), "models")
	if err := os.MkdirAll(store, 0o755); err != nil {
		t.Fatal(err)
	}
	plain, err := NewManager(Options{DataDir: store, Manifest: m, VerifySignature: true, Attempts: 2, BackoffBase: 0})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := plain.Ensure(context.Background(), "kws-fixture"); err != nil {
		t.Fatalf("clean first install: %v", err)
	}

	t.Run("cache hit", func(t *testing.T) {
		hook, count, restoreSwap := swapAtPhaseDone(t, store, "kws-fixture", entry.InstalledFiles()[0].Path)
		mgr, err := NewManager(Options{DataDir: store, Manifest: m, VerifySignature: true, Attempts: 2, BackoffBase: 0, Progress: hook})
		if err != nil {
			t.Fatal(err)
		}
		_, err = mgr.Ensure(context.Background(), "kws-fixture")
		restoreSwap()
		if err == nil {
			t.Fatalf("AC#2/R-109-4: the cache-hit return path handed back a store whose %s was swapped after "+
				"its verification (swaps=%d)", entry.InstalledFiles()[0].Path, count())
		}
		if count() == 0 {
			t.Fatalf("the swap never fired, so this leg proved nothing (err=%v)", err)
		}
		t.Logf("cache-hit return path refused: %v", err)
	})

	t.Run("local_override", func(t *testing.T) {
		side := filepath.Join(t.TempDir(), "sideload")
		if err := os.MkdirAll(filepath.Join(side, "dict"), 0o755); err != nil {
			t.Fatal(err)
		}
		for _, want := range entry.InstalledFiles() {
			body, rerr := os.ReadFile(filepath.Join(store, "kws-fixture", filepath.FromSlash(want.Path)))
			if rerr != nil {
				t.Fatal(rerr)
			}
			if werr := os.WriteFile(filepath.Join(side, filepath.FromSlash(want.Path)), body, 0o644); werr != nil {
				t.Fatal(werr)
			}
		}
		swapped := false
		var restore func()
		hook := func(ev ProgressEvent) {
			if ev.Phase != PhaseDone || swapped {
				return
			}
			swapped = true
			restore = swapOneInstalledFile(t, side, entry.InstalledFiles()[0].Path)
		}
		mgr, err := NewManager(Options{
			DataDir: store, Manifest: m, VerifySignature: true, Attempts: 2, BackoffBase: 0,
			LocalOverride: map[string]string{"kws-fixture": side}, Progress: hook,
		})
		if err != nil {
			t.Fatal(err)
		}
		_, err = mgr.Ensure(context.Background(), "kws-fixture")
		if restore != nil {
			restore()
		}
		if err == nil {
			t.Fatalf("AC#2/R-109-4: the local_override return path announced a sideloaded model available although "+
				"1 byte was swapped after its verification (swapped=%v)", swapped)
		}
		if !swapped {
			t.Fatalf("the override swap never fired, so this leg proved nothing (err=%v)", err)
		}
		t.Logf("local_override return path refused: %v", err)
	})
}

// TestAC23ReverseLegCleanHandOffsAreNotBroken is AC#2's reverse half for the new
// landing: every success path must still hand back the directory, and the second
// instance must still reuse the bytes instead of re-fetching them. Ticket 95's
// reason for not sealing this class is the thing this must not disturb.
func TestAC23ReverseLegCleanHandOffsAreNotBroken(t *testing.T) {
	srv := newCountingServer(t, map[string][]byte{
		"/kws-fixture/tiny-model.tar.bz2": mustRead(t, filepath.Join("testdata", "tiny-model.tar.bz2")),
	})
	m := manifestWithURLs(t, srv.URL())
	store := filepath.Join(t.TempDir(), "models")
	if err := os.MkdirAll(store, 0o755); err != nil {
		t.Fatal(err)
	}
	newMgr := func() *Manager {
		mgr, err := NewManager(Options{DataDir: store, Manifest: m, VerifySignature: true, Attempts: 2, BackoffBase: 0})
		if err != nil {
			t.Fatal(err)
		}
		return mgr
	}
	first, second := newMgr(), newMgr()
	ctx := context.Background()

	hits0 := servedHits(srv)
	dir1, err := first.Ensure(ctx, "kws-fixture")
	if err != nil {
		t.Fatalf("fresh install through the guarded return path: %v", err)
	}
	hits1 := servedHits(srv)
	dir2, err := second.Ensure(ctx, "kws-fixture")
	if err != nil {
		t.Fatalf("second instance reusing the cache: %v", err)
	}
	hits2 := servedHits(srv)
	t.Logf("server hits: install=%d, second-instance reuse=%d", hits1-hits0, hits2-hits1)
	if dir1 != dir2 {
		t.Errorf("the guarded return path handed back a different directory: %s vs %s", dir1, dir2)
	}
	if hits2 != hits1 {
		t.Errorf("AC#2 reverse leg: the guard turned reuse into a re-download (%d extra hits)", hits2-hits1)
	}
	if err := second.VerifyInstalled("kws-fixture"); err != nil {
		t.Errorf("VerifyInstalled on the same clean tree: %v", err)
	}
	if _, err := WireDownloading(second, newWalkMachine()).Run(ctx, "kws-fixture"); err != nil {
		t.Errorf("the hand-off point on a clean shared cache: %v", err)
	}
}
