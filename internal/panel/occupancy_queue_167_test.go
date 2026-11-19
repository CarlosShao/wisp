package panel

// Ticket 167 AC#2 + AC#3's wire nails (2026-11-19 leg).
//
// The two teeth the dispatch names are aimed at THIS package's render rule:
//   - tooth 1: an unknown numerator must never render as 0 (the pump fills no
//     number when the reader says unknown - mutation-filling it with 0 reds
//     TestOccupancyUnknownNumeratorIsANamedStateNotZero);
//   - tooth 2: an unknown denominator must never be replaced by an invented
//     number (mutation-inventing one reds
//     TestOccupancyUnknownDenominatorNeverInventsOne).
//
// The Q-51 nails for the queue keys are here too: a verdict with no
// correlation id carries NEITHER key on the wire (正控 ⇒ absent), and a
// numbered one carries the panel-cast seq plus the host-read place verbatim
// (the stable number travels through the pump's own mint since AC#7 re-homed
// the cast to the panel side).

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// composerKeysOf marshals a pump's snapshot and returns the composer section's
// JSON keys, so the absence assertions below read the WIRE, not the struct (a
// struct's zero would hide an omitempty lie).
func composerKeysOf(t *testing.T, src PumpSources) (map[string]json.RawMessage, string) {
	t.Helper()
	data, err := NewSnapshotPump(src).Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var wire map[string]json.RawMessage
	if err := json.Unmarshal(data, &wire); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return wire, string(wire["composer"])
}

func TestOccupancyBothKnownTravelVerbatim(t *testing.T) {
	wire, composer := composerKeysOf(t, PumpSources{
		Occupancy: func() *OccupancyReading {
			return &OccupancyReading{UsedTokens: 5200, WindowTokens: 128000}
		},
	})
	var cs map[string]json.RawMessage
	if err := json.Unmarshal(wire["composer"], &cs); err != nil {
		t.Fatal(err)
	}
	if v, ok := cs["contextOccupancyKnown"]; !ok || string(v) != "true" {
		t.Errorf("contextOccupancyKnown = %s (present %v), want true on the wire", v, ok)
	}
	if v, ok := cs["contextTokens"]; !ok || string(v) != "5200" {
		t.Errorf("contextTokens = %s (present %v), want the measured 5200", v, ok)
	}
	if v, ok := cs["contextWindow"]; !ok || string(v) != "128000" {
		t.Errorf("contextWindow = %s (present %v), want the resolved 128000", v, ok)
	}
	_ = composer
}

// TestOccupancyUnknownNumeratorIsANamedStateNotZero is AC#2's named-unknown
// half, numerator side: a host with a window but no usage event yet says
// "unknown" by name and puts NO token number on the wire - the missing numerator
// must not degenerate into 0 (nor the bar into 0%).
func TestOccupancyUnknownNumeratorIsANamedStateNotZero(t *testing.T) {
	wire, _ := composerKeysOf(t, PumpSources{
		Occupancy: func() *OccupancyReading {
			// The honest shape for "window read, nothing measured yet": nil.
			return nil
		},
	})
	var cs map[string]json.RawMessage
	if err := json.Unmarshal(wire["composer"], &cs); err != nil {
		t.Fatal(err)
	}
	if v, ok := cs["contextOccupancyKnown"]; !ok || string(v) != "false" {
		t.Fatalf("contextOccupancyKnown = %s (present %v), want the named unknown state false", v, ok)
	}
	for _, k := range []string{"contextTokens", "contextWindow"} {
		if v, ok := cs[k]; ok {
			t.Errorf("%s = %s on the wire while the dimension is unknown; an absent fact must not travel as a number", k, v)
		}
	}
}

