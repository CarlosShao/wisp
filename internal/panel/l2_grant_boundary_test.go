package panel

// ---------------------------------------------------------------------------
// The panel-side L2 "allow" ban, nailed on the GO BOUNDARY.
//
// AGENTS.md §1.2 lists as an absolutely-forbidden code shape:
//
//	"由面板侧来源的 L2「允许」" - a web panel control may never be the thing
//	that grants a high-risk approval. The panel may display state and INITIATE a
//	request; the verdict comes from the native side (D33/F2, ruling R20).
//
// WHY THIS FILE EXISTS (the incident it is a nail against, re-derived at this
// anchor rather than retold): commit 53a1359 (2026-09-25 07:41) deleted 7 lines
// from frontend/src/components/l2-approval-card.tsx. In the parent commit
// (git show 53a1359^) those lines were a live button:
//
//	:161-167  a bg-accent Button, onClick={() => send("grant")}
//	:96-100   const send = (outcome: ApprovalOutcome) => requestApprovalResolution(...)
//	   ->      frontend/src/lib/panel.ts:169-186
//	   ->      bridge.postMessage(JSON.stringify({method:"panel.approval.request", correlationId, outcome}))
//
// The only thing that made it inert was Go-side: panel.approval.request is not an
// answered inbound method (facet 1), and no answered method can carry an outcome
// (facet 2). Deleting the button removed one DRAWING of the shape and nothing
// that would catch the next one - and frontend/src/lib/panel.ts:50 still exports
// ApprovalOutcome with "grant" in it, so the wire vocabulary is still alive on the
// other side of the boundary.
//
// WHY THE STRONGEST EXISTING INSTRUMENT MISSED IT:
// TestTheRendererHoldsExactlyOneDoorToTheHost (composer_test.go:502) pins which
// file may talk to the host, that routes are literals, and that every "panel.*"
// literal is one Go answers - and composer_test.go:417 puts
// "panel.approval.request" in its ALLOWED set. The violating call site and the
// legitimate ones differ only by which member of ApprovalOutcome is passed at
// argument position (send("grant") at :163 vs send("refuse") at :112 and :155).
// A ruler that reads which door was used cannot see what was carried through it.
//
// WHAT THIS FILE PINS (four facets, all Go-side, all green at this anchor):
//
//	(1) the set of inbound methods Go actually answers. It is derived from this
//	    package's own route guard rather than copied into a test, it goes LOUD
//	    when a case label cannot be resolved, and it is cross-checked against
//	    knownComposerMethod in BOTH directions - plus the wider oracle, which
//	    takes every route-shaped name written anywhere in the package and asks the
//	    running guard about it. No answered name is approval-shaped. Facet 1 has a
//	    second, STANDING half that reads no source at all:
//	    TestRealGuardRefusesEveryAssemblableApprovalRouteName assembles its
//	    candidate names out of grantRouteWords and asks the running guard about
//	    each one, so an approval-shaped route that Go answers is a red test even
//	    when the package never spells its name (F-R2-3);
//	(2) the envelope every answered method decodes into cannot BIND a key whose
//	    name carries an approval verdict, recursively through embedded and nested
//	    types. "The envelope" is found by locating this package's own JSON decode
//	    calls, not by recognising a route field, and it is checked three ways: an
//	    AST scan that can name file:line, reflection over the production types,
//	    and encoding/json itself via DisallowUnknownFields;
//	(3) behaviour: the historical wire shape is refused by ParseComposerRequest,
//	    and for every answered route an envelope that adds outcome:"grant" parses
//	    into a value that no longer holds it;
//	(4) teeth: the same AST instrument, pointed at a COPY of this package with a
//	    grant-carrying envelope, an answered approval route, a renamed route key,
//	    an unreadable case label and a second route chain planted, must go red or
//	    say why it could not read the tree.
//
// NOTHING IS WIRED TODAY - THIS PIN PROTECTS NOTHING YET. No production code
// calls ParseComposerRequest, so no request from the page reaches these routes at
// all: at this anchor `grep -rn ParseComposerRequest --include=*.go .` yields the
// definition (bridge.go:77) plus test call sites, and every other hit is a
// comment - including cmd/wisp/run.go:225, which says so itself ("Nothing calls it
// yet - the WebView2 \"event -> ParseComposerRequest\" hop does not exist in this
// tree"), and `grep -c webview go.mod` is 0, so the module cannot even receive a
// postMessage today. A green run here therefore does NOT mean "the panel cannot
// grant"; it means "this boundary cannot be WIRED to grant without a red test".
// Read the two rounds in docs/evidence/s1/ before treating any of it as covered.
//
// WHAT THIS FILE DOES *NOT* COVER - read this before citing it as coverage:
// the ban's wording is wider than this instrument. A violation drawn entirely in
// JSX (a button whose onClick passes "grant" into a legitimate
// requestApprovalResolution call) is invisible here, because this reads Go only.
// facet 4 measures that blindness rather than pretending it away: the verbatim
// historical button lines are planted into the same tree and the instrument
// reports nothing about them. Registered as zero-instrument coverage on the UI
// side in docs/evidence/s1/panel-l2-grant-nail-r1.md.
//
// Two holes stay open after the r2 fixes, and both are rules about what is
// WRITTEN rather than what can be MADE to happen: a route name assembled at
// runtime (out of config, or off the wire) appears in no literal pool and no
// static scan reaches it, and a verdict field whose spelling is not in
// grantFieldWords ("proceed", "yes", "ok") is a name this file has never heard
// of. The vocabulary is deliberately narrow - see its own comment - so these are
// judgement calls, not oversights, and widening them is a slice card rather than
// a quiet edit here.
//
// One sentence of qualification on the first hole, because it is now part-static:
// TestRealGuardRefusesEveryAssemblableApprovalRouteName asks the RUNNING guard
// about every name grantRoutePrefixes x grantRouteWords x grantRouteSuffixes can
// build (528 candidates at this anchor), and it reads no source, so an answered
// name inside that grid is red even when the package never spells it. The hole
// that stays open is the grid's outside: another namespace, another vocabulary,
// or a route that only exists once config is read.
//
// The size of that grid is itself nailed now (F-ACC-1, docs/evidence/s1/panel-
// l2-grant-nail-fix-r3-accept-r1.md §1.1): the helper multiplies the three lists
// and refuses to report a verdict when the number of names it asked is not that
// product, so emptying any one of them, or dropping one of the two plural forms,
// goes red inside the sweep instead of quietly printing asked=0 over a grid that
// stopped existing.
//
// Two more nails, both from docs/evidence/s1/panel-l2-grant-nail-fix-r4-accept-r1.md
// §10/§11. (F-R4-1) The ask-the-guard step is now one function, sweepAssembledNames,
// and the positive control reaches the guard through it rather than from a line of
// its own - so short-circuiting the sweep's ask or accumulate step reddens the
// control, and "hits=[]" can no longer mean "never asked". (F-R4-2) A factor list
// can no longer be narrowed quietly: every member of grantRoutePrefixes and
// grantRouteSuffixes has a witness row that the sweep must still be able to
// assemble, and wantGrantRoute* writes the three counts down as numbers, because
// asked= always equals the lists as they now stand. What stays open is the
// three-piece edit inside this file - member, witness row and written count moved
// together - which this round's own rig measured green (X-list+row+pin); it is the
// same family grantRouteWordWitnesses has always shared, and the backstop is the
// human roster diff, not a ruler here.
//
// ONE CONSTRAINT THIS FILE PUTS ON PRODUCTION CODE, DECLARED HERE BECAUSE
// NOTHING ELSE WRITES IT DOWN (F-R2-1, docs/evidence/s1/panel-l2-grant-nail-
// accept-r2.md §2.4 + §7.4): the inbound half of this instrument starts from the
// JSON decode calls it finds in this package, so a SECOND decode destination in
// internal/panel has to be a struct in this package. A decode into a
// map[string]any or a json.RawMessage - both perfectly legal Go, neither
// enumerable by a key scan - makes four of this file's tests abort with "a
// JSON decode destination ... cannot be enumerated": facet 1, facet 2, the
// three-way key test, and facet 4's pristine-snapshot precheck. That direction is
// deliberate: the file refuses to call the boundary clean over bytes it cannot
// read, and the message carries its own way out. That way out is ONE PATH IN
// TWO STEPS, not a choice between two (F-ACC-3, docs/evidence/s1/panel-l2-grant-
// nail-fix-r3-accept-r1.md §5.2; the wording here used to say "or", and the
// second step is what the first one runs into): step one, route the bytes
// through a struct in this package, which is the only thing that makes the
// destination enumerable by a key scan; step two, register that type in
// inboundTypeRegistry, which is what lets the reflection half of the three-way
// key test see it. Step one ALONE leaves this file red - one test,
// TestJSONKeyDerivationAgreesWithEncodingJSON, "decode destination <name> has no
// reflection twin in inboundTypeRegistry" - so a second decode destination that
// enters the package has to enter the registry in the same commit, and the
// registered type is then judged by both instruments at once. Loosening the
// scan is not a test edit either, it is the same slice-card decision as the
// vocabulary above.
//
// NOTHING here keys on frontend/src/lib/panel.ts:50's ApprovalOutcome union: that
// union still contains "grant" by the frontend session's explicit choice, so an
// enum-keyed check would be red today and would not be a Go-boundary nail.
// ---------------------------------------------------------------------------

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// grantFieldWords are name fragments that turn a JSON key into an approval
// verdict carrier. The test is on the SHAPE OF THE NAME (lowercased, separators
// removed, then substring-matched), never on rendered text and never on a value:
// a panel that can SEND a verdict is the violation, whichever string it sends.
//
// Deliberately narrow: every entry here would give an inbound envelope the power
// to say "allow". Words that name an approval-adjacent concept without carrying a
// verdict (confirm, accept, intent) are out, and
// TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes pins that the narrowness
// still catches the shapes that matter - a vocabulary too wide to be true gets
// widened away by the next person who hits a false red.
var grantFieldWords = []string{
	"approve", "approval", "autoapprove", "grant", "granted",
	"allow", "allowed", "allowonce", "permit", "permitted",
	"ratify", "authorize", "authorised", "decision", "decide",
	"verdict", "outcome", "bypass", "override",
}

// grantRouteWords is the same idea for a ROUTE name: a method whose own name
// states an approval decision is a panel-side allow door.
var grantRouteWords = []string{
	"approve", "approval", "grant", "allow", "permit",
	"ratify", "authorize", "authorised", "decide", "decision", "verdict",
}

// grantRouteWordWitnesses is the roster behind F-ACC-2 (docs/evidence/s1/panel-
// l2-grant-nail-fix-r3-accept-r1.md §1.2): one route name per member of
// grantRouteWords, so that deleting a word from that list cannot quietly remove
// the 48 candidate names that word contributes to the sweep.
//
// It is not "a table asserting the table contains itself". Each row is judged by
// two things that live outside it - the REAL predicate over the REAL vocabulary
// (carriers must equal exactly the set of grantRouteWords entries that flag
// route today, so both shrinking and growing that list breaks a row that is now
// describing a different vocabulary), and the RUNNING production guard, which
// must refuse every name here: that refusal is the live evidence this file claims
// for each word, and it is a fact about bridge.go, not about this list.
//
// carriers is the set of grantRouteWords entries that flag route, written out
// rather than derived: no word in this vocabulary is a substring of another
// ("approve" is not inside "approval", "decide" is not inside "decision", and
// authorize/authorised diverge before the ending), so each row names exactly one
// word today. Recording the set anyway is the point - if the vocabulary grows a
// member that also matches one of these names, the row that did not ask for it
// goes red and says so.
//
// route names here are deliberate, not templated: "panel.review.ratify" is the
// spelling acceptance r3 planted as its M16, and "panel.approval.request" is the
// spelling the removed 53a1359^ button emitted. Widening grantRouteWords means
// adding a row in the same edit, and TestGrantVocabularyIsNotSatisfiedByTheReal
// Envelopes goes red when the two lists disagree in either direction.
var grantRouteWordWitnesses = []struct {
	word     string
	route    string
	carriers []string
}{
	{"approve", "panel.mode.approve", []string{"approve"}},
	{"approval", "panel.approval.request", []string{"approval"}},
	{"grant", "panel.grant", []string{"grant"}},
	{"allow", "panel.l2.allow", []string{"allow"}},
	{"permit", "panel.permits", []string{"permit"}},
	{"ratify", "panel.review.ratify", []string{"ratify"}},
	{"authorize", "panel.authorize.request", []string{"authorize"}},
	{"authorised", "panel.authorised.request", []string{"authorised"}},
	{"decide", "panel.decide", []string{"decide"}},
	{"decision", "panel.decision", []string{"decision"}},
	{"verdict", "panel.verdict", []string{"verdict"}},
}

