//go:build windows

package main

// Ticket 253, leg 253-r1 - the lexical/AST HALF of AC#1 that no ruler on disk covers:
// "the door name the dispatch site binds is the same door name the page-side forwarding
// script addresses, and both come from one reviewed constant" - the binding literal
// tied to the roster, checked by READING THIS PACKAGE'S OWN SOURCE ON DISK.
//
// WHY THIS CELL AND NOT THE CAPABILITY ONE. 票 253 AC#1 (:16) asked to move the
// reachability judgement from a word ruler to a capability ruler. The orchestrator's
// 2026-10-08 ruling (.scratch/wisp/probes/253 panel narrowing, ledger A717 section 5,
// census .scratch/wisp/probes/ruler-dedup-1/ruler-dedup-1.md row R5) measured that the
// capability half is ALREADY on disk: cmd/wisp/panel_transport_35r1_test.go:195 calls the
// real (*PanelManager).installPanelTransport (:203) and reddens when no page->door hook is
// registered. Rebuilding that would be the duplicate labour issues/README.md:80 forbids.
// What the census found genuinely absent is the narrow word/AST shape:
//
//	grep -rn "panelDispatchBinding" --include=*.go cmd internal
//
// hits production at exactly three sites - the constant (cmd/wisp/panel_host_windows.go
// :80), the page-side forwarding script built from it (:792), and the Bind call in
// installPanelTransport (:802) - and NO ruler asserts those three move together. Rename
// the constant's VALUE and all three shift in lockstep (that is the design the comment at
// :777-779 asks for); rename it at ONE site only and the page ends up calling a door Go
// never bound - the exact dead-door shape 票 253 AC#2 later refused to let "a name with no
// answer" count as registered. That lockstep is this file's whole subject.
//
// WHAT THE ROSTER IS. doorRoster253r1 below is the reviewed closed set of page<->host
// inbound transport doors this package binds inside installPanelTransport - today exactly
// one name, "wispDispatch". The ruler judges four things against it, all read off disk:
//
//	(1) the constant named doorBindingConst253r1 is declared, with a non-empty string value;
//	(2) installPanelTransport exists EXACTLY ONCE in this package's production sources (its
//	    own doc comment calls it "the ONE place the page<->host inbound transport is wired");
//	(3) the Bind call(s) inside it bind exactly the rostered names (bidirectional: an
//	    unrostered bound name is red, a rostered name nobody binds is red, a roster with no
//	    live Bind behind it cannot rot quietly - the same shape as the Locked-suffix roster in
//	    cmd/wisp/panel_locked_naming_33r11_windows_test.go);
//	(4) the script handed to w.Init inside installPanelTransport is built (via fmt.Sprintf)
//	    from the SAME constant, so the JS door the page posts to and the Go door the transport
//	    binds are the one reviewed name, not two names free to drift.
//
// WHY A SOURCE CENSUS AND NOT ONE MORE CALL. The capability ruler already calls the real
// function on a fake control; what it CANNOT see is a rename at one of the three sites,
// because in that world the fake control never exercises the real page script at all. This
// census reads the bytes the page script and the binding are made of, so its reach is
// general and it stays honest about the transient probe door (w.Bind("wispProbeRT") at
// :852) which is deliberately NOT part of the transport roster.
//
// THIS RULER READS THE DISK, SO ITS GREEN IS NOT PROVEN BY -overlay. go test -overlay
// swaps the compiler's eyes, never os.ReadFile's (ledger A713; the ruling the 255r6 ruler
// restates at cmd/wisp/panel_geometry_255r6_range_windows_test.go:10-13). The self-proofs
// for the production subject therefore plant the wrong shape ON DISK and restore it with
// git show HEAD:<file> after hashing both sides (.scratch/wisp/probes/253/r1/). The
// positive-control fixtures in this file parse source strings this file OWNS: that branch
// is overlay-visible and only demonstrates that the matching logic bites a bad name and
// spares an honest one - it is never offered as proof that the disk census reads anything.

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// doorBindingConst253r1 is the package-level constant whose value names the transport
// door. The census requires the page script to be built from THIS identifier, so renaming
// the identifier is a reviewed act that must move the roster and this name together.
const doorBindingConst253r1 = "panelDispatchBinding"

// doorTransportFunc253r1 is the one function that wires the page<->host transport. The
// census finds it by name across the whole package, so it does not hardcode the file it
// currently lives in (cmd/wisp/panel_host_windows.go today).
const doorTransportFunc253r1 = "installPanelTransport"

