//go:build windows

// Tickets 342 + 348, leg 342-r1, form 乙 on Windows: the pinned chain the private write route now
// stands on. This is the first Windows-side handle pinning in this package - before it, the route's
// act was by name in both directions (the open resolved the spelling, and the refusal-path unlink
// resolved it again), which is what the two tickets booked.
//
// WHY NEW CODE AND NOT REUSE. The family form 乙 reuses on POSIX - pinPlacement, sameObjectPinned,
// pinBase, chmodPinned - lives behind //go:build !windows in winsec_other.go, and this file's
// sealHandle goes by name (winsec_windows.go). So on this platform 乙 means "write the chain", and
// the chain is the stdlib's own os.Root descent, rooted at the volume, taking each component with
// OBJ_DONT_REPARSE so the kernel answers "is this component a link" INSIDE the open of the next one,
// instead of in a probe that can go stale.
//
// WHAT THE CHAIN BUYS, and the trap that comes with it. os.Root pins ANCESTORS and does NOT pin the
// LEAF: rootOpenFileNolog hands the leaf to the same doInRoot that resolves an in-root relative
// symlink by FOLLOWING it (Go's own openat maps ELOOP/ENOTDIR back into errSymlink and restarts).
// That is measured, not inferred: 342-a1's W-ROOTSYM row reports `followed=true` for a relative FILE
// link inside the root, three runs. So a leaf identity re-check is part of this landing and not an
// option - see openPinnedLeaf, and the residual it does NOT claim.
//
// WHAT THIS DOES NOT BUY, stated because the next reader will be asked to trust it:
//   - the leaf leg's own residual from ticket 322 stays open on this platform and stays named:
//     preOpenLeafGuard still answers about a NAME before the open, and an answer about a name is not a
//     lock on it. What is new here is that after the open the route no longer takes the followed
//     object on faith - it re-checks the name against the handle and refuses;
//   - an ancestor that is swapped between two states the walk can resolve (a real directory and a
//     NON-parse point) is a rename of the object itself, and no descriptor-based answer can refuse
//     that: pinning fixes WHICH object the act stands on, which is exactly what the two tickets ask,
//     and nothing wider;
//   - sealHandle on this platform still applies the descriptor by name with no same-object re-check.
//     That is ticket 120's open Windows leg, named at winsec_windows.go, and this landing does not
//     close it or claim to.
package winsec

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/windows"
)

// leafPin is this platform's face of the ancestor pin: the pinned handle of the directory that
// carries the leaf name, plus that directory's spelling for error text. It is a Root rather than an
// *os.File because the act-time primitives here are Root.OpenRoot / Root.OpenFile / Root.Lstat /
// Root.Remove, all of which resolve their argument relative to the handle, never through a spelling.
type leafPin struct {
	root *os.Root
	dir  string
}

// pinnedParent walks the caller's resolved spelling from the root of its own volume, pinning every
// component and re-verifying each one against the name it was reached through, and hands back the
// pin for the directory carrying the leaf.
//
// Order at each level is pin, then verify - never verify, then pin - which is the rule
// pinPlacement states for POSIX and the reason that function exists (Lstat before the open is the
// check-then-act shape: the racer plants a link after the check and the open follows it). Pin-then-
// verify closes both directions of that interleaving: either the name carries a link now, or it
// carries a different object than the one just pinned, and each of those is refused by
// sameObjectPinnedWindows below.
//
// The parent descriptor lives until the route releases it (see winsec.go's deferred close), because
// the refusal-path unlink is relative to it. Closing it as soon as the leaf is open would leave
// ticket 348's removal with nothing to be relative to - the coupling 342-v1's answer (v) named.
func pinnedParent(resolved string) (leafPin, error) {
	comps := pathComponents(resolved)
	leaf := "."
	ancestors := comps
	if n := len(comps); n > 0 {
		leaf = comps[n-1]
		ancestors = comps[:n-1]
	}
	cur, err := os.OpenRoot(pinVolumeRoot(resolved))
	if err != nil {
		return leafPin{}, wrapPath(resolved, err)
	}
	for _, comp := range ancestors {
		next, err := cur.OpenRoot(comp)
		if err != nil {
			_ = cur.Close()
			return leafPin{}, pinRefusal(resolved, comp, err)
		}
		pinned, err := next.Stat(".")
		if err != nil {
			_ = next.Close()
			_ = cur.Close()
			return leafPin{}, wrapPath(resolved, err)
		}
		if err := sameObjectPinnedWindows(resolved, comp, cur, pinned); err != nil {
			_ = next.Close()
			_ = cur.Close()
			return leafPin{}, err
		}
		_ = cur.Close()
		cur = next
	}
	if leaf == "." && len(comps) == 0 {
		_ = cur.Close()
		return leafPin{}, wrapPath(resolved, errors.New("the spelling names no directory to pin"))
	}
	return leafPin{root: cur, dir: filepath.Dir(resolved)}, nil
}

