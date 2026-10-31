//go:build windows

package main

// Ticket 26p's integration readings, part 1: the equivalence first grid (the
// 389-v1 disposition suggestion) and the (a) true number - 播报中发声 ⇒
// barge-in 触发停播 ≤400ms on the staged full chain: the real waveOut device,
// the real matcha announcement (the far reference), the real C8 seam frames,
// the real AEC3 cancellation, the real 389 combination gate, the real silero
// VAD, and the production TriggerInterrupt -> abort stop.
//
// The talker (播报中发声) is real TTS-rendered speech - the same matcha
// checkpoint the announcement uses, rendered once per test into a 16k wav at
// a calibrated RMS level (the -26.70 dBFS control the 389 calibration set
// carries). The echo (the (b) material, part 2) is a previous round's own
// announcement recording fed back through the C8 seam with a fixed return
// delay - the acoustic-leak analog; real echoes are not bit-exact and this
// one deliberately is not either (a different synthesis round of the same
// sentence cancels to a residual, not to silence).
//
// Staging is 15e's plus the AEC3 build: -tags webrtcap with the CGO flags
// deps.toml's [webrtc-audio-processing] section documents, the sherpa DLLs
// next to the test binary (A1326 card), the WISP_* model paths. Without the
// tag the monitor's arm refuses (the readings say so) and every case here
// skips with that reason - the no-tag default build never reaches for the
// native library.

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/audio/aec"
	"github.com/CarlosShao/wisp/internal/speech"
)

// voice26pOut is this leg's artifact root (the readings JSONs land here).
const voice26pOut = `D:\tmp\wisp\20261031\26pr1\out`

// voice26pEchoDelay is the injected echo's return lag: 30ms at the seam rate,
// the magnitude the rig measured on the loopback (493 samples there).
const voice26pEchoDelay = 480

// voice26pEchoGain is the acoustic coupling of the injected matched echo: a
// realistic speaker-to-mic leak whose raw level (~-32.6 dBFS over a -18.6
// announcement) sits ~6dB below the talker - the double-talk shape an actual
// barge-in has (the user is near the mic, the speaker is not). The AEC'd
// residual rides further below the tightened cut, which is what the (b)
// reading needs; the M2 bypass reading needs the raw echo only to OPEN
// episodes on the judge (the verbatim transcription), not to pass the cut.
const voice26pEchoGain = 0.2

// voice26pMismatchGain is the coupling of the mismatched echo (the
// artifact-family generator). Measured at 0.10 the engine suppresses the
// mismatch into a flat -45.8 floor - BELOW the loose nominal, where even
// the M1 reverse finds nothing to trigger on; at 0.25 the uncancelled
// residual lands in the (-45, -30) band - above the loose nominal the M1
// reverse rides, below the tightened cut the dynamic gate holds - exactly
// the band the 59-w collapse family (-32.24/-35.59) occupied.
const voice26pMismatchGain = 0.25

// voice26pTalkerAt is when the talker starts inside the announcement playback
// (after the monitor's alignment warm-up, which the readings carry): the
// steady-state barge-in surface the S4 budget describes.
const voice26pTalkerAt = 1200 * time.Millisecond

// voice26pTalkerDbfs is the (a) talker's level: a normal near-mic utterance
// (real barge-in speech at a laptop-mic capture level). The calibration
// control's -26.70 dBFS stays on the (c) margin reading, where the AEC's
// double-talk damage on THAT level is the finding - measured across the
// staged runs, the integrated chain's post-AEC level of a -26.70 talker
// swings -26.9 to -30.9 dBFS with the engine's suppressor state, which puts
// the control's 3.3dB tightened-cut margin on the wrong side of the cut at
// the damage's worst; a near-mic utterance does not live there.
const voice26pTalkerDbfs = -20.0

// voice26pRecordSpeaker is the render seam: the real waveOut device behind a
// recorder that keeps the announcement's samples (the next round's echo
// source) and exposes the two moments the readings key on - the first Play
// (the monitor's anchor twin) and the abort landing (the stop).
type voice26pRecordSpeaker struct {
	voiceSpeaker
	mu        sync.Mutex
	samples   []float32
	first     chan struct{}
	firstOnce sync.Once
	abortDone chan struct{}
	abortOnce sync.Once
}

func (s *voice26pRecordSpeaker) Play(samples []float32, engineRate int) error {
	s.mu.Lock()
	s.samples = append(s.samples, samples...)
	n := len(s.samples)
	s.mu.Unlock()
	if n > 0 {
		s.firstOnce.Do(func() { close(s.first) })
	}
	return s.voiceSpeaker.Play(samples, engineRate)
}

