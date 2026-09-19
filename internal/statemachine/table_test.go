package statemachine

import (
	"slices"
	"testing"
)

// TestAllFortyD43RowsPresent pins the reconciliation between the Go table and
// the authoritative SPEC-08 §3 numbering: the 40 frozen D43 rows (1..40) are
// all encoded, plus the two rows SPEC-08 §3 appends to them (#41 barge-in,
// #42 task intent) - 42 distinct numbers total, none invented.
func TestAllFortyD43RowsPresent(t *testing.T) {
	seen := map[int]bool{}
	for _, r := range Table {
		if r.D43 < 1 || r.D43 > 42 {
			t.Errorf("row %d outside numbering 1..42 (%s + %s)", r.D43, r.From, r.Event)
		}
		seen[r.D43] = true
	}
	for n := 1; n <= 42; n++ {
		if !seen[n] {
			t.Errorf("row %d missing from the table", n)
		}
	}
	if len(seen) != 42 {
		t.Errorf("got %d distinct rows, want 42 (40 D43 + #41/#42)", len(seen))
	}
}

// TestEveryStateVocabularyCovered asserts the table's From/To columns only
// reference the 20 frozen states plus the wildcard.
func TestEveryStateVocabularyCovered(t *testing.T) {
	for _, r := range Table {
		if r.From != AnyState && !Valid(r.From) {
			t.Errorf("row %d: unknown From %q", r.D43, r.From)
		}
		if !Valid(r.To) {
			t.Errorf("row %d: unknown To %q", r.D43, r.To)
		}
	}
}

// d43Case is one legal dispatch mirroring one line of SPEC-08 §3. The facts
// satisfy the row's guard; want is the D43 To column.
type d43Case struct {
	name    string
	from    State
	event   Event
	facts   *Facts
	want    State
	phase   Phase // expected Machine.Phase() after the dispatch
	skipPh  bool  // skip the phase assertion
	effects []string
}

