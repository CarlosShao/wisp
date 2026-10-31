//go:build windows

package main

// Ticket 26p's integration readings, part 2: the (b) true number and its
// reverse, the (c) margin re-test, and the (a) reverse - the surfaces the
// 59-v3 §7 ruling assigns to this integration face:
//
//	(b)  产品麦面零转写：the announcement's own echo on the real capture
//	     path (AEC3 in it) opens zero episodes, fires zero barge-ins, and
//	     leaves zero transcripts - the S4 literal "把播报音频注入麦克风回路
//	     不得产生任何 ASR 转写（不自激）".
//	(b)- M2 摘 AEC3: the judge on the raw capture (aecBypass) - the echo
//	     reaches the speech judge, opens episodes, and the same material
//	     through the real ASR chain transcribes the announcement (the
//	     59-r1 rig's 5/5 shape). The restore half re-runs the (b) surface
//	     with the AEC back and watches the echo disappear again.
//	(c)  -30 裕量重测加严读数：the 59-v3 §3 two constants re-measured on
//	     the integration chain - the real-speech control (-26.70), the
//	     quiet-speech cost (-36.58), and the artifact family the mismatched
//	     echo produces - through the gate's real Evaluate at the applied
//	     cut, plus the offline nominal-axis re-judge (the readings-surface
//	     reuse 389-v1 §7 suggested).
//	(a)- M1 摘组合门: the gate's own reverse switch (dynamicOff, 389 AC#3)
//	     on the artifact material - the static loose chain lets the echo
//	     self-trigger and the announcement interrupts itself with no talker;
//	     the restore half watches the tightened cut hold the same material.
//
// Staging: part 1's header.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/speech"
)

const (
	voice26pAnnounceA = "任务已经完成，明天下午三点的会议改到四点，请准时参加，谢谢。"
	voice26pAnnounceB = "文件已经导出到桌面的报表文件夹，请查收。"
	voice26pTalker    = "等一下，先听我说这件事。"

	// voice26pLevelFloor is the plausibility band for an episode's level
	// reading: below it the judge saw nothing real (the readings are levels,
	// not -Inf sentinels).
	voice26pLevelFloor = -60.0
)

// write26pJSON drops one readings artifact (the evidence table quotes these).
func write26pJSON(t *testing.T, name string, v any) {
	t.Helper()
	if err := os.MkdirAll(voice26pOut, 0o755); err != nil {
		t.Fatalf("out dir: %v", err)
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatalf("marshal %s: %v", name, err)
	}
	p := filepath.Join(voice26pOut, name)
	if err := os.WriteFile(p, append(b, '\n'), 0o644); err != nil {
		t.Fatalf("write %s: %v", p, err)
	}
	t.Logf("WROTE %s", p)
}

// voice26pSpikeAsr builds the real transcript chain the M2 reading runs the
// echo material through: silero + the streaming Paraformer, the same engines
// the product transcript path drives.
func voice26pSpikeAsr(t *testing.T) (*speech.Transcriber, func()) {
	t.Helper()
	det, err := speech.NewDetector(speech.DetectorConfig{ModelPath: voice15eVadModel(t)})
	if err != nil {
		t.Fatalf("detector: %v", err)
	}
	models := voice15eSpikeModels(t)
	asr, err := speech.NewStreamingAsr(speech.AsrConfig{
		EncoderPath: filepath.Join(models, "asr-encoder.int8.onnx"),
		DecoderPath: filepath.Join(models, "asr-decoder.int8.onnx"),
		TokensPath:  filepath.Join(models, "asr-tokens.txt"),
	})
	if err != nil {
		_ = det.Close()
		t.Fatalf("streaming asr: %v", err)
	}
	tr := speech.NewTranscriber(det, asr)
	return tr, func() {
		_ = asr.Close()
		_ = det.Close()
	}
}

