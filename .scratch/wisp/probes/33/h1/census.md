# 33-h1 — 面板宿主「窗口真开出来」落地普查（只读，零产码，零 go.mod 改动）

**腿号**：`33-h1`　**派单**：编排者 2026-09-30 23:1x 收件后的只读普查腿　**票面**：`.scratch/wisp/issues/33-panel-host-c27.md`
**开工钟点**：`2026-09-30 23:16 +0800`　**分支**：`dev`　**起手 HEAD**：`77a3b442`（与取数同发，见 §8 尺 1）
**写点**：本文件一处（`.scratch/wisp/probes/33/h1/census.md`）。票面／产码／测试／配置／冻结件**零写面**。

---

## ① 依赖：`jchv/go-webview2` 的今日现状与本机可拉性

**尺**：S8·S9·S10·S11·S12·S13·S14·S15·S16·S18·S19·S45 ＋ 上游 tag 两发网络尺（**均失败，见下**）。

1. **库在不在**：本机的 Go 模块缓存里**两枚都在、且完整**——`go-webview2@v0.0.0-20260205173254-56598839c808` 与 `go-winloader@v0.0.0-20250406163304-c1995be93bd1`（S8），各自的 `cache/download/.../@v/` 目录里 `.info .mod .zip .ziphash` 齐备（S9）。`cd scripts/spike && go list -m all` **不联网即列出这两枚**（S19）。
2. **⚠ 上游"最新 tag"我没核到，具名登记**：`WebFetch https://github.com/jchv/go-webview2` **fetch failed**（本机网络），`WebSearch` 只回到第三方博客、无 tag 页。⇒ 本普查只能说"**仓内钉的这枚 pseudo-version（2026-02-05 提交）本机可解析、可编译成 spike 产物**"，**不能**说"这就是上游最新"。仓内那枚 pseudo-version 逐字在 `scripts/spike/go.mod:9`（`02-spike-report.md:20` 同一条）。**这一条要补的话得给网络或给 `go list -m -versions`（要出网，且本程未跑）。**
3. **纯 Go 还是 cgo**：**纯 Go，零 cgo**。`import "C"` 零命中（S11，退码 1）、`#cgo`／`go:cgo_import_dynamic` 零命中（S12）。库全程用 `golang.org/x/sys/windows` ＋ COM vtable 直调（`pkg/edge/comproc.go` 那族 `NewComProc`）。⚠ 两个必须写清的细节：① `go-winloader` 里**有一枚 C 源件** `tinydll/tinydll.c`（S12），但它**不参编**（无 `import "C"` 对偶），是"被嵌成字节"的样板；② 库**把 3 枚 `WebView2Loader.dll` 用 `//go:embed` 打进二进制**（S10：`webviewloader/sdk/{x64,x86,arm64}/`；`webviewloader/module_amd64.go` 逐字 `//go:embed sdk/x64/WebView2Loader.dll`），磁盘取不到就 `winloader.LoadFromMemory` 走内嵌那份（`module.go` 的 `loadFromMemory`）。⇒ **"不需要 C 工具链"成立；"运行期不需要外部 DLL"也成立；但"二进制里多了 3 枚 BSD 许可的原生件"这一条要落进打包/许可账（见第 6 点）。**
4. **它依赖哪些包 / 要加哪几名**（⛔ 本程 `go.mod`／`go.sum` **一字节未动**，只普查名单）：
   - `go.mod` 需新增 **2 枚 module**：`github.com/jchv/go-webview2`（直接）＋ `github.com/jchv/go-winloader`（**indirect**，S14/S15）。
   - `go.sum` 需新增 **4 行**（2 枚 × {`h1:` 行, `/go.mod h1:` 行}）。逐字模板已在仓里：`scripts/spike/go.sum:13-16`（S18）；根 `go.sum` 现量 `grep -c jchv`＝**0**（S18）。
   - **`golang.org/x/sys` 不用动**：两枚库分别要求 `v0.0.0-20210218…`（S14）与 `v0.0.0-20200810…`（S15），根 `go.mod` 已是 `v0.48.0`（S5），MVS 取高 ⇒ 不抬版本、不加新行。
   - **包级新增名册 9 枚**（S16，`go list -deps github.com/jchv/go-webview2/pkg/edge` 总 89 枚包、外部 10 名、减去已有的 `x/sys/windows`）：`go-webview2/pkg/edge`·`go-webview2/internal/w32`·`go-webview2/webviewloader`·`go-winloader`·`go-winloader/internal/{loader,pe,vmem,memloader,winloader}`。若改用高层 `github.com/jchv/go-webview2`（S24）再加 1 名＝**10 枚**。
5. **本机能不能离线拉到**：**能**（S8/S9/S19）。但两条要具名的现实：① 本仓**无 `vendor/`**（S4/S6，`.gitignore:10` 还写着 `vendor/`）⇒ 取依赖走的是 module cache，CI 那侧靠 `actions/setup-go` 的缓存 ＋ `ci.yml:441/:543` 的 `Cache third_party`；② **GOPROXY 默认是镜像**：`go env GOPROXY`＝`https://goproxy.cn,direct`（S7），且 `scripts/build.ps1:119` 逐字 `if (-not $env:GOPROXY) { $env:GOPROXY = 'https://goproxy.cn,direct' }`。⇒ 真缺件时会**先从镜像站拿 `go.sum` 的 hash**。这与禁止清单里的 ban #5 `mirror-hash` 是不是同一件事，**我判不动 ⇒ ⑩ J5**（今天 `scripts/spike/go.sum` 那 4 行就是从同一 proxy 拿的，S18/S44 语境）。
6. **许可与登记落点**：库 `LICENSE`＝MIT、`webviewloader/sdk/LICENSE.txt`＝3-Clause BSD（S10）。⚠ **本仓没有"Go 模块许可登记册"这个东西**：`deps.toml` 头注逐字是"pinned **native** dependencies (SPEC-11 §2.1)"、全文 `license =` 只 3 处（S44），尺 `grep -rln "LICENSES|licenses|NOTICE" scripts tools docs/specs`＝零命中。⇒ 派单/前程说的"`deps.toml` 许可登记"这一动作**在本票射程里没有现成格子**；这是**登记面的缺口**，不是我可以顺手补的东西（**⑩ J6 附属**）。
7. **与规格互证**：`docs/specs/SPEC-08-ui-ball-panel.md:145` 逐字"`jchv/go-webview2`（MIT，纯 Go 无 cgo）"（S45）⇒ **票面 `:25` 的 "pure Go" 与规格、与我这三把尺（S11/S12/S13）四方一致**，这一条不是推断。

**结论**：依赖侧**没有拦路虎**——库存在、纯 Go、离线可解、只多 2 枚 module／9 枚包名、`x/sys` 不动。**真正的依赖风险不在"能不能拿到"，在"它把 3 枚 BSD DLL 嵌进二进制"与"许可登记无现成格"这两条记账面。**

---

## ② 宿主落点：新包 vs 装配根文件，两形各列出新增的包级依赖边

**尺**：S16·S17·S24·S38·S49·S53·S61（`PanelManager` 名册）·S62（`cmd/wisp/panel_assets.go` 头注）。

**先摆两件既有事实（不是我的推广）**：
- **C27 的契约原文**（`docs/PLAN.md:1377`，S61）："`PanelManager` 单例，**唯一** WebView2 窗口持有者；一会话至多一个面板窗口，**隐藏而非销毁**；面板不可用 → L2 降级为原生最简确认卡，**L2 能力不得消失**；多任务共用，按 correlationId 分区渲染"。⇒ 契约只规定**所有权形状**，**没规定它落在哪个 Go 包**。
- **仓里已经写明宿主该用什么**：`cmd/wisp/panel_assets.go:13` 逐字 **"Ticket 33's WebView2 host is expected to call the same panel.Assets API this command prints."**（S62）⇒ 供给面复用 `panel.Assets`（S25：`Resolve:78`／`Built:57`／`EntryFile:31`／`Check:158`），**宿主不必新造读字节的腿**。
- **`internal/panel` 今天是一枚零平台分叉包**：`find internal/panel -name "*_windows*.go"`＝**0 枚**（S53），而它在 **ubuntu CI 的 core scope 里有分母**（`portable-tests.sh` `core_pin` 含 `internal/panel`，S38 `:141`）。⇒ **谁往 `internal/panel` 里塞第一枚 `_windows.go`，谁就把这包一半代码搬出 ubuntu 分母**——这是本票落点选择的**硬约束**，与两条既有裁定无关。

**形甲：新包 `internal/panel/host`（宿主独占，`//go:build windows`）**
新增包级依赖边（相对 S17 的现量名册）：
- `internal/panel/host` → `github.com/jchv/go-webview2/pkg/edge` ＋ `webviewloader` ＋ `internal/w32` ＋ `go-winloader` 家族（**9 枚新包名，S16**）
- `internal/panel/host` → `github.com/CarlosShao/wisp/internal/panel`（取 `Assets`／`ComposerDispatch`；**panel 的名册不涨**，151 保持）
- `internal/panel/host` → `internal/observe`（`Registry.Spawn`＋`Timeout`，S35/S47；observe 已在 panel 名册里，不算新增名）
- `cmd/wisp` → `internal/panel/host`（**1 枚新边，装配根 → 宿主**）
- ⇒ `internal/panel` 名册**零变化**；`internal/panel` 的 ubuntu 分母**零损失**。
- 代价：`host` 包**在 ubuntu CI 上不参与测试**（对照 S70：`internal/ball` 30 枚 `.go` 里 **10 枚带 `_windows`**、`cmd/wisp` **9 枚带 `_windows`** ⇒ 那两包本来就是平台分叉包，只有 `internal/panel` 是零，S53），`host` 能拿到的 CI 分母只有 staticcheck 的 `GOOS=windows` 一发（S49）＋ `./internal/panel/...` 那枚 glob 把它收进 core scope 的**编译**射程（S38 `:179`）。

**形乙：装配根文件 `cmd/wisp/panel_host_windows.go`**
新增包级依赖边：
- `cmd/wisp` → 那 **9 枚 jchv 包名**（S16）＋ `cmd/wisp` → `internal/panel`（**已有**，S62 的 `panel_assets.go` 就 import 它）
- ⇒ **`internal/panel` 同样零变化**；但**原生 COM／unsafe 代码进装配根**，`cmd/wisp` 的 windows-tagged 面从"注入缝"扩成"宿主本体"。
- 代价：`cmd/wisp` 的测试二进制（`cli` scope，S38 `:195`）**今天就要在装了 sherpa DLL 的机器上才跑得动**（`portable-tests.sh:191-194` 注释逐字："the test binary dies at LOAD time without the sherpa DLLs on PATH (ticket 98)"）⇒ 宿主代码压进这枚包，等于把宿主的测试分母**并到那台已经最窄的载具上**。

