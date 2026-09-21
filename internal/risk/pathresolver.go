package risk

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

// C26 PathResolver (SPEC-06 §4): the ONLY path normalization entry point in
// the codebase. Pipeline order is fixed and must not be reordered:
//
//	expand (env / ~) -> absolutize -> Clean -> handle-based real path
//	(GetFinalPathNameByHandle VOLUME_NAME_DOS) -> FILE_ATTRIBUTE_REPARSE_POINT
//	detection with default DENY -> 8.3 short-name expansion -> UNC normalization
//
// "Parse-then-continue" on reparse points is forbidden (SPEC-06 §4): a
// junction whose source sits outside and target inside (or vice versa) is
// only caught by refusing reparse traversal itself, so the resolver denies
// any path that traverses a reparse point unless that exact component is
// listed in the user's [fs] reparse_point_exceptions.
//
// D22: filepath.Clean / filepath.Abs are sanctioned ONLY inside this file.
// Any other use for fs decisions is a CI-failing violation (see
// scripts/check-pathclean-ban.sh).
//
// The same pipeline also serves the ANCHOR side of a security comparison
// (formsOf / anchorForms below): classification must not depend on which
// spelling a path happened to arrive in.

// ErrReparseDenied is returned when a path traverses a junction/symlink that
// is not covered by reparse_point_exceptions. Fail-closed by design.
var ErrReparseDenied = errors.New("risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions")

// Result is the outcome of resolving one path through the C26 pipeline.
type Result struct {
	// Canonical is the real path after handle resolution (8.3 expanded,
	// reparse targets resolved, UNC and \\?\ prefixes normalized). When the
	// path does not exist yet, it is the lexically cleaned spelling —
	// classification still applies so nonexistent sensitive targets cannot
	// be probed into existence.
	Canonical string
	// Reparse reports whether the INPUT traversed any reparse point that
	// was covered by an explicit exception (non-excepted traversal returns
	// ErrReparseDenied instead).
	Reparse bool
	// Resolved reports whether handle-based resolution succeeded (the path
	// exists and the OS gave us its final form).
	Resolved bool
}

// Resolve runs the fixed C26 pipeline on input. exceptions carries the
// user-configured [fs] reparse_point_exceptions (explicit paths only, never
// prefixes). Nonexistent paths resolve lexically and still classify; only a
// non-exempted reparse traversal is an error.
func Resolve(input string, exceptions []string) (Result, error) {
	p := expandInput(input)
	p = lexCanonical(p)      // the one sanctioned lexical step
	p = normalizeLocalUNC(p) // \\localhost\c$\... and \\?\UNC\localhost\c$\... -> c:\...
	p = lexCanonical(p)      // re-clean after prefix rewriting
	exems := make(map[string]bool, len(exceptions))
	for _, e := range exceptions {
		ee := normalizeLocalUNC(lexCanonical(expandInput(e)))
		exems[strings.ToLower(ee)] = true
	}

	res := Result{}
	if comps := reparseComponents(p); len(comps) > 0 {
		for _, c := range comps {
			if exems[strings.ToLower(c)] {
				res.Reparse = true
				continue
			}
			return res, ErrReparseDenied
		}
	}

	if final, ok := resolveHandle(p); ok {
		final = stripExtendedPrefix(final)
		final = normalizeLocalUNC(final)
		if final != "" {
			res.Canonical = final
			res.Resolved = true
			return res, nil
		}
	}
	// Path does not exist (or handle query unsupported): lexical fallback,
	// classification still enforced — fail-closed.
	res.Canonical = p
	return res, nil
}

