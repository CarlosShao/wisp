# 35 — PanelBridge (C17): method whitelist, correlationId routing, resync statelessness

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 33-panel-host-c27, 34-frontend-scaffold
**Parallel slots:** ≤2 sub-agents (A: bridge transport + dispatch + whitelist; B: resync model +
event stream + frontend client SDK)
**Spec refs:** SPEC-08 §5.2, C17, D31 routing, F2 server-side allow rule, D38(d) backpressure

## What to build
The frontend↔Go bidirectional bridge: `invoke(method,args)→result` with the frozen method
whitelist, correlationId request routing, Go→frontend event push (task deltas, tool chips,
approval requests, ball state, cost ticks) with bounded-queue merge, and `panel.resync` full-
state push making the frontend provably stateless.

## Key constraints
- Method whitelist per SPEC-08 §5.2 table (panel.resync, tasks.list/detail, history.query,
  transcript.get, approval.current/queue, approval.decide, config.get/set, grants.list/revoke,
  privacy.purge/export, cost.summary, models.list/delete, diagnostics.export). Unlisted method →
  reject + log (test with fuzzed names).
- **`approval.decide`: `allow` REJECTED for any panel-sourced call — enforced server-side
  (Go) regardless of UI** (F2 layer 3); `reject` allowed from panel. Method metadata table
  carries required capability + native-authorization flag (consumed by dispatcher).
- correlationId routing: concurrent invokes + out-of-order replies land on the right caller;
  10-way concurrency test.
- Event push: bounded queue per panel; overflow MERGES increments (never drops content);
  backpressure counters visible.
- `panel.resync` on EVERY show: Go pushes full state (tasks, approvals, drafts of current
  results, config view, costs); frontend asserts "no cached state across show" — resync-only
  hydration pattern encoded in the client SDK; stale-state test (hide → mutate → show shows
  fresh).
- Frontend client SDK: typed invoke/events, auto-reconnect on hide/show, error toasts via
  Sonner with badge-fallback rule comment (D42#8).
- Payload hygiene: no secrets/keys in any bridge payload (scan test); long strings truncated
  per log rules.

## Out of scope
- Page implementations (36–40); native card UI (21 already built).

## Acceptance criteria
- [ ] Whitelist fuzz: 50 random/forbidden method names → all rejected + logged; each listed
      method dispatches.
- [ ] Concurrency routing: 10 parallel invokes with shuffled replies → correct pairing.
- [ ] Forged panel allow: direct bridge call `approval.decide{allow:true}` from the webview
      context → server rejects (decision test at Go boundary + browser-context e2e).
- [ ] Resync: hide → tool state changes → show → UI reflects fresh state; SDK forbids cross-show
      caches (architecture test/lint).
- [ ] Backpressure: flood events under blocked consumer → merges, no unbounded memory (heap cap
      asserted), content integrity kept.
- [ ] No secret leakage scan across bridge payloads.

## Progress log (append-only, newest last)
