package session

// ---------------------------------------------------------------------------
// SessionScope (C31) - the keep-alive owner (ticket 28)
//
// doc.go has listed four responsibilities for this package since the boundary
// was frozen, and until this file only the last one (the D45 grants ledger) had
// an implementation: ticket 166 brought the cross-restart conversation outlets
// and ticket 224 the session identity, but "the wake -> explicit end / 90s-idle
// lifecycle" and "model reference counting" were prose. This file is those two.
//
// The shape is deliberately small and dependency-free (stdlib only), because
// every other rule in the ticket is an ordering rule, and an ordering rule that
// lives in a package which owns audio devices or onnxruntime sessions can only
// be tested by those devices:
//
//	Warm         models HOT, mic CLOSED, nothing unloaded. The next wake must do
//	             ZERO loads (B4) - that is the whole point of the state, and
//	             Stats().LoadsAfterWarm is the counter that can go red.
//	Settling     3 s static window (D43 #32) - engines are STILL hot here, which
//	             is why a wake during it cancels the settle (D43 #33) with no
//	             load either.
//	Dispose      one ordered cascade: models unload, panel hide-or-destroy,
//	             FreeOSMemory, grants expiry (D43 #31). Each step is a hook so
//	             the order is an assertion instead of a comment.
//	Conversation Path C's host-side gating: default OFF, first enable costs an
//	             L2-grade privacy confirm that is LOGGED, later enables do not
//	             prompt, and 30 s of no speech falls back to Warm with the mic
//	             closed FIRST.
//
// What this file does not do, on purpose:
//
//   - no state transitions. The D43 table lives in internal/statemachine and
//     this type is DRIVEN by its effects (session.scope-create,
//     session.zero-load-resume, session.dispose-scope). Keeping the two apart is
//     what lets the timer values here (90 s / 30 s / 3 s) be the SAME numbers
//     statemachine.DefaultTimeouts carries without either package importing the
//     other's guts.
//   - no audio device ownership, no onnxruntime. Engines arrive through Loader
//     and leave through Engine.Close; that is the entire surface.
//   - no configuration key. WarmTimeout is an OPTION because the production
//     composition root already has the key ([session] warm_timeout_sec,
//     internal/config/schema.go:188, default 90) - this package reads nothing
//     itself and the battery-saver pair below is a constant of the power policy
//     (D42#6), not of a file.
//   - no wall-clock arithmetic. Every deadline is a Clock timer; the only
//     duration this type ever computes is "how long the timer was asked to wait".
//     Tests inject a fake clock, so the 90 s and 30 s rows are exercised in
//     microseconds without sleeping (D42#9's discipline, and the reason the
//     timer surfaces here rather than in a test's own select).
// ---------------------------------------------------------------------------

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// The five engine families whose lifetime the scope owns (SPEC-05 section 5,
// C31: "ASR/TTS/VAD/punct/KWS acquire-release through SessionScope only"). They
// are strings rather than an iota because they reach logs and audit lines, and
// a renumbered constant is a silent change in what an old line means.
type EngineKind string

const (
	EngineASR   EngineKind = "asr"
	EngineTTS   EngineKind = "tts"
	EngineVAD   EngineKind = "vad"
	EnginePunct EngineKind = "punct"
	EngineKWS   EngineKind = "kws"
)

// engineOrder is the order engines are closed in at dispose. The cascade's own
// first step ("models unload") has no per-engine order in D43 #31, so this is
// the ticket's enumeration order frozen into one slice - one list, so two
// readers cannot disagree about what "all models unloaded" means.
var engineOrder = []EngineKind{EngineASR, EngineTTS, EngineVAD, EnginePunct, EngineKWS}

// Engine is one loaded engine handle. Close RELEASES it (C11): sessions are
// released, never unloaded, which is why Warm can keep the handle hot and why
// Close is called at dispose and nowhere else.
type Engine interface {
	Close() error
}

