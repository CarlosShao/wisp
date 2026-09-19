package config

import (
	"strings"
	"testing"
)

// catalogBase returns a valid config with one provider (acme) exposing two
// models, ready for chain/role/realtime reference tests.
func catalogBase(t *testing.T) *Config {
	t.Helper()
	c := validConfig(t)
	c.LLM.Providers = map[string]Provider{
		"acme": {
			Protocol: ProtocolOpenAIChat,
			BaseURL:  "https://acme.example/v1",
			Models: map[string]ModelSpec{
				"m1": {Billing: BillingPayPerToken},
				"m2": {Billing: BillingPayPerToken},
			},
		},
	}
	return c
}

func TestValidateTextChainValid(t *testing.T) {
	c := catalogBase(t)
	c.LLM.TextChain = []string{"acme/m1", "acme/m2"}
	if err := validate(c); err != nil {
		t.Fatalf("valid chain must pass: %v", err)
	}
}

func TestValidateTextChainUnknownElement(t *testing.T) {
	cases := []string{
		"nosuch/m1",    // unknown provider
		"acme/unknown", // unknown model on a known provider
		"noslash",      // malformed element
	}
	for _, el := range cases {
		c := catalogBase(t)
		c.LLM.TextChain = []string{"acme/m1", el}
		err := validate(c)
		if err == nil {
			t.Fatalf("chain element %q must be rejected", el)
		}
		if !strings.Contains(err.Error(), el) {
			t.Fatalf("error %q must name the failing element %q (SPEC-03 sec 3.1)", err, el)
		}
		if !strings.Contains(err.Error(), "text_chain") {
			t.Fatalf("error %q must name text_chain", err)
		}
	}
}

func TestValidateVoiceCloudChains(t *testing.T) {
	t.Run("valid cloud chains", func(t *testing.T) {
		c := catalogBase(t)
		c.Voice.CloudASRChain = []string{"acme/m1"}
		c.Voice.CloudTTSChain = []string{"acme/m2"}
		if err := validate(c); err != nil {
			t.Fatalf("valid cloud chains must pass: %v", err)
		}
	})
	t.Run("unknown element named", func(t *testing.T) {
		c := catalogBase(t)
		c.Voice.CloudASRChain = []string{"acme/ghost"}
		err := validate(c)
		if err == nil || !strings.Contains(err.Error(), "acme/ghost") ||
			!strings.Contains(err.Error(), "cloud_asr_chain") {
			t.Fatalf("cloud_asr_chain must reject unknown element with its name, got: %v", err)
		}
		c = catalogBase(t)
		c.Voice.CloudTTSChain = []string{"ghostcorp/tts-1"}
		err = validate(c)
		if err == nil || !strings.Contains(err.Error(), "ghostcorp/tts-1") ||
			!strings.Contains(err.Error(), "cloud_tts_chain") {
			t.Fatalf("cloud_tts_chain must reject unknown element with its name, got: %v", err)
		}
	})
	t.Run("local-sherpa never in chains", func(t *testing.T) {
		c := catalogBase(t)
		c.Voice.CloudASRChain = []string{"local-sherpa/parakeet"}
		err := validate(c)
		if err == nil || !strings.Contains(err.Error(), "local-sherpa") {
			t.Fatalf("local-sherpa must be rejected from cloud chains (terminal fallback, never a member), got: %v", err)
		}
	})
}

