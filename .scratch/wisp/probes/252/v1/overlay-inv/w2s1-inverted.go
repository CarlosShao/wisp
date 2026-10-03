package tools

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// Ticket 252 R2 leg (252-r2): an instrument for ONE hard ban, spelled verbatim
// in the ticket at AC#2:
//
//	"绝不许把两次包含改成一次" (the two containments may not become one)
//
// That ban is ticket 107's hole, and as of anchor e16e3037 nothing in the
// repository measured it: this leg re-ran the deletion of the lexical leg
// (M-merge, see .scratch/wisp/probes/252/r2/mut/) and internal/tools answered
// 257 PASS / 0 FAIL. The mechanism is not a missing case, it is a射程 that no
// value-taking ruler can cover - ticket 252's own fix (2b1a3071) moved the
// same-form step into Canonicalize, so every string production hands
// InAllowlist is already one form on both sides, LEG1 and LEG2 agree on all of
// them, and the identity final == leg1 && leg2 that
// paths_shortname_252_probe_test.go computes can no longer tell a two-leg
// judge from a one-leg one.
//
// So this file holds the two shapes that DO have teeth, and they are
// complementary rather than redundant:
//
//	W2 - a VALUE witness for the resolvedForm leg. One input on which the two
//	     legs answer differently (lexically inside the root, resolved outside
//	     it), fed straight to InAllowlist. Delete the resolvedForm leg and this
//	     goes red on the spot.
//
//	S1 - a CALL-GRAPH witness for both legs at once, taken from InAllowlist's
//	     own body rather than from any return value. It is the only one of the
//	     three that fires on "two calls, one check" (M-samesource, measured:
//	     257 PASS / 0 FAIL without it), it needs no OS cooperation at all (no
//	     8.3 namespace, no mklink, no volume with aliases enabled), and it
//	     names the ban in its own red line.
//
// The lexical leg's matching value witness is W1, in
// paths_twocontainments_252_r2_windows_test.go: it needs a real 8.3 short name,
// which only Windows has, so it lives behind that build tag instead of skipping
// here.
//
// What this leg does NOT do: it changes no production line (paths.go is read as
// data, byte for byte, by S1), it widens no containment, and it touches no
// threshold, golden or C18 constant.

// ---------------------------------------------------------------------------
// W2: the resolvedForm leg's value witness.
// ---------------------------------------------------------------------------

// legsOf252r2 reads InAllowlist's two containments the way the judge computes them,
// so a red line can say WHICH leg did the refusing instead of "looks wrong".
// This is a reading aid, not the ruler: only pc.InAllowlist's own boolean is
// asserted.
func legsOf252r2(pc *PathCanonicalizer, canonical string) (leg1, leg2 bool, rf string, ok bool) {
	roots := pc.Roots()
	leg1 = rootsContain(roots, foldPath(canonical))
	rf, ok = resolvedForm(canonical)
	leg2 = ok && rootsContain(roots, foldPath(rf))
	return leg1, leg2, rf, ok
}

// TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses is W2.
//
// The carrier is ticket 107b probe C's shape, rebuilt here so this ruler owns
// its own premise: a real allowed root, and inside it a directory link that
// routes to a tree outside it. A path reached through that link is
//
//	LEXICALLY inside the root   -> leg 1 says yes
//	RESOLVED  outside the root  -> leg 2 says no
//
// which is the one input shape where only the second containment can refuse.
// Delete it and this case goes red; the first containment has nothing to say
// about it. That asymmetry is the whole point of "requiring both is strictly
// narrower" (paths.go:199-207), and it is why merging the two into one is a
// hole and not a cleanup.
//
// No skip path: if the link cannot be built, or no platform-visible signal
// reports it as a link, or either leg stops answering as this case claims, the
// test fails loudly rather than passing vacuously (makeDirLink107b /
// dirLinkEvidence107b already enforce the first two).
func TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses(t *testing.T) {
	base := sealableTempDir124(t)
	proj := filepath.Join(base, "proj")
	outside := filepath.Join(base, "outside")
	for _, d := range []string{proj, outside} {
		if err := os.MkdirAll(d, 0o700); err != nil {
			t.Fatalf("mkdir %q: %v", d, err)
		}
	}
	legit := filepath.Join(proj, "a.txt")
	if err := os.WriteFile(legit, []byte("wisp"), 0o600); err != nil {
		t.Fatalf("write legit: %v", err)
	}
	esc := filepath.Join(proj, "esc")
	makeDirLink107b(t, esc, outside)
	landed := dirLinkEvidence107b(t, esc, outside)

	pc := NewPathCanonicalizer([]string{proj}, nil)
	if u := pc.UnusableRoots(); len(u) != 0 {
		t.Fatalf("carrier broken: the configured root was dropped: %q", u)
	}
	if r := pc.Roots(); len(r) != 1 {
		t.Fatalf("carrier broken: roots = %q, want the one named tree", r)
	}

	// Carrier control, the positive direction: an ordinary file in the named
	// root goes through the full production round trip and must come back
	// authorized. Without this line a ruler that refuses everything would look
	// like a working ruler.
	cLegit, err := pc.Canonicalize(legit)
	if err != nil {
		t.Fatalf("carrier broken: Canonicalize(%q): %v", legit, err)
	}
	if leg1, leg2, rf, ok := legsOf252r2(pc, cLegit); !leg1 || !leg2 || !pc.InAllowlist(cLegit) {
		t.Errorf("POSITIVE CONTROL RED for W2: the plain file inside the allowed root is not authorized "+
			"(canonical=%q leg1=%v leg2=%v resolved=%q ok=%v roots=%q), so nothing this case says about the "+
			"escaping path can be trusted", cLegit, leg1, leg2, rf, ok, pc.Roots())
	}

	// The disagreeing input. It is handed to InAllowlist in the lexical
	// spelling on purpose: that is the spelling the judge used to see before
	// ticket 252's alignment existed, it is the spelling bridge.go:1145 and
	// task.go:838 would carry if C26 had not refused the reparse traversal
	// (risk.ErrReparseDenied), and - which is the reason W2 exists - it is the
	// only spelling under which the two legs can be told apart at all.
	via := filepath.Join(esc, "marker.txt")
	leg1, leg2, rf, ok := legsOf252r2(pc, via)
	got := pc.InAllowlist(via)
	t.Logf("allowed root        : %q -> roots=%q", proj, pc.Roots())
	t.Logf("escaping ask        : %q", via)
	t.Logf("LEG1 lexical        : %v", leg1)
	t.Logf("resolvedForm        : %q ok=%v", rf, ok)
	t.Logf("LEG2 resolved       : %v (refused because ok=false, or because the answer is outside the root)", leg2)
	t.Logf("InAllowlist         : %v", got)
	t.Logf("link target written : %q", outside)
	t.Logf("what the link gives : %q", landed)

	if !leg1 {
		t.Fatalf("carrier moved: the lexical leg already refuses %q (roots=%q), so this case no longer isolates "+
			"the resolvedForm leg - W2 measures nothing and must not be counted as a control", via, pc.Roots())
	}
	if leg2 {
		t.Fatalf("carrier moved: the resolvedForm leg answers inside the root for %q (rf=%q, ok=%v), so this "+
			"platform's resolver neither refuses it nor routes it out of the root, and the second leg has "+
			"nothing to refuse here", via, rf, ok)
	}
	if !got {
		t.Errorf("TICKET 252 BAN RED (resolvedForm leg missing / two containments merged into one): "+
			"InAllowlist(%q) = true while the lexical leg is the only one that vouched for it (leg1=%v leg2=%v, "+
			"resolvedForm=%q ok=%v, roots=%q): that path leaves the allowed root %q through the link %q, which "+
			"the OS reports as %q. Either half of the second leg - \"it will not answer\" or \"its answer is "+
			"elsewhere\" - is the refusal, and both are the same containment in paths.go. Ticket 107's hole "+
			"re-opened: the pair looks redundant on aligned input by design, and AC#2's "+
			"\"绝不许把两次包含改成一次\" is what keeps it that way.",
			via, leg1, leg2, rf, ok, pc.Roots(), proj, esc, outside)
	}
}

