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

// TestLedgerCountsMatchAnIndependentWalk is AC#4's falsifier: the printed numbers
// are compared against a SECOND, independent walk computed in the test. Without
// it a ledger that counts the wrong thing (every .go including _test.go, or a
// directory skipped by mistake) stays internally consistent and still prints a
// confident number - "examined N" only means something if N is the truth.
func TestLedgerCountsMatchAnIndependentWalk(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Skipf("not inside the wisp repo: %v", err)
	}
	s := scanFixture(t, root)

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
			name: "ban 6 tree gone while declared live exits 2",
			setup: func(t *testing.T, root string) {
				if err := os.RemoveAll(filepath.Join(root, "frontend")); err != nil {
					t.Fatal(err)
				}
			},
			wantRC:  2,
			wantAll: []string{"ban #6 frontend/", "empty instrument"},
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
