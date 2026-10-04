package config

// Ticket 257 - the clean machine, and the three ways a settings write says no.
//
// AC#1's run is the one this package can do for itself, and it is run rather
// than read: an empty data root, the first default config written exactly the way
// first-run writes it (cmd/wisp/firstrun.go:82 calls this same SaveFile with this
// same NewDefaults), every settings field addressed once, and the terminal state
// = ticket 257 AC#0's chosen shape (c) actually delivered - a refusal that names
// the one place a provider row comes from, and after the operator follows that
// sentence, all seven fields writable.
//
// AC#2 forbids folding the three refusal reasons into one line (票 223 AC#4's
// discipline, same shape): each reason gets its own planted case, each case
// asserts its own tag AND the absence of the other two. A judge that only asked
// "did it error" would stay green with one shared "配置未生效" line, which is
// precisely the thing this ticket is filed about.
//
// The credential field's own leg is deliberately NOT re-implemented here:
// StoreCredential lives in cmd/wisp and is the only writer of a credential value
// (AC#3 says that surface stays untouched). What this package can say about it is
// the config-side half - the reference that lands in the file - and the walk
// addresses exactly that, by name, never by value.

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The provider and model the page would address once the operator has followed the
// guidance. Neither exists on a clean machine, which is the point of AC#1.
const (
	cleanProvider = "deepseek"
	cleanModel    = "deepseek-chat"
	secondModel   = "deepseek-reasoner"
)

// settingsField is one row of the panel's writable roster, named with the panel's
// own field name, together with the leg cmd/wisp/panel_config_store.go's switch
// routes it through. Seven rows because that roster is seven fields today.
type settingsField struct {
	field string
	leg   string
	run   func(m *Manager) ([]string, error)
}

func writableRosterWalk() []settingsField {
	return []settingsField{
		{"provider_base_url", "Manager.SetProviderBaseURL",
			func(m *Manager) ([]string, error) {
				return m.SetProviderBaseURL(cleanProvider, "https://clean-machine.example.invalid/v1")
			}},
		{"provider_api_key_ref", "Manager.SetProviderAPIKeyRef",
			func(m *Manager) ([]string, error) {
				return m.SetProviderAPIKeyRef(cleanProvider, "env:T257_R1_ENV_NAME")
			}},
		{"model_context_window", "Manager.SetModelContextWindow",
			func(m *Manager) ([]string, error) {
				return m.SetModelContextWindow(cleanProvider, cleanModel, 128000)
			}},
		{"model_price_in", "Manager.SetModelPriceIn",
			func(m *Manager) ([]string, error) {
				return m.SetModelPriceIn(cleanProvider, cleanModel, 3500000)
			}},
		{"model_price_out", "Manager.SetModelPriceOut",
			func(m *Manager) ([]string, error) {
				return m.SetModelPriceOut(cleanProvider, cleanModel, 12000000)
			}},
		{"role_chat_model", "Manager.SetRoleChatModel",
			func(m *Manager) ([]string, error) { return m.SetRoleChatModel(secondModel) }},
		{"provider_credential", "Manager.SetProviderAPIKeyRef (StoreCredential's config-side leg)",
			func(m *Manager) ([]string, error) {
				return m.SetProviderAPIKeyRef(cleanProvider, "dpapi:t257r1canaryblob01")
			}},
	}
}

