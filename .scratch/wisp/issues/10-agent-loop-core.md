# 10 — Agent loop core: ReAct, context assembly, budgets, LoopGuard (single-task)

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 05-config-model, 09-llm-provider-openai-mockllm
**Parallel slots:** ≤2 sub-agents (A: loop + tool plumbing + truncation rule; B: context
assembly + spill + compression + cost hooks)
**Spec refs:** SPEC-05 §2, §4, §8, D11, D15, D39, D21, C22, C23(data), S1

## What to build
The `agent` module core: single-task ReAct loop over the C5 seam — session-control regex layer
(D11), context assembly with cache-prefix-ordered prompt sections, two-tier tool injection +
`list_tools`, long-output spill, history compression, token/round budgets, gradient LoopGuard,
and `stopReason=length` → `failToolCallsFromTruncatedMessage`. Tool execution goes through a
`ToolProvider` interface with a trivial echo/test provider (real tools arrive at 20+).

## Key constraints
- Loop: control layer (local regex: 停/取消/重说/大声点/确认 — no LLM, tens of ms) → context
  assembly → provider stream → consume C6 events → tool calls (via provider; real gating at 21)
  → results → loop. Termination: no tool call + text done / cancel / budget / 50-round floor.
- Context assembly (D39, budgets fixed): cache prefix = ① identity(150) + ⑤ safety(100) +
  ⑥ style(50) + ② resident tools(~1200); suffix = ③ L1 profile(400) + ② BM25 top-K ≤5(≤300) +
  ④ time/focus scene(≤100, LAST). Total ≤2300. Provider capability flag for cache breakpoints
  (Anthropic explicit; OpenAI implicit) — interface field now, consumed by 11.
- Tool injection: builtin tools always resident (~1200 tok); third-party via local BM25/keyword
  top-K ≤5, zero LLM round-trips; miss → `list_tools` meta-tool (normal in-loop call).
- Spill (D15③): single tool result >4000 tok → write `artifacts\tool-output-<id>.txt`
  (HOST-INTERNAL write, never a gated Tool — D34 note②), context keeps head 500 + tail 200 +
  total length + path; >1MB raw → truncate + `truncated=true`. ALL thresholds scale by
  provider `context_window` (never hardcoded).
- History >12000 tok → compress oldest rounds, keep last 3 raw; never drop tool_call ids/results
  (C25 chain + truncation rule depend on them). Compression runs OUTSIDE response path (in Warm
  window hook; synchronous fallback for S1 acceptable but flagged).
- LoopGuard C22: duplicate-call detection thresholds [3,5,8] → graduated reminders injected;
  ≥8 → `Stuck` state + explicit user-visible message (never silent stop); per-tool timeoutMs;
  token budget 200k default (per-task, scaled by context_window); 50 rounds = last-resort floor.
- `stopReason=max_tokens` → ALL unclosed tool calls in that message fail (never execute truncated
  args) — D21 must-steal rule; tested at golden level.
- Steering queue (runtime inserts) minimal; root ctx cancellation end-to-end; task completion =
  WaitGroup drained (goroutine budget test).
- Writes task_log row + tool_call rows (04 tables) with error_class/decision fields.

## Out of scope
- Anthropic/Responses adapters (11); approval gating (21) — ToolProvider contract returns
  "unclassified" risk which this ticket routes as L0-pass-through behind a flag; multi-task
  scheduler (47); memory extraction (29).

## Acceptance criteria
- [ ] Golden-driven loop tests: pure-text reply; single tool call; parallel tool calls; loop with
      tool result feeding next turn; budget exhaustion → Stuck with explicit message.
- [ ] max_tokens golden → all unclosed calls failed (assertion on tool-call statuses).
- [ ] Spill tests: 4k-token boundary, 1MB hard cap, artifacts file content + context stub shape,
      context_window scaling (set tiny window → thresholds shrink).
- [ ] Compression tests: >12k history → oldest compressed, last-3 raw kept, ids preserved.
- [ ] LoopGuard: 3/5/8 ladder injects reminders; 8th → Stuck + visible message; per-tool timeout
      fires.
- [ ] Control-layer regex: control words short-circuit without LLM call (mockllm request count
      = 0).
- [ ] Prompt-section order test: assembly output has cache-prefix byte-stability across turns
      (④ and BM25 segments last).

## Progress log (append-only, newest last)
