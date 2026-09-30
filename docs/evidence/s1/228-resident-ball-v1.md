# 228-v1 — 非实现者对抗验收：票 228 的 AC#1（宿主存在性有生产点数）与 AC#7（三处打给人看／写进日志的字符串）

验收腿：`228-v1`（只读＋突变台件，零产码改动）。锚点自取：`git rev-parse --short HEAD` = **`5fcdb088`**。
被验收实现腿的四枚 commit：`90cf65b2`（骨架）→ `1739451d`（产码＋三枚测试）→ `15d2b2cd`（AC#7 文案）→ `39bcc958`（收尾）。
射程**只有两格**：AC#1 与 AC#7（含 09-30 的射程更正）。AC#2／AC#3／AC#4／AC#5／AC#6 本轮不在射程，**本表不裁它们，也不得被读成裁过**。

---

## §0 逐格 1:1 表（判语）

| 格 | 判语 | 我自己的命令（不是编排者给的基线） | 我的读数（时刻 09-30 16:3x–17:0x +08） |
|---|---|---|---|
| **AC#1** 宿主存在性有生产点数 | **成立**（附下面两枚具名条件，都不翻成"不成立"） | `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp -count=1`；尺 `grep -rln "wisp/internal/ball" --include=*.go . \| grep -v .scratch`；生产调用者逐行现读 `cmd/wisp/resident_windows.go:27-135` | 整包 **rc=0／零枚 `--- FAIL`／148.168s**（默认档，未带 `-tags winlive`）；尺＝3 命中：`./cmd/balldebug/main.go`、`./cmd/wisp/resident_ball_228_test.go`（尺自己在 `_test.go` 里的字面量，`resident_ball_228_test.go:43`）、`./cmd/wisp/resident_ball_windows.go`。宿主调用点是 **`resident_windows.go:118 rb := startResidentBall(rt.Registry)`**，在 `runResident()` 体内、`rt.RunEventLoop()`（`:132`）**之前**，而 `runResident()` 的唯一生产调用者是 **`cmd/wisp/main.go:64`**（`main()` 的 len(args)==0 分支，`:58`）。⇒ **不是"字段填了真值"，是"构造并跑了"**：跑的证据见 §1 第二把尺（真起 `wisp.exe` 无参数、读它自己写进磁盘的记账）、§4/M5 后我自己在真机拉的 7 发探针（console 原文含"hotkeys live 4/4"、盘上 `ball-created`/`ball-stopped` 都在册），以及 §4 的 **M1**（摘掉那一行，`TestAC228ResidentLegIsTheBallHost` 与 `TestAC228ResidentLegReportsAndBooksItsBall` **两枚点名红**）。<br>条件＝下面 §0.1 的 **F-2／F-5／F-6／F-7** 四枚（都不把这一格翻成"不成立"）＋§1 末那枚"语义欠账"。 |
| **AC#7** 三处人读字符串 | **成立**（附一枚**具名阻塞欠账 F-1**：`cmd/wisp/main.go:8` 那行注释今天已是错事实，见 §2 末） | 逐处现读改后正文：`cmd/wisp/main.go:24-29`／`cmd/wisp/resident_windows.go:128-131`／`internal/proc/boot_windows.go:127-134`；对照 `git show 15d2b2cd` 与 `git show 1739451d -- cmd/wisp/resident_windows.go` 的改前原文 | 三处**都**已不含 `ticket 07`，也都不留将来式：<br>① `main.go:24-29`（`wisp -h` 正文）现文逐字 "hosts the floating ball window, its tray icon and its four global hot keys in this same process, then parks in an event loop that takes no task yet"＋把"托盘 Exit 无停机路径"与"球起不来只上报不停 boot"两条**负面**也写进去了 ⇒ 说实话。<br>② `resident_windows.go:131` 运行时行现文 "wisp: empty event loop running (no task is taken by this loop yet); %s (Ctrl+C exits cleanly)"，`%s` 由 `rb.statusLine()`（`resident_ball_windows.go:122-127`）填，两枚句子各在写它的那同一条语句里落成（成功 `:111-112`／失败 `:94`）⇒ 句子与事实同源，不能各自漂。<br>③ `boot_windows.go:132-134` 的 `slog.Info` 已拆成不带票号＋两条属性（`handled_here`／`bring_to_front`），逐字承认"这条循环自己没有窗口可提"⇒ 说实话。<br>⚠ **F-1 不撤销上面这一格**，理由是票面 ② 逐字「其余 25 处是注释……⛔ 不塞进本格（那会让 AC#7 永远勾不上），另立票 243」——见 §2 末的逐字对撞与我的裁法。<br>另附**具名条件 F-6**（同一枚字符串的第二张脸）：失败支那句被**说两遍**（`resident_ball_windows.go:96` 的 printf ＋ `resident_windows.go:129`/`:131` 那条 boot 报告各一遍），成功支只一遍 ⇒ **两支不对称**，撞红该包自己的"恰好说 1 次"不变式（读数 `console up=0 absent=2`，§4/M5）。两遍说的都是**同一句真话**，所以这一条**不撤 AC#7**，但必须与 F-1 一起进下一程的活单。 |

门禁四数（我自跑，不是抄实现腿）：`gofmt -l` 七枚被验文件＝**空**；`bash scripts/d22scan.sh`＝**PASS／clean**（`ban #8 cmd/ 69 枚`、`internal/ 473 枚`，⚠ 我第一次用 `go run ./tools/d22scan` 打的是**不支持的姿势**，根模块报 `does not contain package`，改用 wrapper 才有效——这条坑记给自己）；`go list -deps ./cmd/wisp` 与协程名册两把尺见 §1 末"零新增依赖边／零新起协程"。

### §0.1 具名欠账一览（全部由本腿现测，无一条来自注释自陈；⛔ 本腿一律未修）

| 号 | 一句话 | 证据（file:line） | 该谁做 |
|---|---|---|---|
| **F-1** | `cmd/wisp/main.go:8` 那句「The floating ball GUI is ticket 07.」今天已是**错事实**（球就在这条腿里），而同文件 `:24-29` 的 user-facing 正文刚被改反 | `cmd/wisp/main.go:8` 对照 `:24-29`；主张出处 `.scratch/wisp/probes/243/c1/census.md:75`（#24 行，逐字"同一文件同一主张的两个受众…只改串不改注释＝留两份真相"，归属列写"票 228"） | 票 228 下一程（一行注释）；⚠ 与 AC#7 射程更正②的"归票 243"**互相矛盾**，归谁要编排者拍（§6） |
| **F-2** | 失败支**没有任何自然覆盖**（`startResidentBall` 里 `ball.New` 写死、无注入面），只有 M5 突变走到过 | `cmd/wisp/resident_ball_windows.go:74`（直接调用）／`:93-98`（那一支）；突变读数 §4/M5 | 票 228 下一程（要注入面才算证过） |
| **F-3** | 失败支那句 "NO floating ball window" 在 `CreateWindowExW` 成功之后才失败的那一形里措辞偏窄（窗口曾短暂存在） | `internal/ball/ball_windows.go:215-243` 的失败次序＋`internal/ball/sta_windows.go:65-71`（返错即不跑泵、线程结束） | **不判红**，仅具名；真修在 `internal/ball`（非我地界） |
| **F-4** | 从 `signal.Notify` 起到函数返回，Ctrl+C／SIGTERM 的**默认动作被全程抑制**：若 D38(e) 某步卡住，**再按 Ctrl+C 救不回来** | `cmd/wisp/resident_windows.go:103-105`＋`:69`；保护面已由 §3 甲/乙形实测 | 具名上报（修法要么加"二次信号强退"，要么给 shutdown 步加超时——都超出两格射程） |
| **F-5** | `TestAC228ExitRequestDuringBootStillLeavesThroughD38E` 的名字承诺"boot 窗口内"，但它**不核自己走的是哪一支**；判据只靠球路径那 ≈235ms 空窗供养 | `cmd/wisp/resident_ball_228_windows_test.go:184-224`（无 "arrived during boot" 断言）；M4 实测红＝**弱钉而非空钉** | 票 228 下一程（把分支钉进断言） |
| **F-6** | **本腿新发现**：失败支那句被**说两遍**（`:96` printf ＋ boot 报告 `%s`），撞该包自己的"恰好说 1 次"不变式 ⇒ **在无桌面机器上这一枚用例会因重复而红** | 读数 `console up=0 absent=2`（§4/M5）；断言 `resident_ball_228_windows_test.go:78-81`；两处 printf `resident_ball_windows.go:96` 与 `resident_windows.go:129`/`:131` | 票 228 下一程（摘 `:96` 那句即可，成功支本来就不 printf） |
| **F-7** | boot 分支那句实话（"arrived during boot; the event loop was never entered"）**没有任何测试钉**，今天只由本腿的探针撑着 | 真机原文 §3 甲形；对照 `resident_ball_228_windows_test.go:82-85`（那一枚要求的是**另一支**的句子） | 票 228 下一程 |

