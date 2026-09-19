package proc

import (
	"context"
	"errors"
	"log/slog"
	"runtime/debug"
	"time"
)

// Shutdown sequence - contract D38(e) / SPEC-01 §4.4. The order is frozen;
// executing steps out of order is an incident, not a style choice (close DB
// before flushing logs -> crash; kill the main process before waiting for
// children -> orphan msedgewebview2.exe resident).
//
//	1. TaskScheduler closes its door (no new tasks)
//	2. hotkey + wake-word listening stops (no new sessions)
//	3. all task root ctxs cancelled -> wait <= 3s (abandon and record the
//	   unfinished tasks beyond that - cgo inference cannot be interrupted,
//	   D37c)
//	4. audio capture thread stops -> device closed   <- BEFORE step 5, so ASR
//	   never sees a half buffer
//	5. ASR/TTS/KWS sessions released (DisposalScope reverse order)
//	6. panel WebView destroyed (if alive) -> children exit <= 2s
//	7. flush logs -> db-writer drained -> wal_checkpoint(TRUNCATE) -> SQLite closed
//	8. debug.FreeOSMemory()
//	9. Job Object closed (children killed by the OS)
//	10. exit (recorded here; the caller performs the actual exit)
//
// The fast path (WM_QUERYENDSESSION, D42#7) skips ONLY the non-critical
// flushes of step 7; steps 4, 5 and 9 are never skipped.

// ShutdownStep numbers the frozen D38(e) sequence, 1..10.
type ShutdownStep int

const (
	StepSchedulerClose ShutdownStep = 1 + iota
	StepStopHotkeyKWS
	StepCancelTasks
	StepStopAudio
	StepReleaseSpeechSessions
	StepDestroyPanel
	StepFlushAndCloseDB
	StepFreeOSMemory
	StepCloseJob
	StepExit
)

var shutdownStepNames = [...]string{
	"scheduler-close",
	"stop-hotkey-kws",
	"cancel-task-roots",
	"stop-audio",
	"release-speech-sessions",
	"destroy-panel-webview",
	"flush-logs-close-db",
	"free-os-memory",
	"close-job-object",
	"exit",
}

// Name returns the audit-trail name of the step.
func (s ShutdownStep) Name() string {
	if s >= StepSchedulerClose && s <= StepExit {
		return shutdownStepNames[s-1]
	}
	return "unknown"
}

// Contract timeouts of the bounded steps (D38e).
const (
	TaskWaitTimeout  = 3 * time.Second // step 3
	PanelWaitTimeout = 2 * time.Second // step 6
)

// ShutdownHooks carry the per-step actions; each is owned by the module named
// in the step. A nil hook means the owning module is not implemented yet and
// is recorded as skipped - the sequence still walks every step so the audit
// trail always shows the full 1..10 order.
type ShutdownHooks struct {
	SchedulerClose        func(ctx context.Context) error // step 1 (agent/scheduler)
	StopHotkeyAndKWS      func(ctx context.Context) error // step 2 (ball)
	CancelTasks           func(ctx context.Context) error // step 3 (agent/scheduler); MUST honor the ctx deadline: exceeding it abandons unfinished tasks
	StopAudio             func(ctx context.Context) error // step 4 (audio)
	ReleaseSpeechSessions func(ctx context.Context) error // step 5 (speech, via DisposalScope)
	DestroyPanel          func(ctx context.Context) error // step 6 (panel); MUST honor the ctx deadline
	FlushAndCloseDB       func(ctx context.Context) error // step 7 (observe+memory); skipped on the fast path
	CloseJob              func(ctx context.Context) error // step 9 (proc.JobScope.Close)
}

// ShutdownOptions parameterize the sequence (zero value = contract defaults).
type ShutdownOptions struct {
	Fast             bool          // WM_QUERYENDSESSION quick variant (D42#7)
	TaskWaitTimeout  time.Duration // default TaskWaitTimeout
	PanelWaitTimeout time.Duration // default PanelWaitTimeout
	FreeOSMemory     func()        // step 8; default debug.FreeOSMemory, injectable for tests (C11 step 4 is the same operation)
}

