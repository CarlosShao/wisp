//go:build windows && cgo && webrtcap

package aec

// The ticket-393 reset readings, on the 392-r1 K=4 measurement shape: per
// arm, four firings of the same L1-shaped stimulus (an announcement over the
// matched echo at coupling 0.2, a -26.70 dBFS talker barging in at 1.2s),
// post-AEC readings min/max/swing, everything synthetic and deterministic so
// the cross-arm comparison is exact rather than a band.
//
// Arms:
//   - fresh: a new Apm per firing, the stimulus straight in - the cold
//     control the production default path builds (NewApm, no Reset).
//   - carry: the same Apm fed the truncated carry window first (render keeps
//     streaming 1.4s past capture end - the unpaired-render cross-window
//     carry 392-r1 measured as "thoroughly mis-adapted on the next window"),
//     then the stimulus with no reset in between: the root-cause
//     reproduction.
//   - reset: the carry window first, then Reset(), then the stimulus.
//   - rebuild-control: a new Apm, then Reset() with no history, then the
//     stimulus - what a just-reset session with no carry reads.
//
// The state-cleanliness judgment is reset == rebuild-control (same
// deterministic engine context, so the comparison is exact): after Reset the
// session is indistinguishable from a session that never carried anything.
// The fresh arm stays the production cold control; the reset arm is also
// read against its band under the domain anchor (echo residual at least 15dB
// under the -24 dBFS near echo - "the second window judges again", the
// feasibility reading the warm-posture fix needs; the resampler/detector
// half of the 392 finding is that ticket's named residual).
//
// Nothing here runs in CI: the WISP_AEC393R2=1 gate applies (the spike-gate
// shape), and the artifact JSON is written only to WISP_AEC393R2_OUT
// (default the OS temp dir).

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"testing"
)

const (
	reset393Rate       = 16000
	reset393DelayMs    = 30 // the aec3r1 rig's measured-lag stand-in (493 samples @16k)
	reset393Coupling   = 0.2
	reset393FarDbfs    = -10.0  // the announcement level; echo lands at -24 dBFS near
	reset393TalkerDbfs = -26.70 // the 392-r1 L1 control level
	reset393WindowS    = 3.4
	reset393FarEndS    = 3.0
	reset393CarryCapS  = 2.0 // capture length of the carry window (render runs the full 3.4s)
	reset393TalkerAtS  = 1.2
	reset393EchoWin0S  = 0.8 // echo-only reading window, post cold-start convergence
	reset393EchoWin1S  = 1.15
	reset393Episode1S  = 1.5 // episode reading window: talker start .. +0.3s
	reset393Firings    = 4
	reset393Jitter     = 32 // echo-delay jitter per firing, samples (2ms steps)
)

// reset393Voice synthesizes a deterministic speech-band waveform: four
// per-tone amplitude-modulated carriers, peak-normalized to peakDbfs. The
// two seeds give the announcement and the talker disjoint tone sets and
// phases so neither masquerades as delayed copies of the other.
func reset393Voice(seed int, seconds float64, peakDbfs float64) []float32 {
	n := int(seconds * float64(reset393Rate))
	buf := make([]float64, n)
	freqs := [4]float64{300, 700, 1300, 2200}
	if seed == 1 {
		freqs = [4]float64{500, 900, 1700, 2600}
	}
	for i := range buf {
		t := float64(i) / float64(reset393Rate)
		v := 0.0
		for k, f := range freqs {
			am := 0.6 + 0.4*math.Sin(2*math.Pi*(2.3+1.1*float64(k))*t+float64(seed)*0.7+float64(k))
			v += am * math.Sin(2*math.Pi*f*t+float64(seed)*1.3+float64(k)*0.9)
		}
		buf[i] = v
	}
	peak := 0.0
	for _, v := range buf {
		if a := math.Abs(v); a > peak {
			peak = a
		}
	}
	gain := math.Pow(10, peakDbfs/20) / peak
	out := make([]float32, n)
	for i, v := range buf {
		out[i] = float32(v * gain)
	}
	return out
}

