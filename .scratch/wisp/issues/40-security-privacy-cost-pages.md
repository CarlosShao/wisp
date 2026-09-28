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
- [ ] **Privacy/tool-call rows must not read better than reality: when a row carries
      `risk_level == "L0"` and `decision == "allow"`, the rendered detail (and the exported
      bundle's human-readable field) must say it was auto-allowed by tier and that nobody was
      asked - not a bare "allow".** Landing point today is the string built in
      `internal/memory/privacy.go` (ruler: `grep -n "r.Decision" internal/memory/privacy.go`);
      the column value itself must stay verbatim (取证面不许改字, 见票 184 AC#5(a)).
      > **09-28 13:2x 编排者追加（来源＝票 184 的只读普查 `docs/evidence/s1/184-audit-vocabulary-census-c1.md`＋台账 `A365`，不是新规矩，是把已量到的事实挂到 Will-wire-it 那张票上）**：
      > `tool_call.decision` 只有五枚冻结值（`internal/agent/journal.go:26-38`，域被 `internal/memory/models.go:136-142` 机器强制），**没有一枚的意思是"档位即放行、没问过人"**；普查逐枚读了五处记 `allow` 的地方，分档＝**从不弹卡纯机器自批 2 处**（`bridge.go:374`／`loop.go:794`）、**靠沉默放行 1 处**（`bridge.go:386`，L1 窗口超时未否决也记 `allow`）、**真人点头 1 处**（`bridge.go:400`）、**生产到不了 1 处**（`loop.go:803`）。唯一能分别"依据哪份授权"的 `grant_id` **生产零填充**（`journal.go:101` 逐字传 `nil`，且 `_ =` 把错误丢掉），而 D45 的 `InsertGrant` **生产零调用方** ⇒ 填它今天是一具空壳。
      > ⚠ **为什么挂在这张票而不是单开立一张写码票**：普查现量 **`ListPrivacy`／`ExportPrivacy` 在生产里零调用方**（尺：`grep -rn "ListPrivacy\|ExportPrivacy" --include=*.go internal/ cmd/ | grep -v _test` ⇒ 只命中 `internal/memory/privacy.go` 自己）⇒ **今天没有任何一行字被摊给人看过**，单开一枚"改文案"的腿会落在一根没接线的管子上。本票开工接那页时，这一格就是它的硬 AC。
      > ⚠ 加一枚 `decision` 值＝**人工批准面**（要同动 `SPEC-02:79`＋`PLAN.md:2703`）；普查量的代价是**零迁移、新老库不分叉**（该列是无 CHECK 的自由 TEXT），但**收益只有"机器可查询地区分"**，文案那一形就能给人话那一份 ⇒ 编排者 13:2x 裁：**暂不摆 owner**，等本票真需要机读区分时再摆（口令「票 184 要摆 Q」）。

## Progress log (append-only, newest last)

> **2026-09-26 进度追记（owner 把前端交 ZCode 直管后的可视壳进度，只追加不勾框）**：
> 本票的可视面已按 beautiful-ui 主题收编第四代重画，演示形态在 `?harness=1`
> 应用全观（左栏三段式/会话更改卡/右栏审查·终端·浏览器）与 `?harness=2` 陈列室。
> **本票 AC 的真数据 e2e 格仍等票 145/35/33 的字段与宿主，一格未勾、不由本程
> 勾**；可视壳的 commit 链见票 77 进度追记（`424ac84`…`5e23d99`）与
> `docs/reports/frontend-session-log-zcode.md` §1–§7。
