//go:build !windows

package winsec

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
)

// oNoFollowPlatform is the leaf-link leg of ticket 322, carried on the open inside
// privateFile (winsec.go). It is a constant rather than a bare syscall.O_NOFOLLOW
// in that shared file because the flag has no Windows spelling: the compile matrix
// in .scratch/wisp/probes/322/a1 measured syscall.O_NOFOLLOW defined across the
// whole unix family and undefined on windows, js and plan9, so the portable file
// has to split it and the Windows file defines it as 0 and answers the same
// question a different way (see winsec_windows.go's preOpenLeafGuard).
//
// What it buys, measured rather than assumed: by-name, every leaf-link spelling
// (relative, absolute, climbing out of the root, absolute with an out-of-root
// target, dangling) answers ELOOP "too many levels of symbolic links" and the
// object behind the link keeps its bytes, its size and its mtime bit for bit.
// What it does not buy, and must never be written as bought: it speaks about the
// LAST component only, so an ancestor planted as a link in the same gap is still
// followed. That leg is open at this anchor.
const oNoFollowPlatform = syscall.O_NOFOLLOW

// preOpenLeafGuard is the POSIX half of the pre-open hook, and it does nothing.
//
// That is not a stub standing in for an unmeasured platform: on this platform the
// question it would ask is already answered INSIDE the open, atomically with it, by
// oNoFollowPlatform above. A probe run before the open would only add a verdict
// that can go stale between the probe and the open - which is the exact shape this
// package exists to delete (see applyDescriptorPOSIX's costs) - so the honest POSIX
// answer is that there is nothing to check by name here.
func preOpenLeafGuard(string) error { return nil }

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

// applyDescriptorPOSIX narrows the mode of the object a spelling names, checked
// and applied on ONE pinned descriptor.
//
// This function used to be os.Chmod(path) followed by os.Stat(path): one check
// and one act on a *name*, which is two resolutions of the same spelling with a
// window in the middle, and os.Chmod follows any link planted in that window.
// Ticket 113's acceptance round measured 154 of 4000 retries landing the seal on
// somebody else's file through it (R-113-D), and ticket 120's r1 round re-measured
// the same window on the committed code with the instrument in
// posix_seal_window_120_other_test.go. The act below is fchmod on the descriptor
// the walk pinned, so the object that passed the identity check is the object that
// is narrowed - check and act are the same object by construction, not by luck
// about how fast the racer is.
//
// The walk is portable stdlib only, and it is NOT the form ticket 120's AC#2 named
// first. The fully openat-based form is unreachable from this file: syscall.Openat
// exists only on linux (r1's probe: undefined on darwin, freebsd, openbsd, netbsd,
// dragonfly and solaris), and os.Root does not enforce this package's floor rule -
// it follows a symbolic link whose target stays inside the root, which for a root
// of "/" means every relative link with no ".." element is followed (measured this
// round: OpenFile on "sib.bin -> target.bin" returned a descriptor of the TARGET).
// So O_NOFOLLOW through Root is not a defence here and must never be written as
// one (r1 measured it defeated by Root.OpenFile), and "Root refuses it" only covers
// absolute targets and ".."-climbing targets, which is not the interesting shape.
// What closes the window is sameObjectPinned: the pinned descriptor's identity
// (fstat) compared against what the name carries right now (lstat at the parent).
//
// Three costs, stated here rather than discovered later.
//
//   - an open is now a permission question. Pinning takes O_RDONLY on the leaf and
//     read plus search on every ancestor, where os.Chmod needed only search. So a
//     file at 0w0 or 000, or a directory without r, goes from "the owner can still
//     repair its mode" to "refused". This is a tightening, and it was accepted by
//     the ledger as A990 ruling 2: the direction is refusal, never relaxation (this
//     repository's iron rule is that nothing is loosened to make a reading green),
//     the leg is !windows-only so the shipped platform is untouched, and the cost
//     was measured at uid 1000 in a container (0200/0000 files and a 0300 directory
//     answer EACCES; under the container's uid 0 all of them answer ok, which is
//     ticket 113's standing note about root distorting permission faces - so the
//     permission faces of this leg are read at uid 1000, not at uid 0).
//   - the open also has to be non-blocking. A FIFO standing where an artifact was
//     named would hang an O_RDONLY open forever (measured: blocked past 1.5s, while
//     the same open with O_NONBLOCK returned in 115 microseconds and fchmod on it
//     worked), so the leaf open carries O_NONBLOCK. os.Chmod never opened anything,
//     so this is a new failure mode the old shape could not have; hanging a seal is
//     worse than refusing one.
//   - one seal is now roughly two syscalls per path component plus the fchmod and
//     its read-back, instead of chmod plus stat. Sealing is not a hot path in this
//     repository (data roots at startup, artifacts per spill), and the alternative
//     is the window above.
//
// And one thing this does NOT fix, named so nobody reads the paragraph above as the
// end of the family: the OPEN inside privateFile (winsec.go) resolves a leaf link
// before this walk runs. Ticket 322 closed the LEAF leg of that open on this
// platform by carrying O_NOFOLLOW on it (oNoFollowPlatform above), so a leaf link
// planted in the gap now makes the open itself answer ELOOP instead of truncating
// the object behind the link. What is still open, and stays open at this anchor,
// is the ANCESTOR leg: O_NOFOLLOW speaks about the last component only, so an
// ancestor swapped for a link after resolveString returned is still followed by
// that open, and its O_TRUNC still lands on a file in somebody else's tree. The
// seal legs below stop the private bytes from arriving there; they do not stop the
// truncation. Nobody may read either sentence as the window being closed.
func applyDescriptorPOSIX(path string, dir bool) error {
	mode := fs.FileMode(0o600)
	if dir {
		mode = 0o700
	}
	target, err := pinPlacement(path)
	if err != nil {
		return err
	}
	defer target.Close()
	return chmodPinned(target, path, mode)
}

