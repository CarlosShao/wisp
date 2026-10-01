# 票 33 · 腿 `33-a2` · 只读普查：面板宿主接进常驻腿的阻塞面（甲／乙／丙选型的前置量尺）

⛔ 本腿零产码／零脚本／零测试／零构建。尺具只有 `Read`/`Grep`/`Glob`/`sed -n`/`find`/`go env`/`git log`。
〔量〕＝我这发跑出来的逐字读数。〔推〕＝我从源码结构推的，**没有真跑过，不算已证**。

## 起手锚点

- `date "+%Y-%m-%d %H:%M:%S %z"` ⇒ `2026-10-01 11:21:14 +0800`〔量〕
- `git log -1 --format=%H` ⇒ `22272ac2b06470b9a03b9acc811eb6d0de944423`〔量〕
- `git branch --show-current` ⇒ `dev`〔量〕
- 依赖版本（`grep -n webview2 go.mod`）⇒ `go.mod:19: github.com/jchv/go-webview2 v0.0.0-20260205173254-56598839c808 // indirect`〔量〕（模块名逐字：`github.com/jchv/go-webview2`）
- module cache（`go env GOMODCACHE`）⇒ `D:\work\base\gopath\pkg\mod`，库根＝
  `D:/work/base/gopath/pkg/mod/github.com/jchv/go-webview2@v0.0.0-20260205173254-56598839c808/`〔量〕
- 本腿只写这一枚文件：`.scratch/wisp/probes/33/a2/census.md`。

---

## ① 库侧到底会发生什么（创建控件 → 泵 → 收到 web message 的线程要求）

### ①-1 链路逐处（全部〔量〕，行号由 `grep -n "" <file> | sed -n 'a,bp'` 反取，行号来自 grep）

库根下**没有**根级 `chromium.go`（我按派单里的路径试了一次：`sed: can't read .../chromium.go: No such file or directory`〔量〕）。
真正的文件是 **`pkg/edge/chromium.go`**；派单说的 `:96-111` 与 `:130-131` 在那枚文件里**号数正好对得上**（见下）。

创建链（高层 API 一形）：

| 处 | 位置 | 干什么 | 线程要求／形状 |
|---|---|---|---|
| a | `webview.go:97` `NewWithOptions` | 入口 | — |
| b | `webview.go:108` | `w.mainthread, _, _ = w32.Kernel32GetCurrentThreadID.Call()` | **调用方线程被记成 `mainthread`**〔量〕 |
| c | `webview.go:109` | `w.CreateWithOptions(...)` | 同步 |
| d | `webview.go:284-293` | 注册窗口类 `"webview"`（**类名是定死的字面量**）、`User32RegisterClassExW` | 返回值被丢（`_, _, _ =`）⇒ 第二枚窗口重复注册同名的失败**不会**被看见〔量〕 |
| e | `webview.go:320-333` | `User32CreateWindowExW` | 窗口**归属调用方线程**〔推，Win32 通识＋b 行〕 |
| f | `webview.go:340` | `w.browser.Embed(w.hwnd)` | ⬇ 落到 edge |
| g | `pkg/edge/chromium.go:87` | `createCoreWebView2EnvironmentWithOptions(nil, datapath, 0, e.envCompleted)` | 异步投递＋回调；`res != 0` 或 `err != nil` ⇒ `return false`（`:88-94`）〔量〕 |
| h | `pkg/edge/chromium.go:95-111` | `for { if atomic.LoadUintptr(&e.inited) != 0 { break }; GetMessageW(&msg, 0, 0, 0); if r == 0 { break }; TranslateMessage; DispatchMessageW }` | **一枚阻塞式嵌套泵**，`hwnd=0`＝取该线程队列上**所有**窗口＋线程消息〔量〕 |
| i | `pkg/edge/chromium.go:224` | `atomic.StoreUintptr(&e.inited, 1)` | 只在 `CreateCoreWebView2ControllerCompleted`（`:186`）里置位 ⇒ **必须靠 h 那枚泵把 COM 回调送下来**〔量＋推〕 |
| j | `pkg/edge/chromium.go:112` | `e.Init("window.external={invoke:...}")` | h 退出后**无条件**执行〔量〕 |
| k | `pkg/edge/chromium.go:130-131` | `_, _, _ = e.webview.vtbl.AddScriptToExecuteOnDocumentCreated.Call(uintptr(unsafe.Pointer(e.webview)), ...)` | **`e.webview` 在这条路径上没有任何 nil 判定**〔量：`:130-136` 全文无 `if`〕 |
| l | `webview.go:113-126` | `chromium.GetSettings()` ⇒ `settings.PutAreDefaultContextMenusEnabled`／`PutAreDevToolsEnabled`，两处 `err != nil` ⇒ **`log.Fatal`** | 见 ①-4 |

web message 侧（收到消息时跑在谁的线程上）：

- `pkg/edge/chromium.go:233-248` `MessageReceived` ⇒ 直接调 `e.MessageCallback(...)`（`:240`）。
  该回调在 `webview.go:103` 被接成 `w.msgcb`。〔量〕
- `webview.go:139-160` `msgcb` ⇒ `w.callbinding(d)`（`:162-220`，**反射直接 `v.Call(args)`，`:191`**），随后所有回执走 `w.Dispatch(func(){ w.Eval(...) })`（`:148`/`:152`/`:156`）。〔量〕
  ⇒ **回调线程＝COM 把消息投下来时正在泵的那枚线程**〔推，形状如上：`MessageReceived` 里没有投递／换线程动作，`w.mainthread` 在 msgcb 路径上未被读〕。