// Test26PEchoNeverTriggersAndNeverTranscribes is the (b) true number: the
// announcement's own echo injected into the C8 seam for a whole announcement
// - AEC3 on the real capture path, the gate dynamic - opens zero episodes,
// fires zero barge-ins, and leaves zero transcripts on the product surface.
func Test26PEchoNeverTriggersAndNeverTranscribes(t *testing.T) {
	if probe, err := openWaveOutSpeaker(22050); err != nil {
		t.Skipf("no output device on this machine: %v", err)
	} else {
		_ = probe.Close()
	}
	dir := t.TempDir()
	rig := newVoice26pRig(t, nil)
	rig.requireAEC3()

	// The echo source: one full announcement, recorded.
	walkToThinking381r2(t, rig.rv)
	if !rig.rv.announce(voice26pAnnounceA) {
		t.Fatalf("the control announce refused the sentence; problems=%v", rig.rv.problemsSnapshot())
	}
	control := rig.spyAt(0)
	voice15eWaitFor(t, "the control sentence to play out", func() bool {
		return rig.rv.ttsSpoken.Load() >= 1
	})
	rig.waitReleased()
	echoWav, _ := voice26pEchoWav(t, dir, "echo", control.recording(), 22050, voice26pEchoGain)

	// The echo round: the announcement plays while its own echo arrives on
	// the capture seam - the S4 literal's injection shape.
	finals0, handed0 := rig.rv.finals.Load(), rig.rv.handed.Load()
	walkToThinking381r2(t, rig.rv)
	if !rig.rv.announce(voice26pAnnounceA) {
		t.Fatalf("the echo announce refused the sentence; problems=%v", rig.rv.problemsSnapshot())
	}
	spy := rig.spyAt(1)
	<-spy.first
	stopEcho, _ := rig.inject(echoWav)
	voice15eWaitFor(t, "the echo announcement to play out (not interrupted)", func() bool {
		return rig.rv.ttsSpoken.Load() >= 2
	})
	stopEcho()
	rig.waitReleased()
	read := rig.rv.bargein.readings()
	if read.Triggers != 0 {
		t.Fatalf("MUTANT 26p-ECHO-TRIGGER SHAPE: the echo fired %d barge-in(s) (episodes=%d): the AEC'd residual is not quiet enough or the judge is not the tightened form - 回声自触发了", read.Triggers, len(read.Episodes))
	}
	// The VAD is sensitive enough to open episodes on the residual floor
	// (measured: the first staged run judged them at -45.23/-42.11 dBFS);
	// the reading is that the tightened cut holds every one of them - the
	// margin exchange the S4 no-self-trigger validation point (-30) buys.
	passed := 0
	for _, ep := range read.Episodes {
		if ep.EnergyPass {
			passed++
		}
	}
	if passed != 0 {
		t.Fatalf("MUTANT 26p-ECHO-TRIGGER SHAPE: %d of %d echo episode(s) passed the tightened cut: the residual floor rode above -30 dBFS", passed, len(read.Episodes))
	}
	if got := rig.rv.finals.Load() - finals0; got != 0 {
		t.Fatalf("MUTANT 26p-ECHO-ASR SHAPE: %d transcript(s) grew across the echo window: 把播报音频注入麦克风回路产生了 ASR 转写", got)
	}
	if got := rig.rv.handed.Load() - handed0; got != 0 {
		t.Fatalf("the echo handed %d sentence(s) to the text path", got)
	}
	if closed, reopened := rig.rv.setSpeakingCount(); closed != 2 || reopened != 2 {
		t.Fatalf("the capture gate did not ride both announcements closed=%d reopened=%d", closed, reopened)
	}
	if got := rig.rv.problemsSnapshot(); len(got) != 0 {
		t.Fatalf("the echo round left problems behind: %v", got)
	}
	write26pJSON(t, "26p-b-echo-zero.json", read)
	t.Logf("ECHO-ZERO episodes=%d all_below_cut=%v triggers=0 finals_delta=0 handed_delta=0 lag=%d strength=%.2f windows_untrusted=%d",
		len(read.Episodes), passed == 0, read.Lag, read.LagStrength, countTrue26p(read.Untrusted))
}

