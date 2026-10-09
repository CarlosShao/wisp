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

用的那一发（具名写清）：harness 的 PATH 逐字按票面/派单给的形状，**只多插了一个 `-v`**，
因为 AC#3 要的是逐用例三数（PASS/FAIL/SKIP），非 `-v` 时 `ok` 包里不列通过项、拿不到 PASS 与 SKIP 的名册：

```
cd "D:\work\workspace\projects plans\Wisp" && \
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -v ./internal/panel/ ./cmd/wisp/ -count=1
```

- 命令：改前那一刻跑，工作树里这 4 枚件＝HEAD 内容（§2 已证 `git diff --stat HEAD` 空），
  别人在飞的改动（`design/**` 若干删除、`.gitignore` 修改、别人的 probes 件）一律未动、未还原。
- 原始输出全量落盘＝本目录 `raw-pre.md`（2427 行）。
- **rc=1**（本来就不绿，见下）。
- 环境红排查：`grep -c "0xc0000135" raw-pre.md` = **0** ⇒ 用例真的跑了（`cmd/wisp` 耗时 460.508s 即证）。

三数（并排）＋包级读数：

| 包 | PASS | FAIL | SKIP | 包级行 |
|---|---|---|---|---|
| `internal/panel` | 124 | 6 | 0 | `FAIL github.com/CarlosShao/wisp/internal/panel 7.019s` |
| `cmd/wisp` | 265 | 6 | 2 | `FAIL github.com/CarlosShao/wisp/cmd/wisp 460.508s` |
| 合计（顶层 `^--- X:`） | **389** | **12** | **2** | rc=1 |
| （含子测试的 `PASS`，另尺） | 573 | 12 | 2 | 尺＝`grep -cE "^[[:space:]]*--- X:"` |

改前红名册（逐名，顶层 12 枚）：

`internal/panel`（6 枚）：
1. `TestApprovalCardViewJSONKeysMatchFrontendTypes`
2. `TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`
3. `TestStreamLogDroppedNamingLedgerIsNotBounded35r8`
4. `TestComposerContractTypesMatchFrontend`
5. `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`
6. `TestC21DesignTokensFourWayAgree`

`cmd/wisp`（6 枚）：
7. `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow`
8. `TestPanelHostRealWindowHopAndLifecycle`
9. `TestAC4FocusReturnToPriorWindowGap33r5`
10. `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`
11. `TestAC14AwaitedBindingReplyReachesThePage`
12. `TestAC14GoSideEvalPushReachesThePage`

★ 与票面现量表的**一处不符，具名报回（不采纳票面、按现跑数走）**：
票 289 `AC#3` 写「`internal/panel` 现量 **5 枚具名红**」并列了 5 个名字；现跑是 **6 枚**——
票面漏了 `TestApprovalCardViewJSONKeysMatchFrontendTypes`，它与票面已点名的
`TestComposerContractTypesMatchFrontend` 同属"前端契约"那一族（同一个被删的 `design/assets/` 与前端类型文件的因）。
⇒ 本腿把改前基线记成 6 枚，并把这一处当作**票面文字的缺陷**上报，不改票面一字。
另：`cmd/wisp` 那 6 枚（真窗口／WebView2 那一族）票面**完全没有列**，本腿照样逐名入册，
因为作差的尺是"名册逐名"，不是"票面列了几枚"。
⛔ 上面 12 枚一枚不当"新增红"处理，也一枚不还原/不判绿；尺只有"改后名册 ⊖ 改前名册＝∅"。

## 5. 改后那一发（指针）

改后三数、逐名作差、那一枚红转绿的具名归因、以及三处"派单/票面与代码原文不符"的上报，
全部在 `30-test-comparison.md`（同一把尺、同一支 harness，只多 `-v`）。
本件 §4 的改前读数一个字没回填改动。


