package panel

// The appearance wave's panel-side rulers: the fifth roster joins the union
// without touching any older pin, the four snapshot keys are reader-gated
// (a pump with no reader sends no key; a pump with one always sends it), and
// the values a reader-shaped pump sends travel verbatim - the same shape the
// 65 (ball_style_65_test.go) and 413/422 (ball_size_413 round-trip pins)
// carriers were pinned with.
//
// The reverse criteria, named per case: deleting a reader from the assembly
// roots (run.go's PumpSources literal, panel_resident_pump.go's
// residentPumpSources65) reddens the assembly-pin test; deleting a roster
// name reddens the union test and the door would refuse the field on a live
// machine; deleting the Snapshot field or the pump fill reddens the
// byte-nail test.

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestWaveRosterJoinsTheUnionAndStaysDisjoint(t *testing.T) {
	if got := WritableFieldsWave(); !reflect.DeepEqual(got, []string{FieldBallClickThrough, FieldBallOpacityIdle, FieldPanelFontSize, FieldPanelWidth}) {
		t.Fatalf("WritableFieldsWave = %v, want the four wave names sorted", got)
	}
	for _, f := range WritableFieldsWave() {
		if !knownWritableField(f) {
			t.Errorf("knownWritableField refuses %q - the union admission is gone and the settings door would refuse the field on a live machine", f)
		}
		if !takesValueFor(f) {
			t.Errorf("takesValueFor(%q) = false: an empty value would slip past the door's value check and surface as a config-side refusal in the wrong shape", f)
		}
		if needsProviderFor(f) || needsModelFor(f) {
			t.Errorf("%q must not carry a provider/model selector - these four take the plain value half", f)
		}
	}
	// Disjointness with every older roster, the same assertion the older
	// carriers keep: a field lives in exactly one roster.
	older := []struct {
		name string
		rows []string
	}{
		{"248", WritableFields()},
		{"313", AdditionalWritableFields()},
		{"187", WritableFields187()},
		{"65", WritableFields65()},
	}
	for _, f := range WritableFieldsWave() {
		for _, o := range older {
			for _, row := range o.rows {
				if row == f {
					t.Errorf("the %s roster carries %q - rosters are disjoint by construction and the wave field must live in its own (config_write_wave.go)", o.name, f)
				}
			}
		}
	}
	// The union's refusal text tells the page about every name the door
	// accepts: each wave name must be in the accepted list.
	accepted := map[string]bool{}
	for _, f := range acceptedWritableFields() {
		accepted[f] = true
	}
	for _, f := range WritableFieldsWave() {
		if !accepted[f] {
			t.Errorf("acceptedWritableFields() lacks %q - the door's refusal would send the page chasing the wrong list", f)
		}
	}
}

func TestWaveSnapshotKeysAreReaderGated(t *testing.T) {
	// No readers: none of the four keys may appear on the wire.
	data, err := NewSnapshotPump(PumpSources{}).Marshal()
	if err != nil {
		t.Fatalf("bare pump marshal: %v", err)
	}
	for _, key := range []string{"ballOpacityIdle", "ballClickThrough", "panelWidth", "panelFontSize"} {
		if strings.Contains(string(data), `"`+key+`"`) {
			t.Fatalf("a pump with no wave readers sent the %q key: %s", key, data)
		}
	}

	// Readers: each key is always sent, verbatim.
	wantOpacity := BallOpacityIdleSection{Opacity: 0.9, Min: 0, Max: 1}
	wantClick := BallClickThroughSection{Enabled: true}
	wantWidth := PanelWidthSection{Width: 700}
	wantFont := PanelFontSizeSection{FontSize: 15}
	data, err = NewSnapshotPump(PumpSources{
		BallOpacityIdle: func() BallOpacityIdleSection { return wantOpacity },
		BallClickThrough: func() BallClickThroughSection {
			return wantClick
		},
		PanelWidth:    func() PanelWidthSection { return wantWidth },
		PanelFontSize: func() PanelFontSizeSection { return wantFont },
	}).Marshal()
	if err != nil {
		t.Fatalf("reader pump marshal: %v", err)
	}
	var generic map[string]any
	if err := json.Unmarshal(data, &generic); err != nil {
		t.Fatalf("unmarshal the reader pump's bytes: %v", err)
	}
	raw, ok := generic["ballOpacityIdle"]
	if !ok {
		t.Fatalf("a pump with a ballOpacityIdle reader sent no ballOpacityIdle key: %s", data)
	}
	sec, _ := json.Marshal(raw)
	var gotOpacity BallOpacityIdleSection
	if err := json.Unmarshal(sec, &gotOpacity); err != nil || gotOpacity != wantOpacity {
		t.Fatalf("ballOpacityIdle = %s (%v), want %+v", sec, err, wantOpacity)
	}
	raw, ok = generic["ballClickThrough"]
	if !ok {
		t.Fatalf("a pump with a ballClickThrough reader sent no ballClickThrough key: %s", data)
	}
	sec, _ = json.Marshal(raw)
	var gotClick BallClickThroughSection
	if err := json.Unmarshal(sec, &gotClick); err != nil || gotClick != wantClick {
		t.Fatalf("ballClickThrough = %s (%v), want %+v", sec, err, wantClick)
	}
	raw, ok = generic["panelWidth"]
	if !ok {
		t.Fatalf("a pump with a panelWidth reader sent no panelWidth key: %s", data)
	}
	sec, _ = json.Marshal(raw)
	var gotWidth PanelWidthSection
	if err := json.Unmarshal(sec, &gotWidth); err != nil || gotWidth != wantWidth {
		t.Fatalf("panelWidth = %s (%v), want %+v", sec, err, wantWidth)
	}
	raw, ok = generic["panelFontSize"]
	if !ok {
		t.Fatalf("a pump with a panelFontSize reader sent no panelFontSize key: %s", data)
	}
	sec, _ = json.Marshal(raw)
	var gotFont PanelFontSizeSection
	if err := json.Unmarshal(sec, &gotFont); err != nil || gotFont != wantFont {
		t.Fatalf("panelFontSize = %s (%v), want %+v", sec, err, wantFont)
	}
}
