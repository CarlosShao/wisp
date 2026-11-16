//go:build windows

package main

// Audit R1/P3's behaviour tooth (2026-11-16): the model capability bits are now
// REACHABLE (read face + write face + page switches), and this file pins the
// behaviour the whole fix exists for - 改位 ⇒ 装配结果跟着变:
//
//	default false ⇒ speech.ResolveCloudChain refuses with the named red sentence
//	("lacks capabilities.audio_in" / "audio_out", cloudcatalog.go:97/:101);
//	page write model_audio_in=true through the REAL settings leg (the ticket-248
//	config.set setter family, no new route) + reload ⇒ the SAME chain assembles.
//
// Nothing here re-types a receipt: the refusal is the production assembly's own
// answer, the flip goes through configStore.ApplySetting (the exact path the
// WebView2 host calls), and the read face is read back through the same
// configStore the page reads. That is the three-way agreement the audit asked
// for: 写面落盘、读面可判、装配面跟随。
//
// Also refused here, each by its own sentence: a capability write against a
// model the catalog does not list (reason 2), and a non-bool word (the leg's
// shape check) - both leave config.toml byte-identical.

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/speech"
)

// model313Caps is a made-up catalog id (A877 ruling 4), living only in this
// case's own config.toml. Both chains point at it: the ASR side needs
// capabilities.audio_in, the TTS side capabilities.audio_out, and the schema's
// probe-first default leaves BOTH false - which is exactly the state a fresh
// catalog row lands in (the 反向假话 the audit caught: a savable chain whose
// boot could never honour it).
const model313Caps = "model-under-test-313-caps"

// caps313Seed: one keyless preset provider (no api_key_ref, so the assembly
// needs no key resolver and this file needs no store), one catalog row with
// every capability bit at its default, and both voice chains naming it.
const caps313Seed = `[llm.providers.stepfun]
protocol = "openai-chat"
base_url = "https://api.stepfun.com/v1"

[llm.providers.stepfun.models.` + model313Caps + `]
enabled = true

[voice]
cloud_asr_chain = ['stepfun/` + model313Caps + `']
cloud_tts_chain = ['stepfun/` + model313Caps + `']
`

func caps313Store(t *testing.T) (*configStore, *config.Manager, string) {
	t.Helper()
	dir := sealableTempDir124(t)
	path := filepath.Join(dir, configFileName)
	if err := os.WriteFile(path, []byte(caps313Seed), 0o600); err != nil {
		t.Fatalf("seed config.toml: %v", err)
	}
	mgr, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("manager: %v", err)
	}
	return newConfigStore(mgr, nil, nil, "test-313-caps"), mgr, path
}

