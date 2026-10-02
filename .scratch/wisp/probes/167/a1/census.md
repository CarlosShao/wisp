# 167-a1 — 面板四枚小出口＋崩溃自救：Go 侧真源取证（只读腿）

> 本文件由只读取证腿 `167-a1` 产出。**零产品码改动**；唯一可写文件即本文件。
> 票面：`.scratch/wisp/issues/167-four-small-outputs-the-panel-needs-occupancy-queue-stop-draft-plus-crash-self-rescue.md`（未改动，AC 框一枚未碰）。
> 搜索根一律显式：`cmd internal tools docs scripts .scratch`。`frontend/**` 与 `design/**` 未读、未引用、结论里不出现其行号。
> 未跑 `go test` / `go build`（两条写腿在飞）。允许动作只有 `grep` / `sed` / `ls` / `git log` / `go list`。
> 骨架先落（本轮只有 §0 与各区标题），§1–§6 由后续 commit 逐节补齐。

## 0. 起手锚（逐字读数）

`date -Iseconds`：

```
2026-10-02T09:29:29+08:00
```

分支（`git rev-parse --abbrev-ref HEAD`）：

```
dev
```

`git log -1 --format='%H%n%ad%n%an%n%s'`：

```
f5f9cc3433385c13e24cb2a7d8054d8210f99cde
Fri Oct 2 09:27:56 2026 +0800
CarlosShao
census(114-a2)：§1 调用者名册＋§2 逐格三态交齐——ParseComposerRequest 生产调用者=1(composer_dispatch.go:155)且 Handle 有两枚产码听众；未勾实测 9 枚不是 8 枚；248-r1 一格都没顺手做掉，反而把 AC#8 的覆盖面削薄(config.* 不带 panel. 前缀、两把尺都看不见)
```

`git status --porcelain internal/panel cmd/wisp`（起手读数，逐字）：

```
 M cmd/wisp/run.go
?? cmd/wisp/firstrun.go
?? cmd/wisp/firstrun_198_test.go
?? cmd/wisp/firstrun_acl_198_windows_test.go
```

起手读数解读（仅记录事实，非结论）：`cmd/wisp` 侧有在飞的写腿产物（`run.go` 已改 + 三枚 `firstrun*198*` 新件），与派单说明的 `198-r1` 占 `cmd/wisp` 相符；`internal/panel` 干净——本腿读到的 `internal/panel` 内容即 HEAD 的内容。

## 1a. 起手锚的世代更正（先落这一节，因为它决定下面所有行号的可信度）

本腿开工时 HEAD＝`f5f9cc34`（§0 已逐字记录）。**读数的途中树前进了**：`git log f5f9cc34..HEAD` 现量 **16 枚**，其中动到 `cmd/wisp/*.go` 的是 `613606c0`〔198-r1〕（`run.go` ＋ 三枚 `firstrun*198*`），动到 `internal/ball/*.go` 的是 `494487df`〔33-r8b〕（`sta_windows.go`）。当前 HEAD＝`5fad9a00`，`git status --porcelain internal/panel cmd/wisp` **此刻为空**（复尺时刻 `2026-10-02T09:41:15+08:00` 之后另取一发，见本节末）。

处置（不靠"应该没漂"）：下面**每一枚 `cmd/wisp/run.go` 行号都在新 HEAD 上重新 `grep -n` 复认过**——
`run.go:700` `Verdicts: rt.liveVerdicts`／`run.go:725` `Out: rt.bookPanelSnapshot`／`run.go:735` `rt.ui.publish = rt.publishPanelSnapshot`／`run.go:998` `consoleSink{…}`／`run.go:1056`＋`:1062` `loop.Budgets().PromptTotal`／`run.go:1109` `agent.NewSpiller(…, loop.Budgets())`／`run.go:1012` `ContextWindow: ep.ContextWindow`／`run.go:1099` `loop.RunAsync(ctx, task)` ＝**逐枚与起手读数一致，零漂移**。
`internal/panel/**` 与 `cmd/wisp` 其余文件在该窗口内**无任何 .go 改动**（尺：`git diff --name-only f5f9cc34..HEAD | grep -E '\.go$'` 现量 8 行，全在 `probes/**`、`cmd/wisp/firstrun*`、`cmd/wisp/run.go`、`internal/ball/sta*`），故下面那些 `file:line` 与起手读数同版。
另：本腿**未读**、也**未引用** `internal/ball/sta_windows.go`（33-r8b 刚动过它），所以那一处改动进不了本表任何一格。

## 1. 名册（五枚出口 × 四问）

> 三态口径（本仓已定过的分形）：**非测试生产者**＝产码里有作者；**仅测试生产者**＝只有 `_test.go` 在填；**盘上没有**＝连字段都没有。
> 面板通路只认两枚结构体：`internal/panel/composer.go:57` `Snapshot`（顶层 6 枚字段）与 `internal/panel/composer.go:235` `ComposerState`（11 枚字段）；泵能读什么只认 `internal/panel/pump.go:121` `PumpSources`（12 枚槽）。

### 1.1 占用（上下文／额度用了多少 ÷ 总共多少）

| 问 | 读数 |
|---|---|
| ① Go 侧生产者 | **分母（非测试生产者，但没人往面板方向读）**：`internal/agent/budgets.go:54` `Budgets.ContextWindow` ← `BudgetsFor()` `internal/agent/budgets.go:86-91`，导出读口 `internal/agent/loop.go:253` `Loop.Budgets()`。该读口**有 3 枚产码调用者**，全在 `cmd/wisp/run.go:1056`／`:1062`／`:1109`，**三枚都只取 `.PromptTotal` 或整包喂给 spiller，没有一枚把 `.ContextWindow` 递出 `internal/agent` 之外**。另一条分母来源是配置叶：`internal/panel/config_handlers.go:60` `FieldModelContextWindow`（`config.set` 可写）→ 解析后的落点 `cmd/wisp/run.go:1012` `ContextWindow: ep.ContextWindow`。**分子（非测试生产者，同样不出包）**：`internal/agent/compress.go:110` `Compressor.TotalTokens(hist)`（调用者只有 `compress.go:120/:175/:212/:248` 与 `internal/agent/loop.go` 自己）；`internal/agent/guard.go:132` `Guard.TokensUsed()`（产码调用者 `internal/agent/guard.go:151/:210/:212`、`internal/agent/loop.go:381`，**全在包内**）；`internal/agent/prompt.go:181` `Assembler.TotalTokens()`（包内）。历史读口 `internal/agent/loop.go:256` `Loop.History()` 是导出的，**产码调用者只有 `internal/agent/loop.go:392/:560/:566`，`cmd/wisp` 零枚**（尺：`grep -rn "\.History()" cmd internal tools --include=*.go | grep -v _test.go:`）。 |
| ② 到面板的通路 | **没有字段**。`Snapshot`（`internal/panel/composer.go:57-92`）＝`pending/results/composer/generatedAt/instructions/tasks` 六枚，无占用；`ComposerState`（`:235-270`）＝11 枚里最近的两枚是 `CurrentModel`/`ModelKnown`（`:255-256`）与 `Credential`/`CredentialKnown`（`:268-269`），**没有 used/total 的任何一格**；`PumpSources`（`internal/panel/pump.go:121-188`）12 枚槽里无占用读口。⚠ 唯一**已在产码里被读的"占用"是另一件事**：`internal/panel/subagent_roster_197.go:152-153` `inFlightSlots`/`poolCap`＝**子代理池槽位**，不是上下文；它已有真源（`cmd/wisp/panel_pump.go` 的 `taskRosterState`，经 `run.go:724` `Tasks: rt.taskRosterState` 接上）。**最近的现成来源**＝`internal/agent/loop.go:253` `Loop.Budgets()`（分母）＋`internal/agent/compress.go:110` `TotalTokens(History())`（分子），两枚都已在产码、都已在 `cmd/wisp/run.go` 的同一枚 `execute()` 里够得着（`run.go:1023` 造出 `loop`，`run.go:1056/:1062/:1109` 三枚既成调用者就是证据）。 |
| ③ 动作方向 | 不适用（只读维度，无入向方法）。 |
| ④ 判据 | **正向钉：无**。**负向钉：无**。尺：`grep -rn -iE "occup|占用|contextwindow|token" internal/panel/*_test.go cmd/wisp/panel*_test.go` 的命中全是无关形状（`go/token` 包、`FieldModelContextWindow` 的**写面**名册 `internal/panel/config_route_248_test.go:272/:296`、C21 设计令牌 `internal/panel/frontend_hygiene_test.go:202/:217`、`internal/panel/instructions_200_test.go:56` 的 `BudgetTokens` 夹具）。⛔ 我没有找到任何一枚用例钉"分子或分母未知 ⇒ 出口必须是具名未知态"（＝票面 AC#2 要的那形）。 |

