//go:build windows

package main

// The three things ticket wire-r1 assembles, in one file because they meet in
// one place: the engine RESIDENCE (internal/session's warm keep-alive plus
// internal/speech's single-family Slot), the ANNOUNCE path (risk.ChTTS gate ->
// TtsEngine.Speak -> a real output device -> the ball's Speaking), and the
// PUNCTUATION attach that makes an ASR final dictation-grade text.
//
// WHY THE RESIDENCE IS TWO OBJECTS AND NOT ONE. They answer different questions
// and the ticket that landed each one said so:
//
//   - internal/session.Scope is C31's "acquire-release through SessionScope
//     only": it owns how long an engine may stay HOT (Warm / 90 s / Settling /
//     dispose cascade). It has no per-kind unload - Release only drops a borrow -
//     so it cannot express "release the ASR because the TTS is about to load".
//   - internal/speech.Slot is D32 16.3.3's single-family rule: one family
//     resident, the old one CLOSED before the new loader runs, so the RSS peak is
//     max(ASR, TTS) and not their sum. 26-r1 measured the addition (436.6MB) and
//     the sum is what 700MB of budget does not survive.
//
// Composing them the other way round (a scope holding both families, a slot on
// top) would be the broken shape: the scope keeps a handle hot until dispose, so
// a slot that closed the ASR underneath it would hand a CLOSED engine back on the
// next Acquire. So the epochs are the composition: an epoch is one scope holding
// exactly ONE family (plus the punctuator, which belongs to the dictation path),
// and a family change is `dispose the epoch, then load the other family through
// the slot`. The slot's own release-before-load is what makes the ordering a
// property of the code rather than of the order these lines happen to run in.
//
// WHAT A FAMILY CHANGE COSTS, said out loud rather than hidden: an announcement
// releases the ASR's bytes, so the next dictation turn pays a cold ASR load
// (28-r1 measured 3.2 s on this host); and a dictation turn releases the TTS, so
// the next announcement pays the TTS load (26-r1 measured 5.8 s cold). That is
// D32 16.3.3's "on demand" fallback with its price attached, and the evidence
// table carries both numbers instead of claiming the swap is free.
//
// THREADING. Everything below runs on the voice chain's single owner goroutine
// (D38b) with one exception that the type system cannot express: the session
// scope's own TIMERS fire on a runtime timer goroutine, and the dispose cascade
// they run closes engines. Closing a sherpa session under a thread that is still
// decoding is the crash 15e's `stop` comment exists to prevent, so the scope is
// given a Clock whose callbacks are MARSHALLED onto the owner goroutine
// (voiceOwnerClock). The deadline is still the scope's own time.Duration and no
// wall clock is ever subtracted (D42#9); only the thread the callback lands on
// changes. That is why voiceResidency never starts a goroutine of its own.

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/session"
	"github.com/CarlosShao/wisp/internal/speech"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// voiceTtsFiles is the matcha checkpoint's file set, spelled as the artifact
// names inside an installed model directory (build/models/manifest.json, id
// tts-matcha-zh-baker). The vocoder is a separate download in the same
// directory: the manifest's second `files` entry has no archive, so both land in
// <store>/<id>/.
type voiceTtsFiles struct {
	Acoustic string
	Vocoder  string
	Lexicon  string
	Tokens   string
	DictDir  string
	RuleFsts []string
}

// complete reports whether every field TtsConfig requires is present. A partial
// set is not a TTS this leg may load: NewMatchaTts refuses an empty Vocoder, and
// a leg that tried anyway would spend a device open on a refusal.
func (f voiceTtsFiles) complete() bool {
	return f.Acoustic != "" && f.Vocoder != "" && f.Lexicon != "" && f.Tokens != "" && f.DictDir != ""
}

// voiceAnnounceScope is the C25 scope the announce gate reads. It is the
// process-lifetime scope this leg opens for its own channel, and it is NOT a
// shadow of a task's scope: risk.Provenance.CheckText fails CLOSED on an
// unregistered id while any other scope holds taints (provenance.go:694-711), so
// a leg that passed a task id it no longer owns would refuse every announcement
// for the wrong reason. What the shift costs is named in the evidence table: a
// task's own taints are dropped when that task's scope closes (risk.Scope.Close),
// so the taint that would matter for "do not read that aloud" is gone by the time
// the reply exists. Carrying it past the task is a change in internal/risk's
// scope lifetime, which is outside this leg's write surface.
const voiceAnnounceScope = "voice-announce"

// voiceAnnounceQueue is how many announcements may wait for the owner. One slot:
// the resident voice chain says one sentence at a time, and a second result
// arriving mid-sentence is refused out loud (counted) rather than queued into a
// backlog nobody asked for.
const voiceAnnounceQueue = 1

