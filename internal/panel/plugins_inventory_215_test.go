package panel

// Ticket 215a's carrier tests. The plugins key is a reader-gated section like
// commands (213) and ballStyle (65), so its tests pin the same two directions
// those tests pin, plus the one thing this ticket exists to make visible:
//
//   - WITHOUT the reader the key is absent from the bytes a real pump
//     marshals - its presence elsewhere is caused by the reader, never by a
//     field that is always on the wire with nothing behind it (the three
//     four-key byte nails keep pinning that shape; this case names the
//     plugins key itself).
//   - WITH the reader the key always arrives: rows travel verbatim (deep
//     equality, not a spot check), nil rows render as [] and never null, and
//     the two honest-state sentences travel verbatim - a packet with rows and
//     a blank EngineState/SkillState is the 文案在、控件不在 family A408
//     names, so blank sentences red here.
//   - AC#2's "配置开着但没生效必须能被看出来" is pinned as a POSITIVE
//     CONTROL in both directions: an enabled=true row and the not-wired
//     sentence must appear on the SAME packet (the boolean never reads as
//     "running"), and the canonical sentence spellings are pinned to the
//     carrier's own constants so no forked prose can silently replace them.
//
// Mutation teeth (run against the pristine tree, each reds, then revert):
//   M1 - drop the pump's Plugins assembly block: the key vanishes from a
//        reader-ful pump's packet (test 1, part 2).
//   M2 - drop EngineState/SkillState from the carrier: the reader-ful wire
//        loses the sentences (tests 1 part 4 / 2).
//   M3 - drop the READER at an assembly root while the config still holds
//        enabled=true: the packet carries no plugins key at all, so a config
//        that says "on" is invisible on the wire - the exact shape the
//        carrier's key-absence rule makes detectable, demonstrated in
//        docs/evidence/s1/215-skills-plugins-inventory.md (a standing test
//        cannot plant a config value in internal/panel: the reader is the
//        caller's injection seam, which is the paradigm).

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestThePluginsKeyTravelsOnlyFromItsReader215(t *testing.T) {
	// (1) No reader, no key: the four-key shape the byte nails pin, and the
	// plugins key named so a regression here says which key came back.
	data, err := NewSnapshotPump(PumpSources{}).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), `"plugins"`) {
		t.Fatalf("a pump with no plugins reader sent the key: %s", data)
	}
	var generic map[string]json.RawMessage
	if err := json.Unmarshal(data, &generic); err != nil {
		t.Fatal(err)
	}
	if len(generic) != 4 {
		t.Errorf("top-level keys = %d, want the four pump_test pins for a reader-less pump", len(generic))
	}

	// (2) With a reader the key arrives, rows travel verbatim in the pump's
	// ID order, and the two sentences travel verbatim.
	reader := PluginInventorySection{
		Tier2Enabled: true,
		Rows: []PluginRow{
			{ID: "zeta", Enabled: true, Capabilities: []string{"fs.read"}, NetAllowlist: []string{"example.com"}, HostAPI: []string{"storage"}},
			{ID: "alpha", Enabled: false, Capabilities: []string{"net.fetch"}, NetAllowlist: []string{"a.example", "b.example"}, HostAPI: nil},
		},
		EngineState: PluginEngineNotWired215,
		SkillState:  SkillsNotInConfigModel215,
	}
	pump := NewSnapshotPump(PumpSources{Plugins: func() PluginInventorySection { return reader }})
	data, err = pump.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &generic); err != nil {
		t.Fatal(err)
	}
	raw, ok := generic["plugins"]
	if !ok {
		t.Fatalf("a pump with a plugins reader sent no plugins key: %s", data)
	}
	var got PluginInventorySection
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if got.Tier2Enabled != reader.Tier2Enabled {
		t.Errorf("tier2Enabled = %v, want %v verbatim", got.Tier2Enabled, reader.Tier2Enabled)
	}
	if len(got.Rows) != 2 {
		t.Fatalf("wire rows = %d, want 2: %s", len(got.Rows), raw)
	}
	if got.Rows[0].ID != "alpha" || got.Rows[1].ID != "zeta" {
		t.Errorf("wire rows not sorted by ID (map order must not leak onto the wire): [%s, %s]", got.Rows[0].ID, got.Rows[1].ID)
	}
	wantAlpha := reader.Rows[1]
	if !reflect.DeepEqual(got.Rows[0], wantAlpha) {
		t.Errorf("wire row alpha = %+v, want %+v (every field travels verbatim, or the menu lies about a grant)", got.Rows[0], wantAlpha)
	}
	wantZeta := reader.Rows[0]
	if !reflect.DeepEqual(got.Rows[1], wantZeta) {
		t.Errorf("wire row zeta = %+v, want %+v", got.Rows[1], wantZeta)
	}
	if got.EngineState != PluginEngineNotWired215 {
		t.Errorf("EngineState = %q, want the canonical sentence verbatim: %q", got.EngineState, PluginEngineNotWired215)
	}
	if got.SkillState != SkillsNotInConfigModel215 {
		t.Errorf("SkillState = %q, want the canonical sentence verbatim: %q", got.SkillState, SkillsNotInConfigModel215)
	}

	// (3) Nil rows from a reader render as [] - the page reads .length on
	// this list like it does on the other sections' lists, and null breaks
	// that.
	pump = NewSnapshotPump(PumpSources{
		Plugins: func() PluginInventorySection {
			return PluginInventorySection{EngineState: PluginEngineNotWired215, SkillState: SkillsNotInConfigModel215}
		},
	})
	data, err = pump.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"rows":[]`) {
		t.Errorf("nil rows did not render as []: %s", data)
	}

	// (4) The pump carries the honest-state sentences VERBATIM and invents
	// neither one: a dishonest reader's blank arrives blank on the wire - a
	// visible hole, never papered over by a pump-side default (the pump
	// normalizes nil lists, it does not author sentences). The assertion is
	// byte-exact: the reader's spelling IS the wire's spelling.
	for _, blank := range []string{"EngineState", "SkillState"} {
		sec := PluginInventorySection{
			Rows:        []PluginRow{{ID: "clip", Enabled: true}},
			EngineState: PluginEngineNotWired215,
			SkillState:  SkillsNotInConfigModel215,
		}
		switch blank {
		case "EngineState":
			sec.EngineState = ""
		case "SkillState":
			sec.SkillState = ""
		}
		pump = NewSnapshotPump(PumpSources{Plugins: func() PluginInventorySection { return sec }})
		data, err = pump.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(data, &generic); err != nil {
			t.Fatal(err)
		}
		raw, ok = generic["plugins"]
		if !ok {
			t.Fatalf("the plugins key vanished for the blank-%s case: %s", blank, data)
		}
		var probe PluginInventorySection
		if err := json.Unmarshal(raw, &probe); err != nil {
			t.Fatal(err)
		}
		got := probe.EngineState
		if blank == "SkillState" {
			got = probe.SkillState
		}
		if got != "" {
			t.Errorf("%s was invented or rewritten on the wire, want the blank carried verbatim: %q", blank, got)
		}
	}
}

func TestPluginInventory215EnabledRowAndNotWiredSentenceRideOnePacket(t *testing.T) {
	// AC#2's positive control, standing half: an enabled=true row with real
	// grant lists and the not-wired sentence must arrive TOGETHER - the
	// switch's fact and the "nothing consumes it" fact are one packet, so the
	// boolean can never read as "running" and the sentence can never be read
	// as absent. The spelling is pinned to the carrier's constants.
	sec := PluginInventorySection{
		Tier2Enabled: true,
		Rows: []PluginRow{{
			ID:           "clip",
			Enabled:      true,
			Capabilities: []string{"fs.read", "net.fetch"},
			NetAllowlist: []string{"example.com"},
			HostAPI:      []string{"storage"},
		}},
		EngineState: PluginEngineNotWired215,
		SkillState:  SkillsNotInConfigModel215,
	}
	data, err := NewSnapshotPump(PumpSources{Plugins: func() PluginInventorySection { return sec }}).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	packet := string(data)
	if !strings.Contains(packet, `"enabled":true`) {
		t.Errorf("the enabled fact did not travel: %s", packet)
	}
	if !strings.Contains(packet, PluginEngineNotWired215) {
		t.Errorf("the not-wired sentence did not travel beside the enabled fact: %s", packet)
	}
	if !strings.Contains(packet, SkillsNotInConfigModel215) {
		t.Errorf("the skill-dimension registration sentence did not travel: %s", packet)
	}
}

func TestPluginInventory215PublishesByteIdenticallyForOneConfig(t *testing.T) {
	// Rows come from a Go map at the config layer, so the reader cannot
	// promise an order; the pump owns the sort, and two publishes of one
	// config must marshal byte-identically (a menu reshuffling every push is
	// an audit defect).
	reader := PluginInventorySection{
		Rows: []PluginRow{
			{ID: "m3", Enabled: true},
			{ID: "m1", Enabled: false},
			{ID: "m2", Enabled: true},
		},
		EngineState: PluginEngineNotWired215,
		SkillState:  SkillsNotInConfigModel215,
	}
	pump := NewSnapshotPump(PumpSources{Plugins: func() PluginInventorySection { return reader }})
	first, err := pump.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	second, err := pump.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if string(first) != string(second) {
		t.Errorf("two publishes of one config differ:\nfirst:  %s\nsecond: %s", first, second)
	}
	var generic map[string]json.RawMessage
	if err := json.Unmarshal(first, &generic); err != nil {
		t.Fatal(err)
	}
	var got PluginInventorySection
	if err := json.Unmarshal(generic["plugins"], &got); err != nil {
		t.Fatal(err)
	}
	if len(got.Rows) != 3 || got.Rows[0].ID != "m1" || got.Rows[1].ID != "m2" || got.Rows[2].ID != "m3" {
		t.Errorf("rows not sorted by ID: %+v", got.Rows)
	}
}
