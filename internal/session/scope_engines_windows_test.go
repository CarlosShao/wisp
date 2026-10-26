//go:build windows && wisp_voice_engines

package session

// Ticket 28's real-engine measurement: the SAME SessionScope that the policy
// cases drive with a fake loader, loaded with the production engines
// (internal/speech's silero Detector and streaming Paraformer), to answer the
// ticket's own question in numbers - "Cold vs Warm wake: second-round first
// token latency measured; Warm wake shows zero engine loads".
//
// WHY A BUILD TAG. This file needs the three staged sherpa DLLs and
// third_party/spike-models, exactly like internal/speech's own tests
// (15e-r1 section 8: PATH is not enough - the Windows loader finds System32's
// onnxruntime 1.17.1 before third_party's 1.28, and internal/speech refuses with
// a NAMED error instead of crashing). The package this file lives in is on the
// portable core scope, so the file must not compile there: `wisp_voice_engines`
// is set only by the staged invocation named in the evidence file. Without the
// tag `go test ./internal/session` is byte-for-byte the portable denominator.
//
// WHAT IS INJECTED AND WHAT IS NOT: only the Loader is a seam - it returns the
// REAL engines. The scope, the warm transitions, the refcounts and the cascade
// are the production code, which is the whole point: a warm path proven against
// a stub proves the stub.

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/speech"
)

const (
	// surrogate* is internal/speech's measured speech surrogate (120 Hz
	// harmonic stack, 3.5 Hz syllable envelope): the one synthetic signal out of
	// five candidates silero classifies as speech AND the streaming Paraformer
	// transcribes. It is duplicated rather than imported because those are
	// unexported test fixtures of another package.
	surrogateFundamental = 120.0
	surrogateHarmonics   = 8
	surrogateSyllableHz  = 3.5
	surrogateAmplitude   = 0.2

	// engineFixtureMS is 3 s of speech: 15-c's fixture length, long enough that a
	// transcript is the chain's output rather than an onset artefact.
	engineFixtureMS = 3000
	// engineFeedChunk is 512 samples = 32 ms at the seam rate, the chunk size the
	// production frame path feeds.
	engineFeedChunk = 512
	// sampleRate is audio.TargetRate, restated so the fixture does not depend on
	// an import for one number.
	sampleRate = 16000
)

// realEngines is the production Loader: the two engines that exist today. TTS,
// KWS and punctuation are reserved slots (tickets 26 / 41 / 27), so asking for
// them must fail rather than invent a handle - that is the same fail-closed
// shape scope_test.go's stub exercises.
type realEngines struct {
	silero     string
	spikeDir   string
	loaded     map[EngineKind]int
	closed     map[EngineKind]int
	vadTraffic int // windows the Detector actually classified, so a load that never ran is visible
	asrFinal   string
}

func (r *realEngines) load(kind EngineKind) (Engine, error) {
	switch kind {
	case EngineVAD:
		d, err := speech.NewDetector(speech.DetectorConfig{ModelPath: r.silero})
		if err != nil {
			return nil, err
		}
		r.count(r.loaded, kind)
		return &vadEngine{Detector: d, traffic: &r.vadTraffic, closed: &r.closed, kind: kind}, nil
	case EngineASR:
		a, err := speech.NewStreamingAsr(speech.AsrConfig{
			EncoderPath: filepath.Join(r.spikeDir, "asr-encoder.int8.onnx"),
			DecoderPath: filepath.Join(r.spikeDir, "asr-decoder.int8.onnx"),
			TokensPath:  filepath.Join(r.spikeDir, "asr-tokens.txt"),
		})
		if err != nil {
			return nil, err
		}
		r.count(r.loaded, kind)
		return &asrEngine{StreamingAsr: a, final: &r.asrFinal, closed: &r.closed, kind: kind}, nil
	}
	return nil, fmt.Errorf("ticket 28: %s 是 reserved slot（票 26/41/27 未落），本腿不造替身", string(kind))
}

func (r *realEngines) count(m map[EngineKind]int, kind EngineKind) {
	if m == nil {
		return
	}
	m[kind]++
}

