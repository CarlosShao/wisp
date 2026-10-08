//go:build windows

package main

// Ticket 255 AC#4 - the number in [panel] reaches the window.
//
// WHAT THIS FILE HAS TO PROVE, and why the shape is what it is. Before this
// ticket the host asked WebView2 for a literal 420x260
// (cmd/wisp/panel_host_windows.go:304-305 at HEAD 67ab595d), so [panel]
// width/height were parsed, validated, hot-applied into memory - and read by
// nobody. The break was "the host never receives config", not "the value is not
// stored", so `go build` and `go vet` answering rc=0 proves nothing here: the
// wiring claim is about a VALUE ARRIVING AT A CALL, and it is measured as a
// number plus a call-site count.
//
// Four rulers, in the order they are read:
//
//  1. TestTicket255WindowOptionsFollowTheConfigSource - the resolved geometry,
//     as data: what the host computes and hands the create block, for a source
//     answering each of the shapes [panel] can actually produce (including
//     height 0, whose "auto from content" is registered as NOT done, and the
//     no-source case that keeps the seven pre-AC#4 test call sites at 420x260).
//  2. TestTicket255AssemblyRootGeometrySourceReachesTheWindowOptions - the
//     production closure over a real config.toml on disk, asserted through
//     windowOptions() rather than through the field: the width that arrives is
//     the width that was written, twice, with a rewrite in between. That second
//     read is the 取值闭包 half of the ruling - a snapshot shape would answer
//     the first number both times and this case would go red.
//  3. TestTicket255PanelHostBuildsItsWindowOptions - the consumption half, and
//     the reason ruler 1 is not enough: this package's production sources carry
//     exactly ONE webview2.NewWithOptions call site, it lives in the host file, it
//     takes its WindowOptions from a call to windowOptions, and no int literal may
//     sit in a WindowOptions block outside that function. Ruler 1 alone is blind.
//  4. TestTicket255HostStillDoesNotParseConfigItself +
//     TestTicket255PanelRosterVerdictIsTheHonestShape - the two things AC#4 was
//     NOT allowed to break: no config import and no disk read inside the host
//     file (the dependency edge stays where the ruling put it), and the receipt
//     wording stays "close it and reopen it", never "drag it and it resizes".
//
// NO REAL WINDOW IS CREATED HERE, and that is deliberate. AC#4's other half -
// the on-screen width following along - is a winlive reading and is recorded in
// .scratch/wisp/probes/255/r3/impl.md, not pretended to by this file: nothing in
// here calls bringUp, Show or Destroy, so nothing in here owes a thread back to
// the pool (票 255's own血账 about a real window poisoning the next case on the
// same M does not apply to a case that opens no window).
//
// PATH warning (ticket 98): this package's test binary links sherpa-onnx and dies
// at load (0xc0000135) unless third_party/sherpa-onnx is on PATH.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/panel"
)

// --- ruler 1: the resolved geometry, as data -----------------------------------

func TestTicket255WindowOptionsFollowTheConfigSource(t *testing.T) {
	cases := []struct {
		name string
		// src is what the assembly root answers with. nil means "no source was
		// handed at all", the shape all seven pre-AC#4 call sites are.
		src       func() (int, int)
		wantWidth uint
		wantHigh  uint
	}{
		{
			name:      "no source at all keeps the host's own constants",
			src:       nil,
			wantWidth: 420,
			wantHigh:  260,
		},
		{
			name:      "a config width and height both arrive",
			src:       func() (int, int) { return 517, 331 },
			wantWidth: 517,
			wantHigh:  331,
		},
		{
			name: "height 0 is NOT auto-height: the ruling is 260, and schema.go's " +
				"PanelSection.Height carries no default tag so 0 is what a " +
				"config saying nothing answers",
			src:       func() (int, int) { return 640, 0 },
			wantWidth: 640,
			wantHigh:  260,
		},
		{
			name:      "width 0 falls back to the host constant, not to a guess",
			src:       func() (int, int) { return 0, 400 },
			wantWidth: 420,
			wantHigh:  400,
		},
		{
			name:      "a negative answer is garbage, and garbage is not a second spelling of auto",
			src:       func() (int, int) { return -1, -1 },
			wantWidth: 420,
			wantHigh:  260,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewPanelManager(panelDispatchForGeometry255(), nil, "", withGeometrySourceOrNil(tc.src))
			got := m.windowOptions()
			if got.Width != tc.wantWidth {
				t.Errorf("windowOptions().Width = %d, want %d: the number the window is created with is not the number this host resolved", got.Width, tc.wantWidth)
			}
			if got.Height != tc.wantHigh {
				t.Errorf("windowOptions().Height = %d, want %d", got.Height, tc.wantHigh)
			}
			if got.Title != panelTitle {
				t.Errorf("windowOptions().Title = %q, want the host's own panelTitle %q - AC#4 rewrote this block and must not have moved the title with it", got.Title, panelTitle)
			}
		})
	}
}

