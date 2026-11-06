package ball

// Ticket-65 Wensity band: pure geometry/motion mapping shared by the
// renderer (wensity_windows.go does the D2D drawing) and the tests. Every
// input/output is a fraction of the orb radius R unless stated; every
// constant is a C21-tabled token (tokens.go). This file is deliberately
// free of Win32/D2D so the deform mapping stays unit-testable on any GOOS
// (same posture as liquid.go).

// wensityBandAxisY returns the band spine height at span position t
// (0 = warm left end, 1 = cold right end) as "R above the orb centre".
// The spine peaks at WensityBandCenterFrac above the centre and the two
// ends differ by WensityBandTiltFrac (the cold end is upturned) - the two
// numbers 65-style-direction §1.1 measured on the reference image.
func wensityBandAxisY(t float32) float32 {
	return WensityBandCenterFrac + WensityBandTiltFrac*(t-0.5)
}

// wensityBandAxisX returns the band spine x offset from the orb centre at
// span position t (fractions of R; negative = toward the warm/left end).
func wensityBandAxisX(t float32) float32 {
	return -WensityBandSpanLeftFrac + (WensityBandSpanLeftFrac+WensityBandSpanRightFrac)*t
}

// wensityBandDeform is the audio-reactive amplitude of the band, in R
// units: the spine shifts and the half-width swells by this much. Full
// envelope maps to the full peak-to-peak travel WensityBandDeformFrac
// (65-style-direction §3.2: "满电平默认 pp≈0.15R"); zero spectrum means the
// static resting frame, byte-for-byte.
func wensityBandDeform(spectrumMid float32) float32 {
	if spectrumMid < 0 {
		spectrumMid = 0
	}
	if spectrumMid > 1 {
		spectrumMid = 1
	}
	return WensityBandDeformFrac * spectrumMid
}

// wensityBandSegmentColor picks the stop colour of span position t: the
// nearest stop along wensityStopTs. No interpolation - the band is a chain
// of soft overlapping spots of PURE token colours, which is also how the
// reference's hue clusters measured (65-style-direction §1.1).
func wensityBandSegmentColor(t float32) Color {
	stops := [6]Color{
		wensityBandTokens.StopRed,
		wensityBandTokens.StopMagenta,
		wensityBandTokens.StopWhite,
		wensityBandTokens.StopBlue,
		wensityBandTokens.StopCyan,
		wensityBandTokens.StopOrange,
	}
	best, bestDist := 0, float32(2)
	for i, ts := range wensityStopTs {
		if d := absf32(t - ts); d < bestDist {
			best, bestDist = i, d
		}
	}
	return stops[best]
}

func absf32(f float32) float32 {
	if f < 0 {
		return -f
	}
	return f
}
