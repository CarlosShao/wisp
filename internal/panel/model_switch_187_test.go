package panel

// Ticket 187's carrier tests: the modelSwitch key is a reader-gated section
// like commands (213's pattern, no second paradigm), so its tests pin the
// same two directions that section's tests pin:
//
//   - WITHOUT the reader the key is absent from the bytes a real pump
//     marshals (the four-key shape the byte nails pin, with this key named);
//   - WITH the reader the key always arrives, values travel verbatim, nil
//     lists render as [] and never null, and the no-offerable-rows shape
//     (Reason set, rows empty) arrives as written instead of being dressed
//     down to an empty menu (AC#4's 宁缺毋造).
//
// A mutation that drops the pump wiring, drops a field from the carrier, or
// fills the key without a reader reds here.

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestTheModelSwitchKeyTravelsOnlyFromItsReader187(t *testing.T) {
	// (1) No reader, no key: the four-key shape the byte nails pin, with this
	// key named so a regression says which key came back.
	data, err := NewSnapshotPump(PumpSources{}).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), `"modelSwitch"`) {
		t.Fatalf("a pump with no modelSwitch reader sent the key: %s", data)
	}
	var generic map[string]json.RawMessage
	if err := json.Unmarshal(data, &generic); err != nil {
		t.Fatal(err)
	}
	if len(generic) != 4 {
		t.Errorf("top-level keys = %d, want the four pump_test pins for a reader-less pump", len(generic))
	}

	// (2) With a reader the key arrives and every field travels verbatim.
	want := ModelSwitchSection{
		Current:         ModelOption{Provider: "deepseek", Model: "deepseek-chat"},
		CurrentThinking: "low",
		Models: []ModelOption{
			{Provider: "deepseek", Model: "deepseek-chat"},
			{Provider: "stepfun", Model: "step-3"},
		},
		ThinkingLevels: []string{"off", "low"},
	}
	pump := NewSnapshotPump(PumpSources{
		ModelSwitch: func() ModelSwitchSection { return want },
	})
	data, err = pump.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &generic); err != nil {
		t.Fatal(err)
	}
	raw, ok := generic["modelSwitch"]
	if !ok {
		t.Fatalf("a pump with a modelSwitch reader sent no modelSwitch key: %s", data)
	}
	var got ModelSwitchSection
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("wire section = %+v, want %+v (every field travels verbatim, or the picker lies about what is switchable)", got, want)
	}

	// (3) Nil lists from a reader render as [] - the page reads .length on
	// these like it does on pending/results, and null breaks that.
	pump = NewSnapshotPump(PumpSources{
		ModelSwitch: func() ModelSwitchSection { return ModelSwitchSection{} },
	})
	data, err = pump.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"models":[]`) || !strings.Contains(string(data), `"thinkingLevels":[]`) {
		t.Errorf("nil lists did not render as []: %s", data)
	}

	// (4) The no-offerable-rows shape travels as written: the explicit
	// sentence and empty rows, never a silently empty menu (AC#4).
	empty := ModelSwitchSection{
		Models: []ModelOption{},
		Reason: "配置里没有任何 [llm.providers] 节：没有可用模型（行在配置文件里建，面板不替你造）",
	}
	pump = NewSnapshotPump(PumpSources{
		ModelSwitch: func() ModelSwitchSection { return empty },
	})
	data, err = pump.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &generic); err != nil {
		t.Fatal(err)
	}
	raw, ok = generic["modelSwitch"]
	if !ok {
		t.Fatalf("the empty-catalog shape lost its key: %s", data)
	}
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if got.Reason != empty.Reason || len(got.Models) != 0 {
		t.Errorf("empty-catalog shape = %+v, want the sentence carried verbatim with empty rows", got)
	}
}

// AC#3's panel-side nail, read off the bytes: the two files this ticket adds
// to the panel carry no permission surface. The model/tier dimension and the
// risk mode are unrelated dimensions (R20's 口径 is already nailed elsewhere);
// the nail keeps the second one from ever arriving through the first's door -
// no import, no Mode vocabulary, no permission field name.
func Test187PanelFilesTouchNoPermissionSurface(t *testing.T) {
	for _, name := range []string{"model_switch_187.go", "config_write_187.go"} {
		// Tests run with the package dir as the working directory; the two
		// files live beside this test.
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		src := string(data)
		for _, banned := range []string{"internal/perm", "internal/risk", "risk.Mode", "Modes.Set", "permission_mode"} {
			if strings.Contains(src, banned) {
				t.Errorf("%s mentions %q - the model/tier write face must stay unrelated to the permission dimension (AC#3)", name, banned)
			}
		}
	}
}

// The field name and the section's JSON keys carry no approval vocabulary -
// the same naming ban the 248 pin enforces on its own roster, applied to the
// 187 additions.
func Test187NamesCarryNoApprovalVocabulary(t *testing.T) {
	grantWords := []string{
		"approve", "approval", "grant", "allow", "permit",
		"ratify", "authorize", "authorised", "decide", "decision", "verdict",
	}
	for _, w := range grantWords {
		if strings.Contains(strings.ToLower(FieldRoleChatThinking), w) {
			t.Errorf("writable field %q carries the approval word %q", FieldRoleChatThinking, w)
		}
	}
	for _, key := range []string{"modelSwitch", "currentThinking", "thinkingLevels"} {
		for _, w := range grantWords {
			if strings.Contains(strings.ToLower(key), w) {
				t.Errorf("section key %q carries the approval word %q", key, w)
			}
		}
	}
}
