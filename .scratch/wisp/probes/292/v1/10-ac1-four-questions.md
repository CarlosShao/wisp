# 292-v1 10 — AC#1 四问逐环（我现跑的尺，⛔ 不引票面读数）

被验三行（HEAD 后逐字，尺＝`sed -n '598,604p' cmd/wisp/subagent_carrier_197_test.go`，量于 `bccafc30`）：

```
598: //  1. blockedOnApproval is never asserted on a production packet. The join has two
599: //     producers (subagent_roster_197.go:214, :219-220): a card's corr, never a task id
600: //     for the loop (loop.go:603), and L1 windows no run feeds; no tool_choice is sent, so mockllm
601: //     never answers with a tool call, so a spawned child never reaches the gate. The
602: //     field is therefore asserted at the pump against a real ApprovalCardView
```

⚠ 改动的行＝`:598-600`（三行），`:601-604` 是**未被本改动触及的既有谓语**（`git show -U0` 只有一个 hunk，见 20 件）。

---

## ⓐ join 的来源枚数 —— 【成立】

我的尺＝**整族**（建表函数体全窗）`sed -n '205,240p' internal/panel/subagent_roster_197.go`，另加 `sed -n '198,200p'` 取签名。

- `:199 func taskRosterSectionFrom(state TaskRosterState, pending []ApprovalCardView, l1Windows []L1WindowWait) *TaskRosterSection`
- `:212 waiting := make(map[string]bool, len(pending)+2*len(l1Windows))`
- `:213-215` 卡片那一半：`:214 indexWaitingKey(waiting, card.CorrelationID)`
- `:216-221` L1 窗口那一半：`:219 indexWaitingKey(waiting, w.CorrelationID)`／`:220 indexWaitingKey(waiting, w.TaskID)`
- 查询处 `:236 BlockedOnApproval: row.TaskID != "" && waiting[row.TaskID]`

⇒ 来源**两枚**（两个 `for` 各一枚），腿注释里点的行号 `:214` 与 `:219-220` **逐字对得上**。
⚠ 一处我要记的口径差：`:219` 用的是窗口的 **corr**、`:220` 用的是窗口的 **task id**，所以"枚数＝2"是按**生产者**数的，按 `indexWaitingKey` 调用点是 **3 枚**。注释写 "two producers (:214, :219-220)" 与两种数法都不冲突。

**进攻（未被注释覆盖的一节）**：`indexWaitingKey` 不止存原键，`:267-268 set[base] = true` 会把 `splitClashSuffix(key)` 的前缀也存进去；`splitClashSuffix`（`:277-288`）认"最后一个 `#` 之后**全是 ASCII 数字**"。
⇒ 对 `callCorr` 第二支 `taskID + "#" + callID`，若 provider 的 callID 是**纯数字**，卡片这一半**确实能点亮"该卡片自己的那个 task"的行**。注释那句只说"corr 本身不是 task id"（成立，见 ⓑ），没断言"卡片那一半永不可能点亮行"——**所以这不是假话，但这是下一位程序员最容易读过头的一格**。本仓实测 mockllm 的 callID＝`tools/mockllm/chat.go:244 "id": "call_mock_1"`（非纯数字）⇒ 树内不可复现，属真 provider 的形状余地。

---

## ⓑ "card corr == CHILD task id" 今天是**不可能**还是**只是不可达** —— 【成立：不可能，且注释写法与之一致】

我的尺＝整族 `sed -n '596,608p' internal/agent/loop.go` ＋全仓 callCorr 调用点。

- `:603 func callCorr(taskID, callID string, index int) string`
- `:604-606 if callID == "" { return fmt.Sprintf("%s#call-%d", taskID, index) }`
- `:607 return taskID + "#" + callID`
- 调用点（整族穷举，尺＝`grep -rn "callCorr(" --include=*.go .`，剥掉 `.scratch`）＝**只有两枚**：`internal/agent/loop.go:603`（定义）与 `:676 TaskID: taskID, CorrelationID: callCorr(taskID, p.call.ID, i)`（生产派发）。
- ⇒ 两支返回值都**强制带一个 `#` 且 `#` 后非空**，故 `callCorr(...) != taskID` **恒成立**（结构性质，与外部能不能构造无关）。
- 注释原文 "a card's corr, never a task id **for the loop** (loop.go:603)" ——"never … for the loop" 是**不可能**口径，不是"外面造不出来"口径 ⇒ **ⓑ 这一节腿改对了**，票面 ⓑ 的要求（把"只是不可达"升级成"不可能"）已落进文字。

