# 票 47 — TaskScheduler + PathLock (C20)：非实现者裁决表（VERDICT）

- **腿名**：47v（非实现者 VERDICT 腿；实现腿＝wave45 W47，commit `1cc67664b`）
- **日期**：2026-11-06 23:0x +08
- **裁决者**：本腿与实现腿不同 agent（D22 双角色；SPEC-12 §4.3 #1/#3）
- **验尸对象**：commit `1cc67664b`（票 47 全部改动所在）。**码等性现证**：`git diff 1cc67664b..HEAD -- internal/ cmd/ scripts/ tools/` 空 ⇒ 提交里的码就是 HEAD 上的码（HEAD=`7393cc604`）。
- **隔离载体**：规则 21 触发属实——主树 `cmd/wisp/` 盘上有两枚 323 腿的未跟踪测试件（`panel_resident_*_323*`）⇒ 全部复跑在 `git clone --shared` 的克隆（`.scratch/wisp/probes/47v/clone`，HEAD 同 `7393cc604`、working tree 0 未跟踪）里做，主树零改动。
- **禁令遵守**：零 GUI、零 wisp.exe、零 *winlive*（本腿跑过的测试面＝scheduler+tools 全套、cmd/wisp 仅 `Test197NoAllowDoor|Test174` 两把源尺、agent+observe 回归）；`docs/reports/pending-and-issues.md` 只读；`frontend/**`、`design/**` 未动一字；无 push。
- **环境**：CGO_ENABLED=1（全程未置 0）；GOTMPDIR/TMPDIR/TEMP/TMP 全程 unset；PATH 前置 `build` 与 `third_party/sherpa-onnx`；go1.27.1 windows/amd64；Git Bash＝`D:\work\soft\Git\bin\bash.exe`。

---

## 0. 新鲜复跑读数（全部本腿现跑，共享克隆内）

| 门 | 命令 | 读数 | 判 |
|---|---|---|---|
| gofmt | `gofmt -l internal/agent/scheduler internal/tools` | 空 | ✅ |
| gofumpt | `gofumpt -l`（同两包，`/d/work/base/gopath/bin/gofumpt.exe`） | 空 | ✅ |
| vet | `go vet ./internal/agent/scheduler ./internal/tools` | rc=0 | ✅ |
| test | `go test ./internal/agent/scheduler ./internal/tools -count=2 -timeout 20m` | `ok scheduler 1.234s`＋`ok tools 66.550s`，rc=0（-count=2 两遍全绿） | ✅ |
| d22scan | `sh scripts/d22scan.sh` | rc=0 `clean - no D22 ban violations`；正控 45 PASS/0 FAIL/SKIP=0；分母 internal/=371、cmd/=58、ban7 tools/=48、ban8 internal/=923 | ✅ |
| build | `go build ./cmd/... ./internal/... ./tools/...` | rc=0 | ✅ |
| build 尾巴 | `go build ... .`（派单原文带尾点） | rc=1 `no Go files` —— **结构上任何 commit 都不可能过**：仓根没有 .go 文件（`ls *.go`＝无）。判据按三棵真包树取，rc=0 | ⚠️ 派单笔误，非本票缺陷 |
| cmd/wisp 尺 | `go test ./cmd/wisp -run 'Test197NoAllowDoor|Test174' -count=1` | `ok 0.137s` rc=0（首次跑 0xc0000135＝克隆无 gitignore 的 `build/` DLL；把三枚 DLL 复制进克隆的 gitignore 目录后照绿，源隔离不变） | ✅ |
| 回归 | `go test ./internal/agent ./internal/observe -count=1` | ok（7.1s/5.0s） | ✅ |

## 1. 变异试验（派单点名的 HARD 项）

**突变**：从 `internal/tools/task.go` 的 `BuiltinTaskEntries` 删掉 `{Tool: taskList{d: d}, Decl: taskListDecl()},` 一行（克隆内；sed 逐行删）。

