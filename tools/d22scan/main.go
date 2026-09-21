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
//	8 emoji               zero emoji in design/ and frontend/ (D23)
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
	failAddOn string // unused placeholder guard
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
	s := &scanner{root: root, allow: map[string]map[string]bool{}}
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
	if err := s.walkEmoji(filepath.Join(root, "design")); err != nil {
		return nil, err
	}
	if err := s.walkEmoji(filepath.Join(root, "frontend")); err != nil {
		return nil, err
	}
	return s.findings, nil
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

func (s *scanner) walkEmoji(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "node_modules" || d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if !isTextFile(path) {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for i, line := range strings.Split(string(src), "\n") {
			if emojiRe.MatchString(line) {
				s.add("emoji", path, i+1, "emoji is banned in design/ and frontend/ (D23)")
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
	findings, err := Scan(abs)
	if err != nil {
		fmt.Fprintln(os.Stderr, "d22scan:", err)
		os.Exit(2)
	}
	fmt.Printf("d22scan: examined %d production Go files under internal/ and cmd/ of %s\n", files, filepath.ToSlash(abs))
	for _, f := range findings {
		fmt.Println(f.String())
	}
	if len(findings) > 0 {
		fmt.Fprintf(os.Stderr, "d22scan: %d finding(s); D22 bans are not negotiable (see PLAN.md D22, tools/d22scan/allowlist.txt)\n", len(findings))
		os.Exit(1)
	}
	fmt.Println("d22scan: clean - no D22 ban violations, no emoji in design/ or frontend/")
}
