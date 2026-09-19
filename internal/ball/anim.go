package ball

import "github.com/CarlosShao/wisp/internal/statemachine"

// Animation discipline (SPEC-08 §2, ticket 07):
//
//   - Sleeping: fully static, ZERO timers, CPU ~= 0. Enforced by
//     AnimationPolicy returning no timer and asserted in tests.
//   - Loop animations run only in the five active states
//     Listening/Thinking/Acting/Speaking/Confirming (Acting's visual is
//     static today, so its policy is "allowed but not needed" - no timer).
//   - Warm: 2.4s opacity breathing 0.55<->0.7 via the compositor-friendly
//     path - a slow (10fps) tick that changes ONLY the layered-window
//     constant alpha (no D2D re-render, no pixel upload per frame).
//   - Settling: the SPEC-08 §2.1 260ms fade 1 -> 0.35 as a bounded one-shot
//     burst at the 30fps cap, then static. Documented SPEC-sanctioned
//     exception to the "five states" list (it is not a loop).
//   - Every period respects the 30fps cap (>= 33ms).
//
// WarmKind lets the window route a tick to the cheap alpha path.

// AnimKind classifies what a timer tick triggers.
type AnimKind uint8

const (
	AnimNone  AnimKind = iota
	AnimFrame          // full D2D re-render at animation phase t
	AnimAlpha          // constant-alpha only (Warm breathing)
	AnimFade           // one-shot Settling fade progress (full re-render)
)

// AnimPolicy is the timer policy of one state.
type AnimPolicy struct {
	Kind     AnimKind
	PeriodMs int // 0 = no timer
}

// AnimationPolicy returns the timer policy for a state. This function is the
// single authority; the window and the tests both read it.
func AnimationPolicy(s statemachine.State) AnimPolicy {
	switch s {
	case statemachine.StateListening:
		return AnimPolicy{AnimFrame, framePeriod(ListeningWaveMs)}
	case statemachine.StateThinking:
		return AnimPolicy{AnimFrame, framePeriod(ThinkingSweepMs)}
	case statemachine.StateConfirming:
		return AnimPolicy{AnimFrame, framePeriod(ConfirmingPulseMs)}
	case statemachine.StateSpeaking:
		return AnimPolicy{AnimFrame, framePeriod(SpeakingBreathMs)}
	case statemachine.StateActing:
		// Static solid core today; the 5-state whitelist ALLOWS a timer, the
		// visual does not need one. Fewer timers = better SLO.
		return AnimPolicy{AnimNone, 0}
	case statemachine.StateWarm:
		// 2.4s breathing, compositor-friendly: constant-alpha path at 10fps.
		return AnimPolicy{AnimAlpha, 1000 / WarmBreathFPS}
	case statemachine.StateSettling:
		// One-shot 260ms fade at the 30fps cap; the window kills the timer
		// once fadeProgress reaches 1.
		return AnimPolicy{AnimFade, framePeriod(SettlingFadeMs)}
	default:
		// Sleeping, Armed, Muted, AwaitingApproval, Conversation,
		// Downloading, Error, NoNetwork, Unconfigured, FirstRun,
		// WatchdogAlert, Queued, Stuck: static.
		// (FirstRun's guide pulse and AwaitingApproval's one-shot pulse land
		// with their wiring tickets; S1 keeps them static per the whitelist.)
		return AnimPolicy{AnimNone, 0}
	}
}

// framePeriod derives a timer period from an animation period so the phase
// advances smoothly but never exceeds the 30fps cap.
func framePeriod(periodMs int) int {
	const steps = 24 // steps per cycle (period/24 ticks)
	p := periodMs / steps
	if p < MinFrameMs {
		p = MinFrameMs
	}
	return p
}
