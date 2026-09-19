# 45 — Diagnostics bundle + disk-full/WAL guard + Focus-Assist degradation conclusion

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 42-watchdog
**Parallel slots:** ≤2 sub-agents (A: diagnostics bundle + privacy preview; B: disk-full guard +
notification degradation conclusion D42#8)
**Spec refs:** §5.1, §14.4, D42#5#8, D35 artifacts quota, S6

## What to build
The diagnostics export pipeline (user-reviews-before-send privacy contract), the disk-full/WAL
guards, and the final Focus-Assist degradation implementation with its "no public API" honest
conclusion.

## Key constraints
- Diagnostics bundle: rolling logs + redacted config.toml (secrets stripped) + SLO samples +
  versions/platform + DEFERRED registry dump + watchdog snapshots; DEFAULT excludes ASR
  transcripts, audio (never collected anyway), keys; user explicitly opts in per-category with
  a preview of what's included BEFORE export (§14.4); export = single zip to user-chosen path.
- Disk-full guard (D42#5): pre-write free-space check (200MB threshold) on data volume; below →
  stop L3/log/artifacts writes + user warning, main response path UNAFFECTED; WAL kept bounded
  (startup TRUNCATE + autocheckpoint 1000 from 04, runtime re-check here).
- Focus-Assist (D42#8): attempt best-effort state read (documented API surface, expected to
  fail on Win11) → regardless of outcome, tier-3+ results ALWAYS dual-channel (ball badge +
  task list) — the degradation is unconditional; honest PRECHECK note recorded.
- Artifacts quota enforcement point here (LRU at 500MB, 04 provides store); over-quota during
  tier-1/2 result → oldest evicted + logged (visible in diagnostics).
- Log rollover under disk pressure never loses the last error line (flusher coordination).

## Out of scope
- Panel pages rendering these (40); watchdog unload actions (42).

## Acceptance criteria
- [ ] Bundle contents exactly as previewed; seeded secrets/transcripts/audio absent by default;
      opt-in transcript inclusion works with per-item preview.
- [ ] Disk-full fixture (small virtual disk / quota): guard trips at threshold, L3+logs stop,
      live task completes normally; recovery when space freed.
- [ ] WAL bounded under sustained write fixture (size cap observed); checkpoint on demand.
- [ ] Focus-Assist: with toast suppression simulated, badge+task-list still updated (trace).
- [ ] Artifacts LRU eviction under quota pressure; evicted files logged.

## Progress log (append-only, newest last)
