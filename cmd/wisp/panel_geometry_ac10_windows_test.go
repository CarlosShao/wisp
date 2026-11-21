//go:build windows

package main

// Ticket 33 AC#10's geometry seam rulers: [panel] width/height/scale really
// reach the window, through a seam a test can watch WITHOUT opening one.
//
// WHAT EACH RULER OWES (the dispatch's three teeth, each with its red shape):
//
//  1. TestAC10PanelWindowGeometryMapsTheThreeKeys - the pure mapping, value by
//     value: the schema's 640 default, a written 720x480, height 0 answering
//     the standing 255 ruling (NOT auto-height), garbage falling to the host
//     constants, and the zoom factor passing through to panelScaledSize. The
//     red shape is any second spelling of these rules downstream: the mapping
//     is single-sourced here and the production path consumes THIS function
//     (panelGeometryTripleSource), so a drifted copy reddens the hot ruler
//     below while this one stays green.
//  2. TestAC10PanelKeysRewrittenInTheSameProcessReachTheNextShow - the hot
//     tooth: the guarded setters rewrite config.toml, and the SAME host, in
//     the SAME process, resolves the new triple on its next resolution and
//     posts the zoomed pair on the next show request (SetSize is the
//     SetWindowPos hop; "一次 SetWindowPos 就体现"). The red shape is a
//     snapshot-at-assembly source or a zoom folded anywhere but the resolution
//     point: either leaves the old numbers in the sink.
//  3. TestAC10RecorderProvesTheProductionPathCallsApplyGeometry - the
//     anti-vacuous tooth. AC#10's disease reading was "生产读取方现量 = 各 0";
//     a recorder injected into the real seam field must see the production
//     path deliver the numbers it resolved (720 x 450 x 1.25 from the file on
//     disk), the recorder count being 0 BEFORE the hop and exactly the
//     resolved triple AFTER it. A test that only grepped the field name would
//     pass with the hop deleted - the AST half below pins the hop into
//     bringUp, so deleting the call reddens the structure ruler and an empty
//     recorder reddens this one.
//
// NO REAL WINDOW IS OPENED ANYWHERE HERE. Every reading is taken from the
// pure mapping, the production source closure over a real config.toml, a
// recorded fake control (the 255r1 sink), or the parsed host file - the same
// headless discipline panel_geometry_255_test.go and
// panel_reshow_255r1_windows_test.go run under. The on-screen half of scale
// (the frame a person sees) is a winlive reading and this file does not
// pretend to have one.
//
// PATH warning (ticket 98): this package's test binary links sherpa-onnx and
// dies at load (0xc0000135) unless third_party/sherpa-onnx is on PATH - which
// is why these cases run through scripts/wisp-cli-staged-tests.sh.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
	webview2 "github.com/jchv/go-webview2"
)

// writePanelGeometryAc10 plants a loadable config.toml with the given [panel]
// rows. schema_version is current so readConfigFile takes no migration branch
// (that branch rewrites the file, and a rewrite from under a re-read test
// would be a second mover in a measurement about who read what when).
func writePanelGeometryAc10(t *testing.T, dir string, rows string) string {
	t.Helper()
	body := "schema_version = 2\n\n[panel]\n" + rows
	path := filepath.Join(dir, configFileName)
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
	return path
}

// --- tooth 1: the pure mapping, value by value ---------------------------------

