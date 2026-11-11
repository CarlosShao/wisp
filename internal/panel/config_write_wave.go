package panel

// The appearance wave's writable-field roster: ball_opacity_idle,
// ball_click_through, panel_width and panel_font_size, riding the settings
// route ticket 248 opened.
//
// WHY A FIFTH ROSTER AND NOT A ROW IN AN OLDER ONE. Ticket 248's table is
// pinned by a frozen contract test whose `want` list is the seven 248 names,
// exactly; 313's, 187's and 65's rosters are pinned the same way, and by the
// disjointness assertions between them. The additions live in their own
// roster, the door asks ALL FIVE (knownWritableField's union), and every
// older pin keeps meaning exactly what it meant.
//
// NO APPROVAL VOCABULARY, NO LOCKED FAMILY. Same rules as every roster since
// 187: the field names carry neither, the locked-family guard before the
// older rosters is unchanged and runs first, the write cannot touch a
// permission档, and the config-side setter owns the value check.
//
// THE TIERS ARE THE WRITE LEG'S ANSWER, NOT A STAMP HERE. ball_style's tier
// is measured at write time (the applier's presence); the same is true of
// this wave's two [ball] rows - the idle-opacity and click-through appliers
// live in cmd/wisp where the ball is, so a write answers now when the ball
// is hosted and restart when it is not. The two [panel] rows have NO live
// applier on purpose and honestly: the panel host re-reads [panel] at every
// window creation (and a re-show posts the resolved pair once more,
// ticket 255-r1), so a written width/font size follows on the next show
// request - no process restart is needed, and the leg says so in the
// receipt's Note instead of stamping a tier word the vocabulary does not
// carry.

import "sort"

// The four field names, spelled the way the config screen's rows spell their
// keys (section_key). Values: the two [ball] rows take the plain value half
// ("0.9" / "true"), the two [panel] rows take the integer px as its string.
const (
	FieldBallOpacityIdle  = "ball_opacity_idle"
	FieldBallClickThrough = "ball_click_through"
	FieldPanelWidth       = "panel_width"
	FieldPanelFontSize    = "panel_font_size"
)

// fieldsWave is the roster this file owns, sorted at the accessor.
var fieldsWave = []string{
	FieldBallClickThrough,
	FieldBallOpacityIdle,
	FieldPanelFontSize,
	FieldPanelWidth,
}

// fieldsWaveTakesValue: all four fields need the envelope's plain value half.
var fieldsWaveTakesValue = map[string]bool{
	FieldBallOpacityIdle:  true,
	FieldBallClickThrough: true,
	FieldPanelWidth:       true,
	FieldPanelFontSize:    true,
}

// WritableFieldsWave returns the appearance wave's writable fields, sorted.
// It is the wave twin of WritableFields187/WritableFields65, and a test reads
// the roster off here rather than copying it.
func WritableFieldsWave() []string {
	out := append([]string(nil), fieldsWave...)
	sort.Strings(out)
	return out
}

// writableFieldWave asks this file's roster, not a copy of it. The door's
// union (knownWritableField, config_handlers.go) consults it - the one place
// the five rosters are asked together.
func writableFieldWave(field string) bool {
	for _, f := range fieldsWave {
		if f == field {
			return true
		}
	}
	return false
}
