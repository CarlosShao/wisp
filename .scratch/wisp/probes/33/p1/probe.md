# 票 33 · 腿 `33-p1` · 探针（**允许真跑**）：面板宿主再入 ui-sta 是哪一支 ＋ "Go→页面"那一跳有没有投递者

本腿只交两问的**读数**，⛔ 不修任何产码、⛔ 不做甲／乙／丙选型、⛔ 不答别的格。
〔量〕＝我这发真跑／真读出来的逐字读数。〔推〕＝我从源码结构推的、**没有读数支撑**。
⚠ 本文件第一枚提交是**骨架**（先写满 ⑤⑥⑦ 三节才动两问正文），后续提交只增不改写已有读数行。

## 起手锚点（同发取数）

- `date` ⇒ `Thu Oct  1 12:23:46 CST 2026`〔量〕
- `git log -1 --format=%H` ⇒ `27bf1fb3f931db7430676987850c46e201885ddb`〔量〕
- `git branch --show-current` ⇒ `dev`〔量〕
- `git status --porcelain` 起手脏项名册 ⇒ **227 行**，逐字落盘 `.scratch/wisp/probes/33/p1/status-start.txt`〔量〕
  （其中 `cmd`／`internal` 两枚包 `git status --porcelain -- cmd internal` 起读数见 §⑥ R2。）
- 代号复核：`ls .scratch/wisp/probes/33/` ⇒ `a1 a2 h1 r1 r2 r3 r4`，**`p1` 未被占用**（起手 `ls .scratch/wisp/probes/33/p1` ⇒ `No such file or directory`）〔量〕
- 本腿写面：`.scratch/wisp/probes/33/p1/**`（本文件＋探针源码＋日志）。⛔ 产码零写面。

---

## ⑤ 我可能写错的条目（对抗我自己）

- **⑤-1 载体若被我写成"产码读数"就是越界**。我的探针线程是**〔复刻〕**：自己 `runtime.LockOSThread()` ＋ `CoInitializeEx(COINIT_APARTMENTTHREADED)` ＋ 自注册窗口类 ＋ `GetMessageW/TranslateMessage/DispatchMessageW` 主泵 ＋ 一枚 `wmAppTask` 任务投递口，逐处照 `internal/ball/sta_windows.go` 的形状抄。
  ⇒ 它证的是**"同一形状下的物理行为"**，**不是**"球的 `OnPanelHotkey` 回调真的会这样"。真常驻那条路上还有 `observe.Registry` 的 `recover`＋`debug.Stack()`＋sink，这一半我**另外单独证**（把复刻线程放进 `Registry.Spawn` 里跑，现成仪器、一字不改）。〔推→将由读数升格〕
- **⑤-2 派单给的依赖路径不存在**。派单写 `pkg/edge/chromium.go:95-111`／`:130-136`——`pkg/edge/chromium.go` **是**真身（`33-a2` 已纠正过派单里的根级 `chromium.go`），行号我本轮自己 `grep -n` 复认（见 §⑥）。凡我引用的行号**一律现取**，不引派单/前人的。
- **⑤-3 `firstRoundTripLocked` 绿 ≠ 回执到了页面**。这一点派单已写死，我复认其形状：`cmd/wisp/panel_host_windows.go:361-368` 的 `done` 在 Go 侧绑定体内关闭，早于 `webview.go:148/152/156` 那三枚 `Dispatch` ⇒ 那枚台件断的是"页→Go 到达"，**一字节都不指向"Go→页到达"**。⚠ 我的探针必须**不靠**这枚绿来下结论，判据只问 JS 那侧拿到 `await` 的值没有。
- **⑤-4 我把 `dispatchq` 的"唯一读者"当事实**——它来自 `grep -rn dispatchq <dep>` 的命中枚数（4 行）〔量，见 §⑥〕，但**"取出＝兑现回执"这一步是推断**；只有真跑知道。若两发读数都收不到回执，先怀疑我的 JS 写法（`window._rpc` 的 stub 形状见 `webview.go:462-479`），不要先怀疑库。
- **⑤-5 桌面真窗读数受负载影响**。`33-r4` 已经量到同一族里"Destroy 后 2s 退净"那一支负载敏感（整包 6 发红 1、`-count=5` 隔离全绿）。我的探针凡涉及"出不去"的判定，**必须用可复现的判据**（线程还在嵌套泵里＝任务计数器还在前进而调用不返回），⛔ 不许只拿"过了 N 秒"当"阻塞"。
- **⑤-6 我只跑自己的探针，不跑 `cmd/wisp` 整包**（今天那两枚红不是我造成的、也不许我改）。凡是"整包现在红几枚"这类话，我只引派单/票面已登记的在册红名册，⛔ 不自行加一枚。
- **⑤-7 `internal/observe` 六枚名册与基线我一字未动**，也⛔没跑它的任何测试。探针里 `Registry.Spawn` 起的线程名**故意不是** `ui-sta`（那枚名册归真常驻）；这会走 `goroutine.go:270-273` 的 `CategoryUnknown ⇒ slog.Warn` 那一支＝**只吵不红**，且只在探针进程里。这一处如果被认为"污染了名册语义"，请指出来，我改判语不改代码。
- **⑤-8 我没有读 `frontend/**` 与 `design/**`**（两层禁令）。凡"页面能不能收到回执"的结论，凭据只有 Go 侧绑定被调次数＋JS 自己报回来的字符串，⛔ 不引前端内容。
- **⑤-9 d22scan 与 `go build ./...` 的 rc 我是**终态**跑的，起手值只作对照。若中途任何一发 rc 不为 0，我会带着那发原文落盘，不擦掉。
- **⑤-10 我自己造过一枚假读数，并且差点把它当成答案**（本文件最要紧的一条自攻）。骨架之后第一发 `-baseline` 我**当场 panic**，看着就像"`33-r1` 那句 panic 被复现了、而且发生在没有再入的情形"。真因＝**我把 `COINIT_APARTMENTTHREADED` 写成 `0x0000`**，而本仓逐字是 `internal/ball/win32_windows.go:149 const coinitApartmentThreaded = 0x2`（`0x0` 是 `COINIT_MULTITHREADED`）。⇒ 那发的 panic 是**我的探针的 bug 的产物**，不是常驻形状的产物。处置＝**不删那发日志**（`logs/reentry-baseline-1.txt`），把它改成 `-com` 三档受控实验（`sta`／`mta`／`none`，R31/R32/R33），于是它变成一枚**真读数**：**MTA 初始化过的线程上创建 WebView2 ⇒ 库在 `chromium.go:175` panic**（见 §B-4）。⇒ 教训＝**先跑对照再主张**，我差一点把对照跑漏。
- **⑤-11 派单第 4 条那句"独立模块**导不进**本仓 `internal/...`（internal 可见性按模块路径算）"——我按它的要求自己验了一次，**验出来的不是这个原因**（R23）。路径前缀套 `github.com/CarlosShao/wisp/scripts/...` 的独立模块**可以**导 `internal/`；`example.com/...` 那种才被 `use of internal package ... not allowed` 拦。今天 `scripts/spike` 用不了的**真**原因是它的 `go.mod` 没 `require`/`replace` 根模块，补上之后撞的是 `missing go.sum entry ...` ⇒ **代价是"要动 `scripts/spike/go.mod`＋`go.sum`＋拉全量依赖"，不是"物理上导不进"**。⛔ 我没有据此改载体（`.scratch` 那一形已经跑通），但**这句话如果进了别人的题面会把路线判错**，所以写在正文里而不只是这里。
- **⑤-12 `tasks_taken_by_nested_pump=5` 在 `-baseline` 里也是 5**（R31/R33）。⇒ **这枚枚数本身不证明"再入把队列洗了"**；证明它的只有 R36 那枚正控（同形、同投递、但不阻塞 ⇒ 0 内／5 后、外层泵 iter 2..6 照跑）＋ §B 里"外层 iter 从进到出冻在 1"那一条。⚠ 如果我只有 baseline 一发就写"嵌套泵吞队列"，那是把**没有第二枚泵时的正常行为**读成缺陷。
- **⑤-13 派单那句"如果走真常驻那条路，栈会自己落盘"——我这发把它测窄了**。`observe` 确实把 `debug.Stack()` **装进** `PanicEvent`（`goroutine.go:302`），但落不落盘取决于有没有人 `SetPanicSink`；现量＝**出货代码里零枚 sink 安装者**（R40，只有 `goroutine_test.go:117-118` 和我的探针）。⇒ 今天真常驻炸的时候，`sink == nil` 那一支走空（`:314`），栈**只在内存里生成、随后丢掉**，能留下的只有 `handle.Err()` 那句 `"panic in goroutine %q (owner %s): %s"`（`:317-320`，**含 recover 值不含栈**）。这一条直接改 `33-r5` 的取数方式：它要么自己装 sink，要么别指望盘上有栈。
- **⑤-14 复刻的两处具名偏差**（⛔ 不许读成"和 `ui-sta` 一模一样"）：① 我建的是 **message-only 窗**（`HWND_MESSAGE`），球是 1×1 的 `WS_POPUP` 顶层窗（`ball_windows.go:196-214`）——同一条线程队列、同一个泵、同一条 `WndProc` 派发路，但**没有桌面表面**；② 我没有 D2D／渲染／托盘／热键那几路消息源，所以"嵌套泵会插走**哪些**消息"我只证到 `wmAppTask` 这一类。⇒ §B 的乱序结论射程＝**该线程队列上的任务投递顺序**，不覆盖"球的出帧会不会被插走"（那一格交 `33-r5` 的球侧用例，见 ⑦-2）。
- **⑤-15 `QUIT=not-sent (host finished on its own)` 是我仪器的正常支**，不是"没测到"。三发 receipt 全走这一支（R25/R26/R27）⇒ 判定不依赖 watchdog。反过来说：如果哪发出现 `QUIT=sent-after-Ns`，那发的**回执读数就作废**，因为泵是被我强制松开的。本腿没有这种发。
- **⑤-16 `WMAPP_DEQUEUED=4` 而不是 3**（R25）。我的解释＝第 4 枚来自页面最后那枚 `wispReport` 自己的回话闭包（3 枚 `wispEcho` ＋ 1 枚 `wispReport`），旁证是 `-evalcheck` 那发**零枚 `wispEcho`、恰好 1 枚**（R27）。⚠ 这是〔推＋两发自洽〕，不是逐条追到的；若有人复算出别的来源，本腿 §A 结论不受影响（**0 枚闭包被取出**这件事与枚数无关）。
- **⑤-17 我只跑了 `-stall 15`／`-deadline 20~25` 这几档**。⇒ "阻塞"那一支在我的读数里是**被排除**的（自然形都在 0.5s 级返回），而不是"我等了很久它没返回"。如果 `33-r5` 在真球上量到长阻塞，那不是我这两发否证过的形状（机器负载、D2D、托盘消息都可能参与）。
- **⑤-18 我全程没有跑 `cmd/wisp` 的任何用例**，所以派单登记的两枚在册红（`TestAC4FocusReturnToPriorWindowGap33r2` 故意的红、`TestPanelHostRealWindowHopAndLifecycle` 负载敏感那一支）**与我无关，我也没复量**。桌面独占期我只开过探针自己的窗，一次一枚，`msedgewebview2` 枚数 14→14（R30）。

