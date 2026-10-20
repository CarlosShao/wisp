//go:build unix

// Tickets 342 + 348, leg 342-r1, form 乙 on the unix family: the ancestor pin and the two
// act-time primitives the private write route stands on.
//
// WHY THIS FILE IS TAGGED `unix` AND NOT `!windows`. The primitives below are openat and
// unlinkat. `syscall.Openat` exists only on linux (this package's own r1 probe, recorded in
// winsec_other.go's costs: undefined on darwin, freebsd, openbsd, netbsd, dragonfly and
// solaris), so the family answer is golang.org/x/sys/unix, which carries them for linux,
// android, darwin, the *bsd set, solaris/illumos, aix and zos - the set where this package
// compiles at all today. js/wasm, plan9 and wasip1 have no openat in either package (the
// stdlib itself falls back to a by-name implementation there: os.root_noopenat.go), so they get
// the named refusal in the sibling file rather than silence. This does not widen the compile
// matrix: winsec_other.go's own syscall.O_NOFOLLOW and syscall.O_NONBLOCK already exclude js,
// plan9 and wasip1 at this anchor (measured with `go build` per GOOS before this file landed).
//
// WHY THE PARENT IS PINNED INSTEAD OF THE CHAIN BEING RE-RESOLVED BY NAME. Both tickets are
// one sentence: resolveString and platformVerifyPlacement answer about a SPELLING, and every act
// after them used to re-resolve that spelling through the kernel. Whoever can rename an ancestor
// in the gap therefore chose which object this open truncated (ticket 342) and which object the
// refusal-path unlink removed (ticket 348). A descriptor opened before the gap cannot be
// re-pointed by it, because openat and unlinkat resolve only the name they are handed, RELATIVE
// to the object the descriptor names. That is the answer this platform can actually give, and it
// is the same answer ticket 120-r2 gave for the seal: reuse the package's own handle pinning
// (pinPlacement) instead of inventing a second walk. Form 乙 therefore adds no new walker, no new
// identity gate, and no new exported name anywhere.
//
// WHAT THIS DOES NOT BUY, stated because the next reader is being asked to trust it:
//
//   - the leaf leg still needs O_NOFOLLOW and it still gets it: the route hands the flag in as
//     `flags|oNoFollowPlatform` exactly as it did on the by-name os.OpenFile, so a link planted at
//     the LEAF name answers ELOOP inside the open instead of being truncated through (ticket
//     322's POSIX half, unchanged by this file);
//   - the pin does not freeze the parent directory. Whoever can swap the ancestor can still rename
//     the leaf. What is closed is object SELECTION: open, truncate, write and unlink all land on a
//     child of the pinned directory, which is the object the caller's spelling named before the
//     gap;
//   - one descriptor is not a lock. Nothing here stops the pinned directory itself being deleted.
//     It stops the open and the unlink being RE-DIRECTED by a rename, which is what the two
//     tickets measured, on both damage shapes.
package winsec

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"syscall"

	"golang.org/x/sys/unix"
)

// leafPin is the unix face of the ancestor pin: the descriptor pinPlacement handed back for the
// directory that carries the leaf name, plus the spelling of that directory for error text.
type leafPin struct {
	dir  *os.File
	name string
}

// pinnedParent is the route's first act-time step: pin the directory that will carry the leaf.
//
// pinPlacement is the function the POSIX seal already stands on (winsec_other.go), called here
// under the contract it already has there: every component is pinned BEFORE it is verified, the
// verification is sameObjectPinned (a name carrying a link, or a name carrying a different object
// than the one that was pinned, both refuse), and what comes back is a descriptor for that object
// opened O_RDONLY|O_NONBLOCK. Reusing it is the whole cost claim of form 乙 over form 甲 on this
// platform, and it means a refusal here carries the same two sentinels and the same reason text
// tickets 113 and 120 already pin.
//
// filepath.Dir is not a normalization (it cuts at the last separator; D22's ban is on Clean/Abs
// making decisions, and this package already calls Dir in privateDirAll). The name handed to the
// open and to the unlink is filepath.Base of the same spelling, so the two are one component by
// construction - the same decomposition pinPlacement applies at its own leaf level.
func pinnedParent(resolved string) (leafPin, error) {
	dir := filepath.Dir(resolved)
	pin, err := pinPlacement(dir)
	if err != nil {
		return leafPin{}, err
	}
	return leafPin{dir: pin, name: dir}, nil
}

// close releases the pin. On this platform the pin is a descriptor and closing it is the whole
// release; the route calls this once, in a defer, after every unlink attempt.
func (p leafPin) close() error {
	if p.dir == nil {
		return nil
	}
	return p.dir.Close()
}

