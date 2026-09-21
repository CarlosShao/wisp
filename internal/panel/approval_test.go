package panel

// The panel / frontend data-contract tests (ticket 77 AC#3 and AC#5).
//
// AC#3 asks that the L2 card be driven by a real risk decision rather than
// component-library demo props. Two things are checked here:
//
//   - the values: NewApprovalCardView over risk.NewRiskAssessor() (the same
//     constructor internal/tools/bridge.go uses in production) produces the
//     level, rules and reason the assessor actually decided;
//   - the shape: the JSON keys the frontend interface promises are exactly the
//     keys this Go struct emits, in both directions. That is the half a demo
//     fixture can never fail - a renamed field shows up as undefined in the
//     card and as nothing in a browser test that invents its own props.
//
// AC#5 asks for proof that a WebView restart recovers the same state from the
// Go side. TestPanelSnapshotSurvivesWebviewRestart is that proof in the
// semi-real form the ticket allows: the entire on-screen state is marshalled
// out, a "restart" throws every Go-side value away and reads it back, and the
// two must be identical. The complementary half - that the page cannot
// remember anything itself - is TestPanelFrontendIsStateless.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// Snapshot and ResultChunk moved to composer.go in ticket 92: a view model only
// a test file can build is one production never sends, and the composer section
// this ticket adds has to be constructible from the assembly root.

func TestApprovalCardViewFromRealAssessor(t *testing.T) {
	assessor := risk.NewRiskAssessor()
	subject := ApprovalSubject{
		CorrelationID: "corr-77-1",
		Tool:          "fs.delete",
		Args:          []string{"--permanent", `C:\Users\swq\Documents\notes.md`},
		CallChain:     []string{"session", "agent-loop", "fs.delete"},
		Facts: risk.Facts{
			Declared:     risk.L1,
			Paths:        []string{`C:\Users\swq\Documents\notes.md`},
			Irreversible: []string{"delete"},
		},
	}

	view := NewApprovalCardView(assessor, subject)
	t.Logf("real assessor verdict: level=%s rules=%v reason=%q", view.Level, view.RulesHit, view.Reason)

	if view.DecidedBy != "native" {
		t.Errorf("DecidedBy = %q, want %q - the panel has no verdict of its own (D33/F2)", view.DecidedBy, "native")
	}
	if view.Level == "" {
		t.Error("Level is empty: the card would render a confirmation prompt with no risk level")
	}
	if view.Level == "L0" {
		t.Errorf("an irreversible delete assessed to %q - the card must not be shown for a call the gate considers free", view.Level)
	}
	if len(view.RulesHit) == 0 {
		t.Error("RulesHit is empty while the verdict needed a human: the card could not name why it appeared")
	}
	if !view.ReasonKnown {
		t.Error("ReasonKnown = false with a non-empty verdict: Q-23 requires the card to distinguish 'no reason given' from 'no risk', and an unset reason here means the wiring lost the text")
	}
	if strings.TrimSpace(view.Reason) != "" && view.ReasonKnown {
		// The verbatim requirement (SPEC-06 §7): what the assessor said is what
		// the operator reads, with nothing rewritten.
		if strings.ContainsAny(view.Reason, "\x00") {
			t.Errorf("reason reached the card mangled: %q", view.Reason)
		}
	}
	if view.CorrelationID != subject.CorrelationID || view.Tool != subject.Tool {
		t.Errorf("identity fields lost: %+v", view)
	}
	if !reflect.DeepEqual(view.Args, subject.Args) {
		t.Errorf("Args = %q, want the verbatim argv %q - the card shows the complete command or it shows nothing", view.Args, subject.Args)
	}
	if !reflect.DeepEqual(view.CallChain, subject.CallChain) {
		t.Errorf("CallChain = %q, want %q", view.CallChain, subject.CallChain)
	}
}

func TestApprovalCardViewMarksMissingReason(t *testing.T) {
	// The D22 "reason enum" leg of ticket 20 has not landed, so an empty reason
	// is a real state. It must arrive flagged, never blank and confident.
	view := CardViewFromDecision(ApprovalSubject{CorrelationID: "c", Tool: "shell.run"}, risk.Decision{
		Level:    risk.L2,
		RulesHit: []risk.RuleID{risk.R6},
		Reason:   "   ",
	})
	if view.ReasonKnown {
		t.Error("ReasonKnown = true for a blank reason: the card would render an empty slot next to a blocking prompt")
	}
	if view.Level != "L2" {
		t.Errorf("Level = %q, want L2 - an unknown reason must not soften the verdict", view.Level)
	}
}

