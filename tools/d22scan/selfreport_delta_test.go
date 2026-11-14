package main

// Ticket 71 AC#4's falsifier, in the shape the acceptance actually asks for:
// plant ONE bad file into ONE scope, and require that scope's self-reported
// number to go up by exactly one AND the same run to go red naming the file.
//
// WHY A DELTA AND NOT A MAGNITUDE. TestLedgerCountsMatchAnIndependentWalk
// already compares every reported number against an independent walk of the
// same tree, which catches a WRONG number. It cannot catch a CONSTANT one: a
// count hardcoded to today's value satisfies an equality test on today's tree
// exactly once and the test is then dead weight. The question "is this number
// computed or printed?" is answered by moving the tree and watching the number
// move - so each case below measures the same scope twice, once before and once
// after exactly one file is planted, and asserts the difference is 1.
//
// WHY ALSO THE PRINTED LINE, not just the scanner's internal map. The map is
// not what CI consumes and not what a reader greps; the line
// `d22scan: scope <label> examined N <kind>` is. A future change that printed a
// stale number while keeping the map honest would pass every map-level check and
// still hand a human a lie, which is the disease ticket 71 was opened for. So
// both the map and the line are asserted for every scope.
//
// WHY THIS RUNS IN CI WITHOUT NEW WIRING: CI's lint job runs
// `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` (the positive control)
// BEFORE the scan, and that script runs every test in this package with -count=1.
// A scope whose number stops moving therefore turns the positive control red
// before the scan is allowed to report anything.

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

// scopeDeltaCase plants exactly one violation file into exactly one declared
// scope. rel is repo-relative and its directory part is unique per case, so the
// "the red line names the planted file" assertion cannot be satisfied by a
// finding about some other file.
type scopeDeltaCase struct {
	scope   string // label in declaredScopes(), matched verbatim
	rel     string // where the file is planted
	body    string // its contents: a violation of wantBan
	wantBan string // the finding that must name rel
}

func scopeDeltaCases() []scopeDeltaCase {
	// The glyph is written as a Go escape on purpose: this file must stay ASCII
	// (it is a Go source the D22 emoji ban could otherwise be asked about) while
	// the SEEDED file must contain the real character, because emojiRe scans raw
	// bytes and an escape sequence is not a glyph.
	emojiText := "smile more often \U0001F600\n"
	goEmoji := "package zz71delta\n\nconst banner = \"smile \U0001F600\"\n"
	goroutine := "package zz71delta\n\nfunc worker() {}\n\nfunc leak() { go worker() }\n"
	return []scopeDeltaCase{
		{
			scope:   "bans #1-5 internal/",
			rel:     "internal/zz71delta/goroutine.go",
			body:    goroutine,
			wantBan: "bare-goroutine",
		},
		{
			scope:   "bans #1-5 cmd/",
			rel:     "cmd/zz71delta/goroutine.go",
			body:    goroutine,
			wantBan: "bare-goroutine",
		},
		{
			// The same seed shape TestBuiltBinaryGoesRedEndToEnd uses for the
			// armed ban #6, for the same reason: a rename that dodged the
			// extension filter would be the instrument hole written into a test.
			scope:   "ban #6 frontend/",
			rel:     "frontend/src/zz71delta/decision.tsx",
			body:    "export function decide() { return approval.decide({ allow: true }); }\n",
			wantBan: "panel-approval",
		},
		{
			scope:   "ban #7 internal/tools/",
			rel:     "internal/tools/zz71delta.go",
			body:    "package tools\n\nvar names = []string{\"spill\", \"internal.logwrite\"}\n",
			wantBan: "internal-artifact-tool",
		},
		{
			scope:   "ban #10 internal/risk/",
			rel:     "internal/risk/zz71delta.go",
			body:    "package risk\n\n// aZz71Delta is cited by this comment and declared by no file in the package.\nfunc zz71delta() {}\n",
			wantBan: "phantom-symbol",
		},
		{
			scope:   "ban #8 design/",
			rel:     "design/zz71delta/ball.md",
			body:    emojiText,
			wantBan: "emoji",
		},
		{
			scope:   "ban #8 frontend/",
			rel:     "frontend/src/zz71delta/banner.txt",
			body:    emojiText,
			wantBan: "emoji",
		},
		{
			scope:   "ban #8 internal/",
			rel:     "internal/zz71delta/emoji.go",
			body:    goEmoji,
			wantBan: "emoji",
		},
		{
			scope:   "ban #8 cmd/",
			rel:     "cmd/zz71delta/emoji.go",
			body:    goEmoji,
			wantBan: "emoji",
		},
	}
}

// findScope returns the declared scope with this exact label.
func findScope(t *testing.T, scopes []scanScope, label string) scanScope {
	t.Helper()
	for _, sc := range scopes {
		if sc.label == label {
			return sc
		}
	}
	t.Fatalf("scope %q is not in declaredScopes() (have %v)", label, labels(scopes))
	return scanScope{}
}

