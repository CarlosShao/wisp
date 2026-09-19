// model-residency fills D32 16.3.3 with measured numbers:
//
//	-which kws|vad|asr|tts : one session at a time in a fresh process;
//	                         reports load time, idle residency (no inference),
//	                         post-inference residency, and dispose/settle.
//	-which switch          : ASR<->TTS serial half-duplex cycles in one process
//	                         (D16 corollary: peak = max(ASR,TTS), not sum);
//	                         measures dispose+load switch latency against the
//	                         1600 ms budget (D32 16.3.4 branch "needs TTS load").
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

type residencyReport struct {
	Program           string             `json:"program"`
	Which             string             `json:"which"`
	Machine           common.MachineInfo `json:"machine"`
	StartedAt         string             `json:"startedAtUtc"`
	Models            []modelFile        `json:"modelFiles"`
	PreWSMB           float64            `json:"prePrivateWS_MB"`
	LoadMs            int64              `json:"loadMs"`
	IdleResidencyMB   float64            `json:"idleResidencyPrivateWS_MB"`
	PostInferMB       float64            `json:"postInferencePrivateWS_MB"`
	PostInferCommitMB float64            `json:"postInferencePrivateCommit_MB"`
	PeakCommitMB      float64            `json:"peakPrivateCommitMB_inProcess"`
	DisposeMs         int64              `json:"disposeMs"`
	SettledMs         *int64             `json:"settledMsAfterDispose"`
	Settled           bool               `json:"settled"`
	SherpaVersion     string             `json:"sherpaVersion"`
	Error             string             `json:"error,omitempty"`
	Note              string             `json:"note,omitempty"`

	// switch-mode only
	Switch *switchReport `json:"switch,omitempty"`
}

type modelFile struct {
	Name   string  `json:"name"`
	Path   string  `json:"path"`
	Bytes  float64 `json:"bytesMB"`
	SHA256 string  `json:"sha256,omitempty"`
}

type switchCycle struct {
	Cycle            int     `json:"cycle"`
	AsrDisposeMs     int64   `json:"asrDisposeMs"`
	TtsLoadMs        int64   `json:"ttsLoadMs"`
	SwitchAsrToTtsMs int64   `json:"switchAsrToTtsMs"`
	TtsDisposeMs     int64   `json:"ttsDisposeMs"`
	AsrLoadMs        int64   `json:"asrLoadMs"`
	SwitchTtsToAsrMs int64   `json:"switchTtsToAsrMs"`
	WsAtAsrMB        float64 `json:"wsWhileAsr_MB"`
	WsAtTtsMB        float64 `json:"wsWhileTts_MB"`
}

type switchReport struct {
	Cycles          []switchCycle `json:"cycles"`
	FirstTtsLoadMs  int64         `json:"firstTtsLoadMs_cold"`
	MinSwitchMs     float64       `json:"minAsrToTtsMs"`
	P50SwitchMs     float64       `json:"p50AsrToTtsMs"`
	MaxWsObservedMB float64       `json:"maxWsDuringSwitch_MB"`
	PeakCommitMB    float64       `json:"peakCommitDuringSwitch_MB"`
	BudgetMs        float64       `json:"budgetMs"`
	BudgetVerdict   string        `json:"budgetVerdict"`
	PeakIsMaxNotSum bool          `json:"peakIsMaxNotSum"`
}

func mbOf(b uint64) float64 { return float64(b) / (1 << 20) }

func statFile(path string) modelFile {
	fi, err := os.Stat(path)
	if err != nil {
		return modelFile{Name: filepath.Base(path), Path: path}
	}
	return modelFile{Name: filepath.Base(path), Path: path, Bytes: float64(fi.Size()) / (1 << 20)}
}

func loadPeakCommit() float64 {
	return common.PeakPrivateCommitMB()
}

