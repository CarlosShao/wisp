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
