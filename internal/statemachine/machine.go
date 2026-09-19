package statemachine

import (
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// IllegalTransitionError is returned (and logged loudly) for any (state,
// event) pair that has no row in the D43 table: undefined = stop (D22 gate
// 3). The machine does NOT change state on rejection.
type IllegalTransitionError struct {
	From  State
	Event Event
}

func (e *IllegalTransitionError) Error() string {
	return fmt.Sprintf("statemachine: illegal transition %s + %s (no D43 row; undefined = stop)",
		e.From, e.Event)
}

// Effect is one fired side-effect hook. Effects are EVENTS: the machine never
// executes them itself - consumers (session, speech, panel, tickets 08+) wire
// a Sink. The zero sink is an explicit no-op.
type Effect struct {
	D43  int    // originating table row
	Name string // stable hook name from the table
}

// Sink receives fired side effects. Must not block and must not re-enter
// Dispatch synchronously (it runs on the dispatching goroutine; re-entry
// deadlocks on the machine mutex - post to the caller's own queue instead).
type Sink func(Effect)

// Options adjusts New (tests, embedders).
type Options struct {
	Initial  State                      // default StateFirstRun
	Sink     Sink                       // nil = explicit no-op
	Timeouts map[timeoutKey]TimeoutSpec // nil = DefaultTimeouts(); tests override
}

// Machine is the table-driven D43 state machine. It is safe for concurrent
// use; Dispatch serializes on a mutex. It owns at most ONE runtime timer at a
// time - the per-state timeout - and holds ZERO timers in states without a
// timeout row (notably Sleeping: the zero-timer discipline, ticket 07).
type Machine struct {
	mu       sync.Mutex
	state    State
	phase    Phase // Listening timeout context
	timeouts map[timeoutKey]TimeoutSpec
	sink     Sink
	timer    *time.Timer // nil = none armed
	gen      uint64      // timer generation; stale callbacks no-op
	closed   bool
}

// New builds a machine at opts.Initial (default FirstRun - boot starts in
// onboarding) with the frozen timeout table.
func New(opts Options) *Machine {
	if opts.Initial == "" {
		opts.Initial = StateFirstRun
	}
	if opts.Sink == nil {
		// Explicit no-op side-effect sink (ticket 07: hooks fire as events,
		// consumed by later tickets).
		opts.Sink = func(Effect) {}
	}
	if opts.Timeouts == nil {
		opts.Timeouts = DefaultTimeouts()
	}
	m := &Machine{state: opts.Initial, timeouts: opts.Timeouts, sink: opts.Sink}
	m.rearmLocked() // the initial state may already carry a timeout row
	return m
}

// State returns the current state.
func (m *Machine) State() State {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.state
}

// Phase returns the current Listening timeout context (PhaseNone outside
// Listening).
func (m *Machine) Phase() Phase {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.phase
}

// TimersAlive reports how many machine timers are armed (0 or 1). The ticket
// 07 acceptance asserts TimersAlive() == 0 in Sleeping.
func (m *Machine) TimersAlive() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.timer == nil {
		return 0
	}
	return 1
}

// Dispatch feeds one event. It returns the destination state. A rejected
// (illegal) transition returns IllegalTransitionError, logs loudly and leaves
// the state unchanged. Guard failures of ALL candidate rows for a legal
// (state,event) pair are also rejected: an event whose row exists but whose
// guard does not hold has no defined outcome.
func (m *Machine) Dispatch(ev Event, facts *Facts) (State, error) {
	if facts == nil {
		facts = &Facts{}
	}

	m.mu.Lock()

	if m.closed {
		m.mu.Unlock()
		return m.state, fmt.Errorf("statemachine: dispatch %s on closed machine", ev)
	}

	from := m.state
	var hit *Row
	for i := range Table {
		r := &Table[i]
		if r.Event != ev {
			continue
		}
		if r.From != from && r.From != AnyState {
			continue
		}
		if r.GuardFn != nil && !r.GuardFn(facts) {
			continue
		}
		hit = r
		break
	}
	if hit == nil {
		m.mu.Unlock()
		err := &IllegalTransitionError{From: from, Event: ev}
		slog.Error("state machine rejected transition", "from", string(from), "event", string(ev),
			"err", err.Error())
		return from, err
	}

	m.state = hit.To
	to := hit.To
	if hit.To == StateListening {
		if hit.Listen != PhaseNone {
			m.phase = hit.Listen
		} else if m.phase == PhaseNone {
			m.phase = PhaseFirstRound // defensive: Listening with no pinned context is first-round
		}
	} else {
		m.phase = PhaseNone
	}
	d43 := hit.D43
	effects := hit.SideEffects

	m.rearmLocked()

	m.mu.Unlock()

	for _, name := range effects {
		m.sink(Effect{D43: d43, Name: name})
	}
	return to, nil
}

// Close disarms any timer. The machine is unusable afterwards.
func (m *Machine) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	m.stopTimerLocked()
}

// rearmLocked stops any running timer and arms the timeout row of the new
// state (if any). Called with m.mu held, on state changes.
func (m *Machine) rearmLocked() {
	m.stopTimerLocked()
	key := timeoutKey{state: m.state, phase: m.phase}
	spec, ok := m.timeouts[key]
	if !ok {
		return // no timeout row: ZERO timers in this state (Sleeping discipline)
	}
	m.gen++
	gen := m.gen
	m.timer = time.AfterFunc(spec.After, func() { m.onTimeout(gen, spec.Event) })
}

// onTimeout runs on the runtime timer goroutine: the armed generation is
// checked under the lock (a state change since arming makes the fire stale),
// then the timeout event is dispatched like any other event.
func (m *Machine) onTimeout(gen uint64, ev Event) {
	m.mu.Lock()
	if m.closed || gen != m.gen {
		m.mu.Unlock()
		return
	}
	m.timer = nil // fired; rearmLocked (inside Dispatch) arms the next one
	m.mu.Unlock()

	_, _ = m.Dispatch(ev, nil)
}

func (m *Machine) stopTimerLocked() {
	if m.timer != nil {
		m.timer.Stop()
		m.timer = nil
	}
}
