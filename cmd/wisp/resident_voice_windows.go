//go:build windows

package main

// Ticket 15e: the voice chain, wired into the resident process. This is the
// third and last of the three seams `15-a1` found sitting unowned in cmd/wisp,
// and it is the only file in this tree where `internal/speech` is imported.
//
// WHAT IT CLOSES, one seam at a time, each named where the old comment said it
// could not be closed from the package that owns it:
//
//  1. models.go:25 - "Whichever ticket lands internal/speech owns adding its own
//     handOffModel call; this file cannot become that reader by wishing."
//     handOffModel is a method on *modelStore in package main (models.go:357),
//     so internal/speech can never call it: the call site has to be here.
//     residentVoiceHandOff below is that call site, and the ordering it creates
//     is the one models.go's AC#5 paragraph asks for - the hand-off walk
//     returns, and only then are the model's bytes opened.
//
//  2. resident_audio_windows.go:57-63 - "The C8 frame channel is created with
//     the seam's own bounded window and nothing drains it." The channel is now
//     handed over (residentAudio.frames) and drainVoiceFrames below is its one
//     consumer. The "nothing drains it" sentence in that file is corrected in
//     the same commit; a comment that says a hole is open when it is closed is
//     the defect ticket 121 shipped this file's own header about.
//
//  3. resident_ball_windows.go:522,596 - "the wake-word listener it would pause
//     lives in internal/speech, which is still only a doc.go" and "the pause-wake
//     half has NOTHING to read yet". The listener now exists (this file's
//     consumer), so the tray's pause-wake item gets an executor and the tray's
//     second checkmark gets a truth source: both are attached from here
//     (attachWakePause) and both answer "no listener" when this leg did not
//     build one.
//
// Plus AC#1's last clause, which needs all three: the ball walks
// Listening -> Thinking. That walk is the D43 table's rows #4 and #11, driven by
// the endpoint events this file consumes - not by a state this file sets
// directly. Row #4 is Sleeping + EvSummon -> Listening and it is the ONLY legal
// way into Listening (row #7 needs a wake word, ticket 41, and the summon hot
// key has no executor in this process), so "the user spoke while the microphone
// was open" is dispatched as the summon. That is a wiring decision, not a table
// change: when ticket 41 lands KWS this dispatch becomes EvWakeWord and the
// table stays byte-identical.
//
// WHY THE HAND-OFF IS GATED ON SPEECH AND NOT ON BOOT (the shape everything
// else here follows). SPEC-04 sec 7.1 says a model is fetched on the FIRST USE
// OF ITS CAPABILITY. A microphone that is open all day must not fetch a 237 MB
// ASR pair because the process started, and it must not do it because silence
// arrived either: silence never opens an utterance (15-b's state machine), so
// silence never reaches the capability. The chain therefore runs in two phases:
//
//	phase 1  a warm silero Detector drains the seam and answers "did an
//	         utterance just open?" - 643 KB, part of listening, no hand-off;
//	phase 2  on the FIRST EndpointStart: D43#4, then the hand-off walk, then
//	         (only on exit 0) the ASR bytes are opened from <store>/<id>/ and the
//	         real speech.Transcriber takes over the same audio.
//
// The switchover replays the audio the ring kept, so the utterance that
// triggered the hand-off is not clipped: phase 2's fresh Detector re-derives the
// same EndpointStart from that contiguous history and the Transcriber replays
// the speech run into the ASR exactly as internal/speech documents. The ring is
// bounded (voiceRingSeconds) and only lives until the first speech, because after
// that the Transcriber owns the audio.
//
// WHAT THIS FILE DOES NOT CLAIM.
//   - Not that the owner hears anything on the first utterance: the hand-off
//     walk hashes the installed model files on its way to the announcement, and
//     C8 frames that arrive during it are dropped by the bounded channel and
//     counted by the D38d meter (audio frames dropped: bounded channel full).
//     The replay ring covers the utterance that triggered the walk, not the
//     duration of the walk. Preloading at boot would remove the stall and cost
//     the "silence never fetches a model" property; that trade is the owner's,
//     not this leg's, and it is written down rather than hidden.
//   - Not that the ASR is usable when the hand-off was refused: exit != 0 means
//     the bytes are NOT opened, the utterance is dropped, and the refusal is
//     said out loud with its class (models.go already prints it).
//   - Not that AC#5's retry policy is implemented: an empty final is recorded
//     as 没听清 and not handed to the agent, which is the semantic AC#5 keys on,
//     but "one retry then an asr error class" is not this slice's (15-c said the
//     same).
//   - Not that the D43 Listening timeout is projected: the machine's own timer
//     (row #12) moves the machine without a ball hop, because the timer fires
//     inside internal/statemachine where this leg has no hook. The next voice
//     event projects the machine's then-current state. Named as a residual in
//     15e-r1's evidence table rather than papered over.

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/plugin"
	"github.com/CarlosShao/wisp/internal/proc"
	"github.com/CarlosShao/wisp/internal/speech"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// The signed manifest's ids and the artifact names INSIDE an installed model
