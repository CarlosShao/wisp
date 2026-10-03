package panel

// Ticket 253 AC#2, shape (a) as ruled in the ticket's section 8: a SECOND roster
// ruler covering the inbound method names the two existing prefix rulers cannot
// see.
//
// What is broken today, stated precisely: internal/panel/git_test.go:394
// (panelMethodRe) and internal/panel/composer_test.go:394 (routeLiteralRe) both
// hard-code the "panel." prefix into their patterns, so the two settings routes of
// bridge.go - "config.get" and "config.set" - fall outside both rosters, and a new
// inbound method name can enter the tree without tripping a single existing ruler
// (ticket 253 section 1 item 2, ledger A519). This file widens neither of those
// rulers and changes nothing about bridge.go's behaviour: it adds one closed
// roster read off the production source by STRUCTURE, never by a namespace prefix,
// because a second prefix-bound pattern would only copy the defect.
//
// What this ruler deliberately does NOT do: it asserts nothing about the C17
// four-method contract (git_test.go:385 and :517 stay the gate on that surface,
// and this leg leaves both alone), and it reads nothing outside internal/panel's
// production sources - cmd/wisp's inbound surface belongs to that package's legs.
//
// Denominators, four independent signals, so one narrowed set cannot hide behind
// another:
//
//   - DECLARED: every package-level string constant of bridge.go whose identifier
//     begins with "Method" OR whose VALUE is a dotted name. The second clause is
//     what makes this prefix-agnostic: "config.peek", "settings.read" and "x.y"
//     are all in scope, and no namespace is ever consulted.
//   - GUARD: the case labels of bridge.go's knownComposerMethod, resolved through
//     that file's constants.
//   - ROUTER: the case labels of composer_dispatch.go's dispatch, resolved through
//     the package's constants - the router is where a name becomes an action, so a
//     name it handles belongs on the roster.
//   - ANSWERED: every dotted literal in the package's production files that the
//     RUNNING guard answers (knownComposerMethod is asked, never re-implemented).
//
// Each must be the same closed set, checked in BOTH directions. The set is written
// down here as string literals, not as references to the Go constants: a roster
// that re-lists the constants it is auditing stays green no matter what gets added
// (the reasoning git_test.go:403-406 already gives for its own extractor).