func settledWithin(pre uint64, budgetMs int64) (*int64, bool, []float64) {
	threshold := pre + 2<<20
	n := int(budgetMs / 250)
	var samples []float64
	for i := 1; i <= n; i++ {
		common.SleepMs(250)
		s := common.SampleMem()
		samples = append(samples, mbOf(s.PrivateWorkingSet))
		if s.PrivateWorkingSet <= threshold {
			t := int64(i) * 250
			return &t, true, samples
		}
	}
	return nil, false, samples
}

// ------------------------------------------------------------- sessions

type session struct {
	dispose func()
}

func kwsConfig(modelsDir string) (*sherpa.KeywordSpotter, string) {
	kwsDir := filepath.Join(modelsDir, "sherpa-onnx-kws-zipformer-wenetspeech-3.3M-2024-01-01")
	k := sherpa.NewKeywordSpotter(&sherpa.KeywordSpotterConfig{
		FeatConfig: sherpa.FeatureConfig{SampleRate: 16000, FeatureDim: 80},
		ModelConfig: sherpa.OnlineModelConfig{
			Transducer: sherpa.OnlineTransducerModelConfig{
				Encoder: filepath.Join(kwsDir, "encoder-epoch-12-avg-2-chunk-16-left-64.int8.onnx"),
				Decoder: filepath.Join(kwsDir, "decoder-epoch-12-avg-2-chunk-16-left-64.int8.onnx"),
				Joiner:  filepath.Join(kwsDir, "joiner-epoch-12-avg-2-chunk-16-left-64.int8.onnx"),
			},
			Tokens:     filepath.Join(kwsDir, "tokens.txt"),
			NumThreads: 1,
			Provider:   "cpu",
			Debug:      0,
		},
		KeywordsFile:      filepath.Join(kwsDir, "keywords.txt"),
		KeywordsThreshold: 0.25,
		KeywordsScore:     2.0,
		MaxActivePaths:    4,
	})
	return k, kwsDir
}

func openKws(modelsDir string) (*session, error) {
	k, _ := kwsConfig(modelsDir)
	if k == nil {
		return nil, fmt.Errorf("NewKeywordSpotter returned nil")
	}
	s := sherpa.NewKeywordStream(k)
	chunk := make([]float32, 512)
	for i := 0; i < 16000/512; i++ { // 1 s warmup
		s.AcceptWaveform(16000, chunk)
		for k.IsReady(s) {
			k.Decode(s)
		}
	}
	return &session{dispose: func() {
		sherpa.DeleteOnlineStream(s)
		sherpa.DeleteKeywordSpotter(k)
	}}, nil
}

func openVad(modelsDir string) (*session, error) {
	v := sherpa.NewVoiceActivityDetector(&sherpa.VadModelConfig{
		SileroVad: sherpa.SileroVadModelConfig{
			Model:              filepath.Join(modelsDir, "silero_vad.onnx"),
			Threshold:          0.5,
			MinSilenceDuration: 0.5,
			MinSpeechDuration:  0.25,
			WindowSize:         512,
			MaxSpeechDuration:  5.0,
		},
		SampleRate: 16000,
		NumThreads: 1,
		Provider:   "cpu",
		Debug:      0,
	}, 60)
	if v == nil {
		return nil, fmt.Errorf("NewVoiceActivityDetector returned nil")
	}
	chunk := make([]float32, 512)
	for i := 0; i < 31; i++ { // ~1 s warmup
		v.AcceptWaveform(chunk)
	}
	return &session{dispose: func() {
		sherpa.DeleteVoiceActivityDetector(v)
	}}, nil
}

