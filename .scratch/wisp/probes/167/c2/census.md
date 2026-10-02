# 167-c2 — 面板占用条／排队序号／停止按钮三枚 Go 侧接线：只读普查

> 本文件由只读普查腿 `167-c2` 产出。**零产品码改动**；唯一可写文件即本文件（`.scratch/wisp/probes/167/c2/census.md`）。
> 票面：`.scratch/wisp/issues/167-four-small-outputs-the-panel-needs-occupancy-queue-stop-draft-plus-crash-self-rescue.md`（AC 勾选框一枚未碰）。
> 前程：`167-r1`/`167-r1b` 两程断流零产码；骨架 `69c5825c`（`docs/evidence/s1/167-panel-widgets-r1.md`，58 行全占位）＝证据件骨架，只当线索、不抄读数。
> 参考件：`.scratch/wisp/probes/167/a1/census.md`（166-c1 那轮 166-a1 普查）——只当线索，本件每格读数均自己现量复核。
> 搜索根一律显式：`cmd internal tools docs scripts .scratch`。`frontend/**` 与 `design/**` 未读、未引用。⛔ 未跑任何 Go 命令。
> 派单要求第 30 轮前填满；本件先落骨架（§0＋各区标题），§1–§3 由后续 commit 逐节填满。

## 0. 起手锚（逐字读数）

`date -Iseconds`：

```
2026-10-02T22:13:31+08:00
```

分支（`git rev-parse --abbrev-ref HEAD`）：

```
dev
```

`git log -1 --format='%H %ad' --date=iso`：

```
0ce6cb91af02d583af58c6bdf67283aef1c7075b 2026-10-02 21:53:10 +0800
```

`git status --porcelain cmd internal tools docs scripts`（起手读数，逐字）：

```
（空）
```

起手读数解读：产码五个根全部干净——本腿读到的 `cmd`/`internal`/`tools` 内容即 HEAD `0ce6cb91` 的内容。`.scratch` 下另有他程遗留的 ` M`/`??`（`probes/152`、`probes/161/r6/logs/flip-*.txt` 等，与本腿写面不相交）。

## 1. 问一：占用条（"任务进行中"那一条的状态数据今天由谁产）

> 票面问的"占用条"本腿按**上下文窗口占用**量（used ÷ total）；D32 的资源带（内存/CPU）是另一族，票面 line 37 自己写"不做成本页（D32 那两个数已有读法）"，两族的分叉 167-a1 §6.2 已推翻过票面 line 4 的归属，本件不重复，只在 §6 给本程新增的一处读数。

### 1.1 快照组装链（谁产 → 谁拼 → 谁发）

链路逐跳（本程现量，2026-10-02 锚 0ce6cb91）：

1. **生产 reader**（都在 `cmd/wisp/panel_pump.go`，全部挂在 `agentRuntime` 上）：
   - `panel_pump.go:58` `liveVerdicts()`——审批卡（`rt.gate.Queue().LiveApprovals()`，:62）；
   - `panel_pump.go:83` `workspaceView()`；`:108` `instructionBundle()`；`:146` `taskRosterState()`；`:228` `gitView()`；`:242` `currentModel()`。
   - **没有一枚叫"占用"的 reader 存在**（尺：`grep -n "func (rt \*agentRuntime)" cmd/wisp/panel_pump.go` → 共 11 枚，逐枚过目，无 occupancy/used/total 形）。
