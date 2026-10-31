package aecspike

import (
	"math"
	"os"
	"path/filepath"
	"testing"
)

// Pure-math pins for the spike instruments. Everything here runs on every
// GOOS in CI; the hardware rig (rig_windows_test.go) is gated separately.

// echoFixture builds far + near = far delayed by lag, scaled by scale, plus
// a small independent noise floor on the near side. The far signal carries
// a slow amplitude envelope on purpose: a steady tone has a FLAT envelope,
// and then the coarse envelope-correlation search is ill-posed (every lag
// ties); real speech always fluctuates, and so must the fixture.
func echoFixture(lag, n int, scale, noise float32) (far, near []float32) {
	far = make([]float32, n)
	for i := range far {
		t := float64(i) / 16000
		env := 0.4 + 0.6*math.Abs(math.Sin(2*math.Pi*0.7*t))*math.Abs(math.Sin(2*math.Pi*2.3*t))
		far[i] = float32(env * (0.6*math.Sin(2*math.Pi*220*t) + 0.3*math.Sin(2*math.Pi*660*t+1.0)))
	}
	near = make([]float32, n)
	for i := range near {
		v := 0.0
		if j := i - lag; j >= 0 && j < len(far) {
			v = float64(scale * far[j])
		}
		v += float64(noise) * math.Sin(2*math.Pi*3000*float64(i)/16000)
		near[i] = float32(v)
	}
	return far, near
}

func TestAlignRecoversKnownLag(t *testing.T) {
	for _, lag := range []int{0, 160, 1600, 3200} {
		far, near := echoFixture(lag, 16000*2, 0.5, 0.001)
		got, strength := AlignByEnvelopeCorrelation(far, near, DefaultMaxLag)
		if diff := got - lag; diff < -envDecim || diff > envDecim {
			t.Fatalf("lag %d: AlignByEnvelopeCorrelation = %d (want within %d)", lag, got, envDecim)
		}
		if strength < 0.99 {
			t.Fatalf("lag %d: strength = %.4f, want ~1 for a pure delayed echo", lag, strength)
		}
	}
}

func TestSubtractEchoDrivesResidualDown(t *testing.T) {
	far, near := echoFixture(1600, 16000*2, 0.5, 0.001)
	lag, strength := AlignByEnvelopeCorrelation(far, near, DefaultMaxLag)
	residual, gain, _ := SubtractEcho(far, near, lag)
	if math.Abs(gain-0.5) > 0.01 {
		t.Fatalf("gain = %.4f, want ~0.5", gain)
	}
	if strength < 0.99 {
		t.Fatalf("strength = %.4f, want ~1", strength)
	}
	erle := ERLEdB(far, near, residual, lag)
	// The noise floor is -60dBFS; a 0.5-scaled echo at full scale cancels to
	// well below it, so ERLE must be large (upper bound > 30dB).
	if erle < 30 {
		t.Fatalf("ERLE = %.2f dB, want > 30 for near-perfect cancellation", erle)
	}
}

func TestSubtractEchoLeavesIndependentSignal(t *testing.T) {
	// near is an INDEPENDENT tone (external-audio shape): subtraction must
	// leave it almost untouched, i.e. ERLE near 0dB, not a fake suppression.
	far, _ := echoFixture(0, 16000, 1, 0)
	near := make([]float32, 16000)
	for i := range near {
		near[i] = float32(0.4 * math.Sin(2*math.Pi*900*float64(i)/16000))
	}
	lag, strength := AlignByEnvelopeCorrelation(far, near, DefaultMaxLag)
	if strength > 0.2 {
		t.Fatalf("independent signals correlated %.3f, want ~0", strength)
	}
	residual, _, _ := SubtractEcho(far, near, lag)
	var nearPow, resPow float64
	for i := range near {
		nearPow += float64(near[i]) * float64(near[i])
		resPow += float64(residual[i]) * float64(residual[i])
	}
	erle := 10 * math.Log10(nearPow/resPow)
	if math.Abs(erle) > 1 {
		t.Fatalf("ERLE = %.2f dB on independent near, want ~0", erle)
	}
}

