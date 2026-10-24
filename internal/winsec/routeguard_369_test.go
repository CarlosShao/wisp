// No build tag on purpose: the claim below is about winsec.go, which is the same file on
// every platform, and the linux leg of the winsec gate has to see it too.
//
// 票 369 (leg 369-r1) put this door here because the route's leaf guard had no
// deterministic behavioural instrument left. The mechanism, measured and stated in
// openleaf_guard_322_windows_test.go's header: privateFile asks the steady state first
// (winsec.go:163 resolveString, whose resolver and whose floor both walk the leaf), so a
// leaf link is refused BEFORE preOpenLeafGuard's call site is reached, and the guard only
// ever answers inside the two-syscall gap after that walk. "Did the guard answer in 400
// raced attempts" is therefore a question about the clock; "is the guard called at all,
// and where" is a question about the code. This door answers the second one by reading the
// call graph - the instrument ticket 132's AC#4 door uses, and for its stated reason: a
// line regex counts prose, and every file in this package carries paragraphs that name
// preOpenLeafGuard by hand.
package winsec

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"testing"
)

// routeGate369 is the sequence privateFile has to keep, in source order: the spelling is
// resolved, the leaf is judged, and only then is the parent pinned and the leaf opened.
// The last two names are not a separate claim - they are what makes "before the pin"
// checkable at all.
var routeGate369 = []string{"resolveString", "preOpenLeafGuard", "pinnedParent", "openPinnedLeaf"}

// routeCall369 is one call in privateFile's body, with the two facts the door needs: the
// name, its source position, and whether it stands at the TOP LEVEL of that body.
//
// "Top level" is the point rather than a detail. A guard moved inside an `if` that never
// runs, a loop, or a closure still appears in a grep and still answers "the name is
// mentioned" while the route no longer asks it. The walk up from the call therefore fails
// on any loop, closure, switch or second `if` between it and the body.
type routeCall369 struct {
	name string
	pos  token.Pos
	top  bool
}

func routeCalls369(src string) ([]routeCall369, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, src, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", src, err)
	}
	var fn *ast.FuncDecl
	for _, decl := range file.Decls {
		if d, ok := decl.(*ast.FuncDecl); ok && d.Name != nil && d.Name.Name == "privateFile" {
			fn = d
			break
		}
	}
	if fn == nil || fn.Body == nil {
		return nil, fmt.Errorf("%s declares no privateFile with a body: the route this door reads has moved "+
			"or been renamed", src)
	}
	parent := map[ast.Node]ast.Node{}
	var stack []ast.Node
	var out []routeCall369
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		if n == nil {
			stack = stack[:len(stack)-1]
			return true
		}
		if len(stack) > 0 {
			parent[n] = stack[len(stack)-1]
		}
		stack = append(stack, n)
		if call, ok := n.(*ast.CallExpr); ok {
			if id, ok := call.Fun.(*ast.Ident); ok {
				out = append(out, routeCall369{name: id.Name, pos: call.Pos(), top: atTopLevel369(parent, call, fn.Body)})
			}
		}
		return true
	})
	return out, nil
}

// atTopLevel369 reports whether call's enclosing statements are nothing but the function
// body and at most one `if` - the guard's own `if err := preOpenLeafGuard(resolved); err
// != nil { return err }`.
func atTopLevel369(parent map[ast.Node]ast.Node, call *ast.CallExpr, body *ast.BlockStmt) bool {
	ifs := 0
	for n := ast.Node(call); n != nil; n = parent[n] {
		switch n.(type) {
		case *ast.ForStmt, *ast.RangeStmt, *ast.FuncLit, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
			return false
		case *ast.IfStmt:
			ifs++
			if ifs > 1 {
				return false
			}
		case *ast.BlockStmt:
			return n == ast.Node(body) && ifs <= 1
		}
	}
	return false
}

// routeOrder369 is one reading of the door: the names in the gate's own order, whether
// that order holds, and the two ways it can fail - a call nested out of the top level, or
// a call that is gone.
type routeOrder369 struct {
	order     []string
	sameOrder bool
	nested    string
	missing   string
}

// routeGateOrder369 is factored out of the test so the checker itself can be run against a
// KNOWN-BAD copy of winsec.go: that is how this door's teeth are demonstrated without
// touching the tree - the mutation lives outside the repository, and the checker is asked
// to name the call it can no longer find.
func routeGateOrder369(src string) (routeOrder369, error) {
	calls, err := routeCalls369(src)
	if err != nil {
		return routeOrder369{}, err
	}
	firstTop := map[string]int{}
	seenAny := map[string]bool{}
	for i, c := range calls {
		seenAny[c.name] = true
		if c.top {
			if _, ok := firstTop[c.name]; !ok {
				firstTop[c.name] = i
			}
		}
	}
	out := routeOrder369{}
	idx := map[string]int{}
	for _, name := range routeGate369 {
		if i, ok := firstTop[name]; ok {
			out.order = append(out.order, name)
			idx[name] = i
			continue
		}
		if seenAny[name] {
			out.nested = name
			out.order = append(out.order, name+"(nested)")
			continue
		}
		out.missing = name
		out.order = append(out.order, name+"(MISSING)")
	}
	if out.nested != "" || out.missing != "" {
		return out, nil
	}
	out.sameOrder = idx["resolveString"] < idx["preOpenLeafGuard"] &&
		idx["preOpenLeafGuard"] < idx["pinnedParent"] &&
		idx["preOpenLeafGuard"] < idx["openPinnedLeaf"]
	return out, nil
}

// TestAC3WindowsTheRouteCallsTheLeafGuardBeforeThePin322r1 is the door. The roster is
// printed on every run, green included, because the reading worth keeping is the order
// itself - a red should be the second thing a reviewer sees, not the first.
func TestAC3WindowsTheRouteCallsTheLeafGuardBeforeThePin322r1(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed, so the package directory cannot be located and this door would have " +
			"nothing to read")
	}
	src := filepath.Join(filepath.Dir(thisFile), "winsec.go")
	got, err := routeGateOrder369(src)
	if err != nil {
		t.Fatalf("the route this door reads is gone: %v", err)
	}
	t.Logf("route order in %s: %v", filepath.Base(src), got.order)
	if got.missing != "" {
		t.Errorf("AC#3 RED: privateFile no longer calls %s at all. Deleting the call is the mutation that has "+
			"to turn this door red; a green here means the leaf leg is on the shelf", got.missing)
	}
	if got.nested != "" {
		t.Errorf("AC#3 RED: privateFile calls %s somewhere other than the top level of its body (inside a "+
			"branch, loop or closure), so a path through the route can reach the open without asking the "+
			"guard. A call that only sometimes runs is what 'on the shelf' looks like in source", got.nested)
	}
	if got.missing == "" && got.nested == "" && !got.sameOrder {
		t.Errorf("AC#3 RED: privateFile calls the gate's legs in the order %v, want %v. The leaf guard has to be "+
			"asked AFTER the spelling is resolved (it judges the resolved leaf) and BEFORE the parent is pinned "+
			"and the leaf opened - that ordering IS the half of the claim 'the guard is on the route' that no "+
			"raced count can state", got.order, routeGate369)
	}
}
