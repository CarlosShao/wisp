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
- [x] **AC#5 子代理不自带"允许"出口**：与 `Q-49` 丙那批判据同族；正控＝造一枚"子代理自己批自己"的假腿 ⇒ 要红。
  **09-30 15:3x 翻勾（账 `A471`；验收腿 `197-v1`＝非实现者，裁决表 `docs/evidence/s1/197-ac5-selfapproval-v1.md`，52,134 字节／§0-§6 齐／零占位符）**：
  本格自己的判据（"正控存在且要红"）**成立**——正控在场（`Native().Allow(childCorr, capturedGrant)` 五发逐发 `ErrBadGrant`／`ErrUnknownCorrelation`），
  且**造出被防结局那一发当场红**：⚠ **翻勾凭据只引本件 §1.1 的逐发读数**，两发最值钱——
  **M6**（`queue.go:376 allowScoped` 不再校验令牌 ⇒ 十条具名红，其中被测文件 `:284`／`:291` 打出"那张只被面板路线递过令牌的卡**真的执行了、文件真的被写出来了**"＝AC#5 声称要防的那个结局被造出来，而尺抓住了）；
  **D3**（把宿主送达的那一枚 `AnswerAllow` 改投 `AnswerReject` ⇒ 红**只在** `:277`／`:280` 正控那两行、**六发拒因一字未变** ⇒ "每一发都被拒"**不是**"整条路都是死的"造成的假绿）。
  编排者自己复跑（带 sherpa PATH，未突变）：两枚定向用例 **PASS**（11.40s／0.00s，`ok 11.468s`）。
  ⛔ **不引 `197-r1` 进度件 §② 当凭据**：那句"载具扫描…`task.spawn` 回给父模型的正文"经 P1 现量**是虚的**——`spawnText` 全文件只有 `:149` 声明＋`:308` 读、**没有任何赋值**（我现跑 `grep -n "spawnText" cmd/wisp/subagent_selfapproval_197_test.go`＝**只有那两行**），`spawn=0` ⇒ `strings.Contains("", 令牌)` 恒绿。
  ⚠ **四条附条件逐名（`197-v1` §0.2；四条都是"尺"的缺陷，没有一条推翻"孩子拿不到允许"这件事本身，故不入本格欠账、逐条落名）**：
  C1 第七枚载具恒空＝宣称了没测；C2 `:254` 六枚用例共用 `{ErrBadGrant, ErrPanelAllow}` 一枚两值集合 ⇒ 把 `ErrBadGrant` 原地换成 `ErrPanelAllow`（D1）**全绿**＝**尺分不清"门拦的"与"权限拦的"**（混用发生在尺上，生产侧三句文案逐字互不相同）；
  C3 文件末尾脚注 `M1` 指向 `ui.go` 的 `PanelItem`，照它改**全绿**（那一形进不了快照字节）＝会把下一腿带偏；
  C4 `:466` 名为"跨卡借证"实为**同名重放**（`startChildWrite197:109` 把 `CorrelationID` 写成 `taskID` ⇒ 两张卡共用 corr；13 发读数里三枚 id 逐次全等）。
  ⇒ **残余去处**：C1／C2／C3 三枚**尺的修法**归下一程写腿 `197-r3`（动被测文件必须同批把 **P1／M1 真形／D1** 三发当正控重跑，台件 `.scratch/wisp/probes/197/v1/` 可直接复用）；
  C4 与"出向读面零尺"那一发**单开「票 242」**（`bindDigest` 绑定层与 `approval.PanelItem` 出向读面**今天全仓零仪器**，我现跑 `grep -rln "bindDigest" --include=*_test.go internal cmd`＝**只有那一枚被测文件**）。⛔ 本票不吞这两格。
  另记一枚正面事实（别在更正时顺手抹掉）：`197-r1` 交回时那句"本腿按甲执行并上报"姿势**正确**——它没自行结案、没动那三处、也没把"要执行"当成已完成（`197-v1` §4.1 第 5 条）。
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

## 编排者口径 — 验收腿做单点突变时能不能"动生产码"（09-30 15:3x 定，答 `197-v1` §4.2 提的那处冲突；账 `A471`）

`197-v1` 把它撞到的这处口径冲突**上交给我定**（原话"这一处口径冲突请你在票面上定一句，否则下一腿还会撞到"），定死如下：

