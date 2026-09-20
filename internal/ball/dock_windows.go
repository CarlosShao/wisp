//go:build windows

package ball

// Window-side edge dock (ticket 62 item 2). The geometry lives in dock.go; this
// file is the Win32 glue and it owns exactly one thing the D32 gate cares
// about: it never calls SetTimer. Every frame of the dock ramp is carried by a
// mouse message or a window-position event, because those are the only sources
// that exist while the orb is doing nothing - and "doing nothing" is the state
// a docked ball spends almost all of its life in (Sleeping), where the frozen
// contract allows zero timers and ~zero CPU.
//
// Message map:
//
//	WM_MOUSEMOVE (dragging)   the push toward the edge IS the animation
//	WM_LBUTTONUP (drag ended)  commit: dock to the edge, or land free and clamped
//	WM_MOUSEMOVE (hover)      pop the tab back out (owner: 悬停…弹回完整球)
//	WM_MOUSELEAVE             close the tab again
//	WM_LBUTTONUP (a click)    pop out and stay out until the next drag
//	WM_DISPLAYCHANGE / WM_DPICHANGED  re-evaluate against the new work area
//
// The ramp's position law (dock.DockPos) keeps the squashed orb tangent to the
// work boundary, so what the user sees is always what the user can hit.

import (
	"time"
	"unsafe"
)

// dockMotion is the UI-thread dock record (no locking, like liquidMotion).
type dockMotion struct {
	edge   Edge    // the edge being docked to (EdgeNone while free)
	p      float32 // ramp level 0..1 (0 = free orb, 1 = tab)
	target float32 // the level the ramp is heading to
	owns   bool    // the dock owns the window position (i.e. it is docked)
	hover  bool    // the pointer is on the tab, pop-out in effect
	track  bool    // a TrackMouseEvent leave notification is outstanding
	supped bool    // the user clicked the tab out: no re-dock until the next drag
}

// dockGeometry gathers what the dock math needs, in physical px at this
// window's own DPI (per-monitor V2: the destination monitor's work area, never
// the primary's, and never the full screen - D42#3).
func (b *Ball) dockGeometry() (work Rect, edgePx, marginPx, orbR, triggerPx int32, x, y int, ok bool) {
	if b.hwnd == 0 || b.rend == nil {
		return work, 0, 0, 0, 0, 0, 0, false
	}
	m, ok := monitorRectsForWindow(b.hwnd)
	if !ok {
		return work, 0, 0, 0, 0, 0, 0, false
	}
	var wr rect
	if r, _, _ := pGetWindowRect.Call(uintptr(b.hwnd), unsafePtr(&wr)); r == 0 {
		return work, 0, 0, 0, 0, 0, 0, false
	}
	scale := float32(getDpiForWindow(b.hwnd)) / 96
	if scale < 1 {
		scale = 1
	}
	edgePx = wr.width()
	marginPx = int32(RingMarginPx*scale + 0.5)
	orbR = edgePx/2 - marginPx
	if orbR < 1 {
		return work, 0, 0, 0, 0, 0, 0, false
	}
	triggerPx = int32(float32(DockTriggerPx)*scale + 0.5)
	return m.Work, edgePx, marginPx, orbR, triggerPx, int(wr.l), int(wr.t), true
}

// dockFrameDt returns the real time since the previous DOCK frame when a new
// frame is allowed (>= MinFrameMs, the 30fps cap) and moves the anchor then; it
// returns 0 without consuming anything while the budget is still open, so a
// 125Hz mouse cannot turn a drag into 125 full D2D frames. This is a
// sample-interval measurement for a ramp step, not timeout logic.
func (b *Ball) dockFrameDt() time.Duration {
	dt := time.Since(b.lastDock)
	if dt < MinFrameMs*time.Millisecond {
		return 0
	}
	b.lastDock = time.Now()
	if dt > 250*time.Millisecond {
		dt = 250 * time.Millisecond // a long gap must not fling the tab around
	}
	return dt
}

