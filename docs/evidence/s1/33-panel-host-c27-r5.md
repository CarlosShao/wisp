# `33-r5` — 面板宿主功能腿：接进常驻＋AC#13 冷启真页面＋AC#4 焦点回还＋AC#14 回执两枚钉＋退净那一支搬 winlive＋两处过期注释

- 程：`33-r5`（**写腿·功能**）｜工单＝票 33 `.scratch/wisp/issues/33-panel-host-c27.md`
- 题面来源＝票面「编排者裁定（10-01 12:12）」九条 ＋「编排者补裁（二）（10-01 13:12）」P1-P5（P1 线程形状／P2 判据形状／P3 STA 显式初始化／P4 缺口另立票 249／P5 dispatchq 不当门）
- 读数来源三件：探针 `.scratch/wisp/probes/33/p1/probe.md` §A／§B；仪器腿 `docs/evidence/s1/33-panel-host-c27-r4.md` §① 格 2／格 3／格 5 与 §⑤ 第 18／19／20 项
- 写面：`cmd/wisp/**`（产码＋测试）＋ `.scratch/wisp/probes/33/r5/**` ＋ 本证据件。⛔ 未动 `internal/ball/**`／`internal/panel/**`／`internal/proc/**`／`internal/observe/**`
- ⛔ 本件**不勾票面任何一枚 `- [ ]`／`- [x]`**（勾归编排者与非实现者验收腿 `33-v2`）；⛔ 不为变绿放宽任何断言、⛔ 不降级成 `t.Logf`

---

## 起手锚点（同发取数，一条命令里跑）

| 尺 | 现量 |
|---|---|
| `date "+%Y-%m-%d %H:%M:%S %z"` | `2026-10-01 13:22:12 +0800` |
| `git log -1 --format=%H` | `7a02b12190322cb2a0937940353e6bdfc123d766` |
| `git rev-parse --abbrev-ref HEAD` | `dev` |
| `git status --porcelain` 起手名册 | **229 行**，逐字落 `.scratch/wisp/probes/33/r5/status-start.txt`（⛔ 闸门口径＝**终态名册等于起手名册＋只我这几枚路径**，不写"必须为空"） |
| `git status --porcelain -- cmd internal` | **0 行**（起手即净） |
| 代号复核 | `ls .scratch/wisp/probes/33/` ⇒ `a1 a2 h1 p1 r1 r2 r3 r4` ⇒ **`r5` 未被占用**（起手 `ls .scratch/wisp/probes/33/r5` ⇒ 不存在） |

> ⚠ 派单那句"起手脏项约 226 行"是编排者取数的时刻值；本腿 13:22:12 现量 **229** 行。按记忆里第 77 条那味（共享树里枚数会漂）这里记两枚各带时刻，不判谁错。

---

## 六件逐件：改前读数 — 改后读数 — 反控读数

> 行号一律本腿自己 `grep -n`／`git grep` 现取（改前那批发在起手锚 `7a02b121` 上，改后那批发在当时的 HEAD 上，每条带时刻）。

### ① 面板接进常驻（生产调用者 0 → 1 枚直调 + 1 枚装配入口）

| 项 | 读数 |
|---|---|
| 改前（尺＝`git grep -n "NewPanelManager" 7a02b121 -- cmd/wisp \| grep -v _test.go`，`14:04:50`） | 命中 **2 行**＝`panel_host_windows.go:125` 注释 + `:131` 定义本身 ⇒ **非 test 调用者 0 枚**（口径：排除定义行与注释行后的调用点枚数，`A1b=0`）。这就是票面 AC#1..AC#4 全都停在"宿主可用"的原因 |
| 改后（同一把尺，HEAD `f549ecdc`，`14:04:28`） | 原始 grep **5 行**（`grep -rn "NewPanelManager" --include=*.go cmd/wisp \| grep -v _test.go`）＝3 行注释 + 1 行定义 + **1 行真调用点**（`cmd/wisp/panel_resident_windows.go:193`）。派单要的那枚"非 test 调用者枚数"＝**1**；⚠ 口径写死：原始 grep 行数（5）不是调用者枚数，本票两处读数都按"定义行与注释行不算"取 |
| 改后（AST 尺＝`TestAC1SessionDisposeHasAProductionTrigger_AC1`，同一发） | 逐字 `AC#1 dispose scan: 2 production constructor(s) [panel_resident_windows.go:193 resident_windows.go:193?]`——实际那发打的是 `[panel_resident_windows.go:193 resident_windows.go:142]`＝**直调 1 枚 + 装配根经由工厂入口 1 枚**（`newResidentPanelManager` 也注册在名册里），`0 manager-teardown site(s)` 那一支是反控读数（见下） |
| 接到哪 | `cmd/wisp/resident_windows.go:142` 建管理器 → `:148 startResidentPanel` → `:149 defer panel.stop()` → `:163 withPanelHost(...)` 注入球的两枚面板手势（`OnPanelHotkey`／`OnTrayPanel`）。⛔ 没投球的 `ui-sta`，⛔ 没动 `internal/ball/**` |
| 两档分开说（⛔ 不许读成"面板能用了"） | **宿主可用＝落地**（真窗能建、回执能到、退出有出口）。**用户能打开＝还差三件**，都不是本腿越界：ⓐ 页面侧那一腿归票 248／界面会话（`frontend/**` 两层禁令，本腿零读零写）；ⓑ 入向白名单仍只有四枚方法、**零枚 config/凭据方法** ⇒ owner 那句"点设置自己录 key"仍不通（票面 `:16` 原话）；ⓒ AC#12 那包真页面归口未落（工作树有产物、入库只 `.gitkeep`）。⚠ 我也⛔没真按过热键去证"看得见窗"（`33-p1` §⑦-1 已判桌面注入拿不到可复核前提），那一格是 owner 手测或 `33-v2` |
| 反控（M3/M3b，两发一起做） | 把 `teardown` 与 `RequestDispose` 里**那两枚对管理器的 `Destroy` 调用删掉** ⇒ ①`TestAC1SessionDisposeHasAProductionTrigger_AC1` 红，逐字 `a production site now constructs the panel host (panel_resident_windows.go:193, resident_windows.go:142) but no non-test file calls Destroy ON THE MANAGER. The 2 .Destroy() call(s) this package does have [panel_host_windows.go:252, panel_host_windows.go:465] release the WebView2 control inside the host's own file...`（＝裁定 4 要的那枚收紧有牙：库内 `w.Destroy()` 不再算会话拆窗）；②`TestPanelThreadIsSTAAndExitsCleanly` 同时红：`the panel thread exited with the window still marked created - teardown did not run on the owning thread` |

