//go:build windows

package ball

// Ticket-65 Wensity band renderer (docs/evidence/s1/65-style-direction.md
// §1.1/§2.1 - the owner-approved quantified direction, A1581).
//
// Baked once per state/size/DPI/look change in setAssets, exactly like every
// other brush family in this package: an animated frame only FILLS with
// existing COM objects (the D32 CPU discipline), never creates one. The
// geometry each frame picks is a pure function of the Visual (wensity.go),
// so the resting frame is reproducible and the audio-reactive deformation
// needs no clock of its own.
//
// Structure kept from drawGlass (ticket 62/68): caustic, halo, shell, state
// core, rim/border. What the band replaces is the three rotating blobs - the
// Wensity look reads through the horizontal rainbow light band, the white
// core, the top white arc highlight, the asymmetric bottom-right bright lip
// and the weak bottom inner glow.

import "math"

// wensityBandSegs is the pre-baked resolution of the band along its span
// (same law as liquidLadder: 16 segments blend into a continuous soft band
// and a frame only picks indexes).
const wensityBandSegs = 16

// setWensity bakes the band brushes for one state visual.
func (r *renderer) setWensity(v Visual) {
	c := r.center()
	R := v.SizePx / 2
	wb := wensityBandTokens

	// Band segments: one soft elliptical spot per span step, sitting on the
	// spine (wensityBandAxisX/Y), coloured by the nearest stop. Half-width
	// from the FWHM token (a 2.355-sigma read of 0.39R with the gradient's
	// soft falloff); horizontal radius wider so neighbours blend seamlessly.
	for i := 0; i < wensityBandSegs; i++ {
		t := (float32(i) + 0.5) / wensityBandSegs
		col := mulA(wensityBandSegmentColor(t), v.Opacity)
		p := d2d1Point2F{
			x: c.x + R*wensityBandAxisX(t),
			y: c.y - R*wensityBandAxisY(t),
		}
		r.band[i] = r.radialBrush([]d2d1GradientStop{
			{0.00, colorF(col)},
			{0.55, colorF(mulA(col, 0.90))},
			{0.85, colorF(mulA(col, 0.32))},
			{1.00, colorF(mulA(col, 0))},
		}, p, d2d1Point2F{0, 0},
			R*WensityBandWidthFrac, R*WensityBandWidthFrac*0.74)
	}

	// The white core: the band's bright centre cluster (measured centroid
	// (−0.08R, −0.06R), spanning −0.84R..+0.63R horizontally).
	core := mulA(wb.StopWhite, v.Opacity)
	r.bandCore = r.radialBrush([]d2d1GradientStop{
		{0.00, colorF(core)},
		{0.58, colorF(mulA(core, 0.72))},
		{1.00, colorF(mulA(core, 0))},
	}, d2d1Point2F{x: c.x - R*0.08, y: c.y - R*WensityBandCenterFrac + R*0.07},
		d2d1Point2F{0, 0}, R*0.55, R*0.40)

	// Weak neutral bottom inner glow (below 0.75R, the "floor" reading).
	glow := mulA(wb.BottomGlow, v.Opacity)
	r.wensityGlow = r.radialBrush([]d2d1GradientStop{
		{0.00, colorF(glow)},
		{1.00, colorF(mulA(wb.BottomGlow, 0))},
	}, d2d1Point2F{x: c.x, y: c.y + R*WensityBottomGlowFrac*0.8},
		d2d1Point2F{0, 0}, R*0.60, R*0.35)
}

// drawWensityBand paints the band layer: bottom glow under the band, the
// colour segments along the spine, then the white core on top. The spine
// offset and the half-width swell carry the audio deformation
// (wensityBandDeform - non-zero only in Listening/Thinking/Speaking, because
// bandDeformState gates the spectrum at the liquid seam), and the stop chain
// drifts along the span with the EXISTING audio-driven LiquidAngle (Sleeping
// holds angle 0, so the resting frame is one deterministic static frame).
func (r *renderer) drawWensityBand(v Visual, c d2d1Point2F, R, dockX, dockY float32) {
	def := R * wensityBandDeform(v.SpectrumMid)
	phase := int(v.LiquidAngle/(2*math.Pi)*float64(wensityBandSegs)) % wensityBandSegs
	if phase < 0 {
		phase += wensityBandSegs
	}

	if r.wensityGlow != nil {
		r.fillEllipse(d2d1Point2F{x: c.x, y: c.y + R*WensityBottomGlowFrac*0.8},
			R*0.60*dockX, R*0.35*dockY, r.wensityGlow)
	}
	for i := 0; i < wensityBandSegs; i++ {
		b := r.band[(i+phase)%wensityBandSegs]
		if b == nil {
			continue
		}
		t := (float32(i) + 0.5) / wensityBandSegs
		ry := R * WensityBandWidthFrac * 0.74 * (1 + wensityBandDeform(v.SpectrumMid))
		r.fillEllipse(d2d1Point2F{
			x: c.x + R*wensityBandAxisX(t)*dockX,
			y: c.y - R*wensityBandAxisY(t)*dockY + def,
		}, R*WensityBandWidthFrac*dockX, ry*dockY, b)
	}
	if r.bandCore != nil {
		r.fillEllipse(d2d1Point2F{
			x: c.x - R*0.08*dockX,
			y: c.y - R*WensityBandCenterFrac*dockY + R*0.07 + def,
		}, R*0.45*dockX, R*0.32*dockY, r.bandCore)
	}
}

// drawWensityEdges paints the wensity-specific edge work AFTER drawGlass's
// shared rim/border pass: the asymmetric bright lip (bottom-right 45° only -
// the other three directions stay the dark rim drawGlass already laid with
// lk.Rim) and the top white arc highlight.
func (r *renderer) drawWensityEdges(v Visual, c d2d1Point2F, R, s, dockX, dockY float32) {
	wb := wensityBandTokens
	r.setSolid(mulA(wb.RimLip, 0.85*v.Opacity))
	r.arcEx(c, R*WensityRimLipFrac*dockX, R*WensityRimLipFrac*dockY,
		15, 75, GlassLipPx*s)
	r.setSolid(mulA(wb.SpecArc, v.Opacity))
	span := float32(WensitySpecArcSpanDeg)
	r.arcEx(c, R*WensitySpecArcRadiusFrac*dockX, R*WensitySpecArcRadiusFrac*dockY,
		-90-span/2, -90+span/2, WensitySpecArcThickFrac*R)
}
