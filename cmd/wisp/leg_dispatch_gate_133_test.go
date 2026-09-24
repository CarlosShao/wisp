package main

// Ticket 133 AC#1 and AC#2: the dispatch-hop gate, an instrument that is NOT
// ticket 131's enumeration door.
//
// WHY A SECOND INSTRUMENT HAS TO EXIST. Ticket 131's gate asks, for every leg it
// enumerated out of main.go's dispatch shape: "does this leg reach
// installLogSink, and if so is it nailed?" Ticket 133's founding reading (R-121-1,
// measured by acceptor-ticket121 and re-measured by this ticket) is a question
// that gate structurally cannot ask: delete the ONE line that calls cmdModels and
// the leg's closure stops reaching installLogSink, so the row degrades from
// "nailed" to "no records" and every case in cmd/wisp and internal/models stays
// green. Two instruments aimed at the same fact would then agree, and the family
// would keep landing in the gap. So the checks below are built on premises that
// do NOT contain the phrase "does this leg reach the sink":
//
//	(a) a symbol that installs the listener is a production entry point and has to
//	    be reachable from main() - read from the sink INWARD (whoever calls
//	    installLogSink, whoever calls them, ... up to main), so deleting the hop
//	    anywhere on that chain is red no matter what the leg's own row says;
//	(b) a dispatched leg has to be covered - nailed, or driven by a test that
//	    touches something belonging to that leg alone, or ruled out in prose next to
//	    the code that owns it. This is the premise that keeps ticket 133's fifth
//	    shot red on its SECOND beat: once the leg's install block is deleted,
//	    install-based readings go quiet and "nobody anywhere verifies this command"
//	    does not.
//	(c) a coverage claim that names a test case has to name a case THIS ROUND'S
//	    TEST BINARY CAN START. Check (b) is read out of parsed sources, and parsed
//	    sources are a superset of what runs: a correctly shaped
//	    `func TestXxx(t *testing.T)` sitting in a file whose name carries an
//	    implicit GOOS=linux constraint is a declaration this platform never
//	    compiles, so it can never print `=== RUN`, and ticket 133's second
//	    acceptance round measured this file booking it as `covered=test ... drives
//	    ...` while the whole package stayed green (R-133-9, family shot M-G). The
//	    roster check below is the one reading here that does not come from
//	    go/parser: it asks the running binary which cases it can start. It is NOT
//	    folded back into the parser, because the point of a second instrument is a
//	    second, independent source of truth.
//
// WHAT EACH OF THE FIVE SHOTS IN AC#1 HITS, stated before the fact and measured
// afterwards (readings in docs/evidence/s1/133-ac1-ac2-instrument.md):
//
//	N-3   delete main.go's one line that calls cmdModels => check (a) is red:
//	      cmdModels still installs, and now nobody above it does.
//	X4    a leg dispatched by an early `if args[0] == "--diag"` before the switch
//	      => red three times over: the leg is enumerated (unlike 131, which reads
//	      argv outside the switch as a stray read), it installs, and it is covered
//	      by nothing.
//	X8    `case "slo":` rewritten as `case sloCmdName131:` => red: a leg is
//	      identified by a literal label, so a label this reader cannot turn into a
//	      command name is reported as a leg of its own instead of folding into
//	      default, and the operator-facing list in the usage block no longer matches
//	      the dispatch list in either direction.
//	X12   `var sinkAlias131 = installLogSink` plus a call through it, on a new leg
//	      => red: package-level function values are walked (that is 131's own fix
//	      for this shape), and the new leg is covered by nothing.
//	X14   install hidden in a struct field's initializer (`var holder131 =
//	      sinkHolder131{open: installLogSink}` + `holder131.open(root)`), new leg,
//	      no test file touched => red on the first beat (field edges are walked, so
//	      the leg owes a nail and has none; it is also covered by nothing) AND on
//	      the second beat, after the field initializer and that call are deleted
//	      (the leg is still dispatched and is still covered by nothing).
//
// WHAT THIS FILE CANNOT SEE, because an unbounded coverage claim is the thing this
// repository has learned not to trust (AC#2 of this ticket records the same list
// on the ticket face):
//
//   - a leg dispatched from a function other than main(), or by a command name
//     built at run time (fmt.Sprintf, a config file, a flag package): the
//     enumeration is over func main's own branches, and a callee that is not a
//     function this directory declares is a leaf. Calls in func main that sit
//     outside every branch this gate classifies are red, so a new dispatch SHAPE
//     costs whoever adds it an edit here.
//   - an install reached through a method value (`v.M` as a value), through an
//     interface field, or from a package other than cmd/wisp: field and map edges
//     here are resolved from package-level composite literals by name. Those
//     unplaced calls are reported as this instrument being blind rather than
//     walked past as leaves, but the reading stays a name reading.
//   - the records side: what a leg actually books into the sink is ticket 131's
//     predicate, not this one. Nothing here gets redder because a leg writes
//     slog.Info without a listener; that row belongs to the other door.
//   - what a covered leg's case ASSERTS. A name this round's binary can start is now
//     a measured fact (runRosterReds135 asks the running test binary for its own case
//     roster), but a case that starts, passes and checks nothing is ticket 131's
//     predicate, not this one; the per-leg body readings live there.
//   - a round narrowed by -test.run. The roster below is the set of cases this binary
//     can start in this round, which is the set `=== RUN` is drawn from; it is not a
//     transcript of the lines this particular process printed, because nothing inside
//     a running case can see the testing package's own output. A development call of
//     the shape `-run '^TestAC1AC2DispatchHopGate133$'` therefore proves "startable
//     here", not "started here" - the disclosure line names the filter in force so the
//     reading says which of the two it is. Every CI reading of this package runs the
//     whole package, where the two are the same set.
//   - the four claims in this gate's OWN registry that live in *_windows_test.go
//     files. On the windows leg they are in the roster and are counted; on a linux run
//     they are not, and that is disclosed in the same line rather than reddened,
//     because a red there would be this instrument claiming a coverage debt on a
//     platform whose cmd/wisp denominator is the other leg's business (ticket 133
//     R-133-2's measured "装不下"). The `covered=test` bucket - the one shot M-G was
//     filed for - has no such exemption.
//
// COST: one go/parser pass over this directory's ~30 files per run, plus ONE child
// process: the test binary re-lists its own cases with `-test.list '.*'` (measured
// 0.098 s for this package on this host, and it runs no test bodies). There is no
// other subprocess and no filesystem beyond this directory. The measurements below are
// milliseconds plus that listing.

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	// sinkFunc133 is the production entry that puts a persistent listener on a leg
	// (logsink.go). Check (a) walks inward from it.
	sinkFunc133 = "installLogSink"
	// handoffRoot133 is ticket 121's capability root, the other production entry
	// this repository has already been caught shipping with no caller. It is a
	// qualified name because it lives in another package.
	handoffRoot133 = "models.WireDownloading"
	// coverageRuling133 is the marker a leg without a driver uses to say so in
	// words, next to the code that owns the leg.
	coverageRuling133 = "WISP-LEG-COVERAGE-RULING:"
	// noArgsLeg133 names func main's `if len(args) == 0` branch, the one leg with
	// no case label to read.
	noArgsLeg133 = "no-args"
	// minLegs133 is the lower bound on the enumerated census. Measured at the
	// commit this file landed: nine literal case labels (run, providers, doctor,
	// secret, models, slo, panel-assets, version, help), the no-args branch, and
	// the switch's default. A dispatch rewrite that parses to nothing has to trip
	// this, not pass quietly - ticket 131's AC#4 names the same hole in its own
	// census, which is computed from main.go's shape rather than from a floor.
	minLegs133 = 11
)

// usageCommand133 reads one documented command line out of the usage block.
var usageCommand133 = regexp.MustCompile(`(?m)^  wisp ([a-z][a-z0-9-]*)`)

// usageBare133 matches the `wisp` line with no command token, which documents the
// no-args (GUI resident) leg.
var usageBare133 = regexp.MustCompile(`(?m)^  wisp {2,}\S`)

// legCover133 is this gate's OWN nail registry: a leg, the compiled test case that
// claims it, and the production symbol that case drives. It deliberately does not
// read legNails131 (ticket 131's registry): the point of ticket 133 is a reading
// that survives `-run` excluding TestAC4EveryLegIsNailedOrRuled, and a shared list
// would make the two doors one door.
type legCover133 struct {
	leg   string
	test  string
	entry string
}

var legCovers133 = []legCover133{
	{leg: "run", test: "TestAC2SealNoticeLandsInTheRunLegLogFile", entry: "runTextTask"},
	{leg: "models", test: "TestAC2ModelsLegBooksItsHandOffVerdictOnDisk", entry: "cmdModels"},
	{leg: "secret", test: "TestAC3SecretLegBooksItsAuditRecordsOnDisk", entry: "cmdSecret"},
	{leg: "no-args", test: "TestAC1ResidentLegInstallsItsLogListenerOnDisk", entry: "runResident"},
}

