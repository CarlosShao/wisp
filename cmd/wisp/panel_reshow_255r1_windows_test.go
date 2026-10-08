//go:build windows

package main

// Ticket 255 AC#4, leg 255-r1 - 甲形 on the LIVE window, measured without a window.
//
// WHAT THIS FILE MEASURES. Show's re-show branch now hands an already-created panel
// window the geometry this host resolves at that instant, through the library's own
// webview2.WebView.SetSize posted on webview2.WebView.Dispatch
// (cmd/wisp/panel_host_windows.go's requestGeometryOnReshow). The subject is a
// CAPABILITY - "which method got called, with which arguments, through which hop" -
// so every assertion here is a recorded call on a fake control or a node in the
// parsed host file, never a string searched for in source text. Two rulers in this
// package were retired for being behaviour-blind; a ruler that grepped the word
// SetSize would pass on a host that never calls it.
//
// FOUR SEPARATE QUESTIONS, each with its own red shape:
//
//	(a) did a SetSize go out at all, and with WHICH pair? The pair is asserted
//	    against the geometry source's answer through the host's own resolve rule, so
//	    a host that resized to its constants while the config said 517 is red, and a
//	    host that skipped the "height 0 is not auto-height" rule is red too.
//	(b) did it go out THROUGH Dispatch? The fake records a SetSize that arrives while
//	    no Dispatch closure is running in a DIFFERENT list, which the tests demand
//	    stay empty. Replacing w.Dispatch(func(){w.SetSize(...)}) with a bare
//	    w.SetSize(...) - the shape 255-a2 refused to endorse because the repository
//	    holds no measurement of a cross-thread SetWindowPos - reddens (b) while
//	    leaving (a) green. That split is the point: they are different claims.
//	(c) WHICH hint? HintNone, asserted as a value at every recorded call AND read off
//	    the call site in the parsed file, because HintFixed would clear WSThickFrame
//	    and WSMaximizeBox and take away the user's ability to drag the window edge -
//	    an interface behaviour change 票 255 does not own.
//	(d) does the SHIPPING entry point (Show) actually reach it? One case drives Show
//	    itself against an attached control, so (a)-(c) are not measurements of a
//	    method nobody calls. The cold branch cannot be driven headlessly at all
//	    (bringUp needs a live WebView2 runtime), so that half is pinned structurally.
//
// WHAT THIS FILE DELIBERATELY DOES NOT CLAIM. It opens no window, so it says nothing
// about the width a person sees. The create's WindowOptions pair is read by
// CreateWindowExW as the OUTER FRAME while SetSize's pair is read as the CLIENT area
// (go-webview2 webview.go:296-320 versus :405 running AdjustWindowRect first), so
// the SAME configured number produces DIFFERENT on-screen widths depending on which
// hop carried it, by a frame delta that is 〔仅本机可量〕 - only a winlive reading can
// name that many pixels, and no case here pretends to have one. Nor does any case
// here claim that saving config.toml moves the panel: the tick that would do that
// does not exist in the resident process (internal/config/tiers.go registers "panel"
// as "hot", so manager.go's OnReload never fires for it, and startConfigReload's only
// non-test caller is cmd/wisp/run.go).
//
// PATH warning (ticket 98): this package's test binary links sherpa-onnx and dies at
// load (0xc0000135) unless third_party/sherpa-onnx is on PATH. That is a build/load
// failure, not a verdict on these judgements - triage on the build rc first.

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"unsafe"

	"github.com/CarlosShao/wisp/internal/panel"
	webview2 "github.com/jchv/go-webview2"
)

// compile-time proof that the fake is a stand-in for the real interface, so a
// library change that adds a method shows up here as a build failure rather than as
// a ruler quietly measuring nothing.
var _ webview2.WebView = (*geoSink255r1)(nil)

// sizeCall255r1 is one recorded SetSize: the pair, the hint, and the hop.
type sizeCall255r1 struct {
	width, height  int
	hint           webview2.Hint
	insideDispatch bool
}

// geoSink255r1 is a webview2.WebView that records geometry requests and can tell
// whether it was inside a Dispatch closure when one arrived. Dispatch runs the
// closure the way the library's Run() drains it, which is what makes "did the resize
// come out through Dispatch" an observable property instead of a text pattern.
type geoSink255r1 struct {
	mu   sync.Mutex
	size []sizeCall255r1 // SetSize reached while inside a Dispatch closure
	bare []sizeCall255r1 // SetSize reached NOT inside a Dispatch closure
	// depth counts the Dispatch closures this sink is running right now.
	depth      int
	dispatches int
	destroyed  int
	docs       []string
	bindings   map[string]interface{}
}