// reset393Near builds the near (capture) signal: the far announcement delayed
// by delaySamples at the coupling, plus the talker from its barge-in point
// onward. delaySamples is the true acoustic delay; ProcessCapture is told the
// matching round of delayMs.
func reset393Near(far, talker []float32, delaySamples int) []float32 {
	near := make([]float32, len(far))
	for i := range near {
		if i >= delaySamples {
			near[i] = float32(reset393Coupling) * far[i-delaySamples]
		}
		if i < len(talker) {
			near[i] += talker[i]
		}
	}
	return near
}

// reset393RunPair feeds one window in lockstep (render frame then capture
// frame, the APM's pinned order). The capture stream may be shorter than the
// render stream - the shorter leg just stops, which is exactly the
// window-end truncation carry.
func reset393RunPair(t *testing.T, apm *Apm, far, near []float32, delayMs int) []float32 {
	t.Helper()
	frame := FrameSamples(reset393Rate)
	nCap := len(near) / frame
	out := make([]float32, 0, nCap*frame)
	for off := 0; off+frame <= len(far); off += frame {
		if err := apm.ProcessRender(far[off : off+frame]); err != nil {
			t.Fatalf("render frame at %d: %v", off, err)
		}
		if off+frame <= nCap*frame {
			got, err := apm.ProcessCapture(near[off:off+frame], delayMs)
			if err != nil {
				t.Fatalf("capture frame at %d: %v", off, err)
			}
			out = append(out, got...)
		}
	}
	return out
}

func reset393RmsDbfs(x []float32, from, to int) float64 {
	if to > len(x) {
		to = len(x)
	}
	sum := 0.0
	for _, v := range x[from:to] {
		sum += float64(v) * float64(v)
	}
	rms := math.Sqrt(sum / float64(to-from))
	if rms < 1e-10 {
		rms = 1e-10
	}
	return 20 * math.Log10(rms)
}

func reset393MaxDiff(a, b []float32) float64 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	m := 0.0
	for i := 0; i < n; i++ {
		if d := math.Abs(float64(a[i]) - float64(b[i])); d > m {
			m = d
		}
	}
	return m
}

type reset393Reading struct {
	Firing      int     `json:"firing"`
	DelayMs     int     `json:"delay_ms"`
	EchoResidDb float64 `json:"echo_resid_dbfs"`
	EpisodeDb   float64 `json:"episode_rms_dbfs"`
}

type reset393Arm struct {
	Name                    string            `json:"name"`
	Readings                []reset393Reading `json:"readings"`
	MinEcho                 float64           `json:"min_echo_dbfs"`
	MaxEcho                 float64           `json:"max_echo_dbfs"`
	SwingEcho               float64           `json:"swing_echo_db"`
	MinEpisode              float64           `json:"min_episode_dbfs"`
	MaxEpisode              float64           `json:"max_episode_dbfs"`
	SwingEpi                float64           `json:"swing_episode_db"`
	MaxDiffVsFresh          float64           `json:"max_diff_vs_fresh"`
	MaxDiffVsRebuildControl float64           `json:"max_diff_vs_rebuild_control"`
}

func reset393Summarize(arm *reset393Arm) {
	first := true
	for _, r := range arm.Readings {
		if first {
			arm.MinEcho, arm.MaxEcho, arm.MinEpisode, arm.MaxEpisode = r.EchoResidDb, r.EchoResidDb, r.EpisodeDb, r.EpisodeDb
			first = false
			continue
		}
		arm.MinEcho = min(arm.MinEcho, r.EchoResidDb)
		arm.MaxEcho = max(arm.MaxEcho, r.EchoResidDb)
		arm.MinEpisode = min(arm.MinEpisode, r.EpisodeDb)
		arm.MaxEpisode = max(arm.MaxEpisode, r.EpisodeDb)
	}
	arm.SwingEcho = arm.MaxEcho - arm.MinEcho
	arm.SwingEpi = arm.MaxEpisode - arm.MinEpisode
}