// Loader loads one engine kind. It is the only way an engine enters the
// process: the scope calls it, the scope owns the result, and no caller ever
// receives a handle it may keep past Release.
//
// An error is returned verbatim and no slot is created, so a missing engine
// (TTS and KWS are reserved slots until tickets 26/41) fails the borrow rather
// than half-creating a session.
type Loader func(kind EngineKind) (Engine, error)

// Timer is the one method this package needs from a timer.
type Timer interface {
	// Stop cancels the timer. It reports whether the call stopped a pending
	// fire, exactly like time.Timer.Stop.
	Stop() bool
}

// Clock is the timer source. The zero value is not usable; NewScope substitutes
// systemClock. It exists so the 90 s / 30 s / 3 s rows can be driven by a test
// in microseconds and so no deadline in this file is ever "how long did that
// take" (D42#9).
type Clock interface {
	AfterFunc(d time.Duration, fn func()) Timer
}

type systemClock struct{}

func (systemClock) AfterFunc(d time.Duration, fn func()) Timer { return time.AfterFunc(d, fn) }

// The frozen timeout values. The first pair is [session] warm_timeout_sec's
// default (90) in its battery-saver shape (30); the second pair is D43 #12's
// Conversation row (30 s) in its battery shape (15); the third is D43 #32's
// static settling window (3 s), which has no battery variant because it is a
// teardown grace, not an idle wait.
const (
	WarmTimeoutDefault      = 90 * time.Second
	WarmTimeoutBattery      = 30 * time.Second
	ConversationIdleDefault = 30 * time.Second
	ConversationIdleBattery = 15 * time.Second
	SettleStaticDefault     = 3 * time.Second
)

// Stage is the within-booot lifecycle position. It mirrors, but does not
// replace, the D43 states: the state machine decides Listening/Speaking; the
// scope only tracks how long the session's engines may stay hot.
type Stage string

const (
	// StageAwake is a live turn: Listening / Thinking / Acting / Speaking all
	// live here. Engines may be borrowed.
	StageAwake Stage = "Awake"
	// StageWarm is keep-alive: engines HOT, mic CLOSED, no timeout shorter than
	// the warm timeout. A wake from here loads nothing (D43 #29).
	StageWarm Stage = "Warm"
	// StageSettling is D43 #32's 3 s static window. Engines are still hot;
	// a wake cancels the settle (D43 #33).
	StageSettling Stage = "Settling"
	// StageDisposed is the terminal position: the cascade has run, every engine
	// is closed, and no borrow is possible.
	StageDisposed Stage = "Disposed"
)

// Hooks are the ordering seams. Every one of them is optional (a nil hook is an
// explicit no-op, never a silent skip of the step) and every one of them is
// called OUTSIDE the scope's mutex, so a hook may log, block on a window
// message or call back into Stage() without deadlocking.
//
// They exist as hooks rather than as imports for two reasons: the panel window
// is ticket 33's, the grant invalidation is ticket 49's, and both would drag a
// compile-time edge into this package for a call it does not own; and the D43
// #31 cascade is an ORDER, which is only provable if each step is separately
// observable.
type Hooks struct {
	// MicOpen / MicClose are the capture device's gate. MicClose is called
	// BEFORE the Warm transition lands - the ticket's "close mic FIRST" - so a
	// hook that reads Stage() during the close still sees Awake.
	MicOpen  func()
	MicClose func()

	// PanelHide is the cascade's second step: hide-or-destroy the panel window
	// (D43 #31). Which of the two is ticket 33's business.
	PanelHide func()

	// FreeOSMemory is the cascade's third step (C11 step 4). The default is
	// debug.FreeOSMemory, supplied by the composition root so this package does
	// not pull runtime/debug into every test.
	FreeOSMemory func()

	// GrantsExpire is the cascade's fourth step. It has no body today BY DESIGN:
	// the D45 rows are keyed by the minted session identity, and the identity
	// dies with the process (grants.go:33-37, A435 clause 2), so "expire" is
	// "stop using the key". The hook exists so the step has a position in the
	// order that ticket 49 can fill without moving anything.
	GrantsExpire func()

	// PrivacyConfirm is the L2-grade native card the first Conversation enable
	// must pass (ticket 28 AC4). Returning false leaves Conversation OFF and is
	// logged. A nil PrivacyConfirm is treated as "no surface to ask on", which
	// fails closed - enabling Conversation without a card would be the silent
	// bypass SPEC-06 exists to prevent.
	PrivacyConfirm func() bool

	// WarmWindow is called once each time the session lands in Warm. The warm
	// window is the ONLY idle time the host has inside a live session, so it is
	// where tickets 29 (async L1 extraction) and 10 (history compression) do
	// their work: this hook is their entry point, and the executor stays theirs.
	// The scope never runs the work itself and never waits for it (D38b: the
	// work has an owner or it does not run).
	WarmWindow func()

	// Logf is the audit sink (the "[audit]" family cmd/wisp already uses). A
	// nil Logf is a no-op; the scope never formats a line nobody asked for.
	Logf func(format string, args ...any)
}

