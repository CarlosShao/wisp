//go:build !windows

package winsec

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// This file is AC#6's answer to "non-Windows platforms must not pass silently".
// On POSIX the mode argument is enforced, so these tests measure the real bits
// rather than declaring the ACL story inapplicable. A skip would be the wrong
// shape here: the promise is the same on both platforms, only the mechanism
// differs, and a filesystem that does not store mode bits has to fail (see
// applyDescriptorPOSIX's read-back).

func TestPOSIXPrivateFileIsReally0600(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "artifact.txt")
	// 0o644 in, 0600 out: the promise is current-user-only, not "whatever the
	// caller asked for", and on this platform that costs nothing.
	if err := PrivateFile(p, []byte("tool output"), 0o644); err != nil {
		t.Fatalf("PrivateFile: %v", err)
	}
	info, err := os.Stat(p)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Errorf("artifact mode is %v, want -rw-------", got)
	}
}

func TestPOSIXPrivateDirIsReally0700(t *testing.T) {
	root := filepath.Join(t.TempDir(), "data", "artifacts")
	if err := PrivateDirAll(root, 0o755); err != nil {
		t.Fatalf("PrivateDirAll: %v", err)
	}
	for _, p := range []string{filepath.Dir(root), root} {
		info, err := os.Stat(p)
		if err != nil {
			t.Fatal(err)
		}
		if got := info.Mode().Perm(); got != 0o700 {
			t.Errorf("%s mode is %v, want drwx------", filepath.Base(p), got)
		}
	}
}

// TestPOSIXSymlinkAtArtifactPositionIsNotRecursed is the CI/Linux-side
// equivalent of the Windows junction leg in reparse_windows_test.go: a symlink
// to a non-empty directory standing where an artifact was expected. POSIX
// constructs it without any privilege, so the construction is stronger here
// than on Windows, while the failure the Windows leg pins (os.Remove cannot
// unlink such an entry) does not exist on this platform - unlink operates on
// the link.
func TestPOSIXSymlinkAtArtifactPositionIsNotRecursed(t *testing.T) {
	outside := filepath.Join(t.TempDir(), "someone-elses-tree")
	deep := filepath.Join(outside, "sub")
	if err := os.MkdirAll(deep, 0o700); err != nil {
		t.Fatal(err)
	}
	innocent := filepath.Join(deep, "keep-me.txt")
	if err := os.WriteFile(innocent, []byte("not ours to delete"), 0o600); err != nil {
		t.Fatal(err)
	}

	root := filepath.Join(t.TempDir(), "data")
	if err := PrivateDirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "artifact")
	if err := os.Symlink(outside, link); err != nil {
		t.Fatalf("constructing the symlink (the thing Windows needs a privilege for): %v", err)
	}
	info, err := os.Lstat(link)
	if err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("%s is not a symlink (info=%v, err=%v)", link, info, err)
	}

	if err := RemoveUnlinked(link); err != nil {
		t.Fatalf("RemoveUnlinked: %v", err)
	}
	if _, err := os.Lstat(link); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("the link survived: %v", err)
	}
	if _, err := os.Stat(innocent); err != nil {
		t.Errorf("TARGET DELETED - removal followed the symlink: %v", err)
	}
}

// TestPOSIXMissingFileIsNotAnError keeps the reclaim loop's idempotence claim
// platform-neutral.
func TestPOSIXMissingFileIsNotAnError(t *testing.T) {
	if err := RemoveUnlinked(filepath.Join(t.TempDir(), "gone")); err != nil {
		t.Errorf("RemoveUnlinked of a non-existent entry: %v", err)
	}
}
