# 票 33 · 腿 `33-a2` · 只读普查：面板宿主接进常驻腿的阻塞面（甲／乙／丙选型的前置量尺）

⛔ 本腿零产码／零脚本／零测试／零构建。尺具只有 `Read`/`Grep`/`Glob`/`sed -n`/`find`/`go env`/`git log`/`git status`。
〔量〕＝我这发跑出来的逐字读数。〔推〕＝我从源码结构推的，**没有真跑过，不算已证**。
本版为**终态版**（骨架版＝提交 `eed229e4`，其中 ⑤⑥⑦ 已先写满；本版补完 ①–④ 与 ⑤ 的 R17–R47）。

## 起手锚点

- `date "+%Y-%m-%d %H:%M:%S %z"` ⇒ `2026-10-01 11:21:14 +0800`〔量〕
- `git log -1 --format=%H` ⇒ `22272ac2b06470b9a03b9acc811eb6d0de944423`〔量〕（起手同发取）
- `git branch --show-current` ⇒ `dev`〔量〕
- 依赖版本（`grep -n webview2 go.mod`）⇒ `go.mod:19: github.com/jchv/go-webview2 v0.0.0-20260205173254-56598839c808 // indirect`〔量〕
- module cache（`go env GOMODCACHE`）⇒ `D:\work\base\gopath\pkg\mod`；库根＝
  `D:/work/base/gopath/pkg/mod/github.com/jchv/go-webview2@v0.0.0-20260205173254-56598839c808/`〔量〕
- 本腿只写这一枚文件：`.scratch/wisp/probes/33/a2/census.md`；只加行到票面 Progress log。

---

## ① 库侧到底会发生什么（创建控件 → 泵 → 收到 web message 的线程要求）

### ①-0 一句先行的更正（影响所有后续引用）

库根下**没有**根级 `chromium.go`（我按派单路径取数得到 `sed: can't read .../chromium.go: No such file or directory`〔量〕）。
真身＝**`pkg/edge/chromium.go`**。派单说的 `:96-111` 与 `:130-131` 在真身里**号数正好成立**（见下表 g/h/k 行），
所以那句自述的**行号是可信的**；但按字面路径复核的人会撞空。⛔ 这不改变"没有原始工件"这一判定（②）。

### ①-1 创建链逐处（行号一律由 `grep -n "" <file> | sed -n 'a,bp'` 反取）

| 处 | 位置 | 逐字形状 | 线程要求／后果 |
|---|---|---|---|
| a | `webview.go:97` `NewWithOptions(options WebViewOptions) WebView` | 入口 | 无 apartment 判定〔量：全依赖里 `grep` 不到 `CoInitialize`/`COINIT_*`/`OleInitialize`/`APARTMENTSTATE`〕 |
| b | `webview.go:108` | `w.mainthread, _, _ = w32.Kernel32GetCurrentThreadID.Call()` | **调用方线程被记死**：此后一切跨线程投递都以它为"主线程"〔量〕 |
| c | `webview.go:109` | `if !w.CreateWithOptions(options.WindowOptions) { return nil }` | 同步；返回 `nil` ⇒ 宿主 `panel_host_windows.go:193-195` 才拿得到 error〔量〕 |
| d | `webview.go:284-293` | 类名 `"webview"`（字面量）、`User32RegisterClassExW` 返回值被丢 | 第二扇窗重复注册同名类的失败**不可见**〔量＋推〕 |
| e | `webview.go:320-333` | `User32CreateWindowExW(...)` | 窗口归属**调用方线程**〔推：Win32 通识＋b〕 |
| f | `webview.go:340` | `if !w.browser.Embed(w.hwnd) { return false }` | ⬇ 落到 edge |
| g | `pkg/edge/chromium.go:87-94` | `createCoreWebView2EnvironmentWithOptions(nil, utf16(dataPath), 0, e.envCompleted)`；`err != nil` ⇒ `log.Printf`+`return false`；`res != 0` ⇒ `log.Printf("Result: %08x")`+`return false` | 异步投递＋COM 回调；**两条失败只打日志、不 panic**〔量〕 |
| h | `pkg/edge/chromium.go:95-111` | `var msg w32.Msg; for { if atomic.LoadUintptr(&e.inited) != 0 { break }; r,_,_ := w32.User32GetMessageW.Call(&msg, 0,0,0); if r == 0 { break }; TranslateMessage; DispatchMessageW }` | **一枚阻塞式嵌套泵**；`hwnd=0` ⇒ 取该线程队列上**所有窗口＋线程消息**〔量：形状；语义见 ①-3〕 |
| i | `pkg/edge/chromium.go:224`（在 `:186 CreateCoreWebView2ControllerCompleted` 内） | `atomic.StoreUintptr(&e.inited, 1)` | 置位**只**发生在 COM 回调里 ⇒ h 必须自己把回调泵下来〔量＋推〕 |
| j | `pkg/edge/chromium.go:112` | `e.Init("window.external={invoke:s=>window.chrome.webview.postMessage(s)}")` | h 出口之后**无条件**执行〔量〕 |
| k | `pkg/edge/chromium.go:130-136` | `func (e *Chromium) Init(script string) { _, _, _ = e.webview.vtbl.AddScriptToExecuteOnDocumentCreated.Call(...) }` | ⚠ **`:130-136` 全段无任何 `if`**；`e.webview` 只在 `:194-197` 被赋值、且那次 `Call` 的返回三元组被丢弃 ⇒ **nil 解引用＝Go panic**〔量：无判定；推：nil ⇒ panic〕 |
| l | `webview.go:113-126` | `settings, err := chromium.GetSettings(); if err != nil { log.Fatal(err) }`；两枚 `Put*Enabled` 各带 `log.Fatal` | **失败即杀进程**（见 ①-4）〔量〕 |
| m | `webview.go:343` | `w.browser.Resize()` | ⇒ `pkg/edge/chromium_amd64.go:12 Resize()`、`:18 e.controller.vtbl.PutBounds.Call(...)`〔量〕 |

### ①-2 收到 web message 时跑在哪枚线程

- `pkg/edge/chromium.go:201-205` 注册 `AddWebMessageReceived(e.webMessageReceived)`；回调落 `:233-248 MessageReceived` ⇒ `:240 e.MessageCallback(...)`，库内接成 `webview.go:103 chromium.MessageCallback = w.msgcb`〔量〕。
- `webview.go:139-160 msgcb` ⇒ `:147 w.callbinding(d)`（**同步反射调用** `:191 res := v.Call(args)`）⇒ 随后三枚分支全部把回执交给 **`w.Dispatch(func(){ w.Eval(...) })`**（`:148`/`:152`/`:156`）〔量〕。
- `webview.go:443-448 Dispatch`：`:445 w.dispatchq = append(w.dispatchq, f)`；**`:447 w32.User32PostThreadMessageW.Call(w.mainthread, w32.WMApp, 0, 0)`**〔量〕。
- `w32.WMApp = 0x8000`（`internal/w32/w32.go:96`）、`WMQuit = 0x0012`（`:92`）〔量〕。
- **`dispatchq` 在全依赖里只被读一次**：`grep -rn dispatchq <dep>` ⇒ 命中只有 `webview.go:58`（字段）、`:445`（写）、**`:362-363`（读，唯一一处，位于 `webview.go:351 Run()` 内部 `case msg.Message == w32.WMApp`）**〔量〕。

⇒ **本小节的产出（这是我最想交给裁的一格）**：
**库的"页 → Go → 页"闭环里，"Go → 页"那一跳只有 `webview.Run()` 能兑现**〔量：`:360-367` 是唯一 drain 点；推：线程消息不经 `DispatchMessageW` 送达窗口过程〕。
而本仓宿主**从不跑 `Run()`**：`grep -n "\.Run()\|Terminate()\|Dispatch(" cmd/wisp/panel_host_windows.go cmd/wisp/panel_host_windows_test.go` ⇒ **零命中**〔量〕；
全仓 `grep -rn "\.Dispatch(" --include=*.go cmd internal`（去测试）⇒ 只命中 statemachine 的 `Dispatch`，**零枚 webview2 的**〔量〕。
宿主自己的泵 `pnlPumpOnce`（`panel_host_windows.go:396-411`）是 `PeekMessageW(…,0,0,0,PM_REMOVE)`＋`DispatchMessageW`〔量〕——
`hwnd=0` 的 `PeekMessageW` **会把这些线程消息取出来**，但取出来后没有任何人读 `dispatchq` ⇒ 〔推〕**回执闭包被静默丢弃、且队列无界增长**。

**这条推断能解释盘上已有的两枚事实**（不是新增证据，只是自洽）：
- `panel_host_windows_test.go:190-197` 的 H3/H10 是**用 Go 直接调 `mgr.dispatchRaw(...)`** 证的，页面从未发过那一枪〔量〕；
- `firstRoundTripLocked` 的 `done` 在 `callbinding` 内就被关闭（`panel_host_windows.go:361-368`），**早于** `msgcb` 里的 `Dispatch`（`webview.go:147`→`:148`）〔量＋推〕 ⇒ 那枚 `coldMs` 证的是"页 → Go"通、**不证**"Go → 页"通〔推〕。
（`v1` 表 §A#28 独立记的"冷启结束窗里是探针文档"是同一条链的另一枚症状。）

### ①-3 一枚**已经在泵**的 STA 上再入创建：三种最坏后果形状

