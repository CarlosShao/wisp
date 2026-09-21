package models

// Ticket 109's other two halves, portable on purpose (no ACL text parsed):
//
//   - AC#2's reverse leg: closing the write window must not undo the reason
//     ticket 95 gave for not sealing this class. A normal install still hands
//     over, and a second Manager on the same store - what another running
//     instance is - still reuses the bytes instead of re-fetching them.
//   - AC95-R3's leg: "wide temp + rename into place" is a sequence every
//     installer here uses, and ticket 95's parse.go case only passed because the
//     descriptor happened to survive the move. What gets pinned is the part that
//     holds whichever way the descriptor lands.

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestAC2NormalInstallAndSecondInstanceReuseAreNotBroken(t *testing.T) {
	srv := newCountingServer(t, map[string][]byte{
		"/kws-fixture/tiny-model.tar.bz2": mustRead(t, filepath.Join("testdata", "tiny-model.tar.bz2")),
	})
	m := manifestWithURLs(t, srv.URL())
	store := filepath.Join(t.TempDir(), "models")
	if err := os.MkdirAll(store, 0o755); err != nil {
		t.Fatal(err)
	}
	newMgr := func() *Manager {
		mgr, err := NewManager(Options{
			DataDir: store, Manifest: m, VerifySignature: true, Attempts: 2, BackoffBase: 0,
		})
		if err != nil {
			t.Fatal(err)
		}
		return mgr
	}
	first, second := newMgr(), newMgr()
	ctx := context.Background()

	hits0 := servedHits(srv)
	if _, err := first.Ensure(ctx, "kws-fixture"); err != nil {
		t.Fatalf("instance 1 install: %v", err)
	}
	hits1 := servedHits(srv)
	if err := first.VerifyInstalled("kws-fixture"); err != nil {
		t.Fatalf("AC#2 reverse leg, normal install: the hand-off check refused a clean model: %v", err)
	}
	if _, err := second.Ensure(ctx, "kws-fixture"); err != nil {
		t.Fatalf("instance 2 reusing instance 1's cache: %v", err)
	}
	hits2 := servedHits(srv)
	if err := second.VerifyInstalled("kws-fixture"); err != nil {
		t.Fatalf("instance 2 hand-off check on the shared cache: %v", err)
	}
	// Reuse has to still be a reuse: the second instance must fetch nothing,
	// otherwise the guard has quietly invalidated the cache ticket 95 kept open.
	t.Logf("server hits: install=%d, second-instance reuse=%d", hits1-hits0, hits2-hits1)
	if hits2 != hits1 {
		t.Errorf("AC#2 reverse leg: the second instance re-downloaded (%d extra hits) - multi-instance reuse is broken", hits2-hits1)
	}
	if _, err := second.Ensure(ctx, "kws-fixture"); err != nil {
		t.Fatalf("instance 2 re-reusing: %v", err)
	}
	// And the hand-off point itself, on a clean shared cache.
	if _, err := WireDownloading(second, newWalkMachine()).Run(ctx, "kws-fixture"); err != nil {
		t.Fatalf("hand-off on a clean shared cache: %v", err)
	}
}

// TestAC3WideTempRenameIsGuardedByTheReverifyNotByTheDescriptor is AC95-R3.
func TestAC3WideTempRenameIsGuardedByTheReverifyNotByTheDescriptor(t *testing.T) {
	srv := newCountingServer(t, map[string][]byte{
		"/kws-fixture/tiny-model.tar.bz2": mustRead(t, filepath.Join("testdata", "tiny-model.tar.bz2")),
	})
	m := manifestWithURLs(t, srv.URL())
	mgr, _ := newTestManager(t, m, nil, nil)
	entry, err := m.FindModel("kws-fixture")
	if err != nil {
		t.Fatal(err)
	}
	dir, err := mgr.Ensure(context.Background(), "kws-fixture")
	if err != nil {
		t.Fatal(err)
	}
	// The wide temp is gone: no second writable position left uncovered.
	if _, err := os.Stat(mgr.stagingDir("kws-fixture")); !os.IsNotExist(err) {
		t.Errorf("the staging temp still stands after install, so the wide-temp half of the sequence is unguarded: %v", err)
	}
	for _, want := range entry.InstalledFiles() {
		st, err := os.Lstat(filepath.Join(dir, filepath.FromSlash(want.Path)))
		if err != nil || !st.Mode().IsRegular() {
			t.Fatalf("installed %s missing after the rename sequence: %v", want.Path, err)
		}
	}
	target := entry.InstalledFiles()[0].Path
	restore := swapOneInstalledFile(t, dir, target)
	err = mgr.VerifyInstalled("kws-fixture")
	restore()
	if err == nil {
		t.Fatalf("AC#3: a file that reached its place through wide-temp+rename was then swapped and the hand-off point said nothing - "+
			"the guard must be the re-verification, not whichever descriptor the rename happened to keep (target %s)", target)
	}
	t.Logf("rename-into-place covered by the hand-off re-verify: %v", err)
	if err := mgr.VerifyInstalled("kws-fixture"); err != nil {
		t.Errorf("after restoring the byte, the clean tree no longer verifies: %v", err)
	}
	// The next Ensure has to re-verify too: the harm ceiling is "one session",
	// not "permanent" (AC#1's third question).
	if _, err := mgr.Ensure(context.Background(), "kws-fixture"); err != nil {
		t.Fatalf("Ensure after the swap was caught: %v", err)
	}
}

func servedHits(cs *countingServer) int {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return cs.hits
}
