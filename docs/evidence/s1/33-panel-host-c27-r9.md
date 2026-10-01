# 票 33 面板宿主 —— 落地腿 `33-r9`：第五形现量 ＋ "被具名拒绝"的出口

**本腿代号**：`33-r9`（写码腿／实现者）。写面：`cmd/wisp/**`（产码两枚：`panel_host_windows.go`＋`panel_resident_windows.go`，
判据件一枚：`panel_resident_windows_test.go`）、本件、`.scratch/wisp/probes/33/r9/**`。
⛔ 票面 `.scratch/wisp/issues/33-panel-host-c27.md` 一字未动（终态现量 `grep -c "^- \[ \]"`＝**13**、`grep -c "^- \[x\]"`＝**1**，与派单给的起手数一致）。
⛔ 勾与不勾归编排者与验收腿，本表不翻任何框。

**起手锚点（同发现取，`2026-10-01 17:25:03 +0800`）**：HEAD＝`426a4d03`（`Thu Oct 1 17:20:32 2026 +0800`，
`docs(evidence): 33-v2 terminal re-run plus the leg's own probe artifacts`），分支 `dev`。
起手 `git status --porcelain`＝**272 行**（逐字存 `.scratch/wisp/probes/33/r9/00-start-git-status.txt`）＝本仓常态
（`design/**` 那批是 owner 自己的未提交件，⛔ 不还原、不提交、不删）；起手 `git status --porcelain -- cmd internal`＝**0 行**。

**这格缺陷的出处（不是本腿自造的题）**：非实现者终裁表 `docs/evidence/s1/33-panel-host-c27-v2.md`
**§5b 问①**（"拒绝之后还有没有恢复路径？⇒ 没有。今天这形是'永久拒绝'"，指到行：`panel_host_windows.go:262-264` 返回 error、
`staleCloseQueued` 是 `PM_NOREMOVE`（`:659-667`）、`panel_resident_windows.go:209-224` 的循环条件
`rp.isCreated() || rp.startUpErr() != nil` 两边都不置、`RequestShow:324-327` 只落日志）；
**§5b 问②**（"整张表里'上一任'永远是同一个人……表里没有第五形"，并写明"谁要做问① 那条恢复出口，必须先补第五形"）。
本腿交的两格＝先量第五形（§2，只读数），再据读数给出口（§3，产码）。

---

## 0. 取数口径（钉死，引用本表任何读数都带上它）

- 名册命令（派单给定那一串，逐字）：
  `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test ./cmd/wisp -count=1 -run 'Panel|AC13|AC14|AC4|BallPanel|BallGesture|ListeningSocket' -v`
- 整包命令（同一串去掉 `-run`、加 `-timeout 25m`）。红绿只认 `--- FAIL`／`--- PASS`／`--- SKIP` 行；`t.Logf` 不算失败。
- ① 的载体＝同包临时件 `cmd/wisp/zz_33r9_probe_windows_test.go`，读数取完后**改名搬出**到
  `.scratch/wisp/probes/33/r9/zz_33r9_probe_windows_test.go.moved`（⛔ 零删除；它不是判据件，⛔ 不进构建图、不进名册）。
  先例＝`33-r7` 的 `probes/33/r7/zz_33r7_probe_windows_test.go.bak`。
- 前人读数具名并注本腿是否复跑：`33-r7` 四形表（`probes/33/r7/modes-grid.txt`）本腿**未复跑**，只照它的形状补那一形；
  `33-v2` §5b 的判断本腿只引"缺第五形"这一句问题定义，其四形读数归 r7。
- `staticcheck`：**未跑**（本机版与 CI 钉版不同，跑了就是假绿）⇒ §5 标〔未复认〕。
- `go mod tidy`／`go get`：**未跑**（HEAD 上 tidy 退 1，票面 `:267`）。
- `frontend/**`／`design/**`：⛔ 未读、结论未引。`internal/**`、`shutdown*.go`、`docs/reports/**`、`docs/PLAN.md`、
  `docs/specs/**`、`docs/SLO.md`、`docs/BUILD.md`、`thresholds.go`／golden／任何阈值：**零改动**（尺在 §5）。
- 本腿在 owner 真机上开真窗（①那几发），这是预期动作；⛔ 未弹授权/提权窗、未装组件、未改系统或注册表设置。
- 桌面卫生（起手／终态各一次）：`tasklist //FI "IMAGENAME eq msedgewebview2.exe" //FO CSV | grep -c msedgewebview2`＝**12 → 12**（零增量，§5 表末；
  逐字与 `33-v2` §2.4 那发的 12 同值）。

## 1. 判据进度表（交件态：每格都有读数，⛔ 无"待填"）

| 格 | 要交的凭据 | 状态 | 读数落在哪 |
|---|---|---|---|
| ① 第五形 F1（活的 WebView2 窗＋投给它的 WM_CLOSE，跑产品自己的 bringUp） | 检查读到谁的 hwnd／拒绝有没有发生／那扇活窗后来怎样／第二次 Show | **交完 3 发** | §2.1 |
| ① 第五形 F1b（同形，载体换成普通顶层窗：检查分不清是谁的窗） | 同上＋"排干＝拆窗"的直接读数 | **交完 3 发** | §2.2 |
| ① 对照形 C（窗已 Destroy、close 已被本线程自己派发干净） | 同格读数 3 发，证两形在检查眼里可区分 | **交完 3 发** | §2.3 |
| ① 线程收摊时那扇活窗会怎样（只量） | 3 发逐字 | **交完** | §2.4 |
| ② 选型（甲／乙） | 选一条，另一条写明被否凭据 | **选甲**，乙四条凭据逐条可查 | §3.1 |
| ② 会响的钉 | 拒绝后 `post()` 立刻 false ＋ 第二次 Show 不再走同一条拒绝路 ＋ 线程沿既有 teardown 退场 | **交完 1 枚用例** | §3.3 |
| ② 定向突变 | 两发：中和新加那一处／摘掉哨兵；逐字红句 | **交完** | §3.4 |
| 名册复跑 | 派单 20 枚逐名对表＋新增枚具名＋整包两发 | **交完**（含一枚负载红，具名归因） | §4 |
| 门禁四数 | build／vet／d22scan／gofumpt 各 rc＋分母 | **交完**（staticcheck 未复认） | §5 |
| ⑤ 我可能写错的条目 | 逐条带"如果错了后果" | **交完 11 条** | §6 |
| ⑥ 我跑了哪些尺 | 逐条带命令 | **交完 18 把** | §7 |
| ⑦ 判不动的地方 | 逐条具名交回 | **交完 8 格** | §8 |

