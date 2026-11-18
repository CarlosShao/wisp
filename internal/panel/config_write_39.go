package panel

// Ticket 39 AC#1's writable-field addition: app_language, the settings route's
// first RESTART-tier key, riding the settings door ticket 248 opened.
//
// WHY A SEVENTH ROSTER AND NOT A ROW IN AN OLDER ONE. Same argument every
// roster since 313 carries (config_write_215.go spells it out): the older
// rosters are each pinned by their own frozen contract tests with exact
// want-lists, so an addition lives in its own roster, the door asks ALL of
// them (knownWritableField's union), and every older pin keeps meaning exactly
// what it meant.
//
// THE ONE FIELD AND ITS TIER. app_language is [app].language (tiers.go's
// per-key "app.language": TierRestart row). The field name is flat (the dotted
// spelling is not admitted here, same rule as plugin_enabled); the value is
// the UI language tag as a plain string. The tier the receipt and the reload
// event carry is asked of internal/config's TierRegistry by the leg
// (cmd/wisp's registryEffectiveTier39) - this door adds no tier vocabulary of
// its own. The write's restart BEHAVIOR (file written, running memory keeps
// the boot value, OnRestartPending says 需要重启 on the next reload) is
// internal/config's tier-aware write path, not this door's.

import "sort"

// FieldAppLanguage is app.language, the UI language tag. Restart tier: a save
// lands on disk and the new tag is read at the next process start.
const FieldAppLanguage = "app_language"

// fields39 is the roster this file owns, sorted at the accessor.
var fields39 = []string{
	FieldAppLanguage,
}

// fields39TakesValue: the field needs the plain value half (the language tag).
var fields39TakesValue = map[string]bool{
	FieldAppLanguage: true,
}

// WritableFields39 returns ticket 39's writable fields, sorted. The 39 twin of
// the other rosters' accessors; a test reads the roster off here rather than
// copying it.
func WritableFields39() []string {
	out := append([]string(nil), fields39...)
	sort.Strings(out)
	return out
}

// writableField39 asks this file's roster, not a copy of it. The door's union
// (knownWritableField, config_handlers.go) consults it.
func writableField39(field string) bool {
	for _, f := range fields39 {
		if f == field {
			return true
		}
	}
	return false
}