// vadEngine / asrEngine add the counters the reading needs on top of the real
// engines. Both embed the production type, so the interface they satisfy is the
// real one and Close is the real release.
type vadEngine struct {
	*speech.Detector
	traffic *int
	closed  *map[EngineKind]int
	kind    EngineKind
}

func (v *vadEngine) Accept(samples []float32, out []speech.Endpoint) []speech.Endpoint {
	if v.traffic != nil {
		*v.traffic++
	}
	return v.Detector.Accept(samples, out)
}

func (v *vadEngine) Close() error {
	if v.closed != nil {
		(*v.closed)[v.kind]++
	}
	return v.Detector.Close()
}

type asrEngine struct {
	*speech.StreamingAsr
	final  *string
	closed *map[EngineKind]int
	kind   EngineKind
}

func (a *asrEngine) Finish() string {
	text := a.StreamingAsr.Finish()
	if a.final != nil {
		*a.final = text
	}
	return text
}

func (a *asrEngine) Close() error {
	if a.closed != nil {
		(*a.closed)[a.kind]++
	}
	return a.StreamingAsr.Close()
}

// TestTicket28ColdVsWarmWithTheRealEngines is AC1's measured half. One cold
// round (both engines loaded) then three warm rounds (Warm -> Resume -> borrow):
//
//   - the counter must stay at ZERO loads after the first warm entry (B4);
//   - the wake's own latency must collapse to the borrow, not the load;
//   - BOTH paths must still produce a real transcript (the positive control:
//     a warm path that quietly broke the cold path would be a worse bug than
//     the one this ticket fixes).
func TestTicket28ColdVsWarmWithTheRealEngines(t *testing.T) {
	silero, spike := engineFixturePaths(t)

	engines := &realEngines{
		silero: silero, spikeDir: spike,
		loaded: map[EngineKind]int{}, closed: map[EngineKind]int{},
	}
	clock := newFakeClock()
	s, err := NewScope(ScopeOptions{Loader: engines.load, Clock: clock, Hooks: Hooks{
		MicOpen:      func() {},
		MicClose:     func() {},
		PanelHide:    func() {},
		FreeOSMemory: func() {},
		GrantsExpire: func() {},
		Logf:         func(string, ...any) {},
	}})
	if err != nil {
		t.Fatalf("NewScope: %v", err)
	}

	speechSamples := synthSpeech(engineFixtureMS * sampleRate / 1000)

	// ---- round 0: COLD ----------------------------------------------------
	coldAcquire, coldFirst, coldFinal := engineRound(t, s, engines, speechSamples)
	if coldFinal == "" {
		t.Fatal("cold round produced an empty transcript: the cold path is the positive control and it is broken")
	}
	if engines.loaded[EngineVAD] != 1 || engines.loaded[EngineASR] != 1 {
		t.Fatalf("cold loads = vad:%d asr:%d, want 1/1", engines.loaded[EngineVAD], engines.loaded[EngineASR])
	}

	// ---- Warm: keep the session, drop the turn ----------------------------
	if err := s.EnterWarm("turn.end"); err != nil {
		t.Fatalf("enter warm: %v", err)
	}
	if got := s.Stats(); got.EnginesHot != 2 || got.LoadsAfterWarm != 0 {
		t.Fatalf("warm stats = %+v, want 2 hot engines and 0 loads after warm", got)
	}

	// ---- rounds 1..3: WARM -------------------------------------------------
	type warmReading struct {
		acquire, first time.Duration
		final          string
	}
	warm := make([]warmReading, 0, 3)
	for i := 0; i < 3; i++ {
		if err := s.Resume(); err != nil {
			t.Fatalf("resume %d: %v", i, err)
		}
		acq, first, final := engineRound(t, s, engines, speechSamples)
		if final == "" {
			t.Fatalf("warm round %d produced an empty transcript: the warm path changed behaviour, not just speed", i)
		}
		// The latency claim, asserted where it happens: a wake must cost at most a
		// TENTH of a cold load. Warm is not "a bit cheaper" - it is the difference
		// between borrowing a pointer and loading 165 MB of int8 encoder.
		if acq > coldAcquire/10 {
			t.Fatalf("warm round %d acquire = %s (%dns), want at most a tenth of the cold load %s (%dns): the wake paid for a load",
				i, acq, acq.Nanoseconds(), coldAcquire, coldAcquire.Nanoseconds())
		}
		warm = append(warm, warmReading{acquire: acq, first: first, final: final})
		if err := s.EnterWarm("turn.end"); err != nil {
			t.Fatalf("re-enter warm %d: %v", i, err)
		}
	}

	st := s.Stats()
	if st.LoadsAfterWarm != 0 {
		t.Fatalf("LoadsAfterWarm = %d after 3 warm wakes, want 0 (B4)", st.LoadsAfterWarm)
	}
	if engines.loaded[EngineVAD] != 1 || engines.loaded[EngineASR] != 1 {
		t.Fatalf("loader ran again during warm rounds: vad:%d asr:%d, want 1/1 ever",
			engines.loaded[EngineVAD], engines.loaded[EngineASR])
	}
	if st.WarmResumes != 3 {
		t.Fatalf("WarmResumes = %d, want 3", st.WarmResumes)
	}
	if engines.vadTraffic == 0 {
		t.Fatal("the VAD handle never classified a window: the engine was loaded but never used, so 'hot' would be a claim about a pointer")
	}

	// ---- dispose: the cascade must close the real engines ------------------
	disposeStart := time.Now()
	if err := s.Dispose("evidence.end"); err != nil {
		t.Fatalf("dispose: %v", err)
	}
	disposeFor := time.Since(disposeStart)
	if engines.closed[EngineVAD] != 1 || engines.closed[EngineASR] != 1 {
		t.Fatalf("real engine closes = vad:%d asr:%d, want 1/1", engines.closed[EngineVAD], engines.closed[EngineASR])
	}
	if clock.Pending() != 0 {
		t.Fatalf("%d timers still armed after dispose", clock.Pending())
	}

	// ---- positive control: the COLD path still works AFTER the warm rounds ---
	// A warm pool that quietly broke the cold path would be a worse bug than the
	// one this ticket fixes, so a brand-new session is opened and must pay the
	// load again and still transcribe. This is also what proves "warm" is a
	// per-session property: nothing global was left loaded behind it.
	s2, err := NewScope(ScopeOptions{Loader: engines.load, Clock: clock, Hooks: Hooks{
		MicClose: func() {}, PanelHide: func() {}, FreeOSMemory: func() {}, GrantsExpire: func() {},
		Logf: func(string, ...any) {},
	}})
	if err != nil {
		t.Fatalf("second session NewScope: %v", err)
	}
	coldAgain, _, finalAgain := engineRound(t, s2, engines, speechSamples)
	if finalAgain == "" {
		t.Fatal("a NEW session after the warm rounds produced an empty transcript: the cold path was poisoned by warming")
	}
	if coldAgain <= coldAcquire/10 {
		t.Fatalf("a new session's acquire = %s, want a cold load again (about %s): a warm pool that leaks across sessions is a leak, not a cache",
			coldAgain, coldAcquire)
	}
	if got := s2.Stats(); got.ColdLoads != 2 || got.LoadsAfterWarm != 0 {
		t.Fatalf("new session stats = %+v, want ColdLoads=2 LoadsAfterWarm=0", got)
	}
	if err := s2.Dispose("control.end"); err != nil {
		t.Fatalf("second session dispose: %v", err)
	}

	payload := len([]rune(coldFinal))
	t.Logf("AC1 cold:   acquire=%s (%dns) first_partial=%s first_request=%s final=%q (%d runes)",
		coldAcquire, coldAcquire.Nanoseconds(), coldFirst, coldAcquire+coldFirst, coldFinal, payload)
	for i, w := range warm {
		t.Logf("AC1 warm %d: acquire=%s (%dns) first_partial=%s first_request=%s final=%q",
			i+1, w.acquire, w.acquire.Nanoseconds(), w.first, w.acquire+w.first, w.final)
	}
	t.Logf("AC1 speedup: first_request cold=%s warm=%s ratio=%.1fx; acquire cold=%dns warm=%dns; loads after warm=%d; dispose cascade=%s",
		coldAcquire+coldFirst, warm[0].acquire+warm[0].first,
		float64(coldAcquire+coldFirst)/float64(warm[0].acquire+warm[0].first),
		coldAcquire.Nanoseconds(), warm[0].acquire.Nanoseconds(),
		st.LoadsAfterWarm, disposeFor)
	t.Logf("positive control: new session after the warm rounds, acquire=%s (%dns), final=%q (cold path still pays the load and still transcribes)",
		coldAgain, coldAgain.Nanoseconds(), finalAgain)
}

