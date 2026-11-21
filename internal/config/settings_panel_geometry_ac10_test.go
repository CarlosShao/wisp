package config

// Ticket 33 AC#10's config-side hot half: the two new [panel] geometry keys
// land through the same guarded writeOneKey path every settings setter uses,
// the live Manager answers the written values in the SAME process, and the
// file a fresh load reads back is the one the per-create geometry source
// (cmd/wisp's panelGeometryTripleSource) consumes. The refusal sentences name
// their key path and move no byte - the same shape SetPanelWidth's rulers
// pin, with the two keys' own legal zero kept legal (0 = auto height / 0 =
// follow system DPI, both the schema's own spellings).
//
// The reverse criterion each refusal names: deleting the setter's guard makes
// the refusal case red; deleting the write makes the hot-read case red.

import (
	"math"
	"path/filepath"
	"strings"
	"testing"
)

func TestAC10PanelHeightScaleSettersAreHotInTheSameProcess(t *testing.T) {
	dir := t.TempDir()
	seed := filepath.Join(dir, "config.toml")
	if err := SaveFile(seed, NewDefaults()); err != nil {
		t.Fatalf("SaveFile the defaults seed: %v", err)
	}
	mgr, err := NewManager(seed, nil)
	if err != nil {
		t.Fatalf("manager on the seed: %v", err)
	}

	// height: lands, the receipt names the key path the FILE changed by, and
	// the live Manager follows in this same process (the guarded writer keeps
	// file and memory equal).
	written, err := mgr.SetPanelHeight(450)
	if err != nil || len(written) != 1 || written[0] != "panel.height" {
		t.Fatalf("SetPanelHeight(450) = %v, %v, want [panel.height], nil", written, err)
	}
	if v := mgr.Config().Panel.Height; v != 450 {
		t.Fatalf("Config().Panel.Height = %d, want 450 - the write must be readable hot in this same process", v)
	}

	// scale: same shapes as height.
	written, err = mgr.SetPanelScale(1.25)
	if err != nil || len(written) != 1 || written[0] != "panel.scale" {
		t.Fatalf("SetPanelScale(1.25) = %v, %v, want [panel.scale], nil", written, err)
	}
	if v := mgr.Config().Panel.Scale; v != 1.25 {
		t.Fatalf("Config().Panel.Scale = %v, want 1.25", v)
	}

	// The file half: a fresh load of the SAME path answers what was written -
	// this is the file the per-create geometry source reads, so the hot chain
	// (write -> next create/re-show resolution) crosses here.
	c, _, err := LoadFile(seed, nil)
	if err != nil {
		t.Fatalf("LoadFile after the writes: %v", err)
	}
	if c.Panel.Height != 450 || c.Panel.Scale != 1.25 {
		t.Fatalf("a fresh load reads height=%d scale=%v, want 450 and 1.25: the guarded write left file and memory unequal", c.Panel.Height, c.Panel.Scale)
	}

	// The keys' legal zeros: 0 = auto height (the host answers from its own
	// constant), 0 = follow system DPI. Writing them is the legitimate "back
	// to the default", never garbage.
	if written, err := mgr.SetPanelHeight(0); err != nil || len(written) != 1 || written[0] != "panel.height" {
		t.Fatalf("SetPanelHeight(0) = %v, %v, want [panel.height], nil - 0 is the schema's own auto spelling", written, err)
	}
	if v := mgr.Config().Panel.Height; v != 0 {
		t.Fatalf("Config().Panel.Height = %d, want 0", v)
	}
	if written, err := mgr.SetPanelScale(0); err != nil || len(written) != 1 || written[0] != "panel.scale" {
		t.Fatalf("SetPanelScale(0) = %v, %v, want [panel.scale], nil - 0 is the schema's own follow-DPI spelling", written, err)
	}
	if v := mgr.Config().Panel.Scale; v != 0 {
		t.Fatalf("Config().Panel.Scale = %v, want 0", v)
	}

	// Refusals: a negative (and the two float values no consumer can act on)
	// is refused by name, and a refusal moves no byte.
	if _, err := mgr.SetPanelHeight(-1); err == nil || !strings.Contains(err.Error(), "panel.height") {
		t.Fatalf("a negative height must be refused by name, got: %v", err)
	}
	if _, err := mgr.SetPanelScale(-0.5); err == nil || !strings.Contains(err.Error(), "panel.scale") {
		t.Fatalf("a negative scale must be refused by name, got: %v", err)
	}
	if _, err := mgr.SetPanelScale(math.NaN()); err == nil || !strings.Contains(err.Error(), "panel.scale") {
		t.Fatalf("a NaN scale must be refused by name, got: %v", err)
	}
	if _, err := mgr.SetPanelScale(math.Inf(1)); err == nil || !strings.Contains(err.Error(), "panel.scale") {
		t.Fatalf("an infinite scale must be refused by name, got: %v", err)
	}
	if v := mgr.Config().Panel.Height; v != 0 {
		t.Fatalf("after the refused writes Height = %d, want the untouched 0 - a refusal moves no byte", v)
	}
	if v := mgr.Config().Panel.Scale; v != 0 {
		t.Fatalf("after the refused writes Scale = %v, want the untouched 0 - a refusal moves no byte", v)
	}
	if _, _, err := LoadFile(seed, nil); err != nil {
		t.Fatalf("the file no longer loads after the refused writes: %v", err)
	}

	// The refusals must also have left the FILE untouched: reload and compare
	// the two keys against the last legal write.
	c, _, err = LoadFile(seed, nil)
	if err != nil {
		t.Fatalf("re-load after the refusals: %v", err)
	}
	if c.Panel.Height != 0 || c.Panel.Scale != 0 {
		t.Fatalf("the refused writes reached the file anyway: height=%d scale=%v, want the last legal 0/0", c.Panel.Height, c.Panel.Scale)
	}
}