### ② 线程形状＝P1 已裁（专用 STA 线程 + 库的 `Run()`；不投 `ui-sta`）

| 项 | 读数 |
|---|---|
| 形状 | 新建 `cmd/wisp/panel_resident_windows.go`：`observe.Registry.Spawn("panel-sta", "panel host (ticket 33)", root, loop)`（ban #1 的正解＝有 owner、有 recover、有现成 panic 仪器，⛔ 不裸 `go func(`）；`loop` 逐字 `runtime.LockOSThread()` ＋ `CoInitializeEx(0, 0x2)` 并**查返回值**（`r != 0 && r != 1` ⇒ 具名拒绝建窗）；建窗只在那条线程上发生（`RequestShow` 把闭包投给那条线程，`post()` 在窗已存在时走 `webview.Dispatch`，也就是 `Run()` 唯一会取的那条队列） |
| COM 起手前（`A2`，`14:04:50`） | `git grep -n "CoInitialize" 7a02b121 -- cmd/wisp \| grep -v _test.go` ⇒ **0 命中**（复认派单那句"今天零枚"，本腿自跑，未采信转述）⇒ 出货宿主与测试 harness 跑的正是 `33-p1` 那档"未初始化" |
| 出口 | `stop()` → 把退出**投到那条线程上执行**（见下一行那枚新读数）→ `Run()` 在 `WM_QUIT` 返回 → 同一条线程上做 `teardown`（`Destroy`）→ `finished` 关闭。可观察读数＝`TestPanelThreadIsSTAAndExitsCleanly`（终态跑通，逐字日志行 `panel thread exited cleanly shows=1`＋`panel thread ending why="the library pump returned" window_opened=true`） |
| ⚠ **本腿量到的一枚新事实（比派单更具体，具名上报）** | 库的 `Terminate()` **就是裸 `PostQuitMessage`**（`webview.go:381-383` 全文只有那一行），而 `PostQuitMessage` 投的是**调用方所在线程**的队列，不是窗所属线程。⇒ 从别的 goroutine 调 `w.Terminate()` 收不掉那条线程：本腿第一版就是这么写的，`TestPanelThreadIsSTAAndExitsCleanly` 当场红 `the panel thread did not exit after stop(): Run() has no way out in this shape`（**15 秒**有界等待用尽，⛔ 不是超时放宽）。修法＝退出请求投进那条线程自己执行（`rp.post(func(){ rp.mgr.terminateOnThisThread() })`），并留 `stopRequested` 标志覆盖"还没进 `Run()` 就被要求退出"那一岔。**探针 `33-p1` R26 能收干净正是因为它在闭包里调**（同一线程），这一点它表里没写 |
| 名册／退出十步（两条停手上报的线，本腿**都没越**） | `git diff --numstat 7a02b121..HEAD -- internal/observe internal/proc internal/ball internal/panel` ＝**空**（逐字读数见 §收尾三把尺）。线程名 `panel-sta` 不进 `ResidentNames` ⇒ 落 `rep.Unknown`，实测日志逐字（本机 `13:59:52` 那发）：`level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"`＝**只吵不红**，与 `33-p1` §⑤-7 同形；钉这条边界的用例＝`TestPanelThreadNameIsNotInResidentRoster`（若有人把名字塞进名册它红）。退出侧走 `runResident` 的 **defer**（票 228 给球那一发同一形），⛔ 没新增第 11 步、⛔ 没碰那张闭集钩子名册 |
| 反控（M2，形状对照） | 把 `w.Run()` 换成手工泵（`for !stopRequested { pnlPumpOnce(); Sleep(5ms) }`）⇒ **两枚 AC#14 用例都红**（各 15s 有界等待用尽）。⚠ 这一发**没有**替我分开两维，原因具名：本腿的推送与回话都要先经 `Dispatch` 才落到那条线程，手工泵取走 `WM_APP` 却不取 `dispatchq`（＝`33-p1` R25 那一形），于是连 `Eval` 那枚也没执行。⇒ **"Eval 推送在自泵形里也到"这一维本腿没有独立复现**，它仍只有 `33-p1` R27 的读数撑着；我不拿 M2 这发去否证它，也不拿它当自己的凭据（详见 §⑤ 第 24 条） |

### ③ AC#13：冷启动那两发 `SetHtml` 不许盖掉真页面

| 项 | 读数 |
|---|---|
| 改前（行号现取于 `7a02b121`） | `panel_host_windows.go:221 serveEntry()` → `:250 w.SetHtml(入口字节)` → `:227 firstRoundTripLocked(...)` → `:374-375` 又一发 `SetHtml(自造探测页)` ⇒ **每次冷启动最终显示探测页**。派单那句时序本腿逐字复认（读 `git grep`，未复跑行为） |
| 改后 | `bringUp` 次序重排：探测在前、`serveEntry` 在后（探测**没删**，它仍是冷启"可用"判定的来源）；另加 `serveNotBuiltNoticeLocked()`——bundle 没建成时明写一枚离线提示文档，⛔ 不再把探测页留在屏上充当"页面"。`firstRoundTripLocked` 的射程在注释里写死＝**"页→Go 到达"**，⛔ 不是 AC#14 的凭据（裁定 P2 逐字） |
| 新断言（问能力） | `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`：Go 侧 `panel.Assets.Resolve(panel.EntryFile)` 取入口字节 → 正则抽 `id="..."` 当探针 → **由页面自己**用 `document.getElementById` 逐枚回答，答案经**产品已有的那道门**（`window.wispDispatch`）送回 Go。终态读数（`14:00` 那发）逐字：`AC#13 probes from the resolved entry (1044 bytes): 1 id(s) [root]` ＋ `AC#13 page answer: "1" - 1 of 1 probe id(s) present in the live document` |
| 反控 M1（派单点名的那一发） | 把两发次序换回去（`serveEntry` 在前、探测在后）⇒ **必须红**：逐字 `after a real cold start the live document contains NONE of the 1 element ids the embedded entry declares (the page itself answered "0")`。⚠ 这条红是**第二版尺**给的：第一版尺在 M1 下**仍然绿**（它自己 re-serve 了入口文档，把突变盖掉了），处置＝重写尺，见 §⑤ 第 22 条 |
| 没做的两件事 | ⛔ 没删探测（票面 `:39` 明禁）；⛔ 没上多文件资源过滤器——`AddWebResourceRequestedFilter` 在 `pkg/edge` 是导出的、`Resize()` 也导出了（⑥ 那处改口），所以"能不能满足 420×260"仍未证、这一形**留给编排者**当产品形状裁，⛔ 我不写成"已解决" |