// directory, spelled exactly as build/models/manifest.json spells them: the
// store's layout is <data dir>\models\<id>\<file> (SPEC-04 sec 7.2), and the
// dev copies under third_party\spike-models carry different names
// (asr-encoder.int8.onnx) which is why this leg cannot borrow those.
const (
	// voiceAsrModelID is purpose "asr-streaming" in the manifest.
	voiceAsrModelID = "asr-streaming-paraformer-zh-en"
	// voiceVadModelID is purpose "vad" in the manifest (a single-artifact
	// entry, so its install file is the base name of its one URL).
	voiceVadModelID = "vad-silero"

	voiceVadFile    = "silero_vad.onnx"
	voiceAsrEncoder = "encoder.int8.onnx"
	voiceAsrDecoder = "decoder.int8.onnx"
	voiceAsrTokens  = "tokens.txt"
)

// voiceRingSeconds is how much audio the pre-hand-off ring keeps: the contiguous
// history phase 2's fresh Detector replays so the utterance that opened is not
// clipped. It is measured, not picked: 15-c's fixture opens its utterance ~1.1 s
// into the speech (silero's onset overshoot), and the first use of the capability
// also pays ① the hand-off walk (hashing 237 MB of installed files - 2.3 s on this
// machine, worst case) and ② the engine load (another ~2.3 s). Eight seconds keeps
// the onset and the tail in one contiguous buffer across both, so the ring - not
// the transport channel - is what the switchover replays. A first-ever DOWNLOAD is
// longer than any ring, and the degradation there is said out loud instead of
// hidden (see this file's header).
const voiceRingSeconds = 8

// voiceRingSamples is voiceRingSeconds on the seam's own clock.
const voiceRingSamples = voiceRingSeconds * audio.TargetRate

// voiceGoroutineName is the D38b roster entry this leg's thread is booked under.
// It is NOT a new name and NOT a resident one: "asr-infer" is the frozen
// temporary-roster name for exactly this work (internal/observe/goroutine.go's
// TemporaryNames), scoped to the voice engine's lifetime through the C11
// DisposalScope below, which is what makes it released rather than resident.
const voiceGoroutineName = "asr-infer"

// voiceHandOffGoroutine is the roster entry the hand-off walk runs on. Like
// asr-infer it is a frozen temporary-roster name (internal/observe/goroutine.go's
// TemporaryNames), and it is the RIGHT one: this is a model operation, and it must
// not borrow the resident roster or the consumer's name.
const voiceHandOffGoroutine = "model-downloader"

// errVoiceChainRefused is openChain's "the walk said no" answer. It is not
// printed anywhere: the sentence that matters (which model, which exit code, and
// where the classified reason is) went through noteProblem already, and a second
// spelling would drift from it.
var errVoiceChainRefused = errors.New("voice: model hand-off refused")

// voiceChainReady is what the model thread hands the consumer once the first
// use has been paid for: the exit code of the hand-off walk, and - only when the
// walk said yes - the engines it opened. Fields are nil on the branches that
// failed; exit is the walk's own verdict, kept even on success so the reading
// "exit=0" is the walk's, not an inference from a nil error.
type voiceChainReady struct {
	exit int
	asr  *speech.StreamingAsr
	det  *speech.Detector
	err  error
}

// voiceHandOff is models.go's handOffModel, narrowed to the two facts this leg
// needs: which model id to make trustworthy, and the exit code that decides
// whether its bytes may be opened. 0 = the walk announced the model available
// (models.go's own code contract: 0 happened, 1 refused, 2 Unconfigured).
type voiceHandOff func(id string, logf func(string, ...any)) int

// voiceTranscriptSink is ticket 12's text path as the resident leg exposes it:
// residentTaskSource.submitTask (resident_task_source_windows.go:458). The
// interface is one method wide on purpose - the voice leg may start a task and
// may not touch anything else the task pipeline owns.
type voiceTranscriptSink interface {
	submitTask(text string) error
}

// voicePlan is what the assembly resolved for one process, kept as one value so
// the boot report, the tests and the running chain all read the same decision.
type voicePlan struct {
	// VadModelPath is the silero file phase 1 opens. It is resolved WITHOUT a
	// hand-off on purpose: VAD is listening infrastructure (643 KB, MIT), and
	// fetching it at boot is the one way this leg could make silence call the
	// hand-off.
	VadModelPath string
	// AsrDir is <store>\<asr id>: the directory the hand-off makes trustworthy
	// and the engine opens only after it returned 0.
	AsrDir string
	// HandOff is the call site models.go:25 reserved for this ticket.
	HandOff voiceHandOff
	// Source names which config table produced the paths ("...\config.toml" or
	// the compiled default table), so a boot line cannot imply a file it did not
	// read.
	Source string
}

// voiceTransition is one D43 edge this leg actually walked. AtSample is the
// AUDIO clock (samples drained from the C8 seam), never a wall-clock delta
// (D42#9): the same fixture reports the same timeline on any machine.
type voiceTransition struct {
	From     statemachine.State
	To       statemachine.State
	Event    statemachine.Event
	AtSample int64
	Say      string
}

// voiceTurn is one utterance, from the seam's point of view.
//
// StartSample is the utterance's first sample on the seam clock, derived the way
// internal/speech derives it (the seam position at the opening endpoint minus
// the run that opened it) - so FirstText-SampleStart is the ball-latency number
// AC#1 asks about, and FirstChunk is the same distance measured by the ASR's own
// audio clock (Transcript.Latency). Both are audio-clock readings.
type voiceTurn struct {
	OpenSample   int64
	StartSample  int64
	FirstTextAt  int64
	FirstChunk   time.Duration
	FinalAt      int64
	FinalLatency time.Duration
	SpeechMS     int
	ClosedBy     string
	Text         string
	Handed       bool
}

