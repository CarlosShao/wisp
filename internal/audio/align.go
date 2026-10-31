package audio

import "math"

// AlignByEnvelopeCorrelation is the product form of the rig's delay
// measurement (the aecspike rig's AlignByEnvelopeCorrelation, the instrument
// behind the 59-w delaytrack.json readings). The ticket-26 integration
// surface needs the render-to-capture lag ON the product capture path - the
// AEC3 stream delay hint and the barge-in gate's LagSamples both key on it -
// and the rig package must not compile into the resident binary, so the
// instrument is productized here. The algorithm is the rig's own: a coarse
// pass over 1ms envelope taps (envDecim samples at the 16k seam rate) and a
// fine pass at sample resolution inside the winning tap, so the answer keeps
// sample resolution while the cost stays linear in the search span.
//
// Convention (same as the rig's): near[i] ~ far[i-lag] - a positive lag means
// the echo of far[j] shows up at near[j+lag] (the echo arrives lag samples
// after the render). Returns (lag, strength); strength is the normalized
// correlation at the winning lag (0..1), and a strength below the caller's
// own trust floor means "no measurable echo relationship" - the caller
// decides what that degrades to, this function only measures.

// alignEnvDecim is the envelope decimation: 16 samples = 1ms at the 16k seam
// rate.
const alignEnvDecim = 16

// AlignDefaultMaxLag is the coarse search span: 300ms each way, the span the
// rig validated for one endpoint's render/capture path latency.
const AlignDefaultMaxLag = 3 * 1600

func AlignByEnvelopeCorrelation(far, near []float32, maxLag int) (lag int, strength float64) {
	if len(far) == 0 || len(near) == 0 {
		return 0, 0
	}
	if maxLag <= 0 {
		maxLag = AlignDefaultMaxLag
	}
	fe, ne := alignEnvelope(far), alignEnvelope(near)
	bestEnv, bestScore := 0, -math.MaxFloat64
	for l := -maxLag / alignEnvDecim; l <= maxLag/alignEnvDecim; l++ {
		if s := alignEnvScore(fe, ne, l); s > bestScore {
			bestScore, bestEnv = s, l
		}
	}
	coarse := bestEnv * alignEnvDecim
	bestLag, bestFine := 0, -math.MaxFloat64
	for l := coarse - alignEnvDecim; l <= coarse+alignEnvDecim; l++ {
		if l < -maxLag || l > maxLag {
			continue
		}
		if s := alignSampleScore(far, near, l); s > bestFine {
			bestFine, bestLag = s, l
		}
	}
	if bestFine < -1 {
		return 0, 0
	}
	return bestLag, bestFine
}

// alignEnvelope folds abs() into alignEnvDecim-wide max pools: the envelope
// of a speech-band signal at 1ms resolution.
func alignEnvelope(s []float32) []float32 {
	out := make([]float32, (len(s)+alignEnvDecim-1)/alignEnvDecim)
	for i := range out {
		hi := (i + 1) * alignEnvDecim
		if hi > len(s) {
			hi = len(s)
		}
		m := float32(0)
		for _, v := range s[i*alignEnvDecim : hi] {
			if v < 0 {
				v = -v
			}
			if v > m {
				m = v
			}
		}
		out[i] = m
	}
	return out
}

// alignEnvScore is the normalized cross-correlation of the two envelopes at
// envelope lag l under near[i] ~ far[i-l].
func alignEnvScore(fe, ne []float32, l int) float64 {
	var ab, an, bn float64
	for i := 0; i < len(ne); i++ {
		j := i - l
		if j < 0 || j >= len(fe) {
			continue
		}
		ab += float64(fe[j]) * float64(ne[i])
		an += float64(fe[j]) * float64(fe[j])
		bn += float64(ne[i]) * float64(ne[i])
	}
	if an == 0 || bn == 0 {
		return 0
	}
	return ab / math.Sqrt(an*bn)
}

// alignSampleScore is alignEnvScore at sample resolution on the raw signals.
func alignSampleScore(far, near []float32, l int) float64 {
	var ab, an, bn float64
	for i := 0; i < len(near); i++ {
		j := i - l
		if j < 0 || j >= len(far) {
			continue
		}
		ab += float64(far[j]) * float64(near[i])
		an += float64(far[j]) * float64(far[j])
		bn += float64(near[i]) * float64(near[i])
	}
	if an == 0 || bn == 0 {
		return 0
	}
	return ab / math.Sqrt(an*bn)
}