### ④ AC#4 焦点回还那一跳（那枚故意的红 → 绿）

| 项 | 读数 |
|---|---|
| 改前（`7a02b121` 逐字复认，⛔ 不采信转述） | `Show`：`:261` 未创建先 `bringUp`（窗已建、前台已被拿走）→ `:269` 才 `m.prevFocus = windows.GetForegroundWindow()` ⇒ 记下来的"来处"就是面板自己；`Hide`：`:311` 只有 `if prev != 0`，除此之外没有任何"回还失败"路径 |
| 改后 | 采样挪到 `bringUp` **之前**（`Show` 里 `prior := windows.GetForegroundWindow()` 在第一句）；记值走 `setPriorFocusLocked`＝**0 不记、面板自己或其子窗不记**（判定用 `GetAncestor(GA_ROOT)`，比"等于面板句柄"更宽）；`Hide` 把两枚 Win32 返回值取进 `lastRestoreTo/lastRestoreSetForeground/lastRestoreSetFocus`（**只记录、不重试、不 sleep**） |
| ⚠ 一处形状改动具名（⛔ 不是放松断言） | 那枚用例的 **setup** 换了两点：(a) 窗由产品的 `Show` 建（原来用 `bringUp` 直建——产品从没那个形，而缺陷恰好是"`Show` 在建窗之后才采样"）；(b) 用例自建的**同进程 editor 窗**当"来处"，且宿主调用一律经 harness 的任务口在**窗所属线程**上执行。理由各有一条读数：非前台进程把焦点还给**别人进程**的窗会被拒（第一发逐字读数：`prevFocus 0x30176`＝终端窗、`Hide attempted restore ... 未生效`、`foreground while hidden 0x0`）；跨线程 `SetForegroundWindow` 在本机静默不生效（那一发 `after Show` 落在别窗上）。四枚断言**一字未动**，包括红句里那句 `owner 33-r2` 的过期归属（改它＝改断言，不是我这一程的权） |
| 改名 | `TestAC4FocusReturnToPriorWindowGap33r2` → `...Gap33r5`（同一枚提交里改标识符 + 注释归属；`grep -rn Gap33r2 cmd/` 终态＝**0 命中**） |
| 逐发句柄表（判"前台锁稳不稳"要的那份，同一用例连跑三发，`13:3x`，HEAD `fd6c9026` 之前的形状；形状定稿后另发三发在 §终跑名册） | 发 1：before `0x30176`｜editor `0xbb078e`｜prior `0x38b0d52`｜afterShow `0xbb078e`｜panel `0x38b0d52`｜recorded `0xbb078e`｜afterHide `0xbb078e` ⇒ **③红**；发 2/3：before 是上一发遗留句柄，③绿 ④绿。定稿三发（`13:4x`，同一条命令 `-count=3`）逐字三发同形：`foreground before any panel 0x30176 \| the ruler's own editor window 0x4e0e9c \| prior 0x4e0e9c \| after Show 0xc5f0c98 \| panel hwnd 0xc5f0c98 \| prevFocus recorded at Show 0x4e0e9c \| after Hide 0x4e0e9c \| Hide attempted restore to 0x4e0e9c (SetForegroundWindow 1, SetFocus 5115548)` ⇒ **三发全绿**（`--- PASS` x3） |
| 另立一枚可确定判定的钉 | `TestAC4PriorFocusSurvivesARefusedPanelSample`：不问前台归属，只问"面板自己当样本时不覆盖诚实的来处 + 诚实的来处活过一次 Hide"。绿；`prevFocus == editor` 那一枚还打出 `the host recorded the ruler's own editor window ... as the prior`。MUT（把 `setPriorFocusLocked` 换回无条件赋值）会当场打红它——这一发我没跑（预算），⛔ 所以我不主张它，只主张 M1/M3/M2 三发跑过的 |
| 终态那一枚用例的状态 | **绿**（整包终跑名册里逐名可核，见 §② 名册）。派单要的"本票唯一那枚故意的红"已消失 |

### ⑤ "Destroy 后 2s 退净"那一支移出默认档、进 `winlive`

- 搬了什么：`TestPanelHostRealWindowHopAndLifecycle` 末尾那一段**有界等待 + `TreeWebview > baseline` 判红**整支迁到新文件 `cmd/wisp/panel_host_windows_live_test.go`（`//go:build windows && winlive`），新用例名 `TestPanelHostWebViewChildrenExitWithinTwoSeconds_WinLive`。**上界仍是 2 秒**（`time.Now().Add(2 * time.Second)`，逐字未动），⛔ 没放宽成 5s/10s、⛔ 没加重试、⛔ 没删、⛔ 没降级成 `t.Logf`；分母仍是**本机进程树**（全机枚数只进日志）。
- 默认档继续断可确定判定的维度：同一枚 HWND 跨 hide→re-show 复用（比对身份）、单窗、`IsCreated/IsShown` 状态、**新增一枚**"窗活着 ⇒ 我们树里 `msedgewebview2 >= 1`"（红句写明"树瞎了"与"没起浏览器"两种都算红）、Destroy 后 `IsCreated` 假 + HWND 归 0。
- ⚠ 代价逐字写进文件与本判决：**`winlive` 在 CI 零岗位 ⇒ 这一支从此〔仅本机可量、CI 永看不见〕**；引用它的任何表必须带这句（撤销口令＝「33 退净断言回默认档」）。
- 本腿**没有**跑 `winlive` 档（跑它要 `-tags winlive` 再开一扇真窗，桌面预算内我把它排在默认档之后、且今天没跑）⇒ 那一支的读数只有 `33-r4` 的（我未复跑，出处＝`33-panel-host-c27-r4.md` §② 第二发终跑）。这条欠账也记在 §⑦ 第 10 条。

