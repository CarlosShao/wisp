package speech

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"

	"github.com/CarlosShao/wisp/internal/audio"
)

// IntraOpThreads is the onnxruntime intra-op thread count the whole voice
// chain runs at: one, pinned (D32, SPEC-04 sec 4, ticket 15 AC#6). It is a
// package constant rather than a knob because the SLO story ("ASR inference
// single-core pinned, no thread-pool fan-out") is only true while nothing can
// raise it.
const IntraOpThreads = 1

// requiredORTMinor is the onnxruntime release the linked sherpa-onnx C API
// asks for. sherpa-onnx 1.13.8 calls OrtGetApiBase()->GetApi(28); ORT 1.28 is
// the release that answers it (deps.toml pins [onnxruntime] 1.28.2). An older
// runtime does not return an error, it returns NULL and the C API dereferences
// it, so the version is checked here, before any session is created (D37b:
// the cgo boundary is where a bad runtime must be caught, not inside C).
const requiredORTMinor = 28

// vadBufferSeconds is the detector's own ring size. It must exceed
// MaxUtterance so a 60s utterance can still be finalized; 120s leaves room for
// one utterance plus the trailing silence that closes it.
const vadBufferSeconds = 120

// DefaultVadThreshold is silero's speech-probability cut when the config
// leaves Threshold zero. 0.5 is the model's published default.
const DefaultVadThreshold float32 = 0.5

// DetectorConfig configures a Detector. There is deliberately no thread-count
// field: see IntraOpThreads.
type DetectorConfig struct {
	// ModelPath is the silero VAD onnx file (third_party/spike-models/
	// silero_vad.onnx in this tree; the product path comes from the model
	// store, SPEC-04 sec 8).
	ModelPath string
	// Threshold overrides DefaultVadThreshold when non-zero.
	Threshold float32
}

// Detector is silero VAD plus the endpoint state machine: frames in, endpoint
// transitions out. It owns one sherpa session for its lifetime (C11: sessions
// are released, never unloaded) and no other runtime object, which is what
// keeps the whole chain on the single process-wide onnxruntime (D5).
//
// Detector is not safe for concurrent use: one owner goroutine feeds it
// (D38b). It starts no goroutines of its own, so it cannot outlive its owner.
type Detector struct {
	vad *sherpa.VoiceActivityDetector
	ep  *Endpointer
	cfg DetectorConfig

	rest []float32 // input not yet filling one WindowSamples window
}

// NewDetector loads the silero model and returns a ready detector.
//
// It refuses, rather than crashes, when the process is running on an
// onnxruntime older than the linked C API: see requiredORTMinor.
func NewDetector(cfg DetectorConfig) (*Detector, error) {
	if cfg.ModelPath == "" {
		return nil, fmt.Errorf("speech: DetectorConfig.ModelPath is empty")
	}
	if err := checkOnnxruntime(); err != nil {
		return nil, err
	}
	thr := cfg.Threshold
	if thr == 0 {
		thr = DefaultVadThreshold
	}
	vad := sherpa.NewVoiceActivityDetector(&sherpa.VadModelConfig{
		SileroVad: sherpa.SileroVadModelConfig{
			Model:              cfg.ModelPath,
			Threshold:          thr,
			MinSilenceDuration: 0, // the Go state machine owns EndpointSilence
			MinSpeechDuration:  0, // ... and MinSpeech; the C side is a classifier
			WindowSize:         WindowSamples,
			MaxSpeechDuration:  float32(MaxUtterance.Seconds()),
		},
		SampleRate: audio.TargetRate,
		NumThreads: IntraOpThreads,
		Provider:   "cpu",
	}, vadBufferSeconds)
	if vad == nil {
		return nil, fmt.Errorf("speech: sherpa refused the VAD config (model=%s): "+
			"the file is unreadable or is not a silero VAD model", cfg.ModelPath)
	}
	return &Detector{vad: vad, ep: NewEndpointer(), cfg: cfg}, nil
}

