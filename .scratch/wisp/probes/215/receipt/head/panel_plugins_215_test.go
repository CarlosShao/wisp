package main

// Ticket 215a's assembly-root tooth: run.go's PumpSources literal hands the
// pump the plugins reader, so the Tier-2 plugin inventory carrier
// (internal/panel/plugins_inventory_215.go, read off the live config's
// [plugins] section) actually reaches the wire. The carrier and the pump half
// landed first (5341ead7a) with cmd/wisp held by another leg, and the evidence
// doc registered this one line as owed; without it the running process sends
// no "plugins" key at all - the exact "carrier exists, nothing reads it" hop
// ticket 200 AC#7 closed for its own dimension, ticket 65 closed for ballStyle
// and ticket 187 closed for modelSwitch, and the shape the reader-gated tests
// in internal/panel cannot see from over there.
//
// Teeth:
//  1. the wiring exists: run.go's literal still carries the reader line
//     (the 65/187 byte-nail shape - a deleted line reds here immediately);
//  2. the wired reading is right: a real Manager loaded from a config with
//     [plugins] populated, read through the same chain the literal uses,
//     marshaled through the production pump, sends the "plugins" key whose
//     fields carry the config's own values verbatim (hardcoded here so the
//     nail is independent of the reader's code) with the rows sorted by ID -
//     the pump owns the order because Entries is a Go map;
//  3. the mutation: with the reader gone the wire bytes carry no "plugins"
//     key at all - the reader-gated 213 shape, asserted on the bytes the
//     panel would receive, not a Go field.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
)

// plugins215Base is the same minimal valid config the 187 tests load, with a
// [plugins] section appended: the master switch on, two entries with real
// grant lists (the engine itself is ticket 50/51 and unbuilt - the config
// layer loads and audits these keys today, unwired.go's disposition rows
// name them "not built", and that honesty is the carrier's sentence, not a
// reason to leave the inventory unread).
const plugins215Base = `schema_version = 2

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

[plugins]
tier2_enabled = true

[plugins.clip]
enabled = true
capabilities = ["fs.read", "net.fetch"]
net_allowlist = ["example.com"]
host_api = ["storage"]

[plugins.zed]
enabled = false
capabilities = ["fs.write"]
`

func Test215AssemblyRootCarriesThePluginsReader(t *testing.T) {
	runGo, err := os.ReadFile("run.go")
	if err != nil {
		t.Fatalf("read run.go: %v", err)
	}
	if !strings.Contains(string(runGo), "Plugins: func() panel.PluginInventorySection") {
		t.Error("run.go's PumpSources literal no longer carries the plugins reader: the snapshot's plugins key would silently vanish from every packet this process pushes (ticket 215a)")
	}
}

// pluginsSection215 mirrors run.go's PumpSources literal mapping verbatim
// (mgr.Config()'s [plugins] section into panel.PluginInventorySection, the
// carrier's own two sentences). The byte-nail above pins the literal's
// presence; this helper lets tooth 2 exercise the same read chain against a
// real Manager. If the literal's mapping ever diverges from this helper,
// tooth 2's hardcoded content still holds for the helper - the divergence
// itself is what the evidence doc asks a reviewer to diff by eye, one loop
// plus four fields.
func pluginsSection215(m *config.Manager) panel.PluginInventorySection {
	c := m.Config()
	sec := panel.PluginInventorySection{
		Tier2Enabled: c.Plugins.Tier2Enabled,
		EngineState:  panel.PluginEngineNotWired215,
		SkillState:   panel.SkillsNotInConfigModel215,
	}
	for id, e := range c.Plugins.Entries {
		sec.Rows = append(sec.Rows, panel.PluginRow{
			ID:           id,
			Enabled:      e.Enabled,
			Capabilities: e.Capabilities,
			NetAllowlist: e.NetAllowlist,
			HostAPI:      e.HostAPI,
		})
	}
	return sec
}