// withGeometrySourceOrNil is the test-side spelling of "the assembly root handed
// nothing" (nil src) versus "it handed a source". withGeometrySource(nil) would
// set the field to nil and mean the same thing, so the guard exists to keep the
// case table readable, not to add behaviour.
func withGeometrySourceOrNil(src func() (int, int)) panelHostOption {
	return func(m *PanelManager) { m.geometry = src }
}

// --- ruler 2: disk to create block, and re-read on every create ----------------

func TestTicket255AssemblyRootGeometrySourceReachesTheWindowOptions(t *testing.T) {
	dir := t.TempDir()
	writePanelConfig255(t, dir, 517, 331)

	calls := 0
	m := NewPanelManager(panelDispatchForGeometry255(), nil, filepath.Join(dir, "webview2"),
		withGeometrySource(func() (int, int) {
			calls++
			return panelGeometrySource(dir)()
		}))

	first := m.windowOptions()
	if first.Width != 517 || first.Height != 331 {
		t.Fatalf("the create block got %dx%d from a config.toml carrying width=517 height=331 - AC#4's whole claim is that these two are the same number",
			first.Width, first.Height)
	}
	if calls != 1 {
		t.Errorf("the geometry source was called %d times for one windowOptions, want exactly 1: two reads of one create means the window can be sized from two different file versions", calls)
	}

	// The 取值闭包 half. A snapshot-at-assembly shape answers 517 again here and
	// this is the line that reddens.
	writePanelConfig255(t, dir, 733, 0)
	second := m.windowOptions()
	if second.Width != 733 {
		t.Errorf("after [panel] width was rewritten to 733 the host still hands the create %d: the value was snapshotted, not re-read per create, which contradicts internal/config/tiers.go:34 registering [panel] as hot", second.Width)
	}
	if second.Height != 260 {
		t.Errorf("height 0 must resolve to the host's constant 260, got %d", second.Height)
	}
	if calls != 2 {
		t.Errorf("the source was called %d times across two windowOptions, want 2 (one per create)", calls)
	}

	// And the unreadable-config branch: no number invented, the host's own
	// constants answer, and the host says so (the slog.Warn is exercised by
	// reaching it, not asserted - what is pinned is the VALUE).
	m2 := NewPanelManager(panelDispatchForGeometry255(), nil, "",
		withGeometrySource(panelGeometrySource(filepath.Join(dir, "no-such-dir"))))
	got := m2.windowOptions()
	if got.Width != 420 || got.Height != 260 {
		t.Errorf("an unreadable config produced %dx%d, want the host constants 420x260", got.Width, got.Height)
	}
}

