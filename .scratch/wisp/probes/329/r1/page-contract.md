# 票 329 `AC#5` — 页面侧接收器合同（Go 半边已落地，这半边交给机主的外部前端 agent）

出处＝编排者裁定 A1009 形甲 ＋ `cmd/wisp/panel_snapshot_deliver_329.go`（本腿落地件）。
本文件只写"页面该怎么接"，⛔ 动 `frontend/**`／`design/**` 一字；本票任何一句都不构成"屏幕上已经会动了"。

## 1. 具名 JS 全局

- 名字：`window.wispOnPanelSnapshot`
- 调用形状：宿主执行一枚 JS 语句——

  `window.wispOnPanelSnapshot(<参数>)`

  `<参数>` 就是泵 marshal 出的那串快照 JSON **原文**（`panelSnapshotEvalJS` 逐字包裹，
  不转义、不截断、不重编码），在 JS 里按对象字面量求值。参数⛔ 是字符串。
- 方向＝**出向推送**（Go 到页面），与 `SPEC-08-ui-ball-panel.md` 里那列"请求＝—"的出向行同族；
  它⛔ 走 `window._rpc` 的 RPC 回执，也⛔ 占用 `frontend/src/lib/panel.ts` 那枚页面到 Go 的
  `wispDispatch` 入向门。入向名册（`internal/panel/bridge.go` 的 `Method*`）一个字没动，
  这个全局名也⛔ 出现在那张名册里。

## 2. 参数形状（就是那串快照 JSON）

顶层键（`internal/panel/composer.go:57-93` 为准，逐字段以 Go 侧 json tag 为权威）：

- `pending`：数组，每条一张待批卡——`correlationId` / `tool` / `args[]` / `level`（`L0|L1|L2|Deny`）/
  `rulesHit[]` / `reason` / `reasonKnown` / `sessionOverrideBlocked` / `callChain[]`（今天恒空＝无生产者，具名缺口）/
  `decidedBy`（恒 `"native"`）。
- `results`：数组，每条形如 `{correlationId, text, done}`（`done=true` 表示该路流已闭合）。
- `composer`：对象，含 `mode`（`current` + `names[]`）、`workspace`（set/unset-with-reason 两态）、
  附件上限等输入行状态。
- `generatedAt`：字符串，这一发成立的那一刻。
- `instructions`：可选键。运行时挂了装载器才出现（票 200 载体），取值五态之一。
- `tasks`：可选键。运行时挂了任务名册才出现（票 197 载体）。

语义＝**整包替换**：每次调用是一枚完整快照，⛔ diff、⛔ 增量合并；接收器要能被调用任意多次并
幂等地把最新一包显示出来。页面若未定义该全局，调用被 WebView2 静默丢弃——这正是今天的事实。

## 3. 触发时机

每次 `wisp` 运行时的面板快照 publish 恰好对应一次调用（泵 marshal → 记账 → 本跳推送）：

- 新确认卡打开（gate Prompt 那一刻，队列里多了一条待批）；
- 卡关闭（被否决 / 被答复进入执行）；
- 任务流闭合（`EvDone`，结果区要停止流式）；
- 其余就是今天既有那几个"状态真动"的驱动点。

★明确⛔ 发生的事：逐 token delta⛔ 推（delta 会随下一发 publish 整包带上）；"多久推一次"的
合并/背压规则仍属票 35 名下（D38d），本跳⛔ 发明任何节流。

## 4. 前提与缺口（两侧都要说全）

- 推送只发生在**本进程拥有面板窗口且窗口已建**的常驻进程里；窗口未开时该发只记账不推送
  （`PanelManager.evalPanelJS` 的 nil-window 分支）。`wisp run` 纯控制台进程＝零推送。
- 页面接收器今天缺位＝现量：`frontend/src/main.tsx:97` 渲染 `<App/>` 不传 snapshot、
  `:110` 只挂一次 ⇒ 落 `EMPTY`；`frontend/src/lib/panel.ts` 声明的全是页面到 Go 的出向门
  （`window.wispDispatch` `:155`、`hostBridge()` `:161`），**没有任何 Go 到页面的接收器/onmessage/订阅**。
  ⇒ 前端 agent 要补的就是定义 `window.wispOnPanelSnapshot`（含初始化时若窗口已建、首包何时到的
  问题——本跳只在 publish 时推，页面自己开窗口不触发补发；冷包语义归票 33/35 的显示面，⛔ 本票范围）。
- 本票交完，端到端仍**不许**被任何人写成"已通"：Go 这半＝"publish 即把 marshal 字节交给宿主"；
  页面那半＝本文件。两半各自成立，合起来才算一屏会动。

## 5. 判据引用（给前端 agent 的验收线索）

Go 半边逐位相等的用例＝`cmd/wisp/panel_snapshot_deliver_329_test.go`
（`TestBookedSnapshotReachesHostBitForBit329`）；页面半边验收＝同一全局名收到同名调用、
参数 `JSON.stringify` 回环后与 Go 侧 `pending/results/composer/generatedAt` 各节一致（前端测试归前端面）。