---

## ⑥ 我跑了哪些尺、每条真实读数（骨架版先登已跑完的）

| # | 命令（完整） | 逐字读数（摘要） |
|---|---|---|
| R1 | `date; git log -1 --format=%H; git status --porcelain; git branch --show-current` | 起手锚点四件（见上）；`git status --porcelain` 重定向到 `status-start.txt` ＝ **227 行** |
| R2 | `git status --porcelain -- cmd internal \| wc -l` | **0**（起手即净：`cmd`／`internal` 两枚包零脏项 ⇒ 本腿任何写面都不该出现在这两棵树，终态复算同尺） |
| R3 | `ls .scratch/wisp/probes/33/`；`ls .scratch/wisp/probes/33/p1` | 名册 `a1 a2 h1 r1 r2 r3 r4`；`p1` 起手不存在 |
| R4 | `head -5 go.mod`；`cat scripts/spike/go.mod` | 根模块＝`module github.com/CarlosShao/wisp`（`go 1.27`，`toolchain go1.27.1`）；spike 模块＝**`module github.com/CarlosShao/wisp/scripts/spike`**（独立 `go.mod`，自带 `go.sum`，require 含 `github.com/jchv/go-webview2 v0.0.0-20260205173254-56598839c808`） |
| R5 | `grep -n webview2 go.mod go.sum` | `go.mod:19: github.com/jchv/go-webview2 v0.0.0-...-56598839c808 // indirect`；`go.sum:9/:10` 两行 hash 在仓 ⇒ **根模块已能解析这枚依赖**（`// indirect` 是 `tidy -diff` 那格的事，不归本腿） |
| R6 | `find .scratch -name "*.go"` | 早在本腿之前，`.scratch/**` 里就**已有一大片 `.go`**（`147/`、`149/`、`151/`、`152/`、`153/`… 均前人台件）⇒ 载体落这棵树不是我先例 |
| R7 | `GOFLAGS= go list ./... \| wc -l`；`GOFLAGS= go list ./... \| grep -ci scratch`；`GOFLAGS= go list ./... \| grep -c spike` | **35** 枚包；含 `scratch` 的枚数＝**0**；含 `spike` 的枚数＝**0** ⇒ `.scratch/**` **根本不在 `./...` 的匹配图里**（Go 的 `...` 匹配跳过以 `.`／`_` 开头的目录），`scripts/spike/**` 因是独立模块也不在。**这是选载体的决定性尺**：探针落 `.scratch/wisp/probes/33/p1/probe/` 天生不进默认构建图，不需要 `//go:build` 手工 tag 去挡 |
| R8 | `grep -n "" <dep>/webview.go \| sed -n '49,59p'` | `:49-59 type webview{hwnd,mainthread,browser,autofocus,maxsz,minsz,m sync.Mutex,bindings map[string]interface{},dispatchq []func()}` ⇒ `dispatchq` **未导出** |
| R9 | `grep -n "" <dep>/webview.go \| sed -n '443,448p'` | `Dispatch(f)`：`:445 w.dispatchq = append(w.dispatchq, f)`＋`:447 User32PostThreadMessageW.Call(w.mainthread, w32.WMApp, 0, 0)` ⇒ **只投不取**〔量，派单那句成立〕 |
| R10 | `grep -n "" <dep>/webview.go \| sed -n '351,379p'` | `Run()`：`:354 GetMessageW(&msg,0,0,0)` → `:360 if msg.Message == w32.WMApp` → `:362-363` 取空 `dispatchq` 并逐个 `v()`（**全依赖唯一读者**）→ `:368 WMQuit` 才 `return` → `:371-377 GetAncestor/IsDialogMessage/Translate/Dispatch` |
| R11 | `grep -rn dispatchq <dep> --include=*.go` | 命中枚数＝**4**：`:58` 字段、`:362` 读、`:363` 清零、`:445` 写（与派单/`33-a2` 同读数） |
| R12 | `grep -n "" <dep>/pkg/edge/chromium.go \| sed -n '72,114p'` | `Embed`：`:87 createCoreWebView2EnvironmentWithOptions(..., e.envCompleted)`；`:95-111` 那枚循环＝`for { if atomic.LoadUintptr(&e.inited)!=0 {break}; r:=GetMessageW(&msg,0,0,0); if r==0 {break}; Translate; Dispatch }` ⇒ **阻塞式嵌套泵**，`hwnd=0` ⇒ 取该线程**全部窗口＋线程消息**；`:112 e.Init(...)` 在出口后**无条件**执行 |
| R13 | `grep -n "" <dep>/pkg/edge/chromium.go \| sed -n '130,136p'` | `Init` 全段**零 `if`**：`e.webview.vtbl.AddScriptToExecuteOnDocumentCreated.Call(...)`；`e.webview` 只在 `:194-197`（`CreateCoreWebView2ControllerCompleted`）赋值 ⇒ `:106 r==0` 那条 break 之后 `e.webview` 仍为 nil ⇒ 〔推〕nil 解引用＝Go panic。**这一支会不会被走到，只有真跑知道** |
| R14 | `grep -n "" <dep>/pkg/edge/chromium.go \| sed -n '186,231p'` | `:224 atomic.StoreUintptr(&e.inited, 1)` 只在这枚 COM 回调里置位 ⇒ `Embed` 必须自己把回调泵下来〔量〕 |
| R15 | `grep -n "" <dep>/webview.go \| sed -n '139,160p'` | `msgcb`：`:147 w.callbinding(d)`（同步反射 `:191 v.Call(args)`）⇒ 三枚分支 `:148/:152/:156` 全部把回执交给 `w.Dispatch(func(){ w.Eval("window._rpc[id].resolve/reject(...)") })` ⇒ **"Go→页"那一跳在绑定回话这条路上非过 `Dispatch` 不可**〔量〕 |
| R16 | `grep -n "" <dep>/webview.go \| sed -n '450,482p'` | `Bind(name,f)`：`:458-460` 塞 `w.bindings[name]`，`:462-479` `w.Init(...)` 注入 JS stub——`window[name]` 返回一枚 `Promise`，`RPC[seq]={resolve,reject}` 挂在 `window._rpc` 上；`window.external.invoke(JSON.stringify({id,method,params}))` 发出 ⇒ **JS 侧唯一的"取回执"入口就是 `window._rpc[seq].resolve`**（我的探针 JS 就 `await window.<binding>()` 本体，判据不问别的） |
| R17 | `grep -n "" cmd/wisp/panel_host_windows.go \| sed -n '344,412p'` | 现形宿主：`firstRoundTripLocked` `:349-394`（`:361-368` 探针绑定里 `close(done)`＝**Go 侧自证**；`:374-375` `SetHtml` 探测页；`:377-393` 循环 `pnlPumpOnce()`＋select `done`＋`deadline := t0.Add(5*time.Second)`）；`pnlPumpOnce` `:396-411`＝`PeekMessageW(&m,0,0,0,PM_REMOVE)`＋`TranslateMessage`＋`DispatchMessageW`，**结构体 `pnlMsg`（`:89-97`）字段被声明但逐字未解码** ⇒ 今天那形**看不见 `msg.message`**（我的探针泵会解码，因为它要数的就是 `WM_APP` 到达枚数） |
| R18 | `grep -n ".Run()\|Terminate()\|Dispatch(" cmd/wisp/panel_host_windows.go` | **零命中**〔量，复认 `33-a2` R17〕⇒ 出货宿主既不跑主泵、也不调 `Dispatch`/`Terminate` |
| R19 | `grep -n "" internal/observe/goroutine.go \| sed -n '262,330p'` | `Spawn:262`（`:269-273` `CategoryUnknown ⇒ slog.Warn`）；`run:286`；`:294-326` `defer` 内 `:296 recover()`、`:302 Stack: string(debug.Stack())`、`:309 panicCount++`、`:314-316 sink(ev)`、`:323-325 root.Cancel()`〔量，复认派单那句"栈会自己落盘"的仪器形状〕 |
| R20 | `grep -n "func (r \*Registry) SetPanicSink\|type PanicEvent" internal/observe/goroutine.go` | `:153 type PanicEvent`、`:241 func (r *Registry) SetPanicSink(s PanicSink)` ⇒ 现成可用，⛔ 本腿不改一字 |
| R21 | `grep -n "" tools/d22scan/main.go \| sed -n '650,680p'`；`sed -n '241,246p'` | Go 射程＝**只两枚目录**：`walkGo(root/internal)`、`walkGo(root/cmd)`；`:660` 目录剪枝＝`testdata`／`.git`／`s.ign.skip`；`:665` 只扫 `.go` 且**跳 `_test.go`** ⇒ `.scratch/**` **根本不在 d22scan 的 Go 射程**（`33-a2` ⑦-D 那句"testdata 被跳过"我复认，但对本腿来说更强的一读是：**我的载体连被扫的资格都没有**） |
| R22 | `cat scripts/d22scan.sh` | 两步步：`runtests.sh -C tools/d22scan ./...`（正控）＋`go run . -root <repo root>`；无 `|| true`〔量〕⇒ 我这发要证的 rc=0 是**这两步合起来**的 rc |

