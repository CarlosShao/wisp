# 228-s2 写腿记录（**实现者自陈，⛔ 不是裁决表**）

- 起手锚：`git rev-parse --short HEAD` = **`aed8467b9`**（`dev`，2026-11-15 10:xx +08；该提交只把派单存档，无产码）。
- 射程：**票 228 切片 2 的四格** —— AC#4（行为类：`close()` 真等派生 goroutine）、AC#8（句子：`main.go` godoc）、AC#9（句子：球缺一句话说两次）、AC#10（句子：三处「我没有任务管线」自述）。
- **本件的读数是"我跑了什么"，不是"谁验收了"**：AC 框一枚未翻，勾归编排者／非实现者裁决腿。
- 写面：`cmd/wisp/**`（本波唯一写腿）、本件、`.scratch/wisp/probes/228/**`、票面 Progress log **只追加**。`internal/**`、`docs/PLAN.md`、`docs/specs/**`、`frontend/**`、`design/**` 一字未动（读数见 §5）。

---

## §0 四格 1:1（判语由裁决腿下，这里只交读数）

| 格 | 我改了什么（`aed8467b9` 起的差分） | 可判事实（我自己现跑的读数） |
|---|---|---|
| **AC#4 行为：`close()` 等齐派生的 goroutine** | `cmd/wisp/run.go`：`agentRuntime.close()` 在两条 `Cancel()` 之后**真等**两枚句柄（`rt.replyHandle`＝`approval-waiter`、`rt.reloadHandle`＝`watchdog`），新增具名预算 `joinOnCloseBudget228 = 500ms` 与 `joinedWithin228`；预算来历写在常量注释里（"一行答复的账＋审计收口"）。`cmd/wisp/approval_reply.go`：`replySurface` 加 `parked atomic.Bool` ＋ `noteParked`／`parkedInRead`，`runReplyLoop` 在**每次阻塞读的两侧**报自己在哪里（读法、语法、逐行 ctx 检查一字未动）。 | 判据是"**摘掉等的两行哪枚红**"（§1 M1）：`TestAC4CloseWaitsForBothGoroutinesItSpawned` 红三条（listener 还在写自己的行／tick 还活着／reply root 仍计 1 枚 pending）；`TestAC4CloseNamesTheParkedListenerInsteadOfWaitingOnIt` 红一条（少了那枚具名句）。**不是**"看起来更稳"。<br>正控（现树）：两枚 PASS，`TestAC4CloseWaits…` 0.12s（等到了那 120ms 的在飞行）、`TestAC4CloseNames…` 0.00s（**不花预算**）。<br>**谁 join 谁**逐名在 `close()` 的注释里；**超时从哪来**＝常量注释；GUI 腿那枚 ui-sta 线程不归这里 join，它由 `internal/ball` 的 `opts.Registry.Spawn("ui-sta","ball",…)`（`ball_windows.go:207`）派、`Ball.Close` 收，走的是 D38 允许的唯一一条派发口。<br>**分母**1→2：`reloadHandle`（`config_reload.go:142`）与 `replyHandle`（`approval_reply.go:464`）**一起等**，不是只补一枚。 |
| **AC#8 句子：`main.go` godoc** | **未改一字**——盘上已经是条件事实句。 | 现读（`cmd/wisp/main.go:5-13`）：`…then exits through the frozen D38(e) 10-step shutdown order. The floating ball is not a pending ticket: the resident leg of this same process assembles it at startup (startResidentBall in resident_ball_windows.go, called from runResident).`<br>使它成真的是 `9aaf601c8`（2026-11-08，`git log -S'The floating ball is not a pending ticket' -- cmd/wisp/main.go` 只有这一枚）。<br>**句子里的指认我自己核过**：`startResidentBall` 定义在 `resident_ball_windows.go:463`、唯一调用点在 `resident_windows.go:462`（`runResident` 体内）。⇒ 这一格交的是"读＋出处"，⛔ 不是"我又改了一遍"（改了才是造假话的反向 churn）。 |
| **AC#9 句子：球缺那句只说一次** | `cmd/wisp/resident_ball_windows.go`：**删掉**失败支里那枚 `fmt.Printf("wisp: %s\n", rb.verdict)`（作者从 2 处降到 1 处）；`cmd/wisp/resident_windows.go`：把启动报告最后一行抽成 `residentLoopReport228`（两条互斥分支逐字节等于原两枚 `Printf`，见注释），它就是那唯一作者；`cmd/wisp/resident_ball_seam_228_windows.go`（新）：`withBallNew228` 注入面，`panelHostHooks.newBall` 为 **nil 时走原样的 `ball.New`**。 | ① 盘上计数：`git show aed8467b9:cmd/wisp/resident_ball_windows.go \| grep -c 'fmt.Printf("wisp: %s\\n", rb.verdict)'` = **1**，现树 = **0**。<br>② 判据（§1 M3）：把那行**放回去**，`TestAC9BallRefusalIsSaidOnceByTheBootReport` 红：`the ball host itself printed the absence claim 1 time(s)`；现树该枚 PASS。<br>③ 失败支**第一次能被自然走到**：同一枚用例用 `withBallNew228` 返回 error 真走产码路径（verdict 赋值、`slog.Error` 记账、启动报告），读数：宿主自己说 **0** 次（`ballAbsentClaim` 计数）、启动报告两形各 **1** 次、且句子仍带 refusal 原因；用例 0.00s，**不需要没有桌面的机器**。<br>④ 有人问过"接缝会不会把 AC#1 的形状尺弄瞎"：会，所以产码那一支仍在本函数里写 `ball.New(opts)`（注释写明理由）；`TestAC228ResidentLegIsTheBallHost`＋`TestAC228BallHostAnswersEveryGesture` 现树 PASS（13 枚手势回调）。 |
| **AC#10 句子：三处「我没有任务管线」** | **未改一字**——三处盘上已是条件事实句。 | ① 文件头（`resident_ball_windows.go:21-42`）：`What this installment deliberately is NOT, re-derived by ticket 228 AC#10 because the ticket-228 wording ("the resident leg has no task pipeline and no microphone") became false on the disk two tickets later: 246-r1's startResidentTaskSource assembles the task pipeline in this very process…`——旧措辞是**被引用的过去**，不是断言；作者 `2327a7e7b`。<br>② 运行期那处（`:601` 的 `"gestures", "recorded only except the ones the assembly root hands an executor for: …"`）：同样条件句，同上作者。<br>③ `ballGestureWhy`（`:684`）：`no executor was handed to this ball host for this gesture; the gestures with one are the …`，作者 `fdee536b3`（票 290）。<br>**余量读数**：`grep -n 'has no task pipeline\|no microphone' cmd/wisp/resident_ball_windows.go` 命中 **2** 处，**两处都在引旧措辞的再推导注释里**（`:22`、`:657`），产码字符串里 **0** 处。 |

