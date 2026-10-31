//go:build windows

package main

// Ticket 41 A2's readings, taken through the production wiring the same way
// 15e/wire-r1 take theirs: the frames go through the real C8 seam
// (audio.WavInjector), the wake-word engine is the real internal/speech
// KwsEngine, and the chain behind the wake is the real hand-off walk plus the
// real streaming ASR. The staging contract is 15e's: the sherpa DLLs sit next
// to the test binary and the model files are reachable, so every engine case is
// run with -test.run=Test41a2r1 and never as part of the whole package (which
// dies at MARK299).
//
// The fixtures are the checkpoint's own test_wavs, and the pairing fact the
// 41-a2 engine leg measured carries over UNCHANGED: 3-6.wav are recordings of
// the demo keyword table (test_keywords.txt) and 0-2.wav are that collection's
// own NEGATIVE samples - the production table (keywords.txt) has no recording
// at all, so "every fixture must hit" would be a lie about the model. The wake
// case below feeds 0.wav first and asserts silence out of the engine, then
// 4.wav and asserts the hit; the Warm case repeats the pairing in the Warm
// posture, which is where the R28-2 landing lives.

import (
	"context"
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/session"
	"github.com/CarlosShao/wisp/internal/speech"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// voice41a2r1KwsDir is the pinned KWS checkpoint directory: WISP_KWS_DIR when
// the staging names one, otherwise the repository's dev copy (the same walk-up
// every other model helper here uses).
func voice41a2r1KwsDir(t *testing.T) string {
	t.Helper()
	if p := os.Getenv("WISP_KWS_DIR"); p != "" {
		return p
	}
	return voice15eRepoFile(t, filepath.Join("third_party", "spike-models",
		"sherpa-onnx-kws-zipformer-wenetspeech-3.3M-2024-01-01"))
}

// voice41a2r1Fixture is one of the checkpoint's own wav fixtures.
func voice41a2r1Fixture(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join(voice41a2r1KwsDir(t), "test_wavs", name)
}

// voice41a2r1WavFrames is how many C8 frames the injector will push for this
// wav: the data chunk's sample count (16k mono enforced - the sherpa fixtures
// are 16k mono, and a resampled fixture would change the reading) ceiled over
// the seam's window. The chunk walk is real, not fixed offsets: 0.wav carries a
// LIST chunk before `data`, and a reader that assumes offset 44 reads the LIST
// size and starves the fixture.
func voice41a2r1WavFrames(t *testing.T, wav string) int {
	t.Helper()
	samples := voice41a2r1WavSamples(t, wav)
	return int((uint32(samples) + audio.FrameSamples - 1) / audio.FrameSamples)
}

// voice41a2r1WavSamples parses a canonical RIFF/WAVE chunk walk and returns the
// data chunk's frame count for the one shape the seam consumes (16k mono PCM16).
func voice41a2r1WavSamples(t *testing.T, wav string) int {
	t.Helper()
	data, err := os.ReadFile(wav)
	if err != nil {
		t.Fatalf("read %s: %v", wav, err)
	}
	if len(data) < 12 || string(data[0:4]) != "RIFF" || string(data[8:12]) != "WAVE" {
		t.Fatalf("%s is not a RIFF/WAVE file", wav)
	}
	rate, chans, bits, samples := 0, 0, 0, -1
	for off := 12; off+8 <= len(data); {
		id := string(data[off : off+4])
		size := int(binary.LittleEndian.Uint32(data[off+4 : off+8]))
		body := off + 8
		switch id {
		case "fmt ":
			chans = int(binary.LittleEndian.Uint16(data[body+2 : body+4]))
			rate = int(binary.LittleEndian.Uint32(data[body+4 : body+8]))
			bits = int(binary.LittleEndian.Uint16(data[body+14 : body+16]))
		case "data":
			n := size
			if body+n > len(data) {
				n = len(data) - body
			}
			samples = n / (bits / 8 * chans)
		}
		off = body + size + (size & 1)
	}
	if samples < 0 {
		t.Fatalf("%s carries no data chunk", wav)
	}
	if rate != audio.TargetRate || chans != 1 || bits != 16 {
		t.Fatalf("%s is %d Hz x %d ch x %d bit, want %d Hz mono 16-bit: the fixture must reach the seam unconverted",
			wav, rate, chans, bits, audio.TargetRate)
	}
	return samples
}

// voice41a2r1Start builds and starts the leg the way the assembly root does,
// with the wake-word plan resolved off the production spelling helper and the
// ASR store installed under the manifest's names. It asserts the ARM the boot
// walk produced (D43#5) before returning: a leg that stayed Sleeping has no
// business in a wake case.
func voice41a2r1Start(t *testing.T, sink voiceTranscriptSink) (*residentVoice, chan []byte) {
	t.Helper()
	dataDir := voice15eDataDir(t)
	voice15eWriteConfig(t, dataDir, nil)
	t.Setenv("WISP_MODELS_MANIFEST", voice15eRepoFile(t, filepath.Join("build", "models", "manifest.json")))
	store := voice15eInstallAsr(t, dataDir)

	kws := voiceKwsFilesFrom(voice41a2r1KwsDir(t))
	if !kws.complete() {
		t.Fatalf("the KWS manifest spelling did not resolve under %s", voice41a2r1KwsDir(t))
	}
	// The keyword table here is the DEMO table, not the production one, and that
	// is the pairing fact the 41-a2 engine leg measured: test_wavs are recordings
	// of test_keywords.txt's words, and NO fixture recording says any of the
	// production table's eight wake words - the production table scoring zero on
	// every fixture is the pinned control (internal/speech's
	// TestKwsEngineProductionTableRejectsTheDemoWavs), not a defect. The wake
	// cases need a table whose words the fixtures actually say; the production
	// spelling itself is covered by the resolution case below and by the
	// pipeline-level control in Test41a2r1ProductionTableRejectsTheDemoWavs.
	kws.Keywords = filepath.Join(voice41a2r1KwsDir(t), "test_wavs", "test_keywords.txt")
	plan := voicePlan{
		VadModelPath:   voice15eVadModel(t),
		AsrDir:         store,
		Source:         "41a2r1 test",
		HandOff:        residentVoiceHandOff(dataDir),
		PunctModelPath: voiceWirePunctModel(t),
		Tts:            voiceWireTtsFiles(t),
		Kws:            kws,
	}
	frames := audio.NewBoundedFrames()
	rv := newResidentVoice(frames, plan, &residentBall{}, sink)
	rv.residency = newVoiceResidency(plan, rv.queue, nil)
	rv.start(nil)
	t.Cleanup(rv.stop)
	if !rv.drainsFrames() {
		t.Fatalf("the leg did not start: verdict=%q", rv.verdict)
	}
	if !rv.kwsLoaded() {
		t.Fatalf("the wake-word engine did not arm; problems=%v", rv.problemsSnapshot())
	}
	if st := rv.machine.State(); st != statemachine.StateArmed {
		t.Fatalf("machine after the boot walk = %s, want Armed (D43#5)", string(st))
	}
	if !voice41a2r1Walked(t, rv, statemachine.StateSleeping, statemachine.EvKwsEnabled, statemachine.StateArmed) {
		t.Fatalf("the timeline has no D43#5 Sleeping->Armed walk: %v", rv.timelineSnapshot())
	}
	return rv, frames
}

// voice41a2r1Feed plays one wav through the real C8 seam and waits until the
// seam clock shows its frames drained, so the caller's counter snapshots are
// taken at a known drain position.
func voice41a2r1Feed(t *testing.T, frames chan<- []byte, wav string, rv *residentVoice) {
	t.Helper()
	inj, err := audio.NewWavInjector(wav)
	if err != nil {
		t.Fatalf("wav injector %s: %v", wav, err)
	}
	drained0 := rv.drained.Load()
	want := voice41a2r1WavFrames(t, wav)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := inj.Start(ctx, frames); err != nil {
		t.Fatalf("wav injector start: %v", err)
	}
	defer func() { _ = inj.Stop() }()
	voice15eWaitFor(t, "the wav to drain ("+filepath.Base(wav)+")", func() bool {
		return rv.drained.Load() >= drained0+uint64(want)
	})
}

// voice41a2r1PushSilence pushes n silent frames through the seam and waits for
// them to drain - the deterministic half of the Warm-gate reading (the counters
// are compared between two known drain positions).
func voice41a2r1PushSilence(t *testing.T, frames chan<- []byte, rv *residentVoice, n int) {
	t.Helper()
	drained0 := rv.drained.Load()
	silent := audio.EncodeFrame(make([]int16, audio.FrameSamples))
	for i := 0; i < n; i++ {
		select {
		case frames <- silent:
		case <-time.After(10 * time.Second):
			t.Fatal("the leg never took a frame: the consumer loop is not draining")
		}
	}
	voice15eWaitFor(t, "the silence frames to drain", func() bool {
		return rv.drained.Load() >= drained0+uint64(n)
	})
}

// voice41a2r1Walked reports whether the timeline carries the one transition the
// reading names.
func voice41a2r1Walked(t *testing.T, rv *residentVoice, from statemachine.State, ev statemachine.Event, to statemachine.State) bool {
	t.Helper()
	for _, tr := range rv.timelineSnapshot() {
		if tr.From == from && tr.Event == ev && tr.To == to {
			return true
		}
	}
	return false
}

// Test41a2r1WakeEndToEndArmedToListening is reading 2: the Armed machine hears
// the demo fixture through the real seam, the engine answers KeywordHit, the
// dispatch walks D43#7 Armed->Listening, the ball follows, and the chain the
// row's side effects name really loads and really picks audio up.
func Test41a2r1WakeEndToEndArmedToListening(t *testing.T) {
	sink := &voice15eSink{}
	rv, frames := voice41a2r1Start(t, sink)

	// The negative half of the pairing, IN the real pipeline: 0.wav is the
	// checkpoint collection's own negative sample. In Armed the engine eats it
	// (the feed counter moves) and no hit fires - the fixture stays a negative
	// here exactly as it was for 41-a2, and the machine stays Armed with no
	// chain started.
	fedBefore := rv.kwsFedSamples.Load()
	voice41a2r1Feed(t, frames, voice41a2r1Fixture(t, "0.wav"), rv)
	if rv.kwsFedSamples.Load() <= fedBefore {
		t.Fatalf("the engine was not fed in Armed (kws_fed=%d before, %d after): the Armed gate is closed",
			fedBefore, rv.kwsFedSamples.Load())
	}
	if got := rv.kwsHits.Load(); got != 0 {
		t.Fatalf("0.wav produced %d hits in the live pipeline, want 0 (it is the collection's own negative)", got)
	}
	if st := rv.machine.State(); st != statemachine.StateArmed {
		t.Fatalf("machine after the negative fixture = %s, want Armed (no legal row out of Armed without a hit)", string(st))
	}
	if rv.handOffs.Load() != 0 {
		t.Fatalf("the chain started without a wake (handoffs=%d)", rv.handOffs.Load())
	}

	// The wake: 4.wav says two demo keywords. The hit must walk #7 and load the
	// chain the row's side effects name.
	hitsBefore := rv.kwsHits.Load()
	voice41a2r1Feed(t, frames, voice41a2r1Fixture(t, "4.wav"), rv)
	voice15eWaitFor(t, "the wake-word hit", func() bool { return rv.kwsHits.Load() > hitsBefore })
	fedAtHit := rv.kwsFedSamples.Load()
	voice15eWaitFor(t, "the D43#7 Armed->Listening walk", func() bool {
		return voice41a2r1Walked(t, rv, statemachine.StateArmed, statemachine.EvWakeWord, statemachine.StateListening)
	})
	if st := rv.machine.State(); st != statemachine.StateListening {
		t.Fatalf("machine after the hit = %s, want Listening", string(st))
	}
	voice15eWaitFor(t, "the ASR chain the row loads (asr_loaded)", func() bool { return rv.asrLoaded.Load() })
	if got := rv.handOffs.Load(); got != 1 {
		t.Fatalf("the chain load ran %d times, want 1", got)
	}
	// The ball followed the walk.
	sawArmed, sawListening := false, false
	for _, s := range rv.ball.voiceStatesSnapshot() {
		if s == statemachine.StateArmed {
			sawArmed = true
		}
		if s == statemachine.StateListening {
			sawListening = true
		}
	}
	if !sawArmed || !sawListening {
		t.Fatalf("the ball saw %v, want Armed then Listening", rv.ball.voiceStatesSnapshot())
	}
	// ASR 开拾: the switchover replayed the ring into the real transcriber, so
	// the chain produces hypotheses (partials; the demo audio also produces a
	// final - the case accepts either, because "the ASR picked up" is the claim
	// and the transcript's TEXT is the model's business, not the wiring's).
	voice15eWaitFor(t, "the ASR to pick up (partials+finals >= 1)", func() bool {
		return rv.partials.Load()+rv.finals.Load() >= 1
	})
	t.Logf("KWS-WAKE kws_fed=%d asr_fed=%d hits=%d handoff_exit=%d drained=%d partials=%d finals=%d machine=%s",
		fedAtHit, rv.asrFedSamples.Load(), rv.kwsHits.Load(), rv.handOffExit.Load(),
		rv.drained.Load(), rv.partials.Load(), rv.finals.Load(), string(rv.machine.State()))
	for _, tr := range rv.timelineSnapshot() {
		t.Logf("KWS-TIMELINE %s + %s -> %s at_sample=%d (%s)", string(tr.From), string(tr.Event), string(tr.To), tr.AtSample, tr.Say)
	}
}

// Test41a2r1WarmGateFeedsKwsNotTheAsr is reading 3, the R28-2 landing: in the
// Warm posture the feed gate is closed for the text path and OPEN for the wake
// word - the KWS count advances while the ASR count stands still, and a hit in
// Warm reopens the gate through the zero-load resume.
func Test41a2r1WarmGateFeedsKwsNotTheAsr(t *testing.T) {
	sink := &voice15eSink{}
	rv, frames := voice41a2r1Start(t, sink)

	// Turn 1: the wake (4.wav) loads the chain, and trailing silence closes the
	// turn - the final walks #11 and the session lands in Warm with the feed
	// closed (the MicClose hook).
	voice41a2r1Feed(t, frames, voice41a2r1Fixture(t, "4.wav"), rv)
	voice15eWaitFor(t, "the wake-word hit", func() bool { return rv.kwsHits.Load() > 0 })
	voice15eWaitFor(t, "the ASR chain to load", func() bool { return rv.asrLoaded.Load() })
	voice41a2r1PushSilence(t, frames, rv, 90) // ~2.9 s of tail: past silero's overshoot
	voice15eWaitFor(t, "the turn to close into Warm", func() bool {
		return rv.warmFeed.Load() && rv.residency.stage() == session.StageWarm
	})
	voice15eWaitFor(t, "the final to reach the text path", func() bool { return rv.finals.Load() >= 1 })

	// The gate, measured: silence frames in Warm advance the KWS count and leave
	// the ASR count frozen. Forty frames is ~1.3 s - enough to be unambiguous,
	// far under the 90 s warm window.
	asrFrozen := rv.asrFedSamples.Load()
	kwsBefore := rv.kwsFedSamples.Load()
	voice41a2r1PushSilence(t, frames, rv, 40)
	if got := rv.kwsFedSamples.Load(); got != kwsBefore+40*audio.FrameSamples {
		t.Fatalf("kws_fed moved %d, want exactly %d: the KWS is not the Warm ears",
			got-kwsBefore, 40*audio.FrameSamples)
	}
	if got := rv.asrFedSamples.Load(); got != asrFrozen {
		t.Fatalf("asr_fed moved %d during Warm (want 0): the feed gate did not stop the ASR - R28-2 is open",
			got-asrFrozen)
	}
	// Still exactly one hit (the wake only - silence stays silent), and the gate
	// stays closed.
	if got := rv.kwsHits.Load(); got != 1 {
		t.Fatalf("silence produced extra hits (%d total), want 1", got)
	}
	if !rv.warmFeed.Load() {
		t.Fatal("the Warm gate opened without a wake")
	}

	// The Warm wake: the same fixture again, through the closed gate. The hit
	// reopens the feed (the wire-r1 wake machinery, now reached from the KWS),
	// the session resumes without loading anything, and the ASR count moves
	// again once frames flow past the reopened gate. The ball stays wherever the
	// frozen table left it (the wire-r1 named Thinking-parking) - the wake's
	// FUNCTIONAL half is the gate and the resume, and that is what is asserted.
	warmWakeBefore := rv.warmWake.Load()
	voice41a2r1Feed(t, frames, voice41a2r1Fixture(t, "4.wav"), rv)
	voice15eWaitFor(t, "the second hit", func() bool { return rv.kwsHits.Load() >= 2 })
	voice15eWaitFor(t, "the Warm wake to reopen the feed", func() bool {
		return rv.warmWake.Load() > warmWakeBefore && !rv.warmFeed.Load()
	})
	st := rv.residency.stats()
	if st.LoadsAfterWarm != 0 {
		t.Fatalf("loads_after_warm=%d after the Warm wake, want 0 (B4)", st.LoadsAfterWarm)
	}
	voice41a2r1PushSilence(t, frames, rv, 8)
	if got := rv.asrFedSamples.Load(); got == asrFrozen {
		t.Fatal("asr_fed did not move after the wake reopened the gate")
	}
	t.Logf("KWS-WARM-GATE kws_fed=%d asr_fed=%d hits=%d warm_wake=%d loads_after_warm=%d acquire_log=%v machine=%s",
		rv.kwsFedSamples.Load(), rv.asrFedSamples.Load(), rv.kwsHits.Load(), rv.warmWake.Load(),
		st.LoadsAfterWarm, rv.residency.acquireLog(), string(rv.machine.State()))
}

// voice41a2r1StubKws is the WakeWordEngine seam with nothing behind it: the
// guard-row case needs a RESIDENT engine (the fact the rows read) and none of
// the D43 dispatches below feed samples, so the stub's Accept is never called.
type voice41a2r1StubKws struct {
	mu     sync.Mutex
	closes int
}

func (s *voice41a2r1StubKws) Accept(samples []float32, out []speech.KeywordHit) []speech.KeywordHit {
	return out
}

func (s *voice41a2r1StubKws) Close() error {
	s.mu.Lock()
	s.closes++
	s.mu.Unlock()
	return nil
}

// Test41a2r1GuardRowsWalkWithTheProductionFact walks the rows that read
// Facts.KwsLoaded (#8/#10) and the settle row #32, on the production machine
// with the production fact getter - the fact is assigned in production here,
// not in a test-only literal. No DLLs: the engine behind the fact is the seam
// stub, whose Accept no dispatch below ever feeds.
func Test41a2r1GuardRowsWalkWithTheProductionFact(t *testing.T) {
	frames := audio.NewBoundedFrames()
	rv := newResidentVoice(frames, voicePlan{}, &residentBall{}, &voice15eSink{})
	rv.residency = newVoiceResidency(voicePlan{}, rv.queue, nil)
	stub := &voice41a2r1StubKws{}
	rv.kwsWake = stub

	if !rv.kwsLoaded() || !rv.kwsFacts().KwsLoaded {
		t.Fatal("the production fact getter does not see a resident engine")
	}

	// #5, then #7, then #13 puts the machine in Warm without any voice work;
	// #31 lands Settling, and #32 walks the kws-loaded branch with the real
	// fact: the engine is resident, so Settling returns to Armed.
	mustWalk := func(ev statemachine.Event, facts *statemachine.Facts, want statemachine.State, why string) {
		t.Helper()
		if !rv.dispatch(ev, facts, why) {
			t.Fatalf("dispatch %s refused from %s (problems=%v)", string(ev), string(rv.machine.State()), rv.problemsSnapshot())
		}
		if st := rv.machine.State(); st != want {
			t.Fatalf("dispatch %s landed %s, want %s", string(ev), string(st), string(want))
		}
	}
	mustWalk(statemachine.EvKwsEnabled, nil, statemachine.StateArmed, "D43#5")
	mustWalk(statemachine.EvWakeWord, nil, statemachine.StateListening, "D43#7")
	mustWalk(statemachine.EvVeto, nil, statemachine.StateWarm, "D43#13")
	mustWalk(statemachine.EvWarmIdle, nil, statemachine.StateSettling, "D43#31")
	mustWalk(statemachine.EvSettleExpired, rv.kwsFacts(), statemachine.StateArmed,
		"D43#32 kws-loaded branch: the engine is resident, so Settling returns to Armed")
	// The mute cycle over the same fact: #8 from Armed, then #10's kws-loaded
	// branch back to Armed. The gesture entry is exercised too (it posts to the
	// owner queue; the case drains the queue the way the owner does).
	rv.onMuteGesture()
	rv.queue.run()
	if st := rv.machine.State(); st != statemachine.StateMuted {
		t.Fatalf("mute gesture from Armed landed %s, want Muted (D43#8)", string(st))
	}
	if !voice41a2r1Walked(t, rv, statemachine.StateArmed, statemachine.EvMuteKey, statemachine.StateMuted) {
		t.Fatal("the timeline has no D43#8 walk")
	}
	rv.onMuteGesture()
	rv.queue.run()
	if st := rv.machine.State(); st != statemachine.StateArmed {
		t.Fatalf("unmute gesture from Muted landed %s, want Armed (D43#10, kws-loaded branch)", string(st))
	}
	// The kws-not-loaded branch of both rows, with the same getter after the
	// engine is gone: #8 still walks Armed->Muted (no guard on it), and #10 then
	// falls to Sleeping. #32 needs the same round trip to Settling.
	rv.kwsWake = nil
	if rv.kwsFacts().KwsLoaded {
		t.Fatal("the production fact getter still claims an engine after it is gone")
	}
	mustWalk(statemachine.EvMuteKey, rv.kwsFacts(), statemachine.StateMuted, "D43#8 (no guard)")
	mustWalk(statemachine.EvMuteKey, rv.kwsFacts(), statemachine.StateSleeping,
		"D43#10 kws-not-loaded branch")
	mustWalk(statemachine.EvKwsEnabled, nil, statemachine.StateArmed, "D43#5 again")
	mustWalk(statemachine.EvWakeWord, nil, statemachine.StateListening, "D43#7 again")
	mustWalk(statemachine.EvVeto, nil, statemachine.StateWarm, "D43#13 again")
	mustWalk(statemachine.EvWarmIdle, nil, statemachine.StateSettling, "D43#31 again")
	mustWalk(statemachine.EvSettleExpired, rv.kwsFacts(), statemachine.StateSleeping,
		"D43#32 kws-not-loaded branch")
}

// Test41a2r1ProductionTableRejectsTheDemoWavs is the pipeline-level form of the
// engine leg's pinned control: the engine built the way PRODUCTION builds it
// (voiceKwsFilesFrom's five manifest paths, keywords.txt among them) scores
// zero hits on the demo fixtures through the C9 seam, because nobody said a
// production wake word in those recordings. This is the reading that keeps the
// wake cases honest: a hit there requires the demo table pairing, and "every
// fixture must wake" would be a lie about the model. No DLL-free shortcut: the
// engine is the real one, the feeding is the seam's 512-sample window.
func Test41a2r1ProductionTableRejectsTheDemoWavs(t *testing.T) {
	kws := voiceKwsFilesFrom(voice41a2r1KwsDir(t))
	if !kws.complete() {
		t.Fatalf("the KWS manifest spelling did not resolve under %s", voice41a2r1KwsDir(t))
	}
	eng, err := speech.NewKwsEngine(kws.kwsConfig())
	if err != nil {
		t.Fatalf("the production-shaped engine refused to build: %v", err)
	}
	defer func() { _ = eng.Close() }()
	for _, name := range []string{"0.wav", "1.wav", "2.wav", "3.wav", "4.wav", "5.wav", "6.wav"} {
		pcm := voice41a2r1WavPCM(t, voice41a2r1Fixture(t, name))
		var hits []speech.KeywordHit
		for off := 0; off < len(pcm); off += audio.FrameSamples {
			end := off + audio.FrameSamples
			if end > len(pcm) {
				end = len(pcm)
			}
			hits = eng.Accept(voiceSamples(pcm[off:end]), hits)
		}
		if len(hits) != 0 {
			t.Fatalf("%s produced %d hit(s) against the production table: nobody said a production "+
				"wake word in the demo recordings, and the pipeline table is the production one",
				name, len(hits))
		}
	}
}

// voice41a2r1FileExists is the test helper's existence check (the production
// one lives on voiceFileExists and is not exported to the test's intent).
func voice41a2r1FileExists(p string) bool {
	st, err := os.Stat(p)
	return err == nil && st.Mode().IsRegular()
}

// Test41a2r1OptInGateResolvesTheManifestSpelling is the opt-in half of the

// voice41a2r1WavPCM returns the wav's data chunk as 16k mono int16, parsed with
// the same chunk walk voice41a2r1WavSamples uses.
func voice41a2r1WavPCM(t *testing.T, wav string) []int16 {
	t.Helper()
	data, err := os.ReadFile(wav)
	if err != nil {
		t.Fatalf("read %s: %v", wav, err)
	}
	for off := 12; off+8 <= len(data); {
		id := string(data[off : off+4])
		size := int(binary.LittleEndian.Uint32(data[off+4 : off+8]))
		body := off + 8
		if id == "data" {
			n := size
			if body+n > len(data) {
				n = len(data) - body
			}
			pcm := make([]int16, n/2)
			for i := range pcm {
				pcm[i] = int16(binary.LittleEndian.Uint16(data[body+2*i:]))
			}
			return pcm
		}
		off = body + size + (size & 1)
	}
	t.Fatalf("%s carries no data chunk", wav)
	return nil
}

// Test41a2r1OptInGateResolvesTheManifestSpelling is the opt-in half of the
// posture plus the store-spelling resolution: enabled=false resolves nothing
// and says why; enabled=true resolves the five manifest files through
// WISP_KWS_DIR. No DLLs: this is the plan, not the engine.
func Test41a2r1OptInGateResolvesTheManifestSpelling(t *testing.T) {
	dataDir := voice15eDataDir(t)
	frames := audio.NewBoundedFrames()
	ra := &residentAudio{frames: frames}

	// Off: the opt-in never asked. Empty set, reason spelled, nothing resolved.
	voice15eWriteConfig(t, dataDir, nil) // shipped defaults: wake_word.enabled=false
	rv := buildResidentVoice(dataDir, ra, nil, nil)
	if rv.plan.Kws.complete() {
		t.Fatalf("the opt-in-off plan resolved a KWS set: %+v", rv.plan.Kws)
	}
	if !strings.Contains(rv.plan.KwsNote, "enabled=false") {
		t.Fatalf("the opt-in-off note = %q, want the enabled=false sentence", rv.plan.KwsNote)
	}

	// On: the hatch carries the checkpoint. The manifest spelling resolves all
	// five files and the note is empty (resolved).
	t.Setenv("WISP_KWS_DIR", voice41a2r1KwsDir(t))
	voice15eWriteConfig(t, dataDir, func(c *config.Config) {
		c.Voice.WakeWord.Enabled = true
	})
	rv = buildResidentVoice(dataDir, ra, nil, nil)
	if !rv.plan.Kws.complete() {
		t.Fatalf("the opt-in-on plan did not resolve the KWS set (note=%q)", rv.plan.KwsNote)
	}
	if !strings.HasSuffix(rv.plan.Kws.Encoder, "encoder-epoch-12-avg-2-chunk-16-left-64.int8.onnx") {
		t.Fatalf("the resolved encoder is not the manifest spelling: %s", rv.plan.Kws.Encoder)
	}
	if rv.plan.KwsNote != "" {
		t.Fatalf("a resolved plan carries a note: %q", rv.plan.KwsNote)
	}
}
