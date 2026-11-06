package main

// Ticket 150 AC#3 (option B): the self-ringing instrument for the
// DEFERRED(Dxx-y) family.
//
// WHAT RINGS AND WHY IT DID NOT EXIST BEFORE. AGENTS.md §1.1 states the rule
// ("DEFERRED(D-xx) 代码标记必须与 SPEC-12 §5 登记表 1:1 双向对得上") and
// SPEC-12 §5 restates it as 双向交叉核对, but no instrument checked it: the
// ban roster in main.go never mentioned the registry, and until this file a
// planted DEFERRED(D99-1) with no section-5 row - or a row whose marker was
// deleted - stayed green everywhere (the ticket's own premise, and
// docs/evidence/s1/150-verdict.md AC#3's independent finding). Option B is
// the shape the ticket approved: a Go check living in this package that reads
// the SPEC text and the code tree and compares the two sides. It adds NOTHING
// to the ban list, the judgment body, PLAN.md or any contract text; it is a
// test, and the red it prints is a `go test` failure.
//
// THE TWO SIDES, EXACTLY:
//
//	code side      - every DEFERRED(D<digits>-<digits>) occurrence in the
//	               PRODUCTION Go files of internal/ and cmd/, walked with the
//	               same conventions scanner.walkGo uses (testdata/ and .git/
//	               pruned, .go only, _test.go excluded, and the index-aware
//	               gitignore skip from gitignore.go so build output can never
//	               become this instrument's denominator - A207's lesson).
//	registry side  - every id cited in the 项 (item) cell of an ACTIVE
//	               `| DEFERRED |` row inside section 5 of
//	               docs/specs/SPEC-12-roadmap-governance.md.
//
// WHY _test.go IS EXCLUDED, AND WHY THAT IS THE RULE AND NOT A HOLE: the
// markers living in test files today (DEFERRED(D34-1), 7 occurrences in
// internal/tools/document_text*_test.go) are ABSENCE PINS - fixtures asserting
// the marker is gone after 24-r5 deleted it (docs/evidence/s1/150-deferred-
// registry-rows-impl.md §1 rejected them from the 1:1 range for exactly that
// reason, and 150-verdict.md independently re-verified). Counting them would
// ring against a debt that was already paid. The same convention every other
// check in this package uses - "production scope internal/ + cmd/ (non-test)"
// - is what keeps them out here, and TestDeferredRegistryPositiveControl/
// test_file_absence_pins_stay_out_of_range pins that exclusion on purpose,
// because a convention carried by an absent test is the blind spot this repo
// keeps paying to kill (R16#4's shape).
//
// WHY ONLY ACTIVE ROWS COUNT: the struck-through `~~DEFERRED~~` rows are
// re-adopted or re-opened history (their type cell is not exactly DEFERRED),
// and RESERVED/REJECTED rows are not this family. A DEFERRED row whose item
// cell carries no Dxx-y id (macOS 平台层, 无障碍, ...) registers a FEATURE, not
// a code marker, so it is out of range in the other direction too. The range
// is the marker family, both sides.
//
// THE CALIBER IS PER-OCCURRENCE, NOT PER-ID: section 5 registers D28-1 with
// three rows because the tree carries three occurrences of that marker
// (compress.go names them 出现处①②③, and row ② states the ruler "按出现处计" -
// ticket 150 AC#1's harvest caliber). So the comparison counts rows against
// occurrences per id: a fourth D28-1 marker without a fourth row rings, and so
// does deleting one marker while its row stays. Rows are NOT paired to
// occurrences by location: the 现量 path:line notes are prose that rots with
// any edit above the cited line, and pinning them would ring on unrelated
// commits. What the registered convention actually promises is the counts,
// and that is what this checks.
//
// WHAT IT PRINTS. A green run logs a DENOMINATOR - scope dirs, files examined,
// marker counts per id with locations, registered-row counts per id with the
// SPEC line numbers - so "1:1 holds" can never mean "nothing was looked at"
// (ticket 71 AC#4's rule at the registry layer), and a red run repeats that
// line inside the failure. An instrument whose two sides are BOTH empty is
// fatal, never vacuously green.

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// deferredSpecRel is the registry's home: SPEC-12 §5. The row side of the
// comparison reads its 推迟项登记表 section and nothing else.
const deferredSpecRel = "docs/specs/SPEC-12-roadmap-governance.md"

