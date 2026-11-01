//go:build windows

package main

// Ticket 392's integration readings: the keepalive posture question (the
// 28-r1 warm-pool shape on the barge-in monitor), the AEC3 double-talk
// tuning face, and their reverses - the surfaces the 26p-r1 design inputs
// assigned to the 28 session-integration leg:
//
//	AC#1 保温：the bare warm posture was tried and measured UNSOUND on this
//	     surface - the second warm window's echo cancellation failed even
//	     with a fresh APM and a fresh measurement (the readings carry the
//	     profiles), so production runs the exact pre-392 byte shape
//	     (keepalive off). The 397 between-window clears revive it: each
//	     held engine walks its public clear at release (APM Reset /
//	     resampler Flush / detector Finalize plus the 1.2s cool segment,
//	     boundary events discarded), and this test pins both halves - the
//	     cold posture's multi-window (a)/(b) shapes hold, and the cleared
//	     warm posture's second window judges again (episodes=1,
//	     triggers=1@talker, the 392-r2 revival-cool readings). The posture
//	     retires with the voice link either way.
//	AC#2 调参：the double-talk damage swing the 26p-r1 L1 reading
//	     carried (-29.4/-38.2 across two firings, 0-8dB with the engine
//	     state) re-measured with K firings per setting: the engine
//	     default, and one AEC3 suppressor parameter per firing set (the
//	     binding layer's EchoControlFactory seam - the webrtc tree stays
//	     pinned and untouched). Each setting's swing width is the
//	     reading; if nothing flattens it, the widths say by how much -
//	     the design input's refinement.
//	AC#3 反向：the cold posture is the default (the pre-392 shape); the
//	     tuning face off -> the default engine returns (AecTuned=false);
//	     dispose -> the engines go down. Each with its restore core.
//
// Staging: part 1's header (the 26p rig's shape) - the same staged binary
// this package builds under -tags webrtcap, the 15e/A1326 staging around it.

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/audio/aec"
)

// voice392Out is this leg's artifact root (the readings JSONs land here).
const voice392Out = `D:\tmp\wisp\20261031\392r1\out`

// voice392EarlyTalkerAt is when the early talker steps in - inside the
// alignment span the cold judge waits out, the surface where the warm
// posture's judging failure shows.
const voice392EarlyTalkerAt = 400 * time.Millisecond

// voice392Rounds is the firings per tuning setting: the swing is the
// phenomenon, so one firing per setting cannot separate a narrowing from
// jitter - each setting gets this many L1-shape rounds and the artifact
// carries every one of them.
const voice392Rounds = 4

// voice392ColdReadyFloor is the readiness a fresh-measuring judge must show:
// the cold arm waits out the alignment spans (300ms of far head plus the
// 300ms search span - 9600 capture samples past the anchor), so no window
// arms its judge before ~9600 samples.
const voice392ColdReadyFloor = 9000

// write392JSON drops one readings artifact (the evidence table quotes these).
func write392JSON(t *testing.T, name string, v any) {
	t.Helper()
	if err := os.MkdirAll(voice392Out, 0o755); err != nil {
		t.Fatalf("out dir: %v", err)
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatalf("marshal %s: %v", name, err)
	}
	p := filepath.Join(voice392Out, name)
	if err := os.WriteFile(p, append(b, '\n'), 0o644); err != nil {
		t.Fatalf("write %s: %v", p, err)
	}
	t.Logf("WROTE %s", p)
}

// voice392RmsProfile buckets one recorded stream's RMS per second - the
// residual's shape over the window (a healthy cancellation keeps the
// residual well below the near; the warm posture's finding shows as the
// residual riding at the near level).
func voice392RmsProfile(x []float32, buckets int) []float64 {
	if len(x) == 0 || buckets <= 0 {
		return nil
	}
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
		var sum float64
		for _, v := range x[lo:hi] {
			sum += float64(v) * float64(v)
		}
		rms := math.Sqrt(sum / float64(hi-lo))
		dbfs := math.Inf(-1)
		if rms > 0 {
			dbfs = math.Round(20*math.Log10(rms)*10) / 10
		}
		out[b] = dbfs
	}
	return out
}

