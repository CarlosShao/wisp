package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func seedFile(t *testing.T, root, rel, content string) {
	t.Helper()
	full := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestScanDetectsAllSeededViolations is the lint self-check demanded by the
// ticket: a seeded violation per ban must turn the scan red (locally
// reproducible red test for the CI lint job).
func TestScanDetectsAllSeededViolations(t *testing.T) {
	root := t.TempDir()
	seedFile(t, root, "tools/d22scan/allowlist.txt", "# empty allowlist for the fixture\n")

	// Two seeds, not one: R16#1 widened ban #1 to every `go <anything>`, so the
	// positive control has to prove the NAMED form is caught as well. A gate
	// that only ever fires on `go func(` is exactly the A26 hole being closed -
	// if the named matcher regressed to "invisible", this fixture would still
	// look green.
	seedFile(t, root, "internal/bad/goroutine.go", `package bad

func leak() {
	go func() { println("unnamed") }()
}

func worker() {}

func leakNamed() {
	go worker()
}
`)
	seedFile(t, root, "internal/bad/paths.go", `package bad

import "path/filepath"

func Normalize(p string) string {
	return filepath.Clean(p)
}
`)
	seedFile(t, root, "internal/bad/key.go", `package bad

const apiKey = "abcd1234efgh5678ijkl"
`)
	seedFile(t, root, "internal/bad/wallclock.go", `package bad

import "time"

func deadline() time.Time {
	start := time.Now()
	_ = start
	return start.Add(time.Minute - start.Sub(time.Now()))
}
`)
	seedFile(t, root, "internal/bad/mirror.go", `package bad

const mirrorSHA256URL = "https://mirror.example/file.sha256"
`)
	seedFile(t, root, "internal/tools/artifact.go", `package tools

var names = []string{"spill", "internal.logwrite"}
`)
	seedFile(t, root, "frontend/src/app.js", "export function decide() { return approval.decide({allow: true}); }\n")
	seedFile(t, root, "design/screens/ball.md", "# Ball\n\nsmile more often \U0001F600\n")

	findings, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]int{}
	for _, f := range findings {
		got[f.Ban]++
	}
	want := map[string]int{
		"bare-goroutine":         2,
		"pathresolver-bypass":    1,
		"plaintext-key":          1,
		"wallclock-timeout":      1,
		"mirror-hash":            1,
		"internal-artifact-tool": 1,
		"panel-approval":         1,
		"emoji":                  1,
	}
	for ban, n := range want {
		if got[ban] != n {
			t.Errorf("ban %q: want %d finding(s), got %d (all: %v)", ban, n, got[ban], findings)
		}
		delete(got, ban)
	}
	for ban := range got {
		t.Errorf("unexpected findings for ban %q", ban)
	}
}

func TestScanCleanRepoIsGreen(t *testing.T) {
	root := t.TempDir()
	seedFile(t, root, "tools/d22scan/allowlist.txt", "# empty\n")
	seedFile(t, root, "internal/ok/ok.go", `package ok

import (
	"context"
	"time"
)

// Start runs work via the registry; clock math is monotonic.
func Start(ctx context.Context) {
	deadline := time.Now().Add(time.Second)
	_ = deadline
	_ = ctx
}
`)
	seedFile(t, root, "design/screens/ball.md", "# Ball states\n\n- idle glow\n")
	findings, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("clean fixture must be green, got: %v", findings)
	}
}

func TestAllowlistSuppressesOnlyListedPaths(t *testing.T) {
	root := t.TempDir()
	al := "pathresolver-bypass\tinternal/bad/paths.go\tfixture reason\n"
	seedFile(t, root, "tools/d22scan/allowlist.txt", al)
	seedFile(t, root, "internal/bad/paths.go", `package bad

import "path/filepath"

func Normalize(p string) string {
	return filepath.Clean(p)
}
`)
	seedFile(t, root, "internal/bad/paths2.go", `package bad

import "path/filepath"

func Normalize2(p string) string {
	return filepath.Abs(p)
}
`)
	findings, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 {
		t.Fatalf("allowlist must suppress only the listed file, got: %v", findings)
	}
	if findings[0].Ban != "pathresolver-bypass" || !strings.HasSuffix(findings[0].Path, "paths2.go") {
		t.Fatalf("wrong surviving finding: %+v", findings[0])
	}
}

// TestCheckRootRejectsBlindRoots is the ticket 67 AC#2 guard: a -root that
// cannot contain any banned pattern must be a loud error, never a clean verdict.
func TestCheckRootRejectsBlindRoots(t *testing.T) {
	cases := []struct {
		name string
		seed func(t *testing.T, root string)
	}{
		{"empty dir", func(_ *testing.T, _ string) {}},
		{"go.mod only", func(t *testing.T, root string) {
			seedFile(t, root, "go.mod", "module x\n")
		}},
		// The exact shape of `cd tools/d22scan && go run .` with the default
		// -root .: the module has its own go.mod and allowlist, so a naive
		// existence check passes, but it holds no internal/ and no cmd/.
		{"the d22scan module itself", func(t *testing.T, root string) {
			seedFile(t, root, "go.mod", "module x\n")
			seedFile(t, root, "tools/d22scan/allowlist.txt", "# empty\n")
			seedFile(t, root, "main.go", "package main\n")
		}},
		// internal/ + cmd/ present but too thin to be the product tree.
		{"repo skeleton with two go files", func(t *testing.T, root string) {
			seedFile(t, root, "go.mod", "module x\n")
			seedFile(t, root, "tools/d22scan/allowlist.txt", "# empty\n")
			seedFile(t, root, "internal/a/a.go", "package a\n")
			seedFile(t, root, "cmd/b/b.go", "package b\n")
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			root := t.TempDir()
			c.seed(t, root)
			files, err := checkRoot(root)
			if err == nil {
				t.Fatalf("mis-invocation accepted: checkRoot(%q) reported %d files, want an error", root, files)
			}
			if files >= minProductionGoFiles {
				t.Errorf("error reported but file count %d already clears the floor", files)
			}
		})
	}
}

// TestScanAloneIsNotAFalsifier pins WHY checkRoot exists: Scan() on a root with
// nothing to scan returns zero findings, i.e. the old "clean" was compatible
// with "the scanner looked at nothing". The positive control must be the pair
// (checkRoot ok, Scan green), never Scan by itself.
func TestScanAloneIsNotAFalsifier(t *testing.T) {
	root := t.TempDir()
	seedFile(t, root, "tools/d22scan/allowlist.txt", "# empty\n")
	seedFile(t, root, "go.mod", "module x\n")
	findings, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("fixture unexpectedly found: %v", findings)
	}
	if _, err := checkRoot(root); err == nil {
		t.Fatal("checkRoot must reject the root that Scan just declared clean")
	}
}

// TestCheckRootAcceptsRealRepo keeps the guard from being a tautology: the same
// predicate that rejects the blind roots must pass on the tree it guards.
func TestCheckRootAcceptsRealRepo(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Skipf("not inside the wisp repo: %v", err)
	}
	files, err := checkRoot(root)
	if err != nil {
		t.Fatalf("real repo root rejected: %v", err)
	}
	if files < minProductionGoFiles {
		t.Fatalf("real repo exposes only %d production .go files", files)
	}
	t.Logf("real repo production Go files in scope: %d", files)
}

func TestScannerSelfScanOfRealRepoIsGreen(t *testing.T) {
	// Walk up from the package dir to the real repo root (the test lives at
	// tools/d22scan/) and scan it: the lint job must be green on HEAD.
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Skipf("not inside the wisp repo: %v", err)
	}
	findings, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		for _, f := range findings {
			t.Errorf("repo HEAD violates: %s", f.String())
		}
	}
}

// TestEmojiBanCoversGoSourcesNotJustDesign is registry A23 judgement ① pinned
// in code: ban #8 must fire inside internal/ AND inside cmd/, on a non-test file
// and in a _test.go file. Before ticket 67 AC#3's coverage step all four of
// these seeds walked past the gate unnoticed, because the only declared scopes
// were design/ and frontend/ - and frontend/ does not exist in this repository,
// so every .go file in the product was invisible to the ban while the footer
// still claimed "no emoji in design/ or frontend/".
//
// Q-46(c) (ticket 141) moved the glyphs in these seeds out of comments and into
// string literals, because comments are exempt now and a comment seed would
// prove nothing about coverage. The exemption has its own two tests below
// (TestBan8CommentExemptionInGoSources, TestBan8CommentExemptionInTextScopes);
// deleting the coverage seeds instead of moving them would reopen A23.
func TestEmojiBanCoversGoSourcesNotJustDesign(t *testing.T) {
	root := t.TempDir()
	seedFile(t, root, "tools/d22scan/allowlist.txt", "# empty\n")

	// U+2713 CHECK MARK: inside emojiRe's \x{2600}-\x{27BF} band, and the exact
	// code point that lived in cmd/wisp/providers.go until fff4cad.
	seedFile(t, root, "internal/ok/decl.go", "package ok\n\nconst doc = \"a verdict reads \u2713 here\"\n")
	seedFile(t, root, "internal/ok/literal.go", "package ok\n\nconst banner = \"ready \u2713\"\n")
	seedFile(t, root, "internal/ok/emoji_test.go", "package ok\n\nvar want = \"\u2713 in a test file\"\nfunc TestNothing(*testing.T) {}\n")
	seedFile(t, root, "cmd/wisp/glyph.go", "package main\n\nconst tail = \"\u2713 in cmd/\"\nfunc main() {}\n")

	// Control for the documented goOnly rule: a non-.go file under internal/ is
	// testdata goldens / leaked test debris, not source, so it is out of scope.
	seedFile(t, root, "internal/ok/testdata/fixture.txt", "ready \u2713\n")

	findings, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"internal/ok/decl.go",       // string literal in a non-test file, internal/
		"internal/ok/literal.go",    // string literal, internal/
		"internal/ok/emoji_test.go", // _test.go
		"cmd/wisp/glyph.go",         // cmd/
	}
	got := map[string]int{}
	for _, f := range findings {
		if f.Ban != "emoji" {
			t.Errorf("unexpected non-emoji finding from the fixture: %s", f.String())
			continue
		}
		got[f.Path]++
	}
	for _, p := range want {
		if got[p] != 1 {
			t.Errorf("ban #8 must report %s exactly once, got %d (all findings: %v)", p, got[p], findings)
		}
		delete(got, p)
	}
	for p := range got {
		if strings.HasPrefix(p, "internal/ok/testdata/") {
			t.Errorf("ban #8 reported the goOnly control %s - if that is the NEW intended scope, delete this control and its comment, do not delete the coverage", p)
			continue
		}
		t.Errorf("unexpected emoji finding for %s (all: %v)", p, findings)
	}
}

// TestBan8MathBandAndRemainingGaps pins emojiRe's character class by code point,
// which is the "pin any change there in scan_test.go" that main.go's ban #8
// header asks for (the generated footer cannot pin a range).
//
// Two directions are pinned on purpose. The positive one is Q-46(c)'s addition:
// U+2265/U+2264/U+2229 live in \x{2200}-\x{22FF}, the band ticket 141 added
// because PLAN.md's own range list does not contain it either while 5 real
// non-comment glyphs do (the approval card's ">=50", the panel's "<= 64.0 MB").
// Removing that band from the class must make these rows fail, which is AC#3's
// reverse proof: the seed is caught by the new band, not by a neighbouring one.
//
// The negative rows are the gaps the approval did NOT close, recorded so a
// later reader cannot mistake them for coverage: U+2192 and U+21D2 (arrow band)
// and U+2460/U+2461 (circled numbers) are still unscanned, which is why the
// plan-literal blast radius (141 lines) and the shipped one (9) differ by the
// 78 comment lines Q-46(c) exempts plus these bands nobody authorised.
func TestBan8MathBandAndRemainingGaps(t *testing.T) {
	cases := []struct {
		name      string
		glyph     string
		wantFired bool
	}{
		{"U+2265 greater-or-equal, band added by Q-46(c)", "\u2265", true},
		{"U+2264 less-or-equal, band added by Q-46(c)", "\u2264", true},
		{"U+2229 intersection, band added by Q-46(c)", "\u2229", true},
		{"U+2212 minus sign, band added by Q-46(c)", "\u2212", true},
		{"U+2713 check, pre-existing U+2600-U+27BF band", "\u2713", true},
		{"U+1F600 emoji, pre-existing U+1F000-U+1FAFF band", "\U0001F600", true},
		{"U+FE0F variation selector, pre-existing", "\uFE0F", true},
		{"U+2192 right arrow, gap kept", "\u2192", false},
		{"U+21D2 rightdouble arrow, gap kept", "\u21D2", false},
		{"U+2500 box drawing, gap kept", "\u2500", false},
		{"U+2460 circled one, gap kept", "\u2460", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			seedFile(t, root, "tools/d22scan/allowlist.txt", "# empty\n")
			seedFile(t, root, "internal/ok/glyph.go", "package ok\n\nconst s = \"value "+tc.glyph+" here\"\n")
			findings, err := Scan(root)
			if err != nil {
				t.Fatal(err)
			}
			fired := false
			for _, f := range findings {
				if f.Ban != "emoji" {
					t.Errorf("unexpected non-emoji finding: %s", f.String())
					continue
				}
				if strings.HasSuffix(filepath.ToSlash(f.Path), "internal/ok/glyph.go") {
					fired = true
				}
			}
			if fired != tc.wantFired {
				t.Errorf("U+%X in a string literal: wantFired=%v got %v (findings %v)", []rune(tc.glyph)[0], tc.wantFired, fired, findings)
			}
		})
	}
}

