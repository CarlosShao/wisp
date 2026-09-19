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
- [ ] 5 concurrent tasks × L2 requests: FIFO order preserved, head-only asserted, all decisions
      routed correctly under shuffled replies.
- [ ] Timeout ladder: 270s warning visible; 300s auto-reject; replay completes the original
      call; independent timers (two items timeout at different times).
- [ ] Fail-closed: kill bridge mid-pending → auto-reject + task survives.
- [ ] Suspend fixture: pending approval voided + notice; no stale approval accepted post-resume.
- [ ] Badge/panel depth sync <100ms; empty-queue transitions correct.

## Progress log (append-only, newest last)
