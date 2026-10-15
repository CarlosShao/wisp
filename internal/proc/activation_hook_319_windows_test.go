//go:build windows

package proc

// Ticket 319 AC#2/AC#3 (N2), the HEADLESS half of "a second launch must be able
// to raise the running instance's UI".
//
// The defect (ticket 319; census section 3): the activation branch in boot's
// event loop booked the signal and cleared it, and did nothing else. Its own log
// line said so and a comment said why - no path connects this leaf package to a
// window. The owner's log carries 12 "activation requested by second launch"
// entries for 2026-10-14, including one at 21:40:38 where the panel thread and
// its window were BOTH alive, so the silence was not the panel being dead: the
// signal genuinely had nowhere to go.
//
// The fix this file nails: Runtime.SetActivationHook is the same injection shape
// RegisterShutdownHook uses (a function value, no new dependency edge - internal
// /proc stays a leaf), and one round of handleActivation runs it exactly once
// and clears the manual-reset event so the next launch can signal again.
//
// Mutation control: delete the hook invocation in handleActivation (leave the
// "booked and cleared only" branch as the one path) and the count assertion
// below goes red - the hook then runs zero times and the case says so by name.

import (
	"sync/atomic"
	"testing"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/observe"
)

// bootForActivation319 boots a production runtime with per-run mutex/event
// names, so this case can set its own activation event without ever touching the
// real dev/prod names (including the owner's running instance).
func bootForActivation319(t *testing.T) *Runtime {
	t.Helper()
	rt, err := Boot(buildinfo.EnvProd, WithLayout(uniqueLayout(t)), WithRegistry(observe.NewRegistry()))
	if err != nil {
		t.Fatalf("Boot: %v", err)
	}
	t.Cleanup(func() { rt.Shutdown(true) })
	if rt.Instance == nil {
		t.Fatalf("this layout owns a mutex, so the instance must be non-nil before the activation cases can mean anything")
	}
	return rt
}

// TestActivationHook319RunsExactlyOnceAndResets is N2: one signal -> one hook
// call -> the event is quiet again.
func TestActivationHook319RunsExactlyOnceAndResets(t *testing.T) {
	rt := bootForActivation319(t)

	var calls atomic.Int32
	rt.SetActivationHook(func() string {
		calls.Add(1)
		return "panel toggle accepted (second-launch)"
	})

	// A quiet round: nothing was signalled, so nothing may be consumed and the
	// hook may not run. This is also the false-positive guard for the counter
	// below.
	if rt.handleActivation() {
		t.Errorf("handleActivation consumed an activation nobody sent")
	}
	if n := calls.Load(); n != 0 {
		t.Fatalf("the activation hook ran %d time(s) with no signal", n)
	}

	// The second launch: SignalExistingInstance is the same call cmd/wisp's
	// ErrAlreadyRunning branch makes.
	if err := SignalExistingInstance(rt.Layout.ActivateEventName); err != nil {
		t.Fatalf("SignalExistingInstance(%q): %v", rt.Layout.ActivateEventName, err)
	}
	if !rt.handleActivation() {
		t.Fatalf("one round of the pump did not consume the activation signal - the loop would never raise anything")
	}
	if n := calls.Load(); n != 1 {
		t.Fatalf("the registered activation hook ran %d time(s) for ONE second launch, want exactly 1", n)
	}

	// ResetActivation ran: a second round must find the event quiet, and the
	// hook must not run again off the same signal.
	if rt.handleActivation() {
		t.Errorf("a second round consumed a signal: ResetActivation did not clear the manual-reset event")
	}
	if n := calls.Load(); n != 1 {
		t.Fatalf("the activation hook ran %d time(s) across two rounds, want 1 - one launch's signal leaked into two handoffs", n)
	}

	// And a second launch really can signal again (the manual-reset event was
	// left armed for the next one).
	if err := SignalExistingInstance(rt.Layout.ActivateEventName); err != nil {
		t.Fatalf("second SignalExistingInstance: %v", err)
	}
	if !rt.handleActivation() {
		t.Fatalf("the event did not signal again after ResetActivation")
	}
	if n := calls.Load(); n != 2 {
		t.Fatalf("after a second launch the hook count is %d, want 2", n)
	}
}

// TestActivationWithoutHook319StillConsumesAndClears is the other half of the
// contract: a runtime nobody handed a hook must keep today's behaviour exactly -
// the signal is booked and cleared, and no invented work happens. It is the
// positive control for the failure branch in handleActivation.
func TestActivationWithoutHook319StillConsumesAndClears(t *testing.T) {
	rt := bootForActivation319(t)
	if err := SignalExistingInstance(rt.Layout.ActivateEventName); err != nil {
		t.Fatalf("SignalExistingInstance(%q): %v", rt.Layout.ActivateEventName, err)
	}
	if !rt.handleActivation() {
		t.Fatalf("the no-hook branch did not consume the activation signal")
	}
	if rt.handleActivation() {
		t.Errorf("the no-hook branch did not clear the event: a second round consumed the same signal")
	}
}
