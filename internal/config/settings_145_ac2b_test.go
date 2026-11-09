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

[llm.providers.deepseek.models.deepseek-chat.price]
in = 2500000
out = 10000000

[llm.providers.deepseek.models.deepseek-reasoner]
enabled = false
context_window = 128000
display = "DeepSeek Reasoner"
thinking_levels = ["off", "low", "medium", "high"]
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
