package panel

// Ticket 65's carrier: the snapshot's ballStyle key - the orb appearance
// family the live config holds and the vocabulary an appearance picker may
// offer (the owner's ask, verbatim: 默认就是玻璃球，切换 Jarvis 才换球的风格).
//
// READER-GATED, THE 213/187 PATTERN, NO SECOND PARADIGM. A pump assembled
// without the reader sends no "ballStyle" key at all - the shape the four-key
// byte nails pin - and a pump that HAS one always sends the key. Nothing here
// opens a door: the write half rides config.set's ball_style field
// (config_write_65.go, admitted by knownWritableField's four-roster union),
// and NO C17 method is added by this carrier.
//
// THE VALUES ARE THE SCHEMA'S, NOT A COPY. The production reader is
// internal/config's BallStyleOptions, injected by run.go's assembly root; the
// vocabulary it returns is the schema's own ballStyles list, the same list
// validate.go judges a load by and SetBallStyle refuses a write by. This file
// spells neither glass nor jarvis.

// BallStyleSection is the snapshot's ballStyle key.
type BallStyleSection struct {
	// Style is ball.style as the live config holds it RIGHT NOW - the family
	// the next boot seeds the orb with, and (when this process hosts the ball
	// with the applier installed) the family already showing. Empty only when
	// no Manager backs the read; the page answers "not configured" rather
	// than being handed an invented family.
	Style string `json:"style"`
	// Vocabulary is the schema's own glass|jarvis list, in schema order, so
	// the picker offers exactly what load accepts and what the write route
	// would refuse anything outside of.
	Vocabulary []string `json:"vocabulary"`
}
