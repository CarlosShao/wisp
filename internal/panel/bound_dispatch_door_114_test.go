package panel

// Ticket 114 AC#10, leg 114-r6: the reachability verdict must recognise the THIRD
// mounting method.
//
// frontend/src/lib/panel.ts:146 names two ways a WebView2 host hands the page a
// bridge ("AddHostObjectToScript / postMessage pipe") and this repository ships
// neither for the page's outbound door: AddHostObjectToScript has zero call sites,
// and the postMessage hook at cmd/wisp/panel_host_windows.go:781 forwards INTO the
// same door instead of being it. What makes window.wispDispatch exist is the
// WebView2 JS binding - w.Bind(panelDispatchBinding, ...) at :842, with the name
// reviewed once at :80. The reachability scan in composer_test.go counts
// .postMessage( CALL SITES (hostBridgeCallRe) and by construction cannot see that
// third method: a page line that reaches the host as window.wispDispatch is
// invisible to it, and nothing on disk ties that read to the name the host binds.
//
// This file pins the third method as a capability read off the bytes, both ways:
//
//	host side   cmd/wisp's production sources carry a w.Bind(<name>, callback) call
//	            whose name is the value of a package-level string constant; the states
//	            where that call is gone, or its name stopped resolving from the
//	            source, are failures, never quiet passes;
//	page side   the constant's value is a window property that
//	            frontend/src/lib/panel.ts reads, and every read of it under
//	            frontend/ lives in panel.ts (the one door file, ticket 92 AC#1).
//
// The tie is the point: rename the constant's value and the host binds a door no
// page line reaches for; rename the page read and the host binds a door nobody
// reads. Either half alone is dead weight, and that shape is exactly the one both
// existing door scans leave open (door 1 reads the page, door 2 reads the page's
// route literals; neither reads the binding).
//
// Its green is not proven by -overlay: this scan reads the disk (os.ReadFile) and
// -overlay swaps the compiler's eyes, never a reader's (ledger A713). The
// self-proofs below therefore build planted trees from the real bytes just read,
// mutated at offsets this instrument found; the planted half only shows that the
// matching logic bites a dead or renamed door and spares the honest bytes, while
// the reading itself is what the real-tree run exercises - its floor refuses an
// empty read, so a scan that looked at nothing cannot pass (ticket 253-r1 made the
// same split for its own census, cmd/wisp/panel_dispatch_binding_roster_253r1_windows_test.go).

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// windowDoorReadRe captures one property access on the window object - the shape in
// which a page reaches for a host door it was not handed a bridge object for.
var windowDoorReadRe = regexp.MustCompile(`\bwindow\s*\.\s*([A-Za-z_$][A-Za-z0-9_$]*)`)

// doorReadSourceFile is the page-side file set a door read can live in: the types
// vite can load and run. It is deliberately NARROWER than rendererSourceFile, and
// that direction cannot hide a call site any route scan was responsible for: a
// property read reaches the host only from code that executes, and a door name in
// prose (.md, .txt, .json) cannot run.
func doorReadSourceFile(p string) bool {
	switch filepath.Ext(p) {
	case ".ts", ".tsx", ".js", ".jsx", ".mjs", ".vue", ".html":
		return true
	}
	return false
}

// doorReadCommentOpener finds the earliest comment opener on one line and returns
// its index and the closer that ends it. Both page-side comment styles are handled
// because both are load-bearing here: the door name already sits in a /* ... */
// block at panel.ts:329, and index.html can carry <!-- ... -->.
func doorReadCommentOpener(line string) (int, string, string) {
	i, j := strings.Index(line, "/*"), strings.Index(line, "<!--")
	switch {
	case i < 0 && j < 0:
		return -1, "", ""
	case i < 0:
		return j, "<!--", "-->"
	case j < 0:
		return i, "/*", "*/"
	case i < j:
		return i, "/*", "*/"
	default:
		return j, "<!--", "-->"
	}
}

// doorReadCodeLines returns the code half of a page-side source, one string per
// line: // comments through codeOnly, and /* ... */ / <!-- ... --> blocks dropped,
// so a door name spelled in prose (panel.ts:329 today) can never stand in for the
// line that actually reads the door. Strings are not modelled: panel.ts carries no
// comment opener inside a string literal, and a stripper that over-cut would redden
// the honest tree, where it would be noticed rather than trusted.
func doorReadCodeLines(text string) []string {
	var out []string
	closer := ""
	for _, raw := range strings.Split(text, "\n") {
		line := raw
		if closer != "" {
			if i := strings.Index(line, closer); i >= 0 {
				line, closer = line[i+len(closer):], ""
			} else {
				out = append(out, "")
				continue
			}
		}
		for {
			i, openWith, closeWith := doorReadCommentOpener(line)
			if i < 0 {
				break
			}
			rest := line[i+len(openWith):]
			if j := strings.Index(rest, closeWith); j >= 0 {
				line = line[:i] + rest[j+len(closeWith):]
				continue
			}
			line, closer = line[:i], closeWith
			break
		}
		out = append(out, codeOnly(line))
	}
	return out
}

