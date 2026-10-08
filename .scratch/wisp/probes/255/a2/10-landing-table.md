# 255-a2 落点表（六问）— 「改面板宽度＝不用重启立刻见效」

- 腿：`255-a2`（只读前置普查）。起手 HEAD `5738661`（地板 `914177e6` 之后，按现量记）。
- 本程零 `go build/vet/test`；零产码写点（唯一写面＝本目录＋工单 255 追加一节）。
- ⛔ 本文件不裁形、不勾框、不批契约。

---

## 问① 窗口句柄与库面

**HWND 从哪来、存哪、何时有效**（尺：`grep -n "hwnd\|HWND" cmd/wisp/panel_host_windows.go` ⇒ 21 命中，rc=0）

| 环节 | 位置 | 现量 |
|---|---|---|
| 取得 | `cmd/wisp/panel_host_windows.go:398` | `hwnd := windows.HWND(w.Window())`（create 成功后立刻取，库侧 `Window()` 见下） |
| 落库 | `:415` | `m.hwnd = hwnd`（与 `m.w = w`、`m.created = true` 同一段锁内） |
| 字段 | `:149` | `hwnd windows.HWND`（`PanelManager` 内，由 `:147 mu sync.Mutex` 护） |
| 清零 | `:637` | Destroy 里 `m.hwnd = 0` ⇒ **有效性窗口＝bringUp 成功之后到 Destroy 之前** |
| 对外读口 | `:290 windowHandle() uintptr` | 「0 before creation」，注释自述给仪器读进程窗口态用，⛔ 不经 WebView2 |
| 控件句柄读口 | `:302 currentWindow() webview2.WebView` | 注释写明它是常驻腿路由请求的依据：nil ⇒ 走线程自己的 channel，非 nil ⇒ 走库的 Dispatch 队列 |
| 另一枚出口 | `:879` | 返回 `true, uintptr(m.hwnd)`（status/diag 形状） |

**库侧有没有现成的改尺寸出口＝有，两枚，都在 `go.mod:19` 那一枚 `github.com/jchv/go-webview2 v0.0.0-20260205173254-56598839c808` 的 Go 包装层**（尺语法：先在 `common.go` 读**接口**的真名，再去 `webview.go` 读**实现**；⛔ 不是 COM vtable 名）：

```
$ grep -rn -A40 "^type WebView interface" <modcache>/github.com/jchv/go-webview2@v0.0.0-20260205173254-56598839c808/*.go
common.go:26  type WebView interface
common.go:39  	Dispatch(f func())
common.go:47  	Window() unsafe.Pointer
common.go:54  	SetSize(w int, h int, hint Hint)
rc=0
```

```
$ grep -rn "SetSize\|MoveWindow\|SetWindowPos\|Dispatch(" <同一枚库> --include=*.go（非 test）
webview.go:405  func (w *webview) SetSize(width int, height int, hints Hint)
webview.go:427  	_, _, _ = w32.User32AdjustWindowRect.Call(...)
webview.go:428  	_, _, _ = w32.User32SetWindowPos.Call(w.hwnd, 0, r.Left, r.Top, r.Right-r.Left, r.Bottom-r.Top,
         		SWPNoZOrder|SWPNoActivate|SWPNoMove|SWPFrameChanged)
webview.go:431  	w.browser.Resize()
webview.go:443  func (w *webview) Dispatch(f func())   // :447 PostThreadMessageW(w.mainthread, WMApp, 0, 0)
rc=0
```

⇒ **`SetSize` 就是那枚出口**：它内部就是 `SetWindowPos`，并且**尾随 `browser.Resize()`**（把 WebView2 controller 跟着窗口客户区重算）；仓里对此的现量口径在 `cmd/wisp/panel_host_windows.go:41-49`（`(*edge.Chromium).Resize()` 是导出的、按自己的 `GetClientRect` 定尺寸、且高层包装在建窗成功后 `webview.go:343` 已经调过一次）。