// pinVolumeRoot is the root the chain starts from: the volume the spelling names, with a separator,
// because os.OpenRoot needs a directory and filepath.VolumeName answers "D:" for a drive. On this
// platform a relative spelling has already been refused by resolveString, so there is an answer here.
func pinVolumeRoot(path string) string {
	if v := filepath.VolumeName(path); v != "" {
		return v + string(os.PathSeparator)
	}
	return string(os.PathSeparator)
}

// sameObjectPinnedWindows is the identity gate, and it is the Windows face of the two halves POSIX's
// sameObjectPinned carries: the name inside parent has to carry a NON-LINK right now, and it has to
// carry the object that was just pinned. Both halves are needed and neither is decoration.
//
// The predicate is the platform's own - FILE_ATTRIBUTE_REPARSE_POINT, the same one ancestorIsLink,
// isReparsePoint and preOpenLeafGuard already judge with - because a junction reports NO ModeSymlink
// to os.Lstat on Windows (winsec_windows.go's own statement). os.SameFile compares volume serial plus
// file index, which is a fact about the object and not about a representation, so it is not fooled by
// an 8.3 alias or a case difference - ticket 115's ruling that these comparisons go by tree.
func sameObjectPinnedWindows(path, name string, parent *os.Root, pinned fs.FileInfo) error {
	now, err := parent.Lstat(name)
	if err != nil {
		return wrapPath(path, err)
	}
	if entryIsLink(now) {
		return placementRefusalWindows(path, name, "the name carries a link, so the pinned object is not what the spelling reaches")
	}
	if !os.SameFile(pinned, now) {
		return placementRefusalWindows(path, name, "the name carries a different object than the one that was pinned")
	}
	return nil
}

// placementRefusalWindows is placementRefusal's Windows twin for a refusal raised by the pinned walk.
// It carries both sentinels through preOpenRefusal, because "a refused write is not a completed write"
// (AC#5) and "this spelling does not name the tree it claims" are properties of this answer too, and
// the reason text is the half that names WHICH leg refused - kept in the same words POSIX's
// sameObjectPinned uses so "is this a link" cannot mean two different things inside one package.
func placementRefusalWindows(path, name, why string) error {
	return preOpenRefusal(path, "pinned object at "+name+" is not what the name carries: "+why)
}

// pinRefusal attributes a component the walk could not even open. On this platform that is usually
// the kernel's own answer to OBJ_DONT_REPARSE (ERROR_REPARSE_POINT_NOT_RESOLVED or STATUS_STOP_ON_
// REPARSE_POINT, which Go surfaces as ELOOP/EPERM), i.e. the act-time refusal this landing came for -
// so it keeps both sentinels rather than reporting a bare OS error.
func pinRefusal(path, at string, err error) error {
	return placementRefusalWindows(path, at, "the walk could not pin that component through the handle of its parent: "+underlyingErrorMessage(err))
}

