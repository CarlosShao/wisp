// Command d22scan is the D22 seven-ban + zero-emoji static scanner
// (ticket 08; SPEC-10 §8 item 5). It is the CI lint gate and the local
// self-check; loosening or skipping it is the D22 run-away mode 6.
//
// Bans (PLAN.md D22 fourth-round additions), production scope internal/ +
// cmd/ (non-test, non-testdata) unless stated:
//
//	1 bare-goroutine      every `go <anything>` - closure literal OR named
//	                      call (widened by R16#1; the pre-R16 matcher saw only
//	                      `go func(` and "clean" therefore meant "no bare
//	                      closures", not "nothing bypasses Spawn").
//	                      Every goroutine via observe.Registry.Spawn; the only
//	                      file-level exemption is internal/observe/goroutine.go,
//	                      which implements Registry itself.
//	2 pathresolver-bypass filepath.Clean / filepath.Abs outside the allowlist
//	                      (C26 PathResolver is the only path comparison)
//	3 plaintext-key       key/token/secret-named identifier assigned a
//	                      string literal (SecretStore C28 refs only)
//	4 wallclock-timeout   wall-clock deltas in timeout/deadline logic
//	                      (D42#9: monotonic clock only)
//	5 mirror-hash         hash material mentioned together with a mirror
//	                      (C29: hashes come from the signed manifest only)
//	6 panel-approval      `approval.decide` in frontend/ (D33/F2: allow
//	                      decisions are native-side only) - scope: frontend/
//	                      (TEXT UNCHANGED - D22 owns it. ARMED LIVE since
//	                      ticket 88: ticket 77 landed frontend/ at 63ef895, so
//	                      the exemption ticket 67 registered ("the tree does not
//	                      exist, keeping this would be pretending to scan")
//	                      became a false claim, and its own drift guard fired.
//	                      The walk now examines 40 text files in this repo
//	                      (measured 2026-09-25 at 467f8a4 on a `git archive HEAD`
//	                      snapshot, and the SAME 40 on a worktree that has run a
//	                      frontend build - ledger A207: before the walks excluded
//	                      gitignored paths a built worktree read 43 here, and the
//	                      ledger therefore carried 37 / 40 / 43 as three mutually
//	                      irreconcilable readings of one scope)
//	                      and an empty scope is fatal - see declaredScopes)
//	7 internal-artifact-tool  host-internal artifact writes implemented as
//	                      gated tool names (D34 note 2) - scope: internal/tools/
//	8 emoji               zero emoji in design/ (every text file), in frontend/
//	                      (EVERY file, ticket 96 - same walk shape as ban #6, no
//	                      suffix allowlist) and in the Go sources of internal/ +
//	                      cmd/ - comments and
//	                      _test.go INCLUDED (D23). This is the one ban whose
//	                      scope is NOT the "production, non-test" default
//	                      above, and it is deliberately so: see emojiScopes
//	                      and walkEmoji for the reasons, and pin any change
//	                      there in scan_test.go rather than in the footer,
//	                      which is generated from the scope list.
//
// Findings are suppressed only via allowlist.txt entries of the form
// "ban-id<TAB>repo-relative path prefix<TAB>reason" (committed, reviewable,
// never wildcards beyond the path prefix).
//
// What the walks DESCEND into is a second, separate question, answered by the
// repository's own .gitignore files PLUS this tree's git index (gitignore.go): a
// path is skipped only when a rule matches it AND git confirms the index does not
// hold it. Before that, "ban #6 frontend/ examined N" meant whatever the machine
// held - a tree with a frontend build in it read 43 where CI reads 40, which is
// why the ledger could cite 37 / 40 / 43 for one scope (A207).
//
// THE TWO SENTENCES BELOW ARE ASSERTED, NOT ASPIRATIONAL (A214/F2 called the
// earlier version of them unfalsifiable, and it was right - nothing pinned them):
//   - it never suppresses a finding on a TRACKED path. skip() asks
//     `git ls-files -z` and `git ls-files -i -c --exclude-standard` and refuses to
//     skip anything either one reports, so a `git add -f` file under `dist/`, or a
//     force-added `internal/build/*.go` under the root `build/` rule, is read and
//     can go red. scan_test.go's TestTrackedPathsAreNeverSkippedByTheIgnoreFilter
//     builds exactly that shape with a real `git init` + `git add -f` and demands
//     the findings AND the denominators move; TestWalksSkipGitIgnoredPaths is the
//     other half (ignored-and-untracked stays silent, non-ignored goes red).
//   - when git cannot be asked at all (no binary, not a repository, -root is a
//     subdirectory of a foreign one, the command fails, times out, or the output
//     does not parse) NO ignore rule is applied: every path is scanned and the
//     run prints a loud line saying so. The degradation is over-coverage, never a
//     quiet clean - see gitignore.go's note().
//
// And it still does not hide untracked-but-not-ignored work - `git add` has not
// happened yet is not the same claim as "this is not part of the deliverable".
//
// Invocation (ticket 67 AC#2 - the shape of this command is load-bearing):
// this package is its OWN Go module (tools/d22scan/go.mod), so from the repo
// root `go run ./tools/d22scan` fails inside the ROOT module and this scanner
// never starts. Worse, running it with the default `-root .` from inside
// tools/d22scan points the walk at a tree that has no internal/ or cmd/, so
// every ban walks zero files and the tool would print "clean" while having
// examined nothing. Both mis-invocations are now fatal: main validates
// -root and refuses to report a verdict it did not actually compute. Run it
// as:
//
//	scripts/d22scan.sh                      # what CI's lint job calls
//	cd tools/d22scan && go run . -root ../../
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Finding is one detected violation.
type Finding struct {
	Ban  string `json:"ban"`
	Path string `json:"path"`
	Line int    `json:"line"`
	What string `json:"what"`
}

func (f Finding) String() string {
	return fmt.Sprintf("%s:%d: [%s] %s", f.Path, f.Line, f.Ban, f.What)
}

// scanner carries the allowlist and the findings.
type scanner struct {
	root      string
	allow     map[string]map[string]bool // ban-id -> path prefix set
	findings  []Finding
	failAddOn string         // unused placeholder guard
	emojiSeen map[string]int // ban #8: scope label -> files actually line-scanned
	examined  map[string]int // bans #1-7: scope key -> files actually walked
	// ign answers "would git decline to track this path?" - the scanned tree's own
	// .gitignore files for the pattern half and `git ls-files` for the index half.
	// Every walk consults it, so the denominators count the files a verdict can be
	// about instead of whatever a developer's machine holds (ledger A207: frontend/
	// read 43 on a built worktree and 40 on a clean checkout, which made 37/40/43
	// indistinguishable from a coverage change), and so that a tracked file is
	// never hidden behind a pattern (A214/F1). See gitignore.go for why the rule is
	// read, not copied, and for the no-index fallback that applies no rule at all.
	ign *gitIgnore
}

