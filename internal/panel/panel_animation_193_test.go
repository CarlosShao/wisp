package panel

// Ticket 193's panel-side rulers: the panelAnimation carrier rides the
// reader's tier word verbatim, a pump without the reader sends NO key at all
// (never a made-up "off"), and the carrier adds no timer and no second
// paradigm. These are the Go half of ticket 193 AC#4's off-proof: the wire
// word the page consumes is the config's own word, so "off" reaching the page
// is a fact the page can be held to - whether the page's rAF handles and CSS
// transitions really stop is the frontend half, and its checklist rides
// docs/evidence/s1/193-animation-opt-r1.md.

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestPanelAnimationCarrierRidesTheReaderVerbatim(t *testing.T) {
	for _, tier := range []string{"off", "reduce", "full"} {
		want := tier
		pump := NewSnapshotPump(PumpSources{
			PanelAnimation: func() PanelAnimationSection { return PanelAnimationSection{Tier: want} },
		})
		data, err := pump.Marshal()
		if err != nil {
			t.Fatalf("marshal the %q pump: %v", tier, err)
		}
		var generic map[string]json.RawMessage
		if err := json.Unmarshal(data, &generic); err != nil {
			t.Fatal(err)
		}
		raw, ok := generic["panelAnimation"]
		if !ok {
			t.Fatalf("tier %q: the pump had the reader but sent no panelAnimation key - a reader that exists always sends the key", tier)
		}
		var sec PanelAnimationSection
		if err := json.Unmarshal(raw, &sec); err != nil {
			t.Fatalf("tier %q: the panelAnimation key is not a PanelAnimationSection: %v", tier, err)
		}
		if sec.Tier != want {
			t.Fatalf("the wire carried %q, want the config's own word %q - the pump derives nothing and rewrites nothing", sec.Tier, want)
		}
	}
}

func TestPanelAnimationCarrierWithNoReaderSendsNoKeyAndNeverInventsOff(t *testing.T) {
	// The reader-gated shape the four-key byte nails pin: without the reader
	// the key is absent. The inversion this case adds for ticket 193 AC#4:
	// the ABSENCE must never read as "off" - an unset switch is not a stopped
	// animation, and a pump that answered "off" with no Manager behind it
	// would be the page's lie to build an off-proof on.
	data, err := NewSnapshotPump(PumpSources{}).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	var generic map[string]json.RawMessage
	if err := json.Unmarshal(data, &generic); err != nil {
		t.Fatal(err)
	}
	if _, ok := generic["panelAnimation"]; ok {
		t.Fatalf("no reader, yet the packet carried panelAnimation - the no-reader shape must send no key at all")
	}
	if strings.Contains(string(data), `"tier"`) {
		t.Fatalf("no reader, yet the packet carried a tier word anywhere - the pump must not invent a tier for a switch nobody read")
	}
}

func TestPanelAnimationCarrierSendsTheUnsetTierItRead(t *testing.T) {
	// A live reader over a nil-Manager config answers "" (the honest empty
	// PanelFontSizeSection's 0 follows): the wire then says {"tier":""} -
	// the page's "unset", never "off", never "reduce".
	pump := NewSnapshotPump(PumpSources{
		PanelAnimation: func() PanelAnimationSection { return PanelAnimationSection{Tier: ""} },
	})
	data, err := pump.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"panelAnimation":{"tier":""}`) {
		t.Fatalf("the unset tier did not ride verbatim: %s", data)
	}
}
