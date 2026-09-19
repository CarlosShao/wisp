package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fakeResolver is an injectable SecretResolver for loader tests.
type fakeResolver struct {
	keys  map[string]string
	procs []string // refs seen, in order
}

func (f *fakeResolver) Resolve(ref string) (string, error) {
	f.procs = append(f.procs, ref)
	if k, ok := f.keys[ref]; ok {
		return k, nil
	}
	return "", errFakeMissing
}

var errFakeMissing = errNotFound{}

type errNotFound struct{}

func (errNotFound) Error() string { return "fake: ref not found" }

// writeConfigFile writes content to a temp config.toml and returns its path.
func writeConfigFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

const validMinimal = `
schema_version = 2

[app]
theme = "light"

[observe]
level = "debug"

[llm.providers.deepseek]
api_key_ref = "env:DEEPSEEK_KEY"

[llm.providers.deepseek.models.deepseek-chat]
context_window = 64000
`

func TestLoadFileMergesDefaultsAndOverrides(t *testing.T) {
	path := writeConfigFile(t, validMinimal)
	c, _, err := LoadFile(path, nil)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if c.App.Theme != "light" {
		t.Errorf("override theme = %q, want light", c.App.Theme)
	}
	if c.Observe.Level != "debug" {
		t.Errorf("override observe.level = %q, want debug", c.Observe.Level)
	}
	if c.Session.WarmTimeoutSec != 90 {
		t.Errorf("default warm_timeout_sec = %d, want 90 (defaults must merge under the file)", c.Session.WarmTimeoutSec)
	}
	if c.SchemaVersion != SchemaVersionCurrent {
		t.Errorf("schema_version = %d, want %d", c.SchemaVersion, SchemaVersionCurrent)
	}
}

func TestLoadFileUnknownKeyErrorNamesLine(t *testing.T) {
	path := writeConfigFile(t, `
schema_version = 2

[app]
theme = "light"
nope = 1

[ball]
size = 56
`)
	_, _, err := LoadFile(path, nil)
	if err == nil {
		t.Fatal("unknown key must be rejected")
	}
	// SPEC-03 sec 2 rule 2: the error names the key AND its line.
	if !strings.Contains(err.Error(), "app.nope") {
		t.Fatalf("error %q must name the full key path app.nope", err)
	}
	if !strings.Contains(err.Error(), "line 6") {
		t.Fatalf("error %q must contain the line number 6", err)
	}
}

func TestLoadFileUnknownKeyLineNumbersAreOneBased(t *testing.T) {
	// First line of the file, no leading newline.
	path := writeConfigFile(t, "schema_version = 2\nbogus_key = true\n")
	_, _, err := LoadFile(path, nil)
	if err == nil {
		t.Fatal("unknown key must be rejected")
	}
	if !strings.Contains(err.Error(), "line 2") {
		t.Fatalf("error %q must report line 2", err)
	}
}

func TestLoadFileWrongTypeErrorHasPosition(t *testing.T) {
	path := writeConfigFile(t, `
schema_version = 2

[app]
theme = 42
`)
	_, _, err := LoadFile(path, nil)
	if err == nil {
		t.Fatal("wrong type must be rejected")
	}
	if !strings.Contains(err.Error(), "line 5") {
		t.Fatalf("type error %q must carry the position (line 5)", err)
	}
}

func TestLoadFileRejectsInvalidValues(t *testing.T) {
	path := writeConfigFile(t, `
schema_version = 2

[audio]
half_duplex = false
`)
	_, _, err := LoadFile(path, nil)
	if err == nil || !strings.Contains(err.Error(), "audio.half_duplex") {
		t.Fatalf("hard-coded read-only must be rejected at load, got: %v", err)
	}
}

func TestLoadFilePresetInheritance(t *testing.T) {
	path := writeConfigFile(t, `
schema_version = 2

[llm.providers.deepseek]
api_key_ref = "env:DEEPSEEK_KEY"

[llm.providers.deepseek.models.deepseek-chat]

[llm.providers.acme]
protocol = "openai-chat"
base_url = "https://acme.example/v1"
`)
	c, _, err := LoadFile(path, nil)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	ds := c.LLM.Providers["deepseek"]
	if ds.Protocol != ProtocolOpenAIChat {
		t.Errorf("preset protocol = %q, want %q (presets fill empty protocol)", ds.Protocol, ProtocolOpenAIChat)
	}
	if ds.BaseURL != "https://api.deepseek.com/v1" {
		t.Errorf("preset base_url = %q, want deepseek endpoint", ds.BaseURL)
	}
	if c.LLM.Providers["acme"].BaseURL != "https://acme.example/v1" {
		t.Errorf("explicit base_url must not be overwritten")
	}
}