// TestBan8CommentExemptionInGoSources is the reading Q-46(c) was signed for,
// and the reason it is not decoration: the SAME glyph in the SAME file goes red
// as a string and green as a comment. Every exempt row here was a finding on
// HEAD before ticket 141 (measured: 78 such lines in 35 files under internal/
// and cmd/ at the plan-literal range), so deleting the exemption re-drives this
// test red rather than leaving it vacuously true.
//
// The rows that must STAY red are the half that keeps the gate honest: a raw
// string whose SQL `--` prose looks like a comment (the exact shape of
// internal/memory/schema.go:29, which the owner ruled a violation, not an
// exemption), a string sharing a line with a comment, and an unparseable .go
// file, where failing to locate the comments buys no exemption at all.
func TestBan8CommentExemptionInGoSources(t *testing.T) {
	cases := []struct {
		name         string
		src          string
		wantFindings int
	}{
		{
			name:         "line comment",
			src:          "package ok\n\n// a verdict reads \u2713 here\nfunc f() {}\n",
			wantFindings: 0,
		},
		{
			name:         "doc comment above a declaration",
			src:          "package ok\n\n// F returns \u2265 when satisfied.\nfunc F() {}\n",
			wantFindings: 0,
		},
		{
			name:         "one-line block comment",
			src:          "package ok\n\nfunc f() { /* \u2192 and \u2264 */ }\n",
			wantFindings: 0,
		},
		{
			name:         "multi-line block comment, glyph on a middle line",
			src:          "package ok\n\n/*\nheader\nline with \u2229 in the middle\nfooter\n*/\nfunc f() {}\n",
			wantFindings: 0,
		},
		{
			name:         "trailing comment after clean code",
			src:          "package ok\n\nvar x = 1 // \u2713\n",
			wantFindings: 0,
		},
		{
			name:         "string on the same line as a trailing comment",
			src:          "package ok\n\nvar x = \"\u2713\" // \u2713\n",
			wantFindings: 1,
		},
		{
			name:         "raw string that reads as a SQL comment",
			src:          "package ok\n\nconst ddl = `-- L1 slot\uFF0c \u2264 20 xing\nCREATE TABLE t (id INTEGER);\n`\n",
			wantFindings: 1,
		},
		{
			name:         "interpreted string containing a comment marker",
			src:          "package ok\n\nconst url = \"http://example.invalid/x\u2713\"\n",
			wantFindings: 1,
		},
		{
			name:         "glyph inside a comment AND inside code in one block comment",
			src:          "package ok\n\n/*\nprose \u2713\n*/\nconst s = \"\u2713\"\n",
			wantFindings: 1,
		},
		{
			name:         "unparseable file gets no exemption at all",
			src:          "package ok\n\nfunc broken( {\n// \u2713\n\nconst s = \"\u2713\"\n",
			wantFindings: 2,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			seedFile(t, root, "tools/d22scan/allowlist.txt", "# empty\n")
			seedFile(t, root, "internal/ok/case.go", tc.src)
			findings, err := Scan(root)
			if err != nil {
				t.Fatal(err)
			}
			n := 0
			for _, f := range findings {
				if f.Ban == "unparseable" {
					// walkGo reports a .go file it cannot parse under its own
					// ban; this fixture has exactly one such row, and it is the
					// point of that row, not noise.
					continue
				}
				if f.Ban != "emoji" {
					t.Errorf("unexpected non-emoji finding: %s", f.String())
					continue
				}
				n++
			}
			if n != tc.wantFindings {
				t.Errorf("ban #8 fired %d time(s), want %d for src %q", n, tc.wantFindings, tc.src)
			}
		})
	}
}

// TestBan8CommentExemptionInTextScopes is the same split for the two scopes that
// are not Go: design/ (all text files) and frontend/ (every file). It is pinned
// separately because the non-Go classifier is a different rule with a different
// failure mode - main.go's textCommentRanges only exempts a line from a marker
// its own indentation-free prefix carries, so `#` (a markdown heading, and the
// seed glyph of TestScanDetectsAllSeededViolations) and `--` (SQL) are NOT
// markers, and code after a closed block comment is still examined.
//
// The frontend rows matter for ticket 141's acceptance criterion 3: 19 of the
// 25 glyph lines the widened ban walks there are box-drawing section dividers
// and they must read green, while the 6 that carry U+2264/U+2212 in rendered JSX
// and fixture text must read red.
func TestBan8CommentExemptionInTextScopes(t *testing.T) {
	cases := []struct {
		name         string
		rel          string
		src          string
		wantFindings int
	}{
		{
			// U+2713 and U+2264 are chosen over the U+2500 these files really
			// carry: box drawing sits in a band the approval did NOT add, so a
			// divider seed here would pass whether or not the exemption works.
			name:         "tsx line comment",
			rel:          "frontend/src/a.tsx",
			src:          "// section \u2713\u2713\u2713\nexport const A = () => null;\n",
			wantFindings: 0,
		},
		{
			name:         "tsx block divider, glyph on the closing line",
			rel:          "frontend/src/b.tsx",
			src:          "/* \u2713\u2713\u2713\n * more \u2264\u2264\u2264\n */\nexport const B = () => null;\n",
			wantFindings: 0,
		},
		{
			name:         "html comment",
			rel:          "frontend/fixtures/c.html",
			src:          "<!-- \u2264 not rendered -->\n<div>ok</div>\n",
			wantFindings: 0,
		},
		{
			name:         "css comment",
			rel:          "frontend/src/d.css",
			src:          "/* \u2713 divider */\n.a { color: red }\n",
			wantFindings: 0,
		},
		{
			name:         "JSX text node with a math glyph",
			rel:          "frontend/src/e.tsx",
			src:          "export const E = () => <span>\u2264 64.0 MB</span>;\n",
			wantFindings: 1,
		},
		{
			name:         "code after a closed block comment on the same line",
			rel:          "frontend/src/f.tsx",
			src:          "/* \u2713 */ export const F = \"\u2713\";\n",
			wantFindings: 1,
		},
		{
			name:         "markdown heading is not a comment",
			rel:          "design/screens/g.md",
			src:          "# Ball \u2713\n",
			wantFindings: 1,
		},
		{
			name:         "SQL-style dash is not a comment in a text scope",
			rel:          "design/screens/h.txt",
			src:          "-- \u2713\n",
			wantFindings: 1,
		},
		{
			name:         "extension-less file, plain text line",
			rel:          "frontend/Procfile",
			src:          "web: node server.js \u2713\n",
			wantFindings: 1,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			seedFile(t, root, "tools/d22scan/allowlist.txt", "# empty\n")
			seedFile(t, root, "internal/ok/ok.go", "package ok\n")
			seedFile(t, root, "cmd/wisp/main.go", "package main\n\nfunc main() {}\n")
			seedFile(t, root, tc.rel, tc.src)
			findings, err := Scan(root)
			if err != nil {
				t.Fatal(err)
			}
			n := 0
			for _, f := range findings {
				if f.Ban != "emoji" {
					t.Errorf("unexpected non-emoji finding: %s", f.String())
					continue
				}
				if !strings.HasSuffix(f.Path, tc.rel) {
					t.Errorf("finding %s is not in the seeded file %s", f.Path, tc.rel)
				}
				n++
			}
			if n != tc.wantFindings {
				t.Errorf("ban #8 fired %d time(s) in %s, want %d for src %q", n, tc.rel, tc.wantFindings, tc.src)
			}
		})
	}
}

// TestDeclaredEmojiScopeCannotWalkZeroFiles pins the disposition of frontend/:
// a scope named in emojiScopes() that line-scans no files is a loud failure,
// not a green run. Deleting the dead entry is what ticket 71 AC#4 allows;
// leaving it to walk nothing while the footer advertises it is what this guard
// makes impossible to repeat.
func TestDeclaredEmojiScopeCannotWalkZeroFiles(t *testing.T) {
	root := t.TempDir()
	seedFile(t, root, "tools/d22scan/allowlist.txt", "# empty\n")
	seedFile(t, root, "design/screens/ball.md", "# Ball\n\n- idle glow\n")
	seedFile(t, root, "internal/ok/ok.go", "package ok\n")
	// Ticket 96: frontend/ is a declared scope now, so this fixture has to seed
	// it too - otherwise "cmd/ is the empty one" would be decided by which empty
	// scope happens to come first in the list, and cmd/ is the tree this test was
	// written to prove the guard names. The second half below removes frontend/
	// instead, so the guard is shown firing for the panel's tree as well and the
	// rule cannot be satisfied by ordering luck.
	seedFile(t, root, "frontend/src/panel.tsx", "export const Panel = () => null;\n")
	// cmd/ deliberately absent.

	s, err := scanWithStats(root)
	if err != nil {
		t.Fatal(err)
	}
	if s.emojiSeen["design/"] == 0 || s.emojiSeen["internal/"] == 0 || s.emojiSeen["frontend/"] == 0 {
		t.Fatalf("fixture broken: scopes backed by files must report work, got %v", s.emojiSeen)
	}
	if s.emojiSeen["cmd/"] != 0 {
		t.Fatalf("fixture broken: cmd/ has no files, saw %d", s.emojiSeen["cmd/"])
	}
	if len(s.findings) != 0 {
		t.Fatalf("fixture must be finding-free so the guard is what trips, got %v", s.findings)
	}
	if got := emptyEmojiScope(emojiScopes(root), s.emojiSeen); got != "cmd/" {
		t.Fatalf("the empty declared scope must be named, got %q (counts %v)", got, s.emojiSeen)
	}

	if err := os.RemoveAll(filepath.Join(root, "frontend")); err != nil {
		t.Fatal(err)
	}
	s2, err := scanWithStats(root)
	if err != nil {
		t.Fatal(err)
	}
	if s2.emojiSeen["frontend/"] != 0 {
		t.Fatalf("fixture broken: frontend/ was deleted, saw %d", s2.emojiSeen["frontend/"])
	}
	seedFile(t, root, "cmd/wisp/main.go", "package main\n\nfunc main() {}\n")
	s3, err := scanWithStats(root)
	if err != nil {
		t.Fatal(err)
	}
	if s3.emojiSeen["cmd/"] == 0 {
		t.Fatalf("fixture broken: cmd/ was seeded, saw %d", s3.emojiSeen["cmd/"])
	}
	if got := emptyEmojiScope(emojiScopes(root), s3.emojiSeen); got != "frontend/" {
		t.Fatalf("the panel's tree must be subject to the same empty-scope rule, got %q (counts %v)", got, s3.emojiSeen)
	}
	out, errOut, code := runVerdict(t, root, s3)
	if code != 2 {
		t.Fatalf("an empty ban #8 scope must exit 2, got rc=%d out=%s err=%s", code, out, errOut)
	}
	if !strings.Contains(errOut, "ban #8 scope frontend/") {
		t.Errorf("guard must name ban #8's frontend/ scope, got %q", errOut)
	}
}

// TestScopeReportMatchesRealCoverage is the "覆盖面与话术必须一致" pin: the green
// footer text is GENERATED from emojiScopes() and its real counts, so it can
// neither advertise a tree the scan does not walk (the old frontend/ claim) nor
// silently drop one it does. This runs against the real repository on purpose -
// the claim under test is about this repo's actual coverage.
func TestScopeReportMatchesRealCoverage(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Skipf("not inside the wisp repo: %v", err)
	}
	s, err := scanWithStats(root)
	if err != nil {
		t.Fatal(err)
	}
	scopes := emojiScopes(root)
	for _, sc := range scopes {
		if s.emojiSeen[sc.label] == 0 {
			t.Errorf("ban #8 scope %s walks 0 files in the real repo: it is a dead entry, give it coverage or delete it (ticket 71 AC#4)", sc.label)
		}
	}
	if got := emptyEmojiScope(scopes, s.emojiSeen); got != "" {
		t.Errorf("real repo must have no empty ban #8 scope, got %q", got)
	}
	report := describeEmojiScopes(scopes, s.emojiSeen)
	for _, sc := range scopes {
		if !strings.Contains(report, sc.label) {
			t.Errorf("the clean verdict line would not name scope %s it actually scanned: %q", sc.label, report)
		}
	}
	// RESTATED BY TICKET 96 (AC#3 + AC#4). This assertion used to read
	// `if strings.Contains(report, "frontend") { t.Errorf(...) }`, because at
	// ticket 71's HEAD emojiScopes() had deliberately NO frontend/ entry and a
	// footer naming it would have been a lie. The entry exists now and walks the
	// panel's files, so the claim is the other way round: a report that does NOT
	// name frontend/ is the lie. Why this has to be a literal string at all, and
	// cannot be left to the loop just above: that loop derives BOTH sides from
	// emojiScopes(), so it is self-consistent and can never notice a missing
	// tree - which is precisely ticket 96's bug shape (the list was never
	// synchronised and every internal check stayed green). Only this literal
	// "frontend/" catches "somebody dropped the entry again".
	if !strings.Contains(report, "frontend/") {
		t.Errorf("the ban #8 report does not name frontend/ at all: emojiScopes() must carry the panel's tree (ticket 96 AC#3) - %q", report)
	}
}

// ---------------------------------------------------------------------------
// ticket 71 AC#4 + AC#2: the per-scope ledger, and a RED proof for every guard
// on it. AC#2's rule is that a guard nobody has seen fail is not a guard, so
// each of the four checks has its own seeded-red test instead of one shared
// happy path.
// ---------------------------------------------------------------------------

