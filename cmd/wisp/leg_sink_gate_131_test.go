package main

// Ticket 131 AC#4: the convergence instrument. Ticket 127 AC#1 put one nail on
// the resident leg, ticket 117 AC#2/AC#4 put three on the run leg, ticket 131
// AC#2/AC#3 put nails on the models and secret legs - and "one nail per leg,
// applied the eighth time" is exactly the shape this repository has been booking
// since ticket 110: whoever remembered to look is the gate. This file is where
// remembering stops being required.
//
// WHAT IT ENUMERATES, and from where. The leg list is read out of this package's
// own source at run time: go/ast over main.go's dispatch - the `switch args[0]`
// in func main, every case label, plus the `if len(args) == 0` branch that is the
// GUI leg - and a name-based transitive walk of the functions those branches
// reference, over every NON-TEST .go file in this directory. Build tags are
// ignored and same-named declarations are unioned, which over-approximates: a leg
// inherits the other platform's shape before it loses its own, and the direction
// of the approximation is the one that asks for more nails, never fewer.
//
// There is no list of leg names in this file. Ask the question AC#4 has to answer
// - "新增一条 CLI 腿时，用例会不会自动要求我给它一枚钉" - and the answer is
// mechanical: a new `case "x":` that reaches installLogSink and carries no
// registration goes red at main.go's own line number, on the next
// `go test ./cmd/wisp/`, without anybody editing this file. That reading is not
// argued from the design: it is measured, twice, in the ticket face (a planted leg
// with no nail is named; the same leg with a nail and a ruling is not).
//
// COMPLETENESS OF THE LIST ITSELF, which is what `acceptor-ticket131-r2`'s
// R-131-1 said nobody was checking, and what the three checks below are for. The
// acceptor produced three shapes where a leg with a real listener left the ledger
// in silence; each is now red, with the leg named in the reading, and each was
// re-measured on a planted leg in a throwaway snapshot (readings in
// docs/evidence/s1/131-followup-1-three-shapes.md):
//
//	形 a  X4, an early `if args[0] == "--diag"` beside the switch:
//	     strayArgvReads131 makes "argv is read nowhere else in func main" the
//	     premise of the enumeration rather than an assumption. Anything that reads
//	     the argv local, or os.Args itself, outside the binding / the len==0 leg /
//	     the switch is red at its own line, quoting the statement and the literal
//	     command it compares against. Measured zero false positives on this tree:
//	     main.go today has exactly one such binding and no stray read.
//	形 b  X8, `case sloCmdName131:` instead of `case "slo":`:
//	     caseLabels131 returns what it could not read instead of returning nothing,
//	     and an unreadable label is red, named, with the constant it resolves to.
//	     The clause also keeps its own ledger row under an "unparsed-label@..." key,
//	     so a leg cannot leave the list by being folded into somebody else's
//	     `default` - which is precisely how X8 lost the slo row, the one real
//	     occupant of the ruling marker.
//	形 c  X11/X12, `var sinkAlias = installLogSink` and a call through it:
//	     loadMainPackage131 reads package-level function values and wires them into
//	     the same call set every other edge goes through (chains included), so the
//	     leg reaches installLogSink again and owes a nail again. An alias this reader
//	     genuinely cannot place is red as blindness rather than walked past as a
//	     leaf, and the edges it did add are logged so a changed graph is never
//	     silent. A package-level closure (`var x = func() {...}`) gets an entry of
//	     its own for the same reason.
//
//	形 d  the blindness half of 形 c, measured on its own: a package-level function
//	     value this directory declares but never defines (`var factory func(string)
//	     (*logSink, error)`, filled in by an init()), plus a leg calling an alias
//	     over it. The walk can follow the alias to the name and no further, so the
//	     edge is red as "I cannot place this" with its callers named, instead of
//	     being read as a call that ends outside the package. Aliases nothing calls
//	     stay leaves, on purpose, and say so in the code.
//	第五形 Y3, ticket 129's acceptor's shape: one line inside the instrument that
//	     makes a whole column read its own expectation instead of the fact. Measured
//	     on this door, the three tautological loosening variants each produced NEW
//	     reds rather than a green, so the only silent one was a per-leg waiver - and
//	     rulingCrossChecks131 is what that costs now: a `ruled` row has to be backed
//	     by the marker sentence in the file its own rulingSites name, read by a
//	     second scan that shares no code with the doc walk, and a marker no row
//	     claims is prose rather than a decision.
//
// What remains out of reach and is NOT claimed closed: an install reached through
// a method value, through a struct field holding a function, or from a package
// other than this one is still invisible to a name walk over this directory, and
// the records predicate still only knows slog.* and observe.InitLog (ticket 131's
// own next= ①, moved to ticket 133 AC#2). The reading to compare against is the
// ledger line, not prose: a leg that is not in it at all now has to be put there
// by one of the three reds above.
//
// THE THREE STATES a leg may be in, and the one that is not allowed:
//
//	installs the listener   -> must have at least one registered nail
//	books records, no sink  -> must carry an explicit WISP-LEG-SINK-RULING
//	neither                 -> no obligation
//
// A leg that books records with no listener is R-117-1 in code form - the
// credential leg's own history before ticket 131 - so the middle row is what
// keeps a future leg from re-creating it. The ruling is a human judgment and this
// file cannot check whether a judgment is right; what it checks is that somebody
// made one, in words, next to the code, and that the words are not contradicted
// by the code (a ruling on a leg that does install a listener is red here too,
// which is what keeps the marker from turning into a way to opt out of the gate).
//
// WHAT THIS GATE DOES NOT PROVE, stated because this repository's own history
// says an unbounded claim is the thing to distrust: it cannot prove a nail is
// non-vacuous. A registered case that reads the sink and asserts nothing would
// satisfy it. Two things close the gap instead of this file alone: a registration
// takes a func(*testing.T) *value*, so the named case has to exist and compile
// into this test binary (a rename or a deletion is a compile error rather than a
// green), and the body check below requires the registered case to reach the disk
// through one of the three shared sink readers.
//
// THE ONE HOLE THAT COMPILING DOES NOT CLOSE, and what this file does about it:
// a build tag can hide an ENTIRE nail file, and then the registrations simply are
// not in the binary and nothing fails to link. That is not hypothetical - the
// four nails this gate reconciles all live in files whose names end in
// _windows_test.go, so on any other GOOS the denominator is zero. An empty
// denominator on a platform that still enumerates install legs would be the worst
// possible reading, greener than having no gate at all, so the check at the
// bottom of the case below reports it as blindness (which legs it could not
// reconcile, and which files the registrations live in) and stays red. Measured:
// on linux/amd64 this gate runs and FAILS with four unreconciled legs plus the
// named zero denominator; on windows it runs with four nails registered. What is
// NOT claimed anywhere in this file is that a green here means "every leg has a
// nail that tests it" - the ledger is about which legs reach the listener and
// which claims name them, not about what each nail asserts.
//
// R-131-2, the account that claim used to leave open: a registered nail only had
// to read A disk, so pointing the models case at the secret leg and the secret
// case at the models leg kept both rows reading `nailed` while neither leg's own
// semantics were asserted by the case claiming it (acceptor's X5, measured green
// end to end). A registration therefore now carries the production entry symbol
// its case drives, and the gate checks the claim against the dispatch closure:
// entry must exist in this package, and the claimed leg's own closure must reach
// it. What is still NOT checked, in words because the reading has to say so: a
// case that drives a real process (ticket 127's resident nail) never calls the
// entry by name, and no walk over this directory can see into that child. Those
// rows are logged as name-only claims, so `nailed` never means "asserted
// in-process" unless the detail line above it says checked.
// The per-nail "拆掉 install 会红" mutations are recorded in the ticket face, and
// they are the reading a reviewer is expected to re-run.
//
// COST, written because ticket 127's nail cost 533 lines and about 13 s: this
// case parses this directory's ~27 files with go/parser and touches no filesystem
// beyond that. Milliseconds, and no subprocess - which is the point of ticket 131
// AC#2's "能进程内驱动的就别起子进程": an enumeration gate that scaled the package's
// wall clock linearly would be a different kind of debt.

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// legSinkRulingMarker is the token that turns "nobody thought about this leg" into
// "somebody decided, here is why". Grep it: every occurrence is a leg whose
// record-keeping was ruled by hand, and the gate below is what requires one to
// exist.
const legSinkRulingMarker = "WISP-LEG-SINK-RULING:"

