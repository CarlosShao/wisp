package config

// Ticket 313's settings-write additions, judged by the file they produce.
//
// The pattern is settings_248_test.go's: every refusal asserts config.toml byte
// for byte, and every acceptance reads the file back the way the next start
// would. The seed is settingsBase (one preset provider, one catalog entry) so the
// new fields have something real to sit next to, plus one variant with no chat
// role at all, which is the state the pair writer exists for.

import (
	"reflect"
	"strings"
	"testing"
)

// settings313NoPairs is settingsBase with the chat role left unset: the state a
// fresh config is in, and the state in which neither single-key role write can
// reach a legal value (validateRoles refuses both half states).
const settings313NoPairs = `schema_version = 2

[ball]
size = 56

[llm]
text_chain = ["deepseek/deepseek-chat"]

[llm.providers.deepseek]
api_key_ref = "env:DEEPSEEK_KEY"

[llm.providers.deepseek.models.deepseek-chat]
context_window = 64000
`

// The explicit row-creation step (A877 ruling 2): one call, provider section plus
// the model entry it was asked for, both landed and both validated.
func Test313CreateProviderRowCreatesAPresetRowAndItsModel(t *testing.T) {
	m, path, _ := newTestManager(t, settingsBase)

	written, err := m.CreateProviderRow("stepfun", "model-under-test-313")
	if err != nil {
		t.Fatalf("CreateProviderRow: %v", err)
	}
	got := readLoaded(t, path)
	row, ok := got.LLM.Providers["stepfun"]
	if !ok {
		t.Fatalf("the provider row never landed: %v", written)
	}
	spec, ok := row.Models["model-under-test-313"]
	if !ok {
		t.Fatalf("the model row never landed: %v", written)
	}
	if !spec.Enabled {
		t.Error("the created model entry is not enabled: an entry without enabled=true is treated as off and is refused at start-up")
	}
	if row.Billing != "pay-per-token" {
		t.Errorf("the created row's billing = %q, want the schema's own default: a zero value here would override the default on the next load", row.Billing)
	}
	// The key path the file actually gained, named on the receipt: the whole
	// section appeared at once, so diffKeyPaths names the provider, not each leaf
	// under it (a map key present on one side only is the key itself).
	if want := []string{"llm.providers.stepfun"}; !reflect.DeepEqual(written, want) {
		t.Errorf("written key paths = %v, want %v", written, want)
	}
	if body := readSettingsFile(t, path); !strings.Contains(body, "[llm.providers.stepfun.models.model-under-test-313]") {
		t.Errorf("the created model's section is not in the file:\n%s", body)
	}
	if got.Ball.Size != 56 {
		t.Errorf("an unrelated key moved: ball.size = %d", got.Ball.Size)
	}
	if _, ok := got.LLM.Providers["deepseek"]; !ok {
		t.Error("the pre-existing provider row disappeared")
	}
}

// The row-creation step refuses what it cannot stand behind: a non-preset name
// (no protocol field here, so the section would not load back) and a row that
// already exists (this step creates, it does not rewrite). Both leave the file
// byte-identical, and the control proves the same call works one name over.
func Test313CreateProviderRowRefusesNonPresetAndExistingRow(t *testing.T) {
	cases := []struct {
		name     string
		provider string
		mustSay  string
	}{
		{"a non-preset provider name", "inventco", "protocol"},
		{"a row the config already declares", "deepseek", "already exists"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m, path, _ := newTestManager(t, settingsBase)
			before := readSettingsFile(t, path)
			written, err := m.CreateProviderRow(tc.provider, "model-under-test-313")
			if err == nil {
				t.Fatalf("the create was accepted (reported %v)", written)
			}
			if written != nil {
				t.Errorf("a refused create reported %v, want none", written)
			}
			if !strings.Contains(err.Error(), tc.mustSay) {
				t.Errorf("the refusal does not say why (%q missing): %v", tc.mustSay, err)
			}
			if after := readSettingsFile(t, path); after != before {
				t.Errorf("config.toml changed under a refused create:\n%s", after)
			}
		})
	}
	// Control: the same call with a preset name and an empty model still creates
	// the provider row alone (the model half is optional).
	m, path, _ := newTestManager(t, settingsBase)
	if _, err := m.CreateProviderRow("stepfun", ""); err != nil {
		t.Fatalf("the provider-only control was refused: %v", err)
	}
	if _, ok := readLoaded(t, path).LLM.Providers["stepfun"]; !ok {
		t.Error("the control's provider row never landed")
	}
}

