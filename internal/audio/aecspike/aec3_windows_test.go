//go:build windows

package aecspike

// The AEC3 re-measurement leg (dispatch aec3-r1): the same loop-back rig,
// the same seam rate and the same measurement instruments as the 59-r1
// baseline, with the webrtc AEC3 engine (deps.toml
// [webrtc-audio-processing]) inserted into the render-to-capture chain.
// Nothing here runs in CI: the WISP_AEC_SPIKE=1 gate applies, and the
// native binding must be part of the test build (-tags webrtcap) or the
// test skips with that fact recorded - the honest answer, not a fake green.
//
// Run shape (extends the 59-r1 recipe; see
// docs/evidence/s1/aec3r1-remeasure.md):
//
//	go test -c -o <out>\aecspike.test.exe ./internal/audio/aecspike/
//	  (with CGO_CXXFLAGS/CGO_LDFLAGS pointing at the fetched tree)
//	set WISP_AEC_SPIKE=1 WISP_AEC_SPIKE_ROOT=<repo>
//	  WISP_AEC_SPIKE_OUT=D:\tmp\wisp\20261031\aec3r1\out
//	aecspike.test.exe -test.run TestSpikeAec3 -test.v

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/audio/aec"
	"github.com/CarlosShao/wisp/internal/observe"
)

// aec3OutDir defaults to THIS leg's dispatcher-named temp root (never the
// 59-r1 one; artifacts are not mixed between legs).
func aec3OutDir(t *testing.T) string {
	t.Helper()
	dir := os.Getenv("WISP_AEC_SPIKE_OUT")
	if dir == "" {
		dir = `D:\tmp\wisp\20261031\aec3r1\out`
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("aec3 out dir: %v", err)
	}
	return dir
}

// aec3Arm is one (far, near) capture pair the engine is measured on.
type aec3Arm struct {
	far    []float32
	near   []float32
	source string
}

// baselineArm returns the 59-r1 self-echo capture pair that produced the
// 16.8dB baseline reading - the fixture-identity arm. A missing or
// unreadable pair drops the arm; the fresh-rig arm still runs.
func baselineArm(t *testing.T) (aec3Arm, bool) {
	t.Helper()
	dir := `D:\tmp\wisp\20261028\59r1`
	_, farPcm, err := ReadWav16(filepath.Join(dir, "selfecho-far-16k.wav"))
	if err != nil {
		t.Logf("59-r1 far fixture unavailable (%v): fixture-identity arm dropped", err)
		return aec3Arm{}, false
	}
	_, nearPcm, err := ReadWav16(filepath.Join(dir, "selfecho-near-16k.wav"))
	if err != nil {
		t.Logf("59-r1 near fixture unavailable (%v): fixture-identity arm dropped", err)
		return aec3Arm{}, false
	}
	return aec3Arm{
		far: PCM16ToFloat(farPcm), near: PCM16ToFloat(nearPcm),
		source: "59-r1 selfecho wavs (the exact capture behind the 16.8dB baseline)",
	}, true
}

// freshRigArm renders the same TTS fixture the 59-r1 baseline used and
// captures the endpoint loopback, at the baseline's 0.5 session volume.
func freshRigArm(t *testing.T) aec3Arm {
	t.Helper()
	farPcm := PCM16ToFloat(spikeTTSFixture(t))
	reading, capturedMix, err := RunLoopback(observe.NewRegistry(), context.Background(),
		"", 0.5, farPcm, audio.TargetRate)
	if err != nil {
		t.Skipf("loopback rig unavailable on this host: %v", err)
	}
	nearPcm := audio.NewResampler(reading.MixRate, audio.TargetRate).Process(capturedMix)
	return aec3Arm{
		far: farPcm, near: PCM16ToFloat(nearPcm),
		source: "same-day loopback rerun of the same TTS fixture at vol 0.5",
	}
}

// aec3Pass feeds the pair through AEC3 in 10ms lockstep and returns the
// cancelled stream. delayMs is the measured render-to-capture delay in
// whole milliseconds; frames whose length does not fill one 10ms block are
// dropped (tail, under 10ms of audio).
func aec3Pass(t *testing.T, apm *aec.Apm, far, near []float32, delayMs int) []float32 {
	t.Helper()
	frame := aec.FrameSamples(16000)
	out := make([]float32, 0, len(near)/frame*frame)
	for lo := 0; lo+frame <= len(near); lo += frame {
		hi := lo + frame
		if err := apm.ProcessRender(far[lo:hi]); err != nil {
			t.Fatalf("AEC3 render frame at %d: %v", lo, err)
		}
		chunk, err := apm.ProcessCapture(near[lo:hi], delayMs)
		if err != nil {
			t.Fatalf("AEC3 capture frame at %d: %v", lo, err)
		}
		out = append(out, chunk...)
	}
	return out
}