> ⚠ 本表在骨架版里已登 22 枚**跑完的**尺；探针真跑产生的读数（编译／`go run`／两问各自的日志）会续编号进同一张表，⛔ 不回头改已有行的读数。

**终态续登（R23–R46，全部真跑）**。所有读数的时刻与 HEAD 见 R45/R46；日志原文在 `.scratch/wisp/probes/33/p1/logs/`。

| # | 命令（完整） | 逐字读数（摘要） |
|---|---|---|
| R23 | 载体实测：仓外临时**独立模块** `v1`＝`module github.com/CarlosShao/wisp/scripts/spikeliketmp` ＋ `require`＋`replace` 指回本仓，import `internal/buildinfo`；`v2`＝`module example.com/p1spiketmp` 同法 | `v1` ⇒ **rc=0，真建出 `spikeliketmp.exe`**；`v2` ⇒ **rc=1**，逐字 `main.go:6:2: use of internal package github.com/CarlosShao/wisp/internal/buildinfo not allowed`〔量〕⇒ **"独立模块导不进 internal"这句被我的读数推翻**：可见性按**路径前缀**算，`…/wisp/scripts/...` 在根里 ⇒ **允许**；今天 `scripts/spike` 导不进的**真原因**是它的 `go.mod` 里没有 `require`/`replace`（我第一发补上去之后撞的是 `missing go.sum entry for module providing package golang.org/x/sys/windows`／`go-toml/v2`）。⇒ 见 §⑤-11 |
| R24 | `GOFLAGS= go build -o build/33p1-receipt.exe ./.scratch/wisp/probes/33/p1/receipt` | 通过（第一发报三枚真错：`invalid operation: cannot indirect c.roundsWant`／`undefined: deadlineS`×2，修完即过）〔量〕⇒ 显式目录路径 `go build`／`go run` 能编这枚探针，**且 `go list ./...` 里仍然没有它**（R7） |
| R25 | `GOFLAGS= ./build/33p1-receipt.exe -selfpump -deadline 20`（12:35:32→12:35:42，`logs/receipt-selfpump-1.txt`） | 逐字：`FORM=selfpump ECHO_CALLS=3 ROUNDS_WANT=3 WMAPP_DEQUEUED=4 WMAPP_WITH_HWND=0 OTHER_DEQUEUED=27 REPORTS=1 REPORT_FIRST="0=NOT_RESOLVED;1=NOT_RESOLVED;2=NOT_RESOLVED" REPORT_LATENCY_MS=6150 QUIT=not-sent (host finished on its own) PANICS=0` ＋ `VERDICT selfpump: receipt DID NOT reach the page (awaited value never resolved in 1 reports)` |
| R26 | `GOFLAGS= ./build/33p1-receipt.exe -run -deadline 20`（12:35:52→12:35:55，`logs/receipt-run-1.txt`） | 逐字：`FORM=run RESULT=RUN-RETURNED (WM_QUIT consumed the library loop)`；`FORM=run ECHO_CALLS=3 ROUNDS_WANT=3 WMAPP_DEQUEUED=0 WMAPP_WITH_HWND=0 OTHER_DEQUEUED=3 REPORTS=1 REPORT_FIRST="0=echo:ping-0;1=echo:ping-1;2=echo:ping-2" REPORT_LATENCY_MS=107 PANICS=0` ＋ `VERDICT run: receipt REACHED the page in 1 of 1 reports` |
| R27 | `GOFLAGS= ./build/33p1-receipt.exe -selfpump -evalcheck -deadline 25`（12:48:38→12:48:43） | 逐字：`FORM=selfpump EVAL_ISSUED at=800ms (direct w.Eval on the UI thread, bypassing Dispatch/dispatchq)`；`FORM=selfpump ECHO_CALLS=0 ... WMAPP_DEQUEUED=1 WMAPP_WITH_HWND=0 OTHER_DEQUEUED=27 REPORTS=1 REPORT_FIRST="EVAL_SEEN title=EVAL-OK-33P1 polls=4" REPORT_LATENCY_MS=907 PANICS=0` ＋ `VERDICT selfpump+evalcheck: a direct Go-side Eval DID reach the page ... even though the binding reply did not -> push and reply are two different hops` |
| R28 | `grep -n "" <dep>/common.go \| sed -n '26,50p'` | `WebView` 接口名册（`:26-50`）＝`Run()`、`Terminate()`、`Dispatch(f func())`、`Destroy()`、`Window()`、`SetTitle()`… ⇒ **零枚可导出 drain／取出 `dispatchq` 的方法**〔量〕 |
| R29 | `grep -rn "func.*[Dd]rain\|dispatchq" <dep> --include=*.go`（去测试） | 命中只有 `webview.go:58`（字段）、`:362`、`:363`（`Run()` 内）、`:445`（`Dispatch` 内）⇒ **仓外无法取出，库内无第二枚读者**〔量，复认 R11〕 |
| R30 | `GOFLAGS= go build -o build/33p1-reentry.exe ./.scratch/wisp/probes/33/p1/reentry`；`tasklist //FI "IMAGENAME eq msedgewebview2.exe" //FO CSV \| grep -ci msedgewebview2` | 编译通过；枚数 **BEFORE=14 → AFTER=14**（探针全程跑完后同尺复量＝**没往这台机器上留孤儿 webview 子进程**）〔量〕 |
| R31 | `./build/33p1-reentry.exe -baseline -com sta -stall 15`（12:44:06，`logs/reentry-baseline-sta.txt`） | `STA tid=16732 locked=true com=sta`；`CoInitializeEx(0x2 COINIT_APARTMENTTHREADED) r=0`；`FINAL mode=baseline create_returned=true create_ms=561 tasks_taken_by_nested_pump=5 tasks_after_callback=0 outer_iters=2 panics_total=0 sink_writes=0`；`VERDICT baseline: bring-up with this thread as the ONLY pump returned in 561ms` |
| R32 | `./build/33p1-reentry.exe -baseline -com mta -stall 15`（12:44:10） | `CoInitializeEx(0x0 COINIT_MULTITHREADED) r=0`；`SINK-EVENT goroutine=probe33p1-usta recovered="runtime error: invalid memory address or nil pointer dereference"`；`STA-HANDLE-RETURNED err=internal: panic in goroutine "probe33p1-usta" (owner 33-p1): ...`；`FINAL ... create_returned=false panics_total=1 sink_writes=1 sink_file=C:\Users\swq\AppData\Local\Temp\wisp33p1\reentry-baseline-mta\panic-sink.txt`；**进程没死，rc=0** |
| R33 | `./build/33p1-reentry.exe -baseline -com none -stall 15`（12:44:15） | `CoInitializeEx SKIPPED (com=none)`；`FINAL mode=baseline create_returned=true create_ms=531 tasks_taken_by_nested_pump=5 panics_total=0` |
| R34 | `./build/33p1-reentry.exe -reentrant -com sta -stall 15`（12:44:38→12:44:42） | `FINAL mode=reentrant create_returned=true create_ms=544 tasks_taken_by_nested_pump=5 tasks_after_callback=0 outer_iters=3 ping=0 panics_total=0 sink_writes=0`；`SEQ CALLBACK enter outer_iters=1`／`SEQ task P1..P5 ran INSIDE the callback ... outer_iters=1`／`SEQ bring-up RETURNED a window after 544ms; still inside the same callback frame=true`／`SEQ Bind from inside the callback OK`／`SEQ SetHtml from inside the callback issued`／`SEQ CALLBACK exit outer_iters=1`；`VERDICT reentrant: branch (c) ORDER-REVERSAL ...` |
| R35 | `./build/33p1-reentry.exe -reentrant -com none -stall 15`（12:44:42→12:44:46） | 同形读数：`create_returned=true create_ms=563 tasks_taken_by_nested_pump=5 tasks_after_callback=0 panics_total=0 sink_writes=0`；`CALLBACK enter outer_iters=1` → `CALLBACK exit outer_iters=1`（**外层泵一整段零前进**） |
| R36 | `./build/33p1-reentry.exe -skipcreate -com sta -stall 15`（12:45:58）**正控** | `FINAL mode=skipcreate-control create_ms=0 tasks_taken_by_nested_pump=0 tasks_after_callback=5 outer_iters=6`；`SEQ task P1 ran after the callback, in order (outer_iter=2)` … `P5 ... (outer_iter=6)`；`VERDICT skipcreate-control: ORDER KEPT - all 5 tasks ran after the callback, the nested pump took none.` ⇒ **枚数仪器不是恒 0／恒 5 的瞎尺**，它分得开两种形 |
| R37 | `./build/33p1-reentry.exe -quitduringcreate -com sta -stall 15`（12:45:58→12:46:04） | `SEQ FORCED WM_QUIT posted to tid=16144 r=1`；`SINK-EVENT ... recovered="runtime error: invalid memory address or nil pointer dereference"`；`FINAL mode=quitduringcreate create_returned=false tasks_taken_by_nested_pump=5 outer_iters=1 panics_total=1 sink_writes=1`；`VERDICT quitduringcreate: branch (a) PANIC reachable, and the observe sink captured debug.Stack()`；**进程存活 rc=0** |
| R38 | `cat %TEMP%\wisp33p1\reentry-quitduringcreate-sta\panic-sink.txt` | 栈顶逐帧（现取）：`runtime/debug.Stack()` → `internal/observe.(*Registry).run.func1() .../goroutine.go:302` → `panic(...)` → **`edge.(*Chromium).Init(...) .../pkg/edge/chromium.go:131`** → `edge.(*Chromium).Embed(...) chromium.go:112` → `(*webview).CreateWithOptions(...) webview.go:340` → `NewWithOptions(...) webview.go:109` → `main.probeCallback.func1()` → `main.(*thread).runTask()` → `main.probeWndProc(...)` ⇒ **`33-r1` 那句散文栈的帧名在这里逐字出现了**（但只在人为造 `WM_QUIT` 那发） |
| R39 | `cat %TEMP%\wisp33p1\reentry-baseline-mta\panic-sink.txt` | 栈顶：`...run.func1() goroutine.go:302` → panic → **`edge.(*Chromium).EnvironmentCompleted(...) pkg/edge/chromium.go:175`** → `_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerInvoke(...) corewebview2.go:266` → `webviewloader.CreateCoreWebView2EnvironmentWithOptions(...) module.go:122` → `edge.(*Chromium).Embed(...) chromium.go:87` → `CreateWithOptions webview.go:340` → `NewWithOptions webview.go:109` ⇒ **`Embed` 的第 87 行里回调被同步调用**（不是我预期的"回调靠泵送下来"），`:171 if int64(res) < 0` 没拦住 ⇒ `:175 env.vtbl.AddRef` 解引用 nil `env`〔量＋推：res 那枚 32 位 HRESULT 当 uintptr 传进来不会符号扩展〕 |
| R40 | `grep -rn "SetPanicSink" --include=*.go .`（去 `goroutine.go` 自身） | 命中＝**2 枚**：`internal/observe/goroutine_test.go:117/:118` ＋ 我这枚探针 `reentry/main.go:220` ⇒ **出货进程里今天没人装 sink**；`goroutine.go:314 if sink != nil` 那一支走空 ⇒ 见 §⑤-13（对派单那句"栈会自己落盘"的更正） |
| R41 | `sh scripts/d22scan.sh`（终态，12:49:32，`logs/d22scan-final.txt`） | **rc=0**；逐字尾行 `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=224, bans #1-5 cmd/=33, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=476, ban #8 cmd/=78`；名册里**零枚 `.scratch` 射程**⇒ 我的探针根本不在扫面内 |
| R42 | `GOFLAGS= go build ./...`（终态，12:48:58） | **rc=0**（无任何输出）；同发 `GOFLAGS= go list ./... \| grep -c scratch`＝**0**（R7 复量） |
| R43 | `git status --porcelain -- cmd internal \| wc -l`（终态复量） | **0**＝起手值（R2）。⚠ 探针的 exe 落在 `build/`，那是**起手脏名册里已有**的构建产物目录（`status-start.txt`），⛔ 我没提交它 |
| R44 | `date` 每发前缀 | 关键发：12:35:32／12:35:52／12:44:06／12:44:10／12:44:15／12:44:38／12:44:42／12:45:58／12:45:58／12:48:38／12:48:58／12:49:32（全为 CST＋8） |
| R45 | 本腿所有读数的锚点：`git log -1 --format="%H %ci"` | **`e7eee03e15ba5820ce8588e1d662235db4768017`（2026-10-01 12:29:37 +0800，＝本腿自己那枚骨架提交）**；上面每一条真跑读数都取自这一枚 HEAD（骨架提交之后、本文件终态提交之前，HEAD 未推进）〔量〕 |
| R46 | `git status --porcelain`（终态复算写面名册） | 见 §收尾三把尺；相对起手 227 行的增量只有 `.scratch/wisp/probes/33/p1/**` |
| R47 | `grep -rn "CoInitialize" cmd/wisp/*.go`；`grep -n "" cmd/wisp/panel_host_windows_test.go \| sed -n '38,62p'` | **零命中**＋harness 逐字＝`:44 go func()`／`:46-50 recover → fmt.Errorf("panel host pump panicked: %v", r)`（**这枚形状天生不含栈**，复认 `33-a2` ②-6）／`:51 runtime.LockOSThread()`／`:52 errCh <- mgr.bringUp(ctx)`／`:53-56 for !stopping { pnlPumpOnce(); Sleep(5ms) }` ⇒ **出货宿主与测试 harness 都不初始化 COM apartment**（§B-4 的凭据） |


