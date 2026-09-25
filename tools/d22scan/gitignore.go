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
// reason nobody can see. `.gitignore` IS the policy this tool reads.
//
// WHAT THAT LEAVES OUT, AND WHAT THE r1 ACCEPTANCE FOUND (ledger A214, F1). Git's
// criterion for "is this path ignored?" is the rule file **AND the index**: a path
// the index already holds is never ignored, which is exactly what `git add -f`
// means - someone declared that file part of the deliverable, and git honours it
// (measured on this machine, `git check-ignore -v <force-added>` exits 1 while the
// same command with `--no-index` exits 0 naming the rule). A pattern-only matcher
// therefore skips files git does not skip, and every one of those is a shipped
// byte the gate stops reading. `skip()` asks both questions now: the rule files in
// this tree, and git's own instruments for what this tree's index holds
// (askGitIndex). The claim this file makes - and the one
// scan_test.go's TestTrackedPathsAreNeverSkippedByTheIgnoreFilter pins by force
// adding a file under an ignored path and demanding it go red - is the narrow,
// testable one: **for a path git reports as tracked, skip() returns false.**
//
// WHY git IS CALLED, and what happens when it cannot be. `git check-ignore` was
// not used in the r1 batch because the snapshot shape this repo's evidence and
// every mutation check runs on is `git archive HEAD | tar -x`, which has NO `.git`
// (measured in d22scan-gitignore-accept-r1.md §4), and gating the filter on the
// presence of `.git` is not a fix either: a real CI checkout is a clone WITH
// `.git`, so that gate leaves the dangerous case - a tracked file matching a rule,
// in the tree CI and every developer actually scans - filtered out and only
// re-blinds the archive copy (A214 ②, which is why the named minimal fix was not
// adopted). This batch does the opposite instead: git is asked, and when it cannot
// answer - no git binary, not a repository, `-root` pointing at a subdirectory of
// somebody else's repository, a failing command, an expired deadline, output that
// does not parse - **no ignore rule is applied at all** and the self-report says
// so loudly (see note()). The failure direction is always MORE scanning, never a
// quiet green: a scanner that could not ask does not get to claim a clean count.
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
//   - a TRAILING "**" matches everything inside its directory and never the
//     directory itself (`dist/**` does not ignore `dist`), measured against
//     `git check-ignore -v dist` -> exit 1. Pruning `dist` would swallow the
//     `!dist/.gitkeep` re-inclusion below it and take CI's 40 down to 39, which is
//     the shape gitignore.go's own r1 header warned about and did not implement
//     (acceptance §2 probe P2, F3 ②). A "**" in the MIDDLE of a pattern still
//     matches zero directories, because `a/**/b` matching `a/b` is git's documented
//     behaviour (`git check-ignore -v a/b` -> exit 0) and was already pinned.
//   - a rule line's LEADING whitespace is part of the pattern; only trailing
//     unquoted whitespace and one trailing CR are removed. `git check-ignore -v`
//     on a file under `dist/` with the rule `  dist/*` exits 1 (that rule names a
//     directory whose name starts with two spaces), while a directory really called
//     `  dist` is ignored by it. Trimming the front, as r1 did, turns a rule that
//     git binds to nothing into a rule that swallows a live tree (F3 ①, probe P1).
//   - an ignored directory takes its contents with it and no negation below it
//     can bring them back (git's documented parent-directory rule).
//   - a path the index holds is never skipped, whatever the rules say, and a
//     directory containing one is never pruned (F1).
// UNSUPPORTED ON PURPOSE: `.git/info/exclude` and `core.excludesFile` as SOURCES
// OF A SKIP (they live outside the scanned tree, so honouring them would make the
// count depend on the machine again - A207's disease), backslash-escaped trailing
// spaces, and comments after a pattern. Each of those can only ever make this
// matcher skip LESS than git does, which is the loud direction: the path stays
// scanned and the worst outcome is a denominator that counts one file too many.
// The one shape that could make it skip MORE than git - not being able to reach
// the index at all - is handled by the rule above it, not by this list: no rule
// is applied then, so the tool over-scans and says why.

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
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
// root, using the .gitignore files inside that tree AND that tree's git index -
// the two halves of git's own criterion. The index half is asked lazily (see
// index), and if it cannot be asked, no rule is applied at all.
type gitIgnore struct {
	root   string
	rules  map[string][]ignoreRule // rule file dir (slash, "" = root) -> its rules
	probed map[string]bool         // rule file dir already looked for
	cache  map[string]ignoreVerdict
	idx    *gitIndexState // git's answer about tracked paths, nil until first asked
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

// gitIndexState is git's own answer about the tree under root: which paths the
// index already holds. It is what turns "the rule matched" into "git would decline
// to track this path", and it is the only thing that can prove the difference -
// no in-tree file states whether a path is tracked.
//
// Three read-only questions are asked (runGitIndex), and they are all-or-nothing:
// a half-loaded index would be a new blindfold, so any failure yields
// ok=false plus one reason, and skip() then applies NO ignore rules at all.
//
//	rev-parse --show-prefix                       is root a repository top?
//	ls-files -z                                   every tracked path (the guard)
//	ls-files -z -i -c --exclude-standard          tracked paths that also match
//	                        an ignore rule (the ticket's named instrument, the set
//	                        the self-report names, and an independent union member
//	                        so one command's parsing bug cannot blind the other)
type gitIndexState struct {
	ok        bool
	why       string          // why it could not be consulted; one line, always shown
	tracked   map[string]bool // every path in the index, slash-separated, relative to root
	dirs      map[string]bool // every ancestor directory of a tracked path
	ignored   map[string]bool // `ls-files -i -c --exclude-standard`
	ignoreSet []string        // ignored, sorted (the note must not depend on map order)
}

// holds is the whole of F1's fix in one predicate: this path, or something inside
// this directory, is in the index, so it is part of the delivery and a pattern may
// not hide it.
func (ix *gitIndexState) holds(rel string, isDir bool) bool {
	if ix == nil || !ix.ok {
		return false
	}
	if isDir {
		return ix.dirs[rel]
	}
	return ix.tracked[rel] || ix.ignored[rel]
}

// indexTimeout bounds one git call. A gate that hangs is not a gate. The deadline
// is a context timer (monotonic), never a wall-clock delta, which AGENTS.md §1.2
// bans as a shape - so nothing here compares two time.Now() readings.
const indexTimeout = 10 * time.Second

// runGitIndex asks the three questions above. Every error path collapses to
// ok=false with a single-line reason, because the caller's only legal response is
// "apply no rules and say so".
func runGitIndex(root string) *gitIndexState {
	fail := func(format string, args ...any) *gitIndexState {
		return &gitIndexState{why: fmt.Sprintf(format, args...)}
	}
	prefix, err := runGitAt(root, "rev-parse", "--show-prefix")
	if err != nil {
		return fail("git cannot be consulted in %s: %s", filepath.ToSlash(root), err)
	}
	if p := strings.TrimRight(prefix, "\r\n"); p != "" {
		// root is a SUBDIRECTORY of somebody else's repository. Its index governs
		// more than this tree and lists paths against the repository top, so the
		// answer cannot be resolved against root-relative rules - treat it as
		// unreachable rather than half-trust it.
		return fail("%s is the subdirectory %q of a git repository, not its top, so that index is not this tree's index", filepath.ToSlash(root), filepath.ToSlash(p))
	}
	all, err := runGitAt(root, "ls-files", "-z")
	if err != nil {
		return fail("git ls-files failed in %s: %s", filepath.ToSlash(root), err)
	}
	withIndex, err := runGitAt(root, "ls-files", "-z", "-i", "-c", "--exclude-standard")
	if err != nil {
		return fail("git ls-files -i -c --exclude-standard failed in %s: %s", filepath.ToSlash(root), err)
	}
	tracked, dirs, err := parseIndexPaths(all)
	if err != nil {
		return fail("cannot parse git ls-files output in %s: %s", filepath.ToSlash(root), err)
	}
	ignored, err := parseIndexPathsList(withIndex)
	if err != nil {
		return fail("cannot parse git ls-files -i -c output in %s: %s", filepath.ToSlash(root), err)
	}
	set := map[string]bool{}
	for _, p := range ignored {
		set[p] = true
	}
	sort.Strings(ignored)
	return &gitIndexState{ok: true, tracked: tracked, dirs: dirs, ignored: set, ignoreSet: ignored}
}

// runGitAt runs one read-only git command with root as the working directory and
// returns its stdout. Errors are collapsed to one line because they end up in the
// self-report.
func runGitAt(root string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), indexTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = root
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return "", fmt.Errorf("timed out after %s running `git %s`", indexTimeout, strings.Join(args, " "))
		}
		reason := firstLine(stderr.String())
		if reason == "" {
			reason = err.Error()
		}
		return "", fmt.Errorf("`git %s`: %s", strings.Join(args, " "), reason)
	}
	return stdout.String(), nil
}

