package tools

import (
	"github.com/CarlosShao/wisp/internal/risk"
)

// The machine-readable refusal classes (ticket 20 frame 7 (a)).
//
// WHY A TYPED VALUE AND NOT A WORDING. Until this file, the two refusal
// families below were told apart only by their human Reason text, and
// risk.ErrReparseDenied was folded into a string on the way out
// (paths.go:125's own comment). That is enough for a card, and not enough for
// a decision a host has to make: the [fs] allowed_dirs first-use ask ("把这个
// 目录加进 [fs] allowed_dirs?") may be derived from ONE of these classes and
// must never be derived from the other. A refusal that reads like prose cannot
// carry that rule, because a host would have to match on wording that is not a
// contract.
//
// THE RULE, IN ONE SENTENCE. Only RefusalOutOfAllowlist may ever be turned
// into the first-use ask; RefusalRedTeamDenied must never reach one. The
// reason is the frame's own: a document injected into the model's context can
// name a path that traverses a junction, and if that refusal looked like
// "just out of scope", driving Wisp to ask the user "authorize this
// directory?" would turn the red-team gate into a social-engineering mouth.
//
// WHICH REFUSALS LAND WHERE (the frozen rulings this mapping follows):
//   - Every canonicalization that the C26 resolver itself refuses (a
//     junction/symlink traversal, and anything else the pipeline answers with
//     an error) is RefusalRedTeamDenied. Fail-closed by construction: a path
//     nobody could judge is never askable.
//   - A Deny verdict (R3 tier A, absolute) is RefusalRedTeamDenied whatever
//     its spelling was. This is where ticket 18's other three red-team shapes
//     land: the 8.3 short name, the UNC spelling and the \\?\ spelling of an
//     A-list file all resolve to the long path and are denied by
//     classification (risk/pathresolver_junction_windows_test.go Cases 2-4).
//   - A path that canonicalized cleanly and simply sits under no allowed root
//     is RefusalOutOfAllowlist. The same ruling covers every benign spelling
//     of such a target: bridge_junction_windows_test.go's
//     TestShortNameSpellingGetsTheSameVerdictAsTheLongOne pins verbatim that
//     "8.3 is deliberately NOT a resolver-refused class", and the long and
//     short spellings of an out-of-scope read must land on the same R2 row.
//     The ask class therefore follows the target, not the spelling.
type RefusalKind string

// The refusal classes. RefusalNone is the zero value: the call was not refused
// (or was refused for a reason with no askable counterpart, e.g. a missing
// capability).
const (
	RefusalNone RefusalKind = ""
	// RefusalOutOfAllowlist is the ONLY class a [fs] allowed_dirs first-use
	// ask may ever be derived from: every path canonicalized, and at least one
	// of them is under no configured root.
	RefusalOutOfAllowlist RefusalKind = "out-of-allowlist"
	// RefusalRedTeamDenied is the class that must NEVER become an ask: the
	// C26 resolver refused a path (junction/symlink traversal without an
	// exception), or the verdict is Deny (R3 tier A, absolute - the landing
	// spot of the 8.3 / UNC / \\?\ spellings of A-list files).
	RefusalRedTeamDenied RefusalKind = "redteam-denied"
)

// refusalKindOf classifies a refusal over the SAME raw paths C19 judged, in
// the SAME order (the bridge hands them to risk.Assess as Facts.Paths, and R2
// canonicalizes each one). Re-running the canonicalization here does not open
// a second judge: the answer this function produces decides nothing about the
// call - the call is already refused - it only names WHICH refusal happened,
// for a host that has to decide whether a first-use ask is admissible at all.
//
// level is the EFFECTIVE level routing refused on (ticket 90's screened
// level), and the Deny limb is checked first so an absolute denial can never
// be re-labelled by path shape. A call with no path parameters has nothing to
// ask about and comes back RefusalNone.
func (b *Bridge) refusalKindOf(rawPaths []string, level risk.Level) RefusalKind {
	if level == risk.Deny {
		return RefusalRedTeamDenied
	}
	if level != risk.L1 && level != risk.L2 {
		// Not a level that consults a gate (L0 never reaches a refusal), i.e.
		// a broken judge's impossible level: fail closed to "nothing askable".
		return RefusalNone
	}
	if b.paths == nil || len(rawPaths) == 0 {
		return RefusalNone
	}
	outside := false
	for _, raw := range rawPaths {
		canonical, err := b.paths.Canonicalize(raw)
		if err != nil {
			// The resolver refused this path. Whatever it was, it is not the
			// "cleanly resolved and merely out of scope" shape that an ask is
			// allowed to offer a fix for.
			return RefusalRedTeamDenied
		}
		if !b.paths.InAllowlist(canonical) {
			outside = true
		}
	}
	if outside {
		return RefusalOutOfAllowlist
	}
	return RefusalNone
}

// refusalToken renders one class for the audit line. The empty class is
// spelled out rather than left blank so a reader can tell "the bridge decided
// there is none" from "this line predates the field".
//
// The distinctness of the two non-empty tokens is deliberately NOT expressed
// as a compile-time trick here - it is pinned by the test suite
// (Test20r2RefusalKindsAreDistinct), because what must stay true is not "the
// constants differ" but "the two REFUSALS carry different values end to end",
// and that is a reading, not a spelling.
func (k RefusalKind) refusalToken() string {
	if k == RefusalNone {
		return "none"
	}
	return string(k)
}