---

## ⑦ 判不动的地方（逐条 甲／乙／不做 ＋ 现量 ＋ 为什么判不了）

- **⑦-1 真常驻那条路（起球 ＋ 从 `OnPanelHotkey` 里创建面板）我走不走得到？**
  现量：`internal/ball/**` **禁写**（`33-r5`／票 228 要用）；`b.fire`（`ball_windows.go:728-732`）未导出、`staThread`／`b.sta`／`PostTask` 全未导出（`33-a2` R20 全目录尺）⇒ **包外没有任何一条路能把闭包投到 `ui-sta`**，也就没有"从球回调里触发"的编程入口；只剩两条真触发＝全局热键（`SendInput` 造键）或托盘菜单点击。
  为什么判不了：这两条都要往**桌面**注入输入，而我拿不到"热键已注册且就是这个组合"的可复核前提（`RebindHotkeys`/`RegisteredHotkeys` 只报注册结果，且真按下去会把球/托盘/别人的前台窗一起搅进来）。
  ⇒ **甲**＝探针里复刻 `ui-sta` 形状（本腿现在这条），另起一发把复刻线程放进 `observe.Registry.Spawn` 拿真 sink；**乙**＝真起球＋`SendInput` 按热键（桌面副作用大、失败模式与"再入"不可区分）；**不做**＝把这一格留成未证。
  ⚠ **本腿选了甲并具名标〔复刻〕**：它给的是物理行为，不是"球那条回调一定这样"。**请编排者据此定 `33-r5` 题面里"证到哪一层"这句。**
