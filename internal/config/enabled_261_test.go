package config

// Ticket 261 leg p1 - INSTRUMENTATION ONLY (zero production-code change).
//
// What these rulers measure, run and not inferred:
//
//  1. schema.go:369-371 promises "Enabled removes the model from
//     discovery/selection ..." and carries `default:"true"`, but applyDefaults
//     (defaults.go:77-78 `case reflect.Map:` -> `// leave nil`) never enters a
//     map, so a HAND-WRITTEN entry without the key decodes to the Go zero
//     value false. Ruler: load a hand-written file with the three shapes
//     (no key / = false / = true) through the real LoadFile and read the
//     decoded field back.
//  2. The settings write side (settings.go:115/153) refuses to mint a catalog
//     row (requireCatalogEntry), so it cannot CREATE the false - but every
//     write goes through mergeWrite -> SaveFile -> MarshalCanonical, which
//     emits "every static field explicitly (no omitempty)" (parse.go:160-163),
//     so one unrelated settings save MATERIALIZES `enabled = false` into the
//     file. Ruler: perform one real SetModelContextWindow write and read the
//     file bytes.
//
// Both tests fail loudly if AC#1 (either shape) ever lands - they are the
// reading of TODAY, pinned so the next leg cannot "fix" the promise by only
// changing one of the two ends.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// enabledBase is a hand-written shape: text_chain names all three entries, so
// anything in the load pipeline that meant to honour `enabled` had its chance
// here and would have rejected the first two.
const enabledBase = `schema_version = 2

[llm]
text_chain = ["mock261/t261-nokey", "mock261/t261-false", "mock261/t261-true"]

[llm.providers.mock261]
protocol = "openai-chat"
base_url = "http://127.0.0.1:9/v1"

[llm.providers.mock261.models.t261-nokey]
context_window = 4096

[llm.providers.mock261.models.t261-false]
context_window = 4096
enabled = false

[llm.providers.mock261.models.t261-true]
context_window = 4096
enabled = true
`

func writeEnabledFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config.toml: %v", err)
	}
	return path
}

// TestTicket261P1MissingEnabledKeyDecodesFalse is ruler 1: the
// `default:"true"` tag has no enforcer inside map entries.
func TestTicket261P1MissingEnabledKeyDecodesFalse(t *testing.T) {
	path := writeEnabledFile(t, enabledBase)

	cfg, _, err := LoadFile(path, nil)
	if err != nil {
		t.Fatalf("LoadFile refused the hand-written file: %v", err)
	}
	models := cfg.LLM.Providers["mock261"].Models
	got := map[string]bool{}
	for id, spec := range models {
		got[id] = spec.Enabled
	}
	// The reading, per shape.
	if got["t261-nokey"] {
		t.Errorf("t261-nokey Enabled = true, want false: applyDefaults never enters maps, so the missing key must land on the Go zero value")
	}
	if got["t261-false"] {
		t.Errorf("t261-false Enabled = true, want false")
	}
	if !got["t261-true"] {
		t.Errorf("t261-true Enabled = false, want true (the control: an explicit true does survive the decode)")
	}
	// The companion half of the promise ("... without deleting its catalog
	// entry" implies validate at least knows about enabled): validation ran on
	// a text_chain that names two non-enabled models and let the file through
	// unchanged - LoadFile succeeded above, and all three elements are still
	// in the chain.
	if len(cfg.LLM.TextChain) != 3 {
		t.Errorf("text_chain = %v, want the hand-written 3 elements intact: nothing in the load pipeline filters or rejects chain elements pointing at non-enabled models", cfg.LLM.TextChain)
	}
}

// TestTicket261P1SettingsWriteMaterializesFalse is ruler 2: the settings side
// cannot invent an entry, but it does write the false onto disk.
func TestTicket261P1SettingsWriteMaterializesFalse(t *testing.T) {
	m, path, _ := newTestManager(t, enabledBase)

	if _, err := m.SetModelContextWindow("mock261", "t261-nokey", 8192); err != nil {
		t.Fatalf("SetModelContextWindow: %v", err)
	}
	raw := readSettingsFile(t, path)
	if !strings.Contains(raw, "enabled = false") {
		t.Errorf("the file after one settings write does not carry `enabled = false`:\n%s", raw)
	}
	// And the written state reloads as false on both sides of the file - the
	// settings path never repairs or rescues the flag.
	if cfg, err := readConfigFile(path); err != nil {
		t.Fatalf("reload after settings write: %v", err)
	} else if cfg.LLM.Providers["mock261"].Models["t261-nokey"].Enabled {
		t.Errorf("t261-nokey reloaded as enabled after the settings write")
	} else if cfg.LLM.Providers["mock261"].Models["t261-nokey"].ContextWindow != 8192 {
		t.Errorf("the owned key moved but not the flag check base: context_window = %d, want 8192",
			cfg.LLM.Providers["mock261"].Models["t261-nokey"].ContextWindow)
	}

	// Control (the ticket's own boundary): the settings side refuses to mint a
	// catalog row, so it is NOT the producer of the false shape - requireCatalog-
	// Entry names that refusal, and it must still fire.
	if _, err := m.SetModelContextWindow("mock261", "t261-ghost", 1); err == nil {
		t.Errorf("settings write invented a catalog row t261-ghost")
	} else if !strings.Contains(err.Error(), "refusing to invent one") {
		t.Errorf("unexpected refusal: %v", err)
	}
}