### ①-2 有没有 apartment／`CoInitialize` 的线程判定？

〔量〕在整棵依赖里 `grep` 到的 COM 初始化**只有** `pkg/edge/comproc*.go` 与 `webviewloader` 的 DLL 侧动作；库侧**没有** `CoInitializeEx` 调用、**没有** `APARTMENTSTATE`／`COINIT_*` 常量、**没有**任何"当前线程是不是 STA"的判定，也**没有**任何 `runtime.LockOSThread`。〔推：所以"STA 还是 MTA"是**调用方的义务**，库不校验、不补救。〕

⇒ 对**本仓的形状**的直接影响：`internal/ball/sta_windows.go:62` 已经用 `pCoInitializeEx.Call(0, coinitApartmentThreaded)` 把 `ui-sta` 初始化成 **STA**〔量〕，所以「STA 要求」这一半**甲形本来就满足**（`ui-sta` 就是 `D38a` 那枚"ONE thread owns all UI COM objects … and later the panel"，注释逐字：`sta_windows.go:5-9`〔量〕）。

### ①-3 那一形再入时，最坏后果的形状（三种，按源码结构排）

⚠ 下面**全是〔推〕**，因为本腿⛔不许跑。**只有真跑才能定案的那一半留在 ⑦-A**。

1. **阻塞而不是 panic（最可能、也最"安静"）**：h 那枚泵把调用方**扣住**，直到 i 置位。若创建链是投在**同一枚线程**上的 COM 回调，则需要泵把回调送下来 ⇒ 单线程内**自洽可解**（嵌套泵自己送），后果＝**调用方在 `DispatchMessageW` 里被再入一层泵**，期间外层泵不再前进。〔推〕
2. **再入导致的消息窃取／顺序错乱（甲形最要命的形状）**：h 用 `GetMessageW(…, 0, 0, 0)`，`hwnd=0` 的语义是**该线程队列上所有窗口**＋线程消息。而球的 `ui-sta` 靠**同一个队列**驱动：`sta_windows.go:127-143 PostTask` 把闭包登记后用 `pPostMessageW.Call(hwnd, wmAppTask, id, 0)` 投到**球窗口**（`sta_windows.go:31 hwnd windows.HWND // ball window`），由 `sta_windows.go:91-92`→`ball_windows.go:549 wndProc`→`ball_windows.go:689 b.sta.runTask(uint64(wParam))` 消费〔量：这几处行号均由 `grep -n` 反取〕。
   ⇒ 所以**嵌套泵会把球的任务消息也吃掉并派发**（`DispatchMessageW` 会送到球的 wndProc，行为上"看起来仍然对"），但同时会把**外层泵本应顺序处理**的输入/定时器消息插到外层调用还没返回的位置上执行——**球的状态机回调可能在 `bringUp` 还没跑完时就跑在面板上**（`Ball.uiRun`/`fire` 的注释逐字承认这个模型：`ball_windows.go:725` "be quick and non-blocking (statemachine dispatch + PostTask); a blocking"、`sta_windows.go:54-55` "a caller already on it must run inline … post-and-wait from inside the pump is a deadlock"〔量〕）。〔推＝后果形状，未真跑〕
3. **panic（那句自述的因）**：h 有**两条**不以 i 为条件的出口：`:106-108 if r == 0 { break }`（`GetMessageW` 返回 0＝该线程队列收到 **WM_QUIT**）与 `r == 0` 之外的 `err`。⚠ **库不区分 `inited` 是否已置位就 `break`**，接着 `:112` 无条件 `e.Init(...)` ⇒ `:130-131` 解引用 `e.webview`（在 `:194-197` 才被赋值，且赋值用的 `Call` 返回值被丢弃）⇒ **nil 指针解引用＝Go runtime panic**。〔推：形状如上，源码里 `:130-136` 无 nil 判定这一点是〔量〕〕
   - **谁能投出 WM_QUIT**：`internal/ball/ball_windows.go:962 b.sta.quit()` ⇒ `sta_windows.go:156-161` 里 `pPostQuitMessage.Call(0)`；库侧 `webview.go:381-383 Terminate()` 也是 `PostQuitMessage`（`webview.go:243` 在 `WMNCLButtonDown…WMDestroy` 分支调它）。〔量：这些都是行内字面〕
   - **另一枚同族出口**：`EnvironmentCompleted`（`chromium.go:171-174`）与 `CreateCoreWebView2ControllerCompleted`（`:186-189`）里 `res < 0` ⇒ **`log.Fatalf`** ⇒ 进程 `os.Exit(1)`，**不是 panic**、也**不可 recover**。〔量＋推（`log.Fatalf` 的行为是 Go 通识）〕

### ①-4 库侧还有两枚会把"失败"升级成"进程没了"的地方

- `webview.go:113-116`：`settings, err := chromium.GetSettings(); if err != nil { log.Fatal(err) }`
- `webview.go:118-126`：`PutAreDefaultContextMenusEnabled` / `PutAreDevToolsEnabled` 两枚 `err != nil` ⇒ `log.Fatal`
- `pkg/edge/chromium.go:283-285`：`WebResourceRequested` 里 `args.GetRequest()` 出错 ⇒ `log.Fatal(err)`（本票的 `r1` 表已点名这条，见 `docs/evidence/s1/33-panel-host-c27-r1.md:28`〔量〕）

