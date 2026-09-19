# 47 — TaskScheduler + PathLock (C20): multi-task concurrency unlocks here

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 10-agent-loop-core, 21-approval-gates-minimal
**Parallel slots:** ≤2 sub-agents (A: scheduler + root ctx + WaitGroup completion; B: PathLock +
conflict queueing + task tools)
**Spec refs:** SPEC-06 §7, D31, C20, D43 #20/#39, D38(c)(d), S7
⚠ Unlocks multi-task: tickets 21–32 ran single-task by design (safety-incomplete period).

## What to build
Concurrent task execution: per-task root contexts with named goroutine trios, WaitGroup-based
completion, tool concurrency cap 4, and PathLock — path-level conflict detection queuing the
later task with a visible reason.

## Key constraints
- Scheduler: FIFO accept + concurrent run; per task: root ctx + goroutines agent-task-<id> /
  tool-exec-<id> / approval-waiter (D38b roster); completion REQUIRES WaitGroup drain (RSS
  settle depends on it); task cancel → root ctx + applied-steps report.
- Tool concurrency ≤4 across a task (D38d); per-task LoopGuard/budgets independent (from 10);
  cross-task LoopGuard isolation tested (task A's repeats don't trip task B).
- PathLock (C20): before path-touching tool calls, register normalized target paths (via C26);
  conflict → later task QUEUED with visible "在等 <task> 释放 <path>" (panel task list + ball
  Queued dot); release → #39 transition; deadlock guard (two tasks mutually waiting → timeout
  + explicit error, no silent deadlock).
- Approval-wait does NOT block other tasks' non-conflicting work (D31 core property).
- `task.list` (L0) / `task.cancel` (L1) builtin tools expose scheduler state to the agent.
- Single-writer SQLite unaffected (04 db-writer serializes).

## Out of scope
- Multi-item approval queue routing (48); session grants (49); plugin runtimes.

## Acceptance criteria
- [ ] Concurrency: 4 tasks × mixed tool loads → all complete, goroutine count returns to
      6 + 3×active baseline, zero leaks (registry test).
- [ ] PathLock: overlapping file ops → second queued with correct reason → resumes after
      release (#20→#39 state trace).
- [ ] Deadlock fixture (A holds p1 wants p2; B holds p2 wants p1) → timeout + explicit error.
- [ ] Approval independence: task A awaiting L2 while task B executes tooling to completion.
- [ ] Cancel under concurrency: mid-flight cancel → applied-steps report; siblings unaffected.
- [ ] task.list/cancel round-trip via agent (mock LLM drives both).

## Progress log (append-only, newest last)
