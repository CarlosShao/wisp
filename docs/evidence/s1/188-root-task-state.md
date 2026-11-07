# 188-r4 · 监控栏状态源 · 第一步：根任务状态维（写腿）

- 时刻：2026-11-07 18:5x +08
- 腿：188 写腿（本波 `cmd/wisp` 唯一写面；188 取证腿 `9ad02eb50` 裁定顺序①＝根任务状态维）
- 起手锚：head `2fec9f8787f76f1dd2bf43f417c1a3202b9e80ad`（dev，`2b3bc6839` 之后）
- 写面：`cmd/wisp/run.go` ＋ 新判据件 `cmd/wisp/task_state_188_test.go` ＋ 本证据件 ＋ 票面 Progress log 一行。⛔ internal/ball／internal/tools／internal/config／frontend／design，零字节。
- 铁律：只 commit 不 push；两处显式 pathspec；临时件只建不删；改动文件 LF。

## 1. 起手三核（现跑读数，逐字；⛔ 抄取证腿）

### 1.1 取证腿读数①「MarkRoot 刻意不落 State」——复核成立

`internal/tools/task.go:417-425`（本程现读，逐字）：

```
// MarkRoot files the root identity for a task the host is running itself.
//
// State is deliberately NOT written here: the root task's producer still files
// no D43 name (ruling A394 keeps that rewiring with ticket 196), so
// StateAnswer keeps saying 「宿主没有登记这一维」 for a root row. What lands is
// the parent link (none) and the kind, which is what makes the tree readable.
func (r *TaskRoster) MarkRoot(taskID, label string) {
	r.Record(taskID, TaskOutput{Label: label, Kind: TaskKindRoot})
}
```

`:419` 起的注释逐字对上取证腿转述；函数体确无 State。生产侧根行 `State == ""` 是刻意形状，注释把生产者改写归 196（A394）；196 已 done，落的只是 memory 侧映射（1.3），名册根行今天仍是「宿主没有登记这一维」。

### 1.2 取证腿读数②「组合根构造 TaskBackfill 不带 State」——成立，行号漂移 +5

取证腿写 `run.go:1855`；现量 **`cmd/wisp/run.go:1860-1863`**（本程现读，逐字）：

```go
_, why := (tools.TaskBackfill{
	Roster: rt.tasks,
	Spills: agent.NewSpiller(filepath.Join(rt.spec.dataDir, "artifacts"), loop.Budgets()),
}).Backfill(bg.Root().Ctx, res.TaskID, res.Text)
```

构造不带 `State`（零值）⇒ 落表的根记录 `State == ""`。语义未变，纯行号漂移：取证腿自己的证据件 `docs/evidence/s1/188-monitor-state-source.md:22`、`:43` 也写 1855——锚点之后树又动了，不是转述错。

结构位置：所在函数 `func (rt *agentRuntime) execute(parent context.Context, task string) int`（`run.go:1622`），即 `runTextTask`（`run.go:204`）的真身，`wisp run` 与常驻腿共用。紧邻 `rt.tasks.MarkRoot(bg.ID, task)` 在 `:1845`、`res := bg.Wait()` 在 `:1846`。

### 1.3 取证腿读数③「TaskStateFor / taskLogState 在场且零自造」——复核成立

- `internal/memory/taskstate.go:145-157`（本程现读）：`func TaskStateFor(s string) (statemachine.State, bool)`——收 memory 旧词与 D43 名两种词；off-roster 回 `ok=false`，不翻译。
- 映射表 `taskstate.go:95-103`（票 196 落的表，本程零自造）：`running→Thinking`、`done→Settling`、`succeeded→Settling`、`cancelled→Muted`、`interrupted→Muted`、`stuck→Stuck`、`error→Error`。
- `cmd/wisp/run.go:2182-2184`（本程现读，逐字）：`func taskLogState(s agent.Status) string { return string(agent.TaskLogState(s)) }`——票 365 的一行委托。
- `agent.TaskLogState`（`internal/agent/loop.go:1145-1157` 现读）：completed/control→`done`、stuck→`stuck`、cancelled→`cancelled`、truncated→`interrupted`、failed/未知→`error`——七种终局**全部落在 196 的名册内**（`taskStateRoster`，`taskstate.go:72-87`），`TaskStateFor` 对它们全部 ok=true，故 `ok=false` 分支在本装配点不可达。

### 1.4 票面判据件指认过期（现量在 internal/tools，不在 panel）

票面写「`internal/panel/task_state_188_test.go`（含 `TestTaskState188AC3PanelHasNoWriteLeg`）」。本程 glob＋grep 现跑：**`internal/panel/task_state_188_test.go` 不存在**。现量＝`internal/tools/task_state_188_test.go`（`TestTaskState188AC3PanelHasNoWriteLeg` 在 `:351`；188-c2 证据件与票 188 Progress log 写的也是 tools 路径）。⇒ 票面此点过期，与 `9ad02eb50` 报的「票面现量整段过期」同族。该件在 internal/tools，⛔ 名单内，本程零触碰。

### 1.5 「红/绿是别人的在飞件」排查

- `git status --porcelain`：cmd/wisp 两枚未跟踪测试件（`panel_resident_323_windows_test.go`、`panel_resident_assembly_323_test.go`）＋ docs/evidence 五枚未跟踪件——全是别家的，本程 pathspec 不带，测试读数也不挂本程头上。
- 本波在飞写面 `internal/ball/**`（wave68＋wave66）：本程零触碰；门禁读数与它分开报。
- 写面互斥核对：`run.go` 无别家未提交改动（干净），panel/config 无本程字节。

## 2. 落地方案（邻件已核，无打红）