### 1.2 排队（第几条／队列长度）

| 问 | 读数 |
|---|---|
| ① Go 侧生产者 | **审批队列：非测试生产者，两枚都在**：深度 `internal/agent/approval/queue.go:133` `Queue.Depth()`，产码调用者 `cmd/wisp/approval_reply.go:377`（`"队头没有待审批的卡片（当前深度 …）"`）；序号 `internal/agent/approval/queue.go:268` `position()`，**由 `internal/agent/approval/pending_read.go:118` 在 `LiveApprovals()` 里逐行填进 `LiveApproval.Position`（字段声明 `internal/agent/approval/pending_read.go:49`）**——这条是真源，不是测试形状。**用户消息队列：盘上没有**。常驻腿**明写拒绝排队**：`cmd/wisp/resident_task_source_windows.go:411-414`（`if src.running { … return errors.New("本进程此刻已有一发任务在跑：常驻腿一次只接一发，等它结束再提交") }`），单槽规则的理由逐字在 `:396-400`。环内那条"忙时插话"的通道**只有方法没有作者**：`internal/agent/loop.go:275` `Loop.Steer()` 导出、往 `l.steer`（`:188`）追加、`:390` `drainSteering()` 消费，⚠ **`Steer` 的调用者＝零枚**（尺：`grep -rn "Steer" cmd internal tools --include=*.go | grep -v _test.go:` 只命中 `SteeringEnabled` 那三枚配置/选项行与定义本体，**无一处调用**；连测试都没调）。另有一枚序号：`cmd/wisp/resident_task_source_windows.go:159` `src.seq`，`:416-417` 自增，**只用于 goroutine 名与审计行**（`:422` root 名、`:425` `Spawn` 名、`:437` 日志 attr），从不外递。 |
| ② 到面板的通路 | **序号：字段不存在，且现成的真源在最后一跳被丢掉**。`internal/panel/pump.go:77-91` `NativeVerdict` 的字段名册＝`CorrelationID/Tool/Args/CallChain/Level/RulesHit/Reason/SessionOverrideBlocked`，**没有 `Position`**；`internal/panel/approval.go:39-56` `ApprovalCardView` 同样没有；产码映射 `cmd/wisp/panel_pump.go:58-74` `liveVerdicts()` 逐字段抄 `it.Decision`，**`it.Position` 一个字节都没抄**。⚠ 那里**已有一段"哪一枚字段没有产源"的自述**，但它点名的是 `CallChain` 不是 `Position`（`cmd/wisp/panel_pump.go:50` 注释＋`internal/panel/pump.go:82-87`）⇒ **这条丢弃至今无人登记过**。**长度：间接在场上**——`Snapshot.Pending` 是数组（`internal/panel/composer.go:58`），页面数得出长度，但没有一具"深度"标量；`PumpSources`（`pump.go:121-188`）无长度读口。 |
| ③ 动作方向 | 取消排队＝**入向方法不存在**。白名单现量 6 枚：`internal/panel/bridge.go:42-45` 四枚 `panel.*` ＋ `internal/panel/bridge.go:66-67` `config.get`/`config.set`；守卫 `internal/panel/bridge.go:146-148` `knownComposerMethod` 的 case 行；派发表 `internal/panel/composer_dispatch.go:176-208` 六个 case。⛔ 别把"入向有了一条线"读成"所有口都接好了"：这条线**今天确实有产码听众**（`cmd/wisp/panel_host_windows.go:551` 与 `cmd/wisp/panel_inbound.go:163` 两枚 `disp.Handle`），但**受理的处理器只有两枚**——`cmd/wisp/panel_inbound.go:228` `newComposerDispatchChain` 产出的字面量（`:271-281`）里 `Workspace`/`Attachment`/`Message` **逐枚写死 nil**（注释具名归属：workspace＝票 186、attachment＝票 92、message＝票 35），只有 `Mode` 与 `Config` 接了真处理器。常驻那条腿走的是**同一枚链**：`cmd/wisp/panel_resident_windows.go:170` `newResidentComposerDispatch` 直接 `return newComposerDispatchChain(dataDir, auditf, residentPanelActor)`，而它由 `cmd/wisp/resident_windows.go:142` 调用。 |
| ④ 判据 | **正向钉：一枚，且钉的是"卡阵列"不是"第几条"**——`cmd/wisp/panel_pump_test.go:141` `TestRunBooksWithASnapshotOfItsLiveQueue`：真跑一发越允许根的工具、等 L2 卡进队列、采样快照、断 `len(live.Pending) == 1`（`:194`）并与该 run 自己触发的账本记录对账。**它没钉 `Position`**。**负向钉：无**（没有一枚用例说"取消排队不许把后面的挤掉"）。尺：`grep -rn "^func Test" internal/panel/*_test.go cmd/wisp/*_test.go | grep -iE "queue"` 除窗口消息队列那形（`cmd/wisp/panel_resident_windows_test.go:523`）外只命中上面这一枚。 |

### 1.3 停止（取消按钮）