var (
	secretNameWords  = map[string]bool{"key": true, "token": true, "secret": true, "password": true, "passwd": true, "credential": true, "passphrase": true}
	secretNameRe     = regexp.MustCompile(`(?i)(?:^|[^a-z])(api[_-]?key|secret|passw(or)?d|credential|passphrase|access[_-]?token|refresh[_-]?token)(?:[^a-z]|$)`)
	wallclockRe      = regexp.MustCompile(`\.Sub\(time\.Now\(\)\)`)
	unixTimeRe       = regexp.MustCompile(`time\.Now\(\)\.(Unix|UnixNano|UnixMilli)\(`)
	timeoutWordRe    = regexp.MustCompile(`(?i)timeout|deadline|expire|\bttl\b|budget|until`)
	mirrorHashRe     = regexp.MustCompile(`(?i)sha3?-?(256|512)|checksum|\bhash`)
	mirrorWordRe     = regexp.MustCompile(`(?i)mirror`)
	approvalPanelRe  = regexp.MustCompile(`approval\.decide`)
	artifactToolRe   = regexp.MustCompile(`(?i)"(spill|internal[._-][a-z0-9_.-]+)"`)
	assignKeyShapeRe = regexp.MustCompile(`^[A-Za-z0-9._~-]{16,}$`)
)

// emojiRe is ban #8's character class.
//
// Q-46(c) widened it by \x{2200}-\x{22FF} (ticket 141): PLAN.md's own list of
// ranges does NOT contain that band either, so the spec is narrower than what
// the ban exists to prevent - of the non-comment occurrences of U+2265,
// U+2264, U+2229 and U+2212 in this repo, 4 are already-rendered panel copy and
// 1 is the approval card's own reason string. Do not "tidy" the bands into one
// span: \x{2190}-\x{21FF} (arrows) and \x{2460}-\x{24FF} (circled numbers) are
// deliberately still absent, because the approval added ONLY the math band, and
// any further band belongs in scan_test.go's pin, not here.
var emojiRe = regexp.MustCompile(`[\x{1F000}-\x{1FAFF}\x{2200}-\x{22FF}\x{2600}-\x{27BF}\x{2B00}-\x{2BFF}\x{FE0F}\x{1F1E6}-\x{1F1FF}]`)

func splitSecretWords(name string) []string {
	// Split CamelCase and non-alphanumeric separators into lowercase words.
	var words []string
	var cur strings.Builder
	for _, r := range name {
		switch {
		case r >= 'A' && r <= 'Z':
			if cur.Len() > 0 {
				words = append(words, cur.String())
				cur.Reset()
			}
			cur.WriteRune(r + 'a' - 'A')
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			cur.WriteRune(r)
		default:
			if cur.Len() > 0 {
				words = append(words, cur.String())
				cur.Reset()
			}
		}
	}
	if cur.Len() > 0 {
		words = append(words, cur.String())
	}
	return words
}

func isSecretNamed(name string) bool {
	words := splitSecretWords(name)
	if len(words) == 0 {
		return false
	}
	// The secret word must be a word of the name (suffix preferred to avoid
	// "PayPerToken"-style false positives on plain "token").
	for i, w := range words {
		if secretNameWords[w] && (i == len(words)-1 || i > 0) {
			return true
		}
	}
	return false
}

// Scan walks the repo at root and returns all non-allowlisted findings.
func Scan(root string) ([]Finding, error) {
	s, err := scanWithStats(root)
	if err != nil {
		return nil, err
	}
	return s.findings, nil
}