// deferredMarkerRe matches one DEFERRED(Dxx-y) code marker. It runs over RAW
// lines (comments included - every marker in this tree is a comment), so a
// marker spelled with a space inside the parens would NOT match; the row side
// then rings for the now-unmatched id, which is the loud direction.
var deferredMarkerRe = regexp.MustCompile(`DEFERRED\((D\d+-\d+)\)`)

// deferredRowIDRe extracts the registered ids from a row's item cell. The
// word-boundary guards keep date stamps (2026-11-06) and plain decision
// numbers (D28, D11③) out of the harvest; only the D<digits>-<digits> shape
// the markers use can match.
var deferredRowIDRe = regexp.MustCompile(`\bD\d+-\d+\b`)

// deferredSide is one side of the ledger: id -> occurrence count, and id ->
// human-readable locations (code: repo-relative path:line; registry: SPEC
// path:line), so every failure message and the denominator can show WHICH
// occurrences are being counted instead of a bare number.
type deferredSide struct {
	count map[string]int
	where map[string][]string
}

func newDeferredSide() *deferredSide {
	return &deferredSide{count: map[string]int{}, where: map[string][]string{}}
}

func (d *deferredSide) add(id, location string) {
	d.count[id]++
	d.where[id] = append(d.where[id], location)
}

// total is the occurrence count across all ids - the number the denominator
// prints and the empty-instrument guard reads.
func (d *deferredSide) total() int {
	n := 0
	for _, c := range d.count {
		n += c
	}
	return n
}

func (d *deferredSide) ids() []string {
	out := make([]string, 0, len(d.count))
	for id := range d.count {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}

// render spells out "ID xN [loc; loc]" per id in sorted order.
func (d *deferredSide) render() string {
	parts := make([]string, 0, len(d.count))
	for _, id := range d.ids() {
		parts = append(parts, fmt.Sprintf("%s x%d [%s]", id, d.count[id], strings.Join(d.where[id], ", ")))
	}
	return strings.Join(parts, "; ")
}

// deferredRowsFromSpec parses the ACTIVE `| DEFERRED |` rows of SPEC-12
// section 5 (root-relative specRel, so locations print repo-relative) and
// harvests the ids each row's 项 (item) cell registers. It returns the ledger
// plus the number of active rows examined - the registry side's own
// denominator, which is deliberately larger than the id count: most deferred
// rows register features without a marker, and that is legal.
//
// The id is read from the ITEM cell only (cells[2] of the split row), not from
// the whole line: the 依据/前置/残缺 cells cite OTHER rows' ids when arguing
// ("与出现处①同一迁移", "既有 REJECTED ... 那行"), and counting citations as
// registrations would inflate the registry side into false reds. A row that
// mentions the same id twice still registers one occurrence (rows count, per
// the section's own 出现处-per-row convention).
func deferredRowsFromSpec(t *testing.T, root, specRel string) (side *deferredSide, activeRows int) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(specRel)))
	if err != nil {
		t.Fatalf("registry side is blind: cannot read %s: %v", specRel, err)
	}
	lines := strings.Split(string(data), "\n")

	start := -1
	for i, ln := range lines {
		if strings.HasPrefix(strings.TrimRight(ln, "\r"), "## 5.") {
			start = i
			break
		}
	}
	if start < 0 {
		t.Fatalf("registry side is blind: no '## 5.' heading in %s - the section this instrument reads moved or was renamed, and reading 0 rows would pass vacuously", specRel)
	}
	end := len(lines)
	for i := start + 1; i < len(lines); i++ {
		if strings.HasPrefix(strings.TrimRight(lines[i], "\r"), "## ") {
			end = i
			break
		}
	}

	side = newDeferredSide()
	for i := start; i < end; i++ {
		line := strings.TrimRight(lines[i], "\r")
		if !strings.HasPrefix(line, "|") {
			continue
		}
		cells := strings.Split(line, "|")
		if len(cells) < 3 {
			continue
		}
		// Exactly "DEFERRED" - which excludes the header (类型), the separator
		// (---), RESERVED/REJECTED rows, and the struck-through
		// `~~DEFERRED~~` rows that record re-adopted or re-opened history.
		if strings.TrimSpace(cells[1]) != "DEFERRED" {
			continue
		}
		activeRows++
		seen := map[string]bool{}
		for _, id := range deferredRowIDRe.FindAllString(cells[2], -1) {
			if seen[id] {
				continue
			}
			seen[id] = true
			side.add(id, fmt.Sprintf("%s:%d", specRel, i+1))
		}
	}
	return side, activeRows
}