---

## §1 装配可达性（编排者第 1 条：import 了 ≠ 装起来了）

### 逐行现读 `runResident()`（`cmd/wisp/resident_windows.go:27-135`）

| 行 | 原文形状 | 是不是生产路径 |
|---|---|---|
| `:27` | `func runResident()` | 唯一调用者 `cmd/wisp/main.go:64`，在 `main()` 的 `if len(args) == 0` 分支（`:58`）内 ⇒ **这就是 owner 双击图标那一形** |
| `:36` | `rt, err := proc.Boot(env)` | 生产 |
| `:60` | `sink, sinkErr := installLogSink(rt.Layout.DataDir)` | 生产 |
| `:69` | `defer func() { records := rt.Shutdown(false) … }()` | 生产（第 2 个 defer） |
| `:103-105` | `bootExit := make(chan os.Signal, 1)`／`signal.Notify(bootExit, os.Interrupt, syscall.SIGTERM)`／`defer signal.Stop(bootExit)` | 生产（第 3 个 defer）——见 §3 |
| **`:118`** | **`rb := startResidentBall(rt.Registry)`** | **生产调用点。它不是构造给测试看的：它的返回值 `rb` 立刻被 `:129` 与 `:131` 两条 `fmt.Printf` 读，又被 `:119 defer rb.stop()` 持有** |
| `:119` | `defer rb.stop()` | 生产（第 4 个 defer，注册最晚 ⇒ LIFO **最先**跑） |
| `:126-133` | `select { case sig := <-bootExit: … default: … reason = rt.RunEventLoop() }` | 生产 |

`startResidentBall`（`cmd/wisp/resident_ball_windows.go:71-117`）**不是把球放进一个字段里就返回**：`:74` 直接 `ball.New(ball.Options{…})`，`internal/ball/ball_windows.go:168-172` 里 `New` 会 `opts.Registry.Spawn("ui-sta", "ball", nil, …)` 起线程并 `b.sta.waitStarted()`（`:173`）**阻塞到 Win32 真返回结果**为止。所以"调用了"与"球在场"之间没有第三次抽象。

### 判据要能区分"填了字段"与"构造了没跑"——我用三把尺，不是一把

1. **静态形状**（`resident_ball_228_test.go`，走 AST，跨平台）：clause(1) `:205` 要求生产文件 import；clause(2) `:211` 要求包内有函数**调用 `ball.New`**；clause(3) `:246` 要求 **Windows 那一枚 `runResident` 的函数体里出现该 host 的调用**；clause(4) `:266` 要求同一函数体里有 `defer ….stop()`。⇒ clause(3) 正是"填字段不跑"这一形的克星：把 host 调用摘掉，`called[]` 里就没有 `startResidentBall`。
2. **运行时**（`resident_ball_228_windows_test.go:53`）：`buildWispForTest` 造**真 `wisp.exe`**、`bootResidentLeg` 用**无参数**命令行拉起（`resident_sink_nail_127_windows_test.go:195-217`，`exec.Command(exe)` 不带 args），再核两件事——控制台那句球姿态**恰好出现 1 次**（`:78-81`）且**与"empty event loop running"同一轮 boot**（`:82-85`），以及**同一个进程写进磁盘的日志里** created／refused **必须异或**（`:117-124`，两边都没有＝`t.Fatalf`）。
3. **窗口在场**（`-tags winlive`，`resident_ball_live_228_windows_test.go:89`）：按 **pid** 用 `FindWindowExW` 认领 `WispBallWindow`，要求恰好 1 枚（`:139-142`）、`hotkeys live 0/4` 判红（`:151-153`），退出后要求该 hwnd 不再是窗口且不再被那个 pid 认领（`:168-177`）。

### M1 定向突变：摘掉宿主那一次调用

见 §4 的 **M1** 一发（改后读数、点名用例、还原逐字）。编排者要我指名"必须红的用例"，我给的是**两枚**：`TestAC228ResidentLegIsTheBallHost`（clause 3）与 `TestAC228ResidentLegReportsAndBooksItsBall`（created/refused 双缺那一支）。若两枚都不红就是 AC#1 假绿——读数落定后此格判语不变，因为实测两枚都红。

### 零新增依赖边／零新起协程（编排者第 6 条，我自跑两把尺）

**依赖边**：
- `git diff 90cf65b2~1 39bcc958 -- go.mod go.sum` ＝ **空** ⇒ 实现腿那句"零新增模块依赖"成立。
- `cmd/wisp` **生产文件**（剥 `_test.go`）的 `internal/*` import 集，逐枚 `git cat-file blob` 前后对照：**16 枚 → 17 枚**，差集**恰好只有** `github.com/CarlosShao/wisp/internal/ball` 一条。⇒ ⚠ **"零新增依赖边"这句要分开判**：对 `internal/proc` 那一处成立（`git show 15d2b2cd -- internal/proc/boot_windows.go` 只改字符串＋注释，无 import 变化，我逐行看过）；对 `cmd/wisp` **不成立、而且必须不成立**——那一条新边**就是 AC#1 的尺本身**。实现腿交件写的是"零新增**模块**依赖"，措辞对得上；⛔ 谁把它读成"零新增依赖边"谁就会去摘那条边。
- `go list -deps ./cmd/wisp` 现跑＝`internal/observe`／`internal/statemachine`／`internal/ball` 全在图内；`internal/ball` 的外部依赖只有 `golang.org/x/sys/windows` 与 `github.com/pelletier/go-toml/v2`（**两枚都已在册**，`internal/ball` 早被 `cmd/balldebug` 用过 ⇒ 二进制层面没有新第三方代码）。