// legalRows is the line-by-line reconciliation of the 40 D43 rows. One entry
// per table edge (guard-branched rows appear per branch).
var legalRows = []d43Case{
	{name: "#1 onboarding done", from: StateFirstRun, event: EvOnboardingCompleted, want: StateSleeping,
		effects: []string{"config.persist", "db.create-schema"}},
	{name: "#2 model missing", from: StateFirstRun, event: EvModelMissing, want: StateDownloading},
	{name: "#3 key missing", from: StateFirstRun, event: EvKeyMissing, want: StateUnconfigured},

	{name: "#4 sleeping summon", from: StateSleeping, event: EvSummon, want: StateListening, phase: PhaseFirstRound,
		effects: []string{"session.scope-create", "speech.load-vad-asr"}},
	{name: "#5 kws enabled", from: StateSleeping, event: EvKwsEnabled, want: StateArmed,
		effects: []string{"kws.load"}},
	{name: "#6 config invalid", from: StateSleeping, event: EvConfigInvalid, want: StateUnconfigured},

	{name: "#7 wake word", from: StateArmed, event: EvWakeWord, want: StateListening, phase: PhaseFirstRound,
		effects: []string{"session.scope-create", "speech.load-vad-asr", "kws.pause"}},
	{name: "#8 mute", from: StateArmed, event: EvMuteKey, want: StateMuted,
		effects: []string{"kws.stop-inference"}},
	{name: "#9 watchdog overrun (self-loop, KWS kept)", from: StateArmed, event: EvWatchdogOverrun,
		want: StateArmed, effects: []string{"kws.keep-alive-alert"}},

	{name: "#10 unmute with KWS", from: StateMuted, event: EvMuteKey, facts: &Facts{KwsLoaded: true},
		want: StateArmed},
	{name: "#10 unmute without KWS", from: StateMuted, event: EvMuteKey, facts: &Facts{},
		want: StateSleeping},

	{name: "#11 vad stop", from: StateListening, event: EvVadStop, facts: &Facts{SpeechMS: 300},
		want: StateThinking, skipPh: true,
		effects: []string{"audio.stop-capture", "asr.punctuate", "input.taint-mark"}},
	{name: "#12 first-round timeout", from: StateListening, event: EvTimeoutFirstRound, want: StateSleeping},
	{name: "#12 in-session timeout", from: StateListening, event: EvTimeoutInSession, want: StateWarm},
	{name: "#12 conversation timeout", from: StateListening, event: EvTimeoutConversation, want: StateWarm},
	{name: "#13 veto", from: StateListening, event: EvVeto, want: StateWarm,
		effects: []string{"audio.discard-buffer"}},
	{name: "#14 audio device lost", from: StateListening, event: EvAudioDeviceLost, want: StateError,
		effects: []string{"error.device-name"}},

	{name: "#15 first token with tool call", from: StateThinking, event: EvFirstToken,
		facts: &Facts{HasToolCall: true}, want: StateActing, effects: []string{"panel.stream-push"}},
	{name: "#15 first token pure text", from: StateThinking, event: EvFirstToken,
		facts: &Facts{}, want: StateSpeaking, effects: []string{"panel.stream-push"}},
	{name: "#16 brain failure", from: StateThinking, event: EvBrainFailed, want: StateError,
		effects: []string{"error.classify-retry"}},

	{name: "#17 L1 approval needed", from: StateActing, event: EvApprovalNeeded,
		facts: &Facts{ApprovalLevel: 1}, want: StateConfirming, effects: []string{"confirm.countdown-start"}},
	{name: "#17 L2 approval needed", from: StateActing, event: EvApprovalNeeded,
		facts: &Facts{ApprovalLevel: 2}, want: StateAwaitingApproval, effects: []string{"approval.enqueue-c18"}},
	{name: "#18 tools done with announcement", from: StateActing, event: EvToolsCompleted,
		facts: &Facts{HasAnnouncement: true}, want: StateSpeaking},
	{name: "#18 tools done silent", from: StateActing, event: EvToolsCompleted,
		facts: &Facts{}, want: StateSettling},
	{name: "#19 repeat below 8 (self-loop)", from: StateActing, event: EvRepeatThreshold,
		facts: &Facts{RepeatLevel: 5}, want: StateActing, effects: []string{"agent.inject-reminder"}},
	{name: "#19 repeat hits 8", from: StateActing, event: EvRepeatThreshold,
		facts: &Facts{RepeatLevel: 8}, want: StateStuck},
	{name: "#20 path conflict", from: StateActing, event: EvPathConflict, want: StateQueued,
		effects: []string{"panel.queue-waiting"}},

	{name: "#21 confirm countdown expired", from: StateConfirming, event: EvConfirmExpired,
		want: StateActing, effects: []string{"tools.execute"}},
	{name: "#22 confirming veto", from: StateConfirming, event: EvVeto, want: StateActing,
		effects: []string{"agent.cancel-tool-call", "ui.voice-cancel-availability"}},

	{name: "#23 approval granted, queue empties", from: StateAwaitingApproval, event: EvApprovalGranted,
		facts: &Facts{QueueDepthAfter: 0}, want: StateActing,
		effects: []string{"approval.dequeue", "badge.decrement"}},
	{name: "#23 approval granted, queue non-empty (stay)", from: StateAwaitingApproval, event: EvApprovalGranted,
		facts: &Facts{QueueDepthAfter: 2}, want: StateAwaitingApproval,
		effects: []string{"approval.dequeue", "badge.decrement"}},
	{name: "#24 approval denied", from: StateAwaitingApproval, event: EvApprovalDenied, want: StateActing,
		effects: []string{"approval.cancel-call", "approval.replay-offer"}},
	{name: "#24 approval timeout 300s", from: StateAwaitingApproval, event: EvApprovalTimeout, want: StateActing,
		effects: []string{"approval.cancel-call", "approval.replay-offer"}},
	{name: "#25 queue drained with announcement", from: StateAwaitingApproval, event: EvQueueDrained,
		facts: &Facts{HasAnnouncement: true}, want: StateSpeaking},
	{name: "#25 queue drained silent", from: StateAwaitingApproval, event: EvQueueDrained,
		facts: &Facts{}, want: StateSettling},

	{name: "#26 speak done (default, mic stays off)", from: StateSpeaking, event: EvSpeakDone,
		facts: &Facts{}, want: StateWarm, effects: []string{"mic.keep-off", "panel.keep-alive"}},
	{name: "#27 interrupt", from: StateSpeaking, event: EvInterrupt, want: StateListening, phase: PhaseInSession,
		effects: []string{"tts.stop", "audio.release-output"}},
	{name: "#41 barge-in (default half-duplex)", from: StateSpeaking, event: EvBargeIn,
		facts: &Facts{}, want: StateListening, phase: PhaseInSession,
		effects: []string{"tts.stop-bargein-400ms", "audio.release-output", "asr.exclude-playback"}},
	{name: "#41 barge-in in conversation", from: StateSpeaking, event: EvBargeIn,
		facts: &Facts{ConversationMode: true}, want: StateListening, phase: PhaseConversation,
		effects: []string{"tts.stop-bargein-400ms", "audio.release-output", "asr.exclude-playback"}},
	{name: "#28 speak done in conversation", from: StateSpeaking, event: EvSpeakDone,
		facts: &Facts{ConversationMode: true}, want: StateListening, phase: PhaseConversation,
		effects: []string{"mic.enable", "conversation.ring-solid"}},

	{name: "#29 warm summon (zero load)", from: StateWarm, event: EvSummon, want: StateListening,
		phase: PhaseInSession, effects: []string{"session.zero-load-resume"}},
	{name: "#30 conversation on", from: StateWarm, event: EvConversationOn, want: StateListening,
		phase:   PhaseConversation,
		effects: []string{"conversation.privacy-confirm-l2", "conversation.entered"}},
	{name: "#31 warm 90s idle", from: StateWarm, event: EvWarmIdle, want: StateSettling,
		effects: []string{"session.dispose-scope", "speech.unload-asr-tts", "mem.free-os-memory", "panel.destroy-or-hide"}},

	{name: "#32 settling done without KWS", from: StateSettling, event: EvSettleExpired,
		facts: &Facts{}, want: StateSleeping, effects: []string{"mem.rss-verify-10s"}},
	{name: "#32 settling done with KWS", from: StateSettling, event: EvSettleExpired,
		facts: &Facts{KwsLoaded: true}, want: StateArmed, effects: []string{"mem.rss-verify-10s"}},
	{name: "#33 settling re-summon", from: StateSettling, event: EvSummon, want: StateListening,
		phase: PhaseFirstRound, effects: []string{"settling.cancel-fallback"}},

	{name: "#42 conversation task intent", from: StateConversation, event: EvTaskIntent, want: StateThinking,
		effects: []string{"conversation.suspend", "ctx.carry-c7-d47"}},

	{name: "#37 download completed", from: StateDownloading, event: EvDownloadCompleted, want: StateFirstRun,
		effects: []string{"model.verify-sha256-signature"}},
	{name: "#37 download failed", from: StateDownloading, event: EvDownloadFailed, want: StateError},

	{name: "#38 error ack", from: StateError, event: EvErrorAck, want: StateWarm,
		effects: []string{"error.ack"}},
	{name: "#39 path lock released", from: StateQueued, event: EvPathLockReleased, want: StateActing},

	{name: "#40 stuck continue", from: StateStuck, event: EvStuckContinue, want: StateActing},
	{name: "#40 stuck abandon", from: StateStuck, event: EvStuckAbandon, want: StateSettling},

	// Wildcard rows #34/#35/#36 - sampled from several source states.
	{name: "#34 network down (from Sleeping)", from: StateSleeping, event: EvNetworkDown, want: StateNoNetwork,
		effects: []string{"task.ctx-keep"}},
	{name: "#34 network down (from Acting)", from: StateActing, event: EvNetworkDown, want: StateNoNetwork,
		effects: []string{"task.ctx-keep"}},
	{name: "#35 watchdog fatal (from Warm)", from: StateWarm, event: EvWatchdogFatal, want: StateWatchdogAlert,
		effects: []string{"ui.one-click-restart"}},
	{name: "#35 watchdog fatal (from Listening)", from: StateListening, event: EvWatchdogFatal, want: StateWatchdogAlert,
		effects: []string{"ui.one-click-restart"}},
	{name: "#36 panic recovered (from Thinking)", from: StateThinking, event: EvPanicRecovered, want: StateError,
		effects: []string{"diag.record-stack"}},
	{name: "#36 panic recovered (from Muted)", from: StateMuted, event: EvPanicRecovered, want: StateError,
		effects: []string{"diag.record-stack"}},
}