- **⑦-2 `33-r5` 若必须"甲形（投 `ui-sta`）＋出货路径"，读数够不够？**
  现量：`J1` 裁定（票面 `:24`）写死"落地腿的第一枚用例就证泵期间球仍能出帧、消息仍被派发"——那是**球侧用例**，本腿的复刻读数只能告诉它"物理上会不会发生"，⛔ 不能替它勾。
  为什么判不了：判它要 `internal/ball` 的写面（禁）。⇒ 甲＝只交物理读数；乙＝我把球侧用例也造了（越界）；**本腿做甲**。
- **⑦-3 `Run()` 那一形与 `ui-sta` 能否共存**（票面 `:43` 的 ⓐ 那一支）。
  现量：`Run()` `:351-379` 与 `Embed` `:95-111` 都是 `GetMessageW(…,0,0,0)` 全队列泵；`Run` 只在 `WMQuit` 返回。
  为什么判不了：**这属"选型"射程**（派单第 41 行 ⛔ 不做选型）。我只量"回执到／不到"这一维，谁泵＝不答。
- **⑦-4 面板窗口的 `WM_DESTROY → Terminate() → PostQuitMessage`**（`webview.go:242-243`、`:381-383`）会把 `WM_QUIT` 投进宿主线程队列。真常驻里那枚 `WM_QUIT` 会**同时**结束球的 `ui-sta` 主泵（`sta_windows.go:78-92`）。
  现量：两处行号本轮自己 `grep -n` 复认〔量〕。
  为什么判不了：本腿**不测 Destroy 那一跳**（票面 `A501` 已把"退净"那一支归 `33-r5`，且 `33-r4` 已登记它负载敏感）。⚠ 但它与 ⑦-1 有耦合（`Embed` 的 `:106 r==0` break 正是靠这枚 `WM_QUIT` 才走得到 nil 解引用），所以我在 §⑦-1 的读数里**会具名说明我有没有造出过 `WM_QUIT`**，免得下一程把我的 panic 归因错。
- **⑦-5 回执那一跳若两形都收不到**：先怀疑谁？
  现量：`window._rpc` stub（R16）＋ `msgcb`（R15）＋ `dispatchq` 唯一读者（R11）。
  为什么判不了：真跑之前无法排除"我的 JS 写法不对"。⇒ 处置＝**每发都带一枚自证**：Go 侧必须看见 JS 报回来的那枚调用（第二枚绑定被调），否则那发的 JS 判定作废、不算读数。
- **⑦-6 本腿结不了的两桩旧账（具名，不是我加的）**：`go mod tidy -diff` 在 HEAD 上 exit 1（票面 `:267`）⇒ 我⛔没跑；`staticcheck` 本机版与 CI 钉版不同（票面 `:262`、`:296`）⇒ 我⛔没跑。

### ⑦-结 本腿真跑之后，上面哪几格被结掉了、哪几格仍然判不动

**结掉的（都带 R 号）**：
- `33-a2` 的 **⑦-A／①-3 那一岔**＝本腿 §B 主文（读数 R31–R39）：**自然形是第 (c) 支**，不是第 1 支也不是第 3 支。
- `33-a2` 的 **⑦-H／⑥-5 那半"只有 `Run()` 才兑现回执"**＝本腿 §A 主文（R25/R26/R27）：两形对照读数齐全，且**那枚通识（线程消息不经 `DispatchMessageW` 送达窗口过程）现在有本机读数了**（`WMAPP_DEQUEUED=4` 而 `WMAPP_WITH_HWND=0`、页面三枚 promise 全不 resolve）。
- ⑦-4 里"`Embed` 的 `:106 r==0 → :131` nil 解引用可达吗"＝**可达**，但要人为造 `WM_QUIT`（R37/R38）；自然形这一发没走到。