// Test215WiredPluginsReaderSendsTheLiveValues: the reader chain the assembly
// root uses, against a real Manager, through the production pump's marshal -
// the key arrives, every field carries the live config's values verbatim, the
// rows are sorted by ID (map order must not leak onto the wire), and the two
// honest-state sentences ride the same packet so enabled=true can never read
// as "running" (no consumer of these keys exists until the 50/51 engine
// lands).
func Test215WiredPluginsReaderSendsTheLiveValues(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, configFileName)
	if err := os.WriteFile(path, []byte(plugins215Base), 0o600); err != nil {
		t.Fatal(err)
	}
	mgr, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("manager: %v", err)
	}
	if live := pluginsSection215(mgr); len(live.Rows) != 2 || !live.Tier2Enabled {
		t.Fatalf("the loaded config does not carry the [plugins] section the reader must read: tier2Enabled=%v rows=%+v", live.Tier2Enabled, live.Rows)
	}

	data, err := panel.NewSnapshotPump(panel.PumpSources{
		Plugins: func() panel.PluginInventorySection { return pluginsSection215(mgr) },
	}).Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var envelope struct {
		Plugins *panel.PluginInventorySection `json:"plugins"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		t.Fatalf("the wire bytes are not the snapshot's JSON: %v", err)
	}
	if envelope.Plugins == nil {
		t.Fatalf("the wired reader sent no plugins key at all:\n%s", data)
	}
	got := *envelope.Plugins

	// Hardcoded from plugins215Base. The map the config layer holds has no
	// order, so the ID order on the wire is the pump's sort doing its job.
	if !got.Tier2Enabled {
		t.Errorf("plugins.tier2Enabled = false, want the config's true verbatim")
	}
	if len(got.Rows) != 2 || got.Rows[0].ID != "clip" || got.Rows[1].ID != "zed" {
		t.Fatalf("plugins rows = %+v, want exactly [clip zed] (sorted by ID)", got.Rows)
	}
	wantClip := panel.PluginRow{
		ID:           "clip",
		Enabled:      true,
		Capabilities: []string{"fs.read", "net.fetch"},
		NetAllowlist: []string{"example.com"},
		HostAPI:      []string{"storage"},
	}
	if !reflect.DeepEqual(got.Rows[0], wantClip) {
		t.Errorf("wire row clip = %+v, want %+v (every field travels verbatim, or the menu lies about a grant)", got.Rows[0], wantClip)
	}
	wantZed := panel.PluginRow{
		ID:           "zed",
		Enabled:      false,
		Capabilities: []string{"fs.write"},
	}
	if !reflect.DeepEqual(got.Rows[1], wantZed) {
		t.Errorf("wire row zed = %+v, want %+v (absent grant lists travel as the config's nil)", got.Rows[1], wantZed)
	}
	if got.EngineState != panel.PluginEngineNotWired215 {
		t.Errorf("EngineState = %q, want the canonical sentence verbatim: %q", got.EngineState, panel.PluginEngineNotWired215)
	}
	if got.SkillState != panel.SkillsNotInConfigModel215 {
		t.Errorf("SkillState = %q, want the canonical sentence verbatim: %q", got.SkillState, panel.SkillsNotInConfigModel215)
	}
	// The menu must agree with what the config layer holds - one read face,
	// not two that can drift.
	if cfg := mgr.Config(); cfg.Plugins.Tier2Enabled != got.Tier2Enabled || len(cfg.Plugins.Entries) != len(got.Rows) {
		t.Errorf("wire section disagrees with the live config (tier2Enabled %v/%v, rows %d/%d): the picker lies about the config",
			got.Tier2Enabled, cfg.Plugins.Tier2Enabled, len(got.Rows), len(cfg.Plugins.Entries))
	}
}

// Test215MutationReaderRemovedSendsNoPluginsKeyAtAll: the assembly root's
// reader gone - the pump the four-key byte nails pin. The key must be absent
// from the WIRE BYTES, not merely nil in Go: the omitempty tag is the last
// mile, and this asserts the bytes the panel receives.
func Test215MutationReaderRemovedSendsNoPluginsKeyAtAll(t *testing.T) {
	data, err := panel.NewSnapshotPump(panel.PumpSources{}).Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if bytes.Contains(data, []byte(`"plugins"`)) {
		t.Errorf("a pump without the Plugins reader still sent the key:\n%s", data)
	}
}
