package config

import (
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Catalog reference validation (SPEC-03 sec 3.1): every chain element, role
// assignment and realtime head must reference an existing
// llm.providers.<provider>.models.<model> pair, and a failing element is
// named in the error. The local sherpa cascade is the terminal voice
// fallback and is never a member of the cloud chains.

// localSherpaProvider is the reserved provider name of the on-device voice
// cascade. It is not a catalog provider: an element naming it in any cloud
// chain is a configuration error (the local cascade is always tried last,
// outside the chains).
const localSherpaProvider = "local-sherpa"

// validateCatalog checks text_chain, the voice cloud chains, the role
// assignments and the realtime head against the provider registry.
func validateCatalog(c *Config) error {
	for _, el := range c.LLM.TextChain {
		if err := checkChainElement("llm.text_chain", el, c.LLM.Providers); err != nil {
			return err
		}
	}
	for _, el := range c.Voice.CloudASRChain {
		if err := checkChainElement("voice.cloud_asr_chain", el, c.LLM.Providers); err != nil {
			return err
		}
	}
	for _, el := range c.Voice.CloudTTSChain {
		if err := checkChainElement("voice.cloud_tts_chain", el, c.LLM.Providers); err != nil {
			return err
		}
	}
	if err := validateRoles(c); err != nil {
		return err
	}
	return validateRealtime(c)
}

// checkChainElement validates one "provider/model" element: split at the
// FIRST '/' (openrouter-style model ids keep their slashes in the catalog
// key), the provider must exist in the registry, the model must exist in
// the provider's catalog, and the reserved local-sherpa name is rejected.
// The error always names the failing element verbatim.
func checkChainElement(chain, el string, providers map[string]Provider) error {
	provider, model, err := splitChainElement(el)
	if err != nil {
		return observe.New(observe.ClassConfig, fmt.Sprintf(
			"config.toml: %s element %q: %v", chain, el, err))
	}
	if provider == localSherpaProvider {
		return observe.New(observe.ClassConfig, fmt.Sprintf(
			"config.toml: %s element %q: the local sherpa cascade is the terminal voice fallback and must not appear in a cloud chain", chain, el))
	}
	p, ok := providers[provider]
	if !ok {
		return observe.New(observe.ClassConfig, fmt.Sprintf(
			"config.toml: %s element %q references unknown provider %q", chain, el, provider))
	}
	if _, ok := p.Models[model]; !ok {
		return observe.New(observe.ClassConfig, fmt.Sprintf(
			"config.toml: %s element %q references unknown model %q of provider %q", chain, el, model, provider))
	}
	return nil
}

// splitChainElement splits "provider/model" at the first '/'.
func splitChainElement(el string) (provider, model string, err error) {
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

// validateRoles checks the four fixed role assignments: provider+model are
// set together (exactly one set is an error) and a set pair must exist in
// the registry.
func validateRoles(c *Config) error {
	r := c.LLM.Roles
	for name, role := range map[string]Role{
		"chat":           r.Chat,
		"memory_extract": r.MemoryExtract,
		"handoff":        r.Handoff,
		"summarize":      r.Summarize,
	} {
		path := "llm.roles." + name
		switch {
		case role.Provider == "" && role.Model == "":
			// unset: the owning feature falls back / disables itself
		case role.Provider == "":
			return observe.New(observe.ClassConfig, fmt.Sprintf(
				"config.toml: %s.model set without %s.provider (set both or neither)", path, path))
		case role.Model == "":
			return observe.New(observe.ClassConfig, fmt.Sprintf(
				"config.toml: %s.provider set without %s.model (set both or neither)", path, path))
		default:
			if _, ok := c.LLM.Providers[role.Provider]; !ok {
				return observe.New(observe.ClassConfig, fmt.Sprintf(
					"config.toml: %s references unknown provider %q", path, role.Provider))
			}
			if !modelExists(c.LLM.Providers, role.Provider, role.Model) {
				return observe.New(observe.ClassConfig, fmt.Sprintf(
					"config.toml: %s references unknown model pair %q", path, role.Provider+"/"+role.Model))
			}
		}
	}
	return nil
}

// validateRealtime checks the realtime voice head; its provider/model pair
// must exist only when realtime is enabled (a disabled head may carry a
// stale pair without breaking load).
func validateRealtime(c *Config) error {
	rt := c.Voice.Realtime
	if !rt.Enabled {
		return nil
	}
	if rt.Provider == "" || rt.Model == "" {
		return observe.New(observe.ClassConfig,
			"config.toml: voice.realtime.enabled=true requires voice.realtime.provider and voice.realtime.model")
	}
	if _, ok := c.LLM.Providers[rt.Provider]; !ok {
		return observe.New(observe.ClassConfig, fmt.Sprintf(
			"config.toml: voice.realtime references unknown provider %q", rt.Provider))
	}
	if !modelExists(c.LLM.Providers, rt.Provider, rt.Model) {
		return observe.New(observe.ClassConfig, fmt.Sprintf(
			"config.toml: voice.realtime references unknown model pair %q", rt.Provider+"/"+rt.Model))
	}
	return nil
}

// modelExists reports whether providers[p] has a catalog entry for model.
func modelExists(providers map[string]Provider, p, model string) bool {
	_, ok := providers[p].Models[model]
	return ok
}
