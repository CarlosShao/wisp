# 188 — 面板「任务监控」栏两堆的 Go 侧来源（写腿，2026-11-13）

- 腿：188 面板状态来源（编排者 A1827/A1828 派，独占 `internal/tools/**`＋`internal/panel/**`＋`cmd/wisp/**`）
- head 锚：`8a8673fe8`（A1827 在树上）
- 本程产码：`internal/tools/task_backfill.go`（1 处门改动）＋两枚判据（tools/cmd 各一枚）
- 本程零改动：`internal/panel/**`（核毕不动，理由见 §3）、`task.go`、`run.go`、`subagent_197.go`、前端、design、契约面

## §1 票面行号 vs 今天真实行号（现跑 grep/read 核对，⛔ 直接采信票面）

| 票面写法 | 票面行号 | 今天真实位置 | 判 |
|---|---|---|---|
| `internal/tools/task.go:110-152`（Record／Look／Count 名册） | 110-152 | `Record`＝`internal/tools/task.go:281`、`Look`＝:334、`Count`＝:349 | **漂**（+171～+199；票 197/221/236/287 往该文件加了大量内容） |
| `Record` 在 `:123`（188-r1/派单写法） | 123 | `task.go:281` | **漂**（+158） |
| `TaskOutput` 的 `State` 字段（188-r2 落的） | （r2 时 :124） | `task.go:196`（字段）、`StateAnswer`＝:224 | **漂**（+72） |
| 「写侧门 `task_backfill.go:113`」 | 113 | `task_backfill.go:130`（本程在 :106-116 加注释、:123 改门后 +17） | **漂**（本程贡献 +17 中的注释部分） |
| 188-c2 记「`run.go:1855` 构造 TaskBackfill 不带 State」 | 1855 | 该构造已带 State（r4 落的）：`cmd/wisp/run.go:2185-2189` | **已过期**（r4 已修） |
| 188-s1 记 spawn 拒绝点「`:254/:269/:276/:281/:288/:302-304/:309/:316`」 | 8 处 | 今天同形状仍在：`subagent_197.go:254/:269/:276/:280-283/:287-291/:302-305/:308-310/:315-318`（多行展开，枚数一致） | **基本没漂** |
| AC#3 判据件「`internal/tools/task_state_188_test.go:341/351`」 | 341/351 | 负向腿＝`TestTaskState188AC3PanelHasNoWriteLeg`，现在 :412（本程在其上插入新判据 +61 行前为 :351→:412 漂自他腿） | **漂** |
| 票面「grep Subagent…0 命中」的 09-28 现量 | 0 命中 | 今天 `grep -riE "subagent" -l internal/ cmd/ --include=*.go` ⇒ **19 枚文件**（188-c2 已改判） | **整段过期**（c2 已记） |

## §2 两堆各自的 Go 侧来源（有来源 / 有但缺字段 / 完全没有）

### 堆一：子代理状态 —— **有来源，全链在场**（本程零改动，核毕）

| 环节 | 位置 | 状态 |
|---|---|---|
| 生产者（运行态） | `internal/tools/subagent_197.go:410-415` `PublishSubagent` → `State=Thinking`；双发布（spawn goroutine :373＋子环路 admission :337-345） | 有 |
| 生产者（终态） | `subagent_197.go:450-465` `finalize`：`Settling`（有正文）/`Muted`（被停）/`Error`（无正文） | 有 |
| 状态判定门 | `internal/tools/task.go:224-234` `StateAnswer`（fail-closed，非 D43 名拒付） | 有 |
| 读侧（面板载体） | `internal/panel/pump.go:202` `Tasks func() TaskRosterState` → `pump.go:447-465` 装进 `snap.Tasks`；行载体 `internal/panel/subagent_roster_197.go:51-84`（TaskRow）→ :108-143（TaskRowView，JSON `status/statusKnown/statusReason`） | 有 |
| 装配根读口 | `cmd/wisp/panel_pump.go:207-271` `taskRosterState()`（:232 过 `StateAnswer`，状态读只经宿主判官）；`cmd/wisp/run.go:1201` `Tasks: rt.taskRosterState` 接线 | 有 |
| 「谁在跑」计数 | `task.go:457-484` `InFlightSubagents`/`RunningSubagentIDs`；wire 侧 `InFlightSlots/PoolCap`（`subagent_roster_197.go:94-106`） | 有 |
| **缺口①**：拒绝行 | spawn 八处拒绝（§1 表）全部早退、零名册行；「已拒绝」该借 D43 哪个名＝**规格真空**（188-s1 卡点，A394 后无人裁） | **无源（规格未裁，本程不碰）** |
| **缺口②**：时间维 | `TaskOutput`/`TaskRow` 无钟点字段；188-c2 裁「立票归 197 射程」 | **无源（归 197）** |
| 缺口③：根行运行态 | `run.go:2157` `MarkRoot` 刻意不落 State（`task.go:417-425` 注释：A394 把改判归 196；运行中根行读「宿主没有登记这一维」）；`cmd/wisp/subagent_carrier_197_test.go:324` 把这一诚实空态钉为**正向断言** ⇒ 补它要改判别人的钉＋动 A394 边界，**本程不碰，如实上报** | 有源但刻意空窗 |

