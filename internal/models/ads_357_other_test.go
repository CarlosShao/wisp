//go:build !windows

package models

// Ticket 357's POSIX side, and it is a CONTROL LEG, not a tooth.
//
// Stated plainly because the distinction is the whole point: POSIX has no
// alternate data streams, so listNamedStreams here is a hard-coded empty answer
// and cannot fail. What this case can establish is only that the report stays
// silent where there is nothing to report - the caller walks a real installed
// tree, the enumeration returns nothing, and no sentence is produced. It says
// NOTHING about whether the Windows enumeration works; that tooth lives in
// ads_357_windows_test.go on real NTFS and cannot be moved onto this leg.
//
// It is here rather than left out for the reason ticket 121's assembly gate gives
// for its own shape: internal/models is in the portable CI scope, and a caller
// added on one platform with no case on the other is an edge nobody looks at
// until it breaks the build. A control leg is cheaper than that surprise, as long
// as nobody mistakes it for evidence.

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func Test357POSIXHasNoStreamsToReport(t *testing.T) {
	srv := newCountingServer(t, map[string][]byte{
		"/kws-fixture/tiny-model.tar.bz2": mustRead(t, filepath.Join("testdata", "tiny-model.tar.bz2")),
	})
	mf := manifestWithURLs(t, srv.URL())
	store := filepath.Join(t.TempDir(), "models")
	if err := os.MkdirAll(store, 0o755); err != nil {
		t.Fatal(err)
	}
	mgr, err := NewManager(Options{
		DataDir: store, Manifest: mf, VerifySignature: true, Attempts: 2, BackoffBase: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	dir, err := mgr.Ensure(context.Background(), "kws-fixture")
	if err != nil {
		t.Fatalf("install the fixture: %v", err)
	}

	gotDir, streams, err := mgr.UnnamedStreams("kws-fixture")
	if err != nil {
		t.Fatalf("ticket 357: the enumeration failed on this platform, which would make the report print a "+
			"failure sentence on every clean install: %v", err)
	}
	if gotDir != dir {
		t.Errorf("ticket 357: enumerated %s, want the install directory %s", gotDir, dir)
	}
	if len(streams) != 0 {
		t.Fatalf("ticket 357: POSIX answered with %v; there is no way to spell a second stream here, so a "+
			"non-empty answer means the !windows half is not the empty one it claims to be", streams)
	}
	// And the tree still verifies: the report never becomes a refusal.
	if err := mgr.VerifyInstalled("kws-fixture"); err != nil {
		t.Errorf("ticket 357: enumerating streams changed the verdict of a clean install: %v", err)
	}
}
