package statemachine

// The authoritative D43 transition table (SPEC-08 §3 / PLAN.md D43): 20
// states, "anything not listed is illegal" (undefined = stop, D22 gate 3).
//
// Numbering reconciliation: the frozen D43 table carries 40 rows; SPEC-08 §3
// (the ticket's line-by-line authority) appends two numbered rows to it -
// #41 (Path C barge-in) and #42 (Conversation task intent) - for 42 rows
// total. All 42 are implemented and pinned by the row-presence test.
//
// Table-driven by contract: this file is DATA. machine.go only walks it. A
// few rows list two destinations in their To column (guard-branched); those
// are encoded as adjacent entries sharing one row number.
//
// SideEffects are event names only: the machine fires them on the EffectSink
// hook; executing them (SessionScope creation, model loading, panel pushes)
// belongs to later tickets. The default sink is an explicit no-op.

// AnyState is the From wildcard of the D43 "any" rows (#34/#35/#36). A
// wildcard row never shadows a specific row: dispatch resolves specific rows
// first (machine.go).
const AnyState State = "*"

// Row is one transition edge. Guard names are stable audit strings; GuardFn
// returning false makes dispatch skip to the next candidate row for the same
// (From, Event). Listen pins the Listening timeout phase a row installs when
// its destination is Listening (zero value = not a Listening entry).
type Row struct {
	D43         int   // row number in SPEC-08 §3 (1..40)
	From        State // AnyState for the wildcard rows
	Event       Event
	To          State
	Guard       string // "" = unconditional; else the D43 guard label
	GuardFn     func(*Facts) bool
	Listen      Phase    // Listening timeout context (see timeouts.go)
	SideEffects []string // effect hook names, fired in order
}

