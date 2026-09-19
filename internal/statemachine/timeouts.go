package statemachine

import "time"

// Phase is the Listening timeout context. The D43 row column "超时（首轮
// 15s / 会话内 90s / Conversation 30s）" is only decidable once we know HOW
// Listening was entered, so the entering table row pins the phase and the
// machine keeps it while Listening lasts.
type Phase uint8

const (
	PhaseNone         Phase = iota // not in Listening (or phaseless entry)
	PhaseFirstRound                // cold session start (#4/#7/#33): 15s
	PhaseInSession                 // warm resume (#29) or interrupt (#27): 90s
	PhaseConversation              // conversation mode (#28/#30): 30s
)

// Per-state timeout table (SPEC-08 §3 timeout column):
//
//	Listening        first-round 15s | in-session 90s | conversation 30s (#12)
//	Confirming       L1 countdown 3s (#21; SPEC-06 L1 window is 2-3s, the
//	                 upper bound is chosen - it is the user-safe side)
//	Warm             90s idle (#31)
//	Settling         3s (#32)
//	AwaitingApproval 300s (#24; deny on expiry, the 30s pre-warning is UI)
//	Error            10s auto-ack (#38)
//
// Every other state (incl. Sleeping - the zero-timer discipline, and
// Armed/Muted/Thinking/Acting/Speaking which are event-driven by the
// pipeline) has NO machine-side timeout. Timeouts are monotonic (time.Timer
// carries a monotonic reading; observe clock discipline D42#9).
type TimeoutSpec struct {
	After time.Duration
	Event Event // the event Dispatch receives when the timer fires
}

type timeoutKey struct {
	state State
	phase Phase
}

// DefaultTimeouts is the frozen timeout table. Tests may override via Options.
func DefaultTimeouts() map[timeoutKey]TimeoutSpec {
	return map[timeoutKey]TimeoutSpec{
		{state: StateListening, phase: PhaseFirstRound}:   {15 * time.Second, EvTimeoutFirstRound},
		{state: StateListening, phase: PhaseInSession}:    {90 * time.Second, EvTimeoutInSession},
		{state: StateListening, phase: PhaseConversation}: {30 * time.Second, EvTimeoutConversation},
		{state: StateConfirming}:                          {3 * time.Second, EvConfirmExpired},
		{state: StateWarm}:                                {90 * time.Second, EvWarmIdle},
		{state: StateSettling}:                            {3 * time.Second, EvSettleExpired},
		{state: StateAwaitingApproval}:                    {300 * time.Second, EvApprovalTimeout},
		{state: StateError}:                               {10 * time.Second, EvErrorAck},
	}
}
