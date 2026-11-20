//go:build windows

package main

// Ticket 193 AC#3's wiring tests, cmd/wisp half: the pump source lines, the
// ApplySetting write branch and the roster union all land in ONE commit, and
// the chain "设置页写 panel_animation=off ⇒ 下一次快照 panelAnimation.tier==off"
// is proven over the PRODUCTION assembly - the boot's own chain
// (newResidentPanelManager), the settings door it really answers with
// (host.disp -> configStore -> SetPanelAnimation), and the resident pump's
// real sources - never a config-layer-only fake agreeing with itself.
//
// THE TEETH (each names the mutation it reddens):
//   - 泵源假 off: the two reader lines ride the tier VERBATIM off
//     PanelAnimationOptions. A wiring that answers "off" when the read has no
//     value - or that strips the reader and lets a default replace absence -
//     reddens Test193PanelAnimationMutationStrippedReaderIsAbsent (key must
//     vanish) and Test65ResidentPumpNilManagerOmitsConfigKeys (nil Manager ⇒
//     no key, no invented tier), and the verbatim assertions below (reduce
//     reads reduce, off reads off, full reads full) redden any invented word.
//   - 名册并集: Test313WritableUnionEqualsTheTwoRosters grew the
//     WritableFields193() append in the same commit (A-313's lesson); a roster
//     added without that sync reddens it.
//   - 写路由错指: an ApplySetting case pointed at a field that is not
//     panel_animation reddens Test193PanelAnimationWriteRidesTheNextSnapshot
//     (the receipt names panel.animation and the disk carries it; a misrouted
//     case writes the wrong key or refuses), and the unlisted-name control
//     proves the door refuses what no roster carries.
//
// Same counting/identity rules as the sibling files: wire shapes read off the
// pushed bytes, refusals judged by error identity, never by log strings.

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/panel"
)

// Test193AssemblyRootsCarryTheAnimationReader pins the two pump-source lines
// at the two PumpSources assembly sites (the same disk reading
// panel_appearance_wave_test.go makes for the wave's four readers): deleting
// either reader line reddens its case, and the key would silently vanish from
// every packet that shape pushes.
func Test193AssemblyRootsCarryTheAnimationReader(t *testing.T) {
	runGo, err := os.ReadFile("run.go")
	if err != nil {
		t.Fatalf("read run.go: %v", err)
	}
	if !strings.Contains(string(runGo), "PanelAnimation: func() panel.PanelAnimationSection") {
		t.Errorf("run.go's PumpSources literal no longer carries the PanelAnimation reader: the panelAnimation key would silently vanish from every packet the pipeline shape pushes (ticket 193 AC#3)")
	}
	residentGo, err := os.ReadFile("panel_resident_pump.go")
	if err != nil {
		t.Fatalf("read panel_resident_pump.go: %v", err)
	}
	if !strings.Contains(string(residentGo), "src.PanelAnimation = func() panel.PanelAnimationSection") {
		t.Errorf("panel_resident_pump.go no longer carries the PanelAnimation reader: the shipping double-click shape would push packets without the key (ticket 193 AC#3)")
	}
}

