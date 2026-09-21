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
//	7 internal-artifact-tool  host-internal artifact writes implemented as
//	                      gated tool names (D34 note 2) - scope: internal/tools/
//	8 emoji               zero emoji in design/ (every text file) and in the
//	                      Go sources of internal/ + cmd/ - comments and
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
	"os"
	"path/filepath"
	"regexp"
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
	emojiRe          = regexp.MustCompile(`[\x{1F000}-\x{1FAFF}\x{2600}-\x{27BF}\x{2B00}-\x{2BFF}\x{FE0F}\x{1F1E6}-\x{1F1FF}]`)
	assignKeyShapeRe = regexp.MustCompile(`^[A-Za-z0-9._~-]{16,}$`)
)

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
	s := &scanner{root: root, allow: map[string]map[string]bool{}, emojiSeen: map[string]int{}}
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
}

// emojiScopes is ban #8's ENTIRE coverage, kept in one place so the tool cannot
// print a footer naming a tree it never walked - the failure shape of A22 /
// A26 / A30, and exactly what this list was before ticket 67 AC#3 landed:
// design/ plus frontend/, and frontend/ does not exist at this HEAD, so the
// ban was blind to all 310 .go files of the product while printing "no emoji in
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
func emojiScopes(root string) []emojiScope {
	return []emojiScope{
		{dir: filepath.Join(root, "design"), label: "design/"},
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

// walkGo runs the Go-level bans over production files of dir.
func (s *scanner) walkGo(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "testdata" || d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
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
func (s *scanner) walkText(dir, ban string, check func(string) (string, bool), goOnly bool) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "testdata" || d.Name() == "node_modules" || d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if goOnly && (!strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go")) {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
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
//   - COMMENTS COUNT, and this is not a policy choice but a measurement: the
//     loop below matches emojiRe against the RAW line and nothing strips
//     comments first, so a glyph in prose is already a violation in this
//     instrument's semantics. Ticket 67 AC#3 proposed the opposite ("comments
//     are not user-visible, exclude them"); implementing that needs new
//     comment-stripping code, i.e. a coverage NARROWING, which R16#4 forbids.
//     The measured blast radius of keeping comments in was 3 production lines
//     (internal/llm/probe_health.go:17/:120/:203, all comments), and they were
//     cleaned to ASCII PASS/FAIL rather than exempted - see
//     docs/evidence/s1/67-emoji-scope-internal-cmd.md. Attribution caveat on
//     purpose: ban #8 is D23's design-language ban, NOT a console-encoding
//     check; the encoding argument was specific to cmd/wisp's verdict column.
//   - _test.go COUNTS. Test sources are where a verdict literal gets
//     copy-pasted FROM, and every expensive false green in this repo's recent
//     history came from a gate blind to a whole class of file (ban #1 saw only
//     `go func(`; ban #8 saw only design/). Cost measured for this change:
//     1 line, in internal/tools/bridge_junction_windows_test.go:444, owned by
//     ticket 20 and deliberately NOT edited here (concurrent same-file writes
//     are fake parallelism); it is registered as the known pending finding.
//
// goOnly covers the remaining file classes under internal/ and cmd/: the only
// non-.go files there are testdata goldens (57 .sse fixtures) and leaked test
// debris under internal/tools/tmp/, i.e. not source - and testdata is skipped
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
			if d.Name() == "node_modules" || d.Name() == ".git" {
				return filepath.SkipDir
			}
			if sc.goOnly && d.Name() == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}
		if sc.goOnly {
			if !strings.HasSuffix(path, ".go") {
				return nil
			}
		} else if !isTextFile(path) {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		s.emojiSeen[sc.label]++
		for i, line := range strings.Split(string(src), "\n") {
			if emojiRe.MatchString(line) {
				s.add("emoji", path, i+1, fmt.Sprintf("ban #8 glyph in scope %s is banned (D23): covers comments and _test.go, not only string literals", sc.label))
			}
		}
		return nil
	})
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
	for _, dir := range []string{filepath.Join(root, "internal"), filepath.Join(root, "cmd")} {
		err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if d.Name() == "testdata" || d.Name() == ".git" {
					return filepath.SkipDir
				}
				return nil
			}
			if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
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
	findings := stats.findings
	fmt.Printf("d22scan: examined %d production Go files under internal/ and cmd/ of %s\n", files, filepath.ToSlash(abs))
	scopes := emojiScopes(abs)
	for _, sc := range scopes {
		fmt.Printf("d22scan: ban #8 scope %-10s examined %d file(s)\n", sc.label, stats.emojiSeen[sc.label])
	}
	for _, f := range findings {
		fmt.Println(f.String())
	}
	// A declared scope that walked nothing is an empty instrument, not a
	// verdict: louder and higher-priority than "N findings" (ticket 71 AC#4).
	if empty := emptyEmojiScope(scopes, stats.emojiSeen); empty != "" {
		fmt.Fprintf(os.Stderr, "d22scan: ban #8 scope %s examined 0 files - it is declared in emojiScopes() but walks nothing."+
			" Point it at a real tree or delete the entry; never leave a scope pretending to scan (ticket 71 AC#4)\n", empty)
		os.Exit(2)
	}
	if len(findings) > 0 {
		fmt.Fprintf(os.Stderr, "d22scan: %d finding(s); D22 bans are not negotiable (see PLAN.md D22, tools/d22scan/allowlist.txt)\n", len(findings))
		os.Exit(1)
	}
	// Generated from the scope list + its real counts, so this sentence cannot
	// outlive a coverage change (the old footer claimed frontend/, which is not
	// a tree in this repo, while all 310 product .go files went unscanned).
	fmt.Println("d22scan: clean - no D22 ban violations; ban #8 emoji coverage: " + describeEmojiScopes(scopes, stats.emojiSeen))
}