// TestAC1AC2DispatchHopGate133 is the gate. Every red below names a leg, a symbol
// and a file:line.
func TestAC1AC2DispatchHopGate133(t *testing.T) {
	pkg, err := loadPackage133(".")
	if err != nil {
		t.Fatalf("AC#1/#2 RED (the instrument, not the code): %v", err)
	}
	legs, reds := pkg.enumerateLegs133()
	for _, r := range reds {
		t.Errorf("AC#1/#2 RED: %s", r)
	}
	if len(legs) < minLegs133 {
		t.Errorf("AC#1/#2 RED: the leg census enumerated %d legs (%d is the floor, measured when this gate landed), so a dispatch rewrite is being read as \"no legs\" rather than reported. Census: %s",
			len(legs), minLegs133, strings.Join(legKeys133(legs), ", "))
	}

	sinks := pkg.sinkCallers133()
	reachable := pkg.reachableFromMain133(legs)

	var orphan []string
	for name := range sinks {
		if !reachable[name] {
			orphan = append(orphan, fmt.Sprintf("%s (%s)", name, pkg.sites133(name)))
		}
	}
	sort.Strings(orphan)
	for _, o := range orphan {
		t.Errorf("AC#1 RED: %s reaches %s (or is a production entry of this package's dispatch), and no chain from func main reaches it any more.\n"+
			"This is the reading ticket 121 measured for models.Manager and ticket 133 for the dispatch hop: the capability is compiled into the binary and nothing in the product walks to it. Restoring the deleted call is the fix; deleting this symbol is not, because the row above says who installed it.",
			o, sinkFunc133)
	}

	for _, r := range pkg.coverageReds133(legs) {
		t.Errorf("AC#1 RED: %s", r)
	}

	// The instrument's own blindness, reported instead of leafed past.
	for _, b := range pkg.blind {
		t.Errorf("AC#1/#2 RED (the instrument, not the code): %s\n"+
			"An edge this reader cannot place is an edge a leg can install the listener through without the ledger noticing, so it is red here rather than counted as a leaf.", b)
	}

	for _, r := range censusVsUsage133(pkg, legs) {
		t.Errorf("AC#1/#2 RED: %s", r)
	}

	// Check (c): every case name this round's ledger credited as coverage has to be a
	// case THIS test binary can start. The roster below is read out of the running
	// binary, which is the only reading in this file that does not come from go/parser
	// - the parser is what shot M-G got past, because a parsed declaration and a
	// compiled case are not the same object.
	roster, rosterN, rosterErr := compiledRunRoster135()
	for _, r := range pkg.runRosterReds135(legs, t.Name(), roster, rosterN, rosterErr) {
		t.Errorf("AC#1/#2 RED: %s", r)
	}

	var lines []string
	for _, leg := range legs {
		lines = append(lines, fmt.Sprintf("  leg %-14s %-24s installs=%-5v handoff=%-5v covered=%-62s entries=%s aliases=%s",
			leg.key, leg.site, leg.installs, leg.handoff, leg.covered, strings.Join(leg.entries, "|"), strings.Join(leg.aliases, "|")))
	}
	report := fmt.Sprintf("dispatch ledger, read out of func main's own branches at run time (%d legs, %d claims in this gate's registry):\n", len(legs), len(legCovers133)) +
		strings.Join(lines, "\n")
	if t.Failed() {
		t.Errorf("the dispatch ledger has a red row:\n%s", report)
	} else {
		t.Logf("%s", report)
	}
	t.Logf("blindness disclosure: %s", pkg.runRosterDisclosure135(legs, roster, rosterN, rosterErr))
}

// ---------------------------------------------------------------------------
// the reader
// ---------------------------------------------------------------------------

// decl133 is one function, method, or package-level closure in this directory.
type decl133 struct {
	name string // bare declaration name ("cmdModels", "logger")
	key  string // site-independent identity: "cmdModels" or "logSink.logger"
	site string // file:line
	file string
	prod bool // declared in a non-test file
	body *ast.BlockStmt
}

// ruling133 is one WISP-LEG-COVERAGE-RULING marker seen by the line scan: the leg
// name it carries, where it sits, and the file that holds it. The file is the
// point - R-133-3 is a ruling that backs a leg from a file that owns nothing about
// that leg, which is the same hole ticket 131's second ruler has as R-131r3-1 and
// that ticket 135 is closing over there. This instrument closes its own copy.
type ruling133 struct {
	leg  string
	site string
	file string
}

// leg133 is one dispatch branch of func main, with the readings the two checks
// above need.
type leg133 struct {
	key       string
	site      string
	isDefault bool
	aliases   []string // other labels of the same clause that fold into this row
	entries   []string
	quals     map[string]bool
	reached   map[string]bool
	path      map[string]string
	installs  bool
	handoff   bool
	covered   string
	// coveredBy holds the case names this row credited through its `covered=test`
	// form - the names a reader of the ledger above takes as "a test runs this leg".
	// It is what runRosterReds135 reconciles against the running binary's case
	// roster, so that the reconciliation reads the same claim the row prints rather
	// than re-deriving it from the parser.
	coveredBy []string
}

// pkg133 is the parsed directory. Build tags are ignored and same-named
// declarations are unioned, so a leg inherits the other platform's shape before it
// loses its own; the approximation runs in the direction that asks for more
// coverage, never less, and it is bounded by check (c): a case name this row
// credits is still reconciled against the roster of cases the running binary can
// start (runRosterReds135), which is where the borrowed shape stops being
// coverage.
type pkg133 struct {
	fset *token.FileSet

	decls   map[string][]*decl133 // production decls by bare name
	methods map[string][]*decl133 // production methods by bare name
	tests   map[string][]*decl133 // runnable cases in *_test.go: top-level func TestXxx(t *testing.T)
	helpers map[string][]*decl133 // everything else a *_test.go declares (helpers, methods, Benchmarks)

	globals   map[string]bool            // package-level var names
	varType   map[string]string          // package-level var -> struct or map type
	funcField map[string]bool            // "Type.field" -> the field is func-typed
	aliasOf   map[string]string          // "var" -> name it holds
	edges     map[string]string          // "var.field" / "var:key" -> name it holds
	unplaced  map[string]bool            // package-level func-typed var with no value
	imports   map[string]map[string]bool // file -> local import names
	consts    map[string]string          // package-level string constants

	mainBody  *ast.BlockStmt
	argvName  string
	blind     []string
	rulings   map[string]string // leg key -> file:line of the first ruling seen for it
	rulingAll []ruling133       // every ruling seen, with the file it sits in (R-133-3)

	knownPkgs map[string]bool         // every import name used anywhere in the directory
	cache     map[*decl133]*calleeSet // resolved call sets, memoised
	loose     map[*decl133]*calleeSet // same, with test-file decls resolved too
}

// calleeSet is one declaration's resolved call set: the declarations it reaches,
// the qualified external names it calls (capability roots in another package live
// there), and for each reached declaration the edge that got there, so a ledger row
// can quote the alias or the field it walked instead of pretending every edge is a
// plain name.
type calleeSet struct {
	decls []*decl133
	quals []string
	via   map[string]string
}

func newPkg133() *pkg133 {
	return &pkg133{
		fset:      token.NewFileSet(),
		decls:     map[string][]*decl133{},
		methods:   map[string][]*decl133{},
		tests:     map[string][]*decl133{},
		helpers:   map[string][]*decl133{},
		globals:   map[string]bool{},
		varType:   map[string]string{},
		funcField: map[string]bool{},
		aliasOf:   map[string]string{},
		edges:     map[string]string{},
		unplaced:  map[string]bool{},
		imports:   map[string]map[string]bool{},
		consts:    map[string]string{},
		rulings:   map[string]string{},
		knownPkgs: map[string]bool{},
		cache:     map[*decl133]*calleeSet{},
		loose:     map[*decl133]*calleeSet{},
	}
}

// loadPackage133 parses every .go file in dir twice over: first the declarations
// and the package-level values a name walk needs, then the bodies.
func loadPackage133(dir string) (*pkg133, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", dir, err)
	}
	p := newPkg133()
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") {
			continue
		}
		file, perr := parser.ParseFile(p.fset, filepath.Join(dir, name), nil, parser.ParseComments)
		if perr != nil {
			return nil, fmt.Errorf("parse %s: %w", name, perr)
		}
		is := localImportNames133(file)
		p.imports[name] = is
		for k := range is {
			p.knownPkgs[k] = true
		}
		p.collectTopLevel133(file, name, testingPkgName133(file))
	}
	p.collectRulings133(dir, entries)
	return p, nil
}

func localImportNames133(f *ast.File) map[string]bool {
	out := map[string]bool{}
	for _, im := range f.Imports {
		path, err := strconv.Unquote(im.Path.Value)
		if err != nil {
			continue
		}
		local := path
		if i := strings.LastIndex(local, "/"); i >= 0 {
			local = local[i+1:]
		}
		if im.Name != nil {
			local = im.Name.Name
		}
		out[local] = true
	}
	return out
}

// exprIdent returns the bare identifier of an expression, if that is what it is.
func exprIdent(e ast.Expr) string {
	if id, ok := e.(*ast.Ident); ok {
		return id.Name
	}
	return ""
}

// qualifiedName renders pkg.Sel for a selector on a known import.
func qualifiedName133(sel *ast.SelectorExpr) string {
	if id, ok := sel.X.(*ast.Ident); ok {
		return id.Name + "." + sel.Sel.Name
	}
	return ""
}

// funcTargetName turns the right-hand side of a function-valued assignment into
// the name it holds ("" when this reader does not model the value).
func funcTargetName133(e ast.Expr) string {
	switch v := e.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return qualifiedName133(v)
	}
	return ""
}

