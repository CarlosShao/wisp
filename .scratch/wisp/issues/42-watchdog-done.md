# 42 — Watchdog: per-state thresholds, safe unloadables, settle enforcement, alerts

**Status:** **DONE-带残余**（2026-10-31 编排者收口：五格全翻；残余=宿主接线票＋R-42r1-1 小票＋225 摘标记） ~~ready-for-agen~~
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 08-observability-slo, 28-session-scope-warm
**Parallel slots:** 1
**Spec refs:** D18 (valid parts), D32 16.3.2 修订③, D42#10, WatchdogAlert, S6

## What to build
The resource watchdog: state-aware thresholds (query the D32 row for the CURRENT state), a safe
unloadable-component list, `Armed`-never-unloads-KWS rule, settle-failure detection with
WatchdogAlert, and the leak-indicator alerts (GDI/USER/handles/goroutines).

## Key constraints
- Threshold per CURRENT state row (Sleeping/Armed/Warm/Conversation/panel/peak) from the SLO
  table (02-measured values); over-limit sustained 60s → unload pass; 3 consecutive settle
  failures → WatchdogAlert (orange ball + one-click restart).
- **Unloadables (whitelist)**: panel WebView if user not viewing, long-output temp buffers,
  goja VM (if enabled), BM25 index cache. NEVER: KWS in Armed (warn + suggest user disable
  instead — D32 rule), audio capture mid-utterance, SQLite db-writer.
- Leak alerts: handles >2000, GDI >5000, USER objects over cap, goroutines over state cap →
  alert + diagnostic snapshot (not auto-unload of arbitrary things).
- Settle enforcement: after Settling→unload, watchdog re-measures ≤10s; failure counted; also
  verifies FreeOSMemory counter (08) >0 per dispose.
- Watchdog loop also owns config-file mtime polling (05 contract) — one loop, both duties.
- Actions logged with before/after metrics; user-visible alert texts per D37 resource class.
- Watchdog self-protection: runs as named goroutine, its own memory negligible, panic-isolated.

## Out of scope
- Diagnostics bundle export UX (45); power-event triggers (43); slo-check script (08).

## Acceptance criteria
- [x] Forced-leak fixtures per state: Sleeping over-cap → unloadables dropped (panel hidden)
      → within cap or alert; Armed over-cap → KWS RETAINED + warning (rule-critical test).
- [x] 3× settle-failure fixture → WatchdogAlert state + restart affordance works.
- [x] Leak indicators: fixture handle/GDI growth crosses caps → alert + snapshot.
- [x] All actions logged with metrics; no unload of non-whitelisted components ever (log audit).
- [x] Config mtime poll fires reload events identical to 05's contract (shared test).

## Progress log (append-only, newest last)

## ★★★终裁落地（2026-10-31 11:4x 编排者，凭非实现者 `42-v1`，件 `docs/evidence/s1/42-v1-verdict.md` 112 行、`934daf6e5`）

★**五格全翻、可带残余结案**（★`42-v1` 自跑：watchdog 12 PASS／observe 82 PASS；★**mutation A 复刻恰 7 枚红、红句与实现件逐字对上**；★**反例探针自加**——注册表只装非白名单 ⇒ 零 Unload 调用＝**白名单∩注册表结构保证成立**，换突变 B ⇒ 恰 4 红＝探针有牙）。
★**D43 核表**：#9＝`table.go:69-72`（Armed 留态＋`kws.keep-alive-alert`、实现的结构性不卸 KWS＋去重告警一致）；#35＝`table.go:256-257`（`WatchdogAlert`＋`ui.one-click-restart`、settle 真机 Dispatch 落 sink 一致）。
★**残余（具名、不阻结案）**：①**宿主接线新票**（真 TreeReader/StateFn/OnAlert→D43 映射/四组件注册＋Start；★**OnAlert 依赖票 273 状态机 sink** ⇒ 与 273 同波或其后）②**`R-42r1-1` 新小票**（★**`382-v1` 同款已判"现在开"**——`manager.go:160-168` 失败支无条件收养 stat ⇒ 瞬时 IO 读失败永久隐藏该次写；修法「解析类收养/IO 类不收养」今天可行：`loader.go:81` read 支 vs `:102` parse 支错误链可分）③doc.go DEFERRED 摘标记（225 写面）④USER record-only/CPU 无卸载＝如实登记非缺口。