// boundDoorBind is one w.Bind call this scan read as a page-door binding: a name
// (a string literal, or a package-level constant) and a callback. The offsets are
// byte spans the planted cases below mutate, so the plants cut the real bytes at
// the span this instrument itself found.
type boundDoorBind struct {
	rel        string // root-relative, slash form
	site       string // rel:line of the call
	name       string // the door name the call binds; "" when unreadable
	arg        string // the first argument as written
	constBound bool   // the name came through a package-level string constant
	problem    string // non-empty: the argument is not a door name this scan can read
	callStart  int
	callEnd    int
	valueRel   string // the file holding the constant's literal ("" for a literal arg)
	valueStart int
	valueEnd   int
}

// boundDoorRead is one window.<name> property access in a page-side source file.
type boundDoorRead struct {
	site    string // frontend/src/lib/panel.ts:350
	name    string // the property name as read
	inPanel bool   // lives in the one door file
}

// boundDoorReport is one run of the scan.
type boundDoorReport struct {
	hostFiles int
	binds     []boundDoorBind
	skipped   []string
	reads     []boundDoorRead
	problems  []string
}

// reviewed returns the bind sites whose name came through a package-level
// constant. A literal-bound door is not reviewed here on purpose: the repository's
// reviewed act is the named constant (cmd/wisp panel_host_windows.go:791-793 says
// why), and the transient probe door binds a literal so it stays out of this set.
func (r boundDoorReport) reviewed() []boundDoorBind {
	var out []boundDoorBind
	for _, b := range r.binds {
		if b.constBound && b.problem == "" {
			out = append(out, b)
		}
	}
	return out
}

func (r boundDoorReport) pageReads(name string) bool {
	for _, rd := range r.reads {
		if rd.inPanel && rd.name == name {
			return true
		}
	}
	return false
}

func (r boundDoorReport) bindsName(name string) bool {
	for _, b := range r.binds {
		if b.name == name {
			return true
		}
	}
	return false
}

// readSites renders the door-relevant reads for a message: every read in the one
// door file plus every read of a name some Bind binds. The other window.* accesses
// under frontend/ (innerWidth, setTimeout, addEventListener ...) are counted but
// never listed - they are not doors, and a message that drowns the door in them
// names nothing.
func (r boundDoorReport) readSites() string {
	var out []string
	for _, rd := range r.reads {
		if rd.inPanel || r.bindsName(rd.name) {
			out = append(out, rd.site+" window."+rd.name)
		}
	}
	if len(out) == 0 {
		return "(none)"
	}
	return strings.Join(out, "; ")
}

func (r boundDoorReport) bindSites() string {
	var out []string
	for _, b := range r.binds {
		name := b.name
		if name == "" {
			name = "<unreadable>"
		}
		out = append(out, b.site+" binds "+strconv.Quote(name)+" via "+b.arg)
	}
	if len(out) == 0 {
		return "(none)"
	}
	return strings.Join(out, "; ")
}

// findings turns one run into the list of failures, each naming the file:line it
// judged.
func (r boundDoorReport) findings() []string {
	out := append([]string{}, r.problems...)
	for _, b := range r.binds {
		if b.problem != "" {
			out = append(out, b.site+": the w.Bind argument "+b.arg+" "+b.problem+" (ticket 114 AC#10)")
		}
	}
	reviewed := r.reviewed()
	if len(reviewed) == 0 {
		out = append(out, "cmd/wisp's production sources bind no door through a w.Bind(<package-level string constant>, callback) call, so the page's door has no host-side receiver: the third mounting method is what makes window.<door> exist, AddHostObjectToScript has zero call sites in this repository, and the .postMessage( reachability scan cannot see a bound door at all (ticket 114 AC#10). w.Bind sites read: "+
			r.bindSites()+"; page window.* reads: "+r.readSites())
	}
	for _, b := range reviewed {
		if !r.pageReads(b.name) {
			out = append(out, fmt.Sprintf("%s binds %q through w.Bind, and frontend/src/lib/panel.ts never reads window.%s - a bound door no page line reaches for is dead weight, and the two-sided name tie is the whole capability (ticket 114 AC#10). The page's window.* reads: %s",
				b.site, b.name, b.name, r.readSites()))
		}
	}
	for _, rd := range r.reads {
		if !rd.inPanel && r.bindsName(rd.name) {
			out = append(out, rd.site+" reads window."+rd.name+", and every read of a door cmd/wisp binds must live in frontend/src/lib/panel.ts (the one door file, ticket 92 AC#1): a page-side reach of the third mounting method outside that file is a second door")
		}
	}
	return out
}

