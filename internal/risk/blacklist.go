package risk

import (
	"os"
	"path/filepath"
	"strings"
)

// A/B sensitive-path blacklists (SPEC-06 §4.1, F1 CRITICAL).
//
//   - A tier: absolutely forbidden, non-overridable. Read AND write denied.
//     No grant, session grant or config relaxation can unlock it.
//   - B tier: default deny; a single-file override grants only after one L2
//     strong confirm plus a log entry (emitted here via Logf).
//
// Classification always operates on the C26-resolved canonical path
// (SPEC-06 §4: whitelist AND blacklist membership both go through the same
// resolver), so junction/8.3/UNC spellings collapse before matching.

// Class is the blacklist classification of a resolved path.
type Class int

const (
	ClassNone Class = iota
	ClassB          // default deny + single-file override (one L2 confirm + log)
	ClassA          // absolute deny, non-overridable
)

func (c Class) String() string {
	switch c {
	case ClassA:
		return "A"
	case ClassB:
		return "B"
	default:
		return "none"
	}
}

// PathDecision is what the gate concludes for one path.
type PathDecision struct {
	Allow  bool
	NeedL2 bool // B tier without override: one strong L2 confirm unlocks this single file
	Class  Class
	Reason string
}

// Logf is the host-injected security log sink (nil = drop). Every A deny and
// every B decision (deny->L2 or single-file override) is logged here.
var Logf func(format string, args ...any)

func logf(format string, args ...any) {
	if Logf != nil {
		Logf(format, args...)
	}
}

// Classify classifies an already-resolved canonical path. Callers MUST pass
// the path through Resolve first; lexically-spelled inputs are matched on a
// best-effort normalized form only (defence in depth, not a bypass route).
//
// "Best effort" since ticket 72 means "the same handle pipeline": the input is
// expanded again (formsOf below) and so is every anchor, because a spelling
// the resolver never saw (a short-form USERPROFILE, say) must not be able to
// move a path out of A. The extra forms only ever turn a miss into a deny:
// the ALLOW side is deliberately narrower (see overrideApplies), because the
// basis for letting something through has to be the one spelling the OS vouches
// for, not any string that happens to look like it.
func Classify(canonical string) Class {
	c, _ := classifyWith(canonical)
	return c
}

// Gate applies the blacklist policy to a canonical path. bOverrides maps the
// normalized path of individually L2-confirmed B-tier files to true.
func Gate(canonical string, bOverrides map[string]bool) PathDecision {
	p := normPath(canonical)
	cand := formsOf(canonical)
	class, rule := classifyForms(cand)
	switch class {
	case ClassA:
		logf("risk: A-list DENY path=%s rule=%s (non-overridable)", p, rule)
		return PathDecision{Class: ClassA, Reason: "A-list sensitive path (non-overridable): " + rule}
	case ClassB:
		if overrideApplies(cand, bOverrides) {
			logf("risk: B-list OVERRIDE (single-file, L2 confirmed) path=%s rule=%s", p, rule)
			return PathDecision{Class: ClassB, Allow: true, Reason: "B-list single-file override (one L2 confirm + logged): " + rule}
		}
		logf("risk: B-list default DENY -> L2 path=%s rule=%s", p, rule)
		return PathDecision{Class: ClassB, NeedL2: true, Reason: "B-list sensitive path; single-file override requires one L2 confirm: " + rule}
	default:
		return PathDecision{Class: ClassNone, Allow: true}
	}
}

// overrideApplies reports whether the operator confirmed this exact file.
//
// Deliberately NARROWER than the deny side (编排者插单 2026-09-21 10:08, and
// the same asymmetry R17 asks for): classification walks every comparison
// form because a miss there silently downgrades A to B, but an override is an
// ALLOW, and an allow must rest on the one form the OS vouches for. Two
// consequences, both intentional:
//
//   - the incoming path must already BE the handle-resolved spelling. A caller
//     that hands in a spelling with a second form (an 8.3 short name, which on
//     a busy volume can alias a DIFFERENT directory, a junction, a \\?\ or UNC
//     variant) is asked to confirm again rather than waved through;
//   - the key must equal that resolved form, so a confirmation recorded for a
//     non-canonical string never unlocks the real file behind it either.
//
// When nothing on the path can be opened (a platform whose resolver is still
// DEFERRED, or a dead volume) there is exactly one form and no evidence of a
// second one, so the lookup falls back to it — that is the pre-ticket-72
// behavior, not a widening.
func overrideApplies(cand pathForms, bOverrides map[string]bool) bool {
	if len(bOverrides) == 0 {
		return false
	}
	if cand.real == "" {
		return cand.raw != "" && bOverrides[cand.raw]
	}
	if cand.raw != cand.real {
		return false // spelled differently from its own real path: confirm again
	}
	return bOverrides[cand.real]
}

