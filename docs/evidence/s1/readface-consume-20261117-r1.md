# readface-consume-20261117-r1 — 前端消费读面波（票 39 AC#1/#4、票 40 AC#3/#4、票 37 AC#5 + render:security40 过期断言修正）

**腿**：agent-109（编排者 20261117 frontend-consume-readface-dispatch.txt，含追加单）
**前置**：`b1d2162dd`（Go 读面 wave 1：sections/errors/reload 事件、CostAggregateRow、ApprovalCardView 五枚 batch 字段）
**本腿 commit**：见 git log（只 commit 不 push）

## 1. 五格落地

### 票 39 AC#1 — sections 全节树＋tier 徽章＋reload 事件重取
- 解析：`frontend/src/lib/config-view.ts`（`ConfigSectionRowView`/`ConfigLoadErrorView` 进 `ConfigView`，`sections?`/`errors?` 双向 omitempty 语义＝缺席＝"腿未答"，不是空文件）；reload 事件 `CONFIG_RELOAD_MARKER = "wisp-config-reload:"`＋`parseConfigReloadEvent()`（refused reply / 非 config.reload 事件 / 非 JSON ⇒ null）。
- 事件总线：`frontend/src/lib/config-actions.ts` `onConfigReloadEvent` / `notifyConfigReloadListeners`；**五个 config.set 写助手**（performConfigWrite / performComposerModelWrite / performComposerThinkingWrite / performPluginToggleWrite / performCredentialWrite）都在 reply 上检测一次并扇出；被拒 reply 零通知。
- 渲染：新件 `frontend/src/components/config-editor.tsx`——`ConfigSectionTables`（每节一张表，key/value 逐字，tier 徽章**逐字** hot/reload/restart/locked；"" ⇒ 「未登记档位」；locked 红 tint）；`ConfigErrorList`（第 N 行；Line=0 ⇒ 「（无行号）」）；`ConfigEditorSection`（useConfigView 自读＋订阅事件整页重读，五态诚实）。
- 挂载：`frontend/src/components/config-screen.tsx` 尾部（外观类脚底；自读不用 snapshot prop，搬去独立 rail 类是页面布局活不是 wire 活）。
- 牙（`frontend/scripts/render-config-screen.tsx`）：喂 18 节 ⇒ 18 张表；摘一节 ⇒ 17 张且旧值零残留；tier 词 hot×12/reload×1/restart×1/locked×4 逐字＋locked 红 tint；未登记档位；errors 第 5 行/（无行号）；reload 事件逐字解析＋被拒零通知（总线级）；config-actions 五助手 notify 计数钉死＝5；编辑器源级钉（sections 缺席≠空文件、errors 缺席≠编造）。

### 票 39 AC#4 — errors 按行号
- 同上 `ConfigErrorList`：line=5 ⇒ 「第 5 行」；Line=0 ⇒ 「（无行号）」，**不折成第 1 行**（Go 侧 loaderrors_456.go 的 0=无行 契约）。

### 票 40 AC#3 — CostAggregateRow 上屏＋cachedKnown 门
- 解析：`frontend/src/lib/cost-view.ts`（`cachedKnown`/`aggregateKnown`/`aggregate` 十字段 `CostAggregateView`；aggregateKnown=true 而块残缺 ⇒ payload-error 具名态，不猜字段；行级 cachedTokens 只在 cachedKnown=true 才算读数）。
- 渲染：`frontend/src/components/cost-settings.tsx`——`CostDayTable` 加 CACHED 列（cachedKnown=false ⇒ 每行「未知」，**结构零不画**）；新 `CostAggregateCard`（月累计 IN/OUT/金额、今日金额、日/月预算——0 ⇒ 「不封顶（0 = 未设预算）」schema 自己的词、计价版本 0/空 ⇒ 「计价版本未知」）；缺面卡改「门与缺」（cachedKnown 恒 false 的构造事实＋aggregateKnown 可选面语义）。

### 票 40 AC#4 — 预算三态条
- `BudgetStateBar`：**只听 aggregate.budgetState 词**（normal/warn/paused 各点亮一段；词表外的词零点亮＋逐字显示＋「条不定格」句）；新任务准入读 `budgetAllowNewTask` 布尔逐字；**页面不复算百分比**（internal/observe 整数基点判定是唯一裁判）。
- 牙（`frontend/scripts/render-security-cost-40.tsx`）：三态各一枚（点亮段＋其余不亮）；「观察中」词 ⇒ 零点亮且逐字；aggregate 十字段喂真 ⇒ 逐字段同值回屏；days-only leg ⇒ aggregateKnown=false 逐字；malformed aggregate ⇒ payload-error；cachedKnown=false ⇒ 「未知」且零 `>0<`；cachedKnown=true ⇒ 77 逐字。