// Test393R2ResetRestoresCleanState is the AC#2 state-cleanliness reading and
// the AC#3 reverse judgments on the K=4 shape: fresh vs carry vs reset vs
// rebuild-control, the reset arm asserted exactly back on the rebuild
// control, the carry arm asserted degraded, and a restore core re-reading
// the fresh shape after everything.
func Test393R2ResetRestoresCleanState(t *testing.T) {
	if os.Getenv("WISP_AEC393R2") != "1" {
		t.Skip("reset reading: set WISP_AEC393R2=1 (synthetic, deterministic, no audio hardware)")
	}
	if _, err := NewApm(reset393Rate); err != nil {
		t.Skipf("AEC3 binding not in this build: %v", err)
	}

	far := reset393Voice(0, reset393WindowS, reset393FarDbfs)
	farEnd := int(reset393FarEndS * float64(reset393Rate))
	for i := farEnd; i < len(far); i++ {
		far[i] = 0
	}
	talker := reset393Voice(1, reset393WindowS-reset393TalkerAtS, reset393TalkerDbfs)

	arms := map[string]*reset393Arm{
		"fresh":           {Name: "fresh"},
		"carry":           {Name: "carry"},
		"reset":           {Name: "reset"},
		"rebuild-control": {Name: "rebuild-control"},
	}
	freshOut := make([][]float32, reset393Firings)
	resetOut := make([][]float32, reset393Firings)

	for k := 1; k <= reset393Firings; k++ {
		delaySamples := 480 + reset393Jitter*(k-1)
		delayMs := reset393DelayMs + 2*(k-1)
		near := reset393Near(far, talker, delaySamples)
		carryNear := reset393Near(far, make([]float32, 0), delaySamples)
		carryCap := carryNear[:int(reset393CarryCapS*float64(reset393Rate))]

		// fresh: the cold control, the production default path shape.
		apm, err := NewApm(reset393Rate)
		if err != nil {
			t.Fatalf("fresh firing %d: %v", k, err)
		}
		outFresh := reset393RunPair(t, apm, far, near, delayMs)
		apm.Close()
		freshOut[k-1] = outFresh

		// carry: the truncated window first, the stimulus right after, no
		// reset - the 392 root-cause shape.
		apm, err = NewApm(reset393Rate)
		if err != nil {
			t.Fatalf("carry firing %d: %v", k, err)
		}
		reset393RunPair(t, apm, far, carryCap, delayMs)
		outCarry := reset393RunPair(t, apm, far, near, delayMs)
		apm.Close()

		// reset: the same carry, then the face, then the stimulus.
		apm, err = NewApm(reset393Rate)
		if err != nil {
			t.Fatalf("reset firing %d: %v", k, err)
		}
		reset393RunPair(t, apm, far, carryCap, delayMs)
		if err := apm.Reset(); err != nil {
			t.Fatalf("reset firing %d: Reset: %v", k, err)
		}
		outReset := reset393RunPair(t, apm, far, near, delayMs)
		apm.Close()
		resetOut[k-1] = outReset

		// rebuild-control: Reset with no history - the just-reset shape.
		apm, err = NewApm(reset393Rate)
		if err != nil {
			t.Fatalf("rebuild-control firing %d: %v", k, err)
		}
		if err := apm.Reset(); err != nil {
			t.Fatalf("rebuild-control firing %d: Reset: %v", k, err)
		}
		outRebuild := reset393RunPair(t, apm, far, near, delayMs)
		apm.Close()

		echo0, echo1 := int(reset393EchoWin0S*float64(reset393Rate)), int(reset393EchoWin1S*float64(reset393Rate))
		ep0, ep1 := int(reset393TalkerAtS*float64(reset393Rate)), int(reset393Episode1S*float64(reset393Rate))
		arms["fresh"].Readings = append(arms["fresh"].Readings, reset393Reading{
			k, delayMs,
			reset393RmsDbfs(outFresh, echo0, echo1), reset393RmsDbfs(outFresh, ep0, ep1),
		})
		arms["carry"].Readings = append(arms["carry"].Readings, reset393Reading{
			k, delayMs,
			reset393RmsDbfs(outCarry, echo0, echo1), reset393RmsDbfs(outCarry, ep0, ep1),
		})
		arms["reset"].Readings = append(arms["reset"].Readings, reset393Reading{
			k, delayMs,
			reset393RmsDbfs(outReset, echo0, echo1), reset393RmsDbfs(outReset, ep0, ep1),
		})
		arms["rebuild-control"].Readings = append(arms["rebuild-control"].Readings, reset393Reading{
			k, delayMs,
			reset393RmsDbfs(outRebuild, echo0, echo1), reset393RmsDbfs(outRebuild, ep0, ep1),
		})
		arms["carry"].MaxDiffVsFresh = math.Max(arms["carry"].MaxDiffVsFresh, reset393MaxDiff(outCarry, outFresh))
		arms["reset"].MaxDiffVsFresh = math.Max(arms["reset"].MaxDiffVsFresh, reset393MaxDiff(outReset, outFresh))
		arms["reset"].MaxDiffVsRebuildControl = math.Max(arms["reset"].MaxDiffVsRebuildControl, reset393MaxDiff(outReset, outRebuild))
	}
	for _, a := range arms {
		reset393Summarize(a)
	}

	for _, name := range []string{"fresh", "carry", "reset", "rebuild-control"} {
		a := arms[name]
		for _, r := range a.Readings {
			t.Logf("393R2-%s firing=%d delay=%dms echo_resid=%.2f episode_rms=%.2f",
				a.Name, r.Firing, r.DelayMs, r.EchoResidDb, r.EpisodeDb)
		}
		t.Logf("393R2-%s echo swing=%.2f (%.2f..%.2f) episode swing=%.2f (%.2f..%.2f) max_diff_vs_fresh=%.3g max_diff_vs_rebuild_control=%.3g",
			a.Name, a.SwingEcho, a.MinEcho, a.MaxEcho, a.SwingEpi, a.MinEpisode, a.MaxEpisode,
			a.MaxDiffVsFresh, a.MaxDiffVsRebuildControl)
	}

	// AC#2, the state-cleanliness judgment: after Reset the session is
	// byte-identical to a just-reset session with no carry history.
	if d := arms["reset"].MaxDiffVsRebuildControl; d > 1e-6 {
		t.Fatalf("AC#2 FAILED: the reset arm's frame-level diff against the rebuild control is %.3g (want <= 1e-6) - Reset did not clear the carried state", d)
	}
	t.Logf("393R2-FINDING reset==rebuild-control holds bit-exact on all %d firings (max frame diff %.3g): the carried state is fully cleared by Reset", reset393Firings, arms["reset"].MaxDiffVsRebuildControl)

	// The domain anchor on the reset arm: the echo residual sits at least
	// 15dB under the -24 dBFS near echo - the "second window judges again"
	// shape (the 392 finding was episodes=0 with the residual riding at the
	// near level).
	for k := 0; k < reset393Firings; k++ {
		if r := arms["reset"].Readings[k].EchoResidDb; r > -39 {
			t.Fatalf("AC#2 FAILED: reset firing %d echo residual %.2f dBFS is not cancelled (want <= -39, i.e. 15dB under the near echo): the second window would not judge", k+1, r)
		}
	}
	// The fresh band comparison, logged as context: the reset context (the
	// state-synced split-band shape) and the production default path (the
	// full-band shape the stale submodule_states leave) read within ~2dB of
	// each other on this stimulus - named in the evidence file, not a
	// pass/fail here.
	t.Logf("393R2-FINDING context readings: fresh %.2f..%.2f vs reset %.2f..%.2f (frame diff %.3g - the full-band vs split-band engine shape, see the evidence file)",
		arms["fresh"].MinEcho, arms["fresh"].MaxEcho, arms["reset"].MinEcho, arms["reset"].MaxEcho, arms["reset"].MaxDiffVsFresh)

	// AC#3, the reverse judgment: the carry arm is measurably degraded
	// against the reset arm - the cross-window carry reproduces and the
	// reset face removes it.
	carryMean, resetMean := 0.0, 0.0
	for k := 0; k < reset393Firings; k++ {
		carryMean += arms["carry"].Readings[k].EchoResidDb
		resetMean += arms["reset"].Readings[k].EchoResidDb
	}
	carryMean /= reset393Firings
	resetMean /= reset393Firings
	if carryMean < resetMean+3 {
		t.Fatalf("AC#3 FAILED: the carry arm is not degraded (echo_resid mean %.2f vs reset %.2f) - this stimulus does not reproduce the cross-window carry; the reset reading would prove nothing",
			carryMean, resetMean)
	}
	t.Logf("393R2-FINDING carry degraded by %.2f dB against the reset arm (mean %.2f vs %.2f): the cross-window carry reproduces at binding isolation and Reset removes it",
		carryMean-resetMean, carryMean, resetMean)

	// The restore core: after everything, a brand-new session re-reads the
	// fresh[0] shape (same delay, deterministic comparison).
	apm, err := NewApm(reset393Rate)
	if err != nil {
		t.Fatalf("restore core: %v", err)
	}
	outRestore := reset393RunPair(t, apm, far, reset393Near(far, talker, 480), reset393DelayMs)
	apm.Close()
	if d := reset393MaxDiff(outRestore, freshOut[0]); d > 1e-6 {
		t.Fatalf("the restore core drifted (frame diff %.3g against fresh firing 1): the instrument chain did not hold the run", d)
	}
	t.Logf("393R2-RESTORE core matches fresh firing 1 (frame diff %.3g)", reset393MaxDiff(outRestore, freshOut[0]))

	reset393WriteJSON(t, "393-aec-reset.json", map[string]any{
		"arms":         []*reset393Arm{arms["fresh"], arms["carry"], arms["reset"], arms["rebuild-control"]},
		"stimulus":     "synthetic: announcement -10 dBFS far over matched echo 0.2 (near echo -24 dBFS), talker -26.70 dBFS at 1.2s, delay 30ms + 2ms/firing jitter, K=4",
		"carry":        "render runs the full 3.4s window, capture stops at 2.0s, the next window follows in the same Apm",
		"finding":      "Reset() (the public Initialize() face) restores the just-reset state bit-exactly (reset==rebuild-control on all firings) and returns the echo residual to cancelled depth; the no-reset carry arm degrades ~8dB - the 392 cross-window carry reproduces and the reset face removes it at the APM half",
		"context_note": "the fresh (production NewApm) engine shape and the reset (Initialize-locked) engine shape differ at frame level (full-band vs split-band capture, a stale submodule_states property of the ApplyConfig-after-Create path in this pinned tree) while reading within ~2dB - named in docs/evidence/s1/393r2-flush-reset.md, untouched by this leg",
	})
}