**协程**：
- `grep -nE "^[[:space:]]*go " cmd/wisp/resident_ball_windows.go cmd/wisp/resident_windows.go` ＝ **零命中** ⇒ "本轮零新起协程"这句**在这两枚文件里成立**。
- 全仓 `Spawn` 名册现算：常驻这条腿真正会跑到的新协程**只有一枚**，就是 `internal/ball/ball_windows.go:170` 的 `opts.Registry.Spawn("ui-sta", "ball", nil, …)`——**`ui-sta` 在 `internal/observe/goroutine.go:44` 的 `ResidentNames` 里在册** ⇒ 不会造出 `RosterReport.Unknown`（`:419`）那枚假泄漏。其余在册/在图的名字（`log-flusher` `internal/observe/logging.go:96`、`audio-capture`、`db-writer`、`watchdog`、`approval-waiter`）本轮**都没被这枚改动新带动**（`watchdog`/`approval-waiter` 属 `run.go` 那条腿：`startConfigReload` 的唯一调用者是 `cmd/wisp/run.go:708`，我现查）。
- ⚠ 一处**读数口径**要写死：`internal/proc/boot_windows.go:108` 的 `ResidentOverBaseline` 检查跑在 `rt.Boot()` 里、**早于**本票新增的球宿主（宿主在 `resident_windows.go:118`），所以 `ui-sta` **不可能**把这枚 boot 门撞红——这是我读次序读出来的，不是实测（实测"常驻腿此刻活着的 resident 名有几枚"本表没做，见 §5 第 3 条）。对照形状：`cmd/balldebug/main.go:33-42` 逐字承认调试器的三枚协程**故意不借在册名**，代价是"one `goroutine outside the D38 roster` WARN per spawn"——**那一形才是假泄漏告警的形状，本票没造它**。

### AC#1 的语义欠账（具名，不藏）

票面 AC#1 原文要求的是"≥1 枚来自 **`cmd/wisp` 那条真正会跑任务的腿**"。今天装球的是**常驻那条腿**，而它**仍不接任务**（`:131` 那句 `no task is taken by this loop yet` 是真话，`RunEventLoop` 只做激活记账）。这一格的口径来自**编排者 09-29 的 AC#0 裁定"走甲"**（票面第 3 行逐字「方向由编排者按 owner 长期原话裁死＝走甲（把审批门与任务通路接进**常驻**那条腿）」）＋ D2 拓扑（`PLAN.md:74/83-88`）：**被提名"真正会跑任务"的那条腿就是常驻腿**，`wisp run` 那一形在一发即退的语义下永远装不了球（票面"为什么不许就在 wisp run 里挂个球"三支）。
⇒ 我按票面 AC#1 的**字面尺**（grep＋生产调用者）判**成立**；但**AC#1 的立格目的**（"用户看得见球的那条进程＝干活的那条"）要到 AC#2／AC#3／AC#4 才闭合，**这一格不替它们背书**。编排者若要连目的一起勾，需要的是那三格的读数，不是这一格。

---

## §2 文案两支与失败路径（编排者第 2 条）

### 失败支今天有没有"测试真走到"？——**没有：只有突变走到过**

`startResidentBall` 的错误支（`resident_ball_windows.go:93-98`）由 `ball.New` 返回非 nil err 才进。`ball.New` 的 err 只来自 `b.sta.waitStarted()`（`internal/ball/ball_windows.go:173`）即 `createOnSTA` 的失败：`ensureFactories`／`registerBallClass`／`CreateWindowExW==0`／`newRenderer`／`addTrayIcon`（`ball_windows.go:182-243`）。本机有桌面 ⇒ 默认档与 winlive 档**都走成功支**。
**没有测试注入面**：`startResidentBall(reg *observe.Registry)` 直接把 `ball.New` 写死在函数体里（`:74`），没有 `func New` 字段、没有 seam、没有第二枚构造签名，所以**产码里不存在任何能让默认档走失败支的入口**。
⇒ 编排者这一条我**坐实为缺陷形状**：票 181 AC#7 栽的那一形（"有一条用例"被当成"失败支被证过"）在这里**部分是**真的——失败支**只被 M5 突变证过**（§4 末发），被自然运行**没证过**；⛔ 而且**它一旦被走到就撞红该包自己的不变式**（读数 `console up=0 absent=2`，红点 `resident_ball_228_windows_test.go:78-81`；成因见 §4/M5 第 2 条与 **F-6**：失败支那句被 `resident_ball_windows.go:96` 与 `resident_windows.go:129/131` 各说一遍，成功支不说第二遍）。⇒ **失败支今天既是死路（无人自然走到）又是生路（走到就会红）**，两件事都不是 AC#1/AC#7 的字面判据，但都必须入账。
⚠ 这不是 AC#1/AC#7 的**判据缺件**：票面两格都没写"失败支要有自然覆盖"，也没写"每条只说一次"。我按具名欠账 **F-2／F-6** 交回，不擅自抬成不成立。CI 那一侧：`internal/ball` 的失败支若在**无桌面 runner** 上跑常驻腿就会自然走到——**那一台上这一枚用例会红，而红因是重复不是缺陷**。本表**未核 CI 环境是否有桌面**（判不动，见 §6）。

### 进程该继续跑还是该退出？——**它今天选"继续跑"，我判这一支与在册定案相合**

`resident_ball_windows.go:93-98` 的做法：填 `verdict`、`slog.Error`、再 `fmt.Printf("wisp: %s\n", rb.verdict)`、`return rb`（**不 `os.Exit`、不 return err**）。函数文档 `:70` 逐字「Errors are returned as a verdict, never as a failure to boot」。
对照两处在册定案：
- **票 117 同腿先例**（`resident_windows.go:61-68`）：日志目录打不开时「Loud, and it does not stop the app: a log directory that will not open must not become a way to keep Wisp from starting」——同一姿势，球这一格是照它抄的（`:34` 文件头也逐字承认「the same stance installLogSink takes above」）。
- **在册常记**（`AGENTS.md`／记忆《安全闸门的范围》）：「无条件拒绝启动会堵死 owner 的联调方式」。球是**输出面**，不是**门**；一台没有桌面的机器（服务会话／CI runner）必须仍能 boot、仍持有 Job Object、仍走 D38(e) 离开。
- **票 128"拒绝启动"那一套**射程是**数据根／密钥解析不出来**（`cmd/wisp/dataroot_128_test.go` 一族），那是"没有安全落盘位置就不能开机"的形；球失败与之不同形（没有隐私后果）。
⇒ **判：选型正确**。⛔ 我不判"应该拒绝启动"，那会造出一条与票 117 相抵触的新规矩，而改契约须人工批准。

### 失败支文案是否说实话（M5 那一发读到的原文）

`rb.verdict = fmt.Sprintf("this process has NO floating ball window: %v", err)`（`:94`）——`%v` 就是 `ball.New` 的原 err，句子把"没有窗口"与"为什么没有"绑在同一条语句里，且**成功支写的是另一枚句子**（`:111`），两枚不可能同真（`resident_ball_228_windows_test.go:74-77` 就是在钉这个）。
⚠ 一处**措辞级不精确**（具名 F-3，不判红）：`createOnSTA` 有可能在 `CreateWindowExW` **成功之后**才失败（`newRenderer` `:235` 或 `addTrayIcon` `:241`），此时那一枚 hwnd 已存在、句子已说 "NO floating ball window"。真实后果有限：`staThread.start` 在 `create` 返错时直接 `return`（`internal/ball/sta_windows.go:65-71`），不再跑消息泵、OS 线程结束 ⇒ 那个隐藏窗口随线程销毁；且失败发生在 `addTrayIcon` 之后已无可失败步骤 ⇒ **不会留下托盘、也不会留下 hotkey**（`registerAll` 在 `:244`，`ShowWindow` 在 `:247`）。所以我判它"措辞偏窄但不错到能骗人"，**不列 AC#7 不成立**；实现腿自己在 `notes.md` 里也具名留了"ball.New 失败路径不销毁半成品窗口在 internal/ball 里没动"——那一处属 `internal/ball` 地界，本表不裁。

