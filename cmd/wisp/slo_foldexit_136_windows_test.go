//go:build windows

package main

// Ticket 136 AC#10 (criteria rewritten 2026-09-24, ticket face :176-193): the
// zero-sample fail-closed guard of internal/observe/sampler.go has its
// in-package nail (AC#1), but the chain "a red gate must reach the CLI" had
// only file:line citations. The preflight auditor measured the catcher count
// as zero:
//
//	grep -rn "run\.Pass" cmd/wisp/*_test.go  =>  0 hits
//
// so a mutation at either load-bearing site changed not one test reading. The
// two sites in cmd/wisp/slo_windows.go are:
//
//   - the fold where the out-of-tree state verdict reaches the run's pass bit
//     (ticket-face :386, on-disk :399): run.Pass = rep.Pass && observer.Pass;
//   - the chain end where a false run.Pass maps to exit code 1 (ticket-face
//     :323-325, on-disk :336-338) - the number scripts/slo-check.ps1 and the
//     CI jobs actually read.
//
// This file is the first catcher on the `go test` / gate face. The second
// catcher stays the CI leak self-test (scripts/slo-check.ps1:594
// leakFlipped, run by ci.yml's slo-smoke job): per the rewritten criteria this
// leg installs "go test first catcher, CI second catcher", not "first
// instrument". It reads THIS package's own production bytes (the
// Test44R3RunGoAppliesTheWiring idiom, costwiring_44r3_test.go:94) and fails
// closed when a needle no longer names the shipped shape, instead of staying
// green over a drifted source.
//
// Deliberately NOT in this file: a real `wisp slo` process run forced into a
// zero-sample window. cmdSLO boots proc.Boot + observe.InitLog and builds its
// samplers directly (observe.NewSampler over proc.NewTreeSampler /
// proc.NewExternalSampler, slo_windows.go :283 and :385) with no injection
// seam, so that shape is unreachable from outside internal/observe without
// adding production code or a test backdoor env - a seam decision under
// AGENTS §1.3 that this ticket registered (:189-191) rather than made. What
// IS real on this host (the rebuilt-exe reading of criteria 1) lives in
// docs/evidence/s1/136-ac7-ac10.md.
//
// No threshold, golden, budget or production byte is touched by this file: it
// reads source shapes and asserts nothing about any SLO number.

import (
	"os"
	"strings"
	"testing"
)

// slo136FoldLine is the one line where the out-of-tree state verdict reaches
// the run's pass bit. Both halves are load-bearing: dropping rep.Pass would
// publish a red state report as a pass, dropping observer.Pass would publish
// a red observer record (the fail-closed report collector of ticket 144) as a
// pass, and turning the conjunction into a disjunction would collapse two
// independent red sources into one.
const slo136FoldLine = "run.Pass = rep.Pass && observer.Pass"

// slo136ExitBlock is the chain end: a false run.Pass must map to exit code 1
// exactly here, where the gate scripts are the readers. Exit 2 stays the
// usage/environment code; this block is the verdict path.
const slo136ExitBlock = "if !run.Pass {\n\t\texitCode = 1\n\t}"

// TestSLO136StatePassFoldReachesRunPass pins the fold at slo_windows.go:399.
// It carries its OWN needles: when this fold moves, the reading names only
// this case (the exit-mapping case stays green over an unrelated fold change -
// attribution was measured on the first M-A landing, which reded both cases
// through a shared helper and was reworked for exactly that reason).
func TestSLO136StatePassFoldReachesRunPass(t *testing.T) {
	src := slo136ReadSelf(t)

	if got := strings.Count(src, slo136FoldLine); got != 1 {
		t.Fatalf("slo_windows.go carries the state pass fold %q %d times, want exactly 1 (a second fold would fork the verdict, none would orphan it)", slo136FoldLine, got)
	}
	for _, half := range []string{"rep.Pass", "observer.Pass"} {
		if !strings.Contains(slo136FoldLine, half) {
			t.Fatalf("the pinned fold line lost its %q half: that is a drift of this instrument, not of the source - re-pin against the shipped bytes before trusting any green here", half)
		}
	}
	if seg := slo136TopLevelSegment(t, src, "func runOutOfTree("); !strings.Contains(seg, slo136FoldLine) {
		t.Fatalf("the state pass fold is no longer inside runOutOfTree - the out-of-tree acceptance path this ticket nailed has moved")
	}

	// The subject-in-tree write is a different seam (the observer-cost record)
	// and must stay exactly one plain assignment.
	if got := strings.Count(src, "run.Pass = rep.Pass\n"); got != 1 {
		t.Fatalf("run.Pass = rep.Pass (the subject-in-tree write) appears %d times, want exactly 1", got)
	}
}

// TestSLO136FalseRunPassExitsOne pins the exit mapping at :336-338 with its
// OWN needles only: an unrelated fold change keeps this case green and this
// case alone reds when the mapping moves.
func TestSLO136FalseRunPassExitsOne(t *testing.T) {
	src := slo136ReadSelf(t)

	if got := strings.Count(src, slo136ExitBlock); got != 1 {
		t.Fatalf("the run.Pass -> exit code 1 block appears %d times in slo_windows.go, want exactly 1", got)
	}

	// Every mode writes its verdict before cmdSLO reads it: the settle fold
	// stays inside cmdSLO, the state fold stays outside it (it belongs to
	// runOutOfTree), and the mapping sits inside cmdSLO after both.
	seg := slo136TopLevelSegment(t, src, "func cmdSLO(")
	if !strings.Contains(seg, "run.Pass = run.Settle.Pass") {
		t.Fatalf("the settle fold left cmdSLO - the settle mode's verdict no longer reaches the exit mapping's prefix")
	}
	if !strings.Contains(seg, slo136ExitBlock) {
		t.Fatalf("the run.Pass -> exit code 1 mapping is no longer inside cmdSLO")
	}
	if strings.Contains(seg, slo136FoldLine) {
		t.Fatalf("the state pass fold moved into cmdSLO - the fold and the exit mapping must stay in their own functions so each mode's verdict path is readable")
	}
}

// slo136ReadSelf reads the production sibling from the package directory (the
// test binary's working directory is cmd/wisp). It grades nothing: every
// verdict lives in the case that owns the pin, so each mutation reds exactly
// the case that watches it.
func slo136ReadSelf(t *testing.T) string {
	t.Helper()
	raw, err := os.ReadFile("slo_windows.go")
	if err != nil {
		t.Fatalf("read cmd/wisp/slo_windows.go: %v - the AC#10 instrument cannot grade a source it cannot read", err)
	}
	return string(raw)
}

// slo136TopLevelSegment returns the bytes from a top-level func declaration to
// just before the next top-level func (indented nested funcs do not end the
// segment). It is a shape helper and carries no verdict of its own.
func slo136TopLevelSegment(t *testing.T, src, decl string) string {
	t.Helper()
	start := strings.Index(src, decl)
	if start < 0 {
		t.Fatalf("%q not found in slo_windows.go - unexpected file shape", decl)
	}
	next := strings.Index(src[start+1:], "\nfunc ")
	if next < 0 {
		t.Fatalf("no following top-level func after %q - unexpected file shape", decl)
	}
	return src[start : start+1+next]
}