⚠ **全节〔推〕**（本腿不许跑）。只有真跑能定案的那一半留在 ⑦-A。

1. **只是阻塞（最安静）**：h 自洽——同线程的嵌套泵能把 i 的 COM 回调送下来。后果＝调用方在**外层 `DispatchMessageW` 里再入一层泵**，外层在此之前不再前进〔推〕。
   触发条件（〔量〕）：`internal/ball/sta_windows.go:62 pCoInitializeEx.Call(0, coinitApartmentThreaded)` ⇒ `ui-sta` 确为 STA；`sta_windows.go:78-92` 确为主泵；`ball_windows.go:665/672/677 b.fire(b.opts.Events.OnPanelHotkey/OnTrayPanel)` ⇒ 那两枚回调**确在泵内**（`fire` 逐字 `:728-732`＝`fn()` 直调，**无投递、无 recover**）〔量〕。
2. **再入窃队列／顺序错乱**：h 的 `GetMessageW(&msg, 0, 0, 0)` 覆盖**整条线程队列**。球的 UI 工作全走同一条队列：`sta_windows.go:127-143 PostTask`（`:142 pPostMessageW.Call(hwnd, wmAppTask, id, 0)`）→ `sta_windows.go:91-92 Translate/Dispatch` → `ball_windows.go:688-689 case wmAppTask: b.sta.runTask(uint64(wParam))`〔量〕。
   `wmAppTask = wmApp + 0x201 = 0x8201`、`wmApp = 0x8000`（`internal/ball/win32_windows.go:99,102`）〔量〕⇒ 与库的 `WMApp(0x8000)` **编号不撞**，但**同队列**：嵌套泵会把外层本该按序处理的消息插到 `bringUp` 返回之前执行〔推〕。
   球自己的两处注释逐字承认这套模型：`sta_windows.go:54-55` "a caller already on it must run inline (see Ball.uiRun) - post-and-wait from inside the pump is a deadlock"；`ball_windows.go:724-727` "CONTRACT: callbacks must be quick and non-blocking (statemachine dispatch + PostTask); a blocking callback stalls the message pump. No goroutine is spawned per gesture (D38b roster discipline)"〔量〕
   ⇒ **甲形若把创建挂在 `OnPanelHotkey` 里，就正面违反这条 `fire` 契约**（创建是最阻塞的那种回调）。〔推〕
3. **panic（那句自述的因）**：h 有一枚**不看 `inited` 就 `break`** 的出口（`:106-108 if r == 0 { break }`＝该线程队列收到 `WM_QUIT`），随后 `:112 → :130-131` 解引用可能仍为 nil 的 `e.webview`〔量：无判定；推：nil ⇒ panic〕。
   - 谁能投 `WM_QUIT`：`ball_windows.go:962 b.sta.quit()` → `sta_windows.go:156-161 pPostQuitMessage.Call(0)`；库侧 `webview.go:381-383 Terminate()`，且 `webview.go:242-243` 在窗口 `WMDestroy` 时**自动**调它〔量〕。
   - ⇒ **最坏后果的形状因此分岔**，而这一岔我能量到、不用跑：
     **出货路径的 panic 不会杀进程**。`ui-sta` 是经 `observe.Registry.Spawn` 起的（`ball_windows.go:170 b.sta = newSTAThread(...); :170 b.sta.handle = opts.Registry.Spawn("ui-sta", "ball", nil, func(_ context.Context) { b.sta.start(...) })`〔量〕），而 `internal/observe/goroutine.go:294-326 run()` 的 `defer` 里逐字有 `if rec := recover(); rec != nil { ... Stack: string(debug.Stack()) ... r.panicCount++ ... sink(ev) ... root.Cancel() }`〔量〕。
     ⇒ 〔推〕**真跑出来的形状是：panic 被 observe 兜住 → 落一枚带 `debug.Stack()` 的 `PanicEvent` → `ui-sta` 那根协程 return → 球窗与面板同归于尽，进程还活着（`Close()` 的 `<-b.sta.handle.Done()`（`ball_windows.go:965-966`）反而会立刻返回）**。
     这条比"进程炸了"更糟也更难查，而且**它意味着复现探针根本不用自己抓栈**：只要走 `Registry.Spawn` 那条路，`observe` 的 `SetPanicSink`（已有仪器 `internal/observe/goroutine_test.go:116-117`）就会把原始栈交出来。〔推〕

### ①-4 库侧把"失败"升级成"进程没了"的枚数（我逐字数过）

`webview.go:115`、`webview.go:120`、`webview.go:125`（三枚 `log.Fatal`）＋`webview.go:141`（`Eval` 里 `UTF16PtrFromString` 失败 ⇒ `log.Fatal`）；
`pkg/edge/chromium.go:173`（环境创建失败 `log.Fatalf`）、`:188`（控制器创建失败 `log.Fatalf`）、`:284`（`WebResourceRequested` 的 `GetRequest` 错误 `log.Fatal`）、`:295`（`AddWebResourceRequestedFilter` 错误 `log.Fatal`）、`:141`（Eval 那枚在 edge 侧）。〔量：以上行号均由 `grep -n` 输出逐条对上〕
⇒ `v1` 表 §J2 已把这条判成"J2 那句'不走会 fatal 的 API'不成立"〔量：`33-panel-host-c27-v1.md:179`〕。**本腿复认其行号，不重判。**

### ①-5 本格欠账（⛔ 不许读成已证）

只有真跑能定的两半：**(i)** 在已在泵的同一枚 STA 上再入创建，实际走 ①-3 的第 1／2／3 哪一支；**(ii)** 走第 3 支时是不是 `edge.(*Chromium).Init` 栈顶。
本腿⛔不许跑（`33-r3` 正在 `cmd/wisp` 里跑测试）。⇒ 见 ⑦-A（那里我把"怎么一发起就同时拿到栈与判定"写成可执行形状）。

---

## ② 那句〔自述〕的工件面

**结论：`.scratch/wisp/probes/33/` 下＝零该发工件。**〔量〕

搜过的路径（逐枚具名）：

1. `ls -R .scratch/wisp/probes/` ⇒ 顶层含 `33`（另有 `132/139/…/248/v5-survey` 等）。〔量〕
2. `ls -R .scratch/wisp/probes/33/` ⇒ **`a1` `h1` `r1` `r2` `r3`** 五枚子目录（本轮新增 `a2`）。〔量〕
3. `find .scratch/wisp/probes/33/ -type f | sort` ⇒ **84 枚**全名单（⑤-R3）。其中：
   - `r1/`＝`composer_dispatch.go.pristine`、`mutate.sh`、`logs/{d22scan,gate-clauses,green-after,green,mutate-run,panel,red-M1..M4,risk}.txt`
   - `a1/logs/`＝`{d22scan,gate-clauses,gotest}.txt`
   - `h1/`＝`census.md`、`msg-1to7.txt`、`msg-8910.txt`、`msg-skeleton.txt`
   - `r2/`＝33 枚（含 `mut-1/`、`mut-2/`；全是 gate/gotest/gofumpt/mutation 的 stdout 与 `.go` 快照）
   - `r3/`＝16 枚（同上形状）
   ⇒ **没有一枚**名字或内容形如 panic dump／复现命令（`r1/mutate.sh` 是变异脚本，不是复现脚本）。〔量＋推：判据见 4〕
4. `grep -rni "panic" .scratch/wisp/probes/33/` ⇒ **15 枚命中全部来自 `gate-clauses.txt` 那五份**（`a1/logs/`、`r1/logs/`、`r2/gate.txt`、`r2/gate-final.txt`、`r3/gate-final.txt`），逐字只有两类：
   `internal/memory/retention.go:100: panic("memory: StartRetentionJob requires a DisposalScope")` 与
   `internal/plugin/disposal_test.go:87: s.DeferNamed("panics", func() { panic("destroy failed: injected") })`
   ⇒ 闸门把仓源码逐行贴进日志留下的**字面回声**，与 webview2／`bringUp`／`ui-sta` 无关；**零命中**含 `webview2`、`GetMessageW`、`goroutine [0-9]+ \[running\]` 之类栈头形状。〔量〕
5. `grep -rn "panic\|⑧" docs/evidence/s1/33-panel-host-c27-r1.md` ⇒ 命中 `:7,28,33,41,45,46,53,66,73`。〔量〕
   - `:73` 逐字："本机现量 `edge.Chromium.Init` nil-webview panic（栈：`ballWndProc→staThread.start DispatchMessageW→…→PanelManager.bringUp→webview2.NewWithOptions→CreateWithOptions(webview.go:340)→Embed→Init(chromium.go:131)`）" ⇒ **散文栈**：中间 `…` 省略、无 goroutine 头、无 `exit status 2`。〔量〕
   - `:46` 逐字自陈："我一度写了这枚再入用例（本机现量 nil-webview panic 栈已入 ②），因它会把 `cmd/wisp` 留红……**未提交**该文件"。〔量〕
   - `:45` 逐字："从它再入跑 go-webview2 的阻塞嵌套泵会 panic（见 ② 硬伤 2）⇒ 这条出货接线本腿**未做**"。〔量〕
6. `v1` 表自己扫过盘并同判：`33-panel-host-c27-v1.md:60` 逐字"我扫盘找凭证……全仓 `grep -rl "panic: runtime error"` 只命中它自己那份证据件（**没有原始输出、没有复现命令**）⇒ 判〔仅自述＋栈形状〕"。〔量〕

