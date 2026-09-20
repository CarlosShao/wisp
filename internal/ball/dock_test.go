package ball

import (
	"testing"
	"time"
)

// Edge-dock geometry (ticket 62 item 2, owner: "就像迅雷那种悬浮球一样"). These
// run on any OS: the numbers are physical px, and the window code feeds them the
// real work area of the monitor the ball is on.

const (
	dockTestOrb    = int32(56) // configured orb at 96 DPI
	dockTestMargin = int32(RingMarginPx)
)

func dockTestEdge() int32 { return int32(WindowEdgePx(int(dockTestOrb), 96)) } // orb + 2*margin
func dockTestR() int32    { return dockTestEdge()/2 - dockTestMargin }

// workRect is a primary monitor with a 48px taskbar at the bottom: the dock must
// measure against the WORK area, never the screen, so a tab can never end up
// under the taskbar (SPEC-08/D42#3).
func dockWork() Rect {
	return Rect{L: 0, T: 0, R: 3440, B: 1392}
}

func TestDockGeometryShape(t *testing.T) {
	if e := dockTestEdge(); e != dockTestOrb+2*dockTestMargin {
		t.Fatalf("window edge %d must be orb+2*margin", e)
	}
	if r := dockTestR(); r != dockTestOrb/2 {
		t.Fatalf("orb radius from the window edge: got %d, want %d", r, dockTestOrb/2)
	}
}

func TestDockTangentMeansGapZero(t *testing.T) {
	work, e, m := dockWork(), dockTestEdge(), dockTestMargin
	for _, edge := range []Edge{EdgeLeft, EdgeRight, EdgeTop, EdgeBottom} {
		tn := DockTangent(work, edge, e, m)
		if got := DockGap(work, tn.X, tn.Y, edge, e, m); got != 0 {
			t.Errorf("%d: gap at the tangent position must be 0, got %d", edge, got)
		}
		// The unsquashed orb must be exactly tangent: its drawn boundary lands
		// on the work edge at ramp level 0, which is what DockPos keeps true.
		pos := DockPos(work, edge, tn.X, tn.Y, e, dockTestR(), 0)
		if pos.X != tn.X || pos.Y != tn.Y {
			t.Errorf("%d: DockPos(p=0) = %+v, want the tangent %+v", edge, pos, tn)
		}
		// And away from the edge the gap is positive in the away direction.
		away := DockGap(work, tn.X+40, tn.Y+40, edge, e, m)
		if edge == EdgeLeft || edge == EdgeTop {
			if away <= 0 {
				t.Errorf("%d: moving away must read a positive gap, got %d", edge, away)
			}
		}
	}
}

// TestDockPosKeepsTheOrbTangentAtEveryRampLevel is the load-bearing law: while
// the tab slides in and the orb shrinks, the drawn edge of the orb stays on the
// work boundary, so nothing the user can see is unreachable.
func TestDockPosKeepsTheOrbTangentAtEveryRampLevel(t *testing.T) {
	work, e, R := dockWork(), dockTestEdge(), dockTestR()
	for _, edge := range []Edge{EdgeLeft, EdgeRight, EdgeTop, EdgeBottom} {
		for _, p := range []float32{0, 0.25, 0.5, 0.75, 1} {
			crossX, crossY := 1200, 600
			pos := DockPos(work, edge, crossX, crossY, e, R, p)
			half := int(e / 2)
			drawn := int(float32(R) * DockSquash(p)) // half the orb along the axis
			var got int
			switch edge {
			case EdgeLeft:
				got = pos.X + half - drawn // orb's left boundary
				if want := work.L; got != want {
					t.Errorf("left p=%.2f: orb left %d, want the work edge %d", p, got, want)
				}
			case EdgeRight:
				got = pos.X + half + drawn
				if want := work.R; got != want {
					t.Errorf("right p=%.2f: orb right %d, want %d", p, got, want)
				}
			case EdgeTop:
				got = pos.Y + half - drawn
				if want := work.T; got != want {
					t.Errorf("top p=%.2f: orb top %d, want %d", p, got, want)
				}
			case EdgeBottom:
				got = pos.Y + half + drawn
				if want := work.B; got != want {
					t.Errorf("bottom p=%.2f: orb bottom %d, want %d (the taskbar line)", p, got, want)
				}
			}
		}
	}
}