// StepRecord is the audit entry of one executed, skipped or failed step.
type StepRecord struct {
	Step    ShutdownStep
	Name    string
	Skipped bool
	Err     error
	Elapsed time.Duration
}

// RunShutdownSequence executes steps 1..10 strictly in D38(e) order and
// returns the audit trail. A failing step is logged and the sequence
// CONTINUES: a broken step must not block the teardown of later resources.
// Step 10 is recorded but not executed; the caller exits.
func RunShutdownSequence(hooks ShutdownHooks, opts ShutdownOptions) []StepRecord {
	if opts.TaskWaitTimeout <= 0 {
		opts.TaskWaitTimeout = TaskWaitTimeout
	}
	if opts.PanelWaitTimeout <= 0 {
		opts.PanelWaitTimeout = PanelWaitTimeout
	}
	if opts.FreeOSMemory == nil {
		opts.FreeOSMemory = debug.FreeOSMemory
	}

	records := make([]StepRecord, 0, 10)

	run := func(step ShutdownStep, h func(context.Context) error, timeout time.Duration) {
		rec := StepRecord{Step: step, Name: step.Name()}
		if h == nil {
			rec.Skipped = true
			records = append(records, rec)
			slog.Info("shutdown step skipped (module not present)", "step", int(step), "name", rec.Name)
			return
		}
		start := time.Now()
		ctx := context.Background()
		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
		err := h(ctx)
		rec.Elapsed = time.Since(start)
		if err != nil {
			rec.Err = err
			if errors.Is(err, context.DeadlineExceeded) {
				slog.Error("shutdown step abandoned wait (deadline exceeded)", "step", int(step), "name", rec.Name)
			} else {
				slog.Error("shutdown step failed", "step", int(step), "name", rec.Name, "err", err)
			}
		}
		records = append(records, rec)
	}

	// 1-2: stop producing work.
	run(StepSchedulerClose, hooks.SchedulerClose, 0)
	run(StepStopHotkeyKWS, hooks.StopHotkeyAndKWS, 0)

	// 3: cancel all task root ctxs, bounded wait (abandon + record beyond).
	run(StepCancelTasks, hooks.CancelTasks, opts.TaskWaitTimeout)

	// 4: audio OFF before ASR release (never sees a half buffer).
	// 5: native speech sessions via DisposalScope.
	// (Steps 4/5/9 are never skipped, not even on the fast path.)
	run(StepStopAudio, hooks.StopAudio, 0)
	run(StepReleaseSpeechSessions, hooks.ReleaseSpeechSessions, 0)

	// 6: panel WebView gone, children waited (bounded).
	run(StepDestroyPanel, hooks.DestroyPanel, opts.PanelWaitTimeout)

	// 7: flushes - the ONLY step the fast path skips.
	if opts.Fast {
		records = append(records, StepRecord{Step: StepFlushAndCloseDB, Name: StepFlushAndCloseDB.Name(), Skipped: true})
		slog.Warn("fast shutdown: skipping non-critical log/db flushes (D38e fast path)")
	} else {
		run(StepFlushAndCloseDB, hooks.FlushAndCloseDB, 0)
	}

	// 8: return Go heap to the OS (C11 step 4; D32 fall-back target).
	start := time.Now()
	opts.FreeOSMemory()
	records = append(records, StepRecord{Step: StepFreeOSMemory, Name: StepFreeOSMemory.Name(), Elapsed: time.Since(start)})

	// 9: Job Object closed - the OS kills whatever children remain.
	run(StepCloseJob, hooks.CloseJob, 0)

	// 10: exit is the caller's job.
	records = append(records, StepRecord{Step: StepExit, Name: StepExit.Name()})
	return records
}