func TestPpmDriftDegradesCancellationMonotonically(t *testing.T) {
	far, near := echoFixture(1600, 16000*3, 0.5, 0.001)
	// Rebuild near from the DRIFTED far so the echo really carries the drift.
	drifted := ApplyPpmDrift(far, 200)
	near = make([]float32, len(drifted))
	copy(near, drifted)
	for i := range near { // undo the scale so subtraction has to work
		near[i] *= 0.5
	}
	good := ApplyPpmDrift(far, 0)
	lag, _ := AlignByEnvelopeCorrelation(good, near, DefaultMaxLag)
	res, _, _ := SubtractEcho(good, near, lag)
	erle := ERLEdB(good, near, res, lag)
	if erle > 10 {
		t.Fatalf("200ppm drift: ERLE = %.1f dB, want heavy degradation (<10)", erle)
	}
	// The zero-drift control must cancel cleanly.
	lag0, _ := AlignByEnvelopeCorrelation(far, ApplyPpmDrift(far, 0), DefaultMaxLag)
	res0, _, _ := SubtractEcho(far, ApplyPpmDrift(far, 0), lag0)
	erle0 := ERLEdB(far, ApplyPpmDrift(far, 0), res0, lag0)
	if erle0 < 30 {
		t.Fatalf("zero-drift control: ERLE = %.1f dB, want > 30", erle0)
	}
}

func TestWindowedCorrelationSeparatesEchoFromIndependent(t *testing.T) {
	far, nearEcho := echoFixture(1600, 16000, 0.5, 0.001)
	_, nearForeign := echoFixture(0, 16000, 1, 0)
	foreign := make([]float32, 16000)
	for i := range foreign {
		foreign[i] = float32(0.4 * math.Sin(2*math.Pi*900*float64(i)/16000))
	}
	_ = nearForeign
	echoCorr := WindowedCorrelation(far, nearEcho, 1600, 1024, 512)
	foreignCorr := WindowedCorrelation(far, foreign, 0, 1024, 512)
	var echoMean, foreignMean float64
	for _, v := range echoCorr {
		echoMean += v
	}
	echoMean /= float64(len(echoCorr))
	for _, v := range foreignCorr {
		foreignMean += v
	}
	foreignMean /= float64(len(foreignCorr))
	if echoMean < 0.9 {
		t.Fatalf("self-echo mean correlation %.3f, want >0.9", echoMean)
	}
	if foreignMean > 0.1 {
		t.Fatalf("external-source mean correlation %.3f, want ~0", foreignMean)
	}
}

func TestFixtureGeneratorsAreDeterministicAndAudible(t *testing.T) {
	for name, fn := range map[string]func() []float32{
		"music":  func() []float32 { return SynthMusic(16000, 1) },
		"video":  func() []float32 { return SynthVideo(16000, 1, SynthSpeechLike(16000, 1)) },
		"speech": func() []float32 { return SynthSpeechLike(16000, 1) },
	} {
		a, b := fn(), fn()
		if len(a) != 16000 {
			t.Fatalf("%s: len = %d, want 16000", name, len(a))
		}
		for i := range a {
			if a[i] != b[i] {
				t.Fatalf("%s: not deterministic at sample %d", name, i)
			}
		}
		peak := 0.0
		for _, v := range a {
			if p := math.Abs(float64(v)); p > peak {
				peak = p
			}
		}
		if peak < 0.1 {
			t.Fatalf("%s: peak %.3f, fixture is near-silent", name, peak)
		}
	}
}

func TestWavRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "fixture.wav")
	samples := make([]int16, 1600)
	for i := range samples {
		samples[i] = int16(i % 1000)
	}
	if err := WriteWav16(path, 16000, samples); err != nil {
		t.Fatalf("WriteWav16: %v", err)
	}
	rate, got, err := ReadWav16(path)
	if err != nil {
		t.Fatalf("ReadWav16: %v", err)
	}
	if rate != 16000 || len(got) != len(samples) {
		t.Fatalf("roundtrip rate=%d len=%d, want 16000/%d", rate, len(got), len(samples))
	}
	for i := range samples {
		if got[i] != samples[i] {
			t.Fatalf("roundtrip sample %d: %d != %d", i, got[i], samples[i])
		}
	}
	if _, _, err := ReadWav16(filepath.Join(dir, "missing.wav")); !os.IsNotExist(err) {
		t.Fatalf("missing file error = %v", err)
	}
}

func TestResampleLinearKeepsLengthAndEnergy(t *testing.T) {
	in := make([]float32, 22050) // 1s at 22050
	for i := range in {
		in[i] = float32(math.Sin(2 * math.Pi * 440 * float64(i) / 22050))
	}
	out := ResampleLinear(in, 22050, 16000)
	if len(out) < 15800 || len(out) > 16200 {
		t.Fatalf("resampled len = %d, want ~16000", len(out))
	}
	var inPow, outPow float64
	for _, v := range in {
		inPow += float64(v) * float64(v)
	}
	for _, v := range out {
		outPow += float64(v) * float64(v)
	}
	if math.Abs(outPow/float64(len(out))/(inPow/float64(len(in)))-1) > 0.05 {
		t.Fatalf("energy ratio drifted beyond 5%%: %f", outPow/float64(len(out))/(inPow/float64(len(in))))
	}
}