import (
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

const (
	// wantFullInboundRosterSize253 and wantNonPanelPrefixedInbound253 are the two
	// numbers a diff has to move on purpose. Pinning them is what stops this file
	// from being satisfied by a silently narrowed roster - and, together with the
	// empty-source case below, by an extractor that found nothing.
	wantFullInboundRosterSize253 = 6
	// wantNonPanelPrefixedInbound253 is AC#2's target count: the names carrying no
	// "panel." prefix, invisible to both old rulers. It is 2 today ("config.get",
	// "config.set"); a third settings-shaped route has to move this number, which
	// is the friction the ticket asked for.
	wantNonPanelPrefixedInbound253 = 2
)

// wantFullInboundRoster253 is the closed roster: every name a page can put in the
// "method" field of a composer envelope that this package's inbound door answers.
var wantFullInboundRoster253 = []string{
	"config.get",
	"config.set",
	"panel.attachment.add",
	"panel.message.send",
	"panel.mode.request",
	"panel.workspace.request",
}

// inboundDotShaped reports whether s is spelled like a namespaced name: two or
// more segments, each at least one character of letters/digits/underscore, a
// hyphen allowed only away from the first position. It asks nothing about WHICH
// namespace - that is the difference between this ruler and the two prefix-bound
// patterns whose reach it replaces. Paths, mimes, printf verbs, struct tags and
// the sender identity ("panel-composer", no dot) fail the shape.
func inboundDotShaped(s string) bool {
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

// inboundDecl is one package-level string constant, kept with the two independent
// reasons it might be an inbound method name and the line it was written on.
type inboundDecl struct {
	ident            string
	value            string
	methodNamedIdent bool
	dottedValue      bool
	line             int
}

// inboundName returns the roster name this declaration contributes, or "" when it
// is an inbound method name by neither signal.
func (d inboundDecl) inboundName() string {
	if d.methodNamedIdent || d.dottedValue {
		return d.value
	}
	return ""
}

// fileInbound is one parsed production file: its package-level string constants
// and the inbound names among them.
type fileInbound struct {
	path       string
	consts     map[string]string
	inbound    []inboundDecl
	parseError string
}

// parseFileInbound reads one Go file and collects its package-level string
// constants plus the inbound names they declare. Only top-level const declarations
// count: a constant inside a function body is not part of any inbound surface.
func parseFileInbound(t *testing.T, path string) fileInbound {
	t.Helper()
	f := fileInbound{path: path, consts: map[string]string{}}
	src, err := os.ReadFile(path)
	if err != nil {
		f.parseError = "read " + path + ": " + err.Error()
		return f
	}
	fset := token.NewFileSet()
	af, perr := parser.ParseFile(fset, filepath.Base(path), src, 0)
	if perr != nil {
		f.parseError = "parse " + path + ": " + perr.Error()
		return f
	}
	for _, decl := range af.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.CONST {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok || len(vs.Values) == 0 {
				continue
			}
			for i, id := range vs.Names {
				if i >= len(vs.Values) {
					continue
				}
				lit, ok := vs.Values[i].(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					continue
				}
				val, uerr := strconv.Unquote(lit.Value)
				if uerr != nil {
					continue
				}
				f.consts[id.Name] = val
				f.inbound = append(f.inbound, inboundDecl{
					ident:            id.Name,
					value:            val,
					methodNamedIdent: strings.HasPrefix(id.Name, "Method"),
					dottedValue:      inboundDotShaped(val),
					line:             fset.Position(lit.Pos()).Line,
				})
			}
		}
	}
	return f
}

// namesOf returns the sorted, deduplicated inbound names of one parsed file.
func (f fileInbound) namesOf() []string {
	seen := map[string]bool{}
	var out []string
	for _, d := range f.inbound {
		n := d.inboundName()
		if n == "" || seen[n] {
			continue
		}
		seen[n] = true
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

// siteOf names where a declaration was written, so a failure can point.
func (f fileInbound) siteOf(name string) string {
	for _, d := range f.inbound {
		if d.inboundName() == name {
			return filepath.Base(f.path) + ":" + strconv.Itoa(d.line) + " const " + d.ident
		}
	}
	return filepath.Base(f.path)
}

// resolveCaseLabel turns one case label expression into the route value it
// carries. It reports a problem instead of staying quiet: an unresolvable label is
// exactly the shape this ruler exists to catch, and dropping it would rebuild the
// blind spot AC#2 was filed over.
func resolveCaseLabel(label ast.Expr, constMaps []map[string]string) (value string, problem string) {
	switch e := label.(type) {
	case *ast.BasicLit:
		if e.Kind != token.STRING {
			return "", "case label " + e.Value + " is a non-string literal"
		}
		v, uerr := strconv.Unquote(e.Value)
		if uerr != nil {
			return "", "case label " + e.Value + " is not an unquotable string literal"
		}
		return v, ""
	case *ast.Ident:
		for _, m := range constMaps {
			if v, ok := m[e.Name]; ok {
				return v, ""
			}
		}
		return "", "case label " + e.Name + " names no package-level string constant of internal/panel: a route spelled out only here"
	default:
		return "", "case label is neither an identifier nor a string literal, so this ruler cannot read which route it answers"
	}
}

// switchCaseValues walks the switch statements of one named function and returns
// the value each case label carries, resolved through the const maps it is given
// (the file's own first, then the package's, so a label may name a constant
// declared in another file).
//
// wantDefaultCall decides what a default branch means: the empty string says this
// function must have none (knownComposerMethod with a default answers every route
// its labels never name), while a name says the function must have exactly one and
// its body must call that function - which is how composer_dispatch.go's
// rosterMismatch backstop stays required instead of being read as a hole. An
// unresolvable case label always comes back as a problem.
func switchCaseValues(t *testing.T, path string, funcName string, wantDefaultCall string, constMaps ...map[string]string) ([]string, []string) {
	t.Helper()
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, []string{"read " + path + ": " + err.Error()}
	}
	fset := token.NewFileSet()
	af, perr := parser.ParseFile(fset, filepath.Base(path), src, 0)
	if perr != nil {
		return nil, []string{"parse " + path + ": " + perr.Error()}
	}
	var fn *ast.FuncDecl
	for _, decl := range af.Decls {
		if d, ok := decl.(*ast.FuncDecl); ok && d.Name.Name == funcName {
			fn = d
			break
		}
	}
	if fn == nil {
		return nil, []string{filepath.Base(path) + " declares no func " + funcName +
			": the inbound surface this ruler reads was renamed, so it is now judging a shape Go does not run"}
	}
	var values, problems []string
	var defaults []string
	seen := map[string]bool{}
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
			at := filepath.Base(path) + ":" + strconv.Itoa(fset.Position(cs.Pos()).Line)
			if len(cs.List) == 0 {
				site := at + " func " + funcName + " default branch"
				if wantDefaultCall == "" {
					problems = append(problems, site+": it answers every route its case labels do not name")
					continue
				}
				if !callsNamed(cs.Body, wantDefaultCall) {
					problems = append(problems, site+": the backstop that refuses an unrostered route (it must call "+
						wantDefaultCall+") is not what runs here")
					continue
				}
				defaults = append(defaults, site+" refuses through "+wantDefaultCall)
				continue
			}
			for _, e := range cs.List {
				v, problem := resolveCaseLabel(e, constMaps)
				if problem != "" {
					problems = append(problems, at+" "+problem)
					continue
				}
				if !seen[v] {
					seen[v] = true
					values = append(values, v)
				}
			}
		}
		return true
	})
	if wantDefaultCall != "" && len(defaults) == 0 {
		problems = append(problems, "func "+funcName+" has no default branch at all: a route off this roster would reach a case that does not exist and be handled by nothing")
	}
	sort.Strings(values)
	return values, problems
}