### 票 37 AC#5 — 聚合卡五字段、10 操作一张卡、⛔ 画 allow
- `frontend/src/lib/panel.ts`：`ApprovalCardView` 增五枚可选键（batchTotal/batchDirs/batchFirst/batchDeferred/batchSummary，对应 `internal/panel/approval.go:81-85` 的 omitempty 平铺字段；门＝batchTotal>0，即 Go `HasBatch()`）。
- `frontend/src/components/l2-approval-card.tsx`：新 `BatchBlock`（batchSummary 原句逐字、本批共 N 项、涉及目录、前 N 项路径预览、未展开余数）；非 batch 卡零 batch chrome；**块内零控件**（no-allow invariant 不变）。
- ⚠ 同 commit 删了 `internal/panel/composer_test.go` composerTSGapLedger 的 `"ApprovalCardView"` 登记——ledger 双向规则自己点的名：page 一声明就 stale-red，条目必须同 commit 删，普通双向比对从今接管五键（`TestApprovalCardViewJSONKeysMatchFrontendTypes` rc=0 复跑过）。
- 牙（`frontend/scripts/render-approval-37.tsx` 4b 节）：10 操作 ⇒ **一张卡**（「需要你的确认」×1）＋五字段逐字；普通卡渲染零 batch 痕；batch 卡零 allow 形控件；panel.ts 五键声明钉；ledger 条目消失钉。

## 2. 追加单 — render:security40 过期断言修正
- 事实：`dd52313e1` 把 `security-settings.tsx` grants 假话改真话后，旧断言 `the grants absence must name the nil grant_id writer`（钉 `internal/agent/journal.go`）一直红——仪器把谎言当真理钉住了。
- 修正（`frontend/scripts/render-security-cost-40.tsx` §4）：**正控**＝新真句在源与渲染面同时钉死（生产在写 `gate.go:859`／`grants.go:201/274`／`tools/bridge.go:1445`；面板侧无 grants 列表读面；出口＝console 动词 `grants`/`revoke-grant`（cmd/wisp/approval_reply.go:415/:474））＋渲染面「生产零写入」零出现；**负控**＝把源里的真句剥掉（模拟改回旧假话）⇒ 探测器必须翻红（探测器是函数，正控负控同函数，牙真的咬）。

## 3. 门禁读数（本机 windows，node24 @ D:\work\server\node24，串行）
- `npm run build` rc=0（vite build＋tsc -b）
- `npm run typecheck` rc=0
- `npm run lint` rc=0（40 warning 0 error；新增文件零 warning，l2-approval-card 的 1 条为存量 ToolIcon 形状）
- `render:security40` rc=0；`render:config-screen` rc=0（orchestrator 手动命令 `vite build --ssr scripts/render-config-screen.tsx`＋node）；`render:approval37` rc=0；回归批 render:l2（喂 fixtures/l2-card-fs-delete.json）／composer／nav／shell77／palette38／palette38kbd／roster／receiver／stream／history36／model-settings 全 rc=0
- `go test ./internal/panel` rc=0（21s，含 TestApprovalCardViewJSONKeysMatchFrontendTypes／readface_wave1 五牙／TestTheRendererHoldsExactlyOneDoorToTheHost）
- `tools/d22scan` rc=0（ban #6/#8 frontend 142 文本件全绿；frontend/dist/assets git-ignored，dist 未入库字节，`git ls-files frontend/dist`＝1 枚 .gitkeep）

## 4. 诚实边界（未做 / 不属本腿）
- 票 40 隐私三支按钮（privacy.delete-one/purge/export 走真路由＋确认卡）：读面已落地但**不在本派单五格**，页也未接——留给下一腿。
- 票 37 AC#1 深度 e2e / AC#2 倒计时 / AC#3 原生卡 / AC#6 人眼签收：仍留白（同前批理由）。
- ConfigEditorSection 现挂外观类脚底；是否升独立 rail 类＝页面布局裁定，未动 settings-page.tsx（不在本腿写面）。
- render:config-screen 没有 npm script（历史如此），复跑命令＝`vite build --ssr scripts/render-config-screen.tsx --outDir node_modules/.tmp/configrender --logLevel warn && node node_modules/.tmp/configrender/render-config-screen.js`（cwd frontend）。
- 票面 AC 框本腿不勾（裁决者≠实现者），勾格建议随证据交编排者。
