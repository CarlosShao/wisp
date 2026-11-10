package panel

// Ticket 413's carrier, landed by ticket 422's panel wiring: the snapshot's
// ballSize key - the orb diameter the live config holds and the band a size
// slider may offer (docs/evidence/s1/413-ball-size-slider-contract.md §2.1,
// the contract the 413 leg named for exactly this leg).
//
// READER-GATED, THE 213/187/65 PATTERN, NO SECOND PARADIGM. A pump assembled
// without the reader sends no "ballSize" key at all - the shape the four-key
// byte nails pin - and a pump that HAS one always sends the key. Nothing here
// opens a door: the write half rides config.set's ball_size field
// (config_write_65.go, admitted by knownWritableField's four-roster union),
// and NO C17 method is added by this carrier.
//
// THE VALUES ARE THE SCHEMA'S, NOT A COPY. The production reader is
// internal/config's BallSizeOptions, injected by run.go's assembly root; the
// band it returns is validate.go's own 44..160 (ballSizeMinPx/ballSizeMaxPx),
// the same constants the load gate judges a hand-edited file by and
// SetBallSize refuses a write by. This file spells neither 44 nor 160.

// BallSizeSection is the snapshot's ballSize key.
type BallSizeSection struct {
	// Size is ball.size in px as the live config holds it RIGHT NOW - the
	// diameter the next boot seeds the orb window with (the window geometry
	// is built once at boot from the seed read; there is no in-process
	// re-size, so the write route's honest tier answer is restart). 0 only
	// when no Manager backs the read; the page answers "not configured"
	// rather than being handed an invented diameter.
	Size int `json:"size"`
	// Min is the schema's own band floor (validate.go's ballSizeMinPx), the
	// same constant the load gate judges by and SetBallSize refuses a write
	// by - a slider offers exactly the band load accepts, never a copy of it.
	Min int `json:"min"`
	// Max see Min.
	Max int `json:"max"`
}
