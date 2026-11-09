package main

// Ticket 187's assembly-root tooth: run.go's PumpSources literal hands the
// pump the modelSwitch reader, so the composer row's switch dimension
// (internal/panel/model_switch_187.go's carrier, read off the live config by
// internal/config's ChatSwitchOptions - the production reader the pump field
// doc names) actually reaches the wire. Without that one line the carrier and
// the write leg shipped packets carrying no "modelSwitch" key at all: the
// exact "carrier exists, nothing reads it" hop ticket 200 AC#7 closed for its
// own dimension and ticket 65 closed for ballStyle, and the shape the
// reader-gated tests in internal/panel cannot see from over there.
//
// Teeth:
//  1. the wiring exists: run.go's literal still carries the reader line
//     (the 65 byte-nail shape - a deleted line reds here immediately);
//  2. the wired reading is right: a real Manager loaded from the same
//     minimal 187 config the store tests use, read through the same chain
//     the literal uses, marshaled through the production pump, sends the
//     "modelSwitch" key whose fields carry ChatSwitchOptions' own values
//     (hardcoded here so the nail is independent of the reader's code);
//  3. the mutation: with the reader gone the wire bytes carry no
//     "modelSwitch" key at all - the reader-gated 213 shape, asserted on
//     the bytes the panel would receive, not a Go field.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
)

func Test187AssemblyRootCarriesTheModelSwitchReader(t *testing.T) {
	runGo, err := os.ReadFile("run.go")
	if err != nil {
		t.Fatalf("read run.go: %v", err)
	}
	if !strings.Contains(string(runGo), "ModelSwitch: func() panel.ModelSwitchSection") {
		t.Error("run.go's PumpSources literal no longer carries the modelSwitch reader: the composer row's switch dimension would silently vanish from every packet this process pushes (ticket 187)")
	}
}

// modelSwitchSection187 exercises the SAME mapping function run.go's
// PumpSources literal calls (modelSwitchSection, ticket 145's refactor of
// the five-field literal) - production and test share one body, so the
// mirror cannot drift. The byte-nail above still pins the literal's
// presence; this helper hands it a nil provider, the shape the helper's
// original five-field assertions were built on.
func modelSwitchSection187(m *config.Manager) panel.ModelSwitchSection {
	return modelSwitchSection(m, nil)
}

// Test187WiredModelSwitchReaderSendsTheLiveValues: the reader chain the
// assembly root uses, against a real Manager, through the production pump's
// marshal - the key arrives and every field carries the live config's values
// verbatim (the picker shows what the next ApplySetting would act on).
func Test187WiredModelSwitchReaderSendsTheLiveValues(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, configFileName)
	if err := os.WriteFile(path, []byte(thinking187Base), 0o600); err != nil {
		t.Fatal(err)
	}
	mgr, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("manager: %v", err)
	}

	data, err := panel.NewSnapshotPump(panel.PumpSources{
		ModelSwitch: func() panel.ModelSwitchSection { return modelSwitchSection187(mgr) },
	}).Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var envelope struct {
		ModelSwitch *panel.ModelSwitchSection `json:"modelSwitch"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		t.Fatalf("the wire bytes are not the snapshot's JSON: %v", err)
	}
	if envelope.ModelSwitch == nil {
		t.Fatalf("the wired reader sent no modelSwitch key at all:\n%s", data)
	}
	got := *envelope.ModelSwitch

	// Hardcoded from thinking187Base: one provider, one enabled model
	// declaring off|low, the chat role pinned to it with "off". Keeping the
	// expectations literal makes this nail independent of ChatSwitchOptions'
	// own code - a config-side regression reds here too.
	want := panel.ModelSwitchSection{
		Current:         panel.ModelOption{Provider: "deepseek", Model: "deepseek-chat"},
		CurrentThinking: "off",
		Models:          []panel.ModelOption{{Provider: "deepseek", Model: "deepseek-chat"}},
		ThinkingLevels:  []string{"off", "low"},
	}
	if got.Current != want.Current || got.CurrentThinking != want.CurrentThinking || got.Reason != "" {
		t.Errorf("modelSwitch current/tier/reason = %+v, %q, %q, want %+v, %q, empty", got.Current, got.CurrentThinking, got.Reason, want.Current, want.CurrentThinking)
	}
	if len(got.Models) != 1 || got.Models[0] != want.Models[0] {
		t.Errorf("modelSwitch models = %+v, want exactly [%+v]", got.Models, want.Models[0])
	}
	if len(got.ThinkingLevels) != 2 || got.ThinkingLevels[0] != "off" || got.ThinkingLevels[1] != "low" {
		t.Errorf("modelSwitch thinkingLevels = %v, want [off low] (the current model's own declared list)", got.ThinkingLevels)
	}
	// The picker must agree with what a write through the same Manager would
	// act on - one read face, not two that can drift.
	if opts := mgr.ChatSwitchOptions(); opts.CurrentProvider != got.Current.Provider || opts.CurrentModel != got.Current.Model || opts.CurrentThinking != got.CurrentThinking {
		t.Errorf("wire section %+v disagrees with ChatSwitchOptions %+v (the picker lies about what the next write acts on)", got, opts)
	}
}

// Test187MutationReaderRemovedSendsNoModelSwitchKeyAtAll: the assembly root's
// reader gone - the pump the four-key byte nails pin. The key must be absent
// from the WIRE BYTES, not merely nil in Go: composer omitempty is the last
// mile, and this asserts the bytes the panel receives.
func Test187MutationReaderRemovedSendsNoModelSwitchKeyAtAll(t *testing.T) {
	data, err := panel.NewSnapshotPump(panel.PumpSources{}).Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if bytes.Contains(data, []byte(`"modelSwitch"`)) {
		t.Errorf("a pump without the ModelSwitch reader still sent the key:\n%s", data)
	}
}
