package panel

// Ticket 193's closing-leg tests, panel half: the settings DOOR now admits
// panel_animation (the ninth roster, config_write_193.go) and the animation
// tier reaches the leg exactly once. The write leg underneath the door (the
// ApplySetting case, the pump-source wiring and the end-to-end snapshot proof)
// is tested for real in cmd/wisp (panel_animation_wiring_193_windows_test.go);
// the config-side setter/reader is tested in internal/config
// (settings_animation_193_test.go). These tests pin the door half:
//
//   - the roster joins the union WITHOUT moving any older pin (the frozen 248
//     table and every older roster's exact want-list keep meaning exactly what
//     they meant);
//   - a config.set naming panel_animation reaches the leg - the tooth that
//     goes red if the union admission (knownWritableField's writableField193
//     term) is ever dropped, which is the mutation this closing leg exists to
//     prevent (the 65 roster's tooth shape, ball_style_65_test.go).
//
// Every judgement is a call count or an error identity, never a log string -
// the same counting rule config_route_248_test.go states for the doors these
// tests share.

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

// Test193RosterJoinsTheUnionWithoutMovingOlderPins: the ninth roster is
// admitted by the union and by both helper unions (the refusal text and the
// empty-value check), and no older roster grew to carry the 193 field.
func Test193RosterJoinsTheUnionWithoutMovingOlderPins(t *testing.T) {
	if got := WritableFields193(); !reflect.DeepEqual(got, []string{FieldPanelAnimation}) {
		t.Fatalf("WritableFields193() = %v, want exactly [panel_animation]", got)
	}
	if !knownWritableField(FieldPanelAnimation) {
		t.Fatalf("knownWritableField refuses %q - the union admission is gone and the settings door would refuse the animation switch on a live machine", FieldPanelAnimation)
	}
	for _, older := range []struct {
		name string
		rows []string
	}{
		{"248", WritableFields()},
		{"313", AdditionalWritableFields()},
		{"187", WritableFields187()},
		{"65", WritableFields65()},
		{"wave", WritableFieldsWave()},
		{"215a", WritableFields215()},
		{"39", WritableFields39()},
		{"40quota", WritableFields40Quota()},
	} {
		if containsString(older.rows, FieldPanelAnimation) {
			t.Errorf("the %s roster carries %q - rosters are disjoint by construction and the 193 field must live in its own (config_write_193.go)", older.name, FieldPanelAnimation)
		}
	}
	if !containsString(acceptedWritableFields(), FieldPanelAnimation) {
		t.Errorf("acceptedWritableFields() = %v, want it to carry the 193 field: the refusal text would send the page chasing the wrong list", acceptedWritableFields())
	}
	if !takesValueFor(FieldPanelAnimation) {
		t.Errorf("takesValueFor(%q) = false: an empty value would slip past the door's value check and surface as a config-side enum refusal in the wrong shape", FieldPanelAnimation)
	}
	if needsProviderFor(FieldPanelAnimation) || needsModelFor(FieldPanelAnimation) || needsPluginFor(FieldPanelAnimation) {
		t.Errorf("panel_animation takes the plain value half; it must not demand a provider, model or plugin selector")
	}
}

// Test193DoorAdmitsPanelAnimationAndRoutesItToTheLeg: config.set naming
// panel_animation clears the door and reaches the leg exactly once with the
// tier word - the MUTATION TOOTH. Drop `|| writableField193(field)` from
// knownWritableField and this test goes red on the refusal error (the door
// answers "不在可写名册内" and the leg's apply count stays 0). Drop the field
// from acceptedWritableFields (the refusal text) and the door still routes -
// but the page is told a list without the name, which the first test's union
// check catches.
func Test193DoorAdmitsPanelAnimationAndRoutesItToTheLeg(t *testing.T) {
	leg := &cfg248Leg{res: SettingWriteResult{
		Field: FieldPanelAnimation, Written: []string{"panel.animation"}, Tier: EffectiveNow,
	}}
	reply, err := cfg248Dispatch(leg).Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigSet, "requestId": "r-193", "source": ComposerRequestSource,
		"configField": FieldPanelAnimation, "configValue": "off",
	}))
	if err != nil {
		t.Fatalf("the door refused panel_animation: %v", err)
	}
	if leg.applyCalls != 1 {
		t.Fatalf("leg apply calls = %d, want 1", leg.applyCalls)
	}
	if leg.lastWrite.Field != FieldPanelAnimation || leg.lastWrite.Value != "off" {
		t.Fatalf("leg saw %+v, want field panel_animation with value off", leg.lastWrite)
	}
	if !strings.Contains(reply, "panel.animation") {
		t.Errorf("the receipt never named the key path: %q", reply)
	}

	// Positive control for the door's value check (the takesValueFor seam): an
	// empty value is refused BEFORE the leg, never written as a default.
	empty := &cfg248Leg{}
	if _, err := cfg248Dispatch(empty).Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigSet, "requestId": "r-193-empty", "source": ComposerRequestSource,
		"configField": FieldPanelAnimation, "configValue": "   ",
	})); err == nil {
		t.Fatalf("an empty panel_animation value must be refused at the door")
	} else if empty.applyCalls != 0 {
		t.Errorf("the refused write reached the leg %d times, want 0", empty.applyCalls)
	}

	// The negative control the positive one is judged against: a name no
	// roster carries is still refused, and the refusal is the door's own.
	// (The dotted spelling panel.animation is also refused: the roster admits
	// the FLAT field name only, the same rule plugin_enabled's roster states.)
	for _, bad := range []string{"orb_animation", "panel.animation"} {
		unlisted := &cfg248Leg{}
		if _, err := cfg248Dispatch(unlisted).Handle(context.Background(), cfg248Raw(t, map[string]any{
			"method": MethodConfigSet, "requestId": "r-193-bad", "source": ComposerRequestSource,
			"configField": bad, "configValue": "off",
		})); err == nil {
			t.Fatalf("the name %q no roster carries must be refused", bad)
		} else if unlisted.applyCalls != 0 {
			t.Errorf("the unlisted write %q reached the leg %d times, want 0", bad, unlisted.applyCalls)
		}
	}
}
