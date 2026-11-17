package config

// Readface wave 1 (2026-11-17) - behavior teeth for the two new read faces:
// SectionRows (ticket 39 AC#1's sections[]{key,value,tier}) and the
// errors[]{line,message} face (AC#4). Both are BEHAVIOR assertions over real
// loads, never source-string pins.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeCfg456(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write config fixture: %v", err)
	}
	return path
}

func sectionRow456(rows []SectionRow, key string) (SectionRow, bool) {
	for _, r := range rows {
		if r.Key == key {
			return r, true
		}
	}
	return SectionRow{}, false
}

// TestSectionRowsCoversEveryD36Section is AC#1's completeness half: every
// top-level section of the schema walks out with its own rows, and the tier
// words are the registry's own words.
func TestSectionRowsCoversEveryD36Section(t *testing.T) {
	cfg := NewDefaults()
	rows := SectionRows(cfg)
	if len(rows) == 0 {
		t.Fatal("SectionRows over a default config answered nothing - the walk is reading nothing")
	}
	// The 17 top-level sections of schema.go's Config (D36's section tree).
	sections := []string{"app", "ball", "hotkey", "session", "voice", "audio",
		"llm", "agent", "risk", "fs", "net", "privacy", "memory", "panel",
		"cost", "plugins", "models", "observe"}
	seen := map[string]bool{}
	for _, r := range rows {
		seen[strings.SplitN(r.Key, ".", 2)[0]] = true
	}
	for _, s := range sections {
		if !seen[s] {
			t.Errorf("section %q produced no rows - the D36 section tree is not fully walked (AC#1)", s)
		}
	}
}

// TestSectionRowsTierWordsComeFromTheRegistry pins the tier half per key:
// per-key registry rows win for [app]/[voice], whole-section rows for the
// rest, and locked sections carry the registry's own "locked" marker.
func TestSectionRowsTierWordsComeFromTheRegistry(t *testing.T) {
	cfg := NewDefaults()
	cfg.Ball.Size = 90
	cfg.App.Language = "en-US"
	cfg.Voice.ASR.Model = "streaming-zh"
	rows := SectionRows(cfg)

	cases := []struct {
		key  string
		tier string
	}{
		{"app.theme", "hot"},
		{"app.language", "restart"},
		{"ball.size", "hot"},
		{"voice.asr.model", "reload"},
		{"voice.tts.speed", "hot"},
		{"privacy.redact_paths", "hot"},
		{"cost.daily_budget", "hot"},
		{"risk.permission_mode", "locked"},
		{"fs.allowed_dirs", "locked"},
		{"net.allowlist", "locked"},
		{"plugins.tier2_enabled", "locked"},
	}
	for _, c := range cases {
		r, ok := sectionRow456(rows, c.key)
		if !ok {
			t.Errorf("no row for %q (AC#1)", c.key)
			continue
		}
		if r.Tier != c.tier {
			t.Errorf("row %q tier = %q, want %q (the registry's own word)", c.key, r.Tier, c.tier)
		}
	}
}

// TestSectionRowsCarriesWrittenValues walks a config with non-default values
// and requires the row value to be the written fact, not the default.
func TestSectionRowsCarriesWrittenValues(t *testing.T) {
	cfg := NewDefaults()
	cfg.Ball.Size = 90
	cfg.Panel.FontSize = 14
	cfg.LLM.Roles.Chat.Model = "step-2-16k"
	cfg.Voice.WakeWord.VetoWords = []string{"取消", "停下"}
	rows := SectionRows(cfg)

	if r, ok := sectionRow456(rows, "ball.size"); !ok || r.Value != "90" {
		t.Errorf("ball.size row = %+v ok=%v, want value 90", r, ok)
	}
	if r, ok := sectionRow456(rows, "panel.font_size"); !ok || r.Value != "14" {
		t.Errorf("panel.font_size row = %+v ok=%v, want value 14", r, ok)
	}
	if r, ok := sectionRow456(rows, "llm.roles.chat.model"); !ok || r.Value != "step-2-16k" {
		t.Errorf("llm.roles.chat.model row = %+v ok=%v, want the written model", r, ok)
	}
	r, ok := sectionRow456(rows, "voice.wake_word.veto_words")
	if !ok {
		t.Fatalf("no row for voice.wake_word.veto_words")
	}
	if !strings.Contains(r.Value, "取消") || !strings.HasPrefix(r.Value, "[") {
		t.Errorf("composite row %q does not render the written array: %q", r.Key, r.Value)
	}
}

