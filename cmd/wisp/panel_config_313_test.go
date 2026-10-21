package main

// Ticket 313 at the assembly root: the settings leg now answers the picker's read
// payload, and the six new write fields reach their real setters.
//
// WHAT DRIVES THIS FILE. Everything below runs the production assembly the same way
// panel_config_248_test.go does - one config.Manager, one secret.Store, the one
// panel.ComposerDispatch the WebView2 host and `wisp panel-inbound` both call - so
// "the write landed" means bytes in config.toml read back through the same leg, not
// a fake agreeing with itself. The only fake is the door-level roster spy in the
// union pin, which is about the roster and nothing about a write.
//
// THE A881 CONDITIONS THIS LEG OWES, and where each is answered:
//   a. the ticket-248 table and its frozen pin are untouched - this leg opens
//      neither file; the live evidence is the probe logs' git diff over both (0).
//   b. the union is computed in one accessor (config_write_313.go's
//      acceptedWritableFields); the door's refusal text is its only externally
//      readable copy, and the pin below reads it live.
//   c. the union-equals-table1-union-table2 pin itself: 313-r1 shipped the accessor
//      with no test reader of the union (its config_313_test.go pins disjointness,
//      not equality), so the pin is added here, in this leg's new file.

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
)

// model313 is a made-up catalog id. A877 ruling 4 forbids inventing real model
// names, and this leg does not either: the catalog here is the test's own row, so
// the id only has to be legal.
const model313 = "model-under-test-313"

// seed313Empty is the state ticket 313 exists for: a config.toml with no provider
// rows and no role pair, so the page has nothing to address until it creates the
// row (A877 ruling 2).
const seed313Empty = `schema_version = 2
# ticket 313 end-to-end: no [llm] section at all - zero provider rows.
`

func empty313Dir(t *testing.T) string {
	t.Helper()
	dir := sealableTempDir124(t)
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(seed313Empty), 0o600); err != nil {
		t.Fatalf("seed config.toml: %v", err)
	}
	return dir
}

// settingsPayload313 reads the machine line the read side appends to config.get.
func settingsPayload313(t *testing.T, reply string) map[string]any {
	t.Helper()
	idx := strings.Index(reply, panel.SettingsPayloadMarker)
	if idx < 0 {
		t.Fatalf("config.get carries no payload line: %q", reply)
	}
	line := strings.SplitN(reply[idx+len(panel.SettingsPayloadMarker):], "\n", 2)[0]
	var payload map[string]any
	if err := json.Unmarshal([]byte(line), &payload); err != nil {
		t.Fatalf("the payload after the marker is not JSON: %v", err)
	}
	return payload
}

// apply313 drives one config.set through the production dispatch and asserts the
// receipt names the key paths that write owns and keeps the restart wording.
func apply313(t *testing.T, disp *panel.ComposerDispatch, env string, want ...string) string {
	t.Helper()
	reply, err := disp.Handle(context.Background(), env)
	if err != nil {
		t.Fatalf("config.set was refused: %v\n%s", err, env)
	}
	if !strings.Contains(reply, "重启") || strings.Contains(reply, "立即生效") || strings.Contains(reply, "保存即生效") {
		t.Errorf("receipt lost the restart wording or claims a live save: %q", reply)
	}
	for _, w := range want {
		if !strings.Contains(reply, w) {
			t.Errorf("receipt does not name %q: %q", w, reply)
		}
	}
	return reply
}

// set313Env builds one settings envelope; every selector the six fields can need is
// always present - an extra one is ignored, a missing one would refuse.
func set313Env(id, field, provider, model, value string) string {
	raw, err := json.Marshal(map[string]string{
		"method": "config.set", "requestId": id, "source": "panel-composer", "configField": field,
		"configProvider": provider, "configModel": model, "configValue": value,
	})
	if err != nil {
		panic(err)
	}
	return string(raw)
}

