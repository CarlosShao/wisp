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
	// rewritten records roots whose spelling C26's expansion step substituted
	// (ticket 102): the authorized tree is not the tree the config spelled, and
	// an operator reading the config would not recognize it. Bounded by the
	// number of [fs] allowed_dirs entries.
	rewritten []string
}

// NewPathCanonicalizer resolves the allowlist roots through C26 once and
// returns the adapter. An EMPTY allowlist authorizes nothing: every fs call
// lands at L2 (R2) rather than passing through unjudged.
func NewPathCanonicalizer(allowedDirs, reparseExceptions []string) *PathCanonicalizer {
	p := &PathCanonicalizer{exceptions: append([]string(nil), reparseExceptions...)}
	for _, d := range allowedDirs {
		res, err := p.resolve(d)
		if err != nil || res.Canonical == "" {
			p.unusable = append(p.unusable, fmt.Sprintf("%q: %v", d, err))
			continue
		}
		// Ticket 102 (fix B) consumption site, authorization leg. Expansion is
		// contract-required for a root spelled ~/... or %VAR%\... (SPEC-06 §4
		// step 1), so a rewrite is not refused outright - but a root that C26
		// MOVED and that the OS cannot confirm as an existing tree is dropped:
		// authorizing a tree nobody can point at is a fail-open, and the config
		// line that produced it is kept for the operator.
		if res.Rewritten {
			p.mu.Lock()
			p.rewritten = append(p.rewritten, fmt.Sprintf("%q -> %s (%s)",
				d, res.Canonical, strings.Join(res.Rewrites, "+")))
			p.mu.Unlock()
			if !res.Resolved {
				p.unusable = append(p.unusable, fmt.Sprintf(
					"%q: C26 expanded it onto %s but that tree is not confirmed on disk, authorizing nothing", d, res.Canonical))
				continue
			}
		}
		p.roots = append(p.roots, foldPath(res.Canonical))
	}
	return p
}

// resolve is the ONE canonicalization call site in this package. It hands back
// the whole C26 Result rather than just a spelling on purpose: the two callers
// read ticket 102's account differently (roots are authorization, so a rewrite
// is recorded and fail-closed; a target's judged tree is the same tree the tool
// opens and prints, so a rewrite is only traceable through Result.Spelling).
func (p *PathCanonicalizer) resolve(raw string) (risk.Result, error) {
	res, err := risk.Resolve(raw, p.exceptions)
	if err != nil {
		return risk.Result{}, err
	}
	return res, nil
}

// Canonicalize implements risk.PathCanonicalizer (and the narrower contract
// builtin tools use to open what the judge looked at). A reparse traversal
// returns risk.ErrReparseDenied, which R2 turns into a fail-closed L2.
//
// Ticket 102 disposition for this leg: NOT a refusal. fs.go / fs_write.go /
// mode.go open exactly what this function returns, and rules_gateway's R2
// verdict and audit line are computed from the same string, so judgement,
// execution and report all name ONE tree even when expansion moved it.
func (p *PathCanonicalizer) Canonicalize(raw string) (string, error) {
	if strings.TrimSpace(raw) == "" {
		return "", fmt.Errorf("tools: empty path")
	}
	res, err := p.resolve(raw)
	if err != nil {
		return "", err
	}
	return res.Canonical, nil
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

// RewrittenRoots lists the [fs] allowed_dirs entries C26's expansion step moved
// onto another tree, as `"spelling" -> tree (constructs)` (ticket 102). The
// audit surface is what makes the (B) account more than a struct field: a root
// that authorizes a tree other than the one written down is operator-visible.
func (p *PathCanonicalizer) RewrittenRoots() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return append([]string(nil), p.rewritten...)
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
