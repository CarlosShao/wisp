# 188-c2 — 监控栏 S1 堆（子代理／后台任务状态）的状态源字段表 ＋ 下一步裁定（只读·零产码）

- 性质：票 188 的**只读取证＋裁定腿**。本波 `internal/panel/**`＋`cmd/wisp/**` 写面归 213 腿独占，故本程⛔产码；写面只有本件＋票面 Progress log 一行。
- 起手锚：`29b3501d5`（分支 dev，2026-11-07 18:10 +08）。上一程只读腿 188-r1（零产码停报）与写腿 188-r2 的锚是 `e9ef94d0`（09-28）；**票面「现量」整段是 09-28 的读数，今天大面积过期**（子代理维已被 197 落成整链），本程逐枚重跑，见 §2。
- 纪律：`internal/ball/**`（在飞）、`frontend/**`、`design/**`、thresholds/golden、票面 AC 框与 Status——零触碰；只 commit 不 push。

## 1. 起手三核

1. **票面点名的串/行**：`internal/tools/task.go:110-152` 的 `Record`／`Look`／`Count`——载体 `TaskRoster`／`TaskOutput` 仍在场，但行号已漂（`Record` 的非测试写入者现共 4 枚：`task.go:411` PublishSubagent、`task.go:424` MarkRoot、`task_backfill.go:141` Backfill、`subagent_197.go:462` 收尾落终态）。票面「全仓没有子代理这一维的任何数据模型」的尺今天打出来 **19 枚文件**（`grep -rniE "Subagent|sub_agent|SpawnAgent|DelegateTask|PublishSubagent" --include=*.go internal/ cmd/ -l`）⇒ 票面现量过期，原因是票 197 落地。
2. **引文出处**：票面引的普查件 `docs/evidence/s1/180-182-panel-fields-census-c1.md` §4–§8（`:119-230`）逐段读毕，其 S1 行「规格真空」判定已被 197 推翻；census 交的 S2–S7 归属（181／189／190／163／145 行 12）本程未逐一重开，只核与本票 S1 有关的部分。
3. **红/绿归属**：本程零 Go 改动、零测试运行 ⇒ 无红绿可归属。起手 porcelain 里本程未碰任何已跟踪文件；在飞写面（65/66 的 `internal/ball`、213 的 `panel+cmd/wisp`）零干涉。

## 2. S1 堆字段表（字段 | 真源 | reader 有无 | 缺什么 | 补它要动哪个文件）

**载体与读口今天已整链在场**：
`internal/panel/composer.go:91`（`Snapshot.Tasks *TaskRosterSection`，omitempty 指针）← `internal/panel/pump.go:355,372`（`p.src.Tasks != nil` 才发段）← `cmd/wisp/run.go:1019`（`Tasks: rt.taskRosterState` 装配）← `cmd/wisp/panel_pump.go:189-247`（宿主侧映射，状态读**只经** `tools.TaskOutput.StateAnswer`）← 真源 `internal/tools` 名册。

| 字段 | 真源 | reader 有无 | 缺什么 | 补它要动哪个文件 |
|---|---|---|---|---|
| `TaskRow.TaskID`／`ParentTaskID`／`Label`／`Kind` | 名册行：spawn 双发布点 `subagent_197.go:338/:373`（经 `PublishSubagent` task.go:410）；根行 `run.go:1840`（经 `MarkRoot` task.go:423） | 有（`panel_pump.go:218` 全量映射；段内按 TaskID 排序） | 无 | — |
| **子代理 Status**（运行中／终态） | 运行中＝`PublishSubagent` 落 `Thinking`（task.go:413）；终态 `Settling`（有答案）／`Muted`（被停）／`Error`（无答案）由 197 收尾路径落（`subagent_197.go:462` Record 带 State；借名词表 `subagent_197.go:99-108`，全部是 D43 既有 20 名） | 有（`StateAnswer` → `Status`／`StatusKnown`／`StatusReason`，pump.go:222-231） | 无 | — |
| **根任务 Status** | **缺**：`MarkRoot` 刻意不落 State（task.go:419-422 注释逐字＝A394 把生产者改写作归 196；196 已 done 但只落 memory 侧映射 `internal/memory/taskstate.go`，**名册根行仍空**）；终局 `Backfill` 本身支持带 State（`task_backfill.go:113` 校验 D43 20 名）但 `run.go:1855` 构造 `TaskBackfill` 时**不带** | reader 在场，值恒走「宿主没有登记这一维」文案 | 生产写入**一行装配** | `cmd/wisp/run.go:1855` 给 `TaskBackfill{}` 补 `State:`——D43 名从 `memory.TaskStateFor(taskLogState(res.Status))` 取（taskstate.go:145，旧词→D43 名的对表已由 196 落好，零自造） |
| **「已拒绝」行** | **缺**：spawn 四处拒绝（`subagent_197.go:254` 名册未接线／`:281` 深度上限／`:288` 池满／`:316` 未接审批钩子）**全都不留名册行**（上面 4 枚写入者无一在拒绝分支）⇒ 截图里「已拒绝 昨天 16:05」那一堆今天无源 | reader 现成（行一落就上页） | 源侧整维 | `internal/tools/subagent_197.go` 拒绝分支补 Record——⚠ **D43 借哪个名（Queued？Muted？）是规格真空**，先裁再写 |
| **「排队中」行** | 半缺：C20 scheduler 的 `TaskListRow`（task.go 的 `TaskListView`，含 State/Note/HeldPaths）是**另一张表**；面板 Tasks 段只读 `TaskRoster` ⇒ scheduler 排队的任务在监控栏不可见 | 无 reader | 泵侧第二源 | `cmd/wisp/panel_pump.go`（加 scheduler 读口）＋`internal/panel/`（段扩张）⇒ 撞 213 写面与 197 载体射程，**另派** |
| **时间维**（截图「已完成 09:41／回复中 10:18」的钟点） | **规格真空**：`TaskOutput`／`TaskRow` 都没有 started/finished 字段，任何一处都不记钟点 | 无 | 源＋载体双侧新字段 | `internal/tools/task.go`＋`internal/panel/subagent_roster_197.go` ⇒ 新字段＝**立票**（197/145 射程），非装配级 |

