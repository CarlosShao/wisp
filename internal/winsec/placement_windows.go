//go:build windows

package winsec

import (
	"fmt"
	"path/filepath"
	"strings"
)

// platformVerifyPlacement is the Windows floor of the ticket 94 guarantee: what
// runs when no C26 pipeline is linked (internal/risk installs the full one, see
// internal/risk/winsec_c26.go). Like the portable half in resolve.go it can only
// refuse, never rewrite, so it is a placement check and not a second PathResolver.
//
// Three shapes it refuses, one per way "sealing A" becomes "modifying B":
//
//   - an extended-length (\\?\, \\.\) or UNC spelling, where the OS applies no
//     normalization at all, so a trailing dot or space is part of the name and
//     the entry created is one that no later ordinary path can reach;
//   - a component ending in '.' or ' ', which the normalizing namespace silently
//     strips, so the caller's string and the object created are different names;
//   - any existing component that carries FILE_ATTRIBUTE_REPARSE_POINT, which is
//     the junction case: everything below it lives in somebody else's tree, and
//     sealing "through" it narrows that tree instead of ours.
//
// Every existing component of the path is walked, prefix by prefix; a component
// that does not exist yet is not a link, which is the correct answer for the
// levels PrivateDirAll is about to create - and those are sealed parent first, by
// this same predicate, so a deeper component cannot silently become a link
// between the walk and the seal without somebody building it on purpose.
func platformVerifyPlacement(path string) (string, error) {
	if strings.HasPrefix(path, `\\?\`) || strings.HasPrefix(path, `\\.\`) {
		return "", fmt.Errorf("%w: %s uses an extended-length prefix, which turns off path normalization",
			ErrUnresolvedPath, path)
	}
	if vol := filepath.VolumeName(path); strings.HasPrefix(vol, `\\`) {
		return "", fmt.Errorf("%w: %s is a UNC share path, where the ACL this package writes may not be the share's own",
			ErrUnresolvedPath, path)
	}
	for _, comp := range pathComponents(path) {
		if comp == "" {
			continue
		}
		if last := comp[len(comp)-1]; last == '.' || last == ' ' {
			return "", fmt.Errorf("%w: component %q ends in %q, which the filesystem strips before it names the object",
				ErrUnresolvedPath, comp, string(last))
		}
	}
	// Ticket 108's AC#3 is this loop: it used to cut the path on the backslash
	// alone, so the same directory reached as "C:/Users/.../link/sub/x" presented
	// exactly one component to the walk, the leaf, which is the one thing the walk
	// does not check (privateDirAll may legitimately be about to create it).
	// PROBE P3 measured the consequence: the seal went through the junction and
	// rewrote a foreign DACL, and it did so in a binary that does not even link
	// internal/risk, so no resolver was involved. pathPieces gives the walk every
	// separator the OS honours, as an exact substring of the input, without
	// normalizing anything.
	for _, prefix := range pathPieces(path) {
		// isReparsePoint is the platform's own predicate, shared with
		// propagatePrivate above, so the walk and the seal agree on what a link
		// is. A missing component reads as "not a link", which is also the right
		// answer for the leg that stops here: everything below it is about to be
		// created by this call, parent first.
		if isReparsePoint(prefix) {
			return "", fmt.Errorf("%w: %s traverses a reparse point at %s, which is not the tree this call names",
				ErrUnresolvedPath, path, prefix)
		}
	}
	return path, nil
}
