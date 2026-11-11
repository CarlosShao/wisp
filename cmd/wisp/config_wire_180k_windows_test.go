//go:build windows

package main

// The 180-census follow-up's wire cases (ticket 423's second wave): the
// census's zero-reader keys whose production readers this leg landed, one
// case per key family, each stating its feed and its reverse. The legs are
// organized like config_wire_423_windows_test.go: deleting the exact feed
// line reddens exactly the case above.
//
// Staging: no sherpa DLLs and no model install are needed except where the
// voice harness is reused (the VAD hatch env var, the same contract the 423
// cases run under).

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// ---------------------------------------------------------------- [observe] /
// [privacy] -> the log pipeline (installLogSink's boot feed)

// readLog180k concatenates every JSONL file the pipeline wrote in dir.
func readLog180k(t *testing.T, dir string) string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read log dir: %v", err)
	}
	var b strings.Builder
	for _, e := range entries {
		if e.IsDir() || !strings.HasPrefix(e.Name(), "wisp-") || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatalf("read %s: %v", e.Name(), err)
		}
		b.Write(raw)
	}
	return b.String()
}

func logSinkConfig180k(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(body), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return dir
}

func Test180KObserveLevelGatesTheFile(t *testing.T) {
	dir := logSinkConfig180k(t, "[observe]\nlevel = \"warn\"\n")
	sink, err := installLogSink(dir)
	if err != nil {
		t.Fatalf("installLogSink: %v", err)
	}
	sink.logger().Info("180k info marker below the gate")
	sink.logger().Warn("180k warn marker above the gate")
	sink.close()
	file := readLog180k(t, filepath.Join(dir, "logs"))
	if strings.Contains(file, "180k info marker") {
		t.Fatal("observe.level=warn still wrote Info records to the file: the level key is not riding the handler's filter")
	}
	if !strings.Contains(file, "180k warn marker above the gate") {
		t.Fatal("the Warn record never landed: the pipeline did not run at all")
	}
}

func Test180KRedactPathsMasksTheFile(t *testing.T) {
	// The positive control beside it (redact off) is what separates "the
	// redactor ate it" from "the record never carried the path".
	on := logSinkConfig180k(t, "[privacy]\nredact_paths = true\n")
	sinkOn, err := installLogSink(on)
	if err != nil {
		t.Fatalf("installLogSink(on): %v", err)
	}
	sinkOn.logger().Warn("see C:\\Users\\alice\\notes.txt for details")
	sinkOn.close()
	masked := readLog180k(t, filepath.Join(on, "logs"))
	// The JSONL layer escapes backslashes, so the honest predicate is the
	// filename: masked means no path fragment survives, not that the exact
	// single-backslash spelling is absent from raw JSON bytes.
	if strings.Contains(masked, "alice") || strings.Contains(masked, "notes.txt") {
		t.Fatal("privacy.redact_paths=true left the raw path in the file: the key is not riding the redactor")
	}
	if !strings.Contains(masked, "<path>") {
		t.Fatalf("the masked record lost the mask shape itself: %q", masked)
	}

	off := logSinkConfig180k(t, "[privacy]\nredact_paths = false\n")
	sinkOff, err := installLogSink(off)
	if err != nil {
		t.Fatalf("installLogSink(off): %v", err)
	}
	sinkOff.logger().Warn("see C:\\Users\\alice\\notes.txt for details")
	sinkOff.close()
	raw := readLog180k(t, filepath.Join(off, "logs"))
	if !strings.Contains(raw, "notes.txt") || !strings.Contains(raw, "alice") {
		t.Fatalf("the redact=false control does not carry the raw path any more: the control is dead and proves nothing (%q)", raw)
	}
}

