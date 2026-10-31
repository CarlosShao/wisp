//go:build windows

package aecspike

// The 59-w-r1 driver: the PRECHECK P15 re-open readings on the aec3-r1 rig
// extended to long windows. Nothing here runs in CI (the WISP_AEC_SPIKE=1
// gate applies); the native binding must be part of the test build
// (-tags webrtcap per docs/evidence/s1/aec3r1-remeasure.md section 7) or the
// test skips with that fact recorded.
//
// Stages, in one pass so the artifacts stay mutually consistent:
//
//   A. long-window runs (60s/120s, TTS far x3 + music far x1): the far
//      loops, the loopback capture streams to disk in segments, one cold
//      AEC3 engine per run processes the whole run in 10ms lockstep, and the
//      per-window ERLE curve plus the collapse detector answer PASS
//      condition (a) - does the 8.9dB cold-start collapse vanish in a long
//      window, or only in the startup segment?
//   B. cold 3.4s runs on three head variants (immediate TTS onset, 0.5s
//      leading silence, music): the collapse-characterization matrix. Each
//      run records the full-window ERLE (the aec3-r1 ruler) and the
//      sub-sample render/loopback phase (FractionalLag) - the double peak's
//      leading suspect (docs/evidence/s1/59-v2-verdict.md section 4).
//   C. the combined barge-in gate (VAD / VAD+energy / VAD+energy+ASR text)
//      on the AEC output of good-peak and collapse-state material, plus a
//      speech positive control - PASS condition (b).
//
// Run shape (the aec3-r1 recipe, this leg's temp root):
//
//	set CGO_CXXFLAGS=-I<D:\tmp\wisp\20261031\aec3r1\webrtc-audio-processing-2.1\webrtc> -I<D:\tmp\wisp\20261031\aec3r1\abseil-cpp-20240722.0>
//	set CGO_LDFLAGS=-LD:\tmp\wisp\20261031\aec3r1\build
//	go test -c -o D:\tmp\wisp\20261031\59wr1\run\aecspike.test.exe -tags webrtcap ./internal/audio/aecspike/
//	cd D:\tmp\wisp\20261031\59wr1\run
//	set WISP_AEC_SPIKE=1 WISP_AEC_SPIKE_ROOT=<repo> WISP_AEC_SPIKE_OUT=D:\tmp\wisp\20261031\59wr1\out
//	set PATH=<repo>\third_party\sherpa-onnx;%PATH%
//	aecspike.test.exe -test.run TestSpikeAec3LongWindow -test.v -test.timeout 40m

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/audio/aec"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/speech"
)

// longWindowOutDefault is THIS leg's dispatcher-named temp root (never the
// aec3-r1 or 59-r1 one; artifacts are not mixed between legs).
const longWindowOutDefault = `D:\tmp\wisp\20261031\59wr1\out`

// archivedCollapseWav is the 59-v2 verdict leg's archived AEC output of a
// collapsed live capture (8.91dB per its JSON, ASR two syllables): the
// collapse-state gate material of last resort when this session's own cold
// runs all land on the good peak.
const archivedCollapseWav = `D:\tmp\wisp\20261031\59v2\out\aec3-fresh-rig-out-16k.wav`

func longWindowOutDir(t *testing.T) string {
	t.Helper()
	dir := os.Getenv("WISP_AEC_SPIKE_OUT")
	if dir == "" {
		dir = longWindowOutDefault
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("long-window out dir: %v", err)
	}
	return dir
}

// curveJSON renders one ERLE for JSON (encoding/json refuses NaN/Inf).
func curveJSON(v float64) any {
	switch {
	case math.IsNaN(v):
		return "inactive"
	case math.IsInf(v, 1):
		return "+Infinity"
	case math.IsInf(v, -1):
		return "-Infinity"
	default:
		return v
	}
}

// rmsDbfs is the RMS of a float32 span in dBFS (0 dBFS = int16 full scale).
func rmsDbfs(samples []float32) float64 {
	var sum float64
	for _, v := range samples {
		sum += float64(v) * float64(v)
	}
	if len(samples) == 0 || sum <= 0 {
		return -120
	}
	return 20 * math.Log10(math.Sqrt(sum/float64(len(samples))))
}

// longFixtureSet holds this leg's far fixtures at the seam rate.
type longFixtureSet struct {
	tts    []float32 // the pinned TTS phrase (the aec3-r1 far)
	ttsGap []float32 // 0.5s leading silence + TTS: the head-onset variant
	music  []float32 // deterministic synth music, 8s seamless loop
}

func longFixtureSetBuild(t *testing.T) longFixtureSet {
	t.Helper()
	tts := PCM16ToFloat(spikeTTSFixture(t))
	gap := make([]float32, audio.TargetRate/2)
	return longFixtureSet{
		tts:    tts,
		ttsGap: append(gap, tts...),
		music:  SynthMusic(audio.TargetRate, 8),
	}
}

// ttsSpeak is spikeTTSFixture with free text: the positive control needs a
// second phrase the far fixture never contained.
func ttsSpeak(t *testing.T, text string) []float32 {
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
	if err := tts.Speak(text, func(samples []float32) error {
		all = append(all, samples...)
		return nil
	}); err != nil {
		t.Fatalf("tts speak %q: %v", text, err)
	}
	return PCM16ToFloat(audio.NewResampler(tts.SampleRate(), audio.TargetRate).Process(FloatToPCM16(all)))
}

// windowPoint is one 1s point of the long-window ERLE curve.
type windowPoint struct {
	Window int `json:"window"`
	ERLE   any `json:"erle_db"`
	Active int `json:"active_frames"`
}

// jsonAlarm is CollapseAlarm with JSON-safe numbers (a NaN baseline is a
// legal value - it means the alarm fired on the absolute floor, not a drop).
func jsonAlarm(a CollapseAlarm) map[string]any {
	base := any("none")
	if !math.IsNaN(a.Baseline) {
		base = a.Baseline
	}
	return map[string]any{
		"window": a.WindowIndex, "reason": a.Reason,
		"window_erle_db": curveJSON(a.WindowERLE), "baseline_db": base,
	}
}

// runReport is one rig run's full reading.
type runReport struct {
	Name              string           `json:"name"`
	Kind              string           `json:"kind"` // "long" | "cold"
	FarContent        string           `json:"far_content"`
	Seconds           float64          `json:"seconds"`
	Reading           *LongRunResult   `json:"reading"`
	AlignLag          int              `json:"align_lag_samples"`
	AlignStrength     float64          `json:"align_strength"`
	FractionalLag     float64          `json:"fractional_lag_samples"`
	DelayMs           int              `json:"aec3_delay_ms"`
	HeadERLE          any              `json:"head_erle_db"`    // [0, 3.4s): the aec3-r1 cold-window ruler
	StartupERLE       any              `json:"startup_erle_db"` // [0, 5s)
	SteadyERLE        any              `json:"steady_erle_db"`  // last 30s
	FullERLE          any              `json:"full_erle_db"`
	ConvergenceWindow int              `json:"convergence_window"` // first window at steady median - 3dB; -1 none
	Alarms            []map[string]any `json:"collapse_alarms"`
	Aec3CpuPercent    float64          `json:"aec3_cpu_percent_of_core"`
	Aec3CpuPerAudio   float64          `json:"aec3_cpu_seconds_per_audio_second"`
	Outcome           string           `json:"outcome"` // "good" | "collapsed" (cold runs; 15dB cut)
	Windows           []windowPoint    `json:"windows"`
	AecSegments       []LongRunSegment `json:"aec_segments"`
}