〔推〕⇒ **甲形落地时，"创建失败"这一支的默认库行为是杀进程，不是返回 error**。`PanelManager.bringUp` 里对 `err` 的处理能不能包住这一支＝`cmd/wisp` 的事（本腿不判）。

### ①-5 本格欠账（⛔ 不许读成已证）

只有真跑能定的两半：**(i)** 在**已在泵的同一枚 STA** 上再入 `NewWithOptions` 究竟走 ①-3 的第 1／2／3 哪一支；**(ii)** 走第 3 支时栈顶是不是 `edge.(*Chromium).Init`。
本腿**未跑、也不能跑**（另一枚写腿 `33-r3` 正在 `cmd/wisp` 里跑测试）。⇒ 见 ⑦-A。

---

## ② 那句〔自述〕的工件面

**结论：`.scratch/wisp/probes/33/` 下＝零该发工件。**〔量〕

搜过的路径（逐枚具名）：

1. `ls -R .scratch/wisp/probes/` ⇒ 顶层含 `33`；`ls -R .scratch/wisp/probes/33/` ⇒ **`a1` `h1` `r1` `r2` `r3`** 五枚子目录（本轮新增 `a2`）。〔量〕
2. `find .scratch/wisp/probes/33/ -type f | sort` ⇒ **共 84 枚文件**〔量，清单见 ⑤-R3 摘要〕。其中 `r1/logs/` 只有 `d22scan.txt gate-clauses.txt gotest.txt?`…（准确名单）：`r1/logs/`＝`d22scan.txt, gate-clauses.txt, green-after.txt, green.txt, mutate-run.txt, panel.txt, red-M1.txt, red-M2.txt, red-M3.txt, red-M4.txt, risk.txt`；`r1/` 另有 `composer_dispatch.go.pristine, mutate.sh`；`a1/logs/`＝`d22scan.txt, gate-clauses.txt, gotest.txt`；`h1/`＝`census.md, msg-1to7.txt, msg-8910.txt, msg-skeleton.txt`；`r2/`＝33 枚（含 `mut-1/`、`mut-2/` 两子目录）；`r3/`＝16 枚。〔量〕
   ⇒ **没有一枚 `.log`/`.txt` 的名字或内容形如 panic dump／repro 脚本**（`r1/mutate.sh` 是变异脚本，不是复现命令）。〔量＋推：判据见 3〕
3. `grep -rni "panic" .scratch/wisp/probes/33/`〔量〕⇒ **全部 15 枚命中都来自 `gate-clauses.txt` 那四份日志**（`a1/logs/`、`r1/logs/`、`r2/gate.txt`、`r2/gate-final.txt`、`r3/gate-final.txt`），命中内容逐字只有两类：
   - `internal/memory/retention.go:100:  panic("memory: StartRetentionJob requires a DisposalScope")`
   - `internal/plugin/disposal_test.go:87: s.DeferNamed("panics", func() { panic("destroy failed: injected") })`
   ⇒ 这些是**闸门把仓源码逐行贴进日志**留下的字面，**与 webview2／`bringUp`／`ui-sta` 无关**；**零命中**含 `webview2`、`Init`、`GetMessageW`、`goroutine [number] \[running\]` 之类栈头形状。〔量〕
4. 唯一存在过的"栈"＝**证据表里的一句散文**：`docs/evidence/s1/33-panel-host-c27-r1.md:73` 逐字含
   "本机现量 `edge.Chromium.Init` nil-webview panic（栈：`ballWndProc→staThread.start DispatchMessageW→…→PanelManager.bringUp→webview2.NewWithOptions→CreateWithOptions(webview.go:340)→Embed→Init(chromium.go:131)`）"〔量〕
   同文件 `:45` 与 `:46` 也各有一句"会 panic"／"本机现量 nil-webview panic 栈已入 ②"，`:46` 逐字承认该用例**"未提交该文件"**。〔量〕
   ⚠ 该栈是**符号顺序的散文**、中间用 `…` 省略、无 goroutine 头／无行号列／无 `exit status 2`，**盘上不存在任何原始 dump** ⇒ 定性仍是〔仅自述＋栈形状〕，与本腿任务书给的定性一致。〔量〕

**本格补一句要紧的**：`r1` 那句自述里的 `webview.go:340` 与 `chromium.go:131` **两处行号我这一发都对上了库源码**（`webview.go:340 = w.browser.Embed(w.hwnd)`；`pkg/edge/chromium.go:131` 是 `e.webview.vtbl.AddScriptToExecuteOnDocumentCreated.Call(` 那一行）〔量〕。
⇒ 这**提高了那句自述的可信度**（栈形状与真实代码形状自洽），但**不等于那发 panic 真发生过**——它是"形状可能对"，不是"有工件"。〔推〕

---

## ③ 甲形的真实射程（导出什么才算"最小一枚"＋会被叫红的现成用例）

### ③-1 现有导出面：够不够？〔量〕

`internal/ball/ball_windows.go` 里 `type Ball` 的**导出**方法与类型名册（`grep -n "^func (b \*Ball)\|^type \|^func " internal/ball/ball_windows.go`）：

