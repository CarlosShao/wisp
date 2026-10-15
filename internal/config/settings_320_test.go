package config

// Ticket 320's settings-write additions, judged by the file they produce.
//
// The pattern is settings_313_test.go's: every refusal asserts config.toml byte
// for byte, and every acceptance reads the file back the way the next start
// would (LoadFile). The seed is settingsBase - one preset provider with one
// catalog entry - because every write this ticket adds either addresses a row
// that exists or creates the row it then writes under.

import (
	"reflect"
	"strings"
	"testing"
)

// AC#1: the billing leaf. Both enum words land (the second one is a real change
// after the first, so this is not one accepted spelling), and both refusal
// directions - a word outside the enum and a provider this config does not
// declare - leave the file byte-identical.
func Test320SetProviderBillingWritesTheEnumLeafAndRefusesOthers(t *testing.T) {
	m, path, _ := newTestManager(t, settingsBase)

	written, err := m.SetProviderBilling("deepseek", BillingPlan)
	if err != nil {
		t.Fatalf("SetProviderBilling: %v", err)
	}
	if want := []string{"llm.providers.deepseek.billing"}; !reflect.DeepEqual(written, want) {
		t.Errorf("written key paths = %v, want %v", written, want)
	}
	if got := readLoaded(t, path).LLM.Providers["deepseek"].Billing; got != BillingPlan {
		t.Errorf("the billing mode never reached the file: %q", got)
	}

	// The other enum word, now a real change (the file carries plan above).
	if _, err := m.SetProviderBilling("deepseek", BillingPayPerToken); err != nil {
		t.Fatalf("the second enum word was refused: %v", err)
	}
	if got := readLoaded(t, path).LLM.Providers["deepseek"].Billing; got != BillingPayPerToken {
		t.Errorf("billing = %q, want %q", got, BillingPayPerToken)
	}

	// Direction one: a third spelling never lands, and the sentence names the
	// field rather than a whole-file validation failure.
	before := readSettingsFile(t, path)
	written, err = m.SetProviderBilling("deepseek", "flat-rate")
	if err == nil {
		t.Fatalf("an out-of-enum billing word was accepted: %v", written)
	}
	if written != nil {
		t.Errorf("a refused billing write reported %v, want none", written)
	}
	if !strings.Contains(err.Error(), "billing") {
		t.Errorf("the refusal does not name the field: %v", err)
	}
	if after := readSettingsFile(t, path); after != before {
		t.Errorf("config.toml changed under a refused billing write:\n%s", after)
	}

	// Direction two: a provider row this config does not declare is reason 2
	// (行不存在), and it too touches nothing.
	written, err = m.SetProviderBilling("ghostco", BillingPlan)
	if err == nil {
		t.Fatalf("a billing write for an undeclared provider was accepted: %v", written)
	}
	if !strings.Contains(err.Error(), "行不存在") {
		t.Errorf("the undeclared-provider refusal is not reason 2 (行不存在): %v", err)
	}
	if after := readSettingsFile(t, path); after != before {
		t.Errorf("config.toml changed under a refused ghost-provider write:\n%s", after)
	}
	// Control: the same call on the declared provider still works after both
	// refusals (the refusals did not wedge the writer).
	if _, err := m.SetProviderBilling("deepseek", BillingPlan); err != nil {
		t.Fatalf("the control write after two refusals: %v", err)
	}
}

