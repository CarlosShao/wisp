package speech

import (
	"fmt"
	"time"

	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"

	"github.com/CarlosShao/wisp/internal/audio"
)

// asrFeatureDim is the mel feature width every sherpa pre-trained streaming
// model is trained on. It belongs to the model file, not to configuration.
const asrFeatureDim = 80

// asrDecodingMethod is the search the streaming Paraformer supports. It is a
// constant for the same reason IntraOpThreads is: the shape of the result (and
// the CPU story) has to be a property of the code, not of whoever called it.
const asrDecodingMethod = "greedy_search"

// StreamingLatencyBudget is the ceiling this session is held to: from the
// moment the VAD closes an utterance to the final transcript, 800ms on the
// streaming path (SPEC-04 sec 8). The offline engines get 1500ms, which is
// exactly why the streaming engine is the one that ships (15-a1). The same
// number bounds the first partial, because a partial that arrives after the
// final is not a partial at all.
const StreamingLatencyBudget = 800 * time.Millisecond

// AsrConfig points at the streaming Paraformer pair and its token table.
// There is deliberately no thread-count and no endpoint-detection field: see
// IntraOpThreads and asrRecognizerConfig.
type AsrConfig struct {
	// EncoderPath is the streaming encoder (asr-encoder.int8.onnx in
	// third_party/spike-models; the product path comes from the model store).
	EncoderPath string
	// DecoderPath is the streaming decoder (asr-decoder.int8.onnx).
	DecoderPath string
	// TokensPath is the token table (asr-tokens.txt).
	TokensPath string
}

// StreamingAsr is one streaming Paraformer session: audio in, a growing
// hypothesis out, a final transcript when the utterance closes.
//
// It owns exactly one sherpa recognizer for its lifetime (loading the int8
// encoder is the expensive part, so it is done once) and leases exactly one
// sherpa stream per utterance. That is the whole cgo surface of the ASR leg:
// no other runtime object is created, which is what keeps the chain on the
// single process-wide onnxruntime (D5).
//
// StreamingAsr starts no goroutines and is not safe for concurrent use: one
// owner goroutine drives it (D38b). It never measures wall-clock time; every
// duration it reports comes from the audio clock (D42#9).
type StreamingAsr struct {
	rec    *sherpa.OnlineRecognizer
	stream *sherpa.OnlineStream
	cfg    AsrConfig
	fed    int // samples fed into the open (or most recently closed) utterance
}

// NewStreamingAsr loads the model pair and returns a session that is ready to
// Begin an utterance.
//
// It refuses, rather than crashes, when the process is running on an
// onnxruntime older than the linked C API gives it: see checkOnnxruntime.
func NewStreamingAsr(cfg AsrConfig) (*StreamingAsr, error) {
	if cfg.EncoderPath == "" || cfg.DecoderPath == "" || cfg.TokensPath == "" {
		return nil, fmt.Errorf("speech: AsrConfig needs EncoderPath, DecoderPath and " +
			"TokensPath (the streaming Paraformer pair plus its token table)")
	}
	if err := checkOnnxruntime(); err != nil {
		return nil, err
	}
	rec := sherpa.NewOnlineRecognizer(asrRecognizerConfig(cfg))
	if rec == nil {
		return nil, fmt.Errorf("speech: sherpa refused the streaming ASR config "+
			"(encoder=%s decoder=%s tokens=%s): a file is unreadable, or the pair is "+
			"not a streaming Paraformer", cfg.EncoderPath, cfg.DecoderPath, cfg.TokensPath)
	}
	return &StreamingAsr{rec: rec, cfg: cfg}, nil
}