func Test313CapabilityBitFlipsAssemblyOutcome(t *testing.T) {
	s, mgr, _ := caps313Store(t)
	ctx := context.Background()

	// HALF 1 - the default-false refusal, verbatim from the production
	// assembly: the chain is savable (the load validates existence only), but
	// no direction can assemble. This is the red sentence the audit named.
	c := mgr.Config()
	if c == nil {
		t.Fatal("the manager never loaded the seed config")
	}
	_, asrErr := speech.ResolveCloudChain(c, speech.ChainASR, nil)
	if asrErr == nil || !strings.Contains(asrErr.Error(), "lacks capabilities.audio_in") {
		t.Fatalf("default-false ASR chain must refuse with the named sentence, got %v", asrErr)
	}
	_, ttsErr := speech.ResolveCloudChain(c, speech.ChainTTS, nil)
	if ttsErr == nil || !strings.Contains(ttsErr.Error(), "lacks capabilities.audio_out") {
		t.Fatalf("default-false TTS chain must refuse with the named sentence, got %v", ttsErr)
	}

	// HALF 1b - the read face answers the default BEFORE the flip: both bits
	// ride as plain false, a real answer the page renders, not an absent key.
	view, err := s.ReadSettings(ctx)
	if err != nil {
		t.Fatalf("ReadSettings: %v", err)
	}
	var row *panel.ProviderSettings
	for i := range view.Providers {
		if view.Providers[i].Name == "stepfun" {
			row = &view.Providers[i]
		}
	}
	if row == nil || row.ModelCapabilities == nil {
		t.Fatalf("the read face did not answer the capability map: %+v", view.Providers)
	}
	bits, ok := row.ModelCapabilities[model313Caps]
	if !ok {
		t.Fatalf("the catalog model %q has no capability entry: %v", model313Caps, row.ModelCapabilities)
	}
	if bits.AudioIn || bits.AudioOut {
		t.Fatalf("a probe-first row answered audio_in=%v audio_out=%v, want false/false", bits.AudioIn, bits.AudioOut)
	}

	// HALF 2 - the page's write: flip audio_in through the ticket-248 setter
	// family (model_audio_in, no new route), then the SAME chain assembles.
	res, err := s.ApplySetting(ctx, panel.SettingWrite{
		Field: panel.FieldModelAudioIn, Provider: "stepfun", Model: model313Caps,
		Value: "true", Actor: "test-313-caps",
	})
	if err != nil {
		t.Fatalf("ApplySetting(model_audio_in=true): %v", err)
	}
	if len(res.Written) != 1 ||
		res.Written[0] != "llm.providers.stepfun.models."+model313Caps+".capabilities.audio_in" {
		t.Fatalf("receipt written = %v, want the capability leaf's dotted key", res.Written)
	}
	if res.Tier != panel.EffectiveRestart {
		t.Fatalf("receipt tier = %v, want restart (the bits feed the boot assembly)", res.Tier)
	}
	if _, err := mgr.CheckAndReload(); err != nil {
		t.Fatalf("CheckAndReload after the flip: %v", err)
	}
	c2 := mgr.Config()
	asrEps, asrErr := speech.ResolveCloudChain(c2, speech.ChainASR, nil)
	if asrErr != nil {
		t.Fatalf("打开 audio_in 后装配仍然拒绝：改位必须让装配结果跟着变: %v", asrErr)
	}
	if len(asrEps) != 1 || asrEps[0].Provider != "stepfun" || asrEps[0].Model != model313Caps {
		t.Fatalf("assembled ASR endpoints = %v, want the chain's own element", asrEps)
	}
	// The other direction is untouched: flipping audio_in must not smuggle the
	// TTS chain past its own bit.
	if _, ttsErr := speech.ResolveCloudChain(c2, speech.ChainTTS, nil); ttsErr == nil ||
		!strings.Contains(ttsErr.Error(), "lacks capabilities.audio_out") {
		t.Fatalf("audio_out must still refuse after the audio_in flip, got %v", ttsErr)
	}
	// The read face answers the NEW truth (the page's switch has a 当前 to show).
	view2, err := s.ReadSettings(ctx)
	if err != nil {
		t.Fatalf("ReadSettings after the flip: %v", err)
	}
	var row2 *panel.ProviderSettings
	for i := range view2.Providers {
		if view2.Providers[i].Name == "stepfun" {
			row2 = &view2.Providers[i]
		}
	}
	if row2 == nil || !row2.ModelCapabilities[model313Caps].AudioIn ||
		row2.ModelCapabilities[model313Caps].AudioOut {
		t.Fatalf("the read face did not follow the flip: %+v", row2)
	}

	// HALF 3 - the TTS twin, through the same leg: audio_out flips, the TTS
	// chain assembles, the receipt names its own leaf.
	res, err = s.ApplySetting(ctx, panel.SettingWrite{
		Field: panel.FieldModelAudioOut, Provider: "stepfun", Model: model313Caps,
		Value: "true", Actor: "test-313-caps",
	})
	if err != nil {
		t.Fatalf("ApplySetting(model_audio_out=true): %v", err)
	}
	if len(res.Written) != 1 ||
		res.Written[0] != "llm.providers.stepfun.models."+model313Caps+".capabilities.audio_out" {
		t.Fatalf("receipt written = %v, want the audio_out leaf's dotted key", res.Written)
	}
	if _, err := mgr.CheckAndReload(); err != nil {
		t.Fatalf("CheckAndReload after the TTS flip: %v", err)
	}
	ttsEps, ttsErr := speech.ResolveCloudChain(mgr.Config(), speech.ChainTTS, nil)
	if ttsErr != nil {
		t.Fatalf("打开 audio_out 后装配仍然拒绝：改位必须让装配结果跟着变: %v", ttsErr)
	}
	if len(ttsEps) != 1 || ttsEps[0].Model != model313Caps {
		t.Fatalf("assembled TTS endpoints = %v, want the chain's own element", ttsEps)
	}
}