func TestAC10PanelWindowGeometryMapsTheThreeKeys(t *testing.T) {
	cases := []struct {
		name          string
		width, height int
		scale         float64
		wantW, wantH  int
		wantScale     float64
	}{
		{
			name:  "the schema's own 640 default maps through with a silent height and scale",
			width: 640, height: 0, scale: 0,
			wantW: 640, wantH: panelHeightPx, wantScale: 0,
		},
		{
			name:  "width, height and scale written: value for value",
			width: 720, height: 480, scale: 1.25,
			wantW: 720, wantH: 480, wantScale: 1.25,
		},
		{
			// The standing 255 ruling, restated where the mapping lives: height 0
			// is NOT auto-height, it answers the host constant with 由内容定高
			// registered as NOT DONE - this mapping invents no content metric.
			name:  "height 0 answers the host constant, not an invented auto-height",
			width: 720, height: 0, scale: 1.25,
			wantW: 720, wantH: panelHeightPx, wantScale: 1.25,
		},
		{
			name:  "a zero width falls to the host constant, not to a guess",
			width: 0, height: 400, scale: 0,
			wantW: panelWidthPx, wantH: 400, wantScale: 0,
		},
		{
			name:  "garbage answers are not a second spelling of auto",
			width: -1, height: -1, scale: -1,
			wantW: panelWidthPx, wantH: panelHeightPx, wantScale: 0,
		},
		{
			name:  "a fraction zoom is legal and passes through: the applier owns the multiply",
			width: 640, height: 480, scale: 0.5,
			wantW: 640, wantH: 480, wantScale: 0.5,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w, h, s := panelWindowGeometry(config.PanelSection{Width: tc.width, Height: tc.height, Scale: tc.scale})
			if w != tc.wantW || h != tc.wantH || s != tc.wantScale {
				t.Fatalf("panelWindowGeometry(%d, %d, %v) = (%d, %d, %v), want (%d, %d, %v): the mapping the production source runs must answer these three keys value by value",
					tc.width, tc.height, tc.scale, w, h, s, tc.wantW, tc.wantH, tc.wantScale)
			}
		})
	}

	// The zoom the applier and the re-show fold carry: a written 1.25 turns the
	// resolved pair into the round(w*s) x round(h*s) frame; scale 0 - follow
	// system DPI - answers the pair unchanged, so every default config's frame
	// stays byte-identical to what it was before AC#10.
	if bw, bh := panelScaledSize(720, 480, 1.25); bw != 900 || bh != 600 {
		t.Fatalf("panelScaledSize(720, 480, 1.25) = %dx%d, want 900x600", bw, bh)
	}
	if bw, bh := panelScaledSize(640, 260, 0); bw != 640 || bh != 260 {
		t.Fatalf("panelScaledSize(640, 260, 0) = %dx%d, want the unchanged 640x260 (follow DPI multiplies nothing)", bw, bh)
	}
	if bw, bh := panelScaledSize(640, 260, -1); bw != 640 || bh != 260 {
		t.Fatalf("panelScaledSize(640, 260, -1) = %dx%d, want the unchanged 640x260 (garbage is not a zoom)", bw, bh)
	}
}

// --- tooth 2: rewritten keys reach the next resolution in the same process -----

func TestAC10PanelKeysRewrittenInTheSameProcessReachTheNextShow(t *testing.T) {
	dir := t.TempDir()
	path := writePanelGeometryAc10(t, dir, "width = 640\n")
	mgr, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("manager over the plant: %v", err)
	}
	m := NewPanelManager(&panel.ComposerDispatch{}, nil, "", withGeometryTriple(panelGeometryTripleSource(dir)))

	opts := m.windowOptions()
	if opts.Width != 640 || opts.Height != panelHeightPx {
		t.Fatalf("the first resolution got %dx%d, want the planted 640 and the silent-height constant %d",
			opts.Width, opts.Height, panelHeightPx)
	}

	// 改键: the guarded setters land width, height AND scale - file and memory
	// kept equal by the guarded writer - and nothing restarts in between.
	if written, err := mgr.SetPanelWidth(720); err != nil || len(written) != 1 || written[0] != "panel.width" {
		t.Fatalf("SetPanelWidth(720) = %v, %v, want [panel.width], nil", written, err)
	}
	if written, err := mgr.SetPanelHeight(480); err != nil || len(written) != 1 || written[0] != "panel.height" {
		t.Fatalf("SetPanelHeight(480) = %v, %v, want [panel.height], nil", written, err)
	}
	if written, err := mgr.SetPanelScale(1.25); err != nil || len(written) != 1 || written[0] != "panel.scale" {
		t.Fatalf("SetPanelScale(1.25) = %v, %v, want [panel.scale], nil", written, err)
	}

	// 同一进程的下一次解析: the source re-reads the file the setters wrote, the
	// mapping answers the new triple - a snapshot-at-assembly source would
	// answer 640x260 again and redden here.
	opts = m.windowOptions()
	if opts.Width != 720 || opts.Height != 480 {
		t.Fatalf("after the writes the next resolution got %dx%d, want 720x480: the keys were snapshotted, not re-read, which contradicts [panel] being registered hot",
			opts.Width, opts.Height)
	}

	// 一次 SetWindowPos 就体现: the re-show hop posts the zoomed pair through
	// the library's SetSize (AdjustWindowRect + SetWindowPos), inside Dispatch,
	// under HintNone - the shape the 255r1 rulers pin for this hop.
	sink := newGeoSink255r1()
	attachGeoSink255r1(m, sink)
	m.requestGeometryOnReshow()
	sized, bare, _ := sink.snapshot()
	if len(bare) != 0 {
		t.Fatalf("%d SetSize call(s) went out without Dispatch: %+v", len(bare), bare)
	}
	if len(sized) != 1 {
		t.Fatalf("the re-show sent %d SetSize request(s), want exactly 1: %+v", len(sized), sized)
	}
	if sized[0].width != 900 || sized[0].height != 600 {
		t.Fatalf("the zoomed re-show asked for %dx%d, want the scaled 900x600 (720x480 at 1.25): scale never reached the SetWindowPos path",
			sized[0].width, sized[0].height)
	}
	if !sized[0].insideDispatch || sized[0].hint != webview2.HintNone {
		t.Fatalf("the zoomed resize went out as %+v, want inside Dispatch under HintNone - the hop's own pinned shape", sized[0])
	}

	// And back to follow-DPI: scale 0 folds nothing, the pair goes out as the
	// resolution answered it.
	if _, err := mgr.SetPanelScale(0); err != nil {
		t.Fatalf("SetPanelScale(0): %v", err)
	}
	m.requestGeometryOnReshow()
	sized, _, _ = sink.snapshot()
	if len(sized) != 2 {
		t.Fatalf("the second re-show sent %d SetSize request(s), want 2: %+v", len(sized), sized)
	}
	if sized[1].width != 720 || sized[1].height != 480 {
		t.Fatalf("the follow-DPI re-show asked for %dx%d, want the unchanged 720x480", sized[1].width, sized[1].height)
	}
}

