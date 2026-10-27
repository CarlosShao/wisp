//go:build windows

package main

// Ticket 26's one needle (26r3), read back: the announce path now drives
// internal/audio's HalfDuplexGate.SetSpeaking, and these cases are the readings
// the evidence table quotes.
//
//	#1  the wire          announce -> SetSpeaking(true) mid-speech -> false
//	                      after; every outcome releases (the defer runs on the
//	                      refused / no-engine / no-device / silent paths too)
//	#2  the self-echo     speech-time frames fed into the seam while the gate
//	     fixture          is closed: the reading is the GATE's state (PathT
//	                      speaking = capture refused) and the transcript delta,
//	                      never a claim that a device was switched off
//	#3  the zero-value    nil-safe counters, so a boot without a capture leg
//	                      cannot fake a reading
//
// The gate here is a RECORDER, not internal/audio's gate: production speak()
// drives whatever sits behind rv.captureGate, and the recorder pins the calls
// plus, through its probe, the chain's own state at the moment of each call.
// The reverse criteria (26r3-M1 removes SetSpeaking(true), 26r3-M2 removes the
// release) are run OUT OF TREE against the same case: #1 carries the two named
// red shapes those overlays must produce, so a mutant fails this file's first
// assertions with the M1/M2 sentence instead of an accident.
//
// Staging: the same contract the wire-r1 cases run under (DLLs next to the
// test binary, WISP_* for the models). Run with -test.run=Test26R3 only.

import (
	"context"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/session"
)

// voice26r3Gate records the SetSpeaking calls the announce path makes: the
// counts, the flag the calls leave behind, the ordered transitions, and - via
// the probe - the chain state at the instant of each call.
type voice26r3Gate struct {
	mu       sync.Mutex
	closed   int
	reopened int
	speaking bool
	muted    bool
	events   []voice26r3GateEvent

	// probe is wired by the case to read the live leg without a race window:
	// SetSpeaking runs ON the owner goroutine inside speak(), so what it reads
	// back is the chain's state mid-speech.
	probe func() (asrGone, asrLoaded bool, drained, finals uint64)

	asrGoneAtClose   bool
	asrLoadedAtClose bool
	drainedAtClose   uint64
	finalsAtClose    uint64
	drainedAtOpen    uint64
	finalsAtOpen     uint64
}

type voice26r3GateEvent struct {
	muted    bool
	speaking bool
}

func (g *voice26r3Gate) SetSpeaking(speaking bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if speaking {
		g.closed++
		if g.probe != nil {
			g.asrGoneAtClose, g.asrLoadedAtClose, g.drainedAtClose, g.finalsAtClose = g.probe()
		}
	} else {
		g.reopened++
		if g.probe != nil {
			_, _, g.drainedAtOpen, g.finalsAtOpen = g.probe()
		}
	}
	g.speaking = speaking
	g.events = append(g.events, voice26r3GateEvent{muted: g.muted, speaking: speaking})
}

func (g *voice26r3Gate) snapshot() (closed, reopened int, speaking bool, events []voice26r3GateEvent) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.closed, g.reopened, g.speaking, append([]voice26r3GateEvent(nil), g.events...)
}

func (g *voice26r3Gate) windows() (drainedDelta, finalsDelta uint64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.drainedAtOpen - g.drainedAtClose, g.finalsAtOpen - g.finalsAtClose
}

