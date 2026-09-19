package audio

import (
	"bytes"
	"context"
	"encoding/binary"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// writeWav encodes PCM16 mono samples as a RIFF/WAVE file (test fixture).
func writeWav(t *testing.T, path string, samples []int16, rate int, chans int) {
	t.Helper()
	dataLen := len(samples) * 2
	var buf bytes.Buffer
	buf.WriteString("RIFF")
	buf.Write(u32le(uint32(36 + dataLen)))
	buf.WriteString("WAVE")
	buf.WriteString("fmt ")
	buf.Write(u32le(16))
	buf.Write(u16le(1)) // PCM
	buf.Write(u16le(uint16(chans)))
	buf.Write(u32le(uint32(rate)))
	buf.Write(u32le(uint32(rate * chans * 2)))
	buf.Write(u16le(uint16(chans * 2)))
	buf.Write(u16le(16))
	buf.WriteString("data")
	buf.Write(u32le(uint32(dataLen)))
	for _, s := range samples {
		buf.Write(u16le(uint16(s)))
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
}

func u16le(v uint16) []byte { var b [2]byte; binary.LittleEndian.PutUint16(b[:], v); return b[:] }
func u32le(v uint32) []byte { var b [4]byte; binary.LittleEndian.PutUint32(b[:], v); return b[:] }

func sineI16At(n, rate int, freq, amp float64) []int16 {
	out := make([]int16, n)
	for i := range out {
		out[i] = int16(math.Round(amp * 32767 * math.Sin(2*math.Pi*freq*float64(i)/float64(rate))))
	}
	return out
}

// expectedFrames chunks the 16k sample stream into seam frames, zero-padding
// the tail (the injector's documented frame contract).
func expectedFrames(samples []int16) [][]byte {
	var frames [][]byte
	for off := 0; off < len(samples); off += FrameSamples {
		end := off + FrameSamples
		pad := 0
		if end > len(samples) {
			pad = end - len(samples)
			end = len(samples)
		}
		b := make([]byte, 0, FrameBytes)
		for _, s := range samples[off:end] {
			b = append(b, byte(s), byte(s>>8))
		}
		b = append(b, make([]byte, pad*2)...)
		frames = append(frames, b)
	}
	return frames
}

// TestWavInjectorFrameExact: acceptance criterion 1 - WavInjector -> channel
// delivers exactly the wav content as 16k/mono/int16 frames.
func TestWavInjectorFrameExact(t *testing.T) {
	samples := sineI16At(1200, TargetRate, 440, 0.5) // 75ms: 2 full frames + tail
	dir := t.TempDir()
	p := filepath.Join(dir, "tone16k.wav")
	writeWav(t, p, samples, TargetRate, 1)

	inj, err := NewWavInjector(p)
	if err != nil {
		t.Fatal(err)
	}
	buf := make(chan []byte, 64)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := inj.Start(ctx, buf); err != nil {
		t.Fatal(err)
	}

	var got [][]byte
	for range len(samples)/FrameSamples + 1 {
		select {
		case f := <-buf:
			got = append(got, f)
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for frames")
		}
	}
	if err := inj.Stop(); err != nil {
		t.Fatal(err)
	}
	if len(got) != len(samples)/FrameSamples+1 {
		t.Fatalf("got %d frames", len(got))
	}
	want := expectedFrames(samples)
	for i := range want {
		if !bytes.Equal(got[i], want[i]) {
			t.Fatalf("frame %d differs (len %d vs %d)", i, len(got[i]), len(want[i]))
		}
	}
	st := inj.Stats()
	if st.FramesSent != uint64(len(got)) || st.FramesDropped != 0 {
		t.Fatalf("stats: %+v", st)
	}
}

// TestWavInjectorRateConversion: a 48kHz stereo wav is downmixed and
// resampled to the 16k mono seam rate inside the injector.
func TestWavInjectorRateConversion(t *testing.T) {
	l := sineI16At(4800, 48000, 1000, 0.4)
	r := sineI16At(4800, 48000, 1000, 0.2)
	stereo := make([]int16, 0, len(l)*2)
	for i := range l {
		stereo = append(stereo, l[i], r[i])
	}
	p := filepath.Join(t.TempDir(), "tone48k.wav")
	writeWav(t, p, stereo, 48000, 2)

	inj, err := NewWavInjector(p)
	if err != nil {
		t.Fatal(err)
	}
	want := ResampleLinear(MonoDownmix(stereo, 2), 48000, TargetRate)
	if len(inj.samples) != len(want) {
		t.Fatalf("injector holds %d samples, want %d", len(inj.samples), len(want))
	}
	for i := range want {
		if inj.samples[i] != want[i] {
			t.Fatalf("sample %d differs", i)
		}
	}
}

// TestWavInjectorBackpressureDropCounted: acceptance criterion 1 second half
// - a slow consumer makes the bounded channel overflow, frames are dropped
// AND counted, and the counter is visible in the log.
func TestWavInjectorBackpressureDropCounted(t *testing.T) {
	var logBuf syncBuffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&logBuf, nil)))
	defer slog.SetDefault(old)

	samples := sineI16At(TargetRate*2, TargetRate, 440, 0.5) // 2s of audio
	p := filepath.Join(t.TempDir(), "long.wav")
	writeWav(t, p, samples, TargetRate, 1)
	inj, err := NewWavInjector(p)
	if err != nil {
		t.Fatal(err)
	}

	buf := NewBoundedFrames() // 6 frames = 192ms
	ctx, cancel := context.WithCancel(context.Background())
	if err := inj.Start(ctx, buf); err != nil {
		t.Fatal(err)
	}

	// Slow consumer: 1 frame per 200ms, then give up after 1s.
	deadline := time.After(1200 * time.Millisecond)
	received := 0