func Test180KRollSizeGatesTheFile(t *testing.T) {
	dir := logSinkConfig180k(t, "[observe]\n\n[observe.roll]\nsize_mb = 1\ndays = 7\n")
	sink, err := installLogSink(dir)
	if err != nil {
		t.Fatalf("installLogSink: %v", err)
	}
	// MaxLoggedString caps one record at 512 body bytes, so the size roll is
	// reached by volume: ~3000 small records clear 1 MB and roll the writer.
	for i := 0; i < 3000; i++ {
		sink.logger().Warn("180k roll volume", "i", i, "pad", strings.Repeat("x", 300))
	}
	sink.close()
	n, err := observe.CountLogFiles(filepath.Join(dir, "logs"))
	if err != nil {
		t.Fatalf("count log files: %v", err)
	}
	if n < 2 {
		t.Fatalf("roll.size_mb=1 produced %d file(s), want at least 2: the writer never rolled, so the size key is not riding the threshold", n)
	}
}

func Test180KRollDaysSweepsTheFile(t *testing.T) {
	dir := logSinkConfig180k(t, "[observe]\n\n[observe.roll]\nsize_mb = 10\ndays = 2\n")
	logs := filepath.Join(dir, "logs")
	// A file named five days ago: inside the compiled default's 7-day window,
	// outside the file's 2-day window. The boot sweep (newRollingWriterClock's
	// startup sweepLocked) deletes it before the first record lands, which is
	// the days key acting on the real sweep - and the 5-vs-7 gap is what makes
	// a gutted feed (the hardcoded default back) red instead of silent.
	oldName := "wisp-" + time.Now().UTC().Add(-5*24*time.Hour).Format("20060102") + "-001.jsonl"
	old := filepath.Join(logs, oldName)
	if err := os.MkdirAll(logs, 0o700); err != nil {
		t.Fatalf("mkdir logs: %v", err)
	}
	if err := os.WriteFile(old, []byte{}, 0o600); err != nil {
		t.Fatalf("plant old file: %v", err)
	}
	sink, err := installLogSink(dir)
	if err != nil {
		t.Fatalf("installLogSink: %v", err)
	}
	sink.close()
	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Fatalf("roll.days=2 kept %s: the retention sweep is not riding the configured window", oldName)
	}
}

// ---------------------------------------------------------------- [memory] /
// [privacy] -> the retention window and the profile cap

func Test180KRetentionWindowHelperPicksTheSmallerNonDefault(t *testing.T) {
	def := config.NewDefaults()
	if def.Memory.L3RetentionDays != 30 || def.Privacy.RetentionDays != 30 {
		t.Fatalf("this case's premise moved: the shipped defaults are %d / %d, want 30 / 30",
			def.Memory.L3RetentionDays, def.Privacy.RetentionDays)
	}
	cases := []struct {
		l3, privacy, want int
	}{
		{30, 30, 30}, // both default
		{10, 30, 10}, // l3 non-default
		{30, 7, 7},   // privacy non-default
		{10, 7, 7},   // both: the smaller (more protective) window wins
		{0, 0, 30},   // zero = not written, keeps the D35 default
		{-5, 90, 30}, // a negative cannot widen or zero the window
	}
	c := config.NewDefaults()
	for _, tc := range cases {
		c.Memory.L3RetentionDays = tc.l3
		c.Privacy.RetentionDays = tc.privacy
		if got := retentionWindowDays180K(c); got != tc.want {
			t.Errorf("l3=%d privacy=%d -> %d, want %d", tc.l3, tc.privacy, got, tc.want)
		}
	}
}

