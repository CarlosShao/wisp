# 228-s1 写腿记录（**实现者自陈，⛔ 不是裁决表**）

- 起手锚：`git rev-parse --short HEAD` = **`27617392b`**（`dev`，2026-11-14 23:0x +08）。
- 射程：**票 228 切片 1 的四格** —— AC#2（托盘四个按钮各对各枚举）、AC#2 范围补充（托盘「退出」是 GUI 进程将来唯一停机源）、AC#11（退出一枚形：只记一次请求、立刻返回、让常驻腿自己的事件循环结束进程）、AC#3（`bookWaitingState` 推进 D43「等人」，⛔ 重做 AC#6a）。
- **本件的读数是"我跑了什么"，不是"谁验收了"**：AC 框一枚未翻，勾归编排者／非实现者裁决腿。
- 写面：`cmd/wisp/**`（本波唯一写腿）、`internal/ball/**`（编排者裁①批准的三处＋名册一行）、本件、`.scratch/wisp/probes/228/**`、票面 Progress log **只追加**。

---

## §0 四格 1:1（判语由裁决腿下，这里只交读数）

| 格 | 我改了什么（`file:line` 锚 `27617392b` 起的差分） | 可判事实（我自己现跑的读数） |
|---|---|---|
| **AC#2 托盘四项↔既有枚举** | `internal/ball/tray_windows.go`：`menuAllowOnce=5`／`menuReject=6`／`menuAllowAlways=7`＋`appendItem` 三枚标签，**追加在「打开面板」分隔线之后**；`internal/ball/ball_windows.go`：`Events` 尾部三枚 `OnTrayAllowOnce/OnTrayReject/OnTrayAllowAlways`（`OnDragEnd` 之后的键一个未动）＋ `wmAppTray` 三个 `case`；`cmd/wisp/resident_ball_windows.go`：`trayAnswerVerb`／`trayAnswerFunc`／`attachTrayAnswers`／`trayAnswerGesture`；`cmd/wisp/resident_approval_windows.go`：`trayAnswerExecutor` 把三个动词**逐一**交给 `surface.allow`＝`Replies.Allow`、`surface.reject`＝`Replies.Reject`、`surface.always`＝`replySurface.always`（**零新枚举**：`approval.Channel` 四枚值一枚未加，`git diff --numstat internal/agent/approval` = **0 行**）。 | ① `-test.run 'TestAC2TrayAllowOnceSpendsTheNativeGrant'` **PASS 1.82s**，审计原文 `approval: ANSWER-ALLOW corr=t228-tray-allow tool=fs.write route=native decision=allow`，且那张卡挡住的 `fs.write` **真的落盘**（用例读文件正文 `written from the tray`）。<br>② 同机制第二发 **PASS 1.93s**，`approval: REPLY ANSWERED corr="t228-tray-reject" tool=fs.write route="native/reject"`，目标文件 `os.Stat` 得 `IsNotExist`（拒绝真的没执行）。<br>③ 第三发 **PASS 0.00s**：只有 `always` 会记 `route="native/always"`，读数里有它、且**没有** `route="native/allow"`。<br>④ 落空那一形（F4）**PASS**：句子里有字面 `no pending card` 与「nothing was allowed, rejected or remembered」，且**不含**「已允许／已拒绝／ANSWERED」；slog 行含 `tray answer: no pending card` 与 `verb=allow-once`。<br>⑤ 未挂执行者那一形 **PASS**：`ball gesture arrived with no executor` ＋ `tray-allow-once`。<br>⑥ 托盘那侧的挂载／派发 **PASS**（AST 尺：id 1..4=1..4 与三枚标签逐字、`menuAllowOnce→[OnTrayAllowOnce]` 等三对**恰好一枚**）。<br>⑦ 名册行 **PASS**：`all 13 gesture callbacks answered by the resident host`。 |
| **AC#2 范围补充（退出是唯一停机源）** | 未新增代码：这一格落在既有 AC#11 执行者上（`residentBall.trayExit` → `rt.RequestStop`）。本波只做**复验**，并把「退出」与三枚答复项放进同一张菜单（`showMenu`）。 | 四枚定向钉枚枚绿（PASS）：`TestAC11StopBitIsBuiltAtBoot`、`TestAC11StopRequestDuringTheBootWindowIsNotSwallowed`（读数 `"stop request (tray Exit item selected in this process) arrived during boot; the event loop was never entered"`）、`TestAC11TrayExitExecutorIsMountedAtTheBallHost`、`TestAC11TrayExitTriggerWalksTheTenSteps`（读数＝进程内**真读满 1..10**：`scheduler-close > stop-hotkey-kws > cancel-task-roots > stop-audio > release-speech-sessions > destroy-panel-webview > flush-logs-close-db > free-os-memory > close-job-object > exit`，`reason` 点名托盘 Exit）。 |
| **AC#11 退出一枚形** | 未改实现（既有腿已落地）。这一格我交的是**"没有第二种形"的读数**。 | ① `cmd/wisp/resident_ball_windows.go` 里 `.Shutdown(` / `os.Exit` 生产命中共 **1** 处，是 **1112 行的注释**（逐字写着那两种形状被拒）；② 同一文件 `trayExit` 只调 `requestStop`（回程立刻返回），进程结束走 `RunEventLoop` 的 `<-stopSignal`；③ 上表四枚钉的 ⓑ-1 一拍读完十步＝「记一次请求、立刻返回、事件循环结束进程」这一形**有牙**（反形见 `228-r3c` 的 M1/M2，本波未重做）。 |
| **AC#3 `bookWaitingState` 推进「等人」** | `cmd/wisp/approval_always.go`：`bookWaitingState` 在**审计行之后**加一行 `rt.projectWaitingState(st)`（只推进"进入"方向）；同文件新增 `waitingStateSink`／`bindWaitingStateSink`／`projectWaitingState`；`cmd/wisp/run.go`：`agentRuntime` 加 `waitingMu`＋`waitingSink`（`run.go` 不 import statemachine，故用具名类型）；`cmd/wisp/resident_ball_windows.go`：`projectWaitingState`（先记 `waitingStates`，再 `rb.b.SetState(s)`）＋`waitingStatesSnapshot`；`cmd/wisp/resident_windows.go`：`src.run.bindWaitingStateSink(rb.projectWaitingState)`；`cmd/wisp/resident_approval_windows.go`：`ballCardUI.Prompt` 在 `cards.Record(p)` 之后、`SetState`／`SetBadge` 之前调 `ra.bookWaitingState("resident-ui-prompt")`（`Update` 那只在落地/开跑后记一次收尾）。 | ① `-test.run 'TestAC3BookWaitingStateDrivesTheBallWaitingState'` **PASS 1.98s**：该用例**自己不调** `bookWaitingState`，只把 `rb.projectWaitingState` 绑上去，然后让**生产路径**（`consoleApprovalUI.Prompt`）举一张真 L2 卡 ⇒ 第一枚到达投影的状态是 `AwaitingApproval`，审计行 `approval: WAITING-STATE state="AwaitingApproval" awaiting=true corr="t228-ac3-corr"` 同时在。<br>② 接线普查 **PASS**：`ballCardUI.Prompt` 调 `bookWaitingState` **恰好 1 次**；`resident_windows.go` 三行锚逐字在（`ra.bindRuntime(src.run)`／`src.run.bindWaitingStateSink(rb.projectWaitingState)`／`rb.attachTrayAnswers(ra.trayAnswerExecutor(src.answers))`）。<br>③ **没有第二条 UI 线程、没有新协程**：我的差分里 `+.*go func(`＝**0** 命中；投影落点是既有 `Ball.SetState`（`PostTask` 到那唯一 ui-sta 线程）。 |