| 问 | 读数 |
|---|---|
| ① Go 侧生产者 | **动词有、链子半截**。`internal/agent/control.go:22` `ControlStop = "stop"`／`:24` `ControlCancel`；解析器 `internal/agent/control.go:51` `MatchControl`（词表 `:36-42`，`"停"`／`"取消"`）；执行体 `internal/agent/loop.go:524-537` `defaultControl`——它真的 `root.Cancel()` 并回 `"已停止当前任务"`（`:531-532`）。⚠ **入口只有一枚**：`internal/agent/loop.go:341` `if verb, ok := MatchControl(input); ok` 位于 `run()` 开头，⇒ **控制词必须以"新提交一发输入"的形态到达**；而常驻腿在任务运行中**拒绝第二发**（`cmd/wisp/resident_task_source_windows.go:411-414`，见 1.2）。**宿主自定义的 ControlHandler 从未接线**：`internal/agent/loop.go:158` `Options.Control ControlHandler` 的产码赋值＝零枚（尺：`grep -rn "Control:\s" cmd internal tools --include=*.go | grep -v _test.go:` ⇒ 只命中 `internal/llm/anthropic/request.go:467/:474/:482` 的 `CacheControl`，与这枚无关）。**异步入口的取消口有定义无调用者**：`internal/agent/loop.go:300` `func (t *RunningTask) Cancel() { t.root.Cancel() }`，产码调用者零枚。**名册级停口存在但只覆盖孩子**：`internal/tools/task.go:447` `TaskRoster.Cancel(taskID)`，产码调用者**唯一一枚是 `internal/tools/task.go:740`——即 `task.cancel` 那枚模型侧工具**（注册名 `internal/tools/task.go:625` `Name() = "task.cancel"`）；而句柄只在派生子代理时挂上（`internal/tools/subagent_197.go:346` `AttachCancel(bg.ID, cancelChild)`），⛔ **根任务从未 `AttachCancel`** ⇒ 对根 id 调 `Cancel` 走 `internal/tools/task.go:458-459` 分支，答 `"任务 %s 没有在跑（kind=… state=…），没有可停的句柄"`。**今天真能按下的停只有两枚，都不在面板上**：① 原生 Esc 快捷键 → `cmd/wisp/resident_approval_windows.go:173-184` `vetoByEsc()`（它只否决**正在等人的那张卡**，`ra.cards.AwaitingHuman()` 为假时回 `"按下的取消键没有可否决的确认项…"` `:178`），注入点 `cmd/wisp/resident_windows.go:163` `startResidentBall(rt.Registry, ra.vetoByEsc, …)`；② 退出序列第 3 步 `cmd/wisp/resident_approval_windows.go:283` `cancelTaskRoots`，注册于 `cmd/wisp/resident_windows.go:186` `RegisterShutdownHook(proc.StepCancelTasks, …)`——那是**关机，不是用户按停**。 |
| ② 到面板的通路 | **无字段、无读口**。"能不能停"这一维今天不从快照出：`Snapshot` 没有 running/cancellable 的任何一格（`internal/panel/composer.go:57-92`）。页面唯一能推断"在忙"的材料是 `Snapshot.Tasks.Rows[].status`（`internal/panel/subagent_roster_197.go:109-115` `status/statusKnown/statusReason`），而**根行的状态今天没人填**——`cmd/wisp/run.go:1105` `rt.tasks.MarkRoot(bg.ID, task)` 只挂 label，该处注释逐字写着"State is deliberately NOT filed here… StateAnswer keeps saying 「宿主没有登记这一维」 for a root row"。⇒ **快照连"这发正在跑"都说不出口，更说不出"能不能停"**。 |
| ③ 动作方向 | **入向那一跳今天落不到任何停止处理器**：`knownComposerMethod`（`internal/panel/bridge.go:146-148`）与派发表（`internal/panel/composer_dispatch.go:176-208`）都没有停的字。⚠ **它不是"再加一枚 case"就完事**——本票面 AC#7 把 `internal/panel/**` 的写面限在"票 145 已批的局部解冻范围（composer/pump/panel_pump，且只到新增字段）"，**入向名册不在那一寸里**；且停车点已有具名裁定：`docs/reports/HANDOVER.md:464`（§4.0v 那一节内）逐字＝「**若 r3 认为要新增入向方法＝C17 契约面，停手上报，不许自己加**」。⇒ **这一格本腿的结论是"上报"，不是"最小改动面"**，详见 §3.4。 |
| ④ 判据 | **正向钉：两枚，钉的都是退出那一支，不是用户按停**——`cmd/wisp/resident_task_source_246_windows_test.go:183` `TestAC246ResidentTaskRootCancelStopsTheModelCall`、`cmd/wisp/resident_task_source_live_246_windows_test.go:238` `TestLive246ExitCancelsARunningTask`（两枚都在 D38(e) 第 3 步那条链上；我读到定义为止，未跑）。另两枚钉快捷键执行器：`cmd/wisp/resident_approval_246_windows_test.go:40` `TestAC246CancelGestureUsesTheInjectedExecutor`／`:132` `TestAC246CancelStepHookRunsOnTheRealShutdownSequence`。**票面 AC#4 要的那形（"必须在助手正忙时也能到达"）无正向钉；"停止后不许留还在写的尾巴"无钉**（`cmd/wisp/panel_pump_test.go`／`internal/panel/*_test.go` 里 `grep -iE "stop"` 零命中）。 |

### 1.4 草稿（输入框里没发出去的那段）

| 问 | 读数 |
|---|---|
| ① Go 侧生产者 | **盘上没有**。尺：`grep -rn -i "draft\|草稿" cmd internal tools docs --include=*.go` 的**全部**命中都是无关物（`internal/config/migrate.go:30/:102` 与 `internal/config/schema.go:22` 的"D36 draft model"＝配置代际名、`internal/tools/fs_edit_ac34_test.go:75/:86` 的桩文本、若干 `first draft`/`草稿` 指证据件自身的改稿）。`cmd/wisp/**` 与 `internal/panel/**` 的**产码**里没有一个"未发出原文"的存放点。发送侧的现成形状只有 `internal/panel/composer.go:295-298` `OutgoingMessage{Text, Attachments}` 与 `internal/panel/bridge.go:91` `ComposerRequest`，二者都是**一次性入向封套**，没有失败分支把原文留下。 |
| ② 到面板的通路 | **无字段**（`Snapshot`／`ComposerState`／`PumpSources` 三处名册逐一看过，见 1.1）。⚠ **落盘面也没有**：`internal/memory/schema.go` 的建表名册＝`schema_meta(:23)/profile(:30)/memory(:40)/task_log(:53)/tool_call(:69)/approval_grant(:89)/cost_daily(:102)/plugin_state(:111)/provider_health(:136)`，**九枚里没有一枚存未发出文本**。⇒ 这一枚不是"建了没接"，是**两层都缺**。 |
| ③ 动作方向 | **无入向方法**（同 1.2③ 的名册：6 枚，无一枚与草稿有关）。⚠ "发送失败的原文要单独存住"若要落库，最小面不是加 case，而是**加表**＝`internal/memory/schema.go` 的 DDL 与其 `SchemaVersionTarget`（`:129` 现量 `= 2`，注释逐字规定 v1/v2 各自来历并写着"DDL 与 `docs/specs/SPEC-02-data-storage.md` §3 **byte-for-byte** 镜像"）⇒ **动它就动规格**，属票面 AC#7 零字节名单（`docs/specs/**`）。本腿判这一枚应走**内存驻留 + 快照字段**，不走落库；若编排者要落库，那是规格变更、要人工批准。 |
| ④ 判据 | **正向钉：无；负向钉：无**。尺：`grep -rn "^func Test" internal/panel/*_test.go cmd/wisp/*_test.go | grep -iE "draft"` 零命中。**最接近的一枚判据不是用例**：`internal/panel/composer.go:300-333` `ForAgent` 的"不吞用户意图"规则（被拒附件逐条回填进模型读的那段文本，`:327-331`），注释自述是 AC#2 的"do not eat the user's intent"——**同一族原则、不同一枚出口，且它钉的是发给模型的文本，不是输入框**。 |

### 1.5 崩溃自救