// liveFixture is a repository shape on which EVERY live scope does real work, so
// a guard tripping in a test built on it can only be the thing under test.
func liveFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	seedFile(t, root, "tools/d22scan/allowlist.txt", "# empty\n")
	seedFile(t, root, "go.mod", "module x\n")
	seedFile(t, root, "internal/ok/ok.go", "package ok\n")
	seedFile(t, root, "internal/tools/ok.go", "package tools\n")
	seedFile(t, root, "cmd/wisp/main.go", "package main\n\nfunc main() {}\n")
	seedFile(t, root, "design/index.html", "<html>ok</html>\n")
	// Ticket 88: ban #6 became live, so "every live scope does real work" now
	// includes frontend/. The file is a .tsx on purpose - that is the extension
	// ticket 88 AC#1 asked about and the one the panel is written in, so a
	// regression that narrowed walkText's filter to .js would turn every test
	// built on this fixture red instead of quietly re-creating an empty
	// instrument. The content is clean: ban #6 must be able to examine files and
	// still be green (a live scope is not a violation).
	seedFile(t, root, "frontend/src/panel.tsx", "export const Panel = () => null;\n")
	return root
}

func scanFixture(t *testing.T, root string) *scanner {
	t.Helper()
	s, err := scanWithStats(root)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// runVerdict calls verdict directly rather than exec'ing the binary: the thing
// under test is the decision, and the subprocess wiring (main turning a returned
// code into a process exit code) has its own end-to-end test below,
// TestBuiltBinaryGoesRedEndToEnd, which is the one CI's step actually depends on.
func runVerdict(t *testing.T, root string, s *scanner) (string, string, int) {
	t.Helper()
	var out, errOut strings.Builder
	code := verdict(&out, &errOut, root, 99, s)
	return out.String(), errOut.String(), code
}

// runInjectedVerdict is runVerdict for the ONE rule that the shipped ledger can
// no longer reach: exemption drift. It calls fixtureVerdict, which appends the
// given scopes to declaredScopes() and nothing else, so the real ledger is still
// what produces every other verdict in this file (see fixtureScope's comment for
// why ticket 88 had to open this seam instead of leaning on ban #6 again).
func runInjectedVerdict(t *testing.T, root string, s *scanner, extra ...scanScope) (string, string, int) {
	t.Helper()
	var out, errOut strings.Builder
	code := fixtureVerdict(&out, &errOut, root, 99, s, fixtureScope{extra: extra})
	return out.String(), errOut.String(), code
}

// labels renders scope labels for failure messages.
func labels(scopes []scanScope) []string {
	out := make([]string, 0, len(scopes))
	for _, sc := range scopes {
		out = append(out, sc.label)
	}
	return out
}

// TestVerdictGreenOnFullyLiveFixture is the control the red tests below are
// measured against: with every live scope populated the scan exits 0 and prints
// exactly one work line per declared scope.
func TestVerdictGreenOnFullyLiveFixture(t *testing.T) {
	root := liveFixture(t)
	s := scanFixture(t, root)
	if len(s.findings) != 0 {
		t.Fatalf("fixture must be finding-free so the exit code means what it says: %v", s.findings)
	}
	scopes := declaredScopes(root)
	out, errOut, code := runVerdict(t, root, s)
	if code != 0 {
		t.Fatalf("fully live fixture must be green, got rc=%d out=%s err=%s", code, out, errOut)
	}
	for _, sc := range scopes {
		if got := strings.Count(out, "scope "+sc.label); got != 1 {
			t.Errorf("self-report line for %s printed %d times, want exactly 1 (report:\n%s)", sc.label, got, out)
		}
		// Ticket 88 inverted this assertion. It used to read "an exempt scope must
		// examine 0 files", which was only ever satisfiable by ban #6 and which
		// died with the exemption. What the fixture now pins is the stronger
		// shape: a live scope that examined nothing cannot be green, so if this
		// trips the run is lying about coverage in the other direction.
		if sc.live && sc.count(s) == 0 {
			t.Errorf("live scope %s examined 0 files in the fully live fixture: the green rc=0 below would be an empty-instrument pass", sc.label)
		}
	}
	if unc := uncoveredScopes(scopes); len(unc) != 0 {
		t.Errorf("ticket 88 armed the last exemption; this fixture is what keeps guard 3 honest, so an exempt scope here must be deliberate: %v", labels(unc))
	}
	if n := strings.Count(out, "d22scan: scope "); n != len(scopes) {
		t.Errorf("printed %d scope lines for %d declared scopes - report and ledger disagree", n, len(scopes))
	}
}

// TestVerdictRedOnEmptyBan7Scope is AC#2 for the generalized guard: ban #7's
// directory exists but holds no production Go file, which the pre-ticket-71 tool
// reported as NOTHING at all (no counter, no line) and so as a pass.
func TestVerdictRedOnEmptyBan7Scope(t *testing.T) {
	root := liveFixture(t)
	// Empty the scope without deleting the directory: the walk runs, sees no
	// production .go file, and must not be allowed to read as "checked".
	if err := os.Remove(filepath.Join(root, "internal", "tools", "ok.go")); err != nil {
		t.Fatal(err)
	}
	s := scanFixture(t, root)
	if s.examined["internal-artifact-tool"] != 0 {
		t.Fatalf("fixture broken: ban #7 examined %d files, want 0", s.examined["internal-artifact-tool"])
	}
	_, errOut, code := runVerdict(t, root, s)
	if code != 2 {
		t.Fatalf("an always-empty live scope must be fatal (rc=2), got rc=%d stderr=%q", code, errOut)
	}
	if !strings.Contains(errOut, "ban #7 internal/tools/") {
		t.Errorf("the fatal message must name the empty scope, got %q", errOut)
	}
}

// TestVerdictRedOnGoScopeWithOnlyTestFiles is the shape the OLD aggregate number
// could not see: a cmd/ holding only _test.go files still left the aggregate
// "examined N production Go files under internal/ and cmd/" large, so the run
// looked healthy while bans #1-5 were blind to every cmd/ file.
func TestVerdictRedOnGoScopeWithOnlyTestFiles(t *testing.T) {
	root := liveFixture(t)
	cmdMain := filepath.Join(root, "cmd", "wisp", "main.go")
	if err := os.Rename(cmdMain, filepath.Join(root, "cmd", "wisp", "main_test.go")); err != nil {
		t.Fatal(err)
	}
	s := scanFixture(t, root)
	cmdKey, internalKey := goScopeKey(filepath.Join(root, "cmd")), goScopeKey(filepath.Join(root, "internal"))
	if s.examined[cmdKey] != 0 {
		t.Fatalf("fixture broken: cmd/ must contribute 0 production files, got %d", s.examined[cmdKey])
	}
	agg := s.examined[internalKey] + s.examined[cmdKey]
	if agg == 0 {
		t.Fatal("fixture broken: the aggregate must be non-zero, that is the point of the test")
	}
	_, errOut, code := runVerdict(t, root, s)
	if code != 2 {
		t.Fatalf("aggregate non-zero (%d) but cmd/ empty: must be fatal, got rc=%d stderr=%q", agg, code, errOut)
	}
	if !strings.Contains(errOut, "bans #1-5 cmd/") {
		t.Errorf("must name bans #1-5 cmd/, got %q", errOut)
	}
}

// TestExemptScopeCannotOutliveItsAbsentTree is the drift guard (guard 3), pinned
// with a SYNTHETIC exempt scope rather than with ban #6.
//
// Until ticket 88 this test rode on a coincidence: ban #6 happened to be
// registered exempt, so `mkdir frontend/` in a fixture tripped the guard. Once
// ban #6 is armed - which is what the tree landing requires - that coincidence
// is gone, and the rule would have gone untested while looking untouched. So the
// scope is now built here: the assertion is about the RULE ("an exemption is a
// claim about a tree, and claims rot"), not about one ban's row in the ledger.
// Both directions are asserted, because a guard that is only ever seen firing -
// or only ever seen silent - is a tautology and not a gate.
func TestExemptScopeCannotOutliveItsAbsentTree(t *testing.T) {
	root := liveFixture(t)

	// Precondition, so the injection below cannot mask a real exemption: the
	// shipped ledger declares none, hence guard 3 has no production target and
	// this synthetic scope is the only thing keeping it proven.
	real := declaredScopes(root)
	if unc := uncoveredScopes(real); len(unc) != 0 {
		t.Fatalf("expected no exempt scope in declaredScopes() after ticket 88, got %v - re-pin this test if that changed", labels(unc))
	}
	if got := driftedAbsentScope(real); got != "" {
		t.Fatalf("the real ledger drifts on its own: %q", got)
	}

	const tree = "exempt-fixture-tree"
	dir := filepath.Join(root, tree)
	exempt := scanScope{
		label: "synthetic exempt/" + tree, dir: dir,
		kind: "text files", examinedKey: "synthetic-exempt-ban", live: false, absentOK: true,
		note: "fixture-only scope: exists so guard 3 stays testable after ticket 88 " +
			"armed ban #6 and emptied the ledger of exemptions",
	}
	scopes := append(append([]scanScope{}, real...), exempt)

	// Direction 1: tree absent -> the exemption is legal, the guard is silent and
	// the verdict stays green. Without this half, "exempt" would read as "always
	// fatal" and nobody could ever register one.
	s := scanFixture(t, root)
	if len(s.findings) != 0 {
		t.Fatalf("fixture must start finding-free so rc means what it says: %v", s.findings)
	}
	if got := driftedAbsentScope(scopes); got != "" {
		t.Fatalf("guard fired while %s is absent - it is not conditioned on the tree, it is a tautology: %q", tree, got)
	}
	out, errOut, code := runInjectedVerdict(t, root, s, exempt)
	if code != 0 {
		t.Fatalf("an exemption whose tree is really absent must be green, got rc=%d out=%s err=%s", code, out, errOut)
	}
	if !strings.Contains(out, "[NOT COVERED]") || !strings.Contains(out, "NOT COVERED: synthetic exempt/"+tree) {
		t.Errorf("the exempt scope must still be reported as not covered, got %q", out)
	}

	// Direction 2: the tree appears. Seed a real ban #6 violation in the same
	// breath so the old test's third assertion survives: the guard outranks the
	// finding (rc=2, "do not trust anything this run printed") while the finding
	// is still on stdout, because a suppression guard must not silence evidence.
	seedFile(t, root, tree+"/holder.md", "# this tree exists now\n")
	seedFile(t, root, "frontend/src/app.js", "export function decide() { return approval.decide({allow: true}); }\n")
	s = scanFixture(t, root)
	var hit bool
	for _, f := range s.findings {
		if f.Ban == "panel-approval" {
			hit = true
		}
	}
	if !hit {
		t.Fatalf("ban #6's matcher is dead - the seeded panel-approval violation was not reported: %v", s.findings)
	}
	if got := driftedAbsentScope(scopes); got != exempt.label {
		t.Fatalf("tree present while the scope is exempt must name %q, got %q", exempt.label, got)
	}
	out, errOut, code = runInjectedVerdict(t, root, s, exempt)
	if code != 2 {
		t.Fatalf("tree present while the scope is exempt must be fatal (rc=2), got rc=%d out=%s err=%s", code, out, errOut)
	}
	if !strings.Contains(errOut, exempt.label) {
		t.Errorf("drift message must name the scope, got %q", errOut)
	}
	if !strings.Contains(out, "panel-approval") {
		t.Errorf("the finding must still be printed before the guard exits, got %q", out)
	}
}

// TestBan6ScopeIsNotNarrowedByAnExtensionFilter pins ticket 88 AC#1's real
// question. ban #6 predates the panel: if walkText filtered by suffix and the
// panel turned out to be .tsx, arming the scope would have produced a gate that
// is live, green, and blind. The fixture therefore plants a violation in every
// extension class the panel actually uses plus the two skip rules, and demands
// the counts and the hits come out as the wording promises.
func TestBan6ScopeIsNotNarrowedByAnExtensionFilter(t *testing.T) {
	root := liveFixture(t)
	if err := os.RemoveAll(filepath.Join(root, "frontend")); err != nil {
		t.Fatal(err)
	}
	seedFile(t, root, "tools/d22scan/allowlist.txt", "# empty\n")
	cases := []string{
		"frontend/src/Card.tsx",    // ticket 77's components
		"frontend/src/lib/util.ts", // ticket 77's helpers
		"frontend/scripts/dev.mjs", // ticket 77's tooling
		"frontend/src/app.css",     // styles
		"frontend/README.md",       // docs
		"frontend/package.json",    // config
		"frontend/Procfile",        // no extension at all
	}
	for _, rel := range cases {
		seedFile(t, root, rel, "export const x = approval.decide();\n")
	}
	// Out of scope by rule, not by suffix: these two directories are skipped for
	// every ban (vendored code is not our source).
	seedFile(t, root, "frontend/node_modules/pk/index.tsx", "approval.decide();\n")
	seedFile(t, root, "frontend/testdata/golden.tsx", "approval.decide();\n")

	s := scanFixture(t, root)
	if got := s.examined["panel-approval"]; got != len(cases) {
		t.Errorf("ban #6 examined %d files, want %d - walkText's filter and this list disagree, and one of them is narrowing the ban", got, len(cases))
	}
	byPath := map[string]int{}
	for _, f := range s.findings {
		if f.Ban == "panel-approval" {
			byPath[f.Path]++
		}
	}
	for _, rel := range cases {
		if byPath[rel] == 0 {
			t.Errorf("ban #6 did not fire in %s: the scope examines it but the matcher is blind to that file class (all findings %v)", rel, s.findings)
		}
		delete(byPath, rel)
	}
	for p := range byPath {
		t.Errorf("unexpected panel-approval finding at %s (vendored/testdata leak?)", p)
	}
	t.Logf("ban #6 examined %d and fired in all %d extension classes", s.examined["panel-approval"], len(cases))
}

// TestVerdictRedOnUndeclaredCounter pins the structural guard: a walk that bumps a
// counter no scope reports is work the self-report does not account for. This is
// exactly the hole opened by adding an s.walkGo(...) call and forgetting the
// ledger entry, which would otherwise print a confident "clean".
func TestVerdictRedOnUndeclaredCounter(t *testing.T) {
	root := liveFixture(t)
	s := scanFixture(t, root)
	s.examined["some-future-ban"] = 7
	out, errOut, code := runVerdict(t, root, s)
	if code == 0 {
		t.Fatalf("an undeclared counter must not yield a verdict, got rc=%d out=%s", code, out)
	}
	if code != 2 {
		t.Fatalf("want rc=2 for undeclared work, got rc=%d err=%q", code, errOut)
	}
	if !strings.Contains(errOut, "some-future-ban") {
		t.Errorf("guard must name the undeclared counter, got %q", errOut)
	}
}

// TestRealRepoBan8CoversFrontendTreeAtBan6sCount is ticket 96 AC#1's number and
// AC#3's falsifier in one place. The two gates now walk the SAME tree with the
// SAME rule - every file, node_modules/ and testdata/ skipped - so their counts
// are not merely "close", they are the same integer, and this test demands
// equality rather than explaining a difference away.
//
// AC#3's half: every expectation below is a LITERAL "ban #8 frontend/". The
// other real-repo ledger tests build their expectations by iterating over
// emojiScopes(), which is exactly why they all stayed green for the whole time
// the panel was uncovered - delete frontend/ from that list and they shrink
// theirselves to three rows and still agree. This test cannot be satisfied that
// way: if the entry goes missing, emojiSeen["frontend/"] is 0, ban #6's is 37+,
// the equality fails, the >0 guard fires and the printed line disappears.
func TestRealRepoBan8CoversFrontendTreeAtBan6sCount(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Skipf("not inside the wisp repo: %v", err)
	}
	s := scanFixture(t, root)

	front := s.emojiSeen["frontend/"]
	if front == 0 {
		t.Fatal("ban #8 examined 0 files under frontend/: the panel's tree is uncovered again (ticket 96), and a run that examines nothing does not get to print \"clean\"")
	}
	if ban6 := s.examined["panel-approval"]; front != ban6 {
		t.Errorf("ban #8 examined %d frontend/ files but ban #6 examined %d - the two gates walk the same tree with the same rule (all files; node_modules/ and testdata/ skipped), so any difference is an exclusion or a filter that one walk grew and the other did not", front, ban6)
	}

	var labels []string
	for _, sc := range declaredScopes(root) {
		labels = append(labels, sc.label)
		if sc.label != "ban #8 frontend/" {
			continue
		}
		if !sc.live {
			t.Errorf("ban #8 frontend/ is registered non-live, so an empty walk there would print a verdict instead of failing")
		}
		if sc.count(s) != front {
			t.Errorf("ledger reports %d for ban #8 frontend/ while emojiSeen says %d", sc.count(s), front)
		}
	}
	if !strings.Contains(strings.Join(labels, ","), "ban #8 frontend/") {
		t.Errorf("declaredScopes() carries no ban #8 frontend/ entry: %v - the ledger must name the tree the ban claims to cover (ticket 96 AC#1)", labels)
	}

	out, errOut, code := runVerdict(t, root, s)
	// HEAD may be red for reasons this ticket does not own, so only the scope
	// lines are claimed here; TestRealRepoLedgerIsHonest owns the verdict.
	//
	// describeScopes pads the label to %-24s, so the line is compared after
	// collapsing runs of spaces - otherwise this test would be pinning the
	// column width rather than the coverage.
	line := ""
	for _, l := range strings.Split(out, "\n") {
		if strings.Contains(l, "scope ban #8 frontend/") {
			line = strings.Join(strings.Fields(l), " ")
		}
	}
	want := fmt.Sprintf("d22scan: scope ban #8 frontend/ examined %d text files", front)
	if line != want {
		t.Errorf("the self-report line for ban #8 frontend/ is %q, want %q (rc=%d err=%s)", line, want, code, errOut)
	}
	if sum := scopeSummary(declaredScopes(root), s); !strings.Contains(sum, fmt.Sprintf("ban #8 frontend/=%d", front)) {
		t.Errorf("the clean line's scope summary omits ban #8 frontend/=%d, got %q", front, sum)
	}
}