// TestOccupancyUnknownDenominatorNeverInventsOne is the second tooth's fixture:
// a numerator that exists and a denominator that does NOT make the whole
// dimension unknown - the pump must not fill the gap with an invented window
// (the 128k reference or any other constant).
func TestOccupancyUnknownDenominatorNeverInventsOne(t *testing.T) {
	wire, _ := composerKeysOf(t, PumpSources{
		Occupancy: func() *OccupancyReading {
			return &OccupancyReading{UsedTokens: 5200, WindowTokens: 0}
		},
	})
	var cs map[string]json.RawMessage
	if err := json.Unmarshal(wire["composer"], &cs); err != nil {
		t.Fatal(err)
	}
	if v, ok := cs["contextOccupancyKnown"]; !ok || string(v) != "false" {
		t.Errorf("contextOccupancyKnown = %s (present %v), want false when the denominator is unread", v, ok)
	}
	for _, k := range []string{"contextTokens", "contextWindow"} {
		if v, ok := cs[k]; ok {
			t.Errorf("%s = %s on the wire while the denominator is unknown; half a reading must not travel as the whole one", k, v)
		}
	}
}

// TestQueueNumbersTravelAndZeroStaysHome pins the card keys both directions
// (Q-51's omitempty rule), now through the pump's own mint (AC#7 re-homed the
// seq cast to the panel side, so the pump is the only route a number exists
// through): a verdict WITH a correlation id carries the ledger-cast number
// and the host-read FIFO place (2 here - a value that must appear, not be
// swallowed by a zero-check), the same cast is stable across rebuilds, and a
// verdict with NO correlation id carries NEITHER key (nothing to number, no
// place read - absent is the honest shape).
func TestQueueNumbersTravelAndZeroStaysHome(t *testing.T) {
	numbered := NativeVerdict{
		CorrelationID: "corr-167-3", Tool: "fs.write", Level: risk.L2,
		Position: 2,
	}
	pump := NewSnapshotPump(PumpSources{
		Verdicts: func() []NativeVerdict { return []NativeVerdict{numbered} },
	})
	data, err := pump.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	var packet struct {
		Pending []map[string]json.RawMessage `json:"pending"`
	}
	if err := json.Unmarshal(data, &packet); err != nil {
		t.Fatal(err)
	}
	if len(packet.Pending) != 1 {
		t.Fatalf("pending cards = %d, want 1", len(packet.Pending))
	}
	card := packet.Pending[0]
	if v, ok := card["seq"]; !ok || string(v) != "1" {
		t.Errorf("seq = %s (present %v), want the panel-cast number 1 (first sighting)", v, ok)
	}
	if v, ok := card["position"]; !ok || string(v) != "2" {
		t.Errorf("position = %s (present %v), want the FIFO place 2", v, ok)
	}
	// The cast is stable across rebuilds: the same id re-seen gets the SAME
	// number, not a new one (the property the 取消不重排 rule rests on).
	data, err = pump.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	var again struct {
		Pending []map[string]json.RawMessage `json:"pending"`
	}
	if err := json.Unmarshal(data, &again); err != nil {
		t.Fatal(err)
	}
	if v, ok := again.Pending[0]["seq"]; !ok || string(v) != "1" {
		t.Errorf("seq on the rebuild = %s (present %v), want the same cast 1 - the ledger never renumbers a seen id", v, ok)
	}

	plain := NativeVerdict{Tool: "fs.write", Level: risk.L2}
	pump2 := NewSnapshotPump(PumpSources{
		Verdicts: func() []NativeVerdict { return []NativeVerdict{plain} },
	})
	data, err = pump2.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	var bare struct {
		Pending []map[string]json.RawMessage `json:"pending"`
	}
	if err := json.Unmarshal(data, &bare); err != nil {
		t.Fatal(err)
	}
	if len(bare.Pending) != 1 {
		t.Fatalf("bare pending cards = %d, want 1", len(bare.Pending))
	}
	raw, err := json.Marshal(bare.Pending[0])
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), `"seq"`) || strings.Contains(string(raw), `"position"`) {
		t.Errorf("an unnumbered verdict carries a zero key: %s (the CLI card path must stay free of the degenerate 0s)", raw)
	}
}