func TestDockSquashIsMonotonicAndClamped(t *testing.T) {
	if got := DockSquash(0); got != 1 {
		t.Fatalf("a free orb must not be squashed, got %v", got)
	}
	if got := DockSquash(1); got != DockOverlapFrac {
		t.Fatalf("full dock must leave DockOverlapFrac, got %v", got)
	}
	prev := DockSquash(-1)
	for _, p := range []float32{0, 0.2, 0.4, 0.6, 0.8, 1, 2} {
		got := DockSquash(p)
		if got > prev {
			t.Fatalf("the squash must shrink as the ramp grows: p=%v gave %v after %v", p, got, prev)
		}
		prev = got
	}
	if DockSquash(1) <= 0 {
		t.Fatal("a docked orb must still be drawn")
	}
}

func TestDockProgressReadsThePush(t *testing.T) {
	const trigger = DockTriggerPx
	if got := dockProgressOf(trigger, trigger); got != 0 {
		t.Fatalf("at the trigger distance the orb is still round, got %v", got)
	}
	if got := dockProgressOf(trigger+40, trigger); got != 0 {
		t.Fatalf("far from the edge nothing squeezes, got %v", got)
	}
	if got := dockProgressOf(0, trigger); got != 1 {
		t.Fatalf("tangent must be fully squeezed, got %v", got)
	}
	if got := dockProgressOf(-30, trigger); got != 1 {
		t.Fatalf("pushed past the edge must clamp at 1, got %v", got)
	}
	prev := dockProgressOf(trigger+500, trigger)
	for _, gap := range []int32{trigger, trigger / 2, trigger / 4, 4, 1, 0, -10} {
		got := dockProgressOf(gap, trigger)
		if got < prev {
			t.Fatalf("the squeeze must grow as the ball is pushed in: gap %d gave %v after %v", gap, got, prev)
		}
		prev = got
	}
	if prev != 1 {
		t.Fatalf("tangent and beyond must clamp at 1, got %v", prev)
	}
}

func TestNearestDockEdgePicksTheEdgeYouPushed(t *testing.T) {
	work, e, m := dockWork(), dockTestEdge(), dockTestMargin
	half := int(e / 2)
	cy := 600
	cases := []struct {
		name string
		x, y int
		want Edge
	}{
		{"left", work.L - int(m), cy, EdgeLeft},
		{"right", work.R - int(e) + int(m), cy, EdgeRight},
		{"top", 1200, work.T - int(m), EdgeTop},
		{"bottom", 1200, work.B - int(e) + int(m), EdgeBottom},
	}
	for _, c := range cases {
		got, gap := NearestDockEdge(work, c.x, c.y, e, m)
		if got != c.want || gap != 0 {
			t.Errorf("%s: NearestDockEdge = (%d, gap %d), want edge %d at gap 0",
				c.name, got, gap, c.want)
		}
	}
	// Mid-screen: the nearest edge is still reported, but the caller only
	// commits a dock inside the trigger distance, which 1200px is not.
	_, gap := NearestDockEdge(work, 1200, 600, e, m)
	if gap <= triggerForTest() {
		t.Fatalf("a mid-screen ball must not read as dockable, gap %d", gap)
	}
	// A corner hands back the edge it is closer to.
	x := work.R - int(e) - 5 // 5px past the right tangent
	y := work.B - int(e) - 2 // 2px past the bottom tangent
	if got, _ := NearestDockEdge(work, x, y, e, m); got != EdgeBottom {
		t.Fatalf("corner: want the closer edge (bottom), got %d (half=%d)", got, half)
	}
}

func triggerForTest() int32 { return DockTriggerPx }

// TestDockedLeavesAHittableStrip is "只留一条可命中区域" as a number: at full
// dock the orb is a sliver, but a sliver wide enough to click.
func TestDockedLeavesAHittableStrip(t *testing.T) {
	work, e, m, R := dockWork(), dockTestEdge(), dockTestMargin, dockTestR()
	free := DockPos(work, EdgeRight, 0, 600, e, R, 0)
	docked := DockPos(work, EdgeRight, 0, 600, e, R, 1)
	if docked.X <= free.X {
		t.Fatalf("docking must slide the window past the edge: free x=%d docked x=%d", free.X, docked.X)
	}
	visibleOrb := int(float32(2*R) * DockSquash(1))
	if visibleOrb < 16 || visibleOrb > int(dockTestOrb)/2+2 {
		t.Fatalf("the docked tab should be a sliver of the orb, got %dpx of %dpx", visibleOrb, dockTestOrb)
	}
	// The strip the pointer can hit is at least the visible orb plus the
	// click-through margin, and it is strictly smaller than the free window.
	strip := work.R - docked.X
	if strip <= visibleOrb || strip > int(e)-int(float32(R)*(1-DockOverlapFrac))+int(m) {
		t.Fatalf("unexpected hit strip width %d for a %dpx window showing %dpx of orb", strip, e, visibleOrb)
	}
}

