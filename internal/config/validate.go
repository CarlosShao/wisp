package config

import (
	"fmt"
	"slices"
	"strings"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/secret"
)

// Semantic validation (SPEC-03 sec 2 rule 2: nothing invalid passes
// silently). Runs after the strict decode, on the merged defaults+file
// Config. Structural errors (unknown keys, wrong types) are already caught
// by decodeStrict with line numbers; this layer checks values: enums,
// ranges, the hard-coded read-only fields and api_key_ref formats.
//
// Every error is an observe.ClassConfig error naming the offending key path
// ("section.key" or "llm.providers.<name>.key"), never a bare value.

// validate checks the whole config semantically.
func validate(c *Config) error {
	for _, err := range []error{
		validateApp(c),
		validateBall(c),
		validateAudio(c),
		validatePrivacy(c),
		validateModels(c),
		validateCost(c),
		validateObserve(c),
		validateNet(c),
		validateAPIKeyRefs(c),
		validateRolesIntensity(c),
		validateCatalog(c),
	} {
		if err != nil {
			return err
		}
	}
	return nil
}

// validateApp checks the [app] enums (theme is the only one).
func validateApp(c *Config) error {
	if !slices.Contains(themes, c.App.Theme) {
		return observe.New(observe.ClassConfig, fmt.Sprintf(
			"config.toml: app.theme %q must be one of %s", c.App.Theme, vocabList(themes)))
	}
	return nil
}

// validateBall checks the [ball] ranges (SPEC-03 sec 3: size 44-72).
func validateBall(c *Config) error {
	if c.Ball.Size < 44 || c.Ball.Size > 72 {
		return observe.New(observe.ClassConfig, fmt.Sprintf(
			"config.toml: ball.size %d out of range [44, 72]", c.Ball.Size))
	}
	if c.Ball.OpacityIdle < 0 || c.Ball.OpacityIdle > 1 {
		return observe.New(observe.ClassConfig, fmt.Sprintf(
			"config.toml: ball.opacity_idle %v out of range [0, 1]", c.Ball.OpacityIdle))
	}
	return nil
}

// validateAudio enforces the hard-coded read-only audio.half_duplex (D36:
// hard-coded true; writing false is an error, never a setting).
func validateAudio(c *Config) error {
	if !c.Audio.HalfDuplex {
		return observe.New(observe.ClassConfig,
			"config.toml: audio.half_duplex is hard-coded true (read-only); writing false is rejected")
	}
	return nil
}

// validatePrivacy enforces the hard-coded read-only keep_transcript /
// keep_audio (D16: no code path may implement keep_audio=true, SPEC-03 sec 7).
func validatePrivacy(c *Config) error {
	if c.Privacy.KeepTranscript {
		return observe.New(observe.ClassConfig,
			"config.toml: privacy.keep_transcript is hard-coded false (read-only); writing true is rejected")
	}
	if c.Privacy.KeepAudio {
		return observe.New(observe.ClassConfig,
			"config.toml: privacy.keep_audio is hard-coded false (read-only); writing true is rejected")
	}
	return nil
}

// validateModels enforces the hard-coded read-only verify_signature (C29:
// signature checks are not disableable).
func validateModels(c *Config) error {
	if !c.Models.VerifySignature {
		return observe.New(observe.ClassConfig,
			"config.toml: models.verify_signature is hard-coded true (read-only, C29); writing false is rejected")
	}
	return nil
}

// validateCost checks the [cost] enum and threshold range.
func validateCost(c *Config) error {
	if !slices.Contains(overBudgetModes, c.Cost.OverBudget) {
		return observe.New(observe.ClassConfig, fmt.Sprintf(
			"config.toml: cost.over_budget %q must be one of %s", c.Cost.OverBudget, vocabList(overBudgetModes)))
	}
	if c.Cost.AlertThreshold < 0 || c.Cost.AlertThreshold > 1 {
		return observe.New(observe.ClassConfig, fmt.Sprintf(
			"config.toml: cost.alert_threshold %v out of range [0, 1]", c.Cost.AlertThreshold))
	}
	return nil
}

// validateObserve checks the [observe] enum.
func validateObserve(c *Config) error {
	if !slices.Contains(observeLevels, c.Observe.Level) {
		return observe.New(observe.ClassConfig, fmt.Sprintf(
			"config.toml: observe.level %q must be one of %s", c.Observe.Level, vocabList(observeLevels)))
	}
	return nil
}

