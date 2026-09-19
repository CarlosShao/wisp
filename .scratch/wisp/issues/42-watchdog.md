# 42 — Watchdog: per-state thresholds, safe unloadables, settle enforcement, alerts

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 08-observability-slo, 28-session-scope-warm
**Parallel slots:** 1
**Spec refs:** D18 (valid parts), D32 16.3.2 修订③, D42#10, WatchdogAlert, S6

## What to build
The resource watchdog: state-aware thresholds (query the D32 row for the CURRENT state), a safe
unloadable-component list, `Armed`-never-unloads-KWS rule, settle-failure detection with
WatchdogAlert, and the leak-indicator alerts (GDI/USER/handles/goroutines).

## Key constraints
- Threshold per CURRENT state row (Sleeping/Armed/Warm/Conversation/panel/peak) from the SLO
  table (02-measured values); over-limit sustained 60s → unload pass; 3 consecutive settle
  failures → WatchdogAlert (orange ball + one-click restart).
- **Unloadables (whitelist)**: panel WebView if user not viewing, long-output temp buffers,
  goja VM (if enabled), BM25 index cache. NEVER: KWS in Armed (warn + suggest user disable
  instead — D32 rule), audio capture mid-utterance, SQLite db-writer.
- Leak alerts: handles >2000, GDI >5000, USER objects over cap, goroutines over state cap →
  alert + diagnostic snapshot (not auto-unload of arbitrary things).
- Settle enforcement: after Settling→unload, watchdog re-measures ≤10s; failure counted; also
  verifies FreeOSMemory counter (08) >0 per dispose.
- Watchdog loop also owns config-file mtime polling (05 contract) — one loop, both duties.
- Actions logged with before/after metrics; user-visible alert texts per D37 resource class.
- Watchdog self-protection: runs as named goroutine, its own memory negligible, panic-isolated.

## Out of scope
- Diagnostics bundle export UX (45); power-event triggers (43); slo-check script (08).

## Acceptance criteria
- [ ] Forced-leak fixtures per state: Sleeping over-cap → unloadables dropped (panel hidden)
      → within cap or alert; Armed over-cap → KWS RETAINED + warning (rule-critical test).
- [ ] 3× settle-failure fixture → WatchdogAlert state + restart affordance works.
- [ ] Leak indicators: fixture handle/GDI growth crosses caps → alert + snapshot.
- [ ] All actions logged with metrics; no unload of non-whitelisted components ever (log audit).
- [ ] Config mtime poll fires reload events identical to 05's contract (shared test).

## Progress log (append-only, newest last)