// Test313CapabilityWritesRefuseByNameAndMoveNoByte: the two refusal shapes the
// page can produce - a model the catalog does not list (reason 2) and a
// non-bool word (the leg's shape check) - both leave config.toml byte-identical,
// and neither is dressed up as a save.
func Test313CapabilityWritesRefuseByNameAndMoveNoByte(t *testing.T) {
	s, _, path := caps313Store(t)
	ctx := context.Background()

	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.ApplySetting(ctx, panel.SettingWrite{
		Field: panel.FieldModelAudioIn, Provider: "stepfun", Model: "no-such-model-313",
		Value: "true", Actor: "test-313-caps",
	}); err == nil || !strings.Contains(err.Error(), "行不存在") {
		t.Fatalf("a capability write against an unknown model returned %v, want the reason-2 refusal", err)
	}
	if _, err := s.ApplySetting(ctx, panel.SettingWrite{
		Field: panel.FieldModelAudioOut, Provider: "stepfun", Model: model313Caps,
		Value: "maybe", Actor: "test-313-caps",
	}); err == nil || !strings.Contains(err.Error(), "true/false") {
		t.Fatalf("a non-bool word returned %v, want the leg's shape refusal", err)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("a refused capability write moved bytes in config.toml")
	}
}

// ---------------------------------------------------------------------------
// R1-recheck residual 2 (2026-11-16): the boot's other two thresholds had no
// page-visible read face. The capability gate rode ModelCapabilities; the
// enabled gate rode NOTHING (ProviderSettings carried the model NAMES only),
// and the key gate's states (CredentialSet/CredentialError) were answered but
// never tied to the boot's refusal. These teeth pin the read faces on the
// real settings leg, against the boot's own refusal sentences.
// ---------------------------------------------------------------------------

const model313Disabled = "model-under-test-313-off"

// enabled313Seed: one catalog row enabled, one disabled (enabled=false - only
// reachable through a hand-edited file today: schema.go:393 defaults true and
// no write face flips it), and the voice chain naming the DISABLED model.
const enabled313Seed = `[llm.providers.stepfun]
protocol = "openai-chat"
base_url = "https://api.stepfun.com/v1"

[llm.providers.stepfun.models.` + model313Caps + `]
enabled = true

[llm.providers.stepfun.models.` + model313Disabled + `]
enabled = false

[voice]
cloud_asr_chain = ['stepfun/` + model313Disabled + `']
`

func enabled313Store(t *testing.T) (*configStore, *config.Manager) {
	t.Helper()
	dir := sealableTempDir124(t)
	path := filepath.Join(dir, configFileName)
	if err := os.WriteFile(path, []byte(enabled313Seed), 0o600); err != nil {
		t.Fatalf("seed config.toml: %v", err)
	}
	mgr, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("manager: %v", err)
	}
	return newConfigStore(mgr, nil, nil, "test-313-enabled"), mgr
}