---

## §1 突变三发（每发：改后点名红 → 发内还原 → 三段哈希对撞）

| 号 | 变异点（内容锚） | 点名红（原文读数） | 哈希 before → mutated → after |
|---|---|---|---|
| **M1** | `cmd/wisp/run.go`：摘掉 `close()` 里**等**的那两段（reply 分支＋reload join） | `TestAC4CloseWaitsForBothGoroutinesItSpawned` **FAIL** 三条：`close() returned while the answer listener (approval-waiter) was still finishing its line`／`…while the config-reload tick (watchdog) was still alive`／`the reply root still counts 1 pending goroutine(s) after close()`；`TestAC4CloseNamesTheParkedListener…` **FAIL**：`close() left the answer listener unjoined WITHOUT recording it` | `400552bc…aaab` → `b5acd21f…7597` → `400552bc…aaab` ✅ |
| **M2** | `cmd/wisp/run.go`：只摘 `parked-in-read` 那一条 `case`（等还在，但那句具名话没了） | `TestAC4CloseNamesTheParkedListener…` **FAIL** 两条：`close() spent 501.0981ms (the whole join budget is 500ms) on a listener that told it it was parked in its read…`（**顺带证了那枚 case 买到的正是"不花预算"**）＋ `left the answer listener unjoined WITHOUT recording it`；审计行变成 `joined=false reason=budget-expired` | `400552bc…aaab` → `9b3827de…d0c3` → `400552bc…aaab` ✅ |
| **M3** | `cmd/wisp/resident_ball_windows.go`：把删掉的 `fmt.Printf("wisp: %s\n", rb.verdict)` **放回失败支**（接缝改完后按**最终树**重跑过一次，下面这组是最终树的读数） | `TestAC9BallRefusalIsSaidOnceByTheBootReport` **FAIL**：`the ball host itself printed the absence claim 1 time(s); the boot report is the single author of that sentence…` | `739b219a…3b0a` → `81e99125…d871` → `739b219a…3b0a` ✅ |