### 堆二：后台任务状态 —— **有来源、有字段，今天真缺口＝「空正文终局整条丢行」**（本程修）

| 环节 | 位置 | 修前 | 修后 |
|---|---|---|---|
| 名册/字段 | `task.go:166-207`（`TaskOutput.State`，类型 `statemachine.State`，20 枚 D43 名） | 有 | 不变 |
| 生产者（根任务） | `cmd/wisp/run.go:2184-2189`：`State: rootState`，`memory.TaskStateFor(taskLogState(res.Status))`（对表读数，零自造；`internal/memory/taskstate.go:95-103`） | 有 | 不变 |
| 写侧门 | `internal/tools/task_backfill.go:117-160` `Backfill` | **门 :112 旧形状 `case text == ""` 整条拒绝** ⇒ failed/stuck/cancelled/truncated 终局（环路的 finish 只发布 message 不发布正文，`internal/agent/loop.go:417/:449/:454/:518/:569/:600/:1066/:1075`）**连状态维一起丢**：名册无行、面板读「宿主没有登记这一维」 | `case text == "" && b.State == ""`（:123）：**带 D43 名的空正文终局照登**（State 真、Text 空、无 artifact）；空态＋空正文维持旧拒；非空正文路径逐字节不变 |
| 读侧（模型） | `task.go:571-642` `taskOutput.Execute` | 不变 | 不变 |
| 读侧（面板） | 同堆一（同一载体/同一泵） | 空正文终局行缺席 | 终局行在场、`statusKnown=true` |
| scheduler 那张表（第二源） | `task.go:107-123`（`TaskListRow` 带 D43 State）→ 装配 `cmd/wisp/run.go:815-819` 接线 `scheduler.NewBridgeDeps` | 有（task.list 读它；面板侧不经它，票面「排队中」堆归 213/197 射程） | 不变 |

## §3 技术选型（铁律⓪⓪⓪，我自己拍）

**选了什么：不新开载体、不画新线，修「回填门把状态维跟正文绑死」这一枚缺陷。**

- 主流 harness 的面板（Claude Code / Codex / Cursor）任务监控的样子：**一行一任务、带状态词与父子链、终局状态常驻可回看**——这正是 197 载体层已落成的形状（rows/status/inFlight/poolCap），泵与 wire 都在场。两堆今天**不是**「没有来源」，子代理堆全链在场、后台任务堆差的是**「失败这类终局在名册里整行蒸发」**——这是唯一一个让面板对终局撒谎的点，修它=让既有链路对得起主流形状。
- 为什么它最有价值：不修它，owner 圈的那一栏里「已拒绝/已失败」这类最有用的行**一枚都出不来**（失败＝无正文＝整条被拒）；修它之后，失败/stuck/取消/截断四种终局的名字（`Error`/`Stuck`/`Muted`）全部从**宿主已有的对表读数**（`memory.TaskStateFor`）进名册，零自造词汇、零新契约面、零新 key。相比之下：补根行运行态要改判 A394＋重钉 `subagent_carrier_197_test.go:324`（别人的正向钉）；补拒绝行要先裁 D43 借名（规格真空，D22③ 停手）；补时间维要动 `TaskOutput`/`TaskRow` 字段（归 197 射程）。三件都碰别人的地界或未裁的规格，价值都不如这一枚高。
- ⛔ 没用「动作小/门复杂度」当理由：不修这一枚的代价＝面板终局栏对失败永远显示「宿主没有登记这一维」，那是**谎话**（host 明明知道终局）；修它不改任何冻结面。

## §4 成对判据（两发，命令逐字＋读数）

**正向**（喂两种不同的状态源输入 ⇒ 面板读到的值随之变）：

1. completed 终局（`done→Settling`，既有判据复跑）：
```
scripts/wisp-cli-staged-tests.sh -test.run 'Test188' -test.v
```
读数：`--- PASS: Test188RootRowFilesTheD43NameTheLoopSpelled (1.83s)`／`--- PASS: Test188FailedRunStillFilesItsD43State (1.68s)`（修后树上）
2. failed 终局（新判据，mockllm `QueueFault(500×3)`）：同上命令。读数（`-v` 逐字）：
```
task_state_188_test.go:179: failed run: roster state="Error" text=0 bytes, panel statusKnown=true status="Error"
--- PASS: Test188FailedRunStillFilesItsD43State (2.78s)
```
⇒ 同一读者、两种状态源（completed/failed），面板读到的值跟着变（Settling / Error），非「字段非零」断言。
3. tools 单元面（直接喂门，三种输入）：
```
go test ./internal/tools/ -run 'TestTaskState188' -count=1 -v
```
读数：`--- PASS: TestTaskState188AC2EmptyBodiedEndingKeepsItsState` 等 **5 枚全 PASS**（`ok ... 0.099s`）。