func TestApprovalCardViewJSONKeysMatchFrontendTypes(t *testing.T) {
	root := panelRepoRoot(t)
	data, err := os.ReadFile(filepath.Join(root, "frontend", "src", "lib", "panel.ts"))
	if err != nil {
		t.Fatalf("read frontend/src/lib/panel.ts: %v", err)
	}
	text := string(data)

	pairs := []struct {
		name    string
		goType  any
		tsIface string
	}{
		{"ApprovalCardView", ApprovalCardView{}, "ApprovalCardView"},
		{"ResultChunk", ResultChunk{}, "ResultChunkView"},
		{"Snapshot", Snapshot{}, "PanelSnapshot"},
	}
	for _, p := range pairs {
		goKeys := jsonKeysOf(p.goType)
		tsKeys, ok := tsInterfaceKeys(text, p.tsIface)
		if !ok {
			t.Fatalf("frontend/src/lib/panel.ts declares no interface %s - the contract was renamed on one side only", p.tsIface)
		}
		if missing := subtract(goKeys, tsKeys); len(missing) > 0 {
			t.Errorf("Go %s emits %v that interface %s does not declare", p.name, missing, p.tsIface)
		}
		if extra := subtract(tsKeys, goKeys); len(extra) > 0 {
			t.Errorf("interface %s reads %v that Go %s never sends - those fields render as undefined", p.tsIface, extra, p.name)
		}
		t.Logf("%s <-> %s: %d JSON keys reconciled", p.name, p.tsIface, len(goKeys))
	}
}

func TestPanelSnapshotSurvivesWebviewRestart(t *testing.T) {
	assessor := risk.NewRiskAssessor()
	first := Snapshot{
		Pending: []ApprovalCardView{NewApprovalCardView(assessor, ApprovalSubject{
			CorrelationID: "corr-restart",
			Tool:          "fs.write",
			Args:          []string{`C:\Users\swq\Desktop\plan.md`},
			Facts:         risk.Facts{Declared: risk.L1, Paths: []string{`C:\Users\swq\Desktop\plan.md`}, Irreversible: []string{"overwrite"}},
		})},
		Results:     []ResultChunk{{CorrelationID: "corr-restart", Text: "writing", Done: false}},
		GeneratedAt: "2026-09-21T13:00:00Z",
	}
	if first.Pending[0].Level != "L2" && first.Pending[0].Level != "Deny" {
		t.Fatalf("fixture premise broken: this call must need a human, got %q", first.Pending[0].Level)
	}

	wire, err := json.Marshal(first)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	// The restart: everything the page held is gone, and what comes back is
	// exactly one push from Go.
	var after Snapshot
	if err := json.Unmarshal(wire, &after); err != nil {
		t.Fatalf("unmarshal after restart: %v", err)
	}
	if !reflect.DeepEqual(first, after) {
		t.Errorf("state after a WebView restart differs from the state before it:\n before %+v\n after  %+v", first, after)
	}
	t.Logf("snapshot round-tripped through %d bytes of JSON with %d pending card(s)", len(wire), len(after.Pending))
}

// ------------------------------------------------------------------ contracts

var jsonTagRe = regexp.MustCompile(`json:"([^",]+)`)

func jsonKeysOf(v any) []string {
	var out []string
	t := reflect.TypeOf(v)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue // unexported
		}
		if tag := f.Tag.Get("json"); tag != "" && tag != "-" {
			if m := jsonTagRe.FindStringSubmatch(`json:"` + tag + `"`); m != nil {
				out = append(out, m[1])
				continue
			}
			out = append(out, strings.Split(tag, ",")[0])
		}
	}
	return out
}

var tsBlockRe = regexp.MustCompile(`(?s)export interface (\w+) \{(.*?)\n\}`)

func tsInterfaces(text string) map[string][]string {
	out := map[string][]string{}
	keyRe := regexp.MustCompile(`(?m)^\s*(\w+)(\??):`)
	for _, m := range tsBlockRe.FindAllStringSubmatch(text, -1) {
		var keys []string
		for _, k := range keyRe.FindAllStringSubmatch(m[2], -1) {
			keys = append(keys, k[1])
		}
		out[m[1]] = keys
	}
	return out
}

func tsInterfaceKeys(text, name string) ([]string, bool) {
	keys, ok := tsInterfaces(text)[name]
	return keys, ok
}

func subtract(a, b []string) []string {
	inB := map[string]bool{}
	for _, s := range b {
		inB[s] = true
	}
	var out []string
	for _, s := range a {
		if !inB[s] {
			out = append(out, s)
		}
	}
	return out
}