// voice392Round is one announcement round's reading: the loudest episode's
// level (the talker's - the residual floor opens quieter ones), the engine
// facts the keepalive and tuning readings key on, and the near/residual RMS
// profiles (1s buckets over the recorded streams, the honest shape of what
// the window's engine produced).
type voice392Round struct {
	round        int
	rms          float64
	threshold    float64
	episodes     int
	triggers     int
	engineBuilt  bool
	engineWarm   bool
	aecTuned     bool
	readySamples int
	lag          int
	nearProfile  []float64
	resProfile   []float64
}

// voice392RoundOut is the round reading's exported artifact shape - the keys
// the evidence quotes. 392-r1's artifact serialized every round as an empty
// object because the working struct's fields are unexported (392-v3 §1 named
// the shell); every reading this test pins lands on disk through this shape.
type voice392RoundOut struct {
	Round        int       `json:"round"`
	Rms          float64   `json:"rms"`
	Threshold    float64   `json:"threshold"`
	Episodes     int       `json:"episodes"`
	Triggers     int       `json:"triggers"`
	EngineBuilt  bool      `json:"engine_built"`
	EngineWarm   bool      `json:"engine_warm"`
	AecTuned     bool      `json:"aec_tuned"`
	ReadySamples int       `json:"ready_samples"`
	Lag          int       `json:"lag"`
	NearProfile  []float64 `json:"near_profile"`
	ResProfile   []float64 `json:"res_profile"`
}

// voice392RoundOutFrom copies one round reading into the artifact shape.
func voice392RoundOutFrom(r voice392Round) voice392RoundOut {
	return voice392RoundOut{
		Round: r.round, Rms: r.rms, Threshold: r.threshold,
		Episodes: r.episodes, Triggers: r.triggers,
		EngineBuilt: r.engineBuilt, EngineWarm: r.engineWarm,
		AecTuned: r.aecTuned, ReadySamples: r.readySamples, Lag: r.lag,
		NearProfile: r.nearProfile, ResProfile: r.resProfile,
	}
}

// voice392RoundsOut copies a round sequence into the artifact shape.
func voice392RoundsOut(rs []voice392Round) []voice392RoundOut {
	out := make([]voice392RoundOut, 0, len(rs))
	for _, r := range rs {
		out = append(out, voice392RoundOutFrom(r))
	}
	return out
}

// voice392Capture rounds one announcement (sentence + matched echo + an
// optional talker at talkerAt; a zero talkerAt runs the (b) shape), waits
// for the window to close, and reads the judge with the stream profiles.
func voice392Capture(t *testing.T, rig *voice26pRig, sentence, echoWav, talkerWav string, talkerAt time.Duration, round int) voice392Round {
	t.Helper()
	aborted0, spoken0 := rig.rv.ttsAborted.Load(), rig.rv.ttsSpoken.Load()
	walkToThinking381r2(t, rig.rv)
	if !rig.rv.announce(sentence) {
		t.Fatalf("round %d: the announce refused the sentence; problems=%v", round, rig.rv.problemsSnapshot())
	}
	spy := rig.spyAt(len(rig.spies))
	<-spy.first
	stopEcho, _ := rig.inject(echoWav)
	var stopTalker func()
	if talkerWav != "" {
		time.Sleep(talkerAt)
		stopTalker, _ = rig.inject(talkerWav)
	}
	voice15eWaitFor(t, fmt.Sprintf("round %d: the announcement to finish (spoken or interrupted)", round), func() bool {
		return rig.rv.ttsAborted.Load() > aborted0 || rig.rv.ttsSpoken.Load() > spoken0
	})
	if stopTalker != nil {
		stopTalker()
	}
	stopEcho()
	rig.waitReleased()
	read := rig.rv.bargein.readings()
	out := voice392Round{
		round:        round,
		episodes:     len(read.Episodes),
		triggers:     read.Triggers,
		engineBuilt:  read.EngineBuilt,
		engineWarm:   read.EngineWarm,
		aecTuned:     read.AecTuned,
		readySamples: read.ReadySamples,
		lag:          read.Lag,
		rms:          math.Inf(-1),
	}
	for i := range read.Episodes {
		if read.Episodes[i].RmsDbfs > out.rms {
			out.rms = read.Episodes[i].RmsDbfs
			out.threshold = read.Episodes[i].ThresholdDbfs
		}
	}
	if near, _, res := rig.rv.bargein.streams(); len(near) > 0 {
		buckets := int(float64(len(near))/float64(audio.TargetRate)) + 1
		out.nearProfile = voice392RmsProfile(near, buckets)
		out.resProfile = voice392RmsProfile(res, buckets)
		t.Logf("ROUND-%d-PROFILE near_dbfs=%v res_dbfs=%v", round, out.nearProfile, out.resProfile)
	}
	return out
}