| 问 | 读数 |
|---|---|
| ① Go 侧生产者 | **现场记录：非测试生产者，在跑**。恢复边界 `internal/observe/goroutine.go:296` 一带的 `recover()` 把 panic 变成结构化记录 `internal/observe/goroutine.go:153-161` `PanicEvent{Goroutine,Owner,RootID,Class,Recovered,Stack,At}`，默认落 `internal/observe/goroutine.go:167-177` `defaultPanicSink` → `slog.Error("goroutine panic recovered", …)`；产码里这条 slog 会进滚动 JSONL 文件（`cmd/wisp/logsink.go:144` `installLogSink`，产码调用者 `cmd/wisp/models.go:284`；目录由 `cmd/wisp/logsink.go:87` `logSinkDir` 决定）。⚠ **专门给"崩溃现场"的那枚钩子没人接**：`internal/observe/goroutine.go:241` `Registry.SetPanicSink` 的调用者**只有 `internal/observe/goroutine_test.go:117` 一枚**（尺：`grep -rn "SetPanicSink\|PanicSink\|PanicEvent" cmd internal tools --include=*.go`，`cmd/` 下**零命中**）⇒ 覆盖不了默认 sink 的自定义出口今天**只在测试里存在**。**"禁用插件后重启"：那一枚能力盘上没有**——`internal/plugin/` 只有 `doc.go`＋`disposal.go` 两枚文件，`doc.go` 末段逐字写着 `DEFERRED(manifest/goja): implemented by tickets 50 (Tier-1) and 51 (Tier-2 goja). This ticket implements only DisposalScope`；库里那枚 `enabled` 列确实存在（`internal/memory/schema.go:114` `enabled INTEGER NOT NULL DEFAULT 1`），DAO 也在（`internal/memory/dao_misc.go:188` `UpsertPluginState`／`:222` `PluginStateByID`／`:243` `ListPluginStates`），⚠ **三枚 DAO 的产码调用者＝零枚**（尺：`grep -rn "UpsertPluginState\|PluginStateByID\|ListPluginStates" cmd internal tools --include=*.go | grep -v "_test.go:" | grep -v "internal/memory/dao_misc.go"` ⇒ 空）。**自救指令文本：无生产者**（尺：`grep -rnE "self.?rescue|fatal-recovery|disablePlugins|自救" cmd internal tools --include=*.go` ⇒ 现量 **0 行**；我第一版把这条写成 BRE 的 `\|` 拼法，那把尺的 `?` 不是量词，已在本节按 ERE 重跑）。 |
| ② 到面板的通路 | **无字段**。`Snapshot`／`ComposerState` 里没有任何崩溃/上次退出维度（名册见 1.1②）；`PumpSources` 亦无。最近的一枚"给人读的诊断出口"是另一条 CLI：`cmd/wisp/doctor.go`（构建链自检，逐行 `pass/fail/info`，`main.go:47` 那条 usage 行之外另有命令表），它**不吃运行时快照、也不落"上一次为什么没了"**。 |
| ③ 动作方向 | **无入向方法**（6 枚名册里没有）。"禁用插件后重启"若将来要做，写面至少跨三枚：`internal/plugin`（装不存在的清单加载器，票 50/51）＋`internal/memory/dao_misc.go`（那三枚没人调的 DAO）＋配置侧的禁用叶；而 `internal/panel/config_handlers.go:108` `lockedFieldFamilies` 把 `"plugins."` 逐字列进**设置页永不许写**的族名册 ⇒ **今天连"从面板禁用一枚插件"的路都是明令封着的**。 |
| ④ 判据 | **正向钉：只钉"panic 变成结构化记录"这一格，且仅在包内**——`internal/observe/goroutine_test.go:116-118` 造了一枚自定义 sink 收 `PanicEvent`（我读到定义为止，未跑）。⚠ **没有一枚用例把"记下来的那发抵达用户可复制的地方"钉住**，也没有一枚钉"自救指令必须存在"；`cmd/wisp/logsink_test.go` 存在（在列，见 `cmd/wisp` 文件清单）但我**未读到它的用例名册**，所以本表不引用它＝§5 那条判不动。**"能禁用插件重启"今天无任何判据，也不该有**（票面 AC#6 自己写了那一格未裁）。 |

## 2. 逐枚"缺哪一环"

> 四词表：**已存在／建了但没接／只有测试里有／盘上没有**（`AGENTS.md` 与 `docs/reports/HANDOVER.md:461`（§4.0v）都要求"没有"必须分形，两种"没有"下一步动作完全不同）。

| 出口 | 判定 | 缺的那一环（具名到函数） |
|---|---|---|
| 占用 | **建了但没接**（分母＋分子两半都在产码里，都没往面板方向递） | 缺的是**读口→`PumpSources` 槽→`ComposerState` 字段**这三跳里的一枚具名 reader：`cmd/wisp` 侧今天已有 `loop.Budgets()` 的产码调用者（`run.go:1056/:1062/:1109`），但**没有一枚读 `.ContextWindow`**，也没有任何一处在任务运行中读 `History()` 的长度（`Loop.History()` 的产码调用者只有 `internal/agent/loop.go:392/:560/:566`）。 |
| 排队（审批卡的第几条） | **建了但没接**，且**丢掉的那一跳至今无人登记** | `internal/agent/approval/pending_read.go:118` 已经在填 `LiveApproval.Position`；断点在 `cmd/wisp/panel_pump.go:58-74` `liveVerdicts()` 未抄 `Position`，因为载体 `internal/panel/pump.go:77-91` `NativeVerdict` 没有那一枚字段。**这是四枚里最便宜的一枚**：真源、账本、用例骨架全都在。 |
| 排队（用户消息队列／"我发的第几条"） | **盘上没有**（并且是**明写拒绝**，不是漏) | 单槽规则在 `cmd/wisp/resident_task_source_windows.go:411-414` 逐字拒绝并发提交；环内的 `Loop.Steer()`（`internal/agent/loop.go:275`）**调用者零枚**。要做这一枚＝先裁"常驻腿到底要不要排队"，那是**产品决定不是写腿能自决的**（见 §5）。 |
| 停止（按钮可达） | **建了但没接**（动词／解析器／执行体三件都在，缺的是"正忙时还能到达"那一跳） | 缺的入口一跳：`MatchControl` 的唯一产码调用点是 `internal/agent/loop.go:341`（`run()` 开头，必须先能提交第二发）；执行体 `internal/agent/loop.go:524-537` `defaultControl` 与取消口 `internal/agent/loop.go:300` `RunningTask.Cancel()` **都没有产码调用者**；根任务**没有停的句柄**（`internal/tools/subagent_197.go:346` 只给孩子挂，`internal/tools/task.go:458-459` 对根 id 明答"没有可停的句柄"）。⚠ 入向那一跳按既有裁定**不许写腿自己加**（见 §3.4）。 |
| 草稿 | **盘上没有**（内存名册、结构体字段、库里表**三层都缺**） | 最近的一枚现成形状是一次性入向封套 `internal/panel/composer.go:295` `OutgoingMessage`；失败分支今天不留原文。落库一支要动 `internal/memory/schema.go` 的 DDL 与 `SchemaVersionTarget`（`:120`），而该文件注释自述与 `docs/specs/SPEC-02` §3 **逐字节镜像** ⇒ 属票面 AC#7 零字节名单。 |
| 崩溃自救（记一笔＋一条指令） | **只有测试里有**（就"专用自救出口"而言）／**建了但没接**（就"现场已能被记下来"而言）——**两种"没有"必须分开写** | 已有的：`internal/observe/goroutine.go:153` `PanicEvent` 经 `:167` `defaultPanicSink` 进 slog，slog 经 `cmd/wisp/logsink.go:144` `installLogSink`（调用者 `cmd/wisp/models.go:284`）真的落文件。缺的：`:241` `Registry.SetPanicSink` **产码零调用者**（唯一调用者 `internal/observe/goroutine_test.go:117`）；"可复制的自救指令"那一句**无任何生产者**；"禁用插件后重启"的插件地基为空（`internal/plugin/doc.go` 具名 `DEFERRED` 50/51），且那三枚 DAO 零调用者。 |

## 3. 最小落地面

> 前提复述（它决定下面每一步的边界）：**装配根是唯一接缝**这一既有裁定在本仓是硬的——`internal/panel/subagent_roster_197.go:43` 逐字写着「internal/panel must not import internal/tools, so the composition root does the mapping」，同一句话也排除了 `internal/panel` 直接 import `internal/agent`。⇒ **下面每一跳读口都必须落在 `cmd/wisp`；凡是要让 panel 直接够到 agent/tools 的形状＝新开依赖边＝上报，不许自决。**
> 顺序按"最便宜的先落"。

### 3.1 审批卡的"第几条"（最便宜的一枚：真源已在产码，只差最后一跳没抄）

按序四步：

