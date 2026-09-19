# 03 — Project skeleton, runtime rules, error model, goroutine registry

**Status:** in-progress
**Claimed by:** orchestrator -> sub-agent T03-impl
**Last update:** 2026-09-19T06:32:38Z
**Blocked by:** 01-build-chain
**Parallel slots:** ≤2 sub-agents (A: package scaffolding + buildinfo; B: goroutine registry +
error model + clock/instance/Job utilities)
**Spec refs:** SPEC-01 §3–§6, D37, D38, C11, C30, §14.10, D22 bans

## What to build
The full `internal/` package skeleton wired so the binary boots an empty event loop, plus the
cross-cutting runtime rules every later ticket depends on: named-goroutine registry with recover
boundaries, 17-class error model, monotonic-clock discipline, per-session single-instance mutex,
Job Object (C30), DisposalScope (C11), and the WISP_ENV selection logic (values themselves in 06).

## Key constraints
- Packages exactly per SPEC-01 §3 (1:1 with frozen module table §1.2): ball, statemachine, audio,
  speech, agent(+approval,+scheduler), llm, tools, plugin, memory, panel, models, config, secret,
  risk, session, proc, watchdog, observe, buildinfo. No extra modules, none missing.
- **Every goroutine has a name/owner/exit condition**; entry `defer recover()` → error `internal`
  + structured log (name+stack) + cancel owning task ctx (D38b, §14.10). Bare `go func(` is a
  CI-failing pattern (AST scan from ticket 08; keep a local grep check meanwhile).
- **Error model D37**: 17 `error_class` enums (config, auth, network, rate_limit, provider, model,
  audio_device, asr, tool, permission_denied, user_rejected, cancelled, budget, loop, injection,
  internal, resource) with state mapping + retryability; `task_log.error_class` column conforms.
- **DisposalScope C11**: `Defer/Dispose` idempotent, reverse order, per-fn recover, total 3s
  timeout, `disposal_incomplete` visible; mandatory tail: cancel ctx → wait WaitGroup → release
  native session → `debug.FreeOSMemory()` → record peak/tail memory to SLO sampler.
- **C30 JobScope**: all child processes into Job (`KILL_ON_JOB_CLOSE`), `TreePrivateBytes()` API
  (this is THE SLO metric), used for single-instance test and later WebView2 children.
- **Monotonic clock rule** (D42#9): timeouts/retention only via `time.Since`-style monotonic;
  wall clock only for persisted timestamps. Provide a tiny wrapper/linter-note.
- Single instance: named mutex `Local\wisp-single-instance` (per-session; env suffix added in 06);
  second launch activates existing instance's ball and exits.
- Confirm module path `github.com/CarlosShao/wisp` in go.mod (P10 account part is settled by the
  created repo; residual npm/domain/trademark checks live in ticket 56).

## Out of scope
- Any state machine transitions beyond an Idle placeholder; config schema; SQLite; audio; UI.

## Acceptance criteria
- [ ] Binary boots all packages (init-time self-checks pass), exits cleanly through the D38(e)
      10-step shutdown order — with an order-audit test asserting the actual sequence.
- [ ] Goroutine registry test: after boot + one no-op task, count returns to resident baseline ±1.
- [ ] Deliberate panic in a worker goroutine → process survives, error=internal logged with stack,
      root ctx cancelled (test).
- [ ] JobScope: spawn a test child → child killed when parent Job closes; TreePrivateBytes returns
      sane numbers vs Task Manager.
- [ ] DisposalScope tests: reverse order, idempotent, one fn panicking doesn't block others,
      3s timeout marks disposal_incomplete.
- [ ] Second-instance test: launching twice → second exits after signalling first.

## Progress log (append-only, newest last)
- [2026-09-19T06:32:38Z] agent=orchestrator claimed=T03-impl did=dispatched (parallel with T01-adv) next=sub-agent works through acceptance criteria
