# 228-v1 — 非实现者对抗验收：票 228 的 AC#1（宿主存在性有生产点数）与 AC#7（三处打给人看／写进日志的字符串）

验收腿：`228-v1`（只读＋突变台件，零产码改动）。锚点自取：`git rev-parse --short HEAD` = **`5fcdb088`**。
被验收实现腿的四枚 commit：`90cf65b2`（骨架）→ `1739451d`（产码＋三枚测试）→ `15d2b2cd`（AC#7 文案）→ `39bcc958`（收尾）。
射程**只有两格**：AC#1 与 AC#7（含 09-30 的射程更正）。AC#2／AC#3／AC#4／AC#5／AC#6 本轮不在射程，**本表不裁它们，也不得被读成裁过**。

---

## §0 逐格 1:1 表（判语）

| 格 | 判语 | 我自己的命令（不是编排者给的基线） | 我的读数（时刻 09-30 16:3x–17:0x +08） |
|---|---|---|---|
| **AC#1** 宿主存在性有生产点数 | **成立**（附下面两枚具名条件，都不翻成"不成立"） | `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp -count=1`；尺 `grep -rln "wisp/internal/ball" --include=*.go . \| grep -v .scratch`；生产调用者逐行现读 `cmd/wisp/resident_windows.go:27-135` | 整包 **rc=0／零枚 `--- FAIL`／148.168s**（默认档，未带 `-tags winlive`）；尺＝3 命中：`./cmd/balldebug/main.go`、`./cmd/wisp/resident_ball_228_test.go`（尺自己在 `_test.go` 里的字面量，`resident_ball_228_test.go:43`）、`./cmd/wisp/resident_ball_windows.go`。宿主调用点是 **`resident_windows.go:118 rb := startResidentBall(rt.Registry)`**，在 `runResident()` 体内、`rt.RunEventLoop()`（`:132`）**之前**，而 `runResident()` 的唯一生产调用者是 **`cmd/wisp/main.go:64`**（`main()` 的 len(args)==0 分支，`:58`）。⇒ **不是"字段填了真值"，是"构造并跑了"**：跑的证据见 §1 第二把尺（真起 `wisp.exe` 无参数、读它自己写进磁盘的记账）与 §4 的 M1（摘掉那一行，两枚用例点名红）。<br>条件①＝见 §1 末"AC#1 的语义欠账"；条件②＝见 §3 的 `bootExit` 单发缓冲。 |
| **AC#7** 三处人读字符串 | **成立**（附一枚**具名阻塞欠账 F-1**：`cmd/wisp/main.go:8` 那行注释今天已是错事实，见 §2 末） | 逐处现读改后正文：`cmd/wisp/main.go:24-29`／`cmd/wisp/resident_windows.go:128-131`／`internal/proc/boot_windows.go:127-134`；对照 `git show 15d2b2cd` 与 `git show 1739451d -- cmd/wisp/resident_windows.go` 的改前原文 | 三处**都**已不含 `ticket 07`，也都不留将来式：<br>① `main.go:24-29`（`wisp -h` 正文）现文逐字 "hosts the floating ball window, its tray icon and its four global hot keys in this same process, then parks in an event loop that takes no task yet"＋把"托盘 Exit 无停机路径"与"球起不来只上报不停 boot"两条**负面**也写进去了 ⇒ 说实话。<br>② `resident_windows.go:131` 运行时行现文 "wisp: empty event loop running (no task is taken by this loop yet); %s (Ctrl+C exits cleanly)"，`%s` 由 `rb.statusLine()`（`resident_ball_windows.go:122-127`）填，两枚句子各在写它的那同一条语句里落成（成功 `:111-112`／失败 `:94`）⇒ 句子与事实同源，不能各自漂。<br>③ `boot_windows.go:132-134` 的 `slog.Info` 已拆成不带票号＋两条属性（`handled_here`／`bring_to_front`），逐字承认"这条循环自己没有窗口可提"⇒ 说实话。<br>⚠ **F-1 不撤销上面这一格**，理由是票面 ② 逐字「其余 25 处是注释……⛔ 不塞进本格（那会让 AC#7 永远勾不上），另立票 243」——见 §2 末的逐字对撞与我的裁法。 |