---

## §1 突变四发（每发：改后点名红 → 还原 → 三段哈希对撞）

| 号 | 变异点（内容锚） | 点名红（原文读数） | 哈希 before → mutated → after |
|---|---|---|---|
| **M1** | `cmd/wisp/approval_always.go`：摘掉 `rt.projectWaitingState(st)` 那一行 | `TestAC3BookWaitingStateDrivesTheBallWaitingState` **FAIL**：`bookWaitingState drove nothing onto the ball projection; the producer only wrote its audit line (AC#3's defect)` | `306af320…9f09` → （变异树）→ `306af320…9f09` ✅ |
| **M2** | `internal/ball/ball_windows.go`：摘掉 `case menuAllowOnce:` ＋ `b.fire(b.opts.Events.OnTrayAllowOnce)` | `TestAC2TrayItemsAreMountedInTheBallPackage` **FAIL**：`wndproc dispatches menuAllowOnce to [], want exactly [OnTrayAllowOnce]` | `fa56c8ad…edd8` → （变异树）→ `fa56c8ad…edd8` ✅ |
| **M3** | `cmd/wisp/resident_approval_windows.go`：把"没卡"那条具名路径换成**静默空串** `return ""` | `TestAC2TrayAnswerWithNoPendingCardIsLoud` **FAIL** 四条：句子为空（不是 `no pending card`）、无 "nothing … remembered"、日志里没有 `tray answer: no pending card`、没有 `verb=allow-once` | `1f1441ce…0d74a` → （变异树）→ `1f1441ce…0d74a` ✅ |
| **M4** | `cmd/wisp/resident_windows.go`：摘掉 `src.run.bindWaitingStateSink(rb.projectWaitingState)` | `TestAC3TheWaitingStateDriveIsWiredIntoTheShippedPath` **FAIL**：`the shipped resident assembly is missing the line "src.run.bindWaitingStateSink(rb.projectWaitingState)"` | `52d45c9d…b5c3` → （变异树）→ `52d45c9d…b5c3` ✅ |