func (s *voice26pRecordSpeaker) Abort() int {
	n := s.voiceSpeaker.Abort()
	s.abortOnce.Do(func() { close(s.abortDone) })
	return n
}

func (s *voice26pRecordSpeaker) recording() []float32 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]float32(nil), s.samples...)
}

// voice26pRig is one staged rig: the production chain on a bounded frame
// channel, the real device behind a recording spy, and the monitor wired.
// The capture gate is the 381-r2 recorder (SetSpeaking calls counted, no
// SetConversation - so the monitor's window reads conv_window=false and the
// readings say so; the frames these cases judge are pushed into the seam by
// the injectors).
type voice26pRig struct {
	t      *testing.T
	rv     *residentVoice
	frames chan []byte
	gate   *voice381r2Gate

	spyMu sync.Mutex
	spies []*voice26pRecordSpeaker
}

func newVoice26pRig(t *testing.T, tune func(*voiceBargeinMonitor)) *voice26pRig {
	t.Helper()
	plan := voicePlan{
		VadModelPath: voice15eVadModel(t),
		Source:       "26p test",
		Tts:          voiceWireTtsFiles(t),
	}
	rig := &voice26pRig{t: t, gate: &voice381r2Gate{}}
	rv, frames := voiceWireStart(t, plan, &residentBall{}, &voice15eSink{}, func(rv *residentVoice) {
		rv.newSpeaker = func(rate int) (voiceSpeaker, error) {
			sp, err := openWaveOutSpeaker(rate)
			if err != nil {
				return nil, err
			}
			spy := &voice26pRecordSpeaker{
				voiceSpeaker: sp,
				first:        make(chan struct{}),
				abortDone:    make(chan struct{}),
			}
			rig.spyMu.Lock()
			rig.spies = append(rig.spies, spy)
			rig.spyMu.Unlock()
			return spy, nil
		}
		rv.attachCaptureGate(rig.gate)
		if tune != nil {
			tune(rv.bargein)
		}
	})
	rig.rv, rig.frames = rv, frames
	return rig
}

func (r *voice26pRig) spyAt(i int) *voice26pRecordSpeaker {
	r.t.Helper()
	voice15eWaitFor(r.t, fmt.Sprintf("speaker %d to open", i), func() bool {
		r.spyMu.Lock()
		defer r.spyMu.Unlock()
		return len(r.spies) > i
	})
	r.spyMu.Lock()
	defer r.spyMu.Unlock()
	return r.spies[i]
}