// callsNamed reports whether a statement list calls a function of this name, as a
// plain call or a method call on a receiver.
func callsNamed(body []ast.Stmt, name string) bool {
	found := false
	for _, s := range body {
		ast.Inspect(s, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			switch fn := call.Fun.(type) {
			case *ast.Ident:
				if fn.Name == name {
					found = true
				}
			case *ast.SelectorExpr:
				if fn.Sel != nil && fn.Sel.Name == name {
					found = true
				}
			}
			return true
		})
	}
	return found
}

// productionPanelFiles lists internal/panel's production Go files, skipping
// _test.go the way internal/panel/l2_grant_boundary_test.go:369 does. That skip is
// load-bearing for this ticket in both directions: it is why adding THIS file
// cannot move any existing ruler's denominator (253-p4 section 2, shape (a)), and
// why the roster below audits the shipped surface rather than test scaffolding.
func productionPanelFiles(t *testing.T, dir string) []string {
	t.Helper()
	var out []string
	err := filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
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
	if err != nil {
		t.Fatalf("walk %s: %v", dir, err)
	}
	sort.Strings(out)
	if len(out) == 0 {
		t.Fatalf("no production Go files under %s - this roster would be judging an empty universe", dir)
	}
	return out
}

// sameNameSet reports the two directions of a set difference. Both are needed: the
// one-directional subset form is the relaxed shape
// TestSubsetOnlyRosterDirectionIsBlindToAPlantedName measures as blind.
func sameNameSet(got, want []string) (missing, extra []string) {
	inWant := map[string]bool{}
	for _, w := range want {
		inWant[w] = true
	}
	inGot := map[string]bool{}
	for _, g := range got {
		inGot[g] = true
	}
	for _, w := range want {
		if !inGot[w] {
			missing = append(missing, w)
		}
	}
	for _, g := range got {
		if !inWant[g] {
			extra = append(extra, g)
		}
	}
	sort.Strings(missing)
	sort.Strings(extra)
	return missing, extra
}

// quoteNames formats a name list for a failure message.
func quoteNames(names []string) string {
	if len(names) == 0 {
		return "(none)"
	}
	quoted := make([]string, 0, len(names))
	for _, n := range names {
		quoted = append(quoted, strconv.Quote(n))
	}
	return strings.Join(quoted, ", ")
}