// scopeLine returns the self-report line for one scope exactly as main() prints
// it. The match uses describeScopes()'s own padding, so a line about a DIFFERENT
// scope that happens to share a prefix ("internal/" appears in three labels)
// cannot satisfy the caller.
func scopeLine(out, label string) string {
	want := "d22scan: scope " + fmt.Sprintf("%-24s", label)
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, want) {
			return line
		}
	}
	return ""
}

// TestEveryDeclaredScopeCountsTheFileItIsGiven is the ticket's "the number is
// not a constant" proof for EVERY declared live scope.
func TestEveryDeclaredScopeCountsTheFileItIsGiven(t *testing.T) {
	cases := scopeDeltaCases()

	// Guard 1: coverage of the case list. This is the half that keeps the test
	// alive as the ledger grows - the next agent that declares a scope and
	// forgets to pin it gets a red here instead of an unpinned number.
	declared := declaredScopes(liveFixture(t))
	liveLabels := map[string]bool{}
	for _, sc := range declared {
		if sc.live {
			liveLabels[sc.label] = true
		}
	}
	cased := map[string]bool{}
	for _, c := range cases {
		if cased[c.scope] {
			t.Fatalf("scope %q has more than one case: the delta would be measured twice against the same scope", c.scope)
		}
		cased[c.scope] = true
		if !liveLabels[c.scope] {
			t.Errorf("case names scope %q, which a fully populated fixture does not declare live (declared: %v)", c.scope, labels(declared))
		}
	}
	for label := range liveLabels {
		if !cased[label] {
			t.Errorf("live scope %q has no plant-one-file case: its self-reported number is unpinned and could be a constant (ticket 71 AC#4)", label)
		}
	}

	for _, c := range cases {
		c := c
		t.Run(strings.NewReplacer("/", "_", " ", "_", "#", "").Replace(c.scope), func(t *testing.T) {
			root := liveFixture(t)

			// BEFORE: the fixture is green and the scope reports n0.
			before := scanFixture(t, root)
			n0 := findScope(t, declaredScopes(root), c.scope).count(before)
			outBefore, errBefore, codeBefore := runVerdict(t, root, before)
			if codeBefore != 0 {
				t.Fatalf("fixture must be green BEFORE the plant, so the red below is the plant's and not the fixture's: rc=%d\nstdout:\n%s\nstderr:\n%s", codeBefore, outBefore, errBefore)
			}
			lineBefore := scopeLine(outBefore, c.scope)
			if lineBefore == "" {
				t.Fatalf("no self-report line for scope %q in the green run:\n%s", c.scope, outBefore)
			}
			if !strings.Contains(lineBefore, fmt.Sprintf("examined %3d", n0)) {
				t.Errorf("the printed line does not carry the counter's value %d: %q", n0, lineBefore)
			}

			// Plant exactly one file, into that scope only.
			seedFile(t, root, c.rel, c.body)

			// AFTER: the same scope must report n0+1, and the run must be red
			// naming that one file under that one ban.
			after := scanFixture(t, root)
			n1 := findScope(t, declaredScopes(root), c.scope).count(after)
			if n1 != n0+1 {
				t.Errorf("planting ONE file at %s moved scope %q from %d to %d examined files; want exactly %d. A self-report that does not move with the tree is a constant, which is the one thing ticket 71 AC#4 forbids",
					c.rel, c.scope, n0, n1, n0+1)
			}
			out, errOut, code := runVerdict(t, root, after)
			if code != 1 {
				t.Errorf("a %s violation planted at %s must make the run red with rc=1, got rc=%d\nstdout:\n%s\nstderr:\n%s", c.wantBan, c.rel, code, out, errOut)
			}
			line := scopeLine(out, c.scope)
			if line == "" {
				t.Fatalf("no self-report line for scope %q in the red run:\n%s", c.scope, out)
			}
			if !strings.Contains(line, fmt.Sprintf("examined %3d", n1)) {
				t.Errorf("the red run's line for %q does not print the new count %d: %q - the number CI and a reader see would then disagree with the scan's own work", c.scope, n1, line)
			}
			wantPath := filepath.ToSlash(c.rel)
			if !strings.Contains(out+errOut, wantPath) {
				t.Errorf("the red output must NAME the planted file %s; stdout:\n%s\nstderr:\n%s", wantPath, out, errOut)
			}
			named := false
			for _, f := range after.findings {
				if f.Ban == c.wantBan && f.Path == wantPath {
					named = true
				}
			}
			if !named {
				t.Errorf("want one %s finding against %s; got %v", c.wantBan, wantPath, after.findings)
			}
			t.Logf("%s: examined %d -> %d, rc 0 -> %d, line %q", c.scope, n0, n1, code, line)
		})
	}
}
