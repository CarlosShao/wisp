package models

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// Ticket 95 AC#3, the "deliberately not sealed" side of the ruling.
//
// Every write site in this package (downloader.go's staging dir, its per
// artifact subdirectory, archive.go's extraction tree and the install dir)
// still creates with 0o755/0o644 and is NOT passed through winsec. That is a
// decision, not an omission, and it rests on exactly two facts:
//
//  1. the bytes are public - every file is hash pinned in the signed manifest
//     (doc.go: C29 verifies the manifest offline against a build-time key
//     before any network use, and hashes come only from the manifest), so a
//     reader on another local account learns nothing it could not fetch from
//     the mirror itself;
//  2. the read path verifies per file, at read time - Ensure's cache hit runs
//     VerifyDir over InstalledFiles() before handing the directory back.
//
// The test below pins fact 2 at its stated coverage, file by file. If any
// installed file ever stops being checked, the premise of "no leak here" is
// gone and this class has to be re-ruled and sealed; that failure is the point
// of the test, not a side effect.

func TestAC3EveryInstalledFileIsReverifiedAtReadTime(t *testing.T) {
	srv := newCountingServer(t, map[string][]byte{
		"/kws-fixture/tiny-model.tar.bz2": mustRead(t, filepath.Join("testdata", "tiny-model.tar.bz2")),
	})
	m := manifestWithURLs(t, srv.URL())
	mgr, _ := newTestManager(t, m, nil, nil)

	dir, err := mgr.Ensure(context.Background(), "kws-fixture")
	if err != nil {
		t.Fatalf("archive install failed: %v", err)
	}
	entry, err := m.FindModel("kws-fixture")
	if err != nil {
		t.Fatal(err)
	}
	want := entry.InstalledFiles()
	if len(want) != len(tinyArchiveMembers) {
		t.Fatalf("coverage of the signed manifest is %d files, the fixture installs %d",
			len(want), len(tinyArchiveMembers))
	}
	for i, f := range want {
		p := filepath.Join(dir, filepath.FromSlash(f.Path))
		orig, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("installed file %s missing: %v", f.Path, err)
		}
		tampered := append([]byte(nil), orig...)
		tampered[0] ^= 0x01
		if err := os.WriteFile(p, tampered, 0o600); err != nil {
			t.Fatalf("tamper write: %v", err)
		}
		if err := mgr.VerifyDir(entry, dir); err == nil {
			t.Errorf("tampering installed file %d/%d (%s) went undetected: the read path "+
				"does not cover every installed file, so ticket 95's no-seal ruling for this "+
				"class is void and it must be re-ruled", i+1, len(want), f.Path)
		}
		if err := os.WriteFile(p, orig, 0o600); err != nil {
			t.Fatalf("restore write: %v", err)
		}
	}
	if err := mgr.VerifyDir(entry, dir); err != nil {
		t.Fatalf("the clean install dir must verify: %v", err)
	}
	// Ensure itself must take the same branch: one more Ensure after the
	// round trip proves the dir is still accepted, i.e. the loop above really
	// was the tamper test and not a leftover half-installed tree.
	if _, err := mgr.Ensure(context.Background(), "kws-fixture"); err != nil {
		t.Fatalf("Ensure after restore: %v", err)
	}
}
