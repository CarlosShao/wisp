package speech

import (
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
)

// This file is the catalog adapter of the ticket 61 cloud cascade (the
// consumption side of ticket 05's catalog schema; SPEC-03 sec 3.1 fixes the
// reference rules, and config/catalog.go already validated them at load).
//
// ResolveCloudChain turns one voice chain from a loaded snapshot into
// runnable endpoints. It is the speech twin of llm.NewResolver +
// ResolveChain, deliberate about the same three things:
//
//   - a chain element names an existing catalog pair (re-checked, not
//     trusted: the caller may hand us a hand-built chain, and a ghost model
//     must fail here, not mid-utterance);
//   - the model must DECLARE the capability bit this direction needs --
//     audio_in for the ASR chain, audio_out for the TTS chain. The declared
//     bit is what the dropdown and this selection read; the probe (ticket
//     11's provider_health) is what verifies it against the wire, and
//     cloudprobe.go is the request definition that verification needs;
//   - an api_key_ref resolves through the injected KeyResolver (secret.Store
//     satisfies it), and the plaintext exists only inside the endpoint --
//     never in a log line, never in an error detail (same rule the C5
//     resolver runs).
//
// A model disabled with enabled=false is refused by name (the ModelSpec
// promise llm/resolver.go spells out applies to the voice chains too: an
// entry the user turned off can never be selected through a chain).

// ResolveCloudChain resolves one voice cloud chain (voice.cloud_asr_chain
// for kind ChainASR, voice.cloud_tts_chain for kind ChainTTS) into ordered,
// fully resolved endpoints. An empty chain resolves to an empty list -- that
// is the "local cascade only" configuration, valid forever.
func ResolveCloudChain(cfg *config.Config, kind CloudChainKind, keys KeyResolver) ([]CloudEndpoint, error) {
	if cfg == nil {
		return nil, observe.New(observe.ClassConfig, "speech: ResolveCloudChain needs a config snapshot")
	}
	var chain []string
	switch kind {
	case ChainASR:
		chain = cfg.Voice.CloudASRChain
	case ChainTTS:
		chain = cfg.Voice.CloudTTSChain
	default:
		return nil, observe.New(observe.ClassConfig, fmt.Sprintf(
			"speech: unknown chain kind %q (want ChainASR or ChainTTS)", string(kind)))
	}
	out := make([]CloudEndpoint, 0, len(chain))
	for _, el := range chain {
		ep, err := resolveChainElement(el, kind, cfg, keys)
		if err != nil {
			return nil, err
		}
		out = append(out, ep)
	}
	return out, nil
}

// resolveChainElement materializes one "provider/model" element. The error
// always names the failing element verbatim -- the same shape
// config/catalog.go's checkChainElement uses at load, so a user matching a
// load error to a settings error sees one vocabulary (ticket 196's lesson
// again: two spellings for one refusal is how a bug hides).
func resolveChainElement(el string, kind CloudChainKind, cfg *config.Config, keys KeyResolver) (CloudEndpoint, error) {
	chainKey := "voice.cloud_" + string(kind) + "_chain"
	provider, model, err := splitElement(el)
	if err != nil {
		return CloudEndpoint{}, observe.New(observe.ClassConfig, fmt.Sprintf(
			"speech: %s element %q: %v", chainKey, el, err))
	}
	p, ok := cfg.LLM.Providers[provider]
	if !ok {
		return CloudEndpoint{}, observe.New(observe.ClassConfig, fmt.Sprintf(
			"speech: %s element %q references unknown provider %q", chainKey, el, provider))
	}
	spec, ok := p.Models[model]
	if !ok {
		return CloudEndpoint{}, observe.New(observe.ClassConfig, fmt.Sprintf(
			"speech: %s element %q references unknown model %q of provider %q", chainKey, el, model, provider))
	}
	if !spec.Enabled {
		return CloudEndpoint{}, observe.New(observe.ClassConfig, fmt.Sprintf(
			"speech: %s element %q names a model disabled in the catalog (enabled=false); re-enable it or remove the element", chainKey, el))
	}
	// The capability bit drives selection (ticket 61: "由 catalog 里
	// capabilities.audio_in/audio_out 的已配置模型驱动"). A model without
	// the bit is a configuration error named element-first, never a silent
	// skip -- a silently dropped chain entry would make the dropdown lie.
	if kind == ChainASR && !spec.Capabilities.AudioIn {
		return CloudEndpoint{}, observe.New(observe.ClassConfig, fmt.Sprintf(
			"speech: %s element %q lacks capabilities.audio_in in its catalog entry", chainKey, el))
	}
	if kind == ChainTTS && !spec.Capabilities.AudioOut {
		return CloudEndpoint{}, observe.New(observe.ClassConfig, fmt.Sprintf(
			"speech: %s element %q lacks capabilities.audio_out in its catalog entry", chainKey, el))
	}

	ep := CloudEndpoint{
		Provider:         provider,
		Model:            model,
		BaseURL:          p.BaseURL,
		ExtraHeaders:     p.Compat.ExtraHeaders,
		Billing:          spec.Billing,
		PriceAudioIn:     spec.Price.AudioIn,
		PriceAudioOut:    spec.Price.AudioOut,
		DeclaresAudioIn:  spec.Capabilities.AudioIn,
		DeclaresAudioOut: spec.Capabilities.AudioOut,
	}
	// TTS-specific knobs come from [voice] tts; the dropdown writes the
	// voice id here and the hot-tier speed rides every Speak.
	if kind == ChainTTS {
		ep.Voice = cfg.Voice.TTS.Voice
		ep.Speed = cfg.Voice.TTS.Speed
	}
	if ref := p.APIKeyRef; ref != "" {
		if keys == nil {
			return CloudEndpoint{}, observe.New(observe.ClassAuth, fmt.Sprintf(
				"speech: provider %q has api_key_ref %q but no key resolver is wired (Unconfigured)", provider, ref))
		}
		key, err := keys.Resolve(ref)
		if err != nil || key == "" {
			return CloudEndpoint{}, observe.New(observe.ClassAuth, fmt.Sprintf(
				"speech: resolving api_key_ref for provider %q failed (Unconfigured)", provider))
		}
		ep.APIKey = key
	}
	return ep, nil
}

// splitElement splits "provider/model" at the FIRST '/' -- openrouter-style
// model ids keep their slashes in the model part, exactly as
// config/catalog.go's splitChainElement parses them at load. Keeping both
// spellings in one file vocabulary (render: ChainKey, parse: splitElement)
// is what makes the round trip testable.
func splitElement(el string) (provider, model string, err error) {
	i := strings.Index(el, "/")
	if i < 0 {
		return "", "", fmt.Errorf(`must be "provider/model"`)
	}
	provider, model = el[:i], el[i+1:]
	if provider == "" || model == "" {
		return "", "", fmt.Errorf(`must be "provider/model"`)
	}
	return provider, model, nil
}
