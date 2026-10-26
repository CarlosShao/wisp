//go:build windows

package main

// Ticket wire-r1's readings, taken through the production assembly:
//
//	#1  TTS announce end-to-end   ChTTS gate -> Speak -> the device -> the ball
//	#2  punctuation end-to-end    ASR final -> the stage -> the text path
//	#3  warm end-to-end           turn end -> Warm -> the next acquire is free
//	#4  the gate's teeth          a live taint refuses the sentence
//	#5  the positive control      silence and empty text never speak
//	#6  the idle dispose          the timer cascade lands on the owner goroutine
//
// The staging contract is 15e's (see resident_voice_15e_windows_test.go's header):
// the sherpa DLLs must sit next to the test binary and the model files must be
// reachable, so every case is run with -test.run=TestWireR1 and never as part of
// the whole package (which dies at MARK299). Nothing here re-implements the
// chain: the frames go through the real C8 seam, the engines are the real
// internal/speech ones (except the one case that needs a KNOWN transcript, which
// uses the exported AsrEngine seam for exactly the reason 27-r1 did), and the
// announce path is the production one with a recording speaker standing in for
// the sound card.

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/session"
	"github.com/CarlosShao/wisp/internal/speech"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// voiceWireSpeaker is the render seam's recorder: it keeps every sample the
// announce path produced, so "it spoke" is a reading (chunks, samples, peak, RMS)
// and not a counter somebody could fake by returning nil.
type voiceWireSpeaker struct {
	mu      sync.Mutex
	rate    int
	samples []float32
	chunks  int
	drains  int
	closes  int
	aborts  int
	playErr error
}

func (s *voiceWireSpeaker) SampleRate() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rate
}

func (s *voiceWireSpeaker) Play(samples []float32, engineRate int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.playErr != nil {
		return s.playErr
	}
	s.chunks++
	s.samples = append(s.samples, samples...)
	return nil
}

func (s *voiceWireSpeaker) Drain() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.drains++
	return nil
}

func (s *voiceWireSpeaker) Abort() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.aborts++
	return 0
}

func (s *voiceWireSpeaker) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closes++
	return nil
}

// energy reports what was really handed to the device: a silent buffer would pass
// a samples>0 assertion, which is why 26-r1 measured RMS.
func (s *voiceWireSpeaker) energy() (chunks, samples int, peak, rms float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var sum float64
	for _, v := range s.samples {
		a := math.Abs(float64(v))
		if a > peak {
			peak = a
		}
		sum += float64(v) * float64(v)
	}
	if n := len(s.samples); n > 0 {
		rms = math.Sqrt(sum / float64(n))
	}
	return s.chunks, len(s.samples), peak, rms
}

func (s *voiceWireSpeaker) written() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.samples)
}

// voiceWireAsr is the exported AsrEngine seam with a fixed transcript. It exists
// for the punctuation and announce cases, which need a KNOWN sentence: the sherpa
// fixture's transcript is filler syllables ("嗯嗯嗯嗯嗯嗯哦哦") and a punctuator has
// nothing to punctuate in it, so driving those cases with the real ASR would
// measure the fixture instead of the wiring.
type voiceWireAsr struct {
	text   string
	open   bool
	fed    int
	closes atomic.Int64
}

func (a *voiceWireAsr) Begin()           { a.open, a.fed = true, 0 }
func (a *voiceWireAsr) Feed(s []float32) { a.fed += len(s) }
func (a *voiceWireAsr) Partial() string  { return "" }
func (a *voiceWireAsr) Finish() string {
	a.open = false
	return a.text
}
func (a *voiceWireAsr) InFlight() bool { return a.open }
func (a *voiceWireAsr) Fed() int       { return a.fed }
func (a *voiceWireAsr) Close() error   { a.closes.Add(1); return nil }

// voiceWireClock is session.Clock without the waiting, so the 90 s / 3 s rows are
// exercised in microseconds (28-r1's own readings did the same). It is wrapped by
// voiceOwnerClock in the cases that fire it, which is what makes the cascade land
// on the owner goroutine exactly as it does in production.
type voiceWireClock struct {
	mu     sync.Mutex
	timers []*voiceWireTimer
}

type voiceWireTimer struct {
	mu    sync.Mutex
	fn    func()
	fired bool
}

func (c *voiceWireClock) AfterFunc(d time.Duration, fn func()) session.Timer {
	t := &voiceWireTimer{fn: fn}
	c.mu.Lock()
	c.timers = append(c.timers, t)
	c.mu.Unlock()
	return t
}

func (t *voiceWireTimer) Stop() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.fired {
		return false
	}
	t.fired = true
	return true
}