// TestFullInboundMethodRosterIsClosed is AC#2's ruler of record: one closed
// roster, five denominators over the real tree, both directions on each one.
func TestFullInboundMethodRosterIsClosed(t *testing.T) {
	root := panelRepoRoot(t)
	dir := filepath.Join(root, "internal", "panel")
	bridgePath := filepath.Join(dir, "bridge.go")
	dispatchPath := filepath.Join(dir, "composer_dispatch.go")

	bridge := parseFileInbound(t, bridgePath)
	if bridge.parseError != "" {
		t.Fatalf("bridge.go: %s", bridge.parseError)
	}
	declared := bridge.namesOf()
	if len(declared) == 0 {
		t.Fatalf("bridge.go declares no inbound method name at all - an extractor that finds nothing must never read as a pass")
	}

	// Denominator 1 (DECLARED): bridge.go's inbound names are exactly the roster.
	missing, extra := sameNameSet(declared, wantFullInboundRoster253)
	for _, name := range extra {
		t.Errorf("ticket 253 AC#2: bridge.go declares the inbound method name %s at %s and the closed roster does not carry it - a new inbound route entered the tree without being registered anywhere. The two prefix rulers cannot see a name that carries no \"panel.\" prefix (panelMethodRe at git_test.go:394, routeLiteralRe at composer_test.go:394), which is the hole this roster exists to close. If the name is meant to stay, register it in wantFullInboundRoster253 and move the count constants beside it",
			name, bridge.siteOf(name))
	}
	for _, name := range missing {
		t.Errorf("ticket 253 AC#2: the closed roster carries %s but bridge.go no longer declares it - the roster and the inbound surface are two different universes, so every check below is auditing a name that may not exist", name)
	}

	// The two written-down numbers.
	if len(wantFullInboundRoster253) != wantFullInboundRosterSize253 {
		t.Errorf("wantFullInboundRoster253 holds %d names, want %d - the roster grew or shrank without its count constant being moved with it",
			len(wantFullInboundRoster253), wantFullInboundRosterSize253)
	}
	var unprefixed []string
	for _, name := range wantFullInboundRoster253 {
		if !strings.HasPrefix(name, "panel.") {
			unprefixed = append(unprefixed, name)
		}
	}
	if len(unprefixed) != wantNonPanelPrefixedInbound253 {
		t.Errorf("the closed roster holds %d names with no \"panel.\" prefix, want %d: that batch is exactly what AC#2 exists to cover (%s). If the count moved, either the prefix-blind rulers got blinder relative to this one or this roster grew, and the diff has to say which",
			len(unprefixed), wantNonPanelPrefixedInbound253, quoteNames(unprefixed))
	}

	// Denominator 2 (RUNNING GUARD): every roster name is answered by Go itself,
	// and the guard's refusals stay refusals. knownComposerMethod is asked here,
	// never re-implemented.
	for _, name := range wantFullInboundRoster253 {
		if !knownComposerMethod(name) {
			t.Errorf("ticket 253 AC#2: the roster carries %s but the running guard knownComposerMethod refuses it - a rostered inbound name Go does not answer is a page-visible dead door", name)
		}
	}

	// Denominator 3 (GUARD AST): the guard's own case labels are the roster.
	guardLabels, guardProblems := switchCaseValues(t, bridgePath, "knownComposerMethod", "", bridge.consts)
	for _, p := range guardProblems {
		t.Errorf("ticket 253 AC#2: %s", p)
	}
	miss, ext := sameNameSet(guardLabels, wantFullInboundRoster253)
	for _, name := range ext {
		t.Errorf("ticket 253 AC#2: knownComposerMethod's case list answers %s and the closed roster does not carry it - an answered route that was never registered", name)
	}
	for _, name := range miss {
		t.Errorf("ticket 253 AC#2: the roster carries %s but knownComposerMethod's case list no longer names it - the guard stopped answering a door this file still claims is open", name)
	}

	// Package view of the constants, for the router's labels.
	pkgConsts := map[string]string{}
	files := productionPanelFiles(t, dir)
	for _, p := range files {
		fi := parseFileInbound(t, p)
		if fi.parseError != "" {
			t.Fatalf("%s", fi.parseError)
		}
		pkgConsts = mergeConsts(pkgConsts, fi.consts)
	}

	// Denominator 4 (ROUTER): what the router routes is on the roster too.
	routerLabels, routerProblems := switchCaseValues(t, dispatchPath, "dispatch", "rosterMismatch", bridge.consts, pkgConsts)
	for _, p := range routerProblems {
		t.Errorf("ticket 253 AC#2: %s", p)
	}
	miss, ext = sameNameSet(routerLabels, wantFullInboundRoster253)
	for _, name := range ext {
		t.Errorf("ticket 253 AC#2: composer_dispatch.go's dispatch routes %s and the closed roster does not carry it - the router is where an inbound name becomes an action, so a name it handles belongs on this roster (register it, or take the case out)", name)
	}
	for _, name := range miss {
		t.Errorf("ticket 253 AC#2: the roster carries %s but composer_dispatch.go's dispatch has no case for it - the guard would answer it and the router would fall through to the roster-mismatch backstop", name)
	}

	// Denominator 5 (ANSWERED literals): any dotted literal in the package's
	// production sources that the RUNNING guard answers must be a roster name. This
	// is the direction that catches a route written straight into an expression,
	// and it asks no namespace.
	offRoster := map[string][]string{}
	for _, p := range files {
		offRoster = collectAnsweredOffRoster(t, p, offRoster)
	}
	for _, name := range sortedKeysOf(offRoster) {
		t.Errorf("ticket 253 AC#2: the running guard answers the dotted literal %q found at %s, and it is on no roster - a route that exists only inside an expression has no registerable spelling, which is how an inbound door gets shipped unnumbered",
			name, strings.Join(offRoster[name], ", "))
	}

	t.Logf("full inbound roster held: %d names declared in bridge.go, %d guard case labels, %d router case labels, %d of them carrying no \"panel.\" prefix (%s)",
		len(declared), len(guardLabels), len(routerLabels), len(unprefixed), quoteNames(unprefixed))
}

