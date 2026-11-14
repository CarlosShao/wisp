package audio

// Audit probe for ticket 449 (AC#1-AC#5). NOT production code: this file lives
// in .scratch and is only ever injected into package audio through
// `go test -overlay`, so the tree under audit is untouched.
//
// Every reading below is taken from the real exported surface
// (LevelOfSamples / FrameLevel / EncodeFrame) of whatever level.go the overlay
// maps in: HEAD's (mean removed) or f9b6d60f7^'s (DC counted). The same file is
// run against both, which is what makes "before" and "after" comparable rather
// than two different probes.

import (
	"math"
	"testing"
)

func ac449ConstFrame(v int16) []int16 {
	s := make([]int16, FrameSamples)
	for i := range s {
		s[i] = v
	}
	return s
}

// AC#1: a constant frame reads exactly 0, through both forms, and the zero is
// POSITIVE zero (Signbit false, all bits clear).
func Test449AC1ConstantFrameIsPositiveZero(t *testing.T) {
	for _, v := range []int16{0, 1, -1, 12345, 32767, -32768} {
		s := ac449ConstFrame(v)
		l := LevelOfSamples(s)
		fl, err := FrameLevel(EncodeFrame(s))
		if err != nil {
			t.Fatalf("AC#1 constant %d: FrameLevel error %v", v, err)
		}
		t.Logf("AC#1 const %6d: LevelOfSamples=%.17g bits=%016x ; FrameLevel=%.9g f32bits=%08x signbit(f64)=%v signbit(f32)=%v",
			v, l, math.Float64bits(l), fl, math.Float32bits(fl), math.Signbit(l), math.Signbit(float64(fl)))
		if l != 0 || math.Float64bits(l) != 0 {
			t.Errorf("AC#1 constant %d: LevelOfSamples = %v (bits %016x), want exactly +0", v, l, math.Float64bits(l))
		}
		if fl != 0 || math.Float32bits(fl) != 0 {
			t.Errorf("AC#1 constant %d: FrameLevel = %v (bits %08x), want exactly +0 (float32)", v, fl, math.Float32bits(fl))
		}
		if math.Signbit(float64(fl)) {
			t.Errorf("AC#1 constant %d: FrameLevel is negative zero (%08x)", v, math.Float32bits(fl))
		}
	}
	// The empty window, for completeness (outside the six shapes the ticket names).
	if l := LevelOfSamples(nil); l != 0 || math.Float64bits(l) != 0 {
		t.Errorf("AC#1 empty window = %v (bits %016x), want exactly +0", l, math.Float64bits(l))
	}
}

// AC#2: adding a constant to every sample does not move the level by one bit -
// Float64bits and Float32bits, not ==.
func Test449AC2TranslationInvarianceBitExact(t *testing.T) {
	base := wholeCycleSine(FrameSamples, 3000, 5)
	want := LevelOfSamples(base)
	t.Logf("AC#2 base (sine 3000, 5 cycles): level=%.17g bits=%016x", want, math.Float64bits(want))
	if want <= 0 {
		t.Fatalf("AC#2 base window reads %v: the tooth would be vacuous", want)
	}
	frame := EncodeFrame(base)
	wantF, err := FrameLevel(frame)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("AC#2 base FrameLevel=%v bits=%08x", wantF, math.Float32bits(wantF))

	shifts := []int16{-12345, -3000, -1, 1, 3000, 12345, 20000}
	for _, c := range shifts {
		shifted := make([]int16, len(base))
		for i, v := range base {
			shifted[i] = v + c
		}
		got := LevelOfSamples(shifted)
		f64same := math.Float64bits(got) == math.Float64bits(want)
		gotF, err := FrameLevel(EncodeFrame(shifted))
		if err != nil {
			t.Fatal(err)
		}
		f32same := math.Float32bits(gotF) == math.Float32bits(wantF)
		t.Logf("AC#2 shift %+6d: level=%.17g bits=%016x f64bits-equal=%v ; FrameLevel=%.9g f32bits=%08x f32bits-equal=%v",
			c, got, math.Float64bits(got), f64same, gotF, math.Float32bits(gotF), f32same)
		if !f64same {
			t.Errorf("AC#2 shift %+d moved Float64bits from %016x to %016x (%v vs %v)",
				c, math.Float64bits(want), math.Float64bits(got), want, got)
		}
		if !f32same {
			t.Errorf("AC#2 shift %+d moved Float32bits from %08x to %08x", c, math.Float32bits(wantF), math.Float32bits(gotF))
		}
	}
}

