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
