package main

// Ticket 143: `wisp panel-assets -l2` could not produce a card carrying the C19
// taint rule at all, because the assessor it builds was bare - internal/risk's
// ruleTaint returns nil when no TaintDetector is wired, and nothing on this
// path wired one. The production bridge does (internal/tools/bridge.go's
// assessorFor), so a demo machine and a real machine judged by different rules.
//
// What these cases pin, in the order the ticket's ACs ask for it:
//
//   - AC#1/AC#2: the flag is an INPUT (which source was read, from where, with
//     what content) and the card's hit, its level and its explanation sentence
//     come out of internal/risk. The load-bearing reading is that the SOURCE
//     NAME printed on the card is the one the flag declared: cmd/wisp never
//     sees the sentence shape and could not have written it.
//   - AC#2's other half: without the flag the printed card is what it was
//     before this file existed - no taint clause anywhere in it.
//   - the over-match guard: declaring a source is not the same as hitting, so a
//     detector that answered "hit" unconditionally goes red here (A218's family
//     of "too eager" mutations, and the reason AC#3(i) alone is not enough).
//   - the P9 red line, made mechanical: no verdict-shaped string literal lives
//     in panel_assets.go at all.
//
// No case here reads the network, the registry or another package's data. Each
// one calls cmdPanelAssets directly with os.Stdout swapped to a temp file, which
// is the same entry func main dispatches for this leg (ticket 133's coverage
// ledger reads that as `covered=test`).

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strconv"
	"strings"
	"testing"
)

// cardSourceFragment is >= 8 normalized runes, which is C25's contract floor
// (internal/risk/taintmatch.go); below it no contiguous fragment exists and the
// rule cannot fire, which test 3 uses.
const (
	taintFragment143   = "合同编号 HT-2026-0731-KX"
	taintOrigin143     = "https://files.example.com/q3-notes.txt"
	taintSourceFact143 = "web.fetch|" + taintOrigin143 + "|" + taintFragment143
)

// cardView143 re-declares the printed keys rather than decoding into
// panel.ApprovalCardView: the reading has to be about the bytes this command
// emits, which is what the panel consumes.
type cardView143 struct {
	CorrelationID          string   `json:"correlationId"`
	Tool                   string   `json:"tool"`
	Level                  string   `json:"level"`
	RulesHit               []string `json:"rulesHit"`
	Reason                 string   `json:"reason"`
	ReasonKnown            bool     `json:"reasonKnown"`
	SessionOverrideBlocked bool     `json:"sessionOverrideBlocked"`
	DecidedBy              string   `json:"decidedBy"`
}

func (c cardView143) hasRule(id string) bool {
	for _, r := range c.RulesHit {
		if r == id {
			return true
		}
	}
	return false
}

// runPanelAssets143 drives the dispatched leg and hands back what it wrote on
// stdout and stderr plus its exit code.
func runPanelAssets143(t *testing.T, args ...string) (string, string, int) {
	t.Helper()
	outFile, err := os.CreateTemp(t.TempDir(), "stdout-*.json")
	if err != nil {
		t.Fatalf("temp stdout: %v", err)
	}
	errFile, err := os.CreateTemp(t.TempDir(), "stderr-*.txt")
	if err != nil {
		outFile.Close()
		t.Fatalf("temp stderr: %v", err)
	}
	savedOut, savedErr := os.Stdout, os.Stderr
	os.Stdout, os.Stderr = outFile, errFile
	defer func() { os.Stdout, os.Stderr = savedOut, savedErr }()

	code := cmdPanelAssets(args)

	os.Stdout, os.Stderr = savedOut, savedErr
	if err := outFile.Close(); err != nil {
		t.Fatalf("close stdout: %v", err)
	}
	if err := errFile.Close(); err != nil {
		t.Fatalf("close stderr: %v", err)
	}
	stdout, err := os.ReadFile(outFile.Name())
	if err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	stderr, err := os.ReadFile(errFile.Name())
	if err != nil {
		t.Fatalf("read stderr: %v", err)
	}
	return string(stdout), string(stderr), code
}

func decodeCard143(t *testing.T, stdout string) cardView143 {
	t.Helper()
	var card cardView143
	if err := json.Unmarshal([]byte(stdout), &card); err != nil {
		t.Fatalf("the card path prints JSON on stdout, decoding it failed: %v\n--- stdout ---\n%s", err, stdout)
	}
	return card
}

