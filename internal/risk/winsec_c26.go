package risk

import "github.com/CarlosShao/wisp/internal/winsec"

// This file is the wiring, not the pipeline: internal/risk owns C26 (SPEC-06 §4)
// and hands it to the one package whose *security decision* is taken from a path
// spelling. internal/winsec cannot import this package - the graph already runs
// risk -> observe -> secret -> winsec (observe/redact.go calls secret.RedactSecret,
// secret/store.go calls winsec.PrivateDirAll), so the reverse edge would close a
// cycle. Installing in this direction keeps C26 a single implementation: winsec
// asks, this package resolves, and no second normalization pipeline is written
// anywhere (which is what D22 ban #2 is actually protecting, not the symbol).
//
// Reparse exceptions are deliberately nil. [fs] reparse_point_exceptions exist so
// a *user-chosen* allowlisted junction can be traversed for reads; nothing about
// that licenses placing the DPAPI blob directory, the SQLite data root or the
// artifact tree behind somebody's link. For placement, default-deny is the only
// policy, and winsec propagates the resulting ErrReparseDenied unchanged so the
// caller can tell "C26 refused" from "this package is unwired" (ErrNoPathResolver).
func init() {
	winsec.SetPathResolver(c26Pipeline{})
}

// c26Pipeline adapts the C26 entry point to the narrow shape winsec needs.
type c26Pipeline struct{}

// Resolve returns the pipeline's canonical spelling: the handle-based real path
// when the object exists (8.3 expanded, reparse targets resolved, \\?\ and UNC
// normalized), the lexical form only when it does not - which is the normal
// state for a directory winsec is about to create, and still after the reparse
// traversal check, so a nonexistent leaf under an existing junction is refused
// rather than created through the link.
//
// This is the one C26 leg that both ACTS on a tree (through winsec, which
// seals whatever it returns) and reports success to its caller, so it goes
// through Result.Actable: an expansion that moved the path off the tree the
// caller named (ticket 102 / PROBE B2) is refused here and propagates through
// winsec.ResolvePath's %w wrapping as a sealing failure - never as "sealed".
func (c26Pipeline) Resolve(input string) (string, error) {
	res, err := Resolve(input, nil)
	if err != nil {
		return "", err
	}
	return res.Actable()
}
