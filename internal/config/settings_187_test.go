package config

// Ticket 187's config-side tests: the thinking-tier write leg and the picker
// read face (settings_187.go).
//
// Every negative case carries its positive control in the same test, and every
// accepted write is read back the way the next start would read it (the same
// rule settings_248_test.go states). The mutation teeth for this half:
//
//	- drop the declared-levels check from SetRoleChatThinking ->
//	  Test187ThinkingWriteRefusesALevelTheModelDoesNotDeclare reds;
//	- drop the enum check -> the same test's "extreme" case reds;
//	- drop the Enabled filter or the Reason sentences in ChatSwitchOptions ->
//	  Test187ChatSwitchOptionsNeverOffersWhatTheResolverWouldRefuse reds.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// settings187Base is a config with one provider, one enabled model that
// declares off|low, and one DISABLED model, so the picker's filter and the
// declared-levels check both have something real to bite.
const settings187Base = `schema_version = 2

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

[llm.providers.deepseek.models.deepseek-reasoner]
enabled = false
context_window = 128000
thinking_levels = ["off", "low", "medium", "high"]
`

// AC#1's write half: the tier write lands on the file, the receipt names the
// key path, and the value is the one the next task's re-resolution would read.
func Test187ThinkingWriteLandsTheKeyAndNamesIt(t *testing.T) {
	m, path, _ := newTestManager(t, settings187Base)

	written, err := m.SetRoleChatThinking("low")
	if err != nil {
		t.Fatalf("SetRoleChatThinking: %v", err)
	}
	if len(written) != 1 || written[0] != "llm.roles.chat.thinking_intensity" {
		t.Errorf("written key paths = %v, want [llm.roles.chat.thinking_intensity]", written)
	}
	got := readLoaded(t, path)
	if got.LLM.Roles.Chat.ThinkingIntensity != "low" {
		t.Errorf("the tier never reached the file: %+v", got.LLM.Roles.Chat)
	}
	// The in-memory config the per-task re-resolution reads moved with the
	// file - this is the fact that makes the tier next_task instead of restart.
	if m.Config().LLM.Roles.Chat.ThinkingIntensity != "low" {
		t.Errorf("memory did not move with the file: %q", m.Config().LLM.Roles.Chat.ThinkingIntensity)
	}
}

// The enum and the declared-levels refusals, each with the file untouched
// afterwards (reason 3: the value would not survive validation).
func Test187ThinkingWriteRefusesALevelTheModelDoesNotDeclare(t *testing.T) {
	m, path, _ := newTestManager(t, settings187Base)
	before := readFile187(t, path)

	// 1. Outside the four-value vocabulary entirely.
	if _, err := m.SetRoleChatThinking("extreme"); err == nil {
		t.Fatal("a tier outside the enum was accepted")
	}
	// 2. Inside the enum but not one the model declares (deepseek-chat says off|low).
	if _, err := m.SetRoleChatThinking("high"); err == nil {
		t.Fatal("a tier the model does not declare was accepted")
	} else if !strings.Contains(err.Error(), "declares thinking_levels") {
		t.Errorf("the refusal does not name the model's own list: %v", err)
	}
	// 3. Empty is not a tier (an absent key is the schema's business, not a
	// panel write's).
	if _, err := m.SetRoleChatThinking("   "); err == nil {
		t.Fatal("an empty tier was accepted")
	}
	if readFile187(t, path) != before {
		t.Errorf("a refused write moved the file:\n%s", readFile187(t, path))
	}

	// Control: a tier inside BOTH the enum and the model's declared list lands.
	if _, err := m.SetRoleChatThinking("low"); err != nil {
		t.Fatalf("the control write was refused too: %v", err)
	}
}

// A model that declares NO thinking_levels constrains nothing: the four-value
// enum alone decides. The write must not invent a stricter rule than the
// schema carries (and the picker side lists nothing for such a model - the
// asymmetry is measured, stated in the evidence doc, and pinned here).
func Test187ThinkingWriteOnAModelThatDeclaresNoLevelsIsEnumOnly(t *testing.T) {
	base := strings.Replace(settings187Base,
		`thinking_levels = ["off", "low"]

[llm.providers.deepseek.models.deepseek-chat.price]`,
		`[llm.providers.deepseek.models.deepseek-chat.price]`, 1)
	m, _, _ := newTestManager(t, base)

	if _, err := m.SetRoleChatThinking("medium"); err != nil {
		t.Errorf("an enum-legal tier was refused on a model that declares no levels: %v", err)
	}
	if _, err := m.SetRoleChatThinking("extreme"); err == nil {
		t.Error("a tier outside the enum was accepted on a model that declares no levels")
	}
}

