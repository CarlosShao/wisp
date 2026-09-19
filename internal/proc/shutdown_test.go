package proc

import (
	"context"
	"testing"
	"time"
)

// TestShutdownOrderAudit is the SPEC-01 §7 order-audit test: the actual
// execution order must equal the frozen D38(e) sequence 1..10.
func TestShutdownOrderAudit(t *testing.T) {
	var order []ShutdownStep
	hooks := ShutdownHooks{}
	mk := func(step ShutdownStep) func(context.Context) error {
		return func(context.Context) error {
			order = append(order, step)
			return nil
		}
	}
	hooks.SchedulerClose = mk(StepSchedulerClose)
	hooks.StopHotkeyAndKWS = mk(StepStopHotkeyKWS)
	hooks.CancelTasks = mk(StepCancelTasks)
	hooks.StopAudio = mk(StepStopAudio)
	hooks.ReleaseSpeechSessions = mk(StepReleaseSpeechSessions)
	hooks.DestroyPanel = mk(StepDestroyPanel)
	hooks.FlushAndCloseDB = mk(StepFlushAndCloseDB)
	hooks.CloseJob = mk(StepCloseJob)

	freeRan := false
	records := RunShutdownSequence(hooks, ShutdownOptions{FreeOSMemory: func() {
		order = append(order, StepFreeOSMemory)
		freeRan = true
	}})

	// Steps 1..9 execute in frozen order; step 10 (exit) is recorded only.
	want := []ShutdownStep{
		StepSchedulerClose, StepStopHotkeyKWS, StepCancelTasks, StepStopAudio,
		StepReleaseSpeechSessions, StepDestroyPanel, StepFlushAndCloseDB,
		StepFreeOSMemory, StepCloseJob,
	}
	if len(order) != len(want) {
		t.Fatalf("executed %d steps, want %d: %v", len(order), len(want), order)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("step %d executed = %s, want %s (full order %v)", i+1, order[i].Name(), want[i].Name(), order)
		}
	}
	if !freeRan {
		t.Fatal("step 8 (FreeOSMemory) did not run")
	}
	if len(records) != 10 {
		t.Fatalf("audit trail has %d records, want 10", len(records))
	}
	wantRecords := append(append([]ShutdownStep{}, want...), StepExit)
	for i, rec := range records {
		if rec.Step != wantRecords[i] || rec.Skipped || rec.Err != nil {
			t.Fatalf("record %d = %+v, want step %s executed cleanly", i, rec, wantRecords[i].Name())
		}
	}
	// The audio->ASR ordering contract: step 4 strictly before step 5.
	if order[3] != StepStopAudio || order[4] != StepReleaseSpeechSessions {
		t.Fatal("audio must stop before speech sessions are released")
	}
	// Step 10 (exit) is recorded, never executed by the sequence.
	if records[9].Step != StepExit {
		t.Fatal("last record must be exit")
	}
}

// TestShutdownFastPathSkipsOnlyStep7 pins the D42#7 rule: WM_QUERYENDSESSION
// skips the non-critical flushes of step 7 ONLY; steps 4/5/9 always run.
func TestShutdownFastPathSkipsOnlyStep7(t *testing.T) {
	var order []ShutdownStep
	mk := func(step ShutdownStep) func(context.Context) error {
		return func(context.Context) error {
			order = append(order, step)
			return nil
		}
	}
	hooks := ShutdownHooks{
		SchedulerClose:        mk(StepSchedulerClose),
		StopHotkeyAndKWS:      mk(StepStopHotkeyKWS),
		CancelTasks:           mk(StepCancelTasks),
		StopAudio:             mk(StepStopAudio),
		ReleaseSpeechSessions: mk(StepReleaseSpeechSessions),
		DestroyPanel:          mk(StepDestroyPanel),
		FlushAndCloseDB:       mk(StepFlushAndCloseDB),
		CloseJob:              mk(StepCloseJob),
	}
	records := RunShutdownSequence(hooks, ShutdownOptions{
		Fast:         true,
		FreeOSMemory: func() { order = append(order, StepFreeOSMemory) },
	})

	var executed, skipped []ShutdownStep
	for _, rec := range records {
		if rec.Skipped {
			skipped = append(skipped, rec.Step)
		} else if rec.Step != StepExit {
			executed = append(executed, rec.Step)
		}
	}
	if len(skipped) != 1 || skipped[0] != StepFlushAndCloseDB {
		t.Fatalf("skipped steps = %v, want only flush-logs-close-db", skipped)
	}
	for _, mustRun := range []ShutdownStep{StepStopAudio, StepReleaseSpeechSessions, StepFreeOSMemory, StepCloseJob} {
		found := false
		for _, s := range executed {
			if s == mustRun {
				found = true
			}
		}
		if !found {
			t.Fatalf("fast path must not skip %s (executed: %v)", mustRun.Name(), executed)
		}
	}
}

// TestShutdownBoundedWaitAbandonsAndContinues pins step 3/6 semantics: a hook
// that honors its ctx deadline gets the contract timeout; exceeding it is
// recorded and the sequence still completes all 10 steps.
func TestShutdownBoundedWaitAbandonsAndContinues(t *testing.T) {
	hooks := ShutdownHooks{
		CancelTasks: func(ctx context.Context) error {
			<-ctx.Done() // tasks that cannot be interrupted (cgo, D37c)
			return ctx.Err()
		},
		DestroyPanel: func(ctx context.Context) error {
			<-ctx.Done()
			return ctx.Err()
		},
	}
	start := time.Now()
	records := RunShutdownSequence(hooks, ShutdownOptions{
		TaskWaitTimeout:  80 * time.Millisecond,
		PanelWaitTimeout: 60 * time.Millisecond,
		FreeOSMemory:     func() {},
	})
	elapsed := time.Since(start)

	if elapsed > 2*time.Second {
		t.Fatalf("sequence took %v; bounded waits must cap it", elapsed)
	}
	for _, want := range []ShutdownStep{StepCancelTasks, StepDestroyPanel} {
		rec := records[want-1]
		if rec.Step != want || rec.Err == nil {
			t.Fatalf("record for %s = %+v, want abandonment error", want.Name(), rec)
		}
	}
	if last := records[9]; last.Step != StepExit {
		t.Fatalf("sequence did not reach step 10: %+v", last)
	}
}

// TestShutdownNilHooksAreRecordedSkipped: a module that has not landed yet
// shows up as a skipped step, and the audit trail is still complete.
func TestShutdownNilHooksAreRecordedSkipped(t *testing.T) {
	records := RunShutdownSequence(ShutdownHooks{}, ShutdownOptions{FreeOSMemory: func() {}})
	if len(records) != 10 {
		t.Fatalf("records = %d, want 10", len(records))
	}
	for _, rec := range records {
		switch rec.Step {
		case StepFreeOSMemory, StepExit:
			if rec.Skipped {
				t.Fatalf("%s must never be skipped", rec.Name)
			}
		default:
			if !rec.Skipped {
				t.Fatalf("%s should be skipped without a hook", rec.Name)
			}
		}
	}
}