## 2. ① 第五形现量（只量读数，⛔ 不改语义、⛔ 没有顺手排干别人的消息再建窗）

`33-v2` 问② 说四形表里"上一任"永远是**已经被自己 Destroy 掉**的那扇窗；本腿量的就是缺的那一形：
**这条线程上有一扇还活着的合法窗 owner ＋ 一枚投给它的 `WM_CLOSE`**，然后跑产品自己的建窗路径。
载体：一条 `runtime.LockOSThread()` 的协程；`owner.bringUp()` 起真窗（⛔ owner 不 Destroy）；
`pumpThreadToQuiet` **只在投 close 之前**把那扇窗自己的流量派到空闲（照 `33-r7` 定式，让植入没有歧义），
**投完 close 之后、建窗之前一次都不排**；所有读数取完才做还池前的排干（那是 `33-r6/r7` 立下的"真窗钉必须自己排干净再还池"纪律）。
台件：`01-fifth-shape-1.txt`（F1＋对照形）、`03-fifth-shape-plain.txt`（F1b）、`02-thread-exit.txt`（§2.4）。

### 2.1 F1 ＝ 活的 WebView2 窗（owner 仍在世）＋投给它的 WM_CLOSE —— 三发，tid 10816／11724／23840

发 1 逐字（三发除 tid 与句柄外逐字相同）：

```
owner bringUp ok: hwnd=0x1E10E00 live=true thread windows=3 | staleCloseQueued=false hwnd=0x0
pumped 0 message(s) to idle BEFORE the plant; head=empty | staleCloseQueued=false hwnd=0x0
plant PostMessageW returned=true lastErr=The operation completed successfully.; IsWindow(0x1E10E00)=true owner.IsCreated=true head=msg 0x10 on hwnd 0x1E10E00
PRE-CREATE CHECK: staleCloseQueued=true hwnd=0x1E10E00 (owner hwnd 0x1E10E00, same=true) head=msg 0x10 on hwnd 0x1E10E00 thread windows=3
bringUp#2: err=panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0x1E10E00. ... created=false in 0ms
AFTER the refusal: IsWindow(owner 0x1E10E00)=true owner.IsCreated=true | staleCloseQueued=true hwnd=0x1E10E00 | head=msg 0x10 on hwnd 0x1E10E00 | thread windows=3 | mask=0x0
SECOND Show on the refused manager: err=panel host: refusing ... queued for hwnd 0x1E10E00. ... created=false (same road? refusal mentions WM_CLOSE=true)
THIRD Show on the refused manager: err=panel host: refusing ... created=false
Show on the LIVE owner: ok, hwnd still 0x1E10E00 live=true
released: dispatched 3, windows left=0, head=empty
RUN 1 [five] tid=10816 panicked=<nil> bringUp#2 err=panel host: refusing ... created=false elapsed=0ms
```

派单点名的四问，逐条回答（三发同读）：

1. **检查读到什么**：`staleCloseQueued=true`，hwnd 报的是**那扇还活着的合法窗自己的句柄**（`owner hwnd 0x1E10E00, same=true`）。
   ⚠ 这**不是** `33-v2` §5 M1 那种 `hwnd 0x0`（那一发队里躺的是 latched `WM_QUIT`）；本形里句柄是真的、窗是活的。
2. **拒绝有没有发生**：发生了。`created=false`、`panicked=<nil>`、`elapsed=0ms`（拒绝走在 `NewWithOptions` 之前，没建、没派、没吃）。
3. **那扇活窗后来怎样**：**还在**（`IsWindow=true`）、owner 仍以为自己在世（`owner.IsCreated=true`）、
   **那枚 close 还在队**（`staleCloseQueued=true`、`head=msg 0x10`）、线程窗枚数不动（`thread windows=3`→`3`）。
   而它自己的 `Show` 仍然成功（`Show on the LIVE owner: ok`）——被拒的不是"那扇窗还能不能用"，是"这条线程上再建一扇"。
4. **第二次 Show 的结果**：同一条拒绝路、同一个具名句柄；第三次也是（逐字见上）。
   这就是问① 那句"此后每一次 Show 都走同一条拒绝路"，本腿在第五形上把它复现了一遍。
5. 附带一枚对选型决定性的读数：**排干＝拆窗**。收尾那次 `released: dispatched 3, windows left=0`；
   F1b 那发更直白：`dispatched 7, windows left=0, head=empty, IsWindow(0x780F3C)=false`。
   ⇒ "把队列派到空"确实把**这扇合法活窗**拆掉了。问② 那句"恢复动作不许拆那扇活窗"从此不是推断，是带句柄的读数。

### 2.2 F1b ＝ 同一形、载体换成一枚普通（非 WebView2）顶层窗 —— 三发，tid 22940／26464／30488

```
plain window live: hwnd=0x780F3C IsWindow=true thread windows=3 | staleCloseQueued=false hwnd=0x0
plant PostMessageW returned=true lastErr=The operation completed successfully.; head=msg 0x31F on hwnd 0x780F3C
PRE-CREATE CHECK: staleCloseQueued=true hwnd=0x780F3C names-the-plain-window=true
bringUp on the plain-carrier thread: err=panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0x780F3C. ... created=false in 0ms (IsWindow(0x780F3C) still true, close still queued true)
released on its own thread: dispatched 7, windows left=0, head=empty, IsWindow(0x780F3C)=false
```

- **检查不在乎那是谁的窗**：一扇跟面板无关的普通窗欠一枚 close，`bringUp` 一样拒、句柄照报（`names-the-plain-window=true`，三发同）。
  ⇒ 这一枚把"拒绝"的语义从"面板上一任没排干净"扩成"**这条线程上有任何未派发的 close**"；
  ② 的出口必须对两种来源都成立，⛔ 不能靠"猜这是谁的窗"来放行。
- 一枚容易被读歪的顺带读数：队首是 `msg 0x31F`，而**过滤后的 close peek 仍然读到了那枚 close**。
  这与 `33-r6` 那句"别的消息排队时 filtered peek 看不见 quit"**不矛盾**——`WM_QUIT` 的可见性另有优先级规则，⛔ 两句不许并成一句。
- 这一形还给②的钉一枚**不产 `msedgewebview2` 子进程**的载体，所以那枚钉不去挤 AC#1 的进程树分母。

