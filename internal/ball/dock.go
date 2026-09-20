package ball

// Edge-dock geometry (ticket 62 item 2, owner: "球在靠近电脑侧边的时候，应该会
// 自动收缩，就像迅雷那种悬浮球一样"). Pure math on rects - no Win32, no clock -
// so any monitor layout, DPI and drag position is testable without a window
// (the pattern position.go established).
//
// The behavior, and where its frames come from:
//
//   - dragging toward a work-area edge squeezes the orb along that axis as it
//     goes (dockProximity: the ramp is read off the gap, so the squeeze is
//     continuous while the user pushes);
//   - releasing at the edge commits the dock and the squashed orb becomes a tab
//     hugging the boundary, with only a hittable strip left showing;
//   - hovering the tab slides the ball back out to a full orb (owner: "悬停或
//     单击 → 丝滑弹回完整球"), and a click on the tab pops it out and keeps it
//     out until the next drag;
//   - the pointer leaving lets the tab close again.
//
// Every step above is carried by a mouse message or a window-position event the
// user is generating anyway. The dock therefore owns NO timer at all, resident
// or bounded, which is what keeps a docked Sleeping orb inside D32's
// zero-timer/zero-CPU promise (ticket 62 constraint: "靠边吸附不得引入常驻定时
// 器：用 WM_NCHITTEST + 窗口位置事件驱动；悬停检测走鼠标消息"). The one place a
// transition lands in a single frame is the retreat after WM_MOUSELEAVE, where
// by definition no further message is coming; that is the price of the promise
// and it is paid on an 8-21px slide with the pointer already gone.

import "time"

// dockRampFrame is the ramp's step law, shared verbatim by every message path
// that moves the tab: the hover pop-out (WM_MOUSEMOVE -> dockHoverMove ->
// dockStep), the retreat after WM_MOUSELEAVE (dockLeave -> dockStep) and the
// drag-end commit. It is pure - level in, level out, dt is the caller's own
// frame-spacing measurement - so the reverse direction (docked -> popped) is
// pinned by dock_test.go without anyone having to own the physical cursor.
// That matters: TrackMouseEvent fires an immediate WM_MOUSELEAVE for a
// synthesised hover, so the live test proves the WIRING only, and this function
// is where the BEHAVIOUR is proven.
func dockRampFrame(p, target float32, dt time.Duration) (next float32, moved bool) {
	if dt <= 0 {
		return p, false
	}
	next = clamp01(approach(p, target, float32(dt)/float32(DockAnimMs*time.Millisecond)))
	return next, next != p
}

// DockSquash is the fraction of the orb still drawn along the dock axis at
// ramp level p (1 = free orb, DockOverlapFrac = full tab). The renderer scales
// the orb with it and DockPos places it, so the two always agree.
func DockSquash(p float32) float32 {
	if p <= 0 {
		return 1
	}
	if p >= 1 {
		return DockOverlapFrac // exact endpoint: no float drift on a landed ramp
	}
	return 1 - (1-DockOverlapFrac)*p
}

// DockTangent returns the window origin at which the UNSQUASHED orb is exactly
// tangent to work-area edge e - the free end of the ramp. Only the dock axis is
// meaningful in the result; the cross axis is left for the caller.
func DockTangent(work Rect, e Edge, edgePx, marginPx int32) PosEntry {
	switch e {
	case EdgeLeft:
		return PosEntry{X: work.L - int(marginPx)}
	case EdgeRight:
		return PosEntry{X: work.R - int(edgePx) + int(marginPx)}
	case EdgeTop:
		return PosEntry{Y: work.T - int(marginPx)}
	case EdgeBottom:
		return PosEntry{Y: work.B - int(edgePx) + int(marginPx)}
	}
	return PosEntry{}
}