// legSinkReaderFuncs are the instruments a nail may read the listener through:
// this ticket's readLegSink131, ticket 117's readSink, ticket 127's
// readResidentSink. A fourth reader would be a fourth pair of eyes on one
// directory, which is how "the file exists" and "the record says what the claim
// needs" drift apart.
var legSinkReaderFuncs = map[string]bool{
	"readLegSink131":   true,
	"readSink":         true,
	"readResidentSink": true,
}

// recordEmittingSelectors are the calls that book a sentence in a logger rather
// than print it on a terminal. Qualified form only (slog.Info, not x.Info): a bare
// selector match would read every err.Error() in the package as an audit record,
// which would turn this gate into noise, and then into a rubber stamp.
var recordEmittingSelectors = map[string]bool{
	"Info": true, "Warn": true, "Error": true, "Debug": true,
}

// installsFuncName is the one production entry point that puts a persistent
// listener on a leg (logsink.go). Reaching it is the obligation enumerated here.
const installsFuncName = "installLogSink"

// subprocessEntryPrefix131 is the spelling a nail's entry claim uses to say "this
// case drives a real process, it does not call the entry in this binary". It buys
// exactly one thing - the gate stops requiring the case to call the entry - and it
// still has to name the claimed leg's own entry, so it is not a way to register a
// nail against somebody else's leg. See entryClaims131 and R-131-2.
const subprocessEntryPrefix131 = "subprocess:"

// noArgsLegKey is the label given to the branch main takes when argv is empty. It
// is the one leg with no case label to read, so the key is a name chosen here -
// and the consequence is checked rather than assumed: if that branch is ever
// removed, the nail registered against "no-args" becomes a stale claim and goes
// red (staleRegistrations131). A made-up key cannot outlive the code it was made
// for.
const noArgsLegKey = "no-args"

// legNail131 is one nail's claim about one leg: "the case whose name
// registerLegNail131 reads out of the compiled function is the reason leg X's
// listener cannot be deleted in silence", plus the symbol that case drives to make
// that true (entry, "" when the registration did not say - which is itself a red
// below, because R-131-2 is exactly the hole where a nail reads the disk of a leg
// it never dispatches).
type legNail131 struct {
	leg      string
	testName string
	entry    string
}

var legNails131 []legNail131

// registerLegNail131 is called from a package-level var in the file that owns the
// nail. Taking func(*testing.T) rather than a string is the load-bearing part: the
// claim is linked at compile time, so a nail that is renamed, moved behind a build
// tag this platform does not compile, or deleted cannot leave a registration
// pointing at nothing.
//
// The optional third argument is the name of the production entry function the case
// drives (`registerLegNail131("models", TestAC2Models..., "cmdModels")`). The gate
// checks the claim two ways and says out loud which of the two it could not check;
// see entryClaims131.
func registerLegNail131(leg string, test func(*testing.T), entry ...string) legNail131 {
	name := "<nil>"
	if test != nil {
		if fn := runtime.FuncForPC(reflect.ValueOf(test).Pointer()); fn != nil {
			name = strings.TrimPrefix(fn.Name(), "main.")
			// A closure or a generic instantiation carries a dotted suffix. A
			// top-level test case in package main does not, and if one ever did
			// the "this function exists in the sources" check in the gate is what
			// reports the mismatch, so the trimming cannot hide a bad claim.
			if i := strings.LastIndex(name, "."); i >= 0 {
				name = name[i+1:]
			}
		}
	}
	claimed := ""
	if len(entry) > 0 {
		claimed = entry[0]
	}
	n := legNail131{leg: leg, testName: name, entry: claimed}
	legNails131 = append(legNails131, n)
	return n
}

