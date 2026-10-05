package main

// This file is ticket 161 AC#2: the gates' own gate.
//
// WHY IT EXISTS. Every ban in main.go is a rule about somebody else's code, and
// until now the only instrument that could tell whether a ban still SEES what its
// own documentation claims it sees was a human reading the regex. 161-r1 measured
// that question for real (docs/evidence/s1/161-gate-blindspot-r1.md §3): all eight
// numbered bans do ring today, and six narrower escapes plus three temperaments
// are registered in that file's §4 / §4.1. Those readings were produced by a
// throwaway shell bench (`probes/161/r1/bench.sh`) that CI never runs, so the
// numbers expired the moment they were written - "the gate is not blind" was a
// claim with zero instruments behind it, which is the exact disease ticket 161 was
// opened for. This file turns that bench into a rule: one table, one pair of
// samples per ban (one that MUST ring, one that MUST stay silent), executed by
// `d22scan -self-test` and by `go test ./...` in this module.
//
// WHAT IT DELIBERATELY DOES NOT TOUCH. Not one ban's scope, band, regex or
// comment-exemption is read from here or modified by it: the matcher logic in
// main.go is frozen by Q-46 / ticket 141 (a contract face, owner-approved only).
// This file only ever FEEDS files to the existing walks and reports what came
// back, so "make the gate match more" cannot be smuggled in as "make the gate
// testable". Where a case's expectation is uncomfortable - ban #6 firing on a
// comment, ban #8 staying quiet on an arrow - the case pins today's behaviour and
// says so in its note, rather than asserting what someone wishes the ban did.
// Flipping any of those pins is a scope change and must read as one.
//
// WHY THE EXPECT-SILENT SIDE NEEDS A COVERAGE COUNTER. "No finding" is what a
// deleted walk, a mistyped tag and a genuinely clean file all print. So every
// case also reads the counter the walk that owns that tag bumped, on the tree the
// case just built, and a silent expectation backed by a zero counter is reported
// as VACUOUS and fails the run (ticket 71 AC#4's rule, applied to the self-test
// itself instead of only to the scan).
//
// THE ROSTER IS READ, NOT COPIED. Which tags must have a pair is computed at
// run time from this tool's own source (the `// Bans` block, and every
// `s.add("<tag>"` literal it contains). A ban added to main.go without a pair
// here, a pair written for a tag main.go no longer emits, or a pair with one
// direction deleted all exit 2 - which is what makes the reverse criterion
// (ticket 161 AC#4) answerable with a reading instead of a sentence: removing one
// sample from this table turns the CI step red, it does not turn it quiet.
//
// Invocation (the source-file argument is resolved from the working directory,
// same constraint as the rest of this module - see main.go's doc comment):
//
//	cd tools/d22scan && go run . -self-test
//	bash tools/d22scan/runtests.sh -C tools/d22scan ./...   # runs it in-process too

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// selfWant is the single assertion every self-test case makes about one tag.
type selfWant int

const (
	// wantRing: the violating sample must produce at least one finding with this tag.
	wantRing selfWant = iota
	// wantSilent: the clean (or deliberately pinned) sample must produce none.
	wantSilent
)

func (w selfWant) String() string {
	if w == wantRing {
		return "ring"
	}
	return "silent"
}

// selfCase is one sample file plus the one thing its scan result must show.
//
// Cases come in pairs per tag by construction: auditSelfCases() refuses to let
// the table run unless every tag main.go can emit has at least one wantRing case
// AND at least one wantSilent case, so "the pair" is a property of the machine,
// not of my diligence.
type selfCase struct {
	tag  string   // the finding tag this case is about, exactly as s.add() spells it
	want selfWant // the direction
	file string   // path inside the throwaway fixture root, repo-relative spelling
	src  string   // bytes written there
	// cover names the counter that must be > 0 for this case's walk to have
	// happened at all. Empty means "resolve by tag" (see selfCoverFor), which is
	// the normal case; it exists so a pin can name a different scope than the one
	// its tag usually lives in.
	cover string
	// allow, when non-empty, is the fixture root's tools/d22scan/allowlist.txt
	// body. Only the allowlist-suppression pin sets it; every other case gets the
	// empty one from the skeleton.
	allow string
	// summary is what a person who has never read this file sees on the line.
	summary string
	// note carries the provenance of the less obvious expectations (registered
	// blind spots, comment-exemption temperaments, deliberate band gaps).
	note string
}