// Test392KeepalivePostureFindingsAndColdDefault is the AC#1 reading. Part A
// (the production default, the exact pre-392 byte shape): a multi-window
// sequence - control, two talker rounds, an echo-only round - keeps the
// (a)/(b) shapes healthy window after window with the engines rebuilt per
// window. Part B (the warm posture, the tests' switch): the same sequence on
// held engines - the first warm window judges, and the second one is the
// 397 revival reading: the held engines walked their between-window clears
// (APM Reset, resampler Flush, detector Finalize plus the 1.2s cool segment,
// boundary events discarded), so the second window judges again - episodes=1,
// triggers=1@talker, the residual at the cancelled depth (the 392-r2
// revival-cool shape, first read in the aecspike replica and pinned here on
// the staged chain). The restore core: the switch back off, the (b) shape
// holds again. The default stays off - the clears revive the posture, the
// default flip is a separate ruling that re-measures the full staged chain.
func Test392KeepalivePostureFindingsAndColdDefault(t *testing.T) {
	if probe, err := openWaveOutSpeaker(22050); err != nil {
		t.Skipf("no output device on this machine: %v", err)
	} else {
		_ = probe.Close()
	}
	dir := t.TempDir()
	talkerWav := voice26pRenderTtsWav(t, dir, "talker-392", voice26pTalker, voice26pTalkerDbfs)

	rig := newVoice26pRig(t, nil)
	rig.requireAEC3()
	if rig.rv.bargein.readings().Keepalive {
		t.Fatalf("the monitor's production default is the warm posture - the measured finding keeps it off")
	}

	// Part A, round 0: the echo-source control. No injection - the rig has
	// no capture device, so the drain sees no frames and the judge never
	// arms here; the window builds the engines (rebuilt at the next arm -
	// the cold posture) and records the announcement.
	walkToThinking381r2(t, rig.rv)
	if !rig.rv.announce(voice26pAnnounceA) {
		t.Fatalf("the control refused the sentence; problems=%v", rig.rv.problemsSnapshot())
	}
	control := rig.spyAt(0)
	voice15eWaitFor(t, "the control to play out", func() bool {
		return rig.rv.ttsSpoken.Load() >= 1
	})
	rig.waitReleased()
	ctrl := rig.rv.bargein.readings()
	if ctrl.Keepalive || !ctrl.EngineBuilt || ctrl.EngineWarm {
		t.Fatalf("the control is not the cold posture: keepalive=%v built=%v warm=%v",
			ctrl.Keepalive, ctrl.EngineBuilt, ctrl.EngineWarm)
	}
	echoWav, _ := voice26pEchoWav(t, dir, "echo-392", control.recording(), 22050, voice26pEchoGain)

	// Part A, rounds 1-2: the (a) shape on two consecutive windows - the -20
	// talker judged through the fresh-measuring judge both times, engines
	// rebuilt per window (the pre-392 shape holds multi-window).
	var partA []voice392Round
	for r := 1; r <= 2; r++ {
		rd := voice392Capture(t, rig, voice26pAnnounceA, echoWav, talkerWav, voice392EarlyTalkerAt, r)
		if rd.triggers != 1 {
			t.Fatalf("part A round %d triggered %d time(s) (episodes=%d): the -20 talker the (a) surface passes was not judged",
				r, rd.triggers, rd.episodes)
		}
		if !rd.engineBuilt || rd.engineWarm {
			t.Fatalf("part A round %d did not rebuild per window: built=%v warm=%v", r, rd.engineBuilt, rd.engineWarm)
		}
		if rd.readySamples < voice392ColdReadyFloor {
			t.Fatalf("part A round %d's judge did not measure fresh (ready=%d samples)", r, rd.readySamples)
		}
		t.Logf("COLD-MULTIWINDOW round=%d triggers=1 ready_samples=%d lag=%d (the (a) shape holds window after window)",
			r, rd.readySamples, rd.lag)
		partA = append(partA, rd)
	}

	// Part A, round 3: the (b) shape on the same posture - the echo-only
	// window holds the tightened cut.
	bCold := voice392Capture(t, rig, voice26pAnnounceA, echoWav, "", 0, 3)
	if bCold.triggers != 0 {
		t.Fatalf("part A round 3's echo self-triggered (triggers=%d): the cold posture does not hold the (b) shape", bCold.triggers)
	}
	for _, ep := range rig.rv.bargein.readings().Episodes {
		if ep.EnergyPass {
			t.Fatalf("part A round 3 let an echo episode through the tightened cut (%.2f dBFS)", ep.RmsDbfs)
		}
	}
	t.Logf("COLD-MULTIWINDOW echo-only triggers=0 episodes=%d all_below_cut=true (the (b) shape holds)", bCold.episodes)

	// Part B: the warm posture on. The engines the last cold window released
	// are rebuilt here and then HELD - the 28-r1 warm shape - and the first
	// warm window judges fine.
	rig.rv.bargein.setKeepalive(true)
	warm1 := voice392Capture(t, rig, voice26pAnnounceA, echoWav, talkerWav, voice392EarlyTalkerAt, 4)
	if warm1.triggers != 1 {
		t.Fatalf("the first warm window triggered %d time(s) (episodes=%d): the -20 talker was not judged",
			warm1.triggers, warm1.episodes)
	}
	if !warm1.engineBuilt || warm1.engineWarm {
		t.Fatalf("the first warm window did not build fresh (the cold round closed the engines): built=%v warm=%v",
			warm1.engineBuilt, warm1.engineWarm)
	}
	t.Logf("WARM-POSTURE round=1 triggers=1 (the first warm window judges)")

	// The revival (ticket 397): the second warm window judges again. The held
	// engines walked their between-window clears at the last release - the
	// APM's Reset, the resampler's Flush, the detector's Finalize plus the
	// 1.2s cool segment whose boundary events are discarded - so the -20
	// talker fires from the held pool exactly as the 392-r2 revival-cool
	// readings pinned: episodes=1, triggers=1@talker, the residual at the
	// cancelled depth (the profiles in the artifact). The clears are jointly
	// load-bearing: removing any one regresses its own directed failure
	// shape (the AC#3 overlay firings name them).
	warm2 := voice392Capture(t, rig, voice26pAnnounceA, echoWav, talkerWav, voice392EarlyTalkerAt, 5)
	if !warm2.engineWarm || warm2.engineBuilt {
		t.Fatalf("the second warm window rebuilt the held engines: built=%v warm=%v", warm2.engineBuilt, warm2.engineWarm)
	}
	if warm2.episodes != 1 || warm2.triggers != 1 {
		t.Fatalf("WARM REVIVAL MOVED: the second warm window read (episodes=%d triggers=%d), want episodes=1 triggers=1@talker (the 392-r2 revival-cool readings) - the between-window clears did not revive the judge or moved past them",
			warm2.episodes, warm2.triggers)
	}
	t.Logf("WARM-POSTURE round=2 triggers=1 episodes=1 (the 397 revival: the held engines judge the second window again; the profiles are in the artifact)")

	// The restore core: the switch back off, the (b) shape holds again.
	rig.rv.bargein.setKeepalive(false)
	bRestore := voice392Capture(t, rig, voice26pAnnounceA, echoWav, "", 0, 6)
	if bRestore.triggers != 0 {
		t.Fatalf("the restore round's echo self-triggered (triggers=%d)", bRestore.triggers)
	}
	for _, ep := range rig.rv.bargein.readings().Episodes {
		if ep.EnergyPass {
			t.Fatalf("the restore round let an echo episode through the tightened cut (%.2f dBFS)", ep.RmsDbfs)
		}
	}
	if !bRestore.engineBuilt {
		t.Fatalf("the restore round did not rebuild (the cold posture is not back): built=%v", bRestore.engineBuilt)
	}
	t.Logf("COLD-RESTORE triggers=0 episodes=%d all_below_cut=true (the pair closed)", bRestore.episodes)
	write392JSON(t, "392-keepalive.json", map[string]any{
		"cold_control": ctrl,
		"cold_rounds":  voice392RoundsOut(partA),
		"cold_echo":    voice392RoundOutFrom(bCold),
		"warm_round1":  voice392RoundOutFrom(warm1),
		"warm_round2":  voice392RoundOutFrom(warm2),
		"cold_restore": voice392RoundOutFrom(bRestore),
		"talker_at":    voice392EarlyTalkerAt.String(),
		"talker_level": voice26pTalkerDbfs,
		"revival":      "the 397 between-window clears revived the warm posture on the staged chain: the second warm window judges again (episodes=1, triggers=1@talker, the residual at the cancelled depth) - the 392-r2 revival-cool shape; the clears are jointly load-bearing (the AC#3 overlay firings), and the production default stays keepalive=false pending the separate default-flip ruling",
	})
}

