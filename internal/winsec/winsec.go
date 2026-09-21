// Package winsec owns the one promise behind every private on-disk write:
// "the current user is the only principal that can read this".
//
// Ticket 89 (A51①) is the reason this package exists: the os.OpenFile /
// os.WriteFile mode argument (0o600, 0o700) is *decorative on Windows*. The
// platform ignores those permission bits and gives the file whatever the
// parent directory's discretionary ACL hands down, so every write in this
// repository that documented itself as owner-only has been claiming something
// the filesystem never enforced. A mode argument cannot be made to work here;
// the enforcement mechanism on NTFS is the security descriptor, and the only
// evidence that counts is what `icacls` prints - not FileInfo.Mode(), which
// reports the same fiction the mode argument was supposed to control.
//
// Two mechanisms, one per platform, behind one API:
//
//   - Windows: set an explicit, non-inherited DACL that grants the current
//     token's SID and nothing else (SE_OBJECT_DACL flags via advapi32).
//   - POSIX: the mode argument really is enforced there, so the same calls
//     resolve to plain 0600/0700 and the tests assert the real bits.
//
// Failure direction is a contract, not a detail (AC#5): when the descriptor
// cannot be applied the write is refused and the half-made file is removed.
// "Could not tighten the ACL, so write it wide anyway" is the failure mode
// this package exists to eliminate, so no function here reports a warning and
// proceeds.
//
// Naming: the package is called winsec rather than acl because what it exposes
// is a promise, not an abstraction over permission systems. There is
// deliberately no API for "some ACL" - the only shape it can produce is
// current-user-only - and the Windows-specific machinery (security
// descriptors, reparse-point-safe unlink) is where the new behavior lives. On
// POSIX there is nothing new to do, so a name that advertises the platform
// the hole is specific to is more honest than one that implies a portable ACL
// model.
package winsec

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

// ErrNotSealable reports that a path could not be made private. Callers must
// treat it as "do not put the bytes here": the failure direction of this
// package is always to tighten, so a failed seal means the write is refused,
// never that the write proceeds with inherited permissions.
var ErrNotSealable = errors.New("winsec: cannot apply a private security descriptor")

// ErrIsReparsePoint reports that the entry standing where a file or directory
// was expected is a link (symlink, junction, or other reparse point). Removal
// refuses to follow it: the link itself is unlinked, and whatever lives behind
// it is left alone. Ticket 79's A51② found that on Windows os.Remove cannot
// unlink such an entry at all when it points at a non-empty directory -
// the API resolves the target first and then reports the target's contents as
// the reason not to delete - so the naive call fails in the wide direction:
// either it clears nothing or, with RemoveAll, it deletes someone else's tree.
var ErrIsReparsePoint = errors.New("winsec: entry is a link to something else, not private data")

// PrivateFile creates (or replaces) path with data and guarantees the
// platform's strongest "current user only" placement before a single content
// byte is written. An existing file is sealed in place; a file this call
// created is unlinked again if sealing fails.
func PrivateFile(path string, data []byte, perm fs.FileMode) error {
	if err := os.WriteFile(path, data, perm); err != nil {
		return err
	}
	return SealFile(path)
}

// PrivateFileExclusive creates path with data and fails with fs.ErrExist if
// any entry already occupies the name. The exclusive create is the point: the
// artifact name is injective in the logical id (ticket 79), so a collision
// means "this id's own earlier artifact", and silently truncating it would let
// one tool call overwrite another's evidence.
func PrivateFileExclusive(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	if err := sealHandle(f); err != nil {
		_ = f.Close()
		_ = os.Remove(path)
		return err
	}
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		_ = os.Remove(path)
		return err
	}
	return f.Close()
}

// SealFile narrows an existing file's permissions to the current user.
// Implemented per platform (see sealFileWindows / sealFileUnix).
func SealFile(path string) error {
	return sealFile(path)
}

// PrivateDirAll creates path and any missing parents, sealing **every level it
// creates** before returning. Sealing after the fact is not enough: a child
// inherits the ACL it had at creation time, so a directory that was created
// wide and tightened later leaves every file written into the window behind.
// The ordering this guarantees is "sealed parent, then any child", which is
// what makes sealing a parent directory once cover files nobody sealed
// individually (SQLite's -wal/-shm, a rename temp file, a leftover retry).
//
// An existing directory is verified, not re-created: if an existing level does
// not already restrict access to the current user, PrivateDirAll reports
// ErrNotSealable rather than quietly accepting it.
func PrivateDirAll(path string, perm fs.FileMode) error {
	if err := os.MkdirAll(path, perm); err != nil {
		return err
	}
	return SealDir(path)
}

// SealDir narrows a directory's permissions to the current user, with new
// children inheriting exactly that and nothing wider.
func SealDir(path string) error {
	return sealDir(path)
}

// RemoveUnlinked deletes the entry at path without following it, so a link
// standing where private data was expected disappears while whatever lives
// behind it is untouched. It never recurses (that is ticket 79's reasoning for
// os.Remove over os.RemoveAll, and RemoveAll would be the larger hole).
//
// When the entry is a link to a non-empty directory that the platform refuses
// to unlink directly, the error wraps ErrIsReparsePoint and *mentions the
// path*: the caller's reclaim loop must surface a name it cannot clear rather
// than continue with a quota it now does not know.
func RemoveUnlinked(path string) error {
	return removeUnlinked(path)
}

func wrapPath(path string, err error) error {
	return fmt.Errorf("%w: %s: %v", ErrNotSealable, path, err)
}