**关于那两条既有裁定——按派单要求具名写成"这是我的推广"**：
- 票 197「装配根 `cmd/wisp` 是唯一的接缝」：我在 ② 里读到的**原文形状**是"`cmd/wisp/panel_assets.go` 里放一枚只读诊断子命令"（S62），它**不是**"一切原生代码都塞 `cmd/wisp`"。**⇒ 我若据此主张"宿主必须落 `cmd/wisp`"，那是我的推广，不是仓里本来就有的规矩。**
- 票 238「正向依赖边一律不开、改注入」：我**没有读到**它对"原生窗口宿主"说过一个字。**⇒ 我若据此主张"宿主必须是独立包、由 `cmd/wisp` 注入"，同样是推广。**
- 另有一条**未定案就在票池里**：`docs/evidence/s1/114-ac1-status-table.md:203` 把"票 114 与票 33 的地界"列成「**未定，需 owner/编排者拍**」。⇒ 落点这件事**今天没有权威答案**，**⑩ J4** 摆甲／乙／不做交裁。

**我的倾向（明确标注为倾向，不是规矩）**：**甲**。只有一条读数支撑它——**S53＋S38 那对因果**（`internal/panel` 是零平台分叉、且它在 ubuntu 有分母；两形里只有"新包"这一形能让宿主不碰那枚名册）。C27 的单例所有权放甲是 `host.Manager`、放乙是 `cmd/wisp` 的包级变量，**契约两边都容得下**（S61 原文未指包）。

**结论**：两形的**依赖增量都是同一批 9 枚 jchv 包名**（S16），差别不在名册长短，在**"这 9 枚挂在哪枚包上、那枚包今天有几份 CI 分母"**（S38/S49/S53/S63）。

---

## ③ STA／消息循环共存（本票最要命的一问）

**尺**：S31·S32·S21·S22·S23·S24·S33·S35·S36·S47（全部是逐行读，无一条"通常没问题"）。

**先摆"今天已经有几枚泵"**：常驻进程的球跑在 `ui-sta` 上——`runtime.LockOSThread()`（`sta_windows.go:52`，S31）→ `CoInitializeEx(0, coinitApartmentThreaded)`（**:62**，S31）→ `for { GetMessageW → TranslateMessage → DispatchMessageW }`（**:77-93**，S31）。投递口是 `PostTask`（**:127-143**，S32/S31 互证），执行口是 `runTask`（**:146-154**），同步等待口是 `Ball.uiRun`（`ball_windows.go:741-752`，S32：`GetCurrentThreadId()==threadID()` 则**就地跑**，否则 PostTask＋`<-done`）。球的退出走 `Close`（**`:935-968`**，S32：PostTask 清理 → `sta.quit()` → `<-handle.Done()` join）。D38a 的原文注释就写在 `sta_windows.go:5-9`（S31）："ONE thread owns all UI COM objects (… and later the panel) … goroutines elsewhere must never touch window/COM state directly - they PostTask"，以及 `:60-61`"D2D single-threaded factory **+ future WebView2 panel** both want STA"。⇒ **"共用 ui-sta"是这枚文件自己预告的形状**，票面 `:30` 的 Key constraint 与它逐字同向。

**下面是四条会真出事的地方，每条给行号。**

**K1｜库的 `Embed` 自带第二枚泵，压在调用线程上。**
`pkg/edge/chromium.go:96-111`（S21）：`for { if e.inited!=0 break; GetMessageW(0,…); if r==0 break; Translate; Dispatch }`，出口 `inited` 只在 `CreateCoreWebView2ControllerCompleted:224` 置。**⇒ 谁调 `Embed`，谁的泵就被这段循环接管，直到环境＋控制器都建好**（本机冷启这段是秒级：S34 里 879.7–1808.9 ms）。若投在 `ui-sta` 上，那几秒里球**不在自己的循环里、在库的循环里**。两枚循环都是 `GetMessageW(hwnd=0)`＝**取本线程队列全部消息**，所以球的 `WM_TIMER`/`wmAppTask` **理论上仍被派发**（这半句是 Win32 线程队列语义的推断，我没跑过 ⇒ ⑨ 第 3 条，落地腿第一枚用例就该证它）。

**K2｜`WM_QUIT` 会被库那份泵吃掉，退出路径有两处 `continue` 语义。**
`sta_windows.go:156-162`（S31）的 `quit()` 用 `PostQuitMessage`；库那份循环看到 `GetMessageW` 返回 0 就 `break`（`chromium.go:106-108`，S21）并**把这条 WM_QUIT 消费掉**。⇒ **如果 `Ball.Close()` 落在冷启窗口内，退出信号可能被 `Embed` 吞掉而不是交给球自己的循环**，`start()` 后半段（`:95 releaseCOM`）与 `Close` 的 `<-handle.Done()`（S32 `:966`）谁先谁后就不可读。**这一条可以纯本机量**（载具先例：`internal/ball/live_windows_test.go` 那一族 `winlive`，S64），我登记给落地腿，不在普查里判它对程序正确性的净影响。

**K3｜WebView2 的回调会直接跑在 `ui-sta` 上，绕开 `runTask` 那把锁的语义。**
STA 的含义就是"这套间里的 COM 对象只由这枚线程调"。库把 4 枚 handler 注册进 WebView2：`AddWebMessageReceived`（`chromium.go:201`）、`AddPermissionRequested`（`:206`）、`AddWebResourceRequested`（`:211`）、`AddNavigationCompleted`（`:216`）（S20/S21）。它们的 `Invoke` 是 COM 调用，**由 `DispatchMessageW` 在创建线程上触发**。⇒ 三条后果，逐条对上本仓既有契约：
- **`Chromium.MessageReceived`（`:233-248`，S21）会在 `ui-sta` 上直接调 `e.MessageCallback`** ⇒ 那就是入向那一跳的手（票面 `:94` 的 H3）。回调里若走 `ComposerDispatch.Handle`（票面 `:142` 逐字给出 `Handle:120`）**必须是快且非阻塞**——这与 `ball_windows.go:724-727`（S32）给 `fire` 写的契约**同一句话**："callbacks must be quick and non-blocking … a blocking callback stalls the message pump"。⇒ **本票落地的第一枚真听众，天生压在一条"不许阻塞 UI 线程"的既有契约上。**这条是我这一问里最硬的读数。
- 更要命：`MessageReceived` 末尾**无条件回灌** `PostWebMessageAsString(sender, message)`（`chromium.go:242-245`，S21）——库把收到的原串**广播回页面**。这与票面 `AC#9` 的"拒答不外泄"形状、与票 35 的回执设计**都会打架**（回执该由 Go 决定内容，不是把 raw 原样回吐）。⇒ **落地腿若用库的默认 handler，得先处理这一支**；这是宿主形状问题、不是桥协议问题，所以我不越界去票 35 的地盘，只登记。
- **`Chromium.WebResourceRequested`（`:281-290`，S21）也在同一线程上跑**，里面 `args.GetRequest()` 失败时**逐字 `log.Fatal(err)`（`:284`）**；而 `AddWebResourceRequestedFilter`（`:292-297`）自己在 err 上也是 **`log.Fatal`（`:295`）**。⇒ **供给腿上任何一次 COM 失败都是整进程退出**，与票面 AC#5「no-crash」直接对撞。

**K4｜库的错误路径有两条，一条安全一条致命，而且我判不动这台机器会走哪条。**
同步支：`createCoreWebView2EnvironmentWithOptions` 报错 → `log.Printf` ＋ `return false`（`chromium.go:87-94`，S22）＝**不崩，可判 unavailable**。异步支：`EnvironmentCompleted` 收到负 HRESULT → **`log.Fatalf`（`:173`）**；`CreateCoreWebView2ControllerCompleted` 负 HRESULT → **`log.Fatalf`（`:188`）**（S22）＝**进程死**。⇒ 票面 `AC#5` 要的"creation failure → `panel.unavailable` event + no-crash"，**在缺 runtime 那一发上取决于失败落在哪一支**，而本机装了运行时（S46）⇒ 本票在这台机器上永远走同步可达路径、**证不到那一支**。⇒ **⑩ J2**，并给出一条现成的预检：`webviewloader.GetInstalledVersion()` 在没装时**返回 `("", nil)`**（`webviewloader/module.go`，S13/S10 同一包），创建前先问它一把，就能把"缺件"这一支从异步 `log.Fatalf` 提前挪到同步可判——⚠ **我没实跑过（⑨ 第 5 条）**。

**K5｜高层 API 必然多一枚线程，因此"守 D38a"＝只能用 `pkg/edge`。**
`webview.go` 的 `NewWithOptions:97` → `CreateWithOptions:269` 自己 `RegisterClassExW`（类名逐字 `"webview"`）＋自己建窗，`Run:351` 是**第三枚**死循环，`Dispatch:443` 用 `PostThreadMessageW(w.mainthread, WMApp)`（S24）。spike 用的正是这一形，且它把 `Run()` 换成手工 `PeekMessageW` 泵，并在 `main.go:113-115` 具名写过坑（S33）。⇒ **要"WebView2 created on the shared ui-sta"（票面 `:30`），高层 API 这条路走不通**（它要独占一枚跑 `Run()` 的线程）；只有 `pkg/edge.Chromium` ＋**我们自己的 HWND**（`Embed(hwnd)` 收现成句柄，`chromium.go:72`，S21）能压在 `ui-sta` 上。

**K6｜"另起一枚常驻 STA 线程"这个名字在名册里今天不存在。**
`observe.ResidentNames` 六枚＝ui-sta·audio-capture·hotkey-listener·db-writer·watchdog·log-flusher（`goroutine.go:44-46`，S35），`ResidentBaseline = 6`（`:39`）。`panel-host` 在 **`TemporaryNames`（`:58`）**，而 `SPEC-01-architecture.md:122`（S36）给它写的注是"**均在对应 DisposalScope 内**"。⇒ **独立线程那一形要么占用一枚不在常驻名册的名字（`Spawn` 会 `slog.Warn("goroutine outside the D38 roster (leak symptom)")`，`goroutine.go:271`，S35），要么把"临时·处置范围内"的名字拿去当常驻**——两条都是**动 D38b 名册＝契约变更＝人工批准**（`AGENTS.md` §1.1）。⇒ **⑩ J1**。

**三种接法与代价（只列，不替谁定）**：
| 接法 | 形状 | 代价（全部带尺号） |
|---|---|---|
| **甲：ui-sta 上跑 `Embed`（票面 `:30` 原样）** | `PostTask(func(){ edge.Embed(panelHwnd) })` | 冷启几秒里泵的**控制流在库里**（K1）；WM_QUIT 有被吞的一面（K2）；回调直接落 ui-sta，听众必须非阻塞（K3）；`log.Fatal` 两支都在（K4）；**名册零变更**；**依赖最少**（S16 那 9 枚） |
| **乙：ui-sta 上自写异步创建，不用 `Embed`** | 自己 `webviewloader.CreateCoreWebView2EnvironmentWithOptions` ＋自己实现 COM handler | **避开 K1/K2**（无嵌套泵，回调照常走球的循环）；但要**自己写未导出的那枚 environment-completed handler**（S23：`iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler` 全小写、`pkg/edge/` 里无 Environment 那枚文件）⇒ **新增约百行 unsafe vtable 代码**，是本仓最贵的一类产码；K3/K4 仍需自行处置 |
| **丙：独立线程跑高层 `webview.NewWithOptions`** | spike 那一形搬进常驻进程 | 与 K5 直接冲突（`Run()` 独占线程）＋**撞 D38b 名册**（K6）⇒ **要人工批准**，本票不许自行选它 |