// carriesGrantWord reports whether name, normalised, contains a verdict word.
func carriesGrantWord(name string, vocab []string) bool {
	norm := strings.ToLower(strings.NewReplacer("_", "", "-", "", " ", "", ".", "").Replace(name))
	for _, w := range vocab {
		if strings.Contains(norm, w) {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// facet 2/3 helper: every JSON key a decoder can bind into a type, recursive
// through embedded, nested, pointer, slice and array types.
// ---------------------------------------------------------------------------

// bindableKey is one wire-addressable field: the JSON name plus a Go path to
// point at when it goes wrong.
type bindableKey struct {
	Key    string
	GoPath string
}

func bindableKeysOf(typ reflect.Type, depth int, prefix string, seen map[reflect.Type]bool) []bindableKey {
	if typ == nil || depth > 6 {
		return nil
	}
	if seen[typ] {
		return nil
	}
	seen[typ] = true
	for typ.Kind() == reflect.Ptr || typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array || typ.Kind() == reflect.Map {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil
	}
	var out []bindableKey
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if !f.IsExported() {
			continue
		}
		tag := f.Tag.Get("json")
		jsonName := ""
		if tag != "" {
			jsonName = strings.Split(tag, ",")[0]
		}
		// The name-part rule, shared with jsonOr and with the AST half: "-" is
		// skipped, an empty name is the Go field name. The first version of this
		// helper compared the WHOLE tag to "-", so `json:"-,omitempty"` fell
		// through as a key named "-" - harmless for this ban (no verdict word
		// lives in that spelling) but it is a second instrument that disagreed.
		if jsonName == "-" {
			continue
		}
		path := prefix + "." + f.Name
		// An embedded type with no JSON name is not addressable by its own name:
		// encoding/json promotes its fields, so only recurse (and do not mint a
		// phantom key that would make this rule about Go type names).
		if f.Anonymous && jsonName == "" {
			out = append(out, bindableKeysOf(f.Type, depth+1, path, seen)...)
			continue
		}
		key := f.Name
		if jsonName != "" {
			key = jsonName
		}
		out = append(out, bindableKey{Key: key, GoPath: path})
		if underlyingStruct(f.Type) {
			out = append(out, bindableKeysOf(f.Type, depth+1, path, seen)...)
		}
	}
	return out
}

func underlyingStruct(typ reflect.Type) bool {
	for typ.Kind() == reflect.Ptr || typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array || typ.Kind() == reflect.Map {
		typ = typ.Elem()
	}
	return typ.Kind() == reflect.Struct
}

func grantCarryingKeys(keys []bindableKey) []bindableKey {
	var hits []bindableKey
	for _, k := range keys {
		if carriesGrantWord(k.Key, grantFieldWords) {
			hits = append(hits, k)
		}
	}
	sort.Slice(hits, func(i, j int) bool { return hits[i].GoPath < hits[j].GoPath })
	return hits
}

func sortedKeys(keys []bindableKey) []string {
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		out = append(out, k.Key)
	}
	sort.Strings(out)
	return out
}

// ---------------------------------------------------------------------------
// The AST instrument: the same rule, but able to name file:line and runnable
// against a snapshot the reflection path cannot look at.
// ---------------------------------------------------------------------------

type boundaryFinding struct {
	File string
	Line int
	What string
}

func (b boundaryFinding) String() string {
	return b.File + ":" + strconv.Itoa(b.Line) + ": " + b.What
}

// goSourceFiles returns the non-test .go files under dir, sorted.
func goSourceFiles(dir string) ([]string, error) {
	var out []string
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
			return nil
		}
		out = append(out, p)
		return nil
	})
	sort.Strings(out)
	return out, err
}

// astField is one struct field as the AST sees it.
type astField struct {
	JSONKey  string
	GoName   string
	Line     int
	Text     string
	TypeName string // same-package struct type it embeds or points at, "" if none
	NoTag    bool
	// Anonymous is the AST's own record of an embedded field, which is how
	// astTopKeys knows not to mint a key for a promoted type.
	Anonymous bool
}

// structShape is one declared struct type.
type structShape struct {
	name   string
	file   string
	line   int
	fields []astField
}

// unresolvedLabel is one case label of the route guard that this instrument
// cannot turn into a string. It used to be dropped on the floor, which made the
// enumeration silently narrower than the guard (acceptance r1, F-1: a route
// reaching through a package-level var or a concatenation is invisible to a scan
// that only reads literals and package-level constants). Recording it lets every
// caller choose to go loud, and lets the teeth plant name the shape.
type unresolvedLabel struct {
	File string
	Line int
	Expr string
}

func (u unresolvedLabel) String() string {
	return u.File + ":" + strconv.Itoa(u.Line) + ": " + u.Expr
}

// boundaryPackage is one parsed directory.
type boundaryPackage struct {
	structs   map[string]structShape
	consts    map[string]string
	answered  map[string]bool
	literals  []routeLiteral
	dropped   []unresolvedLabel
	guardFile string
	guardLine int
	guardSeen bool

	// pkgVars holds a package-level var's declared type, so a decode destination
	// written at file scope still resolves.
	pkgVars map[string]ast.Expr
	// decodes is every JSON decode this package performs and what it lands in;
	// inboundSeeds and decodeProblems are derived from it after the parse.
	decodes       []decodeSite
	inboundSeeds  map[string]bool
	decodeProblem string
}

// decodeSite is one JSON decode call in the package: the expression it decodes
// into, and the same-package struct type that expression resolves to ("" when it
// resolves to nothing this instrument can enumerate).
type decodeSite struct {
	File     string
	Line     int
	DstExpr  string
	TypeName string
	TypeText string
}

// parsedFile keeps one file's tree so the decode pass can run after every
// package-level var in the directory is known.
type parsedFile struct {
	af          *ast.File
	fset        *token.FileSet
	rel         string
	importsJSON bool
}

// importsEncodingJSON reports whether a file names encoding/json, which is what
// makes a bare ".Decode(&x)" call worth treating as a JSON decode at all.
func importsEncodingJSON(af *ast.File) bool {
	for _, imp := range af.Imports {
		if unquoteGo(imp.Path.Value) == "encoding/json" {
			return true
		}
	}
	return false
}

// collectPackageVars records package-level var types across one file. Only the
// single-name, explicitly-typed spelling matters here, because that is how a
// hidden module-level destination gets written.
func collectPackageVars(af *ast.File, pkg *boundaryPackage) {
	for _, decl := range af.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, ident := range vs.Names {
				if vs.Type != nil {
					pkg.pkgVars[ident.Name] = vs.Type
					continue
				}
				if vi := indexOf(vs.Names, ident); vi < len(vs.Values) {
					pkg.pkgVars[ident.Name] = valueExprType(vs.Values[vi])
				}
			}
		}
	}
}

func indexOf(names []*ast.Ident, want *ast.Ident) int {
	for i, n := range names {
		if n == want {
			return i
		}
	}
	return -1
}

// valueExprType narrows an initialiser to the type expression it implies. It
// understands new(T), &T{} and a bare ident; anything else yields nil, which
// makes the site loud rather than guessed at.
func valueExprType(expr ast.Expr) ast.Expr {
	switch t := expr.(type) {
	case *ast.UnaryExpr:
		if t.Op == token.AND {
			return valueExprType(t.X)
		}
	case *ast.CompositeLit:
		return t.Type
	case *ast.CallExpr:
		if id, ok := t.Fun.(*ast.Ident); ok && id.Name == "new" && len(t.Args) == 1 {
			return t.Args[0]
		}
	}
	return nil
}

// collectDecodeSites walks every function body in one file and records the
// destination of each JSON decode. The seed criterion for "what is an inbound
// envelope" used to be "a struct that binds a method key", which a renamed route
// field walks straight past (F-4 in the acceptance: {Cmd, Outcome} with a real
// json.Unmarshal, five tests green). A decode call cannot be walked past that
// way - it is the thing the ban is about.
func collectDecodeSites(pf parsedFile, pkg *boundaryPackage) {
	add := func(fn *ast.FuncType, body *ast.BlockStmt) {
		if body == nil {
			return
		}
		scope := map[string]ast.Expr{}
		conflicts := map[string]bool{}
		record := func(name string, typ ast.Expr) {
			if typ == nil {
				return
			}
			if prev, seen := scope[name]; seen && exprText(pf.fset, prev) != exprText(pf.fset, typ) {
				conflicts[name] = true
			}
			scope[name] = typ
		}
		if fn != nil && fn.Params != nil {
			for _, f := range fn.Params.List {
				for _, n := range f.Names {
					record(n.Name, f.Type)
				}
			}
		}
		ast.Inspect(body, func(n ast.Node) bool {
			switch st := n.(type) {
			case *ast.DeclStmt:
				gd, ok := st.Decl.(*ast.GenDecl)
				if !ok || gd.Tok != token.VAR {
					return true
				}
				for _, spec := range gd.Specs {
					vs, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}
					for i, ident := range vs.Names {
						if vs.Type != nil {
							record(ident.Name, vs.Type)
						} else if i < len(vs.Values) {
							record(ident.Name, valueExprType(vs.Values[i]))
						}
					}
				}
			case *ast.AssignStmt:
				if st.Tok != token.DEFINE {
					return true
				}
				for i, lhs := range st.Lhs {
					id, ok := lhs.(*ast.Ident)
					if !ok || id.Name == "_" || i >= len(st.Rhs) {
						continue
					}
					if len(st.Rhs) == 1 {
						record(id.Name, valueExprType(st.Rhs[0]))
						continue
					}
					record(id.Name, valueExprType(st.Rhs[i]))
				}
			}
			return true
		})
		ast.Inspect(body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			var dst ast.Expr
			switch {
			case sel.Sel.Name == "Unmarshal" && identName(sel.X) == "json" && len(call.Args) >= 2:
				dst = call.Args[1]
			case sel.Sel.Name == "Decode" && pf.importsJSON && len(call.Args) >= 1 &&
				strings.Contains(exprText(pf.fset, sel.X), "json"):
				dst = call.Args[0]
			default:
				return true
			}
			site := decodeSite{
				File: pf.rel, Line: pf.fset.Position(call.Pos()).Line,
				DstExpr: exprText(pf.fset, dst),
			}
			site.TypeName, site.TypeText = resolveDestination(dst, scope, pkg.pkgVars, pf.fset, conflicts)
			pkg.decodes = append(pkg.decodes, site)
			return true
		})
	}
	for _, decl := range pf.af.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok {
			add(fd.Type, fd.Body)
		}
	}
	ast.Inspect(pf.af, func(n ast.Node) bool {
		if fl, ok := n.(*ast.FuncLit); ok {
			add(fl.Type, fl.Body)
		}
		return true
	})
}

func identName(e ast.Expr) string {
	if id, ok := e.(*ast.Ident); ok {
		return id.Name
	}
	return ""
}

func exprText(fset *token.FileSet, e ast.Expr) string {
	if e == nil {
		return ""
	}
	return renderExpr(fset, e)
}

// resolveDestination turns the second argument of a decode call into the type it
// lands in. Pointer-of, composite-literal-of and a var or parameter name are all
// read; a name that is not declared anywhere in the file's package scope, and a
// name rebound to two different types, both come back unresolvable.
func resolveDestination(dst ast.Expr, scope, pkgVars map[string]ast.Expr, fset *token.FileSet, conflicts map[string]bool) (string, string) {
	if u, ok := dst.(*ast.UnaryExpr); ok && u.Op == token.AND {
		dst = u.X
	}
	var typ ast.Expr
	switch t := dst.(type) {
	case *ast.CompositeLit:
		typ = t.Type
	case *ast.Ident:
		if conflicts[t.Name] {
			return "", "conflicting local declarations of " + t.Name
		}
		typ = scope[t.Name]
		if typ == nil {
			typ = pkgVars[t.Name]
		}
		if typ == nil {
			return "", "no declaration found for " + t.Name
		}
	default:
		return "", renderExpr(fset, dst)
	}
	return structNameOf(typ), renderExpr(fset, typ)
}

