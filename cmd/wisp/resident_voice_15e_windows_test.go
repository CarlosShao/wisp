//go:build windows

package main

// Ticket 15e's readings, driven through the production assembly rather than
// around it: the frames go through the real C8 seam (audio.WavInjector, the
// project's sanctioned injection point), the chain is the real internal/speech
// chain (silero + streaming Paraformer), and the hand-off is the real
// residentVoiceHandOff -> models.go modelStore.handOffModel walk against the
// signed manifest in build/models. What is injected is only what the assembly
// root itself injects in production: the frame channel, the plan, the ball host
// and the task sink.
//
// The synthetic speech surrogate is the one internal/speech's own fixtures
// measured (120 Hz harmonic stack, 3.5 Hz syllable envelope): of five candidate
// signals it is the one silero accepts as speech AND the streaming Paraformer
// transcribes. It is duplicated here because a test fixture may not be imported
// across packages as an unexported helper, and this slice may not touch
// internal/speech.
//
// These cases need the sherpa DLLs staged next to the test binary and the model
// files reachable (see 15e-r1 §staging): `go test ./cmd/wisp` alone will skip the
// chain cases by failing the model load, so the evidence table always names the
// staged invocation.

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

const (
	// voice15eLeadMS is the silence in front of the speech in the end-to-end
	// fixture. It is the seam position the ball latency is measured from, so it
	// is a constant of this file rather than a number buried in a sample count.
	voice15eLeadMS = 250
	// voice15eSpeechMS is 15-c's fixture length: 3 s of the surrogate is what it
	// measured partials and a final on, and a shorter run would make the final's
	// existence depend on silero's onset overshoot instead of on the chain.
	voice15eSpeechMS = 3000
	// voice15eTailMS is the trailing silence. It is generous on purpose, twice
	// over: 15-b measured ~400 ms of silero overshoot on synthetic pauses, so a
	// 700 ms tail would not reliably reach EndpointSilence; and the switchover's
	// replay decodes synchronously, during which the bounded channel drops the
	// NEWEST frames (~0.5-0.8 s measured), so a tail that is still long enough
	// after that loss is what keeps this case from depending on which of the two
	// happened first.
	voice15eTailMS = 2500
)

// voice15eSink is ticket 12's text path as a recorder: the voice leg may only
// call submitTask, and this is the whole of that contract.
type voice15eSink struct {
	mu  sync.Mutex
	got []string
}

func (s *voice15eSink) submitTask(text string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.got = append(s.got, text)
	return nil
}

func (s *voice15eSink) texts() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.got...)
}

// voice15eDataDir is a data root on the SAME volume as the repository, because
// the model store is installed by hardlinking the pinned spike models into the
// manifest's file names - a link across volumes would silently become a 237 MB
// copy on every run. It honours WISP_15E_TMP and falls back to the system temp
// directory (with a copy) for a machine that has no D: drive.
func voice15eDataDir(t *testing.T) string {
	t.Helper()
	if base := os.Getenv("WISP_15E_TMP"); base != "" {
		if dir, err := os.MkdirTemp(base, "case-"); err == nil {
			return dir
		}
	}
	if dir, err := os.MkdirTemp(`D:\tmp\wisp\20261026\15er1`, "case-"); err == nil {
		return dir
	}
	return t.TempDir()
}

// voice15eVadModel is the silero file the plan opens. It honours the same env
// var 15-b/15-c's fixtures use, then the repository tree, so a staged binary run
// from anywhere still finds it.
func voice15eVadModel(t *testing.T) string {
	t.Helper()
	if p := os.Getenv("WISP_SILERO_MODEL"); p != "" {
		return p
	}
	return filepath.Join(voice15eRepoFile(t, filepath.Join("third_party", "spike-models", "silero_vad.onnx")))
}

// voice15eSpikeModels is the pinned model directory (the same env var convention
// 15-c's cases use).
func voice15eSpikeModels(t *testing.T) string {
	t.Helper()
	if p := os.Getenv("WISP_SPIKE_MODELS"); p != "" {
		return p
	}
	return voice15eRepoFile(t, filepath.Join("third_party", "spike-models"))
}

