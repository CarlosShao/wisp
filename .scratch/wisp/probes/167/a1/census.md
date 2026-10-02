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
| ① Go 侧生产者 | **现场记录：非测试生产者，在跑**。恢复边界 `internal/observe/goroutine.go:296` 一带的 `recover()` 把 panic 变成结构化记录 `internal/observe/goroutine.go:153-161` `PanicEvent{Goroutine,Owner,RootID,Class,Recovered,Stack,At}`，默认落 `internal/observe/goroutine.go:167-177` `defaultPanicSink` → `slog.Error("goroutine panic recovered", …)`；产码里这条 slog 会进滚动 JSONL 文件（`cmd/wisp/logsink.go:144` `installLogSink`，产码调用者 `cmd/wisp/models.go:284`；目录由 `cmd/wisp/logsink.go:87` `logSinkDir` 决定）。⚠ **专门给"崩溃现场"的那枚钩子没人接**：`internal/observe/goroutine.go:241` `Registry.SetPanicSink` 的调用者**只有 `internal/observe/goroutine_test.go:117` 一枚**（尺：`grep -rn "SetPanicSink\|PanicSink\|PanicEvent" cmd internal tools --include=*.go`，`cmd/` 下**零命中**）⇒ 覆盖不了默认 sink 的自定义出口今天**只在测试里存在**。**"禁用插件后重启"：那一枚能力盘上没有**——`internal/plugin/` 只有 `doc.go`＋`disposal.go` 两枚文件，`doc.go` 末段逐字写着 `DEFERRED(manifest/goja): implemented by tickets 50 (Tier-1) and 51 (Tier-2 goja). This ticket implements only DisposalScope`；库里那枚 `enabled` 列确实存在（`internal/memory/schema.go:114` `enabled INTEGER NOT NULL DEFAULT 1`），DAO 也在（`internal/memory/dao_misc.go:188` `UpsertPluginState`／`:222` `PluginStateByID`／`:243` `ListPluginStates`），⚠ **三枚 DAO 的产码调用者＝零枚**（尺：`grep -rn "UpsertPluginState\|PluginStateByID\|ListPluginStates" cmd internal tools --include=*.go | grep -v "_test.go:" | grep -v "internal/memory/dao_misc.go"` ⇒ 空）。**自救指令文本：无生产者**（`grep -rn -i "self.?rescue\|fatal-recovery\|disablePlugins" cmd internal tools` 零命中）。 |
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

（后续 commit）

## 4. 我可能判错的条目

（后续 commit）

## 5. 判不动的地方

（后续 commit）

## 6. 我推翻票面/派单哪一句

（后续 commit）