// AC#3: DC 12345 + whole-cycle sine 6000 reads A/sqrt(2)/FullScale.
func Test449AC3DCPlusSineClosedForm(t *testing.T) {
	const (
		dc   = 12345
		amp  = 6000
		cycs = 8
	)
	sine := wholeCycleSine(FrameSamples, amp, cycs)
	shifted := make([]int16, len(sine))
	for i, v := range sine {
		shifted[i] = v + dc
	}
	withDC := LevelOfSamples(shifted)
	withoutDC := LevelOfSamples(sine)
	want := float64(amp) / math.Sqrt2 / LevelFullScale
	diff := math.Abs(withDC - want)
	t.Logf("AC#3 DC %d + sine %d: read=%.10f (%.17g) closed=%.10f (%.17g) diff=%.3e (= %.2f x 1.4e-6) tol=%.1e",
		dc, amp, withDC, withDC, want, want, diff, diff/1.4e-6, SineLevelTolerance)
	t.Logf("AC#3 same sine without DC: %.17g bits=%016x ; with DC bits=%016x ; bit-equal=%v",
		withoutDC, math.Float64bits(withoutDC), math.Float64bits(withDC),
		math.Float64bits(withoutDC) == math.Float64bits(withDC))
	t.Logf("AC#3 old (DC-counting) reading to compare against: the ticket says 0.3983676962; this run's value for the pre-change arithmetic is in the paired OLD-code log")
	t.Logf("AC#3 distance from the old reading: |read-0.398367696| = %.6f", math.Abs(withDC-0.398367696))
	if diff > SineLevelTolerance {
		t.Errorf("AC#3 read %.17g, want %.17g within %v (off by %.3e)", withDC, want, SineLevelTolerance, diff)
	}
}

