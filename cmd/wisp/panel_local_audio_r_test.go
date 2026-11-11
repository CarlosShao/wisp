package main

// The 本机音频 alignment at the assembly root (the 2026-11-11 S-leg follow-up,
// leg R). The backend already had real readers/guards for voice.aec.enabled
// (the barge-in monitor's switch), audio.input_device (the capture open's
// resolver) and audio.sample_rate (the load-time 16000 guard) while the panel
// had no channel for any of them; this leg adds the channel on both halves and
// these are its teeth:
//
//  1. the snapshot keys ride the LIVE reader - a real Manager, the real
//     configStore, and a hot-applied edit changing the answer without a
//     restart (the wire-byte pin lives panel-side, where the JSON is built:
//     internal/panel's Test313PayloadKeySetIsFrozenAndCarriesNoValueBearingKey);
//  2. the write fields route: ApplySetting lands voice.aec.enabled and
//     audio.input_device in the FILE through the guarded setter, and after the
//     guarded write's reload the live reader answers the new values;
//  3. the refusals: a non-bool word for the switch is refused before any byte
//     moves, and audio_sample_rate - the guard key - has no write leg at all.
//
// The 摘 reader => 无键 reverse is pinned panel-side
// (Test313SilentAndUnreadableLegsRenderAsStatesNotAsFacts); the 摘写 => 红
// reverse is the first Fatalf below (delete a case line in
// panel_config_store.go and the default branch refuses the field).

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
)

// localAudioRBase is a minimal valid file the decoder merges over the compiled
// defaults: AEC off (the non-default half, so "off" is a real answer a zero
// value could not fake), a named device, the chain rate.
const localAudioRBase = `[voice.aec]
enabled = false

[audio]
input_device = "Blue Yeti"
sample_rate = 16000
`

func localAudioRStore(t *testing.T) (*configStore, *config.Manager, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, configFileName)
	if err := os.WriteFile(path, []byte(localAudioRBase), 0o600); err != nil {
		t.Fatal(err)
	}
	mgr, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("manager: %v", err)
	}
	return newConfigStore(mgr, nil, nil, "test-local-audio"), mgr, path
}

