package panel

// Ticket 313's panel-side half: the second writable roster, and the structured
// payload the settings page reads off config.get.
//
// Every judgement here is a call count, an error identity or the payload's own
// keys - never a log string (the rule config_route_248_test.go states, inherited).
// The frozen ticket-248 pin lives in that file and is NOT touched: this file
// asserts the two rosters stay disjoint, which is the property that keeps the pin
// meaning what it meant.

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// fields313Case is one 313 field with the selectors it needs to clear the gate.
type fields313Case struct {
	field string
	extra map[string]any
}

func fields313Cases() []fields313Case {
	return []fields313Case{
		{field: FieldVoiceCloudASRChain, extra: map[string]any{"configValue": "stepfun/model-under-test-313"}},
		{field: FieldVoiceCloudTTSChain, extra: map[string]any{"configValue": "stepfun/model-under-test-313"}},
		{field: FieldVoiceRealtimeProvider, extra: map[string]any{"configValue": "stepfun"}},
		{field: FieldVoiceRealtimeModel, extra: map[string]any{"configValue": "model-under-test-313"}},
		{field: FieldProviderRowCreate, extra: map[string]any{"configProvider": "stepfun"}},
		{field: FieldProviderBilling, extra: map[string]any{"configProvider": "stepfun", "configValue": "plan"}},
		{field: FieldProviderRowCreateCustom, extra: map[string]any{
			"configProvider": "local-gw", "configValue": "openai-chat,http://127.0.0.1:8080/v1",
		}},
		{field: FieldModelRowCreate, extra: map[string]any{"configProvider": "stepfun", "configModel": "model-under-test-313"}},
		{field: FieldRoleChatPair, extra: map[string]any{"configProvider": "stepfun", "configModel": "model-under-test-313"}},
	}
}

// AC#2's panel-side half: every 313 field clears the door and reaches the leg
// exactly once, and the ticket-248 pin still describes the ticket-248 seven.
func Test313AdditionalFieldsReachTheLegAndThe248PinIsUntouched(t *testing.T) {
	// The pin's meaning, restated here as a positive assertion: the 248 roster is
	// the 248 seven, exactly.
	want248 := []string{
		FieldModelContextWindow, FieldModelPriceIn, FieldModelPriceOut,
		FieldProviderAPIKeyRef, FieldProviderBaseURL, FieldProviderCredential, FieldRoleChatModel,
	}
	if got := WritableFields(); !reflect.DeepEqual(got, want248) {
		t.Fatalf("WritableFields() = %v, want the ticket-248 seven %v: the frozen pin's meaning moved", got, want248)
	}
	listed := map[string]bool{}
	for _, f := range WritableFields() {
		listed[f] = true
	}
	for _, f := range AdditionalWritableFields() {
		if listed[f] {
			t.Errorf("%q is in both rosters: the 313 additions must not restate a 248 name", f)
		}
	}

	for _, tc := range fields313Cases() {
		leg := &cfg248Leg{res: SettingWriteResult{Field: tc.field, Tier: EffectiveRestart}}
		env := map[string]any{
			"method": MethodConfigSet, "requestId": "r-313-" + tc.field, "source": ComposerRequestSource,
			"configField": tc.field,
		}
		for k, v := range tc.extra {
			env[k] = v
		}
		if _, err := cfg248Dispatch(leg).Handle(context.Background(), cfg248Raw(t, env)); err != nil {
			t.Errorf("field %q was refused by the door: %v", tc.field, err)
		}
		if leg.applyCalls != 1 {
			t.Errorf("field %q reached the leg %d time(s), want 1", tc.field, leg.applyCalls)
		}
		if leg.lastWrite.Field != tc.field {
			t.Errorf("the leg received field %q, want %q", leg.lastWrite.Field, tc.field)
		}
	}

	// Negative with its control: a name on neither roster never reaches the leg,
	// and the refusal names the union so the page is not sent to half a list.
	leg := &cfg248Leg{}
	_, err := cfg248Dispatch(leg).Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigSet, "requestId": "r-313-nope", "source": ComposerRequestSource,
		"configField": "voice_cloud_tts_chain_nope", "configValue": "x",
	}))
	if !errors.Is(err, ErrConfigFieldRefused) {
		t.Errorf("an unlisted 313-shaped name returned %v, want %v", err, ErrConfigFieldRefused)
	}
	if err != nil && !strings.Contains(err.Error(), FieldVoiceCloudTTSChain) {
		t.Errorf("the refusal does not list the accepted backlog: %v", err)
	}
	if leg.applyCalls+leg.readCalls+leg.credCalls != 0 {
		t.Errorf("a refused name still reached the leg: %+v", leg)
	}
}

