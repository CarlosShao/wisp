package main

// Ticket 321 AC#1: the tray / ball / hot key hop's EXECUTOR, pinned at the AST layer
// (case A, the wiring) and at the capability layer (case B, what that executor posts),
// never at the word layer - a nail that only compares the string "toggle" is still a
// word nail: renaming show's body to the name toggle would satisfy it.
//
// The gap: resident_windows.go hands the ball host withPanelHost(func(via string) bool {
// return holder.toggle(via) }) and no case asserted that executor is the holder's TOGGLE;
// rewiring it to holder.show reddened nothing, though ticket 319-r3 gave show to the
// SECOND LAUNCH and to nothing else precisely because a toggle hides a panel already on
// screen. The 反面 this file is NOT: TestBallPanelGesturesReachThePanelThread builds its
// own panelHostHooks and asks whether a gesture REACHES the panel thread - as true of
// show as of toggle, so it stays green under this ticket's mutation (logs:
// .scratch/wisp/probes/321/r1/). Why a source walk: runResident never returns
// (resident_ball_228_test.go:17 names it) and the closure lives in its body, so all a
// test can hold is the AST; this reuses the package's extractors (runResidentBody258,
// productionFiles228) instead of growing a third walker, the rule ticket 294's header
// sets for this same hop shape. Case A reads only the CLOSURE spelling: a method value
// (withPanelHost(holder.toggle)) or a ticket-296-style hoisted executor reddens it by
// name, and whoever makes that change re-adjudicates this file, not the assertion.

import (
	"go/ast"
	"go/token"
	"slices"
	"strconv"
	"testing"
)

const (
	// panelHolderCtor321 is what the holder variable is bound to: the receiver clause
	// names the constructor, not the spelling "holder", so renaming the variable cannot
	// move the ruler. The hook must sit among the installer's arguments, or a
	// withPanelHost call that reaches no ball host would count.
	panelHolderCtor321   = "newPanelHolder"
	ballHostInstaller321 = "startResidentBall"
	panelHostHook321     = "withPanelHost"
	toggleExecutor321    = "toggle"
	toggleRequest321     = "RequestToggle"
)

// panelHop321 is one reading of the assembly hop.
type panelHop321 struct {
	installs int      // startResidentBall calls in runResident
	hops     int      // withPanelHost calls among their arguments
	argShape string   // how the executor argument was spelled, for the red text
	receiver string   // holder variable the closure drives ("" = none found)
	driven   []string // holder executor method names found, sorted
}

// readPanelHop321 walks runResident's AST. No line numbers anywhere: every clause reads
// shape, so an unrelated edit above the assembly cannot churn it.
func readPanelHop321(t *testing.T, body *ast.BlockStmt) panelHop321 {
	t.Helper()
	holderVars := map[string]bool{}
	ast.Inspect(body, func(n ast.Node) bool { // which name holds the holder this boot built
		d, ok := n.(*ast.AssignStmt)
		if !ok || d.Tok != token.DEFINE || len(d.Lhs) != 1 || len(d.Rhs) != 1 {
			return true
		}
		call, ok := d.Rhs[0].(*ast.CallExpr)
		if !ok {
			return true
		}
		id, ok1 := call.Fun.(*ast.Ident)
		lhs, ok2 := d.Lhs[0].(*ast.Ident)
		if ok1 && ok2 && id.Name == panelHolderCtor321 {
			holderVars[lhs.Name] = true
		}
		return true
	})
	if len(holderVars) == 0 {
		t.Fatalf("AC#1 RED: runResident binds no variable from %s(...), so this boot has no panel holder the ball's gestures could drive - the hop is gone, not merely rewired", panelHolderCtor321)
	}
	out := panelHop321{}
	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if id, ok := call.Fun.(*ast.Ident); !ok || id.Name != ballHostInstaller321 {
			return true
		}
		out.installs++
		for _, arg := range call.Args {
			ast.Inspect(arg, func(m ast.Node) bool {
				hop, ok := m.(*ast.CallExpr)
				if !ok {
					return true
				}
				if id, ok := hop.Fun.(*ast.Ident); !ok || id.Name != panelHostHook321 {
					return true
				}
				out.hops++
				if len(hop.Args) != 1 {
					out.argShape = "a call with " + strconv.Itoa(len(hop.Args)) + " argument(s)"
					return true
				}
				lit, ok := hop.Args[0].(*ast.FuncLit)
				if !ok {
					out.argShape = "not a function literal"
					return true
				}
				out.argShape = "a function literal"
				out.receiver, out.driven = drivesHolder321(lit.Body, holderVars)
				return true
			})
		}
		return true
	})
	slices.Sort(out.driven)
	out.driven = slices.Compact(out.driven)
	return out
}