// The cloud voice chains (A877 ruling 1): refused while the pair is not in the
// catalog (reason 2, file untouched), accepted once the rows exist, and the local
// sherpa cascade is never a chain member.
func Test313VoiceCloudChainNeedsTheCatalogRowAndRefusesLocalSherpa(t *testing.T) {
	m, path, _ := newTestManager(t, settingsBase)

	before := readSettingsFile(t, path)
	written, err := m.SetVoiceCloudASRChain([]string{"stepfun/model-under-test-313"})
	if err == nil {
		t.Fatalf("a chain element naming a missing provider row was accepted: %v", written)
	}
	if written != nil {
		t.Errorf("a refused chain write reported %v, want none", written)
	}
	if !strings.Contains(err.Error(), "行不存在") {
		t.Errorf("the refusal is not reason 2 (行不存在): %v", err)
	}
	if after := readSettingsFile(t, path); after != before {
		t.Errorf("config.toml changed under a refused chain write:\n%s", after)
	}

	if _, err := m.CreateProviderRow("stepfun", "model-under-test-313"); err != nil {
		t.Fatalf("CreateProviderRow: %v", err)
	}
	written, err = m.SetVoiceCloudASRChain([]string{"stepfun/model-under-test-313"})
	if err != nil {
		t.Fatalf("the chain write after the row existed: %v", err)
	}
	if want := []string{"voice.cloud_asr_chain"}; !reflect.DeepEqual(written, want) {
		t.Errorf("written key paths = %v, want %v", written, want)
	}
	got := readLoaded(t, path)
	if !reflect.DeepEqual(got.Voice.CloudASRChain, []string{"stepfun/model-under-test-313"}) {
		t.Errorf("the chain did not land: %v", got.Voice.CloudASRChain)
	}
	if len(got.Voice.CloudTTSChain) != 0 {
		t.Errorf("the ASR write touched the TTS chain: %v", got.Voice.CloudTTSChain)
	}

	// The local cascade is the terminal fallback, never a chain member - the same
	// ruler validateCatalog applies at load.
	before = readSettingsFile(t, path)
	if _, err := m.SetVoiceCloudTTSChain([]string{"local-sherpa/anything"}); err == nil {
		t.Error("a local-sherpa element was accepted into a cloud chain")
	}
	if after := readSettingsFile(t, path); after != before {
		t.Errorf("a refused local-sherpa write changed the file:\n%s", after)
	}
	// And an empty element is refused rather than skipped.
	if _, err := m.SetVoiceCloudTTSChain([]string{"stepfun/model-under-test-313", " "}); err == nil {
		t.Error("an empty chain element was silently accepted")
	}
	// Clearing the chain is a legitimate write.
	if _, err := m.SetVoiceCloudTTSChain(nil); err != nil {
		t.Fatalf("clearing the TTS chain: %v", err)
	}
}

