package config

// Ticket 145 AC#2b①'s config-side value teeth: ChatSwitchOptions now carries
// each catalog entry's own human-facing Display name (ModelSpec.Display,
// verbatim) and the CURRENT model's display, off the same walk that fills
// Models and ThinkingLevels. The assertions are hardcoded from the config
// text below, so a reader that invents a name - or drops the real one - reds
// here instead of on the page.

import (
	"strings"
	"testing"
)

const settings145DisplayBase = `schema_version = 2

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
display = "DeepSeek Chat"
thinking_levels = ["off", "low"]

[llm.providers.deepseek.models.deepseek-chat.capabilities]
thinking = true

[llm.providers.deepseek.models.deepseek-chat.price]
in = 2500000
out = 10000000

[llm.providers.deepseek.models.deepseek-reasoner]
enabled = false
context_window = 128000
display = "DeepSeek Reasoner"
thinking_levels = ["off", "low", "medium", "high"]

[llm.providers.deepseek.models.deepseek-reasoner.capabilities]
thinking = true
`

func Test145ChatSwitchOptionsCarriesTheCatalogDisplayNames(t *testing.T) {
	m, _, _ := newTestManager(t, settings145DisplayBase)

	opts := m.ChatSwitchOptions()
	if len(opts.Models) != 1 {
		t.Fatalf("offered rows = %+v, want exactly the enabled deepseek-chat", opts.Models)
	}
	if opts.Models[0].Display != "DeepSeek Chat" {
		t.Errorf("display = %q, want the catalog entry's own \"DeepSeek Chat\" verbatim", opts.Models[0].Display)
	}
	if opts.CurrentDisplay != "DeepSeek Chat" {
		t.Errorf("currentDisplay = %q, want the current model's own display off the same walk", opts.CurrentDisplay)
	}
}

// A model that declares no display name travels with none - the page falls
// back to the model id, and this reader invents no name. The disabled row's
// display (deepseek-reasoner) must not leak onto the enabled row's answer.
func Test145ChatSwitchOptionsUndeclaredDisplayStaysEmpty(t *testing.T) {
	base := strings.Replace(settings145DisplayBase, `display = "DeepSeek Chat"`+"\n", "", 1)
	m, _, _ := newTestManager(t, base)

	opts := m.ChatSwitchOptions()
	if len(opts.Models) != 1 || opts.Models[0].Display != "" {
		t.Errorf("offered rows = %+v, want the enabled row with no invented display", opts.Models)
	}
	if opts.CurrentDisplay != "" {
		t.Errorf("currentDisplay = %q, want empty when the entry declares none", opts.CurrentDisplay)
	}
}

// settings145TwoModelBase is AC#2b's own criterion fixture: a TWO-MODEL tree,
// each entry naming TWO tiers, both enabled. The existing two-model fixtures
// all leave the second entry enabled=false and assert len==1, which measures
// the FILTER direction; this one measures the PLURAL direction - the list has
// to come back with both rows, in sorted order, word for word.
//
// capabilities.thinking = true is declared on both entries because AC#2b's
// existence question is answered by that bit (see thinkingDimensionDeclared);
// the no-bit and no-tiers variants below derive from this text rather than
// duplicating it.
const settings145TwoModelBase = `schema_version = 2

[llm]
text_chain = ["acme/first", "acme/second"]

[llm.roles.chat]
provider = "acme"
model = "first"
thinking_intensity = "low"

[llm.providers.acme]
protocol = "openai-responses"
api_key_ref = "env:ACME_KEY"

[llm.providers.acme.models.first]
enabled = true
context_window = 128000
display = "Acme First"
thinking_levels = ["low", "medium"]

[llm.providers.acme.models.first.capabilities]
thinking = true

[llm.providers.acme.models.second]
enabled = true
context_window = 64000
display = "Acme Second"
thinking_levels = ["low", "high"]

[llm.providers.acme.models.second.capabilities]
thinking = true
`