recv:
	for {
		select {
		case <-buf:
			received++
			time.Sleep(200 * time.Millisecond)
		case <-deadline:
			break recv
		}
	}
	cancel()
	if err := inj.Stop(); err != nil {
		t.Fatal(err)
	}
	// Frames delivered into the bounded buffer before the deadline but not
	// yet received still count as sent - drain them.
drain:
	for {
		select {
		case <-buf:
			received++
		default:
			break drain
		}
	}
	st := inj.Stats()
	if st.FramesDropped == 0 {
		t.Fatal("slow consumer must produce drops, got none")
	}
	if st.FramesSent != uint64(received) {
		t.Fatalf("FramesSent %d != received %d", st.FramesSent, received)
	}
	if st.BytesDropped != st.FramesDropped*FrameBytes {
		t.Fatalf("byte counter inconsistent: %+v", st)
	}
	// Counter visible in logs (D38d: silent drops are forbidden).
	if !strings.Contains(logBuf.String(), "audio frames dropped") {
		t.Fatalf("drop warning missing from log:\n%s", logBuf.String())
	}
	if !strings.Contains(logBuf.String(), "dropped_frames_total=") {
		t.Fatalf("drop total missing from log:\n%s", logBuf.String())
	}
}

// TestWavInjectorCtxCancel: ctx cancellation stops the source promptly.
func TestWavInjectorCtxCancel(t *testing.T) {
	samples := sineI16At(TargetRate*10, TargetRate, 440, 0.5) // 10s
	p := filepath.Join(t.TempDir(), "verylong.wav")
	writeWav(t, p, samples, TargetRate, 1)
	inj, err := NewWavInjector(p)
	if err != nil {
		t.Fatal(err)
	}
	buf := make(chan []byte, 4)
	ctx, cancel := context.WithCancel(context.Background())
	if err := inj.Start(ctx, buf); err != nil {
		t.Fatal(err)
	}
	time.Sleep(150 * time.Millisecond)
	cancel()
	if err := inj.Stop(); err != nil {
		t.Fatal(err)
	}
	if inj.Stats().FramesSent == 0 {
		t.Fatal("expected some frames before cancel")
	}
}

// TestWavInjectorBadFile pins parse errors for non-wav / unsupported formats.
func TestWavInjectorBadFile(t *testing.T) {
	if _, err := NewWavInjector(filepath.Join(t.TempDir(), "missing.wav")); err == nil {
		t.Fatal("missing file must error")
	}
	p := filepath.Join(t.TempDir(), "garbage.wav")
	if err := os.WriteFile(p, []byte("not a wav at all"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := NewWavInjector(p); err == nil {
		t.Fatal("garbage must error")
	}
	// 8-bit PCM is out of contract: reject loudly, never misinterpret.
	p2 := filepath.Join(t.TempDir(), "bits8.wav")
	bad := append([]byte("RIFF"), u32le(50)...)
	bad = append(bad, []byte("WAVEfmt ")...)
	bad = append(bad, u32le(16)...)
	bad = append(bad, u16le(1)...)
	bad = append(bad, u16le(1)...)
	bad = append(bad, u32le(8000)...)
	bad = append(bad, u32le(8000)...)
	bad = append(bad, u16le(1)...)
	bad = append(bad, u16le(8)...)
	bad = append(bad, []byte("data")...)
	bad = append(bad, u32le(4)...)
	bad = append(bad, 1, 2, 3, 4)
	if err := os.WriteFile(p2, bad, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := NewWavInjector(p2); err == nil {
		t.Fatal("8-bit wav must error")
	}
}

// TestBoundedWindowInvariants pins the D38d sizing rules.
func TestBoundedWindowInvariants(t *testing.T) {
	if FrameDuration != 32*time.Millisecond {
		t.Fatalf("FrameDuration = %v", FrameDuration)
	}
	if FrameBytes != 1024 {
		t.Fatalf("FrameBytes = %d", FrameBytes)
	}
	if BoundedWindow > 200*time.Millisecond {
		t.Fatalf("bounded window %v exceeds the 200ms cap (D38d)", BoundedWindow)
	}
	if BoundedFrameCapacity() != int(BoundedWindow/FrameDuration) {
		t.Fatalf("capacity inconsistent with window")
	}
	if got := cap(NewBoundedFrames()); got != BoundedFrameCapacity() {
		t.Fatalf("NewBoundedFrames cap = %d", got)
	}
}

type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}