### 2.3 对照形 C ＝ 窗已 Destroy、close 已被本线程自己派发干净 —— 三发，tid 11860／25816／31912

```
owner bringUp ok: hwnd=0x6F0BF4 live=true thread windows=3 | staleCloseQueued=false hwnd=0x0
owner.Destroy() then pumped 4 to quiet (the close was DISPATCHED by the thread's own owner); IsWindow(0x6F0BF4)=false thread windows=0
PRE-CREATE CHECK: staleCloseQueued=false hwnd=0x0 (owner hwnd 0x6F0BF4, same=false) head=empty thread windows=0
bringUp#2: err=<nil> created=true in 928ms
SECOND Show on the refused manager: err=<nil> created=true (same road? refusal mentions WM_CLOSE=false)
released: dispatched 4, windows left=0, head=empty
```
三发第二次建窗 928／899／1273 ms，全部 `created=true`、`panicked=<nil>`。

**两形在检查眼里可区分**：F1 → `staleCloseQueued=true`＋`thread windows=3`＋拒绝；C → `staleCloseQueued=false`＋`head=empty`＋建窗成功。
分界是"队列里那枚 close 派发掉了没有"，⛔ 不是"窗还活着没有"——检查看不见后者。这正是问② 担心的那一维：
任何"据这张表实现的自动排干"都是在替**别人**决定哪些窗可以死。

### 2.4 线程收摊时那扇活窗会怎样（只量，三发）

```
EXIT-RUN 89197154: window 0x5510A62 owned by a locked thread that has ended: IsWindow=false (it owned 3 window(s) while alive)
EXIT-RUN 89262690: window 0x5520A62 owned by a locked thread that has ended: IsWindow=false (it owned 3 window(s) while alive)
EXIT-RUN 89328226: window 0x5530A62 owned by a locked thread that has ended: IsWindow=false (it owned 3 window(s) while alive)
```
（逐字里那串十进制就是句柄本身：`0x5510A62`＝89197154；日志把句柄同时当"第几发"打了一次，见 §6 第 8 条。）
一条 `LockOSThread` 的线程结束时（协程返回、不 unlock），**它拥有的窗跟着没**（`IsWindow=false`，3/3）。
这一枚是 §3.1 否掉**乙**的最硬的一条凭据：换线程＝先把这条线程上的窗全拆掉，包括那扇合法活窗。
⚠ 载体是一枚普通窗（问的是 Win32 的线程收摊规则）；"WebView2 的 controller 跟着线程没之后 Go 侧那半会怎样"＝**本腿未量**，见 §8 R6。

## 3. ② 把"被具名拒绝"变成有出口的（产码，最小面）

### 3.1 选 **甲**；**乙** 被否的四条凭据

**甲**＝拒绝即把该线程标成致命：`showOnThread` 认出哨兵错误 ⇒ `rp.setStartUp(err)` ⇒
`post()` 当场回 false（它本来就先看 `startUpErr()`，`:293-295`）、`loop` 沿**已有**的
`if rp.isCreated() || rp.startUpErr() != nil { break }` → `handOverPump`（取到 nil 窗）→ `teardown("no window to pump")` 退场、
`statusLine()` 把原因带出来。面板本进程内保持关闭，但**状态是诚实的**。
落点：`cmd/wisp/panel_resident_windows.go:330-374`（`RequestShow` 的闭包抽成 `showOnThread`＋那一条新分支）
与 `cmd/wisp/panel_host_windows.go:68-77,295`（哨兵＋`%w` 包装）。⛔ `staleCloseQueued`／`drainStaleQuitBeforeCreate`／检查的时机与条件一字节未动。

**乙**（允许一次 `stop()` 之后重起一条干净线程）被否，凭据逐条可在盘上查：

1. **换句柄要动第三枚产码件**。`cmd/wisp/resident_windows.go:142-147`：
   `rp, rpErr := newResidentPanelManager(...)` → `panel = startResidentPanel(rt.Registry, rp)` → `defer panel.stop()`；
   `:163-165` 的球手势闭包捕获的是那枚 `panel` 变量。`residentPanel.stop()` 被 `closeOnce` 挡着
   （`panel_resident_windows.go:367`），而 `defer panel.stop()` 在登记那一刻绑的是**旧句柄的值**
   ⇒ 原地换一枚新 `residentPanel`，旧的永远没人 stop（会话拆窗那一半漏掉，C31/31-hook 那一格）、
   新的不在 boot 的退出序里。要接对就得改 `resident_windows.go`＝**派单给的产码面只有两枚文件**，本腿不越面。
2. **原地重起要重建一次性状态**：`entered`／`finished`／`stopNow`（`:108-110`）都是一次性通道
   （`finished` 在 `:200` 被 close、`entered` 在 `:281`、`stopNow` 在 `:379`/`:391`），而 `post()`（`:293-311`）与
   `stop()`（`:394-400`）正在 select 旧通道；重建它们＝新造一枚竞态面，今天**零读数**支持它安全。
3. **§2.4 的读数**：重起前那次 `stop()` 会让这条线程收摊，而线程收摊拆掉它拥有的**全部**窗（3/3 `IsWindow=false`），
   包括第五形里那扇"owner 在世、自己的 Show 还能用"的合法窗（§2.1 第 3 条）。
   ⇒ 乙不是"不拆活窗"，是**换了扇门拆**，而且拆完立刻向用户重新承诺"面板能开"。问② 要的那句断言（恢复动作不许拆那扇活窗）它过不了。
4. **原因不可分**：§2.2 显示检查只问"这条线程有没有未派发的 close"，不在乎那是面板上一任欠的还是一扇普通窗欠的。
   甲把这句原话交给线程状态并停手；乙要在同样的信息量下自动重投一次——把"我不知道这是谁的窗"变成"再试一次"。
   另：新线程的干净**不是构造保证**而是运气——`33-r7` 的 `TestAC13BringUpSurvivesAReusedThreadQuit` 存在的理由就是"池里的线程可能带毒"，
   所以乙要么接受"重起后仍然可能被拒"，要么加一枚重试环——而重试环正是问① 里那枚"永久拒绝"的另一种写法。
5. ⛔ 第三条形（`bringUp` 自己排干/派发别人的队列）派单已禁，且 §2.1 第 5 条给了它现场读数（排干＝那扇活窗没）。

### 3.2 产码那一处（逐字）

