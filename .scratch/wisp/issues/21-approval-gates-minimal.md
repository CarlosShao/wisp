# 21 — Approval gates minimal: L1 block window, L2 native card, queue trivial, batch D45-1

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 17-risk-assessor-c19
**Parallel slots:** ≤2 sub-agents (A: L1 window + veto channels; B: L2 native card + queue +
batch aggregation)
**Spec refs:** SPEC-06 §2, §7, §8, D4, D31, D45-1, B1, C18 (trivial impl), 16.9#5

## What to build
Single-task approval mechanics with the simplest viable UI (native, pre-panel): L1 pre-execution
block window (2–3s countdown strip) with its four veto channels, L2 native confirm card
(full params + rules_hit + reason; allow ONLY from native side), a single-task ApprovalQueue
(C18 semantics at trivial scale), and batch aggregation for homogeneous L1 operations.

## Key constraints
- L1 window (B1-corrected): 2–3s pre-execution BLOCK (not "undo"); countdown ring on ball strip;
  veto = click ball / Esc hotkey / panel reject (no-op until 37) / KWS veto word (channel
  present, functional at 41; when KWS not loaded UI MUST say 「语音取消不可用」 — never fake it).
  Timeout → execute; veto → call cancelled, failure returned to LLM.
- L2: native confirm card (C27 fallback-grade): tool name + L2 badge + FULL params (not summary)
  + rules_hit reasons + reject button + "click ball to approve" guide. **Allow action accepted
  ONLY from native sources: ball click / native card button / global hotkey. Panel-sourced allow
  is structurally rejected at the server-side API (F2 layer 3) — implement the source check
  NOW even before panel exists.**
- Queue (single-task trivial): one pending item; correlation_id assigned; **timeout 300s →
  auto-REJECT** (never infinite wait); last-30s prominent warning; reject → task root ctx NOT
  cancelled; one-click replay available. Host-unreachable/queue-broken → fail-closed.
- Batch aggregation (D45-1): ≥3 homogeneous L1 ops in ONE tool call → single confirm (total +
  affected-dirs summary + first 5 + expandable); L2 NEVER aggregated; R7 (≥50 files) upgrades
  the batch to L2 (unaggregated).
- Cancel non-atomicity surfaced: cancel/veto after start → applied-steps report (from 20).
- Ball visuals: Confirming pulse + AwaitingApproval depth badge (depth=1 now).
- **D47 scope note: all approval mechanics live on the TEXT loop (Path T + handoff). Path C
  realtime brain has ZERO tool permissions — it cannot produce approval requests; if the user
  asks it to "do X", the C32 handoff returns the task to the text loop where full gating applies.**

## Out of scope
- Multi-task queue routing (48); panel-side card (37); session grants (49); input.type binding
  (23 + 49).

## Acceptance criteria
- [ ] L1: countdown executes on timeout; each working veto channel cancels; non-loaded-KWS
      message asserted; veto returns failure to LLM (loop continues).
- [ ] L2: card shows full params + rules_hit; native allow executes; simulated panel-source
      allow REJECTED at API (test with forged source); 300s timeout → reject + warning at 270s
      + replay works.
- [ ] Batch: 10 L1 ops → one confirm; L2 mixed in → no aggregation; 50+ → L2 (R7).
- [ ] Cancel-mid-execution → applied-steps report rendered.
- [ ] State machine: Confirming/AwaitingApproval transitions match D43 #17,21–24.

## Progress log (append-only, newest last)
