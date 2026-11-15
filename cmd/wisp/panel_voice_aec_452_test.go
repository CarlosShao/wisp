package main

// Ticket 452 AC#2's panel half, at the assembly root: the settings payload must
// tell the page that THIS BINARY cannot cancel echo, on a key of its own, beside
// (not instead of) the user's switch.
//
// WHY THIS TEST IS THE TEETH AND NOT A SOURCE-STRING PIN (A1880). It builds the
// production leg (config.Manager + secret.Store + panel.ComposerDispatch, via
// newRealSettingsLeg), drives `config.get` through the door, and reads the JSON
// that actually arrives at the page. Delete the `AecBuildAvailable:
// &aecBuildAvailable` line in cmd/wisp/panel_config_store.go's ReadSettings and
// the key leaves the payload - this fails on the wire, not on a grep of the
// source. Flip the bit to come from config (i.e. fold it into AecEnabled) and
// the write half below turns red, because the user switching AEC on would drag
// the capability answer with it.
//
// THE READINGS ARE COMPARED AGAINST aec.BindingBuilt(), never against a literal:
// on a `-tags webrtcap` build the same test must be green with true. The host
// this ticket was written on has no binding (scripts/build.ps1 passes no -tags),
// so the value here is false - which is exactly why the assertion cannot be
// "the key reads false" (a zero value would pass it in silence) but must be
// "the key is PRESENT and equals the build's own answer".

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/audio/aec"
	"github.com/CarlosShao/wisp/internal/panel"
)

// aec452Seed is settingsSeed with AEC switched OFF by hand: "off" is then a real
// answer a zero value could not fake, and the write below flips it to on.
const aec452Seed = settingsSeed + `
[voice.aec]
enabled = false
`

func TestAC452PanelVoiceTellsThePageWhatTheBuildCanDo(t *testing.T) {
	dir := sealableTempDir124(t)
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(aec452Seed), 0o600); err != nil {
		t.Fatalf("seed config.toml: %v", err)
	}
	var audit []string
	disp, _ := newRealSettingsLeg(t, dir, &audit, nil)
	ctx := context.Background()

	wantCapability := aec.BindingBuilt()

	// Half 1: the answered read carries the capability, beside a switch the user
	// left off.
	reply, err := disp.Handle(ctx, `{"method":"config.get","requestId":"r-452-get-0","source":"panel-composer"}`)
	if err != nil {
		t.Fatalf("config.get: %v", err)
	}
	voice := voicePayload452(t, reply)
	if got, ok := voice["aecEnabled"]; !ok || got != false {
		t.Fatalf("voice.aecEnabled = %v (present=%v), want the answered false from config.toml: the writable switch half is not on the wire", got, ok)
	}
	capability, ok := voice["aecBuildAvailable"]
	if !ok {
		t.Fatalf("the settings payload's voice half carries NO aecBuildAvailable key: the page cannot tell that this build cannot cancel echo, so it prints the user's switch as if the capability existed - which is ticket 452 AC#2's undone half (voice = %v)", voice)
	}
	if capability != wantCapability {
		t.Fatalf("voice.aecBuildAvailable = %v, want %v (= aec.BindingBuilt(), the build's own answer): the panel is answering something other than what this binary can do", capability, wantCapability)
	}

	// Half 2: the user switches AEC on through the same door. The switch moves;
	// the capability does not - it is not a config value and must not be derived
	// from one.
	apply313(t, disp, set313Env("r-452-set-aec", panel.FieldVoiceAecEnabled, "", "", "true"), "voice.aec.enabled")

	reply, err = disp.Handle(ctx, `{"method":"config.get","requestId":"r-452-get-1","source":"panel-composer"}`)
	if err != nil {
		t.Fatalf("config.get after the AEC write: %v", err)
	}
	voice = voicePayload452(t, reply)
	if got := voice["aecEnabled"]; got != true {
		t.Fatalf("voice.aecEnabled = %v after the write, want true: the switch is not a live reader, so this half of the test cannot prove independence", got)
	}
	if got, ok := voice["aecBuildAvailable"]; !ok || got != wantCapability {
		t.Fatalf("voice.aecBuildAvailable = %v (present=%v) after the user switched AEC on, want %v: the capability bit is being derived from the user's switch instead of the build, which is the one-boolean shape ticket 452 forbids", got, ok, wantCapability)
	}
}

// voicePayload452 pulls the voice object out of one config.get reply's machine
// line. It fails loudly if the marker or the object is missing, so the two
// assertions above never rest on a silently empty map.
func voicePayload452(t *testing.T, reply string) map[string]any {
	t.Helper()
	idx := strings.Index(reply, panel.SettingsPayloadMarker)
	if idx < 0 {
		t.Fatalf("config.get carries no settings payload line: %q", reply)
	}
	line := strings.SplitN(reply[idx+len(panel.SettingsPayloadMarker):], "\n", 2)[0]
	var payload map[string]any
	if err := json.Unmarshal([]byte(line), &payload); err != nil {
		t.Fatalf("the payload after the marker is not JSON: %v", err)
	}
	raw, ok := payload["voice"]
	if !ok {
		t.Fatalf("the payload carries no voice object: %q", line)
	}
	voice, ok := raw.(map[string]any)
	if !ok {
		t.Fatalf("the payload's voice half is not an object: %v", raw)
	}
	return voice
}
