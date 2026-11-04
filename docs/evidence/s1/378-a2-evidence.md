# 378 AC#2 落地证据（实现者＝编排者留，供 378-v1 非实现者裁决）

## 实现
- 新增 `cmd/wisp/child_task_log_378_test.go`（gofumpt 净，package main）：
  - **P1** `TestTicket378ChildTaskRowParentsItsToolCalls`：用 `taskLifecycleJournal`（`run.go` 装给子环路的**同一类型**）对**真 `memory.Store`** 调 `StartTaskLog` 落子行 → `store.InsertToolCall` 写入 bridge 那一枚 tool_call → 用**独立新连接**跑 `PRAGMA foreign_key_check` 断言 **=0**（票373 同款尺）。
  - **P2** `TestTicket378OrphanToolCallIsWhatTheRulerSees`：无父行的 ghost tool_call ⇒ fkc **=1**（"0 有意义"的反向/正控）。
  - 三读数：基线绿 `ok 0.339s / 0.353s / 0.458s`（-run TestTicket378）。
- 生产码**未改**：修复早已在 `run.go:1573-1585`〔乙〕（`rt.loopOpt = loopOpts` 携 `taskLifecycleJournal`）。

## 门禁（编排者跑尺，可复核）
- `go build ./cmd/wisp` rc=0 · `go vet ./cmd/wisp` rc=0 · `go test ./cmd/wisp -run TestTicket378 -count=1` rc=0(ok) · `gofumpt -l` 空 · `d22scan` clean rc=0。
- `go build ./...` 既有假红（`build/mut_rvw.go` untracked staging 件，非本票）——不在本票射程。

## ★范围诚实声明（AC#2 的边界）
- 本测是 **run.go 占位缝**（`taskLifecycleJournal→Store`）＋独立 `foreign_key_check` 尺的**组件级集成测**；**不驱动真 SSE `task.spawn` 全跑**——因该 e2e 在本机不稳（腿遗留测试 `golden requests=[]` 恒红，非 fkc 断言）。
- 接线（子环路真拿到 journal）以**代码核**：`run.go:1008` `BaseOptions` 返 `rt.loopOpt`；`:1573-1585` 附 `Journal=taskLifecycleJournal`；`subagent_197.go:306` `opt := base`（字面拷贝 base＝rt.loopOpt）；`taskLifecycleJournal.StartTaskLog→store.StartTaskLog`。
- ★本测**不是"只测多一行 task_log"**——它断言 `foreign_key_check=0`，符合 AC#2 验收判据的后半句。

## 待 378-v1（非实现者）
复核上述证据＋自跑测试 → 翻 AC#0–4 → `git mv` 票为 `-done`（commit 旧+新名两 pathspec）→ 编排者 push。