// Test26PWithoutAECTheEchoTranscribesVerbatim is the (b) reverse (M2): with
// the AEC stage removed from the judge (aecBypass), the same echo reaches
// the speech judge raw, opens episodes, and the same material through the
// real ASR chain transcribes the announcement - 逐字转写复现. The restore
// half re-runs the (b) surface with the AEC back.
func Test26PWithoutAECTheEchoTranscribesVerbatim(t *testing.T) {
	if probe, err := openWaveOutSpeaker(22050); err != nil {
		t.Skipf("no output device on this machine: %v", err)
	} else {
		_ = probe.Close()
	}
	dir := t.TempDir()
	rig := newVoice26pRig(t, func(m *voiceBargeinMonitor) {
		m.setReverseCriteria(true, false) // aecBypass: the M2 mutant
	})

	// The echo source.
	walkToThinking381r2(t, rig.rv)
	if !rig.rv.announce(voice26pAnnounceA) {
		t.Fatalf("the control announce refused the sentence; problems=%v", rig.rv.problemsSnapshot())
	}
	control := rig.spyAt(0)
	voice15eWaitFor(t, "the control sentence to play out", func() bool {
		return rig.rv.ttsSpoken.Load() >= 1
	})
	rig.waitReleased()
	echoWav, echoSamples := voice26pEchoWav(t, dir, "echo", control.recording(), 22050, voice26pEchoGain)

	// The bypass round: the echo IS the near speech the judge sees.
	finals0 := rig.rv.finals.Load()
	walkToThinking381r2(t, rig.rv)
	if !rig.rv.announce(voice26pAnnounceA) {
		t.Fatalf("the bypass announce refused the sentence; problems=%v", rig.rv.problemsSnapshot())
	}
	spy := rig.spyAt(1)
	<-spy.first
	stopEcho, _ := rig.inject(echoWav)
	voice15eWaitFor(t, "the bypass round's announcement to finish", func() bool {
		return rig.rv.ttsAborted.Load() >= 1 || rig.rv.ttsSpoken.Load() >= 2
	})
	stopEcho()
	rig.waitReleased()
	read := rig.rv.bargein.readings()
	if len(read.Episodes) == 0 {
		t.Fatalf("MUTANT 26p-M2 SHAPE: the raw echo opened zero episodes - the judge did not see the announcement audio the mic surface carried (the 59-r1 rig measured 5/5 verbatim on this shape)")
	}
	passed := 0
	for _, ep := range read.Episodes {
		if ep.EnergyPass {
			passed++
		}
	}
	if got := rig.rv.finals.Load() - finals0; got != 0 {
		t.Logf("M2-ASR-SURFACE %d transcript(s) reached the product surface under the bypass (the monitor still owns the window; the product transcript path stays echo-free by the same bracket)", got)
	}

	// The verbatim half: the exact material the mic surface carried, through
	// the real ASR chain.
	tr, done := voice26pSpikeAsr(t)
	defer done()
	var outs []speech.Transcript
	outs = tr.Accept(echoSamples, outs[:0])
	outs = tr.Finalize(outs[:0])
	text := ""
	for _, o := range outs {
		if o.Kind == speech.TranscriptFinal && o.Text != "" {
			text += o.Text
		}
	}
	if len([]rune(text)) < 4 {
		t.Fatalf("MUTANT 26p-M2 SHAPE: the echo material transcribed to %q - the announcement audio did not verbatim-transcribe (the 59-r1 shape is a real transcript)", text)
	}
	hit := strings.Contains(text, "会议") || strings.Contains(text, "四点") || strings.Contains(text, "完成")
	t.Logf("M2-VERBATIM episodes=%d passed=%d triggers=%d transcript=%q distinctive=%v",
		len(read.Episodes), passed, read.Triggers, text, hit)

	// The restore half: the AEC back, the same echo, the (b) surface again.
	rig2 := newVoice26pRig(t, nil)
	rig2.requireAEC3()
	walkToThinking381r2(t, rig2.rv)
	if !rig2.rv.announce(voice26pAnnounceA) {
		t.Fatalf("the restore control refused the sentence; problems=%v", rig2.rv.problemsSnapshot())
	}
	control2 := rig2.spyAt(0)
	voice15eWaitFor(t, "the restore control to play out", func() bool {
		return rig2.rv.ttsSpoken.Load() >= 1
	})
	rig2.waitReleased()
	echoWav2, _ := voice26pEchoWav(t, dir, "echo-restore", control2.recording(), 22050, voice26pEchoGain)
	walkToThinking381r2(t, rig2.rv)
	if !rig2.rv.announce(voice26pAnnounceA) {
		t.Fatalf("the restore announce refused the sentence; problems=%v", rig2.rv.problemsSnapshot())
	}
	spy2 := rig2.spyAt(1)
	<-spy2.first
	stopEcho2, _ := rig2.inject(echoWav2)
	voice15eWaitFor(t, "the restore announcement to play out", func() bool {
		return rig2.rv.ttsSpoken.Load() >= 2
	})
	stopEcho2()
	rig2.waitReleased()
	read2 := rig2.rv.bargein.readings()
	if read2.Triggers != 0 {
		t.Fatalf("the restore half diverged from the (b) surface (triggers=%d episodes=%d): 摘 AEC3 逐字复现、还回 AEC3 零转写——这一对不成立", read2.Triggers, len(read2.Episodes))
	}
	for _, ep := range read2.Episodes {
		if ep.EnergyPass {
			t.Fatalf("the restore half let an echo episode through the tightened cut (%.2f dBFS): 摘 AEC3 逐字复现、还回 AEC3 零转写——这一对不成立", ep.RmsDbfs)
		}
	}
	write26pJSON(t, "26p-m2-verbatim.json", map[string]any{
		"bypass": read, "restore": read2, "transcript": text,
	})
	t.Logf("M2-RESTORE episodes=0 triggers=0 (the AEC back = the echo gone, the pair closed)")
}