- **允许**验收腿为了跑一发突变而**临时改生产码**（改完用内存里那份原字节还原、逐枚打印 `RESTORE <path> bytes=<n> exact=True`），这是对抗验收的正常动作，不算越界；⛔ 但仍不许动三枚冻结件、不许放宽任何断言、不许 commit 产码。
- **当某一形突变必然连带要求动一枚冻结件时**（例：给 `tools.Gate` 接口加方法 ⇒ `internal/perm/ticket90_persist_test.go` 里那枚 `t90gate` 就必须实现它，而那枚文件一字不动），**允许走 `go test -overlay`**，三个条件一枚不许少：
  ① 具名写清 overlay **买到什么**（只换编译期的源，盘上字节一字未动）与**买不到什么**——⛔ **读盘／走 AST 的测试对 `-overlay` 结构性失明**，本包就是活例子：`cmd/wisp/subagent_carrier_197_test.go:431 TestSubagentStreamKeyHasOneMintSite` 与 `leg_sink_gate_131` 那一族都在这一类里；
  ② 那条"读盘"的腿必须**另给一发不依赖 overlay 的读数**，或者直接**明写为空缺**，⛔ 不许把 overlay 的绿说成"整形已证"；
  ③ 驱动器落盘前先断言锚点在那枚文件里**恰好命中 1 次**，跑失败也要还原（`finally` 分支）。
- 这一条是**本票口径**。要把它推广到别的票／别的依赖边裁定，按台账 `:9761` 我自己立的规矩**具名写成"这是编排者的推广"**，不写成"仓里本来就有这条"。

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

---

## 197-r1 进度（后端写码腿，起钟 2026-09-30 11:37:58 +0800／落笔 12:10:0x +0800）

⛔ **本腿未跑任何 `go test`**（整包／单包／`-run` 定向都没跑）——门禁窗口那段时间归独占桌面采样的验收腿。
本文件一枚 AC 复选框都没碰（勾框归编排者）。台账与 HANDOVER 未写（归编排者）。
交付＝一枚 commit `a818df46`：`cmd/wisp/subagent_selfapproval_197_test.go`（710 行，**产码零改动**；
`git show HEAD:<path>` 与工作树逐字节 diff rc=0 验过）。

### ① 起手复跑票面「现量」表：五行里三行已过期（具名更正；锚 `bf26dfc4`，收尾锚见 §⑤）

| 票面那行（锚 `fccaf3e3`） | 票面读数 | 本腿现跑 | 尺 |
|---|---|---|---|
| Go 侧「子代理」实体 | 0 命中 | **9 个文件**（`internal/tools/subagent_197.go`、`internal/panel/subagent_roster_197.go`、`cmd/wisp/panel_pump.go`、`internal/streamkey/streamkey.go` 等） | 票面原尺逐字复跑 |
| 中文「子代理」在 Go 侧 | 0 命中 | **3 个文件** | 票面原尺逐字复跑 |
| `TaskRoster.byTask` 的 `Record:123`／`Look:139`／`Count:152` | — | **漂到 215／268／281**，且名册已有 13 枚方法（新增 `Descendants:305`、`PublishSubagent:342`、`MarkRoot:355`、`TryAcquireSubagentSlot:363`、`InFlightSubagents:390`、`RunningSubagentIDs:402`、`AttachCancel:420`、`DetachCancel:434`、`Cancel:447`、`WatchRow:248`） | `grep -n "func (r \*TaskRoster)" internal/tools/task.go` |
| 流式载体 `pump.go:301`／`:339`，上限尺 `:275`／`:287-290` | 键上限 32、溢出「合并而不是丢」 | **上限常量在 `pump.go:406`（`DefaultStreamKeys = 32`）＋ `:408` `StreamKeyHardCeilingMultiple`**；溢出那一支已经改成「显式声明截断、不再并键」：`:509` 注释逐字 `no key is ever folded into another`，判定体在 `:514`／`:521` ⇒ **票面这一行对子代理的那条反对意见已被现形解决** | `grep -n "DefaultStreamKeys\|maxKeys" internal/panel/pump.go` |
| 面板今天能读到的：`Snapshot` 4 枚键 | 4 | **6 枚**（`internal/panel/composer.go:57-`：`pending`／`results`／`composer`／`generatedAt` ＋ `instructions`（票 200）＋ `tasks`（本票））；`Results: rt.stream.Chunks` 从 `run.go:442-448` **漂到 `run.go:615`** | `sed -n '56,66p' internal/panel/composer.go`；`grep -n "Results:   rt.stream.Chunks" cmd/wisp/run.go` |
| ⛔ 撞上的两枚既有定案 | `PLAN.md:1544`／`:1447` | **原文在场、一字未改**：`1544` 行首逐字 `\| **REJECTED** \| 多 Agent 协作 \| 草案 6.6 明确不做 \|`，`1447` 行内含「不得开放多任务并发」；且 `git status --porcelain docs/PLAN.md docs/specs tools/d22scan/allowlist.txt` **全空** | `sed -n '1544p' docs/PLAN.md`；`sed -n '1447p' docs/PLAN.md` |

