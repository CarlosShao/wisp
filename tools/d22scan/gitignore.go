package main

// Gitignore-aware path skipping for every d22scan walk (ledger A207, the
// 2026-09-25 fix).
//
// WHY THIS FILE EXISTS. The per-scope numbers this tool prints are the
// instrument's own readings: "ban #6 frontend/ examined 40 text files" is what
// makes "0 findings" mean anything (ticket 67 AC#2, ticket 71 AC#4). Those
// denominators were computed from the FILESYSTEM, so they counted whatever the
// machine happened to hold - and `frontend/dist/` is build output, ignored by
// `frontend/.gitignore:12` (`dist/*`) and `.gitignore:24` (`frontend/dist/*`),
// with `dist/.gitkeep` re-included by a negation in both files because
// go:embed needs the directory to exist in a clean checkout (ticket 77 AC#1).
// Measured at anchor 7771409 with one binary on one source tree:
//
//	real worktree (someone ran a build)      ban #6 / ban #8 frontend/ = 43
//	`git archive HEAD` snapshot (no build)   ban #6 / ban #8 frontend/ = 40
//
// The 3-path gap is `frontend/dist/index.html` plus the two hashed bundles under
// `frontend/dist/assets/`. That is the whole explanation for the ledger citing
// frontend/=37, 40 and 43 for one scope (A207 ①, A209 ②): the ruler was standing
// on moving ground, so a real coverage change and someone running `npm run
// build` were indistinguishable. Nothing an ignored path contains can ever be
// part of what a verdict is about - CI's checkout has no ignored files, so a ban
// that "examines" one is counting work it will not do there.
//
// WHY THE RULE IS READ FROM .gitignore RATHER THAN BEING HARD-CODED HERE. This
// repository keeps getting burned by one rule living in three copies where only
// one gets updated (A207 ②/`U1`/`U2` for the emoji character class is the same
// disease). Writing a second literal list of {"dist","node_modules"} would
// re-open it the day somebody ignores a new directory inside a scanned scope:
// the scanner would keep counting it and the baseline would move again for a
// reason nobody can see. `.gitignore` IS the policy; the scanner follows it, so
// the instrument and the repository's ignore rules cannot disagree.
//
// WHY `git check-ignore` IS NOT USED, even though it is the authoritative
// implementation. The snapshot shape CI's evidence and every mutation check
// here runs on is `git archive HEAD | tar -x`, which has NO `.git` directory
// (measured: `ls <snap>/.git` -> "No such file or directory"), so a git call
// would either fail or silently return "nothing is ignored" - making the reading
// depend on whether the tree is a checkout, which is exactly the defect being
// fixed. A lint gate must not require a git binary to run.
//
// SEMANTICS IMPLEMENTED (the gitignore(5) subset that can bind a path inside
// design/ frontend/ internal/ cmd/):
//   - one rule file per directory: <root>/.gitignore plus a nested .gitignore in
//     any ancestor of the path; the DEEPEST file holding a matching rule decides,
//     and inside one file the LAST matching line wins. The deepest-file rule is
//     load-bearing for the counts, not tidiness: it is what makes
//     `!frontend/dist/.gitkeep` survive, and .gitkeep is one of the 40 files CI
//     reports, so an implementation that pruned the ignored `dist` directory
//     wholesale would move CI's number too (40 -> 39) and break the "no behavior
//     change on CI's shape" requirement.
//   - trailing "/" means directory-only; a pattern carrying a "/" anywhere but
//     at the end is anchored to its own rule file's directory, otherwise it
//     matches the base name at any depth.
//   - "*" and "?" do not cross "/" (git matches with FNM_PATHNAME); "**" does.
//     "[...]" classes support a leading "!" or "^" negation and "a-z" ranges.
//   - an ignored directory takes its contents with it and no negation below it
//     can bring them back (git's documented parent-directory rule).
// UNSUPPORTED ON PURPOSE: `.git/info/exclude` and `core.excludesFile` (they live
// outside the scanned tree, so honouring them would make the count depend on the
// machine again), backslash-escaped trailing spaces, and comments after a
// pattern. Every unsupported shape fails in the loud direction: the path stays
// scanned, so the worst this matcher can do is examine too much, never less.

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"
)

// ignoreRule is one .gitignore line, pre-split for matching.
type ignoreRule struct {
	segs     []string // pattern split on '/'; the directory-only marker removed
	dirSegs  int      // leading path segments consumed by the rule file's own directory
	negate   bool     // "!" - re-include
	dirOnly  bool     // trailing "/"
	anchored bool     // matched against the path below the rule file, not the base name
}

// ignoreVerdict is a memoised decision plus the rule file that produced it, kept
// so the self-report can name WHERE an exclusion came from (a number that moved
// with no traceable source is the failure this whole file exists to end).
type ignoreVerdict struct {
	ignored bool
	origin  string // repo-relative rule file, "" when nothing matched
}

