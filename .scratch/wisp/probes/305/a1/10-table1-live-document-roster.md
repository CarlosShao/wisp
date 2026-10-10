# `305-a1` · `10` 表①：冷启动这条装配序里"活文档"名册

尺＝**逐枚调用点 `文件:行`（HEAD 对象层 `a81c2980`，28 枚锚逐枚 `sed` 核过全 MATCH，件 `logs/l7`）＋调用序（静态读码）**。
⛔ 按注释读——本表里注释与代码**确实不一致**两处，具名在 §D。
射程＝`cmd/wisp/**`＋`internal/**` 的工作树（该两目录 porcelain 0 行 ⇒ 工作树＝blob）；字节数＝`logs/l6`／`logs/l8`。

## A. 装配序（自线程起，到最后一枚 `SetHtml` 止）

| 序 | 调用点（blob 行） | 这一跳做什么 | 对"活文档"的影响 |
|---|---|---|---|
| 1 | `cmd/wisp/panel_resident_windows.go:262` `runtime.LockOSThread()`；`:264` `CoInitializeEx(COINIT_APARTMENTTHREADED)` | 面板专用 STA 线程 | ⛔ |
| 2 | `cmd/wisp/panel_resident_windows.go:299` `w.Run()` | 库自己的泵占住这条线程；它是库私有 dispatch 队列的**唯一读者** | ⛔ |
| 3 | `cmd/wisp/panel_resident_windows.go:349 func (rp *residentPanel) post` → `:360 w.Dispatch(fn)` | 窗口存在后一切请求进 `Dispatch` 队列 | ⛔ |
| 4 | `cmd/wisp/panel_host_windows.go:505 Show` → `:512 m.bringUp(ctx)` | 未创建则建窗 | ⛔ |
| 5 | `cmd/wisp/panel_host_windows.go:318 bringUp` → `:386 webview2.NewWithOptions(...)` | 建控件（`WindowOptions` 由 `:392 windowOptions()` 现读几何） | **未定性**，见 §C-注 |
| 6 | `:408 installPanelTransport(w, ctx)` → `:829 w.Bind(panelDispatchBinding, …)` | 绑 `wispDispatch` 那道门（先于任何文档） | ⛔ 换文档 |
| 7 | `:835 w.Init(panelPostMessageForwardInit)` | 库的 `AddScriptToExecuteOnDocumentCreated`（`logs/l3`）⇒ **每份文档重建都重跑** | ⛔ 换文档 |
| 8 | `:421 rtMs := m.coldStartPageHandover(ctx, t0)` | 文档交接那一族的唯一入口 | — |
| 9 | `:448 m.firstRoundTrip(ctx, t0)` → 函数体 `:867` | `:879` 再绑一枚 `wispProbeRT`；**`:892 w.SetHtml(行内探针文档)`** | **文档 #1 ＝ 探针页（136 字节）**；随后 `:895-911` 用 `pnlPumpOnce()` **手排**到 `done` 或 5s 死线 |
| 10 | `:450 m.serveEntry()` → 函数体 `:472` | `:476 m.assets.Resolve(panel.EntryFile)` 取字节；**`:486 w.SetHtml(string(data))`** | **文档 #2（成功支）＝内嵌入口（母仓台面 1044 字节）＝代码序上的最后一次** |
| 11 | `:451 m.serveNotBuiltNotice()`（**只在 `:450` 返回错误时才跑**）→ 函数体 `:460` | **`:467 w.SetHtml(行内公告文档)`** | **文档 #2（失败支）＝公告页（181 字节）＝该支的最后一次** |

⇒ **"最后一次是谁"**：成功支＝`:486`（入口字节）；失败支＝`:467`（公告）。
**两支都⛔ 是 `:892` 那枚探针文档。** 这条是硬结论，尺是名册本身（下表 B）。

## B. 名册尺（非测试面的 `SetHtml`／导航调用点，逐字）
尺＝`grep -rn "SetHtml" --include=*.go cmd internal`（件 `logs/l2`）；**同一把尺里"导航调用"那一节是空集**
（`grep -rn -E "\.Navigate\(|NavigateToString"` 在 `cmd internal` 非测试面 ⇒ **0 行 ⇒ 本仓生产面⛔ 任何直接导航调用点，文档只经 `SetHtml` 换**；空集＝结论，⛔ 读成"没查到"）。

非测试面 `SetHtml` 真调用点＝**三枚**（`l2` 里另 7 行⛔ 注释）：`:467`／`:486`／`:892`。

| 文档 | 调用点 | 字节从哪来 | 字节数（`logs/l8`） | `document.title` | `getElementById('root')` |
|---|---|---|---|---|---|
| 探针页 | `panel_host_windows.go:892` | **源内行内字面量**（两段拼接），与 embed／dist 无关 | **136** | `""`（⛔ `<title>`，hit 0） | **null**（⛔ 任何 `id=`，hit 0） |
| 内嵌入口 | `:486` | `m.assets.Resolve(panel.EntryFile)` ← `//go:embed all:dist` 的 `frontend/dist/index.html` | **1044**（母仓台面） | `"Wisp"` | **存在**（`id="root"` hit 1，该行逐字 `<div id="root"></div>`） |
| 公告页 | `:467` | 源内行内字面量 | **181** | `"panel assets unavailable"` | null |
| "空串"那一支 | —— | **不存在**：尺＝对 HEAD blob 数 `SetHtml("")` 那形 ＝ **0 枚**（`l8` 末节逐字） | — | — | — |