// TestSpikeAec3ErleRemeasure is the aec3-r1 dispatch reading: AEC3 on vs
// off on the same fixture, the ASR self-transcription check, and the
// barge-in discriminator (corr<0.5) readings. The PASS/NO-GO judgment
// stays with the orchestrator; this test only records numbers.
func TestSpikeAec3ErleRemeasure(t *testing.T) {
	spikeGate(t)
	out := aec3OutDir(t)
	if _, err := aec.NewApm(16000); err != nil {
		if errors.Is(err, aec.ErrUnavailable) {
			t.Skipf("AEC3 binding not in this build: %v", err)
		}
		t.Fatalf("NewApm: %v", err)
	}

	arms := map[string]aec3Arm{}
	if arm, ok := baselineArm(t); ok {
		arms["fixture-59r1"] = arm
	}
	arms["fresh-rig"] = freshRigArm(t)

	asr := spikeASR(t)
	vad := spikeNewVad(t)
	type armReport struct {
		Source          string  `json:"source"`
		AlignLag        int     `json:"align_lag_samples"`
		AlignStrength   float64 `json:"align_strength"`
		DelayMs         int     `json:"aec3_delay_ms"`
		SubtractERLEdB  float64 `json:"subtract_erle_db"`
		Aec3ERLEdB      float64 `json:"aec3_erle_db"`
		Aec3Frames      int     `json:"aec3_frames"`
		Aec3CpuPercent  float64 `json:"aec3_cpu_percent_of_core"`
		Aec3CpuPerAudio float64 `json:"aec3_cpu_seconds_per_audio_second"`
		CorrMeanRaw     float64 `json:"corr_mean_raw"`
		CorrMeanAec3    float64 `json:"corr_mean_aec3"`
		VadRaw          int     `json:"vad_starts_raw"`
		VadAec3         int     `json:"vad_starts_aec3"`
		AsrFar          string  `json:"asr_far_reference"`
		AsrRaw          string  `json:"asr_raw_playback"`
		AsrAec3         string  `json:"asr_aec3_output"`
	}
	report := map[string]armReport{}
	for name, arm := range arms {
		// The engine consumes whole 10ms blocks; measure everything on the
		// truncated window (both sides, the shorter one bounds it) so the
		// on/off ERLE numbers cover identical data.
		frame := aec.FrameSamples(16000)
		pair := len(arm.far)
		if len(arm.near) < pair {
			pair = len(arm.near)
		}
		n := pair / frame * frame
		far, near := arm.far[:n], arm.near[:n]

		lag, strength := AlignByEnvelopeCorrelation(far, near, DefaultMaxLag)
		if lag < 0 {
			lag = 0
		}
		residual, _, _ := SubtractEcho(far, near, lag)
		subErle := ERLEdB(far, near, residual, lag)

		delayMs := lag * 1000 / audio.TargetRate
		// A fresh engine per arm: an engine reused across arms would carry
		// the previous arm's converged filter into the next reading (and
		// map iteration order would decide who benefits). Cold start is the
		// product-shaped number the baseline compares against.
		apm, err := aec.NewApm(16000)
		if err != nil {
			t.Fatalf("NewApm for arm %s: %v", name, err)
		}
		cpu0, wall0 := spikeCPUTimes(t), time.Now()
		aecOut := aec3Pass(t, apm, far, near, delayMs)
		cpu, wall := spikeCPUTimes(t).sub(cpu0), time.Since(wall0)
		_ = apm.Close()
		aecErle := ERLEdB(far, near, aecOut, lag)

		rep := armReport{
			Source: arm.source, AlignLag: lag, AlignStrength: strength,
			DelayMs: delayMs, SubtractERLEdB: subErle, Aec3ERLEdB: aecErle,
			Aec3Frames:      len(aecOut) / frame,
			Aec3CpuPercent:  cpuPercent(cpu, wall),
			Aec3CpuPerAudio: (float64(cpu.total()) / 1e9) / (float64(len(aecOut)) / audio.TargetRate),
			CorrMeanRaw:     meanCorr(WindowedCorrelation(far, near, lag, 1024, 512)),
			CorrMeanAec3:    meanCorr(WindowedCorrelation(far, aecOut, lag, 1024, 512)),
			VadRaw:          spikeVadStarts(t, vad, near),
			VadAec3:         spikeVadStarts(t, vad, aecOut),
			AsrFar:          spikeTranscribe(t, asr, far),
			AsrRaw:          spikeTranscribe(t, asr, near), AsrAec3: spikeTranscribe(t, asr, aecOut),
		}
		report[name] = rep
		_ = WriteWav16(filepath.Join(out, fmt.Sprintf("aec3-%s-out-16k.wav", name)),
			audio.TargetRate, FloatToPCM16(aecOut))
		b, err := json.MarshalIndent(rep, "", "  ")
		if err != nil {
			t.Fatalf("marshal %s: %v", name, err)
		}
		if err := os.WriteFile(filepath.Join(out, fmt.Sprintf("aec3-%s.json", name)), append(b, '\n'), 0o644); err != nil {
			t.Fatalf("write %s report: %v", name, err)
		}
		t.Logf("[%s] lag=%d strength=%.3f delay=%dms subtractERLE=%.1fdB AEC3-ERLE=%.1fdB cpu=%.2f%% corr(raw/aec3)=%.3f/%.3f vad(raw/aec3)=%d/%d",
			name, lag, strength, delayMs, subErle, aecErle, rep.Aec3CpuPercent,
			rep.CorrMeanRaw, rep.CorrMeanAec3, rep.VadRaw, rep.VadAec3)
		t.Logf("[%s] asr far=%q raw=%q aec3=%q", name, rep.AsrFar, rep.AsrRaw, rep.AsrAec3)
	}
}