// doorRoster253r1 is the reviewed closed set of transport door names bound inside
// installPanelTransport. Adding a name here is the reviewed act; leaving one behind when
// the Bind is gone is caught by the same case (bidirectional, like lockedNamingRoster33r11).
var doorRoster253r1 = map[string]string{
	"wispDispatch": "installPanelTransport's one inbound door: the page's postMessage is folded into window.<this name> by panelPostMessageForwardInit, and its Go callback feeds (*panel.ComposerDispatch).Handle via dispatchRaw",
}

// bindingDoor253r1 is one recognised door name bound inside the transport, with the site
// that bound it, so a red can name file:line instead of only a count.
type bindingDoor253r1 struct {
	name string
	site string
}

// transportSite253r1 is the census result for one installPanelTransport declaration.
type transportSite253r1 struct {
	file       string
	line       int
	bindCalls  int
	doors      []bindingDoor253r1
	initDoorId string // constant name the w.Init script was fmt.Sprintf'd from ("" if none)
	initSites  []string
}

// pkgSymbols253r1 is the package-level string table the census resolves identifiers
// through: constants' string values and, for every package-level var, the identifiers fed
// to any fmt.Sprintf that builds it (the page script's door reference lives there).
type pkgSymbols253r1 struct {
	consts     map[string]string   // const name -> string literal value
	varSprintf map[string][]string // var name -> identifier names inside its fmt.Sprintf args
}

func newPkgSymbols253r1() *pkgSymbols253r1 {
	return &pkgSymbols253r1{consts: map[string]string{}, varSprintf: map[string][]string{}}
}

// absorbFile records the package-level declarations of one parsed file. It reads only
// top-level const/var GenDecls, because the door constant and the forwarding script are
// both package-level by design (:80 and :781).
func (p *pkgSymbols253r1) absorbFile(file *ast.File) {
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gen.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			switch gen.Tok {
			case token.CONST:
				for i, n := range vs.Names {
					if i < len(vs.Values) {
						if lit, ok := vs.Values[i].(*ast.BasicLit); ok && lit.Kind == token.STRING {
							if v, err := strconv.Unquote(lit.Value); err == nil {
								p.consts[n.Name] = v
							}
						}
					}
				}
			case token.VAR:
				for _, n := range vs.Names {
					var idents []string
					for _, v := range vs.Values {
						ast.Inspect(v, func(x ast.Node) bool {
							call, ok := x.(*ast.CallExpr)
							if !ok {
								return true
							}
							if !isFmtSprintf253r1(call) {
								return true
							}
							for _, a := range call.Args {
								if id, ok := a.(*ast.Ident); ok {
									idents = append(idents, id.Name)
								}
							}
							return true
						})
					}
					if len(idents) > 0 {
						p.varSprintf[n.Name] = append(p.varSprintf[n.Name], idents...)
					}
				}
			}
		}
	}
}

// isFmtSprintf253r1 recognises fmt.Sprintf(...) as a selector call (an aliased fmt import
// is still fmt's Sprintf here: the selector name is what carries the meaning).
func isFmtSprintf253r1(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Sprintf" {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	return ok && pkg.Name == "fmt"
}

// censusTransport253r1 records, for one installPanelTransport body, every Bind first
// argument and every Init first argument. A door is "recognised" only when its first
// argument is a non-empty string literal or an identifier that resolves through the
// package constant table to a non-empty string; a Bind to a non-constant expression or an
// empty name is counted as a call but yields no door, which is what makes an empty roster
// a failure rather than a quiet pass.
func censusTransport253r1(fset *token.FileSet, file *ast.File, name string, syms *pkgSymbols253r1) []transportSite253r1 {
	var sites []transportSite253r1
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != doorTransportFunc253r1 || fn.Body == nil {
			continue
		}
		site := transportSite253r1{file: name, line: fset.Position(fn.Pos()).Line}
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			pos := fmt.Sprintf("%s:%d", name, fset.Position(call.Pos()).Line)
			switch sel.Sel.Name {
			case "Bind":
				site.bindCalls++
				if len(call.Args) == 0 {
					return true
				}
				door, _ := resolveDoorExpr253r1(call.Args[0], syms)
				if door != "" {
					site.doors = append(site.doors, bindingDoor253r1{name: door, site: pos})
				}
			case "Init":
				site.initSites = append(site.initSites, pos)
				if len(call.Args) == 0 {
					return true
				}
				if id, ok := call.Args[0].(*ast.Ident); ok {
					for _, fed := range syms.varSprintf[id.Name] {
						if fed == doorBindingConst253r1 {
							site.initDoorId = id.Name
						}
					}
				}
			}
			return true
		})
		sites = append(sites, site)
	}
	return sites
}

