package speech

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/CarlosShao/wisp/internal/audio"
)

// The streaming Paraformer triple, as pinned in third_party/spike-models.
const (
	spikeDriverEncoder = "asr-encoder.int8.onnx"
	spikeDriverDecoder = "asr-decoder.int8.onnx"
	spikeDriverTokens  = "asr-tokens.txt"
)

// spikeModelsDir locates the pinned model directory. It honours
// WISP_SPIKE_MODELS and otherwise walks up from the working directory, so the
// same test binary works under "go test" (cwd = the package directory) and when
// run from a staging directory that holds the sherpa DLLs.
func spikeModelsDir(t *testing.T) string {
	t.Helper()
	if p := os.Getenv("WISP_SPIKE_MODELS"); p != "" {
		return p
	}
	rel := filepath.Join("third_party", "spike-models")
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(dir, rel)); err == nil {
			return filepath.Join(dir, rel)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatalf("spike models not found: set WISP_SPIKE_MODELS, or run from inside the repo "+
		"(looked for %s)", rel)
	return ""
}

// asrConfigForTest is the production config pointed at the spike models.
func asrConfigForTest(t *testing.T) AsrConfig {
	t.Helper()
	d := spikeModelsDir(t)
	return AsrConfig{
		EncoderPath: filepath.Join(d, spikeDriverEncoder),
		DecoderPath: filepath.Join(d, spikeDriverDecoder),
		TokensPath:  filepath.Join(d, spikeDriverTokens),
	}
}

// newTestAsr opens a recognizer, or skips by name when this host is running on
// a shadowed onnxruntime (see checkOnnxruntime). A skip is not a pass: the
// staging recipe in docs/evidence/s1/15c-r1-impl.md is what makes these run.
func newTestAsr(t *testing.T) *StreamingAsr {
	t.Helper()
	a, err := NewStreamingAsr(asrConfigForTest(t))
	if err != nil {
		t.Skipf("streaming ASR session unavailable on this host: %v", err)
	}
	t.Cleanup(func() { _ = a.Close() })
	return a
}

// newTestTranscriber opens the whole voice half: silero endpointing plus the
// streaming recognizer, wired by Transcriber.
func newTestTranscriber(t *testing.T) (*Transcriber, *StreamingAsr) {
	t.Helper()
	asr := newTestAsr(t)
	det, err := NewDetector(DetectorConfig{ModelPath: sileroModelPath(t)})
	if err != nil {
		t.Skipf("VAD session unavailable on this host: %v", err)
	}
	t.Cleanup(func() { _ = det.Close() })
	return NewTranscriber(det, asr), asr
}

// feedWindow pushes one seam window worth of pcm through the transcriber and
// returns the transcripts it produced.
func feedWindow(tr *Transcriber, pcm []int16) []Transcript {
	return tr.Accept(pcmToFloat(pcm), nil)
}

// TestAsrRecognizerConfigPinsTheShape is the wiring receipt for the two
// decisions that must not drift: one intra-op thread, and endpointing owned by
// the Go state machine rather than by sherpa.
//
// It asserts the literal 1 and the literal 0 on purpose. Asserting against
// IntraOpThreads would make this test agree with any value that constant ever
// takes, which is exactly how a pin stops being a pin: mutating the constant to
// 4 must turn this red (receipt: reverse criterion B in 15c-r1-impl.md).
func TestAsrRecognizerConfigPinsTheShape(t *testing.T) {
	cfg := AsrConfig{
		EncoderPath: "enc.onnx",
		DecoderPath: "dec.onnx",
		TokensPath:  "tokens.txt",
	}
	got := asrRecognizerConfig(cfg)

	if n := got.ModelConfig.NumThreads; n != 1 {
		t.Fatalf("ModelConfig.NumThreads = %d, want the literal 1: ASR inference must be "+
			"single-core pinned (AC#6, D32), and the pin is only real while nothing can raise it", n)
	}
	if sr := got.FeatConfig.SampleRate; sr != 16000 {
		t.Fatalf("FeatConfig.SampleRate = %d, want 16000 (SPEC-04 sec 3 / audio.TargetRate)", sr)
	}
	if d := got.FeatConfig.FeatureDim; d != 80 {
		t.Fatalf("FeatConfig.FeatureDim = %d, want 80 (the width every sherpa streaming model is trained on)", d)
	}
	if got.EnableEndpoint != 0 {
		t.Fatalf("EnableEndpoint = %d, want 0: endpointing is owned by Endpointer on the audio "+
			"clock; sherpa's own would be a second silent source of truth", got.EnableEndpoint)
	}
	if m := got.DecodingMethod; m != asrDecodingMethod {
		t.Fatalf("DecodingMethod = %q, want %q", m, asrDecodingMethod)
	}
	if p := got.ModelConfig.Paraformer; p.Encoder != cfg.EncoderPath || p.Decoder != cfg.DecoderPath {
		t.Fatalf("paraformer pair = (%q, %q), want (%q, %q)", p.Encoder, p.Decoder,
			cfg.EncoderPath, cfg.DecoderPath)
	}
	if tok := got.ModelConfig.Tokens; tok != cfg.TokensPath {
		t.Fatalf("tokens = %q, want %q", tok, cfg.TokensPath)
	}
	if prov := got.ModelConfig.Provider; prov != "cpu" {
		t.Fatalf("provider = %q, want cpu", prov)
	}
}