// deferredMarkersInDirs walks the PRODUCTION Go files of the named root-relative
// directories and harvests every DEFERRED(Dxx-y) occurrence. The skip chain is
// scanner.walkGo's (main.go), verbatim: testdata/ and .git/ pruned; .go only;
// _test.go excluded (the absence-pin convention, see this file's header); and
// gitignore.go's index-aware skip on directories and files so ignored build
// output never enters the denominator while a force-added tracked file still
// does (A207/A214). Returns the ledger and the number of files examined.
func deferredMarkersInDirs(t *testing.T, root string, dirs []string) (side *deferredSide, examined int) {
	t.Helper()
	ign := newGitIgnore(root)
	side = newDeferredSide()
	for _, d := range dirs {
		dir := filepath.Join(root, filepath.FromSlash(d))
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			// Absent scope dir: stays 0, and the caller's empty-instrument
			// guard turns "0" into a refusal, never into a silent clean.
			continue
		}
		err := filepath.WalkDir(dir, func(path string, ent os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if ent.IsDir() {
				if ent.Name() == "testdata" || ent.Name() == ".git" || ign.skip(path, true) {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			if ign.skip(path, false) {
				return nil
			}
			src, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			examined++
			rel, rerr := filepath.Rel(root, path)
			if rerr != nil {
				rel = path
			}
			rel = filepath.ToSlash(rel)
			for i, line := range strings.Split(string(src), "\n") {
				for _, m := range deferredMarkerRe.FindAllStringSubmatch(line, -1) {
					side.add(m[1], rel+":"+strconv.Itoa(i+1))
				}
			}
			return nil
		})
		if err != nil {
			t.Fatalf("code side is blind: walking %s for DEFERRED markers: %v", dir, err)
		}
	}
	return side, examined
}

// deferredDenominator renders the scope sentence both sides of every verdict
// is built from: what was walked, what each side holds and where. It exists so
// the number printed on green and the number asserted on red are generated by
// the same function and cannot tell two stories.
func deferredDenominator(dirs []string, examined, activeRows int, registered, marked *deferredSide) string {
	scope := make([]string, 0, len(dirs))
	for _, d := range dirs {
		scope = append(scope, filepath.ToSlash(d)+"/")
	}
	return fmt.Sprintf("DENOMINATOR scope=[%s] examined=%d production Go file(s); code markers=%d over %d id(s) [%s]; registry: activeRows=%d in SPEC-12 section 5, id registrations=%d over %d id(s) [%s]",
		strings.Join(scope, " "), examined,
		marked.total(), len(marked.count), marked.render(),
		activeRows, registered.total(), len(registered.count), registered.render())
}

// deferredBidirectionalProblems is the 1:1 assertion itself, as a pure
// function over the two ledgers: empty means every code marker has a row,
// every row has a code marker, and the per-occurrence counts agree. Keeping it
// pure is what lets TestDeferredRegistryPositiveControl feed it planted
// shapes without touching the shared tree.
func deferredBidirectionalProblems(registered, marked *deferredSide) []string {
	var problems []string
	for _, id := range deferredUnionIDs(registered, marked) {
		r, m := registered.count[id], marked.count[id]
		switch {
		case r == 0:
			problems = append(problems, fmt.Sprintf(
				"%s: DEFERRED(%s) marker(s) at [%s] have NO section-5 row (registered 0) - plant the row with all five fields, or retire the marker",
				id, id, strings.Join(marked.where[id], ", ")))
		case m == 0:
			problems = append(problems, fmt.Sprintf(
				"%s: section-5 row(s) at [%s] register DEFERRED(%s) but NO marker exists in the production Go tree - a marker was removed without retiring its row, or the row was written ahead of a marker that never landed",
				id, strings.Join(registered.where[id], ", "), id))
		case r != m:
			problems = append(problems, fmt.Sprintf(
				"%s: %d code marker occurrence(s) [%s] vs %d registered row(s) [%s] - the ruler counts OCCURRENCES (section 5's own 出现处 caliber), so each extra marker needs its own row and each extra row its own marker",
				id, m, strings.Join(marked.where[id], ", "), r, strings.Join(registered.where[id], ", ")))
		}
	}
	return problems
}

func deferredUnionIDs(a, b *deferredSide) []string {
	set := map[string]bool{}
	for id := range a.count {
		set[id] = true
	}
	for id := range b.count {
		set[id] = true
	}
	out := make([]string, 0, len(set))
	for id := range set {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}

// TestDeferredRegistryIsBidirectional is ticket 150 AC#3 on the real tree: the
// DEFERRED(Dxx-y) markers of internal/ + cmd/ and the ids registered in
// SPEC-12 section 5 must match 1:1 in BOTH directions, per occurrence.
// Measured at delivery (HEAD after 9b3cb56d7 landed the four rows): 4 markers
// against 4 rows - D28-1 x3, D11-3 x1 - so this PASSES on the current tree,
// and it has teeth because TestDeferredRegistryPositiveControl runs the same
// collectors and the same comparison over planted shapes.
func TestDeferredRegistryIsBidirectional(t *testing.T) {
	root := repoRoot(t)
	dirs := []string{"internal", "cmd"}

	registered, activeRows := deferredRowsFromSpec(t, root, deferredSpecRel)
	marked, examined := deferredMarkersInDirs(t, root, dirs)
	denominator := deferredDenominator(dirs, examined, activeRows, registered, marked)
	t.Log("d22scan: deferred-registry " + denominator)

	// Empty-instrument guards, same rule every live scope in this package
	// obeys (ticket 71 AC#4): work that did not happen outranks a verdict.
	if examined == 0 {
		t.Fatalf("code side walked 0 production Go files under %v - the walk or the scope dirs broke; an empty instrument is not a verdict (%s)", dirs, denominator)
	}
	if registered.total() == 0 && marked.total() == 0 {
		t.Fatalf("deferred-registry instrument is EMPTY on both sides (0 markers, 0 registered ids over %d active rows) - the 1:1 claim would be vacuous. If the Dxx-y family was legitimately retired, retire this test in the same commit; do not let it pass by reading nothing (%s)", activeRows, denominator)
	}

	if problems := deferredBidirectionalProblems(registered, marked); len(problems) > 0 {
		t.Fatalf("DEFERRED(Dxx-y) <-> SPEC-12 section 5 is NOT 1:1 (%d problem(s)):\n  %s\n%s",
			len(problems), strings.Join(problems, "\n  "), denominator)
	}
}

// deferredFixtureSpec renders a SPEC-12-shaped document whose section 5 holds
// exactly the given rows (plus neighbours the section boundary must respect),
// so the positive control exercises the REAL parser instead of a stub.
func deferredFixtureSpec(rows ...string) string {
	head := "## 4. 治理规程\n\n散文。\n\n## 5. 推迟项登记表（fixture）\n\n" +
		"| 类型 | 项 | 为什么现在不做（依据） | 完成判据 | 前置 | 当前残缺表现 |\n" +
		"|---|---|---|---|---|---|\n"
	body := ""
	if len(rows) > 0 {
		body = strings.Join(rows, "\n") + "\n"
	}
	return head + body + "\n## 6. 文档交付物清单\n"
}

// deferredFixtureRow is one active DEFERRED row registering id under item text,
// shaped like section 5's real 出现处 rows.
func deferredFixtureRow(id, item string) string {
	return "| DEFERRED | " + item + "（`" + id + "` 出现处①） | 依据 | 判据 | 前置 | 残缺 |"
}

// deferredFixture collects one planted tree plus its two read sides. The
// production-count guard lives here (not at the call sites) so a broken
// fixture fails as a broken fixture at the seam that broke it, never as a
// passing or failing 1:1 claim. _test.go seeds are excluded from that count
// by the same convention the walk applies - they exist to prove the
// exclusion, not to be counted.
func deferredFixture(t *testing.T, spec string, files map[string]string) (registered, marked *deferredSide) {
	t.Helper()
	root := t.TempDir()
	seedFile(t, root, deferredSpecRel, spec)
	production := 0
	for rel, body := range files {
		seedFile(t, root, rel, body)
		if !strings.HasSuffix(rel, "_test.go") {
			production++
		}
	}
	registered, _ = deferredRowsFromSpec(t, root, deferredSpecRel)
	marked, examined := deferredMarkersInDirs(t, root, []string{"internal"})
	if examined != production {
		t.Fatalf("fixture broken: walked %d production file(s), planted %d - the walk does not see what the case plants", examined, production)
	}
	return registered, marked
}

// TestDeferredRegistryPositiveControl is AC#3's mutation self-proof: the same
// collectors and the same comparison that judge the real tree must RING on the
// two shapes the ticket names, and must stay silent on the matching shape. A
// guard only ever observed passing proves nothing (the rule behind
// TestExemptScopeCannotOutliveItsAbsentTree and the self-test roster), and the
// ticket's premise is that today NOTHING rings - so every case here asserts
// the problem names the seeded id, not just that some problem exists.
func TestDeferredRegistryPositiveControl(t *testing.T) {
	t.Run("matched_fixture_is_the_control", func(t *testing.T) {
		// One marker, one row: the clean shape. Without this case the ring
		// below could be the comparison being always-on.
		registered, marked := deferredFixture(t,
			deferredFixtureSpec(deferredFixtureRow("D28-1", "压缩调用迁入 hook")),
			map[string]string{"internal/agent/fake.go": "// DEFERRED(D28-1): seeded matched marker\npackage agent\n"})
		if marked.count["D28-1"] != 1 || registered.count["D28-1"] != 1 {
			t.Fatalf("fixture broken: registered=%v marked=%v", registered.count, marked.count)
		}
		if problems := deferredBidirectionalProblems(registered, marked); len(problems) != 0 {
			t.Fatalf("the matched control must be silent, got: %v", problems)
		}
	})

	t.Run("seeded_marker_without_a_row_rings", func(t *testing.T) {
		// The ticket's literal mutation: a fake DEFERRED(D99-1) with no
		// section-5 row. Before this file landed, this was green everywhere.
		registered, marked := deferredFixture(t,
			deferredFixtureSpec(),
			map[string]string{"internal/agent/fake.go": "// DEFERRED(D99-1): seeded marker with no registered row\npackage agent\n"})
		if marked.count["D99-1"] != 1 {
			t.Fatalf("fixture broken: the walk did not harvest the seeded marker: %v", marked.count)
		}
		problems := deferredBidirectionalProblems(registered, marked)
		if len(problems) != 1 || !strings.Contains(problems[0], "D99-1") || !strings.Contains(problems[0], "NO section-5 row") {
			t.Fatalf("a DEFERRED(D99-1) marker with no row must ring naming D99-1, got: %v", problems)
		}
	})

	t.Run("seeded_row_without_a_marker_rings", func(t *testing.T) {
		// The other direction the ticket names: a registered row whose marker
		// is not in the code (deleted out from under the row, or written
		// ahead of a landing that never happened).
		registered, marked := deferredFixture(t,
			deferredFixtureSpec(deferredFixtureRow("D98-2", "表先行、码未落")),
			map[string]string{"internal/agent/fake.go": "// no marker here\npackage agent\n"})
		if registered.count["D98-2"] != 1 {
			t.Fatalf("fixture broken: the parser did not harvest the seeded row: %v", registered.count)
		}
		problems := deferredBidirectionalProblems(registered, marked)
		if len(problems) != 1 || !strings.Contains(problems[0], "D98-2") || !strings.Contains(problems[0], "NO marker exists") {
			t.Fatalf("a section-5 row for D98-2 with no marker must ring naming D98-2, got: %v", problems)
		}
	})

	t.Run("extra_occurrence_without_an_extra_row_rings", func(t *testing.T) {
		// The per-occurrence caliber's own teeth: an id that already HAS a
		// row is not a free pass for a second occurrence - section 5 registers
		// D28-1 three times because the code carries it three times.
		registered, marked := deferredFixture(t,
			deferredFixtureSpec(deferredFixtureRow("D28-1", "压缩调用迁入 hook")),
			map[string]string{
				"internal/agent/one.go": "// DEFERRED(D28-1): occurrence one, registered\npackage one\n",
				"internal/agent/two.go": "// DEFERRED(D28-1): occurrence two, NOT registered\npackage two\n",
			})
		if marked.count["D28-1"] != 2 {
			t.Fatalf("fixture broken: two planted occurrences must harvest as 2, got %v", marked.count)
		}
		problems := deferredBidirectionalProblems(registered, marked)
		if len(problems) != 1 || !strings.Contains(problems[0], "D28-1") || !strings.Contains(problems[0], "OCCURRENCES") {
			t.Fatalf("a second D28-1 marker without a second row must ring on the occurrence count, got: %v", problems)
		}
	})

	t.Run("test_file_absence_pins_stay_out_of_range", func(t *testing.T) {
		// The convention AC#3 leans on: absence-pins in *_test.go are not
		// markers (internal/tools/document_text*_test.go hold 7 DEFERRED(D34-1)
		// fixtures asserting a deleted marker STAYS deleted - counting them
		// would ring against a paid debt), and struck-through rows are
		// re-opened history, not registrations. Both shapes planted here; the
		// check must see neither, and must say so with explicit zeros instead
		// of an inferred silence.
		registered, marked := deferredFixture(t,
			deferredFixtureSpec(
				deferredFixtureRow("D28-1", "压缩调用迁入 hook"),
				"| ~~DEFERRED~~ → **重开·已落地** | 历史行（`D97-2` 出现处①） | 依据 | 判据 | 前置 | 残缺 |",
			),
			map[string]string{
				"internal/agent/one.go":        "// DEFERRED(D28-1): the registered occurrence\npackage one\n",
				"internal/agent/pin_test.go":   "package one\n\n// DEFERRED(D97-1) must be deleted; this pin asserts it never returns.\n",
				"internal/agent/other_test.go": "package one\n\nvar _ = \"DEFERRED(D97-1)\"\n",
			})
		if marked.total() != 1 || marked.count["D97-1"] != 0 {
			t.Fatalf("the walk must harvest exactly the 1 production marker and exclude the _test.go absence-pins, got %v", marked.count)
		}
		if registered.total() != 1 || registered.count["D97-2"] != 0 {
			t.Fatalf("the struck-through row must register nothing, got %v", registered.count)
		}
		if problems := deferredBidirectionalProblems(registered, marked); len(problems) != 0 {
			t.Fatalf("with the pins and the struck row excluded the fixture is 1:1, got: %v", problems)
		}
	})

	t.Run("citation_cells_do_not_register", func(t *testing.T) {
		// The row-side harvest reads the item cell only. A 依据 cell that
		// cites another row's id as an argument must not inflate the registry
		// count into a false red against a marker that exists exactly once.
		registered, marked := deferredFixture(t,
			deferredFixtureSpec(
				deferredFixtureRow("D28-1", "压缩调用迁入 hook"),
				"| DEFERRED | 迁移需同步摘除别处注释 | 与 `D28-1` 同一迁移（引用，非登记） | 判据 | 前置 | 残缺 |",
			),
			map[string]string{"internal/agent/one.go": "// DEFERRED(D28-1): single registered occurrence\npackage one\n"})
		if registered.count["D28-1"] != 1 {
			t.Fatalf("a cross-citation in the 依据 cell must not register a second occurrence, got %v", registered.count)
		}
		if problems := deferredBidirectionalProblems(registered, marked); len(problems) != 0 {
			t.Fatalf("the cross-citation fixture must be silent, got: %v", problems)
		}
	})
}