// residentVoice is this process's hold on the voice chain: the consumer, the
// D43 machine that decides what the ball shows, and the counters that make every
// claim in the boot report a reading. It is nil-safe everywhere, like
// residentAudio and residentBall.
type residentVoice struct {
	frames <-chan []byte
	plan   voicePlan
	ball   *residentBall
	tasks  voiceTranscriptSink

	machine *statemachine.Machine
	scope   *plugin.DisposalScope
	reg     *observe.Registry

	paused atomic.Bool

	// handOffStarted and handOffDone are the first-use trigger's own state, owned
	// by the consumer goroutine: the walk AND the engine load run on the roster's
	// model thread while the loop keeps draining, and the buffered channel is how
	// the ready engines (or the reason there are none) come back.
	handOffStarted bool
	chainReady     chan voiceChainReady

	// Counters, all atomic: the asr-infer goroutine writes them, the boot report
	// and the tray gesture read them.
	drained     atomic.Uint64
	pausedDrops atomic.Uint64
	partials    atomic.Uint64
	finals      atomic.Uint64
	discarded   atomic.Uint64
	handed      atomic.Uint64
	handOffs    atomic.Uint64
	handOffExit atomic.Int64
	asrLoaded   atomic.Bool
	fed         atomic.Int64

	mu       sync.Mutex
	timeline []voiceTransition
	turns    []voiceTurn
	problems []string
	askedFor []string // every model id this leg asked the store to hand over
	verdict  string
	stopped  bool
}

// voiceSamples converts one C8 frame to the float32 silero and the ASR want;
// int16 full scale maps to +-1.0, which is the same scale internal/speech's own
// pcmToFloat uses. That function is unexported and this slice may not change
// internal/speech, so the conversion is written once here, in the one place in
// cmd/wisp that crosses this seam - and a test pins the two implementations to
// the same arithmetic rather than trusting the comment.
func voiceSamples(pcm []int16) []float32 {
	out := make([]float32, len(pcm))
	for i, s := range pcm {
		out[i] = float32(s) / 32768
	}
	return out
}

// voiceRingAppend appends a window to the replay ring and keeps it bounded.
func voiceRingAppend(ring, samples []float32) []float32 {
	ring = append(ring, samples...)
	if len(ring) > voiceRingSamples {
		n := copy(ring, ring[len(ring)-voiceRingSamples:])
		ring = ring[:n]
	}
	return ring
}

// newResidentVoice builds the handle with everything but the thread. It is the
// seam the evidence table's end-to-end case drives, which is why the frame
// channel, the plan and the two sinks are parameters instead of being read from
// the process.
func newResidentVoice(frames <-chan []byte, plan voicePlan, ball *residentBall, tasks voiceTranscriptSink) *residentVoice {
	return &residentVoice{
		frames:  frames,
		plan:    plan,
		ball:    ball,
		tasks:   tasks,
		machine: statemachine.New(statemachine.Options{Initial: statemachine.StateSleeping}),
	}
}

// start warms the VAD, then books the consumer on the D38b roster through the
// C11 scope. A voice leg that has no frame channel (no capture leg) or no plan is
// not started at all, and the boot report prints the verdict that says so.
//
// The VAD session is created HERE, not on the consumer thread, and that ordering
// is measured rather than tidy: building the silero session costs ~200 ms, which
// is longer than the seam's 192 ms bounded window, so a consumer that built it
// inside its own loop would drop the first frames of whatever the microphone was
// already saying - and a dropped frame before a speech onset moves every number
// the ball-latency reading is built from. A VAD that will not build is therefore a
// boot-time refusal with a sentence, not a silent thread that exits.
func (rv *residentVoice) start(reg *observe.Registry) {
	if rv == nil || rv.frames == nil {
		return
	}
	if rv.plan.VadModelPath == "" {
		rv.verdict = "语音链未构造：VAD 模型路径为空"
		return
	}
	phase1, err := speech.NewDetector(speech.DetectorConfig{ModelPath: rv.plan.VadModelPath})
	if err != nil {
		rv.verdict = "语音链未构造：VAD 没起来（" + err.Error() + "）"
		rv.noteProblem(rv.verdict)
		return
	}
	rv.reg = reg
	opts := []plugin.Option{}
	if reg != nil {
		opts = append(opts, plugin.WithRegistry(reg))
	}
	// The scope is the engine's lifetime (C11): Dispose cancels the consumer's
	// ctx, joins it, and the engines themselves are released by the defers
	// inside run - cancel, then join, then release, which is the mandatory tail
	// order.
	rv.scope = plugin.NewDisposalScope("voice-15e", nil, opts...)
	rv.scope.Go(voiceGoroutineName, func(ctx context.Context) { rv.run(ctx, phase1) })
	slog.Info("voice: consumer booked on the D38b roster",
		"name", voiceGoroutineName, "category", "temporary (C11 scope)", "ring_seconds", voiceRingSeconds,
		"vad", rv.plan.VadModelPath, "asr_dir", rv.plan.AsrDir, "config_source", rv.plan.Source)
}