// classifyDecodes is run once every struct in the package is known: it decides
// which decode destinations are enumerable (a same-package struct), and which are
// a hole in this instrument (everything else, including map[string]any, which can
// bind ANY key a page sends).
func (pkg *boundaryPackage) classifyDecodes() {
	pkg.inboundSeeds = map[string]bool{}
	var holes []string
	for _, d := range pkg.decodes {
		typeText := d.TypeText
		if typeText == "" {
			typeText = "(unknown)"
		} else {
			typeText = strconv.Quote(typeText)
		}
		if d.TypeName == "" {
			holes = append(holes, d.File+":"+strconv.Itoa(d.Line)+
				": decodes into "+d.DstExpr+" (type "+typeText+"), which is not a same-package struct: "+
				"a destination this instrument cannot enumerate keys for is not a clean boundary")
			continue
		}
		if _, ok := pkg.structs[d.TypeName]; !ok {
			holes = append(holes, d.File+":"+strconv.Itoa(d.Line)+
				": decodes into "+strconv.Quote(d.TypeName)+", a type from another package or not a struct: "+
				"its keys are not visible to this scan")
			continue
		}
		pkg.inboundSeeds[d.TypeName] = true
	}
	if len(holes) > 0 {
		sort.Strings(holes)
		pkg.decodeProblem = "a JSON decode destination in " + strconv.Itoa(len(holes)) +
			" place(s) cannot be enumerated: " + strings.Join(holes, "; ") +
			". Judge it in a test that can see that type, or route the bytes through a same-package struct - do not let this file report the boundary clean"
	}
}

func quoteOrEmpty(s string) string {
	if s == "" {
		return "(unknown)"
	}
	return strconv.Quote(s)
}

// parseBoundaryPackage reads a directory's struct types, string constants and
// inbound route guard.
func parseBoundaryPackage(dir string) (boundaryPackage, error) {
	pkg := boundaryPackage{
		structs:  map[string]structShape{},
		consts:   map[string]string{},
		answered: map[string]bool{},
		pkgVars:  map[string]ast.Expr{},
	}
	files, err := goSourceFiles(dir)
	if err != nil {
		return pkg, err
	}
	var parsed []parsedFile
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			return pkg, err
		}
		lines := strings.Split(string(data), "\n")
		fset := token.NewFileSet()
		af, err := parser.ParseFile(fset, path, data, 0)
		if err != nil {
			return pkg, err
		}
		rel := filepath.ToSlash(filepath.Base(path))

		collectRouteLiterals(af, fset, rel, &pkg)
		collectPackageVars(af, &pkg)
		parsed = append(parsed, parsedFile{
			af: af, fset: fset, rel: rel, importsJSON: importsEncodingJSON(af),
		})

		ast.Inspect(af, func(n ast.Node) bool {
			switch decl := n.(type) {
			case *ast.GenDecl:
				if decl.Tok != token.CONST {
					return true
				}
				for _, spec := range decl.Specs {
					vs, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}
					for vi, ident := range vs.Names {
						if vi >= len(vs.Values) {
							continue
						}
						if lit, ok := vs.Values[vi].(*ast.BasicLit); ok && lit.Kind == token.STRING {
							pkg.consts[ident.Name] = unquoteGo(lit.Value)
						}
					}
				}
				return true

			case *ast.FuncDecl:
				// The guard is what decides "answered by Go". If it is ever
				// renamed or replaced, this instrument must go loud rather than
				// report a clean boundary it has stopped looking at - that is
				// scanGrantBoundary's job, via guardSeen.
				if decl.Name.Name != "knownComposerMethod" {
					return true
				}
				pkg.guardSeen = true
				if pkg.guardFile == "" {
					pkg.guardFile = rel
					pkg.guardLine = fset.Position(decl.Pos()).Line
				}
				names, dropped := routeLabelsInFunc(decl, pkg.consts, fset, rel)
				for _, name := range names {
					pkg.answered[name] = true
				}
				pkg.dropped = append(pkg.dropped, dropped...)
				return true

			case *ast.TypeSpec:
				st, ok := decl.Type.(*ast.StructType)
				if !ok {
					return true
				}
				sh := structShape{name: decl.Name.Name, file: rel, line: fset.Position(decl.Pos()).Line}
				for _, fld := range st.Fields.List {
					raw := ""
					if fld.Tag != nil {
						raw = unquoteGo(fld.Tag.Value)
					}
					jsonName, tagged := "", false
					if raw != "" {
						if v, ok := reflect.StructTag(raw).Lookup("json"); ok {
							tagged = true
							jsonName = strings.Split(v, ",")[0]
						}
					}
					line := fset.Position(fld.Pos()).Line
					text := ""
					if line-1 < len(lines) {
						text = strings.TrimSpace(lines[line-1])
					}
					typeName := structNameOf(fld.Type)
					anon := len(fld.Names) == 0
					for _, nm := range fld.Names {
						sh.fields = append(sh.fields, astField{
							JSONKey: jsonOr(nm.Name, jsonName, tagged), GoName: nm.Name,
							Line: line, Text: text, TypeName: typeName, NoTag: !tagged,
						})
					}
					if anon {
						sh.fields = append(sh.fields, astField{
							JSONKey: jsonName, GoName: "embedded:" + typeName,
							Line: line, Text: text, TypeName: typeName, NoTag: !tagged,
							Anonymous: true,
						})
					}
				}
				pkg.structs[decl.Name.Name] = sh
				return true
			}
			return true
		})
	}
	// Second pass: a decode destination may be a package-level var declared in
	// another file, and classifying it needs every struct in the directory, so
	// neither half can run inside the loop above.
	for _, pf := range parsed {
		collectDecodeSites(pf, &pkg)
	}
	pkg.classifyDecodes()
	return pkg, nil
}

// jsonOr is the AST half's JSON-key rule, and it has to be the SAME rule
// encoding/json runs, because the reflection half in this file gets it right by
// asking the real decoder. encoding/json resolves a tag whose name part is empty
// (`json:",omitempty"`) to the Go field name, not to "no key". The first version
// of this helper returned "" for that spelling, which made the AST half blind to
// exactly the shape acceptance r1 planted as M7/M8: the field bound "Outcome" on
// the wire, reflection went red, the AST reported zero findings, and the two
// instruments that claim to state one rule in two ways disagreed in the
// under-reporting direction (their F-5).
func jsonOr(goName, jsonName string, tagged bool) string {
	if tagged && jsonName != "" {
		return jsonName
	}
	return goName
}

func unquoteGo(s string) string {
	if r, err := strconv.Unquote(s); err == nil {
		return r
	}
	return strings.Trim(s, "`\"")
}

// structNameOf resolves a field type to a same-package struct type name.
func structNameOf(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return structNameOf(t.X)
	case *ast.ArrayType:
		return structNameOf(t.Elt)
	}
	return ""
}

// routeLabelsInFunc reads a switch guard and resolves each case label to a
// string, through the package's own constants. Anything it cannot resolve is
// RETURNED, not dropped: the caller decides how loud to be, and the teeth plant
// needs to read the drop list without aborting its own subtest.
func routeLabelsInFunc(fn *ast.FuncDecl, consts map[string]string, fset *token.FileSet, rel string) ([]string, []unresolvedLabel) {
	var names []string
	var dropped []unresolvedLabel
	ast.Inspect(fn, func(n ast.Node) bool {
		sw, ok := n.(*ast.SwitchStmt)
		if !ok {
			return true
		}
		for _, stmt := range sw.Body.List {
			cs, ok := stmt.(*ast.CaseClause)
			if !ok {
				continue
			}
			for _, e := range cs.List {
				if s, ok := routeNameOf(e, consts); ok {
					names = append(names, s)
					continue
				}
				dropped = append(dropped, unresolvedLabel{
					File: rel, Line: fset.Position(e.Pos()).Line, Expr: renderExpr(fset, e),
				})
			}
		}
		return true
	})
	return names, dropped
}

// renderExpr prints one expression back the way it was written, so a loud
// failure can name the thing it could not read instead of a line number alone.
func renderExpr(fset *token.FileSet, e ast.Expr) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, e); err != nil || buf.Len() == 0 {
		// No fmt import for one fallback: name the shape that defeated the
		// resolver rather than hiding it behind an empty string.
		return "<unprintable " + reflect.TypeOf(e).String() + ">"
	}
	return buf.String()
}

// routeNameOf resolves one case label to the route string it stands for. It
// understands exactly two shapes - a string literal, and an identifier naming a
// package-level string constant - because those are the two the guard uses
// today. Everything else is a drop, and every drop is reported.
func routeNameOf(e ast.Expr, consts map[string]string) (string, bool) {
	switch t := e.(type) {
	case *ast.BasicLit:
		if t.Kind == token.STRING {
			return unquoteGo(t.Value), true
		}
	case *ast.Ident:
		if v, ok := consts[t.Name]; ok {
			return v, true
		}
	}
	return "", false
}

// joinUnresolved renders a drop list for a failure message.
func joinUnresolved(items []unresolvedLabel) string {
	out := make([]string, 0, len(items))
	for _, it := range items {
		out = append(out, "\n  "+it.String())
	}
	sort.Strings(out)
	return strings.Join(out, "")
}

// instrumentBlindnessProblem names why this tree's answered-set or inbound-set
// cannot be trusted, or "" when both can. It is a function rather than a pile of
// t.Fatal calls so the teeth plants in facet 4 can assert the *message* exists
// without aborting their own subtest.
//
// The two families are the two ways an enumeration gets narrower than the thing
// it reads: a case label that cannot be resolved (acceptance r1 F-1, where a
// route reached through a package-level var or a concatenation was dropped
// without a word), and a JSON decode whose destination cannot be named (F-4's
// descendants - an inbound type this scan cannot list the keys of is an inbound
// type it must not certify).
func instrumentBlindnessProblem(dir string, pkg boundaryPackage) string {
	if !pkg.guardSeen {
		return "no inbound route guard (knownComposerMethod) under " + dir +
			": the instrument can no longer tell which methods Go answers, so a clean report here would be silence, not safety"
	}
	if len(pkg.dropped) > 0 {
		return "the inbound route guard under " + dir + " has " + strconv.Itoa(len(pkg.dropped)) +
			" case label(s) this instrument cannot resolve to a string: " + joinUnresolved(pkg.dropped) +
			". This scan reads string literals and package-level string constants; a route reaching through a var, a call, or a concatenation is a door facet 1 cannot see, and it will not be reported as clean. Give the label a const, or teach routeNameOf the expression - and if it was deliberate, say so in the test that owns it (D33/F2, AGENTS.md ban #6)"
	}
	if len(pkg.answered) == 0 {
		return "the route guard under " + dir + " resolved to 0 answered routes - the enumeration is broken, not clean"
	}
	if len(pkg.structs) == 0 {
		return "no struct types under " + dir + " - a scan that parses nothing reports nothing"
	}
	if len(pkg.decodes) == 0 {
		return "no JSON decode call under " + dir + ": the inbound-envelope scan has no destination to start from, so it judges nothing and reports that as a clean boundary"
	}
	if pkg.decodeProblem != "" {
		return pkg.decodeProblem + " (under " + dir + ")"
	}
	return ""
}

// requireReadableInstrument turns instrumentBlindnessProblem into an abort.
func requireReadableInstrument(t *testing.T, dir string, pkg boundaryPackage) {
	t.Helper()
	if problem := instrumentBlindnessProblem(dir, pkg); problem != "" {
		t.Fatalf("%s", problem)
	}
}

// inboundEnvelopes returns the struct types this package treats as an INBOUND
// envelope, found by structure and not by a name typed in here. The seed is
// "a type some JSON decode in this package lands in", which replaced the older
// "a struct that binds a method key" - a route field renamed to "cmd" walks past
// the second criterion and straight into the ban (F-4 in the acceptance, plant G
// in facet 4). Types nested or embedded inside a destination are pulled in too,
// because that is where a smuggled verdict gets parked.
func inboundEnvelopes(pkg boundaryPackage) map[string]bool {
	inbound := map[string]bool{}
	for name := range pkg.inboundSeeds {
		if _, ok := pkg.structs[name]; ok {
			inbound[name] = true
		}
	}
	for changed := true; changed; {
		changed = false
		for name := range inbound {
			for _, f := range pkg.structs[name].fields {
				if f.TypeName == "" {
					continue
				}
				if _, ok := pkg.structs[f.TypeName]; !ok {
					continue // another package's type: not this envelope's face
				}
				if !inbound[f.TypeName] {
					inbound[f.TypeName] = true
					changed = true
				}
			}
		}
	}
	return inbound
}

// routeLiteral is one string literal in this package's Go sources that is
// spelled like a route, and the line it was written on.
type routeLiteral struct {
	File  string
	Line  int
	Value string
}