### ⑥ 两处过期注释

| 位置 | 改前（`7a02b121` 现取） | 改后 |
|---|---|---|
| `panel_host_windows.go:33-45` 段 | 「`PutBounds` 吃模块私有 `w32.Rect` ⇒ 没法设尺寸 ⇒ 只能换依赖／fork，所以这一形属依赖边界」 | 改口：`AddWebResourceRequestedFilter` 是导出的；`(*edge.Chromium).Resize()`（`pkg/edge/chromium_amd64.go:12`）也是导出的、高层 `webview.go:343` 在 `Embed` 成功后自己就调它 ⇒ **那不是依赖边界，是产品形状决定**。⛔ 同时明写"能不能满足 420×260 的面板尺寸**仍未证**"，没写成已解决 |
| `panel_host_windows.go:101` | 「see bringUp's `runtime.LockOSThread`」——**幻影指认**（起手 `git grep -n "LockOSThread\|UnlockOSThread" 7a02b121 -- cmd/wisp/panel_host_windows.go` ⇒ 只命中 `:101` 这行注释自己） | 改口并**补上真锁**：`PanelManager` 那段现在逐字说"这个文件里没有 LockOSThread，锁在唯一那个调用方＝`cmd/wisp/panel_resident_windows.go` 的面板线程（同时显式 STA + 把泵交给 `Run()`）"。现量：`grep -rn "LockOSThread" cmd/wisp` ⇒ 产码命中 `panel_resident_windows.go:202`（真锁）＋ 注释两行（不再指认不存在的东⻄）＋ 既有的 `notify_windows.go:138` |

---

## 门禁四数（终态，逐条真实 rc）

| 门 | 命令 | 读数 |
|---|---|---|
| build | `GOFLAGS= go build ./...` | **rc=0**（`14:03:23`，HEAD `ee5d167d`，无输出） |
| vet | `GOFLAGS= go vet ./...` | **rc=0**（同发，无输出）；另加一发 `GOFLAGS= go vet -tags winlive ./cmd/wisp/` ⇒ **rc=0**（新搬进去的那一档也编得过） |
| d22scan | `sh scripts/d22scan.sh` | **rc=0**，逐字尾行 `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=224, bans #1-5 cmd/=34, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=476, ban #8 cmd/=81`。⚠ 与 `33-p1` 终态那发（`cmd/=33`、`ban #8 cmd/=78`）差 **1 枚／3 枚＝本腿新增两枚产码文件里被扫进分母的枚数**（分母＝文件枚数口径，⛔ 不是违规数）；`internal/` 两栏**一字未动** |
| gofumpt | `gofumpt -l cmd/` | 第一发列出 3 枚（`panel_host_windows_test.go`／`panel_resident_windows.go`／`panel_resident_windows_test.go`）⇒ `-w` 之后**空列表**，单独提交 `f549ecdc style(cmd/wisp): gofumpt the 33-r5 panel host files` |
| staticcheck | — | **〔未复认〕**：本机版与 CI 钉版不同（票面 `:262`/`:296`，裁定 8），本腿未跑 |

禁区自证（同发取，逐条空输出＝一字节未动）：`git diff --numstat 7a02b121..HEAD -- internal/ docs/PLAN.md docs/specs docs/SLO.md docs/BUILD.md tools/d22scan/allowlist.txt go.mod go.sum` 的**结果与本腿的归因**落在 §收尾三把尺；三枚冻结件与 `internal/**` 全树未动；`go mod tidy`／`go get` 零次；票面 `- [ ]`／`- [x]` 零枚触碰。

---

## ⑤ 我可能写错的条目（对抗我自己）

