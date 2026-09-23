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
//   - which nail a covered leg runs: this gate reconciles claims against compiled
//     test sources. Whether a nailed case asserts anything real is checked by
//     ticket 131's body check and by the per-leg mutations in ticket 131's face.
//
// COST: one go/parser pass over this directory's ~30 files per run, no
// subprocess, no filesystem beyond that. The measurements below are milliseconds.

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
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

	var lines []string
	for _, leg := range legs {
		lines = append(lines, fmt.Sprintf("  leg %-14s %-24s installs=%-5v handoff=%-5v covered=%-62s entries=%s",
			leg.key, leg.site, leg.installs, leg.handoff, leg.covered, strings.Join(leg.entries, "|")))
	}
	report := fmt.Sprintf("dispatch ledger, read out of func main's own branches at run time (%d legs, %d claims in this gate's registry):\n", len(legs), len(legCovers133)) +
		strings.Join(lines, "\n")
	if t.Failed() {
		t.Errorf("the dispatch ledger has a red row:\n%s", report)
	} else {
		t.Logf("%s", report)
	}
	t.Logf("blindness disclosure: this gate reconciles its nail claims against this directory's test sources, so on GOOS=%s the *_windows_test.go cases it names are parsed but not compiled; scripts/wisp-cli-tests.sh is the windows leg where they run.", goos133())
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

// leg133 is one dispatch branch of func main, with the readings the two checks
// above need.
type leg133 struct {
	key       string
	site      string
	isDefault bool
	entries   []string
	quals     map[string]bool
	reached   map[string]bool
	path      map[string]string
	installs  bool
	handoff   bool
	covered   string
}

// pkg133 is the parsed directory. Build tags are ignored and same-named
// declarations are unioned, so a leg inherits the other platform's shape before it
// loses its own; the approximation runs in the direction that asks for more
// coverage, never less.
type pkg133 struct {
	fset *token.FileSet

	decls   map[string][]*decl133 // production decls by bare name
	methods map[string][]*decl133 // production methods by bare name
	tests   map[string][]*decl133 // decls in *_test.go, by bare name

	globals   map[string]bool            // package-level var names
	varType   map[string]string          // package-level var -> struct or map type
	funcField map[string]bool            // "Type.field" -> the field is func-typed
	aliasOf   map[string]string          // "var" -> name it holds
	edges     map[string]string          // "var.field" / "var:key" -> name it holds
	unplaced  map[string]bool            // package-level func-typed var with no value
	imports   map[string]map[string]bool // file -> local import names
	consts    map[string]string          // package-level string constants

	mainBody *ast.BlockStmt
	argvName string
	blind    []string
	rulings  map[string]string // leg key -> file:line of its coverage ruling

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
		p.collectTopLevel133(file, name)
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

func (p *pkg133) collectTopLevel133(f *ast.File, name string) {
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
			case isTest:
				p.tests[fd.Name.Name] = append(p.tests[fd.Name.Name], decl)
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
// is still recognised as that case driving the leg.
func (p *pkg133) lookup133(name string, loose bool) []*decl133 {
	out := append([]*decl133{}, p.decls[name]...)
	out = append(out, p.methods[name]...)
	if loose {
		out = append(out, p.tests[name]...)
	}
	return out
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
			if fs := p.tests[fn.Name]; len(fs) > 0 {
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

	addLeg := func(key, site string, n ast.Node, isDefault bool, extraReds ...string) {
		reds = append(reds, extraReds...)
		ranges = append(ranges, [2]int{int(n.Pos()), int(n.End())})
		leg := &leg133{key: key, site: site, isDefault: isDefault}
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
			key, kind := p.classifyCond133(s.Cond)
			site := p.site133(s.Pos())
			switch kind {
			case condNoArgs133:
				addLeg(noArgsLeg133, site, s.Body, false)
			case condLiteral133:
				addLeg(key, site, s.Body, false)
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
					addLeg("default", p.site133(cc.Pos()), cc, true)
					continue
				}
				key := ""
				var labelReds []string
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
						name := "unparsed-label@" + p.site133(cc.Pos()) + ":" + expr
						if key == "" {
							key = name
						}
						labelReds = append(labelReds, fmt.Sprintf("%s: the case label %s is not a string literal%s. A leg is identified by a literal command name; a label that has to be resolved through a constant is a leg that can leave the census without saying so, which is ticket 131's second shape. This row carries it as %q instead of folding it into default.",
							p.site133(l.Pos()), expr, resolved, name))
						continue
					}
					if key == "" {
						if v, err := strconv.Unquote(bl.Value); err == nil {
							key = v
						}
					}
				}
				if key == "" {
					key = "unparsed-label@" + p.site133(cc.Pos())
				}
				addLeg(key, p.site133(cc.Pos()), cc, false, labelReds...)
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

type condKind133 int

const (
	condUnknown133 condKind133 = iota
	condNoArgs133
	condLiteral133
)

// classifyCond133 reads the two argv conditions this gate dispatches on: the
// `len(argv) == 0` GUI branch, and any branch comparing an argv slot to a literal
// (X4's `if len(args) > 0 && args[0] == "--diag"` is the second form).
func (p *pkg133) classifyCond133(e ast.Expr) (string, condKind133) {
	key := ""
	found := false
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
					key, found = v, true
				}
			}
			return true
		}
		if call, ok := bin.X.(*ast.CallExpr); ok && exprIdent(call.Fun) == "len" && len(call.Args) == 1 && bin.Op == token.EQL {
			if id, ok := call.Args[0].(*ast.Ident); ok && id.Name == p.argvName {
				if bl, ok := bin.Y.(*ast.BasicLit); ok && bl.Value == "0" && !found {
					key, found = noArgsLeg133, true
				}
			}
		}
		return true
	})
	if !found {
		return "", condUnknown133
	}
	if key == noArgsLeg133 {
		return key, condNoArgs133
	}
	return key, condLiteral133
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
// file's own registry, a test case that drives a symbol belonging to that leg alone,
// or a WISP-LEG-COVERAGE-RULING sentence in a production file naming the leg.
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
			reds = append(reds, fmt.Sprintf("this gate claims leg %q is nailed by %q, which is not a test function in this directory's sources. A renamed or build-tag-hidden case has to be red somewhere, and here is where.", c.leg, c.test))
			continue
		}
		claimedBy[c.leg] = append(claimedBy[c.leg], c)
	}

	for _, leg := range legs {
		nails := claimedBy[leg.key]
		ruling := p.rulings[leg.key]
		driver := ""
		if len(nails) == 0 {
			driver = p.drivenBy133(leg, shared)
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
			leg.covered = "test " + driver
		case ruling != "":
			leg.covered = "ruling " + ruling
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
	sort.Strings(reds)
	return reds
}

// drivenBy133 answers "which Test case in this directory reaches a symbol that
// belongs to this leg and to no other leg". Evidence is restricted to plain
// functions, not methods: a method name is resolved here without receiver types,
// and a coincidental `.stop()` in somebody else's case is not a claim about this
// leg. Test names are walked in sorted order so the row a reading quotes is the
// same row the next run quotes.
func (p *pkg133) drivenBy133(leg *leg133, shared map[string]int) string {
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
				return name + " drives " + k
			}
		}
	}
	return ""
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
			if _, seen := p.rulings[leg]; !seen {
				p.rulings[leg] = fmt.Sprintf("%s:%d", name, i+1)
			}
		}
	}
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
