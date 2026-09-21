package main

import (
	"os"
	"path/filepath"
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
// in code: ban #8 must fire inside internal/ AND inside cmd/, on a COMMENT
// line, on a string literal and in a _test.go file. Before ticket 67 AC#3's
// coverage step all four of these seeds walked past the gate unnoticed, because
// the only declared scopes were design/ and frontend/ - and frontend/ does not
// exist in this repository, so every .go file in the product was invisible to
// the ban while the footer still claimed "no emoji in design/ or frontend/".
func TestEmojiBanCoversGoSourcesNotJustDesign(t *testing.T) {
	root := t.TempDir()
	seedFile(t, root, "tools/d22scan/allowlist.txt", "# empty\n")

	// U+2713 CHECK MARK: inside emojiRe's \x{2600}-\x{27BF} band, and the exact
	// code point that lived in cmd/wisp/providers.go until fff4cad.
	seedFile(t, root, "internal/ok/comment.go", "package ok\n\n// a verdict reads \u2713 here\n")
	seedFile(t, root, "internal/ok/literal.go", "package ok\n\nconst banner = \"ready \u2713\"\n")
	seedFile(t, root, "internal/ok/emoji_test.go", "package ok\n\n// \u2713 in a test file comment\nfunc TestNothing(*testing.T) {}\n")
	seedFile(t, root, "cmd/wisp/glyph.go", "package main\n\n// \u2713 in cmd/\nfunc main() {}\n")

	// Control for the documented goOnly rule: a non-.go file under internal/ is
	// testdata goldens / leaked test debris, not source, so it is out of scope.
	seedFile(t, root, "internal/ok/testdata/fixture.txt", "ready \u2713\n")

	findings, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"internal/ok/comment.go",    // comment line, internal/
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
	// cmd/ deliberately absent.

	s, err := scanWithStats(root)
	if err != nil {
		t.Fatal(err)
	}
	if s.emojiSeen["design/"] == 0 || s.emojiSeen["internal/"] == 0 {
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
	if strings.Contains(report, "frontend") {
		t.Errorf("the verdict line mentions frontend/ while emojiScopes() has no such entry: %q", report)
	}
}