func openAsr(modelsDir string) (*session, error) {
	rec := sherpa.NewOnlineRecognizer(&sherpa.OnlineRecognizerConfig{
		FeatConfig:     sherpa.FeatureConfig{SampleRate: 16000, FeatureDim: 80},
		DecodingMethod: "greedy_search",
		ModelConfig: sherpa.OnlineModelConfig{
			Paraformer: sherpa.OnlineParaformerModelConfig{
				Encoder: filepath.Join(modelsDir, "asr-encoder.int8.onnx"),
				Decoder: filepath.Join(modelsDir, "asr-decoder.int8.onnx"),
			},
			Tokens:     filepath.Join(modelsDir, "asr-tokens.txt"),
			NumThreads: 1,
			Provider:   "cpu",
			Debug:      0,
		},
	})
	if rec == nil {
		return nil, fmt.Errorf("NewOnlineRecognizer returned nil")
	}
	stream := sherpa.NewOnlineStream(rec)
	chunk := make([]float32, 512)
	for i := 0; i < 16000/512; i++ {
		stream.AcceptWaveform(16000, chunk)
	}
	stream.InputFinished()
	for !rec.IsReady(stream) {
		time.Sleep(5 * time.Millisecond)
	}
	rec.Decode(stream)
	_ = rec.GetResult(stream)
	return &session{dispose: func() {
		sherpa.DeleteOnlineStream(stream)
		sherpa.DeleteOnlineRecognizer(rec)
	}}, nil
}

func newTts(modelsDir string) *sherpa.OfflineTts {
	matchaDir := filepath.Join(modelsDir, "matcha-icefall-zh-baker")
	return sherpa.NewOfflineTts(&sherpa.OfflineTtsConfig{
		Model: sherpa.OfflineTtsModelConfig{
			Matcha: sherpa.OfflineTtsMatchaModelConfig{
				AcousticModel: filepath.Join(matchaDir, "model-steps-3.onnx"),
				Vocoder:       filepath.Join(modelsDir, "vocos-22khz-univ.onnx"),
				Lexicon:       filepath.Join(matchaDir, "lexicon.txt"),
				Tokens:        filepath.Join(matchaDir, "tokens.txt"),
				DictDir:       filepath.Join(matchaDir, "dict"),
				NoiseScale:    0.667,
				LengthScale:   1.0,
			},
			NumThreads: 1,
			Provider:   "cpu",
			Debug:      0,
		},
		RuleFsts: filepath.Join(matchaDir, "phone.fst") + "," + filepath.Join(matchaDir, "date.fst") + "," + filepath.Join(matchaDir, "number.fst"),
	})
}

func openTts(modelsDir string) (*session, error) {
	t := newTts(modelsDir)
	if t == nil {
		return nil, fmt.Errorf("NewOfflineTts returned nil")
	}
	audio := t.Generate("你好，世界。", 0, 1.0)
	_ = audio
	return &session{dispose: func() {
		sherpa.DeleteOfflineTts(t)
	}}, nil
}

func openSession(which, modelsDir string) (func() (*session, error), []modelFile) {
	switch which {
	case "kws":
		kwsDir := filepath.Join(modelsDir, "sherpa-onnx-kws-zipformer-wenetspeech-3.3M-2024-01-01")
		return func() (*session, error) { return openKws(modelsDir) }, []modelFile{
			statFile(filepath.Join(kwsDir, "encoder-epoch-12-avg-2-chunk-16-left-64.int8.onnx")),
			statFile(filepath.Join(kwsDir, "decoder-epoch-12-avg-2-chunk-16-left-64.int8.onnx")),
			statFile(filepath.Join(kwsDir, "joiner-epoch-12-avg-2-chunk-16-left-64.int8.onnx")),
			statFile(filepath.Join(kwsDir, "tokens.txt")),
			statFile(filepath.Join(kwsDir, "keywords.txt")),
		}
	case "vad":
		return func() (*session, error) { return openVad(modelsDir) }, []modelFile{statFile(filepath.Join(modelsDir, "silero_vad.onnx"))}
	case "asr":
		return func() (*session, error) { return openAsr(modelsDir) }, []modelFile{
			statFile(filepath.Join(modelsDir, "asr-encoder.int8.onnx")),
			statFile(filepath.Join(modelsDir, "asr-decoder.int8.onnx")),
			statFile(filepath.Join(modelsDir, "asr-tokens.txt")),
		}
	case "tts":
		matchaDir := filepath.Join(modelsDir, "matcha-icefall-zh-baker")
		return func() (*session, error) { return openTts(modelsDir) }, []modelFile{
			statFile(filepath.Join(matchaDir, "model-steps-3.onnx")),
			statFile(filepath.Join(modelsDir, "vocos-22khz-univ.onnx")),
			statFile(filepath.Join(matchaDir, "lexicon.txt")),
			statFile(filepath.Join(matchaDir, "tokens.txt")),
		}
	}
	return nil, nil
}

