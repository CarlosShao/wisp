# 220 — **名册上"卡在等批准"那一格今天只看得见 L2，而同一枚任务的第二张卡会静默读成"没被卡住"**（Go 侧送出去的布尔本身就是假的，不是界面没画）

- Status: **待派**。来源＝只读腿 `197-c2` 的 ④ 节＋编排者顺着 `queue.go` 读出的第二形；现读证据件＝`.scratch/wisp/probes/197/r5c/census.md` ④（**每一行都由编排者自己在现树重跑过尺**）。
- 与票 197 的关系：票 197 的 AC#4（阻塞态可见）**今天不许算已满足**——它勾上的那半是"载体里有这一格"，不是"这一格读得对"。本票是那半之后的"读得对"。
- 与票 219 的关系：票 219 AC#6 早就指过"等批准时名册恒读没被卡住"这一族并写"要么一起修、要么具名留洞"⇒ **本票就是那枚具名洞**。

## 现量（起手逐条复算，别信这里的行号）

| 事实 | 现读 | 尺 |
|---|---|---|
| 那一格怎么算出来的 | `internal/panel/subagent_roster_197.go:190-210`：先 `waiting[card.CorrelationID] = true`（**键＝correlation id**），再 `BlockedOnApproval: row.TaskID != "" && waiting[row.TaskID]`（**查＝task id**） | `sed -n '185,215p'` |
| 为什么"一张卡"时对得上 | corr 缺省会回落到 taskID：`internal/agent/approval/gate.go:247`／`:471` 的 `orDefaultText(d.CorrelationID, d.TaskID)` | 现读 |
| 卡的来源只覆盖 L2 | 泵读 `p.src.Verdicts()` → `Queue().LiveApprovals()`（`internal/panel/pump.go:196-204`、`:263`；`cmd/wisp/panel_pump.go:62`），而 `LiveApprovals()`（`internal/agent/approval/pending_read.go:106-121`）**只走 `q.pending`**＝L2 队列 | 现读 |
| ⛔ L1 短窗口今天不可枚举 | L1 住在 `Gate.windows`（`gate.go:302-316` 开／关，`:377` 只有 `Veto` 按名查），**全仓没有任何枚举口** | `grep -rn 'g\.windows\|windows\[' \| grep -v _test` ＝只剩那四行自己用 |
| ⛔ 同任务第二张卡会掉格 | `internal/agent/approval/queue.go:152-154`：`if _, clash := q.byID[corr]; clash { corr = fmt.Sprintf("%s#%d", d.CorrelationID, q.seq) }` ⇒ 改写后的 corr 变成 `原id#序号`，**用 task id 去查必然落空** | 现读 |
| 态名不能顺手新造 | `blockedOnApproval` 今天是一枚**布尔**（`subagent_roster_197.go:127`），不是 D43 的态；D43 名册里 `Confirming`／`AwaitingApproval` 存在但生产零使用〔枚数＝腿报，未复核〕 | `statemachine/states.go:21-22` |

## 为什么这不是措辞问题

一个正在等用户点头的任务，在名册上读起来和"根本没在等"一模一样 ⇒ **"它卡住了"这件事对用户不可见**，而票 219 的三枚答复按钮、票 201 的答复听众全都建在"看得见有哪枚在等"之上。⚠ 更要紧的是**第二形**：同一枚任务连着两次要批准时，第二张卡**静默**不计数 ⇒ 越忙的任务越读起来越像没在等，这一形今天没有任何尺看得见。

## 判据（每格都要现跑读数）

- [ ] **AC#1 先把 join 读对（不新增任何导出方法）**：把 `blockedOnApproval` 的查键改成 **corr 与 taskID 双查**，并把 `#序号` 那一形纳入匹配（同一任务的多张卡只要有一张 pending 就算"在等"）。正控＝种"同任务两枚卡"样本 ⇒ 改前**那一格为假**、改后**为真**。
- [ ] **AC#2 让 L1 可见**：**落点必须先裁，二选一并写明选了哪一支**——(甲) 给 `Gate` 补一枚**只读枚举口**把 `windows` 暴露给泵（⚠ 这是新增方法名，票 197 的禁区写"不新增方法名"、名册补齐属票 194 的射程 ⇒ 走这一支要先在台账落一枚批准记录）；(乙) 让 L1 窗口在等的时候也进 `LiveApprovals()` 的视图（不动方法名，动的是那枚视图的语义，**要一并复查票 146 对 `LiveApprovals` 的既有判据不被顶红**）。
  ⇒ 判据本身：种一枚正在 L1 短窗口里等的任务 ⇒ 名册那一格**必须为真**；**改前读数要当场跑出来贴回票上**（不许写"应该会红"）。
- [ ] **AC#3 两形各有正控**：L1 等待中／同任务第二张卡，**逐枚**"种 X 必响"；本仓规矩＝负向尺必配正控。
- [ ] **AC#4 不许把布尔改成态、不许新造 D43 名**：改状态拼写＝契约级（D43 转移表冻结），本票射程外。
- [ ] **AC#5 整包终态读数**：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/...` 跑到终态、逐名比红名集合（**`-run` 单跑不算终态**），并**具名写清"这一格修完之后界面上还差哪一跳"**——那一跳属界面地界，只写票面、不代做。

## 禁区

`frontend/**`／`design/**` 零写面（连内容都不转述）；`PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`／C17 白名单既有名字一字不动；不新增方法名（要走 AC#2 甲形先落批准记录）；`internal/panel/tokens_fourway_test.go` 一字不动；票 201 的写面（`internal/agent/approval`＋`internal/tools`＋`cmd/wisp`）**未空出之前不许派本票**（同文件互斥）。