**结论（一句）**：**会互抢的不是"两个套间"，是"同一枚线程上的两份泵"**——`ui-sta` 本身能容 WebView2（STA 正是它想要的，S31 `:60-61` 自己写了），但 `pkg/edge.Embed` 在冷启期间**把那份泵换成库里的一份**（K1/K2），并且**把 4 类回调直接投在球的家门口**（K3）。要避开就只剩乙那一形（自写异步创建，代价＝unsafe 代码量），或者接受甲并把"退出顺序"钉成一发用例。**我这一问没有一条能靠"通常没问题"交差，也确实交不出实测——实测需要载具，属落地腿。**

---

## ④ embed 供给：路径存在性 ＋ "空目录会怎样"

**尺**：S25·S26·S27·S28·S29·S30·S52·S63（Go 官方 `embed` 文档原句）＋ 库侧供给链。

1. **票面那条路径今天不存在，两种口径都不存在**：`ls -d assets assets/web assets/web/dist` ⇒ **三枚全部 `No such file`**（S27）；`ls -d internal/panel/assets internal/panel/assets/web` ⇒ **同样不存在**（S28）；`find . -maxdepth 3 -type d -name assets` ⇒ 全树只有 `./design/old/assets` 与 `./frontend/dist/assets`（S52）。**⚠ 更深一枚：`.gitignore:24` 逐字写着 `assets/web/dist/`，整条忽略、连锚文件豁免都没有**（S29；`git diff` 现量证明这条**是 HEAD 里的既有规则**，不是别人今天新加的，S50）。⇒ 票面那句 `//go:embed assets/web/dist` 是**一条被 gitignore 封住、且没有任何落点的路径**。
2. **仓里活的 embed 不是它**：`internal/panel/assets.go:38` 的 `BuiltinAssets()` 走 `fs.Sub(frontend.Dist(), "dist")`（S25），即嵌在 `frontend` 那枚 Go 包里；二手引仓内证据件：`docs/evidence/s1/156-close-and-load-bearing-r1-accept-r1.md:762` 逐字 `frontend/embed.go:19 → //go:embed all:dist`（S30）。**我按两层禁令没读 `frontend/**` 一字**，只取了存在性／大小／路径形状：`frontend/dist` **在**，4 枚文件——`.gitkeep` 0 B／`assets/index-BRKj5OIJ.css` 49943 B／`assets/index-BVKlegVD.js` 553469 B／`index.html` 1044 B，合计 **604456 B**（S26）。
3. **"embed 一个空目录会怎样"——答案是构建期硬失败**，Go 官方文档三句原话（S63，`pkg.go.dev/embed`）：
   - "If a pattern names a directory, all files in the subtree rooted at that directory are embedded (recursively), **except that files with names beginning with '.' or '_' are excluded**."
   - "**Matches for empty directories are ignored.**"
   - "After that, each pattern in a `//go:embed` line must match at least one file or **non-empty directory**. If any patterns are invalid or have invalid matches, **the build will fail.**"
   ⇒ 三条合起来对本票是**一个陷阱**：锚文件叫 `.gitkeep`（以 `.` 开头），**裸 `//go:embed assets/web/dist` 收不到它**（第一条排除），于是清检出该目录匹配不到任何件 ⇒ **`go build` 直接红**（第三条）。仓里现用的那枚写的是 **`all:dist`**（带 `all:` 前缀，才把 `.gitkeep` 收进匹配），这正是 `frontend/.gitignore:12-13` 与 `.gitignore:22-23` 那对"`dist/*` ＋ `!dist/.gitkeep`"能成立的前提（S29）。⇒ **票面那条 embed 一旦原样落进产码，就把"清检出构建"打成红的**——这正是派单问"这决定构建期会不会硬失败"的那一刀，**答案是会**。
4. **宿主侧的供给链已经齐了，而且全在库的导出面上**（尺：S20·S21·库侧 grep 结果）：
   - 过滤器：`Chromium.AddWebResourceRequestedFilter(filter, ctx)`（`chromium.go:292`）＋ vtable 绑定 `ICoreWebView2.AddWebResourceRequestedFilter`（`corewebview2.go:383`）。
   - 事件回调字段：`Chromium.WebResourceRequestedCallback func(request *ICoreWebView2WebResourceRequest, args *ICoreWebView2WebResourceRequestedEventArgs)`（`chromium.go:42`）。
   - **取 URI**：`ICoreWebView2WebResourceRequest.GetUri()`（导出，`ICoreWebView2WebResourceRequest.go:29`）。
   - **造响应（含头）**：`ICoreWebView2Environment.CreateWebResourceResponse(content []byte, statusCode int, reasonPhrase string, headers string)`（导出，`corewebview2.go:169`）＋ `Chromium.Environment()`（`chromium.go:299`）。
   - **交响应**：`ICoreWebView2WebResourceRequestedEventArgs.PutResponse(...)`（导出，`ICoreWebView2WebResourceRequestedEventArgs.go:27`）。
   ⇒ **完整闭环 `GetUri → panel.Assets.Resolve → CreateWebResourceResponse(bytes, 200, "OK", headers) → PutResponse`，四步全在导出 API 上，不需反射不需新依赖。**
5. **票面要的 CSP 注入点＝就是上面那枚 `headers` 形参**（`CreateWebResourceResponse(content, code, reason, **headers**)`）。⇒ "provide the injection hook ＋ default strict CSP now"这一格**今天有位置可放，不必新造接口**；⚠ 另一条**别选错**：库也导出了 `ICoreWebView2_3.SetVirtualHostNameToFolderMapping(hostName, folderPath, accessKind)`（`ICoreWebView2_3.go:22`）——**它要的是磁盘真实目录**，与 embed 供给是两条路。票面与 D29 要的是 filter 那一支（S45），**virtual-host 那一支会把产物落回磁盘**，本普查具名提醒别混。
6. **"无 localhost server"这一侧现状**：全仓 `net.Listen` 的非 test 命中只有 `tools/mockllm/server.go:243`（测试用 mock，不在产码路径），`internal/panel` **零 `net` 依赖**（S39/S17）。⇒ 票面 D21/D29 那条反模式约束**今天是天然满足的**，本票要做的是**让它继续可证**（见 ⑤）。

**结论**：**票面 `:28` 那条 embed 路径与仓里活的那枚不是同一枚**（S27/S28/S30），而"空目录"问题的答案是**硬失败**（S63）。落点二选一（复用 `panel.Assets` 还是新立 `assets/web/dist`）⇒ **⑩ J6**。

---

## ⑤ 「无监听端口」这条判据怎么机读 ＋ CI 能不能真跑到

**尺**：S37·S38·S39·S40·S49·S53·S64。

**判据形状：分两层，一层能在 CI 跑、一层只有真机能跑，且都必须具名。**

**L1（CI 可跑，能力型，非词面）——本普查给的主张**：断言**宿主包不引入监听能力**，问的是 import 名册而不是注释词面：
- 谓词＝"扫描 `internal/panel/host`（或 `cmd/wisp` 的宿主文件）的 **`go/ast` import 声明**，出现 `net`／`net/http`／`google.golang.org/grpc` 即红"。
- ⚠ 为什么不是 `grep -i netstat`：**全仓 `netstat` 字样现量＝0 命中**（S39），拿词面做门只会永远绿——那正是本仓禁的"恒真判据"那一族。
- ⚠ 为什么它**不是**扫注释：`go/ast` 只取 `ImportSpec`，注释与字符串不入射程（这一口径与本票 `AC#9` 那枚改造后的能力尺同源，票面 `:89-90` 已把"注释与字符串字面量不入射程"写成在册口径）。
- 分母：`internal/panel` 在 ubuntu 的 `core_pin` 里（S38 `:141`），若宿主是新包 `internal/panel/host`，它**被 `./internal/panel/...` 那枚 glob 直接覆盖**（S38 `:179` 逐字 `./internal/panel/... ./internal/ball/`）⇒ **这条判据在 ubuntu-latest 上有真分母**，不需要真机。⚠ 前提是新包落点选 ② 的甲形。
- 补一发**正控**（沿用票面 `:92` 已确立的"尺自带正控"形状）：在仓外临时目录建一棵含 `import "net/http"` 的假宿主 ⇒ 必红；把那条 import 摘掉 ⇒ 必绿。

**L2（真机才成立那一层）——要机读，但不能只留它**：枚举**本进程（含进程树）持有的监听 socket**：
- 现成能借的只有一枚，且**不好用**：`internal/proc/treemetrics_windows.go:269` 的 `treeTCPConnections(inTree map[uint32]bool) int`（S40）——它**未导出**、用 `GetExtendedTcpTable(AF_INET, TCP_TABLE_OWNER_PID_ALL)`（`:47`）但**只数行数、不读 `state` 字段**（偏移只取 `tcpOwningPidOffset`），失败重试耗尽后 `return 0`。⇒ **"有 TCP 行"≠"在监听"**，直接拿它会判错。
- 能力型形状要的是：**同一张表里过滤 `dwState == MIB_TCP_STATE_LISTEN(10)`，并且 AF_INET6 也扫一发**（`GetExtendedTcpTable` 的 `ulAF` 参数），断言**计数＝0**。⇒ 要么在 `internal/proc` 导出/新增一枚读 state 的枚举（**动的是票 156/246 那一族的地界**），要么宿主测试自带一份枚举（新增 unsafe/偏移代码）。两条都在派单外，**我只报形状，不选**。
- ⚠ 这条真机判据的**机器归因**必须写在读数旁边：WebView2 自己的子进程（`msedgewebview2.exe`，S46 那两枚版本目录里的可执行件）**本来就可能持有 socket**（DevTools、网络服务），所以"树内 PID"这一维**不能简单排除 WebView2 子进程**，否则判据测的是"我看不看得到它"。**这条我判不动它该怎么划界**（属票面未覆盖）⇒ 记进下面 CI 段的待裁。

**CI 到底能不能跑到——三条现量，一条比一条窄**（S37/S38/S49/S53/S64）：
| 岗位 | runs-on | 宿主代码/用例能不能跑到 | 现量凭据 |
|---|---|---|---|
| lint（staticcheck 等） | `ubuntu-latest`（`ci.yml:66`） | **代码能**（类型检查两平台形状，逐字"34 (GOOS=linux) / 78 (GOOS=windows)"，`ci.yml:199`，S49）／**用例不能** | S37·S49 |
| portable-core | `ubuntu-latest`（`ci.yml:278`） | `./internal/panel/...` 在 scope（S38 `:179`）⇒ **L1 那层能真跑**；但 `_windows.go` 在 linux 不参编 ⇒ 宿主用例零分母 | S38·S53 |
| **test-windows** | `windows-latest`（`ci.yml:388`） | **`internal/panel` 不在 `win_pin`**（S38 `:149-158` 那 8 枚里没有它）⇒ **今天零分母**；要它就得改门禁名册 | S37·S38 |
| build wisp.exe ＋ SLO smoke | `windows-latest`（`ci.yml:533`） | **链接得到**（`scripts/build.ps1:123` 逐字 `go build -trimpath -ldflags $ldflags -o build/wisp.exe ./cmd/wisp`，S65）⇒ 编译期红能抓，**用例仍不跑** | S37·S65 |
| SLO full | `[self-hosted, wisp-slo]`（`ci.yml:591`） | **只有开机那台机器**⇒正是"同形坑"那一族 | S37 |
| **`-tags winlive`** | —— | **`grep -rn winlive .github/workflows scripts/*.sh scripts/*.ps1` 现量＝零命中**；用它的是 6 枚文件（`internal/ball/{live_windows,interaction_live,hotkey_live,live_guard_windows}_test.go`、`cmd/wisp/resident_ball_live_228_windows_test.go`、`cmd/wisp/resident_approval_live_246_windows_test.go`）（S64）⇒ **本仓没有任何 CI 岗位跑 winlive** | S64 |