// voiceAnnounceMaxRunes bounds what an announcement may contain. Ticket 26's
// AC1 says "short result -> TTS audible": a 4000-rune transcript read aloud is
// minutes of speech, and the D10 notification is the delivery path for those.
// It is a package constant, not a configuration key (no [voice] key may be added
// by this leg; the punctuation switch needs its own ticket and D36).
const voiceAnnounceMaxRunes = 120

// ---- owner-goroutine marshalling ------------------------------------------

// voiceOwnerQueue posts work onto the voice chain's owner goroutine. It is the
// only bridge into that goroutine from anywhere else, and it starts nothing: the
// owner drains it from its own select, so an enqueue can never outlive the owner
// (a post after the owner returned is a counted drop, not a goroutine).
type voiceOwnerQueue struct {
	mu     sync.Mutex
	jobs   []func()
	wake   chan struct{}
	closed bool
}

func newVoiceOwnerQueue() *voiceOwnerQueue {
	return &voiceOwnerQueue{wake: make(chan struct{}, 1)}
}

// post enqueues fn and nudges the owner. It never blocks the caller.
func (q *voiceOwnerQueue) post(fn func()) bool {
	if q == nil || fn == nil {
		return false
	}
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return false
	}
	q.jobs = append(q.jobs, fn)
	q.mu.Unlock()
	select {
	case q.wake <- struct{}{}:
	default:
	}
	return true
}

// wakeChan is the case the owner's select waits on.
func (q *voiceOwnerQueue) wakeChan() <-chan struct{} {
	if q == nil {
		return nil
	}
	return q.wake
}

// run executes queued work on the owner goroutine, in order, until the queue is
// empty. It is called from the owner's select and from nowhere else.
func (q *voiceOwnerQueue) run() {
	if q == nil {
		return
	}
	for {
		q.mu.Lock()
		if len(q.jobs) == 0 {
			q.mu.Unlock()
			return
		}
		fn := q.jobs[0]
		q.jobs = q.jobs[1:]
		q.mu.Unlock()
		fn()
	}
}

// close stops accepting work. Called by the owner on its way out, so a late post
// is refused instead of growing a queue nobody will drain.
func (q *voiceOwnerQueue) close() {
	if q == nil {
		return
	}
	q.mu.Lock()
	q.closed = true
	q.jobs = nil
	q.mu.Unlock()
}

// voiceOwnerClock is session.Clock with one difference: the callback is run ON
// THE VOICE OWNER, not on the runtime timer goroutine. session's own timers are
// exactly right (they are the D42#9-clean way to express 90 s and 30 s), and the
// only thing this wrapper changes is which thread the dispose cascade lands on.
type voiceOwnerClock struct {
	queue *voiceOwnerQueue
	base  session.Clock // nil = the system clock
}

type voiceBaseTimer struct{ stop func() bool }

func (t voiceBaseTimer) Stop() bool {
	if t.stop == nil {
		return false
	}
	return t.stop()
}

func (c voiceOwnerClock) AfterFunc(d time.Duration, fn func()) session.Timer {
	post := func() {
		if c.queue != nil {
			c.queue.post(fn)
			return
		}
		fn()
	}
	if c.base != nil {
		return c.base.AfterFunc(d, post)
	}
	t := time.AfterFunc(d, post)
	return voiceBaseTimer{stop: t.Stop}
}

// ---- the engine handle ----------------------------------------------------

// voiceEngineHandle is one engine family as the session scope holds it. Its Close
// is what releases the family: it empties the SLOT (which closes the resident
// engine exactly once and counts it), so "the engine the scope released" and "the
// family the slot no longer holds" cannot disagree.
type voiceEngineHandle struct {
	kind  session.EngineKind
	slot  *speech.Slot
	asr   speech.AsrEngine
	tts   speech.TtsEngine
	punct speech.Punctuator

	residency *voiceResidency
	closed    atomic.Int64
}

func (h *voiceEngineHandle) Asr() speech.AsrEngine         { return h.asr }
func (h *voiceEngineHandle) Tts() speech.TtsEngine         { return h.tts }
func (h *voiceEngineHandle) Punctuator() speech.Punctuator { return h.punct }

// freeOSMemoryNow is the cascade's third step. It is a function value rather than
// a direct call so the hook reads as one of the four frozen steps instead of an
// inline import.
func freeOSMemoryNow() { debug.FreeOSMemory() }