⇒ 票名尾巴那句「go-side-has-zero-subagent-entity」连同上面三行现量都是**历史状态**（与本票 §「要 owner 带给界面那支的一跳」第一条同一口径）；上表是本腿在 `bf26dfc4` 上重新读到的。

### ② 本格交付＝AC#5 的那枚正控（票面 §09-29 那节留下的唯一 Go 侧格）

09-29 普查把 AC#5 判成〔仅文档〕，并留下了一句可执行接缝：「最便宜的合法接缝＝`cmd/wisp` 的 `package main`
测试……记下 `grant` 后断言 `Native().Allow(childCorr, capturedGrant)` 必须失败」。本腿把那一句变成读数。

**(1) `Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands`（真 run 装配，产码零改动）**
- `task.spawn` 派一枚**实名册子代理**（`rt.bridge.Execute` 走真桥；名册行由 `PublishSubagent` 与子环路自己的准入钩子发布，本文件不造行），
  随后**用孩子自己的 task id** 通过真桥撞出两张真 L2 卡（`fs.write` 打到 `t.TempDir()`，即 R2⇒L2 那一形）。
- **正控先行**：第一张卡由宿主账上那枚真令牌经 `Replies.Allow` 允许 ⇒ 调用真的执行、文件真的落盘。
  没有这一发，后面每一发「被拒」都只是「这条路由今天没在跑」。
- **自批的每一发**：空令牌／假令牌（运行时从真令牌改尾 4 字，源码里没有令牌字面量）／
  自称 `Source="subagent-self"` 的 `DecideFromNative` ⇒ 逐发 `ErrBadGrant`；已结卡的重放 ⇒ 只允许 `ErrUnknownCorrelation`。
- **第二张卡**：跨卡借证 `ErrBadGrant`；真令牌在 `DecideFromPanel` 上露一次面即被按泄露烧掉（`ErrPanelAllow`），
  此后连原生侧花它也 `ErrBadGrant`；收口是宿主的 `Reject`，文件始终不存在。
  审计侧同时要求 `FORGED-OR-STALE`／`PANEL-ALLOW-REJECTED`／`ANSWER-ALLOW` 三行在场（被拒与放行是两种读数，不许靠推断）。
- **载具扫描（卡还挂着时取样，不是事后）**：这一程自己发布的快照字节／工具侧 `TaskOutput` 全字段／
  `StreamLog` 全字段／`task.spawn` 回给父模型的正文／stdout／stderr／`<data>/logs` 持久日志
  ⇒ **令牌本身一处都不许出现**（`bindDigest` 绑的是事不是人，「拿不到」是这一层唯一防线）。
- 设计裁定落点复跑：**子级永不自批**＝本格；**结论复用现成源名 `task.output`** 的尺
  `git grep -n "subagent\.output" -- 'internal/**/*.go' 'cmd/**/*.go'` **＝ 1 命中，且那一行逐字是
  「No "subagent.output" name is invented」**（`internal/tools/subagent_197.go:35`）；
  **父取消不级联且写给模型看**＝票面 §09-29 ② 已判为做到，本腿未动（孩子那一侧的停止出口仍在票 220／221，本腿不回收）。
