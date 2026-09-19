package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// walkFields visits every struct field reachable from t (through structs,
// maps, slices, pointers) exactly once.
func walkFields(t reflect.Type, visit func(name, tag string)) {
	seen := map[reflect.Type]bool{}
	var rec func(t reflect.Type)
	rec = func(t reflect.Type) {
		if t == nil || seen[t] {
			return
		}
		seen[t] = true
		switch t.Kind() {
		case reflect.Struct:
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				visit(f.Name, f.Tag.Get("toml"))
				rec(f.Type)
			}
		case reflect.Map, reflect.Slice, reflect.Array, reflect.Ptr, reflect.Interface:
			rec(t.Elem())
		}
	}
	rec(t)
}

// TestBoundaryNoPlaintextSecretFields pins D36 rule 5 structurally: the only
// key-bearing form is APIKeyRef ("dpapi:<id>"/"env:NAME"); a plaintext
// api_key field can never be added without this test failing.
func TestBoundaryNoPlaintextSecretFields(t *testing.T) {
	ct := reflect.TypeOf(Config{})
	walkFields(ct, func(name, tag string) {
		if name == "APIKeyRef" {
			return
		}
		if strings.Contains(strings.ToLower(name), "apikey") || strings.Contains(strings.ToLower(name), "api_key") {
			t.Errorf("field %q violates D36 rule 5 (only APIKeyRef may carry key material)", name)
		}
		if tag == "api_key" {
			t.Errorf("a field carries toml tag %q - plaintext api_key is banned", tag)
		}
	})
}

// TestBoundaryNoRuntimeObservationFields pins the SPEC-03 sec 3.1 storage
// split: probe results / health / latency sampling live in SQLite
// provider_health (SPEC-02 schema v2, ticket 09) and must have no
// representation in the config structs (config = declared catalog; SQLite =
// measured runtime state). The provider_health DDL itself lands with
// ticket 09 - deliberately absent from memory schema v1.
func TestBoundaryNoRuntimeObservationFields(t *testing.T) {
	banned := []string{"health", "probe", "latency", "lasterror", "last_error", "successrate", "success_rate", "verified"}
	ct := reflect.TypeOf(Config{})
	walkFields(ct, func(name, tag string) {
		lower := strings.ToLower(name)
		for _, b := range banned {
			if strings.Contains(lower, b) {
				t.Errorf("field %q violates the storage split (runtime observation state belongs to SQLite provider_health)", name)
			}
		}
	})
}

