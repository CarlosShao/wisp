package panel

// Ticket 145 AC#2b's two half-items at the carrier level: the wire honesty of
// the fields this ticket added to ModelOption and ModelSwitchSection. The
// VALUE teeth (display names off the catalog, tier words off the adapter's
// own mapping table) live where their sources live - internal/config's
// settings test and internal/llm's vocabulary test - and the assembly
// mapping is pinned in cmd/wisp. This file pins what only the carrier owns:
// a populated field travels on the bytes under its wire name, an empty one
// stays absent (omitempty is the last mile, and "no display name declared"
// / "this adapter maps no tier" must not render as invented values).

import (
	"encoding/json"
	"strings"
	"testing"
)

func Test145ModelSwitchNewFieldsTravelOnTheWire(t *testing.T) {
	sec := ModelSwitchSection{
		Current:               ModelOption{Provider: "deepseek", Model: "deepseek-chat", Display: "DeepSeek Chat"},
		CurrentThinking:       "off",
		Models:                []ModelOption{{Provider: "deepseek", Model: "deepseek-chat", Display: "DeepSeek Chat"}, {Provider: "openai", Model: "gpt-x"}},
		ThinkingLevels:        []string{"off", "low"},
		Adapter:               "openai-responses",
		AdapterThinkingLevels: []string{"high", "low", "medium"},
	}
	data, err := json.Marshal(sec)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`"display":"DeepSeek Chat"`, `"adapter":"openai-responses"`,
		`"adapterThinkingLevels":["high","low","medium"]`,
	} {
		if !strings.Contains(string(data), want) {
			t.Errorf("the wire lost %s:\n%s", want, data)
		}
	}

	// An entry with no declared display name and an adapter that maps no
	// tier travel absent - the honest shapes, not empty strings dressed up.
	quiet := ModelSwitchSection{
		Current:         ModelOption{Provider: "deepseek", Model: "deepseek-chat"},
		CurrentThinking: "off",
		Models:          []ModelOption{{Provider: "deepseek", Model: "deepseek-chat"}},
	}
	quietData, err := json.Marshal(quiet)
	if err != nil {
		t.Fatal(err)
	}
	for _, banned := range []string{`"display"`, `"adapter"`, `"adapterThinkingLevels"`} {
		if strings.Contains(string(quietData), banned) {
			t.Errorf("an undeclared %s must stay absent, never invented:\n%s", banned, quietData)
		}
	}
}
