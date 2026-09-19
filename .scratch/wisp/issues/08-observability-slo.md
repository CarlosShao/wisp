# 08 — Observability: logs/redaction, SLO sampler, slo-check.ps1, CI pipeline

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 03-skeleton-runtime-rules, 07-ball-state-machine-core
**Parallel slots:** ≤2 sub-agents (A: observe module + sampler; B: slo-check script + CI workflows)
**Spec refs:** SPEC-10 §3(+3.1), §5.1, D32, §5.1 redaction, D42#10, D22 static bans, SLO gates

## What to build
`observe` module (slog JSONL rolling logs with hard redaction rules; SLO sampler for all six
states; diagnostics-pack API stubbed for 45) and `scripts/slo-check.ps1` (tree-private-memory +
CPU + GDI/User/handles/goroutine/thread sampling per state, settle check) wired into GitHub
Actions PR gates.

## Key constraints
- Logs: JSONL, size+daily rolling, 7-day retention (`[observe]` config), `log-flusher` goroutine.
- Redaction (hard-coded, non-disableable): API keys → last 4; audio buffers never logged;
  `web.fetch` bodies not logged; file contents not logged; long arg strings truncated; file paths
  optionally redactable (`[privacy] redact_paths`).
- Sampler metrics per state (Sleeping/Armed/Warm/Conversation/panel-open/work-peak):
  tree private bytes via C30 `TreePrivateBytes()`, CPU 1-min mean, GDI objects, USER objects,
  handles, goroutines, threads. Also captures state-transition timestamps.
- `scripts/slo-check.ps1`: drives the app into each state (test hooks via env/CLI), samples N
  seconds, emits pass/fail JSON per D32 16.3.2 row (see SPEC-10 §3.1 table, thresholds Y/X per
  spike verdict from 02); settle check: after Settling, tree memory back under state cap ≤10s and
  asserts `FreeOSMemory` was invoked (instrumentation counter).
- CI workflows: lint (gofumpt/vet/staticcheck + D22 seven-ban grep/AST scans + emoji scan over
  design/+frontend/), test-core (ubuntu+compose), test-windows (junction + cgo smoke),
  slo-smoke (windows-latest memory/handle subset). Full audio/WebView SLO = self-hosted/local
  pre-merge (documented; script identical). No job may be skippable (D22 mode-6 ban).
- Docker: bring up `docker/compose.test.yml` skeleton with mock-llm service stub if 09 not yet
  merged (else reuse).

## Out of scope
- Watchdog reactions to metrics (42); diagnostics bundle export UI/API (45); CER latency gates (16).

## Acceptance criteria
- [ ] Sampler reports six states with all 7 metrics; Sleeping shows zero periodic disk writes +
      no long-lived network connections (assertions from process telemetry).
- [ ] slo-check.ps1 full run on Windows → JSON pass/fail matching the state table; forced leak
      fixture (allocate 100MB) flips to fail; settle check fails if FreeOSMemory counter is zero.
- [ ] CI green on PR: all four jobs; introducing a bare `go func(` or emoji fails lint.
- [ ] Log redaction tests: seeded key/audio/fetch-body/long-arg → none appear in rolling logs.
- [ ] Diagnostics: sampler data attachable to a bundle (45 completes the UX).

## Progress log (append-only, newest last)
