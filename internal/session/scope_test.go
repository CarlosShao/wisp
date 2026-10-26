package session

// Ticket 28's portable readings: the keep-alive policy itself (warm/hot
// engines, the 90 s and 30 s timers, the Conversation privacy gate).
//
// Everything here runs on a fake Clock, on purpose: the ticket's timer rows are
// 90 s / 30 s / 3 s, and a test that slept through them would be a test nobody
// runs. The clock is injected at the same seam production uses, so what is
// exercised is the real transition code - only the "how long" part is virtual.
// The REAL engines (silero VAD + streaming Paraformer) are measured separately
// in scope_engines_windows_test.go, which needs the staged DLLs.

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"
)

// ---- fake clock -----------------------------------------------------------

// fakeClock is a deterministic timer source. Advance fires everything due in
// (due time, insertion order) and runs callbacks synchronously, so a test can
// assert the exact state after "90 seconds" without sleeping.
type fakeClock struct {
	mu     sync.Mutex
	now    time.Duration
	nextID int
	timers map[int]*fakeTimer
	armed  []time.Duration // every duration ever asked for, for value assertions
}

type fakeTimer struct {
	clock *fakeClock
	id    int
	due   time.Duration
	fn    func()
	dead  bool
}

func (t *fakeTimer) Stop() bool {
	c := t.clock
	c.mu.Lock()
	defer c.mu.Unlock()
	if t.dead {
		return false
	}
	t.dead = true
	delete(c.timers, t.id)
	return true
}

func newFakeClock() *fakeClock {
	return &fakeClock{timers: make(map[int]*fakeTimer)}
}

func (c *fakeClock) AfterFunc(d time.Duration, fn func()) Timer {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.nextID++
	t := &fakeTimer{clock: c, id: c.nextID, due: c.now + d, fn: fn}
	c.timers[t.id] = t
	c.armed = append(c.armed, d)
	return t
}

// Advance moves the virtual clock and fires every timer that becomes due.
func (c *fakeClock) Advance(d time.Duration) {
	c.mu.Lock()
	target := c.now + d
	for {
		var due *fakeTimer
		for _, t := range c.timers {
			if t.dead || t.due > target {
				continue
			}
			if due == nil || t.due < due.due || (t.due == due.due && t.id < due.id) {
				due = t
			}
		}
		if due == nil {
			break
		}
		c.now = due.due
		due.dead = true
		delete(c.timers, due.id)
		fn := due.fn
		c.mu.Unlock()
		fn()
		c.mu.Lock()
	}
	c.now = target
	c.mu.Unlock()
}

// Pending reports armed timers (a disposed scope must leave zero).
func (c *fakeClock) Pending() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	n := 0
	for _, t := range c.timers {
		if !t.dead {
			n++
		}
	}
	return n
}

// Armed returns the durations the scope asked for, sorted ascending.
func (c *fakeClock) Armed() []time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := append([]time.Duration(nil), c.armed...)
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

// ---- fake engines ---------------------------------------------------------

type fakeEngine struct {
	kind    EngineKind
	closes  int
	log     *eventLog
	onClose func()
}

func (e *fakeEngine) Close() error {
	e.closes++
	if e.onClose != nil {
		e.onClose()
	}
	if e.log != nil {
		e.log.add("engine.close:" + string(e.kind))
	}
	return nil
}

// engineStub is the Loader half: it counts loads and closes per kind, keeps the
// latest handle it made (so a test can prove dispose closed it) and can be told
// to fail a kind, which is how the "TTS is a reserved slot" shape is exercised.
type engineStub struct {
	mu     sync.Mutex
	loads  map[EngineKind]int
	closes map[EngineKind]int
	made   map[EngineKind]*fakeEngine
	failOn map[EngineKind]error
	log    *eventLog
}

func newEngineStub(log *eventLog) *engineStub {
	return &engineStub{
		loads: map[EngineKind]int{}, closes: map[EngineKind]int{},
		made: map[EngineKind]*fakeEngine{}, failOn: map[EngineKind]error{}, log: log,
	}
}

