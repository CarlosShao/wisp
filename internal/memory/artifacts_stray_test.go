package memory

// Ticket 79, defect 2 — nothing inside artifacts\ may be invisible.
//
// listArtifactsDir used to do `if e.IsDir() { continue }`. That single line made
// a subdirectory invisible to ListArtifacts, to PurgeArtifacts AND to the 500 MB
// quota at once, which is how "flat directory, 500MB cap" stops describing the
// disk: bytes parked one level down never counted, never got LRU'd, and survived
// the privacy "one-click clear". Ticket 76 even recorded the symptom honestly
// (its `nested\` canary outlived a purge).
//
// The shape chosen here is the ticket's option (a): the listing measures
// non-conforming entries recursively, the quota counts them in the same LRU queue
// as regular artifacts, and both the quota job and a purge reclaim them. Option
// (b) - report an error and stop - was rejected because it fails in the two ways
// this ticket exists to close: a quota that can be *told* about an entry but
// cannot evict it is still unenforceable, and a privacy clear that returns an
// error while the bytes sit there is worse than one that removes them.
//
// Every test in this file is red against the `continue`: the property asserted is
// always about bytes or entries the old code could not see.

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

// stray79Plant writes a file (creating parents) and back-dates it, returning its
// path. The back-dating matters: LRU order here is decided by mtime, and a test
// that left everything at "now" would be asserting nothing about the queue.
func stray79Plant(t *testing.T, path string, size int, age time.Duration) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, make([]byte, size), 0o644); err != nil {
		t.Fatal(err)
	}
	past := time.Now().Add(-age)
	if err := os.Chtimes(path, past, past); err != nil {
		t.Fatal(err)
	}
}

// stray79TreeBytes measures the WHOLE artifacts tree independently of the code
// under test: this is the number the quota is supposed to be bounded by, so it
// must not come from the same listing that is being audited.
func stray79TreeBytes(t *testing.T, dir string) (int64, int) {
	t.Helper()
	var bytes int64
	var entries int
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if p != dir {
				entries++
			}
			return nil
		}
		fi, iErr := d.Info()
		if iErr != nil {
			return iErr
		}
		bytes += fi.Size()
		entries++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return bytes, entries
}

// TestStraySubdirectoryCannotHideBytesFromTheQuota is AC#4's headline: 4 KB of
// payload parked one directory down used to be invisible to the cap, so a quota
// of 20 bytes was "satisfied" by a 4 KB tree. Under the old `continue` this fails
// twice over - nothing is reclaimed, and the tree still holds the bytes.
func TestStraySubdirectoryCannotHideBytesFromTheQuota(t *testing.T) {
	logger, buf := bufLogger(t)
	s := openTestStore(t, WithLogger(logger))
	ctx := context.Background()
	arts := s.ArtifactsDir()

	stray79Plant(t, filepath.Join(arts, "flat.txt"), 10, time.Hour)
	stray79Plant(t, filepath.Join(arts, "stray", "deep", "hidden.bin"), 4096, 2*time.Hour)

	before, beforeEntries := stray79TreeBytes(t, arts)
	if before < 4106 {
		t.Fatalf("fixture holds %d bytes (%d entries), want the 4 KB stray plus the flat file", before, beforeEntries)
	}
	if got, err := s.ListArtifacts(ctx); err != nil || len(got) != 1 {
		t.Fatalf("ListArtifacts = %d entries, %v; want the 1 conforming artifact", len(got), err)
	}
	if !strings.Contains(buf.String(), "non-conforming") || !strings.Contains(buf.String(), "4096") {
		t.Errorf("the stray subtree was not reported: %s", buf.String())
	}

	// Quota 20 against 4106 bytes on disk: enforcement must bring the TREE under
	// the cap, not just the slice of it the listing happens to see.
	files, freed, err := s.enforceArtifactsQuota(ctx, 20)
	if err != nil {
		t.Fatalf("enforceArtifactsQuota: %v", err)
	}
	if freed < 4096 {
		t.Errorf("quota enforcement freed %d bytes against a 4096-byte hidden subtree (files=%d): "+
			"the stray entry is still invisible to the cap", freed, files)
	}
	after, _ := stray79TreeBytes(t, arts)
	if after > 20 {
		t.Errorf("artifacts tree still holds %d bytes against a 20-byte quota (freed=%d): "+
			"the cap was applied to a listing that could not see them", after, freed)
	}
	if _, err := os.Stat(filepath.Join(arts, "stray")); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("the stray directory survived quota enforcement: err=%v", err)
	}
	if !strings.Contains(buf.String(), "artifacts LRU cleanup") {
		t.Errorf("LRU cleanup not logged: %s", buf.String())
	}
}