func TestValidateRoleReferences(t *testing.T) {
	t.Run("both empty is unset-ok", func(t *testing.T) {
		c := catalogBase(t)
		if err := validate(c); err != nil {
			t.Fatalf("empty role must pass: %v", err)
		}
	})
	t.Run("provider without model", func(t *testing.T) {
		c := catalogBase(t)
		c.LLM.Roles.Chat.Provider = "acme"
		err := validate(c)
		if err == nil || !strings.Contains(err.Error(), "llm.roles.chat") {
			t.Fatalf("provider-only role must be rejected naming the role, got: %v", err)
		}
	})
	t.Run("model without provider", func(t *testing.T) {
		c := catalogBase(t)
		c.LLM.Roles.Chat.Model = "m1"
		err := validate(c)
		if err == nil || !strings.Contains(err.Error(), "llm.roles.chat") {
			t.Fatalf("model-only role must be rejected naming the role, got: %v", err)
		}
	})
	t.Run("unknown provider named", func(t *testing.T) {
		c := catalogBase(t)
		c.LLM.Roles.Chat.Provider = "ghost"
		c.LLM.Roles.Chat.Model = "m1"
		err := validate(c)
		if err == nil || !strings.Contains(err.Error(), "ghost") {
			t.Fatalf("unknown role provider must be named in the error, got: %v", err)
		}
	})
	t.Run("unknown model named", func(t *testing.T) {
		c := catalogBase(t)
		c.LLM.Roles.Chat.Provider = "acme"
		c.LLM.Roles.Chat.Model = "ghost"
		err := validate(c)
		if err == nil || !strings.Contains(err.Error(), "acme/ghost") {
			t.Fatalf("unknown role model pair must be named in the error, got: %v", err)
		}
	})
	t.Run("valid pair", func(t *testing.T) {
		c := catalogBase(t)
		c.LLM.Roles.Chat.Provider = "acme"
		c.LLM.Roles.Chat.Model = "m1"
		c.LLM.Roles.Chat.ThinkingIntensity = ThinkingMedium
		if err := validate(c); err != nil {
			t.Fatalf("valid role pair must pass: %v", err)
		}
	})
}

func TestValidateRealtimeReferences(t *testing.T) {
	t.Run("disabled with dangling pair tolerated", func(t *testing.T) {
		c := catalogBase(t)
		c.Voice.Realtime = RealtimeConfig{Enabled: false, Provider: "ghost", Model: "m1"}
		if err := validate(c); err != nil {
			t.Fatalf("disabled realtime must not be reference-checked: %v", err)
		}
	})
	t.Run("enabled requires existing pair", func(t *testing.T) {
		c := catalogBase(t)
		c.Voice.Realtime = RealtimeConfig{Enabled: true, Provider: "acme", Model: "ghost"}
		err := validate(c)
		if err == nil || !strings.Contains(err.Error(), "acme/ghost") {
			t.Fatalf("enabled realtime with unknown pair must be rejected naming it, got: %v", err)
		}
	})
	t.Run("enabled without pair rejected", func(t *testing.T) {
		c := catalogBase(t)
		c.Voice.Realtime = RealtimeConfig{Enabled: true}
		err := validate(c)
		if err == nil || !strings.Contains(err.Error(), "voice.realtime") {
			t.Fatalf("enabled realtime without pair must be rejected, got: %v", err)
		}
	})
	t.Run("enabled valid pair", func(t *testing.T) {
		c := catalogBase(t)
		c.Voice.Realtime = RealtimeConfig{Enabled: true, Provider: "acme", Model: "m1"}
		if err := validate(c); err != nil {
			t.Fatalf("enabled realtime with valid pair must pass: %v", err)
		}
	})
}

func TestSplitChainElement(t *testing.T) {
	provider, model, err := splitChainElement("acme/m1")
	if err != nil || provider != "acme" || model != "m1" {
		t.Fatalf("simple split: %q/%q %v", provider, model, err)
	}
	// First '/' splits: openrouter-style model ids keep their slashes.
	provider, model, err = splitChainElement("openrouter/meta-llama/llama-3-70b")
	if err != nil || provider != "openrouter" || model != "meta-llama/llama-3-70b" {
		t.Fatalf("nested model id: %q/%q %v", provider, model, err)
	}
	if _, _, err = splitChainElement("noslash"); err == nil {
		t.Fatal("element without '/' must be rejected")
	}
	if _, _, err = splitChainElement("/m1"); err == nil {
		t.Fatal("element without provider must be rejected")
	}
	if _, _, err = splitChainElement("acme/"); err == nil {
		t.Fatal("element without model must be rejected")
	}
}
