# 票 197 残余判据现读普查（取消语义 ＋ "卡在等批准"那一格）

> ⚠ **本件由编排者代落（09-29 10:5x）**。原只读腿 `197-c2` 交回了结论但**没写出文件**——因为我把它派成了**只读代理类型**（那类型没有写文件的工具），**这是我的派单错**，不是它的漏。
> ⇒ 下面每一行都**由编排者自己在当前工作树重跑过 `grep`/`sed -n`** 才落；每条都标了是哪一枚尺。**〔腿报，编排者未复核〕** 的格子一律不许当判据用。

- 起手锚点：`HEAD=83fd5883`（`de204ff1`＝账 `A418` 之后的两枚 `HANDOVER` 更正）
- 现树脏项：另有 **1 枚写腿正在动** `internal/agent/approval`／`internal/tools`／`cmd/wisp` ⇒ **本件不含金红读数**，一条 `go test` 都没跑（跑了也是脏的）。

---

## ① `task.spawn` 的执行体长什么样

| 事实 | 现读 | 尺 |
|---|---|---|
| 有没有裸 `go func(` | **零命中**（`internal/tools`／`internal/panel`／`internal/agent/approval`／`cmd/wisp` 非 `_test.go`） | `grep -rn 'go func('` |
| 孩子怎么起来 | `agent.New(opt)` → `child.RunAsync(childCtx, prompt)` ＝ **`internal/tools/subagent_197.go:303`／`:328`** | `sed -n` 逐行 |
| 起的口子有没有 owner | `RunAsync` 内部走 `internal/agent/loop.go:321-327`（先铸 `agent-task-<id>` 再 `reg.Spawn`），recover 边界在 `internal/observe/goroutine.go:287-330` | 〔腿报，编排者未复核那两段的行号；只复用了"不是裸 `go func(`"这条我已 `grep` 过的结论〕 |
| 收尾 watcher | `observe.Default.Spawn("subagent-finish-"+bg.ID, "tools", …)` ＝ **`subagent_197.go:345-354`**（父上下文是 `context.WithoutCancel`，见下） | `grep -n WithoutCancel` |
| **父取消会不会级联给孩子** | **不会**，且是刻意的：`subagent_197.go:327` `childCtx, cancelChild := context.WithCancel(context.WithoutCancel(ctx))`，上面 `:309`／`:323` 两处注释逐字「Parent cancellation must not cascade (ticket 197 §0)」 | `sed -n '305,330p'` |
| 池帽 | **`MaxConcurrentSubagents = 4`＝`subagent_197.go:78`**（占位口 `:254` `TryAcquireSubagentSlot(MaxConcurrentSubagents)`，超一枚硬拒带可读理由 `:259`） | `grep -n MaxConcurrentSubagents` |
| ⚠ 行号漂移更正 | 编排者旧账（`A418` 前后）记的是 `:73`，**现树真值在 `:78`**；枚数没变（仍是 4），**变的只有行号** | 同一枚 `grep` |
| 深度帽 | 两形：名册侧检查（`:248-252`）＋**结构性剥掉孩子身上的 `task.spawn`**（`:504-542`） | 〔腿报，编排者未复核；`Description()` 里那句「深度 `MaxSubagentDepth`」在 `:192` 有对得上〕 |

---

## ② 父取消不级联：这一条做到了，**但孩子一枚都没有停止出口**

**做到了的两半**（都在盘上）：
1. 级联确实不断——见 ① 的 `:327`。
2. **"没级联"这件事写给模型看了**：`subagent_197.go:360-365` 那段 `Result.Text` 里逐字含「没有被级联取消」，并经由 `internal/agent/loop.go:712` 进历史；`Description()`（`:190` 一带）也带同一句。
   ⇒ 编排者旧账里"要写给模型看"这一条**已落地**，不必再派。〔`loop.go:712` 那一跳＝腿报，编排者未复核〕

**没做的两半（＝票 197 AC#6 的真实缺口）**：
1. **`TaskRoster.Cancel` 存在但生产零调用者**：定义＝**`internal/tools/task.go:442`** `func (r *TaskRoster) Cancel(taskID string) (bool, string)`；全仓非测试对它**唯一命中是 `subagent_197.go:39` 的一句注释**（"back by TaskRoster.Cancel"），**没有一处调用**。
   `grep -rn '\.Cancel(' internal cmd \| grep -v _test` 交回来的全是别的东西（`RunningTask.Cancel`＝`loop.go:300`、上下文 root 的 `root.Cancel()`＝`loop.go:531`／`goroutine.go:324`／`logging.go:116`／`cmd/balldebug`）——**名册那一枚没人按**。
2. **`task.cancel` 这枚工具今天不存在**：`internal/tools/task.go:20-27` 头部注释逐字把 `task.list`／`task.cancel` 写成「DEFERRED with five fields, PLAN.md §7 :1531」，而 `BuiltinTaskEntries`（**`:595`**）**只返回 `task.output` 一枚**。
3. ⛔ **同时，`subagent_197.go:189` 当着模型许诺了它**：逐字「它会在任务名册里留下一行有父子关系与状态的记录，**可以用 `task.cancel` 单独停它**；」。
   ⇒ **"名册有一行、能单独停它"这两句一起构成了对模型的谎报**：名册那行有（②-1 的定义在），**停它的那枚工具没注册，而唯一能停它的函数没有生产调用者**。

