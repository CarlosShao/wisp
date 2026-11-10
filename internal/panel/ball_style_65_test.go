package panel

// Ticket 65's closing-leg tests: the settings DOOR now admits ball_style and
// the ballStyle key is a reader-gated section like commands (213) and
// modelSwitch (187). The write leg underneath the door (config setter, boot
// seed, applier) is tested for real in cmd/wisp (ball_style_65_windows_test.go)
// and internal/config (settings_65_test.go); these tests pin the panel half:
//
//   - the roster joins the union WITHOUT moving any older pin (the frozen 248
//     table, the 313 and 187 rosters and the disjointness between them keep
//     meaning exactly what they meant);
//   - a config.set naming ball_style reaches the leg - the tooth that goes red
//     if the union admission (knownWritableField's writableField65 term) is
//     ever dropped, which is the mutation this closing leg exists to prevent;
//   - the wire key travels only from its reader, values verbatim, nil
//     vocabulary rendered as [].
//
// Every judgement is a call count or an error identity, never a log string -
// the same counting rule config_route_248_test.go states for the doors these
// tests share.

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// Test65RosterJoinsTheUnionWithoutMovingOlderPins: the fourth roster is
// admitted by the union and by both helper unions (the refusal text and the
// empty-value check), and none of the three older rosters grew to carry it.
func Test65RosterJoinsTheUnionWithoutMovingOlderPins(t *testing.T) {
	// Ticket 422 grew this roster to two fields (ball_size joined ball_style
	// here, per the 413 contract's landing note); the pin stays exact-set.
	if got := WritableFields65(); !reflect.DeepEqual(got, []string{FieldBallSize, FieldBallStyle}) {
		t.Fatalf("WritableFields65() = %v, want exactly [ball_size ball_style]", got)
	}
	if !knownWritableField(FieldBallStyle) {
		t.Fatalf("knownWritableField refuses %q - the union admission is gone and the settings door would refuse the appearance switch on a live machine", FieldBallStyle)
	}
	for _, older := range []struct {
		name string
		rows []string
	}{
		{"248", WritableFields()},
		{"313", AdditionalWritableFields()},
		{"187", WritableFields187()},
	} {
		if containsString(older.rows, FieldBallStyle) {
			t.Errorf("the %s roster carries %q - rosters are disjoint by construction and the 65 field must live in its own (config_write_65.go)", older.name, FieldBallStyle)
		}
	}
	if !containsString(acceptedWritableFields(), FieldBallStyle) {
		t.Errorf("acceptedWritableFields() = %v, want it to carry the 65 field: the refusal text would send the page chasing the wrong list", acceptedWritableFields())
	}
	if !takesValueFor(FieldBallStyle) {
		t.Errorf("takesValueFor(%q) = false: an empty value would slip past the door's value check and surface as a config-side enum refusal in the wrong shape", FieldBallStyle)
	}
	if needsProviderFor(FieldBallStyle) || needsModelFor(FieldBallStyle) {
		t.Errorf("ball_style takes the plain value half; it must not demand a provider or model selector")
	}
}

// Test65DoorAdmitsBallStyleAndRoutesItToTheLeg: config.set naming ball_style
// clears the door and reaches the leg exactly once with the parsed value - the
// MUTATION TOOTH. Drop `|| writableField65(field)` from knownWritableField and
// this test goes red on the refusal error (the door answers "不在可写名册内"
// and the leg's apply count stays 0).
func Test65DoorAdmitsBallStyleAndRoutesItToTheLeg(t *testing.T) {
	leg := &cfg248Leg{res: SettingWriteResult{
		Field: FieldBallStyle, Written: []string{"ball.style"}, Tier: EffectiveNow,
	}}
	reply, err := cfg248Dispatch(leg).Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigSet, "requestId": "r-65", "source": ComposerRequestSource,
		"configField": FieldBallStyle, "configValue": "jarvis",
	}))
	if err != nil {
		t.Fatalf("the door refused ball_style: %v", err)
	}
	if leg.applyCalls != 1 {
		t.Fatalf("leg apply calls = %d, want 1", leg.applyCalls)
	}
	if leg.lastWrite.Field != FieldBallStyle || leg.lastWrite.Value != "jarvis" {
		t.Fatalf("leg saw %+v, want field ball_style with value jarvis", leg.lastWrite)
	}
	if !strings.Contains(reply, "ball.style") {
		t.Errorf("the receipt never named the key path: %q", reply)
	}

	// Positive control for the door's value check (the takesValueFor seam): an
	// empty value is refused BEFORE the leg, never written as a default.
	empty := &cfg248Leg{}
	if _, err := cfg248Dispatch(empty).Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigSet, "requestId": "r-65-empty", "source": ComposerRequestSource,
		"configField": FieldBallStyle, "configValue": "",
	})); err == nil {
		t.Fatalf("an empty ball_style value must be refused at the door")
	} else if empty.applyCalls != 0 {
		t.Errorf("the refused write reached the leg %d times, want 0", empty.applyCalls)
	}

	// The negative control the positive one is judged against: a name no
	// roster carries is still refused, and the refusal is the door's own.
	unlisted := &cfg248Leg{}
	if _, err := cfg248Dispatch(unlisted).Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigSet, "requestId": "r-65-bad", "source": ComposerRequestSource,
		"configField": "orb_theme", "configValue": "dark",
	})); err == nil {
		t.Fatalf("a field no roster carries must be refused")
	} else if unlisted.applyCalls != 0 {
		t.Errorf("the unlisted write reached the leg %d times, want 0", unlisted.applyCalls)
	}
}

// TestTheBallStyleKeyTravelsOnlyFromItsReader65 pins the same two directions
// the 213 and 187 carrier tests pin: no reader, no key; a reader that exists
// always sends the key with every field verbatim and a nil vocabulary as [].
func TestTheBallStyleKeyTravelsOnlyFromItsReader65(t *testing.T) {
	// (1) No reader, no key: the four-key shape the byte nails pin, with this
	// key named so a regression says which key came back.
	data, err := NewSnapshotPump(PumpSources{}).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), `"ballStyle"`) {
		t.Fatalf("a pump with no ballStyle reader sent the key: %s", data)
	}
	var generic map[string]json.RawMessage
	if err := json.Unmarshal(data, &generic); err != nil {
		t.Fatal(err)
	}
	if len(generic) != 4 {
		t.Errorf("top-level keys = %d, want the four pump_test pins for a reader-less pump", len(generic))
	}

	// (2) With a reader the key arrives and both fields travel verbatim.
	want := BallStyleSection{Style: "jarvis", Vocabulary: []string{"glass", "jarvis"}}
	data, err = NewSnapshotPump(PumpSources{
		BallStyle: func() BallStyleSection { return want },
	}).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &generic); err != nil {
		t.Fatal(err)
	}
	raw, ok := generic["ballStyle"]
	if !ok {
		t.Fatalf("a pump with a ballStyle reader sent no ballStyle key: %s", data)
	}
	var got BallStyleSection
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("wire section = %+v, want %+v (every field travels verbatim, or the picker lies about what is offered)", got, want)
	}

	// (3) A nil vocabulary renders as [] - the page reads the list like it
	// reads every other roster, and null breaks that.
	data, err = NewSnapshotPump(PumpSources{
		BallStyle: func() BallStyleSection { return BallStyleSection{Style: "glass"} },
	}).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"vocabulary":[]`) {
		t.Errorf("a nil vocabulary did not render as []: %s", data)
	}
}