// Stats is the instrumentation AC1 asks for: a counter that can go RED when the
// warm path stops being warm.
//
//	LoadsAfterWarm  every engine load that happened AFTER the first EnterWarm.
//	                B4's claim is that this is 0 for a Warm wake, so a wake that
//	                reloads shows up here as 1 or more - no timing needed.
//	WarmResumes     how many Warm wakes this scope served.
//	EnginesHot      engines currently held (hot in Warm too).
//	Refs            live borrows (acquire minus release).
type Stats struct {
	Stage          Stage
	ColdLoads      int
	LoadsAfterWarm int
	WarmResumes    int
	EnginesHot     int
	Refs           int
	Conversation   bool
	PanelReserved  bool
	Disposals      int
	Prompts        int
}

// slot is one engine the scope holds. refs is the borrow count; the handle
// stays hot at refs == 0 until dispose, which is exactly what Warm means.
type slot struct {
	eng  Engine
	refs int
}

// ScopeOptions configures a Scope. Loader is the only required field.
type ScopeOptions struct {
	// Loader loads engines. Required: a scope that cannot load anything could
	// never answer a borrow, and a nil loader would turn that into a crash.
	Loader Loader
	// Clock is the timer source. nil = the system clock.
	Clock Clock
	// Hooks are the ordering seams. The zero value is all no-ops.
	Hooks Hooks
	// WarmTimeout is [session] warm_timeout_sec. 0 = WarmTimeoutDefault (90 s).
	WarmTimeout time.Duration
	// ConversationIdle is D43 #12's Conversation row. 0 = 30 s.
	ConversationIdle time.Duration
	// SettleStatic is D43 #32's static window. 0 = 3 s.
	SettleStatic time.Duration
	// BatterySaver selects D42#6's shapes: Warm 30 s / Conversation 15 s and
	// no KWS residency. It overrides the two idle values above when set, because
	// "battery saver" is a power-policy fact and not another timeout to tune.
	BatterySaver bool
}

// Scope is the session's keep-alive owner. Every method is safe for concurrent
// use; the mutex is the state.
//
// It starts no goroutines of its own: its only asynchronous surface is the
// Clock's timers, each of whose callbacks re-enters through the same mutex and
// is guarded by a generation counter, so a stale fire (one that lost a race
// with Stop) is a no-op rather than a state change.
type Scope struct {
	mu     sync.Mutex
	loader Loader
	clock  Clock
	hooks  Hooks

	warmTimeout      time.Duration
	conversationIdle time.Duration
	settleStatic     time.Duration

	stage          Stage
	slots          map[EngineKind]*slot
	everWarm       bool
	coldLoads      int
	loadsAfterWarm int
	warmResumes    int
	disposals      int
	prompts        int
	conversation   bool
	confirmed      bool
	panelReserved  bool

	warmTimer   Timer
	convTimer   Timer
	settleTimer Timer
	warmGen     uint64
	convGen     uint64
	settleGen   uint64
}

