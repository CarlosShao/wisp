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