⇒ **结论**：**只留 L2（真机 TCP 枚举）＝第四枚同形坑**，我明确不推荐单独交它。**要交的是 L1 进 ubuntu 分母 ＋ L2 具名"读数出自哪台机器、哪个 tag"，且跑不到的那一格必须登台账 `A##`**——改门禁分母这件事本仓有在册前例：`A374` 里编排者已自决"暂不把 `internal/panel/` 加进 windows scope"（票面 `:139` 第 ② 条逐字），⇒ **加它要单独落一枚 `A##`，不是写腿能自己做的**。⇒ **⑩ J7**。

---

## ⑥ 延迟两数（冷 ≤1500ms／热 ≤200ms）的测量方法普查

**尺**：S33·S34·S47·S53·S64·S66（`thresholds.go`／`sampler.go` 的 panel 维度）。⛔ **本程对 `docs/SLO.md`、`internal/observe/thresholds.go`、golden、`tools/d22scan/allowlist.txt` 零字改动**（S48/S66 为自证尺）。

**1. 谁起秒表——必须是宿主自己在 UI 线程上取单调钟。**
`internal/observe/clock.go:10-19`（S47）逐字："Every timeout, deadline, retention window and interval comparison MUST be computed on the monotonic clock… Computing intervals from wall-clock readings … **is BANNED**"；具名机制＝`observe.Timeout`（`NewTimeout:34`／`Elapsed:39`，S47）。⇒ 两数**不能**用"两次 `time.Now()` 的墙钟差"实现（d22scan ban #4 `wallclock-timeout`，S41），也不能用带 `Round(0)`／编解码过的值。spike 的 `time.Since(t0)`（`main.go:139`，S33）就是单调口径，**沿用它的写法**。

**2. 定义必须逐字沿用 S0 那一套，否则读数不可比。**
`02-spike-report.md:82-84`（S34）：**cold** ＝ `NewWithOptions`（建窗＋show＋`Embed` 阻塞至 Environment/Controller 就绪）＋ `SetHtml` ＋ **首次浏览器往返**（`Bind("spikePing")` ＋ `Eval("spikePing()")` 回调到达）；**hot** ＝ `ShowWindow(SW_HIDE→SW_SHOW)` ＋ 一次泵 ＋ 往返；**recreate** ＝ Destroy ＋ 完整再建（同进程重建 Environment）。
⚠ **本票落地会多出两件 S0 口径没覆盖的成本**：① 走 `AddWebResourceRequestedFilter` ＋ `Assets.Resolve` 取字节（④ 第 4 那条闭环，替代 `SetHtml`），② CSP 头注入。⇒ 我的建议：**另立第三枚口径 `cold-embed`**，与 S0 的 `cold` 并列报，**不替换原口径**；否则"换了定义把 1500 线蒙过去"正是本仓那族假绿。**这是测量方法层面的主张，不是数。**

**3. 几次取样、P50/P95 放哪。**
- 取样：票面 `AC#2` 逐字"**10-run P50/P95**"；spike 的先例是 12 枚子进程各测一次真冷 ＋ 一记 warmup 弃样（`02-spike-report.md:85-87`，S34；载具 `main.go:259-279`，S33）⇒ 冷启**必须跨进程/跨会话取**，同进程内第二次不是冷（那叫 recreate：P50/P95 955.7/1221.3 与 859.0/955.5，S34）。热那一发在同窗口内 hide→show 循环 10 次即可（`main.go:209-224`，S33）。
- 存放：**只进证据件**（`docs/evidence/s1/33-*.md`）。现量支持这个选择：`internal/observe/thresholds.go` 里与 panel 有关的**只有内存与 CPU 两枚**——`memCapPanel = 600<<20`（`:23`）／`cpuLimitPanel = 10.0`（`:30`）／`SLOPanelOpen` 分支（`:55`/`:72`）（S66），`internal/observe/sampler.go:46/:51` 里 `SLOPanelOpen = "PanelOpen"` 是六枚 SLO 态之一——**没有任何延迟毫秒字段**。⇒ **延迟两数今天不在阈值文件里，本票也不该提议把它加进去**：加阈值＝动 SLO＝**人工批准**（`AGENTS.md` §1.1"SLO 阈值…一字节都不许动"）。
- ⛔ 我**不动** `docs/SLO.md`——并具名登记它**不存在**（`AGENTS.md` §4 末段逐字列它为"截至锚点 `4e66817` 在仓里不存在"的交付物之一），所以票面 `AC#2` 说的"SLO appendix"今天**没有落点文件**⇒ 落点这件事本身要裁（我把它并进 ⑩ J3 的注里）。

**4. 本机现状（转引在册读数，非本程实测）。**
S34 两 run：cold P50/P95 ＝ **879.7/1041.6** 与 **1125.7/1256.4 ms**（`02-spike-report.md:148`）；hot show P50/P95 ＝ **25.7/49.5** 与 **71.4/79.5 ms**（`:150`）；hot 浏览器往返 p50 **3.5 ms**（`:151`）。**⚠ 同表 `:149` 还有一枚"全进程内首次 create（run 内）＝1808.9 ms（run1，紧跟 12 子进程后、竞争态）／1163.7 ms（run2）"**——这一枚才是与本票落地形态同形的（常驻进程里首次建窗，球＋音频＋DB 都在跑）。**本程未复跑（只读腿无载具，且派单禁 build/test），以上全是转引在册读数，引用前须重跑。**

**5. P11（未定义即停）这条我怎么处理。**
`AGENTS.md` §2 与 `docs/specs/SPEC-12-roadmap-governance.md:49` 逐字列着："WebView2 冷拉起 >2s → 重评 L2 卡是否回原生（P11，S0）"。按 S34 现量：**cold P95 最高 1256.4 ms、同形那枚 1808.9 ms，两者都 <2000 ms ⇒ 按现有读数未触发 P11**；但 **1808.9 > 1500 已越过 D32 的冷启线**，而它是本票落地形态的最近似读数。⇒ **我不改任何契约、不勾任何格**，只把它作为"量到/推断可能越线"的候选**摆进 ⑩ J3 交裁**（派单 `⑥` 末段的要求正是这个动作）。⚠ 另据票面 `:129-130`（`33-a1` 具名）："这条待定项在真相源台账里**没有对应 `A##`**"——**这一条我未独立复核，属转引**，落地前该先入册。

---

## ⑦ 先后与互斥：动同一枚包的必须串行

**尺**：S3（工作树脏项现量，取数 23:17）、S38（scope 名册）、S53（panel 零平台分叉）、S65（`scripts/build.ps1` 的构建形状段）、票池目录名册（`ls .scratch/wisp/issues | grep -E "^(33|35|228|244|246|248)-"` ⇒ **六枚全在、无 `-done` 后缀**，S65 同批）。

| 票/腿 | 要动的写面 | 现量凭据 | 与本票关系 |
|---|---|---|---|
| **`246-r2`（在飞）** | `cmd/wisp/`（工作树此刻 `M cmd/wisp/{main,run,resident_windows,approval_reply}.go` ＋ 1 枚 test） | S3（23:17 现量 5 枚 `M`） | **互斥**：本票落地腿要动 `cmd/wisp` 的注入点（`resident_ball_windows.go:128` 的 `OnTrayPanel` 今天只 `recordBallGesture("tray-open-panel")`，S42）⇒ **必须等 246-r2 交件**（票面 `:15` 也逐字写了"⛔ 不与 246-r2 并发（同占 `cmd/wisp`）"） |
| **本票 33（宿主落地腿）** | `go.mod`＋`go.sum`（**唯一动它们的票**，S18 根 `go.sum` 现量 jchv 0 命中）＋新包或 `cmd/wisp` 宿主文件＋`cmd/wisp` 注入点 | S18·S42·S53 | — |
| **248（设置路由）** | `internal/panel/bridge.go`（白名单四枚，`:42-45` 逐字，S41）＋处理器落点 | 票面 `:16`＋票 248 头部逐字"入向通道今天的白名单**只有四枚方法**，里面没有任何一枚是配置或凭据" | **可并行**（见下面"关键读数"） |
| **35（桥的出站）** | `internal/panel/pump.go`（`Snapshot:191`／`Marshal:273`／`Publish:284`）＋出站通道 | 票面 `:54`／`:58`；`35-panel-bridge-c17.md` 在册未 `-done` | **依赖 33**：AC#7 要的"第一枚**生产** `Marshal()` 调用者"就是宿主（票面 `:56-57` 逐字） |
| **244（GUI subsystem 构建）** | `scripts/build.ps1`（现量：`$ldflags` 在 `:103-111`、**逐字不含 `-H=windowsgui`**，构建命令 `:123`；票 244 头部自陈"今天 `build/wisp.exe` 的 PE 子系统我读到的是 CUI(3)"） | S65＋票 244 现量表 | **同域不同文件**：票面 `:15` 逐字"⛔ 不与票 244（GUI subsystem 构建）同批改"⇒ 串行；两枚都在改"产物形态/构建链" |

**关键读数（这条决定排程能不能松一档）**：**本票落地腿可以做到零 `internal/panel/` 写面**——供给复用 `panel.Assets`（S25，且 `cmd/wisp/panel_assets.go:13` 已逐字写明宿主就该调它，S62），入向复用 `ComposerDispatch.Handle`（票面 `:142`／`:86-87` 已在册），白名单那四枚一字不动（S41）。⇒ **只要落点选 ② 的甲形（新包 `internal/panel/host`）或乙形（`cmd/wisp` 文件），33 与 248 打的不是同一枚包，可以并行；两形里若是甲，`internal/panel` 的现有名册与 ubuntu 分母都零损失（S17/S38/S53）。**
⚠ 反过来说：**若有人把宿主塞进 `internal/panel` 本包（第一枚 `_windows.go`，S53 现量今天 0 枚），33 与 248 立刻互斥，且 panel 的 CI 分母当场变形**——这就是 ② 那条硬约束的排程后果。

