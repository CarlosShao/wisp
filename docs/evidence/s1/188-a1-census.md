# 188-a1 — 票 188 只读普查：任务监控栏两堆（子代理状态 / 后台任务状态）从哪取、197 覆盖到哪、剩下那一格落在哪

- 腿：`188-a1`（只读设计核·零产码）。工单＝票 188 `AC#1`／`AC#2`／`AC#3`／`AC#4` 的落地前料刷新。
- 锚：`HEAD c7820bebd`（分支 `dev`），2026-11-04 16:0x-16:2x +08，本程现跑。
- 产出：本文件一枚（仓库内唯一写点）。临时件在 `D:\tmp\wisp\20261104\188a1\`（只建不删）。
- ⚠ 本程**零产码**、未改任何 `*.go`、未改票面、未动 `docs/PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`／C17 白名单。
- 下文所有 `file:line` 都是本程现读；标〔现跑〕的是本程跑过命令/测试的读数，标〔读断言〕的是只读源码/断言文本、未执行。

---

## §0 结论摘要（先给答案，证据在后）

1. **两堆都从同一枚现成结构取，不需要新读面**：`panel.TaskRosterState`（`internal/panel/subagent_roster_197.go:94`），
   由 `cmd/wisp/panel_pump.go:191 taskRosterState()` 从 `rt.tasks`（`*tools.TaskRoster`）读出，经
   `panel.TaskRosterSectionFrom`（`:211`）渲染成快照的 `tasks` 节（`internal/panel/composer.go:91`）。
   子代理＝`rows` 里 `kind == "subagent"` 的行；后台任务（本轮 `wisp run` 那枚根任务）＝同一 `rows` 里 `kind == "root"` 的行。
   **不重造**：状态词、状态判据、行、池数、流截断数、被丢键、点击进页的 `streamKey` 全在。
2. **188 已被票 197 覆盖到「字段＋载体＋反向门」**：`AC#1` 归口 197（`Q-71` 乙）；`AC#2` 的字段（`TaskOutput.State`）、
   枚举（D43 二十枚，`statemachine.Valid`）、进快照（`tasks.rows[].status / statusKnown / statusReason`）、
   「值来自真名册」那一发判据（`internal/tools/task_state_188_test.go:118`）**全部在场**；`AC#3` 反向门**现跑 PASS**。
3. **只剩一格**：`AC#2` 的**生产写入者**——`MarkRoot` 明确不写状态（`internal/tools/task.go:349-357`），
   `TaskBackfill` 的唯一生产构造点没带 `State`（`cmd/wisp/run.go:1674-1677`，字段在 `internal/tools/task_backfill.go:75`）。
   ⇒ 今天线上每一枚**根行**读起来都是 `statusKnown=false`＋「宿主没有登记这一维」，只有子代理行有真 D43 名。
   这一格就是编排者在 `A398` 排的 **`188-r3`＝生产者那一跳**，它当年卡在「任务态 -> D43 名字」是规格真空。
   ⚠ **那个真空今天已经被票 196 补上了**：`internal/memory/taskstate.go:95-103`（memory 词 -> D43 的**全映射**）、
   `:149 TaskStateFor`、`:169 (TaskLog).D43State`，并且 schema v4 的 CHECK 由同一张名册生成。
   ⇒ `188-r3` 今天不再需要新裁定，只是 cmd/wisp 里的一跳。
4. **判语**：**188 可封（归口 197／已实现）**，但封之前要两件事同时落：
   (a) `AC#2` 那一格补上生产者（≈6-24 行，见 §3；两个案型里 **乙** 有价值）；
   (b) `AC#4` 的分账**改写**（原文写"载体归票 145"，现实是 197 开的 `tasks` 键，145 在这一维上**零欠账**）。
   `AC#1`／`AC#3` 已闭；`AC#5` 只要按 §3 的落点（只碰 `cmd/wisp/run.go`）就成立；`AC#6` 只看门禁读数（§4 已知红照实记）。
5. **不新增任何顶层快照键、不新增 C17 方法**：`tasks` 键已在（`composer.go:91`），点进去靠行自带 `streamKey` 去同一包的
   `results` 里取正文——`internal/panel/subagent_roster_197.go:31-39` 逐字写明这**不是** C17 变更。
   ⚠ 界面侧 `tasks` 这一跳**已经落地**（§4 第 7 条现跑读数），票 197 面里「界面那一跳仍红」那段是历史状态。

