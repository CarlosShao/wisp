// Package scheduler owns the TaskScheduler and PathLock (C20) (SPEC-01 §3).
//
// Responsibilities:
//   - concurrent task management: FIFO accept + concurrent run, one root ctx
//     per task (D38c) with every task goroutine spawned under it through the
//     D38b registry, completion REQUIRES the join counter to drain
//   - PathLock: path-conflict detection routes conflicting tasks to Queued
//     (D43 edge #20) instead of racing, with the visible 「在等 <task> 释放
//     <path>」 reason, and un-blocks them on release (#39); deadlock guard
//     gives an explicit error after the ceiling, never a silent hang
//   - scheduler door: refuses new tasks once shutdown started (D38e step 1)
//     and the CancelAll step-3 shape (cancel all roots, bounded drain)
//   - the task.list / task.cancel scheduler surface through the tools'
//     SchedulerDeps faces (bridge.go)
//
// Non-responsibilities:
//   - no task execution itself (agent loop), no queue UI (panel), no risk
//     decisions (the bridge/C19 judge every call before the lock is asked)
//
// DEFERRED(scheduler) was implemented by ticket 47 (this package now carries
// the scheduler, the PathLock and the tools adapters; the composition wiring
// into cmd/wisp's run.go assembly is the assembly ticket's write face). The
// marker is lifted in the same commit as the code - the same marker/wiring
// same-landing rule ruling A434 set for task.cancel.
package scheduler
