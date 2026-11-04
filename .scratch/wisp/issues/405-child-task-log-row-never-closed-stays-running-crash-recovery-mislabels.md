# 405 — 子任务 `task_log` 行开了却从不关：child 恒 `running`，crash_recovery 会把成功的也标 `interrupted`

- Status: **open**
- 立案：2026-11-04 源起票 **378-v1 R3**（378 的修复让子任务获得 `task_log` 行，但 `FinishTaskLog` 惰性 ⇒ 子行永不关）
- 关联：票 **378**（子任务获父行、fkc=0，已解决"无父"）· 票 **373**（级联）· 票 **404**（R4 前置 bug，已修）

## 这是什么（人话）

`cmd/wisp/run.go:1324-1328` 的 `taskLifecycleJournal.FinishTaskLog` 是**有意惰性**（`return nil`）。子环路用的是同一个 journal ⇒ 子任务 `StartTaskLog` 落一行 `running`、但 `FinishTaskLog` 不写 ⇒ **子任务行永远停在 `running`**。

全仓只有 `run.go` 在 `execute()` 收尾为 **root** 任务关闭行（`res.TaskID`）；**子任务没有任何关闭点**。而 `cmd/wisp/crash_recovery.go:145-165` 会把仍是 `running` 的行改判 `interrupted` ⇒ **成功完成的子任务被误标"中断"**。

## 影响

用户可见的隐私页"任务日志"里：**每个子任务恒显 `running` 或被误标 `interrupted`**，与实际成败不符。

## 结案判据

- [ ] **AC#1 裁形**（编排者拍）：子任务终态写手（主流 harness：子任务完成应有终态）。候选：①`FinishTaskLog` 对子任务写入终态（task 行半，不动 tool_call 惰性）；②`run.go` 侧在子任务收尾处关闭；③按 `SPEC-12 §5` 五字段登记推迟。
- [ ] **AC#2 牙**：一发 `task.spawn` 成功后该子行 `state=done`（或等价终态）；反向＝摘终态写手 ⇒ 用例转红逐名报。
- [ ] **AC#3 门禁** / **AC#4 纪律**（同上票规范）。

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
