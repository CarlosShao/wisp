package main

// Ticket 258: the assembly-root shape walk, in the untagged file the shape
// rulers of this package live in (the walk reads source, so it needs no
// Windows build).
//
// Test258AssemblyRootWiresTheChainAndTheBridge reads runResident's body the way
// resident_ball_228_test.go does and requires the [hotkey] chain to be part of
// the production call: startResidentBall must be called with at least four
// arguments (the 3rd and 4th being the hotCfg / hotReload closures), and
// resident_windows.go must still read config.toml for them. Deleting either
// closure from the assembly root reddens the matching clause - the ball cannot
// silently drop back to compiled literals.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func runResidentBody258(t *testing.T) (*ast.BlockStmt, string) {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filepath.Join(".", "resident_windows.go"), nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("parsing resident_windows.go (the instrument, not the code): %v", err)
	}
	for _, d := range f.Decls {
		if fn, ok := d.(*ast.FuncDecl); ok && fn.Recv == nil && fn.Name.Name == "runResident" {
			return fn.Body, "resident_windows.go"
		}
	}
	return nil, ""
}

func Test258AssemblyRootWiresTheChainAndTheBridge(t *testing.T) {
	body, file := runResidentBody258(t)
	if body == nil {
		t.Fatalf("no runResident found in resident_windows.go (the instrument, not the code)")
	}
	// (a) startResidentBall is called with four positional arguments: reg, the
	// cancel executor, the hotCfg closure, the hotReload closure. A call that
	// dropped back to the two-argument pre-258 shape is the defect this walk
	// exists to catch.
	calls, hotArgs := 0, 0
	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if id, ok := call.Fun.(*ast.Ident); !ok || id.Name != "startResidentBall" {
			return true
		}
		calls++
		if len(call.Args) >= 4 {
			hotArgs++
		}
		return true
	})
	if calls == 0 {
		t.Fatalf("runResident (%s) no longer calls startResidentBall: ticket 228's own walk will be red too", file)
	}
	if hotArgs == 0 {
		t.Fatalf("AC#1 RED: runResident (%s) calls startResidentBall without the [hotkey] chain and the bridge src (arguments 3 and 4) - the ball is back on compiled literals", file)
	}

	// (b) The closures read config.toml - the [hotkey] source the chain is named
	// for. A closure that reads nothing proves nothing.
	src, err := os.ReadFile(filepath.Join(".", "resident_windows.go"))
	if err != nil {
		t.Fatalf("read resident_windows.go (the instrument): %v", err)
	}
	if !strings.Contains(string(src), "config.LoadFile") {
		t.Errorf("resident_windows.go no longer reads config.toml for the [hotkey] chain: the construction source moved, re-adjudicate this walk")
	}
}