// AC#1/AC#2 through the real leg: an empty config becomes a stepfun row, the two
// cloud chains, the realtime head and the chat pair - one receipt per field, each
// naming its own key path - and the read side then answers the picker with what is
// now true.
func Test313EmptyConfigToStepfunRowChainAndPair(t *testing.T) {
	dir := empty313Dir(t)
	var audit []string
	disp, _ := newRealSettingsLeg(t, dir, &audit, nil)
	ctx := context.Background()

	// Before: presets answered, "nothing configured" answered as a known false.
	reply, err := disp.Handle(ctx, `{"method":"config.get","requestId":"r-313-get-0","source":"panel-composer"}`)
	if err != nil {
		t.Fatalf("config.get on an empty config: %v", err)
	}
	payload := settingsPayload313(t, reply)
	if payload["presetsKnown"] != true || payload["configuredKnown"] != true || payload["configured"] != false {
		t.Errorf("the empty-config flags are wrong: %v", payload)
	}
	presets, _ := payload["presets"].([]any)
	var gotPresets []string
	for _, p := range presets {
		gotPresets = append(gotPresets, p.(string))
	}
	if want := config.PresetNames(); !reflect.DeepEqual(gotPresets, want) {
		t.Errorf("presets = %v, want the built-in list %v", gotPresets, want)
	}
	// Zero rows travel as an empty array, not null and not absent: the settings
	// page reads providers with unguarded find/some at render (ticket 313's
	// acceptance leg 313-v1 measured the old null shape against the live page;
	// it throws exactly on the fresh-config state AC#1 targets).
	if rows, ok := payload["providers"].([]any); !ok || len(rows) != 0 {
		t.Errorf("zero provider rows must render as the empty array [], got %v", payload["providers"])
	}

	// The six writes in dependency order (the row comes first: both the chains and
	// the pair name the model it creates).
	apply313(t, disp, set313Env("r-313-set-1", panel.FieldProviderRowCreate, "stepfun", model313, ""), "llm.providers.stepfun")
	apply313(t, disp, set313Env("r-313-set-2", panel.FieldVoiceCloudASRChain, "stepfun", model313, "stepfun/"+model313), "voice.cloud_asr_chain")
	apply313(t, disp, set313Env("r-313-set-3", panel.FieldVoiceCloudTTSChain, "stepfun", model313, "stepfun/"+model313), "voice.cloud_tts_chain")
	apply313(t, disp, set313Env("r-313-set-4", panel.FieldVoiceRealtimeProvider, "stepfun", model313, "stepfun"), "voice.realtime.provider")
	apply313(t, disp, set313Env("r-313-set-5", panel.FieldVoiceRealtimeModel, "stepfun", model313, model313), "voice.realtime.model")
	// The one field that owns two dotted keys (A881 ruling 1): both must be named.
	apply313(t, disp, set313Env("r-313-set-6", panel.FieldRoleChatPair, "stepfun", model313, ""),
		"llm.roles.chat.provider", "llm.roles.chat.model")

	// On disk, not only in memory.
	body, err := os.ReadFile(filepath.Join(dir, configFileName))
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	for _, frag := range []string{
		"[llm.providers.stepfun]", "[llm.roles.chat]", "cloud_asr_chain", "cloud_tts_chain",
		"voice.realtime", "stepfun/" + model313,
	} {
		if !strings.Contains(string(body), frag) {
			t.Errorf("config.toml does not carry %q:\n%s", frag, body)
		}
	}

	// Read back: every one of the six writes is visible to the page's data source.
	reply, err = disp.Handle(ctx, `{"method":"config.get","requestId":"r-313-get-1","source":"panel-composer"}`)
	if err != nil {
		t.Fatalf("config.get after the writes: %v", err)
	}
	payload = settingsPayload313(t, reply)
	if payload["configured"] != true {
		t.Errorf("after the row write configured=%v, want true", payload["configured"])
	}
	if payload["chatProvider"] != "stepfun" || payload["chatModel"] != model313 {
		t.Errorf("chat pair read back as %v/%v", payload["chatProvider"], payload["chatModel"])
	}
	rows, _ := payload["providers"].([]any)
	if len(rows) != 1 {
		t.Fatalf("providers read back as %v", payload["providers"])
	}
	row := rows[0].(map[string]any)
	models, _ := row["models"].([]any)
	if row["name"] != "stepfun" || len(models) != 1 || models[0] != model313 || row["modelCount"] != float64(1) {
		t.Errorf("provider row read back as %v", row)
	}
	voice, _ := payload["voice"].(map[string]any)
	asr, _ := voice["cloudASRChain"].([]any)
	tts, _ := voice["cloudTTSChain"].([]any)
	if len(asr) != 1 || asr[0] != "stepfun/"+model313 || len(tts) != 1 || tts[0] != "stepfun/"+model313 {
		t.Errorf("voice chains read back as %v / %v", voice["cloudASRChain"], voice["cloudTTSChain"])
	}
	if voice["realtimeProvider"] != "stepfun" || voice["realtimeModel"] != model313 {
		t.Errorf("realtime head read back as %v/%v", voice["realtimeProvider"], voice["realtimeModel"])
	}
	if voice["realtimeEnabled"] != false {
		t.Errorf("realtime enabled flipped without a write: %v", voice["realtimeEnabled"])
	}
}

