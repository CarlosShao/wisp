# `305-a1` · `40` 表④：代价表＋必答（动②撞⛔ 撞票 33 那条"专用 STA ＋ 库 `Run()` 泵"）

## A. 调用点尺先跑（派单要求：尺名＋射程必写）

**尺 1（生产调用者）**＝`grep -rn -E "\.Eval\(" --include=*.go cmd internal | grep -v _test.go`，射程＝工作树 `cmd/`＋`internal/`（这两目录 porcelain 0 行 ⇒ ＝blob），件 `logs/l2` 逐字自证行：
**`EVAL_PROD_HITS=0`** ⇒ **`Eval` 今天有 0 枚生产调用者。**
**尺 2（控制面上的调用形状）**＝`grep … | grep -oE "\.(SetHtml|Eval|Dispatch|Init|Bind|…)\(" | sort | uniq -c`，件 `logs/l9`，非测试面逐字计数：
`.Dispatch(=8`、`.SetHtml(=3`、`.Bind(=2`、`.Init(=2`、`.Run(=1`、`.SetSize(=1`、`.Terminate(=2`、`.Destroy(=2`、`.Window(=1`、**`.Eval(=0`**（名册里没有它）。
**尺 3（测试面调用者，另一把尺）**＝`grep -rn -E "\.Eval\(|func .*\) Eval\(" --include=*.go cmd`（`l2`）＝7 行，其中**真调用点只有两枚**：
`panel_resident_windows_test.go:175`（`//go:build windows`，CI 真跑）＋`panel_transport_live_35v2_windows_test.go:159`（`windows && winlive`，**⛔ CI job**，逐字凭据 `panel_host_windows_live_test.go:18` `winlive has NO CI job`）；
另两行（`docSink33r10:103`／`geoSink255r1:127`）是**假控制器的方法定义**，⛔ 调用者。⇒ **两把尺⛔ 相加。**
⚠ `Eval` 确实在导出接口上（凭据是类型而非计数：`currentWindow()` 返回 `webview2.WebView`（`panel_host_windows.go:302`），`panel_resident_windows_test.go:170→:175` 就在它上面调 `Eval` ⇒ 编得过 ⇒ 接口含 `Eval`）。

⇒ **这一格要写的话**：修②**⛔ 有任何现存生产路径受益**；它能被"修好"今天唯一的意义＝让 nail2 那枚只有测试在测的维度可用，而票面 `AC#2` 的硬约束正写着⛔ 放宽那句判据 ⇒ **"为变绿而动产码"那一族风险这格是真的**，交我裁，本腿⛔ 自便。

## B. 代价表

| 格 | 要动的东西（对象层锚点） | 今天的实际形状 | 代价 | 撞哪条已裁 |
|---|---|---|---|---|
| **修①（AC13，入口交接）** | `panel_host_windows.go:447 coldStartPageHandover`（`:448`/`:450`/`:451` 三句） | **次序已经是对的**（`:486` 入口在最后）⇒ 真缺口是"最后一次 `SetHtml` 的文档⛔ 落地成 Eval 看到的那份"（表①§E／表③§B-(ii)） | ①要一枚"文档落地"的确认 ⇒ 落进下面 C 那两种形状之一；②`SetHtml`＝库的 `NavigateToString`（`webview.go:394`），⛔ 完成信号、⛔ 错误返回（`pkg/edge/chromium.go:123-128` 写作 `_, _, _ =`） | ⛔ 放宽 AC13 判据（票面 `AC#1` 硬约束）；⛔ 动探针的存在性（它本身就是"冷可用"的度量点，`:330` 逐字写着 "the probe stays"） |
| **修②（nail2，`Eval` 推）** | 推那一跳：`panel_host_windows.go` 里⛔ 任何 `Eval` 调用者（尺 1＝0）⇒ 落地形状必然是"新增一枚受门控的推"或"在 `coldStartPageHandover` 里等就绪" | 现有唯一形状＝测试的 `post→Dispatch→Run→Eval`（表③§A 跳 3-6） | 见 C | `A502` P1（本表 §C） |
| **★修①买⛔ 到的东西（一票必须具名的成本）** | 入口文档 `frontend/dist/index.html:19` `<script type="module" crossorigin src="./assets/index-B8yINMF1.js">`＋`:20` `<link … ./assets/index-yy8KMgdf.css>`（`logs/l6` 逐字） | 台账 `A798` 已写：那两条相对 URL 取 js/css 的路**"在本仓从未注册过"**，且 `:15` 那枚冻结 CSP 是 `default-src 'none'; script-src 'self'`（`A798` 实测"那枚冻结的 CSP `<meta>` 在 `SetHtml` 产出的文档里是真执行的"） | ⇒ **即便 AC13 变绿，屏上仍是 1044 字节那份里一枚空 `<div id="root"></div>`**：`id="root"` 会present（判据过的形状），页面内容⛔ 存在 | 这一格归〔待人拍板〕**`Q-84`**（`A798` 三栏，默认＝**不做**，口令「299 走甲」／「299 走乙」／「面板先不管」）⇒ **⛔ 塞进本票**，但**修①的收益边界**必须由票面写清，否则下一腿会以为 AC13 绿＝面板能看 |

## C. ★必答：动②撞⛔ 撞"面板线程＝专用 STA ＋ 库 `Run()` 泵"（`Eval` 只能在泵线程上跑）