// goSrc returns the exact source text of one AST node inside text.
func goSrc(text string, fset *token.FileSet, n ast.Node) string {
	start, end := fset.Position(n.Pos()).Offset, fset.Position(n.End()).Offset
	if start < 0 || end < start || end > len(text) {
		return "?"
	}
	return text[start:end]
}

// scanBoundDoor reads one repository root: every production .go of cmd/wisp
// (parsed, with the package-level string constants resolved), and every loadable
// page-side source under frontend/ that accesses a window.<name> property.
func scanBoundDoor(t *testing.T, root string) boundDoorReport {
	t.Helper()
	rep := boundDoorReport{}
	fset := token.NewFileSet()
	type hostFile struct {
		rel  string
		text string
		file *ast.File
	}
	var hosts []hostFile
	consts := map[string]boundDoorBind{}

	err := filepath.WalkDir(filepath.Join(root, "cmd", "wisp"), func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case "testdata", "logs":
				return fs.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".go") || strings.HasSuffix(d.Name(), "_test.go") {
			return nil
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		f, err := parser.ParseFile(fset, p, data, 0)
		if err != nil {
			rep.problems = append(rep.problems, rel+" does not parse ("+err.Error()+"), so the door binding in it cannot be read and the host half of the tie is blind")
			return nil
		}
		hosts = append(hosts, hostFile{rel: rel, text: string(data), file: f})
		for _, decl := range f.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.CONST {
				continue
			}
			for _, spec := range gen.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for i, n := range vs.Names {
					if i >= len(vs.Values) {
						continue
					}
					lit, ok := vs.Values[i].(*ast.BasicLit)
					if !ok || lit.Kind != token.STRING {
						continue
					}
					v, err := strconv.Unquote(lit.Value)
					if err != nil || v == "" {
						continue
					}
					consts[n.Name] = boundDoorBind{
						rel: rel, name: v,
						valueStart: fset.Position(lit.Pos()).Offset,
						valueEnd:   fset.Position(lit.End()).Offset,
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		rep.problems = append(rep.problems, "cmd/wisp is not walkable: "+err.Error())
	}
	rep.hostFiles = len(hosts)

	// Door bind calls: a .Bind whose first argument is a plain name (a literal or
	// an identifier) and which carries a callback too - the WebView2 Bind(name,
	// f) shape. Anything else that happens to be called .Bind (ws.Bind(run.paths),
	// one argument) is recorded as skipped, not judged.
	for _, h := range hosts {
		ast.Inspect(h.file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "Bind" {
				return true
			}
			b := boundDoorBind{
				rel:       h.rel,
				callStart: fset.Position(call.Pos()).Offset,
				callEnd:   fset.Position(call.End()).Offset,
				site:      h.rel + ":" + strconv.Itoa(fset.Position(call.Pos()).Line),
			}
			if len(call.Args) < 2 {
				arg := "(no arguments)"
				if len(call.Args) == 1 {
					arg = goSrc(h.text, fset, call.Args[0])
				}
				rep.skipped = append(rep.skipped, b.site+" .Bind("+arg+") carries one argument, not the name+callback pair a page door needs")
				return true
			}
			b.arg = goSrc(h.text, fset, call.Args[0])
			switch a := call.Args[0].(type) {
			case *ast.BasicLit:
				if a.Kind != token.STRING {
					b.problem = "is not a string literal"
				} else if v, err := strconv.Unquote(a.Value); err != nil || v == "" {
					b.problem = "is a string literal this scan cannot read as a door name"
				} else {
					b.name = v
				}
			case *ast.Ident:
				c, ok := consts[a.Name]
				if !ok || c.name == "" {
					b.problem = "names " + a.Name + ", which no package-level string constant in cmd/wisp declares - the door name is not stated anywhere this scan can read"
				} else {
					b.name, b.constBound = c.name, true
					b.valueRel, b.valueStart, b.valueEnd = c.rel, c.valueStart, c.valueEnd
				}
			default:
				b.problem = "is not a plain name (a string literal or a package-level constant), so the door it binds is not stated anywhere a page-side pin can judge"
			}
			rep.binds = append(rep.binds, b)
			return true
		})
	}

	pageRel := "frontend/src/lib/panel.ts"
	err = filepath.WalkDir(filepath.Join(root, "frontend"), func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case "node_modules", ".git", "dist":
				return fs.SkipDir
			}
			return nil
		}
		if !doorReadSourceFile(p) {
			return nil
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		for i, line := range doorReadCodeLines(string(data)) {
			for _, m := range windowDoorReadRe.FindAllStringSubmatch(line, -1) {
				rep.reads = append(rep.reads, boundDoorRead{site: rel + ":" + strconv.Itoa(i+1), name: m[1], inPanel: rel == pageRel})
			}
		}
		return nil
	})
	if err != nil {
		rep.problems = append(rep.problems, "frontend is not walkable: "+err.Error())
	}
	return rep
}

