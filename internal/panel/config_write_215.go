package panel

// Ticket 215a's writable-field addition: plugin_enabled, the Tier-2 plugin
// switch, riding the settings route ticket 248 opened.
//
// WHY A SIXTH ROSTER AND NOT A ROW IN AN OLDER ONE. Same argument every
// roster since 313 carries: the older rosters are each pinned by their own
// frozen contract tests with exact want-lists, so an addition lives in its
// own roster, the door asks ALL of them (knownWritableField's union), and
// every older pin keeps meaning exactly what it meant.
//
// THE ONE FIELD AND ITS TWO SELECTORS.
//
//	plugin_enabled  plugins.<id>.enabled, ONE row per call. The plugin id
//	  rides the envelope's configPlugin selector (a NEW selector, added with
//	  this roster, the same move ticket 248 made when config.set needed
//	  provider/model selectors); the value is the plain "true"/"false".
//	  The dotted spelling ("plugins.clip.enabled") is NOT admitted here and
//	  never will be: it falls into the locked-family guard before this
//	  roster is even consulted, and that guard is where the family's
//	  security sentence lives. The flat name passes the family guard by
//	  construction while still naming exactly one guarded write.
//
// WHERE THE SECURITY LIVES. The panel's gate (this roster + the family
// guard) only decides WHO MAY ASK; the verdict that decides what the write
// MAY DO is internal/config's section-write entry, reached through
// Manager.SetPluginEnabled: switching OFF is a tighten and writes
// immediately; switching ON is a LOOSENING of the locked [plugins] section
// and is refused for every panel-sourced attempt (the panel is a permission
// input, R20 - it can ask the native side for a loosening, never perform
// one). The toggle is audited at both layers: the CONFIG-ROUTE line here
// (record, every outcome) and the SECTION-WRITE / PLUGIN-TOGGLE lines in
// internal/config (every outcome, including the neutral one).
//
// NO APPROVAL VOCABULARY, NO LOCKED FAMILY. Same rules as every roster
// since 187: the field name carries neither (asserted by the same frozen
// instruments), and the config-side setter owns the direction verdict.

import "sort"

// FieldPluginEnabled is plugins.<id>.enabled, the Tier-2 plugin switch.
// Value = "true"/"false"; the plugin id rides configPlugin. The direction
// verdict (enable = loosening) is internal/config's, not this door's.
const FieldPluginEnabled = "plugin_enabled"

// fields215 is the roster this file owns, sorted at the accessor.
var fields215 = []string{
	FieldPluginEnabled,
}

// fields215NeedsPlugin: the selector column this roster introduces. The
// plugin id is the write's address - without it the gate cannot tell which
// row the page is toggling, and guessing one is not an option.
var fields215NeedsPlugin = map[string]bool{
	FieldPluginEnabled: true,
}

// fields215TakesValue: the field needs the plain value half (the on/off
// word); the bool shape is checked by the leg, not here.
var fields215TakesValue = map[string]bool{
	FieldPluginEnabled: true,
}

// WritableFields215 returns ticket 215a's writable fields, sorted. The 215
// twin of the other rosters' accessors; a test reads the roster off here
// rather than copying it.
func WritableFields215() []string {
	out := append([]string(nil), fields215...)
	sort.Strings(out)
	return out
}

// writableField215 asks this file's roster, not a copy of it. The door's
// union (knownWritableField, config_handlers.go) consults it.
func writableField215(field string) bool {
	for _, f := range fields215 {
		if f == field {
			return true
		}
	}
	return false
}