// rosterSpy313 counts what the door handed down. The union pin is about the roster
// the door accepts, not about any write landing, so the leg just counts.
type rosterSpy313 struct{ calls int }

func (l *rosterSpy313) ReadSettings(context.Context) (panel.SettingsView, error) {
	return panel.SettingsView{Readable: true}, nil
}

func (l *rosterSpy313) ApplySetting(_ context.Context, w panel.SettingWrite) (panel.SettingWriteResult, error) {
	l.calls++
	return panel.SettingWriteResult{Field: w.Field, Tier: panel.EffectiveRestart}, nil
}

func (l *rosterSpy313) StoreCredential(context.Context, string) (panel.SettingWriteResult, error) {
	l.calls++
	return panel.SettingWriteResult{Field: panel.FieldProviderCredential, Tier: panel.EffectiveRestart}, nil
}

// A881 condition c: one pin over the union. 313-r1 shipped acceptedWritableFields
// (condition b) but no test read the union itself, so this is the pin the ruling
// asked the landing leg for. Two halves: the union the door REPORTS is exactly
// WritableFields() + AdditionalWritableFields(), and every member of that union is
// then ACCEPTED by the door and reaches the leg exactly once - so a one-sided field
// drop moves this test whichever table it happens in.
func Test313WritableUnionEqualsTheTwoRosters(t *testing.T) {
	spy := &rosterSpy313{}
	audit := func(string, ...any) {}
	disp := &panel.ComposerDispatch{
		Config: &panel.ConfigWriteHandler{Store: spy, Audit: audit, Actor: "test-313"},
		Audit:  audit,
	}

	// Half 1: the refusal's list is the union accessor, read live.
	_, err := disp.Handle(context.Background(),
		`{"method":"config.set","requestId":"r-313-union","source":"panel-composer","configField":"no_such_field_313"}`)
	if !errors.Is(err, panel.ErrConfigFieldRefused) {
		t.Fatalf("an unlisted name returned %v, want %v", err, panel.ErrConfigFieldRefused)
	}
	const head = "可写的只有 "
	i := strings.Index(err.Error(), head)
	if i < 0 {
		t.Fatalf("the refusal does not name the accepted list: %v", err)
	}
	rest := err.Error()[i+len(head):]
	if j := strings.Index(rest, "（requestId="); j >= 0 {
		rest = rest[:j]
	}
	reported := strings.Split(rest, "/")
	want := append(append([]string{}, panel.WritableFields()...), panel.AdditionalWritableFields()...)
	sort.Strings(want)
	if !reflect.DeepEqual(reported, want) {
		t.Errorf("the door reports the union %v; the two rosters say %v (A881: union == table1 + table2)", reported, want)
	}

	// Half 2: every union member really reaches the leg through the door, once.
	for _, field := range want {
		spy.calls = 0
		reply, herr := disp.Handle(context.Background(), set313Env("r-313-union-"+field, field, "stepfun", "m-313", "x"))
		if herr != nil {
			t.Errorf("union member %q was refused by the door: %v", field, herr)
			continue
		}
		if spy.calls != 1 {
			t.Errorf("union member %q reached the leg %d time(s), want exactly 1 (reply %q)", field, spy.calls, reply)
		}
	}
}