---

## §1 ① 两堆分别从哪个现有结构取

### 1.1 逐堆来源表

| 堆 | 现成结构（真相源） | 读面（谁读、在哪渲染） | 需要新读面吗 | 判语 |
|---|---|---|---|---|
| **子代理状态** | `tools.TaskRoster` 的行（`internal/tools/task.go:176`），行值 `TaskOutput`（`:100`，状态维 `:130`，判据 `:158 StateAnswer`）；身份/父链由 `PublishSubagent`（`:342`）与子环路自己的准入钩子两次发布（`internal/tools/subagent_197.go:338`、`:373`），终态由 `finalize` 落同一行（`:462-465`） | `cmd/wisp/panel_pump.go:191 taskRosterState()` → `panel.TaskRosterState{Rows}`（`subagent_roster_197.go:94`）→ `TaskRosterSectionFrom`（`:211`）→ `Snapshot.Tasks`（`composer.go:91`），装配点 `cmd/wisp/run.go:930` | **不需要**（复用） | 归口票 197，已实现 |
| **后台任务状态** | **同一枚** `tools.TaskRoster` 的**根行**：`cmd/wisp/run.go:1672 rt.tasks.MarkRoot(bg.ID, task)`（`internal/tools/task.go:355`），终态入口 `internal/tools/task_backfill.go:105 Backfill`（唯一生产调用 `run.go:1674`） | **同一枚** `tasks` 节、同一 `rows` 数组（`kind == "root"`），字段同 §1.2 | **不需要**（复用）；缺的是**写值**不是读面 | 载体已覆盖，生产值缺（§2 AC#2-e） |
| （旁证，不是本票来源）后台任务的**历史＋时间** | `memory.TaskLog`（`internal/memory/models.go:58-70`：`StartedAt`／`EndedAt`／`State`／`Cost*`），读面 `internal/memory/dao_tasklog.go:95 ListTaskLogs`，D43 门 `internal/memory/taskstate.go:169` | `internal/panel/**` 今天**零命中**（`grep TaskLog internal/panel` 为空） | 要用就得新读面 | 见 §3.3：**188 的 AC 文字没有要求它**；截图那列时间是外部产品，不是规格（`docs/evidence/s1/182-rail-stacks-census-a1.md:163` 逐字） |

### 1.2 `tasks` 行今天已经带了什么（逐字段，取自 `internal/panel/subagent_roster_197.go`）

- 行：`TaskID`／`Label`／`Kind`（`root`／`subagent`）／`ParentTaskID`（根为空串）／
  `Status`（D43 名，非 D43 名一律清空）／`StatusKnown`／`StatusReason`／`StreamKey`（子代理＝`subagent:<id>`，根＝`<id>`）／
  `BlockedOnApproval`（与同包 `pending`、L1 窗口 join）／`StreamTruncated`／`StreamElidedRunes`／`StreamDropped`（`:111-143`，渲染 `:268-272`）。
- 节：`rows`（总是发、可空、按 taskId 排序）／`inFlightSlots`／`poolCap`（＝`tools.MaxConcurrentSubagents`）／
  `streamTruncated`／`streamElidedRunes`／`droppedStreamKeys`／`droppedCount`（`:155-184`，装配 `:223-231`）。
- 读端「未登记」的两种形状：空值回「宿主没有登记这一维」；非 D43 名回带原因的拒绝（`internal/tools/task.go:158-168`），
  载体在 `statusKnown=false` 时把 `Status` 清空并保底一句原因（`subagent_roster_197.go:273-282`）。

⇒ **"子代理状态＋后台任务状态"这两堆今天共用一枚键、一枚节、一张行表。** 界面侧只能读**一个** `tasks`，
不许再开第二枚键：`182-a1` 当年建议的「进 K2 的 `Feed`、别开两枚键」（`182-rail-stacks-census-a1.md:209-210`）
在形状上已被 197 落成 `tasks`，那两条的**精神**（一颗键、两堆同源）成立，只是键名不是 `Feed`。

---

## §2 ② 188 逐格 vs 现有 `TaskRosterState`：是不是已被 197 部分/全部覆盖

### 2.1 逐格对比表