func Test180KMemoryKnobsReachTheAssemblyRoot(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	var live *agentRuntime
	h.rtHook = func(rt *agentRuntime) { live = rt }
	cfgPath := filepath.Join(h.dir, configFileName)
	body, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("read config.toml: %v", err)
	}
	// Non-default, chosen against the compiled defaults (cap 20, window 30) so
	// every slot can only carry the file's numbers.
	hand := "\n[memory]\nl1_max = 7\n\n[privacy]\nretention_days = 7\n"
	if err := os.WriteFile(cfgPath, []byte(string(body)+hand), 0o600); err != nil {
		t.Fatalf("plant the [memory]/[privacy] tail: %v", err)
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("wisp run exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	if live == nil || live.store == nil {
		t.Fatal("the runtime hook never ran, so nothing here read the assembled host")
	}
	if got := live.store.ProfileCap(); got != 7 {
		t.Errorf("store.ProfileCap() = %d, want the file's 7 (memory.l1_max feeding memory.WithProfileCap; the compiled default is 20)", got)
	}
	if live.retention == nil {
		t.Error("assembleRuntime started no retention job: memory.StartRetentionJob's production owner is missing")
	}
	want := retentionBootLine180K(7)
	if !strings.Contains(h.err.String(), "task_log 窗口 7 天") {
		t.Errorf("the boot audit does not carry the retention window: want substring of %q, stderr=%q", want, h.err.String())
	}
}

func Test180KL1EnabledGateGatesTheAttach(t *testing.T) {
	// The disabled shape: the key's answer is that the pipeline is NOT
	// attached, and the attach path (which would noteProblem on this bare
	// runtime) must not run at all.
	off := logSinkConfig180k(t, "[memory]\nl1_enabled = false\n")
	rv := &residentVoice{residency: newVoiceResidency(voicePlan{}, newVoiceOwnerQueue(), nil)}
	attachMemoryL1IfEnabled180K(rv, &agentRuntime{}, off)
	if probs := rv.problemsSnapshot(); len(probs) != 0 {
		t.Fatalf("l1_enabled=false still reached the attach path (problems %v): the gate does not read the key", probs)
	}
	if warm := rv.residency.l1Warm; warm != nil {
		t.Fatal("the residency carries a warm entry with l1_enabled=false: the gate wired the pipeline anyway")
	}

	// The enabled control: the same bare runtime reaches the attach and its
	// named refusal ("no store") is booked - proof the gate let it through.
	on := logSinkConfig180k(t, "[memory]\nl1_enabled = true\n")
	rv2 := &residentVoice{residency: newVoiceResidency(voicePlan{}, newVoiceOwnerQueue(), nil)}
	attachMemoryL1IfEnabled180K(rv2, &agentRuntime{}, on)
	found := false
	for _, p := range rv2.problemsSnapshot() {
		if strings.Contains(p, "记忆抽取未装配") {
			found = true
		}
	}
	if !found {
		t.Fatalf("l1_enabled=true never reached the attach path (problems %v): the enabled control is dead", rv2.problemsSnapshot())
	}
}

// ------------------------------------------------- [voice] conversation_mode

func Test180KConversationModeReachesTheSpeakDoneRow(t *testing.T) {
	dataDir := voice15eDataDir(t)
	t.Setenv("WISP_SILERO_MODEL", voice15eVadModel(t))
	voice15eWriteConfig(t, dataDir, func(c *config.Config) {
		c.Voice.ConversationMode = true
	})
	ra := &residentAudio{frames: make(chan []byte)}
	rv := buildResidentVoice(dataDir, ra, nil, nil, nil)
	if rv.verdict != "语音链已装配" {
		t.Fatalf("the voice leg did not assemble: verdict=%q (problems=%v)", rv.verdict, rv.problemsSnapshot())
	}
	if !rv.plan.ConversationMode {
		t.Fatal("the boot plan carries ConversationMode=false off a true file: the feed does not ride the boot LoadFile")
	}
	// The machine half: the exact facts the announcement end dispatches walk
	// the D43 row the config selects - #28 into Listening (mic stays open,
	// PLAN.md:1326's "Speaking 结束直接回 Listening") instead of #26 into Warm
	// (mic closed). A machine born in Speaking needs no announce machinery to
	// judge the row.
	facts := rv.speakDoneFacts()
	if !facts.ConversationMode {
		t.Fatal("speakDoneFacts lost the plan value: the dispatch would walk #26 over a conversation_mode=true file")
	}
	m := statemachine.New(statemachine.Options{Initial: statemachine.StateSpeaking})
	to, err := m.Dispatch(statemachine.EvSpeakDone, facts)
	if err != nil {
		t.Fatalf("EvSpeakDone refused the conversation facts: %v", err)
	}
	if to != statemachine.StateListening {
		t.Fatalf("conversation_mode=true landed in %s, want Listening (D43 #28)", to)
	}

	// The shipped-default control: facts false, row #26, Warm.
	off := &residentVoice{plan: voicePlan{ConversationMode: false}}
	to2, err := statemachine.New(statemachine.Options{Initial: statemachine.StateSpeaking}).
		Dispatch(statemachine.EvSpeakDone, off.speakDoneFacts())
	if err != nil {
		t.Fatalf("EvSpeakDone refused the default facts: %v", err)
	}
	if to2 != statemachine.StateWarm {
		t.Fatalf("conversation_mode=false landed in %s, want Warm (D43 #26 - the pre-wire shape)", to2)
	}
}

