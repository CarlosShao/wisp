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
func Classify(canonical string) Class {
	c, _ := classifyWith(canonical)
	return c
}

// Gate applies the blacklist policy to a canonical path. bOverrides maps the
// normalized path of individually L2-confirmed B-tier files to true.
func Gate(canonical string, bOverrides map[string]bool) PathDecision {
	p := normPath(canonical)
	class, rule := classifyWith(canonical)
	switch class {
	case ClassA:
		logf("risk: A-list DENY path=%s rule=%s (non-overridable)", p, rule)
		return PathDecision{Class: ClassA, Reason: "A-list sensitive path (non-overridable): " + rule}
	case ClassB:
		if bOverrides != nil && bOverrides[p] {
			logf("risk: B-list OVERRIDE (single-file, L2 confirmed) path=%s rule=%s", p, rule)
			return PathDecision{Class: ClassB, Allow: true, Reason: "B-list single-file override (one L2 confirm + logged): " + rule}
		}
		logf("risk: B-list default DENY -> L2 path=%s rule=%s", p, rule)
		return PathDecision{Class: ClassB, NeedL2: true, Reason: "B-list sensitive path; single-file override requires one L2 confirm: " + rule}
	default:
		return PathDecision{Class: ClassNone, Allow: true}
	}
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
	p := normPath(canonical)
	if p == "" {
		return ClassNone, ""
	}
	home := normDir(userHomeDir())
	appData := normDir(envOr("APPDATA", filepath.Join(userHomeDir(), "AppData", "Roaming")))
	localAppData := normDir(envOr("LOCALAPPDATA", filepath.Join(userHomeDir(), "AppData", "Local")))

	// ---- A tier (absolute, non-overridable) ----
	switch {
	case home != "" && p == home+`\.git-credentials`:
		return ClassA, "~/.git-credentials"
	case hasGitConfigSegment(p):
		return ClassA, ".git/config"
	case isUnder(p, normDir(appData+`\wisp`)) && p == normDir(appData+`\wisp\config.toml`):
		return ClassA, "%APPDATA%\\wisp\\config.toml"
	case home != "" && isUnder(p, home+`\.ssh`):
		return ClassA, "~/.ssh/**"
	case home != "" && p == home+`\.aws\credentials`:
		return ClassA, "~/.aws/credentials"
	case home != "" && p == home+`\.kube\config`:
		return ClassA, "~/.kube/config"
	case isBrowserCredentialStore(baseName(p)):
		return ClassA, "browser credential store"
	case appData != "" && isUnder(p, appData+`\microsoft\protect`):
		return ClassA, "%APPDATA%\\Microsoft\\Protect\\** (DPAPI master keys)"
	case localAppData != "" && isUnder(p, localAppData+`\microsoft\credentials`):
		return ClassA, "%LOCALAPPDATA%\\Microsoft\\Credentials\\**"
	}

	// ---- B tier (default deny + single-file override) ----
	base := strings.ToLower(baseName(p))
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
	return ClassNone, ""
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