// Close releases the family. C11: idempotent, and the second call moves nothing.
func (h *voiceEngineHandle) Close() error {
	if h == nil {
		return nil
	}
	if h.closed.Add(1) != 1 {
		return nil
	}
	var err error
	if h.punct != nil {
		if cerr := h.punct.Close(); cerr != nil {
			err = errors.Join(err, cerr)
		}
		h.punct = nil
	}
	if h.slot != nil {
		// Only the family this handle loaded may be released through the slot: a
		// handle that closed the OTHER family would trade a live session for a
		// silent one.
		want := speech.FamilyAsr
		if h.kind == session.EngineTTS {
			want = speech.FamilyTts
		}
		if h.slot.Resident() == want {
			if cerr := h.slot.Close(); cerr != nil {
				err = errors.Join(err, cerr)
			}
		}
	}
	h.asr, h.tts = nil, nil
	if h.residency != nil {
		h.residency.noteEngineClosed(h.kind)
	}
	return err
}

// ---- the residency -------------------------------------------------------

// voiceResidency owns the voice chain's engine lifetime: one epoch (one session
// scope) per resident family, the slot underneath it, and the counters every
// claim in the evidence table is read from.
type voiceResidency struct {
	plan  voicePlan
	queue *voiceOwnerQueue
	clock session.Clock

	mu    sync.Mutex
	slot  *speech.Slot
	scope *session.Scope
	epoch int
	kind  session.EngineKind
	eng   *voiceEngineHandle

	// Counters. The two that matter are LastAcquire (0 when the warm keeps the
	// engine) and FamilyLoads/FamilyCloses (one of each per swap).
	epochs       int
	familyLoads  int
	familyCloses int
	coldLoads    int
	warmEntries  int
	warmResumes  int
	micCloses    int
	disposals    int
	lastAcquire  time.Duration
	acquires     []time.Duration
	stageAtClose []string
	problems     []string

	// asrLoader is the ASR factory. Production is the sherpa streaming session;
	// it is a field so a case can drive the dictation path with an engine that
	// returns known text (27-r1 measured its punctuation examples that way, and the
	// wiring under test here - the epoch, the slot, SetPunctuator, the final's text
	// reach - is the production one either way).
	asrLoader func() (speech.AsrEngine, error)
	// loaderCalls counts how many times the factory above ran. It is the strong form
	// of "the second turn loaded nothing": a timing can be fast by accident, a
	// counter cannot.
	loaderCalls int

	// onMicClose is the mic half of "Warm closes the microphone FIRST". This leg's
	// microphone is not a device it may switch off: the C8 frames that carry the
	// wake in are the same frames the capture leg owns, and a chain that closed
	// the device in Warm would never hear the next utterance (ticket 41's KWS,
	// which is the only other wake source, is not implemented). So the hook closes
	// the FEED - what reaches the agent - and the evidence table says exactly that
	// (see the 采集门 row) instead of claiming a device close it did not do.
	onMicClose func()
}

func newVoiceResidency(plan voicePlan, queue *voiceOwnerQueue, clock session.Clock) *voiceResidency {
	r := &voiceResidency{plan: plan, queue: queue, clock: clock, slot: speech.NewSlot()}
	if clock == nil {
		clock = voiceOwnerClock{queue: queue}
	}
	r.clock = clock
	return r
}

func (r *voiceResidency) noteProblem(why string) {
	if r == nil {
		return
	}
	r.mu.Lock()
	r.problems = append(r.problems, why)
	r.mu.Unlock()
}

func (r *voiceResidency) problemsSnapshot() []string {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.problems...)
}

func (r *voiceResidency) noteEngineClosed(kind session.EngineKind) {
	r.mu.Lock()
	r.familyCloses++
	r.mu.Unlock()
}