// scanWithStats is Scan plus the per-scope work counts: "no findings" and
// "nothing was looked at" stay distinguishable (ticket 67 AC#2 bought the same
// property for the Go scope with checkRoot; main.go enforces it for ban #8).
func scanWithStats(root string) (*scanner, error) {
	s := &scanner{
		root: root, allow: map[string]map[string]bool{},
		emojiSeen: map[string]int{}, examined: map[string]int{},
		ign: newGitIgnore(root),
	}
	if err := s.loadAllowlist(filepath.Join(root, "tools", "d22scan", "allowlist.txt")); err != nil {
		return nil, err
	}
	// Production Go scope.
	if err := s.walkGo(filepath.Join(root, "internal")); err != nil {
		return nil, err
	}
	if err := s.walkGo(filepath.Join(root, "cmd")); err != nil {
		return nil, err
	}
	// Directory-scoped bans.
	if err := s.walkText(filepath.Join(root, "frontend"), "panel-approval", s.panelCheck, false); err != nil {
		return nil, err
	}
	if err := s.walkText(filepath.Join(root, "internal", "tools"), "internal-artifact-tool", s.artifactCheck, true); err != nil {
		return nil, err
	}
	for _, sc := range emojiScopes(root) {
		if err := s.walkEmoji(sc); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// emojiScope is one tree ban #8 walks.
type emojiScope struct {
	dir    string // absolute path
	label  string // repo-relative label, quoted in findings and in the self-report
	goOnly bool   // true: only .go files count (see walkEmoji's scope note)
	// everyFile true scans EVERY file under dir, with no suffix allowlist -
	// the walk shape ban #6 uses on the same tree (walkText). Ticket 96 added
	// it because the alternative for frontend/ was the goOnly:false branch,
	// which filters through isTextFile(): measured on `git ls-files frontend`,
	// that whitelist drops 5 of 37 tracked files, including all three
	// scripts/*.mjs (the vendoring tooling R18 exists to police). Reusing it
	// would have made ban #8's view of frontend/ a strict SUBSET of ban #6's,
	// i.e. the "gate blind to a whole class of file" shape R16#4 forbids
	// re-narrowing a ban into. Mutually exclusive with goOnly.
	everyFile bool
}

// scanScope is ONE DECLARED WALK of the whole scan, whichever ban it serves.
//
// Before ticket 71 AC#4 the tool only kept a ledger for ban #8 (emojiSeen), so
// "examined 200 production Go files" was the only Go-side number and it was an
// AGGREGATE: a run where cmd/ walked nothing still printed 200 and read as
// healthy, and the two directory-scoped bans (#6 frontend/, #7 internal/tools/)
// printed no number at all. A scope nobody counts is a scope nobody notices
// going empty - which is precisely how ban #6's frontend/ survived as a
// permanent 0-file walk while the footer advertised it (registry A22/A26/A30).
//
// Every walk now has an entry here, the self-report prints its real count, and
// the guards below decide what a 0 means. The count is read from whichever map
// the walk writes into (emojiSeen for ban #8, examined for the rest) rather
// than being duplicated, so there is no second bookkeeping to drift - and it is
// read by KEY, not by closure, so undeclaredKeys() can verify that every
// counter a walk bumped is actually in the report.
type scanScope struct {
	label       string // printed verbatim in the self-report; also the guard's message
	dir         string // absolute path the walk targets (used by the absent-tree drift guard)
	kind        string // wording: what one "file" means here
	seenKey     string // ban #8's emojiSeen key; empty when examinedKey is used
	examinedKey string // bans #1-7's examined key; empty when seenKey is used
	live        bool   // true: 0 examined files is fatal; false: see absentOK below
	absentOK    bool   // true only for a scope whose tree is KNOWN ABSENT (no such scope in the ledger today - see driftedAbsentScope)
	note        string // printed with the scope so a disclaimer cannot be dropped
}

// count reads this scope's work counter from the map the walk bumped.
func (sc scanScope) count(s *scanner) int {
	if sc.examinedKey != "" {
		return s.examined[sc.examinedKey]
	}
	return s.emojiSeen[sc.seenKey]
}

// key names the counter this scope consumes, qualified by the map it lives in,
// so undeclaredKeys can compare the ledger against what the walks wrote.
func (sc scanScope) key() string {
	if sc.examinedKey != "" {
		return "examined:" + sc.examinedKey
	}
	return "emojiSeen:" + sc.seenKey
}

// goScopeKey is the ledger key walkGo bumps for a directory it walks. It is a
// function so declaredScopes and walkGo cannot each invent their own spelling -
// a mismatch there would read as "scope walked 0 files".
func goScopeKey(dir string) string { return filepath.ToSlash(dir) + "/" }

// undeclaredKeys returns the counters a walk bumped that no declared scope
// reads ("map:key" form): work that happened but is missing from the
// self-report. Its counterpart (a declared scope whose walk was deleted) is
// caught by emptyLiveScope, so between the two, coverage can be neither added
// nor removed without the tool refusing to print a verdict (ticket 71 AC#4).
func undeclaredKeys(s *scanner, scopes []scanScope) []string {
	declared := map[string]bool{}
	for _, sc := range scopes {
		declared[sc.key()] = true
	}
	var out []string
	for k := range s.examined {
		if !declared["examined:"+k] {
			out = append(out, "examined:"+k)
		}
	}
	for k := range s.emojiSeen {
		if !declared["emojiSeen:"+k] {
			out = append(out, "emojiSeen:"+k)
		}
	}
	sort.Strings(out)
	return out
}

// declaredScopes is the single list of everything the scan claims to look at.
// main() renders its self-report from this and only this, so a new walk that is
// not registered here cannot run silently, and a registered walk cannot hide
// its count.
//
// WHY ban #6 IS REGISTERED AS absentOK INSTEAD OF LIVE (and why its text is
// untouched): D22 puts the prohibition WORDING out of an agent's reach, and the
// orchestrator's 2026-09-21 handover repeats that: the "no approval.decide in
// the panel" rule stays exactly as written. What ticket 71 AC#4 governs is the
// *instrument*, and the instrument's two legal states are "has coverage" and
// "says out loud that it has none". frontend/ does not exist at this HEAD
// (`ls frontend/` -> No such file or directory, measured 2026-09-21), so the
// walk can only ever examine 0 files. Making that fatal would make CI red
// forever for a reason no agent may fix; making it silent is the original lie.
// So it is declared non-live WITH a note, and the drift guard below fires the
// moment frontend/ reappears while the entry is still exempt - the exemption
// cannot outlive the fact that justified it. Re-arm it (flip live:true) in the
// same commit that lands ticket 34's panel scaffold.
//
// SUPERSEDED BY TICKET 88 (2026-09-21), kept verbatim above because the paragraph
// above is a measurement claim that has since become false, and a correction that
// erases the claim leaves no trace of the mistake. What changed: ticket 77 created
// frontend/ (63ef895), so "the tree does not exist" stopped being a fact and the
// exemption became the lie it was invented to prevent - driftedAbsentScope fired
// exactly as designed (measured: `sh scripts/d22scan.sh` on a `git archive HEAD`
// snapshot exits 2 naming ban #6). The entry is now live:true, the ban TEXT is
// still untouched (D22: not an agent's to shorten), and the empty-scope rule that
// made the flip cost 5 restated ledger tests is still the rule nobody may soften:
// an instrument that examines 0 files does not get to print a verdict.
func declaredScopes(root string) []scanScope {
	internalGo, cmdGo := goScopeKey(filepath.Join(root, "internal")), goScopeKey(filepath.Join(root, "cmd"))
	return append([]scanScope{
		{
			label: "bans #1-5 internal/", dir: filepath.Join(root, "internal"),
			kind: "production Go files", examinedKey: internalGo, live: true,
		},
		{
			label: "bans #1-5 cmd/", dir: filepath.Join(root, "cmd"),
			kind: "production Go files", examinedKey: cmdGo, live: true,
		},
		{
			label: "ban #6 frontend/", dir: filepath.Join(root, "frontend"),
			kind: "text files", examinedKey: "panel-approval", live: true,
			note: "armed by ticket 88 on 2026-09-21, in the batch after ticket 77 landed " +
				"frontend/ at 63ef895 - ticket 67 had registered this scope exempt on the " +
				"ground that the tree did not exist, and that ground is gone. live:true means " +
				"exactly one thing here: 0 examined files is fatal. The ban TEXT is unchanged " +
				"(D22: not an agent's to shorten), and walkText gives it every file under " +
				"frontend/ that is not node_modules/, testdata/ or a gitignored path " +
				"(A207) - .tsx, .ts, .mjs, .css, " +
				".md, .json and extension-less files included, no suffix allowlist to shrink " +
				"the ban into (35 files as measured 2026-09-21; 40 as measured 2026-09-25 at " +
				"467f8a4, and that 40 is now the reading on a worktree that HAS run a frontend " +
				"build too, which it was not: a built tree read 43 here until the walks started " +
				"excluding what .gitignore excludes).",
		},
		{
			label: "ban #7 internal/tools/", dir: filepath.Join(root, "internal", "tools"),
			kind: "production Go files", examinedKey: "internal-artifact-tool", live: true,
		},
	}, ban8Scopes(root)...)
}

// ban8Scopes adapts emojiScopes() (ban #8's own truth source, unchanged) into
// the ledger. Kept separate so ticket 67's emojiScopes tests keep passing and
// there is still exactly one list naming ban #8's trees.
func ban8Scopes(root string) []scanScope {
	out := make([]scanScope, 0, 3)
	for _, sc := range emojiScopes(root) {
		kind := "text files"
		if sc.goOnly {
			kind = "Go files, comments and _test.go included"
		}
		out = append(out, scanScope{
			label: "ban #8 " + sc.label, dir: sc.dir, kind: kind,
			seenKey: sc.label, live: true,
		})
	}
	return out
}

// emptyLiveScope returns the label of the first live scope that walked zero
// files ("" if none). This generalizes emptyEmojiScope to every scope in the
// ledger, which is the delta ticket 71 AC#4 asks for: the guard may no longer
// be something only ban #8 is subject to.
func emptyLiveScope(scopes []scanScope, s *scanner) string {
	for _, sc := range scopes {
		if sc.live && sc.count(s) == 0 {
			return sc.label
		}
	}
	return ""
}

// driftedAbsentScope returns the label of a scope registered as "tree absent"
// whose tree is now PRESENT ("" if none). It is the counterpart guard: without
// it, an absentOK entry is a hole that stays open after the reason closes.
//
// NO SCOPE IN declaredScopes() SETS absentOK ANY MORE (ticket 88 armed ban #6,
// the last exemption). This guard is NOT dead code to delete: it is what makes
// an exemption self-expiring, so the next agent that registers a scope as
// "absent, therefore exempt" cannot leave it there after the tree appears.
// Because the production ledger no longer exercises it, its teeth are pinned in
// scan_test.go against a SYNTHETIC exempt scope (TestExemptScopeCannotOutliveItsAbsentTree),
// in both directions - a guard only ever observed passing is not a guard.
func driftedAbsentScope(scopes []scanScope) string {
	for _, sc := range scopes {
		if !sc.live && sc.absentOK {
			if _, err := os.Stat(sc.dir); err == nil {
				return sc.label
			}
		}
	}
	return ""
}

// uncoveredScopes lists the exempt scopes, so the verdict line can state what
// it did NOT check instead of implying total coverage (A30's "clean was a lie"
// shape).
func uncoveredScopes(scopes []scanScope) []scanScope {
	var out []scanScope
	for _, sc := range scopes {
		if !sc.live {
			out = append(out, sc)
		}
	}
	return out
}

// describeScopes renders the per-scope work lines. main() prints exactly these
// lines and builds its verdict from them, so the sentence can never name a
// scope the ledger does not contain (the ticket 67 AC#3 footer bug, generalized).
func describeScopes(scopes []scanScope, s *scanner) []string {
	out := make([]string, 0, len(scopes))
	for _, sc := range scopes {
		line := fmt.Sprintf("d22scan: scope %-24s examined %3d %s", sc.label, sc.count(s), sc.kind)
		if !sc.live {
			line += "  [NOT COVERED]"
		}
		out = append(out, line)
	}
	return out
}

// emojiScopes is ban #8's ENTIRE coverage, kept in one place so the tool cannot
// print a footer naming a tree it never walked - the failure shape of A22 /
// A26 / A30, and exactly what this list was before ticket 67 AC#3 landed:
// design/ plus frontend/, and frontend/ does not exist at this HEAD, so the
// ban was blind to every .go file of the product while printing "no emoji in
// design/ or frontend/".
//
// WHY frontend/ WAS DELETED INSTEAD OF FIXED (ticket 71 AC#4's two allowed
// outcomes are "give it real coverage" or "delete it"; only one is available):
//   - It cannot have coverage today: `ls frontend/` is "No such file or
//     directory" (measured 2026-09-21), so the entry could only ever walk zero
//     files. Keeping it was pretending to scan.
//   - Making an absent scope non-fatal was the lie; so the empty-scope check in
//     main is fatal (exit 2). Re-add frontend/ to this list in the SAME commit
//     that lands ticket 34's scaffold - emojiScopes() is now the single source
//     of truth for ban #8's coverage, so "the panel tree is not covered" reads
//     off one function body instead of hiding behind a walk that never runs.
//   - NOT touched here: ban #6 (panel-approval) still walks frontend/ via
//     walkText above and is still an always-0 scope. Deleting it would narrow a
//     ban I do not own, which R16#4 forbids an agent to do; ticket 71 AC#4 owns
//     its disposal. It is listed as a known residual in
//     docs/evidence/s1/67-emoji-scope-internal-cmd.md rather than left unsaid.
//
// SUPERSEDED IN PART BY TICKET 71 AC#4 (2026-09-21), kept above verbatim because
// a correction that erases the original claim leaves no trace of the mistake:
// ban #6 IS disposed of now, but NOT by deleting it - the walk stays, the ban
// TEXT stays (D22: not an agent's to shorten), and the scope is registered in
// declaredScopes() as exempt-with-a-note so the run prints "NOT COVERED
// ban #6 frontend/" on every invocation. What ticket 71 closed is the silence:
// before it no number at all was printed for bans #1-5 per directory, ban #6 or
// ban #7, so "the walk ran empty" was indistinguishable from "the walk found
// nothing". See scanScope for the ledger and driftedAbsentScope for why the
// exemption expires by itself.
//
// SUPERSEDED IN PART BY TICKET 96 (2026-09-21): frontend/ is back in this list
// under the OTHER allowed outcome - real coverage, not a deletion. The paragraphs
// above stand as the reason it was deleted once, and that reason ("the tree does
// not exist") is exactly what ticket 77 changed at 63ef895: ticket 88 armed ban
// #6 over frontend/ and measures 40 text files on this HEAD, yet this list still
// had three entries, so ban #8 read ZERO bytes of the panel while the run printed
// three confident scope lines. Different shape from ticket 67's - no exemption
// was ever registered here, so nothing could drift and no guard could fire: the
// list was simply never synchronised. That is why the coverage is pinned by
// literal label in scan_test.go rather than by inspecting this function body.
// frontend/ is declared everyFile, NOT goOnly:false, because the goOnly:false
// branch filters through isTextFile(), whose suffix whitelist drops 5 of the 37
// tracked files under frontend/ (all three scripts/*.mjs, .gitignore,
// frontend/dist/.gitkeep) and would make ban #8's view of this tree a strict
// subset of ban #6's over the same files. Widening a ban is allowed; handing it a
// filter is not (R16#4).
//
// SUPERSEDED IN PART BY A207 (2026-09-25), kept above because the paragraph is a
// measurement claim and a correction that erases the claim leaves no trace: the
// "40 text files" it quotes for ban #6 was true ONLY of a tree with no build
// output in it. Measured at 467f8a4 the same binary read 43 on the worktree where
// ticket 96's ban #8 frontend/ entry was being argued and 40 on a
// `git archive HEAD` snapshot, so every count in this paragraph was machine-
// dependent, which is what a denominator must never be. The walks now exclude
// what .gitignore excludes (gitignore.go), and 40 is the reading of both shapes.
func emojiScopes(root string) []emojiScope {
	return []emojiScope{
		{dir: filepath.Join(root, "design"), label: "design/"},
		{dir: filepath.Join(root, "frontend"), label: "frontend/", everyFile: true},
		{dir: filepath.Join(root, "internal"), label: "internal/", goOnly: true},
		{dir: filepath.Join(root, "cmd"), label: "cmd/", goOnly: true},
	}
}

// describeEmojiScopes renders the per-scope work counts. main() builds its
// "clean" line out of this string, so the sentence cannot claim a scope the
// list does not contain (the ticket 67 AC#3 footer bug).
func describeEmojiScopes(scopes []emojiScope, seen map[string]int) string {
	parts := make([]string, 0, len(scopes))
	for _, sc := range scopes {
		kind := "text files"
		if sc.goOnly {
			kind = "Go files, comments and _test.go included"
		}
		parts = append(parts, fmt.Sprintf("%s %d %s", sc.label, seen[sc.label], kind))
	}
	return strings.Join(parts, "; ")
}

// emptyEmojiScope returns the label of the first declared ban #8 scope that
// line-scanned zero files, or "" when every scope did real work. A declared
// scope with no files behind it is an empty instrument (ticket 71 AC#4): main
// treats it as fatal instead of letting it read as green.
func emptyEmojiScope(scopes []emojiScope, seen map[string]int) string {
	for _, sc := range scopes {
		if seen[sc.label] == 0 {
			return sc.label
		}
	}
	return ""
}

func (s *scanner) loadAllowlist(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for i, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) < 3 {
			return fmt.Errorf("allowlist line %d malformed (want ban-id\\tpath\\treason): %q", i+1, line)
		}
		if s.allow[parts[0]] == nil {
			s.allow[parts[0]] = map[string]bool{}
		}
		s.allow[parts[0]][filepath.ToSlash(parts[1])] = true
	}
	return nil
}