func (s *engineStub) load(kind EngineKind) (Engine, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.failOn[kind]; err != nil {
		return nil, err
	}
	s.loads[kind]++
	e := &fakeEngine{kind: kind, log: s.log, onClose: func() {
		s.mu.Lock()
		s.closes[kind]++
		s.mu.Unlock()
	}}
	s.made[kind] = e
	return e, nil
}

func (s *engineStub) closesOf(kind EngineKind) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closes[kind]
}

func (s *engineStub) loadsOf(kind EngineKind) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loads[kind]
}

func (s *engineStub) totalLoads() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	for _, v := range s.loads {
		n += v
	}
	return n
}

func (s *engineStub) engine(kind EngineKind) *fakeEngine {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.made[kind]
}

// ---- event log ------------------------------------------------------------

// eventLog records the ordered seams (mic close, cascade steps, log lines) so
// an ORDER can be asserted rather than described. It is safe for concurrent use
// because hook calls are permitted from any goroutine.
type eventLog struct {
	mu     sync.Mutex
	events []string
}

func (l *eventLog) add(e string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.events = append(l.events, e)
}

func (l *eventLog) all() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]string(nil), l.events...)
}

func (l *eventLog) index(e string) int {
	for i, got := range l.all() {
		if got == e {
			return i
		}
	}
	return -1
}

// indexPrefix finds the first event that starts with p, which is how log lines
// are matched: the format string carries counters the test does not want to
// restate.
func (l *eventLog) indexPrefix(p string) int {
	for i, got := range l.all() {
		if len(got) >= len(p) && got[:len(p)] == p {
			return i
		}
	}
	return -1
}

func (l *eventLog) countPrefix(p string) int {
	n := 0
	for _, got := range l.all() {
		if len(got) >= len(p) && got[:len(p)] == p {
			n++
		}
	}
	return n
}

// scopeBuilder keeps the three AC tests short: each one names only the seams it
// cares about and gets the rest recorded in one shared log.
type scopeBuilder struct {
	stub    *engineStub
	clock   *fakeClock
	log     *eventLog
	hooks   Hooks
	opts    ScopeOptions
	confirm func() bool
}

func newScopeBuilder() *scopeBuilder {
	log := &eventLog{}
	return &scopeBuilder{stub: newEngineStub(log), clock: newFakeClock(), log: log}
}

func (b *scopeBuilder) build(t *testing.T) (*Scope, *fakeClock, *engineStub, *eventLog) {
	t.Helper()
	b.hooks = Hooks{
		MicOpen:        func() { b.log.add("mic.open") },
		MicClose:       func() { b.log.add("mic.close") },
		PanelHide:      func() { b.log.add("panel.hide") },
		FreeOSMemory:   func() { b.log.add("free-os-memory") },
		GrantsExpire:   func() { b.log.add("grants.expire") },
		PrivacyConfirm: b.confirm,
		Logf:           func(f string, a ...any) { b.log.add("log:" + fmt.Sprintf(f, a...)) },
	}
	b.opts.Loader = b.stub.load
	b.opts.Clock = b.clock
	b.opts.Hooks = b.hooks
	s, err := NewScope(b.opts)
	if err != nil {
		t.Fatalf("NewScope: %v", err)
	}
	return s, b.clock, b.stub, b.log
}

// ---- AC1: cold vs warm, zero loads ----------------------------------------