> 还原不是"我 echo 了已还原"：每发都用**发内 `Copy-Item` 原备份**，发后重算 sha256 与发前同串；`git status` 里除本波的九件（见 §5）无第十件。

---

## §2 编排者那两问（接线到底接上没有，说清楚）

问：`runSpec.replyVeto` 产码赋值点＝0、`run.go` 某行给门递空 `approval.NewChannels()` —— 接上了还是没接上。

答（三句，都可现跑）：

1. **托盘三枚答复项走的**不是 `runSpec.replyVeto` **这条路**，也**不需要**它。它们走的是 `src.answers`＝`run.newReplySurface(src.root, ra.consoleVetoChannel(), false)`（`cmd/wisp/resident_task_source_windows.go`，本波把这只 surface 从"只在有控制台时建"改成"**有 runtime 就建**"），而 `newReplySurface` 把 `rt.liveCards`（常驻腿**开工时就**把门与两个路由器绑好的那张账：`resident_approval_windows.go:580-586` 的 `Attach(HostBinding{Gate, VetoChannel: ChannelEsc, NativeSource, PanelSource})`）递进 surface。⇒ `allow/reject/always` 三个动词花的是**常驻门自己的**一次性令牌。
2. `runSpec.replyVeto` 的产码赋值点现读 **1 枚**（`cmd/wisp/main.go:247`，只在 `reply != nil` 的**控制台**那一支），零产码写者不是事实——但**本波没有加第 2 枚**，也没改它。`approval.NewChannels()` 的产码出现点现读 **2 枚**：`run.go:1066`（`wisp run` 的门）与 `resident_approval_windows.go:556`（常驻门**构造时**）。两处"空"都是**构造时的空**，装载走诚实声明口：常驻腿在球真加载了 Esc 钩子时 `SetLoaded(ChannelEsc,true)`（`resident_approval_windows.go:761`），本波 `newReplySurface(..., ra.consoleVetoChannel(), ...)` 再按同一事实声明一次（`escLoad` 为假就传空，`newReplySurface` 里 `if vetoChannel != ""` 才装）。
3. 结论一句：**接上了，接在常驻门那张账上；不是"在 `wisp run` 的空通道上假装接上"**。判据＝§0 AC#2 的 ①② 两发：真卡的令牌真被花掉、被拒后真没执行——如果接的是一条空通道，这两个读数都拿不到。

