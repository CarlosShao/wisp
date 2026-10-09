# 票 289 / 腿 289-v1 — AC#1 判语：四句新注释逐字读＋每一枚锚现跑复量＋全仓名册尺

**判语：成立。**（四句逐字读下来每一句都是实话，锚全部现跑命中；另现读到**票面"四句"低估了这一族**——见 §4，那一处不属本票 AC#1 写面，本腿不据此判失败，具名上报。）

尺的起点＝本腿自己跑的 `git show -U0 --format="" 73c30dfe`（全量落盘＝本目录 `raw-v1-diff-u0.md`），⛔ 不引用 r1 的输出。
四个 hunk 头逐字：`@@ -20,4 +20,4 @@`／`@@ -155,2 +155,2 @@`／`@@ -56,3 +56,3 @@`／`@@ -466,5 +466,5 @@`
⇒ 写面＝**4 枚件、4 个 hunk**，没有第五处，也没有"顺手改这 4 枚件里的其它注释"（四个 hunk 全在 `corr==taskID` 那一段里）。

## 1. 本腿现跑的共同事实底座（每枚锚都自己量过）

| 锚 | 现跑命令（逐字） | 现跑输出（逐字） |
|---|---|---|
| `loop.go:603` 定义 | `grep -n "func callCorr" internal/agent/loop.go` | `603:func callCorr(taskID, callID string, index int) string {` |
| `loop.go:676` 使用 | `grep -n "CorrelationID: callCorr" internal/agent/loop.go` | `676:			TaskID: taskID, CorrelationID: callCorr(taskID, p.call.ID, i),` |
| corr 永不等 taskID | `sed -n '603,608p'`（函数体） | `	if callID == "" {` / `		return fmt.Sprintf("%s#call-%d", taskID, index)` / `	}` / `	return taskID + "#" + callID` / `}` ⇒ 两支返回都带 `#` 后缀 |
| 造这一跳的提交 | `git log --oneline -1 bd124b2a`（rc=0） | `bd124b2a 242-corrland-1 落地：loop 每次工具调用铸独立 corr（callCorr=taskID#callID，任务级溯源仍走 taskID）…` |
| `internal/tools/subagent_197.go:262` | `sed -n '259,263p'` | `259:	// The task-level id of the caller: since ticket 242 the correlation id is` … `262:	parentID := TaskID(ctx)` |
| 同一个 parentID 落进名册行 | `grep -n "ParentTaskID" internal/tools/subagent_197.go` | `429:		ParentTaskID: parentID, Label: label, Kind: TaskKindSubagent,` |
| `cmd/wisp/subagent_carrier_197_test.go:173` | `sed -n '171,174p'` | `172:		out, eerr := rt.bridge.Execute(ctx, agent.ToolRequest{` / `173:			TaskID: rootID, CorrelationID: rootID, CallID: fmt.Sprintf("spawn-197-%02d", i),` |
| join 的两半建表 | `sed -n '212,221p' internal/panel/subagent_roster_197.go` | `212:	waiting := make(map[string]bool, len(pending)+2*len(l1Windows))` / `214:		indexWaitingKey(waiting, card.CorrelationID)` / `219:		indexWaitingKey(waiting, w.CorrelationID)` / `220:		indexWaitingKey(waiting, w.TaskID)` |
| 行侧查询 | `sed -n '236p' internal/panel/subagent_roster_197.go` | `236:			BlockedOnApproval: row.TaskID != "" && waiting[row.TaskID],` |

## 2. 四句逐字判真伪（`+` 行＝`git show 73c30dfe` 原样）

### 件 1 `internal/panel/subagent_roster_197.go:56-58`（产码件，`ParentTaskID` 字段的注释）

```
	// The entity leg files the deriving call's task id (TaskID(ctx),
	// internal/tools/subagent_197.go:262), never its correlation id: since
	// bd124b2a loop.go's callCorr mints one per call (<taskID>#<callID>, :603/:676).
```

逐枚核：`TaskID(ctx)` 在 `:262` ✓；同一 `parentID` 在 `:429` 写成 `ParentTaskID` ✓ ⇒ "entity leg files the deriving call's **task id**" 为真；"never its correlation id" 为真（`:263-266` 还在 `parentID==""` 时 fail-closed，没有任何 corr 回落）；`callCorr` 形状 `<taskID>#<callID>` 由函数体两支直接坐实 ✓；`:603/:676` 两枚锚现跑命中 ✓。**判：实话。**