**建议顺序与理由（理由是我的，不是仓里的规矩）**：
1. **`246-r2` 交件**（它先走完，`cmd/wisp` 才是稳定态；S3 现量它正在写）。
2. **本票 33 宿主落地腿**（go.mod/go.sum ＋ 宿主 ＋ `cmd/wisp` 注入点 ＋ `OnTrayPanel` 从"只记一笔"接到真 Show）。⇒ 排在队列头这件事票面 `:15` 已定，不是我推的。
3. **票 248（设置路由）**，与 2 并行开写可以、**合批交件不行**（若 33 选了碰 `internal/panel` 的落点，退化为串行）。理由：owner 那句"我要点击设置，自己配置模型这些参数"**要的是 2 与 3 都落地**，缺一个都点不出来（票面 `:16` 逐字"本票的验收腿**不许**因为窗口开出来了就宣称'设置可用'"）。
4. **票 244（GUI subsystem）**。理由：它改的是**产物外观**（双击带不带控制台），与功能链无耦合，放最后不动任何语义；但**它必须与 2 分开批**（票面 `:15` 的明文）。
5. **票 35 的出站**放在 2 之后（AC#7/AC#8 的 hook 由 2 产生，票面 `:56`／`:59` 逐字）。

**⛔ 本程不勾任何 AC 框**：票面 AC 框现量 `grep -c "^- \[ \]"`＝**10 格未勾**、`grep -c "^- \[x\]"`＝**0 格已勾**（S68，取数 23:32，逐枚行号 `:44 :46 :47 :48 :49 :50 :51 :59 :68 :207`＝AC#1..#6 六枚 ＋ AC#7/#8/#9/#10 四枚）。本普查是只读腿，勾与不勾都由非实现者裁（`AGENTS.md` §0.3）。
⚠ **这一条我写错过一次并当场改回**：初稿我写"9 格"，是**按记忆数的**；`grep -c` 现量是 **10**（差在 `:47/:48/:49/:50` 那四枚各算一格，我把 AC#1..#6 想成了五枚）。⇒ 定式：**报枚数必须同批发 `grep -c`，不许报脑内数。**

## ⑧ 我跑了哪些尺，每条真实读数（原样命令）

> 所有命令均在 `D:/work/workspace/projects plans/Wisp`（分支 `dev`）由 Git Bash 跑；
> `M` = `D:/work/base/gopath/pkg/mod/github.com/jchv/go-webview2@v0.0.0-20260205173254-56598839c808`。
> **禁跑闸门遵守情况：`go build`／`go vet`／`go test`／任何带 `./...` 的命令＝零枚。** `go list` 只指名单包。