func (c *voiceWireClock) fireLatest() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := len(c.timers) - 1; i >= 0; i-- {
		t := c.timers[i]
		t.mu.Lock()
		if t.fired {
			t.mu.Unlock()
			continue
		}
		t.fired = true
		fn := t.fn
		t.mu.Unlock()
		fn()
		return true
	}
	return false
}

// voiceWirePunctModel is the CT-Transformer file, from the store-style tree the
// repository already carries (27-r1 measured this file's sha256 against the
// manifest pin) or from the same WISP_PUNC_MODEL the staging uses.
func voiceWirePunctModel(t *testing.T) string {
	t.Helper()
	if p := os.Getenv("WISP_PUNC_MODEL"); p != "" {
		return p
	}
	return voice15eRepoFile(t, filepath.Join("third_party", "model-fetch", "x-punc",
		"sherpa-onnx-punct-ct-transformer-zh-en-vocab272727-2024-04-12-int8", "model.int8.onnx"))
}

// voiceWireTtsFiles is the matcha checkpoint, from WISP_TTS_DIR when the staging
// names one and from the pinned spike-models tree otherwise.
func voiceWireTtsFiles(t *testing.T) voiceTtsFiles {
	t.Helper()
	dir := os.Getenv("WISP_TTS_DIR")
	if dir == "" {
		dir = voice15eRepoFile(t, filepath.Join("third_party", "spike-models", "matcha-icefall-zh-baker"))
	}
	files := voiceTtsFilesIn(dir)
	if !files.complete() {
		t.Fatalf("the matcha checkpoint is incomplete under %s: %+v", dir, files)
	}
	return files
}

// voiceWireStart builds and starts the leg the way the assembly root does, and
// lets the case adjust the seams (speaker, clock, ASR factory) before the thread
// is booked - the same order the production root has, where the plan is resolved
// before start().
func voiceWireStart(t *testing.T, plan voicePlan, ball *residentBall, sink voiceTranscriptSink, adjust func(*residentVoice)) (*residentVoice, chan []byte) {
	t.Helper()
	frames := audio.NewBoundedFrames()
	rv := newResidentVoice(frames, plan, ball, sink)
	// buildResidentVoice replaces the residency with one built from the resolved
	// plan; this constructor path already has the plan, so keep the constructor's
	// and only re-point it at the same queue.
	rv.residency = newVoiceResidency(plan, rv.queue, nil)
	if adjust != nil {
		adjust(rv)
	}
	rv.start(nil)
	t.Cleanup(rv.stop)
	if !rv.drainsFrames() {
		t.Fatalf("the leg did not start: verdict=%q", rv.verdict)
	}
	return rv, frames
}