// The chat role pair: both keys land together, a half never lands, and the pair
// writer exists because the model-half setter cannot reach a legal state on its
// own (the control below is that refusal).
func Test313RoleChatPairWritesBothKeysOrNeither(t *testing.T) {
	m, path, _ := newTestManager(t, settings313NoPairs)

	// The reason the pair writer exists: from an unset role, the model half alone
	// is refused (there is no provider to pair it with).
	if written, err := m.SetRoleChatModel("deepseek-chat"); err == nil {
		t.Fatalf("the model half alone was accepted from an unset role: %v", written)
	}

	before := readSettingsFile(t, path)
	if written, err := m.SetRoleChatPair("stepfun", "deepseek-chat"); err == nil {
		t.Fatalf("a pair naming a provider this config does not declare was accepted: %v", written)
	}
	if after := readSettingsFile(t, path); after != before {
		t.Errorf("a refused pair write changed the file:\n%s", after)
	}
	if written, err := m.SetRoleChatPair("deepseek", ""); err == nil {
		t.Fatalf("a half pair was accepted: %v", written)
	}

	written, err := m.SetRoleChatPair("deepseek", "deepseek-chat")
	if err != nil {
		t.Fatalf("SetRoleChatPair: %v", err)
	}
	want := []string{"llm.roles.chat.model", "llm.roles.chat.provider"}
	if !reflect.DeepEqual(written, want) {
		t.Errorf("written key paths = %v, want %v (both halves named)", written, want)
	}
	got := readLoaded(t, path)
	if got.LLM.Roles.Chat.Provider != "deepseek" || got.LLM.Roles.Chat.Model != "deepseek-chat" {
		t.Errorf("the pair did not land: %+v", got.LLM.Roles.Chat)
	}
	if got.Ball.Size != 56 {
		t.Errorf("an unrelated key moved: ball.size = %d", got.Ball.Size)
	}
}

// The realtime head's two leaves (A877 ruling 1's β shape): a half is legal on a
// disabled head (validateRealtime's own rule, which is the only ruler these two
// writes run against), and on an enabled head a pair that names a row this config
// does not list is refused with the file untouched.
func Test313RealtimeHalfPairFollowsTheEnabledRule(t *testing.T) {
	m, path, _ := newTestManager(t, settingsBase)

	if written, err := m.SetVoiceRealtimeProvider("stepfun"); err != nil {
		t.Fatalf("a provider-only realtime write on a disabled head was refused: %v (%v)", err, written)
	}
	if got := readLoaded(t, path).Voice.Realtime.Provider; got != "stepfun" {
		t.Errorf("the realtime provider did not land: %q", got)
	}
	if written, err := m.SetVoiceRealtimeModel("model-under-test-313"); err != nil {
		t.Fatalf("the model half on a disabled head was refused: %v (%v)", err, written)
	} else if want := []string{"voice.realtime.model"}; !reflect.DeepEqual(written, want) {
		t.Errorf("written key paths = %v, want %v", written, want)
	}
	if got := readLoaded(t, path).Voice.Realtime; !got.Enabled && (got.Provider != "stepfun" || got.Model != "model-under-test-313") {
		t.Errorf("the disabled head's stale pair did not survive reload: %+v", got)
	}

	// The enabled head: the pair must exist in the catalog, so a ghost model or a
	// provider the config does not declare is refused and the file is untouched.
	// The catalog carries a second model so the legal-write control is a real
	// change and not a rewrite of the value the file already holds.
	enabled := settingsBase + `
[llm.providers.deepseek.models.deepseek-reasoner]
context_window = 32000

[voice.realtime]
enabled = true
provider = "deepseek"
model = "deepseek-chat"
`
	m2, path2, _ := newTestManager(t, enabled)
	before := readSettingsFile(t, path2)
	if written, err := m2.SetVoiceRealtimeModel("ghost-model"); err == nil {
		t.Fatalf("an enabled head accepted a pair naming a model its catalog does not list: %v", written)
	}
	if written, err := m2.SetVoiceRealtimeProvider("stepfun"); err == nil {
		t.Fatalf("an enabled head accepted a provider the config does not declare: %v", written)
	}
	if after := readSettingsFile(t, path2); after != before {
		t.Errorf("a refused enabled-head write changed the file:\n%s", after)
	}
	written, err := m2.SetVoiceRealtimeModel("deepseek-reasoner")
	if err != nil {
		t.Fatalf("the enabled head's legal model write: %v", err)
	}
	if want := []string{"voice.realtime.model"}; !reflect.DeepEqual(written, want) {
		t.Errorf("written key paths = %v, want %v", written, want)
	}
	if !strings.Contains(readSettingsFile(t, path2), "model = 'deepseek-reasoner'") {
		t.Errorf("the legal write did not land:\n%s", readSettingsFile(t, path2))
	}
}