| 188 格 | 票面要求（逐字要点） | 现在有什么（`file:line`） | 判语 |
|---|---|---|---|
| **AC#1** 子代理这一维要不要进产品 | 先答要不要 | 已答＝乙（`Q-71`／`A397`／`A399`），实体＋流＋载体三层入库（`internal/tools/subagent_197.go`、`panel/subagent_roster_197.go`、`internal/streamkey/**`） | **覆盖，归口 197**；188 本格已勾 |
| **AC#2-a** 给 `Record` 补状态维 | `TaskOutput` 要有状态字段 | `internal/tools/task.go:130 State statemachine.State`；写侧 `task_backfill.go:127`；`Record` 的身份保留语义 `task.go:211-234` | **覆盖**（188-r2） |
| **AC#2-b** 枚举逐字对齐 D43，不许自造 | 二十枚，逐枚过 `statemachine.Valid` | `internal/statemachine/states.go:11-31`（二十枚）＋`:39 Valid`；判据 `internal/tools/task_state_188_test.go:76`（枚数 20、逐枚 Valid、十一枚旧词必须被拒） | **覆盖** |
| **AC#2-c** 送进快照 | 状态要进包 | `tasks.rows[].status / statusKnown / statusReason`（`subagent_roster_197.go:116-125`，渲染 `:259-282`，装配 `cmd/wisp/run.go:930` + `panel/pump.go:322-339`） | **覆盖**（载体是 197 的 `tasks`，**不是** 145） |
| **AC#2-d** 「值来自真名册，不许硬编码」 | 判据要钉住 | `task_state_188_test.go:118`（同表三条记录三个答案，`Record` 换名读数跟着走）＋读侧 `cmd/wisp/panel_pump.go:216`（只走 `StateAnswer`） | **覆盖** |
| **AC#2-e** 生产里真有一枚后台任务带状态 | 票面 `现量` 表的"缺什么"那一列 | **不成立**：`internal/tools/task.go:349-357 MarkRoot` 逐字「State is deliberately NOT written here」；`cmd/wisp/run.go:1674-1677` 构造 `TaskBackfill` 只给 `Roster`／`Spills`，`State` 恒零（字段注释 `task_backfill.go:60-67` 自己写着"wiring a real name is one line in the composition root"） | **唯一未覆盖格＝`188-r3` 生产者那一跳** |
| **AC#3** 反向判据：面板不许写这一维 | 一枚常驻钉＋正控 | `internal/tools/task_state_188_test.go:351`；〔现跑〕**PASS**：分母 **20** 枚 `internal/panel` 非测试源、**0** 写腿、五枚正控逐枚红、反向控（只读腿）不红 | **覆盖，现跑仍绿** |
| **AC#4** 与票 145 分账 | 快照载体扩张归 145，本票只交"源" | 现实：载体是票 **197** 的 `tasks` 键（`composer.go:91`），`internal/panel/composer.go` 的六键里没有 145 新增这一维；145 在这一维上**零欠账** | **需改账**（不是产码缺口）：把"欠 145"改成"197 已交载体、188 欠生产者" |
| **AC#5** 契约轴 | `PLAN.md`／specs／阈值／golden／C17 白名单一枚不许动 | 本格要求的落点（生产者）只碰 `cmd/wisp/run.go`；`statemachine/states.go`、`bridge.go:42-45`、`thresholds.go`、golden 都不必动 | **成立**（前提：按 §3 落） |
| **AC#6** 门禁 | 三包 `go test`／`d22scan`／`gate-clauses`／面板两枚已知红照实记 | §4（本程只现跑了三条最相关的尺，其余按已知读数记，未冒充复跑） | **按 §4 记，不许当绿** |

### 2.2 判语（逐字回答票面问的那句）

> **「188 是不是已被票 197 的 `taskRosterState` 部分/全部覆盖？」**
> **答：功能面上是「全部覆盖」，生产值上是「差一格」。**
> - 字段、枚举、载体、反向门、与界面那一跳的键集：**全部已由 197 落地，`taskRosterState()` 就是 188 要的那枚读面，不要重造**。
> - 唯一还差的**不是**新结构，而是**没有人往根行里写值**：今天 `tasks.rows` 里根行的 `statusKnown=false`。
> - ⇒ **188 可封**，封的条件是：① 补 `AC#2-e`（`188-r3` 生产者那一跳，落点见 §3）；② 把 `AC#4` 的分账从"145 载体"改写为"197 载体＋188 生产者"（改账归编排者）。
> - ⚠ 与 197 的**重叠**（避免双记账）：`AC#1` 已经整体搬进 197（`Q-71` 乙），188 不再自留；`AC#2-a..d`＋`AC#3` 是 188 的**历史交付**（r2），
>   197 只是把它们的值搬上线；`AC#2-e` 是 188 独有的欠账，197 的 AC 里没有一格要求"根行有状态"
>   （197 `AC#2` 只要求"子代理状态只能是 D43 那一集"，`AC#4` 只要求阻塞态可见）。