// voiceWireFeed pushes one wav through the real C8 seam and returns once the
// wanted number of finals has been consumed.
func voiceWireFeed(t *testing.T, frames chan<- []byte, wav string, rv *residentVoice, wantFinals uint64) {
	t.Helper()
	inj, err := audio.NewWavInjector(wav)
	if err != nil {
		t.Fatalf("wav injector: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := inj.Start(ctx, frames); err != nil {
		t.Fatalf("wav injector start: %v", err)
	}
	defer func() { _ = inj.Stop() }()
	voice15eWaitFor(t, fmt.Sprintf("%d final(s) to reach the text path", wantFinals), func() bool {
		return rv.finals.Load() >= wantFinals
	})
}

// TestWireR1AnnounceEndToEnd is reading #1 plus the turn-end half of #3: one real
// utterance through the real ASR, the final handed to the text path, the session
// landing in Warm, and then one announcement that walks every hop the ticket asks
// for.
func TestWireR1AnnounceEndToEnd(t *testing.T) {
	dataDir := voice15eDataDir(t)
	voice15eWriteConfig(t, dataDir, nil)
	t.Setenv("WISP_MODELS_MANIFEST", voice15eRepoFile(t, filepath.Join("build", "models", "manifest.json")))
	store := voice15eInstallAsr(t, dataDir)

	speaker := &voiceWireSpeaker{rate: 22050}
	var speakerMu sync.Mutex
	prov := risk.NewProvenance(risk.ProvOptions{NoProbe: true})

	sink := &voice15eSink{}
	ball := &residentBall{}
	plan := voicePlan{
		VadModelPath:   voice15eVadModel(t),
		AsrDir:         store,
		Source:         "wire-r1 test",
		HandOff:        residentVoiceHandOff(dataDir),
		PunctModelPath: voiceWirePunctModel(t),
		Tts:            voiceWireTtsFiles(t),
	}
	rv, frames := voiceWireStart(t, plan, ball, sink, func(rv *residentVoice) {
		rv.newSpeaker = func(rate int) (voiceSpeaker, error) {
			speakerMu.Lock()
			speaker.rate = rate
			speakerMu.Unlock()
			return speaker, nil
		}
		rv.attachProvenance(prov)
	})
	voice15eArm(t, frames, rv)

	wav := filepath.Join(dataDir, "turn1.wav")
	voice15eWriteWav(t, wav, voice15eSynth(voice15eLeadMS, voice15eSpeechMS, voice15eTailMS))
	voiceWireFeed(t, frames, wav, rv, 1)

	// The turn really happened: a transcript reached the text path.
	texts := sink.texts()
	if len(texts) != 1 || strings.TrimSpace(texts[0]) == "" {
		t.Fatalf("the text path received %v (problems=%v)", texts, rv.problemsSnapshot())
	}
	if got := rv.handOffs.Load(); got != 1 {
		t.Fatalf("the model hand-off ran %d times, want 1", got)
	}
	// The punctuation stage attached at the one call site that exists.
	if !rv.punctAttached.Load() {
		t.Fatalf("the punctuation stage did not attach; problems=%v", rv.problemsSnapshot())
	}
	// The turn ended into WARM, not into a teardown.
	stage := rv.residency.stage()
	if stage != session.StageWarm && stage != session.StageSettling {
		t.Fatalf("stage after the turn = %q, want Warm (engines kept hot, ticket 28)", string(stage))
	}
	st := rv.residency.stats()
	if st.EnginesHot < 1 {
		t.Fatalf("engines hot after the turn = %d, want the ASR epoch still resident", st.EnginesHot)
	}
	t.Logf("WARM-TURN-END stage=%s engines_hot=%d loads_after_warm=%d refs=%d turn_end=%d cold_acquire=%v",
		string(st.Stage), st.EnginesHot, st.LoadsAfterWarm, st.Refs, rv.warmTurnEnd.Load(), rv.residency.acquireLog()[0])

	// Reading #2: the announcement. It is handed over from the task goroutine in
	// production; here it is the same call.
	if !rv.announce("任务已完成：明天下午三点的会已经改到四点。") {
		t.Fatalf("the announce path refused the sentence; problems=%v", rv.problemsSnapshot())
	}
	voice15eWaitFor(t, "the announcement to be spoken", func() bool { return rv.ttsSpoken.Load() == 1 })

	// HOP 1: the ChTTS gate ran and passed.
	if got := rv.ttsGatedOk.Load(); got != 1 {
		t.Fatalf("the ChTTS gate passed %d sentences, want 1", got)
	}
	if got := rv.ttsGated.Load(); got != 0 {
		t.Fatalf("the ChTTS gate refused %d sentences with no taint in scope", got)
	}
	// HOP 2: the engine came from the slot, which released the ASR first.
	if resident := rv.residency.slot.Resident(); resident != speech.FamilyTts {
		t.Fatalf("slot resident = %q after the announcement, want tts", string(resident))
	}
	if loads, closes := rv.residency.slot.Loads(), rv.residency.slot.Closes(); loads != 2 || closes != 1 {
		t.Fatalf("slot lifecycle loads=%d closes=%d, want loads=2 (asr, tts) closes=1 (the asr released before the tts loaded)",
			loads, closes)
	}
	// HOP 3 and 4: real audio, and not silence.
	chunks, samples, peak, rms := speaker.energy()
	if chunks == 0 || samples == 0 {
		t.Fatalf("the announcement emitted nothing: chunks=%d samples=%d (problems=%v)", chunks, samples, rv.problemsSnapshot())
	}
	if peak < 0.05 || rms < 0.005 {
		t.Fatalf("the announcement emitted near-silence: peak=%.4f rms=%.4f (a silent buffer would pass samples>0)", peak, rms)
	}
	if got := rv.ttsSamples.Load(); got != uint64(samples) {
		t.Fatalf("counted %d samples, the device got %d", got, samples)
	}
	// HOP 5: the ball walked Thinking -> Speaking -> Warm, and nothing else.
	states := ball.voiceStatesSnapshot()
	want := []statemachine.State{
		statemachine.StateListening, statemachine.StateThinking,
		statemachine.StateSpeaking, statemachine.StateWarm,
	}
	if len(states) != len(want) {
		t.Fatalf("the ball was projected with %v, want %v", states, want)
	}
	for i := range want {
		if states[i] != want[i] {
			t.Fatalf("the ball was projected with %v, want %v", states, want)
		}
	}
	timeline := rv.timelineSnapshot()
	if len(timeline) != 4 {
		t.Fatalf("the D43 timeline has %d edges, want 4 (%+v)", len(timeline), timeline)
	}
	if timeline[2].Event != statemachine.EvFirstToken || timeline[2].To != statemachine.StateSpeaking {
		t.Fatalf("third edge %+v, want Thinking+llm.first-token->Speaking (D43#15)", timeline[2])
	}
	if timeline[3].Event != statemachine.EvSpeakDone || timeline[3].To != statemachine.StateWarm {
		t.Fatalf("fourth edge %+v, want Speaking+tts.speak-done->Warm (D43#26)", timeline[3])
	}
	// The session landed back in Warm after speaking, with the TTS still hot.
	if st := rv.residency.stage(); st != session.StageWarm {
		t.Fatalf("stage after the announcement = %q, want Warm", string(st))
	}
	t.Logf("ANNOUNCE spoken=%d chunks=%d samples=%d peak=%.4f rms=%.4f rate=%d gate_ok=%d gate_refused=%d slot=%s loads=%d closes=%d",
		rv.ttsSpoken.Load(), chunks, samples, peak, rms, speaker.rate,
		rv.ttsGatedOk.Load(), rv.ttsGated.Load(), string(rv.residency.slot.Resident()),
		rv.residency.slot.Loads(), rv.residency.slot.Closes())
	t.Logf("ANNOUNCE-ACQUIRE cold(turn1)=%v announce(tts)=%v", rv.residency.acquireLog()[0], rv.residency.acquireLog()[1])
	t.Logf("TIMELINE %+v", timeline)

	// Positive control inside the same run: punctuation-only text is refused by the
	// ENGINE (26-r1's M2 guard), so nothing more reaches the device.
	before := speaker.written()
	if rv.announce("。，、！？") {
		voice15eWaitFor(t, "the silent announcement to be refused by the engine", func() bool {
			return rv.ttsSilent.Load() == 1
		})
	}
	if after := speaker.written(); after != before {
		t.Fatalf("a punctuation-only sentence emitted %d samples; the engine's nothing-to-say guard did not hold", after-before)
	}
	t.Logf("SILENT-TEXT silent=%d samples_after=%d", rv.ttsSilent.Load(), speaker.written())
}

// TestWireR1WarmSecondTurnAcquiresNothing is reading #3: the second turn's engine
// is already hot, and "hot" is proved twice - the acquire's own duration and the
// loader's call count, because a timing can be lucky and a counter cannot.
func TestWireR1WarmSecondTurnAcquiresNothing(t *testing.T) {
	dataDir := voice15eDataDir(t)
	voice15eWriteConfig(t, dataDir, nil)
	t.Setenv("WISP_MODELS_MANIFEST", voice15eRepoFile(t, filepath.Join("build", "models", "manifest.json")))
	store := voice15eInstallAsr(t, dataDir)

	sink := &voice15eSink{}
	ball := &residentBall{}
	plan := voicePlan{
		VadModelPath: voice15eVadModel(t),
		AsrDir:       store,
		Source:       "wire-r1 test",
		HandOff:      residentVoiceHandOff(dataDir),
	}
	rv, frames := voiceWireStart(t, plan, ball, sink, func(rv *residentVoice) {
		rv.newSpeaker = func(rate int) (voiceSpeaker, error) { return &voiceWireSpeaker{rate: rate}, nil }
	})
	voice15eArm(t, frames, rv)

	wav1 := filepath.Join(dataDir, "turn1.wav")
	voice15eWriteWav(t, wav1, voice15eSynth(voice15eLeadMS, voice15eSpeechMS, voice15eTailMS))
	wav2 := filepath.Join(dataDir, "turn2.wav")
	voice15eWriteWav(t, wav2, voice15eSynth(voice15eLeadMS, voice15eSpeechMS, voice15eTailMS))

	voiceWireFeed(t, frames, wav1, rv, 1)
	// The turn ended into Warm before the second utterance starts.
	voice15eWaitFor(t, "the first turn to land in Warm", func() bool {
		return rv.warmTurnEnd.Load() == 1
	})
	firstAcquire := rv.residency.acquireLog()[0]
	if firstAcquire <= 0 {
		t.Fatalf("the first acquire was %v: a cold ASR load cannot take no time", firstAcquire)
	}
	afterFirst := rv.residency.loaderCallCount()
	if afterFirst != 1 {
		t.Fatalf("the ASR factory ran %d times after one turn, want 1", afterFirst)
	}

	voiceWireFeed(t, frames, wav2, rv, 2)
	voice15eWaitFor(t, "the second turn to wake the session", func() bool {
		return rv.warmWake.Load() >= 1
	})
	log := rv.residency.acquireLog()
	if len(log) < 2 {
		t.Fatalf("only %d acquires were recorded, want the cold one and the warm one (%v)", len(log), log)
	}
	warmAcquire := log[1]
	if got := rv.residency.loaderCallCount(); got != 1 {
		t.Fatalf("the ASR factory ran %d times over two turns, want 1: the warm path reloaded an engine (B4)", got)
	}
	if st := rv.residency.stats(); st.LoadsAfterWarm != 0 {
		t.Fatalf("LoadsAfterWarm = %d after a Warm wake, want 0 (B4's claim, from the session's own counter)", st.LoadsAfterWarm)
	}
	if st := rv.residency.stats(); st.WarmResumes < 1 {
		t.Fatalf("WarmResumes = %d, want at least 1: the second turn did not wake the session", st.WarmResumes)
	}
	if closes := rv.residency.micCloseCount(); closes < 1 {
		t.Fatalf("the Warm transition never closed the mic half (mic.keep-off), closes=%d", closes)
	}
	if warmAcquire > firstAcquire/10 {
		t.Fatalf("the warm acquire took %v against the cold %v: the second turn paid for a load", warmAcquire, firstAcquire)
	}
	t.Logf("WARM-TWO-TURNS cold_acquire=%v warm_acquire=%v loader_calls=%d loads_after_warm=%d warm_resumes=%d mic_closes=%d stages_at_mic_close=%v",
		firstAcquire, warmAcquire, rv.residency.loaderCallCount(), rv.residency.stats().LoadsAfterWarm,
		rv.residency.stats().WarmResumes, rv.residency.micCloseCount(), rv.residency.micCloseStages())
	t.Logf("WARM-TWO-TURNS finals=%d partials=%d handed=%d warm_turn_end=%d warm_wake=%d drained=%d",
		rv.finals.Load(), rv.partials.Load(), rv.handed.Load(), rv.warmTurnEnd.Load(), rv.warmWake.Load(), rv.drained.Load())
}

// TestWireR1PunctuationReachesTheAgent is reading #2: the SAME text with the stage
// detached and with it attached. The "before" half is taken by reading the raw
// engine string out of the fake engine, so the two halves differ in exactly one
// thing - whether SetPunctuator was called.
func TestWireR1PunctuationReachesTheAgent(t *testing.T) {
	const raw = "今天下午三点提醒我开会"
	dataDir := voice15eDataDir(t)
	voice15eWriteConfig(t, dataDir, nil)
	t.Setenv("WISP_MODELS_MANIFEST", voice15eRepoFile(t, filepath.Join("build", "models", "manifest.json")))
	store := voice15eInstallAsr(t, dataDir)

	sink := &voice15eSink{}
	ball := &residentBall{}
	asr := &voiceWireAsr{text: raw}
	plan := voicePlan{
		VadModelPath:   voice15eVadModel(t),
		AsrDir:         store,
		Source:         "wire-r1 test",
		HandOff:        residentVoiceHandOff(dataDir),
		PunctModelPath: voiceWirePunctModel(t),
	}
	rv, frames := voiceWireStart(t, plan, ball, sink, func(rv *residentVoice) {
		rv.residency.asrLoader = func() (speech.AsrEngine, error) { return asr, nil }
		rv.newSpeaker = func(rate int) (voiceSpeaker, error) { return &voiceWireSpeaker{rate: rate}, nil }
	})
	voice15eArm(t, frames, rv)

	wav := filepath.Join(dataDir, "dictation.wav")
	voice15eWriteWav(t, wav, voice15eSynth(voice15eLeadMS, voice15eSpeechMS, voice15eTailMS))
	voiceWireFeed(t, frames, wav, rv, 1)

	got := sink.texts()
	if len(got) != 1 {
		t.Fatalf("the text path received %v, want exactly one final (problems=%v)", got, rv.problemsSnapshot())
	}
	if !rv.punctAttached.Load() {
		t.Fatalf("the stage never attached, so this case cannot say anything about punctuation: %v", rv.problemsSnapshot())
	}
	if got[0] == raw {
		t.Fatalf("the final reached the agent UNPUNCTUATED (%q) with a stage attached: SetPunctuator is not on the path", got[0])
	}
	if got[0] == voiceWireStripPunct(got[0]) {
		t.Fatalf("the final %q carries no punctuation at all", got[0])
	}
	if stripped := voiceWireStripPunct(got[0]); stripped != raw {
		t.Fatalf("the stage changed the words: %q stripped is %q, the raw transcript is %q", got[0], stripped, raw)
	}
	if rv.punctFinals.Load() != 1 {
		t.Fatalf("finals reporting Punctuated = %d, want 1: downstream cannot tell dictation from a raw hypothesis", rv.punctFinals.Load())
	}
	turns := rv.turnsSnapshot()
	if len(turns) != 1 || turns[0].Text != got[0] {
		t.Fatalf("the turn record %+v does not match what the agent got %q", turns, got[0])
	}
	t.Logf("PUNCTUATION raw=%q punctuated=%q punctuated_flag=%d attached=%v",
		raw, got[0], rv.punctFinals.Load(), rv.punctAttached.Load())
}

// TestWireR1AnnounceGateRefusesTaintedText is reading #4: the teeth. A live taint
// in the announce scope must stop the sentence BEFORE synthesis; the reverse
// criterion (mutation A removes the CheckText call) turns this case red.
func TestWireR1AnnounceGateRefusesTaintedText(t *testing.T) {
	dataDir := voice15eDataDir(t)
	voice15eWriteConfig(t, dataDir, nil)
	t.Setenv("WISP_MODELS_MANIFEST", voice15eRepoFile(t, filepath.Join("build", "models", "manifest.json")))
	store := voice15eInstallAsr(t, dataDir)

	speaker := &voiceWireSpeaker{rate: 22050}
	prov := risk.NewProvenance(risk.ProvOptions{NoProbe: true})
	const secret = "AKIAIOSFODNN7EXAMPLE"

	sink := &voice15eSink{}
	ball := &residentBall{}
	plan := voicePlan{
		VadModelPath: voiceVadNone(),
		AsrDir:       store,
		Source:       "wire-r1 test",
		Tts:          voiceWireTtsFiles(t),
	}
	// No frames at all: this case is about the gate, not about the microphone, and
	// a leg with no frame channel still owns an announce path.
	frames := make(chan []byte)
	rv := newResidentVoice(frames, plan, ball, sink)
	rv.residency = newVoiceResidency(plan, rv.queue, nil)
	rv.newSpeaker = func(rate int) (voiceSpeaker, error) { return speaker, nil }
	// Attach FIRST, then mark: the taint has to land in the scope the leg actually
	// reads, or the case would prove a scope mismatch instead of the gate.
	rv.attachProvenance(prov)
	if !prov.Mark(voiceAnnounceScope, risk.SrcFSRead, `D:\work\secret.txt`, "凭据是 "+secret) {
		t.Fatal("Mark refused the taint")
	}
	t.Cleanup(rv.stop)

	if !rv.announce("读到的内容：凭据是 " + secret) {
		t.Fatalf("the queue refused the sentence, so the gate was never asked: %v", rv.problemsSnapshot())
	}
	// The queue's drain is what the owner goroutine runs from, and running it on
	// this goroutine is the same call the owner makes (no goroutine of the leg's is
	// involved - the leg starts none, D38b).
	rv.queue.run()
	voice15eWaitFor(t, "the gate to refuse the sentence", func() bool {
		return rv.ttsGated.Load() == 1
	})
	if got := rv.ttsGatedOk.Load(); got != 0 {
		t.Fatalf("the gate PASSED %d tainted sentences; the ChTTS check is not on the path", got)
	}
	if rv.ttsSpoken.Load() != 0 || speaker.written() != 0 {
		t.Fatalf("a tainted sentence reached the device: spoken=%d samples=%d", rv.ttsSpoken.Load(), speaker.written())
	}
	if got := rv.residency.slot.Loads(); got != 0 {
		t.Fatalf("the TTS engine loaded %d times for a refused sentence; the gate must run BEFORE synthesis", got)
	}
	found := false
	for _, p := range rv.problemsSnapshot() {
		if strings.Contains(p, string(risk.ChTTS)) {
			found = true
		}
	}
	if !found {
		t.Fatalf("the refusal was not recorded with its channel: %v", rv.problemsSnapshot())
	}
	t.Logf("GATE refused=%d passed=%d spoken=%d samples=%d slot_loads=%d",
		rv.ttsGated.Load(), rv.ttsGatedOk.Load(), rv.ttsSpoken.Load(), speaker.written(), rv.residency.slot.Loads())
}

// TestWireR1SilenceAndEmptyTextNeverSpeak is reading #5, the positive control:
// silence produces no announcement, empty text is refused before anything is
// loaded, and a punctuation-only sentence is refused by the ENGINE rather than
// played as a near-silent blip.
func TestWireR1SilenceAndEmptyTextNeverSpeak(t *testing.T) {
	dataDir := voice15eDataDir(t)
	voice15eWriteConfig(t, dataDir, nil)
	t.Setenv("WISP_MODELS_MANIFEST", voice15eRepoFile(t, filepath.Join("build", "models", "manifest.json")))
	store := voice15eInstallAsr(t, dataDir)

	speaker := &voiceWireSpeaker{rate: 22050}
	sink := &voice15eSink{}
	ball := &residentBall{}
	plan := voicePlan{
		VadModelPath: voice15eVadModel(t),
		AsrDir:       store,
		Source:       "wire-r1 test",
		HandOff:      residentVoiceHandOff(dataDir),
		Tts:          voiceWireTtsFiles(t),
	}
	rv, frames := voiceWireStart(t, plan, ball, sink, func(rv *residentVoice) {
		rv.newSpeaker = func(rate int) (voiceSpeaker, error) { return speaker, nil }
		rv.attachProvenance(nil)
	})
	voice15eArm(t, frames, rv)

	// Silence only: 120 frames, pushed the way the capture leg pushes them.
	silent := audio.EncodeFrame(make([]int16, audio.FrameSamples))
	for i := 0; i < 120; i++ {
		select {
		case frames <- silent:
		case <-time.After(10 * time.Second):
			t.Fatal("the leg stopped draining during the silence run")
		}
	}
	voice15eWaitFor(t, "the silence run to be drained", func() bool { return rv.drained.Load() >= 121 })
	if rv.finals.Load() != 0 || rv.handed.Load() != 0 {
		t.Fatalf("silence produced a final: finals=%d handed=%d", rv.finals.Load(), rv.handed.Load())
	}
	if rv.announceAsked.Load() != 0 || rv.ttsSpoken.Load() != 0 || speaker.written() != 0 {
		t.Fatalf("silence produced an announcement: asked=%d spoken=%d samples=%d",
			rv.announceAsked.Load(), rv.ttsSpoken.Load(), speaker.written())
	}

	// Empty and whitespace: refused before any engine is touched.
	for _, text := range []string{"", "   ", "\t\n"} {
		if rv.announce(text) {
			t.Fatalf("the announce path accepted %q", text)
		}
	}
	if got := rv.announceEmpty.Load(); got != 3 {
		t.Fatalf("empty announcements counted %d, want 3", got)
	}
	if got := rv.residency.slot.Loads(); got != 0 {
		t.Fatalf("an empty announcement loaded %d engines", got)
	}

	// Punctuation-only: non-empty, so it reaches the engine, and the ENGINE's own
	// nothing-to-say guard is what stops it (26-r1's M2 case).
	if !rv.announce("。，、！？") {
		t.Fatalf("the punctuation-only sentence was refused by the leg, not by the engine")
	}
	voice15eWaitFor(t, "the engine to refuse the punctuation-only sentence", func() bool {
		return rv.ttsSilent.Load() == 1
	})
	if speaker.written() != 0 {
		t.Fatalf("a punctuation-only sentence emitted %d samples", speaker.written())
	}
	if rv.ttsSpoken.Load() != 0 {
		t.Fatalf("a punctuation-only sentence was counted as spoken")
	}
	t.Logf("POSITIVE-CONTROL finals=%d handed=%d announce_asked=%d empty=%d silent=%d spoken=%d samples=%d",
		rv.finals.Load(), rv.handed.Load(), rv.announceAsked.Load(), rv.announceEmpty.Load(),
		rv.ttsSilent.Load(), rv.ttsSpoken.Load(), speaker.written())
}

// TestWireR1WarmIdleDisposesOnTheOwner is the other half of #3: the 90 s idle row
// and the 3 s settle window really end the session, the cascade closes the engines
// exactly once, and - the part that is a crash if it is wrong - it does so on the
// voice owner goroutine, never on the runtime timer goroutine.
func TestWireR1WarmIdleDisposesOnTheOwner(t *testing.T) {
	const raw = "打开收音机"
	dataDir := voice15eDataDir(t)
	voice15eWriteConfig(t, dataDir, nil)
	t.Setenv("WISP_MODELS_MANIFEST", voice15eRepoFile(t, filepath.Join("build", "models", "manifest.json")))
	store := voice15eInstallAsr(t, dataDir)

	clock := &voiceWireClock{}
	asr := &voiceWireAsr{text: raw}
	sink := &voice15eSink{}
	ball := &residentBall{}
	plan := voicePlan{
		VadModelPath: voice15eVadModel(t),
		AsrDir:       store,
		Source:       "wire-r1 test",
		HandOff:      residentVoiceHandOff(dataDir),
	}
	rv, frames := voiceWireStart(t, plan, ball, sink, func(rv *residentVoice) {
		rv.residency.asrLoader = func() (speech.AsrEngine, error) { return asr, nil }
		rv.residency.clock = voiceOwnerClock{queue: rv.queue, base: clock}
		rv.newSpeaker = func(rate int) (voiceSpeaker, error) { return &voiceWireSpeaker{rate: rate}, nil }
	})
	voice15eArm(t, frames, rv)
	wav := filepath.Join(dataDir, "idle.wav")
	voice15eWriteWav(t, wav, voice15eSynth(voice15eLeadMS, voice15eSpeechMS, voice15eTailMS))
	voiceWireFeed(t, frames, wav, rv, 1)
	voice15eWaitFor(t, "the turn to land in Warm", func() bool { return rv.warmTurnEnd.Load() == 1 })
	if st := rv.residency.stage(); st != session.StageWarm {
		t.Fatalf("stage = %q, want Warm", string(st))
	}

	// The 90 s row: Warm -> Settling, engines still hot.
	if !clock.fireLatest() {
		t.Fatal("no warm timer was armed")
	}
	voice15eWaitFor(t, "the warm idle row to land in Settling", func() bool {
		return rv.residency.stage() == session.StageSettling
	})
	if st := rv.residency.stats(); st.EnginesHot < 1 {
		t.Fatalf("Settling unloaded the engines (engines_hot=%d); the 3 s window is still hot", st.EnginesHot)
	}
	// The 3 s row: Settling -> dispose, and the cascade closes the engines.
	if !clock.fireLatest() {
		t.Fatal("no settle timer was armed")
	}
	voice15eWaitFor(t, "the settle window to dispose the session", func() bool {
		return rv.residency.stage() == session.StageDisposed
	})
	voice15eWaitFor(t, "the engine to be released", func() bool { return asr.closes.Load() >= 1 })
	if st := rv.residency.stats(); st.EnginesHot != 0 {
		t.Fatalf("engines hot after dispose = %d, want 0", st.EnginesHot)
	}
	if got := rv.residency.micCloseCount(); got < 2 {
		t.Fatalf("mic close counted %d times, want at least 2 (the Warm row and the cascade's defence in depth)", got)
	}
	if rv.drainsFrames() != true {
		t.Fatalf("the leg stopped draining, so this case cannot separate the dispose from a crash")
	}
	t.Logf("WARM-IDLE stage=%s engines_hot=%d mic_closes=%d stages_at_mic_close=%v engine_closes=%d",
		string(rv.residency.stage()), rv.residency.stats().EnginesHot, rv.residency.micCloseCount(),
		rv.residency.micCloseStages(), asr.closes.Load())
}

// voiceVadNone is the empty VAD path: a case that drives the announce path with no
// microphone has no detector to open, and an empty path is how start() says so.
func voiceVadNone() string { return "" }

// voiceWireStripPunct removes the four marks this checkpoint's alphabet can add
// (punctuationAlphabet in internal/speech/punctuate.go). It is how the case says
// "the words were not touched" without depending on WHERE a comma landed.
func voiceWireStripPunct(s string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case '。', '，', '？', '、':
			return -1
		}
		return r
	}, s)
}