// 313-r3: the comma split of a chain value, pinned at the seam that performs it.
//
// WHAT WAS MISSING. The panel sends a whole chain as ONE string (configValue), and
// panel_config_store.go's splitChainValue is the only place in the repository that
// turns it into the element slice the config setter takes. Until this pin, every
// door-level chain write in the tests carried a single element, so 313-v1's
// mutation - splitting on ";" instead of "," - left all three Test313 packages
// green: a two-element chain would have been refused (the whole comma string read
// as one model id) with nothing in the repository to say so. This is that reader.
//
// THE NAMES. The two model ids below are the ones the repository's own fixtures
// already spell (settingsSeed's "deepseek-chat"; "deepseek-reasoner" is the second
// model the 313/257 seeds carry, e.g. internal/config's enabled-head variant). No
// name is invented here.
//
// THE VALUE carries the space a user types after the comma on purpose:
// splitChainValue trims each element (its own contract), and the anchored line
// below is what keeps that trimming from being dropped silently along with the
// split.
func Test313MultiElementCommaChainKeepsEveryElement(t *testing.T) {
	// settingsSeed declares deepseek/deepseek-chat; the appended entry gives the
	// chain a second legal fallback head.
	const seed = settingsSeed + `
[llm.providers.deepseek.models.deepseek-reasoner]
enabled = true
context_window = 32000
`
	dir := sealableTempDir124(t)
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(seed), 0o600); err != nil {
		t.Fatalf("seed config.toml: %v", err)
	}
	var audit []string
	disp, _ := newRealSettingsLeg(t, dir, &audit, nil)

	// One config.set for the whole chain; apply313 carries the receipt-side
	// requirement (names voice.cloud_asr_chain, keeps the restart wording).
	const chainValue = "deepseek/deepseek-chat, deepseek/deepseek-reasoner"
	apply313(t, disp, set313Env("r-313-2elt-asr", panel.FieldVoiceCloudASRChain, "", "", chainValue), "voice.cloud_asr_chain")

	// Surface 1 - the file: both elements, in order, trimmed. Anchored to the array
	// line itself because the seed's own model section also spells
	// "deepseek-reasoner": a fragment search over the whole body would pass even if
	// the chain landed half of itself.
	body, err := os.ReadFile(filepath.Join(dir, configFileName))
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	const wantLine = "cloud_asr_chain = ['deepseek/deepseek-chat', 'deepseek/deepseek-reasoner']"
	if !strings.Contains(string(body), wantLine) {
		t.Errorf("config.toml does not carry the two comma-separated elements of %q as %q:\n%s", chainValue, wantLine, body)
	}

	// Surface 2 - the read side the picker calls: the same two elements, in order.
	reply, err := disp.Handle(context.Background(), `{"method":"config.get","requestId":"r-313-2elt-get","source":"panel-composer"}`)
	if err != nil {
		t.Fatalf("config.get after the chain write: %v", err)
	}
	payload := settingsPayload313(t, reply)
	voice, _ := payload["voice"].(map[string]any)
	asr, _ := voice["cloudASRChain"].([]any)
	if len(asr) != 2 || asr[0] != "deepseek/deepseek-chat" || asr[1] != "deepseek/deepseek-reasoner" {
		t.Errorf("voice.cloudASRChain read back as %v, want both elements of %q in order", voice["cloudASRChain"], chainValue)
	}
	// The write owns the ASR chain alone.
	if tts, _ := voice["cloudTTSChain"].([]any); len(tts) != 0 {
		t.Errorf("the ASR-only write moved the TTS chain: %v", voice["cloudTTSChain"])
	}
}