// voice15eRepoFile walks up from the test binary's working directory to the repo
// root (the same shape internal/speech's fixtures use).
func voice15eRepoFile(t *testing.T, rel string) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for i := 0; i < 8; i++ {
		candidate := filepath.Join(dir, rel)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatalf("cannot find %s above the working directory (runtime.GOOS=%s)", rel, runtime.GOOS)
	return ""
}

// voice15eInstallAsr makes the data root's model store look like a completed
// install of the streaming ASR: the manifest's own file names, the manifest's own
// bytes (the pinned spike copies hash-match it - measured, see 15e-r1 §2).
func voice15eInstallAsr(t *testing.T, dataDir string) string {
	t.Helper()
	if err := voice15eInstallAsrFiles(dataDir); err != nil {
		t.Fatalf("install the streaming ASR into %s: %v", dataDir, err)
	}
	return filepath.Join(dataDir, "models", voiceAsrModelID)
}

// voice15eInstallAsrFiles is the same install without a *testing.T, because the
// end-to-end case performs it INSIDE the hand-off seam: in production the walk is
// what makes the bytes available (Ensure downloads and verifies them), so the
// fixture makes the walk the only thing that puts them there too - which is what
// makes "delete the hand-off call" a product failure and not just a counter one.
func voice15eInstallAsrFiles(dataDir string) error {
	store := filepath.Join(dataDir, "models", voiceAsrModelID)
	if err := os.MkdirAll(store, 0o755); err != nil {
		return err
	}
	src := os.Getenv("WISP_SPIKE_MODELS")
	if src == "" {
		return fmt.Errorf("WISP_SPIKE_MODELS is not set")
	}
	for _, pair := range [][2]string{
		{voiceAsrEncoder, "asr-encoder.int8.onnx"},
		{voiceAsrDecoder, "asr-decoder.int8.onnx"},
		{voiceAsrTokens, "asr-tokens.txt"},
	} {
		from := filepath.Join(src, pair[1])
		to := filepath.Join(store, pair[0])
		if _, err := os.Stat(to); err == nil {
			continue
		}
		if err := os.Link(from, to); err == nil {
			continue
		}
		blob, err := os.ReadFile(from)
		if err != nil {
			return err
		}
		if err := os.WriteFile(to, blob, 0o644); err != nil {
			return err
		}
	}
	return nil
}

// voice15eWriteConfig lays a config.toml down from the shipped default table.
func voice15eWriteConfig(t *testing.T, dir string, mutate func(*config.Config)) {
	t.Helper()
	c := config.NewDefaults()
	if mutate != nil {
		mutate(c)
	}
	if err := config.SaveFile(filepath.Join(dir, configFileName), c); err != nil {
		t.Fatalf("write %s: %v", configFileName, err)
	}
}

// voice15eSynth renders the fixture the way internal/speech's fixtures do, and
// returns 16k mono int16.
func voice15eSynth(leadMS, speechMS, tailMS int) []int16 {
	total := (leadMS + speechMS + tailMS) * audio.TargetRate / 1000
	out := make([]int16, total)
	lead := leadMS * audio.TargetRate / 1000
	speech := speechMS * audio.TargetRate / 1000
	for i := 0; i < speech && lead+i < total; i++ {
		tt := float64(i) / float64(audio.TargetRate)
		env := 0.5 + 0.5*math.Sin(2*math.Pi*3.5*tt)
		var acc float64
		for h := 1; h <= 8; h++ {
			acc += math.Sin(2*math.Pi*120*float64(h)*tt) / float64(h)
		}
		out[lead+i] = int16(math.Max(-1, math.Min(1, 0.2*env*acc)) * 32767)
	}
	return out
}

