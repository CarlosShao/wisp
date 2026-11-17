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
- [x] Cost numbers reconcile with cost_daily/task_log to the micro-unit (reconciliation test).
- [x] Budget states: <80 normal, ≥80 warn, ≥100 paused (fixture by editing aggregates).
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

---

## 编排者追加事实（2026-11-16 20:3x，来自只读审计腿 `census-unwired`，件 `docs/evidence/s1/census-unwired-20261116-2017.md`）

**这一段是"现读坐标"，比票面正文新；按它做，⛔ 按票面旧行号。**
- **成本格是"假免接线"**：`frontend/src/components/settings-page.tsx:114` 说"没收到用量数据"，但 **Go 侧读口早接好了**——`cmd/wisp/panel_cost_194.go:51` 读真表 ＋ `panel_inbound.go:354` 装配 ＋ `internal/panel/cost_handlers.go:120` 出 `wisp-cost-summary:` payload 行（票 194 七格全勾 -done），**而 `frontend/` 里 `cost.summary` 零命中** ⇒ **数据在门口没人取**。本票必须把它取上屏（行为牙：喂真 payload ⇒ 逐字段上屏；喂不到 ⇒ 显式"未判"）。
- **安全格同样是"假免接线"**：`settings-page.tsx:104` 说"没收到策略快照"，但档位**每包都推**（`cmd/wisp/run.go:1187`→`internal/panel/pump.go:372-373`）⇒ 接上即可。
- **`settings-page.tsx` 12 格里 9 格是 waiting 格**（左栏 9 枚灰点＋「（还没有数据）」aria 名）——这就是机主说的"一堆未接线"的本体。本票逐格判：**有数据的接上、真没有的显式说缺哪台引擎**，⛔ 保持一个含糊的"还没有数据"。
- **MCP 格不接是架构裁定的**（D13），⛔ 顺手接；但**文案要写明"按设计不接"**而不是让它看起来像忘了接。

---

## 进度追记（2026-11-17 09:3x，前端写腿 agent-101，票 40 本轮；只追加不勾框）

**本轮交付（前端面，全部 `frontend/**` 写面；未动 cmd/wisp 一字节）：**

1. **成本格接线（census ①的"假免接线"清账）**：`frontend/src/lib/panel.ts` 新增
   `requestCostSummary()`（`cost.summary` 页面字面量落在唯一门文件 panel.ts）；
   `frontend/src/lib/cost-view.ts` 新增 `wisp-cost-summary:` 解析器＋`useCostSummary()` hook
   （五态：loading/no-host/ready/unreadable/payload-error，拒绝句逐字上屏不解释）；
   `frontend/src/components/cost-settings.tsx` 新组件：按天表（day/IN/OUT/costMicros/tasks
   逐字上屏，IN/OUT 微标签，行序照读口），残缺行**计数丢弃不补零**；页面不派生任何聚合。
   行为牙＝`frontend/scripts/render-security-cost-40.tsx`（npm script `render:security40`，
   rc=0）：喂真 payload ⇒ 逐字段同值回屏；refusal ⇒ 原话逐字；payloadError ⇒ 具名态。
2. **安全格接线（census ②的"假免接线"清账）**：`frontend/src/components/security-settings.tsx`
   新组件，读 snapshot `composer.mode`（每包都推）：当前档位（unknown ⇒ 档位未知＋"本泵没读到
   档位读数"，⛔ 借最安全档名）、词表（risk.ModeNames 逐字）、L2 强确认切换名单逐字；**本页不画
   切档控件**（切档只走 composer 行的 panel.mode.request 原生 L2 通路）。五枚缺引擎面（阈值/
   会话授权撤销/B-list/taint 历史/审批审计）逐枚具名"未接线"＋缺哪台引擎（含
   internal/risk/thresholds.go、internal/agent/journal.go）。
3. **技能格接线（census 判"技能半假"清账）**：`plugins-screen.tsx` 抽出 `PluginInventoryList`
   （加号浮层与设置页共挂同一份清单，DOM 不复制），`frontend/src/components/skills-settings.tsx`
   新组件挂 snapshot `plugins` 节：行（id/name/version/state/能力/网络白名单/宿主面）＋
   engineState/skillState **逐字**；开关走 config.set 的 plugin_enabled（同一设置门）；无读取器
   （键缺）⇒ 具名 PumpSources.Plugins 缺失句，不是空清单。