- `panel_host_windows.go`：
  `var errPanelRefusedThread = errors.New("panel host: refusing to create the panel window")`，
  拒绝那句改为 `fmt.Errorf("%w: this thread still has an undispatched WM_CLOSE queued for hwnd 0x%X. ...", errPanelRefusedThread, closedHwnd)`
  ⇒ **组合出的字符串与 HEAD 逐字相同**（哨兵文本就是原句的头一段）。MUT2 那发把这层包装摘掉、原文字拼回，
  `33-r7` 的具名拒绝钉 `TestAC13BringUpRefusesAThreadWithAQueuedClose` 照旧 `--- PASS (1.04s)` ⇒ 文案未变是有读数的（§3.4）。
  `runtime missing` 那枚 `w == nil` 的 error **不带**哨兵：它可重试、归 AC#5，不在本格射程（§6 第 4 条把这一取舍的后果写明）。
- `panel_resident_windows.go`：`RequestShow` 的闭包抽成 `showOnThread(via string)`（成功/失败两条打印语句逐字保留），
  末尾多一支 `if errors.Is(err, errPanelRefusedThread) { rp.setStartUp(err); slog.Error("panel thread: retiring without a panel window after a named refusal", ...); fmt.Printf("wisp: panel thread will take no further requests: %v\n", err) }`。
- 两枚文件各自的头部注释补了这一格的出处与读数指认（⛔ 不是"改进检查"，是给那声拒绝登记后果）。
- ⛔ 本腿**没有**造 `lastShowErr`（编排者 16:42 第 4 格那枚欠账照旧留给下一枚动 `cmd/wisp` 的腿），
  ⛔ 没动 `showAndWait`／`post()` 的语义，⛔ 没动 `stop()`，⛔ 没新增线程、没新增泵、没新增名册名。

### 3.3 那枚会响的钉：`TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever`

位置 `cmd/wisp/panel_resident_windows_test.go:590-765`（既有判据件内，⛔ 不新增 `.go` 文件 ⇒ d22scan 分母不动）。
形状：在**生产那条面板线程**上、用**生产自己的投递口** `rp.post` 植第五形（普通窗＋`PostMessageW(WM_CLOSE)`；载体的正当性＝§2.2），
再走生产自己的门 `rp.RequestShow`，然后用**同步读数**判定：

1. `startUpErr()` 必须是那枚具名拒绝（`errors.Is(err, errPanelRefusedThread)` 且含 "WM_CLOSE"）；
2. `rp.post(noop)` **立刻** false；
3. 第二次 `rp.RequestShow(...)` **立刻** false（`time.Since` 记下答复用时；那是仪器上界，⛔ 不是阈值、不在 `thresholds.go`）；
4. `mgr.IsCreated()` 仍 false（被拒之后不许凭空多出 C27 的那一扇窗）；
5. 线程沿既有 `teardown` 退场（有界等＋具名红句，排在 1–3 之后，所以摘掉改动时先响的是 1–3 那三条具名断言）。

绿色态读数（台件 `04-nail-green.txt`，逐字，`--- PASS (0.13s)`）：

```
fifth-shape plant on the panel thread: hwnd=0xB6708D4 live=true PostMessageW(WM_CLOSE) returned=true lastErr=The operation completed successfully. | check reads queued=true names-the-live-window=true
level=ERROR msg="panel host: show failed on the panel thread" via=33r9-nail err="panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0xB6708D4. ..."
wisp: panel could not open (panel host: refusing to create the panel window: ... 0xB6708D4. ...): 33r9-nail
level=ERROR msg="panel thread: retiring without a panel window after a named refusal" via=33r9-nail err="..." shows=1
wisp: panel thread will take no further requests: panel host: refusing to create the panel window: ... 0xB6708D4. ...
level=INFO msg="panel thread ending" why="no window to pump" window_opened=false
refusal-to-retirement: settle=6.0518ms second RequestShow accepted=false in 0s | thread retired=true isFinished=true | the planted live window 0xB6708D4 IsWindow=true | startUpErr=panel host: refusing to create the panel window: ... | statusLine="panel thread could not start: panel host: refusing to create the panel window: ..."
--- PASS: TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever (0.13s)
```
读出来的形状：拒绝发生 → 落成线程状态 → `post()` 与第二次 `RequestShow` 当场回 false（`in 0s`）→ 线程沿 `"no window to pump"` 退场 →
boot 那句 `statusLine` 从此带具名原因。
⚠ 那一发的 `the planted live window IsWindow=true`：协程已返回、runtime 结束线程有先后，本用例**没对它下断言**
（"线程收摊最终拆窗"的读数是 §2.4 那三发）；这一格刻意只当读数，见 §6 第 2 条与 §8 R1。

### 3.4 定向突变两发（⛔ 零枚突变体进提交；跑完逐枚还原）

统一流程：`md5sum` 取起手 → 改坏 → `grep -c MUT` 证落地 → 指名用例单跑 → 用台件里的 `.pristine` 副本还原 → 再取 `md5sum` 逐字相同 → `grep -c MUT` 归零。

| # | 突变（改坏的那一处） | 指名用例 | 逐字红句 |
|---|---|---|---|
| **M1** | `panel_resident_windows.go:367` 那支改成 `if false && errors.Is(...) {`（中和本腿新加的出口） | `TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever` | `panel_resident_windows_test.go:715: after a show request refused by name, startUpErr() = <nil>: the named refusal left no consequence on the panel thread. That is ticket 33 / 33-v2 问① - the thread stays up taking every later request and refusing each one, so the process can never open a panel and never says so outside a log line. Fix: showOnThread must recognise the sentinel bringUp returns (errPanelRefusedThread) and mark the thread fatal`；同发 `:724: post() still accepts work on a thread that has refused this shape by name - the caller's 15-second wait is the only thing that ever reports the failure (33-v2 §9.1), which is the shape this case exists to end`；`:733: a second show request was accepted onto a thread that already refused this shape by name - it will be refused again, which is the permanent-refusal road; RequestShow must answer false once the refusal is recorded`；`:756: the panel thread that refused by name is still up 5s later: it holds a queue whose only message is somebody else's WM_CLOSE and a task port nobody should be posting to` ⇒ `--- FAIL: TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever (11.09s)`／`FAIL github.com/CarlosShao/wisp/cmd/wisp 11.178s`（台件 `05-mutation-1-no-outlet.txt`） |
| **M2** | `panel_host_windows.go:295` 把 `%w` 摘掉、把原句文字逐字拼回（哨兵不再进错误链） | 同一枚钉 ＋ `33-r7` 那枚具名拒绝钉 | `--- PASS: TestAC13BringUpRefusesAThreadWithAQueuedClose (1.04s)`（拒绝文案钉照旧绿＝文字未变）同发 `--- FAIL: TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever (11.09s)`，红句与 M1 的前三条逐字相同（台件 `06-mutation-2-unwrapped-refusal.txt`）⇒ 哨兵这一半**承重** |

