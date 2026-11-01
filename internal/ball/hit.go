package ball

import "math"

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

// RingMarginPx is the non-interactive rendering margin around the orb. It
// sizes the window but NOT the clickable circle, which is the orb body plus a
// small tolerance, extended outward to the halo's visible edge (R15#5).
//
// The margin has to contain the widest thing a frame draws, which is the
// Listening wave train at R+5px+16px of travel (+ a 2px badge beyond the rim):
// R+21.75. The frozen ball.html inset (-5, Conversation -7) sized this for the
// pre-ticket-62 ring, so the prototype's expanding waves were sliced by the
// window rect and the ball read as a framed square (owner 09-30 sign-off).
// 22 is the smallest value that fits the waves at every size in 44..72.
const RingMarginPx = 22.0

// ClickTolerancePx extends the clickable circle slightly beyond the orb so
// the drawn ring belongs to the hit region (DPI-scaled at call sites). It is
// the FLOOR of the clickable circle, not its whole story: since R15#5 (owner,
// 2026-11-01) the circle never stops inside the visible outer edge either -
// see HitTest.
const ClickTolerancePx = 4.0

// HitTest returns HTClient when (x, y) - window-client coordinates in
// physical pixels - lies within the clickable circle, HTTRANSPARENT
// otherwise. The circle is the larger of the orb body plus the click
// tolerance and the halo's visible outer edge (1.5*R, the fill radius
// haloExtents draws - unclipped at every supported size since RingMarginPx
// went to 22): a pixel the user can see must be clickable (R15#5 "对齐可见
// 外沿"), so the Sleeping rest body's 21px circle became the dispatched 26px.
// sizePx is the window edge in physical pixels, dpi the window's monitor DPI
// (per-monitor V2) and orbPx the configured orb size at 96 DPI. All inputs
// are injected: tests simulate any DPI/size without a window.
func HitTest(x, y int32, sizePx int32, orbPx int, dpi uint32) int {
	scale := float64(dpi) / 96.0
	r := float64(orbPx)*scale/2 + ClickTolerancePx*scale
	if edge := math.Ceil(float64(orbPx) * scale * 0.75); edge > r {
		r = edge
	}
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
