package main

// Ticket 350's own tests: the bounded wordlist gate (ban #10), AC#1③.
//
// Everything here runs on THROWAWAY fixtures (t.TempDir()) or on blobs read
// out of this repository's git objects - never by editing the shared tree: the
// start-anchor control (AC#1①) materialises the committed bytes of
// internal/risk/blacklist.go at the START anchor 29d258f9 into a fixture,
// re-scans, and requires the red to be EXACTLY aEquiv; the delivery-side
// control (AC#1②) is the real tree's own scan, which must examine >0 risk
// files and stay green.
//
// THE CONTROL TREE, stated plainly (measured, not assumed): these tests were
// written and run inside a `git clone --shared` of the main worktree at HEAD
// d0a230547 (rule-21 isolation; the main tree carried no in-flight file under
// tools/d22scan or docs/evidence/s1/350-* at leg start). The git-object
// controls therefore read internal/risk/blacklist.go AS COMMITTED at HEAD and
// at 29d258f9 - which is exactly what AC#1① wants judged (the committed
// tree's own state) - and they stay stable no matter what other legs hold
// uncommitted in the shared tree.

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// repoRoot returns the repository root the way scan_test.go's real-repo tests
// do: two levels up from this package directory. In the rule-21 clone that is
// the clone itself; in a detached tools/d22scan checkout it is the repo above.
func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Skipf("not inside the wisp repo: %v", err)
	}
	return root
}

// writeBlobAt materialises <base>/<rel> with the EXACT bytes of <rev>:<rel>
// from the repository at repoRoot (git show), creating parent dirs. The base
// is a t.TempDir() fixture, never the repository: the blob is planted in a
// throwaway tree and the shared tree is never edited for a control.
func writeBlobAt(t *testing.T, repoRoot, rev, rel, base string) {
	t.Helper()
	out, err := exec.Command("git", "-C", repoRoot, "show", rev+":"+rel).Output()
	if err != nil {
		t.Fatalf("git show %s:%s: %v", rev, rel, err)
	}
	dst := filepath.Join(base, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dst, out, 0o644); err != nil {
		t.Fatal(err)
	}
}

// gateFixture is the minimum declaration surface the scanner needs to build a
// verdict at all (every live scope doing real work), plus whatever the caller
// plants under internal/risk/. The risk file is the caller's variable, so the
// ban #10 scope's work comes from the case, not from a constant.
func gateFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	seedFile(t, root, "tools/d22scan/allowlist.txt", "# empty\n")
	seedFile(t, root, "go.mod", "module x\n")
	seedFile(t, root, "internal/ok/ok.go", "package ok\n")
	seedFile(t, root, "internal/tools/ok.go", "package tools\n")
	seedFile(t, root, "cmd/wisp/main.go", "package main\n\nfunc main() {}\n")
	seedFile(t, root, "design/index.html", "<html>ok</html>\n")
	seedFile(t, root, "frontend/src/panel.tsx", "export const Panel = () => null;\n")
	return root
}

// checkWalkRan is the anti-hollow gate for these cases: if the fixture fails
// to seed, ban #10's counter staying 0 would read as "clean" when it means
// "blind" - the same vacuous-reading rule the self-test's cover counters
// enforce (ticket 71 AC#4's shape, at the test layer).
func checkWalkRan(t *testing.T, s *scanner) {
	t.Helper()
	if got := s.examined["phantom-symbol"]; got == 0 {
		t.Fatal("fixture broken: ban #10 examined 0 risk files - an empty instrument must not read as a verdict")
	}
}

// phantomFindings collects one scan's phantom-symbol findings, rendered.
func phantomFindings(s *scanner) []string {
	var got []string
	for _, f := range s.findings {
		if f.Ban == "phantom-symbol" {
			got = append(got, fmt.Sprintf("%s:%d %s", filepath.ToSlash(f.Path), f.Line, f.What))
		}
	}
	return got
}