2. **装配**：`cmd/wisp/run.go:699-726` `panel.NewSnapshotPump(panel.PumpSources{…})`——现填 9 枚 reader 槽（`Verdicts/Mode/Workspace/Git/Model/Credential/Results/Instructions/Tasks`）＋`Out`。`PumpSources` 全部槽名册＝`internal/panel/pump.go:121-188`（12 槽＋Out＋Now＋AttachmentMax），**没有占用读口**。
3. **拼包**：`internal/panel/pump.go:207-293` `Snapshot()`——把各 reader 的值经 `NewSnapshot`/`NewComposerState` 填进 `Snapshot`。
4. **出口**：`internal/panel/pump.go:311` `Publish()` → `cmd/wisp/panel_pump.go:319` `bookPanelSnapshot()`——**今天唯一生产出口是账本一行摘要**（`panelSnapshotSummary`，`panel_pump.go:331-362`，带 depth/pending/mode/ws/results/sha256），整包字节留在内存 `rt.lastSnap`（`cmd/wisp/run.go:369-370`）。transport（WebView2 推送）不存在，`pump.go:16-23` 注释自述。
5. **驱动**：`cmd/wisp/run.go:735` `rt.ui.publish = rt.publishPanelSnapshot`（审批态变化时发）＋`cmd/wisp/run.go:998` `consoleSink{… publish: rt.publishPanelSnapshot}`（工具起止/流结束/出错时发，`cmd/wisp/run.go:1260-1300`）。

### 1.2 载体字段名册（占用条有没有地方可放）

- 顶层 `Snapshot`（`internal/panel/composer.go:57-92`）：6 枚字段 `pending/results/composer/generatedAt/instructions/tasks`，无占用。
- `ComposerState`（`internal/panel/composer.go:235-270`）：14 枚 JSON 字段，逐枚：`mode/workspace/attachments/acceptedAttachmentMimes/maxAttachmentBytes/attachmentError/git/currentModel/modelKnown/credentialState/credentialKnown`（`sed -n 235,270p` + `grep json:` 现量）——**没有 used/total 任何一格**。
- `PumpSources`（`internal/panel/pump.go:121-188`）：无占用 reader 槽。

### 1.3 三态（空闲/占用/排队）各由哪个变量区分

**关键读数：快照里不存在"空闲/占用/排队"三态的任何变量。** 一枚一格都没有。要"说得出"只能从这些间接材料推：

| 三态 | 今天的间接材料 | file:line |
|---|---|---|
| 排队（有审批卡等人） | `Snapshot.Pending` 数组长度（页面可数），但**无一具"深度"标量** | `internal/panel/composer.go:58`；标量真源 `internal/agent/approval/queue.go:133` `Queue.Depth()` 有导出、有产码调用者（`cmd/wisp/approval_reply.go:377`），但 pump 装配 `cmd/wisp/run.go:700` 只接了 `Verdicts`，**没接 Depth** |
| 占用（任务在跑） | `Snapshot.Tasks.Rows[]`（子代理行有 status，`internal/panel/subagent_roster_197.go:109-115`）；**根行状态今天没人填**——`cmd/wisp/run.go:1105` `rt.tasks.MarkRoot(bg.ID, task)` 只挂 label+kind，`internal/tools/task.go:355-357` `MarkRoot` 逐字注释"State is deliberately NOT written here"（票 196 归属）⇒ `StateAnswer`（`internal/tools/task.go:158-168`）对根行答"宿主没有登记这一维" | 根行状态为空的机制＝`internal/tools/task.go:349-357`＋`cmd/wisp/run.go:1100-1105` 注释逐字 |
| 空闲 | 无字段。页面只能从"pending 空 + tasks 无 running 行"推断，而根行 status 恒 unknown ⇒ **"空闲"与"根任务在跑"今天在快照上不可区分**（这是本腿自己的读数，167-a1 §1.3② 写了同一件事） | — |

### 1.4 生产调用点枚数（这些字段有没有真的被填）

- 占用字段本身：**0 枚调用点**（字段不存在，谈不到填）。
- 分母候选：`internal/agent/loop.go:253` `Loop.Budgets()` 导出读口，`cmd/wisp` 产码调用者 3 枚——`cmd/wisp/run.go:1056`、`:1062`、`:1109`（本程复认，与 167-a1 一致），三枚都只取 `.PromptTotal` 或整包喂 `agent.NewSpiller`，**没有一枚取 `.ContextWindow` 递出**。`Budgets.ContextWindow` 声明在 `internal/agent/budgets.go:54`。
- 分子候选：`internal/agent/compress.go:110` `Compressor.TotalTokens(hist)`（调用者 `:120/:175` 等，全在 `internal/agent` 包内）；`internal/agent/loop.go:256` `Loop.History()` 导出，cmd/wisp 产码调用者 **0 枚**（尺：`grep -rn "\.History()" cmd --include=*.go` 零命中）；`internal/agent/guard.go:132` `Guard.TokensUsed()`（跨轮累积，包内）。
- ⚠ **但 loop 变量不是 runtime 常驻对象**：`loop` 在 `cmd/wisp/run.go:1023` `execute()` 里造出，是**每任务局部变量**；`agentRuntime` 结构（`cmd/wisp/run.go:266-377`）没有任何字段握着 loop 或其 Budgets/History。任务结束后 loop 无处可读 ⇒ 占用条若要"任务进行中"实时读，得先有人把 loop 的读口挂到 runtime/装配根上——这是 167-a1 §3.2 没有点破的一跳（它只说"在 execute() 里够得着"，但 execute 是单次任务周期，不是常驻 reader 的家）。

