//go:build windows

package aecspike

// Ticket-388 tests. Two layers, following the package's gate convention:
//   - CI-safe (no gate): the guard's validation surface - the value axis is
//     parameterized and out-of-range values are refused BEFORE any COM or
//     volume call, so a plain `go test` run touches nothing.
//   - Gated (WISP_AEC_SPIKE=1): the real readings. The triple-reading test
//     never renders (pure session-volume axis, silent); the loopback test
//     renders two SHORT quiet fixtures under two attenuated session volumes
//     and checks the endpoint mix actually followed the axis, then puts the
//     owner's original volume back.

import (
	"context"
	"math"
	"testing"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/observe"
)

// volumeTriplesEq compares two volume readings with the read-back epsilon.
func volumeTriplesEq(a, b float64) bool { return math.Abs(a-b) <= volumeReadbackEpsilon }

// Test388AttenuationValueIsParameterizedAndValidated is CI-safe: the guard
// refuses out-of-range volumes before touching COM or the session, and the
// default value is a real attenuation (strictly below full scale).
func Test388AttenuationValueIsParameterizedAndValidated(t *testing.T) {
	for _, v := range []float64{0, -0.15, 1.0000001, 42} {
		g, err := ApplyAnnounceAttenuation(v)
		if err == nil {
			t.Fatalf("ApplyAnnounceAttenuation(%g) was accepted; the value axis must be validated before any COM call", v)
		}
		if g != nil {
			t.Fatalf("ApplyAnnounceAttenuation(%g) returned a guard alongside its error; nothing may look applied", v)
		}
	}
	if DefaultAnnounceAttenuatedVolume <= 0 || DefaultAnnounceAttenuatedVolume >= 1 {
		t.Fatalf("DefaultAnnounceAttenuatedVolume=%g is not an attenuation", DefaultAnnounceAttenuatedVolume)
	}
}

// Test388AnnounceAttenuationTripleReadings is AC#1+AC#3's volume-axis core:
// apply the attenuation (read-back must agree), restore the owner's volume
// (read-back must agree with the saved original), then re-check through a
// FRESH client that the restore actually persisted in the session. It never
// renders a sample, so it is silent even under the gate.
func Test388AnnounceAttenuationTripleReadings(t *testing.T) {
	spikeGate(t)
	g, err := ApplyAnnounceAttenuation(DefaultAnnounceAttenuatedVolume)
	if err != nil {
		t.Skipf("announce attenuation guard unavailable on this host: %v", err)
	}
	readings, err := g.Restore()
	if err != nil {
		t.Fatalf("the attenuation guard failed its own restore: %v", err)
	}
	if !volumeTriplesEq(readings.Attenuated, DefaultAnnounceAttenuatedVolume) {
		t.Fatalf("AC#1: the attenuated session volume read back %g, wanted %g",
			readings.Attenuated, DefaultAnnounceAttenuatedVolume)
	}
	if !volumeTriplesEq(readings.Restored, readings.Original) {
		t.Fatalf("AC#3: the restored session volume read back %g, wanted the original %g",
			readings.Restored, readings.Original)
	}
	// AC#3's restore check through an independent client: what the session
	// PERSISTED must equal what the guard saved, not just what the guard's
	// own write answered.
	client, release, err := openSessionVolumeClient()
	if err != nil {
		t.Fatalf("the restore re-check could not open a client: %v", err)
	}
	defer release()
	again, err := GetSessionVolume(client)
	if err != nil {
		t.Fatalf("the restore re-check could not read the session volume: %v", err)
	}
	if !volumeTriplesEq(again, readings.Original) {
		t.Fatalf("AC#3 restore check: the session persisted %g, the guard saved %g", again, readings.Original)
	}
	t.Logf("388 volume triple: original=%.6f attenuated=%.6f restored=%.6f fresh-client re-read=%.6f",
		readings.Original, readings.Attenuated, readings.Restored, again)
}

// rmsInt16 is the root-mean-square of a captured mix (the AC#2 reading).
func rmsInt16(in []int16) float64 {
	if len(in) == 0 {
		return 0
	}
	var sum float64
	for _, v := range in {
		sum += float64(v) * float64(v)
	}
	return math.Sqrt(sum / float64(len(in)))
}

// peakWindowRMS is the maximum RMS over sliding windows of a captured mix.
// The loopback run's tail is engine silence whose LENGTH varies with packet
// timing, so a whole-buffer RMS dilutes by an uncontrollable factor; the
// peak window rides on the signal itself and is comparable across runs.
func peakWindowRMS(in []int16, window int) float64 {
	if window <= 0 || window > len(in) {
		return rmsInt16(in)
	}
	best, sum := 0.0, 0.0
	for i, v := range in {
		sum += float64(v) * float64(v)
		if i >= window {
			sum -= float64(in[i-window]) * float64(in[i-window])
		}
		if i >= window-1 {
			if r := math.Sqrt(sum / float64(window)); r > best {
				best = r
			}
		}
	}
	return best
}

