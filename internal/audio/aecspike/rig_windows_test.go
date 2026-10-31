//go:build windows

package aecspike

// The gated manual driver. Nothing in this file runs in CI: every
// measurement test skips unless WISP_AEC_SPIKE=1, and writes artifacts ONLY
// into the workspace temp root (WISP_AEC_SPIKE_OUT, default the root the
// dispatcher named). The ungated tests at the bottom are the CI-safe smoke
// that the package still enumerates and probes without failing anywhere.
//
// Run shape (proven PATH prefix so the sherpa runtime DLLs resolve):
//
//	PATH="$PWD/third_party/sherpa-onnx:$PATH" WISP_AEC_SPIKE=1 \
//	  go test ./internal/audio/aecspike/ -run 'TestSpike' -v

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/speech"
)

// spikeVolumes is the volume axis of the false-trigger table: quiet room,
// normal, loud.
var spikeVolumes = []float64{0.2, 0.5, 0.8}

// spikeDriftCurve is the synthetic cross-endpoint drift axis (ppm).
var spikeDriftCurve = []int{0, 25, 50, 100, 200, 400}

// procGetProcessTimes is the CPU-accounting syscall of the CPU cell.
var procGetProcessTimes = windows.NewLazySystemDLL("kernel32.dll").NewProc("GetProcessTimes")

func spikeGate(t *testing.T) {
	t.Helper()
	if os.Getenv("WISP_AEC_SPIKE") != "1" {
		t.Skip("hardware rig: set WISP_AEC_SPIKE=1 (renders synthetic audio on the default endpoint)")
	}
}

func spikeOutDir(t *testing.T) string {
	t.Helper()
	dir := os.Getenv("WISP_AEC_SPIKE_OUT")
	if dir == "" {
		dir = `D:\tmp\wisp\20261028\59r1`
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("spike out dir: %v", err)
	}
	return dir
}

// spikeRepoRoot walks up from the package dir to the module root (same
// idiom as internal/speech tests; no filepath.Abs anywhere: the D22 ban on
// path decisions outside risk.PathResolver stands in tests too). The staged
// test binary (DLL-colocated run shape) sets WISP_AEC_SPIKE_ROOT because
// its cwd is the temp root, far above the module.
func spikeRepoRoot(t *testing.T) string {
	t.Helper()
	if root := os.Getenv("WISP_AEC_SPIKE_ROOT"); root != "" {
		if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
			return root
		}
		t.Fatalf("WISP_AEC_SPIKE_ROOT=%s has no go.mod", root)
	}
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatalf("module root not found above %s", dir)
	return ""
}

func spikeModelsDir(t *testing.T) string {
	t.Helper()
	dir := filepath.Join(spikeRepoRoot(t), "third_party", "spike-models")
	if _, err := os.Stat(filepath.Join(dir, "silero_vad.onnx")); err != nil {
		t.Skipf("spike models not staged on this host (%v)", err)
	}
	return dir
}

func spikeVadStarts(t *testing.T, det *speech.Detector, samples []float32) int {
	t.Helper()
	var eps []speech.Endpoint
	eps = det.Accept(samples, eps)
	eps = det.Finalize(eps)
	starts := 0
	for _, e := range eps {
		if e.Kind == speech.EndpointStart {
			starts++
		}
	}
	return starts
}

// spikeNewVad loads one silero detector (once per caller, OUTSIDE any CPU
// accounting window: model load must not pollute the sustained-CPU cell).
func spikeNewVad(t *testing.T) *speech.Detector {
	t.Helper()
	det, err := speech.NewDetector(speech.DetectorConfig{
		ModelPath: filepath.Join(spikeModelsDir(t), "silero_vad.onnx"),
	})
	if err != nil {
		t.Skipf("silero vad unavailable: %v", err)
	}
	t.Cleanup(func() { _ = det.Close() })
	return det
}

func spikeTranscribe(t *testing.T, asr *speech.StreamingAsr, samples []float32) string {
	t.Helper()
	asr.Begin()
	asr.Feed(samples)
	return asr.Finish()
}