// TestSectionRowsRefusesNilConfig keeps the honest-empty contract: nil is
// "did not read", never an empty tree.
func TestSectionRowsRefusesNilConfig(t *testing.T) {
	if rows := SectionRows(nil); rows != nil {
		t.Errorf("SectionRows(nil) = %d rows, want nil", len(rows))
	}
}

// TestSectionRowsWalksPluginEntries covers the dynamic half: per-plugin
// sub-tables surface as plugins.<id>.<key> rows with the locked marker.
func TestSectionRowsWalksPluginEntries(t *testing.T) {
	cfg := NewDefaults()
	cfg.Plugins.Tier2Enabled = true
	cfg.Plugins.Entries = map[string]PluginEntry{
		"demo": {Enabled: true, Capabilities: []string{"net.fetch"}},
	}
	rows := SectionRows(cfg)
	if r, ok := sectionRow456(rows, "plugins.tier2_enabled"); !ok || r.Value != "true" || r.Tier != "locked" {
		t.Errorf("plugins.tier2_enabled row = %+v ok=%v", r, ok)
	}
	if r, ok := sectionRow456(rows, "plugins.demo.enabled"); !ok || r.Value != "true" || r.Tier != "locked" {
		t.Errorf("plugins.demo.enabled row = %+v ok=%v (the dynamic sub-table must walk)", r, ok)
	}
}

// TestParseErrorsNamesUnknownKeysWithLines is AC#4's core: an externally
// edited file with an unknown key reports that key's line.
func TestParseErrorsNamesUnknownKeysWithLines(t *testing.T) {
	body := "schema_version = 2\n\n[app]\ntheme = \"dark\"\nbogus_key = 1\n"
	path := writeCfg456(t, body)
	rows := ParseErrorsFile(path)
	if len(rows) == 0 {
		t.Fatal("a file with an unknown key produced no errors")
	}
	var found bool
	for _, r := range rows {
		if strings.Contains(r.Message, "bogus_key") {
			found = true
			if r.Line != 5 {
				t.Errorf("unknown key row line = %d, want 5 (the key's own 1-based line)", r.Line)
			}
		}
	}
	if !found {
		t.Errorf("no row names the unknown key: %v", rows)
	}
}

// TestParseErrorsNamesTypeErrorsWithLines: a wrong type carries its line too.
func TestParseErrorsNamesTypeErrorsWithLines(t *testing.T) {
	body := "schema_version = 2\n\n[ball]\nsize = \"huge\"\n"
	path := writeCfg456(t, body)
	rows := ParseErrorsFile(path)
	if len(rows) == 0 {
		t.Fatal("a file with a wrong-typed value produced no errors")
	}
	if rows[0].Line != 4 {
		t.Errorf("type-error row line = %d, want 4 (the offending value's line)", rows[0].Line)
	}
	if rows[0].Message == "" {
		t.Error("type-error row carries an empty message")
	}
}

// TestParseErrorsCleanFileAnswersNil: a file that loads clean is not an error
// list with zero rows of noise - it is nil.
func TestParseErrorsCleanFileAnswersNil(t *testing.T) {
	path := writeCfg456(t, "schema_version = 2\n\n[app]\ntheme = \"dark\"\n")
	if rows := ParseErrorsFile(path); rows != nil {
		t.Errorf("a clean file produced %d error rows, want nil", len(rows))
	}
}

// TestManagerParseErrorsReadsItsOwnPath pins the Manager read face over the
// same bytes: a broken file on disk yields line-numbered rows while the
// Manager's in-memory config is untouched.
func TestManagerParseErrorsReadsItsOwnPath(t *testing.T) {
	good := writeCfg456(t, "schema_version = 2\n")
	m, err := NewManager(good, nil)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	if rows := m.ParseErrors(); rows != nil {
		t.Errorf("clean file via Manager produced %d rows, want nil", len(rows))
	}
	bad := writeCfg456(t, "schema_version = 2\n\n[aap]\nnot_a_key = true\n")
	if err := os.Rename(bad, good); err != nil {
		t.Fatalf("replace: %v", err)
	}
	rows := m.ParseErrors()
	if len(rows) == 0 {
		t.Fatal("Manager.ParseErrors answered nothing for a broken file")
	}
	if !strings.Contains(rows[0].Message, "aap") || rows[0].Line != 3 {
		// go-toml reports an unknown TABLE at the table header's own line.
		t.Errorf("Manager.ParseErrors = %+v, want the unknown key with its line", rows)
	}
}