// Test193PanelAnimationWriteRidesTheNextSnapshot is the end-to-end tooth:
// drive config.set panel_animation=off through the REAL settings door of the
// REAL resident chain, then take the NEXT snapshot the resident pump pushes
// and read the tier word off the wire. Every hop is production: the door is
// the chain's ConfigWriteHandler over the chain's own Manager, the write is
// SetPanelAnimation over the guarded path, the packet is residentPumpSources65's
// readers over the SAME Manager. No config-layer-only fake anywhere.
func Test193PanelAnimationWriteRidesTheNextSnapshot(t *testing.T) {
	dir := settingsDir248(t) // settingsSeed has no [panel] section: the schema default "reduce" is what the file means
	host, cfgMgr, _, drafts, err := newResidentPanelManager(dir, nil)
	if err != nil {
		t.Fatalf("newResidentPanelManager refused to assemble: %v", err)
	}

	// tierOnWire publishes one fresh packet (one recorder per read, so the
	// recorder's mutex is never contended across reads) and reads the tier
	// word off the exact bytes the page would have received.
	tierOnWire := func() string {
		t.Helper()
		rec := &snapshotHostRecorder65{}
		pump := newResidentPanelPump(nil, cfgMgr, host.disp, drafts, rec)
		pump.publish()
		calls := rec.calls()
		if len(calls) != 1 {
			t.Fatalf("one publish must push exactly one packet; got %d", len(calls))
		}
		top := decodedSnapshot65(t, calls[0])
		raw, ok := top["panelAnimation"]
		if !ok {
			t.Fatalf("the resident packet carries no panelAnimation key; the assembly-root reader line is gone or gated wrong; top-level keys: %v", keysOf65(top))
		}
		var sec panel.PanelAnimationSection
		if err := json.Unmarshal(raw, &sec); err != nil {
			t.Fatalf("panelAnimation does not decode as the carrier: %v", err)
		}
		return sec.Tier
	}

	// Before the write: the schema default rides verbatim (a wiring that
	// invents "off" for an unset row reddens HERE, before anything was written).
	if got := tierOnWire(); got != "reduce" {
		t.Fatalf("the pre-write snapshot's panelAnimation.tier = %q, want the schema default %q read verbatim off the config the chain holds", got, "reduce")
	}

	// The write, through the production settings door.
	reply, err := host.disp.Handle(context.Background(), set313Env("r-193-set-off", panel.FieldPanelAnimation, "", "", "off"))
	if err != nil {
		t.Fatalf("config.set panel_animation=off was refused: %v", err)
	}
	// The receipt carries the D36 tier sentence for the hot tier: 已立即生效
	// column, never a restart sentence (TierRegistry's whole-section "panel"
	// row = hot; the Note names the snapshot cadence honestly).
	if !strings.Contains(reply, "立即生效") {
		t.Errorf("the panel_animation receipt lost the hot-tier wording: %q", reply)
	}
	if strings.Contains(reply, "需要重启") {
		t.Errorf("the panel_animation receipt claims a restart for a hot-tier key: %q", reply)
	}
	if !strings.Contains(reply, "panel.animation") {
		t.Errorf("the receipt never names the key path it wrote: %q", reply)
	}
	// On disk, not only in memory.
	body, err := os.ReadFile(filepath.Join(dir, configFileName))
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if !strings.Contains(string(body), "animation = 'off'") {
		t.Errorf("config.toml does not carry the written animation tier:\n%s", body)
	}

	// THE ASSERTION THE TICKET ASKS FOR: the NEXT snapshot runs the new tier.
	if got := tierOnWire(); got != "off" {
		t.Fatalf("the next snapshot's panelAnimation.tier = %q, want %q - the write route, the pump reader and the carrier must agree (真链路)", got, "off")
	}

	// The other two vocabulary words ride verbatim too: a wiring that maps
	// every tier to "off" reddens on this leg.
	if _, err := host.disp.Handle(context.Background(), set313Env("r-193-set-full", panel.FieldPanelAnimation, "", "", "full")); err != nil {
		t.Fatalf("config.set panel_animation=full was refused: %v", err)
	}
	if got := tierOnWire(); got != "full" {
		t.Fatalf("after writing full the next snapshot's panelAnimation.tier = %q, want the word read verbatim", got)
	}

	// MUTATION TOOTH (写路由错指): a field name no roster carries is refused by
	// the door, reaches no leg, and moves nothing - the write route only ever
	// runs the names the rosters admit.
	before, rerr := os.ReadFile(filepath.Join(dir, configFileName))
	if rerr != nil {
		t.Fatalf("read config before the refused write: %v", rerr)
	}
	if _, err := host.disp.Handle(context.Background(), set313Env("r-193-bad-route", "panel_animation_typo", "", "", "off")); err == nil {
		t.Fatalf("a config.set naming a field no roster carries must be refused")
	} else if !strings.Contains(err.Error(), "不在可写名册内") {
		t.Errorf("the unlisted-name refusal is not the door's own roster sentence: %v", err)
	}
	after, rerr := os.ReadFile(filepath.Join(dir, configFileName))
	if rerr != nil {
		t.Fatalf("read config after the refused write: %v", rerr)
	}
	if string(before) != string(after) {
		t.Fatalf("a refused write moved config.toml anyway")
	}
	if got := cfgMgr.PanelAnimationOptions().Tier; got != "full" {
		t.Fatalf("the refused write moved the live tier: %q, want %q", got, "full")
	}
}

// Test193PanelAnimationMutationStrippedReaderIsAbsent is the reader-gated
// tooth at the wiring level (the 65 ballStyle mutation's twin): take the REAL
// sources the boot wiring builds, strip exactly the PanelAnimation reader, and
// the key must vanish from the wire - never replaced by a default word. A
// green here with the reader stripped proves the key's presence came from the
// reader, and a wiring that "answered off where there is no reading" instead
// of leaving the reader off reddens the nil-Manager pin
// (Test65ResidentPumpNilManagerOmitsConfigKeys) alongside it.
func Test193PanelAnimationMutationStrippedReaderIsAbsent(t *testing.T) {
	dir := settingsDir248(t)
	host, cfgMgr, _, drafts, err := newResidentPanelManager(dir, nil)
	if err != nil {
		t.Fatalf("newResidentPanelManager refused to assemble: %v", err)
	}
	rec := &snapshotHostRecorder65{}
	src := residentPumpSources65(nil, cfgMgr, host.disp, drafts, rec)
	if src.PanelAnimation == nil {
		t.Fatalf("the unmutated sources carry no PanelAnimation reader; there is nothing to strip and the mutation would be vacuous")
	}
	src.PanelAnimation = nil

	mutated := panel.NewSnapshotPump(src)
	snap, _, err := mutated.Publish()
	if err != nil {
		t.Fatalf("the mutated pump refused to publish: %v", err)
	}
	if snap.PanelAnimation != nil {
		t.Fatalf("the snapshot still carries a PanelAnimation section with the reader stripped - the key is not reader-gated")
	}
	for _, js := range rec.calls() {
		top := decodedSnapshot65(t, js)
		if _, ok := top["panelAnimation"]; ok {
			t.Fatalf("the pushed bytes carry a panelAnimation key with the reader stripped: the wire is not reader-gated")
		}
	}
}