**仍然判不动的（⛔ 不许读成我答了）**：
- **⑦-结-1 真球那条回调到底怎样**。我这发是〔复刻〕（⑤-14 列了两处偏差）。把 `bringUp` 放进**真的** `OnPanelHotkey` 里跑，需要 `internal/ball` 的导出投递面（今天零枚，`33-a2` R20）⇒ 本腿⛔禁写那一棵，做不到。⇒ **这一格请照 §B 的结论按"形状物理"采信，不要当"球一定这样"**。
- **⑦-结-2 球的出帧／渲染会不会被嵌套泵插走**。我只证了 `wmAppTask` 这一类消息的顺序（R34/R36）。D2D 的 `WM_TIMER` 出帧、托盘回调、热键消息**一枚都没测**（复刻里没有那些消息源）。`33-a2` F6 那一族行为型钉（`internal/ball/live_windows_test.go`、`live_guard_windows_test.go:79-80` 数 `wmTimer`）**会不会红，仍然只有球侧用例知道**。
- **⑦-结-3 MTA panic 的 `res` 真值**。R39 那发我只拿到"回调被同步调用、`env` 为 nil、`int64(res) < 0` 没拦住"这三件；**那枚 HRESULT 具体是什么我拿不到**——要拿就得自己实现 `ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler`（依赖里那些 handler 类型全是未导出的 `iCore...Impl`/`Vtbl`），本腿没做。⇒ 任何腿引用"STA/MTA 决定创建成败"时，请带上这句：**成因未定值，只定了现象**。
- **⑦-结-4 回执那一跳在 `Run()` 形里是否与球的泵共存**。`Run()` 与 `ui-sta` 主泵同线程互斥这一条仍是〔推〕（`33-a2` ④ 末行），本腿⛔不做选型、也没测"两枚泵同线程"。（我只测了"单独一线程跑 `Run()` 时回执到"。）
- **⑦-结-5 `dispatchq` 无界增长**我没**直接**量（未导出、拿不到 len）。间接读数是有的：R25 里 4 枚 `WM_APP` 被取出、0 枚闭包被执行 ⇒ 那 4 枚一定还挂在 `:58` 那枚切片上。⛔ "泄漏"这个词我只用到"闭包没被取出"这一层，没用到"内存涨到多少"。


---

## 载体形状（具名 ＋ 为 ②③ 两条跑过的命令）

**用了哪一形**：`.scratch/wisp/probes/33/p1/probe/main.go` ＋ `go run` 显式文件路径。
为什么**不是** `scripts/spike/**`：**两个理由，其中派单给的那个被我的读数否掉了**。
- ⛔ 派单那句"独立模块**导不进**本仓 `internal/...`（internal 可见性按模块路径算）"＝**不成立**（R23 实测：路径前缀套在 `github.com/CarlosShao/wisp/` 里面的独立模块 ＋ `require` ＋ `replace` ⇒ **rc=0 真建出 exe**；只有路径在根外面的才吃 `use of internal package ... not allowed`）。
- ✅ 成立的代价理由：`scripts/spike/go.mod` 今天**没有** `require`/`replace` 指回根模块，硬要用它就得改 `scripts/spike/go.mod`＋`go.sum`（我补上 `require` 后第一发吃的是 `missing go.sum entry for module providing package golang.org/x/sys/windows`／`github.com/pelletier/go-toml/v2`）。那是在**出厂树**里留 diff，且要为探针拉一遍依赖 ⇒ 弃用。理由已在 §⑤-11 具名更正。
为什么**不是** `cmd/wisp/testdata/<name>/main.go`：那仍在 `cmd/wisp/` 这棵树下，本票下一枚验收腿（`33-v2`）要逐名比 `cmd/wisp` 名册；派单第 1 条写死"⛔ 不给 `cmd/wisp` 测试包添文件"，我按**更严**的读法走（连 `cmd/wisp/**` 都不落）。
为什么**不是** `//go:build` 手工 tag 门：R7 实测 `GOFLAGS= go list ./...` ＝ **35 枚包、零枚含 `scratch`**，`.scratch/**` 天生不在 `./...` 的匹配图里 ⇒ **不需要**再加一枚 tag 去挡；我也⛔没测过"全文件被约束排除的包会不会让 `go build ./...` 报 `build constraints exclude all Go files`"（没必要冒这个风险去测）。
为什么**选** `.scratch/wisp/probes/33/p1/<name>/main.go` ＋ 显式目录：① 它在**根模块内** ⇒ 能 `import "github.com/CarlosShao/wisp/internal/observe"`（我的探针真这么写了 panic sink，R37/R38 的栈就是从这枚现成仪器里出来的）；② 它不在构建图里（R7）；③ 它不在 d22scan 的 Go 射程里（R21：Go 只走 `internal/`＋`cmd/`；R41 的 scope 名册里**零枚 `.scratch`**）；④ 编译/运行用显式路径 `go build -o build/... ./.scratch/.../receipt`，不污染 `./...`。

为这一形跑过的命令（两枚闸门都要 rc=0，终态复跑）：
- `GOFLAGS= go build ./...`
- `GOFLAGS= go list ./...`（证明探针**不在**构建图里）
- `sh scripts/d22scan.sh`

---

## §A 问题② · "Go→页面"那一跳有没有投递者（票 33 AC#14／`docs/PLAN.md:1367` C17）

**一句话结论**：**自然形（票 33 现在这一形自泵）＝回执到不了页面；`Run()` 形＝到。** 而"Go 主动 `Eval` 推一段 JS"在自泵形里**也能到**——这两件事是**两维**，不许合成一句（下面各有一发读数）。

### §A-1 三发读数（同一枚探针、同一枚页面协议、只差泵的形状）

页面干的事（`receipt/main.go` 的 `page()`，逐字在源码里）：连续 **3 次** `await window.wispEcho("ping-i")`，每次跟一枚 2 秒 `setTimeout` 赛跑，然后把三次拿到的值拼成字符串，**用第二枚绑定 `wispReport` 报回 Go**。
⇒ 判据**只在页面嘴里**：拿到 `echo:ping-i` ＝ 回执到了；拿到 `NOT_RESOLVED` ＝ 回执没到。**"Go 侧绑定被调过"我只当旁证看**（那正是今天 `firstRoundTripLocked` 会绿的那一枚）。

| 发 | 形状 | `ECHO_CALLS`（页→Go 旁证） | 页面报回来的值 | `WMAPP_DEQUEUED` / `WITH_HWND` | 判定 |
|---|---|---|---|---|---|
| R25 | **自泵**＝`PeekMessageW(0,0,0,PM_REMOVE)`＋`Translate`＋`Dispatch`（逐字复刻 `panel_host_windows.go:400-411`，只多解码两个字段） | **3**（入向通） | `0=NOT_RESOLVED;1=NOT_RESOLVED;2=NOT_RESOLVED`（6150ms） | **4 / 0** | **回执没到页面** |
| R26 | **`Run()`**＝宿主线程进库自己的主泵（`webview.go:351-379`），拿到报告后由 Go 侧 `w.Terminate()` 收尾 | **3** | `0=echo:ping-0;1=echo:ping-1;2=echo:ping-2`（**107ms**） | 0 / 0（取 `WM_APP` 的是库，不是我） | **回执到了页面** |
| R27 | **自泵 ＋ 只走 `Eval`**（Go 在 800ms 时在这条 UI 线程上直接 `w.Eval("document.title='EVAL-OK-33P1'")`，页面每 200ms 轮询自己的 title 并报回） | 0（这发不调 echo） | `EVAL_SEEN title=EVAL-OK-33P1 polls=4`（907ms） | 1 / 0 | **推得到**（`Eval` 不经 `dispatchq`，`webview.go:439-440` 直调 `browser.Eval`） |

### §A-2 这三发把哪几格从〔推〕升成〔量〕

