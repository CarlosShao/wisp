package tools

import (
	"context"
	"encoding/json"
	"time"

	"github.com/CarlosShao/wisp/internal/risk"
)

// Decision is one fully resolved bridge verdict: what C19 concluded, why, and
// what the caller must show the user. The confirmation card (tickets 21/37)
// renders RulesHit and Reason VERBATIM (ticket 17's frozen-contract note), so
// this struct, not a string, is the unit that leaves the bridge.
type Decision struct {
	Tool         string
	Provider     Kind
	Params       map[string]any
	Args         json.RawMessage
	Level        risk.Level
	RulesHit     []risk.RuleID
	Reason       string
	Paths        []string // canonical, as judged by R2/R3
	Capabilities []Capability
	// SessionOverrideBlocked is R4's "no D45 grant covers this" flag.
	SessionOverrideBlocked bool
	// Mode is the user-facing permission mode this call was screened under
	// (ticket 90, R20). It travels with the verdict because "why was I not
	// asked?" is a question the audit trail has to answer from the record, not
	// from the config file's current contents.
	Mode risk.Mode
	// ModeSilenced is true when the mode removed a question this verdict would
	// otherwise have asked. It is NOT an allow: the assessed Level below stays
	// the judge's own conclusion, which is what the tool_call.risk_level column
	// books (SPEC-02 §3).
	ModeSilenced bool
	// ModeKept names the red line that refused to silence the verdict, empty
	// when the mode applied as asked. Audit text, rendered verbatim.
	ModeKept string
	// Blacklist is risk.Gate's reading of this call's paths (tier A/B plus
	// which B-tier files a single L2 confirmation could unlock). See
	// BlacklistNote in mode.go: it is evidence for the card, never authority.
	Blacklist BlacklistNote
	// DecisionColumn is the tool_call.decision value the routing resolved to
	// (agent.DecisionAllow / Reject / Timeout). Empty until routing runs.
	DecisionColumn string
	// CorrelationID routes the decision to the right approval card (C18);
	// the same value lands in the tool_call row.
	CorrelationID string
	TaskID        string
	CallID        string
	Timeout       time.Duration
}

// LevelString renders the level in the shape the tool_call.risk_level column
// takes ('L0'|'L1'|'L2'). risk.Deny has no column value - it is booked as L2
// with decision='reject', which is the only honest mapping the frozen SPEC-02
// DDL allows.
func (d Decision) LevelString() string {
	switch d.Level {
	case risk.L1:
		return "L1"
	case risk.L2, risk.Deny:
		return "L2"
	default:
		return "L0"
	}
}

// Answer is what a gate callback concludes. The vocabulary is deliberately
// narrower than "error": a gate that could not be reached is a distinct,
// fail-closed event from a user who said no.
type Answer string

// The gate answers.
const (
	// AnswerAllow: execute.
	AnswerAllow Answer = "allow"
	// AnswerReject: user/queue said no.
	AnswerReject Answer = "reject"
	// AnswerVeto: the L1 pre-execution window was vetoed (B1). Distinguished
	// from Reject because the forensics and the user-facing wording differ.
	AnswerVeto Answer = "veto"
	// AnswerTimeout: the L1 window ran out unopposed, which per SPEC-06 §2
	// MEANS EXECUTE; for the L2 queue it means auto-reject (ticket 21:
	// "timeout 300s -> auto-REJECT, never infinite wait"). The branch is the
	// caller's, never the gate's.
	AnswerTimeout Answer = "timeout"
)

// Gate is the ticket-21 seam: the two callbacks the L0/L1/L2 routing needs.
// The bridge owns the ROUTING (which level goes to which callback, what each
// answer means, what gets booked); ticket 21 owns the MECHANICS behind these
// two methods (countdown strip, veto channels, native card, queue).
//
// Implementations must return without blocking past the call's own deadline
// and must treat a missing/unreachable UI as AnswerReject: the fail-closed
// direction is "did not run", never "assumed approved".
type Gate interface {
	// PendingWindow is the L1 route: 2-3s pre-execution block window with
	// veto channels. Timeout -> execute.
	PendingWindow(ctx context.Context, d Decision) (Answer, string)
	// PendingApproval is the L2 route: the C18 approval queue, answered only
	// from native sources (F2).
	PendingApproval(ctx context.Context, d Decision) (Answer, string)
}

// GateFuncs adapts two functions to Gate; the test harness and any host that
// only has callbacks use it.
type GateFuncs struct {
	Window   func(context.Context, Decision) (Answer, string)
	Approval func(context.Context, Decision) (Answer, string)
}

// PendingWindow implements Gate. A nil callback fails closed.
func (g GateFuncs) PendingWindow(ctx context.Context, d Decision) (Answer, string) {
	if g.Window == nil {
		return AnswerReject, "L1 确认窗口不可用（票 21 未接入），已拒绝执行"
	}
	return g.Window(ctx, d)
}

// PendingApproval implements Gate. A nil callback fails closed.
func (g GateFuncs) PendingApproval(ctx context.Context, d Decision) (Answer, string) {
	if g.Approval == nil {
		return AnswerReject, "L2 审批通道不可用（票 21 未接入），已拒绝执行"
	}
	return g.Approval(ctx, d)
}

// NoGate is the default: nothing is wired to ask the user, so anything above
// L0 is refused. This is the same fail-closed posture the loop already runs
// with its PassThroughUnclassifiedRisk flag off (agent/tools.go:20-24), now
// expressed where the decision is actually made.
type NoGate struct{}

// PendingWindow implements Gate.
func (NoGate) PendingWindow(context.Context, Decision) (Answer, string) {
	return AnswerReject, "L1 确认窗口尚未接入（票 21），已拒绝执行"
}

// PendingApproval implements Gate.
func (NoGate) PendingApproval(context.Context, Decision) (Answer, string) {
	return AnswerReject, "L2 审批通道尚未接入（票 21），已拒绝执行"
}