func TestLoadFileUnknownProviderNeedsProtocol(t *testing.T) {
	path := writeConfigFile(t, `
schema_version = 2

[llm.providers.acme]

[llm.providers.acme.models.m1]
`)
	_, _, err := LoadFile(path, nil)
	if err == nil || !strings.Contains(err.Error(), "llm.providers.acme.protocol") {
		t.Fatalf("custom provider without protocol must be rejected naming the key, got: %v", err)
	}
}

func TestLoadFileResolvesSecrets(t *testing.T) {
	path := writeConfigFile(t, `
schema_version = 2

[llm.providers.deepseek]
api_key_ref = "env:DEEPSEEK_KEY"

[llm.providers.deepseek.models.deepseek-chat]

[voice.realtime]
enabled = false
api_key_ref = "dpapi:blob-1"
`)
	res := &fakeResolver{keys: map[string]string{
		"env:DEEPSEEK_KEY": "sk-live-1",
		"dpapi:blob-1":     "rt-key",
	}}
	c, got, err := LoadFile(path, res)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if got.ProviderKeys["deepseek"] != "sk-live-1" {
		t.Errorf("provider key = %q, want sk-live-1", got.ProviderKeys["deepseek"])
	}
	if got.RealtimeKey != "rt-key" {
		t.Errorf("realtime key = %q, want rt-key", got.RealtimeKey)
	}
	// The resolved plaintext must never be written back into the config
	// (storage boundary: refs only).
	if strings.Contains(c.LLM.Providers["deepseek"].APIKeyRef, "sk-live-1") {
		t.Error("plaintext leaked into config struct")
	}
	if !contains(res.procs, "env:DEEPSEEK_KEY") || !contains(res.procs, "dpapi:blob-1") {
		t.Errorf("resolver must see every non-empty ref, saw %v", res.procs)
	}
}

func contains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func TestLoadFileResolveFailureIsError(t *testing.T) {
	path := writeConfigFile(t, `
schema_version = 2

[llm.providers.deepseek]
api_key_ref = "env:MISSING_KEY"

[llm.providers.deepseek.models.deepseek-chat]
`)
	res := &fakeResolver{keys: map[string]string{}}
	_, _, err := LoadFile(path, res)
	if err == nil {
		t.Fatal("resolve failure must surface (spec 4.1: failure -> Unconfigured state)")
	}
	if !strings.Contains(err.Error(), "env:MISSING_KEY") {
		t.Fatalf("error %q must name the failing ref", err)
	}
}

func TestLoadFileNilResolverSkipsResolution(t *testing.T) {
	path := writeConfigFile(t, validMinimal)
	c, got, err := LoadFile(path, nil)
	if err != nil {
		t.Fatalf("load with nil resolver: %v", err)
	}
	if got == nil || len(got.ProviderKeys) != 0 {
		t.Errorf("nil resolver must yield empty resolution, got %+v", got)
	}
	if c.LLM.Providers["deepseek"].APIKeyRef != "env:DEEPSEEK_KEY" {
		t.Errorf("refs must survive a nil resolver untouched")
	}
}

func TestLoadFileNewerSchemaVersionRejected(t *testing.T) {
	path := writeConfigFile(t, "schema_version = 99\n")
	_, _, err := LoadFile(path, nil)
	if err == nil {
		t.Fatal("file from a newer build must be rejected")
	}
	if !strings.Contains(err.Error(), "99") {
		t.Fatalf("error %q must mention the offending version", err)
	}
}

func TestSaveFileRoundTripsThroughLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	c := NewDefaults()
	c.SchemaVersion = SchemaVersionCurrent
	c.App.Theme = "light"
	c.Ball.Size = 60
	c.LLM.Providers = map[string]Provider{
		"deepseek": {Protocol: ProtocolOpenAIChat, BaseURL: "https://api.deepseek.com/v1",
			APIKeyRef: "env:K", Models: map[string]ModelSpec{"deepseek-chat": {}}},
	}
	c.Plugins.Entries = map[string]PluginEntry{
		"clip": {Enabled: true, Capabilities: []string{"clipboard"}},
	}
	if err := SaveFile(path, c); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, _, err := LoadFile(path, nil)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if got.App.Theme != "light" || got.Ball.Size != 60 {
		t.Errorf("saved values lost: theme=%q size=%d", got.App.Theme, got.Ball.Size)
	}
	if got.LLM.Providers["deepseek"].APIKeyRef != "env:K" {
		t.Errorf("saved api_key_ref lost")
	}
	if !got.Plugins.Entries["clip"].Enabled {
		t.Errorf("saved plugin entry lost")
	}
	if got.SchemaVersion != SchemaVersionCurrent {
		t.Errorf("saved schema_version = %d", got.SchemaVersion)
	}
}
