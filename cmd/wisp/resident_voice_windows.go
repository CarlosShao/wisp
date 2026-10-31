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
// directly. Row #4 is Sleeping + EvSummon -> Listening and, for a boot with no
// wake-word engine, it is the ONLY legal way into Listening, so "the user spoke
// while the microphone was open" is dispatched as the summon - a wiring
// decision, not a table change. Ticket 41 landed KWS and the reserved sentence
// came true: an ARMED boot wakes through row #7 instead (beginKwsWake dispatches
// EvWakeWord, the table byte-identical), Sleeping keeps the summon for the
// re-arm ordering, and the header comment 41-a1 read is now the shipped shape.
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
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/plugin"
	"github.com/CarlosShao/wisp/internal/proc"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/session"
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

	// voicePunctModelID / voicePunctFile are the punctuation entry's manifest id
	// (purpose "punctuation", status ok) and the artifact name inside its archive.
	// 27-r1 measured the file this store path resolves to: 75,519,198 B, sha256
	// 65a3fb9f5ad7bfb9..., which is the manifest's own pin.
	voicePunctModelID = "punc-ct-transformer-zh-en-vocab272727"
	voicePunctFile    = "model.int8.onnx"

	// voiceTtsModelID / voiceTtsAcoustic ... are the matcha checkpoint's manifest
	// id and the artifact names inside it. The vocoder is the manifest's second
	// `files` entry (no archive), so it lands in the same directory. Status:
	// blocked-p3 - this is a dev/measurement TTS, not the product's default (R5,
	// ticket 61); nothing here selects an engine.
	voiceTtsModelID  = "tts-matcha-zh-baker"
	voiceTtsAcoustic = "model-steps-3.onnx"
	voiceTtsVocoder  = "vocos-22khz-univ.onnx"
	voiceTtsLexicon  = "lexicon.txt"
	voiceTtsTokens   = "tokens.txt"
	voiceTtsDict     = "dict"

	// Ticket 41's wake-word entry: the manifest id (purpose "kws", status ok,
	// license Apache-2.0, quant int8) and the artifact names its archive pins.
	// The store spelling is epoch-12-avg-2; the dev copies under third_party
	// carry an epoch-99 spelling whose int8 triple hashes bit-identical to these
	// (measured before this leg landed), so both spellings build the same engine
	// the 41-a2 engine leg measured. The store paths are filled into
	// speech.KwsConfig's five fields directly; the engine's own package constants
	// are never consulted for the store layout.
	voiceKwsModelID  = "kws-zipformer-wenetspeech-3.3M-2024-01-01"
	voiceKwsEncoder  = "encoder-epoch-12-avg-2-chunk-16-left-64.int8.onnx"
	voiceKwsDecoder  = "decoder-epoch-12-avg-2-chunk-16-left-64.int8.onnx"
	voiceKwsJoiner   = "joiner-epoch-12-avg-2-chunk-16-left-64.int8.onnx"
	voiceKwsTokens   = "tokens.txt"
	voiceKwsKeywords = "keywords.txt"
)

// voiceTtsRuleFsts are the three text-normalisation graphs the checkpoint ships.
// Missing ones are dropped rather than fatal: a date/number rule FST changes how
// "2026" is read, not whether the engine speaks at all.
var voiceTtsRuleFsts = []string{"date.fst", "number.fst", "phone.fst"}

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
	// asr is the C9 seam, not the concrete session: since wire-r1 the engine comes
	// from the session scope, which hands back an AsrEngine. NewTranscriber takes
	// the interface, so nothing downstream needs the concrete type.
	asr speech.AsrEngine
	det *speech.Detector
	err error
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
	// PunctModelPath is the CT-Transformer onnx file the dictation path attaches
	// through Transcriber.SetPunctuator (ticket 27). Empty means the store carries
	// no punctuation model: the chain still runs and its finals are raw, which is
	// the "stage switched off" state 27-r1 pinned.
	PunctModelPath string
	// Tts is the matcha checkpoint's file set (ticket 26). Incomplete means this
	// process cannot speak: the announcement path says so and the D10 notification
	// remains the delivery path.
	Tts voiceTtsFiles
	// Kws is the wake-word checkpoint's file set (ticket 41), resolved off the
	// same store layout when [voice] wake_word.enabled asks for it. Incomplete
	// means this process arms no wake-word engine and stays in Sleeping: the
	// chain runs exactly as ticket 41 found it, and KwsNote says which of the
	// two reasons (opt-in off, model missing) produced the empty set.
	Kws     voiceKwsFiles
	KwsNote string
	// WarmTimeout is [session].warm_timeout_sec: how long the engines stay hot
	// before the Settling window starts. 0 means the session package's own default
	// (90 s), which is also the config default.
	WarmTimeout time.Duration
}

// voiceKwsFiles is the KWS transducer triple plus its token and keyword tables,
// spelled as the five paths speech.KwsConfig wants. The manifest pins each
// file's sha256; this leg checks existence only, because the store's own
// verification (hashes, signature) is models.go's - the same split the VAD, the
// punctuation stage and the TTS checkpoint already follow.
type voiceKwsFiles struct {
	Encoder  string
	Decoder  string
	Joiner   string
	Tokens   string
	Keywords string
}

// complete reports whether every field KwsConfig requires is present. A partial
// set is not a KWS this leg may load: NewKwsEngine refuses an empty path, and a
// leg that tried anyway would spend an engine construction on a refusal.
func (f voiceKwsFiles) complete() bool {
	return f.Encoder != "" && f.Decoder != "" && f.Joiner != "" && f.Tokens != "" && f.Keywords != ""
}