1. **"接进常驻"可能被我自己写成"接进测试"**。判据是 `grep -rn "NewPanelManager" --include=*.go cmd/wisp | grep -v _test.go` 的**非 test 调用者枚数**。如果我只在新建的 `*_test.go` 里造它，那一枚尺仍然给 0＝本件未落地。⇒ 终态那一把尺逐字落表，⛔ 不引别人表里的数。
2. **`Run()` 形可能被我误当成"回执一定到"**。`33-p1` R26 是在**探针自己的线程**上跑 `Run()` 量到的；本腿把 `Run()` 搬进 `cmd/wisp` 之后是**新的一条线程、新的初始化路径**（`CoInitializeEx(0x2)` 显式、`Registry.Spawn` 的 owner/recover、`bringUp` 在 `Run()` 之前）。如果接完之后回执仍然不到，那是**新证据**，具名上报，⛔ 不许改判据（派单逐字）。
3. **两枚钉可能被我又并成一句**。P2 写死：`Eval` 主动推送与绑定回话是**两维**（R25 vs R27），各要各的凭据。我在代码与用例命名里必须能逐名指出哪一枚钉问哪一维；⛔ 一句 `t.Errorf` 里不许同时塞两维。
4. **AC#13 那枚断言可能是"问 SetHtml 被调用过"的伪牙**。判据是**最终文档含 embed 入口的真内容**，问能力。我的取数路径＝Go 侧 `panel.Assets.Resolve(panel.EntryFile)` 拿字节 → 从字节里抽稳定特征 → 由**页面自己**回答特征在不在（绑定报回 Go）。⛔ 不打开一枚 `frontend/**` 文件（两层禁令＝零读零写零转述；字节只经 `internal/panel` 的 Go API）。风险具名：**DOM 序列化会改写属性引号／补 `html/head/body`**，所以我抽的特征必须是**元素可寻址的东西（id）**，不是原始子串匹配；如果入口里连一枚 id 都没有，我的尺会退化——那一格我按"量不到＝失败测量"处理，⛔ 不当绿。
5. **反控可能在"另一枚口径"下不响**。`git archive HEAD` 抽的仓外副本里 `frontend/dist` 是 **gitignored**（`33-r4` §④#8 现量：入库只 1 枚 `.gitkeep`），那枚副本 `Assets.Built()=false` ⇒ AC#13 那一发只会 skip、不会红。⇒ 反控载体＝**带上工作树 dist 的仓外副本**（先例＝`33-r3` §③"产品树级复现：internal＋go.mod/go.sum＋frontend 复制到仓外"）；搬运是**字节搬运**，⛔ 不打开、不转述任何前端文件。若那一发仍不红，我就具名写"这枚反控在我的载体下量不到"，不假装量过。
6. **我把 `prevFocus` 采样挪到建窗之前，可能顺手改了那四枚断言**。派单逐字：⛔ **四枚断言一字不动**，只改标识符与注释里的归属代号。我对测试文件的改动必须逐枚能说出"这是 setup，不是断言"；断言块若被我的编辑工具碰到，我用 `git diff` 的**删除列**逐名核对并落表。
7. **⚠ 一枚必须让编排者看见的形状改动**：那枚焦点用例今天用**测试 harness**（`hostThreadHarness.bringUp`）先把窗建起来，再 `Hide`→`Show`。产品路径不是这个形（产品的 `Show` 自己建窗）。⇒ 若我把 setup 改成走**常驻线程＋产品的 `Show`**，那是让被测形状等于产品形状，**不是**放松断言；代价＝这条用例从此依赖本腿新建的产码线程。我把这一处单独记在 §六件④ 里，并保留 harness 那两枚既有用例的用法不动，免得 `33-v2` 以为我把"读产码非导出字段"那一格偷偷换成了别的观察口（裁定 5 的耦合我接受）。
8. **前台锁那一枚可能仍然红**（`33-r4` §⑤ 第 5 项：非前台进程时 `SetForegroundWindow` 可能被拒）。派单⛔ 不许加 sleep 重试、⛔ 不许降级。⇒ 若它红，我把**每一发的三枚句柄逐名落表**（before／recorded／afterShow／afterHide＋HEAD＋时刻＋整包或隔离两形）交编排者判，⛔ 不自己宣布"这是环境"。
9. **`winlive` 搬档可能被我搬成"默认档不再有退净判定"而没人知道**。那一支上界仍写 **2 秒**，⛔ 不放宽、⛔ 不删、⛔ 不加重试。代价逐字写进 §六件⑤：`winlive` 在 CI **零岗位** ⇒ 这一支从此〔仅本机可量、CI 永看不见〕。
10. **新线程的名字可能撞上契约面**。`observe.ResidentNames` 六枚名册／`ResidentBaseline` 我**一字未动**；不在名册里的名字落进 `rep.Unknown`（`goroutine.go:270-273` 那一支 `slog.Warn`＝只吵不红）。⇒ 如果哪枚既有门因为我多了一枚 unknown 而红，我**停下上报**，⛔ 不为了"看起来在册"去改名册或改线程名。
11. **退出十步是闭集**（`internal/proc/shutdown.go`／`shutdown_hooks.go`，注释逐字"冻结的是顺序"，可挂名册 8 枚）。我那枚线程的出口今天用 **defer**（`rb.stop()` 同形先例＝票 228 把球那一发挂在 `runResident` 的 defer 上，⛔ 没新增第 31 步）。⇒ 如果我发现"非挂进十步不可"，那就是要动 `internal/proc`＝禁写面 ⇒ **停手上报**，写清要动哪几行。这一格本腿判**不需要**越界，理由与读数在 §六件②。
12. **`Run()` 会占住线程 ⇒ 出货路径的出口可能不干净**。探针 R26 用的是 `w.Terminate()`（投 `WM_QUIT`，库的 `Run()` 在 `WMQuit` 分支 `return`）。我用同一形，并给"线程收干净"一枚**可观察读数**（`done` 通道＋有界等待，等待用的是 monotonic deadline，⛔ 不用墙钟时间差判超时——那是 ban #4）。⚠ 具名风险：`WM_QUIT` 是**线程消息**，如果哪天有人把面板窗挪进球的线程，那枚 `WM_QUIT` 会同时结束球的泵（`33-p1` ⑦-4 已具名这一格，我没测）。
13. **ban #1 的形状**：我新起的协程必须**有 owner ＋ 有 recover**。姿势＝`observe.Registry.Spawn(name, owner, root, fn)`（现成仪器，含 `recover`＋`debug.Stack()`），⛔ 不裸 `go func(`。测试 harness 那一枚既有 `go func()` 有注释 owner＋recover，我不动它。
14. **ban #8 会扫 `_test.go`，注释豁免、字符串字面量不豁免**。射程含 `cmd/`（`33-r4` §③ 现量它就是这样被抓的）。⇒ 本腿所有新增 Go 字符串**只用 ASCII**，⛔ 不写 `⛔`／`✓`／`≤`／任何 U+1F000-U+1FAFF。终态 `sh scripts/d22scan.sh` 必须 rc=0。
15. **`go mod tidy` 一字节不许跑**（它在 HEAD 上 exit 1，会造出不属本腿的 diff）。我用的是库里**已有**的 webview2 依赖，⛔ 零新增依赖、⛔ 不改 `go.mod`／`go.sum`。
16. **桌面**：真窗起完就收、一次一枚。⚠ 我这发会新增**会起真窗的用例**（AC#13／AC#14／线程那一族），默认档的 `cmd/wisp` 整包时间会涨；`msedgewebview2` 枚数起手／终态各量一次（`33-p1` R30 那把尺），⛔ 不给这台机器留孤儿子进程。
17. **⚠ 撞钉预检的漏计风险（记忆里第 64／70 条那一味）**。我起手跑过这些尺并把读数抄在下面 §⑥，但"今天绿的用例名册"要**整包跑一遍才算量到**；如果我在写完之前没跑起跑名册，那我可能在改语义时撞红别家的钉（本票已知的三枚：`internal/panel/composer_dispatch_test.go` 的反转钉（要求"必须有宿主"，本腿只会让它更成立）、`cmd/wisp/resident_ball_228_test.go` F1 的十枚回调名册（我只改 body、保留键名）、`cmd/wisp/panel_host_gate_test.go` 的 dispose AST 尺（裁定 4 要求本腿**同时收紧**它，见 §六件①）。
18. **⛔ 我没有做的事**（免得被读成做了）：没动 `internal/**` 一字；没翻票面任何一枚框；没跑 `go mod tidy`／`go get`；没读没写 `frontend/**`／`design/**`；没动 `PLAN.md`／`specs`／`SLO.md`／`BUILD.md`／`thresholds.go`／golden／`allowlist.txt`／三枚冻结件；没加 `Allow` 方法、没造第 5 枚 veto 通道、没造"页面点一下→Go 判成 L2 批准"的任何变体；没有把 artifacts 写入做成受门控的 Tool。
19. **可能被我写成"面板能用了"的两处夸大**。本腿落地之后仍然**不成立**的三件，我在结论里逐字保留：ⓐ 页面侧那一腿（`frontend/**`，票 248 与界面侧会话）今天仍没有人点；ⓑ 入向白名单只有四枚方法、**零枚 config/凭据方法** ⇒ owner 那句"点设置自己录 key"仍差票 248；ⓒ AC#12 的"有没有一包真页面"归口未落（工作树有产物、入库只 `.gitkeep`）。⇒ "用户能打开面板"这一句的成立条件是**上面三件之外都齐**，我在 §六件① 里按"宿主可用／用户可打开"两档分开写。
20. **一处归属卫生**：本件的行号一律本腿自己 `grep -n` 现取（⛔ 不照抄派单或前人表里的行号；引用别人读数时写明"出处＝谁的哪一节，本腿未复跑"）。
21. **⚠ 我的起跑名册被我自己污染的，而且是我自己的尺抓出来的**：起跑那发整包（`13:26:50` 收钟，`logs/baseline-start.txt`）＝**145 PASS／3 FAIL／1 SKIP**，其中一枚 `TestSecretRealBinaryRefusesValueFlag` 红句逐字 `go build ./cmd/wisp failed (is mingw on PATH?): exit status 1` ＋ 六行 `m.lastRestoreTo undefined...`＝**那一枚用例在运行时真去 `go build ./cmd/wisp`，读的是我改到一半的工作树**，不是"包里有别人的红"。⇒ 三条结论写死：① 那枚红**不是起手即在**、也不是别人造成，是本腿的中间态；② 派单给的在册名册（"故意的红 + 负载敏感那一支 + boot Ctrl+C flake"）只解释了另外两枚；③ **以后凡动 `cmd/wisp`，起跑名册要在动笔之前跑完**，否则起跑读数不可用。第二枚红 `TestPanelHostRealWindowHopAndLifecycle` 的退净那一支逐字 `our tree went baseline 0 -> now 2 (tree pids 2)`＋同发 `machine-wide 14 -> 16`＝两口径同向（复认 `33-r4` 的负载敏感判定，不是我换分母造的）。
22. **我写过一枚没牙的尺，是反控当场把它打回原形的**：AC#13 第一版为了能让页面回话，注册了一枚 ruler 绑定并**重新 serve 了一次入口文档**——MUT-A（把两发 `SetHtml` 换回旧次序）跑下去它**仍然绿**。绿的原因是尺自己把最终文档又换成了入口页，**恰好把被检的缺陷盖掉**。处置＝重写（回话改走产品已有的 `window.wispDispatch` 那道门，尺不再碰文档生命周期），MUT-A 当场红、逐字读数在 §六件③。⇒ 这一格进 §⑤ 而不是只进 commit message：**"反控不红"和"产品没缺陷"是两件事，前者只说明我的尺瞎**。
23. **一次我自己读错的形状**：`showAndWait` 第一版等的是 `IsCreated`，而 `bringUp` 在建窗那一刻就把 `created` 置真、`Show` 到末尾才置 `shown` ⇒ 两枚用例（线程那一枚与手势那一枚）读到 `created=true / IsShown=false`。那不是产品缺陷，是我的等待条件挑错了状态；改成等 `IsShown` 后同形读数消失。⚠ 反过来这条读数**有信息量**：产品的 `created` 与 `shown` 之间确实隔着探测往返与入口交付，冷启"可用"判定点落在这段里，`33-v2` 若拿 `IsCreated` 当"面板打开了"会读到半程状态。
24. **MUT-C（手工泵）没有把它声称要分开的两维分开**，原因具名：本腿的两枚 AC#14 用例都把 JS 投递走 `rp.post → webview.Dispatch`，而手工泵根本不取 `dispatchq`（＝`33-p1` R25 那一形），于是**连推送那一枚也红**。⇒ 我能主张的只有"回话那一维在手工泵下不到、在 `Run()` 下三发全到"；**"Go 主动 `Eval` 推送在自泵形里也到"这一维本腿没独立复现**，它仍只有 R27 那一发（出处＝`33-p1` §A，本腿未复跑）。两枚用例仍是两枚钉（判据形状不同、问的事实不同），但⛔ 别把 M2 读成"两维已被我这发分离"。
25. **`Terminate()` 的语义我是靠一枚红学会的，不是先读会的**：见 §六件② 那行新读数。写这行的同时我核了一遍自己有没有把它写成因果——没有：我只报"`PostQuitMessage` 投调用方队列 + 从别的 goroutine 收不掉"这两件现象，⛔ 不主张库设计好坏。
26. **焦点那一格我改的是 setup，不是断言；这句话必须能被尺复核**，否则它就是自我声明。复核法：`git diff 7a02b121..HEAD -- cmd/wisp/panel_host_windows_test.go` 里那四枚 `if` 块（含 `t.Errorf` 全文）逐字不变，变的只有函数名、注释、以及建窗/调用所在线程那几行。若 `33-v2` 读出的不是这个形状，请以它为准并把这一处当本腿的缺陷。
27. **一处可能被判"越权"的取舍**：我把 `prevFocus` 在 `Hide` 之后**保留**（不清零）。理由有读数（清了之后下一发 Hide 没有目标，实测 `after Hide == panel`）；代价是"陈旧句柄"仍可能被再次尝试——`SetForegroundWindow` 对已销毁窗返回 0，我把返回值取进 `lastRestoreTo/lastRestoreSetForeground` 而不假装成功。这一处属产品形状微裁，编排者若要相反的形状，撤销点＝`PanelManager.Hide` 的那段注释。

