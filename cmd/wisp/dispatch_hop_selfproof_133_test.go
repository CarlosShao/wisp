package main

// Ticket 133 AC#3 and AC#4, second instrument file. The gate in
// leg_dispatch_gate_133_test.go reads THIS directory; this file makes the gate's
// readers prove, on planted inputs, that they can answer "no".
//
// AC#3 ("a sanity leg may not be dumb") has a measured precedent in this
// repository: R-121-2, where the assembly-reachability gate's membership test was
// degraded to a substring match and the whole case stayed green, including its own
// notARealPackage sanity leg, because that leg was fed through the same degraded
// predicate. A sanity leg that cannot fail is a comment. So the legs below do not
// assert that the readers "work": they hand each reader a fixture whose answer is
// known to be bad, and assert the reader reports it. If a reader is degraded to
// always-answer-yes, the corresponding assertion here fails, which is the only
// acceptable shape of a self-proof leg.
//
// AC#4 ("silent shortening has to be audible") asks for three readings, and the
// mapping onto this instrument is written down here because the wording on the
// ticket face comes from another ruler (ticket 129's Y1/Y2/Y3):
//
//	the floor: minLegs133 is the census floor. It is exercised through
//	floorReds133, the same predicate the gate calls, and the fixture below (six
//	legs, measured) proves the predicate can answer "too few". N is the number
//	written next to minLegs133: 12, the census measured on this directory when
//	this floor was raised (ticket 133 leg r1, 2026-10-14).
//
//	the two-way roster difference: censusVsUsage133 (the gate) reconciles the
//	dispatch census against the operator-facing usage block in BOTH directions;
//	the fixture plants one command in each direction (a leg documented and not
//	dispatched, a leg dispatched and not documented) and asserts both reds.
//
//	the per-leg distinguishability: distinctnessReds133, called by the gate and
//	by this file, refuses two legs wearing the same row and refuses a table whose
//	coverage column says the same thing for every leg at once (a table that is
//	vacuous is the Y3 shape: refused, not counted). The fixture plants a pair of
//	twin legs and asserts the predicate names both.
//
// The pair half of a nail claim (test must drive the entry it is registered
// against) lives in coverageReds133: ticket 131 measured the same hole over there
// as R-131-2 (two nails whose legs were swapped, package green), and this gate's
// registry has the same two halves.
//
// WHAT THIS FILE CANNOT SEE: it proves the readers answer "no" on fixtures whose
// badness was chosen here. It says nothing about a shape nobody thought of; the
// gate's own header lists that limit, and the five AC#1 shots are where new
// shapes get added. The fixture package is written under t.TempDir(), so no
// production file is touched and the reading is repeatable on any platform.

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// the shared predicates (the gate calls both; this file exercises both)
// ---------------------------------------------------------------------------

// floorReds133 is the census floor as a predicate, so the number below is
// enforced by the gate and provable by the self-proof legs from one place.
func floorReds133(legs []*leg133) []string {
	if len(legs) >= minLegs133 {
		return nil
	}
	return []string{fmt.Sprintf("the leg census enumerated %d legs (%d is the floor written down next to minLegs133), so a dispatch rewrite is being read as \"no legs\" rather than reported. Census: %s",
		len(legs), minLegs133, strings.Join(legKeys133(legs), ", "))}
}

// coverageKind133 names the form of a coverage credit, which is the column the
// distinctness reading below compares.
func coverageKind133(covered string) string {
	switch {
	case covered == "":
		return "none"
	case strings.HasPrefix(covered, "RED"):
		return "RED"
	case strings.HasPrefix(covered, "nail"):
		return "nail"
	case strings.HasPrefix(covered, "test"):
		return "test"
	case strings.HasPrefix(covered, "ruling"):
		return "ruling"
	}
	return "other"
}

// legSignature133 renders one ledger row the way the ledger prints it. Two legs
// with the same signature are two names for one answer, which is the Y2 shape
// (the count stays right and the attribution is wrong).
func legSignature133(leg *leg133) string {
	return fmt.Sprintf("installs=%v handoff=%v covered=%s entries=%s aliases=%s",
		leg.installs, leg.handoff, leg.covered, strings.Join(leg.entries, "|"), strings.Join(leg.aliases, "|"))
}