// DockGap measures the window at (x, y) against the tangent position of edge e
// along that edge's axis, POSITIVE AWAY from the edge: 0 means the orb is
// tangent, negative means it has been pushed past the boundary.
func DockGap(work Rect, x, y int, e Edge, edgePx, marginPx int32) int32 {
	t := DockTangent(work, e, edgePx, marginPx)
	switch e {
	case EdgeLeft:
		return int32(x - t.X)
	case EdgeRight:
		return int32(t.X - x)
	case EdgeTop:
		return int32(y - t.Y)
	case EdgeBottom:
		return int32(t.Y - y)
	}
	return 0
}

// NearestDockEdge returns the work-area edge the window hugs most closely and
// the signed gap to its tangent position. Every edge is measured against the
// WORK area of the monitor the window is on, never the full screen, so the ball
// docks above the taskbar and stays inside the monitor it belongs to
// (SPEC-08/D42#3: no docking under the taskbar or off a secondary screen).
func NearestDockEdge(work Rect, x, y int, edgePx, marginPx int32) (Edge, int32) {
	best, bestGap := EdgeNone, int32(0)
	for _, e := range []Edge{EdgeLeft, EdgeRight, EdgeTop, EdgeBottom} {
		g := DockGap(work, x, y, e, edgePx, marginPx)
		if best == EdgeNone || g < bestGap {
			best, bestGap = e, g
		}
	}
	return best, bestGap
}

// dockProgressOf reads the squeeze off the gap: nothing beyond triggerPx of the
// edge, full squeeze once the orb is tangent or has been pushed past it, linear
// in between - so the ball visibly compresses as it is pressed toward the edge.
func dockProgressOf(gap, triggerPx int32) float32 {
	if triggerPx <= 0 {
		if gap <= 0 {
			return 1
		}
		return 0
	}
	if gap >= triggerPx {
		return 0
	}
	if gap <= 0 {
		return 1
	}
	return 1 - float32(gap)/float32(triggerPx)
}

// DockPos is the position law of the ramp: it returns the window origin that
// keeps the ORB, at its current squash level, tangent to edge e. Because the
// law follows the squash instead of lerping between two endpoints, the ball
// never shows a part of itself the pointer cannot reach, and the tab it becomes
// is the squashed orb tucked against the boundary with its transparent margin
// riding past the edge.
//
// crossX/crossY are the drag's own off-axis coordinates, clamped so the window
// stays inside the work area on that axis.
func DockPos(work Rect, e Edge, crossX, crossY int, edgePx, orbR int32, p float32) PosEntry {
	off := int(float32(orbR) * DockSquash(p))
	half := int(edgePx / 2)
	maxX, maxY := work.R-int(edgePx), work.B-int(edgePx)
	x, y := crossX, crossY

	// Cross axis first: a tab must never sit off the top/bottom (or side) of
	// the monitor it docked on. DockNone clamps both axes, which is exactly
	// what a free window needs to stay reachable.
	if e == EdgeNone || e == EdgeTop || e == EdgeBottom {
		if x > maxX {
			x = maxX
		}
		if x < work.L {
			x = work.L
		}
	}
	if e == EdgeNone || e == EdgeLeft || e == EdgeRight {
		if y > maxY {
			y = maxY
		}
		if y < work.T {
			y = work.T
		}
	}

	// Dock axis last, so the clamp cannot undo it: the window deliberately
	// rides across the work boundary by exactly the amount the orb shrank.
	switch e {
	case EdgeLeft:
		x = work.L + off - half
	case EdgeRight:
		x = work.R - off - half
	case EdgeTop:
		y = work.T + off - half
	case EdgeBottom:
		y = work.B - off - half
	case EdgeNone:
	}
	return PosEntry{X: x, Y: y}
}

// DockFreePos clamps a free (undocked) window into the work area, which is what
// keeps "near an edge but not docked" reachable instead of half off-screen.
func DockFreePos(work Rect, x, y int, edgePx int32) PosEntry {
	return DockPos(work, EdgeNone, x, y, edgePx, 0, 0)
}
