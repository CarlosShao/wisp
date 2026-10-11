# 305-a2 · 20 表②：两枚症状各自的「文档交接／`Eval` 那一跳」逐笔归属

⚠ 本件⛔ 任何颜色读数：本腿⛔ `go test`／`go build`／`go vet`／真窗（车道）。下面每一条都是**盘上读到的调用序**，凡推到"运行时到底谁在场"那一层，一律标〔读码推到〕，判死归 30 号件的配方。

## 1. 生产面「活文档」名册：`bringUp` → `firstRoundTrip` → `serveEntry` 的装配序

尺＝逐枚调用点（HEAD blob `cmd/wisp/panel_host_windows.go`／`panel_resident_windows.go`；工作树对 HEAD 干净：`git status --porcelain -- cmd internal`＝**0 行**、`git diff --quiet HEAD -- cmd/wisp/panel_host_windows.go`＝**rc=0**、`git diff --quiet HEAD -- cmd/wisp/panel_resident_windows_test.go`＝**rc=0**，件 `logs/20-bringup-chain.txt`／`logs/20-firstroundtrip.txt`／`logs/27-transport-and-post.txt`，各 rc=0）。

调用序（⛔ 按注释读，全按语句读）：

| 步 | 语句锚（文件＋该句内容） | 这一跳把哪份文档变成活的 | 那份文档的字节与指纹 |
|---|---|---|---|
| 0 | `bringUp`（`panel_host_windows.go:318`）→ `installPanelTransport(w, ctx)`；体内逐字 `w.Bind(panelDispatchBinding, func(raw string) string { reply, _ := m.dispatchRaw(ctx, raw); return reply })` 然后 `w.Init(panelPostMessageForwardInit)` | ⛔ 新文档；绑定与 Init 脚本装在这枚控件上，Init **逐字**注释"re-runs on every document so it survives later SetHtml rebuilds" | — |
| 1 | `m.created = true`；`rtMs := m.coldStartPageHandover(ctx, t0)`（`:421`） | — | — |
| 2 | `coldStartPageHandover`（`:447`）第一句 `rtMs := m.firstRoundTrip(ctx, t0)` | — | — |
| 3 | `firstRoundTrip`（`:867`）：`w.Bind("wispProbeRT", …)` 后 `w.SetHtml(`<!doctype html>…<script>if(window.wispProbeRT)window.wispProbeRT();</script>…`)`（体内那一行，锚 `:892`） | **探针** | **136 B**（编排者现算，票 305 §1 表①）；⛔ `<title>` ⇒ `document.title == ""`；⛔ 任何 `id=` |
| 4 | `firstRoundTrip` 的循环体里 `pnlPumpOnce()`（锚 `:897`；这条链里**唯一**的手工泵点），直到 `done` 关闭／5s 有界等待／ctx 结束 | 仍是探针（泵让它活起来并被页面回话） | — |
| 5 | 回到 `coldStartPageHandover`：`if err := m.serveEntry(); err != nil { m.serveNotBuiltNotice() }` | **入口**（若 embed 有货） | `serveEntry`（`:472`）体内逐字 `if m.assets == nil \|\| !m.assets.Built() { return fmt.Errorf(…not built…) }` → `m.assets.Resolve(panel.EntryFile)` → `w.SetHtml(string(data))`（锚 `:486`）。本机那份＝**1044 B／`sha256 9b7856b9…074f`／mtime `2026-10-10 08:51:10.5 +0800`**，内含 `<title>Wisp</title>` 与**唯一一枚** `id="root"`（尺 `logs/29-entry-surface.txt`，rc=0） |
| 6 | `serveEntry` 报错那一支：`serveNotBuiltNotice()`（`:460`）体内 `w.SetHtml(`…<title>panel assets unavailable</title>…`)`（锚 `:467`） | **公告** | **181 B**；`document.title == "panel assets unavailable"`；⛔ 任何 `id=` |
| 7 | `coldStartPageHandover` 与 `bringUp` **return**；体内⛔ 任何后续泵或等待 | ⛔ — | ★**按语句：最后一次 `SetHtml` 之后⛔ 再无任何手工泵点**（第 4 步那个泵在 `serveEntry` **之前**） |

**逐字结论（静态）**：装配序上"最后一份文档"⛔ 是探针，应是入口（步 5）或公告（步 6）。
**读数结论（票面现量，非我跑的）**：两枚症状都说"运行时在 Eval/探针那一刻，活着的文档⛔ 是第 5、第 6 步写进去的那一份"。

## 2. 生产面「谁把文档换掉」只有这三枚 `SetHtml`（我重跑过）

- 尺A：`grep -rn --include=*.go 'SetHtml' cmd internal`（HEAD 工作树）＝生产面**只有 3 枚调用点**：`panel_host_windows.go:467`／`:486`／`:892`。其余命中＝注释（`:52`/`:95`/`:314`/`:457`/`:824`/`:825`/`:890`）与两枚**测试替身**的方法定义（`panel_pageover_33r10_windows_test.go:72`、`panel_reshow_255r1_windows_test.go:130`，即无头驱动用的 `docSink`）。件 `logs/`（同把尺的原始输出在 `20-handover` 取数那次，rc=0）。
- 尺B：`grep -rn --include=*.go '\.Navigate(\|NavigateToString' cmd internal \| grep -v '_test\.go'` ＝ **0 行**（rc=1，grep 无命中）。⇒ **与票 305 §1 表① 那枚"生产面⛔ 导航调用点"的读数一致**（我这发独立复跑）。⇒ 活文档**只**由上面那 3 枚 `SetHtml` 决定。

