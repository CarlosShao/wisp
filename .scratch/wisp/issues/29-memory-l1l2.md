# 29 — Memory system: L1 profile extraction, L2 explicit memory, feedback loop

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 04-sqlite-core, 10-agent-loop-core, 28-session-scope-warm
**Parallel slots:** ≤2 sub-agents (A: extraction pipeline + prompt injection; B: explicit memory
tools + feedback UX primitives)
**Spec refs:** SPEC-05 §6, D20, C13, D35 profile/memory tables, S4

## What to build
The memory layer users actually feel: async L1 profile extraction (≤20 slots, full-injection,
"记住了：X" flash + one-click undo), L2 explicit memory via `memory.save`/`memory.recall`, and
the data-ops APIs the privacy page (40) will render.

## Key constraints
- L1 extraction: post-task async in the Warm window (NEVER in the response path — asserted by
  latency test); small LLM call via C5; batch-commit profile changes (cache-invalidation
  discipline, §14.3); slot enum per SPEC-02 (pref.language, pref.tone, habit.work_hours,
  fact.family, …); >20 → LRU by updated_at + eviction log; feedback: ball flash "记住了：X" +
  undo (native affordance now; panel polish at 40). Mis-extraction undo must be one interaction.
- L1 injection: full ~400-token profile into prompt suffix (③ section); toggle
  `[memory] l1_enabled`.
- L2 explicit: `memory.save` (L1-risk per D34) / `memory.recall` (L0) keyword+recency ranking
  (last_hit_at/hit_count); NOT auto-injected into prompt; recall results taint-considerate
  (recall output may become sensitive source for 19 — mark it).
- Privacy ops APIs: list/purge-all/delete-one/export for profile + memory (40 renders).
- Extraction failures → silent skip + debug log (never disturb user); extraction quality
  harness: scripted conversations → expected slots (golden set).
- No embedding, no auto-conversation-summaries (REJECTED — do not add).

## Out of scope
- Privacy page UI (40); image/text UX polish of the flash (37/40); cross-session persistent
  grants (RESERVED).

## Acceptance criteria
- [ ] Latency-isolation test: extraction runs post-response; first-token latency unchanged
      (statistical assertion vs baseline).
- [ ] Golden extraction set: 10 scripted conversations → expected profile slots/values;
      batch-commit produces one prompt-cache invalidation per task (not per fact).
- [ ] LRU: 21st fact evicts oldest + log line + flash reflects what was evicted... (flash on
      SAVE; eviction visible in log + panel later) — assert log.
- [ ] Undo: saving then undoing restores prior profile state exactly (row-level diff).
- [ ] memory.save/recall: risk levels, ranking order (recency+hits), no auto-injection
      (prompt snapshot asserts absence).
- [ ] Privacy APIs covered by tests; export format documented (JSON schema in docs).

## Progress log (append-only, newest last)
