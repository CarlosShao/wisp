//go:build windows

package audio

import (
	"context"
	"math"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// The hotplug tests inject fake devices at the ENUMERATION seam
// (DeviceWatcher + streamOpener) while the REAL capture loop, resampler and
// bounded push run - the audio-data seam is never mocked (SPEC-04 sec 9).
// Stream data comes from generated wav content through the real pipeline.

// fakeWatcher enumerates scripted devices and fires change events.
type fakeWatcher struct {
	mu     sync.Mutex
	devs   []DeviceDescriptor // scripted default capture devices, in order
	call   int
	failAt int // 1-based enumeration call that fails (0 = never)
	cb     func()
}

func (f *fakeWatcher) DefaultCaptureDevice() (DeviceDescriptor, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.call++
	if f.failAt > 0 && f.call >= f.failAt {
		return DeviceDescriptor{}, observe.New(observe.ClassAudioDevice, "no capture endpoint present")
	}
	idx := f.call - 1
	if idx >= len(f.devs) {
		idx = len(f.devs) - 1
	}
	return f.devs[idx], nil
}

func (f *fakeWatcher) DefaultRenderDevice() (DeviceDescriptor, error) {
	return DeviceDescriptor{ID: "render-fake", Name: "Fake Speakers"}, nil
}

func (f *fakeWatcher) OnDefaultCaptureChanged(cb func()) (func(), error) {
	f.mu.Lock()
	f.cb = cb
	f.mu.Unlock()
	return func() {
		f.mu.Lock()
		f.cb = nil
		f.mu.Unlock()
	}, nil
}

func (f *fakeWatcher) fireChange() {
	f.mu.Lock()
	cb := f.cb
	f.mu.Unlock()
	if cb != nil {
		cb()
	}
}

func (f *fakeWatcher) Close() error { return nil }

// fakeStream replays mono chunks at `rate` in `period` intervals, looping.
// failAfter > 0 makes WaitEvent return a device-invalidated error once that
// many periods elapsed (the D42#2 stale-handle scenario).
type fakeStream struct {
	dev       DeviceDescriptor
	rate      int
	chunks    [][]int16
	period    time.Duration
	failAfter int

	pos    int
	fails  int
	mu     sync.Mutex
	closed atomic.Bool
}

func (s *fakeStream) Descriptor() DeviceDescriptor { return s.dev }
func (s *fakeStream) Rate() int                    { return s.rate }

func (s *fakeStream) WaitEvent(timeoutMs uint32) (bool, error) {
	if s.closed.Load() {
		return false, DeviceError(hrAUDCLNTDeviceInv, s.dev, "WaitEvent")
	}
	if s.failAfter > 0 {
		s.mu.Lock()
		s.fails++
		fail := s.fails > s.failAfter
		s.mu.Unlock()
		if fail {
			return false, DeviceError(hrAUDCLNTDeviceInv, s.dev, "stream invalidated mid-capture")
		}
	}
	time.Sleep(s.period)
	return true, nil
}

func (s *fakeStream) Drain() ([]int16, error) {
	if len(s.chunks) == 0 {
		return make([]int16, 320), nil // 20ms of digital silence @16k
	}
	s.mu.Lock()
	chunk := s.chunks[s.pos%len(s.chunks)]
	s.pos++
	s.mu.Unlock()
	out := make([]int16, len(chunk))
	copy(out, chunk)
	return out, nil
}

func (s *fakeStream) Close() error {
	s.closed.Store(true)
	return nil
}

// fakeOpener hands out scripted streams / errors per open call.
type fakeOpener struct {
	mu      sync.Mutex
	streams []*fakeStream
	errAt   map[int]error // 1-based open call that must fail
	opens   []DeviceDescriptor
}

func (o *fakeOpener) Open(dev DeviceDescriptor) (deviceStream, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.opens = append(o.opens, dev)
	n := len(o.opens)
	if err, ok := o.errAt[n]; ok {
		return nil, err
	}
	if len(o.streams) == 0 {
		// Endless silence keeps the loop alive without extra fixtures.
		return &fakeStream{dev: dev, rate: TargetRate, period: 20 * time.Millisecond}, nil
	}
	s := o.streams[0]
	o.streams = o.streams[1:]
	return s, nil
}

func (o *fakeOpener) openCount() int {
	o.mu.Lock()
	defer o.mu.Unlock()
	return len(o.opens)
}

// makeChunks turns a mono sample slice into fakeStream device periods.
func makeChunks(samples []int16, periodSamples int) [][]int16 {
	var chunks [][]int16
	for off := 0; off < len(samples); off += periodSamples {
		end := off + periodSamples
		if end > len(samples) {
			end = len(samples)
		}
		chunks = append(chunks, samples[off:end])
	}
	return chunks
}

func toneSamples(n int, freq float64, amp float64) []int16 {
	out := make([]int16, n)
	for i := range out {
		out[i] = int16(amp * 32767 * float64(0.5+0.5*sin2pi(freq*float64(i)/float64(TargetRate))))
	}
	return out
}

func sin2pi(x float64) float64 { return math.Sin(2 * math.Pi * x) }

// drainFrames receives n frames or fails.
func drainFrames(t *testing.T, buf chan []byte, n int, within time.Duration) [][]byte {
	t.Helper()
	var got [][]byte
	deadline := time.After(within)
	for len(got) < n {
		select {
		case f := <-buf:
			got = append(got, f)
		case <-deadline:
			t.Fatalf("only %d/%d frames within %v", len(got), n, within)
		}
	}
	return got
}

// TestHotplugReenumerateOnce: acceptance criterion 4 - an injected default-
// device change triggers exactly ONE re-enumeration + reopen, the new
// endpoint's data flows through the real pipeline, and the counter shows it.
func TestHotplugReenumerateOnce(t *testing.T) {
	dev1 := DeviceDescriptor{ID: "dev-1", Name: "Fake Mic A"}
	dev2 := DeviceDescriptor{ID: "dev-2", Name: "Fake Mic B"}
	watcher := &fakeWatcher{devs: []DeviceDescriptor{dev1, dev2}}

	// Stream 1: 48kHz device rate (exercises the resampler), stream 2: 16k.
	stream1 := &fakeStream{
		dev: dev1, rate: 48000, period: 10 * time.Millisecond,
		chunks: makeChunks(toneSamples(4800, 440, 0.8), 480), // 100ms @48k
	}
	stream2 := &fakeStream{
		dev: dev2, rate: TargetRate, period: 20 * time.Millisecond,
		chunks: makeChunks(toneSamples(3200, 880, 0.8), 320), // 100ms @16k
	}
	opener := &fakeOpener{streams: []*fakeStream{stream1, stream2}}

	mic := newWASAPIMicrophoneWith(watcher, opener)
	buf := make(chan []byte, 64)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := mic.Start(ctx, buf); err != nil {
		t.Fatal(err)
	}
	if got := opener.openCount(); got != 1 {
		t.Fatalf("initial opens = %d, want 1", got)
	}
	drainFrames(t, buf, 2, 3*time.Second)

	// Inject OnDefaultDeviceChanged (D42#2).
	watcher.fireChange()
	drainFrames(t, buf, 2, 3*time.Second) // frames keep flowing after the swap

	if got := opener.openCount(); got != 2 {
		t.Fatalf("opens after change = %d, want exactly 2 (re-enumerate ONCE)", got)
	}
	opener.mu.Lock()
	secondDev := opener.opens[1]
	opener.mu.Unlock()
	if secondDev != dev2 {
		t.Fatalf("reopen used %v, want the NEW default %v", secondDev, dev2)
	}
	if st := mic.Stats(); st.Reopens != 1 {
		t.Fatalf("Reopens = %d, want 1", st.Reopens)
	}
	// fire the change again: same default -> another reopen (still once per event)
	watcher.fireChange()
	time.Sleep(100 * time.Millisecond)
	if st := mic.Stats(); st.Reopens != 2 || opener.openCount() != 3 {
		t.Fatalf("second change not handled once: reopens=%d opens=%d", st.Reopens, opener.openCount())
	}
	if err := mic.Stop(); err != nil {
		t.Fatal(err)
	}
	if err := mic.Err(); err != nil {
		t.Fatalf("unexpected terminal error: %v", err)
	}
}

// TestHotplugStaleHandleFailsLoudly: acceptance criterion 4 second half -
// when the reopen fails, the supervisor surfaces Error(audio_device) with
// the device name instead of silently keeping the stale handle.
func TestHotplugStaleHandleFailsLoudly(t *testing.T) {
	dev1 := DeviceDescriptor{ID: "dev-stale", Name: "Dying Mic"}
	watcher := &fakeWatcher{devs: []DeviceDescriptor{dev1}}

	stream1 := &fakeStream{
		dev: dev1, rate: TargetRate, period: 10 * time.Millisecond,
		chunks:    makeChunks(toneSamples(1600, 440, 0.8), 160),
		failAfter: 8, // device invalidated after 8 periods (2 frames emitted first)
	}
	opener := &fakeOpener{
		streams: []*fakeStream{stream1},
		errAt: map[int]error{
			2: DeviceError(hrAUDCLNTDeviceInUse, dev1, "Open capture"), // reopen occupied
		},
	}
	mic := newWASAPIMicrophoneWith(watcher, opener)
	buf := make(chan []byte, 32)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := mic.Start(ctx, buf); err != nil {
		t.Fatal(err)
	}
	drainFrames(t, buf, 1, 3*time.Second)

	// Wait for the stale-handle + failed reopen cycle.
	deadline := time.After(5 * time.Second)
	for {
		if err := mic.Err(); err != nil {
			if got, ok := observe.ClassOf(err); !ok || got != observe.ClassAudioDevice {
				t.Fatalf("terminal error class = %v (ok=%v), want audio_device", got, ok)
			}
			msg := err.Error()
			if !strings.Contains(msg, "Dying Mic") {
				t.Fatalf("error must name the device: %s", msg)
			}
			if !strings.Contains(msg, inUseGuidance) {
				t.Fatalf("error must carry the in-use guidance: %s", msg)
			}
			break
		}
		select {
		case <-deadline:
			t.Fatal("no terminal error within 5s")
		case <-time.After(10 * time.Millisecond):
		}
	}
	if st := mic.Stats(); st.LastError == "" {
		t.Fatal("stats.LastError must record the failure (not silent)")
	}
	if opener.openCount() != 2 {
		t.Fatalf("opens = %d, want exactly 2 (one reopen attempt)", opener.openCount())
	}
	// Frames must have stopped: drain what was buffered pre-failure, then
	// the channel must stay empty (the loop is done - no silent stale flow).
drainAll:
	for {
		select {
		case <-buf:
		default:
			break drainAll
		}
	}
	select {
	case f := <-buf:
		t.Fatalf("frames still flowing after terminal failure: %d bytes", len(f))
	case <-time.After(300 * time.Millisecond):
	}
	if err := mic.Stop(); err != nil {
		t.Fatal(err)
	}
}

// TestHotplugReenumerateFailureNamesDevice: enumeration itself failing on
// reopen is terminal and names the device we refuse to keep silently.
func TestHotplugReenumerateFailureNamesDevice(t *testing.T) {
	dev1 := DeviceDescriptor{ID: "dev-gone", Name: "Unplugged Mic"}
	watcher := &fakeWatcher{devs: []DeviceDescriptor{dev1}, failAt: 2} // 2nd enum fails

	stream1 := &fakeStream{
		dev: dev1, rate: TargetRate, period: 10 * time.Millisecond,
		chunks:    makeChunks(toneSamples(1600, 440, 0.8), 160),
		failAfter: 6, // device invalidated after 6 periods (1 frame emitted first)
	}
	opener := &fakeOpener{streams: []*fakeStream{stream1}}
	mic := newWASAPIMicrophoneWith(watcher, opener)
	buf := make(chan []byte, 32)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := mic.Start(ctx, buf); err != nil {
		t.Fatal(err)
	}
	drainFrames(t, buf, 1, 3*time.Second)

	deadline := time.After(5 * time.Second)
	for {
		if err := mic.Err(); err != nil {
			msg := err.Error()
			if !strings.Contains(msg, "Unplugged Mic") {
				t.Fatalf("re-enumerate failure must name the (previous) device: %s", msg)
			}
			if got, ok := observe.ClassOf(err); !ok || got != observe.ClassAudioDevice {
				t.Fatalf("class = %v (ok=%v)", got, ok)
			}
			break
		}
		select {
		case <-deadline:
			t.Fatal("no terminal error within 5s")
		case <-time.After(10 * time.Millisecond):
		}
	}
	if err := mic.Stop(); err != nil {
		t.Fatal(err)
	}
}

// TestOpenOccupiedAndPermissionDenied: acceptance criterion 5 - the open
// failure fixtures map to Error(audio_device) with the right guidance text.
func TestOpenOccupiedAndPermissionDenied(t *testing.T) {
	dev := DeviceDescriptor{ID: "dev-locked", Name: "Busy Mic"}

	cases := []struct {
		name     string
		err      error
		contains string
	}{
		{
			name:     "occupied (exclusive mode)",
			err:      DeviceError(hrAUDCLNTDeviceInUse, dev, "Initialize(shared mode)"),
			contains: inUseGuidance,
		},
		{
			name:     "permission denied",
			err:      DeviceError(hrEAccessDenied, dev, "Initialize(shared mode)"),
			contains: privacyGuidance,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			watcher := &fakeWatcher{devs: []DeviceDescriptor{dev}}
			opener := &fakeOpener{errAt: map[int]error{1: tc.err}}
			mic := newWASAPIMicrophoneWith(watcher, opener)
			buf := make(chan []byte, 8)
			err := mic.Start(context.Background(), buf)
			if err == nil {
				_ = mic.Stop()
				t.Fatal("Start must fail when the device cannot be opened (D42#12)")
			}
			if got, ok := observe.ClassOf(err); !ok || got != observe.ClassAudioDevice {
				t.Fatalf("class = %v (ok=%v), want audio_device", got, ok)
			}
			msg := err.Error()
			if !strings.Contains(msg, "Busy Mic") {
				t.Fatalf("error must name the device: %s", msg)
			}
			if !strings.Contains(msg, tc.contains) {
				t.Fatalf("error must contain guidance %q: %s", tc.contains, msg)
			}
			if st := mic.Stats(); st.LastError == "" {
				t.Fatal("stats.LastError must record the open failure")
			}
		})
	}
}

