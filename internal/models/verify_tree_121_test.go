package models

// Ticket 121 AC#4, first half: R-109-2 from acceptor-ticket109b, "未点名的文件/
// 子目录不参与哈希".
//
// Ticket 109's hand-off guard re-hashes the files the signed manifest names.
// That is the right answer to "were these bytes swapped" and no answer at all to
// "was a file added that nobody named" - the walk in VerifyDir simply never saw
// it. acceptor-ticket109b proved the consequence, and the reading reproduced
// byte-for-byte in the real binary before this file's fix: an installed, fully
// hash-correct model directory accepted a stray zz-unnamed-extra.onnx and an
// extra sub/ directory and `wisp models verify` still printed 与已验签清单一致
// with rc=0.
//
// Why that is a poisonable shape rather than a tidiness complaint: the point of
// Ensure's return value is a directory a loader opens by name (SPEC-04 §7.2's
// "落 models\<id>\"), so bytes standing in that directory under a plausible name
// are exactly what the next ticket's engine would read, and nothing between here
// and there ever asked whether those bytes were in the manifest.
//
// Portable on purpose: nothing here parses ACL text, so the same case runs under
// Linux (the ticket 109 POSIX container leg is the precedent), and the
// link-shaped sibling case is in verify_tree_121_other_test.go.

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// poisonTree installs the kws fixture (four named files, one of them inside a
// named subdirectory) into a fresh store and hands back the manager and the
// install directory.
func poisonTree(t *testing.T) (store string, m *Manager, dir string) {
	t.Helper()
	srv := newCountingServer(t, map[string][]byte{
		"/kws-fixture/tiny-model.tar.bz2": mustRead(t, filepath.Join("testdata", "tiny-model.tar.bz2")),
	})
	mf := manifestWithURLs(t, srv.URL())
	store = filepath.Join(t.TempDir(), "models")
	if err := os.MkdirAll(store, 0o755); err != nil {
		t.Fatal(err)
	}
	m, err := NewManager(Options{
		DataDir: store, Manifest: mf, VerifySignature: true, Attempts: 2, BackoffBase: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	dir, err = m.Ensure(context.Background(), "kws-fixture")
	if err != nil {
		t.Fatalf("install the fixture: %v", err)
	}
	return store, m, dir
}

// TestAC41CleanInstallStillVerifies is the reverse leg. A check that always fails
// would pass every other case in this file, so this one exists to make that
// shape impossible: the SAME tree, touched by nothing, must still verify -
// including the named subdirectory dict/, which the "nothing unnamed" walk is
// not allowed to report as an intruder.
func TestAC41CleanInstallStillVerifies(t *testing.T) {
	_, m, dir := poisonTree(t)
	if err := m.VerifyInstalled("kws-fixture"); err != nil {
		t.Fatalf("AC#4 reverse leg: a clean install no longer verifies, so the unnamed-entry check "+
			"is rejecting the tree it is supposed to accept: %v (dir %s)", err, dir)
	}
	if _, err := os.Stat(filepath.Join(dir, "dict", "inner.txt")); err != nil {
		t.Fatalf("the fixture did not install the named subdirectory, so this case never tested it: %v", err)
	}
	t.Logf("clean install verifies: %s holds exactly the manifest's 4 named files plus dict/", dir)
}

// TestAC42UnnamedFileAndDirAreRefused is the R-109-2 leg, in the three shapes it
// was measured sliding through: a file at the root, a directory at the root, and
// a file inside a directory the manifest does name.
func TestAC42UnnamedFileAndDirAreRefused(t *testing.T) {
	cases := []struct {
		name  string
		setup func(t *testing.T, dir string)
		want  string
	}{
		{
			name: "loose file next to the model",
			setup: func(t *testing.T, dir string) {
				if err := os.WriteFile(filepath.Join(dir, "zz-unnamed-extra.onnx"), []byte("poison"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			want: "zz-unnamed-extra.onnx",
		},
		{
			name: "empty directory nobody named",
			setup: func(t *testing.T, dir string) {
				if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o755); err != nil {
					t.Fatal(err)
				}
			},
			want: "sub/",
		},
		{
			name: "file added inside a named directory",
			setup: func(t *testing.T, dir string) {
				if err := os.WriteFile(filepath.Join(dir, "dict", "extra.txt"), []byte("poison"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			want: "dict/extra.txt",
		},
		{
			// The loader-relevant one: the manifest's own model.onnx is still
			// byte-correct here, so a hash-only check has nothing to report.
			name: "plausible sibling of the model file",
			setup: func(t *testing.T, dir string) {
				if err := os.WriteFile(filepath.Join(dir, "model.int8.onnx"), []byte("newer bytes nobody signed"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			want: "model.int8.onnx",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, m, dir := poisonTree(t)
			tc.setup(t, dir)
			err := m.VerifyInstalled("kws-fixture")
			if err == nil {
				t.Fatalf("AC#4/R-109-2: the hand-off verified a model directory holding bytes the signed "+
					"manifest does not name (%s). Every named file still hashes correctly, which is exactly "+
					"why hashing the named subset cannot see this.", tc.want)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("the refusal must name the offender %q or the operator cannot act on it, got: %v", tc.want, err)
			}
			t.Logf("refused and named: %v", err)
		})
	}
}

// TestAC43HandOffRefusesAnUnnamedFile is the same hole seen from the only
// production caller there is: DownloadingBridge.Run (internal/models/bridge.go),
// i.e. the guard ticket 109 built. Without AC#4's fix the walk below announced
// the model available, because Ensure's re-download path and the hand-off
// re-verification both looked only at named files.
func TestAC43HandOffRefusesAnUnnamedFile(t *testing.T) {
	store, _, dir := poisonTree(t)
	srv := newCountingServer(t, map[string][]byte{
		"/kws-fixture/tiny-model.tar.bz2": mustRead(t, filepath.Join("testdata", "tiny-model.tar.bz2")),
	})
	mf := manifestWithURLs(t, srv.URL())
	mgr, err := NewManager(Options{
		DataDir: store, Manifest: mf, VerifySignature: true, Attempts: 2, BackoffBase: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "extra.bin"), []byte("added after install"), 0o644); err != nil {
		t.Fatal(err)
	}
	state, err := WireDownloading(mgr, newWalkMachine()).Run(context.Background(), "kws-fixture")
	if err == nil {
		t.Fatalf("AC#4/R-109-2: the hand-off announced a model available whose directory holds an unnamed "+
			"file (state=%s, dir=%s)", state, dir)
	}
	if !strings.Contains(err.Error(), "extra.bin") {
		t.Errorf("the refusal names nothing actionable: %v", err)
	}
	t.Logf("hand-off refused at %s: %v", state, err)
}
