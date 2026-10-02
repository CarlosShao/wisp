package config

// Ticket 248 - the settings write path: validate BEFORE the file, one guarded key
// per field, and a receipt that names the key paths it wrote.
//
// AC#7's whole claim is that validate() is not on the write path today, so every
// case below asserts the file, read back the way the next start would read it.
// The pair that carries the load is the same shape ticket 226 used: a guarded write
// that keeps a foreign hand edit (AC#11's first clause), plus the control that runs
// the same scenario through the older snapshot-serializer shape and shows the hand
// edit vanishing. Without the second, the first could be green because the scenario
// never actually carried a second key.

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

// settingsBase is a config with one preset provider and one catalog entry, so the
// settings fields this ticket exposes have something real to address.
const settingsBase = `schema_version = 2

[ball]
size = 56

[llm]
text_chain = ["deepseek/deepseek-chat"]

[llm.roles.chat]
provider = "deepseek"
model = "deepseek-chat"

[llm.providers.deepseek]
api_key_ref = "env:DEEPSEEK_KEY"

[llm.providers.deepseek.models.deepseek-chat]
context_window = 64000

[llm.providers.deepseek.models.deepseek-chat.price]
in = 2500000
out = 10000000
`

func readSettingsFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read config.toml: %v", err)
	}
	return string(data)
}

// AC#11's positive form: one field, one guarded write, and the report says which
// key path the file actually changed by.
func TestAC11SettingWriteReportsTheKeyPathItChanged(t *testing.T) {
	m, path, _ := newTestManager(t, settingsBase)

	written, err := m.SetProviderBaseURL("deepseek", "https://gateway.example.invalid/v1")
	if err != nil {
		t.Fatalf("SetProviderBaseURL: %v", err)
	}
	if want := []string{"llm.providers.deepseek.base_url"}; !reflect.DeepEqual(written, want) {
		t.Errorf("written key paths = %v, want %v: the receipt must name the file's key, not the field's label",
			written, want)
	}
	got := readLoaded(t, path)
	if got.LLM.Providers["deepseek"].BaseURL != "https://gateway.example.invalid/v1" {
		t.Errorf("the value never reached the file: %+v", got.LLM.Providers["deepseek"])
	}
	if got.Ball.Size != 56 {
		t.Errorf("an unrelated key moved: ball.size = %d, want 56", got.Ball.Size)
	}
}

// AC#11's second clause: saving several fields is several writes, each reported by
// its own key path, and a field that fails contributes nothing to the list.
func TestAC11MultiFieldSaveIsOneWritePerKeyAndStopsOnFailure(t *testing.T) {
	m, path, _ := newTestManager(t, settingsBase)

	first, err := m.SetModelContextWindow("deepseek", "deepseek-chat", 200000)
	if err != nil {
		t.Fatalf("first write: %v", err)
	}
	if !reflect.DeepEqual(first, []string{"llm.providers.deepseek.models.deepseek-chat.context_window"}) {
		t.Errorf("first write reported %v", first)
	}
	// The second field of the same save is refused by AC#7's gate (a plaintext key
	// is not a reference), so the receipt for it must be empty and the file must
	// still carry the first.
	second, err := m.SetProviderAPIKeyRef("deepseek", "sk-plaintext-not-a-reference-9f3a2b1c")
	if err == nil {
		t.Fatalf("a plaintext credential value was accepted into config.toml: %v", second)
	}
	if second != nil {
		t.Errorf("a refused write reported key paths %v, want none: nothing landed", second)
	}
	body := readSettingsFile(t, path)
	if !strings.Contains(body, "context_window = 200000") {
		t.Errorf("the accepted field vanished when the next one failed:\n%s", body)
	}
	if strings.Contains(body, "sk-plaintext") {
		t.Errorf("the refused value reached the file:\n%s", body)
	}
	if !strings.Contains(body, `api_key_ref = 'env:DEEPSEEK_KEY'`) {
		t.Errorf("a refused write rewrote the reference it did not accept:\n%s", body)
	}
}

// AC#7: a value that would not survive the start-up validation is refused BEFORE
// anything is written. The file is the assertion, byte for byte.
func TestAC7InvalidValueLeavesTheFileByteIdentical(t *testing.T) {
	cases := []struct {
		name string
		run  func(m *Manager) ([]string, error)
	}{
		{"plaintext credential in the reference field", func(m *Manager) ([]string, error) {
			return m.SetProviderAPIKeyRef("deepseek", "sk-abcdef1234567890ABCDEF")
		}},
		{"negative context window", func(m *Manager) ([]string, error) {
			return m.SetModelContextWindow("deepseek", "deepseek-chat", -1)
		}},
		{"negative price", func(m *Manager) ([]string, error) {
			return m.SetModelPriceIn("deepseek", "deepseek-chat", -0.5)
		}},
		{"chat model that is not in the catalog", func(m *Manager) ([]string, error) {
			return m.SetRoleChatModel("no-such-model-9")
		}},
		{"a provider the config does not declare", func(m *Manager) ([]string, error) {
			return m.SetProviderBaseURL("inventco", "https://x.example.invalid")
		}},
		{"a model the provider's catalog does not list", func(m *Manager) ([]string, error) {
			return m.SetModelContextWindow("deepseek", "ghost-model", 32000)
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m, path, _ := newTestManager(t, settingsBase)
			before := readSettingsFile(t, path)
			written, err := tc.run(m)
			if err == nil {
				t.Fatalf("the write was accepted (reported %v): AC#7's gate did not run", written)
			}
			if written != nil {
				t.Errorf("a refused write reported key paths %v, want none", written)
			}
			if after := readSettingsFile(t, path); after != before {
				t.Errorf("config.toml changed under a refused write.\nbefore:\n%s\nafter:\n%s", before, after)
			}
			// The next start must still come up: that is the whole reason AC#7
			// exists (a settings page that can brick the next boot).
			if got := readLoaded(t, path); got.LLM.Providers["deepseek"].Models["deepseek-chat"].ContextWindow != 64000 {
				t.Errorf("the file this process refused is not the file the next start reads: %+v", got.LLM.Providers)
			}
			// And the in-memory config the runtime acts on cannot diverge either.
			if got := m.Config().LLM.Roles.Chat.Model; got != "deepseek-chat" {
				t.Errorf("memory diverged from the file under a refused write: chat model = %q", got)
			}
		})
	}
}