// passOutcome carries what the offline AEC streaming pass produced.
type passOutcome struct {
	acc      *ERLEAccum
	alarms   []CollapseAlarm
	headFar  []float32
	headNear []float32
	headOut  []float32
	segs     []LongRunSegment
	cpuPct   float64
	cpuPA    float64
}

// headWindowSamples is 3.4s at the seam rate: the aec3-r1 cold-window ruler.
const headWindowSamples = 16000 * 34 / 10

// longAecPass is the offline streaming pass over one long-window capture:
// one cold AEC3 engine per run (the product's cold start), the capture
// resampled segment by segment (the streaming resampler is chunk-invariant),
// the far reference from the phase-continuous FarClock, per-window ERLE
// accumulation, the collapse detector on each closed window, and the AEC
// output itself written to disk in segments for the gate stage. Memory stays
// at one segment plus the 3.4s head, whatever the run length.
func longAecPass(t *testing.T, res *LongRunResult, lag, delayMs int, outDir, name string, segSeconds float64) *passOutcome {
	t.Helper()
	frame := aec.FrameSamples(16000)
	apm, err := aec.NewApm(16000)
	if err != nil {
		if err == aec.ErrUnavailable {
			t.Skipf("AEC3 binding not in this build: %v", err)
		}
		t.Fatalf("NewApm: %v", err)
	}
	defer func() { _ = apm.Close() }()

	rs := audio.NewResampler(res.MixRate, audio.TargetRate)
	clock := NewFarClock(res.MixFar, res.MixRate)
	acc := NewERLEAccum(frameRMS(res.MixFar))
	det := NewCollapseDetector()

	// The lag-shifted far stream (the echo-source-time far that decides
	// frame activity, the baseline ERLEdB convention): a ring of lag samples
	// delays the render clock into the activity input. lag==0 shortcuts.
	actBuf := make([]float32, frame)
	delay := make([]float32, lag)
	delayPos := 0

	farScratch := make([]float32, frame)
	var frameBuf []int16
	var outSeg []int16
	segFrames16 := int(segSeconds * audio.TargetRate)
	if segFrames16 <= 0 {
		segFrames16 = math.MaxInt32
	}
	windowsSeen := 0
	out := &passOutcome{acc: acc}

	cpu0, wall0 := spikeCPUTimes(t), time.Now()
	for _, seg := range res.Segments {
		rate, pcm, err := ReadWav16(seg.File)
		if err != nil {
			t.Fatalf("read capture segment %s: %v", seg.File, err)
		}
		if rate != res.MixRate {
			t.Fatalf("capture segment %s rate %d, want %d", seg.File, rate, res.MixRate)
		}
		frameBuf = append(frameBuf, rs.Process(pcm)...)
		for len(frameBuf) >= frame {
			near := PCM16ToFloat(frameBuf[:frame])
			if lag == 0 {
				for k := 0; k < frame; k++ {
					v := clock.Next()
					farScratch[k] = v
					actBuf[k] = v
				}
			} else {
				for k := 0; k < frame; k++ {
					v := clock.Next()
					farScratch[k] = v
					actBuf[k] = delay[delayPos]
					delay[delayPos] = v
					delayPos = (delayPos + 1) % lag
				}
			}
			if err := apm.ProcessRender(farScratch); err != nil {
				t.Fatalf("AEC3 render frame at window %d: %v", windowsSeen, err)
			}
			aecOut, err := apm.ProcessCapture(near, delayMs)
			if err != nil {
				t.Fatalf("AEC3 capture frame at window %d: %v", windowsSeen, err)
			}
			acc.Observe(actBuf, near, aecOut)
			for len(acc.Windows) > windowsSeen {
				w := acc.Windows[windowsSeen]
				windowsSeen++
				out.alarms = append(out.alarms, det.Observe(w)...)
			}
			if len(out.headNear) < headWindowSamples {
				room := headWindowSamples - len(out.headNear)
				take := frame
				if take > room {
					take = room
				}
				out.headFar = append(out.headFar, farScratch[:take]...)
				out.headNear = append(out.headNear, near[:take]...)
				out.headOut = append(out.headOut, aecOut[:take]...)
			}
			outSeg = append(outSeg, FloatToPCM16(aecOut)...)
			for len(outSeg) >= segFrames16 {
				out.segs = append(out.segs, writeAecSegment(t, outDir, name, out.segs, outSeg[:segFrames16]))
				rest := copy(outSeg, outSeg[segFrames16:])
				outSeg = outSeg[:rest]
			}
			frameBuf = frameBuf[:copy(frameBuf, frameBuf[frame:])]
		}
	}
	if len(outSeg) > 0 {
		out.segs = append(out.segs, writeAecSegment(t, outDir, name, out.segs, outSeg))
	}
	cpu := spikeCPUTimes(t).sub(cpu0)
	wall := time.Since(wall0)
	audioSec := float64(res.CapturedFrames) / float64(res.MixRate)
	out.cpuPct = cpuPercent(cpu, wall)
	out.cpuPA = (float64(cpu.total()) / 1e9) / audioSec
	return out
}

// writeAecSegment flushes one 16k AEC-output segment to disk.
func writeAecSegment(t *testing.T, outDir, name string, prev []LongRunSegment, samples []int16) LongRunSegment {
	t.Helper()
	idx := len(prev)
	start := 0
	if idx > 0 {
		start = prev[idx-1].StartFrame + prev[idx-1].Frames
	}
	file := filepath.Join(outDir, fmt.Sprintf("%s-aec-%03d.wav", name, idx))
	if err := WriteWav16(file, audio.TargetRate, samples); err != nil {
		t.Fatalf("write aec segment: %v", err)
	}
	return LongRunSegment{File: file, StartFrame: start, Frames: len(samples)}
}