| 尺 | 突变下读数 | 判 |
|---|---|---|
| **236r3 名集钉** `Test236R3TaskListDeferralIsPinnedByTheRosterNotByAComment` | **FAIL**：「票 47 把 task.list 接进名册，而现在的注册名册里没有它（现名册：task.output / task.cancel）」——派单要的那枚 RED，红了 | ✅ 有牙 |
| 271 枚数哨兵 `Test221TaskCancelIsRegisteredAtItsFrozenLevel` | FAIL：「task 家族注册了 2 枚」（≠3） | ✅ 有牙 |
| 票 47 两枚接线尺（`Test47TaskListDeclMatchesTheFrozenRow` / `Test47TaskListAnswersRowsAndFailsClosed`） | FAIL（未注册 / `未知工具 task.list`） | ✅ 有牙 |
| 172 两把尺（`Test172FrozenTierMatchesTheTierTheGateGives` / `Test172DocCopyThatLiesReddensTheRuler`） | FAIL（`probe for task.list but it is not registered`） | ✅ 有牙 |
| AC#6 真环 `Test47AC6TaskListRoundTripThroughTheRealLoop` | FAIL（模型上下文里没有 task.list 回执；调用被门拒） | ✅ 真环真依赖注册 |
| 221 标记尺 `Test221DeferredMarkerForCancelLiftedButListStillMarked` | ok——**预期内**：该尺管的是「task.go 里不许留 DEFERRED＋task.list 注释行」，本突变只摘注册行不动标记行；分工与 236r3 头注写的一致 | （说明） |
| 271 B-3 正控 `Test271B3...` | FAIL——**预期内**：正控假设干净的三枚家族（内层哨兵报「4 枚」），复合突变下家族只有 3 枚故断言不中；干净树上它每跑都绿（-count=2 套件与下述还原后重跑均绿） | （说明） |

**还原**：`cp` 回原文件后 `git hash-object internal/tools/task.go` ＝ `git rev-parse HEAD:internal/tools/task.go` ＝ `518235a0873d4b1bc0d33a9a1401ada9dc57c719`，逐字节相等；还原后八枚尺（236r3＋221 两枚＋47 三枚＋271 B-3＋172）同批 `ok` rc=0。突变只建不删、还原逐字节、主树零触碰。

## 2. 五枚被翻的尺——逐枚核「翻 ≠ 放宽」（diff 逐字读＋行为实测）

1. **236r3 名集钉**（`tasklist_deferred_236r3_teeth_test.go`）：断言从「名册里**不许**出现 task.list」翻成「**必须**出现」。这不是删尺：negative 半没丢——`sawOutput/sawCancel` 支原样保留（:171-184），「别的名字误注册从它手下绿过去」的假红半明确移交 271 枚数哨兵＋B-3 正控（头注逐字写明分工），且本腿突变实测它照样红。**未放宽**。
2. **221 哨兵 2→3**（`task_cancel_221_legs_test.go`）：`len(names) != 2` → `!= 3`。数字**向上**走＝更严不是更松；红句更新且明文禁止「就地把 != 3 放宽」；工序注释按 271 甲形更到「第四枚先回 271 重选形」。**未放宽**。
3. **221 标记尺**（同文件 `Test221DeferredMarkerForCancelLiftedButListStillMarked`）：`listMarked == 0 → Errorf` 翻成 `listMarked != 0 → Errorf`。方向＝随能力落地**反转**（能力在册后，留一行说它不可用的注释才是谎话），cancel 半支（cancelMarked）与整尺结构未动。这正是 A434 为 task.cancel 立的「标记与接线同一次核销」同规。**未放宽**。
4. **271 B-3 正控 3→4**（`taskfamily_roster_271_tripwire_test.go`）：突变期望「task 家族注册了 3 枚」→「4 枚」——正控的期望随合法第三枚**向上**跟踪，哨兵红/名集钉绿的分工复证每跑都在（干净树 -count=2 套件两次绿）。**未放宽**。
5. **175r2 分类表＋172 分母 9→10**：`task.list` 以「unstamped：调度器自己的行表，无外来 body」入 175r2 分类表（**加行不是删行**；删行即红的账不变）；172 的 `tier172LandedNumb` 9→10、`tier172Landed` 增 `task.list`、machine 表加一枚 canonical call——分母**向上**（AC#3 防的正是分母缩水），且本腿突变下两把 172 尺实测照红。**未放宽**。

