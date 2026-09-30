package proc

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
)

// The hook registration seam D38(e) left open on purpose.
//
// WHY THIS FILE EXISTS. shutdown.go:76-89 declares ShutdownHooks and says a nil
// hook "means the owning module is not implemented yet and is recorded as
// skipped". boot_windows.go's comment above Shutdown said, verbatim, "Later
// tickets register their hooks on the same sequence - the order itself is frozen
// and audited" - but no register existed: Runtime.Shutdown built a fresh
// ShutdownHooks{} every call and filled exactly one slot (CloseJob), so seven of
// the ten steps had no producer and the audit could only ever say "skipped".
// This is that registration entry, and nothing else: it adds no step, renumbers
// no step and cannot reorder anything, because RunShutdownSequence still owns the
// order (shutdown.go:112, pinned by TestShutdownOrderAudit).
//
// WHAT A REGISTRATION MAY NOT DO.
//   - overwrite silently: the first hook for a slot wins and the second is an
//     error. Two owners of one step is the shape where a teardown quietly stops
//     happening, and a step that lies about being executed is worse than a step
//     honestly recorded as skipped.
//   - claim a slot that does not exist: step 8 (FreeOSMemory) and step 10 (exit)
//     are not hookable - the first is this package's own call, the second is the
//     caller's os.Exit.
//   - arrive after the sequence has been read: Hooks() seals the set, so a hook
//     registered during or after shutdown is refused rather than silently missed.

// ErrHookSlotUnknown means a caller named a step that carries no hook.
var ErrHookSlotUnknown = errors.New("proc: 该退出步骤没有钩子位（第 8 步由本包执行，第 10 步由调用方退出）")

// ErrHookAlreadyRegistered means the slot already has an owner. It is returned
// to the second registrant; the first hook keeps running.
var ErrHookAlreadyRegistered = errors.New("proc: 该退出步骤已注册钩子，不接受第二个所有者")

// ErrHookSetSealed means the shutdown sequence has already read this set.
var ErrHookSetSealed = errors.New("proc: 退出序列已开始，钩子注册面已封闭")

// hookableSteps is the closed roster of slots ShutdownHooks carries, keyed by the
// step each one runs at. It is derived from the frozen sequence's own constants,
// so it cannot drift into offering a ninth hookable step.
var hookableSteps = map[ShutdownStep]bool{
	StepSchedulerClose:        true,
	StepStopHotkeyKWS:         true,
	StepCancelTasks:           true,
	StepStopAudio:             true,
	StepReleaseSpeechSessions: true,
	StepDestroyPanel:          true,
	StepFlushAndCloseDB:       true,
	StepCloseJob:              true,
}

// ShutdownHookSet holds the hooks one runtime registered. Its zero value is
// ready to use.
type ShutdownHookSet struct {
	mu     sync.Mutex
	hooks  ShutdownHooks
	filled map[ShutdownStep]bool
	sealed bool
}

// Register attaches one hook to the frozen D38(e) sequence at the step that
// owns it. The hook is handed the step's contract deadline (3s for step 3, 2s
// for step 6) and MUST honor it - shutdown.go's field comment is the authority
// on that, and RunShutdownSequence records a blown deadline as an abandonment,
// not as a hang.
func (hs *ShutdownHookSet) Register(step ShutdownStep, hook func(ctx context.Context) error) error {
	if hook == nil {
		// A nil "hook" is indistinguishable from "nobody registered this", which
		// is the skipped record this whole file exists to make honest.
		return fmt.Errorf("proc: 第 %d 步的钩子为空，注册无意义（不注册就是 skipped）", int(step))
	}
	hs.mu.Lock()
	defer hs.mu.Unlock()
	if hs.sealed {
		return fmt.Errorf("proc: 注册第 %d 步失败: %w", int(step), ErrHookSetSealed)
	}
	if !hookableSteps[step] {
		return fmt.Errorf("proc: 注册第 %d 步失败: %w", int(step), ErrHookSlotUnknown)
	}
	if hs.filled[step] {
		return fmt.Errorf("proc: 注册第 %d 步（%s）失败: %w", int(step), step.Name(), ErrHookAlreadyRegistered)
	}
	if hs.filled == nil {
		hs.filled = map[ShutdownStep]bool{}
	}
	hs.filled[step] = true
	hs.setSlotLocked(step, hook)
	return nil
}

// setSlotLocked writes one field of the struct. The switch is exhaustive over
// hookableSteps by construction: a step in that roster that has no case here is
// a compile-visible gap, and TestShutdownHookSetCoversEveryNamedStep keeps it
// from becoming a silent "registered but never runs".
func (hs *ShutdownHookSet) setSlotLocked(step ShutdownStep, hook func(context.Context) error) {
	switch step {
	case StepSchedulerClose:
		hs.hooks.SchedulerClose = hook
	case StepStopHotkeyKWS:
		hs.hooks.StopHotkeyAndKWS = hook
	case StepCancelTasks:
		hs.hooks.CancelTasks = hook
	case StepStopAudio:
		hs.hooks.StopAudio = hook
	case StepReleaseSpeechSessions:
		hs.hooks.ReleaseSpeechSessions = hook
	case StepDestroyPanel:
		hs.hooks.DestroyPanel = hook
	case StepFlushAndCloseDB:
		hs.hooks.FlushAndCloseDB = hook
	case StepCloseJob:
		hs.hooks.CloseJob = hook
	}
}

// Hooks snapshots the registered set and seals it against further registration.
// It is what Runtime.Shutdown reads; RunShutdownSequence then walks the frozen
// order over it.
func (hs *ShutdownHookSet) Hooks() ShutdownHooks {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	hs.sealed = true
	return hs.hooks
}

// Registered lists the steps that currently have an owner, ascending. It is the
// read-only answer to "which of the ten steps does this process really run", so
// a boot report can state it instead of implying it.
func (hs *ShutdownHookSet) Registered() []ShutdownStep {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	out := make([]ShutdownStep, 0, len(hs.filled))
	for s := range hs.filled {
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}
