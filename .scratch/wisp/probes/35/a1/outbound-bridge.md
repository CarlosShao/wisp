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
