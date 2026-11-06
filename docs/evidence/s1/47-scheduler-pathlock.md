# 票 47 — TaskScheduler + PathLock (C20)：实现腿裁决表

- **腿名**：W47（wave45 写腿）
- **日期**：2026-11-06 22:3x +08
- **HEAD（起手）**：`aa7159b29`（票面锚 `ac7d14c7f` 的后一枚；本腿全部判据在起手 HEAD 现测）
- **前置**：10-agent-loop-core `-done`（盘上 `10-agent-loop-core-done.md`）；21-approval-gates-minimal `-done`（`21-approval-gates-minimal-done.md`）——与台账 `A1588` 的"前置核查（现跑）"一致，票真解锁。
- **写面**：`internal/agent/scheduler/**`（owning package）＋ `internal/tools`（task 家族接 scheduler 面，票面 AC#6 的工具半边）＋ 本件。
- **禁改面**：`docs/reports/pending-and-issues.md`（未动一字）、`frontend/**`、`design/**`（未动一字）、PLAN.md / specs（未动一字）。
- **GUI/进程**：零 GUI、零 wisp.exe、零 *winlive*（本腿全程未跑任何窗口进程）。

---

## 0. 起手三核（逐条）

1. **票面全文读毕**（`.scratch/wisp/issues/47-task-scheduler-pathlock.md`，44 行）：What to build / Key constraints / Out of scope / 6 条 AC 全部过目。
2. **具名串/行现测**：D43 边 #20/#39 在 `internal/statemachine/table.go:146/:272`（真行号现数）；`EvPathConflict`/`EvPathLockReleased` 在 `events.go:38/:61`；D38c 全文在 `PLAN.md:2838-2842`（"每任务 root context / WaitGroup"段）；D38d（工具并发 =4）在 `PLAN.md:2844-2847`；SPEC-06 §7 PathLock 行在 `docs/specs/SPEC-06-security-gatekeeping.md:109`；「在等 <task> 释放 <path>」措辞按票面逐字落进 `waitingNote`。行号全部现跑现引，无一处沿用票面旧号。
3. **未定义即停检查**：无。票面 6 条 AC 全部可实现且不撞任何 OPEN 项（`AwaitingApproval` 挂起行为、web.search 等均与本票无涉）。

## 1. 落了什么（对 AC 逐条）

| AC | 判据（票面逐字） | 落点 | 测试 |
|---|---|---|---|
| AC#1 并发 | 4 任务 × 混合工具负载全完成；协程数回到基线、零泄漏（registry test） | `scheduler.go`：FIFO accept＋并发跑＋每任务 root ctx＋registry Spawn | `Test47AC1FourTasksRunConcurrentlyAndRosterDrains`（4 任务×trio，收尾 `RosterReport` 全零）；`Test47AC1FIFOAcceptOrderIsTheQueueOrder`（单窗 FIFO 序）；`Test47AC1DrainIsRequiredBeforeCompletion`（D38c 排干法） |
| AC#2 PathLock | 重叠文件操作 → 第二个排队＋正确原因 → 释放后恢复（#20→#39 态迹） | `pathlock.go`：Acquire 阻塞等待＋`waitingNote`＝「在等 <task> 释放 <path>」＋Release 重发队 | `Test47AC2OverlappingPathQueuesSecondWithReasonAndResumes`（排队行 Queued＋原因逐字；释放后完成）；`Test47AC2SameTaskReacquireIsNotAConflict`（幂等半） |
| AC#3 死锁 | A 持 p1 要 p2；B 持 p2 要 p1 → 超时＋显式错误 | `pathlock.go`：`SetWaitCeiling`＋单调 `observe.NewTimeout` 判限＋`ErrDeadlock` | `Test47AC3MutualWaitTimesOutWithExplicitError`（票面原样夹具，150ms 缩窗）；`Test47AC3CancelUnblocksAPathWaiter`（等待可停半：ctx 取消即出队） |
| AC#4 批准独立 | A 等 L2 时 B 照常跑完工具 | per-task `approval-waiter` goroutine 形状＋`MaxConcurrent=2` 并发窗 | `Test47AC4ApprovalWaitDoesNotBlockOtherTasks`（B 完成时 A 的卡还开着） |
| AC#5 并发下取消 | 中飞取消 → applied-steps 报告；兄弟不受影响 | `Cancel`（排队：就地结算；在跑：root ctx）＋行注带 D31 台账措辞 | `Test47AC5MidFlightCancelCarriesAppliedStepsSiblingsUnaffected`（行注含步骤＋「未自动回退」）；`Test47AC5CancelAQueuedTaskSettlesItWithoutRunning` |
| AC#6 往返 | mock LLM 驱动 task.list/cancel | `bridge.go`（scheduler→tools 适配器）＋ `tools/task.go`（task.list 工具＋task.cancel 调度器分支） | `Test47AC6TaskListRoundTripThroughTheRealLoop`（真 loop＋真 bridge，行进模型上下文）；`Test47AC6TaskCancelRoundTripStopsASchedulerRow`（父停子＋外人拒）；`Test47AC6TaskListUnwiredFailsClosed` |