// resolveDoorExpr253r1 turns a Bind's first argument into the door name it denotes: a
// string literal yields itself, an identifier yields the package constant's value. An
// empty or unresolvable expression yields "" (a bind call that names no reviewed door).
func resolveDoorExpr253r1(e ast.Expr, syms *pkgSymbols253r1) (string, bool) {
	switch v := e.(type) {
	case *ast.BasicLit:
		if v.Kind == token.STRING {
			if s, err := strconv.Unquote(v.Value); err == nil {
				return s, s != ""
			}
		}
		return "", false
	case *ast.Ident:
		if val, ok := syms.consts[v.Name]; ok {
			return val, val != ""
		}
		return "", false
	default:
		return "", false
	}
}

// packageProductionSources253r1 is the census subject, shape-for-shape from
// packageProductionSources255r6 and packageSourceFiles33r11: ReadDir of the package
// directory (the directory `go test` runs in - THIS RULER READS THE WORKING TREE, not a
// git snapshot), files only, .go only, _test.go excluded, sorted. An empty list is a
// failure, never a clean pass: a census that read no files proves nothing about any door.
func packageProductionSources253r1(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("253r1 census could not read its own package directory %s: %v", dir, err)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)
	if len(names) == 0 {
		t.Fatalf("253r1 census found zero .go files in %s: the ruler has no subject", dir)
	}
	return names
}

// parseProductionFiles253r1 parses every production source into ONE shared FileSet (a
// position read through a different FileSet is a number nobody can trust) and builds the
// package symbol table from them before any transport is judged, because the transport's
// Bind and Init reference package-level constants and vars declared in the same package.
// A file that cannot be read or parsed reddens through t.Fatalf, never a silent skip.
func parseProductionFiles253r1(t *testing.T, dir string) (*token.FileSet, []*ast.File, *pkgSymbols253r1, []string) {
	t.Helper()
	fset := token.NewFileSet()
	syms := newPkgSymbols253r1()
	var files []*ast.File
	names := packageProductionSources253r1(t, dir)
	for _, name := range names {
		path := filepath.Join(dir, name)
		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("253r1 census could not read %s: %v", path, err)
		}
		file, err := parser.ParseFile(fset, name, src, parser.SkipObjectResolution)
		if err != nil {
			t.Fatalf("253r1 census could not parse %s - a source the ruler cannot read is not an exempt source: %v", name, err)
		}
		syms.absorbFile(file)
		files = append(files, file)
	}
	return fset, files, syms, names
}

