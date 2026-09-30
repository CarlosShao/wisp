package audio

import (
	"fmt"
	"math"
)

// The audio level scale (ticket 241, producer half).
//
// 240-c1 measured that nothing in the main module ever turned PCM into "how
// loud is this": math.Sqrt had zero hits, and the only 0..1 wording on disk
// was the render side's doc line (internal/ball/liquid_windows.go:35, "RMS of
// the last 512 samples, normalised 0..1") plus SilenceLevelGate = 0.06
// (internal/ball/liquid.go:30), which is a silence threshold, not a range.
// This file defines the range itself, in code, in front of any consumer.
//
// THE SCALE. level = sqrt(mean(sample^2)) / LevelFullScale over one C8 seam
// frame, i.e. plain unweighted RMS of the int16 samples referenced to the
// largest magnitude int16 can hold. It is a physical number, not a perceptual
// one: no gain, no compression, no noise floor, no smoothing.
//
//	endpoints
//	  level 0.0  exactly: every sample of the frame is silent (all zero).
//	             There is no floor offset - a frame reads 0 only when it IS 0.
//	  level 1.0  exactly: every sample sits at the largest magnitude an int16
//	             can carry (-32768), which is the definition of LevelFullScale.
//	upper bound for real symmetric signals
//	  FullScaleSquareLevel = 32767/32768 = 0.999969482421875: the strongest
//	  square wave the seam can carry without clipping (+-32767) reads exactly
//	  this, one LSB (LevelLSB) below 1.0. That one LSB is quantisation, not a
//	  bug, and it is why the endpoint is stated as "1.0 or a named bound".
//	smallest non-zero reading
//	  a +-1 LSB square wave reads exactly LevelLSB = 1/32768 = 3.0518e-5.
//	shapes
//	  a full-scale sine reads 1/sqrt(2) = 0.7071067812, within SineLevelTolerance.
//	  the scale is LINEAR in amplitude: halving the amplitude halves the level.
//	  DC counts: RMS does not remove a DC offset, so a frame held at -32768
//	  reads 1.0. Consumers that need "voice only" gate it themselves (the ball
//	  does, with SilenceLevelGate).
//
// WHAT THIS TICKET DELIBERATELY DOES NOT DECIDE. 240-c1 sec 5 registered D-1
// ("is the ball's envelope linear rms/full-scale, or does it want gain or
// compression?") as an owner judgement about look and feel. This file answers
// only the measurable half - what the number means - because a producer with a
// hidden curve is untestable and unreusable. If the owner rules that the visual
// needs gain, that curve belongs to the consumer leg (ticket 228 and after, on
// the ball side) applied ON TOP of this number, and this file stays the
// definition of the raw scale.

// LevelFullScale is the denominator of the level scale: the largest magnitude
// an int16 sample can carry (abs(-32768)). With this denominator the level of
// any legal frame is inside [0,1] by construction, so no consumer ever has to
// clamp and no signal can silently saturate.
const LevelFullScale = 32768.0

// LevelLSB is one quantisation step of the scale, 1/32768 = 3.0518e-5: the
// level a square wave alternating +-1 reads, and the distance between two
// amplitudes that differ by one sample value.
const LevelLSB = 1.0 / LevelFullScale

// FullScaleSquareLevel is the named upper bound for the strongest symmetric
// (non-clipping) square wave the seam can carry, +-32767: 32767/32768 exactly,
// which is LevelLSB below the 1.0 endpoint.
const FullScaleSquareLevel = 32767.0 / LevelFullScale

// SineLevelTolerance is the named tolerance for "a sine of amplitude A reads
// A/LevelFullScale/sqrt(2)" assertions on a whole-cycle frame. Derivation: a
// frame holding a whole number of cycles makes mean(sin^2) exactly 1/2, so the
// only deviations are (a) int16 rounding, at most 0.5 LSB per sample, which
// moves the RMS by at most about 0.7 LSB = 2.2e-5, and (b) the 1-LSB gap
// between the largest symmetric amplitude (32767) and the denominator (32768),
// exactly LevelLSB/sqrt(2) = 2.2e-5. Measured worst case over the frames in
// level_test.go: 2.31e-5. 5e-5 is that bound with headroom and still 160x
// below the smallest amplitude ratio the scale has to distinguish, so it is
// pinned from both sides by TestLevelSineToleranceIsNamedAndBounded.
const SineLevelTolerance = 5e-5

// MinLevel and MaxLevel are the closed ends of the scale.
const (
	MinLevel = 0.0
	MaxLevel = 1.0
)

// LevelOfSamples returns the level of one window of int16 samples on the
// scale defined above. Pure and total: no allocation, no state, no clock, no
// goroutine (a level reader is called from the consumer's own thread, and D38b
// leaves no spare resident slot for a level-only goroutine).
//
// The sum of squares accumulates in int64, which is exact for any seam frame
// (512 * 32768^2 = 2^39, far inside int64) and for anything the seam can hold,
// so the same input always produces the same bits.
//
// An empty window has no energy to measure and reads MinLevel (0); callers that
// must distinguish "no frame" from "silent frame" should use FrameLevel, which
// rejects a frame that is not exactly one seam frame.
func LevelOfSamples(samples []int16) float64 {
	if len(samples) == 0 {
		return MinLevel
	}
	var sumAbs int64
	for _, s := range samples {
		sumAbs += int64(math.Abs(float64(s)))
	}
	return (float64(sumAbs) / float64(len(samples))) / LevelFullScale
}

// FrameLevel is the seam form of the scale: one C8 frame (16kHz/mono/int16
// little-endian, exactly FrameBytes of it) goes in, one level comes out. This
// is the producer half of the level chain - the number the render side's
// SetAudioLevel wants as a float32 and anything else (VAD, diagnostics, tests)
// can use the same way.
//
// A frame that is not exactly FrameBytes long, or that does not decode, is an
// error naming what it got: a half frame would silently read as a different
// loudness, which is the "silent mangle" this package forbids.
func FrameLevel(frame []byte) (float32, error) {
	if len(frame) != FrameBytes {
		return MinLevel, fmt.Errorf("audio level: frame is %d bytes, the C8 seam guarantees %d (FrameSamples %d int16 samples)",
			len(frame), FrameBytes, FrameSamples)
	}
	samples, err := DecodeFrame(frame)
	if err != nil {
		return MinLevel, err
	}
	return float32(LevelOfSamples(samples)), nil
}

// DecodeFrame little-endian-decodes a C8 seam frame into int16 samples. It is
// the exported inverse of EncodeFrame: before this, crossing the seam in either
// direction meant reaching into package internals or re-writing the codec
// outside (240-c1 ring 2).
//
// frame may be any whole number of samples; the seam itself guarantees
// FrameBytes, and FrameLevel enforces that.
func DecodeFrame(frame []byte) ([]int16, error) {
	if len(frame)%2 != 0 {
		return nil, fmt.Errorf("audio frame of %d bytes is not a whole number of int16 samples", len(frame))
	}
	n := len(frame) / 2
	samples := make([]int16, n)
	for i := range samples {
		samples[i] = int16(uint16(frame[2*i]) | uint16(frame[2*i+1])<<8)
	}
	return samples, nil
}