// TestEndpointsPairQueryable pins the D47 note: the capture/render pair is
// queryable through the watcher for the P15 spike.
func TestEndpointsPairQueryable(t *testing.T) {
	watcher := &fakeWatcher{devs: []DeviceDescriptor{{ID: "c1", Name: "Mic X"}}}
	opener := &fakeOpener{}
	mic := newWASAPIMicrophoneWith(watcher, opener)
	cap1, render, err := mic.Endpoints()
	if err != nil {
		t.Fatal(err)
	}
	if cap1.Name != "Mic X" || render.Name != "Fake Speakers" {
		t.Fatalf("endpoints = %v / %v", cap1, render)
	}
}

// TestPinnedThreadStable10s: acceptance criterion 2 - the capture loop's OS
// thread is stable across 10s (no M-migration) and no other goroutine in the
// process lands on that thread (LockOSThread exclusivity, D38a).
func TestPinnedThreadStable10s(t *testing.T) {
	if testing.Short() {
		t.Skip("10s pinned-thread window (acceptance criterion) skipped in -short")
	}
	dev := DeviceDescriptor{ID: "dev-pin", Name: "Pinned Mic"}
	watcher := &fakeWatcher{devs: []DeviceDescriptor{dev}}
	stream := &fakeStream{
		dev: dev, rate: TargetRate, period: 20 * time.Millisecond,
		chunks: makeChunks(toneSamples(320, 440, 0.8), 320),
	}
	opener := &fakeOpener{streams: []*fakeStream{stream}}
	mic := newWASAPIMicrophoneWith(watcher, opener)
	buf := make(chan []byte, 16)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := mic.Start(ctx, buf); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = mic.Stop() }()

	pinned := mic.ThreadID()
	if pinned == 0 {
		t.Fatal("capture thread id not published")
	}

	// Other goroutines sample their own OS thread id; none may equal the
	// pinned capture thread while the loop lives.
	otherIDs := make(chan uint32, 64)
	stopSampling := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stopSampling:
					return
				default:
					select {
					case otherIDs <- currentThreadID():
					default:
					}
				}
				time.Sleep(20 * time.Millisecond)
			}
		}()
	}

	ids := map[uint32]bool{}
	start := time.Now()
	frames := 0
	for time.Since(start) < 10*time.Second {
		ids[mic.ThreadID()] = true
		select {
		case <-buf:
			frames++
		default:
		}
		select {
		case id := <-otherIDs:
			if id == pinned {
				t.Fatalf("another goroutine ran on the pinned capture thread %d", pinned)
			}
		default:
		}
		time.Sleep(50 * time.Millisecond)
	}
	close(stopSampling)
	wg.Wait()

	if len(ids) != 1 {
		t.Fatalf("capture thread migrated: distinct OS thread ids = %v", ids)
	}
	if frames == 0 {
		t.Fatal("no frames captured during the 10s window")
	}
}