// TestTicket28WarmWakeLoadsNothing is AC1's portable half: the state is Warm,
// the engines are still hot, and the wake that follows touches the Loader ZERO
// times. The counter AC1 names is Stats().LoadsAfterWarm, and it is the thing
// the mutation in the evidence file makes go red.
func TestTicket28WarmWakeLoadsNothing(t *testing.T) {
	b := newScopeBuilder()
	s, clock, stub, _ := b.build(t)
	defer func() { _ = s.Dispose("test-end") }()

	for _, k := range []EngineKind{EngineVAD, EngineASR} {
		if _, err := s.Acquire(k); err != nil {
			t.Fatalf("cold acquire %s: %v", k, err)
		}
	}
	cold := s.Stats()
	if cold.ColdLoads != 2 || cold.LoadsAfterWarm != 0 {
		t.Fatalf("cold stats = %+v, want ColdLoads=2 LoadsAfterWarm=0", cold)
	}
	if err := s.Release(EngineVAD); err != nil {
		t.Fatalf("release vad: %v", err)
	}
	if err := s.Release(EngineASR); err != nil {
		t.Fatalf("release asr: %v", err)
	}

	if err := s.EnterWarm("vad.stop"); err != nil {
		t.Fatalf("enter warm: %v", err)
	}
	warm := s.Stats()
	if warm.Stage != StageWarm {
		t.Fatalf("stage = %s, want Warm", warm.Stage)
	}
	if warm.EnginesHot != 2 || warm.Refs != 0 {
		t.Fatalf("warm stats = %+v, want EnginesHot=2 Refs=0 (hot but unborrowed)", warm)
	}

	// Five wake -> warm rounds: this is B4's claim, repeated, not asserted once.
	for i := 0; i < 5; i++ {
		if err := s.Resume(); err != nil {
			t.Fatalf("resume %d: %v", i, err)
		}
		for _, k := range []EngineKind{EngineVAD, EngineASR} {
			if _, err := s.Acquire(k); err != nil {
				t.Fatalf("warm acquire %s (round %d): %v", k, i, err)
			}
			if err := s.Release(k); err != nil {
				t.Fatalf("warm release %s (round %d): %v", k, i, err)
			}
		}
		if err := s.EnterWarm("turn.end"); err != nil {
			t.Fatalf("re-enter warm %d: %v", i, err)
		}
	}
	st := s.Stats()
	if st.LoadsAfterWarm != 0 {
		t.Fatalf("LoadsAfterWarm = %d after 5 warm wakes, want 0 (B4: the warm wake loads nothing)", st.LoadsAfterWarm)
	}
	if stub.totalLoads() != 2 {
		t.Fatalf("loader was called %d times, want 2 (one per engine, ever)", stub.totalLoads())
	}
	if st.WarmResumes != 5 {
		t.Fatalf("WarmResumes = %d, want 5", st.WarmResumes)
	}
	// The warm timer was re-armed per warm entry and is still armed; a dispose
	// must leave none (AC6's timer half, asserted here because it is cheap).
	if clock.Pending() == 0 {
		t.Fatal("no warm timer armed while Warm: the 90s row would never fire")
	}
}

// TestTicket28WarmTimerIsNinetySeconds pins the VALUE, not just the behaviour:
// a Warm that expires at 89.9 s would still "work" and would still be wrong.
func TestTicket28WarmTimerIsNinetySeconds(t *testing.T) {
	b := newScopeBuilder()
	s, clock, _, _ := b.build(t)
	defer func() { _ = s.Dispose("test-end") }()
	if _, err := s.Acquire(EngineASR); err != nil {
		t.Fatalf("acquire: %v", err)
	}
	_ = s.Release(EngineASR)
	if err := s.EnterWarm("turn.end"); err != nil {
		t.Fatalf("enter warm: %v", err)
	}
	armed := clock.Armed()
	if len(armed) != 1 || armed[0] != WarmTimeoutDefault {
		t.Fatalf("armed = %v, want exactly [%s]", armed, WarmTimeoutDefault)
	}
	if WarmTimeoutDefault != 90*time.Second {
		t.Fatalf("WarmTimeoutDefault = %s, want 90s (ticket 28 key constraint)", WarmTimeoutDefault)
	}
}

// ---- AC2: timers ----------------------------------------------------------

