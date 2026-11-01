//go:build windows && cgo && webrtcap

package aecspike

// The ticket-392-r2 keepalive isolation and revival readings, on the
// deterministic replica of the monitor's announcement-window judge chain
// (cmd/wisp's bargeinDrain shape: anchor, envelope-aligned lag, APM3
// render/capture lockstep at the measured delay, the barge-in gate at the
// tightened cut, the residual noise gate, silero VAD episode starts judged
// through Evaluate). The 392-r1 measurements ran on the real staged rig in
// cmd/wisp; this leg's write surface is internal/audio only, so the same
// measurement shape is replicated here with synthetic deterministic
// scheduling (393-r2's precedent: the readings compare across arms, so the
// arms must be bit-reproducible, not wall-clock-real).
//
// The 392-r1 finding this leg isolates: under the keepalive posture (the
// silero session and the streaming resampler held across announcement
// windows, the APM rebuilt per window) the SECOND warm window judges nothing
// - episodes=0, triggers=0 - while the cold posture stays healthy on the
// same sequence. Two warm pieces exist, so the isolation matrix holds one
// warm and one fresh at a time:
//
//	cold-control   everything rebuilt per window (the production default)
//	warm-392       detector + resampler warm, APM fresh (the 392-r1 shape)
//	iso-det-fresh  resampler warm ONLY (detector fresh per window)
//	iso-res-fresh  detector warm ONLY (resampler fresh per window)
//
// The discriminator is two-fold per window: the episode/trigger counts read
// the detector half (an endpointer state machine cannot fire a start while
// an utterance is open), and the echo-only residual depth reads the
// resampler half (the far tap is the warm resampler's output; a corrupted
// far reference cancels to nothing). The depth judgments are relative to
// the cold control's same-window reading - the APM is fresh in both, so the
// convergence head cancels out of the comparison - and to the near echo
// itself (at least 8dB of real cancellation).
//
// The revival arm is the 393-v1 "Reset as the unifier" shape applied to the
// whole warm pool: every held piece gets its public between-windows clear
// face - the APM gets Reset() (393-r2: bit-exact state restore), the
// resampler gets Flush() (its documented stream-state reset), the detector
// gets Finalize(nil) (closes the open utterance the abort left open; the
// endpointer returns to idle with the silero session still warm):
//
//	revival             det warm + Finalize, resampler warm + Flush, APM warm + Reset
//	revival-cool        revival plus the detector's silence cool-down (the RNN
//	                    carries the abort's speech context across the idle
//	                    gap; the probe measures a hangover start firing on
//	                    post-speech silence, and the first warm window's
//	                    trigger needs the decay)
//	revival-no-reset    revival minus the APM Reset (the carried-APM reverse)
//	revival-no-finalize revival minus the detector Finalize (the 392 finding returns)
//
// Windows are W1 (engines built), W2 (the discriminator), W3 (persistence).
// The whole matrix runs twice and must read identically (deterministic
// synthesis, deterministic engines). Nothing here runs in CI: the
// WISP_AEC392R2=1 gate applies, and the artifact JSON is written only to
// WISP_AEC392R2_OUT.

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/audio/aec"
	"github.com/CarlosShao/wisp/internal/speech"
)

// The monitor's shape constants, mirrored (cmd/wisp
// resident_voice_bargein_windows.go): all audio-clock quantities, the same
// values the staged readings ran under.
const (
	keep392AlignFar      = 4800 // 300ms of far head the delay measurement runs on
	keep392AlignMaxLag   = 4800 // the search span around the anchor
	keep392AlignStrength = 0.2  // the trust floor; below it the lag answers zero
	keep392Frame         = 512  // the C8 seam frame (32ms at the 16k seam rate)
	keep392ApmFrame      = 160  // the 10ms APM frame
	keep392OnsetRing     = 2 * 16000
	keep392FloorWindows  = 64 // the residual noise gate: 2s of 32ms windows
	keep392GateStep      = 2.0
	keep392GateHold      = 16 // 512ms

	// The (a) material: the 26p echo return lag and coupling, the talker at
	// the -20 dBFS near-mic level starting 400ms into the announcement - the
	// exact (a) shape the 392-r1 keepalive readings judged. The 400ms
	// pre-talker span is load-bearing for the isolation: it is shorter than
	// the endpointer's 700ms close threshold, so the utterance the abort
	// leaves open survives into the next window.
	keep392EchoDelay  = 480
	keep392EchoGain   = 0.2
	keep392TalkerAt   = 6400 // 400ms at the seam rate
	keep392TalkerDbfs = -20.0

	// The 393-r2 far level (peak -10 dBFS): the matched echo lands at -24
	// dBFS peak, the same ruler the 393-r2 domain anchor was cut on.
	keep392FarDbfs = -10.0
)

const (
	keep392AnnounceText = "任务已经完成，明天下午三点的会议改到四点，请准时参加，谢谢。"
	keep392TalkerText   = "等一下，我有问题要说。"
)

// keep392EchoOnlyLate is the late echo-only reading span [200ms, 400ms) of
// processed residual: inside the echo-only span (the talker starts at
// 400ms), past the first convergence instants, and fully covered by every
// window (a trigger cannot fire before ~700ms of audio has processed - the
// talker needs MinSpeech through the gate first). The APM is fresh in every
// arm's W2, so the span's convergence head is common to all arms and the
// cross-arm comparison reads the far path, not the head.
func keep392EchoOnlyLate() (int, int) { return 3200, keep392TalkerAt }

// keep392MinCancelDb is the "real cancellation happened" margin: the
// echo-only residual must sit this far under the near echo itself.
const keep392MinCancelDb = 8