// mergeConsts folds one file's constants into the package view and refuses to
// resolve a name two files declare with different values: an ambiguous spelling
// must not be silently settled by a roster check.
func mergeConsts(dst, src map[string]string) map[string]string {
	for k, v := range src {
		if old, ok := dst[k]; ok && old != v {
			dst[k] = "\x00ambiguous:" + old + "/" + v
			continue
		}
		if _, ok := dst[k]; !ok {
			dst[k] = v
		}
	}
	return dst
}

// collectAnsweredOffRoster records, per file, the dotted string literals that the
// running guard answers while the closed roster does not carry them.
func collectAnsweredOffRoster(t *testing.T, path string, acc map[string][]string) map[string][]string {
	t.Helper()
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	fset := token.NewFileSet()
	af, perr := parser.ParseFile(fset, filepath.Base(path), src, 0)
	if perr != nil {
		t.Fatalf("parse %s: %v", path, perr)
	}
	onRoster := map[string]bool{}
	for _, name := range wantFullInboundRoster253 {
		onRoster[name] = true
	}
	ast.Inspect(af, func(n ast.Node) bool {
		lit, ok := n.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}
		v, uerr := strconv.Unquote(lit.Value)
		if uerr != nil || !inboundDotShaped(v) {
			return true
		}
		if !knownComposerMethod(v) {
			return true
		}
		if onRoster[v] {
			return true
		}
		acc[v] = append(acc[v], filepath.Base(path)+":"+strconv.Itoa(fset.Position(lit.Pos()).Line))
		return true
	})
	return acc
}