// --- tooth 3: the recorder watches the production path deliver -----------------

func TestAC10RecorderProvesTheProductionPathCallsApplyGeometry(t *testing.T) {
	dir := t.TempDir()
	writePanelGeometryAc10(t, dir, "width = 720\nheight = 450\nscale = 1.25\n")
	var got []panelResolved
	m := NewPanelManager(&panel.ComposerDispatch{}, nil, "",
		withGeometryTriple(panelGeometryTripleSource(dir)),
		withApplyGeometry(func(w, h int, scale float64) {
			got = append(got, panelResolved{w: w, h: h, scale: scale, ok: true})
		}))

	// The =0 reading, kept as a live assertion: BEFORE the production hops run
	// the recorder count is 0. If the seam wiring were deleted this stays 0
	// and this case reddens instead of passing vacuously.
	if len(got) != 0 {
		t.Fatalf("the recorder saw %d delivery(ies) before any production hop ran: %+v", len(got), got)
	}

	// The create path's two hops, in the order bringUp makes them: the create's
	// geometry argument resolves the triple (and caches it), the tail hop hands
	// it to the seam.
	m.windowOptions()
	m.applyResolvedGeometry()

	if len(got) != 1 {
		t.Fatalf("the recorder saw %d delivery(ies), want exactly 1: the production path never crossed the seam", len(got))
	}
	if got[0].w != 720 || got[0].h != 450 || got[0].scale != 1.25 {
		t.Fatalf("the production path delivered (%d, %d, %v), want (720, 450, 1.25) - the numbers AC#10 resolved from the file on disk must be the numbers the seam receives",
			got[0].w, got[0].h, got[0].scale)
	}

	// The default seam is the production value, not a nil the hop skips around:
	// a manager nobody handed an applier to still crosses into Win32.
	m2 := NewPanelManager(&panel.ComposerDispatch{}, nil, "")
	if m2.applyGeometry == nil {
		t.Fatal("NewPanelManager left applyGeometry nil: in production the create tail's hop would call nothing at all")
	}

	// Structural half, parsed not grepped: bringUp really makes BOTH hops and
	// the production value really reaches SetWindowPos through the zoom helper.
	// Deleting either reddens here, naming the hop that moved.
	file := parseHostFile255r1(t)
	bringUp := funcDecl255r1(t, file, "bringUp")
	if n := countCall255r1(bringUp, "applyResolvedGeometry"); n != 1 {
		t.Fatalf("bringUp calls applyResolvedGeometry %d time(s), want exactly 1: the create tail's applier hop was moved or removed and tooth 3 now measures a method nothing calls", n)
	}
	if m.geometryTriple == nil {
		t.Fatal("the test manager lost its triple source; the reading above would not be aimed")
	}
	// applyGeometryWin32 lives in the seam file, so the structural half parses
	// THAT file the way parseHostFile255r1 parses the host's. The counter here
	// matches BOTH call shapes: bare identifiers (panelScaledSize(...)) and
	// the package var the Win32 hop rides (pnlSetWindowPos.Call(...)).
	seamRaw, err := os.ReadFile(filepath.Join(hostPkgDir255(t), "panel_geometry_ac10.go"))
	if err != nil {
		t.Fatalf("read the seam file: %v", err)
	}
	seamFile, perr := parser.ParseFile(token.NewFileSet(), "panel_geometry_ac10.go", seamRaw, parser.AllErrors)
	if perr != nil {
		t.Fatalf("parse panel_geometry_ac10.go: %v", perr)
	}
	applier := funcDecl255r1(t, seamFile, "applyGeometryWin32")
	countCallAc10 := func(root ast.Node, name string) int {
		n := 0
		ast.Inspect(root, func(node ast.Node) bool {
			if call, ok := node.(*ast.CallExpr); ok {
				switch fun := call.Fun.(type) {
				case *ast.Ident:
					if fun.Name == name {
						n++
					}
				case *ast.SelectorExpr:
					if fun.Sel != nil && fun.Sel.Name == name {
						n++
					}
					if id, ok := fun.X.(*ast.Ident); ok && id.Name == name {
						n++
					}
				}
			}
			return true
		})
		return n
	}
	if n := countCallAc10(applier, "pnlSetWindowPos"); n != 1 {
		t.Fatalf("applyGeometryWin32 rides pnlSetWindowPos %d time(s), want exactly 1: scale's Win32 exit moved", n)
	}
	if n := countCallAc10(applier, "panelScaledSize"); n != 1 {
		t.Fatalf("applyGeometryWin32 calls panelScaledSize %d time(s), want exactly 1: the zoom multiply must be the single-sourced helper", n)
	}
}