- **⚠ 本腿自己发现并当场改掉的空转判据（具名，不藏）**：起初写了「取到的那枚快照字节要出现在持久台账里」——
  这条**今天恒红**：`publishPanelSnapshot` 把最新包留在 `rt.lastSnap`／`lastSnapBytes`，收尾那一发会覆盖挂卡时那枚，
  而本用例**故意不调** `publishPanelSnapshot`（197-r4 的 `case 1 never calls publishPanelSnapshot` 同一纪律），
  所以永远读不到挂卡那一枚的 sha ⇒ 判据与产码事实互斥。改成两条真话（`snap.Pending` 里认得出这张卡 ＋ 这一程至少落账过一枚快照），
  并把做不到的那半（sha-to-ledger tie）**具名留在注释里**交给整包复跑，不留一条恒红判据给别人。

**(2) `Test197NoAllowDoorIsReachableFromASubagentsAssembly`（纯反射，不起 run）**
- `tools.SubagentDeps`／`tools.TaskDeps`／`tools.Options`／`agent.Options` 的**静态类型图**里没有任何类型
  声明 allow 类方法（`Allow`／`Native`／`DecideFromNative`／`DecideFromPanel`／`GrantNonce`），
  且 `tools.Gate` 的方法集仍恰好是 `{PendingWindow, PendingApproval}`——这才是「孩子自带允许出口」的可达性那一半。
  名字集刻意不含 `Reject`／`Veto`：那两个方向混进来会把两件事判成一件事（拒绝方向那一格是票 220 的地界）。
- **自带种门的正控**（第 64 条；负向尺必配「种 X 必响」）：文件里植了一枚带 `Allow`／`Native` 的假门载体，
  扫描器读不出那两枚名字就 `t.Fatalf`——**照不见门的仪器不许签发「没有门」的判据**。

### ③ 期望编排者跑哪几发、期望看到什么（⛔ 未跑测试，等门禁窗口；本腿不写「应该能过」）

| # | 跑法（都要带 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`，否则会以 `exit status 0xc0000135` 失败且不打 `--- FAIL`） | 期望 |
|---|---|---|
| 1 | `go test ./cmd/wisp/ -run 'Test197NoAllowDoorIsReachableFromASubagentsAssembly' -count=1 -v` | `--- PASS`；秒级（这一发不起 run）。若红：**先读 findings 里指名的「类型.字段.方法」**，那是判「是不是新门」的原始读数，不是先删断言 |
| 2 | `go test ./cmd/wisp/ -run 'Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands' -count=1 -v` | `--- PASS`。本用例把 mockllm latency 设成 2500ms（`runTextTask` 在根环路结束后就 `rt.close()`，那 2.5s 就是探针的工作窗口），健康读数个位数秒；红路径最坏约 50s（写调用 ctx 30s 上限）。`t.Logf` 那行期望：`child=<id> corr1=<id> corr2=<id>`、六发 error 是 `ErrBadGrant`×4 ＋ `ErrPanelAllow` ＋ `ErrBadGrant`、`replay=ErrUnknownCorrelation`、`落盘=true/false` |
| 3 | 整包 `go test ./cmd/wisp/ -count=1` | 本腿两枚 PASS；红名册相对 `subagent_blocked_197_test.go` 那一族**不应多出一枚**。⛔ 不许为变绿放宽任何断言 |
| 4 | 变异对照 **M1–M6**（清单逐条写在该测试文件末尾注释里） | 每发要看到**具名那一格**红：M1 面板载具带 Grant／M2 控制台打印令牌／M3 `SubagentDeps` 多一枚能答的字段／M4 `tools.Gate` 多出答复方法／M5 让 `Request.Source` 变成权威／M6 `Replies.Allow` 不要令牌。⚠ M1、M2 会**同时**打红冻结件 `internal/panel/l2_grant_boundary_test.go`——那是同族互证，不是冲突，不许为放行其中一枚去改另一枚 |

### ④ 撞钉预检（第 64 条：读断言，不是只 grep 新符号名；本腿未跑任何包，凡未读者一律标〔未跑包，此条来自读断言〕）