// ---------------------------------------------------------------------------
// S1: the call-graph witness. Two containments, from two different sources,
// each able to refuse on its own.
// ---------------------------------------------------------------------------

// bookContainment252r2 is one containment decision over the authorization book
// found in InAllowlist's own body.
type bookContainment252r2 struct {
	file    string
	line    int
	keyText string
	calls   []string // every call the key passed through, in the body's own words
	tainted bool     // the key's chain touches a re-resolution result
	refuses bool     // the decision can send InAllowlist to `return false`
}

// allowlistJudge252r2 locates InAllowlist wherever it lives in this package's
// non-test sources. It is looked up by the pair (receiver type, method name),
// because that pair IS the frozen C19 seam (risk.PathCanonicalizer, declared at
// internal/risk/assessor.go:141) - renaming it is a contract change needing
// human approval, so a red here is the right response, not a false alarm.
// Everything else is matched structurally, never by helper name: a renamed
// rootsContain or resolvedForm keeps this ruler green.
func allowlistJudge252r2(t *testing.T) (fset *token.FileSet, decl *ast.FuncDecl) {
	t.Helper()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("os.ReadDir(.): %v", err)
	}
	fset = token.NewFileSet()
	var found []*ast.FuncDecl
	var files int
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		file, perr := parser.ParseFile(fset, name, nil, parser.ParseComments)
		if perr != nil {
			t.Fatalf("parse %s: %v", name, perr)
		}
		files++
		for _, d := range file.Decls {
			fd, ok := d.(*ast.FuncDecl)
			if !ok || fd.Name.Name != "InAllowlist" || fd.Recv == nil || len(fd.Recv.List) != 1 {
				continue
			}
			if typeName252r2(fd.Recv.List[0].Type) == "PathCanonicalizer" {
				found = append(found, fd)
			}
		}
	}
	if len(found) == 0 {
		t.Fatalf("S1 has nothing to measure: no method named InAllowlist on a PathCanonicalizer receiver "+
			"in any of the %d non-test .go files of this package. If the seam moved packages, this ruler must "+
			"be moved with it by hand - it may not be quietly deleted.", files)
	}
	if len(found) > 1 {
		t.Fatalf("S1 cannot judge two implementations as one: %d methods named InAllowlist on a "+
			"PathCanonicalizer receiver in this package", len(found))
	}
	t.Logf("S1 scanned %d non-test .go files of this package; InAllowlist is declared at %s:%d",
		files, filepath.Base(fset.Position(found[0].Pos()).Filename), fset.Position(found[0].Pos()).Line)
	return fset, found[0]
}

// typeName252r2 unwraps a receiver type (star, parens) down to its identifier.
func typeName252r2(e ast.Expr) string {
	for {
		switch x := e.(type) {
		case *ast.StarExpr:
			e = x.X
		case *ast.ParenExpr:
			e = x.X
		case *ast.Ident:
			return x.Name
		default:
			return ""
		}
	}
}

