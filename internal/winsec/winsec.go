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
//
// What this package deliberately does not decide. The placement checks behind
// those words are about a spelling: whether any component of the path handed to
// them is a link, whether the ancestor chain stays inside the tree the call
// names, and whether an installed resolver's answer still names that same tree
// (resolve.go). None of them asks whose tree it is. Given a foreign absolute
// path that the caller names directly, with no link anywhere in its ancestor
// chain, SealFile succeeds - on Windows that strips an explicit S-1-1-0 grant
// standing on the file, on POSIX it narrows the mode - and nothing in this
// package can tell that the tree was never the caller's to seal. That is a
// ruled boundary, not a gap (ticket 113 AC#6, answering R-108-2): which roots
// this process may write under belongs to the caller's data-root discipline
// (tickets 76/95), and folding that policy into the floor would turn the guard
// into a second argument about intent instead of the one check that cannot be
// argued out of. A hard link sits exactly on this line - it shares an inode
// under a clean spelling, so it is invisible for the same reason - and
// winsec_other.go's platformVerifyPlacement carries the matching statement at
// the function that has to enforce it.
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
//
// The path is resolved before it is opened, not after: writing to a name the
// filesystem reads as some other object would place private bytes in public
// (ticket 94, same defect class as the directory leg below).
func privateFile(path string, data []byte, flags int, perm fs.FileMode) error {
	resolved, err := resolveString(path)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(resolved, flags, perm)
	if err != nil {
		return err
	}
	if err := sealHandle(f); err != nil {
		_ = f.Close()
		_ = os.Remove(resolved)
		return sealError(resolved, err)
	}
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		_ = os.Remove(resolved)
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
func SealFile(path string) error {
	resolved, err := resolveString(path)
	if err != nil {
		return err
	}
	return sealError(resolved, sealFile(resolved))
}

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
//
// path is resolved through the C26 PathResolver before any level is created or
// sealed, and a path C26 refuses to traverse is refused here too
// (risk.ErrReparseDenied propagates). Both halves of that sentence are the
// point: which tree this creates and seals *is* the security decision, so it may
// not be taken from a lexical spelling - a junction mid-path, an 8.3 short name,
// a \\?\ prefix or a trailing dot would make "sealing A" modify B, and the wide
// tree would then be reported to the caller as private. See resolve.go.
func PrivateDirAll(path string, perm fs.FileMode) error {
	dir, err := ResolvePath(path)
	if err != nil {
		return err
	}
	return privateDirAll(dir, perm)
}

// privateDirAll is the sealing walk, stated over a ResolvedPath so that no path
// reaching os.Mkdir or SealDir here can have skipped C26: the only mint for that
// type is ResolvePath (see resolve.go), which is what keeps a future caller from
// quietly reintroducing filepath.Abs.
func privateDirAll(dir ResolvedPath, perm fs.FileMode) error {
	abs := dir.String()
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
			return fmt.Errorf("winsec: no existing ancestor above %s", abs)
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
		// Each level is sealed through the platform primitive directly: it is a
		// prefix of an already-resolved path, so re-resolving would only add
		// syscalls, and skipping the seal is the hole AC#5 refuses to leave.
		if err := sealError(missing[i], sealDir(missing[i])); err != nil {
			return err
		}
	}
	return sealError(abs, sealDir(abs))
}

// SealDir narrows a directory to the current user, with new children inheriting
// exactly that and nothing wider.
func SealDir(path string) error {
	resolved, err := resolveString(path)
	if err != nil {
		return err
	}
	return sealError(resolved, sealDir(resolved))
}

// RemoveUnlinked deletes the entry at path without following it, so a link
// standing where private data was expected disappears while whatever lives
// behind it is untouched. It never recurses - that is ticket 79's reasoning for
// os.Remove over os.RemoveAll, and RemoveAll would be the larger hole.
//
// When the entry is a link the platform refuses to unlink, the error wraps
// ErrIsReparsePoint and names the path: a reclaim loop that cannot clear a
// subtree must surface it rather than continue with a quota it no longer knows.
//
// This is the one entry point that deliberately does not resolve its argument
// through C26. Its subject *is* the link: resolving first would deny the call on
// the very reparse point it exists to unlink, so the entry would stay occupied
// and the quota would keep counting a tree nobody owns. The safety comes from the
// operation instead of from a normalized spelling - it opens with
// FILE_FLAG_OPEN_REPARSE_POINT and never recurses - plus the ancestor check
// below, which is ticket 103's AC#2: PROBE F measured that a spelling *reaching*
// the leaf through a junction made this call delete a file inside somebody else's
// tree and return nil, because the only guard was on the last component.
//
// Because no resolver runs here, the argument has to name its own tree, so a
// relative spelling is refused (ticket 108's per-shape reading of AC#2): it
// resolves against wherever the process happens to be standing, which is a tree
// no caller named, and "removed" reported for an unlink the platform answered
// with "no such file" is the same report-success-while-doing-nothing shape this
// repository has already booked twice. The reclaim callers all pass paths joined
// onto a resolved data root.
func RemoveUnlinked(path string) error {
	if !filepath.IsAbs(path) {
		return fmt.Errorf("%w: %s is relative, so no tree names it and this call cannot tell which one it would remove from",
			ErrUnresolvedPath, path)
	}
	if link := firstLinkAncestor(path); link != "" {
		return fmt.Errorf("%w %s: the spelling reaches it through the link at %s, and whatever lives behind that link is not this tree's data to delete",
			ErrIsReparsePoint, path, link)
	}
	return removeUnlinked(path)
}

