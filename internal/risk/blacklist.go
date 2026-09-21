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
// move a path out of A. The added forms only ever turn a miss into a deny;
// Gate's B-tier override matches them too, which is the same-file question
// asked the other way round — the confirmation was for a file, not for one
// spelling of it.
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
// Every comparison form of the path is checked, so a B grant recorded under
// the handle-resolved spelling still covers the caller's own spelling of the
// same file (and vice versa) — it is the same file either way.
func overrideApplies(cand pathForms, bOverrides map[string]bool) bool {
	if len(bOverrides) == 0 {
		return false
	}
	for _, s := range cand.spellings() {
		if s != "" && bOverrides[s] {
			return true
		}
	}
	return false
}

// normPath folds a path into the comparison form: forward slashes unified,
// lowercased, trailing separators trimmed.
func normPath(p string) string {
	u := strings.ReplaceAll(p, "/", `\`)
	return strings.TrimSuffix(strings.ToLower(u), `\`)
}

// normDir normalizes an anchor directory for prefix matching.
func normDir(p string) string {
	if d := normPath(p); d != "" {
		return d
	}
	return p
}

func isUnder(path, dir string) bool {
	return path == dir || strings.HasPrefix(path, dir+`\`)
}

func baseName(p string) string {
	if i := strings.LastIndex(p, `\`); i >= 0 {
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

	// ---- A tier (absolute, non-overridable) ----
	switch {
	case homeS != "" && aEq(cand, anchorForms(homeS+`\.git-credentials`)):
		return ClassA, "~/.git-credentials"
	case aAnyForm(cand, hasGitConfigSegment):
		return ClassA, ".git/config"
	case aUnder(cand, anchorForms(appDataS+`\wisp`)) &&
		aEq(cand, anchorForms(appDataS+`\wisp\config.toml`)):
		return ClassA, "%APPDATA%\\wisp\\config.toml"
	case homeS != "" && aUnder(cand, anchorForms(homeS+`\.ssh`)):
		return ClassA, "~/.ssh/**"
	case homeS != "" && aEq(cand, anchorForms(homeS+`\.aws\credentials`)):
		return ClassA, "~/.aws/credentials"
	case homeS != "" && aEq(cand, anchorForms(homeS+`\.kube\config`)):
		return ClassA, "~/.kube/config"
	case aAnyForm(cand, func(s string) bool { return isBrowserCredentialStore(baseName(s)) }):
		return ClassA, "browser credential store"
	case appDataS != "" && aUnder(cand, anchorForms(appDataS+`\microsoft\protect`)):
		return ClassA, "%APPDATA%\\Microsoft\\Protect\\** (DPAPI master keys)"
	case localAppDataS != "" && aUnder(cand, anchorForms(localAppDataS+`\microsoft\credentials`)):
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
func hasGitConfigSegment(p string) bool {
	segs := strings.Split(p, `\`)
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
