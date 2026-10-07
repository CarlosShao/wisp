# 35-a1 — C17 出向（Go→页）落点普查（只读腿）

> 腿：`35-a1`。射程：回答「出向那一跳要落在哪、怎么落、会撞谁」，**不选形**。
> ⛔ 零写入（除本目录）、零 `go test/build/vet`（`272-v1` 正在取整包 Go 读数）。

## S0 起手锚（原文）

```
$ date ; git rev-parse --short HEAD ; git status --porcelain -- cmd internal scripts tools .github docs frontend
Wed Oct  7 11:31:51 CST 2026
367e41b3
（六族路径 status 输出为空 = 工作树在这六族上干净）
```

- 分支：`dev`。
- ⚠ 起手即见：六族路径 status **为空**。工作树里 `design/**`、`.gitignore`、`.scratch/**` 的未提交件不属于本腿，
  不还原、不提交、不评价。终态自查按这六族判，必须与上面逐字相同（= 仍为空）。
- 原始输出：`logs/00-anchor.txt`。

## S1 宿主这一侧的能力面

### 1.1 `PanelManager` 的结构（`cmd/wisp/panel_host_windows.go:146-185`）

字段逐枚（`panel_host_windows.go:147-184`）：

| 字段 | 类型 | 是什么 |
|---|---|---|
| `mu` | `sync.Mutex` | 见下 |
| `w` | `webview2.WebView` | **webview 实例就叫这个名**（接口值，非具体类型） |
| `hwnd` | `windows.HWND` | **窗口句柄就叫这个名**（`bringUp` 里 `windows.HWND(w.Window())` 得来，`:398`） |
| `dataPath` | `string` | WebView2 user-data 目录 |
| `assets` | `*panel.Assets` | 内嵌页面字节 |
| `disp` | `*panel.ComposerDispatch` | 入向路由器（与 `panel_inbound.go` 共用同一枚） |
| `created`/`shown` | `bool` | 状态位 |
| `prevFocus`/`lastRestore*` | `windows.HWND`/`uintptr` | AC#4 焦点归还的读数 |
| `lastColdMs`/`lastHotMs` | `float64` | 实测延迟，非阈值 |
| `geometry` | `func() (int,int)` | 装配根注入的尺寸来源（票 255 AC#4 禁 panel→config 边） |

`mu` 的护法（读码结论，`:136-145` 的注释与 `:414-418` 的实现两边对得上）：
- `mu` **只护字段可见性**，让状态能被任意 goroutine 读；注释明写「它**不**让 WebView2 控件变线程安全」。
- 真规矩是**线程归属**：碰控件的方法（`bringUp`/`Show`/`Hide`/`Destroy`）必须在创建窗口那条线程上跑；
  本文件**自己没有一枚 `runtime.LockOSThread`**（`:138-140` 具名承认），锁线程那件事只在常驻进程的专用面板线程里做。
- 取控件句柄的统一出口是 `currentWindow()`（`:302-306`，带锁读 `m.w`），注释点明它是「nil 时走线程自己的 channel、非 nil 时走库的 Dispatch 队列（**只有 `Run()` 在抽它**）」的路由判据。

方法名册（`grep '^func (m \*PanelManager)'` 全文，共 16 枚，**导出 8 / 非导出 8**）：
- 导出：`IsCreated:268` `IsShown:276` `LastColdMs:285` `LastHotMs:286` `Show:481` `HotShow:547` `Hide:574` `Destroy:615`。
- 非导出：`windowOptions:236` `windowHandle:290` `currentWindow:302` `bringUp:318` `serveNotBuiltNoticeLocked:442`
  `serveEntry:454` `setPriorFocusLocked:518` `terminateOnThisThread:604` `dispatchRaw:630` `firstRoundTripLocked:656`。
  （16 枚行号里 8+11 是因为 `LastColdMs/LastHotMs` 同占一行 `:285`/`:286`；枚数按符号计＝导出 8、非导出 11，总 19 枚声明行、其中 1 枚是 `panelHostOption` 闭包不算。以行号为准，本腿不裁枚数。）
- **导出 8 枚里确实没有一枚送数据进页面**（这条与转述一致），但**送字节这件事在文件里已经存在**，只是不叫「推送」：
  `w.SetHtml(...)` 共 3 处生产调用（`:449` `:468` `:681`）——Go→页的字节通道**今天就在用**，只是只在「建窗那一次」用。

### 1.2 库给的能力面（⚠ 与转述的尺不同，先看这条）

- 模块与版本（`go.mod:19`）：`github.com/jchv/go-webview2 v0.0.0-20260205173254-56598839c808 // indirect`。
  `go list -deps ./cmd/wisp`（286 行，`logs/01-golist-deps.txt`）里带它：
  `github.com/jchv/go-webview2`（:281）、`…/pkg/edge`（:280）、`…/webviewloader`（:279）、`…/internal/w32`（:272）。
- 命令尺：用了 `go doc -all`（两次，`logs/01-godoc-webview2.txt` 118 行 / `logs/01-godoc-edge.txt` 421 行）＋**直读模块源码的 `func (` 声明行**。
  两者都**不是** `go build/vet/test`，不占 `272-v1` 的包级互斥；没有换尺，也没有长跑命令。