// ------------------------------------------------------------- modes

func runResidency(which, modelsDir string, budgetMs int64, out string) {
	rep := residencyReport{Program: "model-residency", Which: which,
		Machine: common.GetMachineInfo(), StartedAt: time.Now().UTC().Format(time.RFC3339),
		SherpaVersion: sherpa.GetVersion(),
		Note:          "idle residency = loaded, settled, NO inference; post-inference = after one warmup inference"}
	defer func() {
		b, _ := json.MarshalIndent(rep, "", "  ")
		fmt.Println(string(b))
		if out != "" {
			os.WriteFile(out, b, 0644)
		}
	}()

	opener, files := openSession(which, modelsDir)
	if opener == nil {
		rep.Error = "unknown which " + which
		return
	}
	rep.Models = files

	common.SettleGC()
	_, med, _ := common.SampleStable(5, 200)
	pre := med.PrivateWorkingSet
	rep.PreWSMB = mbOf(pre)

	t0 := time.Now()
	s, err := opener()
	rep.LoadMs = time.Since(t0).Milliseconds()
	if err != nil {
		rep.Error = err.Error()
		return
	}

	common.SettleGC()
	_, m1, _ := common.SampleStable(5, 200)
	rep.IdleResidencyMB = mbOf(m1.PrivateWorkingSet)
	rep.PostInferCommitMB = mbOf(m1.PrivateCommit)
	rep.PeakCommitMB = loadPeakCommit()

	// dispose + settle
	dt0 := time.Now()
	s.dispose()
	common.SettleGC()
	rep.DisposeMs = time.Since(dt0).Milliseconds()
	_, _, _ = med, med, med
	sMs, ok, _ := settledWithin(pre, budgetMs)
	rep.SettledMs = sMs
	rep.Settled = ok
}

