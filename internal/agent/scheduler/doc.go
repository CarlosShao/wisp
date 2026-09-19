// Package scheduler owns the TaskScheduler and PathLock (C20) (SPEC-01 §3).
//
// Responsibilities:
//   - concurrent task management: creates each task's root ctx (D38c), 4-way
//     tool concurrency ceiling (D38d) belongs to tools but is configured here
//   - PathLock: path-conflict detection routes conflicting tasks to Queued
//     (D43 edge #20) instead of racing
//   - scheduler door: refuses new tasks once shutdown started (D38e step 1)
//
// Non-responsibilities:
//   - no task execution itself (agent loop), no queue UI (panel)
//
// DEFERRED(scheduler): implemented by ticket 47. This ticket only freezes the
// package boundary.
package scheduler