⚠ **两条路的数值语义不一样，这条最重要**（尺：`sed -n '288,330p' <库>/webview.go`，rc=0）：
- **建窗**：`webview.go:296-305` 把 `opts.Width/Height` **当作 CreateWindowExW 的窗口外框尺寸**直接用（`:320`，style `0xCF0000 = WS_OVERLAPPEDWINDOW`），0 值兜底 640×480。
- **改尺寸**：`webview.go:405` 先把入参当**客户区**做 `AdjustWindowRect` 再 `SetWindowPos`。
⇒ 同一个数字走这两条路得到的**外框宽度不同**（差一枚边框＋padding）。今天 `windowOptions()` 喂的是前者；甲形若上 `SetSize`，**它的入参是后者**。仪器与文案都不许把两个 420 当同一个 420。

`Hint` 的射程（尺：`grep -rn -B2 -A8 "type Hint" <库>/*.go`，rc=0；`common.go:8-17`）：`HintNone=0`/`HintFixed`/…；实现里 **`HintMax`/`HintMin` 只写 `w.maxsz/w.minsz` 而根本不碰 `SetWindowPos`**（`webview.go:415-420`），只有其余值进 `:421-432` 的改尺寸分支；而 `HintFixed` 会顺手 `style &^= (WSThickFrame|WSMaximizeBox)`（`:409`）＝**把窗口变成用户拖不动**。⇒ 选 hint 不是免费的，本腿不选。

**宿主今天有没有改尺寸出口**（复核编排者的转述，结论一致但行号已漂）：`cmd/wisp/config_readers_255.go:157` 注释逐字仍在写 `resizes": this host contains no MoveWindow / SetWindowPos / SetBounds`，`panel_host_windows.go:234` 注释同样写「this host has no MoveWindow/SetWindowPos/SetBounds call at all」。尺：`grep -rn "panelWidthPx|panelHeightPx|MoveWindow|SetWindowPos|SetBounds" cmd/ internal/`（含 test）⇒ `cmd/wisp/panel_host_windows.go` **零枚真调用**，只有 `:234` 那句注释；全仓 `SetWindowPos` 的真调用都在别处（`internal/ball/win32_windows.go:38`＋`ball_windows.go:600`／`monitors_windows.go:97`、`cmd/balldebug/shot_windows.go:35/:145`、`cmd/wisp/testdata/esclistener/main.go:69/:229`、库内 `webview.go:428`）。rc=0

---

## 问② 像素那套换算在哪

| 符号 | 位置 | 类型全名／调用形状 |
|---|---|---|
| `panelWidthPx = 420` / `panelHeightPx = 260` | `panel_host_windows.go:89-90`（const 块 `:79-91`） | 语义已被 AC#4 改成「没人给几何源时的缺省」（`:81-88` 注释逐字） |
| 读者 1 | `:237` | `width, height := panelWidthPx, panelHeightPx`（`windowOptions` 起点） |
| 读者 2 | `panel_resident_windows.go:204` | `slog.Warn(... "default", fmt.Sprintf("%dx%d", panelWidthPx, panelHeightPx))`（**只在读盘失败的分支**） |
| 几何字段 | `panel_host_windows.go:184` | `geometry func() (width, height int)`，注释明写 nil⇒用常量、构造后不再改所以不加锁 |
| 唯一消费点 | `:244` | `gw, gh := m.geometry()` ⇒ `if gw > 0 { width = gw }` / `if gh > 0 { height = gh }`（`:245-258`，`<=0` 一律回落到常量；`gh==0`＝「auto from content 没实现」那格，`:248-255` 具名登记） |
| 注入口 | `:199-201` | `func withGeometrySource(src func() (width, height int)) panelHostOption`；`panelHostOption` 定义 `:194`（变参 option，`:187-193` 自述理由＝7 处旧测试构造点不改） |
| 装配根调用形状 | `panel_resident_windows.go:253` | `NewPanelManager(disp, assets, dataPath, withGeometrySource(panelGeometrySource(dataDir)))` |
| 几何源本身 | `panel_resident_windows.go:200-209` | `func panelGeometrySource(dataDir string) func() (width, height int)` ⇒ 内体 `cfg, _, err := config.LoadFile(cfgPath, nil)`（`:202`）⇒ `return cfg.Panel.Width, cfg.Panel.Height`（`:208`）；失败 `return 0, 0`＋Warn |
| schema | `internal/config/schema.go:539-549` | `PanelSection`：`Width int` default `640`（`:542`）／`Height int`「0 = auto from content」（`:543-544`）／`Scale float64`「0 = follow system DPI」（`:547-548`） |

