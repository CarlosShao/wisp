# 票 292 腿 292-r1 — 起手锚 + AC#1 逐环裁真伪 + 用过的尺

本件零产码改动（AC#1 明文「只读」）。所有引文为落笔前在盘上现跑的逐字原文。

## 1. 起手锚

- 钟点（`date '+%Y-%m-%d %H:%M:%S %z'` 的 stdout）：`2026-10-09 17:01:54 +0800`
- HEAD（`git log -1 --format='%H %ad %s' --date=iso`）：
  `a3782ff067ffe77b064621e55578c215bbb468b9 2026-10-09 16:58:02 +0800 290-v1 收口：台账 A785＋停车点 §4.0ap（四格判语入账＋记我自己两处纪律缺陷）`
- 分支：`dev`
- 工作树既有未入库改动（非本腿的活，未碰、未提交、未还原）：`.gitignore` 与 `.scratch/wisp/probes/**` 若干（152/161/242/268）为 ` M`；`design/**` 16 枚为 ` D`。
  派单说「17 枚 = 16 枚 ` D design/**` + ` M .gitignore`」；盘上 `git status --porcelain | head -30` 的前 30 行里 ` M` 的 probes 件有 12 枚（152 一枚、161 九枚、242 一枚、268 一枚），` M .gitignore` 一枚，` D design/**` 16 枚 —— 见本件 §6 不符点 N1。

## 2. 那三行在不在、行号对不对（`sed -n '598p;599p;600p' cmd/wisp/subagent_carrier_197_test.go`）

```
:598  //  1. blockedOnApproval is never asserted on a production packet. The join needs a card
:599  //     whose correlation id equals a CHILD task id, and nothing in this tree can make
:600  //     that happen from the outside: the agent loop never sends tool_choice, so mockllm
```

行号与票面「现量」逐字一致。票面引文省略了行首的 `//` 空格排版，本处为盘上原文。

## 3. AC#1 四问逐问（每问带 `文件:行` 逐字引文）

### 问 ⓐ — join 的来源枚数（对 `internal/panel/subagent_roster_197.go:212-221`）

判语：**票面说的两枚来源成立；原句 `:598-599` 只说了半个 join。**

建表处逐字（同文件行号已现场点数核对）：

```
:212 	waiting := make(map[string]bool, len(pending)+2*len(l1Windows))
:213 	for _, card := range pending {
:214 		indexWaitingKey(waiting, card.CorrelationID)
:215 	}
:216 	for _, w := range l1Windows {
:217 		// Both halves of one window: an L1 window is keyed by its correlation, and
:218 		// the row it blocks is named by its task.
:219 		indexWaitingKey(waiting, w.CorrelationID)
:220 		indexWaitingKey(waiting, w.TaskID)
:221 	}
```

查询处逐字：

```
:236 			BlockedOnApproval: row.TaskID != "" && waiting[row.TaskID],
```

`:212` 的容量式 `len(pending)+2*len(l1Windows)` 本身就写着「两枚来源、L1 那枚每个占两把钥匙」。 ⇒ `:220` 那一半允许**只凭 task id** 点亮行，比原句「needs a card whose correlation id equals a CHILD task id」宽。原句**结论还对、这一节的射程说窄了**。

### 问 ⓑ — 「card corr == CHILD task id」今天是不可能还是只是不可达（对 `internal/agent/loop.go:603/:676`）

判语：**不可能（结构上），原句写成了「只是从外面做不到」——这一节也错了。**

```
internal/agent/loop.go:603 func callCorr(taskID, callID string, index int) string {
internal/agent/loop.go:604 	if callID == "" {
internal/agent/loop.go:605 		return fmt.Sprintf("%s#call-%d", taskID, index)
internal/agent/loop.go:606 	}
internal/agent/loop.go:607 	return taskID + "#" + callID
internal/agent/loop.go:608 }
```

```
internal/agent/loop.go:676 		TaskID: taskID, CorrelationID: callCorr(taskID, p.call.ID, i),
```

两支返回值都带 `#` 分隔的追加段（`#call-N` 或 `#<callID>`），对任何非空 `taskID` 恒不等于 `taskID`。 ⇒ 由 loop 派发的卡片，其 corr **结构上不可能**等于任何 task id。

### 问 ⓒ — 「nothing in this tree can make that happen」的射程（仓内 fixture vs 生产装配）

