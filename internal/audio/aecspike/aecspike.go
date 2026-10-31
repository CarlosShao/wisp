package aecspike

import "math"

// The measurement math of the spike: reference alignment, residual echo and
// ERLE-class metrics, synthetic clock drift, and the cross-correlation gate
// that separates "near looks like our reference" (self echo) from "near is a
// third-party source" (external audio). Everything is deterministic pure
// DSP over mono float32 samples, so the unit tests can pin every constant.

// envDecim is the envelope decimation for the coarse correlation search:
// 16 samples = 1ms at the 16k seam rate. The fine pass then refines within
// one envelope tap, so the result keeps sample resolution.
const envDecim = 16

// DefaultMaxLag is the coarse search span: 300ms each way, enough for the
// shared-mode engine's render/capture path latency on one endpoint.
const DefaultMaxLag = 3 * 1600

// AlignByEnvelopeCorrelation returns the lag that best aligns near against
// far under the convention near[i] ~ far[i-lag] (the echo of the reference
// arrives lag samples later), plus the normalized correlation strength at
// that lag (0..1). The coarse pass runs on 1ms envelopes, the fine pass at
// sample resolution, so the cost stays linear in the search span.
func AlignByEnvelopeCorrelation(far, near []float32, maxLag int) (lag int, strength float64) {
	if len(far) == 0 || len(near) == 0 {
		return 0, 0
	}
	if maxLag <= 0 {
		maxLag = DefaultMaxLag
	}
	fe, ne := envelope(far), envelope(near)
	bestEnv, bestScore := 0, -math.MaxFloat64
	for l := -maxLag / envDecim; l <= maxLag/envDecim; l++ {
		if s := envScore(fe, ne, l); s > bestScore {
			bestScore, bestEnv = s, l
		}
	}
	coarse := bestEnv * envDecim
	bestLag, bestFine := 0, -math.MaxFloat64
	for l := coarse - envDecim; l <= coarse+envDecim; l++ {
		if l < -maxLag || l > maxLag {
			continue
		}
		if s := sampleScore(far, near, l); s > bestFine {
			bestFine, bestLag = s, l
		}
	}
	if bestFine < -1 {
		return 0, 0
	}
	return bestLag, bestFine
}

// envelope folds abs() into envDecim-wide max pools: the envelope of a speech
// band signal, 1ms resolution.
func envelope(s []float32) []float32 {
	out := make([]float32, (len(s)+envDecim-1)/envDecim)
	for i := range out {
		hi := (i + 1) * envDecim
		if hi > len(s) {
			hi = len(s)
		}
		m := float32(0)
		for _, v := range s[i*envDecim : hi] {
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

// envScore is the normalized cross-correlation of the two envelopes at
// envelope lag l under near[i] ~ far[i-l].
func envScore(fe, ne []float32, l int) float64 {
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

// sampleScore is envScore at sample resolution on the raw signals.
func sampleScore(far, near []float32, l int) float64 {
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

// SubtractEcho aligns near against far and removes the least-squares scaled
// reference: residual[i] = near[i] - gain*far[i-lag]. With a single gain
// (no adaptive filter) this is the spike's reference-subtraction stand-in for
// the AEC-cleaned signal; the real AEC reading comes from the DMO engine in
// rig_windows.go when the OS component is present.
func SubtractEcho(far, near []float32, lag int) (residual []float32, gain float64, strength float64) {
	var ab, ref float64
	for i := 0; i < len(near); i++ {
		j := i - lag
		if j < 0 || j >= len(far) {
			continue
		}
		ab += float64(far[j]) * float64(near[i])
		ref += float64(far[j]) * float64(far[j])
	}
	gain = 0
	if ref > 0 {
		gain = ab / ref
	}
	residual = make([]float32, len(near))
	for i := range residual {
		j := i - lag
		if j < 0 || j >= len(far) {
			residual[i] = near[i]
			continue
		}
		residual[i] = near[i] - float32(gain)*far[j]
	}
	return residual, gain, strength
}

// ERLEdB is the echo-return-loss-enhancement of the residual against the
// aligned near signal, computed over far-active 32ms frames only (silence
// frames would inflate the metric with division by near-zero energy). The
// farActive reference is the shifted reference the residual was built from.
func ERLEdB(far, near, residual []float32, lag int) float64 {
	const frame = 512 // 32ms at the 16k seam rate
	var nearPow, resPow float64
	threshold := frameRMS(far) * 0.1
	for start := 0; start+frame <= len(near); start += frame {
		fi := start - lag
		if fi < 0 || fi+frame > len(far) {
			continue
		}
		if frameRMS(far[fi:fi+frame]) <= threshold {
			continue
		}
		for k := 0; k < frame; k++ {
			n := float64(near[start+k])
			r := float64(residual[start+k])
			nearPow += n * n
			resPow += r * r
		}
	}
	if resPow <= 0 {
		return math.Inf(1) // perfect cancellation over the active frames
	}
	if nearPow <= 0 {
		return math.Inf(-1)
	}
	return 10 * math.Log10(nearPow/resPow)
}

// frameRMS is the RMS of one frame.
func frameRMS(f []float32) float64 {
	var sum float64
	for _, v := range f {
		sum += float64(v) * float64(v)
	}
	return math.Sqrt(sum / float64(maxInt(1, len(f))))
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ApplyPpmDrift resamples far by 1+ppm/1e6 with linear interpolation: a
// synthetic stand-in for the cross-endpoint clock mismatch (render endpoint
// clock vs capture endpoint clock) that the real Bluetooth-speaker case
// would impose. The sensitivity curve built over several ppms is the
// "thresholds or restrictions" data of ticket 59 AC#4 until the owner runs
// the real cross-endpoint rig on hardware this host does not have.
func ApplyPpmDrift(far []float32, ppm int) []float32 {
	if ppm == 0 || len(far) < 2 {
		return append([]float32(nil), far...)
	}
	ratio := 1 + float64(ppm)/1e6
	n := int(float64(len(far)) / ratio)
	out := make([]float32, n)
	for i := range out {
		x := float64(i) * ratio
		lo := int(x)
		if lo >= len(far)-1 {
			out[i] = far[len(far)-1]
			continue
		}
		frac := x - float64(lo)
		out[i] = far[lo] + float32(frac)*(far[lo+1]-far[lo])
	}
	return out
}

// WindowedCorrelation is the barge-in gate instrument: the normalized
// cross-correlation of near against far over win-wide windows every hop
// samples. Self echo tracks the reference (high correlation); a third-party
// source does not (low correlation). This is the "direction of AEC residual"
// evidence shape ticket 59 names for the external-audio case.
func WindowedCorrelation(far, near []float32, lag, win, hop int) []float64 {
	if win <= 0 || hop <= 0 {
		return nil
	}
	var out []float64
	for start := 0; start+win <= len(near); start += hop {
		var ab, an, bn float64
		for k := 0; k < win; k++ {
			j := start + k - lag
			if j < 0 || j >= len(far) {
				continue
			}
			f := float64(far[j])
			n := float64(near[start+k])
			ab += f * n
			an += f * f
			bn += n * n
		}
		if an == 0 || bn == 0 {
			out = append(out, 0)
			continue
		}
		out = append(out, ab/math.Sqrt(an*bn))
	}
	return out
}
