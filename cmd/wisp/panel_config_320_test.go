package main

// Ticket 320 at the assembly root: the settings writes this ticket adds reach
// their real setters through the production dispatch, and the rows they create
// are what the read side the settings page already calls then answers.
//
// The shape is panel_config_313_test.go's: one config.Manager, one secret.Store,
// the one panel.ComposerDispatch - so "the write landed" means bytes in
// config.toml read back through the same leg, not a fake agreeing with itself.
// The two helpers this file borrows (settingsDir248, newRealSettingsLeg, apply313,
// set313Env) are the ticket-248/313 fixtures, used unchanged.

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
)

// readConfigBody320 reads config.toml off a settings dir. The refusals below
// compare the whole body, because "nothing was written" is the claim.
func readConfigBody320(t *testing.T, dir string) string {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(dir, configFileName))
	if err != nil {
		t.Fatalf("read config.toml: %v", err)
	}
	return string(body)
}

// AC#1 through the real leg: the billing word lands, the receipt names the key
// path and keeps the restart wording, and both refusal directions leave the file
// byte-identical (a word outside the enum, a provider the config does not
// declare).
func Test320ProviderBillingThroughTheRealLeg(t *testing.T) {
	dir := settingsDir248(t)
	var audit []string
	disp, _ := newRealSettingsLeg(t, dir, &audit, nil)
	ctx := context.Background()

	apply313(t, disp, set313Env("r-320-billing", panel.FieldProviderBilling, "deepseek", "", config.BillingPlan),
		"llm.providers.deepseek.billing")
	if body := readConfigBody320(t, dir); !strings.Contains(body, "billing = 'plan'") {
		t.Errorf("config.toml does not carry the billing mode:\n%s", body)
	}

	before := readConfigBody320(t, dir)
	if _, err := disp.Handle(ctx, set313Env("r-320-billing-bad", panel.FieldProviderBilling, "deepseek", "", "flat-rate")); err == nil {
		t.Fatal("a billing word outside the enum cleared the leg")
	} else if !strings.Contains(err.Error(), "billing") {
		t.Errorf("the refusal does not name the field: %v", err)
	}
	if after := readConfigBody320(t, dir); after != before {
		t.Errorf("config.toml changed under a refused billing value:\n%s", after)
	}

	if _, err := disp.Handle(ctx, set313Env("r-320-billing-ghost", panel.FieldProviderBilling, "ghostco", "", config.BillingPlan)); err == nil {
		t.Fatal("a billing write for a provider the config does not declare cleared the leg")
	} else if !strings.Contains(err.Error(), "行不存在") {
		t.Errorf("the undeclared-provider refusal is not reason 2 (行不存在): %v", err)
	}
	if after := readConfigBody320(t, dir); after != before {
		t.Errorf("config.toml changed under a refused ghost-provider write:\n%s", after)
	}

	// Control: after both refusals the same field still writes through the real
	// leg (the refusals did not wedge the door or the writer).
	apply313(t, disp, set313Env("r-320-billing-ctl", panel.FieldProviderBilling, "deepseek", "", config.BillingPayPerToken),
		"llm.providers.deepseek.billing")
	if body := readConfigBody320(t, dir); !strings.Contains(body, "billing = 'pay-per-token'") {
		t.Errorf("the control write did not land:\n%s", body)
	}
}

// customModel320 is the made-up catalog id AC#2's row carries as its optional
// first model. The name is this test's own; no real model is invented (A877
// ruling 4's rule, inherited).
const customModel320 = "custom-model-under-test-320"

// providerRow320 finds one provider row in the config.get payload.
func providerRow320(t *testing.T, payload map[string]any, name string) map[string]any {
	t.Helper()
	rows, _ := payload["providers"].([]any)
	for _, r := range rows {
		row, _ := r.(map[string]any)
		if row["name"] == name {
			return row
		}
	}
	t.Fatalf("the payload carries no provider row %q: %v", name, payload["providers"])
	return nil
}

// AC#2 through the real leg: a non-preset row appears with the protocol and
// base_url the value carried (the space after the comma is trimmed by the leg's
// split, the wire format's own courtesy), the first model lands with it, the read
// side the page already calls sees both, and a preset name is refused with the
// file untouched.
func Test320CustomProviderRowThroughTheRealLeg(t *testing.T) {
	dir := settingsDir248(t)
	var audit []string
	disp, _ := newRealSettingsLeg(t, dir, &audit, nil)
	ctx := context.Background()

	apply313(t, disp, set313Env("r-320-custom-row", panel.FieldProviderRowCreateCustom, "local-gw", customModel320,
		"openai-chat, http://127.0.0.1:8080/v1"), "llm.providers.local-gw")
	body := readConfigBody320(t, dir)
	for _, frag := range []string{
		"[llm.providers.local-gw]",
		"protocol = 'openai-chat'",
		"base_url = 'http://127.0.0.1:8080/v1'",
		"[llm.providers.local-gw.models." + customModel320 + "]",
	} {
		if !strings.Contains(body, frag) {
			t.Errorf("config.toml does not carry %q:\n%s", frag, body)
		}
	}

	// Read back: the row and its catalog entry reach the page's own data source.
	reply, err := disp.Handle(ctx, `{"method":"config.get","requestId":"r-320-custom-get","source":"panel-composer"}`)
	if err != nil {
		t.Fatalf("config.get after the custom row: %v", err)
	}
	row := providerRow320(t, settingsPayload313(t, reply), "local-gw")
	models, _ := row["models"].([]any)
	if row["modelCount"] != float64(1) || len(models) != 1 || models[0] != customModel320 {
		t.Errorf("the custom row read back as %v", row)
	}
	if row["hasBaseURL"] != true {
		t.Errorf("the read side does not see the base_url: %v", row)
	}

	// Refusal direction: a built-in preset name belongs to provider_row_create,
	// and the whole value must be present - both leave the file byte-identical.
	before := readConfigBody320(t, dir)
	if _, err := disp.Handle(ctx, set313Env("r-320-custom-preset", panel.FieldProviderRowCreateCustom, "stepfun", "",
		"openai-chat,https://api.stepfun.com/v1")); err == nil {
		t.Fatal("a preset name cleared the custom-row step")
	} else if !strings.Contains(err.Error(), "provider_row_create") {
		t.Errorf("the preset-name refusal does not name the step that owns preset rows: %v", err)
	}
	if _, err := disp.Handle(ctx, set313Env("r-320-custom-novalue", panel.FieldProviderRowCreateCustom, "local-gw-two", "",
		"openai-chat")); err == nil {
		t.Fatal("a value with no second half cleared the custom-row step")
	} else if !strings.Contains(err.Error(), "base_url") {
		t.Errorf("the half-value refusal does not name the missing half: %v", err)
	}
	if after := readConfigBody320(t, dir); after != before {
		t.Errorf("config.toml changed under the two refused custom-row writes:\n%s", after)
	}

	// Control: a second custom row still lands after both refusals.
	apply313(t, disp, set313Env("r-320-custom-ctl", panel.FieldProviderRowCreateCustom, "local-gw-four", "",
		"openai-responses,http://127.0.0.1:4/v1"), "llm.providers.local-gw-four")
	if body := readConfigBody320(t, dir); !strings.Contains(body, "[llm.providers.local-gw-four]") {
		t.Errorf("the control custom row did not land:\n%s", body)
	}
}