// NewScope opens one session's scope at StageAwake: a session exists because
// something woke it (D43 #4/#7/#33), and minting one in any other position
// would be a lifecycle nobody dispatched.
func NewScope(o ScopeOptions) (*Scope, error) {
	if o.Loader == nil {
		return nil, errors.New("session: NewScope 需要一枚 Loader（没有装载器就没有引擎所有权）")
	}
	clock := o.Clock
	if clock == nil {
		clock = systemClock{}
	}
	warm := o.WarmTimeout
	if warm <= 0 {
		warm = WarmTimeoutDefault
	}
	conv := o.ConversationIdle
	if conv <= 0 {
		conv = ConversationIdleDefault
	}
	settle := o.SettleStatic
	if settle <= 0 {
		settle = SettleStaticDefault
	}
	if o.BatterySaver {
		warm = WarmTimeoutBattery
		conv = ConversationIdleBattery
	}
	return &Scope{
		loader:           o.Loader,
		clock:            clock,
		hooks:            o.Hooks,
		warmTimeout:      warm,
		conversationIdle: conv,
		settleStatic:     settle,
		stage:            StageAwake,
		slots:            make(map[EngineKind]*slot),
	}, nil
}

// Stage reports the current lifecycle position.
func (s *Scope) Stage() Stage {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stage
}

// Stats snapshots the instrumentation. It never loads anything.
func (s *Scope) Stats() Stats {
	s.mu.Lock()
	defer s.mu.Unlock()
	st := Stats{
		Stage:          s.stage,
		ColdLoads:      s.coldLoads,
		LoadsAfterWarm: s.loadsAfterWarm,
		WarmResumes:    s.warmResumes,
		EnginesHot:     len(s.slots),
		Conversation:   s.conversation,
		PanelReserved:  s.panelReserved,
		Disposals:      s.disposals,
		Prompts:        s.prompts,
	}
	for _, sl := range s.slots {
		st.Refs += sl.refs
	}
	return st
}

// ReservePanel records that this session owns a panel window, which is the
// ticket's "panel keep-alive reservation": the window's lifetime is the
// session's, so it must not outlive StageDisposed. The window itself is ticket
// 33's (this package never imports internal/panel); the reservation is what
// makes "whose window is this" answerable, and the cascade's panel step is what
// makes it true.
//
// Reserving twice is not an error (a re-shown window is the same window).
func (s *Scope) ReservePanel() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.stage == StageDisposed {
		return errors.New("session: 会话已结束，不能再为面板预留寿命")
	}
	s.panelReserved = true
	return nil
}

// PanelReserved reports whether a panel window is reserved by this session.
func (s *Scope) PanelReserved() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.panelReserved
}