func spikeTTSFixture(t *testing.T) []int16 {
	t.Helper()
	dir := spikeModelsDir(t)
	m := filepath.Join(dir, "matcha-icefall-zh-baker")
	cfg := speech.TtsConfig{
		AcousticModel: filepath.Join(m, "model-steps-3.onnx"),
		Vocoder:       filepath.Join(dir, "vocos-22khz-univ.onnx"),
		Lexicon:       filepath.Join(m, "lexicon.txt"),
		Tokens:        filepath.Join(m, "tokens.txt"),
		DictDir:       filepath.Join(m, "dict"),
		RuleFsts: []string{
			filepath.Join(m, "phone.fst"),
			filepath.Join(m, "date.fst"),
			filepath.Join(m, "number.fst"),
		},
	}
	tts, err := speech.NewMatchaTts(cfg)
	if err != nil {
		t.Skipf("matcha tts unavailable on this host: %v", err)
	}
	defer func() { _ = tts.Close() }()
	var all []float32
	if err := tts.Speak("你好，我是威斯普，今天天气不错。", func(samples []float32) error {
		all = append(all, samples...)
		return nil
	}); err != nil {
		t.Fatalf("tts speak: %v", err)
	}
	return audio.NewResampler(tts.SampleRate(), audio.TargetRate).Process(FloatToPCM16(all))
}

func spikeASR(t *testing.T) *speech.StreamingAsr {
	t.Helper()
	dir := spikeModelsDir(t)
	asr, err := speech.NewStreamingAsr(speech.AsrConfig{
		EncoderPath: filepath.Join(dir, "asr-encoder.int8.onnx"),
		DecoderPath: filepath.Join(dir, "asr-decoder.int8.onnx"),
		TokensPath:  filepath.Join(dir, "asr-tokens.txt"),
	})
	if err != nil {
		t.Skipf("streaming asr unavailable on this host: %v", err)
	}
	t.Cleanup(func() { _ = asr.Close() })
	return asr
}

func writeJSON(t *testing.T, path string, v any) {
	t.Helper()
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatalf("marshal %s: %v", path, err)
	}
	if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func meanCorr(vs []float64) float64 {
	if len(vs) == 0 {
		return 0
	}
	var sum float64
	for _, v := range vs {
		sum += v
	}
	return sum / float64(len(vs))
}

// TestSpikeEnumerateIsSafeOnHeadless runs everywhere on windows CI: the
// enumeration must answer (possibly with zero endpoints) without failing.
func TestSpikeEnumerateIsSafeOnHeadless(t *testing.T) {
	inv, err := EnumerateEndpoints()
	if err != nil {
		t.Fatalf("EnumerateEndpoints: %v", err)
	}
	t.Logf("endpoints=%d default_render=%q default_capture=%q",
		len(inv.Endpoints), inv.DefaultRender, inv.DefaultCapture)
}

// TestSpikeEndpointsAndDmoInventory is AC#1/AC#2 evidence: what the OS
// exposes on THIS machine, including the voice-AEC DMO presence probe
// (build+link+run of the cgo shim on the pinned toolchain) and the DMO's
// registered DLL path (the "exact source/build" of the license verdict).
func TestSpikeEndpointsAndDmoInventory(t *testing.T) {
	spikeGate(t)
	inv, err := EnumerateEndpoints()
	if err != nil {
		t.Fatalf("EnumerateEndpoints: %v", err)
	}
	dmo := DmoProbe()
	report := struct {
		OS             string         `json:"os"`
		Endpoints      []EndpointInfo `json:"endpoints"`
		DefaultRender  string         `json:"default_render"`
		DefaultCapture string         `json:"default_capture"`
		Dmo            DmoInfo        `json:"dmo"`
		DmoDLL         string         `json:"dmo_dll"`
	}{
		Endpoints: inv.Endpoints, DefaultRender: inv.DefaultRender,
		DefaultCapture: inv.DefaultCapture, Dmo: dmo,
	}
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err == nil {
		defer k.Close()
		pn, _, _ := k.GetStringValue("ProductName")
		cb, _, _ := k.GetStringValue("CurrentBuildNumber")
		dv, _, _ := k.GetStringValue("DisplayVersion")
		report.OS = pn + " build " + cb + " (" + dv + ")"
	}
	if ck, err := registry.OpenKey(registry.CLASSES_ROOT,
		`CLSID\{745057C7-F353-4F2D-A7EE-5843427A0EF4}\InprocServer32`, registry.QUERY_VALUE); err == nil {
		defer ck.Close()
		report.DmoDLL, _, _ = ck.GetStringValue("")
	}
	writeJSON(t, filepath.Join(spikeOutDir(t), "inventory.json"), report)
	t.Logf("os=%s endpoints=%d dmo={present=%v hr=%s in=%d out=%d} dll=%q",
		report.OS, len(inv.Endpoints), dmo.Present, dmo.HRESULT, dmo.InStreams, dmo.OutStreams, report.DmoDLL)
}