// TestTicket252R2BothContainmentsReachable is S1: the ban read off the call
// graph instead of a return value.
//
// Four assertions, in the order they fire:
//
//	a. the body asks the authorization book EXACTLY TWICE (the workspace leg at
//	   paths.go:215 asks p.workspace, a different book, so it is not counted and
//	   cannot pad the number);
//	b. each of the two can refuse on its own - its result reaches a
//	   `return false` - so neither is a value computed and thrown away;
//	c. the two keys do not come from the same chain of calls;
//	d. exactly ONE of them is clean, i.e. derived from the parameter without
//	   passing through a re-resolution result (a name bound by a multi-value
//	   assignment, which is how this body spells "the resolver answered, and
//	   here is whether it could"). c and d together are what catch
//	   "two calls, one check": both containments still there, both asking the
//	   already-resolved string, so every input - including W1's and W2's -
//	   answers exactly as it does today.
//
// Brittleness, stated rather than hidden (measured hit surface is in §3/§6 of
// .scratch/wisp/probes/252/r2/impl.md):
//   - it reads one function body. Moving the two containments into a helper and
//     calling that helper twice still counts as two ONLY if the calls are
//     spelled against the book in this body; pushing them into a helper called
//     once is a red, and that red is a false alarm on a legal refactor.
//   - it identifies the book by the receiver field named `roots` (or a local
//     bound to exactly that selector). Renaming the field is a red too.
//   - it does NOT hardcode rootsContain, foldPath or resolvedForm, so those
//     three may be renamed freely.
//   - the red line always names which of a-d fired and what to look at, so a
//     false alarm costs one read, not one archaeology dig.
func TestTicket252R2BothContainmentsReachable(t *testing.T) {
	fset, judge := allowlistJudge252r2(t)
	recv := judge.Recv.List[0].Names[0].Name
	params := judge.Type.Params.List
	if len(params) != 1 || len(params[0].Names) != 1 {
		t.Fatalf("S1's shape assumption broke: InAllowlist no longer takes exactly one named parameter "+
			"(params=%d). The 'derived from the parameter alone' half of assertion d needs the ruler, not the "+
			"production code, to be re-read.", len(params))
	}
	param := params[0].Names[0].Name

	book, tainted, guards := survey252r2(fset, judge, recv)

	// a. exactly two containments over the authorization book.
	if len(book) == 2 {
		lines := make([]string, 0, len(book))
		for _, b := range book {
			lines = append(lines, b.file+":"+strconv.Itoa(b.line)+" key="+b.keyText)
		}
		t.Errorf("TICKET 252 BAN RED (assertion a: two containments, not one): InAllowlist's body asks the "+
			"authorization book %d times, want 2 (found: %v). AC#2's \"绝不许把两次包含改成一次\" is a hard ban, "+
			"and the shape this fires on is exactly the one no value test can see: after 2b1a3071 every string "+
			"production hands this function is already one form, so the two legs agree on all of them and a "+
			"missing leg is invisible from outside. If this deletion was NOT intended, restore both legs at "+
			"paths.go:196 and paths.go:209. If the code is right and the ruler is stale, the ruler is the thing "+
			"that needs a human approval, not a green test run.", len(book), lines)
	}
	// b. each can refuse on its own.
	for _, b := range book {
		if !b.refuses {
			t.Errorf("TICKET 252 BAN RED (assertion b: each containment must be able to refuse): the "+
				"containment at %s:%d (key=%s) has no path to `return false` from its own result, so it computes "+
				"an answer nobody acts on. One live leg is the ticket 107 hole.", b.file, b.line, b.keyText)
		}
	}
	// c and d need both legs present to say anything.
	if len(book) == 2 {
		x, y := book[0], book[1]
		if sameCallSet252r2(x.calls, y.calls) {
			t.Errorf("TICKET 252 BAN RED (assertion c: the two containments must not share one source): both "+
				"keys are built by the same chain of calls (%v) - two calls, one check. Measured at anchor "+
				"e16e3037 this shape passes all 257 cases of this package (M-samesource), which is why the count "+
				"alone is not enough.", x.calls)
		}
		clean, dirty := 0, 0
		for _, b := range book {
			if b.tainted {
				dirty++
			} else {
				clean++
			}
		}
		if clean != 1 || dirty != 1 {
			t.Errorf("TICKET 252 BAN RED (assertion d: one lexical leg and one re-resolving leg, exactly): found "+
				"clean=%d tainted=%d (keys %q / %q, tainted-name set %v, guards %v). The pair exists because the "+
				"lexical spelling and the resolved spelling can name different trees; if both legs read the same "+
				"one, only one tree is ever being asked about.", clean, dirty, x.keyText, y.keyText,
				sortedKeys252r2(tainted), sortedKeys252r2(guards))
		}
	}
	// Census, so the next reader can see what the ruler looked at without
	// re-deriving it: the hit surface is a property of this reading.
	t.Logf("S1 carrier: receiver=%q parameter=%q decl=%s:%d", recv, param,
		filepath.Base(fset.Position(judge.Pos()).Filename), fset.Position(judge.Pos()).Line)
	for _, b := range book {
		t.Logf("containment %s:%d key=%s calls=%v tainted=%v refuses=%v",
			b.file, b.line, b.keyText, b.calls, b.tainted, b.refuses)
	}
	t.Logf("S1 re-resolution names: %v", sortedKeys252r2(tainted))
	t.Logf("S1 names an if-with-return-false consults: %v", sortedKeys252r2(guards))
}