## 2. 问二：排队序号（"第几项"这个数能不能拿到面板）

> G3 警告逐字有效：**不许建议加 Loop 方法**。本节只量"现有什么缝可用"；`Loop` 的导出方法面一枚不碰。真正的序号真源也不在 `Loop` 上，在 `approval.Queue` 上——下面逐枚量。

### 2.1 "第几项"这个数：真源在、已在产码被填、最后一跳被丢

| 层 | 现量 | file:line |
|---|---|---|
| **真源①：每卡序号** | `position(it)`＝1-based FIFO 位次，`LiveApprovals()` 逐行填进 `LiveApproval.Position`（字段声明 `:49`，填充 `:118`）。**这是产码非测试路径** | `internal/agent/approval/queue.go:268-275`；`internal/agent/approval/pending_read.go:104-122` |
| **真源②：队列深度** | `Queue.Depth()`＝`len(q.pending)`，导出；产码调用者 1 枚（`cmd/wisp/approval_reply.go:377` 的 CLI 提示行）。pump 侧**没接**（`cmd/wisp/run.go:699-726` 装配无 Depth 槽） | `internal/agent/approval/queue.go:133-137` |
| **真源③：卡上深度** | `Prompt.Depth`＝`g.q.Depth()`（gate 举卡时填）；`PanelItem.Depth`＝`q.position(it)`（`viewLocked` 里 `:531`）。**`PanelItem` 一枚都不进 Snapshot**——A551 已核（平行路），唯一生产消费者是 console 的 `head/view`（`cmd/wisp/approval_reply.go:375/:384`） | `internal/agent/approval/gate.go:529`；`internal/agent/approval/queue.go:524-533` |
| **载体** | `NativeVerdict`（pump 读队列的载体）字段名册＝`CorrelationID/Tool/Args/CallChain/Level/RulesHit/Reason/SessionOverrideBlocked`——**没有 Position**（`sed -n 77,91p` 现量） | `internal/panel/pump.go:77-91` |
| **断点** | `liveVerdicts()` 逐字段抄 `it.Decision`/`it.Corr`，**`it.Position` 一个字节没抄**。因为载体没那枚字段 | `cmd/wisp/panel_pump.go:58-76` |
| **卡面** | `ApprovalCardView` 也没有 position 键（字段册＝correlationId/tool/args/level/rulesHit/reason/reasonKnown/sessionOverrideBlocked/callChain/decidedBy，`internal/panel/approval.go:39-59`） | `internal/panel/approval.go:39-59` |

**结论**：`LiveApproval.Position` 就是"第几项"，真源在产码、非测试形状；从真源到面板断在两处——`NativeVerdict` 无字段（`pump.go:77`）＋`liveVerdicts()` 未抄（`panel_pump.go:65-73`）。这与 167-a1 §1.2② 判读一致（本程逐行复核无误）。

### 2.2 用户消息队列（"我发的第几条"）：明写拒绝，不是漏

