//go:build windows

// Production-path evidence for ticket 89: the write that matters is the one the
// application performs, not the one a test performs on its behalf. Everything
// under a data root created by memory.Open / secret.NewStore must come out
// private even though SQLite and DPAPI are the ones writing bytes.
package winsec_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/secret"
)

func TestAC3ProductionDataRootIsPrivateEndToEnd(t *testing.T) {
	// A parent that grants Everyone read: this is what a shared volume, a
	// roaming profile or an operator's "make it work" fix looks like, and it is
	// the only way to see whether inheritance widened anything.
	base := wideParent(t, "root")
	data := filepath.Join(base, "data")

	store, err := memory.Open(data)
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	t.Logf("data root contains: %v", mustReadDir(t, data))

	st, err := secret.NewStore(data)
	if err != nil {
		t.Fatalf("secret.NewStore: %v", err)
	}
	if err := st.Store("dpapi:endtoend", "a-secret-nobody-else-should-try"); err != nil {
		t.Fatal(err)
	}

	// Probe while the store is open: wisp.db-wal and wisp.db-shm exist only in
	// this window, and they are written by SQLite, which has never heard of
	// winsec - which is precisely why AC#2 chose the directory route.
	for _, name := range []string{"wisp.db", "wisp.db-wal", "wisp.db-shm"} {
		p := filepath.Join(data, name)
		if _, statErr := os.Stat(p); statErr != nil {
			t.Logf("%s not present while open: %v", name, statErr)
			continue
		}
		sids, names := aclSIDs(t, p)
		if msg := privateACLError(sids, currentSID(t)); msg != "" {
			t.Errorf("SQLite-written %s is not private: %s (principals %v)", name, msg, names)
		}
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}

	// Sweep: every entry of the data root, recursively, with its real icacls
	// text logged, so the claim can be re-checked against the filesystem.
	sweepPrivate(t, data)
}

// sweepPrivate walks a tree and asserts privacy on every directory and file. It
// does not follow links: a link is not ours to walk through (ticket 79's rule,
// kept here for the same reason).
func sweepPrivate(t *testing.T, root string) {
	t.Helper()
	dirs := []string{root}
	checked := 0
	for len(dirs) > 0 {
		dir := dirs[len(dirs)-1]
		dirs = dirs[:len(dirs)-1]
		assertPrivateACL(t, dir)
		checked++
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Errorf("cannot walk %s: %v", dir, err)
			continue
		}
		for _, e := range entries {
			full := filepath.Join(dir, e.Name())
			info, lErr := os.Lstat(full)
			if lErr != nil {
				t.Errorf("lstat %s: %v", full, lErr)
				continue
			}
			switch {
			case info.Mode()&os.ModeSymlink != 0:
				// Never traverse a link, and never judge what is behind it.
			case e.IsDir():
				dirs = append(dirs, full)
			default:
				assertPrivateACL(t, full)
				checked++
			}
		}
	}
	t.Logf("sweep checked %d entries under %s, all private", checked, root)
}