// The selector columns of the new roster, each with its missing-selector case and
// the same shape with the selector present as the control.
func Test313SelectorRefusalsOnTheAdditionalRoster(t *testing.T) {
	cases := []struct {
		name string
		env  map[string]any
	}{
		{"provider_row_create without a provider", map[string]any{"configField": FieldProviderRowCreate}},
		{"role_chat_pair without a provider", map[string]any{"configField": FieldRoleChatPair, "configModel": "m"}},
		{"role_chat_pair without a model", map[string]any{"configField": FieldRoleChatPair, "configProvider": "stepfun"}},
		{"voice_cloud_asr_chain with an empty value", map[string]any{"configField": FieldVoiceCloudASRChain, "configValue": "   "}},
		{"provider_billing without a provider", map[string]any{"configField": FieldProviderBilling, "configValue": "plan"}},
		{"provider_billing with an empty value", map[string]any{"configField": FieldProviderBilling, "configProvider": "stepfun", "configValue": "   "}},
		{"provider_row_create_custom without a provider", map[string]any{
			"configField": FieldProviderRowCreateCustom, "configValue": "openai-chat,http://127.0.0.1:8080/v1",
		}},
		{"provider_row_create_custom with an empty value", map[string]any{
			"configField": FieldProviderRowCreateCustom, "configProvider": "local-gw", "configValue": "   ",
		}},
		{"model_row_create without a provider", map[string]any{
			"configField": FieldModelRowCreate, "configModel": "m-320",
		}},
		{"model_row_create without a model", map[string]any{
			"configField": FieldModelRowCreate, "configProvider": "stepfun",
		}},
	}
	for _, tc := range cases {
		leg := &cfg248Leg{}
		env := map[string]any{"method": MethodConfigSet, "requestId": "r-313-sel", "source": ComposerRequestSource}
		for k, v := range tc.env {
			env[k] = v
		}
		if _, err := cfg248Dispatch(leg).Handle(context.Background(), cfg248Raw(t, env)); !errors.Is(err, ErrConfigFieldRefused) {
			t.Errorf("%s: got %v, want %v", tc.name, err, ErrConfigFieldRefused)
		}
		if leg.applyCalls+leg.credCalls != 0 {
			t.Errorf("%s reached the leg anyway", tc.name)
		}
	}
}

// The naming bans, machine-readable over the new roster: no 313 field carries an
// approval word, and none reaches into a locked family. The control in the same
// test is a locked-family name, which must still be refused.
func Test313AdditionalFieldNamesCarryNoApprovalWordAndNoLockedFamily(t *testing.T) {
	grantWords := []string{
		"approve", "approval", "grant", "allow", "permit",
		"ratify", "authorize", "authorised", "decide", "decision", "verdict",
	}
	for _, f := range AdditionalWritableFields() {
		for _, w := range grantWords {
			if strings.Contains(strings.ToLower(f), w) {
				t.Errorf("313 field %q carries the approval word %q (AGENTS.md §1.2)", f, w)
			}
		}
		if family, why := refusedLockedFamily(f); family != "" {
			t.Errorf("313 field %q fell into the locked family %s (%s)", f, family, why)
		}
	}
	leg := &cfg248Leg{}
	if _, err := cfg248Dispatch(leg).Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigSet, "requestId": "r-313-lock", "source": ComposerRequestSource,
		"configField": "permission_mode", "configValue": "auto_approve", "configProvider": "stepfun",
	})); !errors.Is(err, ErrConfigFieldRefused) {
		t.Errorf("the locked-family control returned %v, want %v", err, ErrConfigFieldRefused)
	}
	if leg.applyCalls != 0 {
		t.Error("the locked-family control reached the write leg")
	}
}