> 还原不是"我 echo 了已还原"：每发都用**发内 `Copy-Item` 原备份**，发后重算 sha256 与发前同串（脚本 `.scratch/wisp/probes/228/s2/mut/run-mutations.ps1`，读数 `.scratch/wisp/probes/228/s2/mut/*.log`）。

---

## §2 AC#4 的两句"说清楚"（编排者会问的）

1. **为什么不把这个监听也"叫醒"？** 叫得醒：`runReplyLoop` 的阻塞读在 `os.Stdin` 上，唯一真能唤醒它的动作是关掉/投喂 stdin。但 `wisp run` 的 stdin 是操作员的控制台，常驻腿另有 `taskSourceConsoleGoroutine` 在同一进程读它；从 teardown 里关 stdin 会把"退出"变成"顺手拆掉共享的输入口"，代价落在别人身上。⇒ 本波选**不假装等待**：等得住就等（in-flight 行），等不住就把**事实具名记进审计**，这不是"放弃"，是把"取消不等于等"这句注释换成一次现读。
2. **500ms 从哪来？** 一行答复收尾要花的钱是"账一写＋审计一追加"（`replySurface.handle` 的 `approval: REPLY ANSWERED` 与 risk sink 的落盘），毫秒量级；500ms 是它的宽限，不是任务预算（`proc.TaskWaitTimeout` 是排空任务的预算，不是这件事）。任务尺度的等待在这里**错两次**：工作不同，而且监听停在读上时再等也只是把"切断"推迟。

---

## §3 AC#5 正控（仍然红／仍然拒）

- `git diff cf3050401..HEAD -- internal/agent/approval` = **0 行**（切片 1 之后一字未动）；`git diff --numstat -- internal/` 本波 = **空**。
- `DecideFromPanel` 产码**真调用**现读 **2 处，全在 `internal/agent/approval/replies.go`**（`:408`、`:437`），`cmd/wisp` 产码 **0** 处，本波新增 **0**。
- 正控读数（定向发，现树）：`TestNativeHostSeamRefusesAPanelSourcedAllow` **PASS 2.19s**、`TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole` **PASS 2.14s**（同一张账上原生侧允许仍成立，所以"拒"不是全拒）。

---

## §4 门禁四数（我自己现跑）

| 门 | 命令 | 读数 |
|---|---|---|
| 格式 | `gofumpt -l cmd/wisp` | **空**，rc=0 |
| 静态 | `go vet -tags winlive ./cmd/wisp` | rc=**0** |
| D22 | `sh scripts/d22scan.sh` | **clean**，rc=0（正控 `runtests.sh` PASS=46 FAIL=0；扫 485 枚产码 Go；cmd/=80） |
| 路径长度 | `sh scripts/check-path-length-budget.sh` | **VERDICT GREEN**（tracked 20361、over-budget=68、roster=68、not-in-roster=0） |

**D38 形状尺（我自己现跑）**：本波差分 `+.*go func(` = **0** 命中、`+.*\.Spawn(` = **0** 命中、新增文件里 `go func(` = **0** ⇒ 一枚新协程都没有；球的 UI 线程仍是 `observe.Registry.Spawn("ui-sta","ball",…)`（`internal/ball/ball_windows.go:207`）。

### §4.1 整包 `go test ./cmd/wisp`（经 `scripts/wisp-cli-staged-tests.sh`，**一次一发**）

**一发，跑完了**（日志 `.scratch/wisp/probes/228/s2/s2-full.log`；载具自带 `-test.timeout=90m`，实际 `1525.00s` = 25m25s）：

| 数 | 读数 |
|---|---|
| `RUN` | **864** |
| `PASS` | **622** |
| `FAIL` | **2** |
| `SKIP` | **17** |

- **红名单对名（与派单给的已知名册逐名比）**：两枚红**都在名册里、且只有这两枚**——
  1. `TestAC4FocusReturnToPriorWindowGap33r5`（5.42s）＝派单点名的宿主环境红（前台被别的进程占着）。
  2. `TestAC246DevLegIgnoresTheTestTaskInjection`（35.81s）＝派单点名的宿主环境红（机主的常驻 `wisp` PID 49348 从 2026-11-14 16:23 起在跑，该枚自己红句写着 `a dev Wisp is already running in this session … Stop it and re-run.`）。
  名册里另外两条（`internal/ball` 1 枚、`internal/audio` 若干）**不在本包**——载具只跑 `./cmd/wisp`，故本发拿不到它们的读数，⛔ 不冒充。
