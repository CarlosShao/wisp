# 40 — Security/Privacy/Cost pages: grants, blacklist status, memory ops, cost summary

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 29-memory-l1l2, 35-panel-bridge-c17, 44-costmeter-c23
**Parallel slots:** ≤2 sub-agents (A: security + privacy pages; B: cost page + models/diagnostics
views)
**Spec refs:** SPEC-08 §5.2, D45-2 registry view, D20/D35 privacy ops, C23, design/screens/
{security,privacy,cost}.html

## What to build
The trust surfaces: Security page (live session grants with revoke, blacklist status, taint-hit
history), Privacy page (profile/memory/task_log/tool_call/artifacts full ops), Cost page (C23
summaries + budgets), plus Models and Diagnostics views.

## Key constraints
- Security page: active session grants (tool, pattern, session) with one-click revoke (49
  lifecycle); B-list overrides listed with their audit rows; taint-hit (R4) history from
  tool_call (source-named); approval-decision audit trail browseable (correlationId joined).
- Privacy page (D35 rule 2 — non-optional): for profile/memory/task_log/tool_call/artifacts —
  list, delete-one, purge-all, export-JSON; transcript rows show masking caveat; export shows
  what's included BEFORE download; confirm dialogs danger-styled.
- Cost page (C23): per-task breakdown (tokens in/out/cached, cost, tools, duration), daily/
  monthly/all-time aggregates, budget bars (80% warn / 100% paused state), pricing-version
  footnote; IN/OUT/CACHED micro labels per §17.5; NO currency emoji, no chart libs (numbers +
  minimal SVG bars).
- Models view: installed models with sizes + delete; download progress reuses Downloading
  state events; verify-signature status display (read-only).
- **Model-management page (2026-09-19 supplement, LLM 接入)**: provider/model catalog CRUD
  (config.toml is the truth source — page is its editor); **probe button** per model with
  "declared ✓ / measured ✗" warnings from `provider_health`; per-model quota editing
  (daily/monthly micro + billing mode plan/pay-per-token + plan credit total); **role assignment
  dropdowns** (chat/memory_extract/handoff/summarize + realtime/cloud_asr/cloud_tts) filtered by
  capability (voice roles list ONLY models with audio_in/audio_out/realtime bits); **text_chain /
  voice chain drag-order editor**; per-provider health/last-error/latency display. All writes go
  through `config.set` (security-section rules N/A here, hot tier).
- Diagnostics view: export bundle button → preview checklist of included items (user reviews
  BEFORE export; §5.1) → produce bundle; "excludes audio/keys/transcripts by default" stated.
- All pages stateless via resync; lists virtualized; empty/error/loading states per
  design/screens/states.html.

## Out of scope
- Grant creation flow (49); cost-meter engine (44 consumes it); watchdog UI states (46).

## Acceptance criteria
- [ ] Every privacy object type: delete-one/purge/export e2e against live DB (counts match).
- [ ] Grants view reflects 49's grants live (integration fixture); revoke takes effect
      immediately on next gated call.
- [ ] Cost numbers reconcile with cost_daily/task_log to the micro-unit (reconciliation test).
- [ ] Budget states: <80 normal, ≥80 warn, ≥100 paused (fixture by editing aggregates).
- [ ] Diagnostics preview lists exactly the bundle contents; exported bundle verified to exclude
      audio/keys/transcript-full.
- [ ] **Model-management: CRUD round-trip into config.toml; probe button → provider_health
      display; capability-filtered role dropdowns (voice roles hide non-voice models); chain
      drag-order persisted; quota edits enforced ranges.**
- [ ] Visual sign-off vs the three design screens; zero-emoji + hex scans green.

## Progress log (append-only, newest last)