// openPinnedLeaf opens the leaf name on the pinned parent descriptor.
//
// ★Ticket 353's third return value is answered with nil HERE, on purpose. The route is untagged
// (winsec.go), so the Windows landing's identity re-check forces this signature onto every build set;
// what it must NOT force onto this platform is behaviour. This file adds no extra syscall and no new
// failure branch, and unlinkPinnedLeaf below keeps the fd-relative unlinkat exactly as ticket 348
// landed it, because the damage 353 measured is a Windows-shape damage: it exists because this
// package's Windows leaf handle carries FILE_SHARE_DELETE
// (Go's at_windows.go passes that mask for a Root open; the os package's by-name open does not), which
// is what lets the leaf name be vacated and
// re-filled while the route holds it. POSIX has no such sharing mask, and this platform's route
// already answers the identity question twice before the removal - in pinPlacement's walk
// (sameObjectPinned) and in sealHandle's inline faces below - so the parameter is taken and ignored,
// and said out loud rather than silently.
//
// unix.O_CLOEXEC is not decoration: os.OpenFile marks the descriptor close-on-exec inside the Go
// runtime and a raw openat does not inherit that behaviour. Without it every private write would
// hand the leaf descriptor to whatever child process this process spawns, and the leak would
// outlive the File - a wider hole than the one being closed here.
//
// The result is a real *os.File over the kernel's own descriptor, so sealHandle's fchmod path and
// the caller's Write both see the object the pinned directory named at the instant of the open,
// not whatever the spelling resolves to afterwards.
//
// The error shape is part of the contract, not cosmetics: a refusal comes back as *os.PathError
// with Op "open" naming the full spelling, which is exactly what the by-name os.OpenFile produced.
// openleaf_322_other_test.go tells ticket 322's leaf leg apart from pinPlacement's chain ELOOP by
// precisely that difference ("open " and not "openat "), so renaming the operator here would
// silence a pin that ticket 322 paid for.
func openPinnedLeaf(p leafPin, base string, flags int, perm fs.FileMode) (*os.File, fs.FileInfo, error) {
	path := filepath.Join(p.name, base)
	fd, err := unix.Openat(int(p.fd()), base, flags|unix.O_CLOEXEC, openModeBits(perm))
	if err != nil {
		return nil, nil, &os.PathError{Op: "open", Path: path, Err: err}
	}
	f := os.NewFile(uintptr(fd), path)
	if f == nil {
		// Unreachable in practice: Openat returns a negative descriptor only together
		// with an error, which the branch above already took. EBADF is named because it
		// is the fact os.NewFile is reporting, and the descriptor is released here
		// rather than leaked.
		_ = unix.Close(fd)
		return nil, nil, &os.PathError{Op: "open", Path: path, Err: syscall.EBADF}
	}
	return f, nil, nil
}

// unlinkPinnedLeaf removes the leaf name through the pinned parent: unlinkat with flag 0, the
// fd-relative form of the removal ticket 348 measured by name.
//
// What this closes is the SECOND damage shape, and it is not a race: the route unlinks on a
// refusal, and `os.Remove(resolved)` re-resolved the whole spelling, so an ancestor standing on a
// link pointed the unlink at a name in somebody else's tree. 342-a1's A3 case measured it
// succeeding (unlink err=<nil>, witness gone) - deterministically reachable, not a window anyone
// had to win.
//
// unlinkat resolves only `base`, inside the object `p` names. There is no spelling left to
// re-resolve, so the ancestor is not in the path at all. The leaf case keeps the behaviour tickets
// 118 and 120-r3 pin: unlinkat on a name carrying a symlink removes the LINK and never its target,
// and that is the kernel answering, not a check that can go stale.
//
// The error is still discarded, deliberately - that is the shape `_ = os.Remove(resolved)` had, and
// privateFile's returned error is the seal failure or the write failure, which is the more specific
// truth the caller has to see. Ticket 348's own现量 records the swallowing as an observed cost;
// changing what this route returns would be a change to C26's entry contract, which is not this
// landing's to make. The one thing this file does is make the removal itself safe.
//
// ★The `pinned` parameter is ticket 353's Windows identity reading, taken here and NOT used, for the
// reason stated on openPinnedLeaf above: the route is untagged, so the signature travels to every
// build set while the damage it guards was measured on the platform whose leaf handle shares DELETE.
// Nothing on this branch changed: same unlinkat, same flag 0, same swallowing.
func unlinkPinnedLeaf(p leafPin, base string, pinned fs.FileInfo) {
	_ = unix.Unlinkat(int(p.fd()), base, 0)
	runtime.KeepAlive(p.dir)
}

// fd hands out the raw parent descriptor and keeps the *os.File alive across the call, so a pin
// that is no longer otherwise reachable cannot be finalized (and its descriptor closed) under an
// in-flight openat.
func (p leafPin) fd() uintptr {
	defer runtime.KeepAlive(p.dir)
	return p.dir.Fd()
}

// openModeBits maps an fs.FileMode onto the mode argument of openat the way the os package maps it
// for OpenFile: permission bits plus the three setuid/setgid/sticky bits when the caller asks for
// them. A private artifact never should, but the by-name open this replaces did honour them, and
// this landing is not allowed to quietly change what a caller's perm argument means.
func openModeBits(perm fs.FileMode) uint32 {
	mode := uint32(perm.Perm())
	if perm&fs.ModeSetuid != 0 {
		mode |= syscall.S_ISUID
	}
	if perm&fs.ModeSetgid != 0 {
		mode |= syscall.S_ISGID
	}
	if perm&fs.ModeSticky != 0 {
		mode |= syscall.S_ISVTX
	}
	return mode
}