// TestTicket28WarmNinetySecondsSettlingDispose walks D43 #31 and #32 on the
// virtual clock: 90 s -> Settling (engines STILL hot), 3 s -> disposed with the
// cascade run. The intermediate assertion is the load-bearing one: an
// implementation that disposes at the 90 s mark passes an end-state check and
// breaks #33, so the settle window must be observable.
func TestTicket28WarmNinetySecondsSettlingDispose(t *testing.T) {
	b := newScopeBuilder()
	s, clock, stub, log := b.build(t)
	if _, err := s.Acquire(EngineASR); err != nil {
		t.Fatalf("acquire asr: %v", err)
	}
	if _, err := s.Acquire(EngineVAD); err != nil {
		t.Fatalf("acquire vad: %v", err)
	}
	_ = s.Release(EngineASR)
	_ = s.Release(EngineVAD)
	if err := s.EnterWarm("turn.end"); err != nil {
		t.Fatalf("enter warm: %v", err)
	}

	clock.Advance(89 * time.Second)
	if got := s.Stage(); got != StageWarm {
		t.Fatalf("stage at 89s = %s, want Warm", got)
	}
	clock.Advance(1 * time.Second)
	if got := s.Stage(); got != StageSettling {
		t.Fatalf("stage at 90s = %s, want Settling (D43 #31)", got)
	}
	if s.Stats().EnginesHot != 2 {
		t.Fatalf("Settling unloaded engines: EnginesHot = %d, want 2 (the 3s window is still warm)", s.Stats().EnginesHot)
	}
	if stub.engine(EngineASR).closes != 0 {
		t.Fatal("engine closed during Settling: the cascade must run at #32, not #31")
	}

	clock.Advance(2999 * time.Millisecond)
	if got := s.Stage(); got != StageSettling {
		t.Fatalf("stage at 92.999s = %s, want still Settling", got)
	}
	clock.Advance(1 * time.Millisecond)
	if got := s.Stage(); got != StageDisposed {
		t.Fatalf("stage at 93s = %s, want Disposed (D43 #32)", got)
	}
	if stub.engine(EngineASR).closes != 1 || stub.engine(EngineVAD).closes != 1 {
		t.Fatalf("engine closes = %d/%d, want 1/1 (models unload is the cascade's first step)",
			stub.engine(EngineASR).closes, stub.engine(EngineVAD).closes)
	}
	if s.Stats().EnginesHot != 0 || s.Stats().Refs != 0 {
		t.Fatalf("post-dispose stats = %+v, want zero hot engines and zero refs", s.Stats())
	}
	if clock.Pending() != 0 {
		t.Fatalf("%d timers still armed after dispose, want 0", clock.Pending())
	}
	assertCascadeOrder(t, log)
}

// TestTicket28ConversationThirtySecondsClosesMicFirst is AC2's second row: 30 s
// of no speech falls back to Warm, and the ORDER (mic closed BEFORE the warm
// transition lands) is asserted twice - by the seam order in the log and by
// reading Stage() from inside the mic-close hook.
func TestTicket28ConversationThirtySecondsClosesMicFirst(t *testing.T) {
	b := newScopeBuilder()
	b.confirm = func() bool { return true }
	stageAtMicClose := Stage("")
	s, clock, _, _ := b.build(t)
	// Rebuild with a stage-reading MicClose: the builder's hooks are frozen at
	// construction, so this hook is added through the scope's own seam.
	scopeWithStageProbe(t, s, &stageAtMicClose, b.log)
	defer func() { _ = s.Dispose("test-end") }()

	if _, err := s.Acquire(EngineASR); err != nil {
		t.Fatalf("acquire: %v", err)
	}
	_ = s.Release(EngineASR)
	if _, err := s.EnableConversation(); err != nil {
		t.Fatalf("enable conversation: %v", err)
	}
	if !s.Stats().Conversation {
		t.Fatal("conversation not on after a granted enable")
	}

	clock.Advance(29999 * time.Millisecond)
	if s.Stage() != StageAwake || !s.Stats().Conversation {
		t.Fatalf("at 29.999s: stage=%s conversation=%v, want Awake/true", s.Stage(), s.Stats().Conversation)
	}
	clock.Advance(1 * time.Millisecond)
	if got := s.Stage(); got != StageWarm {
		t.Fatalf("stage at 30s = %s, want Warm (D43 #12 conversation row)", got)
	}
	if s.Stats().Conversation {
		t.Fatal("conversation still on after the 30s fallback, want OFF (the row ends the mode)")
	}
	micIdx := b.log.index("mic.close")
	warmIdx := b.log.indexPrefix("log:session: SCOPE-WARM reason=conversation.idle-30s")
	if micIdx < 0 || warmIdx < 0 {
		t.Fatalf("missing seams: mic.close=%d warm=%d log=%v", micIdx, warmIdx, b.log.all())
	}
	if micIdx > warmIdx {
		t.Fatalf("mic.close at %d came AFTER the warm transition at %d: the ticket says close the mic FIRST", micIdx, warmIdx)
	}
	if stageAtMicClose != StageAwake {
		t.Fatalf("Stage() inside MicClose = %s, want Awake (the transition must not have landed yet)", stageAtMicClose)
	}
	if s.Stats().EnginesHot != 1 {
		t.Fatalf("EnginesHot after the fallback = %d, want 1 (Warm keeps models hot)", s.Stats().EnginesHot)
	}
}