func (s *scanner) allowed(ban, path string) bool {
	rel := filepath.ToSlash(path)
	if !strings.HasPrefix(rel, "/") && !strings.Contains(rel, ":") {
		// relative already
	} else {
		rel = strings.TrimPrefix(strings.TrimPrefix(rel, filepath.ToSlash(s.root)), "/")
	}
	rel = strings.TrimPrefix(rel, "/")
	for prefix := range s.allow[ban] {
		if strings.HasPrefix(rel, prefix) {
			return true
		}
	}
	return false
}

func (s *scanner) add(ban, path string, line int, what string) {
	if s.allowed(ban, path) {
		return
	}
	rel, err := filepath.Rel(s.root, path)
	if err != nil {
		rel = path
	}
	s.findings = append(s.findings, Finding{Ban: ban, Path: filepath.ToSlash(rel), Line: line, What: what})
}

// walkGo runs the Go-level bans over production files of dir. Every file it
// hands to scanGoFile is counted into the ledger under the repo-relative dir,
// which is what makes bans #1-5 report per-directory work instead of one
// aggregate that hides an empty directory (ticket 71 AC#4).
func (s *scanner) walkGo(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}
	key := goScopeKey(dir)
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "testdata" || d.Name() == ".git" || s.ign.skip(path, true) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		// A .go file this walk skips is one git itself would decline to track: build
		// output or a scratch copy. A tracked file under an ignore rule (the
		// force-added `internal/build/leak.go` shape, A214/F1) never reaches this
		// line, which is the whole reason the skip is index-aware and not a pattern
		// match - counting it is what made the denominator move between machines
		// (A207), reading it is what makes a delivered byte visible.
		if s.ign.skip(path, false) {
			return nil
		}
		s.examined[key]++
		return s.scanGoFile(path)
	})
}

