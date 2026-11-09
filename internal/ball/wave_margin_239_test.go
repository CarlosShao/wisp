package ball

import "testing"

// Ticket 239 residual 1: the Listening wave train used to rely on a bare
// literal (16) at the renderer's wave draw site, so nothing in the repo
// could catch "someone raises the walk distance or lowers RingMarginPx and
// the window rect starts slicing the waves again" - the exact defect the
// owner signed off on 09-30 (ball read as a framed square). This ruler pins
// the geometry analytically: at every supported orb size the train's outer
// edge (ring offset + full travel + half stroke width) must stay strictly
// inside the window half, and the margin that behaviour rests on must keep
// a floor of 0.25px at 96 DPI (radius-relative, so it scales with size and
// never goes negative at any DPI).
//
// The ring offset (R+5) is still a bare literal at the renderer's ring draw
// sites; it is mirrored here as a named local with a comment so the judge
// reads the same shape the renderer draws. Folding that literal into the
// C21 roster too is deliberately out of this ticket's scope (ticket 239
// residual names only the walk distance).
func Test239WaveTrainStaysInsideTheWindowWithAMarginFloor(t *testing.T) {
	const ringOffsetPx = 5.0 // renderer drawFrame: ring + waves + progress all spawn at R+5*s
	outer := ringOffsetPx + WaveTravelPx + RingStrokePx/2

	if floor := RingMarginPx - outer; floor < 0.25 {
		t.Fatalf("RingMarginPx %.2f leaves only %.3fpx beyond the wave train outer edge %.2f (ring %.1f + travel %.1f + half stroke %.2f): the 09-30 square-ball defect floor is 0.25px; raising the walk distance or lowering the margin re-slices the waves with the window rect",
			RingMarginPx, floor, outer, ringOffsetPx, WaveTravelPx, RingStrokePx/2)
	}

	// WindowEdgePx is radius-relative: half = R + RingMarginPx*scale. The
	// ratio between train outer edge and window half is DPI-independent, so
	// the 96 DPI pass below guards every DPI; the per-size sweep guards the
	// configured range (44..160, ceiling widened by A1636).
	for size := BallSizeMinPx; size <= BallSizeMaxPx; size++ {
		half := float64(WindowEdgePx(size, 96)) / 2
		reach := float64(size)/2 + outer
		if reach >= half {
			t.Fatalf("orb size %d: wave train outer edge %.2f is not strictly inside the window half %.2f - the window rect would clip the waves and the ball reads as a square (ticket 239)",
				size, reach, half)
		}
	}
}