## 3. ★一枚可判死的静态推论（这两枚症状"部分同源"的形状）

三枚候选文档的指纹（§1 表最后一列）是**互不相同的可读数**：

| 那一刻活着的文档 | `document.title` | 入口那枚 `id="root"` 在否 |
|---|---|---|
| 探针（`:892`） | `""` | ⛔（答 `"0"`） |
| 公告（`:467`） | `"panel assets unavailable"` | ⛔（答 `"0"`） |
| 入口（`:486`） | `"Wisp"` | ✓（答 `"1"`） |

⇒ **nail2 的红句逐字 `the page reports its title as ""` 自己就把"哪一份文档在场"报了**：⛔ 公告、⛔ 入口，**＝探针**。（这一格是**读码＋字面量比对**推的，⛔ 我跑出来的色；判它只要一发 nail2，配方见 30 号件 §2。）
⇒ 而 AC13 的判据只问 `getElementById` 的位图，**答不出**"探针还是公告"（两者都答 `"0"`）⇒ 两枚症状共享的上游＝**第 5/6 步那次 `SetHtml` 没成为那一刻在场的文档**；症状①的判据⛔ 能区分"⛔ 执行"与"执行了但⛔ 落地"，症状②的红句能。这与编排者"部分同源、⛔ 合并"的裁**相容**，并多给一格：**修①如果只让 `serveEntry` 被执行、⛔ 让它落地，nail2 会照旧红**（跨维度负控那发的判读就靠这一格）。

## 4. `Eval` 那一路：调用链与它落到哪份文档实例

尺＝`grep -rn --include=*.go '\.Eval(' cmd internal`（件 `logs/12-eval-sites.txt`，rc=0）：

- **生产面 `Eval` 调用者＝0 枚**（与票 305 §1 表④ `EVAL_PROD_HITS=0` 独立复现一致）。全部命中＝注释与两枚测试。
-  failing 的那一路：`TestAC14GoSideEvalPushReachesThePage`（`panel_resident_windows_test.go:855`）连发两枚 `evalOnPanelThread(t, rp, …)`：第一枚逐字 `evalOnPanelThread(t, rp, fmt.Sprintf("document.title = %q;", pushed))`，第二枚 `evalOnPanelThread(t, rp, "window.wispDispatch("+reportJSEnv("ac14-push", "document.title")+");")`；判据逐字 `if answer != pushed { t.Errorf("Go's Eval push did not reach the document: …") }`。
- `evalOnPanelThread`（`:166` 起）体内逐字：`posted := rp.post(func() { w := rp.mgr.currentWindow(); if w == nil { done <- fmt.Errorf("no window on the panel thread"); return }; w.Eval(js); done <- nil })`，等待用 `panelThreadWait`，报的是"线程⛔ 跑这次请求"或"⛔ window"，⛔ "脚本没落地"。
- `currentWindow()`（`panel_host_windows.go:302`）返回 `m.w`——**一个 manager 只有一个控件实例**⇒ ⛔ "落到另一枚控件"这一形；只剩"落到同一控件的**另一份文档**"或"脚本被丢"两形。
- `residentPanel.post`（`panel_resident_windows.go:349`）逐字：`if rp.wRef != nil { w := rp.wRef; …; w.Dispatch(fn); return true }` ⇒ **建窗之后每一枚 Eval 闭包都走库的 `Dispatch` 队列**；该文件头注释逐字："webview.Dispatch only appends a closure to the library's private queue and posts a thread message, and **Run() is that queue's** reader"（另见 `panel_host_windows.go:35-38`、`:143`）。库的泵交接口＝`panel_resident_windows.go:299` 那一枚 `w.Run()`，而 `loop()`（`:259`）在 `for { rp.drainTasks(); if rp.isCreated() \|\| rp.startUpErr() != nil { break } … }` **之后**才进 `Run()` ⇒ 〔读码推到〕首扇窗的 `bringUp`（含步 5 那次 `SetHtml`）跑在 `drainTasks()` 里、`Run()` **尚未接管**，而那一次 `SetHtml` 之后这条链⛔ 手工泵（§1 步 7）。
- **两形分不开的那一截＝依赖侧，⛔ 本票射程**：库的 `Eval` 在第三方模块里把回程回调参数传 0（`…/go-webview2@v0.0.0-20260205173254-56598839c808/pkg/edge/chromium.go` 内 `e.webview.vtbl.ExecuteScript.Call(…, _script, 0,)`——该枚文件⛔ 在仓里，仓⛔ `vendor/`，出处＝票 305 §1 表③ 与我这一单）。**要回程＝fork／`replace`＝换依赖＝人工批准，⛔ 本票动⛔ 到。** 本腿只登记"它在依赖里"。

## 5. 三条边界（⛔ 再往上撞）

1. ⛔ 动依赖面（§4 末条）。
2. ⛔ 提议动"面板线程＝专用 STA ＋ 库 `Run()` 泵"那套形状（撞已定案面＋〔待人拍板〕`Q-84` 同格）。
3. `70b00885` ＝**已否证**（`git show 70b00885 -- cmd/wisp/panel_host_windows.go`＝`bringUp` 尾段四行原样搬进 `coldStartPageHandover`，同语句同顺序，除注释外⛔ 语义变化）。我在名册里也未列它——它⛔ 在 `416d9d56..f718e9b6` 尺内（该窗 S2＝37 枚的清单里⛔ 此号），⛔ 需要再判一次。