// keep392Materials holds the window-independent stimulus: the announcement
// at the TTS engine rate in its real synthesis chunking (the far tap's
// input), the 16k talker, and the seam-rate near stream (matched echo +
// talker) built from the ONE-SHOT resample of the same announcement - the
// echo path is independent of the far tap, exactly as the played audio is
// independent of the monitor's resampled reference. Real matcha speech:
// silero must hear speech, and the matcha render is deterministic.
type keep392Materials struct {
	engineRate     int
	announceChunks [][]float32 // engine-rate chunks, peak-normalized to keep392FarDbfs
	near           []float32   // seam-rate capture: echo (delay 480, gain 0.2) + talker
	talker         []float32   // the seam-rate talker alone, -20 dBFS RMS
}

func keep392TtsConfig(t *testing.T) speech.TtsConfig {
	t.Helper()
	dir := spikeModelsDir(t)
	m := filepath.Join(dir, "matcha-icefall-zh-baker")
	return speech.TtsConfig{
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
}

func keep392Speak(t *testing.T, cfg speech.TtsConfig, text string, wantChunks bool) ([][]float32, []float32, int) {
	t.Helper()
	eng, err := speech.NewMatchaTts(cfg)
	if err != nil {
		t.Skipf("matcha tts unavailable: %v", err)
	}
	defer func() { _ = eng.Close() }()
	var all []float32
	var chunks [][]float32
	if err := eng.Speak(text, func(samples []float32) error {
		all = append(all, samples...)
		if wantChunks && len(samples) > 0 {
			chunks = append(chunks, append([]float32(nil), samples...))
		}
		return nil
	}); err != nil {
		t.Fatalf("tts speak %q: %v", text, err)
	}
	if len(all) == 0 {
		t.Fatalf("tts speak %q produced no samples", text)
	}
	return chunks, all, eng.SampleRate()
}

func keep392MaterialsBuild(t *testing.T) *keep392Materials {
	t.Helper()
	cfg := keep392TtsConfig(t)
	chunks, all, rate := keep392Speak(t, cfg, keep392AnnounceText, true)
	peak := 0.0
	for _, v := range all {
		if a := math.Abs(float64(v)); a > peak {
			peak = a
		}
	}
	gain := math.Pow(10, keep392FarDbfs/20) / peak
	for i := range all {
		all[i] *= float32(gain)
	}
	for _, c := range chunks {
		for i := range c {
			c[i] *= float32(gain)
		}
	}

	// The seam-rate far the echo is built from: the one-shot resample of the
	// normalized announcement, the independent path.
	far16k := PCM16ToFloat(audio.ResampleLinear(audio.FloatToPCM16(all), rate, audio.TargetRate))
	near := make([]float32, len(far16k)+keep392EchoDelay)
	for i, v := range far16k {
		near[i+keep392EchoDelay] = keep392EchoGain * v
	}

	_, talkerRaw, talkerRate := keep392Speak(t, cfg, keep392TalkerText, false)
	talker := PCM16ToFloat(audio.ResampleLinear(audio.FloatToPCM16(talkerRaw), talkerRate, audio.TargetRate))
	var sum float64
	for _, v := range talker {
		sum += float64(v) * float64(v)
	}
	tGain := math.Pow(10, keep392TalkerDbfs/20) / math.Sqrt(sum/float64(len(talker)))
	for i := range talker {
		talker[i] *= float32(tGain)
		near[keep392TalkerAt+i] += talker[i]
	}
	return &keep392Materials{engineRate: rate, announceChunks: chunks, near: near, talker: talker}
}

// keep392Posture is one arm's engine policy. The names are the readings'
// keys; the booleans mirror exactly which pieces ride the keepalive and
// which between-windows clear face each held piece gets.
type keep392Posture struct {
	name string

	detWarm     bool // the silero detector held across windows
	detFinalize bool // ...with Finalize(nil) at the boundary (drops the open utterance)
	detCooldown bool // ...plus a silence cool-down through the session (decays the RNN context the abort left)
	resWarm     bool // the streaming resampler held across windows
	resFlush    bool // ...with Flush() at the boundary
	apmWarm     bool // the APM held across windows
	apmReset    bool // ...with Reset() at the boundary
}

func keep392Arms() []keep392Posture {
	return []keep392Posture{
		{name: "cold-control"},
		{name: "warm-392", detWarm: true, resWarm: true},
		{name: "iso-det-fresh", resWarm: true},
		{name: "iso-res-fresh", detWarm: true},
		{name: "revival", detWarm: true, detFinalize: true, resWarm: true, resFlush: true, apmWarm: true, apmReset: true},
		{name: "revival-cool", detWarm: true, detFinalize: true, detCooldown: true, resWarm: true, resFlush: true, apmWarm: true, apmReset: true},
		{name: "revival-no-reset", detWarm: true, detFinalize: true, resWarm: true, resFlush: true, apmWarm: true},
		{name: "revival-no-finalize", detWarm: true, resWarm: true, resFlush: true, apmWarm: true, apmReset: true},
	}
}

func (p keep392Posture) policy() map[string]bool {
	return map[string]bool{
		"det_warm": p.detWarm, "det_finalize": p.detFinalize, "det_cooldown": p.detCooldown,
		"res_warm": p.resWarm, "res_flush": p.resFlush,
		"apm_warm": p.apmWarm, "apm_reset": p.apmReset,
	}
}

// keep392CooldownSamples is the detector cool-down: 1.2s of digital silence
// through the session at the window boundary. The silero RNN carries the
// abort's speech context across an idle gap (no input, no decay), and the
// probe measures a start firing on post-speech silence - the hangover feeds
// the endpointer's speech run. The cool-down lets the context decay (and,
// if a hangover start opens an utterance mid-cool-down, the remaining
// silence exceeds the 700ms close threshold, so the boundary ends idle).
const keep392CooldownSamples = 19200

type keep392WindowReading struct {
	Window          int       `json:"window"`
	Episodes        int       `json:"episodes"`
	Triggers        int       `json:"triggers"`
	TriggerAtSample int       `json:"trigger_at_sample"` // -1 = none fired
	Lag             int       `json:"lag"`
	LagStrength     float64   `json:"lag_strength"`
	ReadySamples    int       `json:"ready_samples"`
	Processed       int       `json:"processed"`
	NearEchoDb      float64   `json:"near_echo_dbfs"`
	EchoResidDb     float64   `json:"echo_only_late_resid_dbfs"`
	NearProfile     []float64 `json:"near_profile_dbfs"`
	ResProfile      []float64 `json:"res_profile_dbfs"`
}

type keep392ArmReading struct {
	Name    string                 `json:"name"`
	Posture map[string]bool        `json:"posture"`
	Windows []keep392WindowReading `json:"windows"`
}

// keep392Engines is one arm's engine pool across its windows.
type keep392Engines struct {
	posture keep392Posture
	det     *speech.Detector
	apm     *aec.Apm
	rs      *audio.Resampler
}

func keep392NewDetector(t *testing.T) *speech.Detector {
	t.Helper()
	det, err := speech.NewDetector(speech.DetectorConfig{
		ModelPath: filepath.Join(spikeModelsDir(t), "silero_vad.onnx"),
	})
	if err != nil {
		t.Skipf("silero vad unavailable: %v", err)
	}
	return det
}

// keep392ArmClose releases whatever the posture still holds (the arm's tail,
// the C11 discipline: sessions are released, never leaked).
func keep392ArmClose(eng *keep392Engines) {
	if eng.apm != nil {
		_ = eng.apm.Close()
		eng.apm = nil
	}
	if eng.det != nil {
		_ = eng.det.Close()
		eng.det = nil
	}
	eng.rs = nil
}

// keep392Judge is the window's drain state, the bargeinDrain mirror: one
// owner (the window loop), no locks.
type keep392Judge struct {
	mat          *keep392Materials
	eng          *keep392Engines
	far          []float32
	anchorSet    bool
	anchor       int
	capCount     int
	pending      []float32
	procSamples  int
	lag          int
	lagStrength  float64
	gateArmed    bool
	readySamples int
	gate         *audio.BargeinGate
	vadBuf       []float32
	onsetRing    []float32
	eps          []speech.Endpoint
	triggered    bool
	triggerAt    int
	episodes     int
	floorRing    []float64
	floorN       int
	gateHold     int
	nearRec      []float32
	resRec       []float32
}

// arm builds or reuses the window's engines (the monitor's arm mirror: what
// the keepalive holds stays, what it does not hold is rebuilt).
func (j *keep392Judge) arm(t *testing.T) {
	p := j.eng.posture
	if !p.detWarm || j.eng.det == nil {
		if j.eng.det != nil {
			_ = j.eng.det.Close()
		}
		j.eng.det = keep392NewDetector(t)
	}
	if !p.apmWarm || j.eng.apm == nil {
		if j.eng.apm != nil {
			_ = j.eng.apm.Close()
		}
		a, aerr := aec.NewApm(audio.TargetRate)
		if aerr != nil {
			t.Skipf("AEC3 binding not in this build: %v", aerr)
		}
		j.eng.apm = a
	}
	if !p.resWarm || j.eng.rs == nil {
		j.eng.rs = audio.NewResampler(j.mat.engineRate, audio.TargetRate)
	}
}

// boundary is the release mirror: the window closed, the held pieces get
// their between-windows clear face, the non-held pieces go down.
func (j *keep392Judge) boundary() {
	p := j.eng.posture
	if j.eng.apm != nil {
		if p.apmWarm {
			if p.apmReset {
				_ = j.eng.apm.Reset()
			}
		} else {
			_ = j.eng.apm.Close()
			j.eng.apm = nil
		}
	}
	if j.eng.rs != nil {
		if p.resWarm {
			if p.resFlush {
				j.eng.rs.Flush()
			}
		} else {
			j.eng.rs = nil
		}
	}
	if j.eng.det != nil {
		if p.detWarm {
			if p.detFinalize {
				// Drop the open utterance the abort left behind: the
				// endpointer returns to idle and the silero session stays
				// warm. This is the detector half's clear face.
				_ = j.eng.det.Finalize(nil)
			}
			if p.detCooldown {
				// Decay the RNN context the abort left (the session sat
				// idle holding it): feed digital silence through the
				// session and discard whatever transitions the hangover
				// produces - boundary events are not window readings.
				silence := make([]float32, speech.WindowSamples)
				var eps []speech.Endpoint
				for n := 0; n+speech.WindowSamples <= keep392CooldownSamples; n += speech.WindowSamples {
					eps = j.eng.det.Accept(silence, eps[:0])
				}
			}
		} else {
			_ = j.eng.det.Close()
			j.eng.det = nil
		}
	}
}

// feedFar runs the far tap: the announcement's synthesis chunks through the
// (warm or fresh) streaming resampler into the window's far ring - the
// speakFar mirror.
func (j *keep392Judge) feedFar() {
	for _, chunk := range j.mat.announceChunks {
		conv := j.eng.rs.Process(audio.FloatToPCM16(chunk))
		for _, v := range conv {
			j.far = append(j.far, float32(v)/32768)
		}
	}
}

// armJudge is the drain's judge arming mirror: wait for the alignment spans,
// measure the lag online, arm the gate at the tightened cut.
func (j *keep392Judge) armJudge() bool {
	if len(j.far) < keep392AlignFar {
		return false
	}
	if len(j.pending) < keep392AlignFar+keep392AlignMaxLag {
		return false
	}
	farHead := j.far[:keep392AlignFar]
	var farSum float64
	for _, v := range farHead {
		farSum += float64(v) * float64(v)
	}
	farRMS := math.Sqrt(farSum / float64(keep392AlignFar))
	measured, strength := audio.AlignByEnvelopeCorrelation(farHead, j.pending[:keep392AlignFar+keep392AlignMaxLag], keep392AlignMaxLag)
	if measured < 0 || strength < keep392AlignStrength {
		measured = 0
	}
	j.lag = measured
	j.lagStrength = strength
	j.gateArmed = true
	j.readySamples = j.capCount - j.anchor
	j.gate = audio.NewBargeinGate(audio.BargeinOptions{
		EnergyThresholdDbfs: audio.BargeinTightenedEnergyDbfs,
		TightenedEnergyDbfs: audio.BargeinTightenedEnergyDbfs,
		DynamicGating:       true,
		LagSamples:          j.lag,
		FarRMS:              farRMS,
	})
	return true
}

// process pushes whole 10ms blocks through the judge while the far ring
// covers the next block and pending capture holds one - the drain's process
// mirror, same order, same delay hint convention.
func (j *keep392Judge) process() {
	for {
		if len(j.far) < j.procSamples+keep392ApmFrame {
			return
		}
		if len(j.pending) < keep392ApmFrame {
			return
		}
		farBlock := j.far[j.procSamples : j.procSamples+keep392ApmFrame]
		nearBlock := j.pending[:keep392ApmFrame]
		if err := j.eng.apm.ProcessRender(farBlock); err != nil {
			return
		}
		res, err := j.eng.apm.ProcessCapture(nearBlock, j.lag*1000/audio.TargetRate)
		if err != nil {
			return
		}
		j.pending = j.pending[keep392ApmFrame:]
		j.procSamples += keep392ApmFrame
		j.gate.Feed(farBlock, nearBlock, res)
		j.nearRec = append(j.nearRec, nearBlock...)
		j.resRec = append(j.resRec, res...)
		j.onsetRing = append(j.onsetRing, res...)
		if len(j.onsetRing) > keep392OnsetRing {
			cut := len(j.onsetRing) - keep392OnsetRing
			copy(j.onsetRing, j.onsetRing[cut:])
			j.onsetRing = j.onsetRing[:keep392OnsetRing]
		}
		j.vadBuf = append(j.vadBuf, res...)
		for len(j.vadBuf) >= speech.WindowSamples {
			win := j.vadBuf[:speech.WindowSamples]
			j.vadBuf = j.vadBuf[speech.WindowSamples:]
			j.eps = j.eng.det.Accept(j.noiseGate(win), j.eps[:0])
			for _, ep := range j.eps {
				if ep.Kind != speech.EndpointStart {
					continue
				}
				j.judge(ep)
			}
		}
	}
}

// noiseGate is the residual segmenter's floor gate, the bargeinDrain mirror:
// the floor is the minimum window RMS over the trailing two seconds, a
// window passes when it clears the floor by the step, and a passing window
// holds the gate open for the hold span.
func (j *keep392Judge) noiseGate(win []float32) []float32 {
	var sum float64
	for _, v := range win {
		sum += float64(v) * float64(v)
	}
	rms := math.Sqrt(sum / float64(len(win)))
	if j.floorRing == nil {
		j.floorRing = make([]float64, keep392FloorWindows)
	}
	j.floorRing[j.floorN%keep392FloorWindows] = rms
	j.floorN++
	floor := math.MaxFloat64
	for i, n := 0, min(j.floorN, keep392FloorWindows); i < n; i++ {
		if j.floorRing[i] < floor {
			floor = j.floorRing[i]
		}
	}
	if j.gateHold > 0 {
		j.gateHold--
	}
	if rms >= floor*keep392GateStep {
		j.gateHold = keep392GateHold
		return win
	}
	if j.gateHold > 0 {
		return win
	}
	return make([]float32, len(win))
}

// judge evaluates one episode start through the gate's real decision path
// and fires the trigger once (the drain's judge mirror).
func (j *keep392Judge) judge(ep speech.Endpoint) {
	n := ep.Total
	if n <= 0 || n > len(j.onsetRing) {
		n = len(j.onsetRing)
	}
	span := j.onsetRing[len(j.onsetRing)-n:]
	var sum float64
	for _, v := range span {
		sum += float64(v) * float64(v)
	}
	rms := math.Sqrt(sum / float64(max(1, len(span))))
	rmsDbfs := math.Inf(-1)
	if rms > 0 {
		rmsDbfs = 20 * math.Log10(rms)
	}
	dd := j.gate.Evaluate(audio.BargeinEpisode{
		StartSample: j.procSamples - ep.Total,
		Samples:     ep.Total,
		RmsDbfs:     rmsDbfs,
		AsrText:     "",
	})
	j.episodes++
	if dd.EnergyPass && !j.triggered {
		j.triggered = true
		j.triggerAt = j.procSamples - ep.Total
	}
}

// feedNear drains one 512-sample seam frame (the drain's frame mirror: the
// anchor is the first frame after far appeared; frames before it are read
// and counted, not judged).
func (j *keep392Judge) feedNear(frame []float32) {
	j.capCount += len(frame)
	if !j.anchorSet {
		if len(j.far) == 0 {
			return
		}
		j.anchorSet = true
		j.anchor = j.capCount - len(frame)
	}
	j.pending = append(j.pending, frame...)
	if !j.gateArmed {
		if !j.armJudge() {
			return
		}
	}
	j.process()
}

func keep392RmsDbfs(x []float32, from, to int) float64 {
	if to > len(x) {
		to = len(x)
	}
	if from >= to {
		return math.Inf(-1)
	}
	var sum float64
	for _, v := range x[from:to] {
		sum += float64(v) * float64(v)
	}
	rms := math.Sqrt(sum / float64(to-from))
	if rms < 1e-10 {
		return math.Inf(-1)
	}
	return 20 * math.Log10(rms)
}

func keep392Profile(x []float32, buckets int) []float64 {
	out := make([]float64, buckets)
	per := len(x) / buckets
	if per == 0 {
		per = 1
	}
	for b := 0; b < buckets; b++ {
		lo, hi := b*per, min((b+1)*per, len(x))
		if lo >= hi {
			break
		}
		out[b] = keep392RmsDbfs(x, lo, hi)
	}
	return out
}

// keep392RunWindow runs one announcement window under the arm's posture and
// returns its readings. The far tap feeds first (the deterministic schedule;
// the APM's processed render/capture pairing is the same lockstep the real
// interleaving produces), then the near frames flow and the judge runs; the
// trigger truncates the window (the abort) exactly as the staged rig's
// announcement abort stops the streams.
func keep392RunWindow(t *testing.T, mat *keep392Materials, eng *keep392Engines, w int) keep392WindowReading {
	t.Helper()
	j := &keep392Judge{mat: mat, eng: eng}
	j.arm(t)
	j.feedFar()
	for off := 0; off+keep392Frame <= len(mat.near); off += keep392Frame {
		if j.triggered {
			break
		}
		j.feedNear(mat.near[off : off+keep392Frame])
	}
	lo, hi := keep392EchoOnlyLate()
	r := keep392WindowReading{
		Window:          w,
		Episodes:        j.episodes,
		TriggerAtSample: -1,
		Lag:             j.lag,
		LagStrength:     j.lagStrength,
		ReadySamples:    j.readySamples,
		Processed:       j.procSamples,
		NearEchoDb:      keep392RmsDbfs(j.nearRec, lo, hi),
		EchoResidDb:     keep392RmsDbfs(j.resRec, lo, hi),
		NearProfile:     keep392Profile(j.nearRec, len(j.nearRec)/audio.TargetRate+1),
		ResProfile:      keep392Profile(j.resRec, len(j.resRec)/audio.TargetRate+1),
	}
	if j.triggered {
		r.Triggers = 1
		r.TriggerAtSample = j.triggerAt
	}
	j.boundary()
	return r
}

func keep392GateCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("WISP_AEC392R2") != "1" {
		t.Skip("keepalive isolation reading: set WISP_AEC392R2=1 (synthetic, deterministic, no audio hardware)")
	}
	if _, err := aec.NewApm(audio.TargetRate); err != nil {
		t.Skipf("AEC3 binding not in this build (-tags webrtcap + the deps.toml CGO flags): %v", err)
	}
}