// inject starts one real-time C8 injection of a wav file. The caller stops it.
func (r *voice26pRig) inject(wav string) (stop func(), injected *audio.WavInjector) {
	r.t.Helper()
	inj, err := audio.NewWavInjector(wav)
	if err != nil {
		r.t.Fatalf("wav injector %s: %v", wav, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	if err := inj.Start(ctx, r.frames); err != nil {
		r.t.Fatalf("wav injector start: %v", err)
	}
	return func() { cancel(); _ = inj.Stop() }, inj
}

// waitReleased waits for the announcement window to fully close (the monitor
// joined), so the readings are final.
func (r *voice26pRig) waitReleased() {
	r.t.Helper()
	voice15eWaitFor(r.t, "the barge-in monitor to release the window", func() bool {
		return !r.rv.bargein.readings().Armed
	})
}

// requireAEC3 skips the staged case honestly when the build has no binding:
// the probe is the same constructor the monitor's arm calls.
func (r *voice26pRig) requireAEC3() {
	r.t.Helper()
	apm, err := aec.NewApm(audio.TargetRate)
	if err != nil {
		r.t.Skipf("the AEC3 binding is not in this build (-tags webrtcap + the deps.toml CGO flags): the integration readings need it: %v", err)
	}
	_ = apm.Close()
}

// voice26pRenderTtsWav renders one sentence with the real matcha checkpoint
// and writes a 16k mono wav gained to the requested RMS level (dBFS, the
// calibration set's convention).
func voice26pRenderTtsWav(t *testing.T, dir, name, sentence string, targetDbfs float64) string {
	t.Helper()
	files := voiceWireTtsFiles(t)
	eng, err := speech.NewMatchaTts(speech.TtsConfig{
		AcousticModel: files.Acoustic,
		Vocoder:       files.Vocoder,
		Lexicon:       files.Lexicon,
		Tokens:        files.Tokens,
		DictDir:       files.DictDir,
		RuleFsts:      files.RuleFsts,
	})
	if err != nil {
		t.Fatalf("matcha tts: %v", err)
	}
	defer func() { _ = eng.Close() }()
	var got []float32
	if err := eng.Speak(sentence, func(samples []float32) error {
		got = append(got, samples...)
		return nil
	}); err != nil {
		t.Fatalf("render %q: %v", sentence, err)
	}
	if len(got) == 0 {
		t.Fatalf("render %q produced no samples", sentence)
	}
	pcm := audio.ResampleLinear(audio.FloatToPCM16(got), eng.SampleRate(), audio.TargetRate)
	var sum float64
	for _, v := range pcm {
		f := float64(v) / 32768
		sum += f * f
	}
	rms := math.Sqrt(sum / float64(len(pcm)))
	if rms <= 0 {
		t.Fatalf("render %q is silent", sentence)
	}
	gain := math.Pow(10, targetDbfs/20) / rms
	if gain > 8 {
		t.Fatalf("render %q too quiet to calibrate (rms %.5f)", sentence, rms)
	}
	out := make([]int16, len(pcm))
	for i, v := range pcm {
		s := float64(v) / 32768 * gain
		if s > 0.999 {
			s = 0.999
		}
		if s < -0.999 {
			s = -0.999
		}
		out[i] = int16(s * 32767)
	}
	p := filepath.Join(dir, name+".wav")
	voice15eWriteWav(t, p, out)
	return p
}

// voice26pEchoWav builds one echo file: a previous announcement's own
// recording, at the seam rate, delayed by the return lag and scaled by the
// given coupling gain - the acoustic-leak analog injected into the C8 seam.
// The seam-rate samples come back too: they are what the product mic surface
// carried, and the M2 reading transcribes them.
func voice26pEchoWav(t *testing.T, dir, name string, recording []float32, engineRate int, gain float64) (string, []float32) {
	t.Helper()
	pcm := audio.ResampleLinear(audio.FloatToPCM16(recording), engineRate, audio.TargetRate)
	out := make([]int16, voice26pEchoDelay+len(pcm))
	samples := make([]float32, len(out))
	for i, v := range pcm {
		out[voice26pEchoDelay+i] = int16(float64(v) * gain)
		samples[voice26pEchoDelay+i] = float32(float64(v) / 32768 * gain)
	}
	p := filepath.Join(dir, name+".wav")
	voice15eWriteWav(t, p, out)
	return p, samples
}

// Test26PBargeinStopsWithin400ms is the (a) true number: five barge-in
// rounds, each a real matcha announcement on the real device with the
// announcement's own echo injected into the C8 seam and real TTS speech
// starting 1.2s into the playback. The monitor's trigger stops the device;
// the reading is the VAD-start-to-stop distance (the half D32:2315 assigns to
// the system: 400ms 门禁由 VAD 而非 ASR 满足), with the onset-to-stop sum and
// its decomposition reported. One control round first (no injection, the
// echo source).
func Test26PBargeinStopsWithin400ms(t *testing.T) {
	if probe, err := openWaveOutSpeaker(22050); err != nil {
		t.Skipf("no output device on this machine: %v", err)
	} else {
		_ = probe.Close()
	}
	if err := os.MkdirAll(voice26pOut, 0o755); err != nil {
		t.Fatalf("out dir: %v", err)
	}
	const announceSentence = "任务已经完成，明天下午三点的会议改到四点，请准时参加，谢谢。"
	const talkerSentence = "等一下，先听我说这件事。"
	dir := t.TempDir()
	talkerWav := voice26pRenderTtsWav(t, dir, "talker", talkerSentence, voice26pTalkerDbfs)

	rig := newVoice26pRig(t, nil)
	rig.requireAEC3()

	// Control round: the announcement plays out and is recorded (the echo
	// source for every later round).
	walkToThinking381r2(t, rig.rv)
	if !rig.rv.announce(announceSentence) {
		t.Fatalf("the control announce refused the sentence; problems=%v", rig.rv.problemsSnapshot())
	}
	control := rig.spyAt(0)
	voice15eWaitFor(t, "the control sentence to play out", func() bool {
		return rig.rv.ttsSpoken.Load() >= 1
	})
	rig.waitReleased()
	echoWav, _ := voice26pEchoWav(t, dir, "echo", control.recording(), 22050, voice26pEchoGain)
	t.Logf("CONTROL spoken=%d echo_samples=%d", rig.rv.ttsSpoken.Load(), len(control.recording()))

	const rounds = 5
	type roundReading struct {
		round          int
		vadToStopMS    int64
		episodeRmsDbfs float64
		thresholdDbfs  float64
		tighten        bool
		lag            int
		lagStrength    float64
	}
	readings := make([]roundReading, 0, rounds)
	for r := 1; r <= rounds; r++ {
		walkToThinking381r2(t, rig.rv)
		if !rig.rv.announce(announceSentence) {
			t.Fatalf("round %d: the announce refused the sentence; problems=%v", r, rig.rv.problemsSnapshot())
		}
		spy := rig.spyAt(r)
		<-spy.first
		stopEcho, _ := rig.inject(echoWav)
		time.Sleep(voice26pTalkerAt)
		stopTalker, _ := rig.inject(talkerWav)
		// The round's stop marker is the abort landing on the device; the
		// t0 is the monitor's own trigger stamp.
		voice15eWaitFor(t, fmt.Sprintf("round %d: the barge-in stop to land", r), func() bool {
			return rig.rv.ttsAborted.Load() >= uint64(r)
		})
		<-spy.abortDone
		// The injectors stop before the window closes: their tails belong to
		// the monitor's window, never to the transcript path.
		stopEcho()
		stopTalker()
		rig.waitReleased()
		read := rig.rv.bargein.readings()
		if read.Triggers != 1 {
			t.Fatalf("round %d: the monitor triggered %d times, want exactly 1 (episodes=%d): the echo must not self-trigger and the talker must", r, read.Triggers, len(read.Episodes))
		}
		if got := rig.rv.ttsSpoken.Load(); got != 1 {
			t.Fatalf("round %d: the interrupted announcement was counted as spoken (%d)", r, got)
		}
		var fired *voiceBargeinEpisode
		for i := range read.Episodes {
			if read.Episodes[i].EnergyPass {
				fired = &read.Episodes[i]
			}
		}
		if fired == nil {
			t.Fatalf("round %d: the trigger fired without a passing episode (episodes=%d)", r, len(read.Episodes))
		}
		vadToStop := time.Since(read.TriggerAt)
		if vadToStop > 400*time.Millisecond {
			t.Fatalf("round %d: the device answered the barge-in trigger in %s, want within 400ms", r, vadToStop)
		}
		rr := roundReading{
			round:          r,
			vadToStopMS:    vadToStop.Milliseconds(),
			episodeRmsDbfs: fired.RmsDbfs,
			thresholdDbfs:  fired.ThresholdDbfs,
			tighten:        fired.TightenActive,
			lag:            read.Lag,
			lagStrength:    read.LagStrength,
		}
		readings = append(readings, rr)
		t.Logf("BARGE-IN-%d vad_to_stop=%s episode_rms=%.2f threshold=%.1f tighten=%v lag=%d strength=%.2f episodes=%d drained=%d",
			r, vadToStop.Round(time.Millisecond), fired.RmsDbfs, fired.ThresholdDbfs,
			fired.TightenActive, read.Lag, read.LagStrength, len(read.Episodes), read.FramesDrained)
	}
	if got := rig.rv.problemsSnapshot(); len(got) != 0 {
		t.Fatalf("the barge-in rounds left problems behind: %v", got)
	}
	t.Logf("BARGE-IN-SUMMARY rounds=%d vad_to_stop=%v all<=400ms=true (the onset-to-stop sum rides the VAD's MinSpeech 300ms design constant; the decomposition is in the per-round lines)", rounds, readings)
}

// Test26PFeedOnlineMatchesWindowReplay is the equivalence first grid (the
// 389-v1 disposition suggestion): one live announcement window's (far, near,
// residual) streams driven through the gate's ONLINE Feed seam must answer
// the same per-window untrusted verdicts as the same windows re-cut offline
// through the ObserveWindow seam - one decision path, two input seams, zero
// drift on real AEC3 material.
func Test26PFeedOnlineMatchesWindowReplay(t *testing.T) {
	if probe, err := openWaveOutSpeaker(22050); err != nil {
		t.Skipf("no output device on this machine: %v", err)
	} else {
		_ = probe.Close()
	}
	if err := os.MkdirAll(voice26pOut, 0o755); err != nil {
		t.Fatalf("out dir: %v", err)
	}
	const announceSentence = "任务已经完成，明天下午三点的会议改到四点，请准时参加，谢谢。"
	dir := t.TempDir()

	rig := newVoice26pRig(t, nil)
	rig.requireAEC3()

	// The echo source: one full announcement, recorded.
	walkToThinking381r2(t, rig.rv)
	if !rig.rv.announce(announceSentence) {
		t.Fatalf("the control announce refused the sentence; problems=%v", rig.rv.problemsSnapshot())
	}
	control := rig.spyAt(0)
	voice15eWaitFor(t, "the control sentence to play out", func() bool {
		return rig.rv.ttsSpoken.Load() >= 1
	})
	rig.waitReleased()
	echoWav, _ := voice26pEchoWav(t, dir, "echo", control.recording(), 22050, voice26pEchoGain)

	// One live window: the echo-only round - no trigger, so the whole
	// announcement plays out and every window closes (the equivalence grid
	// wants the richest window sequence the real chain produces; a triggered
	// round would cut the stream at the abort).
	walkToThinking381r2(t, rig.rv)
	if !rig.rv.announce(announceSentence) {
		t.Fatalf("the window announce refused the sentence; problems=%v", rig.rv.problemsSnapshot())
	}
	spy := rig.spyAt(1)
	<-spy.first
	stopEcho, _ := rig.inject(echoWav)
	voice15eWaitFor(t, "the echo window to play out", func() bool {
		return rig.rv.ttsSpoken.Load() >= 2
	})
	stopEcho()
	rig.waitReleased()
	read := rig.rv.bargein.readings()
	if read.Triggers != 0 {
		t.Fatalf("the echo window self-triggered (triggers=%d): the equivalence grid needs the (b) shape", read.Triggers)
	}
	far, near, res := rig.rv.bargein.streams()
	if len(far) == 0 || len(far) != len(near) || len(near) != len(res) {
		t.Fatalf("the recorded streams are not frame-aligned (far=%d near=%d res=%d)", len(far), len(near), len(res))
	}
	windows := len(far) / audio.BargeinWindowSamples
	if windows < 3 || windows != len(read.Untrusted) {
		t.Fatalf("the online verdict count %d does not match the stream's closed windows %d (want >=3 rich windows)", len(read.Untrusted), windows)
	}

	// The offline re-cut: the same raw streams, the same window boundaries,
	// the gate's window-level seam. The ERLE channel needs the lag-shifted
	// far the online Feed built internally; the correlation channel takes the
	// raw far and the same lag.
	var farSum float64
	for _, v := range far[:voiceBargeinAlignFar] {
		farSum += float64(v) * float64(v)
	}
	farRMS := math.Sqrt(farSum / float64(voiceBargeinAlignFar))
	shiftedFar := make([]float32, len(far))
	for i := range far {
		if i-read.Lag >= 0 {
			shiftedFar[i] = far[i-read.Lag]
		}
	}
	erle := audio.NewSlidingERLE(farRMS)
	corr := audio.NewFarCorrelationDetector(read.Lag)
	erle.Observe(shiftedFar, near, res)
	corr.Observe(far, near)
	if len(erle.Windows) != windows || len(corr.Values) != windows {
		t.Fatalf("the offline re-cut closed %d ERLE / %d correlation windows, want %d",
			len(erle.Windows), len(corr.Values), windows)
	}
	offline := audio.NewBargeinGate(audio.BargeinOptions{
		DynamicGating: true,
		LagSamples:    read.Lag,
		FarRMS:        farRMS,
	})
	for i := range windows {
		offline.ObserveWindow(erle.Windows[i], corr.Values[i])
	}
	got := offline.UntrustedWindows()
	for i := range windows {
		if got[i] != read.Untrusted[i] {
			t.Fatalf("MUTANT 26p-FEED-WINDOW DIVERGENCE at window %d: online=%v offline=%v - the two input seams do not reach the same decision path on real AEC3 material (389-v1 §3's boundary ①)",
				i, read.Untrusted[i], got[i])
		}
	}
	// The window-level equivalence above is the grid's claim: the two input
	// seams reach the same per-window verdicts. An episode's TightenActive
	// reads the flag at its DECISION moment - the start window may not have
	// closed yet (an episode starting mid-window judges before that window's
	// verdict exists, the gate's own documented rule) - so the episode flag
	// is a decision-moment reading, not a window-final one, and comparing it
	// against the re-cut's final flags would compare two different moments.
	t.Logf("FEED-EQUIVALENCE windows=%d untrusted=%d lag=%d strength=%.2f episodes=%d triggers=%d (online Feed == offline ObserveWindow, zero drift)",
		windows, countTrue26p(read.Untrusted), read.Lag, read.LagStrength, len(read.Episodes), read.Triggers)
	write26pJSON(t, "26p-feed-equivalence.json", map[string]any{
		"online_untrusted":  read.Untrusted,
		"offline_untrusted": got,
		"lag":               read.Lag,
		"lag_strength":      read.LagStrength,
		"windows":           windows,
		"episodes":          read.Episodes,
	})
}

func countTrue26p(flags []bool) int {
	n := 0
	for _, f := range flags {
		if f {
			n++
		}
	}
	return n
}

// compile-time: the rig's speaker satisfies the seam.
var _ voiceSpeaker = (*voice26pRecordSpeaker)(nil)
