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
- [x] Concurrency: 4 tasks × mixed tool loads → all complete, goroutine count returns to
      6 + 3×active baseline, zero leaks (registry test).
- [x] PathLock: overlapping file ops → second queued with correct reason → resumes after
      release (#20→#39 state trace).
- [x] Deadlock fixture (A holds p1 wants p2; B holds p2 wants p1) → timeout + explicit error.
- [x] Approval independence: task A awaiting L2 while task B executes tooling to completion.
- [x] Cancel under concurrency: mid-flight cancel → applied-steps report; siblings unaffected.
- [x] task.list/cancel round-trip via agent (mock LLM drives both).

## Progress log (append-only, newest last)

- 2026-11-06 22:4x +08 W47（wave45 写腿）：6 AC 全落。写面＝internal/agent/scheduler（scheduler.go/pathlock.go/bridge.go/goos_*，14 测试）＋internal/tools（task.list 注册＋task.cancel 调度器分支＋TaskDeps.Scheduler 窄面，3 测试＋五枚尺同批翻转）。AC#6 的 task.list 按 A434 同规（标记与接线同一次核销）解冻：done-判据=PLAN.md:1537 行，D34 行 L0 本就冻结、契约文字零改。D32 零新增 timer（死锁上限=等待任务自身 goroutine 的单调 deadline，可停）；D38b 唯一 spawn=agent-task-<id>（229 普查登记随落）。cmd/wisp 装配不动（盘上有 323 未跟踪件，规则 21 不触发；装配缝=NewBridgeDeps/ExecAdapter 留给装配票）。门禁：gofumpt 空／vet rc=0／scheduler+tools -count=2 ok（1.2s+55.5s）／d22scan clean rc=0／agent+observe 回归 ok／cmd/wisp allow-door+174 尺 ok。件 docs/evidence/s1/47-scheduler-pathlock.md。零 GUI、零 wisp.exe、零 winlive、零 push。
