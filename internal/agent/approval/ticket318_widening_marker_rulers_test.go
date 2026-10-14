package approval

// Ticket 318's rulers: the widening chain must refuse a card whose Paths carry
// the C26 "无法规范化" marker. These are unit-level pins for the same shape
// zz_201_a1_junction_marker_windows_test.go measures end-to-end (real junction →
// real bridge → real Gate/Queue/Replies → real WideningRule): here the card is
// built directly, which keeps the reading on every platform and names the exact
// string the production marker is.
//
// Why the card face can carry the marker at all: internal/tools/bridge.go:1028
// renders a raw path C26 refused as `p + " (无法规范化: " + err.Error() + ")"`,
// and that string is what reaches Prompt.Paths → ReplyCard.Paths verbatim
// (gate.go promptFor → replies.go Record → cmd/wisp/run.go's live.record).
// WideningRule's dirOf (batch.go:91) splits by the LAST '/' or '\', and the
// junction marker's own text ("… (junction/symlink) …", internal/risk/
// pathresolver.go:36) carries a slash, so before this ticket the split fell
// inside the marker and "…(junction" came back as a directory (live reading:
// .scratch/wisp/probes/201/a1/logs/run3-full.txt, READING-5/6).

import (
	"strings"
	"testing"
)

// t318ReparseMarker is byte-for-byte what bridge.go:1028 appends for the one
// canonicalization error today's C26 emits for a refused reparse traversal:
// risk.ErrReparseDenied.Text() (internal/risk/pathresolver.go:36).
const t318ReparseMarker = " (无法规范化: risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions)"

// TestTicket318AMarkerCardOffersNoWideningTarget is the ticket's 复现形 at card
// level, three cases: the junction card 201-a1 measured (forward slashes), the
// same chain in backslash spelling (the pre-fix probe proved both spellings
// fragmented identically), and the empty-path marker, whose string has no slash
// at all - there dirOf handed back the WHOLE marker as a "directory" (AC#2's
// first boundary, red before this fix for that other reason).
func TestTicket318AMarkerCardOffersNoWideningTarget(t *testing.T) {
	cases := []struct {
		name string
		path string
	}{
		{
			"junction card as 201-a1 measured it (forward slashes)",
			"C:/proj/allowed/jn/victim.txt" + t318ReparseMarker,
		},
		{
			"the same chain in backslash spelling",
			`C:\proj\allowed\jn\victim.txt` + t318ReparseMarker,
		},
		{
			"empty-path marker (tools: empty path), no slash in the string",
			" (无法规范化: tools: empty path)",
		},
	}
	for _, tc := range cases {
		card := ReplyCard{CorrelationID: tc.name, Level: "L2", Tool: "fs.read", Paths: []string{tc.path}}
		dir, rule, ok := card.WideningRule()
		if ok {
			t.Errorf("%s: marked path still offered as a widening base: dir=%q rule=%q", tc.name, dir, rule)
		}
		if dir != "" || rule != "" {
			t.Errorf("%s: a refusal must clear dir and rule, got (%q, %q)", tc.name, dir, rule)
		}
	}
}

// TestTicket318AMarkedPathRefusesTheWholeCard pins the refusal's scope. This
// function's contract is "exactly one distinct directory yields one line, and
// anything wider is refused rather than guessed at": a card that is only
// reducible by IGNORING the path C26 could not canonicalize is a card whose
// rule would be derived from a remainder the operator did not see named, so the
// whole card is refused - the any-path shape session.Ledger.Covering already
// uses (internal/session/grants.go:193-197: one unrecordable path answers the
// whole question false) and bridge.scopeVerdictOf repeats (internal/tools/
// bridge.go:1273). The mixed card below also keeps its reading under the
// pre-fix tree (two distinct "dirs" already refused), so this pin's teeth are
// the per-path mutation in probes/318/r1/mutation (a `continue` here flips it
// red).
func TestTicket318AMarkedPathRefusesTheWholeCard(t *testing.T) {
	mixed := ReplyCard{Paths: []string{
		"C:/proj/allowed/jn/victim.txt" + t318ReparseMarker,
		"C:/proj/outside/victim.txt",
	}}
	if dir, rule, ok := mixed.WideningRule(); ok {
		t.Errorf("one marked path must refuse the whole card, got dir=%q rule=%q", dir, rule)
	}
	twoMarked := ReplyCard{Paths: []string{
		"C:/proj/allowed/jn/victim.txt" + t318ReparseMarker,
		`C:\proj\allowed\jn\other.txt` + t318ReparseMarker,
	}}
	if _, _, ok := twoMarked.WideningRule(); ok {
		t.Errorf("a card whose every path is a marker can offer no rule")
	}
}

// TestTicket318TheSamePathWithoutTheMarkerStillDerivesItsParent is the
// no-regression control, sharpened to the same path as the first case: drop the
// marker and the rule must come back exactly as before (matrix: with marker →
// refuse; without → true parent). The end-to-end half of this control - a real
// out-of-scope read through the real bridge - stays in the zz probe test.
func TestTicket318TheSamePathWithoutTheMarkerStillDerivesItsParent(t *testing.T) {
	card := ReplyCard{Paths: []string{"C:/proj/allowed/jn/victim.txt"}}
	dir, rule, ok := card.WideningRule()
	if !ok {
		t.Fatalf("an unmarked single-directory card must still offer a rule, got (%q,%q,%v)", dir, rule, ok)
	}
	if dir != "C:/proj/allowed/jn" {
		t.Errorf("dir = %q, want the card's own parent segment", dir)
	}
	if rule != `[fs] allowed_dirs += "C:/proj/allowed/jn"` {
		t.Errorf("rule text = %q, want the line exactly as config.toml will hold it", rule)
	}
}

// TestTicket318ShortNamesAreNotTheMarker is AC#2's second boundary, named as a
// non-case: an 8.3 short name is NOT a canonicalization failure (C26 folds it
// and displayPaths returns the long form - ticket 252's paths.go:139-155), so no
// marker is ever produced for it and this guard is blind to the string by
// design. Pinned so "refuses every Windows-looking spelling" cannot be read as
// this fix's shape.
func TestTicket318ShortNamesAreNotTheMarker(t *testing.T) {
	card := ReplyCard{Paths: []string{"C:/PROGRA~1/notes/x.txt"}}
	dir, _, ok := card.WideningRule()
	if !ok || dir != "C:/PROGRA~1/notes" {
		t.Errorf("a short-name spelling without the marker must behave exactly as before, got (%q,%v)", dir, ok)
	}
	if strings.Contains(dir, "无法规范化") {
		t.Errorf("control fixture is polluted: %q", dir)
	}
}
