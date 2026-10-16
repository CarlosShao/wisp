# 328-a1 只读普查 件① — `callCorr` 调用源名册 + 每形 ⇄ 名册可解析性

起手锚 `44f51fc8`。尺全部在 dev 工作树现跑；`internal/agent`＋`internal/panel` 在 `44f51fc8..HEAD` 之间 `git diff --stat` 为空
（⇒ 工作树＝锚点字节，读盘安全；并发腿未碰这两面）。go1.27.1（`go version` 自报，⛔ 未跑 test/build/vet/门禁）。
行号只当读数。"算的"＝按代码那行拼接逐字推，非跑出来。

## ① 调用源名册（今天真会走到 `callCorr` 的）

尺：`grep -n 'callCorr' internal/agent/*.go`（rc=0）。全仓 **`callCorr` 只有一处函数体、一处生产调用点**：

- 定义 `internal/agent/loop.go:603`，两支：`:605` `fmt.Sprintf("%s#call-%d", taskID, index)`（callID 空）／`:607` `taskID + "#" + callID`（callID 非空）。
- **唯一生产调用点 `internal/agent/loop.go:676`** `TaskID: taskID, CorrelationID: callCorr(taskID, p.call.ID, i)`（在 `executeCalls` 内）。其余 `callCorr` 命中全在 `_test.go`。

父链（逐枚读回，尺 `grep -rn 'executeCalls' internal/ --include=*.go|grep -v _test` rc=0）：

- `executeCalls` 被且只被 `loop.go:504` 调（在 `run` 的 turn 循环里，`l.append(assistantTurnMessage(turn))` 之后）。
- `run`（`loop.go:346`）由两枚入口进：`RunAsync`（`:328`→`:333 l.run(c,id,input)`）与 `Run`（`:339`→`:340 l.run(ctx,newTaskID(),input)`）。`taskID`＝`newTaskID()`（`:1157`）＝UUIDv4 形（`:1169` `h[0:8]+"-"+…`，hex+连字符，**不含 `#`**）。
- `Run/RunAsync` 的**生产**入口共 2 枚（尺 `grep -rn '\.RunAsync(\|\.Run(' internal/ cmd/ --include=*.go|grep -v _test|grep -v 'func (l'` rc=0）：
  1. `cmd/wisp/run.go:1164` `bg := loop.RunAsync(ctx, task)` — 半双工干活路径（CLI/文本循环）。
  2. `internal/tools/subagent_197.go:367` `bg := child.RunAsync(childCtx, prompt)` — 子代理 child，仍是文本循环、仍半双工。
  （`balldebug`/`panel_resident`/`resident_ball` 的 `.Run(` 是 `bridge.Run` 热键/球桥，不建 `agent.Loop`，不数。）

两支各走什么（**算的**；`p.call.ID`＝`ToolCall.ID`＝`ev.ToolCallID`，`llm/events.go:227/235`，是 provider 线值的逐字透传）：

| 入口 | 走的支 | tail（末 `#` 后那串）来源 | 尺 |
|---|---|---|---|
| 干活路径·模型给了调用号 | 支 A `loop.go:607` | provider call id 原样 | openaichat `toolasm.go:59 acc.id=f.ID`；空号回退 `:114 call_unnamed_%d`；responses `stream.go:447-450` 取 `call_id` 否则 `id`；anthropic `stream.go:258 st.id=blk.ID` |
| 干活路径·模型没给号 | 支 B `loop.go:605` | 字面 `call-<index>`（含字母与 `-`） | 代码逐字 |

两条路径到不到那里（D47）：
- **半双工干活路径**：到（上表 1/2）。
- **陪聊全双工路径（Path C）**：**到不了**。`run.go`/`subagent_197.go` 是文本循环；Path C（`internal/audio/gate.go:17`）不建工具调用，进不了 `executeCalls`。
- **realtime 那支**：**到不了**。`approval/gate.go:182-194` 只有文本循环持 `AdmitTextTask`（`loop.go:374-375 opt.AdmitTask`）；`admitCheck`（`gate.go:237-238`）对 Path C 明拒"实时语音（Path C）不具备工具权限"。全仓 `grep -rln realtime internal cmd|grep -v _test`（rc=0）只命中 config/approval/panel 的注释，**无独立 realtime Go 包**。⇒ realtime 大脑零工具权限，永不产 corr。