## 3. 契约锚点现核

- **D34 行**：现跑真行在 **`PLAN.md:2570`**（`task.list / task.cancel | 查看/取消任务 | L0 / L1 | —`）。实现腿证据表引 `:2564`，**行号烂了 6 行**（内容对得上，读数以 :2570 为准）——记录为引号腐烂，非缺陷。
- **§7 登记行**：`PLAN.md:1537` 逐字核对，done-判据＝「注册表里 task.list／task.cancel 各有一枚真实现并过契约测试」。该行今天仍标 DEFERRED——台账那半归 owner（票 221 后同状），本腿不认它为阻塞，因为收件三把尺看的是盘上注册＋实现＋契约测试，三样全在（`BuiltinTaskEntries` 三枚、AC#6 真环、172 对账）。
- **D43 #20/#39**：`internal/statemachine/table.go:146`（Acting→Queued, EvPathConflict）／`:272`（Queued→Acting, EvPathLockReleased），两行真在。
- **台账先例**：`A434`（task.cancel 甲形批准，批准记录同先例自落）；`A1588`（wave44 前置核查解锁 47/23/33）；`A1589`（wave45 收 47/23/33 三腿皆 PASS）。本票 task.list 的解冻沿用 A434「标记与接线同一次核销」先例，PLAN.md/specs 零字改动（实现腿证据表如此声称，本腿 `git show 1cc67664b --stat` 复核：写面只在 internal/ 两包＋票面进度两行＋证据表，属实）。

## 4. 六条 AC 逐条 1:1 判决

| AC | 票面判据 | 盘上证据（本腿在克隆内读过源码＋现跑） | 判 |
|---|---|---|---|
| 1 并发 | 4 任务×混合工具负载全完成；协程数回基线、零泄漏（registry test） | `Test47AC1FourTasksRunConcurrentlyAndRosterDrains`：4 任务并发，每任务带 `tool-exec-<id>`＋`approval-waiter` 两枚 per-task 助手（agent-task-<id> 是根本体 Spawn 名），收尾 `RosterReport` 全零（Total=0/ResidentOverBaseline=false/Unknown=0）；另有 FIFO 序尺＋「排干才算完成」尺（helper 未退场 Wait 不返回）。判据读的是本测试自己的 Registry（隔离名册），「6＋3×active」基线按 RosterReport 全零等价兑现；「混合工具负载」以 gated 助手形模拟，真工具负载由 AC#6 真环补 | **PASS** |
| 2 PathLock | 重叠文件操作→第二个排队＋正确原因→释放后恢复（#20→#39 态迹） | `Test47AC2OverlappingPathQueuesSecondWithReasonAndResumes`：B 行 Queued＋行注含「在等 <A的id> 释放」逐字（`waitingNote`＝`"在等 "+owner+" 释放 "+path`，scheduler.go:614-616）；释放后 B completed、HolderOf 空。#20/#39 两边都在冻结转移表（§3）。幂等半（同任务重复 Acquire 不自堵）另有一尺 | **PASS** |
| 3 死锁 | A 持 p1 要 p2；B 持 p2 要 p1 → 超时＋显式错误 | `Test47AC3MutualWaitTimesOutWithExplicitError`：票面原样夹具，`SetWaitCeiling` 150ms 缩窗，两边收尾都含 `ErrDeadlock` 哨兵措辞；上限用单调 `observe.NewTimeout`（无墙钟差，D42#9 习惯）；另有「取消解锁等待者」半尺（ErrWaitCancelled） | **PASS** |
| 4 批准独立 | A 等 L2 时 B 照常跑完工具 | `Test47AC4ApprovalWaitDoesNotBlockOtherTasks`：A 的 approval-waiter 卡开着，B 的工具执行到 completed；批准落下后 A 才完成；waiter 退场读数=0。D31 核心性质被测试直接钉住 | **PASS** |
| 5 并发下取消 | 中飞取消→applied-steps 报告；兄弟不受影响 | `Test47AC5MidFlightCancelCarriesAppliedStepsSiblingsUnaffected`：B 在 A 还在跑时独立完成（先证）；Cancel A 后行注含已生效步骤逐字＋「未自动回退」（D31 措辞）；排队取消半尺（未起跑的行就地结算、Exec 零执行） | **PASS** |
| 6 往返 | task.list/cancel 经 agent（mock LLM 驱动两头） | `Test47AC6TaskListRoundTripThroughTheRealLoop`＋`Test47AC6TaskCancelRoundTripStopsASchedulerRow`：**真 agent.Loop＋真 bridge＋真注册表＋真 C19 assessor**，mock 只在冻结的 C5 LlmProvider 缝（AGENTS §1.3 合规）；task.list 回执（行数头/id/标题/行注）实测进了模型上下文（capture 读 req.Messages）；task.cancel 父停子经真 bridge、外人被拒；未接线 fail-closed 另有尺。tools 侧三枚形状尺（L0 声明/三态应答/权威分支）在 tools 包内同绿 | **PASS** |

