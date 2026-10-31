package aecspike

import (
	"math"
	"math/rand"
	"testing"

	"github.com/CarlosShao/wisp/internal/audio"
)

// The ticket-26 integration surface productized this package's delay
// instrument as audio.AlignByEnvelopeCorrelation (this package must not
// compile into the resident binary, so the product form lives in
// internal/audio). This file pins the two to the same answers on the same
// input, so the integration-surface lag readings stay comparable with the
// rig's delaytrack readings (493 samples on the 59-w loopback).

func align26pFixture(n int, seed int64) []float32 {
	r := rand.New(rand.NewSource(seed))
	out := make([]float32, n)
	for i := range out {
		t := float64(i) / float64(audio.TargetRate)
		env := 0.5 + 0.5*math.Sin(2*math.Pi*3.5*t)
		var acc float64
		for h := 1; h <= 8; h++ {
			acc += math.Sin(2*math.Pi*120*float64(h)*t+float64(h)) / float64(h)
		}
		out[i] = float32(0.2*env*acc + 0.02*r.NormFloat64())
	}
	return out
}

func TestAlign26pProductFormMatchesRig(t *testing.T) {
	far := align26pFixture(16000, 4)
	near := make([]float32, 0, len(far)+4000)
	near = append(near, align26pFixture(1600, 5)...)
	for _, v := range far {
		near = append(near, 0.25*v)
	}
	near = append(near, align26pFixture(1600, 6)...)
	lag, strength := audio.AlignByEnvelopeCorrelation(far, near, audio.AlignDefaultMaxLag)
	rlag, rstrength := AlignByEnvelopeCorrelation(far, near, DefaultMaxLag)
	if lag != rlag || strength != rstrength {
		t.Fatalf("product aligner (%d, %.6f) diverged from the rig instrument (%d, %.6f)",
			lag, strength, rlag, rstrength)
	}
}