// gitIgnore answers "would git decline to track this path?" for the tree under
// root, using only the .gitignore files inside that tree.
type gitIgnore struct {
	root   string
	rules  map[string][]ignoreRule // rule file dir (slash, "" = root) -> its rules
	probed map[string]bool         // rule file dir already looked for
	cache  map[string]ignoreVerdict
	// The three counters below describe the TREE, not the number of walk visits:
	// four walks share one matcher (walkGo, walkText, walkEmoji x4 scopes), so
	// without `counted` a file sitting in a scope two bans read - which is exactly
	// frontend/, the tree this fix is about - would be reported twice, and a note
	// that over-counts is as untrustworthy as one that omits.
	counted map[string]bool // path -> already accounted for in the note
	files   int             // excluded files
	pruned  map[string]bool // excluded directory -> pruned
	origins map[string]int  // rule file -> how many distinct paths it excluded
}

func newGitIgnore(root string) *gitIgnore {
	return &gitIgnore{
		root:    root,
		rules:   map[string][]ignoreRule{},
		probed:  map[string]bool{},
		cache:   map[string]ignoreVerdict{},
		counted: map[string]bool{},
		pruned:  map[string]bool{},
		origins: map[string]int{},
	}
}

// skip is the walks' entry point: true means "do not count and do not read this
// path". A directory returning true is a filepath.SkipDir, a file returning true
// is a plain return nil.
func (g *gitIgnore) skip(path string, isDir bool) bool {
	if g == nil {
		return false
	}
	rel := g.rel(path)
	if rel == "" {
		// The walk's own starting directory (rel == ""), and anything the root
		// prefix did not survive: never skipped, because skipping it would abort
		// the walk before it examined anything.
		return false
	}
	v := g.decide(rel, isDir)
	if !v.ignored {
		return false
	}
	if !g.counted[rel] {
		g.counted[rel] = true
		if isDir {
			g.pruned[rel] = true
		} else {
			g.files++
		}
		if v.origin != "" {
			g.origins[v.origin]++
		}
	}
	return true
}

// note renders one self-report line about what the matcher excluded, or "" when
// it excluded nothing. It is deliberately NOT part of the scanScope ledger:
// guard 0 (undeclaredKeys) makes a counter that no scope reports fatal, and this
// is not coverage - it is provenance for a coverage number.
func (g *gitIgnore) note() string {
	if g == nil || (g.files == 0 && len(g.pruned) == 0) {
		return ""
	}
	dirs := make([]string, 0, len(g.pruned))
	for d := range g.pruned {
		dirs = append(dirs, d+"/")
	}
	sort.Strings(dirs)
	src := make([]string, 0, len(g.origins))
	for o, n := range g.origins {
		src = append(src, fmt.Sprintf("%s (%d path(s))", filepath.ToSlash(o), n))
	}
	sort.Strings(src)
	return fmt.Sprintf("d22scan: skipped as git-ignored: %d file(s) under %d ignored director(ies) [%s], decided by %s",
		g.files, len(g.pruned), strings.Join(dirs, ", "), strings.Join(src, ", "))
}

// rel turns a walk path into a slash path relative to root ("" when it is root).
func (g *gitIgnore) rel(path string) string {
	rel, err := filepath.Rel(g.root, path)
	if err != nil {
		return ""
	}
	rel = filepath.ToSlash(rel)
	if rel == "." {
		return ""
	}
	return rel
}

func (g *gitIgnore) decide(rel string, isDir bool) ignoreVerdict {
	key := rel
	if isDir {
		key += "/"
	}
	if v, ok := g.cache[key]; ok {
		return v
	}
	v := g.match(rel, isDir)
	g.cache[key] = v
	return v
}

// match applies the rule files. Parent directories are consulted first because
// git prunes an excluded directory and never looks inside it, so a negation
// below an ignored dir cannot re-include anything.
func (g *gitIgnore) match(rel string, isDir bool) ignoreVerdict {
	segs := strings.Split(rel, "/")
	acc := ""
	for i := 0; i < len(segs)-1; i++ {
		if acc == "" {
			acc = segs[i]
		} else {
			acc += "/" + segs[i]
		}
		if v := g.ruleMatch(acc, true); v.ignored {
			// The path is inside an excluded directory; report that directory's
			// rule file as the origin so the note stays traceable.
			return ignoreVerdict{ignored: true, origin: v.origin}
		}
	}
	return g.ruleMatch(rel, isDir)
}

// ruleMatch consults the rule files from the deepest to the root and returns the
// decision of the first file that has a matching line.
func (g *gitIgnore) ruleMatch(rel string, isDir bool) ignoreVerdict {
	segs := strings.Split(rel, "/")
	name := segs[len(segs)-1]
	for d := len(segs) - 1; d >= 0; d-- {
		dir := strings.Join(segs[:d], "/")
		rules := g.rulesFor(dir)
		if len(rules) == 0 {
			continue
		}
		ignored, matched := false, false
		for _, r := range rules {
			if r.dirOnly && !isDir {
				continue
			}
			target := name
			if r.anchored {
				if len(segs) <= r.dirSegs {
					continue // the pattern needs a path below its own directory
				}
				target = strings.Join(segs[r.dirSegs:], "/")
			}
			if pathMatch(r.segs, strings.Split(target, "/")) {
				matched, ignored = true, !r.negate
			}
		}
		if matched {
			origin := ".gitignore"
			if dir != "" {
				origin = dir + "/.gitignore"
			}
			return ignoreVerdict{ignored: ignored, origin: origin}
		}
	}
	return ignoreVerdict{}
}