// TestBan8FrontendScopeIsNotNarrowedByAnExtensionFilter is ticket 96 AC#2's
// in-repo twin and AC#3's decision spelled out: reusing the existing
// goOnly:false branch would have routed frontend/ through isTextFile(), which
// drops .mjs, extension-less files and dotfiles - measured 5 of the 37 tracked
// files in frontend/ today, including all three scripts/*.mjs, i.e. the
// vendoring tooling R18 exists to police. Every class below therefore carries
// the glyph as plain text, not in a comment, because Q-46(c) (ticket 141) made
// walkEmoji's comment rule exempt prose; a comment seed would prove coverage of
// nothing.
//
// The `.md` row answers the ticket's third question explicitly: yes, non-.go
// text is scanned, and the reason is that the panel's most likely emoji are UI
// copy and notes, and D23 governs design language, which lives in prose. The
// last two rows are the exclusion rule (AC#2's "排除了什么"): node_modules/ and
// testdata/ are out for every ban in this tool, so an emoji there is reported
// neither by ban #6 nor by ban #8.
func TestBan8FrontendScopeIsNotNarrowedByAnExtensionFilter(t *testing.T) {
	root := liveFixture(t)
	if err := os.RemoveAll(filepath.Join(root, "frontend")); err != nil {
		t.Fatal(err)
	}
	seedFile(t, root, "tools/d22scan/allowlist.txt", "# empty\n")
	// U+2713 CHECK MARK, inside emojiRe's \x{2600}-\x{27BF} band.
	cases := []struct {
		rel, why string
	}{
		{"frontend/src/Card.tsx", "the panel's own components"},
		{"frontend/src/lib/util.ts", "helpers"},
		{"frontend/scripts/vendor.mjs", "the vendoring tooling isTextFile() has no .mjs for"},
		{"frontend/src/app.css", "styles"},
		{"frontend/VENDORED.md", "docs: this is the .md answer - .md IS scanned"},
		{"frontend/package.json", "config"},
		{"frontend/Procfile", "no extension at all, also scanned"},
	}
	for _, c := range cases {
		seedFile(t, root, c.rel, "ready \u2713\n")
	}
	seedFile(t, root, "frontend/node_modules/pk/index.tsx", "ready \u2713\n")
	seedFile(t, root, "frontend/testdata/golden.tsx", "ready \u2713\n")

	s := scanFixture(t, root)
	if got := s.emojiSeen["frontend/"]; got != len(cases) {
		t.Errorf("ban #8 examined %d frontend/ files, want %d - walkEmoji's everyFile branch and this list disagree, and one of them is narrowing the ban", got, len(cases))
	}
	byPath := map[string]int{}
	for _, f := range s.findings {
		if f.Ban == "emoji" && strings.HasPrefix(filepath.ToSlash(f.Path), "frontend/") {
			byPath[filepath.ToSlash(f.Path)]++
		}
	}
	for _, c := range cases {
		want := filepath.ToSlash(c.rel)
		if byPath[want] == 0 {
			t.Errorf("ban #8 did not fire in %s (%s): the scope counted it but the matcher is blind to that file class (all findings %v)", c.rel, c.why, s.findings)
		}
		delete(byPath, want)
	}
	for p := range byPath {
		t.Errorf("unexpected ban #8 hit at %s: node_modules/ or testdata/ leaked into the emoji scope", p)
	}
	t.Logf("ban #8 examined %d and fired in all %d frontend/ file classes, incl. .mjs, .md and extension-less", s.emojiSeen["frontend/"], len(cases))
}

// TestLedgerCountsMatchAnIndependentWalk is AC#4's falsifier: the printed numbers
// are compared against a SECOND, independent walk computed in the test. Without
// it a ledger that counts the wrong thing (every .go including _test.go, or a
// directory skipped by mistake) stays internally consistent and still prints a
// confident number - "examined N" only means something if N is the truth.
//
// A207 ADDS A SECOND COPY OF THE IGNORE RULE HERE, ON PURPOSE. The scanner reads
// .gitignore itself (gitignore.go); the walk below carries a hand-written literal
// of what that file says today. That is the only way this test keeps its teeth:
// if both sides called the same function, an ignore rule that swallowed a live
// scope's files would be counted by the instrument and confirmed by its own
// oracle. A disagreement between the two copies is the signal A207 ② demands -
// the ignore policy now reaches inside a scanned scope, which moves a baseline
// number, and that must be a deliberate ledger-visible decision rather than
// something the instrument absorbed silently. So: when someone adds an ignore
// rule that binds design/, frontend/, internal/ or cmd/, THIS LIST MUST MOVE IN
// THE SAME BATCH and the scope baselines in the ledger get re-stated with it.
//
// A214/F1 EXTENDS THAT DUTY TO THE INDEX, IN THE SAME COMMIT AS THE CODE. The copy
// above now takes the tracked set and a "were any rules applied at all" flag, the
// same two facts gitignore.go's skip() uses. That is not decoration: an oracle that
// keeps only the pattern half agrees with the instrument on the one shape the
// acceptance built by hand, so the equality guard endorses the blindfold instead of
// catching it (d22scan-gitignore-accept-r1.md §1.4, §6.8 - four guards, all green,
// on a tree with a force-added `<=` in it).
func TestLedgerCountsMatchAnIndependentWalk(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Skipf("not inside the wisp repo: %v", err)
	}
	s := scanFixture(t, root)

	// A214/F1: the oracle has to know what the index holds, or it agrees with the
	// instrument on the one shape that matters and both go green over a delivered
	// byte. Same three questions gitignore.go asks, asked here by hand.
	tracked, indexOK, indexWhy := gitTrackedIndex(root)
	if !indexOK {
		t.Logf("independent walk: no index to consult (%s) - applying no ignore rule, exactly as the instrument must", indexWhy)
	}

	count := func(dir string, accept func(string) bool) int {
		n := 0
		err := filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if d.Name() == "testdata" || d.Name() == ".git" || d.Name() == "node_modules" {
					return filepath.SkipDir
				}
				return nil
			}
			if ignoredLikeGit(relToRepo(root, p), tracked, indexOK) {
				return nil
			}
			if accept(p) {
				n++
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		return n
	}
	goFile := func(p string) bool { return strings.HasSuffix(p, ".go") }
	prodGo := func(p string) bool { return goFile(p) && !strings.HasSuffix(p, "_test.go") }

	cases := []struct {
		label  string
		want   int
		report int
	}{
		{"bans #1-5 internal/", count(filepath.Join(root, "internal"), prodGo), s.examined[goScopeKey(filepath.Join(root, "internal"))]},
		{"bans #1-5 cmd/", count(filepath.Join(root, "cmd"), prodGo), s.examined[goScopeKey(filepath.Join(root, "cmd"))]},
		{"ban #7 internal/tools/", count(filepath.Join(root, "internal", "tools"), prodGo), s.examined["internal-artifact-tool"]},
		{"ban #8 internal/", count(filepath.Join(root, "internal"), goFile), s.emojiSeen["internal/"]},
		{"ban #8 cmd/", count(filepath.Join(root, "cmd"), goFile), s.emojiSeen["cmd/"]},
		// Ticket 88: ban #6 used to be missing here because it always walked 0
		// files and needed no falsifier. Now that its number decides CI it gets
		// one, and `anything` is the point - the independent walk accepts every
		// file the way walkText does, so a future suffix filter in either place
		// shows up as a mismatch instead of a smaller number nobody noticed.
		{"ban #6 frontend/", count(filepath.Join(root, "frontend"), func(string) bool { return true }), s.examined["panel-approval"]},
		// Ticket 96: ban #8 now walks the SAME tree with the SAME rule (every
		// file, node_modules/ and testdata/ skipped), so this scope's falsifier
		// is the identical `anything` predicate as ban #6's two lines above. A
		// suffix filter reappearing in either walk shows up here as a mismatch
		// against this independent count, and a frontend/ entry deleted from
		// emojiScopes() shows up as report=0.
		{"ban #8 frontend/", count(filepath.Join(root, "frontend"), func(string) bool { return true }), s.emojiSeen["frontend/"]},
	}
	for _, c := range cases {
		if c.report != c.want {
			t.Errorf("scope %s reported examining %d files, independent walk says %d - the number in the self-report is wrong, not just small",
				c.label, c.report, c.want)
		}
		if c.report == 0 {
			t.Errorf("scope %s examined 0 files in the real repo", c.label)
		}
		t.Logf("verified %s: %d files", c.label, c.report)
	}
}

