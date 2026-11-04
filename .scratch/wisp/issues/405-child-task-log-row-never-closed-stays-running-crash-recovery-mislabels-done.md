# 405 — 子任务 `task_log` 行开了却从不关：child 恒 `running`，crash_recovery 会把成功的也标 `interrupted`

- Status: **done**
- 立案：2026-11-04 源起票 **378-v1 R3**（378 的修复让子任务获得 `task_log` 行，但 `FinishTaskLog` 惰性 ⇒ 子行永不关）
- 关联：票 **378**（子任务获父行、fkc=0，已解决"无父"）· 票 **373**（级联）· 票 **404**（R4 前置 bug，已修）

## 这是什么（人话）

`cmd/wisp/run.go:1324-1328` 的 `taskLifecycleJournal.FinishTaskLog` 是**有意惰性**（`return nil`）。子环路用的是同一个 journal ⇒ 子任务 `StartTaskLog` 落一行 `running`、但 `FinishTaskLog` 不写 ⇒ **子任务行永远停在 `running`**。

全仓只有 `run.go` 在 `execute()` 收尾为 **root** 任务关闭行（`res.TaskID`）；**子任务没有任何关闭点**。而 `cmd/wisp/crash_recovery.go:145-165` 会把仍是 `running` 的行改判 `interrupted` ⇒ **成功完成的子任务被误标"中断"**。

## 影响

用户可见的隐私页"任务日志"里：**每个子任务恒显 `running` 或被误标 `interrupted`**，与实际成败不符。

## 结案判据

- [x] **AC#1 裁形**（编排者拍）：子任务终态写手（主流 harness：子任务完成应有终态）。候选：①`FinishTaskLog` 对子任务写入终态（task 行半，不动 tool_call 惰性）；②`run.go` 侧在子任务收尾处关闭；③按 `SPEC-12 §5` 五字段登记推迟。
- [x] **AC#2 牙**：一发 `task.spawn` 成功后该子行 `state=done`（或等价终态）；反向＝摘终态写手 ⇒ 用例转红逐名报。
- [x] **AC#3 门禁**（同上票规范）：build/vet/两条 `-run` 用例/gofumpt/d22scan。
- [x] **AC#4 纪律**（同上票规范）。

## Progress log

- [2026-11-04] agent=`编排者` did=立案（票378-v1 R3 具名上报：FinishTaskLog 惰性⇒子行恒 running、crash_recovery 误标 interrupted；独立于 378 的"无父"核心，故另立本票）。

## AC#1 裁形（编排者 2026-11-04）：Option A — journal 自写终态（主流：环路拥有自己 task 行的 lifecycle）
- 机制：loop 经 internal/agent/journal.go:137-145 taskJournal.finishTask 调 j.FinishTaskLog(state,summary,tokensIn,tokensOut,costMicro,errorClass)，落到 taskLifecycleJournal.FinishTaskLog（现 inert，run.go:1324-1328）。root 另由 run.go:1742 关闭；child 无闭点 ⇒ 恒 running。
- 定形：①taskLifecycleJournal.FinishTaskLog 由 inert 改为调 j.store.FinishTaskLog 写终态，加 rowFinished atomic.Bool 置位；②run.go:1742 的 root 关闭加 guard：若 journal 已 rowFinished 则跳过（防双写），否则照旧（fallback：journal 未开行时）；③loop 传的 state/summary/tokens/cost 来自 loop 自己的 finishTask，与 res 应一致（以测试核）。
- 牙：child task.spawn 成功后子行 state=done；反向：摘 FinishTaskLog 写 ⇒ 子行 running、crash_recovery 语境红。
- 风险：动 root 关闭路径，须回归现有 root 收口用例；cmd/wisp 共享面单腿。
- 状态：本轮仅裁形+勘查落档；实现留 fresh/后续（编排者上下文近枯，⛔ 半成品动 run.go）。

## Progress log（续）
- [2026-11-04] agent=编排者 did=AC#1 裁形=Option A + 勘查 loop/journal/run.go 关闭三处机制。

## 实现（编排者 2026-11-04，待 v1 非实现者裁）
- AC#1 已落 Option A：run.go taskLifecycleJournal.FinishTaskLog 由 inert 改为调 store.FinishTaskLog 写终态（5s bound，同 StartTaskLog 风格）。root 双写值等价幂等（loop finishTask 与 run.go:1742 同源 cost）。
- AC#2 牙：把 cmd/wisp/tasklog_fkcheck_378_test.go 的 [B] 段由「只记录不断言 state」升级为断言 childRow.State==done。修后 e2e ok 0.285s、Test233 ok 2.600s、build rc0、gofumpt 空。反向：把 store 写回改 inert ⇒ e2e 红 [405] 子任务 state=running want done；还原 build rc0。
- Status 待 v1 裁定后置 done。

