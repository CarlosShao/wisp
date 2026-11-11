package config

// The appearance wave's config-side rulers: the four setters land their keys
// through the same guarded writeOneKey path every settings setter uses, the
// read faces answer what the live Manager holds, and the refusals are the
// single-sourced sentences - opacity's band is validate.go's own
// ballOpacityIdleRangeError (the same string a bad file fails at boot with),
// and the two [panel] setters refuse the one shape that is meaningless in
// every consumer (non-positive px) instead of inventing a band nobody froze.
//
// The reverse criterion each case names: deleting its ApplySetting case
// (cmd/wisp/panel_config_store.go) reddens the store half of the wire test;
// deleting the setter body reddens here.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// appearanceWavePlant builds a minimal valid config.toml carrying one [ball]
// and one [panel] value each, non-default so a default leaking into any
// reading cannot fake the answer.
func appearanceWavePlant(t *testing.T, dir string) string {
	t.Helper()
	body := "[ball]\nopacity_idle = 0.7\nclick_through = false\n\n[panel]\nwidth = 700\nfont_size = 14\n"
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write the wave plant: %v", err)
	}
	return path
}

func TestAppearanceWaveSettersLandAndRefuse(t *testing.T) {
	dir := t.TempDir()
	seedPath := filepath.Join(dir, "config.toml")
	if err := SaveFile(seedPath, NewDefaults()); err != nil {
		t.Fatalf("SaveFile the defaults seed: %v", err)
	}
	mgr, err := NewManager(seedPath, nil)
	if err != nil {
		t.Fatalf("manager on the seed: %v", err)
	}

	// opacity_idle, in band: lands, the receipt names the key path the FILE
	// changed by, and the memory view follows (the guarded writer keeps file
	// and memory equal).
	written, err := mgr.SetBallOpacityIdle(0.9)
	if err != nil {
		t.Fatalf("SetBallOpacityIdle(0.9): %v", err)
	}
	if len(written) != 1 || written[0] != "ball.opacity_idle" {
		t.Fatalf("receipt = %v, want [ball.opacity_idle]", written)
	}
	if v := mgr.BallOpacityIdleOptions(); v.Opacity != 0.9 {
		t.Fatalf("BallOpacityIdleOptions().Opacity = %v, want 0.9 - the write must be readable hot in this same process", v.Opacity)
	}

	// opacity_idle, out of band: refused with the validate.go sentence
	// VERBATIM (the single-source shape ballSizeRangeError keeps), and the
	// file's value does not move.
	_, err = mgr.SetBallOpacityIdle(1.5)
	if err == nil {
		t.Fatalf("1.5 must be refused by the setter's band guard")
	}
	if got := err.Error(); !strings.Contains(got, "config.toml: ball.opacity_idle 1.5 out of range [0, 1]") {
		t.Fatalf("the 1.5 refusal is not the single-sourced validate sentence, got: %v", err)
	}
	if v := mgr.BallOpacityIdleOptions(); v.Opacity != 0.9 {
		t.Fatalf("after the refused write Opacity = %v, want the untouched 0.9 - a refusal moves no byte", v.Opacity)
	}

	// The load gate and the write gate are the SAME sentence by construction:
	// a hand-edited file outside the band fails at boot with the string
	// ballOpacityIdleRangeError formats, spelled here so the equality is
	// checkable rather than asserted in a comment.
	plant := appearanceWavePlant(t, t.TempDir())
	raw, err := os.ReadFile(plant)
	if err != nil {
		t.Fatalf("read the plant: %v", err)
	}
	bad := strings.Replace(string(raw), "opacity_idle = 0.7", "opacity_idle = 1.5", 1)
	if bad == string(raw) {
		t.Fatalf("the band plant changed nothing")
	}
	if err := os.WriteFile(plant, []byte(bad), 0o600); err != nil {
		t.Fatalf("rewrite the plant out of band: %v", err)
	}
	_, _, err = LoadFile(plant, nil)
	if err == nil {
		t.Fatalf("a 1.5 file must fail at load")
	}
	if !strings.Contains(err.Error(), "config.toml: ball.opacity_idle 1.5 out of range [0, 1]") {
		t.Fatalf("the load refusal is not the same sentence the write gate hands out: %v", err)
	}

	// click_through: a bool has no band; the write lands and the view follows.
	if written, err := mgr.SetBallClickThrough(false); err != nil || len(written) != 1 || written[0] != "ball.click_through" {
		t.Fatalf("SetBallClickThrough(false) = %v, %v, want [ball.click_through], nil", written, err)
	}
	if v := mgr.BallClickThroughOptions(); v.Enabled {
		t.Fatalf("BallClickThroughOptions().Enabled = true, want false after the write")
	}

	// panel.width: lands; non-positive is refused by name and nothing moves.
	if written, err := mgr.SetPanelWidth(500); err != nil || len(written) != 1 || written[0] != "panel.width" {
		t.Fatalf("SetPanelWidth(500) = %v, %v, want [panel.width], nil", written, err)
	}
	if v := mgr.PanelWidthOptions(); v.Width != 500 {
		t.Fatalf("PanelWidthOptions().Width = %d, want 500", v.Width)
	}
	_, err = mgr.SetPanelWidth(0)
	if err == nil || !strings.Contains(err.Error(), "panel.width") {
		t.Fatalf("a zero width must be refused by name, got: %v", err)
	}
	_, err = mgr.SetPanelWidth(-3)
	if err == nil || !strings.Contains(err.Error(), "panel.width") {
		t.Fatalf("a negative width must be refused by name, got: %v", err)
	}
	if v := mgr.PanelWidthOptions(); v.Width != 500 {
		t.Fatalf("after the refused widths Width = %d, want the untouched 500", v.Width)
	}

	// panel.font_size: the NEW key, same shapes as width.
	if written, err := mgr.SetPanelFontSize(15); err != nil || len(written) != 1 || written[0] != "panel.font_size" {
		t.Fatalf("SetPanelFontSize(15) = %v, %v, want [panel.font_size], nil", written, err)
	}
	if v := mgr.PanelFontSizeOptions(); v.FontSize != 15 {
		t.Fatalf("PanelFontSizeOptions().FontSize = %d, want 15", v.FontSize)
	}
	_, err = mgr.SetPanelFontSize(0)
	if err == nil || !strings.Contains(err.Error(), "panel.font_size") {
		t.Fatalf("a zero font size must be refused by name, got: %v", err)
	}
	if v := mgr.PanelFontSizeOptions(); v.FontSize != 15 {
		t.Fatalf("after the refused font size FontSize = %d, want the untouched 15", v.FontSize)
	}
}