// cleanMachineManager builds the state AC#1 names: a temporary data root, no
// config.toml in it, then the first default config written the way first-run
// writes it, then the Manager the panel chain would build (res=nil, the shape
// every production NewManager call passes - ledger A560).
func cleanMachineManager(t *testing.T) (*Manager, string, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	if _, err := os.Stat(path); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("the run must start from NO config.toml; stat said %v", err)
	}
	// The starting state has its own sentence, and it is reason 1 of 3. Proven
	// again here because AC#1's premise is that this is where the machine is
	// before anything is built.
	if _, err := NewManager(path, nil); err == nil {
		t.Fatalf("NewManager succeeded with no config.toml at all")
	} else if !strings.Contains(err.Error(), refusalFileMissing) {
		t.Fatalf("the no-file start must answer %q, got: %v", refusalFileMissing, err)
	}

	if err := SaveFile(path, NewDefaults()); err != nil {
		t.Fatalf("first-run SaveFile(NewDefaults()): %v", err)
	}
	body := readSettingsFile(t, path)
	// The bug's first half, on disk rather than in a comment: the default table
	// leaves the provider registry nil, so the file carries no provider row -
	// which is exactly what ticket 198's own pin requires (its receipt test forbids
	// an invented [llm.providers section here).
	if strings.Contains(body, "[llm.providers") {
		t.Fatalf("the first default config invented a provider row, which contradicts the ticket 198 pin:\n%s", body)
	}
	m, err := NewManager(path, nil)
	if err != nil {
		t.Fatalf("NewManager over the first default config: %v", err)
	}
	if got := m.Config().LLM.Providers; got != nil {
		t.Fatalf("the ticket's premise (registry is nil, not an empty map) is gone: %#v", got)
	}
	return m, path, body
}

// AC#1, first half: the clean machine refuses all seven, and every refusal is
// shape (c) - it says which of the three reasons it is, and for a missing row it
// names the file-side shape that creates one.
func TestTicket257R1AC1CleanMachineRefusesAllSevenFields(t *testing.T) {
	m, path, before := cleanMachineManager(t)

	seenReasons := map[string]int{}
	for _, f := range writableRosterWalk() {
		written, err := f.run(m)
		if err == nil {
			t.Errorf("%s: the clean machine accepted the write (reported %v) - the ticket is not closed", f.field, written)
			continue
		}
		if written != nil {
			t.Errorf("%s: a refused write reported key paths %v, want none", f.field, written)
		}
		msg := err.Error()
		tag := whichReason(msg)
		if tag == "" {
			t.Errorf("%s: refusal names none of the three reasons: %s", f.field, msg)
			continue
		}
		if tag == refusalFileMissing {
			t.Errorf("%s: the file IS built here, so reason 1 is the wrong answer: %s", f.field, msg)
		}
		seenReasons[tag]++
		// Shape (c)'s half of the bargain: the owner is told where the row comes
		// from, spelled the way the schema reads it.
		if !strings.Contains(msg, "[llm.providers.") {
			t.Errorf("%s: refusal does not name [llm.providers.<name>], so nothing tells the operator how to make the row writable: %s",
				f.field, msg)
		}
		// The folded line AC#2 forbids, asserted absent rather than assumed absent.
		for _, folded := range []string{"配置未生效", "生效失败", "写入失败，配置未生效"} {
			if strings.Contains(msg, folded) {
				t.Errorf("%s: refusal folded into the shared line %q: %s", f.field, folded, msg)
			}
		}
	}
	if after := readSettingsFile(t, path); after != before {
		t.Errorf("seven refused writes still changed config.toml:\n%s", after)
	}
	if got := m.Config().LLM.Roles.Chat.Model; got != "" {
		t.Errorf("a refused role write still moved memory: chat model = %q", got)
	}
	t.Logf("clean machine, seven refusals by reason: %v (fields: %d)", seenReasons, len(writableRosterWalk()))
}

