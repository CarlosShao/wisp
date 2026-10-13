//go:build windows

package main

// Ticket 255 r6 - the range behind ruler 3's "exactly one" claim, on disk.
//
// 255-v2 measured that the claim outran the ruler: cmd/wisp/panel_geometry_255_test.go
// demanded exactly one webview2.NewWithOptions call site while reading ONE file
// (panel_host_windows.go), and the same one-file blindness was ruled for the ticket
// family's other "only one" rulers. A713's method ruling says a source-reading ruler
// is only proven by planting the wrong shape ON DISK - `go test -overlay` swaps the
// compiler's eyes, never the ruler's - so the proof lives in
// .scratch/wisp/probes/255/r6/impl.md with mutation logs next to it.
//
// RETOOLED BY 299-r1 (2026-10-13): ticket 299 AC#1 moved the create site off the
// library wrapper - the host now creates its own window (newPanelControl,
// cmd/wisp/panel_control_windows.go), so the "exactly one create site" census counts
// THAT call instead of webview2.NewWithOptions. The claim is unchanged and so is the
// range 255-v2 ruled for it: exactly one place in this package's shipped code where
// a panel window is born. The census matches Ident and Selector callees both, so
// neither `newPanelControl(...)` nor a receiver-qualified spelling slips past it.
//
// The range chosen here is this package's PRODUCTION SOURCES (every .go in the
// directory minus _test.go), the same census subject
// panel_locked_naming_33r11_windows_test.go:packageSourceFiles33r11 judges. Two
// named limits, both deliberate: _test.go files are out (shipped code owns the
// promise, and fixtures would red-lock the ruler the way 33r11 says), and files in
// other packages are out (scripts/spike/webview2-latency/main.go:149 keeps its own
// wrapper create, so a repo-wide "exactly 1" would be baseline-red - the ruler would
// then have to ship red or lie about its subject; see impl.md section 1).

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

// panelCreateCall299r1 is the callee name of the ONE window-create site this
// package ships, read by panel_geometry_255_test.go's ruler 3 through this file's
// census. It names the function ticket 299-r1 moved the create into; before that
// ticket the same ruler counted webview2.NewWithOptions.
const panelCreateCall299r1 = "newPanelControl"

// packageProductionSources255r6 is the census subject, borrowed shape for shape
// from packageSourceFiles33r11: ReadDir of the package directory, files only,
// .go only, _test.go excluded, sorted, and an empty list is a failure - a census
// that read no files proves nothing and must never report a clean zero.
func packageProductionSources255r6(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("the 255r6 census could not read its own package directory %s: %v", dir, err)
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
		t.Fatalf("the 255r6 census found zero .go files in %s: the ruler has no subject", dir)
	}
	return names
}

// calleeName299r1 answers the name a call expression goes through, for both
// `name(...)` (this package's create call today) and `recv.name(...)`. An empty
// answer means "not a named call", which the census never counts.
func calleeName299r1(call *ast.CallExpr) string {
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		return fun.Name
	case *ast.SelectorExpr:
		if fun.Sel != nil {
			return fun.Sel.Name
		}
	}
	return ""
}

// createSitesInPackage299r1 returns every call site whose callee name is
// panelCreateCall299r1 across the package's production sources, as "file:line" so a
// red names the intruding file instead of only a count. It parses into the CALLER's
// fset so positions stay comparable with the host parse already in hand. Two shape
// decisions, both load-bearing: a source that cannot be read or parsed reddens the
// census through t.Fatalf instead of being skipped (a silently skipped file is how a
// range ruler goes blind), and the match is on the callee's final name only - for a
// selector call the receiver identifier is not checked - so an aliased or
// receiver-qualified spelling lands inside the range rather than slipping past it.
func createSitesInPackage299r1(t *testing.T, fset *token.FileSet, dir string) []string {
	t.Helper()
	var sites []string
	for _, name := range packageProductionSources255r6(t, dir) {
		path := filepath.Join(dir, name)
		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("the 255r6 census could not read %s: %v", path, err)
		}
		file, err := parser.ParseFile(fset, path, src, parser.SkipObjectResolution)
		if err != nil {
			t.Fatalf("the 255r6 census could not parse %s - a source the ruler cannot read is not an exempt source: %v", path, err)
		}
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			if calleeName299r1(call) == panelCreateCall299r1 {
				sites = append(sites, fmt.Sprintf("%s:%d", name, fset.Position(call.Pos()).Line))
			}
			return true
		})
	}
	sort.Strings(sites)
	return sites
}