// routeShapedName reports whether s is spelled like a route name: two or more
// dot-separated segments of letters, digits and underscores (a leading hyphen in
// a segment is not allowed). Struct tags, printf verbs, paths, mimes and file
// names with a directory all fail this shape, which is what keeps the pool from
// drowning in noise.
//
// Being route-shaped is only the filter. What decides anything is whether the
// real guard answers the name - see poolJudgedByRealGuard - so a stray match
// here costs nothing.
func routeShapedName(s string) bool {
	segments := strings.Split(s, ".")
	if len(segments) < 2 {
		return false
	}
	for _, seg := range segments {
		if seg == "" {
			return false
		}
		for i, r := range seg {
			switch {
			case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_':
			case r == '-' && i > 0:
			default:
				return false
			}
		}
	}
	return true
}

// collectRouteLiterals records every route-shaped string literal in one parsed
// file, INCLUDING the ones sitting inside functions this instrument never reads.
// That is the whole point: the guard's own switch is not the only place a route
// name is written down, and acceptance r1's M4 planted its door in a second
// function whose case list the old enumeration had no business seeing.
func collectRouteLiterals(af *ast.File, fset *token.FileSet, rel string, pkg *boundaryPackage) {
	ast.Inspect(af, func(n ast.Node) bool {
		lit, ok := n.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}
		v := unquoteGo(lit.Value)
		if !routeShapedName(v) {
			return true
		}
		pkg.literals = append(pkg.literals, routeLiteral{
			File: rel, Line: fset.Position(lit.Pos()).Line, Value: v,
		})
		return true
	})
}

// routeNamePool is every candidate route name this package spells out: each
// route-shaped literal, plus each package-level string constant with the same
// shape (a guard may answer a name that only ever appears in a const). It maps a
// name to the places it was found, so a failure can point.
func routeNamePool(pkg boundaryPackage) map[string][]string {
	pool := map[string][]string{}
	for _, l := range pkg.literals {
		key := l.File + ":" + strconv.Itoa(l.Line)
		pool[l.Value] = append(pool[l.Value], key)
	}
	for name, val := range pkg.consts {
		if routeShapedName(val) {
			pool[val] = append(pool[val], "const "+name)
		}
	}
	return pool
}

// poolJudgedByRealGuard is the oracle. It asks the RUNNING guard, not the AST,
// about every name the package spells out, and reports two things:
//
//	answeredElsewhere - the guard answers a name the derived enumeration never
//	  saw. This is F-2: a second switch/if chain that knownComposerMethod hands
//	  its decision to is invisible to a scan pointed at that one function, and it
//	  is exactly the shape a handler registry takes.
//	answeredGrantDoor - the guard answers a name whose own spelling is an
//	  approval decision. This is the ban itself (D33/F2, AGENTS.md ban #6).
//
// The direction matters and is written down so nobody re-reads this as a
// two-way proof: names come from what the package WRITES. A route assembled at
// runtime from a config value or the network is not in the pool and cannot be,
// and no static ruler reaches it either.
func poolJudgedByRealGuard(pool map[string][]string, answered map[string]bool) (answeredElsewhere, answeredGrantDoor []string) {
	for name, sites := range pool {
		if !knownComposerMethod(name) {
			continue
		}
		sort.Strings(sites)
		where := strings.Join(sites, ", ")
		if !answered[name] {
			answeredElsewhere = append(answeredElsewhere, strconv.Quote(name)+" (written at "+where+")")
		}
		if carriesGrantWord(name, grantRouteWords) {
			answeredGrantDoor = append(answeredGrantDoor, strconv.Quote(name)+" (written at "+where+")")
		}
	}
	sort.Strings(answeredElsewhere)
	sort.Strings(answeredGrantDoor)
	return answeredElsewhere, answeredGrantDoor
}

// answeredOutsidePool names the degenerate case the other two directions cannot
// see: a route the enumeration resolved that the package never spells in one
// piece (a case label built from a concatenation whose pieces are each
// route-shaped). If that happened, the pool would be a smaller universe than the
// guard and every check built on it would be decoration.
func answeredOutsidePool(pool map[string][]string, answered map[string]bool) []string {
	var out []string
	for name := range answered {
		if _, ok := pool[name]; !ok {
			out = append(out, strconv.Quote(name))
		}
	}
	sort.Strings(out)
	return out
}

// scanGrantBoundary runs the AST instrument over one Go directory and reports
// every grant-carrying face: an inbound envelope field, or an answered route
// whose own name is an approval decision.
func scanGrantBoundary(t *testing.T, dir string) []boundaryFinding {
	t.Helper()
	pkg, err := parseBoundaryPackage(dir)
	if err != nil {
		t.Fatalf("parse %s: %v", dir, err)
	}
	requireReadableInstrument(t, dir, pkg)

	var findings []boundaryFinding
	for _, route := range sortedSet(pkg.answered) {
		if carriesGrantWord(route, grantRouteWords) {
			findings = append(findings, boundaryFinding{
				File: pkg.guardFile, Line: pkg.guardLine,
				What: "Go answers the inbound route " + strconv.Quote(route) + ", whose own name is an approval decision - a panel-side allow door (AGENTS.md §1.2 ban #6, D33/F2)",
			})
		}
	}

	inbound := inboundEnvelopes(pkg)
	if len(inbound) == 0 {
		t.Fatalf("no inbound envelope type (a JSON decode destination, or anything nested in one) under %s: the instrument claims to check what a panel request can carry but cannot find the request type", dir)
	}
	for name := range inbound {
		sh := pkg.structs[name]
		for _, f := range sh.fields {
			if f.JSONKey == "-" {
				continue
			}
			if carriesGrantWord(f.JSONKey, grantFieldWords) {
				what := "inbound envelope " + name + " can bind the JSON key " + strconv.Quote(f.JSONKey) +
					" - a verdict a page can set, which D33/F2 and R20 keep native-side only"
				if f.Text != "" {
					what += ": " + f.Text
				}
				findings = append(findings, boundaryFinding{File: sh.file, Line: f.Line, What: what})
			}
		}
	}
	sort.Slice(findings, func(i, j int) bool {
		if findings[i].File != findings[j].File {
			return findings[i].File < findings[j].File
		}
		return findings[i].Line < findings[j].Line
	})
	return findings
}

