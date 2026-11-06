package main

// Ban #10 machinery (ticket 350: the bounded wordlist gate).
//
// WHAT THIS FILE HOLDS: the pure halves of the gate - the word shape, the
// declaration-roster collector, and the comment citation scanner. The WALK that
// feeds them, the counter it bumps, the emission site (the literal s.add tag
// selftest.go's readRosters() reads out of main.go) and the ledger row all live
// in main.go, because every other emission site does and readRosters() reads
// exactly one source file.
//
// THE BOUNDS ARE THE APPROVAL (ticket 350 AC#0, owner 2026-11-06, ledger A1571;
// the shape was measured and defended by 343-r1 / 343-v1, verdict docs
// docs/evidence/s1/343-v1-verdict.md):
//
//   - the WORD SHAPE is `a[A-Z][A-Za-z0-9_]*` at a word boundary, and nothing
//     wider. The general forms were measured and REJECTED before this gate was
//     approved: loose call shapes ring 142 times on prose ("the candidate (",
//     "side ("), and any-camelCase rings ~50% false positives (hostPathBox /
//     assessorFor are real symbols of OTHER packages, macOS is platform prose).
//     Generalizing this regex re-opens ticket 350 (AC#2) and is not a tidy-up.
//   - the TREE is internal/risk/ production Go files (non-_test.go, non-
//     testdata, git-ignored paths excluded like every other walk). _test.go is
//     outside the range by the same default as bans #1-5 - measured, not
//     assumed: at HEAD every a[A-Z] token the test files' comments name (aEq,
//     aRuleReason, aUnder) resolves to a production declaration, so the
//     exclusion hides nothing today (recorded in
//     docs/evidence/s1/350-bounded-wordlist-gate.md).
//   - the ROSTER is the same package's own declarations, collected by go/ast
//     (collectDeclaredIdents below). It is deliberately EVERY identifier the
//     package declares - package-level funcs/types/consts/vars, methods,
//     struct and interface fields, parameters and named results, := / range /
//     var locals, labels - and not only package-level symbols. Two reasons:
//     every `a[A-Z]...` name is necessarily UNEXPORTED (lowercase first
//     letter), so it cannot legally refer to another package's symbol and the
//     same-package roster is the whole truth; and a comment citing a local or
//     a method must not be convicted for the citation being real. The measured
//     false-positive surface at HEAD is 0 and the five false-positive CLASSES
//     the two census legs enumerated for the wider shapes (other-package real
//     symbol, platform prose, marker convention, vendor name, deliberately
//     deleted historical name) are recorded in the evidence doc - none of them
//     is implementable as a free-text exemption here, by AC#1's own hard
//     criterion: an allowlist would be exactly the degeneration the bounded
//     shape exists to avoid.

import (
	"go/ast"
	"go/token"
	"regexp"
)

// aWordRe is ban #10's bounded word shape. The \b is load-bearing on the left
// (macOS must not become aCOS-shaped nonsense: the 'a' of "macOS" is preceded
// by a word character, so \b refuses it) and the trailing [A-Za-z0-9_]* is
// greedy, so a token is read as the maximal identifier: a comment naming
// aEquivFoo convicts aEquivFoo, never a prefix of it.
//
// Why comments and only comments: the ban exists because a comment names a
// symbol as a navigation beacon for the next leg (the same reason ban #9
// exists), and `blacklist.go`'s aEquiv hung for a month inside one. Identifiers
// in CODE are checked by the compiler; a name in a comment is checked by
// nobody but this gate.
var aWordRe = regexp.MustCompile(`\ba[A-Z][A-Za-z0-9_]*`)

// collectDeclaredIdents records every identifier the parsed file DECLARES into
// the roster set. Declaration positions only: a name that appears merely as a
// USE must not buy a comment citation a clean bill, because a use without a
// declaration does not compile - if the gate ever needs to be fooled by a use,
// the tree is already red for a better reason.
//
// The switch enumerates go/ast's declaration carriers explicitly. ImportSpec
// is deliberately absent: an import name is not a symbol of this package.
func collectDeclaredIdents(f *ast.File, into map[string]bool) {
	ast.Inspect(f, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.FuncDecl:
			// Plain funcs and methods alike: a method name is exactly as
			// citable in a comment as a function name, and both are unexported
			// under the word shape this gate watches.
			into[v.Name.Name] = true
		case *ast.TypeSpec:
			into[v.Name.Name] = true
		case *ast.ValueSpec:
			// const and var at every level, including the single-name form
			// `var aThing = ...` inside a function body.
			for _, name := range v.Names {
				into[name.Name] = true
			}
		case *ast.Field:
			// struct fields, interface methods, function parameters and
			// named results - including the ones inside func literals.
			for _, name := range v.Names {
				into[name.Name] = true
			}
		case *ast.AssignStmt:
			if v.Tok == token.DEFINE {
				for _, lhs := range v.Lhs {
					if id, ok := lhs.(*ast.Ident); ok {
						into[id.Name] = true
					}
				}
			}
		case *ast.RangeStmt:
			if v.Tok == token.DEFINE {
				if id, ok := v.Key.(*ast.Ident); ok {
					into[id.Name] = true
				}
				if id, ok := v.Value.(*ast.Ident); ok {
					into[id.Name] = true
				}
			}
		case *ast.LabeledStmt:
			into[v.Label.Name] = true
		}
		return true
	})
}

// riskCitation is one a[A-Z]... token a comment names, with where it sits. The
// walk in main.go turns these into findings (or silence, for a token the
// roster holds).
type riskCitation struct {
	Token string
	Path  string
	Line  int
}

// riskCommentCitations scans one parsed file's comments for the bounded word
// shape and returns every hit, whether or not the roster holds it - the caller
// owns the roster check, so tests can pin the census (every token a comment
// names) separately from the verdict (which of them are phantom).
//
// Comments come from f.Comments, which only a parser.ParseComments run
// populates - the same guarantee ban #9 leans on: the token is inside a
// comment BY CONSTRUCTION, not by a line regex that also eats string
// literals. One citation per (file, token) pair is the caller's dedupe.
func riskCommentCitations(fset *token.FileSet, f *ast.File, path string) []riskCitation {
	var out []riskCitation
	for _, cg := range f.Comments {
		for _, c := range cg.List {
			pos := fset.Position(c.Pos())
			for _, loc := range aWordRe.FindAllStringIndex(c.Text, -1) {
				out = append(out, riskCitation{
					Token: c.Text[loc[0]:loc[1]],
					Path:  path,
					Line:  pos.Line,
				})
			}
		}
	}
	return out
}