// TestTicket28ConversationIdleResetsOnSpeech proves the 30 s window is idle
// time, not session age: a turn in the middle restarts it.
func TestTicket28ConversationIdleResetsOnSpeech(t *testing.T) {
	b := newScopeBuilder()
	b.confirm = func() bool { return true }
	s, clock, _, _ := b.build(t)
	defer func() { _ = s.Dispose("test-end") }()

	if _, err := s.EnableConversation(); err != nil {
		t.Fatalf("enable: %v", err)
	}
	clock.Advance(20 * time.Second)
	s.NoteSpeech()
	clock.Advance(29 * time.Second)
	if s.Stage() != StageAwake {
		t.Fatalf("stage after 49s total but 29s since speech = %s, want Awake", s.Stage())
	}
	clock.Advance(1 * time.Second)
	if s.Stage() != StageWarm {
		t.Fatalf("stage at 30s since speech = %s, want Warm", s.Stage())
	}
}

// TestTicket28BatterySaverShapesAreThirtyAndFifteen is D42#6: the same two rows
// with the power-policy values (Warm 30 s, Conversation 15 s).
func TestTicket28BatterySaverShapesAreThirtyAndFifteen(t *testing.T) {
	b := newScopeBuilder()
	b.confirm = func() bool { return true }
	b.opts.BatterySaver = true
	s, clock, _, _ := b.build(t)
	defer func() { _ = s.Dispose("test-end") }()

	if _, err := s.EnableConversation(); err != nil {
		t.Fatalf("enable: %v", err)
	}
	clock.Advance(14999 * time.Millisecond)
	if s.Stage() != StageAwake {
		t.Fatalf("battery conversation expired early: stage=%s at 14.999s", s.Stage())
	}
	clock.Advance(1 * time.Millisecond)
	if s.Stage() != StageWarm {
		t.Fatalf("battery conversation at 15s: stage=%s, want Warm", s.Stage())
	}
	clock.Advance(29999 * time.Millisecond)
	if s.Stage() != StageWarm {
		t.Fatalf("battery warm expired early: stage=%s at 29.999s", s.Stage())
	}
	clock.Advance(1 * time.Millisecond)
	if s.Stage() != StageSettling {
		t.Fatalf("battery warm at 30s: stage=%s, want Settling", s.Stage())
	}
	if WarmTimeoutBattery != 30*time.Second || ConversationIdleBattery != 15*time.Second {
		t.Fatalf("battery constants = %s/%s, want 30s/15s (D42#6)",
			WarmTimeoutBattery, ConversationIdleBattery)
	}
}