// TestSpikeSelfEchoPathRig is AC#3/AC#4's same-endpoint reading: render the
// TTS fixture, capture the endpoint loopback, align, subtract, and measure
// path ERLE, drift ratio, the VAD/ASR suppression readings and the drift
// sensitivity curve.
func TestSpikeSelfEchoPathRig(t *testing.T) {
	spikeGate(t)
	out := spikeOutDir(t)
	reg := observe.NewRegistry()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	far16 := spikeTTSFixture(t)
	far := PCM16ToFloat(far16)
	reading, capturedMix, err := RunLoopback(reg, ctx, "", 0.5, far, audio.TargetRate)
	if err != nil {
		t.Skipf("loopback rig unavailable on this host: %v", err)
	}
	if reading == nil || reading.MixRate <= 0 {
		t.Fatalf("rig returned reading=%+v captured=%d samples; cannot continue",
			reading, len(capturedMix))
	}
	near16 := audio.NewResampler(reading.MixRate, audio.TargetRate).Process(capturedMix)
	near := PCM16ToFloat(near16)

	drift := map[string]any{}
	for _, ppm := range spikeDriftCurve {
		dFar := ApplyPpmDrift(far, ppm)
		l, _ := AlignByEnvelopeCorrelation(dFar, near, DefaultMaxLag)
		res, _, _ := SubtractEcho(dFar, near, l)
		drift[fmt.Sprint(ppm)] = erleJSON(ERLEdB(dFar, near, res, l))
	}

	asr := spikeASR(t)
	lag, strength := AlignByEnvelopeCorrelation(far, near, DefaultMaxLag)
	// CPU cell (AC#6): the sustained per-audio-second cost of the
	// reference-subtraction stand-in plus the product's VAD, on this host.
	// The OS AEC engine itself is not in this process (the DMO probe answers
	// that), so this is the in-process floor, not the whole Path C row.
	subCpu0, subWall0 := spikeCPUTimes(t), time.Now()
	residual, gain, _ := SubtractEcho(far, near, lag)
	corr := WindowedCorrelation(far, near, lag, 1024, 512)
	subCpu, subWall := spikeCPUTimes(t).sub(subCpu0), time.Since(subWall0)
	pathERLE := ERLEdB(far, near, residual, lag)
	vad := spikeNewVad(t)
	vadCpu0, vadWall0 := spikeCPUTimes(t), time.Now()
	vadRaw := spikeVadStarts(t, vad, near)
	vadRes := spikeVadStarts(t, vad, residual)
	vadCpu, vadWall := spikeCPUTimes(t).sub(vadCpu0), time.Since(vadWall0)
	audioSec := float64(len(near)) / audio.TargetRate

	report := struct {
		Reading            *RigReading    `json:"reading"`
		AlignLag           int            `json:"align_lag_samples"`
		AlignStrength      float64        `json:"align_strength"`
		PathERLE           any            `json:"path_erle_db"`
		SubtractGain       float64        `json:"subtract_gain"`
		CorrMean           float64        `json:"corr_mean"`
		VadRaw             int            `json:"vad_starts_raw"`
		VadResidual        int            `json:"vad_starts_residual"`
		AsrFar             string         `json:"asr_far_reference"`
		AsrRaw             string         `json:"asr_raw_playback"`
		AsrResidual        string         `json:"asr_residual"`
		DriftCurveERLE     map[string]any `json:"drift_curve_erle_db"`
		AudioSeconds       float64        `json:"audio_seconds"`
		SubtractCpuPercent float64        `json:"subtract_cpu_percent_of_core"`
		VadCpuPercent      float64        `json:"vad_both_passes_cpu_percent_of_core"`
	}{
		reading, lag, strength, erleJSON(pathERLE), gain, meanCorr(corr),
		vadRaw, vadRes,
		spikeTranscribe(t, asr, far), spikeTranscribe(t, asr, near),
		spikeTranscribe(t, asr, residual), drift,
		audioSec, cpuPercent(subCpu, subWall), cpuPercent(vadCpu, vadWall),
	}
	writeJSON(t, filepath.Join(out, "selfecho.json"), report)
	_ = WriteWav16(filepath.Join(out, "selfecho-far-16k.wav"), audio.TargetRate, far16)
	_ = WriteWav16(filepath.Join(out, "selfecho-near-16k.wav"), audio.TargetRate, near16)
	_ = WriteWav16(filepath.Join(out, "selfecho-residual-16k.wav"), audio.TargetRate, FloatToPCM16(residual))
	t.Logf("reading=%+v lag=%d strength=%.3f pathERLE=%.1fdB corr=%.3f vadRaw=%d vadRes=%d",
		reading, lag, strength, pathERLE, report.CorrMean, report.VadRaw, report.VadResidual)
	t.Logf("cpu: subtract=%.2f%% vad(2 passes)=%.2f%% of one core over %.1fs audio",
		report.SubtractCpuPercent, report.VadCpuPercent, audioSec)
	t.Logf("asr far=%q raw=%q residual=%q", report.AsrFar, report.AsrRaw, report.AsrResidual)
}