func newGeoSink255r1() *geoSink255r1 {
	return &geoSink255r1{bindings: map[string]interface{}{}}
}

func (c *geoSink255r1) Run()       {}
func (c *geoSink255r1) Terminate() {}

func (c *geoSink255r1) Dispatch(f func()) {
	if f == nil {
		return
	}
	c.mu.Lock()
	c.dispatches++
	c.depth++
	c.mu.Unlock()
	f()
	c.mu.Lock()
	c.depth--
	c.mu.Unlock()
}

func (c *geoSink255r1) Destroy() {
	c.mu.Lock()
	c.destroyed++
	c.mu.Unlock()
}

func (c *geoSink255r1) Window() unsafe.Pointer { return nil }
func (c *geoSink255r1) SetTitle(string)        {}
func (c *geoSink255r1) Init(string)            {}
func (c *geoSink255r1) Eval(string)            {}
func (c *geoSink255r1) Navigate(string)        {}

func (c *geoSink255r1) SetHtml(html string) {
	c.mu.Lock()
	c.docs = append(c.docs, html)
	c.mu.Unlock()
}

func (c *geoSink255r1) Bind(name string, f interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.bindings == nil {
		c.bindings = map[string]interface{}{}
	}
	c.bindings[name] = f
	return nil
}

func (c *geoSink255r1) SetSize(w, h int, hint webview2.Hint) {
	rec := sizeCall255r1{width: w, height: h, hint: hint}
	c.mu.Lock()
	if c.depth > 0 {
		rec.insideDispatch = true
		c.size = append(c.size, rec)
	} else {
		c.bare = append(c.bare, rec)
	}
	c.mu.Unlock()
}

func (c *geoSink255r1) snapshot() (sized, bare []sizeCall255r1, dispatches int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]sizeCall255r1(nil), c.size...),
		append([]sizeCall255r1(nil), c.bare...),
		c.dispatches
}

// attachGeoSink255r1 points a manager at a recorded control the way bringUp does
// after a successful create (m.w = w, created = true), and creates nothing. Same
// injection the 33-r10 ruler uses for the page-handover question.
func attachGeoSink255r1(m *PanelManager, c *geoSink255r1) {
	m.mu.Lock()
	m.w = c
	m.created = true
	m.mu.Unlock()
}

// managerOnRecordedWindow builds a host that already has a window, sized by src.
// src == nil is "nobody handed this host a geometry source", the pre-AC#4 shape.
func managerOnRecordedWindow(t *testing.T, src func() (int, int), calls *int) (*PanelManager, *geoSink255r1) {
	t.Helper()
	var wrapped func() (int, int)
	if src != nil {
		wrapped = func() (int, int) {
			*calls++
			return src()
		}
	}
	m := NewPanelManager(&panel.ComposerDispatch{}, nil, "", withGeometrySourceOrNil(wrapped))
	sink := newGeoSink255r1()
	attachGeoSink255r1(m, sink)
	return m, sink
}

