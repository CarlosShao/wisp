package tools

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/CarlosShao/wisp/internal/risk"
)

// D22/CI-banned patterns (filepath.Clean, filepath.Abs) appear NOWHERE in this
// file: every canonicalization below delegates to the frozen C26 pipeline in
// risk.Resolve, which is the codebase's only sanctioned normalization entry
// (SPEC-06 §4). That is also what makes the R2/R3 verdicts and the bytes a
// tool actually opens describe the same path.

// PathCanonicalizer is the bridge's C19 PathCanonicalizer (ticket 17's frozen
// seam) over ticket 18's C26 resolver plus the [fs] allowed_dirs allowlist.
//
// Until this type existed, nothing in the repository implemented the seam:
// risk.Resolve and risk.Classify were both landed and tested, but the assessor
// still had R2/R3 DORMANT because no production composition fed them. The
// bridge is the choke point, so the adapter belongs here, and wiring it is
// what makes "out-of-allowlist read -> L2" (D34) a real verdict instead of a
// documented intention.
type PathCanonicalizer struct {
	exceptions []string
	mu         sync.RWMutex
	// roots are the canonicalized [fs] allowed_dirs entries.
	roots []string
	// unusable records roots that could not be canonicalized: they authorize
	// nothing (fail-closed) and the reason is kept for the audit line.
	unusable []string
}

// NewPathCanonicalizer resolves the allowlist roots through C26 once and
// returns the adapter. An EMPTY allowlist authorizes nothing: every path is
// out-of-scope, so every fs call lands at L2 (R2) rather than passing through
// unjudged.
func NewPathCanonicalizer(allowedDirs, reparseExceptions []string) *PathCanonicalizer {
	p := &PathCanonicalizer{exceptions: append([]string(nil), reparseExceptions...)}
	for _, d := range allowedDirs {
		c, err := p.resolve(d)
		if err != nil || c == "" {
			p.unusable = append(p.unusable, fmt.Sprintf("%q: %v", d, err))
			continue
		}
		p.roots = append(p.roots, foldPath(c))
	}
	return p
}

// resolve is the ONE canonicalization call site in this package.
func (p *PathCanonicalizer) resolve(raw string) (string, error) {
	res, err := risk.Resolve(raw, p.exceptions)
	if err != nil {
		return "", err
	}
	return res.Canonical, nil
}

// Canonicalize implements risk.PathCanonicalizer (and the narrower contract
// builtin tools use to open what the judge looked at). A reparse traversal
// returns risk.ErrReparseDenied, which R2 turns into a fail-closed L2.
func (p *PathCanonicalizer) Canonicalize(raw string) (string, error) {
	if strings.TrimSpace(raw) == "" {
		return "", fmt.Errorf("tools: empty path")
	}
	return p.resolve(raw)
}

// InAllowlist implements risk.PathCanonicalizer against already-canonical
// input, component-boundary and case-folded (Windows is the only shipping
// platform, and SPEC-06 §4 treats canonical paths as case-insensitive).
func (p *PathCanonicalizer) InAllowlist(canonical string) bool {
	f := foldPath(canonical)
	if f == "" {
		return false
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, r := range p.roots {
		if f == r || strings.HasPrefix(f, r+pathSep) {
			return true
		}
	}
	return false
}

// Roots returns the canonical allowlist roots in effect (audit surface).
func (p *PathCanonicalizer) Roots() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return append([]string(nil), p.roots...)
}

// UnusableRoots lists configured roots that could not be canonicalized; they
// authorize nothing, which is the fail-closed direction.
func (p *PathCanonicalizer) UnusableRoots() []string {
	return append([]string(nil), p.unusable...)
}

// pathSep is the comparison separator: the PLATFORM one. C26's canonical is the
// real path of the real file, which on POSIX is '/'-shaped and on Windows
// '\'-shaped (ticket 75); a hard-coded '\\' here used to make every comparison
// string on Linux name a path no OS call can open.
const pathSep = string(filepath.Separator)

// unifySeparators folds '/' into '\' on Windows, where the OS treats the two as
// one, and leaves the string untouched on POSIX, where '\' is an ordinary
// filename character (folding it there would merge two different files into one
// comparison key). The branch is a compile-time constant on each platform, so
// Windows output is byte-identical to the pre-ticket-75 code.
func unifySeparators(p string) string {
	if filepath.Separator == '\\' {
		return strings.ReplaceAll(p, "/", `\`)
	}
	return p
}

// foldPath puts a canonical path into the comparison form: alternate separators
// unified into the platform one, lowercased, trailing separators trimmed. It is
// a COMPARISON fold, not a normalization - it neither cleans nor absolutizes,
// and calling it on user input instead of Canonicalize would be the C26 bypass
// D22 forbids.
func foldPath(p string) string {
	u := strings.ToLower(unifySeparators(p))
	for len(u) > 1 && strings.HasSuffix(u, pathSep) {
		u = u[:len(u)-1]
	}
	return u
}

// sensitiveClassifier is the C19 SensitiveClassifier seam over ticket 18's
// A/B blacklist (risk.Classify).
type sensitiveClassifier struct{}

// NewSensitiveClassifier wires R3 against the SPEC-06 §4.1 blacklist.
func NewSensitiveClassifier() risk.SensitiveClassifier { return sensitiveClassifier{} }

// Classify implements risk.SensitiveClassifier. The tier mapping is direct:
// risk.ClassA is the absolute deny tier, ClassB the default-deny/one-L2-confirm
// tier.
func (sensitiveClassifier) Classify(canonical string) risk.Tier {
	switch risk.Classify(canonical) {
	case risk.ClassA:
		return risk.TierA
	case risk.ClassB:
		return risk.TierB
	default:
		return risk.TierNone
	}
}
