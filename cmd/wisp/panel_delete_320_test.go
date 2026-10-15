package main

// The taking-back steps the owner asked for the moment the three writes landed
// ("added models cannot be deleted at all"): a catalog entry can be removed, a
// row can be removed, and a row that something still points at is refused BY
// NAME instead of quietly breaking the file's references.
//
// The shape is the 320/313 fixture's: one config.Manager over the seeded
// config.toml, one real ComposerDispatch, and "it was deleted" means the bytes
// in config.toml no longer carry the name - not a fake agreeing with itself.

import (
	"context"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/panel"
)

func Test320ModelRowDeleteThroughTheRealLeg(t *testing.T) {
	dir := settingsDir248(t)
	var audit []string
	disp, _ := newRealSettingsLeg(t, dir, &audit, nil)
	ctx := context.Background()

	apply313(t, disp, set313Env("r-320-del-add", panel.FieldModelRowCreate, "deepseek", "model-to-drop", ""),
		"llm.providers.deepseek.models.model-to-drop")
	if body := readConfigBody320(t, dir); !strings.Contains(body, "model-to-drop") {
		t.Fatalf("the entry the delete case removes was never created:\n%s", body)
	}

	apply313(t, disp, set313Env("r-320-del", panel.FieldModelRowDelete, "deepseek", "model-to-drop", ""),
		"llm.providers.deepseek.models.model-to-drop")
	if body := readConfigBody320(t, dir); strings.Contains(body, "model-to-drop") {
		t.Errorf("the deleted entry is still in config.toml:\n%s", body)
	}

	// The second click must be refused by name: "already gone" is a sentence, not
	// a silent success.
	if _, err := disp.Handle(ctx, set313Env("r-320-del-again", panel.FieldModelRowDelete, "deepseek", "model-to-drop", "")); err == nil {
		t.Error("deleting an entry that is not there cleared the leg")
	} else if !strings.Contains(err.Error(), "不在这个目录里") {
		t.Errorf("the refusal does not say the entry is missing: %v", err)
	}
}

func Test320ProviderRowDeleteRefusesWhileReferenced(t *testing.T) {
	dir := settingsDir248(t)
	var audit []string
	disp, _ := newRealSettingsLeg(t, dir, &audit, nil)
	ctx := context.Background()

	// The seed's chat role points at deepseek, so that row must not be removable -
	// and the refusal has to name the sentence that still points at it.
	before := readConfigBody320(t, dir)
	if _, err := disp.Handle(ctx, set313Env("r-320-del-row", panel.FieldProviderRowDelete, "deepseek", "", "")); err == nil {
		t.Fatal("the leg removed a provider row the chat role still points at")
	} else if !strings.Contains(err.Error(), "默认聊天模型") {
		t.Errorf("the refusal does not name the reference that blocks the delete: %v", err)
	}
	if after := readConfigBody320(t, dir); after != before {
		t.Errorf("config.toml changed under a refused row delete:\n%s", after)
	}

	// A row nothing points at goes away whole.
	apply313(t, disp, set313Env("r-320-del-row-add", panel.FieldProviderRowCreateCustom, "local-gw", "",
		"openai-chat,http://127.0.0.1:8080/v1"), "llm.providers.local-gw")
	apply313(t, disp, set313Env("r-320-del-row-go", panel.FieldProviderRowDelete, "local-gw", "", ""),
		"llm.providers.local-gw")
	if body := readConfigBody320(t, dir); strings.Contains(body, "local-gw") {
		t.Errorf("the deleted row is still in config.toml:\n%s", body)
	}
}
