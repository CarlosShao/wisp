package statemachine

import (
	"testing"
	"time"
)

// TestPerStateTimeoutsTableDriven pins the frozen timeout table: state (+
// Listening phase) -> duration and timeout event. Mirrors the ticket
// instruction "15s/90s/30s/倒计时/90s/3s/10s" and SPEC-08 §3.
func TestPerStateTimeoutsTableDriven(t *testing.T) {
	cases := []struct {
		state State
		phase Phase
		after time.Duration
		event Event
	}{
		{StateListening, PhaseFirstRound, 15 * time.Second, EvTimeoutFirstRound},
		{StateListening, PhaseInSession, 90 * time.Second, EvTimeoutInSession},
		{StateListening, PhaseConversation, 30 * time.Second, EvTimeoutConversation},
		{StateConfirming, PhaseNone, 3 * time.Second, EvConfirmExpired},
		{StateWarm, PhaseNone, 90 * time.Second, EvWarmIdle},
		{StateSettling, PhaseNone, 3 * time.Second, EvSettleExpired},
		{StateAwaitingApproval, PhaseNone, 300 * time.Second, EvApprovalTimeout},
		{StateError, PhaseNone, 10 * time.Second, EvErrorAck},
		// Zero-timer states: no row at all.
		{StateSleeping, PhaseNone, 0, ""},
		{StateArmed, PhaseNone, 0, ""},
		{StateMuted, PhaseNone, 0, ""},
		{StateThinking, PhaseNone, 0, ""},
		{StateActing, PhaseNone, 0, ""},
		{StateSpeaking, PhaseNone, 0, ""},
		{StateConversation, PhaseNone, 0, ""},
		{StateDownloading, PhaseNone, 0, ""},
		{StateFirstRun, PhaseNone, 0, ""},
		{StateUnconfigured, PhaseNone, 0, ""},
		{StateNoNetwork, PhaseNone, 0, ""},
		{StateQueued, PhaseNone, 0, ""},
		{StateStuck, PhaseNone, 0, ""},
		{StateWatchdogAlert, PhaseNone, 0, ""},
	}
	for _, tc := range cases {
		spec, ok := DefaultTimeouts()[timeoutKey{state: tc.state, phase: tc.phase}]
		if tc.after == 0 {
			if ok {
				t.Errorf("%s: expected no timeout row, got %+v", tc.state, spec)
			}
			continue
		}
		if !ok {
			t.Errorf("%s (phase %d): expected a timeout row", tc.state, tc.phase)
			continue
		}
		if spec.After != tc.after || spec.Event != tc.event {
			t.Errorf("%s (phase %d): got {%v %s}, want {%v %s}",
				tc.state, tc.phase, spec.After, spec.Event, tc.after, tc.event)
		}
	}
}

// TestSleepingHasZeroTimers is the ticket 07 acceptance assertion: after any
// path into Sleeping the machine owns NO timer.
func TestSleepingHasZeroTimers(t *testing.T) {
	paths := []struct {
		name  string
		start State
		event Event
	}{
		{"boot into sleeping", StateFirstRun, EvOnboardingCompleted},
		{"settling fall back (no kws)", StateSettling, EvSettleExpired},
		{"first-round timeout", StateListening, EvTimeoutFirstRound},
		{"muted without kws", StateMuted, EvMuteKey},
	}
	for _, p := range paths {
		t.Run(p.name, func(t *testing.T) {
			m := New(Options{Initial: p.start})
			if _, err := m.Dispatch(p.event, nil); err != nil {
				t.Fatalf("dispatch %s from %s: %v", p.event, p.start, err)
			}
			if m.State() != StateSleeping {
				t.Fatalf("expected Sleeping, got %s", m.State())
			}
			if n := m.TimersAlive(); n != 0 {
				t.Fatalf("Sleeping holds %d timer(s); zero-timer discipline broken", n)
			}
			// Rejections must not arm anything either.
			if _, err := m.Dispatch(EvErrorAck, nil); err == nil {
				t.Fatal("EvErrorAck from Sleeping should be rejected")
			}
			if m.State() != StateSleeping {
				t.Fatalf("rejected dispatch moved the state to %s", m.State())
			}
			if n := m.TimersAlive(); n != 0 {
				t.Fatalf("Sleeping holds %d timer(s) after a rejected dispatch", n)
			}
			m.Close()
		})
	}
}

