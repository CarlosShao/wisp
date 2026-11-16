# 48 — ApprovalQueue full (C18): FIFO multi-task, head-only display, fail-closed routing

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 37-approval-ui-l2, 47-task-scheduler-pathlock
**Parallel slots:** 1
**Spec refs:** SPEC-06 §7, C18, D31, D43 #23–25, §16.11#10, S7

## What to build
The production ApprovalQueue: global FIFO across concurrent tasks, correlationId-routed
decisions, head-only display with depth, 300s timeout→reject with 270s warning, fail-closed on
host unavailability, replay-after-reject — upgrading 21's trivial queue.

## Key constraints
- Every request: correlationId + task id + tool + FULL params + risk + reason; replies routed
  strictly by correlationId (mis-routed reply = test failure); queue depth on ball badge
  (2+ digits correct); head-only in every UI surface.
- Timeout 300s → auto-REJECT (never default-allow); 270s prominent warning; reject → that tool
  call fails, task root ctx SURVIVES, replay re-executes the call (idempotency note where
  relevant); timeouts independent per item (multi-pending).
- Fail-closed: bridge/queue internal error or host unavailable → treat as reject (never allow).
- Suspend behavior final decision (§16.11#10): on system suspend → pending approvals VOIDED +
  auto-reject + "错过 N 个待确认操作" notice (43 wires the event; implement + document as the
  D43 table addendum; human review at 54).
- AwaitingApproval depth transitions per D43 #23–25; queue empty + task done → Speaking/
  Settling.
- Decision persistence: allow/reject/timeout rows with correlationId + grant linkage for 49.

## Out of scope
- Grant session-auth options in the card (49).

## Acceptance criteria
- [x] 5 concurrent tasks × L2 requests: FIFO order preserved, head-only asserted, all decisions
      routed correctly under shuffled replies.
- [x] Timeout ladder: 270s warning visible; 300s auto-reject; replay completes the original
      call; independent timers (two items timeout at different times).
- [x] Fail-closed: kill bridge mid-pending → auto-reject + task survives.
- [x] Suspend fixture: pending approval voided + notice; no stale approval accepted post-resume.
- [x] Badge/panel depth sync <100ms; empty-queue transitions correct.

## Progress log (append-only, newest last)

---

## 编排者收件（2026-11-16 21:1x，`f622b5f29`，件 `docs/evidence/s1/48-approval-queue-c18-r1.md` 12788 B）

**尺1** `f622b5f29` 在（3 件：两枚测试＋证据）。**尺2** 件在 12788 B。**尺3 我自己复跑（＝非实现者复量）**：`go test ./internal/agent/approval -run TestTicket48` ⇒ **ok 0.224s**；载具 `-test.run '^Test48'` ⇒ **2/2 PASS rc=0**（`Test48SuspendFixtureVoidsPendingAndResumeRefusesStale` 0.03s／`Test48PanelRosterJoinDepthSyncUnder100ms` 0.02s）。
**关键读数（腿交、含注入时钟，⛔ 墙钟）**：FIFO 5 任务次序 1..5、只显队头、乱序答复逐枚回原任务、卡1 grant 花到卡3＝`ErrBadGrant`；超时梯 270s 警告→300s 自动拒→replay 新 corr/grant；fail-closed 红句「任务上下文已结束，审批请求已作废并按拒绝处理」；挂起 voided=2＋恢复后死 grant/死 corr 均 `ErrUnknownCorrelation`；深度 1→12 逐层、**额外轮询=0**。
**整包**：`go test ./cmd/wisp` 经载具一发 **RUN=512/PASS=378/FAIL=0/SKIP=1** 后在 WebView2 区段 **exit 139**（＝**票 450 已知**既定死法，非本腿引入）；A1893 复量截断点那枚 `TestPanelThreadIsSTAAndExitsCleanly` 单跑 **1.29s PASS**；**A1882**：终点锚整包未命中（死在锚点区段前），载具单跑三枚锚 **3/3 PASS**，新测试 Test48* 两枚在截断点前 PASS。环境红 0 枚。
⇒ 5/5，**本票结**。零生产码改动、未涉契约文字。