// TestAC4EveryLegIsNailedOrRuled is the gate.
func TestAC4EveryLegIsNailedOrRuled(t *testing.T) {
	pkg, err := loadMainPackage131(".")
	if err != nil {
		t.Fatalf("AC#4 RED (the instrument, not the code): %v", err)
	}
	legs, err := pkg.enumerateLegs131()
	if err != nil {
		t.Fatalf("AC#4 RED (the instrument, not the code): %v", err)
	}
	if len(legs) == 0 {
		t.Fatal("AC#4 RED: the enumeration found zero legs, so every row below would be vacuous")
	}
	sort.Slice(legs, func(i, j int) bool { return legs[i].key < legs[j].key })

	byLeg := map[string][]legNail131{}
	for _, n := range legNails131 {
		byLeg[n.leg] = append(byLeg[n.leg], n)
	}

	var lines []string
	for _, leg := range legs {
		nails := byLeg[leg.key]
		state := ""
		switch {
		case leg.installs && len(nails) == 0:
			state = "RED listener installed, no nail"
			t.Errorf("AC#4 RED: leg %q (%s) reaches %s on %d line(s) of this package and no registered nail names it.\n"+
				"Delete that install block and this package stays green - the reading ticket 127 measured for models, and ticket 117 for resident.\n"+
				"Fix: write the nail, and claim it with registerLegNail131(%q, TestYourCase) in the file that owns it.",
				leg.key, leg.site, installsFuncName, leg.installSites, leg.key)
		case leg.installs:
			state = "nailed"
			for _, n := range nails {
				fn, ok := pkg.tests[n.testName]
				if !ok {
					state = "RED nail function not found"
					t.Errorf("AC#4 RED: leg %q is registered against %q, which is not a Test function declared in this package's test sources.\n"+
						"A build tag can hide a whole file; the registration still compiled, so the sources are where this has to be caught.",
						leg.key, n.testName)
					continue
				}
				if !callsAny131(fn, legSinkReaderFuncs) {
					state = "RED nail does not read the sink"
					t.Errorf("AC#4 RED: nail %q for leg %q calls none of the shared sink readers (%s), so it cannot be the case that goes red when the install block is deleted (leg site %s).",
						n.testName, leg.key, strings.Join(sortedKeys131(legSinkReaderFuncs), ", "), leg.site)
				}
			}
		case leg.emits && !leg.ruled:
			state = "RED records with no listener and no ruling"
			t.Errorf("AC#4 RED: leg %q (%s) books records (%s) through the process logger, installs no listener, and carries no %s sentence.\n"+
				"That is R-117-1's shape: the notices exist, they die with the terminal, and nothing in the tree says anybody chose that.\n"+
				"Fix: install the sink and nail it, or write the ruling next to the code that decided it.",
				leg.key, leg.site, strings.Join(leg.emitSites, ", "), legSinkRulingMarker)
		case leg.emits && leg.ruled:
			state = "ruled"
		default:
			state = "no records"
		}
		if leg.installs && leg.ruled {
			state = "RED ruling contradicts the install"
			t.Errorf("AC#4 RED: leg %q carries a %s sentence (%v) AND reaches %s (%s).\n"+
				"A ruling is how a leg opts out of a listener, not how a missing nail is waived; delete one of the two.",
				leg.key, legSinkRulingMarker, leg.rulingSites, installsFuncName, leg.site)
		}
		lines = append(lines, fmt.Sprintf("  leg %-12s %-26s install=%-5v records=%-5v ruled=%-5v nails=%-58s -> %s",
			leg.key, leg.site, leg.installs, leg.emits, leg.ruled, nailNames131(nails), state))
	}

	for _, stale := range staleRegistrations131(legs, legNails131) {
		t.Errorf("AC#4 RED: %s", stale)
	}

	// R-131-2: whose leg does this nail actually drive.
	entryReds, entryDetail := entryClaims131(pkg, legs, legNails131)
	for _, r := range entryReds {
		t.Errorf("AC#4 RED: %s", r)
	}
	if len(entryDetail) > 0 {
		sort.Strings(entryDetail)
		t.Logf("nail entry claims, checked against the dispatch closure:\n  %s", strings.Join(entryDetail, "\n  "))
	}

	// R-131-1 形 b: a case label this reader cannot turn into a string is a leg it
	// cannot name, and a leg it cannot name is a leg that leaves the ledger in
	// silence. Both halves of that are red here: the label is named with the line
	// it sits on, and the clause keeps its own row under an unmistakable synthetic
	// key, so it cannot collide with the real default and be read as covered.
	for _, u := range pkg.unparsedLabels {
		t.Errorf("AC#4 RED: the case label at %s is %s, which is not a string literal. %s\n"+
			"This gate enumerates legs by their labels, so a label it cannot read is a leg it cannot name: the row above carries it as %q instead of dropping it or calling it default.\n"+
			"Fix: dispatch on the literal, or teach caseLabels131 the form and make it prove the reading on a planted leg.",
			u.site, u.expr, u.resolved, u.key)
	}

	// R-131-1 形 a: the ledger is only complete if argv is read nowhere else in
	// main(). An if-branch that tests argv before the switch dispatches a real leg
	// that never enters the enumeration at all - which is the same hole as "no nail
	// for a new leg", except that no nail could ever be registered for it either.
	for _, s := range pkg.strayArgvReads {
		t.Errorf("AC#4 RED: %s reads %s in func main outside every branch this gate classifies (the argv binding, the `if len(%s) == 0` leg, and the `switch %s[...]`). %s\n"+
			"Statement: %s\n"+
			"A dispatch that lives beside the switch is invisible to the leg list, so the leg it selects can install a listener with no nail and no row. Fix: move it inside the switch, or teach this gate the new shape and make it prove the reading on a planted leg.",
			s.site, s.what, pkg.argvName, pkg.argvName, s.literals, s.stmt)
	}

	// R-131-1 形 c, the half that is not an edge: a package-level function value
	// naming something this package declares but the walk cannot place is a call
	// the leg list would otherwise walk past as a leaf. Edges that WERE wired are
	// logged, because a reading that silently changes the graph is the thing this
	// whole file is about.
	for _, e := range pkg.blindEdges131 {
		t.Errorf("AC#4 RED (the instrument, not the code): %s: this reader cannot place what the alias holds, "+
			"and the functions named above call through it.\n"+
			"This gate walks names; an edge it cannot resolve is an edge a leg can reach installLogSink through without the ledger noticing. Say which it is: point the alias at a function this package declares, or teach loadMainPackage131 the form and prove the reading on a planted leg.",
			e)
	}
	if len(pkg.aliasEdges) > 0 {
		var edges []string
		for k := range pkg.aliasEdges {
			edges = append(edges, k)
		}
		sort.Strings(edges)
		t.Logf("call edges added through package-level function values (R-131-1 形 c): %s", strings.Join(edges, ", "))
	}

	// 第五形 (ticket 129's Y3 shape, aimed at this door): every `ruled` row has to
	// be backed by the marker sentence in the file it names, and every file that
	// carries the sentence has to be claimed by a row. This is the same fact read
	// by a second instrument, which is what makes loosening the first one loud.
	for _, r := range rulingCrossChecks131(pkg, legs) {
		t.Errorf("AC#4 RED: %s", r)
	}

	if len(legNails131) == 0 {
		// An empty registration list is not an empty obligation list. The legs are
		// read off the sources, which every GOOS in this repository compiles the
		// same way, while the nails are compiled into this binary and every one of
		// them now sits behind a _windows_test.go file name. On a platform whose
		// denominator is zero the rows above say "no nail" about legs that do have
		// nails nobody compiled here - and the same silence, one edit later, would
		// be a green ledger. Report the blindness by name and stay red.
		var blind []string
		for _, leg := range legs {
			if leg.installs {
				blind = append(blind, leg.key)
			}
		}
		sort.Strings(blind)
		t.Errorf("AC#4 RED (the instrument, not the code): zero nails registered in this test binary, so the gate has nothing to reconcile and would pass on an empty list.\n"+
			"GOOS=%s compiled no nail file: every registerLegNail131 call lives in a *_windows_test.go file, so the %d leg(s) in the ledger that reach %s (%s) are UNREAD on this platform, not unnailed - and no row of this ledger is coverage here, in either direction.\n"+
			"Fix: run this package where the nails compile (scripts/wisp-cli-tests.sh, the windows leg), or give this GOOS its own nail file and register it. This reading stays red on purpose: a gate that cannot see is not a gate that has seen nothing to complain about.",
			runtime.GOOS, len(blind), installsFuncName, strings.Join(blind, ", "))
	}

	report := "leg ledger, enumerated from source at run time:\n" + strings.Join(lines, "\n")
	if t.Failed() {
		t.Errorf("the ledger has a red row:\n%s", report)
	} else {
		t.Logf("%s", report)
	}
}

// leg131 is one enumerated dispatch branch of main(), with the three readings the
// gate needs and the source site each came from. reached is the whole set of
// function names the branch's transitive closure touched, which is what lets a
// registered nail's entry claim be checked against the dispatch rather than
// against prose (R-131-2).
type leg131 struct {
	key          string
	site         string
	installs     bool
	installSites int
	emits        bool
	emitSites    []string
	ruled        bool
	rulingSites  []string
	reached      map[string]bool
}