func keep392OutDir(t *testing.T) string {
	t.Helper()
	dir := os.Getenv("WISP_AEC392R2_OUT")
	if dir == "" {
		dir = `D:\tmp\wisp\20261101\392r2b\out`
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("out dir: %v", err)
	}
	return dir
}

func keep392JSON(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal readings: %v", err)
	}
	return string(b)
}

// Test392R2KeepalivePoisonIsolationAndResetRevival is the leg's AC#1/AC#2/
// AC#3 reading: the piece-isolation matrix, the revival arm, the one-factor
// reverses, and the restore core, all on the deterministic replica, the
// whole matrix twice and identical.
func Test392R2KeepalivePoisonIsolationAndResetRevival(t *testing.T) {
	keep392GateCheck(t)
	mat := keep392MaterialsBuild(t)
	arms := keep392Arms()
	const windows = 3
	const runs = 2

	runArms := func() []keep392ArmReading {
		out := make([]keep392ArmReading, 0, len(arms))
		for _, p := range arms {
			eng := &keep392Engines{posture: p}
			ar := keep392ArmReading{Name: p.name, Posture: p.policy()}
			for w := 1; w <= windows; w++ {
				ar.Windows = append(ar.Windows, keep392RunWindow(t, mat, eng, w))
			}
			keep392ArmClose(eng)
			out = append(out, ar)
		}
		return out
	}

	var allRuns [][]keep392ArmReading
	for run := 1; run <= runs; run++ {
		allRuns = append(allRuns, runArms())
	}
	read := func(armIdx, w int) keep392WindowReading {
		return allRuns[0][armIdx].Windows[w-1]
	}
	armIdx := func(name string) int {
		for i, ar := range allRuns[0] {
			if ar.Name == name {
				return i
			}
		}
		t.Fatalf("arm %s missing", name)
		return 0
	}

	for run, rs := range allRuns {
		for _, ar := range rs {
			for _, r := range ar.Windows {
				t.Logf("392R2-R%d %-20s W%d episodes=%d triggers=%d trigger_at=%d lag=%d strength=%.3f ready=%d processed=%d near=%.2f resid=%.2f",
					run+1, ar.Name, r.Window, r.Episodes, r.Triggers, r.TriggerAtSample,
					r.Lag, r.LagStrength, r.ReadySamples, r.Processed, r.NearEchoDb, r.EchoResidDb)
				t.Logf("392R2-R%d %-20s W%d near_profile=%v res_profile=%v",
					run+1, ar.Name, r.Window, r.NearProfile, r.ResProfile)
			}
		}
	}

	// Instrument determinism: the two runs read identically.
	if b1, b2 := keep392JSON(t, allRuns[0]), keep392JSON(t, allRuns[1]); b1 != b2 {
		t.Fatalf("the two matrix runs disagree (the instrument is not deterministic): the readings cannot isolate anything")
	}

	// Instrument health: the cold posture's three windows are fresh-engine
	// windows on identical stimuli, so they read identically; and every
	// arm's W1 is the same fresh-engine window, so it reads identically to
	// the cold W1.
	// Instrument health: the cold posture's three windows are fresh-engine
	// windows on identical stimuli, so they read identically once the window
	// index is out of the key; and every arm's W1 is the same fresh-engine
	// window, so it reads identically to the cold W1.
	cold := allRuns[0][armIdx("cold-control")]
	readingKey := func(r keep392WindowReading) string {
		r.Window = 0
		return keep392JSON(t, r)
	}
	coldW1 := readingKey(cold.Windows[0])
	if readingKey(cold.Windows[1]) != coldW1 || readingKey(cold.Windows[2]) != coldW1 {
		t.Fatalf("the cold posture's identical fresh-engine windows do not read identically: the instrument is not deterministic\nW1=%s\nW2=%s\nW3=%s",
			keep392JSON(t, cold.Windows[0]), keep392JSON(t, cold.Windows[1]), keep392JSON(t, cold.Windows[2]))
	}
	for _, ar := range allRuns[0] {
		if readingKey(ar.Windows[0]) != coldW1 {
			t.Fatalf("arm %s W1 does not read as a fresh-engine window (against the cold W1): the arms are not comparable\narm=%s\ncold=%s",
				ar.Name, keep392JSON(t, ar.Windows[0]), keep392JSON(t, cold.Windows[0]))
		}
	}
	t.Logf("392R2-HEALTH every arm's W1 reads identically (fresh engines, identical stimuli): the matrix is comparable")

	coldW2 := cold.Windows[1]
	w392, isoDet, isoRes := armIdx("warm-392"), armIdx("iso-det-fresh"), armIdx("iso-res-fresh")
	rev, revCool, revNR, revNF := armIdx("revival"), armIdx("revival-cool"), armIdx("revival-no-reset"), armIdx("revival-no-finalize")

	// The instrument's own cancellation health: the cold W2's echo-only
	// residual must sit well under the near echo (the APM is fresh and the
	// far is clean there - if this fails the depth readings are blind and
	// the resampler half cannot be judged with them).
	if coldW2.EchoResidDb > coldW2.NearEchoDb-keep392MinCancelDb {
		t.Logf("392R2-DEPTH-BLIND cold W2 cancels only %.2f dB at the late echo-only span (near %.2f, resid %.2f): the depth readings are weak evidence on this span",
			coldW2.NearEchoDb-coldW2.EchoResidDb, coldW2.NearEchoDb, coldW2.EchoResidDb)
	}

	// AC#1, the finding reproduced in-rig: the 392 warm shape (detector +
	// resampler held, APM fresh) does not judge in W2. W3 recovering is the
	// mechanism's own signature, not a contradiction: W2 ran the full
	// window (no trigger, no abort), its long gated-silence tail exceeds
	// the endpointer's 700ms close threshold, and the carried utterance
	// closes - the poison is the open utterance's lifetime, not a permanent
	// warm-state corruption.
	if r := read(w392, 2); r.Episodes != 0 || r.Triggers != 0 {
		t.Fatalf("POISON NOT REPRODUCED: warm-392 W2 judged (episodes=%d triggers=%d) - the 392-r1 finding does not reproduce on this replica; the matrix cannot isolate and the root-cause chain needs re-measurement on the staged rig",
			r.Episodes, r.Triggers)
	}
	if r := read(w392, 3); r.Episodes == 0 {
		t.Fatalf("MECHANISM MOVED: warm-392 W3 still judges nothing (episodes=%d) - the held utterance should have closed during W2's full-window silence tail; the open-utterance lifetime story does not hold",
			r.Episodes)
	}
	t.Logf("392R2-AC1 FINDING reproduced in-rig: warm-392 W2 judges nothing (episodes=0, triggers=0, full window) while its residual sits at cancelled depth (%.2f vs cold %.2f) - the judge, not the APM, is what died; W3 recovers (episodes=%d) once W2's silence tail closes the carried utterance",
		read(w392, 2).EchoResidDb, coldW2.EchoResidDb, read(w392, 3).Episodes)

	// AC#1, the resampler half isolated: with only the detector held the
	// window still does not judge (the detector carries the poison,
	// resampler-independent), and with only the resampler held the judge
	// fires and the far path cancels as well as the fresh resampler's (same
	// fresh-APM convergence schedule, so the comparison reads the far
	// reference, not the head).
	if r := read(isoRes, 2); r.Episodes != 0 || r.Triggers != 0 {
		t.Fatalf("ISOLATION READS BOTH: iso-res-fresh (detector warm only) W2 judged (episodes=%d triggers=%d) - the detector's held state alone is not sufficient; re-read the matrix before naming the poison",
			r.Episodes, r.Triggers)
	}
	if r := read(isoDet, 2); r.Triggers != 1 {
		t.Fatalf("ISOLATION READS BOTH: iso-det-fresh (resampler warm only) W2 did not judge (episodes=%d triggers=%d) - the resampler's held state blocks judging; the poison is not the detector alone and a third piece must be named",
			r.Episodes, r.Triggers)
	}
	rIsoDet := read(isoDet, 2)
	if rIsoDet.EchoResidDb > coldW2.EchoResidDb+1.5 || rIsoDet.EchoResidDb > rIsoDet.NearEchoDb-keep392MinCancelDb {
		t.Fatalf("RESAMPLER READS DIRTY: iso-det-fresh W2 echo-only residual %.2f dBFS (cold %.2f, near echo %.2f) - the held resampler corrupts the far reference and the resampler half is guilty too",
			rIsoDet.EchoResidDb, coldW2.EchoResidDb, rIsoDet.NearEchoDb)
	}
	t.Logf("392R2-AC1 ISOLATION: detector warm (resampler fresh) still does not judge; resampler warm (detector fresh) judges with the residual at %.2f dBFS (cold %.2f, near echo %.2f) - the poison is the held detector's cross-window state, the held resampler's far output is clean",
		rIsoDet.EchoResidDb, coldW2.EchoResidDb, rIsoDet.NearEchoDb)

	// AC#2, the revival. The naive Finalize revival's W2 reading carries a
	// named transient the probe measures directly: the abort leaves the
	// silero RNN in a speech context, the session sat idle (no input, no
	// decay), and the hangover fires a start on the fresh window's head -
	// the episode fires (the judge is alive again, the open-utterance block
	// is gone) but it consumes the endpointer before the talker arrives, so
	// the trigger misses. The cool-down completes the clear face: the same
	// posture with the silence cool-down re-judges W2 AND W3 at the talker.
	rRev := read(rev, 2)
	if rRev.Episodes == 0 {
		t.Fatalf("REVIVAL FAILED: revival W2 produced no episode (episodes=%d) - the detector clear face did not return the judge",
			rRev.Episodes)
	}
	if rRev.EchoResidDb > coldW2.EchoResidDb+1.5 || rRev.EchoResidDb > rRev.NearEchoDb-keep392MinCancelDb {
		t.Fatalf("REVIVAL FAILED: revival W2 echo-only residual %.2f dBFS (cold %.2f, near echo %.2f) - the second window does not hear a cancelled stream",
			rRev.EchoResidDb, coldW2.EchoResidDb, rRev.NearEchoDb)
	}
	t.Logf("392R2-AC2 REVIVAL W2 (Finalize only): the judge is alive again (episodes=%d, the open-utterance block is gone) and the residual is at cancelled depth (%.2f vs cold %.2f); the trigger missed (triggers=%d) on the RNN-hangover transient the probe pins - the first warm window fired a start on the post-abort context and the talker's onset never opened",
		rRev.Episodes, rRev.EchoResidDb, coldW2.EchoResidDb, rRev.Triggers)
	if r := read(revCool, 2); r.Triggers != 1 || r.Episodes == 0 {
		t.Fatalf("REVIVAL FAILED: revival-cool W2 did not re-judge at the talker (episodes=%d triggers=%d) - the silence cool-down did not complete the detector's clear face",
			r.Episodes, r.Triggers)
	}
	if r := read(revCool, 2); r.EchoResidDb > coldW2.EchoResidDb+1.5 || r.EchoResidDb > r.NearEchoDb-keep392MinCancelDb {
		t.Fatalf("REVIVAL FAILED: revival-cool W2 echo-only residual %.2f dBFS (cold %.2f, near echo %.2f) - the second window does not hear a cancelled stream",
			r.EchoResidDb, coldW2.EchoResidDb, r.NearEchoDb)
	}
	if r := read(revCool, 3); r.Triggers != 1 {
		t.Fatalf("REVIVAL FAILED: revival-cool W3 did not judge (episodes=%d triggers=%d) - the revival does not persist",
			r.Episodes, r.Triggers)
	}
	t.Logf("392R2-AC2 REVIVAL holds: keepalive explicitly on, every held piece cleared at the boundary (det Finalize + silence cool-down / resampler Flush / APM Reset) - W2 and W3 re-judge at the talker (episodes=%d triggers=%d trigger_at=%d, residual %.2f vs cold %.2f), W3 (episodes=%d triggers=%d)",
		read(revCool, 2).Episodes, read(revCool, 2).Triggers, read(revCool, 2).TriggerAtSample,
		read(revCool, 2).EchoResidDb, coldW2.EchoResidDb, read(revCool, 3).Episodes, read(revCool, 3).Triggers)

	// AC#3, the one-factor reverses. Removing the detector's clear face
	// brings the 392 finding back (episodes=0). Removing the APM's Reset
	// brings the residual death back: the carried APM's W2 echo-only
	// residual sits at the near level (392-r1's candidate-2 "thoroughly
	// mis-adapted" shape, reproduced in-chain) - the Reset is load-bearing
	// for the revival's cancellation health, not belt-and-braces.
	if r := read(revNF, 2); r.Episodes != 0 || r.Triggers != 0 {
		t.Fatalf("REVERSE FAILED: revival-no-finalize W2 judged (episodes=%d triggers=%d) - the detector clear face is not load-bearing on this replica",
			r.Episodes, r.Triggers)
	}
	t.Logf("392R2-AC3 REVERSE det: revival minus Finalize returns the finding (episodes=0, triggers=0) - the clear face is load-bearing")
	rNR := read(revNR, 2)
	if rNR.EchoResidDb <= rRev.EchoResidDb+3 {
		t.Fatalf("REVERSE FAILED: revival-no-reset W2 residual %.2f is not degraded against the revival arm (%.2f) - the APM Reset is not load-bearing on this shape; report the reading instead",
			rNR.EchoResidDb, rRev.EchoResidDb)
	}
	t.Logf("392R2-AC3 REVERSE apm: revival-no-reset W2 residual %.2f vs revival %.2f (degraded %.2f dB) with the near echo at %.2f - the carried APM cancels to nothing (392-r1 candidate-2's residual-at-near-level shape), the Reset removes it",
		rNR.EchoResidDb, rRev.EchoResidDb, rNR.EchoResidDb-rRev.EchoResidDb, rNR.NearEchoDb)

	// The restore core: after everything, a fresh cold window re-reads the
	// instrument's opening shape bit-exactly.
	eng := &keep392Engines{posture: keep392Posture{name: "restore-core"}}
	restore := keep392RunWindow(t, mat, eng, 1)
	keep392ArmClose(eng)
	if readingKey(restore) != coldW1 {
		t.Fatalf("the restore core drifted against the opening cold W1: the instrument chain did not hold the run")
	}
	t.Logf("392R2-RESTORE core matches the opening cold W1 bit-exactly")

	artifact := struct {
		Stimulus string                `json:"stimulus"`
		Runs     [][]keep392ArmReading `json:"runs"`
		Finding  string                `json:"finding"`
	}{
		Stimulus: "real matcha speech: the production announcement (peak -10 dBFS) in its synthesis chunking through the streaming resampler far tap, matched echo 0.2 at the 26p return lag (480 samples), talker -20 dBFS RMS at 400ms; W1 builds, W2 discriminates, W3 persists; the whole matrix twice",
		Runs:     allRuns,
		Finding:  "the 392 keepalive poison is the held detector's cross-window state: the abort leaves the endpointer mid-utterance and no start can fire while an utterance is open, so the next warm window judges nothing (episodes=0) while its APM residual stays at cancelled depth - the held resampler's far output is clean (iso-det-fresh judges with the residual identical to cold). The 392-r1 residual-at-near-level shape belongs to the APM-carry half: carrying the APM without Reset degrades the next window's cancellation to ~0dB (residual at the near echo), and Reset() at the boundary removes it. The revival is the hold-and-clear discipline - det Finalize (+ silence cool-down: the RNN carries the abort's speech context across the idle gap and the probe measures a hangover start firing on post-speech silence) + resampler Flush + APM Reset - and W2/W3 re-judge at the talker with the residual at cancelled depth.",
	}
	b, err := json.MarshalIndent(artifact, "", "  ")
	if err != nil {
		t.Fatalf("marshal artifact: %v", err)
	}
	p := filepath.Join(keep392OutDir(t), "392r2-keepalive-matrix.json")
	if err := os.WriteFile(p, append(b, '\n'), 0o644); err != nil {
		t.Fatalf("write %s: %v", p, err)
	}
	t.Logf("WROTE %s", p)
}