// normPath folds a path into the comparison form: alternate separators unified
// into the platform one, lowercased, trailing separator trimmed. It is a
// COMPARISON fold and stays inside this package; the shape it produces must
// never be handed to the OS or printed as a path (ticket 75) - that shape is
// C26's canonical, which is platform-native by construction.
func normPath(p string) string {
	u := strings.ToLower(unifySeparators(p))
	// Trimming the root itself (POSIX "/", Windows "\") would erase the only
	// thing that makes the path absolute, so keep it.
	if t := strings.TrimSuffix(u, sepStr); t != "" || u != sepStr {
		return t
	}
	return u
}

// normDir normalizes an anchor directory for prefix matching.
func normDir(p string) string {
	if d := normPath(p); d != "" {
		return d
	}
	return p
}

// anchorPath puts an anchor directory and its literal child segments together
// in the comparison shape. Both sides arrive separator-unified (normPath /
// normDir above), so a plain sepStr join is all this is - and it has to be the
// platform separator, because the result goes into formsOf, which absolutizes
// and Lstats it. A backslash glued onto a POSIX home names a file that can
// never exist (ticket 75).
func anchorPath(dir string, parts ...string) string {
	out := dir
	for _, p := range parts {
		if out == sepStr {
			out += p // never build "//" off the POSIX root
			continue
		}
		out += sepStr + p
	}
	return out
}

func isUnder(path, dir string) bool {
	return path == dir || strings.HasPrefix(path, dir+sepStr)
}

func baseName(p string) string {
	if i := strings.LastIndex(p, sepStr); i >= 0 {
		return p[i+1:]
	}
	return p
}

// classifyWith matches the canonical path against the A and B tier rules.
// Anchors come from the live environment so tests can relocate the profile.
func classifyWith(canonical string) (Class, string) {
	return classifyForms(formsOf(canonical))
}

// classifyForms does the actual matching, on the comparison forms of one path
// (see pathresolver.go: formsOf / pathForms).
//
// Ticket 72 invariant, in one line: **the verdict must not depend on the
// spelling a path arrived in.** Both sides of every A-anchor comparison below
// are therefore handle-expanded, and a comparison whose provenance is not
// provably expanded fails CLOSED (aEquiv / aUnder report such a miss as a
// hit) instead of quietly dropping a path from A into B.
func classifyForms(cand pathForms) (Class, string) {
	p := cand.raw
	if p == "" {
		return ClassNone, ""
	}
	homeS := normDir(userHomeDir())
	appDataS := normDir(envOr("APPDATA", filepath.Join(userHomeDir(), "AppData", "Roaming")))
	localAppDataS := normDir(envOr("LOCALAPPDATA", filepath.Join(userHomeDir(), "AppData", "Local")))

	// Anchor side of the comparison: the same C26 handle pipeline the
	// candidate went through, so RUNNER~1 and runneradmin become one shape.
	// Each rule asks for its own child anchor (`<home>\.ssh` and friends),
	// because the expansion has to be proven per path, not per env var.
	// The join goes through anchorPath, i.e. the PLATFORM separator: these
	// strings are handed to formsOf, which absolutizes and Lstats them, so a
	// backslash glued onto a POSIX home names a file that cannot exist
	// (ticket 75). The human-readable rule strings below stay as written.

	// ---- A tier (absolute, non-overridable) ----
	switch {
	case homeS != "" && aEq(cand, anchorForms(anchorPath(homeS, `.git-credentials`))):
		return ClassA, "~/.git-credentials"
	case aAnyForm(cand, hasGitConfigSegment):
		return ClassA, ".git/config"
	case aUnder(cand, anchorForms(anchorPath(appDataS, `wisp`))) &&
		aEq(cand, anchorForms(anchorPath(appDataS, `wisp`, `config.toml`))):
		return ClassA, "%APPDATA%\\wisp\\config.toml"
	case homeS != "" && aUnder(cand, anchorForms(anchorPath(homeS, `.ssh`))):
		return ClassA, "~/.ssh/**"
	case homeS != "" && aEq(cand, anchorForms(anchorPath(homeS, `.aws`, `credentials`))):
		return ClassA, "~/.aws/credentials"
	case homeS != "" && aEq(cand, anchorForms(anchorPath(homeS, `.kube`, `config`))):
		return ClassA, "~/.kube/config"
	case aAnyForm(cand, func(s string) bool { return isBrowserCredentialStore(baseName(s)) }):
		return ClassA, "browser credential store"
	case appDataS != "" && aUnder(cand, anchorForms(anchorPath(appDataS, `microsoft`, `protect`))):
		return ClassA, "%APPDATA%\\Microsoft\\Protect\\** (DPAPI master keys)"
	case localAppDataS != "" && aUnder(cand, anchorForms(anchorPath(localAppDataS, `microsoft`, `credentials`))):
		return ClassA, "%LOCALAPPDATA%\\Microsoft\\Credentials\\**"
	}

	// ---- B tier (default deny + single-file override) ----
	// The first form is exactly the spelling the caller handed in (which is
	// what this function has always matched), so no B verdict can be lost;
	// later forms only add hits, which for this tier means "asks first".
	for _, s := range cand.spellings() {
		base := strings.ToLower(baseName(s))
		switch {
		case strings.HasPrefix(base, ".env"):
			return ClassB, ".env*"
		case strings.HasSuffix(base, ".pem"):
			return ClassB, "*.pem"
		case strings.HasSuffix(base, ".p12"), strings.HasSuffix(base, ".pfx"):
			return ClassB, "*.p12|*.pfx"
		case strings.HasPrefix(base, "id_"):
			return ClassB, "id_*"
		case strings.HasPrefix(base, "secrets."):
			return ClassB, "secrets.*"
		case strings.Contains(base, "credentials") && strings.HasSuffix(base, ".json"):
			return ClassB, "*credentials*.json"
		}
	}
	return ClassNone, ""
}