---

## §3 ③ 落点：文件、约行数、C17、顶层键

### 3.1 唯一落点

**`cmd/wisp/run.go` 的 `execute()`（`:1449`）里那三行装配（`:1666-1680`）**，不碰 `internal/panel/**`、不碰 `internal/tools/**`、
不碰 `internal/memory/**`（只用 `memory.TaskStateFor` 这个**已有**函数做词表桥）。

两个案型（都**不**新增方法、**不**新增字段、**不**新增顶层键）：

| 案型 | 写什么 | 约行数 | 线上读数 | 红 |
|---|---|---|---|---|
| **甲（终态补丁）** | 只给 `TaskBackfill{...}` 带 `State:`：`memory.TaskStateFor(string(agent.TaskLogState(res.Status)))` | **6-8 行**（无需新增 import：`memory` 已在 `run.go:53`，`TaskStateFor` 的返回类型就是 `statemachine.State`，赋值处不出现类型名） | 跑完才有状态；**跑的过程中根行仍是"未登记"** ⇒ 栏里"回复中"这一形**画不出来** | 197 的根行钉（`cmd/wisp/subagent_carrier_197_test.go:317-327`）**不红**——末枚留存包的构建时点早于 Backfill（`run.go:2031-2041` 的 EvDone 发布 / `:1672` 的 MarkRoot / `:1673` 的 Wait 三者时序）〔读断言〕 |
| **乙（运行态＋终态，推荐）** | ① MarkRoot 之后一行 `Record(bg.ID, TaskOutput{State: Thinking})`（身份字段由 `Record` 保留，`task.go:224-234`）；② `TaskBackfill{State: 终态}`；③ **取消/空正文那一支的兜底**（见 §4 雷 3）：`Backfill` 回 `why != ""` 时补一枚 `Record(bg.ID, TaskOutput{State: 终态})`（此时表里还没有正文，抹不掉东西） | **18-24 行**（三处写点＋注释；import 情况同甲） | 开始＝`Thinking`，结束＝`Settling`／`Muted`／`Error`／`Stuck`，取消＝`Muted` | **会打红** `subagent_carrier_197_test.go:317-327` 那枚钉（它钉的是"根行状态必须空＋带原因"，写于 `A394` 的"生产者还没有"）；同批要把它**改判**成新真话（根行必须有一枚 `statemachine.Valid` 的名字且 `statusKnown=true`）——⚠ 这是**改 197 的钉**，按规矩要有编排者/`A##` 的具名一句，不是实现腿自己动手 |

**为什么推荐乙（按机主铁律：最有价值的、最主流 harness 的做法，不是风险最小的）**：
任务监控栏的**最低可用**语义是"现在有谁在跑"，不是"跑完了才出现一行"。甲只会让栏里永远只有子代理在动、
根任务恒显示"状态不知道"——这正是 `PLAN.md:1127` 记的那类"被阻塞/被卡住的东西在所有界面上不可见"的同族缺陷。
乙的代价是**一条旧钉的改判**和**三个写点**，买到的是"根任务从起跑到收口都有 D43 名"，这是主流 harness 任务栏的标准形状。
且 `A394` 当时不让根行写状态的理由是"词表还没桥"（`internal/tools/task.go:353-354` 逐字"that producer rewiring is ticket 196's"），
**票 196 已落（`memory/taskstate.go`）⇒ 那条理由今天已经不存在**，不是被绕过。

### 3.2 C17／顶层键／契约测试

- **C17：不加任何方法名**。点击进页＝行自带 `streamKey`（`subagent_roster_197.go:126-131`），正文在同一包的 `results` 里
  按 `correlationId == streamKey` 取；`subagent_roster_197.go:31-39` 与票 197 面 §「要 owner 带给界面那支的一跳」末条都逐字写着
  "入向那一跳今天不需要"。若哪天要"宿主记住当前格装的是哪枚会话"，那才是 C17 变更，要另落 `A##`。
