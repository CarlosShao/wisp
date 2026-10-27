package speech

import (
	"fmt"
	"os"
	"path/filepath"

	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"

	"github.com/CarlosShao/wisp/internal/audio"
)

// Ticket 41 (engine leg): the local implementation of the C9 wake-word seam.
// The seam itself -- WakeWordEngine, KeywordHit and the [Accept Close] roster
// -- is frozen in engines.go and is not touched here: KwsEngine only conforms
// to it, and the compile-time assertion below is what keeps that true.
//
// The model is the pinned kws-zipformer-wenetspeech 3.3M transducer
// (models/manifest.json id kws-zipformer-wenetspeech-3.3M-2024-01-01, license
// Apache-2.0, quant int8, status ok). Like every engine in this package it
// crosses cgo exactly once per process: one KeywordSpotter plus one keyword
// stream, both released by Close (C11: sessions are released, never unloaded,
// because D5 keeps one onnxruntime for the whole chain).
//
// The one sherpa clause this file exists to honor is verbatim in the binding
// (sherpa_onnx.go:2466): "You MUST call it right after detecting a keyword" --
// Reset. Accept calls it on every hit, and the test suite proves the line is
// load-bearing by removing it and watching the engine go red.

// kwsModelDirName is the pinned model directory under third_party/spike-models.
// The product path comes from the model store (SPEC-04 sec 8); until that
// exists the package constant is the path, exactly as 15-b and 15-c decided
// for ASR and VAD.
const kwsModelDirName = "sherpa-onnx-kws-zipformer-wenetspeech-3.3M-2024-01-01"

// The five manifest files the spotter is built from, int8 quantization
// (manifest archive.files; each file's sha256 is pinned there and was verified
// bit-for-bit before this leg landed).
const (
	kwsEncoderFile  = "encoder-epoch-99-avg-1-chunk-16-left-64.int8.onnx"
	kwsDecoderFile  = "decoder-epoch-99-avg-1-chunk-16-left-64.int8.onnx"
	kwsJoinerFile   = "joiner-epoch-99-avg-1-chunk-16-left-64.int8.onnx"
	kwsTokensFile   = "tokens.txt"
	kwsKeywordsFile = "keywords.txt"
)

// kwsFeatureDim is the mel feature width the KWS zipformer is trained on (80,
// the width every sherpa streaming model shares). It belongs to the model
// file, not to configuration.
const kwsFeatureDim = 80

// The keyword search's two knobs, as package constants rather than
// configuration keys (thresholds and word lists eat constants until a wiring
// leg consumes the ticket-255 keys; a config key this package cannot honor is
// debt, not flexibility). The values are the ones the sherpa-onnx c-api
// example documents for this checkpoint family (keywords_score 3.0,
// keywords_threshold 0.1) AND the ones this host measured against the
// checkpoint's own test_wavs: at the API defaults (1.0/0.25) the same audio
// scores zero hits everywhere, at 3.0/0.1 every demo keyword fires and the
// negative controls stay silent. Raising the threshold trades misses for
// false triggers; ticket 41's P14 report is where a measured production value
// would replace this one.
const (
	// DefaultKwsThreshold is the keyword probability cut.
	DefaultKwsThreshold float32 = 0.1
	// DefaultKwsScore is the boosting score applied to every keyword in the
	// table.
	DefaultKwsScore float32 = 3.0
	// DefaultKwsMaxActivePaths is the beam width of the keyword search (the
	// c-api example's value; not re-measured).
	DefaultKwsMaxActivePaths = 4
)

// kwsVetoWords is the veto-word set KeywordHit.Veto keys on. Both kinds of
// word live in ONE keywords file (ticket 41: veto words share the same model
// for the B1 voice-cancel channel), so the classification is data, not code:
// any display word listed here that the spotter fires is reported with
// Veto=true. The shipped table (manifest keywords sha256 46de8c68) contains
// eight wake words and no veto word, so the set is empty today; the B leg adds
// veto words by editing the keywords file, and they are classified the moment
// they appear here. An engine hit can therefore never claim a veto word that
// is not in the table.
var kwsVetoWords = map[string]bool{}

// KwsConfig points at the KWS transducer triple, its token table and its
// keyword list. There is deliberately no thread-count field (see
// IntraOpThreads) and no keywords-buffer field (the file is the source of
// truth).
type KwsConfig struct {
	// EncoderPath is the streaming transducer encoder (int8).
	EncoderPath string
	// DecoderPath is the streaming transducer decoder (int8).
	DecoderPath string
	// JoinerPath is the streaming transducer joiner (int8).
	JoinerPath string
	// TokensPath is the token table (tokens.txt).
	TokensPath string
	// KeywordsPath is the keyword table (keywords.txt, one pinyin-spelled
	// keyword per line, display form after '@').
	KeywordsPath string
	// Threshold overrides DefaultKwsThreshold when non-zero.
	Threshold float32
}