**本格独立增量（我没白读的两件）**：
- 那句散文里的两处行号**我独立对上库源码**：`webview.go:340 = if !w.browser.Embed(w.hwnd) {`〔量〕、`pkg/edge/chromium.go:131 = e.webview.vtbl.AddScriptToExecuteOnDocumentCreated.Call(`〔量〕。⇒ 栈形状与真实代码自洽，可信度上升；**仍不等于那发 panic 发生过**。〔推〕
- 盘上找不到栈这件事，**和我量到的 recover 形状不矛盾、反而互相解释**：那条用例若走的是**测试 harness**（`panel_host_windows_test.go:39-62`），harness 自己的 `recover` 只产 `fmt.Errorf("panel host pump panicked: %v", r)`（`:45`）——**这个形状天生不含栈**〔量〕。⇒ "没有原始栈"可能是**取数方式**决定的，不是"没跑过"。这条改变 ⑦-A 的最省复现形状：换走 `Registry.Spawn` 那条就有 `debug.Stack()`。〔推〕

---

## ③ 甲形的真实射程（导出什么才算"最小一枚"＋会被叫红的现成用例）

### ③-1 现有导出面：够不够？〔量〕

**全目录尺**（补 R12 的单文件漏计，尺＝`grep -rn "^func (b \*Ball) [A-Z]\|^func (s \*staThread) [A-Z]" internal/ball/*.go`）：

- 导出方法全名册（**跨文件**）：`SetState:311`、`SetBadge:316`、`SetProgress:324`、`SetBadgeText:332`、`TimersAlive:441`、`DebugTimersAlive:450`、`RebindHotkeys:813`、`HotkeyReport:830`、`ConfiguredHotkeys:840`、`RegisteredHotkeys:849`、`TakeEscForCancel:873`、`ReleaseEscAfterSession:895`、`EscTakenOver:910`、`SetTrayChecks:919`、`SetTrayTip:927`、`Close:935`、`DebugHWND:1003`（`ball_windows.go`）＋ `DebugDock:255`、`DebugPop:291`、`Docked:313`（`dock_windows.go`）＋ `SetAudioLevel:42`（`liquid_windows.go`）。
- 唯一命中大写方法名的 `staThread` 行＝`sta_windows.go:127 func (s *staThread) PostTask(fn func())`〔量〕——**receiver 类型本身未导出**（`sta_windows.go:24 type staThread struct`）⇒ 包外不可用。〔量〕
- 未导出的还有：`Ball.sta`（`ball_windows.go:84 sta *staThread`）、`Ball.fire:728`、`Ball.uiRun:741`、`newSTAThread:38`、`start:47`、`waitStarted:111`、`threadID:120`、`runTask:146`、`quit:156`、`releaseCOM:103`、`threadID:120`。
- 导出类型只有 `EventKind:33`、`Events:49`、`Options:63`、`Ball:82` ＋ `New:143`。
⇒ **今天包外没有任何一条路能把闭包投到 `ui-sta`**〔量，全目录尺，不再是单文件推〕。
⇒ **`PostTask` 确实存在**（`sta_windows.go:127`，25+ 处内部调用点：`ball_windows.go:312,317,325,333,747,920,928,940`、`dock_windows.go:260,294,320`、`liquid_windows.go:46`）——验收腿那句补充成立。〔量〕

### ③-2 "最小一枚"的候选形状（只列形状，⛔ 不选形）

| 候选 | 形状 | 射程／代价 |
|---|---|---|
| **甲-1** | `func (b *Ball) PostUITask(fn func())`＝`b.sta.PostTask` 一层包装 | **一枚方法、零新类型、零新契约**。但⚠ 拿不到错误、拿不到完成通知，而 `sta_windows.go:134-141` 在 `hwnd == 0` 时 `delete`＋`return`＝**静默丢任务**〔量〕⇒ 面板创建失败无人知道。 |
| **甲-2** | 一枚导出接口／类型（`STAHandler`）＋`b.STA()` | 一枚方法＋一枚类型＝把 `staThread` 的形状升成跨包契约面（触 D38 叙述的可能性最高）。 |
| **甲-3** | `func (b *Ball) RunOnUIThread(fn func())`（等待完成，＝导出 `uiRun:741`） | ⛔ `sta_windows.go:54-55` 与 `ball_windows.go:724-727` 两处注释明写这套模型的边界（post-and-wait 在泵内＝死锁；回调必须非阻塞）。甲-3 与"在回调里创建"互斥。〔量：注释原文；推：互斥〕 |
| **甲-4** | 在 `Events`（`ball_windows.go:49-61`）加一枚宿主 hook | ⚠⚠ **必叫红一枚在册名册钉**，见 ③-3 第一族。 |

⇒ "最小一枚"在**形状度量**上是甲-1（一枚方法、零新类型）〔量＋推〕。**⛔ 我不选形。**
**⚠ 但甲-1 的射程不足以兑现 AC#1 的"回执到达页面"**：①-2 已量到库侧 Go→页那一跳只有 `Run()` 会 drain `dispatchq`〔量〕⇒ 任何"只导出投递面"的形都还得解掉那一跳，这件事我在 ⑦-H 单列给裁。

### ③-3 ⚠ 会被叫红的现成用例，逐枚点名（本发的头号产出）

#### 第一族｜`cmd/wisp` 的**词面型名册钉**（只被 甲-4 打红）

**F1 `TestAC228BallHostAnswersEveryGesture`（`cmd/wisp/resident_ball_228_test.go:275-324`）**
- 名册原文（`:50-53`）：
  `var ballEventCallbacks228 = []string{ "OnClickBall", "OnSummonHotkey", "OnMuteHotkey", "OnCancelHotkey", "OnPanelHotkey", "OnTrayPanel", "OnTrayMute", "OnTrayPauseWake", "OnTrayExit", "OnDragEnd" }`〔量〕
- 设计意图原文（`:45-49`）：**"It is a literal list on purpose: when internal/ball grows an eleventh gesture, the case below goes red and whoever added it has to say what the resident leg does with it."**〔量〕
- 两枚断言原文：
  - `:313-317` `if len(missing) > 0 { t.Fatalf("AC#1 RED: the resident ball host leaves %d of internal/ball's %d gesture callbacks unset: %v. ...") }`
  - `:319-321` `if len(set) > len(ballEventCallbacks228) { t.Fatalf("AC#1 RED (the instrument, not the code): the host sets %d callbacks but this file lists %d; internal/ball.Events grew and ballEventCallbacks228 did not.") }`〔量〕
- **为什么会被叫红**：它用 `ast.Inspect` 找 `Events` 复合字面量（`:293 sel.Sel.Name != "Events"`）把**被设置的键名**收进 `set`，再和字面名册比。⇒ **往 `ball.Events` 加任何一枚新回调＝`set` 变 11、名册仍 10＝`:319` 当场红**；红句自己就说"这是仪器、不是码"。〔量＋推〕
- **定性＝词面型名册钉**（钉的是名字集合与枚数）。⇒ **这一枚必须在派单里当场裁**：甲-4 要落就得同时改这枚名册（改名册＝动别人的在册钉，需编排者授权），**不该等写完被打红**。甲-1/2/3 一枚回调都不加 ⇒ **不被叫红**〔推，依据＝`Events` 字段名册 `:49-61` 不变〕。

**F2 `TestAC228ResidentLegIsTheBallHost`（同文件 `:196-269`）**
- `:205-208` 要求 `cmd/wisp` 至少一枚生产文件 import `internal/ball`；`:211-213` 要求有人真调 `ball.New`（现量＝`resident_ball_windows.go:114` 唯一一处〔量〕）；`:246-249` 要求 `runResident` 的函数体 walk 到那枚 host；`:255-267` 要求 `runResident` 里 `defer ...stop()`。
- **定性＝行为型＋形状混合钉**（钉的是"常驻腿真的把球宿主接上了"这一条不变式）。
- **判定**：甲形**只做投递面**⇒ 四条都还成立〔推〕。⚠ 但如果接面板时**顺手把 `ball.New` 或那句 `defer ... .stop()` 挪出 `runResident`**，`:246`／`:265` 会当场红。派单要写死"这两条不许动"。〔推〕

#### 第二族｜`internal/panel` 的**能力型正向钉**（甲形若搬家会把它掏空）

**F3 `TestPanelHostIsAttachedAndNamesTheWindowHops`（`internal/panel/composer_dispatch_test.go:433-447`）**
- `:438-442` 原文判据：`if !prodSignal { t.Errorf("no native host / WebView2 message channel is attached in this tree: no production source carries a CoreWebView2/WebView2/WebMessage identifier or a webview import. Ticket 33 landed H2/H3/H10 - a tree that links no host while the panel window ACs read done is exactly the false-green this nail now forbids") }`〔量〕
- `:443-447` 另一枚 `depSignal`（`go.mod`/`go.sum` 必须载 webview 路径）〔量〕。
- 它自带三枚正控（`:449-493` carrier A/B/C），其中 `:482-487` 逐字说明 C 是"分辨腿"：**只把 `ComposerDispatch` 递给 `Handle`、没有窗口的形状必须保持沉默**〔量〕。
- **定性＝能力型正向钉**（比词面型高一档：扫标识符/import，不扫注释）。
- **判定**：甲形只要**保留一枚生产源载着 webview import／CoreWebView2 标识符**就不红〔推〕。⚠ 风险形状＝把创建从 `panel_host_windows.go` 整枚搬走、或搬进 `_test.go`（`:51-53` 的正控明说非测试文件才算数）。派单该钉一句"载体不许变成只有测试"。〔推〕

