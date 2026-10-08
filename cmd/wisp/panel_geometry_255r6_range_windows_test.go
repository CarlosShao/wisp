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
// The range chosen here is this package's PRODUCTION SOURCES (every .go in the
// directory minus _test.go), the same census subject
// panel_locked_naming_33r11_windows_test.go:packageSourceFiles33r11 judges. Two
// named limits, both deliberate: _test.go files are out (shipped code owns the
// promise, and fixtures would red-lock the ruler the way 33r11 says), and files in
// other packages are out (scripts/spike/webview2-latency/main.go:149 is a real
// second call site, so a repo-wide "exactly 1" would be baseline-red - the ruler
// would then have to ship red or lie about its subject; see impl.md section 1).

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

// newWithOptionsSitesInPackage255r6 returns every call site whose selector name
// is NewWithOptions across the package's production sources, as "file:line" so a
// red names the intruding file instead of only a count. It parses into the CALLER's
// fset so positions stay comparable with the host parse already in hand. Two shape
// decisions, both load-bearing: a source that cannot be read or parsed reddens the
// census through t.Fatalf instead of being skipped (a silently skipped file is how a
// range ruler goes blind), and the match is on the selector's final name only - the
// receiver identifier is not checked - so an aliased import (wv2.NewWithOptions
// with webview2 dodged) lands inside the range rather than slipping past it.
func newWithOptionsSitesInPackage255r6(t *testing.T, fset *token.FileSet, dir string) []string {
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
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if sel.Sel.Name == "NewWithOptions" {
				sites = append(sites, fmt.Sprintf("%s:%d", name, fset.Position(call.Pos()).Line))
			}
			return true
		})
	}
	sort.Strings(sites)
	return sites
}
