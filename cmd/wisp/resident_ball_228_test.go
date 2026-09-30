package main

// Ticket 228 AC#1: the floating ball must be hosted by a production leg, not
// merely compiled.
//
// The defect this file exists for is stated in the ticket as a reading, not an
// opinion: before ticket 228 the ruler
//
//	grep -rln "wisp/internal/ball" --include=*.go . | grep -v .scratch
//
// returned exactly one file, cmd/balldebug/main.go - a debug harness. Every
// ball feature that "exists" today therefore exists in a process nobody
// launches. AC#1's own wording forbids the cheap answer: "不许用能编译充当我
// 有人起它". So the claim checked here is a walk: entry point -> host call ->
// ball.New, read out of the source of this package's non-test files.
//
// Why a source walk and not a runtime call: runResident() never returns (it
// parks in the event loop and may os.Exit), so an in-process call would mean a
// test reimplementing its surroundings - the judgment ticket 127's nail file
// already made for the same leg (resident_sink_nail_127_windows_test.go, header).
// What the runtime side of the claim is worth is nailed there and in the two
// companion files: the subprocess case reads the sentence and the log records
// the shipped no-args process produces, and the winlive case looks for the real
// window. This file's job is the shape: delete the call site, or drop a
// gesture's executor, and it goes red by name on every platform - which matters
// because the code it walks is //go:build windows and the POSIX half of this
// repository can otherwise only measure instruments like this one.

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// ballImportPath228 is the AC#1 ruler's subject, spelled the way an import
// block spells it.
const ballImportPath228 = "github.com/CarlosShao/wisp/internal/ball"

// ballEventCallbacks228 names the consumer surface of internal/ball. It is a
// literal list on purpose: when internal/ball grows an eleventh gesture, the
// case below goes red and whoever added it has to say what the resident leg
// does with it. A nil Events entry is a gesture that vanishes without a word,
// which is the same silence this family keeps counting.
var ballEventCallbacks228 = []string{
	"OnClickBall", "OnSummonHotkey", "OnMuteHotkey", "OnCancelHotkey", "OnPanelHotkey",
	"OnTrayPanel", "OnTrayMute", "OnTrayPauseWake", "OnTrayExit", "OnDragEnd",
}

// ballHost228 is one production function that creates the ball.
type ballHost228 struct {
	file string
	name string
}

// productionFiles228 parses every non-test .go file of this package. Test files
// are excluded by name because an importer in a _test.go file would satisfy a
// grep and satisfy nothing else - that is exactly the reading AC#1 refuses.
func productionFiles228(t *testing.T) map[string]*ast.File {
	t.Helper()
	names, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("reading this package's directory: %v", err)
	}
	fset := token.NewFileSet()
	out := map[string]*ast.File{}
	for _, e := range names {
		n := e.Name()
		if e.IsDir() || !strings.HasSuffix(n, ".go") || strings.HasSuffix(n, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(".", n), nil, parser.ParseComments)
		if err != nil {
			t.Fatalf("parsing %s (the instrument, not the code): %v", n, err)
		}
		out[n] = f
	}
	if len(out) == 0 {
		t.Fatalf("this package holds no production .go file, so every claim below is an empty instrument")
	}
	return out
}

// hostsOfBall228 returns the production functions that call ball.New, plus the
// production files that import internal/ball at all.
func hostsOfBall228(files map[string]*ast.File) ([]ballHost228, []string) {
	var hosts []ballHost228
	var importers []string
	for name, f := range files {
		imports := false
		for _, imp := range f.Imports {
			if strings.Trim(imp.Path.Value, `"`) == ballImportPath228 {
				imports = true
			}
		}
		if !imports {
			continue
		}
		importers = append(importers, name)
		for _, d := range f.Decls {
			fn, ok := d.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || fn.Body == nil {
				continue
			}
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "ball" && sel.Sel.Name == "New" {
					hosts = append(hosts, ballHost228{file: name, name: fn.Name.Name})
				}
				return true
			})
		}
	}
	sort.Strings(importers)
	sort.Slice(hosts, func(i, j int) bool { return hosts[i].name < hosts[j].name })
	return hosts, importers
}

