package panel

// Ticket 65's writable-field addition: ball_style, the orb appearance family,
// riding the settings route ticket 248 opened.
//
// WHY A FOURTH ROSTER AND NOT A ROW IN AN OLDER ONE. Ticket 248's table is
// pinned by a frozen contract test whose `want` list is the seven 248 names,
// exactly; 313's roster and 187's roster are pinned the same way, and by the
// disjointness assertions between them. The addition lives in its own roster,
// the door asks ALL FOUR (knownWritableField's union), and every older pin
// keeps meaning exactly what it meant.
//
// WIRING SPLIT, NOW CLOSED (docs/evidence/s1/65-style-switch-config.md §2/§3).
// Ticket 187's in-flight leg once owned uncommitted hunks in
// config_handlers.go, so this roster's admission rode a separate leg: the
// one-line union admission (`|| writableField65(field)` in
// knownWritableField, config_handlers.go) landed with the closing leg
// together with the snapshot key (ball_style_65.go's carrier, composer.go's
// field, pump.go's reader) and run.go's injection - the same commit family
// that carries cmd/wisp's ApplySetting case. The settings door admits
// ball_style from here on.
//
// THE TIER IS THE WRITE LEG'S ANSWER, NOT A STAMP HERE. 187 could carry a
// static per-field tier because its reader's behavior is fixed (per-task
// re-resolution). ball_style's tier is measured at write time by the process
// that hosts the ball: when the runtime applier is installed the write takes
// effect NOW (ball.ApplyStyle swaps the family and re-navigates the orb page);
// when this process has no ball there is nothing to apply it to and the
// honest answer is restart (the boot seed reads [ball].style at startup). No
// static map here, and no tier word invented in this file.
//
// NO APPROVAL VOCABULARY, NO LOCKED FAMILY. Same rules as every roster since
// 187: the field name carries neither, the locked-family guard before the
// older rosters is unchanged and runs first, the write cannot touch a
// permission档, and the config-side setter owns the vocabulary check.

import "sort"

// FieldBallStyle is ball.style, the orb appearance family. Value = the family
// word itself (glass|jarvis); the enum is enforced in internal/config before
// anything is written, the same way provider_billing's enum is.
const FieldBallStyle = "ball_style"

// FieldBallSize is ball.size, the orb's diameter in px. Value = the integer
// px as its string ("90"); the 44..160 band is enforced in internal/config
// before anything is written - SetBallSize refuses with validate.go's own
// ballSizeRangeError, the single-sourced sentence the load gate uses, so the
// write refusal and the load failure can never diverge. Ticket 422 grew this
// roster from one field to two (the 413 contract named the landing here, not
// in a fifth roster, and the disjointness pins stay untouched by the same
// argument the file comment above makes).
const FieldBallSize = "ball_size"

// fields65 is the roster this file owns, sorted at the accessor.
var fields65 = []string{
	FieldBallSize,
	FieldBallStyle,
}

// fields65TakesValue: both fields need the envelope's plain value half.
var fields65TakesValue = map[string]bool{
	FieldBallSize:  true,
	FieldBallStyle: true,
}

// WritableFields65 returns ticket 65's writable fields, sorted. It is the 65
// twin of WritableFields187, and a test reads the roster off here rather than
// copying it.
func WritableFields65() []string {
	out := append([]string(nil), fields65...)
	sort.Strings(out)
	return out
}

// writableField65 asks this file's roster, not a copy of it. The door's union
// (knownWritableField, config_handlers.go) consults it - the one place the
// four rosters are asked together.
func writableField65(field string) bool {
	for _, f := range fields65 {
		if f == field {
			return true
		}
	}
	return false
}