// Test26R3AnnounceClosesTheCaptureGate is reading #1 plus the self-echo
// fixture: one real dictation turn, then one real announcement with the gate
// driven by the production speak() body, and speech-time frames fed into the
// seam while the gate is closed.
func Test26R3AnnounceClosesTheCaptureGate(t *testing.T) {
	dataDir := voice15eDataDir(t)
	voice15eWriteConfig(t, dataDir, nil)
	t.Setenv("WISP_MODELS_MANIFEST", voice15eRepoFile(t, filepath.Join("build", "models", "manifest.json")))
	store := voice15eInstallAsr(t, dataDir)

	speaker := &voiceWireSpeaker{rate: 22050}
	gate := &voice26r3Gate{}
	sink := &voice15eSink{}
	ball := &residentBall{}
	plan := voicePlan{
		VadModelPath: voice15eVadModel(t),
		AsrDir:       store,
		Source:       "26r3 test",
		HandOff:      residentVoiceHandOff(dataDir),
		Tts:          voiceWireTtsFiles(t),
	}
	rv, frames := voiceWireStart(t, plan, ball, sink, func(rv *residentVoice) {
		rv.newSpeaker = func(rate int) (voiceSpeaker, error) { return speaker, nil }
		rv.attachCaptureGate(gate)
		// The probe reads the live leg. SetSpeaking fires on the owner
		// goroutine while run() is parked inside speakQueued, so these loads
		// are mid-speech facts, not post-hoc snapshots.
		gate.probe = func() (asrGone, asrLoaded bool, drained, finals uint64) {
			return rv.asrGone.Load(), rv.asrLoaded.Load(), rv.drained.Load(), rv.finals.Load()
		}
	})
	voice15eArm(t, frames, rv)

	// Turn 1: a real dictation through the real seam, so the announcement below
	// really does swap the engine family out from under a live chain.
	wav1 := filepath.Join(dataDir, "turn1.wav")
	voice15eWriteWav(t, wav1, voice15eSynth(voice15eLeadMS, voice15eSpeechMS, voice15eTailMS))
	voiceWireFeed(t, frames, wav1, rv, 1)
	voice15eWaitFor(t, "the first turn to land in Warm", func() bool {
		return rv.warmTurnEnd.Load() == 1
	})

	if !rv.announce("任务已完成：26r3 的采集门接线已经落上。") {
		t.Fatalf("the announce path refused the sentence; problems=%v", rv.problemsSnapshot())
	}
	voice15eWaitFor(t, "the announcement to be spoken", func() bool {
		return rv.ttsSpoken.Load() == 1
	})
	voice15eWaitFor(t, "the capture gate to be reopened after the sentence", func() bool {
		return rv.speakGateOff.Load() == 1
	})

	closed, reopened, speaking, events := gate.snapshot()
	// THE TWO REVERSE SHAPES. 26r3-M1 (the out-of-tree overlay removes
	// SetSpeaking(true)) lands here: closed=0, the announce path never closes
	// the capture gate. 26r3-M2 (the overlay removes the release) lands on the
	// wait above timing out AND here: reopened=0, speaking stuck true - capture
	// stays closed after the sentence, so no later utterance can ever be
	// captured (the observable consequence of being stuck in Speaking).
	if closed != 1 {
		t.Fatalf("MUTANT 26r3-M1 SHAPE: the announce path never closed the capture gate (closed=%d reopened=%d): SetSpeaking(true) is not on the path", closed, reopened)
	}
	if reopened != 1 || speaking {
		t.Fatalf("MUTANT 26r3-M2 SHAPE: the capture gate was left in Speaking after the sentence (reopened=%d speaking=%v): SetSpeaking(false) is not on the path; capture stays closed and no later utterance can ever be captured", reopened, speaking)
	}
	if got1, got2 := rv.setSpeakingCount(); got1 != uint64(closed) || got2 != uint64(reopened) {
		t.Fatalf("the leg's own counters say closed=%d reopened=%d, the gate saw %d/%d", got1, got2, closed, reopened)
	}
	if len(events) != 2 || !events[0].speaking || events[1].speaking {
		t.Fatalf("the gate events are %v, want SetSpeaking(true)-then-SetSpeaking(false) while not muted", events)
	}
	if !gate.asrGoneAtClose || !gate.asrLoadedAtClose {
		t.Fatalf("at the close the chain said asrGone=%v asrLoaded=%v: the gate was not closed mid-speech, before the phase-1 re-arm", gate.asrGoneAtClose, gate.asrLoadedAtClose)
	}
	_, samples, _, _ := speaker.energy()
	if samples == 0 {
		t.Fatalf("the announcement emitted no samples; the gate wiring must not break the speak path")
	}
	if st := rv.residency.stage(); st != session.StageWarm {
		t.Fatalf("stage after the announcement = %q, want Warm (the needle must not disturb the turn-end path)", string(st))
	}
	_, finalsDelta := gate.windows()
	if finalsDelta != 0 {
		t.Fatalf("%d transcript(s) appeared while the gate was closed: the announcement's own window produced ASR output (self-echo)", finalsDelta)
	}
	t.Logf("GATE-26R3 closed=%d reopened=%d events=%d mid_speech_transcripts=%d spoken_samples=%d gate_on=%d gate_off=%d",
		closed, reopened, len(events), finalsDelta, samples, rv.speakGateOn.Load(), rv.speakGateOff.Load())

	// The self-echo feed: frames arrive at the seam WHILE the gate is closed.
	// The reading is the gate's state and the refusal, not a device claim: this
	// rig has no device, and on PathT a speaking gate refuses capture
	// (effectiveOpen=false) - the frames below stand for the echo the microphone
	// would hear, and the gate being closed is what would have kept them out.
	gate.SetSpeaking(true)
	fedFrames := 0
	silent := audio.EncodeFrame(make([]int16, audio.FrameSamples))
	for i := 0; i < 6; i++ { // exactly the bounded window's capacity
		select {
		case frames <- silent:
			fedFrames++
		case <-time.After(10 * time.Second):
			t.Fatal("the seam never took the echo frames: the consumer stopped draining")
		}
	}
	voice15eWaitFor(t, "the echo frames to drain after the gate reopened", func() bool {
		return rv.drained.Load() >= gate.drainedAtOpen+uint64(fedFrames)
	})
	gate.SetSpeaking(false)
	_, _, speakingNow, _ := gate.snapshot()
	if speakingNow {
		t.Fatal("the probe window left the gate closed")
	}
	if _, fl := gate.windows(); fl != 0 {
		t.Fatalf("%d transcript(s) grew across the echo window: silence fed behind a closed gate must not transcribe", fl)
	}
	t.Logf("ECHO-26R3 gate_state_during_feed=speaking(PathT capture refused) frames_fed=%d drained_total=%d transcripts_after=%d",
		fedFrames, rv.drained.Load(), rv.finals.Load())

	// Reopen really reopens: one more real utterance is captured end to end with
	// the gate back open.
	wav2 := filepath.Join(dataDir, "turn2.wav")
	voice15eWriteWav(t, wav2, voice15eSynth(voice15eLeadMS, voice15eSpeechMS, voice15eTailMS))
	inj, err := audio.NewWavInjector(wav2)
	if err != nil {
		t.Fatalf("wav injector: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := inj.Start(ctx, frames); err != nil {
		t.Fatalf("wav injector start: %v", err)
	}
	voice15eWaitFor(t, "the reopened turn to reach the text path", func() bool {
		return rv.finals.Load() >= 2
	})
	_ = inj.Stop()
	voice15eWaitFor(t, "the second turn to wake the session", func() bool {
		return rv.warmWake.Load() >= 1
	})
	if _, _, speakingAtWake, _ := gate.snapshot(); speakingAtWake {
		t.Fatal("the gate was closed during a live utterance after the announcement ended")
	}
	t.Logf("REOPEN-26R3 finals=%d warm_wake=%d stage=%s gate_speaking=false",
		rv.finals.Load(), rv.warmWake.Load(), string(rv.residency.stage()))
}

// Test26R3GateNeedleZeroValueSafety pins the nil-safety: a boot that built no
// capture leg leaves the counters at zero and the speak body never
// dereferences anything.
func Test26R3GateNeedleZeroValueSafety(t *testing.T) {
	var nilRV *residentVoice
	if closed, reopened := nilRV.setSpeakingCount(); closed != 0 || reopened != 0 {
		t.Fatalf("a nil leg reported gate counts %d/%d", closed, reopened)
	}
	nilRV.attachCaptureGate(nil)
	rv := &residentVoice{}
	rv.attachCaptureGate(nil)
	if closed, reopened := rv.setSpeakingCount(); closed != 0 || reopened != 0 {
		t.Fatalf("a boot without a gate reported counts %d/%d", closed, reopened)
	}
	if got := rv.stopLine(); got == "" {
		t.Fatal("the stop line went empty on a leg without a scope")
	}
	t.Logf("ZEROVALUE closed=0 reopened=0 stop_line_ok=1")
}