// TestTicket28WakeDuringSettlingCancelsWithZeroLoad is D43 #33: a new wake
// during the 3 s window cancels the settle and loads nothing, because the
// cascade had not run yet.
func TestTicket28WakeDuringSettlingCancelsWithZeroLoad(t *testing.T) {
	b := newScopeBuilder()
	s, clock, stub, _ := b.build(t)
	defer func() { _ = s.Dispose("test-end") }()

	if _, err := s.Acquire(EngineASR); err != nil {
		t.Fatalf("acquire: %v", err)
	}
	_ = s.Release(EngineASR)
	if err := s.EnterWarm("turn.end"); err != nil {
		t.Fatalf("enter warm: %v", err)
	}
	clock.Advance(90 * time.Second)
	if s.Stage() != StageSettling {
		t.Fatalf("stage = %s, want Settling", s.Stage())
	}
	if err := s.EndSettle(); err != nil {
		t.Fatalf("end settle: %v", err)
	}
	if s.Stage() != StageAwake {
		t.Fatalf("stage after cancel = %s, want Awake (D43 #33)", s.Stage())
	}
	if _, err := s.Acquire(EngineASR); err != nil {
		t.Fatalf("acquire after cancel: %v", err)
	}
	if stub.totalLoads() != 1 || s.Stats().LoadsAfterWarm != 0 {
		t.Fatalf("loads = %d, loadsAfterWarm = %d; cancelling a settle must load nothing",
			stub.totalLoads(), s.Stats().LoadsAfterWarm)
	}
	if stub.engine(EngineASR).closes != 0 {
		t.Fatal("engine closed by a cancelled settle")
	}
	// The settle timer must be gone: advancing past it must not dispose.
	clock.Advance(3 * time.Second)
	if s.Stage() != StageAwake {
		t.Fatalf("stage after the cancelled settle window = %s, want Awake", s.Stage())
	}
}

// TestTicket28ExplicitEndDisposesImmediately is the other teardown route: an
// explicit end (hotkey long-press / menu) does NOT wait for 90 s + 3 s.
func TestTicket28ExplicitEndDisposesImmediately(t *testing.T) {
	b := newScopeBuilder()
	s, _, stub, log := b.build(t)
	if _, err := s.Acquire(EngineASR); err != nil {
		t.Fatalf("acquire: %v", err)
	}
	if err := s.Dispose("hotkey.long-press"); err != nil {
		t.Fatalf("dispose: %v", err)
	}
	if s.Stage() != StageDisposed {
		t.Fatalf("stage = %s, want Disposed", s.Stage())
	}
	if stub.engine(EngineASR).closes != 1 {
		t.Fatalf("engine closes = %d, want 1", stub.engine(EngineASR).closes)
	}
	assertCascadeOrder(t, log)
	// Idempotent (C11): a second end is not a second cascade.
	if err := s.Dispose("again"); err != nil {
		t.Fatalf("second dispose: %v", err)
	}
	if stub.engine(EngineASR).closes != 1 {
		t.Fatalf("second dispose closed the engine again (%d closes)", stub.engine(EngineASR).closes)
	}
	if s.Stats().Disposals != 1 {
		t.Fatalf("Disposals = %d, want 1", s.Stats().Disposals)
	}
}

// assertCascadeOrder is AC5's order assertion as one function: models unload
// (the engine close event) then panel, then the OS memory return, then the
// grants - and the mic close is deliberately last.
func assertCascadeOrder(t *testing.T, log *eventLog) {
	t.Helper()
	models := log.indexPrefix("engine.close:")
	panel := log.index("panel.hide")
	free := log.index("free-os-memory")
	grants := log.index("grants.expire")
	if models < 0 || panel < 0 || free < 0 || grants < 0 {
		t.Fatalf("cascade step missing: models=%d panel=%d free=%d grants=%d log=%v",
			models, panel, free, grants, log.all())
	}
	if !(models < panel && panel < free && free < grants) {
		t.Fatalf("cascade order = models:%d panel:%d free:%d grants:%d, "+
			"want models then panel then free-os-memory then grants (D43 #31)",
			models, panel, free, grants)
	}
}

// ---- AC4: the Conversation privacy gate -----------------------------------

