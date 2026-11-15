package panel

// Ticket 452 AC#2, the panel half: the voice section of the settings payload
// carries TWO separate facts about echo cancellation, and this file pins that
// they cannot be folded back into one.
//
// WHY A SECOND FIELD EXISTS AT ALL. The shipped build passes no `-tags
// webrtcap` to `go build` (scripts/build.ps1), so internal/audio/aec compiles
// to its stub and aec.BindingBuilt() answers false: that binary cannot cancel
// echo. Meanwhile voice.aec.enabled stays a user-writable switch, and its
// default is on. A page that only had the switch would print "回声消除：开" on a
// build that cannot do it - the mismatch ticket 452 exists to stop. So the
// payload answers both questions on separate keys:
//
//	voice.aecEnabled       - does the user want AEC (config value, writable)
//	voice.aecBuildAvailable - can this binary cancel echo (compile-time fact)
//
// The cases below drive renderSettingsPayload - the same function the settings
// page's data source is built from - and read the JSON back, so they fail if
// either key is renamed, merged or silently dropped. The three-state gate
// (nil = the leg did not answer) is pinned in the last case, the same
// convention the 本机音频 half already uses.

import (
	"encoding/json"
	"strings"
	"testing"
)

// voiceOf452 parses the voice object out of one rendered settings payload.
func voiceOf452(t *testing.T, v SettingsView) map[string]json.RawMessage {
	t.Helper()
	line := renderSettingsPayload(v)
	if !strings.HasPrefix(line, SettingsPayloadMarker) {
		t.Fatalf("the payload line does not start with the marker: %q", line)
	}
	var top map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimPrefix(line, SettingsPayloadMarker)), &top); err != nil {
		t.Fatalf("payload: %v", err)
	}
	raw, ok := top["voice"]
	if !ok {
		t.Fatalf("the payload carries no voice object at all: %q", line)
	}
	var voice map[string]json.RawMessage
	if err := json.Unmarshal(raw, &voice); err != nil {
		t.Fatalf("voice: %v", err)
	}
	return voice
}

func Test452VoiceCapabilityAndTheUserSwitchAreTwoIndependentKeys(t *testing.T) {
	cannot, can, on, off := false, true, true, false

	// 1. THE CASE THE TICKET IS ABOUT: the user wants AEC on, the build cannot
	//    cancel echo. Both facts must reach the page, on two keys, with two
	//    different values - this is the exact combination that a single merged
	//    boolean would destroy.
	v := voiceOf452(t, SettingsView{
		Readable: true,
		Voice:    VoiceSettingsView{AecEnabled: &on, AecBuildAvailable: &cannot},
	})
	if got := string(v["aecEnabled"]); got != "true" {
		t.Fatalf("voice.aecEnabled = %s, want the user's answered true: the writable switch disappeared", got)
	}
	got, ok := v["aecBuildAvailable"]
	if !ok {
		t.Fatalf("the voice half carries no aecBuildAvailable key: the page cannot say that this build cannot cancel echo, so \"回声消除：开\" prints on a binary with no AEC binding (voice keys: %v)", keysOf452(v))
	}
	if string(got) != "false" {
		t.Errorf("voice.aecBuildAvailable = %s, want the answered false (the shipped build has no -tags webrtcap binding): the capability bit is reporting something other than the build's own answer", got)
	}

	// 2. THE REVERSE COMBINATION: a capable build with the user's switch off.
	//    Together with case 1 this is the two-way proof that neither key is
	//    derived from the other - flipping one leaves the other where it was.
	v = voiceOf452(t, SettingsView{
		Readable: true,
		Voice:    VoiceSettingsView{AecEnabled: &off, AecBuildAvailable: &can},
	})
	if got := string(v["aecEnabled"]); got != "false" {
		t.Errorf("voice.aecEnabled = %s, want the answered false", got)
	}
	if got := string(v["aecBuildAvailable"]); got != "true" {
		t.Errorf("voice.aecBuildAvailable = %s, want the answered true for a capable build: it is being derived from the user's switch instead of the build", got)
	}

	// 3. THE GATE: a leg that answered nothing puts neither claim on the wire.
	//    "absent" (nobody asked) and "false" (asked, the build cannot) are
	//    different states and must not render the same way; production fills the
	//    capability on every answered read (cmd/wisp's ReadSettings), so absent
	//    means a leg-side defect, never a claim about the binary.
	v = voiceOf452(t, SettingsView{Readable: true})
	for _, k := range []string{"aecEnabled", "aecBuildAvailable"} {
		if _, ok := v[k]; ok {
			t.Errorf("a leg that never answered carries %q: a zero value must not ride the wire as a fact (voice keys: %v)", k, keysOf452(v))
		}
	}
}

// keysOf452 lists the voice keys of one payload, for the failure messages above
// (a red that says what it DID see saves the next reader a rerun).
func keysOf452(v map[string]json.RawMessage) []string {
	out := make([]string, 0, len(v))
	for k := range v {
		out = append(out, k)
	}
	return out
}