| # | 原样命令 | 真实读数 | 钟点 |
|---|---|---|---|
| S1 | `date "+%Y-%m-%d %H:%M %z"` | `2026-09-30 23:16 +0800` | 23:16 |
| S2 | `git rev-parse --abbrev-ref HEAD` / `git log --oneline -3` | `dev` ／ `77a3b442`·`bb37fac2`·`457a824e`（起手 HEAD） | 23:17 |
| S3 | `git status --short \| head -40` | 脏项：`M .gitignore`、`M .scratch/wisp/issues/33-panel-host-c27.md`、`M cmd/wisp/{approval_reply,main,resident_sink_nail_127_windows_test,resident_windows,run}.go`、`D design/**`（16 枚）、`M docs/evidence/s1/152-*.md`／`246-*.md`、`?? .scratch/.scratch/`。**全部未动、未提交** | 23:17 |
| S4 | `ls`（仓根） | 有 `go.mod go.sum internal cmd scripts tools docs frontend third_party deps.toml`；**无 `vendor/`、无 `assets/`** | 23:17 |
| S5 | `cat go.mod` | 6 枚直接依赖（sherpa-onnx-go / go-toml v2 / x\*crypto / x\*sys v0.48.0 / x\*term / modernc.org/sqlite）＋10 枚 indirect；**`jchv` 零命中** | 23:17 |
| S6 | `ls -d vendor` | `No such file or directory`（退码 2）⇒ 本仓**不做 vendor** | 23:17 |
| S7 | `go env GOMODCACHE GOFLAGS GOPROXY GOSUMDB` | `D:\work\base\gopath\pkg\mod` ／ 空 ／ **`https://goproxy.cn,direct`** ／ `sum.golang.org` | 23:17 |
| S8 | `ls $GOMODCACHE/github.com/jchv` | `go-webview2@v0.0.0-20260205173254-56598839c808` ＋ `go-winloader@v0.0.0-20250406163304-c1995be93bd1`（**两枚都在本机模块缓存里**） | 23:17 |
| S9 | `ls $GOMODCACHE/cache/download/github.com/jchv/go-webview2/@v/`（同构跑 go-winloader） | 各含 `.info .mod .zip .ziphash .lock list` ⇒ **离线可完整解析，不需网络** | 23:20 |
| S10 | `find $M -type f` | 库共 **40 枚 `.go`** ＋ 3 枚 `.dll`（`webviewloader/sdk/{x64,x86,arm64}/WebView2Loader.dll`）＋ `LICENSE` ＋ `webviewloader/sdk/LICENSE.txt`（3-Clause BSD） | 23:21 |
| S11 | `grep -rln '"C"' $M --include=*.go` | **退码 1＝零命中**（无 cgo） | 23:17 |
| S12 | `grep -rln "go:cgo_import_dynamic\|#cgo" ...jchv` | **零命中**；唯一 C 源件＝`go-winloader@…/tinydll/tinydll.c`（不参编，被嵌成字节） | 23:22 |
| S13 | `grep -rn "^//go:build" $M --include=*.go` | `webview.go:1`／`pkg/edge/corewebview2.go:1`／`chromium*.go:1` 全 `windows`；`internal/w32` 按架构分 | 23:17 |
| S14 | `cat $M/go.mod` | `go 1.16`；require `go-winloader v0.0.0-20250406163304…` ＋ `golang.org/x/sys v0.0.0-20210218145245-beda7e5e158e` | 23:18 |
| S15 | `cat go-winloader@…/go.mod` | `go 1.14`；require `golang.org/x/sys v0.0.0-20200810151505-1b9f1253b3ed` ⇒ 两枚库要求的 x/sys 都**低于**本机 `v0.48.0`，MVS 不会抬高根 go.mod 的 x/sys | 23:22 |
| S16 | `go list -deps github.com/jchv/go-webview2/pkg/edge`（在 `scripts/spike/` 内，单包名） | 总 **89** 枚包；外部名 **10** 枚＝`golang.org/x/sys/windows`（已有）＋ **9 枚全新增**：`go-webview2/pkg/edge`·`go-webview2/internal/w32`·`go-webview2/webviewloader`·`go-winloader`·`go-winloader/internal/{loader,pe,vmem,memloader,winloader}` | 23:20 |
| S17 | `go list -deps ./internal/panel` / `./internal/ball`（单包名） | panel **151**（其中 `github.com/CarlosShao/wisp/*` **8**：winsec·secret·observe·risk·projctx·streamkey·frontend·panel）；ball **145**（wisp 内部 **5**：winsec·secret·observe·statemachine·ball）；两包外部非 wisp 依赖只有 `go-toml/v2` 家族 | 23:20 |
| S18 | `grep -c "jchv" go.sum` / `grep -n "jchv" scripts/spike/go.sum` | 根 `go.sum`＝**0**；spike `go.sum:13-16`＝**4 行**（两枚 module 各带 `h1:` 与 `/go.mod h1:`） | 23:20 |
| S19 | `cd scripts/spike && go list -m all` | 列出 `go-webview2 v0.0.0-20260205173254-56598839c808` ＋ `go-winloader v0.0.0-20250406163304-c1995be93bd1`，**未联网即成**（缓存命中） | 23:20 |
| S20 | `grep -rn "^func" $M/pkg/edge/chromium.go` | `NewChromium:47`·**`Embed:72`**·`Navigate:116`·`Show:151`·`Hide:155`·`EnvironmentCompleted:171`·`CreateCoreWebView2ControllerCompleted:186`·`MessageReceived:233`·`WebResourceRequested:281`·**`AddWebResourceRequestedFilter:292`**·`Environment:299`·`GetSettings:324`·`Focus:355` | 23:18 |
| S21 | `sed -n '72,114p' $M/pkg/edge/chromium.go` | `Embed` 内 **`:96-111` 自带一枚 `for { GetMessageW / TranslateMessage / DispatchMessageW }` 嵌套泵**，出口条件只有 `atomic.LoadUintptr(&e.inited)!=0`（`:97`）或 `GetMessageW` 返回 0（`:106`）；`inited` 在 `CreateCoreWebView2ControllerCompleted:224` 才置 | 23:18 |
| S22 | `grep -n "log.Fatal" $M/pkg/edge/*.go` | **`chromium.go:173` `log.Fatalf("Creating environment failed with %08x")`**、**`:188` `log.Fatalf("Creating controller failed with %08x")`**、`:284`（WebResourceRequested 取 request 失败）、`webview.go:141`（Eval 里 `log.Fatal`）；同步失败路 `:89/:92` 只是 `log.Printf`＋`return false` | 23:19 |
| S23 | `grep -n "func newI\|iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler" $M/pkg/edge/` | 该 handler 类型与构造函数**全小写未导出**（`corewebview2.go:238`、`chromium.go:23/:60`）；`ls $M/pkg/edge \| grep -i Environment`＝**零枚文件** | 23:22 |
| S24 | `grep -n "^func\|^type" $M/webview.go` | `New:87`·`NewWindow:92`·**`NewWithOptions:97`**·`CreateWithOptions:269`（自己 `RegisterClassExW` 类名 `"webview"`、自己 `CreateWindowExW`）·`Destroy:347`·**`Run:351`（又一枚 `GetMessageW` 死循环）**·`Terminate:381`·`Dispatch:443`（`PostThreadMessageW(w.mainthread,…)`）·`Bind:450` | 23:18 |
| S25 | `cat internal/panel/assets.go` | 现成供给缝：`BuiltinAssets:38` 走 `fs.Sub(frontend.Dist(), "dist")`；`EntryFile:31`＝`index.html`；`errNotBuilt:35`；`Built:57`；**`Resolve:78`**（返回 bytes＋Content-Type）；`Manifest:113`；`Check:158`；`contentTypeOf:196` | 23:18 |
| S26 | `ls -d frontend/dist` / `find frontend/dist -type f -printf "%s\t%p\n"` / `du -sb frontend/dist` | **存在**，4 枚文件：`.gitkeep` 0 B／`assets/index-BRKj5OIJ.css` 49943 B／`assets/index-BVKlegVD.js` 553469 B／`index.html` 1044 B；合计 **604456 B**。**只取名字与字节数，未读任何内容**（`frontend/**` 两层禁令） | 23:21 |
| S27 | `ls -d assets assets/web assets/web/dist` | **三枚全部 `No such file or directory`（退码 2）** | 23:18 |
| S28 | `ls -d internal/panel/assets internal/panel/assets/web` | 同样**不存在**（退码 2）⇒ 票面那条 `//go:embed assets/web/dist` 无论按仓根还是按包目录解释，今天都没有落点 | 23:22 |
| S29 | `sed -n '14,30p' .gitignore` | **`.gitignore:24` 逐字＝`assets/web/dist/`**；而 `:22-23` 是 `frontend/dist/*` ＋ `!frontend/dist/.gitkeep`（注释 `:18-21` 具名写"go:embed needs the directory to exist in a clean checkout (ticket 77 AC#1)"）⇒ **那枚路径被整条 ignore、连锚文件豁免都没有** | 23:23 |
| S30 | `grep -rn "all:dist" docs scripts internal cmd tools .github` | `docs/evidence/s1/156-…:762` 逐字 `frontend/embed.go:19 → //go:embed all:dist`；`docs/evidence/s1/85-preflight-staticcheck.md:256` 同；`docs/reports/frontend-handoff.md:24`／`frontend-handover-to-zcode.md:34` ⇒ **现用 embed 是 `frontend/embed.go` 的 `all:dist`，不是票面那条**（二手引仓内证据件，未读 frontend 一字） | 23:23 |
| S31 | `Read internal/ball/sta_windows.go`（163 行逐行） | `runtime.LockOSThread():52`／`CoInitializeEx(0, coinitApartmentThreaded):62`／泵 `:77-93`（`GetMessageW`→`TranslateMessage`→`DispatchMessageW`）／`releaseCOM→CoUninitialize:105`／**`PostTask:127`：`hwnd==0` 时 `:134-141` 只 `delete(tasks,id)` 后 `return`**／`runTask:146`／`quit():156-162`（`PostQuitMessage`） | 23:17 |
| S32 | `Read internal/ball/ball_windows.go`（1004 行逐行） | `New:143`→`Registry.Spawn("ui-sta","ball",nil,fn):170`→`b.sta(start(createOnSTA)):171`；`createOnSTA:182`（`ensureFactories:183`→`registerBallClass:187`→`CreateWindowExW:208`→`s.hwnd=…:219`）；**`uiRun:741-752`＝PostTask＋`<-done` 同步等**；`Close:935-968`（PostTask 内 DestroyWindow＋`sta.quit()`，随后 `<-b.sta.handle.Done()` join）；`fire:728` 契约注释 `:724-727`"callbacks must be quick and non-blocking" | 23:17 |
| S33 | `Read scripts/spike/webview2-latency/main.go`（328 行） | 用的是**高层** `webview2.NewWithOptions:149`（自建窗口），不是 `pkg/edge`；`pumpOnce:83-101`＝**手工 `PeekMessageW` 泵**；`:113-115` 具名坑："`w.Dispatch` 投递的是线程消息，只有 go-webview2 自己的 `Run()` 循环处理"；`runDriver:259` spawn 12 子进程测真冷 | 23:19 |
| S34 | `grep -n "webview\|P95\|cold\|hot" docs/evidence/s0/02-spike-report.md` ＋ `sed -n '144,152p'` | §3.4 两 run：cold P50/P95/max＝`879.7/1041.6/1116.8` 与 `1125.7/1256.4/(n=12)`；**全进程内首次 create＝1808.9 ms（run1，紧跟 12 子进程后、竞争态）／1163.7 ms（run2）**；hot show P50/P95＝`25.7/49.5` 与 `71.4/79.5`；hot 浏览器往返 p50＝3.5 ms；recreate P50/P95＝`955.7/1221.3` 与 `859.0/955.5`。§2.5 时序定义逐字见 `:82-89` | 23:19 |
| S35 | `sed -n '27,100p' internal/observe/goroutine.go` | `ResidentBaseline = 6`（`:39`）／`ResidentNames`（`:44-46`）＝ui-sta·audio-capture·hotkey-listener·db-writer·watchdog·log-flusher／**`TemporaryNames`（`:57-59`）里逐字含 `"panel-host"`（`:58`）**／`Spawn:262`，未入册名 `slog.Warn("goroutine outside the D38 roster (leak symptom)"):271` | 23:20 |
| S36 | `grep -rn "panel-host" internal cmd docs tools scripts` | 除 `goroutine.go:58`／`goroutine_test.go:255`（`CategoryTemporary`）外，产码**零调用方**；规格侧同名条目＝`docs/PLAN.md:2835`、`docs/specs/SPEC-01-architecture.md:122`（后者注明"**均在对应 DisposalScope 内**"） | 23:22 |
| S37 | `ls .github/workflows` ＋ `grep -n "runs-on:\|- name:" .github/workflows/ci.yml` | 两枚 workflow（`ci.yml`／`slo-fresh.yml`）；`ci.yml` 的 `runs-on`＝`ubuntu-latest`(`:66` lint、`:278` portable-core、`:665` frontend) ／ **`windows-latest`（`:388` test-windows、`:533` build＋SLO smoke）** ／ **`[self-hosted, wisp-slo]`（`:591` SLO full）** | 23:20 |
| S38 | `sed -n '120,200p' scripts/portable-tests.sh` | `core_pin`（`:125-148`）**含 `github.com/CarlosShao/wisp/internal/panel`（`:141`）**；`win_pin`（`:149-158`）＝proc·secret·config·perm·plugin·ball·risk·llmrecord ⇒ **`internal/panel` 不在 windows scope**；`core` 的 glob 里 `./internal/panel/...`（`:179`）、`windows` 的 glob 无 panel（`:185-189`） | 23:20 |
| S39 | `grep -rn "netstat\|GetExtendedTcpTable\|net.Listen" --include=*.go internal cmd tools scripts` | 非 test 命中 **2 处**：`internal/proc/treemetrics_windows.go:47`（`procGetExtendedTcpTable = modiphlpapi.NewProc("GetExtendedTcpTable")`）＋`tools/mockllm/server.go:243`（mock LLM 自己 listen，不在产码路径）。**全仓零 `netstat` 字样** | 23:20 |
| S40 | `grep -n "^func \|^type " internal/proc/treemetrics_windows.go` ＋ `sed -n '255,300p'` | 能跑的那把是**未导出**的 `treeTCPConnections(inTree map[uint32]bool) int`（`:269`）：`TCP_TABLE_OWNER_PIDAll`、AF_INET、**只数行数、不读 `state` 字段**（偏移只取 `tcpOwningPidOffset`），失败重试耗尽后 `return 0` | 23:20 |
| S41 | `sed -n '42,45p' internal/panel/bridge.go` ＋ `sed -n '1,60p' tools/d22scan/main.go` | 白名单四枚逐字（`:42-45`）＝mode/workspace/attachment/message，**零枚 config**（与票面 `:16` 一致）；d22scan 八禁：`1 bare-goroutine`（**每一枚 `go <anything>`**，唯一文件豁免＝`internal/observe/goroutine.go`）、`2 pathresolver-bypass`、`3 plaintext-key`、`4 wallclock-timeout`、`5 mirror-hash`、`6 panel-approval(frontend/)`、`7 internal-artifact-tool`、`8 emoji`（`internal/`＋`cmd/` 含注释与 `_test.go`） | 23:21 |
| S42 | `grep -rn "ball.New(" --include=*.go cmd internal` ＋ `grep -n "OnTrayPanel\|OnPanelHotkey" cmd/wisp/resident_ball_windows.go` | 常驻进程的球在 **`cmd/wisp/resident_ball_windows.go:114`** 建；`:127` `OnPanelHotkey` 与 `:128` `OnTrayPanel` **今天都只 `recordBallGesture(...)`**（后者串名 `"tray-open-panel"`），无宿主；另一枚 `cmd/balldebug/main.go:199` 逐字打印 `"tray: open panel (stub, ticket 33)"` | 23:21 |
| S43 | `grep -n "Width int\|Height int\|Scale float64\|KeepAliveInSession" internal/config/schema.go` | `[panel]`＝`PanelSection` 起 `:528`；`Width:532`（`default:"640"`）／`Height:534`（无默认＝0 auto）／`KeepAliveInSession:536`（`default:"true"`）／`Scale:538`（无默认＝0 跟随系统 DPI）；节注释 `:528` 具名 **hot-tier** | 23:21 |
| S44 | `sed -n '1,40p' deps.toml` ＋ `grep -c "license =" deps.toml` | `deps.toml` 头注逐字＝"pinned **native** dependencies (SPEC-11 §2.1)"，消费者＝`scripts/fetch-deps.ps1` ＋ `wisp doctor`；全文 `license =` **3 处**（都是 native 条目）。`ls third_party`＝`model-fetch·sherpa-onnx·spike-models`。**未找到任何"Go 模块许可登记册"**（尺＝`grep -rln "LICENSES\|licenses\|NOTICE" scripts tools docs/specs`＝零命中） | 23:21 |
| S45 | `sed -n '140,180p' docs/specs/SPEC-08*.md` | §5.1（`:143-154`）逐条：单例 PanelManager／**"jchv/go-webview2（MIT，纯 Go 无 cgo）"**／`AddWebResourceRequestedFilter` 从 `embed.FS` 喂、**不起 localhost**／前端无状态、每次 show 推 `panel.resync`／Runtime 缺失→无面板模式＋L2 原生降级卡＋**不得自动下载安装器**／多任务按 correlationId 分区。**五条里没有一条写"消息接收"**（与 `33-a1` §结论② 一致，本轮复核成立） | 23:21 |
| S46 | `ls "C:/Program Files (x86)/Microsoft/EdgeWebView/Application"` ＋ `find … -maxdepth 2 -name msedgewebview2.exe` | `153.0.4234.48`／`154.0.4258.37`／`SetupMetrics`；两枚版本目录**各含 `msedgewebview2.exe`**（路径逐字 `…/153.0.4234.48/msedgewebview2.exe`、`…/154.0.4258.37/msedgewebview2.exe`）⇒ **编排者 09-30 那枚读数我复认成立** | 23:21 |
| S47 | `sed -n '1,40p' internal/observe/clock.go` | 单调钟纪律（D42#9）逐字：超时/截止**必须**走单调钟，`observe.Timeout`（`NewTimeout:34`／`Elapsed:39`）是具名 sanctioned 机制；墙钟差值做超时常量＝**禁** | 23:21 |
| S48 | `grep -n "go func\|filepath.Clean\|Allowlist" tools/d22scan/main.go` ＋ `sed -n '1,60p'` | 豁免只有一条通道：`tools/d22scan/allowlist.txt`，格式 `ban-id<TAB>repo-relative path prefix<TAB>reason`，**从不通配**；`loadAllowlist:238`。**本程对该文件零字改动** | 23:22 |
| S49 | `sed -n '185,215p' .github/workflows/ci.yml` | staticcheck 步骤注明它跑 **`GOOS=linux`（34 枚 finding）与 `GOOS=windows`（78 枚 finding）** 两个平台形状 ⇒ **`_windows.go` 代码在 ubuntu CI 上确实被类型检查过（但不跑测试）**；该步 `if: ${{ !cancelled() }}` | 23:22 |
| S50 | `git diff -- .gitignore`（只读别人那枚脏项） | 别人正在加的只有 `.worktrees/` 一段（@@ -44,6 ＋15 行）⇒ **`assets/web/dist/`（`:24`）是 HEAD 里的既有规则，不是本轮新加的** | 23:23 |
| S51 | `wc -l internal/proc/treemetrics_windows.go` | **298 行**（⑨ 第 9 条引的就是这个数） | 23:26 |
| S52 | `find . -maxdepth 3 -type d -name assets` | 全树 3 层内只有 **`./design/old/assets`** 与 **`./frontend/dist/assets`** 两枚 ⇒ 票面 `assets/web/dist` 在任何既有目录下都没有落点 | 23:26 |
| S53 | `ls internal/panel \| grep -i "host\|windows"` ／ `find internal/panel -name "*_windows*.go" \| wc -l` | 前者**退码 1＝零命中**、后者＝**0 枚** ⇒ `internal/panel` 今天是**一枚零平台分叉包**（与 S38「它在 ubuntu 的 core scope 有分母」互为因果；②的落点问题就压在这一条上） | 23:26 |
| S61 | `grep -rn "PanelManager" --include=*.go --include=*.md internal cmd docs/specs docs/PLAN.md` | **产码零命中**；规格侧 3 处：`docs/PLAN.md:1377`（C27 契约原文，逐字含"唯一 WebView2 窗口持有者／隐藏而非销毁／面板不可用→L2 降级为原生最简确认卡，L2 能力不得消失"）、`PLAN.md:3113`（S5 切片行）、`docs/specs/SPEC-08-ui-ball-panel.md:145` | 23:26 |
| S62 | `sed -n '1,40p' cmd/wisp/panel_assets.go` | **`:13` 逐字："Ticket 33's WebView2 host is expected to call the same panel.Assets API this command prints."**；`:15` `WISP-LEG-COVERAGE-RULING: panel-assets is dispatched by main…ruled`；import 名册＝`encoding/json·flag·fmt·os·strings·internal/panel·internal/risk`（**零 net**） | 23:26 |
| S63 | `WebFetch https://pkg.go.dev/embed` | 三句原话（④第 3 点引的就是它们）："…except that files with names beginning with '.' or '_' are excluded"／"**Matches for empty directories are ignored.**"／"each pattern in a //go:embed line must match at least one file or non-empty directory. If any patterns are invalid or have invalid matches, **the build will fail**" | 23:24 |
| S64 | `grep -rn "winlive" .github/workflows scripts/*.sh scripts/*.ps1` ／ `grep -rln "go:build windows && winlive" --include=*.go .` | 前者 **零命中**（⇒ 没有任何 CI 岗位跑 winlive）；后者 **6 枚文件**：`internal/ball/{live_windows,interaction_live,hotkey_live,live_guard_windows}_test.go` ＋ `cmd/wisp/resident_ball_live_228_windows_test.go` ＋ `cmd/wisp/resident_approval_live_246_windows_test.go` | 23:27 |
| S65 | `sed -n '100,130p' scripts/build.ps1` | `$ldflags`（`:103-111`）＝6 枚 `-X …buildinfo`，**逐字不含 `-H=windowsgui`**；构建命令 `:123` `go build -trimpath -ldflags $ldflags -o build\wisp.exe ./cmd/wisp`；`:116-119` `CGO_ENABLED='1'`／`CC=$cc`（mingw）／`GOOS='windows'`／`if (-not $env:GOPROXY) { $env:GOPROXY = 'https://goproxy.cn,direct' }` | 23:27 |
| S66 | `grep -n "panel\|cold\|hot\|1500\|200" internal/observe/thresholds.go` ＋ `grep -rn "SLOPanelOpen" --include=*.go internal cmd tools` | thresholds.go 里与 panel 有关的只有 `memCapPanel int64 = 600<<20`（`:23`）／`cpuLimitPanel = 10.0`（`:30`）／`gdiLimitAll = 200`（`:33`，GDI 不是毫秒）／`case SLOPanelOpen:`（`:55`/`:72`）⇒ **零枚延迟毫秒字段**；态名在 `internal/observe/sampler.go:46 SLOPanelOpen SLOState = "PanelOpen"`、`:51` 六枚态之一。**该文件本程零字改动** | 23:27 |
| S67 | `ls .scratch/wisp/issues \| grep -E "^(33\|35\|228\|244\|246\|248)-"` | 六枚**全在、零枚带 `-done`**（`-done` 是防重领唯一键，`AGENTS.md` §1.5）⇒ 33／35／228／244／246／248 **全是活账** | 23:32 |
| S68 | `grep -c "^- \[ \]" / grep -c "^- \[x\]" .scratch/wisp/issues/33-panel-host-c27.md` | **10 ／ 0**（逐枚行号 `:44 :46 :47 :48 :49 :50 :51 :59 :68 :207`）⇒ ⑦ 末段引的就是这个数；**我初稿写的"9 格"是脑内数，已当场改回并在 ⑦ 具名自陈** | 23:32 |
| S69 | `WebFetch https://github.com/jchv/go-webview2` ／ `WebSearch "jchv/go-webview2 latest tag …"` | **两发都拿不到上游 tag**：前者 `fetch failed`（本机网络），后者只回第三方博客（CSDN/HN），**无 tag／无 releases 页** ⇒ ① 第 2 点"最新 tag 未核到"就是这两发的读数 | 23:31 |