// aEq and aUnder are the two A-anchor relations. They are spelled separately
// from isUnder / == so that the fail-closed clause below applies to security
// comparisons only and never to a plain path test.
func aEq(cand, anchor pathForms) bool { return aRelates(cand, anchor, eqPath) }

func aUnder(cand, anchor pathForms) bool { return aRelates(cand, anchor, isUnder) }

func eqPath(p, dir string) bool { return p == dir }

// aRelates matches every comparison form of the candidate against every
// comparison form of the anchor, and — when that finds nothing — asks whether
// the miss is even trustworthy (uncertainAnchorMiss).
func aRelates(cand, anchor pathForms, rel func(p, dir string) bool) bool {
	for _, c := range cand.spellings() {
		if c == "" {
			continue
		}
		for _, a := range anchor.spellings() {
			if a == "" {
				continue
			}
			if rel(c, a) {
				return true
			}
		}
	}
	return uncertainAnchorMiss(cand, anchor)
}

// aAnyForm applies a shape test (path segments, base name) to every
// comparison form the path can take.
func aAnyForm(cand pathForms, test func(string) bool) bool {
	for _, s := range cand.spellings() {
		if s != "" && test(s) {
			return true
		}
	}
	return false
}

// uncertainAnchorMiss is R17's fail-closed clause: "if a side cannot be proven
// expanded, the only allowed behavior is fail-closed". A miss between an A
// anchor and a candidate is only evidence of absence when BOTH sides are
// proven handle-resolved real paths. If one side had to stop its expansion
// below a component that exists but refused to open (a deny-filtered
// directory, a dangling reparse point), the two paths might be the same tree
// and the answer is "treat it as protected".
//
// The clause is bounded on purpose, so it stays a fail-closed tie-breaker
// rather than a blanket deny:
//   - it needs a resolved root on both sides. A path on a dead volume, or one
//     on a platform whose resolver is still DEFERRED (pathresolver_other.go),
//     has no root and cannot hide inside a live anchor tree;
//   - the roots must be ancestry-related, i.e. one tree could contain the
//     other. Unrelated trees are a proven miss.
func uncertainAnchorMiss(cand, anchor pathForms) bool {
	if cand.root == "" || anchor.root == "" {
		return false
	}
	if cand.certain() && anchor.certain() {
		return false
	}
	return isUnder(cand.root, anchor.root) || isUnder(anchor.root, cand.root)
}

// hasGitConfigSegment reports whether the path is exactly <...>/.git/config.
// It is fed comparison forms (normPath output), so one separator is enough.
func hasGitConfigSegment(p string) bool {
	segs := strings.Split(p, sepStr)
	for i, s := range segs {
		if s == ".git" && i+2 == len(segs) && segs[i+1] == "config" {
			return true
		}
	}
	return false
}

// isBrowserCredentialStore matches the credential databases of Chromium-based
// browsers by file name (profiles move; the fail-safe direction is a
// base-name match anywhere).
func isBrowserCredentialStore(base string) bool {
	switch base {
	case "login data", "cookies", "web data", "local state":
		return true
	}
	return false
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func userHomeDir() string {
	if h, err := os.UserHomeDir(); err == nil && h != "" {
		return h
	}
	return ""
}