### F-1：`cmd/wisp/main.go:8` 那行注释（编排者第 5 条，我按逐字对撞裁）

现状**逐字**：`:5-10` 的包文档注释末句仍是「The floating ball GUI is ticket 07.」，而同一枚文件 `:24-29` 的用户正文已改成「hosts the floating ball window … in this same process」。
"注释与紧邻的用户文案是同一句话的两面"这一主张——**不是我提的，也不是实现腿提的**：票 243 的只读普查表 `243-c1` 第 24 行逐字写着「这句与 §1 的 `main.go:25` usage 串是**同一文件同一主张的两个受众**（godoc 面／用户面）——只改串不改注释＝留两份真相」，判定列是「**还在承诺未来**（兼：错事实）」，归属列写的是「**票 228**」（`.scratch/wisp/probes/243/c1/census.md:75`，另见 `:89` 把它列进第二优先）。
⇒ **主张我判"成立"**。但**本格判语仍是 AC#7 成立**，理由只有一条、且是票面自己的话：AC#7 的射程更正 ② 逐字把注释那一堆划出本格（「其余 25 处是注释……逐处判"历史出处 vs 过期承诺"⛔ 不塞进本格（那会让 AC#7 永远勾不上），另立票 243」），并把 ① 的"完成判据"限死在三枚人读字符串上，而那三枚**全改对了**。
**我的裁法（不和稀泥）**：AC#7 **成立**；F-1 作为**本票名下欠账**入账，措辞按编排者能办事的那形给死——**若翻 AC#7 的勾，必须同时把 F-1 记成 228 名下"未清的一行"并派给下一程写腿**（一行注释、`main.go:8`，改法就是把末句换成带条件的事实句，与 `:24-29` 同真），**否则这一枚勾读起来就是在替"球要等票 07"这枚今天已假的事实背书**。⚠ 与此同时 243-c1 的归属列写的是票 228，与射程更正②的"归票 243"**互相矛盾**——这一处**不是我造成的，也不是我能定的**（改票面正文＝契约面），我按未定义即停在 §6 具名上报，**不擅自选边补一条规矩**。

---

## §3 退出接管与 0xc000013a（编排者第 3、4 条）

### `bootExit` 有没有人读？

**有，且只有一个读者、只读一次**：`resident_windows.go:126-133` 的 `select { case sig := <-bootExit: … default: … }`。它跑在 `startResidentBall` 返回之后，取的是"非阻塞看一眼"的形（`default` 分支）。
**之后没有人再读它**：`RunEventLoop` 自己 `signal.NotifyContext`（`internal/proc/boot_windows.go:121`），`sigCtx` 是另一枚通道；`bootExit` 直到函数返回才被 `defer signal.Stop(bootExit)`（`:105`）摘掉。
**缓冲满之后第二发信号会怎样**：`bootExit` 是 `make(chan os.Signal, 1)`（`:103`），首发包进缓冲、无人取；此后 `os/signal` 对所有已注册通道做**非阻塞发送**，通道满则**丢弃该发**——**不会**回落到默认动作、**不会**把进程打死。⇒ 这一条我**不是只靠推理**：§3 末的甲形探针（5/5 发 exit 0、`0xc000013a` 零枚）量的就是它。**这一处编排者怀疑的三个后果，实测形状是：不崩、不二次死、但第二发在 `bootExit` 这一路是哑的**；真正让进程离开的是 `sigCtx`，而它由 `RunEventLoop` 在读。**不判红**，理由：`default` 那一支已经把进程交给 `RunEventLoop`，此后一切退出请求由 `sigCtx` 承接；`bootExit` 哑掉的是**同一发的重复投递**，不是唯一出口。
⚠ 具名代价（F-4，交回不擅自修）：从 `signal.Notify` 起，**SIGTERM／os.Interrupt 的默认动作在整个 `runResident` 期间被抑制**，包括 D38(e) 正在跑的那段（`:69` 的 defer 与 `signal.Stop` 之间）。若 shutdown 序列里某步卡住，**再按 Ctrl+C 救不回来**（`bootExit` 已注册不读／`sigCtx` 的 `stop()` 已随 `RunEventLoop` 返回而执行）。今天 10 步序列有没有能卡的步＝本表未证（§6）。

### 接管窗口与正规 shutdown defer 之间会不会重复触发？

不会重复**触发关闭**：`rt.Shutdown` 只在 `:69` 那一枚 defer 里被调；`rb.stop()` 只在 `:119` 那一枚 defer；`sink.close()` 只在 `:67`。信号本身**不直接**触发任何一枚，它只让 `RunEventLoop` 返回。
LIFO 次序现读（注册序 `:67`→`:69`→`:105`→`:119` ⇒ 执行序 **`rb.stop()` → `signal.Stop(bootExit)` → `rt.Shutdown(false)` → `sink.close()`**）：`:109-113` 注释自称"球这一格做的是 D38(e) 第 2 步的活（`internal/proc/shutdown.go:18` hotkey + wake-word listening stops），跑在第 9 步 Job 关闭之前，且两边的日志都还能落到 sink"——三条我都现读对上：`shutdown.go:18` 逐字含 step 2 的措辞；`sink.close()` 确实最后；`Ball.Close()` 确实 join `ui-sta`（`internal/ball/ball_windows.go:935 <-b.sta.handle.Done()`）。
⚠ 一处**注释与代码的次序说法**要具名：`:113` 说 defer 是"registered last so LIFO runs it FIRST"——**成立**（`:119` 是本函数最后一枚 defer）。但**同一枚文件头 `:16-20`** 与 `:58-59` 把 `signal.Stop` 排在 `rb.stop()` 之后跑这件事没写（无害：`stop()` 不发日志之外的事件），我复核后判**不影响正确性**。

### 0xc000013a 那一形有没有真机证据？——**有，但被钉住的方式比实现腿的说法弱一档**