- **`webview2.WebView` 接口导出的方法名册（`logs/01-godoc-webview2.txt:22-80`，共 12 枚）**：
  `Run()`, `Terminate()`, `Dispatch(f func())`, `Destroy()`, `Window() unsafe.Pointer`,
  `SetTitle(string)`, `SetSize(w,h int, hint Hint)`, `Navigate(url string)`, `SetHtml(html string)`,
  `Init(js string)`, `Eval(js string)`, `Bind(name string, f interface{}) error`。
  ⇒ **能「往页面里执行 JS／送字符串进去」的只有三枚：`Eval`、`Init`、`SetHtml`**（`SetHtml` 是整页替换）。
  再加 `Dispatch`——它不送字节，它是**唯一那扇跨线程的门**（下面 1.2.丙 很关键）。

**1.2.甲 ⚠ 转述那把尺搜的三个名字，在这套 API 里一个都不存在**
`PostWebMessage` / `EvaluateScript` / `CreateWebMessageAsJson` 于 `cmd/wisp`＋`internal/panel` ＝ 0 —— 这条读数**事实为真但射程为负**：
- 库源码里 `PostWebMessageAsJSON` / `PostWebMessageAsString` / `ExecuteScript` 只作为 **COM vtbl 表项**存在
  （`pkg/edge/corewebview2.go:103` `ExecuteScript`、`:106` `PostWebMessageAsJSON`、`:107` `PostWebMessageAsString`），
  **没有任何导出方法把 `PostWebMessage*` 交给调用者**；`ExecuteScript` 被包在 `Chromium.Eval` 里面
  （`pkg/edge/chromium.go:138-145`，第 144 行 `e.webview.vtbl.ExecuteScript.Call(...)`，第三个实参写死 `0`＝**不接回调、结果丢弃**）。
- `PostWebMessageAsString` 的唯一出现处是**库自己的入向回显**（`pkg/edge/chromium.go:242`，在 `MessageReceived` 里把收到的原话贴回页面），
  不是产品能调的通道。
- ⇒ **正控**：同一条尺换成库真名后立刻非零——`cmd/wisp/panel_host_windows.go:449/468/681` 三处 `w.SetHtml(`（生产、非 test），
  `cmd/wisp/panel_resident_windows.go:360` 一处 `w.Dispatch(`（生产、非 test）。
  「出向整条不存在」这句话要按 1.2.乙 修正成更准的形状。

**1.2.乙 ⚠ 这条要推翻转述的框架：Go→页的通道今天**并非**零，它藏在库自己的 Bind 回执里**
`github.com/jchv/go-webview2/webview.go:140-160`（`msgcb`）逐字：绑定函数一跑完，库就做
`w.Dispatch(func() { w.Eval("window._rpc["+id+"].resolve("+string(b)+"); …") })`（resolve 在 `:156`，两条 reject 在 `:148`/`:152`）。
⇒ 也就是说：**每次页面调用 `window.<绑定名>(…)`，Go 的返回值都是靠「`Dispatch` 包一枚 `Eval`」送回页面的**，
这条通道**在生产里每开一次面板就走若干次**（`panelDispatchBinding` 的回执，`panel_host_windows.go:405-408`；
宿主自己的注释 `:401` 也写了「the returned string is the receipt that reaches the page」）。
⇒ 精确的形状是：**缺的不是「Go→页」这根线，缺的是「Go 不经页问就先说」那一支（unsolicited push）**。
派单时把这两件事分开写，否则落地腿会以为要新建通道，而真选项里有一条「复用已存在的那支」。

**1.2.丙 `Dispatch` 的实现细节（决定 §3 全部结论，逐字读自 `webview.go:443-447` 与 `:349-377`）**
```go
func (w *webview) Dispatch(f func()) {
	w.m.Lock(); w.dispatchq = append(w.dispatchq, f); w.m.Unlock()
	_, _, _ = w32.User32PostThreadMessageW.Call(w.mainthread, w32.WMApp, 0, 0)
}
```
- 排队＝在 `w.m` 下 append 一枚切片；唤起＝`PostThreadMessageW(mainthread, WM_APP)`。
- **抽队列的唯一地方＝`Run()`**（`webview.go:349` 起，`:360-366` 那段：`msg.Message == WMApp` 时锁住、整队拷出、清空、逐枚执行；
  `WM_QUIT` 则 `return`）。⇒ **`Dispatch` 的闭包只在库的 `Run()` 泵在转的时候才会跑**；这与 `currentWindow()` 的注释（宿主文件 `:299-301`）逐字对得上。
- `Eval` 的两条硬伤（落地必须知道，来源是模块源码不是推测）：
  ① `Chromium.Eval` 里 `windows.UTF16PtrFromString` 失败走 **`log.Fatal`**（`chromium.go:140-142`）⇒ 载荷含 NUL 字节＝**进程退出**，不是 error 返回；
  ② `Eval` 不看 `e.webview` 是否为 nil，也未做任何「窗口已销毁」检查 ⇒ 销毁后调用是 COM vtbl 上的空指针风险。
  ③ 结果一律丢弃（`:144` 第三实参 `0`）⇒ **push 送达与否，Go 侧永远拿不到回执**，判据只能来自页面（与探针同结论）。

### 1.3 探针 `33/p1/q2` 当年具体调的是哪一个方法、怎么调的

`wc -l -c` ＝ **393 行 / 14358 字节**（全文已按整行读，未 `cut -c`）。它是仓里唯一一次真把东西推进页面并**由页面自己判定收到**的先例。