// pinPlacement walks path from the root of its own volume, pinning every component
// as a descriptor and re-verifying each one against the name it was reached through,
// and hands back the pinned descriptor of the object path names. Closing every
// parent on the way down is safe and deliberate: a child descriptor is independent
// of the directory it was opened from, so the chain that the act stands on is the
// sequence of objects, not the spelling.
//
// Order at each level is pin, then verify - never verify, then pin. Lstat before the
// open is the check-then-act shape this function exists to delete: the racer plants
// a link after the Lstat and the open follows it. Pin-then-verify closes both
// directions of that interleaving. Either way the object that comes back is refused
// unless (a) the name still carries a non-link and (b) the name's object is the
// pinned object, by dev and inode (os.SameFile). A link planted and then swapped back
// for a *different* real object still fails (b); a link still standing fails (a).
//
// Nothing after the last check re-resolves this spelling by name. That is the whole
// claim of the function: the act is fchmod on what is returned here, so a swap that
// happens after the walk cannot move it (measured: with the ancestor renamed away and
// replaced by a link, the seal still landed on the original object and the foreign
// witness kept its mode). The only by-name read left below is fstypeOf, which runs on
// an already-refused path to name a filesystem type in an error message; it cannot
// select an object because nothing is acted on after it.
func pinPlacement(path string) (*os.File, error) {
	if !filepath.IsAbs(path) {
		// Not a new refusal in practice: every caller reaches this through
		// builtinVerifier.Resolve, which refuses a relative spelling before any
		// seal runs. It is here because pinning has to name the root it starts
		// from, and a relative spelling has no such root - it has wherever the
		// process happens to be standing, which is a tree this call did not name.
		return nil, placementRefusal(path, path, "the spelling is not absolute, so no root pins it")
	}
	comps := pathComponents(path)
	leaf := "."
	ancestors := comps
	if n := len(comps); n > 0 {
		leaf = comps[n-1]
		ancestors = comps[:n-1]
	}
	cur, err := os.OpenRoot(pinBase(path))
	if err != nil {
		return nil, wrapPath(path, err)
	}
	for _, comp := range ancestors {
		next, err := cur.OpenRoot(comp)
		if err != nil {
			_ = cur.Close()
			return nil, wrapPath(path, err)
		}
		pinned, err := next.Stat(".")
		if err != nil {
			_ = next.Close()
			_ = cur.Close()
			return nil, wrapPath(path, err)
		}
		if err := sameObjectPinned(path, comp, cur, pinned); err != nil {
			_ = next.Close()
			_ = cur.Close()
			return nil, err
		}
		_ = cur.Close()
		cur = next
	}
	f, err := cur.OpenFile(leaf, os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		_ = cur.Close()
		return nil, wrapPath(path, err)
	}
	info, err := f.Stat()
	if err != nil {
		_ = f.Close()
		_ = cur.Close()
		return nil, wrapPath(path, err)
	}
	if err := sameObjectPinned(path, leaf, cur, info); err != nil {
		_ = f.Close()
		_ = cur.Close()
		return nil, err
	}
	_ = cur.Close()
	return f, nil
}