// TestWalksSkipGitIgnoredPaths is A207's mutation pair inside one fixture: the
// SAME glyph and the SAME forbidden call, once on a gitignored path and once on a
// path the repository would track. The ignored copy must move no number and raise
// no finding; the non-ignored copy must do both. The second half is the load-bearing
// one - an "ignore" rule that also swallowed tracked files would turn this gate
// into a blindfold wearing a denominator fix, and this repository has already been
// caught twice with a scope that was "armed but blind" (ticket 88, A207 ①'s
// 37/40/43).
//
// The .gitignore contents mirror the real ones (root .gitignore `frontend/dist/*`
// + `!frontend/dist/.gitkeep`, frontend/.gitignore `dist/*` + `!dist/.gitkeep`)
// so this pins the rule SHAPE the scanner reads, not a convenient invention. The
// negation matters for the counts: `frontend/dist/.gitkeep` is one of the 40 files
// CI reports, so excluding it would move CI's number as well as the local one.
func TestWalksSkipGitIgnoredPaths(t *testing.T) {
	root := liveFixture(t)
	seedFile(t, root, ".gitignore", "# fixture copy of the repository's own rules\n\nfrontend/dist/*\n!frontend/dist/.gitkeep\nbuild/\n*.generated.go\n")
	seedFile(t, root, "frontend/.gitignore", "dist/*\n!dist/.gitkeep\n")

	// U+2264 (inside emojiRe's math band, Q-46(c)/ticket 141) plus, where the ban
	// owns one, a forbidden call shape. Ignored side first:
	seedFile(t, root, "frontend/dist/index.html", "<p>ready \u2264</p>\n<script>approval.decide({allow: true})</script>\n")
	seedFile(t, root, "frontend/dist/probe.txt", "ready \u2264\n")
	seedFile(t, root, "internal/build/probe.go", "package build\n\nconst doc = \"ready \u2264\"\n\nfunc worker() {}\n\nfunc leak() { go worker() }\n")
	seedFile(t, root, "internal/pkg/x.generated.go", "package pkg\n\nconst doc = \"ready \u2264\"\n")
	// Non-ignored side, including the re-included anchor and a plain tracked file:
	seedFile(t, root, "frontend/dist/.gitkeep", "ready \u2264\n")
	seedFile(t, root, "frontend/src/live.tsx", "export const doc = \"ready \u2264\";\n")
	seedFile(t, root, "frontend/src/decide.tsx", "export function d() { return approval.decide({allow: true}); }\n")
	seedFile(t, root, "internal/pkg/keep.go", "package pkg\n\nconst doc = \"ready \u2264\"\n\nfunc worker() {}\n\nfunc leak() { go worker() }\n")

	// A214/F1 turned this fixture into a repository. "Ignored" is not a fact a
	// .gitignore file alone carries - git's criterion is the rule AND the index -
	// so the matcher now asks git, and a tree whose index cannot be asked gets no
	// rules at all (see TestIgnoreRulesAreOffWhenTheIndexCannotBeAsked for that
	// half). `git add -A` is what makes the claim below mean something: everything
	// the rules do NOT cover becomes tracked, everything they DO cover stays
	// untracked, which is the state a real worktree is in after a build.
	gitIndexFixture(t, root)

	s := scanFixture(t, root)

	// 1. the denominators. frontend/ = panel.tsx (liveFixture) + .gitkeep +
	// live.tsx + decide.tsx + frontend/.gitignore itself, which is a tracked file
	// in the tree and both frontend walks take EVERY file (no suffix allowlist -
	// ticket 96 AC#2), so it is one of the counted five; the two dist/ seeds are
	// out. internal/'s Go count is ok/ok.go + tools/ok.go + pkg/keep.go;
	// build/probe.go and pkg/x.generated.go are out. Ban #6 and ban #8 read the
	// same tree with the same rule, so their numbers must be ONE integer (the
	// invariant the real-repo test above also demands).
	if got, want := s.emojiSeen["frontend/"], 5; got != want {
		t.Errorf("ban #8 examined %d frontend/ files, want %d", got, want)
	}
	if got, want := s.examined["panel-approval"], 5; got != want {
		t.Errorf("ban #6 examined %d frontend/ files, want %d", got, want)
	}
	if ban6, ban8 := s.examined["panel-approval"], s.emojiSeen["frontend/"]; ban6 != ban8 {
		t.Errorf("the two frontend/ walks disagree (%d vs %d): one of them grew an exclusion the other did not", ban6, ban8)
	}
	if got, want := s.examined[goScopeKey(filepath.Join(root, "internal"))], 3; got != want {
		t.Errorf("bans #1-5 examined %d internal/ production Go files, want %d", got, want)
	}
	if got, want := s.emojiSeen["internal/"], 3; got != want {
		t.Errorf("ban #8 examined %d internal/ Go files, want %d", got, want)
	}

	// 2. the findings, by path.
	byBan := map[string]map[string]int{}
	for _, f := range s.findings {
		if byBan[f.Ban] == nil {
			byBan[f.Ban] = map[string]int{}
		}
		byBan[f.Ban][f.Path]++
	}
	for ban, want := range map[string][]string{
		"emoji":          {"frontend/dist/.gitkeep", "frontend/src/live.tsx", "internal/pkg/keep.go"},
		"panel-approval": {"frontend/src/decide.tsx"},
		"bare-goroutine": {"internal/pkg/keep.go"},
	} {
		got := byBan[ban]
		if len(got) != len(want) {
			t.Errorf("ban %q fired at %v, want exactly %v", ban, got, want)
		}
		for _, p := range want {
			if got[p] == 0 {
				t.Errorf("ban %q did NOT fire at the tracked path %s - an ignore rule that hides non-ignored files is a blindfold, not a fix (all findings %v)", ban, p, s.findings)
			}
		}
	}
	for _, f := range s.findings {
		for _, leak := range []string{"frontend/dist/index.html", "frontend/dist/probe.txt", "internal/build/probe.go", "internal/pkg/x.generated.go"} {
			if f.Path == leak {
				t.Errorf("ban %q fired at the gitignored path %s - the exclusion did not take (all findings %v)", f.Ban, leak, s.findings)
			}
		}
	}

	// 3. provenance: a number that moved must say which rule file moved it, or
	// A207's "three readings, no explanation" is back. The fixture excludes 3 files
	// (dist/index.html, dist/probe.txt, pkg/x.generated.go) plus one pruned
	// directory (internal/build/), decided by two rule files.
	note := s.ign.note()
	if note == "" {
		t.Fatal("the matcher excluded 4 paths but the self-report says nothing about it")
	}
	for _, want := range []string{"3 file(s)", "internal/build/", "frontend/.gitignore"} {
		if !strings.Contains(note, want) {
			t.Errorf("note %q must name %q", note, want)
		}
	}
	t.Logf("note: %s", note)
}

// gitCommand runs one git command inside dir and fails the test OUT LOUD if git
// is not there. No t.Skip on this path, on purpose: runtests.sh counts SKIP as not
// a pass (ticket 71 AC#3), and a branch that silently stops being exercised is
// the same shape this file exists to catch. The scanner itself does NOT require a
// git binary - it degrades to scanning everything and says so - but proving the
// tracked-path guarantee needs the index to be real, and a real index needs git.
func gitCommand(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	var stdout, stderr strings.Builder
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("git %s in %s: %v (%s) - the tracked-path guarantees of A214 cannot be built without a git binary",
			strings.Join(args, " "), dir, err, firstLineOf(stderr.String()))
	}
	return stdout.String()
}

func firstLineOf(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	for _, line := range strings.Split(s, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			return line
		}
	}
	return "(no stderr)"
}

// gitIndexFixture turns a fixture tree into a repository and stages everything
// git would take on its own - i.e. everything the ignore rules do NOT cover. It
// does not commit, because `git ls-files` reads the index, and the index is the
// only thing the matcher asks.
func gitIndexFixture(t *testing.T, root string) {
	t.Helper()
	gitCommand(t, root, "init", "-q", "-b", "main")
	gitCommand(t, root, "add", "-A")
}

// gitForceAdd is `git add -f`: the one documented way a path ends up tracked while
// an ignore rule still covers it, which is the shape the r1 acceptance built by
// hand outside the repository (d22scan-gitignore-accept-r1.md §1.3) and the shape
// this batch was退回 for not defending against.
func gitForceAdd(t *testing.T, root string, rels ...string) {
	t.Helper()
	gitCommand(t, root, append([]string{"add", "-f", "--"}, rels...)...)
}

// TestTrackedPathsAreNeverSkippedByTheIgnoreFilter is A214/F1's mutation pair and
// the assertion that makes the two sentences F2 called unfalsifiable checkable:
// main.go's "it does not suppress a finding on a TRACKED path" and gitignore.go's
// "for a path git reports as tracked, skip() returns false".
//
// ONE PAIR, ONE DIFFERENCE. Both halves hold the same bytes on disk under the same
// ignore rules; the only thing that changes is whether the index holds them. The
// untracked half must stay out of the denominators (that is A207, and a fix that
// loses it is a revert wearing a coat); the tracked half must be counted, must go
// red, and must be named by the self-report. Before this batch the tracked half
// read `frontend/=40` / `internal/`=203 and rc=0, hiding a `<=` in a `.tsx` under
// `dist/` and a bare `go worker()` under the root `build/` rule - a glyph ban and a
// SECURITY ban, both invisible, in the tool's own positive-control layer too.
func TestTrackedPathsAreNeverSkippedByTheIgnoreFilter(t *testing.T) {
	seed := func(t *testing.T) string {
		t.Helper()
		root := liveFixture(t)
		seedFile(t, root, ".gitignore", "frontend/dist/*\n!frontend/dist/.gitkeep\nbuild/\n")
		seedFile(t, root, "frontend/.gitignore", "dist/*\n!dist/.gitkeep\n")
		seedFile(t, root, "frontend/dist/tracked-or-not.tsx", "export const doc = \"ready \u2264\";\n")
		seedFile(t, root, "internal/build/leak.go", "package build\n\nconst doc = \"ready \u2264\"\n\nfunc worker() {}\n\nfunc leak() { go worker() }\n")
		gitIndexFixture(t, root)
		return root
	}
	untracked := seed(t)
	tracked := seed(t)
	gitForceAdd(t, tracked, "frontend/dist/tracked-or-not.tsx", "internal/build/leak.go")

	// git's own instruments first, so the test cannot pass on a split the scanner
	// invented: the same two paths, untracked in one tree and tracked-while-matched
	// in the other. `git check-ignore` is the same question git itself answers, and
	// the index-aware form must say "not ignored" for a force-added path while the
	// --no-index form (what the r1 matcher implemented) says "ignored".
	if got := strings.TrimSpace(gitCommand(t, untracked, "ls-files", "-i", "-c", "--exclude-standard")); got != "" {
		t.Fatalf("the untracked half must have nothing tracked-and-ignored, got %q", got)
	}
	both := gitCommand(t, tracked, "ls-files", "-i", "-c", "--exclude-standard")
	for _, want := range []string{"frontend/dist/tracked-or-not.tsx", "internal/build/leak.go"} {
		if !strings.Contains(both, want) {
			t.Fatalf("git ls-files -i -c --exclude-standard must name %s, got %q", want, both)
		}
		if out := gitCommand(t, tracked, "check-ignore", "--no-index", "-v", want); !strings.Contains(out, want) {
			t.Fatalf("git check-ignore --no-index must still match %s (the pattern half), got %q", want, out)
		}
	}

	sA := scanFixture(t, untracked)
	sB := scanFixture(t, tracked)
	goKey := func(root string) string { return goScopeKey(filepath.Join(root, "internal")) }

	// 1. the denominators, as a DELTA between the halves so a future change that
	//    moves both together cannot pass by shrinking the pair.
	if a, b := sA.examined["panel-approval"], sB.examined["panel-approval"]; b != a+1 {
		t.Errorf("ban #6 counted %d frontend/ files tracked and %d untracked - the force-added .tsx must add exactly 1", b, a)
	}
	if a, b := sA.emojiSeen["frontend/"], sB.emojiSeen["frontend/"]; b != a+1 {
		t.Errorf("ban #8 counted %d frontend/ files tracked and %d untracked, want +1", b, a)
	}
	if a, b := sA.examined[goKey(untracked)], sB.examined[goKey(tracked)]; b != a+1 {
		t.Errorf("bans #1-5 counted %d internal/ production Go files tracked and %d untracked, want +1 (leak.go under the root `build/` rule)", b, a)
	}
	if a, b := sA.emojiSeen["internal/"], sB.emojiSeen["internal/"]; b != a+1 {
		t.Errorf("ban #8 counted %d internal/ Go files tracked and %d untracked, want +1", b, a)
	}
	if b, a := sB.emojiSeen["frontend/"], sB.examined["panel-approval"]; b != a {
		t.Errorf("the two frontend/ walks disagree on a tracked-ignored tree (%d vs %d)", b, a)
	}

	// 2. the findings. Tracked means visible: both bans must fire at both paths.
	byPath := map[string]map[string]bool{}
	for _, f := range sB.findings {
		if byPath[f.Path] == nil {
			byPath[f.Path] = map[string]bool{}
		}
		byPath[f.Path][f.Ban] = true
	}
	for path, bans := range map[string][]string{
		"frontend/dist/tracked-or-not.tsx": {"emoji"},
		"internal/build/leak.go":           {"emoji", "bare-goroutine"},
	} {
		for _, ban := range bans {
			if !byPath[path][ban] {
				t.Errorf("ban %q did NOT fire at the tracked path %s - a rule that hides a delivered byte is a blindfold, not a fix (all findings %v)", ban, path, sB.findings)
			}
		}
	}
	for _, f := range sA.findings {
		if strings.Contains(f.Path, "tracked-or-not.tsx") || strings.Contains(f.Path, "build/leak.go") {
			t.Errorf("ban %q fired at %s in the UNTRACKED half: A207's exclusion was lost, ignored build output is back in the denominator (%v)", f.Ban, f.Path, sA.findings)
		}
	}

	// 3. provenance, both directions. The untracked half names what it skipped;
	//    the tracked half must name what it REFUSED to skip, because "the count went
	//    up by one and nothing said why" is the A207 ① complaint all over again.
	if n := sA.ign.note(); !strings.Contains(n, "skipped as git-ignored") {
		t.Errorf("the untracked half must report its exclusions, got %q", n)
	}
	nb := sB.ign.note()
	for _, want := range []string{"TRACKED and matching an ignore rule", "frontend/dist/tracked-or-not.tsx", "internal/build/leak.go"} {
		if !strings.Contains(nb, want) {
			t.Errorf("note %q must name %q", nb, want)
		}
	}

	// 4. the exit codes: red where a delivered byte breaks a ban, green where the
	//    same byte is untracked build output.
	if _, errOut, code := runVerdict(t, tracked, sB); code != 1 {
		t.Errorf("verdict on the tracked half must be 1 (findings), got %d err=%s", code, errOut)
	}
	if _, errOut, code := runVerdict(t, untracked, sA); code != 0 {
		t.Errorf("verdict on the untracked half must stay 0, got %d err=%s", code, errOut)
	}
}