**⇒ 这一处是编排者旧账没写到的新洞**（票 197 的 AC 表里只有"取消语义怎么收"，没有"文案对模型许诺了不存在的工具"）。已单独立成 **票 221**。

---

## ③ 名册那一格能不能表达"被卡住 / 被取消"

`TaskRowView` 的 JSON 键（`internal/panel/subagent_roster_197.go`，逐枚现读）：

| 键 | 行 | 能不能表达今天关心的事 |
|---|---|---|
| `taskId` `label` `kind` `parentTaskId` | `:102` `:103` `:104` `:105` | 父子关系有 ⇒ 界面上"谁生的谁"能画 |
| `status` | `:109` | 只有 `subagent_197.go:99-104` 那四枚态在册：**Thinking / Settling / Muted / Error** |
| `statusKnown`／`statusReason` | `:110`／`:115` | "不知道"与"知道是没有"分得开（票 176/179 那一族的延长线） |
| `streamKey` | `:121` | 「点进去看它自己的流」的寻址键，单点铸造在 `internal/streamkey` |
| `blockedOnApproval` | `:127` | ⚠ **一枚布尔，不是一枚态**——见 ④，它今天只看得见 L2 |
| `streamTruncated`／`streamElidedRunes`／`streamDropped` | `:130`／`:131`／`:132` | 流的三段诚实旗 |

**⛔ D43 里那两枚"等批准"的态今天在生产上零使用**：`internal/statemachine` 的名册里有 `Confirming`（`states.go:21`）与 `AwaitingApproval`（`:22`），但〔腿报：全仓 `grep` 零生产使用；编排者未复核枚数，**只确认 `blockedOnApproval` 是布尔、不是态**这一条』。
⇒ **"被卡住"今天不是 D43 的一枚态，是一枚外挂布尔**。这不算违规（不许新造态名），但意味着界面画"卡在等批准"只能靠那枚布尔，而它按 ④ 只有一半是真。

---

## ④ ⚠ 我那条"L1 短窗口里等待的任务在名册上恒读'没被卡住'"——**证实，而且比我想的更漏**

**判据链，逐跳现读**：
1. `subagent_roster_197.go:190-210`：先把待批卡收成 `waiting[card.CorrelationID] = true`（**键是 correlation id**），再在 `:210` 用 **`row.TaskID`** 去查：`BlockedOnApproval: row.TaskID != "" && waiting[row.TaskID]`。
2. 卡从哪来：泵读 `p.src.Verdicts()` → `cmd/wisp/run.go` 的 `rt.liveVerdicts` → `Queue().LiveApprovals()`（`internal/panel/pump.go:196-204`、`:263`；`cmd/wisp/panel_pump.go:62`）。
3. `LiveApprovals()`（**`internal/agent/approval/pending_read.go:106-121`**）**只走 `q.pending`**——那是 **L2 队列**。
4. **L1 短窗口不住在队列里**：它住在 `Gate.windows`（`gate.go:302-316` 开／关；`:377` 只有 `Veto` 按名查），**全仓没有任何枚举口**（`grep -rn 'g\.windows\|windows\['` 非测试只剩那四行自己用）。
   ⇒ **结论：今天一枚正在 L1 短窗口里等的任务，在名册上读起来是"没被卡住"。这一格是真漏，不是措辞问题。**

**⚠ 第二枚漏（腿没报，编排者顺着 `queue.go` 读出来的）**：同一枚任务**第二次**排队时，correlation id 会被**改写**成 `原id#序号`（`internal/agent/approval/queue.go:152-154`：`if _, clash := q.byID[corr]; clash { corr = fmt.Sprintf("%s#%d", d.CorrelationID, q.seq) }`）。
而 join 用的是 `row.TaskID` ⇒ **"一张卡"时对得上（因为 corr 缺省回落到 taskID：`gate.go:247`／`:471` 的 `orDefaultText(d.CorrelationID, d.TaskID)`），"同一枚任务第二张卡"时对不上 ⇒ 那一格会静默读成"没被卡住"。**

⇒ 这两枚都不是"界面没画"，是 **Go 侧送出去的布尔本身就是假的**。已单独立成 **票 220**（含两形正控）。票 219 的 AC#6 早就指过这一族，但把它写成"要么一起修、要么具名留洞"——现在有具名票了。

---

## ⑤ AC#5：子代理能不能自己批自己——**今天"不能"，但没有任何尺会因为有人做到而红**

