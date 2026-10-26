package speech

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
)

// These tests drive the real silero model through the real C8 seam
// (audio.WavInjector -> EncodeFrame -> DecodeFrame), on synthetic fixtures
// only. The exact boundary numbers live in endpoint_test.go, where they can be
// probed to the sample; here the question is whether the model is actually
// wired to the state machine and whether a pause is split or held.
//
// The fixtures are calibrated by measurement, not by arithmetic: silero keeps
// reporting speech for a while after a synthetic surrogate stops (measured
// ~720ms of hangover for this one), so a nominal pause is NOT the pause the
// state machine sees. Each assertion below carries the measured margin.

// newTestDetector loads silero or skips with the reason.
//
// The skip is not a convenience: on a host whose C:\Windows\System32 holds an
// OS-provided onnxruntime.dll, the Windows loader resolves the sherpa import
// to that copy instead of the one on PATH, and creating a session there
// segfaults the process (the C API gets a NULL OrtApi and dereferences it).
// NewDetector refuses that runtime up front, so the suite reports a named skip
// instead of dying. Run the staged binary (see the receipt) for the full run.
func newTestDetector(t *testing.T) *Detector {
	t.Helper()
	d, err := NewDetector(DetectorConfig{ModelPath: sileroModelPath(t)})
	if err != nil {
		t.Skipf("silero VAD unavailable in this process: %v", err)
	}
	t.Cleanup(func() { _ = d.Close() })
	return d
}

// driveSeam replays pcm through a WAV file and the C8 injector, exactly as the
// product's capture side would, and returns the resulting transitions.
func driveSeam(t *testing.T, d *Detector, pcm []int16) []Endpoint {
	t.Helper()
	p := filepath.Join(t.TempDir(), "fixture.wav")
	writePCM16Wav(t, p, pcm)
	inj, err := audio.NewWavInjector(p, audio.WithInjectorPace(0))
	if err != nil {
		t.Fatalf("wav injector: %v", err)
	}
	frames := seamFrameCount(len(pcm))
	buf := make(chan []byte, frames+4)
	if err := inj.Start(context.Background(), buf); err != nil {
		t.Fatalf("injector start: %v", err)
	}
	var out []Endpoint
	for i := 0; i < frames; i++ {
		select {
		case frame := <-buf:
			pcm, err := audio.DecodeFrame(frame)
			if err != nil {
				t.Fatalf("decode frame %d: %v", i, err)
			}
			out = d.Accept(pcmToFloat(pcm), out)
		case <-time.After(60 * time.Second):
			t.Fatalf("seam stalled at frame %d of %d", i, frames)
		}
	}
	out = d.Finalize(out)
	if err := inj.Stop(); err != nil {
		t.Fatalf("injector stop: %v", err)
	}
	if st := inj.Stats(); st.FramesDropped != 0 {
		t.Fatalf("fixture lost %d frame(s) to the bounded channel; the audio clock is broken",
			st.FramesDropped)
	}
	return out
}

// closings returns the endpoints that closed an utterance.
func closings(out []Endpoint) []Endpoint {
	var c []Endpoint
	for _, e := range out {
		if e.Kind != EndpointStart {
			c = append(c, e)
		}
	}
	return c
}

// kinds renders the transition log for the evidence table.
func kinds(out []Endpoint) []string {
	got := make([]string, 0, len(out))
	for _, e := range out {
		got = append(got, e.Kind.String())
	}
	return got
}

// speakPauseSpeak builds lead + body + gap + body + tail.
func speakPauseSpeak(leadMS, bodyMS, gapMS, tailMS int) []int16 {
	out := silence(msSamples(leadMS))
	out = append(out, synthSpeech(msSamples(bodyMS))...)
	out = append(out, silence(msSamples(gapMS))...)
	out = append(out, synthSpeech(msSamples(bodyMS))...)
	out = append(out, silence(msSamples(tailMS))...)
	return out
}

