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

// platformVerifyPlacement is the POSIX half of the floor the Windows file
// implements, and ticket 113 is the reason it is no longer `return path, nil`.
//
// That stub used to justify itself by pointing at internal/risk's DEFERRED
// non-Windows resolver: "the platform's own PathResolver detects no reparse
// traversal, so the floor must not invent a stricter rule than the resolver
// it stands under." The comparison was with the wrong side. This function is the
// floor for binaries that link no resolver at all (internal/secret and
// internal/memory do not), and os.Chmod on this platform follows symlinks, so
// `return path, nil` did not mean "this platform has nothing to check" - it meant
// the sealing entry points had no link leg here, and PROBE P3's outcome
// reappeared one platform away: SealFile through a symlink returned nil while the
// chmod landed on somebody else's file (measured: -rw-rw-rw- becomes
// -rw-------). A guarantee that holds only where the second platform's API
// happens to be stricter is not a guarantee; "Windows 上守住了" was never evidence
// that this invariant holds.
//
// So the walk is the same walk as placement_windows.go's, on this platform's own
// facts: every prefix pathPieces offers is Lstat'ed, and the predicate is
// ancestorIsLink - the same one RemoveUnlinked uses above, so the reclaim route
// and the seal route cannot disagree about what a link is. pathPieces is reused
// rather than re-cut here because it is the piece ticket 108's AC#2 proved
// platform-correct: on POSIX it cuts on '/' alone, and a backslash is an ordinary
// character in a file name. Treating the backslash as a separator on this
// platform is not conservatism - it either refuses the real directory named `a\b`
// or, once the pieces are re-joined, Lstats `x/y` while the link is `x\y`, which
// is the cross-directory fail-open this repository already booked as A74(3).
//
// Like the Windows leg it can only refuse, never rewrite (D22 ban #2), and it
// checks the leaf as well as the ancestors: a symlink standing where the caller
// named a file is not "the tree this call names" either, and it is precisely the
// case where os.Chmod would seal the target. Missing prefixes read as "not a
// link", which is what lets PrivateDirAll create its own levels parent first.
//
// Two costs, both stated rather than discovered later:
//
//   - a data root that reaches itself through a symlink starts refusing. That is
//     the already-booked R-103-7 trade for RemoveUnlinked, now applied to
//     sealing too, and the reason is the same: refusing loudly is the only
//     direction available to a floor. The shapes are not hypothetical, and
//     ticket 113 registered only the first of them (R-113-B, measured in a Linux
//     container: 81 failing lines across this package and internal/config,
//     internal/agent, internal/memory with TMPDIR behind a symlink). Three
//     shapes, and only the first was registered by ticket 113: macOS, where /tmp
//     and /var are symlinks so any root under them is one; the test env's data
//     root, which is os.TempDir() plus a pid suffix; and Linux, where dotfiles
//     managers commonly make $HOME/.config a symlink, which is the user config
//     dir this repository's data root hangs off.
//
//     Ticket 119's answer is option 2 of the three it was offered: the layer that
//     asks the OS resolves what the OS answered. This function's rule is
//     unchanged, deliberately: a per-platform containment test inside the floor
//     would make the POSIX leg weaker than the Windows one, which is the exact
//     asymmetry ticket 113 was opened to close. What it still refuses, and must
//     keep refusing, is a root nobody resolved - the link is then the only thing
//     that says which tree the bytes land in, and this package does not read
//     intent out of a spelling (see the package doc, R-108-2).
//
//     Which roots are resolved today is a list of this repository's callers, not
//     a property this package can check: proc.TestDataDir (os.TempDir),
//     proc.DefaultLayout (os.UserConfigDir), cmd/wisp's resolveDataDir and - since
//     ticket 119's rework, which is what R-119-1 was - cmd/wisp's
//     resolveSecretLayout. The fourth was missing, and measured on one and the
//     same binary: `wisp doctor` printed a resolved tree while `wisp secret list`
//     returned rc=1 through a symlinked HOME, .config or XDG_CONFIG_HOME, because
//     that command reads os.UserConfigDir() itself and never passes through
//     DefaultLayout. So the honest claim is "every data root this repository ships
//     today is resolved before it reaches this floor", never "a root handed to
//     this floor names a real tree": the next sealing site has to do the same and
//     nothing here notices if it forgets. The rule that decides which root gets
//     resolved is R-119-3's ruling, written on proc.TestDataDir - an answer to an
//     OS question is resolved, a value this process was handed by name is left
//     exactly as declared.
//
//   - what option 2 buys back is a decision this leg used to force. Once the
//     caller resolves, whoever owns TMPDIR / HOME / XDG_CONFIG_HOME decides which
//     tree a seal lands in, where the answer used to be a refusal (that is
//     R-119-4's account, and it is a landing-point ownership question rather than
//     a traversal one: a link planted inside an already-resolved data root is
//     still refused here, pinned on that route by ticket 119's
//     TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119 and, for the
//     resolveSecretLayout route, by
//     TestAC3POSIXSecretRouteLinkInsideItsDataRootStillRefused119).
//
// Hard links stay out of this leg and were never the second cost: a regular file
// that happens to share its inode with somebody else's name has no symlink
// anywhere in its spelling, so nothing in this package can see it; that is
// R-108-2's question ("whose tree is this"), which lives with the caller's
// data-root discipline (tickets 76/95), not in a placement check.
//
// Deliberately out of scope, and named so nobody reads this function as the end
// of the family: R-108-3. Once a resolver is installed, the use-time leg trusts
// that resolver's own `rewritten` account plus this floor re-run on its answer,
// so a resolver that answers each tree with a clean, mutually contained spelling
// and truthfully reports no rewrite can still move a tree. That leg is a property
// of the seam in resolve.go, it is reachable only from inside this package today
// (no release path, latch proven race-free by ticket 108's AC#1), and closing it
// here would be a second answer to one question.
func platformVerifyPlacement(path string) (string, error) {
	for _, prefix := range pathPieces(path) {
		if ancestorIsLink(prefix) {
			return "", fmt.Errorf("%w: %s reaches it through the link at %s, which is not the tree this call names",
				ErrUnresolvedPath, path, prefix)
		}
	}
	return path, nil
}

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

// ancestorIsLink is RemoveUnlinked's half of ticket 103's AC#2, and since ticket
// 113 the predicate platformVerifyPlacement walks with. The two routes need it for
// mirror-image reasons: os.Remove never follows the leaf, but a symlink in the
// middle of a spelling still redirects the unlink into somebody else's tree, and
// os.Chmod follows the leaf as well, so the seal route has to refuse one more
// component than the unlink route does.
func ancestorIsLink(prefix string) bool {
	info, err := os.Lstat(prefix)
	if err != nil {
		return false // an ancestor that is not there cannot be a link
	}
	return info.Mode()&os.ModeSymlink != 0
}