⚠ **顶回编排者转述的一处：DPI 那层换算今天根本不存在，没有任何 Go 侧代码做逻辑像素→物理像素。** 两把尺：

```
$ grep -rn "DPI|dpi|GetDpiForWindow|scaleFactor|LogicalToPhysical|dpiscale" --include=*.go cmd/ internal/config/ internal/panel/（排 *_test.go）
internal/config/schema.go:547  // Scale is the UI zoom factor; 0 = follow system DPI.   ← 唯一命中，且是注释
rc=0

$ grep -rn "\.Scale\b" --include=*.go cmd/ internal/（排 *_test.go）
internal/agent/budgets.go:94   s := func(ref int) int { return scaleInt(ref, b.Scale) }   ← 别个类型的 Scale，非 PanelSection
rc=0
```

⇒ `PanelSection.Scale` 的**非测试读者＝0 枚**；`cfg.Panel.Width` 是被**原样**塞进 `webview2.WindowOptions{Width: uint(width)}`（`:260-264`）。真实缩放由 Win32 的进程 DPI awareness 承担，仓里对它的唯一现量口径是**间接量出来的**：`cmd/wisp/panel_geometry_255_winlive_test.go:171-176` 用第一枚真窗反推 `scale := float64(got1)/float64(asked1)`，再按 `want2 := float64(asked2)*scale` 断言（文件头 `:29-34` 逐字讨论 DPI）。⇒ 甲形的落点表里，"谁做换算"这一格的正确回答是**「今天没人做，且量出来是 1.0 时才等于配置值」**，不是「在某某文件里做」。

---

## 问③ 线程归属

- **归属定案在码里的形状**：`panel_resident_windows.go:262 runtime.LockOSThread()`（在 `rp.loop` 内，全仓该文件唯一一枚）→ `:264` `CoInitializeEx(COINIT_APARTMENTTHREADED)` 非 STA 即拒绝建窗 → `:299 w.Run()`（库泵占住这条线程）→ `:300 rp.teardown(...)`。名字与 owner：`:82 panelSTAName = "panel-sta"`／`:83 panelSTAOwner = "panel host (ticket 33)"`，注册点 `:154 rp.thread = reg.Spawn(panelSTAName, panelSTAOwner, root, rp.loop)`。
- 宿主自己**不**锁线程（`panel_host_windows.go:138-145` 逐字：`this file deliberately contains no runtime.LockOSThread of its own - the rule is that every method touching the control (bringUp/Show/Hide/Destroy) runs on the thread that created the window`）。
- ⇒ **改尺寸的调用必须落在那条专用 STA 线程上**（与 `Show/Hide/Destroy` 同一条规则；`SetSize` 内含 `SetWindowLong`/`SetWindowPos`＋`browser.Resize()`，属「touching the control」）。

**现成的「回到那条线程」机制＝`residentPanel.post`，两路由，已存在，不许新造**：

```
$ grep -n "func (rp \*residentPanel) post" -A 22 cmd/wisp/panel_resident_windows.go
:347 func (rp *residentPanel) post(fn func()) bool
:355 	rp.mu.Lock(); :356 if rp.wRef != nil { ... :358 w.Dispatch(fn); return true }   ← 库队列路（Run() 是唯一读者）
:362 	case rp.tasks <- fn:                                                            ← 建窗前路（chan func() 容量 16，:150）
:366 	default: rp.failedPost.Add(1); return false
rc=0
```

