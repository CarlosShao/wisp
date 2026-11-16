package panel

// Ticket 215a's door tests: the settings DOOR now admits plugin_enabled with
// its own selector, and the write leg underneath the door (the direction
// verdict, the audit, the receipt) is tested for real in internal/config
// (settings_plugins_215_test.go) and cmd/wisp; these tests pin the panel half:
//
//   - the roster joins the union WITHOUT moving any older pin (the frozen 248
//     table, the 313/187/65/wave rosters and the disjointness between them
//     keep meaning exactly what they meant);
//   - a config.set naming plugin_enabled reaches the leg exactly once, with
//     the plugin id carried on the envelope's configPlugin selector - the
//     tooth that goes red if the union admission
//     (knownWritableField's writableField215 term) is ever dropped;
//   - the selector is REQUIRED at the door: a plugin-shaped write with no
//     plugin id is refused before the leg, never guessed;
//   - the dotted spelling ("plugins.clip.enabled") is still refused by the
//     locked-family guard with its own security sentence - the flat name
//     admits one guarded write, never a key path.

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

// Test215ARosterJoinsTheUnionWithoutMovingOlderPins: the sixth roster is
// admitted by the union and by the helper unions (the refusal text, the
// value check and the new plugin-selector check), and none of the older
// rosters grew to carry it.
func Test215ARosterJoinsTheUnionWithoutMovingOlderPins(t *testing.T) {
	if got := WritableFields215(); !reflect.DeepEqual(got, []string{FieldPluginEnabled}) {
		t.Fatalf("WritableFields215() = %v, want exactly [plugin_enabled]", got)
	}
	if !knownWritableField(FieldPluginEnabled) {
		t.Fatalf("knownWritableField refuses %q - the union admission is gone and the plugin switch would die at the door on a live machine", FieldPluginEnabled)
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
	} {
		if containsString(older.rows, FieldPluginEnabled) {
			t.Errorf("the %s roster carries %q - rosters are disjoint by construction and the 215a field must live in its own (config_write_215.go)", older.name, FieldPluginEnabled)
		}
	}
	if !containsString(acceptedWritableFields(), FieldPluginEnabled) {
		t.Errorf("acceptedWritableFields() = %v, want it to carry the 215a field: the refusal text would send the page chasing the wrong list", acceptedWritableFields())
	}
	if !takesValueFor(FieldPluginEnabled) {
		t.Errorf("takesValueFor(%q) = false: an empty value would slip past the door's value check", FieldPluginEnabled)
	}
	if !needsPluginFor(FieldPluginEnabled) {
		t.Errorf("needsPluginFor(%q) = false: the door would hand the leg a write with no plugin id to address", FieldPluginEnabled)
	}
	if needsProviderFor(FieldPluginEnabled) || needsModelFor(FieldPluginEnabled) {
		t.Errorf("plugin_enabled takes the plugin selector, not the provider/model selectors")
	}
}

// Test215ADoorAdmitsPluginEnabledAndRoutesItToTheLeg: config.set naming
// plugin_enabled with a configPlugin selector clears the door and reaches
// the leg exactly once with the id and the value - the MUTATION TOOTH. Drop
// `|| writableField215(field)` from knownWritableField and this test goes
// red on the refusal error.
func Test215ADoorAdmitsPluginEnabledAndRoutesItToTheLeg(t *testing.T) {
	leg := &cfg248Leg{res: SettingWriteResult{
		Field: FieldPluginEnabled, Written: []string{"plugins.clip.enabled"}, Tier: EffectiveNow,
	}}
	reply, err := cfg248Dispatch(leg).Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigSet, "requestId": "r-215a", "source": ComposerRequestSource,
		"configField": FieldPluginEnabled, "configPlugin": "clip", "configValue": "false",
	}))
	if err != nil {
		t.Fatalf("the door refused plugin_enabled: %v", err)
	}
	if leg.applyCalls != 1 {
		t.Fatalf("leg apply calls = %d, want 1", leg.applyCalls)
	}
	if leg.lastWrite.Field != FieldPluginEnabled || leg.lastWrite.Plugin != "clip" || leg.lastWrite.Value != "false" {
		t.Fatalf("leg saw %+v, want plugin_enabled addressed to clip with value false", leg.lastWrite)
	}
	if !strings.Contains(reply, "plugins.clip.enabled") {
		t.Errorf("the receipt never named the key path: %q", reply)
	}

	// The selector is required at the door: no plugin id, no write - and the
	// refusal names the selector so the page knows which half it dropped.
	empty := &cfg248Leg{}
	_, err = cfg248Dispatch(empty).Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigSet, "requestId": "r-215a-noid", "source": ComposerRequestSource,
		"configField": FieldPluginEnabled, "configValue": "false",
	}))
	if err == nil {
		t.Fatalf("a plugin_enabled write with no configPlugin must be refused at the door")
	}
	if !strings.Contains(err.Error(), "configPlugin") {
		t.Errorf("the no-id refusal does not name the selector: %v", err)
	}
	if empty.applyCalls != 0 {
		t.Errorf("the refused write reached the leg %d times, want 0", empty.applyCalls)
	}

	// The negative control: a name no roster carries is still refused, and
	// the refusal is the door's own.
	unlisted := &cfg248Leg{}
	if _, err := cfg248Dispatch(unlisted).Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigSet, "requestId": "r-215a-bad", "source": ComposerRequestSource,
		"configField": "plugin_uninstall", "configPlugin": "clip",
	})); err == nil {
		t.Fatalf("a field no roster carries must be refused")
	} else if unlisted.applyCalls != 0 {
		t.Errorf("the unlisted write reached the leg %d times, want 0", unlisted.applyCalls)
	}
}

// The family guard holds: the DOTTED spelling of the same key is refused by
// the locked-family guard with its security sentence, before the roster
// union is even consulted. The flat field name admits exactly one guarded
// write; it never opens the family.
func Test215ADottedPluginSpellingIsStillRefusedByTheFamilyGuard(t *testing.T) {
	leg := &cfg248Leg{}
	_, err := cfg248Dispatch(leg).Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigSet, "requestId": "r-215a-dotted", "source": ComposerRequestSource,
		"configField": "plugins.clip.enabled", "configValue": "true",
	}))
	if err == nil {
		t.Fatalf("the dotted plugins.* spelling must stay refused by the family guard")
	}
	if leg.applyCalls != 0 {
		t.Errorf("the family-refused write reached the leg %d times, want 0", leg.applyCalls)
	}
	family, why := refusedLockedFamily("plugins.clip.enabled")
	if family != "plugins." || why == "" {
		t.Errorf("refusedLockedFamily = (%q, %q), want the plugins family with its security sentence", family, why)
	}
}
