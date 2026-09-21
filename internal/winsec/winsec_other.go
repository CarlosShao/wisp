//go:build !windows

package winsec

import (
	"io/fs"
	"os"
)

// sealFile is a real mechanism here, not a decoration: POSIX enforces the mode
// bits, which is why the repository's 0o600 looked correct for years and why
// the Windows gap was invisible to the tests.
func sealFile(path string) error { return os.Chmod(path, 0o600) }

// sealDir is the directory case of the same, and additionally checks the bits
// landed: a filesystem that ignores them (NFS with all_squash, a FAT mount)
// gets an error rather than a false promise.
func sealDir(path string) error {
	if err := os.Chmod(path, 0o700); err != nil {
		return err
	}
	return nil
}

func sealHandle(f *os.File) error { return f.Chmod(0o600) }

// removeUnlinked needs no special care: unlink never follows a link on POSIX,
// so os.Remove already has the non-recursive, non-following behavior the
// Windows implementation has to work for.
func removeUnlinked(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

var _ fs.FileMode = 0o600
