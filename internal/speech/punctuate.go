package speech

import (
	"fmt"
	"io"
	"os"
	"time"

	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
)

// This file is ticket 27's punctuation stage: the CT-Punc-class CT-Transformer
// that turns the streaming ASR's unpunctuated final transcript into
// dictation-grade text before it reaches the agent.
//
// Where it sits, and why there: asr.go's Begin already says "punctuation and
// stitching is a later stage, ticket 27". This is that stage. It sits on the
// text seam rather than the audio one because the model consumes the
// recognizer's hypothesis, not the waveform -- so it can be swapped, measured
// and degraded independently of whichever engine produced the text.
//
// What it deliberately does not punctuate: partials. They are display-only and
// a hypothesis still growing every 32ms would be re-punctuated for nothing,
// while the budget that matters (300ms from final to punctuated text) is spent
// once per utterance, not once per frame.
//
// Two of the ticket's clauses are structural here rather than documented:
//
//   - A nil stage is the toggle-off path AND the load-failure path, and both
//     degrade to the raw transcript. There is no error return on the seam, so
//     punctuation can never fail an utterance -- the single signal AC#5's retry
//     policy keys on (an empty final) cannot be produced by this stage.
//   - AddPunct returns a string, not (string, error), for the same reason
//     AsrEngine.Finish returns a bare string: at this seam a model that refused
//     the text and a model that had nothing to add are the same event, which is
//     exactly the degrade the ticket asks for.
//
// The interface lives here rather than in engines.go on purpose: engines.go is
// the frozen C9 roster (four names, pinned by c9Rosters), and punctuation is
// not one of them. Adding a fifth name there would move a frozen file for no
// gain, so the seam is declared beside its implementation instead.

// punctuationThreads pins the CT-Transformer to the chain's single intra-op
// thread (D32). It is a package constant, not a knob, for the same reason
// IntraOpThreads is: the CPU story is only true while nothing can raise it.
const punctuationThreads = IntraOpThreads

// PunctuationLatencyBudget is the stage ceiling: from the ASR final to the
// punctuated text, 300ms (D32 16.3.4, ticket 27).
//
// It is a number to measure against, not a deadline enforced with a clock.
// D42#9 forbids wall-clock timeouts, and a stage that abandoned the sentence
// half way would silently drop the punctuation it exists to add -- the failure
// mode would be a plausible-looking unpunctuated transcript, which is the one
// outcome this stage must never manufacture on its own.
const PunctuationLatencyBudget = 300 * time.Millisecond

// punctuationAlphabet is every mark the shipped model can emit, read off its
// own config.yaml punc_list (<unk>, _, comma, full stop, question mark,
// enumeration comma). The ASCII underscore is the model's "no punctuation"
// label and never reaches output; it is listed so the alphabet is the label set
// rather than a guess about it.
//
// A different model changes this constant, not the seam. The positive controls
// use it to state precisely what "no invented punctuation" means.
const punctuationAlphabet = "_\uFF0C\u3002\uFF1F\u3001"

// Punctuator is the punctuation seam: unpunctuated transcript in, punctuated
// transcript out.
//
// The method set is the shipped model's shape, extracted rather than designed,
// for the same reason AsrEngine's is (engines.go): the local implementation is
// the one carrying measurements, so an engine that cannot express this shape is
// the thing that is wrong, not the shape. A cloud punctuation service and an
// English-only model both conform to it unchanged.
//
// There is no Reset and no per-utterance state: the CT-Transformer is stateless
// across calls, which is what lets one resident instance serve every utterance
// without the stitching problem a streaming punctuator would bring.
type Punctuator interface {
	// AddPunct returns text with punctuation restored. It never returns an
	// error: a stage that cannot add punctuation returns its input, and the
	// caller sees an unpunctuated transcript rather than a failed utterance.
	AddPunct(text string) string
	// Close releases the model (C11). Idempotent.
	Close() error
}