// The read face: only enabled rows are offered (the resolver refuses disabled
// ones - offering them would be offering a request that cannot resolve), the
// declared levels come from the CURRENT model, and the empty shapes say so
// instead of rendering as a blank menu (AC#4's 宁缺毋造).
func Test187ChatSwitchOptionsNeverOffersWhatTheResolverWouldRefuse(t *testing.T) {
	m, _, _ := newTestManager(t, settings187Base)

	opts := m.ChatSwitchOptions()
	if opts.Reason != "" {
		t.Errorf("a readable catalog with offerable rows carried a Reason: %q", opts.Reason)
	}
	if opts.CurrentProvider != "deepseek" || opts.CurrentModel != "deepseek-chat" {
		t.Errorf("current pair = %s/%s, want deepseek/deepseek-chat", opts.CurrentProvider, opts.CurrentModel)
	}
	if opts.CurrentThinking != "off" {
		t.Errorf("current thinking = %q, want the file's \"off\"", opts.CurrentThinking)
	}
	if len(opts.Models) != 1 || opts.Models[0].Provider != "deepseek" || opts.Models[0].Model != "deepseek-chat" {
		t.Errorf("offered rows = %+v, want exactly the enabled deepseek-chat (the disabled deepseek-reasoner must not be offered)", opts.Models)
	}
	if strings.Join(opts.ThinkingLevels, ",") != "off,low" {
		t.Errorf("thinkingLevels = %v, want the CURRENT model's declared off|low", opts.ThinkingLevels)
	}
}

// The two empty shapes, each with its own sentence: no providers at all, and
// providers whose rows are all disabled. An empty Models with no Reason would
// render as "nothing to choose" over a truth the user can fix - the same
// lying-key family the commands section's Reason exists for.
func Test187ChatSwitchOptionsEmptyShapesSaySo(t *testing.T) {
	noProviders, _, _ := newTestManager(t, "schema_version = 2\n")
	opts := noProviders.ChatSwitchOptions()
	if len(opts.Models) != 0 || opts.Reason == "" || !strings.Contains(opts.Reason, "没有可用模型") {
		t.Errorf("no-providers shape = %+v, want empty rows and the explicit sentence", opts)
	}

	allDisabled := strings.Replace(settings187Base,
		"[llm.providers.deepseek.models.deepseek-chat]\nenabled = true\ncontext_window = 64000",
		"[llm.providers.deepseek.models.deepseek-chat]\nenabled = false\ncontext_window = 64000", 1)
	disabled, _, _ := newTestManager(t, allDisabled)
	opts = disabled.ChatSwitchOptions()
	if len(opts.Models) != 0 || opts.Reason == "" || !strings.Contains(opts.Reason, "enabled=false") {
		t.Errorf("all-disabled shape = %+v, want empty rows and the re-enable sentence", opts)
	}
	// The declared levels still name the current model even when nothing is
	// offerable: the pair is readable, only the offer is empty.
	if strings.Join(opts.ThinkingLevels, ",") != "off,low" {
		t.Errorf("thinkingLevels = %v, want the current model's declared list even in the empty-offer shape", opts.ThinkingLevels)
	}
}

// An unset chat role reads as the unset state - empty pair, empty levels -
// and does not panic hunting for a model that is not named.
func Test187ChatSwitchOptionsUnsetRoleReadsAsUnset(t *testing.T) {
	base := strings.Replace(settings187Base,
		"provider = \"deepseek\"\nmodel = \"deepseek-chat\"\nthinking_intensity = \"off\"",
		"", 1)
	m, _, _ := newTestManager(t, base)

	opts := m.ChatSwitchOptions()
	// The schema default for the role tier is applied at load, so an unset
	// role still reads "off" - the file's own default, not an invention here.
	if opts.CurrentProvider != "" || opts.CurrentModel != "" || opts.CurrentThinking != "off" {
		t.Errorf("unset role read as %+v, want the empty pair plus the schema default tier", opts)
	}
	if len(opts.Models) != 1 {
		t.Errorf("the catalog is still offerable without a role: %+v", opts.Models)
	}
	if len(opts.ThinkingLevels) != 0 {
		t.Errorf("no current model means no declared levels: %v", opts.ThinkingLevels)
	}
}

// AC#3's config-side nail, read off the bytes: this file's write leg carries
// no permission surface of any kind. The thinking tier and the risk mode are
// unrelated dimensions, and the nail keeps the second one from ever arriving
// through the first's door.
func Test187ConfigWriteLegTouchesNoPermissionSurface(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("settings_187.go"))
	if err != nil {
		t.Fatalf("read settings_187.go: %v", err)
	}
	src := string(data)
	for _, banned := range []string{"internal/perm", "internal/risk", "risk.Mode", "Modes.Set", "permission_mode"} {
		if strings.Contains(src, banned) {
			t.Errorf("settings_187.go mentions %q - the tier write leg must stay unrelated to the permission dimension (AC#3)", banned)
		}
	}
}

func readFile187(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