装配＝`run.go:1860` 的 `TaskBackfill{}` 补一枚 `State:` 字段：

```go
rootState, _ := memory.TaskStateFor(taskLogState(res.Status))
_, why := (tools.TaskBackfill{
	Roster: rt.tasks,
	Spills: ...,
	State: rootState,
}).Backfill(bg.Root().Ctx, res.TaskID, res.Text)
```

**import 增量：1 枚**＝`github.com/CarlosShao/wisp/internal/statemachine`（run.go 现量未引；`memory`/`tools` 已在——票面预估「1-2 行 import」的下限）。

为什么这个形状：

- D43 名只能从 196 的对表取（⛔ 自造）：`TaskStateFor` 是该表唯一读口；返回类型 `statemachine.State` 本身把取值集钉在 D43 20 名上（`internal/statemachine/states.go:11-31`）；`Backfill` 落表前还有 `statemachine.Valid` 门（`task_backfill.go:113-121`），Off-roster 的词整条拒绝并带原因。
- 恒有值：`taskLogState` 七种终局全在 196 名册内（1.3），本装配点 `ok=false` 不可达；即便将来某词掉册，`TaskStateFor` 回零值 State，`Backfill` 的 Valid 门拒绝落表——第三个词不会从这个口子漏进这一维。

### 2.1 诚实边界（现跑核出，如实记）

- **无正文的终局落不进名册，是既有 fail-closed 门，本程不修**：`Backfill` 的 `text == ""` 分支（`task_backfill.go:111-112`）在带 State 之前整条拒绝；`finish`（`loop.go:1053-1092` 现读）不清 `res.Text` 但 failed/stuck/cancelled 路径 `res.Text` 通常为空 ⇒ 这些终局在名册里**没有行**（连「未登记」也没有）。本程只把「有正文的终局」（completed）从「未登记」改口为真 D43 名（`Settling`）。要覆盖无正文终局须动 Backfill 语义（internal/tools）或旁路补写，均非本步射程。
- **终态补丁、无运行态**：`MarkRoot`（`:1845`）到 Backfill（`:1863`）之间根行仍是「未登记」。取证腿裁定本步只做终态一行装配；运行态＋取消兜底是 188-c2 的方案乙（18-24 行三处写点），需要改判 197 的根行钉（`subagent_carrier_197_test.go:317-327`），按规矩要编排者/`A##` 具名一句，本程不动。
- **Backfill `why != ""`（含取消）分支不补写**：同上，归方案乙。
- **`wisp run` 有 store 无 store 都落名册**：名册是进程内 `tools.TaskRoster`，与 `rt.store` 无关；本装配不依赖 store。

### 2.2 邻件不打红（逐枚现跑核过）

- **197 根行钉 `cmd/wisp/subagent_carrier_197_test.go:324-327`**：断言根行 `StatusKnown=false && Status=="" && StatusReason!=""`。取材＝**最后一枚留存包**（`driven.lastPanelSnapshot()`，`:249`）。快照只在循环运行中经 `consoleSink.Publish`（`run.go:2230-2280`，changed 事件驱动）与审批卡/常驻展示点发布（`run.go:1662` 构造 sink；`resident_present_windows.go:197` 属常驻路径）；循环内最后一次 publish 是 EvDone（`run.go:2269-2275`）——全部早于 `bg.Wait()`（`:1846`）返回、早于 Backfill 的 Record（`:1863`）。⇒ 末枚包里根行仍是「未登记」，钉不红。本程新用例在 f.run 返回后不再触发任何 publish，同样不产生更晚的包。
- **`Record` 保留语义**（`task.go:271-301` 现读）：身份字段（ParentTaskID/Label/Kind）空值保留，**State 不保留**——Backfill 落 State 不会抹掉 MarkRoot 已落的 Kind/Label；反向（MarkRoot 不带 State 重落）也不会抹掉 Backfill 的 State，但本程不动 MarkRoot。
- **面板读侧零改**：`taskRosterState`（`panel_pump.go:216`）经 `StateAnswer` 读根行，状态值从「宿主没有登记这一维」变真名——纯数据变化，panel 包零字节；AC#3 钉（`TestTaskState188AC3PanelHasNoWriteLeg`）扫的是面板包写腿形状，与读数值无关。

## 3. 用例与牙（步3 设计）

- 新判据件 `cmd/wisp/task_state_188_test.go`（package main）：
  - 正例 `Test188RootRowFilesTheD43NameTheLoopSpelled`：真组合根（`newRunFixture`＋mockllm）跑一个 Completed 终局（不派子代理、不卡审批）→ 经 `rtHook` 拿 runtime，run 收口后 `rt.tasks.Look(rootID)`：`State == Settling`（D43 名，非 memory 旧词、非空）；同记录过 `StateAnswer()` 得 `(Settling, "")`；再从 `rt.taskRosterState()` 的根行得 `statusKnown=true` 且 `status=="Settling"`（面板读数同一句真话）。
  - 牙 M1：文件级突变——把 `run.go` 装配里 `State: rootState` 那行摘掉 → 正例必须红（根行回到「未登记」）。红行逐字记录；sha256 核对还原字节一致；还原后复跑转绿。
- 反向判据（AC#3）不新增：`TestTaskState188AC3PanelHasNoWriteLeg`（tools 包现量 `:351`）在场，面板包本程零触碰，判据照旧成立。
- 预定红登记：`go test ./internal/panel` 里的 2 枚前端 TS 声明红（`Go Snapshot emits [commands]` 族）不是本程面，逐字登记不修。

## 4. 门禁读数

（步4 后填）

## 5. 诚收

（收件时填）
