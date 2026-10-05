package config

import (
	"fmt"
	"strings"
	"testing"
)

// validConfig returns a defaults-based config that must pass validation
// (schema_version stamped as a load would).
func validConfig(t *testing.T) *Config {
	t.Helper()
	c := NewDefaults()
	c.SchemaVersion = SchemaVersionCurrent
	if err := validate(c); err != nil {
		t.Fatalf("defaults must validate: %v", err)
	}
	return c
}

func TestValidateDefaultsPass(t *testing.T) {
	validConfig(t)
}

func TestValidateEnumErrors(t *testing.T) {
	cases := []struct {
		name string
		path string // key path that must appear in the error
		mut  func(*Config)
	}{
		{"theme", "app.theme", func(c *Config) { c.App.Theme = "blue" }},
		{"observe.level", "observe.level", func(c *Config) { c.Observe.Level = "verbose" }},
		{"proxy.mode", "net.proxy.mode", func(c *Config) { c.Net.Proxy.Mode = "auto" }},
		{"over_budget", "cost.over_budget", func(c *Config) { c.Cost.OverBudget = "stop" }},
		{"protocol", "llm.providers.acme.protocol", func(c *Config) {
			c.LLM.Providers = map[string]Provider{
				"acme": {Protocol: "grpc", BaseURL: "https://acme.example"},
			}
		}},
		{"provider billing", "llm.providers.acme.billing", func(c *Config) {
			c.LLM.Providers = map[string]Provider{
				"acme": {Protocol: ProtocolOpenAIChat, Billing: "free"},
			}
		}},
		{"model billing", "llm.providers.acme.models.m1.billing", func(c *Config) {
			c.LLM.Providers = map[string]Provider{
				"acme": {Protocol: ProtocolOpenAIChat, Models: map[string]ModelSpec{
					"m1": {Billing: "per-char"},
				}},
			}
		}},
		{"thinking_intensity", "llm.roles.chat.thinking_intensity", func(c *Config) {
			c.LLM.Roles.Chat.ThinkingIntensity = "maximum"
		}},
		{"thinking_levels element", "llm.providers.acme.models.m1.thinking_levels", func(c *Config) {
			c.LLM.Providers = map[string]Provider{
				"acme": {Protocol: ProtocolOpenAIChat, Models: map[string]ModelSpec{
					"m1": {ThinkingLevels: []string{ThinkingLow, "extreme"}},
				}},
			}
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := validConfig(t)
			tc.mut(c)
			err := validate(c)
			if err == nil {
				t.Fatalf("expected error for %s", tc.name)
			}
			if !strings.Contains(err.Error(), tc.path) {
				t.Fatalf("error %q must name key path %q", err, tc.path)
			}
		})
	}
}

func TestValidateEnumAcceptsEveryListedValue(t *testing.T) {
	for _, theme := range themes {
		c := validConfig(t)
		c.App.Theme = theme
		if err := validate(c); err != nil {
			t.Fatalf("theme %q must validate: %v", theme, err)
		}
	}
	for _, lvl := range observeLevels {
		c := validConfig(t)
		c.Observe.Level = lvl
		if err := validate(c); err != nil {
			t.Fatalf("observe.level %q must validate: %v", lvl, err)
		}
	}
	for _, mode := range proxyModes {
		c := validConfig(t)
		c.Net.Proxy.Mode = mode
		if err := validate(c); err != nil {
			t.Fatalf("proxy.mode %q must validate: %v", mode, err)
		}
	}
	for _, mode := range overBudgetModes {
		c := validConfig(t)
		c.Cost.OverBudget = mode
		if err := validate(c); err != nil {
			t.Fatalf("over_budget %q must validate: %v", mode, err)
		}
	}
	for _, proto := range protocols {
		c := validConfig(t)
		c.LLM.Providers = map[string]Provider{"acme": {Protocol: proto}}
		if err := validate(c); err != nil {
			t.Fatalf("protocol %q must validate: %v", proto, err)
		}
	}
	for _, b := range billingModes {
		c := validConfig(t)
		c.LLM.Providers = map[string]Provider{"acme": {Protocol: ProtocolOpenAIChat, Billing: b}}
		if err := validate(c); err != nil {
			t.Fatalf("billing %q must validate: %v", b, err)
		}
	}
	for _, ti := range thinkingIntensity {
		c := validConfig(t)
		c.LLM.Roles.Chat.ThinkingIntensity = ti
		if err := validate(c); err != nil {
			t.Fatalf("thinking_intensity %q must validate: %v", ti, err)
		}
	}
}

func TestValidateHardcodedReadOnly(t *testing.T) {
	cases := []struct {
		name string
		path string
		mut  func(*Config)
	}{
		{"audio.half_duplex=false", "audio.half_duplex", func(c *Config) { c.Audio.HalfDuplex = false }},
		{"privacy.keep_transcript=true", "privacy.keep_transcript", func(c *Config) { c.Privacy.KeepTranscript = true }},
		{"privacy.keep_audio=true", "privacy.keep_audio", func(c *Config) { c.Privacy.KeepAudio = true }},
		{"models.verify_signature=false", "models.verify_signature", func(c *Config) { c.Models.VerifySignature = false }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := validConfig(t)
			tc.mut(c)
			err := validate(c)
			if err == nil {
				t.Fatalf("%s must be rejected", tc.name)
			}
			if !strings.Contains(err.Error(), tc.path) {
				t.Fatalf("error %q must name %q", err, tc.path)
			}
		})
	}
}