// expandInput expands environment variables (%VAR% and $VAR) and a leading
// ~ (the user home). Unknown constructs pass through untouched.
func expandInput(input string) string {
	p := strings.TrimSpace(input)
	if p == "" {
		return p
	}
	// Windows-style %VAR%.
	for start := strings.Index(p, "%"); start >= 0; start = strings.Index(p[start+1:], "%") + start + 1 {
		end := strings.Index(p[start+1:], "%")
		if end < 0 {
			break
		}
		name := p[start+1 : start+1+end]
		if v := os.Getenv(name); v != "" && !strings.Contains(name, "%") {
			p = p[:start] + v + p[start+1+end+1:]
		}
	}
	p = os.ExpandEnv(p)
	// Leading ~ -> home. Only bare ~, ~/ or ~\ (never ~user).
	if p == "~" || strings.HasPrefix(p, `~\`) || strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil && home != "" {
			p = home + p[1:]
		}
	}
	return p
}

// lexCanonical is the single sanctioned filepath.Abs + filepath.Clean call
// site in the codebase (SPEC-06 §4 pipeline step 2-3).
func lexCanonical(p string) string {
	if p == "" {
		return p
	}
	if a, err := filepath.Abs(p); err == nil {
		p = a
	}
	return filepath.Clean(p)
}

// sepStr is the platform's path separator as a string. C26's canonical form is
// the real path of the real file (SPEC-06 §4), and on a POSIX system that shape
// uses '/'; '\' is merely a legal character inside a POSIX file name. Anything
// that hands a path to the OS (Lstat, Open, a re-joined ancestor chain) or to a
// human (an audit line) must build it with sepStr, never with a literal '\'.
const sepStr = string(filepath.Separator)

// unifySeparators folds the alternate separator into the platform one. On
// Windows the OS treats '/' and '\' as interchangeable, so both have to be
// folded or "C:/a" and "C:\a" compare unequal; on POSIX '\' is an ordinary
// filename character and folding it away would merge two different files into
// one comparison string, which for the B-tier override map is a cross-file
// unlock. The branch is a compile-time constant on each platform.
func unifySeparators(p string) string {
	if filepath.Separator == '\\' {
		return strings.ReplaceAll(p, "/", `\`)
	}
	return p
}

// normalizeLocalUNC maps UNC spellings of local drives back to drive paths:
// \\localhost\c$\x\y, \\127.0.0.1\c$\x\y and \\?\UNC\localhost\c$\x\y
// become c:\x\y. Other UNC paths are left as canonical \\server\share form.
//
// It normalizes UNC spellings and NOTHING ELSE: a path that is not shaped like
// a UNC path is returned untouched. Rewriting separators unconditionally used
// to live here (ticket 75), and on Linux that turned every absolute path into
// a backslash string which (a) is not absolute anymore, so the next
// lexCanonical re-anchored it on the process working directory, and (b) no OS
// call can open. It never mattered on Windows because filepath.Clean had
// already folded '/' to '\' two lines earlier, which is exactly why the defect
// survived every local gate.
func normalizeLocalUNC(p string) string {
	if !strings.HasPrefix(p, `\\`) && !strings.HasPrefix(p, `//`) {
		return p // not a UNC spelling: this function has no business here
	}
	u := strings.ReplaceAll(p, "/", `\`)
	for _, server := range []string{"localhost", "127.0.0.1"} {
		for _, prefix := range []string{`\\?\UNC\` + server + `\`, `\\\?\UNC\` + server + `\`, `\\` + server + `\`} {
			if len(u) > len(prefix)+2 && strings.HasPrefix(strings.ToLower(u), strings.ToLower(prefix)) {
				share := u[len(prefix):]
				i := strings.Index(share, `\`)
				if i < 0 {
					continue
				}
				drive, rest := share[:i], share[i+1:]
				if len(drive) == 2 && drive[1] == '$' && drive[0] >= 'a' && drive[0] <= 'z' || len(drive) == 2 && drive[1] == '$' && drive[0] >= 'A' && drive[0] <= 'Z' {
					return strings.ToUpper(drive[:1]) + `:\` + rest
				}
			}
		}
	}
	return u
}

// stripExtendedPrefix removes the \\?\ (and \\?\UNC\) prefix from a
// GetFinalPathNameByHandle result, yielding the canonical user-visible form.
func stripExtendedPrefix(p string) string {
	if strings.HasPrefix(p, `\\?\UNC\`) {
		return `\\` + p[len(`\\?\UNC\`):]
	}
	if strings.HasPrefix(p, `\\?\`) {
		return p[len(`\\?\`):]
	}
	return p
}

// ---------------------------------------------------------------------------
// Comparison forms: the ANCHOR side of a security comparison (ticket 72).
// ---------------------------------------------------------------------------
//
// Resolve guarantees that a candidate path is a handle-resolved real path. It
// says nothing about the anchors of the A/B tables, which are built from the
// environment (USERPROFILE / APPDATA / LOCALAPPDATA) and therefore keep
// whatever spelling the process was started with. On GitHub's windows-latest
// runner USERPROFILE is the 8.3 short form (C:\Users\RUNNER~1) while
// GetFinalPathNameByHandle hands back the long one (C:\Users\runneradmin), so
// "canonical is under ~/.ssh" compared raw-vs-resolved misses and an A-tier
// path silently degrades to B — one L2 confirm away from being allowed.
//
// SPEC-06 §4 puts "展开 8.3 短名" inside the resolution pipeline and PLAN.md
// C26 requires the handle's real path, so the anchor side owes the same
// treatment: expand both sides through this pipeline and only then compare.
// Where a side cannot be PROVEN expanded, a miss is not evidence of absence,
// and the only permitted answer is fail-closed (see pathForms.certain).

// pathForms are the spellings one path can take in a security comparison.
type pathForms struct {
	// raw is the spelling as given, comparison-folded. Always populated
	// (except for an empty input) and always compared, so no verdict can be
	// lost by moving to real paths.
	raw string
	// real is raw with every openable component replaced by its handle
	// final path; components below the deepest openable one are re-appended
	// verbatim (they do not exist yet). Empty when nothing could be opened:
	// a non-Windows build (see pathresolver_other.go), a dead volume, or a
	// malformed spelling.
	real string
	// root is the handle-resolved prefix real was built on ("" when nothing
	// opened). Two sides with roots that are not ancestry-related cannot hide
	// a hit from each other, which bounds the fail-closed rule below.
	root string
	// partial marks an untrustworthy real: the walk stopped below a component
	// that exists but refused to open (deny-filtered directory, dangling
	// reparse point), so that component's true spelling is unknown.
	partial bool
}

// certain reports whether this side is proven to be a handle-expanded real
// path. R17: an A-anchor comparison may only conclude "no hit" when both
// sides are certain.
func (f pathForms) certain() bool { return f.real != "" && !f.partial }

// spellings lists the forms this path must be matched on. The list is a union
// on purpose: matching an A anchor on any form can only ADD a deny, never
// remove one, so the change is monotone toward fail-closed.
func (f pathForms) spellings() []string {
	if f.real == "" || f.real == f.raw {
		return []string{f.raw}
	}
	return []string{f.raw, f.real}
}

// maxAnchorWalk bounds the ancestor walk so a pathological spelling (or a
// filesystem that keeps refusing to open) cannot make Classify spin.
const maxAnchorWalk = 64

// formsOf expands one spelling into its comparison forms using the C26 handle
// pipeline only: no string heuristics, no case or short-name "alignment".
func formsOf(spelling string) pathForms {
	raw := normPath(spelling)
	f := pathForms{raw: raw}
	if raw == "" {
		return f
	}
	p := lexCanonical(spelling)
	cut := p
	for depth := 0; depth < maxAnchorWalk; depth++ {
		if final, ok := resolveHandle(cut); ok {
			final = normalizeLocalUNC(stripExtendedPrefix(final))
			if final != "" {
				f.root = normPath(final)
				tail := strings.TrimPrefix(p, cut) // cut is a literal prefix of p
				if tail == "" {
					f.real = f.root
				} else {
					f.real = normPath(lexCanonical(final + tail))
				}
				f.partial = tailExistsBelow(cut, tail)
				return f
			}
		}
		parent := filepath.Dir(cut)
		if parent == cut {
			break
		}
		cut = parent
	}
	return f
}

// anchorCache memoizes anchor forms. Env anchors are a handful of spellings
// per process and they are read on every classification, so without this the
// gate would pay a CreateFile + GetFinalPathNameByHandle per rule per call.
// Keys are folded spellings; a cached entry stays correct when a protected
// leaf is created later, because its not-yet-existing tail components are
// literal names re-appended onto an already-expanded prefix.
var (
	anchorCache    sync.Map // folded spelling -> pathForms
	anchorCacheLen atomic.Int64
)

const anchorCacheMax = 128

// anchorForms is formsOf for an anchor spelling, with the cache above.
func anchorForms(spelling string) pathForms {
	key := normPath(spelling)
	if key == "" {
		return pathForms{}
	}
	if v, ok := anchorCache.Load(key); ok {
		return v.(pathForms)
	}
	f := formsOf(spelling)
	if anchorCacheLen.Load() < anchorCacheMax {
		if _, loaded := anchorCache.LoadOrStore(key, f); !loaded {
			anchorCacheLen.Add(1)
		}
	}
	return f
}

// tailExistsBelow reports whether the first re-appended component exists.
// resolveHandle already refused every prefix of the tail, so an existing
// first component means the walk skipped something that is really there and
// whose true spelling we therefore do not know (partial). If it does not
// exist, nothing below it can exist either, so the verbatim tail is exact.
func tailExistsBelow(cut, tail string) bool {
	if tail == "" {
		return false
	}
	first := strings.TrimPrefix(tail, sepStr)
	if i := strings.Index(first, sepStr); i >= 0 {
		first = first[:i]
	}
	if first == "" {
		return false
	}
	return pathExists(cut + sepStr + first)
}

// pathExists is a fail-toward-exists stat: only a clean "no such file" reads
// as absent, any other error (access denied, sharing violation, bad path)
// leaves the component's spelling unproven.
func pathExists(p string) bool {
	_, err := os.Lstat(p)
	if err == nil {
		return true
	}
	return !errors.Is(err, os.ErrNotExist)
}