⇒ **名册结论**：今天真产 corr 的调用点只有 `loop.go:676` 一枚，只在半双工文本干活路径的支 A/支 B 两支被走到；陪聊全双工与 realtime 两支到不了。

## ② 每形 corr 实例 ⇄ 名册可解析性（**算的**）

名册把 L2 卡片点亮任务行的判据（尺 `sed -n '212,236p' internal/panel/subagent_roster_197.go` rc=0）：
`:213-214` 对每张卡只 `indexWaitingKey(waiting, card.CorrelationID)`；`:236` `BlockedOnApproval: row.TaskID!="" && waiting[row.TaskID]`。
`indexWaitingKey`（`:262-270`）把 key 逐字入库，再当 `splitClashSuffix(key)` ok 时把 base 也入库。
`splitClashSuffix`（`:277-288`，唯一读者是 `:267`）＝末 `#` 之后**至少一个且全 ASCII 数字**才 ok，返回 base＝末 `#` 前全部。
对比：L1 窗口在 `:216-221` **两枚都索引**（`w.CorrelationID`＋`w.TaskID`），所以 L1 能点亮、L2 卡不能——差的就是 L2 卡不带任务号那枚。

设真实任务号 `T = 8f14e45f-ceea-467a-9a1b-3c2d4e5f6a7b`（UUID 形，无 `#`）。逐形跑那把尺：

| 形（算的） | 产生处 | 末`#`后那串 | 尺判 ok | base | 与 `waiting[T]` join？ |
|---|---|---|---|---|---|
| `T#call-0`（支B，index=0） | `loop.go:605` | `call-0`（含 `c`） | false（`:283` 首字符 `c<'0'`） | — | **断**：仅逐字 `T#call-0` 入集，无 base，`waiting[T]` 不被置 |
| `T#call_abc123`（支A，openaichat） | `loop.go:607` | `call_abc123` | false（`_`/字母） | — | **断**（同上） |
| `T#call_unnamed_0`（支A，openaichat 空号回退） | `toolasm.go:114`→`loop.go:607` | `call_unnamed_0` | false | — | **断** |
| `T#toolu_xyz…`（支A，anthropic） | `stream.go:258`→`:607` | `toolu_…` | false | — | **断** |
| `T#call_…`（支A，responses 取 `call_id`） | `stream.go:447`→`:607` | `call_…`/`id` | false | — | **断** |
| 队列撞名改写（叠在上面任一形上） | `queue.go:168` | 末段 `=<seq>` 数字 | **true** | base＝改写入 `d.CorrelationID`＝`T#call-0` 或 `T#call_x`（本身带 `#`）≠ T | **仍断**：base≠T，逐字串≠T |
| 队列 replay 撞名 | `queue.go:616` | `=<seq>` 数字 | true | base＝`corr#replay` 去尾…仍≠T | **断** |

⇒ **join 得上唯一的充要条件**（算的）：卡片 corr 必须＝`T`（逐字），或＝`T#<全数字>`（恰好一个 `#`，剥出 base=T）。
`callCorr` 两支**都不落在这两个形上**：支 B 永远带 `call-` 前缀；支 A 的 tail＝provider 号逐字，三个真在用的 provider 都给带前缀的号（`call_`/`toolu_`/`call_unnamed_`），全非纯数字。
残余口子（⛔ 代码不保证、仅经验）：若某 provider 真回一个**纯 ASCII 数字且不含 `#`** 的 call id，支 A 会拼出 `T#<数字>` → 那把尺能剥出 base=T → 会点亮；但仓内没有任何代码造这种 id，且配置里的三家 provider 不发这种。

⛔ 未拿测试夹具 `child-1#9`/`task#2` 当生产形（那是名册用例 `card197()` 手造的撞名形，不经 `callCorr`；见件②/件③计数）。