⚠ **红因不是 15 秒超时**：M1/M2 先响的是 `:715`／`:724`／`:733` 三条**具名**断言（都在线程还在世时读的同步断言）；
`:756` 那条收摊断言排在它们之后，用的是本包既有的 5 秒档（`refusalSettleBudget` 与 `panelThreadExitBudget` 同值，本腿不新造时间阈值）。
M1 那发的 `settle=5.0031717s` 是"没置状态所以轮询到上界"的读数，红句本身点名的是缺陷（`startUpErr() = <nil>`），不是"超时了"。
⚠ 另外 M1/M2 读数里能看到 `via=33r9-nail-second` 那一行**第二次拒绝**的日志——突变体下"每次 Show 都被拒"这件事被现场复现了一次。

还原自证（终态）：`panel_host_windows.go` 还原后 `md5 3008f4a552ede10b6b3597b1e8c667af`＝起手值；
`panel_resident_windows.go` 还原后 `md5 b2e0bee9de472cd3d8024fd49d72d0a3`＝起手值；
`grep -c "MUT-33r9\|MUT2-33r9" cmd/wisp/*.go` **非零者＝零枚**；
`git diff --numstat -- cmd internal`＝**恰三枚路径**（下表 §5）。
⚠ 一枚顺序披露：`gofumpt -w` 跑在突变还原**之后**，它把 `panel_resident_windows.go` 与判据件各重排了一次
（前者从 `b2e0bee9…` → 终态 `8bc5c0de27c8fd074c949318d5694721`，后者 → `46403f59df121858f444fa868d4e2e60`），
`panel_host_windows.go` 未被重排（终态＝起手值）。三枚终态副本另存为 `.post-gofumpt`，⛔ 突变期的 `.pristine` 副本留着不改名。

## 4. 名册复跑

### 4.1 派单给的那一串（同命令、同 `-v`），终态

`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test ./cmd/wisp -count=1 -run 'Panel|AC13|AC14|AC4|BallPanel|BallGesture|ListeningSocket' -v`
逐字末行 `ok  	github.com/CarlosShao/wisp/cmd/wisp	10.005s`，`rc=0`；
`grep -c '^--- PASS'`＝**21**、`grep -c '^--- FAIL\|^--- SKIP'`＝**0**。台件 `07-roster-1.txt`／名册 `07-roster-names.txt`。

派单那 20 枚**逐名对表**（右列＝本发结果，全部 PASS）：

| # | 用例名 | 本发 |
|---|---|---|
| 1 | TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject | PASS |
| 2 | TestNativeHostSeamRefusesAPanelSourcedAllow | PASS |
| 3 | TestTicket223PanelInboundSaysHotReloadIsDisabled | PASS |
| 4 | TestAC4EveryLegIsNailedOrRuled | PASS |
| 5 | TestPanelHostOpensNoListeningSocketL1 | PASS |
| 6 | TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12 | PASS |
| 7 | TestPanelHostRealWindowHopAndLifecycle | PASS |
| 8 | TestAC4FocusReturnToPriorWindowGap33r5 | PASS |
| 9 | TestPanelHostLatencyPercentilesAC2 | PASS |
| 10 | TestAC3ListeningSocketRulerSeesItsOwnListener | PASS |
| 11 | TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe | PASS |
| 12 | TestAC13BringUpSurvivesAReusedThreadQuit | PASS |
| 13 | TestAC13BringUpRefusesAThreadWithAQueuedClose | PASS |
| 14 | TestAC14AwaitedBindingReplyReachesThePage | PASS |
| 15 | TestAC14GoSideEvalPushReachesThePage | PASS |
| 16 | TestPanelThreadIsSTAAndExitsCleanly | PASS |
| 17 | TestPanelThreadNameIsNotInResidentRoster | PASS |
| 18 | TestBallPanelGesturesReachThePanelThread | PASS |
| 19 | TestBallGestureWithoutPanelHostStillRecords | PASS |
| 20 | TestAC4PriorFocusSurvivesARefusedPanelSample | PASS |

**＋1 枚新增（本腿的钉）**：`TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever`（`--- PASS (0.06s)`）。
⇒ 差集＝**＋1／−0**：没有用例静悄悄消失，也没有别的枚被带进来。⛔ 没有一枚从绿变红。

### 4.2 整包两发（同 `33-v2` §2.2 的口径：`-count=1 -v -timeout 25m`），其中一枚具名负载红

| 尺 | 发 1（`09-fullpack-1.txt`，`17:53:xx`→`17:56:5x`，HEAD `28f7ba8c`） | 发 2（`11-fullpack-2.txt`，`18:0x`→`18:0x`，同码） |
|---|---|---|
| 末行 | `FAIL github.com/CarlosShao/wisp/cmd/wisp 184.620s`，`rc=1` | （见下：本腿落笔时该发仍在跑，跑完把末行逐字补进这一格——⛔ 不拿发 1 冒充两发、⛔ 不拿名册那一发冒充整包） |
| `^--- PASS`（只顶层） | **159** | — |
| `^--- FAIL`／`^--- SKIP` | **1／0** | — |
| 含子项 `--- PASS` | **240** | — |
| 逐名差集 vs `33-v2` §2.3 那份 `top-final.txt`（159 名） | **＋1 枚**（本腿的钉 `TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever`）／**−0 枚**；另 **1 枚由绿转红**（下一行那枚） | — |

**那一枚红的具名归因**（⛔ 不压、不降级、不加重试、不改判据）：

