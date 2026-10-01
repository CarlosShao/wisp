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

---

## 载体形状（具名 ＋ 为 ②③ 两条跑过的命令）

**用了哪一形**：`.scratch/wisp/probes/33/p1/probe/main.go` ＋ `go run` 显式文件路径。
为什么**不是** `scripts/spike/**`：那是**独立模块**（R4 逐字 `module github.com/CarlosShao/wisp/scripts/spike`），本仓 `internal/...` 的可见性按**模块路径根**算 ⇒ 独立模块导不进来。〔推→本腿实测见 §⑥ 续表；测不出来就不许写这一句〕
为什么**不是** `cmd/wisp/testdata/<name>/main.go`：那仍在 `cmd/wisp/` 这棵树下，本票下一枚验收腿（`33-v2`）要逐名比 `cmd/wisp` 名册；派单第 1 条写死"⛔ 不给 `cmd/wisp` 测试包添文件"，我按**更严**的读法走（连 `cmd/wisp/**` 都不落）。
为什么**不是** `//go:build` 手工 tag 门：`go build ./...` 遇到"全文件被约束排除"的包会不会报 `build constraints exclude all Go files`，我没有读数支撑 ⇒ 不冒险（实测见 §⑥ 续表）。

为这一形跑过的命令（两枚闸门都要 rc=0，终态复跑）：
- `GOFLAGS= go build ./...`
- `GOFLAGS= go list ./...`（证明探针**不在**构建图里）
- `sh scripts/d22scan.sh`

---

## 问题② · "Go→页面"那一跳的投递者（正文待读数落盘后填写）

## 问题① · 在正在泵的线程回调里再入创建：panic／阻塞／乱序（正文待读数落盘后填写）

## 收尾三把尺