**"不能"是真的，且是结构性的**：
- 递给工具的那枚 `Gate` 接口面**只有两枚方法**：`PendingWindow`／`PendingApproval`（`internal/agent/approval/gate.go:99-106`）——**都是"发问"，没有"答复"**。
- `grant` 只往一个方向流：交给 `approval.UI.Prompt`（`gate.go:484`、`ui.go:42`），而宿主那枚 `consoleApprovalUI` **根本不把它印出来**（`cmd/wisp/run.go:1012-1035`）。
- 子代理拿到的依赖束被钉死：`SubagentDeps` 的字段集合在 `internal/tools/subagent_197_test.go:827`（`reflect.TypeOf`）＋`:872` 那条 `t.Errorf`（红句逐字「SubagentDeps 字段集合变了：%v（多出来的每一枚都要问一句是不是新的允许出口）」）。

**但"没有尺"也是真的（＝AC#5 今天＝〔仅文档〕）**：
- `grant` 绑的是**事**不是**人**：`bindDigest(corr, taskID, tool, level string, seq uint64, args []byte)`＝`approval.go:253`，摘要里**没有"谁在批"这一维**。
- **`Request.Source` 今天不参与判定**：`approval.go:231` 注释逐字「caller-settable `Request.Source` string is **logged and never consulted**」，`gate.go:611`／`:618` 也只把它塞进日志（`claimed_source=%q` 这个写法本身就写着"未核实的自称"）。
- 已有的 `Test197SubagentHasNoSelfApprovalOutlet`（`subagent_197_test.go:838`）**射程只有两格**：`AdmitTask == nil` 时拒绝 ＋ 字段枚数。
  ⇒ **造一枚"孩子自批"的假腿，今天没有任何一枚尺会红**——这正是票 197 AC#5 要的正控（"假腿 ⇒ 要红"）现在**还不存在**。

**最便宜的合法接缝（不新造导出方法）**：`cmd/wisp` 的 `package main` 测试——**只有那一处同时握着 concrete `*Gate` 和 spawn 那条路**：把 `rt.ui` 换成一枚会记下 `grant` 的假 `approval.UI`，然后断言 `rt.gate.Native().Allow(childCorr, capturedGrant)` **必须失败**。
⚠ 派写腿前**必做撞钉预检**：`ui.go:155` 的 `PanelAPI` 没有 `Allow`＝合法负向锚；`gate.go:572` `Native()` 目前生产零消费者（账 `A418` ③）⇒ 这枚正控**会第一次给 `Native()` 一个调用者**（在测试里），要检查它会不会把票 145/146 那族名册差集尺打红。

---

## ⑥ 名册差集（票 197 的 AC 逐格，只标"盘上有没有"）

| AC | 现在 | 缺什么 |
|---|---|---|
| AC#1 实体 | 〔已证〕 | — |
| AC#2 每枚自己的流 | 〔已证〕 | — |
| AC#3 载体（名册进快照） | 〔已证〕 | 界面那一跳不在本编队地界 |
| AC#4 阻塞态可见 | 〔建了但没接〕 | **只报 L2、且第二张卡静默掉格**＝④ 两枚洞 ⇒ **票 220**；这一格**不许当已满足** |
| AC#5 子代理不自带"允许"出口 | **〔仅文档〕** | 结构性不能，但**零尺会红** ⇒ 正控未存在＝⑤ |
| AC#6 取消语义 | 〔建了但没接〕 | `TaskRoster.Cancel`（`task.go:442`）生产零调用者；`task.cancel` 未注册（`task.go:595`）；而 `subagent_197.go:189` 已对模型许诺 ⇒ **票 221** |

**下一枚最便宜的写腿（编排者自己判，不派到别的票上）**：
1. `internal/agent/approval/gate.go`：给 L1 窗口**补一枚只读枚举口**（照 `LiveApprovals()` 的形状，不动状态机）；
2. `cmd/wisp/panel_pump.go`：把那些窗口喂进 `TaskRosterSectionFrom`，并把 join 的键**改成 corr 与 taskID 双查**（修 ④ 的第二形）；
3. `internal/tools/subagent_197.go`：`Description()` 那句 `task.cancel` 许诺**要么兑现要么删**（**不许留"许诺了但没接"**）；
4. 一枚 `cmd/wisp` 自批正控测试（⑤ 的接缝）。
⇒ 这枚腿**不新增任何导出方法名**（枚举口是新方法——⚠ **这一条要先裁**：票 194 是名册补齐的射程，票 197 的禁区写着"不新增方法名"；**编排者处置＝把它拆成两枚**：220 只做"把已有数据读对"（不新增方法），新增枚举口那一步**并进票 194 的批准面**再去要一句话）。

---

## ⛔ 本件丢弃的一格（合规，不是质量）

原腿在 ⑥ 里交了**一行来自禁读目录的引用**。本编队的规矩是那两棵目录**不读、不引、不转述**，所以**那一格整行删除，不进任何票面或台账**；它原本在说什么，我不在这里复述（复述本身就是违规的那一半）。界侧要什么，由 owner 自己带给他用的那枚 agent。
⇒ **错在我**：派单里我只写了"别读那两棵目录"，**没写"你的结论也不许引它"**——派单模板已据此改（后续只读腿一律写"结论里也不许出现那两棵目录"）。