// The read side: config.get still answers ticket 248's two sentences verbatim and
// now carries ticket 313's machine line after them.
func Test313ConfigGetCarriesThePickerPayloadAfterTheFrozenProse(t *testing.T) {
	view := SettingsView{
		Readable: true, Credential: CredentialNoneDeclared,
		ChatModel: "model-under-test-313", ChatModelKnown: true, ChatProvider: "stepfun",
		Presets: []string{"anthropic", "deepseek", "stepfun"}, PresetsKnown: true,
		Configured: true, ConfiguredKnown: true,
		Providers: []ProviderSettings{{
			Name: "stepfun", APIKeyRef: "dpapi:313canaryblob01", HasBaseURL: true,
			ModelCount: 1, Models: []string{"model-under-test-313"},
		}},
		Voice: VoiceSettingsView{
			CloudASRChain: []string{"stepfun/model-under-test-313"}, CloudTTSChain: []string{"stepfun/model-under-test-313"},
			RealtimeEnabled: true, RealtimeProvider: "stepfun", RealtimeModel: "model-under-test-313",
		},
	}
	leg := &cfg248Leg{readView: view}
	reply, err := cfg248Dispatch(leg).Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigGet, "requestId": "r-313-get", "source": ComposerRequestSource,
	}))
	if err != nil {
		t.Fatalf("config.get: %v", err)
	}
	// Ticket 248's two sentences, verbatim, still the beginning of the answer.
	for _, frozen := range []string{"配置可读。服务商 1 家；凭据状态 no_ref_declared。", "当前聊天模型 model-under-test-313。"} {
		if !strings.Contains(reply, frozen) {
			t.Errorf("the frozen sentence %q left the config.get answer: %q", frozen, reply)
		}
	}
	idx := strings.Index(reply, SettingsPayloadMarker)
	if idx < 0 {
		t.Fatalf("the structured payload is not in the answer: %q", reply)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(strings.SplitN(reply[idx+len(SettingsPayloadMarker):], "\n", 2)[0]), &payload); err != nil {
		t.Fatalf("the payload after the marker is not JSON: %v", err)
	}
	for _, key := range []string{"readable", "credential", "presetsKnown", "presets", "configuredKnown", "configured", "chatProvider", "chatModel", "chatModelKnown", "providers", "voice"} {
		if _, ok := payload[key]; !ok {
			t.Errorf("the payload lost the key %q: %v", key, payload)
		}
	}
	if got, _ := payload["chatProvider"].(string); got != "stepfun" {
		t.Errorf("the payload does not carry the chat provider half: %v", payload["chatProvider"])
	}
	if got, _ := payload["configured"].(bool); !got {
		t.Error("the configured flag did not survive into the payload")
	}
	providers, _ := payload["providers"].([]any)
	if len(providers) != 1 {
		t.Fatalf("payload providers = %v", payload["providers"])
	}
	row, _ := providers[0].(map[string]any)
	models, _ := row["models"].([]any)
	if len(models) != 1 || models[0] != "model-under-test-313" {
		t.Errorf("the provider row's model NAMES are not in the payload: %v", row["models"])
	}
}

// A leg that answered nothing still produces an honest payload: not "no presets",
// but "presetsKnown:false, presets:null". The unreadable state also keeps the
// prose's own wording (no "未录入" appears where the config could not be read).
func Test313SilentAndUnreadableLegsRenderAsStatesNotAsFacts(t *testing.T) {
	silent := &cfg248Leg{readView: SettingsView{Readable: true}}
	reply, err := cfg248Dispatch(silent).Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigGet, "requestId": "r-313-silent", "source": ComposerRequestSource,
	}))
	if err != nil {
		t.Fatalf("config.get: %v", err)
	}
	idx := strings.Index(reply, SettingsPayloadMarker)
	if idx < 0 {
		t.Fatalf("a silent leg produced no payload line: %q", reply)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(strings.SplitN(reply[idx+len(SettingsPayloadMarker):], "\n", 2)[0]), &payload); err != nil {
		t.Fatalf("payload: %v", err)
	}
	if known, _ := payload["presetsKnown"].(bool); known {
		t.Error("a leg that never answered reported presetsKnown=true")
	}
	if payload["presets"] != nil {
		t.Errorf("a leg that never answered produced preset values: %v", payload["presets"])
	}

	unreadable := &cfg248Leg{readView: SettingsView{Readable: false, Credential: CredentialConfigUnreadable}}
	reply, err = cfg248Dispatch(unreadable).Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigGet, "requestId": "r-313-unread", "source": ComposerRequestSource,
	}))
	if err != nil {
		t.Fatalf("unreadable config.get: %v", err)
	}
	if !strings.Contains(reply, "不可读") || strings.Contains(reply, "未录入") {
		t.Errorf("the unreadable answer changed shape: %q", reply)
	}
	if !strings.Contains(reply, SettingsPayloadMarker) {
		t.Errorf("the unreadable answer carries no state payload at all: %q", reply)
	}
}

