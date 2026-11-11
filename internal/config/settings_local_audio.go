package config

// The 本机音频 setters (the 2026-11-11 S-leg snapshot/write alignment, leg R):
// the two leaves the settings page's local-audio card can write.
//
// WHY THEY EXIST. The backend grew real readers/guards for voice.aec.enabled
// (the barge-in monitor's config switch, fed at the resident boot plan),
// audio.input_device (the capture open's device resolver) and
// audio.sample_rate (the load-time 16000 chain guard) before the panel had any
// channel for them; the frontend card honestly rendered a WAITING state instead
// of guessing. These two setters are the write half of that channel - the
// snapshot half lives in cmd/wisp/panel_config_store.go's ReadSettings and
// internal/panel's voicePayload (reader-gated keys).
//
// THE SINGLE SOURCE OF REFUSAL. Both ride writeOneKey like every other
// settings setter: validate() runs on the base the write would actually
// produce before a byte is written, and neither key adds a second vocabulary
// on top. voice.aec.enabled is a bool with no band; audio.input_device is a
// free name whose real gate is at open time (internal/audio's
// selectCaptureDevice refuses an unmatched name loudly and lists the present
// endpoints) - device presence is a boot-time fact, not a config-layer
// vocabulary, so pre-judging it here would be a second gate that can lie.
//
// WHY THERE IS NO SetAudioSampleRate. audio.sample_rate is not a setting, it
// is a load-time guard: validateAudio refuses every value except the chain's
// own 16000 (the capture resample target, the frame geometry and the VAD/ASR
// engines' contract are all baked to that one rate), so a settings field for
// it could only ever write the number the file already carries or be refused
// by that same sentence. A control that is either a no-op or a guaranteed
// refusal is ticket 83's lie wearing a slider; the snapshot carries the fact
// and the panel door refuses the field name (pinned in internal/panel's
// config_313_test.go).

import "strings"

// SetVoiceAecEnabled writes voice.aec.enabled (the barge-in monitor's switch).
// The value is the boolean itself; the guarded write path is the only gate.
func (m *Manager) SetVoiceAecEnabled(on bool) ([]string, error) {
	return m.writeOneKey("voice.aec.enabled", func(base *Config) bool {
		base.Voice.AEC.Enabled = on
		return true
	}, nil)
}

// SetAudioInputDevice writes audio.input_device (the capture open's device
// name). "default" and the empty string both mean "follow the OS default" per
// internal/audio's own convention, so neither is refused here; any other name
// is written as given (trimmed) and resolved - or loudly refused - at the
// next capture open.
func (m *Manager) SetAudioInputDevice(name string) ([]string, error) {
	v := strings.TrimSpace(name)
	return m.writeOneKey("audio.input_device", func(base *Config) bool {
		base.Audio.InputDevice = v
		return true
	}, nil)
}