// drivesHolder321 reads which executor methods of the holder a closure body drives.
func drivesHolder321(body *ast.BlockStmt, holderVars map[string]bool) (string, []string) {
	receiver := ""
	var driven []string
	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok2 := call.Fun.(*ast.SelectorExpr)
		if !ok2 {
			return true
		}
		id, ok3 := sel.X.(*ast.Ident)
		if ok3 && holderVars[id.Name] {
			receiver, driven = id.Name, append(driven, sel.Sel.Name)
		}
		return true
	})
	return receiver, driven
}

// TestAC321PanelHopDrivesTheHolderToggle is case A: the wiring.
func TestAC321PanelHopDrivesTheHolderToggle(t *testing.T) {
	body, file := runResidentBody258(t)
	if body == nil {
		t.Fatalf("no runResident found in resident_windows.go (the instrument, not the code)")
	}
	hop := readPanelHop321(t, body)
	// A1 the counts, their own assertion (the 294 rule: a for-every-found clause over
	// the empty set is the always-green fake). hops 0 = no executor reaches the host.
	if hop.installs != 1 || hop.hops != 1 {
		t.Fatalf("AC#1 RED: runResident (%s) has %d %s call(s) reaching %d %s call(s), want 1 and 1 - the panel hop either hands the ball host no executor or hands it two, and this pin cannot say which the gestures drive", file, hop.installs, ballHostInstaller321, hop.hops, panelHostHook321)
	}
	// A2 an executor is driven at all: a closure answering with a constant is the
	// silence every gesture hanging off this line would fall into.
	if len(hop.driven) == 0 {
		t.Fatalf("AC#1 RED: the panel hop in runResident (%s) is %s that drives NO executor on %s (receiver %q) - the tray item, the ball and the panel hot key now execute nothing, whatever this file is named", file, hop.argShape, panelHolderCtor321, hop.receiver)
	}
	// A3 exactly one executor, and it is the toggle. This is the ticket's claim.
	if len(hop.driven) != 1 || hop.driven[0] != toggleExecutor321 {
		t.Fatalf("AC#1 RED: this hop is NOT toggle semantics - runResident (%s) hands the ball host the executor set %v on %s (receiver %q, argument %s), want exactly [%s]. A panel hop that shows instead of flipping is the rewiring ticket 319-r3 split the two executors to forbid: it hid the operator's window on a double click (probes/319/v/03-ac2-live-adjudication.txt shape A).", file, hop.driven, panelHolderCtor321, hop.receiver, hop.argShape, toggleExecutor321)
	}
}

// TestAC321NamedExecutorIsTheToggleRoad is case B: the capability. Case A's method name
// is resolved across this package's production files with the count pinned to 1 (no
// decoy toggle), and that body must post RequestToggle - so a renamed show road reddens.
func TestAC321NamedExecutorIsTheToggleRoad(t *testing.T) {
	files := productionFiles228(t)
	body, file := runResidentBody258(t)
	if body == nil {
		t.Fatalf("no runResident found in resident_windows.go (the instrument, not the code)")
	}
	hop := readPanelHop321(t, body)
	if len(hop.driven) != 1 {
		t.Fatalf("AC#1 RED: the panel hop (%s) drives %v executor(s), want exactly one toggle - no single method here has a capability that could be checked", file, hop.driven)
	}
	name := hop.driven[0]
	posts, decls := []string{}, 0
	for _, f := range files {
		for _, d := range f.Decls {
			fn, ok := d.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || fn.Name.Name != name || fn.Body == nil {
				continue
			}
			decls++
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				if sel, ok := n.(*ast.CallExpr); ok {
					if e, ok := sel.Fun.(*ast.SelectorExpr); ok {
						posts = append(posts, e.Sel.Name)
					}
				}
				return true
			})
		}
	}
	if decls != 1 {
		t.Fatalf("AC#1 RED: %d method(s) named %s in this package's production files, want exactly 1 - the executor the panel hop drives is not a single defined thing, so its capability cannot be attributed", decls, name)
	}
	if !slices.Contains(posts, toggleRequest321) {
		t.Fatalf("AC#1 RED: this hop is NOT toggle semantics by capability either - the executor %s the panel hop drives posts %v, which lacks %s: the tray item, the ball and the panel hot key no longer flip the panel's state", name, posts, toggleRequest321)
	}
}
