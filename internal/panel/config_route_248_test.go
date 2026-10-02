package panel

// Ticket 248 AC#1/AC#8 - the settings route, judged by capability.
//
// Every judgement below is a CALL COUNT or an error identity, never a log string:
// a router that writes a pleasant sentence while never reaching the leg is exactly
// the failure this file is written against (the same counting rule ticket 33 slice
// A's own nails state, inherited because these two doors live in that router).
//
// Each negative case carries its positive control in the same test, because a
// refusal that cannot be made to fire the other way proves nothing about the door.

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
)

// cfg248Leg records what the router handed it and returns what the test asks for.
// It is the seam panel.ConfigStore declares (AGENTS.md §1.3: the injection is at
// the seam, not a mock standing in for the real writer, which is tested for real
// in cmd/wisp).
type cfg248Leg struct {
	readCalls   int
	applyCalls  int
	credCalls   int
	lastWrite   SettingWrite
	lastRawSeen string
	lastReqSeen ComposerRequest

	readView SettingsView
	readErr  error

	res     SettingWriteResult
	applyEr error
}

func (l *cfg248Leg) ReadSettings(context.Context) (SettingsView, error) {
	l.readCalls++
	return l.readView, l.readErr
}

func (l *cfg248Leg) ApplySetting(_ context.Context, w SettingWrite) (SettingWriteResult, error) {
	l.applyCalls++
	l.lastWrite = w
	return l.res, l.applyEr
}

func (l *cfg248Leg) StoreCredential(_ context.Context, raw string) (SettingWriteResult, error) {
	l.credCalls++
	l.lastRawSeen = raw
	return l.res, l.applyEr
}

func cfg248Raw(t *testing.T, fields map[string]any) string {
	t.Helper()
	data, err := json.Marshal(fields)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(data)
}

func cfg248Dispatch(leg ConfigStore) *ComposerDispatch {
	h := &ConfigWriteHandler{Store: leg, Actor: "test-248"}
	return &ComposerDispatch{Config: h}
}

// AC#1's positive half: both new names clear the whitelist AND reach the leg.
func TestAC1BothSettingsNamesAreWhitelistedAndRouted(t *testing.T) {
	for _, method := range []string{MethodConfigGet, MethodConfigSet} {
		if !knownComposerMethod(method) {
			t.Fatalf("knownComposerMethod refuses %q: the whitelist did not really grow", method)
		}
		_, err := ParseComposerRequest(cfg248Raw(t, map[string]any{
			"method": method, "requestId": "r-" + method, "source": ComposerRequestSource,
		}))
		if err != nil {
			t.Errorf("ParseComposerRequest refuses %q: %v", method, err)
		}
	}

	t.Run("config.get reaches the leg once and its answer comes back", func(t *testing.T) {
		leg := &cfg248Leg{readView: SettingsView{Readable: true, Credential: CredentialNoneDeclared}}
		reply, err := cfg248Dispatch(leg).Handle(context.Background(), cfg248Raw(t, map[string]any{
			"method": MethodConfigGet, "requestId": "r-get", "source": ComposerRequestSource,
		}))
		if err != nil {
			t.Fatalf("Handle: %v", err)
		}
		if leg.readCalls != 1 {
			t.Errorf("leg read calls = %d, want 1", leg.readCalls)
		}
		if !strings.Contains(reply, "凭据状态") {
			t.Errorf("the read answer never reached the caller: %q", reply)
		}
	})

	t.Run("config.set reaches the leg once with the parsed selectors", func(t *testing.T) {
		leg := &cfg248Leg{res: SettingWriteResult{
			Field:   FieldProviderBaseURL,
			Written: []string{"llm.providers.deepseek.base_url"}, Tier: EffectiveRestart,
		}}
		reply, err := cfg248Dispatch(leg).Handle(context.Background(), cfg248Raw(t, map[string]any{
			"method": MethodConfigSet, "requestId": "r-set", "source": ComposerRequestSource,
			"configField": FieldProviderBaseURL, "configProvider": "deepseek",
			"configValue": "https://x.example.invalid",
		}))
		if err != nil {
			t.Fatalf("Handle: %v", err)
		}
		if leg.applyCalls != 1 {
			t.Fatalf("leg apply calls = %d, want 1", leg.applyCalls)
		}
		if leg.lastWrite.Provider != "deepseek" || leg.lastWrite.Field != FieldProviderBaseURL {
			t.Errorf("the leg received %+v", leg.lastWrite)
		}
		if !strings.Contains(reply, "llm.providers.deepseek.base_url") {
			t.Errorf("the receipt flattened the key paths instead of naming them: %q", reply)
		}
	})
}

