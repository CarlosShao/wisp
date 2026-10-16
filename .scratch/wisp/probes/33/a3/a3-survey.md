# 33-a3 只读普查 — 快照"最后一公里"到底缺哪一跳（锚 c315980b，⛔ 判语/选形/排程）

腿 `33-a3` · 角色＝只读普查 · 起手锚＝具名提交 `c315980b`（读数全程 `git show/grep c315980b`，⛔ 读工作树在飞字节）。
尺逐字输出见 `logs/a3-rulers.txt`。本件只量、⛔ 修、⛔ 替编排者排程、⛔ 选修法形状。
目录内 `00-anchor.md`/`verdict.md` 是上一波 33-a3（另一题、锚 `e6c3be17`）已提交旧件，本腿不覆盖、不引用为凭。

## ① 出向：Go 把字节送进 WebView 的每一形

| 形 | 调用点 (file:line) | 生产调用者枚数 | 送的是快照还是回执 |
|---|---|---|---|
| `PostWebMessageAsJson` | 仅测试桩 `internal/panel/composer_dispatch_test.go:461` | **0** | — |
| `AddHostObjectToScript` | 仅注释 `bound_dispatch_door_114_test.go:7/8/240` | **0** | —（页面门走 `Bind`，非此 API）|
| `NavigateToString`(=`SetHtml`) | 定义 `panel_control_windows.go:665`；调用 `panel_host_windows.go:479` 与 `:968` | 2，均 bootstrap/兜底 HTML | ⛔ 非快照（真页面走 `WebResourceRequested` URL 过滤器）|
| `Eval` | 定义 `:667`；调用 `:737/:744/:748/:752` 全在 `msgcb` | 4 | **RPC 回执**：`window._rpc[id].resolve/reject`，回页面主动发的帧 |
| 快照推送（读 `rt.lastSnapBytes` 再送页面）| **不存在** | **0** | — |

快照字节实测落点：`pump.go Publish()` → `run.go:732 Out: rt.bookPanelSnapshot` → `panel_pump.go:360-363` 存 `rt.lastSnap/lastSnapBytes` ＋ 一行 ledger 摘要（`panelSnapshotSummary`）。
唯一读者 `lastPanelSnapshot()`（`panel_pump.go:424`）**生产调用者 0 枚**（`git grep ':!*_test.go' ':!.scratch/**'` 只剩定义行；与票 35 `:93`/`:391` 普查独立同色）。
正控（证明不是"整条没接"）：`publishPanelSnapshot` 生产接线在 `run.go:742`（ui.publish）＋`run.go:1063`（consoleSink），由 `run.go:1363/1452/1481` 在状态真动时调用 ⇒ **Publish() 在生产里跑、字节确被 book**；只缺"从 book 到页面"这一跳。

## ② 页面那边等着什么

- `main.tsx:97` 渲染 `product: node: <App />` **不传 `snapshot` prop**；`main.tsx:110` `createRoot().render(...)` 只挂一次 ⇒ `App.tsx:91 { snapshot = EMPTY }` 落 `EMPTY`（`App.tsx:70` 四键空壳）。
- `lib/panel.ts` 声明的入口全是**页面→Go 出向门**：`window.wispDispatch`（`:155`）、`hostBridge()=window.wispBridge ?? chrome.webview`（`:161`）。**无任何 Go→页面接收器**：`panel.ts` 里 `subscribe|setInterval|onmessage|addEventListener` **0 命中**（正控＝同尺 `window./wisp` 命中 4）。全树唯一 `setInterval` 在 `ai-native/agent-screen.tsx:218`（演示时钟，非快照供数）。
- 意图与现状自陈：`main.tsx:10-14` 注释逐字 "Go pushes snapshots through the C17 bridge … this file renders `<App/>` with its empty frame"。

⇒ 回答：**设计意图＝等推；实现现状＝只用 props 初始值（EMPTY）挂一次，既不拉（无轮询/订阅）也无接推（无接收回调）。** 这一节只报告形状，⛔ 本腿不选修法。

## ③ 桥里缺哪一枚名

`bridge.go:43-75` 入向方法常量 **7 枚**（正控＝同尺 `grep 'Method… ='` 数到 7）：
`panel.mode.request` / `panel.workspace.request` / `panel.attachment.add` / `panel.message.send` / `config.get` / `config.set` / `config.models.list`。
**无 `panel.snapshot.*` 这一族**（`git grep panel\.snapshot` 非测试 0 命中，仅 `composer.go:249` 注释）。
方向说明（事实、非裁）：这 7 枚全是**入向**（页面→Go 请求）；快照是**出向**（Go→页面），`SPEC-08:163/174` 把 `panel.resync`/事件推送标在"请求列＝—"＝出向、本不入入向名册。

若**新增一枚入向方法名**会撞的既有钉（逐枚现量 file:line）：
1. `internal/panel/git_test.go:387`（`whitelistMethodsFromSource` 与四枚 `panel.*` 常量 `equalStrings`；正控形另在 `:517 len!=4`）。⚠ **派单写的 `cmd/wisp/git_test.go` 路径不存在**（`git grep … cmd/wisp/git_test.go` rc=1），真身在 `internal/panel/`。此尺只数带 `panel.` 前缀字面量，`config.*` 不占数。
2. `internal/panel/inbound_roster_253_test.go:63 wantFullInboundRosterSize253 = 7`（断言体 `:449`/`:791`）⇒ 新增入向名让名册涨 ⇒ 当场红。
3. `docs/specs/SPEC-08-ui-ball-panel.md:150/158/163/174` 四处已把出向推送命名入契约（`panel.resync`＋事件推送行），改契约＝人工批准。
（⛔ 本腿不裁该不该新增、⛔ 不排程。）

## 必答题：从 `Publish()` 那串字节到用户屏幕，缺的具体是哪一跳？

**断点＝`cmd/wisp/panel_pump.go:424`。** `bookPanelSnapshot`（`panel_pump.go:360`）把快照存进 `rt.lastSnap/lastSnapBytes` ＋ 一行日志后就**止步**；紧接其下唯一能把这些字节交给宿主的读者 `lastPanelSnapshot()`（`:424`）**生产调用者 0 枚**（`git grep ':!*_test.go'` 只剩定义行）。⇒ 缺的是"**没有一条生产代码调用 `lastPanelSnapshot()` 并经宿主 `c.Eval`/postMessage 把快照推给页面**"这一跳。

⛔ 不说"整条都没接"：射程更窄的事实是——**RPC 回执那一跳是活的**（`msgcb`→`Eval window._rpc.resolve`，`panel_control_windows.go:752`；且 `Run()` 已在 `panel_resident_windows.go:320` 被泵，AC#14 早先"Run/Dispatch 0 调用者"在锚点已过期）。死的只有"出向推快照"这一跳。

⚠ 页面侧还有**第二处未接**（不属上面这枚断点、属第②节）：即便 Go 侧补出推送，`main.tsx` 目前把 `<App/>` 以 EMPTY 单挂、无接收器，`frontend/**` 写权归机主外部 agent（票 35 `:97/:140` 默认 ⓐ）。本腿只并排报出这两处事实，⛔ 选谁先、⛔ 定修法。