判语：**这句话今天仍成立，但成立的方式变了；原句把根因归在卡片那一半上，归错了。**

尺 1（仓内 L1 生产者枚数）＝`grep -rn "L1Windows:" --include=*.go .` ⇒ **1 命中**，且是测试件：

```
./internal/panel/subagent_blocked_220_test.go:81:		L1Windows: func() []L1WindowWait { return waits },
```

尺 2（生产装配点里 L1Windows 有无值）＝`grep -c L1Windows cmd/wisp/run.go` ⇒ **0**。装配点逐字：

```
cmd/wisp/run.go:699 	rt.pump = panel.NewSnapshotPump(panel.PumpSources{
```

其字段实参（`:700-712` 现量）只有 `Verdicts` / `Mode` / `Workspace` / `Git` / `Model` / `Credential` / `Results` …，**无 `L1Windows`**。

那枚 fixture 逐字（票面 ⛔ 明令不得读成「生产能造出来」，此处照办）：

```
internal/panel/subagent_blocked_220_test.go:151 	data := blockedPacket(t, nil, []L1WindowWait{{CorrelationID: "corr-host-9", TaskID: "child-1"}})
internal/panel/subagent_blocked_220_test.go:152 	assertBlocked(t, data, "child-1", true)
```

⇒ 「仓里能造（测试内注入）、生产里没人造」为真。**根因＝给 `PumpSources.L1Windows` 供数据那一跳还没人做（票 220 落地腿的活），不是那句写的「卡片造不出来」。**

### 问 ⓓ — 「so mockllm never answers with a tool call, so a spawned child never reaches the gate」两环今天真伪

判语：**两环今天都仍为真，但第一环的机制要补一句才说实话；而本票这条路径上的 spawn 本身并不来自 mockllm 的 tool call。**

环 1「the agent loop never sends tool_choice」＝真。尺＝`grep -rn "ToolChoice" internal/agent/*.go | grep -v "_test"` ⇒ **0 命中（rc=1）**；`grep -rn "ToolChoice" internal/agent/` 全量亦 **0 命中**。全仓非测试的赋值点只有 probe，不是 loop：

```
internal/llm/probe.go:91 					ToolChoice: &ToolChoice{Mode: ToolChoiceRequired},
```

环 2「so mockllm never answers with a tool call」＝真，且机制是 mockllm 自己的门：

```
tools/mockllm/chat.go:172 	if len(req.Tools) > 0 && req.ToolChoice != nil {
```

⇒ `tool_choice` 缺席时 `req.ToolChoice == nil`，整个 tool-call 决策块不进，`toolName` 恒为 `""`，`finish` 落在 `"stop"`。注意 **`auto` + `[tool]`  marker 那一支（`chat.go:178-182`）也要求 `tool_choice` 非 nil**，所以「loop 不发 tool_choice」这一条对 mockllm 是**充分**的；仓里那些 `tool-call.sse` golden 是适配器级（`internal/llm/adaptertest/harness.go:86 ScToolCall ScenarioID = "tool-call"`），由测试直接喂，不经 loop。
 ⇒ 原句把「不发 tool_choice」说成 mockllm 举不起工具的唯一理由，**结论对、措辞偏窄**（真机制是 mockllm 把 tool-call 整个决策挂在 `ToolChoice != nil` 上）。

环 3「so a spawned child never reaches the gate」＝真（子任务自己那一轮无 tool call ⇒ 无 ToolRequest ⇒ 无卡）。
 ⚠ 但本件的 spawn **不是**模型 tool call 派发的，是测试把一条 ToolRequest 直接放到被装配好的桥上（票面 §15-19 已自陈「nothing may call the spawner by hand」的取舍）：

```
cmd/wisp/subagent_carrier_197_test.go:173 			TaskID: rootID, CorrelationID: rootID, CallID: fmt.Sprintf("spawn-197-%02d", i),
cmd/wisp/subagent_carrier_197_test.go:174 			Name: "task.spawn", Args: args,
```

 ⇒ 「卡片那一半在本路径上没人派发」这条**不因**「loop 不发 tool_choice」而成立，它成立是因为测试没往桥上放需要门控的调用。票面把这一环交给我裁，我裁：**结论真，因果链在此件里是多余的**，故改写后的注释**不再沿用这条因果**，只陈述观察到的事实。