// DefaultKwsConfig resolves the pinned model directory and returns the
// production config pointed at its five manifest files. The directory is found
// from WISP_SPIKE_MODELS -- which may name either the spike-models root or the
// KWS directory itself -- and otherwise by walking up from the working
// directory, so the same code runs under "go test" (cwd = the package
// directory), from a staging directory that holds the sherpa DLLs, and later
// from the composition root.
func DefaultKwsConfig() (KwsConfig, error) {
	dir, err := kwsModelDir()
	if err != nil {
		return KwsConfig{}, err
	}
	return KwsConfig{
		EncoderPath:  filepath.Join(dir, kwsEncoderFile),
		DecoderPath:  filepath.Join(dir, kwsDecoderFile),
		JoinerPath:   filepath.Join(dir, kwsJoinerFile),
		TokensPath:   filepath.Join(dir, kwsTokensFile),
		KeywordsPath: filepath.Join(dir, kwsKeywordsFile),
	}, nil
}

// kwsModelDir locates the pinned KWS model directory. WISP_SPIKE_MODELS wins
// when it is set, and it is accepted in either spelling (the spike-models root
// or the KWS directory itself); otherwise the walk-up finds a repository
// checkout the way spikeModelsDir does for ASR.
func kwsModelDir() (string, error) {
	if p := os.Getenv("WISP_SPIKE_MODELS"); p != "" {
		if isKwsModelDir(p) {
			return p, nil
		}
		inner := filepath.Join(p, kwsModelDirName)
		if isKwsModelDir(inner) {
			return inner, nil
		}
		return "", fmt.Errorf("speech: WISP_SPIKE_MODELS=%s holds no %s (looked for %s in it and "+
			"one level down); point it at the spike-models root or at the KWS directory", p, kwsKeywordsFile, kwsKeywordsFile)
	}
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("speech: getwd: %w", err)
	}
	return kwsWalkUpFrom(dir)
}

// kwsWalkUpFrom looks for the pinned model directory in start and up to seven
// levels above it. Split out so the walk is testable from a directory the
// caller names (under staging the working directory is not in the repository).
func kwsWalkUpFrom(start string) (string, error) {
	rel := filepath.Join("third_party", "spike-models", kwsModelDirName)
	dir := start
	for i := 0; i < 8; i++ {
		if isKwsModelDir(filepath.Join(dir, rel)) {
			return filepath.Join(dir, rel), nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("speech: KWS model directory not found: set WISP_SPIKE_MODELS, or run "+
		"from inside the repository (looked for %s)", rel)
}

// isKwsModelDir reports whether dir holds the keyword table -- the one file
// the spotter cannot open without.
func isKwsModelDir(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, kwsKeywordsFile))
	return err == nil
}

// KwsEngine is the local wake-word engine: samples in, keyword hits out, for
// as long as the engine is alive. It implements WakeWordEngine.
//
// It owns exactly one sherpa KeywordSpotter and one keyword stream for its
// lifetime. That is the whole cgo surface of the KWS leg: no other runtime
// object is created, which keeps the chain on the single process-wide
// onnxruntime (D5).
//
// KwsEngine starts no goroutines and is not safe for concurrent use: one owner
// goroutine feeds it (D38b). It never measures wall-clock time; At comes off
// the audio clock (D42#9), accumulated over everything the engine has been
// fed, not over one utterance -- the caller's clock, not the engine's.
type KwsEngine struct {
	spotter *sherpa.KeywordSpotter
	stream  *sherpa.OnlineStream
	cfg     KwsConfig
	fed     int // samples fed since the engine opened (the At clock)
}

// The compile-time assertion: if *KwsEngine stops matching the frozen C9
// wake-word seam the package stops building, which is the only kind of
// conformance check that cannot rot.
var _ WakeWordEngine = (*KwsEngine)(nil)