// engineRound borrows both engines (timing the borrow), runs one utterance
// through the ASR while timing the first partial, and returns the final
// transcript. The borrow is what Warm makes cheap; the decode is what it must
// NOT change.
func engineRound(t *testing.T, s *Scope, engines *realEngines, samples []float32) (time.Duration, time.Duration, string) {
	t.Helper()
	start := time.Now()
	vadRaw, err := s.Acquire(EngineVAD)
	if err != nil {
		t.Fatalf("acquire vad: %v", err)
	}
	asrRaw, err := s.Acquire(EngineASR)
	if err != nil {
		t.Fatalf("acquire asr: %v", err)
	}
	acquire := time.Since(start)
	defer func() {
		if err := s.Release(EngineVAD); err != nil {
			t.Fatalf("release vad: %v", err)
		}
		if err := s.Release(EngineASR); err != nil {
			t.Fatalf("release asr: %v", err)
		}
	}()

	vad, ok := vadRaw.(*vadEngine)
	if !ok {
		t.Fatalf("vad handle is %T, want the real Detector wrapper", vadRaw)
	}
	asr, ok := asrRaw.(*asrEngine)
	if !ok {
		t.Fatalf("asr handle is %T, want the real StreamingAsr wrapper", asrRaw)
	}
	// Feed the detector one window so a loaded-but-broken session is not "hot".
	vad.Accept(samples[:len(samples)/2], nil)

	asr.Begin()
	var first time.Duration
	for off := 0; off < len(samples); off += engineFeedChunk {
		end := off + engineFeedChunk
		if end > len(samples) {
			end = len(samples)
		}
		asr.Feed(samples[off:end])
		if first == 0 && asr.Partial() != "" {
			first = time.Since(start) - acquire
		}
	}
	final := asr.Finish()
	return acquire, first, final
}