// TestFullTrackedListCoversWhatTheNarrowListCannot is the third half of the
// ingredient pair (acceptance d22scan-gitignore-fix-accept-r1.md §4 and its N3,
// 2026-09-25). runGitIndex asks git TWO listing questions and the union feeds
// holds(): `ls-files -z` (every tracked path) and `ls-files -z -i -c
// --exclude-standard` (the tracked paths that also match a rule). Until this
// test existed, nothing required the first one - removing it
// (the `all` argument of the parseIndexPaths pair inside runGitIndex) left
// the whole delivered roster green, which acceptance mutant M1 measured and this
// file does not accept as "unneeded". Deliberately no bare line number here:
// this same batch moved that call site twice within one commit, so a number
// would rot faster than the fact it points at (acceptance r1 §4.1, 退回枚 1).
//
// WHY THE FULL LIST IS LOAD-BEARING AND THE NARROW ONE STRUCTURALLY CANNOT
// REPLACE IT: `-i -c` can only report a path git itself calls matched, so the
// class it is blind to is "this matcher skips a tracked path that git does not
// match". That class is not hypothetical - F3 ① and F3 ② (a rule line's leading
// whitespace trimmed into a different pattern, a trailing `**` pruning its own
// parent directory) are two instances of it fixed in the same batch that added
// the guard, so a third instance is guarded against by data, not by faith.
//
// THE SEED: `frontend/.gitignore` carries `  weird/*` - TWO leading spaces, which
// git keeps as pattern data (measured: `git check-ignore -v
// frontend/weird/inside.tsx` exits 1 in this shape, and the same rule binds a
// directory literally named `  weird`). A plain `git add -A` therefore tracks the
// file while `-i -c` names nothing for it. Both halves are asserted: the byte must
// be scanned and named, and holds() must know about it from the full list. The
// force-added `frontend/dist/` sibling inside the SAME tree is the positive
// control for the zero reading - the exact-`-i -c`-content assertion below is what
// stops "no hits" here from meaning "the instrument was not looking".
func TestFullTrackedListCoversWhatTheNarrowListCannot(t *testing.T) {
	seed := func(t *testing.T, withWeird bool) string {
		t.Helper()
		root := liveFixture(t)
		seedFile(t, root, ".gitignore", "frontend/dist/*\n!frontend/dist/.gitkeep\nbuild/\n")
		fe := "dist/*\n!dist/.gitkeep\n"
		if withWeird {
			fe += "  weird/*\n" // two leading spaces ARE the pattern's first characters
		}
		seedFile(t, root, "frontend/.gitignore", fe)
		seedFile(t, root, "frontend/dist/tracked-or-not.tsx", "export const doc = \"ready \u2264\";\n")
		if withWeird {
			seedFile(t, root, "frontend/weird/inside.tsx", "export const doc = \"ready \u2264\";\n")
		}
		gitIndexFixture(t, root) // plain `git add -A`: this path is tracked WITHOUT a force-add
		gitForceAdd(t, root, "frontend/dist/tracked-or-not.tsx")
		return root
	}
	base := seed(t, false)
	root := seed(t, true)

	// git's own instruments, both readings in one assertion so the empty half is
	// only ever read against a non-empty control on the same tree and command.
	if got := strings.TrimSpace(gitCommand(t, root, "ls-files", "-i", "-c", "--exclude-standard")); got != "frontend/dist/tracked-or-not.tsx" {
		t.Fatalf("the premise of this case is `-i -c` naming the force-added path and NOTHING else, got %q - "+
			"if the weird/ path appears here the seed no longer pins the full-list leg, and if the dist/ path vanished "+
			"the instrument is not reading this tree at all", got)
	}
	if out := gitCommand(t, root, "ls-files"); !strings.Contains(out, "frontend/weird/inside.tsx") {
		t.Fatalf("a plain `git add -A` must track frontend/weird/inside.tsx (git matches no rule on it), ls-files=%q", out)
	}

	sB := scanFixture(t, base)
	s := scanFixture(t, root)

	// 1. counted: exactly one more file in both frontend/ walks than the same tree
	//    without that seed and rule.
	if a, b := sB.examined["panel-approval"], s.examined["panel-approval"]; b != a+1 {
		t.Errorf("ban #6 counted %d frontend/ files without the seed and %d with it, want +1 - the whitespace rule swallowed a tracked delivered byte", a, b)
	}
	if a, b := sB.emojiSeen["frontend/"], s.emojiSeen["frontend/"]; b != a+1 {
		t.Errorf("ban #8 counted %d frontend/ files without the seed and %d with it, want +1", a, b)
	}
	// 2. read: a tracked byte that breaks a ban has to be named, not just counted.
	saw := false
	for _, f := range s.findings {
		if f.Path == "frontend/weird/inside.tsx" && f.Ban == "emoji" {
			saw = true
		}
	}
	if !saw {
		t.Errorf("ban #8 did NOT fire at the tracked path frontend/weird/inside.tsx - that is the invisible-delivered-byte shape of A214/F1, and only the full `git ls-files` list can see it (findings=%v)", s.findings)
	}
	// 3. the ingredient itself: holds() must answer this path from the FULL list.
	//    This is the assertion that mutant M1 (loading `holds()`' guard from
	//    `-i -c` instead of `ls-files`) has to break on, because nothing else in
	//    either list can know about a path git does not call matched.
	ix := s.ign.index()
	if !ix.ok {
		t.Fatalf("this fixture must be a readable repository for the guard leg to be testable: %s", ix.why)
	}
	if ix.ignored["frontend/weird/inside.tsx"] {
		t.Fatal("the narrow list now reports this path too, so this case no longer pins the full-list leg - re-seed with a rule git declines to match")
	}
	if !ix.holds("frontend/weird/inside.tsx", false) {
		t.Error("holds() does not know a path the index plainly holds: the guard leg is being loaded from `git ls-files -i -c --exclude-standard` instead of `git ls-files -z`, which cannot see a tracked path git does not match a rule against - precisely the matcher-over-skips class F3 ① and F3 ② are two instances of")
	}
	// 4. and the tree goes red, so a regression here is a loud failure and not a
	//    silently unchanged denominator.
	if _, errOut, code := runVerdict(t, root, s); code != 1 {
		t.Errorf("verdict must be 1 (the tracked byte breaks ban #8), got %d err=%s", code, errOut)
	}
}

// gitCommitAll commits whatever the fixture's index holds. The identity and the
// gpgsign flag are passed on the command line instead of trusted from the machine:
// on a runner with no `user.email` configured a bare `git commit` dies with
// "Please tell me who you are", and a fixture whose commit does not land is a
// fixture that quietly stops testing the thing it was written for.
func gitCommitAll(t *testing.T, dir, msg string) {
	t.Helper()
	gitCommand(t, dir, "-c", "user.name=d22scan fixture", "-c", "user.email=fixture@invalid",
		"-c", "commit.gpgsign=false", "commit", "-q", "-m", msg)
}

// TestIndexChangingBetweenTheTwoGitReadsAppliesNoRules is ticket 142: runGitIndex
// asks git TWO listing questions in two independent spawns, so a commit landing
// between them leaves the second answer naming a tracked path the first answer
// never saw - and since `dirs` is derived from the first answer only, holds(dir)
// is FALSE for a directory the index really holds. A rule then prunes a live tree
// and the scan says nothing, which is the exact direction this repository has been
// killing since A207 (measured by a non-implementer at anchor 310816f, quoted in
// the ticket: holds(file)=true, holds(dir)=false, one-directional under-scan).
//
// WHY THIS SHAPE AND NOT A TIMED RACE. The window is ~one subprocess spawn wide;
// a test that commits "in the hope" of hitting it is a test that passes when the
// machine is idle and flakes when it is not, which is worse than no test here
// because the guard it claims to pin is the loud kind. So the two reads are taken
// deliberately, across a commit that is certainly in between, and handed to the
// SAME function runGitIndex uses to assemble them (buildGitIndexState). What the
// test does that production does not is *not* run the two spawns back to back;
// everything after that is production code, including the walks in block 5.
//
// CONTROLS, because three of the readings below are zeroes: block 1 demands that
// the narrow read really names the path and the full read really does not (if the
// seed stops landing in the window the test fails OUT rather than passing on an
// empty premise); block 3 demands that the same narrow stream against a matching
// full read assembles fine (so "refused" is not "always refused"); block 4 demands
// that the ignore rule being declined is LIVE on this directory (so "skipped
// nothing" is not "there was nothing to skip"); and block 5 keeps an untracked
// ignored file in the tree so the quiescent scan still skips something - this guard
// may not pass by having ignore filtering switched off globally.
func TestIndexChangingBetweenTheTwoGitReadsAppliesNoRules(t *testing.T) {
	const stray = "frontend/dist/assets/shipped.tsx"
	root := liveFixture(t)
	seedFile(t, root, ".gitignore", "frontend/dist/*\n!frontend/dist/.gitkeep\nbuild/\n")
	seedFile(t, root, "frontend/.gitignore", "dist/*\n!dist/.gitkeep\n")
	// Tracked from the start and re-included by the negation, so the walk enters
	// frontend/dist/ in every reading of this tree: the directory that gets pruned
	// is the one BELOW it, which is the shape the acceptance measured.
	seedFile(t, root, "frontend/dist/.gitkeep", "")
	gitIndexFixture(t, root)
	gitCommitAll(t, root, "the tree as it stands while the first read is taken")

	// READ 1, word for word what runGitIndex issues.
	read1, err := runGitAt(root, "ls-files", "-z")
	if err != nil {
		t.Fatalf("read 1: %s", err)
	}

	// THE COMMIT THAT LANDS IN THE WINDOW. `shipped.tsx` is force-added because
	// that is the only way a path under a live rule enters the index (A214/F1), and
	// it breaks ban #8, so "was it scanned?" has a visible answer. `build-output.tsx`
	// is the same bytes and stays untracked: a normal run must still exclude it.
	seedFile(t, root, stray, "export const doc = \"ready \u2264\";\n")
	seedFile(t, root, "frontend/dist/assets/build-output.tsx", "export const doc = \"ready \u2264\";\n")
	gitForceAdd(t, root, stray)
	gitCommitAll(t, root, "lands between the two reads")

	// READ 2.
	read2, err := runGitAt(root, "ls-files", "-z", "-i", "-c", "--exclude-standard")
	if err != nil {
		t.Fatalf("read 2: %s", err)
	}

	// 1. THE PREMISE, measured on the two raw streams before any skip logic sees
	//    them. If either half stops holding, the whole case below is vacuous and
	//    this block is what says so.
	tracked1, dirs1, err := parseIndexPaths(read1)
	if err != nil {
		t.Fatalf("read 1 does not parse: %s", err)
	}
	narrow2, err := parseIndexPathsList(read2)
	if err != nil {
		t.Fatalf("read 2 does not parse: %s", err)
	}
	sawStray := false
	for _, p := range narrow2 {
		if p == stray {
			sawStray = true
		}
	}
	if !sawStray {
		t.Fatalf("the narrow read does not name %s (got %v) - the seed no longer produces a tracked path under a live rule, "+
			"so there is no disagreement for the guard to find and this test would pass on an empty premise", stray, narrow2)
	}
	if tracked1[stray] || dirs1["frontend/dist/assets"] {
		t.Fatalf("read 1 already holds the new path (tracked=%v dirs=%v) - the commit did not land in the window, "+
			"so this pair is not the race shape and a passing test here would prove nothing", tracked1[stray], dirs1["frontend/dist/assets"])
	}

	// 2. THE GUARD: those two reads describe two different trees, and assembling
	//    them is not this tool's call to make. Red before buildGitIndexState
	//    checked the subset relation (ticket 142 AC#1/AC#2).
	racy := buildGitIndexState(root, read1, read2)
	if racy.ok {
		t.Errorf("an index that moved between the two reads was ACCEPTED (holds(%q, dir)=%v, holds(file)=%v): the walks "+
			"prune a directory the index holds on the strength of a read that never saw it - A207/A214's under-scan "+
			"direction, silently", "frontend/dist/assets", racy.holds("frontend/dist/assets", true), racy.holds(stray, false))
	}
	for _, want := range []string{"changed between the two", stray} {
		if !strings.Contains(racy.why, want) {
			t.Errorf("the refusal must name %q, got %q", want, racy.why)
		}
	}
	if strings.Contains(racy.why, "\n") {
		t.Errorf("the reason goes on the scan's first line, so it must be one line, got %q", racy.why)
	}

	// 3. CONTROL AGAINST "ALWAYS REFUSES": the same narrow stream against a full
	//    read taken AFTER the commit is a quiescent pair, it assembles, and its
	//    holds() answers are the ones A214/F1 demands. This is also the reading
	//    that shows what the first read was missing.
	read3, err := runGitAt(root, "ls-files", "-z")
	if err != nil {
		t.Fatalf("read 3: %s", err)
	}
	quiet := buildGitIndexState(root, read3, read2)
	if !quiet.ok {
		t.Fatalf("a quiescent pair was refused (%q) - the subset check fires on more than the race and every real "+
			"scan would lose its ignore filter", quiet.why)
	}
	if !quiet.holds("frontend/dist/assets", true) || !quiet.holds(stray, false) {
		t.Errorf("holds() is wrong on a quiescent index (dir=%v file=%v) - A214/F1's guarantee is the baseline this "+
			"ticket adds to, not the one it replaces", quiet.holds("frontend/dist/assets", true), quiet.holds(stray, false))
	}

	// 4. THE DECISION skip() makes with each state, and the live rule the racy one
	//    refuses to apply.
	racyMatcher := newGitIgnore(root)
	racyMatcher.idx = racy
	if v := racyMatcher.decide("frontend/dist/assets", true); !v.ignored {
		t.Fatal("no rule matches frontend/dist/assets in this seed, so block 4's 'skipped nothing' would be a tautology")
	}
	if racyMatcher.skip(filepath.Join(root, "frontend/dist/assets"), true) {
		t.Errorf("skip(dir) pruned a directory on a torn index - the guard was not reached")
	}
	if racyMatcher.skip(filepath.Join(root, stray), false) {
		t.Errorf("skip(file) excluded a delivered byte on a torn index")
	}
	n := racyMatcher.note()
	if !strings.HasPrefix(n, "d22scan: gitignore rules NOT APPLIED") {
		t.Errorf("the mechanism must be the FIRST line of the self-report (A218(7)), got %q", n)
	}
	for _, want := range []string{"changed between the two", "every path in every scope is being scanned"} {
		if !strings.Contains(n, want) {
			t.Errorf("note must say %q, got %q", want, n)
		}
	}
	if strings.Contains(n, "\n") {
		t.Errorf("this run skipped nothing, so the note must be exactly the one loud line, got %q", n)
	}

	// 5. END TO END, over the walks main() actually runs, same tree, two states.
	//    The quiescent half is production's own path (scanFixture -> runGitIndex),
	//    so the pair below is the ticket's claim in one line: a torn read now costs
	//    the same coverage as a clean one, and says so.
	sQuiet := scanFixture(t, root)
	sRacy, err := scanWithIgnore(root, racyMatcher)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := sQuiet.examined["panel-approval"], 4; got != want {
		t.Fatalf("the quiescent scan counted %d frontend/ text files, want %d (panel.tsx + frontend/.gitignore + "+
			"dist/.gitkeep + shipped.tsx, with the untracked build output excluded) - the seed moved, so the deltas "+
			"below measure nothing", got, want)
	}
	if got, want := sRacy.examined["panel-approval"], sQuiet.examined["panel-approval"]+1; got != want {
		t.Errorf("ban #6 counted %d frontend/ files on the torn index and %d on a clean one, want +1 - the torn read "+
			"must fail toward MORE scanning, never fewer", got, want)
	}
	if got, want := sRacy.emojiSeen["frontend/"], sQuiet.emojiSeen["frontend/"]+1; got != want {
		t.Errorf("ban #8 counted %d frontend/ files on the torn index and %d on a clean one, want +1", got, want)
	}
	paths := func(s *scanner) map[string]bool {
		out := map[string]bool{}
		for _, f := range s.findings {
			if f.Ban == "emoji" {
				out[f.Path] = true
			}
		}
		return out
	}
	qf, rf := paths(sQuiet), paths(sRacy)
	if !qf[stray] {
		t.Errorf("ban #8 did not fire at %s on a QUIESCENT index (findings=%v) - that is A214/F1's guarantee, and it has to survive this ticket", stray, sQuiet.findings)
	}
	if qf["frontend/dist/assets/build-output.tsx"] {
		t.Errorf("the quiescent scan read the untracked build output; A207's exclusion is gone (findings=%v)", sQuiet.findings)
	}
	if !rf[stray] || !rf["frontend/dist/assets/build-output.tsx"] {
		t.Errorf("the torn-index scan must read BOTH delivered and untracked bytes under the pruned directory, got %v", sRacy.findings)
	}
	if _, errOut, code := runVerdict(t, root, sRacy); code != 1 {
		t.Errorf("verdict on the torn index must be 1 (it scanned more, and more is what it found), got %d err=%s", code, errOut)
	}
	if _, errOut, code := runVerdict(t, root, sQuiet); code != 1 {
		t.Errorf("verdict on the clean index must also be 1 (the delivered byte is a finding either way), got %d err=%s", code, errOut)
	}
	// And the clean run must not complain about this ticket's guard: the loud line
	// is for torn or unreachable indexes only, or every ordinary scan would carry
	// it and nobody would read it.
	if qn := sQuiet.ign.note(); strings.Contains(qn, "NOT APPLIED") || !strings.Contains(qn, "skipped as git-ignored") {
		t.Errorf("the quiescent scan must skip what it has always skipped and say only that, got %q", qn)
	}
}

