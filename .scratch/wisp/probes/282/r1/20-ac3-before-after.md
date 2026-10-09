# 票 282 `AC#3` 深度判定与父子停机 — 改读访问器前后的同形状对拉（腿 282-r1，2026-10-09）

票面判据句（`sed -n '18p'` 逐字，未改一字；本腿不翻框）：

```
- [ ] **AC#3 深度判定与父子停机**：读 `subagent_197.go` 深度判定那一跳，给"改读访问器前后"的各一发读数（同一形状），判"静默失效"是真被修掉还是只是换了写法。
```

命令形状（四发只有盘上那一行不同）：`go test ./internal/tools/ ./internal/agent/ -count=1 -v`（⛔ 带 `-count=1`，无缓存重放）。

## 0. 那一跳的原文（现读，`internal/tools/subagent_197.go`）

```
262		parentID := TaskID(ctx)
263		if parentID == "" {
264			return Result{Text: "这条调用没有宿主给的任务 id（fail-closed：不知道父任务是谁，" +
267		if rec, ok := t.d.Roster.Look(parentID); ok && rec.Kind == TaskKindSubagent {
270					"子代理不许再派子代理。", MaxSubagentDepth, parentID, rec.Kind), IsError: true}, nil
```

（行号由 `grep -n` 现量复量，非引用票面。）

⇒ 深度判定吃的是 `Roster.Look(parentID)` 的 `ok`：`parentID` 读错 ⇒ 查不到行 ⇒ `ok=false` ⇒ **整支静默跳过**（这就是"静默失效"的形状）。
`internal/tools/task.go:708` 那一跳同形（`caller` 读错 ⇒ `:732 if rec.ParentTaskID != caller` ⇒ 走到"拒绝停止"那一支，父停不了子）。

## 1. 对拉一 · `subagent_197.go:262`（spawn 侧，深度判定所在的那一跳）

| 形状 | 盘上那一行 | 深度那枚用例（逐字） | 名册父子那枚断言（逐字） | 两包三数 |
|---|---|---|---|---|
| **后**（现产码＝读访问器） | `	parentID := TaskID(ctx)` | `--- PASS: Test197ChildCannotDeriveSubagent (0.00s)` | `--- PASS: Test283IdentityChainThroughTheRealLoop (0.06s)`（红句缺席） | tools PASS=209/FAIL=0/SKIP=0（`ok ... 13.947s`／`13.803s` 两跑一致量级）· agent PASS=84/FAIL=1/SKIP=0 |
| **前**（种回 corr） | `	parentID := CorrelationID(ctx)` | `--- PASS: Test197ChildCannotDeriveSubagent (0.00s)` | `--- FAIL: Test283IdentityChainThroughTheRealLoop (0.06s)`<br>`    internal/tools/ticket283_corr_identity_rulers_test.go:194: 子代理行的 ParentTaskID = "caf08939-ea19-4564-bbc5-deb69566e3b2#call-283-spawn", want 宿主 task id "caf08939-ea19-4564-bbc5-deb69566e3b2"（task.spawn 读的是 TaskID(ctx)，不是 per-call corr）` | tools PASS=208/FAIL=1/SKIP=0（`FAIL ... 13.798s`）· agent PASS=84/FAIL=1/SKIP=0（那枚既有红同下） |

两跑共有的既有红（不归本腿、归票 286，未修未 Skip）：`--- FAIL: TestGoldenSingleToolCall (0.00s)`。

读数比对：**深度那枚用例两发同绿**（PASS/PASS，逐字相同）；**同一条链上的身份登记一发绿一发红**。

## 2. 对拉二 · `task.go:708`（父停子侧）

| 形状 | 盘上那一行 | 停机那组用例（逐字） | 真链那枚断言（逐字） | 两包三数 |
|---|---|---|---|---|
| **后**（现产码＝读访问器） | `	caller := TaskID(ctx)` | `--- PASS: Test221ParentStopsItsOwnChildRowAndStreamSettle (0.06s)`<br>`--- PASS: Test221SubagentCannotStopSiblingOrItself (0.00s)`<br>`--- PASS: Test236R2TaskCancelRefusesWhenHostGaveNoCallerID (0.00s)` | `--- PASS: Test283IdentityChainThroughTheRealLoop (0.06s)` | tools PASS=209/FAIL=0/SKIP=0 · agent PASS=84/FAIL=1/SKIP=0 |
| **前**（种回 corr） | `	caller := CorrelationID(ctx)` | 同三枚**仍逐字 PASS**（直接派发夹具 `corr == taskID`，分不出） | `--- FAIL: Test283IdentityChainThroughTheRealLoop (0.07s)`<br>`    internal/tools/ticket283_corr_identity_rulers_test.go:206: task.cancel 的回执没有走到权威判定通过那一支: "拒绝停止：f4f7abb1-9a1c-456e-a959-5fa17112b75f 是任务 cbe4af20-760f-464e-9ba7-cc48f6bd25d3 派生的孩子，不是调用者 cbe4af20-760f-464e-9ba7-cc48f6bd25d3#call-283-cancel 的孩子——只有父任务能停自己的孩子，兄弟之间、旁支之间都停不了（票 221 AC#3 的正控）。"`<br>`    internal/tools/ticket283_corr_identity_rulers_test.go:209: 父任务停自己的孩子被身份判定拒了（调用者读到的不是 task id）: （同上一句文案，逐字重复）` | tools PASS=208/FAIL=1/SKIP=0（`FAIL ... 14.059s`）· agent PASS=84/FAIL=1/SKIP=0 |