// TestVadPauseSplitPoint is AC#3 clause 1 through the real model, measured as a
// sweep rather than asserted at one nominal gap.
//
// Why a sweep: the pause the state machine sees is not the pause written into
// the fixture. silero keeps reporting speech for a while after this surrogate
// stops and takes a while to report it again, so a nominal pause loses ~400ms
// to the model before EndpointSilence starts counting. The sweep records where
// the split actually happens on this build, and the assertions sit one step
// either side of it so the fixtures cannot drift onto a knife edge.
func TestVadPauseSplitPoint(t *testing.T) {
	const leadMS, bodyMS, tailMS = 400, 1500, 1500
	splits := make(map[int]int)
	for _, gap := range []int{500, 900, 1100, 1500, 2000} {
		d := newTestDetector(t)
		closed := closings(driveSeam(t, d, speakPauseSpeak(leadMS, bodyMS, gap, tailMS)))
		splits[gap] = len(closed)
		var detail []string
		for _, e := range closed {
			detail = append(detail, fmt.Sprintf("%s/%dms(speech=%dms)", e.Kind,
				e.Duration().Milliseconds(), e.SpeechDuration().Milliseconds()))
		}
		t.Logf("gap=%4dms utterances=%d %v", gap, len(closed), detail)
	}

	// Held: the pause is bridged, so both bodies are one utterance.
	if n := splits[900]; n != 1 {
		t.Fatalf("900ms pause gave %d utterance(s), want 1 (bridged)", n)
	}
	// Split: the pause is long enough that the state machine's 700ms of
	// detected silence elapses.
	if n := splits[1500]; n != 2 {
		t.Fatalf("1500ms pause gave %d utterance(s), want 2 (split)", n)
	}
	// Monotonicity: a longer pause never re-merges what a shorter one split.
	for _, gap := range []int{1100, 1500, 2000} {
		if splits[gap] < splits[500] {
			t.Fatalf("splits are not monotonic: gap %d gave %d, gap 500 gave %d",
				gap, splits[gap], splits[500])
		}
	}
}

// TestVadPauseDetail records the per-utterance shape of the two settled
// fixtures, so a wrong split (one utterance that swallowed both bodies, or one
// that cut inside a body) is visible rather than merely counted.
func TestVadPauseDetail(t *testing.T) {
	d := newTestDetector(t)
	out := driveSeam(t, d, speakPauseSpeak(400, 1500, 1500, 1500))
	t.Logf("split fixture transitions: %v", kinds(out))
	closed := closings(out)
	if len(closed) != 2 {
		t.Fatalf("got %d utterance(s), want 2: %v", len(closed), kinds(out))
	}
	for i, e := range closed {
		if e.Kind != EndpointEnd {
			t.Fatalf("utterance %d closed by %s, want end-of-silence: %v", i, e.Kind, kinds(out))
		}
		if ms := e.SpeechDuration().Milliseconds(); ms < 600 || ms > 3000 {
			t.Fatalf("utterance %d carries %dms of speech, want one 1500ms body", i, ms)
		}
	}

	held := newTestDetector(t)
	out = driveSeam(t, held, speakPauseSpeak(400, 1500, 900, 1500))
	t.Logf("held fixture transitions: %v", kinds(out))
	closed = closings(out)
	if len(closed) != 1 {
		t.Fatalf("900ms pause gave %d utterance(s), want 1 held: %v", len(closed), kinds(out))
	}
	// Both bodies in one utterance: the speech inside it is the sum, not one body.
	if ms := closed[0].SpeechDuration().Milliseconds(); ms < 2500 {
		t.Fatalf("held utterance carries only %dms of speech, so the pause split it after all", ms)
	}
}

// TestVadPureSilenceProducesNoEndpoint is the positive control through the real
// model: three seconds of digital silence is not a sentence, not even a
// too-short one.
func TestVadPureSilenceProducesNoEndpoint(t *testing.T) {
	d := newTestDetector(t)
	out := driveSeam(t, d, silence(msSamples(3000)))
	if len(out) != 0 {
		t.Fatalf("3s of pure silence produced %v, want nothing", kinds(out))
	}
}

// TestVadForcedAt60s is AC#3 clause 3 through the real model: a run-on
// utterance is cut at the 60s guard rather than growing forever.
func TestVadForcedAt60s(t *testing.T) {
	d := newTestDetector(t)
	out := driveSeam(t, d, synthSpeech(msSamples(65_000)))
	var forced []Endpoint
	for _, e := range out {
		if e.Kind == EndpointMaxDuration {
			forced = append(forced, e)
		}
	}
	t.Logf("transitions: %v", kinds(out))
	if len(forced) != 1 {
		t.Fatalf("got %d max-duration close(s), want 1: %v", len(forced), kinds(out))
	}
	if d := forced[0].Duration().Milliseconds(); d != 60_000 {
		t.Fatalf("forced utterance reported %dms, want the 60000ms cap", d)
	}
	if n := len(forced[0].Samples); n != maxUtteranceSamples {
		t.Fatalf("forced utterance carried %d samples, want the %d sample cap", n, maxUtteranceSamples)
	}
}

