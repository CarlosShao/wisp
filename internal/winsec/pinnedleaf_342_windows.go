//go:build windows

// Tickets 342 + 348, leg 342-r1 — TRANSITIONAL FILE, and it says so on purpose.
//
// This is the POSIX-half commit (cap order in the dispatch: ① POSIX code + resident pins + the
// container reading, ② the Windows pinned chain, ③ the self-statement and the ticket tails). The
// portable route in winsec.go now calls three platform primitives, so this platform has to answer
// them in the same commit or the shipped platform stops compiling. What it answers HERE is the
// by-name shape the route used before this landing, byte for byte:
//
//   - the open is os.OpenFile(resolved, flags|oNoFollowPlatform, perm) - ticket 322's leaf leg, via
//     preOpenLeafGuard plus the zero-valued platform constant, exactly as shipped;
//   - the two refusal-path removals are os.Remove(resolved), exactly as shipped.
//
// So no Windows behaviour changes in this commit, no Windows pin changes colour, and nothing here
// claims a window is closed. The ANCESTOR leg stays open on this platform and the comment in
// winsec.go says so; the pinned-handle chain that closes it (os.Root descent plus the leaf identity
// re-verification this platform has never had) is the next commit's work, and the instruments that
// measure it have to be measured there too.
//
// WHY NOT A REFUSAL HERE: refusing would be the wide direction's opposite failure - it would stop
// every PrivateFile and PrivateFileExclusive caller dead on the shipped platform, break the four
// overwrite-promise pins for no security gain, and turn a work-in-progress split into a real
// regression. The honest interim is the old code, labelled as the old code.
package winsec

import (
	"io/fs"
	"os"
)

// leafPin is the transitional Windows pin: it carries the caller's spelling and holds no handle.
type leafPin struct {
	full string
}

// pinnedParent hands back the spelling unchanged. No syscall here on purpose: an extra probe would
// be a new answer this platform has not been measured with.
func pinnedParent(resolved string) (leafPin, error) {
	return leafPin{full: resolved}, nil
}

// close is a no-op because there is no handle to release yet.
func (p leafPin) close() error { return nil }

// openPinnedLeaf is the by-name open, verbatim from the route before this landing. `base` is
// unused in the transitional shape; the real chain uses it.
func openPinnedLeaf(p leafPin, _ string, flags int, perm fs.FileMode) (*os.File, error) {
	return os.OpenFile(p.full, flags, perm)
}

// unlinkPinnedLeaf is the by-name removal, verbatim from the route before this landing. This is the
// call ticket 348 booked: it re-resolves the spelling, so a swapped ancestor redirects it.
func unlinkPinnedLeaf(p leafPin, _ string) {
	_ = os.Remove(p.full)
}