// selfBansSourceFlag is the default path of the file the roster is read from.
// The binary is invoked from this module directory (see the doc comment above and
// scripts/d22scan.sh), so a relative default is the honest one: it names the same
// file the walks name.
const selfBansSourceFlag = "main.go"

var (
	// numberedBanRe matches the tab-indented `<N> <tag>` lines of main.go's
	// `// Bans` block. The tab is load-bearing: without it the pattern also reads
	// the prose of other doc comments - measured on this file, `// 1 is the approval
	// card's own reason string` (main.go:160) and `//     141 put the opposite to the
	// owner` (main.go:891) both look like roster entries to a looser regex, and a
	// roster parse that invents ban #141 is worse than one that misses a ban.
	numberedBanRe = regexp.MustCompile(`(?m)^//\t([1-9][0-9]*)[ \t]+([a-z][a-z0-9-]+)`)
	// emittedTagRe matches a tag printed straight into an s.add() call.
	emittedTagRe = regexp.MustCompile(`s\.add\("([a-z][a-z0-9-]+)"`)
	// walkTextTagRe matches the tag of a directory-scoped ban: walkText() takes the
	// ban id as a PARAMETER and calls s.add(ban, ...) with it (main.go:875), so bans
	// #6 and #7 have no literal s.add() site at all. Measured, not assumed: reading
	// only literals reported "ban #6 panel-approval ... no s.add() can emit it",
	// which is this file being wrong about the tool, not the tool being blind.
	walkTextTagRe = regexp.MustCompile(`walkText\([^;]*?,[ \t]*"([a-z][a-z0-9-]+)",[ \t]*s\.[A-Za-z0-9_]+`)
	// indirectAddRe is the ONE emission site whose tag is a variable. Anything that
	// is neither a literal tag nor this known shape stops the run: a new channel for
	// reporting a finding must be added to the roster deliberately (the same
	// "coverage and reporting move together" rule verdict() applies to counters).
	indirectAddRe = regexp.MustCompile(`s\.add\(ban[),]`)
	addSiteRe     = regexp.MustCompile(`s\.add\(`)
)

// bansBlock returns the doc-comment group that DECLARES the bans, i.e. the run of
// consecutive `//` lines containing the `// Bans (` header. The roster is read
// from inside that group only, so neither the package prose above it nor the
// emoji-band discussion below it can be mistaken for a ban line.
func bansBlock(text string) (string, error) {
	lines := strings.Split(text, "\n")
	anchor := -1
	for i, l := range lines {
		if strings.HasPrefix(l, "// Bans (") {
			anchor = i
			break
		}
	}
	if anchor < 0 {
		return "", fmt.Errorf("no `// Bans (` header: the self-test reads its roster out of that block, so a move or a reword there must come with a re-read of this file, not a silent empty list")
	}
	start := anchor
	for start > 0 && strings.HasPrefix(lines[start-1], "//") {
		start--
	}
	end := anchor
	for end+1 < len(lines) && strings.HasPrefix(lines[end+1], "//") {
		end++
	}
	return strings.Join(lines[start:end+1], "\n"), nil
}

// banEntry is one numbered line of the `// Bans` block.
type banEntry struct {
	num int
	tag string
}

// findingTypeTags are emitted capabilities that the `// Bans` block does not
// number, and therefore must not be reported as "a ban nobody documented".
// `unparseable` is the one ticket 161 AC#1 names explicitly: a finding type, not
// a ban - excluded from the ban count, never excluded from the samples.
var findingTypeTags = map[string]bool{"unparseable": true}