// AC#2: the non-preset row. One guarded write lands the section with the protocol
// and base_url it was given plus the optional first model; five refusals - the
// four the design named and the empty name - each say why and leave the file
// byte-identical. The preset-name refusal names the step that DOES own preset
// rows, so the page is sent to the right field instead of to the file.
func Test320CreateCustomProviderRowLandsAndRefusesByName(t *testing.T) {
	m, path, _ := newTestManager(t, settingsBase)

	written, err := m.CreateCustomProviderRow("local-gw", ProtocolOpenAIChat, "http://127.0.0.1:8080/v1", "local-model-x")
	if err != nil {
		t.Fatalf("CreateCustomProviderRow: %v", err)
	}
	if want := []string{"llm.providers.local-gw"}; !reflect.DeepEqual(written, want) {
		t.Errorf("written key paths = %v, want %v", written, want)
	}
	got := readLoaded(t, path)
	row, ok := got.LLM.Providers["local-gw"]
	if !ok {
		t.Fatalf("the custom row never landed: %+v", got.LLM.Providers)
	}
	if row.Protocol != ProtocolOpenAIChat || row.BaseURL != "http://127.0.0.1:8080/v1" {
		t.Errorf("the row landed with protocol %q / base_url %q", row.Protocol, row.BaseURL)
	}
	if row.Billing != BillingPayPerToken {
		t.Errorf("the row's billing = %q, want the schema's own default: a zero value here would override the default on the next load", row.Billing)
	}
	spec, ok := row.Models["local-model-x"]
	if !ok {
		t.Fatalf("the first model never landed: %+v", row.Models)
	}
	if !spec.Enabled {
		t.Error("the created model entry is not enabled: an entry without enabled=true is treated as off and is refused at start-up")
	}
	if body := readSettingsFile(t, path); !strings.Contains(body, "[llm.providers.local-gw]") {
		t.Errorf("the section is not in the file:\n%s", body)
	}
	if got.Ball.Size != 56 {
		t.Errorf("an unrelated key moved: ball.size = %d", got.Ball.Size)
	}
	if _, ok := got.LLM.Providers["deepseek"]; !ok {
		t.Error("the pre-existing provider row disappeared")
	}

	// The model half is optional (the control): a provider with no catalog entry
	// yet is a legitimate state.
	if _, err := m.CreateCustomProviderRow("local-gw-bare", ProtocolAnthropic, "http://127.0.0.1:9090", ""); err != nil {
		t.Fatalf("the provider-only control was refused: %v", err)
	}
	if bare, ok := readLoaded(t, path).LLM.Providers["local-gw-bare"]; !ok || len(bare.Models) != 0 {
		t.Errorf("the provider-only control landed as %+v", bare)
	}

	// The refusals, each naming its reason, each leaving the file byte-identical.
	before := readSettingsFile(t, path)
	cases := []struct {
		name     string
		provider string
		protocol string
		baseURL  string
		mustSay  string
	}{
		{"a preset name belongs to provider_row_create", "stepfun", ProtocolOpenAIChat, "https://api.stepfun.com/v1", "provider_row_create"},
		{"an empty base_url", "local-gw-two", ProtocolOpenAIChat, "   ", "base_url"},
		{"a protocol outside the enum", "local-gw-three", "grpc", "http://127.0.0.1:1", "protocol"},
		// The row this test created above: a non-preset name that already exists,
		// so this is refusal 4 and not the preset-name refusal one row up.
		{"a row the config already declares", "local-gw", ProtocolOpenAIChat, "http://127.0.0.1:8080/v1", "already exists"},
		{"no name at all", "   ", ProtocolOpenAIChat, "http://127.0.0.1:1", "行不存在"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			written, err := m.CreateCustomProviderRow(tc.provider, tc.protocol, tc.baseURL, "")
			if err == nil {
				t.Fatalf("the create was accepted (reported %v)", written)
			}
			if written != nil {
				t.Errorf("a refused create reported %v, want none", written)
			}
			if !strings.Contains(err.Error(), tc.mustSay) {
				t.Errorf("the refusal does not say why (%q missing): %v", tc.mustSay, err)
			}
			if after := readSettingsFile(t, path); after != before {
				t.Errorf("config.toml changed under a refused create:\n%s", after)
			}
		})
	}

	// Control: after five refusals the same call on a fresh name still lands, and
	// the third protocol word is a legal row too.
	if _, err := m.CreateCustomProviderRow("local-gw-four", ProtocolOpenAIResponses, "http://127.0.0.1:4/v1", ""); err != nil {
		t.Fatalf("the control write after five refusals: %v", err)
	}
	if p, ok := readLoaded(t, path).LLM.Providers["local-gw-four"]; !ok || p.Protocol != ProtocolOpenAIResponses {
		t.Errorf("the control row landed as %+v", p)
	}
}