它的既有用户（同一形状的现例，全部只递闭包、不等结果）：`:450 rp.post(func() { rp.mgr.Destroy() })`（`RequestDispose`）、`:482 rp.post(func() { rp.mgr.terminateOnThisThread() })`（`stop`），注释 `:344-346` 逐字「False means "this thread will not run your request" ... which callers report instead of pretending the window moved」。发布 `wRef` 的那一步是 `:331 handOverPump()`（注释 `:317-323` 明写「post() either sees wRef empty and lands in the channel this drain reads, or sees wRef set and lands in the queue the pump reads」）。库侧证据：`webview.go:443-447`（`Dispatch` 追加 dispatchq ＋ `PostThreadMessageW(w.mainthread, WMApp)`），仓内口径 `panel_host_windows.go:35-39`／`panel_resident_windows.go:21-25`。

**「在别的线程调 `SetWindowPos` 会是什么形状的后果」＝仓里没有凭据。** 尺：`grep -rn "SetWindowPos" --include=*.md docs/ .scratch/wisp/probes/33 .scratch/wisp/issues/` ⇒ 5 命中（rc=0），逐枚都是**「宿主没有 SetWindowPos」那一类断言**（`docs/evidence/s1/180-182-panel-fields-census-c1.md:50`、`docs/reports/pending-and-issues.md:11104`、票 180 `:138/:187`、票 255 `:57`），**没有一枚是跨线程调用它的实测**。⇒ ⛔ 不许拿「user32 跨线程安全/不安全」的常识冒充实测；本腿只登记这一格为**无凭据**。最接近的一条实测**是另一枚 API**，别混：`panel_resident_windows.go:474-481` 具名写的是 `Terminate` 的裸 `PostQuitMessage`「posts WM_QUIT to the queue of the thread that CALLS it」，那才是本机量过的跨线程事故。

---

## 问④ 热加载那一跳挂不挂得上

**现量的热加载驱动**（尺：`grep -rn "startConfigReload|reloadOnce" --include=*.go cmd/wisp/*.go` 排 test ⇒ 4 命中，rc=0）
- `cmd/wisp/config_reload.go:105 func (rt *agentRuntime) startConfigReload()`（1s `time.Ticker` ⇒ `:153 rep, err := rt.mgr.CheckAndReload()`，`reloadOnce` 在 `:152`）。
- 唯一非测试调用者：**`cmd/wisp/run.go:813 rt.startConfigReload()`**。
- ⚠ **常驻腿不 tick**：`grep -rn "startConfigReload|reloadOnce" cmd/wisp/resident*.go` ⇒ **0 命中（rc=1，即空）**。这条是票 255 `:163`／`residentPanelHotReloadNote`（`panel_resident_windows.go:163`）与 `residentPanelGeometryNote`（`:177`，正文 `:178-179` 逐字「本腿不轮询 config.toml…已建好的窗口也不会自己改大小（今天没有 resize 路）」）已经写下的裁语，本腿只复核形状、不改它。

**有没有现成回调能让面板知道「[panel] 那一段变了」＝没有，而且是结构性没有**（三把尺）：
1. `internal/config/manager.go:54 OnReload func(sections []string)`，触发条件 `:198 if len(rep.Reload) > 0 && m.OnReload != nil` ⇒ **只对 Reload 档开火**。
2. `[panel]` 登记在 **hot**：`internal/config/tiers.go:34 "panel": "hot"` ⇒ 它进的是 `rep.Hot`，**永远不进 `rep.Reload`** ⇒ `OnReload` 这条路对 `[panel]` 天生不开火（`internal/ball/hotkey_reload.go:26-30` 已经把同一件事写成了口径：「Why the bridge polls the section instead of riding config.Manager.OnReload: [hotkey] is HOT-tier … OnReload only fires for RELOAD-tier」）。
3. `OnReload` 的非测试赋值点全仓 1 枚：`cmd/balldebug/main.go:244 mgr.OnReload = bridge.OnReload()`（票 255 禁区 `:26` 已把这格派给票 42）。

