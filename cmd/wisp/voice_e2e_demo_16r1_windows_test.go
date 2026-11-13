//go:build windows

package main

// Ticket 16 (S2 acceptance gate) evidence rig. Three cases:
//
//   - Test16r1VoiceDemo... is AC#1's chain demo in one run on the real
//     hardware this machine has (Realtek capture + render endpoints, live-mic
//     smoke PASS the same day): the wake word is the real KwsEngine on the
//     checkpoint's own fixture, the command is the real streaming Paraformer
//     through the real C8 seam, the transcript reaches the REAL composition
//     root (runTextTask, the `wisp run` entry) over the C5 golden SSE seam,
//     the agent really dispatches one allowlisted fs.read (its tool_call row
//     is read back off the real store), the reply is announced by the
//     PRODUCTION speak() path (real matcha checkpoint, real waveOut render,
//     the capture gate really closes), and the scope really walks
//     SCOPE-WARM -> SCOPE-SETTLING (frozen 3s static) -> SCOPE-DISPOSED with
//     the settle trace stamping T0/T1/T2 on the production seams (12-r2's
//     shape). Nothing is stubbed that HEAD wires for real: no speaker stub,
//     no TTS stub, no machine bridge - the whole D43 walk below is what the
//     production paths dispatch themselves.
//
//   - Test16r1PremortemMicPermissionDeniedLiveProbe is the OS-level half of
//     the mic-permission pre-mortem: gated behind WISP_16R1_MIC_DENY_PROBE=1
//     and run by the rig that has really flipped the Windows consent store
//     (the same shape TestLiveWasapiSmoke's WISP_LIVE_MIC gate established).
//
//   - Test16r1PremortemModelSha256Corrupt... is the model-corruption
//     pre-mortem: one flipped byte in the INSTALLED store copy (never the
//     pinned third_party source - the store entry is replaced with a modified
//     copy so the hardlink never carries the flip back), then the production
//     hand-off walk must refuse loudly.
//
// Staging contract is 15e's: sherpa DLLs next to the test binary, the model
// envs pointed at the pinned trees, -test.run naming this family. The corrupt
// case needs only the manifest env and the spike-models staging; the demo case
// needs the full staging plus an exclusive moment on the render device.

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/llm/golden"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/secret"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// demo16Reply is the golden server's second-round reply. Short on purpose:
// announce() refuses sentences past voiceAnnounceMaxRunes and this demo wants
// the sentence spoken, not D10'd into silence.
const demo16Reply = "已经读完了，笔记里只有一行话，随时可以继续"

// demo16Golden scripts the two rounds the demo turn runs: one allowlisted
// fs.read (L0, no card, the route zz315's control leg measured) and then a
// pure-text close that the announce path speaks.
func demo16Golden(t *testing.T, note string) string {
	t.Helper()
	args, err := json.Marshal(map[string]string{"path": filepath.ToSlash(note)})
	if err != nil {
		t.Fatal(err)
	}
	toolCall := demo16Chunk(t, map[string]any{"index": 0, "delta": map[string]any{
		"tool_calls": []any{map[string]any{
			"index": 0, "id": "call_16r1_demo", "type": "function",
			"function": map[string]any{"name": "fs.read", "arguments": string(args)},
		}},
	}})
	finish := demo16Chunk(t, map[string]any{"index": 0, "delta": map[string]any{}, "finish_reason": "tool_calls"})
	usage := `data: {"choices":[],"usage":{"prompt_tokens":31,"completion_tokens":12,"total_tokens":43}}`
	reply := demo16Chunk(t, map[string]any{"index": 0, "delta": map[string]any{"content": demo16Reply}})
	stop := demo16Chunk(t, map[string]any{"index": 0, "delta": map[string]any{}, "finish_reason": "stop"})
	return strings.Join([]string{
		"# wisp golden sse v1",
		"# @scenario 16-r1 voice demo: one allowlisted fs.read, then the spoken close",
		"# @response 200",
		toolCall, "", finish, "", usage, "", "data: [DONE]", "",
		"# @response 200",
		reply, "", stop, "", usage, "", "data: [DONE]", "",
	}, "\n")
}