func (s *scanner) scanGoFile(path string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		// Unparseable production code is a lint failure in itself.
		s.add("unparseable", path, 1, err.Error())
		return nil
	}

	// Line maps for the regex bans (skip pure comment lines).
	src, _ := os.ReadFile(path)
	lines := strings.Split(string(src), "\n")

	ast.Inspect(f, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.GoStmt:
			// R16#1: the ban covers EVERY `go <anything>`, not just the closure
			// literal. `go probeReader()` spawns exactly as unsafely as
			// `go func(){}`, so matching only FuncLit left the gate blind to
			// named calls. The one sanctioned spawn point is exempted BY FILE
			// PATH (internal/observe/goroutine.go - it implements Registry),
			// never by call shape; see tools/d22scan/allowlist.txt.
			pos := fset.Position(v.Pos())
			if _, ok := v.Call.Fun.(*ast.FuncLit); ok {
				s.add("bare-goroutine", path, pos.Line,
					"bare `go func(` is banned (D22/D38b): use observe.Registry.Spawn (named, owner, recover boundary)")
				break
			}
			s.add("bare-goroutine", path, pos.Line,
				"bare `"+goStmtText(v)+"` is banned (D22/D38b, R16: named calls count too): use observe.Registry.Spawn (named, owner, recover boundary)")

		case *ast.CallExpr:
			if sel, ok := v.Fun.(*ast.SelectorExpr); ok {
				if id, ok := sel.X.(*ast.Ident); ok && id.Name == "filepath" {
					if sel.Sel.Name == "Clean" || sel.Sel.Name == "Abs" {
						pos := fset.Position(v.Pos())
						s.add("pathresolver-bypass", path, pos.Line,
							"filepath."+sel.Sel.Name+" outside the C26 PathResolver is banned (D22); see tools/d22scan/allowlist.txt for the sanctioned exceptions")
					}
				}
			}
		case *ast.ValueSpec:
			for _, name := range v.Names {
				if !isSecretNamed(name.Name) {
					continue
				}
				for _, val := range v.Values {
					if lit, ok := val.(*ast.BasicLit); ok && lit.Kind == token.STRING {
						if raw, err := strconv.Unquote(lit.Value); err == nil && looksLikePlaintextKey(raw) {
							pos := fset.Position(lit.Pos())
							s.add("plaintext-key", path, pos.Line,
								"identifier "+name.Name+" holds a plaintext key literal (D22/C28: SecretStore refs only)")
						}
					}
				}
			}
		case *ast.AssignStmt:
			for i, lhs := range v.Lhs {
				id, ok := lhs.(*ast.Ident)
				if !ok || i >= len(v.Rhs) {
					continue
				}
				if !isSecretNamed(id.Name) {
					continue
				}
				if lit, ok := v.Rhs[i].(*ast.BasicLit); ok && lit.Kind == token.STRING {
					if raw, err := strconv.Unquote(lit.Value); err == nil && looksLikePlaintextKey(raw) {
						pos := fset.Position(lit.Pos())
						s.add("plaintext-key", path, pos.Line,
							"identifier "+id.Name+" assigned a plaintext key literal (D22/C28: SecretStore refs only)")
					}
				}
			}
		}
		return true
	})

	for i, line := range lines {
		code := strings.TrimSpace(line)
		if code == "" || strings.HasPrefix(code, "//") {
			continue
		}
		if wallclockRe.MatchString(code) {
			s.add("wallclock-timeout", path, i+1,
				"wall-clock delta (Sub(time.Now())) in what must be monotonic logic (D42#9/D22)")
		} else if unixTimeRe.MatchString(code) && timeoutWordRe.MatchString(code) {
			s.add("wallclock-timeout", path, i+1,
				"unix timestamp on a timeout/deadline line - monotonic clock required (D42#9/D22)")
		}
		if mirrorWordRe.MatchString(code) && mirrorHashRe.MatchString(code) {
			s.add("mirror-hash", path, i+1,
				"hash material mentioned together with a mirror (C29/F3: hashes come from the signed manifest only)")
		}
	}
	return nil
}

// goStmtText renders a compact "go f()" spelling of a go statement for the
// finding message, so a named call reads distinctly from a closure literal.
func goStmtText(v *ast.GoStmt) string {
	var b strings.Builder
	b.WriteString("go ")
	switch fn := v.Call.Fun.(type) {
	case *ast.Ident:
		b.WriteString(fn.Name)
	case *ast.SelectorExpr:
		if id, ok := fn.X.(*ast.Ident); ok {
			b.WriteString(id.Name + "." + fn.Sel.Name)
		} else {
			b.WriteString(fn.Sel.Name)
		}
	case *ast.FuncLit:
		b.WriteString("func()")
	default:
		b.WriteString("call")
	}
	b.WriteString("(...)")
	return b.String()
}

// looksLikePlaintextKey filters ref kinds, placeholders and low-entropy words.
func looksLikePlaintextKey(raw string) bool {
	if raw == "" {
		return false
	}
	lower := strings.ToLower(raw)
	for _, prefix := range []string{"env:", "dpapi:", "placeholder", "ref:", "${"} {
		if strings.HasPrefix(lower, prefix) {
			return false
		}
	}
	if strings.ContainsAny(raw, " \t\n") {
		return false
	}
	if !assignKeyShapeRe.MatchString(raw) {
		return false
	}
	// Require at least one digit: distinguishes key material from words
	// like "pay-per-token".
	hasDigit := strings.ContainsAny(raw, "0123456789")
	return hasDigit
}