// stop is D38(e) step 4's second half: the capture leg's own hook calls this
// after the gate is closed. Dispose cancels the consumer, joins it (3 s budget)
// and only then are the engines released - a session released under a thread
// that is still decoding is the crash this order exists to prevent.
func (rv *residentVoice) stop() {
	if rv == nil {
		return
	}
	rv.mu.Lock()
	if rv.stopped {
		rv.mu.Unlock()
		return
	}
	rv.stopped = true
	rv.mu.Unlock()

	if rv.scope != nil {
		res := rv.scope.Dispose()
		if res.Incomplete || res.TimedOut {
			slog.Warn("voice: disposal did not complete cleanly", "scope", res.Scope, "timed_out", res.TimedOut)
		}
	}
	if rv.machine != nil {
		rv.machine.Close()
	}
	line := rv.stopLine()
	slog.Info(line)
	fmt.Printf("wisp: %s\n", line)
}

// run is the consumer thread: it drains the C8 seam for the whole life of the
// scope. Phase 1 (a warm VAD) lasts until the first utterance opens; from then on
// the real Transcriber owns the audio and the ASR is reused for every later
// utterance. It never returns before its engines are released (the defers below),
// and it is cancelled by the scope's ctx - not by a wall clock (D42#9).
//
// The hand-off walk runs on ITS OWN roster name ("model-downloader", a frozen
// temporary-roster entry) while this loop KEEPS DRAINING into the ring. That is
// not tidiness, it is the difference between "the first sentence is transcribed"
// and "the first sentence is lost to a hash walk": the walk hashes 237 MB of
// installed model files on its way to the announcement, and a consumer that
// blocked on it would let the bounded channel drop every frame of the very
// utterance that triggered it. The ring is what the chain is replayed from once
// the walk returns, so the audio is continuous across the switchover as long as
// the walk finishes inside voiceRingSeconds.
func (rv *residentVoice) run(ctx context.Context, phase1 *speech.Detector) {
	defer phase1.Close()

	var (
		tr   *speech.Transcriber
		asr  *speech.StreamingAsr
		eps  []speech.Endpoint
		outs []speech.Transcript
	)
	defer func() {
		if asr != nil {
			// C11: sessions are released, never unloaded.
			_ = asr.Close()
		}
	}()

	ring := make([]float32, 0, voiceRingSamples)
	pausedClosed := false
	for {
		if rv.paused.Load() {
			if !pausedClosed {
				// The tray paused the listener: close whatever was open and
				// DISCARD it (D43#13's side effect is audio.discard-buffer, so
				// a paused-then-closed utterance must not reach the agent).
				if tr != nil {
					outs = tr.Finalize(outs[:0])
					rv.discard(outs)
				}
				pausedClosed = true
			}
			if !rv.drainWhilePaused(ctx) {
				return
			}
			continue
		}
		pausedClosed = false

		if tr == nil {
			// Phase 1 (and the first-use load, if one is in flight). chainReady
			// is nil until beginHandOff creates it, and a nil case in a select is
			// never taken - so this is one loop, not two.
			select {
			case <-ctx.Done():
				rv.abandonChain()
				return
			case r := <-rv.chainReady:
				var adoptErr error
				tr, asr, adoptErr = rv.adopt(r)
				phase1.Close()
				if adoptErr != nil {
					return
				}
				// The switchover: the ring is the contiguous history that makes
				// the replay exact, and it already holds every frame up to now -
				// including the ones that arrived while the walk and the engine
				// load were running.
				outs = tr.Accept(ring, outs[:0])
				rv.consume(outs)
			case frame, ok := <-rv.frames:
				if !ok {
					return
				}
				pcm, derr := audio.DecodeFrame(frame)
				if derr != nil {
					rv.noteProblem("C8 帧解不开：" + derr.Error())
					continue
				}
				rv.drained.Add(1)
				rv.fed.Add(int64(len(pcm)))
				samples := voiceSamples(pcm)
				ring = voiceRingAppend(ring, samples)
				if rv.handOffStarted {
					continue
				}
				eps = phase1.Accept(samples, eps[:0])
				if len(eps) > 0 {
					rv.beginHandOff(eps[0])
				}
			}
			continue
		}

		select {
		case <-ctx.Done():
			return
		case frame, ok := <-rv.frames:
			if !ok {
				outs = tr.Finalize(outs[:0])
				rv.consume(outs)
				return
			}
			pcm, derr := audio.DecodeFrame(frame)
			if derr != nil {
				rv.noteProblem("C8 帧解不开：" + derr.Error())
				continue
			}
			rv.drained.Add(1)
			rv.fed.Add(int64(len(pcm)))
			outs = tr.Accept(voiceSamples(pcm), outs[:0])
			rv.consume(outs)
		}
	}
}