// synthSpeech is the measured surrogate, in the float shape AsrEngine wants.
func synthSpeech(n int) []float32 {
	out := make([]float32, n)
	for i := range out {
		tt := float64(i) / float64(sampleRate)
		env := 0.5 + 0.5*math.Sin(2*math.Pi*surrogateSyllableHz*tt)
		var acc float64
		for h := 1; h <= surrogateHarmonics; h++ {
			acc += math.Sin(2*math.Pi*surrogateFundamental*float64(h)*tt) / float64(h)
		}
		out[i] = float32(math.Max(-1, math.Min(1, surrogateAmplitude*env*acc)))
	}
	return out
}

// engineFixturePaths finds the pinned models. The staged invocation exports both
// variables; the walk-up is the same convenience internal/speech's fixtures
// offer so the binary also works when run from the repository root.
func engineFixturePaths(t *testing.T) (silero, spikeDir string) {
	t.Helper()
	spikeDir = os.Getenv("WISP_SPIKE_MODELS")
	silero = os.Getenv("WISP_SILERO_MODEL")
	if spikeDir == "" {
		spikeDir = walkUpFor(t, filepath.Join("third_party", "spike-models"))
	}
	if silero == "" {
		silero = filepath.Join(spikeDir, "silero_vad.onnx")
	}
	for _, p := range []string{
		silero,
		filepath.Join(spikeDir, "asr-encoder.int8.onnx"),
		filepath.Join(spikeDir, "asr-decoder.int8.onnx"),
		filepath.Join(spikeDir, "asr-tokens.txt"),
	} {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("ticket 28 staging: %v (set WISP_SILERO_MODEL / WISP_SPIKE_MODELS, or run with cwd inside the repo; runtime.GOOS=%s)", err, runtime.GOOS)
		}
	}
	return silero, spikeDir
}

func walkUpFor(t *testing.T, rel string) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for i := 0; i < 8; i++ {
		cand := filepath.Join(dir, rel)
		if _, err := os.Stat(cand); err == nil {
			return cand
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatalf("ticket 28 staging: %s not found above the working directory", rel)
	return ""
}