// kwsConfig fills speech.KwsConfig's five fields directly off the resolved
// paths. The threshold stays at the package constant (DefaultKwsThreshold
// 0.1): the API defaults (1.0/0.25) score zero hits on the engine leg's own
// fixtures, and 41-a2 measured 0.1/3.0 against the checkpoint's test_wavs.
func (f voiceKwsFiles) kwsConfig() speech.KwsConfig {
	return speech.KwsConfig{
		EncoderPath:  f.Encoder,
		DecoderPath:  f.Decoder,
		JoinerPath:   f.Joiner,
		TokensPath:   f.Tokens,
		KeywordsPath: f.Keywords,
	}
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

	// The wire-r1 half: the owner queue (work for THIS goroutine), the engine
	// residence (session warm epochs + the single-family slot), and the C25 scope
	// the announce gate reads.
	queue     *voiceOwnerQueue
	residency *voiceResidency
	prov      *risk.Provenance
	provScope string
	// newSpeaker is the render seam. Production is the waveOut device; the tests
	// hand in a recorder so the samples an announcement produced can be read back
	// without a sound card, and one live case uses the device itself.
	newSpeaker func(rate int) (voiceSpeaker, error)

	// TICKET 26's one-needle seam (26r3): the capture half-duplex gate this
	// process already mounts (residentAudio's *audio.HalfDuplexGate, PathT),
	// handed in by the assembly root. On PathT SetSpeaking(true) closes capture -
	// D16's "关闭麦克风" - and SetSpeaking(false) reopens it with zero model
	// reload (internal/audio/gate.go). Before 26r3 the gate's own comment said it
	// is "driven by the TTS side" and the TTS side never drove it: production
	// callers were zero. The field is an anonymous one-method interface so the
	// wiring stays on the cmd/wisp side of the seam (this leg may not edit
	// internal/audio) and the
	// test can drive the production speak() against a recording gate; nil means
	// this boot built no capture leg (the same nil-safety every other handle on
	// this struct carries).
	captureGate interface {
		SetSpeaking(speaking bool)
	}
	// speakGateOn / speakGateOff are the needle's two observability counters:
	// how many announcements closed the gate and how many times the chain
	// released it (one speak() = at most one of each, whatever its outcome). The
	// boot report's announce reading and the stop line both quote them, so
	// "the gate was driven" is a reading and not an inference from a code path.
	speakGateOn  atomic.Uint64
	speakGateOff atomic.Uint64

	// warmFeed is D43 #26's mic.keep-off, and since ticket 41 a REAL feed gate:
	// the frames keep draining, but with a wake-word engine resident nothing they
	// produce reaches the text path until a wake resumes the session - the KWS
	// is the ears instead (R28-2's landing: Warm 期间 ASR 不再吃帧). Without a
	// resident KWS the gate keeps the wire-r1 latch shape: the transcript path
	// stays the only wake source and eats the frames exactly as it did before.
	// Set by the session scope's MicClose hook, cleared on Resume.
	warmFeed atomic.Bool

	// kwsWake is ticket 41's wake-word engine, loaded once in start() when the
	// opt-in asks for it and the model resolved (the same boot-thread precedent
	// the VAD session follows: building it inside the consumer loop would drop
	// the first frames of whatever the microphone was already saying). The
	// consumer goroutine is its only Accept caller (D38b) and closes it on its
	// way out; arming and pausing are the state machine's business, expressed
	// here as the kwsFeed gate below - never as engine methods.
	kwsWake speech.WakeWordEngine

	announceMu      sync.Mutex
	announcePending []voiceAnnounceRequest

	// handOffStarted and chainReady are the first-use trigger's own state, owned
	// by the consumer goroutine: the walk AND the engine load run on the roster's
	// model thread while the loop keeps draining, and the buffered channel is how
	// the ready engines (or the reason there are none) come back.
	handOffStarted bool
	chainReady     chan voiceChainReady

	// TICKET 381 R26-1: interrupt channel for播报中打断。Produced by three sources:
	// hotkey/cancel key, ball click, KWS veto words. Consumed by speak() to stop TTS
	// within ≤400ms and dispatch EvInterrupt. Buffered =1 so a request during the
	// final Play() frame still fires without blocking the caller.
	interruptCh chan struct{}

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

	// wire-r1 counters: one per hop of the three paths, so every sentence in the
	// evidence table is a reading and not an inference.
	announceAsked   atomic.Uint64
	announceEmpty   atomic.Uint64
	announceTooLong atomic.Uint64
	announceDropped atomic.Uint64
	ttsGated        atomic.Uint64
	ttsGatedOk      atomic.Uint64
	ttsNoEngine     atomic.Uint64
	ttsNoDevice     atomic.Uint64
	ttsSilent       atomic.Uint64
	ttsFailed       atomic.Uint64
	ttsSpoken       atomic.Uint64
	ttsChunks       atomic.Uint64
	ttsSamples      atomic.Uint64
	warmTurnEnd     atomic.Uint64
	warmWake        atomic.Uint64
	warmDiscards    atomic.Uint64
	punctAttached   atomic.Bool
	punctFinals     atomic.Uint64
	// asrGone is set by the announce path when a family swap released the ASR the
	// running transcriber still holds; the loop re-arms phase 1 on it.
	asrGone atomic.Bool

	// Ticket 41's feed counters. kwsFedSamples counts every sample handed to the
	// wake-word engine, asrFedSamples every sample handed to the transcriber's
	// Accept; the Warm reading is the PAIR of them (the KWS count advances while
	// the ASR count stands still), which is what makes "the gate really gates" a
	// reading instead of an inference. kwsHits/kwsVetos count the engine's two
	// hit kinds separately, because they walk different D43 rows.
	kwsFedSamples atomic.Int64
	asrFedSamples atomic.Int64
	kwsHits       atomic.Uint64
	kwsVetos      atomic.Uint64

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
	queue := newVoiceOwnerQueue()
	return &residentVoice{
		frames:      frames,
		plan:        plan,
		ball:        ball,
		tasks:       tasks,
		machine:     statemachine.New(statemachine.Options{Initial: statemachine.StateSleeping}),
		queue:       queue,
		residency:   newVoiceResidency(plan, queue, nil),
		provScope:   voiceAnnounceScope,
		newSpeaker:  openWaveOutSpeaker,
		interruptCh: make(chan struct{}, 1),
	}
}

// TICKET 381 R26-1: interrupt trigger API. Call this from any of the three
// interrupt sources (hotkey/cancel key, ball click, KWS veto word) to fire an
// EvInterrupt during播报。Non-blocking send into a buffered chan =1 guarantees
// no caller blocks even if speak() is mid-Pause(). The speak() consumer checks
// this channel in its hot loop and will dispatch EvInterrupt within ≤400ms.
func (rv *residentVoice) TriggerInterrupt() {
	if rv == nil || rv.interruptCh == nil {
		return
	}
	select {
	case rv.interruptCh <- struct{}{}:
	default:
		// Channel full = interrupt already queued for this播报句，drop silently.
	}
}

// attachProvenance hands the leg the process's C25 engine and opens the channel
// scope the announce gate reads. A nil engine means this boot has no provenance
// engine at all; the leg then owns one with the same options panel_assets.go uses
// for its own process-lifetime scope, so the gate is never "skipped because there
// was no engine" (CheckText on a nil engine would be a nil dereference, and a
// skipped check is the fail-open direction).
func (rv *residentVoice) attachProvenance(p *risk.Provenance) {
	if rv == nil {
		return
	}
	owned := false
	if p == nil {
		p = risk.NewProvenance(risk.ProvOptions{NoProbe: true})
		owned = true
	}
	scope := p.OpenScope(rv.provScope)
	if scope == nil {
		if !owned {
			// The id is already registered by somebody else: take a private engine
			// rather than check against a scope this leg does not own.
			p = risk.NewProvenance(risk.ProvOptions{NoProbe: true})
			scope = p.OpenScope(rv.provScope)
		}
	}
	if scope == nil {
		rv.noteProblem("播报闸未接线：C25 作用域打不开，播报路径会拒绝出声（fail-closed）")
		return
	}
	rv.prov = p
	slog.Info("voice: ChTTS gate wired", "scope", rv.provScope, "channel", string(risk.ChTTS))
}

