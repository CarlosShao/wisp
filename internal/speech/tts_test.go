package speech

import (
	"errors"
	"math"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"
)

// The Matcha TTS quadruple, as pinned in third_party/spike-models.
const (
	spikeTtsModelDir = "matcha-icefall-zh-baker"
	spikeTtsVocoder  = "vocos-22khz-univ.onnx"
)

// matchaConfigForTest is the production config pointed at the spike models.
func matchaConfigForTest(t *testing.T) TtsConfig {
	t.Helper()
	d := spikeModelsDir(t)
	m := filepath.Join(d, spikeTtsModelDir)
	cfg := TtsConfig{
		AcousticModel: filepath.Join(m, "model-steps-3.onnx"),
		Vocoder:       filepath.Join(d, spikeTtsVocoder),
		Lexicon:       filepath.Join(m, "lexicon.txt"),
		Tokens:        filepath.Join(m, "tokens.txt"),
		DictDir:       filepath.Join(m, "dict"),
		RuleFsts: []string{
			filepath.Join(m, "phone.fst"),
			filepath.Join(m, "date.fst"),
			filepath.Join(m, "number.fst"),
		},
	}
	for _, p := range []string{cfg.AcousticModel, cfg.Vocoder, cfg.Lexicon, cfg.Tokens} {
		if _, err := os.Stat(p); err != nil {
			t.Skipf("matcha tts assets are not staged on this host (%v); "+
				"the staging recipe is docs/evidence/s1/26-r1-impl.md", err)
		}
	}
	return cfg
}

// newTestMatchaTts opens the engine, or skips by name when this host cannot.
// A skip is not a pass (same rule as 15-c's newTestAsr).
func newTestMatchaTts(t *testing.T) *MatchaTts {
	t.Helper()
	tts, err := NewMatchaTts(matchaConfigForTest(t))
	if err != nil {
		t.Skipf("matcha tts unavailable on this host: %v", err)
	}
	t.Cleanup(func() { _ = tts.Close() })
	return tts
}

// pcmStats26 is peak absolute amplitude and RMS: the two numbers that separate
// "samples exist" from "the samples are speech".
func pcmStats26(samples []float32) (peak, rms float64) {
	if len(samples) == 0 {
		return 0, 0
	}
	var sumSquares float64
	for _, v := range samples {
		a := math.Abs(float64(v))
		if a > peak {
			peak = a
		}
		sumSquares += a * a
	}
	return peak, math.Sqrt(sumSquares / float64(len(samples)))
}

// Test26TheTtsEngineEmitsRealAudio is ticket 26's "TTS really makes sound"
// reading: not "the interface is implemented" but samples, duration, rate and
// amplitude of a real synthesis.
func Test26TheTtsEngineEmitsRealAudio(t *testing.T) {
	tts := newTestMatchaTts(t)

	// Hold it through the seam, not through the concrete type: what ships is the
	// interface, so the reading has to come from the interface.
	var eng TtsEngine = tts

	const line = "你好，我是威斯普。"
	var chunks, total, firstLen int
	var firstAt time.Duration
	var all []float32

	start := time.Now()
	err := eng.Speak(line, func(samples []float32) error {
		chunks++
		if chunks == 1 {
			firstLen = len(samples)
			firstAt = time.Since(start)
		}
		total += len(samples)
		all = append(all, samples...)
		for _, v := range samples {
			if math.IsNaN(float64(v)) || math.Abs(float64(v)) > 1.5 {
				return errors.New("speech: tts emitted a sample outside the normalised range")
			}
		}
		return nil
	})
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("Speak(%q): %v", line, err)
	}
	if chunks == 0 || total == 0 {
		t.Fatalf("Speak(%q) emitted %d chunk(s) / %d samples, want real audio", line, chunks, total)
	}

	rate := eng.SampleRate()
	if rate != 22050 {
		t.Fatalf("SampleRate() = %d, want 22050 (the vocos 22kHz vocoder)", rate)
	}
	peak, rms := pcmStats26(all)
	t.Logf("tts real audio: text %q chunks=%d samples=%d rate=%d duration=%.3fs peak=%.4f rms=%.4f firstChunk=%d@%s total=%s",
		line, chunks, total, rate, float64(total)/float64(rate), peak, rms, firstLen, firstAt, elapsed)

	// Silence would satisfy "samples > 0" -- it would not satisfy "speak".
	if peak < 0.05 {
		t.Fatalf("synthesised audio is silent: peak=%.4f rms=%.4f over %d samples", peak, rms, total)
	}
	if rms < 0.005 {
		t.Fatalf("synthesised audio is near-silent: rms=%.5f over %d samples", rms, total)
	}
}

