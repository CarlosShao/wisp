package main

// Ticket 187 at the assembly root, the write-leg tooth the 187 leg left owed:
// the panel's third writable roster carries role_chat_thinking, and this file
// proves the store's ApplySetting switch actually routes it to
// config.Manager.SetRoleChatThinking - the one case line the panel door's
// roster alone cannot prove (a roster entry without a case lands in the
// default branch and is refused, which is the honest shape, but then the
// field the page offers would never write).
//
// Teeth, in order:
//  1. the write lands: a real Manager, the real guarded setter, the receipt
//     names the one key, and the FILE on disk carries the new tier;
//  2. the case is load-bearing: comment out the FieldRoleChatThinking case in
//     panel_config_store.go and the first test goes red (default branch
//     refuses) - the mutation was run and recorded in
//     docs/evidence/s1/187-apply-setting-case.md;
//  3. the refusal path: a word outside the four-value vocabulary is refused
//     by the setter before any byte moves, and neither memory nor file
//     changes.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
)

// thinking187Base is the same minimal valid config the config side's own 187
// tests load: one provider, one enabled model declaring off|low, the chat
// role pinned to it with thinking_intensity = "off".
const thinking187Base = `schema_version = 2

[llm]
text_chain = ["deepseek/deepseek-chat"]

[llm.roles.chat]
provider = "deepseek"
model = "deepseek-chat"
thinking_intensity = "off"

[llm.providers.deepseek]
api_key_ref = "env:DEEPSEEK_KEY"

[llm.providers.deepseek.models.deepseek-chat]
enabled = true
context_window = 64000
thinking_levels = ["off", "low"]

[llm.providers.deepseek.models.deepseek-chat.price]
in = 2500000
out = 10000000
`

func thinking187Store(t *testing.T) (*configStore, *config.Manager, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, configFileName)
	if err := os.WriteFile(path, []byte(thinking187Base), 0o600); err != nil {
		t.Fatal(err)
	}
	mgr, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("manager: %v", err)
	}
	return newConfigStore(mgr, nil, nil, "test-187"), mgr, path
}

// Test187ApplySettingThinkingLandsTheKey: the case line routes the panel's
// value into SetRoleChatThinking, the receipt names the one dotted key, the
// in-memory config the next task's re-resolution reads carries the tier, and
// the file on disk says the same. Delete the case in panel_config_store.go
// and this dies on the first Fatalf (the default branch refuses the field).
func Test187ApplySettingThinkingLandsTheKey(t *testing.T) {
	s, mgr, path := thinking187Store(t)

	res, err := s.ApplySetting(t.Context(), panel.SettingWrite{
		Field: panel.FieldRoleChatThinking, Value: "low", Actor: "test-187",
	})
	if err != nil {
		t.Fatalf("ApplySetting(%s=low): %v", panel.FieldRoleChatThinking, err)
	}
	if len(res.Written) != 1 || res.Written[0] != "llm.roles.chat.thinking_intensity" {
		t.Fatalf("receipt written = %v, want [llm.roles.chat.thinking_intensity]", res.Written)
	}
	if res.Tier != panel.EffectiveRestart {
		// The store stamps the [llm] blanket restart; the panel door corrects
		// this roster's field to next_task via fieldTier187. Here we pin the
		// store's own stamp so a future leg cannot quietly re-tier it.
		t.Fatalf("store tier = %v, want restart (the panel door re-tiers next_task)", res.Tier)
	}
	if got := mgr.Config().LLM.Roles.Chat.ThinkingIntensity; got != "low" {
		t.Fatalf("in-memory thinking = %q, want low (what the next task would re-resolve)", got)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	// The writer normalizes the whole file (quoting style included), so the
	// disk check scans for the key line, not a byte-exact original phrase.
	landed := false
	for _, line := range strings.Split(string(raw), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "thinking_intensity") && strings.Contains(line, "low") {
			landed = true
		}
	}
	if !landed {
		t.Fatalf("file on disk does not carry the write:\n%s", raw)
	}
}

// Test187ApplySettingThinkingRefusesOutsideVocabulary: the setter owns the
// enum (ticket 320's ruling shape) and refuses before any byte moves; the
// config stays "off" in memory and on disk, and the error names the field.
func Test187ApplySettingThinkingRefusesOutsideVocabulary(t *testing.T) {
	s, mgr, path := thinking187Store(t)

	res, err := s.ApplySetting(t.Context(), panel.SettingWrite{
		Field: panel.FieldRoleChatThinking, Value: "maximum", Actor: "test-187",
	})
	if err == nil {
		t.Fatalf("ApplySetting(maximum) must be refused, got receipt %+v", res)
	}
	if res.Tier != panel.EffectiveNotApplied {
		t.Fatalf("refused tier = %v, want not_applied", res.Tier)
	}
	if !strings.Contains(err.Error(), "llm.roles.chat.thinking_intensity") {
		// The setter's refusal names the dotted key path (its own vocabulary
		// sentence), not the panel field name; pin the key path.
		t.Fatalf("refusal must name the key path, got: %v", err)
	}
	if got := mgr.Config().LLM.Roles.Chat.ThinkingIntensity; got != "off" {
		t.Fatalf("memory changed to %q on a refused write", got)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), `"maximum"`) {
		t.Fatalf("file on disk changed on a refused write:\n%s", raw)
	}
}