// TestTicket28ConversationPrivacyConfirmOnce is AC4: the first enable costs an
// L2-grade confirm and is LOGGED; later enables do not prompt.
func TestTicket28ConversationPrivacyConfirmOnce(t *testing.T) {
	prompts := 0
	b := newScopeBuilder()
	b.confirm = func() bool { prompts++; return true }
	s, _, _, log := b.build(t)
	defer func() { _ = s.Dispose("test-end") }()

	if ok, err := s.EnableConversation(); err != nil || !ok {
		t.Fatalf("first enable = (%v, %v), want (true, nil)", ok, err)
	}
	if prompts != 1 {
		t.Fatalf("prompts = %d after the first enable, want 1", prompts)
	}
	if log.countPrefix("log:session: CONVERSATION-PRIVACY-CONFIRM granted") != 1 {
		t.Fatalf("the grant was not logged once: %v", log.all())
	}
	if s.Stats().Prompts != 1 {
		t.Fatalf("Stats.Prompts = %d, want 1", s.Stats().Prompts)
	}

	// Fall back to Warm (the mode ends) and enable again: no second card.
	if err := s.EnterWarm("conversation.idle"); err != nil {
		t.Fatalf("enter warm: %v", err)
	}
	if err := s.Resume(); err != nil {
		t.Fatalf("resume: %v", err)
	}
	if ok, err := s.EnableConversation(); err != nil || !ok {
		t.Fatalf("second enable = (%v, %v), want (true, nil)", ok, err)
	}
	if prompts != 1 {
		t.Fatalf("prompts = %d after the second enable, want 1 (subsequent enables must not prompt)", prompts)
	}
	if log.countPrefix("log:session: CONVERSATION-PRIVACY-CONFIRM granted") != 1 {
		t.Fatalf("the confirm line appeared %d times, want 1", log.countPrefix("log:session: CONVERSATION-PRIVACY-CONFIRM granted"))
	}
}

// TestTicket28ConversationDeniedOrUnwiredStaysOff: a denied card and a missing
// card surface both leave the mode OFF. Fail-closed, and the denial is logged.
func TestTicket28ConversationDeniedOrUnwiredStaysOff(t *testing.T) {
	b := newScopeBuilder()
	b.confirm = func() bool { return false }
	s, _, _, log := b.build(t)
	if ok, err := s.EnableConversation(); err != nil || ok {
		t.Fatalf("denied enable = (%v, %v), want (false, nil)", ok, err)
	}
	if s.Stats().Conversation {
		t.Fatal("conversation is on after a denied card")
	}
	if log.countPrefix("log:session: conversation 首次开启被拒") != 1 {
		t.Fatalf("the denial was not logged: %v", log.all())
	}
	_ = s.Dispose("test-end")

	b2 := newScopeBuilder() // no confirm hook at all
	s2, _, _, _ := b2.build(t)
	defer func() { _ = s2.Dispose("test-end") }()
	if ok, err := s2.EnableConversation(); err == nil || ok {
		t.Fatalf("unwired enable = (%v, %v), want an error (a mode that keeps a mic open must not enable without a card)", ok, err)
	}
	if s2.Stats().Conversation {
		t.Fatal("conversation is on with no privacy confirm surface")
	}
}

// ---- reservation / fail-closed surface ------------------------------------

// TestTicket28BorrowFailsClosed covers the three ways a borrow must not invent
// an engine: an unknown kind, a loader refusal (TTS is a reserved slot until
// ticket 26), and a borrow after dispose.
func TestTicket28BorrowFailsClosed(t *testing.T) {
	b := newScopeBuilder()
	s, _, stub, _ := b.build(t)

	if _, err := s.Acquire("nosuch"); err == nil {
		t.Fatal("Acquire on an unknown kind succeeded")
	}
	stub.failOn[EngineTTS] = errors.New("ticket 26 is not landed")
	if _, err := s.Acquire(EngineTTS); err == nil {
		t.Fatal("Acquire on a refusing loader succeeded")
	}
	if s.Stats().EnginesHot != 0 {
		t.Fatalf("a failed load created %d slots, want 0", s.Stats().EnginesHot)
	}
	if err := s.Release(EngineTTS); err == nil {
		t.Fatal("Release of a never-borrowed engine succeeded")
	}
	if err := s.Dispose("test-end"); err != nil {
		t.Fatalf("dispose: %v", err)
	}
	if _, err := s.Acquire(EngineASR); err == nil {
		t.Fatal("Acquire after dispose succeeded")
	}
	if err := s.Resume(); err == nil {
		t.Fatal("Resume after dispose succeeded")
	}
	if err := s.EnterWarm("late"); err == nil {
		t.Fatal("EnterWarm after dispose succeeded")
	}
}