// checkOnnxruntime rejects an onnxruntime older than the C API expects, which
// is how a shadowed onnxruntime.dll (for example the OS-provided copy in
// C:\Windows\System32, which the Windows loader prefers over PATH) shows up.
func checkOnnxruntime() error {
	ver := sherpa.GetOnnxruntimeVersion()
	major, minor, ok := parseMajorMinor(ver)
	if !ok {
		return fmt.Errorf("speech: cannot parse onnxruntime version %q", ver)
	}
	if major > 1 || (major == 1 && minor >= requiredORTMinor) {
		return nil
	}
	return fmt.Errorf("speech: onnxruntime %s is older than the 1.%d the linked sherpa-onnx "+
		"C API requires; a second onnxruntime.dll is shadowing the one on PATH "+
		"(on Windows the loader searches System32 before PATH)", ver, requiredORTMinor)
}

// parseMajorMinor reads the leading "MAJOR.MINOR" of a version string.
func parseMajorMinor(v string) (major, minor int, ok bool) {
	parts := strings.SplitN(v, ".", 3)
	if len(parts) < 2 {
		return 0, 0, false
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, false
	}
	minor, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, false
	}
	return major, minor, true
}

// Accept feeds samples (any length, any chunking) and appends the endpoint
// transitions they produce. Samples are buffered until they fill a whole
// silero window, because the C API drops a partial tail.
func (d *Detector) Accept(samples []float32, out []Endpoint) []Endpoint {
	for len(samples) > 0 {
		need := WindowSamples - len(d.rest)
		if need > len(samples) {
			need = len(samples)
		}
		d.rest = append(d.rest, samples[:need]...)
		samples = samples[need:]
		if len(d.rest) < WindowSamples {
			return out
		}
		out = d.acceptWindow(d.rest, out)
		d.rest = d.rest[:0]
	}
	return out
}

// acceptWindow runs one window through silero and the state machine.
func (d *Detector) acceptWindow(win []float32, out []Endpoint) []Endpoint {
	d.vad.AcceptWaveform(win)
	isSpeech := d.vad.IsSpeech()
	// Drop sherpa's own segment queue: this package owns segmentation, and an
	// undrained queue would keep a second copy of every utterance alive.
	for !d.vad.IsEmpty() {
		d.vad.Pop()
	}
	return d.ep.Accept(win, isSpeech, out)
}

// Finalize zero-pads a partial trailing window, classifies it, then closes any
// open utterance as EndpointEndOfStream. Call it once when the stream ends.
func (d *Detector) Finalize(out []Endpoint) []Endpoint {
	if n := len(d.rest); n > 0 {
		pad := make([]float32, WindowSamples)
		copy(pad, d.rest)
		out = d.acceptWindow(pad, out)
		d.rest = d.rest[:0]
	}
	return d.ep.Finalize(out)
}

// Consume drains seam frames until the channel closes or ctx is cancelled,
// decoding each to samples and feeding the detector. It runs in the caller's
// goroutine by design: the voice chain starts no detached goroutines (D38b),
// so cancellation is the caller's ctx and shutdown is the caller's wait.
//
// A closed channel is end of stream: the open utterance is finalized.
func (d *Detector) Consume(ctx context.Context, frames <-chan []byte, out []Endpoint) ([]Endpoint, error) {
	for {
		select {
		case <-ctx.Done():
			return out, ctx.Err()
		case frame, ok := <-frames:
			if !ok {
				return d.Finalize(out), nil
			}
			pcm, err := audio.DecodeFrame(frame)
			if err != nil {
				return out, fmt.Errorf("speech: %w", err)
			}
			out = d.Accept(pcmToFloat(pcm), out)
		}
	}
}

// Close releases the sherpa session. It is idempotent (C11).
func (d *Detector) Close() error {
	if d.vad != nil {
		sherpa.DeleteVoiceActivityDetector(d.vad)
		d.vad = nil
	}
	return nil
}

// pcmToFloat converts seam samples to the float32 silero wants. int16 full
// scale maps to +-1.0.
func pcmToFloat(pcm []int16) []float32 {
	out := make([]float32, len(pcm))
	for i, s := range pcm {
		out[i] = float32(s) / 32768
	}
	return out
}
