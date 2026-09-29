# 197 — 界面要看得见每个子代理的状态，点进去能看到它**各自的流式工作页**；⚠ 现量：Go 侧今天**一枚"子代理"实体都没有**（不是缺视图，是缺那一层）

- Status: **在派（批准已落，见下面那一节）**。⚠ **功能票，owner 09-28 18:0x 原话逐字**：「**子代理必须看到状态，而且点击某个子代理，能看到它们各自的流式工作页面，能明白吗？这是主流 harness 必做的，不要偷懒**」⇒ 这一句同时是 `Q-71` 的裁定＝**乙（子代理上界面）**，票 188 的 `AC#1` 已按此勾。
- 现量（锚 `fccaf3e3`，尺可复制）：
  | 事实 | 读数 | 尺 |
  |---|---|---|
  | Go 侧"子代理"实体 | **0 命中**（非测试码，`internal/`＋`cmd/`＋`tools/`） | `grep -rln "subagent\|Subagent\|SubAgent" --include=*.go internal/ cmd/ tools/ \| grep -v _test.go \| wc -l` |
  | 中文词"子代理"在 Go 侧 | **0 命中** | `grep -rln "子代理" --include=*.go internal/ cmd/ \| grep -v _test.go \| wc -l` |
  | 多任务名册**已在** | `TaskRoster.byTask[taskID]`（`Record:123`／`Look:139`／`Count:152`） | `grep -n "func (r \*TaskRoster)" internal/tools/task.go` |
  | 流式载体**已在** | `StreamLog.Append(key,text)`／`Chunks()`（`pump.go:301`／`:339`），**键上限 32、溢出是"合并而不是丢"** | `grep -n "DefaultStreamKeys\|maxKeys" internal/panel/pump.go`（`:275`、`:287-290`） |
  | 面板今天能读到的 | `Snapshot` **4 枚键**；`PumpSources` 已接 `Results: rt.stream.Chunks` | `sed -n '57,62p' internal/panel/composer.go`；`sed -n '442,448p' cmd/wisp/run.go` |
  | ⛔ **撞上的既有定案（两处，都要人工批准才能改判）** | ① `PLAN.md:1544` 行＝**REJECTED「多 Agent 协作」——草案 6.6 明确不做——"不是债务"——单 Agent，无编排**；② `PLAN.md:1447`＝**S7 之前不得开放多任务并发（`TaskScheduler` 只允许单任务）** | `sed -n '1544p' docs/PLAN.md`；`sed -n '1447p' docs/PLAN.md` |

## 批准已落（09-28 18:1x，账 `A399`）——这一节原来是"等他一句话"，他当场给了，而且比我要的那句更宽
- **他 18:1x 原话逐字**：「**子代理这么重要的东西，说不要就不要了，还尼玛问问问，别人怎么做的，你就怎么做，别鸡吧天天问我，还想让我给你砍一刀，专门做简单了？？？必须特么做完整功能，明白吗？**」
  ＋「**从现在开始，别鸡吧看那个初始方案了**……别看初始方案，或者什么这个那个的门禁，把什么都自己主动舍弃了」。
- ⇒ **改判的两枚既有定案**（射程写死）：`PLAN.md:1544` 的 **REJECTED「多 Agent 协作——单 Agent，无编排」**，与 `PLAN.md:1447` 的 **「S7 之前不得开放多任务并发」**——**由这句当场改判，本票不再等批准**。
- ⚠ **`PLAN.md` 与 `docs/specs/**` 的文字仍一字不动**：他放权的是"**别拿它们当不做的理由**"，不是"把它们改掉"——改判**只落 `A399` 这一枚账**。要动冻结文字须另有一句批准（`AGENTS.md` §1.1 他没撤）。
- ⚠ **本票不许交"简单版"**：三层（实体／流／载体）都要落地，界面那层由 owner 带给界面那支（`A383` 定案，本编队 `frontend/**` 零写面）。**做不全就具名写清"三之一没做、挂在哪一格"，不许默默收窄**——这是 `A399` 里我认的那笔账的正面写法。

