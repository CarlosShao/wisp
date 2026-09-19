# 12 — S1 gate: `wisp run` text path end-to-end, notify/list_tools, settle, RSS gate

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 04-sqlite-core, 06-secretstore-envs, 07-ball-state-machine-core,
08-observability-slo, 10-agent-loop-core
**Parallel slots:** ≤2 sub-agents (A: CLI + notify/list_tools tools; B: end-to-end wiring +
SLO verification run)
**Spec refs:** S1 done criteria (§4/D44), D12, D34 (notify, list_tools), D32 SLO, C21 native

## What to build
The S1 vertical slice: `wisp run "..."` CLI (skips voice per D12), `notify` and `list_tools` as
the first two builtin tools, hotkey→(placeholder text capture via CLI for S1)→agent→TTS-less
result via system notification, 3s Settling→unload, and the measured idle SLO gate (tree RSS
≤25MB path-Y / ≤40MB path-X per spike verdict in 02). This is the first demoable product
moment: type → reply → notification → memory returns to idle.

## Key constraints
- CLI: `wisp run "task text"` attaches console, prints streaming text result + final status +
  cost line; exit code reflects error_class. CLI shares the agent loop (no voice, no ball pop
  unless GUI mode).
- `notify` tool (L0, D10 dependency): posts a Windows toast; Focus-Assist degradation contract
  (always also update ball badge) stubbed at ball-level now (real rule at 45).
- `list_tools` (L0): returns resident + third-party tool directory (D15② fallback path).
- Settle: result presentation done → 3s → DisposalScope session teardown (incl. FreeOSMemory) →
  RSS back to idle cap ≤10s.
- SLO measured with 08's sampler on a real Windows desktop session: idle tree-private ≤ spike
  verdict value; goroutines ≤6; handles <300; GDI <200; zero periodic disk writes; zero long
  network connections in idle.
- C21 native-side token table complete (ball + future panel share); tokens doc updated.
- task_log row per CLI task; tool_call rows for notify/list_tools (correlationId present).
- GUI mode end-to-end smoke: hotkey → ball Listening (text captured via CLI in this slice;
  mic arrives S2) → reply → notify → Warm/Settling → Sleeping.

## Out of scope
- Voice anything (13+); panel (33+); gating on non-notify tools (21+); memory (29).

## Acceptance criteria
- [ ] `wisp run "总结一下…"` against mockllm: streamed reply, notification posted, exit 0;
      with fail_next(3) → proper error_class + non-zero exit + user-visible message.
- [ ] Idle SLO gate PASS recorded into docs/SLO.md appendix (numbers + machine + date), using
      tree-private metric; settle ≤10s verified with FreeOSMemory counter >0.
- [ ] Ball state walk Sleeping→Listening→Thinking→Acting→Speaking/notify→Warm→Settling→Sleeping
      matches D43 rows #4,12,15,18,26,31,32 (state log captured).
- [ ] task_log + tool_call rows written with correlation_id; visible via sqlite query.
- [ ] Zero emoji scan over new UI strings passes; tokens doc updated with any additions.
- [ ] Human visual acceptance of ball states (user signs off screenshots — D29 rule).

## Progress log (append-only, newest last)