// Acquire borrows one engine, loading it only when the scope does not already
// hold it. The returned handle is valid until Release (and, after dispose,
// never); the handle is not a reference the caller may keep.
//
// The load happens under the mutex. That is deliberate: the voice pipeline is
// one owner goroutine (D38b), so contention is not a real cost, while a load
// outside the lock would let two goroutines race the same missing engine and
// produce two handles - the "no other module may keep models alive" hole this
// method exists to close.
func (s *Scope) Acquire(kind EngineKind) (Engine, error) {
	if !knownEngine(kind) {
		return nil, fmt.Errorf("session: 未知引擎类型 %q（C31 名册只有 asr/tts/vad/punct/kws）", string(kind))
	}
	loaded := false
	s.mu.Lock()
	if s.stage == StageDisposed {
		s.mu.Unlock()
		return nil, fmt.Errorf("session: 会话已结束，拒绝再借引擎 %s", string(kind))
	}
	sl := s.slots[kind]
	if sl == nil {
		eng, err := s.loader(kind)
		if err != nil {
			s.mu.Unlock()
			return nil, fmt.Errorf("session: 装载引擎 %s 失败: %w", string(kind), err)
		}
		if eng == nil {
			s.mu.Unlock()
			return nil, fmt.Errorf("session: Loader 对 %s 返回了 nil 引擎和 nil 错误（这项能力今天不可用）", string(kind))
		}
		s.slots[kind] = &slot{eng: eng}
		s.coldLoads++
		if s.everWarm {
			s.loadsAfterWarm++
			loaded = true
		}
		sl = s.slots[kind]
	}
	sl.refs++
	afterWarm := s.loadsAfterWarm
	s.mu.Unlock()
	if loaded {
		// B4's failure mode, said out loud: a wake that had to reload an engine
		// is exactly what a broken warm path looks like from the counter's side.
		s.log("session: SCOPE-COLD-LOAD-AFTER-WARM kind=%s loads=%d - 热路径掉了一枚引擎（B4 会红）",
			string(kind), afterWarm)
	}
	return sl.eng, nil
}

// Release returns one borrow. It never closes the engine: releasing the last
// reference while the session is alive is the NORMAL path into Warm, and
// unloading there would be exactly the reload B4 forbids.
func (s *Scope) Release(kind EngineKind) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	sl := s.slots[kind]
	if sl == nil || sl.refs == 0 {
		return fmt.Errorf("session: Release(%s) 没有对应的借出（重复归还或从未借出）", string(kind))
	}
	sl.refs--
	return nil
}

// EnterWarm ends the turn but keeps the session alive: the mic is closed FIRST
// and every engine stays hot. The order is the assertion in the ticket ("30s
// Conversation -> Warm with mic closed FIRST"), so it is one function for both
// the ordinary Warm row and the Conversation fallback.
//
// A call from StageWarm is idempotent (the warm timer restarts, the mic is
// already closed).
func (s *Scope) EnterWarm(reason string) error {
	s.mu.Lock()
	if s.stage == StageDisposed {
		s.mu.Unlock()
		return errors.New("session: 会话已结束，不能进入 Warm")
	}
	if s.stage == StageSettling {
		s.mu.Unlock()
		return errors.New("session: Settling 里没有 Warm 这一步（要么被唤醒取消，要么过期 dispose）")
	}
	s.stopConvTimerLocked()
	wasConversation := s.conversation
	s.conversation = false
	mute := s.hooks.MicClose
	s.mu.Unlock()

	// Step 1, outside the lock: close the mic. A hook that reads Stage() here
	// still sees Awake, which is what makes "closed FIRST" observable.
	if mute != nil {
		mute()
	}

	s.mu.Lock()
	s.stage = StageWarm
	s.everWarm = true
	s.armWarmTimerLocked()
	warmWork := s.hooks.WarmWindow
	s.mu.Unlock()

	// The warm window is the host's only idle time inside a live session, so the
	// async work tickets 29/10 own starts HERE and not on a timer of its own.
	if warmWork != nil {
		warmWork()
	}

	if wasConversation {
		s.log("session: conversation 落回 Warm（30s 无语音，麦克风已先关）reason=%s", reason)
	}
	s.log("session: SCOPE-WARM reason=%s engines_hot=%d warm_timeout=%s", reason, s.Stats().EnginesHot, s.warmTimeout)
	return nil
}