// Table is the frozen D43 table, in row order. Read-only by contract.
var Table = []Row{
	// ---- boot / onboarding ----
	{
		D43: 1, From: StateFirstRun, Event: EvOnboardingCompleted, To: StateSleeping,
		SideEffects: []string{"config.persist", "db.create-schema"},
	},
	{D43: 2, From: StateFirstRun, Event: EvModelMissing, To: StateDownloading},
	{D43: 3, From: StateFirstRun, Event: EvKeyMissing, To: StateUnconfigured},

	// ---- session entry ----
	{
		D43: 4, From: StateSleeping, Event: EvSummon, To: StateListening, Listen: PhaseFirstRound,
		SideEffects: []string{"session.scope-create", "speech.load-vad-asr"},
	},
	{
		D43: 5, From: StateSleeping, Event: EvKwsEnabled, To: StateArmed,
		SideEffects: []string{"kws.load"},
	},
	{D43: 6, From: StateSleeping, Event: EvConfigInvalid, To: StateUnconfigured},

	{
		D43: 7, From: StateArmed, Event: EvWakeWord, To: StateListening, Listen: PhaseFirstRound,
		SideEffects: []string{"session.scope-create", "speech.load-vad-asr", "kws.pause"},
	},
	{
		D43: 8, From: StateArmed, Event: EvMuteKey, To: StateMuted,
		SideEffects: []string{"kws.stop-inference"},
	},
	{
		D43: 9, From: StateArmed, Event: EvWatchdogOverrun, To: StateArmed,
		Guard:       "watchdog-overrun",
		SideEffects: []string{"kws.keep-alive-alert"},
	}, // D32③: KWS must NOT be unloaded

	{
		D43: 10, From: StateMuted, Event: EvMuteKey, To: StateArmed,
		Guard: "kws-loaded", GuardFn: func(f *Facts) bool { return f.KwsLoaded },
	},
	{
		D43: 10, From: StateMuted, Event: EvMuteKey, To: StateSleeping,
		Guard: "kws-not-loaded", GuardFn: func(f *Facts) bool { return !f.KwsLoaded },
	},

	// ---- voice turn ----
	{
		D43: 11, From: StateListening, Event: EvVadStop, To: StateThinking,
		Guard: "speech-ge-300ms", GuardFn: func(f *Facts) bool { return f.SpeechMS >= 300 },
		SideEffects: []string{"audio.stop-capture", "asr.punctuate", "input.taint-mark"},
	},
	{D43: 12, From: StateListening, Event: EvTimeoutFirstRound, To: StateSleeping},
	{D43: 12, From: StateListening, Event: EvTimeoutInSession, To: StateWarm},
	{D43: 12, From: StateListening, Event: EvTimeoutConversation, To: StateWarm},
	{
		D43: 13, From: StateListening, Event: EvVeto, To: StateWarm,
		SideEffects: []string{"audio.discard-buffer"},
	},
	{
		D43: 14, From: StateListening, Event: EvAudioDeviceLost, To: StateError,
		SideEffects: []string{"error.device-name"},
	},
	{
		D43: 15, From: StateThinking, Event: EvFirstToken, To: StateActing,
		Guard: "first-token-has-tool-call", GuardFn: func(f *Facts) bool { return f.HasToolCall },
		SideEffects: []string{"panel.stream-push"},
	},
	{
		D43: 15, From: StateThinking, Event: EvFirstToken, To: StateSpeaking,
		Guard:       "first-token-pure-text",
		GuardFn:     func(f *Facts) bool { return !f.HasToolCall },
		SideEffects: []string{"panel.stream-push"},
	},
	{
		D43: 16, From: StateThinking, Event: EvBrainFailed, To: StateError,
		SideEffects: []string{"error.classify-retry"},
	},

	// ---- tool execution / approvals ----
	{
		D43: 17, From: StateActing, Event: EvApprovalNeeded, To: StateConfirming,
		Guard: "level-L1", GuardFn: func(f *Facts) bool { return f.ApprovalLevel == 1 },
		SideEffects: []string{"confirm.countdown-start"},
	},
	{
		D43: 17, From: StateActing, Event: EvApprovalNeeded, To: StateAwaitingApproval,
		Guard: "level-L2", GuardFn: func(f *Facts) bool { return f.ApprovalLevel == 2 },
		SideEffects: []string{"approval.enqueue-c18"},
	},
	{
		D43: 18, From: StateActing, Event: EvToolsCompleted, To: StateSpeaking,
		Guard: "has-announcement", GuardFn: func(f *Facts) bool { return f.HasAnnouncement },
	},
	{
		D43: 18, From: StateActing, Event: EvToolsCompleted, To: StateSettling,
		Guard: "silent-completion", GuardFn: func(f *Facts) bool { return !f.HasAnnouncement },
	},
	{
		D43: 19, From: StateActing, Event: EvRepeatThreshold, To: StateActing,
		Guard: "repeat-below-8", GuardFn: func(f *Facts) bool { return f.RepeatLevel < 8 },
		SideEffects: []string{"agent.inject-reminder"},
	},
	{
		D43: 19, From: StateActing, Event: EvRepeatThreshold, To: StateStuck,
		Guard:   "repeat-hit-8",
		GuardFn: func(f *Facts) bool { return f.RepeatLevel >= 8 },
	},
	{
		D43: 20, From: StateActing, Event: EvPathConflict, To: StateQueued,
		SideEffects: []string{"panel.queue-waiting"},
	},
	{
		D43: 21, From: StateConfirming, Event: EvConfirmExpired, To: StateActing,
		SideEffects: []string{"tools.execute"},
	}, // fs.write MUST be temp+rename (C19)
	{
		D43: 22, From: StateConfirming, Event: EvVeto, To: StateActing,
		SideEffects: []string{"agent.cancel-tool-call", "ui.voice-cancel-availability"},
	},
	{
		D43: 23, From: StateAwaitingApproval, Event: EvApprovalGranted, To: StateActing,
		Guard: "queue-empty-after-dequeue", GuardFn: func(f *Facts) bool { return f.QueueDepthAfter <= 0 },
		SideEffects: []string{"approval.dequeue", "badge.decrement"},
	},
	{
		D43: 23, From: StateAwaitingApproval, Event: EvApprovalGranted, To: StateAwaitingApproval,
		Guard: "queue-still-non-empty", GuardFn: func(f *Facts) bool { return f.QueueDepthAfter > 0 },
		SideEffects: []string{"approval.dequeue", "badge.decrement"},
	},
	{
		D43: 24, From: StateAwaitingApproval, Event: EvApprovalDenied, To: StateActing,
		SideEffects: []string{"approval.cancel-call", "approval.replay-offer"},
	},
	{
		D43: 24, From: StateAwaitingApproval, Event: EvApprovalTimeout, To: StateActing,
		SideEffects: []string{"approval.cancel-call", "approval.replay-offer"},
	},
	{
		D43: 25, From: StateAwaitingApproval, Event: EvQueueDrained, To: StateSpeaking,
		Guard: "has-announcement", GuardFn: func(f *Facts) bool { return f.HasAnnouncement },
	},
	{
		D43: 25, From: StateAwaitingApproval, Event: EvQueueDrained, To: StateSettling,
		Guard: "silent-completion", GuardFn: func(f *Facts) bool { return !f.HasAnnouncement },
	},

	// ---- session tail / modes ----
	{
		D43: 26, From: StateSpeaking, Event: EvSpeakDone, To: StateWarm,
		Guard: "not-conversation", GuardFn: func(f *Facts) bool { return !f.ConversationMode },
		SideEffects: []string{"mic.keep-off", "panel.keep-alive"},
	},
	{
		D43: 27, From: StateSpeaking, Event: EvInterrupt, To: StateListening, Listen: PhaseInSession,
		SideEffects: []string{"tts.stop", "audio.release-output"},
	},
	{
		D43: 41, From: StateSpeaking, Event: EvBargeIn, To: StateListening, Listen: PhaseInSession,
		// SPEC-08 §3 row #41 (Path C, AEC): user speech detected over the
		// announcement. TTS stops and the output is released within 400ms;
		// the announcement audio must never reach ASR (D47).
		Guard:       "not-conversation",
		GuardFn:     func(f *Facts) bool { return !f.ConversationMode },
		SideEffects: []string{"tts.stop-bargein-400ms", "audio.release-output", "asr.exclude-playback"},
	},
	{
		D43: 41, From: StateSpeaking, Event: EvBargeIn, To: StateListening, Listen: PhaseConversation,
		Guard:       "conversation",
		GuardFn:     func(f *Facts) bool { return f.ConversationMode },
		SideEffects: []string{"tts.stop-bargein-400ms", "audio.release-output", "asr.exclude-playback"},
	},
	{
		D43: 28, From: StateSpeaking, Event: EvSpeakDone, To: StateListening, Listen: PhaseConversation,
		Guard: "conversation", GuardFn: func(f *Facts) bool { return f.ConversationMode },
		SideEffects: []string{"mic.enable", "conversation.ring-solid"},
	},
	{
		D43: 29, From: StateWarm, Event: EvSummon, To: StateListening, Listen: PhaseInSession,
		SideEffects: []string{"session.zero-load-resume"},
	}, // B4: zero model loading
	{
		D43: 30, From: StateWarm, Event: EvConversationOn, To: StateListening, Listen: PhaseConversation,
		// To column reads "Conversation -> Listening": the L2 privacy confirm
		// gates the pass through Conversation; the state machine itself moves
		// Warm straight to Listening with the conversation phase (and the
		// renderer shows Conversation while the confirm is pending).
		SideEffects: []string{"conversation.privacy-confirm-l2", "conversation.entered"},
	},
	{
		D43: 31, From: StateWarm, Event: EvWarmIdle, To: StateSettling,
		SideEffects: []string{"session.dispose-scope", "speech.unload-asr-tts", "mem.free-os-memory", "panel.destroy-or-hide"},
	},
	{
		D43: 32, From: StateSettling, Event: EvSettleExpired, To: StateSleeping,
		Guard: "kws-not-loaded", GuardFn: func(f *Facts) bool { return !f.KwsLoaded },
		SideEffects: []string{"mem.rss-verify-10s"},
	},
	{
		D43: 32, From: StateSettling, Event: EvSettleExpired, To: StateArmed,
		Guard:       "kws-loaded",
		GuardFn:     func(f *Facts) bool { return f.KwsLoaded },
		SideEffects: []string{"mem.rss-verify-10s"},
	},
	{
		D43: 33, From: StateSettling, Event: EvSummon, To: StateListening, Listen: PhaseFirstRound,
		SideEffects: []string{"settling.cancel-fallback"},
	},
	{
		D43: 42, From: StateConversation, Event: EvTaskIntent, To: StateThinking,
		SideEffects: []string{"conversation.suspend", "ctx.carry-c7-d47"},
	},

	// ---- ambient / failure / recovery ("任意" rows) ----
	{
		D43: 34, From: AnyState, Event: EvNetworkDown, To: StateNoNetwork,
		SideEffects: []string{"task.ctx-keep"},
	},
	{
		D43: 35, From: AnyState, Event: EvWatchdogFatal, To: StateWatchdogAlert,
		SideEffects: []string{"ui.one-click-restart"},
	},
	{
		D43: 36, From: AnyState, Event: EvPanicRecovered, To: StateError,
		SideEffects: []string{"diag.record-stack"},
	},
	{
		D43: 37, From: StateDownloading, Event: EvDownloadCompleted, To: StateFirstRun,
		SideEffects: []string{"model.verify-sha256-signature"},
	},
	{D43: 37, From: StateDownloading, Event: EvDownloadFailed, To: StateError},
	{
		D43: 38, From: StateError, Event: EvErrorAck, To: StateWarm,
		SideEffects: []string{"error.ack"},
	}, // restore target resolved by the error.ack consumer
	{D43: 39, From: StateQueued, Event: EvPathLockReleased, To: StateActing},
	{D43: 40, From: StateStuck, Event: EvStuckContinue, To: StateActing},
	{D43: 40, From: StateStuck, Event: EvStuckAbandon, To: StateSettling},
}
