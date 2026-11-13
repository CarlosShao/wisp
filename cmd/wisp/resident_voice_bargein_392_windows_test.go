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
//	     profiles). The 397 between-window clears revive it: each held
//	     engine walks its public clear at release (APM Reset / resampler
//	     Flush / detector Finalize plus the 1.2s cool segment, boundary
//	     events discarded), and the 401 flip made the CLEARED posture the
//	     production default - this test pins the default path (the held
//	     windows judge again window after window: episodes=1,
//	     triggers=1@talker, the 392-r2 revival-cool readings) beside the
//	     switch-off cold face (the exact pre-392 shapes). The posture
//	     retires with the voice link either way.
//	AC#2 调参：the double-talk damage swing the 26p-r1 L1 reading
//	     carried (-29.4/-38.2 across two firings, 0-8dB with the engine
//	     state) re-measured with K firings per setting: the engine
//	     default, and one AEC3 suppressor parameter per firing set (the
//	     binding layer's EchoControlFactory seam - the webrtc tree stays
//	     pinned and untouched). Each setting's swing width is the
//	     reading; if nothing flattens it, the widths say by how much -
//	     the design input's refinement.
//	AC#3 反向：the cold posture is the switch-off face (the pre-392 shape);
//	     the tuning face off -> the default engine returns (AecTuned=false);
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
	"reflect"
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
// Non-finite readings are clamped first: JSON has no Infinity, so one -Inf
// anywhere in the artifact fails the whole marshal and the readings never
// land (397-v2's residual 2 - a fully gate-clean (b) round red the run's
// tail with every behavioral assertion passed).
func write392JSON(t *testing.T, name string, v any) {
	t.Helper()
	if err := os.MkdirAll(voice392Out, 0o755); err != nil {
		t.Fatalf("out dir: %v", err)
	}
	safe, clamped := voice392JSONSafe(v)
	if clamped > 0 {
		t.Logf("JSON-SENTINEL %s: clamped %d non-finite reading(s) to %.0f dBFS (the gate-clean -Inf shape)", name, clamped, voice392NoSignalDbfs)
	}
	b, err := json.MarshalIndent(safe, "", "  ")
	if err != nil {
		t.Fatalf("marshal %s: %v", name, err)
	}
	p := filepath.Join(voice392Out, name)
	if err := os.WriteFile(p, append(b, '\n'), 0o644); err != nil {
		t.Fatalf("write %s: %v", p, err)
	}
	t.Logf("WROTE %s", p)
}

// voice392NoSignalDbfs is the artifact value for a non-finite reading. -Inf
// here means "no signal" - the no-episode round's rms sentinel and a
// pure-silence profile bucket are digital silence, not measurements - and
// -999 sits far below every real reading this suite produces (the quietest
// band bottoms out near -59) while staying a number the evidence table can
// quote. A finite clamp (not null) keeps the artifact schema untouched:
// float fields stay floats, a finite artifact marshals bit-for-bit as
// before, and the sentinel only fires where -Inf did. Both signs and NaN
// take the same value - the artifact's job is to say "this reading does not
// exist", not to guess a magnitude, and 0 would misread as full scale.
const voice392NoSignalDbfs = -999.0

// voice392JSONSafe returns v with every non-finite float clamped to
// voice392NoSignalDbfs, and reports how many readings it clamped. The walk
// covers maps, slices, pointers, interfaces, and exported struct fields -
// both artifact shapes route through here, plus the raw control readings.
// Containers are copied only when something inside them was clamped; a
// value with nothing to clamp comes back as-is.
func voice392JSONSafe(v any) (any, int) {
	out, clamped := voice392JSONSafeWalk(reflect.ValueOf(v))
	if clamped == 0 {
		return v, 0
	}
	return out.Interface(), clamped
}