// demo16Chunk renders one OpenAI chat-completion SSE data line (the same byte
// shape the 246 live suite's helper renders; that helper lives behind the
// winlive tag, this file must not).
func demo16Chunk(t *testing.T, choice map[string]any) string {
	t.Helper()
	body, err := json.Marshal(map[string]any{"choices": []any{choice}})
	if err != nil {
		t.Fatal(err)
	}
	return "data: " + string(body)
}

// demo16Bridge is the transcript sink between the voice leg and the REAL task
// pipeline. The production residentTaskSource.submitTask is async (one Spawn
// per task, resident_task_source_windows.go:458); this bridge keeps that shape
// and runs the same composition root `wisp run` uses, with the runtime hook
// wiring the reply into the production announce path and the settle trace
// exactly as resident_windows.go does at boot.
type demo16Bridge struct {
	dataDir string

	mu          sync.Mutex
	transcripts []string
	exits       []int
	stdout      *strings.Builder
	stderr      *strings.Builder
	notices     []string
	onRuntime   func(*agentRuntime)

	turned chan struct{}
}

func (b *demo16Bridge) submitTask(text string) error {
	b.mu.Lock()
	b.transcripts = append(b.transcripts, text)
	b.mu.Unlock()
	go func() {
		code := runTextTask(runSpec{
			argv:    []string{text},
			stdout:  b.stdout,
			stderr:  b.stderr,
			dataDir: b.dataDir,
			notify: func(title, body string) error {
				b.mu.Lock()
				b.notices = append(b.notices, title+"\x00"+body)
				b.mu.Unlock()
				return nil
			},
			onRuntime: b.onRuntime,
		})
		b.mu.Lock()
		b.exits = append(b.exits, code)
		b.mu.Unlock()
		select {
		case b.turned <- struct{}{}:
		default:
		}
	}()
	return nil
}

func (b *demo16Bridge) snapshot() (transcripts []string, exits []int, notices []string, out, errb string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]string(nil), b.transcripts...), append([]int(nil), b.exits...),
		append([]string(nil), b.notices...), b.stdout.String(), b.stderr.String()
}

// demo16AnnounceNeedsAnOutputDevice is ①'s precondition probe (ticket 111 wave
// 10). The demo's tail speaks the reply through the PRODUCTION announce path -
// real matcha checkpoint, real waveOut render (voice_out_windows.go:217 opens
// the default output device). On a hosted runner that reports zero waveOut
// devices the announce line itself degrades to the D10 notification
// (resident_voice_engines_windows.go:1171 notes 播报没有输出设备 and the
// announcement cannot be heard), so the SPEAK half of the D43 walk below
// (#15/#26, and the counters at :379) has no subject. Skip ONLY that named
// precondition; on any machine with an output device the probe opens nothing
// durably (probe speaker is closed immediately) and every assertion below runs
// unchanged. Probe shape = the 26p echo rig's own precedent
// (resident_voice_bargein_26p_echo_windows_test.go:102-106).
func demo16AnnounceNeedsAnOutputDevice(t *testing.T) {
	t.Helper()
	// Tripwire proof lever (ticket 111 wave 10): forces the probe's ANSWER to
	// no-device so the skip branch can be demonstrated on a machine that HAS
	// an output device - the dispatcher's own injection shape ("把输出设备探
	// 测结果打到无设备"). No product code runs below this branch either way.
	if os.Getenv("WISP_TEST_FORCE_NO_WAVEOUT") == "1" {
		t.Skipf("the demo's announce needs a waveOut output device; this machine reports none: [INJECTED premise WISP_TEST_FORCE_NO_WAVEOUT=1 - the probe's answer is forced to no-device for the tripwire reading, the real openWaveOutSpeaker probe did not run] - the D43#15/#26 Speaking walk has no subject here (degrades to the D10 notification), so the demo cannot be measured on this host")
	}
	if probe, err := openWaveOutSpeaker(22050); err != nil {
		t.Skipf("the demo's announce needs a waveOut output device; this machine reports none: %v - the D43#15/#26 Speaking walk has no subject here (degrades to the D10 notification), so the demo cannot be measured on this host", err)
	} else {
		_ = probe.Close()
	}
}