func (p *pkg133) collectTopLevel133(f *ast.File, name string, testingPkg string) {
	isTest := strings.HasSuffix(name, "_test.go")
	for _, d := range f.Decls {
		switch fd := d.(type) {
		case *ast.FuncDecl:
			if fd.Body == nil {
				continue
			}
			key := fd.Name.Name
			if fd.Recv != nil {
				key = recvName133(fd.Recv.List[0].Type) + "." + fd.Name.Name
			}
			decl := &decl133{name: fd.Name.Name, key: key, site: p.site133(fd.Pos()), file: name, prod: !isTest, body: fd.Body}
			switch {
			case isTest && isRunnableCase133(fd, testingPkg):
				// R-133-1: the coverage bucket holds ONLY declarations shaped the way
				// Go's testing package runs them - a top-level func TestXxx(t
				// *testing.T) with no receiver and no results. Before this, any
				// declaration in a _test.go file landed here under its bare name, so a
				// method (which the testing package never calls) satisfied check (b)'s
				// "a test case drives this leg", and the ledger printed
				// `covered=test <name>` for a name with no === RUN anywhere.
				//
				// Shape is still not existence. This walk reads the directory, not the
				// build, so a correctly shaped case in a file this platform does not
				// compile (a `_linux_test.go` suffix, a false //go:build line) lands
				// here exactly as a real one does - that is R-133-9, family shot M-G,
				// and it is why a name leaving this bucket as coverage is reconciled
				// against the running binary's own case roster in runRosterReds135
				// before the ledger row is believed.
				p.tests[fd.Name.Name] = append(p.tests[fd.Name.Name], decl)
			case isTest:
				// Helpers, methods, benchmarks: still walkable in loose mode, so a
				// case that drives a leg THROUGH a helper keeps being recognised as
				// that case driving the leg. They are not cases, and they can no
				// longer cover a leg by being named after one.
				p.helpers[fd.Name.Name] = append(p.helpers[fd.Name.Name], decl)
			case fd.Recv != nil:
				p.methods[fd.Name.Name] = append(p.methods[fd.Name.Name], decl)
			default:
				p.decls[fd.Name.Name] = append(p.decls[fd.Name.Name], decl)
			}
			if !isTest && fd.Recv == nil && fd.Name.Name == "main" {
				p.mainBody = fd.Body
			}
		case *ast.GenDecl:
			if isTest {
				continue
			}
			p.collectSpecs133(fd, name)
		}
	}
}

// testingPkgName133 is the local name `testing` answers to in one file, which is
// what makes the signature check below readable when a file imports it under an
// alias.
func testingPkgName133(f *ast.File) string {
	for _, im := range f.Imports {
		path, err := strconv.Unquote(im.Path.Value)
		if err != nil || path != "testing" {
			continue
		}
		if im.Name != nil {
			return im.Name.Name
		}
		return "testing"
	}
	return "testing"
}

// isRunnableCase133 is R-133-1's judgement, spelled the way ticket 131's gate
// spells its own (leg_sink_gate_131_test.go:617 checks `fn.Recv == nil &&
// HasPrefix(fn.Name.Name, "Test")` before a test declaration is a case): no
// receiver, a Test prefix, one parameter of type *<testing>.T, and no results.
// Go's testing package runs exactly that shape and nothing else, so this is the
// difference between "a name that looks like a case" and "a name shaped like a
// case". Whether a name shaped like one is a case this round can start is the next
// question down, and it is answered from the binary rather than from here.
func isRunnableCase133(fd *ast.FuncDecl, testingPkg string) bool {
	if fd.Recv != nil || !strings.HasPrefix(fd.Name.Name, "Test") {
		return false
	}
	if fd.Type == nil || fd.Type.Params == nil || len(fd.Type.Params.List) != 1 {
		return false
	}
	if fd.Type.Results != nil && len(fd.Type.Results.List) > 0 {
		return false
	}
	if len(fd.Type.Params.List[0].Names) > 1 {
		return false
	}
	return isTestingTPtr133(fd.Type.Params.List[0].Type, testingPkg)
}

func isTestingTPtr133(e ast.Expr, testingPkg string) bool {
	star, ok := e.(*ast.StarExpr)
	if !ok {
		return false
	}
	sel, ok := star.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	id, ok := sel.X.(*ast.Ident)
	return ok && id.Name == testingPkg && sel.Sel.Name == "T"
}

func (p *pkg133) collectSpecs133(gd *ast.GenDecl, file string) {
	for _, s := range gd.Specs {
		switch spec := s.(type) {
		case *ast.TypeSpec:
			if st, ok := spec.Type.(*ast.StructType); ok {
				for _, fld := range st.Fields.List {
					if _, isFunc := fld.Type.(*ast.FuncType); !isFunc {
						continue
					}
					for nm := range fieldNames133(fld) {
						p.funcField[spec.Name.Name+"."+nm] = true
					}
				}
			}
		case *ast.ValueSpec:
			if gd.Tok == token.CONST {
				p.collectConst133(spec)
			}
			p.collectValueSpec133(spec, file)
		}
	}
}

// collectConst133 keeps package-level string constants, which is what the usage
// block is made of and what a non-literal case label resolves to.
func (p *pkg133) collectConst133(vs *ast.ValueSpec) {
	if len(vs.Values) != len(vs.Names) {
		return
	}
	for i, nm := range vs.Names {
		if bl, ok := vs.Values[i].(*ast.BasicLit); ok && bl.Kind == token.STRING {
			if v, err := strconv.Unquote(bl.Value); err == nil {
				p.consts[nm.Name] = v
			}
		}
	}
}

func (p *pkg133) collectValueSpec133(vs *ast.ValueSpec, file string) {
	for _, nm := range varNames133(vs) {
		p.globals[nm.Name] = true
	}
	if len(vs.Values) == 0 {
		// `var holder131 sinkHolder131` - declared, never placed here. A call
		// through one of its func fields is reported as blindness, not as a leaf.
		for _, nm := range varNames133(vs) {
			if t := typeName133(vs.Type); t != "" {
				p.varType[nm.Name] = t
				p.unplaced[nm.Name] = true
			}
		}
		return
	}
	if len(vs.Names) != 1 || len(vs.Values) != 1 {
		return
	}
	name := vs.Names[0].Name
	switch v := vs.Values[0].(type) {
	case *ast.Ident, *ast.SelectorExpr:
		if t := funcTargetName133(vs.Values[0]); t != "" {
			p.aliasOf[name] = t
		}
	case *ast.FuncLit:
		p.decls[name] = append(p.decls[name], &decl133{
			name: name, key: name + "#closure", site: p.site133(v.Pos()),
			file: file, body: v.Body,
		})
	case *ast.CompositeLit:
		t := typeName133(v.Type)
		if t != "" {
			p.varType[name] = t
		}
		for _, el := range v.Elts {
			kv, ok := el.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			target := funcTargetName133(kv.Value)
			if target == "" {
				continue
			}
			switch k := kv.Key.(type) {
			case *ast.Ident:
				p.edges[name+"."+k.Name] = target
			case *ast.BasicLit:
				if key, err := strconv.Unquote(k.Value); err == nil {
					p.edges[name+":"+key] = target
				}
			}
		}
	}
}

// recordBodies133 is gone: the call walk is lazy (callees133), so it runs after
// every package-level name in the directory is known.

// ---------------------------------------------------------------------------
// the call walk
// ---------------------------------------------------------------------------

// site133 renders a position as file:line, the form every red below quotes.
func (p *pkg133) site133(pos token.Pos) string {
	posn := p.fset.Position(pos)
	return filepath.Base(posn.Filename) + ":" + strconv.Itoa(posn.Line)
}

func recvName133(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.StarExpr:
		return recvName133(t.X)
	case *ast.Ident:
		return t.Name
	case *ast.IndexExpr:
		return t.X.(*ast.Ident).Name
	case *ast.IndexListExpr:
		return t.X.(*ast.Ident).Name
	}
	return "?"
}

func fieldNames133(f *ast.Field) map[string]bool {
	out := map[string]bool{}
	for _, id := range f.Names {
		out[id.Name] = true
	}
	return out
}

func varNames133(vs *ast.ValueSpec) []*ast.Ident { return vs.Names }

// typeName133 names the type of a var declaration or composite literal, star and
// qualifiers stripped, which is the form funcField keys are stored in.
func typeName133(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.StarExpr:
		return typeName133(t.X)
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return qualifiedName133(t)
	case *ast.ArrayType:
		return ""
	case *ast.MapType:
		return "map"
	}
	return ""
}

var builtins133 = map[string]bool{
	"append": true, "cap": true, "close": true, "complex": true, "copy": true,
	"delete": true, "imag": true, "len": true, "make": true, "max": true,
	"min": true, "new": true, "panic": true, "print": true, "println": true,
	"real": true, "recover": true, "clear": true,
}