// mainPackage131 is the parsed package: production functions (callees, doc,
// site), test functions (callees), and func main's body. The four fields at the
// bottom are what the loader noticed about main() and about the package's
// function values, so the case above can report them without re-walking: a leg
// list is only trustworthy if the reader says what it could not read.
type mainPackage131 struct {
	fset     *token.FileSet
	funcs    map[string]*funcInfo131
	tests    map[string]*funcInfo131
	mainBody *ast.BlockStmt

	// argvName is the local func main binds os.Args to ("" when it binds none).
	argvName string
	// aliases maps a package-level function value to the name it holds
	// (`var sink = installLogSink`), remembering whether that name was qualified
	// (`var f = os.UserConfigDir`, which ends outside this package by definition).
	// aliasEdges records the edges those aliases added to a body's call set, keyed
	// "func:via:to", blindEdges131 the ones this reader could not place.
	aliases       map[string]aliasTarget131
	aliasEdges    map[string]bool
	declared      map[string]bool
	stringConst   map[string]string
	blindEdges131 []string

	unparsedLabels []unparsedLabel131
	strayArgvReads []argvRead131
	// rulingFiles names the files on disk whose own text carries the ruling
	// marker, collected by a plain line scan that shares no code with the AST doc
	// walk. It exists so `ruled=true` can be cross-checked against the source
	// rather than only against the reader's own predicate; see 第五形 in the
	// header and rulingCrossChecks131.
	rulingFiles map[string]bool
}

// rulingCrossChecks131 is the second reading over the same fact, which is what
// makes a one-line loosening of the first reading loud instead of silent
// (ticket 129's Y3 shape, aimed at this door: a single swapped predicate made a
// whole table read its own expectation). For every leg the ledger books as ruled,
// the file behind one of its ruling sites has to carry the marker sentence in its
// text, and every file that carries the marker has to be claimed by some ruled
// row. Either direction failing is red: the first is "the ledger says ruled
// because its own predicate says so, not because anybody wrote the ruling down",
// the second is "a ruling was written where no dispatch reaches it".
func rulingCrossChecks131(pkg *mainPackage131, legs []leg131) []string {
	var reds []string
	claimed := map[string]bool{}
	for _, leg := range legs {
		if !leg.ruled {
			continue
		}
		if len(leg.rulingSites) == 0 {
			reds = append(reds, fmt.Sprintf("leg %q is booked as ruled and names no ruling site at all: nothing on disk says anybody chose this.", leg.key))
			continue
		}
		hit := ""
		for _, site := range leg.rulingSites {
			// A site is one file:line, or several joined when the same function is
			// declared in per-platform files (cmdSLO lives in slo_windows.go and
			// slo_other.go, so its site is "slo_other.go:34 | slo_windows.go:185").
			// Each half gets checked; reading only the first one called an honest
			// ruling a false one, which this run measured on the unmutated tree.
			for _, one := range strings.Split(site, " | ") {
				file := one
				if i := strings.Index(file, ":"); i >= 0 {
					file = file[:i]
				}
				if pkg.rulingFiles[file] {
					hit = file
					claimed[file] = true
				}
			}
		}
		if hit == "" {
			reds = append(reds, fmt.Sprintf("leg %q is booked as ruled at %v, but the marker sentence %q is in none of those files: the reading came from the instrument's own predicate, not from a ruling anybody wrote down.",
				leg.key, leg.rulingSites, legSinkRulingMarker))
		}
	}
	for file := range pkg.rulingFiles {
		if !claimed[file] {
			reds = append(reds, fmt.Sprintf("%s carries a %q sentence that no ruled ledger row claims: a ruling nobody's dispatch reaches is prose, not a decision.", file, legSinkRulingMarker))
		}
	}
	sort.Strings(reds)
	return reds
}

// aliasTarget131 is the right-hand side of a package-level `var x = <name>`.
type aliasTarget131 struct {
	name      string
	qualified bool
}

// unparsedLabel131 is one case label that is not a string literal, with the
// package-level constant it resolves to when it resolves at all.
type unparsedLabel131 struct {
	site     string
	expr     string
	resolved string
	key      string
}

// argvRead131 is one read of argv in func main that sits outside every branch the
// gate classifies, carrying the statement it was found in and the string
// literals that statement mentions (the operand a dispatch comparison is made
// against, which is what names the leg the reading is about).
type argvRead131 struct {
	site     string
	what     string
	stmt     string
	literals string
}

type funcInfo131 struct {
	name    string
	site    string
	doc     string
	calls   map[string]bool
	slogSel map[string]bool
}