1. `internal/panel/pump.go:77` `NativeVerdict` —— 追加一枚 `Position int`（载体，不是判断；它的兄弟字段全是从队列抄来的）。
2. `cmd/wisp/panel_pump.go:58` `(*agentRuntime).liveVerdicts()` —— 在已有的那段逐字段抄写里补一行 `Position: it.Position`。真源：`internal/agent/approval/pending_read.go:49` 字段声明＋`:118` 填充（`q.position(it)`，定义 `internal/agent/approval/queue.go:268`，1-based FIFO）。
3. `internal/panel/approval.go:39` `ApprovalCardView` —— 追加 `Position int json:"position"`，并在纯函数半边 `internal/panel/approval.go:72` `CardViewFromDecision` 里带上它；`NativeVerdict.CardView()`（`internal/panel/pump.go:99`）是唯一的产码入口，所以只有这一处会填。
4. 判据**只能进新文件**：票 194 的普查件 `docs/evidence/s1/194-method-roster-census-c1.md:140` 已裁定「AC#6 点名 `tokens_fourway_test.go`／`composer_test.go`／`l2_grant_boundary_test.go` 三枚不许动 ⇒ 判据只能进新文件」。票面 AC#3 后半"取消排队不许把后面的挤掉"今天无对象（没有用户消息队列，见 §2），只能钉前半。

⚠ **这一枚撞 `Q-51`，且撞得最直接**：第 3 步加的是 `ApprovalCardView` 的 JSON 键，正落在 `internal/panel/approval_test.go:105` `TestApprovalCardViewJSONKeysMatchFrontendTypes` 与 `internal/panel/composer_test.go:49` `TestComposerContractTypesMatchFrontend` 那两把**双向**尺的射程里（我读到断言本体：`internal/panel/composer_test.go:73-78` 两个方向各一次 `t.Errorf`）。`Q-51`（台账 `docs/reports/pending-and-issues.md:1086`）**仍未答**，口径在 `A383` 那条及其后果句（`docs/reports/pending-and-issues.md:8497` 逐字：加字段很可能让那把尺更红 ⇒ **照实记进证据件、不修、不放宽、不动前端**）。
⇒ **给编排者的形状**：要么按 3.1 走完并在同一枚 commit 里带页面那份声明（＝`Q-51` 甲支，需 owner 点头），要么走一条不动契约的替代——**页面自己数 `Snapshot.Pending` 的下标**；但那是页面推出来的数，不是票面 AC#3 要的"Go 侧可读出口"，选它就得有人明写"AC#3 按页面自数结"。⛔ 这不是本腿能替 owner 选的。
**新开依赖边：零。** `NativeVerdict` 与 `ApprovalCardView` 同包、`cmd/wisp` 已握着队列（装配点 `cmd/wisp/run.go:700` `Verdicts: rt.liveVerdicts` 是既有的）。

### 3.2 占用（两个数齐之前整枚不显示那一枚）

1. **先裁"说哪一件"**（§5 第 2 项）。裁完才谈下面。
2. 分子与分母的读法**都不需要新接口**——本节最值钱的一句：
   - 分母：`internal/agent/loop.go:253` `Loop.Budgets()` 已是导出读口，`cmd/wisp` 已有三枚产码调用者（`run.go:1056`／`:1062`／`:1109`），`.ContextWindow`（`internal/agent/budgets.go:54`）就在那枚返回值里，**只是没人取**。
   - 分子：`internal/agent/loop.go:256` `Loop.History()` 是导出的，`internal/agent/budgets.go:150` `ApproxTokensOfMessage` 也是导出的 ⇒ **在 `cmd/wisp` 侧用这两枚现成的就能算出"当前历史成本"**，不必给 `Loop` 加方法。
   ⚠ 反过来才是雷：停车点登记过「给 `internal/agent.Loop` 加一个"收 `taskID` 的导出方法"会把 `gate-clauses.sh` 的 G3 腿从安静打成 BAD ⇒ 裁定＝入口放 tools 侧，非放 Loop 不可就停手上报」（`docs/reports/HANDOVER.md:495`，18:5x 那一段）。占用这一枚用不到 taskID，**但任何"给 Loop 加导出方法"的实现形状都必须先按那条停手上报**。
3. `internal/panel/composer.go:235` `ComposerState` —— 新加**必须是指针＋`omitempty` 的一节**，照 `internal/panel/composer.go:74` `Instructions`／`:91` `Tasks` 那两枚的先例；否则 `internal/panel/pump_test.go:123`／`:291`、`internal/panel/subagent_roster_197_test.go:214`、`internal/panel/subagent_stream_197_test.go:130` 那四枚"无 reader 恰好发四枚键"的字节钉会红（这四枚我逐枚读到断言字面才写，非转述他人读数）。
4. **"未知"必须是具名状态，不许退化成 0／空串**（票面 AC#2）——本仓已有现成形状可抄：`internal/panel/composer.go:255-256` `CurrentModel`+`ModelKnown`、`:268-269` `Credential`+`CredentialKnown`（"没读过"与"读出来是空"是两件事）。分子分母各配一枚 `…Known`，或一节里放 `Occupied/Total/BothKnown`，都落在"新增字段"那一寸（票面 AC#7 逐字：`internal/panel/**` 写面"只到新增字段为止"）。
5. `cmd/wisp/run.go:699` `rt.pump = panel.NewSnapshotPump(panel.PumpSources{…})` —— 在这一枚已有字面量里补 reader 槽（`internal/panel/pump.go:121` `PumpSources` 现 12 枚槽、`:699` 处已在填 8 枚）。**这是"装配根是唯一接缝"的既有位置，不是新边。**
⚠ 同 3.1：新增 `ComposerState` 键也在那把双向尺射程里 ⇒ 同一句 `Q-51`。

### 3.3 停止（今天能在 Go 侧落的只有"看得见"，不是"按得下"）

- 能在 `internal/panel` 解冻那一寸里落的：**把"此刻有没有一发在跑／它能不能被停"送进快照**。现成源两枚都在 `cmd/wisp`：`cmd/wisp/resident_task_source_windows.go:154` `src.running`（正是 `:411-414` 用来拒绝第二发的那枚）与 `:159` `src.seq`（`:416-417` 自增，今天只进 goroutine 名与审计行 `:422/:425/:437`）。再带上"根任务没有停的句柄"这一枚事实（`internal/tools/task.go:458-459` 会答"没有可停的句柄"），快照就能诚实说"现在不能停"，而不是画一枚按下去没反应的按钮——**这就是票面硬约束"宁缺毋造"落在停止这一枚上的形状**。
- 不能自决的两枚，具名上报：① 给**根任务**挂停的句柄＝动 `internal/tools`（`TaskRoster.AttachCancel`（`internal/tools/task.go:420`）今天唯一产码调用者是孩子侧 `internal/tools/subagent_197.go:346`），**不在票 167 AC#7 的 `internal/panel` 解冻范围内**；② 让面板"按停"＝新增 C17 入向方法，见 3.4。

### 3.4 必须上报、写腿不许自决的三枚

1. **任何新入向方法**（停止按钮、草稿存取都算）。依据＝`docs/reports/HANDOVER.md:464` 逐字「若 r3 认为要新增入向方法＝C17 契约面，停手上报，不许自己加」，加上票面 AC#7 把 `internal/panel/**` 写面限定在"票 145 已批的局部解冻范围（composer/pump/panel_pump，且只到新增字段）"——**入向名册不在那一寸里**。
   ⚠ **同时给编排者一枚可能改变裁定的读数**：那枚会数名册的冻结尺 `internal/panel/l2_grant_boundary_test.go:2030` `guardRosterOf` **不是硬编码名单**——它从 `bridge.go` 的 `Method*` 常量与 `knownComposerMethod` 的 case 标签**双向推导**，函数头注释逐字写着「which survives the roster growing by any number of names and still refuses the shape it was written for」（`internal/panel/l2_grant_boundary_test.go:2028`，函数体 `:2030`）。⇒ 加一枚方法**不需要改那枚冻结件**，只要：新常量声明在 `internal/panel/bridge.go` 的 const 块（名字以 `Method` 开头、值为字面串），且 `internal/panel/bridge.go:148` 那枚 `case` **仍写成一行逗号分隔**（`internal/panel/l2_grant_boundary_test.go:2128-2130` 有一条"case 列表不是一行就判脏"的形状钉）。但"技术上不用改冻结件"≠"治理上过得了"——**这一枚仍按停手上报走，我只是把成本量清给编排者**。另注意 `docs/evidence/s1/248-settings-write-path-r1.md:155` 已登记的代价：`guardRosterOf` 只读 `bridge.go` 一枚文件的 const 块，把新常量挪去同包另一枚文件会**假红**。