// beginHandOff is the first-use trigger: it walks D43#4 and starts the hand-off
// AND the engine load on the roster's own model thread. It runs once per process,
// on the asr-infer goroutine, and it deliberately does NOT wait - waiting is what
// would drop the audio (see run).
func (rv *residentVoice) beginHandOff(start speech.Endpoint) {
	rv.handOffStarted = true
	rv.recordTurnOpen(start, rv.fed.Load())

	rv.dispatch(statemachine.EvSummon, nil,
		"麦克风开着时的第一次语音：D43#4（唤醒词是票 41、召唤热键在本进程没有执行器，所以语音就是 summon）")

	rv.handOffs.Add(1)
	rv.mu.Lock()
	rv.askedFor = append(rv.askedFor, voiceAsrModelID)
	rv.mu.Unlock()

	ready := make(chan voiceChainReady, 1)
	rv.chainReady = ready
	hook := rv.plan.HandOff
	if hook == nil {
		// Nothing to wait for: adopt turns a non-zero exit into the refused
		// sentence, and -1 is not a code models.go can produce.
		ready <- voiceChainReady{exit: -1}
		return
	}
	// The walk and the load go on the roster's model thread through the same C11
	// scope, so disposal joins them like everything else this leg started.
	rv.scope.Go(voiceHandOffGoroutine, func(ctx context.Context) {
		exit := hook(voiceAsrModelID, rv.logf)
		if exit != 0 {
			ready <- voiceChainReady{exit: exit}
			return
		}
		asr, err := speech.NewStreamingAsr(speech.AsrConfig{
			EncoderPath: filepath.Join(rv.plan.AsrDir, voiceAsrEncoder),
			DecoderPath: filepath.Join(rv.plan.AsrDir, voiceAsrDecoder),
			TokensPath:  filepath.Join(rv.plan.AsrDir, voiceAsrTokens),
		})
		if err != nil {
			ready <- voiceChainReady{exit: exit, err: err}
			return
		}
		det, err := speech.NewDetector(speech.DetectorConfig{ModelPath: rv.plan.VadModelPath})
		if err != nil {
			_ = asr.Close()
			ready <- voiceChainReady{exit: exit, err: err}
			return
		}
		ready <- voiceChainReady{exit: exit, asr: asr, det: det}
	})
	slog.Info("voice: model hand-off and engine load started on the roster's model thread",
		"goroutine", voiceHandOffGoroutine, "id", voiceAsrModelID, "dir", rv.plan.AsrDir)
}

// abandonChain releases engines the model thread finished opening after the
// consumer was already leaving. Without this a cancellation that lands between
// "engines ready" and "engines adopted" would leak a sherpa session for the rest
// of the process (C11: sessions are released, never leaked).
func (rv *residentVoice) abandonChain() {
	if rv == nil || rv.chainReady == nil {
		return
	}
	select {
	case r := <-rv.chainReady:
		if r.asr != nil {
			_ = r.asr.Close()
		}
		if r.det != nil {
			_ = r.det.Close()
		}
	default:
	}
}

// adopt is the switchover's verdict half: the walk's exit code decides whether
// the engines may exist at all, and the order models.go's AC#5 paragraph asks for
// is preserved by construction - the engines were opened by the same thread that
// ran the walk, after it returned.
//
// A refusal ends the chain for this process (the consumer returns and the defer
// releases what exists): retrying a refused hand-off is a policy this slice does
// not own, and a leg that kept running with no ASR would make the boot report's
// "frames have a consumer" claim a half-truth.
func (rv *residentVoice) adopt(r voiceChainReady) (*speech.Transcriber, *speech.StreamingAsr, error) {
	rv.handOffExit.Store(int64(r.exit))
	if r.exit != 0 {
		rv.noteProblem(fmt.Sprintf("模型交还被拒绝（id=%s，exit=%d）：ASR 的字节不打开，这一句不会被转写；"+
			"models.go 已经把分类原因印在 stderr 上。本进程的语音链到此为止（重试策略不是本片的）",
			voiceAsrModelID, r.exit))
		return nil, nil, errVoiceChainRefused
	}
	if r.err != nil {
		rv.noteProblem("ASR 没起来（交还已通过，但字节打不开）：" + r.err.Error())
		return nil, nil, r.err
	}
	rv.asrLoaded.Store(true)
	// The old endpointer is released by the caller (run's phase1.Close): the
	// transcriber brought its own device state, so re-feeding the history cannot
	// double-classify anything phase 1 already saw.
	slog.Info("voice: the chain is live", "asr", voiceAsrModelID, "vad", rv.plan.VadModelPath,
		"goroutine", voiceGoroutineName, "window_samples", speech.WindowSamples)
	fmt.Printf("wisp: 语音链接管：模型交还已通过（id=%s），ASR 字节从 %s 打开，VAD 句柄换成第二枚会话接管回放\n",
		voiceAsrModelID, rv.plan.AsrDir)
	return speech.NewTranscriber(r.det, r.asr), r.asr, nil
}

// consume fans the seam's transcripts into this leg's readings, the D43 walk and
// the text path.
func (rv *residentVoice) consume(outs []speech.Transcript) {
	for _, t := range outs {
		switch t.Kind {
		case speech.TranscriptPartial:
			rv.partials.Add(1)
			rv.recordPartial(t)
		case speech.TranscriptFinal:
			rv.finals.Add(1)
			rv.finish(t)
		default:
			rv.noteProblem(fmt.Sprintf("未知的转写类型 %v", t.Kind))
		}
	}
}

// discard counts the transcripts a paused listener threw away.
func (rv *residentVoice) discard(outs []speech.Transcript) {
	for range outs {
		rv.discarded.Add(1)
	}
}

