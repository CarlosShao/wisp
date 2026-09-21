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
