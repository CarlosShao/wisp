# 44 — CostMeter (C23): usage accounting, pricing, budgets, pause-at-limit

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 04-sqlite-core, 10-agent-loop-core
**Parallel slots:** ≤2 sub-agents (A: meter + aggregates; B: budgets + pricing table + pause)
**Spec refs:** SPEC-05 §8, §14.1, C23, D35 cost_daily, S6

## What to build
Per-task and cumulative cost accounting from C6 Usage events: token/cost rollups, day/month
aggregates, budget thresholds with warn-at-80% and pause-new-tasks-at-100%, all persisted to
the cost tables.

## Key constraints
- Meter consumes C6 Usage{in,out,cached} per stream + per-task rollup → task_log columns +
  cost_daily upsert; tool-call counts and duration from tool_call rows.
- Pricing: per-provider/model price table (in-repo JSON, versioned; editable via config
  override) → cost in micro-units (integer math; no float); unknown price → tokens-only +
  "价格未知" marker (never fabricate amounts).
- Budgets `[cost]`: daily/monthly caps; ≥80% → warn toast + badge; 100% → PAUSE new tasks
  (default; `over_budget=warn` alternative); pause behavior explicit + resumable next day;
  running tasks unaffected mid-flight.
- Distinct layers documented: per-task token budget (200k, anti-loop, from 10) vs money budget
  (this ticket) — both enforced independently.
- Cached-token accounting shown (IN/OUT/CACHED labels for 40); cached priced at cache rates.
- All math integer micro-units; reconciliation test vs task_log sums exact.

- **Quota system (2026-09-19 supplement, SPEC-03 §3.1)**: three layers — per-model
  (`quota_daily_micro/quota_monthly_micro`), per-provider (`plan_credit_total_micro` for
  billing=plan: user-entered total, C23 bookkeeping, remaining marked "estimated"; no public
  balance API assumed — provider balance APIs probed later if any), global (`[cost]`).
  **Anti-negative = PRE-DISPATCH check**: before each LLM/streaming dispatch (and each realtime
  turn, ticket 60) C23 gate verifies remaining budget on the role's model → over-quota model is
  SKIPPED via the text_chain / voice_chain (never dispatched into debt); 80% warn, 100% hard-stop
  with a visible fallback event (ball + panel). Quota type aware: plan credit vs pay-per-token
  tracked separately per model.
- Fallback-chain integration: quota exhaustion / 401/403 / provider errors all funnel through the
  same "advance text_chain" mechanism (ticket 11 owns stepping; C23 owns the quota signal).

## Out of scope
- Cost page UI (40); provider price research beyond a starter table (documented, user-editable);
  chain stepping logic (11).

## Acceptance criteria
- [ ] Multi-turn golden task: per-task tokens/cost match hand-computed fixture to the unit.
- [ ] Aggregates: day rollover + month rollover correct (fixture clocks); reconciliation with
      task_log exact.
- [ ] Budget: 80% warn event; 100% → new task rejected with explicit message + running task
      finishes; warn-only mode alternative verified.
- [ ] Unknown-price provider → tokens-only + marker (no invented currency amounts).
- [ ] Quota: per-model daily/monthly gates; plan-credit bookkeeping (manual total − consumption,
      "estimated" marker); **pre-dispatch anti-negative test (fixture: remaining 0 → dispatch
      skipped, chain advanced, visible event); 80%/100% events; plan vs pay-per-token tracked
      separately per model.**
- [ ] Pricing-table version bump → recorded, old aggregates unchanged.

## Progress log (append-only, newest last)
