//go:build windows

package winsec

import (
	"io/fs"
	"os"
)

// sealFile is the checkpoint state of this package, written deliberately
// before any security-descriptor code exists: os.Chmod is exactly what the
// repository has been doing all along (the mode argument is already applied by
// os.WriteFile, re-applying it here changes nothing), so the tests in this
// directory are expected to be RED while it is the only mechanism. The real
// implementation replaces this file's body; see ticket 89 AC#1/AC#3 for the
// icacls evidence that says so.
func sealFile(path string) error { return os.Chmod(path, 0o600) }

// sealDir is the same placeholder for directories.
func sealDir(path string) error { return os.Chmod(path, 0o700) }

func sealHandle(f *os.File) error { return f.Chmod(0o600) }

// removeUnlinked is the naive removal, i.e. os.Remove with no reparse-point
// handling: the A51② behavior this package has to fix.
func removeUnlinked(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

var _ fs.FileMode = 0o600 // keep the import used while the body is a placeholder