// Resume serves the Warm wake (D43 #29). It loads NOTHING: it cancels the warm
// timer, returns to Awake and uses the handles that are already hot. The
// counter AC1 keys on is checked here rather than asserted in prose - if any
// load happens after this call, Stats().LoadsAfterWarm says so.
func (s *Scope) Resume() error {
	s.mu.Lock()
	if s.stage == StageDisposed {
		s.mu.Unlock()
		return errors.New("session: 会话已结束，不能唤醒")
	}
	if s.stage != StageWarm && s.stage != StageSettling {
		s.mu.Unlock()
		return fmt.Errorf("session: Resume 只服务 Warm/Settling 的唤醒（当前 %s）", string(s.stage))
	}
	s.stopWarmTimerLocked()
	s.stopSettleTimerLocked()
	s.stage = StageAwake
	s.warmResumes++
	resumes := s.warmResumes
	hot := len(s.slots)
	s.mu.Unlock()

	s.log("session: SCOPE-WARM-RESUME resumes=%d engines_hot=%d loads_after_warm=%d",
		resumes, hot, s.Stats().LoadsAfterWarm)
	return nil
}

// EnableConversation turns the Conversation mode on. The first enable in the
// life of the scope MUST pass PrivacyConfirm (an L2-grade native card) and is
// logged; every later enable passes without a prompt (ticket 28 AC4).
//
// It fails closed when the confirm hook is nil: a build with no card surface
// must not silently enable a mode that keeps a microphone open.
func (s *Scope) EnableConversation() (bool, error) {
	s.mu.Lock()
	if s.stage == StageDisposed {
		s.mu.Unlock()
		return false, errors.New("session: 会话已结束，不能开启 Conversation")
	}
	if s.conversation {
		s.mu.Unlock()
		return true, nil
	}
	first := !s.confirmed
	confirm := s.hooks.PrivacyConfirm
	s.mu.Unlock()

	if first {
		if confirm == nil {
			s.log("session: conversation 首次开启被拒（没有隐私确认面，fail-closed）")
			return false, errors.New("session: 首次开启 Conversation 需要 L2 级隐私确认（PrivacyConfirm 未接线）")
		}
		if !confirm() {
			s.log("session: conversation 首次开启被拒（隐私确认未通过），仍是 OFF")
			return false, nil
		}
	}

	s.mu.Lock()
	if s.stage == StageDisposed {
		s.mu.Unlock()
		return false, errors.New("session: 会话已结束，不能开启 Conversation")
	}
	s.conversation = true
	s.confirmed = true
	if first {
		s.prompts++
	}
	s.armConvTimerLocked()
	s.mu.Unlock()

	if first {
		s.log("session: CONVERSATION-PRIVACY-CONFIRM granted（首次开启，L2 卡片已确认，已记入审计）")
	}
	s.log("session: conversation ON（idle_timeout=%s）", s.conversationIdle)
	return true, nil
}

// NoteSpeech resets the Conversation idle timer. It is the pipeline's "the user
// is talking" signal: 30 s of NO speech is what falls back to Warm, so every
// detected utterance restarts the window. Outside Conversation it is a no-op.
func (s *Scope) NoteSpeech() {
	s.mu.Lock()
	if s.stage == StageDisposed || !s.conversation {
		s.mu.Unlock()
		return
	}
	s.armConvTimerLocked()
	s.mu.Unlock()
}

// Settle begins the 3 s static window (D43 #32). The cascades do not run here:
// engines stay hot, which is what makes a wake during settling free of loads.
func (s *Scope) Settle(reason string) error {
	s.mu.Lock()
	if s.stage != StageWarm {
		s.mu.Unlock()
		return fmt.Errorf("session: Settle 只从 Warm 进入（当前 %s）", string(s.stage))
	}
	s.stopWarmTimerLocked()
	s.stage = StageSettling
	s.armSettleTimerLocked()
	s.mu.Unlock()
	s.log("session: SCOPE-SETTLING reason=%s static=%s engines_hot=%d", reason, s.settleStatic, s.Stats().EnginesHot)
	return nil
}