- 常驻腿单槽规则：`cmd/wisp/resident_task_source_windows.go:411-414` `if src.running { … return errors.New("本进程此刻已有一发任务在跑：常驻腿一次只接一发，等它结束再提交") }`；理由逐字在 `:396-400`（"a queue of tasks in a leg that cannot show a page would be a queue nobody can read"）。
- 环内插话通道 `Loop.Steer()`（`internal/agent/loop.go:275-283`）**调用者＝0 枚**（尺：`grep -rn "Steer" cmd internal tools --include=*.go` 只命中定义、`SteeringEnabled` 配置行与测试；本程复认）。⚠ 提醒：`Steer` 不是 `Loop` 的新增方法——它已存在，问题是零调用。
- 序号旁证：`cmd/wisp/resident_task_source_windows.go:159` `src.seq`（`:416-417` 自增），只进 goroutine 名与审计行（`:422/:425/:437`），不外递。
- **判定**：用户消息队列＝**盘上没有**（是决定不是缺口），推翻它属产品裁定，本腿不选形。

### 2.3 现有可用的缝（不加 Loop 方法、不加新依赖边）

1. `rt.gate.Queue()` 已是 cmd/wisp 产码既有的读法（`cmd/wisp/panel_pump.go:62` `rt.gate.Queue().LiveApprovals()`、`cmd/wisp/approval_reply.go:377` `s.gate.Queue().Depth()`）——**队列深度 reader 不需要任何新口**，`Depth()` 是现成导出方法，在 cmd/wisp 侧直接可调。
2. 每卡 `Position` 已随 `LiveApprovals()` 的返回值递到 `panel_pump.go:64` 的 `it` 手里——**抄一行就到**，前提是载体给字段（载体改动见 §4）。
3. `PumpSources` 加一枚 reader 槽是既有形状（`internal/panel/pump.go:121-188` 12 槽先例；`cmd/wisp/run.go:699-726` 已在填 9 枚）。
4. ⛔ G3 红线复核：`.scratch/wisp/probes/154/gate-clauses.sh:363-367` G3 腿的尺是 `^func \(l \*Loop\) [A-Z][A-Za-z0-9]*\(.*taskID`（`internal/agent`，want quiet want_n 0）——本腿三枚接线的任何形状都**不落在 internal/agent**，G3 不涉；但任何未来"给 Loop 加收 taskID 导出方法"的形状都先停手上报（HANDOVER 18:5x 裁定在册）。

## 3. 问三：停止按钮（入向停口现在有几个、各缺哪一跳）

> 票面原话点名 `requestTaskStop`／`RequestWorkspaceSwitch` 那族"带判据零调用方"的入向。**本程现量：`requestTaskStop` 在本仓任何根（cmd internal tools docs scripts）＝0 命中**（尺：`grep -rn "requestTaskStop" cmd internal tools docs scripts`）——那是票面写作参考的外部项目符号，本仓没有对应物。本仓同族（"处理函数写好了、生产调用者零枚"）的入向停口/请求口，本程枚到 **4 枚**，逐一列下。

### 3.1 名册：4 枚"已写好、零生产调用方"的入向/停口

| # | 口 | 真源 | 生产调用者 | file:line |
|---|---|---|---|---|
| 1 | `Loop.Steer(text)`（忙时插话，非停但同族"写好了没人调"） | 已存在 | **0 枚** | `internal/agent/loop.go:275-283` |
| 2 | `RunningTask.Cancel()`（异步任务的取消口） | 已存在 | **0 枚**（cmd/wisp 全树 `grep "\.Cancel()"` 只命中 observe/root 侧，无 RunningTask） | `internal/agent/loop.go:300` |
| 3 | `RequestWorkspaceSwitch`（工作区切换请求） | 已存在 | **0 枚产码**——2 枚全在测试（`cmd/wisp/panel_pump_test.go:357`、`cmd/wisp/instructions_200r2_test.go:124`）；它前面的 socket `ComposerDispatch.Workspace` 在产码装配里**写死 nil** | `internal/panel/workspace.go:76`；socket=`internal/panel/composer_dispatch.go:76-78`；nil 装配=`cmd/wisp/panel_inbound.go:275` |
| 4 | `Options.Control ControlHandler`（宿主控制词执行体） | 已存在 | **0 枚产码赋值**（尺：`grep -rn "Control:\s" cmd internal tools --include=*.go` ⇒ 只命中 `internal/llm/anthropic/request.go` 的 CacheControl，无关） | `internal/agent/loop.go:158` |