// sameObjectPinned is the identity gate: the name inside parent has to carry a
// non-link right now, and it has to carry the object that was just pinned. Both
// halves are needed and neither is decoration - the first refuses the link that
// Root follows happily, the second refuses the swap-back that leaves no link
// standing by the time anyone looks.
func sameObjectPinned(path, name string, parent *os.Root, pinned fs.FileInfo) error {
	now, err := parent.Lstat(name)
	if err != nil {
		return wrapPath(path, err)
	}
	if now.Mode()&os.ModeSymlink != 0 {
		return placementRefusal(path, name, "the name carries a link, so the pinned object is not what the spelling reaches")
	}
	if !os.SameFile(pinned, now) {
		return placementRefusal(path, name, "the name carries a different object than the one that was pinned")
	}
	return nil
}

// pinBase is the root the pinning walk starts from. On this platform
// filepath.VolumeName is the empty string for every spelling, so this is "/", and
// the statement is spelled out rather than assumed because the walk's whole meaning
// is "every component below this root is checked as an object".
func pinBase(path string) string {
	if v := filepath.VolumeName(path); v != "" {
		return v
	}
	return string(os.PathSeparator)
}

// placementRefusal is the act-time face of the rule platformVerifyPlacement states
// before the act, and it carries both sentinels on purpose: ErrUnresolvedPath keeps
// a refusal attributable to "this spelling does not name the tree it claims" - which
// is what the ticket 120 instrument counts as refused-at-check - and ErrNotSealable
// keeps this package's contract ("a refused seal is not a completed write", AC#5)
// independent of whether sealError remembered to relabel it.
func placementRefusal(path, at, why string) error {
	return fmt.Errorf("%w: %s: pinned object at %s is not what the name carries: %s (%w)",
		ErrUnresolvedPath, path, at, why, ErrNotSealable)
}

// chmodPinned narrows the mode through an already pinned descriptor and reads it back
// from that SAME descriptor. The read-back is still AC#6's requirement, not a
// formality: a filesystem that stores no usable mode bits answers fchmod with success
// and reports something else, and that has to become ErrNotSealable rather than a
// promise.
func chmodPinned(f *os.File, path string, mode fs.FileMode) error {
	if err := f.Chmod(mode); err != nil {
		return wrapPath(path, err)
	}
	info, err := f.Stat()
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
// PrivateFileExclusive. It consults applyDescriptor first, so the seam AC#5's
// injection replaces is on this route on both platforms (private_fail_test.go is
// untagged and asserts that a refused handle seal stops the bytes), and then adds
// the question a seal by name cannot ask: is the descriptor the caller is about to
// write through the object that was just narrowed?
//
// The identity check is os.SameFile(f.Stat(), os.Lstat(f.Name())) - the caller's
// open and the name, not two opens of the name. Lstat does not follow the leaf, so
// this leg needs no walk, and it refuses both shapes a by-name seal cannot see: the
// name carrying a link, and the name having been moved to a different real object
// since the caller opened it. In the second case the pinned walk seals the object
// the name denotes now, which is NOT the object the pending write goes to, so
// without this leg private bytes would land in a file that was never narrowed.
//
// What it cannot do is undo the damage of the open that came before it; see the
// ticket 322 sentence in applyDescriptorPOSIX's costs.
func sealHandle(f *os.File) error {
	if err := applyDescriptor(f.Name(), false); err != nil {
		return err
	}
	name := f.Name()
	pinned, err := f.Stat()
	if err != nil {
		return wrapPath(name, err)
	}
	now, err := os.Lstat(name)
	if err != nil {
		return wrapPath(name, err)
	}
	if now.Mode()&os.ModeSymlink != 0 {
		return placementRefusal(name, name, "the name carries a link, so this descriptor is not the object the caller named")
	}
	if !os.SameFile(pinned, now) {
		return placementRefusal(name, name, "this descriptor is not the object the name carries now, so the bytes would not land in what was just sealed")
	}
	return chmodPinned(f, name, 0o600)
}

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