#### 第三族｜**负向钉的"射程漏孔"**（甲形不打红它，但会在它脚下开个洞）

**F4 `TestPanelHostOpensNoListeningSocketL1`（`cmd/wisp/panel_host_gate_test.go:26-38`）**
- `:27` 逐字：`src := filepath.Join("panel_host_windows.go")`；`:38` 逐字：`t.Errorf("panel host imports %q, which can open a listening socket; D29/AC#3 forbid a localhost server", path)`〔量〕
- **这是全仓唯一按文件名钉住宿主的那枚尺**（尺＝`grep -rn "panel_host_windows.go" --include=*.go .` 去证据面 ⇒ 只这一处命中）〔量〕。
- **定性＝词面型（按文件名寻址）负向钉**。
- **判定**：甲形**不会叫红**它，但**只要把窗口创建挪进新文件（例如 `panel_host_ui_windows.go`），L1 的扫描面就自动少一枚文件＝假绿**〔推，依据＝路径写死〕。`v1` 表 §A#32 已经独立抓到过这枚钉"只看一个文件"的口径问题。**这一条我建议进派单判据（不是选形，是漏孔）。**

#### 第四族｜行为型钉：球自己的取数用例（甲形的真实代价面，需真跑才知红不红）

**F5** 包内测试**逐枚 `package ball`**（尺＝对 `internal/ball/*_test.go` 逐枚 `grep -m1 "^package "` ⇒ 11 枚全是 `package ball`）〔量〕
⇒ **"导出一枚方法会让它们编译不过"这条** **不成立**〔量＋推：包内可见性〕。**骨架版 ⑥-6 挂的欠账就地结清。**

**F6** 但**顺序/存活**类断言是行为型钉，甲形若引入 ①-3 第 2 支（再入乱序）就可能打红：
- `internal/ball/live_windows_test.go:224,229,438,529,534,555,589,658,673`（逐字形状 `b.sta.PostTask(func() { done <- … })`，如 `:529 done <- b.dock.p`、`:658 liquid <- b.liquidTimerActive`）〔量〕
- `internal/ball/live_guard_windows_test.go:102,114,253,322`（`:322 b.sta.PostTask(func() { done <- b.curState })`）；`:72 type timerCount struct{ wmTimer, total atomic.Uint64 }`、`:79-80 if msg == wmTimer { p.wmTimer.Add(1) }`〔量〕
- `internal/ball/hotkey_live_test.go:60,327-334`（`:327 h := b.sta; h.PostTask(...)`），以及 `:393,411,416,418` 的 `cnt.wmTimer.Load()` 计数断言〔量〕
- **定性＝行为型钉**（钉"投上去必须按序跑到／WM_TIMER 必须被数到"）。**甲-1/2/3 都不新增协程、不改名册 ⇒ 名册类不红；会不会按序乱＝只有真跑知道**（⑦-A）。〔推〕

**F7** `cmd/wisp/hostThreadHarness`（`panel_host_windows_test.go:30-62`）＋ `:166-167` 的 `t.Fatalf("bringUp on the test thread: %v ...")`〔量〕
- 注释原文（`:30-34`）逐字："runs the host on a private locked OS thread for the local measurement / round-trip tests. It is deliberately a TEST harness - **ban #1 targets production goroutines, and production bringUp runs on the caller's ui-sta thread, never here**"〔量〕
- **定性＝行为型钉＋叙述契约**：harness 的"我起的线程＝宿主线程"这一前提在甲形下不再真；`:43-47` 的 `recover` 会把 panic 变成 error，`:48 runtime.LockOSThread()` 会与"面板必须建在 `ui-sta`"矛盾。⇒ **甲形若不改 harness，`:166` 那枚 Fatalf 会在真跑时先红**（因为它把"宿主线程是测试线程"当成立前提）。〔推〕

#### 第五族｜⛔ 只有**乙形**会打红的冻结名册（本腿只登记代价，不建议）

**F8** `internal/observe/goroutine.go:40 const ResidentBaseline = 6`、`:44 "ui-sta", "audio-capture", "hotkey-listener", "db-writer", "watchdog", "log-flusher"`、`:422 rep.ResidentOverBaseline = rep.Resident > ResidentBaseline`〔量〕
**F9** `internal/observe/goroutine_test.go:268-269 if ResidentBaseline != len(ResidentNames) { t.Fatalf(...) }`、`:243 "ui-sta": CategoryResident`〔量〕
- **定性＝词面型名册钉**。乙形多一枚常驻 ⇒ `ResidentOverBaseline` 为真＋名册枚数不等 ⇒ 红〔推〕。**这一族属人工批准那一档**，按任务书我只登记代价。〔量＋推〕
- 甲形**一枚协程都不起**（`sta_windows.go` 不动 `Spawn` 次数）⇒ 这族不被叫红〔推〕。

**本格结论**：甲形的**代码面**代价确实可能只是"一枚导出方法"（甲-1），**但它的射程不覆盖 ①-2 那枚"Go→页只有 `Run()` 能兑现"的库事实**；
它的**测试面**代价按候选分岔得很清楚：**只有甲-4 会撞在册名册钉 F1（且那枚红的注释自己写明"是仪器不是码"）**；F3/F4 是"搬家会把钉掏空"的漏孔型代价；F6/F7 是需要真跑才知道的行为型代价。

---

## ④ 有没有第四形（只量不发明）

| 现成形状 | 位置〔量〕 | 能不能承住 WebView2 的泵要求 | 依据 |
|---|---|---|---|
| 球的 `ui-sta` 泵 | `sta_windows.go:47-96`（`:52 LockOSThread`、`:62 CoInitializeEx(STA)`、`:78-92 GetMessageW/Translate/Dispatch`） | **能泵；STA 也满足**；⛔ 但它是库侧嵌套泵的**唯一候选共存者**（①-3 第 2 支未证） | 它就是 D38a 那枚"ONE thread owns all UI COM objects (… and later the panel)"（`sta_windows.go:5-9` 逐字）〔量〕 |
| `wmAppTask` 投递链 | `sta_windows.go:127-143`→`ball_windows.go:688-689` | 甲形底座（原语存在、未导出） | ③-1 |
| **`webview.Run()`** | `webview.go:351-379` | **能、而且是唯一能兑现回执的形状**（`:362-363` 是 `dispatchq` 的唯一读者） | ⛔ 它是**该线程唯一主泵**的形态（`:368 WMQuit` 才 return）⇒ 与 `ui-sta` 泵同线程互斥〔推〕，另起线程＝乙形名册钉（F8/F9） |
| **`webview.Dispatch`** | `webview.go:443-448`（`:447 PostThreadMessageW(mainthread, WMApp)`） | ⛔ 它**只投不取**：取只发生在 `Run()` | ①-2；且本仓零枚调用者（`grep` 现量〔量〕） |
| `notify_windows.go` 那套 | `cmd/wisp/notify_windows.go:137-186`（`:138 LockOSThread`、`:149 CreateWindowExW(...hwndMessage...)`、`:155 defer DestroyWindow`、`:183 time.Sleep`） | ⛔ **不能**：它**根本不泵** | `grep -n "GetMessageW\|PeekMessageW\|DispatchMessageW" cmd/wisp/notify_windows.go` ⇒ **零命中**〔量〕。它建完 message-only 窗就用 `time.Sleep` 守气球通知 ⇒ 没有消息循环，WebView2 的 COM 回调无处送达 |
| 球的 tray | `internal/ball/tray_windows.go:30 addTrayIcon(hwnd ...)`、`:55 setTip`、`:62 remove`、`:72 showMenu` | ⛔ 不是第二枚泵 | 它复用**球窗口自己的 `hwnd`**（`ball_windows.go:928 b.sta.PostTask(func(){ b.tray.setTip(tip) })` 投回 `ui-sta`）〔量〕 |
| `staThread` 第二实例 | `sta_windows.go:24-36` 结构只有**一枚** `hwnd`；`newSTAThread:38-44` 只造一枚；`ball_windows.go:166` 全仓唯一调用点 | ⛔ 现成形状不支持"同一名下第二枚泵" | 〔量＋推：真要第二枚＝乙形〕 |
| spike 测量腿 | `scripts/spike/webview2-latency/main.go:164 bringUp` | 一次性进程，不是常驻形状 | `panel_host_windows.go:29-31` 逐字："a test harness, not the shipping resident topology, and it registers no goroutine name" |

⇒ **第四形的真实候选只有一枚**：**让面板宿主拥有它自己的泵（`webview.Run()` 那一形）**，而它和 `ui-sta` **不能同线程共存**〔推，依据＝两处 `GetMessageW` 都取整条线程队列＋`Run` 只被 `WM_QUIT` 结束〕。
⇒ **所以"甲／乙／丙"三形之间，盘上还有一个未列的形状**：**甲＋库侧改造**（在 `ui-sta` 上让宿主自己 drain `dispatchq`——但 `dispatchq` 未导出 ⇒ 只能靠 replace/fork 或换绑定＝**丙的形状之一**）。〔量＋推〕这一条请直接读 ⑦-H。