// Test392AecTuningDoubleTalkSwing is the AC#2 reading: the double-talk
// damage swing re-measured with K firings per setting on the L1 shape (the
// -26.70 control over the matched echo, the 26p-r1 margin re-test's). The
// settings: the engine default, and one AEC3 suppressor parameter per cold
// firing set - the binding layer's parameter face, one parameter at a time.
// The swing widths are the reading; nothing here asserts a flattening - if
// a parameter does not narrow the swing, the widths say by how much, and
// that is the design input's refinement.
func Test392AecTuningDoubleTalkSwing(t *testing.T) {
	if probe, err := openWaveOutSpeaker(22050); err != nil {
		t.Skipf("no output device on this machine: %v", err)
	} else {
		_ = probe.Close()
	}
	dir := t.TempDir()
	talkerL1 := voice26pRenderTtsWav(t, dir, "talker-l1-392", voice26pTalker, -26.70)

	rig := newVoice26pRig(t, nil)
	rig.requireAEC3()

	// The matched echo source: one full announcement, recorded.
	walkToThinking381r2(t, rig.rv)
	if !rig.rv.announce(voice26pAnnounceA) {
		t.Fatalf("the echo-source control refused the sentence; problems=%v", rig.rv.problemsSnapshot())
	}
	ctrl := rig.spyAt(0)
	voice15eWaitFor(t, "the echo-source control to play out", func() bool {
		return rig.rv.ttsSpoken.Load() >= 1
	})
	rig.waitReleased()
	echoMatch, _ := voice26pEchoWav(t, dir, "echo-match-392", ctrl.recording(), 22050, voice26pEchoGain)

	type setting struct {
		name   string
		tuning aec.Tuning
	}
	settings := []setting{
		{name: "baseline-default"},
		{name: "enr-threshold-0.45", tuning: aec.Tuning{DominantNearendEnrThreshold: 0.45}},
		{name: "trigger-threshold-6", tuning: aec.Tuning{DominantNearendTriggerThreshold: 6}},
		{name: "nearend-max-dec-lf-0.5", tuning: aec.Tuning{NearendMaxDecFactorLf: 0.5}},
	}
	type settingReadings struct {
		Name    string             `json:"name"`
		Tuning  aec.Tuning         `json:"tuning"`
		Rounds  []voice392RoundOut `json:"rounds"`
		MinRms  float64            `json:"min_rms"`
		MaxRms  float64            `json:"max_rms"`
		SwingDb float64            `json:"swing_db"`
	}
	out := make([]settingReadings, 0, len(settings)+1)
	for _, s := range settings {
		rig.rv.bargein.setAECTuning(s.tuning)
		sr := settingReadings{Name: s.name, Tuning: s.tuning}
		first := true
		for r := 1; r <= voice392Rounds; r++ {
			round := voice392Capture(t, rig, voice26pAnnounceA, echoMatch, talkerL1, voice26pTalkerAt, r)
			if round.episodes == 0 || round.rms <= voice26pLevelFloor {
				t.Fatalf("%s round %d: no episode level to read (episodes=%d rms=%.2f) - the L1 shape did not reach the judge",
					s.name, r, round.episodes, round.rms)
			}
			t.Logf("TUNING-%s round=%d rms=%.2f threshold=%.1f episodes=%d triggers=%d tuned=%v ready=%d",
				s.name, r, round.rms, round.threshold, round.episodes, round.triggers,
				round.aecTuned, round.readySamples)
			sr.Rounds = append(sr.Rounds, voice392RoundOutFrom(round))
			if first {
				sr.MinRms, sr.MaxRms = round.rms, round.rms
				first = false
			} else {
				sr.MinRms = min(sr.MinRms, round.rms)
				sr.MaxRms = max(sr.MaxRms, round.rms)
			}
		}
		sr.SwingDb = sr.MaxRms - sr.MinRms
		t.Logf("TUNING-%s swing=%.2fdB (min=%.2f max=%.2f)", s.name, sr.SwingDb, sr.MinRms, sr.MaxRms)
		out = append(out, sr)
	}

	// The restore core: the tuning face off - the default engine is what the
	// next firing builds (AecTuned=false in its reading), and its level is
	// logged beside the baseline distribution.
	rig.rv.bargein.setAECTuning(aec.Tuning{})
	restore := voice392Capture(t, rig, voice26pAnnounceA, echoMatch, talkerL1, voice26pTalkerAt, voice392Rounds+1)
	if restore.aecTuned {
		t.Fatalf("the restore round still reads the tuning face as active: the default engine did not come back")
	}
	if restore.episodes == 0 || restore.rms <= voice26pLevelFloor {
		t.Fatalf("the restore round has no level to read (episodes=%d rms=%.2f)", restore.episodes, restore.rms)
	}
	t.Logf("TUNING-RESTORE tuned=false rms=%.2f (the default engine back, beside the baseline %.2f..%.2f)",
		restore.rms, out[0].MinRms, out[0].MaxRms)
	write392JSON(t, "392-aec-tuning.json", map[string]any{
		"settings":  out,
		"restore":   voice392RoundOutFrom(restore),
		"talker":    "control-26 (-26.70 dBFS)",
		"echo":      "matched, coupling 0.2",
		"talker_at": voice26pTalkerAt.String(),
	})
}