func voice392JSONSafeWalk(rv reflect.Value) (reflect.Value, int) {
	switch rv.Kind() {
	case reflect.Float32, reflect.Float64:
		if f := rv.Float(); math.IsInf(f, 0) || math.IsNaN(f) {
			return reflect.ValueOf(voice392NoSignalDbfs), 1
		}
	case reflect.Interface, reflect.Pointer:
		if !rv.IsNil() {
			return voice392JSONSafeWalk(rv.Elem())
		}
	case reflect.Map:
		clamped := 0
		out := reflect.MakeMapWithSize(rv.Type(), rv.Len())
		iter := rv.MapRange()
		for iter.Next() {
			cv, n := voice392JSONSafeWalk(iter.Value())
			clamped += n
			out.SetMapIndex(iter.Key(), cv)
		}
		if clamped > 0 {
			return out, clamped
		}
	case reflect.Slice:
		clamped := 0
		out := reflect.MakeSlice(rv.Type(), rv.Len(), rv.Len())
		for i := 0; i < rv.Len(); i++ {
			cv, n := voice392JSONSafeWalk(rv.Index(i))
			clamped += n
			out.Index(i).Set(cv)
		}
		if clamped > 0 {
			return out, clamped
		}
	case reflect.Struct:
		clamped := 0
		out := reflect.New(rv.Type()).Elem()
		for i := 0; i < rv.NumField(); i++ {
			if !rv.Type().Field(i).IsExported() {
				continue
			}
			cv, n := voice392JSONSafeWalk(rv.Field(i))
			clamped += n
			out.Field(i).Set(cv)
		}
		if clamped > 0 {
			return out, clamped
		}
	}
	return rv, 0
}

