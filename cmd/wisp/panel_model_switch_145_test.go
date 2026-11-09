package main

// Ticket 145 AC#2b's assembly-root teeth: the modelSwitchSection mapping
// carries each catalog row's own Display name (AC#2b①, read off
// ChatSwitchOptions' walk) and the SENDABLE tier set - the model's declared
// tiers intersected with the live adapter's own census (AC#2b②, the
// orchestrator's 方案 B: Go computes the intersection once, so every reader
// receives a list it can send as-is instead of two lists to combine itself).
// The values are hardcoded from the adapters' own tables and the config text
// below, so a mapping that invents a name, drops a real one, hands the
// config's four-value enum out as "what every adapter supports", or ships the
// declaration un-intersected reds here.

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/panel"
)

func Test145ModelSwitchMappingCarriesTheDisplayNames(t *testing.T) {
	// One display-bearing entry, one without, the current pair on the first.
	base := strings.Replace(thinking187Base,
		`[llm.providers.deepseek.models.deepseek-chat]
enabled = true
context_window = 64000`,
		`[llm.providers.deepseek.models.deepseek-chat]
enabled = true
context_window = 64000
display = "DeepSeek Chat"`, 1)

	dir := t.TempDir()
	path := filepath.Join(dir, configFileName)
	if err := os.WriteFile(path, []byte(base), 0o600); err != nil {
		t.Fatal(err)
	}
	mgr, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("manager: %v", err)
	}

	sec := modelSwitchSection(mgr, nil)
	if len(sec.Models) != 1 || sec.Models[0].Display != "DeepSeek Chat" {
		t.Errorf("models = %+v, want the enabled row carrying its own display verbatim", sec.Models)
	}
	if sec.Current.Display != "DeepSeek Chat" {
		t.Errorf("current.display = %q, want the current model's own display off the same walk", sec.Current.Display)
	}
	// No provider built: no adapter word, no vocabulary - the honest unset
	// shape, never a guessed protocol.
	if sec.Adapter != "" || len(sec.AdapterThinkingLevels) != 0 {
		t.Errorf("adapter = %q, vocabulary = %v, want both empty with no provider built", sec.Adapter, sec.AdapterThinkingLevels)
	}

	// The wire carries the name the page will render, byte-level.
	data, err := panel.NewSnapshotPump(panel.PumpSources{
		ModelSwitch: func() panel.ModelSwitchSection { return modelSwitchSection(mgr, nil) },
	}).Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(data), `"display":"DeepSeek Chat"`) {
		t.Errorf("the wire lost the catalog's display name:\n%s", data)
	}
}

func Test145ModelSwitchMappingCarriesTheAdapterVocabulary(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, configFileName)
	if err := os.WriteFile(path, []byte(thinking187Base), 0o600); err != nil {
		t.Fatal(err)
	}
	mgr, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("manager: %v", err)
	}

	// A real provider of the responses protocol: the vocabulary is that
	// adapter's own registered words (effortLevels), sorted - not the
	// config enum's four values, and not "off".
	prov, err := llm.NewProvider(llm.EndpointOptions{Protocol: "openai-responses", Provider: "deepseek", Model: "deepseek-chat"})
	if err != nil {
		t.Fatalf("provider: %v", err)
	}
	sec := modelSwitchSection(mgr, prov)
	if sec.Adapter != "openai-responses" {
		t.Errorf("adapter = %q, want the built provider's own protocol word", sec.Adapter)
	}
	if want := []string{"high", "low", "medium"}; !reflect.DeepEqual(sec.AdapterThinkingLevels, want) {
		t.Errorf("adapterThinkingLevels = %v, want %v (the responses adapter's own tier words)", sec.AdapterThinkingLevels, want)
	}
	// The SENDABLE set: thinking187Base declares off|low, the responses
	// adapter maps high|low|medium, so exactly "low" survives. "off" drops out
	// because no adapter registers it - it is the request for NO thinking, not
	// a tier. A mapping that shipped the declaration verbatim would offer
	// "off" to an adapter that cannot map it; one that shipped the census
	// verbatim would offer "high"/"medium" to a model that never declared them.
	if want := []string{"low"}; !reflect.DeepEqual(sec.ThinkingLevels, want) {
		t.Errorf("thinkingLevels = %v, want %v (declared off|low intersected with the adapter's high|low|medium)", sec.ThinkingLevels, want)
	}

	// The chat-completions adapter maps NO tier: the intersection is empty
	// even though the model declares two - the adapter's own doc, not a
	// guessed default.
	chatProv, err := llm.NewProvider(llm.EndpointOptions{Protocol: "openai-chat", Provider: "deepseek", Model: "deepseek-chat"})
	if err != nil {
		t.Fatalf("chat provider: %v", err)
	}
	chatSec := modelSwitchSection(mgr, chatProv)
	if chatSec.Adapter != "openai-chat" || len(chatSec.AdapterThinkingLevels) != 0 {
		t.Errorf("openai-chat mapping = %q/%v, want the protocol word with no invented tiers",
			chatSec.Adapter, chatSec.AdapterThinkingLevels)
	}
	if len(chatSec.ThinkingLevels) != 0 {
		t.Errorf("thinkingLevels = %v, want empty: the model declares off|low but this adapter maps no tier at all", chatSec.ThinkingLevels)
	}

	// Byte-level: the census travels under its wire name for the tiered
	// adapter and stays absent for the zero-tier one; the offerable set is
	// the intersection, never the declaration.
	data, err := panel.NewSnapshotPump(panel.PumpSources{
		ModelSwitch: func() panel.ModelSwitchSection { return sec },
	}).Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(data), `"adapter":"openai-responses"`) ||
		!strings.Contains(string(data), `"adapterThinkingLevels":["high","low","medium"]`) {
		t.Errorf("the wire lost the adapter census:\n%s", data)
	}
	if !strings.Contains(string(data), `"thinkingLevels":["low"]`) {
		t.Errorf("the wire does not carry the intersected offer:\n%s", data)
	}
	chatData, err := panel.NewSnapshotPump(panel.PumpSources{
		ModelSwitch: func() panel.ModelSwitchSection { return chatSec },
	}).Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(chatData), "adapterThinkingLevels") {
		t.Errorf("a zero-tier adapter must not ship a census key:\n%s", chatData)
	}
	if !strings.Contains(string(chatData), `"thinkingLevels":[]`) {
		t.Errorf("a zero-tier adapter must ship an empty offer, not the declaration:\n%s", chatData)
	}
}