// funcBody228 finds one package-level function body by name.
func funcBody228(files map[string]*ast.File, want string) (*ast.BlockStmt, string) {
	for name, f := range files {
		for _, d := range f.Decls {
			fn, ok := d.(*ast.FuncDecl)
			if ok && fn.Recv == nil && fn.Name.Name == want {
				return fn.Body, name
			}
		}
	}
	return nil, ""
}

// isWindowsBuild228 reports whether Go would build this file on Windows: the
// implicit filename rule, or an explicit build comment that mentions windows
// without negating it. This is how the walk tells the leg under test apart from
// resident_other.go, whose whole contract is refusing to start.
func isWindowsBuild228(name string, f *ast.File) bool {
	if strings.HasSuffix(name, "_windows.go") {
		return true
	}
	for _, g := range f.Comments {
		for _, c := range g.List {
			if !strings.HasPrefix(c.Text, "//go:build") {
				continue
			}
			if strings.Contains(c.Text, "windows") && !strings.Contains(c.Text, "!windows") {
				return true
			}
		}
	}
	return false
}

// windowsResidentLeg228 is clause (3)'s subject: the body of the Windows
// runResident.
func windowsResidentLeg228(files map[string]*ast.File) (*ast.BlockStmt, string) {
	for name, f := range files {
		if !isWindowsBuild228(name, f) {
			continue
		}
		for _, d := range f.Decls {
			if fn, ok := d.(*ast.FuncDecl); ok && fn.Recv == nil && fn.Name.Name == "runResident" {
				return fn.Body, name
			}
		}
	}
	return nil, ""
}

// filesWithRunResident228 names what the instrument did see, so a rename of the
// resident leg reads as a finding rather than as a silent zero.
func filesWithRunResident228(files map[string]*ast.File) []string {
	var out []string
	for name, f := range files {
		for _, d := range f.Decls {
			if fn, ok := d.(*ast.FuncDecl); ok && fn.Recv == nil && fn.Name.Name == "runResident" {
				out = append(out, fmt.Sprintf("%s(build-windows=%v)", name, isWindowsBuild228(name, f)))
			}
		}
	}
	sort.Strings(out)
	return out
}

// TestAC228ResidentLegIsTheBallHost is AC#1's shape claim, in four clauses that
// each fail for a different deletion.
func TestAC228ResidentLegIsTheBallHost(t *testing.T) {
	files := productionFiles228(t)
	hosts, importers := hostsOfBall228(files)

	// (1) The AC#1 ruler, as an instrument: at least one production file of the
	// leg that runs tasks imports the ball package. Before ticket 228 this was
	// zero for cmd/wisp and one for cmd/balldebug.
	if len(importers) == 0 {
		t.Fatalf("AC#1 RED: no production file of cmd/wisp imports %s, so the ball is compiled into nothing "+
			"anybody launches. That is the exact reading ticket 228 opened with.", ballImportPath228)
	}
	// (2) Somebody in this package actually creates it, rather than importing
	// the name for a type it never boots.
	if len(hosts) == 0 {
		t.Fatalf("AC#1 RED: %v import internal/ball but no function in this package calls ball.New.", importers)
	}
	t.Logf("ball host functions (production): %s", hostNames228(hosts))

	// (3) The resident leg - func main's no-args branch, the one D2 draws as the
	// resident main process - walks to that host. Delete the call in
	// resident_windows.go and this clause is the red name.
	//
	// Two files declare runResident, and only one of them is the leg under
	// test: resident_other.go is //go:build !windows and its whole contract is
	// refusing to start (exit 2, ticket 78), so holding it to "hosts a ball"
	// would be a claim about a platform that has no layered window to host.
	// The selector below is Go's own: a file named *_windows.go, or one whose
	// build comment mentions windows, is the Windows leg.
	body, file := windowsResidentLeg228(files)
	if body == nil {
		t.Fatalf("AC#1 RED (the instrument, not the code): no Windows build of runResident found among this "+
			"package's production files, so the walk below has no subject (saw %v).", filesWithRunResident228(files))
	}
	called := map[string]bool{}
	ast.Inspect(body, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if id, ok := call.Fun.(*ast.Ident); ok {
				called[id.Name] = true
			}
		}
		return true
	})
	hits := 0
	for _, h := range hosts {
		if called[h.name] {
			hits++
		}
	}
	if hits == 0 {
		t.Fatalf("AC#1 RED: runResident (%s) calls none of the ball hosts %v, so the leg that stays alive still "+
			"hosts nothing and the ball is back to living only in a debug harness.", file, hostNames228(hosts))
	}

	// (4) The same walk has to reach a teardown: the gestures stop mattering at
	// exit, and D38(e) step 2 is "hotkey + wake-word listening stops". A host
	// with no deferred stop is a leaked ui-sta thread and a tray icon nobody
	// removed.
	deferred := false
	for _, st := range body.List {
		d, ok := st.(*ast.DeferStmt)
		if !ok {
			continue
		}
		if sel, ok := d.Call.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "stop" {
			deferred = true
		}
	}
	if !deferred {
		t.Fatalf("AC#1 RED: runResident (%s) defers no ball teardown, so the hot keys and the tray icon it "+
			"installed are left to process death. See internal/proc/shutdown.go:18 (D38(e) step 2).", file)
	}
}