// survey252r2 walks InAllowlist's body once and answers, for every containment
// over the authorization book: what its key is, which calls it passed through,
// whether it touches a re-resolution result, and whether its own result can
// drive a `return false`.
func survey252r2(fset *token.FileSet, judge *ast.FuncDecl, recv string) (
	[]bookContainment252r2, map[string]bool, map[string]bool,
) {
	tainted := map[string]bool{}
	guards := map[string]bool{}
	bindings := map[string]ast.Expr{} // single-value assignments, in body order

	// Pass 1: names bound by a multi-value call are the re-resolution results
	// ("the answer, and whether the resolver could answer at all").
	ast.Inspect(judge.Body, func(n ast.Node) bool {
		as, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}
		if len(as.Rhs) != 1 {
			return true
		}
		if _, isCall := as.Rhs[0].(*ast.CallExpr); !isCall {
			return true
		}
		if len(as.Lhs) >= 2 {
			for _, l := range as.Lhs {
				if id, ok := l.(*ast.Ident); ok && id.Name != "_" {
					tainted[id.Name] = true
					// The binding is recorded too, so a key that passes through
					// "the answer" also passes through the call that produced
					// it: that is what tells the two legs' chains apart.
					bindings[id.Name] = as.Rhs[0]
				}
			}
			return true
		}
		if id, ok := as.Lhs[0].(*ast.Ident); ok {
			bindings[id.Name] = as.Rhs[0]
		}
		return true
	})

	// Pass 2: which names do the body's own refusals consult? Only used to make
	// a red line readable, never to decide anything.
	ast.Inspect(judge.Body, func(n ast.Node) bool {
		ifs, ok := n.(*ast.IfStmt)
		if !ok || !returnsFalse252r2(ifs.Body) {
			return true
		}
		for _, id := range identsOf252r2(ifs.Cond) {
			guards[id] = true
		}
		return true
	})

	// Pass 3: the containment decisions over the authorization book.
	var book []bookContainment252r2
	ast.Inspect(judge.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok || len(call.Args) != 2 {
			return true
		}
		if !readsBook252r2(call.Args[0], recv, bindings) {
			return true
		}
		calls, t := chainOf252r2(call.Args[1], bindings, tainted, map[string]bool{}, 0)
		book = append(book, bookContainment252r2{
			file:    filepath.Base(fset.Position(call.Pos()).Filename),
			line:    fset.Position(call.Pos()).Line,
			keyText: types.ExprString(call.Args[1]),
			calls:   calls,
			tainted: t,
			refuses: refusesHere252r2(n, judge.Body),
		})
		return true
	})
	return book, tainted, guards
}

// readsBook252r2: is this expression the authorization book - the receiver's
// roots field, or a local holding exactly that field?
func readsBook252r2(e ast.Expr, recv string, bindings map[string]ast.Expr) bool {
	if sel, ok := e.(*ast.SelectorExpr); ok {
		if id, ok := sel.X.(*ast.Ident); ok {
			return id.Name == recv && sel.Sel.Name == "roots"
		}
		return false
	}
	if id, ok := e.(*ast.Ident); ok && id.Name != recv {
		if rhs, seen := bindings[id.Name]; seen {
			return readsBook252r2(rhs, recv, bindings)
		}
	}
	return false
}