func (s *scanner) panelCheck(line string) (string, bool) {
	if approvalPanelRe.MatchString(line) {
		return "`approval.decide` in frontend/ is banned (D33/F2: allow decisions are native-side only)", true
	}
	return "", false
}

func (s *scanner) artifactCheck(line string) (string, bool) {
	if artifactToolRe.MatchString(line) {
		return "host-internal artifact write looks implemented as a gated tool name (D34 note 2/D22)", true
	}
	return "", false
}

// walkText runs a line-based check over all files of dir (all extensions).
// The per-ban counter is what ticket 71 AC#4's ledger reads: ban #6 and ban #7
// each walk exactly one directory, so "files examined for this ban" is a
// well-defined number, and an always-0 one is now visible in the self-report
// instead of being inferred from the absence of findings.
func (s *scanner) walkText(dir, ban string, check func(string) (string, bool), goOnly bool) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "testdata" || d.Name() == "node_modules" || d.Name() == ".git" || s.ign.skip(path, true) {
				return filepath.SkipDir
			}
			return nil
		}
		if goOnly && (!strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go")) {
			return nil
		}
		// The same exclusion ban #8 applies to this very tree (walkEmoji's
		// everyFile branch over frontend/): scan_test.go demands the two counts be
		// one integer, so the ignore rule has to be consulted by both walks or the
		// equality guard goes red on a built worktree.
		if s.ign.skip(path, false) {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		s.examined[ban]++
		for i, line := range strings.Split(string(src), "\n") {
			if what, hit := check(line); hit {
				s.add(ban, path, i+1, what)
			}
		}
		return nil
	})
}

// walkEmoji line-scans one ban #8 scope.
//
// Two scope decisions live here because both were argued in the ticket and
// both are things a later reader will otherwise "fix" the wrong way:
//
//   - COMMENTS ARE EXEMPT AS OF Q-46(c), AND THAT IS A SIGNATURE, NOT AN
//     IMPROVEMENT. The bullet that used to sit here read "COMMENTS COUNT, and
//     this is not a policy choice but a measurement", because the loop matched
//     emojiRe against the RAW line and nothing stripped comments first. Ticket
//     141 put the opposite to the owner and he answered「按推荐」= option (c),
//     "注释豁免、字符串从严", in the same breath as adding the math band
//     U+2200-U+22FF - so the exemption and the widening travel together and
//     neither stands alone: (c) without the band leaves the gate blind to
//     exactly the symbols users see, and the band without (c) is a 78-line
//     comment rewrite. R16#4 forbids an agent narrowing a ban; what happened
//     here is a human authorising it, which is why the reason lives in the
//     ticket and in AGENTS.md's ban list rather than only in this comment, and
//     why the split itself is pinned in scan_test.go (TestBan8CommentExemption*).
//     _test.go still counts: a test source is a string-literal carrier (verdict
//     copy-paste), and Q-46(c) exempts prose, not files.
//   - _test.go COUNTS. Test sources are where a verdict literal gets
//     copy-pasted FROM, and every expensive false green in this repo's recent
//     history came from a gate blind to a whole class of file (ban #1 saw only
//     `go func(`; ban #8 saw only design/). Cost as counted 2026-09-21: 1 line,
//     a comment in internal/tools/bridge_junction_windows_test.go, owned by
//     ticket 20 and deliberately NOT edited here (concurrent same-file writes
//     are fake parallelism); it is registered as the known pending finding in
//     docs/evidence/s1/67-emoji-scope-internal-cmd.md (S7).
//
// goOnly covers the remaining file classes under internal/ and cmd/: the only
// non-.go files there are testdata goldens (the .sse fixtures under
// internal/{agent,llm,models}/testdata) and leaked test-debris directories such
// as internal/tools/tmp/, i.e. not source - and testdata is skipped
// below exactly as walkGo/walkText skip it for every other ban, so ban #8 gains
// no file class the other bans lack and drops none either. design/ stays
// all-text (isTextFile) because its HTML/CSS/JS mockups ARE the surface D23
// governs.
func (s *scanner) walkEmoji(sc emojiScope) error {
	if _, err := os.Stat(sc.dir); os.IsNotExist(err) {
		// Absent tree: the count stays 0 and emptyEmojiScope in main makes
		// that fatal. Not an error here, so Scan stays usable from fixtures
		// that seed only part of the repo.
		return nil
	}
	return filepath.WalkDir(sc.dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "node_modules" || d.Name() == ".git" || s.ign.skip(path, true) {
				return filepath.SkipDir
			}
			// testdata is skipped for the scopes whose sibling bans skip it
			// (goOnly: internal/ + cmd/, everyFile: frontend/, whose walk must
			// match ban #6's rule for=for). design/ keeps walking its testdata if
			// any appears: that tree predates the rule and narrowing it now would
			// drop a file class from a ban this ticket may not shrink.
			if (sc.goOnly || sc.everyFile) && d.Name() == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}
		switch {
		case sc.goOnly:
			if !strings.HasSuffix(path, ".go") {
				return nil
			}
		case sc.everyFile:
			// No filter at all: this is ban #6's walkText shape, chosen so the
			// emoji gate cannot read a smaller tree than the panel-approval gate
			// reads next door. Files that are not text still get their BYTES
			// searched - see the binary note above.
		default:
			if !isTextFile(path) {
				return nil
			}
		}
		// Git-ignored paths are out of ban #8's denominator too, and this is the
		// scope the drift actually bit: `frontend/dist/` is vite output, so on a
		// machine that ran a build the count read 43 where CI reads 40 (A207).
		// The `.gitignore` files themselves and the `frontend/dist/.gitkeep`
		// anchor stay in, because no rule excludes them.
		if s.ign.skip(path, false) {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		s.emojiSeen[sc.label]++
		// Q-46(c): comments are blanked before matching, strings are not. The
		// residue of a line is what a reader or a renderer can see; a comment
		// is what nobody can. blankCommentRanges() returns nil for a file with
		// no comments, so the common path is one map lookup per hit line.
		ranges := commentRangesFor(path, string(src))
		for i, line := range strings.Split(string(src), "\n") {
			probe := line
			if rs := ranges[i+1]; len(rs) > 0 {
				probe = removeRanges(line, rs)
			}
			if emojiRe.MatchString(probe) {
				s.add("emoji", path, i+1, fmt.Sprintf("ban #8 glyph in scope %s is banned (D23): non-comment text, string literals included; comments are exempt per Q-46(c)", sc.label))
			}
		}
		return nil
	})
}