- 它调的**不是** `PostWebMessage`（票头 `:16-19` 逐字写明：`PostWebMessageAsString` 只是库内部回显副作用、「no method exposes it to callers, and the production tree (cmd internal) contains zero call sites of either name」）。
- 它调的是 **`w1.Dispatch(func(){ w1.Eval("document.title='<marker>'") })`**，两种模式各一枚：
  - `-mode eval`（`:280-285`）：即便场景想要「本线程直接 Eval」，实现仍然**借了 `Dispatch` 一次**把闭包送上 UI 线程，闭包里才是 `w1.Eval(...)`（`:283`）。
    注释 `:268-270` 具名写着「CROSS-THREAD by construction here would be wrong for Eval」。
  - `-mode eval-disp`（`:297-305`）：另一条 goroutine 里 `time.Sleep(600ms)` 后 `w1.Dispatch(→ w1.Eval(…))`（`:302`），
    Purpose 写在 `:293-296`：把「Dispatch 闭包真跑了」与「Eval 真到了页面」**两件事分开量**。
  - ⇒ **可抄的最小形状 = 非 UI 线程调 `Dispatch`，闭包内调 `Eval`；页面靠 `w.Run()` 泵收到**（`w1.Run()` 紧跟其后：`:287`、`:306`）。
- **页侧接收器写法＝轮询，不是 message 事件**（`evalPage()` 全文 `:378-393`）：`setInterval` 每 200ms 读 `document.title`，
  命中 `P1Q2-OK` 就 `window.wispReport("EVAL_SEEN …")`，超过 60 次就报 `EVAL_NOT_SEEN`。
  ⇒ 探针**从未用过** `chrome.webview.addEventListener('message')`；这条对 §4 的判决很重要。