// falseTriggerRow is one music/speech/video x volume cell of AC#5's table.
type falseTriggerRow struct {
	Fixture          string  `json:"fixture"`
	Volume           float64 `json:"volume"`
	DurationSec      float64 `json:"duration_sec"`
	VadStartsRaw     int     `json:"vad_starts_raw"`
	RawPer10Min      float64 `json:"raw_per_10min"`
	VadStartsResid   int     `json:"vad_starts_residual"`
	ResidualPer10Min float64 `json:"residual_per_10min"`
	// CorrMean is the near-vs-TTS-reference correlation AT THE ALIGNED LAG:
	// self echo scores ~1, third-party sources ~0. It is the recommended
	// barge-in gate feature (AC#5's "direction of AEC residual").
	CorrMean      float64 `json:"corr_mean_tts_ref_aligned"`
	AlignLag      int     `json:"align_lag_samples"`
	AlignStrength float64 `json:"align_strength"`
	RatioPpm      int     `json:"ratio_ppm"`
}

// TestSpikeFalseTriggerTable is AC#5: the external-audio table. The
// fixtures are synthetic (music/video) plus the real TTS voice for speech;
// the capture is the endpoint loopback (digital, no room reverberation),
// which makes every rate here a LOWER bound of the physical-room case.
func TestSpikeFalseTriggerTable(t *testing.T) {
	spikeGate(t)
	out := spikeOutDir(t)
	reg := observe.NewRegistry()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tts16 := spikeTTSFixture(t)
	ttsRef := PCM16ToFloat(tts16)
	music := SynthMusic(audio.TargetRate, 8)
	video := SynthVideo(audio.TargetRate, 8, ttsRef)
	fixtures := []struct {
		name string
		sig  []float32
	}{{"speech-tts", ttsRef}, {"music-synth", music}, {"video-synth", video}}
	vad := spikeNewVad(t)

	var rows []falseTriggerRow
	for _, fx := range fixtures {
		for _, vol := range spikeVolumes {
			reading, capturedMix, err := RunLoopback(reg, ctx, "", vol, fx.sig, audio.TargetRate)
			if err != nil {
				t.Skipf("loopback rig unavailable: %v", err)
			}
			near16 := audio.NewResampler(reading.MixRate, audio.TargetRate).Process(capturedMix)
			near := PCM16ToFloat(near16)
			// The external-audio rows subtract the TTS reference (which does
			// NOT contain the fixture): AEC alone cannot suppress a
			// third-party source, and that is the finding the table records.
			lag, strength := AlignByEnvelopeCorrelation(ttsRef, near, DefaultMaxLag)
			residual, _, _ := SubtractEcho(ttsRef, near, lag)
			corr := WindowedCorrelation(ttsRef, near, lag, 1024, 512)
			dur := float64(len(near)) / audio.TargetRate
			raw := spikeVadStarts(t, vad, near)
			res := spikeVadStarts(t, vad, residual)
			row := falseTriggerRow{
				Fixture: fx.name, Volume: vol, DurationSec: dur,
				VadStartsRaw: raw, RawPer10Min: per10Min(raw, dur),
				VadStartsResid: res, ResidualPer10Min: per10Min(res, dur),
				CorrMean: meanCorr(corr), AlignLag: lag, AlignStrength: strength,
				RatioPpm: reading.RatioPpm,
			}
			rows = append(rows, row)
			t.Logf("%-12s vol=%.1f raw=%d (%.1f/10min) resid=%d (%.1f/10min) corr=%.3f align=%d@%.3f",
				fx.name, vol, raw, row.RawPer10Min, res, row.ResidualPer10Min,
				row.CorrMean, lag, strength)
		}
	}
	writeJSON(t, filepath.Join(out, "false-trigger.json"), rows)
}