func TestAppearanceWaveReadFacesAndTheFile(t *testing.T) {
	// The read faces answer the FILE through a fresh Manager: every planted
	// value arrives as the file's own answer, and the [0, 1] band rides the
	// schema's own constants (this case spells the band so a copied band
	// cannot pass by accident).
	dir := t.TempDir()
	mgr, err := NewManager(appearanceWavePlant(t, dir), nil)
	if err != nil {
		t.Fatalf("manager over the plant: %v", err)
	}
	ov := mgr.BallOpacityIdleOptions()
	if ov.Opacity != 0.7 || ov.Min != 0 || ov.Max != 1 {
		t.Fatalf("BallOpacityIdleOptions = %+v, want opacity 0.7 over the band [0, 1]", ov)
	}
	if mgr.BallClickThroughOptions().Enabled {
		t.Fatalf("BallClickThroughOptions().Enabled = true, want the planted false")
	}
	if v := mgr.PanelWidthOptions(); v.Width != 700 {
		t.Fatalf("PanelWidthOptions().Width = %d, want the planted 700", v.Width)
	}
	if v := mgr.PanelFontSizeOptions(); v.FontSize != 14 {
		t.Fatalf("PanelFontSizeOptions().FontSize = %d, want the planted 14", v.FontSize)
	}

	// A nil Manager answers "not configured" states with the band intact -
	// never a made-up value.
	nilView := (*Manager)(nil).BallOpacityIdleOptions()
	if nilView.Opacity != 0 || nilView.Min != 0 || nilView.Max != 1 {
		t.Fatalf("a nil Manager's opacity view = %+v, want the band with no opacity", nilView)
	}
	if (*Manager)(nil).BallClickThroughOptions().Enabled {
		t.Fatalf("a nil Manager must not answer Enabled=true - 'not configured' is a state, not a switch position")
	}
	if v := (*Manager)(nil).PanelFontSizeOptions(); v.FontSize != 0 {
		t.Fatalf("a nil Manager's font view = %d, want 0", v.FontSize)
	}

	// The schema default: a canonical file written from NewDefaults carries
	// font_size 12 - the panel document's de-facto body size - so the wiring
	// changes nothing at the default, and the round-trip keeps the key.
	def := NewDefaults()
	if def.Panel.FontSize != 12 {
		t.Fatalf("Panel.FontSize default = %d, want 12", def.Panel.FontSize)
	}
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := SaveFile(path, def); err != nil {
		t.Fatalf("SaveFile defaults: %v", err)
	}
	c, _, err := LoadFile(path, nil)
	if err != nil {
		t.Fatalf("LoadFile of the canonical file: %v", err)
	}
	if c.Panel.FontSize != 12 {
		t.Fatalf("the canonical file round-trips font_size %d, want 12", c.Panel.FontSize)
	}
}