// readRosters parses this tool's own source for (a) the numbered ban list and
// (b) the set of tags the code can emit. Both are required to be non-empty: a
// parse that finds nothing is an instrument hole, and it exits 2 rather than
// letting the self-test "pass" against an empty roster.
func readRosters(path string) ([]banEntry, []string, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot read the roster source %q: %w", path, err)
	}
	text := string(src)

	block, err := bansBlock(text)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", path, err)
	}

	var numbered []banEntry
	seenNum := map[int]bool{}
	for _, m := range numberedBanRe.FindAllStringSubmatch(block, -1) {
		n, err := strconv.Atoi(m[1])
		if err != nil {
			return nil, nil, fmt.Errorf("ban number %q in %s does not parse: %w", m[1], path, err)
		}
		if seenNum[n] {
			return nil, nil, fmt.Errorf("ban #%d appears twice in %s's roster block", n, path)
		}
		seenNum[n] = true
		numbered = append(numbered, banEntry{num: n, tag: m[2]})
	}
	if len(numbered) == 0 {
		return nil, nil, fmt.Errorf("%s declares no numbered bans: the `// Bans` block moved or changed shape, and the self-test refuses to audit an empty roster", path)
	}

	var emitted []string
	seenTag := map[string]bool{}
	// Comment lines are stripped before looking for emission sites, and this is not
	// cosmetic: prose that merely MENTIONS s.add() would otherwise be read as a
	// capability (this file's own flag help did exactly that on the first run).
	codeText := stripCommentLines(text)
	for _, m := range emittedTagRe.FindAllStringSubmatch(codeText, -1) {
		if seenTag[m[1]] {
			continue
		}
		seenTag[m[1]] = true
		emitted = append(emitted, m[1])
	}
	for _, m := range walkTextTagRe.FindAllStringSubmatch(codeText, -1) {
		if seenTag[m[1]] {
			continue
		}
		seenTag[m[1]] = true
		emitted = append(emitted, m[1])
	}
	// Every emission site must be accounted for by one of the two shapes above.
	for i, line := range strings.Split(codeText, "\n") {
		n := len(addSiteRe.FindAllString(line, -1))
		if n == 0 {
			continue
		}
		covered := len(emittedTagRe.FindAllString(line, -1)) + len(indirectAddRe.FindAllString(line, -1))
		if covered < n {
			return nil, nil, fmt.Errorf("%s line %d has %d s.add() site(s) of which only %d are recognisable (%q): a new way to report a finding must be taught to readRosters() on purpose, never absorbed silently", path, i+1, n, covered, strings.TrimSpace(line))
		}
	}
	if len(emitted) == 0 {
		return nil, nil, fmt.Errorf("%s contains no recognised s.add(\"<tag>\") call: nothing could ever be reported, so a green self-test would be a lie", path)
	}
	sort.Strings(emitted)
	return numbered, emitted, nil
}

// stripCommentLines drops whole-line comments. It is line-based on purpose: this
// reads a source file to enumerate what the author could print, and the failure
// mode being guarded against is prose being mistaken for code, not vice versa.
func stripCommentLines(text string) string {
	lines := strings.Split(text, "\n")
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		if strings.HasPrefix(strings.TrimSpace(l), "//") {
			out = append(out, "")
			continue
		}
		out = append(out, l)
	}
	return strings.Join(out, "\n")
}