⇒ **挂不挂得上：挂不上现成的；`[panel]` 那一档既不在 OnReload 射程、常驻腿也不跑那个 tick。** 面板今天拿到新值的唯一一跳是**自己那枚 per-create 现读闭包**（`panelGeometrySource`，`panel_resident_windows.go:200-209`），它对「已建好的窗」无能为力。

**最少要动的三处（只指认，不选形）**：
- (a) `cmd/wisp/panel_host_windows.go`：一枚「读 `m.geometry()` 一次 ⇒ 对活控件调 `SetSize`」的方法，落点就在 `:236 windowOptions` 与 `:302 currentWindow` 之间这一段面上；⛔ 不许在这枚文件里读盘（`:177-184` 是 AC#4 自己立的规矩）。
- (b) `cmd/wisp/panel_resident_windows.go`：一枚 `rp.post(func(){ ... (a) ... })` 形状的请求，落点紧邻 `:439-450 RequestDispose`（同一族，同一把「false 就说出来」的话）。
- (c) **触发者**——这才是真正的一跳，也是本腿认为最贵的一格：今天没有任何人在 `[panel] 变化`那一刻去按 (b)。候选只有两枚，**都在别人的写面上**：`cmd/wisp/config_reload.go:105/:153` 那枚 tick（它现在只在 `run.go:813` 起来，常驻进程 0 命中），或 `cmd/wisp/resident_windows.go:152-158` 的装配段（面板线程在这儿被 `startResidentPanel` 立起）。⚠ 把 tick 引进常驻进程 = 动 223/258 那族刚落的 `hotReload258`（`cmd/wisp/resident_windows.go:205`，注释 `:201-202` 逐字写「the exact shape panelGeometrySource has been since ticket 255 AC#4」＝**那条腿刻意选了现读闭包而不是 tick**）⇒ 这一处会碰到别人正在写的面，须编排者排程，不是腿自己并进。
- ⛔ **同批撞钉预检（票 255 `:58` 那条规矩的现场版）**：`cmd/wisp/config_readers_255.go:161` 的判词今天逐字写着「…（`cmd/wisp/panel_host_windows.go:262` [Width: uint(width)], reached from the create at `:392`）**面板关窗再开即跟上新值**…」，而 `config_readers_255.go:17/:141` 与 `cmd/wisp/config_receipt_255_test.go:440` 还在 cite 已漂的 `panel_host_windows.go:304`。⇒ **甲形落地那一次必须同批改这三处文案/引用**，否则「关窗再开即跟上」会变成一句新的吹牛（它已经在本票 10-08 补格 `:67-71` 被抓过一次）。

---

## 问⑤ 仪器能钉到哪一格（无窗那半）

**能钉死（CI 有载体，零真窗）**：