// asrRecognizerConfig is the one place the recognizer's shape is decided, and
// it is a pure function so the shape can be asserted without a model file.
//
// Two decisions are load-bearing:
//
//   - NumThreads is IntraOpThreads, never a parameter (AC#6, D32).
//   - EnableEndpoint is 0. Endpointing is owned by Endpointer in Go, on the
//     audio clock, where tickets can measure it; sherpa's own endpointing is
//     the C-side twin of the same decision and would be a second, silent
//     source of truth exactly like silero's min_silence/min_speech would be
//     (see NewDetector, which pins those to 0 for the same reason).
func asrRecognizerConfig(cfg AsrConfig) *sherpa.OnlineRecognizerConfig {
	return &sherpa.OnlineRecognizerConfig{
		FeatConfig: sherpa.FeatureConfig{
			SampleRate: audio.TargetRate,
			FeatureDim: asrFeatureDim,
		},
		ModelConfig: sherpa.OnlineModelConfig{
			Paraformer: sherpa.OnlineParaformerModelConfig{
				Encoder: cfg.EncoderPath,
				Decoder: cfg.DecoderPath,
			},
			Tokens:     cfg.TokensPath,
			NumThreads: IntraOpThreads,
			Provider:   "cpu",
		},
		DecodingMethod: asrDecodingMethod,
		EnableEndpoint: 0,
	}
}

// Begin opens an utterance: a fresh decoder state, so nothing from the previous
// utterance can leak into this one (no cross-utterance context, by design --
// 15-c emits one transcript per endpoint, and punctuation/stitching is a later
// stage, ticket 27).
func (a *StreamingAsr) Begin() {
	if a.rec == nil {
		return
	}
	if a.stream != nil {
		sherpa.DeleteOnlineStream(a.stream)
	}
	a.stream = sherpa.NewOnlineStream(a.rec)
	a.fed = 0
}

// Feed appends utterance audio (16k mono, any length) and runs the decoder as
// far as the stream is ready. It is a no-op when no utterance is open, so a
// caller that feeds a whole recording keeps its own pacing and still gets one
// transcript per Begin/Finish pair.
func (a *StreamingAsr) Feed(samples []float32) {
	if a.rec == nil || a.stream == nil || len(samples) == 0 {
		return
	}
	a.stream.AcceptWaveform(audio.TargetRate, samples)
	a.fed += len(samples)
	for a.rec.IsReady(a.stream) {
		a.rec.Decode(a.stream)
	}
}

// Partial is the hypothesis for the open utterance so far. It is empty until
// the encoder has consumed one full chunk, which is why the first partial
// cannot appear at 0ms; after that it grows and never rewinds, because the
// decoder state only ever moves forward within an utterance.
func (a *StreamingAsr) Partial() string {
	if a.rec == nil || a.stream == nil {
		return ""
	}
	return a.rec.GetResult(a.stream).Text
}

// Finish flushes the feature extractor, decodes what is left, returns the
// utterance's final transcript and releases the stream. It is safe to call with
// no utterance open (returns "").
//
// The C API has no error return here, so neither does Finish: a failed decode
// shows up as a short or empty transcript, which is the signal AC#5's retry
// policy keys on (that policy is not this slice's).
func (a *StreamingAsr) Finish() string {
	if a.rec == nil || a.stream == nil {
		return ""
	}
	a.stream.InputFinished()
	for a.rec.IsReady(a.stream) {
		a.rec.Decode(a.stream)
	}
	text := a.rec.GetResult(a.stream).Text
	sherpa.DeleteOnlineStream(a.stream)
	a.stream = nil
	return text
}

// InFlight reports whether an utterance is open.
func (a *StreamingAsr) InFlight() bool { return a.stream != nil }

// Fed is how many samples the open -- or, right after Finish, the most
// recently closed -- utterance streamed. It is the audio clock the transcriber
// reports latencies on, and it is the observable behind the invariant that the
// ASR saw exactly the utterance the endpoint closed.
func (a *StreamingAsr) Fed() int { return a.fed }

// Close releases the recognizer. It is idempotent (C11).
func (a *StreamingAsr) Close() error {
	if a.stream != nil {
		sherpa.DeleteOnlineStream(a.stream)
		a.stream = nil
	}
	if a.rec != nil {
		sherpa.DeleteOnlineRecognizer(a.rec)
		a.rec = nil
	}
	return nil
}