---

## §3 AC#5 正控（仍然红／仍然拒）

- 生产面 `DecideFromPanel` 的调用点现读 **2 处，全在 `internal/agent/approval/replies.go`**（`decide(viaPanel=true)` 的拒绝路与 `PanelAllow` 的拒绝路）；`cmd/wisp` 产码里 **0** 处；本波**没有**新增一处。`git diff --stat -- internal/agent/approval` = **空**。
- 正控读数：`-test.run 'TestNativeHostSeamRefusesAPanelSourcedAllow'` **PASS 1.84s**；同一发里 `TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole` **PASS 1.62s**（正控的反面：同一张账上原生侧允许仍然成立，所以"拒"不是全拒）。
- 另：托盘三枚项**没有**任何一枚能走面板路——`PanelAPI` 没有 `Allow` 方法（`replies.go` 的 `PanelAllow` 只有拒绝支），`DecideFromPanel` 的调用面本波逐字未动。

---

## §4 门禁四数（我自己现跑）

| 门 | 命令 | 读数 |
|---|---|---|
| 格式 | `gofumpt -l cmd/wisp internal/ball` | **空**，rc=0 |
| 静态 | `go vet -tags winlive ./cmd/wisp` | rc=**0** |
| D22 | `sh scripts/d22scan.sh` | **clean**，rc=0（正控 `runtests.sh` PASS=46 FAIL=0；扫 484 枚产码 Go） |
| 路径长度 | `sh scripts/check-path-length-budget.sh` | **VERDICT GREEN**（tracked 20342、over-budget=68、roster=68、not-in-roster=0） |

### §4.1 整包 `go test ./cmd/wisp`（经 `scripts/wisp-cli-staged-tests.sh`，一次一发，串行两发）

**两发都没能跑完，且两发都不是断言红**——本机在**真面板窗外那一区**整进程死亡（`228-r3b` 记过同一族"真窗区整进程死亡/非确定性"）：

| 发 | 读数 | 死在哪 |
|---|---|---|
| 第 1 发 | `--- PASS` **339** 枚、`--- FAIL` **0** 枚；`FAIL github.com/CarlosShao/wisp/cmd/wisp 416.00s`（载具按 rc 打的行） | 死在 `Test145ArgsJSONComesFromTheBookedForensicsRow`（`cmd/wisp/panel_run_145_test.go`）**执行到一半**：该枚无 `--- 结果` 行 |
| 第 2 发 | `--- PASS` **334** 枚、`--- FAIL` **0** 枚；`FAIL … 403.00s` | 死在 `TestAC4PriorFocusSurvivesARefusedPanelSample`（同为真 WebView2 窗那一区）**执行到一半** |

⇒ 两发都只覆盖到文件序 `p` 区（`panel_*`）。**`r`..`v` 区（`resident_*` / `ticket*` / `voice_*`）在两发里根本没开始跑**，所以本波对那一区的覆盖改由**定向发**给出，⛔ 不把"整包跑过"写成绿：

`sh scripts/wisp-cli-staged-tests.sh -test.run '…'`（`-test.v`，exit 0，`ok github.com/CarlosShao/wisp/cmd/wisp 26.00s`，**19 PASS / 0 FAIL**，日志 `.scratch/wisp/probes/228/s1-sweep-final.log`）：
`TestAC228ResidentLegIsTheBallHost`(0.05s)／`TestAC228BallHostAnswersEveryGesture`(0.06s)／`TestAC228ResidentLegReportsAndBooksItsBall`(**9.48s，真子进程拉起**)／`TestAC228ExitRequestDuringBootStillLeavesThroughD38E`(**8.58s**)／`TestAC2TrayAllowOnceSpendsTheNativeGrant`／`TestAC2TrayRejectRefusesOnTheNativeRoute`／`TestAC2TrayAllowAlwaysRoutesIntoTheWideningVerb`／`TestAC2TrayAnswerWithNoPendingCardIsLoud`／`TestAC2TrayAnswerWithoutAnExecutorSaysSo`／`TestAC2TrayItemsAreMountedInTheBallPackage`／`TestAC3BookWaitingStateDrivesTheBallWaitingState`／`TestAC3TheWaitingStateDriveIsWiredIntoTheShippedPath`／四枚 `TestAC11*`／`TestNativeHostSeamRefusesAPanelSourcedAllow`／`TestTicket21S2BadgeWiringCensus`／`TestTicket21S2BadgePostsAreOrderedStateThenDepth`。