2. **草稿若要落库**＝动 `internal/memory/schema.go` 的 DDL 与 `SchemaVersionTarget`（`:129` 现量 `= 2`），而该文件注释自述与 `docs/specs/SPEC-02-data-storage.md` §3 逐字节镜像 ⇒ 属票面 AC#7 的 `docs/specs/**` 零字节名单。⇒ **本腿建议：草稿走内存＋快照，不走落库**；要走落库是规格变更，需人工批准。
3. **崩溃自救里"禁用插件后重启"那一支**——票面 AC#6 已明写未裁，本腿量到三重空：`internal/plugin/` 只有 `doc.go`＋`disposal.go`（清单加载器具名 `DEFERRED` 50/51）、那三枚 `plugin_state` DAO 产码调用者 0（`internal/memory/dao_misc.go:188/:222/:243`）、而 `internal/panel/config_handlers.go:108` `lockedFieldFamilies` 把 `"plugins."` 逐字列进设置页**永不许写**的族。⇒ **明写不做**（票面硬约束 4 同向）。

### 3.5 崩溃自救今天能落的那半（记一笔＋给一条可复制的指令）

- "记一笔"**已经成立**：`internal/observe/goroutine.go:236` `NewRegistry` 就装了 `defaultPanicSink`（`internal/observe/goroutine.go:167-176` 带 `"stack", ev.Stack`）→ `slog` 默认被 `cmd/wisp/logsink.go:144` `installLogSink`（调用者 `cmd/wisp/models.go:284`）做成"脱敏 JSONL 文件＋stderr"两路。⇒ 落地腿**不要新造记录通路**，只做"让人拿到那个文件"：路径由 `cmd/wisp/logsink.go:87` `logSinkDir(dataDir)` 决定，把它连同一条指令文案送进快照或一条 CLI 出口即可。
- ⚠ **会咬人的一格**（本腿现量）：生产进程里 `observe.Registry` **有三枚实例**——`internal/observe/goroutine.go:230` `var Default = NewRegistry()`、`cmd/wisp/run.go:1009` 给根环路的 `Registry: observe.NewRegistry()`（私有）、`internal/tools/subagent_197.go:294` 给每个孩子又是一枚私有。三枚都装了 `defaultPanicSink` ⇒ 栈都能落盘，但**只往 `Default` 上装自定义出口会漏掉所有 agent 任务与所有子代理的 panic**（今天动过这枚钩子的只有 `internal/observe/goroutine_test.go:117`）。⇒ 要装就得逐枚装，或在装配根收一处。
- **零新开依赖边**：`Registry`、`installLogSink`、`logSinkDir` 都已在 `cmd/wisp` 被调用。

### 3.6 若 owner 明天就要用上那四个数：一句实话

按上面的量法，**明天能用上的只有 3.1（第几条）与 3.2（占用）两枚，而两枚都卡在 `Q-51` 那一寸**（新 JSON 键要有人写页面那份声明）。停止那一枚今天能给的是"显示不能停"；草稿那一枚三层都缺（Go 侧字段、库表、入向方法）。⇒ 如果"四个数明天都要"是硬期限，要先答的是 `Q-51`，不是派落地腿。

## 4. 我可能判错的条目

> 每条附「如果错了后果」。本节是本腿的自我对抗面，不是免责声明：判错的代价我写成可核的动作。

1. **"零调用者"这类断言全部靠 `grep` 的形状匹配，接口值／函数值／方法表能溜过去。**
   射程内我扫的是名字本身：`Steer`／`.History()`／`TotalTokens(`／`RunningTask.Cancel`／`Control:`／`UpsertPluginState|PluginStateByID|ListPluginStates`／`SetPanicSink`。若某处把 `Loop.Steer` 取成一枚 `func(string) bool` 字段、或经接口方法表分派，我的尺看不见。
   **如果错了后果**＝§2 那行"用户消息队列＝盘上没有"里"`Steer` 调用者零枚"这一支撑是假否证，真正缺的只是把它接到面板；落地腿会去新建一条已经存在的链。**缓解**：这一枚的**另一半证据不依赖数调用者**——`cmd/wisp/resident_task_source_windows.go:411-414` 逐字在运行中拒绝第二发，那是产码里的显式分支，不是零命中。⇒ 即使 `Steer` 被我读错，**"没有面向用户的排队"仍成立**。

2. **我把"分母不出包"的措辞写重了。** `cmd/wisp/run.go:1109` `agent.NewSpiller(filepath.Join(rt.spec.dataDir, "artifacts"), loop.Budgets())` 已经**把整枚 `Budgets` 结构体递到 `internal/agent` 之外**，只是没有一格被命名为"分母"。
   **如果错了后果**＝实现程读了我那句"没有一枚把 `.ContextWindow` 递出去"，以为要先加一层导出面，白建一处 API。**正确读法**：`Loop.Budgets()`（`internal/agent/loop.go:253`）已经是导出读口且已有 3 枚产码调用者（`run.go:1056/:1062/:1109`），**新加一字段就够，不需要新接口**。

3. **占用那一枚，我列了三枚候选分子而没有替 owner 选。** `internal/agent/guard.go:132` `TokensUsed()`＝**跨轮累积**的 in+out；`internal/agent/compress.go:110` `TotalTokens(hist)`＝**当前历史**成本；`internal/agent/prompt.go:181` `Assembler.TotalTokens()`＝**段预算上限之和**（D39 那 2300 那一族，压根不是"已用"）。三枚语义互不等价，我在 §1.1 里都给了出处，但**"占用条该说哪一件"不是读数能答的**。
   **如果错了后果**＝按"还剩多少空间"画出来的条用累积消耗做分子，会在第 8 轮显示 95% 而当下历史只有 30%（反之亦然），用户照着它决定要不要"继续"，而那条数**是假的**——这正撞票面硬约束 1"不许填个假数"。**处置**＝这一条落进 §5 第 2 项，摆给 owner 一个词。

4. **分母可能永远不是用户填的那个数。** `internal/panel/config_handlers.go:60` `FieldModelContextWindow` 是目录里 per-provider+model 的可写叶，而 `cmd/wisp/run.go:1012` 读的是**解析后的** `ep.ContextWindow`；`internal/agent/budgets.go:86-87` 逐字写着 `if ctxWindow <= 0 || ctxWindow > ReferenceContextWindow { ctxWindow = ReferenceContextWindow }`，`ReferenceContextWindow = 128000`（`internal/agent/budgets.go:18`）。⇒ 目录没声明、或声明值**大于** 128000 时，缩放用的都是 128000，不是用户填的那个。
   **如果错了后果**＝我若把这写成"分母现成"，占用条会给用户显示一个他从来没配过的分母；我若把它写成"分母是配置叶"，实现程会把 128000 当成假数。**处置**＝§5 第 2 项一并摆。