门禁四数（我自跑，不是抄实现腿）：`gofmt -l` 七枚被验文件＝**空**；`bash scripts/d22scan.sh`＝**PASS／clean**（`ban #8 cmd/ 69 枚`、`internal/ 473 枚`，⚠ 我第一次用 `go run ./tools/d22scan` 打的是**不支持的姿势**，根模块报 `does not contain package`，改用 wrapper 才有效——这条坑记给自己）；`go list -deps ./cmd/wisp` 见 §6。

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

### AC#1 的语义欠账（具名，不藏）

票面 AC#1 原文要求的是"≥1 枚来自 **`cmd/wisp` 那条真正会跑任务的腿**"。今天装球的是**常驻那条腿**，而它**仍不接任务**（`:131` 那句 `no task is taken by this loop yet` 是真话，`RunEventLoop` 只做激活记账）。这一格的口径来自**编排者 09-29 的 AC#0 裁定"走甲"**（票面第 3 行逐字「方向由编排者按 owner 长期原话裁死＝走甲（把审批门与任务通路接进**常驻**那条腿）」）＋ D2 拓扑（`PLAN.md:74/83-88`）：**被提名"真正会跑任务"的那条腿就是常驻腿**，`wisp run` 那一形在一发即退的语义下永远装不了球（票面"为什么不许就在 wisp run 里挂个球"三支）。
⇒ 我按票面 AC#1 的**字面尺**（grep＋生产调用者）判**成立**；但**AC#1 的立格目的**（"用户看得见球的那条进程＝干活的那条"）要到 AC#2／AC#3／AC#4 才闭合，**这一格不替它们背书**。编排者若要连目的一起勾，需要的是那三格的读数，不是这一格。

---

## §2 文案两支与失败路径（编排者第 2 条）

### 失败支今天有没有"测试真走到"？——**没有：只有突变走到过**

`startResidentBall` 的错误支（`resident_ball_windows.go:93-98`）由 `ball.New` 返回非 nil err 才进。`ball.New` 的 err 只来自 `b.sta.waitStarted()`（`internal/ball/ball_windows.go:173`）即 `createOnSTA` 的失败：`ensureFactories`／`registerBallClass`／`CreateWindowExW==0`／`newRenderer`／`addTrayIcon`（`ball_windows.go:182-243`）。本机有桌面 ⇒ 默认档与 winlive 档**都走成功支**。
**没有测试注入面**：`startResidentBall(reg *observe.Registry)` 直接把 `ball.New` 写死在函数体里（`:74`），没有 `func New` 字段、没有 seam、没有第二枚构造签名，所以**产码里不存在任何能让默认档走失败支的入口**。
⇒ 编排者这一条我**坐实为缺陷形状**：票 181 AC#7 栽的那一形（"有一条用例"被当成"失败支被证过"）在这里**部分是**真的——失败支**只被 M5 突变证过**（§4 末发），被自然运行**没证过**。
⚠ 这不是 AC#1/AC#7 的**判据缺件**：票面两格都没写"失败支要有自然覆盖"。但它是**文案两枚句子里较弱的那一枚的实证缺口**，我按具名欠账 F-2 交回，不擅自抬成不成立。CI 那一侧：`internal/ball` 的失败支若在**无桌面 runner** 上跑常驻腿就会自然走到——本表**未核 CI 环境是否有桌面**（判不动，见 §6）。

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
**缓冲满之后第二发信号会怎样**：`bootExit` 是 `make(chan os.Signal, 1)`（`:103`），首发包进缓冲、无人取；此后 `os/signal` 对所有已注册通道做**非阻塞发送**，通道满则**丢弃该发**——**不会**回落到默认动作、**不会**把进程打死。这正是"两发 Ctrl+C 不会二次硬杀"的原因，也是 §3 末那枚"硬杀能力被自己关掉"的代价来源。⇒ **这一处编排者怀疑的三个后果，实测形状是：不崩、不二次死、但第二发在 `bootExit` 这一路是哑的**；真正让进程离开的是 `sigCtx`，而它由 `RunEventLoop` 在读。**不判红**，理由：`default` 那一支已经把进程交给 `RunEventLoop`，此后一切退出请求由 `sigCtx` 承接；`bootExit` 哑掉的是**同一发的重复投递**，不是唯一出口。
⚠ 具名代价（F-4，交回不擅自修）：从 `signal.Notify` 起，**SIGTERM／os.Interrupt 的默认动作在整个 `runResident` 期间被抑制**，包括 D38(e) 正在跑的那段（`:69` 的 defer 与 `signal.Stop` 之间）。若 shutdown 序列里某步卡住，**再按 Ctrl+C 救不回来**（`bootExit` 已注册不读／`sigCtx` 的 `stop()` 已随 `RunEventLoop` 返回而执行）。今天 10 步序列有没有能卡的步＝本表未证（§6）。