⚠ 入口那 1044 字节是**母仓工作树台面读数**，那枚文件在 HEAD ⛔ 受跟踪（`git ls-tree -r --name-only HEAD -- frontend/dist` ⇒ 只有 `.gitkeep`，`l6` 逐字可见），
且盘上 mtime＝**2026-10-10 08:51** ⇒ 它带着 `:13-16` 那枚冻结 CSP `<meta>`（`default-src 'none'; script-src 'self'`）、
`:19 <script type="module" src="./assets/index-B8yINMF1.js">`、`:20 <link ... ./assets/index-yy8KMgdf.css>`。
⇒ **"台面"在这张票里有两个轴**：① 编排者票面写的 母仓↔clone（bundle 有无）；② **同一枚母仓里那枚未跟踪 dist 的内容自己漂过**（10-01 那发绿读数用的是旧字节）。

## C. `//go:build` / `-tags winlive` 射程（哪几枚今天只有真窗／特定档才求值）
- `cmd/wisp/panel_host_windows.go` 首行＝`//go:build windows`（`l2` 末节逐字头名册）⇒ **§B 那三枚 `SetHtml` 全在 windows 档**；⛔ windows ⇒ 根本不编进去。
  ⇒ 三枚调用点都**只在真开了 WebView2 控件时求值**——本腿⛔ 编译面，所以本表**没有任何一枚的新颜色**。
- `cmd/wisp/panel_resident_windows_test.go`＝`//go:build windows`（**⛔ winlive**）⇒ `TestAC13…`／`TestAC14Awaited…`／`TestAC14GoSideEvalPush…` 在 CI 的 `test-windows` 里真跑（凭据＝`probes/303/orch/r9-ci-after-red-roster.txt` 那两句逐字出自 job `114194107792`）。
- `-tags winlive` 射程只有三枚文件：`panel_host_windows_live_test.go`、`panel_geometry_255_winlive_test.go`、`panel_transport_live_35v2_windows_test.go`；
  其中 **`panel_transport_live_35v2_windows_test.go:159` 是本仓第二枚真·`Eval` 调用点**（形＝`wv.Dispatch(func() { wv.Eval(script) })`），
  而 `panel_host_windows_live_test.go:18` 逐字写着 **`winlive has NO CI job`** ⇒ 那一档今天⛔ CI 名册、⛔ 默认档 ⇒ 表②③ 判"Eval 推不推得进"时**它是唯一没被人跑过的第二次测量面**（具名要那一发，见 `20-` §D-R4）。
- 注（§A 第 5 行为什么留"未定性"）：尺＝`logs/l3`（依赖 blob `webview.go`／`pkg/edge/chromium.go` 的 grep）——库内**只有 `webview.go:390` 那枚 `browser.Navigate(url)`**，本仓 `cmd internal` 非测试面⛔ 调用者（`l2` 空集）。
  ⇒ `NewWithOptions` 自己给不给一份初始文档、给的是哪份，**静态读码判不死**，能量它的只有真窗读数；本腿⛔ 猜。

## D. 注释比代码乐观／过期（两处，具名）
1. `cmd/wisp/panel_resident_windows_test.go:330` 的**判据文字本身**逐字写着
   `the round-trip probe page is the last document shown, so the user sees a stub instead of the panel. Product side: bringUp must hand the entry over AFTER the probe`——
   而 §A 的序（`:448` 先探针、`:450`/`:467` 后交接）在 `a81c2980` 上**已经是"入口/公告在后"**。
   ⇒ 那句话是 33-r10 之前那一族的诊断文字，被逐字抄进了 `t.Errorf` 的串里；**它现在是一句待验断言，⛔ 既成事实**（尺＝`logs/l5` 里 `70b00885` 的逐字 diff）。
   ⇒ ⚠ 本腿⛔ 因此判 AC13 是"判据写错了"——盘上四件读数都真给了 `"0"`／`""`（`20-`），所以要判的是**"最后一次 `SetHtml` 的那份文档为什么不是 Eval 看到的那份"**。
2. `cmd/wisp/panel_host_windows.go:864-866` 注释逐字写着 `since AC#13 was fixed, is no longer the last word: bringUp runs this BEFORE serveEntry, so the user ends on the embedded entry`——
   "**so the user ends on the embedded entry**"这一半是**推论⛔ 读数**；母仓那发（`r2`）页面自己答的是 `"0"`（入口那 1 枚 id ⛔ 在场）。同一枚注释的另一半（序）与代码一致。

## E. 一句话结论（表①）
**代码序上的"最后活文档"是内嵌入口（成功支 `:486`）或公告页（失败支 `:467`），⛔ 探针页；而两枚症状里页面自报的签名（`""`＋root 缺席）在三份候选文档里唯一匹配 136 字节的探针页（`:892`）⇒ 冷启动"最后一次 `SetHtml`"的那份文档，⛔ 是测试的 `Eval` 实际跑在其上的那份。**
