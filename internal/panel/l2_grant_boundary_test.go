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
//	(1) the set of inbound methods Go actually answers - derived from this
//	    package's own route guard, not from a list copied into a test - contains
//	    no approval-shaped name, and agrees with knownComposerMethod at runtime;
//	(2) the envelope type every answered method decodes into cannot BIND a key
//	    whose name carries an approval verdict, recursively through embedded and
//	    nested types - checked by reflection over the production type and by an
//	    AST scan that can name file:line;
//	(3) behaviour: the historical wire shape is refused by ParseComposerRequest,
//	    and for every answered route an envelope that adds outcome:"grant" parses
//	    into a value that no longer holds it;
//	(4) teeth: the same AST instrument, pointed at a COPY of this package with a
//	    grant-carrying envelope, an answered approval route, and a throwaway
//	    handler planted, must name the planted lines.
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
// NOTHING here keys on frontend/src/lib/panel.ts:50's ApprovalOutcome union: that
// union still contains "grant" by the frontend session's explicit choice, so an
// enum-keyed check would be red today and would not be a Go-boundary nail.
// ---------------------------------------------------------------------------

import (
	"bytes"
	"encoding/json"
	"errors"
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
		if tag == "-" {
			continue
		}
		jsonName := ""
		if tag != "" {
			jsonName = strings.Split(tag, ",")[0]
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
	dropped   []unresolvedLabel
	guardFile string
	guardLine int
	guardSeen bool
}

// parseBoundaryPackage reads a directory's struct types, string constants and
// inbound route guard.
func parseBoundaryPackage(dir string) (boundaryPackage, error) {
	pkg := boundaryPackage{structs: map[string]structShape{}, consts: map[string]string{}, answered: map[string]bool{}}
	files, err := goSourceFiles(dir)
	if err != nil {
		return pkg, err
	}
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
					for _, nm := range fld.Names {
						sh.fields = append(sh.fields, astField{
							JSONKey: jsonOr(nm.Name, jsonName, tagged), GoName: nm.Name,
							Line: line, Text: text, TypeName: typeName, NoTag: !tagged,
						})
					}
					if len(fld.Names) == 0 {
						sh.fields = append(sh.fields, astField{
							JSONKey: jsonName, GoName: "embedded:" + typeName,
							Line: line, Text: text, TypeName: typeName, NoTag: !tagged,
						})
					}
				}
				pkg.structs[decl.Name.Name] = sh
				return true
			}
			return true
		})
	}
	return pkg, nil
}

func jsonOr(goName, jsonName string, tagged bool) string {
	if tagged {
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

// guardReadabilityProblem names why this tree's answered-set enumeration cannot
// be trusted, or "" when it can. It is a function rather than a pile of t.Fatal
// calls so the teeth plant in facet 4 can assert the *message* exists without
// aborting its own subtest.
//
// The two ways the enumeration is narrower than the guard it reads are both here:
// no guard at all, and a guard carrying case labels this scan cannot resolve
// (F-1 in docs/evidence/s1/panel-l2-grant-nail-accept-r1.md: a route reaching
// through a package-level var or a concatenation used to be dropped without a
// word, and the file then reported the boundary as clean).
func guardReadabilityProblem(dir string, pkg boundaryPackage) string {
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
	return ""
}

// requireReadableGuard turns guardReadabilityProblem into an abort.
func requireReadableGuard(t *testing.T, dir string, pkg boundaryPackage) {
	t.Helper()
	if problem := guardReadabilityProblem(dir, pkg); problem != "" {
		t.Fatalf("%s", problem)
	}
}

// inboundEnvelopes returns the struct types this package treats as an INBOUND
// envelope, found by structure and not by a name typed in here: a type that
// binds a "method" key is a request the route lives on. Types nested inside one
// are pulled in too, because that is where a smuggled verdict would be parked.
func inboundEnvelopes(pkg boundaryPackage) map[string]bool {
	inbound := map[string]bool{}
	for name, sh := range pkg.structs {
		for _, f := range sh.fields {
			if strings.EqualFold(f.JSONKey, "method") {
				inbound[name] = true
			}
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

// scanGrantBoundary runs the AST instrument over one Go directory and reports
// every grant-carrying face: an inbound envelope field, or an answered route
// whose own name is an approval decision.
func scanGrantBoundary(t *testing.T, dir string) []boundaryFinding {
	t.Helper()
	pkg, err := parseBoundaryPackage(dir)
	if err != nil {
		t.Fatalf("parse %s: %v", dir, err)
	}
	requireReadableGuard(t, dir, pkg)

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
		t.Fatalf("no inbound envelope type (a struct binding a \"method\" key) under %s: the instrument claims to check what a panel request can carry but cannot find the request type", dir)
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
// this package's own route guard, and cross-checked against the function that
// runs, so a door opened in one place cannot be invisible in the other.
func TestAnsweredPanelRoutesCarryNoApprovalDecision(t *testing.T) {
	root := panelRepoRoot(t)
	dir := filepath.Join(root, "internal", "panel")
	pkg, err := parseBoundaryPackage(dir)
	if err != nil {
		t.Fatalf("parse %s: %v", dir, err)
	}
	requireReadableGuard(t, dir, pkg)

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
			t.Errorf("Go answers %q but no Method* constant declares it - a route written straight into the guard, past the naming gate TestComposerMethodNamesMatchFrontend", name)
		}
	}
	t.Logf("answered inbound routes = %v; %d declared Method* constants; 0 approval-shaped names on either side",
		answered, len(declared))
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
		problem := guardReadabilityProblem(dir, pkg)
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
