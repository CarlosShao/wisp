//go:build windows

package main

// Ticket 247, the assembly-level half of AC#4 / AC#5 / AC#6.
//
// These run the production assembly function (assembleCapture -> the real
// WASAPI source constructor, the real gate, the real D38(e) registration)
// against a config.toml the test wrote. They are NOT AC#2: the reading that
// "the ball's number comes from sound" needs a live microphone and lives in
// resident_audio_247_live_windows_test.go behind WISP_LIVE_MIC=1.

import (
	"context"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/proc"
)

// levelTap counts the scalars the leg handed toward the ball and keeps them for
// the readings.
type levelTap struct {
	mu     sync.Mutex
	levels []float32
}

func (t *levelTap) call(level float32) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.levels = append(t.levels, level)
}

func (t *levelTap) count() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.levels)
}

// bootAudioRuntime brings up the test runtime with a private registry and
// returns a closer that runs the frozen sequence exactly once (the shutdown
// half of AC#5 needs the records, so it cannot be left to t.Cleanup alone).
func bootAudioRuntime(t *testing.T) (*proc.Runtime, func() []proc.StepRecord) {
	t.Helper()
	rt, err := proc.Boot(buildinfo.EnvTest, proc.WithRegistry(observe.NewRegistry()))
	if err != nil {
		t.Fatalf("proc.Boot(test): %v", err)
	}
	var once sync.Once
	shut := func() []proc.StepRecord {
		var recs []proc.StepRecord
		once.Do(func() { recs = rt.Shutdown(false) })
		return recs
	}
	t.Cleanup(func() { _ = shut() })
	return rt, shut
}

// writeAudioConfig lays a config.toml down in a temp dir from the shipped
// default table, optionally mutated, and returns the dir. The assembly reads it
// through config.LoadFile, the same path production takes.
func writeAudioConfig(t *testing.T, mutate func(*config.Config)) string {
	t.Helper()
	dir := t.TempDir()
	c := config.NewDefaults()
	if mutate != nil {
		mutate(c)
	}
	if err := config.SaveFile(filepath.Join(dir, configFileName), c); err != nil {
		t.Fatalf("write %s: %v", configFileName, err)
	}
	return dir
}

// TestAC247ShippedDefaultsAreTheOnesThisLegReads is AC#4's "not one notch"
// check taken at the moment this leg reads them: voice.enabled true,
// wake_word.enabled false, audio.mic_muted_default true. The wiring may not
// move a default, and it may not quietly decide what to do when the table it is
// reading is not the one it expected.
func TestAC247ShippedDefaultsAreTheOnesThisLegReads(t *testing.T) {
	c := config.NewDefaults()
	if !c.Voice.Enabled {
		t.Fatalf("[voice] enabled default = %v, want true (schema.go's tagged default)", c.Voice.Enabled)
	}
	if c.Voice.WakeWord.Enabled {
		t.Fatalf("[voice wake_word] enabled default = %v, want false", c.Voice.WakeWord.Enabled)
	}
	if !c.Audio.MicMutedDefault {
		t.Fatalf("[audio] mic_muted_default = %v, want true", c.Audio.MicMutedDefault)
	}
}

// TestAC247VoiceDisabledBuildsNoCollector is P1 form 甲's second half: a
// disabled voice path means the collector is not constructed at all, so there
// is no device handle, no capture thread and no step-4 owner to pretend about.
func TestAC247VoiceDisabledBuildsNoCollector(t *testing.T) {
	rt, _ := bootAudioRuntime(t)
	dir := writeAudioConfig(t, func(c *config.Config) { c.Voice.Enabled = false })

	tap := &levelTap{}
	ra := assembleCapture(rt, dir, tap.call, newRealCaptureSource)

	if ra.gate != nil || ra.mic != nil {
		t.Fatal("voice.enabled=false still built a capture stack")
	}
	if ra.started {
		t.Fatal("a leg with no collector reports itself running")
	}
	if tap.count() != 0 {
		t.Fatalf("the ball seam received %d levels from a leg that built nothing", tap.count())
	}
	if n := rt.Registry.CountByName("audio-capture"); n != 0 {
		t.Fatalf("a leg with no collector booked %d audio-capture threads", n)
	}
	for _, s := range rt.RegisteredShutdownSteps() {
		if s == proc.StepStopAudio {
			t.Fatal("step 4 has an owner in a process that owns no audio thread")
		}
	}
	if !strings.Contains(ra.posture(), "未构造") {
		t.Fatalf("boot posture hides the reason: %q", ra.posture())
	}
}

// TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice is AC#4's machine
// reading: at the shipped defaults the gate is mounted WITH
// mic_muted_default handed into it (gate.go:57's "map it here at boot wiring",
// the wiring that did not exist before this leg), and because the gate is
// muted the inner source is never started - so the synchronous device open
// inside wasapimic_windows.go never runs. Double-clicking this build captures
// nothing.
func TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice(t *testing.T) {
	rt, _ := bootAudioRuntime(t)
	dir := writeAudioConfig(t, nil) // the shipped table, untouched

	tap := &levelTap{}
	ra := assembleCapture(rt, dir, tap.call, newRealCaptureSource)

	if ra.gate == nil {
		t.Fatalf("no gate at the default (voice is enabled there): %q", ra.verdict)
	}
	if !ra.gate.Muted() {
		t.Fatal("WithStartMuted did not receive [audio] mic_muted_default: the mic would open at boot")
	}
	if ra.gate.Open() {
		t.Fatal("the gate handed the device to the capture thread at the default")
	}
	if ra.started {
		t.Fatal("posture claims a running capture leg with the gate muted")
	}
	if st := ra.mic.Stats(); st.FramesSent != 0 || st.BytesSent != 0 {
		t.Fatalf("the microphone delivered frames at the default: %+v", st)
	}
	if n := rt.Registry.CountByName("audio-capture"); n != 0 {
		t.Fatalf("a muted boot still spawned the capture thread (%d booked)", n)
	}
	if tap.count() != 0 {
		t.Fatalf("the ball seam got %d levels while muted", tap.count())
	}
	if !strings.Contains(ra.posture(), "设备未打开") {
		t.Fatalf("the boot line does not say the device is closed: %q", ra.posture())
	}
	if !strings.Contains(ra.posture(), "mic_muted_default") {
		t.Fatalf("the boot line does not name the key that decided it: %q", ra.posture())
	}
}

// TestAC247UnmuteReachesTheBallSeam is the capability half of AC#3 and AC#1:
// the one thing this leg hands the ball is a float32, and the way to make the
// device deliver it is an unmuted gate, which is what the operator's mute hot
// key does. The device itself is not opened here - the source is the failing
// one AC#6 needs - so this pair of assertions lives with the seam, not with the
// live reading.
func TestAC247UnmuteReachesTheBallSeam(t *testing.T) {
	rt, _ := bootAudioRuntime(t)
	dir := writeAudioConfig(t, nil)

	tap := &levelTap{}
	ra := assembleCapture(rt, dir, tap.call, newRealCaptureSource)
	if ra.gate == nil {
		t.Fatal("no collector at the default")
	}
	// levelOut is the sink the capture thread calls. Its signature is the whole
	// of what crosses: one scalar in, nothing back.
	sink := ra.levelOut(tap.call)
	sink(0.25)
	if tap.count() != 1 {
		t.Fatalf("the ball seam saw %d levels, want 1", tap.count())
	}
	if got := ra.levels.Load(); got != 1 {
		t.Fatalf("levels counter = %d, want 1 (the exit line reads it)", got)
	}
	// A nil destination (no ball in this process) must not lose the count or
	// panic: the resident leg can be assembled on a machine whose window failed.
	before := tap.count()
	silent := ra.levelOut(nil)
	silent(0.5)
	if tap.count() != before {
		t.Fatalf("a sink with no destination invented %d levels", tap.count()-before)
	}
	// And a nil handle is nil-safe, because runResident can print a posture for
	// a leg that was never assembled.
	var nilLeg *residentAudio
	nilLeg.levelOut(nil)(0.1)
	if nilLeg.posture() == "" {
		t.Fatal("a nil capture handle has no posture sentence to print")
	}
}

