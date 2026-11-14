package audio

// Ticket 291: the level scale strips DC. Four teeth, all stated as facts that
// can be re-run, never as "looks lower":
//
//  1. a frame held at ONE constant value is pure DC and reads exactly 0, for
//     any offset (the negative of what this scale used to say: all -32768 read
//     1.0 because the offset was counted);
//  2. the scale is TRANSLATION INVARIANT: level(s) == level(s + c), bit for
//     bit, for any constant c that keeps the samples inside int16. This is the
//     sufficient-and-necessary face of "DC does not take part";
//  3. a DC offset plus a whole-cycle sine reads the sine's own closed-form RMS,
//     A/sqrt(2)/LevelFullScale, not the quadrature sum with the offset;
//  4. the closed top of the scale is still respected, and it is now strictly
//     unreachable: the hottest legal frame is the anti-symmetric square
//     +32767/-32768, whose level is the closed form below.
//
// Measured against the pre-291 arithmetic on this tree, all four were red:
// tooth 1 (constant 1 LSB read 3.0517578125e-05), tooth 2 (adding a DC of
// -12345 moved a sine frame from 0.06473924235500596 to 0.38226145742487416),
// tooth 3 (0.3983676962 against 0.1294751186, off by 0.269), tooth 4
// (0.9999847413273546 vs the closed form 0.9999847412109375 - off only in the
// last bits, because the old RMS still counted the -0.5 mean of the
// anti-symmetric square). Tooth 4 is NOT a DC witness: it pins the new ceiling,
// and it must not be cited as one.

import (
	"math"
	"testing"
)

// ac291ConstFrame builds one seam window held at a single sample value.
func ac291ConstFrame(v int16) []int16 {
	s := make([]int16, FrameSamples)
	for i := range s {
		s[i] = v
	}
	return s
}

// TestLevelDCFrameReadsZero: a constant frame carries no level, exactly zero,
// through both the sample form and the seam form.
func TestLevelDCFrameReadsZero(t *testing.T) {
	for _, v := range []int16{0, 1, -1, 12345, 32767, -32768} {
		s := ac291ConstFrame(v)
		if l := LevelOfSamples(s); l != MinLevel {
			t.Fatalf("constant %d frame = %v, want exactly %v: a frame that does not change carries no level", v, l, MinLevel)
		}
		got, err := FrameLevel(EncodeFrame(s))
		if err != nil {
			t.Fatal(err)
		}
		if got != float32(MinLevel) {
			t.Fatalf("FrameLevel(constant %d) = %v, want exactly %v", v, got, MinLevel)
		}
		if math.Signbit(float64(got)) {
			t.Fatalf("constant %d frame level is negative zero (%08x)", v, math.Float32bits(got))
		}
	}
	// The ticket's own fixture value, stated as a reading: the old scale read
	// 12345/32768 = 0.376740 here.
	if l := LevelOfSamples(ac291ConstFrame(12345)); l != 0 {
		t.Fatalf("the DC fixture (12345) reads %v, want exactly 0", l)
	}
}

// TestLevelTranslationInvariant: shifting every sample by a constant may not
// move the level by one bit. Bit equality, not "close enough".
func TestLevelTranslationInvariant(t *testing.T) {
	base := wholeCycleSine(FrameSamples, 3000, 5)
	want := LevelOfSamples(base)
	if want <= 0 {
		t.Fatalf("the base window is silent (%v), the tooth would be vacuous", want)
	}
	for _, c := range []int16{-12345, -3000, -1, 1, 3000, 12345, 20000} {
		shifted := make([]int16, len(base))
		for i, v := range base {
			shifted[i] = v + c
		}
		got := LevelOfSamples(shifted)
		if math.Float64bits(got) != math.Float64bits(want) {
			t.Fatalf("adding a constant DC of %d moved the level from %v to %v: the scale still counts DC", c, want, got)
		}
	}
	// Same fact through the seam form, float32 and all.
	frame := EncodeFrame(base)
	shiftedFrame := EncodeFrame(func() []int16 {
		out := make([]int16, len(base))
		for i, v := range base {
			out[i] = v + 12345
		}
		return out
	}())
	a, err := FrameLevel(frame)
	if err != nil {
		t.Fatal(err)
	}
	b, err := FrameLevel(shiftedFrame)
	if err != nil {
		t.Fatal(err)
	}
	if math.Float32bits(a) != math.Float32bits(b) {
		t.Fatalf("FrameLevel moved under a constant offset: %v vs %v", a, b)
	}
}