### 3.2 今天真正能"按停"的两枚，都不在面板上

1. **Esc 快捷键否决卡**：`cmd/wisp/resident_approval_windows.go:171-184` `vetoByEsc()`——**只否决正在等人的那张卡**（`ra.cards.AwaitingHuman()` 为假时答"按下的取消键没有可否决的确认项"，`:177`）；注入点 `cmd/wisp/resident_windows.go:163`。⚠ **它停的是审批卡，不是正在跑的任务**。
2. **退出序列第 3 步**：`cmd/wisp/resident_approval_windows.go:283` `cancelTaskRoots`，注册于 `cmd/wisp/resident_windows.go:186`——那是关机，不是用户按停。
3. 模型侧 `task.cancel`（`internal/tools/task.go:625`）是 C4 工具不是面板入向，且**根任务从不 `AttachCancel`**（产码唯一调用者＝孩子侧 `internal/tools/subagent_197.go:346`）；对根 id 调 `TaskRoster.Cancel` 走 `internal/tools/task.go:458-459` 答"没有可停的句柄"。

### 3.3 面板按钮缺哪几跳（逐枚口各缺的跳）

**共同缺的上半截（面板→Go）**：入向名册 6 枚（`internal/panel/bridge.go:42-45` 四枚 `panel.*`＋`:66-67` 两枚 config），**无一枚与停止有关**；`knownComposerMethod` 的 case（`bridge.go:146-152`）与派发表（`internal/panel/composer_dispatch.go:175-210`）都没有停的字。新增入向方法＝C17 契约面，HANDOVER `:464` 逐字裁定"停手上报，不许自己加"——本腿不选形，只把成本量清（§4）。

**枚 1（Steer）**：缺①调用者（谁在忙时把用户的话送进来——这先要面板有 message 通路，而 `ComposerDispatch.Message` socket 产码也写死 nil，`cmd/wisp/panel_inbound.go:277`）；且 Steer 只插话不停任务，**与停止按钮不是同一枚出口**。
**枚 2（RunningTask.Cancel）**：缺两跳——①没有人握着 `*RunningTask`：`bg` 是 `cmd/wisp/run.go:1099` `execute()` 局部变量，任务结束后不可达，`agentRuntime`（`run.go:266-377`）无字段存它；②没有入向方法把它递进面板链。
**枚 3（RequestWorkspaceSwitch）**：只缺一跳——`ComposerDispatch.Workspace` 处理器（interface `composer_dispatch.go:76-78` 已声明，native leg `workspace.go:76` 已写好）；在 `newComposerDispatchChain`（`cmd/wisp/panel_inbound.go:271-281`）里把 nil 换成真处理器即通。**这是四枚里唯一"下一跳已具名到函数"的**。但它接的是工作区切换，不是停止。
**枚 4（Options.Control）**：缺①装配根一行赋值（`defaultControl`（`loop.go:524-537`）已能对 `l.current` 做 stop/cancel，宿主执行体是可选增强）；②入口：`MatchControl` 唯一产码调用点在 `internal/agent/loop.go:341` `run()` 开头，控制词必须以**新任务输入**的形态到达，而常驻腿忙时拒绝第二发（`resident_task_source_windows.go:411-414`）⇒ **正忙时控制词到不了 loop**，这正是票面 AC#4"停止必须正忙时可达"的机制性障碍。
**根任务停的句柄**：`TaskRoster.AttachCancel`（`internal/tools/task.go:420`）存在，但根任务从未挂（`cmd/wisp/run.go:1099-1105` 只 `RunAsync`+`MarkRoot`，没 `AttachCancel(bg.ID, …)`）。给根挂＝在 `cmd/wisp/run.go` 装配处补一行（`bg.Root().Cancel` 或闭包），不动 `internal/tools` 一字——这是"让名册知道根可停"的最小一跳，但它是**出向能力**，离"面板按钮"仍隔整个入向契约面。

### 3.4 判据钉（现量）

