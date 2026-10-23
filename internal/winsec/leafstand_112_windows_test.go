//go:build windows

// Ticket 112, leg 112-r1: the deterministic denominator for a leg the gate currently
// measures by winning a race.
//
// WHY THIS FILE EXISTS (measured, not guessed). The windows gate's step went red on the
// runner on 2026-10-23 in two cases, and the one that belongs to the leaf leg is
// openleaf_guard_322_windows_test.go:240: it needs the leaf guard to speak at least once
// inside 400 raced attempts through PrivateFile. On the runner it spoke 0 times and the
// case printed "nothing below proves the guard is on this route". The same case speaks
// 8 to 15 times in 400 on the laptop. The difference is not a product byte: it is which
// leg answers first. On the runner the installed C26 resolver's steady-state refusal
// answered 288/400 of the same shape (sibling case ancestoropen_342_windows_test.go:48),
// on the laptop 205/400. A number that moves with the machine's speed cannot be the
// proof that a check is wired in, and it cannot be the thing that decides the gate.
//
// WHAT IS PINNED HERE. The same claim the race was hoping to sample, asked of the
// production route with no racer at all: when a link is ALREADY standing at the leaf,
// PrivateFile must refuse it before a byte is written and must name the leg that
// answered. That is the steady-state half of ticket 322's leaf leg (it is what makes a
// race result interpretable at all: the race only asks about the gap between the answer
// and the act), and unlike a race it cannot pass by luck and cannot pass on a machine
// that is merely idle.
//
// WHY IT IS NOT A "SHELF" TEST. Every call below goes through PrivateFile, the exported
// route, so deleting the guard or the floor from winsec.go's privateFile turns this file
// red by name. The direct guard call lives in the 322 file and proves the function only.
//
// WHAT IT DOES NOT PROVE, stated rather than glossed: that the guard is reached in the
// gap after the floor has answered. That claim is 322's race and it stays rate-bound;
// proving it by construction needs an injection point in winsec.go, which is a product
// change and ticket 346's form 乙 (another ticket plus human approval). Nothing here is
// read as that.
package winsec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The three legs that may answer "a link already stands at this name", each pinned by its
// own verbatim rather than by "it returned an error": 120-r3 measured that a wrong leg
// answering reads green when the assertion is only "err != nil".
const (
	floorLeafRefusal112r1 = "traverses a reparse point at"   // placement_windows.go's pathPieces walk
	c26LeafRefusal112r1   = "path traverses a reparse point" // internal/risk, steady state before any act
	guardLeafRefusal112r1 = whyGuardRefusal322r1             // preOpenLeafGuard, in the gap after the answer
	leafWitnessBody112r1  = "bytes a truncating open would have taken"
)

// passthrough112r1 answers the caller's spelling unchanged, which is what a steady-state
// resolver does for a path that has not moved. It is installed through the package's own
// test seam (export_test.go) so the floor behind it is measured on its own: with the real
// C26 resolver in place the floor's leaf walk is masked by the resolver's own refusal, and
// a masked leg is the exact quantity this ticket is about.
type passthrough112r1 struct{}

func (passthrough112r1) Resolve(input string) (string, error) { return input, nil }

// ResolveAccounted is the seam's second question (ticket 108 AC#4). A resolver that
// cannot answer it may not be installed at all, and ResolvePath refuses to act on its
// answer even if something got it past that rule - measured here: without this method
// every call below was refused with "does not implement RewriteAccounted", which is
// itself the guard ticket 108 added, answering in the shape of a fixture mistake.
func (passthrough112r1) ResolveAccounted(input string) (string, bool, error) {
	return input, false, nil
}

// legThatAnswered112r1 attributes a refusal instead of asserting "some error happened".
// The guard's reason is checked first because it is a substring of nothing else; a leg
// outside the three is reported as such rather than counted as a pass.
func legThatAnswered112r1(err error) (string, bool) {
	if err == nil {
		return "none: the call returned nil", false
	}
	msg := err.Error()
	switch {
	case strings.Contains(msg, guardLeafRefusal112r1):
		return "preOpenLeafGuard (the gap after the answer)", true
	case strings.Contains(msg, c26LeafRefusal112r1):
		return "the installed C26 resolver's steady-state refusal", true
	case strings.Contains(msg, floorLeafRefusal112r1):
		return "the placement floor's own reparse walk", true
	}
	return "some other failure: " + msg, false
}