// resolve133 turns one call expression into the declarations it can reach. The
// forms are the ones ticket 133's five shots are made of, and each is named in the
// ledger's via column so a reading can be traced to the edge it came from:
//
//	name            a function or closure var declared in this directory
//	pkg.Name        an imported package (a leaf; capability roots live here)
//	x.M             a method declared here, receivers unioned by name
//	holder.field()  a package-level struct var whose field initializer named a func
//	table["k"]()    the same shape over a package-level map literal of funcs
//	f := <name>     a local alias for one of the above, inside the body being walked
//
// The direction of every guess here is the one that asks for more coverage: a
// method call resolves to every receiver of that name, and an edge this reader
// cannot place is reported in pkg.blind instead of being walked past as a leaf.
// lookup133 finds declarations by bare name. In loose mode (used only for the
// test-side evidence walk) it also sees the functions declared in this
// directory's test files, which is how a case that drives a leg through a helper
// is still recognised as that case driving the leg. Both test buckets are in
// scope here - runnable cases and helpers - while only p.tests counts as coverage
// evidence anywhere else in this file.
func (p *pkg133) lookup133(name string, loose bool) []*decl133 {
	out := append([]*decl133{}, p.decls[name]...)
	out = append(out, p.methods[name]...)
	if loose {
		out = append(out, p.tests[name]...)
		out = append(out, p.helpers[name]...)
	}
	return out
}

// testSide133 is every declaration a *_test.go file contributes to the loose walk.
func (p *pkg133) testSide133(name string) []*decl133 {
	out := append([]*decl133{}, p.tests[name]...)
	return append(out, p.helpers[name]...)
}

func (p *pkg133) resolve133(c *ast.CallExpr, file string, locals map[string]string, loose bool) (ds []*decl133, disp string, qual string) {
	// edgeTarget resolves a name written on the right-hand side of an edge.
	edgeTarget := func(target string) ([]*decl133, string) {
		if i := strings.LastIndex(target, "."); i > 0 && p.knownPkgs[target[:i]] {
			return nil, target
		}
		return p.lookup133(target, loose), ""
	}
	switch fn := ast.Unparen(c.Fun).(type) {
	case *ast.Ident:
		if builtins133[fn.Name] {
			return nil, "", ""
		}
		if fs := p.lookup133(fn.Name, false); len(fs) > 0 {
			return fs, fn.Name, ""
		}
		if t, ok := locals[fn.Name]; ok {
			got, q := edgeTarget(t)
			return got, fn.Name + " -> " + t, q
		}
		if t, ok := p.aliasOf[fn.Name]; ok {
			if p.unplaced[fn.Name] {
				p.blind = append(p.blind, fmt.Sprintf("%s calls %s, a package-level function value this directory declares but never defines, so the walk can follow the name and no further", p.site133(c.Pos()), fn.Name))
				return nil, fn.Name, ""
			}
			got, q := edgeTarget(t)
			return got, fn.Name + " -> " + t, q
		}
		if p.unplaced[fn.Name] {
			p.blind = append(p.blind, fmt.Sprintf("%s calls %s, a package-level function value this directory declares but never defines", p.site133(c.Pos()), fn.Name))
			return nil, fn.Name, ""
		}
		if loose {
			if fs := p.testSide133(fn.Name); len(fs) > 0 {
				return fs, fn.Name, ""
			}
		}
		return nil, fn.Name, ""

	case *ast.SelectorExpr:
		x := exprIdent(fn.X)
		if x != "" && p.imports[file][x] {
			qn := qualifiedName133(fn)
			return nil, qn, qn
		}
		if x != "" && p.globals[x] {
			field := x + "." + fn.Sel.Name
			if t, ok := p.edges[field]; ok {
				got, q := edgeTarget(t)
				return got, field + " -> " + t, q
			}
			if p.varType[x] != "" && p.funcField[p.varType[x]+"."+fn.Sel.Name] {
				p.blind = append(p.blind, fmt.Sprintf("%s calls through field %s of package-level var %s (type %s), whose value this reader cannot place: a function carried in a struct field is the shape ticket 131's own header lists as its remaining window",
					p.site133(c.Pos()), fn.Sel.Name, x, p.varType[x]))
				return nil, field, ""
			}
			if ms := p.methods[fn.Sel.Name]; len(ms) > 0 {
				return ms, fn.Sel.Name, ""
			}
			return nil, field, ""
		}
		if ms := p.methods[fn.Sel.Name]; len(ms) > 0 {
			return ms, fn.Sel.Name, ""
		}
		return nil, "", ""

	case *ast.IndexExpr:
		base := exprIdent(ast.Unparen(fn.X))
		if base == "" || !p.globals[base] {
			return nil, "", ""
		}
		bl, ok := fn.Index.(*ast.BasicLit)
		if !ok {
			return nil, "", ""
		}
		key, kerr := strconv.Unquote(bl.Value)
		if kerr != nil {
			return nil, "", ""
		}
		table := base + ":" + key
		if t, ok := p.edges[table]; ok {
			got, q := edgeTarget(t)
			return got, table + " -> " + t, q
		}
		return nil, table, ""
	}
	return nil, "", ""
}

// callees133 is one declaration's resolved call set, walked lazily and memoised.
// Loose mode is only used from the test side, so the production graph never picks
// up a helper name.
func (p *pkg133) callees133(d *decl133, loose bool) *calleeSet {
	cache := p.cache
	if loose {
		cache = p.loose
	}
	if cs, ok := cache[d]; ok {
		return cs
	}
	cs := &calleeSet{via: map[string]string{}}
	cache[d] = cs
	if d.body == nil {
		return cs
	}
	locals := map[string]string{}
	ast.Inspect(d.body, func(n ast.Node) bool {
		switch st := n.(type) {
		case *ast.AssignStmt:
			// A local alias for a function value, the one-edit sibling of the
			// package-level form ticket 131 registered as its third shape.
			if st.Tok == token.DEFINE && len(st.Lhs) == 1 && len(st.Rhs) == 1 {
				if nm := exprIdent(st.Lhs[0]); nm != "" {
					if t := funcTargetName133(st.Rhs[0]); t != "" {
						locals[nm] = t
					}
				}
			}
		case *ast.CallExpr:
			got, disp, qual := p.resolve133(st, d.file, locals, loose)
			if qual != "" {
				cs.quals = append(cs.quals, qual)
			}
			for _, target := range got {
				cs.decls = append(cs.decls, target)
				if _, seen := cs.via[target.key]; !seen {
					cs.via[target.key] = disp
				}
			}
		}
		return true
	})
	return cs
}

// allDecls133 is every production declaration in the directory, methods included.
func (p *pkg133) allDecls133() []*decl133 {
	var out []*decl133
	for _, ds := range p.decls {
		out = append(out, ds...)
	}
	for _, ds := range p.methods {
		out = append(out, ds...)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].key < out[j].key })
	return out
}

// closureOf133 walks a set of seeds to a fixpoint. reached keys are declaration
// identities; path remembers the edges taken, which is what the ledger prints.
func (p *pkg133) closureOf133(seeds []*decl133) (map[string]bool, map[string]string, map[string]bool) {
	reached := map[string]bool{}
	path := map[string]string{}
	quals := map[string]bool{}
	seen := map[*decl133]bool{}
	var queue []*decl133
	enqueue := func(d *decl133, via string) {
		reached[d.key] = true
		if _, ok := path[d.key]; !ok {
			path[d.key] = via
		}
		if seen[d] {
			return
		}
		seen[d] = true
		queue = append(queue, d)
	}
	for _, s := range seeds {
		enqueue(s, "dispatch")
	}
	for len(queue) > 0 {
		d := queue[0]
		queue = queue[1:]
		cs := p.callees133(d, false)
		for _, q := range cs.quals {
			quals[q] = true
		}
		for _, c := range cs.decls {
			via := path[d.key] + " -> " + c.name
			if v := cs.via[c.key]; v != "" {
				via = path[d.key] + " -> " + v
			}
			enqueue(c, via)
		}
	}
	return reached, path, quals
}

// sites133 renders every file:line a declaration key sits at, the form a red quotes
// for a name declared in two platform files.
func (p *pkg133) sites133(key string) string {
	var sites []string
	for _, d := range p.allDecls133() {
		if d.key == key {
			sites = append(sites, d.site)
		}
	}
	if len(sites) == 0 {
		return "not declared in this directory"
	}
	return strings.Join(sites, " | ")
}

func (p *pkg133) hasDecl133(name string) bool {
	return len(p.decls[name]) > 0 || len(p.methods[name]) > 0
}

// sinkCallers133 is check (a)'s premise, read from the capability inward: every
// production symbol that reaches installLogSink, or the ticket 121 hand-off root,
// transitively. Deleting any hop on the way out to main() is what this set makes
// visible, which is why it does not care what the leg's own row says.
func (p *pkg133) sinkCallers133() map[string]bool {
	preds := map[string]map[string]bool{}
	add := func(callee, caller string) {
		if preds[callee] == nil {
			preds[callee] = map[string]bool{}
		}
		preds[callee][caller] = true
	}
	for _, d := range p.allDecls133() {
		cs := p.callees133(d, false)
		for _, c := range cs.decls {
			add(c.key, d.key)
		}
		for _, q := range cs.quals {
			add(q, d.key)
		}
	}
	out := map[string]bool{}
	seen := map[string]bool{sinkFunc133: true, handoffRoot133: true}
	queue := []string{sinkFunc133, handoffRoot133}
	for len(queue) > 0 {
		k := queue[0]
		queue = queue[1:]
		for caller := range preds[k] {
			if seen[caller] {
				continue
			}
			seen[caller] = true
			out[caller] = true
			queue = append(queue, caller)
		}
	}
	return out
}