4. **设置页左栏 12 格逐格重判**（`harness/settings-page.tsx`）：安全/成本/技能三格转实控件；
   通用/记忆/MCP/隐私/引导/悬浮球六格 waiting 文案逐格改写——每格写明缺哪台引擎（MCP 格明确
   写"按设计不接（D13）"）；aria 后缀「（还没有数据）」改「（未接线）」（census 点名的含糊本体）。
5. **TIER_NAMES 收拢**：`frontend/src/lib/tier-names.ts`（一处拼写，App 档位控件与安全格共用）。

**Go 侧需要（⛔ 本腿不造，留给编排者派 cmd/wisp／internal/panel 腿；file:line 现读）：**
- 成本聚合/cached/预算/计价版本：字段加在 `internal/panel/cost_handlers.go:30-36`（CostDayRow）
  与 `:124-127`（costPayload），页面读取随后补——页面侧已具名待接。
- 隐私 D35 四操作＋config 读口：`internal/panel/bridge.go:41-107` 名册加 privacy 路由 +
  `internal/panel/composer_dispatch.go` 派发表加 case + `cmd/wisp` 装配（接
  `internal/memory/privacy.go` 的 ListPrivacy/ExportPrivacy——生产零调用方，census A365）；
  工具调用行 tier 自批的人话文案（AC 最后一格）在该读面落地时一并做（列值保持逐字）。
- 会话授权（D45）读面＋撤销写面、B-list 审计行、taint-hit 历史、审批审计 correlationId 关联：
  全部无快照节无路由，安全格已逐枚具名。
- 诊断导出束、模型管理 CRUD 的探针/角色下拉/拖序（票 40 其余 AC）：无 Go 读/写面，未开工。

**Go 侧唯一改动（本腿，同 commit 铁则）**：`internal/panel/bridge_test.go` 的
`rosterPageGapLedger194` 删除 `cost.summary` 条目——页面字面量已落 panel.ts，双向审计规定
"页面开始携带 ⇒ 同 commit 删条目"；`go test ./internal/panel` 三把镜像尺
（TestEveryWhitelistedInboundNameIsAlsoSentByThePage／TestFrontendComposerRequestsMatchTheEnvelope
／TestTheRendererHoldsExactlyOneDoorToTheHost）rc=0。

**AC 格全部未勾**，理由：真数据 e2e（live DB 对账、budget fixture、诊断束验证）要上面的
Go 读/写面先落地；可视对照（design 三屏）不在本腿射程。

**门禁读数**：`npm run build` rc=0；`typecheck` rc=0；`lint` rc=0（38 条既有 warning，0 error）；
`render:security40` rc=0；`render:nav` / `render:palette38` / `render:composer` / `render:shell77`
全部 rc=0；dist 已删回 `.gitkeep`；⛔ 未 push；未动 thresholds/golden/allowlist/PLAN/specs。

---

## 编排者收件·第一批（2026-11-17 09:2x，`2f173996e`，仪器 `render:security40`）

**尺1** `2f173996e` 在。**尺3 我自己复跑**：`npm run render:security40` **rc=0**（`…consume the real feeds (cost.summary payload verbatim …)`）；build/typecheck/lint 全 rc=0。
**census 点名的三处假免接线已清账（我 grep 复核）**：成本 `cost-settings.tsx:11` 真发 `cost.summary`（`lib/panel.ts requestCostSummary`）；安全 `security-settings.tsx:83` 真读 `composer.mode`；技能复用 `PluginInventoryList`（`skills-settings.tsx:27`）。
**0/6 格勾（诚实）**：AC#1 隐私对象 delete/purge/export 对活 DB —— 缺 privacy 路由（`bridge.go:41-107` 八法无它＋`composer_dispatch.go` 派发表无它）；AC#2 grants 视图随 49 活（49 在跑）；AC#3 成本**对账到微单位**缺聚合字段（`cost_handlers.go:30-36/:124-127` CACHED/月累计/预算/计价版本 4 枚）；AC#4 预算三态同缺；AC#5 诊断包清单缺 Go 导出面；AC#6 视觉签收（机主）。

---

## 编排者追加·读面已落地（2026-11-17 10:3x，`b1d2162dd`）