// Test392DisposeReleasesTheEngines is the third reverse: the engines retire
// with the voice link, on the C11 cascade's release step. With the warm
// posture holding them across windows, stop() must take them down; the
// readings say EnginesHot=false after it.
func Test392DisposeReleasesTheEngines(t *testing.T) {
	if probe, err := openWaveOutSpeaker(22050); err != nil {
		t.Skipf("no output device on this machine: %v", err)
	} else {
		_ = probe.Close()
	}
	rig := newVoice26pRig(t, nil)
	rig.requireAEC3()
	rig.rv.bargein.setKeepalive(true)

	walkToThinking381r2(t, rig.rv)
	if !rig.rv.announce(voice26pAnnounceA) {
		t.Fatalf("the control refused the sentence; problems=%v", rig.rv.problemsSnapshot())
	}
	voice15eWaitFor(t, "the control to play out", func() bool {
		return rig.rv.ttsSpoken.Load() >= 1
	})
	rig.waitReleased()
	if hot := rig.rv.bargein.readings().EnginesHot; !hot {
		t.Fatalf("the engines are not held after the window (EnginesHot=false): the warm posture did not hold them")
	}
	t.Logf("DISPOSE engines_hot=true (the warm posture holds the model engine between windows)")
	rig.rv.stop()
	if hot := rig.rv.bargein.readings().EnginesHot; hot {
		t.Fatalf("the engines are still hot after the voice link stopped - the C11 release step did not reach them")
	}
	t.Logf("DISPOSE engines_hot=false (the link's release step took them down; the cascade closed)")
}