func sortedSet(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

// TestAnsweredPanelRoutesCarryNoApprovalDecision is facet 1: the inbound methods
// Go actually answers contain no approval-shaped name. The set is read out of
// this package's own route guard (and goes loud when a label cannot be read),
// and then cross-checked in BOTH directions against the function that runs: the
// enumeration may not name a door the guard refuses, and the guard may not
// answer a door the enumeration never saw.
func TestAnsweredPanelRoutesCarryNoApprovalDecision(t *testing.T) {
	root := panelRepoRoot(t)
	dir := filepath.Join(root, "internal", "panel")
	pkg, err := parseBoundaryPackage(dir)
	if err != nil {
		t.Fatalf("parse %s: %v", dir, err)
	}
	requireReadableInstrument(t, dir, pkg)

	answered := sortedSet(pkg.answered)
	for _, name := range answered {
		if !knownComposerMethod(name) {
			t.Errorf("the guard declares %q but knownComposerMethod refuses it: the answered set read here is not the one Go runs", name)
		}
		if carriesGrantWord(name, grantRouteWords) {
			t.Errorf("Go answers the inbound route %q, whose own name is an approval decision (D33/F2, R20, AGENTS.md §1.2 ban #6)", name)
		}
	}

	// Nothing approval-shaped may be answered. These are the spellings a wiring
	// commit would plausibly use, starting with the one the removed button
	// already emitted.
	for _, cand := range []string{
		"panel.approval.request", // verbatim from frontend/src/lib/panel.ts:181
		"panel.approval.grant",
		"panel.approval.resolve",
		"panel.approval.decide",
		"panel.l2.allow",
		"panel.l2.grant",
		"panel.grant",
		"panel.allow",
		"panel.mode.approve",
		"approval.decide",
		"approval.grant",
	} {
		if knownComposerMethod(cand) {
			t.Errorf("knownComposerMethod answers %q - an approval decision addressed from the panel", cand)
		}
	}

	// The oracle, pointed the other way: every name this package spells out gets
	// asked of the RUNNING guard, not of the AST. The loop above can only ever
	// report what its own enumeration found, so a second chain that
	// knownComposerMethod hands its decision to (acceptance r1's M4, and the M13
	// that combines it with a renamed route key) used to leave this test green
	// while ParseComposerRequest accepted the route.
	pool := routeNamePool(pkg)
	if len(pool) < len(answered) {
		t.Fatalf("the route-name pool holds %d names but the guard answers %d - the pool is smaller than the thing it is supposed to audit, so every finding below is decoration",
			len(pool), len(answered))
	}
	answeredElsewhere, answeredGrantDoor := poolJudgedByRealGuard(pool, pkg.answered)
	for _, hit := range answeredElsewhere {
		t.Errorf("knownComposerMethod answers %s but the guard's own case list, which is what this file enumerates, never named it: a route answered through another function, or another expression, is a door no vocabulary in this package reviewed - the answered set read here is INCOMPLETE, not clean (D33/F2, R20)", hit)
	}
	for _, hit := range answeredGrantDoor {
		t.Errorf("knownComposerMethod answers %s, an approval decision addressed from the panel, and the name was found written in this package rather than guessed here (D33/F2, R20, AGENTS.md §1.2 ban #6)", hit)
	}
	if strays := answeredOutsidePool(pool, pkg.answered); len(strays) > 0 {
		t.Errorf("the guard's case list resolves to %v, which this package never spells as one route-shaped string: the pool and the enumeration are not the same universe and the audit above is not measuring the guard", strays)
	}

	// The answered set and the declared route constants must be one vocabulary:
	// a method known to only one of them is a door nobody reviewed.
	declared := map[string]bool{}
	for name, val := range pkg.consts {
		if strings.HasPrefix(name, "Method") {
			declared[val] = true
		}
	}
	for _, name := range answered {
		if !declared[name] {
			t.Errorf("Go answers %q but no Method* constant declares it: a route written straight into the guard's case list. Nothing else in this package would notice - TestFrontendComposerRequestsMatchTheEnvelope (bridge_test.go:131) is the test that reads the frontend for these names, and it only ever iterates the four declared constants, so a fifth route that never became a constant is invisible to it. This line is the only gate on that shape", name)
		}
	}
	t.Logf("answered inbound routes = %v; %d declared Method* constants; 0 approval-shaped names on either side",
		answered, len(declared))
}

// grantRoutePrefixes are the namespaces a wiring commit would address a panel-side
// allow through. They are NAMESPACES, not route names: this file hard-codes no
// candidate it later asks about, because the shape being nailed here is exactly
// the one where a name is never written down in one piece (F-R2-3). The set is a
// judgement call; growing it is allowed but is a deliberate two-line edit - the
// member goes into grantRoutePrefixWitnesses (F-R4-2) and the count into
// wantGrantRoutePrefixes below, which is what makes a NARROWING impossible to
// write quietly, even together with its witness row.
var grantRoutePrefixes = []string{
	"panel", "panel.review", "panel.approval", "panel.l2", "panel.mode",
	"panel.workspace", "panel.attachment", "panel.message",
}

// grantRouteSuffixes are the verb tails the four declared routes actually use
// ("request", "add", "send" - the first is a suffix here and the others are not,
// which is why "now" is along: a wiring commit reaches for whatever reads well).
var grantRouteSuffixes = []string{"", ".request", ".now"}

// wantGrantRoute* are the grid's three factors as a written-down number, which is
// the part the product assertion cannot do: asked= always equals the lists as
// they now stand, so a member deleted together with its witness row leaves the
// product true. Acceptance r4 measured the same silence one level up (its
// §11: prefix-shrink asked=462, suffix-shrink asked=352, both green). The three
// lengths are named separately so the red says which factor moved. Widening the
// grid means moving these numbers on purpose; that friction is the whole design.
const (
	wantGrantRoutePrefixes = 8
	wantGrantRouteWords    = 11
	wantGrantRouteSuffixes = 3
)

// gridFactorWitness is one member of a grid factor plus one route name that only
// that member lets the sweep assemble.
type gridFactorWitness struct {
	factor string
	route  string
}

// grantRoutePrefixWitnesses and grantRouteSuffixWitnesses are the floor under the
// other two factors of the product (F-R4-2, docs/evidence/s1/panel-l2-grant-
// nail-fix-r4-accept-r1.md §11). grantRouteWords has had a witness per member
// since F-ACC-2; the two lists here had none, and acceptance r4 measured what
// that costs: deleting "panel.mode" (asked 528 -> 462) or the ".now" tail (asked
// 528 -> 352) leaves the package green, so a production guard that answers an
// approval route spelled through the pruned piece is silently unasked - its own
// M-now reading was hits=[panel.review.grant.now] red with the list intact and
// clean with ".now" gone. The product assertion cannot see this: asked= always
// equals the lists as they now stand, so it nails "a factor that stopped being
// ranged", never "a factor that got narrower".
//
// This is not a table asserting that it contains itself, for the same reason the
// word roster is not (F-ACC-2): each row is judged by three things outside the
// list it defends - the names the sweep really assembled (checkGridFactor
// Witnesses), the real predicate over the real vocabulary, and the RUNNING
// production guard.
//
// A row is one name, not the whole grid: it says "this piece is still in play",
// so the lists may still gain members freely - widening the grid is safe and
// needs no new row - while narrowing it reds. Rows are spelled out rather than
// templated, and each names a door the guard must refuse, panel.review.grant.now
// being the spelling acceptance r4 planted.
var (
	grantRoutePrefixWitnesses = []gridFactorWitness{
		{"panel", "panel.grant"},
		{"panel.review", "panel.review.approve"},
		{"panel.approval", "panel.approval.grant"},
		{"panel.l2", "panel.l2.approve"},
		{"panel.mode", "panel.mode.approve"},
		{"panel.workspace", "panel.workspace.approve"},
		{"panel.attachment", "panel.attachment.approve"},
		{"panel.message", "panel.message.approve"},
	}
	grantRouteSuffixWitnesses = []gridFactorWitness{
		{"", "panel.decide"},
		{".request", "panel.authorize.request"},
		{".now", "panel.review.grant.now"},
	}
)

// checkGridFactorWitnesses is the ruler over one factor list: four claims per
// row, and a deleted member breaks the first two at once - the piece is gone
// from the list, and the name this file calls its evidence is no longer among
// the names the sweep put to the guard.
func checkGridFactorWitnesses(t *testing.T, listName string, list []string, rows []gridFactorWitness, names []string, perMember int) {
	t.Helper()
	for _, row := range rows {
		inList := false
		for _, f := range list {
			if f == row.factor {
				inList = true
			}
		}
		if !inList {
			t.Errorf("%s no longer carries %q while this file still keeps %q as its witness: the sweep stopped asking the %d names that piece built, and a panel-side door spelled that way is now unasked rather than refused (F-R4-2)",
				listName, row.factor, row.route, perMember)
		}
		asked := false
		for _, n := range names {
			if n == row.route {
				asked = true
			}
		}
		if !asked {
			t.Errorf("the standing sweep never asked %q, the witness name for the %s entry %q: it is assembled out of grantRoutePrefixes x grantRouteWords x grantRouteSuffixes, so one of those three lists has narrowed away from the evidence written next to it (F-R4-2)",
				row.route, listName, row.factor)
		}
		if !carriesGrantWord(row.route, grantRouteWords) {
			t.Errorf("the %s entry %q is witnessed by %q, which the real predicate does not read as an approval door: this row is evidence of nothing", listName, row.factor, row.route)
		}
		if knownComposerMethod(row.route) {
			t.Errorf("the running guard answers %q, the witness name for the %s entry %q: this row stands as proof that the grid that piece builds is refused, and it is not (D33/F2, R20, AGENTS.md §1.2 ban #6)", row.route, listName, row.factor)
		}
	}
}

// sweepAssembledNames is the ONE place this file puts the assembled name to the
// RUNNING guard. It exists so the positive control and the negative assertion
// below cannot drift apart onto separate execution points (F-R4-1, docs/
// evidence/s1/panel-l2-grant-nail-fix-r4-accept-r1.md §10): while the control
// called knownComposerMethod on a line of its own, short-circuiting this
// function's ask step (if false && knownComposerMethod(name)) or its accumulate
// step left the control green and the sweep reporting hits=[] as a clean
// boundary, asked=528 and all. It returns how many names it built, the names
// themselves (so a narrowed factor list is detectable by membership, F-R4-2),
// and the names the guard answered, sorted.
func sweepAssembledNames(prefixes, words, suffixes []string) (int, []string, []string) {
	var names, hits []string
	for _, p := range prefixes {
		for _, w := range words {
			for _, suffix := range suffixes {
				for _, form := range []string{w, w + "s"} {
					name := p + "." + form + suffix
					names = append(names, name)
					if knownComposerMethod(name) {
						hits = append(hits, name)
					}
				}
			}
		}
	}
	sort.Strings(hits)
	return len(names), names, hits
}

// declaredRoutePrefixes / declaredRouteWords / declaredRouteSuffixes are the
// positive control's grid: the same three kinds of pieces the sweep multiplies,
// aimed at the namespaces of the four routes this package really declares.
// Exactly four of its 24 names are answered today, so the control says both
// halves of the fact - the guard still answers its own routes, and the
// neighbourhood around them is still refused.
var (
	declaredRoutePrefixes = []string{"panel"}
	declaredRouteWords    = []string{"mode", "workspace", "attachment", "message"}
	declaredRouteSuffixes = []string{".request", ".add", ".send"}
)

// knownComposerRefusesAssembledGrantNames is the behavioural sweep behind
// TestRealGuardRefusesEveryAssemblableApprovalRouteName. It builds
// len(grantRoutePrefixes) x len(grantRouteWords) x len(grantRouteSuffixes) x 2
// plural/no-plural names by concatenating at run time and asks the RUNNING guard
// about each one. Nothing here reads a file: a route the guard answers is enough,
// however the production code got that name.
//
// That product is an ASSERTION and not merely the shape of the loops (F-ACC-1).
// Before it, asked= was only a t.Logf reading, so emptying grantRouteSuffixes -
// the one factor the old control did not name - or dropping the plural ring made
// the sweep ask 0 or 264 names and still report a clean boundary. All three
// factors are now named in the empty-vocabulary control below, and the returned
// count is checked against the product, because a factor that stops being ranged
// is a grid that silently shrank and hits=[] over it is a different fact.
//
// The control that proves "hits=[] means refused" and not "means never asked" is
// sweepAssembledNames' own ask line, reached from both sides at once (F-R4-1).
func knownComposerRefusesAssembledGrantNames(t *testing.T) (int, []string) {
	t.Helper()
	if len(grantRouteWords) == 0 || len(grantRoutePrefixes) == 0 || len(grantRouteSuffixes) == 0 {
		t.Fatalf("the sweep vocabulary is empty (grantRouteWords=%d, grantRoutePrefixes=%d, grantRouteSuffixes=%d): a sweep that asks nothing answers clean forever",
			len(grantRouteWords), len(grantRoutePrefixes), len(grantRouteSuffixes))
	}
	if len(grantRoutePrefixes) != wantGrantRoutePrefixes || len(grantRouteWords) != wantGrantRouteWords || len(grantRouteSuffixes) != wantGrantRouteSuffixes {
		t.Fatalf("the sweep grid is %d prefixes x %d words x %d suffixes per plural form, not the %d x %d x %d written down in wantGrantRoute*: a factor list moved without the pin moving with it. The product assertion below cannot see this, because asked= always equals the lists as they now stand - a narrowing paired with the deletion of its own witness row is what this line is for (F-R4-2)",
			len(grantRoutePrefixes), len(grantRouteWords), len(grantRouteSuffixes),
			wantGrantRoutePrefixes, wantGrantRouteWords, wantGrantRouteSuffixes)
	}
	// Anti-vacuity, positive control: the same sweep line that must find no
	// approval-shaped name below has to find the four routes this package really
	// declares, in a grid of its own. It is the same line on purpose - see
	// sweepAssembledNames.
	wantAnswered := []string{MethodModeRequest, MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend}
	sort.Strings(wantAnswered)
	_, _, legit := sweepAssembledNames(declaredRoutePrefixes, declaredRouteWords, declaredRouteSuffixes)
	if !reflect.DeepEqual(legit, wantAnswered) {
		t.Fatalf("the running guard answers %v of the %d names sweepAssembledNames builds from the four declared route namespaces (got %v, want %v): a sweep whose ask-or-accumulate step has been short-circuited reports hits=[] and reads as a clean boundary, which is not the same fact (F-R4-1)",
			len(legit), len(declaredRoutePrefixes)*len(declaredRouteWords)*len(declaredRouteSuffixes)*2, legit, wantAnswered)
	}
	asked, _, hits := sweepAssembledNames(grantRoutePrefixes, grantRouteWords, grantRouteSuffixes)
	if want := len(grantRoutePrefixes) * len(grantRouteWords) * len(grantRouteSuffixes) * 2; asked != want {
		t.Fatalf("the sweep asked %d names but the three lists it reads multiply to %d (%d prefixes x %d words x %d suffixes x 2 plural forms): a factor that stopped being ranged is a grid that silently shrank, and hits=%v over the smaller grid is not the verdict this file documents (F-ACC-1)",
			asked, want, len(grantRoutePrefixes), len(grantRouteWords), len(grantRouteSuffixes), hits)
	}
	return asked, hits
}

// TestRealGuardRefusesEveryAssemblableApprovalRouteName is the standing half of
// facet 1's wider oracle, and it exists because that oracle only reads what this
// package WRITES. Acceptance r2 (§7.2, F-R2-3) found that the one assertion
// refusing "the real knownComposerMethod answers panel.review.allow" was living
// inside facet 4's plant F subtest as a negative control: rewrite that subtest,
// Skip it, or rename it and the ban on the guard answering an approval route
// through a runtime-assembled name disappears with no other test going red.
//
// So the check is a top-level test now, and it asks no name written into this
// file: it concatenates prefixes out of grantRoutePrefixes with the words in
// grantRouteWords and puts the question to the running guard. That is what makes
// it cover the shape a literal pool cannot see - production code that assembles
// "panel.review.allow" out of pieces is answering a name no scan ever reads -
// and it is also why the sweep must stay free of false reds: it is only worth
// anything while a clean tree answers none of it.
//
// What it does NOT cover is the same line the file header draws: a name outside
// grantRoutePrefixes x grantRouteWords ("wisp.review.ok", or one read off config)
// is not asked here either. This sweep narrows that hole, it does not close it.
func TestRealGuardRefusesEveryAssemblableApprovalRouteName(t *testing.T) {
	asked, hits := knownComposerRefusesAssembledGrantNames(t)
	for _, hit := range hits {
		t.Errorf("knownComposerMethod answers %q, an approval-shaped route name this file assembled at run time out of grantRoutePrefixes and grantRouteWords: Go answers a panel-side allow door whose spelling never appears in the package, so no literal pool, no case-list enumeration and no vocabulary scan reaches it (D33/F2, R20, AGENTS.md §1.2 ban #6)", hit)
	}
	// plant F's pool half depends on this name being unanswered but never written
	// into the guard, so say out loud which half of that pair this test owns.
	if strings.Contains(strings.Join(hits, " "), "panel.review.allow") {
		t.Errorf("the refused name plant F plants is now answered by the guard, so facet 4's plant F subtest and facet 1's answeredElsewhere half are describing different trees")
	}
	t.Logf("behavioural sweep of the running guard: asked=%d assembled route names, answeredByRealGuard=%d, hits=%v",
		asked, len(hits), hits)
}

// TestNoInboundEnvelopeCanBindAnApprovalVerdict is facet 2 on the production
// types: the single envelope every answered method decodes into, and everything
// reachable from it, cannot bind a key whose name carries an allow.
func TestNoInboundEnvelopeCanBindAnApprovalVerdict(t *testing.T) {
	keys := bindableKeysOf(reflect.TypeOf(ComposerRequest{}), 0, "ComposerRequest", map[reflect.Type]bool{})
	if len(keys) == 0 {
		t.Fatal("ComposerRequest binds no JSON keys at all - this check would pass by seeing nothing")
	}
	if hits := grantCarryingKeys(keys); len(hits) > 0 {
		var lines []string
		for _, h := range hits {
			lines = append(lines, h.GoPath+" binds "+strconv.Quote(h.Key))
		}
		t.Errorf("a panel request can carry an approval verdict through %d field(s); D33/F2 and R20 keep that verdict native-side only, and AGENTS.md §1.2 bans the shape:\n  %s",
			len(hits), strings.Join(lines, "\n  "))
	}

	// The same rule over the payload shapes the answered routes resolve into.
	for _, typ := range []reflect.Type{reflect.TypeOf(ModeRequest{}), reflect.TypeOf(AttachmentPayload{})} {
		if hits := grantCarryingKeys(bindableKeysOf(typ, 0, typ.Name(), map[reflect.Type]bool{})); len(hits) > 0 {
			var lines []string
			for _, h := range hits {
				lines = append(lines, h.GoPath+" binds "+strconv.Quote(h.Key))
			}
			t.Errorf("%s carries a verdict field, which makes it an allow payload:\n  %s", typ.Name(), strings.Join(lines, "\n  "))
		}
	}

	// The AST twin, over the real package: same rule, able to name a line.
	root := panelRepoRoot(t)
	findings := scanGrantBoundary(t, filepath.Join(root, "internal", "panel"))
	if len(findings) > 0 {
		t.Errorf("the panel's inbound Go boundary has a grant-carrying face:%s", joinFindings(findings))
	}
	if hits := grantCarryingKeys(keys); len(hits) == 0 && len(findings) == 0 {
		t.Logf("ComposerRequest binds %d JSON keys %v; 0 verdict-shaped; AST scan of internal/panel: 0 findings",
			len(keys), sortedKeys(keys))
	}
}

// TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes is the predicate's own
// teeth, kept inside a green run: the planted grant shapes below must be flagged
// and the four legitimate routes must not, or the vocabulary is either decorative
// or so wide it gets widened away.
//
// Its last block answers for the SWEEP vocabulary rather than the field one:
// every member of grantRouteWords carries a witness route that the RUNNING guard
// refuses, and the roster and the vocabulary are compared in both directions
// (F-ACC-2). Acceptance r3 measured that the four historical names above cover
// four of eleven words, so deleting, say, "ratify" from grantRouteWords left 48
// candidate names unasked with nothing in the package going red.
func TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes(t *testing.T) {
	type Decision struct{ Verdict string }
	type plantedOutcomeEnvelope struct {
		Method  string `json:"method"`
		To      string `json:"to"`
		Outcome string `json:"outcome"` // the field panel.ts:169-186 sends
	}
	type plantedAllowFlag struct {
		Method string `json:"method"`
		Allow  bool   `json:"allowOnce"`
	}
	type plantedVerdictIn struct {
		Method string `json:"method"`
		Reason string `json:"reason"`
		Decision
	}

	cases := []struct {
		name string
		typ  reflect.Type
	}{
		{"outcome on the envelope", reflect.TypeOf(plantedOutcomeEnvelope{})},
		{"allowOnce on the envelope", reflect.TypeOf(plantedAllowFlag{})},
		{"a verdict smuggled inside an embedded struct", reflect.TypeOf(plantedVerdictIn{})},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hits := grantCarryingKeys(bindableKeysOf(tc.typ, 0, tc.typ.Name(), map[reflect.Type]bool{}))
			if len(hits) == 0 {
				t.Fatalf("planted a grant-carrying envelope (%s) and the predicate reported nothing - the check is decoration", tc.name)
			}
			t.Logf("flagged as required: %s binds %s", hits[0].GoPath, strconv.Quote(hits[0].Key))
		})
	}

	for _, route := range []string{"panel.approval.request", "panel.l2.allow", "approval.decide", "panel.grant"} {
		if !carriesGrantWord(route, grantRouteWords) {
			t.Errorf("route %q reads as an innocent name - the answered-set nail cannot catch a wired grant door", route)
		}
	}
	for _, route := range []string{MethodModeRequest, MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend} {
		if carriesGrantWord(route, grantRouteWords) {
			t.Errorf("legitimate composer route %q reads as an approval door - the vocabulary would cry wolf", route)
		}
	}

	// F-ACC-2, the sweep vocabulary's own witnesses. grantRouteWords is multiplied
	// out into the grid by TestRealGuardRefusesEveryAssemblableApprovalRouteName,
	// so a word that disappears from it takes 48 candidate names with it - and
	// until this block nothing in the package could tell that had happened
	// (acceptance r3 §1.2 deleted "ratify", stacked its own M16 on top of that, and
	// read 53 top-level PASSes). Each row makes three claims and only the last one
	// is about this list: the REAL predicate over the REAL vocabulary must flag
	// route by exactly the recorded set of words, the name assembled from the word
	// at run time must read as an approval door and be refused by the RUNNING
	// guard, and the written-out route must be refused by it too. That refusal is
	// a fact about bridge.go, not about the roster.
	for _, row := range grantRouteWordWitnesses {
		var flagged []string
		for _, word := range grantRouteWords {
			if carriesGrantWord(row.route, []string{word}) {
				flagged = append(flagged, word)
			}
		}
		want := append([]string(nil), row.carriers...)
		sort.Strings(flagged)
		sort.Strings(want)
		if !reflect.DeepEqual(flagged, want) {
			t.Errorf("witness route %q is flagged by %v under grantRouteWords as it stands, and this file wrote down %v for it: the vocabulary and the evidence attached to it are no longer the same list, which is what a silently narrowed sweep looks like from outside (F-ACC-2)", row.route, flagged, want)
		}
		if knownComposerMethod(row.route) {
			t.Errorf("the running guard answers %q, the witness name for the vocabulary word %q: this row stands as evidence that a panel-side allow door spelled that way is refused, and it is not (D33/F2, R20, AGENTS.md §1.2 ban #6)", row.route, row.word)
		}
		// The same claim over a name no list hands this file: built from the word
		// at run time, so a widened vocabulary puts its own word in front of the
		// guard instead of waiting for someone to remember the roster.
		assembled := "panel." + row.word + ".request"
		if !carriesGrantWord(assembled, grantRouteWords) {
			t.Errorf("the route name assembled here from the vocabulary word %q (%q) does not read as an approval door - that word no longer carries the meaning the sweep assumes when it multiplies it into a grid", row.word, assembled)
		}
		if knownComposerMethod(assembled) {
			t.Errorf("the running guard answers %q, assembled at run time from the vocabulary word %q (D33/F2, R20, AGENTS.md §1.2 ban #6)", assembled, row.word)
		}
	}

	// Roster and vocabulary, both directions: a word with no witness is a hole in
	// the grid nobody asks about, a witness with no word is a claim about a name
	// nothing sweeps any more. Deleting a member of grantRouteWords is the shape
	// F-ACC-2 is about, and it lands on the second message.
	inVocab := map[string]bool{}
	for _, word := range grantRouteWords {
		inVocab[word] = true
	}
	rostered := map[string]bool{}
	for _, row := range grantRouteWordWitnesses {
		if !inVocab[row.word] {
			t.Errorf("grantRouteWords no longer carries %q while %q still stands here as its witness: the sweep stopped asking the %d names that word used to build (len(grantRoutePrefixes) x len(grantRouteSuffixes) x 2 plural forms), and this roster is the only thing in the package that noticed (F-ACC-2)", row.word, row.route, len(grantRoutePrefixes)*len(grantRouteSuffixes)*2)
		}
		rostered[row.word] = true
	}
	for _, word := range grantRouteWords {
		if !rostered[word] {
			t.Fatalf("grantRouteWords carries %q with no witness row in grantRouteWordWitnesses: no name built from it has been put to the running guard here, so this file cannot claim the grid it multiplies that word into is clean", word)
		}
	}

	// F-R4-2, the same rule over the two factors that had no witnesses: the names
	// below are taken from the standing sweep's own output, so this is the sweep
	// reporting what it asked rather than this file re-deriving a product.
	_, names, _ := sweepAssembledNames(grantRoutePrefixes, grantRouteWords, grantRouteSuffixes)
	checkGridFactorWitnesses(t, "grantRoutePrefixes", grantRoutePrefixes, grantRoutePrefixWitnesses, names,
		len(grantRouteWords)*len(grantRouteSuffixes)*2)
	checkGridFactorWitnesses(t, "grantRouteSuffixes", grantRouteSuffixes, grantRouteSuffixWitnesses, names,
		len(grantRoutePrefixes)*len(grantRouteWords)*2)
}

