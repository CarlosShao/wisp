# 43 — Power & lifecycle events: suspend/resume, session end, autostart, crash recovery

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 28-session-scope-warm
**Parallel slots:** ≤2 sub-agents (A: power + session-end paths; B: autostart + crash recovery +
update hooks)
**Spec refs:** §6, D42#1#7, D38(e) fast path, D43 pending-approval open item note, S6

## What to build
Windows lifecycle robustness: suspend/resume handling (dispose session, reopen audio), fast
shutdown on WM_QUERYENDSESSION, autostart opt-in, crash-recovery semantics (interrupted tasks
visible, never auto-rerun), and the update-pending hook consumed by 56.

## Key constraints
- Suspend (PBT_APMSUSPEND): dispose session scope (models released) + stop capture BEFORE
  sleep; resume: reopen audio device, re-enter Sleeping/Armed per config, resync monotonic
  timers; pending L2 approvals at suspend → follow the §16.11#10 lean: void + auto-reject +
  "错过 N 个待确认操作" notice (decision logged; final wording at 54 review).
- Session end (WM_QUERYENDSESSION/WM_ENDSESSION): D38 fast path ≤5s — skip non-critical flushes,
  NEVER skip steps 4/5/9 (audio release, native sessions, Job close).
- Autostart: opt-in registry Run key; default OFF; toggle in tray menu + config.
- Crash recovery: start-after-crash marker → previous tasks marked 中断 in panel/CLI listing;
  never auto-rerun (side effects); diagnostics hint shown.
- Clock resync after resume: monotonic timers unaffected, wall-derived schedules recomputed
  (reminders).
- Update hook: on boot, if `update-pending.json` present → apply-handoff event (56 completes);
  N-1 rollback counters start here.

## Out of scope
- Updater binary itself (56); battery-saver policy internals (28 handles shrink; this ticket
  ensures events reach it).

## Acceptance criteria
- [ ] Real suspend/resume (psshutdown/manual on runner) → session disposed pre-sleep (no models
      resident during sleep — sampler), audio restored post-wake, one resume-e2e task works.
- [ ] Session-end: broadcast WM_QUERYENDSESSION → process exit ≤5s with steps 4/5/9 done (order
      audit log), Job children gone.
- [ ] Suspend with pending L2 → voided + counted notice (D43 note updated with final decision).
- [ ] Crash-recovery fixture: kill -9 mid-task → next boot lists task as 中断; no rerun.
- [ ] Autostart toggle writes/removes Run key; OFF default verified on fresh config.

## Progress log (append-only, newest last)