- 正向钉（退出支，非用户按停）：`cmd/wisp/resident_task_source_246_windows_test.go:183`、`cmd/wisp/resident_task_source_live_246_windows_test.go:238`（读到定义为止，未跑）。
- 票面 AC#4 要的"正忙时可到达"与"停止后不留还在写的尾巴"：**无钉**（尺：`grep -in "stop" internal/panel/*_test.go cmd/wisp/panel*_test.go` 零命中）。

## 4. 最小改动面（不选形）

> 前提（167-a1 §3 已立的硬规矩照抄有效）：装配根 cmd/wisp 是唯一接缝（`internal/panel/subagent_roster_197.go:43` 逐字），internal/panel 不 import internal/agent/tools；票面 AC#7 把 internal/panel 写面限在"票 145 已批解冻范围（composer/pump/panel_pump），只到新增字段为止"。下列每枚只描述**改动面**，选不选、选哪形归编排者。

**① 排队序号（最小、真源已在产码）**
- `internal/panel/pump.go`：`NativeVerdict` 加一枚 `Position int`（新增字段，在解冻范围内）。
- `cmd/wisp/panel_pump.go`：`liveVerdicts()` 逐字段抄写里补一行 `Position: it.Position`（`it` 已在手，`:64`）。
- `internal/panel/approval.go`：`ApprovalCardView` 加 `Position int json:"position"` ＋ `CardViewFromDecision` 带上（`:72` 一处）。
- ⚠ 撞 `Q-51`：`ApprovalCardView` 的 JSON 键在双向尺射程内（`internal/panel/approval_test.go:105`、`internal/panel/composer_test.go:49`）——加键需要页面声明同批，而 `frontend/**` 本票禁写。**这是三枚里唯一需要 Q-51 先答的形状**。
- 队列深度标量（"共几张"）：`PumpSources` 加 `QueueDepth func() int`＋`cmd/wisp/run.go:699` 装配 `QueueDepth: func() int { return rt.gate.Queue().Depth() }`——真源 `queue.go:133` 是现成导出方法，**零新口、零新依赖边**；落点是 ComposerState 或顶层新节（同为新增字段）。

**② 占用条（两枚前置、然后是纯装配根活）**
- 前置裁夺（产品裁定，167-a1 §5.2 已摆，本腿不再开）：分子说"当前历史成本"还是"跨轮累积"；分母用缩放后窗口（`budgets.go:54`）还是目录叶。不答则**整枚不做**（宁缺毋造）。
- 若做，改动面：`cmd/wisp` 侧需要 loop 的常驻读口——今天 `loop` 是 `execute()` 局部变量，最小面是 execute() 里装配 `PumpSources` 之前把两个读数闭包交给 runtime（或给 runtime 加字段握 loop 的 reader），**在 cmd/wisp 内完成，不动 internal/agent**。`loop.Budgets()`（`loop.go:253`）与 `loop.History()`（`:256`）＋`agent.ApproxTokensOfMessage`（`budgets.go:150`）全是现成导出——**零 Loop 新方法，G3 不涉**。
- `ComposerState` 加占用节（指针＋omitempty 照 `composer.go:74/:91` 先例）＋分子分母各配 `…Known`（"未知"具名态，票面 AC#2）。

**③ 停止按钮（分"看得见"与"按得下"两半）**
- **"看得见"半**（快照诚实说"有没有一发在跑/能不能停"）：真源＝`cmd/wisp/resident_task_source_windows.go:154` `src.running`＋`:159` `src.seq`，但那是常驻腿对象，`wisp run` 腿的对应物是 `agentRuntime` 自身没有的 running 维。最小面：run.go 装配处给 runtime 记 bg.ID/running 态（`MarkRoot` 之后补挂 `rt.tasks.AttachCancel(bg.ID, bg.Root().Cancel)`——一行，让根可停成为事实），快照经由既有 `Tasks` reader 的 row status/statusReason 表达（根行 status 今天恒 unknown 的洞顺带补上，机制=`internal/tools/task.go:349-357`）。全部在 cmd/wisp＋PumpSources 既有槽内。
- **"按得下"半**：新入向方法＝C17 契约面（HANDOVER `:464` 死令）＋`guardRosterOf` 派发链（`internal/panel/l2_grant_boundary_test.go:2030` 双向推导尺，case 须保持一行逗号分隔 `:2128-2130`）——**本腿判：上报，不自决**。技术上不动冻结件可行（A481 形），但治理归 owner。

