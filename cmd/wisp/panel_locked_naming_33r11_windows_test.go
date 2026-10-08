//go:build windows

package main

// Ticket 33, item (c) of the orchestrator's ruling in
// .scratch/wisp/issues/33-panel-host-c27.md (leg 33-r11, 2026-10-08): a ruler over
// this package's OWN SOURCE TEXT that pins "the Locked suffix and the locking
// behaviour agree".
//
// What the suffix means in this repo: setPriorFocusLocked in
// cmd/wisp/panel_host_windows.go documents "The caller must hold m.mu" and its one
// call site sits inside a locked section of Show. A method carrying that suffix
// which takes the mutex itself is the opposite shape, and sync.Mutex is not
// reentrant: the next person who trusts the name and calls it while holding the lock
// parks that thread on the spot. That is why the two methods this leg renamed
// (firstRoundTrip, serveNotBuiltNotice) no longer carry the suffix - they read m.w
// under m.mu themselves, so they must not be called with the lock held.
//
// What this file pins, in two independent ways:
//
//	(1) census: every function or method in this package's production sources whose
//	    name ends in "Locked" must (a) be named in the roster below and (b) take no
//	    lock anywhere in its own body. A new Locked-suffixed func is red until it is
//	    reviewed and named here, and a roster name that no longer exists is red too,
//	    so neither list can rot quietly.
//	(2) behaviour: the rostered method that is reachable without a window is called
//	    for real while the caller holds m.mu, inside a bounded wait, and must return.
//	    The bounded-wait shape is the one internal/config/manager_223_test.go already
//	    uses for deadlock probes (completedWithin there).
//
// Why a source census and not only calls: asking "does this Locked-suffixed method
// take the lock itself?" of a method we do not yet have means NOT calling it - a yes
// answer would hang the very run that asks. The census is the part with general
// reach; the behavioural case is the part proving the census describes code that
// really runs.
//
// Positive controls live in this file, parsed from source strings this file owns and
// never planted in production: they show the census bites on a bad name, stays off an
// honest one, ignores the honest-no-suffix shape (which is ticket 33's item (b)
// family, out of this ruler's range), and that the bounded wait reports a parked
// call. Without them this file would be a ruler that could pass by reading nothing.
//
// A bare `go` appears in this file on purpose (completedWithin33r11). tools/d22scan's
// bans #1-5 skip *_test.go by design (walkGo), and this case probes a deadlock, so it
// must start a goroutine without joining an ownership registry that would itself be
// part of what is under test - the same reason internal/config/manager_223_test.go
// gives. Production gains no goroutine here.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/sys/windows"
)

// lockedSuffix33r11 is the naming convention under guard.
const lockedSuffix33r11 = "Locked"

// lockedNamingRoster33r11 is the named exemption list: every Locked-suffixed
// function this package is allowed to carry, reviewed one by one as honestly
// requiring the caller to hold the mutex. Adding a name here is the reviewed act;
// leaving one behind when the func is gone is caught by the same case.
//
//	PanelManager.setPriorFocusLocked: its own doc comment says "The caller must hold
//	m.mu", its body takes no lock, and its single call site runs inside Show's
//	m.mu.Lock() / m.mu.Unlock() section.
//
// The two methods this leg renamed are deliberately ABSENT: they take m.mu inside
// their own bodies, so they must not carry the Locked suffix and no roster entry
// could make that honest.
var lockedNamingRoster33r11 = map[string]string{
	"PanelManager.setPriorFocusLocked": "requires the caller to hold m.mu; body takes no lock",
}

// lockMethodNames33r11 are the mutex and rwmutex operations whose presence in a
// body means that body touches a lock itself. runtime.LockOSThread is deliberately
// not in the set: it is not this convention's mutex.
var lockMethodNames33r11 = map[string]bool{
	"Lock": true, "Unlock": true, "RLock": true, "RUnlock": true,
	"TryLock": true, "TryRLock": true,
}

// lockedFunc33r11 is one census hit.
type lockedFunc33r11 struct {
	qualified string
	file      string
	line      int
	lockCalls []string
}