// (a) + (b) + (c), asked of the method that owns the decision.
func TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch(t *testing.T) {
	cases := []struct {
		name       string
		width      int // the geometry source's answer; unused when noSource is set
		height     int
		noSource   bool
		wantWidth  int
		wantHeight int
	}{
		{
			// The pair below goes out in CLIENT-area terms; that is NOT the width a
			// create would have produced from the same number. See the header.
			name: "the configured pair goes out as it stands", width: 517, height: 331,
			wantWidth: 517, wantHeight: 331,
		},
		{
			// internal/config/schema.go's PanelSection.Height carries no default, so
			// 0 is what a config saying nothing answers, and 票 255's ruling lands
			// that as 260 with "由内容定高" registered as NOT done.
			name: "height 0 keeps the host's 260, the same rule the create uses", width: 640, height: 0,
			wantWidth: 640, wantHeight: 260,
		},
		{
			name: "width 0 falls back to the host constant, not to a guess", width: 0, height: 400,
			wantWidth: 420, wantHeight: 400,
		},
		{
			// panelGeometrySource answers 0,0 on its unreadable-config branch, so a
			// re-show moves the window to what a fresh create would have been sized
			// at. That consequence is spelled out in requestGeometryOnReshow rather
			// than left for the operator to infer.
			name: "an unreadable config answers the host constants", width: 0, height: 0,
			wantWidth: 420, wantHeight: 260,
		},
		{
			// A garbage answer is not a second spelling of auto.
			name: "a negative answer resolves the way the create resolves it", width: -1, height: -1,
			wantWidth: 420, wantHeight: 260,
		},
		{
			name:     "a host nobody sized sends nothing at all",
			noSource: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			var src func() (int, int)
			if !tc.noSource {
				w, h := tc.width, tc.height
				src = func() (int, int) { return w, h }
			}
			m, sink := managerOnRecordedWindow(t, src, &calls)

			m.requestGeometryOnReshow()

			sized, bare, dispatches := sink.snapshot()
			if tc.noSource {
				if len(sized) != 0 || len(bare) != 0 {
					t.Errorf("a host with no geometry source sent %d resize(s) through Dispatch and %d outside it, want 0 and 0 - AC#4's own invariant is that a host nobody sized behaves exactly as it did before this ticket: %+v / %+v",
						len(sized), len(bare), sized, bare)
				}
				if calls != 0 {
					t.Errorf("the absent source was consulted %d times, want 0", calls)
				}
				return
			}
			if len(bare) != 0 {
				t.Errorf("%d SetSize call(s) went out WITHOUT Dispatch: %+v - the repository holds no measurement of a cross-thread SetWindowPos on a WebView2 window (255-a2 says so by name), so the resize has to ride the library's own thread hop", len(bare), bare)
			}
			if len(sized) != 1 {
				t.Fatalf("the re-show sent %d SetSize request(s) through Dispatch, want exactly 1: %+v", len(sized), sized)
			}
			got := sized[0]
			if !got.insideDispatch {
				t.Errorf("the recorded SetSize is not marked as having come out of a Dispatch closure, so question (b) measured nothing: %+v", got)
			}
			if dispatches < 1 {
				t.Errorf("Dispatch ran %d time(s) yet a SetSize was recorded - those cannot both be true", dispatches)
			}
			if got.width != tc.wantWidth || got.height != tc.wantHeight {
				t.Errorf("the live window was asked for %dx%d, want %dx%d - what goes out must be the pair THIS host resolves from its source at this instant",
					got.width, got.height, tc.wantWidth, tc.wantHeight)
			}
			if got.hint != webview2.HintNone {
				t.Errorf("the resize went out with hint %d, want webview2.HintNone (%d): HintFixed would clear WSThickFrame/WSMaximizeBox and take away the user's ability to drag the window edge, which 票 255 does not own",
					int(got.hint), int(webview2.HintNone))
			}
			if calls != 1 {
				t.Errorf("the geometry source was consulted %d times for one re-show, want exactly 1: two reads could size the window from two different versions of the file", calls)
			}
		})
	}
}

// (d) the shipping entry point, not just the method: a Show on an already-created
// window must be what produces the request above. Without this case, (a)-(c) would
// measure a function nothing calls - the exact gap ruler 3 of
// panel_geometry_255_test.go exists to close on the create side.
func TestTicket255r1ShowOnAnExistingWindowSendsTheResize(t *testing.T) {
	calls := 0
	m, sink := managerOnRecordedWindow(t, func() (int, int) { return 517, 331 }, &calls)

	if err := m.Show(context.Background()); err != nil {
		t.Fatalf("Show on an already-created window returned %v - the re-show branch is the one AC#4 just grew", err)
	}

	sized, bare, _ := sink.snapshot()
	if len(bare) != 0 {
		t.Errorf("Show sent %d resize(s) outside Dispatch: %+v", len(bare), bare)
	}
	if len(sized) != 1 {
		t.Fatalf("Show on an existing window sent %d resize(s), want exactly 1: %+v", len(sized), sized)
	}
	if sized[0].width != 517 || sized[0].height != 331 || sized[0].hint != webview2.HintNone {
		t.Errorf("Show asked for %+v, want 517x331 under webview2.HintNone", sized[0])
	}
	if calls != 1 {
		t.Errorf("Show consulted the geometry source %d times, want 1 (the branch it took reads the source exactly once)", calls)
	}
}

// (b) self-check: the sink really separates the two shapes, so the "must go through
// Dispatch" ruler above is not constant-true. A sink that put every SetSize in one
// list would make that assertion vacuous.
func TestTicket255r1SinkTellsDispatchFromABareCall(t *testing.T) {
	sink := newGeoSink255r1()
	sink.SetSize(111, 222, webview2.HintNone) // the shape the product must NOT use
	sink.Dispatch(func() { sink.SetSize(333, 444, webview2.HintNone) })

	sized, bare, _ := sink.snapshot()
	if len(bare) != 1 || bare[0].width != 111 {
		t.Errorf("a SetSize called outside Dispatch landed in %+v, want exactly one 111-wide record in `bare` - the ruler above cannot tell the two shapes apart", bare)
	}
	if len(sized) != 1 || sized[0].width != 333 || !sized[0].insideDispatch {
		t.Errorf("a SetSize inside Dispatch landed as %+v, want one insideDispatch record of width 333", sized)
	}
}