// AC#11's first clause: a foreign hand edit of a DIFFERENT key, already on disk when
// the settings write happens, must still be there verbatim afterwards.
func TestAC11SettingWriteKeepsAForeignHandEditedKey(t *testing.T) {
	m, path, mutate := newTestManager(t, settingsBase)
	mutate(strings.Replace(settingsBase, "size = 56", "size = 60", 1))

	before := readSettingsFile(t, path)
	if !strings.Contains(before, "size = 60") {
		t.Fatalf("the fixture never carried the hand edit: %s", before)
	}
	if _, err := m.SetProviderBaseURL("deepseek", "https://kept.example.invalid/v1"); err != nil {
		t.Fatalf("SetProviderBaseURL: %v", err)
	}
	after := readSettingsFile(t, path)
	if !strings.Contains(after, "size = 60") {
		t.Errorf("the settings write reverted the operator's [ball] size:\n%s", after)
	}
	if !strings.Contains(after, "https://kept.example.invalid/v1") {
		t.Errorf("the settings write did not land its own key:\n%s", after)
	}
}

// TestAC11ControlSnapshotWriteIsWhatRevertsTheHandEdit is the positive control for
// the case above: the same file, the same hand edit, and the write made the way a
// settings route would make it WITHOUT the guarded merge (hand the in-memory
// snapshot to the serializer). It must revert the hand edit, because that is the
// hole AC#11 forbids this ticket from reopening - and if this control ever goes
// green, the scenario above has stopped being able to detect anything.
func TestAC11ControlSnapshotWriteIsWhatRevertsTheHandEdit(t *testing.T) {
	m, path, mutate := newTestManager(t, settingsBase)
	mutate(strings.Replace(settingsBase, "size = 56", "size = 60", 1))

	// Nothing re-read here: this is memory, which still holds size = 56, straight
	// to the serializer - ticket 226's measured old behaviour.
	if _, err := m.SetProviderBaseURL("deepseek", "https://guarded.example.invalid/v1"); err != nil {
		t.Fatalf("guarded write as a control precondition: %v", err)
	}
	stale := m.Config()
	if err := SaveFile(path, stale); err != nil {
		t.Fatalf("control write: %v", err)
	}
	if body := readSettingsFile(t, path); strings.Contains(body, "size = 60") {
		t.Errorf("the control no longer reverts a hand edit, so TestAC11SettingWriteKeepsAForeignHandEditedKey " +
			"is green for a scenario that cannot fail - re-check it before believing it")
	}
}

// A reference the schema accepts must land, and the value it names must not be in
// the file: the settings path writes refs, the store keeps values.
func TestAC3SettingWriteStoresARefNeverAValue(t *testing.T) {
	m, path, _ := newTestManager(t, settingsBase)

	written, err := m.SetProviderAPIKeyRef("deepseek", "dpapi:248canaryblob01")
	if err != nil {
		t.Fatalf("SetProviderAPIKeyRef: %v", err)
	}
	if want := []string{"llm.providers.deepseek.api_key_ref"}; !reflect.DeepEqual(written, want) {
		t.Errorf("written = %v, want %v", written, want)
	}
	body := readSettingsFile(t, path)
	if !strings.Contains(body, `api_key_ref = 'dpapi:248canaryblob01'`) {
		t.Errorf("the reference did not land:\n%s", body)
	}
	// Nothing on this path can hold the secret itself, because the schema has no
	// such field: assert the one credential-bearing leaf keeps a reference form.
	got := readLoaded(t, path)
	ref := got.LLM.Providers["deepseek"].APIKeyRef
	if !strings.HasPrefix(ref, "dpapi:") && !strings.HasPrefix(ref, "env:") {
		t.Errorf("api_key_ref holds %q, which is neither reference form", ref)
	}
}

// The already-carries-this-value shape must not rewrite the file (mergeWrite's
// second outcome, inherited by every setter that rides it).
func TestAC11WritingTheValueTheFileAlreadyHoldsChangesNoBytes(t *testing.T) {
	m, path, _ := newTestManager(t, settingsBase)
	before := readSettingsFile(t, path)

	written, err := m.SetProviderBaseURL("deepseek", "https://gateway.example.invalid/v1")
	if err != nil {
		t.Fatalf("first write: %v", err)
	}
	if len(written) != 1 {
		t.Fatalf("first write reported %v, want one key path", written)
	}
	second, err := m.SetProviderBaseURL("deepseek", "https://gateway.example.invalid/v1")
	if err != nil {
		t.Fatalf("second identical write: %v", err)
	}
	if len(second) != 0 {
		t.Errorf("an identical write reported changed keys %v, want none", second)
	}
	if after := readSettingsFile(t, path); strings.Count(after, "https://gateway.example.invalid/v1") != 1 {
		t.Errorf("the value was written twice into the same file:\n%s", after)
	}
	if readSettingsFile(t, path) == before {
		t.Errorf("the file never changed at all, so the first write did not land: %s", before)
	}
}