5. **四"已存在"的字段名册我只对 `Snapshot`/`ComposerState`/`PumpSources` 三处核过，没有核 `internal/observe` 的 SLO 出口。** `internal/watchdog/doc.go` 逐字写着它 non-responsibility 是"no sampling backend itself (observe provides SLO plumbing)"，而 `cmd/wisp/slo_windows.go`/`internal/observe/thresholds.go` 那条 D32 资源线是**另一族占用**（内存/CPU）。票面 line 4 断言"这枚占用条的数据源就是 D32"，我在 §1.1 里按**上下文窗口**回答了这一枚。
   **如果错了后果**＝owner 要的其实是"D32 那两枚资源数上屏"，那我这整节都答在另一枚对象上——而那一枚**已有读法**（票面 line 37 自己写"不做成本页（D32 那两个数已有读法）"⇒ 我据此判票面 line 4 那句与 line 37 自相矛盾，见 §6 第 2 项）。**缓解**：两条我都量了，编排者按对象挑。

6. **新增快照字段会不会多打红一枚尺，我判不了，只能给出射程。** 我知道的：`internal/panel/pump_test.go:123`／`:291`、`internal/panel/subagent_roster_197_test.go:214`、`internal/panel/subagent_stream_197_test.go:130` 逐字钉"无 reader 的泵恰好发四枚键"（`composer,generatedAt,pending,results`），所以新字段**必须**是指针＋`omitempty` 才不动这四枚（先例＝`internal/panel/composer.go:74` `Instructions`／`:91` `Tasks`）；我也读到 `internal/panel/composer_test.go:49` 与 `internal/panel/approval_test.go:105` 那两把**双向**尺会拿 Go 的 JSON 键去对界面声明，逐键比。
   **如果错了后果**＝落地腿加了字段、那两把尺红了，它以为是"契约变更被拦"，其实那是 `docs/reports/HANDOVER.md:467` 那张红名册里**已在册的历史红**（该节同时自述"红名册本轮不复量"、引用要带锚点并当过期）。⇒ 我不引用那四枚具体用例名作结论，改在 §5 第 1 项登记"今天还在不在，判不了"。

7. **`plugin_state` 的"零调用者"不等于"禁用插件的地基缺"。** 表列在（`internal/memory/schema.go:114`）、DAO 三枚在（`internal/memory/dao_misc.go:188/:222/:243`），我把"缺"落在 `internal/plugin` 的清单加载器（`internal/plugin/doc.go` 末段具名 `DEFERRED` 50/51）。
   **如果错了后果**＝我把 §2 那一行写重了：真实形状是"存储与 DAO 已存在、只有没人调用；缺的是插件本体"——那和"盘上没有"的下一步动作不同（前者只要接，后者要先建）。⇒ 我在 §2 那行**没有**用"盘上没有"评整枚崩溃自救，只评"自救指令文本"那一支。

8. **崩溃自救这一枚我踩在编排者自己刚收回过的那同一格上。** `docs/reports/HANDOVER.md` §4.0y 第 1 笔逐字登记：编排者曾写"生产没装 panic sink ⇒ 真崩时盘上零栈"并据此立票 249，后被现量推翻（`internal/observe/goroutine.go:236` `NewRegistry` 就装了 `defaultPanicSink`、`:167-176` 带 `"stack", ev.Stack`、`cmd/wisp/logsink.go:145-162` 把默认 slog 做成"脱敏 JSONL 文件＋stderr"两路），并立了新规矩：**"'盘上没有 X'这类负向句，落笔前先读'谁产出／谁投递／谁落盘'三处；能力问题要跑一遍回答，不许靠数调用者回答"。** 我按这条规矩把 §1.5① 写成"记录已存在（产出→投递→落盘三处逐枚读到）、只有'自定义出口'与'自救指令文本'缺"，**没有**写"真崩时零栈"。
   **如果错了后果**＝我仍可能低估了一格：`cmd/wisp/run.go:1009` 给根环路的是**私有 registry**（`observe.NewRegistry()`），`internal/tools/subagent_197.go:294` 给每个孩子又是一枚私有 registry，而 `observe.Default`（`internal/observe/goroutine.go:230`）是第三枚。三枚都在 `NewRegistry` 里装了 `defaultPanicSink` ⇒ **栈都能落盘**，但"给崩溃装一枚自定义自救出口"这件事必须**逐枚装**，只装 `Default` 会漏掉所有 agent 任务的 panic。这一条我写进 §3.5，因为它是会咬落地腿的形状。

9. **我没有排除"草稿其实归界面自己管"这一支。** 禁令所限我未读 `frontend/**`，所以"Go 侧盘上没有"是我的全部读数；如果未发出文本今天由页面自己存着，那 §2 那行在 Go 侧仍成立、但**不是缺口**。
   **如果错了后果**＝编排者按我这张表派一枚 Go 侧草稿落地腿，做出来的东西和页面自己那套**各存一份**，用户看到的"发不出去的那条"会有两个来源——这正是票面 AC#5 要防的"用户以为丢了半段话"的镜像。**处置**＝归 §5 第 3 项，只有 owner/他那枚界面 agent 能一句话定。

## 5. 判不动的地方

> 甲＝谁补得上、怎么补；乙＝补不上，明写不做。