### 票面 :603 那句引文顺带核了一下，是真话

```
internal/panel/subagent_roster_197_test.go:103 func TestTheRosterReaderPutsSubagentsOnTheWire(t *testing.T) {
```

（`:603` 原文称它在 internal/panel 的 pump 上对真 `ApprovalCardView` 断言 —— 件与函数名都在 HEAD 上，非假凭据。）

## 4. 搜过的字样名册（AC#1 完成判据要求「不许只答没找到」）

| 字样 | 尺 | 结果 |
|---|---|---|
| `L1Windows:` | `grep -rn "L1Windows:" --include=*.go .` | 1 命中（测试件 `:81`） |
| `L1Windows` | `grep -c L1Windows cmd/wisp/run.go` | 0 |
| `NewSnapshotPump(panel.PumpSources` | `grep -n` on `cmd/wisp/run.go` | 1 命中 `:699` |
| `ToolChoice` | `grep -rn internal/agent/*.go` 去 `_test` | 0 命中 |
| `ToolChoice` | `grep -rn internal/agent/`（含测试） | 0 命中 |
| `ToolChoice` | `grep -rn --include=*.go .` 去 `_test` | 命中全在 `internal/llm/**`（content.go / probe.go / anthropic/request.go），无一在 loop |
| `ToolCall` / `tool_call` | `grep -n` on `internal/llm/adaptertest/mockllm.go` | 0 命中（该件不产 tool call；tool call 场景在 `harness.go` 与 `tools/mockllm/chat.go`） |
| `tool_choice` / `ToolChoice` | `grep -n tools/mockllm/*.go` | 命中 `chat.go:64/:172/:174/:183/:193/:230`、`messages.go`、`capability_test.go`、`main.go:28` |
| `TestTheRosterReaderPutsSubagentsOnTheWire` | `git grep -n` | 定义在 `internal/panel/subagent_roster_197_test.go:103` |
| `func indexWaitingKey` / `type L1WindowWait` | `git grep -n` | `subagent_roster_197.go:262` / `pump.go:127` |
| 引用件名存在性 | `git ls-tree -r --name-only HEAD \| grep -qx` | `subagent_roster_197.go`、`subagent_blocked_220_test.go`、`loop.go`、`run.go`、`probe.go`、`chat.go` 六枚全 OK |

## 5. 本腿接下来的尺（预告，读数在 10/20/30 三件）

AC#2 只改 `:598-600` 三行注释；不动 `:311-313` 那句 `t.Error`；不动本件其它注释。
AC#3 整包 harness 逐字 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1 -v`，改前 2 发、改后 2 发，每态取交集后逐名作差，每发前 `tasklist` 现跑 `wisp.exe`/`balldebug.exe` 枚数。

## 6. 派单/票面与盘上不符之处（具名）

- **N1**：派单说工作树「本来就有别人未入库的 17 枚改动：16 枚 ` D design/**` ＋ ` M .gitignore` ＋ `.scratch/wisp/probes/{152,161,242,268}/**`」；`head -30` 截断后我只能确证 16 枚 ` D design/**`（前 30 行内）＋ ` M .gitignore` ＋ 12 枚 ` M` probes。枚数口径取决于是否被 `head -30` 截掉，**我没有把它当 17 枚这个数去核**，只做了一句：这些件我一枚没碰。（不构成票面缺陷，属派单转述的计数口径。）
- **N2**：票面 `:11` 与 `:12` 的引文把 `//` 前缀与行内空格压平成 `//  1.` / `//     ` 形状，与盘上逐字一致，**无冲突**；票面 `:13` 把 roster 的 `:212/:214/:219/:220/:236` 五处行号写出，我在本件 §3 问 ⓐ 现场点数复核，**五个行号全对**。票面 `:14` 的 `loop.go:603/:676` 同样复核为**全对**。⇒ 本票未发现票面与盘上的实质冲突。
- **N3**（给编排者的信息，不是不符）：票面 `:16` 说「再下一环 `so mockllm never answers with a tool call / so a spawned child never reaches the gate` 我没裁＝归本票 AC#1」。我在 §3 问 ⓓ 裁了，并额外发现本件的 spawn 不经模型 tool call（`carrier_197_test.go:173-174`），这一条票面没写，我把它的**逐字读数**与**由此得出的「因果链多余」判语**一并交回，不擅自扩写注释之外的任何东西。