⚠ 射程限定我必须写清：`internal/agent/approval/gate.go:282 corr := orDefaultText(d.CorrelationID, d.TaskID)`（尺＝`sed -n '268,282p'`）⇒ **不经 loop 的** request（corr 为空时）其窗口/卡片**就以 task id 为键**；`internal/panel/subagent_roster_197.go:180` 那节注释逐字也这么说（"a decision with no correlation id inherits the task id"）。注释的 "for the loop" 一词把这层挡住了 ⇒ **不是假话**；但它**没有**告诉下一位程序员"非 loop 的派发者可以 corr==task id"。票面 ⓑ 那句"结构上不可能"同样只在 loop 射程内成立。**这是票面与注释共有的射程沉默，记账、不算 292 的缺陷。**

---

## ⓒ "nothing in this tree can make that happen" 的射程（仓内 fixture vs 生产装配）—— 【成立】

尺一＝整族穷举 `grep -rn "L1Windows" --include=*.go .`（输出全量，未截断）。逐字分类：

- **定义/读取处（产码，非装配）**：`internal/panel/pump.go:200 L1Windows func() []L1WindowWait`、`:324 if p.src.L1Windows != nil {`、`:325 waits = p.src.L1Windows()`
- **供数侧新能力（票 220 AC#2 已落地的读侧）**：`internal/agent/approval/window_read.go:63 func (g *Gate) LiveL1Windows() []L1Window`、`:73 out = append(out, L1Window{CorrelationID: w.corr, TaskID: w.taskID, Tool: w.tool})`
- **仓内唯一的生产者＝测试**：`internal/panel/subagent_blocked_220_test.go:81 L1Windows: func() []L1WindowWait { return waits },`（`L1Windows:` 这一**装配形状**全仓**只有 1 命中**，且它在 `_test.go` 里）

尺二＝`grep -c L1Windows cmd/wisp/run.go` ⇒ **0**。我另把 `run.go:699` 那枚 `PumpSources{` 的字面量整块读了（`sed -n '690,715p'`）：列出的字段是 `Verdicts / Mode / Workspace / Git / Model / Credential / Results`（+:715 之后未完），**没有 `L1Windows`**。

⇒ 注释那一节 "and L1 windows **no run feeds**" **成立**：生产装配点没供数，仓内只有 fixture 供。
⚠ 两点我要具名：① 这一节的**根因**（票 220 落地腿那一跳没派）注释没写，票面 ⓒ 写了——注释不需要写根因，**不算假话**。② `window_read.go:29-30` 逐字 "internal/panel receives a copy of these ids through a READER the composition root supplies (PumpSources.L1Windows)" ——那是**规定应当如此**，不是"已经如此"；我的尺二证明 composition root 今天**没有**供。⇒ 注释没有把它写成"已经能点亮"（AC#2 的红线），**这一条我按字面过**。

---

## ⓓ "no tool_choice is sent, so mockllm never answers with a tool call, so a spawned child never reaches the gate" 两环 —— 【成立（第二环带一条我补的射程条件）】

**环 1（loop 不发 tool_choice）＝真。** 尺＝`grep -rn "ToolChoice" internal/agent/*.go` ⇒ **0 命中**（含测试；`wc -l`＝0）。
尺＝整族 `grep -rn "ToolChoice:" --include=*.go .` 剥 `_test.go` 剥 `.scratch` ⇒ 全仓非测试赋值点**只有 1 枚**：`internal/llm/probe.go:91 ToolChoice: &ToolChoice{Mode: ToolChoiceRequired}`（probe，不是 loop）。

**环 2（所以 mockllm 不答 tool call）＝在这把尺的形状下为真，但因果条件是"字段为 nil"，不是"惯例"。** 尺＝`sed -n '171,215p' tools/mockllm/chat.go`：