1. **`33-a2` ⑥-5 那枚"最重的自我怀疑"结清**（那句是 Win32 通识、当时盘上零读数）：R25 直接量到 **自泵把 4 枚 `WM_APP` 从队列里取走了，`msg.hwnd` 全为 0（线程消息），`DispatchMessageW` 把它们送到了谁也不在的地方**，而**页面侧 3 枚 promise 一枚都没 resolve**。⇒ 不是"慢"，是**没有投递者**：`Dispatch`（`webview.go:443-448`）只 append＋投线程消息；`dispatchq` 全依赖只有 4 处引用（R29），取出点只有 `Run()`（`:362-363`）；`WebView` 接口名册里**没有任何一枚可导出的 drain**（R28）⇒ **仓外形状无法兑现那一跳，只有 `Run()` 形可以**（R26 是现场证据）。
2. **`firstRoundTripLocked` 的射程就此定死**：它等的 `done` 在 Go 侧绑定体内关闭（`panel_host_windows.go:361-368`），**R25 证明同一枚页面协议里 JS 自己什么也没收到**。⇒ 那枚用例绿＝"页→Go 到达"，**它不构成 AC#14 的任何凭据**（连"回执"这一维的旁证都不是，因为两维在 R25 里同时给出了相反值）。
3. **ⓐ／ⓑ／ⓒ 三形里"ⓑ 只能救推送、救不了绑定回话"这句（票面 `:43`）现在有读数了**：R27（推得到）与 R25（回话到不了）**是同一枚自泵形状里的两发**。写方案时那两件事不许再并成一句。
4. 附带量到：**自泵形里那几枚闭包不会自己消失**——4 枚 `WM_APP` 被取走、0 枚执行 ⇒ 闭包仍挂在 `:58` 那枚切片上。⚠ 我没直接读到长度（未导出），这一句是〔量＋推〕，见 ⑦-结-5。

### §A-3 本格欠 `33-r5` 的，只有形状读数，没有选型

我没有答"该用 ⓐ 还是 ⓒ"。我只交出可复算的三件事：**(i)** 自泵形里那一跳**今天确实不通**，**(ii)** 库里唯一能兑现它的形状是 `Run()`，**(iii)** `Eval` 这条路只能兑现"主动推送"那一半。**判据长什么样我有读数支持**：会响的断言必须**由页面自己报**（我这发的 `wispReport` 那一形），⛔ 断 `Dispatch 被调过`／断 `Eval 没报错`／断 `Go 侧 done 关闭` 三枚都在 R25 里同时成立却仍然**没回执**——那正是票面 `:44` 禁的两种假绿的活样本。

---

## §B 问题① · 在正在泵的线程的窗口回调里再入创建：panic／阻塞／乱序

**一句话结论**：**是第 (c) 支——返回了，但顺序被打乱（嵌套泵把外层队列里的任务先执行掉了）；不 panic、也不永久阻塞**〔复刻，见 ⑤-14〕。自然形的耗时 **544ms**，与"没有外层泵"的对照 **561ms** 同一量级 ⇒ 不是慢；外层泵在整个回调期间 `iter` **冻在 1**（一步没走），而它本来该按序处理的 5 枚任务**全部提前在回调内部被执行**。

### §B-1 五发读数（同一枚复刻线程，只差两个变量：是否再入、COM 怎么初始化）

| 发 | 再入？ | COM 档 | 创建返回 | `nested`（回调内执行的任务枚数） | `after`（回调后按序执行） | 外层泵 iter | panic |
|---|---|---|---|---|---|---|---|
| R31 `-baseline -com sta` | ⛔ 否（直接调） | **STA `0x2`** | **561ms** | 5 | 0 | 0→2（**当时还没有外层泵**） | 无 |
| R33 `-baseline -com none` | ⛔ 否 | 未初始化 | **531ms** | 5 | 0 | 0→2 | 无 |
| R32 `-baseline -com mta` | ⛔ 否 | **MTA `0x0`** | **没返回** | — | — | 0 | **有**，帧＝`chromium.go:175`（见 §B-4） |
| R34 `-reentrant -com sta` | ✅ 是（`WndProc` 内的任务里） | **STA `0x2`** | **544ms** | **5** | 0 | **1 → 1（零前进）** | 无 |
| R35 `-reentrant -com none` | ✅ 是 | 未初始化 | **563ms** | **5** | 0 | **1 → 1** | 无 |
| R36 `-skipcreate -com sta`（**正控**） | ✅ 是，但**不阻塞** | STA | — | **0** | **5**（`outer_iter=2..6`） | 1→6 | 无 |
| R37 `-quitduringcreate -com sta`（**人为造 `WM_QUIT`**） | ✅ 是 | STA | **没返回** | 5 | — | 1 | **有**，帧＝`chromium.go:131`（见 §B-3） |

### §B-2 三支结局的逐条判定（自然形＝R34/R35，对照＝R31/R33/R36）

- **(a) panic**：⛔ **自然形两发都没有**（`panics_total=0`、`sink_writes=0`）。⇒ **`33-r1` 那句"从球的 ui-sta 再入 ⇒ panic"在我的复刻形里不复现**；本仓 `.scratch/wisp/probes/33/` 从今天起多了一枚**反证读数**（R34/R35），而它先前只有散文（`33-a2` ②-5）。
- **(b) 永久阻塞**：⛔ 不是。判据不是"我觉得快"，而是**三件同时成立**：① `create_returned=true`；② 时长与对照同一量级（544 vs 561ms）；③ 返回之后外层泵继续走（`iter 1→3`）并且窗口能被正常 `Bind`／`SetHtml`（`SEQ Bind from inside the callback OK`、`SetHtml ... issued`）。
  我还给"真的出不去"预备了判据并跑过一遍那支持（R37/R32 那两发 `create_returned=false` ⇒ watchdog 会打出 `tasks_taken_by_nested_pump` 与 ping 增量再强制 `WM_QUIT`）；**自然形一次都没触发 watchdog**（三发 receipt 与 R34/R35 都是 `QUIT=not-sent`）。
- **(c) 返回了但消息乱序／回调重入**：✅ **就是这一支**。证据链三条，各自独立：
  1. **外层泵零前进**：`SEQ CALLBACK enter outer_iters=1` … `SEQ CALLBACK exit outer_iters=1`（R34/R35 两发同值）。外层那枚 `DispatchMessageW` 还没返回，线程就又在里面泵了一整轮。
  2. **5/5 提前**：`task P1..P5 ran INSIDE the callback ... outer_iters=1`；枚数 `nested=5 / after=0`。
  3. **正控证明这枚仪器分得开**：R36 同形、同投递、但不阻塞 ⇒ `nested=0 / after=5`，且那 5 枚分别落在 `outer_iter=2,3,4,5,6`。**顺序在"不阻塞"时是被守住的，在"阻塞"时是被插走的** ⇒ 这不是仪器恒 0／恒 5。
  ⚠ **`-baseline` 那两发也是 `nested=5`**（⑤-12）：那时**外层泵还没启动**，嵌套泵是唯一的泵，所以"在回调里执行"是**正常**的，不构成乱序证据。**只有配上第 1 条（外层已 iter=1 且冻住）才算**——这就是我先跑对照的理由。
- 直接后果（对出货形状的说法，⛔ 不是我的选型）：`sta_windows.go:54-55` 与 `ball_windows.go:724-727` 那两处契约自述（"回调必须快、非阻塞；阻塞的回调会停住消息泵；一条协程不为手势起"）在**这一形**里被同时兑现成两件事：**`bringUp` 期间球的泵整段停 0.55 秒**，且**队列里已投递的 `wmAppTask` 会被嵌套泵按"嵌套泵的取数顺序"提前跑掉**。`33-r5` 如果走"投 `ui-sta`"那形（票面 `:24` J1 甲形），它要证的"泵期间球仍能出帧、消息仍被派发"这一格——**我这发给的是坏消息：消息会被派发，但是被插着派发**；而"出帧"那一维我没测（⑦-结-2）。

### §B-3 人为造出来的那支 panic：`33-a2` ①-3 第 3 支**可达**，且栈的形状与散文预测逐字对上