// plantLeafLink112r1 builds the shape once and checks the fixture itself before it is
// used: a link that is not a reparse point would make every reading below meaningless,
// and mkFileSymlink already states that as a SKIP that acceptance has to count.
func plantLeafLink112r1(t *testing.T) (dir, link, witness, body string) {
	t.Helper()
	root := SealableTempDirForTest124(t)
	dir, outside := filepath.Join(root, "tree"), filepath.Join(root, "outside-tree")
	// Joined rather than concatenated: t.TempDir() inherits GOTMPDIR's separator style,
	// and a forward-slash spelling reaches cmd's mklink as `Invalid switch` (322 file).
	for _, d := range []string{dir, outside} {
		if err := os.MkdirAll(d, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	witness = filepath.Join(outside, "someone-elses-file-112r1.bin")
	body = leafWitnessBody112r1
	if err := os.WriteFile(witness, []byte(body), 0o666); err != nil {
		t.Fatal(err)
	}
	link = filepath.Join(dir, "leaf-link-112r1.bin")
	mkFileSymlink(t, link, witness)
	if !reparseAt(link) {
		t.Fatalf("the leaf %s is not a reparse point, so this fixture measured nothing", filepath.Base(link))
	}
	return dir, link, witness, body
}

// witnessIntact112r1 is the damage half: a refusal that still truncated the target would
// be a refusal after the act, which is the shape ticket 322 exists to remove.
func witnessIntact112r1(t *testing.T, witness, body string) {
	t.Helper()
	got, err := os.ReadFile(witness)
	if err != nil || string(got) != body {
		t.Errorf("AC#1 RED: the leaf call damaged the object behind the link: bytes=%q err=%v, want %q",
			got, err, body)
	}
}

// TestAC1AStandingLeafLinkIsRefusedByTheFloorBeforeAnyByte112r1 asks the floor on its own.
// With a pass-through resolver installed the resolver cannot answer, so the only legs
// between the caller and the open are the floor's leaf walk and the guard. The floor is
// asked first by winsec.go's ResolvePath (it re-checks the resolver's own answer), so its
// verbatim - not the guard's - is what this leg pins. Removing the leaf from that walk
// moves the answer to the guard and this case goes red naming it.
func TestAC1AStandingLeafLinkIsRefusedByTheFloorBeforeAnyByte112r1(t *testing.T) {
	_, link, witness, body := plantLeafLink112r1(t)
	restore := SetSeamForTest(passthrough112r1{})
	defer restore()

	err := PrivateFile(link, []byte("the route's own payload"), 0o600)
	leg, named := legThatAnswered112r1(err)
	if !named {
		t.Fatalf("AC#1 RED: the standing leaf link was answered by %s; want one of the three named legs", leg)
	}
	if !strings.Contains(err.Error(), floorLeafRefusal112r1) {
		t.Errorf("AC#1 RED: the leg that answered was %s, want the placement floor's reparse walk (%q);"+
			" verbatim: %v", leg, floorLeafRefusal112r1, err)
	}
	if !errors.Is(err, ErrUnresolvedPath) {
		t.Errorf("AC#1 RED: the refusal lost ErrUnresolvedPath: %v", err)
	}
	witnessIntact112r1(t, witness, body)
	t.Logf("AC#1 leg that answered = %s; verbatim: %v", leg, err)
}

// TestAC2TheSameRouteStillWritesARealName112r1 is the other direction, and it is the one
// that keeps the case above from passing on a route that refuses everything: the same
// pass-through resolver, the same tree, an ordinary name. The four production callers
// (internal/secret, internal/config, internal/agent) depend on this branch staying open.
func TestAC2TheSameRouteStillWritesARealName112r1(t *testing.T) {
	dir, _, witness, body := plantLeafLink112r1(t)
	restore := SetSeamForTest(passthrough112r1{})
	defer restore()

	real := filepath.Join(dir, "plain-112r1.bin")
	const payload = "the tree's own bytes"
	if err := PrivateFile(real, []byte(payload), 0o600); err != nil {
		t.Fatalf("AC#2 RED: PrivateFile refused an ordinary name (%v), which would break every caller", err)
	}
	if got, err := os.ReadFile(real); err != nil || string(got) != payload {
		t.Errorf("AC#2 RED: the file holds %q (err %v), want %q", got, err, payload)
	}
	witnessIntact112r1(t, witness, body)
}

// TestAC3TheShippedResolverAlsoRefusesAStandingLeafLink112r1 runs the same shape against
// whatever is installed (internal/risk's C26 pipeline when it is linked, the floor when it
// is not). It pins the conclusion the gate's red run needed: on the shipped configuration
// a link that already stands is refused before any byte, and the answering leg is named in
// the log rather than assumed. It does not pin WHICH leg answers, because both answers are
// "before any act" and 112 job log 120189170762 shows the runner's leg is the resolver.
func TestAC3TheShippedResolverAlsoRefusesAStandingLeafLink112r1(t *testing.T) {
	_, link, witness, body := plantLeafLink112r1(t)

	err := PrivateFile(link, []byte("the route's own payload"), 0o600)
	leg, named := legThatAnswered112r1(err)
	if !named {
		t.Fatalf("AC#3 RED: with the shipped resolver installed the standing leaf link was answered by %s", leg)
	}
	// No ErrUnresolvedPath assertion on this leg, and that is a reading rather than a
	// relaxation: resolveAccounted propagates a resolver's error AS-IS on purpose ("so a
	// refusal stays attributable to C26 rather than to this package's opinion"), so this
	// leg carries internal/risk's own sentinel and not winsec's. The floor leg above does
	// carry it, which is what distinguishes the two answers.
	witnessIntact112r1(t, witness, body)
	t.Logf("AC#3 shipped configuration: leg that answered = %s; seamLatched=%v; verbatim: %v",
		leg, SeamLatchedForTest(), err)
}

// TestAC4TheLegThisPinsIsTheOneTheRunnerLost112r1 is the bridge between this file and the
// gate: it reads the same 400-attempt race the runner failed, but it does not judge on it.
// What it pins is the structure that makes the runner's 0/400 interpretable - a standing
// link is always answered by a steady-state leg (measured above), so a race that never
// reaches the guard is a statement about the gap, not about the guard. The gate's colour
// stays 322's to set; this case fails only if the steady-state answer disappears, which is
// the one change that would make an abstention mean "no check".
func TestAC4TheLegThisPinsIsTheOneTheRunnerLost112r1(t *testing.T) {
	dir, _, _, _ := plantLeafLink112r1(t)
	answered := map[string]int{}
	for i := 0; i < 20; i++ {
		link := filepath.Join(dir, fmt.Sprintf("standing-%02d-112r1.bin", i))
		target := filepath.Join(dir, "ordinary-target-112r1.bin")
		if err := os.WriteFile(target, []byte("an ordinary object"), 0o600); err != nil {
			t.Fatal(err)
		}
		mkFileSymlink(t, link, target)
		leg, named := legThatAnswered112r1(PrivateFile(link, []byte("payload"), 0o600))
		if !named {
			t.Fatalf("AC#4 RED: attempt %d: the standing link was answered by %s", i, leg)
		}
		answered[leg]++
	}
	// A denominator, not a total: 20 attempts were made and 20 answers were attributed,
	// which is the guard against "the number got bigger" being read as progress.
	if answered["none: the call returned nil"] != 0 {
		t.Fatalf("AC#4 RED: %d of 20 standing links were accepted", answered["none: the call returned nil"])
	}
	for leg, n := range answered {
		t.Logf("AC#4 standing-link refusals: %-55s %2d/20", leg, n)
	}
}