1. **「发出了一次改尺寸调用」这一形状读得到**，且**不需要新造接缝**：`cmd/wisp/panel_pageover_33r10_windows_test.go` 里的假控件**已经实现了这两枚方法**——`:123 func (c *docSink33r10) Dispatch(f func())`、`:130 func (c *docSink33r10) SetSize(w, h int, hint webview2.Hint) {}`（现量 `grep -n "func (.*docSink33r10)"` ⇒ 14 枚方法，rc=0；全接口已满足，注入点 `:154 m.w = c`）。⇒ 把 `:130` 那枚空体改成**记录 (w,h,hint) 的 sink**，一台无窗机器就能问：宿主到底发没发、发的是不是几何源那一次的答案、hint 是哪一枚。⚠ 另一枚同族假控件在 `cmd/wisp/panel_transport_live_35v2_windows_test.go`（`grep -rln "webview2.WebView\b" --include=*_test.go cmd/wisp` ⇒ 2 命中，rc=0），加方法时它是第二处要看的。
2. **能力尺写法照 33-r10 那一族**（文件头 `:19-33` 逐字：断言问的是「最终文档的能力」，⛔ 不问 `SetHtml` 调用枚数，⛔ 不锚行号，⛔ 不锚字面量）。同到改尺寸上＝**问「控件被 handed 的那一对数，等不等于几何源此刻答案（± 那枚 `AdjustWindowRect` 的边框差）」，别问「SetSize 被调了几次」**；问次数会被「合并成一次」这种形状骗过去。
3. **线程归属也能无窗钉**：甲形若走 `rp.post`，假控件的 `Dispatch` 就是现成的读数点——仪器可以断言「改尺寸是作为闭包**经 Dispatch 那一队列**到达控件的」，而不是在调用者线程上直接落到 `SetSize`（两枚记录点已经分好：`:123` vs `:130`）。
4. **既有那枚不许漂的钉**：`cmd/wisp/panel_geometry_255_test.go:227/:239/:266` 用 AST 解析本文件，要求**恰好 1 枚 `webview2.NewWithOptions`**、且 `WindowOptions:` 的值必须是 `windowOptions()` 这枚**调用**而不是字面量（`:279` 那条报错逐字：「the fallback must go through `panelWidthPx`/`panelHeightPx` so this file has one place the default lives」）。⇒ 甲形加 `SetSize` **不得**顺手加第二枚建窗路，这条红会自己响；`grep -n "geometry" panel_geometry_255_test.go` ⇒ 现量数建窗次数的断言在同文件（`:202-203` 自述四把尺的射程＝「数字到了 create」，⛔ 不含「Win32 真把窗开成那么宽」）。

**只能〔仅本机可量〕那一半（⛔ 不许报成 CI 测过）**：

- **真窗改完到底多宽**：只能问 HWND 几何。仓里同类先例已存在并且**只此一枚**属面板：`cmd/wisp/panel_geometry_255_winlive_test.go:48/:57`（`GetWindowRect`）＋ `:173` 那条 `t.Logf`（逐字打印 `asked → GetWindowRect … implied scale … geometry source read %d time(s) for 2 creates`）。甲形需要的是**同一形状的第二发**：不 dispose、只 post 一次改尺寸，再 `GetWindowRect` 比前后两样。
- **DPI 那枚 factor 的真实值**：`scale` 是从第一枚真窗**反推**的（`:169-171`），换不来无窗仪器。
- **`SetSize` 的 `AdjustWindowRect` 边框差到底几像素**：本机才量得出，而且它决定无窗仪器那条 ±容差该写多少（本腿没量，⛔ 不猜）。
- **winlive 家族的 CI 载体＝编译查、零执行**：尺（票 255 `:60` 已记）＝编排者现量 `ci.yml` 里 `tags`／`GOFLAGS` **零命中**；本腿只数文件：`grep -rln "go:build.*winlive" --include=*.go cmd internal` ⇒ **12 枚**（rc=0，清单见 commit 日志）。⇒ 这 12 枚**默认 `go test ./cmd/wisp/` 根本不编译**（`panel_geometry_255_winlive_test.go:9-11` 自述这一句），⛔ 不许读成「CI 有载体」，也不许为了让它们进 CI 去摘标签。

---

## 问⑥ 雷区与契约（只指认，不自批，不选形）

**这条链上有没有面板侧发起的写/批准动作＝没有**（三把尺）：
1. 入站白名单是固定的六枚，⛔ 没有一枚是几何：`internal/panel/bridge.go:42-45`（`panel.mode.request`／`panel.workspace.request`／`panel.attachment.add`／`panel.message.send`）＋ `:66-67`（`config.get`／`config.set`）。宿主侧的门在 `panel_host_windows.go:400-411`（`installPanelTransport`，注释 `:402-404` 逐字「No new inbound method name is invented here」）。
2. 甲形的触发源是**磁盘上的 config.toml + Go 侧装配根**，方向是 Go→Win32，**页面不在这条路上**。⇒ 一旦有人把「页面拖个滑块就改宽」做进来，那是**新增一枚入站方法名＝C17 契约面**，票 255 禁区 `:25` 已写明「凡要新增面板快照字段／新方法名＝C17 契约面 ⇒ 先停手上报」。**本腿把它指认为雷，不上报为待定案（它不在 AGENTS §2 那六条里，也不该被顺手做掉）。**
3. 审批面在这条链上是零：常驻面板的 `Confirm` 刻意留 nil（`panel_resident_windows.go:220-229` 注释逐字「Confirm stays nil on purpose … the panel side never produces an "allow"（internal/agent/approval/ui.go: PanelAPI has no Allow method）」），`UI.go` 那句话本腿复核到注释级、非实现级。