// PunctuationConfig points at the CT-Transformer model file.
//
// One path, not three: unlike the streaming Paraformer the punctuator carries
// its token table in the onnx metadata, which is why the model's archive ships
// tokens.json for the training scripts and no separate tokens file for
// inference. Adding a TokensPath field here would be a field nothing reads.
type PunctuationConfig struct {
	// ModelPath is the CT-Transformer onnx file. The product path comes from
	// the model store (build/models/manifest.json, id
	// punc-ct-transformer-zh-en-vocab272727).
	ModelPath string
}

// CTTransformerPunctuator is one resident CT-Punc session: text in, punctuated
// text out.
//
// It owns exactly one sherpa punctuation object for its lifetime and no other
// runtime object, which is what keeps the chain on the single process-wide
// onnxruntime (D5). That is also the memory story: one model, resident for the
// session, released at Close -- not one per utterance (C11: sessions are
// released, never unloaded).
//
// It starts no goroutines and is not safe for concurrent use: one owner
// goroutine drives the voice chain (D38b). It measures no time.
type CTTransformerPunctuator struct {
	punc *sherpa.OfflinePunctuation
	cfg  PunctuationConfig
}

// onnxModelTag is the first byte of a serialized ONNX ModelProto: field 1
// (ir_version) with wire type 0 (varint), so the tag byte is 0x08. Protobuf
// writers emit fields in tag order, so this is the first byte of every ONNX
// model carrying an ir_version -- which is every model ONNX Runtime will load.
// All three models already in this tree start 08 07 (ir_version = 7).
const onnxModelTag = 0x08

// checkOnnxModelFile refuses a file that is not an ONNX model, before anything
// is handed to the C API.
//
// This is not defensive decoration, it is the only place the failure can become
// an error. Handing sherpa a file that is not a model does not return an error
// and does not return nil: onnxruntime throws a C++ exception (0xe06d7363 on
// Windows) and the process dies before Go gets a say -- measured, see
// docs/evidence/s1/27-r1-impl.md. Go cannot recover a C++ throw, so the check
// has to come first.
//
// It is also exactly the shape of the bug this tree already carries:
// third_party/model-fetch/punc-model.int8.onnx is a 29-byte HTTP-401 payload
// ("Invalid username or password."), and a loader without this check aborts the
// host on it -- which is what an interrupted download looks like.
//
// What it cannot see: a well-formed ONNX file that is the wrong model, or one
// corrupted after the header. The model store's manifest hash answers that;
// this answers "the download never happened".
func checkOnnxModelFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("speech: punctuation model %q: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("speech: punctuation model %q is not a regular file", path)
	}
	if info.Size() == 0 {
		return fmt.Errorf("speech: punctuation model %q is empty: a 0-byte file is what an "+
			"interrupted download leaves behind", path)
	}
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("speech: punctuation model %q: %w", path, err)
	}
	defer func() { _ = f.Close() }()
	var head [1]byte
	if _, err := io.ReadFull(f, head[:]); err != nil {
		return fmt.Errorf("speech: punctuation model %q: cannot read its header: %w", path, err)
	}
	if head[0] != onnxModelTag {
		return fmt.Errorf("speech: punctuation model %q is not an ONNX model: it starts with "+
			"byte %#02x (%q) and every ONNX ModelProto starts with %#02x (field 1, ir_version). "+
			"An error page saved under a model name, or a truncated archive, is the usual cause",
			path, head[0], string(head[0]), onnxModelTag)
	}
	return nil
}

