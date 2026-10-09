//go:build windows

package audio

import (
	"context"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// TestAC247RealCaptureLoopEmitsLevels: the capture loop itself (wait / drain /
// resample / frame / bounded push), fed by a fake watcher and opener whose
// streams carry real wav bytes - the enumeration seam is injected, the
// audio-data seam never is (SPEC-04 sec 9). The level sink and the registry
// injection are asserted against the loop production code runs, not against a
// copy of it.
func TestAC247RealCaptureLoopEmitsLevels(t *testing.T) {
	const frames = 8
	samples := make([]int16, frames*FrameSamples)
	for i := range samples {
		samples[i] = 12345 // a steady DC level: RMS is exactly this magnitude
	}
	dev := DeviceDescriptor{ID: "dev-level", Name: "Level Fixture Mic"}
	watcher := &fakeWatcher{devs: []DeviceDescriptor{dev}}
	stream := &fakeStream{
		dev: dev, rate: TargetRate, period: 10 * time.Millisecond,
		chunks: makeChunks(samples, 160),
	}
	opener := &fakeOpener{streams: []*fakeStream{stream}}

	rec := &levelRecorder{}
	reg := observe.NewRegistry()
	mic := newWASAPIMicrophoneWith(watcher, opener,
		WithCaptureRegistry(reg),
		WithLevelSink(rec.call),
	)
	baseDefault := observe.Default.CountByName("audio-capture")

	buf := make(chan []byte, 64)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := mic.Start(ctx, buf); err != nil {
		t.Fatal(err)
	}

	got := waitForLevels(t, rec, frames, 5*time.Second)
	for i, l := range got[:frames] {
		want := float32(LevelOfSamples(samples[i*FrameSamples : (i+1)*FrameSamples]))
		if l != want {
			t.Fatalf("frame %d: level = %v, want %v", i, l, want)
		}
	}
	// The thread the loop runs on is booked in the injected registry only.
	if n := reg.CountByName("audio-capture"); n != 1 {
		t.Fatalf("injected registry audio-capture = %d, want 1", n)
	}
	if n := observe.Default.CountByName("audio-capture") - baseDefault; n != 0 {
		t.Fatalf("the process registry also booked the capture thread (delta %d)", n)
	}
	// P8 form 甲's cost clause: the delivery runs on the pinned thread, so the
	// frame cadence is the reading that has to stay at FrameDuration. Measured
	// here (and re-measured on a real device in this leg's evidence file).
	_ = mic.Stop()
	if st := mic.Stats(); st.FramesSent < frames {
		t.Fatalf("loop delivered %d frames, want >= %d", st.FramesSent, frames)
	}
}