// validateNet checks the [net] proxy mode enum.
func validateNet(c *Config) error {
	if !slices.Contains(proxyModes, c.Net.Proxy.Mode) {
		return observe.New(observe.ClassConfig, fmt.Sprintf(
			"config.toml: net.proxy.mode %q must be one of %s", c.Net.Proxy.Mode, vocabList(proxyModes)))
	}
	return nil
}

// validateAPIKeyRefs checks every api_key_ref against the frozen reference
// grammar (D36 rule 5: "dpapi:<blob-id>" or "env:NAME"; plaintext is
// structurally impossible in the schema but a pasted plaintext key would
// land in the ref field, so it must be rejected at load). Empty = keyless
// provider, allowed.
func validateAPIKeyRefs(c *Config) error {
	for name, p := range c.LLM.Providers {
		if p.APIKeyRef != "" {
			if _, _, err := secret.ParseRef(p.APIKeyRef); err != nil {
				return observe.New(observe.ClassConfig, fmt.Sprintf(
					"config.toml: llm.providers.%s.api_key_ref: %v", name, err))
			}
		}
		if p.Billing != "" && !slices.Contains(billingModes, p.Billing) {
			return observe.New(observe.ClassConfig, fmt.Sprintf(
				"config.toml: llm.providers.%s.billing %q must be one of %s",
				name, p.Billing, vocabList(billingModes)))
		}
		if p.Protocol != "" && !slices.Contains(protocols, p.Protocol) {
			return observe.New(observe.ClassConfig, fmt.Sprintf(
				"config.toml: llm.providers.%s.protocol %q must be one of %s",
				name, p.Protocol, vocabList(protocols)))
		}
		for id, m := range p.Models {
			if err := validateModelSpec(name, id, m); err != nil {
				return err
			}
		}
	}
	if c.Voice.Realtime.APIKeyRef != "" {
		if _, _, err := secret.ParseRef(c.Voice.Realtime.APIKeyRef); err != nil {
			return observe.New(observe.ClassConfig, fmt.Sprintf(
				"config.toml: voice.realtime.api_key_ref: %v", err))
		}
	}
	return nil
}

// validateModelSpec checks one catalog entry: enums only (model billing and
// thinking_levels vocabulary; SPEC-03 sec 3.1).
func validateModelSpec(provider, id string, m ModelSpec) error {
	path := fmt.Sprintf("llm.providers.%s.models.%s", provider, id)
	if m.Billing != "" && !slices.Contains(billingModes, m.Billing) {
		return observe.New(observe.ClassConfig, fmt.Sprintf(
			"config.toml: %s.billing %q must be one of %s", path, m.Billing, vocabList(billingModes)))
	}
	for _, lvl := range m.ThinkingLevels {
		if !slices.Contains(thinkingIntensity, lvl) {
			return observe.New(observe.ClassConfig, fmt.Sprintf(
				"config.toml: %s.thinking_levels element %q must be one of %s",
				path, lvl, vocabList(thinkingIntensity)))
		}
	}
	return nil
}

// validateRolesIntensity checks the thinking_intensity enum of all four
// fixed roles (provider/model reference checks live in validateCatalog,
// which needs the provider registry to exist first).
func validateRolesIntensity(c *Config) error {
	r := c.LLM.Roles
	for name, role := range map[string]Role{
		"chat":           r.Chat,
		"memory_extract": r.MemoryExtract,
		"handoff":        r.Handoff,
		"summarize":      r.Summarize,
	} {
		if err := validateRoleIntensity("llm.roles."+name, role.ThinkingIntensity); err != nil {
			return err
		}
	}
	return nil
}

// validateRoleIntensity checks one role's thinking_intensity enum.
func validateRoleIntensity(path, intensity string) error {
	if !slices.Contains(thinkingIntensity, intensity) {
		return observe.New(observe.ClassConfig, fmt.Sprintf(
			"config.toml: %s.thinking_intensity %q must be one of %s",
			path, intensity, vocabList(thinkingIntensity)))
	}
	return nil
}

// vocabList renders an enum vocabulary for error messages.
func vocabList(v []string) string {
	return strings.Join(v, "|")
}