---

## ⑥ 我跑了哪些尺（逐条真实读数；起手与预检先登，真跑读数续编号）

| # | 命令（完整） | 逐字读数（摘要） |
|---|---|---|
| S1 | `ls .scratch/wisp/probes/33/` | `a1 a2 h1 p1 r1 r2 r3 r4` ⇒ **`r5` 不在名册**（撞名即停的预检，通过） |
| S2 | `date; git log -1 --format=%H; git rev-parse --abbrev-ref HEAD` | `2026-10-01 13:22:12 +0800`／`7a02b12190322cb2a0937940353e6bdfc123d766`／`dev`（同发，落 `logs/anchor-start.txt`） |
| S3 | `git status --porcelain \| wc -l`；`git status --porcelain -- cmd internal \| wc -l` | **229**；**0**（名册逐字＝`status-start.txt`） |
| S4 | `grep -rn "NewPanelManager" --include=*.go cmd/wisp \| grep -v _test.go`（起手那发） | 见 §六件① 的"改前读数"表（本腿自己现取，⛔ 不引票面 `:132` 那句） |
| S5 | `grep -n "CoInitialize" cmd/wisp/*.go`（排 test 另尺） | 见 §六件② 改前读数（派单那句"今天零枚"本腿复跑，⛔ 不采信转述） |
| S6 | 撞钉预检：`grep -rn "no panel host\|panel-hotkey\|tray-open-panel" --include=*.go cmd internal` | 命中 4 枚，全在 `cmd/wisp`（`resident_ball_windows.go:22/:127/:128/:203`）＋`approval_reply.go:428` 一枚注释 ⇒ **没有任何一枚断言钉死那句"no panel host"文案**；`ballGestureWhy` 那句改口不会撞词面钉 |
| S7 | 撞钉预检：`sed -n '40,60p;300,340p' cmd/wisp/resident_ball_228_test.go` | F1＝**词面型名册钉**：`ballEventCallbacks228` 十枚**回调键名**，判的是"宿主有没有把每一枚键设上"（`len(missing)`／`len(set) > len(list)`）。⇒ 本腿只改 `OnPanelHotkey`／`OnTrayPanel` 的 **body**，两枚键名保留 ⇒ 不红 |
| S8 | 撞钉预检：`grep -rn "runResident\|startResidentBall(" --include=*.go cmd/wisp` | `startResidentBall` 既有**三枚** live 用例＋一枚 `resident_approval_246_windows_test.go:314` 两枚参数调用 ⇒ 我改签名会撞编译；本腿**用变参**（`hooks ...ballHostHook`）保两枚参数形仍然合法，⛔ 不改任何既有用例的一行 |
| S9 | `grep -n "func (r \*Registry) Spawn" internal/observe/goroutine.go` | `:262 Spawn(name, owner string, root *Root, fn func(ctx context.Context)) *Handle` ⇒ 现成的 owner＋recover 形状，本腿用它是 ban #1 的正解（⛔ 不裸 `go func(`） |
| S10 | `grep -n "ResidentNames\|ResidentOverBaseline\|CategoryUnknown" internal/observe/*.go`（排 test） | `:35 CategoryUnknown`／`:43 ResidentNames`／`:65 slices.Contains`／`:81 return CategoryUnknown`／`:195 ResidentOverBaseline`／`:270 if cat == CategoryUnknown`／`:422 rep.ResidentOverBaseline = rep.Resident > ResidentBaseline` ⇒ 名册与基线我**一字未动**；新线程名落 `Unknown` |
| S11 | 必读三件全文已读 | 票面裁定九条＋补裁（二）P1-P5／`probes/33/p1/probe.md` 263 行（§A／§B）／`docs/evidence/s1/33-panel-host-c27-r4.md` 307 行。**引用它们的地方都标出处，⛔ 没把任何一枚前人读数当本腿读数** |
| S12 | `wc -l` 被测四件 | `cmd/wisp/panel_host_windows.go` **411 行**、`panel_host_windows_test.go` **723 行**、`resident_windows.go` 205 行、`resident_ball_windows.go` 271 行（全文逐行读过） |

