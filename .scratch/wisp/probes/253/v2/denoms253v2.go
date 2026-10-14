// 253-v2 acceptance leg: an INDEPENDENT reader of the five denominators the
// shipped ruler (internal/panel/inbound_roster_253_test.go) claims to reconcile.
// Written from the shipped ruler's documented intent, not by importing it.
// Run: go run denoms253v2.go <repoRoot>
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func dotShaped(s string) bool {
	segs := strings.Split(s, ".")
	if len(segs) < 2 {
		return false
	}
	for _, seg := range segs {
		if seg == "" {
			return false
		}
		for i, r := range seg {
			switch {
			case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_':
			case r == '-' && i > 0:
			default:
				return false
			}
		}
	}
	return true
}

type parsed struct {
	consts map[string]string // ident -> value (package-level string consts)
	lines  map[string]int    // ident -> line
}

func parseProd(path string) (*ast.File, *parsed, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	fset := token.NewFileSet()
	af, err := parser.ParseFile(fset, filepath.Base(path), src, 0)
	if err != nil {
		return nil, nil, err
	}
	p := &parsed{consts: map[string]string{}, lines: map[string]int{}}
	for _, d := range af.Decls {
		gd, ok := d.(*ast.GenDecl)
		if !ok || gd.Tok != token.CONST {
			continue
		}
		for _, sp := range gd.Specs {
			vs, ok := sp.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for i, id := range vs.Names {
				if i >= len(vs.Values) {
					continue
				}
				lit, ok := vs.Values[i].(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					continue
				}
				v, err := strconv.Unquote(lit.Value)
				if err != nil {
					continue
				}
				p.consts[id.Name] = v
				p.lines[id.Name] = fset.Position(lit.Pos()).Line
			}
		}
	}
	return af, p, nil
}

func findFunc(af *ast.File, name string) *ast.FuncDecl {
	var fn *ast.FuncDecl
	for _, d := range af.Decls {
		if fd, ok := d.(*ast.FuncDecl); ok && fd.Name.Name == name {
			fn = fd
		}
	}
	return fn
}

// cases returns resolved case values, raw labels that could not resolve, and
// the called names of any default branch, for the first switch in fn.
func cases(fset *token.FileSet, fn *ast.FuncDecl, maps []*parsed) (vals, unres []string, defaultCalls []string, hasDefault bool) {
	ast.Inspect(fn, func(n ast.Node) bool {
		sw, ok := n.(*ast.SwitchStmt)
		if !ok {
			return true
		}
		for _, st := range sw.Body.List {
			cc, ok := st.(*ast.CaseClause)
			if !ok {
				continue
			}
			if len(cc.List) == 0 {
				hasDefault = true
				for _, s := range cc.Body {
					ast.Inspect(s, func(m ast.Node) bool {
						if c, ok := m.(*ast.CallExpr); ok {
							switch f := c.Fun.(type) {
							case *ast.Ident:
								defaultCalls = append(defaultCalls, f.Name)
							case *ast.SelectorExpr:
								defaultCalls = append(defaultCalls, f.Sel.Name)
							}
						}
						return true
					})
				}
				continue
			}
			for _, e := range cc.List {
				switch l := e.(type) {
				case *ast.BasicLit:
					if l.Kind == token.STRING {
						v, err := strconv.Unquote(l.Value)
						if err == nil {
							vals = append(vals, v)
							continue
						}
					}
					unres = append(unres, l.Value)
				case *ast.Ident:
					found := false
					for _, p := range maps {
						if v, ok := p.consts[l.Name]; ok {
							vals = append(vals, v)
							found = true
							break
						}
					}
					if !found {
						unres = append(unres, l.Name)
					}
				default:
					unres = append(unres, fmt.Sprintf("%T", e))
				}
			}
		}
		return true
	})
	sort.Strings(vals)
	return
}