---

## §5 交付清单（本波九件＋一件测试）

```
M cmd/wisp/approval_always.go                 +51 -0
M cmd/wisp/resident_approval_windows.go       +170 -0
M cmd/wisp/resident_ball_228_test.go          +6 -0     （名册只扩列表）
M cmd/wisp/resident_ball_windows.go           +144 -2
M cmd/wisp/resident_task_source_windows.go    +16 -1
M cmd/wisp/resident_windows.go                +24 -0
M cmd/wisp/run.go                             +10 -0
M internal/ball/ball_windows.go               +22 -1    （删的 1 行是 OnDragEnd 原文搬家）
M internal/ball/tray_windows.go               +28 -0
?? cmd/wisp/resident_tray_answer_228_windows_test.go  （8 枚用例）
```

「只加不改」可核：`internal/ball/tray_windows.go` **0 删除行**⇒ id 1..4 与四枚既有标签逐字未动；`Events` 里既有四枚 `OnTray*` ＋ `OnDragEnd` 逐字未动（差分只有 `+` 与搬家）。

---

## §6 没做／做不到（具名，⛔ 不粉饰）

1. **一次真的鼠标点托盘菜单**：需要通知区域里一枚真图标＋人手（`//go:build winlive` 那一族），本波没跑。替代读数＝§0 AC#2 的 ⑥（AST 尺：菜单真建三枚、wndproc 真派发到三枚回调）＋ ⑦（宿主真把这 13 个手势都挂上）。**这不是"跑过了"**。
2. **AC#11 的 ⓑ-2 真机格**（`TestLive228R3TrayExitStopsTheShippedProcessWithCodeZero`）仍拦在 `resident_tray_stop_live_228_windows_test.go:184`（本机开不出可注入鼠标点选的交互托盘世界），归编排者约真窗；本波未重跑该枚，也未把它写成绿。
3. **两枚"已知环境红"的现况（与本波无关，未修未删）**：
   - `TestAC4FocusReturnToPriorWindowGap33r5`：`228-r3b` 记它是"前台被别的进程占着"而红；**本波第 2 发整包里它 PASS（1.27s）**⇒ 它是**环境相关**而不是恒红，本波既没修它也没把它写成绿。
   - `TestAC246DevLegIgnoresTheTestTaskInjection`：由编排者的常驻球进程（PID 49348）在跑而红；**本波两发整包都没跑到它（`r` 区）**，故另跑一发定向读数（`.scratch/wisp/probes/228/s1-known-red-ac246.log`）：`--- FAIL: TestAC246DevLegIgnoresTheTestTaskInjection (36.56s)`，红句原文 `resident_task_source_246_windows_test.go:484: AC#7 RED (desktop state, not our code): a dev Wisp is already running in this session, so this leg handed its activation over and exited. Stop it and re-run.` ⇒ **红在桌面状态、不在本波代码**。

---

## §7 撤销口令

一句话复原：`git checkout 27617392b -- cmd/wisp internal/ball` ＋ 删掉未跟踪的 `cmd/wisp/resident_tray_answer_228_windows_test.go`（本波全部产码都在 `cmd/wisp/**` 与 `internal/ball/**` 两处，`internal/agent/approval`／`frontend`／`design`／`docs/specs` 一字未动）。