**反向**（把来源那条线摘掉 ⇒ 必须有可观察变化，⛔ 静默）：把门临时改回旧形状 `case text == "":`（mutation，sha 备份 `D:\tmp\wisp\20261113\188-panel\task_backfill.go.keep`，改前 `15AC30EB…7013F0`，复原后同值）：

1. tools 侧：
```
go test ./internal/tools/ -run 'TestTaskState188AC2EmptyBodiedEndingKeepsItsState' -count=1 -v
```
读数（红，逐字）：
```
task_state_188_test.go:245: an empty-bodied ending with a D43 state was refused whole:
"这个任务没有留下正文（空正文不进名册，免得「查不到」和「没打印」被读成同一件事）"
--- FAIL: TestTaskState188AC2EmptyBodiedEndingKeepsItsState (0.00s)
```
2. cmd/wisp 侧（真组合根）：
```
scripts/wisp-cli-staged-tests.sh -test.run 'Test188FailedRunStillFilesItsD43State' -test.v
```
读数（红，逐字）：
```
task_state_188_test.go:151: roster state = "", want Error ...
task_state_188_test.go:158: StateAnswer() = ("", "宿主没有登记这一维（fail-closed：不替任务编一个状态）"), want (Error, empty)
task_state_188_test.go:167: panel row statusKnown = false (reason "宿主没有登记这一维…")
task_state_188_test.go:170: panel row status = "", want "Error"
--- FAIL: Test188FailedRunStillFilesItsD43State (1.48s)
```
⇒ 摘掉来源＝名册行丢、面板从 `statusKnown=true/"Error"` 变回 `false/「宿主没有登记这一维」`——**红在明处，不静默**。

## §5 门禁读数（全绿）

| 门 | 命令 | 读数 |
|---|---|---|
| build | `go build ./cmd/wisp/` | rc=0 |
| vet | `go vet ./cmd/wisp/` | rc=0 |
| test | `go test -count=1 ./internal/tools/ ./internal/panel/` | `ok ... tools 37.335s`／`ok ... panel 20.997s`（**panel 全绿**，票 23 的 2-4 枚已知红今天不在） |
| gofumpt | `gofumpt -l internal cmd tools scripts frontend`（走 `D:\work\base\gopath\bin\gofumpt.exe`） | 无输出（首跑点名 `cmd\wisp\task_state_188_test.go`，已按 fumpt 形状重写后复跑为空） |
| d22scan | `sh scripts/d22scan.sh`（走 Git bash） | **clean**；`bans #1-5 internal/=404, cmd/=77, ban #8 internal/=1005, cmd/=276` |
| ⚠ 票 98 | `go test ./cmd/wisp/` 直跑 | `exit status 0xc0000135`（DLL 不在 PATH，按派单不算红）；改走 `scripts/wisp-cli-staged-tests.sh`（sherpa DLL staged，读数见 §4） |

## §6 AC#4 与票 145／197／196 的分账（谁欠谁）

- 本票交的＝**源**（回填门修好，终局状态进名册）；**载体**（`TaskRosterSection` 十六字段）197 已交、在 wire；**出向传输**（`PostWebMessage|EvaluateScript|CreateWebMessageAsJson`）仍零命中非 test——**票 35 欠页面那半**；**页面读取**（`.tasks` 消费）仍零文件——**前端欠**（机主外部 agent 地界，188-c2 已记，本程复跑口径不变）。
- **票 196**：两套词表的并表已在其名下落（`memory/taskstate.go`），本程只消费其对表读数，未统一任何词表。
- **本票余格**（不归本程）：①根行**运行态**空窗（要改判 A394＋重钉 `subagent_carrier_197_test.go:324`，编排者具名）；②拒绝行（先裁 D43 借名）；③时间维（归 197）。

## §7 诚实边界

- `state TEXT NOT NULL` 之外的 memory 侧一行未动；`task.output` 模型可读 stub 一字未动；`internal/panel/` 本程零字节（核毕结论：该包今天没有缺陷，缺的在生产者那一跳）。
- 撞见「整进程死亡」现象：**0**（`0xc0000135` 是票 98 已定性的加载环境现象，staged 载体绕过，非票 314 形状）。
- 纪律：写面＝`internal/tools/task_backfill.go`＋两枚判据件＋本证据件＋票面 Progress log；零 push；AC 框一枚未勾（勾要非实现者）；临时件只建不删（`D:\tmp\wisp\20261113\188-panel\`）。