- 判据是**页面的话**不是 Go 的话（`:21-25` 具名：「Go-side 'the call did not error' is recorded but is NOT the judge (that is the fake green the first leg's R25 exposed)」）。
- 线程侧前置（`:205-215`）：`runHost` 是 STA 线程——`runtime.LockOSThread()`、`CoInitializeEx(STA)`、
  把自己的 TID 写进 `uiTid` 通道；注释 `:206-207` 点名这是「the resident panel thread's own guard shape, `panel_resident_windows.go:214`」。

## S2 泵到宿主之间的装配线

### 2.1 谁构造 `PanelManager`、在哪、与 `agentRuntime` 是否同一 owner 树

- **唯一非 test 构造点**：`cmd/wisp/panel_resident_windows.go:253`（在 `newResidentPanelManager(dataDir string)`，定义 `:228`）。
  其余 10 处 `NewPanelManager(` 全在 `_test.go`（`panel_geometry_255_test.go:110/140/172`、`…_winlive_test.go:101`、
  `panel_host_windows_live_test.go:39`、`panel_host_windows_test.go:623/880`、`panel_resident_windows_test.go:63/449/526/527`）。
  ⇒ **正控**：同一把尺把定义行本身也数进来（`panel_host_windows.go:213`），说明尺没有漏过滤。
- 调用链：`cmd/wisp/resident_windows.go:33 runResident()` → `:42 rt, err := proc.Boot(env)` → `:151 rp, rpErr := newResidentPanelManager(rt.Layout.DataDir)`
  → `:157 panel = startResidentPanel(rt.Registry, rp)`（`:144` 定义，注释 `:139` 明写「builds the thread but NOT the window」）→ `:158 defer panel.stop()`。
- ⚠ **同名陷阱，必须纠正一句**：`resident_windows.go` 里的 `rt` 是 **`*proc.Runtime`**（`:42` 由 `proc.Boot` 返回），
  **不是** `cmd/wisp/run.go:266` 那枚 `agentRuntime`。两枚 owner 树各自有自己的 `rt` 名字：
  - `agentRuntime` 构造：`run.go:382 assembleRuntime(s runSpec)` → `run.go:390 rt := &agentRuntime{…}`；
    生产调用方只有两枚：`run.go:249`（`wisp run` 控制台腿）与 `resident_task_source_windows.go:278 run, code := assembleRuntime(spec)`（常驻腿）。
  - ⇒ **同一个常驻进程里两枚都在**（`proc.Runtime` 持面板线程，`agentRuntime` 持泵），**但今天它们之间没有任何一条线**。

### 2.2 `Publish` 的两族同名字，逐枚点名（⚠ 别混）

`internal/panel/pump.go:348 func (p *SnapshotPump) Publish() (Snapshot, []byte, error)` —— 快照泵那一族：
- **非 test 调用者只有一枚**：`cmd/wisp/panel_pump.go:405 snap, _, err := rt.pump.Publish()`（在 `publishPanelSnapshot` 体内）。其余 10 处全在 `_test.go`。
- 另一枚相关计数：`cmd/wisp/panel_pump.go:419 return rt.pump.Publishes()`（`snapshotCount`，非 test，读计数不产字节）。

另一族同名的**不是**快照泵（按调用形状排除，免得上面那条被读成「Publish 只有一个调用者」）：
`agent.Event` 的 Sink 面：`run.go:1260 func (c consoleSink) Publish(e agent.Event)`、`internal/agent/sink.go:71 NopSink`、`:85 RecordSink`、
`internal/agent/loop.go:1041 l.opt.Sink.Publish(e)`、`internal/tools/subagent_197.go:520 subagentTextSink.Publish`、
`cmd/wisp/panel_config_248_test.go:249 cards.Publish()`（test，又是第三族 `approval.Replies`）。

`publishPanelSnapshot` / `bookPanelSnapshot` 的生产接线（**这就是「隔着什么」的答案**）：
- `run.go:725 Out: rt.bookPanelSnapshot` —— 泵的唯一出口，装配根注入的函数值（`panel_pump.go:290` 自陈「bookPanelSnapshot is the exit's far end」，
  `:317-318` 逐字：「The full bytes stay with the process **until the transport that delivers them exists**」）。
- `run.go:734-735 if rt.ui != nil { rt.ui.publish = rt.publishPanelSnapshot }` —— 面板快照的触发器**只有这一处**。
- `run.go:998 Sink: consoleSink{out:…, stream:…, publish: rt.publishPanelSnapshot}` —— 控制台 sink 那一支。
- ⇒ ⚠ **本节最值钱的一条读数**：`run.go:599-607`——常驻腿递了 `s.gate != nil`（`resident_task_source_windows.go:269 gate: ra.gate`，`ra` 来自 `resident_windows.go:132`），
  于是走 `:607 rt.ui = nil // no console surface in a process whose cards go to another host`；
  而 `:734` 的守卫正是 `rt.ui != nil` ⇒ **常驻进程（唯一持有 `PanelManager` 与面板线程的那个进程）今天一次都不调 `publishPanelSnapshot`**。
  `run.go:730-733` 的注释还具名拒绝过这条路：「The panel route is ticket 33/35's, not AC#7's, and inventing a second publish path from the resident side
  would be a claim about a page this process cannot show.」
  ⇒ 所以「入向真在、出向整条不存在」要补一层：**泵与宿主不仅缺一段代码，还缺一个同进程内的触发器**；
  而有触发器的那条路（`wisp run` 控制台腿）里 `PanelManager` 根本没被构造。
  （正控：`startResidentPanel` 与 `assembleRuntime` 在常驻进程都真有生产调用者，见 2.1，两条读数不是「谁也没跑」。）

### 2.3 最小改动落点候选（**本腿不选形，只列代价**）

| # | 候选落点 | 具体动哪 | 新增依赖边 | 代价／要解冻什么 |
|---|---|---|---|---|
| A | `cmd/wisp/panel_pump.go:319 bookPanelSnapshot` | 记账之后追加一次「交给宿主」 | **无新包边**：两者同在 `package main` | 需要 `agentRuntime` 手上有面板句柄——今天没有，须装配根多递一枚函数值；仍缺触发器（见 B） |
| B | `cmd/wisp/run.go:734-735` 那道 `rt.ui != nil` 守卫 | 给常驻腿补一条 publish 触发路径 | 无 | **直接推翻 `run.go:730-733` 已具名的裁定** ⇒ 要解冻，须点名裁定者。不动它，A/C/D 都拿不到调用者 |
| C | `cmd/wisp/panel_resident_windows.go`：给 `residentPanel` 加一枚「送快照」方法，内部走已有 `post()` | 复用 `post()`（`:349`）的两态路由，不新造跨线程机制 | `residentPanel→PanelManager` 已存在（`:109`）；`agentRuntime→residentPanel` **不存在**，须装配根注入函数值（正斜方向） | 与仓库既有注入形状同构（`withPanelHost`、票 246 form 乙），先例最多 |
| D | `PanelManager` 加第 9 枚导出方法承接字节 | 需新增库调用（`Eval`/`Init`） | 无新包边 | 若走「页侧轮询取包」则必须新增绑定名 ⇒ 撞 §5① C17 白名单加名＝人工批准 |
| E | `internal/panel` 直接持有 webview | — | `internal/panel → 库`，且被 `cmd/wisp` 依赖 ⇒ **把 WebView2 拉进库层** | ⛔ **不是「不合适」，是点名禁过**：`panel_host_windows.go:180-183`（票 255 AC#4 forbids a panel->config dependency edge）、`resident_windows.go:113-125`（票 246 裁定 form 乙＋票 238 cut-1「两枚正向依赖边一律不开、改注入」；尺＝`GOOS=windows go list -deps ./cmd/wisp` 前后对拉） |

**依赖边方向的规矩（具名出处）**：装配根 `run.go:382 assembleRuntime`（`resident_task_source_windows.go:36` 称它
「the only assembly root in the repository」）负责**把函数值递下去**；先例逐字 `resident_windows.go:134-139`：
「It is BUILT here, by the assembly root, and handed to the ball host as a function value … so `resident_ball_windows.go`
keeps knowing nothing about WebView2 and this file keeps knowing nothing about the approval UI.」
⇒ `cmd/wisp` 包内新增一条 `agentRuntime → residentPanel` 的注入**不算跨包反向边**；
`internal/panel → cmd/wisp` 与 `internal/panel → webview2` 那两形才是禁的（tools→panel 同族）。

## S3 线程／再入这一族（按现状读码，不按注释读）

### 3.1 今天谁创建／拥有那条 STA 线程

- **就是 `cmd/wisp/panel_resident_windows.go`**，票头逐字（`:5-6`）：「The resident panel thread: the ONE place in this process that owns
  the panel window's OS thread (ticket 33, orchestrator ruling P1 at 10-01 13:12)」。
- 持有者 `residentPanel`（`:108` 定义，`:109 mgr *PanelManager`）；`:144 startResidentPanel` 建线程不建窗；`loop` 在 `:259`；
  `:297-300` 是「库泵返回即 teardown」；`handOverPump()`（`:322`）把 `mgr.currentWindow()` 存进 `rp.wRef`。
- ⚠ `PanelManager` **自己没有线程**：`bringUp`（`panel_host_windows.go:318`）跑在**调用方**线程上；
  文件头 `:138-145` 具名承认本文件「deliberately contains no runtime.LockOSThread of its own」，锁线程只发生在上面那条线程。

### 3.2 推送那一跳从泵（另一条 goroutine）发起会撞上什么——逐条带证据

1. **库侧硬事实**：`Dispatch` 只做两件事——`w.m` 下 append 进私有 `dispatchq`，再 `PostThreadMessageW(mainthread, WM_APP)`（`webview.go:443-447`）；
   **唯一抽队列处在 `Run()` 内**（`webview.go:349` 起，`:360-366`）。
2. **本仓已量过这条**（非推测）：`panel_resident_windows.go:21-28` 逐字——
   「`webview.Dispatch` only appends a closure to the library's private queue and posts a thread message, and **Run() is that queue's ONLY reader,
   so a host that pumps by hand can never deliver a Go -> page reply**: 33-p1 §A measured the page's awaited binding replies staying unresolved with the
   hand pump and **arriving in ~107ms with Run()**.」
   ⇒ 好消息：常驻线程**今天确实把泵交给了 `Run()`**（`handOverPump`），这扇门上着；
   ⇒ 约束：若哪天回到手写泵（`pnlPeekMessageW`/`pnlDispatchMsgW`，`panel_host_windows.go:104-105`），**Go→页整条会静默失效**，票面必须钉住。
3. **⛔ 不得投到 `ui-sta`**：同文件 `:8-19` 记着为什么——裁定 J1 先选「post it to the existing ui-sta」，`33-r1` 量到坏事；
   带正控复量（probe 33-p1，R34/R35＋正控 R36）：从一条**已在泵的线程**的窗口回调里 bring-up **不 panic**、约半秒返回，
   但它**在库的嵌套泵里时外层泵迭代计数器冻结** ⇒「up to five tasks that were already queued on that thread get executed early by the nested pump.
   That is **order reversal on the ball's own gesture queue**, which is worse to ship than a crash.」panic 那支需强制 `WM_QUIT`（R37）或 MTA（R32/R39）。
   ⇒ **面板归专用线程＋库 `Run()`；球泵永不重入。** 出向那一跳必须落在面板线程，不能搭球那条。
4. **`bringUp` 的线程前置**（`panel_host_windows.go:327-339`，33-r7）：创建窗口那条线程上不能有未派发的 `WM_CLOSE`。
   实测三种清理形状：只摘 `WM_QUIT`（33-r6 窄形）＝**PANIC 3/3**；先派发 close 再跑 quit 抽干＝**PANIC 3/3**；整队 purge 不派发＝前一枚窗口存活、线程窗口数 3→4（无界泄漏）。
   ⇒ 出向实现若顺手加「推之前先清队列」，这一族是它自己会撞的第一枚雷。
5. **再入形状（⚠ 读码推的，未真机量过）**：`Eval` 闭包在 `Run()` 抽队时执行，且队列已在锁外逐项调用（`webview.go:361-366`）；
   库自己的绑定回执同样走 `Dispatch→Eval`（`webview.go:148/152/156`）⇒ 从绑定回调里再 `Dispatch` **不构成自死锁**；
   真正要避的是**在 `Eval` 闭包内阻塞等待**——它会占住那条唯一能抽队的线程。派单若以此为前提，应要求落地腿实测。
6. **`Eval` 的两条硬伤要写进票面**（模块源码）：① 载荷含 NUL ⇒ `UTF16PtrFromString` 失败 ⇒ **`log.Fatal`＝进程退出**（`chromium.go:140-142`）；
   ② 结果丢弃（`chromium.go:144` 第三实参写死 `0`）⇒ **送达与否 Go 永远不知道**，判据只能来自页面（与探针同结论）。

### 3.3 仓里可抄的「跨线程投球」先例（点名一枚）

**`cmd/wisp/panel_resident_windows.go:349-372 func (rp *residentPanel) post(fn func()) bool`** —— 目前唯一一枚生产在用的
「非面板线程 → 面板线程」投递器，两态路由逐字：
```go
rp.mu.Lock()
if rp.wRef != nil { w := rp.wRef; rp.mu.Unlock(); w.Dispatch(fn); return true }
select { case rp.tasks <- fn: rp.mu.Unlock(); return true
         default: rp.mu.Unlock(); rp.failedPost.Add(1); return false }
```
- `wRef` 由 `handOverPump()`（`:322-346`）在交泵那一刻从 `mgr.currentWindow()` 取；注释 `:318-321` 写明这一排序使路由**无缝**：
  「post() either sees wRef empty and lands in the channel this drain reads, or sees wRef set and lands in the queue the pump reads.」
- **失败不假装**：`false` ＝「这条线程不会跑你的请求」（队列满或正在退出），`rp.failedPost` 计数，调用方须如实上报（`:346-348`）。
- 生产调用方：`RequestShow:378` / `RequestToggle:427` / `RequestDispose:442`（都是「把 Show/Hide/Destroy 搬上线程」那一族，见 `:408 showOnThread`）
  ⇒ **形状可抄，载荷还没人用过**（没有一枚调用者送过数据）。
- teardown 侧：`:495` 「the published handle so nothing can Dispatch into a dead pump」——推送路径要与这条退出口对拍，否则会有「往死泵 Dispatch」的窗口。

## S4 页面接收器有哪几种合法形状

⛔ 全程走 git 对象层（`git show dsh/feat/frontend-p0-v2:<path>`／`git grep … -- 'frontend/src'`，ref 解析＝`16c2f038`）；
**一步都没进 `D:/wt/fe`，仓内 `frontend/**` 一个字节都没落**（终态自查见 S6）。
⛔ 负判决一律用目录形式 `-- 'frontend/src'`，没用那支会跳过 `App.tsx` 的 `*.tsx` glob。

### 4.1 `frontend/src/lib/panel.ts` 今天暴露的桥面（311 行）

**桥面只有一个方法，且方向是页→Go**：
```ts
140: interface WispHostBridge {
141:   postMessage(message: string): void;
142: }
144: declare global {
145:   interface Window {
146:     /** Installed by WebView2's AddHostObjectToScript / postMessage pipe. */
147:     wispBridge?: WispHostBridge;
148:     chrome?: { webview?: WispHostBridge };
152: function hostBridge(): WispHostBridge | null {
154:   return window.wispBridge ?? window.chrome?.webview ?? null;
```
- `window.wispBridge` 与 `window.chrome?.webview` **两枚都只声明了同一枚方法 `postMessage(message: string)`**（`:141`）；
  取值顺序 `wispBridge ?? chrome?.webview ?? null`（`:154`）。
- 页侧出口函数共 6 枚，全部走 `bridge.postMessage(`（`:179`、`:218`）：
  `requestApprovalResolution:169`、`requestModeSwitch:245`、`requestWorkspaceChange:254`、`sendAttachmentBytes:266`、
  `submitAttachment:281`、`sendMessage:299`；外加 `hostedByNative:158`（＝`hostBridge() !== null`）、`sendRequest:211`、`nextRequestId:205`。
- 类型侧的快照形状：`PanelSnapshot:129-137`（`pending / results / composer / generatedAt`），
  注释 `:128` 逐字：**「What the host pushes on every state change (ticket 35 owns the pump).」**
  ⇒ 页面的类型模型**从一开始就假定有一枚 host→页的 push**，缺的从来不是类型。

### 4.2 WebView2 页侧接收 Go 推送的合法写法——以仓里真出现过的写法为据

⛔ 不凭记忆编 API。仓里/探针里**真出现过**的只有下面三形，第四形（`chrome.webview.addEventListener('message')`）**在本仓零证据**：

| 形 | 本仓的证据 | 状态 |
|---|---|---|
| 甲 `Eval` 写页面全局状态/DOM，页面**轮询**自己读到 | 探针 `-mode eval`／`eval-disp`：Go `Dispatch(Eval("document.title='…'"))`（`probes/33/p1/q2/main.go:283`、`:302`），页侧 `setInterval` 读 `document.title` 命中即 `wispReport("EVAL_SEEN")`（`:378-393`） | **实测收到过**（票头 `:21-25` 明确以页面的话为判据） |
| 乙 `Eval` 调页面上的**全局函数**（`window.__x = …`／`fn(json)`） | 库自己的绑定回执就是这一形：`webview.go:156 Eval("window._rpc["+id+"].resolve("+string(b)+"); …")` | **生产在跑**（每次绑定调用都走），但载荷是绑定回执 |
| 丙 `Init(js)` 在每次文档创建前注入一段脚本 | 库导出（`godoc :64`；实现 `chromium.go:130-136 AddScriptToExecuteOnDocumentCreated`），且 `Bind` 本身就靠 `Init` 装 `window[name]`（`webview.go:461` 起） | 通道存在，**本仓无生产 push 用过它** |
| 丁 `chrome.webview.addEventListener('message')` ＋ Go 侧 `PostWebMessageAsJSON` | **零证据**：`PostWebMessage*` 在库里只是 vtbl 表项（`corewebview2.go:106-107`），无导出方法（票头 `:16-19` 已具名）；页侧 `addEventListener('message')`／`onmessage`／`MessageEvent` 于该 ref 的 `frontend/src` **共 0 命中** | **取不到仓内写法** ⇒ 不是「不行」，是**本仓没有先例、也没有可抄的正控**，选它＝从零验证 |

⇒ 甲/乙/丙三形**在 Go 侧都是同一扇门**（`Eval`，区别只在页那头怎么读）；它们都要求 §3.2 那条：`Run()` 在抽队。

### 4.3 页面要把收到的 JSON 灌进 `App({snapshot})`，缺的是哪一环

**缺的那一环点名到文件:行：`frontend/src/main.tsx:97`**
```
 97:   product: { area: "面板", node: <App /> },
110:   createRoot(document.getElementById("root")!).render(
```
- `<App />` **不带任何 prop** 被挂载 ⇒ `App.tsx:91 export default function App({ snapshot = EMPTY }: { snapshot?: PanelSnapshot })`
  的默认值 **`EMPTY`（`App.tsx:70-73` 定义）就是页面永远看到的唯一值**。
  ⇒ **`snapshot` 的上游不是任何人**：`main.tsx:97` 是挂载点，`App.tsx:70/91` 是兜底，中间**没有 React state、没有 receiver、没有 store**。
- 而 `main.tsx:12-13` 的注释逐字写着「inside Wisp, WebView2 serves the embedded dist and **Go pushes snapshots through the C17 bridge** -
  no query string, so this file renders `<App />`」⇒ 页面作者**以为**这条 push 存在，实际挂载点没有 prop。
- ⇒ 补齐「那一环」的最小清单（**不是选形，是缺件名册**）：①一枚 receiver（4.2 甲/乙/丙 任一的页侧读法），
  ②一枚把 receiver 的值交给 React 的 state/订阅（`main.tsx` 或 `App.tsx` 内 `useState`＋`useEffect`），
  ③`<App snapshot={…}/>` 的实参。今天 **①②③ 三件全无**。

### 4.4 ⚠ 本节最该上报的一条：入向（页→Go）也不是「真在」，它只在对 Go 半截成立

- Go 半截（与转述一致）：`panel_host_windows.go:80 panelDispatchBinding = "wispDispatch"`、`:405 w.Bind(panelDispatchBinding, …)` → `:630 dispatchRaw` → `internal/panel/composer_dispatch.go:181 HandleModeRequest(`。✅
- **页半截不成立**：页面发的唯一形状是 `bridge.postMessage(JSON.stringify({method:"panel.mode.request", …}))`（`panel.ts:179-185`、`:218`），
  走的是 **WebView2 原生 `chrome.webview.postMessage`**；而库把所有收到的消息一律按 **RPC 信封**解析
  （`webview.go:131-135 type rpcMessage{ID;"id"; Method;"method"; Params;"params"}` → `:139 msgcb` → `:162 callbinding` 查 `w.bindings[d.Method]`）。
  绑定名册（Go 生产只有两枚）＝ `wispDispatch`、`wispProbeRT` ⇒ `bindings["panel.mode.request"]` 查不到 ⇒ `callbinding` 回 `(nil, nil)`，
  随后 `:156` 会对不存在的 `window._rpc[0]` 做 `.resolve`（TypeError 被 Eval 吞掉）。
- **尺与正控**：`git grep -E "wispDispatch|wispProbeRT" dsh/feat/frontend-p0-v2 -- 'frontend/src'` ＝ **0 命中**；
  同一把尺在 Go 树命中 5＋ 处（`panel_host_windows.go:80`、`panel_resident_windows.go:44`、`panel_resident_windows_test.go:199/301/801`，
  其中 `:301` 与 `:801` 是产品自己的 AC#13/AC#14 探针**真调过** `window.wispDispatch(`）⇒ 尺能命中真东西，0 不是尺坏。
- 另：dev 工作树那份 `frontend/src` 拷贝同尺同结果（`frontend/src/lib/panel.ts:147` 只有 `wispBridge`，全树 0 处 `wispDispatch`）。
- `window.wispBridge` 这个名字在 Go 树里**没有任何一处安装它**（唯一出现是测试夹具 `internal/panel/composer_test.go:101`）。
- ⇒ **精确表述**：Go→页缺；页→Go 是**两半各缺一块**——Go 那半的门开在 `window.wispDispatch(…)` 上，
  页那半从来不叫这个名字，而是把信封直接 `postMessage` 给库的 RPC 解析器。**这座桥两头都没接上，只是 Go 这头有半根线。**

## S5 禁区与代价表（可派单形状）

### 5.1 四条硬禁区（逐条给出处）

1. **⛔ 不新造 `C##`、不往 C17 方法白名单加名字**。白名单现量（`internal/panel/bridge.go:41-45` 的 const 块）＝ 6 枚：
   `MethodModeRequest "panel.mode.request"`(:42)、`MethodWorkspaceRequest "panel.workspace.request"`(:43)、
   `MethodAttachmentAdd "panel.attachment.add"`(:44)、`MethodMessageSend "panel.message.send"`(:45)、
   外加 `bridge.go:46-49` 具名的票 248 两枚 `MethodConfigGet/MethodConfigSet`（**故意不带 `panel.` 前缀**，`:54-56` 写着这「load-bearing twice over」，
   且 `git_test.go` 的 `whitelistMethodsFromSource` 在数 `panel.*` 字面量 ⇒ 加名＝同时踩白名单与反漂移钉）。
   加名＝人工批准。**注**：本节 4.4 那条「页该改叫 `wispDispatch`」的形状**不加白名单名**，但会新增一枚 JS 全局名，派单要分清这两件事。
2. **⛔ 面板只是「显示＋发起请求」的口，不许经它给权限/批准**。出处：`panel_resident_windows.go:41-47`
   「it grants nothing: the only door it exposes to the page is the one … whose widening branches are fail-closed, and it adds no approval channel of its own」。
3. **⛔ 宁缺毋造**：`App.tsx:214-215` 的注释声称「Reads its strings from the snapshot」，紧接 `:216-219` 的实参是四个字面量
   （`view={{ current: "", isRepo: false, local: [], worktrees: [] }}`、`onCheckout={() => undefined}`、`onNewFrame` 同形）——**转述这条读数在盘上逐字成立，本腿复核通过**。
   这座桥落地时同族那一形不许再被引入一次：拿不到真值的栏位只许画「没有」。
4. **⛔ 页侧「记住一份答案」被明令禁止**（这一条是我在 4.1 读到的，不是转述）：`panel.ts:162-168` 逐字
   「responses arrive as a fresh PanelSnapshot push, **never as a return value**, because a panel that keeps a copy of the answer would be a second state holder (PLAN.md:1044)」
   ⇒ **S2.3 表里「页侧轮询绑定函数、从 Promise resolve 里读 JSON」那一形（我标为 A/D 的候选）被这行票面直接否掉**；
   它技术上今天唯一跑得通，但要动 `panel.ts` 的既有契约文字 ⇒ 属**契约面**，不是实现细节。

### 5.2 逐枚落点的代价表（选形归编排者）

| 形 | 动的文件 | 新增依赖边 | 碰 STA 约束？ | 要解冻的既有钉／裁定 |
|---|---|---|---|---|
| **push 走 `Dispatch→Eval`**（复用 `post()`） | `cmd/wisp/panel_resident_windows.go`（新增方法）、`cmd/wisp/panel_pump.go:319`（记账后调它）、`cmd/wisp/run.go:734` 或 `:607/:725` 附近（把句柄递进 `agentRuntime`） | 包内注入，**无跨包边** | **是**：必须落 `post()` 的两态路由＋`wRef`（`S3.3`），且 `Eval` 闭包内不许阻塞 | `run.go:730-733` 那段「不给常驻腿第二条 publish 路」的裁定 |
| **push 走 `Init` 预置回调**（Go 只 `Eval("window.__wispPush(json)")`） | 同上，外加页侧一处 receiver | 无 | 是（同上） | 无契约钉；但页侧新增全局名 |
| **页→Go 改叫 `window.wispDispatch(env)`**（修 4.4） | 只动 `frontend/src/lib/panel.ts`（`hostBridge`/`sendRequest`） | 无 | **否**（入向不经新线程） | 动的是页面契约文字；`panel.ts:146` 那句「AddHostObjectToScript」是过期描述 |
| **页侧轮询取包（绑定返回快照）** | Go 加第 9 枚绑定 ＋ 页侧 `useEffect` 轮询 | 无 | 否（走已有绑定回执） | ⛔ **撞 5.1-④**（`panel.ts:162-168` 明令「never as a return value」）＋ 可能需新绑定名 ⇒ 撞 5.1-① |
| **`internal/panel` 直接握 webview** | — | ⛔ `internal/panel → 库` | — | ⛔ 点名禁过（票 255 AC#4 / 票 246 form 乙 / 票 238 cut-1），见 S2.3-E |
| **新增 `PostWebMessageAsJSON` 通道** | 需在 Wisp 侧直调 COM vtbl，或换/patch 依赖 | 无 Go 包边，但**新增仓内 COM unsafe 代码** | 是 | ⛔ 会改 `go.mod` 依赖面（`go list -deps` 前后对拉那把尺会被触发）；且页侧**零先例**（4.2-丁） |

**共同代价（不分形）**：载荷含 NUL ⇒ `log.Fatal` 退进程（`chromium.go:140-142`）；`Eval` 无回执 ⇒ 判据必须来自页面（探针 R25 教训）；
`dispatchRaw` 的返回字符串今天**已经**经 `:156` 那条形回页 ⇒ 任何「页侧自己再造一条回执通道」都要与这条对拍，不要造出两份状态。

## S6 终态自证

见本件末尾「S6 输出原文」小节（与 `git status`／`git show --stat HEAD` 实测同批写入）。

## 附：本腿推翻／收窄了派单里的哪几句

| # | 派单原话 | 盘上实测 | 严重度 |
|---|---|---|---|
| 1 | 「`PostWebMessage\|EvaluateScript\|CreateWebMessageAsJson` 于 `cmd/wisp`＋`internal/panel` ＝ 0」 | 读数对，**尺为负**：这三个名字里 `EvaluateScript`/`CreateWebMessageAsJson` **根本不是这套 API 的拼写**，`PostWebMessage*` 只在库里作 COM vtbl 表项存在、无导出方法。库的真出向名册是 `Eval`/`Init`/`SetHtml` ＋跨线程门 `Dispatch` | 高——按原尺永远得不出「有/没有」的判决 |
| 2 | 「出向（Go→页）**整条不存在**」 | **不准确**：库在每次绑定调用后都用 `Dispatch(Eval("window._rpc[i].resolve(…)"))` 把 Go 的返回值送回页面（`webview.go:148/152/156`），本仓生产每次面板交互都走它；`SetHtml` 更有 3 处生产调用。**缺的是 unsolicited push（Go 不经页问就先说）那一支，不是「Go→页这根线」** | 高——决定了落地腿是「新建通道」还是「给已有通道加发起端」 |
| 3 | 「入向（页→Go）真在」 | **只对 Go 半截成立**：门在 `window.wispDispatch(…)`（`panel_host_windows.go:80`），页面从来不叫这个名字，它发 `chrome.webview.postMessage({method:"panel.mode.request",…})`（`panel.ts:179-185`），而库把一切收到的消息按 RPC 信封解析并查绑定表 ⇒ 查不到 ⇒ 空转。`git grep wispDispatch -- 'frontend/src'`＝0（正控：Go 测试里 2 处真调过） | **本腿最值钱**——两头都没接上 |
| 4 | 「`func (m *PanelManager)` 导出方法共 8 枚」 | 复核**通过**：`IsCreated/IsShown/LastColdMs/LastHotMs/Show/HotShow/Hide/Destroy` 逐枚对上（`:268/276/285/286/481/547/574/615`）；非导出另有 11 枚声明行 | 低（确认） |
| 5 | 「`lastPanelSnapshot()` 带括号数＝10 处调用全在 `_test.go`、非 test 0 处」 | 复核**通过**并加一条：非 test 唯一命中是 `panel_pump.go:383` 的**定义行本身**，不是调用 | 低（确认） |
| 6 | 「泵在跑、只差 transport」 | 更硬的一层：`run.go:607 rt.ui = nil`（常驻腿递了 gate）＋ `run.go:734 if rt.ui != nil` ⇒ **持有 `PanelManager` 的那个进程今天一次都不调 `publishPanelSnapshot`**；而调得到它的控制台腿里 `PanelManager` 没被构造。`run.go:730-733` 还具名拒绝过补这条触发器 | 高——这是「隔着什么」的真答案 |
| 7 | 「`App.tsx:216` 注释声称读快照、实参是四个字面量」 | 复核**通过**，逐字（`:214-215` 注释 vs `:216-219` `view={{ current:"", isRepo:false, local:[], worktrees:[] }}`） | — |
| 8 | （派单未问，但会挡住选形） | `panel.ts:162-168` 已把「响应只能作为新快照 push 到达、**绝不用返回值**」写成契约（引 `PLAN.md:1044`）⇒ 「页侧轮询绑定函数读 Promise resolve」这一形**技术上今天唯一通、契约上被明令否掉** | 高 |