// AC#4: the hottest legal frame is the anti-symmetric square +32767/-32768,
// reads 32767.5/32768 = MaxLevel - LevelLSB/2, and MaxLevel is unreachable.
func Test449AC4Ceiling(t *testing.T) {
	hot := make([]int16, FrameSamples)
	for i := range hot {
		if i%2 == 0 {
			hot[i] = 32767
		} else {
			hot[i] = -32768
		}
	}
	got := LevelOfSamples(hot)
	want := (32767.0 + 32768.0) / 2 / LevelFullScale
	gotF, err := FrameLevel(EncodeFrame(hot))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("AC#4 hottest frame: level=%.17g want=%.17g equal=%v ; FrameLevel=%.9g bits=%08x ; MaxLevel=%v attained=%v",
		got, want, got == want, gotF, math.Float32bits(gotF), MaxLevel, got == MaxLevel)
	t.Logf("AC#4 MaxLevel-LevelLSB/2 = %.17g ; LevelLSB = %.17g ; FullScaleSquareLevel = %.17g",
		MaxLevel-LevelLSB/2, LevelLSB, FullScaleSquareLevel)
	if got != want {
		t.Errorf("AC#4 hottest legal frame = %.17g, want exactly %.17g", got, want)
	}
	if got == MaxLevel || !(got < MaxLevel) {
		t.Errorf("AC#4 hottest legal frame %v attains MaxLevel %v", got, MaxLevel)
	}

	// Independent upper bound: every sample is inside [-32768, 32767] and so is
	// the mean, hence |s-mean| <= 32767.5 for every sample, hence RMS <= 32767.5
	// and level <= 32767.5/32768 < 1.0. No legal frame can reach 1.0.
	bound := 32767.5 / LevelFullScale
	t.Logf("AC#4 analytic bound: |s-mean| <= 32767.5 => level <= %.17g < 1.0 ; gap to 1.0 = %.17g (= LevelLSB/2)",
		bound, MaxLevel-bound)

	// A search over the extreme families plus random frames: nothing reaches the
	// bound's value except the anti-symmetric square, nothing reaches 1.0.
	maxSeen, maxName := got, "anti-symmetric +32767/-32768"
	try := func(name string, s []int16) {
		if l := LevelOfSamples(s); l > maxSeen {
			maxSeen, maxName = l, name
		}
	}
	for _, a := range []int16{0, 1, 2, 100, 3000, 16384, 32766, 32767} {
		sq := make([]int16, FrameSamples)
		for i := range sq {
			if i%2 == 0 {
				sq[i] = a
			} else {
				sq[i] = -a - 1
			}
		}
		try("bipolar square", sq)
	}
	for _, a := range []int16{32767, 32766, 32760} {
		alt := make([]int16, FrameSamples)
		for i := range alt {
			if i%2 == 0 {
				alt[i] = a
			} else {
				alt[i] = -32768
			}
		}
		try("extreme asymmetric square", alt)
	}
	// deterministic pseudo-random frames, fixed seed by construction
	state := uint64(0x449449)
	rnd := func() int16 {
		state ^= state << 13
		state ^= state >> 7
		state ^= state << 17
		return int16(state >> 17)
	}
	for k := 0; k < 200000; k++ {
		s := make([]int16, FrameSamples)
		for i := range s {
			s[i] = rnd()
		}
		try("random frame", s)
	}
	t.Logf("AC#4 search: max level over extreme families + 200000 random frames = %.17g (%s) ; MaxLevel attained=%v",
		maxSeen, maxName, maxSeen == MaxLevel)
	if maxSeen > bound {
		t.Errorf("AC#4 search exceeded the analytic bound: %.17g > %.17g", maxSeen, bound)
	}
	if maxSeen == MaxLevel {
		t.Errorf("AC#4 search attained MaxLevel = %v", MaxLevel)
	}
}

// The three sine errors TestLevelSineAgainstRootTwo rides on, with their bits,
// so "the tolerance was not widened" can be checked as a reading.
func Test449SineToleranceErrors(t *testing.T) {
	rootTwo := 1 / math.Sqrt2
	for _, c := range []struct {
		amp int64
		want float64
	}{
		{32767, rootTwo},
		{16384, rootTwo / 2},
		{8192, rootTwo / 4},
	} {
		l := LevelOfSamples(wholeCycleSine(FrameSamples, c.amp, 8))
		err := math.Abs(l - c.want)
		t.Logf("SINE amp %5d: level=%.17g bits=%016x want=%.17g err=%.17e bits=%016x passes=%.1e=%v",
			c.amp, l, math.Float64bits(l), c.want, err, math.Float64bits(err), SineLevelTolerance, err <= SineLevelTolerance)
	}
}

// The tolerance the three sine pins in level_test.go ride on, printed so the
// audit can show it was not widened (its bytes are also diffed outside).
func Test449ToleranceReading(t *testing.T) {
	t.Logf("AC-note SineLevelTolerance = %.17g (5e-5) ; MinLevel=%.17g MaxLevel=%.17g LevelFullScale=%.17g",
		SineLevelTolerance, MinLevel, MaxLevel, LevelFullScale)
}

// AC#5 arithmetic: why a missing guard passes silently. |+Inf - +Inf| is NaN and
// every ordered comparison with NaN is false, so the ratio assertion below it
// never fires.
func Test449InfInfIsNaN(t *testing.T) {
	diff := math.Abs(math.Inf(1) - math.Inf(1))
	t.Logf("AC#5 Inf-Inf: math.Abs(+Inf - +Inf) = %v isNaN=%v bits=%016x ; (diff > 0.01) = %v",
		diff, math.IsNaN(diff), math.Float64bits(diff), diff > 0.01)
	if !math.IsNaN(diff) || diff > 0.01 {
		t.Fatalf("the silent-pass premise does not hold: diff=%v isNaN=%v", diff, math.IsNaN(diff))
	}
}