// TestWireR1TheRealDeviceTakesTheAudio is the other half of reading #1: the four
// hops above prove the samples leave the engine with real energy, and this one
// proves the renderer they leave through is the machine's own device and not dead
// code. It is skipped where there is no output device (a CI box has none), which is
// the same discipline internal/audio's live capture case uses.
func TestWireR1TheRealDeviceTakesTheAudio(t *testing.T) {
	sp, err := openWaveOutSpeaker(22050)
	if err != nil {
		t.Skipf("no output device on this machine: %v", err)
	}
	defer func() { _ = sp.Close() }()
	if rate := sp.SampleRate(); rate <= 0 {
		t.Fatalf("the speaker opened at rate %d", rate)
	}
	// 200 ms of a 440 Hz tone: audible, and short enough that a case nobody is
	// watching does not hold the device.
	const ms = 200
	n := sp.SampleRate() * ms / 1000
	tone := make([]float32, n)
	for i := range tone {
		tone[i] = float32(0.25 * math.Sin(2*math.Pi*440*float64(i)/float64(sp.SampleRate())))
	}
	if err := sp.Play(tone, sp.SampleRate()); err != nil {
		t.Fatalf("the device refused the buffer: %v", err)
	}
	if err := sp.Drain(); err != nil {
		t.Fatalf("the device did not finish the buffer: %v", err)
	}
	t.Logf("DEVICE rate=%d samples=%d ms=%d abort=%d", sp.SampleRate(), n, ms, sp.Abort())
}