// firstLinkAncestor returns the shortest ancestor of path that is a link, or ""
// when none is. path itself is deliberately excluded: unblocking a link standing
// where an artifact was expected is RemoveUnlinked's whole purpose, and refusing
// on the leaf would break the reclaim route ticket 79 built.
//
// The prefixes come from pathPieces below, which is the whole of ticket 108's
// AC#2: ticket 103 cut the input on filepath.Separator only, so on Windows the
// spelling "C:\data\link/sub/keep-me.txt" looked like ONE component to the guard
// (PROBE P2 measured it deleting a file inside somebody else's tree and
// returning nil, as did the all-forward-slash spelling). Separator handling is
// platform-correct rather than uniform, and pathPieces says why.
func firstLinkAncestor(path string) string {
	prefixes := pathPieces(path)
	if len(prefixes) == 0 {
		return ""
	}
	for _, prefix := range prefixes[:len(prefixes)-1] {
		if ancestorIsLink(prefix) {
			return prefix
		}
	}
	return ""
}

// pathPieces returns every prefix of path that names one more component than the
// last, longest last, each one an *exact substring of the input*. The caller
// treats the final element as the leaf and inspects the rest as ancestors.
//
// Two rules, both load-bearing:
//
//  1. Which characters are separators is the platform's own answer. On Windows
//     the OS reads '/' as a separator exactly like '\', so both cut - refusing to
//     split there is what let P2/P3 through. On POSIX only '/' cuts: a backslash
//     is an ordinary character in a file name, so folding it into a separator
//     would either refuse a real ancestor that is not a link (the reclaim route
//     broken for no reason) or, once the pieces are re-joined, check "a/b" while
//     the real ancestor is the directory named `a\b` - a link at `a\b` then goes
//     unseen and the unlink lands in somebody else's tree. That second direction
//     is the fail-open this repository already booked as A74(3), which is why
//     prefixes are substrings of the input and are never rebuilt by joining.
//  2. A volume (`C:`, and the leading separators of a UNC or extended-length
//     spelling) is not an ancestor: Lstat("C:") names whatever directory the
//     process happens to be standing in, which is somebody else's link to no
//     purpose of this call.
//
// It never normalizes: no Clean, no Abs, no case folding (D22 ban #2, and the
// same reasoning as ticket 103's comment above) - a doubled separator yields the
// same ancestor twice, and a "." component yields a prefix the OS resolves to
// the object the previous one already named, which is exactly what an Lstat is
// allowed to be told.
func pathPieces(path string) []string {
	nativeIsBackslash := os.PathSeparator == '\\'
	isSep := func(c byte) bool {
		return c == os.PathSeparator || (nativeIsBackslash && c == '/')
	}
	vol := filepath.VolumeName(path)
	var out []string
	i := len(vol)
	for i < len(path) && isSep(path[i]) {
		i++
	}
	for i < len(path) {
		start := i
		for i < len(path) && !isSep(path[i]) {
			i++
		}
		if i > start {
			out = append(out, path[:i])
		}
		for i < len(path) && isSep(path[i]) {
			i++
		}
	}
	return out
}

// pathComponents returns the component names of path, using the same
// platform-correct separator rule as pathPieces, so a walk that checks names for
// trailing '.' or ' ' cannot disagree with the walk that Lstats their prefixes.
func pathComponents(path string) []string {
	nativeIsBackslash := os.PathSeparator == '\\'
	isSep := func(c byte) bool {
		return c == os.PathSeparator || (nativeIsBackslash && c == '/')
	}
	var out []string
	i := len(filepath.VolumeName(path))
	for i < len(path) && isSep(path[i]) {
		i++
	}
	for i < len(path) {
		start := i
		for i < len(path) && !isSep(path[i]) {
			i++
		}
		if i > start {
			out = append(out, path[start:i])
		}
		for i < len(path) && isSep(path[i]) {
			i++
		}
	}
	return out
}

func wrapPath(path string, err error) error {
	return fmt.Errorf("%w: %s: %v", ErrNotSealable, path, err)
}

// sealError normalizes any sealing failure to ErrNotSealable at the API
// boundary. The platform implementations already wrap, but the promise
// "a refused seal is not a completed write" is a contract for callers, so it
// must not depend on whether some implementation remembered to label its error.
func sealError(path string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrNotSealable) {
		return err
	}
	return wrapPath(path, err)
}