---

## ⑤ 我跑了哪些尺、每条真实读数（终态：R1–R47）

| # | 命令（完整） | 逐字读数（摘要） |
|---|---|---|
| R1 | `date "+%Y-%m-%d %H:%M:%S %z"; git log -1 --format=%H; git status --short \| head -40; git branch --show-current` | `2026-10-01 11:21:14 +0800`／`22272ac2b06470b9a03b9acc811eb6d0de944423`／`dev`；脏项＝` M .gitignore`、` M .scratch/wisp/probes/{152,161,228}/…`、`design/**` 一大片 `D`/`M`、`?? .scratch/.scratch/`、`?? .scratch/commit-msg-*.txt` 多枚。⛔ 一律不动、不提交、不还原 |
| R2 | `ls -R .scratch/wisp/probes/` | 顶层名册含 `33`（另 `132 139 145 … 248 done-* v5-survey` 等） |
| R3 | `ls -R .scratch/wisp/probes/33/` | `a1 h1 r1 r2 r3`（＋本轮新建 `a2`）；子文件名单见 ②-2/②-3 |
| R4 | `grep -n "webview2\|jchv" go.mod` | `19: github.com/jchv/go-webview2 v0.0.0-20260205173254-56598839c808 // indirect`；`20: github.com/jchv/go-winloader v0.0.0-20250406163304-c1995be93bd1 // indirect` |
| R5 | `go env GOMODCACHE`；`ls ~/go/pkg/mod/github.com/jchv/` | `D:\work\base\gopath\pkg\mod`；名册＝`go-webview2@v0.0.0-20260205173254-56598839c808`、`go-winloader@v0.0.0-20250406163304-c1995be93bd1` |
| R6 | `find <dep> -name "*.go" \| sed 's\|.*56598839c808/\|\|\' \| sort` | 40 枚：`cmd/demo/main.go`、`common.go`、`internal/w32/{w32,w32_386,w32_64bit}.go`、`pkg/edge/*`（`chromium.go`、`chromium_386/amd64/arm64.go`、`corewebview2.go`、`comproc*.go`、`ICoreWebView2*.go`、`guid.go`）、`webview.go`、`webviewloader/module*.go`。⛔ **无根级 `chromium.go`** |
| R7 | `sed -n '60,150p' <dep>/chromium.go` | `sed: can't read .../chromium.go: No such file or directory`（⇒ 派单路径不成立，真身 `pkg/edge/`） |
| R8 | `grep -n "" <dep>/pkg/edge/chromium.go \| head -200` 与 `\| sed -n '200,320p'` | `:17-45` 结构、`:47-70 NewChromium`、`:72-114 Embed`（`:87`、`:95-111`、`:106-108`、`:112`）、`:116-136`（`:130-136` 无 `if`）、`:138-149 Eval`、`:151-157 Show/Hide`、`:171-184`（`:173 log.Fatalf`）、`:186-231`（`:188 log.Fatalf`、`:194-197 GetCoreWebView2`、`:201-220` 四枚 `Add*`、`:224 inited=1`）、`:233-248 MessageReceived`、`:258-279 PermissionRequested`、`:281-297`（`:284`/`:295 log.Fatal`）、`:306-320 AcceleratorKeyPressed` |
| R9 | `grep -n "" <dep>/webview.go \| head -120`／`\| sed -n '120,420p'` | `:21-36 windowContext`、`:38-47 browser 接口`、`:49-59 webview{hwnd,mainthread,...,dispatchq}`、`:87/:92-94/:97`、`:102-111`、`:113-126`（三枚 `log.Fatal`）、`:131-137 rpcMessage`、`:139-160 msgcb`、`:162-220 callbinding`（`:191 v.Call`）、`:222-261 wndproc`（`:240 WMClose`、`:242-243 WMDestroy→Terminate`）、`:263-267`、`:269-345 CreateWithOptions`（`:284`、`:293`、`:320-333`、`:334`、`:336-338`、`:340 Embed`、`:343 Resize`）、`:347-349 Destroy`、`:351-379 Run`（`:354`、`:360-367`、`:368`、`:371-372`）、`:381-383 Terminate`、`:385-387 Window`、`:389-395 Navigate/SetHtml`、`:405-420 SetSize` |
| R10 | `grep -n "func (w \*webview) Dispatch" -A 14 <dep>/webview.go` | `:443-448`，`:447 User32PostThreadMessageW.Call(w.mainthread, w32.WMApp, 0, 0)` |
| R11 | `grep -rn "dispatchq" <dep> --include=*.go` | **仅 4 枚**：`:58`、`:362`、`:363`、`:445` ⇒ `Run()` 是唯一读者 |
| R12 | `grep -rn "WMApp\|WMQuit\|WM_APP\|WM_QUIT" <dep> --include=*.go` | `internal/w32/w32.go:92 WMQuit = 0x0012`、`:96 WMApp = 0x8000`、`webview.go:368` |
| R13 | `grep -n "" <dep>/common.go` | `:26-84 WebView 接口`：`:30 Run()`、`:34 Terminate()`、`:39 Dispatch(f func())`、`:41 Destroy()`、`:47 Window()`、`:51 SetTitle`（"Must be called from the UI thread"）、`:54 SetSize`、`:59 Navigate`、`:63 SetHtml`、`:68 Init`、`:73 Eval`、`:83 Bind` |
| R14 | `grep -rn "func (e \*Chromium) Resize\|PutBounds" <dep> --include=*.go` | `pkg/edge/chromium_386.go:11`、`chromium_amd64.go:12`（`:18 e.controller.vtbl.PutBounds.Call`）、`chromium_arm64.go:12`；`ICoreWebView2Controller.go:15`、`:59 func (i *ICoreWebView2Controller) PutBounds(bounds w32.Rect) error`、`ICoreWebView2Controller2.go:14` |
| R15 | `grep -rn bringUp --include=*.go .` | 本仓 3 枚文件命中（`cmd/wisp/panel_host_windows.go:22,25,101,164,174,261`；`cmd/wisp/panel_host_windows_test.go:32,39,49,166,167`；`scripts/spike/webview2-latency/main.go:163,164,200,233,249`）；**库侧零命中** |
| R16 | `grep -n "" cmd/wisp/panel_host_windows.go \| sed -n '1,120p'`／`'150,300p'`／`'300,430p'` | 411 行文件全貌：`:5-45` 头部叙述（`:22-31` 再入自述、`:33-45` `w32.Rect` 缺口）、`:55` webview import、`:62` 绑定名、`:70-84` user32 proc 名册、`:89-97 pnlMsg`、`:103-120 PanelManager`、`:153-162 LastColdMs/LastHotMs/windowHandle`、`:164-233 bringUp`（`:183-192 NewWithOptions`、`:193-195 nil→error`、`:197 Window()`、`:204-211 Bind`、`:221 serveEntry`、`:227 firstRoundTripLocked`）、`:236-252 serveEntry`（`:250 SetHtml`）、`:256-296 Show/HotShow`、`:300-330 Hide/Destroy`、`:334-342 dispatchRaw`、`:349-394 firstRoundTripLocked`（`:361-368` 探针绑定、`:372-375` 探针 `SetHtml`、`:377-393` 泵循环＋`:377 deadline := t0.Add(5 * time.Second)`）、`:396-411 pnlPumpOnce` |
| R17 | `grep -n "\.Run()\|Terminate()\|Dispatch(" cmd/wisp/panel_host_windows.go cmd/wisp/panel_host_windows_test.go` | **零命中**（`rc=1`）⇒ 宿主从不跑库的主泵、从不调 `Dispatch`、不调 `Terminate` |
| R18 | `grep -rn "\.Dispatch(" --include=*.go cmd internal \| grep -v _test` | 只命中 statemachine／models 的 `Dispatch`；**webview2 的 `Dispatch` 零调用者** |
| R19 | `grep -rn "staThread\|PostTask\|type staThread\|b.sta" internal/ball/*.go` | 未导出面与 25+ 调用点（③-1） |
| R20 | `grep -rn "^func (b \*Ball) [A-Z]\|^func (s \*staThread) [A-Z]" internal/ball/*.go` | **全目录**导出名册（③-1）；`staThread.PostTask:127` 是唯一 receiver 未导出者 |
| R21 | `for f in internal/ball/*_test.go; do grep -m1 "^package " "$f"; done` | 11 枚**全部** `package ball` ⇒ F5 结清 |
| R22 | `grep -n "^func (b \*Ball)\|^type \|^func " internal/ball/ball_windows.go` | `:33 EventKind`、`:49 Events`、`:63 Options`、`:82 Ball`、`:143 New`、`:182 createOnSTA`、`:540 ballWndProc`、`:549 wndProc`、`:728 fire`、`:741 uiRun`、`:935 Close`、`:1003 DebugHWND` |
| R23 | `grep -n "" internal/ball/sta_windows.go`（全文 162 行） | `:5-9` D38a 契约注释（"…and later the panel"）、`:24-36`（`:31 hwnd windows.HWND // ball window`）、`:38-44`、`:47-96`（`:52-53 LockOSThread`、`:54-55` post-and-wait 死锁警告、`:62 CoInitializeEx(coinitApartmentThreaded)`、`:65-72 create/close(started)`、`:78-93` 主泵）、`:103-108 releaseCOM`、`:111-116 waitStarted`、`:120-124 threadID`、`:127-143 PostTask`（`:134-141 hwnd==0 → delete+return`）、`:146-154 runTask`、`:156-161 quit`（`:161 pPostQuitMessage.Call(0)`） |
| R24 | `grep -n "" internal/ball/ball_windows.go \| sed -n '720,760p'` | `:724-727` `fire` 契约原文＋`:728-732`（无投递、无 recover）；`:734-740 uiRun` 注释原文＋`:741-752`（`:742 GetCurrentThreadId()==b.sta.threadID()`） |
| R25 | `grep -rn "wmAppTask\|wmNull =\|wmApp =" internal/ball/*.go`；`grep -n "" internal/ball/win32_windows.go`（局部） | `win32_windows.go:83 wmNull = 0x0000`、`:85 wmTimer = 0x0113`、`:99 wmApp = 0x8000`、`:102 wmAppTask = wmApp + 0x201`；消费点 `ball_windows.go:688-689` |
| R26 | `grep -rn "OnPanelHotkey\|OnTrayPanel\|opts.Events\." internal/ball/*.go \| grep -v _test` | `ball_windows.go:54-55` 字段声明；`:659-683` 逐枚 `b.fire(...)`；面板两枚＝`:665`、`:672`/`:677` |
| R27 | `grep -n "" cmd/wisp/resident_ball_windows.go \| sed -n '115,140p'`；`grep -n "ball.New(" cmd/wisp/*.go \| grep -v _test` | `:114 b, err := ball.New(ball.Options{`（全仓唯一产码调用点）；`:122-133 Events` 字面量（`:127 OnPanelHotkey: func() { recordBallGesture("panel-hotkey") }`、`:128 OnTrayPanel: func() { recordBallGesture("tray-open-panel") }`）；`:206-211 recordBallGesture` |
| R28 | `grep -rn "ui-sta" --include=*.go internal/observe/`；`grep -rn ResidentBaseline --include=*.go .` | `goroutine.go:38`（"D38b: ui-sta, …"）、`:40 const ResidentBaseline = 6`、`:44` 六枚名册、`:422 ResidentOverBaseline`；`goroutine_test.go:243`、`:268-269` |
| R29 | `grep -n "func (r \*Registry) Spawn" -A 30 internal/observe/goroutine.go`；`sed -n '292,330p'` | `:262-283 Spawn`（`:269 ClassifyGoroutine`、`:270-273 CategoryUnknown ⇒ slog.Warn`、`:281 go r.run(...)`）；`:286 run`；**`:294-326 defer/recover`**（`:296 if rec := recover(); rec != nil`、`:302 Stack: string(debug.Stack())`、`:309 r.panicCount++`、`:314-316 sink(ev)`、`:317-320` 错误文本 `"panic in goroutine %q (owner %s): %s"`、`:323-325 root.Cancel()`） |
| R30 | `grep -rn "panicCount\|PanicEvent" --include=*.go .`（局部） | 在册仪器只有 `internal/observe/goroutine_test.go:116-117 reg.SetPanicSink(...)` ⇒ **全仓没有一枚"⛔ 不许有 panic"的正向钉**〔量：零额外命中〕 |
| R31 | `grep -rn "panic" .scratch/wisp/probes/33/`；`find .scratch/wisp/probes/33/ -type f \| sort` | 15 枚命中全是 `gate-clauses.txt` 的源码回声；84 枚文件全名单（②-3） |
| R32 | `grep -rn "panic\|⑧" docs/evidence/s1/33-panel-host-c27-r1.md` | `:7,28,33,41,45,46,53,66,73`（②-5 逐字） |
| R33 | `grep -n "格\|AC#" docs/evidence/s1/33-panel-host-c27-v1.md`（219 行表） | `:154-172` 六格判语逐条（AC#1 部分成立／AC#2 数值成立形状不成立／AC#3 不成立／AC#4 不成立／AC#11 成立／AC#12 今天勾不了）；`:215`、`:216` 总结；`:60`＝#18 那句 panic 定性；`:87`＝#26；`:89`＝#28；`:90`＝#29；`:105`＝#31；`:135`、`:149`＝§C；`:179`＝J2；`:183`＝J6 |
| R34 | `grep -n "func TestAC228..." -A 60 cmd/wisp/resident_ball_228_test.go`；`sed -n '186,275p'`；`sed -n '276,325p'` | F1/F2 全文（`:45-53` 名册与设计意图；`:205-208/:211-213/:246-249/:255-267` 四条款；`:288-305` Events 字面量 walk；`:313-318`、`:319-322` 两枚 Fatalf 原文） |
| R35 | `grep -n "func TestPanelHostIsAttachedAndNamesTheWindowHops" -A 60 internal/panel/composer_dispatch_test.go` | F3 全文：`:433-447` 两枚 Errorf 原文；`:449-480` carrier A/B 正控；`:482-493` carrier C 分辨腿原文 |
| R36 | `grep -n "ReadFile\|ast\.\|func Test" cmd/wisp/panel_host_windows_test.go`；`sed -n '20,70p'`；`sed -n '140,240p'` | 测试侧**一枚** `Test` 函数（`:149`）＋无源码扫；`:30-64` harness 全文（`:41 go func`、`:43-47 recover`、`:48 LockOSThread`、`:51 pnlPumpOnce`、`:59 time.After(30 * time.Second)`）；`:158`＝全仓唯一 `NewPanelManager` 调用者；`:190-197`＝H3/H10 用 Go 直调 `dispatchRaw` |
| R37 | `grep -n "ReadFile\|panel_host_windows.go\|func Test" cmd/wisp/panel_host_gate_test.go` | F4：`:26-38`（`:27 src := filepath.Join("panel_host_windows.go")`、`:38` Errorf 原文）；`:52`、`:96` 另两枚 gate 用例 |
| R38 | `grep -rln "ReadFile\|source walk\|ast.Parse" cmd/wisp/*_test.go` | 20 枚源码走查类用例名册；其中与票 33 直接相关＝`panel_host_gate_test.go`、`panel_inbound_33_test.go`、`resident_ball_228_test.go` |
| R39 | `grep -n "func Test\|t.Errorf" cmd/wisp/panel_inbound_33_test.go` | `:182-196 TestAC9ComposerDispatchHasAProductionCaller`（要求非测试文件送 `Handle`＋`main.go` 有 `panel-inbound` case＋usage 块写该行）⇒ 甲形不触〔推〕 |
| R40 | `grep -n "^func \|LockOSThread\|CreateWindowEx\|Shell_NotifyIcon" cmd/wisp/notify_windows.go`；`sed -n '130,195p'` | `:36 hwndMessage = ^windows.HWND(2)`；`postSystemNotification:137-186`＝`:138 LockOSThread`、`:149 CreateWindowExW(message-only)`、`:155 defer DestroyWindow`、`:172-179 Shell_NotifyIconW`、`:183 time.Sleep(notifyBalloonTimeout)`；**泵调用零命中**（`grep -n "GetMessageW\|PeekMessageW\|DispatchMessageW"` ⇒ 无输出） |
| R41 | `grep -rn "^func \|CreateWindowEx\|Shell_NotifyIcon" internal/ball/tray_windows.go`；`ls internal/ball/` | `:30 addTrayIcon(hwnd windows.HWND, tip string)`、`:55 setTip`、`:62 remove`、`:72 showMenu` ⇒ tray **复用球窗句柄**；包内 29 枚文件名册（含 `sta_windows.go`、`win32_windows.go`、`renderer_windows.go`、`d2d_windows.go`） |
| R42 | `grep -rn "\.Run()" <dep> --include=*.go` | **唯一命中 `cmd/demo/main.go:27`** ⇒ 库内不自跑主泵；结清 ⑥-6 |
| R43 | `grep -rn "PostThreadMessage\|User32PostMessageW\|PostMessageW" <dep> --include=*.go` | 4 枚命中＝`internal/w32/w32.go:36`（proc 声明）、`:38`（proc 声明）、`webview.go:348`（`PostMessageW(hwnd, WMClose)`）、`webview.go:447`（`PostThreadMessageW(mainthread, WMApp)`）⇒ **COM 回调路径零枚投递**；结清 ⑥-3 的前提 |
| R44 | `grep -rn "createCoreWebView2EnvironmentWithOptions" <dep> --include=*.go`；`grep -n "" <dep>/webviewloader/module.go \| head -60` | 声明＝`pkg/edge/corewebview2.go:50`（唯一调用点 `pkg/edge/chromium.go:87`）；loader 侧＝`webviewloader/module.go:13-14 nativeModule = windows.NewLazyDLL("WebView2Loader")`、`nativeCreate = ...NewProc("CreateCoreWebView2EnvironmentWithOptions")`、`:18-23` 双路径（磁盘 DLL / `go-winloader` 从内存）〔量〕 |
| R45 | `grep -on 's\.add("[a-z0-9-]*"' tools/d22scan/main.go \| sed 's/.*add(//' \| sort -u`；`grep -n "\.go\"\|HasSuffix\|filepath.WalkDir\|skip" tools/d22scan/main.go` | **规则全名册 7 枚**：bare-goroutine／emoji／mirror-hash／pathresolver-bypass／plaintext-key／unparseable／wallclock-timeout；扫描面＝`:655 WalkDir`、`:660` 跳 `testdata`/`.git`/ignore、`:665` 与 `:858` **`!HasSuffix(".go") \|\| HasSuffix("_test.go")` ⇒ 只看非测试 `.go`**〔量〕；结清 ⑦-D |
| R46 | `grep -n "OnPanelHotkey\|OnTrayPanel\|ball.Events" docs/PLAN.md docs/specs/*.md`；`grep -n "^| C\|^| \*\*C" docs/PLAN.md`（挑球/面板/STA 行） | 前者**零命中**；后者命中 `:1362 C12 BallState`、`:1367 C17 PanelBridge`、`:1368 C18 ApprovalQueue`、`:1370 C20`、`:1373 C23`、`:1374 C24`、`:1377 C27 PanelManager`、`:1381 C31 SessionScope` 等〔量，逐字引文见 ⑦-F〕；结清 ⑦-F |
| R47 | `grep -rn "panicCount\|PanicEvent" --include=*.go .`（全仓，非局部） | 除 `internal/observe/goroutine{,_test}.go` 自身外**零枚在册仪器断言"不许有 panic"** ⇒ 甲形若在 `ui-sta` 上炸，**今天没有一枚钉会主动叫红**（只有 `resident_ball_228_test.go` 的存活/记账类断言可能间接受影响）。〔量＋推〕 |