// commentRangesFor returns the byte ranges of one file that are COMMENT text,
// keyed by 1-based line number, half-open [lo,hi) in byte offsets of that line.
//
// This is the "comment vs string" split Q-46(c) was approved for (ticket 141:
// owner answer「按推荐」, so the exemption is signed, not assumed). Two
// classifiers, because the two file classes disagree about what a comment is in
// exactly the place that matters:
//
//   - .go files go through go/ast, because Go's own parser is the only
//     authority on where a comment ends and a string begins. The line
//     `-- L1 slot, <= 20 rows` inside a raw string is a SQL comment and a Go
//     string at once; a lexical `--`/prefix rule would have exempted it, and the
//     owner's own ruling on that exact pair (internal/memory/schema.go:29) is
//     that it counts as a violation. ast cannot be tricked by string contents,
//     so a glyph inside any Go string literal always survives.
//     A .go file that fails to parse yields nil, i.e. NO exemption: a broken
//     file must not buy a free pass.
//   - everything else (design/'s mockups, frontend/'s whole tree) has no parser
//     here, so the rule is deliberately coarse and errs toward reporting: a line
//     is comment text only from a marker that its own non-whitespace prefix
//     carries (`//`, `/*`, `*` while inside an open block, `<!--`, `-->`), and
//     only from that marker onward. Code can therefore share a line with an
//     exempted comment and still be matched. `#` and `--` are NOT markers: the
//     first is a markdown heading and a CSS-independent selector, the second is
//     SQL, and mistaking either for a comment is the blind spot this ticket is
//     closing, not a feature.
func commentRangesFor(path, src string) map[int][][2]int {
	if strings.HasSuffix(path, ".go") {
		return goCommentRanges([]byte(src))
	}
	return textCommentRanges(src)
}

// goCommentRanges blanks every byte go/ast classifies as a comment.
func goCommentRanges(src []byte) map[int][][2]int {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "ban8.go", src, parser.ParseComments)
	if err != nil {
		return nil
	}
	lines := strings.Split(string(src), "\n")
	out := map[int][][2]int{}
	for _, group := range file.Comments {
		for _, c := range group.List {
			start, end := fset.Position(c.Pos()), fset.Position(c.End())
			for ln := start.Line; ln <= end.Line; ln++ {
				if ln < 1 || ln > len(lines) {
					continue
				}
				lo, hi := 0, len(lines[ln-1])
				if ln == start.Line {
					lo = start.Column - 1
				}
				if ln == end.Line {
					hi = end.Column - 1
				}
				if lo < 0 {
					lo = 0
				}
				if hi > lo {
					out[ln] = append(out[ln], [2]int{lo, hi})
				}
			}
		}
	}
	return out
}

// textCommentRanges is the marker-plus-block-state rule described on
// commentRangesFor. open is the kind of block comment spanning the current line
// ("" when none).
func textCommentRanges(src string) map[int][][2]int {
	out := map[int][][2]int{}
	open := ""
	for ln, line := range strings.Split(src, "\n") {
		trimmed := strings.TrimLeft(line, " \t")
		indent := len(line) - len(trimmed)
		mark := func(lo, hi int) {
			if hi > lo {
				out[ln+1] = append(out[ln+1], [2]int{lo, hi})
			}
		}
		if open != "" {
			closer := "*/"
			if open == "<!--" {
				closer = "-->"
			}
			if k := strings.Index(line, closer); k >= 0 {
				mark(0, k+len(closer))
				open = ""
				// Whatever follows a closed block on the same line is code, and
				// code is never exempted, so the line stops being examined here.
				continue
			}
			mark(0, len(line))
			continue
		}
		switch {
		case strings.HasPrefix(trimmed, "//"):
			mark(indent, len(line))
		case strings.HasPrefix(trimmed, "/*"):
			if k := strings.Index(line[indent:], "*/"); k >= 0 {
				mark(indent, indent+k+2)
			} else {
				open = "/*"
				mark(indent, len(line))
			}
		case strings.HasPrefix(trimmed, "<!--"):
			if k := strings.Index(line[indent:], "-->"); k >= 0 {
				mark(indent, indent+k+3)
			} else {
				open = "<!--"
				mark(indent, len(line))
			}
		case strings.HasPrefix(trimmed, "*/"):
			// A line that only closes a block this walk never opened (the block
			// started and ended on one line, or the file was truncated): exempt
			// the closing marker and nothing else.
			mark(indent, indent+2)
		}
	}
	return out
}

// removeRanges returns line with the given byte ranges replaced by spaces, so
// the column of every surviving glyph is unchanged and the finding still points
// at the right place.
func removeRanges(line string, ranges [][2]int) string {
	buf := []byte(line)
	for _, r := range ranges {
		lo, hi := r[0], r[1]
		if lo < 0 {
			lo = 0
		}
		if hi > len(buf) {
			hi = len(buf)
		}
		for i := lo; i < hi && i < len(buf); i++ {
			buf[i] = ' '
		}
	}
	return string(buf)
}

func isTextFile(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".md", ".txt", ".html", ".css", ".js", ".ts", ".tsx", ".jsx", ".json", ".yaml", ".yml", ".toml", ".go", ".svg", ".vue":
		return true
	}
	return false
}

// minProductionGoFiles is how many non-test .go files a real wisp root must
// expose under internal/ + cmd/ for a clean verdict to mean anything. The
// repo is far above this; the floor exists so a mis-pointed -root cannot walk
// an empty tree and report "clean" (same guard shape as
// internal/observe/nobarego_test.go).
const minProductionGoFiles = 10

// checkRoot proves root is a wisp repository the bans can actually be
// evaluated against, and returns the number of production Go files in scope.
// It is the falsifiability guard of ticket 67 AC#2: without it, "0 findings"
// and "the scanner ran nowhere" are indistinguishable on stdout.
func checkRoot(root string) (int, error) {
	for _, required := range []string{
		"go.mod",
		filepath.Join("tools", "d22scan", "allowlist.txt"),
		filepath.Join("internal"),
		filepath.Join("cmd"),
	} {
		if _, err := os.Stat(filepath.Join(root, required)); err != nil {
			return 0, fmt.Errorf("%q is not a wisp repository root: missing %s", root, filepath.ToSlash(required))
		}
	}
	files := 0
	// The same ignore rule the walks use, so this headline number cannot drift
	// from the per-scope denominators it is supposed to corroborate (A207: the
	// whole point of printing the number is that it means the same thing on every
	// machine). Its own matcher means its own index probe (three read-only git
	// calls, cached for its lifetime - see gitignore.go's index()); main() runs one
	// checkRoot and one scan, so a whole run asks git twice, always the same
	// questions, and a disagreement between the two answers would show up as the
	// headline number and the scope counts parting company.
	ign := newGitIgnore(root)
	for _, dir := range []string{filepath.Join(root, "internal"), filepath.Join(root, "cmd")} {
		err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if d.Name() == "testdata" || d.Name() == ".git" || ign.skip(path, true) {
					return filepath.SkipDir
				}
				return nil
			}
			if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") && !ign.skip(path, false) {
				files++
			}
			return nil
		})
		if err != nil {
			return 0, err
		}
	}
	if files < minProductionGoFiles {
		return files, fmt.Errorf("%q has only %d production .go files under internal/ and cmd/, want >= %d: the ban scan would be a no-op",
			root, files, minProductionGoFiles)
	}
	return files, nil
}