### 件 2 `internal/panel/subagent_roster_197_test.go:466-470`

```
//  3. The blockedOnApproval join is asserted against a card this file made. On
//     the production path it pairs the row's task id, not the corr: since
//     bd124b2a callCorr mints one corr per call (loop.go:603, used :676),
//     shape <taskID>#<callID>, never == taskID. A host that paired on the
//     corr alone would miss every row; only cmd/wisp's real gate catches it.
```

"pairs the row's task id, not the corr" ✓（`:236` 的查询键就是 `row.TaskID`）；"a host that paired on the corr alone would miss every row" ✓（真实 corr 带 `#` 后缀，`:236` 拿 `row.TaskID` 去查一张只按 corr 建键的表必然全空）；两枚 loop 锚 ✓。**判：实话。**

### 件 3 `internal/panel/subagent_blocked_220_test.go:155-156`

```
	// Since bd124b2a the loop's corr is per call (<taskID>#<callID>, loop.go:603,
	// used :676), never == taskID: the corr-only pairing below is this fixture's.
```

"the corr-only pairing below is this fixture's" ✓：紧接其下的 `:157` 现读逐字
`	data = blockedPacket(t, nil, []L1WindowWait{{CorrelationID: "child-1", TaskID: ""}})`
——只有 corr 一枚、`TaskID` 为空 ⇒ 那是**夹具自己的形状** ✓。同函数 `:151` 的
`{CorrelationID: "corr-host-9", TaskID: "child-1"}` + `:152 assertBlocked(t, data, "child-1", true)` 现读还在，
⇒ "永不等 taskID" 与"这条夹具是 corr-only"两句同时坐实。**判：实话。**

### 件 4 `cmd/wisp/subagent_carrier_197_test.go:20-23`

```
// bridge with the fixture's OWN ids (TaskID == CorrelationID at :173), which is
// NOT the loop's: since bd124b2a callCorr mints one corr per call
// (<taskID>#<callID>, internal/agent/loop.go:603, used :676). The spawn waits
// for the ROOT's own roster row - proof it is live - then uses rt.bridge.
```

`TaskID == CorrelationID at :173` ✓（`:173` 现读把 `rootID` 同时塞给两个字段）；"NOT the loop's" ✓；末两句的叙述也在盘上：`:153 rootID, err := waitRootRow197(rt)` → `:172 out, eerr := rt.bridge.Execute(...)` ⇒ "先等根的行、再走 bridge" 为真 ✓。**判：实话。**

⚠ 一处**证据件级**（非产码注释级）的小错要具名：r1 的 `00-anchor-and-baseline.md:48` 与 `10-four-comments.md:56` 把行侧查询写成 `subagent_roster_197.go:235`，本腿现读是 **`:236`**（`:235` 逐字 `			StreamKey:         row.StreamKey,`）。该行锚**没有**被写进四句注释里，所以不影响 AC#1 的判定，只是 r1 证据件的引用偏一枚。

## 3. 名册尺一（与 r1 同尺，两枚拼法）——第五枚在册件的**落点**

命令逐字：`grep -rn --include='*.go' -e 'CorrelationID == TaskID' -e 'TaskID == CorrelationID' .`
原样输出（去掉 probes 快照）：

```
./cmd/wisp/subagent_carrier_197_test.go:20:// bridge with the fixture's OWN ids (TaskID == CorrelationID at :173), which is
./internal/tools/ticket283_corr_identity_rulers_test.go:7:// 一律 TaskID == CorrelationID（bridge.go 的 orDefault 回落也让缺 corr 的
./.scratch/wisp/probes/197/r3b/pre/subagent_roster_197.go:52 / :173（别人票 197 的改前存档，不在册）
```