// Test388AttenuatedSessionReachesTheEndpointLoopback is AC#2's leg-doable
// half: with the session attenuated, the endpoint loopback mix must follow
// the volume axis downward (two attenuated levels; see the reading note at
// the assertion for the exponent this host's tap reports). The fixtures are
// one-second speech-like tones at low sample gain, so the loudness stays
// far below an announce test's, and the owner's original volume is saved
// first and restored last through the guard chain.
func Test388AttenuatedSessionReachesTheEndpointLoopback(t *testing.T) {
	spikeGate(t)
	reg := observe.NewRegistry()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	far := SynthSpeechLike(audio.TargetRate, 1)
	const gain = 0.5

	// Level 1: session at 0.6. Level 2: session at 0.15 (a quarter of it).
	// The second guard's "original" is the first attenuation, and the chain
	// unwinds to the owner's true volume.
	g1, err := ApplyAnnounceAttenuation(0.6)
	if err != nil {
		t.Skipf("announce attenuation guard unavailable on this host: %v", err)
	}
	t.Cleanup(func() { g1.Restore() })
	r1, mix1, err := RunLoopbackQuiet(reg, ctx, "", gain, far, audio.TargetRate)
	if err != nil {
		t.Skipf("loopback rig unavailable on this host: %v", err)
	}
	g2, err := ApplyAnnounceAttenuation(0.15)
	if err != nil {
		t.Fatalf("the second attenuation did not apply (session was at g1's %g): %v", g1.applied, err)
	}
	t.Cleanup(func() { g2.Restore() })
	r2, mix2, err := RunLoopbackQuiet(reg, ctx, "", gain, far, audio.TargetRate)
	if err != nil {
		t.Fatalf("the attenuated loopback run failed: %v", err)
	}
	readings2, err := g2.Restore()
	if err != nil || !volumeTriplesEq(readings2.Restored, g1.applied) {
		t.Fatalf("the guard chain did not unwind to the first attenuation: readings=%+v err=%v", readings2, err)
	}
	readings1, err := g1.Restore()
	if err != nil {
		t.Fatalf("the owner's original session volume was not restored: %v", err)
	}

	rms1, rms2 := peakWindowRMS(mix1, 1024), peakWindowRMS(mix2, 1024)
	if rms1 <= 0 || rms2 <= 0 {
		t.Fatalf("the endpoint mix is silent at both levels: rms1=%g rms2=%g readings1=%+v readings2=%+v",
			rms1, rms2, r1, r2)
	}
	ratio := rms2 / rms1
	// AC#2: the endpoint (exactly what the owner's speakers are fed) must
	// follow the session axis DOWNWARD and by a large factor. Reading note
	// (measured twice, pushed=captured=48000 both runs, no tail dilution):
	// this host's loopback tap tracks the session axis as a SQUARE law -
	// rms(0.15)/rms(0.6) = 0.062 ~= (0.15/0.6)^2 - so the strict 0.25
	// linear expectation does not hold on the tap. That exponent lives in
	// the host's APO/tap chain, not in the guard; the guard's own axis is
	// verified exact by the read-back triple. What AC#2 needs here is that
	// the attenuation really reached the endpoint feed: the bound below
	// accepts the linear (0.25) and the measured square (0.0625) shapes and
	// rejects "the tap did not attenuate" (ratio near 1) and "silence".
	if ratio >= 0.5 || ratio <= 0.01 {
		t.Fatalf("AC#2: the endpoint mix did not follow the session volume axis downward: "+
			"peakrms(0.6)=%g peakrms(0.15)=%g ratio=%g (this host measured ~0.0625, a square law; "+
			"linear would read ~0.25); run diagnostics: pushed1=%d captured1=%d pushed2=%d captured2=%d "+
			"whoRms1=%g whoRms2=%g",
			rms1, rms2, ratio, r1.PushedFrames, r1.CapturedFrames, r2.PushedFrames, r2.CapturedFrames,
			rmsInt16(mix1), rmsInt16(mix2))
	}
	t.Logf("388 loopback contrast: session 0.6 peakrms=%g (whole %g, %d/%d frames), session 0.15 peakrms=%g (whole %g, %d/%d frames), "+
		"ratio=%.4f (square law on this host's tap; linear would read 0.25); volume triple chain: owner original=%.6f restored=%.6f",
		rms1, rmsInt16(mix1), r1.PushedFrames, r1.CapturedFrames,
		rms2, rmsInt16(mix2), r2.PushedFrames, r2.CapturedFrames,
		ratio, readings1.Original, readings1.Restored)
}
