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
// WIRING SPLIT, HONESTLY STATED (docs/evidence/s1/65-style-switch-config.md §2).
// Ticket 187's in-flight leg owns uncommitted hunks in config_handlers.go, so
// this leg must not edit that file's union - the one-line admission
// (`|| writableField65(field)` in knownWritableField) lands with the next leg,
// together with the snapshot key (composer.go/pump.go are the same conflict
// zone). Until then the settings door refuses ball_style, while the write leg
// underneath it (cmd/wisp's ApplySetting case and the guarded config setter)
// is already wired and exercised by its own tests.
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

// fields65 is the roster this file owns, sorted at the accessor.
var fields65 = []string{
	FieldBallStyle,
}

// fields65TakesValue: the one field needs the envelope's plain value half.
var fields65TakesValue = map[string]bool{
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
// does not consult it YET (see the wiring-split note above).
func writableField65(field string) bool {
	for _, f := range fields65 {
		if f == field {
			return true
		}
	}
	return false
}
