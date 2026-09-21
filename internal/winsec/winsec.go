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
// Two mechanisms, one API:
//
//   - Windows (winsec_windows.go): an explicit DACL carrying only the current
//     user's SID plus SYSTEM and Administrators, marked SE_DACL_PROTECTED so no
//     inherited grant can widen it, read back and verified before any byte is
//     written.
//   - POSIX (winsec_other.go): the mode argument really is enforced there, so
//     the same calls resolve to 0600/0700 - and are checked to have landed,
//     which is what makes a filesystem that ignores them (FAT, some NFS mounts)
//     an error instead of a false promise.
//
// Failure direction is a contract, not a detail (AC#5): when the descriptor
// cannot be applied the write is refused and the half-made file is removed.
// "Could not tighten the ACL, so write it wide anyway" is the failure mode this
// package exists to eliminate, so nothing here reports a warning and proceeds.
//
// Naming: winsec rather than acl, because what it exposes is a promise, not an
// abstraction over permission systems. There is deliberately no API for "some
// ACL" - the only shape it can produce is current-user-only - and the new
// machinery (security descriptors, reparse-point-safe unlink) is Windows
// specific. On POSIX there is nothing new to do, so a name that advertises the
// platform the hole is specific to is more honest than one implying a portable
// ACL model.
package winsec

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// ErrNotSealable reports that a path could not be made private. Callers must
// treat it as "do not put the bytes here": the failure direction of this
// package is always to tighten, so a failed seal means the write is refused,
// never that the write proceeds with inherited permissions.
var ErrNotSealable = errors.New("winsec: cannot apply a private security descriptor")

// ErrIsReparsePoint reports that the entry standing where a file or directory
// was expected is a link (symlink, junction, or other reparse point). Removal
// refuses to follow it: the link itself is unlinked and whatever lives behind
// it is left alone. Ticket 79's A51② found that on Windows os.Remove cannot
// unlink such an entry at all when the target is a non-empty directory - the
// API resolves the target first and then cites the target's contents as the
// reason not to delete - so the naive call fails in the wide direction: either
// it clears nothing or, with RemoveAll, it deletes someone else's tree.
var ErrIsReparsePoint = errors.New("winsec: entry is a link to something else, not private data")

// PrivateFile creates or replaces path with data and guarantees the platform's
// strongest "current user only" placement *before* a single content byte is
// written. perm is what the mode means on POSIX and is ignored on Windows,
// where the descriptor decides - callers pass 0o600 for the readers of the
// source, not because Windows obeys it.
func PrivateFile(path string, data []byte, perm fs.FileMode) error {
	return privateFile(path, data, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
}

// PrivateFileExclusive creates path with data and fails with fs.ErrExist if any
// entry already occupies the name. The exclusive create is the point: artifact
// names are injective in the logical id (ticket 79), so a collision means "this
// id's own earlier artifact", and silently truncating it would let one tool call
// overwrite another's evidence. Deciding what a same-id retry means stays the
// caller's job - see agent.Spiller.Prepare's documented last-writer-wins swap.
func PrivateFileExclusive(path string, data []byte) error {
	return privateFile(path, data, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
}

// privateFile creates with the given flags, seals the empty file, and only then
// writes. A seal that fails leaves no content behind: the entry is removed
// before the error is returned, so the failure is "this artifact was not
// written", not "this artifact is world-readable".
func privateFile(path string, data []byte, flags int, perm fs.FileMode) error {
	f, err := os.OpenFile(path, flags, perm)
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
	if err := f.Close(); err != nil {
		// The bytes are complete and sealed, so the artifact is not a leak; a
		// failed Close on Windows can still mean the write did not flush.
		return err
	}
	return nil
}

// SealFile narrows an existing file to the current user.
func SealFile(path string) error { return sealFile(path) }

// PrivateDirAll creates path and any missing parents, sealing **every level it
// creates** before descending into it. Sealing afterwards is not enough: a
// child inherits the ACL it had at creation time, so a directory created wide
// and tightened later leaves everything written during that window wide. The
// order this guarantees is "sealed parent, then any child", which is what makes
// one seal cover files nobody seals individually - SQLite's -wal/-shm and a
// pre-rename retry temp.
//
// Ancestors above path are never touched: path is usually a directory the user
// chose (a data root under a profile), and sealing a profile directory on the
// way through would be a far larger change than the one being asked for. An
// existing path *is* sealed, because that is the repair path for a tree that
// predates this package.
func PrivateDirAll(path string, perm fs.FileMode) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("winsec: resolve %s: %w", path, err)
	}
	// Climb to the first existing ancestor; everything from there down is
	// missing and gets created + sealed in one pass, parent first.
	var missing []string
	cur := abs
	for {
		info, err := os.Lstat(cur)
		if err == nil {
			if !info.IsDir() {
				return fmt.Errorf("winsec: %s is not a directory", cur)
			}
			// Existing level: stop climbing. Its parents are not ours to seal,
			// and if cur is the requested path the SealDir below repairs it.
			break
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			return fmt.Errorf("winsec: no existing ancestor above %s", path)
		}
		missing = append(missing, cur)
		cur = parent
	}
	for i := len(missing) - 1; i >= 0; i-- {
		if err := os.Mkdir(missing[i], perm); err != nil {
			if existing, sErr := os.Lstat(missing[i]); sErr == nil && existing.IsDir() {
				continue // raced with another creator; seal it below
			}
			return err
		}
		if err := SealDir(missing[i]); err != nil {
			return err
		}
	}
	return SealDir(abs)
}

// SealDir narrows a directory to the current user, with new children inheriting
// exactly that and nothing wider.
func SealDir(path string) error { return sealDir(path) }

// RemoveUnlinked deletes the entry at path without following it, so a link
// standing where private data was expected disappears while whatever lives
// behind it is untouched. It never recurses - that is ticket 79's reasoning for
// os.Remove over os.RemoveAll, and RemoveAll would be the larger hole.
//
// When the entry is a link the platform refuses to unlink, the error wraps
// ErrIsReparsePoint and names the path: a reclaim loop that cannot clear a
// subtree must surface it rather than continue with a quota it no longer knows.
func RemoveUnlinked(path string) error { return removeUnlinked(path) }

func wrapPath(path string, err error) error {
	return fmt.Errorf("%w: %s: %v", ErrNotSealable, path, err)
}