// distinctnessReds133 is AC#4's per-leg reading. It must be called after
// coverageReds133 has filled each leg's covered column (the gate's order).
func distinctnessReds133(legs []*leg133) []string {
	var reds []string
	first := map[string]string{}
	for _, leg := range legs {
		sig := legSignature133(leg)
		if prev, ok := first[sig]; ok {
			reds = append(reds, fmt.Sprintf("legs %q and %q carry the same ledger row: %s\nTwo dispatched commands with one answer are two names for one fact, and whoever moved a leg's witness onto its twin would leave the count right and the attribution wrong - ticket 129's Y2 reading, filed against a leg table that could not tell its rows apart.",
				prev, leg.key, sig))
			continue
		}
		first[sig] = leg.key
	}
	kinds := map[string]bool{}
	installs, driven := 0, 0
	for _, leg := range legs {
		kinds[coverageKind133(leg.covered)] = true
		if leg.installs {
			installs++
		}
		if strings.HasPrefix(leg.covered, "test ") {
			driven++
		}
	}
	if len(kinds) < 2 {
		reds = append(reds, fmt.Sprintf("every leg in the census carries the same coverage reading %v, so one answer explains the whole table: a decision that is true for every row at once cannot tell a covered leg from an uncovered one, which is the Y3 shape (the judgement replaced by its expectation) arriving at this column.",
			sortedKeys133(kinds)))
	}
	if installs == 0 {
		reds = append(reds, "no leg in the census reaches "+sinkFunc133+", so the install-based half of this gate is vacuous: either the dispatch was rewritten so nothing installs (in which case the deletions this gate exists for just landed), or the walk stopped resolving calls.")
	}
	if driven == 0 {
		reds = append(reds, "no leg in the census is credited through the `covered=test` form, so the run-roster reconciliation (runRosterReds135) has an empty denominator: every reading it makes is about zero cases.")
	}
	sort.Strings(reds)
	return reds
}

// ---------------------------------------------------------------------------
// the fixture package
// ---------------------------------------------------------------------------

// writeFixture133 lays a whole miniature cmd/wisp package under a temp dir and
// returns it. loadPackage133 reads it exactly as it reads this directory.
func writeFixture133(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	var names []string
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(files[name]), 0o600); err != nil {
			t.Fatalf("write fixture %s: %v", name, err)
		}
	}
	return dir
}

const fixtureMain133 = `package main

import "os"

const usage = "wisp - fixture for ticket 133's self-proof legs\n\nUsage:\n  wisp             the no-args branch\n  wisp alpha  documented and dispatched\n  wisp beta   documented, installs the listener, no nail\n  wisp gamma  documented; its only would-be driver is not a case\n  wisp ghost  documented and never dispatched\n"

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		runResident()
		return
	}
	switch args[0] {
	case "alpha":
		os.Exit(cmdAlpha(args[1:]))
	case "beta":
		os.Exit(cmdBeta(args[1:]))
	case "gamma":
		os.Exit(cmdGamma(args[1:]))
	case "delta":
		os.Exit(cmdDelta(args[1:]))
	default:
		os.Exit(cmdUnknown())
	}
}

func cmdAlpha(args []string) int { return len(args) }

func cmdBeta(args []string) int {
	if _, err := installLogSink(""); err != nil {
		return 1
	}
	return 0
}

func cmdGamma(args []string) int { return len(args) }

func cmdDelta(args []string) int { return len(args) }

func cmdUnknown() int { return 2 }

func installLogSink(root string) (*logSink, error) { return &logSink{root: root}, nil }

type logSink struct{ root string }

func runResident() {}
`

const fixtureWrongFile133 = `package main

// WISP-LEG-COVERAGE-RULING: alpha is verified by the fixture case in this directory.

func wrongOwner133() {}
`

const fixtureTest133 = `package main

import "testing"

type phantomRecv133 struct{}

// TestPhantomMethod133 has a name starting with Test and a testing.T parameter,
// and it is still not a case: Go's testing package only runs top-level functions.
// The gate must not let it cover anything, and this fixture leg is where that is
// asserted rather than asserted-in-a-comment.
func (phantomRecv133) TestPhantomMethod133(t *testing.T) {
	cmdGamma(nil)
}

func TestFixtureAlphaDrivesTheLeg(t *testing.T) {
	cmdAlpha(nil)
}
`

const fixtureBlindMain133 = `package main

import "os"

// func main binds no local to os.Args at all: the reader has to say it cannot
// see the dispatch, not hand back an empty census that passes every row below.
func main() {
	os.Exit(0)
}
`

