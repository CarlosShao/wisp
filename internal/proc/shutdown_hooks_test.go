package proc

import (
	"context"
	"errors"
	"testing"
)

// Ticket 246 AC#4's registration seam, measured on the sequence itself.
//
// The claim under test is narrow and it is not about order: RunShutdownSequence
// already pinned the 1..10 order (TestShutdownOrderAudit), and this file must not
// move it either. What had no test and no producer before ticket 246 is the word
// "register" in boot_windows.go's own sentence ("Later tickets register their
// hooks on the same sequence"): ShutdownHooks existed, ShutdownHooks could be
// executed, and nothing could be FILLED by a caller.

// allHookSteps is the closed roster ShutdownHooks carries, in sequence order -
// the same eight steps hookableSteps names. Written out here rather than derived
// so a future ninth field has to be added in two places on purpose.
var allHookSteps = []ShutdownStep{
	StepSchedulerClose, StepStopHotkeyKWS, StepCancelTasks, StepStopAudio,
	StepReleaseSpeechSessions, StepDestroyPanel, StepFlushAndCloseDB, StepCloseJob,
}

// TestShutdownHookSetRunsEveryRegisteredStepInOrder is the capability case, and
// its positive control at the same time: every one of the eight slots is filled,
// the sequence must EXECUTE all eight, in the frozen order, with nothing skipped.
// A registry that accepted a hook and then dropped it would fail here, and so
// would one that reordered anything.
func TestShutdownHookSetRunsEveryRegisteredStepInOrder(t *testing.T) {
	var hs ShutdownHookSet
	var order []ShutdownStep
	for _, step := range allHookSteps {
		s := step
		if err := hs.Register(s, func(context.Context) error {
			order = append(order, s)
			return nil
		}); err != nil {
			t.Fatalf("Register(step %d): %v", int(s), err)
		}
	}
	if got := hs.Registered(); len(got) != len(allHookSteps) {
		t.Fatalf("Registered() = %v, want all %d hookable steps", got, len(allHookSteps))
	}

	records := RunShutdownSequence(hs.Hooks(), ShutdownOptions{FreeOSMemory: func() {}})
	if len(records) != 10 {
		t.Fatalf("audit trail = %d records, want 10", len(records))
	}
	if len(order) != len(allHookSteps) {
		t.Fatalf("executed %d hooks, want %d (order %v)", len(order), len(allHookSteps), order)
	}
	for i, want := range allHookSteps {
		if order[i] != want {
			t.Fatalf("hook %d ran as step %s, want %s (full order %v)", i+1, order[i].Name(), want.Name(), order)
		}
		rec := records[want-1]
		if rec.Skipped || rec.Err != nil {
			t.Fatalf("record for %s = %+v, want executed and clean", want.Name(), rec)
		}
	}
}

// TestShutdownHookSetRefusesEveryShapeThatWouldLie pins the four refusals. Each
// one is a way the audit trail could end up describing a teardown that did not
// happen, or losing one that did.
func TestShutdownHookSetRefusesEveryShapeThatWouldLie(t *testing.T) {
	t.Run("nil hook is not a registration", func(t *testing.T) {
		var hs ShutdownHookSet
		if err := hs.Register(StepCancelTasks, nil); err == nil {
			t.Fatal("Register(nil hook) succeeded: a nil hook is indistinguishable from skipped")
		}
		if got := hs.Registered(); len(got) != 0 {
			t.Fatalf("Registered() = %v after a refused registration", got)
		}
	})

	t.Run("steps without a slot are refused", func(t *testing.T) {
		var hs ShutdownHookSet
		for _, step := range []ShutdownStep{StepFreeOSMemory, StepExit, ShutdownStep(0), ShutdownStep(11)} {
			err := hs.Register(step, func(context.Context) error { return nil })
			if !errors.Is(err, ErrHookSlotUnknown) {
				t.Fatalf("Register(step %d) err = %v, want ErrHookSlotUnknown", int(step), err)
			}
		}
	})

	t.Run("one slot has one owner", func(t *testing.T) {
		var hs ShutdownHookSet
		ran := 0
		if err := hs.Register(StepCancelTasks, func(context.Context) error { ran++; return nil }); err != nil {
			t.Fatalf("first Register: %v", err)
		}
		if err := hs.Register(StepCancelTasks, func(context.Context) error {
			t.Error("the second registrant's hook ran: an overwritten owner is a silent loss of teardown")
			return nil
		}); !errors.Is(err, ErrHookAlreadyRegistered) {
			t.Fatalf("second Register err = %v, want ErrHookAlreadyRegistered", err)
		}
		RunShutdownSequence(hs.Hooks(), ShutdownOptions{FreeOSMemory: func() {}})
		if ran != 1 {
			t.Fatalf("first hook ran %d times, want 1", ran)
		}
	})

	t.Run("registration after the sequence started is refused", func(t *testing.T) {
		var hs ShutdownHookSet
		hs.Hooks() // what Runtime.Shutdown does
		if err := hs.Register(StepCancelTasks, func(context.Context) error { return nil }); !errors.Is(err, ErrHookSetSealed) {
			t.Fatalf("Register after Hooks() err = %v, want ErrHookSetSealed", err)
		}
	})
}

// TestShutdownHookSetZeroValueIsUsable is the cheap half of "no hidden setup":
// the set is a field of a struct Boot builds, and a zero value that panicked
// would turn a forgotten initialisation into a crash at the worst possible moment
// (inside shutdown).
func TestShutdownHookSetZeroValueIsUsable(t *testing.T) {
	var empty ShutdownHookSet
	if got := empty.Registered(); len(got) != 0 {
		t.Fatalf("zero value Registered() = %v, want empty", got)
	}
	if h := empty.Hooks(); h.CancelTasks != nil || h.CloseJob != nil || h.SchedulerClose != nil {
		t.Fatal("zero value Hooks() reported hooks nobody registered")
	}

	var hs ShutdownHookSet
	if err := hs.Register(StepStopAudio, func(context.Context) error { return nil }); err != nil {
		t.Fatalf("Register on a zero value: %v", err)
	}
	if got := hs.Registered(); len(got) != 1 || got[0] != StepStopAudio {
		t.Fatalf("Registered() = %v, want [stop-audio]", got)
	}
	if h := hs.Hooks(); h.StopAudio == nil {
		t.Fatal("a registered hook disappeared from the snapshot")
	}
}