// TestLevelDCPlusSineReadsClosedFormRMS is the closed-form face: a window that
// holds a constant DC plus a whole-cycle sine reads the sine's own RMS,
// A/sqrt(2)/LevelFullScale - not the quadrature sum of the two. Measured before
// ticket 291: 0.3983676962 (DC and sine added in quadrature, red by 0.269).
// Measured after: 0.1294765192 against the closed form 0.1294751186, i.e. off
// by 1.4e-6 from int16 rounding of the sine samples, well inside the package's
// named SineLevelTolerance (5e-5).
func TestLevelDCPlusSineReadsClosedFormRMS(t *testing.T) {
	const (
		dc   = 12345
		amp  = 6000
		cycs = 8
	)
	base := wholeCycleSine(FrameSamples, amp, cycs)
	shifted := make([]int16, len(base))
	for i, v := range base {
		shifted[i] = v + dc
	}
	want := float64(amp) / math.Sqrt2 / LevelFullScale
	withDC := LevelOfSamples(shifted)
	t.Logf("DC %d + sine %d reads %.10f, closed form A/sqrt(2)/FullScale = %.10f (the DC-inclusive reading was 0.3983676962)",
		dc, amp, withDC, want)
	if diff := math.Abs(withDC - want); diff > SineLevelTolerance {
		t.Fatalf("DC %d + sine %d reads %v, want %v (A/sqrt(2)/FullScale) within %v (off by %.3e)",
			dc, amp, withDC, want, SineLevelTolerance, diff)
	}
	// The same sine with the offset taken off is the same number, bit for bit:
	// the closed form above is not a near miss on two nearby values.
	if withoutDC := LevelOfSamples(base); withoutDC != withDC {
		t.Fatalf("the same sine read %v without the DC and %v with it", withoutDC, withDC)
	}
	// And it is nowhere near the old reading, so nobody can call this a rounding
	// coincidence.
	if math.Abs(withDC-0.398367696) < 0.1 {
		t.Fatalf("DC+sine reads %v, which is still the DC-inclusive value", withDC)
	}
}

// TestLevelHottestLegalFrameStaysBelowMaxLevel pins the new ceiling: with the
// window's own mean removed, the hottest legal frame is the anti-symmetric
// square +32767/-32768. Its level is sqrt((32767^2+32768^2)/2)/32768, which is
// 32767.5/32768 = MaxLevel - LevelLSB/2 exactly - so MaxLevel stays the closed
// upper bound of the scale but no legal frame reaches it any more.
func TestLevelHottestLegalFrameStaysBelowMaxLevel(t *testing.T) {
	hot := make([]int16, FrameSamples)
	for i := range hot {
		if i%2 == 0 {
			hot[i] = 32767
		} else {
			hot[i] = -32768
		}
	}
	want := (32767.0 + 32768.0) / 2 / LevelFullScale // = MaxLevel - LevelLSB/2 exactly
	if got := LevelOfSamples(hot); got != want {
		t.Fatalf("hottest legal frame = %v, want exactly %v", got, want)
	}
	if l := LevelOfSamples(hot); !(l > FullScaleSquareLevel && l < MaxLevel) {
		t.Fatalf("hottest legal frame %v is not strictly inside (FullScaleSquareLevel=%v, MaxLevel=%v)", l, FullScaleSquareLevel, MaxLevel)
	}
}