- `:172 if len(req.Tools) > 0 && req.ToolChoice != nil {` ——**整个 tool-call 决策块被 `req.ToolChoice != nil` 守着**；不发送 ⇒ 不进块 ⇒ `toolName` 恒空。
- `:179-181 case "auto": if strings.Contains(lastUser.text, "[tool]") { toolName = ... }` ——连 `[tool]` 提示词那条旁路**也在 `tool_choice == "auto"` 之内**，没有 tool_choice 就走不到。
- `:200-202 if toolName != "" { toolArgs = ... }`／`:207-213 finish = "tool_calls"` 仅在 `toolName != ""` 时——⇒ 环 1 推环 2 的因果**在 chat 方言上逐字成立**（`messages.go:120` 同样写 "Forced tool selection (**tool_choice** any | tool) with tools declared"）。

⚠ **但我必须记一条射程例外**：mockllm 有**第二条与 tool_choice 无关**的出话路径＝golden 回放（`tools/mockllm/chat.go:33 if name := goldenName(req.Model, r.Header.Get("X-Wisp-Golden")); name != "" {`），而仓里带 tool_call 的 golden **成族存在**（尺＝`grep -ln "tool_call" internal/agent/testdata/golden/*` ⇒ `budget-loop / disconnect-mid-toolcall / max-tokens-toolcall / parallel-tools / repeat-echo / six-tools / slow-tool / spill-tool / tool-then-text`＝**11 枚里 9 枚命中**）。
⇒ 精确说法应是"**这条路径**上（本件用 mockllm 的 echo，见 `subagent_carrier_197_test.go:61/67` 那两行逐字 "carrier197RootTask is long enough that mockllm's echo comes back in several text" / "mockllm's echo of the message the loop sent"）不发 tool_choice 就不会有 tool call"。注释把这条写在 "WHAT THIS FILE DID NOT MEASURE" 标题下、且 `:600` 尾部本就锁在 `:601`（编排者 C3 已裁"不许动 `:601`"）⇒ **判：注释未失真，因果链在本件射程内成立；"golden 回放可绕过 tool_choice"这一节属**射程外的知识**，不是本三行的假话。**

**环 3（所以被举的孩子到不了门）＝真，但我要把"本件的 spawn 本身不经模型"记在明处。**
尺＝`sed -n '172,175p' cmd/wisp/subagent_carrier_197_test.go`：

```
173: out, eerr := rt.bridge.Execute(ctx, agent.ToolRequest{
174:   TaskID: rootID, CorrelationID: rootID, CallID: fmt.Sprintf("spawn-197-%02d", i),
```

⇒ 孩子的**举起**不经 tool call（是测试直接把 `ToolRequest` 放上桥），孩子的**回合**里 loop 不发 tool_choice ⇒ mockllm 不答 tool call ⇒ 孩子自己**不会**再产生任何门控请求。`:174` 那一枚 `CorrelationID: rootID` 是**测试内注入的 corr==task id 形状**（票面 ⛔ 警告同族形状，出处 `internal/panel/subagent_blocked_220_test.go:151-152`）——它把的是**根**的 spawn 卡，`rootID` 不是子行的 `TaskID`，故**不会**点亮 `:236` 那把查询；**结论不受它影响**，但它是"仓里能造、生产里没人造"的第二枚证据（票面 ⓒ 的读法）。

---

## AC#1 小结

| 问 | 判语 | 关键逐字凭据 |
|---|---|---|
| ⓐ | 成立 | `subagent_roster_197.go:214`／`:219-220`／`:236`；`indexWaitingKey:267-268` 是我另量到的一节 |
| ⓑ | 成立（不可能，非仅不可达） | `loop.go:603/:604-607`＋调用点唯一 `:676`；射程外例外 `gate.go:282 orDefaultText` |
| ⓒ | 成立 | `grep -c L1Windows cmd/wisp/run.go`＝0；`L1Windows:` 装配形状全仓唯一命中＝`subagent_blocked_220_test.go:81` |
| ⓓ | 成立（三环在本件射程内） | `internal/agent/*.go` ToolChoice 0 命中；`tools/mockllm/chat.go:172` 守门；例外＝`chat.go:33` golden 回放可绕开 tool_choice |

腿在 AC#1 额外交的那一项（本件 spawn 不经模型 tool call，出处 `:173-174`）我**独立复现成立**，票面 ⓓ 没写它，这是**加法不是越界**。