- 导出方法（逐条）：`SetState:311`、`SetBadge:316`、`SetProgress:324`、`SetBadgeText:332`、`TimersAlive:441`、`DebugTimersAlive:450`、`RebindHotkeys:813`、`HotkeyReport:830`、`ConfiguredHotkeys:840`、`RegisteredHotkeys:849`、`TakeEscForCancel:873`、`ReleaseEscAfterSession:895`、`EscTakenOver:910`、`SetTrayChecks:919`、`SetTrayTip:927`、`Close:935`、`DebugHWND:1003`。
- 导出类型：`EventKind:33`、`Events:49`、`Options:63`、`Ball:82`。导出函数：`New:143`。
- ⛔ **投递面一枚都不导出**：`staThread`（`sta_windows.go:24`）、`newSTAThread:38`、`start:47`、`waitStarted:111`、`threadID:120`、`PostTask:127`、`runTask:146`、`quit:156` 全是小写；`Ball.sta`（`ball_windows.go:84 sta *staThread`）也是小写；`Ball.uiRun:741`、`Ball.fire:728` 也是小写。
  ⇒ **今天从包外没有任何一条路能把闭包投到 `ui-sta`**〔量〕。库侧的 `webview.Dispatch`（`webview.go:148/152/156`）是**库内部**的跨线程投递，且 `w.mainthread` 是它自己的线程，不认球的线程〔量＋推〕。

### ③-2 "最小一枚"的候选面（只列形状，不选形）〔量＋推〕

| 候选 | 形状 | 射程 |
|---|---|---|
| 甲-1 | `func (b *Ball) PostUITask(fn func())`（＝`b.sta.PostTask` 的一层包装） | **一枚方法、零新类型**。调用方只拿到"投上去就走"，**拿不到错误、拿不到完成通知** ⇒ 面板创建失败无法回报（`staThread.PostTask:134-141` 在 `hwnd == 0` 时**静默丢弃**、`:135` 注释逐字 "run inline as last resort" 与代码不符〔量：实际是 `delete` 后 `return`，不跑 inline ⇒ 这是一枚**注释与行为不符**，见 ⑥-4〕）。 |
| 甲-2 | 导出 `STAHandler`／`PostTask` 一枚**接口**＋`b.STA()` getter | 一枚方法＋一枚类型。会把 `staThread` 的形状变成契约面。 |
| 甲-3 | `func (b *Ball) RunOnUIThread(fn func()) error`（**等待完成**） | ⛔ **正是 `sta_windows.go:54-55` 与 `ball_windows.go:725` 两处注释明写的 deadlock 形状**（"post-and-wait from inside the pump is a deadlock"），面板创建又要跑在里面 ⇒ 甲-3 与"从回调里创建"互斥。〔推〕 |
| 甲-4 | 什么都不导出，改在 `Ball` 上加一枚**宿主注册口**（如 `Events` 里加一个"在 ui-sta 上调你"的 hook） | 动 `Events:49` 结构＝动已钉的回调契约面（`Events` 现有字段谁在用要另量）。 |

⇒ **"最小一枚"在形状上＝甲-1**〔量：唯一一枚"一枚方法、零新类型、零新契约"〕。**⛔ 我不选形。**

### ③-3 ⚠ 本格的真正产出：会被叫红的现成用例，逐枚点名

（**判据**：行为型钉＝钉的是"跑起来必须得到什么结果"；词面型名册钉＝钉的是"某枚名字/常量/清单必须等于这个字符串集合"。）

**第一族：`ui-sta` 名册钉（乙形会打红，甲形按上面四枚候选都不会）**