// writePanelConfig255 puts a loadable config.toml with one [panel] row in dir.
// schema_version is current so readConfigFile takes no migration branch (that
// branch rewrites the file, and a rewrite from under a re-read test would be a
// second mover in a measurement about who read what when).
func writePanelConfig255(t *testing.T, dir string, width, height int) {
	t.Helper()
	body := "schema_version = 2\n\n[panel]\n"
	if width > 0 {
		body += "width = " + strconv.Itoa(width) + "\n"
	}
	if height > 0 {
		body += "height = " + strconv.Itoa(height) + "\n"
	}
	path := filepath.Join(dir, configFileName)
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// --- ruler 3: the create site actually consumes it -----------------------------

// TestTicket255PanelHostBuildsItsWindowOptions is the wiring ruler: it reads the
// host's source and demands that the geometry reaching webview2.NewWithOptions is
// a CALL to windowOptions, that this package's production sources hold exactly one
// NewWithOptions site and that it is this host file's, and that no int literal sits
// in a WindowOptions block outside windowOptions' own body. Ruler 1 measures a
// function; this one measures the create built from it - 255-r6 widened the count.
func TestTicket255PanelHostBuildsItsWindowOptions(t *testing.T) {
	src := readHostFileForGeometry255(t, filepath.Join(hostPkgDir255(t), "panel_host_windows.go"))
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "panel_host_windows.go", src, parser.AllErrors)
	if err != nil {
		t.Fatalf("parse cmd/wisp/panel_host_windows.go: %v", err)
	}

	var newWithOptions []*ast.CallExpr
	var windowOptionsBodies []*ast.FuncDecl
	windowsOptionsLits := map[*ast.CompositeLit]*ast.FuncDecl{} // literal -> enclosing func
	var inFunc *ast.FuncDecl
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.FuncDecl:
			inFunc = node
			if node.Name.Name == "windowOptions" {
				windowOptionsBodies = append(windowOptionsBodies, node)
			}
		case *ast.CallExpr:
			if sel, ok := node.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "NewWithOptions" {
				newWithOptions = append(newWithOptions, node)
			}
		case *ast.CompositeLit:
			if isWindowOptionsLit255(node.Type) {
				windowsOptionsLits[node] = inFunc
			}
		}
		return true
	})
	sites255 := newWithOptionsSitesInPackage255r6(t, fset, hostPkgDir255(t))
	if len(newWithOptions) != 1 || len(sites255) != 1 {
		t.Fatalf("cmd/wisp production sources carry %d webview2.NewWithOptions call sites [%s] and cmd/wisp/panel_host_windows.go itself carries %d, want 1 and 1: AC#4's claim is about the geometry THIS create takes, and a second create anywhere in this package's shipped code - not just in this file - would be a second answer to check", len(sites255), strings.Join(sites255, " "), len(newWithOptions))
	}
	if len(windowOptionsBodies) != 1 {
		t.Fatalf("found %d declarations of windowOptions, want 1", len(windowOptionsBodies))
	}

	// The create's argument block must name windowOptions for its geometry.
	foundGeometryCall := false
	ast.Inspect(newWithOptions[0], func(n ast.Node) bool {
		kv, ok := n.(*ast.KeyValueExpr)
		if !ok || !isIdent255(kv.Key, "WindowOptions") {
			return true
		}
		call, ok := kv.Value.(*ast.CallExpr)
		if !ok {
			t.Errorf("the create's WindowOptions field is a %T, not a call: the geometry is not being resolved at create time", kv.Value)
			return false
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "windowOptions" {
			t.Errorf("the create's WindowOptions comes from %s, want a call to windowOptions - AC#4's whole shape is that ONE function answers the window size", exprText255(fset, call.Fun, src))
			return false
		}
		foundGeometryCall = true
		return false
	})
	if !foundGeometryCall {
		t.Errorf("the one NewWithOptions call has no WindowOptions field fed by windowOptions(); the number measured in ruler 1 never reaches the create")
	}

	// No int literal may size a window anywhere except inside windowOptions,
	// where the fallback constants are named, not spelled.
	for lit, fn := range windowsOptionsLits {
		if fn != nil && fn.Name.Name == "windowOptions" {
			for _, e := range lit.Elts {
				kv, ok := e.(*ast.KeyValueExpr)
				if !ok || !(isIdent255(kv.Key, "Width") || isIdent255(kv.Key, "Height")) {
					continue
				}
				if bl, ok := kv.Value.(*ast.BasicLit); ok {
					t.Errorf("windowOptions sizes %s from the literal %s: the fallback must go through panelWidthPx/panelHeightPx so this file has one place the default lives",
						exprText255(fset, kv.Key, src), bl.Value)
				}
			}
			continue
		}
		who := "<package level>"
		if fn != nil {
			who = fn.Name.Name
		}
		t.Errorf("a webview2.WindowOptions composite literal in %s is a second place the panel geometry could come from; only windowOptions() may build it", who)
	}

	// bringUp is the only caller, and it is the create function AC#4's ruling names.
	bringUpText := funcBodyText255(t, fset, file, "bringUp", src)
	if !strings.Contains(bringUpText, "m.windowOptions()") {
		t.Errorf("bringUp does not call m.windowOptions(); its create block reads:\n%s", bringUpText)
	}
}