// TestLiveWasapiSmoke is the real-machine hook (真机冒烟): gated behind
// WISP_LIVE_MIC=1, it exercises the actual WASAPI COM path end to end. CI
// has no microphone; the human smoke pass belongs to ticket 16 (S2).
func TestLiveWasapiSmoke(t *testing.T) {
	if os.Getenv("WISP_LIVE_MIC") != "1" {
		t.Skip("live WASAPI smoke requires WISP_LIVE_MIC=1 and a real microphone (真机冒烟待票 16)")
	}
	mic := NewWASAPIMicrophone()
	cap1, render, err := mic.Endpoints()
	if err != nil {
		t.Fatalf("enumerate endpoints: %v", err)
	}
	t.Logf("capture=%q render=%q", cap1.String(), render.String())

	buf := NewBoundedFrames()
	ctx, cancel := context.WithCancel(context.Background())
	if err := mic.Start(ctx, buf); err != nil {
		t.Fatalf("live Start: %v (device occupied? privacy settings?)", err)
	}
	// Sustained capture: ~1s of audio must arrive at the 32ms frame cadence.
	const wantFrames = 25
	frames := 0
	start := time.Now()
	deadline := time.After(2 * time.Second)
	for frames < wantFrames {
		select {
		case f := <-buf:
			if len(f) != FrameBytes {
				t.Fatalf("frame size %d, want %d", len(f), FrameBytes)
			}
			frames++
		case <-deadline:
			t.Fatalf("only %d frames in 2s (mic muted or silent?)", frames)
		}
	}
	if elapsed := time.Since(start); elapsed > 1500*time.Millisecond {
		t.Fatalf("%d frames took %v - capture cadence broken", wantFrames, elapsed)
	}
	cancel()
	if err := mic.Stop(); err != nil {
		t.Fatal(err)
	}
	if tid := mic.ThreadID(); tid != 0 {
		t.Fatal("thread id not reset after Stop")
	}
	t.Logf("live smoke OK: %d frames in %v, stats %+v", frames, time.Since(start).Round(time.Millisecond), mic.Stats())
}

