//go:build windows

package main

// Ticket 65's teeth for the resident leg's snapshot read face
// (cmd/wisp/panel_resident_pump.go): the packet the shipping -H=windowsgui
// wisp.exe pushes must carry the reader-gated keys this leg REALLY has sources
// for, and must NOT carry the ones it does not. The assembly under test is the
// boot's own chain (newResidentPanelManager - the same call resident_windows.go
// makes), never a hand-built lookalike, and the wire shape is read off the
// bytes the 329 hop would Eval into the page - bit-for-bit what the page would
// have received.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
)

// snapshotHostRecorder65 captures the JS the pump's Out hands to the host, so
// the assertions read the exact bytes a page would have received.
type snapshotHostRecorder65 struct {
	mu sync.Mutex
	js []string
}

func (r *snapshotHostRecorder65) evalPanelJS(js string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.js = append(r.js, js)
}

func (r *snapshotHostRecorder65) calls() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.js...)
}

// decodedSnapshot65 unwraps one `window.wispOnPanelSnapshot(<json>)` call into
// the packet's top-level keys, preserving presence/absence exactly as the wire
// carried it (a map of RawMessage, not a struct with zero values).
func decodedSnapshot65(t *testing.T, js string) map[string]json.RawMessage {
	t.Helper()
	const global = "window.wispOnPanelSnapshot("
	if !strings.HasPrefix(js, global) || !strings.HasSuffix(js, ")") {
		t.Fatalf("the pushed JS is not the named-global call the 329 contract spells: %q", js[:min(80, len(js))])
	}
	raw := js[len(global) : len(js)-1]
	var top map[string]json.RawMessage
	if err := json.Unmarshal([]byte(raw), &top); err != nil {
		t.Fatalf("the pushed bytes are not the pump's JSON: %v", err)
	}
	return top
}

// Test65ResidentPumpCarriesConfigTruth is tooth ① + ②: the resident shape's
// packet carries ballStyle with the config's CURRENT style, plus modelSwitch/
// commands/plugins - each read off the chain Manager the settings door writes
// through, never off a test constant. The seed writes jarvis through the REAL
// setter BEFORE the assembly runs, so the chain reads a file that names jarvis
// (the same order residentSeed323 uses for its base URL: seeding the file
// after assembly would test a config the assembly never read).
func Test65ResidentPumpCarriesConfigTruth(t *testing.T) {
	dir := settingsDir248(t)
	seed, err := config.NewManager(filepath.Join(dir, configFileName), nil)
	if err != nil {
		t.Fatalf("config.NewManager over the seed: %v", err)
	}
	if _, err := seed.SetBallStyle("jarvis"); err != nil {
		t.Fatalf("SetBallStyle(jarvis): %v", err)
	}

	// The boot's own assembly: five returns now - the host, the chain Manager
	// the pump reads, the workspace holder, the ingest ledger, the error.
	host, cfgMgr, _, drafts, err := newResidentPanelManager(dir, nil)
	if err != nil {
		t.Fatalf("newResidentPanelManager refused to assemble: %v", err)
	}

	rec := &snapshotHostRecorder65{}
	pump := newResidentPanelPump(nil, cfgMgr, host.disp, drafts, rec)
	pump.publish()
	calls := rec.calls()
	if len(calls) != 1 {
		t.Fatalf("one publish must push exactly one packet to the host; got %d", len(calls))
	}
	top := decodedSnapshot65(t, calls[0])

	// ① ballStyle: present, and the style is the config's current value.
	raw, ok := top["ballStyle"]
	if !ok {
		t.Fatalf("the resident packet carries no ballStyle key; top-level keys: %v", keysOf65(top))
	}
	var bs panel.BallStyleSection
	if err := json.Unmarshal(raw, &bs); err != nil {
		t.Fatalf("ballStyle does not decode as the carrier: %v", err)
	}
	if bs.Style != "jarvis" {
		t.Errorf("ballStyle.style = %q, want jarvis (what the chain Manager holds after the setter wrote it)", bs.Style)
	}
	if len(bs.Vocabulary) == 0 {
		t.Errorf("ballStyle.vocabulary is empty; the schema's own vocabulary must travel verbatim")
	}

	// ② modelSwitch: present, and the current pair is the config's chat role.
	raw, ok = top["modelSwitch"]
	if !ok {
		t.Fatalf("the resident packet carries no modelSwitch key; top-level keys: %v", keysOf65(top))
	}
	var ms panel.ModelSwitchSection
	if err := json.Unmarshal(raw, &ms); err != nil {
		t.Fatalf("modelSwitch does not decode as the carrier: %v", err)
	}
	want := cfgMgr.ChatSwitchOptions()
	if ms.Current.Provider != want.CurrentProvider || ms.Current.Model != want.CurrentModel {
		t.Errorf("modelSwitch.current = %s/%s, want the config's %s/%s",
			ms.Current.Provider, ms.Current.Model, want.CurrentProvider, want.CurrentModel)
	}

	// ② commands: present (the production roster, derived from internal/command).
	if _, ok := top["commands"]; !ok {
		t.Fatalf("the resident packet carries no commands key; top-level keys: %v", keysOf65(top))
	}

	// ② plugins: present, and the inventory names the honest engine state.
	raw, ok = top["plugins"]
	if !ok {
		t.Fatalf("the resident packet carries no plugins key; top-level keys: %v", keysOf65(top))
	}
	var pl panel.PluginInventorySection
	if err := json.Unmarshal(raw, &pl); err != nil {
		t.Fatalf("plugins does not decode as the carrier: %v", err)
	}
	if pl.EngineState != panel.PluginEngineNotWired215 {
		t.Errorf("plugins.engineState = %q, want the carrier's honest not-wired sentence constant", pl.EngineState)
	}

	// ④ the run-shaped dimensions: ABSENT, on purpose. This leg runs no task
	// loop, so there is no RunFeed (run), no task roster (tasks) and no
	// instruction loader (instructions) behind it - a reader invented for any
	// of those would be a constant dressed as a read. Absence is the honest
	// wire shape pump.go pins for "this host has no source for the dimension";
	// the console shape's pipeline pump is the one that carries them.
	// (results is deliberately NOT on this list: it is not reader-gated - the
	// pump always sends the key, nil rows rendered empty - so its presence here
	// says "nothing has streamed", which is the truth in this shape.)
	for _, absent := range []string{"run", "tasks", "instructions"} {
		if _, ok := top[absent]; ok {
			t.Errorf("the resident packet carries a %q key, but this leg has no source for it - a faked reader would be a constant dressed as a read", absent)
		}
	}
}

