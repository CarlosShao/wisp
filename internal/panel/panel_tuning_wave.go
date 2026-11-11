package panel

// The appearance wave's carriers, panel half: the snapshot's panelWidth and
// panelFontSize keys - the two [panel] rows the config screen left locked
// ("窗口尺寸在宿主手里" / "快照里没有这个字段"), now with the read halves
// this page's rows need. The width half reads what the panel host's
// per-creation geometry source has been reading since ticket 255 AC#4
// (cfg.Panel.Width); the font-size half reads the NEW schema key the same
// wave's config leg added (panel.font_size, default 12 - the panel
// document's de-facto body size, so the wiring changes nothing at the
// default).
//
// READER-GATED, THE 213/187/65/413 PATTERN, NO SECOND PARADIGM. A pump
// assembled without the reader sends no key at all; a pump that HAS one
// always sends it. The write halves ride config.set's panel_width /
// panel_font_size fields (config_write_wave.go, admitted by
// knownWritableField's five-roster union), and NO C17 method is added by
// these carriers.
//
// THE WIDTH KEY CARRIES NO BAND, ON PURPOSE. The schema freezes no range for
// panel.width (and this carrier invents none): the page's control offers
// what the write route actually accepts (any positive px value), and the
// panel host's own fallback constants stay where they have always lived.

// PanelWidthSection is the snapshot's panelWidth key.
type PanelWidthSection struct {
	// Width is panel.width in px as the live config holds it RIGHT NOW - the
	// number the panel host's geometry source re-reads at every window
	// creation (and that a re-show posts once more, ticket 255-r1), so a
	// written width follows on the next show request without a process
	// restart. 0 only when no Manager backs the read.
	Width int `json:"width"`
}

// PanelFontSizeSection is the snapshot's panelFontSize key.
type PanelFontSizeSection struct {
	// FontSize is panel.font_size in px as the live config holds it RIGHT
	// NOW - the body font size the panel page consumes with the frontend
	// wiring leg. 0 only when no Manager backs the read.
	FontSize int `json:"fontSize"`
}