## AC#0 的裁定（不再摆他，理由落在冻结文本上）
- **落在 S5（Agent 内核）为家、S7（并发）为闸**：派生入口是"模型调一枚工具"＝`D34` 工具表与 `SPEC-05` 的射程；"同时能跑几枚"才是 `D31`／`SPEC-12` 的 S7 射程。
  ⚠ **这两节今天已经都不是单任务**：`TaskRoster`（`internal/tools/task.go:123/139/152`）与票 176 那根起跑口早已在树里跑多任务名册——**`:1447` 那行描述的是早已过去的一条时序，不是一条还在生效的禁令**。
- ⇒ 所以本票**不需要"先开放多任务并发"这一步**：名册在场，缺的是**"子代理"这一层实体**（现量 0 命中）。这一条按 `A382`/`A391` 那把尺分类＝**差一整层，不是差一根线**。

## 规格早就预见过这两个坑（直接当本票的实现约束，不用重新发明）
- `PLAN.md:1123`（DSH 对照行）：**子代理根本没有向人提问的通道**（审批通道缺失＝结构性取消问题）。⇒ 本票 AC 里必须钉：**子代理不许自带"允许"出口**；它遇到的 L2 只能**上抛给同一张面板队列**（面板只能拒绝／查看，`SPEC-06:19`），或者被父任务取消。**"面板侧来源的 L2 允许"这条铁律不许为子代理开口子**（`Q-49` 丙那批判据在场）。
- `PLAN.md:1127`：**"被阻塞的子代理在所有界面上都不可见"是别人的真实事故**（Issue #1723）。⇒ 本票 AC 里必须钉：**blocked／waiting-approval 这一态要在竖条上看得见**，不许只显示"在跑"。

## 要建什么（三层，按序；每层都要非实现者裁）
1. **实体层（Go 侧，全新）**：一枚"派生子代理作为任务"的能力——每个子代理拿到**自己的 `taskID`**、进 `TaskRoster` 名册（复用现成的 `Record`／`Look`），状态维取自 D43 那一集（`internal/statemachine/states.go:11-31`，与票 188／196 同一套词表纪律）。⚠ **不许新造第二套状态词**（票 196 就是防这个）。
2. **流层**：每个子代理一条自己的流式记录，键形状具名（建议 `subagent:<taskID>`）；⚠ **32 枚键上限与"溢出合并而不是丢"这一形对子代理是错的**——两个子代理的输出被合并成一条＝界面在说谎。⇒ 本票要**重裁溢出语义**（要么上限随子代理数走、要么溢出显式标"已截断"），并给一发**结构性判据**：造 N＋1 个子代理，**没有任何一条会被并进别人的**。
3. **界面层（归别家）**：竖条上每枚子代理一行＋点进去是它自己的工作页。⇒ **本编队只出载体与 Go 侧判据**，把该给的键集、形状、状态枚举逐字写进证据件，**由 owner 自己带给前端那枚 agent**（`frontend/**` 零写面、不读不引；见 `A383` 的定案）。

## 验收判据（草，逐格要 `file:line` 与正控）
- [ ] **AC#0 先答"这一层归哪个切片"**：`PLAN.md:1544`／`:1447` 两行改判之后，子代理派生落在 S 几（S3？S5 Agent 内核？S7 并发？）——**落在哪一节决定它要不要现在就能开多任务**；答不出就摆 owner，不许我自己挑一节塞进去。
- [ ] **AC#1 实体层有生产调用者**：派生入口在非测试码里有真听众（⚠ 今天"零调用方有两种成因：差一根线／差一整层"＝`A382`/`A391`，本票起手就要分清是哪一种）。
- [ ] **AC#2 名册与状态同源**：子代理状态**只能**是 D43 那一集；判据要能区分"填了真状态"与"字段恒空"（票 181 `AC#7` 那枚空转教训）。
- [ ] **AC#3 每子代理一条流、不许合并**：上面第 2 层那发结构判据。
- [x] **AC#4 阻塞态可见**：`blocked`／`waiting-approval` 在载体里有一格，且界面拿得到（`PLAN.md:1127` 那起事故的对应物）。
  **09-28 22:4x 翻勾（账 `A416`）**：生产路径判据 `TestRunPacketMarksTheRosterRowACardIsHolding`（`cmd/wisp/subagent_blocked_197_test.go`，commit `237e64f4`，**产码零改动**），
  我整包复跑到终态 `ok 109.672s`、零 FAIL；五发单点变异逐字在 `docs/evidence/s1/197-blocked-carrier-r4.md`。
  ⚠ **只勾"有一格且拿得到"这一半**：今天这枚旗**只报 L2**（L1 短窗口里等的孩子在名册上恒读 `false`，因为 `LiveApprovals` 只遍历 `q.pending`）——**那一半留在 `A416` 待裁，不算本格**。
