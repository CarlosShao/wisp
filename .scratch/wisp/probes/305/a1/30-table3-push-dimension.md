# `305-a1` · `30` 表③：nail1 绿／nail2 红是不是同一条判据——push 那一维今天缺哪一跳

## A. 两枚维度各自的路径（逐跳，文件:行＝HEAD blob，核过 `logs/l7`）

**nail1＝awaited binding reply（今天绿：`r2`/`r3`/`r4`/`r9` 四发都逐字 `REPLIED,REPLIED,REPLIED`）**
页的 `await window.wispDispatch(env)` → 库的 bound stub（`webview.go:450 Bind` 生成的那份）→ `window.external.invoke(frame)`
→ `chrome.webview.postMessage(frame)`（库自己的 init 脚本，`pkg/edge/chromium.go:112` 逐字
`e.Init("window.external={invoke:s=>window.chrome.webview.postMessage(s)}")`）→ Go `msgcb` → `dispatchRaw`
（`panel_host_windows.go:841`）→ `m.disp.Handle`（`:848`）→ **返回值**沿库的回程（`pkg/edge/chromium.go:242` 那枚
`PostWebMessageAsString`）送回页 ⇒ 页的 promise 解开。

**nail2＝Go-side Eval push（今天红：四发逐字 `title=""`）**
| 跳 | 调用点（逐字） | 在哪个线程／队列上 | 返回什么 |
|---|---|---|---|
| 1 | `panel_resident_windows_test.go:863` `evalOnPanelThread(t, rp, fmt.Sprintf("document.title = %q;", pushed))`（`pushed` ＝ `"PUSHED-33R5-OK"`，同用例 `:862`） | 测试 goroutine，**只负责投** | — |
| 2 | `panel_resident_windows_test.go:166 func evalOnPanelThread` → `:169 rp.post(func(){ … })` | 同上 | `posted bool` |
| 3 | `panel_host_windows_test`→ `panel_resident_windows.go:349 func (rp *residentPanel) post` → `:360 w.Dispatch(fn)`（`wRef` 非 nil 支）；⛔ nil 支走 `rp.tasks` 通道 | — | bool |
| 4 | 库 `webview.go:443 func (w *webview) Dispatch(f func())`：`dispatchq = append(…)` ＋ `User32PostThreadMessageW(w.mainthread, WMApp, 0, 0)` | **队列的唯一读者＝`Run()`**（`panel_resident_windows.go:299`，线程＝`:262 LockOSThread` ＋ `:264 COINIT_APARTMENTTHREADED` 那条**专用 STA**） | 无 |
| 5 | 库 `webview.go:439 func (w *webview) Eval(js string)` → `w.browser.Eval(js)` | **泵线程上** | **void** |
| 6 | 依赖 `pkg/edge/chromium.go:138 func (e *Chromium) Eval(script string)` → `e.webview.vtbl.ExecuteScript.Call(webview, script,`**`0`**`)`（`logs/l3` 逐字：第三个实参是 `0`＝**NULL callback**；返回值写作 `_, _, _ =`） | 泵线程上的 COM 调用 | **⛔ HRESULT、⛔ COM 错误、⛔ 脚本结果、⛔ 完成信号**——全丢 |
| 7 | 判据那一跳：`:864` 第二发 `evalOnPanelThread(… "window.wispDispatch("+reportJSEnv("ac14-push", "document.title")+")")` → 门 → `awaitReport`（`:866`）→ 红句 `:867`／`:869` | 走 nail1 那条回程 | 页面自报的 title |

⇒ **页面的 `document.title` 今天是谁设的**：
- 生产代码**⛔ 一处**。尺＝`logs/l2`（`cmd internal` 非测试面 `Eval` 调用者 **0 枚**）＋`logs/l3`
  （库的 `webview.go:397 SetTitle` 走 `SetWindowTextW`＝**Win32 窗口标题**，`panelTitle`＝`"Wisp panel"` 只到窗口，⛔ 进 DOM）。
- ⇒ `document.title` 唯一的来源是**文档自己字节里那枚 `<title>`**（`logs/l8` 逐字量过）：入口 `"Wisp"`／公告 `"panel assets unavailable"`／探针 **⛔ `<title>` ⇒ `""`**。

## B. 必答："push 那一维今天缺的是哪一跳"——**两截，分开写**

**缺跳 (i)＝库的形状：`ExecuteScript` 那一发没有回程。**
`Chromium.Eval` 把 callback 实参写死成 `0`（`logs/l3`），所以 Go **拿不到**"这一发脚本跑没跑、在哪个文档上跑、抛没抛"。
⇒ 缺的⛔ 是通道（第二发 Eval 走同一通道**确实到达了页面**——`awaitReport` 真收到了那枚 env，红句里带的是页面自己的字），
缺的是**"第一发的副作用进不到第二发读到的那份 DOM"**这一格：`document.title = "PUSHED-33R5-OK"` 之后，同文档上的第二发读到 `""`。
⇒ **静态读码判不出机制**（脚本被丢／跑在另一份文档实例／文档在两发之间被重建，三种形状都⛔ 这枚通道）。
**能量它的只有 R1（同台面 `-count=3`，看红是不是窗数/顺序依赖）＋R4（`-tags winlive`，`panel_transport_live_35v2_windows_test.go:159` 那枚本仓唯一另一处真窗 `Dispatch`+`Eval` 面）。**