// TestIgnoreRulesAreOffWhenTheIndexCannotBeAsked pins the failure direction of

// A214/F1's fix. `git archive HEAD | tar -x` - the shape every evidence run and
// mutation check in this repository is taken on - has no `.git`, so the tool
// cannot know whether the bytes in front of it are tracked. Guessing "ignored"
// there is what made a force-added deliverable invisible (acceptance §1.3, and the
// second instance at §6.9 where a bare goroutine under the root `build/` rule
// turned rc=1 into rc=0). So: no rule is applied, EVERY path is scanned, and the
// run says out loud that it could not ask. A quiet green in this shape is a failed
// cell, which is why the last block below demands the note even when the tree is
// clean.
func TestIgnoreRulesAreOffWhenTheIndexCannotBeAsked(t *testing.T) {
	// No gitIndexFixture anywhere in this test: the absence of an index IS the
	// subject. Same rules and the same two shapes as the test above.
	root := liveFixture(t)
	seedFile(t, root, ".gitignore", "frontend/dist/*\n!frontend/dist/.gitkeep\nbuild/\n")
	seedFile(t, root, "frontend/.gitignore", "dist/*\n!dist/.gitkeep\n")
	seedFile(t, root, "frontend/dist/copied-delivery.tsx", "export const doc = \"ready \u2264\";\n")
	seedFile(t, root, "internal/build/copied-leak.go", "package build\n\nconst doc = \"ready \u2264\"\n\nfunc worker() {}\n\nfunc leak() { go worker() }\n")

	s := scanFixture(t, root)

	// 1. scanned, not skipped: both ignored paths are in the denominators.
	if got, want := s.examined["panel-approval"], 3; got != want {
		t.Errorf("ban #6 examined %d frontend/ files, want %d - with no index to consult nothing may be excluded", got, want)
	}
	if got, want := s.emojiSeen["frontend/"], 3; got != want {
		t.Errorf("ban #8 examined %d frontend/ files, want %d", got, want)
	}
	if got, want := s.examined[goScopeKey(filepath.Join(root, "internal"))], 3; got != want {
		t.Errorf("bans #1-5 examined %d internal/ production Go files, want %d (internal/build/copied-leak.go included)", got, want)
	}

	// 2. and reported as findings, because counting a file and not reading it would
	//    be the same blindfold wearing a bigger number.
	saw := map[string]bool{}
	for _, f := range s.findings {
		saw[f.Path] = true
	}
	for _, want := range []string{"frontend/dist/copied-delivery.tsx", "internal/build/copied-leak.go"} {
		if !saw[want] {
			t.Errorf("the index-less tree must go red at %s, findings=%v", want, s.findings)
		}
	}

	// 3. the loud self-report: what could not be asked, and what that did to the
	//    numbers. One line, naming the reason.
	n := s.ign.note()
	for _, want := range []string{"gitignore rules NOT APPLIED", "not a git repository", "every path in every scope is being scanned"} {
		if !strings.Contains(n, want) {
			t.Errorf("note must say %q, got %q", want, n)
		}
	}
	if strings.Contains(n, "\n") {
		t.Errorf("the index note must be one line, got %q", n)
	}
	if _, _, code := runVerdict(t, root, s); code != 1 {
		t.Errorf("verdict must be 1 on findings, got %d", code)
	}

	// 4. the shape a quiet green would take: a clean index-less tree with paths the
	//    rules would have skipped. Green is allowed; silence is not.
	clean := liveFixture(t)
	seedFile(t, clean, ".gitignore", "frontend/dist/*\n!frontend/dist/.gitkeep\n")
	seedFile(t, clean, "frontend/dist/bundle.js", "export const Panel = () => null;\n")
	sc := scanFixture(t, clean)
	out, errOut, code := runVerdict(t, clean, sc)
	if code != 0 {
		t.Fatalf("a tree with no banned shape must still be 0, got %d out=%s err=%s", code, out, errOut)
	}
	if nc := sc.ign.note(); !strings.Contains(nc, "gitignore rules NOT APPLIED") {
		t.Errorf("an index-less run that reports clean must still say the rules were off, got %q", nc)
	}
	if got := sc.examined["panel-approval"]; got != 2 {
		t.Errorf("ban #6 examined %d frontend/ files, want 2 (panel.tsx + bundle.js counted, nothing skipped)", got)
	}
}

// TestGitIgnoreRuleSemantics pins the matcher against gitignore(5) directly, path
// by path, because the counts depend on the fine print: a matcher that excluded
// the `dist` DIRECTORY instead of its contents would drop
// `frontend/dist/.gitkeep` and move CI's 40 to 39, and one that could not read a
// negation would do the same thing for the wrong reason.
func TestGitIgnoreRuleSemantics(t *testing.T) {
	root := t.TempDir()
	seedFile(t, root, ".gitignore", strings.Join([]string{
		"# a comment line",
		"",
		"dist/*",
		"!dist/.gitkeep",
		"build/",
		"logs",
		"*.log",
		"**/gen/*.go",
		"secret?.txt",
		"a/b/mid.txt",
		"notes/blocked.md",
		// A214/F3's two shapes, both of which r1 got in the SILENT direction: a
		// trailing `**` that pruned its own parent directory, and a rule line whose
		// leading whitespace was trimmed into a different pattern.
		"deep/**",
		"!deep/.gitkeep",
		"  lead/*",
	}, "\n")+"\n")
	// The deepest rule file wins, so this negation re-includes a path the root
	// file excludes.
	seedFile(t, root, "notes/.gitignore", "!blocked.md\n")

	g := newGitIgnore(root)
	cases := []struct {
		rel   string
		isDir bool
		want  bool
		why   string
	}{
		{"dist/index.html", false, true, "`dist/*` excludes build output"},
		{"dist/.gitkeep", false, false, "the negation keeps the go:embed anchor - THIS is what preserves CI's 40"},
		{"dist", true, false, "the anchor's parent directory is NOT itself excluded, so the walk must enter it"},
		{"dist/assets", true, true, "`dist/*` matches the directory, so it is pruned"},
		{"dist/assets/b.js", false, true, "an ignored directory takes its contents with it"},
		{"web/build", true, true, "non-anchored `build/` matches a directory at any depth"},
		{"web/build/x.go", false, true, "via the excluded parent"},
		{"src/build", false, false, "`build/` is directory-only, so a FILE named build stays scanned"},
		{"web/logs", true, true, "`logs` matches the base name at depth"},
		{"x.log", false, true, "`*.log` at the root"},
		{"sub/y.log", false, true, "`*.log` at depth"},
		{"sub/gen/a.go", false, true, "`**/gen/*.go` crosses directories"},
		{"gen/a.go", false, true, "a leading `**/` matches zero directories"},
		{"secret1.txt", false, true, "`?` matches exactly one character"},
		{"secret.txt", false, false, "`?` requires a character there"},
		{"a/b/mid.txt", false, true, "a pattern with an internal `/` is anchored to its own directory"},
		{"c/a/b/mid.txt", false, false, "and therefore does NOT match at depth"},
		{"notes/blocked.md", false, false, "the deeper notes/.gitignore negation outranks the root exclusion"},
		{"comment", false, false, "a `#` line is a comment, not a rule"},
		{"src/app.tsx", false, false, "nothing here touches a tracked source file"},
		// A214/F3 ②: `deep/**` must not be allowed to prune `deep` itself.
		{"deep", true, false, "a TRAILING `**` never matches its own parent directory (measured: `git check-ignore -v dist` exits 1 for the rule `dist/**`) - pruning it swallows the re-included anchor below and moves CI's 40 to 39"},
		{"deep/.gitkeep", false, false, "the negation still reaches it, because nothing above pruned the directory"},
		{"deep/x", true, true, "`deep/**` does ignore the directories INSIDE deep"},
		{"deep/x/y.js", false, true, "and everything below them"},
		// A214/F3 ①: leading whitespace in a rule line is data.
		{"lead/inside.txt", false, false, "git does NOT strip leading whitespace, so `  lead/*` is a pattern beginning with two spaces and binds nothing here (r1's TrimSpace made it swallow live/lead/inside.txt)"},
		{"  lead/inside.txt", false, true, "and it does bind the directory whose name really starts with those two spaces"},
	}
	for _, c := range cases {
		if got := g.decide(c.rel, c.isDir); got.ignored != c.want {
			t.Errorf("decide(%q, isDir=%v) = ignored %v, want %v (%s)", c.rel, c.isDir, got.ignored, c.want, c.why)
		}
	}

	// And the origin is reported, not just the verdict: `dist/index.html` is
	// excluded by the ROOT file here, and the note has to say so.
	if v := g.decide("dist/index.html", false); v.origin != ".gitignore" {
		t.Errorf("origin = %q, want .gitignore", v.origin)
	}
	if v := g.decide("notes/anything-else.md", false); v.ignored {
		t.Errorf("notes/.gitignore only re-includes blocked.md; got %v", v)
	}

	// nil receiver and the walk's own root must never skip: skipping the starting
	// directory would abort the walk before it examined anything, which is the
	// empty-instrument failure shape ticket 71 AC#4 made fatal.
	if (*gitIgnore)(nil).skip(filepath.Join(root, "dist"), true) {
		t.Error("a nil matcher must not skip anything")
	}
	if newGitIgnore(root).skip(root, true) {
		t.Error("the walk's starting directory must never be skipped")
	}
}

