package config

// R-145-3's write-face teeth (ticket 416): SetRoleChatThinking does NOT consult
// Capabilities.Thinking, and that is now a DECIDED ruling, not an oversight -
// the reasoning lives in the setter's own comment (settings_187.go). These
// tests pin the ruling's observable shape:
//
//   - a model that never declared the bit still takes a tier write, and the
//     round trip keeps what was written and drops only the OFFER (the read
//     face's bit gate) - write what you read back, and name what is lost;
//   - the write face's refusals stay VALUE-shaped: they track the entry's own
//     declared vocabulary, never its capability bits;
//   - the tier rides across a model switch, and re-asserting it on a no-bit
//     model lands - the preference the read face hides is not dead weight.
//
// The mutation this file exists to catch: give requireDeclaredThinkingLevel a
// capability-bit conjunct (or delete the declared-list judgment) and
// Test187R1453WriteFaceRefusalsTrackTheVocabularyNotTheBit reds; delete the
// whole declared-list check and the same test plus
// Test187ThinkingWriteRefusesALevelTheModelDoesNotDeclare red.

import (
	"strings"
	"testing"
)

// settings187R1453Base: one provider, two enabled models. "probed" is the
// catalog shape the read face offers from - bit declared, two tiers. "imported"
// is the shape ImportDiscovered leaves behind when the probe never ran
// (internal/llm/discover.go imports entries with every capability bit unset,
// SPEC-03 section 3.1's all-false default): levels declared, NO capabilities
// section at all. The whole R-145-3 question lives in that asymmetry.
const settings187R1453Base = `schema_version = 2

[llm]
text_chain = ["acme/probed"]

[llm.roles.chat]
provider = "acme"
model = "probed"
thinking_intensity = "off"

[llm.providers.acme]
protocol = "openai-responses"
api_key_ref = "env:ACME_KEY"

[llm.providers.acme.models.probed]
enabled = true
context_window = 128000
thinking_levels = ["off", "low"]

[llm.providers.acme.models.probed.capabilities]
thinking = true

[llm.providers.acme.models.imported]
enabled = true
context_window = 32000
thinking_levels = ["off", "low"]
`

// Test187R1453NoBitModelStillTakesADeclaredTierWrite is the ruling's round
// trip: a tier write on a model that never declared the capability bit lands,
// the file and memory both hold exactly what was written, and the read back
// names what the trip drops - the OFFER, not the value. Both no-bit shapes are
// measured: the bit absent entirely (the imported entry) and the bit explicitly
// false (an entry the probe ran on and answered no for) - the write face
// treats them identically because it never asks the question.
func Test187R1453NoBitModelStillTakesADeclaredTierWrite(t *testing.T) {
	// Shape 1: the bit absent entirely.
	m, path, _ := newTestManager(t, settings187R1453Base)
	if _, err := m.SetRoleChatModel("imported"); err != nil {
		t.Fatalf("switch to the imported entry: %v", err)
	}

	written, err := m.SetRoleChatThinking("low")
	if err != nil {
		t.Fatalf("a declared tier was refused on a model that never declared the bit: %v", err)
	}
	if len(written) != 1 || written[0] != "llm.roles.chat.thinking_intensity" {
		t.Errorf("written key paths = %v, want [llm.roles.chat.thinking_intensity]", written)
	}
	if got := readLoaded(t, path); got.LLM.Roles.Chat.ThinkingIntensity != "low" {
		t.Errorf("the tier never reached the file: %+v", got.LLM.Roles.Chat)
	}
	if m.Config().LLM.Roles.Chat.ThinkingIntensity != "low" {
		t.Errorf("memory did not move with the file: %q", m.Config().LLM.Roles.Chat.ThinkingIntensity)
	}

	// The round trip's drop, named: the stored VALUE survives verbatim
	// (CurrentThinking reads back "low"); the read face still answers the
	// existence question with the bit and offers NOTHING for this model.
	// Both halves are the ruling - neither face may answer the other's
	// question, and this is what that looks like on one trip.
	opts := m.ChatSwitchOptions()
	if opts.CurrentThinking != "low" {
		t.Errorf("currentThinking = %q, want the stored tier verbatim (write what you read back)", opts.CurrentThinking)
	}
	if len(opts.ThinkingLevels) != 0 {
		t.Errorf("thinkingLevels = %v, want EMPTY: the read face's bit gate still hides the dimension the write face stored", opts.ThinkingLevels)
	}

	// Shape 2: the bit explicitly false. Same write, same answer - the setter
	// does not ask, so false and absent cannot diverge here.
	explicitFalse := strings.Replace(settings187R1453Base,
		`[llm.providers.acme.models.imported]
enabled = true
context_window = 32000
thinking_levels = ["off", "low"]
`,
		`[llm.providers.acme.models.imported]
enabled = true
context_window = 32000
thinking_levels = ["off", "low"]

[llm.providers.acme.models.imported.capabilities]
thinking = false
`, 1)
	m2, _, _ := newTestManager(t, explicitFalse)
	if _, err := m2.SetRoleChatModel("imported"); err != nil {
		t.Fatalf("switch to the bit-false entry: %v", err)
	}
	if _, err := m2.SetRoleChatThinking("low"); err != nil {
		t.Errorf("a declared tier was refused on a model whose probe answered thinking=false: %v", err)
	}
}