func main() {
	root := os.Args[1]
	panel := filepath.Join(root, "internal", "panel")
	bridgePath := filepath.Join(panel, "bridge.go")
	dispPath := filepath.Join(panel, "composer_dispatch.go")

	bridgeAF, bridgeP, err := parseProd(bridgePath)
	if err != nil {
		panic(err)
	}
	fmt.Println("== denominator 1 (my read): bridge.go package-level string consts ==")
	var declared []string
	names := make([]string, 0, len(bridgeP.consts))
	for n := range bridgeP.consts {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		v := bridgeP.consts[n]
		hit := strings.HasPrefix(n, "Method") || dotShaped(v)
		fmt.Printf("  %-24s = %-24q MethodPrefix=%-5v dotShaped=%-5v inbound=%v (line %d)\n",
			n, v, strings.HasPrefix(n, "Method"), dotShaped(v), hit, bridgeP.lines[n])
		if hit {
			declared = append(declared, v)
		}
	}
	sort.Strings(declared)
	fmt.Printf("  DECLARED set = %q\n", declared)

	// all package consts (for the router's resolution).
	allFiles, _ := filepath.Glob(filepath.Join(panel, "*.go"))
	var pkgConsts *parsed = &parsed{consts: map[string]string{}, lines: map[string]int{}}
	var prodFiles []string
	literalHits := map[string][]string{}
	for _, f := range allFiles {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		prodFiles = append(prodFiles, filepath.Base(f))
		for k, v := range func() map[string]string { _, p, e := parseProd(f); if e != nil { panic(e) }; return p.consts }() {
			if old, ok := pkgConsts.consts[k]; !ok {
				pkgConsts.consts[k] = v
			} else if old != v {
				pkgConsts.consts[k] = "\x00ambiguous"
			}
		}
		// denominator 5 numerator: every dot-shaped string literal in production files.
		src, _ := os.ReadFile(f)
		fset := token.NewFileSet()
		af, perr := parser.ParseFile(fset, filepath.Base(f), src, 0)
		if perr != nil {
			panic(perr)
		}
		ast.Inspect(af, func(n ast.Node) bool {
			lit, ok := n.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}
			v, uerr := strconv.Unquote(lit.Value)
			if uerr != nil || !dotShaped(v) {
				return true
			}
			literalHits[v] = append(literalHits[v], fmt.Sprintf("%s:%d", filepath.Base(f), fset.Position(lit.Pos()).Line))
			return true
		})
	}
	sort.Strings(prodFiles)
	fmt.Printf("\n== package-level string consts merged from %d production files: %d names ==\n", len(prodFiles), len(pkgConsts.consts))

	fmt.Println("\n== denominator 3 (my read): knownComposerMethod case labels (bridge.go) ==")
	guard := findFunc(bridgeAF, "knownComposerMethod")
	fset := token.NewFileSet()
	// re-parse for position-based walk consistency: just use fn from bridgeAF
	pass1 := []*parsed{bridgeP}
	gv, gu, gd, gHas := cases(fset, guard, pass1)
	fmt.Printf("  labels resolved via bridge.go consts = %q\n", gv)
	fmt.Printf("  unresolvable labels = %q ; default present = %v ; default calls = %q\n", gu, gHas, gd)

	fmt.Println("\n== denominator 4 (my read): composer_dispatch.go dispatch case labels ==")
	dispAF, _, err := parseProd(dispPath)
	if err != nil {
		panic(err)
	}
	disp := findFunc(dispAF, "dispatch")
	dv, du, dd, dHas := cases(fset, disp, []*parsed{bridgeP, pkgConsts})
	fmt.Printf("  labels resolved via bridge.go+pkg consts = %q\n", dv)
	fmt.Printf("  unresolvable labels = %q ; default present = %v ; default calls = %q\n", du, dHas, dd)

	fmt.Println("\n== denominator 5 (my read): dot-shaped literals in production files, and whether the DECLARED six answer them ==")
	keys := make([]string, 0, len(literalHits))
	for k := range literalHits {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	six := map[string]bool{"panel.mode.request": true, "panel.workspace.request": true, "panel.attachment.add": true, "panel.message.send": true, "config.get": true, "config.set": true}
	for _, k := range keys {
		fmt.Printf("  %-24q answeredByTheSix=%-5v sites=%v\n", k, six[k], literalHits[k])
	}
}