// rulesFor parses <root>/<dir>/.gitignore once per directory.
func (g *gitIgnore) rulesFor(dir string) []ignoreRule {
	if rs, ok := g.rules[dir]; ok {
		return rs
	}
	if g.probed[dir] {
		return nil
	}
	g.probed[dir] = true
	path := filepath.Join(g.root, filepath.FromSlash(dir), ".gitignore")
	data, err := os.ReadFile(path)
	if err != nil {
		// Absent is the normal case. Unreadable is treated the same way, which
		// means NO rules and therefore no exclusion: the path stays scanned.
		g.rules[dir] = nil
		return nil
	}
	rs := parseIgnoreFile(string(data), strings.Count(dir, "/")+lenNonEmpty(dir))
	g.rules[dir] = rs
	return rs
}

// lenNonEmpty is 1 for a non-empty directory name and 0 for "" (the root), so
// dirSegs counts the segments a relative path must drop before matching.
func lenNonEmpty(dir string) int {
	if dir == "" {
		return 0
	}
	return 1
}

func parseIgnoreFile(src string, dirSegs int) []ignoreRule {
	var out []ignoreRule
	for _, raw := range strings.Split(src, "\n") {
		line := strings.TrimRight(strings.TrimSpace(raw), " \t")
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		r := ignoreRule{dirSegs: dirSegs}
		if strings.HasPrefix(line, "!") {
			r.negate = true
			line = line[1:]
		}
		if strings.HasSuffix(line, "/") {
			r.dirOnly = true
			line = strings.TrimSuffix(line, "/")
		}
		if strings.HasPrefix(line, "/") {
			r.anchored = true
			line = line[1:]
		} else if strings.Contains(line, "/") {
			r.anchored = true
		}
		if line == "" {
			continue
		}
		r.segs = strings.Split(line, "/")
		out = append(out, r)
	}
	return out
}

// pathMatch matches a slash path against a pattern, segment by segment: '*' and
// '?' stop at '/', '**' crosses them.
func pathMatch(pat, name []string) bool {
	if len(pat) == 0 {
		return len(name) == 0
	}
	if pat[0] == "**" {
		for i := 0; i <= len(name); i++ {
			if pathMatch(pat[1:], name[i:]) {
				return true
			}
		}
		return false
	}
	if len(name) == 0 {
		return false
	}
	if !segMatch(pat[0], name[0]) {
		return false
	}
	return pathMatch(pat[1:], name[1:])
}

// segMatch is fnmatch without FNM_PATHNAME for one path segment: '*' (runs of
// asterisks count as one, matching git's "two or more adjacent asterisks are
// treated as a single '*'), '?', "[class]" and backslash escapes.
func segMatch(pat, s string) bool {
	for len(pat) > 0 {
		switch pat[0] {
		case '*':
			pat = strings.TrimLeft(pat, "*")
			if pat == "" {
				return true
			}
			for i := 0; i <= len(s); i++ {
				if segMatch(pat, s[i:]) {
					return true
				}
			}
			return false
		case '?':
			if s == "" {
				return false
			}
			_, size := utf8.DecodeRuneInString(s)
			pat, s = pat[1:], s[size:]
		case '[':
			end := strings.Index(pat[1:], "]")
			if end < 0 || s == "" {
				// An unterminated class matches a literal '[' (git treats a
				// malformed pattern's bracket as data).
				if s == "" || s[0] != '[' {
					return false
				}
				pat, s = pat[1:], s[1:]
				continue
			}
			if !classMatch(pat[1:1+end], s[0]) {
				return false
			}
			pat, s = pat[end+2:], s[1:]
		case '\\':
			if len(pat) < 2 || s == "" || s[0] != pat[1] {
				return false
			}
			pat, s = pat[2:], s[1:]
		default:
			if s == "" || s[0] != pat[0] {
				return false
			}
			pat, s = pat[1:], s[1:]
		}
	}
	return s == ""
}

func classMatch(class string, c byte) bool {
	negate := false
	if strings.HasPrefix(class, "!") || strings.HasPrefix(class, "^") {
		negate = true
		class = class[1:]
	}
	hit := false
	for i := 0; i < len(class); i++ {
		if i+2 < len(class) && class[i+1] == '-' {
			if class[i] <= c && c <= class[i+2] {
				hit = true
			}
			i += 2
			continue
		}
		if class[i] == c {
			hit = true
		}
	}
	return hit != negate
}