// measureRun runs one rig capture and its offline pass end to end and files
// the report.
func measureRun(t *testing.T, reg *observe.Registry, outDir string, cfg LongRunConfig, kind, farContent string) *runReport {
	t.Helper()
	res, err := RunLoopbackLong(reg, context.Background(), cfg)
	if err != nil {
		t.Fatalf("long-window rig %s: %v", cfg.Name, err)
	}
	if len(res.Segments) == 0 || res.MixFar == nil {
		t.Fatalf("long-window rig %s returned no segments", cfg.Name)
	}

	// Head alignment on the first 3.4s of capture - the aec3-r1 protocol
	// shape, so lag and strength stay on the baseline ruler.
	rate, headPcm, err := ReadWav16(res.Segments[0].File)
	if err != nil {
		t.Fatalf("read first segment: %v", err)
	}
	if rate != res.MixRate {
		t.Fatalf("first segment rate %d, want %d", rate, res.MixRate)
	}
	nearHead := PCM16ToFloat(audio.NewResampler(res.MixRate, audio.TargetRate).Process(headPcm))
	if len(nearHead) > headWindowSamples {
		nearHead = nearHead[:headWindowSamples]
	}
	farHead := make([]float32, len(nearHead))
	NewFarClock(res.MixFar, res.MixRate).NextFrame(farHead)
	lag, strength := AlignByEnvelopeCorrelation(farHead, nearHead, DefaultMaxLag)
	if lag < 0 {
		lag = 0
	}
	frac := FractionalLag(farHead, nearHead, lag)
	delayMs := lag * 1000 / audio.TargetRate

	pass := longAecPass(t, res, lag, delayMs, outDir, cfg.Name, cfg.SegmentSeconds)

	rep := &runReport{
		Name: cfg.Name, Kind: kind, FarContent: farContent, Seconds: cfg.Seconds,
		Reading: res, AlignLag: lag, AlignStrength: strength, FractionalLag: frac,
		DelayMs: delayMs, AecSegments: pass.segs,
		Aec3CpuPercent: pass.cpuPct, Aec3CpuPerAudio: pass.cpuPA,
		ConvergenceWindow: -1,
	}
	if len(pass.headFar) == headWindowSamples {
		rep.HeadERLE = curveJSON(ERLEdB(pass.headFar, pass.headNear, pass.headOut, lag))
	}
	n := len(pass.acc.Windows)
	rep.StartupERLE = curveJSON(pass.acc.SpanERLE(0, minInt(5, n)))
	steadyFrom := n - 30
	if steadyFrom < 5 {
		steadyFrom = 0
	}
	rep.SteadyERLE = curveJSON(pass.acc.SpanERLE(steadyFrom, n))
	rep.FullERLE = curveJSON(pass.acc.SpanERLE(0, n))
	for _, w := range pass.acc.Windows {
		rep.Windows = append(rep.Windows, windowPoint{Window: w.Index, ERLE: curveJSON(w.ERLEdB()), Active: w.ActiveFrames})
	}
	for _, a := range pass.alarms {
		rep.Alarms = append(rep.Alarms, jsonAlarm(a))
	}
	// Convergence: the first window that sits within 3dB of the steady-state
	// median (the "how long until the engine is there" reading).
	if med, ok := steadyMedian(pass.acc, steadyFrom, n); ok {
		for _, w := range pass.acc.Windows {
			if db := w.ERLEdB(); !math.IsNaN(db) && db >= med-3 {
				rep.ConvergenceWindow = w.Index
				break
			}
		}
	}
	// Cold runs get an outcome classification on the aec3-r1 ruler: the
	// double peak separates at 15dB (good 24.5-30.7, collapsed 8.88-8.91).
	if kind == "cold" {
		if f, ok := rep.FullERLE.(float64); ok && !math.IsNaN(f) && f < 15 {
			rep.Outcome = "collapsed"
		} else {
			rep.Outcome = "good"
		}
	}

	writeJSON(t, filepath.Join(outDir, cfg.Name+".json"), rep)
	writeJSON(t, filepath.Join(outDir, cfg.Name+"-curve.json"), rep.Windows)
	t.Logf("[%s] lag=%d(+%.2f) strength=%.3f delay=%dms headERLE=%v startup=%v steady=%v full=%v convWin=%d alarms=%d cpu=%.2f%% ratio=%dppm segs=%d",
		cfg.Name, lag, frac, strength, delayMs, rep.HeadERLE, rep.StartupERLE, rep.SteadyERLE, rep.FullERLE,
		rep.ConvergenceWindow, len(rep.Alarms), rep.Aec3CpuPercent, res.RatioPpm, len(res.Segments))
	return rep
}

