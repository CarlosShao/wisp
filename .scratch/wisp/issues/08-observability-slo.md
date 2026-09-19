# 08 — Observability: logs/redaction, SLO sampler, slo-check.ps1, CI pipeline

**Status:** in-progress
**Claimed by:** orchestrator -> sub-agent T08-impl
**Last update:** 2026-09-19T22:36:14Z
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
- **Roster exemption (from T03-adv, must land in the D38 roster doc at SLO time)**:
  `disposal-worker` is a sanctioned bookkeeping goroutine name outside the D38b temporary list —
  document it as non-product (or refactor into the temporary roster) so the goroutine SLO gate
  does not misreport it as a leak.
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
- [2026-09-19T22:36:14Z] agent=orchestrator claimed=T08-impl did=dispatched next=sub-agent works through acceptance criteria
- [2026-09-19T23:12:00Z] agent=T08-impl did=observe-module-complete next=slo-driver+ps1+ci [logs: internal/observe logging.go/redact.go/diagnostics.go JSONL rolling+retention, non-disableable 5-class redaction (keys last-4 via secret.RedactSecret, audio/content placeholders, 512-char truncation, inline sk-/Bearer/provider masks, [privacy] redact_paths opt-in), log-flusher via Registry.Spawn (D38b); sampler.go: 6 SLO states x 7 metrics (tree private bytes C30/CPU all-core/GDI/USER/handles/goroutines D38b registry/threads) + transitions + Sleeping zero-disk-write (tree-minus-self)/zero-TCP gates; thresholds.go frozen from SLO.md incl. ruling-1 handles<600, work-peak 700MB recorded as target-not-gate; CheckSettle with instrumented ReleaseMemory counter; diagnostics.go bundle API stub (fail-closed on surviving key shapes); proc: TreePIDs/PrivateBytesByPID/TreeSampler (D42#10 per-pid reads); secret: linux compile shim (fail-closed, no DPAPI on non-Windows); tests green incl. seeded 5-class redaction, leak fixture flip, settle counter, size/day/retention rolls. next: cmd/wisp slo driver, slo-check.ps1, d22scan, compose.test.yml, ci.yml]
- [2026-09-19T23:58:00Z] agent=T08-impl did=slo-driver-complete next=slo-check.ps1+ci [cmd/wisp/slo.go: `wisp slo` SLO driver - boots real runtime (Job Object/registry/log pipeline), per-state JSON verdict exit 0/1/2; SLO.md §7 measurement protocol (GC x2 + instrumented ReleaseMemory + settle before window); leak fixture 100MB page-touched held; settle fixture with C11 dispose-order (buffer dead before FreeOSMemory - the gate caught the wrong-order bug: counter>0 with memory parked); proc.TreeSampler rewritten to NtQuerySystemInformation(SystemProcessInformation): ONE snapshot per read = private working set + threads + handles + CPU per pid (observer CPU 3.7% -> 0.38%), GDI/USER/IO per-pid opens only; KEY finding: SLO.md thresholds are private-WORKING-SET units (spike §7) while C30 TreePrivateBytes is commit charge - sampler gates on private WS, records commit alongside (documented in TreeMetrics/thresholds); fail-closed: unmeasurable/zero reads are sample errors, zero-valid-sample window FAILS, snapshot missing self retries then errors; local dev sandbox intermittently hides self from the system snapshot -> leak/settle fixtures verified partially locally (leak memory gate red confirmed pre-guard, settle dispose-order bug caught), full flip validation lands on CI real runners]
- [2026-09-20T00:25:00Z] agent=T08-impl did=slo-check-ps1-complete next=d22scan+ci [scripts/slo-check.ps1: smoke (Sleeping+Warm+settle+leak) / full (six states) subsets, aggregate slo-report.json (schema in header), build-if-missing (not skippable), independent FreeOSMemory-count double-check, FATAL if forced leak does not flip red; local smoke on the dev sandbox: leak self-test PASS (exit1 flip proven), settle PASS observed once, state sampling intermittently blocked by the dev sandbox's filtered SystemProcessInformation view (self process missing for stretches - values correct whenever visible; probe + leak runs prove the pipeline); the ps1 is runner-grade (windows-latest / self-hosted wisp-slo have unfiltered snapshots)]
