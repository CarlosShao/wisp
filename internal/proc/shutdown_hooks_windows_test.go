//go:build windows

package proc

import (
	"context"
	"errors"
	"testing"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Ticket 246 AC#4, runtime half: the D38(e) audit for step 3 has to describe
// what this process did.
//
// Before this seam existed, Runtime.Shutdown built a fresh ShutdownHooks{} on
// every call and filled one field (CloseJob). The consequence is the sentence
// 246-a1 §2③ measured: a step a module never handed over is recorded as skipped,
// and a cancel performed outside the sequence (in a caller's defer) is NOT in the
// trail at all - behaviour and audit forked, and only a mutation could see it.
// The two cases below are the same ruler run twice: registered means executed and
// named in the record, unregistered means still honestly skipped.

func bootRuntimeForHookTest(t *testing.T) *Runtime {
	t.Helper()
	rt, err := Boot(buildinfo.EnvTest, WithRegistry(observe.NewRegistry()))
	if err != nil {
		t.Fatalf("Boot(test): %v", err)
	}
	if rt.Job == nil {
		t.Fatal("Boot returned a runtime with no Job Object")
	}
	return rt
}

// TestRuntimeCancelStepExecutesWhenRegistered is the capability reading.
func TestRuntimeCancelStepExecutesWhenRegistered(t *testing.T) {
	rt := bootRuntimeForHookTest(t)

	var ranWithDeadline bool
	if err := rt.RegisterShutdownHook(StepCancelTasks, func(ctx context.Context) error {
		// Step 3's contract is a bounded wait (shutdown.go:83): the hook is
		// handed a deadline, and a hook that never got one could hang the exit.
		_, ok := ctx.Deadline()
		ranWithDeadline = ok
		return nil
	}); err != nil {
		t.Fatalf("RegisterShutdownHook(step 3): %v", err)
	}
	if got := rt.RegisteredShutdownSteps(); len(got) != 1 || got[0] != StepCancelTasks {
		t.Fatalf("RegisteredShutdownSteps() = %v, want exactly [3]", got)
	}

	records := rt.Shutdown(false)
	if len(records) != 10 {
		t.Fatalf("audit trail = %d records, want the frozen 10", len(records))
	}
	for i, rec := range records {
		if rec.Step != ShutdownStep(i+1) {
			t.Fatalf("record %d = step %d (%s): registering a hook must not reorder anything",
				i, int(rec.Step), rec.Name)
		}
	}
	rec := records[StepCancelTasks-1]
	if rec.Skipped {
		t.Fatalf("step 3 still recorded skipped after registration: %+v", rec)
	}
	if rec.Err != nil {
		t.Fatalf("step 3 reported an error: %+v", rec)
	}
	if !ranWithDeadline {
		t.Fatal("step 3's hook was handed a context with no deadline (the <=3s wait is contractual)")
	}
	// proc's own step 9 is untouched by this: the runtime still closes the Job,
	// which is the one hook it filled before ticket 246 existed.
	if jobRec := records[StepCloseJob-1]; jobRec.Skipped || jobRec.Err != nil || !rt.Job.Closed() {
		t.Fatalf("job close record = %+v, closed = %v", jobRec, rt.Job.Closed())
	}
}

// TestRuntimeCancelStepIsSkippedWhenUnregistered is the positive control: the
// ruler must be able to see the old state, or the reading above proves nothing.
func TestRuntimeCancelStepIsSkippedWhenUnregistered(t *testing.T) {
	rt := bootRuntimeForHookTest(t)
	records := rt.Shutdown(false)
	if got := rt.RegisteredShutdownSteps(); len(got) != 0 {
		t.Fatalf("RegisteredShutdownSteps() = %v on a runtime that registered nothing", got)
	}
	rec := records[StepCancelTasks-1]
	if !rec.Skipped {
		t.Fatalf("step 3 = %+v: an unregistered step must still be recorded skipped", rec)
	}
	if records[StepCloseJob-1].Skipped {
		t.Fatal("the runtime's own CloseJob hook stopped running")
	}
}

// TestRuntimeRegistrationRefusalsReachTheCaller pins that a half-wired teardown
// is an error the assembly root can print, not a quiet no-op.
func TestRuntimeRegistrationRefusalsReachTheCaller(t *testing.T) {
	rt := bootRuntimeForHookTest(t)
	if err := rt.RegisterShutdownHook(StepExit, func(context.Context) error { return nil }); !errors.Is(err, ErrHookSlotUnknown) {
		t.Fatalf("RegisterShutdownHook(step 10) err = %v, want ErrHookSlotUnknown", err)
	}
	if err := rt.RegisterShutdownHook(StepCancelTasks, func(context.Context) error { return nil }); err != nil {
		t.Fatalf("first RegisterShutdownHook: %v", err)
	}
	if err := rt.RegisterShutdownHook(StepCancelTasks, func(context.Context) error { return nil }); !errors.Is(err, ErrHookAlreadyRegistered) {
		t.Fatalf("duplicate RegisterShutdownHook err = %v, want ErrHookAlreadyRegistered", err)
	}
	rt.Shutdown(false)
	if err := rt.RegisterShutdownHook(StepStopAudio, func(context.Context) error { return nil }); !errors.Is(err, ErrHookSetSealed) {
		t.Fatalf("RegisterShutdownHook after Shutdown err = %v, want ErrHookSetSealed", err)
	}
}