// Test187R1453WriteFaceRefusalsTrackTheVocabularyNotTheBit pins the half the
// ruling keeps: the declared-list check is a VALUE judgment. On the no-bit
// entry the four-value enum and the entry's own two-tier list BOTH apply, in
// that order, and neither refusal mentions or depends on capability bits.
func Test187R1453WriteFaceRefusalsTrackTheVocabularyNotTheBit(t *testing.T) {
	m, path, _ := newTestManager(t, settings187R1453Base)
	if _, err := m.SetRoleChatModel("imported"); err != nil {
		t.Fatalf("switch to the imported entry: %v", err)
	}
	before := readFile187(t, path)

	// Enum-legal but not one the entry names: refused, by the entry's own
	// list - with no bit anywhere in the argument.
	if _, err := m.SetRoleChatThinking("medium"); err == nil {
		t.Fatal("a tier outside the no-bit entry's declared list was accepted")
	} else if !strings.Contains(err.Error(), "declares thinking_levels") {
		t.Errorf("the refusal does not name the entry's own list (the only authority a value refusal has): %v", err)
	}
	// Outside the enum entirely: refused, by the schema's vocabulary.
	if _, err := m.SetRoleChatThinking("extreme"); err == nil {
		t.Fatal("a tier outside the enum was accepted on the no-bit entry")
	}
	if readFile187(t, path) != before {
		t.Errorf("a refused write moved the file:\n%s", readFile187(t, path))
	}

	// Control: the tier the entry DOES name lands, on the same no-bit entry.
	if _, err := m.SetRoleChatThinking("low"); err != nil {
		t.Errorf("the control write was refused too: %v", err)
	}
}

// Test187R1453TierPreferenceSurvivesAModelSwitch pins the cost side of the
// ruling: the model switcher never touches thinking_intensity, so a tier rides
// across models - and the re-assert write on the no-bit destination must land
// too. A bit veto at the setter would strand exactly this flow: the stored
// value stays, but the user could never write anything back, not even "off".
func Test187R1453TierPreferenceSurvivesAModelSwitch(t *testing.T) {
	m, path, _ := newTestManager(t, settings187R1453Base)

	if _, err := m.SetRoleChatThinking("low"); err != nil {
		t.Fatalf("tier write on the bit-declared entry: %v", err)
	}
	if _, err := m.SetRoleChatModel("imported"); err != nil {
		t.Fatalf("switch to the no-bit entry: %v", err)
	}
	if got := readLoaded(t, path).LLM.Roles.Chat.ThinkingIntensity; got != "low" {
		t.Errorf("the model switch disturbed the stored tier: %q", got)
	}
	// Re-assert on the no-bit destination: lands, because the setter asks the
	// vocabulary, not the bit.
	if _, err := m.SetRoleChatThinking("low"); err != nil {
		t.Errorf("re-asserting the tier on a no-bit model was refused: %v", err)
	}
	if opts := m.ChatSwitchOptions(); opts.CurrentThinking != "low" {
		t.Errorf("currentThinking = %q, want the riding tier verbatim after the switch", opts.CurrentThinking)
	}
}