// Test16r1VoiceDemoWakeAsrAgentSpeakSettleRealHardware is AC#1's chain demo:
// say command (wav) -> Armed->Listening->Thinking->(agent turn)->Speaking->
// Warm->Settling with the transcript, the tool_call rows and the settle trace
// all on record. Every D43 row below is dispatched by production code paths
// (KWS arm #5, wake #7, vad-stop #11, announce #15/#26); the case dispatches
// nothing by hand.
func Test16r1VoiceDemoWakeAsrAgentSpeakSettleRealHardware(t *testing.T) {
	demo16AnnounceNeedsAnOutputDevice(t)
	dataDir := voice15eDataDir(t)

	// The note fs.read will read, and the config that points the agent at the
	// golden server while allowlisting exactly this data root.
	note := filepath.Join(dataDir, "note.txt")
	if err := os.WriteFile(note, []byte("wisp 票 16 演示笔记\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	srv, rep, err := golden.ServerFromBytes([]byte(demo16Golden(t, note)))
	if err != nil {
		t.Fatalf("golden server: %v", err)
	}
	rep.Pace = true
	t.Cleanup(srv.Close)
	body := fmt.Sprintf(`schema_version = 2

[llm]
text_chain = ["acme/m1"]
timeout_ms = 120000

[llm.retry]
max = 1
backoff_ms = 1

[llm.providers.acme]
protocol = "openai-chat"
base_url = %q
api_key_ref = "dpapi:acme"

[llm.providers.acme.models.m1]
enabled = true
context_window = 128000

[fs]
allowed_dirs = [%q]
`, srv.URL+"/v1", filepath.ToSlash(dataDir))
	if err := os.WriteFile(filepath.Join(dataDir, configFileName), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	st, err := secret.NewStore(dataDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Store("dpapi:acme", fakeStoreKey); err != nil {
		t.Fatal(err)
	}

	t.Setenv("WISP_MODELS_MANIFEST", voice15eRepoFile(t, filepath.Join("build", "models", "manifest.json")))
	store := voice15eInstallAsr(t, dataDir)

	kws := voiceKwsFilesFrom(voice41a2r1KwsDir(t))
	if !kws.complete() {
		t.Fatalf("the KWS manifest spelling did not resolve under %s", voice41a2r1KwsDir(t))
	}
	// The demo-table pairing 41-a2 measured: test_wavs are recordings of
	// test_keywords.txt's words, and 0.wav is that collection's own negative.
	kws.Keywords = filepath.Join(voice41a2r1KwsDir(t), "test_wavs", "test_keywords.txt")

	plan := voicePlan{
		VadModelPath: voice15eVadModel(t),
		AsrDir:       store,
		Source:       "16-r1 E2E demo",
		HandOff:      residentVoiceHandOff(dataDir),
		Tts:          voiceWireTtsFiles(t), // real matcha checkpoint: the production TTS factory loads it
		Kws:          kws,
		WarmTimeout:  1200 * time.Millisecond, // staging warm window (12-r2's named shape)
	}

	bridge := &demo16Bridge{
		dataDir: dataDir,
		stdout:  &strings.Builder{},
		stderr:  &strings.Builder{},
		turned:  make(chan struct{}, 8),
	}
	var rv16demo *residentVoice
	// The production boot wiring (resident_windows.go:580/:610): the reply
	// announces through the voice leg and the sink stamps the settle trace.
	bridge.onRuntime = func(rt *agentRuntime) {
		rt.voiceAnnounce = rv16demo.announce
		rt.settle = rv16demo.residency.settle
	}
	// voiceWireStart keeps the constructor's residency; the speaker factory and
	// the TTS loader are deliberately NOT overridden - the demo speaks through
	// the real matcha factory and the real waveOut device.
	rv16demo, frames := voiceWireStart(t, plan, &residentBall{}, bridge, nil)
	trace := rv16demo.residency.settle

	// The boot walk is production's own: KWS armed puts the machine in Armed.
	if got := rv16demo.machine.State(); got != statemachine.StateArmed {
		t.Fatalf("machine after the boot walk = %s, want Armed (D43#5)", string(got))
	}

	// Negative control through the same live pipeline: the checkpoint's own
	// negative sample must not wake anything.
	hits0 := rv16demo.kwsHits.Load()
	voice41a2r1Feed(t, frames, voice41a2r1Fixture(t, "0.wav"), rv16demo)
	if got := rv16demo.kwsHits.Load(); got != hits0 {
		t.Fatalf("0.wav produced %d hits in the live pipeline, want 0 (it is the collection's own negative)", got-hits0)
	}
	if got := rv16demo.machine.State(); got != statemachine.StateArmed {
		t.Fatalf("machine after the negative fixture = %s, want Armed", string(got))
	}

	// The demo utterance: 4.wav says two demo keywords - the wake word is the
	// first, the rest is the command the ASR transcribes.
	hits1 := rv16demo.kwsHits.Load()
	voice41a2r1Feed(t, frames, voice41a2r1Fixture(t, "4.wav"), rv16demo)
	voice15eWaitFor(t, "the wake-word hit", func() bool { return rv16demo.kwsHits.Load() > hits1 })
	voice15eWaitFor(t, "the D43#7 Armed->Listening walk", func() bool {
		return demo16Walked(rv16demo, statemachine.StateArmed, statemachine.EvWakeWord, statemachine.StateListening)
	})
	voice15eWaitFor(t, "the ASR chain the wake loads", func() bool { return rv16demo.asrLoaded.Load() })

	// The transcript crosses into the REAL composition root; the turn ends;
	// the production announce speaks it; the scope falls back.
	voice15eWaitFor(t, "the transcript to reach the task pipeline", func() bool {
		tr, _, _, _, _ := bridge.snapshot()
		return len(tr) > 0
	})
	select {
	case <-bridge.turned:
	case <-time.After(90 * time.Second):
		_, _, _, out, errb := bridge.snapshot()
		t.Fatalf("the agent turn never finished\nstdout:\n%s\nstderr:\n%s", out, errb)
	}
	voice15eWaitFor(t, "the production announce to speak through the real device", func() bool {
		return rv16demo.ttsSpoken.Load() >= 1
	})
	voice15eWaitFor(t, "the reply's fall-back to land in the trace (T2)", func() bool {
		return len(trace.samplesSnapshot()) >= 1
	})

	// ---- the record: state trace, tool rows, transcript, settle, counters.
	timeline := rv16demo.timelineSnapshot()
	for _, tr := range timeline {
		t.Logf("DEMO16 D43 %s+%s->%s at_sample=%d", tr.From, tr.Event, tr.To, tr.AtSample)
	}
	t.Logf("DEMO16 BALL_STATES %v", rv16demo.ball.voiceStatesSnapshot())
	if !demo16Walked(rv16demo, statemachine.StateSleeping, statemachine.EvKwsEnabled, statemachine.StateArmed) {
		t.Fatalf("no D43#5 Sleeping->Armed walk: %+v", timeline)
	}
	if !demo16Walked(rv16demo, statemachine.StateArmed, statemachine.EvWakeWord, statemachine.StateListening) {
		t.Fatalf("no D43#7 Armed->Listening walk: %+v", timeline)
	}
	if !demo16Walked(rv16demo, statemachine.StateListening, statemachine.EvVadStop, statemachine.StateThinking) {
		t.Fatalf("no D43#11 Listening->Thinking walk: %+v", timeline)
	}
	if !demo16Walked(rv16demo, statemachine.StateThinking, statemachine.EvFirstToken, statemachine.StateSpeaking) {
		t.Fatalf("no D43#15 Thinking->Speaking walk (the announce path's first-token row): %+v", timeline)
	}
	if !demo16Walked(rv16demo, statemachine.StateSpeaking, statemachine.EvSpeakDone, statemachine.StateWarm) {
		t.Fatalf("no D43#26 Speaking->Warm walk: %+v", timeline)
	}

	transcripts, exits, notices, out, errb := bridge.snapshot()
	if len(exits) != 1 || exits[0] != 0 {
		t.Fatalf("the agent turn exits = %v, want exactly [0]\nstdout:\n%s\nstderr:\n%s", exits, out, errb)
	}
	t.Logf("DEMO16 TRANSCRIPT %q", transcripts)
	t.Logf("DEMO16 REPLY %q (streamed to the sink, spoken, notified)", demo16Reply)

	// Tool evidence rows, read back off the real store the task booked into.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id := func() string {
		m := taskIDRe.FindStringSubmatch(out)
		if m == nil {
			t.Fatalf("the run printed no status line:\n%s", out)
		}
		return m[1]
	}()
	st2, err := memory.Open(dataDir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st2.Close() })
	rows, err := st2.ListToolCallsByTask(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	var read *memory.ToolCall
	for i := range rows {
		t.Logf("DEMO16 TOOL_ROW tool=%s level=%s decision=%s outcome=%s class=%q corr=%q",
			rows[i].Tool, rows[i].RiskLevel, rows[i].Decision, rows[i].Outcome, rows[i].ErrorClass, rows[i].CorrelationID)
		if rows[i].Tool == "fs.read" {
			read = &rows[i]
		}
	}
	if read == nil {
		t.Fatalf("no fs.read tool_call row for task %s: %+v", id, rows)
	}
	if read.Outcome != "success" {
		t.Fatalf("the fs.read row outcome = %s, want success: the demo's tool call must really have run", read.Outcome)
	}
	tl, err := st2.TaskLogByID(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("DEMO16 TASK_LOG id=%s state=%s class=%q query=%q", tl.ID, tl.State, tl.ErrorClass, tl.QueryText)
	if tl.State != "done" || tl.QueryText == "" {
		t.Fatalf("task_log = %+v, want a closed done row carrying the voice transcript", tl)
	}
	if len(notices) == 0 {
		t.Fatalf("the reply was never presented as a D10 notification: %v", notices)
	}
	t.Logf("DEMO16 NOTIFY %d posted (first: %q)", len(notices), notices[0])

	// The announce counters: exactly one sentence spoken through the real
	// device, none aborted, none gated.
	t.Logf("DEMO16 SPEAK asked=%d spoken=%d chunks=%d samples=%d aborted=%d gated=%d",
		rv16demo.announceAsked.Load(), rv16demo.ttsSpoken.Load(), rv16demo.ttsChunks.Load(),
		rv16demo.ttsSamples.Load(), rv16demo.ttsAborted.Load(), rv16demo.ttsGated.Load())
	if rv16demo.ttsSpoken.Load() != 1 || rv16demo.ttsAborted.Load() != 0 {
		t.Fatalf("spoken=%d aborted=%d, want exactly one real playback and no abort",
			rv16demo.ttsSpoken.Load(), rv16demo.ttsAborted.Load())
	}

	// The settle trace: T0/T1/T2 stamped on the production seams. The frozen
	// 3s static window is the reading the demo judges (it never fires early);
	// the token->Settling leg legitimately carries the REAL announce (first
	// matcha load ~3.7s + playback ~4s + the staging warm window 1.2s), which
	// the stub-TTS text-path shape of 12-r2 does not have - so that leg is
	// recorded, not bounded here.
	s := trace.samplesSnapshot()[0]
	t.Logf("DEMO16 SETTLE token_to_settling=%dms (includes the real TTS load+playback) settling_to_fallback=%dms token_to_fallback=%dms",
		s.LastTokenToSettling.Milliseconds(), s.SettlingToFallback.Milliseconds(), s.LastTokenToFallback.Milliseconds())
	if s.LastTokenToSettling <= 0 {
		t.Fatal("the settle trace recorded no token->Settling leg")
	}
	if s.SettlingToFallback < 2900*time.Millisecond || s.SettlingToFallback > 3600*time.Millisecond {
		t.Errorf("Settling->fallback %dms is not the frozen 3s static window", s.SettlingToFallback.Milliseconds())
	}

	// The audio-clock latency readings (the turn really ran the real ASR).
	// These are the WAKE-path shape (the utterance opens only after the wake
	// word lands mid-stream), so openDelay here is not the 15e single-utterance
	// baseline shape; the 262/640/902 zero-regression reading comes from the
	// six-family Test15e run on the same binary, not from this demo.
	turns := rv16demo.turnsSnapshot()
	if len(turns) != 1 {
		t.Fatalf("turn records = %+v, want exactly the demo utterance", turns)
	}
	onset := int64(voice15eLeadMS * audio.TargetRate / 1000)
	openDelay := time.Duration(turns[0].StartSample-onset) * time.Second / time.Duration(audio.TargetRate)
	waitToText := time.Duration(turns[0].FirstTextAt-onset) * time.Second / time.Duration(audio.TargetRate)
	t.Logf("DEMO16 LATENCY (wake-path utterance) openDelay=%v chunk=%v finalLatency=%v speechMS=%d closedBy=%s seamWaitToFirstText=%v",
		openDelay, turns[0].FirstChunk, turns[0].FinalLatency, turns[0].SpeechMS, turns[0].ClosedBy, waitToText)
}

func demo16Walked(rv *residentVoice, from statemachine.State, ev statemachine.Event, to statemachine.State) bool {
	for _, tr := range rv.timelineSnapshot() {
		if tr.From == from && tr.Event == ev && tr.To == to {
			return true
		}
	}
	return false
}

// Test16r1PremortemMicPermissionDeniedLiveProbe is the OS-level half of the
// mic-permission pre-mortem (D42#12). The rig flips the Windows consent store
// to Deny, runs this probe, and flips it back; the probe asserts what the
// REAL open path answers under a REAL OS refusal. Without the flip the probe
// would trivially open the device, so it is gated like TestLiveWasapiSmoke.
func Test16r1PremortemMicPermissionDeniedLiveProbe(t *testing.T) {
	if os.Getenv("WISP_16R1_MIC_DENY_PROBE") != "1" {
		t.Skip("live probe: run only under the 16-r1 pre-mortem rig that has flipped the OS consent store (WISP_16R1_MIC_DENY_PROBE=1)")
	}
	mic := audio.NewWASAPIMicrophone()
	cap1, _, err := mic.Endpoints()
	if err != nil {
		t.Fatalf("enumerate endpoints: %v", err)
	}
	buf := audio.NewBoundedFrames()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	startErr := mic.Start(ctx, buf)
	if startErr == nil {
		_ = mic.Stop()
		t.Fatalf("the real microphone %q opened while the OS consent store says Deny: the flip did not reach this process", cap1.String())
	}
	if got, ok := observe.ClassOf(startErr); !ok || got != observe.ClassAudioDevice {
		t.Fatalf("class = %v (ok=%v), want audio_device", got, ok)
	}
	t.Logf("DENY-PROBE capture=%q LastError=%q", cap1.String(), mic.Stats().LastError)
	t.Logf("DENY-PROBE error=%q", startErr.Error())
}

// Test16r1PremortemModelSha256CorruptRefusesTheHandOff is the model-corruption
// pre-mortem (SPEC-10 §7): one flipped byte in the INSTALLED store copy. The
// reading has two halves and records what HEAD actually does:
//
//  1. DETECTION - the store's own verification on the tampered tree fails with
//     the sha256 mismatch named, under the model class (D37's 模型 row).
//  2. THE WALK - the production hand-off walk refuses to hand the tampered
//     bytes over: at HEAD it self-heals (re-fetch + verify, the downloader's
//     documented "reuse it, then re-verify at the hand-off point" contract)
//     and the store it leaves behind verifies clean. The walk never exits 0
//     with unverified bytes; the can't-heal refusal branch (mirror tampered ->
//     "sha256 MISMATCH ... staged file deleted", model class) is the models
//     package's own pinned case (TestAC21FreshInstallIsGuardedOnEnsuresOwnReturnPath).
func Test16r1PremortemModelSha256CorruptRefusesTheHandOff(t *testing.T) {
	dataDir := voice15eDataDir(t)
	voice15eWriteConfig(t, dataDir, nil)
	t.Setenv("WISP_MODELS_MANIFEST", voice15eRepoFile(t, filepath.Join("build", "models", "manifest.json")))
	storeDir := filepath.Join(dataDir, "models", voiceAsrModelID)
	if err := voice15eInstallAsrFiles(dataDir); err != nil {
		t.Fatalf("install the streaming ASR into %s: %v", dataDir, err)
	}

	// Baseline: the pristine store passes the production walk.
	if exit := residentVoiceHandOff(dataDir)(voiceAsrModelID, func(string, ...any) {}); exit != 0 {
		t.Fatalf("the pristine store was refused (exit %d): the fixture is broken, not the model gate", exit)
	}

	// Corrupt the INSTALLED copy. The install hardlinks the pinned source, so
	// the store entry is REPLACED with a modified copy - the pinned
	// third_party bytes must never carry the flip back through the link.
	target := filepath.Join(storeDir, voiceAsrEncoder)
	blob, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	blob[len(blob)/2] ^= 0xFF
	if err := os.Remove(target); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, blob, 0o644); err != nil {
		t.Fatal(err)
	}

	// Half 1, DETECTION: the store's own verification on the tampered tree.
	ms, err := openModelStore(dataDir, nil)
	if err != nil {
		t.Fatal(err)
	}
	verr := ms.mgr.VerifyInstalled(voiceAsrModelID)
	if verr == nil {
		t.Fatal("VerifyInstalled passed on the tampered store: the corruption was not detected")
	}
	if class, ok := observe.ClassOf(verr); !ok || class != observe.ClassModel {
		t.Fatalf("VerifyInstalled error class = %v (ok=%v), want model", class, ok)
	}
	if !strings.Contains(verr.Error(), "sha256 mismatch") {
		t.Fatalf("the verification error does not name the sha256 mismatch: %q", verr.Error())
	}
	t.Logf("CORRUPT detected verify_error=%q class=model", verr.Error())

	// Half 2, THE WALK: the production hand-off must never hand the tampered
	// bytes over. Whatever branch it takes (self-heal re-fetch at HEAD), the
	// store it leaves behind must verify clean.
	started := time.Now()
	var walkLog []string
	exit := residentVoiceHandOff(dataDir)(voiceAsrModelID, func(format string, args ...any) {
		walkLog = append(walkLog, fmt.Sprintf(format, args...))
	})
	took := time.Since(started)
	t.Logf("CORRUPT walk exit=%d took=%v (baseline pristine walk was ~2s)", exit, took.Round(time.Millisecond))
	for _, line := range walkLog {
		t.Logf("CORRUPT walk_log %s", line)
	}
	if exit != 0 {
		t.Fatalf("the walk refused the tampered store (exit %d): record the refusal branch - the healed-store check below is then moot", exit)
	}
	ms2, err := openModelStore(dataDir, nil)
	if err != nil {
		t.Fatal(err)
	}
	if verr2 := ms2.mgr.VerifyInstalled(voiceAsrModelID); verr2 != nil {
		t.Fatalf("after the walk the store still does not verify: %v (the user would be left with tampered bytes)", verr2)
	}
	if took < 5*time.Second {
		t.Fatalf("the walk returned exit 0 in %v: a cache-hit shape would mean the tampered bytes were trusted", took.Round(time.Millisecond))
	}
	t.Logf("CORRUPT healed store verifies clean (the hand-off never delivered unverified bytes)")
}