// EndSettle is D43 #33: a new wake during Settling cancels the settle. Like
// Resume it loads nothing - the cascade had not run, so the handles are still
// the ones the previous turn used.
func (s *Scope) EndSettle() error {
	s.mu.Lock()
	if s.stage != StageSettling {
		s.mu.Unlock()
		return fmt.Errorf("session: EndSettle 只服务 Settling（当前 %s）", string(s.stage))
	}
	s.stopSettleTimerLocked()
	s.stage = StageAwake
	s.mu.Unlock()
	s.log("session: SCOPE-SETTLE-CANCELLED（唤醒取消了 3s 静态期，零装载）")
	return nil
}

// cascade is one dispose's frozen work: the engines to close, the hooks to walk
// and the reason. It is built under the mutex (so the state transition is
// atomic) and executed outside it (so a hook can log or call back).
type cascade struct {
	reason string
	slots  map[EngineKind]*slot
	hooks  Hooks
}

// Dispose runs the D43 #31 cascade in its frozen order and is terminal:
//
//  1. models unload        (every hot engine, engineOrder)
//  2. panel hide-or-destroy
//  3. FreeOSMemory         (C11 step 4, so the D32 10 s fall-back stays reachable)
//  4. grants expiry        (the session key stops being used; ticket 49 fills the body)
//
// It is the explicit end (hotkey long-press / menu) AND the Settling expiry's
// action - one function, so the two routes cannot drift apart.
func (s *Scope) Dispose(reason string) error {
	s.mu.Lock()
	c, ok := s.beginDisposeLocked(reason)
	s.mu.Unlock()
	if !ok {
		return nil // idempotent (C11): a second end is not a second cascade
	}
	return s.runCascade(c)
}

// beginDisposeLocked is the atomic half: it stops every timer, marks the scope
// disposed and hands back the work. Called with s.mu held, so a timer fire and
// a wake cannot both believe they own the transition.
func (s *Scope) beginDisposeLocked(reason string) (*cascade, bool) {
	if s.stage == StageDisposed {
		return nil, false
	}
	s.stopWarmTimerLocked()
	s.stopConvTimerLocked()
	s.stopSettleTimerLocked()
	s.stage = StageDisposed
	s.disposals++
	s.conversation = false
	s.panelReserved = false
	c := &cascade{reason: reason, slots: s.slots, hooks: s.hooks}
	s.slots = make(map[EngineKind]*slot)
	return c, true
}

// runCascade is the ordered half, and the order IS the assertion: models first,
// then the window, then the OS memory return, then the grants. It runs with no
// lock held.
func (s *Scope) runCascade(c *cascade) error {
	var errs []error
	for _, kind := range engineOrder {
		sl := c.slots[kind]
		if sl == nil || sl.eng == nil {
			continue
		}
		if cerr := sl.eng.Close(); cerr != nil {
			errs = append(errs, fmt.Errorf("session: 卸载引擎 %s 失败: %w", string(kind), cerr))
			s.log("session: SCOPE-DISPOSE engine-close-failed kind=%s err=%v", string(kind), cerr)
		}
	}
	if c.hooks.PanelHide != nil {
		c.hooks.PanelHide()
	}
	if c.hooks.FreeOSMemory != nil {
		c.hooks.FreeOSMemory()
	}
	if c.hooks.GrantsExpire != nil {
		c.hooks.GrantsExpire()
	}
	if c.hooks.MicClose != nil {
		// Defence in depth, deliberately AFTER the four frozen steps: an explicit
		// end can arrive from Speaking with the mic open, and a destroyed session
		// must never leave the device armed. It is not part of D43 #31's order.
		c.hooks.MicClose()
	}
	s.log("session: SCOPE-DISPOSED reason=%s cascade=models-unload,panel-hide,free-os-memory,grants-expire engines=%d",
		c.reason, len(c.slots))
	return errors.Join(errs...)
}

// Close is the C11 spelling callers already use elsewhere in the tree. It is
// Dispose with the caller's reason, and idempotent the same way.
func (s *Scope) Close(reason string) error { return s.Dispose(reason) }