// Test65ResidentPumpMutationStrippedBallStyleIsAbsent is tooth ③: take the REAL
// sources the boot wiring builds, strip exactly one reader, and the key must
// vanish from the wire - the reader-gated shape the pump pins, biting this
// file's own injection. A green here with the reader stripped proves the key's
// presence in the unmutated packet above came from the reader, not from
// plumbing that always writes the key.
func Test65ResidentPumpMutationStrippedBallStyleIsAbsent(t *testing.T) {
	dir := settingsDir248(t)
	host, cfgMgr, _, drafts, err := newResidentPanelManager(dir, nil)
	if err != nil {
		t.Fatalf("newResidentPanelManager refused to assemble: %v", err)
	}

	rec := &snapshotHostRecorder65{}
	src := residentPumpSources65(nil, cfgMgr, host.disp, drafts, rec)
	if src.BallStyle == nil {
		t.Fatalf("the unmutated sources carry no BallStyle reader; there is nothing to strip and the mutation would be vacuous")
	}
	src.BallStyle = nil

	mutated := panel.NewSnapshotPump(src)
	snap, _, err := mutated.Publish()
	if err != nil {
		t.Fatalf("the mutated pump refused to publish: %v", err)
	}
	if snap.BallStyle != nil {
		t.Fatalf("the snapshot still carries a BallStyle section with the reader stripped - the key is not reader-gated")
	}
	for _, js := range rec.calls() {
		top := decodedSnapshot65(t, js)
		if _, ok := top["ballStyle"]; ok {
			t.Fatalf("the pushed bytes carry a ballStyle key with the reader stripped: the wire is not reader-gated")
		}
	}
}

// Test65ResidentPumpNilManagerOmitsConfigKeys is tooth ④'s constructor face: a
// boot whose chain never assembled (cfg nil) leaves the config-derived readers
// off the pump instead of backing them with zero values - a packet that says
// "no reading" instead of one that invents an empty style or an empty menu.
func Test65ResidentPumpNilManagerOmitsConfigKeys(t *testing.T) {
	rec := &snapshotHostRecorder65{}
	pump := newResidentPanelPump(nil, nil, nil, newAttachmentDraftLedger(), rec)
	pump.publish()
	calls := rec.calls()
	if len(calls) != 1 {
		t.Fatalf("one publish must push exactly one packet; got %d", len(calls))
	}
	top := decodedSnapshot65(t, calls[0])
	for _, absent := range []string{"ballStyle", "modelSwitch", "plugins"} {
		if _, ok := top[absent]; ok {
			t.Errorf("with no config Manager the packet carries a %q key; a zero-value section would read as a made-up answer", absent)
		}
	}
}

// Test65ResidentShowTriggersResyncPins pins the resync-on-show wiring at its
// two seams, in the source-pin style this package's other assemblies use:
//
//	showOnThread - after a successful Show, the publisher must fire, or the
//	               frozen contract (PLAN.md:2982, SPEC-08 §5.1 "每次 show Go 侧
//	               推全量状态") stays a sentence with no mechanism;
//	resident_windows.go - the boot must select exactly one publisher, and the
//	                       no-pipeline branch must be the resident pump.
func Test65ResidentShowTriggersResyncPins(t *testing.T) {
	thread := readSource65(t, "panel_resident_windows.go")
	if !strings.Contains(thread, "rp.publishOnShow()") {
		t.Fatalf("showOnThread no longer fires the resync-on-show publisher after a successful Show; the frozen resync contract has no mechanism")
	}
	if !strings.Contains(thread, "func (rp *residentPanel) publishOnShow()") {
		t.Fatalf("residentPanel.publishOnShow is gone; the show trigger reads a field nobody owns")
	}

	boot := readSource65(t, "resident_windows.go")
	if !strings.Contains(boot, "holder.setSnapshotPublish(src.run.publishPanelSnapshot)") {
		t.Fatalf("the pipeline branch of the publisher selection is gone; the console shape's show would push nothing or the wrong packet")
	}
	if !strings.Contains(boot, "newResidentPanelPump(ra.gate, cfgMgr65, rp.disp, attachmentDrafts, rp).publish") {
		t.Fatalf("the no-pipeline branch of the publisher selection is gone; the shipping -H=windowsgui shape's show would push nothing again")
	}
	if !strings.Contains(boot, "h.stampPublishLocked()") {
		t.Fatalf("panelHolder no longer stamps freshly spawned threads with the publisher; a 319-style recreate would silently lose the resync")
	}
}

// readSource65 reads one production source file of this package for the pins.
func readSource65(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return string(b)
}

func keysOf65(top map[string]json.RawMessage) []string {
	out := make([]string, 0, len(top))
	for k := range top {
		out = append(out, k)
	}
	return out
}
