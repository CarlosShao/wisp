//go:build windows

package risk

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Ticket 105, AC#3 — the BEHAVIOR case for ticket 102's disposition row 4
// (internal/risk/syncdirs.go resolveTarget's SECOND Resolve, the one on the
// nearest existing ancestor, plus its Actable() leg).
//
// What row 4 pinned before this file was a symbol: a test asserted that the
// ancestor result is read through Actable(). It never showed that the leg does
// anything. This file shows it: a write target whose ANCESTOR is a real OS
// junction, registered as a reparse_point_exception (so the first Resolve
// cannot deny it), and a not-yet-existing leaf (so resolveHandle on the whole
// spelling fails and the ancestor walk is what produces the comparison
// candidate). The only way the verdict can name the sync root is for the
// ancestor to have been re-resolved THROUGH the mount point onto the real tree.
//
// Mutation the acceptance side should run: in syncdirs.go replace lines
// 222-232 with `anceCanon := anc` (drop the re-resolve, its deny leg and its
// Actable leg together). The candidate then keeps the JUNCTION's spelling,
// isUnder(root) misses, and the first case below goes red on the verdict
// (Sync=false, i.e. the write is released un-suspected) - not on a symbol.
//
// Everything created here lives in t.TempDir(); nothing outside it is touched,
// and the junction is removed before its target directory is.

// junctionTargetAndLink builds <tmp>/tree (the real directory), <tmp>/proj as a
// REAL junction onto it, and <tmp>/tree/Notes inside the target.
func junctionTargetAndLink(t *testing.T) (target, link, notes string) {
	t.Helper()
	target = t.TempDir()
	link = filepath.Join(t.TempDir(), "proj")
	notes = filepath.Join(target, "Notes")
	if err := os.MkdirAll(notes, 0o755); err != nil {
		t.Fatal(err)
	}
	mkJunction(t, link, target)
	// Precondition, so a red below cannot be blamed on mklink silently failing:
	// the link exists, and stat-ing THROUGH it reaches the target.
	st, err := os.Stat(filepath.Join(link, "Notes"))
	if err != nil || !st.IsDir() {
		t.Fatalf("junction %s -> %s is not traversable: %v", link, target, err)
	}
	t.Cleanup(func() { _ = os.Remove(link) }) // remove the link, never its target
	return target, link, notes
}

// TestSyncAncestorLegReResolvesThroughJunctionToRealTree is the pinned
// behaviour: the ancestor leg anchors the candidate on the tree the junction
// points AT, which is the registered sync root, and not on the link's spelling.
func TestSyncAncestorLegReResolvesThroughJunctionToRealTree(t *testing.T) {
	target, link, notes := junctionTargetAndLink(t)
	home := t.TempDir() // deliberately unrelated: the suspect net must not be what answers

	p := NewProvenance(ProvOptions{
		NoProbe: true, HomeDir: home,
		SyncRoots:         []SyncRoot{{Provider: "Ticket105Sync", Path: target, Source: "env"}},
		ReparseExceptions: []string{link},
	})
	if !p.SyncDetectionComplete() {
		t.Fatalf("precondition: the sync root must be confirmed+canonical, otherwise the "+
			"blanket under-profile net answers every case and this test proves nothing (roots=%v)",
			p.SyncRoots())
	}

	// Entry condition of the ancestor leg, asserted rather than assumed: the
	// whole spelling does NOT resolve (so resolveTarget must fall through to
	// deepestExistingAncestor + the second Resolve) and it DOES traverse a
	// reparse point that is exempted.
	written := filepath.Join(link, "Notes", "new.md")
	res, err := Resolve(written, []string{link})
	if err != nil {
		t.Fatalf("precondition: exempted junction must not be denied by the first Resolve: %v", err)
	}
	if res.Resolved {
		t.Fatalf("precondition: %s must NOT be handle-resolvable as a whole, otherwise the ancestor leg never runs", written)
	}
	if !res.Reparse {
		t.Fatalf("precondition: %s must traverse an exempted reparse point, otherwise this is not the ancestor case", written)
	}

	st := p.IsSyncPath(written)
	if !st.Sync {
		t.Errorf("IsSyncPath(%q) = Sync:false (why=%q): the ancestor leg did not anchor this write "+
			"onto the sync root %q, so a path that lands inside a sync tree THROUGH a junction would be "+
			"released un-suspected", written, st.Why, target)
	} else {
		if got := filepath.Clean(st.Root.Path); !strings.EqualFold(got, filepath.Clean(target)) {
			t.Errorf("matched root = %q, want the junction TARGET %q: matching through the link's own "+
				"spelling means the ancestor was never re-resolved", st.Root.Path, target)
		}
		if st.Root.Provider != "Ticket105Sync" || st.Root.Source != "env" {
			t.Errorf("matched root = %+v, want the registered env root (not the suspect fallback)", st.Root)
		}
	}
	// The same target reached by its real spelling must agree, which is what
	// makes this one verdict rather than two coincidences.
	if direct := p.IsSyncPath(filepath.Join(notes, "new.md")); !direct.Sync ||
		!strings.EqualFold(filepath.Clean(direct.Root.Path), filepath.Clean(target)) {
		t.Errorf("IsSyncPath(direct spelling) = %+v, want a match on %q", direct, target)
	}
}

// TestSyncAncestorLegNotExemptedStaysFailClosed is the other half of the same
// leg: the moment the junction is NOT registered, the verdict must come back as
// fail-closed sync-suspect naming the reparse denial, and it must not leak a
// matched root.
func TestSyncAncestorLegNotExemptedStaysFailClosed(t *testing.T) {
	target, link, _ := junctionTargetAndLink(t)
	home := t.TempDir()
	p := NewProvenance(ProvOptions{
		NoProbe: true, HomeDir: home,
		SyncRoots: []SyncRoot{{Provider: "Ticket105Sync", Path: target, Source: "env"}},
	})
	st := p.IsSyncPath(filepath.Join(link, "Notes", "new.md"))
	if !st.Sync {
		t.Fatal("a non-exempted reparse traversal must never answer Sync:false")
	}
	if st.Root.Source != "suspect-fallback" {
		t.Errorf("verdict root = %+v, want the fail-closed suspect fallback and no matched sync root", st.Root)
	}
	if !strings.Contains(st.Why, "reparse") {
		t.Errorf("Why = %q, want it to name the reparse denial the operator would act on", st.Why)
	}
}
