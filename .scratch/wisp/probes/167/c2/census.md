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

（骨架占位——待填。）

## 3. 问三：停止按钮（入向停口现在有几个、各缺哪一跳）

（骨架占位——待填。）

## 4. 最小改动面（不选形）

（骨架占位——待填。）

## 5. 量不到的格子

（骨架占位——待填。）

## 6. 推翻前人哪句

（骨架占位——待填。）