// TestAC1AC2TaintSourceLegProducesAJudgedR4 is the hole this ticket exists to
// close: the flag declares a sensitive-source read, the outgoing call carries a
// fragment of it, and the taint clause on the card is produced by internal/risk.
func TestAC1AC2TaintSourceLegProducesAJudgedR4(t *testing.T) {
	stdout, stderr, code := runPanelAssets143(t,
		"-taint-source", taintSourceFact143,
		"-l2", "notify", "把 "+taintFragment143+" 发到远端")
	if code != 0 {
		t.Fatalf("exit code = %d, want 0 (stderr: %s)", code, stderr)
	}
	card := decodeCard143(t, stdout)

	if !card.hasRule("R4") {
		t.Fatalf("AC#1 RED: rulesHit = %v, want the taint rule among them.\n"+
			"The declared source is carried verbatim by the call, so this is the rule's own judgement being missing, not the input's.\nstdout:\n%s", card.RulesHit, stdout)
	}
	if card.Level != "L2" {
		t.Errorf("AC#1 RED: level = %q, want L2 (the taint rule's verdict is L2 and fusion takes the max)", card.Level)
	}
	if !card.SessionOverrideBlocked {
		t.Errorf("AC#1 RED: sessionOverrideBlocked = false; the taint rule's hit is the one a D45 session authorization may not cover (SPEC-06 §8.3), and only internal/risk sets this")
	}
	if !strings.Contains(card.Reason, taintOrigin143) {
		t.Errorf("AC#1 RED: reason = %q, which does not name the declared source %q.\n"+
			"The source string reaches the card through C25's provenance record and nowhere else: if this fails while R4 is present, something started writing its own explanation.", card.Reason, taintOrigin143)
	}
	if !strings.HasPrefix(card.Reason, "R1: ") || !strings.Contains(card.Reason, "R4") {
		t.Errorf("reason = %q, want both judged clauses (fusion names every rule that fired)", card.Reason)
	}
	if !card.ReasonKnown {
		t.Errorf("reasonKnown = false, want true: a card with a reason must say it has one")
	}
	if card.DecidedBy != "native" {
		t.Errorf("decidedBy = %q, want native (D33/F2)", card.DecidedBy)
	}
}

// TestAC2TaintFlagDoesNotLeakIntoTheDefaultCard is the in-process half of "the
// new capability did not seep into the default path"; the byte-for-byte half is
// the external cmp recorded in docs/evidence/s1/143-panel-assets-r4-leg-r1.md.
func TestAC2TaintFlagDoesNotLeakIntoTheDefaultCard(t *testing.T) {
	stdout, _, code := runPanelAssets143(t, "-l2", "notify", "把 "+taintFragment143+" 发到远端")
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	card := decodeCard143(t, stdout)
	if strings.Contains(stdout, "R4") {
		t.Fatalf("AC#2 RED: the card printed without -taint-source contains \"R4\" anyway:\n%s", stdout)
	}
	if len(card.RulesHit) != 1 || card.RulesHit[0] != "R1" {
		t.Errorf("AC#2 RED: rulesHit = %v, want exactly [R1] - the declared floor is all a bare assessor can judge", card.RulesHit)
	}
	if card.SessionOverrideBlocked {
		t.Errorf("AC#2 RED: sessionOverrideBlocked = true with no taint declared")
	}
}

// TestAC3TaintSourceThatTheCallDoesNotCarryIsNotAHit is the over-match guard.
// A wired detector must hit on the fragment and on nothing else; both shapes
// below fail for a detector that answers "hit" unconditionally (mutation
// AC#3(ii)), which is why this case is not a duplicate of the previous one.
func TestAC3TaintSourceThatTheCallDoesNotCarryIsNotAHit(t *testing.T) {
	t.Run("declared but absent from the outgoing call", func(t *testing.T) {
		stdout, _, code := runPanelAssets143(t,
			"-taint-source", taintSourceFact143,
			"-l2", "notify", "今天没有要外发的内容")
		if code != 0 {
			t.Fatalf("exit code = %d, want 0", code)
		}
		card := decodeCard143(t, stdout)
		if card.hasRule("R4") {
			t.Fatalf("RED: rulesHit = %v, want no taint rule - the call carries none of the declared source.\n%s", card.RulesHit, stdout)
		}
		if strings.Contains(stdout, taintOrigin143) {
			t.Errorf("RED: the declared source appears on a card that should not name it:\n%s", stdout)
		}
		if card.Level != "L1" {
			t.Errorf("level = %q, want L1 (the declared floor and nothing above it)", card.Level)
		}
	})

	t.Run("fragment below the contract floor cannot match", func(t *testing.T) {
		// 7 normalized characters: C25 documents that a source shorter than the
		// >=8-char floor produces no matchable fragment (provenance.go's
		// "deliberate limits"). A detector that hit here would be matching
		// something the contract does not name.
		stdout, _, code := runPanelAssets143(t,
			"-taint-source", "fs.read|D:/notes/short.txt|abc1234",
			"-l2", "notify", "abc1234")
		if code != 0 {
			t.Fatalf("exit code = %d, want 0", code)
		}
		card := decodeCard143(t, stdout)
		if card.hasRule("R4") {
			t.Fatalf("RED: rulesHit = %v for a 7-char source, below the 8-char contract floor:\n%s", card.RulesHit, stdout)
		}
	})
}

