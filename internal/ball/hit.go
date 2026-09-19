package ball

// Hit-testing / click-through (SPEC-08 §2: "透明像素区域 WS_EX_TRANSPARENT 或
// hit-test 返回 HTTRANSPARENT；球体本体可点").
//
// The ball window is a square sized orb+2*margin (WindowEdgePx) so the ring
// and glow never clip. Everything inside the INSCRIBED CIRCLE is the ball
// (body + ring zone) and clickable; the four corners are transparent and
// must fall through - WM_NCHITTEST returns HTTRANSPARENT there, which is the
// D2D-equivalent of WS_EX_TRANSPARENT without giving up per-pixel control.

// Windows hit-test codes (subset).
const (
	HTClient      = 1
	HTTransparent = -1
	HTNowhere     = 0
)

// RingMarginPx is the non-interactive rendering margin around the orb (the
// ring/glow zone; ball.html inset -5px, Conversation -7px). It sizes the
// window but NOT the clickable circle, which stays the orb body + a small
// tolerance so the visible ring still counts as "the ball".
const RingMarginPx = 8.0

// ClickTolerancePx extends the clickable circle slightly beyond the orb so
// the drawn ring belongs to the hit region (DPI-scaled at call sites).
const ClickTolerancePx = 4.0

// HitTest returns HTClient when (x, y) - window-client coordinates in
// physical pixels - lies within the clickable circle (orb radius scaled to
// the monitor DPI + tolerance), HTTRANSPARENT otherwise. sizePx is the
// window edge in physical pixels, dpi the window's monitor DPI (per-monitor
// V2) and orbPx the configured orb size at 96 DPI. All inputs are injected:
// tests simulate any DPI/size without a window.
func HitTest(x, y int32, sizePx int32, orbPx int, dpi uint32) int {
	scale := float64(dpi) / 96.0
	r := float64(orbPx)*scale/2 + ClickTolerancePx*scale
	cx, cy := float64(sizePx)/2, float64(sizePx)/2
	dx, dy := float64(x)-cx, float64(y)-cy
	if dx*dx+dy*dy <= r*r {
		return HTClient
	}
	return HTTransparent
}

// WindowEdgePx returns the physical window edge (orb + ring margins on both
// sides) for a configured orb size and DPI.
func WindowEdgePx(configuredOrbPx int, dpi uint32) int {
	scale := float64(dpi) / 96.0
	orb := float64(configuredOrbPx) * scale
	return int(orb + 2*RingMarginPx*scale + 0.5)
}

// SleepWindowEdgePx returns the window edge for the Sleeping micro dot. The
// dot is 12px but the window keeps a minimum usable hit area so the user can
// still find and click it (the extra area renders fully transparent and
// falls through to the desktop).
func SleepWindowEdgePx(dpi uint32) int {
	return WindowEdgePx(SleepingDotPx, dpi)
}