// Test392JSONSentinelClampsNonFiniteReadings pins the writer's JSON safety
// sentinel at the exact face 397-v2's residual 2 measured: a fully gate-clean
// (b) round (episodes=0) leaves the round reading's rms at the capture's -Inf
// sentinel and a pure-silence profile bucket at -Inf, and one -Inf failed the
// whole artifact marshal (the run red at the tail with every behavioral
// assertion passed). The sentinel must also stay silent on finite readings -
// a finite artifact marshals byte-for-byte as before.
func Test392JSONSentinelClampsNonFiniteReadings(t *testing.T) {
	finite := voice392RoundOutFrom(voice392Round{
		round: 3, rms: -47.2, threshold: -30,
		episodes: 1, triggers: 0, engineBuilt: true,
		readySamples: voice392ColdReadyFloor, lag: 490,
		nearProfile: []float64{-22.6, -24.3},
		resProfile:  []float64{-47.2, -58.9},
	})
	safe, clamped := voice392JSONSafe(finite)
	if clamped != 0 {
		t.Fatalf("the finite reading was clamped %d time(s): the sentinel must stay silent on finite values", clamped)
	}
	before, err := json.Marshal(finite)
	if err != nil {
		t.Fatalf("marshal finite: %v", err)
	}
	after, err := json.Marshal(safe)
	if err != nil {
		t.Fatalf("marshal finite through the sentinel: %v", err)
	}
	if string(before) != string(after) {
		t.Fatalf("the finite artifact changed through the sentinel:\nbefore:  %s\nafter:   %s", before, after)
	}

	// The (b) artifact face: no episode fired, so rms keeps the value
	// voice392Capture leaves (its -Inf sentinel), and the residual's first
	// second is digital silence (its profile bucket is -Inf the same way).
	gateClean := voice392RoundOutFrom(voice392Round{
		round: 3, rms: math.Inf(-1),
		episodes: 0, triggers: 0, engineBuilt: true,
		readySamples: voice392ColdReadyFloor, lag: 536,
		nearProfile: []float64{-22.6, -24.3},
		resProfile:  []float64{math.Inf(-1), -58.9},
	})
	safe, clamped = voice392JSONSafe(gateClean)
	if clamped != 2 {
		t.Fatalf("clamped=%d, want 2 (the no-episode rms plus the silent bucket)", clamped)
	}
	b, err := json.Marshal(safe)
	if err != nil {
		t.Fatalf("marshal gate-clean through the sentinel: %v", err)
	}
	var back map[string]any
	if err := json.Unmarshal(b, &back); err != nil {
		t.Fatalf("the clamped artifact is not valid JSON: %v", err)
	}
	if got := fmt.Sprintf("%v", back["rms"]); got != "-999" {
		t.Fatalf("the gate-clean rms reads %s, want the no-signal sentinel -999", got)
	}
	profile, ok := back["res_profile"].([]any)
	if !ok || len(profile) != 2 || fmt.Sprintf("%v", profile[0]) != "-999" {
		t.Fatalf("the silent profile bucket reads %v, want the no-signal sentinel -999 in slot 0", back["res_profile"])
	}
	write392JSON(t, "392-json-sentinel.json", gateClean)
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
// (the production default since the 401 flip): a multi-window sequence -
// control, two talker rounds, an echo-only round - on the held engines: the
// between-window clears (APM Reset, resampler Flush, detector Finalize plus
// the 1.2s cool segment, boundary events discarded) put the engines back to
// their just-created state at every release, so the talker windows judge
// again window after window - episodes=1, triggers=1@talker, the residual at
// the cancelled depth (the 392-r2 revival-cool shape, first read in the
// aecspike replica, pinned here on the staged chain, now live under the
// default with no switch) and the echo-only window holds (no false-start
// leak - the 398-v3 decay curve's clearing). Part B (the switch-off face,
// the exact pre-392 shape): the cold posture rebuilds per window and keeps
// the (a)/(b) shapes. The restore core: the switch back on, the (b) shape
// holds on the default's fresh build again.
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
	if ka := rig.rv.bargein.readings().Keepalive; !ka {
		t.Fatalf("the monitor's default is not the warm posture - the 401 flip made the cleared posture the default, got keepalive=%v", ka)
	}

	// Part A, round 0: the echo-source control. No injection - the rig has
	// no capture device, so the drain sees no frames and the judge never
	// arms here; the window builds the engines (the default's first build)
	// and records the announcement.
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
	if !ctrl.Keepalive || !ctrl.EngineBuilt || ctrl.EngineWarm {
		t.Fatalf("the control is not the default warm posture's first build: keepalive=%v built=%v warm=%v",
			ctrl.Keepalive, ctrl.EngineBuilt, ctrl.EngineWarm)
	}
	if hot := rig.rv.bargein.readings().EnginesHot; !hot {
		t.Fatalf("the default posture did not hold the engines past the control window (EnginesHot=false)")
	}
	echoWav, _ := voice26pEchoWav(t, dir, "echo-392", control.recording(), 22050, voice26pEchoGain)

	// Part A, rounds 1-2: the (a) shape on the held engines - the between-
	// window clears ran at the previous release, so the judge measures fresh
	// and the -20 talker is judged window after window with nothing rebuilt:
	// the revival is live under the DEFAULT (the WARM REVIVAL MOVED guard
	// runs here without any switch - the 401 flip's pin).
	var warmRounds []voice392Round
	for r := 1; r <= 2; r++ {
		rd := voice392Capture(t, rig, voice26pAnnounceA, echoWav, talkerWav, voice392EarlyTalkerAt, r)
		if rd.episodes != 1 || rd.triggers != 1 {
			t.Fatalf("WARM REVIVAL MOVED: the held window round %d read (episodes=%d triggers=%d), want episodes=1 triggers=1@talker (the 392-r2 revival-cool readings, live under the default) - the between-window clears did not revive the judge or moved past them",
				r, rd.episodes, rd.triggers)
		}
		if rd.engineBuilt || !rd.engineWarm {
			t.Fatalf("default posture round %d rebuilt instead of judging on the held engines: built=%v warm=%v", r, rd.engineBuilt, rd.engineWarm)
		}
		if rd.readySamples < voice392ColdReadyFloor {
			t.Fatalf("default posture round %d's judge did not measure fresh (ready=%d samples)", r, rd.readySamples)
		}
		t.Logf("WARM-DEFAULT round=%d triggers=1 episodes=1 ready_samples=%d lag=%d (the revival holds window after window on the held engines)",
			r, rd.readySamples, rd.lag)
		warmRounds = append(warmRounds, rd)
	}

	// Part A, round 3: the (b) shape on the held engines - the echo-only
	// window holds the tightened cut, and no false start leaks through the
	// clears (the 398-v3 decay curve: the Finalize closes the open utterance
	// and the cool segment absorbs the hangover - the first fed silence eats
	// any false start and the discarded events never reach a reading).
	bWarm := voice392Capture(t, rig, voice26pAnnounceA, echoWav, "", 0, 3)
	if bWarm.triggers != 0 {
		t.Fatalf("the default posture's echo-only window self-triggered (triggers=%d): the clears did not hold the (b) shape on the held engines", bWarm.triggers)
	}
	if bWarm.engineBuilt || !bWarm.engineWarm {
		t.Fatalf("the default posture's echo-only window rebuilt instead of judging on the held engines: built=%v warm=%v", bWarm.engineBuilt, bWarm.engineWarm)
	}
	for _, ep := range rig.rv.bargein.readings().Episodes {
		if ep.EnergyPass {
			t.Fatalf("the default posture's echo-only window let an echo episode through the tightened cut (%.2f dBFS)", ep.RmsDbfs)
		}
	}
	t.Logf("WARM-DEFAULT echo-only triggers=0 episodes=%d all_below_cut=true (the (b) shape holds on the held engines; no false-start leak)", bWarm.episodes)

	// Part B: the switch-off face - the exact pre-392 shape. The cold
	// posture rebuilds per window and keeps the (a)/(b) shapes on the same
	// material (the pre-392 default's own reading, still the reverse gear).
	rig.rv.bargein.setKeepalive(false)
	aCold := voice392Capture(t, rig, voice26pAnnounceA, echoWav, talkerWav, voice392EarlyTalkerAt, 4)
	if aCold.triggers != 1 {
		t.Fatalf("the cold face's talker round triggered %d time(s) (episodes=%d): the (a) shape does not hold on the rebuilt engines",
			aCold.triggers, aCold.episodes)
	}
	if !aCold.engineBuilt || aCold.engineWarm {
		t.Fatalf("the cold face did not rebuild per window: built=%v warm=%v", aCold.engineBuilt, aCold.engineWarm)
	}
	if aCold.readySamples < voice392ColdReadyFloor {
		t.Fatalf("the cold face's judge did not measure fresh (ready=%d samples)", aCold.readySamples)
	}
	t.Logf("COLD-FACE round=4 triggers=1 ready_samples=%d lag=%d (the pre-392 shape holds on the switch-off face)",
		aCold.readySamples, aCold.lag)
	bCold := voice392Capture(t, rig, voice26pAnnounceA, echoWav, "", 0, 5)
	if bCold.triggers != 0 {
		t.Fatalf("the cold face's echo self-triggered (triggers=%d): the cold posture does not hold the (b) shape", bCold.triggers)
	}
	if !bCold.engineBuilt || bCold.engineWarm {
		t.Fatalf("the cold face's echo-only window did not rebuild: built=%v warm=%v", bCold.engineBuilt, bCold.engineWarm)
	}
	for _, ep := range rig.rv.bargein.readings().Episodes {
		if ep.EnergyPass {
			t.Fatalf("the cold face's echo-only window let an echo episode through the tightened cut (%.2f dBFS)", ep.RmsDbfs)
		}
	}
	t.Logf("COLD-FACE echo-only triggers=0 episodes=%d all_below_cut=true built=true (the (b) shape holds on the rebuilt engines)", bCold.episodes)

	// The restore core: the switch back on - the default's first window
	// after the cold rounds builds fresh (the cold face closed the engines)
	// and the (b) shape holds again.
	rig.rv.bargein.setKeepalive(true)
	bRestore := voice392Capture(t, rig, voice26pAnnounceA, echoWav, "", 0, 6)
	if bRestore.triggers != 0 {
		t.Fatalf("the restore round's echo self-triggered (triggers=%d)", bRestore.triggers)
	}
	if !bRestore.engineBuilt || bRestore.engineWarm {
		t.Fatalf("the restore round did not rebuild (the cold face closed the engines; the default's first window back builds fresh): built=%v warm=%v",
			bRestore.engineBuilt, bRestore.engineWarm)
	}
	for _, ep := range rig.rv.bargein.readings().Episodes {
		if ep.EnergyPass {
			t.Fatalf("the restore round let an echo episode through the tightened cut (%.2f dBFS)", ep.RmsDbfs)
		}
	}
	t.Logf("DEFAULT-RESTORE triggers=0 episodes=%d all_below_cut=true built=true (the default is back and the pair closed)", bRestore.episodes)
	write392JSON(t, "392-keepalive.json", map[string]any{
		"control_first_build": ctrl,
		"warm_rounds":         voice392RoundsOut(warmRounds),
		"warm_echo":           voice392RoundOutFrom(bWarm),
		"cold_round":          voice392RoundOutFrom(aCold),
		"cold_echo":           voice392RoundOutFrom(bCold),
		"restore_echo":        voice392RoundOutFrom(bRestore),
		"talker_at":           voice392EarlyTalkerAt.String(),
		"talker_level":        voice26pTalkerDbfs,
		"default":             "the 401 flip made the cleared warm posture the production default: the held engines judge again window after window under the between-window clears (the WARM REVIVAL MOVED guard is live on the default path, no switch), the echo-only window holds on held engines with no false-start leak (the 398-v3 decay curve), and the switch-off face still carries the exact pre-392 cold shapes with the restore core bringing the default back",
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
		SwingDB float64            `json:"swing_db"`
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
		sr.SwingDB = sr.MaxRms - sr.MinRms
		t.Logf("TUNING-%s swing=%.2fdB (min=%.2f max=%.2f)", s.name, sr.SwingDB, sr.MinRms, sr.MaxRms)
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