// Test26TheChunksAreChunksAndNotOneWholeClip is the seam's own property: the
// engine must hand over audio in playback order, starting before the whole
// utterance exists. A whole-clip wrapper satisfies TtsEngine's signature and
// fails this.
func Test26TheChunksAreChunksAndNotOneWholeClip(t *testing.T) {
	tts := newTestMatchaTts(t)

	const three = "你好。今天天气不错。我们开始吧。"
	var chunks, total, firstLen int
	var firstAt time.Duration
	start := time.Now()
	if err := tts.Speak(three, func(samples []float32) error {
		chunks++
		if chunks == 1 {
			firstLen = len(samples)
			firstAt = time.Since(start)
		}
		total += len(samples)
		return nil
	}); err != nil {
		t.Fatalf("Speak(three sentences): %v", err)
	}
	elapsed := time.Since(start)
	t.Logf("tts chunking: chunks=%d firstChunk=%d@%s totalSamples=%d totalWall=%s", chunks, firstLen, firstAt, total, elapsed)

	if chunks < 2 {
		t.Fatalf("Speak emitted %d chunk(s) for a three-sentence utterance; "+
			"a single buffer is the whole-clip shape the seam forbids", chunks)
	}
	if firstLen >= total {
		t.Fatalf("the first chunk was the whole utterance (%d of %d samples): not streaming", firstLen, total)
	}
	if chunks*1 < 3 {
		t.Fatalf("three sentences produced %d chunk(s); the engine's sentence stages were batched", chunks)
	}
}

// Test26BargeInStopsTheSentenceMidFlight is AC3's mechanism at this seam: an emit
// error must abort synthesis where it happens (the cooperative abort C5 uses),
// not after the sentence has already been spoken.
func Test26BargeInStopsTheSentenceMidFlight(t *testing.T) {
	tts := newTestMatchaTts(t)

	const three = "你好。今天天气不错。我们开始吧。"
	var full int
	if err := tts.Speak(three, func(s []float32) error { full += len(s); return nil }); err != nil {
		t.Fatalf("baseline Speak: %v", err)
	}

	stop := errors.New("host stopped playback")
	var got, chunks int
	err := tts.Speak(three, func(s []float32) error {
		chunks++
		got += len(s)
		return stop
	})
	if !errors.Is(err, stop) {
		t.Fatalf("Speak returned %v, want the emit error verbatim", err)
	}
	t.Logf("tts barge-in: emittedSamples=%d ofFull=%d chunks=%d", got, full, chunks)
	if chunks != 1 {
		t.Fatalf("abort after the first chunk still asked for %d chunks", chunks)
	}
	if got >= full {
		t.Fatalf("aborted run emitted %d samples, the full utterance is %d: nothing was stopped", got, full)
	}
}

// Test26EmptyAndPunctuationAreNotSilentSuccess is the positive control: nothing
// to say must be an error, not a zero-length success.
//
// The second half is why the guard is in Speak and not in the caller: the engine
// itself reports success for punctuation-only text and hands back a short buffer
// of near-silence. That is the shape the guard exists to convert into a named
// error.
func Test26EmptyAndPunctuationAreNotSilentSuccess(t *testing.T) {
	tts := newTestMatchaTts(t)

	for _, text := range []string{"", "   ", "\t\n", "。，、！？", "…––"} {
		emitted := 0
		err := tts.Speak(text, func(s []float32) error { emitted += len(s); return nil })
		if !errors.Is(err, ErrNothingToSpeak) {
			t.Fatalf("Speak(%q) = %v, want ErrNothingToSpeak", text, err)
		}
		if emitted != 0 {
			t.Fatalf("Speak(%q) emitted %d samples for text with nothing to say", text, emitted)
		}
	}

	// Counter-reading: the engine underneath does NOT agree, which is why the
	// guard cannot be dropped in favour of "the engine returns nothing".
	audio := tts.tts.Generate("。，、！？", ttsSpeakerID, ttsSpeed)
	if audio == nil || len(audio.Samples) == 0 {
		t.Fatalf("the guard is unobservable on this checkpoint: the engine returned no audio for " +
			"punctuation, so the silent-success trap this test names does not exist here")
	}
	peak, _ := pcmStats26(audio.Samples)
	t.Logf("engine-only punctuation reading: samples=%d duration=%.3fs peak=%.4f (guard turns this into an error)",
		len(audio.Samples), float64(len(audio.Samples))/float64(tts.SampleRate()), peak)

	if err := tts.Speak("你好", nil); err == nil {
		t.Fatalf("Speak accepted a nil emit callback")
	}
}