// loadMainPackage131 parses every .go file in dir, whatever its build constraints
// say, and that is argued from a measurement rather than assumed. The assumption
// this file used to rest on - "package main has no Linux build, because main.go's
// sherpa import makes `GOOS=linux go vet ./cmd/wisp/` rc=1" - was only ever true
// of cross compiling FROM windows, where the sherpa constraint error comes first
// and hides everything behind it. On a real linux runner the same command
// compiles this package, test files included (CGO_ENABLED=1, rc=0), so
// ignoring build tags here is not covered by "nothing else compiles anyway":
// it unions same-named declarations across platforms on purpose, and the
// denominator that CAN go empty on another platform is the registration list,
// which the blindness check at the bottom of the gate reports instead of
// swallowing.
func loadMainPackage131(dir string) (*mainPackage131, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", dir, err)
	}
	pkg := &mainPackage131{
		fset:        token.NewFileSet(),
		funcs:       map[string]*funcInfo131{},
		tests:       map[string]*funcInfo131{},
		aliases:     map[string]aliasTarget131{},
		aliasEdges:  map[string]bool{},
		declared:    map[string]bool{},
		stringConst: map[string]string{},
		rulingFiles: map[string]bool{},
	}
	var parsed int
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") {
			continue
		}
		f, err := parser.ParseFile(pkg.fset, filepath.Join(dir, name), nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", name, err)
		}
		parsed++
		isTest := strings.HasSuffix(name, "_test.go")
		// 第五形's second instrument: does this file's own text carry the ruling
		// marker? Read off the bytes, not off the AST doc the ledger's ruled column
		// is built from, so the two readings can disagree and say so. Test files are
		// out of the scan: this very file states the marker in a const, which is an
		// instrument's string, not a leg's ruling.
		if !isTest {
			if src, srcErr := os.ReadFile(filepath.Join(dir, name)); srcErr == nil && bytes.Contains(src, []byte(legSinkRulingMarker)) {
				pkg.rulingFiles[name] = true
			}
		}
		for _, d := range f.Decls {
			fn, ok := d.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			info := &funcInfo131{
				name:    fn.Name.Name,
				site:    fmt.Sprintf("%s:%d", name, pkg.fset.Position(fn.Pos()).Line),
				doc:     fn.Doc.Text(),
				calls:   map[string]bool{},
				slogSel: map[string]bool{},
			}
			walkCallees131(fn.Body, info)
			if fn.Recv == nil {
				pkg.declared[fn.Name.Name] = true
			}
			if !isTest && fn.Recv == nil && fn.Name.Name == "main" {
				pkg.mainBody = fn.Body
			}
			into := pkg.funcs
			if isTest && fn.Recv == nil && strings.HasPrefix(fn.Name.Name, "Test") {
				into = pkg.tests
			} else if isTest {
				// Test helpers are not production code: they stay out of the
				// closure a leg walks, because a helper that calls
				// installLogSink directly (several do) would otherwise hand its
				// own fixture to a leg that never runs it.
				continue
			}
			if prev, ok := into[fn.Name.Name]; ok {
				prev.merge(info)
			} else {
				into[fn.Name.Name] = info
			}
		}
		// R-131-1 形 c, first half: the package's named values. A production file
		// can take the shape `var sinkAlias = installLogSink` and call the alias,
		// and a name-based walk that only knows funcDecl names reads that edge as a
		// leaf - which is how one line of refactoring (换实现、按平台分派、注入替身)
		// put a live leg outside the ledger. Collect the alias and the string
		// constants (the labels caseLabels131 cannot read are usually one of those);
		// the edges themselves are added once every declaration in every file is
		// known, below.
		for _, d := range f.Decls {
			gd, ok := d.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range gd.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					pkg.declared[ts.Name.Name] = true
					continue
				}
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				// Every package-level name is declared here, whether or not this
				// file says what holds it. `var factory func(string) (*logSink, error)`
				// with no initializer is a name a leg can call through - the shape
				// a platform-split or an injected fake grows into - and a reader
				// that recorded only initialised names would resolve an alias over
				// it to "not this package's business" and walk past it as a leaf.
				for _, id := range vs.Names {
					pkg.declared[id.Name] = true
				}
				if len(vs.Values) == 0 || len(vs.Names) != len(vs.Values) {
					continue
				}
				for i, id := range vs.Names {
					if gd.Tok == token.CONST {
						if bl, ok := vs.Values[i].(*ast.BasicLit); ok && bl.Kind == token.STRING {
							if un, err := strconv.Unquote(bl.Value); err == nil {
								pkg.stringConst[id.Name] = un
							}
						}
						continue
					}
					if gd.Tok != token.VAR || id.Name == "_" {
						continue
					}
					switch v := vs.Values[i].(type) {
					case *ast.FuncLit:
						// `var x = func(...) { ... }` at package level: the body is
						// a real function of this package that no FuncDecl
						// declaration ever names, so without an entry here every
						// call through it - including an installLogSink - is
						// invisible. Registered under the var's own name.
						if _, taken := pkg.funcs[id.Name]; !taken {
							fi := &funcInfo131{
								name:    id.Name,
								site:    fmt.Sprintf("%s:%d", name, pkg.fset.Position(v.Pos()).Line),
								calls:   map[string]bool{},
								slogSel: map[string]bool{},
							}
							walkCallees131(v.Body, fi)
							pkg.funcs[id.Name] = fi
						}
					default:
						if target, qualified := calleeName131(vs.Values[i]); target != "" {
							pkg.aliases[id.Name] = aliasTarget131{name: target, qualified: qualified}
						}
					}
				}
			}
		}
	}
	// R-131-1 形 c, second half: an alias is a call edge, so add it wherever it is
	// called. Both readings the gate takes from a body (does this leg reach
	// installLogSink, does this nail read the sink off disk) are computed from this
	// call set, so the edge has to be in the set before either is asked.
	//
	// What stays a leaf on purpose: an alias over something this directory does not
	// declare at all (`var now = time.Now`, `var dial = net.Dial`). Those calls end
	// outside the package the gate walks, which is the same fact a direct call to
	// them states, and adding a red for every stdlib function this command touches
	// would be noise that gets silenced. What is NOT a leaf: a name this package
	// declares that the walk cannot place - that one goes on blindEdges131 and the
	// gate says it cannot see the edge.
	for alias, at := range pkg.aliases {
		name, kind := pkg.resolveAlias131(at)
		switch kind {
		case aliasExternal131:
			continue
		case aliasBlind131:
			callers := pkg.callersOf131(alias)
			if len(callers) == 0 {
				// Nothing in this package calls it, so no leg can reach the sink
				// through it either: "leaf" is the honest reading, and a red for
				// every package-level value this reader cannot place would be the
				// kind of noise that gets silenced.
				continue
			}
			pkg.blindEdges131 = append(pkg.blindEdges131, fmt.Sprintf(
				"var %s = %s, called by %s", alias, at.name, strings.Join(callers, ", ")))
			continue
		}
		for _, fi := range pkg.funcs {
			if fi.calls[alias] && !fi.calls[name] {
				fi.calls[name] = true
				pkg.aliasEdges[fmt.Sprintf("%s:%s->%s", fi.name, alias, name)] = true
			}
		}
		for _, fi := range pkg.tests {
			if fi.calls[alias] && !fi.calls[name] {
				fi.calls[name] = true
				pkg.aliasEdges[fmt.Sprintf("%s:%s->%s", fi.name, alias, name)] = true
			}
		}
	}
	sort.Strings(pkg.blindEdges131)
	if parsed == 0 {
		return nil, fmt.Errorf("no .go files in %s: empty instrument, and an empty instrument is red", dir)
	}
	if pkg.mainBody == nil {
		return nil, fmt.Errorf("no func main found in %s: the dispatch this gate enumerates is not there to be read", dir)
	}
	return pkg, nil
}

// walkCallees131 records both readings the gate takes from a body: every called
// name (functions and methods, by bare name) and every qualified record-booking
// call (slog.Info and friends).
func walkCallees131(body ast.Node, fi *funcInfo131) {
	ast.Inspect(body, func(n ast.Node) bool {
		ce, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		switch fn := ce.Fun.(type) {
		case *ast.Ident:
			fi.calls[fn.Name] = true
		case *ast.SelectorExpr:
			fi.calls[fn.Sel.Name] = true
			if x, ok := fn.X.(*ast.Ident); ok {
				switch {
				case x.Name == "slog" && recordEmittingSelectors[fn.Sel.Name]:
					fi.slogSel[fn.Sel.Name] = true
				case x.Name == "observe" && fn.Sel.Name == "InitLog":
					// A leg that resolves the pipeline by hand has a listener
					// question to answer too, and `wisp slo` is the one leg that
					// does it without installLogSink. It is reported through the
					// same reading as a slog call because the obligation is the
					// same one: say what keeps the records, or say why nothing
					// does.
					fi.slogSel["InitLog"] = true
				}
			}
		}
		return true
	})
}

// merge folds a second declaration of the same name into the first. package main
// has runResident in resident_windows.go and a fail-closed counterpart in
// resident_other.go, so the walk unions them (see loadMainPackage131's note on
// build constraints): a leg inherits every same-named shape in the tree before it
// inherits none of them, which is the direction that asks for more nails.
func (fi *funcInfo131) merge(other *funcInfo131) {
	for c := range other.calls {
		fi.calls[c] = true
	}
	for s := range other.slogSel {
		fi.slogSel[s] = true
	}
	if fi.doc == "" {
		fi.doc = other.doc
	} else if other.doc != "" {
		// Both platform halves of a leg may carry prose (cmdSLO does), and a
		// ruling marker has to be findable wherever it was written: keeping only
		// the first declaration's comment would make "where did you put the
		// ruling?" a question about file-name ordering.
		fi.doc = fi.doc + "\n" + other.doc
	}
	fi.site = fi.site + " | " + other.site
}

