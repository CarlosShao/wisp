package audio

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Ticket 247 AC#1/AC#3/AC#5's deterministic half: the two capture-side seams
// the wiring leg added (which registry owns the audio-capture thread, and where
// the per-frame level goes) exercised through the C8 test seam, the wav
// injector, which runs the same bounded push and the same pump the real
// microphone does. These are NOT the ticket's AC#2 reading (AC#2 wants a real
// device); they are the proof that the plumbing carries exactly one scalar per
// frame and books exactly one registry.

type levelRecorder struct {
	mu     sync.Mutex
	levels []float32
}

func (r *levelRecorder) call(level float32) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.levels = append(r.levels, level)
}

func (r *levelRecorder) snapshot() []float32 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]float32(nil), r.levels...)
}

// waitForLevels polls the recorder until n arrivals or the deadline.
func waitForLevels(t *testing.T, r *levelRecorder, n int, within time.Duration) []float32 {
	t.Helper()
	deadline := time.Now().Add(within)
	for time.Now().Before(deadline) {
		if got := r.snapshot(); len(got) >= n {
			return got
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("only %d levels arrived within %v, want %d", len(r.snapshot()), within, n)
	return nil
}

// TestAC247LevelSinkGetsOneScalarPerSeamFrame: the injector's pump must hand
// the sink exactly one float32 per frame it pushes, and each one must be the
// level of that very frame (FrameLevel of the encoded frame), not an
// approximation of some other window.
func TestAC247LevelSinkGetsOneScalarPerSeamFrame(t *testing.T) {
	const frames = 12
	samples := make([]int16, frames*FrameSamples)
	for i := range samples {
		// A square wave at the strongest symmetric amplitude: its level is a
		// named constant of this package (FullScaleSquareLevel), so the test
		// pins the scale too, not just the count.
		samples[i] = 32767
		if i%2 == 0 {
			samples[i] = -32767
		}
	}
	path := t.TempDir() + "/levels.wav"
	writeWav(t, path, samples, TargetRate, 1)

	rec := &levelRecorder{}
	inj, err := NewWavInjector(path,
		WithInjectorPace(0), // flood: the level path must not depend on pacing
		WithInjectorCapture(WithLevelSink(rec.call)),
	)
	if err != nil {
		t.Fatal(err)
	}
	buf := make(chan []byte, 64) // wide enough that no frame is dropped here
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := inj.Start(ctx, buf); err != nil {
		t.Fatal(err)
	}

	got := waitForLevels(t, rec, frames, 5*time.Second)
	if len(got) != frames {
		t.Fatalf("levels = %d, want exactly %d frames", len(got), frames)
	}
	for i, l := range got {
		if l < MinLevel || l > MaxLevel {
			t.Fatalf("level %d = %v, outside the closed scale [%v,%v]", i, l, MinLevel, MaxLevel)
		}
		want := float32(LevelOfSamples(inj.samples[i*FrameSamples : (i+1)*FrameSamples]))
		if l != want {
			t.Fatalf("level %d = %v, want the same frame's FrameLevel %v", i, l, want)
		}
	}
	if last := got[0]; last != FullScaleSquareLevel {
		t.Fatalf("square wave at +-32767 reads %v, want FullScaleSquareLevel %v", last, FullScaleSquareLevel)
	}
	if err := inj.Stop(); err != nil {
		t.Fatal(err)
	}
}

// TestAC247LevelSinkRunsEvenWithNoFrameConsumer: the resident leg has no PCM
// consumer yet (internal/speech is unwritten and is not this ticket), so the
// bounded channel fills and every frame is dropped. The level must still flow:
// it is a fact about the microphone, not about how fast a reader drains.
func TestAC247LevelSinkRunsEvenWithNoFrameConsumer(t *testing.T) {
	samples := sineI16At(40*FrameSamples, TargetRate, 440, 0.5)
	path := t.TempDir() + "/unread.wav"
	writeWav(t, path, samples, TargetRate, 1)

	rec := &levelRecorder{}
	inj, err := NewWavInjector(path,
		WithInjectorPace(2*time.Millisecond),
		WithInjectorCapture(WithLevelSink(rec.call)),
	)
	if err != nil {
		t.Fatal(err)
	}
	buf := NewBoundedFrames() // capacity 6: nobody reads it, 34 frames must drop
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := inj.Start(ctx, buf); err != nil {
		t.Fatal(err)
	}
	want := 40
	got := waitForLevels(t, rec, want, 10*time.Second)
	if len(got) < want {
		t.Fatalf("levels = %d, want %d with an unread channel", len(got), want)
	}
	if st := inj.Stats(); st.FramesDropped == 0 {
		t.Fatalf("the channel was never full, so this test proved nothing: %+v", st)
	} else if st.FramesSent >= uint64(want) {
		t.Fatalf("a bounded channel with no reader delivered %d frames", st.FramesSent)
	}
	if err := inj.Stop(); err != nil {
		t.Fatal(err)
	}
}

// TestAC247CaptureRegistryReplacesTheDefault: the injected registry becomes the
// ONLY owner of the audio-capture booking. The process registry gains nothing
// for that thread - two registries counting the same goroutine is the false
// green shape this ruling (P4 form 甲) refuses, because a roster read off the
// wrong one would say "no capture thread is running" while it runs.
func TestAC247CaptureRegistryReplacesTheDefault(t *testing.T) {
	samples := sineI16At(6*FrameSamples, TargetRate, 440, 0.4)
	path := t.TempDir() + "/registry.wav"
	writeWav(t, path, samples, TargetRate, 1)

	reg := observe.NewRegistry()
	baseDefault := observe.Default.CountByName("audio-capture")

	inj, err := NewWavInjector(path,
		WithInjectorPace(20*time.Millisecond),
		WithInjectorCapture(WithCaptureRegistry(reg)),
	)
	if err != nil {
		t.Fatal(err)
	}
	buf := make(chan []byte, 64)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := inj.Start(ctx, buf); err != nil {
		t.Fatal(err)
	}
	if got := reg.CountByName("audio-capture"); got != 1 {
		t.Fatalf("injected registry counts %d audio-capture, want 1", got)
	}
	if got := observe.Default.CountByName("audio-capture") - baseDefault; got != 0 {
		t.Fatalf("the process registry also booked the thread (delta %d): two registries for one goroutine", got)
	}
	rep := reg.RosterReport()
	if rep.Resident != 1 || rep.ResidentOverBaseline || len(rep.Unknown) != 0 {
		t.Fatalf("roster report = %+v, want one resident name and nothing unknown", rep)
	}
	if err := inj.Stop(); err != nil {
		t.Fatal(err)
	}
	if got := reg.CountByName("audio-capture"); got != 0 {
		t.Fatalf("audio-capture still counted after Stop: %d", got)
	}
}

// TestAC247CaptureConfigDefaultsKeepTheOldShape: with no options the sources
// behave exactly as before this seam existed (the process registry, no sink),
// and a nil registry cannot clear it. That is what makes the injection additive
// rather than a semantic change to a symbol ticket 13/241 already delivered.
func TestAC247CaptureConfigDefaultsKeepTheOldShape(t *testing.T) {
	c := newCaptureConfig()
	if c.Registry != observe.Default {
		t.Fatal("the no-option default is not the process registry")
	}
	if c.OnLevel != nil {
		t.Fatal("a source built with no level sink must emit nothing, not guess a destination")
	}
	c2 := newCaptureConfig(WithCaptureRegistry(nil))
	if c2.Registry != observe.Default {
		t.Fatal("a nil registry cleared the process default (a thread with no roster)")
	}
	// emitLevel with no sink must be a no-op and must not panic.
	newCaptureConfig().emitLevel(make([]byte, FrameBytes))
	// A frame that is not one seam frame is refused out loud and emits nothing.
	rec := &levelRecorder{}
	newCaptureConfig(WithLevelSink(rec.call)).emitLevel(make([]byte, FrameBytes-2))
	if len(rec.snapshot()) != 0 {
		t.Fatal("a half frame reached the sink: it would read as a different loudness")
	}
}
