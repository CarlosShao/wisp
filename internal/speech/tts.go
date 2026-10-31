package speech

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
)

// This file is the local implementation of the C9 TTS seam (ticket 26,
// 26-r1): sherpa-onnx MatchaTTS in Chinese (matcha-icefall-zh-baker plus the
// vocos 22kHz vocoder) behind the TtsEngine interface that 15-a froze in
// engines.go. The seam is not touched here -- an implementation that needs a
// different interface is the thing that is wrong.
//
// SHIPPING STATUS: this is a MEASUREMENT-GRADE REFERENCE IMPLEMENTATION, not
// the product's default TTS. The weights it loads are recorded in
// models/manifest.json as tts-matcha-zh-baker, status "blocked-p3", license
// "NON-COMMERCIAL-ONLY ... must not ship; dev/measurement use only", and the
// owner's ruling recorded in the ledger as R5 keeps the default voice path on
// cloud TTS (StepFun, ticket 61). Nothing here selects an engine: engine choice
// stays with the composition root, and this file adds no configuration key, so
// this implementation cannot become the default by being compiled.
//
// Why the numbers below are constants rather than fields: 15-b and 15-c made the
// same call for the ASR leg (the result shape and the CPU story have to be
// properties of the code, not of whoever called it). A TTS speed knob would let
// a caller change what a "duration" reading means between two runs of the same
// evidence, which is the one thing the latency budget cannot tolerate.
const (
	// ttsNoiseScale and ttsLengthScale are the Matcha defaults upstream
	// documents for this model family; they are not taste, they are the values
	// the checkpoint was trained against.
	ttsNoiseScale  = 0.667
	ttsLengthScale = 1.0

	// ttsSpeakerID is the only speaker this checkpoint has. NumSpeakers is read
	// back in the tests, so a multi-speaker model would be noticed rather than
	// silently speaking with voice 0.
	ttsSpeakerID = 0

	// ttsSpeed 1.0 is the checkpoint's own rate; the honest way to make speech
	// longer is to speak more text, not to resample the budget.
	ttsSpeed = 1.0

	// ttsNumThreads pins synthesis to one core for the same reason
	// IntraOpThreads pins ASR (ticket 15 AC6): a thread pool fan-out is what
	// makes a latency reading depend on the host's core count.
	ttsNumThreads = 1

	// ttsMaxNumSentences 1 is what turns the engine's own sentence stages into
	// the chunks the seam emits: with a queue depth above one the engine is free
	// to batch a whole paragraph and hand back a single long buffer, which is
	// precisely the shape the seam's chunked contract exists to forbid.
	ttsMaxNumSentences = 1
)

// The two budgets this leg is measured against (SPEC-04 section 8, "first token
// to TTS first frame"): 800ms when the model is already resident, 1600ms when it
// has to be loaded first. They are exported so the readings and the tests quote
// the same numbers the spec does instead of re-typing them.
const (
	TtsFirstFrameBudgetWarm = 800 * time.Millisecond
	TtsFirstFrameBudgetCold = 1600 * time.Millisecond
)

// ErrNothingToSpeak is what Speak returns for text that has nothing to say --
// empty, whitespace only, or punctuation only.
//
// This error is the whole point of the guard behind it. The engine does not
// agree that punctuation is silent input: synthesising the cluster "。" returns
// 3840 samples (174ms) of near-silence with a nil error, so an engine that just
// forwarded the text would report success, emit a silent blip, and let the host
// believe it spoke. A caller cannot tell that apart from speaking too quietly.
// Naming it as an error is what lets the announce path fall back to text instead
// of pretending (ticket 26's "text-only fallback still delivers the result").
var ErrNothingToSpeak = errors.New("speech: nothing to speak (empty or punctuation-only text)")

// TtsConfig points at one MatchaTTS checkpoint and its front-end assets.
//
// There is no thread count, no speed and no sentence limit here on purpose: those
// are the package constants above, for the reasons written there.
type TtsConfig struct {
	// AcousticModel is the Matcha acoustic model (model-steps-3.onnx in
	// third_party/spike-models/matcha-icefall-zh-baker; the product path comes
	// from the model store).
	AcousticModel string
	// Vocoder is the neural vocoder that turns the acoustic model's frames into
	// samples (vocos-22khz-univ.onnx). Matcha does not produce audio without it.
	Vocoder string
	// Lexicon is the grapheme-to-phoneme table (lexicon.txt).
	Lexicon string
	// Tokens is the phoneme token table (tokens.txt).
	Tokens string
	// DictDir is the jieba-style word-segmentation dictionary directory (dict/).
	DictDir string
	// RuleFsts are the text-normalisation rule FSTs (phone/date/number.fst),
	// applied in order: numbers and dates have to be verbalised before the
	// lexicon can spell them.
	RuleFsts []string
}

// MatchaTts is one MatchaTTS session: text in, chunks of 22.05kHz mono float32
// out, in playback order.
//
// It owns exactly one sherpa TTS instance for its lifetime (loading the fp32
// acoustic model is 5.8s on the spike host, so it is done once) and it starts no
// goroutines: one owner goroutine drives it, like every other engine in this
// package (D38b). It is not safe for concurrent use.
type MatchaTts struct {
	tts        *sherpa.OfflineTts
	sampleRate int
	closed     bool
}

