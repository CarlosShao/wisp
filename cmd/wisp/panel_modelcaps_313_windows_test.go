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
