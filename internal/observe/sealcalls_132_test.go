package observe

// Ticket 132's AC#4 enumeration door: the one that stops "the capability exists
// and nothing calls it" from coming back.
//
// The defect this ticket was opened on is not a missing function. winsec.SealDir
// was written, tested sixteen times, and called by nothing in the product, so the
// log tree shipped world-readable while every test about sealing was green. A
// behavioural case cannot catch that on its own: it asserts what happened on the
// objects it made, and a future leg that stops sealing produces objects nobody
// asserts on. So this door reads the CALL GRAPH instead, and states the claim the
// ticket makes: the sealing entry points are CALLED from outside the package that
// defines them, and the log tree's own package is one of the callers.
//
// Read through go/parser, not a line regex, because a regex counts prose. This
// repository writes long comments that name winsec.SealFile by hand (the header of
// logging.go is one), so a grep-shaped door would stay green after the call was
// deleted - checklist item 3 of the ruler-scope list, "am I counting code or
// comments". An AST selector is only ever code.
//
// A CALL, not a mention, and that distinction was measured rather than assumed:
// the first version of this door counted any selector, and it stayed GREEN with
// the file seal deleted, because the seam declaration `var sealLogFile =
// winsec.SealFile` still mentioned the name while nothing called it (mutation M2 in
// docs/evidence/s1/132-r1-impl.md). That is the exact shape this door is for, so an
// alias is now followed: a package-level `var X = winsec.SealFile` counts only when
// something in the same package calls X, and an uncalled alias is a red with its
// own name in the message.
//
// What it does NOT prove, stated because an unbounded claim is the thing to
// distrust: it proves a call exists, not that the call is on a path the product
// runs. The behavioural cases next to this file are what say the objects are
// sealed; this door is what says the sealing is still wired to a caller. Nor does
// it follow an alias carried in a struct field, through an interface, or across
// packages - the three shapes this repository's own leg-sink gate
// (cmd/wisp/leg_sink_gate_131_test.go) had to grow readers for, and the same
// blindness is stated there.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
)

// sealEntry132 is one reference to a winsec sealing entry point found in
// production code.
type sealEntry132 struct {
	dir   string // package directory, slashed, relative to the module root
	name  string // the sealing entry point's own name
	call  bool   // the selector is the callee, directly or through a package-level alias
	alias string // the alias the call went through, when it was not a direct call
}

// sealers132 is the family this door guards: the directory seals and the file
// seal. PrivateDirAll and SealDir are one family because they answer one question
// (make this directory private); a leg may use either.
var sealers132 = map[string]bool{"PrivateDirAll": true, "SealDir": true, "SealFile": true}

func dirSealer132(name string) bool { return name == "PrivateDirAll" || name == "SealDir" }

// TestSealEntryPointsHaveProductionCallers132r1 is the door. Its roster is printed
// on every run, green included, because the reading worth keeping is "who calls
// the sealer today" - a red should be the second thing a reviewer sees, not the
// first.
func TestSealEntryPointsHaveProductionCallers132r1(t *testing.T) {
	moduleRoot := moduleRoot132(t)

	direct, aliases, scanned := scanSealers132(t, moduleRoot)
	if scanned < 200 {
		t.Fatalf("the scan read only %d of the ~290 production files under internal/ and cmd/; the layout this door enumerates has moved, so a green here would be a green over nothing", scanned)
	}
	entries := append(direct, aliasCalls132(t, moduleRoot, aliases)...)

	var roster []string
	for _, e := range entries {
		shape := "value (nothing calls it)"
		switch {
		case e.alias != "":
			shape = "call through alias " + e.alias
		case e.call:
			shape = "call"
		}
		roster = append(roster, e.dir+" -> winsec."+e.name+" ("+shape+")")
	}
	sort.Strings(roster)
	t.Logf("ticket132 door: %d sealing references outside internal/winsec, from %d production files:\n  %s",
		len(entries), scanned, strings.Join(roster, "\n  "))

	// The log tree's package must CALL both halves: a directory seal and a file
	// seal. Asserted per capability rather than per function name so replacing
	// PrivateDirAll with an equivalent walk is not a red, while removing the seal
	// is.
	var haveDirSeal, haveFileSeal bool
	for _, e := range entries {
		if e.dir != "internal/observe" || !e.call {
			continue
		}
		if dirSealer132(e.name) {
			haveDirSeal = true
		}
		if e.name == "SealFile" {
			haveFileSeal = true
		}
	}
	if !haveDirSeal {
		t.Errorf("internal/observe does not CALL winsec.PrivateDirAll or winsec.SealDir: the log directory is created unsealed again, which is the defect ticket 132 landed here. The roster this door did find is logged above.")
	}
	if !haveFileSeal {
		t.Errorf("internal/observe does not CALL winsec.SealFile: a rolled log file can be created and written without ever being sealed, which is the defect ticket 132 landed here. The roster this door did find is logged above.")
	}

	// And the family-wide claim: no caller at all outside internal/winsec would
	// mean the API is a shelf again, whatever this package happens to do.
	called := 0
	for _, e := range entries {
		if e.call {
			called++
		}
	}
	if called == 0 {
		t.Error("no package outside internal/winsec calls a winsec sealing entry point: the sealer is unreachable from the product")
	}
}