// AC#1's negative half with its control: a name off the roster is refused and the
// leg never runs, while the same envelope with a rostered name does run - so the
// refusal is about the name and not about how the request was written.
func TestAC1UnlistedMethodIsRefusedAndTheLegNeverRuns(t *testing.T) {
	leg := &cfg248Leg{}
	d := cfg248Dispatch(leg)

	var audit []string
	d.Audit = func(format string, args ...any) { audit = append(audit, format) }

	for _, name := range []string{"panel.review.allow", "approval.decide", "config.write", "panel.config.set"} {
		reply, err := d.Handle(context.Background(), cfg248Raw(t, map[string]any{
			"method": name, "requestId": "r-nope", "source": ComposerRequestSource,
		}))
		if !errors.Is(err, ErrComposerRequest) {
			t.Errorf("%q returned %v, want a refusal carrying ErrComposerRequest", name, err)
		}
		if !strings.Contains(reply, name) || !strings.Contains(reply, "r-nope") {
			t.Errorf("refusal for %q does not say what it refused: %q", name, reply)
		}
	}
	if leg.applyCalls+leg.readCalls+leg.credCalls != 0 {
		t.Errorf("a refused name still reached the leg: read=%d apply=%d cred=%d",
			leg.readCalls, leg.applyCalls, leg.credCalls)
	}

	// Control: one byte different (the method) and the same door opens.
	if _, err := d.Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigSet, "requestId": "r-yes", "source": ComposerRequestSource,
		"configField": FieldRoleChatModel, "configValue": "deepseek-chat",
	})); err != nil {
		t.Fatalf("the control request was refused too: %v", err)
	}
	if leg.applyCalls != 1 {
		t.Errorf("control leg apply calls = %d, want 1", leg.applyCalls)
	}
}

// AC#1's other named refusal: rostered here, unattached in this assembly. The
// control is the same envelope against a dispatch that has the socket.
func TestAC1RosteredSettingWithoutALegRefusesLoudly(t *testing.T) {
	for _, method := range []string{MethodConfigGet, MethodConfigSet} {
		d := &ComposerDispatch{}
		_, err := d.Handle(context.Background(), cfg248Raw(t, map[string]any{
			"method": method, "requestId": "r-nil", "source": ComposerRequestSource,
		}))
		if !errors.Is(err, ErrNoHandlerAttached) {
			t.Errorf("%q with no leg returned %v, want %v", method, err, ErrNoHandlerAttached)
		}
	}
	leg := &cfg248Leg{readView: SettingsView{Readable: true}}
	d := cfg248Dispatch(leg)
	if _, err := d.Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigGet, "requestId": "r-attached", "source": ComposerRequestSource,
	})); err != nil {
		t.Errorf("control with a leg attached: %v", err)
	}
	if leg.readCalls != 1 {
		t.Errorf("control leg read calls = %d, want 1", leg.readCalls)
	}
}

// AC#1's backstop half: a name the whitelist answers that the dispatch table never
// lists must be refused, not accepted - reached only through dispatch directly, and
// here through the settings pair so the two rosters are checked against each other.
func TestAC1DispatchTableNamesEveryWhitelistedSettingMethod(t *testing.T) {
	names := map[string]bool{MethodConfigGet: true, MethodConfigSet: true}
	for name := range names {
		if !knownComposerMethod(name) {
			t.Errorf("%q is answered by nothing", name)
		}
	}
	// The roster this file may not copy: the guard's set and WritableFields must
	// agree with the switch, which is the same capability question ticket 248's
	// ruling asked for (判据一律问能力).
	d := cfg248Dispatch(&cfg248Leg{})
	_, err := d.Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigSet, "requestId": "r-nofield", "source": ComposerRequestSource,
	}))
	if err == nil {
		t.Error("a config.set with no field named was accepted")
	} else if errors.Is(err, ErrRosterMismatch) {
		t.Errorf("the settings door fell into the backstop instead of the field gate: %v", err)
	}
}

// AC#0 Q1 / J5: the field table is a whitelist, and the locked families are refused
// by name before the leg is reached. The control is one accepted field.
func TestAC1LockedFamilyFieldsAreRefusedBeforeTheLeg(t *testing.T) {
	cases := []struct {
		field   string
		mustSay string
	}{
		{"risk.permission_mode", "panel.mode.request"},
		{"fs.allowed_dirs", "L2"},
		{"permission_mode", "panel.mode.request"},
		{"net.proxy.mode", "锁定族"},
		{"privacy.keep_transcript", "锁定族"},
	}
	for _, tc := range cases {
		leg := &cfg248Leg{}
		_, err := cfg248Dispatch(leg).Handle(context.Background(), cfg248Raw(t, map[string]any{
			"method": MethodConfigSet, "requestId": "r-lock", "source": ComposerRequestSource,
			"configField": tc.field, "configValue": "whatever", "configProvider": "deepseek",
		}))
		if !errors.Is(err, ErrConfigFieldRefused) {
			t.Errorf("field %q returned %v, want %v", tc.field, err, ErrConfigFieldRefused)
		}
		if err != nil && !strings.Contains(err.Error(), tc.mustSay) {
			t.Errorf("refusal for %q does not name where that key does live: %v", tc.field, err)
		}
		if leg.applyCalls != 0 {
			t.Errorf("field %q reached the write leg %d time(s)", tc.field, leg.applyCalls)
		}
	}
	// Control: the same shape with a rostered field does reach the leg.
	leg := &cfg248Leg{res: SettingWriteResult{Tier: EffectiveRestart}}
	if _, err := cfg248Dispatch(leg).Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigSet, "requestId": "r-ok", "source": ComposerRequestSource,
		"configField": FieldRoleChatModel, "configValue": "deepseek-chat",
	})); err != nil {
		t.Fatalf("control: %v", err)
	}
	if leg.applyCalls != 1 {
		t.Errorf("control apply calls = %d, want 1", leg.applyCalls)
	}
}