// AC#3: the catalog entry. One guarded write adds the entry seeded from the
// schema's own default tags, the two refusals (an undeclared provider, an entry
// that already exists) leave the file byte-identical, and the leaves the page
// writes afterwards - context_window and the two price halves - address the entry
// this step created.
func Test320CreateModelRowSeedsFromDefaultsAndRefusesTwice(t *testing.T) {
	m, path, _ := newTestManager(t, settingsBase)

	written, err := m.CreateModelRow("deepseek", "deepseek-new-under-test-320")
	if err != nil {
		t.Fatalf("CreateModelRow: %v", err)
	}
	if want := []string{"llm.providers.deepseek.models.deepseek-new-under-test-320"}; !reflect.DeepEqual(written, want) {
		t.Errorf("written key paths = %v, want %v", written, want)
	}
	got := readLoaded(t, path)
	spec, ok := got.LLM.Providers["deepseek"].Models["deepseek-new-under-test-320"]
	if !ok {
		t.Fatalf("the catalog entry never landed: %+v", got.LLM.Providers["deepseek"].Models)
	}
	if !spec.Enabled {
		t.Error("the created entry is not enabled: an entry without enabled=true is treated as off and is refused at start-up")
	}
	if spec.Billing != BillingPayPerToken {
		t.Errorf("the entry's billing = %q, want the schema's own default", spec.Billing)
	}
	if body := readSettingsFile(t, path); !strings.Contains(body, "[llm.providers.deepseek.models.deepseek-new-under-test-320]") {
		t.Errorf("the entry's section is not in the file:\n%s", body)
	}
	if got.Ball.Size != 56 {
		t.Errorf("an unrelated key moved: ball.size = %d", got.Ball.Size)
	}

	// The step exists so the three leaf fields have something to address: the
	// control proves the entry the create made is the one they then write.
	if _, err := m.SetModelContextWindow("deepseek", "deepseek-new-under-test-320", 128000); err != nil {
		t.Fatalf("the leaf write against the created entry: %v", err)
	}
	if _, err := m.SetModelPriceIn("deepseek", "deepseek-new-under-test-320", 1500000); err != nil {
		t.Fatalf("the price write against the created entry: %v", err)
	}
	if after := readLoaded(t, path).LLM.Providers["deepseek"].Models["deepseek-new-under-test-320"]; after.ContextWindow != 128000 || after.Price.In != 1500000 {
		t.Errorf("the leaves did not land on the created entry: %+v", after)
	}

	// Two refusals, the file byte-identical under each.
	before := readSettingsFile(t, path)
	cases := []struct {
		name     string
		provider string
		model    string
		mustSay  string
	}{
		{"an undeclared provider", "ghostco", "m-320", "行不存在"},
		{"an entry that already exists", "deepseek", "deepseek-chat", "already exists"},
		{"no model id at all", "deepseek", "   ", "行不存在"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			written, err := m.CreateModelRow(tc.provider, tc.model)
			if err == nil {
				t.Fatalf("the create was accepted (reported %v)", written)
			}
			if written != nil {
				t.Errorf("a refused create reported %v, want none", written)
			}
			if !strings.Contains(err.Error(), tc.mustSay) {
				t.Errorf("the refusal does not say why (%q missing): %v", tc.mustSay, err)
			}
			if after := readSettingsFile(t, path); after != before {
				t.Errorf("config.toml changed under a refused create:\n%s", after)
			}
		})
	}

	// Control: after the refusals a second entry still lands, and the first one
	// (with its leaves) is untouched by that write.
	if _, err := m.CreateModelRow("deepseek", "deepseek-second-under-test-320"); err != nil {
		t.Fatalf("the control create after three refusals: %v", err)
	}
	after := readLoaded(t, path).LLM.Providers["deepseek"].Models
	if _, ok := after["deepseek-second-under-test-320"]; !ok {
		t.Fatalf("the control entry never landed: %+v", after)
	}
	if after["deepseek-new-under-test-320"].ContextWindow != 128000 {
		t.Errorf("the control write disturbed the first entry: %+v", after["deepseek-new-under-test-320"])
	}
}