// attachCaptureGate hands the announce path the capture leg's half-duplex gate
// (ticket 26's one needle, 26r3). It is the only writer of rv.captureGate and
// the assembly root calls it once, right after the voice leg exists - the same
// ordering attachVoice has, because both are "the leg is told what already
// exists" hops: the gate was mounted by assembleCapture before the voice chain
// was built, and this hand-off is what makes the TTS side finally drive it.
// nil-safe like every other attach here; a boot with no capture leg leaves the
// announce path unmuted-but-honest (speakGateOn/off stay zero and the reading
// says so).
func (rv *residentVoice) attachCaptureGate(gate interface {
	SetSpeaking(speaking bool)
},
) {
	if rv == nil {
		return
	}
	rv.captureGate = gate
	if gate != nil {
		slog.Info("voice: capture gate wired to the announce path",
			"seam", "audio.HalfDuplexGate.SetSpeaking", "path", string(audio.PathT))
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
	// The mic half of "Warm closes the microphone FIRST". What this chain can
	// honestly close is the FEED, not the device: the C8 frames that carry the wake
	// in are the capture leg's, and a chain that switched the device off in Warm
	// would never hear the next utterance. Since ticket 41 the Warm ears are the
	// wake-word engine (the frames keep feeding IT while the text path's feed
	// stays closed), which is what makes the feed gate below a real gate instead
	// of the latch wire-r1 could honestly build.
	rv.residency.onMicClose = func() { rv.warmFeed.Store(true) }
	// Ticket 41's arm: the opt-in asked for the wake word and the model resolved,
	// so the engine is built HERE - on the boot thread, before the consumer is
	// booked, for the same measured reason the VAD session is (a construction
	// inside the loop would drop the first frames of whatever the microphone was
	// already saying). A successful build walks D43#5: Sleeping + EvKwsEnabled ->
	// Armed, with the engine as the kws.load side effect this wiring executes. A
	// refusal is a boot sentence and a Sleeping process, never a half-armed one.
	if rv.plan.Kws.complete() {
		k, kerr := speech.NewKwsEngine(rv.plan.Kws.kwsConfig())
		if kerr != nil {
			rv.noteProblem("KWS 没装载（这一进程不进 Armed，唤醒源仍只有 summon 一条）：" + kerr.Error())
		} else {
			rv.kwsWake = k
			rv.dispatch(statemachine.EvKwsEnabled, nil,
				"唤醒词引擎已装载：D43#5 Sleeping 到 Armed（KWS 是这一进程的常驻唤醒源）")
			slog.Info("voice: wake-word engine loaded and armed",
				"model", voiceKwsModelID, "keywords", rv.plan.Kws.Keywords,
				"goroutine", voiceGoroutineName, "threshold", speech.DefaultKwsThreshold)
		}
	}
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
	// The queue is closed AFTER the consumer is joined: a post that lands between
	// the two is refused (counted) instead of growing work nobody will run.
	if rv.queue != nil {
		rv.queue.close()
	}
	// D38(e) step 4's engine half, now the session's cascade instead of a defer:
	// whatever Warm was still holding is released here, and only here.
	if rv.residency != nil {
		if err := rv.residency.dispose("process-exit"); err != nil {
			slog.Warn("voice: 引擎驻留没有干净收口", "err", err.Error())
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
	// phase1 is a variable, not the parameter, because the announce path can swap
	// the engine family and re-arm a fresh endpointer (see the loop head). The defer
	// reads the CURRENT one, so a rebuilt detector is released too.
	defer func() {
		if phase1 != nil {
			_ = phase1.Close()
		}
	}()

	var (
		tr   *speech.Transcriber
		asr  speech.AsrEngine
		eps  []speech.Endpoint
		outs []speech.Transcript
	)
	defer func() {
		if asr != nil {
			// C11: sessions are released, never unloaded.
			_ = asr.Close()
		}
	}()
	// The wake-word engine is the consumer goroutine's too: it is the only Accept
	// caller (D38b), so it is the one that releases it (C11).
	defer func() {
		if rv.kwsWake != nil {
			_ = rv.kwsWake.Close()
			rv.kwsWake = nil
		}
	}()

	ring := make([]float32, 0, voiceRingSamples)
	pausedClosed := false
	for {
		// The family swap. An announcement loads the TTS through the slot, and the
		// slot releases the ASR BEFORE it does (D32 16.3.3: one family resident, so
		// the peak is max and not the sum). The transcriber this loop was driving
		// therefore holds a released session, and the honest move is to drop it and
		// re-arm phase 1: the next utterance pays the ASR's cold load and walks the
		// hand-off again, which is D32's on-demand price and is counted, not hidden.
		if rv.asrGone.Load() && tr != nil {
			tr, asr = nil, nil
			rv.asrGone.Store(false)
			rv.handOffStarted = false
			rv.asrLoaded.Store(false)
			if det, derr := speech.NewDetector(speech.DetectorConfig{ModelPath: rv.plan.VadModelPath}); derr != nil {
				rv.noteProblem("播报换族后重建 VAD 失败：" + derr.Error())
				return
			} else {
				phase1 = det
			}
			slog.Info("voice: the announce swapped the engine family; phase 1 re-armed",
				"reason", "tts-resident", "asr_released", true)
		}
		// The session's own cascade (90 s warm idle + 3 s settling) released the
		// engines UNDER a chain this loop still holds: the transcriber's session
		// is closed, and a live chain that keeps feeding it would silence the
		// wake forever (a released session's Feed is the contract's no-op). The
		// honest move is the announce branch's: drop the transcriber, forget the
		// in-flight chain, re-arm phase 1 - the next wake starts a fresh chain
		// and pays the load it really owes. The guard is the LIVE transcriber,
		// nothing else: a chain still walking (tr == nil, the first epoch not
		// open yet) must NOT fire this - stage() answers Disposed for a residency
		// that never opened a scope, and resetting handOffStarted mid-walk is
		// what turns one wake into four loads. A chainless boot self-heals
		// instead: every wake path's acquireTurn opens a fresh epoch when the
		// scope is gone. This runs on the owner goroutine between queue jobs, so
		// it can never observe a half-finished family change (those complete
		// inside one queued job).
		if tr != nil && rv.residency != nil && rv.residency.stage() == session.StageDisposed {
			tr, asr = nil, nil
			rv.handOffStarted = false
			rv.asrLoaded.Store(false)
			if det, derr := speech.NewDetector(speech.DetectorConfig{ModelPath: rv.plan.VadModelPath}); derr != nil {
				rv.noteProblem("会话级联释放后重建 VAD 失败：" + derr.Error())
				return
			} else {
				phase1 = det
			}
			slog.Info("voice: the session cascade disposed the engines; phase 1 re-armed",
				"reason", "settling-expiry", "asr_released", true)
		}
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
			// never taken - so this is one loop, not two. The owner queue is the
			// third case: session timers and announcements land HERE (see
			// resident_voice_engines_windows.go's threading note), which is what
			// keeps the dispose cascade off the runtime timer goroutine.
			select {
			case <-ctx.Done():
				rv.abandonChain()
				return
			case <-rv.queue.wakeChan():
				rv.queue.run()
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
				// Ticket 41's insertion point: the frame pipeline, after the ring.
				// A hit inside kwsFeed dispatches its own D43 row (see onKwsHit);
				// while the machine is Armed the wake-word engine is the ONLY
				// legal ears - the phase-1 endpoint has no row out of Armed (row
				// #4 needs a Sleeping machine), so the VAD is not fed here and a
				// plain utterance without the wake word stays unheard. That is
				// the privacy property Armed exists for, not a defect.
				if rv.kwsFeed(samples) && rv.machine.State() == statemachine.StateArmed {
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
		case <-rv.queue.wakeChan():
			rv.queue.run()
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
			samples := voiceSamples(pcm)
			// The ring is kept for the whole life of the chain since wire-r1, not
			// only until the first speech. With the ticket-41 gate the Warm wake
			// no longer replays it: the frames that closed the feed went to the
			// wake-word engine, and the audio a KWS wake woke on is the WAKE WORD
			// itself, which is not a command - replaying it would hand 「你好军哥」
			// to the agent as a task. The ring stays the first-use switchover's
			// replay history, and the wake-word-utterance handling is named in
			// 41a2r1's evidence table rather than papered over.
			ring = voiceRingAppend(ring, samples)
			if rv.warmFeed.Load() && rv.kwsFeed(samples) {
				// The feed gate, ticket 41 (R28-2's landing): Warm closed the
				// text path and the wake-word engine ate the frame as the wake
				// source. The transcriber sees nothing until a wake reopens the
				// gate; without a resident KWS the gate keeps the wire-r1 latch
				// shape and the transcript path stays the wake source.
				continue
			}
			outs = tr.Accept(samples, outs[:0])
			rv.asrFedSamples.Add(int64(len(samples)))
			rv.consume(outs)
		}
	}
}

// beginHandOff is the first-use trigger: it walks D43#4 and starts the hand-off
// AND the engine load on the roster's own model thread. It runs once per process,
// on the asr-infer goroutine, and it deliberately does NOT wait - waiting is what
// would drop the audio (see run).
//
// Since ticket 41 this is the VAD endpoint's entry only; the wake-word entry is
// beginKwsWake below. Both share startChainLoad - the walk and the load are the
// same machinery, only the D43 row that justifies them differs.
func (rv *residentVoice) beginHandOff(start speech.Endpoint) {
	rv.handOffStarted = true
	rv.recordTurnOpen(start, rv.fed.Load())

	rv.dispatch(statemachine.EvSummon, nil,
		"麦克风开着时的第一次语音：D43#4（本进程没武装唤醒词时语音就是 summon；KWS 武装后 Sleeping 只有重武装一条路，见 onKwsHit）")

	rv.startChainLoad()
}

// startChainLoad is the walk-and-load half the two wake entries share: the model
// hand-off and the ASR engine load run on the roster's model thread while the
// consumer keeps draining into the ring, and the buffered channel is how the
// ready engines (or the reason there are none) come back.
func (rv *residentVoice) startChainLoad() {
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
	// scope, so disposal joins them like everything else this leg started. Since
	// wire-r1 the load itself goes through the session scope's ASR epoch
	// (resident_voice_engines_windows.go): that is what makes the engine the
	// session keeps hot, and it is why the walk still runs BEFORE the bytes are
	// opened - acquireTurn is the call that opens them.
	rv.scope.Go(voiceHandOffGoroutine, func(ctx context.Context) {
		exit := hook(voiceAsrModelID, rv.logf)
		if exit != 0 {
			ready <- voiceChainReady{exit: exit}
			return
		}
		handle, _, err := rv.residency.acquireTurn(session.EngineASR, "first-use")
		if err != nil {
			ready <- voiceChainReady{exit: exit, err: err}
			return
		}
		det, err := speech.NewDetector(speech.DetectorConfig{ModelPath: rv.plan.VadModelPath})
		if err != nil {
			ready <- voiceChainReady{exit: exit, err: err}
			return
		}
		ready <- voiceChainReady{exit: exit, asr: handle.Asr(), det: det}
	})
	slog.Info("voice: model hand-off and engine load started on the roster's model thread",
		"goroutine", voiceHandOffGoroutine, "id", voiceAsrModelID, "dir", rv.plan.AsrDir)
}

// kwsFeed hands one frame's samples to the wake-word engine and answers whether
// it ate them. It is the ARming gate: the engine runs when the machine is Armed
// (D43#5 made it the resident wake source) or when the text feed is closed in the
// Warm posture (ticket 41's ruling: the Warm gate stays open and the KWS is the
// ears). Every other state is D43#7's kws.pause - Listening's floor belongs to
// the ASR, Speaking's to the announcement, and a released engine answers
// silence, so the gate degrades instead of panicking.
//
// It runs on the consumer goroutine only (the one Accept caller, D38b); the hit
// handling happens inline because the D43 dispatches below are owner-goroutine
// work.
func (rv *residentVoice) kwsFeed(samples []float32) bool {
	if rv == nil || rv.kwsWake == nil {
		return false
	}
	if !rv.warmFeed.Load() && rv.machine.State() != statemachine.StateArmed {
		return false
	}
	rv.kwsFedSamples.Add(int64(len(samples)))
	var hits []speech.KeywordHit
	hits = rv.kwsWake.Accept(samples, hits)
	for _, h := range hits {
		rv.onKwsHit(h)
	}
	return true
}

// onKwsHit routes one keyword detection through the rows the frozen table
// already has. The veto half is the B1 voice-cancel channel: the interrupt
// channel 381-r1 reserved for KWS veto words fires first (Speaking has no EvVeto
// row - the announcement is stopped through the speak loop's interrupt check),
// and the two states WITH a row walk it. The shipped keyword table carries no
// veto word (kwsVetoWords is empty in internal/speech), so this branch is wired
// but cannot fire from this tree's model - said here rather than left implied.
//
// The wake half branches on where the machine really is:
//
//	Armed       row #7, EvWakeWord -> Listening, plus the chain load the row's
//	            speech.load-vad-asr side effect names (beginKwsWake).
//	Sleeping    row #5 first: the machine's own timers (#12) can land it here
//	            while the engine is still loaded, and #7 needs Armed - the
//	            re-arm walk is the table's own answer, then #7.
//	Warm/Settling with a live chain   the wire-r1 wake: open the feed, resume,
//	            borrow (wakeFromWarm). #29/#33's row exists for EvSummon and the
//	            warm wake dispatches it from inside.
//	Warm/Settling with no chain       the post-announce shape (the announce's
//	            family swap re-armed phase 1): beginKwsWake's Warm branch walks
//	            #29/#33 and starts the chain load.
//
// A hit in any other state is counted and dropped: the table has no row, and a
// leg that guesses an edge is worse than one that says so.
func (rv *residentVoice) onKwsHit(hit speech.KeywordHit) {
	if hit.Veto {
		rv.kwsVetos.Add(1)
		rv.TriggerInterrupt()
		switch rv.machine.State() {
		case statemachine.StateListening, statemachine.StateConfirming:
			rv.dispatch(statemachine.EvVeto, nil,
				"KWS veto 词命中：D43#13/#22（B1 语音取消通道，与唤醒词共用同一模型）")
		}
		return
	}
	rv.kwsHits.Add(1)
	if rv.asrLoaded.Load() && !rv.asrGone.Load() {
		// The chain is live (asrLoaded is reset by exactly the two branches that
		// retire a chain - the announce swap and the settling dispose), so the
		// wake is the wire-r1 one: open the feed, resume the session, take the
		// borrow back. This runs before the in-flight guard below ON PURPOSE:
		// handOffStarted stays true for the whole life of a live chain, and a
		// Warm hit must not be swallowed by the wake that started it.
		rv.wakeFromWarm("kws-wake")
		return
	}
	if rv.handOffStarted {
		// A chain load is already in flight for an earlier wake; a second hit
		// must not start a second walk over it.
		return
	}
	state := rv.machine.State()
	if state == statemachine.StateSleeping {
		rv.dispatch(statemachine.EvKwsEnabled, nil,
			"KWS 已装载而机器停在 Sleeping：D43#5 重新武装（机器自己的 #12 定时器会落这里）")
		state = rv.machine.State()
	}
	switch state {
	case statemachine.StateArmed:
		rv.beginKwsWake(hit, statemachine.EvWakeWord,
			"唤醒词命中：D43#7 Armed 到 Listening（KWS 暂停、ASR 按行装载）")
	case statemachine.StateWarm, statemachine.StateSettling:
		_, _ = rv.residency.resume("kws-wake")
		rv.beginKwsWake(hit, statemachine.EvSummon,
			"Warm 里唤醒词命中：D43#29/#33（表无 Warm+EvWakeWord 行，走 summon 的零装载唤醒）")
	default:
		// The table has no row out of this state for a wake (the wire-r1 named
		// Thinking-parking is the common one): the chain still restarts - the
		// pipeline work a wake justifies - and the ball stays where the frozen
		// table left it. Counted here, named in the evidence table.
		slog.Info("voice: wake-word hit in a state with no wake row; the chain restarts, the ball stays",
			"state", string(state), "word", hit.Word, "at_sample", hit.At)
		rv.beginKwsWake(hit, "", "")
	}
}

// beginKwsWake is the wake-word entry into the chain load: the turn opens at the
// hit (on the seam clock, the same clock every other turn reading uses), the row
// the caller chose is walked, and the walk-and-load machinery starts. An empty
// ev means the machine had no row to walk (the default branch above): the load
// starts anyway, and the ball stays where the frozen table put it.
func (rv *residentVoice) beginKwsWake(hit speech.KeywordHit, ev statemachine.Event, say string) {
	rv.handOffStarted = true
	rv.mu.Lock()
	rv.turns = append(rv.turns, voiceTurn{
		OpenSample:  rv.fed.Load(),
		StartSample: rv.fed.Load(),
	})
	rv.mu.Unlock()
	if ev != "" {
		rv.dispatch(ev, nil, say)
	}
	slog.Info("voice: wake-word hit opens the chain load",
		"word", hit.Word, "engine_at_sample", hit.At, "event", string(ev))
	rv.startChainLoad()
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
func (rv *residentVoice) adopt(r voiceChainReady) (*speech.Transcriber, speech.AsrEngine, error) {
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
	tr := speech.NewTranscriber(r.det, r.asr)
	// TICKET 27's call site. This is the one place a Transcriber exists in
	// production, so it is the only place SetPunctuator can be called from: the
	// driver punctuates in final() and nowhere else, which is what makes
	// "partials stay raw, finals are dictation grade" a property of the code.
	// A nil punct (no model resolved, or a stage that refused to load) is the
	// "switched off" state 27-r1 pinned: the final is the raw hypothesis and
	// Transcript.Punctuated says so, so nothing downstream has to guess.
	if p := rv.residency.punctuator(); p != nil {
		tr.SetPunctuator(p)
		rv.punctAttached.Store(true)
		slog.Info("voice: punctuation stage attached", "model", rv.plan.PunctModelPath,
			"attach_point", "Transcriber.SetPunctuator", "scope", "the ASR epoch (released with it)")
		fmt.Printf("wisp: 标点级已接上：ASR 的 final 过 %s 标点模型（%s），partial 结构性不带标点\n",
			"CT-Transformer", rv.plan.PunctModelPath)
	} else {
		// "no model on this machine" and "a model that would not attach" are
		// different facts: the first is a store without the punctuation entry and is
		// not a problem, the second is a stage this leg failed to wire and is. Only
		// the second goes through noteProblem.
		if rv.plan.PunctModelPath == "" {
			slog.Info("voice: no punctuation model resolved; finals stay raw (the switched-off state)")
			fmt.Printf("wisp: 标点级未接：本机模型库没有标点模型（这一进程的 final 是原始假设，不带标点）\n")
		} else {
			rv.noteProblem("标点级未接上（这一进程的 final 是原始假设，不带标点）：" + rv.plan.PunctModelPath)
		}
	}
	return tr, r.asr, nil
}

// transcriptSink wraps submitTask and adds EvInterrupt dispatch for veto words.
// It is passed to the Transcriber constructor so veto words can trigger interrupts.
// consume fans the seam's transcripts into this leg's readings, the D43 walk and
// the text path.
//
// TICKET 381 R26-1: also checks for KWS veto word interrupts。如果 VAD 停止触发时检测到
// 否决词（"停"/"取消"/"别"），则 dispatch EvInterrupt 而不是送任务流。
func (rv *residentVoice) consume(outs []speech.Transcript) {
	for _, t := range outs {
		// The Warm wake. A transcript after Warm proves the microphone has
		// something for the agent, whatever the D43 row that got us here was: the
		// feed opens, the borrow is taken back and the session resumes without
		// loading anything. Both kinds count, because a short utterance can reach
		// its final without ever emitting a partial.
		if rv.warmFeed.Load() && (t.Kind == speech.TranscriptPartial || t.Kind == speech.TranscriptFinal) {
			rv.wakeFromWarm("utterance-after-warm")
		}
		switch t.Kind {
		case speech.TranscriptPartial:
			rv.partials.Add(1)
			rv.recordPartial(t)
		case speech.TranscriptFinal:
			// 381-r1:播报中 Veto 词打断。简单的关键词匹配作为 KWS 模型的降级方案。
			isVeto := strings.Contains(t.Text, "停") || strings.Contains(t.Text, "取消") ||
				strings.Contains(t.Text, "别")
			if isVeto {
				slog.Info("voice: interrupted by veto word", "word", t.Text)
				rv.dispatch(statemachine.EvInterrupt, nil, "Veto 词打断：D43#27 回到 Listening（用户否决词）")
				continue
			}
			rv.finals.Add(1)
			if t.Punctuated {
				rv.punctFinals.Add(1)
			}
			rv.finish(t)
		default:
			rv.noteProblem(fmt.Sprintf("未知的转写类型 %v", t.Kind))
		}
	}
}

// wakeFromWarm is D43 #29 (Warm) / #33 (Settling) plus session Resume in one
// place, and it is the reading ticket 28 AC1 keys on: the acquire it takes is
// supposed to cost nothing, and the session's own LoadsAfterWarm counter refutes
// that claim if the warm path lost an engine. Since ticket 41 it has two
// callers: the transcript wake (a final after Warm - the shape that exists when
// no wake-word engine is resident) and the KWS wake (a keyword hit while the
// feed is closed); the reason names which.
func (rv *residentVoice) wakeFromWarm(reason string) {
	rv.warmFeed.Store(false)
	rv.warmWake.Add(1)
	resumed := false
	if ok, err := rv.residency.resume(reason); err == nil && ok {
		resumed = true
	}
	_ = resumed
	switch rv.machine.State() {
	case statemachine.StateWarm:
		rv.dispatch(statemachine.EvSummon, nil, "Warm 里听见一句：D43#29 Warm 到 Listening（零装载）")
	case statemachine.StateSettling:
		rv.dispatch(statemachine.EvSummon, nil, "Settling 里听见一句：D43#33 取消静置，回 Listening")
	}
	if _, d, err := rv.residency.acquireTurn(session.EngineASR, "warm-wake"); err != nil {
		rv.noteProblem("唤醒后借不到 ASR：" + err.Error())
	} else {
		slog.Info("voice: warm wake took its ASR borrow back", "acquire", d.String(),
			"stage", string(rv.residency.stage()), "loads_after_warm", rv.residency.stats().LoadsAfterWarm)
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

	// TICKET 28's turn end, and the line the whole warm池 exists for: the engines
	// are NOT stopped. The turn returns its borrow (Release never unloads) and the
	// session lands in Warm, which closes the mic half FIRST and then keeps every
	// engine hot for [session].warm_timeout_sec. Before wire-r1 the only thing
	// that ever released an engine was process exit; ticket 28-r1 could not wire
	// this because the seam is here.
	rv.residency.releaseTurn(session.EngineASR)
	if err := rv.residency.enterWarm("turn-end"); err == nil {
		rv.warmTurnEnd.Add(1)
		slog.Info("voice: turn ended into Warm (engines stay hot)",
			"reason", "turn-end", "stage", string(rv.residency.stage()),
			"loads_after_warm", rv.residency.stats().LoadsAfterWarm,
			"engines_hot", rv.residency.stats().EnginesHot)
	}
}

// dispatch walks one D43 edge and projects the result to the ball. An illegal
// transition is recorded, never panicked: the table is frozen and a leg that
// guesses an edge is worse than one that says the table refused it.
//
// It reports whether the edge was really walked. The announce path needs that
// answer: "the ball moved into Speaking" is a reading, and a hop that the table
// refused (an announcement arriving in a state with no legal row into Speaking,
// which is every state outside Thinking/Acting/AwaitingApproval) must not be
// reported as one.
func (rv *residentVoice) dispatch(ev statemachine.Event, facts *statemachine.Facts, say string) bool {
	if rv == nil || rv.machine == nil {
		return false
	}
	from := rv.machine.State()
	to, err := rv.machine.Dispatch(ev, facts)
	if err != nil {
		rv.noteProblem(fmt.Sprintf("D43 拒绝 %s + %s：%v", from, ev, err))
		return false
	}
	if to == from {
		return false
	}
	rv.mu.Lock()
	rv.timeline = append(rv.timeline, voiceTransition{
		From: from, To: to, Event: ev, AtSample: rv.fed.Load(), Say: say,
	})
	rv.mu.Unlock()
	rv.ball.setVoiceState(to)
	slog.Info("voice: D43 walked", "from", string(from), "event", string(ev), "to", string(to), "why", say)
	return true
}

// kwsLoaded is the truth source behind the Facts.KwsLoaded guard (rows #10 and
// #32): an engine object is resident or it is not. Before ticket 41 the fact was
// assigned nowhere in production; this is the one producer every dispatch that
// walks a KwsLoaded-guarded row reads.
func (rv *residentVoice) kwsLoaded() bool {
	return rv != nil && rv.kwsWake != nil
}

// kwsFacts is the guard input for those dispatches. Built per call, like every
// other Facts literal in this leg.
func (rv *residentVoice) kwsFacts() *statemachine.Facts {
	return &statemachine.Facts{KwsLoaded: rv.kwsLoaded()}
}

// onMuteGesture is the voice machine's half of the two mute gestures (the orb's
// global hot key and the tray's item). The gate flip itself is the capture leg's
// (ticket 290's ruling, unchanged); what ticket 41 adds is the D43 walk the
// rows that EXIST make legal: from Armed, #8 (Armed -> Muted, kws.stop-inference
// - the kwsFeed gate answers it, the engine object stays resident); from Muted,
// #10 with the real KwsLoaded fact (Armed when the engine is resident, Sleeping
// when it is not). The Sleeping posture stays dispatch-free exactly as 290
// ruled: the table has no Sleeping -> Muted edge, and the resting state of a
// boot without KWS is still Sleeping. The walk runs on the owner goroutine via
// the queue - the gesture fires on the ui-sta thread, and the ball hop is owner
// work (D38b).
func (rv *residentVoice) onMuteGesture() {
	if rv == nil || rv.machine == nil || rv.queue == nil {
		return
	}
	from := rv.machine.State()
	if from != statemachine.StateArmed && from != statemachine.StateMuted {
		return
	}
	loaded := rv.kwsLoaded()
	if !rv.queue.post(func() {
		rv.dispatch(statemachine.EvMuteKey, rv.kwsFacts(),
			"静音手势：D43#8/#10（guard=KwsLoaded，Muted 里 KWS 停推理留模型）")
	}) {
		slog.Warn("voice: mute gesture arrived after the owner queue closed; the D43 walk is dropped",
			"from", string(from), "kws_loaded", loaded)
	}
}

// drainWhilePaused drains one frame and drops it. It returns false when the
// stream is over, so a paused leg still exits with its ctx.
func (rv *residentVoice) drainWhilePaused(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	case <-rv.queue.wakeChan():
		rv.queue.run()
		return true
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
		"标点=%s（带标点的 final %d 条）；唤醒词=%s；播报=%s；%s；配置来源 %s",
		voiceGoroutineName, "voice-15e", rv.plan.VadModelPath, voiceAsrModelID, rv.plan.AsrDir,
		rv.drained.Load(), rv.partials.Load(), rv.finals.Load(), rv.handed.Load(),
		rv.handOffs.Load(), rv.handOffExit.Load(),
		punctState(rv.plan.PunctModelPath, rv.punctAttached.Load()), rv.punctFinals.Load(),
		rv.kwsReading(),
		rv.announceReading(), rv.residency.reading(), rv.plan.Source)
}

// kwsReading is the wake-word half of the posture line, spelled so the three
// "no KWS" shapes cannot be read as one: the opt-in never asked, the opt-in
// asked and the model did not resolve, and the engine that was asked for and
// refused to build (that one is in problems, where every other refusal is).
func (rv *residentVoice) kwsReading() string {
	if rv == nil {
		return "未装配"
	}
	if rv.kwsLoaded() {
		return fmt.Sprintf("已武装（Armed；命中 %d／veto %d／喂入 %d 样本；ASR 吃 %d 样本）",
			rv.kwsHits.Load(), rv.kwsVetos.Load(), rv.kwsFedSamples.Load(), rv.asrFedSamples.Load())
	}
	if rv.plan.KwsNote != "" {
		return "未武装（" + rv.plan.KwsNote + "）"
	}
	return "未武装"
}

// punctState spells the two "no punctuation" states apart, because they are
// different facts: no model on this machine, versus a model that is there and a
// stage that did not attach.
func punctState(model string, attached bool) string {
	switch {
	case attached:
		return "已接上（" + model + "）"
	case model == "":
		return "未接（本机模型库没有标点模型）"
	default:
		return "未接（有模型但没有接上：" + model + "）"
	}
}

// announceReading is the announce path's posture: every hop counted separately, so
// "it spoke" and "it was refused by the gate" cannot be read as the same thing.
func (rv *residentVoice) announceReading() string {
	if rv == nil {
		return "未装配"
	}
	if rv.residency == nil {
		return "未装配"
	}
	if !rv.plan.Tts.complete() && rv.ttsSpoken.Load() == 0 {
		return "没有 TTS 模型（这一进程不出声，结果只走 D10 通知）"
	}
	return fmt.Sprintf("闸=%d/%d（拦下 %d）引擎缺=%d 设备缺=%d 无内容=%d 失败=%d 出声=%d 句／%d 块／%d 样本 门闭=%d 放开=%d",
		rv.ttsGatedOk.Load(), rv.ttsGated.Load()+rv.ttsGatedOk.Load(), rv.ttsGated.Load(),
		rv.ttsNoEngine.Load(), rv.ttsNoDevice.Load(), rv.ttsSilent.Load(), rv.ttsFailed.Load(),
		rv.ttsSpoken.Load(), rv.ttsChunks.Load(), rv.ttsSamples.Load(),
		rv.speakGateOn.Load(), rv.speakGateOff.Load())
}

// setSpeakingCount is the gate needle's reading for tests and the boot report:
// on/off in one call, taken off the same atomics the stop line quotes. The
// two must be read together - "the gate was closed 3 times" says nothing
// without "it was reopened 3 times".
func (rv *residentVoice) setSpeakingCount() (closed, reopened uint64) {
	if rv == nil {
		return 0, 0
	}
	return rv.speakGateOn.Load(), rv.speakGateOff.Load()
}

// stopLine is the teardown record: the counters the exit line carries, taken off
// the same atomics the posture read.
func (rv *residentVoice) stopLine() string {
	if rv == nil || rv.scope == nil {
		return "voice leg: not assembled in this process"
	}
	return fmt.Sprintf("voice leg stopped: frames_drained=%d paused_drops=%d partials=%d finals=%d discarded=%d handed=%d handoffs=%d handoff_exit=%d asr_loaded=%v punct_attached=%v punctuated_finals=%d announce_asked=%d spoken=%d gated=%d samples=%d warm_turn_end=%d warm_wake=%d gate_closed=%d gate_reopened=%d kws_fed=%d kws_hits=%d kws_vetos=%d asr_fed=%d",
		rv.drained.Load(), rv.pausedDrops.Load(), rv.partials.Load(), rv.finals.Load(),
		rv.discarded.Load(), rv.handed.Load(), rv.handOffs.Load(), rv.handOffExit.Load(),
		rv.asrLoaded.Load(), rv.punctAttached.Load(), rv.punctFinals.Load(),
		rv.announceAsked.Load(), rv.ttsSpoken.Load(), rv.ttsGated.Load(), rv.ttsSamples.Load(),
		rv.warmTurnEnd.Load(), rv.warmWake.Load(), rv.speakGateOn.Load(), rv.speakGateOff.Load(),
		rv.kwsFedSamples.Load(), rv.kwsHits.Load(), rv.kwsVetos.Load(), rv.asrFedSamples.Load())
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
//
// 26r3 found and fixed the one defect that made the whole announce path
// unreachable in production: the literal below never built the owner queue,
// so rv.queue was nil on this - the PRODUCTION - constructor, announce()
// refused every sentence at its first line (the queue==nil guard), the owner
// select's queue case never fired, and the residency's clock marshalled
// session timers onto the runtime timer goroutine (voiceOwnerClock with a nil
// queue runs the callback inline). Every test that ever walked the announce
// path came through newResidentVoice, which is why the staged readings were
// green while the production shape was dead. The queue is built here, for
// every branch including the early returns, so a degraded boot still has the
// honest refused-post shape instead of a nil dereference.
func buildResidentVoice(dataDir string, ra *residentAudio, rb *residentBall, tasks voiceTranscriptSink) *residentVoice {
	rv := &residentVoice{
		ball:        rb,
		tasks:       tasks,
		queue:       newVoiceOwnerQueue(),
		machine:     statemachine.New(statemachine.Options{Initial: statemachine.StateSleeping}),
		provScope:   voiceAnnounceScope,
		newSpeaker:  openWaveOutSpeaker,
		interruptCh: make(chan struct{}, 1),
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
	plan.WarmTimeout = time.Duration(c.Session.WarmTimeoutSec) * time.Second
	if why != "" {
		rv.verdict = why
		return rv
	}
	rv.plan = plan
	rv.frames = ra.frames
	// The residency is built with the plan it will load from: the voicePlan the
	// constructor got is the zero one on this path, and a residency holding no
	// paths would answer "语音链：ASR 目录为空" for every turn.
	rv.residency = newVoiceResidency(plan, rv.queue, nil)
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
	// Ticket 41's opt-in gate: the wake word is off by default and the empty set
	// (with the reason spelled) is what the boot line reports. Only an explicit
	// [voice] wake_word.enabled=true resolves the model set.
	kws, kwsNote := voiceKwsFiles{}, ""
	if !c.Voice.WakeWord.Enabled {
		kwsNote = "[voice] wake_word.enabled=false（opt-in 没问，本进程停在 Sleeping）"
	} else {
		kws, kwsNote = voiceKwsFilesIn(store)
	}
	return voicePlan{
		VadModelPath: vad,
		AsrDir:       filepath.Join(store, voiceAsrModelID),
		HandOff:      residentVoiceHandOff(dataDir),
		// Ticket 27's model and ticket 26's checkpoint, resolved off the SAME store
		// layout 15e's ASR uses: <store>/<manifest id>/<file name inside the
		// archive>. Both are optional on purpose - a store without them is a chain
		// that cannot punctuate or speak, which the boot line says rather than a
		// chain that refuses to start.
		PunctModelPath: voicePunctModelPath(store),
		Tts:            voiceTtsFilesIn(filepath.Join(store, voiceTtsModelID)),
		// Ticket 41's wake-word set, resolved only when the opt-in asks for it:
		// an unasked-for engine would move the machine into Armed without the
		// user's consent, which is the exact posture the opt-in exists to gate.
		Kws:     kws,
		KwsNote: kwsNote,
	}, ""
}

// voiceKwsFilesIn resolves the KWS checkpoint's five files, in the order the
// other model families use: the store's own directory (the manifest's file
// spelling) first, then the WISP_KWS_DIR escape hatch (the same convention
// WISP_SILERO_MODEL / WISP_PUNC_MODEL / WISP_TTS_DIR follow for staging), then
// the engine package's own resolution (WISP_SPIKE_MODELS or the walk-up, which
// is what lets a plain `go test ./cmd/wisp` in a repository checkout find the
// dev copies without any env var). Every step requires all five files; the first
// complete set wins. The bytes are NOT re-verified here - the manifest's sha256
// pins are models.go's business, the same split the VAD already follows.
func voiceKwsFilesIn(store string) (voiceKwsFiles, string) {
	dir := filepath.Join(store, voiceKwsModelID)
	if f := voiceKwsFilesFrom(dir); f.complete() {
		return f, ""
	}
	if env := os.Getenv("WISP_KWS_DIR"); env != "" {
		if f := voiceKwsFilesFrom(env); f.complete() {
			return f, ""
		}
	}
	if cfg, err := speech.DefaultKwsConfig(); err == nil {
		f := voiceKwsFiles{
			Encoder:  cfg.EncoderPath,
			Decoder:  cfg.DecoderPath,
			Joiner:   cfg.JoinerPath,
			Tokens:   cfg.TokensPath,
			Keywords: cfg.KeywordsPath,
		}
		if f.complete() {
			return f, ""
		}
	}
	return voiceKwsFiles{}, fmt.Sprintf("KWS 模型五件不齐（%s 与 WISP_KWS_DIR 都不完整）", dir)
}

// voiceKwsFilesFrom fills the manifest spelling off one directory, checking each
// of the five files the spotter cannot open without.
func voiceKwsFilesFrom(dir string) voiceKwsFiles {
	f := voiceKwsFiles{
		Encoder:  filepath.Join(dir, voiceKwsEncoder),
		Decoder:  filepath.Join(dir, voiceKwsDecoder),
		Joiner:   filepath.Join(dir, voiceKwsJoiner),
		Tokens:   filepath.Join(dir, voiceKwsTokens),
		Keywords: filepath.Join(dir, voiceKwsKeywords),
	}
	for _, p := range []string{f.Encoder, f.Decoder, f.Joiner, f.Tokens, f.Keywords} {
		if !voiceFileExists(p) {
			return voiceKwsFiles{}
		}
	}
	return f
}

// voicePunctModelPath resolves the CT-Transformer file (manifest purpose
// "punctuation", id voicePunctModelID). WISP_PUNC_MODEL is the same escape hatch
// 27-r1's staging used and WISP_SILERO_MODEL already is for the VAD; it is an
// environment variable, not a configuration key (adding [voice] punctuation would
// move D36 and belongs to its own ticket).
func voicePunctModelPath(store string) string {
	p := filepath.Join(store, voicePunctModelID, voicePunctFile)
	if voiceFileExists(p) {
		return p
	}
	return os.Getenv("WISP_PUNC_MODEL")
}

// voiceTtsFilesIn resolves the matcha checkpoint's file set. WISP_TTS_DIR is the
// staging escape hatch; production is the store directory, whose layout is the
// archive's own (model-steps-3.onnx, lexicon.txt, tokens.txt, the three FSTs,
// dict/) plus the vocoder the manifest lists as a separate file.
func voiceTtsFilesIn(dir string) voiceTtsFiles {
	if env := os.Getenv("WISP_TTS_DIR"); env != "" {
		dir = env
	}
	if dir == "" {
		return voiceTtsFiles{}
	}
	files := voiceTtsFiles{
		Acoustic: filepath.Join(dir, voiceTtsAcoustic),
		Vocoder:  filepath.Join(dir, voiceTtsVocoder),
		Lexicon:  filepath.Join(dir, voiceTtsLexicon),
		Tokens:   filepath.Join(dir, voiceTtsTokens),
		DictDir:  filepath.Join(dir, voiceTtsDict),
	}
	// The vocoder is a SEPARATE download in the manifest, so a tree that keeps the
	// checkpoint's archive contents in their own directory (third_party's
	// spike-models carries matcha-icefall-zh-baker/ and vocos-22khz-univ.onnx side
	// by side) has it one level up. The store layout has both in <store>/<id>/, and
	// the first candidate wins there.
	if !voiceFileExists(files.Vocoder) {
		if up := filepath.Join(filepath.Dir(dir), voiceTtsVocoder); voiceFileExists(up) {
			files.Vocoder = up
		}
	}
	for _, fst := range voiceTtsRuleFsts {
		files.RuleFsts = append(files.RuleFsts, filepath.Join(dir, fst))
	}
	// A path that is not there is not a file set: the announce path refuses a
	// partial checkpoint with one named sentence instead of a sherpa refusal.
	for _, p := range []string{files.Acoustic, files.Vocoder, files.Lexicon, files.Tokens} {
		if !voiceFileExists(p) {
			return voiceTtsFiles{}
		}
	}
	if st, err := os.Stat(files.DictDir); err != nil || !st.IsDir() {
		return voiceTtsFiles{}
	}
	kept := files.RuleFsts[:0]
	for _, fst := range files.RuleFsts {
		if voiceFileExists(fst) {
			kept = append(kept, fst)
		}
	}
	files.RuleFsts = kept
	return files
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
