package audio

import (
	"context"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func gateWav(t *testing.T, seconds int) *WavInjector {
	t.Helper()
	samples := sineI16At(TargetRate*seconds, TargetRate, 440, 0.5)
	p := filepath.Join(t.TempDir(), "gate.wav")
	writeWav(t, p, samples, TargetRate, 1)
	inj, err := NewWavInjector(p)
	if err != nil {
		t.Fatal(err)
	}
	return inj
}

// waitForFrame receives one frame or fails after a bounded wait.
func waitForFrame(t *testing.T, buf chan []byte) {
	t.Helper()
	select {
	case <-buf:
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for a frame")
	}
}

// assertNoFrames proves the capture side is actually closed: no frame within
// 400ms (>= 12 frame periods at 32ms).
func assertNoFrames(t *testing.T, buf chan []byte) {
	t.Helper()
	select {
	case f := <-buf:
		t.Fatalf("frame delivered while gate closed: %d bytes", len(f))
	case <-time.After(400 * time.Millisecond):
	}
}

// TestGateHalfDuplexClosedZeroFrames: acceptance criterion 6 - gate closed
// (TTS speaking) means ZERO frames delivered, and the inner capture is
// actually stopped (device off, D16), not merely discarded; reopen restores.
func TestGateHalfDuplexClosedZeroFrames(t *testing.T) {
	inner := gateWav(t, 30) // 30s: long enough to outlive the test
	gate := NewHalfDuplexGate(inner, PathT)
	buf := make(chan []byte, 64)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := gate.Start(ctx, buf); err != nil {
		t.Fatal(err)
	}
	waitForFrame(t, buf)

	gate.SetSpeaking(true)
	assertNoFrames(t, buf)
	// The D16 contract is "capture OFF": the inner source must have stopped
	// producing (its FramesSent counter frozen), not a discard-only gate.
	frozen := inner.Stats().FramesSent
	time.Sleep(150 * time.Millisecond)
	if got := inner.Stats().FramesSent; got != frozen {
		t.Fatalf("inner capture still running while speaking: %d -> %d", frozen, got)
	}

	gate.SetSpeaking(false)
	waitForFrame(t, buf)
	if st := gate.Stats(); st.FramesSent <= frozen {
		t.Fatalf("no frames after reopen (stats %+v)", st)
	}
	if err := gate.Stop(); err != nil {
		t.Fatal(err)
	}
	if frozen2 := inner.Stats().FramesSent; inner.Stats().FramesSent != frozen2 {
		t.Fatal("inner still running after gate Stop")
	}
}

// TestGatePathCFullDuplex: the D47 exception - on Path C the gate must NOT
// close capture while speaking (AEC is tickets 26/59); events still fire.
func TestGatePathCFullDuplex(t *testing.T) {
	inner := gateWav(t, 10)
	var evs []GateEvent
	gate := NewHalfDuplexGate(inner, PathC, WithGateEvents(func(e GateEvent) {
		evs = append(evs, e)
	}))
	buf := make(chan []byte, 64)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := gate.Start(ctx, buf); err != nil {
		t.Fatal(err)
	}
	waitForFrame(t, buf)
	gate.SetSpeaking(true)
	waitForFrame(t, buf) // capture keeps running while speaking
	waitForFrame(t, buf)
	gate.SetSpeaking(false)
	if st := inner.Stats(); st.FramesDropped != 0 {
		t.Fatalf("Path C gate must not disturb capture: %+v", st)
	}
	// Events observed even though capture stayed open.
	var sawSpeaking, sawStopped bool
	for _, e := range evs {
		if e.Kind == EventSpeakingStarted {
			sawSpeaking = true
		}
		if e.Kind == EventSpeakingStopped {
			sawStopped = true
		}
	}
	if !sawSpeaking || !sawStopped {
		t.Fatalf("Path C speaking events missing: %+v", evs)
	}
	if err := gate.Stop(); err != nil {
		t.Fatal(err)
	}
}

// TestGateMutedEvents: the mute hotkey path (state 07 Muted) - muted means
// zero frames on BOTH paths, and Muted state events fire for the machine.
func TestGateMutedEvents(t *testing.T) {
	inner := gateWav(t, 10)
	var mu sync.Mutex
	var evs []GateEvent
	gate := NewHalfDuplexGate(inner, PathT, WithGateEvents(func(e GateEvent) {
		mu.Lock()
		evs = append(evs, e)
		mu.Unlock()
	}))
	buf := make(chan []byte, 64)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := gate.Start(ctx, buf); err != nil {
		t.Fatal(err)
	}
	waitForFrame(t, buf)

	gate.SetMuted(true)
	assertNoFrames(t, buf)
	// Muted wins even if TTS state flaps underneath.
	gate.SetSpeaking(true)
	gate.SetSpeaking(false)
	assertNoFrames(t, buf)

	gate.SetMuted(false)
	waitForFrame(t, buf)

	mu.Lock()
	defer mu.Unlock()
	var mutedN, unmutedN int
	for _, e := range evs {
		switch e.Kind {
		case EventMuted:
			mutedN++
		case EventUnmuted:
			unmutedN++
		}
	}
	if mutedN != 1 || unmutedN != 1 {
		t.Fatalf("want exactly one muted/unmuted event pair, got %d/%d: %+v", mutedN, unmutedN, evs)
	}
	if err := gate.Stop(); err != nil {
		t.Fatal(err)
	}
}

// TestGateStartMuted: `[audio] mic_muted_default` respected on boot - a gate
// started muted captures nothing until the user unmutes.
func TestGateStartMuted(t *testing.T) {
	inner := gateWav(t, 5)
	gate := NewHalfDuplexGate(inner, PathT, WithStartMuted(true))
	buf := make(chan []byte, 8)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := gate.Start(ctx, buf); err != nil {
		t.Fatal(err)
	}
	assertNoFrames(t, buf)
	if gate.Muted() != true || gate.Open() != false {
		t.Fatalf("gate state wrong: muted=%v open=%v", gate.Muted(), gate.Open())
	}
	gate.SetMuted(false)
	waitForFrame(t, buf)
	if err := gate.Stop(); err != nil {
		t.Fatal(err)
	}
}

// TestGateStopIdempotent: Stop without Start, double Stop, Stop after ctx
// cancel - all safe (D38e step 4 must never hang shutdown).
func TestGateStopIdempotent(t *testing.T) {
	inner := gateWav(t, 2)
	gate := NewHalfDuplexGate(inner, PathT)
	if err := gate.Stop(); err != nil { // never started
		t.Fatal(err)
	}
	buf := make(chan []byte, 4)
	ctx, cancel := context.WithCancel(context.Background())
	if err := gate.Start(ctx, buf); err != nil {
		t.Fatal(err)
	}
	gate.SetSpeaking(true) // close while stopping
	cancel()
	if err := gate.Stop(); err != nil {
		t.Fatal(err)
	}
	if err := gate.Stop(); err != nil { // double stop
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		_ = gate.Stop()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Stop blocked past 2s")
	}
}