// auditSelfCases is the self-test's own denominator check. It returns one string
// per hole; any non-empty result means the run is an instrument failure (exit 2),
// not a scan finding.
//
// The four rules are deliberately asymmetric-looking: they cover "a capability
// with no sample", "a sample with no capability", "a sample with one direction
// only" and "a numbered ban the code cannot emit". Removing any single ingredient
// from selfCases trips rule 3 for that tag, which is ticket 161 AC#4's answer.
func auditSelfCases(cases []selfCase, numbered []banEntry, emitted []string) []string {
	var holes []string

	// Rule 0: the roster must be numbered 1..N with no gap. AGENTS.md's own D-table
	// lesson ("47 枚、无缺号") is that a missing number is invisible unless someone
	// counts: a ban deleted from the middle of `// Bans` would otherwise pass every
	// other rule here as long as its s.add() survived.
	for i, b := range numbered {
		if b.num != i+1 {
			holes = append(holes, fmt.Sprintf("the `// Bans` roster is not contiguous: position %d is numbered #%d (expected #%d) - a renumbering or a deleted line has to be read, not inherited", i+1, b.num, i+1))
		}
	}

	inEmitted := map[string]bool{}
	for _, t := range emitted {
		inEmitted[t] = true
	}
	inNumbered := map[string]string{}
	for _, b := range numbered {
		inNumbered[b.tag] = fmt.Sprintf("#%d", b.num)
	}

	// Rule 4 first: the documented roster and the emitted one must be the same
	// list, apart from the tags declared as finding types above.
	for _, b := range numbered {
		if !inEmitted[b.tag] {
			holes = append(holes, fmt.Sprintf("ban %s %s is documented in the `// Bans` block but no emission site (a literal finding tag, or a walkText call naming it) can print it - the walk was removed, or the tag was renamed in only one of the two places",
				inNumbered[b.tag], b.tag))
		}
	}
	for _, t := range emitted {
		if _, ok := inNumbered[t]; ok {
			continue
		}
		if !findingTypeTags[t] {
			holes = append(holes, fmt.Sprintf("an emission site reports %q but `// Bans` does not number it - document the ban or delete the finding, never leave an unnumbered capability in the scanner", t))
		}
	}

	// Rules 1-3: the pairs.
	byTag := map[string]*[2]int{}
	unknown := []string{}
	for _, t := range emitted {
		byTag[t] = &[2]int{}
	}
	for _, c := range cases {
		pair, ok := byTag[c.tag]
		if !ok {
			unknown = append(unknown, c.tag)
			continue
		}
		if c.want == wantRing {
			pair[0]++
			continue
		}
		pair[1]++
	}
	for _, t := range emitted {
		pair := byTag[t]
		if pair[0] == 0 && pair[1] == 0 {
			holes = append(holes, fmt.Sprintf("tag %q has NO sample at all - the gate can report it, and nothing here proves it still sees anything", t))
			continue
		}
		if pair[0] == 0 {
			holes = append(holes, fmt.Sprintf("tag %q has only an expect-silent sample - a ban whose violating sample was deleted cannot be distinguished from a ban that was never implemented", t))
			continue
		}
		if pair[1] == 0 {
			holes = append(holes, fmt.Sprintf("tag %q has only an expect-ring sample - a gate that is always ringing has never been shown to be able to stay quiet, which is the tautological-check shape this repo has rejected twice", t))
			continue
		}
	}
	for _, t := range dedupe(unknown) {
		if !inEmitted[t] {
			holes = append(holes, fmt.Sprintf("a self-test case names tag %q, which no emission site in this tool can print - the sample is testing a ban that does not exist", t))
		}
	}

	// Duplicate samples would let one file stand in for two directions.
	seen := map[string]string{}
	for _, c := range cases {
		k := c.hash()
		if prev, ok := seen[k]; ok {
			holes = append(holes, fmt.Sprintf("two cases seed the same path with the same bytes (%s / %s) - one of them is not testing anything new", prev, c.summary))
		}
		seen[k] = c.summary
	}
	return holes
}

// hash identifies a case by what the scanner actually sees.
func (c selfCase) hash() string { return c.file + "\x00" + c.src + "\x00" + c.allow }

// selfSkeleton is the throwaway repository a case is scanned in. It is a whole
// repo on purpose: declaredScopes() marks every scope live, and a fixture missing
// a tree would make an expect-silent reading ambiguous between "the sample is
// clean" and "the walk had nowhere to walk".
func selfSkeleton() map[string]string {
	return map[string]string{
		"go.mod":                      "module selftest.invalid/d22scan-fixture\n\ngo 1.24\n",
		"internal/ok/ok.go":           "package ok\n",
		"cmd/probe/probe.go":          "package main\n\nfunc main() {}\n",
		"internal/tools/ok.go":        "package tools\n",
		"frontend/index.html":         "<html><body>ok</body></html>\n",
		"design/index.html":           "<html><body>ok</body></html>\n",
		"tools/d22scan/allowlist.txt": "# empty allowlist: this fixture tests the bans, not the exemptions\n",
		// Ticket 212's wantSilent case cites these two; the fixture must seed
		// them so those citations EXIST and the gate stays silent on them.
		"docs/readings.md":         "readings\n",
		"internal/probe/roster.md": "roster\n",
	}
}