**会不会碰到需要人工批准的契约面**：

| 形 | 最少落点（枚数） | 碰不碰 C27 | 主要风险（都具名） |
|---|---|---|---|
| **甲形：对活窗直接改尺寸** | **3 枚**：`panel_host_windows.go`（新方法，紧邻 `:236`）、`panel_resident_windows.go`（一枚 `rp.post` 请求，紧邻 `:439-450`）、`config_readers_255.go:161`＋`config_receipt_255_test.go:440`＋`config_readers_255.go:17/:141` 的**同批文案/行号**（算一枚写面、三处点）。⚠ **第 4 处才是难点**：谁去按那枚请求（见问④(c)），今天**不存在** | ⛔ **不碰**：窗口不被销毁、不新建、`Show/Hide` 的隐藏语义一字不动 ⇒ C27 的「一会话至多一个面板窗口、隐藏而非销毁」与「单窗口复用绕开 Environment」(`docs/PLAN.md:1377/:2352/:3113`) 全部原样成立 | (i) `SetSize` 的入参是**客户区**语义（问①⚠），与 `windowOptions` 的**外框**语义不同 ⇒ 同一数字两处不等；(ii) hint 选 `HintFixed` 会顺手把窗改成拖不动（库 `webview.go:409`）；(iii) 跨线程调用的后果**仓里无凭据**（问③），必须走 `rp.post`；(iv) 改文案若不同批 ⇒ 直接把本票 AC#1「不许吹牛」那格踩红 |
| **乙形：保留隐藏语义、下次 Show 时重建** | **2 枚**（`Show`/`HotShow` 里加脏判定 + 建窗后回填），但**要复活一枚死方法**：`RequestDispose` 非测试调用者现量 **0 枚**（尺：`grep -rn "RequestDispose" --include=*.go cmd internal` ⇒ 2 命中，逐字是 `panel_resident_windows.go:439` 注释与 `:442` 定义本体，rc=0） | ⚠ **按字面会碰**：Show 时重建＝把一个「隐藏中」的窗**销毁**再换一枚新的 ⇒ 与 `PLAN.md:1377` 「隐藏而非销毁」的字面冲突；是否算违背属编排者裁，本腿不自批 | **代价（这是编排者点名要的那格）**：① 每次改宽之后的第一次 Show 都从热路径掉回**冷路径** ⇒ 直接烧 D32 那两行（`panel_host_windows.go:170-175` 逐字记着 `cold <=1500ms, hot <=200ms`，阈值本体在 `internal/observe/thresholds.go`，⛔ 一字节不许动）；② `NewWithOptions()` **每窗口新建一枚 WebView2 Environment 且不暴露共享入口**（`docs/PLAN.md:1568`/`:2352`，P11 的整条理由）⇒ 重建＝重新付一次 Environment＋一套 `msedgewebview2` 子进程，C27 单窗口复用存在的意义就是绕开它；③ 重建还要过 `bringUp` 的建窗前检查（`panel_host_windows.go:381-384`）：那条线程上有一枚**未派发的 WM_CLOSE** 就直接 `errPanelRefusedThread` 拒绝，且 33-r7/r9 在本机量过「purge 不派发 ⇒ 泄漏前一枚并 1/3 挂在 `GetMessageW`」「同一枚检查分不清孤儿 close 与活主人的 close」两种事故（`:338-346`、`:362-380` 逐字）⇒ 乙形是把这条**已量到危险**的路径从「零调用者」变成「每次改宽都走」 |