// The compile-time conformance assertion C9 asks for: if this engine ever stops
// matching the seam, the package stops building. It is the same shape 15-a used
// for *StreamingAsr in engines.go, kept next to the implementation so that a
// change to either side fails here first.
var _ TtsEngine = (*MatchaTts)(nil)

// NewMatchaTts loads the checkpoint and asserts the one property the seam's
// contract depends on: that the engine can name its own sample rate. A model
// that loads but reports 0 cannot be resampled for the render side, so it is
// rejected here instead of at the first playback.
//
// It refuses, rather than crashes, when the process is running on an
// onnxruntime older than the linked C API gives it: see checkOnnxruntime.
func NewMatchaTts(cfg TtsConfig) (*MatchaTts, error) {
	switch {
	case cfg.AcousticModel == "":
		return nil, errors.New("speech: matcha tts needs an acoustic model path")
	case cfg.Vocoder == "":
		return nil, errors.New("speech: matcha tts needs a vocoder path")
	case cfg.Lexicon == "" || cfg.Tokens == "":
		return nil, errors.New("speech: matcha tts needs a lexicon and a token table")
	}
	if err := checkOnnxruntime(); err != nil {
		return nil, err
	}

	model := sherpa.OfflineTtsModelConfig{
		Matcha: sherpa.OfflineTtsMatchaModelConfig{
			AcousticModel: cfg.AcousticModel,
			Vocoder:       cfg.Vocoder,
			Lexicon:       cfg.Lexicon,
			Tokens:        cfg.Tokens,
			DictDir:       cfg.DictDir,
			NoiseScale:    ttsNoiseScale,
			LengthScale:   ttsLengthScale,
		},
		NumThreads: ttsNumThreads,
	}
	c := sherpa.OfflineTtsConfig{
		Model:           model,
		RuleFsts:        strings.Join(cfg.RuleFsts, ","),
		MaxNumSentences: ttsMaxNumSentences,
	}

	tts := sherpa.NewOfflineTts(&c)
	if tts == nil {
		// The C API returns nil for every load failure (missing file, unreadable
		// vocabulary, a shadowing onnxruntime), so the error names the inputs
		// rather than inventing a cause it cannot see.
		return nil, fmt.Errorf("speech: matcha tts failed to load (acoustic model %q, vocoder %q)",
			cfg.AcousticModel, cfg.Vocoder)
	}
	rate := tts.SampleRate()
	if rate <= 0 {
		sherpa.DeleteOfflineTts(tts)
		return nil, fmt.Errorf("speech: matcha tts reports sample rate %d", rate)
	}
	return &MatchaTts{tts: tts, sampleRate: rate}, nil
}

// Speak synthesises text and emits its audio in playback order, starting with
// the first chunk as soon as one exists.
//
// The engine's progress callback is the seam's callback: each chunk the engine
// finishes is handed to emit, and an emit error is reported back to the engine as
// "stop", which is how a sentence is abandoned mid-flight. No goroutine, no
// timer and no wall clock is involved -- the cooperative abort C5's LlmProvider
// already uses is the same contract here.
func (m *MatchaTts) Speak(text string, emit func(samples []float32) error) error {
	if emit == nil {
		return errors.New("speech: Speak needs an emit callback")
	}
	if m == nil || m.tts == nil || m.closed {
		return errors.New("speech: Speak on a closed or unloaded tts engine")
	}
	if !speakable(text) {
		// Bail before synthesis, not after: see ErrNothingToSpeak.
		return ErrNothingToSpeak
	}

	var emitErr error
	audio := m.tts.GenerateWithProgressCallback(text, ttsSpeakerID, ttsSpeed,
		func(chunk []float32, _ float32) bool {
			if emitErr != nil {
				return false
			}
			if err := emit(chunk); err != nil {
				emitErr = err
				return false
			}
			return true
		})
	if emitErr != nil {
		return emitErr
	}
	if audio == nil {
		return fmt.Errorf("speech: matcha tts produced no audio for %d rune(s)",
			utf8.RuneCountInString(text))
	}
	return nil
}

// SampleRate is the rate of the samples Speak emits (the render side resamples
// from it). It is a property of the checkpoint, read back at load time.
func (m *MatchaTts) SampleRate() int {
	if m == nil {
		return 0
	}
	return m.sampleRate
}

// Close releases the engine (C11). Idempotent: the release is what frees the
// model's memory, so calling it twice must not free it twice.
func (m *MatchaTts) Close() error {
	if m == nil || m.closed {
		return nil
	}
	m.closed = true
	if m.tts != nil {
		sherpa.DeleteOfflineTts(m.tts)
		m.tts = nil
	}
	return nil
}

// speakable reports whether text carries anything a voice can say. Letters and
// digits count (Han characters are letters); whitespace, punctuation and symbols
// do not.
func speakable(text string) bool {
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return true
		}
	}
	return false
}