我复跑（载具一发）：`Test40CostReconcileAgreesToTheMicroUnit` 0.20s／`Test40BudgetThreeStatesFromTheAdapter` 0.57s（4 子例 normal/warn/paused/warn）／`Test40CostSummaryPayloadCarriesAggregateOnChain` 0.16s／`Test40PrivacyDeleteOneDeletesExactlyOneRow` 0.19s／`Test40PrivacyPurgeClearsTheDomainAndCounts` 0.18s／`Test40PrivacyExportCarriesTheDomainJSON` 0.14s／`Test40PrivacyWiringNailInProductionChain` 0.19s／`Test40PrivacyRefusalsNameTheDomainVocabulary` 0.20s ⇒ **8/8 PASS rc=0**（另有内部三包绿）。
- **AC#3/AC#4**：`internal/panel/cost_handlers.go:54 CostAggregateRow`（月累计/预算/三态/计价版本）＋`:92 CostAggregateReader`；`cmd/wisp/panel_cost_194.go:94`（月累计＝cost_daily 真行求和；预算走 `observe.BudgetFrom/Judge` 整数基点）＋`:159 budgetStateWord`。**CACHED 是话不是读数**：`CostDayRow.CachedTokens` 存在但 cost_daily 无 cached 列 ⇒ payload `cachedKnown=false`，⛔ 画 0。
- **AC#1**：名册 `internal/panel/bridge.go:141`（`privacy.delete-one`/`privacy.purge`/`privacy.export`）＋派发表 `composer_dispatch.go:174`＋处理器 `cmd/wisp/panel_privacy_40.go:89`＋接线 `panel_inbound.go:411`。三支计数对真 DB。
⇒ 前端剩：预算条读 aggregate（100→paused/80→warn）、隐私三支按钮走真路由（删除/清空/导出，需确认卡）。

---

## Progress log（续，2026-11-17 11:4x，本票前端消费腿 agent-109 · 读面消费波；只追加不勾框）