// runSelfCase builds a fixture root, runs the SAME walks production runs
// (scanWithStats), and reports what the tag did.
func runSelfCase(c selfCase) (ok bool, evidence string, vacuous bool, err error) {
	root, err := os.MkdirTemp("", "d22scan-selftest-")
	if err != nil {
		return false, "", false, fmt.Errorf("fixture root: %w", err)
	}
	defer func() {
		// Go's own test fixtures get the same treatment from t.TempDir(); this is
		// the non-test path, and the directory is this call's own, in the OS temp
		// tree, never in the repository.
		_ = os.RemoveAll(root)
	}()

	files := selfSkeleton()
	if c.allow != "" {
		files["tools/d22scan/allowlist.txt"] = c.allow
	}
	files[c.file] = c.src
	relPaths := make([]string, 0, len(files))
	for rel := range files {
		relPaths = append(relPaths, rel)
	}
	sort.Strings(relPaths)
	for _, rel := range relPaths {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			return false, "", false, fmt.Errorf("seed dir for %s: %w", rel, err)
		}
		if err := os.WriteFile(full, []byte(files[rel]), 0o644); err != nil {
			return false, "", false, fmt.Errorf("seed %s: %w", rel, err)
		}
	}

	s, err := scanWithStats(root)
	if err != nil {
		return false, "", false, fmt.Errorf("scan: %w", err)
	}

	count, coverLabel, err := selfCounter(s, c, root)
	if err != nil {
		return false, "", false, err
	}
	var hits []string
	all := 0
	for _, f := range s.findings {
		all++
		if f.Ban == c.tag {
			hits = append(hits, f.String())
		}
	}

	switch {
	case count == 0:
		// Both directions are meaningless if the walk that owns this tag read
		// nothing: report it as vacuous instead of letting it read as a pass.
		return false, fmt.Sprintf("%s examined 0 files, so %q says nothing (%d finding(s) in this fixture)",
			coverLabel, c.want.String(), all), true, nil
	case c.want == wantRing && len(hits) == 0:
		return false, fmt.Sprintf("did NOT ring: no %q finding anywhere in the fixture although %s examined %d file(s); other findings: %s",
			c.tag, coverLabel, count, tagList(s.findings)), false, nil
	case c.want == wantSilent && len(hits) > 0:
		return false, fmt.Sprintf("rang on a sample that must stay silent: %s", strings.Join(hits, " | ")), false, nil
	case c.want == wantRing:
		return true, fmt.Sprintf("ranged on %d finding(s): %s", len(hits), hits[0]), false, nil
	default:
		return true, fmt.Sprintf("stayed silent while %s examined %d file(s)%s", coverLabel, count, otherFindingsSuffix(s.findings, c.tag)), false, nil
	}
}

// otherFindingsSuffix keeps a silent reading honest: if a DIFFERENT tag fired in
// the same fixture, a human sees that too instead of inferring a clean tree.
func otherFindingsSuffix(findings []Finding, tag string) string {
	var others []string
	for _, f := range findings {
		if f.Ban != tag {
			others = append(others, f.String())
		}
	}
	if len(others) == 0 {
		return "; the fixture produced zero findings at all"
	}
	return "; other tags did fire here: " + strings.Join(others, " | ")
}

// tagList renders "tag x N" per tag, for the did-not-ring message.
func tagList(findings []Finding) string {
	counts := map[string]int{}
	for _, f := range findings {
		counts[f.Ban]++
	}
	if len(counts) == 0 {
		return "none"
	}
	out := make([]string, 0, len(counts))
	for tag, n := range counts {
		out = append(out, fmt.Sprintf("%s=%d", tag, n))
	}
	sort.Strings(out)
	return strings.Join(out, " ")
}

// selfCounter resolves the work counter that must be non-zero for a case, and the
// label printed for it. The keys are the ones the walks themselves bump
// (goScopeKey / the walkText labels / emojiScopes' labels), so a rename there
// fails the self-test rather than silently making every silent case vacuous -
// which is the point.
func selfCounter(s *scanner, c selfCase, root string) (int, string, error) {
	key := c.cover
	if key == "" {
		key = selfCoverForTag[c.tag]
	}
	switch key {
	case "internalGo":
		return s.examined[goScopeKey(filepath.Join(root, "internal"))], "bans #1-5 + #8 internal/ Go walk", nil
	case "cmdGo":
		return s.examined[goScopeKey(filepath.Join(root, "cmd"))], "bans #1-5 cmd/ walk", nil
	case "panel":
		return s.examined["panel-approval"], "ban #6 frontend/ walk", nil
	case "artifact":
		return s.examined["internal-artifact-tool"], "ban #7 internal/tools/ walk", nil
	case "emojiInternal":
		return s.emojiSeen["internal/"], "ban #8 internal/ walk", nil
	case "emojiCmd":
		return s.emojiSeen["cmd/"], "ban #8 cmd/ walk", nil
	case "emojiFrontend":
		return s.emojiSeen["frontend/"], "ban #8 frontend/ walk", nil
	case "emojiDesign":
		return s.emojiSeen["design/"], "ban #8 design/ walk", nil
	case "":
		return 0, "", fmt.Errorf("case %q declares no coverage counter and no default for tag %q", c.summary, c.tag)
	default:
		return 0, "", fmt.Errorf("case %q names unknown coverage counter %q", c.summary, key)
	}
}

