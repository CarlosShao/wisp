package panel

// Ticket 187's gate tests: role_chat_thinking rides the settings route
// ticket 248 opened, through the third roster (config_write_187.go), with the
// frozen 248 pin and the 313 roster untouched.
//
// Every judgement is a call count, an error identity or the receipt's own
// sentence - never a log string (the rule config_route_248_test.go states).
// The mutation teeth:
//
//	- drop FieldRoleChatThinking from fields187 -> the gate refuses and
//	  Test187RoleChatThinkingClearsTheGate reds;
//	- drop the tier correction in config_handlers.go's write() -> the receipt
//	  falls back to the leg's blanket restart stamp and
//	  Test187ThinkingReceiptNamesTheMeasuredTier reds (that is the wrong-tier
//	  lie AC#8 exists to prevent: the tier is next_task by measurement).

import (
	"context"
	"strings"
	"testing"
)

// The new field clears the door and reaches the leg exactly once, with the
// value carried; an empty value is refused before the leg (the gate's own
// empty-value rule, already table-driven - this pins it for the 187 roster).
func Test187RoleChatThinkingClearsTheGate(t *testing.T) {
	leg := &cfg248Leg{res: SettingWriteResult{
		Field: FieldRoleChatThinking, Written: []string{"llm.roles.chat.thinking_intensity"},
		Tier: EffectiveRestart, // the leg's blanket stamp; the door corrects it
	}}
	env := map[string]any{
		"method": MethodConfigSet, "requestId": "r-187-think", "source": ComposerRequestSource,
		"configField": FieldRoleChatThinking, "configValue": "low",
	}
	if _, err := cfg248Dispatch(leg).Handle(context.Background(), cfg248Raw(t, env)); err != nil {
		t.Fatalf("field %q was refused by the door: %v", FieldRoleChatThinking, err)
	}
	if leg.applyCalls != 1 {
		t.Errorf("leg apply calls = %d, want 1", leg.applyCalls)
	}
	if leg.lastWrite.Field != FieldRoleChatThinking || leg.lastWrite.Value != "low" {
		t.Errorf("the leg received %+v", leg.lastWrite)
	}

	// Empty value: refused at the gate, the leg never runs.
	leg = &cfg248Leg{}
	env["configValue"] = "   "
	if _, err := cfg248Dispatch(leg).Handle(context.Background(), cfg248Raw(t, env)); err == nil {
		t.Error("an empty tier value was accepted")
	}
	if leg.applyCalls != 0 {
		t.Errorf("an empty-value write reached the leg %d time(s)", leg.applyCalls)
	}
}

// The tier correction: the leg stamps restart on every [llm] write, but the
// thinking tier's measured tier is next_task (the per-task re-resolution over
// the live config), so the receipt must say the next-task sentence. The
// correction applies to the 187 roster's names only - the control proves an
// older field's receipt keeps the leg's own stamp.
func Test187ThinkingReceiptNamesTheMeasuredTier(t *testing.T) {
	leg := &cfg248Leg{res: SettingWriteResult{
		Field: FieldRoleChatThinking, Written: []string{"llm.roles.chat.thinking_intensity"},
		Tier: EffectiveRestart,
	}}
	reply, err := cfg248Dispatch(leg).Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigSet, "requestId": "r-187-tier", "source": ComposerRequestSource,
		"configField": FieldRoleChatThinking, "configValue": "high",
	}))
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	if !strings.Contains(reply, tierSentence(EffectiveNextTask)) {
		t.Errorf("the receipt lost the measured next-task tier: %q", reply)
	}
	if strings.Contains(reply, tierSentence(EffectiveRestart)) {
		t.Errorf("the receipt claims a restart for a next-task field (the wrong-tier lie AC#8 forbids): %q", reply)
	}
	if !strings.Contains(reply, "llm.roles.chat.thinking_intensity") {
		t.Errorf("the receipt flattened the key path instead of naming it: %q", reply)
	}

	// Control: a 248-roster field keeps the leg's own tier stamp untouched.
	leg2 := &cfg248Leg{res: SettingWriteResult{
		Field: FieldRoleChatModel, Written: []string{"llm.roles.chat.model"},
		Tier: EffectiveRestart,
	}}
	reply, err = cfg248Dispatch(leg2).Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigSet, "requestId": "r-187-ctrl", "source": ComposerRequestSource,
		"configField": FieldRoleChatModel, "configValue": "deepseek-chat",
	}))
	if err != nil {
		t.Fatalf("control Handle: %v", err)
	}
	if !strings.Contains(reply, tierSentence(EffectiveRestart)) {
		t.Errorf("the control field's receipt lost the leg's restart stamp: %q", reply)
	}
}

// The rosters stay disjoint and the frozen pins keep their meaning: the 248
// seven, the 313 roster, and the 187 roster share no name, and the door asks
// all three (a field added to a roster but not to the union would be
// invisible to the refusal text's own list).
func Test187RosterIsDisjointAndTheUnionAsksAllThree(t *testing.T) {
	listed := map[string]bool{}
	for _, f := range WritableFields() {
		listed[f] = true
	}
	additional := map[string]bool{}
	for _, f := range AdditionalWritableFields() {
		additional[f] = true
	}
	for _, f := range WritableFields187() {
		if listed[f] || additional[f] {
			t.Errorf("%q is in more than one roster: the 187 addition must not restate an older name", f)
		}
	}
	if !knownWritableField(FieldRoleChatThinking) {
		t.Fatal("the door does not ask the 187 roster: the field would be invisible to the gate")
	}
	if !takesValueFor(FieldRoleChatThinking) {
		t.Error("the door does not ask the 187 value column: the empty-value refusal would not fire")
	}
	if needsProviderFor(FieldRoleChatThinking) || needsModelFor(FieldRoleChatThinking) {
		t.Error("the tier field needs no provider/model selector; the gate would chase one that is not there")
	}
	for _, f := range acceptedWritableFields() {
		if !knownWritableField(f) {
			t.Errorf("the refusal text's own list carries %q, which the door would refuse", f)
		}
	}
	if !containsString(acceptedWritableFields(), FieldRoleChatThinking) {
		t.Errorf("acceptedWritableFields() = %v, want it to carry the 187 field: the refusal text would send the page chasing the wrong list", acceptedWritableFields())
	}
}