func firstLine(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	for _, line := range strings.Split(s, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			return line
		}
	}
	return ""
}

// parseIndexPaths reads a NUL-separated `git ls-files` stream into the tracked set
// plus the set of every directory above it (a tracked file anywhere below a
// directory means that directory may never be pruned).
func parseIndexPaths(z string) (map[string]bool, map[string]bool, error) {
	list, err := parseIndexPathsList(z)
	if err != nil {
		return nil, nil, err
	}
	tracked := make(map[string]bool, len(list))
	dirs := map[string]bool{}
	for _, p := range list {
		tracked[p] = true
		acc := ""
		segs := strings.Split(p, "/")
		for i := 0; i < len(segs)-1; i++ {
			if acc == "" {
				acc = segs[i]
			} else {
				acc += "/" + segs[i]
			}
			dirs[acc] = true
		}
	}
	return tracked, dirs, nil
}

func parseIndexPathsList(z string) ([]string, error) {
	var out []string
	for _, p := range strings.Split(z, "\x00") {
		if p == "" {
			continue
		}
		// Nothing git prints here may be a path outside the tree; if it is, this
		// matcher is reading a stream it does not understand and must stop trusting
		// it rather than resolve it into a skip decision.
		if strings.HasPrefix(p, "/") {
			return nil, fmt.Errorf("absolute path %q in git output", p)
		}
		for _, seg := range strings.Split(p, "/") {
			if seg == ".." {
				return nil, fmt.Errorf("path %q in git output reaches above the root", p)
			}
		}
		out = append(out, p)
	}
	return out, nil
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
		// the walk before it examined anything. Asked BEFORE the index, so this
		// branch costs no git call and cannot be reached by a failed one.
		return false
	}
	ix := g.index()
	if !ix.ok {
		// No rule is applied when git could not be asked whether the path is
		// tracked: the only safe answer to "would git decline to track this?" is
		// then "unknown", and unknown must not become a skip.
		return false
	}
	if ix.holds(rel, isDir) {
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

// index loads git's answer once per matcher and keeps it: the walks are
// sequential and four of them share one matcher, so re-asking per path would be
// thousands of subprocess calls, and re-asking mid-walk would let a path skipped
// early disagree with the same path seen again under another ban.
func (g *gitIgnore) index() *gitIndexState {
	if g.idx == nil {
		g.idx = runGitIndex(g.root)
	}
	return g.idx
}

// note renders the self-report lines about what this matcher did and did not get
// to decide, "" when there is nothing to say. It is deliberately NOT part of the
// scanScope ledger: guard 0 (undeclaredKeys) makes a counter that no scope reports
// fatal, and this is not coverage - it is provenance for a coverage number.
//
// Three lines are possible, and the first two are the ones that keep this tool
// honest after A214: an index it could not read must be announced even when the
// run is otherwise green (a quiet "clean" from a scanner that could not ask is the
// false green this repository keeps catching), and a tracked path sitting under an
// ignore rule must be named, because that is precisely the file the r1 batch made
// invisible.
func (g *gitIgnore) note() string {
	if g == nil {
		return ""
	}
	var lines []string
	if g.idx != nil && !g.idx.ok {
		lines = append(lines, fmt.Sprintf("d22scan: gitignore rules NOT APPLIED - %s; every path in every scope is being scanned, "+
			"so the counts below may include build output (A207's machine-dependent denominator). This is the loud direction: "+
			"a scanner that cannot ask git which paths are tracked does not get to skip any.", g.idx.why))
	}
	if g.files > 0 || len(g.pruned) > 0 {
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
		lines = append(lines, fmt.Sprintf("d22scan: skipped as git-ignored: %d file(s) under %d ignored director(ies) [%s], decided by %s",
			g.files, len(g.pruned), strings.Join(dirs, ", "), strings.Join(src, ", ")))
	}
	if g.idx != nil && len(g.idx.ignoreSet) > 0 {
		lines = append(lines, fmt.Sprintf("d22scan: %d path(s) git reports as TRACKED and matching an ignore rule (git ls-files -i -c --exclude-standard) "+
			"- a tracked path is part of the delivery, so none of them was skipped: %s",
			len(g.idx.ignoreSet), strings.Join(g.idx.ignoreSet, ", ")))
	}
	return strings.Join(lines, "\n")
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
		// Strip a trailing CR and then trailing unquoted whitespace ONLY. Git keeps
		// leading whitespace as pattern data (measured: the rule `  dist/*` ignores
		// a directory whose name really begins with two spaces and binds nothing
		// else), so TrimSpace here silently widens a rule that git binds to nothing
		// into one that swallows a live tree - A214/F3 ①, probe P1, where
		// `  dist/*` swallowed dist/index.html.
		line := strings.TrimRight(raw, " \t\r")
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
		// A `**` in the middle of a pattern may match zero directories (`a/**/b`
		// matches `a/b` - pinned by TestGitIgnoreRuleSemantics and confirmed against
		// git), but a `**` at the END of one must consume at least one segment:
		// `dist/**` matches everything inside dist and never dist itself, so it
		// cannot prune the directory that the `!dist/.gitkeep` negation lives under
		// (A214/F3 ②, probe P2, where it took CI's 40 to 39).
		start := 0
		if len(pat) == 1 {
			start = 1
		}
		for i := start; i <= len(name); i++ {
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