- **票 127 那枚 `TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink`（`resident_sink_nail_127_windows_test.go:543`）**：我读了全文断言（`:552-619`），**一字未放宽**：它要求 record 0 是 `residentEarlyResolverMsg`、install 记账**必须落在 index 1**（`:597`）、install 之后必须存在 `shutdown step` 轨迹（`:607`）、**文件最后一条必须是 shutdown 记录**（`:611-614`）、stdout 必须有 D38(e) 报告（`:615`）。`git log -- cmd/wisp/resident_sink_nail_127_windows_test.go` 末次 commit = **`9b5d64d8`（票 130）**，**早于本票四枚** ⇒ 实现腿**没碰过这枚钉**，也没有"为变绿改断言"这回事。⚠ 关键：它**并不断言信号落在 boot 窗口内**——它先等 sink 文件出现（`:552-555`）再 `breakToLoop()`，而 sink 文件在 `installLogSink`（`resident_windows.go:60`）就出现，`Notify` 在 `:104`，只差几条语句。所以这枚用例在修复**之前**红，是因为**球路径那 200+ms 里没有任何处理器**；修复之后它绿。它绿**不等于**"boot 窗口的退出被证过"。
- **真正冲这一形的是新用例 `TestAC228ExitRequestDuringBootStillLeavesThroughD38E`（`resident_ball_228_windows_test.go:184`）**：姿势与 127 同（sink 文件出现即 `breakToLoop`），断言是 exit code 必须 0（`:205-209`，`0xc000013a here means nothing was registered for the signal while the ball came up`）＋ install 记录＋shutdown 轨迹（`:211-220`）。
  ⚠ **这一枚的名字比它的断言强**：它**不核自己走的是哪一支**（既不要求 `:128` 那句 "arrived during boot" 出现，也不禁止 "empty event loop running" 出现）。⇒ 若某一发的调度让 break 落在 `RunEventLoop` 装上 `sigCtx` 之后，这一枚**照样全绿**，而它名字承诺的那一段根本没被走。突变 M4 之所以还打得中，是因为球路径本身给了约 200ms 的空窗（实现腿现测 268ms，我复跑 §4 读到的实际耗时见那发）；**空窗由球路径长度供养，不由任何一枚断言供养**。这是具名缺陷 **F-5**：判据对"摘掉 Notify"敏感，对"这一发究竟落在哪个窗口"**不敏感**，所以它是**弱钉**而非空钉（M4 实测红，见 §4）。
- **`resident_windows.go:85-102` 那段注释里给的数**（268ms／sink installed 16:05:20.7746／ball booked 16:05:21.0424）：来自实现腿自己那台机的一次日志读数，**本表不复现该毫秒数**（性能类判据要单包安静复跑才入账，本轮我不做性能判定），我只复现"**这段窗口里 Ctrl+C 会不会丢 D38(e) 轨迹**"这一**布尔**判据。⚠ 顺带读到本机同一形的时间戳差 ≈**235ms**（sink installed 16:49:18.290 → ball booked 16:49:18.525，我这发探针的原始输出），**只用作"空窗确实存在且量级够大"的旁证，不作性能判定入账**。

### 两发 Ctrl+C 的真机读数（编排者第 3 条要的那一发，我自己写的探针）

台件：`.scratch/wisp/probes/228/v1/twobreak/`（自带 `go.mod`＝独立模块，根模块的 `./...` 不会编它，共享树的门禁不受它影响）。姿势：`wisp.exe` 造到 `$TEMP/wispprobe228v1/`，`CREATE_NEW_PROCESS_GROUP` 拉起、`WISP_ENV=test`＋`WISP_TEST_DATA_DIR` 指向临时目录；**sink 文件一出现（或再等 `-file-delay`）发第一发 `CTRL_BREAK_EVENT`，隔 2–4.6ms 再发第二发**，然后核退出码／控制台原文／盘上日志。每一发跑前 `tasklist //FI "IMAGENAME eq wisp.exe"`＝**0 枚**。

**甲形（`-file-delay 0`，共 5 发：3＋2）＝两发都落在球路径里、`:126` 那枚 `select` 之前**：

| 读数项 | 结果 |
|---|---|
| 第一发／第二发 `GenerateConsoleCtrlEvent` | 两次都 **rc=ok**（间隔 3.397／3.588／4.611／3.75ms） |
| 退出码 | **5/5 发 exit 0**，`0xc000013a` **零枚** ⇒ **第二发被丢掉、没有把进程硬杀** |
| 走的是哪一支 | console 逐字 `wisp: exit request (interrupt) arrived during boot; the event loop was never entered; the floating ball window is up in this process (tray icon added, hotkeys live 4/4: summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc live, panel=Ctrl+Alt+P live)` |
| 时间戳旁证 | sink installed `16:49:18.290` → ball booked `16:49:18.525`（**≈235ms 的空窗**，只当旁证、不作性能判定入账） |
| 收口 | `ball: resident leg destroyed the ball window…` 在 `16:49:18.530`，随后 7 条 `shutdown step skipped (module not present)`，末行 `wisp: exited through the D38(e) shutdown order (10 steps, 0 failed)` |

**乙形（`-file-delay 400ms -gap ~3ms`，共 2 发）＝第一发落在循环里、第二发落在 `rb.stop()`＋`rt.Shutdown()` 正在跑的时候**——**这一形才是"接管窗口与正规 shutdown defer 之间会不会重复触发"那一问的正面对象**（甲形的第二发其实还落在球路径里，这一点我在第一次写表时判错过、已按乙形改写）：

| 读数项 | 结果 |
|---|---|
| 第一行报告 | console 逐字 `wisp: empty event loop running (no task is taken by this loop yet); the floating ball window is up in this process (…)` ⇒ **default 支**，进程已进 `RunEventLoop` |
| 第二发（收口进行中，间隔 2.907／3.02ms） | rc=ok；**退出码 exit 0**；`exited through the D38(e) shutdown order: true`；`shutdown-step lines=7`；`ball-created=true`／`ball-stopped=true` |
| 重复触发？ | **无**：没有第二次关闭、没有半途退出、没有丢轨迹 |

⇒ 五件事落定：
1. **`bootExit` 那枚 1 格缓冲的后果是"丢"，不是"死"**（甲形证的正是这一条：第二发无人读、缓冲满、进程照走 D38(e)）。
2. **接管窗口与正规 shutdown defer 之间不会重复触发**（乙形直接证：第二发落在收口进行中，仍 exit 0、十步 0 failed、球与轨迹的记录都在）。
3. **`resident_windows.go:126-133` 的 boot 分支真能走到，且那句话说的是实话**——"arrived during boot; the event loop was never entered"＋同一枚球的事实句都在，`event loop ending (…)` 复用的还是同一枚 reason ⇒ AC#7 的"不先说空转、再决定不进循环"这一条**修对了**，我用真机 console 原文证的，不是注释。
4. ⚠ **这一形在测试里没有钉**：**没有任何一枚用例断言那句 "arrived during boot"**（`TestAC228ExitRequestDuringBootStillLeavesThroughD38E` 只核退出码与两条记录；`TestAC228ResidentLegReportsAndBooksItsBall` 反而要求 `empty event loop running` 在场，走的是另一支）。⇒ 具名 **F-7**：boot 分支的文案今天**只由我这发探针撑着**；把它钉进用例属下一程（本腿不写用例、不改产码）。
5. ⚠ **F-4 的保护面成立、代价面仍未证**：第二发被丢＝**再按 Ctrl+C 也救不出一个卡住的 shutdown**。今天十步在这台机上 0 failed，**没有卡住的形可测** ⇒ 代价面留在 §6 判不动。

---

## §4 突变逐发读数（改前／改后＋指名用例）