// TestTransportDoorBindingMatchesRoster253r1 is the disk census: the transport's bound door
// names, the door constant, the page script's door reference and the reviewed roster must
// all name the same door. This is the case the mutation self-proofs redden on disk.
func TestTransportDoorBindingMatchesRoster253r1(t *testing.T) {
	dir := "."
	fset, files, syms, names := parseProductionFiles253r1(t, dir)

	constValue, constFound := syms.consts[doorBindingConst253r1]
	if !constFound {
		t.Fatalf("package main production sources declare no string constant %q: the transport door has no reviewed name to bind, so this ruler has no subject (dir=%s, %d files read)",
			doorBindingConst253r1, dir, len(names))
	}
	if constValue == "" {
		t.Fatalf("constant %q resolves to the empty string: a door named by nothing is bound by nobody and addressed by no page script",
			doorBindingConst253r1)
	}

	var sites []transportSite253r1
	for _, f := range files {
		sites = append(sites, censusTransport253r1(fset, f, fset.File(f.Pos()).Name(), syms)...)
	}
	sort.Slice(sites, func(a, b int) bool {
		if sites[a].file == sites[b].file {
			return sites[a].line < sites[b].line
		}
		return sites[a].file < sites[b].file
	})

	if len(sites) == 0 {
		t.Fatalf("no function named %s exists in this package's production sources (read %d files): the one documented transport-wiring host is gone, so nothing here ties a door name to a Bind",
			doorTransportFunc253r1, len(names))
	}
	if len(sites) > 1 {
		var where []string
		for _, s := range sites {
			where = append(where, fmt.Sprintf("%s:%d", s.file, s.line))
		}
		t.Fatalf("found %d functions named %s across the package (%v); its own doc comment calls it the ONE place the transport is wired, and two of them is a split the roster cannot cover",
			len(sites), doorTransportFunc253r1, where)
	}

	site := sites[0]
	t.Logf("253-r1 census: door constant %q = %q; %s at %s:%d made %d Bind call(s) and %d Init call(s); roster holds %d name(s) over %d production file(s)",
		doorBindingConst253r1, constValue, doorTransportFunc253r1, site.file, site.line, site.bindCalls, len(site.initSites), len(doorRoster253r1), len(names))

	// (3a) an empty roster of recognised doors is a failure, not a pass: installPanelTransport
	//     must actually bind at least one name the census can read.
	if len(site.doors) == 0 {
		t.Fatalf("%s (%s:%d) issued %d Bind call(s) but bound zero readable door names: the transport roster is empty, and a page posting to any door lands on nothing",
			doorTransportFunc253r1, site.file, site.line, site.bindCalls)
	}

	boundNames := map[string]string{} // door name -> first site that bound it
	for _, d := range site.doors {
		boundNames[d.name] = d.site
	}

	// (3b) bidirectional roster. A bound name nobody reviewed is red; a reviewed name nobody
	//     bound is red (that is the dead-door 票 253 refuses to count as registered).
	for _, d := range site.doors {
		if _, reviewed := doorRoster253r1[d.name]; !reviewed {
			t.Errorf("AC#1 RED: %s binds transport door %q at %s, which is not in doorRoster253r1 - a new inbound page-addressable door entered the tree without being reviewed here",
				doorTransportFunc253r1, d.name, d.site)
		}
	}
	for reviewed := range doorRoster253r1 {
		if _, bound := boundNames[reviewed]; !bound {
			t.Errorf("AC#1 RED: doorRoster253r1 reviews door %q, but %s (%s:%d) binds no such name (bound: %v) - a rostered door with no Bind behind it is the dead door 票 253 AC#2 refused to call registration",
				reviewed, doorTransportFunc253r1, site.file, site.line, keysOf253r1(boundNames))
		}
	}

	// (1) the door constant's own value must be a reviewed name, so renaming the value
	//     without the roster is caught, not laundered through the identifier.
	if _, reviewed := doorRoster253r1[constValue]; !reviewed {
		t.Errorf("AC#1 RED: constant %q resolves to %q (%s), which is not in doorRoster253r1 - the transport door's reviewed name and its declared value have come apart",
			doorBindingConst253r1, constValue, site.file)
	}

	// (4) the page script handed to w.Init must be built from the SAME door constant, so the
	//     JS door and the Go door cannot be renamed independently.
	if site.initDoorId == "" {
		t.Errorf("AC#1 RED: %s calls w.Init(%v) with no script built (via fmt.Sprintf) from %q - the page-side door the forwarding hook addresses is not the constant the transport binds, so the two can drift apart and leave the page posting at an unbound door",
			doorTransportFunc253r1, site.initSites, doorBindingConst253r1)
	} else {
		t.Logf("253-r1 lockstep: w.Init(%s) at %s derives its window.<door> reference from %q, the same constant the Bind uses",
			site.initDoorId, strings.Join(site.initSites, ", "), doorBindingConst253r1)
	}
}

func keysOf253r1(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// A goodBindingFixture253r1 is a self-contained source the census logic is asked to bless:
// the constant, a page script built from it by fmt.Sprintf, and an installPanelTransport
// that binds that constant and passes that script to w.Init.
const goodBindingFixture253r1 = `package fixture

import "fmt"

const panelDispatchBinding = "wispDispatch"

var panelPostMessageForwardInit = fmt.Sprintf("window.%[1]s(message)", panelDispatchBinding)

func installPanelTransport(w pageTransport) error {
	if err := w.Bind(panelDispatchBinding, func(raw string) string { return "" }); err != nil {
		return err
	}
	w.Init(panelPostMessageForwardInit)
	return nil
}
`

// A driftedBindFixture253r1 is the M1 shape parsed from a string: the constant and the page
// script still name wispDispatch, but the transport binds a literal that matches nothing.
// The door the page posts to has no Go-side answer behind it.
const driftedBindFixture253r1 = `package fixture

import "fmt"

const panelDispatchBinding = "wispDispatch"

var panelPostMessageForwardInit = fmt.Sprintf("window.%[1]s(message)", panelDispatchBinding)

func installPanelTransport(w pageTransport) error {
	if err := w.Bind("wispStaleDoor", func(raw string) string { return "" }); err != nil {
		return err
	}
	w.Init(panelPostMessageForwardInit)
	return nil
}
`

// emptyBindFixture253r1 is the M2 shape parsed from a string: the transport issues a Bind
// but binds the empty name, so no reviewed door survives the census.
const emptyBindFixture253r1 = `package fixture

import "fmt"

const panelDispatchBinding = "wispDispatch"

var panelPostMessageForwardInit = fmt.Sprintf("window.%[1]s(message)", panelDispatchBinding)

func installPanelTransport(w pageTransport) error {
	if err := w.Bind("", func(raw string) string { return "" }); err != nil {
		return err
	}
	w.Init(panelPostMessageForwardInit)
	return nil
}
`

// untiedInitFixture253r1 is the fourth face: the transport binds the reviewed door, but the
// page script is built from a DIFFERENT literal than the constant, so the JS half and the
// Go half are free to drift apart.
const untiedInitFixture253r1 = `package fixture

import "fmt"

const panelDispatchBinding = "wispDispatch"

var panelPostMessageForwardInit = fmt.Sprintf("window.%s(message)", "wispOther")

func installPanelTransport(w pageTransport) error {
	if err := w.Bind(panelDispatchBinding, func(raw string) string { return "" }); err != nil {
		return err
	}
	w.Init(panelPostMessageForwardInit)
	return nil
}
`

// runFixtureCensus253r1 runs the same symbol-table + transport census the disk case uses,
// over one source string this file owns. It returns the resolved constant value and the
// single transport site, or fails if the fixture is not shaped like the transport.
func runFixtureCensus253r1(t *testing.T, name, src string) (string, transportSite253r1) {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, name, src, parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("fixture %s did not parse, so the census was never asked to bite: %v", name, err)
	}
	syms := newPkgSymbols253r1()
	syms.absorbFile(file)
	constValue, found := syms.consts[doorBindingConst253r1]
	if !found {
		t.Fatalf("fixture %s declares no %q constant", name, doorBindingConst253r1)
	}
	sites := censusTransport253r1(fset, file, name, syms)
	if len(sites) != 1 {
		t.Fatalf("fixture %s declared %d %s funcs, want exactly 1", name, len(sites), doorTransportFunc253r1)
	}
	return constValue, sites[0]
}