// TestAC2TaintSourceIsVisibleInTheUsageBlock: AC#2 asks that the flag be
// findable from -h and that its help say which side of the line it is on.
func TestAC2TaintSourceIsVisibleInTheUsageBlock(t *testing.T) {
	_, stderr, code := runPanelAssets143(t, "-h")
	if code != 2 {
		t.Fatalf("exit code for -h = %d, want 2 (flag.ContinueOnError's help path)", code)
	}
	if !strings.Contains(stderr, "-taint-source") {
		t.Fatalf("AC#2 RED: -h never mentions -taint-source, so an operator cannot find it:\n%s", stderr)
	}
	if !strings.Contains(stderr, "INPUT FACT") {
		t.Errorf("AC#2 RED: the help does not say the flag is an input fact: %s", stderr)
	}
	if !strings.Contains(stderr, "no verdict") {
		t.Errorf("AC#2 RED: the help does not disclaim that the flag carries a conclusion: %s", stderr)
	}
}

// TestAC1MalformedTaintSourceIsRefused keeps the flag honest about its shape:
// a half-declared source is a missing fact, and missing input is an error here
// rather than a card that silently omits a rule.
func TestAC1MalformedTaintSourceIsRefused(t *testing.T) {
	for _, tc := range []struct{ name, value string }{
		{"no separators", "web.fetch"},
		{"two parts", "web.fetch|origin only"},
		{"empty tool", "|origin|content"},
		{"empty content", "web.fetch|origin|   "},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, code := runPanelAssets143(t, "-taint-source", tc.value, "-l2", "notify", "hello")
			if code != 2 {
				t.Errorf("exit code = %d, want 2 for %q\nstdout:\n%s\nstderr:\n%s", code, tc.value, stdout, stderr)
			}
			if !strings.Contains(stderr, "-taint-source") {
				t.Errorf("stderr does not name the flag it refused: %s", stderr)
			}
			if strings.TrimSpace(stdout) != "" {
				t.Errorf("a refused flag must print no card, got:\n%s", stdout)
			}
		})
	}
}

// TestAC1CmdSideEmitsNoVerdictTokens is the P9 red line as an instrument: the
// strings this package hands the card path may name sources, never verdicts.
// Comments are excluded on purpose - prose about the rules is not a rule
// output - but every string literal is read, including the flag descriptions.
func TestAC1CmdSideEmitsNoVerdictTokens(t *testing.T) {
	banned := []string{"R4", "R9", "包含来自", "rulesHit", "sessionOverrideBlocked", "L2:"}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "panel_assets.go", nil, 0)
	if err != nil {
		t.Fatalf("parse panel_assets.go: %v", err)
	}
	var hits []string
	ast.Inspect(file, func(n ast.Node) bool {
		lit, ok := n.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}
		value, err := strconv.Unquote(lit.Value)
		if err != nil {
			t.Errorf("unquote %s: %v", lit.Value, err)
			return true
		}
		for _, b := range banned {
			if strings.Contains(value, b) {
				hits = append(hits, fset.Position(lit.Pos()).String()+" holds "+b)
			}
		}
		return true
	})
	if len(hits) > 0 {
		t.Fatalf("AC#1 RED: panel_assets.go emits verdict-shaped text out of its own strings, so a card could be printed by this file instead of by internal/risk: %s",
			strings.Join(hits, ", "))
	}
}