// finish closes one turn: the D43 row #11 walk, then the transcript's
// destination.
func (rv *residentVoice) finish(t speech.Transcript) {
	speechMS := t.Endpoint.Speech * 1000 / audio.TargetRate
	handed := false
	if t.Text != "" && rv.tasks != nil {
		if err := rv.tasks.submitTask(t.Text); err != nil {
			rv.noteProblem("转写交给文本路失败：" + err.Error())
		} else {
			handed = true
			rv.handed.Add(1)
		}
	} else if t.Text != "" {
		rv.noteProblem("转写没有接收方：本进程没有任务入口（票 12 的文本路），这一句只留在读数里")
	} else {
		rv.noteProblem("这一句没听清：ASR 对判停句给出空文本（AC#5 的语义；重试与错误类映射不在本片）")
	}
	rv.recordTurnFinal(t, speechMS, handed)
	rv.dispatch(statemachine.EvVadStop, &statemachine.Facts{SpeechMS: speechMS},
		fmt.Sprintf("判停（%s）且语音 %dms 不小于 300ms：D43#11 让球走 Thinking", t.Endpoint.Kind, speechMS))
}

// dispatch walks one D43 edge and projects the result to the ball. An illegal
// transition is recorded, never panicked: the table is frozen and a leg that
// guesses an edge is worse than one that says the table refused it.
func (rv *residentVoice) dispatch(ev statemachine.Event, facts *statemachine.Facts, say string) {
	if rv == nil || rv.machine == nil {
		return
	}
	from := rv.machine.State()
	to, err := rv.machine.Dispatch(ev, facts)
	if err != nil {
		rv.noteProblem(fmt.Sprintf("D43 拒绝 %s + %s：%v", from, ev, err))
		return
	}
	if to == from {
		return
	}
	rv.mu.Lock()
	rv.timeline = append(rv.timeline, voiceTransition{
		From: from, To: to, Event: ev, AtSample: rv.fed.Load(), Say: say,
	})
	rv.mu.Unlock()
	rv.ball.setVoiceState(to)
	slog.Info("voice: D43 walked", "from", string(from), "event", string(ev), "to", string(to), "why", say)
}

// drainWhilePaused drains one frame and drops it. It returns false when the
// stream is over, so a paused leg still exits with its ctx.
func (rv *residentVoice) drainWhilePaused(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	case _, ok := <-rv.frames:
		if !ok {
			return false
		}
		rv.pausedDrops.Add(1)
		return true
	}
}

// togglePause is the tray's "暂停唤醒" executor and the orb's mirror of it. It
// answers the sentence it said and whether it really paused something; false
// means this process owns no listener to pause, which is the honest answer for a
// boot with no capture leg, no models, or [voice] off.
func (rv *residentVoice) togglePause() (string, bool) {
	if rv == nil || rv.scope == nil {
		if rv != nil && rv.verdict != "" {
			return rv.verdict + "；这一按什么都没停", false
		}
		return "本进程没有语音链可以暂停（没有采集帧，或 [voice] 关着）：这一按什么都没停", false
	}
	now := !rv.paused.Load()
	rv.paused.Store(now)
	if now {
		// D43#13 is the legal edge for "listening stopped by the user" and its
		// freeze side effect is audio.discard-buffer, which is exactly what the
		// consumer does with the open utterance. Only Listening has that row.
		if rv.machine.State() == statemachine.StateListening {
			rv.dispatch(statemachine.EvVeto, nil, "托盘暂停唤醒：D43#13 丢弃缓冲")
		}
		slog.Info("voice: listening paused by the user", "drained", rv.drained.Load())
		return "已暂停唤醒监听：C8 帧照旧被排空但不再进 VAD/ASR（丢弃计数照记），托盘勾选跟着这一位走", true
	}
	if rv.machine.State() == statemachine.StateWarm {
		rv.dispatch(statemachine.EvSummon, nil, "托盘恢复唤醒：D43#29 Warm 到 Listening")
	}
	slog.Info("voice: listening resumed by the user", "paused_drops", rv.pausedDrops.Load())
	return "已恢复唤醒监听：下一句语音重新进 VAD/ASR", true
}

// pausedForTray is the tray checkmark's READ half: whether this leg really owns a
// listener (scope != nil) and, if so, whether the user paused it. The second
// return is "is there a truth source at all", which is what keeps the projection
// from writing a false into a checkmark it cannot support.
func (rv *residentVoice) pausedForTray() (bool, bool) {
	if rv == nil || rv.scope == nil {
		return false, false
	}
	return rv.paused.Load(), true
}

// noteProblem records and says one thing this leg could not do. Every branch that
// gives up goes through here so the boot report and the log cannot disagree.
func (rv *residentVoice) noteProblem(why string) {
	if rv == nil {
		return
	}
	rv.mu.Lock()
	rv.problems = append(rv.problems, why)
	rv.mu.Unlock()
	slog.Warn("voice: " + why)
	fmt.Printf("wisp: 语音链：%s\n", why)
}

// logf is what the hand-off walk calls with its progress ticks (models.go's
// logf parameter). It lands in the process log and never on stdout: a
// per-second model tick must not be mistaken for a user-facing line.
func (rv *residentVoice) logf(format string, args ...any) {
	slog.Info("voice: model hand-off: " + fmt.Sprintf(format, args...))
}