### 接管窗口与正规 shutdown defer 之间会不会重复触发？

不会重复**触发关闭**：`rt.Shutdown` 只在 `:69` 那一枚 defer 里被调；`rb.stop()` 只在 `:119` 那一枚 defer；`sink.close()` 只在 `:67`。信号本身**不直接**触发任何一枚，它只让 `RunEventLoop` 返回。
LIFO 次序现读（注册序 `:67`→`:69`→`:105`→`:119` ⇒ 执行序 **`rb.stop()` → `signal.Stop(bootExit)` → `rt.Shutdown(false)` → `sink.close()`**）：`:109-113` 注释自称"球这一格做的是 D38(e) 第 2 步的活（`internal/proc/shutdown.go:18` hotkey + wake-word listening stops），跑在第 9 步 Job 关闭之前，且两边的日志都还能落到 sink"——三条我都现读对上：`shutdown.go:18` 逐字含 step 2 的措辞；`sink.close()` 确实最后；`Ball.Close()` 确实 join `ui-sta`（`internal/ball/ball_windows.go:935 <-b.sta.handle.Done()`）。
⚠ 一处**注释与代码的次序说法**要具名：`:113` 说 defer 是"registered last so LIFO runs it FIRST"——**成立**（`:119` 是本函数最后一枚 defer）。但**同一枚文件头 `:16-20`** 与 `:58-59` 把 `signal.Stop` 排在 `rb.stop()` 之后跑这件事没写（无害：`stop()` 不发日志之外的事件），我复核后判**不影响正确性**。

### 0xc000013a 那一形有没有真机证据？——**有，但被钉住的方式比实现腿的说法弱一档**

- **票 127 那枚 `TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink`（`resident_sink_nail_127_windows_test.go:543`）**：我读了全文断言（`:552-619`），**一字未放宽**：它要求 record 0 是 `residentEarlyResolverMsg`、install 记账**必须落在 index 1**（`:597`）、install 之后必须存在 `shutdown step` 轨迹（`:607`）、**文件最后一条必须是 shutdown 记录**（`:611-614`）、stdout 必须有 D38(e) 报告（`:615`）。`git log -- cmd/wisp/resident_sink_nail_127_windows_test.go` 末次 commit = **`9b5d64d8`（票 130）**，**早于本票四枚** ⇒ 实现腿**没碰过这枚钉**，也没有"为变绿改断言"这回事。⚠ 关键：它**并不断言信号落在 boot 窗口内**——它先等 sink 文件出现（`:552-555`）再 `breakToLoop()`，而 sink 文件在 `installLogSink`（`resident_windows.go:60`）就出现，`Notify` 在 `:104`，只差几条语句。所以这枚用例在修复**之前**红，是因为**球路径那 200+ms 里没有任何处理器**；修复之后它绿。它绿**不等于**"boot 窗口的退出被证过"。
- **真正冲这一形的是新用例 `TestAC228ExitRequestDuringBootStillLeavesThroughD38E`（`resident_ball_228_windows_test.go:184`）**：姿势与 127 同（sink 文件出现即 `breakToLoop`），断言是 exit code 必须 0（`:205-209`，`0xc000013a here means nothing was registered for the signal while the ball came up`）＋ install 记录＋shutdown 轨迹（`:211-220`）。
  ⚠ **这一枚的名字比它的断言强**：它**不核自己走的是哪一支**（既不要求 `:128` 那句 "arrived during boot" 出现，也不禁止 "empty event loop running" 出现）。⇒ 若某一发的调度让 break 落在 `RunEventLoop` 装上 `sigCtx` 之后，这一枚**照样全绿**，而它名字承诺的那一段根本没被走。突变 M4 之所以还打得中，是因为球路径本身给了约 200ms 的空窗（实现腿现测 268ms，我复跑 §4 读到的实际耗时见那发）；**空窗由球路径长度供养，不由任何一枚断言供养**。这是具名缺陷 **F-5**：判据对"摘掉 Notify"敏感，对"这一发究竟落在哪个窗口"**不敏感**，所以它是**弱钉**而非空钉（M4 实测红，见 §4）。
