package audio

import (
	"fmt"
	"math"
	"testing"
)

// 241-v1 measurement probe (adjudicator leg). Injected with -overlay only; it
// never lands on disk under internal/audio. It reports readings, it asserts
// nothing about the ticket's own nails.

func TestV1ProbeReadings(t *testing.T) {
	fs := float64(LevelFullScale)

	// Shot 3: the two endpoint claims.
	hot := make([]int16, FrameSamples)
	for i := range hot {
		hot[i] = -32768
	}
	lHot := LevelOfSamples(hot)
	fHot, _ := FrameLevel(EncodeFrame(hot))

	sq := func(amp int16) []int16 {
		out := make([]int16, FrameSamples)
		for i := range out {
			if i%2 == 0 {
				out[i] = amp
			} else {
				out[i] = -amp
			}
		}
		return out
	}
	l32767 := LevelOfSamples(sq(32767))
	f32767, _ := FrameLevel(EncodeFrame(sq(32767)))
	l1 := LevelOfSamples(sq(1))

	t.Logf("V1END all-minus-32768   LevelOfSamples=%v FrameLevel(float32)=%v", lHot, fHot)
	t.Logf("V1END plus-minus-32767  LevelOfSamples=%v FrameLevel(float32)=%v FullScaleSquareLevel=%v", l32767, f32767, FullScaleSquareLevel)
	t.Logf("V1END plus-minus-1      LevelOfSamples=%v LevelLSB=%v", l1, LevelLSB)

	// Shot 2: a frame carrying exactly one +-1 sample.
	for _, v := range []int16{1, -1} {
		one := make([]int16, FrameSamples)
		one[7] = v
		l := LevelOfSamples(one)
		f, err := FrameLevel(EncodeFrame(one))
		t.Logf("V1ZERO single-sample=%d LevelOfSamples=%v FrameLevel=%v err=%v exactZero=%v belowSilenceGate=%v",
			v, l, f, err, l == 0.0, l < 0.06)
	}
	// The true floor over all non-zero frames: one LSB in an otherwise zero frame.
	t.Logf("V1ZERO theoretical smallest nonzero = %v (LevelLSB/sqrt(512) = %v)", math.Sqrt(1.0/float64(FrameSamples))/fs, LevelLSB/math.Sqrt(float64(FrameSamples)))

	// Is a silent frame ever negative zero?
	silent := make([]int16, FrameSamples)
	lz := LevelOfSamples(silent)
	t.Logf("V1ZERO silent LevelOfSamples bits=%v signbit=%v", math.Float64bits(lz), math.Signbit(lz))
	fz, _ := FrameLevel(make([]byte, FrameBytes))
	t.Logf("V1ZERO silent FrameLevel bits=%v signbit=%v", math.Float32bits(fz), math.Signbit(float64(fz)))

	// Relationship of the declared scale to the ball's SilenceLevelGate = 0.06.
	// What peak amplitude of a whole-cycle sine is needed to clear 0.06?
	for _, amp := range []int64{1000, 2000, 2780, 3000, 4096, 8192, 16384} {
		s := wholeCycleSineProbe(FrameSamples, amp, 8)
		l := LevelOfSamples(s)
		t.Logf("V1GATE sine amp=%d level=%v clears-0.06=%v", amp, l, l > 0.06)
	}
	t.Logf("V1GATE 20*log10(0.06) = %v dBFS; 0.06 as RMS maps to peak int16 = %v",
		20*math.Log10(0.06), 0.06*fs*math.Sqrt2)

	// DC claim in the doc comment: "a frame held at -32768 reads 1.0".
	dc := make([]int16, FrameSamples)
	for i := range dc {
		dc[i] = -32768
	}
	t.Logf("V1DC constant-minus-32768 = %v", LevelOfSamples(dc))
	dcPos := make([]int16, FrameSamples)
	for i := range dcPos {
		dcPos[i] = 20000
	}
	t.Logf("V1DC constant-20000 = %v (pure DC, no AC at all)", LevelOfSamples(dcPos))

	// Full-scale square above the named bound: is MaxLevel reachable by a
	// symmetric (non-DC) signal at all?
	mixed := make([]int16, FrameSamples)
	for i := range mixed {
		if i < FrameSamples/2 {
			mixed[i] = -32768
		} else {
			mixed[i] = 32767
		}
	}
	t.Logf("V1END half-minus32768-half-plus32767 = %v", LevelOfSamples(mixed))

	// Non-whole-cycle sine: is the named tolerance still honoured?
	for _, cyc := range []int{0, 1, 3, 7, 8, 9, 64, 255, 256} {
		s := wholeCycleSineProbe(FrameSamples, 32767, cyc)
		l := LevelOfSamples(s)
		t.Logf("V1SINE cycles=%d level=%v dev-from-1/sqrt2=%v within-5e-5=%v",
			cyc, l, math.Abs(l-1/math.Sqrt2), math.Abs(l-1/math.Sqrt2) <= SineLevelTolerance)
	}

	// Exported EncodeFrame handed more than FrameSamples: documented behaviour?
	over := make([]int16, FrameSamples+1)
	t.Logf("V1API EncodeFrame(len=%d) reached; see panic probe", len(over))
}

func wholeCycleSineProbe(n int, amp int64, cycles int) []int16 {
	out := make([]int16, n)
	for i := range out {
		out[i] = int16(math.Round(float64(amp)*math.Sin(2*math.Pi*float64(cycles)*float64(i)/float64(n))))
	}
	return out
}

// EncodeFrame with an over-length slice: the exported half of the seam codec.
func TestV1ProbeEncodeFrameOverLength(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("V1API EncodeFrame panics on over-length input: %v", r)
		}
	}()
	samples := make([]int16, FrameSamples+1)
	for i := range samples {
		samples[i] = 1
	}
	got := EncodeFrame(samples)
	fmt.Sprint(len(got))
	t.Logf("V1API EncodeFrame(len=%d) returned %d bytes, no panic", len(samples), len(got))
}

// DecodeFrame/EncodeFrame asymmetry: DecodeFrame takes any even length,
// EncodeFrame silently documents none.
func TestV1ProbeCodecAsymmetry(t *testing.T) {
	samples, err := DecodeFrame(make([]byte, 2048))
	t.Logf("V1API DecodeFrame(2048 bytes) -> %d samples err=%v (accepts non-seam length)", len(samples), err)
}
