package panel

// Ticket 193 AC#3's writable-field addition: panel_animation, the animation
// switch's write half, riding the settings route ticket 248 opened.
//
// WHY A NINTH ROSTER AND NOT A ROW IN AN OLDER ONE. Same argument every
// roster since 313 carries (config_write_215.go spells it out): the older
// rosters are each pinned by their own frozen contract tests with exact
// want-lists, so an addition lives in its own roster, the door asks ALL of
// them (knownWritableField's union), and every older pin keeps meaning exactly
// what it meant. A881's union pin (cmd/wisp/panel_config_313_test.go) grew the
// same one-append way, and this roster's cmd/wisp half (the ApplySetting case
// in panel_config_store.go) lands in the SAME commit as this file - a roster
// without its leg is "门 admitting、腿拒绝" (the 193 evidence file §③'s own
// wording), and a leg without its roster is refused at the door.
//
// THE ONE FIELD AND ITS TIER. panel_animation is [panel].animation (the
// schema key ticket 193's config leg added, default "reduce" - the
// owner-pending placeholder). The field name is flat (the dotted spelling is
// not admitted here, same rule as plugin_enabled); the value is the tier word
// itself, and the vocabulary check owns exactly ONE source:
// internal/config's SetPanelAnimation refuses out-of-vocabulary words with
// validate.go's own panelAnimationEnumError - the same sentence the load gate
// refuses a hand-edited file with (the 65/413 single-sentence shape). Nothing
// here spells full, reduce or off. The tier the receipt carries is asked of
// internal/config's TierRegistry by the leg (config_tier_39.go's mapping, the
// whole-section "panel" row = hot, D36's [panel] 行) - this door adds no tier
// vocabulary of its own.

import "sort"

// FieldPanelAnimation is panel.animation, the panel page's animation switch.
// Hot tier: the write lands in the process's memory, and the snapshot pump's
// reader re-reads the value per publish, so the page's next packet runs the
// new tier with no process restart.
const FieldPanelAnimation = "panel_animation"

// fields193 is the roster this file owns, sorted at the accessor.
var fields193 = []string{
	FieldPanelAnimation,
}

// fields193TakesValue: the field needs the plain value half (the tier word).
var fields193TakesValue = map[string]bool{
	FieldPanelAnimation: true,
}

// WritableFields193 returns ticket 193's writable fields, sorted. The 193 twin
// of the other rosters' accessors; a test reads the roster off here rather than
// copying it.
func WritableFields193() []string {
	out := append([]string(nil), fields193...)
	sort.Strings(out)
	return out
}

// writableField193 asks this file's roster, not a copy of it. The door's union
// (knownWritableField, config_handlers.go) consults it.
func writableField193(field string) bool {
	for _, f := range fields193 {
		if f == field {
			return true
		}
	}
	return false
}