// TestTheBoundDispatchDoorIsReachableFromThePage is ticket 114 AC#10's capability
// nail for the third mounting method: the door the page reads (window.wispDispatch,
// panel.ts:350) must be the door cmd/wisp binds through w.Bind, by the one reviewed
// constant - and the page-side read must live in the one door file. The planted
// half runs the same scan over trees built from the real bytes with the door
// killed or renamed, so a green above cannot come from a scan that looks at nothing.
func TestTheBoundDispatchDoorIsReachableFromThePage(t *testing.T) {
	root := panelRepoRoot(t)
	rep := scanBoundDoor(t, root)
	for _, f := range rep.findings() {
		t.Error(f)
	}
	if t.Failed() {
		return
	}
	reviewed := rep.reviewed()
	t.Logf("bound door held: %d cmd/wisp production file(s), %d w.Bind door site(s) [%s], %d skipped non-door .Bind call(s); %d window.* read site(s) under frontend/, door-relevant: [%s]",
		rep.hostFiles, len(rep.binds), rep.bindSites(), len(rep.skipped), len(rep.reads), rep.readSites())

	t.Run("a dead or renamed door reddens the same scan", func(t *testing.T) {
		b := reviewed[0]
		readFile := func(rel string) string {
			data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
			if err != nil {
				t.Fatalf("read %s: %v", rel, err)
			}
			return string(data)
		}
		plant := func(files map[string]string) []string {
			dir := t.TempDir()
			for rel, text := range files {
				p := filepath.Join(dir, filepath.FromSlash(rel))
				if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(p, []byte(text), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			return scanBoundDoor(t, dir).findings()
		}
		pageRel := "frontend/src/lib/panel.ts"
		hostText, pageText := readFile(b.rel), readFile(pageRel)
		base := map[string]string{b.rel: hostText, pageRel: pageText}
		if b.valueRel != b.rel {
			base[b.valueRel] = readFile(b.valueRel)
		}
		if got := plant(base); len(got) != 0 {
			t.Errorf("the real bytes, planted unmutated, reddened the scan: %v", got)
		}
		dup := func() map[string]string {
			m := map[string]string{}
			for k, v := range base {
				m[k] = v
			}
			return m
		}

		// (1) the w.Bind call is gone. The span this scan found is replaced by
		// error(nil), so the file still parses and the door is simply never bound.
		gone := dup()
		gone[b.rel] = hostText[:b.callStart] + "error(nil)" + hostText[b.callEnd:]
		got := plant(gone)
		if len(got) == 0 {
			t.Error("with the w.Bind call removed the scan stayed green")
		} else if joined := strings.Join(got, "\n"); !strings.Contains(joined, "bind no door through a w.Bind") {
			t.Errorf("the removed door did not name the missing bind: %v", got)
		}

		// (2) the binding name is renamed at the one reviewed constant: the Bind
		// call and the page are untouched, so only the two-sided name tie can see it.
		valueText := base[b.valueRel]
		renamed := dup()
		renamed[b.valueRel] = valueText[:b.valueStart] + strconv.Quote(b.name+"X") + valueText[b.valueEnd:]
		got = plant(renamed)
		if len(got) == 0 {
			t.Error("with the constant's value renamed the scan stayed green")
		} else if joined := strings.Join(got, "\n"); !strings.Contains(joined, "binds \""+b.name+"X\"") || !strings.Contains(joined, "window."+b.name+"X") {
			t.Errorf("the renamed binding did not name both halves of the tie: %v", got)
		}

		// (3) the page's read is renamed instead: the host still binds b.name and
		// no line under frontend/ reads it any more.
		goneRead := dup()
		goneRead[pageRel] = regexp.MustCompile(`window\s*\.\s*`+regexp.QuoteMeta(b.name)).ReplaceAllString(pageText, "window."+b.name+"X")
		got = plant(goneRead)
		if len(got) == 0 {
			t.Error("with the page's read renamed the scan stayed green")
		} else if joined := strings.Join(got, "\n"); !strings.Contains(joined, "never reads window."+b.name+" -") {
			t.Errorf("the renamed page read did not name the dead door: %v", got)
		}
	})
}
