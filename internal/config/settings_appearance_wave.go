package config

// The appearance-wave settings additions: the four keys the panel's config
// screen left locked (ball.opacity_idle, ball.click_through, panel.width and
// the new panel.font_size), each getting the guarded-write pair the 65/413
// shape established - one setter on the writeOneKey path and one read face
// the snapshot carrier reads.
//
// WHAT ALREADY EXISTED: the [ball] keys are schema originals (schema.go's
// BallSection, opacity validated 0..1 by validate.go), and panel.width is a
// schema original with a production reader since ticket 255 AC#4
// (panelGeometrySource, per window creation). panel.font_size is NEW: the
// schema key lands with this leg because nothing could read or write a key
// that did not exist, with the default (12) taken from what the panel
// document's body text de-facto renders at, so the wiring changes nothing at
// the default. The contract-text note on the schema field applies here too:
// PLAN.md's [panel] row does not name the key yet, and that row is not this
// leg's to edit.
//
// SAME SHAPE AS THE SETTERS BEFORE IT (settings_413.go says it best): one
// owned key path, writeOneKey, validate() on the base the write would
// produce, refusal before a byte moves, receipt names the key path the FILE
// changed by. opacity_idle's band is checked HERE and not only at load, the
// same ruling 320/65/413 made for their enums and bands: the refusal is
// validate.go's OWN ballOpacityIdleRangeError, never a second spelling.
// width and font_size carry no load-time band (the schema has none to
// single-source), so their setters refuse the one shape that is meaningless
// in every consumer - a non-positive value - with a sentence naming the
// field, instead of inventing a band nobody froze.
//
// THE TIER IS NOT DECIDED HERE. internal/config must not import internal/ball
// or the panel host; the write legs answer the measured tier the way the
// style leg did (cmd/wisp/panel_config_store.go). The readers below answer
// what the live Manager holds right now - the file and memory are kept equal
// by the guarded writer, and a snapshot must show what the next write would
// act on.

import (
	"fmt"

	"github.com/CarlosShao/wisp/internal/observe"
)

const (
	keyPathBallOpacityIdle  = "ball.opacity_idle"
	keyPathBallClickThrough = "ball.click_through"
	keyPathPanelWidth       = "panel.width"
	keyPathPanelFontSize    = "panel.font_size"
)

// SetBallOpacityIdle writes ball.opacity_idle through the same guarded
// key-by-key path every settings setter uses. The value is checked against
// the schema's own [0, 1] band BEFORE the write (the same helper validateBall
// refuses a hand-edited file with); anything outside is refused with the
// field's name and the band, leaving config.toml byte-identical.
func (m *Manager) SetBallOpacityIdle(v float64) ([]string, error) {
	if v < 0 || v > 1 {
		return nil, ballOpacityIdleRangeError(v)
	}
	return m.writeOneKey(keyPathBallOpacityIdle, func(base *Config) bool {
		base.Ball.OpacityIdle = v
		return true
	}, nil)
}

// SetBallClickThrough writes ball.click_through. A bool has no band to
// single-source; the guarded writer's validate() on the merged base is the
// gate.
func (m *Manager) SetBallClickThrough(on bool) ([]string, error) {
	return m.writeOneKey(keyPathBallClickThrough, func(base *Config) bool {
		base.Ball.ClickThrough = on
		return true
	}, nil)
}

// SetPanelWidth writes panel.width. A non-positive width is refused by name:
// the panel host builds its window edge from this number, and zero would ask
// Win32 for a degenerate window - refusing beats silently clamping to a band
// nobody froze in the schema.
func (m *Manager) SetPanelWidth(px int) ([]string, error) {
	if px <= 0 {
		return nil, observe.New(observe.ClassConfig, fmt.Sprintf(
			"config: panel.width: %d 不是有效的正整数 px 宽度；这一枚只收正整数", px))
	}
	return m.writeOneKey(keyPathPanelWidth, func(base *Config) bool {
		base.Panel.Width = px
		return true
	}, nil)
}

// SetPanelFontSize writes panel.font_size (the key this wave's schema leg
// added). Same refusal shape as SetPanelWidth: a non-positive font size is
// meaningless in every consumer, so it is refused by name rather than
// clamped against a band nobody froze.
func (m *Manager) SetPanelFontSize(px int) ([]string, error) {
	if px <= 0 {
		return nil, observe.New(observe.ClassConfig, fmt.Sprintf(
			"config: panel.font_size: %d 不是有效的正整数 px 字号；这一枚只收正整数", px))
	}
	return m.writeOneKey(keyPathPanelFontSize, func(base *Config) bool {
		base.Panel.FontSize = px
		return true
	}, nil)
}

// BallOpacityIdleView is the read face the opacity row needs: the idle
// opacity the live config holds right now, and the [0, 1] band it may offer
// - validate.go's own band, the same constants the load gate judges a
// hand-edited file by and the setter refuses a write by.
type BallOpacityIdleView struct {
	Opacity float64
	Min     float64
	Max     float64
}

// BallOpacityIdleOptions is the production reader for the opacity row
// (the reader-gated shape BallStyleOptions documents). A nil Manager answers
// an empty view with the band still intact.
func (m *Manager) BallOpacityIdleOptions() BallOpacityIdleView {
	v := BallOpacityIdleView{Min: 0, Max: 1}
	if m == nil {
		return v
	}
	c := m.Config()
	if c == nil {
		return v
	}
	v.Opacity = c.Ball.OpacityIdle
	return v
}

// BallClickThroughView is the read face the click-through row needs: the
// switch the live config holds right now.
type BallClickThroughView struct {
	Enabled bool
}

// BallClickThroughOptions is the production reader for the switch row. A nil
// Manager answers Enabled=false - "not configured" is a state, not a made-up
// switch position.
func (m *Manager) BallClickThroughOptions() BallClickThroughView {
	if m == nil {
		return BallClickThroughView{}
	}
	c := m.Config()
	if c == nil {
		return BallClickThroughView{}
	}
	return BallClickThroughView{Enabled: c.Ball.ClickThrough}
}

// PanelWidthView is the read face the width row needs: the width the live
// config holds right now, in px. 0 only when no Manager backs the read.
type PanelWidthView struct {
	Width int
}

// PanelWidthOptions is the production reader for the width row.
func (m *Manager) PanelWidthOptions() PanelWidthView {
	v := PanelWidthView{}
	if m == nil {
		return v
	}
	c := m.Config()
	if c == nil {
		return v
	}
	v.Width = c.Panel.Width
	return v
}

// PanelFontSizeView is the read face the font-size row needs: the body font
// size in px the live config holds right now.
type PanelFontSizeView struct {
	FontSize int
}

// PanelFontSizeOptions is the production reader for the font-size row.
func (m *Manager) PanelFontSizeOptions() PanelFontSizeView {
	v := PanelFontSizeView{}
	if m == nil {
		return v
	}
	c := m.Config()
	if c == nil {
		return v
	}
	v.FontSize = c.Panel.FontSize
	return v
}