// verdict renders the self-report and decides the exit code. It is a function
// instead of inline main() code so the positive-control tests (ticket 71 AC#2's
// "prove the gate can go red") can assert the RED paths - exit 2 on an empty
// live scope, exit 1 on a seeded finding - against captured output without
// exec'ing the binary. main() is a thin argument-parsing shell over it, so a
// test cannot pass by calling a helper main() never uses - which is why the
// ledger is read here through exactly one call, verdict() -> declaredScopes(),
// and the parameterized body is only reachable from a test that names the
// scopes it injects (see fixtureVerdict).
//
// Priority order is deliberate: an empty instrument outranks a finding, and a
// finding outranks the clean sentence, because a 2 is "do not trust anything
// this run printed".
func verdict(out, errOut io.Writer, root string, goFiles int, s *scanner) int {
	return verdictWithScopes(out, errOut, root, goFiles, s, declaredScopes(root))
}

// fixtureScope is injected by the tests ONLY (ticket 88 AC#3): an extra exempt
// scope with a tree a fixture can create at will.
//
// WHY A SEAM EXISTS: ticket 67 registered ban #6 exempt, so the drift guard was
// testable end-to-end by accident - creating frontend/ in a fixture tripped it.
// Ticket 88 armed ban #6 (correctly: the tree is real now), which means
// declaredScopes() contains no exemption left, and "accidental testability" is
// gone with it. The rule still has to be enforced, so it gets tested the way
// every other rule here is: seed a synthetic exempt scope, show the guard stays
// silent while that tree is absent, and show it fires the moment the tree
// appears. main() never calls this - main() goes through verdict(), which passes
// the real ledger - so the fixture cannot make the shipped decision weaker.
type fixtureScope struct{ extra []scanScope }

// fixtureVerdict exists so a test that injects scopes says so in the function
// name it calls; nothing outside _test.go uses it.
func fixtureVerdict(out, errOut io.Writer, root string, goFiles int, s *scanner, fx fixtureScope) int {
	return verdictWithScopes(out, errOut, root, goFiles, s, append(declaredScopes(root), fx.extra...))
}

func verdictWithScopes(out, errOut io.Writer, root string, goFiles int, s *scanner, scopes []scanScope) int {
	emo := emojiScopes(root)

	fmt.Fprintf(out, "d22scan: examined %d production Go files under internal/ and cmd/ of %s\n", goFiles, filepath.ToSlash(root))
	for _, line := range describeScopes(scopes, s) {
		fmt.Fprintln(out, line)
	}
	for _, sc := range uncoveredScopes(scopes) {
		fmt.Fprintf(out, "d22scan: NOT COVERED %s: %s\n", sc.label, sc.note)
	}
	for _, f := range s.findings {
		fmt.Fprintln(out, f.String())
	}

	// Guard 0 (structural): a walk bumped a counter that no scope reports, so
	// this run did work the self-report does not account for.
	if extra := undeclaredKeys(s, scopes); len(extra) > 0 {
		fmt.Fprintf(errOut, "d22scan: %d counter(s) were written by a walk but declared by no scope in declaredScopes(): %s."+
			" Add the scope (and its live/absent policy) - coverage and reporting must move together\n",
			len(extra), strings.Join(extra, ", "))
		return 2
	}
	// Guard 1 (ticket 67's, most specific message first): a declared ban #8
	// scope that line-scanned nothing.
	if empty := emptyEmojiScope(emo, s.emojiSeen); empty != "" {
		fmt.Fprintf(errOut, "d22scan: ban #8 scope %s examined 0 files - it is declared in emojiScopes() but walks nothing."+
			" Point it at a real tree or delete the entry; never leave a scope pretending to scan (ticket 71 AC#4)\n", empty)
		return 2
	}
	// Guard 2 (ticket 71 AC#4's generalization): the same rule for every other
	// live scope, so bans #1-5/#6/#7 cannot walk 0 files and read as green.
	if empty := emptyLiveScope(scopes, s); empty != "" {
		fmt.Fprintf(errOut, "d22scan: scope %s examined 0 files but is declared live in declaredScopes() -"+
			" an empty instrument is not a verdict (ticket 71 AC#4)\n", empty)
		return 2
	}
	// Guard 3: an exemption is a claim about the tree, and claims rot. A scope
	// registered as "its tree is absent" whose tree now exists means the ledger
	// is lying in the other direction (silent 0-coverage became real coverage
	// that nobody armed).
	if drifted := driftedAbsentScope(scopes); drifted != "" {
		fmt.Fprintf(errOut, "d22scan: scope %s is registered as absent-but-exempt while its directory EXISTS."+
			" Flip it to live:true in the same commit that creates the tree, or the scan under-reports coverage\n", drifted)
		return 2
	}
	if len(s.findings) > 0 {
		fmt.Fprintf(errOut, "d22scan: %d finding(s); D22 bans are not negotiable (see PLAN.md D22, tools/d22scan/allowlist.txt)\n", len(s.findings))
		return 1
	}
	// Generated from the ledger + its real counts, so this sentence cannot
	// outlive a coverage change (the old footer hardcoded "design/ or frontend/",
	// and frontend/ is not a tree in this repo while every product .go file went
	// unscanned).
	fmt.Fprintf(out, "d22scan: clean - no D22 ban violations; live scope work: %s; ban #8 emoji coverage: %s",
		scopeSummary(scopes, s), describeEmojiScopes(emo, s.emojiSeen))
	if unc := uncoveredScopes(scopes); len(unc) > 0 {
		names := make([]string, 0, len(unc))
		for _, sc := range unc {
			names = append(names, sc.label)
		}
		fmt.Fprintf(out, "; NOT COVERED: %s", strings.Join(names, ", "))
	}
	fmt.Fprintln(out)
	return 0
}

// scopeSummary renders "label=N" for every live scope, used in the clean line.
func scopeSummary(scopes []scanScope, s *scanner) string {
	parts := make([]string, 0, len(scopes))
	for _, sc := range scopes {
		if sc.live {
			parts = append(parts, fmt.Sprintf("%s=%d", sc.label, sc.count(s)))
		}
	}
	return strings.Join(parts, ", ")
}

func main() {
	root := flag.String("root", ".", "repository root to scan")
	flag.Parse()
	abs, err := filepath.Abs(*root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "d22scan:", err)
		os.Exit(2)
	}
	// A green verdict must be earned by a scan that could see the code.
	files, err := checkRoot(abs)
	if err != nil {
		fmt.Fprintln(os.Stderr, "d22scan:", err)
		fmt.Fprintln(os.Stderr, "d22scan: this tool is its own Go module; run `scripts/d22scan.sh`"+
			" or `cd tools/d22scan && go run . -root ../../` (see tools/d22scan/main.go doc comment)")
		os.Exit(2)
	}
	stats, err := scanWithStats(abs)
	if err != nil {
		fmt.Fprintln(os.Stderr, "d22scan:", err)
		os.Exit(2)
	}
	// Provenance for the numbers below, printed only when the matcher actually
	// excluded something: "frontend/ went 43 -> 40" must never again be a
	// mystery that costs a ticket to explain (A207). Deliberately outside
	// verdict(), whose output is generated from the scanScope ledger and pinned
	// line-for-line by the tests in scan_test.go.
	if n := stats.ign.note(); n != "" {
		fmt.Println(n)
	}
	os.Exit(verdict(os.Stdout, os.Stderr, abs, files, stats))
}