// TestCaptureLoopWavIntegrity: the REAL capture loop (wait/drain/resample/
// frame/pump) fed by a real wav fixture parsed with the production parser -
// the frames coming out must match the injector's frame contract exactly.
func TestCaptureLoopWavIntegrity(t *testing.T) {
	samples := sineI16At(TargetRate, TargetRate, 440, 0.6) // 1s @16k mono
	p := t.TempDir() + "/fixture.wav"
	writeWav(t, p, samples, TargetRate, 1)
	raw, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	parsed, rate, chans, err := parseWav(raw)
	if err != nil || rate != TargetRate || chans != 1 {
		t.Fatalf("parseWav: %v rate=%d chans=%d", err, rate, chans)
	}

	dev := DeviceDescriptor{ID: "dev-wav", Name: "Wav Fixture Mic"}
	watcher := &fakeWatcher{devs: []DeviceDescriptor{dev}}
	stream := &fakeStream{
		dev: dev, rate: TargetRate, period: 10 * time.Millisecond,
		chunks: makeChunks(parsed, 160), // 10ms periods
	}
	opener := &fakeOpener{streams: []*fakeStream{stream}}
	mic := newWASAPIMicrophoneWith(watcher, opener)
	buf := make(chan []byte, 128)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := mic.Start(ctx, buf); err != nil {
		t.Fatal(err)
	}

	// The endless capture path emits only FULL frames - the zero-padded tail
	// is the finite-source (WavInjector) contract, not the microphone's.
	full := len(samples) / FrameSamples
	want := expectedFrames(samples)[:full]
	got := drainFrames(t, buf, full, 5*time.Second)
	_ = mic.Stop()
	for i := range want {
		if len(got[i]) != FrameBytes || string(got[i]) != string(want[i]) {
			t.Fatalf("frame %d differs (len %d)", i, len(got[i]))
		}
	}
	if st := mic.Stats(); st.FramesDropped != 0 {
		t.Fatalf("unexpected drops: %+v", st)
	}
}