- `internal/observe/goroutine.go:44` 逐字：`"ui-sta", "audio-capture", "hotkey-listener", "db-writer", "watcher…"〔量：完整名册＝`ui-sta`, `audio-capture`, `hotkey-listener`, `db-writer`, `watchdog`, `log-flusher`〕；`goroutine.go:40 const ResidentBaseline = 6`；`goroutine.go:422 rep.ResidentOverBaseline = rep.Resident > ResidentBaseline`。〔量〕
- `internal/observe/goroutine_test.go:268-269`：`if ResidentBaseline != len(ResidentNames) { t.Fatalf(...) }`〔量：词面型名册钉〕
- `internal/observe/goroutine_test.go:243`：`"ui-sta": CategoryResident,`〔量：名册钉的一环〕
- **判定**：**甲形一枚新协程都不起**（复用现成 `ui-sta`）⇒ 这一族**不被叫红**〔推，依据是甲-1/甲-2 的调用链里没有 `Spawn`〕。
  **乙形必叫红**（新增一枚常驻 ⇒ `Resident > 6` ⇒ `ResidentOverBaseline` 为真），且**这一族属人工批准那一档**，本腿不建议。〔推〕

**第二族：`internal/ball` 包内对 `sta`/`PostTask` 的直接依赖**

- 用例（〔量〕名册）：`hotkey_live_test.go:60,327-334,471,490`、`live_guard_windows_test.go:102,114,253,322`、`live_windows_test.go:224,229,383,438,529,534,555,589,658,673`、`interaction_live_test.go:392`（注释）、`dock_test.go:60`（无关）。
- **判定**：这些用例是 `package ball`（包内）⇒ **导出一枚方法不会让它们红**〔推：Go 包内可见性；本腿没读它们的 `package` 行以外证据 ⇒ 该推成立的前提是"文件里没写 `package ball_test`"，我量一下再入账，见 ⑦-C〕。

**第三族：⚠ 会真的叫红的两枚（甲形的代价面，派单里必须当场裁）**

- **3-1 `cmd/wisp/panel_host_windows_test.go:39-49 hostThreadHarness.bringUp`**：测试自持一枚锁线程跑 `bringUp`（逐字注释 `:32` "targets production goroutines, and production bringUp runs on the caller's"）。〔量〕
  **为什么会被叫红**：甲形一旦把创建改到 `ui-sta`，`bringUp` 的**"跑在调用方所在线程"**契约就变了 ⇒ 这枚 harness 的"我起的线程就是宿主线程"这一断不再真。〔推〕
  **性质**：**行为型钉**（钉的是"窗口与泵必须同线程"这一条不变式）。**它不钉名字、不钉清单**，所以改法只有"把 harness 换成投到真 `ui-sta`"或"harness 保留但另立一枚 ui-sta 共驻用例"。
- **3-2 `cmd/wisp/panel_host_windows.go:22-31` 那段注释**（逐字含 "a re-entrant bringUp from inside the ball's ui-sta DispatchMessageW loop **panics on this box**"）与 `:101`（"see bringUp's `runtime.LockOSThread`"）。〔量〕
  **为什么会被叫红**：⛔ 注释**不是断言、不会红**。但 `panel_host_windows_test.go` 里若存在**扫源码词面**的负向钉（例如"生产文件里不许出现 `LockOSThread`"），甲形会当场红。⇒ 这一枚**我还没量到**，见 ⑦-B。**必须派单里当场裁**。
- **3-3 球侧的行为型钉（甲形真正的风险）**：`ball_windows.go:549 wndProc`/`:689 b.sta.runTask(uint64(wParam))` 这一条消费链，若甲形让面板的**嵌套泵**（①-3 第 2 支）插在中间，则**球侧那些"PostTask 后必须按序看到状态"的用例**（`live_windows_test.go:529/534/555/589/658/673` 全部是 `done <- …` 的取数形状、`live_guard_windows_test.go:322 done <- b.curState`）可能被乱序叫红。**性质**：**行为型钉**（钉的是"投上去的闭包在 UI 线程上按序跑到"）。〔推：需要真跑定案〕
- **3-4 `tools/d22scan` 的形状扫**〔量：`.scratch/wisp/probes/33/a1/logs/d22scan.txt` 存在，说明该票每条腿都跑了这枚尺〕：甲形**不新增裸 `go func`**、不新增 `filepath.Clean|Abs`、不动 emoji ⇒ 〔推〕不该红。⛔ 但我**没读 `tools/d22scan/main.go` 的规则全集**（只读了票面提到的 `emojiRe` 那一条来自 AGENTS.md），所以"确实不红"是〔推〕——见 ⑦-D。

**本格结论**：甲形的**代码面代价可能真的一枚导出方法就够**（甲-1），但它会把 `cmd/wisp` 那枚 `hostThreadHarness` 的**契约叙述**（行为型钉）与球侧 `PostTask` 取数用例的**顺序不变式**（行为型钉）**推上桌**；词面型名册钉里只有**乙形**会打红 `internal/observe` 那三处。**"改这处导出面不会叫红任何用例"这句话不成立，但成立的射程取决于 ⑦-B/⑦-C/⑦-D 三枚我还没量的尺。**

---

## ④ 有没有第四形（只量不发明）

| 现成形状 | 位置〔量〕 | 能不能承住 WebView2 的泵要求 | 依据 |
|---|---|---|---|
| 球的 `ui-sta` 泵 | `internal/ball/sta_windows.go:47-96` | **能泵、但和库的嵌套泵共存形状未证** | 它就是那枚已 `CoInitializeEx(STA)`＋`LockOSThread`＋`GetMessageW/Translate/Dispatch` 的常驻泵（`:52-53`、`:62`、`:78-92`）；冲突点＝①-3 第 2 支的再入。 |
| `wmAppTask` 投递链 | `sta_windows.go:127-154`、`ball_windows.go:689` | 甲形的**底座**（投递原语已存在） | 验收腿那句补充成立：`PostTask` 真存在、只是未导出。〔量：见 ③-1〕 |
| message-only 窗口那套 | `ball_windows.go:540 ballWndProc`/`:549 wndProc`；`sta_windows.go:31` 注释逐字 "ball window (created on the thread)" | ⛔ **不是第二枚可复用的泵**：它挂在**同一枚** `ui-sta` 上（`s.hwnd` 就是球窗口，`PostTask` 也投到它） | `staThread` 只有一个 `hwnd` 字段〔量：`sta_windows.go:31`〕；`newSTAThread` 只造一枚〔`:38-44`〕 |
| `cmd/wisp/notify_windows.go` | 待量（本腿尚未读它的 pump 形状） | **尚未判** | 见 ⑦-E |
| 库自带 `webview.Run()` | `webview.go:351-379` | **形态上＝第二枚主泵**（`GetMessageW(…,0,0,0)`＋`IsDialogMessage`＋`WMApp` 分发 `dispatchq`＋`WMQuit` 退出） | 若用它就必须**由它当该线程唯一泵** ⇒ 与球的 `ui-sta` 泵**同线程不可共存**（①-3 第 2 支），不同线程则撞乙形名册钉。〔推〕 |
| 库侧 `webview.Dispatch` | `webview.go:148/152/156`＋`dispatchq`（`:58`） | 只能把活儿投回**库自己的 `mainthread`**（`:108` 在 `NewWithOptions` 里取的） | 谁先调 `NewWithOptions`，谁就是面板窗口的线程 ⇒ **甲形必须让 `ui-sta` 成为那个调用方**。〔推〕 |
| spike 那枚测量腿 | `scripts/spike/webview2-latency/main.go:163-164 bringUp` | 一次性进程，**不是常驻形状** | 票面注释逐字：`panel_host_windows.go:29-31` "The local measurement tests own a locked OS thread … a test harness, not the shipping resident topology, and it registers no goroutine name." |

⇒ **第四形＝有，但不是"新泵"**：④-2 `wmAppTask` 投递链＋④-5 库 `Run()`＋④-6 `Dispatch` 三枚合起来只指向"让 `ui-sta` 当那枚唯一调用方"这一条形状；仓里**不存在**已经现成的"第二枚能独立泵的窗口线程"（除了测试 harness 那枚锁线程，而它不登记名册、不进出货）。〔量＋推〕

---

## ⑤ 我跑了哪些尺、每条真实读数

> ⚠ 本轮为**骨架 commit**（死腿预防：本族已死四次，⑤⑥⑦ 先落盘）。**终态补满在预算闸门**。

- **R1** `date "+%Y-%m-%d %H:%M:%S %z"; git log -1 --format=%H; git status --short | head -40; git branch --show-current`
  ⇒ `2026-10-01 11:21:14 +0800` / `22272ac2b06470b9a03b9acc811eb6d0de944423` / `dev`；脏项含 ` M .gitignore`、大量 `design/**`（D/M）与 `?? .scratch/**`（**本腿一律不动、不提交、不还原**）。
- **R2** `ls -R .scratch/wisp/probes/` ⇒ 顶层 84+ 枚目录名册（`33`、`132`、`228`、`v5-survey` …）。
- **R3** `ls -R .scratch/wisp/probes/33/` ⇒ `a1 h1 r1 r2 r3`；各子目录文件名单见 ②-2。
- **R4** `grep -n "webview2\|jchv" go.mod` ⇒ `19: github.com/jchv/go-webview2 v0.0.0-20260205173254-56598839c808 // indirect`、`20: github.com/jchv/go-winloader v0.0.0-20250406163304-c1995be93bd1 // indirect`。
- **R5** `go env GOMODCACHE` ⇒ `D:\work\base\gopath\pkg\mod`；`ls .../github.com/jchv/` ⇒ 两枚目录名对上版本号。
- **R6** `find <dep> -name "*.go" | sed 's|.*56598839c808/||' | sort` ⇒ 名册 40 枚：`cmd/demo/main.go`、`common.go`、`internal/w32/w32.go(194 行) w32_386.go w32_64bit.go`、`pkg/edge/*`（含 `chromium.go`、`corewebview2.go`、`comproc.go`、`ICoreWebView2*.go`）、`webview.go`、`webviewloader/module*.go`。⛔ **无根级 `chromium.go`**。
- **R7** `grep -n "" <dep>/pkg/edge/chromium.go | head -200` 与 `| sed -n '200,320p'` ⇒ ①-1 表里 `:72-114`、`:116-136`、`:151-157`、`:171-184`、`:186-231`、`:233-248`、`:281-297` 的逐字读数。⛔ 我第一次按派单给的 `chromium.go`（根级）取数 ⇒ `sed: can't read ... No such file or directory`〔量，已在 ①-1 记下〕。
- **R8** `grep -n "" <dep>/webview.go | head -120` 与 `| sed -n '120,420p'` ⇒ `:21-59`（`windowContext`/`webview` 结构＋`dispatchq`）、`:97-128`、`:139-160`、`:162-220`、`:222-261`、`:263-345`、`:347-383`、`:385-420` 逐字。
- **R9** `grep -rn bringUp --include=*.go .` ⇒ `cmd/wisp/panel_host_windows.go:22,25,101,164,174,261`；`cmd/wisp/panel_host_windows_test.go:32,39,49,166-167`；`scripts/spike/webview2-latency/main.go:163,164,200,233,249`。⛔ **库侧零命中** ⇒ `bringUp` 是本仓函数，不是依赖的（见 ⑥-1）。
- **R10** `grep -n "" cmd/wisp/panel_host_windows.go | sed -n '1,120p'` ⇒ 文件头注释 `:5-45` 逐字（含 `:22-31` 的 ui-sta 再入自述与 `:33-45` 的 `w32.Rect` 库缺口叙述）＋ `:47-120` 的 import／`pnl*` proc 名册／`pnlMsg:89-97`／`PanelManager:103-120`。
- **R11** `grep -rn "staThread\|PostTask\|type staThread\|b.sta" internal/ball/*.go` ⇒ ③-1 的未导出名册与全部调用点（含 25+ 处 `b.sta.PostTask`）。
- **R12** `grep -n "^func (b \*Ball)\|^type \|^func " internal/ball/ball_windows.go` ⇒ ③-1 的导出面名册。
- **R13** `grep -n "" internal/ball/sta_windows.go`（全文 162 行）⇒ `:5-9` 契约注释、`:24-36` 结构、`:47-96 start`、`:62 CoInitializeEx(STA)`、`:78-92` 泵、`:127-143 PostTask`、`:146-154 runTask`、`:156-161 quit`。
- **R14** `grep -rn "ui-sta" --include=*.go internal/observe/` ⇒ `goroutine.go:38,44`、`goroutine_test.go:243`；`grep -rn ResidentBaseline --include=*.go .` ⇒ `goroutine.go:40,422`、`goroutine_test.go:268,269`。
- **R15** `grep -rni "panic" .scratch/wisp/probes/33/` ⇒ 15 枚命中**全在 `gate-clauses.txt`**（内容逐字＝`internal/memory/retention.go:100` 与 `internal/plugin/disposal_test.go:87` 两类的源码回声）；`find .scratch/wisp/probes/33/ -type f | sort` ⇒ 84 枚文件全名单。
- **R16** `grep -rn "panic\|⑧" docs/evidence/s1/33-panel-host-c27-r1.md` ⇒ `:7,:28,:33,:41,:45,:46,:53,:66,:73` 逐字（②-4 引用了 `:73` 与 `:46` 的原句）；`ls docs/evidence/s1/ | grep -i 33` ⇒ 10 枚（`33-panel-host-c27-r1.md`、`33-panel-host-c27-v1.md`、`33-minimal-inbound-hop-r1.md`、`33-inbound-listener-r2/r3.md`、`33-inbound-hop-design-a1.md`、`33-35-preflight.md`、`133-*` 四枚另属他票）。

**尚未跑、已排队的尺**（终态前补齐）：`docs/evidence/s1/33-panel-host-c27-v1.md` 全文；`cmd/wisp/notify_windows.go` 的泵形状；`cmd/wisp/panel_host_windows_test.go` 全文（查词面型负向钉）；`internal/ball/*_test.go` 的 `package` 行；`tools/d22scan/main.go` 规则全集；`git cat-file blob HEAD:internal/ball/ball_windows.go` 与 `Events` 字段的消费者。

---

## ⑥ 我可能写错的条目（对抗我自己）

- **⑥-1 派单把 `bringUp` 说成"那枚依赖的泵模型"——它是本仓的函数**。〔量：`grep -rn bringUp` 库侧零命中，命中的是 `cmd/wisp/panel_host_windows.go:174`〕。⇒ 我在 ① 里按"依赖的 `Embed` 泵＋本仓的 `bringUp`"分开写；若我把两者混着叙述，那条叙述是错的。
- **⑥-2 派单给的库路径写作 `.../go-webview2@…/chromium.go`**。〔量：该文件不存在，真身 `pkg/edge/chromium.go`；而 `:96-111`/`:130-131` 两处号数在真身里恰好成立〕⇒ 我照真身取数；若将来有人按派单的字面路径复核会撞 `No such file`。
- **⑥-3 ①-3 的三支全是〔推〕**，我没有任何一发运行证据。特别是**第 1 支（自洽可解、只是阻塞）**：我**没有**逐条读完 `pkg/edge/comproc.go` 去证"COM 回调一定能被同线程的嵌套泵送达"。若库侧有跨线程投递（例如回调 `PostMessage` 到另一枚 message-only 窗口），第 1 支就不成立、第 2 支会变重。**这一条足以推翻 ① 的排序。**
- **⑥-4 我把 `sta_windows.go:134-141` 判成"注释与行为不符"**。〔量：注释逐字 "Window not created yet (or gone): run inline as last resort so callers cannot deadlock on a dead thread."，代码逐字：`s.mu.Lock(); delete(s.tasks, id); s.mu.Unlock(); return`〕⇒ 若我漏看了 `delete` 之前还有一次 `fn()` 调用，这条就写错了。**复核尺：`grep -n "" internal/ball/sta_windows.go | sed -n '126,145p'` 逐字看。** 我认为我读对了，但这句"不符"是我下的判语、不是原文。
- **⑥-5 `ResidentNames` 那枚名册我只读了 `goroutine.go:44` 一行**，把它当成六枚全名册〔量：该行确实含 6 个字符串〕。若名册跨行拼接，我 ③-3 第一族的"六枚"叙述会偏。
- **⑥-6 ③-3 第二族我据"包内测试⇒导出不打红"这条 Go 规则推**，**还没逐枚读 `package` 行**（`hotkey_live_test.go` 等有可能就是 `package ball_test`？名字里带 `_live_test` 的通常是包内）。若其中有 `package ball_test`，则它们对 `b.sta` 的直接访问根本不可能编译 ⇒ 反证它们必是包内。**这条推我不该写成〔量〕，已写成〔推〕。**
- **⑥-7 ④ 表里"message-only 窗口那套＝不是第二枚泵"**：我只量了 `staThread` 结构里只有一枚 `hwnd`。**仓里可能另有 message-only 窗口实现**（例如 `cmd/wisp/notify_windows.go` 或 `internal/.../tray`）。我**尚未证伪**，已在 ⑦-E 挂欠账。若那里真有一枚独立泵线程，④ 的结论要改。
- **⑥-8 凡我写"未导出"的地方，判据都是 `grep` 到小写名字**。**没量过**是否存在同名的**大写转发**（例如 `internal/ball` 里另有 `func (b *Ball) PostTask`，我只在 `ball_windows.go` 上跑了 `^func (b \*Ball)`，**没跑 `internal/ball/*.go` 全目录**〔量：R12 的 glob 只有 `ball_windows.go`〕⇒ ③-1 的"零导出投递面"可能漏计，**这是本格最危险的错。** 补尺见 ⑦-C。
- **⑥-9 我把 `log.Fatalf` 说成"进程没了、不可 recover"**：这是 Go 通识、不是本仓源码读数〔推〕。
- **⑥-10 行数／枚数**：②-2 里"84 枚文件"是 `find | sort` 的行数〔量〕，含 `a2/census.md`？⇒ 不含（该尺跑在我建目录之前）。这一句是提醒：我引用的枚数都带取数时刻。

---

## ⑦ 判不动的地方（逐条 甲／乙／不做 ＋现量＋为什么判不了）

- **⑦-A（最要紧）：再入到底 panic、还是阻塞、还是乱序？**
  现量：`pkg/edge/chromium.go:95-111` 阻塞嵌套泵〔量〕、`:130-136` 无 nil 判定〔量〕、`sta_windows.go:78-92` 同队列主泵〔量〕、`.scratch/wisp/probes/33/` 零该发工件〔量，②〕、`33-panel-host-c27-r1.md:73` 只有散文栈〔量〕。
  为什么判不了：**这一支只有真跑能定案**，而本腿⛔禁跑（`33-r3` 正在 `cmd/wisp` 跑测试）。
  ⇒ 交裁：**甲**＝派一枚"允许跑"的探针腿（在**专用**包／不碰 `cmd/wisp` 写面）复现那再入一次并留原始 dump；**乙**＝接受"未证"就按甲形写（把 ①-3 第 2 支的乱序风险留给用例去接）；**不做**＝继续让生产调用者＝0。
- **⑦-B：`cmd/wisp` 里有没有词面型负向钉会当场打红甲形？**
  现量：`panel_host_windows.go:101` 与 `panel_host_windows_test.go:32` 两处叙述承认"跑在调用方所在线程"；`r1` 表 `:66` 逐字 "⛔ 不许改成独立线程绕过"。
  为什么判不了：我**还没读 `panel_host_windows_test.go` 全文**（它是否含"扫源码词面"的钉、以及 `NewPanelManager` 的调用者钉）。
  ⇒ 甲＝我先读完再报；乙＝派单里预授权"注释＋harness 叙述随甲改"；不做＝按现有叙述不动。
- **⑦-C：③-1 的"零导出投递面"是不是漏计？**
  现量：R12 只扫了 `ball_windows.go` 一枚文件。
  为什么判不了：`internal/ball/` 其余文件（`dock_windows.go`、`liquid_windows.go`、`hotkey*.go`、`tray*.go`）里可能有 `*Ball` 的导出转发。
  ⇒ 甲＝补一枚全目录尺（**我这轮就要补，见交件判语**）；乙／不做＝不该选（这是纯静态、我能量）。
- **⑦-D：`tools/d22scan` 的规则全集会不会扫到甲形的新导出？**
  现量：只从 AGENTS.md 得知 `emojiRe`；票 33 每条腿都跑过该尺（`a1/logs/d22scan.txt` 等）〔量〕。
  为什么判不了：我没读 `tools/d22scan/main.go`；⛔ 且 AGENTS.md 明确"仪器射程≠规格射程"。
  ⇒ 甲＝我读规则名册再点名；乙＝派单里写明"过 d22scan 为判据"；不做＝赌它不扫。
- **⑦-E：仓里有没有第二枚独立泵的窗口线程（message-only / tray / notify）？**
  现量：`cmd/wisp/notify_windows.go` **存在**（派单点名）、`ball_windows.go:927 b.tray.setTip(tip)` 说明有 `tray` 组件。
  为什么判不了：两枚文件的泵形状**我一行都还没读**。若 `notify_windows.go` 自持一枚线程＋泵，它就是**第四形**，①／③／④ 的"最小代价"结论会整体下移。
  ⇒ 甲＝我读完两枚文件（下一轮就做）；乙／不做＝会叫 ④ 变成没答的格子。
- **⑦-F：甲形要不要动 `Events`（`ball_windows.go:49`）的回调契约？动了算不算改 `C*/D*` 契约？**
  现量：`Events` 是导出结构；`ball_windows.go:166-174` 显示 `New` 内部起 `ui-sta` 并经 `observe.Registry.Spawn("ui-sta", …)` 登记。
  为什么判不了：`Events` 的字段名册与"它是冻结契约吗"**我还没量**（要查 `docs/PLAN.md` 的 C/D 表——**只读**、可查，但我这轮没查）。
  ⇒ 甲＝我查 PLAN.md 的 C 表确认 `Events` 是否在冻结面；乙＝编排者直接裁"不算契约"；不做＝风险留给验收腿。
- **⑦-G：`r1` 表 `:33-45` 那句"`PutBounds` 参数是内部 `w32.Rect` ⇒ 无法尺寸化 ⇒ 必须换依赖"是不是真否证？**
  现量：`internal/w32` 是**依赖内部包**〔量：R6 名册路径 `internal/w32/w32.go`〕；Go 规则⇒外部不可 import〔通识〕。
  为什么判不了：我**没读 `pkg/edge` 是否另有导出的尺寸化入口**（`Resize` 是谁调、`ICoreWebView2Controller` 的 `PutBounds` 包装是否导出）。这一条直接决定**丙形是否必需**，而它在我这发之前已被记为"待验"（`r1` 表把它当停手上报）。
  ⇒ 甲＝我读 `pkg/edge/ICoreWebView2Controller.go` 与 `common.go` 全文；乙＝编排者判"丙暂不必"；不做＝把面板留成不可尺寸化。

---

## 交件判语（骨架版）

射程：①（库侧链路与三种最坏形状，含"只有真跑才能定"那一半已明留欠账）／②（零工件，具名 6 枚搜索路径）／③（导出面名册、四枚候选形状、三族会被叫红的用例与行为型／词面型定性）／④（7 枚现成形状逐枚判）／⑤（R1–R16 真实读数）／⑥（10 条自攻）／⑦（A–G 七条 甲／乙／不做＋现量）已落盘。

**没答的格子（不许读成答了）**：⑦-A 的定性（panic／阻塞／乱序）＝本腿**结构推理**，非证；⑦-C/⑦-D/⑦-E/⑦-F/⑦-G **未量完**；③-3 第三族 3-2（词面型钉）尚未定位。**这些由我后续轮次补齐或明写"没答"。**
