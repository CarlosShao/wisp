# 46 — S6 acceptance: idle/KWS gates, settle, power pre-mortems, cost verification

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 41-kws-wake-word, 42-watchdog, 43-power-lifecycle-events, 44-costmeter-c23,
45-diagnostics-guards
**Parallel slots:** 1
**Spec refs:** S6 done (D44), D32 gates, §10 #7, SPEC-10 §3/§7

## What to build
The S6 slice gate: measured evidence that the resident app truly sleeps — Sleeping and Armed
budgets on real hardware, settle-return ≤10s across workload types, full power-event pre-mortem
matrix, and CostMeter correctness sign-off.

## Key constraints (evidence, not assertion)
- SLO runs (sampler, real desktop): Sleeping ≤25/40MB (per spike verdict) CPU ≤0.5%, handles
  <300, GDI <200, goroutines ≤6, zero periodic writes, zero long connections; Armed ≤90/110MB
  CPU ≤2% with KWS live; Warm ≤350/400MB.
- Settle: after task types {pure-text, tool-heavy, voice, panel-used, model-download} → tree
  RSS back under state cap ≤10s each (FreeOSMemory counter >0 every time).
- Pre-mortems real-triggered: suspend/resume ×2 cycles with mid-task; battery-saver mode
  (simulated) → Warm 30s/Conversation 15s + KWS off; clock jump +5h mid-session; GDI-leak
  fixture crossing 5000; WebView2 masked (no-panel mode + native L2); mic exclusive-hold;
  disk-full (45) — each with correct D37 classes.
- Cost: 7-day simulated ledger fixture reconciles; budget pause e2e.
- Watchdog: Armed-no-KWS-unload rule exercised; settle-failure → WatchdogAlert drill.
- Adversarial acceptance (different agent): stub scan, SLO-integrity (tree metric), contract
  conformance, DEFERRED cross-check, seven-ban scans.

## Out of scope
- Concurrency (S7); plugin runtime; signing/dist.

## Acceptance criteria
- [ ] All SLO rows measured and archived in docs/SLO.md appendix with machine + date.
- [ ] Five settle scenarios each ≤10s with counters.
- [ ] Seven pre-mortem scenarios executed with evidence + correct error classes.
- [ ] Cost ledger reconciliation green; pause e2e green.
- [ ] Adversarial acceptance report committed; gap audit run (per-slice rule).

## Progress log (append-only, newest last)
