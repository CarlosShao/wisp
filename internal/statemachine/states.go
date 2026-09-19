package statemachine

// State is one of the 20 BallState values (D43; C12 frozen). The type is
// vocabulary only in ticket 03: the transition table, guards and per-state
// timeouts are implemented by ticket 07.
type State string

// The 20 states, names exactly as in D43 / SPEC-08. Order follows the D43
// table's narrative (boot -> idle -> session -> failure/terminal helpers).
const (
	StateFirstRun         State = "FirstRun"
	StateSleeping         State = "Sleeping"
	StateArmed            State = "Armed"
	StateMuted            State = "Muted"
	StateListening        State = "Listening"
	StateThinking         State = "Thinking"
	StateActing           State = "Acting"
	StateSpeaking         State = "Speaking"
	StateWarm             State = "Warm"
	StateConversation     State = "Conversation"
	StateConfirming       State = "Confirming"       // L1 quick confirm
	StateAwaitingApproval State = "AwaitingApproval" // L2 queue
	StateSettling         State = "Settling"
	StateDownloading      State = "Downloading"
	StateUnconfigured     State = "Unconfigured"
	StateNoNetwork        State = "NoNetwork"
	StateError            State = "Error"
	StateQueued           State = "Queued"
	StateStuck            State = "Stuck"
	StateWatchdogAlert    State = "WatchdogAlert"
)

// Idle is the ticket-03 placeholder state required by the ticket 03 contract
// ("no state machine transitions beyond an Idle placeholder"). It aliases
// Sleeping, the documented idle state, so nothing new is invented here.
const Idle = StateSleeping

// Valid reports whether s is one of the 20 frozen states.
func Valid(s State) bool {
	switch s {
	case StateFirstRun, StateSleeping, StateArmed, StateMuted,
		StateListening, StateThinking, StateActing, StateSpeaking,
		StateWarm, StateConversation, StateConfirming, StateAwaitingApproval,
		StateSettling, StateDownloading, StateUnconfigured, StateNoNetwork,
		StateError, StateQueued, StateStuck, StateWatchdogAlert:
		return true
	}
	return false
}
