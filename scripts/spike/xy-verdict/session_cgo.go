//go:build cgo_sherpa

package main

// Speech session lifecycle for the unload-test role (cgo build only).

import (
	"errors"
	"path/filepath"
	"time"

	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
)

var errAsrCreate = errors.New("ASR session creation failed")

func openSession(kind, modelsDir string) (*sessionHandle, error) {
	switch kind {
	case "asr":
		return openAsr(modelsDir)
	case "tts", "kws", "vad":
		return nil, nil // measured by model-residency / speech-baseline programs
	}
	return nil, nil
}

func openAsr(modelsDir string) (*sessionHandle, error) {
	enc := filepath.Join(modelsDir, "asr-encoder.int8.onnx")
	dec := filepath.Join(modelsDir, "asr-decoder.int8.onnx")
	tok := filepath.Join(modelsDir, "asr-tokens.txt")

	t0 := time.Now()
	rec := sherpa.NewOnlineRecognizer(&sherpa.OnlineRecognizerConfig{
		FeatConfig:     sherpa.FeatureConfig{SampleRate: 16000, FeatureDim: 80},
		DecodingMethod: "greedy_search",
		ModelConfig: sherpa.OnlineModelConfig{
			Paraformer: sherpa.OnlineParaformerModelConfig{Encoder: enc, Decoder: dec},
			Tokens:     tok,
			NumThreads: 1, // D32: intra_op_num_threads pinned to 1
			Provider:   "cpu",
			Debug:      0,
		},
	})
	loadMs := time.Since(t0).Milliseconds()
	if rec == nil {
		return nil, errAsrCreate
	}

	wt := time.Now()
	stream := sherpa.NewOnlineStream(rec)
	chunk := make([]float32, 512)
	for i := 0; i < 16000/512; i++ { // 1 s of audio
		stream.AcceptWaveform(16000, chunk)
	}
	stream.InputFinished()
	for !rec.IsReady(stream) {
		time.Sleep(5 * time.Millisecond)
	}
	rec.Decode(stream)
	_ = rec.GetResult(stream)
	warmupMs := time.Since(wt).Milliseconds()

	return &sessionHandle{
		loadMs:   loadMs,
		warmupMs: warmupMs,
		dispose: func() {
			sherpa.DeleteOnlineStream(stream)
			sherpa.DeleteOnlineRecognizer(rec)
		},
	}, nil
}