// NewCTTransformerPunctuator loads the model and returns a ready stage.
//
// It refuses, rather than crashes, on an onnxruntime older than the linked C
// API gives it (checkOnnxruntime), on a path that is not there, and on a file
// that is not an ONNX model (checkOnnxModelFile -- the C API aborts the process
// on that one, so it is checked before sherpa is given anything).
func NewCTTransformerPunctuator(cfg PunctuationConfig) (*CTTransformerPunctuator, error) {
	if cfg.ModelPath == "" {
		return nil, fmt.Errorf("speech: PunctuationConfig.ModelPath is empty " +
			"(the CT-Transformer onnx file is not optional: a stage with no model " +
			"is not a stage, it is a promise)")
	}
	if err := checkOnnxruntime(); err != nil {
		return nil, err
	}
	if err := checkOnnxModelFile(cfg.ModelPath); err != nil {
		return nil, err
	}
	punc := sherpa.NewOfflinePunctuation(punctuationConfig(cfg))
	if punc == nil {
		return nil, fmt.Errorf("speech: sherpa refused the punctuation config (model=%s): "+
			"the file passes the ONNX header check but is not a CT-Transformer punctuation "+
			"model (the archive's model.int8.onnx, not its tokens.json)", cfg.ModelPath)
	}
	return &CTTransformerPunctuator{punc: punc, cfg: cfg}, nil
}

// punctuationConfig is the one place the model's shape is decided, and it is a
// pure function so the shape can be asserted without a model file -- the same
// arrangement asrRecognizerConfig has, for the same reason (a pin that needs a
// 75MB file to run is a pin that gets skipped).
//
// Three decisions are load-bearing:
//
//   - NumThreads is punctuationThreads, never a parameter (D32).
//   - Debug is 0. The C side prints to stdout when it is not, and a stage that
//     chats during inference is untestable.
//   - Provider is "cpu". The int8 CT-Transformer has no GPU path worth the
//     init cost for 300ms of text work, and one provider for the whole chain is
//     the same decision D5 makes for the runtime.
func punctuationConfig(cfg PunctuationConfig) *sherpa.OfflinePunctuationConfig {
	return &sherpa.OfflinePunctuationConfig{
		Model: sherpa.OfflinePunctuationModelConfig{
			CtTransformer: cfg.ModelPath,
			NumThreads:    punctuationThreads,
			Debug:         0,
			Provider:      "cpu",
		},
	}
}

// The local sherpa implementation of Punctuator, asserted at compile time: if
// *CTTransformerPunctuator stops matching this seam the package stops building.
var _ Punctuator = (*CTTransformerPunctuator)(nil)

// AddPunct returns text with punctuation restored.
//
// Empty and whitespace-only input are returned unchanged without touching the
// model: an empty transcript is not an invitation to punctuate (the final that
// carries it means "nothing sayable", AC#5), and asking a model what mark
// follows nothing is how hallucinated marks get into a transcript.
//
// A nil receiver or a released session returns text unchanged, so a caller that
// kept the stage across a Close gets its transcript back rather than a crash.
func (p *CTTransformerPunctuator) AddPunct(text string) string {
	if p == nil || p.punc == nil || isBlank(text) {
		return text
	}
	return p.punc.AddPunct(text)
}

// Close releases the model. It is idempotent (C11).
func (p *CTTransformerPunctuator) Close() error {
	if p == nil || p.punc == nil {
		return nil
	}
	sherpa.DeleteOfflinePunc(p.punc)
	p.punc = nil
	return nil
}

// OpenPunctuation is the stage's whole lifecycle in one call: it returns a
// ready Punctuator, or nil when the stage is switched off.
//
// The nil,nil return is the contract that makes both of the ticket's degrade
// paths the same path. Switched off and failed to load are indistinguishable to
// every caller above this line, and both mean "hand the raw transcript on" --
// which is why an utterance can never fail because punctuation did. The error
// is still returned on a load failure so the caller can log it and show the UI
// hint the ticket asks for; it is the caller's decision to keep going, and
// every caller should.
//
// Note the explicit nil: returning the typed nil *CTTransformerPunctuator would
// produce a non-nil interface and defeat the whole arrangement.
func OpenPunctuation(cfg PunctuationConfig, enabled bool) (Punctuator, error) {
	if !enabled {
		return nil, nil
	}
	p, err := NewCTTransformerPunctuator(cfg)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// isBlank reports whether s has nothing a punctuator could punctuate.
func isBlank(s string) bool {
	for _, r := range s {
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' && r != '\u3000' {
			return false
		}
	}
	return true
}