- **对票 450 的基线逐数对撞**（那是本机唯一一次跑完的整包读数：`RUN=861 / PASS=619 / FAIL=2 / SKIP=17`）：本发 `RUN +3`、`PASS +3`、`FAIL ±0`、`SKIP ±0`，而本波**恰好新增 3 枚用例**（`TestAC4CloseWaits…` 0.12s、`TestAC4CloseNames…` 0.00s、`TestAC9BallRefusal…` 0.00s，三枚都在本发里 `--- PASS`）⇒ 增量对得上，没有"多出来的红"也没有"少跑的片"。
- **本波改动面的邻接枚在本发里也逐枚绿**：`TestAC228ResidentLegIsTheBallHost`／`…BallHostAnswersEveryGesture`／`…ResidentLegReportsAndBooksItsBall`(10.82s，真子进程)／`…ExitRequestDuringBootStillLeavesThroughD38E`(10.00s)、`TestAC11*` 四枚、`TestAC2Tray*` 六枚、`TestAC3Book…`／`…TheWaitingStateDrive…`、`TestReplyListener*` 两枚、`TestNativeHostSeam*` 两枚、`Test170*` 五枚——**全 PASS**。
- **⚠ 不拿"两发 FAIL 0"当覆盖证明**：这是**一发**，且它自己 `FAIL`（那两枚红）。覆盖的判据是上面"名册逐名对上＋对票 450 基线逐数对上"，⛔ 不是"红了但都是环境的"这种糊法。

> 一句如实的话：这一发**跑了 25m25s**，比票 450 记的 20–22m 长，多出来的时间是 `Test16r1PremortemModelSha256CorruptRefusesTheHandOff` 里那段**模型 re-ensure 走了真网络**（11:31:42 起 3 条 `Established …:443` 连接、日志静默约 7 分钟、进程 CPU 几乎不动，之后自行继续）；它最终 PASS，不是挂死。

---

## §5 交付清单（本波四改＋三新）

```
M cmd/wisp/approval_reply.go                          +39 -1
M cmd/wisp/resident_ball_windows.go                   +34 -3
M cmd/wisp/resident_windows.go                        +21 -3
M cmd/wisp/run.go                                     +87 -7
?? cmd/wisp/close_join_228_test.go                     （2 枚用例，AC#4）
?? cmd/wisp/resident_ball_refusal_228_windows_test.go  （1 枚用例，AC#9）
?? cmd/wisp/resident_ball_seam_228_windows.go          （AC#9 的注入面，仅测试调用）
```

「产码行为面只在 AC#4」可核：`git diff -- cmd/wisp` 里 `ball.New` 的唯一调用仍在 `startResidentBall` 体内（结构改成 `if hs.newBall != nil { … } else { ball.New(opts) }`，nil 默认 ⇒ 逐字节同前）；启动报告那两行是**搬家**，抽出函数的两个 `Sprintf` 形状逐字等于原 `Printf`。

---

## §6 没做／做不到（具名，⛔ 不粉饰）

1. **没有真机人手点托盘/按热键**：本波零 GUI 交互面（`winlive` 族），读数全来自产码路径＋注入面。
2. **`winlive` 族整族未跑**（载具不带 `-tags winlive`）：那一族里的常驻窗/前台断言仍归编排者约真窗。
3. **`TestAC246DevLegIgnoresTheTestTaskInjection` 的环境红**：本机仍有编排者的常驻球进程（PID 49348，2026-11-14 16:23 起）⇒ 该枚红在**桌面状态**而不是本波代码；本波没有修它、也不把它写成绿。整包那一发里它是 `--- FAIL: … (35.81s)`，红句原文 `a dev Wisp is already running in this session, so this leg handed its activation over and exited. Stop it and re-run.`（见 §4.1）。

---

## §7 撤销口令

一句话复原：`git revert --no-commit <AC#4 提交> <AC#9 提交>`（或手工 `git checkout aed8467b9 -- cmd/wisp/run.go cmd/wisp/approval_reply.go cmd/wisp/resident_ball_windows.go cmd/wisp/resident_windows.go`）＋ 删掉未跟踪的三件 `cmd/wisp/close_join_228_test.go`、`cmd/wisp/resident_ball_refusal_228_windows_test.go`、`cmd/wisp/resident_ball_seam_228_windows.go`。
