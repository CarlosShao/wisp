# 10 — Agent loop core: ReAct, context assembly, budgets, LoopGuard (single-task)

**Status:** review
**Claimed by:** agent-ticket10-fix2
**Last update:** 2026-09-20T08:05Z
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
- D11(3) force_tool/force_chat rule table — DEFERRED(D11-3), owned by the routing ticket, not
  built here (registered explicitly so it cannot evaporate from the acceptance surface).

## Acceptance criteria
- [x] Golden-driven loop tests: pure-text reply; single tool call; parallel tool calls; loop with
      tool result feeding next turn; budget exhaustion → Stuck with explicit message.
- [x] max_tokens golden → all unclosed calls failed (assertion on tool-call statuses).
- [x] Spill tests: 4k-token boundary, 1MB hard cap, artifacts file content + context stub shape,
      context_window scaling (set tiny window → thresholds shrink).
- [x] Compression tests: >12k history → oldest compressed, last-3 raw kept, ids preserved.
- [x] LoopGuard: 3/5/8 ladder injects reminders; 8th → Stuck + visible message; per-tool timeout
      fires.
- [x] Control-layer regex: control words short-circuit without LLM call (mockllm request count
      = 0).
- [x] Prompt-section order test: assembly output has cache-prefix byte-stability across turns
      (④ and BM25 segments last).