---

## ⑥ 我可能写错的条目（对抗我自己／终态）

- **⑥-1 派单把 `bringUp` 说成"那枚依赖的泵模型"**。〔量：库侧 `grep -rn bringUp` 零命中；真身＝`cmd/wisp/panel_host_windows.go:174`〕⇒ ① 全程分开叙述。若混读，那条叙述错。
- **⑥-2 派单给的库路径 `.../go-webview2@…/chromium.go` 不存在**（真身 `pkg/edge/chromium.go`），**而号数 `:96-111`/`:130-131` 在真身里成立**。⇒ 复核者按字面路径会撞空，可能误判"自述不实"。
- **⑥-3 ①-3 的三支排序是我的推断，不是读数**。⚠ **同轮已补两枚支撑读数**（R42/R43）：`createCoreWebView2EnvironmentWithOptions` 的真身在 `pkg/edge/corewebview2.go:50`，它把 `e.envCompleted` **按指针交给原生 `WebView2Loader.dll`**（`webviewloader/module.go:13-14 nativeModule/nativeCreate`）〔量〕；而全依赖里**仅有的两枚跨线程投递**是 `webview.go:348 PostMessageW(hwnd, WMClose)` 与 `webview.go:447 PostThreadMessageW(mainthread, WMApp)`——**COM 回调路径上一枚都没有**〔量：`grep -rn "PostThreadMessage|User32PostMessageW|PostMessageW" <dep>` ⇒ 只 4 行命中，含 `w32.go:36/:38` 的 proc 声明〕。
  ⇒ 所以**第 1 支（同线程嵌套泵自洽）有源码支撑**：回调只能靠该线程的泵送达。但**"再入一定不死锁"仍没证**（原生 loader 是否把回调投给它自己的窗口，我读不到）。这一条不再"足以推翻 ①"，降级为"仍能推翻 ①-3 的第 3 支触发条件"。〔量＋推〕