## 5. 量不到的格子

1. **两把双向尺今天红不红**——⛔ 禁跑 Go 命令（本票硬规矩）＋`frontend/src/lib/panel.ts` 两层禁令，`Q-51` 那面镜子的当前颜色判不了；在册读数（HANDOVER `:467` 红名册）自述"引用要带锚点、本轮不复量"⇒ 当过期。后果：§4① 加键"会不会多红"只能给射程不能给结论。
2. **`Snapshot.Tasks` 根行 status 若补填，会不会撞枚"无 reader 恰好发四键"字节钉**——那四枚钉（167-a1 §3.2 点名 `pump_test.go:123/:291` 等）钉的是顶层键集；根行 status 是 section 内的值不是键，理论不撞，但未实测（不能跑），只记射程。
3. **resident 腿的 `src.running` 能否被 `wisp run` 腿的快照读到**——两腿各有独立 `assembleRuntime` 实例（`resident_task_source_windows.go:265`），常驻腿的快照 pump 挂在 resident 自己的 runtime 上；跨腿状态共享不存在，本腿只量到"两腿各自组装"这一层，**常驻腿面板（resident panel 线程，`cmd/wisp/panel_resident_windows.go:183-205`）与 resident runtime 快照 pump 之间今天没有装配连线**——常驻面板的 disp 链（`newResidentComposerDispatch`，`:166-171`）只有入向六门，**没有接任何快照 pump**。这一格本腿现量到：常驻进程的面板今天既不发快照也不显示快照内容，票 33 只接了窗口与入向。
4. **`Q-51`/`A383` 台账原文逐字**——台账在 `docs/reports/pending-and-issues.md`（千行级），本腿只取了 A551-A553 三条（§6 引用），未逐行检索 Q-51 段；其内容以 167-a1 §3.1 的转述为准（该件读过原文），本腿不再复读。

## 6. 推翻前人哪句

1. **推翻/修正 167-a1 §3.2 的一句**："两枚都已在 `cmd/wisp/run.go` 的同一枚 `execute()` 里够得着（`run.go:1023` 造出 loop…）"。读数属实但不完整：`loop` 是 `execute()` 的**局部变量**（`run.go:1099` 一带），execute 每任务调一次、返回即不可达；`agentRuntime`（`run.go:266-377` 全字段过目）没有任何字段握 loop/Budgets/History。⇒ "够得着"只在 execute() 存续期成立；占用条若要快照 reader（快照在 execute 外也被 `consoleApprovalUI.Prompt` 等驱动，`run.go:1344-1389`），**必须先有一跳把 loop 的读口挂上 runtime**——167-a1 的四步方案缺这枚具名跳。这不是推翻它的结论（占用分子分母读口导出是真的），是补它漏掉的第 0 步。
2. **修正 167-a1 §1.1② 的一句**："最近的现成来源＝…两枚都已在产码"里分母 `Loop.Budgets()` 的产码调用者 3 枚——本程复认**属实但三枚全部只取 `.PromptTotal`**（`run.go:1056/:1062/:1109` 逐枚看过），其中 `:1109` 是把**整包 Budgets** 递给 `agent.NewSpiller`（a1 §4.2 自己在第 4 节修正过这句）。维持其修正后读法：`ContextWindow` 从未被单独递出，但整包已出过包一次。
3. **对本仓派单链的一枚新读数（前人未量）**：`requestTaskStop` 在本仓 0 命中——票面 line 4 把它与 `RequestWorkspaceSwitch` 并列成"那族"，但后者本仓真有（`workspace.go:76`）、前者纯是外部项目符号。下一程若照票面字面 grep `requestTaskStop` 会白找；正确的同族名册是 §3.1 那 4 枚。
4. **不推翻的**：167-a1 对审批序号"最便宜的一枚"判读（§2 复核属实）、停止"Esc 只否决卡不停任务"判读（§3.2 复核属实）、单槽拒绝非漏（§2.2 复核属实）——三格本程逐行独立复核，维持原文。
