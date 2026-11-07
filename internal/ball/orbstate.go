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

// The jarvis page (ticket 65 second visual style, the cyber1443 component)
// ships FOUR states of its own - idle|thinking|success|alert (its
// STATE_TARGETS keys) - so the same 20 frozen ball states map onto the richer
// native vocabulary by the SAME work-shape partition, with two extensions the
// page's own states buy:
//
//	"thinking" - the identical orbActiveStates set as the siri mapping (the
//	             work family; one partition, two vocabularies).
//	"alert"    - the failure/attention family (Error, NoNetwork,
//	             WatchdogAlert, Stuck, Unconfigured): the page's own warning
//	             mood, which the siri page does not have. The D2D fallback
//	             still renders its dedicated failure recipes when the orb is
//	             not attached - nothing downstream changes.
//	"success"  - Settling, the one-shot post-work rest transition: the
//	             page's own "job done" mood.
//	"idle"     - everything else (standby, ambient, Downloading) - unchanged
//	             from the siri partition.
//
// jarvisStateFor's fallback is "idle" (the entry IGNORES unknown names, but
// the mapping itself never produces one - pinned by tests over all 20 states).
func jarvisStateFor(s statemachine.State) string {
	if orbActiveStates[s] {
		return "thinking"
	}
	switch s {
	case statemachine.StateError,
		statemachine.StateNoNetwork,
		statemachine.StateWatchdogAlert,
		statemachine.StateStuck,
		statemachine.StateUnconfigured:
		return "alert"
	case statemachine.StateSettling:
		return "success"
	}
	return "idle"
}

// orbAudioFor reports whether the orb should be audio-driven in s.
func orbAudioFor(s statemachine.State) bool { return orbAudioStates[s] }