// fullConfig builds a Config with every section carrying non-default values
// (the hard-coded read-onlys stay at their legal values).
func fullConfig(t *testing.T) *Config {
	t.Helper()
	c := NewDefaults()
	c.SchemaVersion = SchemaVersionCurrent
	c.App.Theme = ThemeLight
	c.App.Language = "en-US"
	c.App.Autostart = true
	c.App.SingleInstance = false
	c.Ball.Size = 60
	c.Ball.Position = BallPos{X: 10, Y: 20, Monitor: "DISPLAY1"}
	c.Ball.OpacityIdle = 0.5
	c.Ball.ClickThrough = false
	c.Ball.HideOnFullscreen = false
	c.Hotkey = HotkeySection{Summon: "Ctrl+Space", Mute: "Ctrl+M", Cancel: "F2", Panel: "Ctrl+P"}
	c.Session.WarmTimeoutSec = 100
	c.Session.SettlingSec = 4
	c.Session.ConversationIdleSec = 45
	c.Voice.Enabled = false
	c.Voice.WakeWord = WakeWord{Enabled: true, Keywords: []string{"hey wisp"}, Thresholds: []float64{0.8}, VetoWords: []string{"stop"}}
	c.Voice.ASR = ASRConfig{Provider: "acme", Model: "m1"}
	c.Voice.TTS = TTSConfig{Provider: "acme", Voice: "v1", Speed: 1.25}
	c.Voice.Punctuation = false
	c.Voice.ConversationMode = true
	c.Voice.AEC = AECConfig{Enabled: false, EchoRef: "loopback"}
	c.Voice.Realtime = RealtimeConfig{Enabled: false, Provider: "acme", Model: "m2", APIKeyRef: "env:RT", BaseURL: "https://rt.example"}
	c.Voice.CloudASRChain = []string{"acme/m2"}
	c.Voice.CloudTTSChain = []string{"acme/m2"}
	c.Audio.InputDevice = "mic-1"
	c.Audio.SampleRate = 48000
	c.Audio.MicMutedDefault = false
	c.LLM.TextChain = []string{"acme/m1", "acme/m2"}
	c.LLM.Roles.Chat = Role{Provider: "acme", Model: "m1", Temperature: 0.3, ThinkingIntensity: ThinkingLow, MaxOutputTokens: 1024}
	c.LLM.Roles.MemoryExtract = Role{Provider: "acme", Model: "m2", Temperature: 0.1, ThinkingIntensity: ThinkingOff}
	c.LLM.Roles.Handoff = Role{Provider: "acme", Model: "m1", ThinkingIntensity: ThinkingOff}
	c.LLM.Roles.Summarize = Role{Provider: "acme", Model: "m2", ThinkingIntensity: ThinkingOff}
	c.LLM.TimeoutMS = 30000
	c.LLM.Retry = RetryConfig{Max: 2, BackoffMS: 500}
	c.LLM.Providers = map[string]Provider{
		"acme": {
			Protocol:             ProtocolOpenAIChat,
			BaseURL:              "https://acme.example/v1",
			APIKeyRef:            "env:ACME",
			Billing:              BillingPlan,
			PlanCreditTotalMicro: 5_000_000,
			Compat:               Compat{Loose: true, AllowMissingUsage: true, ExtraHeaders: map[string]string{"X-A": "b"}},
			RPM:                  60,
			TPM:                  100000,
			Models: map[string]ModelSpec{
				"m1": {
					Display:         "Acme One",
					Capabilities:    Capabilities{Text: true, Vision: true, FC: true, Stream: true},
					ContextWindow:   128000,
					MaxOutputTokens: 4096,
					ThinkingLevels:  []string{ThinkingOff, ThinkingLow},
					Billing:         BillingPayPerToken,
					Price:           Price{In: 1.5, Out: 6, Cached: 0.75, AudioIn: 10, AudioOut: 20},
					QuotaDailyMicro: 1_000_000, QuotaMonthlyMicro: 20_000_000,
					Enabled: true,
				},
				"m2": {Capabilities: Capabilities{Text: true}},
			},
		},
	}
	c.Agent.MaxRounds = 10
	c.Agent.TokenBudget = 1000
	c.Agent.PerToolTimeoutMS = 5000
	c.Agent.LoopGuard = LoopGuard{RepeatThresholds: []int{2, 4, 6}}
	c.Agent.SteeringEnabled = false
	c.Risk.ConfirmTimeoutSec = 60
	c.Risk.L1WindowSec = 5
	c.Risk.ShellEnabled = true
	c.Risk.AllowShellString = true
	c.Risk.ShellAllowlist = []string{"git", "go"}
	c.Risk.BlacklistOverrides = []string{"B-1"}
	c.FS.AllowedDirs = []string{"D:\\x"}
	c.FS.ReparsePointExceptions = []string{"C:\\link"}
	c.FS.DeleteEnabled = true
	c.Net.Allowlist = []string{"api.example.com"}
	c.Net.BlockPrivateRanges = false
	c.Net.Proxy = ProxyConfig{Mode: ProxyManual, URL: "http://proxy:8080"}
	c.Privacy.RedactPaths = true
	c.Privacy.DiagnosticsOptIn = true
	c.Privacy.RetentionDays = 7
	c.Memory.L1Enabled = false
	c.Memory.L1Max = 5
	c.Memory.L3RetentionDays = 10
	c.Memory.ExtractModel = "acme/m2"
	c.Panel.Enabled = false
	c.Panel.Width = 800
	c.Panel.Height = 600
	c.Panel.KeepAliveInSession = false
	c.Panel.Scale = 1.5
	c.Cost.DailyBudget = 1_000_000
	c.Cost.MonthlyBudget = 20_000_000
	c.Cost.AlertThreshold = 0.5
	c.Cost.OverBudget = OverBudgetWarn
	c.Plugins.Tier2Enabled = true
	c.Plugins.Entries = map[string]PluginEntry{
		"clip": {Enabled: true, Capabilities: []string{"clipboard"}, NetAllowlist: []string{"a.com"}, HostAPI: []string{"fs"}},
	}
	c.Models.Dir = "models"
	c.Models.Mirror = []string{"http://mirror"}
	c.Models.LocalOverride = map[string]string{"m1": "D:\\local\\m1"}
	c.Observe.Level = ObserveWarn
	c.Observe.Roll = RollConfig{SizeMB: 5, Days: 3}
	c.Observe.SLOSampleIntervalSec = 2
	if err := validate(c); err != nil {
		t.Fatalf("full config must be valid: %v", err)
	}
	return c
}

// TestRoundTripLoadMarshalLoad pins the GUI double-source guard (D36 rule 4,
// ticket AC): save -> load -> save -> load yields identical structs.
func TestRoundTripLoadMarshalLoad(t *testing.T) {
	c := fullConfig(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := SaveFile(path, c); err != nil {
		t.Fatalf("save: %v", err)
	}
	c2, _, err := LoadFile(path, nil)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if err := SaveFile(path, c2); err != nil {
		t.Fatalf("save2: %v", err)
	}
	c3, _, err := LoadFile(path, nil)
	if err != nil {
		t.Fatalf("reload2: %v", err)
	}
	if !reflect.DeepEqual(c2, c3) {
		t.Fatal("second round-trip diverged (GUI double-source guard)")
	}
	// c2 must equal the saved struct bit-for-bit (defaults merged, nothing
	// added or lost).
	if !reflect.DeepEqual(c, c2) {
		t.Fatalf("loaded struct differs from saved struct")
	}
}

// TestRoundTripBytes pins the canonical marshal/decode symmetry directly
// (no file involved).
func TestRoundTripBytes(t *testing.T) {
	c := fullConfig(t)
	data, err := MarshalCanonical(c)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	c2 := NewDefaults()
	if err := decodeStrict(data, c2); err != nil {
		t.Fatalf("decode: %v", err)
	}
	normalizeConfig(c2)
	if !reflect.DeepEqual(c, c2) {
		t.Fatal("marshal->decode is not the identity on a full config")
	}
}

// TestResolvedNeverPersists pins that resolution is output-only: SaveFile
// writes refs, never resolved plaintext.
func TestResolvedNeverPersists(t *testing.T) {
	c := fullConfig(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := SaveFile(path, c); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	if strings.Contains(string(data), "plaintext-secret") {
		t.Fatal("resolved plaintext must never reach the file")
	}
	if !strings.Contains(string(data), `api_key_ref = 'env:ACME'`) {
		t.Fatal("refs must be written verbatim")
	}
}
