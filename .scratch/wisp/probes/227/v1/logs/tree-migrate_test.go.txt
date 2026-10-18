package config

import (
	"os"
	"slices"
	"strings"
	"testing"
)

// equalStrings is a slices.Equal alias for test readability.
func equalStrings(a, b []string) bool { return slices.Equal(a, b) }

// v1Fixture is a pre-catalog (schema v1) config: [llm] with
// default_provider/fallback_provider/timeout/temperature and flat
// providers.<name>.{base_url, model, api_key_ref} (see SchemaVersionCurrent
// history in schema.go).
const v1Fixture = `
schema_version = 1

[app]
theme = "dark"

[ball]
size = 60

[risk]
l1_window_sec = 5

[fs]
allowed_dirs = ["D:\\keep"]

[llm]
default_provider = "openai"
fallback_provider = "deepseek"
timeout = 45000
temperature = 0.5

[llm.providers.openai]
base_url = "https://api.openai.com/v1"
model = "gpt-4o-mini"
api_key_ref = "env:OPENAI_KEY"

[llm.providers.deepseek]
base_url = "https://api.deepseek.com/v1"
model = "deepseek-chat"
`

func TestMigrateV1Fixture(t *testing.T) {
	path := writeConfigFile(t, v1Fixture)
	before, _ := os.ReadFile(path)
	c, _, err := LoadFile(path, nil)
	if err != nil {
		t.Fatalf("migrating load: %v", err)
	}

	// Backup exists and carries the pre-migration original.
	backup, err := os.ReadFile(path + ".bak-1")
	if err != nil {
		t.Fatalf("backup config.toml.bak-1 must exist: %v", err)
	}
	if string(backup) != string(before) {
		t.Fatal("backup must be byte-identical to the pre-migration file")
	}

	// The file itself was upgraded in place.
	after, _ := os.ReadFile(path)
	if !strings.Contains(string(after), "schema_version = 2") {
		t.Fatalf("migrated file must carry schema_version = 2, got:\n%s", after)
	}

	// v1 values mapped into the v2 model (nothing silently reset).
	if c.LLM.TimeoutMS != 45000 {
		t.Errorf("llm.timeout -> timeout_ms = %d, want 45000", c.LLM.TimeoutMS)
	}
	if c.LLM.Roles.Chat.Temperature != 0.5 {
		t.Errorf("llm.temperature -> roles.chat.temperature = %v, want 0.5", c.LLM.Roles.Chat.Temperature)
	}
	if c.LLM.Roles.Chat.Provider != "openai" || c.LLM.Roles.Chat.Model != "gpt-4o-mini" {
		t.Errorf("roles.chat = %s/%s, want openai/gpt-4o-mini", c.LLM.Roles.Chat.Provider, c.LLM.Roles.Chat.Model)
	}
	wantChain := []string{"openai/gpt-4o-mini", "deepseek/deepseek-chat"}
	if !equalStrings(c.LLM.TextChain, wantChain) {
		t.Errorf("text_chain = %v, want %v (default first, fallback second)", c.LLM.TextChain, wantChain)
	}
	m := c.LLM.Providers["openai"].Models["gpt-4o-mini"]
	if !m.Enabled {
		t.Errorf("migrated flat model must become a catalog entry (enabled), got %+v", m)
	}
	if c.LLM.Providers["openai"].APIKeyRef != "env:OPENAI_KEY" {
		t.Errorf("api_key_ref must survive migration")
	}
	if !equalStrings(c.FS.AllowedDirs, []string{"D:\\keep"}) {
		t.Errorf("unrelated sections must pass through untouched, got %v", c.FS.AllowedDirs)
	}
	if c.Ball.Size != 60 || c.App.Theme != "dark" {
		t.Errorf("non-llm values must survive migration")
	}
	if c.Risk.L1WindowSec != 5 {
		t.Errorf("[risk] must pass through, got l1_window_sec=%d", c.Risk.L1WindowSec)
	}

	// A reload of the migrated file loads cleanly (idempotent).
	m2, _, err := LoadFile(path, nil)
	if err != nil {
		t.Fatalf("reload migrated file: %v", err)
	}
	if m2.LLM.TimeoutMS != 45000 {
		t.Errorf("reloaded value lost")
	}
}

func TestMigrateNoVersionKeyAssumedV1(t *testing.T) {
	noVer := strings.Replace(v1Fixture, "schema_version = 1\n", "", 1)
	path := writeConfigFile(t, noVer)
	if _, _, err := LoadFile(path, nil); err != nil {
		t.Fatalf("file without schema_version is treated as v1: %v", err)
	}
	if _, err := os.ReadFile(path + ".bak-1"); err != nil {
		t.Fatalf("backup must exist: %v", err)
	}
}

func TestMigrateCorruptFileUntouched(t *testing.T) {
	corrupt := "schema_version = 1\n[llm\nbroken ==="
	path := writeConfigFile(t, corrupt)
	before, _ := os.ReadFile(path)
	_, _, err := LoadFile(path, nil)
	if err == nil {
		t.Fatal("corrupt file must fail the load (caller maps to Unconfigured + guidance)")
	}
	if !strings.Contains(err.Error(), "migrat") {
		t.Fatalf("error must come from the migration path with explicit guidance, got: %v", err)
	}
	after, _ := os.ReadFile(path)
	if string(before) != string(after) {
		t.Fatal("unmigratable file must be left untouched (silently resetting would wipe allowlists - security bug)")
	}
	if _, err := os.Stat(path + ".bak-1"); err == nil {
		t.Fatal("no backup must be created for a failed migration")
	}
}

func TestMigrateUnsupportedVersionUntouched(t *testing.T) {
	path := writeConfigFile(t, "schema_version = 1\n[app]\ntheme = \"dark\"\n")
	// Direct unit test of the chain: no step registered for 7.
	_, err := applyMigrations(path, []byte("schema_version = 7\n"), 7)
	if err == nil {
		t.Fatal("missing migration step must be an explicit error")
	}
	if !strings.Contains(err.Error(), "7") {
		t.Fatalf("error must name the stuck-at version, got: %v", err)
	}
}

func TestMigrateHigherIntermediateVersionRejected(t *testing.T) {
	// A file claiming a version beyond the registry cannot be migrated
	// (LoadFile rejects > current before dispatch; this pins the chain too).
	path := writeConfigFile(t, "schema_version = 99\n")
	_, _, err := LoadFile(path, nil)
	if err == nil || !strings.Contains(err.Error(), "newer build") {
		t.Fatalf("newer schema_version must be rejected with guidance, got: %v", err)
	}
}