// TestDockFreePosKeepsTheBallReachable: a drag that ends half off-screen must
// land fully inside the work area when it is not close enough to dock.
func TestDockFreePosKeepsTheBallReachable(t *testing.T) {
	work, e := dockWork(), dockTestEdge()
	hang := DockPos(work, EdgeRight, 0, 600, e, dockTestR(), 0) // tangent: 8px past the edge
	pos := DockFreePos(work, hang.X+900, hang.Y+9000, e)
	if pos.X > work.R-int(e) || pos.Y > work.B-int(e) || pos.X < work.L || pos.Y < work.T {
		t.Fatalf("free clamp left the work area: %+v (max %d,%d)", pos, work.R-int(e), work.B-int(e))
	}
	// A window pushed far past the left/top boundary comes back inside too.
	pos = DockFreePos(work, work.L-500, work.T-500, e)
	if pos.X != work.L || pos.Y != work.T {
		t.Fatalf("left/top clamp: got %+v, want (%d,%d)", pos, work.L, work.T)
	}
}

// TestDockOnSecondaryMonitorWithNegativeOrigin: the dock math is per-monitor, so
// a screen left of the primary docks against ITS own edges, not the primary's.
func TestDockOnSecondaryMonitorWithNegativeOrigin(t *testing.T) {
	work := Rect{L: -1920, T: 0, R: -100, B: 1000}
	e, m, R := dockTestEdge(), dockTestMargin, dockTestR()
	tn := DockTangent(work, EdgeLeft, e, m)
	if got := DockGap(work, tn.X, tn.Y, EdgeLeft, e, m); got != 0 {
		t.Fatalf("negative-origin monitor: gap at tangent %d", got)
	}
	pos := DockPos(work, EdgeLeft, 9999, 9999, e, R, 1)
	if orbLeft := pos.X + int(e/2) - int(float32(R)*DockSquash(1)); orbLeft != work.L {
		t.Fatalf("secondary left edge: orb left %d, want %d", orbLeft, work.L)
	}
	// The cross-axis clamp keeps it inside THIS monitor, not the primary.
	if pos.Y > work.B-int(e) || pos.Y < work.T {
		t.Fatalf("cross clamp left the secondary work area: %+v", pos)
	}
}

// dockEdgeName labels an edge for failure messages (Edge has no String method).
func dockEdgeName(e Edge) string {
	switch e {
	case EdgeLeft:
		return "left"
	case EdgeRight:
		return "right"
	case EdgeTop:
		return "top"
	case EdgeBottom:
		return "bottom"
	}
	return "none"
}

// orbAxisSpan is how much orb the screen shows along the dock axis.
func orbAxisSpan(r Rect, e Edge) int {
	if e == EdgeLeft || e == EdgeRight {
		return r.R - r.L
	}
	return r.B - r.T
}

// awayFromEdge is the window's travel along the dock axis, positive when it moved
// away from the boundary it was docked on.
func awayFromEdge(now, before PosEntry, e Edge) int {
	switch e {
	case EdgeLeft:
		return now.X - before.X
	case EdgeRight:
		return before.X - now.X
	case EdgeTop:
		return now.Y - before.Y
	case EdgeBottom:
		return before.Y - now.Y
	}
	return 0
}

// insideWork reports whether a rect sits wholly inside the monitor's work area,
// taskbar line included.
func insideWork(work, r Rect) bool {
	return r.L >= work.L && r.T >= work.T && r.R <= work.R && r.B <= work.B
}