// steadyMedian is the median window ERLE of the steady span (ok=false when
// the span has no far-active windows).
func steadyMedian(acc *ERLEAccum, from, to int) (float64, bool) {
	if to <= from {
		return 0, false
	}
	var vals []float64
	for _, w := range acc.Windows[from:to] {
		if db := w.ERLEdB(); !math.IsNaN(db) && !math.IsInf(db, 0) {
			vals = append(vals, db)
		}
	}
	if len(vals) == 0 {
		return 0, false
	}
	return median(vals), true
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// TestSpikeAec3LongWindow is the 59-w-r1 reading. It only records numbers;
// the PASS/NO-GO judgment stays with the orchestrator/59-v2.
func TestSpikeAec3LongWindow(t *testing.T) {
	spikeGate(t)
	out := longWindowOutDir(t)
	fix := longFixtureSetBuild(t)
	reg := observe.NewRegistry()

	// Stage A: the long-window runs (PASS condition (a)).
	longDefs := []struct {
		name   string
		far    []float32
		farLbl string
		secs   float64
	}{
		{"long-tts-60a", fix.tts, "tts", 60},
		{"long-tts-60b", fix.tts, "tts", 60},
		{"long-music-60", fix.music, "music", 60},
		{"long-tts-120", fix.tts, "tts", 120},
	}
	longReps := map[string]*runReport{}
	for _, r := range longDefs {
		longReps[r.name] = measureRun(t, reg, out, LongRunConfig{
			Volume: 0.5, Far: r.far, FarRate: audio.TargetRate,
			Seconds: r.secs, SegmentSeconds: 5, OutDir: out, Name: r.name,
		}, "long", r.farLbl)
	}

	// Stage B: cold 3.4s runs on three head variants (collapse
	// characterization; also fresh samples for the gate's collapse state).
	coldDefs := []struct {
		name   string
		far    []float32
		farLbl string
		secs   float64
	}{
		{"cold-tts-1", fix.tts, "tts", 3.4},
		{"cold-tts-2", fix.tts, "tts", 3.4},
		{"cold-ttsgap-1", fix.ttsGap, "tts-gap0.5", 3.9},
		{"cold-ttsgap-2", fix.ttsGap, "tts-gap0.5", 3.9},
		{"cold-music-1", fix.music, "music", 3.4},
		{"cold-music-2", fix.music, "music", 3.4},
	}
	coldReps := map[string]*runReport{}
	coldTable := map[string]any{}
	for _, r := range coldDefs {
		rep := measureRun(t, reg, out, LongRunConfig{
			Volume: 0.5, Far: r.far, FarRate: audio.TargetRate,
			Seconds: r.secs, SegmentSeconds: 5, OutDir: out, Name: r.name,
		}, "cold", r.farLbl)
		coldReps[r.name] = rep
		coldTable[r.name] = map[string]any{
			"far": rep.FarContent, "full_erle_db": rep.FullERLE,
			"fractional_lag": rep.FractionalLag, "outcome": rep.Outcome,
		}
	}

	// Stage C: the combined barge-in gate (PASS condition (b)).
	gates := gateStage(t, out, longReps, coldReps)

	writeJSON(t, filepath.Join(out, "longrun-summary.json"), map[string]any{
		"long_runs": longReps, "cold_runs": coldTable, "gates": gates,
	})
}

// gateMaterial is one audio span the gate is measured on.
type gateMaterial struct {
	Name    string
	State   string // "good" | "collapse" | "control"
	Source  string
	Samples []float32
}

// gateEpisodeReport is one VAD episode's gate reading.
type gateEpisodeReport struct {
	StartSample int     `json:"start_sample"` // approximate: the feeding chunk quantizes it
	Samples     int     `json:"samples"`
	RmsDbfs     float64 `json:"rms_dbfs"`
	AsrText     string  `json:"asr_text"`
}

// gateReport is one material's three-tier gate table.
type gateReport struct {
	Material     string                   `json:"material"`
	State        string                   `json:"state"` // good | collapse | control
	Source       string                   `json:"source"`
	DurationSec  float64                  `json:"duration_sec"`
	Starts       int                      `json:"vad_starts"`
	StartsPer10  float64                  `json:"vad_starts_per_10min"`
	Episodes     []gateEpisodeReport      `json:"episodes"`
	WholeAsr     string                   `json:"whole_stream_asr"`
	Tiers        map[string]GateTierCount `json:"tiers_by_energy_threshold_dbfs"`
	RateAllPer10 map[string]float64       `json:"tier3_per_10min_by_threshold"`
}

// gateEnergyThresholds is the energy-gate axis the tiers are tallied on: a
// catch-all floor, the middle of the usual noise-gate range, and a strict
// cut. The episode RMS distributions in the same report let the reader move
// the cut anywhere else.
var gateEnergyThresholds = []float64{-60, -45, -30}

// spanFromAecSegments reads samples [from, to) (16k indices) out of a run's
// AEC-output segments on disk.
func spanFromAecSegments(t *testing.T, segs []LongRunSegment, from, to int) []float32 {
	t.Helper()
	var out []float32
	for _, seg := range segs {
		segEnd := seg.StartFrame + seg.Frames
		if segEnd <= from {
			continue
		}
		if seg.StartFrame >= to {
			break
		}
		_, pcm, err := ReadWav16(seg.File)
		if err != nil {
			t.Fatalf("read aec segment %s: %v", seg.File, err)
		}
		lo := from - seg.StartFrame
		if lo < 0 {
			lo = 0
		}
		hi := to - seg.StartFrame
		if hi > len(pcm) {
			hi = len(pcm)
		}
		if hi > lo {
			out = append(out, PCM16ToFloat(pcm[lo:hi])...)
		}
	}
	return out
}

// runSpanSeconds is a run's captured audio length in seconds.
func runSpanSeconds(rep *runReport) float64 {
	return float64(rep.Reading.CapturedFrames) / float64(rep.Reading.MixRate)
}

// gateStage builds the good/collapse/control materials from the stage A/B
// artifacts and runs the three-tier gate on each.
func gateStage(t *testing.T, out string, longReps, coldReps map[string]*runReport) map[string]gateReport {
	t.Helper()
	var mats []gateMaterial

	// Good-state long-window steady spans: [10s, min(70s, end-2s)) of every
	// long run. No ERLE quality filter here: the gate is measured on the
	// product-shaped steady material as it actually is, whatever its ERLE.
	longNames := make([]string, 0, len(longReps))
	for name := range longReps {
		longNames = append(longNames, name)
	}
	sort.Strings(longNames)
	for _, name := range longNames {
		rep := longReps[name]
		from := 10 * audio.TargetRate
		to := int((runSpanSeconds(rep) - 2) * audio.TargetRate)
		if cap := from + 60*audio.TargetRate; to > cap {
			to = cap
		}
		if to <= from {
			continue
		}
		mats = append(mats, gateMaterial{
			Name: name + "-steady", State: "good",
			Source:  fmt.Sprintf("%s AEC output, steady span", name),
			Samples: spanFromAecSegments(t, rep.AecSegments, from, to),
		})
	}

	// Cold runs: each one whole, state = its measured outcome.
	coldNames := make([]string, 0, len(coldReps))
	for name := range coldReps {
		coldNames = append(coldNames, name)
	}
	sort.Strings(coldNames)
	for _, name := range coldNames {
		rep := coldReps[name]
		mats = append(mats, gateMaterial{
			Name: name, State: rep.Outcome,
			Source:  fmt.Sprintf("%s AEC output, whole cold run (fullERLE=%v)", name, rep.FullERLE),
			Samples: spanFromAecSegments(t, rep.AecSegments, 0, 1<<30),
		})
	}

	// The archived collapse material: the 59-v2 verdict leg's AEC output of
	// a collapsed live capture, kept as the collapse-state row of last
	// resort (labelled as archived, never as this session's own run).
	if _, err := os.Stat(archivedCollapseWav); err == nil {
		_, pcm, err := ReadWav16(archivedCollapseWav)
		if err != nil {
			t.Logf("archived collapse wav unreadable (%v); collapse-state rows stand on this session's runs", err)
		} else {
			mats = append(mats, gateMaterial{
				Name: "archived-59v2-collapse", State: "collapse",
				Source:  "59-v2 verdict leg archived AEC output (fresh-rig run, 8.91dB, ASR two syllables)",
				Samples: PCM16ToFloat(pcm),
			})
		}
	}

	// The positive control: real TTS speech added to a good-state material
	// at two absolute levels. The gate must fire on these (a dead gate would
	// also fire on nothing - this row is what keeps the false-trigger table
	// honest). Levels are absolute RMS targets, not SNRs: a well-converged
	// AEC output is near-silent, so an SNR definition would degenerate.
	controlBase, hasControl := "", false
	for _, m := range mats {
		if m.State == "good" {
			controlBase, hasControl = m.Name, true
			break
		}
	}
	if hasControl {
		var base []float32
		for _, m := range mats {
			if m.Name == controlBase {
				base = m.Samples
				break
			}
		}
		speech16 := ttsSpeak(t, "今天天气怎么样，帮我看一下日历。")
		for _, lvl := range []struct {
			name     string
			targetDb float64
		}{{"control-26", -26}, {"control-36", -36}} {
			gain := math.Pow(10, lvl.targetDb/20) / math.Max(1e-9, speechRMS(speech16))
			mixed := make([]float32, len(base))
			copy(mixed, base)
			for i, v := range speech16 {
				if i >= len(mixed) {
					break
				}
				mixed[i] += float32(gain) * v
			}
			mats = append(mats, gateMaterial{
				Name: lvl.name, State: "control",
				Source:  fmt.Sprintf("%s + TTS speech at %.0f dBFS RMS", controlBase, lvl.targetDb),
				Samples: mixed,
			})
		}
	}

	asr := spikeASR(t)
	reports := map[string]gateReport{}
	for _, m := range mats {
		rep := runGate(t, m, asr)
		reports[m.Name] = rep
		t.Logf("[gate %-24s] state=%-8s dur=%.1fs starts=%d (%.1f/10min) wholeASR=%q episodes=%d",
			m.Name, m.State, rep.DurationSec, rep.Starts, rep.StartsPer10, rep.WholeAsr, len(rep.Episodes))
		for _, e := range rep.Episodes {
			t.Logf("[gate %-24s]   ep@%d rms=%.1fdBFS asr=%q", m.Name, e.StartSample, e.RmsDbfs, e.AsrText)
		}
		for _, thr := range gateEnergyThresholds {
			k := fmt.Sprintf("%g", thr)
			t.Logf("[gate %-24s]   thr=%s tiers=%+v tier3/10min=%.2f", m.Name, k, rep.Tiers[k], rep.RateAllPer10[k])
		}
	}
	writeJSON(t, filepath.Join(out, "gates.json"), reports)
	return reports
}

// speechRMS is the RMS of a float32 span (linear, for gain math).
func speechRMS(samples []float32) float64 {
	var sum float64
	for _, v := range samples {
		sum += float64(v) * float64(v)
	}
	if len(samples) == 0 {
		return 0
	}
	return math.Sqrt(sum / float64(len(samples)))
}

// runGate feeds one material through the product gate chain: Detector +
// Endpointer (the product's segmentation) chunks in 0.5s, the closing
// event's utterance audio is the episode, and the three tiers are tallied
// per energy threshold.
func runGate(t *testing.T, m gateMaterial, asr *speech.StreamingAsr) gateReport {
	t.Helper()
	det := spikeNewVad(t)
	const chunk = 8000 // 0.5s
	var eps []speech.Endpoint
	fed, lastStart := 0, 0
	starts := 0
	var episodes []gateEpisodeReport
	collect := func(events []speech.Endpoint) {
		for _, e := range events {
			switch e.Kind {
			case speech.EndpointStart:
				starts++
				lastStart = fed
			case speech.EndpointEnd, speech.EndpointMaxDuration, speech.EndpointEndOfStream:
				if len(e.Samples) > 0 {
					episodes = append(episodes, gateEpisodeReport{
						StartSample: lastStart, Samples: len(e.Samples),
						RmsDbfs: rmsDbfs(e.Samples),
						AsrText: spikeTranscribe(t, asr, e.Samples),
					})
				}
			}
		}
	}
	for lo := 0; lo < len(m.Samples); lo += chunk {
		hi := minInt(lo+chunk, len(m.Samples))
		eps = det.Accept(m.Samples[lo:hi], eps[:0])
		fed = hi
		collect(eps)
	}
	eps = det.Finalize(eps[:0])
	collect(eps)

	rep := gateReport{
		Material: m.Name, State: m.State, Source: m.Source,
		DurationSec: float64(len(m.Samples)) / audio.TargetRate,
		Starts:      starts,
		WholeAsr:    spikeTranscribe(t, asr, m.Samples),
		Tiers:       map[string]GateTierCount{}, RateAllPer10: map[string]float64{},
	}
	rep.StartsPer10 = Rate10Min(rep.Starts, rep.DurationSec)
	gateEps := make([]GateEpisode, len(episodes))
	for i, e := range episodes {
		gateEps[i] = GateEpisode{StartSample: e.StartSample, Samples: e.Samples, RmsDbfs: e.RmsDbfs, AsrText: e.AsrText}
	}
	rep.Episodes = episodes
	for _, thr := range gateEnergyThresholds {
		k := fmt.Sprintf("%g", thr)
		g := GateTiers(gateEps, thr)
		rep.Tiers[k] = g
		rep.RateAllPer10[k] = Rate10Min(g.AllPass, rep.DurationSec)
	}
	return rep
}

// probeRunJSON is the slice of the on-disk run report the probe needs.
type probeRunJSON struct {
	Name       string `json:"name"`
	FarContent string `json:"far_content"`
	AlignLag   int    `json:"align_lag_samples"`
	FullERLE   any    `json:"full_erle_db"`
	Reading    struct {
		MixRate      int              `json:"mix_rate"`
		TargetFrames int              `json:"target_frames"`
		Segments     []LongRunSegment `json:"segments"`
	} `json:"reading"`
	AecSegments []LongRunSegment `json:"aec_segments"`
}

// probeWindow is one 1s probe reading: the ERLE from the run's own curve,
// the capture's RMS and correlation against the far reference, and the AEC
// output's RMS. The correlation is the pollution discriminator: a capture
// that still IS the echo (engine-side failure) keeps corr high; foreign
// energy in the endpoint mix decorrelates it.
type probeWindow struct {
	Window     int     `json:"window"`
	ERLE       any     `json:"erle_db"`
	CaptureRms float64 `json:"capture_rms_dbfs"`
	Corr       float64 `json:"corr_vs_far"`
	AecRms     float64 `json:"aec_out_rms_dbfs"`
	Class      string  `json:"class"` // ok | event-capture | event-engine | inactive
}

// TestSpikeLongWindowProbe is the offline second look at the stage A runs
// (no rig): it rebuilds the far fixtures (deterministic - the probe
// self-checks on each run's window 0, whose in-run alignment strength is
// known), correlates every capture window against the far reference, and
// splits the ERLE events into capture-side (foreign energy) and engine-side.
// It then re-runs the gate with the full material set, adding the long-run
// event windows as collapse-state materials. Requires the stage A artifacts
// in the out dir (run TestSpikeAec3LongWindow first).
func TestSpikeLongWindowProbe(t *testing.T) {
	spikeGate(t)
	out := longWindowOutDir(t)
	names := []string{"long-tts-60a", "long-tts-60b", "long-music-60", "long-tts-120"}
	probes := map[string]any{}
	eventSpans := map[string][][2]int{}
	for _, name := range names {
		b, err := os.ReadFile(filepath.Join(out, name+".json"))
		if err != nil {
			continue // run not present; the driver decides what is mandatory
		}
		var run probeRunJSON
		if err := json.Unmarshal(b, &run); err != nil {
			t.Fatalf("parse %s.json: %v", name, err)
		}
		// Rebuild the far fixture (the same deterministic pipeline the rig
		// used) and the capture stream.
		var farMix []float32
		switch run.FarContent {
		case "tts":
			farMix = ResampleLinear(PCM16ToFloat(spikeTTSFixture(t)), audio.TargetRate, run.Reading.MixRate)
		case "music":
			farMix = ResampleLinear(SynthMusic(audio.TargetRate, 8), audio.TargetRate, run.Reading.MixRate)
		default:
			t.Fatalf("unknown far content %q for %s", run.FarContent, name)
		}
		var nearMix []int16
		for _, seg := range run.Reading.Segments {
			rate, pcm, err := ReadWav16(seg.File)
			if err != nil {
				t.Fatalf("read %s: %v", seg.File, err)
			}
			if rate != run.Reading.MixRate {
				t.Fatalf("%s: segment rate %d, want %d", seg.File, rate, run.Reading.MixRate)
			}
			nearMix = append(nearMix, pcm...)
		}
		near16 := PCM16ToFloat(audio.NewResampler(run.Reading.MixRate, audio.TargetRate).Process(nearMix))
		far16 := NewFarClock(farMix, run.Reading.MixRate)
		farStream := make([]float32, len(near16))
		farStream = farStream[:0]
		buf := make([]float32, 4096)
		for len(farStream) < len(near16) {
			n := len(near16) - len(farStream)
			if n > len(buf) {
				n = len(buf)
			}
			farStream = farStream[:len(farStream)+n]
			far16.NextFrame(farStream[len(farStream)-n:])
		}
		lag := run.AlignLag
		if lag < 0 {
			lag = 0
		}
		corr := WindowedCorrelation(farStream, near16, lag, audio.TargetRate, audio.TargetRate)

		// AEC-output per-window RMS from the aec segments.
		aecTotal := spanFromAecSegments(t, run.AecSegments, 0, 1<<30)
		var curve []windowPoint
		if cb, err := os.ReadFile(filepath.Join(out, name+"-curve.json")); err == nil {
			_ = json.Unmarshal(cb, &curve)
		}
		var rows []probeWindow
		var events [][2]int
		nWin := len(near16) / audio.TargetRate
		for w := 0; w < nWin; w++ {
			lo, hi := w*audio.TargetRate, (w+1)*audio.TargetRate
			row := probeWindow{Window: w, CaptureRms: rmsDbfs(near16[lo:hi]), AecRms: rmsDbfs(aecTotal[lo:hi])}
			if w < len(corr) {
				row.Corr = corr[w]
			}
			if w < len(curve) {
				row.ERLE = curve[w].ERLE
			}
			if f, ok := row.ERLE.(float64); ok && !math.IsNaN(f) {
				switch {
				case f >= 20:
					row.Class = "ok"
				case row.Corr < 0.8:
					row.Class = "event-capture"
				default:
					row.Class = "event-engine"
				}
				if row.Class != "ok" {
					if n := len(events); n > 0 && events[n-1][1] == w {
						events[n-1][1] = w + 1
					} else {
						events = append(events, [2]int{w, w + 1})
					}
				}
			} else {
				row.Class = "inactive"
			}
			rows = append(rows, row)
		}
		eventSpans[name] = events
		probes[name] = map[string]any{"windows": rows, "events": events}
		nCapture := 0
		nEngine := 0
		for _, r := range rows {
			if r.Class == "event-capture" {
				nCapture++
			}
			if r.Class == "event-engine" {
				nEngine++
			}
		}
		t.Logf("[probe %s] windows=%d events=%d (capture-side=%d engine-side=%d) corr[0]=%.3f rms[0]=%.1f",
			name, nWin, len(events), nCapture, nEngine, corr[0], rows[0].CaptureRms)
	}
	writeJSON(t, filepath.Join(out, "probe.json"), probes)

	// The gate rerun with the full material set: every long run's steady
	// span, every long run's event windows (collapse-state material the
	// cold runs cannot provide - a collapse DURING a converged session),
	// the cold runs, and the archived collapse wav.
	longReps := map[string]*runReport{}
	coldReps := map[string]*runReport{}
	for _, name := range names {
		b, err := os.ReadFile(filepath.Join(out, name+".json"))
		if err != nil {
			continue
		}
		var rep runReport
		if err := json.Unmarshal(b, &rep); err != nil {
			t.Fatalf("parse %s.json: %v", name, err)
		}
		longReps[name] = &rep
	}
	for _, name := range []string{"cold-tts-1", "cold-tts-2", "cold-ttsgap-1", "cold-ttsgap-2", "cold-music-1", "cold-music-2"} {
		b, err := os.ReadFile(filepath.Join(out, name+".json"))
		if err != nil {
			continue
		}
		var rep runReport
		if err := json.Unmarshal(b, &rep); err != nil {
			t.Fatalf("parse %s.json: %v", name, err)
		}
		coldReps[name] = &rep
	}

	var mats []gateMaterial
	for _, name := range []string{"long-tts-60a", "long-tts-60b", "long-music-60", "long-tts-120"} {
		rep, ok := longReps[name]
		if !ok {
			continue
		}
		from := 10 * audio.TargetRate
		to := int((runSpanSeconds(rep) - 2) * audio.TargetRate)
		if cap := from + 60*audio.TargetRate; to > cap {
			to = cap
		}
		if to > from {
			mats = append(mats, gateMaterial{
				Name: name + "-steady", State: "good",
				Source:  name + " AEC output, steady span",
				Samples: spanFromAecSegments(t, rep.AecSegments, from, to),
			})
		}
		for k, span := range eventSpans[name] {
			lo, hi := span[0]*audio.TargetRate, span[1]*audio.TargetRate
			if hi > to {
				hi = to
			}
			if hi <= lo {
				continue
			}
			mats = append(mats, gateMaterial{
				Name: fmt.Sprintf("%s-event%d", name, k), State: "collapse",
				Source:  fmt.Sprintf("%s AEC output, ERLE event window [%d, %d)", name, span[0], span[1]),
				Samples: spanFromAecSegments(t, rep.AecSegments, lo, hi),
			})
		}
	}
	for _, name := range []string{"cold-tts-1", "cold-tts-2", "cold-ttsgap-1", "cold-ttsgap-2", "cold-music-1", "cold-music-2"} {
		rep, ok := coldReps[name]
		if !ok {
			continue
		}
		mats = append(mats, gateMaterial{
			Name: name, State: rep.Outcome,
			Source:  fmt.Sprintf("%s AEC output, whole cold run (fullERLE=%v)", name, rep.FullERLE),
			Samples: spanFromAecSegments(t, rep.AecSegments, 0, 1<<30),
		})
	}
	if _, err := os.Stat(archivedCollapseWav); err == nil {
		if _, pcm, err := ReadWav16(archivedCollapseWav); err == nil {
			mats = append(mats, gateMaterial{
				Name: "archived-59v2-collapse", State: "collapse",
				Source:  "59-v2 verdict leg archived AEC output (fresh-rig run, 8.91dB)",
				Samples: PCM16ToFloat(pcm),
			})
		}
	}

	asr := spikeASR(t)
	reports := map[string]gateReport{}
	sorted := make([]string, 0, len(mats))
	for _, m := range mats {
		sorted = append(sorted, m.Name)
	}
	sort.Strings(sorted)
	for _, m := range mats {
		rep := runGate(t, m, asr)
		reports[m.Name] = rep
		t.Logf("[gate2 %-24s] state=%-8s dur=%.1fs starts=%d (%.1f/10min) wholeASR=%q episodes=%d allPass(-45)=%d",
			m.Name, m.State, rep.DurationSec, rep.Starts, rep.StartsPer10, rep.WholeAsr, len(rep.Episodes),
			rep.Tiers["-45"].AllPass)
		for _, e := range rep.Episodes {
			t.Logf("[gate2 %-24s]   ep@%d rms=%.1fdBFS asr=%q", m.Name, e.StartSample, e.RmsDbfs, e.AsrText)
		}
	}
	writeJSON(t, filepath.Join(out, "gates-probe.json"), reports)
}

// autocorrRunJSON is the self-correlation probe's per-run output.
type autocorrRunJSON struct {
	Name        string          `json:"name"`
	Period      int             `json:"loop_period_frames"`
	PeriodScore float64         `json:"period_search_score"`
	Windows     []autocorrWin   `json:"windows"`
	Events      []autocorrEvent `json:"events"`
}

type autocorrWin struct {
	Window     int     `json:"window"`
	ERLE       any     `json:"erle_db"`
	CaptureRms float64 `json:"capture_rms_dbfs"`
	SelfCorr   float64 `json:"self_corr_at_period"`
	Class      string  `json:"class"` // ok | event-capture | event-pure | inactive
}

type autocorrEvent struct {
	StartWin    int     `json:"start_window"`
	EndWin      int     `json:"end_window"`
	MinERLE     float64 `json:"min_erle_db"`
	MinSelfCorr float64 `json:"min_self_corr"`
	ToFlushSec  float64 `json:"to_nearest_flush_sec"`
	Class       string  `json:"class"`
}

// TestSpikeLongWindowAutocorr is the reference-free pollution discriminator
// for the stage A runs (no rig, no fixture rebuild): the capture is the echo
// of a seamlessly LOOPED fixture, so a clean capture window repeats itself
// exactly one loop period earlier - self-correlation at the period stays at
// the echo path's own fidelity (about 0.98 here), while ANY foreign energy
// in the endpoint mix decorrelates it. Combined with the run's ERLE curve
// this splits every mid-run event into capture-side (the input changed:
// pollution) and capture-pure (the input is still the echo: engine-side or
// AEC-pass artifact). The period itself is found from the data on each
// run's cleanest stretch, so nothing about the fixture needs to be
// reproduced. Requires the stage A artifacts (run TestSpikeAec3LongWindow
// first).
func TestSpikeLongWindowAutocorr(t *testing.T) {
	spikeGate(t)
	out := longWindowOutDir(t)
	results := map[string]autocorrRunJSON{}
	for _, name := range []string{"long-tts-60a", "long-tts-60b", "long-music-60", "long-tts-120"} {
		b, err := os.ReadFile(filepath.Join(out, name+".json"))
		if err != nil {
			t.Fatalf("stage A artifacts missing (%s): %v", name, err)
		}
		var run probeRunJSON
		if err := json.Unmarshal(b, &run); err != nil {
			t.Fatalf("parse %s.json: %v", name, err)
		}
		var near []float32
		for _, seg := range run.Reading.Segments {
			rate, pcm, err := ReadWav16(seg.File)
			if err != nil {
				t.Fatalf("read %s: %v", seg.File, err)
			}
			if rate != run.Reading.MixRate {
				t.Fatalf("%s: segment rate %d, want %d", seg.File, rate, run.Reading.MixRate)
			}
			near = append(near, PCM16ToFloat(pcm)...)
		}
		var curve []windowPoint
		if cb, err := os.ReadFile(filepath.Join(out, name+"-curve.json")); err == nil {
			_ = json.Unmarshal(cb, &curve)
		}
		erleOf := func(w int) (float64, bool) {
			if w < len(curve) {
				if f, ok := curve[w].ERLE.(float64); ok && !math.IsNaN(f) {
					return f, true
				}
			}
			return 0, false
		}
		// Period search: the cleanest stretch is the highest-ERLE 1s window;
		// the loop period is the self-correlation peak within +-1000 frames
		// of the nominal fixture length (3.4s TTS or 8s music).
		nominal := run.Reading.MixRate * 34 / 10
		if run.FarContent == "music" {
			nominal = run.Reading.MixRate * 8
		}
		bestWin, bestERLE := 0, math.Inf(-1)
		// The reference window [wlo, wlo+1s) needs a full period of capture
		// before it: start the search after nominal + 1000 frames.
		for w := (nominal+1000)/run.Reading.MixRate + 1; w < len(near)/run.Reading.MixRate; w++ {
			if f, ok := erleOf(w); ok && f > bestERLE {
				bestERLE, bestWin = f, w
			}
		}
		wlo, whi := bestWin*run.Reading.MixRate, minInt((bestWin+1)*run.Reading.MixRate, len(near))
		// The period lives around the nominal fixture length; the reference
		// window must start no earlier than one full period.
		searchLo := nominal - 1000
		if searchLo < 0 {
			searchLo = 0
		}
		searchHi := nominal + 1000
		if searchHi > wlo { // the reference index i-p must stay >= 0
			searchHi = wlo
		}
		if searchHi > len(near) {
			searchHi = len(near)
		}
		bestP, bestScore := 0, math.Inf(-1)
		for p := searchLo; p < searchHi; p++ {
			var ab, aa, bb float64
			for i := wlo; i < whi; i++ {
				a := float64(near[i])
				bv := float64(near[i-p])
				ab += a * bv
				aa += a * a
				bb += bv * bv
			}
			if aa > 0 && bb > 0 {
				if s := ab / math.Sqrt(aa*bb); s > bestScore {
					bestScore, bestP = s, p
				}
			}
		}
		if bestP == 0 {
			t.Fatalf("%s: loop period not found", name)
		}

		res := autocorrRunJSON{Name: name, Period: bestP, PeriodScore: bestScore}
		nWin := len(near) / run.Reading.MixRate
		for w := 0; w < nWin; w++ {
			lo, hi := w*run.Reading.MixRate, minInt((w+1)*run.Reading.MixRate, len(near))
			row := autocorrWin{Window: w, CaptureRms: rmsDbfs(near[lo:hi]), Class: "inactive"}
			f, ok := erleOf(w)
			if !ok || lo < bestP {
				res.Windows = append(res.Windows, row)
				continue
			}
			row.ERLE = f
			var ab, aa, bb float64
			for i := lo; i < hi; i++ {
				a := float64(near[i])
				bv := float64(near[i-bestP])
				ab += a * bv
				aa += a * a
				bb += bv * bv
			}
			if aa > 0 && bb > 0 {
				row.SelfCorr = ab / math.Sqrt(aa*bb)
			}
			switch {
			case f >= 20:
				row.Class = "ok"
			case row.SelfCorr < 0.8:
				row.Class = "event-capture"
			default:
				row.Class = "event-pure"
			}
			if n := len(res.Events); row.Class != "ok" && n > 0 &&
				res.Events[n-1].EndWin == w && res.Events[n-1].Class == row.Class {
				res.Events[n-1].EndWin = w + 1
				if f < res.Events[n-1].MinERLE {
					res.Events[n-1].MinERLE = f
				}
				if row.SelfCorr < res.Events[n-1].MinSelfCorr {
					res.Events[n-1].MinSelfCorr = row.SelfCorr
				}
			} else if row.Class != "ok" {
				res.Events = append(res.Events, autocorrEvent{
					StartWin: w, EndWin: w + 1, MinERLE: f, MinSelfCorr: row.SelfCorr, Class: row.Class,
				})
			}
			res.Windows = append(res.Windows, row)
		}
		nCapture, nPure := 0, 0
		for i := range res.Events {
			ev := &res.Events[i]
			ev.ToFlushSec = math.Min(float64(ev.StartWin%5), float64(5-ev.StartWin%5))
			if ev.Class == "event-capture" {
				nCapture++
			} else {
				nPure++
			}
		}
		results[name] = res
		t.Logf("[autocorr %s] period=%d (score %.4f) events=%d capture-side=%d capture-pure=%d",
			name, bestP, bestScore, len(res.Events), nCapture, nPure)
		for _, ev := range res.Events {
			t.Logf("[autocorr %s]   event [%d,%d)s class=%s minERLE=%.1f minSelfCorr=%.3f flushDist=%.1fs",
				name, ev.StartWin, ev.EndWin, ev.Class, ev.MinERLE, ev.MinSelfCorr, ev.ToFlushSec)
		}
	}
	writeJSON(t, filepath.Join(out, "autocorr.json"), results)
}

// delayTrackRunJSON is the delay-trajectory probe's per-run output.
type delayTrackRunJSON struct {
	Name     string            `json:"name"`
	RefBase  int               `json:"ref_base_window"`
	Track    []delayTrackPoint `json:"track"`
	FarTrack []delayTrackPoint `json:"far_track_when_exact_far_available"`
}

type delayTrackPoint struct {
	Window   int     `json:"window"`
	Lag      int     `json:"lag_samples"` // relative to the extracted reference
	Strength float64 `json:"strength"`
	ERLE     any     `json:"erle_db"`
}

// TestSpikeLongWindowDelayTrack answers the last open question about the
// capture-pure ERLE events: did the render-to-capture DELAY step (a rig
// real-time artifact: the reference misaligns and the engine re-adapts), or
// did the capture delay stay constant while the engine lost the echo anyway
// (an AEC3-internal instability on periodic content)? The reference is the
// capture's own cleanest stretch (the echo is a scaled delayed copy of the
// far, and correlation is scale-invariant), so no fixture rebuild is
// needed; the music run additionally tracks against the exact deterministic
// SynthMusic far as a cross-check. Requires the stage A artifacts.
func TestSpikeLongWindowDelayTrack(t *testing.T) {
	spikeGate(t)
	out := longWindowOutDir(t)
	results := map[string]delayTrackRunJSON{}
	for _, name := range []string{"long-tts-60a", "long-tts-60b", "long-music-60", "long-tts-120"} {
		b, err := os.ReadFile(filepath.Join(out, name+".json"))
		if err != nil {
			t.Fatalf("stage A artifacts missing (%s): %v", name, err)
		}
		var run probeRunJSON
		if err := json.Unmarshal(b, &run); err != nil {
			t.Fatalf("parse %s.json: %v", name, err)
		}
		var nearMix []int16
		for _, seg := range run.Reading.Segments {
			_, pcm, err := ReadWav16(seg.File)
			if err != nil {
				t.Fatalf("read %s: %v", seg.File, err)
			}
			nearMix = append(nearMix, pcm...)
		}
		near16 := PCM16ToFloat(audio.NewResampler(run.Reading.MixRate, audio.TargetRate).Process(nearMix))
		var curve []windowPoint
		if cb, err := os.ReadFile(filepath.Join(out, name+"-curve.json")); err == nil {
			_ = json.Unmarshal(cb, &curve)
		}
		erleOf := func(w int) (float64, bool) {
			if w < len(curve) {
				if f, ok := curve[w].ERLE.(float64); ok && !math.IsNaN(f) {
					return f, true
				}
			}
			return 0, false
		}
		// The reference: two seconds of the cleanest (highest-ERLE) stretch,
		// advanced by the measured lag so ref[j] ~ far[j].
		lag16 := run.AlignLag
		bestWin, bestERLE := 0, math.Inf(-1)
		for w := 0; w+2 < len(near16)/audio.TargetRate; w++ {
			if f, ok := erleOf(w); ok && f > bestERLE {
				if f2, ok2 := erleOf(w + 1); ok2 && f2 > 20 {
					bestERLE, bestWin = f, w
				}
			}
		}
		base := bestWin*audio.TargetRate + lag16
		if base+2*audio.TargetRate > len(near16) {
			t.Fatalf("%s: no clean 2s stretch for a reference", name)
		}
		ref := append([]float32(nil), near16[base:base+2*audio.TargetRate]...)

		res := delayTrackRunJSON{Name: name, RefBase: bestWin}
		for w := 0; w*audio.TargetRate+audio.TargetRate <= len(near16); w++ {
			lo := w * audio.TargetRate
			hi := lo + audio.TargetRate
			lag, strength := AlignByEnvelopeCorrelation(ref, near16[lo:hi], 3200)
			pt := delayTrackPoint{Window: w, Lag: lag, Strength: strength}
			if f, ok := erleOf(w); ok {
				pt.ERLE = f
			} else {
				pt.ERLE = "inactive"
			}
			res.Track = append(res.Track, pt)
		}
		// Cross-check on the music run: the exact deterministic far.
		if run.FarContent == "music" {
			farMix := ResampleLinear(SynthMusic(audio.TargetRate, 8), audio.TargetRate, run.Reading.MixRate)
			far16 := NewFarClock(farMix, run.Reading.MixRate)
			farStream := make([]float32, len(near16))
			buf := make([]float32, 4096)
			for done := 0; done < len(farStream); {
				n := len(farStream) - done
				if n > len(buf) {
					n = len(buf)
				}
				far16.NextFrame(buf[:n])
				copy(farStream[done:], buf[:n])
				done += n
			}
			for w := 0; w*audio.TargetRate+audio.TargetRate <= len(near16); w++ {
				lo := w * audio.TargetRate
				lag, strength := AlignByEnvelopeCorrelation(farStream[lo:lo+audio.TargetRate], near16[lo:lo+audio.TargetRate], 3200)
				res.FarTrack = append(res.FarTrack, delayTrackPoint{Window: w, Lag: lag, Strength: strength})
			}
		}
		results[name] = res
		// Compact trajectory log: lag per window where it moves or ERLE dips.
		shifts := 0
		prev := res.Track[len(res.Track)/2].Lag
		line := ""
		for _, pt := range res.Track {
			if pt.Lag != prev {
				shifts++
				line += fmt.Sprintf(" w%d:%d(%s)", pt.Window, pt.Lag, fmtERLE(pt.ERLE))
				prev = pt.Lag
			}
		}
		t.Logf("[delaytrack %s] refWin=%d lag-shift points=%d:%s", name, bestWin, shifts, line)
	}
	writeJSON(t, filepath.Join(out, "delaytrack.json"), results)
}

// fmtERLE renders an ERLE any for a log line.
func fmtERLE(v any) string {
	if f, ok := v.(float64); ok {
		return fmt.Sprintf("%.0f", f)
	}
	return fmt.Sprint(v)
}