// Test145ChatSwitchOptionsReadsATwoModelTreeWordForWord is AC#2b's stated
// criterion: feed a two-model / two-tier fake config tree in, and the list the
// snapshot reads back must match that tree verbatim - not "the field is
// non-empty", which the AC forbids as a criterion.
//
// The mutation this nail is built to catch is "offer only the first row"
// (settings_187.go: append `if len(out.Models) >= 1 { return out }` inside the
// walk): every other model-switch test in the repo asserts len==1 against a
// one-enabled-model tree, so that mutation left the whole repo green before
// this test existed. Measured red/green is recorded in
// docs/evidence/s1/145-ac2b-r1.md section 4.
func Test145ChatSwitchOptionsReadsATwoModelTreeWordForWord(t *testing.T) {
	m, _, _ := newTestManager(t, settings145TwoModelBase)

	opts := m.ChatSwitchOptions()
	if opts.Reason != "" {
		t.Errorf("a readable catalog with offerable rows carried a Reason: %q", opts.Reason)
	}

	// Word for word, in the order the reader promises (provider then model):
	// both rows, their own display names, nothing invented, nothing dropped.
	wantRows := []ChatModelOption{
		{Provider: "acme", Model: "first", Display: "Acme First"},
		{Provider: "acme", Model: "second", Display: "Acme Second"},
	}
	if len(opts.Models) != len(wantRows) {
		t.Fatalf("offered rows = %+v, want exactly the tree's %d enabled models", opts.Models, len(wantRows))
	}
	for i, want := range wantRows {
		if opts.Models[i] != want {
			t.Errorf("row %d = %+v, want %+v verbatim from the tree", i, opts.Models[i], want)
		}
	}

	// The current pair and its display ride the same walk.
	if opts.CurrentProvider != "acme" || opts.CurrentModel != "first" {
		t.Errorf("current pair = %s/%s, want the tree's acme/first", opts.CurrentProvider, opts.CurrentModel)
	}
	if opts.CurrentDisplay != "Acme First" {
		t.Errorf("currentDisplay = %q, want the current entry's own name", opts.CurrentDisplay)
	}
	if opts.CurrentThinking != "low" {
		t.Errorf("currentThinking = %q, want the file's \"low\"", opts.CurrentThinking)
	}

	// The current model's own two tiers, verbatim - not the four-value enum.
	if got := strings.Join(opts.ThinkingLevels, ","); got != "low,medium" {
		t.Errorf("thinkingLevels = %q, want the current entry's own \"low,medium\" (its two, not the enum's four)", got)
	}

	// The OTHER model's tiers must not leak onto the current model's answer:
	// second declares low|high, and the read face answers about first only.
	for _, lvl := range opts.ThinkingLevels {
		if lvl == "high" {
			t.Errorf("thinkingLevels = %v carries the SECOND model's tier - the reader answered about the wrong entry", opts.ThinkingLevels)
		}
	}
}

// Test145ChatSwitchOptionsThinkingBitDecidesTheDimension is ticket 145 AC#2b
// defect (1)'s nail: the existence of the thinking dimension is decided by the
// declaration vocabulary AND the capability bit, and the bit VETOES. Three
// shapes off one tree, so the only thing that moves between them is the bit:
//
//	(a) bit true  + two tiers  -> the two tiers;
//	(b) bit false + two tiers  -> NOTHING (levels declared under a false bit
//	    are dropped, never trusted);
//	(c) bit true  + no tiers   -> NOTHING (the entry names no tier, and the
//	    reader does not substitute the four-value enum).
//
// In all three the rows and the current display stay put: the bit governs the
// TIER dimension only, never the model list.
func Test145ChatSwitchOptionsThinkingBitDecidesTheDimension(t *testing.T) {
	// (a) the declared shape - the control, so a reader that always answers
	// empty cannot pass this test by being uniformly wrong.
	with, _, _ := newTestManager(t, settings145TwoModelBase)
	if got := strings.Join(with.ChatSwitchOptions().ThinkingLevels, ","); got != "low,medium" {
		t.Fatalf("control (bit declared) thinkingLevels = %q, want the entry's own low,medium", got)
	}

	// (b) the bit flipped to false, the two tiers left in place. This is the
	// measured defect: before AC#2b the reader sent ["low","medium"] here.
	noBit := strings.Replace(settings145TwoModelBase,
		`[llm.providers.acme.models.first.capabilities]
thinking = true`,
		`[llm.providers.acme.models.first.capabilities]
thinking = false`, 1)
	m, _, _ := newTestManager(t, noBit)
	opts := m.ChatSwitchOptions()
	if len(opts.ThinkingLevels) != 0 {
		t.Errorf("thinkingLevels = %v, want EMPTY: the entry declares two tiers but never declares the capability, and the bit vetoes", opts.ThinkingLevels)
	}
	// The veto is scoped to the tier dimension: the catalog is still offered.
	if len(opts.Models) != 2 || opts.CurrentDisplay != "Acme First" {
		t.Errorf("rows = %+v / currentDisplay = %q, want the catalog untouched by the tier veto",
			opts.Models, opts.CurrentDisplay)
	}

	// (c) the bit declared, the tiers absent - "thinking exists, and I name no
	// tier for it" must offer nothing rather than the four-value enum.
	noTiers := strings.Replace(settings145TwoModelBase,
		`thinking_levels = ["low", "medium"]`+"\n", "", 1)
	m, _, _ = newTestManager(t, noTiers)
	opts = m.ChatSwitchOptions()
	if len(opts.ThinkingLevels) != 0 {
		t.Errorf("thinkingLevels = %v, want EMPTY: the entry names no tier and the reader invents none", opts.ThinkingLevels)
	}
	if len(opts.Models) != 2 {
		t.Errorf("rows = %+v, want both entries still offered", opts.Models)
	}
}
