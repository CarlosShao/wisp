package config

// Ticket 193's animation switch, config half: the guarded-write pair the
// appearance-wave shape established (settings_appearance_wave.go says it
// best) - one setter on the writeOneKey path and one read face the snapshot
// carrier reads. panel.animation is the schema key this ticket's config leg
// added (schema.go's PanelSection.Animation, default "reduce" - the
// owner-pending placeholder; the ticket says 默认值由 owner 定 and nothing
// here may treat the placeholder as decided).
//
// SAME SHAPE AS THE SETTERS BEFORE IT: one owned key path, writeOneKey,
// validate() on the base the write would produce, refusal before a byte
// moves, receipt names the key path the FILE changed by. The tier vocabulary
// is checked HERE and not only at load, the same ruling 320/65/413 made for
// their enums: the refusal is validate.go's OWN panelAnimationEnumError,
// never a second spelling.
//
// THE TIER IS HOT (TierRegistry's whole-section "panel" row, tiers.go; D36's
// [panel] row). There is no Go-side applier on purpose and honestly - the
// same shape the appearance wave's [panel] rows ride: the value is readable
// hot through the readers below the moment a write lands, the snapshot
// carrier (internal/panel's PanelAnimationSection) re-reads per pump cycle,
// and the page consumes the tier with the frontend wiring leg. The route
// that turns a panel config.set write into SetPanelAnimation (cmd/wisp's
// ApplySetting case plus the ninth roster) is the wiring leg's to land; this
// file is the config seam it calls.

import "slices"

const keyPathPanelAnimation = "panel.animation"

// SetPanelAnimation writes panel.animation through the same guarded key-by-key
// path every settings setter uses. The value is checked against the schema's
// own vocabulary BEFORE the write (the same helper validatePanelAnimation
// refuses a hand-edited file with); anything else is refused with the field's
// name and the three words, leaving config.toml byte-identical.
func (m *Manager) SetPanelAnimation(tier string) ([]string, error) {
	if !slices.Contains(panelAnimationVocabulary, tier) {
		return nil, panelAnimationEnumError(tier)
	}
	return m.writeOneKey(keyPathPanelAnimation, func(base *Config) bool {
		base.Panel.Animation = tier
		return true
	}, nil)
}

// PanelAnimationView is the read face the animation row needs: the tier the
// live config holds right now, and the three-word vocabulary it may offer -
// validate.go's own slice, the same words the load gate judges a hand-edited
// file by and the setter refuses a write by.
type PanelAnimationView struct {
	Tier string
	// Options is the vocabulary in schema order: full, reduce, off.
	Options []string
}

// PanelAnimationOptions is the production reader for the animation row (the
// reader-gated shape PanelWidthOptions documents). A nil Manager answers an
// empty view with the vocabulary still intact - "not configured" is a state,
// not a made-up tier word.
func (m *Manager) PanelAnimationOptions() PanelAnimationView {
	v := PanelAnimationView{Options: append([]string(nil), panelAnimationVocabulary...)}
	if m == nil {
		return v
	}
	c := m.Config()
	if c == nil {
		return v
	}
	v.Tier = c.Panel.Animation
	return v
}