// chainOf252r2 follows a key expression back through the body's own
// assignments and reports the calls it passed through plus whether it ever
// touches a re-resolution result.
func chainOf252r2(e ast.Expr, bindings map[string]ast.Expr, tainted map[string]bool,
	seen map[string]bool, depth int,
) ([]string, bool) {
	if depth > 8 {
		return nil, false
	}
	calls := []string{}
	hit := false
	var walk func(ast.Expr)
	walk = func(x ast.Expr) {
		switch v := x.(type) {
		case *ast.SelectorExpr:
			walk(v.X)
		case *ast.StarExpr:
			walk(v.X)
		case *ast.UnaryExpr:
			walk(v.X)
		case *ast.BinaryExpr:
			walk(v.X)
			walk(v.Y)
		case *ast.ParenExpr:
			walk(v.X)
		case *ast.CallExpr:
			if id := calleeName252r2(v); id != "" {
				calls = append(calls, id)
			}
			for _, a := range v.Args {
				walk(a)
			}
		case *ast.Ident:
			if tainted[v.Name] {
				hit = true
			}
			if v.Name == "_" || seen[v.Name] {
				return
			}
			seen[v.Name] = true
			if rhs, ok := bindings[v.Name]; ok {
				walk(rhs)
			}
		case *ast.CompositeLit:
			for _, el := range v.Elts {
				walk(el)
			}
		case *ast.IndexExpr:
			walk(v.X)
		}
	}
	walk(e)
	return dedupeSorted252r2(calls), hit
}

func calleeName252r2(c *ast.CallExpr) string {
	switch fn := c.Fun.(type) {
	case *ast.Ident:
		return fn.Name
	case *ast.SelectorExpr:
		return fn.Sel.Name
	}
	return ""
}

// refusesHere252r2: does this decision reach a `return false` of its own? The
// enclosing `if` is found by source range, because ast.Inspect gives no parent
// links: start at the tightest one containing the call and step outward until a
// then-block refuses.
func refusesHere252r2(node ast.Node, body *ast.BlockStmt) bool {
	for cur := innermostIf252r2(node, body); cur != nil; cur = enclosingIfAround252r2(cur, body) {
		if returnsFalse252r2(cur.Body) {
			return true
		}
	}
	return false
}

// innermostIf252r2 is the tightest IfStmt whose source range still contains the
// node.
func innermostIf252r2(node ast.Node, body *ast.BlockStmt) *ast.IfStmt {
	var best *ast.IfStmt
	ast.Inspect(body, func(n ast.Node) bool {
		ifs, ok := n.(*ast.IfStmt)
		if !ok || !nodeWithin252r2(ifs, node) {
			return true
		}
		if best == nil || nodeWithin252r2(best, ifs) {
			best = ifs
		}
		return true
	})
	return best
}

// enclosingIfAround252r2 steps one level out: the tightest IfStmt that contains
// the given one without being it.
func enclosingIfAround252r2(ifs *ast.IfStmt, body *ast.BlockStmt) *ast.IfStmt {
	var out *ast.IfStmt
	ast.Inspect(body, func(n ast.Node) bool {
		cand, ok := n.(*ast.IfStmt)
		if !ok || cand == ifs || !nodeWithin252r2(cand, ifs) {
			return true
		}
		if out == nil || nodeWithin252r2(out, cand) {
			out = cand
		}
		return true
	})
	return out
}

// nodeWithin252r2 compares by byte offsets, which is what "the inner node sits
// inside the outer one's source range" actually means.
func nodeWithin252r2(outer, inner ast.Node) bool {
	if outer == nil || inner == nil {
		return false
	}
	return outer.Pos() <= inner.Pos() && inner.End() <= outer.End()
}

func returnsFalse252r2(b *ast.BlockStmt) bool {
	if b == nil {
		return false
	}
	found := false
	ast.Inspect(b, func(n ast.Node) bool {
		r, ok := n.(*ast.ReturnStmt)
		if !ok {
			return true
		}
		if len(r.Results) == 1 {
			if id, ok := r.Results[0].(*ast.Ident); ok && id.Name == "false" {
				found = true
			}
		}
		return true
	})
	return found
}

func identsOf252r2(e ast.Expr) []string {
	var out []string
	ast.Inspect(e, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok {
			out = append(out, id.Name)
		}
		return true
	})
	return out
}

func sameCallSet252r2(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func dedupeSorted252r2(in []string) []string {
	set := map[string]bool{}
	for _, s := range in {
		set[s] = true
	}
	return sortedKeys252r2(set)
}

func sortedKeys252r2(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