// ---- timers ---------------------------------------------------------------

func (s *Scope) armWarmTimerLocked() {
	s.stopWarmTimerLocked()
	s.warmGen++
	gen := s.warmGen
	s.warmTimer = s.clock.AfterFunc(s.warmTimeout, func() { s.onWarmIdle(gen) })
}

func (s *Scope) stopWarmTimerLocked() {
	if s.warmTimer != nil {
		s.warmTimer.Stop()
		s.warmTimer = nil
	}
	s.warmGen++ // a fire that already passed Stop is stale now
}

func (s *Scope) armConvTimerLocked() {
	s.stopConvTimerLocked()
	s.convGen++
	gen := s.convGen
	s.convTimer = s.clock.AfterFunc(s.conversationIdle, func() { s.onConversationIdle(gen) })
}

func (s *Scope) stopConvTimerLocked() {
	if s.convTimer != nil {
		s.convTimer.Stop()
		s.convTimer = nil
	}
	s.convGen++
}

func (s *Scope) armSettleTimerLocked() {
	s.stopSettleTimerLocked()
	s.settleGen++
	gen := s.settleGen
	s.settleTimer = s.clock.AfterFunc(s.settleStatic, func() { s.onSettleExpired(gen) })
}

func (s *Scope) stopSettleTimerLocked() {
	if s.settleTimer != nil {
		s.settleTimer.Stop()
		s.settleTimer = nil
	}
	s.settleGen++
}

// onWarmIdle is D43 #31's trigger (warm.idle-90s): Warm -> Settling. The
// cascade is NOT run here; the 3 s static window is what follows.
func (s *Scope) onWarmIdle(gen uint64) {
	defer s.recoverTimer("warm-idle")
	s.mu.Lock()
	if gen != s.warmGen || s.stage != StageWarm {
		s.mu.Unlock()
		return
	}
	s.warmTimer = nil
	s.mu.Unlock()
	_ = s.Settle("warm.idle")
}

// onConversationIdle is D43 #12's Conversation row (30 s no speech): close the
// mic FIRST, then land in Warm. EnterWarm owns that order, so this is one call.
func (s *Scope) onConversationIdle(gen uint64) {
	defer s.recoverTimer("conversation-idle")
	s.mu.Lock()
	if gen != s.convGen || !s.conversation || s.stage != StageAwake {
		s.mu.Unlock()
		return
	}
	s.convTimer = nil
	s.mu.Unlock()
	_ = s.EnterWarm("conversation.idle-30s")
}

// onSettleExpired is D43 #32's trigger (settling.expired-3s): the cascade. The
// transition is taken atomically under the mutex, so a wake racing this fire
// either cancels the settle first or loses to it - never both.
func (s *Scope) onSettleExpired(gen uint64) {
	defer s.recoverTimer("settle-expired")
	s.mu.Lock()
	if gen != s.settleGen || s.stage != StageSettling {
		s.mu.Unlock()
		return
	}
	s.settleTimer = nil
	c, ok := s.beginDisposeLocked("settling.expired")
	s.mu.Unlock()
	if !ok {
		return
	}
	_ = s.runCascade(c)
}

// recoverTimer keeps a panicking hook from taking the process down through a
// runtime timer goroutine, which cannot be recovered by anyone above it. The
// panic is logged, not swallowed.
func (s *Scope) recoverTimer(which string) {
	if r := recover(); r != nil {
		s.log("session: SCOPE-TIMER-PANIC which=%s panic=%v", which, r)
	}
}

// ---- helpers --------------------------------------------------------------

func knownEngine(k EngineKind) bool {
	switch k {
	case EngineASR, EngineTTS, EngineVAD, EnginePunct, EngineKWS:
		return true
	}
	return false
}

func (s *Scope) log(format string, args ...any) {
	if s.hooks.Logf == nil {
		return
	}
	s.hooks.Logf(format, args...)
}
