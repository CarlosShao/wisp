//go:build !windows

package winsec

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"syscall"
)

// applyDescriptor is the same seam the Windows file exposes, and it is what the
// AC#5 failure-injection test replaces on either platform.
//
// On POSIX the mode argument is real, so this is not the decorative call it
// would be on Windows - but it is still not trusted: a filesystem that stores
// no mode bits (FAT/exFAT, some NFS exports with all_squash) accepts chmod and
// then reports something else, and the verification below turns that into
// ErrNotSealable instead of a false promise. That is the AC#6 requirement in
// code: the non-Windows leg measures something, it does not skip.
var applyDescriptor = applyDescriptorPOSIX

func applyDescriptorPOSIX(path string, dir bool) error {
	mode := fs.FileMode(0o600)
	if dir {
		mode = 0o700
	}
	if err := os.Chmod(path, mode); err != nil {
		return wrapPath(path, err)
	}
	info, err := os.Stat(path)
	if err != nil {
		return wrapPath(path, err)
	}
	if got := info.Mode().Perm(); got != mode {
		return wrapPath(path, fmt.Errorf("mode is %v, want %v (%s stores no usable permission bits)",
			got, mode, fstypeOf(path)))
	}
	return nil
}

// sealHandle is the pre-content seal used by PrivateFile and
// PrivateFileExclusive: same mechanism as SealFile, by name.
func sealHandle(f *os.File) error { return applyDescriptor(f.Name(), false) }

func fstypeOf(path string) string {
	var st syscall.Statfs_t
	if err := syscall.Statfs(path, &st); err != nil {
		return "unknown filesystem"
	}
	return fmt.Sprintf("fs type 0x%x", st.Type)
}

func sealFile(path string) error { return applyDescriptor(path, false) }

// platformVerifyPlacement mirrors the DEFERRED leg of internal/risk's
// pathresolver_other.go: on this platform C26 itself resolves lexically and does
// not detect reparse traversal, so the floor here stops at the portable shape
// checks in resolve.go rather than inventing a stricter rule that would make
// production refuse paths its own PathResolver accepts (and would break a macOS
// install whose /tmp or HOME sits behind a symlink). The gap is stated, not
// papered over: when that DEFERRED leg lands, this function is where it plugs in.
func platformVerifyPlacement(path string) (string, error) { return path, nil }

// sealDir narrows a directory. It deliberately does not walk the existing
// subtree the way the Windows implementation does: on POSIX a child never
// inherited its parent's mode in the first place, so "seal the tree" here would
// be a chmod over files whose permissions somebody else set on purpose, and the
// hole being closed does not exist on this platform.
func sealDir(path string) error { return applyDescriptor(path, true) }

// removeUnlinked needs no special care on POSIX: unlink operates on the link
// itself and never follows it, so os.Remove already has the non-recursive,
// non-following behavior the Windows implementation has to work for.
func removeUnlinked(path string) error {
	if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return nil
}