// scopeWithStageProbe wires a MicClose that reads the scope's Stage(). The
// builder's hooks are fixed at construction, so this is done by replacing the
// hook through a mutex-guarded wrapper the scope already holds - the same seam
// production uses, not a test-only private door.
func scopeWithStageProbe(t *testing.T, s *Scope, seen *Stage, log *eventLog) {
	t.Helper()
	s.mu.Lock()
	prev := s.hooks.MicClose
	s.hooks.MicClose = func() {
		if prev != nil {
			prev()
		}
		*seen = s.Stage()
		log.add("mic.close.probe")
	}
	s.mu.Unlock()
}

// TestTicket28WarmWindowFiresOncePerWarm is the hook tickets 29/10 run their
// async work on: it fires once per warm entry, and it fires AFTER the stage has
// landed (so the work can read Stage() == Warm and know the window is open).
func TestTicket28WarmWindowFiresOncePerWarm(t *testing.T) {
	log := &eventLog{}
	stub := newEngineStub(log)
	clock := newFakeClock()
	stageSeen := Stage("")
	var s *Scope
	s, err := NewScope(ScopeOptions{Loader: stub.load, Clock: clock, Hooks: Hooks{
		MicClose: func() {},
		WarmWindow: func() {
			stageSeen = s.Stage()
			log.add("warm.window")
		},
		Logf: func(string, ...any) {},
	}})
	if err != nil {
		t.Fatalf("NewScope: %v", err)
	}
	defer func() { _ = s.Dispose("test-end") }()

	for i := 0; i < 3; i++ {
		if err := s.EnterWarm("turn.end"); err != nil {
			t.Fatalf("enter warm %d: %v", i, err)
		}
		if err := s.Resume(); err != nil {
			t.Fatalf("resume %d: %v", i, err)
		}
	}
	if got := log.countPrefix("warm.window"); got != 3 {
		t.Fatalf("WarmWindow fired %d times over 3 warm entries, want 3", got)
	}
	if stageSeen != StageWarm {
		t.Fatalf("Stage() inside WarmWindow = %s, want Warm (the window must be open before the work starts)", stageSeen)
	}
}

// TestTicket28PanelReservationDiesWithTheSession is the "panel keep-alive
// reservation" half: the session owns the window's lifetime, so the reservation
// cannot survive a dispose, and a disposed session cannot take a new one.
func TestTicket28PanelReservationDiesWithTheSession(t *testing.T) {
	b := newScopeBuilder()
	s, _, _, log := b.build(t)
	if err := s.ReservePanel(); err != nil {
		t.Fatalf("reserve panel: %v", err)
	}
	if !s.PanelReserved() || !s.Stats().PanelReserved {
		t.Fatal("panel reservation did not stick")
	}
	if err := s.EnterWarm("turn.end"); err != nil {
		t.Fatalf("enter warm: %v", err)
	}
	if !s.PanelReserved() {
		t.Fatal("Warm dropped the panel reservation: the window stays alive through Warm (that is its keep-alive)")
	}
	if err := s.Dispose("menu.end"); err != nil {
		t.Fatalf("dispose: %v", err)
	}
	if s.PanelReserved() {
		t.Fatal("the panel reservation outlived the session")
	}
	if log.index("panel.hide") < 0 {
		t.Fatalf("the cascade did not hide the panel: %v", log.all())
	}
	if err := s.ReservePanel(); err == nil {
		t.Fatal("a disposed session accepted a new panel reservation")
	}
}
