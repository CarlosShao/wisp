package panel

// The L2 confirmation card's render model (ticket 77 AC#3).
//
// This is the only place the panel learns anything about risk, and it learns it
// by ASKING internal/risk - the same assessor the tool bridge runs (see
// internal/tools/bridge.go, which builds risk.NewRiskAssessor()). Nothing here
// re-implements a rule, re-derives a level, or trusts a level handed in from
// outside: SPEC-06 §1 keeps the verdict native, and Q-17/Q-23 keep the card
// honest by showing what the verdict was, which rules produced it, the whole
// command, and the chain that got there.
//
// The field set of ApprovalCardView is a contract with the frontend:
// TestApprovalCardViewJSONKeysMatchFrontendTypes compares these JSON keys with
// the ApprovalCardView interface in frontend/src/lib/panel.ts in both
// directions, so a rename on one side goes red rather than rendering "undefined".

import (
	"strings"

	"github.com/CarlosShao/wisp/internal/risk"
)

// ApprovalSubject is the tool call the card is about. The agent loop owns the
// correlationId (C17); the caller supplies the same Facts it gave the assessor
// so the card can show what was actually judged.
type ApprovalSubject struct {
	CorrelationID string
	Tool          string
	Args          []string
	// CallChain is outermost-first, e.g. ["session", "agent-loop", "shell.run"].
	CallChain []string
	// Facts is the judged input. It is kept out of the JSON: the card shows the
	// verdict and the command, not the internal evidence bundle.
	Facts risk.Facts
}

// ApprovalCardView is what the panel renders for one confirmation.
type ApprovalCardView struct {
	CorrelationID string   `json:"correlationId"`
	Tool          string   `json:"tool"`
	Args          []string `json:"args"`
	// Level is risk.Level.String(): "L0" | "L1" | "L2" | "Deny".
	Level string `json:"level"`
	// RulesHit lists the frozen C19 rule ids that fired, in the assessor's order.
	RulesHit []string `json:"rulesHit"`
	// Reason is Decision.Reason verbatim.
	Reason string `json:"reason"`
	// ReasonKnown is false when the assessor produced no reason text. The card
	// must then say it lacks information - Q-23: an absent reason is never
	// allowed to render as "no risk".
	ReasonKnown bool `json:"reasonKnown"`
	// SessionOverrideBlocked mirrors Decision.SessionOverrideBlocked (R4 taint).
	SessionOverrideBlocked bool     `json:"sessionOverrideBlocked"`
	CallChain              []string `json:"callChain"`
	// DecidedBy states where the verdict came from. It is always "native" while
	// the risk gate is native (D33/F2), and the panel has no way to set it.
	DecidedBy string `json:"decidedBy"`
}

// NewApprovalCardView runs the real assessor over the subject and renders its
// verdict. An assessor is required: a nil one would silently produce an empty
// card, which is the failure shape this whole file exists to prevent.
func NewApprovalCardView(assessor *risk.RiskAssessor, subject ApprovalSubject) ApprovalCardView {
	decision := assessor.Assess(subject.Tool, paramsFromArgs(subject.Args), subject.Facts)
	return CardViewFromDecision(subject, decision)
}

// CardViewFromDecision is the pure half of NewApprovalCardView, split out so a
// test can feed it a decision produced anywhere (including a fixture replayed
// from SQLite history) without reconstructing the assessor.
func CardViewFromDecision(subject ApprovalSubject, decision risk.Decision) ApprovalCardView {
	rules := make([]string, 0, len(decision.RulesHit))
	for _, r := range decision.RulesHit {
		rules = append(rules, string(r))
	}
	view := ApprovalCardView{
		CorrelationID:          subject.CorrelationID,
		Tool:                   subject.Tool,
		Args:                   append([]string(nil), subject.Args...),
		Level:                  decision.Level.String(),
		RulesHit:               rules,
		Reason:                 decision.Reason,
		ReasonKnown:            strings.TrimSpace(decision.Reason) != "",
		SessionOverrideBlocked: decision.SessionOverrideBlocked,
		CallChain:              append([]string(nil), subject.CallChain...),
		DecidedBy:              "native",
	}
	if view.Args == nil {
		view.Args = []string{}
	}
	if view.RulesHit == nil {
		view.RulesHit = []string{}
	}
	if view.CallChain == nil {
		view.CallChain = []string{}
	}
	return view
}

// paramsFromArgs re-presents an argv vector as the map the rules read. The tool
// bridge builds the real map from the manifest (internal/tools/bridge.go); this
// covers the panel-side replay path, where only the rendered command is known.
// The mapping mirrors what the bridge feeds for the fields the rules consult,
// and nothing else: a missing key means the rule for it does not fire, which is
// the documented behaviour of risk.Facts ("all fields are optional").
func paramsFromArgs(args []string) map[string]any {
	params := map[string]any{}
	if len(args) > 0 {
		params["command"] = strings.Join(args, " ")
		params["argv"] = args
	}
	return params
}