// TestAC247DeviceFailureShapesStillBootAndSayTheLoss is AC#6 at the assembly
// level, in two of the three forms the ticket lists (occupied, permission
// denied) plus the third it names (no device). The error objects are the
// package's own canonical D42 mapping, so the guidance text under test is the
// text in internal/audio/device.go:86/:89/:103, not a paraphrase.
//
// What each subtest asserts is the shape the ticket asked for: the leg returns
// (the boot goes on), the class and the guidance are said out loud, nothing is
// pushed into the state machine, and the teardown is still owned.
func TestAC247DeviceFailureShapesStillBootAndSayTheLoss(t *testing.T) {
	cases := []struct {
		name        string
		hr          uintptr
		device      string
		wantClass   string
		wantSaidSub string
	}{
		{"occupied", 0x8889000A, "Busy Mic", "audio_device", "exclusive mode"},
		{"permission denied", 0x80070005, "Built-in Mic", "audio_device", "Privacy & security"},
		{"no device", 0x88890004, "Gone Headset", "audio_device", "device invalidated"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rt, shut := bootAudioRuntime(t)
			// An unmuted default in this test is deliberate and it is not the
			// shipped table: it is the only way to make the assembly reach the
			// device open at all, which is the failure under examination.
			dir := writeAudioConfig(t, func(c *config.Config) { c.Audio.MicMutedDefault = false })

			dev := audio.DeviceDescriptor{ID: "id-" + tc.name, Name: tc.device}
			want := audio.DeviceError(tc.hr, dev, "capture open")
			fake := &failingSource{err: want, lastErr: want.Error()}

			tap := &levelTap{}
			ra := assembleCapture(rt, dir, tap.call, func(...audio.CaptureOption) captureSource { return fake })

			if ra.gate == nil {
				t.Fatal("a failed device stopped the assembly from returning")
			}
			if ra.started {
				t.Fatal("posture claims the capture leg is running over a failed device")
			}
			if !strings.Contains(ra.verdict, tc.wantClass) {
				t.Fatalf("verdict does not name the class %q: %q", tc.wantClass, ra.verdict)
			}
			if !strings.Contains(ra.verdict, tc.wantSaidSub) || !strings.Contains(ra.verdict, dev.Name) {
				t.Fatalf("verdict lost the guidance or the device name: %q", ra.verdict)
			}
			if !strings.Contains(ra.posture(), tc.wantSaidSub) {
				t.Fatalf("the boot line does not carry the loss: %q", ra.posture())
			}
			// The refusal to start stays refused: this leg never exits, never
			// calls the state machine, and pushes no event. The only thing the
			// process did with the failure was say it.
			if tap.count() != 0 {
				t.Fatalf("a failed device still delivered %d levels", tap.count())
			}
			owns := false
			for _, s := range rt.RegisteredShutdownSteps() {
				if s == proc.StepStopAudio {
					owns = true
				}
			}
			if !owns {
				t.Fatal("the gate this leg started has no step-4 teardown")
			}
			if recs := shut(); len(recs) != 10 {
				t.Fatalf("audit trail = %d records, want the frozen 10", len(recs))
			}
		})
	}
}

// TestAC247CaptureLegOwnsStepFourOfTheFrozenOrder is AC#5: the hook lands on
// the slot D38(e) already has, the ten steps keep their order and their names,
// and no eleventh step appeared.
func TestAC247CaptureLegOwnsStepFourOfTheFrozenOrder(t *testing.T) {
	rt, shut := bootAudioRuntime(t)
	dir := writeAudioConfig(t, nil)

	ra := assembleCapture(rt, dir, (&levelTap{}).call, newRealCaptureSource)
	if ra.gate == nil {
		t.Fatal("no collector to own step 4")
	}
	recs := shut()
	if len(recs) != 10 {
		t.Fatalf("audit trail = %d records, want 10 (no eleventh step)", len(recs))
	}
	wantNames := []string{
		"scheduler-close", "stop-hotkey-kws", "cancel-task-roots", "stop-audio",
		"release-speech-sessions", "destroy-panel-webview", "flush-logs-close-db",
		"free-os-memory", "close-job-object", "exit",
	}
	for i, want := range wantNames {
		if recs[i].Name != want {
			t.Fatalf("step %d name = %q, want %q (the frozen order)", i+1, recs[i].Name, want)
		}
	}
	step4 := recs[proc.StepStopAudio-1]
	if step4.Skipped {
		t.Fatalf("step 4 ran skipped while this leg owned it: %+v", step4)
	}
	if step4.Err != nil {
		t.Fatalf("step 4 reported an error on a muted leg: %+v", step4)
	}
	if ra.levels.Load() != 0 {
		t.Fatal("a muted leg delivered levels")
	}
	// The thread is gone from the roster: Stop joined it (or it never started).
	if n := rt.Registry.CountByName("audio-capture"); n != 0 {
		t.Fatalf("audio-capture still booked after the sequence: %d", n)
	}
}

// failingSource is the shape of a device that refuses: Start answers the
// canonical audio_device error, Stats carries its text and Err carries the
// object - the same three reads the real WASAPIMicrophone makes available after
// a failed open (wasapimic_windows.go: meter.fail + postErr).
type failingSource struct {
	err     error
	lastErr string
}

func (f *failingSource) Start(context.Context, chan<- []byte) error { return f.err }
func (f *failingSource) Stop() error                                { return nil }
func (f *failingSource) Stats() audio.Stats                         { return audio.Stats{LastError: f.lastErr} }
func (f *failingSource) Err() error                                 { return f.err }