（续编号 R1.. 留给起跑名册、终跑名册、门禁四数、反控读数。）

---

## ⑦ 判不动的地方（逐条 甲／乙／不做 ＋ 现量 ＋ 为什么判不了）

1. **球侧那一格：面板线程与球的 `ui-sta` 同时活着时，球的出帧会不会被插走**。现量：`33-p1` §⑦-结-2 逐字"我只证了 `wmAppTask` 这一类消息的顺序，D2D 的 `WM_TIMER` 出帧一枚都没测"；裁定 P1 也写明那一格〔没测〕。本腿**不往球的线程上投任何东西**（P1 的裁形），所以这一维在本腿的改动下**不新增暴露面**，但我⛔ 不能说它安全。甲＝具名交回；乙＝去动 `internal/ball/**`（禁写，归票 228）；**本腿做甲**。
2. **`Run()` 与别的泵在同一枚线程上共存**。现量：`33-p1` §⑦-结-4〔推〕。本腿的形状＝**面板线程只有 `Run()` 一枚泵**（`bringUp` 里那发 `pnlPumpOnce` 自泵在 `Run()` **之前**、同一条线程、栈不重叠），所以这一格在本腿**不需要**被裁；但"同一线程两枚泵"这一形我没有读数。⛔ 不许把它读成"已证共存无事"。
3. **MTA 那发的 `res` 真值**。现量：`33-p1` §⑦-结-3＝拿不到（要自己实现那枚未导出的 environment handler）。本腿只钉"显式 STA(`0x2`) ＋ 初始化返回值必须查"，⛔ 不写"STA 就一定成功"这种因果（P3 逐字：成因未定值、只定了现象）。
4. **`dispatchq` 只涨不取那一维**。现量：`33-p1` R25 间接读数（4 枚 `WM_APP` 被取走、0 枚闭包执行），直接长度未导出。裁定 **P5＝记欠账、不当门**，且 P1 选定 `Run()` 之后这一维由投递者本身消掉。⇒ 本腿⛔ 不为此单开闸门，只在表里说"今天的形状是 `Run()` 在取"。
5. **CI 那枚 windows cli 腿上有没有 WebView2 Runtime**。现量：未量（那是 `33-v2` 的格，票面 `:310` 逐字"判法＝读 `.github/workflows/ci.yml` 的 runner 标签＋那一步的命令"）。本腿新增的真窗用例**沿用既有默认档形状**（`//go:build windows`、起不来＝红不是 skip），所以这一格的暴露面**随本腿变大**（多几枚真窗用例）。⇒ 甲＝具名上报"我把默认档真窗用例从 N 枚涨到 M 枚"（读数在 §六件终态）；乙＝自己把它们搬进 `winlive`（⛔ 那是编排者的档级裁定，撤销口令都写好了的那一枚只有退净那一支归我）；本腿做**甲**。
6. **"焦点回还失败"该不该由产品上报**。现量：`Hide` 今天只有 `:311` 那一句 `if prev != 0`，`SetForegroundWindow` 的返回值**被丢掉**；前台锁（`33-r4` §⑤#5）会让那一跳真的失败。⇒ 修法要不要"检查返回值并响亮报出"属产品形状；本腿**加读数、不改判据**（把返回值取进日志，⛔ 不据此放宽断言④）。如果编排者要我把它升成断言，那是**新裁**，我不自己定。
7. **`Assets.Check()` 那一层的牙不在本腿射程**。现量：`33-r4` §① 格 5 逐字"本尺把半包的完备性判定委托给 `Assets.Check()`"，而 `internal/panel/**` 是本腿禁写面。⇒ AC#13 那枚"最终文档含真内容"的断言问的是**送进窗口的那包字节**，⛔ 不背书那包字节自洽。
8. **`go mod tidy -diff` 在 HEAD 上 exit 1／`staticcheck` 本机版与 CI 钉版不同**。现量：票面 `:267`／`:262`（本腿未复跑，⛔ 不跑）。⇒ 门禁四数里 staticcheck 标〔未复认〕，tidy 一字节不跑。
9. **"接进常驻"接到哪一枚触发口**。现量：`internal/ball` 的 `OnPanelHotkey`／`OnTrayPanel` 两枚回调今天只 `recordBallGesture`（`resident_ball_windows.go:127-128`），包外**没有**任何导出投递口（`33-a2` R20）。⇒ 本腿的接法＝在这两枚回调的 **body** 里投给面板线程（回调契约"快、非阻塞"＝一次 channel send／`Dispatch`，⛔ 不等窗建好）。⚠ 我没有真按过热键（那要桌面注入，`33-p1` §⑦-1 已判"做不到可复核前提"）⇒ 我这一格交的是**代码路径＋用例读数**，⛔ 不交"真按 Ctrl+Alt+P 看到了窗"。那一格留给 owner 手测或 `33-v2`。