**缺跳 (ii)＝与 AC13 共用的那一截：最后一次 `SetHtml` 的那份文档⛔ 是 Eval 跑在其上的那份。**
凭据＝三份签名互斥（`logs/l8`）× 两枚症状页面自报的答（`"0"` ＋ `title=""`）唯一匹配 136 字节探针页（`panel_host_windows.go:892`）。

## C. 同源／⛔ 同源（明写，并给凭据；这条直接决定票面分层假设⛔ 作废）

**判：两枚症状⛔ 是同一跳缺；它们**共享上游 (ii)，nail2 还多一枚下游 (i)**。**
★那条**跨台面判据**是本表最硬的一枚，只靠已有读数＋对象层就能算，⛔ 新色：
- **干净 clone**（`r3`/`r4`）上 `serveEntry` 必然失败（逐字 `:315` `Resolve(entry): panel: embedded assets are not built (run npm run build in frontend/)`），
  ⇒ 那一发的**最后一次 `SetHtml` 是 `:467` 那 181 字节公告页**（`document.title` 应为 `"panel assets unavailable"`）。
- **母仓**（`r2`）上入口解析成功 ⇒ 最后一次 `SetHtml` 是 `:486` 那 1044 字节入口页（`document.title` 应为 `"Wisp"`）。
- **两把台面的"最后文档"⛔ 同，nail2 的红一句没变**（四发逐字都是 `title=""`）。
⇒ 所以 nail2 的判据⛔ 吃"入口交接成功"那一枚修复；把它修好⛔ 自动把 nail2 带绿（在**判据层**上）。
⇒ **票面的分层假设⛔ 作废**，但要把话写准：两枚症状在"最后一份文档⛔ 落地"这一截上**会一起动**，
nail2 另外那截 (`i`) 是**它自己独有的**，AC13 的判据（单发 Eval、只读 id）**根本覆盖⛔ 它**。
⇒ **能裁"同源/⛔ 同源"的唯一一发＝票面 `AC#3` 那枚跨维度负控**（修①那一笔落地后看 nail2 的色）；本腿⛔ 产码，只把判据结构交回来，⛔ 替这格下判语。

## D. 与派单/票 33 裁语的一处具名冲突（尺＝现跑台账）
派单写"票 33 的 `A475` 那批裁过'两维⛔ 混一枚'"。**现量⛔ 对不上**：
`grep -n "A475" docs/reports/pending-and-issues.md` ⇒ `9890:## A475 — 09-30 **16:28:17** …：收「228-r1」⇒ 球第一次真住进了会跑任务的那条腿…新立「票 245」`——**那格是球/线程与票 245 的事，⛔ 面板两维**。
"两维⛔ 混一枚"那条**实际落在 10-01 13:12 那道裁定（P2）**上，凭据两处逐字：
- `cmd/wisp/panel_resident_windows_test.go:11-16` 用例头注释：`(i) the reply to an awaited JS binding, which only the library's Run() delivers, and (ii) a Go-side Eval push, which arrives even without it. One instrument may not stand in for both (orchestrator ruling P2)`
- `cmd/wisp/panel_host_windows.go:860-862`（`firstRoundTrip` 注释）：`AC#14's receipt hop has its own ruler … and a second one for Go-side Eval push; the two are separate dimensions and stay separate tests.`
- 台账那一批的锚＝`10287:## A502 — 10-01 **13:13:52** …：收 33-p1 … 票 33 的 AC#14 与"再入是哪一支"就此定案`；
  `A502` 裁定 **P1** 逐字（由 `A798` 抄回）＝「面板归专用一条 STA 线程、泵用依赖库自己的 `Run()`；⛔ 不许把 `bringUp` 投进球的 `ui-sta`」。
⇒ ⚠ 顺带一处**注释与读数冲突**具名：`:14` 那句 "a Go-side Eval push, **which arrives even without it**"（＝⛔ Run() 也到）与 P2 的另一半在今天的盘上**⛔ 复现得了**（缺跳 (i) 的机制判不死）⇒ 引用它当既成事实＝把注释读得比代码乐观。

## E. 一句话结论（表③）
**两枚维度⛔ 同一条判据（A502/P2 立的两枚用例各自钉一维，四发读数里 nail1 逐字绿、nail2 逐字红）；push 那一维今天缺两截——库侧 `ExecuteScript(callback=0)` 那枚没有回程的跳（`pkg/edge/chromium.go:144`，Go 侧 `Eval` 返回 void）＋与 AC13 共享的"最后一次 `SetHtml` 的文档⛔ 是 Eval 落到的那份"。凭据：clone（最后文档＝181B 公告）与母仓（最后文档＝1044B 入口）两把台面下 nail2 一句没变 ⇒ 分层假设⛔ 作废，但 nail2 独有那一截 AC13 判据覆盖⛔ 了，能裁它的只有票面 `AC#3` 那枚跨维度负控。**