// TestStraySubdirectorySharesTheQuotasLRUQueue pins the merged-queue decision: a
// stray dir is not a special class that gets immunity, it competes for eviction
// by the same last-write rule. Here the hidden bytes are the OLDEST thing on
// disk, so they must be the first thing reclaimed and the fresh artifact must
// survive.
func TestStraySubdirectorySharesTheQuotasLRUQueue(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	arts := s.ArtifactsDir()

	stray79Plant(t, filepath.Join(arts, "stray", "old.bin"), 500, 5*time.Hour)
	stray79Plant(t, filepath.Join(arts, "fresh.txt"), 100, time.Minute)
	// The stray's LRU proxy is the newest write ANYWHERE under it, and a
	// directory's own mtime moves when an entry is added - so back-dating the
	// fixture means back-dating the directory too, or the tree really is the
	// freshest thing on disk and the ordering below is not testing what it names.
	past := time.Now().Add(-5 * time.Hour)
	if err := os.Chtimes(filepath.Join(arts, "stray"), past, past); err != nil {
		t.Fatal(err)
	}

	// 600 bytes on disk, quota 150: 450 over, so the oldest 500-byte entry (the
	// stray tree) goes and the 100-byte fresh artifact stays.
	files, freed, err := s.enforceArtifactsQuota(ctx, 150)
	if err != nil {
		t.Fatalf("enforceArtifactsQuota: %v", err)
	}
	if files != 1 || freed != 500 {
		t.Errorf("enforcement reclaimed (%d entries, %d bytes), want the single oldest 500-byte stray",
			files, freed)
	}
	if _, err := os.Stat(filepath.Join(arts, "fresh.txt")); err != nil {
		t.Errorf("the newest artifact was evicted instead of the older stray: %v", err)
	}
	if _, err := os.Stat(filepath.Join(arts, "stray")); err == nil {
		t.Error("the oldest entry on disk was the stray dir and it survived: the queue is not merged")
	}
	after, _ := stray79TreeBytes(t, arts)
	if after > 150 {
		t.Errorf("tree holds %d bytes over a 150-byte quota", after)
	}
}

// TestPurgeArtifactsReclaimsStraySubdirectory is the privacy half: one-click
// clear must not leave bytes behind because they sat one level down. Containment
// is asserted in the same breath - reclaiming recursively must stay inside the
// artifacts dir.
func TestPurgeArtifactsReclaimsStraySubdirectory(t *testing.T) {
	logger, buf := bufLogger(t)
	s := openTestStore(t, WithLogger(logger))
	ctx := context.Background()
	arts := s.ArtifactsDir()
	dataDir := s.Dir()

	stray79Plant(t, filepath.Join(arts, "nested", "sub", "a.bin"), 1024, time.Hour)
	stray79Plant(t, filepath.Join(arts, "top.txt"), 64, time.Hour)
	// Canaries that a recursive reclaim has no business touching.
	stray79Plant(t, filepath.Join(dataDir, "canary-outside.txt"), 8, time.Hour)
	stray79Plant(t, filepath.Join(filepath.Dir(arts), "sibling", "keep.txt"), 8, time.Hour)

	removed, err := s.PurgeArtifacts(ctx)
	if err != nil {
		t.Fatalf("PurgeArtifacts: %v", err)
	}
	if removed != 2 {
		t.Errorf("PurgeArtifacts removed %d entries, want 2 (the file + the stray subtree)", removed)
	}
	if left, entries := stray79TreeBytes(t, arts); left != 0 || entries != 0 {
		t.Errorf("after a purge the artifacts tree still holds %d bytes in %d entries", left, entries)
	}
	for _, keep := range []string{
		filepath.Join(dataDir, "canary-outside.txt"),
		filepath.Join(filepath.Dir(arts), "sibling", "keep.txt"),
	} {
		if _, err := os.Stat(keep); err != nil {
			t.Errorf("recursive purge destroyed %q, which is outside artifacts\\: %v", keep, err)
		}
	}
	if _, err := os.Stat(arts); err != nil {
		t.Errorf("the artifacts dir itself must survive its own purge: %v", err)
	}
	if !strings.Contains(buf.String(), "stray_entries=1") {
		t.Errorf("the purge log must say it reclaimed a stray: %s", buf.String())
	}
}

