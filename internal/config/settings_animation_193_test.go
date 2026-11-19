package config

// Ticket 193's config-side rulers: panel.animation lands as a real schema key
// (default "reduce" - the owner-pending placeholder, NOT a decided value),
// an unknown tier word is a load error by name, the guarded setter and the
// load gate hand out the SAME single-sourced refusal sentence, a written tier
// is readable hot in the same process (the [panel] section's TierRegistry row
// is hot), and the read face carries the vocabulary the page may offer.
//
// The reverse criterion this test names: deleting the setter body or the
// validatePanelAnimation walk reddens here; deleting the cmd/wisp ApplySetting
// case (the wiring leg's work) reddens the store half of the wire test there.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// animationPlant builds a minimal valid config.toml carrying panel.animation
// non-default, so a default leaking into any reading cannot fake the answer
// (the appearance wave's plant shape, one [panel] key each).
func animationPlant(t *testing.T, tier string) string {
	t.Helper()
	body := "[panel]\nwidth = 700\nanimation = \"" + tier + "\"\n"
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write the animation plant (%q): %v", tier, err)
	}
	return path
}

func TestPanelAnimationDefaultIsTheOwnerPendingPlaceholder(t *testing.T) {
	def := NewDefaults()
	if def.Panel.Animation != "reduce" {
		t.Fatalf("Panel.Animation default = %q, want %q - the dispatch-ordered most conservative PLACEHOLDER", def.Panel.Animation, "reduce")
	}
	// The placeholder is NOT a decision: the ticket says 默认值由 owner 定 and
	// the evidence file carries that line verbatim. If the owner ever picks a
	// different default, this test and the schema tag move in the same commit
	// and the evidence file's owner-pending line is retired with them.
	v := (*Manager)(nil).PanelAnimationOptions()
	if v.Tier != "" || len(v.Options) != 3 {
		t.Fatalf("nil-Manager read = %+v, want Tier %q with the three-word vocabulary intact - not configured is a state, not a tier", v, "")
	}
}

func TestPanelAnimationLoadRefusesUnknownTierByTheSingleSentence(t *testing.T) {
	for _, word := range []string{"loud", "OFF", "none", ""} {
		plant := animationPlant(t, word)
		if _, _, err := LoadFile(plant, nil); err == nil {
			t.Fatalf("animation = %q must fail at load, not fall back to a tier", word)
		} else if got := err.Error(); !strings.Contains(got, "panel.animation") {
			t.Fatalf("the %q load refusal does not name panel.animation: %v", word, err)
		}
	}
}

func TestPanelAnimationLoadsOffAndReadsHot(t *testing.T) {
	plant := animationPlant(t, "off")
	mgr, err := NewManager(plant, nil)
	if err != nil {
		t.Fatalf("manager on the off plant: %v", err)
	}
	v := mgr.PanelAnimationOptions()
	if v.Tier != "off" {
		t.Fatalf("PanelAnimationOptions().Tier = %q, want %q read straight off the file", v.Tier, "off")
	}
	if strings.Join(v.Options, ",") != "full,reduce,off" {
		t.Fatalf("the vocabulary the row may offer = %v, want the schema's own three words in order", v.Options)
	}
	// The tier word the snapshot carrier sends is the config's own word: the
	// Go side derives nothing and schedules nothing (panel_animation_193.go).
	if tier, ok := TierOf("panel"); !ok || tier != "hot" {
		t.Fatalf("TierOf(panel) = %q, %v - the switch rides the whole-section hot row, a written tier must not need a restart", tier, ok)
	}
}

func TestSetPanelAnimationWritesAndRefusesByTheSingleSentence(t *testing.T) {
	dir := t.TempDir()
	seedPath := filepath.Join(dir, "config.toml")
	if err := SaveFile(seedPath, NewDefaults()); err != nil {
		t.Fatalf("SaveFile the defaults seed: %v", err)
	}
	mgr, err := NewManager(seedPath, nil)
	if err != nil {
		t.Fatalf("manager on the seed: %v", err)
	}

	// A real tier lands; the receipt names the key path the FILE changed by
	// and the memory view follows (the guarded writer keeps file and memory
	// equal - the hot tier's read half).
	written, err := mgr.SetPanelAnimation("off")
	if err != nil || len(written) != 1 || written[0] != "panel.animation" {
		t.Fatalf("SetPanelAnimation(off) = %v, %v, want [panel.animation], nil", written, err)
	}
	if v := mgr.PanelAnimationOptions(); v.Tier != "off" {
		t.Fatalf("after the write Tier = %q, want %q - the write must be readable hot in this same process", v.Tier, "off")
	}
	postWrite, err := os.ReadFile(seedPath)
	if err != nil {
		t.Fatalf("read the post-write bytes: %v", err)
	}

	// An unknown word is refused by the validate.go sentence VERBATIM and the
	// file's bytes stay exactly as the accepted write left them.
	_, err = mgr.SetPanelAnimation("loud")
	if err == nil {
		t.Fatalf("loud must be refused by the setter's vocabulary guard")
	}
	if got := err.Error(); !strings.Contains(got, "panel.animation") || !strings.Contains(got, "full") {
		t.Fatalf("the loud refusal is not the single-sourced validate sentence, got: %v", err)
	}
	afterRefusal, err := os.ReadFile(seedPath)
	if err != nil {
		t.Fatalf("read the post-refusal bytes: %v", err)
	}
	if string(afterRefusal) != string(postWrite) {
		t.Fatalf("the refused write moved the file - a refusal moves no byte")
	}
	if strings.Contains(string(afterRefusal), "'loud'") {
		t.Fatalf("the refused word reached the file: animation = 'loud' appears in config.toml")
	}
	if v := mgr.PanelAnimationOptions(); v.Tier != "off" {
		t.Fatalf("after the refused write Tier = %q, want the untouched %q", v.Tier, "off")
	}
}