// enumerateLegs131 reads main()'s dispatch: the argv-empty branch and every case
// of the switch over args[0]. A leg's obligation is decided by what its own
// branch can reach, transitively, through this package's production files.
func (pkg *mainPackage131) enumerateLegs131() ([]leg131, error) {
	var sw *ast.SwitchStmt
	var empty *ast.BlockStmt
	var emptyIf *ast.IfStmt
	var emptyPos token.Position
	argvName := ""
	ast.Inspect(pkg.mainBody, func(n ast.Node) bool {
		switch s := n.(type) {
		case *ast.AssignStmt:
			// The local main() binds os.Args to. Reading it instead of
			// hard-coding "args" is what lets this keep working when the local is
			// renamed, and stop working when the switch stops being a dispatch
			// over argv at all - which is the shape an empty ledger would
			// otherwise disguise as "no legs, so nothing to nail".
			if len(s.Lhs) == 1 && len(s.Rhs) == 1 {
				if id, ok := s.Lhs[0].(*ast.Ident); ok && isOsArgs131(s.Rhs[0]) {
					argvName = id.Name
				}
			}
		case *ast.SwitchStmt:
			// The tag this package dispatches on is args[0]: an IndexExpr over
			// that local. Matching the identifier rather than the whole
			// expression follows a re-index; requiring it to be the argv local is
			// what keeps some other switch in main from being read as the
			// dispatch.
			if argvName != "" && sw == nil {
				if id := tagIdentOf131(s.Tag); id != nil && id.Name == argvName {
					sw = s
				}
			}
		case *ast.IfStmt:
			if s.Body == nil || empty != nil {
				return true
			}
			if isLenZeroOf131(s.Cond, argvName) {
				empty = s.Body
				emptyIf = s
				emptyPos = pkg.fset.Position(s.Pos())
			}
		}
		return true
	})
	if argvName == "" {
		return nil, fmt.Errorf("func main never binds a local to os.Args, so this gate cannot tell a dispatch over argv from some other switch")
	}
	if sw == nil {
		return nil, fmt.Errorf("func main has no `switch %s[...]`: either the dispatch moved, in which case this reader has to be pointed at where it went, or it is gone. Neither is a green", argvName)
	}
	pkg.argvName = argvName

	// R-131-1 形 a: the ledger is complete only if nothing else in main() reads
	// argv. Recorded here, next to the two nodes that define what "classified"
	// means, so the definition cannot drift from the dispatch it describes.
	pkg.strayArgvReads = pkg.strayArgvReads131(argvName, sw, emptyIf)

	var legs []leg131
	for _, cl := range sw.Body.List {
		cc, ok := cl.(*ast.CaseClause)
		if !ok {
			continue
		}
		keys, unparsed := caseLabels131(cc, pkg)
		site := fmt.Sprintf("main.go:%d", pkg.fset.Position(cc.Pos()).Line)
		for _, u := range unparsed {
			legs = append(legs, pkg.legFrom131(u.key, u.site, cc.Body))
		}
		if len(keys) == 0 && len(unparsed) == 0 {
			keys = []string{"default"}
		}
		for _, k := range keys {
			legs = append(legs, pkg.legFrom131(k, site, cc.Body))
		}
	}
	if empty == nil {
		// Not fatal on its own: the nail registered against "no-args" turns the
		// absence into a red stale claim, which is the reading with a name on it.
	} else {
		legs = append(legs, pkg.legFrom131(noArgsLegKey, fmt.Sprintf("main.go:%d", emptyPos.Line), empty.List))
	}
	return legs, nil
}

// isOsArgs131 recognises os.Args and os.Args[...], the two spellings a main() may
// bind its argv local with.
func isOsArgs131(e ast.Expr) bool {
	switch x := e.(type) {
	case *ast.IndexExpr:
		return isOsArgs131(x.X)
	case *ast.SliceExpr:
		return isOsArgs131(x.X)
	case *ast.SelectorExpr:
		id, ok := x.X.(*ast.Ident)
		return ok && id.Name == "os" && x.Sel.Name == "Args"
	}
	return false
}

// tagIdentOf131 is the identifier a switch dispatches on, through the index
// (`args[0]`) this package's dispatch uses.
func tagIdentOf131(e ast.Expr) *ast.Ident {
	switch x := e.(type) {
	case *ast.IndexExpr:
		if id, ok := x.X.(*ast.Ident); ok {
			return id
		}
	case *ast.Ident:
		return x
	}
	return nil
}

// legFrom131 turns one branch's statements into a leg reading by walking the
// package's call graph from every function the branch names.
func (pkg *mainPackage131) legFrom131(key, site string, stmts []ast.Stmt) leg131 {
	leg := leg131{key: key, site: site, reached: map[string]bool{}}
	seen := map[string]bool{}
	queue := referencedFuncs131(stmts)
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		if seen[name] {
			continue
		}
		seen[name] = true
		fi, ok := pkg.funcs[name]
		if !ok {
			continue
		}
		leg.reached[name] = true
		if fi.calls[installsFuncName] {
			leg.installs = true
			leg.installSites++
		}
		if len(fi.slogSel) > 0 {
			leg.emits = true
			for s := range fi.slogSel {
				leg.emitSites = append(leg.emitSites, fmt.Sprintf("slog.%s@%s", s, fi.site))
			}
		}
		if strings.Contains(fi.doc, legSinkRulingMarker) {
			leg.ruled = true
			leg.rulingSites = append(leg.rulingSites, fi.site)
		}
		for c := range fi.calls {
			if !seen[c] {
				queue = append(queue, c)
			}
		}
	}
	sort.Strings(leg.emitSites)
	sort.Strings(leg.rulingSites)
	return leg
}

// staleRegistrations131 catches the other direction: a nail claimed for a leg that
// no longer exists, or that stopped installing a listener. Both are "the
// instrument is green because it is comparing against a leg nobody dispatches to".
func staleRegistrations131(legs []leg131, nails []legNail131) []string {
	byKey := map[string]leg131{}
	for _, l := range legs {
		byKey[l.key] = l
	}
	var out []string
	for _, n := range nails {
		leg, ok := byKey[n.leg]
		if !ok {
			out = append(out, fmt.Sprintf("nail %q claims leg %q, which main.go does not dispatch to any more", n.testName, n.leg))
			continue
		}
		if !leg.installs {
			out = append(out, fmt.Sprintf("nail %q claims leg %q, which no longer installs the persistent sink (%s): the nail is now proving something else, or nothing",
				n.testName, n.leg, leg.site))
		}
	}
	return out
}

