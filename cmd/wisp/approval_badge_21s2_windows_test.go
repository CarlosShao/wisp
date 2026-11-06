//go:build windows

package main

// Ticket 21 segment 2: the ball's approval-queue depth badge.
//
// TWO HALVES, TWO RULERS, stated apart:
//
//   - the mapping (approvalDepthBadge) is a pure function and is tested
//     directly below;
//   - the POSTING (ballCardUI.Prompt calling ball.SetBadge after SetState,
//     Update calling syncBadge on a dismissal) needs a real Win32 orb window
//     to observe - a live rig this leg must not run (the 260-r3 precedent:
//     real-window verification belongs to the owner-present winlive pass).
//     What IS provable headless is that the wiring exists and stays wired,
//     which is the AST census below, the same instrument
//     resident_approval_risk_256_windows_test.go uses for approval.Options.
//     Deleting either posting turns the census red; that is its tooth.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent/approval"
)

// TestTicket21S2ApprovalDepthBadgeMapping is the pure half of SPEC-06 §7's
// 「深度角标」: only an L2 queue card has a depth, and a non-positive depth
// prints no mark at all rather than a "0".
func TestTicket21S2ApprovalDepthBadgeMapping(t *testing.T) {
	cases := []struct {
		name string
		p    approval.Prompt
		want int
	}{
		{"L2 depth 1", approval.Prompt{Level: "L2", Depth: 1}, 1},
		{"L2 depth 3", approval.Prompt{Level: "L2", Depth: 3}, 3},
		{"L2 depth 0 shows nothing", approval.Prompt{Level: "L2", Depth: 0}, 0},
		{"L2 negative depth shows nothing", approval.Prompt{Level: "L2", Depth: -2}, 0},
		{"L1 window has no queue depth", approval.Prompt{Level: "L1", Depth: 4}, 0},
		{"empty level shows nothing", approval.Prompt{Level: "", Depth: 2}, 0},
	}
	for _, c := range cases {
		if got := approvalDepthBadge(c.p); got != c.want {
			t.Errorf("%s: approvalDepthBadge = %d, want %d", c.name, got, c.want)
		}
	}
}

// census21S2 walks one file and returns, per top-level function, the selector
// names it calls (SetBadge / SetState / syncBadge / RenderCard) and how often.
func census21S2(t *testing.T, path string) map[string]map[string]int {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	out := map[string]map[string]int{}
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		calls := map[string]int{}
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			switch sel := call.Fun.(type) {
			case *ast.SelectorExpr:
				calls[sel.Sel.Name]++
			case *ast.Ident:
				calls[sel.Name]++
			}
			return true
		})
		out[fn.Name.Name] = calls
	}
	return out
}

// TestTicket21S2BadgeWiringCensus is the headless tooth: the orb's Prompt must
// still post the depth badge AFTER the state entry (applyStateLocked zeroes
// the badge on every state change, and the STA queue runs posted tasks in
// order - so the order below is load-bearing, not stylistic), Update must
// re-sync on the two card-ending events, and the render call must still be the
// shared renderer.
func TestTicket21S2BadgeWiringCensus(t *testing.T) {
	census := census21S2(t, "resident_approval_windows.go")

	prompt := census["Prompt"]
	if prompt["SetBadge"] != 1 {
		t.Errorf("ballCardUI.Prompt calls SetBadge %d times, want exactly 1 (the depth post)", prompt["SetBadge"])
	}
	if prompt["SetState"] != 1 {
		t.Errorf("ballCardUI.Prompt calls SetState %d times, want exactly 1", prompt["SetState"])
	}
	if prompt["RenderCard"] != 1 {
		t.Errorf("ballCardUI.Prompt calls RenderCard %d times, want exactly 1 - the card is the shared renderer's, not a second copy", prompt["RenderCard"])
	}
	update := census["Update"]
	if update["syncBadge"] != 1 {
		t.Errorf("ballCardUI.Update calls syncBadge %d times, want exactly 1 (the dismissal re-sync)", update["syncBadge"])
	}
}

// TestTicket21S2BadgePostsAreOrderedStateThenDepth is the order tooth: inside
// Prompt, the SetBadge call must appear AFTER the SetState call in statement
// order, because the state entry zeroes the badge and the depth post must win.
// A refactor that swaps them produces a badge that flashes and disappears.
func TestTicket21S2BadgePostsAreOrderedStateThenDepth(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "resident_approval_windows.go", nil, 0)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	var setStateLine, setBadgeLine int
	ast.Inspect(file, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "Prompt" || fn.Recv == nil {
			return true
		}
		ast.Inspect(fn.Body, func(m ast.Node) bool {
			call, ok := m.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			line := fset.Position(call.Pos()).Line
			switch sel.Sel.Name {
			case "SetState":
				setStateLine = line
			case "SetBadge":
				setBadgeLine = line
			}
			return true
		})
		return false
	})
	if setStateLine == 0 || setBadgeLine == 0 {
		t.Fatalf("premise broke: Prompt no longer posts both (SetState@%d SetBadge@%d)", setStateLine, setBadgeLine)
	}
	if setBadgeLine < setStateLine {
		t.Errorf("SetBadge (line %d) fires before SetState (line %d); the state entry zeroes the badge, "+
			"so the depth post must come after it", setBadgeLine, setStateLine)
	}
}