- **顶层键：不加**。`tasks` 已在 `internal/panel/composer.go:91`；`Snapshot` 今天六键（`pending`／`results`／`composer`／
  `generatedAt`／`instructions`／`tasks`）。§3.1 的两案都**不改键集** ⇒ 不撞 `TestComposerContractTypesMatchFrontend` 的顶层两向比对。
- **嵌套字段**：§3.1 也不需要加（状态三兄弟已在）。**只有当栏要"开始时间/结束时间"时**才需要动形状，而那是另一件事（§3.3）。
- 装配侧（`cmd/wisp/run.go:930 Tasks: rt.taskRosterState`）与泵侧（`internal/panel/pump.go:322-339`）**一行不用改**：
  读者在、键在、行在。

### 3.3 「截图里那列时间」（09:41／昨天 16:05）——不在 188 AC 内，别顺手做

- 票面 `现量` 表把截图当**示意数据**（`180-c1` 的定性）；`182-a1:163` 逐字：owner 那张截图是**外部产品（qoder）的界面，不是本仓规格**。
- 仓内**没有任何一处**要求任务栏显示时间（`182-a1` 表③给 K9 的那一格只写 `status`）。
- 若哪天要，唯一现成源是 `memory.TaskLog`（`StartedAt`／`EndedAt`／`State`／`Cost*`，`ListTaskLogs` 按 `started_at DESC`，
  `internal/memory/dao_tasklog.go:95-113`），D43 门已在 `memory.TaskStateFor`／`(TaskLog).D43State`（`taskstate.go:149/169`）。
  **代价也具名**：那张表**没有** `label`／`parentTaskID`／`streamKey`／正文（只有 `SummaryText`），
  和 `tasks.rows` 拼起来要一次 join（按 task id），且它按"一次 run 一行"（含子任务行，`run.go:1558-1574`）落库；
  要把它搬进快照＝**新读面＋新嵌套字段**（顶层键可仍复用 `tasks`），属新票，不属于 188 的 `AC#2`。

---

## §4 ④ 雷区与已知红（逐条带 `file:line`）

### 4.1 会立刻咬人的五条

1. **197 的根行状态钉**（本案型乙的头号代价）：`cmd/wisp/subagent_carrier_197_test.go:317-327`
   `if root.StatusKnown || root.Status != "" || root.StatusReason == "" { ... want the empty stated WITH a reason }`。
   注释自己写着理由来源是 `A394`（"filed by nobody"）。乙会让它红 ⇒ **改判要具名**（改成"根行必须有一枚 D43 名、`statusKnown=true`"），
   ⛔ 不许为过尺把它删掉或放宽（`AGENTS.md` §1.1）。同族还要顺带看 `:301-307`（子代理那半，`statemachine.Valid` 钉）与 `:353-359`（收口后 `Settling` 钉）——这两条别动。
2. **`Record` 是 last-writer-wins，且只保留身份三字段**：`internal/tools/task.go:211-234` 逐字
   "Text, ArtifactPath and State are NOT preserved"。⇒ 终态必须**与正文同一次**落表（走 `Backfill` 的 `State` 字段，
   `task_backfill.go:127` 就是那一次写），或**只在 `Backfill` 拒绝的那一支**用 `Record` 补（那时表里还没有正文，抹不掉东西）；
   ⛔ 决不许在 `Backfill` 成功之后再补一枚只带 `State` 的 `Record`——那会连正文与续读路径一起抹掉（`task.output` 立刻读空）。
3. **`Backfill` 有两条不落表的路，兜底不做就是"永远回复中"**：`internal/tools/task_backfill.go:111-125`
   `text == ""` 与 `Stopped(ctx)`（被取消）都**拒绝写**，且都不是错误、只回一句 `why`（`run.go:1678-1680` 只打审计）。
   ⇒ 取消的根任务、空正文的根任务在乙里会停在 `Thinking`——**必须在 `why != ""` 那一支补一枚状态写**，否则栏上是一句谎。
   ⚠ 这条同样撞 `task_backfill.go:111-112` 的既有语义（"空正文不进名册，免得「查不到」和「没打印」被读成同一件事"）：
   根行本来就由 `MarkRoot` 建在册里，补写不破坏"空正文不**新建**行"这条；但**不许**把空正文写成 `Settling`（那才是把"没打印"说成"完成"）。