// selfCoverForTag maps a tag to the walk that owns it, so cases do not have to
// repeat the counter for every sample.
var selfCoverForTag = map[string]string{
	"bare-goroutine":         "internalGo",
	"pathresolver-bypass":    "internalGo",
	"plaintext-key":          "internalGo",
	"wallclock-timeout":      "internalGo",
	"mirror-hash":            "internalGo",
	"unparseable":            "internalGo",
	"panel-approval":         "panel",
	"internal-artifact-tool": "artifact",
	"emoji":                  "emojiInternal",
	"phantom-citation":       "internalGo",
}

// runSelfTest is the body of the -self-test flag and of the Go test that calls it
// in-process, so neither can drift from the other.
//
// Exit codes, stated in the same priority order verdict() uses:
//
//	2 the self-test itself is hollow (unreadable roster, roster/sample mismatch,
//	  a vacuous case, or a fixture the tool refused to scan)
//	1 a direction failed: a ban that should have rung stayed silent, or a clean
//	  sample rang
//	0 every case showed what it claims to show
func runSelfTest(out, errOut io.Writer, srcPath string) int {
	numbered, emitted, err := readRosters(srcPath)
	if err != nil {
		fmt.Fprintf(errOut, "d22scan -self-test: FATAL %v\n", err)
		return 2
	}
	holes := auditSelfCases(selfCases, numbered, emitted)
	if len(holes) > 0 {
		fmt.Fprintf(errOut, "d22scan -self-test: FATAL the table does not cover the tool (%d hole(s)); an unrun self-test is not a green self-test\n", len(holes))
		for _, h := range holes {
			fmt.Fprintf(errOut, "d22scan -self-test:   HOLE %s\n", h)
		}
		return 2
	}

	labels := map[string]string{}
	for _, b := range numbered {
		labels[b.tag] = fmt.Sprintf("ban #%d", b.num)
	}
	for tag := range findingTypeTags {
		if _, ok := labels[tag]; !ok {
			labels[tag] = "type"
		}
	}

	fmt.Fprintf(out, "d22scan -self-test: roster read from %s = %d numbered ban(s) [%s] + %d finding type(s); %d cases, %d tag(s) covered, both directions required per tag\n",
		srcPath, len(numbered), strings.Join(banList(numbered), ", "), len(findingTypeTags), len(selfCases), len(emitted))
	fmt.Fprintf(out, "d22scan -self-test: %-26s %-6s %-6s %s\n", "gate", "want", "result", "what happened")

	failed, vacuous, errored := 0, 0, 0
	for _, c := range selfCases {
		ok, evidence, vac, err := runSelfCase(c)
		if err != nil {
			fmt.Fprintf(out, "d22scan -self-test: %-26s %-6s %-6s %s: %v\n",
				c.tag, c.want, "ERROR", c.summary, err)
			errored++
			continue
		}
		result := "OK"
		if !ok {
			result = "FAIL"
			failed++
			if vac {
				vacuous++
				result = "VACUOUS"
			}
		}
		gate := labels[c.tag] + " " + c.tag
		fmt.Fprintf(out, "d22scan -self-test: %-26s %-6s %-6s %s | %s\n", gate, c.want, result, c.summary, evidence)
		if c.note != "" {
			fmt.Fprintf(out, "d22scan -self-test: %-26s %-6s %-6s pinned because: %s\n", "", "", "", c.note)
		}
	}

	if errored > 0 || vacuous > 0 {
		fmt.Fprintf(errOut, "d22scan -self-test: %d case(s) could not be judged (%d error, %d vacuous) - %d/%d passed\n",
			errored+vacuous, errored, vacuous, len(selfCases)-errored-vacuous, len(selfCases))
		return 2
	}
	if failed > 0 {
		fmt.Fprintf(errOut, "d22scan -self-test: %d direction(s) failed, %d/%d passed - the gate does not see what it claims\n",
			failed, len(selfCases)-failed, len(selfCases))
		return 1
	}
	fmt.Fprintf(out, "d22scan -self-test: clean - all %d direction checks passed (%d expect-ring, %d expect-silent)\n",
		len(selfCases), countWant(selfCases, wantRing), countWant(selfCases, wantSilent))
	return 0
}

func banList(numbered []banEntry) []string {
	parts := make([]string, 0, len(numbered))
	for _, b := range numbered {
		parts = append(parts, fmt.Sprintf("%d %s", b.num, b.tag))
	}
	return parts
}

func countWant(cases []selfCase, w selfWant) int {
	n := 0
	for _, c := range cases {
		if c.want == w {
			n++
		}
	}
	return n
}

func dedupe(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		if seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}