// recordTurnOpen starts the turn this endpoint opened, on the audio clock.
func (rv *residentVoice) recordTurnOpen(start speech.Endpoint, at int64) {
	if start.Kind != speech.EndpointStart {
		return
	}
	rv.mu.Lock()
	defer rv.mu.Unlock()
	rv.turns = append(rv.turns, voiceTurn{
		OpenSample:  at,
		StartSample: at - int64(start.Total),
	})
}

// recordPartial stamps the first text of the open turn.
func (rv *residentVoice) recordPartial(t speech.Transcript) {
	rv.mu.Lock()
	defer rv.mu.Unlock()
	if len(rv.turns) == 0 {
		return
	}
	turn := &rv.turns[len(rv.turns)-1]
	if turn.FirstTextAt != 0 {
		return
	}
	turn.FirstTextAt = rv.fed.Load()
	turn.FirstChunk = t.Latency
}

// recordTurnFinal closes the open turn.
func (rv *residentVoice) recordTurnFinal(t speech.Transcript, speechMS int, handed bool) {
	rv.mu.Lock()
	defer rv.mu.Unlock()
	if len(rv.turns) == 0 {
		return
	}
	turn := &rv.turns[len(rv.turns)-1]
	turn.FinalAt = rv.fed.Load()
	turn.FinalLatency = t.Latency
	turn.SpeechMS = speechMS
	turn.ClosedBy = t.Endpoint.Kind.String()
	turn.Text = t.Text
	turn.Handed = handed
}

// timelineSnapshot and turnsSnapshot are the reading surfaces (copies, so a
// caller cannot mutate the leg's record).
func (rv *residentVoice) timelineSnapshot() []voiceTransition {
	if rv == nil {
		return nil
	}
	rv.mu.Lock()
	defer rv.mu.Unlock()
	return append([]voiceTransition(nil), rv.timeline...)
}

func (rv *residentVoice) turnsSnapshot() []voiceTurn {
	if rv == nil {
		return nil
	}
	rv.mu.Lock()
	defer rv.mu.Unlock()
	return append([]voiceTurn(nil), rv.turns...)
}

// problemsSnapshot is what this leg could not do, in order.
func (rv *residentVoice) problemsSnapshot() []string {
	if rv == nil {
		return nil
	}
	rv.mu.Lock()
	defer rv.mu.Unlock()
	return append([]string(nil), rv.problems...)
}

// handOffIDsSnapshot is every id this leg asked the model store to hand over. It
// is the "was the hook called, and with what" reading, kept separable from the
// counters so a test can name the argument rather than only its count.
func (rv *residentVoice) handOffIDsSnapshot() []string {
	if rv == nil {
		return nil
	}
	rv.mu.Lock()
	defer rv.mu.Unlock()
	return append([]string(nil), rv.askedFor...)
}

// drainsFrames reports whether this handle really owns the seam's consumer. It is
// what the capture leg's posture line reads, so "the frames have a consumer" is
// answered by the wiring instead of by a ticket boundary.
func (rv *residentVoice) drainsFrames() bool {
	return rv != nil && rv.scope != nil
}

// absentReason is why there is no consumer, in one sentence built from the same
// verdict the boot report prints.
func (rv *residentVoice) absentReason() string {
	if rv == nil {
		return "本进程没有语音链"
	}
	if rv.verdict != "" {
		return rv.verdict
	}
	return "本进程没有语音链"
}

// posture is the boot report's 语音链 field, built from what this leg really
// resolved rather than from a hope.
func (rv *residentVoice) posture() string {
	if rv == nil {
		return "语音链未装配：本进程没有任何语音路径"
	}
	if rv.scope == nil {
		if rv.verdict != "" {
			return rv.verdict
		}
		return "语音链未装配：本进程没有任何语音路径"
	}
	return fmt.Sprintf("C8 帧通道已接上 1 个消费者（D38b 名册名 %s，C11 作用域 %q）；"+
		"VAD=%s；ASR 模型 id=%s，字节只在第一次开句时由模型交还通过后从 %s 打开；"+
		"目前排空 %d 帧／部分结果 %d 条／判停 %d 句／交给文本路 %d 句／模型交还 %d 次（exit=%d）；"+
		"配置来源 %s",
		voiceGoroutineName, "voice-15e", rv.plan.VadModelPath, voiceAsrModelID, rv.plan.AsrDir,
		rv.drained.Load(), rv.partials.Load(), rv.finals.Load(), rv.handed.Load(),
		rv.handOffs.Load(), rv.handOffExit.Load(), rv.plan.Source)
}

// stopLine is the teardown record: the counters the exit line carries, taken off
// the same atomics the posture read.
func (rv *residentVoice) stopLine() string {
	if rv == nil || rv.scope == nil {
		return "voice leg: not assembled in this process"
	}
	return fmt.Sprintf("voice leg stopped: frames_drained=%d paused_drops=%d partials=%d finals=%d discarded=%d handed=%d handoffs=%d handoff_exit=%d asr_loaded=%v",
		rv.drained.Load(), rv.pausedDrops.Load(), rv.partials.Load(), rv.finals.Load(),
		rv.discarded.Load(), rv.handed.Load(), rv.handOffs.Load(), rv.handOffExit.Load(),
		rv.asrLoaded.Load())
}