// --- ruler 4: what AC#4 was not allowed to break ------------------------------

// TestTicket255HostStillDoesNotParseConfigItself pins the two forbidden shapes:
// a panel->config dependency edge, and the host reading config.toml on its own.
// The host file is the one that must stay ignorant - the parsing lives in
// cmd/wisp/panel_resident_windows.go's assembly root, which cmd/wisp already
// imported before this ticket (cmd/wisp/run.go:409), so this package's import
// graph gains nothing.
func TestTicket255HostStillDoesNotParseConfigItself(t *testing.T) {
	// WHERE this package's own files live. go test's cwd is the package dir and
	// panel_host_gate_test.go's production scan relies on that, but this ruler is
	// reading three other files by name, so it pins the dir off runtime.Caller the
	// way config_receipt_255_test.go's repoRootForRoster does. The build tags that
	// make those three exist are windows-only, so reading them from anywhere else
	// would be a measurement of a file that is not in the build.
	dir := hostPkgDir255(t)
	hostPath := filepath.Join(dir, "panel_host_windows.go")
	hostSrc := readHostFileForGeometry255(t, hostPath)
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, hostPath, hostSrc, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("parse the import block of %s: %v", hostPath, err)
	}
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		if strings.Contains(path, "internal/config") {
			t.Errorf("cmd/wisp/panel_host_windows.go now imports %q: 票 255 AC#4 forbids a panel->config edge, and the geometry is handed as a function value for exactly that reason", path)
		}
	}
	// Text-level, because an import is not the only way to go to disk: the host
	// must also not be reading a config file by name.
	for _, banned := range []string{"config.LoadFile", "config.NewManager", "configFileName", "config.toml"} {
		for i, line := range strings.Split(hostSrc, "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "//") {
				continue
			}
			if strings.Contains(trimmed, banned) {
				t.Errorf("panel_host_windows.go:%d names %q outside a comment - the host does not read config.toml itself, the assembly root hands the value", i+1, banned)
			}
		}
	}
	// And the same ruler over the resident assembly root, whose file is the one
	// holding the new read site: it must read through config, never by hand.
	root := readHostFileForGeometry255(t, filepath.Join(hostPkgDir255(t), "panel_resident_windows.go"))
	if !strings.Contains(root, "config.LoadFile") {
		t.Errorf("cmd/wisp/panel_resident_windows.go no longer reads [panel] through config.LoadFile - the read site the roster cites has moved")
	}
	if !strings.Contains(root, "withGeometrySource") {
		t.Errorf("cmd/wisp/panel_resident_windows.go no longer hands the host a geometry source: AC#4's wiring was removed and the host is back to its own constants")
	}
}