台件姿势（每一发同一套）：
- 起手名册（**等于终态名册**，共享树不写绝对式）：`git status --porcelain -- internal cmd` 于 **09-30 16:30**＝**空**；锚点 `5fcdb088`。
- 每一发：先 `grep -c` 断言锚点在该文件里**恰好命中 1 次**，再落盘，跑**指定的那几枚用例**（`-run`＋`-v`，判红**只认 `--- FAIL`**），跑完用**内存里那份原字节**还原并打印 `RESTORE <path> bytes=<n> exact=True`。
- ⛔ 全程未用 `checkout`／`stash`／`reset`／`--amend`／`clean`；⛔ 未用 `-overlay`（本包有走 AST 的用例，overlay 对它结构性失明——`resident_ball_228_test.go:77` 用 `parser.ParseFile` **读盘**，overlay 买不到）。
- 必带 PATH：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`（漏了会 `exit status 0xc0000135` 且**没有 `--- FAIL` 行＝一枚用例都没跑**）。

### 改前基线（未突变，默认档）

```
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp -count=1
=> ok  github.com/CarlosShao/wisp/cmd/wisp  148.168s   rc=0   `--- FAIL` 计数 = 0
```
桌面三连（16:30，突变前）：`tasklist //FI "IMAGENAME eq balldebug.exe"`＝**0 枚**；`… eq wisp.exe`＝**0 枚**；`gh run list --limit 3`＝最近三枚全 `completed`（`slo-fresh` success／`ci` 两枚 failure），**本机无 `slo-full` 在跑** ⇒ 本轮可读数；但**本表不采任何性能／句柄类读数**（§3 已具名 268ms 不复现）。

### M1 — 摘掉宿主那一次调用（AC#1 的命门）

改法：`cmd/wisp/resident_windows.go:118` 由 `rb := startResidentBall(rt.Registry)` 改成 `rb := (*residentBall)(nil)`（锚点 `rb := startResidentBall(rt.Registry)` 改前 `grep -c`＝**1**；`:119 defer rb.stop()` 保留，因为它对 nil 安全，`resident_ball_windows.go:141`）。

改后读数（`go test ./cmd/wisp -count=1 -v -run 'TestAC228ResidentLegIsTheBallHost|TestAC228ResidentLegReportsAndBooksItsBall'`，默认档）：

```
resident_ball_228_test.go:214: ball host functions (production): resident_ball_windows.go:startResidentBall
resident_ball_228_test.go:247: AC#1 RED: runResident (resident_windows.go) calls none of the ball hosts …
--- FAIL: TestAC228ResidentLegIsTheBallHost (0.01s)
resident_ball_228_windows_test.go:121: AC#1 RED: neither "ball: the resident leg created the floating ball window"
    nor "ball: the resident leg could not create …" is in the file this process wrote (records: [winsec: sealing path
    resolver installed / wisp: persistent log sink installed / shutdown step skipped (module not present) x7]) …
--- FAIL: TestAC228ResidentLegReportsAndBooksItsBall (3.31s)
FAIL github.com/CarlosShao/wisp/cmd/wisp 3.380s
```

⇒ **两枚点名用例都红，AC#1 不是假绿**。红的位置正是编排者要的那一形：clause(3)（`:246-249`）与运行时那一枚的"created/refused 双缺"支（`:120-124`）。
⚠ **一条只有这一发才看得见的细节**：M1 下**控制台那一枚断言没红**——`(&nil).statusLine()`（`resident_ball_windows.go:123-125`）回的是 "this process has NO floating ball window (the ball host never ran)"，**含** `ballAbsentClaim` ⇒ 句子层"说了话"，真正把假绿打掉的是**磁盘记账**那一支（`:120-124`）。⇒ §1 那三把尺里承重的是第二把（记账），不是第一句；实现腿若只核句子，M1 会溜过去。
还原：`git cat-file blob HEAD:cmd/wisp/resident_windows.go > …` ⇒ `RESTORE cmd/wisp/resident_windows.go bytes=6109 exact=True`（`git diff --quiet` 空）。

### M2 — 摘掉 `defer rb.stop()`（实现腿自称会红，我独立复跑）

改法：删 `cmd/wisp/resident_windows.go:119` 那一行（锚点 `defer rb.stop()` 改前 `grep -c`＝**1**）。

改后读数（同一命令）：

```
resident_ball_228_test.go:266: AC#1 RED: runResident (resident_windows.go) defers no ball teardown, so the hot keys
    and the tray icon it installed are left to process death. See internal/proc/shutdown.go:18 (D38(e) step 2).
--- FAIL: TestAC228ResidentLegIsTheBallHost (0.01s)
resident_ball_228_windows_test.go:154: AC#1 RED: this process created a ball (record 2) and never booked destroying it …
--- FAIL: TestAC228ResidentLegReportsAndBooksItsBall (3.44s)
FAIL github.com/CarlosShao/wisp/cmd/wisp 3.509s
```

⇒ **实现腿自称的那枚 M2 我独立复现，两枚都红**；"record 2＝created、teardown 缺席"与 §3 末正常态读数（created=2／stopped=3）**同形互证**，说明 LIFO 那一套不是注释里的说法而是盘上真有序列。
还原：`RESTORE cmd/wisp/resident_windows.go bytes=6109 exact=True`。

### M3 — 缺一枚手势执行者

改法：删 `cmd/wisp/resident_ball_windows.go:90` 的 `OnDragEnd: func() { recordBallDragEnd() },`（锚点改前唯一命中，Edit 以 1 处替换执行）。

改后读数（`-run 'TestAC228BallHostAnswersEveryGesture'`）：

```
resident_ball_228_test.go:315: AC#1 RED: the resident ball host leaves 1 of internal/ball's 10 gesture callbacks
    unset: [OnDragEnd]. A nil callback is a click or a hot key that reaches no one and logs nothing.
--- FAIL: TestAC228BallHostAnswersEveryGesture (0.02s)
```

⚠ **本腿此处犯错并当场纠正，记录不抹**：我第二次"还原"时误把该行**插了两遍**（`:90`/`:91` 重复键＝编译期错），是改用 `git cat-file blob HEAD:` 复原才真正还原——`grep -c "OnDragEnd"` 从 2 变 1、`git diff --quiet` 空。⇒ 本仓记忆里"还原只用 `git cat-file blob HEAD:`、不用 checkout"那条，我自己踩了一遍才对上；此后每一发都按它做。
还原后复跑：`--- PASS: TestAC228ResidentLegIsTheBallHost (0.01s)`／`--- PASS: TestAC228BallHostAnswersEveryGesture (0.01s)`，`ok … 0.071s`。

### M4 — 摘掉 boot 窗口的 `signal.Notify`

改法：把 `cmd/wisp/resident_windows.go:104` 那一行注释掉（锚点 `signal.Notify(bootExit, os.Interrupt, syscall.SIGTERM)` 改前 `grep -c`＝**1**；`:103`/`:105` 保留以维持 `select` 形状）。
⚠ 第一次尝试直接注释导致 `vet: cmd\wisp\resident_windows.go:10:2: "syscall" imported and not used` ⇒ 输出是 `[build failed]`，**不是 `--- FAIL`**——按本仓仪器规矩它**既不算红也不算绿**；补 `_ = bootExit`／`_ = syscall.SIGTERM` 两行让编译只丢 Notify 这一件事之后才拿到下面的读数。**这条坑记给下一位**：摘一行引发的 build failed 会被误读成"红了"。

改后读数（`-run 'TestAC228ExitRequestDuringBootStillLeavesThroughD38E|TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink'`）：

```
resident_ball_228_windows_test.go:206: AC#1 RED: the child answered a boot-time Ctrl+C with exit status 0xc000013a,
    want a clean exit 0 through the D38(e) order. …
resident_ball_228_windows_test.go:214: AC#1 RED: a boot-time exit produced no "wisp: persistent log sink installed" record (records: [])
--- FAIL: TestAC228ExitRequestDuringBootStillLeavesThroughD38E (3.73s)
resident_sink_nail_127_windows_test.go:569: the child exited through the shutdown path with exit status 0xc000013a, want 0
resident_sink_nail_127_windows_test.go:583: … holds a pipeline file with no record in it …; the leg exited exit code 3221225786 (exit status 0xc000013a)
--- FAIL: TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink (3.38s)
FAIL github.com/CarlosShao/wisp/cmd/wisp 7.168s
```