// loader is the session scope's only way to bring an engine in. It is where C31
// and D32 16.3.3 meet: the scope asked for a kind, and the family that has to go
// is released by the SLOT before the loader runs.
func (r *voiceResidency) loader(kind session.EngineKind) (session.Engine, error) {
	switch kind {
	case session.EngineASR:
		if r.plan.AsrDir == "" {
			return nil, errors.New("语音链：ASR 目录为空")
		}
		factory := r.asrLoader
		if factory == nil {
			factory = func() (speech.AsrEngine, error) {
				return speech.NewStreamingAsr(speech.AsrConfig{
					EncoderPath: filepath.Join(r.plan.AsrDir, voiceAsrEncoder),
					DecoderPath: filepath.Join(r.plan.AsrDir, voiceAsrDecoder),
					TokensPath:  filepath.Join(r.plan.AsrDir, voiceAsrTokens),
				})
			}
		}
		eng, err := r.slot.LoadAsr(func() (speech.AsrEngine, error) {
			r.mu.Lock()
			r.loaderCalls++
			r.mu.Unlock()
			return factory()
		})
		if err != nil {
			return nil, err
		}
		return &voiceEngineHandle{kind: kind, slot: r.slot, asr: eng, residency: r}, nil

	case session.EngineTTS:
		if !r.plan.Tts.complete() {
			return nil, fmt.Errorf("语音链：TTS 模型文件不全（%s）", r.plan.Tts.Acoustic)
		}
		files := r.plan.Tts
		eng, err := r.slot.LoadTts(func() (speech.TtsEngine, error) {
			return speech.NewMatchaTts(speech.TtsConfig{
				AcousticModel: files.Acoustic,
				Vocoder:       files.Vocoder,
				Lexicon:       files.Lexicon,
				Tokens:        files.Tokens,
				DictDir:       files.DictDir,
				RuleFsts:      files.RuleFsts,
			})
		})
		if err != nil {
			return nil, err
		}
		return &voiceEngineHandle{kind: kind, slot: r.slot, tts: eng, residency: r}, nil

	case session.EnginePunct:
		if r.plan.PunctModelPath == "" {
			return nil, errors.New("语音链：标点模型路径为空")
		}
		p, err := speech.OpenPunctuation(speech.PunctuationConfig{ModelPath: r.plan.PunctModelPath}, true)
		if err != nil {
			return nil, err
		}
		return &voiceEngineHandle{kind: kind, punct: p, residency: r}, nil
	}
	return nil, fmt.Errorf("语音链：本腿不装载引擎族 %q", string(kind))
}

// openScope creates one epoch's session scope. The hooks are the ordering seams
// the ticket asks for, each one observable rather than decorative:
//
//	MicClose    the Warm transition closes the feed BEFORE the stage lands
//	Logf        the [audit] family cmd/wisp already prints
//	FreeOSMemory the cascade's third step (C11 step 4)
func (r *voiceResidency) openScope(reason string) (*session.Scope, error) {
	hooks := session.Hooks{
		Logf: func(format string, args ...any) {
			slog.Info("voice[scope]: " + fmt.Sprintf(format, args...))
		},
		MicClose: func() {
			// Read the stage HERE: EnterWarm calls this hook before it lands the
			// transition, so a hook that sees Awake is the proof that "the mic was
			// closed FIRST" is an order and not a sentence.
			stage := "unknown"
			if r != nil {
				stage = string(r.stage())
			}
			r.mu.Lock()
			r.micCloses++
			r.stageAtClose = append(r.stageAtClose, stage)
			hook := r.onMicClose
			r.mu.Unlock()
			if hook != nil {
				hook()
			}
		},
		FreeOSMemory: func() { freeOSMemoryNow() },
	}
	return session.NewScope(session.ScopeOptions{
		Loader: r.loader,
		Clock:  r.clock,
		Hooks:  hooks,
		// [session].warm_timeout_sec, handed over as an OPTION exactly as
		// session's own comment asks: the package reads no configuration, and 0
		// means its default (90 s).
		WarmTimeout: r.plan.WarmTimeout,
	})
}

// acquireTurn returns the engine for one turn, loading it only when the epoch is
// not already the right family. The returned duration is the reading AC1 asks
// for: 0 when Warm kept the engine, the cold load when a family swap happened.
func (r *voiceResidency) acquireTurn(kind session.EngineKind, reason string) (*voiceEngineHandle, time.Duration, error) {
	if r == nil {
		return nil, 0, errors.New("语音链：没有引擎驻留")
	}
	start := time.Now()
	r.mu.Lock()
	scope, curKind, eng := r.scope, r.kind, r.eng
	r.mu.Unlock()
	if scope != nil && curKind == kind && eng != nil {
		// The warm path: one borrow on an engine that is already hot. A disposed
		// scope answers its own refusal here, which is the answer AC1 wants when
		// the 90 s window really did expire.
		if _, err := scope.Acquire(kind); err != nil {
			return nil, time.Since(start), err
		}
		d := time.Since(start)
		r.noteAcquire(d)
		return eng, d, nil
	}

	// A family change: the resident family must go before the other one comes in.
	if err := r.closeEpoch("family-change:" + string(kind) + ":" + reason); err != nil {
		r.noteProblem("上一族引擎没有干净释放：" + err.Error())
	}

	scope, err := r.openScope(reason)
	if err != nil {
		return nil, time.Since(start), err
	}
	loaded, err := scope.Acquire(kind)
	if err != nil {
		_ = scope.Dispose("load-failed:" + string(kind))
		return nil, time.Since(start), err
	}
	handle, _ := loaded.(*voiceEngineHandle)
	if handle == nil {
		_ = scope.Dispose("bad-handle:" + string(kind))
		return nil, time.Since(start), fmt.Errorf("语音链：scope 交回的不是本腿的引擎句柄（%T）", loaded)
	}

	// The punctuator belongs to the dictation path: it rides in the ASR epoch so
	// it is released by the same cascade and cannot outlive the ASR it annotates.
	// A punctuator that will not load is a degraded dictation path, never a failed
	// utterance (27-r1's AC4), so the failure is recorded and the turn goes on.
	if kind == session.EngineASR && r.plan.PunctModelPath != "" {
		peng, perr := scope.Acquire(session.EnginePunct)
		if perr != nil {
			r.noteProblem("标点级没起来（这一句不带标点）：" + perr.Error())
		} else if ph, ok := peng.(*voiceEngineHandle); ok && ph != nil {
			handle.punct = ph.punct
			// The borrow is returned at once: the handle keeps the punctuator, and a
			// reference nobody releases would show up as a leak in Stats().Refs.
			if rerr := scope.Release(session.EnginePunct); rerr != nil {
				r.noteProblem("标点级的借出没有归还：" + rerr.Error())
			}
		}
	}

	d := time.Since(start)
	r.mu.Lock()
	r.scope, r.kind, r.eng = scope, kind, handle
	r.epoch++
	r.familyLoads++
	r.epochs = r.epoch
	r.mu.Unlock()
	r.noteAcquire(d)
	return handle, d, nil
}