// TestNewStreamingAsrRejectsIncompleteConfig: a half-specified config must be an
// error, not a session that decodes garbage.
func TestNewStreamingAsrRejectsIncompleteConfig(t *testing.T) {
	for _, cfg := range []AsrConfig{
		{},
		{EncoderPath: "enc.onnx"},
		{EncoderPath: "enc.onnx", DecoderPath: "dec.onnx"},
	} {
		if _, err := NewStreamingAsr(cfg); err == nil {
			t.Fatalf("NewStreamingAsr(%+v) returned a session, want an error", cfg)
		}
	}
}

// TestStreamingAsrLifecycle: one session, one utterance, release.
func TestStreamingAsrLifecycle(t *testing.T) {
	asr := newTestAsr(t)
	if asr.InFlight() {
		t.Fatalf("a fresh session reports an utterance in flight")
	}
	if got := asr.Finish(); got != "" {
		t.Fatalf("Finish with no utterance open = %q, want empty", got)
	}

	asr.Begin()
	if !asr.InFlight() {
		t.Fatalf("Begin did not open an utterance")
	}
	pcm := synthSpeech(msSamples(1500))
	for i := 0; i+WindowSamples <= len(pcm); i += WindowSamples {
		asr.Feed(pcmToFloat(pcm[i : i+WindowSamples]))
	}
	if fed := asr.Fed(); fed != len(pcm)/WindowSamples*WindowSamples {
		t.Fatalf("Fed = %d samples, want %d", fed, len(pcm)/WindowSamples*WindowSamples)
	}
	partial := asr.Partial()
	final := asr.Finish()
	if final == "" {
		t.Fatalf("a 1.5s utterance produced an empty final transcript")
	}
	if partial == "" {
		t.Logf("note: no partial had appeared yet at 1.5s")
	}
	t.Logf("lifecycle: partial=%q final=%q fed=%d", partial, final, asr.Fed())
	if asr.InFlight() {
		t.Fatalf("Finish left the utterance open")
	}
	if err := asr.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if err := asr.Close(); err != nil {
		t.Fatalf("second Close must be idempotent (C11): %v", err)
	}
}

// TestStreamingAsrKeepsUtterancesApart: a second Begin must not inherit the
// first utterance's decoder state.
func TestStreamingAsrKeepsUtterancesApart(t *testing.T) {
	asr := newTestAsr(t)
	pcm := synthSpeech(msSamples(1200))

	run := func() string {
		asr.Begin()
		for i := 0; i+WindowSamples <= len(pcm); i += WindowSamples {
			asr.Feed(pcmToFloat(pcm[i : i+WindowSamples]))
		}
		return asr.Finish()
	}
	first := run()
	second := run()
	if first == "" || second == "" {
		t.Fatalf("a repeated identical utterance went empty: first=%q second=%q", first, second)
	}
	if first != second {
		t.Fatalf("the same audio decoded to two different texts (%q then %q), so decoder "+
			"state is leaking across utterances", first, second)
	}
	t.Logf("two identical utterances: %q / %q", first, second)
}

// TestStreamingAsrOnSilenceStaysSilent is the recognizer-level half of the
// hallucination control: the model is handed three seconds of digital silence
// with no VAD in the way, and must not invent a word.
func TestStreamingAsrOnSilenceStaysSilent(t *testing.T) {
	asr := newTestAsr(t)
	asr.Begin()
	pcm := silence(msSamples(3000))
	for i := 0; i+WindowSamples <= len(pcm); i += WindowSamples {
		asr.Feed(pcmToFloat(pcm[i : i+WindowSamples]))
		if p := asr.Partial(); p != "" {
			t.Fatalf("partial %q on digital silence at %d samples in: the recognizer is "+
				"hallucinating", p, asr.Fed())
		}
	}
	if final := asr.Finish(); final != "" {
		t.Fatalf("final %q on 3s of digital silence, want empty", final)
	}
}

// TestConsumeDrainsTheSeamUntilTheChannelCloses: the C8 frames go in through
// the real seam and come out as transcripts.
func TestConsumeDrainsTheSeamUntilTheChannelCloses(t *testing.T) {
	tr, _ := newTestTranscriber(t)

	pcm := synthSpeech(msSamples(1500))
	n := len(pcm) / WindowSamples
	// The channel is pre-filled and closed, so Consume can be driven
	// synchronously: the package starts no goroutines, and this test does not
	// need one either to prove the drain works.
	frames := make(chan []byte, n)
	for i := 0; i+WindowSamples <= len(pcm); i += WindowSamples {
		frames <- audio.EncodeFrame(pcm[i : i+WindowSamples])
	}
	close(frames)

	out, err := tr.Consume(t.Context(), frames, nil)
	if err != nil {
		t.Fatalf("Consume: %v", err)
	}
	partials, finals := countTranscripts(out)
	if finals != 1 {
		t.Fatalf("%d finals, want 1 (the stream end must finalize the open utterance)", finals)
	}
	if partials == 0 {
		t.Fatalf("no partials crossed the seam")
	}
	last := out[len(out)-1]
	if last.Kind != TranscriptFinal {
		t.Fatalf("last transcript is a %v, want the final", last.Kind)
	}
	if last.Text == "" {
		t.Fatalf("the final transcript is empty")
	}
	if last.Endpoint.Kind != EndpointEndOfStream {
		t.Fatalf("closing endpoint is %v, want end-of-stream", last.Endpoint.Kind)
	}
	t.Logf("seam: %d partials then final %q over %d sample(s)", partials, last.Text, last.Endpoint.Total)
}

// countTranscripts splits a transcript run into its two kinds.
func countTranscripts(in []Transcript) (partials, finals int) {
	for _, tr := range in {
		switch tr.Kind {
		case TranscriptPartial:
			partials++
		case TranscriptFinal:
			finals++
		}
	}
	return partials, finals
}