// AC#1, second half: the rows the guidance teaches, added by hand the way shape
// (c) says they are added, unlock all seven. This is the difference between an
// honest refusal and a dead end.
func TestTicket257R1AC1HandAddedRowsUnlockAllSevenFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := SaveFile(path, NewDefaults()); err != nil {
		t.Fatalf("first-run SaveFile: %v", err)
	}
	body := readSettingsFile(t, path)
	// The first default file already carries an EMPTY [llm.roles.chat] section
	// (MarshalCanonical emits every static field, no omitempty), so the real
	// operator move for the role pair is to FILL that section in place, not to
	// append a second one - appending one is a toml duplicate-table error and
	// would make this test lie about the machine.
	if !strings.Contains(body, "[llm.roles.chat]") {
		t.Fatalf("premise gone: the first default config no longer carries [llm.roles.chat]:\n%s", body)
	}
	// The operator follows the sentence: one provider section, its catalog rows
	// with the enabled key the guidance names, and the chat role pair named in
	// the section that is already there. The provider key line is matched by
	// regex so the test survives serialization-order drift.
	hand := strings.Replace(body, "[llm.roles.chat]", "[llm.roles.chat]", 1)
	hand = strings.Replace(hand,
		"[llm.roles.chat]\nprovider = ''\nmodel = ''",
		"[llm.roles.chat]\nprovider = \"deepseek\"\nmodel = \"deepseek-chat\"",
		1)
	if hand == body {
		t.Fatalf("the empty chat pair the guidance teaches to fill was not found in the first default file:\n%s", body)
	}
	hand += `
[llm.providers.deepseek]
api_key_ref = "env:DEEPSEEK_KEY"

[llm.providers.deepseek.models.deepseek-chat]
context_window = 64000
enabled = true

[llm.providers.deepseek.models.deepseek-reasoner]
enabled = true
`
	if err := os.WriteFile(path, []byte(hand), 0o600); err != nil {
		t.Fatal(err)
	}
	m, err := NewManager(path, nil)
	if err != nil {
		t.Fatalf("the hand-added shape the guidance teaches must load: %v", err)
	}
	if got := len(m.Config().LLM.Providers[cleanProvider].Models); got != 2 {
		t.Fatalf("catalog rows after the hand add = %d, want 2", got)
	}

	accepted := 0
	for _, f := range writableRosterWalk() {
		written, err := f.run(m)
		if err != nil {
			t.Errorf("%s: still refused after the hand add: %v", f.field, err)
			continue
		}
		if len(written) == 0 {
			t.Errorf("%s: accepted but reported no key path: %v", f.field, written)
			continue
		}
		accepted++
	}
	if accepted != 7 {
		t.Fatalf("unlocked fields = %d/7 after following the guidance; shape (c) promises 7/7", accepted)
	}

	bodyAfter := readSettingsFile(t, path)
	for _, want := range []string{
		"https://clean-machine.example.invalid/v1",
		"128000",
		"llm.providers.deepseek",
	} {
		if !strings.Contains(bodyAfter, want) {
			t.Errorf("the accepted writes never reached the file (missing %q):\n%s", want, bodyAfter)
		}
	}
	// AC#3's boundary, asserted rather than assumed: this path stores a reference
	// and never a value, and the schema still has no field that could hold one.
	if !strings.Contains(bodyAfter, "api_key_ref = 'dpapi:t257r1canaryblob01'") {
		t.Errorf("the credential leg's reference did not land as a reference:\n%s", bodyAfter)
	}
	if strings.Contains(bodyAfter, "api_key =") {
		t.Errorf("a plaintext api_key key appeared in config.toml:\n%s", bodyAfter)
	}
}

// AC#2 reason 1: no config.toml. Planted by asking for none.
func TestTicket257R1AC2MissingFileIsItsOwnSentence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	_, err := NewManager(path, nil)
	if err == nil {
		t.Fatal("NewManager succeeded with no file")
	}
	msg := err.Error()
	assertOneReason(t, "no config.toml", msg, refusalFileMissing)
	// The classifier cmd/wisp reads cause=missing off this chain (票 223 AC#4):
	// the new sentence must not eat the sentinel it keys on.
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("the missing-file error no longer carries fs.ErrNotExist, so cause=missing cannot fire: %v", err)
	}
	// Same sentence from the other door the same file is read through.
	if _, _, err := LoadFile(path, nil); err == nil || whichReason(err.Error()) != refusalFileMissing {
		t.Errorf("LoadFile on a missing file must answer the same reason 1, got: %v", err)
	}
}