// TestStrayRemovalDoesNotFollowLinks pins why removal walks with os.Remove and
// never os.RemoveAll: os.Remove does not recurse, so reclaiming a listed entry
// cannot be steered into deleting a tree the store does not own. The assertions
// are branch-agnostic on purpose - a symlink may land in the stray set, in the
// artifact set (an entry whose Type is a link is not a directory), or nowhere at
// all (Windows needs SeCreateSymbolicLinkPrivilege or developer mode); whatever
// the outcome, the target outside artifacts\ has to be standing at the end.
func TestStrayRemovalDoesNotFollowLinks(t *testing.T) {
	s := openTestStore(t)
	arts := s.ArtifactsDir()
	outside := filepath.Join(s.Dir(), "not-an-artifact")
	if err := os.MkdirAll(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	stray79Plant(t, filepath.Join(outside, "precious.txt"), 32, time.Hour)

	link := filepath.Join(arts, "stray-link")
	if err := os.Symlink(outside, link); err != nil {
		t.Logf("symlink unavailable in this environment (%v): the same assertions still run", err)
	}
	stray79Plant(t, filepath.Join(arts, "stray-dir", "in.bin"), 128, time.Hour)

	listed, strays, err := listArtifactsDir(s.artifactsDir)
	if err != nil {
		t.Fatalf("listArtifactsDir: %v", err)
	}
	// A directory is never an artifact...
	for _, a := range listed {
		if a.Name == "stray-dir" {
			t.Errorf("the stray dir was listed as an artifact: %+v", listed)
		}
		// ...and nothing nested may ride along in a name, because the public
		// delete route only accepts bare names (ticket 76's guard).
		if filepath.Base(a.Name) != a.Name {
			t.Errorf("listArtifactsDir exposed the nested name %q", a.Name)
		}
	}
	var strayNames []string
	for _, st := range strays {
		strayNames = append(strayNames, st.name)
		if filepath.Base(st.name) != st.name {
			t.Errorf("stray name %q is not a single component", st.name)
		}
	}
	if !slices.Contains(strayNames, "stray-dir") {
		t.Fatalf("stray-dir is not in the non-conforming set %v: the listing is blind again", strayNames)
	}
	for _, st := range strays {
		if err := s.removeStray(st.name); err != nil {
			t.Errorf("removeStray(%q): %v", st.name, err)
		}
	}
	for _, a := range listed {
		if err := s.DeleteArtifact(context.Background(), a.Name); err != nil {
			t.Logf("DeleteArtifact(%q) after the link experiment: %v", a.Name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(outside, "precious.txt")); err != nil {
		t.Errorf("reclaiming an entry of artifacts\\ took outside/precious.txt with it: %v", err)
	}
	if _, err := os.Stat(outside); err != nil {
		t.Errorf("reclaiming an entry of artifacts\\ removed the directory it pointed at: %v", err)
	}
	if _, err := os.Stat(filepath.Join(arts, "stray-dir")); err == nil {
		t.Error("stray-dir survived removeStray")
	}
}