**Key constraints 附加核读**（非 AC 行，如实记录）：
- **每任务 LoopGuard/budgets 独立**：结构性成立——`internal/agent/loop.go:415` 每 loop `NewGuard`（每任务根跑自己的 loop ⇒ 天然 per-task）；既有 guard_test.go 8 枚尺照绿。本 commit 未新增「跨任务 LoopGuard 隔离」专项尺；判为沿用票 10 的既有形状，非放宽（无既有尺被删/被松）。
- **工具并发 ≤4（D38d）**：落点在 `internal/agent/budgets.go:47` `MaxToolConcurrency=4`（`guard.go:107` 钳位），pre-existing、本票未动、回归绿。
- **D32 零新增 timer**：pathlock 等待＝自身 goroutine 内 `ctx.Done()` select＋单调 `NewTimeout` 判限；全包无 `time.AfterFunc`。属实。
- **D38b 名册**：唯一新 spawn 形＝`"agent-task-"+t.id`，229 普查表加行复审（加行）。属实。

## 5. closeable 与余债

- **六条 AC 全 PASS（6/6）**，五尺翻转核为「随能力落地的同批反转、无一放宽」，变异 RED／还原逐字节/复绿三段全闭环，五道门禁新鲜复跑全绿 ⇒ **票 47 本体判 closeable**（按收件三把尺：提交真在 `1cc67664b`、件真在 9 文件＋测试、本腿自跑读数全上）。
- **cmd/wisp 装配仍然欠着（Owed，不阻塞本票关闭）**：现跑证据＝`cmd/wisp/run.go:653` 仍是 `tools.TaskDeps{Roster: rt.tasks, Paths: rt.paths}`——**没有 `Scheduler:` 字段**，全 `cmd/wisp/` 无任何 `agent/scheduler` 引用。即生产进程今天 task.list 会答「调度器未接线（fail-closed）」、task.cancel 走名册分支。实现腿也如此具名（disposition：装配缝＝`NewBridgeDeps`/`ExecAdapter` 留给装配票）。这条债**必须在台账里挂名**：不挂，AC#6 的「经 agent 往返」就只剩测试环、没有生产环。
- **面板半边**（panel task list＋球 Queued 点）：票面把它挂在 PathLock 约束句里，UI 载体按实现腿处置归票 145；`Rows()` 读缝已备。记录，不判本票 FAIL。
- **引号卫生**：实现腿证据表的 `PLAN.md:2564`（真行 :2570）行号腐烂一条，已在 §3 记录。

## 6. 撤销口令

- 撤本裁决表＝不认这枚 commit（`docs/evidence/s1/47-verdict.md` 单文件，revert/删皆一句话的事）；
- 实现腿的撤销口令沿用实现腿证据表：revert `1cc67664b` 一枚即全退（注册＋实现＋五尺翻转同批）。
