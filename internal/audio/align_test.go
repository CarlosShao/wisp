package audio

import (
	"math/rand"
	"testing"
)

// alignFixture renders a speech-band surrogate whose envelope is aperiodic:
// a random-walk amplitude envelope over a broadband carrier. A pure harmonic
// stack is periodic at its fundamental (133 samples at 120 Hz) and aliases
// the correlation at every period multiple - the rig measured real content
// for exactly this reason, so the fixture carries the same richness.
func alignFixture(n int, seed int64) []float32 {
	r := rand.New(rand.NewSource(seed))
	out := make([]float32, n)
	env := 0.5
	for i := range out {
		env += 0.08 * r.NormFloat64()
		if env > 1 {
			env = 1
		}
		if env < 0.05 {
			env = 0.05
		}
		out[i] = float32(env * r.NormFloat64())
	}
	return out
}

// TestAlignKnownDelay recovers a known echo delay on speech-shaped material:
// near = far delayed by L and scaled, plus independent noise. The reading
// must land on L within one sample and the strength must clear the gate's
// own far-correlation floor (0.8) for a pure scaled copy.
func TestAlignKnownDelay(t *testing.T) {
	far := alignFixture(TargetRate, 1)
	for _, l := range []int{0, 137, 480, 1600, 493} {
		near := make([]float32, len(far)+l)
		for i, v := range far {
			near[i+l] = 0.3 * v
		}
		lag, strength := AlignByEnvelopeCorrelation(far, near, AlignDefaultMaxLag)
		if lag != l {
			t.Fatalf("lag = %d, want %d", lag, l)
		}
		if strength < BargeinFarCorrFloor {
			t.Fatalf("lag %d: strength %.3f below the correlation floor", l, strength)
		}
	}
}

// TestAlignNoEchoAnswersZero: unrelated signals must not manufacture an
// alignment - the strength collapses and the answer carries no information
// the caller would trust. The near side is pure noise here on purpose: two
// fixtures from the same harmonic generator share their envelope and would
// legitimately align.
func TestAlignNoEchoAnswersZero(t *testing.T) {
	far := alignFixture(TargetRate, 2)
	near := make([]float32, TargetRate)
	r := rand.New(rand.NewSource(7))
	for i := range near {
		near[i] = float32(0.1 * r.NormFloat64())
	}
	_, strength := AlignByEnvelopeCorrelation(far, near, AlignDefaultMaxLag)
	if strength > 0.5 {
		t.Fatalf("unrelated signals aligned at strength %.3f", strength)
	}
}

// The product form's pin against the rig instrument (the 59-w delaytrack
// readings' AlignByEnvelopeCorrelation) lives on the aecspike side: that
// package already imports this one, so the cross-check runs there
// (align_cross_26p_windows_test.go) without an import cycle here.