// NewKwsEngine loads the transducer triple and opens the keyword stream.
//
// It refuses, rather than crashes, when the process is running on an
// onnxruntime older than the linked C API gives it (see checkOnnxruntime) and
// when sherpa refuses the config -- NewKeywordSpotter returns nil instead of
// an error, so the nil is turned into a named error here and the caller never
// sees a half-built engine.
func NewKwsEngine(cfg KwsConfig) (*KwsEngine, error) {
	if cfg.EncoderPath == "" || cfg.DecoderPath == "" || cfg.JoinerPath == "" ||
		cfg.TokensPath == "" || cfg.KeywordsPath == "" {
		return nil, fmt.Errorf("speech: KwsConfig needs EncoderPath, DecoderPath, JoinerPath, " +
			"TokensPath and KeywordsPath (the KWS transducer triple plus its token and keyword tables)")
	}
	if err := checkOnnxruntime(); err != nil {
		return nil, err
	}
	spotter := sherpa.NewKeywordSpotter(kwsSpotterConfig(cfg))
	if spotter == nil {
		return nil, fmt.Errorf("speech: sherpa refused the KWS config (encoder=%s decoder=%s joiner=%s "+
			"tokens=%s keywords=%s): a file is unreadable, or it is not a KWS zipformer transducer",
			cfg.EncoderPath, cfg.DecoderPath, cfg.JoinerPath, cfg.TokensPath, cfg.KeywordsPath)
	}
	stream := sherpa.NewKeywordStream(spotter)
	if stream == nil {
		sherpa.DeleteKeywordSpotter(spotter)
		return nil, fmt.Errorf("speech: sherpa refused to open a KWS stream (keywords=%s)", cfg.KeywordsPath)
	}
	return &KwsEngine{spotter: spotter, stream: stream, cfg: cfg}, nil
}

// kwsSpotterConfig is the one place the spotter's shape is decided, and it is
// a pure function so the shape can be asserted without a model file.
//
// Load-bearing decisions:
//
//   - NumThreads is IntraOpThreads, never a parameter (AC#6, D32), the same
//     pin every engine in this package runs under.
//   - The keyword table is the file, and the three search knobs are the
//     package constants above -- nothing here reads configuration.
func kwsSpotterConfig(cfg KwsConfig) *sherpa.KeywordSpotterConfig {
	thr := cfg.Threshold
	if thr == 0 {
		thr = DefaultKwsThreshold
	}
	return &sherpa.KeywordSpotterConfig{
		FeatConfig: sherpa.FeatureConfig{
			SampleRate: audio.TargetRate,
			FeatureDim: kwsFeatureDim,
		},
		ModelConfig: sherpa.OnlineModelConfig{
			Transducer: sherpa.OnlineTransducerModelConfig{
				Encoder: cfg.EncoderPath,
				Decoder: cfg.DecoderPath,
				Joiner:  cfg.JoinerPath,
			},
			Tokens:     cfg.TokensPath,
			NumThreads: IntraOpThreads,
			Provider:   "cpu",
		},
		MaxActivePaths:    DefaultKwsMaxActivePaths,
		KeywordsFile:      cfg.KeywordsPath,
		KeywordsScore:     DefaultKwsScore,
		KeywordsThreshold: thr,
	}
}

// Accept feeds samples (16kHz mono float32, any length, any chunking) and
// appends the keyword hits they triggered. It is a no-op after Close, so a
// caller that races its own shutdown degrades to silence instead of panicking
// (the same contract StreamingAsr.Feed keeps).
//
// The decode loop is sherpa's own recipe for the keyword spotter: decode while
// ready, read the result after each decode, and -- the clause the binding
// states verbatim -- Reset the stream right after a keyword fires, so the
// search state goes back to the start of the table and the next utterance is
// searched from scratch. A hit's At is the running sample count at the decode
// that fired it: the audio clock, accumulated across the engine's whole life.
func (k *KwsEngine) Accept(samples []float32, out []KeywordHit) []KeywordHit {
	if k.spotter == nil || k.stream == nil || len(samples) == 0 {
		return out
	}
	k.stream.AcceptWaveform(audio.TargetRate, samples)
	k.fed += len(samples)
	for k.spotter.IsReady(k.stream) {
		k.spotter.Decode(k.stream)
		res := k.spotter.GetResult(k.stream)
		if res == nil || res.Keyword == "" {
			continue
		}
		out = append(out, KeywordHit{
			Word: res.Keyword,
			Veto: kwsVetoWords[res.Keyword],
			At:   k.fed,
		})
		k.spotter.Reset(k.stream)
	}
	return out
}

// Close releases the stream and the spotter. It is idempotent (C11): the
// second call finds both nil and does nothing.
func (k *KwsEngine) Close() error {
	if k.stream != nil {
		sherpa.DeleteOnlineStream(k.stream)
		k.stream = nil
	}
	if k.spotter != nil {
		sherpa.DeleteKeywordSpotter(k.spotter)
		k.spotter = nil
	}
	return nil
}