| S70 | `ls internal/ball \| grep -c "_windows\.go$"` ／ `ls internal/ball \| grep -c "\.go$"` ／ `ls cmd/wisp \| grep -c "_windows\.go$"` | **10 ／ 30 ／ 9** ⇒ ball 与 cmd/wisp 今天都是平台分叉包，**只有 `internal/panel` 是零分叉**（S53）⇒ ② 那条硬约束的对照面 | 23:33 |

**尺数合计：70 枚**（S1–S70）。**未跑清单（具名）**：`go build`／`go vet`／`go test`／`./...` 形态／`go get`／`go mod tidy`／`GetInstalledVersion()` 实调／WebView2 真拉起／`go list -m -versions`（要出网）——前三类是派单硬禁，其余是只读腿无载具或网络不通（S69），全部登记在 ⑥⑨⑩。

---

## ⑨ 我可能写错的条目（对抗我自己，必交）

> 逐条摆"我打算这么说／哪一发的证据其实没那么硬／被推翻的代价"。**不许写"无"。**

1. **`cmd/wisp` 的依赖名册我没量。** 派单禁止因 `246-r2` 在写它而跑 `./...`，我把这条精神推广到"也不对它跑 `go list -deps`"（S3 现量：`cmd/wisp` 那 5 枚文件此刻是 `M` 态）。⇒ ②里凡是"落在 `cmd/wisp` 会多出哪几名"的算法，**是从库侧（S16）推的，不是从 `cmd/wisp` 测的**。若有人现在偷跑，读数可能是写到一半的脏集。
2. **"纯 Go 无 cgo"我是三把尺拼的，不是构建验的。** S11（`import "C"` 零命中）＋S12（`#cgo`／`go:cgo_import_dynamic` 零命中）＋S10（40 枚 `.go`／3 枚 `.dll`）。`go-winloader` 里**确有一枚 C 源文件** `tinydll/tinydll.c`（S12），但它没有对应的 `import "C"`，是**被嵌成字节**的样板。⇒ 结论"不需要 C 工具链"我按"无 cgo 指令"下，**没按"真的 `go build` 过"下**；`02-spike-report.md:211` 反而具名写 spike 全套"需 Go 1.27.1 ＋ **mingw64 gcc**"（那是 sherpa-onnx cgo 那一支要的，不是 go-webview2 要的）——**这两件事极容易被我读成一件事，请复核。**
3. **我把 `Embed` 的嵌套泵说成"会吃掉球的泵"，这一步有推断成分。** 直接证据只有 S21（`:96-111` 那段循环是库里的、`GetMessageW(hwnd=0)` 取本线程全部消息）。**"所以球仍能出帧"是我从 Win32 消息队列是线程级这一条推的，我没有跑过**（只读腿无载具，S22/S31 都没法证）。⇒ ③里凡是"实测会/不会卡"的句子我都写成"待量＋怎么量"，没写成结论。
4. **`log.Fatalf`＝进程死，我按 Go 语义说，未在本机验证。** S22 显示 `chromium.go:173/:188` 用 `log.Fatalf`。我的推断"异步失败即整进程退出、与票面 AC#5『no-crash』直接冲突"成立的前提是这两支真的被走到；**缺 runtime 时到底走同步支（`:89/:92` return false）还是异步支（`:173` Fatalf），S46 那种目录存在性证明不了**。⇒ 这一条在 ⑩ 里升格成待裁（我拿不定，且它决定逃生通道怎么写）。
5. **`GetInstalledVersion()` 能当预检我用的是源码读，不是现量。** `webviewloader/module.go` 的 `GetInstalledVersion` 在 S46 那台机器上必然返回非空串，我**没跑过**；而且它依赖 `WebView2Loader.dll` 可载（磁盘或内嵌字节，`module.go` 的 `loadFromMemory` 两路），**组策略禁用 WebView2 时返回什么我没查**。⇒ ③④里"预检再创建"的接法我标了"未实测"。
6. **"本机离线可拉"我只证到缓存里有 zip/mod/info（S9）＋`go list -m all` 不联网即成（S19）。** 我**没证**"并进根 `go.mod` 后 `go build` 也不联网"：那要 `GOFLAGS=-mod=mod` 改契约文件（严禁），且 `GOPROXY=https://goproxy.cn,direct`（S7）意味着**真缺件时会去镜像站**（见第 12 条）。
7. **`assets/web/dist` 我扫了两个解释口径，仍可能有第三个。** S27（仓根不存在）＋S28（`internal/panel/` 下不存在）。若那枚 embed  intends 落在**新包**目录（如 `internal/panel/host/assets/web/dist`），我这条尺没伸到，因为**目录根本不存在，无父可举**（**S52**：`find . -maxdepth 3 -type d -name assets` 现量＝只 `./design/old/assets` 与 `./frontend/dist/assets` 两枚；另 **S53**：`ls internal/panel | grep -i "host\|windows"` 退码 1、`find internal/panel -name "*_windows*.go"`＝**0 枚** ⇒ `internal/panel` 今天是**零平台分叉包**）。⇒ ④的结论我限定成"票面那条路径按仓根与按 `internal/panel` 两种口径都不存在"。
8. **我没读 `frontend/**` 一字，所以"现用 embed 是 `all:dist`"是二手。** 出处＝S30 那批仓内证据件（`156-…:762` 逐字写 `frontend/embed.go:19 → //go:embed all:dist`）。两层禁令我按派单走（不读、结论也不引到它身上）：④只报**文件在不在、多大、路径对不对**（S26），不报内容、不判它好坏。⚠ 但**"dist 里那两枚 hashed 文件名对得上 `index.html` 引用"这种话我说不了**——那是 `panel.Assets.Check()`（S25 `:158`）该干的事，不是普查该干的。
9. **`treeTCPConnections` 我读成"不能直接拿来当『无监听端口』的尺"，是基于它不读 `state` 字段（S40）。** 如果它其实通过别的入口暴露了 state，我就把话说重了。我只看了 `:269-300` 那一段＋枚 `^func` 清单（S40），**没通读 `treemetrics_windows.go` 全文（S51：`wc -l`＝298 行）**。⇒ ⑤里我把建议写成"新增/导出一枚读 state 的枚举"，没写成"这仓做不到"。
10. **"internal/panel 不在 windows CI scope"我只核了 `portable-tests.sh` 的 `win_pin`（S38）与 `ci.yml` 的 runs-on（S37）。** `ci.yml:533` 那枚 job 叫 "Build wisp.exe ＋ SLO smoke gate"，**它 `go build` 出来的 exe 里含 panel**，build 失败也会红。⇒ 所以准确说法是：**panel 的 windows-tagged 用例在 CI 零分母，但 panel 的 windows-tagged 代码在 CI 有两处分母（`build wisp.exe` 的链接、staticcheck 的 `GOOS=windows`，S49）**。我把这两件事说混＝⑤⑥会误导排程。
11. **延迟两数我用的是 09-19 的两轮读数（S34），未复跑，且时序口径是 spike 自己定的。** `:82-84` 的 cold 定义含 `SetHtml` ＋**首次浏览器往返**；本票改成"embed 供给＋`AddWebResourceRequestedFilter`"后，**往返里多了一趟 Go 侧读 embed 的分支**，这条 delta 我**没有数**。⇒ ⑥只交方法不交数。
12. **`GOPROXY=https://goproxy.cn,direct`（S7）与 d22scan ban #5 `mirror-hash`（S41）会不会打架，我判不动。** ban #5 的射程是"**hash material** mentioned together with a mirror"，而 `go.sum` 的 `h1:` 就是 hash、goproxy.cn 就是镜像——**但 `go.sum` 不在 d22scan 的 Go 产码射程里，且今天 S18 那 4 行 hash 就是 spike 从同一 proxy 拿的**。⇒ 我不下结论，摆给编排者（⑩第 5 条）。
13. **超预算风险自陈**：派单没给硬顶，但三枚先例（`33-a1` 42/35、`33-r2` ~50/35、`33-r3` ~47/30）都超。本程到 S70 为止的调用数我会随 commit 报，**不美化**；①-⑦的结论若在 100 轮内写不满，按派单闸门"先把⑧⑨⑩写满并 commit"办，**不拿放宽结论凑数**。
14. **我确实写错过一枚数，并已当场改回（不是假设）**：⑦ 末段初稿写"AC 框 9 格未勾"，那是**脑内数**；`grep -c "^- \[ \]"` 现量＝**10**（S68，逐枚行号也列了）。⇒ 若我没在交件前补这一发，这就是本票里第 22 枚"空输出先怀疑仪器/先怀疑自己的数"同形事故。**已具名写进 ⑦ 末尾，不静改。**
15. **① 第 2 点"最新 tag 未核到"是真的没核到，不是我看走了眼。** `WebFetch` 上游仓库页 **fetch failed**、`WebSearch` 只回第三方博客（S69）⇒ 我在 ① 里给的"库存在、纯 Go、离线可拉"三条**全部建立在仓内钉的 pseudo-version ＋ 本机 module cache 上**，**没有一条建立在上游当前状态上**。⇒ 若编排者要"今天最新 tag"这一格，得给一次能出网的岗位或允许 `go list -m -versions`；**我不拿"缓存里有"冒充"上游就长这样"。**