// Test313ModelsEnabledReadFaceAnswersTheDisabledGate: the boot refuses the
// chain element naming an enabled=false model with its own named sentence
// (cloudcatalog.go:87-90), and the read face answers the same truth per model
// (the new ProviderSettings.ModelsEnabled, filled in ReadSettings) - so the
// page can judge the gate off the same read the boot runs instead of saying
// 已接入 over a chain the boot will refuse.
func Test313ModelsEnabledReadFaceAnswersTheDisabledGate(t *testing.T) {
	s, mgr := enabled313Store(t)
	ctx := context.Background()

	// Boot half: the production assembly's own refusal, verbatim.
	_, err := speech.ResolveCloudChain(mgr.Config(), speech.ChainASR, nil)
	if err == nil || !strings.Contains(err.Error(), "names a model disabled in the catalog (enabled=false)") {
		t.Fatalf("the disabled-model chain must refuse with the boot's named sentence, got %v", err)
	}

	// Read-face half: the leg answers enabled per model - false for the
	// disabled one, true for the enabled one, both real answers.
	view, err := s.ReadSettings(ctx)
	if err != nil {
		t.Fatalf("ReadSettings: %v", err)
	}
	var row *panel.ProviderSettings
	for i := range view.Providers {
		if view.Providers[i].Name == "stepfun" {
			row = &view.Providers[i]
		}
	}
	if row == nil || row.ModelsEnabled == nil {
		t.Fatalf("the read face did not answer the enabled map: %+v", view.Providers)
	}
	got, ok := row.ModelsEnabled[model313Disabled]
	if !ok {
		t.Fatalf("the disabled model has no enabled entry: %v", row.ModelsEnabled)
	}
	if got {
		t.Errorf("ModelsEnabled[%q] = true, want the catalog's enabled=false", model313Disabled)
	}
	if !row.ModelsEnabled[model313Caps] {
		t.Errorf("ModelsEnabled[%q] = false, want the catalog's enabled=true", model313Caps)
	}
}

// keyref313Seed: the provider DECLARES api_key_ref, and this leg wires no
// credential store - the read face must answer the "cannot affirm a
// resolvable credential" state instead of a fake 已就绪. The model declares
// capabilities.audio_in on purpose: the boot checks the capability gate
// BEFORE the key gate (cloudcatalog.go:95 before :121), and a bit-less model
// would refuse there first, never reaching the key sentence this tooth pins.
const keyref313Seed = `[llm.providers.stepfun]
protocol = "openai-chat"
base_url = "https://api.stepfun.com/v1"
api_key_ref = "dpapi:313canaryblob01"

[llm.providers.stepfun.models.` + model313Caps + `]
enabled = true
capabilities = { audio_in = true, audio_out = false }

[voice]
cloud_asr_chain = ['stepfun/` + model313Caps + `']
`

func keyref313Store(t *testing.T) (*configStore, *config.Manager) {
	t.Helper()
	dir := sealableTempDir124(t)
	path := filepath.Join(dir, configFileName)
	if err := os.WriteFile(path, []byte(keyref313Seed), 0o600); err != nil {
		t.Fatalf("seed config.toml: %v", err)
	}
	mgr, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("manager: %v", err)
	}
	return newConfigStore(mgr, nil, nil, "test-313-keyref"), mgr
}

// Test313CredentialReadFaceNamesTheUnresolvableKeyRef: the key gate's read
// surface (already riding the row: APIKeyRef / CredentialSet / CredentialError)
// answers "cannot affirm" when the ref cannot be checked - the state the
// voice card now renders - and the boot side refuses a declared ref it cannot
// resolve (cloudcatalog.go:121-131, the PROBE-KEYREF sentence family). The
// page's gate reads exactly these two fields; nothing is invented client-side.
func Test313CredentialReadFaceNamesTheUnresolvableKeyRef(t *testing.T) {
	s, mgr := keyref313Store(t)
	ctx := context.Background()

	// Boot half: a declared api_key_ref with no resolvable credential is a
	// named auth refusal, never a silent pass.
	_, err := speech.ResolveCloudChain(mgr.Config(), speech.ChainASR, nil)
	if err == nil || !strings.Contains(err.Error(), "api_key_ref") {
		t.Fatalf("a chain over an unresolved api_key_ref must refuse, got %v", err)
	}

	// Read-face half: CredentialSet=false (not affirmed present) and
	// CredentialError names why - the two states the card renders verbatim.
	view, err := s.ReadSettings(ctx)
	if err != nil {
		t.Fatalf("ReadSettings: %v", err)
	}
	var row *panel.ProviderSettings
	for i := range view.Providers {
		if view.Providers[i].Name == "stepfun" {
			row = &view.Providers[i]
		}
	}
	if row == nil {
		t.Fatalf("the provider row is missing: %+v", view.Providers)
	}
	if row.CredentialSet {
		t.Errorf("CredentialSet = true, want false (the ref could not be affirmed)")
	}
	if row.CredentialError == "" {
		t.Errorf("CredentialError is empty, want the cannot-check state named (store not wired)")
	}
}
