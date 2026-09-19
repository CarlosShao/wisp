# 54 — S7 acceptance: concurrency red-team, lifecycle pre-mortems, self-use week 1

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 47-task-scheduler-pathlock, 48-approval-queue-full, 49-session-grants-d45-2,
53-larkcli-plugin-e2e
**Parallel slots:** ≤2 sub-agents (A: concurrency/approval matrix; B: lifecycle pre-mortems +
rollback)
**Spec refs:** S7 done (D44), §10 #7, D42#3, D40 sequences, S7 close-out

## What to build
The S7 slice gate: concurrency behaviors under adversarial schedules, the full lifecycle
pre-mortem matrix (display topology, DPI, crash, suspend-with-pending-approval, updater
round-trip with N-1 rollback), and the final pending-decision close-outs.

## Key constraints
- Concurrency matrix: 2–5 tasks with interleaved L2 requests → FIFO + correlationId routing
  correct under randomized decision timing; PathLock conflicts queue visibly; timeouts fire
  independently; LoopGuard fires per-task; a task awaiting approval never starves others.
- Pre-mortems (real triggers): monitor unplug/replug mid-task (ball repositioned, D2D rebuilt,
  no crash); DPI change; kill -9 during update staging (rollback path intact); suspend with
  pending approval (void+notice per 48 decision); GPU/driver reset fixture (WebView recovery);
  two-instance race on boot.
- Update round-trip rehearsal (56-prep): staging → verify → swap → N-1 rollback drill with
  crash-counter fixture.
- Close-out decisions due: §16.11#10 suspend-approval final wording (48 implemented it —
  ratify in D43 doc); spike-driven SLO table marked "verified"; any remaining OPEN items
  re-registered with owners.
- Scenario↔slice matrix final bidirectional check (D44(b)); all four scenarios re-run once
  post-concurrency (regression guard: enabling concurrency didn't break single-task paths).
- Adversarial acceptance + gap audit by different agents; self-use week begins only after
  this gate (final acceptance is user's week-long real usage per §10 #9).

## Out of scope
- S8 items (signing/macOS/i18n/SDK — deferred tickets).

## Acceptance criteria
- [ ] Concurrency matrix green with evidence; no cross-task interference found.
- [ ] All pre-mortems executed with recordings/logs + correct error classes.
- [ ] Update+rollback drill: staged update applies; forced-crash ×2 → auto-rollback to N-1.
- [ ] Open-decision close-outs documented (D43 addendum ratified or escalated).
- [ ] Four-scenario regression re-run green; adversarial + gap-audit reports committed.

## Progress log (append-only, newest last)
