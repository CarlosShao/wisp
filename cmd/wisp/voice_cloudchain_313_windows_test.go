//go:build windows

package main

// Ticket 313's behaviour teeth for the cloud voice chains, driven through the
// production assembly exactly as config_wire_423_windows_test.go drives the
// [session] feed: buildResidentVoice against a config.toml this case wrote,
// then the residency's own loader - the same factory a session scope's
// Acquire walks - and the chain engine it hands back.
//
// THE TWO MUTATIONS THE DISPATCH NAMES, and where each turns red:
//
//   - 改配置 ⇒ 装配跟着变: Test313CloudVoiceChainAssemblesIntoTheResidentLoaders
//     pins the chain head the loader built against the config element verbatim;
//     Test313CloudVoiceChainFollowsTheConfig pins a SECOND build over a
//     DIFFERENT config element and demands the head moved with it.
//   - 删调用点 ⇒ 红: delete the ResolveCloudChain call in buildResidentVoice
//     (or the NewCloud*Chain branch in the loaders) and the two tests above
//     die - the loaders fall back to the local factories, which on this
//     fixture's store refuse ("语音链：ASR 目录为空" / NewStreamingAsr's
//     missing-files error), and the first red sentence names the missing
//     chain family: "拆掉调用点会红：ASR 族里装的不是云端链（got %T）".
//     The degraded-loud case pins the other half: a chain the catalog refuses
//     must never reach the loaders either, and the refusal must be named.
//
// Fixtures: the same voice15e helpers the rest of the voice family uses (the
// staged carrier rosters this file through voice15eVadModel, so a machine
// without third_party skips these cases by name instead of failing on them).

import (
	"errors"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/session"
	"github.com/CarlosShao/wisp/internal/speech"
)

// Made-up catalog ids (A877 ruling 4: no real model names are invented; these
// rows live only in the case's own config.toml).
const (
	model313ChainAsr  = "model-under-test-313-chain-asr"
	model313ChainAsrB = "model-under-test-313-chain-asr-b"
	model313ChainTts  = "model-under-test-313-chain-tts"
)

// seed313ChainConfig lays a config.toml with a stepfun row whose models
// declare the capability bits the two chain directions need, and the two
// chains pointing where asrChainModel says. A model WITHOUT the bit is the
// other leg's refusal shape, exercised separately below.
func seed313ChainConfig(t *testing.T, dataDir, asrChainModel string) {
	t.Helper()
	t.Setenv("WISP_SILERO_MODEL", voice15eVadModel(t))
	voice15eWriteConfig(t, dataDir, func(c *config.Config) {
		if c.LLM.Providers == nil {
			c.LLM.Providers = map[string]config.Provider{}
		}
		c.LLM.Providers["stepfun"] = config.Provider{
			Protocol: config.ProtocolOpenAIChat,
			BaseURL:  "https://api.stepfun.com/v1",
			Models: map[string]config.ModelSpec{
				model313ChainAsr:  {Enabled: true, Capabilities: config.Capabilities{AudioIn: true}},
				model313ChainAsrB: {Enabled: true, Capabilities: config.Capabilities{AudioIn: true}},
				model313ChainTts:  {Enabled: true, Capabilities: config.Capabilities{AudioOut: true}},
			},
		}
		c.Voice.CloudASRChain = []string{"stepfun/" + asrChainModel}
		c.Voice.CloudTTSChain = []string{"stepfun/" + model313ChainTts}
	})
}

// chain313Of extracts the cascade a family loader built. With no local model
// set on the machine the cascade holds no fallback and the loader returns the
// bare chain; with one, it returns the ownership-joining wrapper. Both are
// the wired shape; anything else is the red line the call point's deletion
// lands on.
func chain313Of(asr speech.AsrEngine) (*speech.CloudASRChain, error) {
	switch a := asr.(type) {
	case *voiceCloudAsrFamily:
		return a.CloudASRChain, nil
	case *speech.CloudASRChain:
		return a, nil
	default:
		return nil, errors.New("not a cloud ASR chain family (this error carries no verdict; the caller's Fatalf names the type)")
	}
}

func ttsChain313Of(tts speech.TtsEngine) (*speech.CloudTTSChain, error) {
	switch t := tts.(type) {
	case *voiceCloudTtsFamily:
		return t.CloudTTSChain, nil
	case *speech.CloudTTSChain:
		return t, nil
	default:
		return nil, errors.New("not a cloud TTS chain family (this error carries no verdict; the caller's Fatalf names the type)")
	}
}

