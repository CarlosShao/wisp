// Package main is a read-only probe for ticket 160 (dispatch 160-c1).
//
// It replicates ONLY the two candidate shapes for "OpenScope hands back a
// closer" over a miniature DisposalScope, so that the compile-time claims in
// the evidence file are measurements instead of assertions:
//
//	甲  the close action lives as a local closure that never leaves the
//	    package that created it (openA registers it and returns nothing).
//	乙  OpenScope hands back a struct carrying Close() (openB).
//
// Run output is the reading; the interesting part is what COMPILES at all.
package main

import "fmt"

// scopeSet is the miniature of risk.Provenance's p.scopes map.
type scopeSet map[string][]string

// disposal is the miniature of plugin.DisposalScope: Defer/Dispose/Incomplete.
type disposal struct {
	fns      []func()
	disposed bool
	late     bool
}

func (d *disposal) Defer(fn func()) bool {
	if d.disposed {
		d.late = true
		return false
	}
	d.fns = append(d.fns, fn)
	return true
}

func (d *disposal) Dispose() {
	d.disposed = true
	for i := len(d.fns) - 1; i >= 0; i-- {
		d.fns[i]()
	}
	d.fns = nil
}

func (d *disposal) Incomplete() bool { return d.late }

// Engine is the shared state both shapes hang on.
type Engine struct {
	scopes scopeSet
}

func NewEngine() *Engine { return &Engine{scopes: scopeSet{}} }

// OpenA is shape 甲: the closer is a local value; open registers it on the
// task's disposal scope and hands the caller NOTHING.
func (e *Engine) OpenA(scopeID string, d *disposal) {
	e.scopes[scopeID] = nil
	d.Defer(func() { e.closeA(scopeID) })
}

// closeA is unexported: with OpenA, no caller outside this package can close
// anything at all, so "only the issuer can close" is free here — and so is
// "nobody outside can close", including the legitimate owner.
func (e *Engine) closeA(scopeID string) { delete(e.scopes, scopeID) }

// ScopeHandle is shape 乙.
type ScopeHandle struct {
	e      *Engine
	id     string
	closed bool
}

// OpenB is shape 乙: open hands back the handle that carries its own closer.
// The return value is DISCARDABLE — Go imposes no use-on-result rule.
func (e *Engine) OpenB(scopeID string, d *disposal, token int) *ScopeHandle {
	e.scopes[scopeID] = nil
	h := &ScopeHandle{e: e, id: scopeID}
	d.Defer(h.Close)
	return h
}

// Close is the identity-checked closer: the id is a private field of the
// handle that was issued, so no caller names an arbitrary scope.
func (h *ScopeHandle) Close() {
	if h == nil || h.closed {
		return
	}
	h.closed = true
	delete(h.e.scopes, h.id)
}

// CloseAny is the MUTANT: AC#3's "drop the identity clause" build. The closer
// is again reachable with any scope id from any package.
func (e *Engine) CloseAny(scopeID string) { delete(e.scopes, scopeID) }

// Taints is the externally visible ledger the census reads.
func (e *Engine) Taints(scopeID string) []string { return e.scopes[scopeID] }

// OpenCount is the externally visible reading that answers "did a scope leak".
func (e *Engine) OpenCount() int { return len(e.scopes) }

func main() {
	// LEG 1: shape 甲, opened and never explicitly closed by the caller.
	d1 := &disposal{}
	e1 := NewEngine()
	e1.OpenA("task-1", d1)
	fmt.Printf("A-after-open  scopes=%d taints=%v\n", e1.OpenCount(), e1.Taints("task-1"))
	d1.Dispose()
	fmt.Printf("A-after-dispose scopes=%d\n", e1.OpenCount())

	// LEG 2: shape 甲 with NO disposal scope wired at all — this is
	// cmd/wisp/panel_assets.go's leg (a CLI probe that owns no task scope).
	e2 := NewEngine()
	e2.OpenA("panel-assets-l2", &disposal{}) // handle of the closer is dropped here, and it still compiles
	fmt.Printf("A-no-dispose-owner scopes=%d\n", e2.OpenCount())

	// LEG 3: shape 乙, the return value discarded outright.
	e3 := NewEngine()
	e3.OpenB("panel-assets-l2", &disposal{}, 1) // discarded: compiles clean
	fmt.Printf("B-discarded scopes=%d\n", e3.OpenCount())

	// LEG 4: AC#3's mutation — identity dropped, so a foreign closer hits the
	// victim's id. The reading that matters is which of these two lines moves.
	e4 := NewEngine()
	e4.OpenB("task-victim", &disposal{}, 1)
	e4.CloseAny("task-victim") // mis-close from an unrelated owner
	fmt.Printf("B-no-identity mis-closed scopes=%d (0 = victim silently lost its scope)\n", e4.OpenCount())

	// LEG 5: the same mis-close with the identity clause IN. There is no
	// expression that names the victim's id at all.
	e5 := NewEngine()
	h5 := e5.OpenB("task-victim", &disposal{}, 1)
	h5.Close()
	fmt.Printf("B-with-identity own-close scopes=%d\n", e5.OpenCount())
}