AC#3 反向判据复核：`TestTaskState188AC3PanelHasNoWriteLeg`（`internal/tools/task_state_188_test.go:351`）在场；面板载体头注自陈服从该钉（`subagent_roster_197.go` 头注「it does not import internal/statemachine and it does not assign to a State field」）；本程静态核对面板包无 State 赋值／Record 调用形状。**复核成立：状态维今天仍没有面板写腿。**

## 3. 与票 145 的分账复核（AC#4 现量）

票面分账「本票交**源**、145 交**载体**」**已被 197 改写**：

- S1 的载体半边（`Snapshot.Tasks`＋`TaskRow`/`TaskRosterSection`＋泵读口＋`run.go:1019` 装配）由 **197 载体腿落掉**（出处见 §2 链）；145 票面自己的记录也是「`tasks`→票 197」。
- 145 今天对 S1 真欠的**不是 Go 字段**：页面类型 `frontend/src/lib/panel.ts:216,253-256` **已声明 `tasks?` 键**（有人落了对接），但页面读取 **0**（`grep -rn "\.tasks" frontend/src --include=*.tsx --include=*.ts -l` ⇒ 0 文件，本程现跑）；且 Go→页出向传输 **0**（`grep -rnE "PostWebMessage|EvaluateScript|CreateWebMessageAsJson" --include=*.go cmd/ internal/` 剔 test ⇒ 0 命中，本程复跑）⇒ 「能显示」的最后一跳归**票 35 的桥**（PostWebMessage/推事件那一跳），不在 145 也不在 188。
- ⇒ 分账今天的形状：**188＝源**（余量三枚：根任务态／拒绝行／时间维）**· 197＝载体＋宿主映射（已落）· 35＝出向传输（未落）· 前端会话＝页面读取（未落）**。票 145 的 AC#2「页面要、Go 没交」清单对 S1 已清空。

## 4. 顺手更正的一笔过期账（归票 220／编排者核收）

`145-a2`（10-13）登记「`PumpSources.L1Windows` 全仓零装配点」⇒ **今天过期**：`cmd/wisp/run.go:1023` 已接 `L1Windows: rt.liveL1WindowWaits`（reader＝`panel_pump.go:108`，读自 gate 本体）。票 220 的 `blockedOnApproval` 第二来源今天有装配。本程不落产码，只在此具名登记。

## 5. 裁定（下一步写腿该补哪几个字段、按什么顺序）

**① 根任务状态维**＝`cmd/wisp/run.go:1855` 一行装配（`State` 取 `memory.TaskStateFor(taskLogState(res.Status))` 的 D43 名，零新词、判据 `task_state_188_test.go` 现成）——**最有价值且最便宜**，把今天根行恒「宿主未登记」的真缺口关掉；本波 `cmd/wisp` 写面归 213，**等其收口后落**。**② 拒绝行**＝`internal/tools/subagent_197.go` 拒绝分支补 Record（本波无包冲突可先行），但 D43 借哪个名是规格真空，**先裁后写**。**③ 时间维**＝源＋载体双侧新字段，**立票归 197 射程**，不许顺手。「排队中」第二源（scheduler 表进段）撞 213/197 射程，随 ③ 一并立票排后。

## 6. 门禁

零 Go 改动 ⇒ gofumpt／vet／d22scan **无对象可扫**，如实登记。零测试运行（不把在飞腿的红绿记到本程头上）。`frontend/**`/`design/**` 只读（两把 grep），零写面。零 push。

## 7. 没做到的（如实）

- 未复跑任何 Go 测试：零改动无对象；且 65/66/213 三腿在飞，整包读数会混入别人的红。
- 「排队中」第二源的段扩张设计、拒绝行的 D43 借名：只登记不裁（规格真空，未定义即停，交编排者）。
- census §6 那句「宿主侧读面 vs 模型侧工具共享实现」的编排者裁题：本票不重开（S2–S7 归各自票）。
- AC 框一枚未勾（勾要非实现者表）。
