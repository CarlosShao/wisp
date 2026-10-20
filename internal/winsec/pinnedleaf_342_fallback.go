//go:build !windows && !unix

// Tickets 342 + 348, leg 342-r1, the fallback for build sets with no openat/unlinkat family:
// js/wasm, plan9 and wasip1.
//
// This file is DEAD CODE AT THIS ANCHOR, and that is the statement rather than a claim of
// coverage: winsec_other.go cannot compile for any of those sets either - it uses syscall.O_NOFOLLOW
// (undefined on js and plan9) and syscall.O_NONBLOCK (undefined on wasip1), measured with `go build`
// per GOOS before this landing, and the r1 probe recorded in winsec_other.go's costs says the same.
// The compile matrix the package ships today is unix plus windows.
//
// It exists because the alternative is worse. Form 乙's claim is that the act is RELATIVE TO A
// PINNED HANDLE; a platform with no openat has no such act, and the tempting shape is to fall back to
// the by-name open and by-name unlink this landing deletes. That would put the truncation window and
// the follow-the-ancestor unlink back on that platform behind a file that reads like it implements
// them. So the fallback refuses, and it refuses in this package's established direction: never wider,
// never silently by name.
//
// If the package is ever widened to one of these platforms, the work is to implement these three
// primitives with that platform's own act-time answer - and the per-platform pins tickets 342 AC#1
// and 348 AC#2 demand have to be measured there, not assumed from the unix rows.
package winsec

import (
	"errors"
	"io/fs"
	"os"
)

// leafPin is the fallback face of the ancestor pin: always an unusable zero value, because there is
// no descriptor to hold on this build set.
type leafPin struct{}

// errNoActTimePrimitive names the refusal below. It is a package-private sentinel inside a build set
// that does not compile today - not an exported API, and not reachable from any shipped binary.
var errNoActTimePrimitive = errors.New("winsec: this platform has no openat/unlinkat, so a private write cannot be pinned to a parent handle")

// pinnedParent refuses rather than pretending to pin.
func pinnedParent(resolved string) (leafPin, error) {
	return leafPin{}, wrapPath(resolved, errNoActTimePrimitive)
}

// close is a no-op: there is never a handle on this build set.
func (leafPin) close() error { return nil }

// openPinnedLeaf refuses: with no pinned descriptor there is no leaf open that is not a by-name one.
//
// ★The signatures of these two primitives are the route's (untagged winsec.go), so ticket 353's
// Windows identity reading and re-check reach every build set. The extra reading is answered with nil
// here because nothing is opened to describe: this build set refuses before any act, which is the
// established direction (never wider, never silently by name), and the removal below still does
// nothing. No behaviour on this build set changed.
func openPinnedLeaf(leafPin, string, int, fs.FileMode) (*os.File, fs.FileInfo, error) {
	return nil, nil, &os.PathError{Op: "open", Path: "unpinned", Err: errNoActTimePrimitive}
}

// unlinkPinnedLeaf does nothing: the by-name removal it would have to stand in for is the damage
// ticket 348 booked, and doing nothing keeps the refused entry standing instead.
func unlinkPinnedLeaf(leafPin, string, fs.FileInfo) {}