```
--- FAIL: TestAC4FocusReturnToPriorWindowGap33r5 (0.50s)
    panel_host_windows_test.go:753: AC#4 focus hook (head 28f7ba8c): foreground before any panel 0x30176 | the ruler's own editor window 0xec067e |
      foreground while hidden (prior) 0x30176 | after Show 0x5ee0978 | panel hwnd 0x5ee0978 | prevFocus recorded at Show 0x30176 | after Hide 0x0 | Hide attempted restore to 0x30176 (...)
    panel_host_windows_test.go:765: Hide did not hand focus back to the window that held it before Show: expected 0x30176, foreground is 0x0
```
- **它走的代码路径与本腿的 diff 零交集**：该用例在 `cmd/wisp/panel_host_windows_test.go`，
  现量 `grep -c "residentPanel\|RequestShow\|showOnThread\|errPanelRefusedThread" cmd/wisp/panel_host_windows_test.go`＝**0**——
  它用 `hostThreadHarness` 直接调 `mgr.Show/Hide`，本腿一句都没改（本腿改的是常驻腿的投递与那句错误的身份）。
- **顺序上也轮不到本腿的钉毒它**：整包逐字行号——`=== RUN TestAC4FocusReturnToPriorWindowGap33r5` 在**第 538 行**，
  `=== RUN TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever` 在**第 614 行**（文件名字序使然），本腿的钉在它**之后**才跑。
- **当场复认（隔离复跑）**：`GOFLAGS= go test ./cmd/wisp -count=2 -v -run 'TestAC4FocusReturnToPriorWindowGap33r5'`
  ⇒ **两枚全 PASS**（台件 `10-ac4-isolated.txt`）。第二发逐字 `prevFocus recorded at Show 0x3a50aa8 | after Hide 0x3a50aa8`（回还成功），
  第一发逐字 `foreground before any panel 0x30176 | ... prevFocus recorded at Show 0x30176 | after Hide 0x30176 | SetForegroundWindow 1, SetFocus 0`
  ——那枚 `0x30176` **不是本进程的窗**（本腿的编辑窗是 `0x13b0782`/`0x3a50aa8`），即"来处"采到了一枚别家进程的窗口句柄。
- **判读（本腿只到"这一形可复现／不可复现"这一层）**：这与 `33-v2` §1 AC#4 那句 M4 结论同源——
  这一族断言吃的是 **Windows 的前台锁**（`SetForegroundWindow` 到一枚别家进程的窗会被拒），
  也与编排者 10-01 12:12 裁定 2 第③件里那句"§⑤ 第 5 条那枚'Windows 前台锁'不稳"对上；
  本发整包跑到这一段时机器上正并发着 `TestCleanCheckoutBuilds_AC11` 的 13.28 秒子构建（逐字在第 530 行）＝负载态。
  ⇒ 具名登记为**负载／前台锁敏感的既有红**，⛔ 不由本腿判它算不算 AC#4 的产码缺陷（§8 R7 交回编排者与 `33-v3`），
  ⛔ 本腿没碰它一字、没有为它改任何门禁。

## 5. 门禁读数（终态，逐条带命令；台件 `.scratch/wisp/probes/33/r9/08-gates.txt`）

| 门 | 命令（逐字） | 读数 |
|---|---|---|
| 构建 | `GOFLAGS= go build ./...` | **rc=0** |
| vet | `GOFLAGS= go vet ./cmd/wisp` | **rc=0**（第一发吃了两枚本腿自己写错的：`t.Errorf` 双 `%X` 单参 ＋ 非常量格式串的 `fmt.Sprintf`；改完复跑归零，逐字红句在 §6 第 6 条旁边——⛔ 没为变绿改任何断言） |
| d22scan | `sh scripts/d22scan.sh` | **rc=0 clean**；分母逐字与 `33-v2` §2.4 同：bans #1-5 `internal/=224` `cmd/=34`；ban #6 `frontend/=85`；ban #7 `internal/tools/=23`；ban #8 `design/=39` `frontend/=85` `internal/=476` `cmd/=81`（⚠ 分母＝射程内文件枚数，本腿只改既有文件、**没新增 `.go`** ⇒ `cmd/=81` 不动＝本腿没把任何探针留在树里）；另 `d22scan: skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/]` |
| gofumpt | `"$(go env GOPATH)/bin/gofumpt.exe" -l cmd/wisp/panel_host_windows.go cmd/wisp/panel_resident_windows.go cmd/wisp/panel_resident_windows_test.go` | 第一发**列出 2 枚**（`panel_resident_windows.go`、`panel_resident_windows_test.go`）⇒ 对这三枚名下件跑 `-w` ⇒ 复列**空输出**；⛔ 未对整目录跑 `-w`（不碰别人的件） |
| staticcheck | 未跑 | 〔未复认〕：本机版与 CI 钉版不同版，跑了拿结论就是假绿（编排者 10-01 12:12 裁定 8 同口径） |
| `go mod tidy`／`go get` | **未跑** | 派单禁（HEAD 上 tidy 退 1，票面 `:267`） |
| 桌面卫生 | `tasklist //FI "IMAGENAME eq msedgewebview2.exe" //FO CSV \| grep -c msedgewebview2` | 终态 **12**（＝`33-v2` 那发的 12，零增量）；`wisp.test.exe`／`go.exe` 在整包跑完后**各 0 枚**（本腿没留孤儿子进程；跑中量的非零是本腿自己那发，见 §6 第 9 条） |

**写面与禁区自证（终态现量）**
- `git diff --numstat HEAD -- cmd internal`＝**恰三枚路径**：
  `cmd/wisp/panel_host_windows.go 33/1`、`cmd/wisp/panel_resident_windows.go 49/7`、`cmd/wisp/panel_resident_windows_test.go 180/0`。
  ⛔ 零枚删除、零枚改名、零枚 `internal/**`；那 1/7 枚删除列全是本腿自己替换掉的行（`RequestShow` 的旧闭包与那句 return）。
- `git diff --numstat HEAD -- .scratch/wisp/issues docs/reports docs/PLAN.md docs/specs internal`＝**空**（票面 13/1 框一枚未碰）。
- 终态 `git status --porcelain -- cmd internal`＝上面那三枚（起手为 0 行）＝**起手名册＋只本腿名下路径**；
  全仓起手 272 行的其余各项（`design/**` 那批、别人名下的 `??`）本腿未动、未提交、未删。
- 本腿名下的另两枚写面：`docs/evidence/s1/33-panel-host-c27-r9.md`（本件）＋ `.scratch/wisp/probes/33/r9/**`
  （18 枚台件：起手名册、四份读数日志、两份突变、两份还原副本、`.moved` 探针载体、gofumpt 后三份副本等）。
- ⛔ 零枚 `t.Skip` 新增、零枚既有用例删除／改名／降级、零处超时放宽；唯一的新时间常量是本用例的仪器上界
  `refusalSettleBudget = 5 * time.Second`（与既有 `panelThreadExitBudget` 同值，不在任何产品阈值文件里）。