---

## ⑩ 判不动的地方（逐条摆甲／乙／不做＋现量再写，交编排者裁）

> 未定义即停（D22 闸门③／`AGENTS.md` §0.1）。**每条都写"现量到哪一步"，缺的都是裁，不是猜。**

**J1（⇒③）WebView2 该投在哪枚线程：`ui-sta` 共用 vs 独立 STA 线程。**
现量：票面 `:30` Key constraints 逐字"WebView2 created on the shared `ui-sta` STA thread (D38a)"；S31 `sta_windows.go:52/:62/:77` 已占死一枚 STA；S35 `panel-host` 在 **`TemporaryNames`（`goroutine.go:58`）**、S36 `SPEC-01:122` 注明"**均在对应 DisposalScope 内**"，而常驻 6 枚名册（`ResidentNames:44-46`）里没有它 ⇒ **"另开一枚常驻 STA 线程"这个名字在 D38b 名册里今天不存在**，硬用 `panel-host` 就是把"临时·在处置范围内"的名字拿去当常驻。
　甲＝按票面走 `ui-sta`（名册零变更，代价＝S21 那枚嵌套泵压在球身上）。
　乙＝新线程（要 D38b 名册变更＝**人工批准**，我不碰）。
　不做＝我自行选乙并起名（＝改契约，禁）。
**请裁：甲，还是先量一遍甲的代价再谈乙。**

**J2（⇒③/④）缺 runtime 时库走同步支还是异步支（决定 AC#5『no-crash』怎么写）。**
现量：S22 两条路都在库里（同步 `:89/:92` `return false`；异步 `:173/:188` `log.Fatalf`）；S46 本机装了运行时 ⇒ **本票在这台机器上永远走不到那支**。
　甲＝本票照 AC#5 做 fixture（改名/遮 loader），并**先在 fixture 上量到底命中哪一支**，若是 `log.Fatalf` 那支 ⇒ 具名上报"必须绕开 `Embed`、自写异步创建路径"。
　乙＝现在就判定它必是 `log.Fatalf`，直接要求自写 COM 创建（＝新增约 100 行 unsafe 产码，且 S23 证明库未导出可复用的 handler）。
　不做＝把"runtime 缺失也不崩"当已成立勾上去。
**请裁：甲（我倾向），但甲要真载具，不属本普查射程。**

**J3（⇒⑥）1808.9 ms 那一形（进程内首次 create、竞争态）算不算本票的冷启口径。**
现量：S34 cold P95 = 1041.6/1256.4 ms（均 <1500 且 <2000）；但同表"全进程内首次 create"**1808.9 ms＞1500**（仍 <2000）；而本票落地形态恰恰就是"在已经跑着球＋音频＋DB 的常驻进程里首次 create"。`AGENTS.md` §2 与 `SPEC-12 §4.2` 把"WebView2 冷拉起 >2s → 重评 L2 卡（P11，S0）"列为**未定义即停**；`33-a1` §④ 又具名报"`grep -n "P11" docs/reports/pending-and-issues.md`＝0 命中，台账无对应 `A##`"。
　甲＝把 1808.9 认作本票冷启口径的候选真值，落地腿**必须**在常驻进程内测一发，超 1500 就具名上报（不放宽 D32、不改 `thresholds.go`，S 见 ⑥ 的硬约束）。
　乙＝沿用 spike 的 1041.6/1256.4 作口径，把 1808.9 记成"竞争态离群"。
　不做＝我顺手把冷启阈值改成 2000（＝改契约，禁）。
**请裁：口径选甲还是乙；并请先给 P11 立一枚 `A##`（台账现量 0 命中这条我没复核到，转引 `33-a1`）。**

**J4（⇒②）"装配根 `cmd/wisp` 是唯一接缝"与"正向依赖边一律不开"能不能推广到宿主。**
现量：这两条是本仓既有裁定（派单 `:②` 具名提醒），但**我没有在 `PLAN.md`／`SPEC-*` 里读到它们对"原生窗口宿主"说过一个字**；`docs/evidence/s1/114-ac1-status-table.md:203` 反而把"票 114 与票 33 的地界"列成**未定需 owner 拍**。⇒ 我这一问是**我自己的推广**，不是仓里本来就有的规矩。
　甲＝宿主落**新包**（如 `internal/panel/host`），`cmd/wisp` 只注入（沿用 238 的形状，但这是我的推广）。
　乙＝宿主落 `cmd/wisp/panel_host_windows.go`（沿用 197 的形状，同样是推广）。
　不做＝我选一个并写成"本仓规矩如此"。
**请裁：甲/乙（②里两形的依赖增量我都算好了，见 ②）。**

**J5（⇒①）从 `goproxy.cn` 拿 go.sum hash，撞不撞 ban #5『mirror-hash』。**
现量：S7 那枚 GOPROXY；S18 spike 已用同 proxy 落了 4 行 hash；S41 ban #5 的原文射程＝"**hash material** mentioned together with a mirror (C29: hashes come from the signed manifest only)"，且扫描射程是 `internal/`＋`cmd/` 的 Go 产码。`deps.toml:22-26` 具名示范过镜像写法（`mirror_prefix = "https://ghfast.top/"`）并靠 **sha256 钉**兜底。
　甲＝按现状加两名 module（hash 由工具链自取），我**不**动手写任何 hash。
　乙＝先请人工确认 go.sum 是否需要与 deps.toml 同一套"官方源 hash ＋ 镜像取件"规矩。
　不做＝我为了绕 ban #5 去手写/改写任何 hash（禁，且本票不改 go.sum）。
**请裁：甲是否即可。**

**J6（⇒④）票面那条 `//go:embed assets/web/dist` 与仓里活的 `frontend/embed.go` `all:dist` 是什么关系。**
现量：S27/S28（路径两种口径都不存在）＋S29（`.gitignore:24` 逐字 `assets/web/dist/`，**整条 ignore、无锚文件豁免**）＋S30（`156-…:762` 具名 `frontend/embed.go:19 //go:embed all:dist`）＋S25（`internal/panel.Assets.Resolve` 已备好供宿主用的读接口）。Go 侧语义（`pkg.go.dev/embed`，原句）："Matches for empty directories are ignored" ＋ "each pattern in a //go:embed line must match at least one file or **non-empty directory**. If any patterns are invalid or have invalid matches, **the build will fail**"。
　甲＝宿主直接复用 `panel.Assets`（S25），票面那句按"已过期的路径写法"处理，**票面我不改一字**，由编排者裁是否留一行更正记录。
　乙＝真要新立 `assets/web/dist` ⇒ 先要动 `.gitignore:24`（那是别人正在弄脏的文件之一，S3/S50）并加锚文件，且**清检出时会硬失败**（上一条 Go 语义）。
　不做＝我改票面或改 `.gitignore`。
**请裁：甲/乙。这是"构建期会不会硬失败"的总开关。**

**J7（⇒⑤）「无监听端口」这条判据要不要进 CI，以及进了谁的scope。**
现量：S39 全仓零 `netstat`；S40 唯一能枚举 TCP 的是 `internal/proc` 里**未导出、不读 state** 的一把；S37/S38/S49 给出三种分母；S38 `internal/panel` 不在 win_pin。⇒ **"该判据只能在真机成立"的风险是实的**，而本仓已有三枚同形坑（派单 `⑤` 具名）。
　甲＝判据写成能力型（本进程 PID ＋ 树内 PID 的 TCP 表里 `state==LISTEN` 计数＝0），**并把"它跑在哪台机器"随读数具名**，跑不到的那格**登进台账**（新增 `A##` 由编排者落）。
　乙＝同时把 `internal/panel/` 追加进 `win_pin`（＝改门禁分母，`A374` 里编排者已自决"暂不加"，加要单独落 `A##`）。
　不做＝我自行往 `portable-tests.sh`/`ci.yml` 里加 scope（本程零写面）。
**请裁：甲必做；乙要不要现在做。**

**J8（⇒⑦）票 248（`internal/panel` 设置路由）与本票落地腿的先后。**
现量：S41 白名单四枚零 config（票面 `:16` 已具名"点设置录 key"归票 248）；owner 09-30 那句"我要点击设置自己配置"要的是**两枚都落地**。两枚同包 ⇒ 串行。
　甲＝先 33 宿主（窗口能开），再 248（能点设置）。
　乙＝先 248（白名单＋路由），再 33（一接上就能用）。
　不做＝并行开两枚写腿打同一包。
**请裁：甲/乙（我给的顺序建议见 ⑦，理由是"宿主是白名单的唯一真听众"，但那是我的理由不是仓里的规矩）。**
