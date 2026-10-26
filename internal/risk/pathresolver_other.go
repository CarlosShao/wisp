//go:build !windows

package risk

import (
	"os"
	"path/filepath"
	"strings"
)

// C26's handle step on the POSIX platforms (ticket 347).
//
// PLAN.md:2385 spells the macOS half of C26 as "realpath + lstat, symlink
// rejected", and D7 (PLAN.md:184-199) puts macOS inside the support matrix, so
// this is a promised behaviour, not a port to be scheduled. Until ticket 347
// this file held two stubs returning ("", false) / nil, which made the A tier's
// fail-closed branch (uncertainAnchorMiss) structurally unreachable on every
// POSIX platform while test-core ran this package there.
//
// The two primitives split the way SPEC-06 §4 splits them:
//
//   - resolveHandle = realpath. EvalSymlinks gives the same thing realpath(3)
//     does: every symlink in the chain resolved and the result freed of "." and
//     "..". It fails exactly when the path cannot be resolved to something that
//     exists (a missing component, a dangling link, a permission refusal), which
//     is what makes "handle resolution succeeded" mean "the OS vouches for this
//     spelling" - the property Result.Resolved and pathForms.root are read for.
//   - reparseComponents = the lstat walk that REFUSES symlinks. Resolution
//     alone would FOLLOW a link, i.e. be more permissive than the lexical
//     pipeline it replaces and a fail-open against SPEC-06 §4's default deny.
//     Resolve runs this walk first and denies any traversed symlink that is not
//     listed in [fs] reparse_point_exceptions, so a link is not followed into
//     an unlisted tree by the step below. Both halves are needed: one without
//     the other is either the old stub or an unconditional follow.

// resolveHandle returns the real path of p through realpath, and false when
// there is nothing to resolve. It follows symlinks; refusing them is
// reparseComponents' job and it runs first (Resolve's pipeline order).
func resolveHandle(p string) (string, bool) {
	if p == "" {
		return "", false
	}
	real, err := filepath.EvalSymlinks(p)
	if err != nil || real == "" {
		// EvalSymlinks("") is ("", nil), hence the guard above; every other
		// failure means "this spelling names nothing the OS can resolve yet",
		// which is the lexical-fallback case, not an error.
		return "", false
	}
	return real, true
}

// reparseComponents walks every component of p from the root down to the leaf
// and returns those that are symlinks, in pipeline spelling. A symlink on ANY
// traversed component triggers SPEC-06 §4's default-deny rule (same contract as
// the Windows FILE_ATTRIBUTE_REPARSE_POINT walk); Resolve is the caller that
// turns a non-empty result into ErrReparseDenied unless the exact component is
// an exception.
//
// Lstat is deliberate: it reports the component AS SPELLED and does not follow
// the link, so the returned spelling is the one the exception table and the
// audit line must carry. A component that does not exist (or cannot be Lstat-ed)
// ends the walk: nothing below it can be a traversed symlink that exists.
//
// p is always an absolute, lexically cleaned spelling here - C26 step 2-3 ran
// before this - so a relative input carries no root to walk from and yields nil.
func reparseComponents(p string) []string {
	if p == "" || !filepath.IsAbs(p) {
		return nil
	}
	var comps []string
	base := ""
	for _, part := range strings.Split(p, sepStr) {
		if part == "" {
			continue // leading separator, or a doubled one Clean already folded
		}
		if base == "" {
			base = sepStr + part
		} else {
			base += sepStr + part
		}
		fi, err := os.Lstat(base)
		if err != nil {
			break
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			comps = append(comps, base)
		}
	}
	return comps
}
