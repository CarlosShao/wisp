# 票 289 / 腿 289-r1 — 锚点、现量复量、改前基线

## 1. 锚

- 工作树根：`D:\work\workspace\projects plans\Wisp`，分支 `dev`。
- 起手记录到的 HEAD：`34a01259`（`git log --oneline -3` 当时顶部）。
- 起手过程中共享工作树的 HEAD 被别的腿推进了一枚，复量时刻 HEAD＝`fd692f7e`
  （`票 290 只读腿 290-a1 骨架件入库`，`git show --stat` 名册只含 `.scratch/wisp/probes/290/a1/00-findings.md`，
  与本腿的 4 枚件无交集）。**下面所有"现量"都是对 `HEAD`=`fd692f7e` 现跑的。**
- 时钟（不手打）：`date -u -d "+8 hours" "+%Y-%m-%d %H:%M +08"` → `2026-10-09 14:06 +08`。

## 2. 四枚件现量复量（`git show HEAD:<path> | sed -n '<N>p'`，逐字）

命令逐字（同一发循环里跑的四枚）：

```
for spec in "internal/panel/subagent_roster_197.go:57" \
            "internal/panel/subagent_roster_197_test.go:468" \
            "internal/panel/subagent_blocked_220_test.go:155" \
            "cmd/wisp/subagent_carrier_197_test.go:20"; do
  p="${spec%:*}"; n="${spec##*:}"; git show "HEAD:$p" | sed -n "${n}p"
done
```

| # | `HEAD:file:行` | 现跑输出逐字 | 与票面现量表一致？ |
|---|---|---|---|
| 1 | `internal/panel/subagent_roster_197.go:57` | `	// loop dispatches with CorrelationID == TaskID (internal/agent/loop.go:647),` | 一致 |
| 2 | `internal/panel/subagent_roster_197_test.go:468` | `//     CorrelationID == TaskID (internal/agent/loop.go:647); a host that used a` | 一致 |
| 3 | `internal/panel/subagent_blocked_220_test.go:155` | `	// The production pairing (the loop dispatches with CorrelationID == TaskID,` | 一致 |
| 4 | `cmd/wisp/subagent_carrier_197_test.go:20` | `// bridge with the ids the agent loop itself uses (TaskID == CorrelationID,` | 一致 |

- 另证：`git diff --stat HEAD -- <这 4 枚件>` 空 ⇒ 工作树里这 4 枚件＝HEAD 内容，改前无别人在飞。

## 3. 正确事实的形状（现读核对，注释按这一族写）

- `git show HEAD:internal/agent/loop.go | grep -n "func callCorr\|CorrelationID: callCorr"`：
  - `603:func callCorr(taskID, callID string, index int) string {`（定义）
  - `676:			TaskID: taskID, CorrelationID: callCorr(taskID, p.call.ID, i),`（使用）
  ⇒ 票面引用的 `:647` 今天已在 `:676`；`loop.go:594` 那句 `// callCorr mints the correlation id of ONE tool call (ticket 242). C18 defines` 与 `:367` 都已是实话（禁区，不动）。
- `callCorr` 函数体（`:603-609` 现读）：`callID == ""` 时回 `fmt.Sprintf("%s#call-%d", taskID, index)`，否则 `taskID + "#" + callID` ⇒ 形如 `<taskID>#<callID>`，**永远不等于 taskID**。
- 铸造那一跳的来路＝`bd124b2a`（`git log --oneline -1 bd124b2a` 现跑：`242-corrland-1 落地：loop 每次工具调用铸独立 corr（callCorr=taskID#callID，任务级溯源仍走 taskID）…连带越格自报：tools 任务级身份改读 TaskID(ctx)（cancel.go 新访问器/bridge 样章/subagent_197:262/task.go:708）`）。
- ★ 关键的一条，写注释前现读确认（否则会把一句假话换成另一句假话）：**`ParentTaskID` 今天仍读作父任务 id，但不是因为 corr==taskID，而是因为 entities 腿改读了 `TaskID(ctx)`**——
  `git show HEAD:internal/tools/subagent_197.go | sed -n '259,262p'` 逐字：
  `	// The task-level id of the caller: since ticket 242 the correlation id is` / `	// minted per call (C18 routes approval replies by it), so it can no longer` / `	// name a task; the roster is keyed by task ids.` / `	parentID := TaskID(ctx)`。
  `internal/tools/task.go:702-706` 同形（`the correlation id went per-call in ticket 242 and no longer names a task`）。
- ★ `blockedOnApproval` 那把 join 现读：`internal/panel/subagent_roster_197.go:216-220` 对每枚 L1 window **两半都建表**
  （`indexWaitingKey(waiting, w.CorrelationID)` 与 `indexWaitingKey(waiting, w.TaskID)`），行侧查询在 `:235`
  `BlockedOnApproval: row.TaskID != "" && waiting[row.TaskID]` ⇒ 生产路径上点亮一行靠的是 window 的 **task 那半**，不是 corr 那半。
- ★ 票面第 4 句描述的夹具现读：`git show HEAD:cmd/wisp/subagent_carrier_197_test.go | sed -n '172,174p'` →
  `out, eerr := rt.bridge.Execute(ctx, agent.ToolRequest{` / `			TaskID: rootID, CorrelationID: rootID, CallID: fmt.Sprintf("spawn-197-%02d", i),`
  ⇒ 夹具把同一个值塞给两个字段，那是**夹具自己的形状**，不是 loop 的形状。

## 4. 改前基线 `go test`（AC#3 的前一半）

待填：见本件 §4 之后的追加（跑完即追加，不改上面任何一行）。