- **`resident_windows.go:85-102` 那段注释里给的数**（268ms／sink installed 16:05:20.7746／ball booked 16:05:21.0424）：来自实现腿自己那台机的一次日志读数，**本表不复现该毫秒数**（性能类判据要单包安静复跑才入账，本轮我不做性能判定），我只复现"**这段窗口里 Ctrl+C 会不会丢 D38(e) 轨迹**"这一**布尔**判据。

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

改法：`cmd/wisp/resident_windows.go:118` 由 `rb := startResidentBall(rt.Registry)` 改成 `rb := (*residentBall)(nil)`（锚点现读命中 1 次；`:119 defer rb.stop()` 保留，因为它对 nil 安全，`resident_ball_windows.go:141`）。⇒ **指名必须红**：`TestAC228ResidentLegIsTheBallHost`（clause 3，`:246`）＋`TestAC228ResidentLegReportsAndBooksItsBall`（created/refused 双缺，`:120-124`）。读数见下（第二次提交落入）。

### M2 — 摘掉 `defer rb.stop()`（实现腿自称会红，我独立复跑）

改法：删 `cmd/wisp/resident_windows.go:119` 那一行（锚点 `defer rb.stop()` 命中 1 次）。⇒ 指名：`TestAC228ResidentLegIsTheBallHost`（clause 4，`:266`）＋`TestAC228ResidentLegReportsAndBooksItsBall`（created≥0 而 stopped<0，`:153-155`）。读数见下。

### M3 — 缺一枚手势执行者

改法：删 `cmd/wisp/resident_ball_windows.go:90` 的 `OnDragEnd: func() { recordBallDragEnd() },`（锚点命中 1 次）。⇒ 指名：`TestAC228BallHostAnswersEveryGesture`（`:313-317`，missing 列表非空）。读数见下。

### M4 — 摘掉 boot 窗口的 `signal.Notify`

改法：把 `cmd/wisp/resident_windows.go:104` 的 `signal.Notify(bootExit, os.Interrupt, syscall.SIGTERM)` 注释掉（锚点命中 1 次；`:103`/`:105` 保留以维持编译与 `select` 形状）。⇒ 指名：`TestAC228ExitRequestDuringBootStillLeavesThroughD38E`（`:205-209` 的 exit-code 断言）＋票 127 的 `TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink`。读数见下。

### M5 — 让 `ball.New` 一定失败（失败支是否真走到／进程选型／文案说实话）

改法：`cmd/wisp/resident_ball_windows.go:93` 的 `if err != nil {` 改成 `if true {`（锚点命中 1 次）⇒ 成功支变成不可达，失败支**每次必走**，且**不动任何产码语义之外的东西**（跑完即还原）。⇒ 指名：默认档 `TestAC228ResidentLegReportsAndBooksItsBall` 应**仍绿**（走 refused 那一支、句子与记账同为"没有球"、teardown 不得在场）；`-tags winlive` 的 `TestLive228ResidentLegOwnsABallWindowOnTheDesktop` 应**红**（`:117-123` 明确在"进程自报无球"时判红）⇒ 这一枚同时证伪"winlive 档是空转"。读数见下。

---

## §5 没做完的（不许为空）