// entryClaims131 is R-131-2: the ledger's nails= column used to answer only "is
// there a case that reads a disk, registered under this leg's name", so swapping
// the two legs a pair of nails claim (acceptor's X5) kept both rows reading
// `nailed` while each case drove the other leg's code. A registration now carries
// the production entry symbol its case drives, and this function checks that claim
// against the dispatch rather than against prose:
//
//	entry must name a function this package's production files declare;
//	the claimed leg's own transitive closure must reach that entry;
//	the claimed case must CALL that entry - unless the claim is spelled
//	"subprocess:<name>", which says out loud that the drive happens in a real
//	process no name walk can follow.
//
// The third rule is what keeps this check from being satisfied by a swapped pair
// that merely names a reachable symbol: registering the secret case against the
// models leg with the models entry claim is red, because that case does not call
// cmdModels. The subprocess form does not opt out of the first two rules either -
// it still has to name this leg's own entry - it only declines the third, and the
// gate logs it as a name-only claim so `nailed` cannot be read as "somebody asserts
// this leg's semantics in-process" (ticket 127's resident nail is that shape, and
// R-131-2's fallback sentence asked for that to be said rather than implied).
func entryClaims131(pkg *mainPackage131, legs []leg131, nails []legNail131) (reds []string, detail []string) {
	byKey := map[string]leg131{}
	for _, l := range legs {
		byKey[l.key] = l
	}
	for _, n := range nails {
		if n.entry == "" {
			reds = append(reds, fmt.Sprintf("nail %q claims leg %q with no entry function, so nothing here can check that the case drives that leg at all (R-131-2 is the hole where it drives a different one). Fix: pass the production entry symbol as registerLegNail131's third argument, spelled \"subprocess:<name>\" when the case drives a real process.",
				n.testName, n.leg))
			continue
		}
		entry, outOfProcess := strings.CutPrefix(n.entry, subprocessEntryPrefix131)
		if !outOfProcess {
			entry = n.entry
		}
		if _, isFunc := pkg.funcs[entry]; !isFunc {
			reds = append(reds, fmt.Sprintf("nail %q claims leg %q and entry %q, which is not a function declared in this package's production files. Either the entry was renamed or it never was this leg's entry.",
				n.testName, n.leg, entry))
			continue
		}
		leg, ok := byKey[n.leg]
		if !ok {
			// staleRegistrations131 already names a leg that is gone; no second red
			// for the same fact.
			continue
		}
		if !leg.reached[entry] {
			reds = append(reds, fmt.Sprintf("nail %q claims leg %q and entry %q, but the dispatch main.go makes for %q never reaches %q (%s). The case drives some other leg's code: the ledger's nailed row for %q is a claim about the wrong function.",
				n.testName, n.leg, entry, n.leg, entry, leg.site, n.leg))
			continue
		}
		fn, known := pkg.tests[n.testName]
		if !known {
			// The nailed-row loop above already reds a case that is not in these
			// sources; this reading does not name the same fact twice.
			continue
		}
		if fn.calls[entry] {
			detail = append(detail, fmt.Sprintf("nail %q -> leg %q, entry %q: checked, the case calls the entry this leg's dispatch reaches",
				n.testName, n.leg, entry))
			continue
		}
		if outOfProcess {
			detail = append(detail, fmt.Sprintf("nail %q -> leg %q, entry %q: declared with the %q prefix, so this row is a name-only claim about a real process and NOT evidence that anybody asserts %q's semantics in-process.",
				n.testName, n.leg, entry, subprocessEntryPrefix131, n.leg))
			continue
		}
		reds = append(reds, fmt.Sprintf("nail %q claims leg %q through entry %q, and this leg's dispatch does reach it, but the registered case never calls %q. Either the case drives a different leg (register it against that one), or it reads a disk this leg never writes; if it drives a real process instead of calling this function, say so by spelling the claim %s%q.",
			n.testName, n.leg, entry, entry, subprocessEntryPrefix131, entry))
	}
	return reds, detail
}

func nailNames131(nails []legNail131) string {
	if len(nails) == 0 {
		return "-"
	}
	names := make([]string, 0, len(nails))
	for _, n := range nails {
		names = append(names, n.testName)
	}
	sort.Strings(names)
	return strings.Join(names, ",")
}

func callsAny131(fn *funcInfo131, set map[string]bool) bool {
	for c := range fn.calls {
		if set[c] {
			return true
		}
	}
	return false
}

func sortedKeys131(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// referencedFuncs131 is a branch's seed set: every name the branch's own
// statements mention, including nested call arguments - os.Exit(cmdModels(...))
// has to seed cmdModels, which is the whole point of reading the dispatch rather
// than a list of names.
func referencedFuncs131(stmts []ast.Stmt) []string {
	var out []string
	for _, s := range stmts {
		ast.Inspect(s, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.Ident:
				out = append(out, x.Name)
			case *ast.SelectorExpr:
				out = append(out, x.Sel.Name)
			}
			return true
		})
	}
	return out
}

// caseLabels131 reads a case clause's string literals: `case "version",
// "--version", "-v":` yields three legs, because those are three spellings an
// operator can type and each has to be answerable.
//
// A label it is NOT a string literal used to fall through as "this clause has no
// labels", which enumerateLegs131 then read as the default row - so naming a
// label (`case sloCmdName131:`) silently moved that leg out of the ledger and onto
// somebody else's key. That is R-131-1's 形 b, and the second return value is what
// it costs now: every unreadable label comes back named, keeps its own row under
// an unmistakable key, and is reported red by the gate.
func caseLabels131(cl *ast.CaseClause, pkg *mainPackage131) ([]string, []unparsedLabel131) {
	var out []string
	var unparsed []unparsedLabel131
	for _, e := range cl.List {
		if bl, ok := e.(*ast.BasicLit); ok && bl.Kind == token.STRING {
			if un, err := strconv.Unquote(bl.Value); err == nil {
				out = append(out, un)
			}
			continue
		}
		site := fmt.Sprintf("main.go:%d", pkg.fset.Position(e.Pos()).Line)
		expr := exprText131(pkg, e)
		u := unparsedLabel131{
			site: site,
			expr: expr,
			key:  fmt.Sprintf("unparsed-label@%s:%s", site, expr),
		}
		if id, ok := e.(*ast.Ident); ok {
			if v, known := pkg.stringConst[id.Name]; known {
				u.resolved = fmt.Sprintf("It is package-level const %s = %q, so the leg this clause dispatches is %q - a name this reader would only know by walking somebody else's declaration, which is the assumption that just got a leg dropped.",
					id.Name, v, v)
			} else {
				u.resolved = fmt.Sprintf("It is the identifier %s, which is not a package-level string constant, so this reader has no way to name the leg at all.", id.Name)
			}
		} else {
			u.resolved = "It is neither a string literal nor a bare identifier, so there is nothing here to resolve a command name out of."
		}
		unparsed = append(unparsed, u)
	}
	if len(unparsed) > 0 {
		pkg.unparsedLabels = append(pkg.unparsedLabels, unparsed...)
	}
	return out, unparsed
}