// The guard against resizing a window that is gone: a Destroy between Show's
// created check and here must produce no call, no panic - and no disk read for a
// request that cannot go anywhere.
func TestTicket255r1ReshowWithNoLiveControlSendsNothing(t *testing.T) {
	calls := 0
	m, sink := managerOnRecordedWindow(t, func() (int, int) { return 517, 331 }, &calls)
	m.Destroy() // the recreate path's first half

	m.requestGeometryOnReshow()

	sized, bare, _ := sink.snapshot()
	if len(sized) != 0 || len(bare) != 0 {
		t.Errorf("a re-show with no live control sent %d resize(s) through Dispatch and %d outside it, want 0 and 0: %+v / %+v",
			len(sized), len(bare), sized, bare)
	}
	if calls != 0 {
		t.Errorf("the host read its geometry source %d time(s) for a window that no longer exists, want 0", calls)
	}
}

// Show must ask for the geometry ONLY on the branch where the window already
// existed. A cold Show sizes the window inside bringUp from the same source, and a
// second hop right after creation would move it by the frame/client delta with
// nothing in the ticket asking for that. This half is structural because the cold
// branch cannot be reached headlessly: bringUp calls webview2.NewWithOptions, which
// needs a live WebView2 runtime - the same reason panel_geometry_255_test.go parses
// this file instead of running it.
func TestTicket255r1ShowAsksForGeometryOnlyOnTheReshowBranch(t *testing.T) {
	file := parseHostFile255r1(t)
	show := funcDecl255r1(t, file, "Show")

	var createBranch, reshowBranch *ast.BlockStmt
	for _, stmt := range show.Body.List {
		ifStmt, ok := stmt.(*ast.IfStmt)
		if !ok || ifStmt.Else == nil {
			continue
		}
		if isNotCreatedTest255r1(ifStmt.Cond) {
			createBranch = ifStmt.Body
			if eb, ok := ifStmt.Else.(*ast.BlockStmt); ok {
				reshowBranch = eb
			}
		}
	}
	if createBranch == nil || reshowBranch == nil {
		t.Fatalf("Show no longer has an `if !created { ... } else { ... }` shape, so the re-show hop has no branch of its own and nothing here can say which Show asks for the geometry")
	}
	if n := countCall255r1(createBranch, "requestGeometryOnReshow"); n != 0 {
		t.Errorf("the `if !created` branch calls requestGeometryOnReshow %d time(s), want 0: a cold Show already sizes the window at create", n)
	}
	if n := countCall255r1(reshowBranch, "requestGeometryOnReshow"); n != 1 {
		t.Errorf("the already-created branch calls requestGeometryOnReshow %d time(s), want exactly 1 - that call IS AC#4's 甲形 on a live window", n)
	}
	if !callIsBareStatement255r1(reshowBranch, "requestGeometryOnReshow") {
		t.Errorf("the re-show branch does not run requestGeometryOnReshow as an unconditional statement, so whether a given show asks is hidden behind a condition this file cannot see")
	}
}

// (c) at the call site: one resize path in the host, and it names HintNone.
func TestTicket255r1HostResizeCallSiteAsksForHintNoneAndIsTheOnlyOne(t *testing.T) {
	file := parseHostFile255r1(t)

	var setSize []*ast.CallExpr
	ast.Inspect(file, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "SetSize" {
				setSize = append(setSize, call)
			}
		}
		return true
	})
	if len(setSize) != 1 {
		t.Fatalf("cmd/wisp/panel_host_windows.go has %d SetSize call site(s), want exactly 1: every extra one is a second place a hint could be picked wrong and a second geometry claim to check", len(setSize))
	}
	call := setSize[0]
	if len(call.Args) != 3 {
		t.Fatalf("the SetSize call site passes %d argument(s), want 3 (width, height, hint)", len(call.Args))
	}
	if !isPackageSelector255r1(call.Args[2], "webview2", "HintNone") {
		t.Errorf("the resize is posted with hint %s, want webview2.HintNone - HintFixed would clear WSThickFrame/WSMaximizeBox and switch off the user dragging the window edge, which is an interface behaviour change outside 票 255", selectorText255r1(call.Args[2]))
	}

	reshow := funcDecl255r1(t, file, "requestGeometryOnReshow")
	if !dispatchWrapsSetSize255r1(reshow) {
		t.Errorf("requestGeometryOnReshow does not post its SetSize from inside a Dispatch closure: 票 255-a2 recorded NO credential for a cross-thread SetWindowPos, and this file must not promote common sense to a measurement")
	}
}