const fixtureTwinMain133 = `package main

import "os"

const usage = "wisp - twin fixture\n\nUsage:\n  wisp             the no-args branch\n  wisp twin1  twin one\n  wisp twin2  twin two\n"

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		runResident()
		return
	}
	switch args[0] {
	case "twin1":
		os.Exit(cmdTwin(args[1:]))
	case "twin2":
		os.Exit(cmdTwin(args[1:]))
	default:
		os.Exit(cmdUnknown())
	}
}

func cmdTwin(args []string) int { return len(args) }

func cmdUnknown() int { return 2 }

func installLogSink(root string) (*logSink, error) { return &logSink{root: root}, nil }

type logSink struct{ root string }

func runResident() {}
`

// TestAC3SelfProof133TheDispatchReadersCanAnswerNo is the self-proof leg: every
// assertion below is a reader being handed a bad input and being required to
// report it. Degrade any of the readers it feeds and this case fails.
func TestAC3SelfProof133TheDispatchReadersCanAnswerNo(t *testing.T) {
	dir := writeFixture133(t, map[string]string{
		"main.go":         fixtureMain133,
		"wrong.go":        fixtureWrongFile133,
		"fixture_test.go": fixtureTest133,
	})
	pkg, err := loadPackage133(dir)
	if err != nil {
		t.Fatalf("the fixture package did not parse: %v", err)
	}
	legs, enumReds := pkg.enumerateLegs133()
	if len(enumReds) != 0 {
		t.Errorf("the fixture's dispatch is readable by construction, and this leg exists to prove the reader reports what it can see; it reported: %v", enumReds)
	}
	keys := legKeys133(legs)
	if !containsFixture133(keys, "alpha") || !containsFixture133(keys, "delta") {
		t.Errorf("the census lost a dispatched fixture leg: %v", keys)
	}
	if containsFixture133(keys, "ghost") {
		t.Errorf("the census names %q, a command the fixture's func main never dispatches: a reader that counts what it did not read has stopped answering \"no\". Census: %v", "ghost", keys)
	}
	if fr := floorReds133(legs); len(fr) == 0 {
		t.Errorf("the fixture's %d legs are below the written floor %d, and the floor predicate said nothing", len(legs), minLegs133)
	}

	coverReds := strings.Join(pkg.coverageReds133(legs), "\n")
	for _, want := range []struct{ needle, why string }{
		{`leg "gamma"`, "gamma's only would-be driver is a method, which is not a case: the leg must stay uncovered"},
		{"covered by nothing", "the uncovered leg red is the second beat of ticket 133's fifth shot"},
		{`leg "beta"`, "a fixture leg that installs the listener and has no nail"},
		{sinkFunc133, "the install-based red has to name the symbol it walked to"},
		{"which is not a test function", "the registry's claims are compared against THIS directory's test sources"},
		{"owns neither this leg's dispatch", "a ruling parked in a file that owns nothing about its leg is prose, not a decision"},
	} {
		if !strings.Contains(coverReds, want.needle) {
			t.Errorf("the coverage reader did not report %s (looked for %q); it reported:\n%s", want.why, want.needle, coverReds)
		}
	}
	if strings.Contains(coverReds, "TestPhantomMethod133") {
		t.Errorf("a method was credited as coverage; the bucket it landed in has to hold only declarations Go's testing package runs:\n%s", coverReds)
	}

	usageReds := strings.Join(censusVsUsage133(pkg, legs), "\n")
	if !strings.Contains(usageReds, `documents command "ghost"`) {
		t.Errorf("the usage reconciliation did not report a documented command that func main never dispatches; it reported:\n%s", usageReds)
	}
	if !strings.Contains(usageReds, `leg "delta"`) || !strings.Contains(usageReds, "documented in no line of the usage block") {
		t.Errorf("the usage reconciliation did not report a dispatched leg the usage block never documents; it reported:\n%s", usageReds)
	}

	self := t.Name()
	cleanRoster := map[string]bool{"TestFixtureAlphaDrivesTheLeg": true, self: true}
	if got := pkg.runRosterReds135(legs, self, cleanRoster, len(cleanRoster), nil); len(got) != 0 {
		t.Errorf("a roster that names the credited case and this case reconciled dirty: %v", got)
	}
	missing := pkg.runRosterReds135(legs, self, map[string]bool{self: true}, 1, nil)
	if !strings.Contains(strings.Join(missing, "\n"), "has no such case") {
		t.Errorf("a roster WITHOUT the credited case has to red the row that credits it (this is the ticket 133 R-133-9 / family shot M-G reading for the `covered=test` bucket); the reader said: %v", missing)
	}
	untrusted := pkg.runRosterReds135(legs, self, map[string]bool{}, 0, nil)
	if !strings.Contains(strings.Join(untrusted, "\n"), "cannot be trusted") {
		t.Errorf("a roster that does not name the running case has to refuse to judge the claims; the reader said: %v", untrusted)
	}

	blindDir := writeFixture133(t, map[string]string{"main.go": fixtureBlindMain133})
	blindPkg, err := loadPackage133(blindDir)
	if err != nil {
		t.Fatalf("the blind fixture did not parse: %v", err)
	}
	blindLegs, blindReds := blindPkg.enumerateLegs133()
	if len(blindLegs) != 0 {
		t.Errorf("a main that binds nothing to os.Args produced %d legs instead of refusing to read the dispatch", len(blindLegs))
	}
	if !strings.Contains(strings.Join(blindReds, "\n"), "binds no local to os.Args") {
		t.Errorf("the reader has to say it cannot see the dispatch; it said: %v", blindReds)
	}

	twinDir := writeFixture133(t, map[string]string{"main.go": fixtureTwinMain133})
	twinPkg, err := loadPackage133(twinDir)
	if err != nil {
		t.Fatalf("the twin fixture did not parse: %v", err)
	}
	twinLegs, twinReds := twinPkg.enumerateLegs133()
	if len(twinReds) != 0 {
		t.Errorf("the twin fixture's dispatch is readable by construction; it reported: %v", twinReds)
	}
	twinPkg.coverageReds133(twinLegs)
	got := strings.Join(distinctnessReds133(twinLegs), "\n")
	if !strings.Contains(got, `"twin1"`) || !strings.Contains(got, `"twin2"`) {
		t.Errorf("the distinctness predicate did not name the planted twin rows (two legs dispatching the same entry, nothing else to tell them apart); it said:\n%s", got)
	}
}