// lockedFuncsIn33r11 walks one parsed file and returns every Locked-suffixed
// function it declares, with the lock operations found anywhere in its body
// (including inside a closure it hands off, which is a worse shape, not a milder one).
func lockedFuncsIn33r11(fset *token.FileSet, file *ast.File, srcName string) []lockedFunc33r11 {
	var hits []lockedFunc33r11
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || !strings.HasSuffix(fn.Name.Name, lockedSuffix33r11) {
			continue
		}
		hit := lockedFunc33r11{
			qualified: fn.Name.Name,
			file:      srcName,
			line:      fset.Position(fn.Pos()).Line,
		}
		if fn.Recv != nil && len(fn.Recv.List) == 1 {
			hit.qualified = recvName33r11(fn.Recv.List[0].Type) + "." + fn.Name.Name
		}
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if lockMethodNames33r11[sel.Sel.Name] {
				hit.lockCalls = append(hit.lockCalls, exprText33r11(sel.X)+"."+sel.Sel.Name)
			}
			return true
		})
		hits = append(hits, hit)
	}
	return hits
}

// recvName33r11 renders a receiver type to the short name used as a roster key.
func recvName33r11(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.StarExpr:
		return recvName33r11(v.X)
	case *ast.Ident:
		return v.Name
	case *ast.IndexExpr:
		return recvName33r11(v.X)
	case *ast.IndexListExpr:
		return recvName33r11(v.X)
	default:
		return "unknownRecv"
	}
}

// exprText33r11 renders the receiver of a lock call (m.mu) so a violation can name
// what it found, without pulling a printer into a test.
func exprText33r11(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return exprText33r11(v.X) + "." + v.Sel.Name
	case *ast.StarExpr:
		return exprText33r11(v.X)
	case *ast.CallExpr:
		return exprText33r11(v.Fun) + "()"
	case *ast.IndexExpr:
		return exprText33r11(v.X) + "[..]"
	default:
		return "?"
	}
}

// packageSourceFiles33r11 lists this package's production sources - the directory
// the test binary runs in, minus _test.go, so the ruler judges shipped naming and not
// test scaffolding (the fixtures below would otherwise red-lock this file). An empty
// list is a failure, not a pass: a census that read no files proves nothing.
func packageSourceFiles33r11(t *testing.T) []string {
	t.Helper()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("the census could not read its own package directory: %v", err)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)
	if len(names) == 0 {
		t.Fatalf("the census found zero .go files in the package directory: the ruler has no subject")
	}
	return names
}

// censusPackage33r11 parses every production source and runs the census over them.
// A file that will not parse reddens the case rather than being skipped, and one
// token.FileSet serves parse and position lookup, because a position read through a
// different FileSet is a number nobody can trust.
func censusPackage33r11(t *testing.T, names []string) []lockedFunc33r11 {
	t.Helper()
	fset := token.NewFileSet()
	var hits []lockedFunc33r11
	for _, name := range names {
		src, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("census could not read %s: %v", name, err)
		}
		file, err := parser.ParseFile(fset, name, src, parser.SkipObjectResolution)
		if err != nil {
			t.Fatalf("census could not parse %s: %v - a source the ruler cannot read is not exempt", name, err)
		}
		hits = append(hits, lockedFuncsIn33r11(fset, file, name)...)
	}
	sort.Slice(hits, func(a, b int) bool {
		if hits[a].qualified == hits[b].qualified {
			return hits[a].file < hits[b].file
		}
		return hits[a].qualified < hits[b].qualified
	})
	return hits
}

// TestLockedSuffixNeverTakesTheLockItself33r11 is ticket 33 item (c): a
// Locked-suffixed func in this package's production sources must be rostered and
// must not take a lock in its own body.
func TestLockedSuffixNeverTakesTheLockItself33r11(t *testing.T) {
	names := packageSourceFiles33r11(t)
	hits := censusPackage33r11(t, names)

	t.Logf("33-r11 census: %d Locked-suffixed func(s) across %d production source file(s) of package main", len(hits), len(names))
	for _, hit := range hits {
		t.Logf("  %s in %s line %d: %d lock operation(s) in its own body %v",
			hit.qualified, hit.file, hit.line, len(hit.lockCalls), hit.lockCalls)
	}

	seen := make(map[string]bool, len(hits))
	for _, hit := range hits {
		seen[hit.qualified] = true
		why, rostered := lockedNamingRoster33r11[hit.qualified]
		if !rostered {
			t.Errorf("%s (%s line %d) carries the %q suffix but no roster entry: a new name that claims the caller holds a mutex must be reviewed and named in lockedNamingRoster33r11, or renamed",
				hit.qualified, hit.file, hit.line, lockedSuffix33r11)
		}
		if len(hit.lockCalls) > 0 {
			t.Errorf("%s (%s line %d) is named as if the caller holds the lock while its own body takes it (%v) - that is the trap ticket 33 recorded: a caller that trusts the name parks on a non-reentrant sync.Mutex. Roster note: %q",
				hit.qualified, hit.file, hit.line, hit.lockCalls, why)
		}
	}

	// An exemption whose func is gone must go with it, or the roster turns into a
	// list of things nobody checks.
	for qualified := range lockedNamingRoster33r11 {
		if !seen[qualified] {
			t.Errorf("lockedNamingRoster33r11 names %s, which this package no longer declares - drop the exemption or restore the method", qualified)
		}
	}
}

