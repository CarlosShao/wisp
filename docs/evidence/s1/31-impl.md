# 票 31 / 31-r3 tools+agent 层 —— 落地件（reminder 三工具 + 进程内单调时钟计时器）

- **提交**：`da9b5fb9c`（票 31 tools+agent 层）
- **范围**：`internal/tools`、`internal/agent`
- **内容**：
  - `internal/tools/reminder.go`：`reminder.set`（L1）/ `reminder.list`（L0）/ `reminder.cancel`（L1）。set 确认文案含 caveat「进程没在跑就不会触发」（16.5.4 诚实规则）；recurring/cron/action-on-fire/OS 计划任务/重启补发 全部显式 REJECTED（代码注释 + 注册 note，非 DEFERRED）。
  - `internal/agent/reminder_timer.go`：进程内 named+owned 计时器；`SystemReminderClock.MonoNow/WallNow` 单调计时，按墙钟目标算 delta（时钟跳变不影响）；到点 `Notifier.Notify` + `MarkReminderFired` + 记 `reminder.fire` 的 `tool_call` 审计行；boot 走 `ReportMissed -> CountMissedArmed`，只报错过条数、不补发。
- **验证**：
  - `go build ./internal/tools ./internal/agent ./cmd/wisp` → `rc=0`
  - `go test ./internal/tools -skip Search -count=1` → `ok ... 32.4s`
  - `go test ./internal/agent -count=1` → `ok ... 6.56s`
  - 窄测 reminder：tools `ok 0.25s`、agent `ok 0.115s`
  - `scripts/d22scan.sh`（D22 门禁，含自测正向对照 + 真仓扫描）→ 修 emoji/phantom 后通过
- **过程说明**：实现腿 `31-r3`(qwen) 产出主体代码后在跑测试时卡死被停；编排者修复测试辅助签名（`itoa(int64)`）并清理 Go 源码字符串内 emoji（`⚠`→`注意`）后全层转绿。