## Progress log (append-only, newest last)
- [2026-09-20T02:20Z] agent=agent-ticket10-loop did=claimed next=loop-core+skeleton-tests
- [2026-09-20T03:25Z] agent=agent-ticket10-loop did=loop core+budgets+inject+spill+compress+guard+journal+sink; AC1 golden loop tests green (text/single/parallel/feeds-next/budget-stuck/cancel/roster/task_log rows) next=AC2 truncation+AC3 spill+AC4 compress+AC5 guard ladder+AC6 control+AC7 prompt order
- [2026-09-20T11:30Z] agent=agent-ticket10-loop resumed=adopting prior WIP (uncommitted loop.go+control/guard/truncation tests+disconnect fixture); build ok, test-build broke on control_test.go unused `observe` import; AC2/5/6 tests present (repair forward), AC3 spill + AC4 compress + AC7 prompt tests not yet written next=fix import, verify AC2/5/6, add AC3/AC4/AC7
- [2026-09-20T11:48Z] agent=agent-ticket10-loop did=AC2/AC5/AC6 green; repairs: removed unused observe import (control_test); loop.go run() now operates under root.Ctx so the D11(1) default handler can cancel a running task (was checking the parent ctx → control-cancels-running failed); steering test switched tool-then-text→slow-tool for a real (400ms-tool) steer window (instant-echo race lost the insert) — assertions unchanged/strengthened (requests==2 proves steering adds no round-trip); truncation asserts fail-ALL (complete+cut calls) via tool-call statuses + persisted rows next=AC3 spill+AC4 compress+AC7 prompt
- [2026-09-20T12:05Z] agent=agent-ticket10-loop did=AC3 spill green (spill_test.go): scaling proven behaviourally (same 600B payload spills only under a 4096 window, never under 128k → no hardcoded 4000), 4000-token boundary, artifact file==full output + head/tail/totals/path stub shape, 1MB raw cap + rune-safe cut, E2E through loop via spill-tool fixture (big window no-spill vs tiny window spill, artifact on disk). Also fixed two -race defects surfaced while gating: RunningTask.Wait() now joins the registry handle (was observing Pending()==1 after Wait, D38e) and harness.requests() counts mutex-guarded bodies instead of racing rep.Requests. next=AC4 compress+AC7 prompt
- [2026-09-20T12:22Z] agent=agent-ticket10-loop did=AC4 compression green (compress_test.go): >12k history folds oldest, exactly last-3 (KeepRawRounds) stay verbatim as live tool_use+tool_result; every folded id survives via the never-drop ledger (asserted by id presence, not token count) incl. a model summarizer that omits ids from prose; compress trigger scales (5-round ~500tok history: 128k no-op, 4k folds); under-threshold no-op. next=AC7 prompt order
- [2026-09-20T12:40Z] agent=agent-ticket10-loop did=AC7 prompt-order green (prompt_test.go) → ALL 7 AC families closed; asserts on BYTES: cache prefix (identity/safety/style/resident-tools) byte-identical across two turns with differing time/focus + BM25 + profile + history, volatile strings never leak into the prefix, non-trivial prefix required (no vacuous equality); canonical D39 order with ④ scene mandatory-last; suffix parts profile→BM25→scene. Full suite: 35 tests green across go test -count=2, go vet, and go test -race; status→review (awaiting adversarial acceptance). Notes for reviewer: (a) D11(3) force_tool/force_chat rule table remains DEFERRED(D11-3) — out of this ticket's scope per §Out-of-scope; (b) compression still called synchronously pre-request → DEFERRED(D28-1) (Warm-window hook owns it); (c) a cancelled task's terminal task_log/tool_call rows are written with the root's cancelled ctx and thus abandoned by memory's writer — pre-existing behaviour, spec-silent, NOT changed here (flag if acceptance wants a detached cleanup ctx).
- [2026-09-20T05:40Z] agent=agent-ticket10-fix did=accepted FAIL verdict, starting MAJOR-1/MAJOR-2 + MINOR set next=cancelled-task-terminal-rows
- [2026-09-20T07:10Z] agent=agent-ticket10-fix2 did=resumed after platform kill; predecessor committed fa82b60+c657f85; orchestrator reverted an uncommitted `if true` in clipToTokens — next=MAJOR-1 cancelled-task terminal rows
- [2026-09-20T08:05Z] agent=agent-ticket10-fix2 did=closed the FAIL set, per-defect disposition. MAJOR-1 (cancelled task's terminal rows): already shipped by the predecessor in 1e5d96f — `terminalWriteCtx` = `context.WithTimeout(context.WithoutCancel(ctx), 2s)` (a ctx deadline, no wall-clock subtraction) used by finish()/executeCalls/failOpenCalls; independently verified by mutation: reverting the helper to `return ctx, func(){}` makes TestCancelledTaskPersistsTerminalRows RED with the reviewer's exact reading (state="running", ended_at=NULL, in-flight rows outcome=""/ended_at NULL), restore GREEN. MAJOR-2 (TOOL_TIMEOUT swallowed by the tool's own error): shipped in acbc455 (dispatch's timeout branch returns nil error); mutation-verified: returning `err` makes TestPerToolTimeoutOfContractHonestToolIsToolClass RED (error_class="internal", text "context deadline exceeded") while the pre-existing echo-based TestPerToolTimeoutFires stays GREEN — exactly MAJOR-2's hiding place; restore GREEN. Note for the orchestrator: both MAJORs were already committed, so the work order's "remaining must-fix" list predated commits 1e5d96f/acbc455/7176655. MINOR-1 (7176655 artifact inside the raw cap), MINOR-2 (fa82b60; mutation: unbooking j.decide → RED on decision=""), MINOR-3 (c657f85; real json decode+re-encode with json.Number), MINOR-4 (acbc455 contract-honest timeout coverage): adopted after reading, all green. MINOR-5 (6aa89d2): the predecessor's uncommitted prompt_test.go passes UNCHANGED against the REAL clipToTokens, so the `if true` hack stays reverted; re-applying it does not turn the suite green (it fails, with the reason below), i.e. the budget guard is load-bearing and now pinned. NEW FINDING shipped in 6aa89d2: re-applying that hack HUNG the run inside enforceTotal (prompt_test.go:277) — the trimmer had no progress guarantee, so any context_window small enough that the seven 1-token section floors exceed the scaled total spun Sections() forever and hung the task goroutine in Build() (measured: windows 1/64/128 stall, 256+ return). That hang is the likelier reason the predecessor reached for `if true`. enforceTotal is now strictly decreasing and TestEnforceTotalTerminatesAtEveryWindow covers windows 1..128000, asserting prefix budgets are never reached into and residual overage is floor-forced only. MINOR-6 (8723570): TestToolExecutionRunsFourAcross holds calls open so four-in-flight is a precondition; with the pool forced to 1 the old ceiling test still PASSES and the new one FAILS (0.01s, deterministic). MINOR-7 (d25f139) dead prefixSummary removed. MINOR-8 (cbbbdf8) spill-tool.sse header now says 128k/4000 and 4096/128, matching the test. MINOR-9 (8a41f1a) RecordSink mutex: forcing the interleaving loses 1887/4000 appends and trips -race (so it was never merely latent), GREEN after. Registered elsewhere, deliberately NOT done here: memory's writeQueue.exec doing WithoutCancel before BeginTx is ticket 04's data-integrity call (review §四.4). Out-of-scope gained the DEFERRED(D11-3) line the review asked for. Gates before every commit: go vet ./internal/agent/, go test -count=2, go test -race; tools/d22scan clean. next=independent re-verification (different agent) of MAJOR-1/MAJOR-2 as committed, plus the new enforceTotal hang fix and its termination test
