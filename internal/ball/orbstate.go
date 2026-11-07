package ball

// The Wisp-to-orb state mapping (ticket 65 integration leg, A1594). The
// upstream orb ships exactly TWO states (orb749/src/orb-states.ts:
// orbStateNames = ["idle", "thinking"]) and the owner's order for this leg is
// to ship the upstream page unmodified, NOT to author new shader states. So
// the 20 frozen ball states (D43, C12) map onto the two upstream names by
// WORK SHAPE:
//
//	"thinking" - the machine is visibly working on the user's behalf: voice
//	             in/out, reasoning, acting, awaiting a human decision, queued
//	             work. The orb's activation profile (brighter, faster bands)
//	             is what the page itself calls its thinking state.
//	"idle"     - everything else: asleep, armed/muted standby, ambient, and
//	             the failure family (which the D2D fallback renders with its
//	             own dedicated recipes when the orb is not attached).
//
// Audio-reactive states push setAudioBands; in every other state the bands
// stay zeroed (the page's own rule: "zero them on stop").

import "github.com/CarlosShao/wisp/internal/statemachine"

// orbActiveStates are the ball states that map to the orb's "thinking".
var orbActiveStates = map[statemachine.State]bool{
	statemachine.StateListening:        true,
	statemachine.StateThinking:         true,
	statemachine.StateActing:           true,
	statemachine.StateSpeaking:         true,
	statemachine.StateConversation:     true,
	statemachine.StateConfirming:       true,
	statemachine.StateAwaitingApproval: true,
	statemachine.StateQueued:           true,
}

// orbAudioStates are the ball states whose audio drives the orb's bands.
var orbAudioStates = map[statemachine.State]bool{
	statemachine.StateListening:    true,
	statemachine.StateSpeaking:     true,
	statemachine.StateConversation: true,
}

// orbStateFor maps one frozen ball state onto the orb page's state name.
// Every valid D43 state lands on exactly one of the two upstream names; the
// fallback for anything unforeseen is "idle" (never a page error - the page
// throws on an unknown name).
func orbStateFor(s statemachine.State) string {
	if orbActiveStates[s] {
		return "thinking"
	}
	return "idle"
}

// orbAudioFor reports whether the orb should be audio-driven in s.
func orbAudioFor(s statemachine.State) bool { return orbAudioStates[s] }