// The payload's wire shape is pinned key by key, so a later addition to
// SettingsView cannot silently start travelling to the page, and none of the keys
// that do travel is a value-bearing name - the read side keeps ticket 248's
// property (states, counts and references only; the one credential-shaped string
// is the reference the schema itself writes).
func Test313PayloadKeySetIsFrozenAndCarriesNoValueBearingKey(t *testing.T) {
	line := renderSettingsPayload(SettingsView{
		Readable: true,
		Providers: []ProviderSettings{{
			Name: "stepfun", APIKeyRef: "env:DEEPSEEK_KEY", Models: []string{"m"},
		}},
	})
	top := `["chatModel","chatModelKnown","chatProvider","configured","configuredKnown","credential","presets","presetsKnown","providers","readable","voice"]`
	var payload map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimPrefix(line, SettingsPayloadMarker)), &payload); err != nil {
		t.Fatalf("payload: %v", err)
	}
	keys := make([]string, 0, len(payload))
	for k := range payload {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if got := `["` + strings.Join(keys, `","`) + `"]`; got != top {
		t.Errorf("the payload's top-level keys = %s, want %s", got, top)
	}
	// The provider row and the voice half are pinned the same way: their keys are
	// the contract the page codes against.
	var rows []map[string]json.RawMessage
	if err := json.Unmarshal(payload["providers"], &rows); err != nil {
		t.Fatalf("providers: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("providers payload is not a one-element array: %s", payload["providers"])
	}
	rowKeys := make([]string, 0, len(rows[0]))
	for k := range rows[0] {
		rowKeys = append(rowKeys, k)
	}
	sort.Strings(rowKeys)
	wantRow := `["apiKeyRef","baseURL","billing","credentialError","credentialSet","hasBaseURL","modelCount","models","name","protocol"]`
	if got := `["` + strings.Join(rowKeys, `","`) + `"]`; got != wantRow {
		t.Errorf("provider row keys = %s, want %s", got, wantRow)
	}
	var voice map[string]json.RawMessage
	if err := json.Unmarshal(payload["voice"], &voice); err != nil {
		t.Fatalf("voice: %v", err)
	}
	voiceKeys := make([]string, 0, len(voice))
	for k := range voice {
		voiceKeys = append(voiceKeys, k)
	}
	sort.Strings(voiceKeys)
	wantVoice := `["cloudASRChain","cloudTTSChain","realtimeEnabled","realtimeModel","realtimeProvider"]`
	if got := `["` + strings.Join(voiceKeys, `","`) + `"]`; got != wantVoice {
		t.Errorf("voice payload keys = %s, want %s", got, wantVoice)
	}
	// No travelled key is a value-bearing name (exact match: apiKeyRef is the
	// reference the schema writes, and it is a different name from any value key).
	valueBearing := map[string]bool{
		"credentialvalue": true, "apikey": true, "apikeyvalue": true,
		"value": true, "secret": true, "token": true, "key": true,
	}
	for _, k := range append(append(append([]string{}, keys...), rowKeys...), voiceKeys...) {
		if valueBearing[strings.ToLower(k)] {
			t.Errorf("the payload travels a value-bearing key %q; the read side carries states, counts and references only", k)
		}
	}
}

// Test313ZeroRowPayloadKeepsProvidersAnArray is the pin for the one shape the
// acceptance leg 313-v1 measured against the live page: the settings page reads
// `providers` as an array with no null guard (find/some at render), so a fresh
// config with zero provider rows must still put [] on the wire. This is the
// AC#1 target state - the owner's very first open of the model category.
func Test313ZeroRowPayloadKeepsProvidersAnArray(t *testing.T) {
	line := renderSettingsPayload(SettingsView{Readable: true})
	if !strings.Contains(line, `"providers":[]`) {
		t.Fatalf("a zero-row SettingsView puts providers on the wire as something other than []: %s", line)
	}
	// Control: one row still travels as an array of objects, not a bare [].
	one := renderSettingsPayload(SettingsView{
		Readable:  true,
		Providers: []ProviderSettings{{Name: "stepfun", APIKeyRef: "dpapi:x", ModelCount: 1, Models: []string{"m"}}},
	})
	if !strings.Contains(one, `"providers":[{"name":"stepfun"`) {
		t.Fatalf("a one-row SettingsView no longer marshals the row as the page expects: %s", one)
	}
}
