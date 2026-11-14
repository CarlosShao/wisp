//go:build windows

package main

// Ticket 395's readings, driven through the production assembly rather than
// around it. The voice leg is the real residentVoice; the turn crosses into the
// REAL composition root `wisp run` uses (runTextTask) over the C5 golden SSE
// seam; the agent really dispatches one allowlisted fs.read (its tool_call row is
// read back off the real store); and the reply is announced through the
// PRODUCTION speak() path (real matcha checkpoint, a recording render seam).
//
// The ONE thing this case injects that HEAD does not wire for real is the ASR
// transcript text (the exported AsrEngine seam wire-r1 already uses), so the case
// does not depend on the KWS/ASR audio fixture or an exclusive sound card. The
// ball walk it measures - Thinking -> Acting -> Speaking -> Warm - is what the
// production dispatch itself produces from the turn's REAL tool-call count.
//
// Before ticket 395 the announce path dispatched D43#15 with HasToolCall:false
// unconditionally, so a turn that really carried a tool call walked the
// pure-text arm (Thinking -> Speaking) and the Acting family never saw a
// production dispatch. This case is the positive reading: the tool arm fires,
// Acting shows on the ball, and the fs.read row proves the turn really had the
// tool call the arm is gated on.
//
// Staging is 15e's: sherpa DLLs next to the test binary, WISP_SPIKE_MODELS /
// WISP_SILERO_MODEL / WISP_TTS_DIR pointed at the pinned trees, and
// -test.run=Test395r1 naming this family.

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/llm/golden"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/secret"
	"github.com/CarlosShao/wisp/internal/speech"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// Test395r1VoiceToolTurnWalksActing is AC#1 and AC#2 in one run: a voice turn
