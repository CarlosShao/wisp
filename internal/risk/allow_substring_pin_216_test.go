package risk

import (
	"os"
	"strconv"
	"strings"
	"testing"
)

// Ticket 216 AC#6 arm A (the operator pin). AC#3 half-sentence
// "a verdict must never change because a name spells 已允许 inside it" had no
// instrument attacking it, so this file pins the ALLOW side of the judgement
// chain structurally: the three places that can wave a call through are read
// off the source by FUNCTION NAME and their bodies must hold zero arbitrary
// substring operators.
//
// Pinned targets (the only three ALLOW-side spellings in this package):
//   - rules_shell.go  shellAllowlisted  - argv[0] vs the shell allowlist
//   - rules_network.go hostAllowlisted  - host vs the domain allowlist
//   - blacklist.go     overrideApplies  - the B-tier single-file override key
//
// Deliberately NOT done here: absolute line numbers (tickets 327/329 moved
// lines in this repo), and a silent pass when a window comes out empty. Both
// failure modes are loud: a renamed or moved function, or a window that does
// not contain the exact-equality construct it is supposed to pin, fails the
// test. Proofs of non-tautology are the mutation readings recorded in
// docs/evidence/s1/216-allow-substring-pin-r2.md, not a claim in this file.

// substringOperators216 is AC#6 arm A's operator set, verbatim from the
// ticket: strings.(Contains|ContainsAny|ContainsRune|Index|IndexAny).
var substringOperators216 = []string{"Contains", "ContainsAny", "ContainsRune", "Index", "IndexAny"}

// allowSidePin216 names one ALLOW-side function and the exact-match construct
// its body must contain, so a window taken from the wrong place (or an empty
// one) cannot read as green.
type allowSidePin216 struct {
	file    string
	fn      string
	markers []string
}

var allowSidePins216 = []allowSidePin216{
	{
		file: "rules_shell.go", fn: "shellAllowlisted",
		markers: []string{"strings.EqualFold"},
	},
	{
		file: "rules_network.go", fn: "hostAllowlisted",
		markers: []string{"host == entry", "strings.HasSuffix"},
	},
	{
		file: "blacklist.go", fn: "overrideApplies",
		markers: []string{"bOverrides["},
	},
}

// TestAllowSideFunctionsCarryNoSubstringOperator is AC#6 arm A.
func TestAllowSideFunctionsCarryNoSubstringOperator(t *testing.T) {
	for _, pin := range allowSidePins216 {
		pin := pin
		t.Run(pin.fn, func(t *testing.T) {
			data, err := os.ReadFile(pin.file)
			if err != nil {
				t.Fatalf("read %s from the package directory: %v (go test runs with the "+
					"package dir as working directory; this pin reads source, it does not skip)",
					pin.file, err)
			}
			window, startLine := functionWindow216(t, string(data), pin.file, pin.fn)
			joined := strings.Join(window, "\n")
			markers := 0
			for _, m := range pin.markers {
				if strings.Contains(joined, m) {
					markers++
				}
			}
			if markers == 0 {
				t.Fatalf("%s: the window taken for %s (%d lines from line %d) holds none of %v - "+
					"the window is not the body this pin is about, so a 0 here would be decoration",
					pin.file, pin.fn, len(window), startLine, pin.markers)
			}
			hits := substringOperatorHits216(window, pin.file, startLine)
			if len(hits) != 0 {
				t.Errorf("%s: %s holds %d arbitrary-substring operator(s), want 0 (AC#6 arm A). "+
					"A name that merely CONTAINS an allowlist entry would be waved through, which is "+
					"exactly the shape AC#3 forbids. Hits:\n\t%s",
					pin.file, pin.fn, len(hits), strings.Join(hits, "\n\t"))
			}
			t.Logf("pinned %s:%s window=%d lines, exact-match markers=%d, substring operators=%d",
				pin.file, pin.fn, len(window), markers, len(hits))
		})
	}
}