// TestRiskWordlistStartAnchorControlIsExactlyAEquiv is AC#1①: the committed
// bytes of internal/risk/blacklist.go at the START anchor 29d258f9 (where the
// ticket's own census measured the real positive) must ring on EXACTLY aEquiv
// - the comment-cited symbol that no declaration in the package defines - and
// on nothing else under this tag.
func TestRiskWordlistStartAnchorControlIsExactlyAEquiv(t *testing.T) {
	root := repoRoot(t)
	fixture := gateFixture(t)
	// blacklist.go at 29d258f9: its comments cite aEquiv (the ticket's real
	// positive) plus the package-declared aAnyForm/aEq/aRelates/aRuleReason/
	// aUnder. Its package is `risk`. The census behind "exactly one phantom
	// citation in this blob" is docs/evidence/s1/350-bounded-wordlist-gate.md.
	writeBlobAt(t, root, "29d258f9", "internal/risk/blacklist.go", fixture)

	s := scanFixture(t, fixture)
	checkWalkRan(t, s)
	hits := phantomFindings(s)
	if len(hits) != 1 {
		t.Fatalf("start-anchor control must ring exactly once under phantom-symbol, got %d: %v", len(hits), hits)
	}
	if !strings.Contains(hits[0], `"aEquiv"`) {
		t.Fatalf("start-anchor control must name aEquiv (the ticket's own positive), got: %q", hits[0])
	}
	if !strings.Contains(hits[0], "package risk") {
		t.Errorf("the finding must name the package, got: %q", hits[0])
	}
}

// TestRiskWordlistDeliveryTreeIsGreen is AC#1② on the rule-21 clone: the real
// internal/risk tree at HEAD must examine >0 production files and produce
// ZERO phantom-symbol findings.
func TestRiskWordlistDeliveryTreeIsGreen(t *testing.T) {
	root := repoRoot(t)
	s := scanFixture(t, root)
	checkWalkRan(t, s)
	if hits := phantomFindings(s); len(hits) != 0 {
		t.Fatalf("delivery tree must be green under ban #10, got %d finding(s): %v", len(hits), hits)
	}
}