func per10Min(starts int, seconds float64) float64 {
	if seconds <= 0 {
		return 0
	}
	return float64(starts) * 600 / seconds
}

// TestSpikeDmoFullDuplexAec is the candidate-A reading: drive the OS Voice
// Capture DSP full duplex with the far reference and the loopback echo and
// measure its residual, plus the sustained processing CPU. Absence of the
// DMO is itself the recorded reading (the shim's hr names it).
func TestSpikeDmoFullDuplexAec(t *testing.T) {
	spikeGate(t)
	out := spikeOutDir(t)
	dmo := DmoProbe()
	if !dmo.Present {
		t.Logf("DMO probe: not present (hr=%s) - recorded as the candidate-A reading", dmo.HRESULT)
		writeJSON(t, filepath.Join(out, "dmo.json"), map[string]any{"dmo": dmo})
		t.Skipf("voice capture DMO not registered on this host (hr=%s)", dmo.HRESULT)
	}
	reg := observe.NewRegistry()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	far16 := spikeTTSFixture(t)
	far := PCM16ToFloat(far16)
	reading, capturedMix, err := RunLoopback(reg, ctx, "", 0.5, far, audio.TargetRate)
	if err != nil {
		t.Skipf("loopback rig unavailable: %v", err)
	}
	near16 := audio.NewResampler(reading.MixRate, audio.TargetRate).Process(capturedMix)
	near := PCM16ToFloat(near16)
	lag, _ := AlignByEnvelopeCorrelation(far, near, DefaultMaxLag)
	if lag < 0 {
		lag = 0 // the DMO feed needs a non-negative pre-alignment offset
	}
	// Rotate the near signal so nearAligned[i] ~ far[i]: the DMO wants
	// time-aligned streams and does its own adaptive alignment on top.
	nearAligned := near
	if lag > 0 && lag < len(near) {
		nearAligned = near[lag:]
		far = far[:len(near)-lag]
		far16 = far16[:len(near)-lag]
	}

	sess, err := DmoBegin()
	if err != nil {
		writeJSON(t, filepath.Join(out, "dmo.json"), map[string]any{
			"dmo": dmo, "begin_error": err.Error(),
		})
		t.Fatalf("DmoBegin: %v", err)
	}
	defer sess.Close()

	const chunk = 1600
	processed := make([]float32, 0, len(nearAligned))
	var produced int
	cpu0 := spikeCPUTimes(t)
	wall0 := time.Now()
	for i := 0; i+chunk <= len(far16); i += chunk {
		outBuf := make([]int16, chunk)
		n, err := sess.Process(far16[i:i+chunk], near16[lag+i:lag+i+chunk], outBuf)
		if err != nil {
			writeJSON(t, filepath.Join(out, "dmo.json"), map[string]any{
				"dmo": dmo, "process_error": err.Error(), "chunks_ok": i / chunk,
			})
			t.Fatalf("DMO process chunk %d: %v", i/chunk, err)
		}
		processed = append(processed, PCM16ToFloat(outBuf[:n])...)
		produced += n
	}
	// Drain: the engine lags its input by its internal buffering; feed
	// silence until it stops producing (bounded, not a timeout).
	for drained := 0; drained < 8; drained++ {
		silence := make([]int16, chunk)
		outBuf := make([]int16, chunk)
		n, err := sess.Process(silence, silence, outBuf)
		if err != nil || n == 0 {
			break
		}
		processed = append(processed, PCM16ToFloat(outBuf[:n])...)
		produced += n
	}
	cpu := spikeCPUTimes(t).sub(cpu0)
	wall := time.Since(wall0)
	audioSec := float64(len(far16)) / audio.TargetRate
	erle := ERLEdB(far, nearAligned, processed, 0) // streams are pre-aligned
	report := map[string]any{
		"dmo":           dmo,
		"chunks":        len(far16) / chunk,
		"produced":      produced,
		"audio_seconds": audioSec,
		"cpu_percent":   cpuPercent(cpu, wall),
		"cpu_per_audio": (float64(cpu.total()) / 1e9) / audioSec,
		"erle_db":       erleJSON(erle),
		"align_lag":     lag,
		"note":          "digital loopback echo path (no room reverberation); near pre-aligned, engine does its own adaptive alignment",
	}
	writeJSON(t, filepath.Join(out, "dmo.json"), report)
	t.Logf("DMO AEC: chunks=%d produced=%d ERLE=%.1fdB cpu=%.2f%% of one core (%.3fs CPU per audio-second)",
		len(far16)/chunk, produced, erle, cpuPercent(cpu, wall), float64(cpu.total())/1e9/audioSec)
}