// TestEveryLegalRowFires drives every legal edge and asserts destination,
// listening phase and fired side-effect names (state+side-effect spy per the
// ticket acceptance).
func TestEveryLegalRowFires(t *testing.T) {
	for _, tc := range legalRows {
		t.Run(tc.name, func(t *testing.T) {
			var fired []Effect
			m := New(Options{Initial: tc.from, Sink: func(e Effect) { fired = append(fired, e) }})
			got, err := m.Dispatch(tc.event, tc.facts)
			if err != nil {
				t.Fatalf("Dispatch(%s) from %s: unexpected rejection: %v", tc.event, tc.from, err)
			}
			if got != tc.want {
				t.Fatalf("Dispatch(%s) from %s: got %s, want %s", tc.event, tc.from, got, tc.want)
			}
			if !tc.skipPh && m.Phase() != tc.phase {
				t.Errorf("phase after dispatch: got %d, want %d", m.Phase(), tc.phase)
			}
			gotNames := make([]string, len(fired))
			for i, e := range fired {
				gotNames[i] = e.Name
			}
			if !slices.Equal(gotNames, tc.effects) {
				t.Errorf("fired effects = %v, want %v", gotNames, tc.effects)
			}
		})
	}
}

// TestIllegalTransitionsRejected walks EVERY (state, event) pair and demands
// rejection for every pair without a table row. The wildcard rows make
// EvNetworkDown / EvWatchdogFatal / EvPanicRecovered legal everywhere; every
// other event must be legal only where the table says so.
func TestIllegalTransitionsRejected(t *testing.T) {
	states := []State{
		StateFirstRun, StateSleeping, StateArmed, StateMuted,
		StateListening, StateThinking, StateActing, StateSpeaking,
		StateWarm, StateConversation, StateConfirming, StateAwaitingApproval,
		StateSettling, StateDownloading, StateUnconfigured, StateNoNetwork,
		StateError, StateQueued, StateStuck, StateWatchdogAlert,
	}
	if len(states) != 20 {
		t.Fatalf("state vocabulary drifted: %d states", len(states))
	}

	events := []Event{
		EvOnboardingCompleted, EvModelMissing, EvKeyMissing,
		EvSummon, EvKwsEnabled, EvConfigInvalid, EvWakeWord, EvMuteKey, EvWatchdogOverrun,
		EvVadStop, EvTimeoutFirstRound, EvTimeoutInSession, EvTimeoutConversation,
		EvVeto, EvAudioDeviceLost, EvFirstToken, EvBrainFailed,
		EvApprovalNeeded, EvToolsCompleted, EvRepeatThreshold, EvPathConflict,
		EvConfirmExpired, EvApprovalGranted, EvApprovalDenied, EvApprovalTimeout, EvQueueDrained,
		EvSpeakDone, EvInterrupt, EvBargeIn, EvConversationOn, EvWarmIdle, EvSettleExpired, EvTaskIntent,
		EvDownloadCompleted, EvDownloadFailed, EvErrorAck, EvPathLockReleased,
		EvStuckContinue, EvStuckAbandon,
	}

	rejected := 0
	for _, s := range states {
		for _, ev := range events {
			legal := false
			for _, r := range Table {
				if r.Event == ev && (r.From == s || r.From == AnyState) {
					legal = true
					break
				}
			}
			if legal {
				continue
			}
			m := New(Options{Initial: s})
			got, err := m.Dispatch(ev, &Facts{KwsLoaded: true, HasToolCall: true, ApprovalLevel: 2,
				HasAnnouncement: true, RepeatLevel: 8, QueueDepthAfter: 1, ConversationMode: true, SpeechMS: 999})
			if err == nil {
				t.Errorf("pair (%s, %s) accepted -> %s; must be rejected (undefined = stop)", s, ev, got)
			} else {
				var ilegal *IllegalTransitionError
				asIlegal := asErr(err, &ilegal)
				if !asIlegal {
					t.Errorf("pair (%s, %s): expected *IllegalTransitionError, got %T", s, ev, err)
				}
				rejected++
			}
			if m.State() != s {
				t.Errorf("pair (%s, %s): state mutated on rejection -> %s", s, ev, m.State())
			}
			if m.TimersAlive() != 0 && s == StateSleeping {
				t.Errorf("rejected dispatch armed a timer in Sleeping")
			}
		}
	}
	if rejected < 20 {
		t.Errorf("suspiciously few rejections (%d); table or matrix drifted", rejected)
	}
}

func asErr(err error, target **IllegalTransitionError) bool {
	if e, ok := err.(*IllegalTransitionError); ok {
		*target = e
		return true
	}
	return false
}

// TestWildcardNeverShadowsSpecificRow: in Error, the wildcard #36 (panic ->
// Error) and the specific #38 (ack -> Warm) coexist; a specific row must win
// over a wildcard for the same event only when events collide, which they
// never do here - pin the dispatch priority anyway.
func TestWildcardNeverShadowsSpecificRow(t *testing.T) {
	// EvErrorAck is specific to Error; from Error it must reach Warm even
	// though wildcard rows exist in the table.
	m := New(Options{Initial: StateError})
	got, err := m.Dispatch(EvErrorAck, nil)
	if err != nil || got != StateWarm {
		t.Fatalf("EvErrorAck from Error: got (%s, %v)", got, err)
	}
	// From any other state the same event is illegal (no wildcard covers it).
	m = New(Options{Initial: StateSleeping})
	if _, err := m.Dispatch(EvErrorAck, nil); err == nil {
		t.Fatal("EvErrorAck from Sleeping must be rejected")
	}
}