**D32**：本包零新增 timer。唯一的等待（死锁上限）是等待任务自己 goroutine 里的单调 deadline 检查（`observe.NewTimeout`＋`Expired()`），与 `observe.Root.Wait` 同一习惯；关停路径 `CancelAll` 亦同。无 `time.AfterFunc`、无新常驻 goroutine。

**D38b 名册**：唯一的 spawn 是任务根本体，名字 `"agent-task-"+t.id`（per_task 前缀族）；229 普查登记随落（`spawn_roster_census_229_test.go` 新增一行复审记录）。

## 2. task.list 的解冻（AC#6 的契约半边，按 A434 同规）

- done-判据＝`PLAN.md:1537` 登记行：「注册表里 task.list／task.cancel 各有一枚真实现并过契约测试」，前置＝票 163（已闭，`A1588` 记）＋票 164 AC#1 裁定表（在册裁定已有）。两枚今天各有一枚真实现：task.cancel（票 221 甲形，A434）＋ task.list（本腿）。
- 落法＝**标记与接线同一次核销**（A434 为 task.cancel 立的同规）：`BuiltinTaskEntries` 注册第三枚、`task.go` 头注释的 `task.list` DEFERRED 行摘除，五个既有尺**同批翻转**，无一放宽：
  1. `tasklist_deferred_236r3_teeth_test.go` 名集钉：从「不许在册」翻成「必须在册」（negative 半交给 271 枚数哨兵，分工不重复）；
  2. `task_cancel_221_legs_test.go` 枚数哨兵：`!= 2` → `!= 3`（票 271 甲形工序明文许可的合法第三枚路径；工序注释同步更到"第四枚先回 271 重选形"）；
  3. 同文件 `Test221DeferredMarkerForCancelLiftedButListStillMarked`：`listMarked == 0` 反向支翻转成 `!= 0` 报红；
  4. `taskfamily_roster_271_tripwire_test.go` B-3 正控：突变期望 3 枚 → 4 枚（哨兵/名集钉分工的每跑复证保持完好）；
  5. `ticket175r2_stamp_live_test.go` 分类表＋`tier_consistency_172_*` 分母：`task.list` 入表（unstamped：调度器自己的行表，无外来 body；172 落地数 9→10、canonical call 加一行）。
- 契约文字零改动：`PLAN.md` D34 行 / §7 登记行 / SPEC-07 §3 镜像一字未动（行内容本就是 task.list L0；本腿只兑现它）。

## 3. 明说不做的（票面 Out of scope 与簿记）

- **cmd/wisp 装配不动**：共享编译面，本腿起手时盘上有两枚 323 腿的未跟踪测试件（`panel_resident_*_323*`），按令不动 cmd/wisp 即不触发规则 21。scheduler 的产品装配（run.go 建 Scheduler＋适配器包住 bridge）留给装配票写面——本包 `bridge.go` 的 `NewBridgeDeps`/`ExecAdapter` 就是那条缝，AC#6 的真环在包内已用**真 loop＋真 bridge** 走通。
- **面板任务列表/球 Queued 点**：票面把 UI 半边归 145 的载体；`Rows()` 是那条读缝。
- **多任务审批队列路由（48）/会话授权（49）/插件运行时**：票面 Out of scope 逐字未越。

## 4. 门禁读数（全部现跑）

| 门 | 读数 |
|---|---|
| `gofumpt -l`（两触及包） | 空（三文件先被 `-w` 修过后复测空） |
| `go vet ./internal/agent/scheduler ./internal/tools` | rc=0 |
| `go test ./internal/agent/scheduler ./internal/tools -count=2 -timeout 20m` | `ok scheduler 1.211s` ＋ `ok tools 55.548s`（rc=0） |
| `sh scripts/d22scan.sh` | `clean - no D22 ban violations`（internal/=371, cmd/=58, tools=48…）；rc=0 |
| `go build ./cmd/... ./internal/... ./tools/...` | rc=0（`build/` 是 gitignore 的本地 tmp，preexisting，非本腿物） |
| `go test ./cmd/wisp -run 'Test197NoAllowDoor|Test174R8|Test174'` | ok（TaskDeps 走查尺＋task.output 指针尺在加字段后照绿） |
| `go test ./internal/agent ./internal/observe -count=1` | ok（loop/registry 回归） |

**诚实告警**：①`go build ./...` 会把 gitignore 的 `build/` 临时目录编进去红（preexisting，与本腿无关，判据用三棵真包树 rc=0）；②`taskDeps` 加了 `Scheduler` 字段，`subagent_selfapproval_197_test.go` 的 allow-door 走查（reflect 走全部字段）实测照绿——`TaskListView`/`TaskCanceler` 两接口不带任何 allow/answer 动词，这正是当初把面收窄的原因。

## 5. 撤销口令

撤本腿＝revert 本提交一枚（全部改动的落点都在一个 commit 里）；task.list 解冻的单独复原＝`git revert` 后 `task.go` 头注释三行 DEFERRED 原样回场。
