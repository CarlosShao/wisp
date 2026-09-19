package statemachine

// Event is one input edge of the D43 transition table (SPEC-08 §3). Names are
// stable vocabulary: the ball wiring, hotkeys, speech pipeline and agent core
// all dispatch these; adding an event without a table row is a compile-time
// no-op and a rejected dispatch (undefined = stop, D22 gate 3).
type Event string

// The D43 events. The comment on each constant lists the D43 row(s) it drives.
const (
	// Boot / onboarding.
	EvOnboardingCompleted Event = "onboarding.completed" // #1 (model+dir auth+key present)
	EvModelMissing        Event = "onboarding.model-missing"
	EvKeyMissing          Event = "onboarding.key-missing"

	// Session entry.
	EvSummon          Event = "session.summon"   // #4 Sleeping, #29 Warm, #33 Settling
	EvKwsEnabled      Event = "kws.enabled"      // #5
	EvConfigInvalid   Event = "config.invalid"   // #6
	EvWakeWord        Event = "kws.wake-word"    // #7
	EvMuteKey         Event = "hotkey.mute"      // #8, #10
	EvWatchdogOverrun Event = "watchdog.overrun" // #9 (Armed self-loop + alert)

	// Voice turn.
	EvVadStop             Event = "asr.vad-stop"           // #11
	EvTimeoutFirstRound   Event = "listen.timeout-first"   // #12 (15s)
	EvTimeoutInSession    Event = "listen.timeout-sess"    // #12 (90s)
	EvTimeoutConversation Event = "listen.timeout-conv"    // #12 (30s)
	EvVeto                Event = "voice.veto"             // #13 Listening, #22 Confirming
	EvAudioDeviceLost     Event = "audio.device-lost"      // #14
	EvFirstToken          Event = "llm.first-token"        // #15
	EvBrainFailed         Event = "llm.timeout-or-network" // #16

	// Tool execution / approvals.
	EvApprovalNeeded  Event = "tools.approval-needed"   // #17
	EvToolsCompleted  Event = "tools.completed"         // #18
	EvRepeatThreshold Event = "tools.repeat-threshold"  // #19
	EvPathConflict    Event = "c20.path-conflict"       // #20
	EvConfirmExpired  Event = "confirm.expired"         // #21 (L1 countdown end)
	EvApprovalGranted Event = "approval.granted-native" // #23 (F2: native side only)
	EvApprovalDenied  Event = "approval.denied"         // #24
	EvApprovalTimeout Event = "approval.timeout-300s"   // #24
	EvQueueDrained    Event = "approval.queue-drained"  // #25

	// Session tail / modes.
	EvSpeakDone      Event = "tts.speak-done"           // #26 default, #28 conversation guard
	EvInterrupt      Event = "voice.interrupt"          // #27
	EvBargeIn        Event = "asr.barge-in-path-c"      // #41 user speech detected over TTS (AEC)
	EvConversationOn Event = "voice.conversation-on"    // #30
	EvWarmIdle       Event = "warm.idle-90s"            // #31
	EvSettleExpired  Event = "settling.expired-3s"      // #32
	EvTaskIntent     Event = "conversation.task-intent" // #42 (Path T hand-back)

	// Ambient / failure / recovery.
	EvNetworkDown       Event = "net.down"                 // #34 (any state)
	EvWatchdogFatal     Event = "watchdog.fatal"           // #35 (any state)
	EvPanicRecovered    Event = "panic.recovered"          // #36 (any state)
	EvDownloadCompleted Event = "model.download-completed" // #37
	EvDownloadFailed    Event = "model.download-failed"    // #37
	EvErrorAck          Event = "error.ack"                // #38 (user confirm / 10s auto)
	EvPathLockReleased  Event = "c20.path-lock-released"   // #39
	EvStuckContinue     Event = "stuck.continue"           // #40
	EvStuckAbandon      Event = "stuck.abandon"            // #40
)

// Facts carries the guard inputs for one dispatch. Nil Facts means the zero
// value (guards see false/0). The fields map 1:1 to the guard annotations in
// SPEC-08 §3; unrelated fields are simply ignored by the fired row's guard.
type Facts struct {
	KwsLoaded        bool // rows #10, #32: KWS loaded?
	HasToolCall      bool // row #15: first token carries a tool call?
	ApprovalLevel    int  // row #17: 1 (L1) or 2 (L2)
	HasAnnouncement  bool // rows #18, #25: TTS announcement pending?
	RepeatLevel      int  // row #19: current repeat threshold hit (3/5/8)
	QueueDepthAfter  int  // row #23: approval queue depth after the dequeue
	ConversationMode bool // row #26/#28: conversation mode active?
	SpeechMS         int  // row #11: voiced speech length in milliseconds
}