- [ ] **AC#5 子代理不自带"允许"出口**：与 `Q-49` 丙那批判据同族；正控＝造一枚"子代理自己批自己"的假腿 ⇒ 要红。
- [ ] **AC#6 取消语义**：父任务取消时子代理怎么收（`TaskRoster` 今天有没有 owner／recover 那条链要现读）；⚠ 裸 `go func(` 而无 owner／recover 是 `AGENTS.md` §1.2 硬禁。

> ⚠ **09-29 10:5x 现读普查交回（只读腿 `197-c2`，证据件 `.scratch/wisp/probes/197/r5c/census.md`，13011 字节；⚠ **该件由编排者代落**——派单把只读腿派成了**没有写工具的类型**，错在我；件里每一行都由编排者自己在现树重跑过尺，标了〔腿报，未复核〕的格子不许当判据用）。三格定案：
> ① **AC#5 今天＝〔仅文档〕**：孩子结构上拿不到"答复"那一面（递给工具的 `Gate` 接口只有 `PendingWindow`／`PendingApproval`，`gate.go:99-106`；`grant` 只交给 `approval.UI.Prompt`，宿主 `consoleApprovalUI` 连印都不印，`run.go:1012-1035`），**但是"造一枚自批假腿 ⇒ 要红"今天没有任何尺会红**——`grant` 绑的是**事**不是**人**（`bindDigest(corr, taskID, tool, level, seq, args)`＝`approval.go:253`），而 `Request.Source` 逐字被注释写着「logged and **never consulted**」（`approval.go:231`，`gate.go:611`／`:618` 只进日志、`claimed_source=%q` 这个写法本身就写着"未核实的自称"）；已有的 `Test197SubagentHasNoSelfApprovalOutlet`（`subagent_197_test.go:838`）射程只有 `AdmitTask == nil` 拒绝＋字段枚数两格。
>   ⇒ **本票 AC#5 不许勾**，直到那枚正控存在。最便宜的合法接缝＝`cmd/wisp` 的 `package main` 测试（只有那一处同时握着 concrete `*Gate` 与 spawn 那条路）：记下 `grant` 后断言 `Native().Allow(childCorr, capturedGrant)` **必须失败**。⚠ 派这枚腿前**必做撞钉预检**——它会是 `Native()`（`gate.go:572`）**第一个调用者**（今天生产零消费者，账 `A418` ③），要先看票 145／146 那族名册差集尺会不会被顶红。
> ② **AC#6 拆成两枚具名票，本票不再自留这两格**：**票 220**＝名册那一格只看得见 L2、且同任务第二张卡静默掉格（`subagent_roster_197.go:190-210` 用 corr 建表却用 taskID 查、`pending_read.go:106-121` 只走 `q.pending`、L1 住在 `Gate.windows` 且**全仓无枚举口**、`queue.go:152-154` 把第二张卡改名成 `id#序号`）；**票 221**＝`subagent_197.go:189` 当着模型许诺「可以用 `task.cancel` 单独停它」而 `task.cancel` 未注册（`task.go:595` 只回 `task.output`，`:20-27` 逐字标 DEFERRED）、`TaskRoster.Cancel`（`task.go:442`）生产零调用者。
>   ⇒ **"父取消不级联、且写给模型看"这一半本票已做到**（`subagent_197.go:327` 的 `context.WithoutCancel` ＋ `:360-365` 结果文本含「没有被级联取消」）；**没做到的是"孩子有停止出口"**。
> ③ **载体那格的更正**：`blockedOnApproval`（`subagent_roster_197.go:127`）今天是一枚**布尔、不是 D43 的态**；名册上在跑的态只有 `subagent_197.go:99-104` 那四枚。**本票不许为"能画阻塞"而新造态名**（D43 转移表冻结），要动先落批准记录。
> ⛔ 另记：原腿交回的内容里有**一行来自本编队禁读目录**，编排者已整行丢弃、不转述（证据件末节记着这件事与"派单只禁了读、没禁引"这一处模板缺陷）。