// sortedKeysOf makes a map's report deterministic.
func sortedKeysOf(m map[string][]string) []string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// TestPlantedUnregisteredInboundMethodNameGoesRed is AC#2's positive control, in
// the judgement shape the ticket writes down: plant a new inbound method name that
// is not on the roster and the ruler must NAME the step that went wrong. Every
// plant goes into a t.TempDir copy, so the real tree is only ever read here, and
// the clean half at the end is the unplanted reading under the same helper.
//
// The plants are chosen so that they trip nothing that is not this ruler: no plant
// carries a "panel." prefix, so git_test.go:385 and :517 - the four-method C17
// contract anchors, one line of which this leg does not touch - stay green by
// construction, and composer_test.go never looks at a tree copy of this file.
func TestPlantedUnregisteredInboundMethodNameGoesRed(t *testing.T) {
	root := panelRepoRoot(t)
	realBridge := filepath.Join(root, "internal", "panel", "bridge.go")
	src, err := os.ReadFile(realBridge)
	if err != nil {
		t.Fatalf("read bridge.go: %v", err)
	}
	unplanted := parseFileInbound(t, realBridge)
	if unplanted.parseError != "" {
		t.Fatalf("real bridge.go: %s", unplanted.parseError)
	}

	cases := []struct {
		name  string
		plant string
		want  string
	}{
		{
			// Plant A: a settings-shaped route whose identifier does NOT begin with
			// "Method" - the shape only a prefix-agnostic and identifier-agnostic
			// ruler sees.
			name:  "plain-identifier-const",
			plant: "\nconst ConfigPeek253Test = \"config.peek253\"\n",
			want:  "config.peek253",
		},
		{
			// Plant B: the same drift with a Method-shaped identifier, still with no
			// "panel." prefix - what the old git anchor misses.
			name:  "Method-named-const",
			plant: "\nconst MethodConfigFoo253Test = \"config.foo253\"\n",
			want:  "config.foo253",
		},
	}
	for _, tc := range cases {
		dir := t.TempDir()
		path := filepath.Join(dir, "bridge.go")
		if err := os.WriteFile(path, append(append([]byte{}, src...), []byte(tc.plant)...), 0o644); err != nil {
			t.Fatalf("plant %s: %v", tc.name, err)
		}
		planted := parseFileInbound(t, path)
		if planted.parseError != "" {
			t.Fatalf("planted bridge.go (%s): %s", tc.name, planted.parseError)
		}
		miss, ext := sameNameSet(planted.namesOf(), wantFullInboundRoster253)
		if len(ext) != 1 || ext[0] != tc.want {
			t.Errorf("planted an unregistered inbound name (%s) and the roster ruler reported extra=%s, want exactly [%s] - this ruler has no teeth for the shape AC#2 was filed over",
				tc.name, quoteNames(ext), strconv.Quote(tc.want))
		}
		if len(miss) != 0 {
			t.Errorf("plant %s moved the missing side too (%s): the plant was supposed to ADD a name, not replace one",
				tc.name, quoteNames(miss))
		}
		if got := planted.siteOf(tc.want); !strings.Contains(got, "const ") {
			t.Errorf("plant %s: the failure message would point at %q, which names no declaration - a red that cannot say where is not a usable verdict", tc.name, got)
		}
	}

	// Plant C: a route answered straight from a literal in the guard's case list -
	// the shape that never became a declarable, registerable name at all. Written
	// into a copy by replacing that copy's single case line.
	dir := t.TempDir()
	path := filepath.Join(dir, "bridge.go")
	guardLine := "\tcase MethodModeRequest, MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend, MethodConfigGet, MethodConfigSet:"
	if !strings.Contains(string(src), guardLine) {
		t.Fatalf("bridge.go's guard case list is no longer written as one comma-separated clause - re-read knownComposerMethod and rebuild this plant from what is on disk rather than replacing nothing")
	}
	zapLine := strings.TrimSuffix(guardLine, ":") + `, "config.zap253":`
	if err := os.WriteFile(path, []byte(strings.Replace(string(src), guardLine, zapLine, 1)), 0o644); err != nil {
		t.Fatalf("plant guard literal: %v", err)
	}
	plantedBridge := parseFileInbound(t, path)
	if plantedBridge.parseError != "" {
		t.Fatalf("planted bridge.go (guard literal): %s", plantedBridge.parseError)
	}
	labels, problems := switchCaseValues(t, path, "knownComposerMethod", "", plantedBridge.consts)
	if len(problems) != 0 {
		t.Errorf("the copied guard plant also produced structural problems %s; this case is about an answered name off the roster and should stay readable", quoteNames(problems))
	}
	_, ext := sameNameSet(labels, wantFullInboundRoster253)
	if len(ext) != 1 || ext[0] != "config.zap253" {
		t.Errorf("planted a route answered straight from a case-list literal and the roster ruler reported extra=%s, want exactly [\"config.zap253\"] - a name written only in a switch is precisely what a prefix ruler cannot see", quoteNames(ext))
	}

	// The honest half under the same helper: the real tree adds nothing.
	miss, ext := sameNameSet(unplanted.namesOf(), wantFullInboundRoster253)
	if len(miss) != 0 || len(ext) != 0 {
		t.Errorf("the plants only mean something if the real tree reads quiet under the same helper; got missing=%s extra=%s",
			quoteNames(miss), quoteNames(ext))
	}
	t.Logf("AC#2's judgement shape holds under this ruler: an unregistered inbound name is named in all three planted shapes (plain const, Method-named const, guard case literal) while the unplanted tree reads quiet")
}