// TestVadWindowBufferingIgnoresChunking: Detector.Accept buffers to whole
// silero windows, so an odd chunking must produce the same utterances as the
// seam's own 512-sample frames.
func TestVadWindowBufferingIgnoresChunking(t *testing.T) {
	d := newTestDetector(t)
	pcm := speakPauseSpeak(400, 1500, 2000, 1500)

	viaSeam := closings(driveSeam(t, d, pcm))

	samples := pcmToFloat(pcm)
	var out []Endpoint
	for off := 0; off < len(samples); off += 300 {
		end := off + 300
		if end > len(samples) {
			end = len(samples)
		}
		out = d.Accept(samples[off:end], out)
	}
	out = d.Finalize(out)
	viaOdd := closings(out)

	if len(viaSeam) != len(viaOdd) {
		t.Fatalf("seam gave %d utterance(s), 300-sample chunking gave %d: %v vs %v",
			len(viaSeam), len(viaOdd), kinds(viaSeam), kinds(viaOdd))
	}
	for i := range viaSeam {
		if viaSeam[i].Kind != viaOdd[i].Kind || viaSeam[i].Duration() != viaOdd[i].Duration() {
			t.Fatalf("utterance %d differs: seam=%s/%dms odd=%s/%dms", i,
				viaSeam[i].Kind, viaSeam[i].Duration().Milliseconds(),
				viaOdd[i].Kind, viaOdd[i].Duration().Milliseconds())
		}
	}
}

// TestVadConsumeStopsOnClosedChannel: Consume is the stream entry point and it
// must end on a closed channel (end of stream), not spin.
func TestVadConsumeStopsOnClosedChannel(t *testing.T) {
	d := newTestDetector(t)
	pcm := synthSpeech(msSamples(200))
	frame := audio.EncodeFrame(pcm[:audio.FrameSamples])
	frames := make(chan []byte, 1)
	frames <- frame
	close(frames)

	out, err := d.Consume(context.Background(), frames, nil)
	if err != nil {
		t.Fatalf("consume: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("200ms of speech should not open an utterance, got %v", kinds(out))
	}

	// Cancellation must also return promptly with ctx's error.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := d.Consume(ctx, make(chan []byte), nil); err == nil {
		t.Fatal("consume on a cancelled ctx returned nil error")
	}
}

// TestPinnedIntraOpThreads states the D32/D5 invariants the package pins, so a
// later leg cannot quietly relax them.
func TestPinnedIntraOpThreads(t *testing.T) {
	if IntraOpThreads != 1 {
		t.Fatalf("IntraOpThreads = %d, D32 pins it to 1", IntraOpThreads)
	}
	if WindowSamples != 512 {
		t.Fatalf("WindowSamples = %d, the silero window at 16k is 512", WindowSamples)
	}
	if audio.TargetRate != 16000 {
		t.Fatalf("seam rate = %d, want 16000", audio.TargetRate)
	}
	// D22 ban #4: no wall-clock arithmetic in deadline logic. Assert the
	// endpointing constants are an audio clock, not a wall clock.
	if samplesOf(EndpointSilence) != 11200 || samplesOf(MinSpeech) != 4800 ||
		samplesOf(MaxUtterance) != 960000 {
		t.Fatalf("policy samples = %d/%d/%d, want 11200/4800/960000",
			samplesOf(EndpointSilence), samplesOf(MinSpeech), samplesOf(MaxUtterance))
	}
}

// TestSileroModelFileIsPinned keeps the model identity visible: a swapped model
// would invalidate every reading in the receipt without failing anything else.
func TestSileroModelFileIsPinned(t *testing.T) {
	p := sileroModelPath(t)
	fi, err := os.Stat(p)
	if err != nil {
		t.Fatalf("stat model: %v", err)
	}
	t.Logf("silero model %s = %d bytes", filepath.Base(p), fi.Size())
	if fi.Size() < 100_000 {
		t.Fatalf("model %s is %d bytes, too small to be silero_vad.onnx", p, fi.Size())
	}
}