## 禁区
`frontend/**`／`design/**` 零写面（连内容都不转述）；`PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`／C17 白名单既有名字一字不动（**改判只落 `A##` 台账，不动冻结文字，除非另有一句批准**）；**不新增方法名**（名册补齐是票 194 的射程，别在这票里顺手加）；`internal/risk/**` 不许绕过（子代理的落盘读回同样走 C25 盖戳那套，票 175／176／183 一族）。

## 要 owner 带给界面那支的一跳（09-28 21:5x，`197-r3b` 落地之后；Go 侧改不了，也不是待拍板项）

- ⚠ **本票文件名尾巴那句"go-side-has-zero-subagent-entity"从今天起是历史状态**：实体层（`b9fa815b`）＋载体层（`335b8d2b`）已入库。
  文件名**不动**（`-done` 后缀才是防重领的唯一键），**引用时以这一节为准**。
- **现象**：`internal/panel` 的 `TestApprovalCardViewJSONKeysMatchFrontendTypes` 今天红因**是两枚键**，红句逐字＝
  `Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`
  （`instructions` 那枚是票 200 带来的，见 `200` 票面同节；**`tasks` 是这票带来的**）。
- **转绿只有一跳**：界面那侧 `frontend/src/lib/panel.ts` 的 `PanelSnapshot` 补 `tasks?:` 一枚**对象**（Go 侧 `internal/panel/composer.go:91` `Tasks *TaskRosterSection`，带 `omitempty`）。
  **不许**为过尺让 Go 少发一个键、也不许改尺或加豁免——那把尺判的就是"加字段必两侧同批移动"。
- **`tasks` 这一节的字段（逐字取自 `internal/panel/subagent_roster_197.go`，21:5x 现读）**：
  节上＝`rows`（数组，**总是发**、可空）· `inFlightSlots` number · `poolCap` number（今天＝4，等于桥的工具并发天花板，见票 211）·
  `streamTruncated` boolean · `streamElidedRunes` number · `droppedStreamKeys` string 数组（**总是发**、可空）；
  `rows` 每枚＝`taskId` · `label` · `kind`（`root`／`subagent`）· `parentTaskId`（根为空串）·
  `status`（**只用 D43 那 20 枚名字**，`statemachine.Valid` 是判据）· `statusKnown` boolean · `statusReason?` ·
  `streamKey`（＝`subagent:<taskId>`；**同一包 `results` 节里 `correlationId === streamKey` 的那一条就是"点进去那一页"的正文**）·
  `blockedOnApproval` · `streamTruncated` · `streamElidedRunes` · `streamDropped`。
- **画的时候两件别做错**：① `statusKnown=false` 要显示"状态还不知道，因为……"（`statusReason` 就是那句原因），**不许把空 `status` 当成"没在跑"**；
  ② `streamTruncated`／`droppedStreamKeys` 是流被截断／被丢的实话——要么显示"这里少了 N 段"，要么明说没显示，**不许静默**。
- **入向那一跳今天不需要**（写腿现读结论，我复核过字节）：名册行自带 `streamKey`、正文走现成 `results` 节 ⇒ **没有新增任何 C17 方法名**。
  只有界面那支坚持要"宿主记住当前这格装的是哪枚会话"时才需要一枚入向方法——**那才是契约变更，要另落 `A##` 再动**。
- **撤销口令**：界面那支若决定这版先不画名册，Go 侧把 `Tasks` 那枚键退回不 marshal（**同批那两枚判据钉一起退**），"子代理上界面"那一格随之退回并具名记账——**不许留"键在、值恒缺"**。