**已裁原文（现跑台账，⛔ 转述）**：`A502` 裁定 **P1** 逐字＝「**面板归专用一条 STA 线程、泵用依赖库自己的 `Run()`；⛔ 不许把 `bringUp` 投进球的 `ui-sta`**」
（凭据＝`docs/reports/pending-and-issues.md:15051` 里 `A798` 抄回的逐字，它自己注明"我 10:3x 现跑 `sed -n '10287,10292p'` 对过原句"；`10287` 那行就是 `A502` 本格）。
盘上对拉（现码）：`panel_resident_windows.go:262 LockOSThread`＋`:264 COINIT_APARTMENTTHREADED`＋`:299 w.Run()`＋`:349/:360 post→Dispatch`；`panel_host_windows.go:136-145` 结构体注释与 `:419-421` 交接点逐字与此一致。
**并且：今天 `Eval` 本来就跑在那根泵线程上**（表③§A 跳 4：`Dispatch` 队列的唯一读者＝`Run()`，`Run()` 占住的就是那根锁定 STA 线程）。

**分形状答（⛔ 一刀切"撞/⛔ 撞"）：**

- **形状甲＝⛔ 换线程、⛔ 换泵、⛔ 动依赖**：在同一根线程、同一套 `post→Dispatch→Run` 上，把推的**时机**改成"由页面先自报它活了（现有那道门 `wispDispatch`，⛔ 新方法名、⛔ 动 C17 白名单）"，Go 收到才 `Eval`。
  ⇒ **⛔ 撞 `A502` P1**（一根⛔ 触碰它的三句：⛔ 新线程、⛔ 自排泵、⛔ 投进 `ui-sta`）。
  ⚠ 付的价：①多一枚"就绪"报文（走既有 `panel.mode.request` 的 `requestId`）⇒ **判据⛔ 被放宽**（"页面自己报 title"那一句原样保留），但**测的是"推得进已就绪的文档"⛔ 是"推得进任何文档"** ⇒ 这一格语义变化必须写进票面 `AC#2`，由我裁；
  ②仍然⛔ 保证写进的那份 DOM 是第二发读到的那份（表③缺跳 (i) 的机制判不死，欠 R1/R4）。
- **形状乙＝要回程**（`ExecuteScript` 带 callback，或订阅 `NavigationCompleted`）⇒ **撞，具名两处，交我裁，本腿⛔ 绕**：
  1. **依赖面**：库的导出接口⛔ 任何"带回调的 Eval"或事件订阅（尺＝`logs/l9`：`common.go:26` 那个 `WebView` interface 只给 `Run/Terminate/Dispatch/Destroy/Window/…`；`Chromium` 的导出方法名册里**⛔ `add_` 那一族**，`NavigationCompleted` 在 `pkg/edge/chromium.go:339` 是**库自己的 COM vtbl 处理器**，⛔ 是给人的订阅点）
     ⇒ 要拿回程只有：改/patch/`replace` 依赖（＝`A798` **乙形**，逐字撞"⛔ 不新增依赖边"；且 `test -d vendor`→**NO**、`git ls-files vendor`→**0 行** ⇒ 真做就是动 `go.mod`/`go.sum`/`deps.toml`/`allowlist.txt` 那一族＝**人工批准面**），
     或从 `cmd/wisp` 直接吃 `pkg/edge`（`GetController()` 导出在 `chromium.go:328`，自持 `ICoreWebView2Controller`→`get_WebView`→`add_NavigationCompleted`）⇒ **包级 import +1**（`A798` 甲形逐字承认"⛔ 不算新增依赖边"但同一格警告：`go.mod:19` 那枚 `// indirect` 注释会被 `go mod tidy` 摘掉 ⇒ `go.mod` 动 1 行）。
  2. **泵面**：`add_NavigationCompleted` 的回调派发依赖那根线程的消息循环 ⇒ 自持 COM 事件＝`A798` **甲形**那一族，逐字代价＝「**票 33 已落并勾过的泵形状要重来一遍**」，`A798` 自己给的台账锚就是 `A502` P1。
     ⇒ **所以形状乙撞的就是这条已裁**，且它与〔待人拍板〕`Q-84` 是**同一格**（默认＝不做）⇒ **⛔ 落地腿自便，⛔ 本腿绕，交编排者裁。**
- **顺带一处依赖面 hazard（现码逐字，给裁量用）**：`pkg/edge/chromium.go:138-142` 里 `windows.UTF16PtrFromString(script)` 失败那支写的是 **`log.Fatal(err)`** ⇒ 任何含 NUL 的脚本串会把进程直接打死。今天唯一给 `Eval` 喂变量的是测试（`:175`）与 winlive 那枚（`:159`）；**一旦②被产品化（往面板推内容），这就是"页面内容能打死宿主进程"那一族形状。**
  现象在哪出现：只在依赖源码 `pkg/edge/chromium.go:141`（库内），**⛔ 在页面上出现**；有没有本机被入侵的证据：**无**（本程只读码，⛔ 任何运行证据）；最坏后果：进程被 `log.Fatal` 直接终止（拒绝服务级，⛔ 任意执行）。⇒ 只登记形状，**未定性**（⛔ 本票射程，也⛔ 本腿能裁）。

## D. 一句话结论（表④）
**`Eval` 今天生产调用者＝0（尺 1 `EVAL_PROD_HITS=0`，射程＝`cmd/`＋`internal/` 非测试面；测试面只有两枚真调用者，另一枚在⛔ CI 的 winlive 档）；动②的"等页面先自报再推"那一形⛔ 撞 `A502` P1（`Eval` 本来就在泵线程上跑），而"要回程"那一形撞——依赖面（库导出接口⛔ callback／⛔ 订阅点）＋泵面（`A798` 甲形逐字"泵形状要重来一遍"，锚＝`A502` P1），且与〔待人拍板〕`Q-84` 同格 ⇒ 具名报回，由编排者裁；★另须写进票面的收益边界：修①变绿⛔ 等于面板有内容，入口那两条相对 URL 的 js/css 路"从未注册过"（`A798`），落点仍是空 `<div id="root">`。**