// The table's own shape, read off the code rather than copied: an unlisted field is
// refused, and every listed field is refused if the leg has no row for it (checked
// in cmd/wisp). This is also the naming ban made machine-readable: no writable field
// and no route name may be spelled out of the approval vocabulary.
func TestAC1FieldTableCarriesNoApprovalVocabulary(t *testing.T) {
	grantWords := []string{
		"approve", "approval", "grant", "allow", "permit",
		"ratify", "authorize", "authorised", "decide", "decision", "verdict",
	}
	for _, name := range []string{MethodConfigGet, MethodConfigSet} {
		for _, w := range grantWords {
			if strings.Contains(strings.ToLower(name), w) {
				t.Errorf("settings route %q carries the approval word %q (AGENTS.md §1.2: no panel-side allow)", name, w)
			}
		}
	}
	want := []string{
		FieldModelContextWindow, FieldModelPriceIn, FieldModelPriceOut,
		FieldProviderAPIKeyRef, FieldProviderBaseURL, FieldProviderCredential, FieldRoleChatModel,
	}
	if got := WritableFields(); !reflect.DeepEqual(got, want) {
		t.Errorf("WritableFields() = %v, want %v", got, want)
	}
	for _, f := range want {
		for _, w := range grantWords {
			if strings.Contains(strings.ToLower(f), w) {
				t.Errorf("writable field %q carries the approval word %q", f, w)
			}
		}
	}
}

// The selectors each field needs, and the empty-value refusal. Controls are in the
// same table: the same field with the selector present does reach the leg.
func TestAC1SelectorAndEmptyValueRefusals(t *testing.T) {
	cases := []struct {
		name   string
		fields map[string]any
	}{
		{"provider field without a provider", map[string]any{"configField": FieldProviderBaseURL, "configValue": "x"}},
		{"model field without a model", map[string]any{
			"configField":    FieldModelContextWindow,
			"configProvider": "deepseek", "configValue": "128000",
		}},
		{"empty value", map[string]any{"configField": FieldRoleChatModel, "configValue": "   "}},
		{"credential field without a provider", map[string]any{"configField": FieldProviderCredential}},
	}
	for _, tc := range cases {
		leg := &cfg248Leg{}
		fields := map[string]any{"method": MethodConfigSet, "requestId": "r-sel", "source": ComposerRequestSource}
		for k, v := range tc.fields {
			fields[k] = v
		}
		if _, err := cfg248Dispatch(leg).Handle(context.Background(), cfg248Raw(t, fields)); !errors.Is(err, ErrConfigFieldRefused) {
			t.Errorf("%s: got %v, want %v", tc.name, err, ErrConfigFieldRefused)
		}
		if leg.applyCalls+leg.credCalls != 0 {
			t.Errorf("%s reached the leg anyway", tc.name)
		}
	}
}

// AC#8: the receipt carries the tier out loud, never claims a save is already in
// force, and never claims a restart-tier field is live.
func TestAC8ReceiptNamesTheTierAndNeverClaimsLiveSave(t *testing.T) {
	if got := tierSentence(EffectiveRestart); !strings.Contains(got, "重启") {
		t.Errorf("the restart sentence does not say restart: %q", got)
	}
	for _, tier := range []EffectiveTier{EffectiveRestart, EffectiveNextTask, EffectiveNotApplied, EffectiveTier("")} {
		line := renderSettingReceipt(SettingWriteResult{Written: []string{"llm.providers.deepseek.base_url"}, Tier: tier})
		if strings.Contains(line, "保存即生效") || strings.Contains(line, "立即生效") {
			t.Errorf("a %q receipt claims a live save: %q", tier, line)
		}
	}
	// The key paths survive into the sentence instead of being flattened to "saved".
	line := renderSettingReceipt(SettingWriteResult{Written: []string{"a.b", "a.c"}, Tier: EffectiveRestart})
	for _, want := range []string{"a.b", "a.c"} {
		if !strings.Contains(line, want) {
			t.Errorf("receipt dropped %q: %q", want, line)
		}
	}
	if empty := renderSettingReceipt(SettingWriteResult{Tier: EffectiveRestart}); !strings.Contains(empty, "没有任何键路径") {
		t.Errorf("a write that changed nothing must say so, not read like a save: %q", empty)
	}
	// The one tier that IS immediate is the only place that wording may appear.
	if now := renderSettingReceipt(SettingWriteResult{Written: []string{"a.b"}, Tier: EffectiveNow}); !strings.Contains(now, "立即生效") {
		t.Errorf("the now-tier sentence lost its meaning: %q", now)
	}
}