// completedWithin33r11 runs fn on its own goroutine and reports whether it returned
// before budget elapsed: the bounded-wait shape of
// internal/config/manager_223_test.go's completedWithin. An abandoned call cannot
// write back into anything the caller reads before releasing the lock, so a parked
// probe leaks one goroutine and nothing else.
func completedWithin33r11(budget time.Duration, fn func()) bool {
	done := make(chan struct{})
	go func() {
		defer close(done)
		fn()
	}()
	select {
	case <-done:
		return true
	case <-time.After(budget):
		return false
	}
}

// TestRosteredLockedMethodReturnsWhileCallerHoldsTheLock33r11 asks the census
// question of real code instead of source text: setPriorFocusLocked, the one rostered
// method reachable without a window, must complete while the caller holds m.mu. A
// self-locking body parks instead, and this case goes red rather than hanging the run.
//
// The sentinel handle is deliberately not a real window: with m.hwnd still 0 the
// method's own refusal check short-circuits before any user32 call, so the case runs
// headlessly - and the prevFocus read-back proves the body really ran instead of
// returning at its first line.
func TestRosteredLockedMethodReturnsWhileCallerHoldsTheLock33r11(t *testing.T) {
	m := NewPanelManager(nil, nil, t.TempDir())
	const sentinel = windows.HWND(0x33f1)

	m.mu.Lock()
	returned := completedWithin33r11(2*time.Second, func() { m.setPriorFocusLocked(sentinel) })
	recorded := m.prevFocus
	m.mu.Unlock()

	if !returned {
		t.Fatalf("setPriorFocusLocked parked while the caller held m.mu: the method takes that lock itself, so its Locked suffix is a trap and any caller that trusts it deadlocks")
	}
	if recorded != sentinel {
		t.Fatalf("the probe returned but recorded nothing (prevFocus=%v, want the sentinel): this case would be green without having called the method", recorded)
	}
	t.Logf("33-r11 behaviour: setPriorFocusLocked returned inside the budget while m.mu was held, and it did record the sample")
}

// fixtureHost33r11 is the behavioural positive control's bench: a method named as if
// the caller held the lock while it takes the lock itself. It lives in a _test.go
// file, which the census above deliberately does not read.
type fixtureHost33r11 struct {
	mu    sync.Mutex
	field int
}

func (f *fixtureHost33r11) parkedLocked() {
	f.mu.Lock()
	f.field++
	f.mu.Unlock()
}

// lockedNamingFixtureSrc33r11 is the source text the census's positive controls are
// parsed from: three names, one bad (Locked suffix plus a self-locking body), one
// honest (Locked suffix, no lock in the body), one outside the suffix's range
// (self-locking, no suffix). The bad name is planted HERE, in a string this file
// owns, never in a production file.
const lockedNamingFixtureSrc33r11 = `package fixture

type fixtureHost struct {
	mu    sync.Mutex
	field int
}

func (f *fixtureHost) parkedLocked() {
	f.mu.Lock()
	f.field++
	f.mu.Unlock()
}

func (f *fixtureHost) honestLocked() {
	f.field++
}

func (f *fixtureHost) selfLocksButSaysNothing() {
	f.mu.Lock()
	f.field++
	f.mu.Unlock()
}
`

// selfLockerOldNameFixtureSrc33r11 is the same question asked of the shape this leg
// removed: a method that still carries the old name and still locks itself.
const selfLockerOldNameFixtureSrc33r11 = `package fixture

type roundTripHost struct {
	mu sync.Mutex
	w  *int
}

func (m *roundTripHost) firstRoundTripLocked() float64 {
	m.mu.Lock()
	w := m.w
	m.mu.Unlock()
	if w == nil {
		return -1
	}
	return 1
}
`