// Test393R2ResetPreservesTunedEngine rides the 392 tuning face through the
// reset: a tuned session carried over the truncated window, Reset, then the
// stimulus, must read identical to a fresh tuned session - the injected
// EchoControlFactory survives Initialize() and re-creates the same tuned
// engine (the property the warm-posture fix needs for tuned engines). The
// knob effect is read in the same engine context on purpose (reset vs
// rebuild-control): the fresh default path runs the other engine shape, so a
// fresh-vs-fresh comparison would confound the knob with the context
// difference named in the evidence file.
func Test393R2ResetPreservesTunedEngine(t *testing.T) {
	if os.Getenv("WISP_AEC393R2") != "1" {
		t.Skip("reset reading: set WISP_AEC393R2=1 (synthetic, deterministic, no audio hardware)")
	}
	if _, err := NewApmTuned(reset393Rate, Tuning{}); err != nil {
		t.Skipf("AEC3 binding not in this build: %v", err)
	}
	tuning := Tuning{DominantNearendEnrThreshold: 0.45}
	far := reset393Voice(0, reset393WindowS, reset393FarDbfs)
	farEnd := int(reset393FarEndS * float64(reset393Rate))
	for i := farEnd; i < len(far); i++ {
		far[i] = 0
	}
	talker := reset393Voice(1, reset393WindowS-reset393TalkerAtS, reset393TalkerDbfs)
	near := reset393Near(far, talker, 480)
	carryNear := reset393Near(far, make([]float32, 0), 480)
	carryCap := carryNear[:int(reset393CarryCapS*float64(reset393Rate))]

	tunedFreshOut := func() []float32 {
		apm, err := NewApmTuned(reset393Rate, tuning)
		if err != nil {
			t.Fatalf("fresh tuned session: %v", err)
		}
		defer apm.Close()
		return reset393RunPair(t, apm, far, near, reset393DelayMs)
	}()
	tunedResetOut := func() []float32 {
		apm, err := NewApmTuned(reset393Rate, tuning)
		if err != nil {
			t.Fatalf("tuned session: %v", err)
		}
		reset393RunPair(t, apm, far, carryCap, reset393DelayMs)
		if err := apm.Reset(); err != nil {
			t.Fatalf("tuned session: Reset: %v", err)
		}
		defer apm.Close()
		return reset393RunPair(t, apm, far, near, reset393DelayMs)
	}()
	dTuned := reset393MaxDiff(tunedResetOut, tunedFreshOut)
	if dTuned > 1e-6 {
		t.Fatalf("the tuned engine did not survive Reset (frame diff %.3g against the fresh tuned session): the injected factory was not retained", dTuned)
	}

	defaultResetOut := func() []float32 {
		apm, err := NewApm(reset393Rate)
		if err != nil {
			t.Fatalf("default session: %v", err)
		}
		if err := apm.Reset(); err != nil {
			t.Fatalf("default session: Reset: %v", err)
		}
		defer apm.Close()
		return reset393RunPair(t, apm, far, near, reset393DelayMs)
	}()
	dKnobs := reset393MaxDiff(defaultResetOut, tunedResetOut)
	echo0, echo1 := int(reset393EchoWin0S*float64(reset393Rate)), int(reset393EchoWin1S*float64(reset393Rate))
	t.Logf("393R2-TUNED reset==fresh-tuned holds (frame diff %.3g); knob effect in the same engine context: frame diff %.3g, echo_resid tuned %.2f vs default %.2f",
		dTuned, dKnobs, reset393RmsDbfs(tunedFreshOut, echo0, echo1), reset393RmsDbfs(defaultResetOut, echo0, echo1))

	reset393WriteJSON(t, "393-aec-reset-tuned.json", map[string]any{
		"tuning":               tuning,
		"reset_vs_fresh_diff":  dTuned,
		"knob_effect_same_ctx": dKnobs,
	})
}

func reset393WriteJSON(t *testing.T, name string, v any) {
	t.Helper()
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatalf("artifact %s: %v", name, err)
	}
	dir := os.Getenv("WISP_AEC393R2_OUT")
	if dir == "" {
		dir = os.TempDir()
	}
	if err := os.WriteFile(filepath.Join(dir, name), append(b, '\n'), 0o644); err != nil {
		t.Fatalf("artifact %s: %v", name, err)
	}
}