⇒ **0xc000013a 那一形我在这台机上亲手复现，两枚用例各一发**：日志文件 0 条记录、sink 没 flush、D38(e) 轨迹全丢——与实现腿自陈后果**逐字对得上**。⛔ 同一发把 §3 的 **F-5 从"可能空转"降级为"本机确实打中"**：摘掉 Notify 会红，说明这一枚真落在球路径那段空窗里。
还原：`RESTORE cmd/wisp/resident_windows.go bytes=6109 exact=True`。

### M5 — 让 `ball.New` 一定失败（失败支是否真走到／进程选型／文案说实话）

改法：`cmd/wisp/resident_ball_windows.go:93` 的 `if err != nil {` 改成 `if true { // MUTATION M5`（锚点 `^	if err != nil {$` 改前 `grep -c`＝**1**）⇒ 成功支变不可达、失败支每次必走；`err` 仍被那一支用（`:94`）所以编译干净；**没动 `internal/ball`**。

**默认档读数**（`-run 'TestAC228ResidentLegReportsAndBooksItsBall'`）：

```
wisp: this process has NO floating ball window: <nil>
wisp: empty event loop running (no task is taken by this loop yet); this process has NO floating ball window: <nil> (Ctrl+C exits cleanly)
resident_ball_228_windows_test.go:168: AC#1 READING: console up=0 absent=2; records created=-1 refused=2 stopped=-1 install=1
--- FAIL: TestAC228ResidentLegReportsAndBooksItsBall (3.53s)
```

⇒ **三件事一次读到**：
1. **失败支真能走到**（`refused=2` 那条记账在场、进程照走 D38(e)、`stopped=-1`＝没球时不误报收口）⇒ "只响亮、不致命"这一选型在真机上成立。⚠ 句末 `: <nil>` 是**我这发突变的产物**（把 err==nil 强推进那一支），**不是产码缺陷**，这一条读数按此打折。
2. **产码缺陷 F-6（本腿新发现，与突变姿势无关）**：`absent=2`——同一枚句子在控制台上**说了两遍**：`resident_ball_windows.go:96` 的 `fmt.Printf("wisp: %s\n", rb.verdict)` 与 `resident_windows.go:129`/`:131` 那条 boot 报告（`%s` 填的正是同一枚 `statusLine()`）各说一遍；成功支只在 `:111-112` 填 `verdict`、**不 printf** ⇒ **两支不对称**。撞红的是**该包自己的**"球姿态恰好说 1 次"不变式（`resident_ball_228_windows_test.go:78-81`）。**后果说死**：在没有桌面的机器上失败支就是自然路径，所以那一台上这枚用例会**因为重复、而不是因为缺陷**判红（§5 第 4 条具名 CI 那侧我未核）。修法（交下一程写腿，本腿不动产码）：摘掉 `:96` 那句 printf（boot 报告已带同一枚 `verdict`），或让两处共用一枚"只说一次"的门。
3. 我**不**把这一条读成 AC#7 不成立：两遍说的是**同一句真话**，AC#7 的射程是"三枚字符串说实话"，不是"每条只说一次"。

**winlive 档读数**（M5 在盘上时 `go test -tags winlive ./cmd/wisp -run 'TestLive228ResidentLegOwnsABallWindowOnTheDesktop'`）：

```
resident_ball_live_228_windows_test.go:121: AC#1 LIVE RED: the resident process reported NO ball ("this process has NO
    floating ball window" in its own console). This tier needs a desktop; the default tier covers the sentences.
--- FAIL: TestLive228ResidentLegOwnsABallWindowOnTheDesktop (3.28s)
```

⇒ **winlive 那一枚不是空转**：进程自报无球时它按 `:117-123` 判红——这正是编排者第 2 条要的**反面**证据（live 档确实盯着窗口，不靠句子放行）。
还原：`RESTORE cmd/wisp/resident_ball_windows.go bytes=9339 exact=True`。

### 全部突变还原后的复跑（每一发都回到绿）

默认档七枚（228 四枚＋127 三枚）：

```
--- PASS: TestAC228ResidentLegIsTheBallHost (0.01s)
--- PASS: TestAC228BallHostAnswersEveryGesture (0.01s)
    AC#1 READING: console up=1 absent=0; records created=2 refused=-1 stopped=3 install=1
--- PASS: TestAC228ResidentLegReportsAndBooksItsBall (3.26s)
    AC#1 READING: boot-time break exited clean; install=1 shutdown trail starts at 4 of 11 record(s); ball records created=2 stopped=3
--- PASS: TestAC228ExitRequestDuringBootStillLeavesThroughD38E (3.20s)
--- PASS: TestAC1ResidentLegInstallsItsLogListenerOnDisk (3.24s)
--- PASS: TestAC1ResidentLegOutlivesItsOwnLogFailure (2.85s)
--- PASS: TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink (3.23s)
ok  github.com/CarlosShao/wisp/cmd/wisp  15.862s
```

winlive 档一枚（本腿**自己**的真机读数，不复用实现腿的 hwnd 号）：

```
resident_ball_live_228_windows_test.go:147: AC#1 LIVE: hwnd=57609096 pid=27084 owns the ball window; hotkeys live 4/4:
    summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc live, …
live228: pid 27084 created and destroyed its WispBallWindow window through the D38(e) order
--- PASS: TestLive228ResidentLegOwnsABallWindowOnTheDesktop (3.03s)
```

⇒ **AC#1"真装起来了"这一层我有三枚独立读数**：桌面窗口按 pid 认领到恰好 1 枚（hwnd=57609096／pid=27084）、四枚热键全 live、退出后该 hwnd 不再是窗口且不再被那个 pid 认领（断言在 `:168-177`）。两次跑前桌面三连都是 `balldebug.exe` 0 枚／`wisp.exe` 0 枚。

---

## §5 没做完的（不许为空）