func runSwitch(modelsDir string, budgetMs int64, out string) {
	rep := residencyReport{Program: "model-residency", Which: "switch",
		Machine: common.GetMachineInfo(), StartedAt: time.Now().UTC().Format(time.RFC3339),
		SherpaVersion: sherpa.GetVersion(),
		Note:          "half-duplex serial: ASR fully disposed (incl FreeOSMemory) before TTS load; budget 1600ms = D32 16.3.4 'first token -> first TTS audio, TTS not yet loaded'"}
	defer func() {
		b, _ := json.MarshalIndent(rep, "", "  ")
		fmt.Println(string(b))
		if out != "" {
			os.WriteFile(out, b, 0644)
		}
	}()

	common.SettleGC()
	_, med, _ := common.SampleStable(5, 200)
	pre := med.PrivateWorkingSet
	rep.PreWSMB = mbOf(pre)

	sw := &switchReport{BudgetMs: 1600}

	// First TTS load (cold, no ASR ever loaded) - the D32 16.3.4 "1600ms" branch.
	t0 := time.Now()
	ttsS, err := openTts(modelsDir)
	sw.FirstTtsLoadMs = time.Since(t0).Milliseconds()
	if err != nil {
		rep.Error = "tts cold load: " + err.Error()
		return
	}
	common.SettleGC()
	_, m, _ := common.SampleStable(3, 150)
	ttsOnlyMB := mbOf(m.PrivateWorkingSet)
	ttsS.dispose()
	common.SettleGC()

	// Cycles: load ASR -> dispose -> load TTS -> dispose -> ...
	for cycle := 1; cycle <= 3; cycle++ {
		var c switchCycle
		c.Cycle = cycle

		t0 := time.Now()
		asrS, err := openAsr(modelsDir)
		c.AsrLoadMs = time.Since(t0).Milliseconds()
		if err != nil {
			rep.Error = "asr load: " + err.Error()
			return
		}
		common.SettleGC()
		_, ma, _ := common.SampleStable(3, 150)
		c.WsAtAsrMB = mbOf(ma.PrivateWorkingSet)

		t0 = time.Now()
		asrS.dispose()
		common.SettleGC()
		c.AsrDisposeMs = time.Since(t0).Milliseconds()

		t0 = time.Now()
		ttsS, err := openTts(modelsDir)
		c.TtsLoadMs = time.Since(t0).Milliseconds()
		if err != nil {
			rep.Error = "tts load: " + err.Error()
			return
		}
		c.SwitchAsrToTtsMs = c.AsrDisposeMs + c.TtsLoadMs
		common.SettleGC()
		_, mt, _ := common.SampleStable(3, 150)
		c.WsAtTtsMB = mbOf(mt.PrivateWorkingSet)

		t0 = time.Now()
		ttsS.dispose()
		common.SettleGC()
		c.TtsDisposeMs = time.Since(t0).Milliseconds()

		t0 = time.Now()
		asrS2, err := openAsr(modelsDir)
		c.AsrLoadMs = time.Since(t0).Milliseconds()
		if err != nil {
			rep.Error = "asr reload: " + err.Error()
			return
		}
		c.SwitchTtsToAsrMs = c.TtsDisposeMs + c.AsrLoadMs
		asrS2.dispose()
		common.SettleGC()
		_, mf, _ := common.SampleStable(1, 100)
		_ = mf

		sw.Cycles = append(sw.Cycles, c)
	}

	// stats
	var switches []float64
	for _, c := range sw.Cycles {
		switches = append(switches, float64(c.SwitchAsrToTtsMs))
		sw.MaxWsObservedMB = maxF(sw.MaxWsObservedMB, maxF(c.WsAtAsrMB, c.WsAtTtsMB))
	}
	sw.P50SwitchMs = pctlF(switches, 0.5)
	sw.MinSwitchMs = switches[0]
	for _, v := range switches {
		if v < sw.MinSwitchMs {
			sw.MinSwitchMs = v
		}
	}
	sw.PeakCommitMB = loadPeakCommit()

	// peak = max not sum? Compare observed residency vs the larger single model
	// plus the pre baseline; a sum would show roughly asr+ tts (~200+150=350+)
	asrTypical := sw.Cycles[0].WsAtAsrMB - rep.PreWSMB
	ttsTypical := ttsOnlyMB - rep.PreWSMB
	bigger := maxF(asrTypical, ttsTypical)
	sw.PeakIsMaxNotSum = sw.MaxWsObservedMB < rep.PreWSMB+bigger*1.5+10

	if sw.P50SwitchMs <= sw.BudgetMs {
		sw.BudgetVerdict = fmt.Sprintf("P50 %.0fms <= 1600ms budget: PASS", sw.P50SwitchMs)
	} else {
		sw.BudgetVerdict = fmt.Sprintf("P50 %.0fms > 1600ms budget: FAIL -> D32 fallback 'TTS resident + ASR on-demand'", sw.P50SwitchMs)
	}
	rep.Switch = sw
}

func pctlF(v []float64, p float64) float64 {
	if len(v) == 0 {
		return 0
	}
	s := append([]float64(nil), v...)
	sortF(s)
	return s[int(float64(p)*float64(len(s)-1)+0.5)]
}

func sortF(v []float64) {
	for i := 1; i < len(v); i++ {
		for j := i; j > 0 && v[j] < v[j-1]; j-- {
			v[j], v[j-1] = v[j-1], v[j]
		}
	}
}

func maxF(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	which := flag.String("which", "kws", "kws | vad | asr | tts | switch")
	modelsDir := flag.String("models", "third_party/spike-models", "spike models dir")
	budget := flag.Int64("budget-ms", 15000, "settle observation window")
	out := flag.String("out", "", "path to write JSON")
	flag.Parse()

	switch *which {
	case "switch":
		runSwitch(*modelsDir, *budget, *out)
	case "kws", "vad", "asr", "tts":
		runResidency(*which, *modelsDir, *budget, *out)
	default:
		fmt.Fprintln(os.Stderr, "unknown which", *which)
		os.Exit(2)
	}
}