// TestRiskCommentCitationsUnit covers the pure halves: the roster collector
// and the comment citation scanner, including the roster's breadth (method,
// field, local) and the shape's boundaries.
func TestRiskCommentCitationsUnit(t *testing.T) {
	const src = `package risk

// aKnown cites the package-level func; aBound cites the method; aLocal cites
// the := local; macOS and "a " prose are not identifier starts; aCustomer is
// read as one maximal token. This line names aVanished and nothing declares it.
type box struct{ aField int }

func (box) aBound() {}

func aKnown() {}

func walkOne() {
	aLocal := 1
	_ = aLocal
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "unit.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	roster := map[string]bool{}
	collectDeclaredIdents(f, roster)
	for _, want := range []string{"box", "aField", "aBound", "aKnown", "walkOne", "aLocal"} {
		if !roster[want] {
			t.Errorf("roster must hold declared ident %q; got %v", want, roster)
		}
	}
	if roster["aVanished"] {
		t.Error("roster must NOT hold aVanished (nothing declares it)")
	}

	cits := riskCommentCitations(fset, f, "unit.go")
	var tokens []string
	for _, c := range cits {
		tokens = append(tokens, c.Token)
	}
	joined := strings.Join(tokens, ",")
	for _, want := range []string{"aKnown", "aBound", "aLocal", "aVanished", "aCustomer"} {
		if !strings.Contains(joined, want) {
			t.Errorf("citations must include %q, got %v", want, tokens)
		}
	}
	if strings.Contains(joined, "aCOS") {
		t.Errorf(`macOS must not be read as an a-prefixed identifier (\b refuses it); got %v`, tokens)
	}
}

// TestRiskWordlistMaximalTokenAndBoundary pins the shape's two edge
// behaviours against a full scan, so a future "simplification" of aWordRe
// shows up as a red direction check rather than as a rereading of the comment:
// boundary prose stays silent, and a citation of a PREFIX of a declared name
// (aCustomer vs the declared aCustomerLong) still rings - the token is read
// as the maximal identifier, and an exact roster miss is a miss.
func TestRiskWordlistMaximalTokenAndBoundary(t *testing.T) {
	root := gateFixture(t)
	seedFile(t, root, "internal/risk/shape.go", "package risk\n\n"+
		"// macOS shipped this spelling and English prose says \"a Cased\" after\n"+
		"// an article; the real citation is aCustomerLong below.\n"+
		"func aCustomerLong() {}\n")
	s := scanFixture(t, root)
	checkWalkRan(t, s)
	if hits := phantomFindings(s); len(hits) != 0 {
		t.Fatalf("boundary prose must stay silent, got: %v", hits)
	}

	seedFile(t, root, "internal/risk/prefix.go", "package risk\n\n"+
		"// the verdict shape is decided by aCustomer, which is not a name this\n"+
		"// package declares.\n"+
		"func other() {}\n")
	s = scanFixture(t, root)
	checkWalkRan(t, s)
	hits := 0
	for _, what := range phantomFindings(s) {
		if strings.Contains(what, `"aCustomer"`) {
			hits++
		}
	}
	if hits != 1 {
		t.Fatalf("citing the prefix aCustomer must ring exactly once, got %d (all: %v)", hits, phantomFindings(s))
	}
}

// TestRiskWordlistOneFindingPerFileAndToken pins the emission dedupe: the same
// phantom token twice in one file is ONE finding; the same token cited from
// the package's other directory is a second finding (per-file dedupe, not
// per-scan - a second package's citation still reports).
func TestRiskWordlistOneFindingPerFileAndToken(t *testing.T) {
	root := gateFixture(t)
	seedFile(t, root, "internal/risk/twice.go", "package risk\n\n"+
		"// aMissing is cited here and\n"+
		"// aMissing is cited again on the next line.\n"+
		"func twice() {}\n")
	seedFile(t, root, "internal/risk/sub/other.go", "package risk\n\n"+
		"// aMissing again from the package's other directory.\n"+
		"func other() {}\n")
	s := scanFixture(t, root)
	checkWalkRan(t, s)
	hits := phantomFindings(s)
	if len(hits) != 2 {
		t.Fatalf("want one finding per (file, token): twice.go x1 + other.go x1 = 2, got %d: %v", len(hits), hits)
	}
}

// TestRiskWordlistSkipsTestFiles pins the range: _test.go is outside ban #10
// the same way it is outside bans #1-5, and the exclusion hides nothing at
// HEAD (measured: every a[A-Z] token the risk test comments name resolves to a
// production declaration). The seeded citation lives ONLY in the test file, so
// the scan must stay green while the counter still moves.
func TestRiskWordlistSkipsTestFiles(t *testing.T) {
	root := gateFixture(t)
	seedFile(t, root, "internal/risk/prod.go", "package risk\n\nfunc present() {}\n")
	seedFile(t, root, "internal/risk/prod_test.go", "package risk\n\n"+
		"// aOnlyInTest is cited by the test file and declared nowhere.\n"+
		"func TestX(*testing.T) {}\n")
	s := scanFixture(t, root)
	checkWalkRan(t, s)
	if hits := phantomFindings(s); len(hits) != 0 {
		t.Fatalf("a citation that exists only in a _test.go file is outside the ban, got: %v", hits)
	}
}

// TestRiskWordlistEmptyTreeIsFatalThroughTheLedger pins the ledger half: with
// internal/risk absent the run must NOT be green - the walk examined 0 files
// and guard 2 (emptyLiveScope) names ban #10's row. Without this pin, deleting
// the tree would read as "the comments got clean" instead of "the instrument
// went blind".
func TestRiskWordlistEmptyTreeIsFatalThroughTheLedger(t *testing.T) {
	root := liveFixture(t)
	// liveFixture seeds internal/risk/ok.go since ticket 350; remove the whole
	// tree and re-scan: guard 2 must fire through declaredScopes().
	if err := os.RemoveAll(filepath.Join(root, "internal", "risk")); err != nil {
		t.Fatal(err)
	}
	s := scanFixture(t, root)
	if s.examined["phantom-symbol"] != 0 {
		t.Fatalf("fixture broken: ban #10 examined %d files after tree removal", s.examined["phantom-symbol"])
	}
	var out, errOut strings.Builder
	code := verdict(&out, &errOut, root, 99, s)
	if code != 2 {
		t.Fatalf("an absent internal/risk with a live ledger row must stop the verdict (rc=2), got rc=%d stderr=%q", code, errOut.String())
	}
	if !strings.Contains(errOut.String(), "ban #10 internal/risk/") {
		t.Errorf("the fatal message must name ban #10's scope, got %q", errOut.String())
	}
}

// TestRiskWordlistEmissionTagIsRecognisable keeps the self-test machinery and
// the emission site in the same universe: the coverage mapping selftest.go
// reads must know the tag the walk emits, or every silent sample reads
// VACUOUS and the direction check is worthless.
func TestRiskWordlistEmissionTagIsRecognisable(t *testing.T) {
	if selfCoverForTag["phantom-symbol"] == "" {
		t.Fatal("selfCoverForTag must map phantom-symbol to a coverage counter, or every silent sample reads VACUOUS")
	}
}