// dockProximity is the drag path: the window stays where the cursor put it and
// the orb squeezes by how far it was pushed. No window move, no timer.
func (b *Ball) dockProximity() {
	if !prototypeVisuals {
		return
	}
	work, edgePx, marginPx, _, triggerPx, x, y, ok := b.dockGeometry()
	if !ok {
		return
	}
	e, gap := NearestDockEdge(work, x, y, edgePx, marginPx)
	p := dockProgressOf(gap, triggerPx)
	if e == b.dock.edge && !quantChanged(b.dock.p, p) {
		return // nothing a pixel can act on
	}
	b.dock.edge, b.dock.p, b.dock.target = e, p, p
	if b.dockFrameDt() == 0 {
		return
	}
	b.renderFrame()
}

// dockCommit runs the drag-end (and startup / topology-change) evaluation: the
// ball either becomes a tab on the nearest edge or lands fully inside the work
// area. It moves the window directly, which is the "window position event" half
// of the ticket's rule, and it repaints once.
func (b *Ball) dockCommit() {
	if !prototypeVisuals {
		return
	}
	work, edgePx, marginPx, orbR, triggerPx, x, y, ok := b.dockGeometry()
	if !ok {
		return
	}
	e, gap := NearestDockEdge(work, x, y, edgePx, marginPx)
	b.dock.hover, b.dock.track = false, false
	b.lastDock = time.Now()

	if e == EdgeNone || gap > triggerPx || b.dock.supped {
		// Free: make sure the orb is entirely on-screen and the squash is gone.
		b.dock.edge, b.dock.p, b.dock.target, b.dock.owns = EdgeNone, 0, 0, false
		b.dockMoveTo(DockFreePos(work, x, y, edgePx))
		return
	}
	b.dock.edge, b.dock.p, b.dock.target, b.dock.owns = e, 1, 1, true
	b.dockMoveTo(DockPos(work, e, x, y, edgePx, orbR, b.dock.p))
}

// dockMoveTo places the window (if that is a move at all) and repaints.
func (b *Ball) dockMoveTo(pos PosEntry) {
	if b.hwnd == 0 {
		return
	}
	var wr rect
	if r, _, _ := pGetWindowRect.Call(uintptr(b.hwnd), unsafePtr(&wr)); r == 0 {
		return
	}
	if int(wr.l) == pos.X && int(wr.t) == pos.Y {
		return // nothing to present either: the pixels did not change
	}
	moveWindow(uintptr(b.hwnd), int32(pos.X), int32(pos.Y), wr.width(), wr.height())
	b.renderFrame()
}

// dockStep advances the ramp toward its target and follows the position law.
func (b *Ball) dockStep(dt time.Duration) {
	if dt <= 0 || b.dock.edge == EdgeNone {
		return
	}
	before := b.dock.p
	b.dock.p = clamp01(approach(b.dock.p, b.dock.target, float32(dt)/float32(DockAnimMs*time.Millisecond)))
	if before == b.dock.p {
		return
	}
	if b.dock.owns {
		if work, edgePx, _, orbR, _, x, y, ok := b.dockGeometry(); ok {
			b.dockMoveTo(DockPos(work, b.dock.edge, x, y, edgePx, orbR, b.dock.p))
		}
	}
	b.renderFrame()
}

// dockHoverMove is a WM_MOUSEMOVE that is not a drag: the pointer is on the tab
// (WM_NCHITTEST returned HTCLIENT to let it through), so the ball pops back out
// - one step per message, which is exactly as smooth as the pointer that asked
// for it.
func (b *Ball) dockHoverMove() {
	if !prototypeVisuals || b.dock.edge == EdgeNone {
		return
	}
	if !b.dock.hover {
		b.dock.hover = true
		b.dockAskLeave()
	}
	b.dock.target = 0
	b.dockStep(b.dockFrameDt())
}

// dockAskLeave requests the one-shot WM_MOUSELEAVE (the event-driven half of
// "no resident timer": the OS tells us when the pointer is gone, we never look).
func (b *Ball) dockAskLeave() {
	if b.dock.track || b.hwnd == 0 {
		return
	}
	var tmi trackMouseInfo
	tmi.size = uint32(unsafe.Sizeof(tmi))
	tmi.flags = tmMouseLeave
	tmi.hwndTrack = b.hwnd
	if r, _, _ := pTrackMouseEvent.Call(unsafePtr(&tmi)); r != 0 {
		b.dock.track = true
	}
}