- Git：只 commit、⛔ 未 push；每枚 commit 带**显式 pathspec**；⛔ 无 `add -A`／`add .`／`commit -a`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；仓内零 worktree、零 `rm`。

## 6. ⑤ 我可能写错的条目（逐条带"如果错了后果"）

1. **"生产路径今天到不了这条拒绝"我只核了一半。** 复认了 `grep -rln 'pumpThreadToQuiet' cmd/wisp` 只命中测试件、
   `resident_windows.go:163-165` 只这一条手势通道；⛔ **没有**逐名量过 `rp.mgr.Show`／`HotShow` 的全部生产调用者。
   如果错了（还有别的路径在这条线程上调 Show）⇒ 后果＝甲的致命覆盖面比我写下的更大，AC#5 那一形（运行库缺失、应可重试）可能被牵连。
2. **甲会让线程收摊，而收摊拆掉那条线程上的活窗（§2.4 三发读数）。** 我把这件事写成读数、没写成判据。
   如果编排者认定"任何形状都不许拆那扇合法活窗"⇒ 后果＝甲也要返工；唯一不拆窗的出口是"置致命但线程 park 不退"，
   那是新线程形状，本腿没造（⇒ §8 R1）。
3. **`refusalSettleBudget = 5 * time.Second` 是本腿新造的仪器上界**，且第 3 条断言里我拿它比了一次答复耗时。
   如果负载把 `in 0s` 抬到 5 s 以上 ⇒ 后果＝多一枚与 AC#1 同族的**负载敏感断言**，`33-v2` §6 那份"真开窗户在 CI"的名册多一行会抖。
   （绿色态读数：`settle=6.0518ms`、`second RequestShow accepted=false in 0s`。）
4. **`errors.Is(err, errPanelRefusedThread)` 只覆盖"具名拒绝"，其它 show 失败照旧不落状态。**
   如果编排者要的是"任何 show 失败都致命"⇒ 后果＝运行库缺失那一形仍是"点了没反应＋只有日志"，本格只算修了一半。
   我按票面 AC#5（那一形要的是降级事件＋L2 标志，`33-v2` 判它不成立）判断它不该并进本格。
5. **我借用了 `startUp` 这枚字段承载"运行中被具名拒绝"，而它的名字是"启动失败"。**
   如果读代码／读日志的人按字面读 ⇒ 后果＝`statusLine()` 那句 `panel thread could not start: ...`
   会让人以为这条线程从没起过，实际是"起过、拒了一声后退场"。改那句措辞超出最小面（且票 246 有枚盯措辞的用例），故只登记（§8 R2）。
6. **"哨兵包装后文案逐字不变"我只用一发突变（M2）＋ `33-r7` 那枚 substring 钉的照旧绿来证。**
   ⛔ 没做逐字节 diff 比对旧句/新句 ⇒ 这一条是〔推＋一发读数〕不是〔量〕。如果错了 ⇒ 后果＝别处逐字引用整句的表/日志对不上。
7. **① 的三形都在同一条锁死线程上量，"合法 owner"我用的是同进程内的另一枚 `PanelManager`／一枚普通窗。**
   如果问② 要的"活 owner"必须是生产里可能出现的形状 ⇒ 后果＝第五形的定义被我读宽或读窄，下一程要重造载体。
   我把 F1（面板窗）与 F1b（普通窗）分开列了，没混成一句。
8. **`02-thread-exit.txt` 那三发的日志把句柄同时当"第几发"打了。** 如果读者把 `EXIT-RUN 89197154` 读成发号
   ⇒ 后果＝误以为跑了八千万发；实际那串就是句柄 `0x5510A62` 的十进制。文本缺陷，不是读数缺陷。
9. **整包两发的负载态不由我保证。** 发 1 里 `TestCleanCheckoutBuilds_AC11` 自带 13.28 秒子构建、
   且本腿自己在同一台机器上连着跑了突变与名册；`tasklist` 在跑中量到 `go.exe`/`wisp.test.exe` 非零＝本腿自己那发。
   如果终态读数被自己的负载污染 ⇒ 后果＝把一枚负载红当成交件红（§4.2 那枚就是这样处理的：先隔离复跑再归因）。
10. **"乙需要第三枚产码件／一次性通道"这枚凭据我只核到 `resident_windows.go:142-147`＋`closeOnce`＋三枚通道。**
    如果错了（存在一种安全的原地重起）⇒ 后果＝本腿否掉了编排者本来想要的出口，选型要重议。
    这是我对乙**最硬**的反对，也是最该被 `33-v3` 复认的一句（§8 R5）。
11. **我把 `RequestShow` 的闭包抽成 `showOnThread` 并改了返回值形状。** 爆炸面＝走 `showAndWait`／`RequestShow`／`RequestToggle` 的那几枚
    （名册 §4.1 逐名绿）；⛔ 但我**没**逐字比对常驻 boot 的 stdout 文本次序与措辞。若措辞变了 ⇒
    后果＝引用 boot 输出的别处证据件（票 228 那一族）需要对一下句子。

## 7. ⑥ 我跑了哪些尺（逐条命令，全有台件）