// strayArgvReads131 is R-131-1's 形 a check: every read of argv inside func main
// has to belong to a branch the gate enumerates. The classified spans are exactly
// the three nodes the enumeration is built from - the binding of os.Args to the
// local, the `if len(args) == 0` leg (condition and body), and the dispatch switch
// - and anything else that touches argv dispatches a leg this ledger will never
// contain. It reports the statement it found the read in, and the string
// literals that statement compares against, because the point of the reading is to
// name the leg, not the line.
func (pkg *mainPackage131) strayArgvReads131(argvName string, sw *ast.SwitchStmt, emptyIf *ast.IfStmt) []argvRead131 {
	type span struct{ lo, hi token.Pos }
	var spans []span
	if sw != nil {
		spans = append(spans, span{sw.Pos(), sw.End()})
	}
	if emptyIf != nil {
		spans = append(spans, span{emptyIf.Pos(), emptyIf.End()})
	}
	ast.Inspect(pkg.mainBody, func(n ast.Node) bool {
		if s, ok := n.(*ast.AssignStmt); ok && len(s.Lhs) == 1 && len(s.Rhs) == 1 {
			if id, ok := s.Lhs[0].(*ast.Ident); ok && id.Name == argvName && isOsArgs131(s.Rhs[0]) {
				spans = append(spans, span{s.Pos(), s.End()})
			}
		}
		return true
	})
	classified := func(pos token.Pos) bool {
		for _, s := range spans {
			if s.lo <= pos && pos < s.hi {
				return true
			}
		}
		return false
	}

	var out []argvRead131
	seenLine := map[int]bool{}
	seenStmt := map[int]bool{}
	ast.Inspect(pkg.mainBody, func(n ast.Node) bool {
		what := ""
		switch x := n.(type) {
		case *ast.Ident:
			if x.Name == argvName {
				what = fmt.Sprintf("%s (the local os.Args was bound to)", argvName)
			}
		case *ast.SelectorExpr:
			if isOsArgs131(x) {
				what = "os.Args"
			}
		}
		if what == "" {
			return true
		}
		if classified(n.Pos()) {
			return true
		}
		line := pkg.fset.Position(n.Pos()).Line
		stmt := enclosingStmt131(pkg, n.Pos())
		// One reading per statement, not per token: the dedup key is the line the
		// enclosing statement starts on, so an `if` that touches argv three times
		// in its own condition costs one red, not three. Measured on the planted
		// 形 a leg (`if len(args) > 0 && args[0] == "--diag" { os.Exit(...) }`):
		// it gave two reds, one for that statement - naming the literal "--diag" -
		// and one for the nested `os.Exit(cmdDiag131(args[1:]))`, which is a
		// different statement reading argv again. Two reds for one leg missing
		// from the ledger is not the duplication this dedup exists to prevent, and
		// each of them points at the line it is about.
		if stmt != nil {
			stmtLine := pkg.fset.Position(stmt.Pos()).Line
			if seenStmt[stmtLine] {
				return true
			}
			seenStmt[stmtLine] = true
		} else if seenLine[line] {
			return true
		} else {
			seenLine[line] = true
		}
		read := argvRead131{
			site:     fmt.Sprintf("main.go:%d", line),
			what:     what,
			stmt:     "<no statement found>",
			literals: "",
		}
		if stmt != nil {
			read.stmt = collapse131(exprText131(pkg, stmt))
			read.literals = literalsIn131(stmt)
		}
		out = append(out, read)
		return true
	})
	return out
}

// enclosingStmt131 is the innermost statement of func main holding pos, which is
// the smallest piece of the dispatch an unread argv read can be described by.
func enclosingStmt131(pkg *mainPackage131, pos token.Pos) ast.Stmt {
	var best ast.Stmt
	ast.Inspect(pkg.mainBody, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		s, ok := n.(ast.Stmt)
		if !ok || s.Pos() > pos || pos >= s.End() {
			return true
		}
		if best == nil || (s.Pos() >= best.Pos() && s.End() <= best.End()) {
			best = s
		}
		return true
	})
	return best
}

// literalsIn131 names the command spellings a statement compares argv against, so
// a reading about a leg says which leg it means ("--diag", not "line 47").
func literalsIn131(n ast.Node) string {
	var found []string
	seen := map[string]bool{}
	ast.Inspect(n, func(x ast.Node) bool {
		bl, ok := x.(*ast.BasicLit)
		if !ok || bl.Kind != token.STRING {
			return true
		}
		un, err := strconv.Unquote(bl.Value)
		if err != nil || seen[un] {
			return true
		}
		seen[un] = true
		found = append(found, strconv.Quote(un))
		return true
	})
	if len(found) == 0 {
		return "The statement names no literal command, so this reader cannot say which leg it dispatches."
	}
	return "Literals in that statement: " + strings.Join(found, ", ")
}

// exprText131 renders a node back to source for a reading, on one line.
func exprText131(pkg *mainPackage131, n ast.Node) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, pkg.fset, n); err != nil {
		return fmt.Sprintf("<unprintable node at %s>", pkg.fset.Position(n.Pos()))
	}
	return strings.TrimSpace(buf.String())
}

// collapse131 flattens a multi-line statement into one line for the reading.
func collapse131(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// calleeName131 is the function name a package-level value declaration points at:
// `var x = installLogSink` names it bare, `var x = os.UserConfigDir` names
// something outside this package (the bool says which), and anything else
// (a closure, a call, a composite literal) is not an alias this reader can wire.
func calleeName131(e ast.Expr) (string, bool) {
	switch x := e.(type) {
	case *ast.Ident:
		return x.Name, false
	case *ast.SelectorExpr:
		return x.Sel.Name, true
	}
	return "", false
}

// callersOf131 names the functions and cases of this package that call a name. It
// is what turns "I cannot place this alias" from a guess into a reading: an alias
// nothing calls cannot carry a leg anywhere, and one something calls can.
func (pkg *mainPackage131) callersOf131(name string) []string {
	var out []string
	for _, fi := range pkg.funcs {
		if fi.calls[name] {
			out = append(out, fi.name)
		}
	}
	for _, fi := range pkg.tests {
		if fi.calls[name] {
			out = append(out, fi.name)
		}
	}
	sort.Strings(out)
	return out
}

// builtin131 is only what this package could plausibly alias at package level; an
// alias over a builtin is a leaf, not an edge to a function of this package.
var builtin131 = map[string]bool{
	"append": true, "cap": true, "close": true, "copy": true, "delete": true,
	"len": true, "make": true, "new": true, "panic": true, "print": true,
	"println": true, "recover": true, "max": true, "min": true, "clear": true,
}

const (
	aliasEdge131 = iota
	aliasExternal131
	aliasBlind131
)

// resolveAlias131 follows a package-level function value to what it actually
// holds, through chains of the same shape, and says which of three things the walk
// found: a function of this package (wire the edge), a call that ends outside it
// (a leaf, exactly what calling it directly would be), or a name this package
// declares that the walk cannot place (blind, and the gate says so).
func (pkg *mainPackage131) resolveAlias131(at aliasTarget131) (string, int) {
	if at.qualified {
		return "", aliasExternal131
	}
	name := at.name
	for depth := 0; depth < 8; depth++ {
		next, isAlias := pkg.aliases[name]
		if !isAlias {
			break
		}
		if next.qualified {
			return "", aliasExternal131
		}
		name = next.name
	}
	if _, ok := pkg.funcs[name]; ok {
		return name, aliasEdge131
	}
	if _, ok := pkg.tests[name]; ok {
		return name, aliasEdge131
	}
	if builtin131[name] {
		return "", aliasExternal131
	}
	if !pkg.declared[name] {
		// Not this package's name at all: nothing to wire, nothing to confess.
		return "", aliasExternal131
	}
	return name, aliasBlind131
}

// isLenZeroOf131 recognises the argv-empty guard (`if len(args) == 0`) so the GUI
// leg is enumerated from the code rather than from a name invented here. The
// identifier it compares against is the local os.Args was bound to, read above.
func isLenZeroOf131(e ast.Expr, argvName string) bool {
	if argvName == "" {
		return false
	}
	be, ok := e.(*ast.BinaryExpr)
	if !ok || be.Op != token.EQL {
		return false
	}
	call, ok := be.X.(*ast.CallExpr)
	if !ok {
		return false
	}
	fn, ok := call.Fun.(*ast.Ident)
	if !ok || fn.Name != "len" || len(call.Args) != 1 {
		return false
	}
	id, ok := call.Args[0].(*ast.Ident)
	if !ok || id.Name != argvName {
		return false
	}
	lit, ok := be.Y.(*ast.BasicLit)
	return ok && lit.Kind == token.INT && lit.Value == "0"
}