// openPinnedLeaf opens the leaf name ON the pinned parent handle, then re-checks the name against the
// handle it just got. The re-check is the part this platform does not get for free: os.Root pins
// ancestors and follows an in-root relative link AT the leaf (W-ROOTSYM, three runs), so "opened
// through a pinned handle" alone would still let the open name an object the caller never named.
//
// Refusing after the open is not cosmetic. privateFile writes only after sealHandle returns nil, so a
// refusal here stops the bytes from landing in the followed object at all, and the removal this route
// then performs is fd-relative, so it cannot delete a foreign name either. What it does NOT undo is a
// truncation the open itself already carried out on the followed object: this package's Windows leaf
// leg keeps that residual, and winsec.go's statement says so instead of claiming the window closed.
//
// The error shape keeps the by-name open's: *os.PathError with Op "open" and the full spelling, so the
// leaf-leg attribution rule ticket 322 paid for ("open " and not "openat ") still means one thing here.
func openPinnedLeaf(p leafPin, base string, flags int, perm fs.FileMode) (*os.File, error) {
	full := filepath.Join(p.dir, base)
	f, err := p.root.OpenFile(base, flags, perm)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: full, Err: underlyingErrorOf(err)}
	}
	pinned, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, &os.PathError{Op: "open", Path: full, Err: underlyingErrorOf(err)}
	}
	// The name is read through the SAME pinned handle, with a non-following open (Root.Lstat), so this
	// question is about the entry the caller named inside the object the route stands on - not about a
	// spelling that can be re-resolved by whoever renamed something in the meantime.
	now, err := p.root.Lstat(base)
	if err != nil {
		_ = f.Close()
		return nil, &os.PathError{Op: "open", Path: full, Err: underlyingErrorOf(err)}
	}
	if entryIsLink(now) {
		_ = f.Close()
		return nil, placementRefusalWindows(full, base, "the name carries a link, so this handle is not the object the caller named")
	}
	if !os.SameFile(pinned, now) {
		_ = f.Close()
		return nil, placementRefusalWindows(full, base, "this handle is not the object the name carries now, so the bytes would not land in what the route pinned")
	}
	return f, nil
}

// unlinkPinnedLeaf removes the leaf name through the pinned parent handle: the fd-relative form of
// the removal ticket 348 booked. os.Root's Remove descends to the parent by handle and calls
// Deleteat, which opens the entry with FILE_OPEN_REPARSE_POINT - so a link standing at that name is
// unlinked, never followed, which is the behaviour tickets 118 and 120-r3 pin on the other platform.
// There is no spelling left to re-resolve, so an ancestor swap cannot move this unlink.
//
// The error is still discarded, deliberately, exactly as `_ = os.Remove(resolved)` did: the route
// returns the seal failure or the write failure, which is the more specific truth the caller has to
// see, and changing what privateFile returns would be a change to C26's entry contract that this
// landing is not authorised to make. Ticket 348's own table records the swallowing as an observed
// cost of that boundary; the one thing this changes is that the removal itself is now safe.
func unlinkPinnedLeaf(p leafPin, base string) {
	if p.root == nil {
		return
	}
	_ = p.root.Remove(base)
}

// close releases the pinned handle. The route calls it once, in a defer, after every unlink attempt.
func (p leafPin) close() error {
	if p.root == nil {
		return nil
	}
	return p.root.Close()
}

// entryIsLink is the platform's own answer to "is this a link", read off a FileInfo this package
// already knows how to get. It duplicates the predicate inside isReparsePoint rather than calling
// that one, because isReparsePoint takes a PATH and re-resolves it, which is the thing being removed.
func entryIsLink(info fs.FileInfo) bool {
	attr, ok := info.Sys().(*syscall.Win32FileAttributeData)
	return ok && attr.FileAttributes&windows.FILE_ATTRIBUTE_REPARSE_POINT != 0
}

// underlyingErrorOf unwraps the *os.PathError the Root returns so the caller sees one PathError naming
// the full spelling and not two nested ones.
func underlyingErrorOf(err error) error {
	var pe *os.PathError
	if errors.As(err, &pe) {
		return pe.Err
	}
	return err
}

func underlyingErrorMessage(err error) string {
	if u := underlyingErrorOf(err); u != nil {
		return u.Error()
	}
	return "no error reported"
}