// TestSubstringOperatorCounter216IsNotBlind is the counter's own positive
// control: it must see each of the five operators by their full identifier,
// and must not be triggered by the neighbouring operators that this package
// legitimately uses (LastIndexAny / ToLower / HasSuffix). Without this, the
// pin above could report 0 because it looks rather than because it sees.
func TestSubstringOperatorCounter216IsNotBlind(t *testing.T) {
	for _, op := range substringOperators216 {
		op := op
		t.Run(op, func(t *testing.T) {
			line := "\t\tif strings." + op + "(argv0, entry) {"
			got := substringOperatorHits216([]string{line}, "fixture", 1)
			if len(got) != 1 {
				t.Fatalf("strings.%s written out in full must be counted, got %d hits %v", op, len(got), got)
			}
		})
	}
	t.Run("lookalikes-are-not-operators", func(t *testing.T) {
		for _, line := range []string{
			"\tif i := strings.LastIndexAny(argv0, ` /\\`); i >= 0 {",
			"\thost = strings.ToLower(strings.TrimSpace(entry))",
			"\tif host == entry || strings.HasSuffix(host, \".\"+entry) {",
			"\treturn bOverrides[cand.real]",
		} {
			if got := substringOperatorHits216([]string{line}, "fixture", 1); len(got) != 0 {
				t.Errorf("lookalike line %q counted %d hit(s), want 0: %v", line, len(got), got)
			}
		}
	})
}

// functionWindow216 returns the source lines of one top-level function, found
// by name. It fails the test when the name is gone: the pin's value depends on
// being attached to the function it names, not to a line number.
func functionWindow216(t *testing.T, src, file, fn string) ([]string, int) {
	t.Helper()
	lines := strings.Split(src, "\n")
	prefix := "func " + fn + "("
	start := -1
	for i, ln := range lines {
		if strings.HasPrefix(ln, prefix) {
			start = i
			break
		}
	}
	if start < 0 {
		t.Fatalf("%s: no top-level function named %q (searched every line for %q). "+
			"The AC#6 arm A pin is bound to this name on purpose: renaming or moving an "+
			"ALLOW-side matcher must fail here, not silently stop being checked.",
			file, fn, prefix)
	}
	depth, opened := 0, false
	for i := start; i < len(lines); i++ {
		for _, r := range lines[i] {
			switch r {
			case '{':
				depth++
				opened = true
			case '}':
				depth--
			}
		}
		if opened && depth <= 0 {
			window := make([]string, 0, i-start+1)
			for _, ln := range lines[start : i+1] {
				window = append(window, strings.TrimRight(ln, "\r"))
			}
			if last := window[len(window)-1]; last != "}" {
				t.Fatalf("%s: the window for %s ends on %q, want a lone closing brace - the "+
					"brace walk lost the body (a brace inside a literal here would do that)",
					file, fn, last)
			}
			return window, start + 1
		}
	}
	t.Fatalf("%s: function %s never closes before end of file: window unterminated", file, fn)
	return nil, start + 1
}

// substringOperatorHits216 lists every strings.<operator> use inside the
// window, one entry per file:line:content hit. An operator is counted only
// when the identifier right after "strings." is exactly one of the five, so
// LastIndexAny and friends are not mistaken for IndexAny.
func substringOperatorHits216(window []string, file string, startLine int) []string {
	var hits []string
	for offset, ln := range window {
		for rest := ln; ; {
			i := strings.Index(rest, "strings.")
			if i < 0 {
				break
			}
			rest = rest[i+len("strings."):]
			end := 0
			for end < len(rest) {
				c := rest[end]
				if c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
					(end > 0 && c >= '0' && c <= '9') {
					end++
					continue
				}
				break
			}
			name := rest[:end]
			for _, op := range substringOperators216 {
				if name == op {
					hits = append(hits, file+":"+strconv.Itoa(startLine+offset)+": "+strings.TrimSpace(ln))
				}
			}
			rest = rest[end:]
		}
	}
	return hits
}