// ------------------------------------------------------- [voice] aec.enabled

func Test180KAECSwitchGatesTheMonitor(t *testing.T) {
	dataDir := voice15eDataDir(t)
	t.Setenv("WISP_SILERO_MODEL", voice15eVadModel(t))
	voice15eWriteConfig(t, dataDir, func(c *config.Config) {
		c.Voice.AEC.Enabled = false
	})
	ra := &residentAudio{frames: make(chan []byte)}
	rv := buildResidentVoice(dataDir, ra, nil, nil, nil)
	if rv.verdict != "语音链已装配" {
		t.Fatalf("the voice leg did not assemble: verdict=%q (problems=%v)", rv.verdict, rv.problemsSnapshot())
	}
	if rv.plan.AECEnabled {
		t.Fatal("the boot plan carries AECEnabled=true off an aec.enabled=false file: the feed does not ride the boot LoadFile")
	}
	if rv.bargein == nil {
		t.Fatal("the monitor was not built at all: the feed moved the construction, not the switch")
	}
	if rv.bargein.arm() {
		t.Fatal("the monitor armed with aec.enabled=false: the switch does not gate the arm")
	}
	if !strings.Contains(rv.bargein.reason, "aec.enabled=false") {
		t.Fatalf("the arm refusal does not name the config switch: %q", rv.bargein.reason)
	}

	// The default control: the shipped aec.enabled=true does not produce the
	// config refusal - the monitor's next reason is its own staging posture
	// (no C11 scope in this bare construction), never the config one.
	def := voice15eDataDir(t)
	ra2 := &residentAudio{frames: make(chan []byte)}
	rv2 := buildResidentVoice(def, ra2, nil, nil, nil)
	if !rv2.plan.AECEnabled {
		t.Fatal("the compiled default carries AECEnabled=false: the schema default is true and the feed must carry it")
	}
	rv2.bargein.arm()
	if strings.Contains(rv2.bargein.reason, "aec.enabled=false") {
		t.Fatalf("the default posture answered the config refusal: %q", rv2.bargein.reason)
	}
}

// -------------------------------------------------------- [audio] input_device

func Test180KInputDeviceReachesTheCaptureConfig(t *testing.T) {
	rt, _ := bootAudioRuntime(t)
	dir := writeAudioConfig(t, func(c *config.Config) { c.Audio.InputDevice = "Blue Yeti" })
	var got audio.CaptureConfig
	fake := &failingSource{}
	ra := assembleCapture(rt, dir, (&levelTap{}).call, func(opts ...audio.CaptureOption) captureSource {
		c := &audio.CaptureConfig{}
		for _, o := range opts {
			o(c)
		}
		got = *c
		return fake
	})
	if ra == nil {
		t.Fatal("assembleCapture returned nothing")
	}
	if got.DeviceName != "Blue Yeti" {
		t.Fatalf("the capture config carries DeviceName %q, want the file's \"Blue Yeti\": audio.input_device is not riding the boot feed", got.DeviceName)
	}
}