⇒ 第五枚在册件＝**`internal/tools/ticket283_corr_identity_rulers_test.go:7`**（r1 报的文件名对、落点对）；
**派单里写的 `internal/panel/ticket283_corr_identity_rulers_test.go:7` 目录错了**（`internal/panel/` 下没有这枚文件）。
r1 的现读结论"那句不是假话"——**本腿独立核过，成立**：那句话的主语是"既有用例的派发 fixture"，
且它说的 `orDefault` 回落机械可核（`internal/tools/bridge.go:269-271` 逐字 `	if req.CorrelationID == "" {` / `		req.CorrelationID = req.TaskID`；
`internal/tools/bridge.go:528` 逐字 `		corr: orDefault(req.CorrelationID, req.TaskID), taskID: req.TaskID, bus: b.cancel,`），
而同文件 `:13-14` 自己就写了真生产链（`真生产链是 agent.Loop 的 callCorr(taskID, callID, i)（loop.go 每枚调用一个 corr，形如 taskID + "#" + callID）`）。

## 4. 名册尺二（本腿自定的宽尺）——⚠ 找到票面没数的**第五句**

尺（多拼法覆盖：大小写、全角等号、`corr == task`、中文"相等/等于/就是"、散文 `correlation id equals … task`；`--include='*.go'` 全仓＋`git grep -I` 追一遍文档，`.scratch` 探针桶单独排除）：

```
grep -rniE --include='*.go' -e 'correlation[[:space:]_]*id?[[:space:]]*(==|＝|!=|is|equals)' \
  -e 'corr[[:space:]]*(==|＝)[[:space:]]*task' -e 'task[[:space:]]*(==|＝)[[:space:]]*corr' \
  -e 'corr( ?id)? ?(等于|就是)' -e 'correlation id (is|as|same as) (the )?task' . | grep -v '^./.scratch/'
```

逐处判（只列真属"corr＝taskID 这一族"的；其余命中是 `CorrelationID == ""`／`!= "pump-corr"` 这类**空值/字面量比较**，不是等值叙述）：

1. `internal/tools/ticket283_corr_identity_rulers_test.go:7` 与 `:13`——已判"不是假话"（§3）。
2. **`cmd/wisp/subagent_carrier_197_test.go:598-600`（在册，⚠ 今天仍是旧说法）**：
   ```
   //  1. blockedOnApproval is never asserted on a production packet. The join needs a card
   //     whose correlation id equals a CHILD task id, and nothing in this tree can make
   //     that happen from the outside: the agent loop never sends tool_choice, so mockllm
   ```
   "The join needs a card whose correlation id equals a CHILD task id"**已过期**：`subagent_roster_197.go:219-220` 对每枚 L1 window 两半都建表，`:236` 用 `row.TaskID` 查 ⇒ 现成的反例就在同一棵树里跑着：`internal/panel/subagent_blocked_220_test.go:151-152` 的窗口 `CorrelationID: "corr-host-9"`／`TaskID: "child-1"` 让 `child-1` 那行读 `true`——**corr 不等于 child 的 task id，行照样亮**。
   ⚠ 它与票 197 的关系：`cmd/wisp/subagent_blocked_197_test.go:7-10` 已经把同一句话当作**被更正的旧说法**引用（逐字 `197-r3's file closed its list with "the join needs a card whose correlation id equals a CHILD task id…" 197-r3b corrected that in its own §⑤1 - the sentence was too strong`）⇒ 那一处**不是假话**（是引文＋当场推翻）；而 `carrier_197_test.go:598-600` 这一处是**自己还在主张**。
   ⇒ 本格处置：这句话落在票 289 写面第 4 枚件里，但 AC#1 明文"⛔ 不许顺手改这 4 枚件里的**其它**注释"，r1 不动它**是守规、不是漏改**；本腿把它作为**票面枚数低估**上报——票说"四句"，这一族今天要算**五句**（其中一句在四枚件内、一句是引文性质）。**不改票面一字。**
3. `internal/panel/subagent_roster_197.go:122` `// section's row with correlationId == streamKey IS the text to render.`——`corr==streamKey`，**另一族**（流键配对），不在本票射程，也非假话（票 197 的 `streamKey = subagent:<taskId>` 那一族）。
4. 文档侧（`docs/reports/*`、`docs/evidence/*`）命中全是**历史记录**（如 `pending-and-issues.md:2208` 明写"今天 corr==taskID 只是 `loop.go:644` 的巧合，不是设计"），⛔ 不是产码假话，本票也不动 `docs/**`。