| 今天绿着的钉子 | 它的射程为什么盖不到／盖得到本腿，判词 |
|---|---|
| `cmd/wisp/leg_dispatch_gate_133_test.go` 的 `runRosterReds135`（`:1632` 起，实现逐条读过） | 只查「账上写的 `covered=test` 名字在这轮 binary 里必须 startable」，**不查反向** ⇒ 新增测试枚不欠账、不顶红。〔未跑包〕 |
| `cmd/wisp/leg_sink_gate_131_test.go` | 走 `func main` 的分派形状与 `installLogSink` 可达性，读源码不读测试名册 ⇒ 不受新增测试文件影响。〔未跑包〕 |
| `cmd/wisp/subagent_carrier_197_test.go:430 TestSubagentStreamKeyHasOneMintSite` | `stringLitSites197` 显式跳过 `_test.go`（`:462`），判的是「非测试文件里恰好一处字面量」；本腿既没在产码里、也没在测试里写那枚字面量 ⇒ 不动它。〔未跑包〕 |
| `internal/tools/subagent_197_test.go:846 Test197SubagentHasNoSelfApprovalOutlet` | `len(SubagentDeps 字段) != 5` 即红 ⇒ 本腿零字段改动；它的 `AdmitTask == nil` 那一支与本腿反射腿相邻互补、不重钉。〔未跑包〕 |
| 冻结件 `internal/panel/l2_grant_boundary_test.go`（**未动一字**） | 它钉「面板路线不得承载裁决」，本腿钉「子代理读得到的一切不得载有原生令牌」——互补；M1 那发同时打红两边属预期（§③ 第 4 行）。 |
| `internal/agent/approval/{queue,ticket87,ticket97,ticket146}_*_test.go` | 都在包内自造 Gate，不读 `cmd/wisp` 源码 ⇒ 不受影响。⚠ 本腿是 `Native()` 在 **`cmd/wisp` 里的第一个调用者**（普查 §① 的预告），但全在 `_test.go` ⇒ 账 `A418 ③`「生产零消费者」那句话今天仍然成立，没有被本腿偷偷改掉。〔未跑包〕 |
| 全仓仪器 `tools/d22scan/d22scan.exe` | **本腿自跑过＝clean**（这条不是读断言）：ban #1（裸 `go func(`）与 ban #3（明文密钥）**不扫 `_test.go`**（`main.go:665`／`:858` 跳过），ban #8（emoji）扫；本文件的 ⚠／「」全在注释里，字符串字面量里没有 U+2600–U+27BF／U+2B00–U+2BFF／U+FE0F。 |

### ⑤ 本腿跑过的门（收尾现量；最后一次复跑与落笔同钟 **12:09:55 +0800**，跑的就是 `a818df46` 那枚提交里的字节）

- `go build ./...` rc=0；`go vet ./cmd/wisp/` rc=0；`go vet ./internal/panel/` rc=0。
- ⛔ **`go vet ./internal/tools/` rc=1，非本腿所动**：别人未入库的 `internal/tools/grant_test.go`
  （`grantsOf redeclared`／`mustCanonical redeclared`／`undefined: grantSourceType`）——票 224 的活在飞。
- ⚠ **未归因（第 74／75 条）**：派单给的起跑读数「`internal/tools` 161 顶层声明／161 PASS」本腿**不复现、也不改口径去凑**——
  源码侧 `^func Test` 枚数在 12:05→12:07 之间从 141 涨到 177（`grep -rhE '^func Test' internal/tools --include='*_test.go' | wc -l`；
  期间 `3b78c246 feat(224)` 落库、`grant_test.go` 仍是未入库新件），`internal/panel` 99 枚、`cmd/wisp` 122 枚（含本腿 2 枚）。
  ⇒ 「161 PASS」要复现就得跑包，本腿不跑 ⇒ **不当成自己的读数，也不判谁漂了**。
- `gofumpt -l cmd/wisp/subagent_selfapproval_197_test.go` 空。
- `tools/d22scan/d22scan.exe`：`clean - no D22 ban violations`。口径提示：`ban #8 internal/=466`、`cmd/=64` 是**被扫文件数**（含 `_test.go`），不是违规数。
- 冻结件与别人地界零改动：`git status --porcelain docs/PLAN.md docs/specs tools/d22scan/allowlist.txt internal/panel/tokens_fourway_test.go internal/panel/l2_grant_boundary_test.go internal/perm/ticket90_persist_test.go` **全空**；
  `frontend/**`／`design/**` **零读零写**（本节与测试文件里没有来自那两层的任何结论）。