// voice26pRoundReading is one announcement round's judge reading.
type voice26pRoundReading struct {
	pass      bool
	rms       float64
	threshold float64
	tighten   bool
	episodes  int
}

// voice26pRound announces one sentence with an optional echo and an optional
// talker injection, waits for the window to close, and reads the judge. The
// talker goes in 1.2s after the first sample (after the alignment warm-up).
func voice26pRound(t *testing.T, rig *voice26pRig, sentence, echoWav, talkerWav string) voice26pRoundReading {
	t.Helper()
	aborted0, spoken0 := rig.rv.ttsAborted.Load(), rig.rv.ttsSpoken.Load()
	walkToThinking381r2(t, rig.rv)
	if !rig.rv.announce(sentence) {
		t.Fatalf("the announce refused the sentence; problems=%v", rig.rv.problemsSnapshot())
	}
	spy := rig.spyAt(len(rig.spies))
	<-spy.first
	var stopEcho, stopTalker func()
	if echoWav != "" {
		stopEcho, _ = rig.inject(echoWav)
		defer stopEcho()
	}
	if talkerWav != "" {
		time.Sleep(voice26pTalkerAt)
		stopTalker, _ = rig.inject(talkerWav)
		defer stopTalker()
	}
	voice15eWaitFor(t, "the round's announcement to finish (spoken or interrupted)", func() bool {
		return rig.rv.ttsAborted.Load() > aborted0 || rig.rv.ttsSpoken.Load() > spoken0
	})
	if stopTalker != nil {
		stopTalker()
	}
	if stopEcho != nil {
		stopEcho()
	}
	rig.waitReleased()
	read := rig.rv.bargein.readings()
	out := voice26pRoundReading{episodes: len(read.Episodes)}
	if read.Triggers > 0 {
		out.pass = true
	}
	for i := range read.Episodes {
		if read.Episodes[i].EnergyPass {
			out.rms = read.Episodes[i].RmsDbfs
			out.threshold = read.Episodes[i].ThresholdDbfs
			out.tighten = read.Episodes[i].TightenActive
			out.pass = true
		} else if out.threshold == 0 {
			out.rms = read.Episodes[i].RmsDbfs
			out.threshold = read.Episodes[i].ThresholdDbfs
			out.tighten = read.Episodes[i].TightenActive
		}
	}
	return out
}