// TestAC4Distinctness133NoTwoLegsWearTheSameRow is AC#4's per-leg reading on
// this directory: called here so the row-level claim has a case of its own, and
// called by the gate so a red lands on the gate's own ledger line.
func TestAC4Distinctness133NoTwoLegsWearTheSameRow(t *testing.T) {
	pkg, err := loadPackage133(".")
	if err != nil {
		t.Fatalf("AC#4 RED (the instrument, not the code): %v", err)
	}
	legs, reds := pkg.enumerateLegs133()
	if len(reds) != 0 {
		t.Errorf("the census of this directory is red, so the rows below cannot be judged: %v", reds)
	}
	if fr := floorReds133(legs); len(fr) != 0 {
		t.Errorf("AC#4 RED: %s", strings.Join(fr, "\n"))
	}
	pkg.coverageReds133(legs)
	if dr := distinctnessReds133(legs); len(dr) != 0 {
		t.Errorf("AC#4 RED: %s", strings.Join(dr, "\n"))
	}
	rows := map[string]bool{}
	kinds := map[string]bool{}
	for _, leg := range legs {
		rows[legSignature133(leg)] = true
		kinds[coverageKind133(leg.covered)] = true
	}
	t.Logf("distinctness ledger: %d legs, %d distinct rows, coverage kinds present: %v", len(legs), len(rows), sortedKeys133(kinds))
}

func containsFixture133(hay []string, needle string) bool {
	for _, s := range hay {
		if s == needle {
			return true
		}
	}
	return false
}

// processNote133 renders the disclosure the ledger owes a registry claim whose
// pair cannot be checked in-process (the marker legCover133.process carries).
func processNote133(c legCover133) string {
	if !c.process {
		return ""
	}
	return " (process claim, name-only: the case boots the real binary instead of calling the entry, so this row is a name claim about a process, not in-process evidence)"
}

// testDrivesEntry133 is the pair half of a registry claim: the named case,
// walked through this directory's test helpers, has to reach the named entry.
// coverageReds133 calls it for every claim not marked process-shaped.
func (p *pkg133) testDrivesEntry133(test, entry string) bool {
	return p.testClosure133(p.tests[test])[entry]
}