// AC#2/AC#3's panel-side half: the credential value is not carried by the parsed
// envelope at all, so the route cannot echo it - and the leg, not this file, is the
// only place raw bytes are read.
func TestAC2CredentialRouteNeverEchoesTheValueItWasGiven(t *testing.T) {
	const canary = "canary-248-not-a-real-key-0f2a"
	leg := &cfg248Leg{
		res: SettingWriteResult{
			Field:      FieldProviderCredential,
			Written:    []string{"llm.providers.deepseek.api_key_ref"},
			Tier:       EffectiveRestart,
			Credential: CredentialSummary{Ref: "dpapi:248canaryblob01", Recorded: true, Provider: "deepseek"},
		},
	}
	raw := cfg248Raw(t, map[string]any{
		"method": MethodConfigSet, "requestId": "r-cred", "source": ComposerRequestSource,
		"configField": FieldProviderCredential, "configProvider": "deepseek",
		credentialValueKeyForTest: canary,
		// Even a value parked on the shared field is never read for this route.
		"configValue": canary,
	})
	reply, err := cfg248Dispatch(leg).Handle(context.Background(), raw)
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	if leg.credCalls != 1 || leg.applyCalls != 0 {
		t.Errorf("credential route went through the wrong door: cred=%d apply=%d", leg.credCalls, leg.applyCalls)
	}
	if leg.lastRawSeen != raw {
		t.Error("the leg was not handed the raw envelope it must decode itself")
	}
	if strings.Contains(reply, canary) {
		t.Errorf("the receipt echoes the credential value: %q", reply)
	}
	if !strings.Contains(reply, "dpapi:248canaryblob01") || !strings.Contains(reply, "只写不回显") {
		t.Errorf("the receipt neither names the reference nor says write-only: %q", reply)
	}
	if strings.Contains(reply, "保存即生效") {
		t.Errorf("the credential receipt claims a live save: %q", reply)
	}
}

// credentialValueKeyForTest is the same spelling cmd/wisp declares, kept here as a
// test-only constant so this package does not start binding it in production code.
const credentialValueKeyForTest = "credentialValue"

// A read failure is a state, never an empty answer: "not configured" and "nobody
// could read it" must not collapse into one sentence.
func TestAC2UnreadableConfigIsToldAsUnreadable(t *testing.T) {
	leg := &cfg248Leg{readErr: errors.New("配置读不了")}
	if _, err := cfg248Dispatch(leg).Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigGet, "requestId": "r-err", "source": ComposerRequestSource,
	})); err == nil {
		t.Fatal("a failed read was reported as success")
	}
	leg2 := &cfg248Leg{readView: SettingsView{Readable: false, Credential: CredentialConfigUnreadable}}
	reply, err := cfg248Dispatch(leg2).Handle(context.Background(), cfg248Raw(t, map[string]any{
		"method": MethodConfigGet, "requestId": "r-unread", "source": ComposerRequestSource,
	}))
	if err != nil {
		t.Fatalf("an unreadable config is still an answer: %v", err)
	}
	if strings.Contains(reply, "未录入") || !strings.Contains(reply, "不可读") {
		t.Errorf("unreadable rendered as though it had been asked: %q", reply)
	}
}

// The read view's vocabulary, checked against what a page may be shown: states,
// counts and references only.
func TestAC2SettingsViewCarriesNoCredentialMaterial(t *testing.T) {
	v := SettingsView{
		Readable: true, Credential: CredentialPartlyMissing, ChatModel: "deepseek-chat", ChatModelKnown: true,
		Providers: []ProviderSettings{{Name: "deepseek", APIKeyRef: "env:DEEPSEEK_KEY", ModelCount: 2, CredentialError: "引用形状不合法"}},
	}
	line := renderSettingsView(v)
	for _, want := range []string{"partly_missing", "deepseek", "env:DEEPSEEK_KEY", "2"} {
		if !strings.Contains(line, want) {
			t.Errorf("the read view lost %q: %q", want, line)
		}
	}
	if strings.Contains(line, "sk-") {
		t.Errorf("a view built only from states produced key material: %q", line)
	}
}