- **⑥-4 `sta_windows.go:134-141` 我判"注释与行为不符"**（注释逐字 "run inline as last resort so callers cannot deadlock on a dead thread"，代码逐字 `delete`＋`return`）。〔量：R23 全文〕⇒ 这条"不符"是**我下的判语**。若我把 `PostTask` 与 `uiRun` 的注释搞混（`uiRun:747` 也调 `PostTask`），这条会过头。**复核尺：只看 `:134-141` 五字节。**
- **⑥-5 "线程消息不经 `DispatchMessageW` 送达窗口过程"**——这是 Win32 文档语义，**不是本仓/本依赖的源码读数**〔标注为通识／推〕。我据它得出"甲形之外还要解 `dispatchq`"。**若这条通识我记反了（例如 `DispatchMessageW` 真的会把 `WM_APP` 送进 `webview` 窗口过程），则 ①-2 与 ④ 的"唯一能兑现回执的形状＝`Run()`"整段作废。** 这是本文件**最重的一枚自我怀疑**，且**能被真跑一次结清**（⑦-H）。
- **⑥-6 我在 ①-2 说"本仓宿主从不跑 `Run()`"**，判据是两枚文件内的 `grep` 零命中〔量：R17〕。**同轮已补全依赖尺**：`grep -rn "\.Run()" <dep>` ⇒ **唯一命中 `cmd/demo/main.go:27`**（依赖自己的 demo，不被本模块 import）〔量：R42〕⇒ `NewWithOptions` 内部不自跑主泵，这条推成立的前提已核。**仍未核的是**"原生 loader 的回调一定投在该线程"（只有 Win32/WebView2 文档支撑＋库侧零投递的反证〔量：R43〕）。
- **⑥-7 `ResidentNames` 我只读了 `goroutine.go:44` 一行**〔量：该行确实含 6 个字符串〕。若名册跨行拼接，"六枚"叙述会偏。`goroutine_test.go:268-269` 恰好钉住 `ResidentBaseline == len(ResidentNames)`，所以**六枚是在册事实**，但我没读枚举全行的原文。
- **⑥-8 骨架版担心的"单文件漏计导出面"已结清**：R20 是全目录尺。⇒ ③-1 的"零导出投递面"现在是〔量〕。
- **⑥-9 骨架版 ⑥-6 的"包内测试会不会是 `package ball_test`"已结清**：R21 ⇒ 11 枚全 `package ball`。
- **⑥-10 F3 的能力钉我判"甲形不红"**，判据是它扫"标识符／import"（`:439-441`）〔量〕。**没量**的是 `hostChannelCapabilityHits` 是否**只看 `cmd/wisp` 与 `internal/panel`**、还是全仓——若它另有"必须落在具名文件"的隐含口径，搬家就可能红。⇒ 复核尺＝读 `hostChannelCapabilityHits`/`hostSignals` 两枚 helper 原文（**我没读**）。
- **⑥-11 ③-3 F6 的行为型红是〔推〕**：我只量到那些用例的取数形状（`PostTask → done <-`），**没读它们的等待/超时与断言全貌**，更没跑。真跑才知道会不会乱序红。
- **⑥-12 已结清**：d22scan 的**规则全名册**我抽出来了（尺＝`grep -on 's\.add("[a-z0-9-]*"' tools/d22scan/main.go`，读数见 R45／⑦-D）⇒ **7 枚规则、扫描面只含非测试 `.go`**（`:665`／`:858` 那两枚 `HasSuffix` 判据），不存在"扫导出命名／扫 `LockOSThread`／扫注释词面"这类规则。**骨架版那句"没读规则全集"的欠账作废。** 仍未量的只剩：我⛔不许跑这枚尺，所以"甲形过 d22scan"是**按规则名册逐条比对**的〔推〕、不是读数。
- **⑥-13 枚数与时刻**：②-3 的"84 枚"跑在建 `a2` 之前（不含本文件）；R38 的"20 枚"是 `cmd/wisp/*_test.go` 里含 `ReadFile|ast.Parse` 的文件数，不是"走查类用例"数。两处口径不同名，别混用。
- **⑥-14 我通篇没读 `frontend/**` 与 `design/**`**（两层禁令）。凡是"页面能不能收到回执"的现实后果，我只从 Go 侧源码说；**不引前端任何内容**，也不据前端形状做判。

---

## ⑦ 判不动的地方（逐条 甲／乙／不做 ＋现量＋为什么判不了）

- **⑦-A（原任务书那一格）：再入到底 panic、阻塞、还是乱序？**
  现量：`pkg/edge/chromium.go:95-111`/`:130-136`〔量〕、`sta_windows.go:78-92`＋`:62`〔量〕、`ball_windows.go:665/672/677`＋`fire:728-732` 无 recover〔量〕、`goroutine.go:294-326` 有 recover＋`debug.Stack()`＋sink〔量〕、`.scratch/wisp/probes/33/` 零工件〔量〕。
  为什么判不了：**只有真跑能定案**，而本腿⛔禁跑（`33-r3` 在 `cmd/wisp` 跑测试）。
  ⇒ **交裁（本腿唯一真正需要你点的格）**：**甲**＝派一枚**允许跑**的探针腿，形状现成：**不写新用例**，直接在**常驻腿自己那条 `ui-sta`**（`Registry.Spawn` 起的）上触发一次面板创建，让 `observe` 的 `SetPanicSink` 把 `debug.Stack()` 落盘 ⇒ **一次跑同时拿到"哪一支"与"原始栈"两半**，且不需要动 `internal/ball` 写面之外多少东西（探针⛔不许进 `cmd/wisp` 测试包，避开 `33-r3`）。**乙**＝接受"未证"直接按甲形写；**不做**＝继续让生产调用者＝0（今天已经是这状态，代价＝owner 那句"我要能看到主面板"继续没有路径）。