// TestConfigLocalAudioSnapshotKeysRideTheLiveReader: the three reader-gated
// keys answer the FILE's values through the live Manager, and a hot-applied
// edit (hand edit + one reload poll) changes the answer without a restart -
// the same per-call reader class the [llm] row's consumed claim rests on,
// which is exactly what cmd/wisp/config_readers_255.go's re-adjudicated
// [audio] row now cites.
func TestConfigLocalAudioSnapshotKeysRideTheLiveReader(t *testing.T) {
	s, mgr, path := localAudioRStore(t)
	ctx := context.Background()

	view, err := s.ReadSettings(ctx)
	if err != nil {
		t.Fatalf("ReadSettings: %v", err)
	}
	if !view.Readable {
		t.Fatal("the leg answered unreadable over a file it just wrote")
	}
	if view.Voice.AecEnabled == nil || *view.Voice.AecEnabled {
		t.Fatalf("AecEnabled = %v, want the answered false (AEC off is a real answer, not a zero value)", view.Voice.AecEnabled)
	}
	if view.Voice.InputDevice == nil || *view.Voice.InputDevice != "Blue Yeti" {
		t.Fatalf("InputDevice = %v, want \"Blue Yeti\"", view.Voice.InputDevice)
	}
	if view.Voice.SampleRate == nil || *view.Voice.SampleRate != 16000 {
		t.Fatalf("SampleRate = %v, want the chain's 16000", view.Voice.SampleRate)
	}

	// The live-reader half: hand-edit the file, one reload poll, the SAME
	// store's next read answers the new values - no restart anywhere.
	edited := strings.Replace(localAudioRBase, "enabled = false", "enabled = true", 1)
	edited = strings.Replace(edited, "input_device = \"Blue Yeti\"", "input_device = \"Realtek Mic Array\"", 1)
	if err := os.WriteFile(path, []byte(edited), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := mgr.CheckAndReload(); err != nil {
		t.Fatalf("CheckAndReload: %v", err)
	}
	view2, err := s.ReadSettings(ctx)
	if err != nil {
		t.Fatalf("ReadSettings after reload: %v", err)
	}
	if view2.Voice.AecEnabled == nil || !*view2.Voice.AecEnabled {
		t.Fatalf("AecEnabled after a hot-applied edit = %v, want true: [voice] aec.enabled has no live snapshot reader", view2.Voice.AecEnabled)
	}
	if view2.Voice.InputDevice == nil || *view2.Voice.InputDevice != "Realtek Mic Array" {
		t.Fatalf("InputDevice after a hot-applied edit = %v, want \"Realtek Mic Array\": [audio] has no live snapshot reader", view2.Voice.InputDevice)
	}
}

// TestConfigLocalAudioWriteFieldsRoundTripAndRefuse: the two whitelist names
// route to the guarded setters, the file carries the dotted keys, the reload
// turns the write into a changed live answer, and the refusal paths move no
// byte. audio_sample_rate stays guard-only: no write leg exists for it.
func TestConfigLocalAudioWriteFieldsRoundTripAndRefuse(t *testing.T) {
	s, mgr, path := localAudioRStore(t)
	ctx := context.Background()

	res, err := s.ApplySetting(ctx, panel.SettingWrite{
		Field: panel.FieldVoiceAecEnabled, Value: "true", Actor: "test-local-audio",
	})
	if err != nil {
		t.Fatalf("ApplySetting(voice_aec_enabled=true): %v", err)
	}
	if len(res.Written) != 1 || res.Written[0] != "voice.aec.enabled" {
		t.Fatalf("receipt written = %v, want [voice.aec.enabled]", res.Written)
	}

	res, err = s.ApplySetting(ctx, panel.SettingWrite{
		Field: panel.FieldAudioInputDevice, Value: "Realtek Mic Array", Actor: "test-local-audio",
	})
	if err != nil {
		t.Fatalf("ApplySetting(audio_input_device): %v", err)
	}
	if len(res.Written) != 1 || res.Written[0] != "audio.input_device" {
		t.Fatalf("receipt written = %v, want [audio.input_device]", res.Written)
	}

	if _, err := mgr.CheckAndReload(); err != nil {
		t.Fatalf("CheckAndReload after the writes: %v", err)
	}
	view, err := s.ReadSettings(ctx)
	if err != nil {
		t.Fatalf("ReadSettings after the writes: %v", err)
	}
	if view.Voice.AecEnabled == nil || !*view.Voice.AecEnabled {
		t.Fatalf("AecEnabled after the write = %v, want true: the write did not reach the live reader", view.Voice.AecEnabled)
	}
	if view.Voice.InputDevice == nil || *view.Voice.InputDevice != "Realtek Mic Array" {
		t.Fatalf("InputDevice after the write = %v, want \"Realtek Mic Array\": the write did not reach the live reader", view.Voice.InputDevice)
	}

	// The refusals: a non-bool word is refused by the leg's shape check and
	// the guard-key name falls through to the default branch; both leave the
	// file byte-identical.
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.ApplySetting(ctx, panel.SettingWrite{
		Field: panel.FieldVoiceAecEnabled, Value: "maybe", Actor: "test-local-audio",
	}); err == nil {
		t.Fatal("voice_aec_enabled=maybe was accepted: the leg must refuse a non-bool word")
	}
	if _, err := s.ApplySetting(ctx, panel.SettingWrite{
		Field: "audio_sample_rate", Value: "48000", Actor: "test-local-audio",
	}); err == nil || !strings.Contains(err.Error(), "没有对应的写入腿") {
		t.Fatalf("audio_sample_rate returned %v, want the default-branch refusal: the guard key has no write leg", err)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("a refused write moved bytes in config.toml")
	}
}