// Test26TheWarmFirstChunkMeetsTheBudgetOverTenRuns is AC4's warm half: ten runs on
// an already-loaded engine, with the percentiles recorded rather than asserted
// away. The assertion is the one that cannot flake: every run produced audio, and
// the first chunk arrived before the utterance finished.
func Test26TheWarmFirstChunkMeetsTheBudgetOverTenRuns(t *testing.T) {
	tts := newTestMatchaTts(t)

	const line = "你好，我是威斯普。"
	const runs = 10
	firstChunk := make([]time.Duration, 0, runs)
	total := make([]time.Duration, 0, runs)
	for i := 0; i < runs; i++ {
		var chunks int
		var firstAt time.Duration
		start := time.Now()
		if err := tts.Speak(line, func(s []float32) error {
			chunks++
			if chunks == 1 {
				firstAt = time.Since(start)
			}
			return nil
		}); err != nil {
			t.Fatalf("run %d: %v", i, err)
		}
		if chunks == 0 {
			t.Fatalf("run %d emitted no audio", i)
		}
		firstChunk = append(firstChunk, firstAt)
		total = append(total, time.Since(start))
	}

	p50, p95 := percentile26(firstChunk, 50), percentile26(firstChunk, 95)
	t.Logf("tts warm first chunk over %d runs: p50=%s p95=%s budget=%s headroom=%.2fx",
		runs, p50, p95, TtsFirstFrameBudgetWarm,
		float64(TtsFirstFrameBudgetWarm)/float64(p50))
	t.Logf("tts warm full sentence over %d runs: p50=%s p95=%s",
		runs, percentile26(total, 50), percentile26(total, 95))

	if p50 <= 0 {
		t.Fatalf("first-chunk latency p50 = %s: no sample arrived before the sentence ended", p50)
	}
	if p50 > TtsFirstFrameBudgetWarm {
		// Not a flake guard: this is the reading AC4 asks for. If it fires, the
		// number to write in the evidence is "over budget", not a smaller run.
		t.Fatalf("warm first chunk p50 = %s, budget %s (%d runs)", p50, TtsFirstFrameBudgetWarm, runs)
	}
}

// percentile26 is nearest-rank: no interpolation, so a single slow run cannot be
// averaged away.
func percentile26(d []time.Duration, p int) time.Duration {
	if len(d) == 0 {
		return 0
	}
	cp := append([]time.Duration(nil), d...)
	sort.Slice(cp, func(i, j int) bool { return cp[i] < cp[j] })
	idx := (p*len(cp) + 99) / 100
	if idx < 1 {
		idx = 1
	}
	if idx > len(cp) {
		idx = len(cp)
	}
	return cp[idx-1]
}

// Test26TheBudgetsAreTheSpecNumbers pins the two numbers the readings are quoted
// against (SPEC-04 section 8). A drift here would silently redefine "meets the
// budget" for every future run.
func Test26TheBudgetsAreTheSpecNumbers(t *testing.T) {
	if TtsFirstFrameBudgetWarm != 800*time.Millisecond {
		t.Fatalf("warm first-frame budget = %s, want 800ms (SPEC-04 section 8)", TtsFirstFrameBudgetWarm)
	}
	if TtsFirstFrameBudgetCold != 1600*time.Millisecond {
		t.Fatalf("cold first-frame budget = %s, want 1600ms (SPEC-04 section 8)", TtsFirstFrameBudgetCold)
	}
}

// Test26TheClosedEngineRefusesToSpeak: Close is idempotent (C11) and a closed
// engine is not silently reusable -- a second Speak after Close must not reach a
// freed pointer.
func Test26TheClosedEngineRefusesToSpeak(t *testing.T) {
	tts, err := NewMatchaTts(matchaConfigForTest(t))
	if err != nil {
		t.Skipf("matcha tts unavailable on this host: %v", err)
	}
	if err := tts.Close(); err != nil {
		t.Fatalf("first Close: %v", err)
	}
	if err := tts.Close(); err != nil {
		t.Fatalf("second Close must be a no-op: %v", err)
	}
	if err := tts.Speak("你好", func([]float32) error { return nil }); err == nil {
		t.Fatalf("Speak after Close returned nil")
	}
}