⇒ **父停子那一支的"静默失效"是真的**：旧形状下父任务停自己的孩子被身份判定硬拒（红句里那句"不是调用者 …#call-283-cancel 的孩子"就是失效现场），
新形状下同一断言绿；两发只有盘上那一行不同 ⇒ 不是换写法。

## 3. 方法学前置（票面逐字那句）—— 本腿现量，不引注释

「夹具必须 `TaskID != CorrelationID`，否则两发必同绿不算读数」：

- 成立的一面：两发红句里都同屏出现两个不相等的值 ——
  `...bd25d3#call-283-spawn` / `...bd25d3#call-283-cancel`（corr）对 `caf08939-ea19-…`／`cbe4af20-…`（task id）；
  且票 283 尺一在干净树 `--- PASS: Test283TaskIDAccessorIsTheTaskIDNotTheCorr (0.00s)`（handle 为 `corr:"corr-283"` 对 `taskID:"task-283"`）。
  尺件 `039ec93c`（227 行新建零删除）／`9a00a890` 已在 `HEAD` 树内（`git show --stat` 逐枚复量）⇒ 前置今天具备，对拉可用。
- 不成立的一面（**这就是深度那一支判不动的原因**）：`Test197ChildCannotDeriveSubagent` 走的不是真 loop 派发。
  它调的是 `childDir.Execute(... agent.ToolRequest{TaskID: "child-197", Name: "task.spawn", ...})`，
  而那枚 `subagentToolProvider`（`internal/tools/subagent_197.go:573-580`）在 `Execute` 入口就对 `task.spawn` 硬拒：
  ```
  576			Text: fmt.Sprintf("拒绝执行：子代理不能再派生子代理（深度上限 %d），"+
  577				"task.spawn 也不在你的工具目录里。", MaxSubagentDepth)
  ```
  ⇒ 根本到不了 `:262`／`:267`；它的 `TaskID`/`CorrelationID` 前置对不对都无意义。

## 4. 深度判定那一支：今天有没有仪器（现量三把）

1. `grep -rn "子代理不许再派子代理" --include=*_test.go internal | wc -l` = **`0`**
   ⇒ `:267-271` 那句运行期拒句**没有任何测试引用**（生产文件里唯一命中＝`internal/tools/subagent_197.go:270`）。
2. `grep -rn "TaskKindSubagent" --include=*_test.go internal` = 5 枚命中，全部是**断言孩子那一行的 `Kind`**
   （`subagent_197_test.go:305-306`／`task_cancel_221_legs_test.go:333`／`ticket283_corr_identity_rulers_test.go:197-198`），
   **没有一枚**把"调用者自己是子代理"当夹具前提去派孙。
3. 种 A（`:262` 读 corr，正是让 `Look(parentID)` 必 miss 的形状）之下，深度那枚用例逐字 `--- PASS: Test197ChildCannotDeriveSubagent (0.00s)`，
   且整跑日志里 `名册里出现第二代` 那句**一次都没出现**（`grep -n "名册里出现第二代\|深度上限" /tmp/w282r1-seedA.txt` = 空）
   ⇒ 深度那一支在四种形状下都不会红 —— 它今天**没有仪器**。

## 5. 本格判语（三句，各管各的范围）

- **身份那一跳（spawn 侧登记／cancel 侧停机）＝真被修掉**，⛔ 不是换了写法：
  证据＝`ticket283_corr_identity_rulers_test.go:194`（前红后绿）与 `:206`/`:209`（前红后绿），两发各只动盘上一行、还原逐字等值。
- **深度判定那一支（`:267-271`）＝判不动**，且这一格**缺的不是种刀而是夹具**：
  缺一枚"调用者那一行在名册里是 `TaskKindSubagent`、派发经真 loop（故 `corr != taskID`）、请求 `task.spawn`"的夹具，
  并且它必须能绕过 `subagentToolProvider` 那道结构性硬拒才够得着 `:267`。
  现状：`MaxSubagentDepth = 1` 由结构（孩子的目录里没有 `task.spawn`）保证，运行期那一跳是 defense-in-depth，**无人看**。
  ⇒ 票面 §AC#3 要求的"第 3／第 4 发前后对读数"里，**深度那一支撑有两发同绿、不构成读数**；本件如实具名，不硬凑红。
- 顺带具名一处与本腿无关的形状：种 A 同发下 `:206` 也红，但红句文案是 `"本用例不带 L2 队列"`（不是"拒绝停止"那一支）
  ⇒ 旧形状下 spawn 身份读错会把后续那枚 `task.cancel` 顶到别的回执上去；本腿不据此下任何额外结论，只记下读数。

## 6. 门禁与还原自检（本件两发共用 `10-ac1-teeth.md` §4 的同一次跑）

- 两枚文件还原后：`git hash-object` 分别 = `8c9266a7dfaf62fbf9418cf8055a5905eb98b8b1`／`3ec8d498f1613f0b7eeadb0b840d4c77cc04ed84`，与各自种前**逐字等值**。
- `git status --porcelain -- internal cmd` = **空**（每发收口时各量一次）。
- `gofmt -l`／`gofumpt -l`（只读档）两枚文件 = 零输出；`GOFLAGS= go build ./...` `rc=0`；`sh scripts/d22scan.sh` `rc=0`。
