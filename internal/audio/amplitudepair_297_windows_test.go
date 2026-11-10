//go:build windows

package audio

// Ticket 297 AC#1: pin "can the scale separate at all" with KNOWN amplitudes
// driven through the injection seam, so the pair reading
//
//	injection face (this test):  two windows of known amplitude, mean ratio
//	                             equals the amplitude ratio
//	real-mic face (ticket 247's six-window run, and the AC#0 live
//	               instrument): the same ruler, real device, ratio 1.04
//
// lands on disk before anyone says "the problem is the ruler". The ticket
// wording: feed two same-length windows of known amplitude (0.05 and 0.4)
// through the C8 AudioSource seam and require mean(windowB)/mean(windowA) to
// equal 8.0 (the amplitude ratio). If that holds here while the real-mic face
// stays at 1.04, the ruler LevelOfSamples is exonerated and the suspicion
// moves to the device/conversion side -- exactly the pair the ticket asks
// for. This also pays the "theoretical RMS" half of ticket 291's named debt
// (B3); it deliberately does NOT tick any ticket 291 box.
//
// Rig: the same shape as capturelevel_windows_test.go -- fake watcher/opener
// at the enumeration seam, REAL capture loop (wait/drain/resample/frame/
// bounded push/emitLevel), real wav bytes parsed by the production parser.
// No goroutine is created by this file; the loop's pinned thread is booked by
// SpawnCapture as in production.
//
// Signal choice: DC windows, not sine. DC makes the level exact (RMS of a
// steady value v is |v|), so the ratio check cannot hide behind a waveform
// argument; the level scale doc (level.go) defines it as plain RMS, and the
// DC-counts clause there says consumers gate DC themselves -- this test is
// the producer-side instrument, not a consumer.

import (
	"context"
	"math"
	"os"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// TestAC297KnownAmplitudePairSeparates drives one capture loop over a wav
// holding window A (amplitude ~0.05, 8 seam frames) then window B
// (amplitude ~0.4, 8 seam frames) and requires the frame levels to match the
// source windows exactly and the mean ratio to be the amplitude ratio.
func TestAC297KnownAmplitudePairSeparates(t *testing.T) {
	const (
		ampLow  = 0.05
		ampHigh = 0.4
		frames  = 8 // seam frames per window (8 * 512 samples = 256ms each)
	)
	// int16 counts: LevelFullScale-1 keeps the largest amplitude inside the
	// symmetric range, matching level.go's endpoint definition.
	lowCounts := int16(math.Round(ampLow * 32767))
	highCounts := int16(math.Round(ampHigh * 32767))

	var samples []int16
	for i := 0; i < frames*FrameSamples; i++ {
		samples = append(samples, lowCounts)
	}
	for i := 0; i < frames*FrameSamples; i++ {
		samples = append(samples, highCounts)
	}

	path := t.TempDir() + "/ac297-pair.wav"
	writeWav(t, path, samples, TargetRate, 1)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	parsed, rate, chans, err := parseWav(raw)
	if err != nil || rate != TargetRate || chans != 1 {
		t.Fatalf("parseWav: err=%v rate=%d chans=%d", err, rate, chans)
	}

	dev := DeviceDescriptor{ID: "dev-ac297", Name: "AC297 Known-Amplitude Mic"}
	watcher := &fakeWatcher{devs: []DeviceDescriptor{dev}}
	stream := &fakeStream{
		dev: dev, rate: TargetRate, period: 10 * time.Millisecond,
		chunks: makeChunks(parsed, 160), // 10ms periods, like the 247 rig
	}
	opener := &fakeOpener{streams: []*fakeStream{stream}}

	rec := &levelRecorder{}
	reg := observe.NewRegistry()
	mic := newWASAPIMicrophoneWith(watcher, opener,
		WithCaptureRegistry(reg),
		WithLevelSink(rec.call),
	)
	buf := make(chan []byte, 64)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := mic.Start(ctx, buf); err != nil {
		t.Fatal(err)
	}

	got := waitForLevels(t, rec, frames*2, 5*time.Second)
	_ = mic.Stop()

	// Frame-by-frame: the real loop's emitLevel must reproduce the source
	// window's level exactly (LevelOfSamples over the same 512 samples).
	for i, l := range got {
		want := float32(LevelOfSamples(samples[i*FrameSamples : (i+1)*FrameSamples]))
		if l != want {
			t.Fatalf("frame %d: level = %v, want %v", i, l, want)
		}
	}

	meanLow := meanF32(got[:frames])
	meanHigh := meanF32(got[frames : frames*2])
	ratio := float64(meanHigh) / float64(meanLow)
	ampRatio := float64(highCounts) / float64(lowCounts)

	t.Logf("AC#1 injection face: window A (amp %d counts) mean=%.6f, window B (amp %d counts) mean=%.6f",
		lowCounts, meanLow, highCounts, meanHigh)
	t.Logf("AC#1 injection face: mean ratio = %.4f, amplitude ratio = %.4f -- the ruler SEPARATES at the injection seam",
		ratio, ampRatio)
	t.Logf("AC#1 pair half: same ruler on the real-mic face = ticket 247 six-window run (sound/quiet-A = 1.0402); paired reading archived in probes/297/")

	if math.Abs(ratio-ampRatio) > 0.01 {
		t.Fatalf("mean ratio = %.4f, amplitude ratio = %.4f (drift > 0.01)", ratio, ampRatio)
	}
	wantLow := math.Abs(float64(lowCounts)) / LevelFullScale
	if diff := math.Abs(float64(meanLow) - wantLow); diff > SineLevelTolerance {
		t.Fatalf("window A mean = %.8f, want %.8f (|counts|/FullScale)", meanLow, wantLow)
	}
}

// meanF32 averages a slice of levels.
func meanF32(vs []float32) float64 {
	var sum float64
	for _, v := range vs {
		sum += float64(v)
	}
	return sum / float64(len(vs))
}