**本轮交付（frontend/**；证据件 `docs/evidence/s1/readface-consume-20261117-r1.md`）：**

1. **AC#3 页半落地**：`frontend/src/lib/cost-view.ts` 解析 `cachedKnown`/`aggregateKnown`/`aggregate`（CostAggregateRow 十字段；aggregateKnown=true 而块残缺 ⇒ payload-error 具名态，不猜字段）；`frontend/src/components/cost-settings.tsx`：CACHED 列上屏但听 cachedKnown 门（false ⇒ 每行「未知」，**结构零不画**）；新 `CostAggregateCard`（月累计 IN/OUT/金额、今日金额、日/月预算——0 ⇒ 「不封顶（0 = 未设预算）」、计价版本 0/空 ⇒ 「计价版本未知」）。
2. **AC#4 页半落地**：新 `BudgetStateBar`——**只听 aggregate.budgetState 词**点亮三态段（normal/warn/paused），词表外的词零点亮＋逐字显示；新任务准入读 budgetAllowNewTask 布尔逐字；**页面不复算百分比**（internal/observe 整数基点判定是唯一裁判，⛔ 第二裁判）。
3. **牙**：`render-security-cost-40.tsx` 重钉——aggregate 十字段喂真⇒逐字段同值回屏；三态各一枚＋词表外词一枚；days-only leg ⇒ aggregateKnown=false 逐字；malformed aggregate ⇒ payload-error；cachedKnown=false ⇒ 「未知」且零 >0<；cachedKnown=true ⇒ 数字逐字。
4. **追加单同轮落地**：修该仪器 §4 的过期 grants 断言（dd52313e1 改真话后旧断言钉着 journal.go 假话一直红）——改钉新真话（生产在写 gate.go:859／grants.go:201/274／tools/bridge.go:1445、面板侧无 grants 列表读面、出口＝console 动词 grants/revoke-grant），**正控**（真句在源＋渲染面同时钉）＋**负控**（剥掉真句模拟改回旧假话⇒探测器必须翻红），断言一枚没删。
5. **[ ] AC#3 / AC#4 勾格留白（裁决者≠实现者）**：页半＋牙已齐，等对抗验收；对账到微单位的真数据 e2e 仍要 live 宿主。

**门禁读数**：build/typecheck/lint rc=0；`render:security40` rc=0；`go test ./internal/panel` rc=0；d22scan rc=0；⛔ 未 push；dist 未入库。
**诚实边界**：隐私三支按钮（privacy.delete-one/purge/export）不在本派单五格，页未接——留给下一腿（路由已落地：bridge.go:141 名册＋composer_dispatch.go:174＋panel_privacy_40.go:89）。

## 编排者收件·第二批（2026-11-17 11:4x，`f392f34a2`）

**尺1** `f392f34a2` 在。**尺3 我自己复跑**：`render:security40` rc=0（成本十字段逐字、`cachedKnown=false` ⇒ CACHED「未知」零 `>0<`）；载具 `Test40CostReconcileAgreesToTheMicroUnit` PASS 0.20s／`Test40BudgetThreeStatesFromTheAdapter` PASS 0.57s（4 子例）；build/typecheck/lint rc=0。
- **勾 AC#3**：`CostAggregateCard` 十字段逐字上屏（月累计 IN 2100/OUT 460/金额 900、日/月预算 0⇒「不封顶」、计价版本 3＋2026-11-01），**对账牙** `1000/1000 Agree=true`、`+1micro ⇒ false 点名 1001/1000`。
- **勾 AC#4**：`BudgetStateBar` 三态各一枚（normal 绿/warn 橙/paused 红，**只听 state 词**、⛔ 不复算百分比）；Go 四子例 700→normal／800→warn／1000+pause→paused／1000+warn→warn 现跑 PASS；词表外词「观察中」⇒ 零点亮＋逐字显示。
- 留白：AC#1 隐私三按钮（路由已备，本单未接 UI）、AC#2 grants 视图（面板无入站读面，已具名）、AC#5 诊断包导出面、mod 管理行、AC#7 视觉签收（机主）。

## Progress log（续，2026-11-17 11:5x，本票隐私页腿 agent-111 · AC#1 页半；只追加不勾框）

**本轮交付（AC#1 隐私三支按钮页半；证据件 `docs/evidence/s1/40-privacy-ui-20261117.md`）：**

1. **三支真按钮落地**：`frontend/src/lib/panel.ts` 新增 `requestPrivacyDeleteOne/Purge/Export`（选择器键照 bridge.go:211-212 JSON tag：privacyDomain/privacyItemId，走既有 settings 门，无第二通道）；`frontend/src/lib/privacy-view.ts` 新解析器（拒绝对照 RefusedEnvelopeForUser 句首分类；export 解析 `wisp-privacy-export:` 载荷行，records=N＝载荷 records 数组长度，envelope 的 count 键故意不读——internal/memory 不填它，第二个数＝第二个裁判；无载荷⇒未判，⛔假 0）；`frontend/src/components/privacy-settings.tsx` 新组件：删除一条（域词＋id 两输入，id 空⇒按钮禁用＝需选 id，页面写明名册无 privacy.list 读面不假造行选器）、清空某域（两步确认卡 role=alertdialog，purge 请求源级单发点只在确认键上，⛔一键清空；门上无「全清」参数页面不发明）、导出某域（prose 逐字＋records=N＋JSON 原行上屏）。域词不复制词表：PrivacyDomains() 拥有词表，未知域的拒绝句点名五域逐字上屏。
2. **设置页隐私格转实控件**（harness/settings-page.tsx `control: "privacy"`）：旧 A365 零调用方 waiting 文案随接线作废（三扇门现在有真调用方——面板这三支）。
3. **仪器**：新 `render:privacy40`（frontend/scripts/render-privacy-40.tsx，rc=0）——喂真 payload⇒三支逐字上屏（delete 点名 id=42／purge 报 store 自报计数 3／export records=2＋行内数据原样）；喂拒绝⇒逐字（含 Test40PrivacyRefusalsNameTheDomainVocabulary 钉的域词句）；喂不到⇒显式未判，records=0 全绿检查；源级牙：三枚字面量落 panel.ts＋台账条目同 commit 已删＋purge 单发点＋delete 需 id＋页面不复制域词表＋零 emoji。`render:security40` 两枚隐私格旧断言改钉新真话（断言一枚没删）。
4. **Go 侧唯一改动（同 commit 双向审计）**：`internal/panel/bridge_test.go` 的 `rosterPageGapLedger194` 删三条 privacy 注册（页面开始携带⇒同 commit 删条目，与 cost.summary 同例；plant 正控保留）；`bridge.go`/`inbound_roster_253_test.go` 相应注释同步（无行为改动；名册计数 11/7 不动）。`go test ./internal/panel` rc=0。

**门禁读数**：build/typecheck/lint rc=0（40 条既有 warning，0 error，新件零命中）；`render:privacy40`/`render:security40` rc=0；`go test ./internal/panel` rc=0；`go build ./...` rc=0；d22scan clean；dist 已删回 `.gitkeep`；⛔ 未 push。
**AC 格未勾**（裁决者≠实现者）；诚实边界：面板无隐私清单读面（id/域词手填，缺面具名）、录音留存/遥测仍无读口不画。