## v1 非实现者裁决（2026-11-04 19:3x +08，agent=`405-v1`）：**PASS / 收口**

- **四数（现跑）**：`go build ./cmd/wisp` rc=0；`go vet ./cmd/wisp` rc=0；`go test ./cmd/wisp -run Test378SpawnedChildToolCallsHaveAParentRow -count=1 -timeout 150s` = **ok 0.398s**（PASS 0.24s；`[378/405] 子任务行 … state=done`、`PRAGMA foreign_key_check = 0`、`task_log 共 2 行（root=1 child=1）`）；`go test ./cmd/wisp -run Test233HotTierAgentReachesTheNextTask -count=1 -timeout 120s` = **ok 4.455s**。`gofumpt -l cmd/wisp/run.go cmd/wisp/tasklog_fkcheck_378_test.go` 空（rc=0）。额外 `sh scripts/d22scan.sh` rc=0、`PASS=35 FAIL=0`、clean。
- **反向读数（自做）**：把 `FinishTaskLog` 的 `return j.store.FinishTaskLog(...)` 换成 `_ = …; return nil`（惰性）⇒ 同一 e2e **rc=1**，逐名报 `[405] 子任务 … state="running" want "done"`，`[378/405]` 行同时打印 `state=running`；**只有这一枚断言红**（fkc 仍 0、root 行仍在）⇒ 牙只咬 405 那一跳。`git checkout -- cmd/wisp/run.go` 还原后 `git diff --stat` 空、`go build ./cmd/wisp` rc=0。
- **AC#1 判据**：`store.FinishTaskLog`（`internal/memory/dao_tasklog.go:49-76`）是**按 id 的单条 UPDATE**（ended_at/state/summary/tokens/cost/error_class），非插入、非累加；root 的两次写（loop `internal/agent/loop.go:1057` → run.go `cmd/wisp/run.go:1751`，中间有 `bg.Wait()` 屏障 `run.go:1682`）**值同源**（`finish()` 先把 `cost.Usage`/`cost.Micros` 写进 `res`，两个 state 都过 `agent.TaskLogState`）⇒ 后写的 run.go 是最终值，**终态/数值与 405 前一致**；loop 先落一次终态还对 CLOSE 失败更稳；无聚合二次计（cost_daily 是 loop `bookCostDay` 每任务一次，不在这两处 Finish 内）。⇒ 定形② 的 `rowFinished` guard **未落但不必要**（幂等替代成立），以本条作**偏离登记**。
- **AC#2 判据**：`tasklog_fkcheck_378_test.go` [B] 段现断言 `childRow.State == memory.TaskStateDone`，反向即红并具名 ⇒ 成立（`[A]` 子任务 fs.read/success 与 `[C]` root 正控仍在，牙不是空 0）。
- **378 / root 收口**：`Test378SpawnedChildToolCallsHaveAParentRow`、`Test378ChildLoopWritesOneTaskRowAndNoToolCallRows`、`TestTicket378ChildTaskRowParentsItsToolCalls`、`TestTicket378OrphanToolCallIsWhatTheRulerSees`、`Test43CrashSweepMarksRunningTasksInterruptedOnce`、`Test370TruncatedRunRowsAnInterruptedTask` 全 PASS；root 行存在且终态由 run.go 收口（Test370 的 `库 task_log.state = "done"` 读数）。
- **残余（登记不挡收口；owner=cmd/wisp）**：三处注释仍写 405 前的单写手理由，与现码相反：`run.go:1290-1298`（"Every other method is deliberately inert … the task row is closed there, once … not by forwarding the loop's own finish"）、`run.go:1717-1721`（"written once"）、`cmd/wisp/tasklog_subloop_378_test.go:127-128`（"StartTaskLog writes, everything else is inert"）。建议随后续 cmd/wisp 小笔更正（裁决者≠实现者，本票不动生产码）。
- **既有红（非本票）**：`Test370TheInterruptedRowNeedsTheCeiling` 在本环境红（`mockllm captured no request on route "chat"`）；用 405 前形状（`FinishTaskLog` 惰性）复跑**同样红** ⇒ 与 405 无关，另账。