// --- helpers ------------------------------------------------------------------

const hostFile255r1 = "panel_host_windows.go"

func parseHostFile255r1(t *testing.T) *ast.File {
	t.Helper()
	// hostPkgDir255 derives this package's directory from runtime.Caller, so the
	// read is anchored to the file's own location rather than to a cwd or a line
	// number that a neighbouring leg can move.
	path := filepath.Join(hostPkgDir255(t), hostFile255r1)
	raw := readHostFileForGeometry255(t, path)
	file, err := parser.ParseFile(token.NewFileSet(), hostFile255r1, raw, parser.AllErrors)
	if err != nil {
		t.Fatalf("parse cmd/wisp/%s: %v", hostFile255r1, err)
	}
	return file
}

func funcDecl255r1(t *testing.T, file *ast.File, name string) *ast.FuncDecl {
	t.Helper()
	for _, d := range file.Decls {
		if fn, ok := d.(*ast.FuncDecl); ok && fn.Name.Name == name && fn.Body != nil {
			return fn
		}
	}
	t.Fatalf("cmd/wisp/%s declares no readable body for %s - the re-show hop was renamed or removed", hostFile255r1, name)
	return nil
}

// isNotCreatedTest255r1 recognises the `if !created` condition by SHAPE (a unary NOT
// over the identifier named created), so rewriting the surrounding prose cannot
// quiet the ruler above while the branch structure stays.
func isNotCreatedTest255r1(cond ast.Expr) bool {
	unary, ok := cond.(*ast.UnaryExpr)
	if !ok || unary.Op.String() != "!" {
		return false
	}
	id, ok := unary.X.(*ast.Ident)
	return ok && id.Name == "created"
}

func countCall255r1(root ast.Node, methodName string) int {
	n := 0
	ast.Inspect(root, func(node ast.Node) bool {
		if call, ok := node.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == methodName {
				n++
			}
		}
		return true
	})
	return n
}

func callIsBareStatement255r1(block *ast.BlockStmt, methodName string) bool {
	for _, stmt := range block.List {
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}
		switch x := exprStmt.X.(type) {
		case *ast.SelectorExpr:
			if x.Sel != nil && x.Sel.Name == methodName {
				return true
			}
		case *ast.CallExpr:
			if sel, ok := x.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == methodName {
				return true
			}
		}
	}
	return false
}

// dispatchWrapsSetSize255r1 reports whether the single SetSize in fn's body sits
// inside a closure handed to a Dispatch call.
func dispatchWrapsSetSize255r1(fn *ast.FuncDecl) bool {
	total := countCall255r1(fn, "SetSize")
	if total != 1 {
		return false
	}
	wrapped := false
	ast.Inspect(fn, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "Dispatch" {
			return true
		}
		for _, arg := range call.Args {
			if closure, isFuncLit := arg.(*ast.FuncLit); isFuncLit && countCall255r1(closure, "SetSize") > 0 {
				wrapped = true
			}
		}
		return true
	})
	return wrapped
}

func isPackageSelector255r1(e ast.Expr, pkg, name string) bool {
	sel, ok := e.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != name {
		return false
	}
	id, ok := sel.X.(*ast.Ident)
	return ok && id.Name == pkg
}

func selectorText255r1(e ast.Expr) string {
	if sel, ok := e.(*ast.SelectorExpr); ok {
		if id, ok := sel.X.(*ast.Ident); ok {
			return id.Name + "." + sel.Sel.Name
		}
		return sel.Sel.Name
	}
	return strings.TrimPrefix(fmtTypeName255r1(e), "*")
}

func fmtTypeName255r1(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.BasicLit:
		return "literal " + t.Value
	case *ast.Ident:
		return t.Name
	default:
		return "<" + typeName255r1(e) + ">"
	}
}

func typeName255r1(e ast.Expr) string {
	switch e.(type) {
	case *ast.CallExpr:
		return "CallExpr"
	case *ast.SelectorExpr:
		return "SelectorExpr"
	case *ast.FuncLit:
		return "FuncLit"
	default:
		return "other"
	}
}
