package config

// Ticket 413's settings additions for [ball].size: the guarded setter and the
// read face the size slider needs, same shape as settings_65.go's style pair.
//
// WHAT ALREADY EXISTED (the ticket's own census, re-measured 2026-11-10): the
// schema key (schema.go:181 Size with default:"56"), the load-time range gate
// (validate.go's validateBall, 44..160, ceiling widened 72 -> 160 with owner
// approval A1636), and the ball option the value rides
// (internal/ball/ball_windows.go:64 SizePx, fed at boot by cmd/wisp's
// ballSizeSeed413 -> withBallSizePx -> SizePx since the 413 AC#1 leg). What did
// NOT exist was anything that could WRITE the key: internal/panel's writable
// roster carries ball_style only (config_write_65.go:39-46), the settings
// write store has no ball.size case (cmd/wisp/panel_config_store.go's
// ApplySetting), and the CLI's closed roster names ball.style only
// (cmd/wisp/config_cmd_65.go:61). A slider with no writer is the half road the
// ticket forbids, so the setter lives here.
//
// SAME SHAPE AS THE SETTERS BEFORE IT (settings_65.go says it best, and
// settings_320.go before that): one owned key path, writeOneKey, validate() on
// the base the write would produce, refusal before a byte moves, receipt names
// the key path the FILE changed by. The range is checked HERE and not only at
// load, the same ruling ticket 320/65 made for their enums: an out-of-band
// slider value must come back as a sentence naming the field and the band -
// and it is validate.go's OWN sentence, because both gates share
// ballSizeRangeError (validate.go), never a second spelling of the band.
//
// THE TIER IS NOT DECIDED HERE, AND THE WINDOW IS NOT RESIZED HERE.
// internal/config must not import internal/ball, and the resident ball's
// window geometry is built once at boot from the seed read (ballSizeSeed413);
// a written size reaches memory hot ([ball] is hot-tier, D36) but the visible
// diameter changes when the ball window is next built (restart). Whoever wires
// the panel route answers the measured tier the way the style leg did
// (cmd/wisp/panel_config_store.go's FieldBallStyle case) - this file only
// guards and writes the key.

// keyPathBallSize is the one dotted key path SetBallSize owns.
const keyPathBallSize = "ball.size"

// SetBallSize writes ball.size through the same guarded key-by-key path every
// settings setter uses. The value is checked against the schema's own band
// BEFORE the write (the same helper validateBall refuses a hand-edited file
// with); anything outside 44..160 is refused with the field's name and the
// band, leaving config.toml byte-identical.
func (m *Manager) SetBallSize(size int) ([]string, error) {
	if size < ballSizeMinPx || size > ballSizeMaxPx {
		return nil, ballSizeRangeError(size)
	}
	return m.writeOneKey(keyPathBallSize, func(base *Config) bool {
		base.Ball.Size = size
		return true
	}, nil)
}

// BallSizeView is the read face a size slider needs: the diameter the live
// config names right now, and the band it may offer. Reads the in-memory
// config (Manager.Config()) the same way BallStyleView does - the file and
// memory are kept equal by the guarded writer, and a reader must answer with
// what this process would act on.
type BallSizeView struct {
	// Size is the current ball.size value in px. 0 only when no Manager backs
	// the read (the caller answers "not configured", never a made-up diameter).
	Size int
	// Min and Max are the schema's own 44..160 band, the same constants
	// validate.go judges a load by and SetBallSize refuses a write by - a
	// slider offers exactly the band load accepts, never a copy of it.
	Min int
	// Max see Min.
	Max int
}

// BallSizeOptions is the production reader for the size slider (the
// reader-gated shape BallStyleOptions documents). A nil Manager answers Size 0
// with the band still intact.
func (m *Manager) BallSizeOptions() BallSizeView {
	v := BallSizeView{Min: ballSizeMinPx, Max: ballSizeMaxPx}
	if m == nil {
		return v
	}
	c := m.Config()
	if c == nil {
		return v
	}
	v.Size = c.Ball.Size
	return v
}
