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

// sinkRebindLineHas258 is the reading rule the winlive rebind ruler
// (cmd/wisp/resident_hotkey_live_258_windows_test.go) uses on the raw sink
// tail. Ticket 258-v1 §5 measured the defect it replaces: the rebind record WAS
// in the sink and the judge still failed, because the old grep demanded the
// console equals shape summon=VALUE while internal/ball/hotkey_reload.go's
// Check books the slog JSON shape ("summon":"VALUE", colon form). The judge
// now names both spellings; a tail that carries neither is read as "not
// rebound", which is what the live case is supposed to fail on.
func sinkRebindLineHas258(tail, summonValue string) bool {
	const rebindMsg = "ball: hotkeys rebound after config change"
	if !strings.Contains(tail, rebindMsg) {
		return false
	}
	return strings.Contains(tail, `"summon":"`+summonValue+`"`) ||
		strings.Contains(tail, "summon="+summonValue)
}

// Test258SinkRebindRulerReadsBothSpellings is AC#2 condition (the rebind
// winlive ruler's format mismatch) pinned in CI: a JSON sink record of the
// exact shape hotkey_reload.go writes must read as a rebind, the console
// equals shape must too, and the negative shapes (old value, missing msg,
// empty tail) must NOT. This is also the static counter-proof for the
// "forever-true" mutation: replace the body with return true and the
// third-block assertions below go red, which is what says the live case's
// verdict is carried by this reading and not by the loop.
func Test258SinkRebindRulerReadsBothSpellings(t *testing.T) {
	jsonLine := `{"time":"2026-10-03T21:48:52.272+08:00","level":"INFO","msg":"ball: hotkeys rebound after config change","summon":"Ctrl+Alt+R","mute":"Ctrl+Alt+M","cancel":"Esc","panel":"Ctrl+Alt+P","live":3}`
	consoleLine := `wisp: ball: hotkeys rebound after config change summon=Ctrl+Alt+R mute=Ctrl+Alt+M cancel=Esc panel=Ctrl+Alt+P live=3`
	if !sinkRebindLineHas258(jsonLine, "Ctrl+Alt+R") {
		t.Errorf("the JSON-colon sink record the shipped process actually writes does not read as a rebind; that is the 258-v1 §5 ruler bug still in place")
	}
	if !sinkRebindLineHas258(consoleLine, "Ctrl+Alt+R") {
		t.Errorf("the console equals shape stopped reading (the shape the old grep alone required)")
	}
	// Negative shapes: the reading must be about the NEW value, not about any
	// line that happens to exist.
	if sinkRebindLineHas258(jsonLine, "Ctrl+Alt+Z") {
		t.Errorf("the matcher answered a summon value the record does not carry")
	}
	if sinkRebindLineHas258(`{"msg":"ball: hotkey reload","binding":"x"}`, "Ctrl+Alt+R") {
		t.Errorf("a non-rebind record read as a rebind")
	}
	if sinkRebindLineHas258("", "Ctrl+Alt+R") {
		t.Errorf("an empty sink tail read as a rebind: the live case would go green on a host whose sink never opened")
	}
}