// dockLeave is WM_MOUSELEAVE: the pointer is gone, so the tab closes. One
// message and no frame source left, so the ramp takes the step it is owed and
// then lands (see the note in dock.go).
func (b *Ball) dockLeave() {
	if !prototypeVisuals {
		return
	}
	b.dock.hover, b.dock.track = false, false
	if b.dock.edge == EdgeNone || b.dock.supped {
		return
	}
	b.dock.target = 1
	b.dockStep(b.dockFrameDt())
	if b.dock.p == 1 {
		return
	}
	b.dock.p = 1
	if work, edgePx, _, orbR, _, x, y, ok := b.dockGeometry(); ok {
		b.dockMoveTo(DockPos(work, b.dock.edge, x, y, edgePx, orbR, 1))
	}
	b.renderFrame()
}

// dockClickOut consumes a click on a docked tab: the ball comes fully out and
// stays out until the next drag (SPEC-08 §2's click is a summon gesture; on a
// tab the click means "give me the ball back", so it is not forwarded).
func (b *Ball) dockClickOut() bool {
	if !prototypeVisuals || b.dock.edge == EdgeNone {
		return false
	}
	b.dock.supped, b.dock.hover, b.dock.target = true, false, 0
	b.dock.p, b.dock.owns, b.dock.edge = 0, false, EdgeNone
	if work, edgePx, _, _, _, x, y, ok := b.dockGeometry(); ok {
		b.dockMoveTo(DockFreePos(work, x, y, edgePx))
	}
	b.renderFrame()
	return true
}

// stampDock puts the live dock on the draw-time copy of the visual. An undocked
// ball (and every frozen-mode frame) leaves the visual exactly as mapped, which
// is what keeps the committed Sleeping frame byte-identical.
func (b *Ball) stampDock(v *Visual) {
	if !prototypeVisuals || !b.dock.owns || b.dock.edge == EdgeNone || b.dock.p <= 0 {
		return
	}
	v.Dock, v.DockProgress = b.dock.edge, b.dock.p
}

// DebugDock is the evidence seam used by cmd/balldebug -dock: it puts the window
// at the work-area edge the way a drag ending there would (tangent, fully
// visible) and runs the very same commit path the mouse uses. No timer is
// started here or by anything below.
func (b *Ball) DebugDock(e Edge) bool {
	if e == EdgeNone {
		return false
	}
	done := make(chan bool, 1)
	b.sta.PostTask(func() {
		work, edgePx, marginPx, _, _, _, _, ok := b.dockGeometry()
		if !ok {
			done <- false
			return
		}
		t := DockTangent(work, e, edgePx, marginPx)
		// Park the cross axis mid-work-area: the seam must ask for THIS edge,
		// not inherit a dock on whichever edge the last test left it hugging.
		x := work.L + (work.R-work.L-int(edgePx))/2
		y := work.T + (work.B-work.T-int(edgePx))/2
		if e == EdgeLeft || e == EdgeRight {
			x = t.X
		} else {
			y = t.Y
		}
		b.dock.supped = false // the seam asks for the behavior, not for its exception
		b.dockMoveTo(DockFreePos(work, x, y, edgePx))
		b.dockCommit()
		done <- b.dock.edge == e && b.dock.owns && b.dock.p == 1
	})
	return <-done
}

// Docked reports the live dock (edge, ramp level, whether the dock owns the
// window position). Safe from any goroutine; used by the evidence harness and
// by the live tests.
func (b *Ball) Docked() (Edge, float32, bool) {
	type dockRead struct {
		e    Edge
		p    float32
		owns bool
	}
	done := make(chan dockRead, 1)
	b.sta.PostTask(func() { done <- dockRead{b.dock.edge, b.dock.p, b.dock.owns} })
	r := <-done
	return r.e, r.p, r.owns
}

// EdgeByName resolves a -dock flag value (left|right|top|bottom).
func EdgeByName(s string) (Edge, bool) {
	switch {
	case sameName(s, "left"):
		return EdgeLeft, true
	case sameName(s, "right"):
		return EdgeRight, true
	case sameName(s, "top"):
		return EdgeTop, true
	case sameName(s, "bottom"):
		return EdgeBottom, true
	case sameName(s, "none"):
		return EdgeNone, true
	}
	return EdgeNone, false
}