// startResidentVoice is the assembly root's call: resolve the plan off this
// process's own data root, build the handle, book the thread, and hand both ends
// of the wake-pause hop to the ball host (ticket 15e's third seam).
//
// It never stops the boot. Every branch that cannot build a voice chain returns a
// handle with its verdict already written, and the boot report prints that
// verdict - the same posture startResidentAudio takes for the microphone.
func startResidentVoice(rt *proc.Runtime, ra *residentAudio, rb *residentBall, tasks voiceTranscriptSink) *residentVoice {
	if rt == nil {
		return nil
	}
	rv := buildResidentVoice(rt.Layout.DataDir, ra, rb, tasks)
	rv.start(rt.Registry)
	if rb != nil {
		// Both ends of the tray's pause-wake hop, attached here because this leg
		// OWNS the listener: toggling is meaningless without a voice chain, and
		// the read half answers "no truth source" for a boot that built none.
		rb.attachWakePause(rv.togglePause, rv.pausedForTray)
		rb.mirrorTrayMute()
	}
	return rv
}

// buildResidentVoice is the same assembly with the data root named, so the
// evidence table can drive the production resolution against a config it wrote.
func buildResidentVoice(dataDir string, ra *residentAudio, rb *residentBall, tasks voiceTranscriptSink) *residentVoice {
	rv := &residentVoice{
		ball: rb, tasks: tasks,
		machine: statemachine.New(statemachine.Options{Initial: statemachine.StateSleeping}),
	}

	cfgPath := filepath.Join(dataDir, configFileName)
	c, _, err := config.LoadFile(cfgPath, nil)
	source := cfgPath
	if err != nil || c == nil {
		c = config.NewDefaults()
		source = "compiled defaults (config.NewDefaults)"
	}

	switch {
	case !c.Voice.Enabled:
		rv.verdict = "语音链未构造：[voice] enabled=false（配置来源 " + source + "）：C8 帧通道没有消费者，帧按 D38d 计数丢弃"
		return rv
	case ra == nil || ra.frames == nil:
		rv.verdict = "语音链未构造：本进程没有采集腿，因此没有 C8 帧通道可排空（配置来源 " + source + "）"
		return rv
	case c.Voice.ASR.Provider != "" && c.Voice.ASR.Provider != "local-sherpa":
		rv.verdict = fmt.Sprintf("语音链未构造：[voice].asr.provider=%q，本片的链路只有本地 sherpa 一条（云 ASR 是 RESERVED，票 61）",
			c.Voice.ASR.Provider)
		return rv
	}

	plan, why := resolveVoicePlan(dataDir, c)
	plan.Source = source
	if why != "" {
		rv.verdict = why
		return rv
	}
	rv.plan = plan
	rv.frames = ra.frames
	rv.verdict = "语音链已装配"
	return rv
}

// resolveVoicePlan finds the two model locations. It answers a reason instead of
// guessing: a missing VAD model is a boot sentence, not a silent fallback to a
// developer's tree (this binary ships without third_party).
//
// The VAD is deliberately NOT resolved through a hand-off: it is the part of
// listening that must be warm before any capability is used, and calling the
// hand-off at boot is exactly what would make silence fetch a model.
func resolveVoicePlan(dataDir string, c *config.Config) (voicePlan, string) {
	store := modelStoreDir(dataDir, c)
	vad := filepath.Join(store, voiceVadModelID, voiceVadFile)
	if !voiceFileExists(vad) {
		if env := os.Getenv("WISP_SILERO_MODEL"); env != "" {
			slog.Info("voice: the store carries no VAD model; using WISP_SILERO_MODEL", "path", env, "store", store)
			vad = env
		} else {
			return voicePlan{}, fmt.Sprintf("语音链未构造：VAD 模型没落地（%s 不存在，也没有 WISP_SILERO_MODEL）；"+
				"wisp models ensure %s 是它的交还入口", vad, voiceVadModelID)
		}
	}
	return voicePlan{
		VadModelPath: vad,
		AsrDir:       filepath.Join(store, voiceAsrModelID),
		HandOff:      residentVoiceHandOff(dataDir),
	}, ""
}

// voiceFileExists is the one existence test this file uses. It answers "is there
// a regular file here" and nothing else: the store's own verification (hashes,
// signature) is models.go's, not a second policy grown here.
func voiceFileExists(p string) bool {
	if p == "" {
		return false
	}
	st, err := os.Stat(p)
	return err == nil && st.Mode().IsRegular()
}

// residentVoiceHandOff is models.go:25's reserved call site, and the only
// production caller of modelStore.handOffModel outside `wisp models ensure`.
//
// It re-reads config.toml per call for the same reason every other hot-tier view
// in this process does: the task pipeline's Manager does not exist at the point
// the voice chain is assembled, and the models section is what names the store
// root and the mirrors.
func residentVoiceHandOff(dataDir string) voiceHandOff {
	return func(id string, logf func(string, ...any)) int {
		cfg, err := loadModelConfig(dataDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "wisp: 语音链：模型交还没做（配置未就绪，Unconfigured）：%v\n", err)
			return 2
		}
		store, err := openModelStore(dataDir, cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "wisp: 语音链：模型库打不开，交还没做：%v\n", err)
			return 2
		}
		slog.Info("voice: model hand-off starting", "id", id, "store", store.root, "manifest", store.src)
		return store.handOffModel(modelsIO{dataDir: dataDir}, id, logf)
	}
}