⇒ **「不动 C27 也能做到」的那条路＝甲形本身**（它不销毁、不新建、不改隐藏语义，改的只是一枚活窗的几何）。乙形并不比甲形更省契约，反而更贵：它把「销毁重建」这件 C27 明令绕开的事，从「今天没人调用」变成「每次改宽都发生」。⛔ 本腿到此为止，不排程、不派腿、不勾框。

---

## 我这把尺的缺陷（哪几处是读码推的、哪几格没凭据）

1. **DPI 那格是「读码＋一把空尺」推的，不是量的**。我说「Go 侧零换算」有两条 grep 支撑（问②），但**「Win32 按进程 awareness 缩放」这半句是从 `panel_geometry_255_winlive_test.go:29-34` 那段注释转述来的**，那是别人本机量过的读码复述，本腿**没有复跑**它（⛔ 无窗约束下也复跑不了）。
2. **`AdjustWindowRect` 的边框差我只证明了两条路的语义不同（码在 `webview.go:427` vs `:320`），没量出差值几像素**。⇒ 问⑤第 2 条那条容差该写多少，本腿给不出来。
3. **`Hint` 各值的行为是读 `webview.go:405-432` 那 28 行推的**（尤其 `HintFixed` 改 style、`HintMax/Min` 不进 `SetWindowPos` 分支）。**零实测**；`browser.Resize()` 到底把 controller 定成多大也只到「`panel_host_windows.go:41-49` 这么说＋库 `:431` 这么调」，本腿没打开过任何窗。
4. **「跨线程 SetWindowPos 的后果」我按派单要求写成了「仓里没有凭据」，而不是补一段常识**。⛔ 这一格是真空的，下一枚写码腿要么走 `rp.post` 绕开这个问题，要么先补一次实测。
5. **`RequestDispose` 非测试调用者＝0 枚：本腿复到了，但尺的射程只有 `cmd internal`**（派单给的显式根）。⚠ 若别的根（`tools/`、`scripts/`）里还有调用者，我这句会窄；`cmd/balldebug/` 我**只**在 `SetWindowPos`/`OnReload` 两把尺里顺带命中过，**没有专门为它起过一把覆盖 `cmd/balldebug` 的尺**（balldebug 有 `debugMoveWindow`，那是**球**窗、不是面板窗，别混）。
6. **问④(c)「谁去按那个请求」我只给了两枚候选位置，没证明第三枚不存在**。尺是 `grep -rn "startConfigReload|reloadOnce" cmd/wisp/resident*.go`＝0 命中＋`OnReload` 的射程推理（`manager.go:198` 那个 `if`）——**后者是读码，不是全仓穷举尺**（我没为「有没有别的面板侧通知机制」起过一把独立尺，例如 `notify_windows.go`／tray 那条路我没读）。
7. **行号一律是「现量」，会漂**。本文件所有 `file:line` 取自 HEAD `5738661`＋本腿工作树；票面里 `:304/:305`、腿报 `:388` 这类**别人写的行号我没有替它们更正**（只具名指出与现量不符：`config_readers_255.go:17/:141`、`config_receipt_255_test.go:440`）。
8. **一处纪律偏差自报**：起手我跑过 **1 枚 `go env GOMODCACHE`**（为拿 vendored 库的真实落点）。它**不是** `go build/vet/test`、不编译、不写仓，但严格读「零 go 命令」这一条它**越了线**；除那一次之外再无任何 `go` 调用。记在这里，不由我自行豁免。
9. **查重让位声明**：`docs/evidence/s1/180-182-panel-fields-census-c1.md:50`（「`cfg.Panel.Width` → 宿主 → `SetWindowPos/尺寸` 整条链缺的不是'一跳'」）与票 255 `:56-61`（`255-c2` 撞钉预检）已给过同一件事的读数，本程**只引用、未重跑**；本文件的增量是**库侧那枚 `SetSize`/`Dispatch` 出口的存在性＋两条路语义不等**这一格（`255-c2` 尺面只覆盖到仓内符号，未下探 modcache）。