1. **winlive 档我只跑了 `TestLive228ResidentLegOwnsABallWindowOnTheDesktop` 一枚**（M5 在盘上一发＝红、还原后一发＝绿，读数在 §4）。**我没有在 winlive 档跑整包**，也**不拿默认档的绿说 winlive 过**。逐档归属写死：§0 的 148.168s 与 §4 每一发的 `--- FAIL`＝**默认档**；§1 第三把尺、§4/M5 的 live 那发、§4 末"还原后复跑"的最后一段＝**`-tags winlive` 档**。⚠ 该文件头（`:20-29`）自己规定它**不可与 `internal/ball` 的 winlive 守卫同屏并跑**，我照它单独 `-run TestLive228`，没跑 `./...`。
2. **没有量球路径的句柄／RSS 增量，也没入账任何性能数**。§3 那个 ≈235ms 只作"空窗存在"的旁证（两枚时间戳之差，单次读数），**不是**性能判定。理由：性能类判据要单包安静复跑才入账，且本票 AC#1/AC#7 两格不含性能判据（那是 AC#6 与 `internal/observe` 的地界）。
3. **`internal/observe` 的名册尺我只跑到"在册名"这一层**（§1 末）：**没有实测"常驻腿此刻活着的 resident 名有几枚"**。要跑得起来＝读一次真实 boot 的 `RosterReport`，本表没有那一枚读数；我只证了次序（`boot_windows.go:108` 的门跑在宿主之前）与名字（`ui-sta` 在 `ResidentNames`）。
4. **F-6 在 CI 上会不会响，今天零读数**。CI 侧我现读到的只有：`cmd/wisp` 的唯一 CI 落点是 **windows leg**（`.github/workflows/ci.yml:455` 那一步＋`scripts/wisp-cli-tests.sh:65` 的 GUARD 自陈"这是 cmd/wisp 门的 windows 半边"）；`-tags winlive` 在 `ci.yml` 里**零命中** ⇒ live 档只在这台机上跑过；`gh run list` 最近三枚（05:43／03:23／前夜 23:09 UTC）**全部早于本票四枚 commit** ⇒ **CI 还没跑过这轮码**。⇒ 留下的是一条**可验的预测**：若那台 runner 没有可建窗口的桌面，失败支会自然走到，`TestAC228ResidentLegReportsAndBooksItsBall` 会**因 F-6 的重复**红在 `:78-81`，而不是因缺陷。**我没有把它写成判语**，因为"CI 那台机器有没有桌面"我查不到（§6）。
5. **票 128"拒绝启动"那套定案我只读了 `cmd/wisp/dataroot_128_test.go` 的族名**，没逐枚读断言就写了"不同形"。这一句的分量按"未逐枚核"打折。
6. **球的手势我今天一次都没真按过**：十枚回调的"有执行者"我只被**静态形状**（§1 clause、M3 实测红）与 ball 自己的 winlive 套件覆盖；我这发的探针只读 boot／退出两段句子，**没有注入一次点击、一枚热键或一次托盘选择**。⇒ `recordBallGesture`／`recordTrayExit` 那几句**今天没有运行时读数**，这一层**不是 AC#1 的字面判据**（票面 AC#2 才裁"点下去有没有执行者"），我**不替 AC#2 背书**，只在此具名。
7. **AC#2／AC#3／AC#4／AC#5／AC#6 一格未裁**（按令不在射程）。§1 末那条"语义欠账"**不是**对 AC#2/#3/#4 的判定，只是说明 AC#1 的尺量到哪儿。
8. **F-1 与射程更正②的归属矛盾**（§2 末）我按未定义即停**上报而未裁**，需要编排者拍：这行注释到底算 228 名下欠账还是算 243 名下判定。

---

## §6 攻不动的地方（判不动就写判不动）

- **GitHub windows runner 有没有能建窗口的桌面**：**判不动**——这是那台机器的事实，不在仓里，我也没有第二次 CI 运行可查（本票四枚 commit 之后还没有 run）。我只能留 §5 第 4 条那枚预测。
- **`bootExit` 哑掉的那一发会不会掩盖某个真实退出路径**：**判不动**，因为要说它必须先证"存在一步能在 D38(e) 里卡住"，而 `internal/proc/shutdown.go` 的 10 步我没逐枚读超时来源（射程外的一枚新账，具名 F-4 已够；保护面已由 §3 甲/乙形实测）。
- **`hotkeys live %d/4` 用的是 `ball.DefaultHotkeys()`（`resident_ball_windows.go:79`）而不是 `[hotkey]` 配置**：句子读起来像"这台机器上活着的四枚"，实现腿自己具名承认"[hotkey]/[ball] 配置未读"。**判不动成不算是缺陷**：票 228 AC#1 的射程里没有"配置生效"这一格，而 `hotkeySummary`（`:189-199`）打印的是每枚槽位的**真实 Win32 结果**（`bd.Status`），不是配置愿望 ⇒ 我不判它谎，但把"配置未读"这一条留在 AC#1 之外、由**读配置那一格**（不属本票）负责。
- **票面 §1.2 那条"规格比仪器宽"的 emoji 缺口**：与本两格无关，未碰。
- **`design/assets/**` 被删导致 `internal/ball TestC21TableColourRowsMatchTokensCSS` 与 `internal/panel` 4 枚常红**：只记归因（别人在工作树里删的），**不判、不修、不引其结论**，且那两枚不是我地界。
- **`resident_ball_228_test.go` 的 AST 走查对 `-overlay` 结构性失明**：这一条**不是**判不动，是**台件约束**——我已按它选姿势（§4 未用 overlay，全部真落盘＋`git cat-file` 还原），写在这里是为了下一位别踩。

---

## §7 终态自证

- **起手名册（09-30 16:30，锚点 `5fcdb088`）**：`git status --porcelain -- internal cmd`＝**空**；桌面三连 `balldebug.exe`／`wisp.exe`＝**0 枚／0 枚**；`gh run list --limit 3`＝三枚全 `completed`（无本机 `slo-full` 在跑）。
- **终态名册（09-30 17:1x，最后一次量）**：`git status --porcelain -- cmd internal`＝**空**（逐字复跑，见本节的"复跑读数"）⇒ **终态等于起手名册**。共享树里这一条**不写"必须为空"**，只写"等于起手"。
- **每一发突变都能还原成绿**：五发（M1／M2／M3／M4／M5）的 `RESTORE <path> bytes=<n> exact=True` 逐枚在 §4；还原姿势一律 `git cat-file blob HEAD:<path> > <path>` ＋ `git diff --quiet` 判空（**M3 那发我自己插重复行犯过一次、按这把尺抓回来的，记录没抹**）。还原后七枚默认档用例＋一枚 winlive 用例**全绿**（§4 末逐名读数，`ok … 15.862s`／`ok … 3.081s`）。
- **复跑读数（本腿收尾自量，命令逐字）**：`git status --porcelain -- cmd internal` 输出为空；`gofmt -l cmd/wisp/resident_windows.go cmd/wisp/resident_ball_windows.go` 输出为空；五枚锚点 `grep -c` 各＝**1**（`rb := startResidentBall(rt.Registry)`／`	defer rb.stop()`／`signal.Notify(bootExit, os.Interrupt, syscall.SIGTERM)`／`OnDragEnd:       func() { recordBallDragEnd() },`／`^	if err != nil {$`）；收尾整包默认档复跑的 rc 与 `--- FAIL` 计数落在本节末"最后一发"那行。
- **产码零改动**：本腿全部落盘＝（a）`docs/evidence/s1/228-resident-ball-v1.md`（本表），（b）`.scratch/wisp/probes/228/v1/twobreak/{go.mod,main.go}`（独立模块的只读探针，**只建不删**，根模块 `./...` 编不到它）。
- **三枚冻结件／`PLAN.md`／`docs/specs/**`／`SLO.md`／`thresholds.go`／golden／`allowlist.txt` 零碰**；票面 `- [ ]`／`- [x]` 框**一枚未动**；`docs/reports/pending-and-issues.md` **未写一格**。
- **收尾整包复跑（本腿自己的最后一发，默认档，17:5x 前）**：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp -count=1` ⇒ **`ok github.com/CarlosShao/wisp/cmd/wisp 137.440s`／rc=0／`--- FAIL` 计数＝0**。与起手那发 148.168s 同侧（**两发都是安静单包**，差 10.7s 属争用型漂移，不入任何性能账）；⚠ 我全程**只用 `--- FAIL` 判红绿**，`build failed` 与 `0xc0000135` 各踩过/避过（§4/M4 与 §4 台件姿势）。
- **只 commit、绝不 push**；每一枚 commit 带显式 pathspec（`git commit --only docs/evidence/s1/228-resident-ball-v1.md`）。
