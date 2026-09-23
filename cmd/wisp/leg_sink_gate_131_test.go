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
// The per-nail "拆掉 install 会红" mutations are recorded in the ticket face, and
// they are the reading a reviewer is expected to re-run.
//
// COST, written because ticket 127's nail cost 533 lines and about 13 s: this
// case parses this directory's ~27 files with go/parser and touches no filesystem
// beyond that. Milliseconds, and no subprocess - which is the point of ticket 131
// AC#2's "能进程内驱动的就别起子进程": an enumeration gate that scaled the package's
// wall clock linearly would be a different kind of debt.

import (
	"fmt"
	"go/ast"
	"go/parser"
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

// noArgsLegKey is the label given to the branch main takes when argv is empty. It
// is the one leg with no case label to read, so the key is a name chosen here -
// and the consequence is checked rather than assumed: if that branch is ever
// removed, the nail registered against "no-args" becomes a stale claim and goes
// red (staleRegistrations131). A made-up key cannot outlive the code it was made
// for.
const noArgsLegKey = "no-args"

// legNail131 is one nail's claim about one leg: "the case whose name
// registerLegNail131 reads out of the compiled function is the reason leg X's
// listener cannot be deleted in silence".
type legNail131 struct {
	leg      string
	testName string
}

var legNails131 []legNail131

// registerLegNail131 is called from a package-level var in the file that owns the
// nail. Taking func(*testing.T) rather than a string is the load-bearing part: the
// claim is linked at compile time, so a nail that is renamed, moved behind a build
// tag this platform does not compile, or deleted cannot leave a registration
// pointing at nothing.
func registerLegNail131(leg string, test func(*testing.T)) legNail131 {
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
	n := legNail131{leg: leg, testName: name}
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
// gate needs and the source site each came from.
type leg131 struct {
	key          string
	site         string
	installs     bool
	installSites int
	emits        bool
	emitSites    []string
	ruled        bool
	rulingSites  []string
}

// mainPackage131 is the parsed package: production functions (callees, doc,
// site), test functions (callees), and func main's body.
type mainPackage131 struct {
	fset     *token.FileSet
	funcs    map[string]*funcInfo131
	tests    map[string]*funcInfo131
	mainBody *ast.BlockStmt
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
		fset:  token.NewFileSet(),
		funcs: map[string]*funcInfo131{},
		tests: map[string]*funcInfo131{},
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
	}
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

	var legs []leg131
	for _, cl := range sw.Body.List {
		cc, ok := cl.(*ast.CaseClause)
		if !ok {
			continue
		}
		keys := caseLabels131(cc)
		if len(keys) == 0 {
			keys = []string{"default"}
		}
		site := fmt.Sprintf("main.go:%d", pkg.fset.Position(cc.Pos()).Line)
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
	leg := leg131{key: key, site: site}
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
func caseLabels131(cl *ast.CaseClause) []string {
	var out []string
	for _, e := range cl.List {
		if bl, ok := e.(*ast.BasicLit); ok && bl.Kind == token.STRING {
			if un, err := strconv.Unquote(bl.Value); err == nil {
				out = append(out, un)
			}
		}
	}
	return out
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