// orbRect is everything the screen shows of the orb at ramp level p, in physical
// px: the drawn circle (squashed along the dock axis) as placed by DockPos. The
// renderer scales by DockSquash and DockPos places, so this single law is what
// both the pop-out and the retreat are judged against.
func orbRect(pos PosEntry, edgePx, R int32, e Edge, p float32) Rect {
	half := int(edgePx / 2)
	cx, cy := pos.X+half, pos.Y+half
	off := int(float32(R) * DockSquash(p))
	switch e {
	case EdgeLeft, EdgeRight:
		return Rect{L: cx - off, T: cy - int(R), R: cx + off, B: cy + int(R)}
	}
	return Rect{L: cx - int(R), T: cy - off, R: cx + int(R), B: cy + off}
}

// TestDockHoverPopBackWalksTheRampHome is the deterministic half of "悬停 ->
// 弹回完整球", and it never skips. The live test (TestBallLiveEdgeDockHover) can
// only prove the WIRING: a synthesised hover is answered by an immediate
// WM_MOUSELEAVE from TrackMouseEvent, and the OS may refuse the pointer to a
// background process outright. Per this repo's standing rule a behavior
// assertion must not live only in a test that can skip, so this test drives the
// REVERSE ramp - docked tab to full orb, then back - through dockRampFrame, the
// exact step law dockHoverMove and dockLeave call on every mouse message, and
// asserts the tab really walks home and lands on a reachable, fully drawn orb.
func TestDockHoverPopBackWalksTheRampHome(t *testing.T) {
	const frame = 33 * time.Millisecond
	e, m, R := dockTestEdge(), dockTestMargin, dockTestR()
	tabSpan := 2 * int(float32(R)*DockOverlapFrac)

	// Two monitor layouts: the primary, and a screen left of it (negative
	// origin), so the pop-out law is pinned per-monitor and not against (0,0).
	layouts := []struct {
		name  string
		work  Rect
		edges []Edge
	}{
		{"primary", dockWork(), []Edge{EdgeLeft, EdgeRight, EdgeTop, EdgeBottom}},
		{"secondary", Rect{L: -1920, T: 120, R: -120, B: 1120}, []Edge{EdgeLeft, EdgeTop}},
	}

	for _, lay := range layouts {
		work := lay.work
		for _, edge := range lay.edges {
			name := lay.name + "/" + dockEdgeName(edge)
			// Where a drag would leave the window: tangent and un-squashed,
			// then committed to a full dock by the same law.
			tn := DockTangent(work, edge, e, m)
			crossX, crossY := tn.X, tn.Y
			if edge == EdgeLeft || edge == EdgeRight {
				crossY = work.T + (work.B-work.T-int(e))/2
			} else {
				crossX = work.L + (work.R-work.L-int(e))/2
			}
			docked := DockPos(work, edge, crossX, crossY, e, R, 1)

			// 1. The docked end is a tab, not a ball - that is what the pop has
			//    to undo, and it is still wholly on its own monitor.
			tab := orbRect(docked, e, R, edge, 1)
			if got := orbAxisSpan(tab, edge); got != tabSpan {
				t.Fatalf("%s: docked tab shows %dpx, want %dpx", name, got, tabSpan)
			}
			if !insideWork(work, tab) {
				t.Fatalf("%s: a docked tab must stay inside its work area: %+v in %+v", name, tab, work)
			}

			// 2. Hover: the handler orders target=0 and the ramp walks back.
			p, target := float32(1), float32(0)
			prevPos, prevSpan, frames := docked, tabSpan, 0
			for p != target {
				if frames >= 100 {
					t.Fatalf("%s: the pop-out never landed, stuck at p=%v", name, p)
				}
				next, moved := dockRampFrame(p, target, frame)
				if !moved {
					t.Fatalf("%s: the ramp froze at p=%v with target=0", name, p)
				}
				if next > p {
					t.Fatalf("%s: popping out must shrink the ramp level, %v -> %v", name, p, next)
				}
				p = next
				pos := DockPos(work, edge, crossX, crossY, e, R, p)
				span := orbAxisSpan(orbRect(pos, e, R, edge, p), edge)
				if span < prevSpan {
					t.Fatalf("%s: the orb must grow as it pops out: %dpx after %dpx", name, span, prevSpan)
				}
				if awayFromEdge(pos, prevPos, edge) < 0 {
					t.Fatalf("%s: the pop-out moved the window back toward the edge at p=%v", name, p)
				}
				prevPos, prevSpan = pos, span
				frames++
			}
			if want := DockAnimMs / MinFrameMs; frames < want-1 || frames > want+2 {
				t.Fatalf("%s: the pop-out took %d frames of ~%d ms, budget is DockAnimMs=%d",
					name, frames, MinFrameMs, DockAnimMs)
			}

			// 3. Landed: a whole orb, fully on its own monitor, so it is
			//    reachable and there is nothing left to hide behind an edge.
			full := orbRect(prevPos, e, R, edge, p)
			if got := orbAxisSpan(full, edge); got != 2*int(R) {
				t.Fatalf("%s: a popped-out orb must be round again, axis span %dpx of %dpx",
					name, got, 2*int(R))
			}
			if !insideWork(work, full) {
				t.Fatalf("%s: the popped-out orb is partly off-screen: %+v in %+v", name, full, work)
			}
			if got := DockSquash(p); got != 1 {
				t.Fatalf("%s: a landed ramp must leave an unsquashed orb, got %v", name, got)
			}
			// The cross-axis clamp holds at the free end too: this is the
			// position dockClickOut hands back when the user wants the ball.
			free := DockFreePos(work, prevPos.X, prevPos.Y, e)
			win := Rect{L: free.X, T: free.Y, R: free.X + int(e), B: free.Y + int(e)}
			if !insideWork(work, win) {
				t.Fatalf("%s: a clicked-out ball may not hang off its monitor: %+v", name, win)
			}

			// 4. Pointer gone mid-pop: the same law, target=1, walks back to the
			//    tab, and one full-budget frame lands it exactly - the forced
			//    landing dockLeave performs when no message will follow.
			for p = float32(0.5); p != 1; {
				next, _ := dockRampFrame(p, 1, frame)
				if next < p {
					t.Fatalf("%s: retreating must grow the ramp level, %v -> %v", name, p, next)
				}
				if orbAxisSpan(orbRect(DockPos(work, edge, crossX, crossY, e, R, next), e, R, edge, next), edge) >
					orbAxisSpan(orbRect(DockPos(work, edge, crossX, crossY, e, R, p), e, R, edge, p), edge) {
					t.Fatalf("%s: retreating must shrink the orb, %v -> %v", name, p, next)
				}
				p = next
				if p > 1 {
					t.Fatalf("%s: the retreat overshot the docked end: %v", name, p)
				}
			}
			if landed, _ := dockRampFrame(0.5, 1, DockAnimMs*time.Millisecond); landed != 1 {
				t.Fatalf("%s: a full-budget retreat must land the tab, got %v", name, landed)
			}

			// 5. Endpoints hold still and never overshoot: a landed ramp costs no
			//    further frames (no jitter, no CPU) and cannot slide past either
			//    end of the ramp.
			if _, moved := dockRampFrame(0, 0, 250*time.Millisecond); moved {
				t.Fatalf("%s: a popped-out ramp must hold still at 0", name)
			}
			if _, moved := dockRampFrame(1, 1, 250*time.Millisecond); moved {
				t.Fatalf("%s: a docked ramp must hold still at 1", name)
			}
			if over, _ := dockRampFrame(0.1, 0, time.Second); over != 0 {
				t.Fatalf("%s: the ramp must clamp at the free end, got %v", name, over)
			}
			if over, _ := dockRampFrame(0.9, 1, time.Second); over != 1 {
				t.Fatalf("%s: the ramp must clamp at the docked end, got %v", name, over)
			}
		}
	}
}

// TestDockRampTakesItsTokenBudget: the ramp the mouse messages carry must walk
// DockAnimMs of frames, not land in one step - and must land exactly.
func TestDockRampWalksItsBudget(t *testing.T) {
	var p float32
	steps := 0
	for p != 1 && steps < 100 {
		p = clamp01(approach(p, 1, float32(33*time.Millisecond)/float32(DockAnimMs*time.Millisecond)))
		steps++
	}
	if p != 1 {
		t.Fatalf("the ramp must reach its target, stopped at %v", p)
	}
	if want := DockAnimMs / MinFrameMs; steps < want-1 || steps > want+2 {
		t.Fatalf("the ramp took %d frames of ~%d ms; DockAnimMs=%d MinFrameMs=%d",
			steps, MinFrameMs, DockAnimMs, MinFrameMs)
	}
}