// voice15eWriteWav writes a canonical 44-byte-header 16k/mono PCM16 RIFF file so
// the fixture travels through the real C8 seam.
func voice15eWriteWav(t *testing.T, path string, samples []int16) {
	t.Helper()
	data := make([]byte, 44+len(samples)*2)
	copy(data[0:4], "RIFF")
	binary.LittleEndian.PutUint32(data[4:8], uint32(36+len(samples)*2))
	copy(data[8:12], "WAVE")
	copy(data[12:16], "fmt ")
	binary.LittleEndian.PutUint32(data[16:20], 16)
	binary.LittleEndian.PutUint16(data[20:22], 1)
	binary.LittleEndian.PutUint16(data[22:24], 1)
	binary.LittleEndian.PutUint32(data[24:28], uint32(audio.TargetRate))
	binary.LittleEndian.PutUint32(data[28:32], uint32(audio.TargetRate*2))
	binary.LittleEndian.PutUint16(data[32:34], 2)
	binary.LittleEndian.PutUint16(data[34:36], 16)
	copy(data[36:40], "data")
	binary.LittleEndian.PutUint32(data[40:44], uint32(len(samples)*2))
	for i, s := range samples {
		binary.LittleEndian.PutUint16(data[44+2*i:], uint16(s))
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// voice15eStartLeg builds and starts the production consumer on a bounded frame
// channel, and stops it at cleanup in the same order D38(e) step 4 uses.
func voice15eStartLeg(t *testing.T, plan voicePlan, ball *residentBall, sink voiceTranscriptSink) (*residentVoice, chan []byte) {
	t.Helper()
	frames := audio.NewBoundedFrames()
	rv := newResidentVoice(frames, plan, ball, sink)
	rv.start(nil) // nil registry = observe.Default, the same default a bare spawn takes
	t.Cleanup(rv.stop)
	if !rv.drainsFrames() {
		t.Fatalf("the leg did not start: verdict=%q", rv.verdict)
	}
	return rv, frames
}

// voice15eWaitFor polls one reading until it is true, so a case never depends on
// a sleep long enough to be a guess.
func voice15eWaitFor(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", what)
}

// voice15eArm proves the consumer loop is live before the fixture starts and
// puts exactly one frame of digital silence in front of it. Both halves matter for
// the latency reading: a leg whose loop has not been scheduled yet drops the first
// frames of whatever the microphone says (the bounded channel holds 192 ms), and a
// dropped frame before a speech onset moves the derived utterance start by exactly
// one window - measured as a 32 ms step between two runs that were otherwise
// identical. The returned value is how much silence this frame put in front of the
// fixture, on the seam clock.
func voice15eArm(t *testing.T, frames chan<- []byte, rv *residentVoice) int64 {
	t.Helper()
	silent := audio.EncodeFrame(make([]int16, audio.FrameSamples))
	select {
	case frames <- silent:
	case <-time.After(10 * time.Second):
		t.Fatal("the leg never took a frame: the consumer loop is not draining")
	}
	voice15eWaitFor(t, "the consumer loop to be live", func() bool { return rv.drained.Load() == 1 })
	return int64(audio.FrameSamples)
}

// Test15eSpeechEndToEndHandsTheTranscriptToTheAgent is reading #1, #2 and #5 in
// one run: audio in -> VAD stop -> ASR final -> the model hand-off really called
// (once, with the manifest id) -> the transcript handed to ticket 12's text path,
// with the ball walked Listening -> Thinking on the way.
func Test15eSpeechEndToEndHandsTheTranscriptToTheAgent(t *testing.T) {
	dataDir := voice15eDataDir(t)
	// NOT installed yet: the store is filled inside the hand-off seam below, so
	// "delete the hand-off call" makes the engine unable to open a model - a
	// product failure, not a counter mismatch.
	store := filepath.Join(dataDir, "models", voiceAsrModelID)
	voice15eWriteConfig(t, dataDir, nil)
	manifest := voice15eRepoFile(t, filepath.Join("build", "models", "manifest.json"))
	t.Setenv("WISP_MODELS_MANIFEST", manifest)

	sink := &voice15eSink{}
	ball := &residentBall{}
	seen := make([]string, 0, 2)
	var seenMu sync.Mutex
	realHandOff := residentVoiceHandOff(dataDir)
	installErr := error(nil)
	var installMu sync.Mutex
	plan := voicePlan{
		VadModelPath: voice15eVadModel(t),
		AsrDir:       store,
		Source:       "test",
		HandOff: func(id string, logf func(string, ...any)) int {
			seenMu.Lock()
			seen = append(seen, id)
			seenMu.Unlock()
			// The walk is what makes the bytes available (in production it is
			// Ensure; here it is the pinned copies, linked under the manifest's
			// names), so this call is the ONLY thing that can produce a model.
			if err := voice15eInstallAsrFiles(dataDir); err != nil {
				installMu.Lock()
				installErr = err
				installMu.Unlock()
				return 2
			}
			return realHandOff(id, logf)
		},
	}
	rv, frames := voice15eStartLeg(t, plan, ball, sink)
	warm := voice15eArm(t, frames, rv)

	wav := filepath.Join(dataDir, "e2e.wav")
	voice15eWriteWav(t, wav, voice15eSynth(voice15eLeadMS, voice15eSpeechMS, voice15eTailMS))
	totalSamples := (voice15eLeadMS + voice15eSpeechMS + voice15eTailMS) * audio.TargetRate / 1000
	wantFrames := (totalSamples+audio.FrameSamples-1)/audio.FrameSamples + 1

	inj, err := audio.NewWavInjector(wav)
	if err != nil {
		t.Fatalf("wav injector: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := inj.Start(ctx, frames); err != nil {
		t.Fatalf("wav injector start: %v", err)
	}
	defer func() { _ = inj.Stop() }()

	voice15eWaitFor(t, "the utterance to be transcribed and handed over", func() bool {
		return rv.finals.Load() >= 1
	})

	// #1 the hand-off: exactly once, with the manifest id, and it really was the
	// production walk (exit 0 means models.go announced the model available).
	if got := rv.handOffs.Load(); got != 1 {
		t.Fatalf("the model hand-off ran %d times, want exactly 1 (a per-utterance hand-off would decompress the whole point)", got)
	}
	if ids := rv.handOffIDsSnapshot(); len(ids) != 1 || ids[0] != voiceAsrModelID {
		t.Fatalf("the hand-off was asked for %v, want exactly [%s]", ids, voiceAsrModelID)
	}
	if exit := rv.handOffExit.Load(); exit != 0 {
		t.Fatalf("the production hand-off walk returned %d, want 0 (the store was installed from the pinned bytes)", exit)
	}
	installMu.Lock()
	if installErr != nil {
		installMu.Unlock()
		t.Fatalf("the seam could not install the pinned model bytes, so the walk's exit 0 means nothing: %v", installErr)
	}
	installMu.Unlock()
	seenMu.Lock()
	seenIDs := append([]string(nil), seen...)
	seenMu.Unlock()
	if len(seenIDs) != 1 || seenIDs[0] != voiceAsrModelID {
		t.Fatalf("the injected seam saw %v, want exactly one call with %s", seenIDs, voiceAsrModelID)
	}
	if !rv.asrLoaded.Load() {
		t.Fatalf("the ASR bytes were never opened after a successful hand-off")
	}

	// #1 (the other end) the transcript reached ticket 12's text path.
	want := sink.texts()
	if len(want) != 1 || strings.TrimSpace(want[0]) == "" {
		t.Fatalf("the text path received %v, want exactly one non-empty transcript (problems=%v)", want, rv.problemsSnapshot())
	}
	turns := rv.turnsSnapshot()
	if len(turns) != 1 || turns[0].Text != want[0] || !turns[0].Handed {
		t.Fatalf("turn record %+v does not match the handed transcript %q", turns, want[0])
	}

	// #2 the ball walked the D43 rows in order, and the projection is on record.
	states := ball.voiceStatesSnapshot()
	if len(states) != 2 || states[0] != statemachine.StateListening || states[1] != statemachine.StateThinking {
		t.Fatalf("the ball was projected with %v, want [Listening Thinking]", states)
	}
	timeline := rv.timelineSnapshot()
	if len(timeline) != 2 {
		t.Fatalf("the D43 timeline has %d edges, want 2 (%+v)", len(timeline), timeline)
	}
	if timeline[0].From != statemachine.StateSleeping || timeline[0].Event != statemachine.EvSummon || timeline[0].To != statemachine.StateListening {
		t.Fatalf("first edge %+v, want Sleeping+session.summon->Listening (D43#4)", timeline[0])
	}
	if timeline[1].From != statemachine.StateListening || timeline[1].Event != statemachine.EvVadStop || timeline[1].To != statemachine.StateThinking {
		t.Fatalf("second edge %+v, want Listening+asr.vad-stop->Thinking (D43#11)", timeline[1])
	}

	// #5 the two latency segments on the audio clock, and their composition.
	onset := warm + int64(voice15eLeadMS*audio.TargetRate/1000)
	openDelay := time.Duration(turns[0].StartSample-onset) * time.Second / time.Duration(audio.TargetRate)
	chunk := turns[0].FirstChunk
	ballLatency := openDelay + chunk
	waitToText := time.Duration(turns[0].FirstTextAt-onset) * time.Second / time.Duration(audio.TargetRate)
	if openDelay <= 0 {
		t.Fatalf("open delay %v: the utterance start was not derived from the seam clock (%+v)", openDelay, turns[0])
	}
	if chunk <= 0 {
		t.Fatalf("chunk segment %v: no first partial was recorded (%+v)", chunk, turns[0])
	}
	t.Logf("LATENCY openDelay=%v chunk=%v ballLatency=%v (openDelay+chunk) seamWaitToFirstText=%v finalLatency=%v speechMS=%d closedBy=%s",
		openDelay, chunk, ballLatency, waitToText, turns[0].FinalLatency, turns[0].SpeechMS, turns[0].ClosedBy)
	t.Logf("COUNTERS drained=%d/%d partials=%d finals=%d handed=%d handoffs=%d asrLoaded=%v problems=%v",
		rv.drained.Load(), wantFrames, rv.partials.Load(), rv.finals.Load(), rv.handed.Load(),
		rv.handOffs.Load(), rv.asrLoaded.Load(), rv.problemsSnapshot())
	t.Logf("TEXT %q", want[0])
	t.Logf("TIMELINE %+v", timeline)

	if rv.drained.Load() < uint64(wantFrames/2) {
		t.Fatalf("only %d of %d frames were drained: the C8 channel is not being consumed", rv.drained.Load(), wantFrames)
	}
}

// Test15eSilenceNeverCallsTheHandOff is reading #7: the positive control. Pure
// silence must not fetch a model, must not open the ASR, must not walk the ball
// and must not produce a transcript - and the frames must still be drained, so a
// green here cannot be "nothing was running".
func Test15eSilenceNeverCallsTheHandOff(t *testing.T) {
	dataDir := voice15eDataDir(t)
	ball := &residentBall{}
	sink := &voice15eSink{}
	calls := 0
	var callsMu sync.Mutex
	plan := voicePlan{
		VadModelPath: voice15eVadModel(t),
		AsrDir:       filepath.Join(dataDir, "models", voiceAsrModelID),
		Source:       "test",
		HandOff: func(id string, logf func(string, ...any)) int {
			callsMu.Lock()
			calls++
			callsMu.Unlock()
			return 0
		},
	}
	rv, frames := voice15eStartLeg(t, plan, ball, sink)

	const pushed = 120 // 3.84s of digital silence
	silent := audio.EncodeFrame(make([]int16, audio.FrameSamples))
	for i := 0; i < pushed; i++ {
		select {
		case frames <- silent:
		case <-time.After(5 * time.Second):
			t.Fatalf("push %d/%d: the leg stopped draining", i+1, pushed)
		}
		voice15eWaitFor(t, fmt.Sprintf("frame %d to be drained", i+1), func() bool {
			return rv.drained.Load() == uint64(i+1)
		})
	}

	callsMu.Lock()
	got := calls
	callsMu.Unlock()
	if got != 0 || rv.handOffs.Load() != 0 {
		t.Fatalf("silence called the model hand-off (%d seam calls, %d leg calls): the capability must be first USE, not first frame",
			got, rv.handOffs.Load())
	}
	if rv.asrLoaded.Load() {
		t.Fatalf("silence opened the ASR bytes")
	}
	if rv.finals.Load() != 0 || len(sink.texts()) != 0 {
		t.Fatalf("silence produced %d finals and handed %v", rv.finals.Load(), sink.texts())
	}
	if states := ball.voiceStatesSnapshot(); len(states) != 0 {
		t.Fatalf("silence walked the ball: %v", states)
	}
	if len(rv.timelineSnapshot()) != 0 || rv.machine.State() != statemachine.StateSleeping {
		t.Fatalf("silence moved the D43 machine: state=%s timeline=%+v", rv.machine.State(), rv.timelineSnapshot())
	}
	t.Logf("SILENCE drained=%d handoffs=%d finals=%d ball=%v state=%s",
		rv.drained.Load(), rv.handOffs.Load(), rv.finals.Load(), ball.voiceStatesSnapshot(), rv.machine.State())
}

// Test15eARefusedHandOffLeavesTheAsrUnopened is the refusal's control: the ball
// still walks into Listening (D43#4 fires before any I/O), the ASR is never
// opened, and the leg says why instead of transcribing with unverified bytes.
func Test15eARefusedHandOffLeavesTheAsrUnopened(t *testing.T) {
	dataDir := voice15eDataDir(t)
	ball := &residentBall{}
	sink := &voice15eSink{}
	plan := voicePlan{
		VadModelPath: voice15eVadModel(t),
		AsrDir:       filepath.Join(dataDir, "models", voiceAsrModelID),
		Source:       "test",
		HandOff:      func(id string, logf func(string, ...any)) int { return 1 },
	}
	rv, frames := voice15eStartLeg(t, plan, ball, sink)

	// A real seam, one short utterance: the refusal happens the moment it opens.
	wav := filepath.Join(dataDir, "refused.wav")
	voice15eWriteWav(t, wav, voice15eSynth(100, 1800, 200))
	inj, err := audio.NewWavInjector(wav)
	if err != nil {
		t.Fatalf("wav injector: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := inj.Start(ctx, frames); err != nil {
		t.Fatalf("injector start: %v", err)
	}
	defer func() { _ = inj.Stop() }()

	voice15eWaitFor(t, "the hand-off refusal to be recorded", func() bool {
		return rv.handOffExit.Load() == 1
	})
	voice15eWaitFor(t, "the leg to give up the chain", func() bool {
		for _, p := range rv.problemsSnapshot() {
			if strings.Contains(p, "模型交还被拒绝") {
				return true
			}
		}
		return false
	})
	if rv.asrLoaded.Load() {
		t.Fatalf("a refused hand-off still opened the ASR bytes")
	}
	if rv.finals.Load() != 0 || len(sink.texts()) != 0 {
		t.Fatalf("a refused hand-off still produced a transcript: %v", sink.texts())
	}
	if states := ball.voiceStatesSnapshot(); len(states) != 1 || states[0] != statemachine.StateListening {
		t.Fatalf("ball states after a refusal = %v, want exactly [Listening] (D43#4 must precede the I/O)", states)
	}
	t.Logf("REFUSED handoffs=%d exit=%d asrLoaded=%v ball=%v problems=%v",
		rv.handOffs.Load(), rv.handOffExit.Load(), rv.asrLoaded.Load(), ball.voiceStatesSnapshot(), rv.problemsSnapshot())
}

// Test15eTheChannelHasNoConsumerWithoutTheLeg is reading #3's control half: the
// bounded channel really does refuse frames when nobody drains it, which is the
// state the old comment described as "nothing drains it" - measured on the same
// constructor production uses.
func Test15eTheChannelHasNoConsumerWithoutTheLeg(t *testing.T) {
	buf := audio.NewBoundedFrames()
	frame := audio.EncodeFrame(make([]int16, audio.FrameSamples))
	accepted := 0
	for i := 0; i < 20; i++ {
		select {
		case buf <- frame:
			accepted++
		default:
		}
	}
	if want := audio.BoundedFrameCapacity(); accepted != want {
		t.Fatalf("an undrained channel accepted %d frames, want exactly its capacity %d", accepted, want)
	}
	t.Logf("NO CONSUMER: capacity=%d accepted=%d refused=%d window=%v",
		audio.BoundedFrameCapacity(), accepted, 20-accepted, audio.BoundedWindow)
}

// Test15ePauseStopsTheListenerAndTheTrayFollows is reading #4: the tray's
// pause-wake item really stops the consumer (frames are drained and dropped, not
// fed), the checkmark's second bit follows that fact, and resuming puts the
// listener back.
func Test15ePauseStopsTheListenerAndTheTrayFollows(t *testing.T) {
	dataDir := voice15eDataDir(t)
	ball := &residentBall{}
	plan := voicePlan{
		VadModelPath: voice15eVadModel(t),
		AsrDir:       filepath.Join(dataDir, "models", voiceAsrModelID),
		Source:       "test",
		HandOff:      func(id string, logf func(string, ...any)) int { return 0 },
	}
	rv, frames := voice15eStartLeg(t, plan, ball, nil)

	// The production attach, exactly as startResidentVoice does it.
	ball.attachWakePause(rv.togglePause, rv.pausedForTray)
	pushed := &trayCheckRec{}
	ball.attachTrayMuteProjection(func() (bool, bool) { return true, true }, pushed.push)

	silent := audio.EncodeFrame(make([]int16, audio.FrameSamples))
	total := uint64(0)
	// pushFrame pushes one frame and waits until the leg has accounted for it
	// (drained or paused-dropped), so no assertion below races the consumer.
	pushFrame := func(what string) {
		t.Helper()
		total++
		select {
		case frames <- silent:
		case <-time.After(5 * time.Second):
			t.Fatalf("%s: the leg is not draining", what)
		}
		voice15eWaitFor(t, what, func() bool {
			return rv.drained.Load()+rv.pausedDrops.Load() >= total
		})
	}

	// Before any pause: frames are drained into the chain.
	pushFrame("the first frame")
	if got := rv.drained.Load(); got != 1 {
		t.Fatalf("drained=%d before the pause, want 1", got)
	}

	// The tray item pauses the listener; the projection must now say so.
	sentence, executed := ball.currentWakePauseToggle()()
	if !executed || !strings.Contains(sentence, "已暂停") {
		t.Fatalf("the tray pause-wake executor answered %q executed=%v, want a pause sentence", sentence, executed)
	}
	if !ball.mirrorTrayMute() {
		t.Fatalf("the tray projection wrote nothing after a pause")
	}
	if last := pushed.last(t); !last.pausedWake {
		t.Fatalf("the tray checkmark did not follow the paused listener: %+v", last)
	}

	// The pause is observed by the consumer on its NEXT turn, so one frame that
	// was already in flight when the user pressed the item can still be fed. That
	// is the honest shape of a pause (it is not a lock), and the case pushes until
	// the drop counter moves instead of pretending the boundary is instantaneous.
	for i := 0; i < 3 && rv.pausedDrops.Load() == 0; i++ {
		pushFrame("a frame until the pause is observed")
	}
	if rv.pausedDrops.Load() == 0 {
		t.Fatalf("three frames after the pause and none was dropped: the consumer is not observing the paused bit")
	}
	fedAtPause := rv.drained.Load()
	dropsAtPause := rv.pausedDrops.Load()

	// Paused: frames are drained and counted, but never fed.
	pushFrame("a second frame while paused")
	pushFrame("a third frame while paused")
	if got := rv.drained.Load(); got != fedAtPause {
		t.Fatalf("drained went %d -> %d while paused: a paused listener fed the chain", fedAtPause, got)
	}
	if got := rv.pausedDrops.Load(); got != dropsAtPause+2 {
		t.Fatalf("paused drops went %d -> %d, want +2", dropsAtPause, got)
	}

	// Resume: the listener takes audio again and the checkmark clears.
	sentence, executed = ball.currentWakePauseToggle()()
	if !executed || !strings.Contains(sentence, "已恢复") {
		t.Fatalf("the tray resume answered %q executed=%v", sentence, executed)
	}
	if !ball.mirrorTrayMute() {
		t.Fatalf("the tray projection wrote nothing after a resume")
	}
	if last := pushed.last(t); last.pausedWake {
		t.Fatalf("the tray checkmark still says paused after a resume: %+v", last)
	}
	// Resume, the same shape in reverse: whichever frame was already inside the
	// paused helper is dropped, and the next one is fed.
	for i := 0; i < 3 && rv.drained.Load() == fedAtPause; i++ {
		pushFrame("a frame until the resume is observed")
	}
	if got := rv.drained.Load(); got != fedAtPause+1 {
		t.Fatalf("drained=%d after the resume, want %d (exactly one frame fed, no more)", got, fedAtPause+1)
	}
	t.Logf("PAUSE drained=%d paused_drops=%d fed_before_pause=%d tray_writes=%d",
		rv.drained.Load(), rv.pausedDrops.Load(), fedAtPause, pushed.count())
}

// Test15eTheProductionHandOffIsTheModelStoreWalk is reading #1's "is it really
// the product's walk" half: residentVoiceHandOff is models.go's handOffModel over
// the signed manifest, and it announces the model available for the pinned bytes
// without touching the network.
func Test15eTheProductionHandOffIsTheModelStoreWalk(t *testing.T) {
	dataDir := voice15eDataDir(t)
	store := voice15eInstallAsr(t, dataDir)
	voice15eWriteConfig(t, dataDir, nil)
	t.Setenv("WISP_MODELS_MANIFEST", voice15eRepoFile(t, filepath.Join("build", "models", "manifest.json")))

	hook := residentVoiceHandOff(dataDir)
	var ticks int
	exit := hook(voiceAsrModelID, func(format string, args ...any) { ticks++ })
	if exit != 0 {
		t.Fatalf("residentVoiceHandOff(%s) = %d, want 0: the pinned bytes hash-match the signed manifest", voiceAsrModelID, exit)
	}
	for _, name := range []string{voiceAsrEncoder, voiceAsrDecoder, voiceAsrTokens} {
		if _, err := os.Stat(filepath.Join(store, name)); err != nil {
			t.Fatalf("the hand-off passed but %s is not in the store: %v", name, err)
		}
	}
	t.Logf("PRODUCTION HAND-OFF id=%s exit=%d progress_ticks=%d store=%s", voiceAsrModelID, exit, ticks, store)

	// The store is what the engine opens, and it is byte-identical to the pinned
	// copies: this is the reading that makes "opens <store>/<id>" safe to claim.
	for _, pair := range [][2]string{
		{voiceAsrEncoder, "asr-encoder.int8.onnx"},
		{voiceAsrDecoder, "asr-decoder.int8.onnx"},
		{voiceAsrTokens, "asr-tokens.txt"},
	} {
		installed, err := os.Stat(filepath.Join(store, pair[0]))
		if err != nil {
			t.Fatalf("stat %s: %v", pair[0], err)
		}
		pinned, err := os.Stat(filepath.Join(voice15eSpikeModels(t), pair[1]))
		if err != nil {
			t.Fatalf("stat %s: %v", pair[1], err)
		}
		if installed.Size() != pinned.Size() {
			t.Fatalf("%s installed size %d != pinned %d", pair[0], installed.Size(), pinned.Size())
		}
	}
}