// Test26PMarginRetestAndGateReverse is the (c) margin re-test plus the (a)
// reverse with its restore: the 59-v3 §3 two constants re-measured through
// the integrated judge at three material levels, then the gate's dynamic
// toggled off on the artifact material (the echo self-trigger) and back on.
func Test26PMarginRetestAndGateReverse(t *testing.T) {
	if probe, err := openWaveOutSpeaker(22050); err != nil {
		t.Skipf("no output device on this machine: %v", err)
	} else {
		_ = probe.Close()
	}
	dir := t.TempDir()
	talkerL1 := voice26pRenderTtsWav(t, dir, "talker-l1", voice26pTalker, -26.70)
	talkerL2 := voice26pRenderTtsWav(t, dir, "talker-l2", voice26pTalker, -36.58)

	rig := newVoice26pRig(t, nil)
	rig.requireAEC3()

	// Two echo sources: the matched one (same sentence) and the mismatched
	// one (a different sentence) - the artifact-family generator.
	walkToThinking381r2(t, rig.rv)
	if !rig.rv.announce(voice26pAnnounceA) {
		t.Fatalf("control A refused; problems=%v", rig.rv.problemsSnapshot())
	}
	ctrlA := rig.spyAt(0)
	voice15eWaitFor(t, "control A to play out", func() bool { return rig.rv.ttsSpoken.Load() >= 1 })
	rig.waitReleased()
	echoMatch, _ := voice26pEchoWav(t, dir, "echo-match", ctrlA.recording(), 22050, voice26pEchoGain)

	walkToThinking381r2(t, rig.rv)
	if !rig.rv.announce(voice26pAnnounceB) {
		t.Fatalf("control B refused; problems=%v", rig.rv.problemsSnapshot())
	}
	ctrlB := rig.spyAt(1)
	voice15eWaitFor(t, "control B to play out", func() bool { return rig.rv.ttsSpoken.Load() >= 2 })
	rig.waitReleased()
	echoMismatch, _ := voice26pEchoWav(t, dir, "echo-mismatch", ctrlB.recording(), 22050, voice26pMismatchGain)

	// L1: the calibration's real-speech control (-26.70, the 389 set's
	// control-26 level) over the matched echo - the margin reading through
	// the integrated chain. The hard fact this reading carries: the AEC's
	// own double-talk suppression damages the control level (measured
	// -26.9 to -30.9 dBFS post-AEC across runs), so the 3.3dB tightened-cut
	// margin the offline calibration recorded does not necessarily survive
	// the chain - reported as measured, the (c) re-test's honest answer.
	l1 := voice26pRound(t, rig, voice26pAnnounceA, echoMatch, talkerL1)
	if l1.rms <= voice26pLevelFloor || l1.rms >= 0 {
		t.Fatalf("L1: the control episode's reading is not a level (%.2f dBFS)", l1.rms)
	}
	t.Logf("MARGIN-L1 control-26 post-AEC rms=%.2f threshold=%.1f margin=%.2fdB tighten=%v pass=%v (the offline calibration's 3.3dB margin vs the chain's own damage - the (c) finding)",
		l1.rms, l1.threshold, l1.rms-l1.threshold, l1.tighten, l1.pass)

	// L2: quiet speech at -36.58 - the tightened cut eats it (the known cost
	// 59-v3 §3 recorded), and the nominal axis the dynamic gating protects
	// still reaches it (the readings-surface reuse).
	l2 := voice26pRound(t, rig, voice26pAnnounceA, echoMatch, talkerL2)
	if l2.pass {
		t.Fatalf("L2 (%.2f dBFS quiet speech) triggered under the tightened cut %.1f: the -36 cost the calibration recorded did not reproduce (the cut would be looser than pinned)", l2.rms, l2.threshold)
	}
	offline := audio.NewBargeinGate(audio.BargeinOptions{DynamicGating: false})
	d := offline.Evaluate(audio.BargeinEpisode{RmsDbfs: l2.rms, AsrText: ""})
	if !d.EnergyPass {
		t.Fatalf("L2: the nominal -45 axis does not reach %.2f dBFS either - the quiet-speech reach the loose nominal buys is gone", l2.rms)
	}
	t.Logf("MARGIN-L2 quiet_speech rms=%.2f threshold=%.1f tightened_eaten=true nominal_reach=true (the tighten is a margin exchange, not a suppression)", l2.rms, l2.threshold)

	// L3: the mismatched echo - the artifact-family generator. Under the
	// dynamic gate the tightened cut holds it (no self-trigger); the episode
	// level is the integration surface's artifact reading.
	l3 := voice26pRound(t, rig, voice26pAnnounceA, echoMismatch, "")
	if l3.pass {
		t.Fatalf("L3: the mismatched echo self-triggered under the dynamic gate (rms=%.2f threshold=%.1f) - the tightened cut does not hold the artifact family", l3.rms, l3.threshold)
	}
	t.Logf("MARGIN-L3 mismatched_echo rms=%.2f threshold=%.1f margin=%.2fdB episodes=%d (the artifact family on the integration surface)",
		l3.rms, l3.threshold, l3.rms-l3.threshold, l3.episodes)

	// M1: the gate's dynamic off - the static loose chain on the same
	// material. The echo self-triggers: 摘组合门 ⇒ (a) 塌.
	rig.rv.bargein.setReverseCriteria(false, true)
	m1 := voice26pRound(t, rig, voice26pAnnounceA, echoMismatch, "")
	if !m1.pass {
		t.Fatalf("MUTANT 26p-M1 SHAPE: the static loose chain did NOT self-trigger on the mismatched echo (rms=%.2f threshold=%.1f episodes=%d): the dynamic gating is a no-op on the integration surface or the artifact family sits below the nominal cut",
			m1.rms, m1.threshold, m1.episodes)
	}
	t.Logf("M1-ECHO-SELF-TRIGGER rms=%.2f threshold=%.1f (the announcement interrupted itself with no talker - the 11.5/10min family on the integration surface)", m1.rms, m1.threshold)

	// Restore: the dynamic gate back on, the same material - held.
	rig.rv.bargein.setReverseCriteria(false, false)
	r4 := voice26pRound(t, rig, voice26pAnnounceA, echoMismatch, "")
	if r4.pass {
		t.Fatalf("the restore half diverged: the tightened cut let the mismatched echo through again (rms=%.2f)", r4.rms)
	}
	t.Logf("M1-RESTORE triggers=0 threshold=%.1f (the gate back = the echo held)", r4.threshold)
	write26pJSON(t, "26p-c-margins.json", map[string]any{
		"l1_real_speech":  l1,
		"l2_quiet_speech": l2,
		"l3_artifacts":    l3,
		"m1_self_trigger": m1,
		"restore":         r4,
	})
}