1. **`internal/panel` 那四枚契约/令牌尺今天红不红。** 判不了的两条腿都硬：我不能跑 `go test`（`198-r1`／`33-r8b` 此刻在飞；本腿交件时两者都已落树，但我仍**没有**跑），也不能读 `frontend/src/lib/panel.ts`（两层禁令）。唯一在册读数是 `docs/reports/HANDOVER.md:467`（§4.0v）那张逐名表，它自己写着"引用要带口径＋锚点 sha"、且同节另一处写"红名册本轮不复量"⇒ **当过期**。
   **甲**＝编排者在编队空时整包复跑（口径已在 `docs/reports/HANDOVER.md:161` 那条：带 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/...`，**不是 `-run` 单跑**）＋逐名比红名集合，出一版带锚点的新表；这一步买到的东西＝§3.1 那一步到底会不会"多红一枚"当场可判。
   **乙**＝不跑就不判。⇒ 本表**不**把任何"红/绿"当结论写；§3 里凡涉及那两把双向尺的地方，我都只写"射程"与"撞 `Q-51`"，不写"会红几枚"。

2. **占用那枚数到底该说什么。** 三枚候选分子（`guard.go:132` 累积 / `compress.go:110` 当前历史 / `prompt.go:181` 段预算上限）与两枚候选分母（`budgets.go:54` 缩放后窗口 / `FieldModelContextWindow` 目录叶）**语义互不等价**，选哪一对是产品裁定。
   **甲**＝编排者把一句话摆给 owner：「占用条要说的是'这一次对话还剩多少'还是'这次任务一共花了多少'？」——两支各对应一组现成读口，答完 §3.2 立刻能定形，**零新依赖边**。
   **乙**＝不答则**不做占用那一枚**。任何腿替 owner 选一对，都是在硬约束 1 上赌一个假数。

3. **常驻腿要不要真的排队。** `cmd/wisp/resident_task_source_windows.go:396-400` 把单槽写成**一条决定**并给了理由（"the loop's per-task budgets, the panel roster and the L1 window were all designed around one operator waiting on one orb, and a queue of tasks in a leg that cannot show a page would be a queue nobody can read"），不是漏。推翻它要连 `internal/agent/budgets.go:47` `MaxToolConcurrency = 4` 与名册槽位一起重算。
   **甲**＝编排者先问"那句 'a queue nobody can read' 的前提（面板不能显示一页）今天还成立吗"——票 33 交完后它**可能已经过期**；过期成立则这一枚从"产品裁定"降级为"落地活"。
   **乙**＝我不判这一枚。§2 那一行只登记"盘上没有"，不写"该建"。

4. **`cmd/wisp/logsink_test.go` 里有没有已经钉住"panic 栈能落盘"的用例。** 文件我 `ls` 到了（在 `cmd/wisp` 清单里），但**未逐枚读它的用例名册**；按派单规矩"要引用一枚用例名必须自己读到它的定义为止"，我没读就不引用。
   **甲**＝任何腿 `grep -n "^func Test" cmd/wisp/logsink_test.go` 再逐枚读定义（只读，不吃 CPU）。
   **乙**＝§1.5④ 那一格维持"无（对自救指令而言）"，并且**不**把它当成"崩溃现场无人钉"的证据。

5. **票面 line 4 那五枚外部路径。** `packages/client/ui-conversation/src/client/context-occupancy.ts:5,21`／`.../queue/QueueDock.tsx`／`stop-shortcut.ts`／`ui-conversation/README.md`／`apps/desktop/src/fatal-recovery.ts`——尺（**逐根分开取数**，因为这把尺的字面文本会命中票面与本文件自己）：`grep -rn "context-occupancy|QueueDock|stop-shortcut|fatal-recovery|disablePlugins"` 按根跑 `cmd:0 internal:0 tools:0 docs:0 scripts:0`，只有 `.scratch` 有命中＝**票面 line 4 那一句＋本件自己**。⇒ 这些路径在本仓无对应物。票面自己说"我读过原文／我读过接口清单"，那是**外部参考项目**，不是本仓文件。
   **甲**＝只有 owner 或编排者能裁"要不要把那份对标转成仓内需求"；若要，请把这些路径的来源仓库/版本写进票面，否则下一程无从复核。
   **乙**＝本腿**不**据外部形状推任何 Go 侧结论；§1/§2 的每格只引用本仓读到的 `file:line`。

## 6. 我推翻票面／派单哪一句

> 派单要求本节不空，且给我的所有断言都标了待验。逐条对，每条给尺与读数。

1. **票面 line 25（AC#1）自述"我写票时没有逐枚现量过这四枚的出口状态"** ⇒ 这一条不是推翻，是**兑现**：本件 §1/§2 就是那一格要的"命令／API／快照字段名逐枚列 有／无／有但没接"。四枚里没有一枚判"已存在"，也没有一枚该"硬造第二套"——**唯一一枚真源已在产码、只差最后一跳没抄的是审批序号**（`internal/agent/approval/pending_read.go:118` 已填 `Position`，`cmd/wisp/panel_pump.go:58` `liveVerdicts()` 未抄，因为载体 `internal/panel/pump.go:77` `NativeVerdict` 没那枚字段）。按 AC#1 末句"量到'已有'就报回、不许硬造第二套"→ **报回**。
2. **推翻票面 line 4 那半句**：「`D32`（资源 SLO——**这枚占用条的数据源就是它**）」。读数＝票面自己在 line 37 写着「不做成本页（`D32` 那两个数已有读法，界面归界面）」，而 D32 那一族是**内存/CPU 资源带**（尺：`internal/watchdog/doc.go:8-11` 逐字——per-state limits、`proc.TreePrivateBytes`、GDI/User objects、handles、goroutine/thread count；`internal/observe/thresholds.go` 是本腿未读的阈值件），跟"这次对话还剩多少空间"不是同一枚对象；上下文那一路的真源在 `internal/agent/budgets.go:16` `ReferenceContextWindow`／`:54` `ContextWindow`，归 `D15` 那一族。**两枚"占用"的数据源不同，而票面 line 4 把它们写成同一枚** ⇒ 落地腿若照 line 4 去接 D32，接到的是一条内存条。⚠ 这条**不推翻 D 表本身**（D32／D15 各是什么以 `docs/PLAN.md` 为准），推翻的是"占用条的数据源就是 D32"那一句归属。
3. **推翻派单与既有普查件共用的那句"入向 `Handle` 生产调用者＝零枚"**：`docs/evidence/s1/194-method-roster-census-c1.md` §5 逐字写着「`ComposerDispatch` 的入向入口 `Handle` `internal/panel/composer_dispatch.go:120` 的生产调用者＝**零枚**…⇒ 本表 18 枚里没有任何一枚今天有真听众」。现量已推翻：**两枚产码调用者**——`cmd/wisp/panel_host_windows.go:551` `return m.disp.Handle(ctx, raw)`（真 WebView2 回调那侧）＋`cmd/wisp/panel_inbound.go:163` `reply, err := disp.Handle(ctx, raw)`（CLI 腿）；常驻那条链也接上了（`cmd/wisp/panel_resident_windows.go:170` → `newComposerDispatchChain`，由 `cmd/wisp/resident_windows.go:142` 调）。
   ⚠ **但派单要我防的那一误同时成立，别翻到另一头**："入向有了一条线"≠"所有口都接好了"——`cmd/wisp/panel_inbound.go:271-281` 那枚字面量里 `Workspace`/`Attachment`/`Message` **逐枚写死 nil**（注释自述归属：workspace＝票 186、attachment＝票 92、message＝票 35），到达即走 `internal/panel/composer_dispatch.go:230` `unattached()` 按名拒绝。
4. **推翻派单里"票 114 那条线接的是 `config.get`/`config.set`"这句的覆盖面**：同一条链里 `Mode: modeWrites`（`cmd/wisp/panel_inbound.go:272`）**也已接**，处理器是 `internal/panel/composer_handlers.go:86` `ModeWriteHandler`；只是 `Confirm` 为 nil，`cmd/wisp/panel_inbound.go:239-242` 注释逐字「no L2 card on this leg, so auto_approve is unreachable and every widening request is refused before Store.Set is ever reached」。⇒ 只提 config 那两枚会**低估**入向现状；本表 §1.2③／§3.4 按"两枚有名处理器（Mode、Config）＋四枚按名拒绝"写。
5. **推翻票面 line 4 末段给人的错觉**：「崩溃时给一个"禁用插件后重启"的出口（`apps/desktop/src/fatal-recovery.ts`，里面真有一枚 `disablePlugins()`，我读过接口清单）」——那是**外部项目**的接口清单。本仓对应物是三重空（§3.4 第 3 条：`internal/plugin/` 只有两枚文件且具名 `DEFERRED` 50/51／那三枚 `plugin_state` DAO 产码调用者 0／`internal/panel/config_handlers.go:108` 把 `"plugins."` 封在设置页永不许写的族里）。票面 AC#6 与硬约束 4 已经把这一支拦住了，所以这条推翻的是**来源描述的口径**（"里面有真东西"读起来像本仓也有），不是推翻判据。
6. **一处推翻不了也证实不了、但必须点名给编排者的**：§3.1／§3.2 依赖的那两把双向尺的**当前颜色**我没判（§5 第 1 项）。如果它们在今天的 HEAD 上其实已绿（＝页面那份声明已被界面那支补上 `instructions`/`tasks`），那"加字段必多红一枚"的顾虑就只剩"声明要同批跟上"，`Q-51` 的实际成本比我写的低。⇒ 这条**不改变 §3 的任何文件清单**，只改变"要不要现在就去催 `Q-51`"。
7. **本腿自己写坏过又改掉的两把尺（记下来，免得下一程复踩）**：① 我第一版把"自救指令文本"那把尺写成 BRE 的转义或形，`?` 在那里不是量词——已按 ERE `grep -rnE "self.?rescue|fatal-recovery|disablePlugins|自救" cmd internal tools --include=*.go` 重跑，现量 **0 行**。② 我第一版把票面那五枚外部路径写成"本仓零命中"，而那把尺的字面文本**会命中票面与本文件自己**（`.scratch` 两根）；已改成**逐根取数**：`cmd:0 internal:0 tools:0 docs:0 scripts:0`。⇒ 这与 `docs/reports/HANDOVER.md` §4.0y 第 1 笔那条新规矩同源（负向句落笔前先读"谁产出／谁投递／谁落盘"），也与 198-r1 那条提交说明同一形（尺的拼法写进结论文件 ⇒ 自查尺当场读到命中是尺文本自己）。