// TestLockedNamingCensusBitesItsOwnFixture33r11 is the census's positive control: the
// same census run against a planted bad name must report it, must leave the honest
// name alone, and must not reach a self-locking method that carries no suffix.
func TestLockedNamingCensusBitesItsOwnFixture33r11(t *testing.T) {
	hits := censusParsedFixture33r11(t, "fixture33r11.go", lockedNamingFixtureSrc33r11)

	byName := make(map[string]lockedFunc33r11, len(hits))
	for _, hit := range hits {
		byName[hit.qualified] = hit
	}
	if len(hits) != 2 {
		t.Fatalf("census found %d Locked-suffixed func(s) in a fixture declaring exactly 2 (%v) - the name match itself is broken", len(hits), hits)
	}

	bad, ok := byName["fixtureHost.parkedLocked"]
	if !ok {
		t.Fatalf("the census did not see the planted bad name fixtureHost.parkedLocked at all: the ruler is blind, not green")
	}
	if len(bad.lockCalls) != 2 {
		t.Errorf("the planted bad name reported %d lock operation(s), want 2 (Lock and Unlock): %v", len(bad.lockCalls), bad.lockCalls)
	}
	if honest, ok := byName["fixtureHost.honestLocked"]; !ok || len(honest.lockCalls) != 0 {
		t.Errorf("the honest Locked-suffixed fixture came out self-locking or missing (%+v): the census would redden every future correct name", honest)
	}
	if _, ok := byName["fixtureHost.selfLocksButSaysNothing"]; ok {
		t.Errorf("the census reached a func with no Locked suffix - ticket 33 item (b) is a different ruling and this ruler must stay out of it")
	}
	t.Logf("33-r11 positive control: planted %s reported with %v, honest name clean, unsuffixed self-locker untouched",
		bad.qualified, bad.lockCalls)
}

// TestLockedNamingCensusRejectsTheOldName33r11 runs the positive control in the exact
// shape of the bug this leg renamed away: the pre-rename name plus the pre-rename
// body must be reported, and must not be exempt.
func TestLockedNamingCensusRejectsTheOldName33r11(t *testing.T) {
	hits := censusParsedFixture33r11(t, "oldname33r11.go", selfLockerOldNameFixtureSrc33r11)
	if len(hits) != 1 {
		t.Fatalf("expected the old name to be found once, got %v", hits)
	}
	if len(hits[0].lockCalls) == 0 {
		t.Fatalf("the old shape (Locked suffix, body takes m.mu) came out clean: the pin cannot bite the very shape ticket 33 recorded")
	}
	if _, rostered := lockedNamingRoster33r11[hits[0].qualified]; rostered {
		t.Fatalf("the pre-rename name is on the exemption roster: item (a) would not have to land")
	}
	t.Logf("33-r11 positive control: the pre-rename shape is reported as %s with %v", hits[0].qualified, hits[0].lockCalls)
}

// TestCompletedWithinReportsAParkedCall33r11 is the behavioural case's positive
// control: a method that takes the very lock the caller holds must come back as "did
// not return", otherwise the rostered-method case could be green on a broken probe.
func TestCompletedWithinReportsAParkedCall33r11(t *testing.T) {
	f := &fixtureHost33r11{}
	f.mu.Lock()
	returned := completedWithin33r11(200*time.Millisecond, func() { f.parkedLocked() })
	ranBeforeRelease := f.field
	f.mu.Unlock()

	if returned {
		t.Fatalf("parkedLocked completed while the caller held the same mutex: no case in this file can detect a self-locking Locked method")
	}
	if ranBeforeRelease != 0 {
		t.Fatalf("the parked call ran its body before parking (field=%d), so what the probe saw was not the deadlock it claims to detect", ranBeforeRelease)
	}
	t.Logf("33-r11 positive control: a self-locking Locked call parked, was reported as parked, and had touched nothing")
}

// censusParsedFixture33r11 parses one fixture source and runs the census on it with
// the same FileSet it reports positions through.
func censusParsedFixture33r11(t *testing.T, name, src string) []lockedFunc33r11 {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, name, src, parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("the positive-control fixture %s did not parse, so the census was never asked to bite: %v", name, err)
	}
	return lockedFuncsIn33r11(fset, file, name)
}