// ---------------------------------------------------------------------------
// The two instruments, held to one rule.
// ---------------------------------------------------------------------------

// inboundTypeRegistry names the production types this file can talk to by
// reflection. The AST reads type names out of the package; reflection cannot
// enumerate a package's types, so the pairing has to be written down - and
// TestJSONKeyDerivationAgreesWithEncodingJSON goes loud when an inbound type is
// missing from it, so the list cannot rot into covering less than the scan.
func inboundTypeRegistry() map[string]reflect.Type {
	return map[string]reflect.Type{
		"ComposerRequest":   reflect.TypeOf(ComposerRequest{}),
		"ModeRequest":       reflect.TypeOf(ModeRequest{}),
		"AttachmentPayload": reflect.TypeOf(AttachmentPayload{}),
		"AttachmentRef":     reflect.TypeOf(AttachmentRef{}),
	}
}

// astTopKeys lists the JSON keys one declared struct can be addressed by at its
// own level: no recursion, no promoted fields, so the three sides of this check
// are comparing one thing. Embedded fields without a name are promoted and are
// left out, exactly as encoding/json leaves them out of the parent's key set.
func astTopKeys(sh structShape) []string {
	var out []string
	for _, f := range sh.fields {
		if f.JSONKey == "-" {
			continue
		}
		if f.Anonymous && f.JSONKey == "" {
			continue // promoted: encoding/json gives the parent no key of its own
		}
		if !f.Anonymous && f.GoName != "" && !token.IsExported(f.GoName) {
			continue // reflection skips unexported fields; so does encoding/json
		}
		out = append(out, f.JSONKey)
	}
	sort.Strings(out)
	return out
}