- 造法（R37）：在窗口回调里**先** `PostThreadMessageW(tid, WM_QUIT)` **再**调 `NewWithOptions` ⇒ `Embed` 的 `chromium.go:100` 那枚 `GetMessageW` 立刻拿到 `r==0`（`:106` break），此时 `:97` 的 `inited` 仍是 0、`e.webview` 从未被 `:194-197` 赋值 ⇒ `:112 e.Init(...)` 落到 `:130-136` 那枚**没有任何 `if` 的**解引用。
- 原始栈落盘（**这就是我该拿的第二样东西**）：`C:\Users\swq\AppData\Local\Temp\wisp33p1\reentry-quitduringcreate-sta\panic-sink.txt`，头四帧逐字：`runtime/debug.Stack()` → `internal/observe.(*Registry).run.func1() .../internal/observe/goroutine.go:302` → `panic(...)` → **`github.com/jchv/go-webview2/pkg/edge.(*Chromium).Init(...) .../pkg/edge/chromium.go:131`** → `Embed chromium.go:112` → `CreateWithOptions webview.go:340` → `NewWithOptions webview.go:109` → 我的回调 → `runTask` → `probeWndProc`。
- 进程行为：⛔ **进程没死**（`rc=0`），`observe` 兜住 ⇒ `PANICS=1`、`STA-HANDLE-RETURNED err=internal: panic in goroutine "probe33p1-usta" (owner 33-p1): runtime error: invalid memory address or nil pointer dereference`、那条协程 return。**这复认了 `33-a2` ①-3 第 3 支的推**（"panic 被 observe 兜住、根 ctx 取消、球窗与面板同归于尽而进程还活着"），并把它从〔推〕升成〔量〕。
- ⚠ **但"栈会自己落盘"这半句要打折**（⑤-13）：栈是我**装了 sink 才有**的。出货代码里 `SetPanicSink` **零枚调用者**（R40），`goroutine.go:314` 的 `if sink != nil` 走空 ⇒ 今天真常驻炸的时候**盘上不会有任何栈**，只有 `:317-320` 那句含 recover 值、不含栈的 error 文本。`33-r5`／`33-v2` 若要拿栈，得**自己装**（我这两行的写法现成可抄）。

### §B-4 一枚我没预期、但会误导下一程的发现：**MTA 线程上创建＝当场 panic**

R32（`-com mta`）：`Embed` 的 `:87` **在调用内部同步**收到环境完成回调，`EnvironmentCompleted`（`chromium.go:171`）的 `if int64(res) < 0` **没拦住**那一次失败，于是 `:175 env.vtbl.AddRef(...)` 对 nil `env` 解引用 ⇒ panic（帧序列逐字见 R39）。三档对照：**STA `0x2` 创建成功（561ms）**、**未初始化也成功（531ms）**、**MTA `0x0` 当场 panic**。
- 对本票的意义（⛔ 选型）：`internal/ball` 的 `ui-sta` 用的正是 `0x2`（`win32_windows.go:149`＋`sta_windows.go:62`）⇒ **出货形状不踩这一支**。但 `cmd/wisp` 今天那枚真窗用例跑在 `hostThreadHarness`（`panel_host_windows_test.go:38-62`）那条形上——**现量尺**＝`grep -rn "CoInitialize" cmd/wisp/*.go` ⇒ **零命中**（R47），且 harness 线程里逐字只有 `:51 runtime.LockOSThread()` ＋ `:52 mgr.bringUp(ctx)` ＋ `:54 pnlPumpOnce()` ⇒ **今天出货宿主与测试 harness 都从未初始化过 COM apartment**，跑的就是我那发 `-com none`（531ms 成功）。而 `33-r5` 把它挪进真 `ui-sta` 之后，创建那一步会**第一次**发生在 STA 线程上（我这发 `-com sta`＝561ms 成功，同样不 panic）。⇒ **两枚线程形状今天不同**，谁把它们混用、或者哪天有腿图省事给宿主线程补一枚 `CoInitializeEx(…, 0x0)`，拿到的就是我 R32 那枚当场 panic。

---

## 交件判语（只两问）

| 问 | 结论 | 支别 | 凭据 |
|---|---|---|---|
| ① 再入 | **返回了但顺序被打乱**（不 panic、不永久阻塞） | **(c)** | R34／R35 ＋ 对照 R31／R33 ＋ 正控 R36 |
| ① 附：栈落盘 | 人造 `WM_QUIT` 那支 panic 的原始栈＝`%TEMP%\wisp33p1\reentry-quitduringcreate-sta\panic-sink.txt`；**出货进程今天装不到它**（零枚 sink 安装者） | — | R37／R38／R40、⑤-13 |
| ② 回执 | **自泵形＝不到；`Run()` 形＝到**（107ms 内拿到三枚回话）；`Eval` 主动推送在自泵形里**也到**（第三维，别混） | — | R25／R26／R27 |

**仍然判不了的（指名，供 `33-r5` 题面用）**：⑦-结-1 真球那条回调本身（复刻射程外）、⑦-结-2 球侧出帧会不会被插走、⑦-结-3 MTA panic 的 `res` 真值、⑦-结-4 `Run()` 与球的泵同线程共存、⑦-结-5 `dispatchq` 的真实长度。

---

## 收尾三把尺

| 尺 | 命令 | 读数 |
|---|---|---|
| ① 写面名册回到起手 | `git status --porcelain \| wc -l`；`git status --porcelain -- cmd internal \| wc -l`；`git diff --name-only e7eee03e..HEAD` | 见本节下方"同发取数"（**终值必须＝起手 227 行 ＋ 只 `.scratch/wisp/probes/33/p1/**` 新块**；`cmd`／`internal` ＝ **0**） |
| ② 主交付文件体量 | `wc -l -c .scratch/wisp/probes/33/p1/probe.md` | 见下 |
| ③ 占位符 | 尺＝`grep -nE -f .scratch/wisp/probes/33/p1/logs/placeholder-pattern.txt .scratch/wisp/probes/33/p1/probe.md`（那枚 pattern 存在**文件里**而不是写在这一行里——写在这一行会让这把自己命中，我第一发就中过，枚数 1→0 的差集就是这一行本身） | pattern 那族＝"骨架期括注 ＋ 四字母与 F 开头的两枚常见未完成标记 ＋ 中文那两枚同义标记"；终态命中＝**0**（`grep -cE -f` 返回退码 1 ＝ 零命中＝好消息） |

**同发取数（12:59:30 CST，HEAD 仍＝`e7eee03e`＝本腿自己那枚骨架提交）**：**逐字读数落在 `logs/final-rulers.txt`**，⛔ 不在这一节里复述文件大小——因为那一枚数是**自指**的（写进本文件就会让它再变一次；我先前给它回填过三轮，每一轮都把下一轮的数字改掉了）。摘要（口径＝那份日志）：
- `git status --porcelain \| wc -l` ＝ 起手 **227** → 终 **231**，`diff status-start.txt logs/status-final.txt` 的差集**恰好只有本腿自己的 5 行**：1 行 ` M` 是本文件，4 行 `??` 是 `p1/logs/`、`p1/receipt/`、`p1/reentry/`、`p1/msg-skeleton.txt`（原先那条折叠的 `?? p1/` 因为树里现在有已跟踪文件而展开成 4 行 ⇒ 净＋4）。⛔ **零枚别的写面**。
- `git status --porcelain -- cmd internal` ＝ **0**（与起手同值，R2/R43）。
- `git diff --name-only e7eee03e..HEAD` ＝ **空** ⇒ 我的全部真跑读数落在同一枚锚点上（R45）。
- `wc -l -c probe.md` ＝ **259 行／52,399 字节**——⚠ 这一枚是**回填前的值**，本文件因为写了这句话又会变大；请以上面那份日志为准，不要拿这一行当终值（这正是我把尺挪到文件外的理由）。
- 占位符尺＝**零命中**（`grep -nE -f logs/placeholder-pattern.txt probe.md` 退码 1）。
- 终态两枚闸门（12:57:03／12:57:23 复跑，原文 `logs/d22scan-final.txt`）：`GOFLAGS= go build ./...` **rc=0**、`sh scripts/d22scan.sh` **rc=0**。

| R48 | 终态复跑两枚闸门（本文件写完之后） | 逐字：`GOFLAGS= go build ./...` ⇒ 无输出、**rc=0**；`sh scripts/d22scan.sh` ⇒ `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=224, bans #1-5 cmd/=33, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=476, ban #8 cmd/=78`、**rc=0**（原文 `logs/d22scan-final.txt`，两步步的退码由脚本 `set -eu` 保证） |


