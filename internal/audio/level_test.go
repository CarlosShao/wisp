package audio

import (
	"bytes"
	"context"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Level scale nails (ticket 241 AC#1). Everything here is integer-input,
// float64-math, so the assertions below are either bit-exact or bounded by the
// one named tolerance in the package (SineLevelTolerance). No wall clock and no
// device is involved: the scale is defined by the samples, not by the machine.

// squareFrame builds one seam frame holding a square wave that alternates
// +amp/-amp (amp <= 32767) - the strongest non-clipping symmetric signal.
func squareFrame(amp int16) []byte {
	samples := make([]int16, FrameSamples)
	for i := range samples {
		if i%2 == 0 {
			samples[i] = amp
		} else {
			samples[i] = -amp
		}
	}
	return EncodeFrame(samples)
}

// wholeCycleSine builds n samples of a sine holding exactly `cycles` whole
// periods inside the window. Whole cycles are what make mean(sin^2) exactly
// 1/2, so the only residual error against A/full-scale/sqrt(2) is int16
// rounding (see the derivation on SineLevelTolerance).
func wholeCycleSine(n int, amp int64, cycles int) []int16 {
	out := make([]int16, n)
	for i := range out {
		// Rounded in integer-safe steps: amp*sin(2*pi*cycles*i/n).
		out[i] = int16(math.Round(float64(amp) * math.Sin(2*math.Pi*float64(cycles)*float64(i)/float64(n))))
	}
	return out
}

// AC#1 first of the three: silence reads exactly 0, not approximately and not
// a fudged floor.
func TestLevelSilentFrameIsExactZero(t *testing.T) {
	silent := make([]byte, FrameBytes)
	got, err := FrameLevel(silent)
	if err != nil {
		t.Fatalf("a full silent seam frame must decode: %v", err)
	}
	if got != MinLevel {
		t.Fatalf("silent frame level = %v, want exactly %v", got, MinLevel)
	}
	if math.Signbit(float64(got)) {
		t.Fatalf("silent frame level is negative zero: %v", math.Float32bits(got))
	}
	if l := LevelOfSamples(make([]int16, FrameSamples)); l != 0 {
		t.Fatalf("LevelOfSamples(all-zero samples) = %v", l)
	}
	// A window with nothing in it has no energy; it reads the floor, never NaN.
	if l := LevelOfSamples(nil); math.IsNaN(l) || l != MinLevel {
		t.Fatalf("LevelOfSamples(nil) = %v, want %v and never NaN", l, MinLevel)
	}
}

// AC#1 second of the three: the loud end. Both the theoretical endpoint (1.0)
// and the strongest symmetric square (the named bound FullScaleSquareLevel)
// are pinned to the bit, plus the smallest non-zero reading of the scale.
func TestLevelFullScaleSquareEndpoints(t *testing.T) {
	// The endpoint itself: every sample at the largest magnitude an int16
	// carries. LevelFullScale is abs(-32768), so this reads exactly 1.
	hot := make([]int16, FrameSamples)
	for i := range hot {
		hot[i] = -32768
	}
	if l := LevelOfSamples(hot); l != MaxLevel {
		t.Fatalf("all-minus-32768 frame = %v, want exactly %v", l, MaxLevel)
	}
	gotHot, err := FrameLevel(EncodeFrame(hot))
	if err != nil {
		t.Fatal(err)
	}
	if float64(gotHot) != MaxLevel {
		t.Fatalf("FrameLevel(all-minus-32768) = %v, want exactly %v", gotHot, MaxLevel)
	}

	// The named upper bound for a real symmetric square wave: 32767/32768,
	// exactly one LevelLSB under 1.0 (quantisation, not a defect).
	got, err := FrameLevel(squareFrame(32767))
	if err != nil {
		t.Fatal(err)
	}
	if float64(got) != FullScaleSquareLevel {
		t.Fatalf("full-scale square = %v, want exactly %v (= 32767/32768)",
			float64(got), FullScaleSquareLevel)
	}
	if diff := MaxLevel - FullScaleSquareLevel; diff != LevelLSB {
		t.Fatalf("named bound is not one LSB below the top: diff = %v, LevelLSB = %v", diff, LevelLSB)
	}

	// One LSB of signal is the smallest thing the scale can see above silence.
	if l := LevelOfSamples(decodeOrFail(t, squareFrame(1))); l != LevelLSB {
		t.Fatalf("plus-minus-1-LSB square = %v, want exactly %v", l, LevelLSB)
	}
	// Half scale reads exactly half, a quarter exactly a quarter: the scale has
	// no gain curve hiding in it.
	for amp, want := range map[int16]float64{16384: 0.5, 8192: 0.25, 4096: 0.125} {
		if l := LevelOfSamples(decodeOrFail(t, squareFrame(amp))); l != want {
			t.Fatalf("square at %d = %v, want exactly %v", amp, l, want)
		}
	}
	// No legal frame escapes [0,1] - so no consumer ever needs to clamp.
	for _, amp := range []int16{-1, 1, 1000, 16384, 32767} {
		frame := squareFrame(amp)
		l, err := FrameLevel(frame)
		if err != nil {
			t.Fatal(err)
		}
		if float64(l) < MinLevel || float64(l) > MaxLevel {
			t.Fatalf("level %v for amp %d is outside the closed scale [%v,%v]", l, amp, MinLevel, MaxLevel)
		}
	}
}

// AC#1 third of the three: a sine against 1/sqrt(2), with the named tolerance.
// The ticket's "half-amplitude sine" is pinned under both readings, because the
// two differ by exactly the sqrt(2) identity: a FULL-scale sine reads 1/sqrt(2),
// a HALF-amplitude sine reads (1/sqrt(2))/2. Both are asserted.
func TestLevelSineAgainstRootTwo(t *testing.T) {
	rootTwo := 1 / math.Sqrt2
	full := wholeCycleSine(FrameSamples, 32767, 8)
	if l := LevelOfSamples(full); math.Abs(l-rootTwo) > SineLevelTolerance {
		t.Fatalf("full-scale sine level = %.10f, want %.10f (1/sqrt(2)) within the named tolerance %v (off by %.3e)",
			l, rootTwo, SineLevelTolerance, math.Abs(l-rootTwo))
	}
	half := wholeCycleSine(FrameSamples, 16384, 8) // 16384 is exactly half of LevelFullScale
	halfWant := rootTwo / 2
	if l := LevelOfSamples(half); math.Abs(l-halfWant) > SineLevelTolerance {
		t.Fatalf("half-amplitude sine level = %.10f, want %.10f (1/(2*sqrt(2))) within %v (off by %.3e)",
			l, halfWant, SineLevelTolerance, math.Abs(l-halfWant))
	}
	// The same numbers through the seam form, float32 and all.
	gotFull, err := FrameLevel(EncodeFrame(full))
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(float64(gotFull)-rootTwo) > SineLevelTolerance {
		t.Fatalf("FrameLevel(full-scale sine) = %.10f, want %.10f within %v", gotFull, rootTwo, SineLevelTolerance)
	}
	gotHalf, err := FrameLevel(EncodeFrame(half))
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(float64(gotHalf)-halfWant) > SineLevelTolerance {
		t.Fatalf("FrameLevel(half-amplitude sine) = %.10f, want %.10f within %v", gotHalf, halfWant, SineLevelTolerance)
	}
	// A quarter amplitude keeps the same identity (16384/8192 are exact powers
	// of two, so this is the linear scale showing, not rounding luck).
	quarter := wholeCycleSine(FrameSamples, 8192, 8)
	if l := LevelOfSamples(quarter); math.Abs(l-rootTwo/4) > SineLevelTolerance {
		t.Fatalf("quarter-amplitude sine level = %.10f, want %.10f within %v", l, rootTwo/4, SineLevelTolerance)
	}
}

// The scale is linear in amplitude: halving the input halves the number. This
// is the nail that would go red if anybody ever folded a gain or a compression
// curve into the producer (240-c1 sec 5 D-1 keeps that decision on the
// consumer side, where it belongs).
func TestLevelIsLinearInAmplitude(t *testing.T) {
	loud := LevelOfSamples(wholeCycleSine(FrameSamples, 16384, 8))
	quiet := LevelOfSamples(wholeCycleSine(FrameSamples, 8192, 8))
	// Each side is already inside SineLevelTolerance of its own ideal, so the
	// pair may not drift apart by more than twice that.
	if drift := math.Abs(2*quiet - loud); drift > 2*SineLevelTolerance {
		t.Fatalf("scale is not linear: level(8192) = %v, level(16384) = %v, drift %v over 2*tolerance %v",
			quiet, loud, drift, 2*SineLevelTolerance)
	}
	// And it is monotone: a louder frame never reads lower. The loop counter is
	// a plain int on purpose - an int16 counter would wrap past 32767 and hand
	// back a negative amplitude, which is a test bug, not a scale bug.
	prev := -1.0
	for amp := 0; amp <= 32767; amp += 512 {
		l := LevelOfSamples(decodeOrFail(t, squareFrame(int16(amp))))
		if l < prev {
			t.Fatalf("level dropped from %v to %v at amplitude %d", prev, l, amp)
		}
		prev = l
	}
}

// The tolerance itself must not be a knob that makes anything pass. Pinned
// from both sides: it has to cover the int16 quantisation it is meant to
// absorb, and it has to stay far under the smallest amplitude gap the scale is
// required to distinguish.
func TestLevelSineToleranceIsNamedAndBounded(t *testing.T) {
	if !(SineLevelTolerance >= LevelLSB) {
		t.Fatalf("SineLevelTolerance = %v is below one quantisation step %v: the sine nails would be asserting a rounding error that cannot exist", SineLevelTolerance, LevelLSB)
	}
	if !(SineLevelTolerance <= 1e-3) {
		t.Fatalf("SineLevelTolerance = %v is above 1e-3: the scale stops distinguishing normal speech levels", SineLevelTolerance)
	}
	// 100 tolerances must still be smaller than the gap between two amplitudes
	// one octave apart - otherwise "within tolerance" says nothing at all.
	gap := LevelOfSamples(wholeCycleSine(FrameSamples, 16384, 8)) -
		LevelOfSamples(wholeCycleSine(FrameSamples, 8192, 8))
	if !(gap > 100*SineLevelTolerance) {
		t.Fatalf("octave gap %v is not more than 100x the tolerance %v (tolerance too wide)", gap, SineLevelTolerance)
	}
}

// "Same input, same number" - pinned to the bit, not to a fuzzy comparison.
func TestLevelSameInputSameBits(t *testing.T) {
	frame := EncodeFrame(wholeCycleSine(FrameSamples, 25000, 5))
	first, err := FrameLevel(frame)
	if err != nil {
		t.Fatal(err)
	}
	want := math.Float32bits(first)
	for i := 0; i < 1000; i++ {
		if got := math.Float32bits(mustLevel(t, frame)); got != want {
			t.Fatalf("iteration %d returned bits %08x, first call %08x: same input must give the same number", i, got, want)
		}
	}
	// And the same samples handed over in a different slice give the same bits.
	samples := decodeOrFail(t, frame)
	clone := make([]int16, len(samples))
	copy(clone, samples)
	if LevelOfSamples(samples) != LevelOfSamples(clone) {
		t.Fatal("level depends on which slice carried the samples")
	}
}

// A non-seam frame is a bug in the caller, and it must be named loudly with the
// byte count it got: a short frame would otherwise read as a different loudness
// and look like the speaker got quiet.
func TestFrameLevelRejectsNonSeamFrames(t *testing.T) {
	cases := []struct {
		name  string
		frame []byte
	}{
		{"nil", nil},
		{"empty", []byte{}},
		{"odd one byte", []byte{0x01}},
		{"half a seam frame", make([]byte, FrameBytes/2)},
		{"one byte short", make([]byte, FrameBytes-1)},
		{"one byte long", make([]byte, FrameBytes+1)},
		{"two frames", make([]byte, 2*FrameBytes)},
	}
	for _, tc := range cases {
		got, err := FrameLevel(tc.frame)
		if err == nil {
			t.Fatalf("%s: FrameLevel must fail, returned %v", tc.name, got)
		}
		if got != MinLevel {
			t.Fatalf("%s: a rejected frame must hand back %v, got %v", tc.name, MinLevel, got)
		}
		if !strings.Contains(err.Error(), strconv.Itoa(len(tc.frame))) {
			t.Fatalf("%s: error must name the byte count it got (%d), got %q", tc.name, len(tc.frame), err.Error())
		}
	}
	// FrameBytes itself is the only accepted length, and the constants say so.
	if FrameBytes != FrameSamples*2 {
		t.Fatalf("FrameBytes = %d, FrameSamples = %d", FrameBytes, FrameSamples)
	}
	if _, err := FrameLevel(make([]byte, FrameBytes)); err != nil {
		t.Fatalf("a full silent seam frame must be accepted: %v", err)
	}
}

// AC#2 ring 2 (encodeFrame was unexported, so the frame codec could not be
// reached from outside the package): the codec is now exported in both
// directions and round-trips, including the zero-padded tail contract.
func TestSeamCodecRoundTrip(t *testing.T) {
	want := []int16{0, 1, -1, 32767, -32768, 1234, -5678}
	frame := EncodeFrame(want)
	if len(frame) != FrameBytes {
		t.Fatalf("EncodeFrame must always produce a full seam frame, got %d bytes", len(frame))
	}
	got, err := DecodeFrame(frame)
	if err != nil {
		t.Fatal(err)
	}
	for i, s := range want {
		if got[i] != s {
			t.Fatalf("sample %d decoded %d, want %d", i, got[i], s)
		}
	}
	for i := len(want); i < FrameSamples; i++ {
		if got[i] != 0 {
			t.Fatalf("tail sample %d must be zero-padded, got %d", i, got[i])
		}
	}
	if _, err := DecodeFrame([]byte{0x00, 0x01, 0x02}); err == nil {
		t.Fatal("an odd-length frame must not decode silently")
	}
	// Little-endian on purpose: pin the byte order, not just the round trip.
	if !bytes.Equal(EncodeFrame([]int16{0x0102, -1})[:4], []byte{0x02, 0x01, 0xff, 0xff}) {
		t.Fatalf("EncodeFrame is not little-endian int16: % x", EncodeFrame([]int16{0x0102, -1})[:4])
	}
}

// AC#2 ring 1 (NewBoundedFrames had no consumer): one real consumer path,
// inside the package, driven through the same C8 seam the microphone uses -
// WavInjector into the bounded channel, drained frame by frame into levels.
// The square wave is at +-16384, exactly half of LevelFullScale, so every
// frame must read exactly 0.5: no tolerance, no fuzz, end to end.
func TestLevelOverBoundedChannelFromWavInjector(t *testing.T) {
	samples := make([]int16, 3*FrameSamples) // three whole seam frames, no tail
	for i := range samples {
		if i%2 == 0 {
			samples[i] = 16384
		} else {
			samples[i] = -16384
		}
	}
	path := filepath.Join(t.TempDir(), "half-scale-square.wav")
	writeWav(t, path, samples, TargetRate, 1)

	inj, err := NewWavInjector(path)
	if err != nil {
		t.Fatal(err)
	}
	buf := NewBoundedFrames() // the channel the census found with zero consumers
	if cap(buf) != BoundedFrameCapacity() {
		t.Fatalf("bounded channel capacity = %d", cap(buf))
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := inj.Start(ctx, buf); err != nil {
		t.Fatal(err)
	}

	levels := make([]float32, 0, len(samples)/FrameSamples)
	// Drains on the test goroutine (the consumer owns its own thread; D38b has
	// no resident slot left for a level-only goroutine).
	for len(levels) < len(samples)/FrameSamples {
		select {
		case frame := <-buf:
			l, err := FrameLevel(frame)
			if err != nil {
				t.Fatalf("frame %d: %v", len(levels), err)
			}
			levels = append(levels, l)
		case <-time.After(5 * time.Second):
			t.Fatalf("timed out after %d of %d frames", len(levels), len(samples)/FrameSamples)
		}
	}
	if err := inj.Stop(); err != nil {
		t.Fatal(err)
	}

	for i, l := range levels {
		if float64(l) != 0.5 {
			t.Fatalf("frame %d level = %v, want exactly 0.5 (square wave at half of LevelFullScale)", i, l)
		}
	}
	st := inj.Stats()
	if st.FramesSent != uint64(len(levels)) || st.FramesDropped != 0 {
		t.Fatalf("stats after a drained bounded channel: %+v", st)
	}
}

// The producer and the seam have to agree on what "one frame" is, or a level
// means a different span of time in every consumer. Pinned against the
// documented 32ms/512-sample contract.
func TestLevelFrameIsOneSeamFrame(t *testing.T) {
	if FrameSamples != 512 || FrameBytes != 1024 {
		t.Fatalf("seam frame shape moved: FrameSamples=%d FrameBytes=%d", FrameSamples, FrameBytes)
	}
	samples := decodeOrFail(t, EncodeFrame(wholeCycleSine(FrameSamples, 32767, 8)))
	if len(samples) != FrameSamples {
		t.Fatalf("decoding a seam frame gave %d samples", len(samples))
	}
	if FrameDuration != 32*time.Millisecond {
		t.Fatalf("one seam frame is %v, the scale is defined per frame", FrameDuration)
	}
}

func decodeOrFail(t *testing.T, frame []byte) []int16 {
	t.Helper()
	samples, err := DecodeFrame(frame)
	if err != nil {
		t.Fatal(err)
	}
	return samples
}

func mustLevel(t *testing.T, frame []byte) float32 {
	t.Helper()
	l, err := FrameLevel(frame)
	if err != nil {
		t.Fatal(err)
	}
	return l
}