// Test392R2DetectorOpenUtteranceCarries pins the isolated mechanism behind
// the AC#1 matrix at unit level: an utterance the abort leaves open blocks
// every later start until 700ms of consecutive VAD silence closes it - and
// the announcement window's echo-only pre-talker span (400ms) is shorter
// than that. Finalize (the detector's clear face) removes the block with
// the silero session still warm. This is the "replace" leg of the
// piece-by-piece isolation: the state is named, read directly, and removed.
// The control also reads the hangover the matrix's revival transient comes
// from: the RNN still holds the talker's speech context (the session sat
// idle, no input, no decay), and a start can fire on post-speech silence.
func Test392R2DetectorOpenUtteranceCarries(t *testing.T) {
	keep392GateCheck(t)
	mat := keep392MaterialsBuild(t)

	countStarts := func(det *speech.Detector, samples []float32) int {
		n := 0
		var eps []speech.Endpoint
		for off := 0; off+speech.WindowSamples <= len(samples); off += speech.WindowSamples {
			eps = det.Accept(samples[off:off+speech.WindowSamples], eps[:0])
			for _, ep := range eps {
				if ep.Kind == speech.EndpointStart {
					n++
				}
			}
		}
		return n
	}
	talker := mat.talker
	silence := make([]float32, keep392TalkerAt) // the echo-only span, gated to digital silence

	// The carry shape: the first talker opens an utterance, the abort leaves
	// it open, 400ms of gated silence is under the 700ms close threshold,
	// and the second talker therefore produces NO new start.
	det := keep392NewDetector(t)
	defer func() { _ = det.Close() }()
	if n := countStarts(det, talker); n == 0 {
		t.Fatalf("the probe talker did not open an utterance: the mechanism probe has no subject")
	}
	countStarts(det, silence)
	if n := countStarts(det, talker); n != 0 {
		t.Fatalf("MECHANISM MOVED: the second talker produced %d start(s) across the carried utterance - the open-utterance carry does not hold", n)
	}
	var eps []speech.Endpoint
	for off := 0; off+speech.WindowSamples <= len(silence); off += speech.WindowSamples {
		eps = det.Accept(silence[off:off+speech.WindowSamples], eps[:0])
	}
	eps = det.Finalize(eps)
	closed := 0
	for _, ep := range eps {
		if ep.Kind == speech.EndpointEndOfStream && ep.Total > 2*len(talker) {
			closed++
		}
	}
	if closed == 0 {
		t.Fatalf("the carried utterance did not close at Finalize with the accumulated total: the carry shape is not the open utterance")
	}
	t.Logf("392R2-MECHANISM the open utterance carried across the 400ms echo-only span and swallowed the second talker's start; Finalize closed it with the accumulated total")

	// The control: the same sequence with the clear face between windows -
	// the block is gone and a start fires again. Where it fires is itself a
	// finding: the RNN still carries the first talker's speech context (the
	// session sat idle, no input, no decay), so the hangover can feed the
	// endpointer's speech run through the SILENCE and the start lands there
	// rather than on the second talker. Either landing proves the endpointer
	// is idle again; the silence landing is the hangover the revival's
	// cool-down exists to absorb.
	det2 := keep392NewDetector(t)
	defer func() { _ = det2.Close() }()
	if n := countStarts(det2, talker); n == 0 {
		t.Fatalf("the control talker did not open an utterance")
	}
	_ = det2.Finalize(nil)
	nSil := countStarts(det2, silence)
	nTalk := countStarts(det2, talker)
	if nSil+nTalk == 0 {
		t.Fatalf("MECHANISM MOVED: after the clear face no start fires anywhere (silence=%d talker=%d) - Finalize does not return the endpointer to idle", nSil, nTalk)
	}
	t.Logf("392R2-MECHANISM control: with Finalize between windows a start fires again (silence=%d talker=%d) - the clear face returns the endpointer to idle; a silence landing is the RNN hangover feeding on post-speech silence", nSil, nTalk)
}
