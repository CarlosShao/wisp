package panel

// Ticket 39 AC#1's tier-events teeth at the settings door (2026-11-18).
// The config layer owns the tier BEHAVIOR (internal/config's tier-aware write
// path) and cmd/wisp owns the route teeth through the real leg; this file pins
// the door's half:
//
//   - the seventh roster: app_language is a known, value-taking, non-locked
//     field, and the refusal text's union carries it (a page told 可写的只有
//     ... must not be sent chasing a list that lacks the field it just used);
//   - the route hands the leg the field and the value, and an empty value is
//     refused before the leg is reached;
//   - the receipt/event wording: a restart-tier receipt carries 需要重启
//     VERBATIM (the page-visible restart-required line this ticket's AC
//     names) and never claims an immediate effect; the reload event's tier
//     field carries the same word the receipt does.
//
// The tier WORDS here are read from the same source the production leg asks
// (internal/config's TierRegistry via the leg's own mapping, asserted by the
// cmd/wisp teeth); this file only pins that the door says what it is handed.

import (
	"context"
	"strings"
	"testing"
)

// stub39Leg captures what the gate handed to the leg and answers a fixed
// result, so the route half is judged without a config leg at all.
type tier39Leg struct {
	got   SettingWrite
	res   SettingWriteResult
	err   error
	calls int
}

func (l *tier39Leg) ReadSettings(ctx context.Context) (SettingsView, error) {
	return SettingsView{Readable: true}, nil
}

func (l *tier39Leg) ApplySetting(ctx context.Context, w SettingWrite) (SettingWriteResult, error) {
	l.calls++
	l.got = w
	return l.res, l.err
}

func (l *tier39Leg) StoreCredential(ctx context.Context, raw string) (SettingWriteResult, error) {
	return SettingWriteResult{}, nil
}

// Test39AppLanguageFieldIsOnTheDoor: the roster, the value column, the family
// guard, and the refusal text's union.
func Test39AppLanguageFieldIsOnTheDoor(t *testing.T) {
	if !writableField39(FieldAppLanguage) {
		t.Fatal("writableField39(app_language) = false: the roster does not carry its own field")
	}
	if !knownWritableField(FieldAppLanguage) {
		t.Fatal("knownWritableField(app_language) = false: the door's union dropped the 39 roster")
	}
	if !takesValueFor(FieldAppLanguage) {
		t.Fatal("takesValueFor(app_language) = false: the gate would refuse a legitimate save as value-less")
	}
	if family, _ := refusedLockedFamily(FieldAppLanguage); family != "" {
		t.Fatalf("refusedLockedFamily(app_language) = %q: [app] is not a locked family", family)
	}
	if !containsString(acceptedWritableFields(), FieldAppLanguage) {
		t.Fatalf("acceptedWritableFields() lacks app_language: the refusal text would send the page chasing the wrong list")
	}
}

// Test39AppLanguageRouteHandsTheFieldAndRefusesTheEmptyValue.
func Test39AppLanguageRouteHandsTheFieldAndRefusesTheEmptyValue(t *testing.T) {
	leg := &tier39Leg{res: SettingWriteResult{
		Field: FieldAppLanguage, Written: []string{"app.language"}, Tier: EffectiveRestart,
	}}
	h := &ConfigWriteHandler{Store: leg, Actor: "test-39"}

	reply, err := h.HandleConfigRequest(context.Background(), ComposerRequest{
		Method: MethodConfigSet, RequestID: "rid-app-lang-39", Source: ComposerRequestSource,
		ConfigField: FieldAppLanguage, ConfigValue: "en-US",
	}, "")
	if err != nil {
		t.Fatalf("config.set app_language was refused: %v", err)
	}
	if leg.calls != 1 || leg.got.Field != FieldAppLanguage || leg.got.Value != "en-US" {
		t.Fatalf("leg got %+v (calls=%d), want one call with field app_language and value en-US", leg.got, leg.calls)
	}
	if !strings.Contains(reply, "设置已写入") {
		t.Fatalf("the accepted save's reply is not a receipt:\n%s", reply)
	}

	// The empty value is refused before the leg is reached (no second call).
	if _, err := h.HandleConfigRequest(context.Background(), ComposerRequest{
		Method: MethodConfigSet, RequestID: "rid-app-lang-39b", Source: ComposerRequestSource,
		ConfigField: FieldAppLanguage, ConfigValue: "   ",
	}, ""); err == nil {
		t.Fatal("an empty app_language save was accepted: the gate must refuse it")
	} else if !strings.Contains(err.Error(), "需要一个值") {
		t.Fatalf("the empty-value refusal does not name the rule: %v", err)
	}
	if leg.calls != 1 {
		t.Fatalf("the refused save reached the leg (calls=%d): the gate did not stop it", leg.calls)
	}
}

// Test39RestartReceiptSaysRequiresRestartVerbatim: AC#8's wording rule at the
// word level this ticket adds - the restart sentence carries 需要重启 verbatim
// and no branch claims a save is already in force.
func Test39RestartReceiptSaysRequiresRestartVerbatim(t *testing.T) {
	sentence := tierSentence(EffectiveRestart)
	if !strings.Contains(sentence, "需要重启") {
		t.Fatalf("the restart tier sentence does not carry 需要重启 verbatim: %q", sentence)
	}
	line := renderSettingReceipt(SettingWriteResult{
		Written: []string{"app.language"}, Tier: EffectiveRestart,
	})
	if !strings.Contains(line, sentence) {
		t.Fatalf("the restart receipt does not carry the tier sentence verbatim:\n%s", line)
	}
	if strings.Contains(line, "立即生效") || strings.Contains(line, "保存即生效") {
		t.Fatalf("a restart receipt claims the save is in force:\n%s", line)
	}
	// The event's tier field carries the same word the receipt does.
	ev := renderReloadEvent(SettingWriteResult{
		Written: []string{"app.language"}, Tier: EffectiveRestart,
	})
	if !strings.HasPrefix(ev, ConfigReloadEventMarker) || !strings.Contains(ev, `"tier":"restart"`) {
		t.Fatalf("the reload event does not carry the restart tier:\n%s", ev)
	}
	// And the hot tier's sentence stays free of restart wording (the other
	// half of the same tooth: a hot receipt must not say 需要重启).
	now := tierSentence(EffectiveNow)
	if strings.Contains(now, "重启") {
		t.Fatalf("the now tier sentence carries restart wording: %q", now)
	}
}
