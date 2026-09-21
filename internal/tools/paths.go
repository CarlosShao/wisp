package tools

import (
	"errors"
	"fmt"
	"os"
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
	// workspace is the narrowed root a panel-side workspace switch set
	// (ticket 92 AC#3). Empty means "nothing was narrowed": the config roots
	// alone decide scope, which is every run that never touched a composer.
	// It can only ever TIGHTEN InAllowlist - a workspace outside the roots is
	// refused at switch time - so no workspace choice can widen authority.
	workspace string
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
		//
		// Ticket 107 (round 1): "the OS cannot confirm" had to be stated in a
		// capability the running platform actually has, because res.Resolved is
		// handle-based and only ever true on Windows (risk.resolveHandle is a
		// stub elsewhere), which dropped EVERY expanded root on POSIX.
		// Ticket 107 (round 2, rejected round 1's answer): confirmation may not
		// be downgraded to "something exists under that name" - os.Stat follows
		// links, so a root spelled %VAR%\proj that is a link onto another tree
		// would authorize a tree the operator never named (acceptance probe A).
		// The release side therefore gets a form that has actually been
		// resolved: treeResolvedAsNamed is true only when the OS says the named
		// path IS the tree on disk, with no link crossed. res.Resolved keeps
		// priority so Windows verdicts stay byte-identical to ticket 102's (a
		// handle-resolved tree is never re-probed). Where a platform cannot
		// resolve, the root is dropped and recorded - stricter, not looser.
		if res.Rewritten {
			p.mu.Lock()
			p.rewritten = append(p.rewritten, fmt.Sprintf("%q -> %s (%s)",
				d, res.Canonical, strings.Join(res.Rewrites, "+")))
			p.mu.Unlock()
			if !res.Resolved && !treeResolvedAsNamed(res.Canonical) {
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
	if !rootsContain(p.roots, f) {
		return false
	}
	// Ticket 107 (round 2): the release side also recognizes the RESOLVED form,
	// and it must contain THAT under the same root. The lexical spelling alone
	// cannot tell "inside the allowed tree" from "inside a link that hangs off
	// it and lands elsewhere" (acceptance probe C: root/proj is real, the target
	// goes out through proj/esc). Requiring both is strictly narrower: every
	// path that passed before either fails to resolve or still resolves inside
	// the root. Unresolvable means unauthorized, so a platform whose resolver
	// cannot see a link (Windows reports no error and no switch for a junction,
	// measured) still refuses instead of vouching for a tree it did not look at.
	rf, ok := resolvedForm(canonical)
	if !ok || !rootsContain(p.roots, foldPath(rf)) {
		return false
	}
	// Ticket 92 AC#3: once a workspace is chosen, scope narrows to it. An empty
	// workspace is "nothing chosen", so every judgement before the first switch
	// is byte-identical to what it was.
	if p.workspace != "" && !rootsContain([]string{p.workspace}, f) {
		return false
	}
	return true
}

// rootsContain is the component-boundary, case-folded containment test over a
// set of already-folded roots. It takes no lock so the locked callers can
// compose it (InAllowlist) without re-entering p.mu.
func rootsContain(roots []string, folded string) bool {
	for _, r := range roots {
		if folded == r || strings.HasPrefix(folded, r+pathSep) {
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

// treeOnDisk answers the WEAKER question: does something exist under this name,
// following links on the way. That is the right strength where a name is only
// being checked for existence (the workspace leg of this package confirms the
// directory a picker chose), and it is the wrong strength for an authorization:
// ticket 107 round 1 used it there and acceptance probe A showed a root that is
// a link onto another tree then authorized the tree nobody named. New callers
// that gate authority must use treeResolvedAsNamed instead.
func treeOnDisk(canonical string) bool {
	st, err := os.Stat(canonical)
	return err == nil && st.IsDir()
}

// treeResolvedAsNamed answers ticket 102's confirmation question with the
// strongest form the release side is allowed to use: the OS says this directory
// exists AND the name the operator wrote is that directory, not a link onto
// something else. It is NOT a normalization - it neither cleans nor absolutizes,
// so it stays outside C26's single-entry rule - and it is only consulted where
// res.Resolved cannot speak (POSIX, where handle-based resolution is deferred;
// see risk.pathresolver_other.go).
//
// Round 1 of this ticket used os.Stat here, which follows links: a root spelled
// %VAR%/proj whose proj is a link onto another tree was accepted and the tree
// the tool then opened was the one nobody named. Equality against the resolved
// spelling is what closes that, and it fails closed on every platform: a machine
// where the resolver cannot answer gets no authorization, only a line in
// UnusableRoots.
func treeResolvedAsNamed(canonical string) bool {
	r, err := filepath.EvalSymlinks(canonical)
	if err != nil || foldPath(r) != foldPath(canonical) {
		return false
	}
	st, err := os.Stat(r)
	return err == nil && st.IsDir()
}

// resolvedForm puts a canonical path into the one shape the release side trusts:
// every link followed. It is only ever used as an ADDITIONAL condition next to
// the lexical comparison, so a platform whose resolver answers nothing usable
// can only lose authorizations, never gain them.
//
// The tail of a path may legitimately not exist yet (a write that creates a new
// file), so a missing component is walked past and re-appended; but "missing"
// and "there, and I could not look through it" report the same ENOENT-looking
// error on both platforms - Windows answers ERROR_PATH_NOT_FOUND for a path
// crossing a junction that EvalSymlinks does not resolve at all. os.Lstat on the
// failing component tells them apart: it succeeds exactly when the path is there
// but unresolvable, which is the shape this function refuses.
func resolvedForm(p string) (string, bool) {
	var tail []string
	cur := p
	for {
		if r, err := filepath.EvalSymlinks(cur); err == nil {
			for i := len(tail) - 1; i >= 0; i-- {
				r = filepath.Join(r, tail[i])
			}
			return r, true
		}
		if _, lerr := os.Lstat(cur); !errors.Is(lerr, os.ErrNotExist) {
			return "", false
		}
		dir := filepath.Dir(cur)
		if dir == cur {
			return "", false
		}
		tail = append(tail, filepath.Base(cur))
		cur = dir
	}
}

// pathSep is the comparison separator: the PLATFORM one. C26's canonical is the
// real path of the real file, which on POSIX is '/'-shaped and on Windows
// '\'-shaped (ticket 75); a hard-coded '\\' here used to make every comparison
// string on Linux name a path no OS call can open.
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