- **⑦-B（骨架版挂的，现已结清一半）**：`cmd/wisp` 的宿主测试里**没有**词面型负向钉（R36：`grep "ReadFile|ast.|os.Open"` 在 `panel_host_windows_test.go` 零命中）。**唯一按文件名钉宿主的是 `panel_host_gate_test.go:27`（F4）**。⇒ 结论改为：**甲形不会因注释被红**，但**会因搬家把 L1 掏成假绿**。这一格我已答，⛔ 不需要你裁，需要你把它写进派单判据。
- **⑦-C（已结清）**：全目录尺 R20/R21 已证"零导出投递面"与"测试全在包内"。⇒ 本格作废，读数搬进 ③-1/F5。
- **⑦-D（已结清）：`tools/d22scan` 会不会扫到甲形的新导出？**
  现量：规则名册（尺＝`grep -on 's\.add("[a-z0-9-]*"' tools/d22scan/main.go | sed 's/.*add(//' | sort -u`）⇒ **全 7 枚**：`bare-goroutine`、`emoji`、`mirror-hash`、`pathresolver-bypass`、`plaintext-key`、`unparseable`、`wallclock-timeout`〔量〕。
  扫描面（`:655 WalkDir`、`:660` 跳 `testdata`/`.git`/ignore、`:665` **`if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") { …`** ⇒ **只看非测试 `.go`**、`:858` 同一条走第二枚 walk）〔量〕。
  ⇒ **判定：甲形（一枚导出方法、不起协程、不碰 `filepath.Clean|Abs`、不写凭据词、不用墙钟差做超时、界面零 emoji）不被任何一枚在册规则叫红**〔量＋推：`bare-goroutine` 的红句原文在 `:707/:711`，只钉 `go <anything>`〕。
  ⚠ 附带量到的一件事与本票有关：**d22scan 不扫 `.go` 以外的文件**⇒ 我这份 `.md` 里满屏的 `⛔`/`⚠`（U+26D4/U+26A0，落在 `emojiRe` 的 `U+2600–U+27BF` 带内）**不会**被仪器抓，且票面规矩本来就只禁 UI 代码里的 emoji（AGENTS.md §1.2）。
- **⑦-E（已结清）**：`cmd/wisp/notify_windows.go` **不泵**（R40）⇒ 不是第四形；tray 复用球窗句柄（R41）⇒ 不是第四形。④ 表已定稿。
- **⑦-F（已结清）**：甲形要不要动 `Events`？动了算不算改契约？
  现量：`Events` 是导出结构（`ball_windows.go:49-61`），10 枚字段与 `resident_ball_228_test.go:50-53` 的字面名册一一对应；`grep -n "OnPanelHotkey\|OnTrayPanel\|ball.Events" docs/PLAN.md docs/specs/*.md` ⇒ **零命中**〔量〕。
  **契约表侧我同轮补读了**（尺＝`grep -n "^| C\|^| \*\*C" docs/PLAN.md` 内挑球/面板/STA 那几行）〔量，PLAN.md 行号〕：
  - `:1362 C12 BallState`（20 态＋D43 转移表）、`:1367 **C17** PanelBridge`（"前端↔Go 双向通道：`invoke(method, args) → result` ＋ **Go→前端事件推送**；**回复必须按 correlationId 路由**；前端必须无状态"）、`:1377 **C27** PanelManager`（"单例，**唯一** WebView2 窗口持有者；一会话至多一个面板窗口，**隐藏而非销毁**；面板不可用 → **L2 降级为原生最简确认卡，L2 能力不得消失**"）、`:1381 C31 SessionScope`（"…面板保活…会话结束必须完整 Dispose"）。
  ⇒ **判定（本格已答）**：`ball.Events` **不在 C1–C32 任何一行的名字里**，也不是 D 表标题项〔量：零命中〕⇒ 动它属**实现面**，不是改契约；**但它会当场叫红 F1 那枚在册名册钉**（词面型、设计如此，红句自己说"the instrument, not the code"）。
  ⇒ **顺带量到两枚对本票更硬的约束**（交我裁时请一并读）：**C27 逐字要求"唯一 WebView2 窗口持有者"**〔量：`:1377`〕——乙形那枚"第二线程自己泵"**并不违反**这句（它约束的是持有者枚数，不是线程数），但**C17 逐字要求"Go→前端事件推送"与"回复必须按 correlationId 路由"**〔量：`:1367`〕⇒ **⑦-H 那枚"回执那一跳今天在这套形状里根本没有投递者"不是实现细节，是 C17 的契约面未兑现**。⚠ 我把这一句写成〔量＋推〕：契约文字是真读数，"⑦-H 属 C17 射程"是我的归口判断。
- **⑦-G（已结清，且与 `v1` 表同判）**：`r1` 表 `:33-45` 那句"`PutBounds` 吃模块内私有 `w32.Rect` ⇒ 给不了尺寸 ⇒ 只能换依赖"——现量：`pkg/edge/chromium_amd64.go:12 func (e *Chromium) Resize()` **是导出的**，`:18` 内部自己调 `e.controller.vtbl.PutBounds.Call(...)`〔量〕；高层 `webview.go:343` 在 `Embed` 成功后就调它。⇒ **`v1` 表 §AC#3（`:169`）那句"理由不成立：`(*edge.Chromium).Resize()` 是导出的"我独立复认成立**〔量〕。本格不再需要甲／乙；**剩下唯一欠的是"Resize 按父窗尺寸走，能不能满足 420×260 的面板尺寸"**＝〔推〕未证，且属丙形射程。
- **⑦-H（本腿新添、我认为比 ⑦-A 更要你先看的一格）：`dispatchq` 那一跳。**
  现量：`webview.go:445/:447` 写＋`PostThreadMessageW(mainthread, WMApp)`；`dispatchq` 唯一读者＝`webview.go:362-363`（在 `Run()` 内）；本仓宿主零枚 `Run()`/`Dispatch()` 调用者（R17/R18）。〔量〕
  为什么判不了：结论依赖"线程消息不会被 `DispatchMessageW` 送进窗口过程"这条 Win32 语义〔推／通识〕，**我没有本机读数**。
  ⇒ **甲**＝派一枚**只读前端无关**的探针：让真页面调一次 `window.wispDispatch(...)` 并 `await`，看回执到不到（这一发同时结清 ⑥-5、⑥-6 与 AC#1 的"Go→页"）；**乙**＝现在就认定"回执那一跳必须靠 `Run()` 形或改依赖"，把丙形的一个子选项（fork/`replace` 加一枚 drain）摆上；**不做**＝先接"看得见窗口"、把回执留成下一格（⚠ 这正是 owner 那句要求的最低标准，代价＝AC#1 的判据仍不成立）。
  ⛔ **我不替你选。**

---

## 交件判语

**射程**：四格全给了带 `file:line` 的读数。
① 库侧：创建链 13 处逐字定位，线程要求（无 apartment 判定／STA 由调用方负责／`mainthread` 在入口线程号被记死／**`dispatchq` 只有 `Run()` 会读**）＋ 三种最坏后果形状（阻塞／再入窃队列／nil 解引用 panic），**并把"panic 在出货路径上会被 `observe.Registry` 兜住、栈能自动落盘"这一枚形状量出来了**（这是本腿对 ⑦-A 最实际的贡献：结清它不需要新仪器，只需要一条真跑的常驻路径）。
② 工件面：**零工件**，具名 6 枚搜索路径；并把"没有原始栈"从"可能没跑过"细化成"走 harness 那条形本来就产不出栈"（`panel_host_windows_test.go:45`）。
③ 甲形射程：全目录证明"零导出投递面"；四枚候选形状；**会被叫红的用例逐枚点名**——F1 `TestAC228BallHostAnswersEveryGesture`（**词面型名册钉，只被甲-4 打红，断言原文两枚 + 设计意图原文**）、F2（行为＋形状混合，禁动 `ball.New`/`defer stop`）、F3 能力型正向钉（搬家要留载体）、F4 **唯一按文件名寻址的负向钉**（甲形不打红它但会把它掏成假绿——这条必须进派单）、F5 结清包内可见性、F6/F7 行为型钉（要真跑）、F8/F9 只有乙形打红的冻结名册（只登记、不建议）。
④ 第四形：7 枚现成形状逐枚判＋依据；`notify_windows.go`/tray/第二枚 `staThread` 三枚**排除**（有零命中读数支撑）；真正未列的形状浮出＝"让宿主拥有自己的泵（`webview.Run()` 那一形）"，而它与 `ui-sta` 同线程互斥。

**没答的格子（⛔ 不许读成答了）**：
⑦-A（panic／阻塞／乱序＝纯〔推〕，非证）、⑦-H 的 Win32 语义那一半（〔推〕：`DispatchMessageW` 对 `hwnd==NULL` 的线程消息不做派发＝通识，不是盘上读数）、⑥-3 剩余那一半（"原生 loader 的回调一定投在调用线程"＝只有 R43 的"库侧零投递"反证）、⑥-10（`hostChannelCapabilityHits`／`hostSignals` 两枚 helper 原文我**没读**，所以 F3"搬家会不会红"仍是〔推〕）。
①-3 整节、⑥-5、⑥-6、⑦-A、⑦-H 都**没有一发运行读数**支撑。
**骨架版挂的四格本轮已结清为〔量〕**：⑦-C（R20/R21）、⑦-E（R40/R41）、⑦-D（R45）、⑦-F（R46）；**F1 那枚在册名册钉＝本发最要紧的产出**（甲-4 必红、红句自称"the instrument, not the code"）。
另：⛔ 本腿**没给任何选型建议**；全文只出现"形状＋代价＋会不会叫红＋只有真跑能定的那一半在哪"。