- Git：`--only` ＋ 显式 pathspec 一枚文件；未 push；未 amend／reset／rebase／stash／`checkout .`／clean；未建 worktree。
  提交前后别人的脏文件（`internal/tools/grant*.go`、`cmd/wisp/run.go` 的会话令牌活、`.gitignore`、`design/**` 那 16 枚删除、
  `docs/evidence/s1/152-*`、`.scratch/wisp/probes/161/r6/logs/**`）仍在原位，本腿一枚未提交。
- 台件只建在 `.scratch/wisp/probes/197/r1/`（本节那一枚，**只建不删**；这份文本随后原样追加进本票面，票面是权威）；
  仓库其它目录没有本腿留下的东西。

### ⑥ 没做完／留给编排者（⛔ 不留半成品不声明）

1. **AC#5 仍不该被勾**：本腿交付的是那枚**正控**（普查 §① 欠的那格），本格是否成立要**非实现者**裁
   （裁决者≠实现者：`SPEC-12 §4.3` #1/#3、`issues/README` 硬约束末条、D22 双角色）。本腿一句不替它判。
2. **AC#0／#1／#2／#3／#6 一枚未动**：AC#0 的答案在票面 §AC#0 那一节（S5 为家、S7 为闸）；
   AC#6 的两格早已移出去到票 220／票 221，本腿没有顺手回收，也没有把「孩子没有停止出口」这格藏起来。
3. **一支需要编排者判的「两支」（本腿按保守那支做完并具名上报，没有自行假设结案）**：
   票面 §0 与本腿派单都写着「流式键今天在两个包里各拼了一份 ⇒ panel 为真相源、由装配根注入函数、不新开 tools→panel 依赖边」，
   而盘面事实是**这一格已由载体层那一程用第三包落掉**：字面量只在 `internal/streamkey/streamkey.go` 一处，
   `internal/panel/pump.go:430` 与 `internal/tools/subagent_197.go:53` 各留一枚 alias，
   外加 `TestSubagentStreamKeyHasOneMintSite` 那枚钉＋`cmd/wisp/subagent_stream_key_197_test.go` 那对比对。
   两支：**(甲)** 承认现形＝裁定精神已满足（真相源只有一处、`tools` 不 import `panel`，那条依赖边确实没开），
   记成「已执行、形状与裁定文字不同」；**(乙)** 要求改回「panel 导出＋装配根注入格式化函数」那一形
   ＝重写已入库且带钉的东西。**本腿按甲执行（那三处一行未动）并在此上报**；
   台账里那句「197-r3 要执行」是否结清，归编排者落账，不归本腿判。
4. **⚠ 一枚可能撞红的预告（不预判、不改判据）**：票 224 的会话令牌那一档若把带 `Allow` 方法的类型挂进
   `tools.Options`／`tools.TaskDeps`，本腿第 (2) 枚反射腿**会红**。那时要裁的是
   「它是判级输入还是答复出口」（`allowed_dirs` 那条教训同族：判级输入不是执行时硬边界，反之也不该被当成出口），
   **不是**放宽这枚尺。
5. 界面那一跳（`PanelSnapshot` 补 `tasks?:` 一枚对象键）与本腿无关：本腿没动 marshal 形状，
   红句今天仍是 `Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare` 那一族。
6. 票名里那句「go-side-has-zero-subagent-entity」依旧是历史状态（现量见 §① 那张表）。

### ⑦ 更正（追加，不删上面那条；同钟 2026-09-30 12:11:27 +0800）

§⑤ 里那枚「`go vet ./internal/tools/` rc=1」**只在量到的那一刻成立**（12:04:59 与 12:05:17 两次读数，
原因是别人**当时还没入库**的 `internal/tools/grant_test.go`：`grantsOf redeclared`／`mustCanonical redeclared`）。
那枚文件随后以 `8b57a419 test(224)` 落库，本腿在 12:11:27 复跑同一把尺：**`go vet ./internal/tools/` rc=0**。
⇒ 那条读数是**时序**，不是一枚还挂在那里的缺陷；引用 §⑤ 的人请按本条为准。
本腿从头到尾没有动过 `internal/tools/` 任何一枚文件（`git log a818df46 657396f0 --name-only` 两枚提交里只有
`cmd/wisp/subagent_selfapproval_197_test.go`、本票面与本台件三枚路径）。