// The 323 reading, extended to the third key: the production assembly wires the
// wide source, so a host born from newResidentPanelManager resolves scale from
// the config on disk - and a silent scale answers follow-DPI (0), never an
// invented zoom. Value asserted against a planted non-default triple so the
// only possible source is the assembly's own wiring.
func TestAC10AssemblyWiresTheTripleSourceScaleFollowsTheDisk(t *testing.T) {
	dir := newResidentAssemblyDir323(t)
	// The 323 plant carries width/height only; extend it with a written scale.
	path := filepath.Join(dir, configFileName)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read the 323 plant: %v", err)
	}
	extended := strings.Replace(string(raw), "height = 534\n", "height = 534\nscale = "+strconv.FormatFloat(1.25, 'f', -1, 64)+"\n", 1)
	if extended == string(raw) {
		t.Fatalf("the scale plant changed nothing")
	}
	if err := os.WriteFile(path, []byte(extended), 0o600); err != nil {
		t.Fatalf("rewrite the 323 plant: %v", err)
	}

	mgr, _, _, _, err := newResidentPanelManager(dir, panelModeConfirm((&countingPermConfirm{}).fn()))
	if err != nil {
		t.Fatalf("the resident assembly failed: %v", err)
	}
	opts := mgr.windowOptions()
	// The create pair stays the config's own 812x534 (the zoom multiplies at the
	// applier/SetWindowPos hop and the re-show fold, not inside the pair) - and
	// the cache the tail hop reads carries the written factor.
	if opts.Width != 812 || opts.Height != 534 {
		t.Fatalf("the assembled host resolves %dx%d, want the planted 812x534: the wide wiring changed the create pair's semantics",
			opts.Width, opts.Height)
	}
	mgr.mu.Lock()
	r := mgr.resolved
	mgr.mu.Unlock()
	if !r.ok || r.w != 812 || r.h != 534 || r.scale != 1.25 {
		t.Fatalf("the assembled host's cached triple = %+v, want (812, 534, 1.25): scale did not ride the one read the assembly wired", r)
	}
}