1. **winlive 档我只在 M5 那一发跑了 `TestLive228…` 一枚**（默认档整包与其余四发都在默认档）。**逐档归属**：§0 的 148.168s＝**默认档**；§1 第二把尺＝默认档；§1 第三把尺与 §4/M5＝**`-tags winlive` 档**。**我没有在 winlive 档跑整包**，也**不拿默认档的绿去说 winlive 过**（实现腿报的 hwnd／hotkeys 4/4 读数我没复现到同一个 hwnd 号，本表不引它当我的读数）。
2. **没有复现 268ms 这个数**，也没有量球路径的句柄／RSS 增量。理由：性能类判据要单包安静复跑才入账，且本票 AC#1/AC#7 两格**不含**性能判据（那是 AC#6 与 `internal/observe` 的地界）。
3. **`internal/observe` 的名册尺我只跑到"在册名"这一层**（§6）：我没验"常驻腿实际活着的 goroutine 数会不会超过 `ResidentBaseline = 6`"。要跑得起来＝读一次真实 boot 的 `RosterReport`，本表没有那一枚读数。
4. **CI 环境有没有桌面**没查 ⇒ §2 那句"无桌面 runner 上失败支会自然走到"是**推论**，不是读数，已在原文里标成推论。
5. **票 128"拒绝启动"那套定案我只读了 `cmd/wisp/dataroot_128_test.go` 的族名**，没逐枚读断言就写了"不同形"。这一句的分量按"未逐枚核"打折。
6. **AC#2／AC#3／AC#4／AC#5／AC#6 一格未裁**（按令不在射程）。§1 末那条"语义欠账"**不是**对 AC#2/#3/#4 的判定，只是说明 AC#1 的尺量到哪儿。
7. **F-1 与射程更正②的归属矛盾**（§2 末）我按未定义即停**上报而未裁**，需要编排者拍：这行注释到底算 228 名下欠账还是算 243 名下判定。

---

## §6 攻不动的地方（判不动就写判不动）

- **`bootExit` 哑掉的那一发会不会掩盖某个真实退出路径**：**判不动**，因为要说它必须先证"存在一步能在 D38(e) 里卡住"，而 `internal/proc/shutdown.go` 的 10 步我没逐枚读超时来源（射程外的一枚新账，具名 F-4 已够）。
- **`hotkeys live %d/4` 用的是 `ball.DefaultHotkeys()`（`resident_ball_windows.go:79`）而不是 `[hotkey]` 配置**：句子读起来像"这台机器上活着的四枚"，实现腿自己具名承认"[hotkey]/[ball] 配置未读"。**判不动成不算是缺陷**：票 228 AC#1 的射程里没有"配置生效"这一格，而 `hotkeySummary`（`:189-199`）打印的是每枚槽位的**真实 Win32 结果**（`bd.Status`），不是配置愿望 ⇒ 我不判它谎，但把"配置未读"这一条留在 AC#1 之外、由**读配置那一格**（不属本票）负责。
- **票面 §1.2 那条"规格比仪器宽"的 emoji 缺口**：与本两格无关，未碰。
- **`design/assets/**` 被删导致 `internal/ball TestC21TableColourRowsMatchTokensCSS` 与 `internal/panel` 4 枚常红**：只记归因（别人在工作树里删的），**不判、不修、不引其结论**，且那两枚不是我地界。
- **`resident_ball_228_test.go` 的 AST 走查对 `-overlay` 结构性失明**：这一条**不是**判不动，是**台件约束**——我已按它选姿势（§4 未用 overlay），写在这里是为了下一位别踩。

---

## §7 终态自证

- **名册**：终态 `git status --porcelain -- internal cmd` **等于起手名册**（09-30 16:30 那一次＝空 ⇒ 终态也应为空，读数在本节末逐字给）。共享树里这一条**不写"必须为空"**，只写"等于起手"。
- **产码零改动**：本腿全部改动=（a）`.scratch/wisp/probes/228/v1/**` 临时件（只建不删），（b）本表。每一发突变的还原逐枚打印 `RESTORE … exact=True`，且**还原后我再跑一次指定用例确认它回到绿**（§4 每发的最后一行就是这个）。
- **三枚冻结件／`PLAN.md`／`docs/specs/**`／`SLO.md`／`thresholds.go`／golden／`allowlist.txt` 零碰**；`- [ ]`／`- [x]` 框**一枚未动**；`docs/reports/pending-and-issues.md` **未写**。
- **锚点**：起手 `5fcdb088`；本表提交时的 HEAD 由 commit 自带，不在此复写。
- **只 commit、绝不 push**；commit 带显式 pathspec（`git commit --only docs/evidence/s1/228-resident-ball-v1.md`）。