---

## 编排者收尾标注（10-01 **14:08:44**，代提人＝编排者本人；⛔ 本节不是我代填的两节，只是把"缺什么"钉死）

**本腿死法具名**：`33-r5` 撞到 **150 轮帽**（`Agent` 回报逐字：`Reached the maximum turn limit (150). The task may be incomplete`，最后一句是 `All readings collected. Writing the six-item body and appending to ⑤⑥⑦:`）。我按死腿 intake 的尺复认：**它已经把这六件正文与 ⑤⑥⑦ 写满**（本件 167 行／39,460 字节，占位符 grep＝**0 命中**），**没写完的是「终跑名册」与「收尾三把尺」两节，以及 §⑥ 里那句"续编号 R1.. 留给起跑名册、终跑名册、门禁四数、反控读数"（`:153`）**。⇒ ⛔ **我不代填那两节**（填了就把"谁做的判"洗混，这是编排者自己的账），它们归下一腿自己取数。

**盘上状态（同发取，不是转述）**：四枚提交逐枚 `git log -1` 真身＝`fd6c9026`（骨架，⑤⑥⑦ 先写满）→ `13acad46`（产码＋用例同一发，九枚路径全在 `cmd/wisp/**`）→ `ee5d167d`（AC#13／AC#14 的尺改走产品那道门）→ `f549ecdc`（gofumpt）。`git status --porcelain -- cmd internal`＝**0 行**；`git diff --name-only 7a02b121..HEAD -- cmd/wisp/shutdown.go cmd/wisp/shutdown_hooks.go internal/observe internal/ball internal/panel internal/proc`＝**空**（⛔ 两条"停手上报"的线它都没越，我在票面上要的那两格成立）。本件那一枚未提交的正文由我代提，提交号落在这段之后的一枚 `git log` 里。

**我自己复跑的那发（补死腿没交的那一维，⛔ 不写进本件当它的读数）**：整包名册读数落在 `.scratch/wisp/probes/orchestrator/33r5-head-roster.txt`（`cmd/wisp`，`-count=1 -v`，带 sherpa DLL 那枚 PATH 形态），逐名红册归 `33-v2` 判、归我入账（台账 `A503`）。**本件的六件正文凡引用"终跑名册"处，请以我这发为准对照，与它不一致就是本腿没写完，⛔ 不是你那把尺错。**

**三件我读出来要在下一腿处理、本标注只登记不动手**：
1. `cmd/wisp/resident_ball_windows.go`（**＋67／−12**）与 `cmd/wisp/resident_windows.go`（＋32／−3）被本腿改过 ⇒ **票 228 AC#2／AC#11 那一发的题面要重锚**（它原计划的几处改动点里，手势回调那两枚 body 已被 `33-r5` 占用；两腿同包 ⇒ 串行，不许并发）。
2. 本腿新增的**真窗用例**把默认档的暴露面扩大了（它自己在 §⑦ 第 5 条具名报了这个数）⇒ 票面 `:310` 那一格（CI 那枚 windows cli 腿上有没有 WebView2 Runtime）从"我该问"变成**必须先答**：`33-v2` 用静态尺（runner 标签＋那一步命令）定案，⛔ 不许用"先推一次看看"来问。
3. §⑤ 第 27 条那处取舍（`Hide` 之后**保留** `prevFocus`、不清零）＝**产品形状微裁，归我**：我追认"保留＋把两枚 Win32 返回值取进可观察字段"这一形，⛔ 不升成断言；撤销点它已写明＝`PanelManager.Hide` 那段注释。撤销口令「**33 焦点清零改回**」。