// reflectionTopKeys is the same list built by asking the field tags the compiler
// sees, with the fallback encoding/json applies to a name-less tag.
func reflectionTopKeys(typ reflect.Type) []string {
	var out []string
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if !f.IsExported() {
			continue
		}
		name := strings.Split(f.Tag.Get("json"), ",")[0]
		if name == "-" {
			continue
		}
		if f.Anonymous && name == "" {
			continue
		}
		if name == "" {
			name = f.Name
		}
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// decoderKnowsKey asks encoding/json itself - DisallowUnknownFields tells the
// difference between "no field answers to that name" and "a field answered and
// the value did not fit", and both of the latter mean the key is bindable.
func decoderKnowsKey(typ reflect.Type, key string) (bool, error) {
	doc := "{" + strconv.Quote(key) + ":\"grant\"}"
	dec := json.NewDecoder(strings.NewReader(doc))
	dec.DisallowUnknownFields()
	err := dec.Decode(reflect.New(typ).Interface())
	switch {
	case err == nil:
		return true, nil
	case strings.Contains(err.Error(), "unknown field"):
		return false, nil
	case strings.Contains(err.Error(), "cannot unmarshal"):
		return true, nil
	default:
		return false, fmt.Errorf("%s: %w", key, err)
	}
}

// TestJSONKeyDerivationAgreesWithEncodingJSON pins F-5 from the other side. The
// header of this file has always claimed the AST scan and the reflection walk are
// "the same rule, written twice"; acceptance r1's M8 showed a spelling where they
// were not, in the direction that hides a field. Three things have to agree now,
// and each one is checked against a different authority:
//
//	(i)   jsonOr's rule, spelled out as a table;
//	(ii)  the AST's key list and the reflect.StructTag key list, per inbound type;
//	(iii) encoding/json's own verdict, via DisallowUnknownFields, on every name
//	      each of the two produced - and on a list of approval-shaped names that
//	      must NOT be bindable.
func TestJSONKeyDerivationAgreesWithEncodingJSON(t *testing.T) {
	t.Run("i the AST key rule matches encoding/json's, spelling by spelling", func(t *testing.T) {
		cases := []struct {
			goName, jsonName string
			tagged           bool
			want             string
			why              string
		}{
			{"Outcome", "", false, "Outcome", "no tag: the Go field name is the key"},
			{"Outcome", "outcome", true, "outcome", "named tag"},
			{"Outcome", "", true, "Outcome", `name-less tag json:",omitempty" - the F-5 shape`},
			{"Outcome", "-", true, "-", `json:"-" is skipped by the caller, not renamed`},
			{"sizeBytes", "sizeBytes", true, "sizeBytes", "the caller already split the options off the tag"},
		}
		for _, tc := range cases {
			if got := jsonOr(tc.goName, tc.jsonName, tc.tagged); got != tc.want {
				t.Errorf("jsonOr(%q,%q,%v) = %q, want %q (%s)", tc.goName, tc.jsonName, tc.tagged, got, tc.want, tc.why)
			}
		}
		// The reflection side must land on the same answers, or the two halves of
		// this file are still two rules.
		type f5Probe struct {
			Named     string `json:"named,omitempty"`
			Nameless  string `json:",omitempty"`
			Untagged  string
			Skipped   string `json:"-"`
			Recursive string `json:"rec,-"`
		}
		want := map[string]bool{"named": true, "Nameless": true, "Untagged": true, "rec": true}
		for _, k := range reflectionTopKeys(reflect.TypeOf(f5Probe{})) {
			if !want[k] {
				t.Errorf("reflection produced key %q, which encoding/json would not address it as", k)
			}
			delete(want, k)
		}
		for k := range want {
			t.Errorf("encoding/json can address %q but the reflection walk did not list it", k)
		}
	})

	root := panelRepoRoot(t)
	dir := filepath.Join(root, "internal", "panel")
	pkg, err := parseBoundaryPackage(dir)
	if err != nil {
		t.Fatalf("parse %s: %v", dir, err)
	}
	requireReadableInstrument(t, dir, pkg)
	reg := inboundTypeRegistry()

	// The registry cannot be allowed to rot into covering less than the scan.
	for seed := range pkg.inboundSeeds {
		if _, ok := reg[seed]; !ok {
			t.Fatalf("decode destination %s has no reflection twin in inboundTypeRegistry: this test would compare the two instruments over different trees, and the AST half would be unchecked", seed)
		}
	}

	t.Run("ii the AST list and the reflection list are one list", func(t *testing.T) {
		for name, typ := range reg {
			sh, ok := pkg.structs[name]
			if !ok {
				t.Fatalf("%s is in the registry but not in the parsed package - the registry names a type that is gone", name)
			}
			astKeys, reflKeys := astTopKeys(sh), reflectionTopKeys(typ)
			if !equalStrings(astKeys, reflKeys) {
				t.Errorf("%s: the AST half reads %v and the reflection half reads %v. Two instruments that claim one rule are disagreeing, which is the F-5 defect class regardless of which side is wrong",
					name, astKeys, reflKeys)
			}
			t.Logf("%s: %d keys agreed by both instruments %v", name, len(astKeys), astKeys)
		}
	})

	t.Run("iii the decoder is the third vote and the verdict words are not bindable", func(t *testing.T) {
		// Names a page would use to hand Go an approval decision. Both spellings of
		// each: encoding/json falls back to the Go field name, so a field called
		// Outcome is addressable as "outcome" and as "Outcome".
		candidates := []string{
			"outcome", "Outcome", "allow", "Allow", "allowOnce", "AllowOnce",
			"approved", "Approved", "grant", "Grant", "verdict", "Verdict",
			"decision", "Decision", "decide", "Decide", "bypass", "Bypass",
			"override", "Override", "permit", "Permit", "authorize", "Authorize",
		}
		for _, name := range sortedRegistryNames(reg) {
			typ := reg[name]
			sh := pkg.structs[name]
			astKeys, reflKeys := astTopKeys(sh), reflectionTopKeys(typ)
			for _, key := range reflKeys {
				known, err := decoderKnowsKey(typ, key)
				if err != nil {
					t.Fatalf("%s: asking encoding/json about %q: %v", name, key, err)
				}
				if !known {
					t.Errorf("%s: both instruments list %q but the decoder refuses it - the key rule in this file is not encoding/json's rule", name, key)
				}
			}
			for _, key := range candidates {
				known, err := decoderKnowsKey(typ, key)
				if err != nil {
					t.Fatalf("%s: asking encoding/json about %q: %v", name, key, err)
				}
				if known {
					t.Errorf("%s binds the wire key %q, so a panel can address an approval decision through it (D33/F2, AGENTS.md ban #6). If %s is inbound, this is the finding; if it is not inbound any more, drop it from inboundTypeRegistry instead of dropping the check", name, key, name)
				}
			}
			if extra := minusStrings(astKeys, reflKeys); len(extra) > 0 {
				t.Errorf("%s: the AST lists %v, which the reflection walk and the decoder both refuse to confirm", name, extra)
			}
			t.Logf("%s: %d listed keys confirmed bindable, %d verdict spellings confirmed not bindable",
				name, len(reflKeys), len(candidates))
		}
	})
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func minusStrings(a, b []string) []string {
	var out []string
	for _, x := range a {
		found := false
		for _, y := range b {
			if x == y {
				found = true
				break
			}
		}
		if !found {
			out = append(out, x)
		}
	}
	return out
}

func sortedRegistryNames(reg map[string]reflect.Type) []string {
	out := make([]string, 0, len(reg))
	for name := range reg {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// TestGrantWireShapesAreRefusedAtTheDoor is facet 3, behavioural: what actually
// happens on the wire today when a request carries a verdict.
func TestGrantWireShapesAreRefusedAtTheDoor(t *testing.T) {
	t.Run("the historical wire shape is refused because Go does not answer its route", func(t *testing.T) {
		// The body l2-approval-card.tsx:163 emitted via panel.ts:181, verbatim,
		// plus the identity fields the current envelope demands.
		for _, raw := range []string{
			`{"method":"panel.approval.request","correlationId":"c-1","outcome":"grant"}`,
			`{"method":"panel.approval.request","requestId":"r-1","source":"panel-composer","correlationId":"c-1","outcome":"grant"}`,
			`{"method":"approval.decide","requestId":"r-1","source":"panel-composer","outcome":"grant"}`,
			`{"method":"panel.l2.allow","requestId":"r-1","source":"panel-composer","allow":true}`,
			`{"method":"panel.grant","requestId":"r-1","source":"panel-composer","approved":true}`,
		} {
			parsed, err := ParseComposerRequest(raw)
			if err == nil {
				t.Errorf("ParseComposerRequest accepted %q (parsed as %+v) - a grant request walked through the door D33/F2 keeps shut", raw, parsed)
				continue
			}
			if !errors.Is(err, ErrComposerRequest) {
				t.Errorf("%q refused with an unexpected error type: %v", raw, err)
			}
		}
	})

	t.Run("every answered route drops a smuggled verdict", func(t *testing.T) {
		smuggled := []string{"outcome", "allowOnce", "decision", "verdict", "approved", "grant"}
		for _, method := range []string{MethodModeRequest, MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend} {
			extra := map[string]any{}
			for _, k := range smuggled {
				extra[k] = "grant"
			}
			raw := okEnvelope(t, method, extra)
			parsed, err := ParseComposerRequest(raw)
			if err != nil {
				// Refusing the whole request is also a safe answer; say which one
				// happened instead of claiming the field went nowhere.
				t.Logf("%s refused outright: %v", method, err)
				continue
			}
			back, err := json.Marshal(parsed)
			if err != nil {
				t.Fatalf("marshal %s: %v", method, err)
			}
			for _, key := range smuggled {
				if strings.Contains(string(back), `"`+key+`"`) {
					t.Errorf("%s: a request arriving with %q comes back out still holding it, so a handler with this value could act on a panel-side allow (re-encoded: %s)",
						method, key, back)
				}
			}
			if hits := grantCarryingKeys(bindableKeysOf(reflect.TypeOf(parsed), 0, "ComposerRequest", map[reflect.Type]bool{})); len(hits) > 0 {
				t.Errorf("%s parsed into a type with a verdict field: %v", method, hits)
			}
		}
	})
}

// TestPlantedGrantWiringGoesRedInASnapshot is facet 4: the teeth. It copies this
// package's Go sources into a temporary directory and plants the shapes the ban
// exists for, then requires the same instrument to name the planted lines.
//
// The tree being mutated is a COPY. Nothing here writes into frontend/, design/,
// or this package's own sources: the plants exist only under t.TempDir(), which
// is outside the repository, and the historical JSX is planted as fixture text,
// not matched as a pattern.
func TestPlantedGrantWiringGoesRedInASnapshot(t *testing.T) {
	root := panelRepoRoot(t)
	src := filepath.Join(root, "internal", "panel")
	files, err := goSourceFiles(src)
	if err != nil {
		t.Fatalf("list %s: %v", src, err)
	}
	if len(files) == 0 {
		t.Fatalf("no Go sources under %s", src)
	}

	bridgeReal, err := os.ReadFile(filepath.Join(src, "bridge.go"))
	if err != nil {
		t.Fatalf("read bridge.go: %v", err)
	}

	// The pristine copy must read clean, or every plant below would be red for
	// the wrong reason.
	if before := scanGrantBoundary(t, copyGoSources(t, src, files, string(bridgeReal), "", "")); len(before) > 0 {
		t.Fatalf("the snapshot of the real package is already dirty before any planting:%s", joinFindings(before))
	}

	fieldAnchor := "\tText string `json:\"text,omitempty\"`"
	if !strings.Contains(string(bridgeReal), fieldAnchor) {
		t.Fatalf("plant A has no anchor: %q is not in bridge.go - the envelope moved, so this mutation no longer tests what it claims to", fieldAnchor)
	}
	// Plant A: the envelope grows the field the deleted button was carrying.
	plantABridge := strings.Replace(string(bridgeReal), fieldAnchor,
		fieldAnchor+"\n\tOutcome string `json:\"outcome,omitempty\"`", 1)

	guardAnchor := "case MethodModeRequest, MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend:"
	if !strings.Contains(string(bridgeReal), guardAnchor) {
		t.Fatalf("plant B has no anchor: the route guard's case list is no longer %q - the enumeration in this file has gone stale", guardAnchor)
	}
	// Plant B: Go starts answering the approval route. Derived from the pristine
	// bridge, so B tests exactly one change.
	plantBBridge := strings.Replace(string(bridgeReal), guardAnchor,
		strings.Replace(guardAnchor, ":", ", \"panel.approval.request\":", 1), 1)

	// Plant E: the F-1 shape. The guard reaches two routes through expressions
	// this scan cannot resolve - a package-level var and a concatenation - which
	// is exactly how a wiring commit hides a door from a literals-and-constants
	// enumeration. Unlike plant B this one never compiles into the snapshot's
	// runtime truth, so what is asserted below is the drop list, not a finding.
	plantEBridge := strings.Replace(string(bridgeReal), guardAnchor,
		strings.Replace(guardAnchor, ":", ", acceptE2Route, \"panel.review.\" + acceptE3Tail:", 1), 1) +
		"\n// planted by facet 4: a route no literals-and-constants scan can name.\n" +
		"var acceptE2Route = \"panel.review.allow\"\n\n" +
		"const acceptE3Tail = \"grant\"\n"

	// Plant C: a throwaway inbound handler - the shape a follow-up commit writes
	// when it "just needs" the panel to pass a verdict through.
	const plantCHandler = "package panel\n\n" +
		"import \"encoding/json\"\n\n" +
		"// throwaway: answers a panel request by trusting what the panel said.\n" +
		"type grantThrough struct {\n" +
		"\tMethod  string `json:\"method\"`\n" +
		"\tOutcome string `json:\"outcome\"`\n" +
		"}\n\n" +
		"func handlePanelGrant(raw string) error {\n" +
		"\tvar g grantThrough\n" +
		"\treturn json.Unmarshal([]byte(raw), &g)\n" +
		"}\n"

	// The verbatim historical button, git show 53a1359^:frontend/src/components/
	// l2-approval-card.tsx lines 161-167, planted as FIXTURE text to measure what
	// this instrument cannot see.
	const plantJSX = "        <Button\n" +
		"          className=\"rounded-control bg-accent text-[var(--accent-fg)] hover:bg-accent-hover\"\n" +
		"          onClick={() => send(\"grant\")}\n" +
		"          size=\"sm\"\n" +
		"        >\n" +
		"          本次允许\n" +
		"        </Button>\n"

	t.Run("A an envelope that grows an outcome field goes red", func(t *testing.T) {
		dir := copyGoSources(t, src, files, plantABridge, "l2-approval-card.tsx", plantJSX)
		hits := scanGrantBoundary(t, dir)
		if len(hits) == 0 {
			t.Fatal("planted an inbound envelope that binds \"outcome\" and NOTHING went red - the nail is decorative")
		}
		if !findingsName(hits, "bridge.go", `"outcome"`) {
			t.Errorf("a finding must name bridge.go and the outcome key, got:%s", joinFindings(hits))
		}
		if findingsName(hits, "bridge.go", `"panel.approval.request"`) {
			t.Errorf("plant A must not also answer the approval route - the snapshot leaked between plants:%s", joinFindings(hits))
		}
		t.Logf("red as required:%s", joinFindings(hits))
	})

	t.Run("B Go answering panel.approval.request goes red", func(t *testing.T) {
		dir := copyGoSources(t, src, files, plantBBridge, "", "")
		hits := scanGrantBoundary(t, dir)
		if !findingsName(hits, "bridge.go", `"panel.approval.request"`) {
			t.Errorf("answering the approval route must be named as an answered grant door, got:%s", joinFindings(hits))
		}
		if findingsName(hits, "bridge.go", `"outcome"`) {
			t.Errorf("plant B adds no field, so an outcome-field finding here means the scan is not reading what it is given:%s", joinFindings(hits))
		}
		t.Logf("red as required:%s", joinFindings(hits))
	})

	// The F-1 plant: the guard's own switch now carries two labels this scan
	// cannot turn into strings. Before this file's fix they were dropped silently
	// and the boundary read clean, which is how acceptance r1 got four facets to
	// pass while ParseComposerRequest accepted a panel route. The drop list is the
	// load-bearing ingredient: guardReadabilityProblem is what turns it into the
	// abort that facet 1 and facet 2 run on every real scan.
	t.Run("E a guard route reached through a var or a concatenation is named, not dropped", func(t *testing.T) {
		dir := copyGoSources(t, src, files, plantEBridge, "", "")
		pkg, err := parseBoundaryPackage(dir)
		if err != nil {
			t.Fatalf("parse %s: %v", dir, err)
		}
		if len(pkg.dropped) != 2 {
			t.Fatalf("planted 2 unreadable case labels, the enumeration reported %d - the drop list is not reading the guard:%s",
				len(pkg.dropped), joinUnresolved(pkg.dropped))
		}
		if got := len(sortedSet(pkg.answered)); got != 4 {
			t.Errorf("the plant must still resolve the 4 declared routes so this run cannot be passing by having stopped reading the guard, got %d", got)
		}
		problem := instrumentBlindnessProblem(dir, pkg)
		if problem == "" {
			t.Fatal("an unreadable case label produced no problem: facet 1 would report this boundary clean")
		}
		for _, want := range []string{"acceptE2Route", `"panel.review." + acceptE3Tail`} {
			if !strings.Contains(problem, want) {
				t.Errorf("the loud failure must name the expression it could not resolve (%s), got: %s", want, problem)
			}
		}
		t.Logf("loud as required, and the names are not invented:%s", joinUnresolved(pkg.dropped))
	})

	// Plant F: the M4 shape - a SECOND chain, in its own function, answering an
	// approval-shaped route. A snapshot is only parsed, never compiled, so what
	// plant F can pin is the half that was actually missing - the pool the oracle
	// reads is built from the whole package, so it sees a name written in any
	// function at all, where the guard's own case list sees nothing.
	// It used to also carry the runtime half ("and the real guard still refuses
	// this name") as its negative control, which acceptance r2 §7.2 recorded as
	// F-R2-3: a ban on production code borrowed another subtest's line, so
	// rewriting or skipping this subtest would have deleted the only thing
	// refusing an answered panel-side allow door. That assertion lives in
	// TestRealGuardRefusesEveryAssemblableApprovalRouteName now, where nothing has
	// to reach it through this plant. The runtime half of the pair was exercised
	// on real trees in docs/evidence/s1/panel-l2-grant-nail-fix-r2.md §6 (M4 and
	// M13 both red).
	const plantFChain = "package panel\n\n" +
		"// acceptFGate is the second chain a wiring commit writes when it does not\n" +
		"// want to touch the guard's own case list.\n" +
		"func acceptFGate(m string) bool {\n" +
		"\tswitch m {\n" +
		"\tcase \"panel.review.allow\":\n" +
		"\t\treturn true\n" +
		"\t}\n" +
		"\treturn false\n" +
		"}\n"

	t.Run("F a route named outside the guard still reaches the pool", func(t *testing.T) {
		dir := copyGoSources(t, src, files, "", "l2-plant-f.go", plantFChain)
		pkg, err := parseBoundaryPackage(dir)
		if err != nil {
			t.Fatalf("parse %s: %v", dir, err)
		}
		requireReadableInstrument(t, dir, pkg)
		pool := routeNamePool(pkg)
		sites, ok := pool["panel.review.allow"]
		if !ok || len(sites) == 0 {
			t.Fatal("the pool never saw a route name written in another function: the oracle would be reading the same one case list as before")
		}
		if !strings.Contains(strings.Join(sites, " "), "l2-plant-f.go") {
			t.Errorf("the pool must name where it found the route, got %v", sites)
		}
		if pkg.answered["panel.review.allow"] {
			t.Errorf("the guard's own case list resolved the planted chain, so the answeredElsewhere half below would be testing nothing")
		}
		t.Logf("pool sees %q at %v while the snapshot's guard case list does not; the half that asks the RUNNING guard about that name - and about every other name grantRoutePrefixes x grantRouteWords can build - is TestRealGuardRefusesEveryAssemblableApprovalRouteName, which is standing rather than a negative control here (F-R2-3)",
			"panel.review.allow", sites)
	})

	// Plant G: the F-4 shape from acceptance r1. An inbound envelope that carries
	// its route on a "cmd" key and its verdict on "outcome", with its own real
	// json.Unmarshal. The old seed criterion - "a struct that binds a method key"
	// - walked right past it, because the whole point of the shape is that it does
	// not bind "method". The new seed is the decode destination itself, which the
	// shape cannot avoid by construction: no decode, nothing received.
	const plantGHandler = "package panel\n\n" +
		"import \"encoding/json\"\n\n" +
		"// throwaway: an inbound envelope that renamed its route field.\n" +
		"type grantViaCmd struct {\n" +
		"\tCmd     string `json:\"cmd\"`\n" +
		"\tOutcome string `json:\"outcome\"`\n" +
		"}\n\n" +
		"func handleGrantViaCmd(raw string) error {\n" +
		"\tvar e grantViaCmd\n" +
		"\treturn json.Unmarshal([]byte(raw), &e)\n" +
		"}\n"

	t.Run("G an inbound envelope with a renamed route key is still seeded from its decode", func(t *testing.T) {
		dir := copyGoSources(t, src, files, string(bridgeReal), "l2-plant-g.go", plantGHandler)
		pkg, err := parseBoundaryPackage(dir)
		if err != nil {
			t.Fatalf("parse %s: %v", dir, err)
		}
		requireReadableInstrument(t, dir, pkg)
		if !pkg.inboundSeeds["grantViaCmd"] {
			t.Fatalf("a type with its own json.Unmarshal did not become an inbound seed - the seed is still reading for a name, not for a decode. seeds=%v",
				sortedSet(pkg.inboundSeeds))
		}
		oldCriterionSawIt := false
		for _, f := range pkg.structs["grantViaCmd"].fields {
			if strings.EqualFold(f.JSONKey, "method") {
				oldCriterionSawIt = true
			}
		}
		if oldCriterionSawIt {
			t.Error("plant G binds a method key, so it no longer separates the two seed criteria and cannot show what it claims to")
		}
		hits := scanGrantBoundary(t, dir)
		if !findingsName(hits, "l2-plant-g.go", `"outcome"`) {
			t.Errorf("plant G is the F-4 shape and produced no finding - the renamed route key is still out of range:%s", joinFindings(hits))
		}
		t.Logf("seeded from its own decode and red as required:%s", joinFindings(hits))
	})

	// Plant H: the F-5 shape written into an inbound envelope. The verdict field
	// carries the name-less tag spelling `json:",omitempty"`, which encoding/json
	// resolves to the GO FIELD NAME. The first version of jsonOr called that "no
	// key", so the AST half of facet 2 saw nothing while the reflection half went
	// red on the same tree - two instruments claiming one rule, disagreeing in the
	// direction that hides a field. What is asserted is the finding keyed on
	// "Outcome", the spelling the wire actually uses.
	const plantHHandler = "package panel\n\n" +
		"import \"encoding/json\"\n\n" +
		"// throwaway: a verdict field with no name in its tag.\n" +
		"type grantByFieldName struct {\n" +
		"\tMethod  string `json:\"method\"`\n" +
		"\tOutcome string `json:\",omitempty\"`\n" +
		"}\n\n" +
		"func handleGrantByFieldName(raw string) error {\n" +
		"\tvar e grantByFieldName\n" +
		"\treturn json.Unmarshal([]byte(raw), &e)\n" +
		"}\n"

	t.Run("H a verdict tagged without a name is still found under its field name", func(t *testing.T) {
		dir := copyGoSources(t, src, files, string(bridgeReal), "l2-plant-h.go", plantHHandler)
		pkg, err := parseBoundaryPackage(dir)
		if err != nil {
			t.Fatalf("parse %s: %v", dir, err)
		}
		requireReadableInstrument(t, dir, pkg)
		var got []string
		for _, f := range pkg.structs["grantByFieldName"].fields {
			got = append(got, f.JSONKey)
		}
		sort.Strings(got)
		// "method" is tagged and keeps its spelling; "Outcome" has no name in its
		// tag and falls back to the Go field name. Before the fix this came out as
		// ["", "method"], which is a key that addresses nothing.
		if !equalStrings(got, []string{"Outcome", "method"}) {
			t.Fatalf("the AST key rule resolved the planted struct to %v, want [Outcome method] - jsonOr is not following encoding/json", got)
		}
		hits := scanGrantBoundary(t, dir)
		if !findingsName(hits, "l2-plant-h.go", `"Outcome"`) {
			t.Errorf("plant H is the F-5 shape and produced no finding for the field-name spelling:%s", joinFindings(hits))
		}
		t.Logf("red as required, on the spelling the wire uses:%s", joinFindings(hits))
	})

	t.Run("C a wired grant door goes red on both halves", func(t *testing.T) {
		// The complete violation: the route is answered AND the envelope the
		// handler decodes carries the verdict. This is a throwaway inbound
		// handler that accepts outcome:"grant", plus the guard that lets it in.
		dir := copyGoSources(t, src, files, plantBBridge, "grant_handler.go", plantCHandler)
		hits := scanGrantBoundary(t, dir)
		if len(hits) == 0 {
			t.Fatal("planted an inbound handler decoding {method, outcome} on an answered approval route and NOTHING went red - the nail is decorative")
		}
		if !findingsName(hits, "grant_handler.go", `"outcome"`) {
			t.Errorf("a finding must name grant_handler.go and the outcome key, got:%s", joinFindings(hits))
		}
		if !findingsName(hits, "bridge.go", `"panel.approval.request"`) {
			t.Errorf("a finding must name the answered route, got:%s", joinFindings(hits))
		}
		if _, err := os.Stat(filepath.Join(dir, "grant_handler.go")); err != nil {
			t.Fatalf("the handler plant is not in the scanned tree: %v", err)
		}
		t.Logf("red as required:%s", joinFindings(hits))
	})

	// The scope pin, and the honest answer to "does this catch the historical
	// shape": NO. Plant C's tree above carries only Go; this run carries the
	// verbatim deleted button in a tree whose Go is clean, and the instrument
	// must report nothing - which is exactly why the UI side is registered as
	// zero-instrument coverage rather than called safe. If this subtest ever
	// goes red the instrument got wider than Go, and the coverage claim in this
	// file's header and in the evidence file has to be rewritten.
	t.Run("D the historical JSX button is invisible to this instrument", func(t *testing.T) {
		dir := copyGoSources(t, src, files, string(bridgeReal), "l2-approval-card.tsx", plantJSX)
		hits := scanGrantBoundary(t, dir)
		for _, f := range hits {
			if strings.HasSuffix(f.File, ".tsx") {
				t.Errorf("the Go-boundary instrument now names a .tsx line (%s): the UI-side hole is closed and this file overstates its own blindness", f)
			}
		}
		if len(hits) != 0 {
			t.Errorf("a clean Go snapshot carrying only the historical button must yield 0 findings from this instrument, got:%s", joinFindings(hits))
		}
		if _, err := os.Stat(filepath.Join(dir, "l2-approval-card.tsx")); err != nil {
			t.Fatalf("the JSX plant is not in the scanned tree (%v) - this run would be passing by not showing the hole", err)
		}
		t.Logf("ZERO INSTRUMENT COVERAGE, measured not assumed: a tree holding the verbatim 53a1359^ button (send(\"grant\") -> panel.approval.request carrying outcome:\"grant\") reads CLEAN to this instrument, because it reads Go and the violation is drawn in JSX.")
	})
}

// copyGoSources writes a snapshot of the package's Go files into a fresh temp
// directory, optionally replacing bridge.go and adding one extra file.
func copyGoSources(t *testing.T, src string, files []string, bridgeContent, extraName, extraBody string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "snapshot")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		name := filepath.Base(f)
		if name == "bridge.go" && bridgeContent != "" {
			data = []byte(bridgeContent)
		}
		if err := os.WriteFile(filepath.Join(dir, name), data, 0o600); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	if extraName != "" {
		if err := os.WriteFile(filepath.Join(dir, extraName), []byte(extraBody), 0o600); err != nil {
			t.Fatalf("plant %s: %v", extraName, err)
		}
	}
	return dir
}

func findingsName(hits []boundaryFinding, fileSuffix, needle string) bool {
	for _, f := range hits {
		if strings.HasSuffix(f.File, fileSuffix) && f.Line > 0 && strings.Contains(f.What, needle) {
			return true
		}
	}
	return false
}

func joinFindings(hits []boundaryFinding) string {
	var lines []string
	for _, f := range hits {
		lines = append(lines, "\n  "+f.String())
	}
	return strings.Join(lines, "")
}