// whose first SSE response carries a tool call (the 16-r1 demo shape) must walk
// D43#15's TOOL arm into Acting and D43#18 out to Speaking for the announcement,
// with the Acting state really projected onto the ball.
func Test395r1VoiceToolTurnWalksActing(t *testing.T) {
	dataDir := voice15eDataDir(t)

	// The note fs.read reads, and the config that points the agent at the golden
	// server while allowlisting exactly this data root.
	note := filepath.Join(dataDir, "note.txt")
	if err := os.WriteFile(note, []byte("wisp 票 395 演示笔记\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	srv, rep, err := golden.ServerFromBytes([]byte(demo16Golden(t, note)))
	if err != nil {
		t.Fatalf("golden server: %v", err)
	}
	rep.Pace = true
	t.Cleanup(srv.Close)
	body := fmt.Sprintf(`schema_version = 2

[llm]
text_chain = ["acme/m1"]
timeout_ms = 120000

[llm.retry]
max = 1
backoff_ms = 1

[llm.providers.acme]
protocol = "openai-chat"
base_url = %q
api_key_ref = "dpapi:acme"

[llm.providers.acme.models.m1]
enabled = true
context_window = 128000

[fs]
allowed_dirs = [%q]
`, srv.URL+"/v1", filepath.ToSlash(dataDir))
	if err := os.WriteFile(filepath.Join(dataDir, configFileName), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	st, err := secret.NewStore(dataDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Store("dpapi:acme", fakeStoreKey); err != nil {
		t.Fatal(err)
	}

	t.Setenv("WISP_MODELS_MANIFEST", voice15eRepoFile(t, filepath.Join("build", "models", "manifest.json")))
	store := voice15eInstallAsr(t, dataDir)

	prov := risk.NewProvenance(risk.ProvOptions{NoProbe: true})
	speaker := &voiceWireSpeaker{rate: 22050}

	// The bridge between the voice leg and the REAL task pipeline: the same
	// composition root `wisp run` uses, with the runtime hooks wired exactly as
	// resident_windows.go wires them at boot - including, since ticket 395, the
	// turn-shape seam.
	bridge := &demo16Bridge{
		dataDir: dataDir,
		stdout:  &strings.Builder{},
		stderr:  &strings.Builder{},
		turned:  make(chan struct{}, 8),
	}
	var rv *residentVoice
	bridge.onRuntime = func(rt *agentRuntime) {
		rt.voiceAnnounce = rv.announce
		rt.voiceTurnTools = rv.noteTurnToolCalls
		rt.settle = rv.residency.settle
	}

	asr := &voiceWireAsr{text: "读一下我的笔记"}
	plan := voicePlan{
		VadModelPath: voice15eVadModel(t),
		AsrDir:       store,
		Source:       "395-r1 test",
		HandOff:      residentVoiceHandOff(dataDir),
		Tts:          voiceWireTtsFiles(t), // real matcha checkpoint: the production TTS factory loads it
	}
	rv, frames := voiceWireStart(t, plan, &residentBall{}, bridge, func(rv *residentVoice) {
		rv.residency.asrLoader = func() (speech.AsrEngine, error) { return asr, nil }
		rv.newSpeaker = func(rate int) (voiceSpeaker, error) { return speaker, nil }
		rv.attachProvenance(prov)
	})
	voice15eArm(t, frames, rv)

	// One utterance: VAD stop -> fake-ASR final -> the real composition root ->
	// the golden fs.read -> the spoken close -> the production announce.
	wav := filepath.Join(dataDir, "turn.wav")
	voice15eWriteWav(t, wav, voice15eSynth(voice15eLeadMS, voice15eSpeechMS, voice15eTailMS))
	voiceWireFeed(t, frames, wav, rv, 1)

	select {
	case <-bridge.turned:
	case <-time.After(120 * time.Second):
		_, _, _, out, errb := bridge.snapshot()
		t.Fatalf("the agent turn never finished\nstdout:\n%s\nstderr:\n%s", out, errb)
	}
	voice15eWaitFor(t, "the production announce to speak through the render seam", func() bool {
		return rv.ttsSpoken.Load() >= 1
	})

	// ---- the record.
	// TICKET 445 AC#4: the wait above is on ttsSpoken, which is incremented one
	// statement BEFORE the #26 dispatch (resident_voice_engines_windows.go:1265
	// vs :1282) - so the record below would be judged while its last row can
	// still be one append away. Wait for that row to land, then take the record.
	waitWalked445(t, rv, statemachine.StateSpeaking, statemachine.EvSpeakDone, statemachine.StateWarm,
		"the announce's D43#26 Speaking -> Warm walk")
	timeline := rv.timelineSnapshot()
	for _, tr := range timeline {
		t.Logf("D43_395 %s+%s->%s at_sample=%d", tr.From, tr.Event, tr.To, tr.AtSample)
	}
	ballStates := rv.ball.voiceStatesSnapshot()
	t.Logf("BALL_STATES_395 %v", ballStates)

	// AC#1: D43#15's TOOL arm was really dispatched (Thinking -> Acting), not the
	// pure-text arm.
	if !demo16Walked(rv, statemachine.StateThinking, statemachine.EvFirstToken, statemachine.StateActing) {
		t.Fatalf("no D43#15 TOOL-arm walk Thinking+EvFirstToken->Acting: the turn carried a tool call but the announce still walked the pure-text arm: %+v", timeline)
	}
	// AC#1/#2: the Acting family was left by its own legal row, #18 Acting ->
	// Speaking (has-announcement), because this sentence IS the announcement.
	if !demo16Walked(rv, statemachine.StateActing, statemachine.EvToolsCompleted, statemachine.StateSpeaking) {
		t.Fatalf("no D43#18 Acting+EvToolsCompleted->Speaking walk: %+v", timeline)
	}
	if !demo16Walked(rv, statemachine.StateSpeaking, statemachine.EvSpeakDone, statemachine.StateWarm) {
		t.Fatalf("no D43#26 Speaking->Warm walk: %+v", timeline)
	}
	// AC#2: Acting really showed on the ball (not merely a machine-bookkeeping
	// edge nobody projected).
	if !ballHasState(ballStates, statemachine.StateActing) {
		t.Fatalf("the ball never projected Acting: states = %v", ballStates)
	}
	// The order is Thinking -> Acting -> Speaking -> Warm, in that sequence.
	if !ballStateSequence(ballStates,
		statemachine.StateThinking, statemachine.StateActing, statemachine.StateSpeaking, statemachine.StateWarm) {
		t.Fatalf("ball state sequence is not Thinking->Acting->Speaking->Warm: %v", ballStates)
	}

	transcripts, exits, _, out, errb := bridge.snapshot()
	if len(exits) != 1 || exits[0] != 0 {
		t.Fatalf("the agent turn exits = %v, want exactly [0]\nstdout:\n%s\nstderr:\n%s", exits, out, errb)
	}
	t.Logf("TRANSCRIPT_395 %q", transcripts)
	t.Logf("REPLY_395 %q (streamed to the sink, spoken, notified)", demo16Reply)

	// The turn REALLY carried a tool call: the fs.read row is in the real store.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	m := taskIDRe.FindStringSubmatch(out)
	if m == nil {
		t.Fatalf("the run printed no status line:\n%s", out)
	}
	st2, err := memory.Open(dataDir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st2.Close() })
	rows, err := st2.ListToolCallsByTask(ctx, m[1])
	if err != nil {
		t.Fatal(err)
	}
	var read *memory.ToolCall
	for i := range rows {
		t.Logf("TOOL_ROW_395 tool=%s level=%s decision=%s outcome=%s class=%q corr=%q",
			rows[i].Tool, rows[i].RiskLevel, rows[i].Decision, rows[i].Outcome, rows[i].ErrorClass, rows[i].CorrelationID)
		if rows[i].Tool == "fs.read" {
			read = &rows[i]
		}
	}
	if read == nil {
		t.Fatalf("no fs.read tool_call row for task %s: %+v", m[1], rows)
	}
	if read.Outcome != "success" {
		t.Fatalf("the fs.read row outcome = %s, want success: the tool arm is gated on a tool call that must really have run", read.Outcome)
	}
}

// ballHasState reports whether the ball was ever projected onto want.
func ballHasState(states []statemachine.State, want statemachine.State) bool {
	for _, s := range states {
		if s == want {
			return true
		}
	}
	return false
}

// ballStateSequence reports whether wants appear as a contiguous subsequence of
// states, in order - the exact Thinking->Acting->Speaking->Warm run this ticket
// exists to make the ball show.
func ballStateSequence(states []statemachine.State, wants ...statemachine.State) bool {
	for i := 0; i+len(wants) <= len(states); i++ {
		ok := true
		for j := range wants {
			if states[i+j] != wants[j] {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}