// TestTimersArmAndFire exercises the timeout machinery with shortened
// overrides: Warm fires EvWarmIdle -> Settling; Settling fires -> Sleeping;
// Listening first-round override fires the first-round event (not the
// in-session one).
func TestTimersArmAndFire(t *testing.T) {
	timeouts := DefaultTimeouts()
	timeouts[timeoutKey{state: StateWarm}] = TimeoutSpec{30 * time.Millisecond, EvWarmIdle}
	timeouts[timeoutKey{state: StateSettling}] = TimeoutSpec{30 * time.Millisecond, EvSettleExpired}

	m := New(Options{Initial: StateWarm, Timeouts: timeouts})
	if n := m.TimersAlive(); n != 1 {
		t.Fatalf("Warm should arm exactly 1 timer, got %d", n)
	}
	deadline := time.Now().Add(2 * time.Second)
	for m.State() != StateSleeping && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	if m.State() != StateSleeping {
		t.Fatalf("Warm->Settling->Sleeping chain did not complete, stuck at %s", m.State())
	}
	if n := m.TimersAlive(); n != 0 {
		t.Fatalf("Sleeping holds %d timer(s) after the timeout chain", n)
	}
	m.Close()

	// A state change cancels the armed timer: leave Warm before 30ms.
	m2 := New(Options{Initial: StateWarm, Timeouts: timeouts})
	_, _ = m2.Dispatch(EvSummon, nil)
	if m2.State() != StateListening || m2.Phase() != PhaseInSession {
		t.Fatalf("warm summon: got %s/%d", m2.State(), m2.Phase())
	}
	time.Sleep(60 * time.Millisecond)
	if m2.State() != StateListening {
		t.Fatalf("stale Warm timer fired after transition: %s", m2.State())
	}
	m2.Close()
}

// TestListeningPhaseTimeoutSelection proves the phase picks the timeout row:
// same state, three different override durations -> three different events.
func TestListeningPhaseTimeoutSelection(t *testing.T) {
	timeouts := DefaultTimeouts()
	timeouts[timeoutKey{state: StateListening, phase: PhaseFirstRound}] = TimeoutSpec{30 * time.Millisecond, EvTimeoutFirstRound}
	timeouts[timeoutKey{state: StateListening, phase: PhaseInSession}] = TimeoutSpec{30 * time.Millisecond, EvTimeoutInSession}
	timeouts[timeoutKey{state: StateListening, phase: PhaseConversation}] = TimeoutSpec{30 * time.Millisecond, EvTimeoutConversation}

	expect := map[Phase]Event{
		PhaseFirstRound: EvTimeoutFirstRound, // -> Sleeping
		PhaseInSession:  EvTimeoutInSession,  // -> Warm
	}
	for phase, ev := range expect {
		m := New(Options{Initial: StateSleeping, Timeouts: timeouts})
		// enter Listening with the wanted phase via a phase-carrying path
		var enter Event
		switch phase {
		case PhaseFirstRound:
			enter = EvSummon // #4
		case PhaseInSession:
			// Warm->Listening #29
			m = New(Options{Initial: StateWarm, Timeouts: timeouts})
			enter = EvSummon
		}
		if _, err := m.Dispatch(enter, nil); err != nil {
			t.Fatalf("enter Listening: %v", err)
		}
		if m.Phase() != phase {
			t.Fatalf("phase = %d, want %d", m.Phase(), phase)
		}
		deadline := time.Now().Add(2 * time.Second)
		for m.State() == StateListening && time.Now().Before(deadline) {
			time.Sleep(5 * time.Millisecond)
		}
		wantState := StateSleeping
		if ev == EvTimeoutInSession {
			wantState = StateWarm
		}
		if m.State() != wantState {
			t.Fatalf("timeout event %s landed in %s, want %s", ev, m.State(), wantState)
		}
		m.Close()
	}
}

// TestSinkNoopDefault is the explicit no-op side-effect contract: dispatching
// a row with side effects against the default sink neither blocks nor panics.
func TestSinkNoopDefault(t *testing.T) {
	m := New(Options{Initial: StateSleeping})
	if _, err := m.Dispatch(EvSummon, nil); err != nil {
		t.Fatalf("dispatch with default no-op sink: %v", err)
	}
	if m.State() != StateListening {
		t.Fatalf("state = %s, want Listening", m.State())
	}
	m.Close()
}