// erleJSON renders an ERLE for JSON: perfect cancellation is +Inf, which
// encoding/json refuses, so it names the limit instead.
func erleJSON(f float64) any {
	if math.IsInf(f, 1) {
		return "+Infinity (perfect cancellation over the active frames)"
	}
	if math.IsInf(f, -1) {
		return "-Infinity (no near energy over the active frames)"
	}
	return f
}

// spikeCPUTimes reads this process's kernel+user CPU time (measurement
// accounting for the CPU cell; never a timeout mechanism).
func spikeCPUTimes(t *testing.T) cpuTimes {
	var creation, exit, kernel, user windows.Filetime
	r, _, _ := procGetProcessTimes.Call(uintptr(windows.CurrentProcess()),
		uintptr(unsafe.Pointer(&creation)), uintptr(unsafe.Pointer(&exit)),
		uintptr(unsafe.Pointer(&kernel)), uintptr(unsafe.Pointer(&user)))
	if r == 0 {
		t.Fatalf("GetProcessTimes failed")
	}
	return cpuTimes{user: time.Duration(user.Nanoseconds()), kernel: time.Duration(kernel.Nanoseconds())}
}

type cpuTimes struct{ user, kernel time.Duration }

func (c cpuTimes) total() time.Duration { return c.user + c.kernel }

func (c cpuTimes) sub(o cpuTimes) cpuTimes {
	return cpuTimes{c.user - o.user, c.kernel - o.kernel}
}

// cpuPercent renders the CPU delta as a percentage of ONE core over the
// wall span (the SLO table's framing: "single-core 5%").
func cpuPercent(delta cpuTimes, wall time.Duration) float64 {
	if wall <= 0 {
		return 0
	}
	return 100 * float64(delta.total()) / float64(wall)
}
