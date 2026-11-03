# 票 31 / 31-r2 内存层 —— 落地件（schema v6 reminders 表 + DAO + 迁移幂等）

- **提交**：`5ebdf6961`（票 31 内存层）
- **范围**：`internal/memory`
- **内容**：
  - schema v6：新增 `reminders` 表（`schema_reminders_v6.go` 的 `buildDDLV6`/`applyV6`），`IF NOT EXISTS` 幂等；`status` 的 CHECK 由 `reminderStatusRoster` 生成（ticket 196「Go 花名册是权威、CHECK 是其一嘴」规则首次用于新枚举）。
  - DAO（`dao_reminders.go`）：`InsertReminder` / `ReminderByID` / `ListReminders` / `MarkReminderFired`（armed→fired 单向）/ `CancelReminder` / `CountMissedArmed`。
  - 版本钉对齐：`schema.go`、`models.go` 及若干 `_test.go` 由 v5→v6；`SchemaVersionTarget = 6`。
- **验证**：`go test ./internal/memory -count=1` → `ok ... 30.35s rc=0`（编排者修两处自写测试 bug 后：recurrence 守卫先剥 `--` 注释再判；幂等测试 `Open` 后补 `t.Cleanup(Close)` 以免 Windows TempDir 清理 file-in-use）。
- **硬排除**（票 31）：无 recurrence/cron 列、无 action/payload 列；boot 只 `CountMissedArmed` 报计数、不逐条补发。
- **spec-note 债**：`SPEC-02 §3` 表册 roster 扩到 v6，待机主按 `D35` 增列（沿用 `A1445`（44b v5）的处理方式）。
