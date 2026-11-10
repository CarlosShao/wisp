package panel

// Ticket 422's closing-leg tests: the settings DOOR now admits ball_size and
// the ballSize key is a reader-gated section like its ball_style twin (65).
// The write leg underneath the door (SetBallSize, the 44..160 band
// single-sourced into validate.go's ballSizeRangeError) is tested for real in
// internal/config (settings_413's own tests) and cmd/wisp
// (ball_size_422_windows_test.go); these tests pin the panel half:
//
//   - ball_size joins the union WITHOUT moving any older pin (the frozen 248
//     table, the 313 and 187 rosters and the disjointness between them keep
//     meaning exactly what they meant);
//   - a config.set naming ball_size reaches the leg - the tooth that goes red
//     if the roster admission (knownWritableField's writableField65 term over
//     the grown fields65) is ever dropped;
//   - the wire key travels only from its reader, values verbatim, a
//     not-yet-configured diameter traveling as Size 0 rather than as a
//     vanished key.
//
// Judgement shape mirrors ball_style_65_test.go: every judgement is a call
// count or an error identity, never a log string.

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// Test422RosterJoinsTheUnionWithoutMovingOlderPins: ball_size is admitted by
// the union and by both helper unions (the refusal text and the empty-value
// check), and none of the three older rosters grew to carry it.
func Test422RosterJoinsTheUnionWithoutMovingOlderPins(t *testing.T) {
	if !containsString(WritableFields65(), FieldBallSize) {
		t.Fatalf("WritableFields65() = %v, want it to carry %q (ticket 422 grew the 65 roster)", WritableFields65(), FieldBallSize)
	}
	if !knownWritableField(FieldBallSize) {
		t.Fatalf("knownWritableField refuses %q - the union admission is gone and the settings door would refuse the size write on a live machine", FieldBallSize)
	}
	for _, older := range []struct {
		name string
		rows []string
	}{
		{"248", WritableFields()},
		{"313", AdditionalWritableFields()},
		{"187", WritableFields187()},
	} {
		if containsString(older.rows, FieldBallSize) {
			t.Errorf("the %s roster carries %q - rosters are disjoint by construction and the 422 field must live in the 65 one (config_write_65.go)", older.name, FieldBallSize)
		}
	}
	if !containsString(acceptedWritableFields(), FieldBallSize) {
		t.Errorf("acceptedWritableFields() = %v, want it to carry the 422 field: the refusal text would send the page chasing the wrong list", acceptedWritableFields())
	}
	if !takesValueFor(FieldBallSize) {
		t.Errorf("takesValueFor(%q) = false: an empty value would slip past the door's value check and surface as a config-side band refusal in the wrong shape", FieldBallSize)
	}
	if needsProviderFor(FieldBallSize) || needsModelFor(FieldBallSize) {
		t.Errorf("ball_size takes the plain value half; it must not demand a provider or model selector")
	}
}

// Test422DoorAdmitsBallSizeAndRoutesItToTheLeg: config.set naming ball_size
// clears the door and reaches the leg exactly once with the parsed value's
// string - the MUTATION TOOTH. Drop FieldBallSize from fields65 and this test
// goes red on the refusal error (the door answers "不在可写名册内" and the
// leg's apply count stays 0).
func Test422DoorAdmitsBallSizeAndRoutesItToTheLeg(t *testing.T) {
	leg := &cfg248Leg{res: SettingWriteResult{
		Field: FieldBallSize, Written: []string{"ball.size"}, Tier: EffectiveRestart,
	}}
	reply, err := cfg248Dispatch(leg).Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigSet, "requestId": "r-422", "source": ComposerRequestSource,
		"configField": FieldBallSize, "configValue": "90",
	}))
	if err != nil {
		t.Fatalf("the door refused ball_size: %v", err)
	}
	if leg.applyCalls != 1 {
		t.Fatalf("leg apply calls = %d, want 1", leg.applyCalls)
	}
	if leg.lastWrite.Field != FieldBallSize || leg.lastWrite.Value != "90" {
		t.Fatalf("leg saw %+v, want field ball_size with value \"90\"", leg.lastWrite)
	}
	if !strings.Contains(reply, "ball.size") {
		t.Errorf("the receipt never named the key path: %q", reply)
	}

	// Positive control for the door's value check (the takesValueFor seam): an
	// empty value is refused BEFORE the leg, never written as a default.
	empty := &cfg248Leg{}
	if _, err := cfg248Dispatch(empty).Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigSet, "requestId": "r-422-empty", "source": ComposerRequestSource,
		"configField": FieldBallSize, "configValue": "",
	})); err == nil {
		t.Fatalf("an empty ball_size value must be refused at the door")
	} else if empty.applyCalls != 0 {
		t.Errorf("the refused write reached the leg %d times, want 0", empty.applyCalls)
	}

	// The negative control the positive one is judged against: a name no
	// roster carries is still refused, and the refusal is the door's own.
	unlisted := &cfg248Leg{}
	if _, err := cfg248Dispatch(unlisted).Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigSet, "requestId": "r-422-bad", "source": ComposerRequestSource,
		"configField": "orb_size", "configValue": "90",
	})); err == nil {
		t.Fatalf("a field no roster carries must be refused")
	} else if unlisted.applyCalls != 0 {
		t.Errorf("the unlisted write reached the leg %d times, want 0", unlisted.applyCalls)
	}
}

// TestTheBallSizeKeyTravelsOnlyFromItsReader422 pins the same two directions
// the 65 carrier test pins: no reader, no key; a reader that exists always
// sends the key with every field verbatim - including a Size 0 section, which
// is the honest "not configured" answer and must not be mistaken for the
// key-vanishing shape a nil reader produces.
func TestTheBallSizeKeyTravelsOnlyFromItsReader422(t *testing.T) {
	// (1) No reader, no key: the four-key shape the byte nails pin, with this
	// key named so a regression says which key came back.
	data, err := NewSnapshotPump(PumpSources{}).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), `"ballSize"`) {
		t.Fatalf("a pump with no ballSize reader sent the key: %s", data)
	}
	var generic map[string]json.RawMessage
	if err := json.Unmarshal(data, &generic); err != nil {
		t.Fatal(err)
	}
	if len(generic) != 4 {
		t.Errorf("top-level keys = %d, want the four pump_test pins for a reader-less pump", len(generic))
	}

	// (2) With a reader the key arrives and all three fields travel verbatim.
	want := BallSizeSection{Size: 90, Min: 44, Max: 160}
	data, err = NewSnapshotPump(PumpSources{
		BallSize: func() BallSizeSection { return want },
	}).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &generic); err != nil {
		t.Fatal(err)
	}
	raw, ok := generic["ballSize"]
	if !ok {
		t.Fatalf("a pump with a ballSize reader sent no ballSize key: %s", data)
	}
	var got BallSizeSection
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("wire section = %+v, want %+v (every field travels verbatim, or the slider lies about what is offered)", got, want)
	}

	// (3) A zero diameter still sends the key - Size 0 is the section's own
	// "no Manager backs the read" sentence (ball_size_413.go), not an absence
	// to be papered over with omitempty semantics.
	data, err = NewSnapshotPump(PumpSources{
		BallSize: func() BallSizeSection { return BallSizeSection{} },
	}).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"ballSize":{"size":0,"min":0,"max":0}`) {
		t.Errorf("a Size 0 section did not travel as the key's honest not-configured shape: %s", data)
	}
}