// rosterRedsFor253r1 is the roster judgement the disk case applies, factored so the
// positive controls assert on it directly instead of re-deciding here.
func rosterRedsFor253r1(constValue string, site transportSite253r1) []string {
	var reds []string
	bound := map[string]bool{}
	for _, d := range site.doors {
		bound[d.name] = true
		if _, ok := doorRoster253r1[d.name]; !ok {
			reds = append(reds, "unrostered bound door "+d.name)
		}
	}
	for reviewed := range doorRoster253r1 {
		if !bound[reviewed] {
			reds = append(reds, "rostered but unbound door "+reviewed)
		}
	}
	if _, ok := doorRoster253r1[constValue]; !ok {
		reds = append(reds, "constant value "+constValue+" not reviewed")
	}
	if site.initDoorId == "" {
		reds = append(reds, "w.Init script not tied to the door constant")
	}
	sort.Strings(reds)
	return reds
}

// TestBindingRosterBitesItsOwnFixtures253r1 is the matching-logic positive control (overlay
// -visible branch, self-owned strings): the honest binding is clean, and each of the three
// wrong shapes is reported. This proves the census would flag them if it read them; it does
// NOT claim the disk census read anything (that proof is the on-disk plants).
func TestBindingRosterBitesItsOwnFixtures253r1(t *testing.T) {
	cases := []struct {
		name    string
		src     string
		wantRed int
	}{
		{"good", goodBindingFixture253r1, 0},
		{"drifted-bind", driftedBindFixture253r1, 2}, // unrostered bound door + rostered unbound door (the constant value itself stays reviewed)
		{"empty-bind", emptyBindFixture253r1, 1},     // only the rostered-unbound door; the empty Bind yields no readable door name
		{"untied-init", untiedInitFixture253r1, 1},   // only the Init-not-tied-to-constant face; the roster itself still agrees
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			constValue, site := runFixtureCensus253r1(t, tc.name+"253r1.go", tc.src)
			reds := rosterRedsFor253r1(constValue, site)
			if tc.name == "empty-bind" && len(site.doors) == 0 {
				t.Logf("empty-bind: census bound zero readable doors (%d Bind call(s)) - the disk case reddens this as the empty-roster Fatalf", site.bindCalls)
			}
			if len(reds) != tc.wantRed {
				t.Fatalf("fixture %q reported %d red(s) %v, want %d", tc.name, len(reds), reds, tc.wantRed)
			}
			t.Logf("fixture %q: const=%q bound=%v initTied=%q reds=%v", tc.name, constValue, keysOf253r1(doorNamesOf253r1(site.doors)), site.initDoorId, reds)
		})
	}
}

func doorNamesOf253r1(doors []bindingDoor253r1) map[string]string {
	m := make(map[string]string, len(doors))
	for _, d := range doors {
		m[d.name] = d.site
	}
	return m
}
