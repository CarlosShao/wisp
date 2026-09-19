# 49 — Session grants (D45-2): (tool, pattern, session) authorization + input.type binding

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 04-sqlite-core, 48-approval-queue-full
**Parallel slots:** ≤2 sub-agents (A: grant store + evaluation in risk path; B: card third
option + revoke UX + audit logging)
**Spec refs:** SPEC-06 §8, D45-2, §15#7 input binding, C18/C20, S7

## What to build
Scoped session authorization: confirmation cards gain "本会话内允许 <tool> 于 <pattern>"
creating (tool, pattern, session_id) grants honored by the risk path, expiring at session end,
revocable live, fully audited — including the input.type target-process binding variant.

## Key constraints
- Grant = (tool, pattern, session_id) row (approval_grant table from 04); evaluation BEFORE
  re-prompting: matching grant → execute at L1-grade without new confirm; every grant use
  writes a tool_call audit row (grant_id linked).
- input.type variant: pattern = target process name; per-process isolation (Word ≠ 终端);
  "本会话内允许向所有应用注入" is danger-styled + expand-confirm (two clicks); swap of target
  app → re-confirm.
- Three inviolable limits (D45-2): L2 NEVER enters grants (an L2 decision can create only an
  L1-grade grant for future L1-class uses of the same tool+pattern — L2 always re-confirms);
  C25/R4 taint upgrades NOT overridden by any grant; 🔒 config-section loosening never via
  grants.
- Session end (31 dispose) → grants expired (rows retained 30d for audit); mid-session revoke
  (40 page + card) effective immediately on next evaluation.
- Confirm-card UX: three options (仅本次 / 本会话内允许此工具+模式【visually primary】 /
  全部应用 variant for input.type only, danger-styled expand).
- Batch interplay: aggregated L1 confirm may carry one grant creation for the whole batch
  pattern (if user opts in).

## Out of scope
- Cross-session persistent scopes (workspace/user — RESERVED, do not build).

## Acceptance criteria
- [ ] Grant lifecycle: create → subsequent calls silent (audit rows written) → session end
      expires → new session re-prompts.
- [ ] Revoke: immediate effect test (next gated call re-prompts).
- [ ] Limits: L2 never auto-allowed by grant; tainted exfil upgrades despite grant; config
      loosening unaffected by grants (three explicit tests).
- [ ] input.type: per-process isolation; all-apps option danger flow; target-swap re-confirm.
- [ ] Batch + grant: 10-op batch with opt-in grant → one card → whole pattern granted.
- [ ] Audit: grant_id linkage present in tool_call rows for every grant-backed execution.

## Progress log (append-only, newest last)
