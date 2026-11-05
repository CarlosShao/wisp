# 51 — Tier2 goja runtime: C24 host API, wall-clock sandbox, default-off gate

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 17-risk-assessor-c19, 20-host-bridge-fs-tools
**Parallel slots:** ≤2 sub-agents (A: goja VM lifecycle + interrupts + scope; B: C24 host API
surface + contract tests)
**Spec refs:** SPEC-07 §5, D3/D19④, D33/F5, C24, C10, S7
⚠ Tier-2 ships DEFAULT-OFF (`[plugins] tier2_enabled=false`, 🔒).

## What to build
The goja JS runtime for Tier2 plugins: per-call DisposalScope isolation, wall-clock interrupt
(5s default), result cap 256KB, the frozen C24 `GojaHostAPI` surface (only legal host entries),
and the enable gate with L2 re-confirmation.

## Key constraints
- goja has NO memory limit — claiming one anywhere (docs/comments) = fake-promise failure
  (F5). Constraints ONLY: wall-clock 5s via `vm.Interrupt()` from a named watchdog goroutine;
  result ≤256KB (truncate + truncated flag); stack via goja limits; recover boundary on the
  calling goroutine; NEVER on ui-sta/audio threads.
- C24 host API = the ONLY entries visible to JS; initial set (SPEC proposal, finalize + freeze
  here): wisp.http.request(net) · wisp.fs.read|write · wisp.clipboard.read|write ·
  wisp.notify.send · wisp.sysinfo.get · wisp.state.get|set (plugin-scoped KV, capability-free).
  No timers, no fetch, no globals beyond the surface — string-scan test proves it.
- Every call: risk via C19 (declarations input-only), capability checks at bridge, taint flow,
  audit rows. JS throw → `tool`-class error back to LLM (self-correct); goja internal panic →
  goroutine-boundary handling.
- Enable gate: tier2_enabled=false default; flipping ON = 🔒 loosening → L2 re-confirm + reload
  event (05); per-plugin enable keys too.
- VM lifecycle: lazily created per enabled plugin (D3 use-to-wake), DisposalScope teardown at
  plugin unload/session end; one plugin's crash must not affect others or host (recover +
  scope tests).
- spike-02's goja ES conclusions recorded here: **T02 results (2026-09-19, final): async/await+Promise ✓; ES modules ✗; async generators ✗ →
      Tier-2 level = ES2020-minus-modules (module-needing plugins must bundle); Proxy/Reflect ok;
      `vm.Interrupt()` precision P50 +0.4–0.6ms / max ~12ms → wall-clock budget ≥250ms granularity
      safe (5s default fine).** supported language level documented in plugin
  docs stub (full SDK at 57).

## Out of scope
- Command-kind tools (52); SDK/example docs (57); memory-sandbox attempts (banned).

## Acceptance criteria
- [x] Interrupt: infinite-loop JS → stopped at ~5s, tool error, host alive (timing assert).
- [x] Result cap: 1MB string return → 256KB truncation + flag.
- [x] C24 surface: all listed entries work e2e (mock plugin each); undeclared global/host
      entry access → undefined + contract-violation log (scan + runtime double-check).
- [x] Panic isolation: JS triggering goja panic → recovered, other plugin unaffected, host alive.
- [x] Enable gate: default off (no VM exists); enabling without confirm impossible (L2 gate
      e2e); per-plugin disable works.
- [x] No-memory-limit wording audit: repo grep "goja memory limit" → zero.

## Progress log (append-only, newest last)

---
**结案 2026-11-05（A1560）**：非实现者裁决 docs/evidence/s1/51-verdict.md —— AC#1 Interrupt PASS(5s 死循环停下来/host 活)、AC#2 结果帽 PASS(512KB>256KB 截断+flag)、AC#3 C24 九入口 e2e PASS、AC#4 panic 隔离 PASS(宿主 Go panic 每调用一 VM 隔离/归因 C24HostPanic+classify/恰一条 C24-HOSTPANIC/calls 记账/不转JS异常)、AC#5 使能门 PASS(默认关/L2)、AC#6 无内存上限措辞审计实质 PASS(repo grep 命中 4 行全在 51-ac5/51-B 证据件自指、产码 0)。37 顶层用例 PASS/0 FAIL；既有红 Test52SchemaViolations/Test52CommandKind 非本票。⇒ AC#1-6 判过、可结。