// TestRegistrationSilencesTheRosterRuler measures AC#2's other half - "the name got
// registered, so the ruler goes quiet" - in a tree copy, where a full registration
// can actually be written. It matters which registration is meant: a name added to
// wantFullInboundRoster253 alone does NOT silence this file, because the roster
// claims the guard answers every name on it. Silence requires the whole pair - the
// declaration and the route the guard answers - which is the point of the two
// count constants above.
func TestRegistrationSilencesTheRosterRuler(t *testing.T) {
	root := panelRepoRoot(t)
	src, err := os.ReadFile(filepath.Join(root, "internal", "panel", "bridge.go"))
	if err != nil {
		t.Fatalf("read bridge.go: %v", err)
	}
	constLine := "\nconst ConfigRegistered253Test = \"config.registered253\"\n"
	guardLine := "\tcase MethodModeRequest, MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend, MethodConfigGet, MethodConfigSet:"
	if !strings.Contains(string(src), guardLine) {
		t.Fatalf("bridge.go's guard case list is not the one comma-separated clause this plant was written against - re-read knownComposerMethod rather than trusting this line")
	}
	registeredLine := strings.TrimSuffix(guardLine, ":") + ", ConfigRegistered253Test:"
	roster := append(append([]string{}, wantFullInboundRoster253...), "config.registered253")

	// Case A: fully registered - declared, answered by the guard, on the roster.
	full := t.TempDir()
	fullPath := filepath.Join(full, "bridge.go")
	both := strings.Replace(string(src)+constLine, guardLine, registeredLine, 1)
	if err := os.WriteFile(fullPath, []byte(both), 0o644); err != nil {
		t.Fatalf("write registered copy: %v", err)
	}
	regFile := parseFileInbound(t, fullPath)
	if regFile.parseError != "" {
		t.Fatalf("registered copy: %s", regFile.parseError)
	}
	miss, ext := sameNameSet(regFile.namesOf(), roster)
	if len(miss) != 0 || len(ext) != 0 {
		t.Errorf("a fully registered inbound name still reads missing=%s extra=%s - the roster ruler would then be shouting at a change that did go through every surface, and shouting at honest work is how a ruler gets muted",
			quoteNames(miss), quoteNames(ext))
	}
	labels, problems := switchCaseValues(t, fullPath, "knownComposerMethod", "", regFile.consts)
	if len(problems) != 0 {
		t.Errorf("the registered copy's guard produced structural problems %s, want none", quoteNames(problems))
	}
	miss, ext = sameNameSet(labels, roster)
	if len(miss) != 0 || len(ext) != 0 {
		t.Errorf("the registered copy's guard case list disagrees with the registered roster (missing=%s extra=%s)", quoteNames(miss), quoteNames(ext))
	}

	// Case B: roster-only registration - the name is on the roster and declared,
	// but nothing answers it. That must stay red, and it must say so.
	partial := t.TempDir()
	partialPath := filepath.Join(partial, "bridge.go")
	if err := os.WriteFile(partialPath, []byte(string(src)+constLine), 0o644); err != nil {
		t.Fatalf("write declared-only copy: %v", err)
	}
	declaredOnly := parseFileInbound(t, partialPath)
	if declaredOnly.parseError != "" {
		t.Fatalf("declared-only copy: %s", declaredOnly.parseError)
	}
	miss, ext = sameNameSet(declaredOnly.namesOf(), roster)
	if len(miss) != 0 || len(ext) != 0 {
		t.Errorf("declared-only registration should already agree on the DECLARED side (missing=%s extra=%s)", quoteNames(miss), quoteNames(ext))
	}
	labels, problems = switchCaseValues(t, partialPath, "knownComposerMethod", "", declaredOnly.consts)
	if len(problems) != 0 {
		t.Errorf("the declared-only copy's guard produced structural problems %s, want none", quoteNames(problems))
	}
	miss, _ = sameNameSet(labels, roster)
	if len(miss) != 1 || miss[0] != "config.registered253" {
		t.Errorf("a name on the roster that nothing answers must read as a dead door: got missing=%s, want [\"config.registered253\"]", quoteNames(miss))
	}
	if knownComposerMethod("config.registered253") {
		t.Errorf("the running guard answers config.registered253, which would mean a test-tree name reached production - this measurement is not measuring what it claims")
	}
	t.Logf("registration read both ways: declared AND answered reads quiet (%d names, no diff); roster-only reads red on the answered side (%s), so adding a name to the list is not how this ruler gets satisfied",
		len(roster), quoteNames(miss))
}