func (r *voiceResidency) noteAcquire(d time.Duration) {
	r.mu.Lock()
	r.lastAcquire = d
	r.acquires = append(r.acquires, d)
	r.mu.Unlock()
}

// acquireLog is every acquire this leg took, in order: [0] is the first turn's
// (cold when nothing was hot) and [1] is the one the evidence table quotes as
// "the second turn cost nothing". A slice, not just the last value, because the
// claim is a comparison and a comparison needs both halves.
func (r *voiceResidency) acquireLog() []time.Duration {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]time.Duration(nil), r.acquires...)
}

// loaderCallCount is how many times the engine factory ran. It is the strong form
// of "the warm path loaded nothing": a duration can be fast by luck.
func (r *voiceResidency) loaderCallCount() int {
	if r == nil {
		return 0
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.loaderCalls
}

// micCloseCount and micCloseStages are the "mic closed FIRST" reading: how many
// times the hook ran and what session.Stage it saw while running. Awake in that
// list is the proof that the close precedes the transition (EnterWarm's order).
func (r *voiceResidency) micCloseCount() int {
	if r == nil {
		return 0
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.micCloses
}

func (r *voiceResidency) micCloseStages() []string {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.stageAtClose...)
}

// waitAcquire returns the duration of the most recent acquire, for the boot line.
func (r *voiceResidency) waitAcquire() time.Duration {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.lastAcquire
}

// releaseTurn returns one borrow. It never unloads: releasing the last reference
// while the session is alive is the NORMAL path into Warm (session Release's own
// contract), which is the whole point of the epoch.
func (r *voiceResidency) releaseTurn(kind session.EngineKind) {
	if r == nil {
		return
	}
	r.mu.Lock()
	scope := r.scope
	r.mu.Unlock()
	if scope == nil {
		return
	}
	if err := scope.Release(kind); err != nil {
		r.noteProblem("归还引擎借出失败：" + err.Error())
	}
}

// closeEpoch disposes the current epoch: the session cascade closes the engines
// in C11's order and the slot is emptied by the handle's Close.
func (r *voiceResidency) closeEpoch(reason string) error {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	scope := r.scope
	r.scope, r.eng = nil, nil
	r.mu.Unlock()
	if scope == nil {
		return nil
	}
	r.mu.Lock()
	r.disposals++
	r.mu.Unlock()
	return scope.Dispose(reason)
}

// enterWarm ends the turn but keeps the engines hot: D43 #26/#31's engine half.
func (r *voiceResidency) enterWarm(reason string) error {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	scope := r.scope
	r.mu.Unlock()
	if scope == nil {
		return nil
	}
	if err := scope.EnterWarm(reason); err != nil {
		// Settling has no Warm step (session says so); saying it happened would be
		// worse than logging the refusal.
		r.noteProblem("落 Warm 被拒（" + reason + "）：" + err.Error())
		return err
	}
	r.mu.Lock()
	r.warmEntries++
	r.mu.Unlock()
	return nil
}

// resume serves the Warm wake. It loads NOTHING: that is the claim, and the
// session scope's own LoadsAfterWarm counter is what refutes it if it is false.
func (r *voiceResidency) resume(reason string) (bool, error) {
	if r == nil {
		return false, nil
	}
	r.mu.Lock()
	scope := r.scope
	r.mu.Unlock()
	if scope == nil {
		return false, nil
	}
	if st := scope.Stage(); st != session.StageWarm && st != session.StageSettling {
		return false, nil
	}
	if err := scope.Resume(); err != nil {
		r.noteProblem("唤醒被拒（" + reason + "）：" + err.Error())
		return false, err
	}
	r.mu.Lock()
	r.warmResumes++
	r.mu.Unlock()
	return true, nil
}

func (r *voiceResidency) stage() session.Stage {
	if r == nil {
		return session.StageDisposed
	}
	r.mu.Lock()
	scope := r.scope
	r.mu.Unlock()
	if scope == nil {
		return session.StageDisposed
	}
	return scope.Stage()
}

func (r *voiceResidency) stats() session.Stats {
	if r == nil {
		return session.Stats{}
	}
	r.mu.Lock()
	scope := r.scope
	r.mu.Unlock()
	if scope == nil {
		return session.Stats{}
	}
	return scope.Stats()
}

// punctuator is the stage the current epoch is holding, for the one call site that
// attaches it (Transcriber.SetPunctuator in the adopt path). nil means the epoch
// carries none: no model resolved, or a stage that refused to load.
func (r *voiceResidency) punctuator() speech.Punctuator {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.eng == nil {
		return nil
	}
	return r.eng.punct
}

// dispose is the process exit: D38(e) step 4's engine half.
func (r *voiceResidency) dispose(reason string) error {
	if r == nil {
		return nil
	}
	err := r.closeEpoch(reason)
	if r.slot != nil {
		if cerr := r.slot.Close(); cerr != nil {
			err = errors.Join(err, cerr)
		}
	}
	return err
}

// reading is the one line the boot report and the evidence table both quote.
func (r *voiceResidency) reading() string {
	if r == nil {
		return "引擎驻留：未装配"
	}
	st := r.stats()
	r.mu.Lock()
	epochs, loads, closes, cold, warmEntries, warmResumes, micCloses, disp :=
		r.epochs, r.familyLoads, r.familyCloses, r.coldLoads, r.warmEntries, r.warmResumes, r.micCloses, r.disposals
	last := r.lastAcquire
	stages := append([]string(nil), r.stageAtClose...)
	r.mu.Unlock()
	slotResident := "none"
	slotLoads, slotCloses := 0, 0
	if r.slot != nil {
		slotResident = string(r.slot.Resident())
		slotLoads, slotCloses = r.slot.Loads(), r.slot.Closes()
	}
	return fmt.Sprintf("引擎驻留：stage=%s epochs=%d 装载=%d 释放=%d 冷装载=%d 上次 acquire=%v Warm=%d 唤醒=%d 麦克风关=%d(%s) dispose=%d；slot=%s loads=%d closes=%d LoadsAfterWarm=%d EnginesHot=%d Refs=%d",
		string(st.Stage), epochs, loads, closes, cold, last, warmEntries, warmResumes, micCloses,
		strings.Join(stages, ","), disp, slotResident, slotLoads, slotCloses,
		st.LoadsAfterWarm, st.EnginesHot, st.Refs)
}

// ---- the announce path ----------------------------------------------------

// voiceAnnounceRequest is one sentence the process wants said out loud, plus the
// C25 scope it came from.
type voiceAnnounceRequest struct {
	text  string
	scope string
}

// announce asks the voice chain to say one sentence. It is called from any
// goroutine (the result presenter is the task goroutine) and returns immediately:
// the sentence is spoken by the chain's owner, which is the only goroutine that
// may touch a TtsEngine (D38b). false means "this leg took no sentence" - no
// voice chain, no TTS model, or a queue that is already full - and each of those
// is counted where it happens.
func (rv *residentVoice) announce(text string) bool {
	if rv == nil || rv.queue == nil || rv.residency == nil {
		return false
	}
	rv.announceAsked.Add(1)
	text = strings.TrimSpace(text)
	if text == "" {
		rv.announceEmpty.Add(1)
		return false
	}
	if n := len([]rune(text)); n > voiceAnnounceMaxRunes {
		rv.announceTooLong.Add(1)
		rv.noteProblem(fmt.Sprintf("这一句不播报：%d 字超过播报上限 %d 字（长文本走 D10 通知，不念出来）", n, voiceAnnounceMaxRunes))
		return false
	}
	scope := rv.provScope
	if scope == "" {
		scope = voiceAnnounceScope
	}
	req := voiceAnnounceRequest{text: text, scope: scope}
	rv.announceMu.Lock()
	if len(rv.announcePending) >= voiceAnnounceQueue {
		rv.announceMu.Unlock()
		rv.announceDropped.Add(1)
		return false
	}
	rv.announcePending = append(rv.announcePending, req)
	rv.announceMu.Unlock()
	if !rv.queue.post(func() { rv.speakQueued() }) {
		rv.announceDropped.Add(1)
		return false
	}
	return true
}

// speakQueued drains this leg's pending sentences on the owner goroutine and
// speaks them in order.
func (rv *residentVoice) speakQueued() {
	if rv == nil {
		return
	}
	for {
		rv.announceMu.Lock()
		if len(rv.announcePending) == 0 {
			rv.announceMu.Unlock()
			return
		}
		req := rv.announcePending[0]
		rv.announcePending = rv.announcePending[1:]
		rv.announceMu.Unlock()
		rv.speak(req)
	}
}

// speak is the announce path itself, and every hop the ticket asks for is a
// separate line below so the evidence table can quote one of them:
//
//  1. risk.ChTTS gate    CheckText before anything is synthesised
//  2. engine             the slot's TTS family (the ASR is released first)
//     2b. capture gate      SetSpeaking(true) once the family swap has released
//     the ASR, so the frames the gate closes behind are
//     already in flight and gate.go's invariant holds
//  3. Speak              the C9 seam, samples out in playback order
//  4. the device         waveOut queues them; Drain waits for the last one
//  5. the ball           D43 #15 into Speaking, then #26 out to Warm
//
// It runs on the owner goroutine, so the audio it plays is not fed to the ASR
// (D47's "the announcement audio must never reach ASR", which the paused and
// warm feeds both implement by not feeding while this runs). Since 26r3 the
// feed is not merely withheld: the capture gate is CLOSED for the whole body -
// including the hops that never reach the device (refused, no engine, no
// speaker, synthesis failure, silence) - and reopened exactly once on the way
// out, so the mic never hears the announcement and gate.go's bounded-channel
// invariant ("frames behind a closed gate were captured before the TTS
// started") is what the ring replay can lean on.
func (rv *residentVoice) speak(req voiceAnnounceRequest) {
	text := req.text
	// HOP 1: the C25 channel gate. A hit is not an error - it is the channel
	// working - so it is counted and reported by source, and NOTHING is
	// synthesised (ticket 26 AC7: "before playback").
	if rv.prov != nil {
		hit, refused := rv.prov.CheckText(req.scope, risk.ChTTS, text)
		if refused {
			rv.ttsGated.Add(1)
			rv.noteProblem("播报被 ChTTS 闸拦下（" + string(risk.ChTTS) + "，来源 " + hit.Source() + "）：这一句不出声")
			return
		}
	}
	rv.ttsGatedOk.Add(1)

	// HOP 2: the engine, through the single-family slot.
	handle, d, err := rv.residency.acquireTurn(session.EngineTTS, "tts.announce")
	if err != nil {
		rv.ttsNoEngine.Add(1)
		rv.noteProblem("播报没有引擎（" + d.String() + "）：" + err.Error())
		return
	}
	defer rv.residency.releaseTurn(session.EngineTTS)
	// The swap released the ASR the running transcriber still holds (D32 16.3.3:
	// the slot closes the resident family before the incoming loader runs). The
	// owner goroutine re-arms phase 1 on this flag, so the next utterance reloads
	// the ASR instead of feeding a released session.
	if rv.asrLoaded.Load() {
		rv.asrGone.Store(true)
	}
	// HOP 2b (ticket 26's one needle, 26r3): close the capture gate BEFORE the
	// first sample can exist. The placement is the placement, not a habit: the
	// family swap above may have Stop()ed the inner source mid-frame, and frames
	// already in the bounded channel were captured BEFORE this sentence began -
	// exactly the invariant internal/audio/gate.go documents for a closed gate.
	// Closing it after Speak's first chunk would be too late by one waveOut
	// buffer and the microphone would hear the announcement. gateRelease below
	// is the ONLY way this body leaves, so every outcome - spoken, refused
	// later, no engine, no device, synthesis failure, the engine's
	// nothing-to-say guard - reopens the gate exactly once.
	gateOpen := false
	gateRelease := func() {
		if gateOpen {
			gateOpen = false
			if rv.captureGate != nil {
				rv.captureGate.SetSpeaking(false)
				rv.speakGateOff.Add(1)
			}
			slog.Info("voice: capture gate reopened (speech ended)",
				"reason", "announce-done", "path", string(audio.PathT))
		}
	}
	defer gateRelease()
	if rv.captureGate != nil {
		rv.captureGate.SetSpeaking(true)
		rv.speakGateOn.Add(1)
		gateOpen = true
		slog.Info("voice: capture gate closed for the announcement",
			"reason", "speaking", "path", string(audio.PathT))
	}
	eng := handle.Tts()
	if eng == nil {
		rv.ttsNoEngine.Add(1)
		rv.noteProblem("播报没有引擎：scope 交回的句柄里没有 TTS")
		return
	}

	// HOP 5a: the ball walks into Speaking BEFORE the first sample, because the
	// orb is the answer to "is it talking" and a hop dispatched after the audio
	// would show it talking late. D43 #15 is the only legal way in from Thinking
	// and its guard is "first token, no tool call" - a pure-text result is the
	// definition of that.
	moved := rv.dispatch(statemachine.EvFirstToken,
		&statemachine.Facts{HasToolCall: false, HasAnnouncement: true},
		"播报：D43#15 首个纯文本 token 让球走 Speaking（TTS 说这一句）")

	// HOP 3 and 4: synthesis and the device.
	speaker, serr := rv.newSpeaker(eng.SampleRate())
	if serr != nil {
		rv.ttsNoDevice.Add(1)
		rv.noteProblem("播报没有输出设备：" + serr.Error() + "；这一句退回 D10 通知")
		if moved {
			rv.dispatch(statemachine.EvInterrupt, nil, "播报失败：D43#27 回到 Listening")
		}
		return
	}
	spoke := 0
	chunks := 0
	// TICKET 381 R26-1:播报中打断通道。每次播放一小段前检查 ev.interrupt，
	// 如果收到 EvInterrupt 则立即返回中断错误，导致 speak() 提前退出并
	// dispatch EvInterrupt (≤400ms)。非阻塞 select，不会卡住播放器线程。
	speakErr := eng.Speak(text, func(samples []float32) error {
		if len(samples) == 0 {
			return nil
		}
		chunks++
		spoke += len(samples)
		rv.ttsSamples.Add(uint64(len(samples)))

		// Check interrupt before playing each chunk. ≤400ms target from hotkey/click/word
		// to TTS pause: this select is in the innermost loop so the next iteration
		// sees the interrupt immediately. Non-blocking default ensures we don't stall.
		select {
		case <-rv.interruptCh:
			// Interrupt fired: return an error that will be treated as "interrupted" not "failed".
			return fmt.Errorf("interrupted mid-speak")
		default:
			// No interrupt yet: play this chunk normally.
		}

		return speaker.Play(samples, eng.SampleRate())
	})
	drainErr := speaker.Drain()
	_ = speaker.Close()
	rv.ttsChunks.Add(uint64(chunks))
	if speakErr != nil {
		// 381-r1: distinguish "interrupted" from "failed"。中断是用户意图（热键/单击/否决词），
		// 失败是技术故障（引擎问题）。中断时只 dispatch EvInterrupt，不报 problem。
		if strings.Contains(speakErr.Error(), "interrupted") {
			if moved {
				rv.dispatch(statemachine.EvInterrupt, nil, "播报被打断：D43#27 回到 Listening（用户触发）")
			}
			return
		}
		if errors.Is(speakErr, speech.ErrNothingToSpeak) {
			// The engine's own guard: empty or punctuation-only text is not a
			// silent success (26-r1's M2 case). Nothing was emitted, so nothing
			// is claimed.
			rv.ttsSilent.Add(1)
			rv.noteProblem("这一句没有可说内容（空或纯标点）：引擎拒绝合成，不播一段近静音")
		} else {
			rv.ttsFailed.Add(1)
			rv.noteProblem("播报合成失败：" + speakErr.Error())
		}
		if moved {
			rv.dispatch(statemachine.EvInterrupt, nil, "播报中断：D43#27 回到 Listening")
		}
		return
	}
	if drainErr != nil {
		rv.ttsFailed.Add(1)
		rv.noteProblem("播报没播完：" + drainErr.Error())
	}
	rv.ttsSpoken.Add(1)
	slog.Info("voice: announcement spoken",
		"runes", len([]rune(text)), "chunks", chunks, "samples", spoke,
		"rate", eng.SampleRate(), "acquire", d.String(), "scope", req.scope)

	// HOP 5b: the sentence is done, so the ball leaves Speaking by the row that
	// exists for it (D43 #26: not-conversation) and the session lands in Warm.
	rv.dispatch(statemachine.EvSpeakDone,
		&statemachine.Facts{ConversationMode: false},
		"播报结束：D43#26 Speaking 到 Warm（会话落 Warm，引擎保持热）")
	if err := rv.residency.enterWarm("announce-end"); err == nil {
		rv.warmTurnEnd.Add(1)
	}
}