// TestTicket255PanelRosterVerdictIsTheHonestShape pins the wording AC#4 owes the
// operator once the value really arrives. Two halves, both from 票 255's ruling:
// the sentence is "close the panel and open it again", never "drag it and it
// resizes" (the resize 票 255-r1 added runs on a show request, never at the instant
// config.toml is saved), and [panel]'s roster row must name the live read site.
func TestTicket255PanelRosterVerdictIsTheHonestShape(t *testing.T) {
	verdict, ok := hotRowClaims["panel"]
	if !ok {
		t.Fatal("hotRowClaims lost its panel row: every TierRegistry hot row owes a verdict (票 255 AC#1)")
	}
	for _, lie := range []string{"拖动即变", "立即改变窗口大小", "自动改大小", "auto from content 已实现"} {
		if strings.Contains(verdict, lie) {
			t.Errorf("the [panel] verdict claims %q; the only resize this host has fires on the NEXT show request, never at the instant config.toml is saved, so the honest shape stays 关窗再开: %q", lie, verdict)
		}
	}
	if !strings.Contains(verdict, "关窗再开") {
		t.Errorf("the [panel] verdict must state the effect the operator can actually get, 关窗再开即跟上新值: %q", verdict)
	}
	// The retired cite must not survive as if it still described the code. AC#4
	// removed the hard-coded geometry literals from bringUp's create block, and a
	// roster pointing at that retired shape is how a fixed bug keeps being
	// reported as present. TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim
	// would catch token drift, but only for a cite that survives the line shuffle,
	// and this one is about the whole claim - both needles below are CONTENT.
	src := readHostFileForGeometry255(t, filepath.Join(hostPkgDir255(t), "config_readers_255.go"))
	for _, retired := range []string{"panel_host_windows.go:304 [", "Width:  420,"} {
		for i, line := range strings.Split(src, "\n") {
			if strings.Contains(line, retired) {
				t.Errorf("config_readers_255.go:%d still carries the retired geometry evidence %q as if it were live: %s",
					i+1, retired, strings.TrimSpace(line))
			}
		}
	}
}

// --- helpers ------------------------------------------------------------------

// panelDispatchForGeometry255 is the empty inbound router. AC#4's measurements are
// about the bytes handed to a create, so no window is ever brought up here and the
// dispatch never runs; a zero value is what the host's own constructor accepts.
func panelDispatchForGeometry255() *panel.ComposerDispatch {
	return &panel.ComposerDispatch{}
}

func isWindowOptionsLit255(e ast.Expr) bool {
	sel, ok := e.(*ast.SelectorExpr)
	return ok && sel.Sel.Name == "WindowOptions"
}

func isIdent255(e ast.Expr, name string) bool {
	id, ok := e.(*ast.Ident)
	return ok && id.Name == name
}

// exprText255 renders an expression back from the source it was parsed out of.
// The offsets are only meaningful for the same bytes that were handed to the
// parser, which is why src travels with it.
func exprText255(fset *token.FileSet, e ast.Expr, src string) string {
	start, end := fset.Position(e.Pos()).Offset, fset.Position(e.End()).Offset
	if start < 0 || end > len(src) || start >= end {
		return "<unprintable>"
	}
	return src[start:end]
}

func funcBodyText255(t *testing.T, fset *token.FileSet, file *ast.File, name, src string) string {
	t.Helper()
	for _, d := range file.Decls {
		fn, ok := d.(*ast.FuncDecl)
		if !ok || fn.Name.Name != name {
			continue
		}
		start := fset.Position(fn.Pos()).Offset
		end := fset.Position(fn.End()).Offset
		if start >= 0 && end <= len(src) && end > start {
			return src[start:end]
		}
	}
	t.Fatalf("no readable body for function %s in cmd/wisp/panel_host_windows.go", name)
	return ""
}

// hostPkgDir255 is this package's directory, derived the same way
// config_receipt_255_test.go's repoRootForRoster derives the repo root. The three
// files read through it are windows-tagged, so a cwd-relative read could measure a
// file that is not in this build; a failed read is reported as a failed
// measurement, never as a zero.
func hostPkgDir255(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// thisFile = <repo>/cmd/wisp/panel_geometry_255_test.go
	return filepath.Dir(thisFile)
}

func readHostFileForGeometry255(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
