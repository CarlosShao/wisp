package audio

import (
	"math"
	"math/rand"
	"testing"
)

// sineI16 renders a sine of the given amplitude (0..1) and frequency at rate.
func sineI16(n, rate int, freq, amp float64) []int16 {
	out := make([]int16, n)
	for i := range out {
		out[i] = int16(math.Round(amp * 32767 * math.Sin(2*math.Pi*freq*float64(i)/float64(rate))))
	}
	return out
}

// snrDB computes the signal-to-noise ratio (dB) of got against want.
func snrDB(t *testing.T, want, got []int16) float64 {
	t.Helper()
	if len(want) != len(got) {
		t.Fatalf("length mismatch: want %d got %d", len(want), len(got))
	}
	var sig, noise float64
	for i := range want {
		w := float64(want[i])
		e := float64(got[i]) - w
		sig += w * w
		noise += e * e
	}
	if noise == 0 {
		return math.Inf(1)
	}
	return 10 * math.Log10(sig/noise)
}

// TestResample48kTo16kLengthExact pins the deterministic length contract for
// the main capture case (48kHz device-native -> 16kHz seam rate).
func TestResample48kTo16kLengthExact(t *testing.T) {
	for _, n := range []int{1, 2, 3, 480, 4800, 48000} {
		in := sineI16(n, 48000, 440, 0.5)
		got := ResampleLinear(in, 48000, 16000)
		// Independent length formula: outLen = ceil((n-1)*to/from), here ceil((n-1)/3).
		wantLen := (n - 1 + 2) / 3
		if len(got) != wantLen {
			t.Fatalf("n=%d: out len %d, want %d", n, len(got), wantLen)
		}
	}
}

// TestResamplerSineSNR48k checks frequency fidelity: a 1kHz sine downsampled
// 48k->16k must match a directly sampled 1kHz sine at 16k with healthy SNR
// (linear interpolation error for this tone is ~ -53dB; the floor guards
// against gross errors like off-by-one phase or DC bias, not perfection).
func TestResamplerSineSNR48k(t *testing.T) {
	const freq = 1000.0
	const n = 48000
	got := ResampleLinear(sineI16(n, 48000, freq, 0.8), 48000, 16000)
	want := sineI16(len(got), TargetRate, freq, 0.8)
	// Skip the first 16 output samples: the linear interpolator has group
	// delay of <1 input sample, but the direct-render reference starts with
	// full amplitude while interpolation slightly smooths the onset.
	sn := snrDB(t, want[16:], got[16:])
	if sn < 40 {
		t.Fatalf("1kHz SNR = %.1f dB, want >= 40 dB", sn)
	}
}

// TestResamplerSineSNR44k1 covers a non-integer ratio (step = 2.75625), the
// actual shape of many laptop microphones.
func TestResamplerSineSNR44k1(t *testing.T) {
	const freq = 440.0
	const n = 44100
	got := ResampleLinear(sineI16(n, 44100, freq, 0.8), 44100, TargetRate)
	want := sineI16(len(got), TargetRate, freq, 0.8)
	sn := snrDB(t, want[16:], got[16:])
	if sn < 35 {
		t.Fatalf("440Hz@44.1k SNR = %.1f dB, want >= 35 dB", sn)
	}
}

// TestResamplerChunkInvariant pins the streaming contract: processing in
// arbitrary chunks must equal the one-shot result bit for bit.
func TestResamplerChunkInvariant(t *testing.T) {
	rng := rand.New(rand.NewSource(13))
	in := make([]int16, 4823)
	for i := range in {
		in[i] = int16(rng.Intn(65536) - 32768)
	}
	oneShot := ResampleLinear(in, 48000, 16000)

	r := NewResampler(48000, 16000)
	var streamed []int16
	for off := 0; off < len(in); {
		size := 1 + rng.Intn(700)
		if off+size > len(in) {
			size = len(in) - off
		}
		streamed = append(streamed, r.Process(in[off:off+size])...)
		off += size
	}
	if len(streamed) != len(oneShot) {
		t.Fatalf("streamed len %d != one-shot len %d", len(streamed), len(oneShot))
	}
	for i := range streamed {
		if streamed[i] != oneShot[i] {
			t.Fatalf("sample %d differs: streamed %d one-shot %d", i, streamed[i], oneShot[i])
		}
	}
}

// TestResamplerLatencyBounded asserts the one-lookahead-sample latency bound:
// every output sample only needs input up to its own position + 1 sample.
// Feeding one sample at a time must still produce the one-shot output.
func TestResamplerLatencyBounded(t *testing.T) {
	in := sineI16(4800, 48000, 1000, 0.5)
	oneShot := ResampleLinear(in, 48000, 16000)
	r := NewResampler(48000, 16000)
	var out []int16
	for _, s := range in {
		out = append(out, r.Process([]int16{s})...)
	}
	if len(out) != len(oneShot) {
		t.Fatalf("sample-by-sample len %d != one-shot %d", len(out), len(oneShot))
	}
	for i := range out {
		if out[i] != oneShot[i] {
			t.Fatalf("sample %d differs", i)
		}
	}
}

// TestResamplerPassthroughAndFlush pins the degenerate cases.
func TestResamplerPassthroughAndFlush(t *testing.T) {
	in := []int16{1, -2, 3, -4}
	got := ResampleLinear(in, 16000, 16000)
	if len(got) != len(in) {
		t.Fatalf("passthrough changed length")
	}
	for i := range in {
		if got[i] != in[i] {
			t.Fatalf("passthrough changed samples")
		}
	}
	r := NewResampler(48000, 16000)
	r.Process(sineI16(1000, 48000, 440, 0.5))
	r.Flush() // must reset: reuse for a fresh signal
	again := r.Process(sineI16(6, 48000, 440, 0.5))
	if len(again) != 2 { // ceil((6-1)/3)
		t.Fatalf("after Flush got %d samples, want 2", len(again))
	}
	if NewResampler(0, 16000) != nil || NewResampler(48000, -1) != nil {
		t.Fatal("non-positive rates must be rejected")
	}
}

// TestMonoDownmix and float conversion cover the WASAPI mix-format path.
func TestMonoDownmixAndFloatConvert(t *testing.T) {
	stereo := []int16{1000, 2000, -1000, -2000}
	mono := MonoDownmix(stereo, 2)
	if len(mono) != 2 || mono[0] != 1500 || mono[1] != -1500 {
		t.Fatalf("downmix wrong: %v", mono)
	}
	same := MonoDownmix(stereo, 1)
	if len(same) != 4 || &same[0] != &stereo[0] {
		t.Fatalf("chans==1 must be a no-op")
	}
	f := []float32{0, 1.0, -1.0, 0.25}
	pcm := FloatToPCM16(f)
	if pcm[0] != 0 || pcm[1] != 32767 || pcm[2] != -32768 || pcm[3] != 8192 {
		t.Fatalf("float->pcm wrong: %v", pcm)
	}
}
