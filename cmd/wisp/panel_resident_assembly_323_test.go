//go:build windows

package main

// Ticket 323 AC#6: the resident host's assembly line gets an end-to-end harness.
//
// WHAT 323 FOUND. The ticket's residual #6 (from the non-implementer acceptance
// leg 320-v1) reads: "常驻宿主那一行只有读码层 —— cmd/wisp/panel_resident_windows.go
// 的装配今天没有端到端台件". The readings that existed before this file were code
// reads plus two neighbouring pins that each hold ONE slice of the assembly:
// panel_mode_confirm_114_windows_test.go drives the chain through
// newResidentComposerDispatch but builds its host with NewPanelManager directly,
// so the production assembly function itself - the BuiltinAssets fallback, the
// geometry source wiring, the dataPath join - was never executed by any test.
// This file calls newResidentPanelManager, the exact call
// cmd/wisp/resident_windows.go:303 makes at boot, and holds the readings that
// assembly owes.
//
// WHAT IS DELIBERATELY NOT HERE. No WebView2 window is created: every reading
// below is taken from objects the assembly returns or methods that answer
// without a window (windowOptions consults the geometry source; dispatchRaw is
// the page's door and runs the chain on the caller's goroutine). The real-window
// hop stays where ticket 33's leg put it, and nothing here claims it.

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
)

// newResidentAssemblyDir323 seeds the data root the resident boot hands
// newResidentPanelManager: a档 the write leg can answer, plus a [panel] pair
// chosen to be NOT the host constants (panel_host_windows.go's 420x260) and NOT
// the schema defaults (640x0), so "these numbers came from the config on disk"
// has exactly one possible source - the assembly's geometry wiring.
func newResidentAssemblyDir323(t *testing.T) string {
	t.Helper()
	dir := sealableTempDir124(t)
	cfg := "[risk]\npermission_mode = \"ask_high_risk\"\n\n[panel]\nwidth = 812\nheight = 534\n"
	path := filepath.Join(dir, configFileName)
	if err := os.WriteFile(path, []byte(cfg), 0o600); err != nil {
		t.Fatalf("cannot seed %s: %v", path, err)
	}
	return dir
}

// modeRequest323 is one panel.mode.request envelope spelled the way the page
// spells it (the same shape panel_mode_confirm_114_windows_test.go sends; the
// copy stands on its own so this file's claim does not lean on 114's helpers
// for its request).
func modeRequest323(rid, to string) string {
	return `{"method":"panel.mode.request","requestId":"` + rid + `","source":"` +
		inboundPanelSource + `","to":"` + to + `"}`
}

// Test323ResidentHostAssemblyEndToEnd is AC#6's harness: the production
// assembly (newResidentPanelManager, not a rebuild of its parts) must return a
// live host plus the scope holder the task source binds later, the host's born
// geometry must be the config on disk (the one config read 票 255 AC#4 gave this
// leg), and the assembled door must carry a real page request through the chain
// onto the disk.
func Test323ResidentHostAssemblyEndToEnd(t *testing.T) {
	ctx := context.Background()
	dir := newResidentAssemblyDir323(t)

	// Pre-check that aims the geometry reading: the config layer itself keeps
	// the seeded [panel] pair. Without this, a red below could not tell
	// "the assembly is blind" from "the config layer dropped the section".
	seeded, _, err := config.LoadFile(filepath.Join(dir, configFileName), nil)
	if err != nil {
		t.Fatalf("the seeded config does not load: %v", err)
	}
	if seeded.Panel.Width != 812 || seeded.Panel.Height != 534 {
		t.Fatalf("the config layer did not keep the seeded [panel] pair (reads %dx%d); this test's geometry reading is not aimed",
			seeded.Panel.Width, seeded.Panel.Height)
	}

	// (1) The assembly itself, exactly the boot's call shape (the L2 leg goes
	// through panelModeConfirm, the same adapter the run seam uses).
	leg := &countingPermConfirm{}
	mgr, _, scope, _, err := newResidentPanelManager(dir, panelModeConfirm(leg.fn()))
	if err != nil {
		t.Fatalf("the resident assembly failed: %v", err)
	}
	if mgr == nil {
		t.Fatal("newResidentPanelManager returned a nil host with no error")
	}
	if scope == nil {
		t.Fatal("newResidentPanelManager returned a nil workspace scope holder; the task source would have nothing to bind (ticket 186's ordering)")
	}

	// (2) The host the assembly returns was born at the config's size, not the
	// constants: windowOptions consults the geometry source the assembly wired.
	opts := mgr.windowOptions()
	if opts.Width != 812 || opts.Height != 534 {
		t.Errorf("the assembled host resolves its window at %dx%d; the config on disk says [panel] 812x534, so the assembly's geometry source is not the one this leg wires (host constants are 420x260, schema defaults 640x0)",
			opts.Width, opts.Height)
	}

	// (3) The door the assembly attached answers a real page request end to
	// end: a narrowing (ask_high_risk -> ask_every_step) needs no card (R9: a
	// tightening raises none), reaches the write leg, and lands on the disk.
	receipt, err := mgr.dispatchRaw(ctx, modeRequest323("r323-assembly", "ask_every_step"))
	if err != nil {
		t.Fatalf("the assembled door refused a narrowing (R9 says a tightening raises no card): %v (receipt %q)", err, receipt)
	}
	if got := readInboundCfg33(t, dir); got != "ask_every_step" {
		t.Errorf("the accepted narrowing never landed: disk reads %q, want ask_every_step", got)
	}
}

// Test323ResidentGeometrySourceUnreadableConfigFallsToZero holds the loud
// branch of the assembly's one config read (panel_resident_windows.go's
// panelGeometrySource): a config the layer cannot parse is answered with the
// 0,0 that makes the host fall back to its own constants and say so in the log,
// never with a made-up size.
func Test323ResidentGeometrySourceUnreadableConfigFallsToZero(t *testing.T) {
	dir := t.TempDir()
	// A directory where the config file should be is the one unreadable shape
	// no parser can talk itself past.
	if err := os.Mkdir(filepath.Join(dir, configFileName), 0o700); err != nil {
		t.Fatalf("cannot plant the unreadable config shape: %v", err)
	}
	width, height := panelGeometrySource(dir)()
	if width != 0 || height != 0 {
		t.Errorf("an unreadable config answered the geometry source with %dx%d; the loud branch returns 0,0 so the host falls back to its own constants and names the failure in the log", width, height)
	}
}