func moduleRoot132(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))
}

// scanSealers132 walks internal/ and cmd/ for production .go files and returns
// every direct reference to a sealing entry point, plus the package-level aliases
// each package declares for one (alias name -> entry point name).
//
// Build-tagged files are read regardless of GOOS on purpose: the door is about
// wiring that has to exist on the shipped platform, and a reading that changed
// with the host would report a Windows-only caller as missing on the Linux CI.
func scanSealers132(t *testing.T, moduleRoot string) ([]sealEntry132, map[string]map[string]string, int) {
	t.Helper()
	var direct []sealEntry132
	aliases := map[string]map[string]string{}
	scanned := 0
	fset := token.NewFileSet()
	for _, root := range []string{"internal", "cmd"} {
		base := filepath.Join(moduleRoot, root)
		err := filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if d.Name() == "testdata" {
					return filepath.SkipDir
				}
				return nil
			}
			if skipFile132(path) {
				return nil
			}
			file, perr := parser.ParseFile(fset, path, nil, 0)
			if perr != nil {
				return perr
			}
			scanned++
			dir, rerr := filepath.Rel(moduleRoot, filepath.Dir(path))
			if rerr != nil {
				return rerr
			}
			dir = filepath.ToSlash(dir)
			callee := calleePositions132(file)
			ast.Inspect(file, func(n ast.Node) bool {
				if spec, isSpec := n.(*ast.ValueSpec); isSpec {
					for i, name := range spec.Names {
						if i >= len(spec.Values) {
							continue
						}
						if entry, ok := winsecSelector132(spec.Values[i]); ok {
							if aliases[dir] == nil {
								aliases[dir] = map[string]string{}
							}
							aliases[dir][name.Name] = entry
						}
					}
				}
				sel, isSel := n.(*ast.SelectorExpr)
				if !isSel {
					return true
				}
				name, ok := winsecSelector132(sel)
				if !ok {
					return true
				}
				direct = append(direct, sealEntry132{dir: dir, name: name, call: callee[sel.Pos()]})
				return true
			})
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", base, err)
		}
	}
	return dropWinsecSelf132(direct), aliases, scanned
}

// aliasCalls132 finds calls to the aliases scanSealers132 recorded. A package's
// aliases are collected from all of its files first, so an alias declared in one
// file and called in another is still seen as a call.
func aliasCalls132(t *testing.T, moduleRoot string, aliases map[string]map[string]string) []sealEntry132 {
	t.Helper()
	var out []sealEntry132
	fset := token.NewFileSet()
	for _, root := range []string{"internal", "cmd"} {
		base := filepath.Join(moduleRoot, root)
		err := filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if d.Name() == "testdata" {
					return filepath.SkipDir
				}
				return nil
			}
			if skipFile132(path) {
				return nil
			}
			file, perr := parser.ParseFile(fset, path, nil, 0)
			if perr != nil {
				return perr
			}
			dir, rerr := filepath.Rel(moduleRoot, filepath.Dir(path))
			if rerr != nil {
				return rerr
			}
			dir = filepath.ToSlash(dir)
			names := aliases[dir]
			if len(names) == 0 {
				return nil
			}
			ast.Inspect(file, func(n ast.Node) bool {
				call, isCall := n.(*ast.CallExpr)
				if !isCall {
					return true
				}
				ident, isIdent := call.Fun.(*ast.Ident)
				if !isIdent {
					return true
				}
				if entry, ok := names[ident.Name]; ok {
					out = append(out, sealEntry132{dir: dir, name: entry, call: true, alias: ident.Name})
				}
				return true
			})
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", base, err)
		}
	}
	return dropWinsecSelf132(out)
}

// calleePositions132 is the set of selector positions used as a call's function,
// so a reference passed on as a value is not read as a call.
func calleePositions132(file *ast.File) map[token.Pos]bool {
	callee := map[token.Pos]bool{}
	ast.Inspect(file, func(n ast.Node) bool {
		if ce, isCall := n.(*ast.CallExpr); isCall {
			if sel, isSel := ce.Fun.(*ast.SelectorExpr); isSel {
				callee[sel.Pos()] = true
			}
		}
		return true
	})
	return callee
}

// winsecSelector132 answers the entry point name when expr is winsec.<sealer>.
func winsecSelector132(expr ast.Expr) (string, bool) {
	sel, isSel := expr.(*ast.SelectorExpr)
	if !isSel || !sealers132[sel.Sel.Name] {
		return "", false
	}
	x, isIdent := sel.X.(*ast.Ident)
	if !isIdent || x.Name != "winsec" {
		return "", false
	}
	return sel.Sel.Name, true
}

func skipFile132(path string) bool {
	return !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go")
}

// dropWinsecSelf132 removes references from internal/winsec: that is where the
// definitions live, and a package talking to itself says nothing about the product.
func dropWinsecSelf132(in []sealEntry132) []sealEntry132 {
	kept := in[:0]
	for _, e := range in {
		if e.dir == "internal/winsec" {
			continue
		}
		kept = append(kept, e)
	}
	return kept
}
