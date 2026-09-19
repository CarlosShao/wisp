// speech-baseline measures RSS baselines 2/3 (cgo binary, sherpa-onnx):
//
//	(2) onnxruntime + sherpa-onnx DLLs mapped at process start, no session
//	(3) + one sherpa session (KWS zipformer wenetspeech 3.3M int8)
//
// Because the Go binding statically imports the C API, the DLLs are mapped
// by the loader at process start - which is exactly path Y's idle shape.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/CarlosShao/wisp/scripts/spike/common"
	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
)

type stage struct {
	Stage           string  `json:"stage"`
	PrivateWS_MB    float64 `json:"privateWS_MB"`
	SharedWS_MB     float64 `json:"sharedWS_MB"`
	TotalWS_MB      float64 `json:"totalWS_MB"`
	PrivateCommitMB float64 `json:"privateCommit_MB"`
	GDIObjects      uint32  `json:"gdiObjects"`
	UserObjects     uint32  `json:"userObjects"`
	Handles         uint32  `json:"handles"`
	Threads         uint32  `json:"threads"`
	DeltaPrev_MB    float64 `json:"deltaPrev_MB"`
	MinPrivateWSMB  float64 `json:"minPrivateWS_MB"`
	MaxPrivateWSMB  float64 `json:"maxPrivateWS_MB"`
	Error           string  `json:"error,omitempty"`
	DurationMs      int64   `json:"stageDurationMs,omitempty"`
}

type result struct {
	Program       string             `json:"program"`
	Kind          string             `json:"kind"`
	Machine       common.MachineInfo `json:"machine"`
	StartedAt     string             `json:"startedAtUtc"`
	SherpaVersion string             `json:"sherpaVersion"`
	OnnxrtVersion string             `json:"onnxruntimeVersion"`
	Model         string             `json:"model"`
	ModelFileMB   float64            `json:"modelFilesMB"`
	LoadMs        int64              `json:"kwsLoadMs,omitempty"`
	Stages        []stage            `json:"stages"`
	GOGC          string             `json:"gogc"`
	Note          string             `json:"note"`
}

func mbOf(b uint64) float64 { return float64(b) / (1 << 20) }

func sampleStage(name string, prev *stage, run func() error) stage {
	st := stage{Stage: name}
	t0 := time.Now()
	var runErr error
	if run != nil {
		runErr = run()
	}
	if runErr != nil {
		st.Error = runErr.Error()
	}
	common.SettleGC()
	minS, medS, maxS := common.SampleStable(5, 200)
	st.PrivateWS_MB = mbOf(medS.PrivateWorkingSet)
	st.SharedWS_MB = mbOf(medS.SharedWorkingSet)
	st.TotalWS_MB = mbOf(medS.WorkingSet)
	st.PrivateCommitMB = mbOf(medS.PrivateCommit)
	st.GDIObjects = medS.GDIObjects
	st.UserObjects = medS.UserObjects
	st.Handles = medS.Handles
	st.Threads = medS.Threads
	st.MinPrivateWSMB = mbOf(minS.PrivateWorkingSet)
	st.MaxPrivateWSMB = mbOf(maxS.PrivateWorkingSet)
	if prev != nil {
		st.DeltaPrev_MB = st.PrivateWS_MB - prev.PrivateWS_MB
	}
	st.DurationMs = time.Since(t0).Milliseconds()
	return st
}

func dirSizeMB(paths ...string) float64 {
	var total int64
	for _, p := range paths {
		if fi, err := os.Stat(p); err == nil {
			total += fi.Size()
		}
	}
	return float64(total) / (1 << 20)
}

func main() {
	modelsDir := flag.String("models", "third_party/spike-models", "spike models dir")
	out := flag.String("out", "", "path to write JSON (also printed to stdout)")
	flag.Parse()

	res := result{
		Program:   "speech-baseline",
		Kind:      "cgo",
		Machine:   common.GetMachineInfo(),
		StartedAt: time.Now().UTC().Format(time.RFC3339),
		GOGC:      os.Getenv("GOGC"),
	}
	res.SherpaVersion = sherpa.GetVersion()
	res.OnnxrtVersion = sherpa.GetOnnxruntimeVersion()

	kwsDir := filepath.Join(*modelsDir, "sherpa-onnx-kws-zipformer-wenetspeech-3.3M-2024-01-01")
	enc := filepath.Join(kwsDir, "encoder-epoch-12-avg-2-chunk-16-left-64.int8.onnx")
	dec := filepath.Join(kwsDir, "decoder-epoch-12-avg-2-chunk-16-left-64.int8.onnx")
	joi := filepath.Join(kwsDir, "joiner-epoch-12-avg-2-chunk-16-left-64.int8.onnx")
	tok := filepath.Join(kwsDir, "tokens.txt")
	kw := filepath.Join(kwsDir, "keywords.txt")
	res.Model = "sherpa-onnx-kws-zipformer-wenetspeech-3.3M-2024-01-01 (int8)"
	res.ModelFileMB = dirSizeMB(enc, dec, joi, tok)

	// Stage 2: DLLs mapped by the loader, no session created.
	res.Stages = append(res.Stages, sampleStage("2-dll-mapped-no-session", nil, nil))
	prev := &res.Stages[0]

	// Stage 3: + one sherpa session (KWS zipformer 3.3M int8).
	var kws *sherpa.KeywordSpotter
	var stream *sherpa.OnlineStream
	st3 := sampleStage("3-kws-session", prev, func() error {
		t0 := time.Now()
		k := sherpa.NewKeywordSpotter(&sherpa.KeywordSpotterConfig{
			FeatConfig: sherpa.FeatureConfig{SampleRate: 16000, FeatureDim: 80},
			ModelConfig: sherpa.OnlineModelConfig{
				Transducer: sherpa.OnlineTransducerModelConfig{Encoder: enc, Decoder: dec, Joiner: joi},
				Tokens:     tok,
				NumThreads: 1, // D32: intra_op_num_threads pinned to 1
				Provider:   "cpu",
				Debug:      0,
			},
			KeywordsFile:      kw,
			KeywordsThreshold: 0.25,
			KeywordsScore:     2.0,
			MaxActivePaths:    4,
		})
		res.LoadMs = time.Since(t0).Milliseconds()
		if k == nil {
			return fmt.Errorf("NewKeywordSpotter returned nil")
		}
		kws = k
		stream = sherpa.NewKeywordStream(kws)
		// feed 3 s of audio in 512-sample chunks and decode (realistic Armed state)
		chunk := make([]float32, 512)
		for i := 0; i < 16000*3/512; i++ {
			stream.AcceptWaveform(16000, chunk)
			for kws.IsReady(stream) {
				kws.Decode(stream)
			}
		}
		return nil
	})
	res.Stages = append(res.Stages, st3)
	prev = &res.Stages[1]

	// Stage 3b: dispose + FreeOSMemory (the C11 unload primitive).
	st3b := sampleStage("3b-after-dispose-freeosmemory", prev, func() error {
		if stream != nil {
			sherpa.DeleteOnlineStream(stream)
		}
		if kws != nil {
			sherpa.DeleteKeywordSpotter(kws)
		}
		return nil
	})
	res.Stages = append(res.Stages, st3b)

	b, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(b))
	if *out != "" {
		os.WriteFile(*out, b, 0644)
	}
}
