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

	// ---- ticket 37 AC#5's aggregate-card fields (readface wave 1, 2026-11-17) ----
	//
	// One tool call that touches many targets renders as ONE card with the
	// batch's display fields on it (D45-1's card shape: total + affected dirs +
	// the first preview rows + the deferred remainder + the summary sentence).
	// They are DISPLAY data read off approval.CardSummary (internal/agent/
	// approval's cardsummary_37.go) by the host's verdict reader; this type
	// carries them and derives nothing - a card that re-derived a count would
	// be a second judge of the batch.
	//
	// Flat fields, not a nested struct, on purpose: the frontend contract ruler
	// (TestApprovalCardViewJSONKeysMatchFrontendTypes) needs no new page-side
	// interface for a shape the page half has not declared yet. Until the
	// frontend leg declares the five keys, the five names ride the judged-gap
	// ledger in composer_test.go - the same registration procedure ticket 415
	// built - so the gap is named every run and closes the moment the page
	// declares them.
	//
	// omitempty keeps a normal card's JSON byte-identical to its pre-field
	// shape: zero/empty batch fields never reach the wire (a card that is not
	// about a batch must not wear batch zeros).
	BatchTotal    int      `json:"batchTotal,omitempty"`
	BatchDirs     []string `json:"batchDirs,omitempty"`
	BatchFirst    []string `json:"batchFirst,omitempty"`
	BatchDeferred int      `json:"batchDeferred,omitempty"`
	BatchSummary  string   `json:"batchSummary,omitempty"`

	// ---- ticket 167 AC#3's two queue numbers ----
	//
	// Seq is the queue's own ADMISSION number (the queue's qitem.Seq, read off
	// LiveApproval.Seq by the host's verdict reader): stamped once when the card
	// was queued, never rewritten - cancelling an earlier card does not renumber
	// the ones behind it, which is the 取消不重排 rule the ticket pins. Position
	// is the card's 1-based place in the FIFO RIGHT NOW (the depth badge), and it
	// moves as the queue drains - that is the honest "当前第几条", deliberately a
	// separate key from the stable one so the page can show both without either
	// lying for the other.
	//
	// Both are omitempty for the reason Q-51 already measured: an int without it
	// always emits 0, and the CLI card path (`wisp panel-assets`) renders cards
	// the queue never numbered - a "seq":0/"position":0 there would be the
	// degenerate shape the ticket forbids. Zero = unread, and the key stays home.
	Seq      uint64 `json:"seq,omitempty"`
	Position int    `json:"position,omitempty"`
}

// HasBatch reports whether this card carries the aggregate-card fields. A
// card rendered from a plain decision never does (the decision had too few
// targets or was tainted), and the page must not render batch chrome from
// zeros.
func (v ApprovalCardView) HasBatch() bool {
	return v.BatchTotal > 0
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