func TestValidateBallSizeRange(t *testing.T) {
	c := validConfig(t)
	for _, bad := range []int{43, 73, 0, -1} {
		c.Ball.Size = bad
		if err := validate(c); err == nil {
			t.Fatalf("ball.size=%d must be rejected (allowed 44-72)", bad)
		}
	}
	for _, good := range []int{44, 56, 72} {
		c.Ball.Size = good
		if err := validate(c); err != nil {
			t.Fatalf("ball.size=%d must validate: %v", good, err)
		}
	}
}

// TestValidateRiskConfirmTimeoutRange pins ticket 267's shape (a): [risk].
// confirm_timeout_sec gets a load-time band, so an out-of-range value is
// refused as "this config is invalid" instead of reaching the approval gate,
// where 99999 would hang an L2 card for 27h46m and 10 would take the C18
// pre-timeout warning away without a word.
//
// The refusal must name the key AND the band, otherwise a user reading the
// error cannot tell which line to fix (the same rule validateBall follows).
func TestValidateRiskConfirmTimeoutRange(t *testing.T) {
	c := validConfig(t)
	band := fmt.Sprintf("out of range [%d, %d]", confirmTimeoutSecMin, confirmTimeoutSecMax)
	for _, bad := range []int{
		confirmTimeoutSecMin - 1, // 30: exactly approval.DefaultApprovalWarning -> C18 lead would be <= 0
		10,                       // a timeout below the warning lead
		confirmTimeoutSecMax + 1, // 3601
		99999,                    // ticket 267's 27h46m seed
		0, -1,                    // unset-shaped values must not read as "use the default"
	} {
		c.Risk.ConfirmTimeoutSec = bad
		err := validate(c)
		if err == nil {
			t.Fatalf("risk.confirm_timeout_sec=%d must be rejected (%s)", bad, band)
		}
		if !strings.Contains(err.Error(), "risk.confirm_timeout_sec") {
			t.Fatalf("error %q must name the key risk.confirm_timeout_sec", err)
		}
		if !strings.Contains(err.Error(), band) {
			t.Fatalf("error %q must state the band %q", err, band)
		}
		if !strings.Contains(err.Error(), fmt.Sprintf("%d out of range", bad)) {
			t.Fatalf("error %q must quote the offending value %d", err, bad)
		}
	}
	// 31 is the floor precisely because the C18 warning needs Timeout() strictly
	// greater than WarningLead(); validate_267_test.go re-derives that floor
	// from approval.DefaultApprovalWarning instead of trusting this line.
	for _, good := range []int{confirmTimeoutSecMin, 300, confirmTimeoutSecMax} {
		c.Risk.ConfirmTimeoutSec = good
		if err := validate(c); err != nil {
			t.Fatalf("risk.confirm_timeout_sec=%d must validate: %v", good, err)
		}
	}
}

// TestValidateRiskConfirmTimeoutBandRejectsWarningLeadPlusOne pins the two
// edges of the band itself rather than a sample of it: one below the floor and
// one at it. Kept separate from the table above so that moving the floor (a
// contract decision, ticket 267's residual 1) fails here with the lead in the
// message rather than in a list of numbers.
func TestValidateRiskConfirmTimeoutBandRejectsWarningLeadPlusOne(t *testing.T) {
	c := validConfig(t)
	c.Risk.ConfirmTimeoutSec = confirmTimeoutSecMin
	if err := validate(c); err != nil {
		t.Fatalf("the floor %d must load: %v", confirmTimeoutSecMin, err)
	}
	c.Risk.ConfirmTimeoutSec = confirmTimeoutSecMin - 1
	if err := validate(c); err == nil {
		t.Fatalf("a timeout equal to the C18 warning lead (%ds) must be rejected: it arms no warning at all", confirmTimeoutSecMin-1)
	}
}

func TestValidateAPIKeyRefFormat(t *testing.T) {
	cases := []struct {
		name   string
		ref    string
		wantOK bool
	}{
		{"empty is keyless", "", true},
		{"dpapi ref", "dpapi:abcd1234", true},
		{"env ref", "env:WISP_TEST_LLM_KEY", true},
		{"plaintext key", "sk-abcdef123456", false},
		{"bare blob id", "abcd1234", false},
		{"traversal blob id", "dpapi:..", false},
		{"empty env name", "env:", false},
		{"unknown scheme", "vault:secret/x", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := validConfig(t)
			c.LLM.Providers = map[string]Provider{
				"acme": {Protocol: ProtocolOpenAIChat, APIKeyRef: tc.ref},
			}
			err := validate(c)
			if tc.wantOK && err != nil {
				t.Fatalf("ref %q must validate: %v", tc.ref, err)
			}
			if !tc.wantOK && err == nil {
				t.Fatalf("ref %q must be rejected", tc.ref)
			}
			if !tc.wantOK && !strings.Contains(err.Error(), "llm.providers.acme.api_key_ref") {
				t.Fatalf("error %q must name llm.providers.acme.api_key_ref", err)
			}
		})
	}
}

func TestValidateRealtimeAPIKeyRefFormat(t *testing.T) {
	c := validConfig(t)
	c.Voice.Realtime.APIKeyRef = "sk-plaintext"
	err := validate(c)
	if err == nil {
		t.Fatal("plaintext realtime api_key_ref must be rejected")
	}
	if !strings.Contains(err.Error(), "voice.realtime.api_key_ref") {
		t.Fatalf("error %q must name voice.realtime.api_key_ref", err)
	}
}
