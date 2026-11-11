package panel

// The appearance wave's carriers: the snapshot's ballOpacityIdle and
// ballClickThrough keys - the two [ball] rows the config screen left locked
// ("那颗球是原生绘制的（D29），页面改不到"), now with the read half this
// page's rows need, the same disk BallStyleSection/BallSizeSection read
// (docs/evidence/s1/413-ball-size-slider-contract.md §2.1, the shape the 65
// and 413 legs established).
//
// READER-GATED, THE 213/187/65/413 PATTERN, NO SECOND PARADIGM. A pump
// assembled without the reader sends no key at all; a pump that HAS one
// always sends it. Nothing here opens a door: the write halves ride
// config.set's ball_opacity_idle / ball_click_through fields
// (config_write_wave.go, admitted by knownWritableField's five-roster union),
// and NO C17 method is added by these carriers.
//
// THE VALUES ARE THE SCHEMA'S, NOT A COPY. The production readers are
// internal/config's BallOpacityIdleOptions / BallClickThroughOptions,
// injected by the assembly roots; the [0, 1] band they return is
// validate.go's own (ballOpacityIdleRangeError), the same sentence the load
// gate judges a hand-edited file by and the setter refuses a write by. This
// file spells neither 0 nor 1 as a band of its own.

// BallOpacityIdleSection is the snapshot's ballOpacityIdle key.
type BallOpacityIdleSection struct {
	// Opacity is ball.opacity_idle as the live config holds it RIGHT NOW -
	// the idle-state opacity the ball's own option carries (the native idle
	// rendering lands with this wave's ball leg; the write route's honest
	// tier answer is now when this process hosts the ball and restart when
	// it does not). 0 only when no Manager backs the read; the page answers
	// "not configured" rather than being handed an invented opacity.
	Opacity float64 `json:"opacity"`
	// Min is the schema's own band floor (validate.go's [0, 1]), the same
	// bound the load gate judges by and the setter refuses a write by - a
	// slider offers exactly the band load accepts, never a copy of it.
	Min float64 `json:"min"`
	// Max see Min.
	Max float64 `json:"max"`
}

// BallClickThroughSection is the snapshot's ballClickThrough key.
type BallClickThroughSection struct {
	// Enabled is ball.click_through as the live config holds it RIGHT NOW -
	// whether the idle ball passes mouse events through (the WS_EX_TRANSPARENT
	// toggle the native ball applies in its idle state). false is also what a
	// nil-backed read answers; "the config says false" and "nobody could read
	// it" are distinguished by the row's own no-reader shape (a pump without
	// the reader sends no key at all).
	Enabled bool `json:"enabled"`
}