1. `date '+%Y-%m-%d %H:%M:%S %z'` ＋ `git log -1 --format='%h %ad %s'` ＋ `git rev-parse --abbrev-ref HEAD` ⇒ 起手锚（表头：`426a4d03`／`dev`／`17:25:03`）。
2. `git status --porcelain > .scratch/wisp/probes/33/r9/00-start-git-status.txt` ⇒ **272 行**。
3. `git status --porcelain -- cmd internal` ⇒ 起手 **0 行**；终态＝§5 那三枚。
4. `GOFLAGS= go test ./cmd/wisp -count=1 -v -run 'TestZzR9Shape' -timeout 20m` ⇒ `01-fifth-shape-1.txt`（F1 三发＋对照形三发；`rc=0`、`ok 15.567s`）。
5. `... -run 'TestZzR9ThreadExitTakesItsWindow'` ⇒ `02-thread-exit.txt`（三发 `IsWindow=false`，`rc=0`）。
6. `... -run 'TestZzR9ShapeFivePlain'` ⇒ `03-fifth-shape-plain.txt`（三发具名拒绝＋"排干即拆窗"，`rc=0`）。
7. `mv cmd/wisp/zz_33r9_probe_windows_test.go .scratch/wisp/probes/33/r9/zz_33r9_probe_windows_test.go.moved` ⇒ 载体搬出（⛔ 未删）。
8. `... -run 'TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever' -v` ⇒ `04-nail-green.txt`（`--- PASS (0.13s)`）。
9. 突变 M1：`sed` 中和 `showOnThread` 那一支 → `grep -n MUT-33r9`（1 处）→ 同条 `-run` ⇒ `05-mutation-1-no-outlet.txt`（FAIL 11.09s）→ 副本还原 → `md5sum` 逐字相同。
10. 突变 M2：`sed` 摘 `%w` → `grep -n MUT2-33r9`（1 处）→ `-run 'TestAC13NamedRefusal...|TestAC13BringUpRefusesAThreadWithAQueuedClose'` ⇒ `06-mutation-2-unwrapped-refusal.txt`（r7 钉 PASS＋本钉 FAIL）→ 还原 → `md5sum` 相同 → `grep -c MUT` 非零者零枚。
11. `"$(go env GOPATH)/bin/gofumpt.exe" -l <三枚名下件>` → `-w` → 复列空。
12. 名册：`... -run 'Panel|AC13|AC14|AC4|BallPanel|BallGesture|ListeningSocket' -v` ⇒ `07-roster-1.txt`＋`07-roster-names.txt`（21 名，逐名见 §4.1）。
13. 门禁：`GOFLAGS= go build ./...`／`GOFLAGS= go vet ./cmd/wisp`／`sh scripts/d22scan.sh`／`gofumpt -l` ⇒ `08-gates.txt`。
14. 整包发 1：`GOFLAGS= go test ./cmd/wisp -count=1 -v -timeout 25m` ⇒ `09-fullpack-1.txt`（159/1/0，§4.2）。
15. 枚举 AC#4 那枚红的射程：`grep -c "residentPanel\|RequestShow\|showOnThread\|errPanelRefusedThread" cmd/wisp/panel_host_windows_test.go` ⇒ **0**。
16. 整包顺序尺：`grep -n "=== RUN   TestAC4FocusReturnToPriorWindowGap33r5\|=== RUN   TestAC13NamedRefusal" 09-fullpack-1.txt` ⇒ **538 / 614**。
17. 隔离复跑那枚红：`GOFLAGS= go test ./cmd/wisp -count=2 -v -run 'TestAC4FocusReturnToPriorWindowGap33r5'` ⇒ `10-ac4-isolated.txt`（两枚 PASS）。
18. 整包发 2：同 14 那条命令 ⇒ `11-fullpack-2.txt`（§4.2 右列）。另：`tasklist` 三枚（msedgewebview2／wisp.test／go.exe）、
    `grep -c "^- \[ \]"`／`grep -c "^- \[x\]"` 票面（13／1）、`git diff --numstat HEAD -- <禁区>`（空）。

## 8. ⑦ 判不动的地方（交编排者裁，⛔ 本腿不代答）

| # | 哪一格 | 为什么本腿判不了 | 本腿立场（不是裁定） |
|---|---|---|---|
| R1 | **甲的"线程收摊"拆掉那扇合法活窗，算不算越线** | §2.4 三发读数说"线程结束窗随它没"，§2.1 说那扇窗的 owner 在世且自己的 Show 还能用；"该不该为出口接受这一拆"是产品形状判断，不是读数 | 若要"绝不许拆"，剩下的唯一形状是"置致命但线程 park 不退"＝新线程形状；本腿选甲并把它登记在此 |
| R2 | **`statusLine()` 那句 "panel thread could not start"（`panel_resident_windows.go:456-458`）** | 语义现在涵盖"运行中被具名拒绝"；改句子＝扩写面，且票 246 的 `TestAC246StatusLineSaysWhatTheLegDoesNot` 在盯措辞 | 记一笔措辞账，归下一枚动 `cmd/wisp` 的腿（与 §6 第 5 条同一条线） |
| R3 | **球手势那句 why 清单**（`cmd/wisp/resident_ball_windows.go:94`："it is starting up, its queue is full, or it already exited"） | 甲落地后这条通道今天会命中"already exited"，但"被具名拒绝"这一因不在清单里；改它＝第三枚产码件，超本腿写面 | 与编排者 16:42 第 4 格那枚 `lastShowErr` 欠账并成一条，派给票 248 或 228 那两发之一 |
| R4 | **第五形在生产里究竟经哪条通道落到过** | 本腿复认了"全仓零枚产码派干该线程队列"，也复认了拒绝在两条载体上都成立；但生产通道的**完整名册**我只对手势/stop 两条核过（§6 第 1 条）⇒ "生产可达"目前是〔v2 表内判断＋本腿半复认＋载体为人工植入〕 | 请 `33-v3` 用一发"装配根＋真常驻 boot"的形状判可达性；⛔ 别拿本腿这两枚人工载体当实证 |
| R5 | **乙是否真的不可实现** | 我的三条凭据都可盘上查（§3.1），但"原地重起一定不安全"这句我是从一次性通道＋`closeOnce` 推的，没跑过任何一形乙的实现 | 如果编排者要乙，需要另派一枚能写 `resident_windows.go` 的腿，并先补"重起后队列真的干净"的读数 |
| R6 | **WebView2 controller 跟着线程没之后 Go 侧那半会怎样** | §2.4 那三发是普通窗（问的是 Win32 线程收摊规则）；面板窗那一形我⛔ 没有对照读数（要再造一扇真窗并等线程死） | 〔未量〕。引用 §2.4 时带上"载体是普通窗"这句 |
| R7 | **整包发 1 那枚 `TestAC4FocusReturnToPriorWindowGap33r5` 的红算不算 AC#4 的产码缺陷** | 本腿只做到"射程零交集＋顺序在其后＋隔离两发全绿＋`0x30176` 是别家进程的窗"这四件事；"前台锁在 CI 镜像里会不会长期吃掉这一维"要 runner 才答得了（`33-v2` §6/N1 同一条线） | 具名登记为**既有负载敏感红**（编排者 12:12 裁定 2 第③件已认过它），⛔ 本腿未放宽、未改判据；发 2 的逐字结果补在 §4.2 |
| R8 | **CI 那两格（windows cli 腿有没有 WebView2 Runtime；干净检出的 `t.Skipf` 会把 `runtests.sh:98` 判红）** | 本机不可判，派单也禁"推一次看看" | 沿用 `33-v2` §6 的判不动，本腿不新增结论；⛔ 本表的绿不延伸到 runner |