// TestFullInboundRosterRefusesAnEmptyDeclarationFile is the tautology guard on the
// other side of the same coin: a ruler that accepted "found nothing" would be green
// on any tree. The decoy parses, declares no string constant, and must read as
// every roster name missing.
func TestFullInboundRosterRefusesAnEmptyDeclarationFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bridge.go")
	decoy := "package panel\n\nconst requestIDBytes253Test = 8\n\nfunc knownComposerMethod(m string) bool {\n\tswitch m {\n\tcase \"config.peek253\":\n\t\treturn true\n\t}\n\treturn false\n}\n"
	if err := os.WriteFile(path, []byte(decoy), 0o644); err != nil {
		t.Fatal(err)
	}
	empty := parseFileInbound(t, path)
	if empty.parseError != "" {
		t.Fatalf("decoy file did not parse: %s", empty.parseError)
	}
	if got := empty.namesOf(); len(got) != 0 {
		t.Fatalf("the decoy declares no inbound name and the extractor found %s - it is reading something other than the file it was given", quoteNames(got))
	}
	miss, ext := sameNameSet(empty.namesOf(), wantFullInboundRoster253)
	if len(ext) != 0 {
		t.Fatalf("sameNameSet reported extra=%s for an empty source; the missing side is the one that must fire", quoteNames(ext))
	}
	if len(miss) != wantFullInboundRosterSize253 {
		t.Errorf("an empty declaration file produced %d missing names, want all %d: if this ruler does not go fully red when the source holds nothing, its green says nothing about the tree",
			len(miss), wantFullInboundRosterSize253)
	}
	// The same decoy's guard: a literal-only case list still has to be read, so a
	// name only the switch knows about comes back rather than silence.
	labels, problems := switchCaseValues(t, path, "knownComposerMethod", "", empty.consts)
	if len(problems) != 0 {
		t.Errorf("the decoy's guard produced structural problems %s, want none - its case label is a plain string literal", quoteNames(problems))
	}
	if len(labels) != 1 || labels[0] != "config.peek253" {
		t.Errorf("reading the decoy's guard gave %s, want exactly [\"config.peek253\"]", quoteNames(labels))
	}
	t.Logf("empty-source refusal held: %d roster names reported missing against a file that declares none", len(miss))
}

// TestSubsetOnlyRosterDirectionIsBlindToAPlantedName answers the tautology question
// in the open instead of leaving it to a reader's judgement. The relaxed form of
// AC#2's assertion - "every roster name is present in the source", subset one way
// only - is run against the plant the ticket names, and the measurement is that it
// has nothing to say. That is what makes this file's credential the two-directional
// equality, and it must never be quietly swapped for a subset check.
func TestSubsetOnlyRosterDirectionIsBlindToAPlantedName(t *testing.T) {
	root := panelRepoRoot(t)
	src, err := os.ReadFile(filepath.Join(root, "internal", "panel", "bridge.go"))
	if err != nil {
		t.Fatalf("read bridge.go: %v", err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "bridge.go")
	plant := "\nconst ConfigBlindnessProbe253Test = \"config.blind253\"\n"
	if err := os.WriteFile(path, append(append([]byte{}, src...), []byte(plant)...), 0o644); err != nil {
		t.Fatal(err)
	}
	planted := parseFileInbound(t, path)
	if planted.parseError != "" {
		t.Fatalf("planted bridge.go: %s", planted.parseError)
	}
	miss, ext := sameNameSet(planted.namesOf(), wantFullInboundRoster253)
	if len(miss) != 0 {
		t.Fatalf("the planted copy reports roster names missing (%s), which means this measurement did not carry the honest half along", quoteNames(miss))
	}
	if len(ext) != 1 || ext[0] != "config.blind253" {
		t.Errorf("the strict form did not name the planted name: extra=%s, want [\"config.blind253\"] - the form being measured as blind only matters while the shipped form has teeth", quoteNames(ext))
	}
	t.Logf("relaxed form measured blind on purpose: with the plant in the source the subset-only reading finds no problem (missing=%s) while the shipped equality names %s. Subset-only is not admissible as AC#2's credential",
		quoteNames(miss), quoteNames(ext))
}