// relToRepo is the walked path as a slash path relative to the repository root,
// which is the form .gitignore patterns are written against.
func relToRepo(root, p string) string {
	rel, err := filepath.Rel(root, p)
	if err != nil {
		return filepath.ToSlash(p)
	}
	return filepath.ToSlash(rel)
}

// ignoredLikeGit is the hand-written second copy of this repository's ignore
// policy for the four scanned trees, used ONLY by the independent walk above.
//
// Measured at anchor 7771409 (`git status --porcelain --ignored=matching` on a
// worktree that had run a frontend build, cross-checked with
// `git ls-files -i -c --exclude-standard` -> empty, i.e. NO tracked file is
// excluded by any rule): inside a scanned scope the only exclusions are
// `frontend/node_modules/` (skipped by directory name above, as it always was)
// and `frontend/dist/*` minus the `!frontend/dist/.gitkeep` go:embed anchor.
// Everything else .gitignore names (build/, third_party/, *.exe, the
// scripts/spike/bin/ dumps) lives outside design/ frontend/ internal/ cmd/.
//
// A214/F1 ADDS THE INDEX HALF TO THIS COPY, IN THIS BATCH, ON PURPOSE. A pattern
// alone was exactly the premise that let the r1 batch hide a force-added
// deliverable, and this copy shares every premise of the thing it is supposed to
// falsify - so as delivered it agreed with the instrument on the very tree the
// acceptance had built by hand (§1.4: "it returns true for
// frontend/dist/tracked-forced.tsx"). Two arguments now bind it to git's criterion
// instead of to a pattern:
//   - rulesOn mirrors gitignore.go's fallback: when the index cannot be consulted
//     the instrument applies NO rule, so this copy must skip nothing either, or the
//     two numbers split on every `git archive` snapshot for the wrong reason.
//   - tracked is the index. A path git holds is never ignored, so this copy may not
//     skip it either. If gitignore.go ever forgets that, the two copies disagree by
//     exactly one file and the failure message is the "the number in the self-report
//     is wrong" line in the test above - the designed friction, not a bug.
//
// What is NOT mirrored and why: the two A214/F3 semantic fixes (leading whitespace,
// trailing `**`) live in the matcher this list does not call. The list encodes
// today's policy, and today's policy has no rule whose meaning those fixes change
// inside design/ frontend/ internal/ cmd/ - so if such a rule is ever added, this
// list moves in the same batch, exactly as the paragraph above has always demanded.
//
// If this list and .gitignore ever disagree, the failure message is
// TestLedgerCountsMatchAnIndependentWalk's "the number in the self-report is
// wrong" line - read it as "the ignore policy now covers a scanned scope", decide
// whether that is intended, move this list, and re-state the scope baselines in
// docs/reports/pending-and-issues.md in the same batch (A207 ②).
func ignoredLikeGit(rel string, tracked map[string]bool, rulesOn bool) bool {
	if !rulesOn {
		return false
	}
	if strings.HasPrefix(rel, "frontend/dist/") {
		if tracked[rel] {
			return false
		}
		return filepath.Base(rel) != ".gitkeep"
	}
	return false
}

// gitTrackedIndex is this file's own hand-asked copy of the three git questions
// gitignore.go's runGitIndex asks, so the oracle above is not secretly calling the
// code under test. It never fails the test: an unreachable index is the normal
// state of a `git archive` snapshot, and both sides then agree on "no rules".
func gitTrackedIndex(root string) (tracked map[string]bool, ok bool, why string) {
	tracked = map[string]bool{}
	for _, step := range []struct {
		name string
		args []string
	}{{"rev-parse", []string{"rev-parse", "--show-prefix"}}, {"ls-files", []string{"ls-files", "-z"}}} {
		cmd := exec.Command("git", step.args...)
		cmd.Dir = root
		var stdout, stderr strings.Builder
		cmd.Stdout, cmd.Stderr = &stdout, &stderr
		if err := cmd.Run(); err != nil {
			return tracked, false, firstLineOf(stderr.String()) + " (`git " + step.name + "`)"
		}
		if step.name == "rev-parse" {
			if p := strings.TrimRight(stdout.String(), "\r\n"); p != "" {
				return tracked, false, "root is the subdirectory " + p + " of a repository, not its top"
			}
			continue
		}
		for _, p := range strings.Split(stdout.String(), "\x00") {
			if p != "" {
				tracked[p] = true
			}
		}
	}
	return tracked, true, ""
}

// TestRealRepoLedgerIsHonest is AC#4 on this repository, restated by ticket 88
// for the armed ban #6: no live scope is empty, no exemption has drifted, no
// counter is undeclared - and there is no longer an "only uncovered scope" to
// confess to, so the honest thing this test pins instead is that ban #6 really
// reads the panel's files. Pre-ticket-88 it asserted the opposite number
// (`examined == 0`), which is what made the flip cost it.
func TestRealRepoLedgerIsHonest(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Skipf("not inside the wisp repo: %v", err)
	}
	s := scanFixture(t, root)
	scopes := declaredScopes(root)

	if got := emptyLiveScope(scopes, s); got != "" {
		t.Errorf("live scope %s walks 0 files in the real repo: give it coverage or delete it (ticket 71 AC#4)", got)
	}
	if got := driftedAbsentScope(scopes); got != "" {
		t.Errorf("scope %s is registered absent while its tree exists - flip it to live:true (ticket 71 AC#4)", got)
	}
	if extra := undeclaredKeys(s, scopes); len(extra) != 0 {
		t.Errorf("walks bumped counters that no scope reports: %v", extra)
	}
	var unc []string
	for _, sc := range uncoveredScopes(scopes) {
		unc = append(unc, sc.label)
	}
	if len(unc) != 0 {
		t.Errorf("ticket 88 armed ban #6, so this HEAD declares no uncovered scope; found %v - an exemption needs a tree that provably does not exist, and one that exists must be live or the run must say it is not covered", unc)
	}
	// AC#1's number, as a gate rather than a log line: this is the assertion that
	// catches "armed but blind" - live:true with a filter that never matches the
	// panel's file classes would still print a verdict.
	if n := s.examined["panel-approval"]; n == 0 {
		t.Error("ban #6 is live yet examined 0 files under frontend/: an armed empty instrument is worse than an honest NOT COVERED")
	} else {
		t.Logf("ban #6 examined %d frontend/ text files in the real repo", n)
	}
	out, errOut, code := runVerdict(t, root, s)
	if code != 0 {
		t.Fatalf("HEAD must be green, rc=%d out=%s err=%s", code, out, errOut)
	}
	if !strings.Contains(out, "clean") {
		t.Errorf("the clean line must be there, got %q", out)
	}
	if strings.Contains(out, "NOT COVERED") {
		t.Errorf("nothing is uncovered on this HEAD, so the verdict must not claim it is: %q", out)
	}
	for _, sc := range scopes {
		if !sc.live {
			continue
		}
		want := fmt.Sprintf("%s=%d", sc.label, sc.count(s))
		if !strings.Contains(out, want) {
			t.Errorf("clean line omits live scope %s with its real count (%q)", sc.label, want)
		}
	}
	if !strings.Contains(out, fmt.Sprintf("ban #6 frontend/=%d", s.examined["panel-approval"])) {
		t.Errorf("the clean line must carry ban #6's real count, got %q", out)
	}
}

// e2eRoot is a fixture that clears checkRoot (which wants >= minProductionGoFiles
// production .go files), so the binary reaches the guards instead of dying on the
// root sanity check first.
func e2eRoot(t *testing.T) string {
	t.Helper()
	root := liveFixture(t)
	for i := 0; i < 12; i++ {
		seedFile(t, root, fmt.Sprintf("internal/pkg%02d/f.go", i), "package pkg\n")
	}
	return root
}

// TestBuiltBinaryGoesRedEndToEnd is the positive control in the shape the CI step
// actually consumes: it compiles the real binary and asserts the PROCESS exit
// code. Everything above tests verdict()'s return value, which is worthless if
// main() ever stops passing it to os.Exit - a `_ = verdict(...)` regression would
// leave every unit test green while the lint job goes green on a red tree, which
// is ticket 71's disease in its purest form.
func TestBuiltBinaryGoesRedEndToEnd(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skipf("no go toolchain to build the scanner with: %v", err)
	}
	bin := filepath.Join(t.TempDir(), "d22scan-e2e"+exeSuffix())
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = "."
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\n%s", err, out)
	}

	cases := []struct {
		name    string
		setup   func(t *testing.T, root string)
		wantRC  int
		wantAll []string
	}{
		{
			name: "seeded violation exits 1",
			setup: func(t *testing.T, root string) {
				seedFile(t, root, "internal/bad/leak.go", "package bad\n\nfunc worker() {}\n\nfunc leak() { go worker() }\n")
			},
			wantRC:  1,
			wantAll: []string{"bare-goroutine"},
		},
		{
			name: "empty live scope exits 2",
			setup: func(t *testing.T, root string) {
				if err := os.Remove(filepath.Join(root, "internal", "tools", "ok.go")); err != nil {
					t.Fatal(err)
				}
			},
			wantRC:  2,
			wantAll: []string{"ban #7 internal/tools/", "empty instrument"},
		},
		// Ticket 88 AC#4: the positive control for the armed ban #6, through the
		// process CI runs. The seed is the same file and the same line the ticket
		// 71 case used (frontend/src/app.js, `approval.decide(...)`) - the panel
		// is .tsx in production, and TestBan6ScopeIsNotNarrowedByAnExtensionFilter
		// covers that class; renaming the seed to dodge a filter question would be
		// the AC#1 hole written into a test.
		{
			name: "armed ban 6 goes red on the panel violation exits 1",
			setup: func(t *testing.T, root string) {
				seedFile(t, root, "frontend/src/app.js", "export function decide() { return approval.decide({allow: true}); }\n")
			},
			wantRC:  1,
			wantAll: []string{"panel-approval", "approval.decide", "ban #6 frontend/", "frontend/src/app.js"},
		},
		// And the other direction of the flip: ban #6 is now subject to the rule
		// it was exempt from. Losing frontend/ is a coverage change, so it stops
		// the verdict instead of reading as "the panel got clean".
		//
		// The case this REPLACED ("exempt scope whose tree appeared exits 2") is
		// gone from the built-binary list because nothing in declaredScopes() is
		// exempt any more; the only way to give the binary an exemption would have
		// been to register a fake one in production code, which is the disease, not
		// the test. Guard 3's teeth moved to
		// TestExemptScopeCannotOutliveItsAbsentTree, which drives the same verdict()
		// decision against a synthetic exempt scope in both directions.
		{
			name: "frontend tree gone while declared live exits 2",
			setup: func(t *testing.T, root string) {
				if err := os.RemoveAll(filepath.Join(root, "frontend")); err != nil {
					t.Fatal(err)
				}
			},
			// RESTATED BY TICKET 96 (AC#4): the assertion used to be
			// {"ban #6 frontend/", "empty instrument"}, i.e. guard 2's wording.
			// It still measures "losing a live tree stops the verdict" - rc=2 and
			// the tree named - but which guard produces that changed, because
			// frontend/ is now declared by ban #8 too and guard 1 (a declared ban
			// #8 scope that walked 0 files) deliberately outranks guard 2. So the
			// case is no longer reachable for ban #6 by deleting the tree; the
			// "empty instrument" wording below keeps that job pinned in a case
			// ban #8 does not shadow.
			wantRC:  2,
			wantAll: []string{"ban #8 frontend/", "examined 0 files"},
		},
		{
			name: "ban 7 tree gone while declared live exits 2",
			setup: func(t *testing.T, root string) {
				if err := os.RemoveAll(filepath.Join(root, "internal", "tools")); err != nil {
					t.Fatal(err)
				}
			},
			wantRC:  2,
			wantAll: []string{"ban #7 internal/tools/", "empty instrument"},
		},
		{
			name:    "fully live fixture exits 0",
			setup:   func(_ *testing.T, _ string) {},
			wantRC:  0,
			wantAll: []string{"clean", "ban #6 frontend/="},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			root := e2eRoot(t)
			c.setup(t, root)
			cmd := exec.Command(bin, "-root", root)
			var stdout, stderr bytes.Buffer
			cmd.Stdout, cmd.Stderr = &stdout, &stderr
			err := cmd.Run()
			rc := 0
			if err != nil {
				var ee *exec.ExitError
				if !errors.As(err, &ee) {
					t.Fatalf("running the scanner failed: %v", err)
				}
				rc = ee.ExitCode()
			}
			t.Logf("rc=%d stdout=%sstderr=%s", rc, stdout.String(), stderr.String())
			if rc != c.wantRC {
				t.Fatalf("exit code: want %d, got %d\nstdout:\n%s\nstderr:\n%s", c.wantRC, rc, stdout.String(), stderr.String())
			}
			joined := stdout.String() + stderr.String()
			for _, want := range c.wantAll {
				if !strings.Contains(joined, want) {
					t.Errorf("output must name %q, got:\n%s", want, joined)
				}
			}
		})
	}
}

func exeSuffix() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}