// Test313CloudVoiceChainAssemblesIntoTheResidentLoaders is the boot-half and
// loader-half of the wiring in one production walk: the verdict names the
// resolved chains (so the boot line cannot overclaim), and each family's
// loader really builds the cascade whose head IS the config element - the
// strong form of "the settings write reaches the runtime path", because the
// loaders are what a live utterance/announcement acquires.
func Test313CloudVoiceChainAssemblesIntoTheResidentLoaders(t *testing.T) {
	dataDir := voice15eDataDir(t)
	seed313ChainConfig(t, dataDir, model313ChainAsr)
	ra := &residentAudio{frames: make(chan []byte)}
	rv := buildResidentVoice(dataDir, ra, nil, nil, nil)

	if !strings.HasPrefix(rv.verdict, "语音链已装配") {
		t.Fatalf("the voice leg did not assemble with a cloud chain configured: verdict=%q (problems=%v)",
			rv.verdict, rv.problemsSnapshot())
	}
	if !strings.Contains(rv.verdict, "云端听写 stepfun/"+model313ChainAsr) ||
		!strings.Contains(rv.verdict, "云端朗读 stepfun/"+model313ChainTts) {
		t.Fatalf("the boot verdict does not name the two resolved chains: verdict=%q", rv.verdict)
	}

	// The ASR family: the loader builds the cascade and its head is the
	// config's first element. Delete the call point and this line is the red
	// sentence: the family stops being a chain at all.
	h, err := rv.residency.loader(session.EngineASR)
	if err != nil {
		t.Fatalf("拆掉调用点会红（loader 落回本地工厂）：ASR loader refused a configured cloud chain: %v", err)
	}
	handle, ok := h.(*voiceEngineHandle)
	if !ok || handle == nil {
		t.Fatalf("the loader handed back %T, want the leg's own engine handle", h)
	}
	chain, err := chain313Of(handle.Asr())
	if err != nil {
		t.Fatalf("拆掉调用点会红：ASR 族里装的不是云端链（got %T）", handle.Asr())
	}
	if got := chain.Active(); got != "stepfun/"+model313ChainAsr {
		t.Fatalf("the ASR chain head is not the config element: got %q, want stepfun/%s", got, model313ChainAsr)
	}
	if err := handle.Asr().Close(); err != nil {
		t.Fatalf("close the ASR family: %v", err)
	}

	// The TTS twin, same walk.
	th, err := rv.residency.loader(session.EngineTTS)
	if err != nil {
		t.Fatalf("拆掉调用点会红（loader 落回本地工厂）：TTS loader refused a configured cloud chain: %v", err)
	}
	thandle, tok := th.(*voiceEngineHandle)
	if !tok || thandle == nil {
		t.Fatalf("the loader handed back %T, want the leg's own engine handle", th)
	}
	tchain, err := ttsChain313Of(thandle.Tts())
	if err != nil {
		t.Fatalf("拆掉调用点会红：TTS 族里装的不是云端链（got %T）", thandle.Tts())
	}
	if got := tchain.Active(); got != "stepfun/"+model313ChainTts {
		t.Fatalf("the TTS chain head is not the config element: got %q, want stepfun/%s", got, model313ChainTts)
	}
	if err := thandle.Tts().Close(); err != nil {
		t.Fatalf("close the TTS family: %v", err)
	}
}

// Test313CloudVoiceChainFollowsTheConfig is the second half of the behaviour
// tooth: a DIFFERENT config element produces a DIFFERENT head on a fresh
// production build. A wiring that resolved once and cached the result would
// pass the first test and die here.
func Test313CloudVoiceChainFollowsTheConfig(t *testing.T) {
	dataDir := voice15eDataDir(t)
	seed313ChainConfig(t, dataDir, model313ChainAsrB)
	ra := &residentAudio{frames: make(chan []byte)}
	rv := buildResidentVoice(dataDir, ra, nil, nil, nil)
	h, err := rv.residency.loader(session.EngineASR)
	if err != nil {
		t.Fatalf("拆掉调用点会红（loader 落回本地工厂）：ASR loader refused the rewritten chain: %v", err)
	}
	bhandle, ok := h.(*voiceEngineHandle)
	if !ok || bhandle == nil {
		t.Fatalf("the loader handed back %T, want the leg's own engine handle", h)
	}
	chain, err := chain313Of(bhandle.Asr())
	if err != nil {
		t.Fatalf("拆掉调用点会红：改配置后 ASR 族里装的不是云端链（got %T）", bhandle.Asr())
	}
	if got := chain.Active(); got != "stepfun/"+model313ChainAsrB {
		t.Fatalf("改配置后装配没有跟着变：head=%q, want stepfun/%s", got, model313ChainAsrB)
	}
	if err := bhandle.Asr().Close(); err != nil {
		t.Fatalf("close the rewritten ASR family: %v", err)
	}
}

// Test313UnresolvableCloudChainSaysSoAndKeepsTheLocalCascade pins the degraded
// half: a chain element whose model lacks the direction's capability bit must
// NOT assemble (the endpoint would lie about what it can hear/speak), must be
// named in the problems, and the boot must keep the plain local verdict -
// never a silent skip, never a boot failure.
func Test313UnresolvableCloudChainSaysSoAndKeepsTheLocalCascade(t *testing.T) {
	dataDir := voice15eDataDir(t)
	t.Setenv("WISP_SILERO_MODEL", voice15eVadModel(t))
	voice15eWriteConfig(t, dataDir, func(c *config.Config) {
		if c.LLM.Providers == nil {
			c.LLM.Providers = map[string]config.Provider{}
		}
		c.LLM.Providers["stepfun"] = config.Provider{
			Protocol: config.ProtocolOpenAIChat,
			BaseURL:  "https://api.stepfun.com/v1",
			Models: map[string]config.ModelSpec{
				// Both bits false: the model exists, the direction it is
				// being asked to serve does not.
				model313ChainAsr: {Enabled: true},
			},
		}
		c.Voice.CloudASRChain = []string{"stepfun/" + model313ChainAsr}
	})
	ra := &residentAudio{frames: make(chan []byte)}
	rv := buildResidentVoice(dataDir, ra, nil, nil, nil)

	if !strings.HasPrefix(rv.verdict, "语音链已装配") {
		t.Fatalf("a refused cloud chain must not fail the boot: verdict=%q", rv.verdict)
	}
	if strings.Contains(rv.verdict, "云端听写") {
		t.Fatalf("the verdict names a cloud ASR head the loaders will not build: verdict=%q", rv.verdict)
	}
	found := false
	for _, p := range rv.problemsSnapshot() {
		if strings.Contains(p, "云端听写链未装配") && strings.Contains(p, model313ChainAsr) {
			found = true
		}
	}
	if !found {
		t.Fatalf("the refusal is not named with its element: problems=%v", rv.problemsSnapshot())
	}
}