4. **词表桥只有一处，别自己 switch**：`memory.TaskStateFor`／`TaskLog.D43State`（`internal/memory/taskstate.go:149/169`）是
   memory 七词 -> D43 的**全映射**（`running`→`Thinking`、`done`/`succeeded`→`Settling`、`cancelled`/`interrupted`→`Muted`、
   `stuck`→`Stuck`、`error`→`Error`，`:95-103`）。把 memory 原词直接塞进 `TaskBackfill.State` 会被 `:113` 拒（并留下"没落表"的读数）。
   ⛔ 也不要顺手改 `task_log.state` 的写法或统一两套词（票 196 已定"两套并存、一处翻译"，`taskstate.go:24-33`）。
5. **D43 那二十枚里没有"已完成""已拒绝"**：`internal/statemachine/states.go:11-31` 是球/会话的态名；任务的终态语义靠上表映射。
   - "回复中"＝`Thinking`；
   - "已完成"＝`Settling`（与子代理 `subagentStateSettled`、`taskstate.go` 的 `done→Settling` 同一口径，`internal/tools/subagent_197.go:107-110`）；
   - "已取消/被停"＝`Muted`；
   - **"已拒绝"没有任务级对应物**：它是**审批裁决**（卡片的 reject/超时），不是任务状态。界面若要那一列，得与同包 `pending`/裁决账 join，
     这属于票 220 那一族（`blockedOnApproval` 是读，不是裁决），⛔ 不许为它新造第 21 个态名（D43 冻结，`188-r1` 的 ① 已量过）。

### 4.2 反向门／不许越的线

6. **`AC#3` 那扇门是"只扫 `internal/panel`"**（`internal/tools/task_state_188_test.go:249-348`，四种形状：
   调 `Record`/`Backfill`、构造带 `State` 键的 `TaskOutput`、给 `.State` 赋值、`Method*` 常量声称 task/state 轴）。
   ⇒ 生产者**必须**留在 `cmd/wisp`（本程 §3 正是），⛔ 别为了"顺手规整"把状态判定/赋值搬进 `internal/panel`：
   那会同时撞 R20（面板只显示＋发起请求）与这枚钉。载体现有的 `view.Status = ""`（`subagent_roster_197.go:278`）是 `.Status` 不是 `.State`，
   〔现跑〕扫描器不咬它；但任何 `.State =` 形状都会咬。
7. **`TestComposerContractTypesMatchFrontend` 今天的红因不是 `tasks`**〔现跑〕：红句逐字
   `Go ComposerState emits [files] that interface ComposerState does not declare`（`internal/panel/composer_test.go:78`），
   同时 `Snapshot <-> PanelSnapshot: 6 JSON keys reconciled` ⇒ **页侧已声明 `tasks` 与 `instructions`**。
   ⇒ 票 197 面里「界面那一跳仍红（`PanelSnapshot` 缺 `tasks?`／`instructions?`）」那段是**历史状态**，别照它去改 `frontend/**`
   （那是机主外部 agent 的写面，本编队零写面）。本程 §3 两案都不碰键集，**与这枚红无关**。
8. **`TestC21DesignTokensFourWayAgree` 今天的红因也不是 tokens.css 被删**〔现跑〕：红句是
   `docs/evidence/s1/c21-native-tokens.md:58-60/:103-105 the table names Palette.GlassRing / GlassHi / GlassHiSoft, but internal/ball/tokens.go holds no such colour literal ...`
   （六行）。⇒ 188 面 `AC#6` 那句"`design/assets/tokens.css` 被别家删了未 staged"也是**历史成因**；照实记、不许顺手修。
9. **面板源分母**：AC#3 的扫描要求 ≥10 枚非测试源（`task_state_188_test.go:356-358`），〔现跑〕今天 **20** 枚 ⇒ 分母不可能因删文件掉到 10 以下；
   但若有人把 `internal/panel` 的源改名成 `_test.go`，这枚门会静默变宽——不是本票射程，只记名。

### 4.3 门禁／工程面