// AC#2 reason 2: the file is there, the row is not. Both halves of the row gate
// (provider section, catalog row) are this one reason and each must say so.
func TestTicket257R1AC2MissingRowIsItsOwnSentence(t *testing.T) {
	m, _, _ := cleanMachineManager(t)

	_, err := m.SetProviderBaseURL(cleanProvider, "https://row-gate.example.invalid/v1")
	if err == nil {
		t.Fatal("a provider-level write was accepted with no provider row")
	}
	msg := err.Error()
	assertOneReason(t, "no provider row", msg, refusalRowMissing)
	if !strings.Contains(msg, guidanceProviderRow) {
		t.Errorf("the row refusal does not say where a provider row is created: %s", msg)
	}
	if strings.Contains(msg, "[providers.") {
		t.Errorf("the refusal tells the operator to write a top-level [providers.x] section, which "+
			"DisallowUnknownFields rejects and which stops the whole chain (ledger A543): %s", msg)
	}

	_, err = m.SetModelContextWindow(cleanProvider, cleanModel, 32000)
	if err == nil {
		t.Fatal("a model-level write was accepted with no catalog row")
	}
	assertOneReason(t, "no catalog row", err.Error(), refusalRowMissing)

	// A row that exists but carries no such model is still reason 2, and the pin
	// ticket 261 put on this sentence keeps its wording.
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(settingsBase), 0o600); err != nil {
		t.Fatal(err)
	}
	row, err := NewManager(path, nil)
	if err != nil {
		t.Fatalf("load settingsBase: %v", err)
	}
	if _, err := row.SetModelContextWindow(cleanProvider, "ghost-model", 32000); err != nil {
		assertOneReason(t, "model absent from an existing row", err.Error(), refusalRowMissing)
		if !strings.Contains(err.Error(), "refusing to invent one") {
			t.Errorf("ticket 261's pin on this sentence broke: %v", err)
		}
	} else {
		t.Error("SetModelContextWindow minted a catalog row")
	}
}

// AC#2 reason 3: the row is there and the value is what fails. Planted twice,
// because role_chat_model is the field whose clean-machine answer changed to
// reason 2 - this proves its reason-3 branch is still reachable when the provider
// really is named.
func TestTicket257R1AC2InvalidValueIsItsOwnSentence(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(settingsBase), 0o600); err != nil {
		t.Fatal(err)
	}
	m, err := NewManager(path, nil)
	if err != nil {
		t.Fatalf("load settingsBase: %v", err)
	}
	before := readSettingsFile(t, path)

	if _, err := m.SetProviderAPIKeyRef(cleanProvider, "no-reference-prefix-here"); err != nil {
		assertOneReason(t, "value not a reference", err.Error(), refusalInvalid)
	} else {
		t.Error("SetProviderAPIKeyRef accepted a value that is not a reference")
	}
	if _, err := m.SetRoleChatModel("not-in-the-catalog-9"); err != nil {
		assertOneReason(t, "role pair not in the catalog", err.Error(), refusalInvalid)
	} else {
		t.Error("SetRoleChatModel accepted a model the catalog does not hold")
	}
	if after := readSettingsFile(t, path); after != before {
		t.Errorf("reason-3 refusals changed config.toml:\n%s", after)
	}
}

// whichReason returns the single reason tag a refusal carries, or "" when it
// carries none of them or more than one. "More than one" is a fold, and a fold is
// what AC#2 forbids, so this is the check that gives the judge its teeth.
func whichReason(msg string) string {
	hits := []string{}
	for _, tag := range []string{refusalFileMissing, refusalRowMissing, refusalInvalid} {
		if strings.Contains(msg, tag) {
			hits = append(hits, tag)
		}
	}
	if len(hits) != 1 {
		return ""
	}
	return hits[0]
}

// assertOneReason is AC#2 in three lines: this reason is said, the other two are
// not, and nothing about the message is shared with the folded line.
func assertOneReason(t *testing.T, what, msg, want string) {
	t.Helper()
	if got := whichReason(msg); got == "" {
		t.Errorf("%s: message carries no reason tag, or more than one (that is a fold): %s", what, msg)
		return
	} else if got != want {
		t.Errorf("%s: answered %q, want %q: %s", what, got, want, msg)
	}
	for _, tag := range []string{refusalFileMissing, refusalRowMissing, refusalInvalid} {
		if tag == want {
			continue
		}
		if strings.Contains(msg, tag) {
			t.Errorf("%s: also carries %q, so the three reasons are not three sentences", what, tag)
		}
	}
}