// reachableFromMain133 is the union of every enumerated leg's closure.
func (p *pkg133) reachableFromMain133(legs []*leg133) map[string]bool {
	out := map[string]bool{"main": true}
	for _, leg := range legs {
		for k := range leg.reached {
			out[k] = true
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// the leg census, read out of func main
// ---------------------------------------------------------------------------

// enumerateLegs133 classifies every branch of func main. A branch it cannot
// classify is red rather than skipped, because the leg list is the census every
// other check in this file is computed over: X4's early return is a leg here
// (ticket 131 reads the same statement as a stray argv read), and a label that is
// not a literal becomes a leg of its own instead of folding into default (X8).
func (p *pkg133) enumerateLegs133() ([]*leg133, []string) {
	var reds []string
	if p.mainBody == nil {
		return nil, []string{"func main is not declared in this directory, so there is no dispatch to read"}
	}
	for _, st := range p.mainBody.List {
		as, ok := st.(*ast.AssignStmt)
		if !ok || as.Tok != token.DEFINE || len(as.Lhs) != 1 || len(as.Rhs) != 1 {
			continue
		}
		// `args := os.Args[1:]` (a slice) and `cmd := os.Args[0]` (an index) are
		// both bindings of a local to the process argv.
		var rhs ast.Expr
		switch v := as.Rhs[0].(type) {
		case *ast.IndexExpr:
			rhs = v.X
		case *ast.SliceExpr:
			rhs = v.X
		}
		if sel, ok := rhs.(*ast.SelectorExpr); ok && qualifiedName133(sel) == "os.Args" {
			p.argvName = exprIdent(as.Lhs[0])
		}
	}
	if p.argvName == "" {
		return nil, []string{"func main binds no local to os.Args, so this gate cannot tell a dispatch from a call"}
	}

	var legs []*leg133
	var ranges [][2]int

	addLeg := func(key, site string, n ast.Node, isDefault bool, aliases []string, extraReds ...string) {
		reds = append(reds, extraReds...)
		ranges = append(ranges, [2]int{int(n.Pos()), int(n.End())})
		leg := &leg133{key: key, site: site, isDefault: isDefault, aliases: dedupe133(aliases)}
		var seeds []*decl133
		var entries []string
		ast.Inspect(n, func(x ast.Node) bool {
			ce, ok := x.(*ast.CallExpr)
			if !ok {
				return true
			}
			got, disp, qual := p.resolve133(ce, "main.go", nil, false)
			if len(got) == 0 {
				if qual == "" && isDispatchableSpelling133(disp) {
					reds = append(reds, fmt.Sprintf("leg %q (%s): %s calls %s, which this directory does not declare as a function. A dispatch this gate resolves to nothing is a leg it cannot read, and \"cannot read\" is not \"nothing to check\".",
						key, site, p.site133(ce.Pos()), disp))
				}
				return true
			}
			for _, d := range got {
				seeds = append(seeds, d)
				entries = append(entries, d.key)
			}
			return true
		})
		leg.entries = dedupe133(entries)
		leg.reached, leg.path, leg.quals = p.closureOf133(seeds)
		if leg.reached[sinkFunc133] {
			leg.installs = true
		}
		if leg.quals[handoffRoot133] {
			leg.handoff = true
		}
		legs = append(legs, leg)
	}

	for _, st := range p.mainBody.List {
		switch s := st.(type) {
		case *ast.IfStmt:
			keys, aliases, kind := p.classifyCond133(s.Cond)
			site := p.site133(s.Pos())
			switch kind {
			case condNoArgs133, condLiteral133:
				// R-133-5: one row per command this condition dispatches.
				for i, k := range keys {
					if i == 0 {
						addLeg(k, site, s.Body, false, aliases)
						continue
					}
					addLeg(k, site, s.Body, false, nil)
				}
			default:
				reds = append(reds, fmt.Sprintf("%s: func main has an if branch this gate cannot classify (%s). A dispatch outside the leg census needs no nail and no ruling from anything here, which is the hole ticket 131 registered as its first shape, so it is red rather than skipped.",
					site, p.render133(s.Cond)))
			}
		case *ast.SwitchStmt:
			if !p.isArgvSwitch133(s.Tag) {
				reds = append(reds, fmt.Sprintf("%s: func main has a switch whose tag is not argv (%s), so its cases are dispatches this census does not read.",
					p.site133(s.Pos()), p.render133(s.Tag)))
				continue
			}
			for _, cl := range s.Body.List {
				cc, ok := cl.(*ast.CaseClause)
				if !ok {
					continue
				}
				if cc.List == nil {
					addLeg("default", p.site133(cc.Pos()), cc, true, nil)
					continue
				}
				// R-133-4: a clause can carry several labels, and every one of them is a
				// command an operator can type. Before this only the FIRST literal label
				// reached the census, so `case "run", "runalt133":` printed one row and the
				// second command name existed nowhere in this reading - not in the ledger,
				// not in the usage reconciliation, not as a leg owing a nail. Each label
				// that reads as a command gets its own row now; only a dash spelling that
				// abbreviates the clause's own command word is folded, and it is folded
				// into that row's aliases column rather than dropped.
				keys, aliases, labelReds := p.caseLabelLegs133(cc)
				for i, k := range keys {
					if i == 0 {
						addLeg(k, p.site133(cc.Pos()), cc, false, aliases, labelReds...)
						continue
					}
					addLeg(k, p.site133(cc.Pos()), cc, false, nil)
				}
			}
		}
	}

	// Every call in func main has to sit inside a branch the census read.
	ast.Inspect(p.mainBody, func(n ast.Node) bool {
		ce, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if builtins133[exprIdent(ce.Fun)] {
			return true
		}
		for _, r := range ranges {
			if int(ce.Pos()) > r[0] && int(ce.Pos()) < r[1] {
				return true
			}
		}
		reds = append(reds, fmt.Sprintf("%s: func main calls %s outside every branch this census classified, so whatever it dispatches never enters the leg list. Move it into a branch the census reads, or teach enumerateLegs133 the shape and make it prove the reading on a planted leg.",
			p.site133(ce.Pos()), p.render133(ce.Fun)))
		return true
	})

	sort.Slice(legs, func(i, j int) bool { return legs[i].key < legs[j].key })
	return legs, reds
}

// caseLabelLegs133 turns one case clause into the legs it dispatches. The rule is
// the one a CLI reader would state out loud:
//
//   - a literal label that reads as a command word is a leg of its own, with its
//     own coverage obligation and its own line owed in the usage block. That is
//     what R-133-4 was: `case "run", "runalt133":` used to enumerate "run" only,
//     so a second command could be wired to any body and appear nowhere.
//   - a literal label that starts with a dash is a flag SPELLING of the clause's
//     command word, and it folds into that row's aliases column - but only when
//     its letters are a prefix of that word ("-v" and "--version" of "version").
//     The four flag spellings in this package's main.go are documented as the
//     commands they abbreviate, so folding them keeps today's reconciliation
//     honest; a dash label that abbreviates nothing ("-x133") is not an alias of
//     anything and becomes a leg that owes both a nail and a usage line.
//   - a label that is not a literal at all is its own synthetic leg, and red, as
//     ticket 131's second shape already required (X8's reading is unchanged).
//
// Returned in source order: keys are the legs to add, aliases belong to keys[0].
func (p *pkg133) caseLabelLegs133(cc *ast.CaseClause) (keys, aliases, reds []string) {
	site := p.site133(cc.Pos())
	var values []string   // the literal labels, in source order
	var unparsed []string // one synthetic row per label this reader cannot read as a name
	for _, l := range cc.List {
		bl, ok := l.(*ast.BasicLit)
		if !ok || bl.Kind != token.STRING {
			expr := p.render133(l)
			resolved := ""
			if nm := exprIdent(l); nm != "" {
				if v, ok2 := p.consts[nm]; ok2 {
					resolved = ", which resolves to the literal " + strconv.Quote(v)
				}
			}
			name := "unparsed-label@" + site + ":" + expr
			unparsed = append(unparsed, name)
			reds = append(reds, fmt.Sprintf("%s: the case label %s is not a string literal%s. A leg is identified by a literal command name; a label that has to be resolved through a constant is a leg that can leave the census without saying so, which is ticket 131's second shape. This row carries it as %q instead of folding it into default.",
				p.site133(l.Pos()), expr, resolved, name))
			continue
		}
		v, err := strconv.Unquote(bl.Value)
		if err != nil {
			continue
		}
		if v == "" {
			continue
		}
		values = append(values, v)
	}
	keys, aliases = splitCommandLabels133(values)
	keys = append(keys, unparsed...)
	if len(keys) == 0 {
		keys = []string{"unparsed-label@" + site}
	}
	return dedupe133Stable133(keys), dedupe133Stable133(aliases), reds
}

// dedupe133Stable133 is dedupe133 without the sort: which label of a clause comes
// first is a fact about the source, and the census row order is what a reader
// compares against the switch.
func dedupe133Stable133(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		if seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}

type condKind133 int

const (
	condUnknown133 condKind133 = iota
	condNoArgs133
	condLiteral133
)

// classifyCond133 reads the two argv conditions this gate dispatches on: the
// `len(argv) == 0` GUI branch, and any branch comparing an argv slot to a literal
// (X4's `if len(args) > 0 && args[0] == "--diag"` is the second form).
//
// R-133-5: it used to return ONE key, and it got that one key two different ways -
// the literal branch overwrote (last match wins), the no-args branch was guarded by
// `!found` (first match wins). So `if args[0] == "x133a" || args[0] == "version"`
// enumerated "version" and dropped "x133a" entirely: that command name appeared
// nowhere in the ledger, owed no nail, and never met the usage block, which is
// ticket 133's own acceptor reading for the case-label shape (R-133-4) arriving a
// second time through the if branch. Every classifiable comparison in the condition
// now yields its own leg, with the same flag-spelling fold R-133-4 uses.
func (p *pkg133) classifyCond133(e ast.Expr) (keys, aliases []string, kind condKind133) {
	var literals []string
	noArgs := false
	ast.Inspect(e, func(n ast.Node) bool {
		bin, ok := n.(*ast.BinaryExpr)
		if !ok {
			return true
		}
		if bin.Op == token.EQL && (p.isArgvSlot133(bin.X) || p.isArgvSlot133(bin.Y)) {
			other := bin.Y
			if p.isArgvSlot133(bin.Y) {
				other = bin.X
			}
			if bl, ok := other.(*ast.BasicLit); ok && bl.Kind == token.STRING {
				if v, err := strconv.Unquote(bl.Value); err == nil {
					literals = append(literals, v)
				}
			}
			return true
		}
		if call, ok := bin.X.(*ast.CallExpr); ok && exprIdent(call.Fun) == "len" && len(call.Args) == 1 && bin.Op == token.EQL {
			if id, ok := call.Args[0].(*ast.Ident); ok && id.Name == p.argvName {
				if bl, ok := bin.Y.(*ast.BasicLit); ok && bl.Value == "0" {
					noArgs = true
				}
			}
		}
		return true
	})
	keys, aliases = splitCommandLabels133(literals)
	switch {
	case len(keys) > 0 && noArgs:
		// A condition that both tests for no arguments and compares an argv slot to
		// a literal dispatches two censuses' worth of commands; neither is dropped.
		return append([]string{noArgsLeg133}, keys...), aliases, condLiteral133
	case len(keys) > 0:
		return keys, aliases, condLiteral133
	case noArgs:
		return []string{noArgsLeg133}, nil, condNoArgs133
	}
	return nil, nil, condUnknown133
}

// splitCommandLabels133 is the one place that decides which of several spellings in
// a single dispatch condition is a command of its own and which is a flag alias of
// the clause's command word. R-133-4 (case labels) and R-133-5 (argv comparisons in
// an if) read the same rule from here so the two cannot drift apart.
func splitCommandLabels133(values []string) (keys, aliases []string) {
	word := ""
	for _, v := range values {
		if !strings.HasPrefix(v, "-") && word == "" {
			word = v
		}
	}
	for _, v := range values {
		if v == "" {
			continue
		}
		if !strings.HasPrefix(v, "-") {
			keys = append(keys, v)
			continue
		}
		short := strings.TrimLeft(v, "-")
		if word != "" && short != "" && strings.HasPrefix(word, short) {
			aliases = append(aliases, v)
			continue
		}
		keys = append(keys, v)
	}
	return dedupe133Stable133(keys), dedupe133Stable133(aliases)
}

func (p *pkg133) isArgvSlot133(e ast.Expr) bool {
	idx, ok := e.(*ast.IndexExpr)
	if !ok {
		return false
	}
	id, ok := idx.X.(*ast.Ident)
	return ok && id.Name == p.argvName
}

func (p *pkg133) isArgvSwitch133(tag ast.Expr) bool {
	if tag == nil {
		return false
	}
	if idx, ok := tag.(*ast.IndexExpr); ok {
		if slice, ok := idx.X.(*ast.Ident); ok {
			return slice.Name == p.argvName
		}
	}
	if id, ok := tag.(*ast.Ident); ok {
		return id.Name == p.argvName
	}
	return false
}

// isDispatchableSpelling133 keeps the unplaced-call reading to lowercase names,
// which is where a command handler lives, and drops the edge labels and type
// conversions this walk also produces.
func isDispatchableSpelling133(disp string) bool {
	if disp == "" || strings.Contains(disp, " -> ") {
		return false
	}
	if strings.Contains(disp, ".") {
		return false
	}
	return disp[0] >= 'a' && disp[0] <= 'z'
}

func (p *pkg133) render133(e ast.Expr) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, p.fset, e); err != nil {
		return "<unrenderable>"
	}
	out := strings.Join(strings.Fields(buf.String()), " ")
	if len(out) > 120 {
		out = out[:120] + "..."
	}
	return out
}

// ---------------------------------------------------------------------------
// coverage, rulings and the usage census
// ---------------------------------------------------------------------------

// coverageReds133 is check (b): a dispatched leg has to be covered by something.
// Three forms count, and the ledger names which one it found: a nail claimed by this
// file's own registry (a name that is a func TestXxx(t *testing.T) in this
// directory's sources, see isRunnableCase133), a test case that drives a symbol
// belonging to that leg alone, or a WISP-LEG-COVERAGE-RULING sentence that both names
// the leg and sits in a file owning that leg's dispatch or one of its entries.
//
// The names behind the `covered=test` form are recorded on the leg (leg133.coveredBy)
// for check (c) to reconcile against the running binary's case roster; this function
// itself cannot see them, because it reads sources and sources say what is declared,
// not what is compiled.
//
// A leg that installs the listener and has no nail is red on its own line, which is
// the reading ticket 133's fifth shot needs on its first beat. The second beat is
// caught one level up, by the same census: once the install block is deleted the leg
// is still dispatched, and nothing drives it.
func (p *pkg133) coverageReds133(legs []*leg133) []string {
	var reds []string
	byLeg := map[string]*leg133{}
	shared := map[string]int{}
	for _, leg := range legs {
		byLeg[leg.key] = leg
	}
	for _, leg := range legs {
		for k := range leg.reached {
			shared[k]++
		}
	}
	claimedBy := map[string][]legCover133{}
	for _, c := range legCovers133 {
		if len(p.tests[c.test]) == 0 {
			// R-133-2: this used to say "not a test function" while the lookup
			// behind it accepted any declaration with a matching bare name,
			// including a helper and including a method. The bucket is tight now
			// (see collectTopLevel133), so this sentence and its predicate are the
			// same claim; the reading below names what was found instead so the
			// next reader does not have to guess whether the name is missing or
			// merely not a case.
			what := "declares no function of that name in this directory's test sources"
			if ds := p.helpers[c.test]; len(ds) > 0 {
				what = fmt.Sprintf("declares %s, which this file sorts with %s: a helper or a method, not a case Go's testing package runs", ds[0].site, "isRunnableCase133")
			}
			reds = append(reds, fmt.Sprintf("this gate claims leg %q is nailed by %q, which is not a test function in this directory's sources: %s. A nail has to name a top-level func TestXxx(t *testing.T); the gate's registry is a list of cases, and the ledger row it produces is a coverage claim. A renamed or deleted case is red here, where the registry is compared with the sources. A case that is in the sources but is not in this round's binary is red one check down for the `covered=test` bucket (runRosterReds135) and disclosed, not reddened, for a registry nail - see this file's header, third limit, and the disclosure line this run prints.", c.leg, c.test, what))
			continue
		}
		claimedBy[c.leg] = append(claimedBy[c.leg], c)
	}

	for _, leg := range legs {
		nails := claimedBy[leg.key]
		ruling := p.rulings[leg.key]
		// R-133-3: only a ruling that sits next to the leg's own code counts as
		// coverage. `ruling` above stays the first-seen site so the two readings
		// that do not credit it (the nail-plus-ruling contradiction, and the
		// orphan-ruling red at the bottom) keep seeing a marker in any file.
		adjacent := p.rulingAdjacentTo133(leg)
		driver, driverProof := "", ""
		if len(nails) == 0 {
			driver, driverProof = p.drivenBy133(leg, shared)
		}
		switch {
		case len(nails) > 0:
			var ns []string
			for _, n := range nails {
				ns = append(ns, n.test+" -> "+n.entry)
				if !p.hasDecl133(n.entry) {
					reds = append(reds, fmt.Sprintf("leg %q's nail %q claims it drives %q, which this directory does not declare.", leg.key, n.test, n.entry))
					continue
				}
				if !leg.reached[n.entry] {
					reds = append(reds, fmt.Sprintf("leg %q (%s) is claimed by nail %q on entry %q, and the dispatch no longer reaches that symbol. This is the dispatch hop itself: the case still compiles and still reads a disk, but the command it names is not wired to the code it claims to cover.",
						leg.key, leg.site, n.test, n.entry))
				}
			}
			leg.covered = "nail " + strings.Join(ns, ", ")
			if ruling != "" {
				reds = append(reds, fmt.Sprintf("leg %q carries a coverage ruling at %s AND a nail claim: say one thing. A ruling is how a leg admits nobody tests it, and that is how a missing nail gets waived, which is the reading this gate exists to refuse.", leg.key, ruling))
			}
		case leg.installs:
			leg.covered = "RED sink with no nail"
			reds = append(reds, fmt.Sprintf("leg %q (%s) reaches %s on this path: %s\nA listener installed on a dispatched leg has to have a nail in this gate's registry, or the block can be deleted in silence, which is the reading ticket 133 was filed for. Fix: write the case, then add {leg: %q, test: TestYourCase, entry: %q} to legCovers133.",
				leg.key, leg.site, sinkFunc133, leg.viaPath133(), leg.key, leg.entryName133()))
		case driver != "":
			leg.covered = "test " + driverProof
			// The bare name travels with the row so the run-roster check below
			// reconciles the same claim the ledger prints, in the same words.
			leg.coveredBy = []string{driver}
		case adjacent != "":
			leg.covered = "ruling " + adjacent
		default:
			leg.covered = "RED nothing"
			reds = append(reds, fmt.Sprintf("leg %q (%s) is dispatched by func main and covered by nothing: no nail in this gate's registry, no test case in this directory that drives a symbol belonging to this leg alone, and no %s sentence naming it.\nThat is the second beat of ticket 133's fifth shot: a leg whose install block was deleted has no install-based obligation left, and \"nobody anywhere verifies this command\" is the fact that survives it. Fix: drive it from a case, or write the ruling next to the code that owns the leg.",
				leg.key, leg.site, coverageRuling133))
		}
	}
	for key, site := range p.rulings {
		if _, ok := byLeg[key]; !ok {
			reds = append(reds, fmt.Sprintf("%s carries a coverage ruling for leg %q, which is not in the dispatch census: a ruling about a command nobody dispatches is prose, not a decision.", site, key))
		}
	}
	reds = append(reds, p.misplacedRulingReds133(legs)...)
	sort.Strings(reds)
	return reds
}

// drivenBy133 answers "which Test case in this directory reaches a symbol that
// belongs to this leg and to no other leg", and returns the bare name next to the
// sentence the ledger prints, so check (c) can reconcile the very claim the row makes.
// Evidence is restricted to plain functions, not methods: the bucket it reads is built
// by collectTopLevel133, which since R-133-1 admits only a top-level func TestXxx(t
// *testing.T) with no receiver and no results - that is, only declarations shaped the
// way the testing package runs them, which is not yet the same set as the cases this
// round's binary can start.
// The prefix filter below is a second, independent read of the same name, because
// a method name would be resolved here without receiver types and a coincidental
// `.stop()` in somebody else's case is not a claim about this leg. Test names are
// walked in sorted order so the row a reading quotes is the same row the next run
// quotes.
func (p *pkg133) drivenBy133(leg *leg133, shared map[string]int) (string, string) {
	var names []string
	for name := range p.tests {
		if strings.HasPrefix(name, "Test") {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	symbols := sortedKeys133(leg.reached)
	for _, name := range names {
		prod := p.testClosure133(p.tests[name])
		for _, k := range symbols {
			if shared[k] > 1 || strings.Contains(k, ".") {
				continue
			}
			if prod[k] {
				return name, name + " drives " + k
			}
		}
	}
	return "", ""
}

// testClosure133 is the production symbols one test decl reaches, through this
// directory's own test helpers.
func (p *pkg133) testClosure133(seeds []*decl133) map[string]bool {
	prod := map[string]bool{}
	seen := map[*decl133]bool{}
	var queue []*decl133
	for _, s := range seeds {
		if !seen[s] {
			seen[s] = true
			queue = append(queue, s)
		}
	}
	for len(queue) > 0 {
		d := queue[0]
		queue = queue[1:]
		if d.prod {
			prod[d.key] = true
		}
		for _, c := range p.callees133(d, true).decls {
			if !seen[c] {
				seen[c] = true
				queue = append(queue, c)
			}
		}
	}
	return prod
}

// ---------------------------------------------------------------------------
// check (c): the run roster - is a credited case name a case this round can start?
// ---------------------------------------------------------------------------

// rosterChildEnv135 marks the listing child this gate starts. Nothing runs a test body
// in `-test.list` mode, so the guard is belt-and-braces against a future where that
// stops being true: a child that re-enters this function says so instead of spawning a
// grandchild.
const rosterChildEnv135 = "WISP_135_RUN_ROSTER_CHILD"

// rosterReadTimeout135 caps the listing child. The measured cost of the call is 0.098 s
// for this package on this host; the cap exists so a wedged child turns into a red
// reading inside this test instead of a hung suite inside the parent's own -timeout.
const rosterReadTimeout135 = 2 * time.Minute

// rosterName135 matches one printed case name. The child's output also carries package
// init noise and its own verdict line, so a name counts only when the WHOLE line is an
// identifier in the testing package's own shape.
var rosterName135 = regexp.MustCompile(`^Test[A-Za-z0-9_]*$`)

// compiledRunRoster135 asks the test binary that is running this reading which
// top-level cases it can start, by re-executing itself with `-test.list '.*'`.
//
// WHY THIS, AND NOT MORE PARSING. Every other reading in this file is a go/parser walk
// over this directory, and a walk over sources answers "what is declared", never "what
// is in the build". That gap is ticket 133's R-133-9, family shot M-G: a correctly
// shaped `func TestXxx(t *testing.T)` in a file whose name carries the implicit
// GOOS=linux constraint is invisible to a GOOS=windows build, was booked by this gate as
// `covered=test ... drives ...`, and left the whole package green at 100 top-level
// passes. Closing it by teaching the parser to evaluate file-name suffixes and
// //go:build expressions was refused on this ticket (AC#8) and on ticket 133: a second
// instrument has to contribute a second, independent source of truth, or the audit is
// the audited object proving itself clean. The roster below comes from the binary's own
// generated test main, which is the exact list this round's `=== RUN` lines are drawn
// from, so it can say something the parser cannot infer.
//
// scripts/portable-tests.sh already reads the same object at the step level for the same
// reason: its -skip ledger "is re-verified against the compiled test binary for THIS
// platform via `go test -list`", and "Rename the test, delete it, or bury it behind a
// build tag, and the entry goes stale and this step goes red".
func compiledRunRoster135() (map[string]bool, int, error) {
	if os.Getenv(rosterChildEnv135) != "" {
		return nil, 0, fmt.Errorf("refused a nested run-roster read from inside the listing child")
	}
	exe, err := os.Executable()
	if err != nil {
		return nil, 0, fmt.Errorf("os.Executable: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), rosterReadTimeout135)
	defer cancel()
	cmd := exec.CommandContext(ctx, exe, "-test.list", ".*")
	cmd.Env = append(os.Environ(), rosterChildEnv135+"=1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, 0, fmt.Errorf("%s -test.list '.*': %w (child output opens: %s)", filepath.Base(exe), err, headRunes135(string(out), 300))
	}
	names := map[string]bool{}
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimRight(line, "\r")
		if rosterName135.MatchString(line) {
			names[line] = true
		}
	}
	if len(names) == 0 {
		return nil, 0, fmt.Errorf("%s -test.list '.*' listed no case at all (child output opens: %s)", filepath.Base(exe), headRunes135(string(out), 300))
	}
	return names, len(names), nil
}

// headRunes135 trims to n runes so a quoted child error can never split a multi-byte
// rune, which would put invalid UTF-8 into the test log.
func headRunes135(s string, n int) string {
	s = strings.TrimSpace(s)
	rs := []rune(s)
	if len(rs) > n {
		return string(rs[:n]) + "..."
	}
	return s
}

// flagValue135 returns the value a flag was given in this process's own command line, in
// either the -flag=value or the -flag value form, or "" when the flag is absent. The
// testing package's own flags are visible here, which is how this gate can say out loud
// whether the round it is reading is a whole-package round or a narrowed one.
func flagValue135(name string) string {
	eq := "-" + name + "="
	for i, a := range os.Args {
		if strings.HasPrefix(a, eq) {
			return strings.TrimPrefix(a, eq)
		}
		if a == "-"+name || a == "--"+name {
			if i+1 < len(os.Args) {
				return os.Args[i+1]
			}
		}
	}
	return ""
}

// runRosterReds135 is check (c): every case name this round's ledger credited through
// its `covered=test` form has to be a case the running binary can start, which is what
// makes a `=== RUN` line for it possible in this round.
func (p *pkg133) runRosterReds135(legs []*leg133, self string, roster map[string]bool, n int, rerr error) []string {
	if rerr != nil {
		return []string{fmt.Sprintf("the run roster this gate has to reconcile its coverage claims against could not be read: %v\n"+
			"Red, not skipped: with no roster every `covered=test` row in the ledger below is unverified, and \"the instrument could not look\" is not evidence that a leg is covered.", rerr)}
	}
	if !roster[self] {
		return []string{fmt.Sprintf("the run roster read lists %d startable cases and does not name %q, the case running this reading. The roster cannot be trusted, so none of the `covered=test` rows below were checked against it.", n, self)}
	}
	var reds []string
	for _, leg := range legs {
		for _, name := range leg.coveredBy {
			if roster[name] {
				continue
			}
			where := "declared nowhere in this directory's test sources"
			if ds := p.tests[name]; len(ds) > 0 {
				where = fmt.Sprintf("declared at %s", ds[0].site)
			}
			reds = append(reds, fmt.Sprintf("leg %q (%s) is booked in the ledger below as covered by the case %q, and this round's test binary has no such case: the roster read from the running binary lists %d startable cases and %q is not one of them, so no \"=== RUN   %s\" line exists in this run or can exist in it. The name is %s, which is the point - the declaration is in the sources and the sources are not the build.\n"+
				"That is ticket 133's R-133-9, family shot M-G: a func TestXxx(t *testing.T) this platform does not compile cannot run, so it reads no disk, books no record and holds no behaviour, and a row crediting it is a coverage claim with no witness. This is also the reading this file cannot get from its own parser, because the parser is what the plant satisfies. Fix: put the case in a file this build takes (the usual shapes are a `_linux_test.go`/`_windows_test.go` suffix and a //go:build line this platform fails), or take the claim out and let the leg be ruled out in words next to the code that owns it.",
				leg.key, leg.site, name, n, name, name, where))
		}
	}
	sort.Strings(reds)
	return reds
}

// runRosterDisclosure135 states, in one measured line, what the roster covered and what
// it could not see. It replaces the sentence this file used to guess ("on GOOS=windows
// the *_windows_test.go cases it names are parsed but not compiled"), which shot M-G
// measured as untrue in one direction and unenforced in the other.
func (p *pkg133) runRosterDisclosure135(legs []*leg133, roster map[string]bool, n int, rerr error) string {
	if rerr != nil {
		return fmt.Sprintf("run-roster disclosure: GOOS=%s, the roster could not be read (%v), so the coverage rows below carry no startable-case reading; that is red above.", goos133(), rerr)
	}
	driven := map[string]bool{}
	for _, leg := range legs {
		for _, name := range leg.coveredBy {
			driven[name] = true
		}
	}
	var missing []string
	for name := range driven {
		if !roster[name] {
			missing = append(missing, name)
		}
	}
	sort.Strings(missing)
	var nailsOut []string
	for _, c := range legCovers133 {
		if !roster[c.test] {
			nailsOut = append(nailsOut, c.test)
		}
	}
	sort.Strings(nailsOut)
	out := fmt.Sprintf("run-roster disclosure: GOOS=%s, %d startable cases read from this binary itself (`-test.list '.*'`); case names this round's ledger credited through `covered=test`: %d distinct, %d of them startable (one outside the roster is red above); this gate's registry: %d claims, %d of them outside this round's roster",
		goos133(), n, len(driven), len(driven)-len(missing), len(legCovers133), len(nailsOut))
	if len(nailsOut) > 0 {
		out += fmt.Sprintf(" - %s. A registry nail the build does not take is disclosed here and not reddened, because the leg it nails is installed and driven on the leg that does compile it: scripts/wisp-cli-tests.sh is the windows leg where these run (this file's header, third limit).", strings.Join(nailsOut, ", "))
	} else {
		out += " - every registry claim this gate makes is a case this round's binary can start."
	}
	if f := flagValue135("test.run"); f != "" && f != "*" {
		out += fmt.Sprintf(" This round was narrowed by -test.run=%s, so the roster proves these cases are startable here, not that each one printed === RUN in this process; a whole-package round (-count=1, no -test.run) is where the two are the same set, which is the shape scripts/wisp-cli-tests.sh runs.", f)
	} else {
		out += " This round is a whole-package round (no -test.run filter), so the roster is the set this run's === RUN lines are drawn from."
	}
	return out
}

// censusVsUsage133 reconciles the enumerated dispatch against the operator-facing
// command list in the usage block, in both directions: a second census over the same
// fact, read from a place the dispatch cannot edit. X4's --diag and X14's sfx131 are
// both missing from that text, and a leg that leaves the dispatch census the way
// ticket 131's X8 left it shows up as a hole here too.
func censusVsUsage133(pkg *pkg133, legs []*leg133) []string {
	usage := pkg.consts["usage"]
	if usage == "" {
		return []string{"the usage block is not a package-level string constant named `usage`, so this gate has no operator-facing list to reconcile against"}
	}
	docs := map[string]bool{}
	for _, m := range usageCommand133.FindAllStringSubmatch(usage, -1) {
		docs[m[1]] = true
	}
	if usageBare133.MatchString(usage) {
		docs[noArgsLeg133] = true
	}
	census := map[string]bool{}
	for _, leg := range legs {
		if leg.isDefault || isSynthetic133(leg.key) {
			continue
		}
		census[leg.key] = true
	}
	var reds []string
	for k := range census {
		if !docs[k] {
			reds = append(reds, fmt.Sprintf("leg %q is dispatched by func main and documented in no line of the usage block: an operator cannot find it, and this gate has no second reading of the census to check it against.", k))
		}
	}
	for k := range docs {
		if !census[k] {
			reds = append(reds, fmt.Sprintf("the usage block documents command %q, which func main does not dispatch: the promise and the dispatch disagree, which is how a leg goes missing from this census without anything noticing.", k))
		}
	}
	sort.Strings(reds)
	return reds
}

// collectRulings133 scans this directory's production sources for coverage rulings.
// It is a plain line scan on purpose: the ledger is built from the AST, and a second
// instrument that shares no code with the first is what makes a one-line loosening of
// either reading loud (ticket 129's fifth shape, aimed at this door).
func (p *pkg133) collectRulings133(dir string, entries []os.DirEntry) {
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		src, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			p.blind = append(p.blind, fmt.Sprintf("cannot read %s to scan it for coverage rulings", name))
			continue
		}
		for i, line := range strings.Split(string(src), "\n") {
			idx := strings.Index(line, coverageRuling133)
			if idx < 0 {
				continue
			}
			rest := strings.TrimSpace(line[idx+len(coverageRuling133):])
			fields := strings.Fields(rest)
			leg := ""
			if len(fields) > 0 {
				leg = fields[0]
			}
			if leg == "" {
				p.blind = append(p.blind, fmt.Sprintf("%s:%d carries %s with no leg name after it", name, i+1, coverageRuling133))
				continue
			}
			site := fmt.Sprintf("%s:%d", name, i+1)
			p.rulingAll = append(p.rulingAll, ruling133{leg: leg, site: site, file: name})
			if _, seen := p.rulings[leg]; !seen {
				p.rulings[leg] = site
			}
		}
	}
}

// rulingAdjacentTo133 returns the site of a ruling that names this leg AND sits in
// a file the leg owns, "" when the only rulings naming it sit elsewhere. "Owns" is
// deliberately narrow (R-133-3): the file the dispatch branch is written in, or a
// file declaring one of the symbols that branch calls directly. The red this
// reading serves says "write the ruling next to the code that owns the leg", and a
// marker that can be parked at the end of any production file in the directory is
// not next to anything - one line in doctor.go would back a coverage claim for a
// leg doctor has nothing to do with.
func (p *pkg133) rulingAdjacentTo133(leg *leg133) string {
	for _, r := range p.rulingAll {
		if r.leg != leg.key {
			continue
		}
		if p.legOwnsFile133(leg, r.file) {
			return r.site
		}
	}
	return ""
}

// legOwnsFile133 is the adjacency test above, spelled against the census row: the
// leg's own site and the files its entries are declared in.
func (p *pkg133) legOwnsFile133(leg *leg133, file string) bool {
	if strings.SplitN(leg.site, ":", 2)[0] == file {
		return true
	}
	for _, d := range p.allDecls133() {
		for _, e := range leg.entries {
			if d.key == e && d.file == file {
				return true
			}
		}
	}
	return false
}

// misplacedRulingReds133 reports a ruling that is not credited anywhere: naming a
// leg from a file that owns neither its dispatch nor any of its entries is prose
// in the wrong place, and it stays red even when the leg is covered by something
// else, so a planted backing line cannot sit in the tree quietly.
func (p *pkg133) misplacedRulingReds133(legs []*leg133) []string {
	byLeg := map[string]*leg133{}
	for _, leg := range legs {
		if _, ok := byLeg[leg.key]; !ok {
			byLeg[leg.key] = leg
		}
	}
	var reds []string
	for _, r := range p.rulingAll {
		leg, ok := byLeg[r.leg]
		if !ok {
			continue // the census already reports a ruling for a leg nobody dispatches
		}
		if p.legOwnsFile133(leg, r.file) {
			continue
		}
		owned := []string{strings.SplitN(leg.site, ":", 2)[0]}
		for _, e := range leg.entries {
			for _, d := range p.allDecls133() {
				if d.key == e {
					owned = append(owned, d.file)
				}
			}
		}
		reds = append(reds, fmt.Sprintf("%s carries %s for leg %q, but %s owns neither this leg's dispatch (%s) nor any symbol it calls (%s). A coverage ruling has to sit next to the code that owns the leg: as placed, one line at the end of any production file in this directory would back the claim that somebody verifies %q, which is the reading ticket 133's own acceptor registered as R-133-3.",
			r.site, coverageRuling133, r.leg, r.file, leg.site, strings.Join(dedupe133(owned), ", "), r.leg))
	}
	return reds
}

// ---------------------------------------------------------------------------
// small helpers
// ---------------------------------------------------------------------------

func (l *leg133) entryName133() string {
	if len(l.entries) == 0 {
		return "cmd" + l.key
	}
	return l.entries[0]
}

// viaPath133 renders the walk from the dispatch to the listener, which is the
// reading ticket 131's ledger states as a bare boolean.
func (l *leg133) viaPath133() string {
	if l.path == nil {
		return "unknown"
	}
	if p := l.path[sinkFunc133]; p != "" {
		return p
	}
	if l.handoff {
		return "qualified leaf " + handoffRoot133
	}
	return "not reached"
}

func isSynthetic133(key string) bool {
	return strings.HasPrefix(key, "unparsed-label@")
}

func legKeys133(legs []*leg133) []string {
	out := make([]string, 0, len(legs))
	for _, l := range legs {
		out = append(out, l.key)
	}
	sort.Strings(out)
	return out
}

func dedupe133(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		if seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

func sortedKeys133(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func goos133() string { return runtime.GOOS }