10. **共享编译面纪律**：`cmd/wisp/**` 与 `internal/panel/**` 每波只许一枚写腿；`go test ./cmd/wisp` 重包同一时刻只跑一个
    （`AGENTS.md` §时代索引③）。票 197 的载体测试是**真 run**（超时 240s 级，`subagent_carrier_197_test.go:159`），
    改根行状态会动到它 ⇒ 落 `188-r3` 时先跑定向 `-run TestRunPacketCarriesTheSubagentItsRosterRowFed`，别整包乱跑。
11. **`MarkRoot` 本身别改**：`internal/tools/task.go:349-357` 的 doc 与 `subagent_depth_fixture_287`／`task_164_ac4_cancel_tail`
   一族都按"根行只有身份"的口径造夹具；改 `MarkRoot` 的写值语义＝改 `internal/tools` 的公共面，收益不抵风险。生产者放 `cmd/wisp` 一侧
    只用现成 `Record` 即可（`Record` 保留 `ParentTaskID/Label/Kind`，`task.go:224-234`）。
12. **`run.go` 现无 `statemachine` 导入**（`:48` 只到 `agent`、`:53` 到 `memory`）——但 §3.1 两案**都不需要加**：
    直接写 `State: memory.TaskStateFor(...)` 即可（返回类型就是 `statemachine.State`，不写类型名就不欠 import）。
    ⛔ 也别在 `cmd/wisp` 里再写一张 memory→D43 的 switch（那会造出第二处翻译，`taskstate.go:24-33` 明写只要一处）。
13. **别动 `task.output` 给模型读的 stub 措辞**：`internal/tools/task.go:555-568` 逐字写着"frozen by PLAN.md:2564 / ticket 164's accepted AC#3"，
    `188-r2` 也一字未动。
14. **`go test ./internal/panel` 今天两枚已知红**（本程未修、未当绿）：
    - `TestComposerContractTypesMatchFrontend`（页侧缺 `ComposerState.files`，见雷 7）；
    - `TestC21DesignTokensFourWayAgree`（文档表 vs `internal/ball/tokens.go` 三个 glass 字段，见雷 8）。
    ⚠ `188-r2` 当年还点名过第三枚 `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`；本程**没有**跑它，不替它报绿/红。

---

## §5 没测到／边界（照实写，不当已证）

1. **没有跑 `cmd/wisp` 任何包**（含 §3 两个案型会撞的那枚根行钉）——落点行数是**读断言后估的**，红/绿是**推断**，标〔读断言〕；
   真值要 `188-r3` 落码那一程现跑。
2. 没有跑 `go test ./internal/tools/ ./internal/agent/ ./internal/panel/` 全量、`d22scan`、`gate-clauses`：
   本程只跑了三条与"这一维"直接相关的定向测试（§4 的〔现跑〕三处）；票面 `AC#6` 的其余读数属实现腿。
3. `tasks.rows` 只有"本进程 admit 过的任务"（`cmd/wisp/panel_pump.go:169-176` 逐字：roster 没有 list-everything 口，
   走 `seenTasks` + `Descendants`）。⇒ **跨重启的历史**、以及"这一程没 admit 过的行"看不见；188 的 AC 没要求，只在这里记成边界。
4. `Record` 的后写覆盖语义（雷 2）在**并发**下还取决于写序；本程只读了单线程语序，没有做并发压测。
5. "已拒绝"这一列怎么画，本程只指出它**不是任务态**；具体 join 归哪张票（220 族）没替它裁。
6. 界面侧字段（`PanelSnapshot.tasks` 的内部行形状）本程**只在 Go 侧**核；`frontend/**` 零读零写（`A383`），
   所以"界面是否已按行渲染 status/statusKnown"这一问**本程无读数**。

---

## §6 纪律面

- 仓库内写点：**只有本文件**。`git commit` 用显式 pathspec（`-- docs/evidence/s1/188-a1-census.md`），未 `push`；
  未 `-A`／`--amend`／`reset`／`stash`／`checkout`；仓内零删除；未建 worktree。
- 临时件：`D:\tmp\wisp\20261104\188a1\readings.md`（原始读数底账，只建不删）。
- `docs/PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`／`internal/panel/bridge.go:42-45`／票面：**一字未动**。
- 本程跑过的读命令：`git rev-parse`／`git branch`／`git status --porcelain`（只读）、三条定向 `go test`（只读，未生成仓内件）、
  `Get-ChildItem` 计数。**未跑** `d22scan`／`gate-clauses`／全量包。