// TestAC228BallHostAnswersEveryGesture is the anti-silence clause: every
// consumer callback internal/ball offers is given a non-nil executor by the
// resident host, because a nil entry is a click or a hot key that disappears
// without a log record.
func TestAC228BallHostAnswersEveryGesture(t *testing.T) {
	files := productionFiles228(t)
	hosts, _ := hostsOfBall228(files)
	if len(hosts) == 0 {
		t.Fatal("AC#1 RED: no ball host in this package to inspect (see TestAC228ResidentLegIsTheBallHost).")
	}

	set := map[string]bool{}
	for _, h := range hosts {
		body, _ := funcBody228(map[string]*ast.File{h.file: fileOf228(t, files, h.file)}, h.name)
		if body == nil {
			continue
		}
		ast.Inspect(body, func(n ast.Node) bool {
			lit, ok := n.(*ast.CompositeLit)
			if !ok {
				return true
			}
			if sel, ok := lit.Type.(*ast.SelectorExpr); !ok || sel.Sel.Name != "Events" {
				return true
			}
			for _, elt := range lit.Elts {
				if kv, ok := elt.(*ast.KeyValueExpr); ok {
					if id, ok := kv.Key.(*ast.Ident); ok {
						set[id.Name] = true
					}
				}
			}
			return true
		})
	}

	var missing []string
	for _, want := range ballEventCallbacks228 {
		if !set[want] {
			missing = append(missing, want)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		t.Fatalf("AC#1 RED: the resident ball host leaves %d of internal/ball's %d gesture callbacks unset: %v. "+
			"A nil callback is a click or hot key that reaches no one and logs nothing.",
			len(missing), len(ballEventCallbacks228), missing)
	}
	if len(set) > len(ballEventCallbacks228) {
		t.Fatalf("AC#1 RED (the instrument, not the code): the host sets %d callbacks but this file lists %d; "+
			"internal/ball.Events grew and ballEventCallbacks228 did not.", len(set), len(ballEventCallbacks228))
	}
	t.Logf("all %d gesture callbacks answered by the resident host", len(ballEventCallbacks228))
}

// fileOf228 keeps the walk above honest: a host name is looked up in the file
// that declared it, never in whichever file the map happened to yield first.
func fileOf228(t *testing.T, files map[string]*ast.File, name string) *ast.File {
	t.Helper()
	f, ok := files[name]
	if !ok {
		t.Fatalf("the walk lost %s between parsing and inspecting (the instrument, not the code)", name)
	}
	return f
}

func hostNames228(hosts []ballHost228) string {
	parts := make([]string, 0, len(hosts))
	for _, h := range hosts {
		parts = append(parts, h.file+":"+h.name)
	}
	return strings.Join(parts, ", ")
}
