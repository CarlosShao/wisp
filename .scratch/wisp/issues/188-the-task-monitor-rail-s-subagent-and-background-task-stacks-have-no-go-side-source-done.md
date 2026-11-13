# 188 — 面板"任务监控"栏里**子代理状态**与**后台任务状态**这两堆今天没有 Go 侧来源（owner 09-28 放权：缺就补）

- Status: **done（2026-11-13 13:0x；修「空正文终局整条丢行」，三包全绿）**。⚠ **功能票**。
- 来路：owner 09-28 15:2x 原话（逐字在台账 `A373`：「我现在给你放权，缺什么就给我补」）＋他圈了右上角那一栏的截图。普查出处＝`docs/evidence/s1/180-182-panel-fields-census-c1.md` §5／§7（`180-c1＋182-c1` 交件，编排者收）。
- 关联：票 182（那一栏的堆数与三档定性＝本票的上游）· 票 145（快照载体十六字段）· 票 163（常驻终端会话，**不同物**）· 票 33／35（面板宿主与桥，硬前置）

## 现量（`180-c1＋182-c1` 本程现跑，锚 `e9ef94d0`；编排者复跑同值）

| 维 | 今天有什么 | 缺什么 | 尺 |
|---|---|---|---|
| **子代理状态** | **规格真空**（不是"有规格没实现"） | 全仓没有"子代理"这一维的任何数据模型 | `grep -rniE "Subagent\|sub_agent\|SpawnAgent\|DelegateTask\|fanout" --include=*.go internal/ cmd/` ⇒ **0 命中** |
| **后台任务状态** | `internal/tools/task.go:110-152` 有 `Record`／`Look`／`Count`（任务名册） | **没有"状态"这一维**（成功／失败／进行中／被拒／取消都取不到） | 同上文件现读 |

⇒ 截图里那一栏今天画的"已完成 09:41／回复中 10:18／已拒绝 昨天 16:05"**在 Go 侧一枚都取不到**（前端是示意数据）。

## AC（每格都要答"这一发在未修码上响不响"）

- [x] **AC#1 先答"子代理这一维到底要不要进产品"**＝**乙（进，owner 09-28 18:0x 原话逐字：「子代理必须看到状态，而且点击某个子代理，能看到它们各自的流式工作页面，能明白吗？这是主流 harness 必做的，不要偷懒」）**；⇒ `Q-71` 已裁（账 `A397`）。⚠ 这一句牵出两行既有定案需要改判入档（`PLAN.md:1544` REJECTED"多 Agent 协作"＋`:1447`"S7 前不得开放多任务并发"），**实体层与流层归口＝新票 197**；本票余格只管**后台任务那一维**（`AC#2` 在飞＝`188-r2`）。⚠ 这一格不是技术题——`180-c1＋182-c1` 现量它是**规格真空**，加它＝新需求（牵 D22 双角色／D31 并发审批／D38 线程模型三节）。⇒ 只读普查：要动哪几节规格、有没有更小的等价维（例："这一轮有几个后台任务在跑"而不是"有几个子代理"）。**结论摆给 owner 拍，不许自己定。**
- [x] **AC#2 后台任务状态这一维先落地**（它不需要新需求，名册已在、缺字段）：给 `Record` 补状态维（枚举要**逐字对齐 D43 状态机权威转移表**，不许自造），并送进快照。⚠ 判据要钉"值来自真名册"，不许硬编码。
- [x] **AC#3 反向判据**：状态维**不许**变成"面板侧可写"（面板只许显示＋发起请求，R20／票 92 口径）；补一枚常驻判据钉住"这一维没有面板写腿"。
- [x] **AC#4 与票 145 分账**：快照载体扩张归票 145（`SPEC-08` 那一节的字段表），本票只交"源"；结题时逐条写明谁欠谁。
- [x] **AC#5 契约轴**：`docs/PLAN.md`（D31／D38／**D43 转移表**）、`docs/specs/**`、`allowlist.txt`、`thresholds.go`／golden 一字节不许动；**C17 白名单一枚不许加**。
- [x] **AC#6 门禁**：逐包 `go test -count=1 ./internal/tools/ ./internal/agent/ ./internal/panel/`；`sh scripts/d22scan.sh`（基线 `ban #8 internal/` **examined=433**）；`gate-clauses.sh` **比红腿名册**（在册只 `G6neg`）；⚠ `./internal/panel/` 今天有 **2 枚已知红**（`design/assets/tokens.css` 被别家会话删了未 staged），照实记不许当绿、不许顺手修。

## 本票**不**解决
不画界面（`frontend/**` 归另一枚会话）；不做"未提交枚数／diff／文件树"（＝票 189／190）；不做面板终端与浏览器（＝票 191／192，两枚都要人先批准）。

## Progress log
- 09-28 15:2x 编排者立票：来路＝owner 放权补票（`A373`）＋`180-c1＋182-c1` 名册里"没票认领"七处中的两堆。未派。
- 09-28 18:1x `188-r1`（写腿，只做 `AC#2`，派单 §E）：**`AC#2` 停手上报，零产码**；证据件 `docs/evidence/s1/188-task-state-r1.md`。三处前提冲突逐条带 `file:line`：
  ① **写面与票面 现量 不重合**——本票 现量 与 `AC#2` 点名的 `Record` 是 `*tools.TaskRoster` 的方法（`internal/tools/task.go:123`，票面行号段 `110-152` 对得上），缺字段的是 `TaskOutput`（`task.go:93-102`，只有 `Text:97`／`ArtifactPath:101`），而派单给的写面是 `internal/memory/**`，那包里**没有名为 `Record` 的载体**（只有 `dao_providerhealth.go:131 RecordProviderError`、`privacy.go:194 Records`）；
  ② **memory 侧不是"缺字段"**：`TaskLog.State` 已在（`internal/memory/models.go:62`；schema `state TEXT NOT NULL` 且刻意无 CHECK，`schema.go:57`／`:14-17`；DAO 只校验非空，`models.go:185-197`），生产者两枚都在写面外——`internal/agent/loop.go:983-998 taskLogState`（`done`/`stuck`/`cancelled`/`interrupted`/`error`，注释自称 free-text）与 `cmd/wisp/run.go:651-670`；这套词**一枚都不在 D43 的 20 个名字里**（那 20 枚的仓内逐字拷贝＝`internal/statemachine/states.go:11-31`，标题行 `docs/PLAN.md:3055`、40 条转移 `:3061-3100`，本程未动一字）；
  ③ **撞钉预检命中即红**：把 `task_log.state` 改成 D43 白名单会让既有断言变红——`internal/memory/dao_test.go:253,262,265,282,307`（用 `running`/`succeeded`/`interrupted`）、`internal/memory/privacy_test.go:28,34`、`internal/agent/forensics_test.go:128-132`（钉 `cancelled`）、`internal/agent/loop_golden_test.go:315-320`（钉 `done`）；两枚在写面外、两枚在写面内，一律不自放宽。
  需要的具名解冻＝编排者二选一：(i) 解冻 `internal/tools/task.go`＋`task_backfill.go`（只到"给 `TaskOutput` 新增状态字段"为止），或 (ii) 改判载体为 `internal/memory.TaskLog` 并同步解冻 `internal/agent/loop.go` 那枚映射与 `cmd/wisp/run.go` 调用点；另外"任务态 -> D43 名字"那张映射表是**规格真空**（D43 讲的是会话/球的 20 态，没有一行讲后台任务），要落 `AC#2` 先定它归谁。
  `AC#3` 反向判据也**未落**：现量 `internal/panel/` 非测试码里 `task_log`／`TaskLog`／`TaskOutput`／background task **零命中**，载体没定就只能交死钉（无正控）；等 (i)/(ii) 定案后同批落，并自带"种一发面板写腿必响"的正控。`AC#1` 子代理那一维**没碰**（`A389` 已摆 owner）。
  门禁读数：`d22scan` **clean**（`ban #8 internal/=441` 枚文件，票面基线 433 已过期＝别家加文件，涨不等于漂移）、`gate-clauses.sh` **BAD 名册只 `G6neg`**、`gofumpt -l internal/memory/` 空、`go test ./internal/memory/` ok 12.810s；`./internal/tools/`／`./internal/agent/`／`./internal/panel/`／`cmd/wisp` **未采读数**（在别人的飞件名下）。`flip-declaration.sh` 未跑。AC 框一枚未勾。
- 09-28 18:0x 编排者收 `188-r1`（五枚 commit，零产码、停手上报）：**载体选 (i)**＝解冻 `internal/tools/task.go`＋`task_backfill.go`，射程只到给 `TaskOutput` 加一枚状态字段，值**只能取自 `internal/statemachine/states.go:11-31`**（D43 逐字拷贝），不许新造、不许塞 `done/cancelled`；**不改判成 `TaskLog.State`**（那要动两处生产写入者＝改既有语义）。⚠ 它挖出的两套词表并存另立**票 196**；本票 `AC#2` 期间**不许与 196 互相顺手统一**。账 `A394`。
- 09-28 18:1x `188-r2`（写腿，派单 §G，只做 `AC#2`＋`AC#3`；commit `62d7d3b4`）：**产码了，两格都落**。逐条自证：
  - **`AC#2` 字段形状**＝`internal/tools/task.go:124` `State statemachine.State`。类型不是 `string`：`statemachine.State` 的取值集就是 D43 的仓内逐字拷贝 `internal/statemachine/states.go:11-31`（`:8` 自陈"names exactly as in D43 / SPEC-08"，出处 `docs/PLAN.md:3055`＋`:3061-3100`，**本程零字节未动**）。**没新造名字**：判据 `TestTaskState188AC2VocabularyIsD43Only`（`internal/tools/task_state_188_test.go:76`）枚数钉 20、逐枚过 `statemachine.Valid`（states.go:39）。**旧词表零渗入**：同尺把 `done/cancelled/running/succeeded/interrupted/error/finished/failed/stuck/queued/pending` 十一枚逐枚要求**被拒**；写侧门＝`task_backfill.go:113`（`statemachine.Valid` 不过 ⇒ 整条记录拒绝落表并带原因），正控＝`TestTaskState188AC2BackfillIsTheGoSideProducer:170` 的 (c) 段（逐枚喂旧词 ⇒ 必须红＋名册 `Count` 不变）。
  - **"值来自真名册，不许硬编码"**＝`TestTaskState188AC2ValueComesFromTheRoster:118`：同一张表里三条记录给出三个不同答案（`Acting`／`Error`／未登记），再 `Record` 换名 ⇒ 读数跟着表走。读侧＝`TaskOutput.StateAnswer`（`task.go:142`），空值与非法值各回一句 fail-closed 文案，**不改 `task.output` 给模型读的那句 stub**（那套措辞按 `task.go:251` 是 PLAN.md:2564／票 164 AC#3 冻结件，一字未动）。
  - **`AC#3` 反向判据**（载体已定，这格现在能落）＝`TestTaskState188AC3PanelHasNoWriteLeg:341`：AST 扫 `internal/panel/` 非测试源码（**分母 11 枚，<10 直接 `Fatalf`**）里"写这一维"的四种形状——调 `Record`／`Backfill`、构造带 `State` 键的 `TaskOutput`、给 `.State` 赋值、`Method*` 白名单常量声称 task/state 轴。现量 **0 命中**。正控在 `t.TempDir()` 里种（**仓内零新文件**）：五枚假腿逐枚必须红（`roster-record-call`／`roster-backfill-call`／`taskoutput-state-key`／`state-field-assign`／`whitelist-method-claim`）；**反向控**＝种一枚"面板只读 `Look`／`Count`"的读腿必须**不红**，否则这扇门会把票 145 的泵读口打进自己的红名册。**未动任何人的钉**：`internal/memory/dao_test.go:253,262,265,282,307`、`privacy_test.go:28,34`、`internal/agent/forensics_test.go:128-132`、`loop_golden_test.go:315-320` 全部 `git diff` 为空，`./internal/memory/` 与 `./internal/agent/` 复跑皆 ok。
  - **门禁读数**：`go test -count=1 ./internal/tools/ ./internal/memory/` ＝ **ok 17.675s / ok 13.319s**（gofumpt 后复跑）；`./internal/agent/` ok 1.784s；`go build ./cmd/wisp/ ./internal/agent/` rc=0（生产写入者未被本程改动仍编得过）；`gofumpt -l internal/tools/ internal/memory/` 空；`sh scripts/d22scan.sh` **clean**，`ban #8 internal/=442`（上一程 441 ⇒ 增量＝本程判据一枚）。`internal/panel/` 三枚在册红**逐名照实记、未修未当绿**：`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`。
  - **没测到／做不到（具名）**：① **这枚字段今天生产写入者仍是零**——`cmd/wisp/run.go:640-647` 构造 `TaskBackfill` 时没带 `State`（本程禁写 `run.go`，§F 正在里面接读口），所以生产路径的记录一律 `State == ""`，读数走"未登记"文案；把它接成实值是一行装配，不归本票。② **送进快照那一半没做**＝票 145 的载体射程（`AC#4` 分账：本票交"源"，145 交"载体"）。③ **`task.output` 的模型可读文本里没有这一维**（stub 措辞冻结）。④ `gate-clauses` 那把尺不在 `scripts/`（`sh scripts/gate-clauses.sh` ⇒ No such file or directory），派单未给路径；本程在 `.scratch/wisp/probes/154/gate-clauses.sh` 找到并跑了：**BAD 名册仍只 `G6neg` 一枚**（与 `A395` 同名的在册红），但它 `声明=ring 基线=1枚 实测=3枚 因=新增未成对`——量的是 `OpenTask`↔`CloseTask` 不过滤任何文件那一味；**本程三枚文件 `OpenTask` 命中 0**（尺：`grep -rln "OpenTask\|CloseTask" --include=*.go internal/ cmd/` ⇒ `internal/tools/bridge.go`／三枚既有 pointer＋ticket175r2 测试件／`cmd/wisp/run.go`／`cmd/wisp/task_scope_close_151_test.go`），那两枚增量不在本程写面，本程**未修不当绿**，也未在"起手锚点减本程改动"的树上复算，所以归属只到"不是我写的字面量"这一层。⑤ 面板侧若以**别的方式**（不经 `Record`／`Backfill`／`TaskOutput` 字面量／`.State` 赋值／`Method*` 常量）触碰这一维，这扇门不覆盖；`AC#1` 子代理那一维**没碰**（票 197 已另立）。⑥ `TaskLog.State` 与 `TaskOutput.State` 的映射表**一张都没画**（A394／票 196 射程，本程按令不统一）。
  - **纪律**：预算 ≤30 枚调用，**实耗约 36 枚＝超 6**（具名：`gate-clauses` 与 `scripts/` 路径定位 3 枚、`gofumpt` 不在 PATH 改走 `go env GOPATH` 重取 2 枚、判据自身两发缺陷各修各验 4 枚〔`ast.Inspect` nil 收尾节点 panic＋`tools.TaskOutput` 限定名让第三枚正控假绿〕、起手名册核对与首枚 commit 失败重来 2 枚）。被拒 0 次。首枚 commit 因未跟踪文件未 `git add` 而失败一次（rc=1，`git log` 当时顶枚仍是 §F 的 `4db5f3f6`＝**没有任何 commit 落地**；补单枚显式 `git add` 后重提，未误带别家 staged 件）。`frontend/**`／`design/**` 写面 0；`go.mod`／`go.sum`／`docs/PLAN.md`／`docs/specs/**`／阈值／golden／`bridge.go:42-45` 一字未动；仓内零删除；只 commit 未 push；**`AC` 框一枚未勾**（勾要非实现者表）。⚠ 更正：本节首稿写"实耗 34＝超 4"是把收尾两枚漏计后的读数，**权威读数以证据件 §6 为准**（`docs/evidence/s1/188-task-state-r2.md`）；同批把 `AC#2` 那一格与表头行之间被误吞的换行补回（一次 Edit 事故，只动我这一程新写的行，别人原文未触）。
- 09-28 18:2x 编排者收 `188-r2`（三枚 commit；`internal/tools/task.go:124` 状态维进树，取值集只能是 `states.go:11-31` 那 **20 枚**，我现数一致；`go test ./internal/tools/` 我复跑 **`ok 16.062s`**）。⚠ **本程落的是"字段＋门"，不是"能显示"**：**生产写入者今天仍为 0**（`cmd/wisp/run.go:640-647` 不带 `State`，本程禁写那枚文件）⇒ 真任务读起来全是"宿主未登记"。**这是本仓今天第三次出现同一形状**（`A382`／`A393`／本轮），所以排 **`188-r3`＝生产者那一跳**（只准填那 20 枚里的名字），排 `145-r3` 交件之后。`AC#2`／`AC#3` 两格**仍不勾**（勾要非实现者表）。账 `A398`。
- 11-07 18:1x `188-c2`（**只读取证＋裁定腿**，本波零产码——`internal/panel/**`＋`cmd/wisp/**` 写面归 213 腿独占；锚 `29b3501d5`）：证据件 `docs/evidence/s1/188-monitor-state-source.md`。三核＝**票面「现量」整段已过期**（09-28 锚 `e9ef94d0` ⇒ 子代理维今天已被票 197 落成整链，票面那把 0 命中的尺现跑 **19 枚文件**）。S1 字段表逐枚核毕：**子代理 Status 已全链在场**（源＝`PublishSubagent` 落 `Thinking`＋197 收尾落 `Settling/Muted/Error` 借名 `subagent_197.go:99-108`；reader＝`StateAnswer`→`cmd/wisp/panel_pump.go:222`，状态读只经宿主判官）；**根任务 Status 是今天真缺口**（`MarkRoot` 刻意不落 State（task.go:419 注释指 A394 归 196），196 已 done 但只落 memory 侧映射 `internal/memory/taskstate.go`，名册根行仍空；`run.go:1855` 构造 `TaskBackfill` 不带 State——补它＝**一行装配**，D43 名从 `memory.TaskStateFor(taskLogState(res.Status))` 取，零自造）；**「已拒绝」无源**（spawn 四处拒绝 `subagent_197.go:254/:281/:288/:316` 全不留名册行；D43 借哪个名＝规格真空，先裁后写）；**时间维规格真空**（`TaskOutput`/`TaskRow` 无钟点字段，立票归 197 射程）；「排队中」第二源（scheduler `TaskListRow` 那张表）撞 213/197 射程另派。**AC#4 与 145 的分账被 197 改写**，今天形状＝188 交源（余根任务态/拒绝行/时间维三枚）·197 交载体＋宿主映射（已落）·**票 35 交出向传输（未落**：`PostWebMessage|EvaluateScript|CreateWebMessageAsJson` 本程复跑非 test **0 命中**）·前端交页面读取（未落：`frontend/src/lib/panel.ts:256` 已声明 `tasks?`，但 `.tasks` 读取 **0 文件**，本程现跑）——145 对 S1 的「页面要、Go 没交」清单已清空。**AC#3 反向判据复核成立**：`TestTaskState188AC3PanelHasNoWriteLeg`（`task_state_188_test.go:351`）在场，面板包无 State 赋值/Record 形状。顺手更正 `145-a2` 一笔过期账：`PumpSources.L1Windows` 今天**已装配**（`cmd/wisp/run.go:1023`→`panel_pump.go:108`），票 220 第二来源不再是黑的（归 220/编排者核收）。裁定＝**下一步写腿顺序：①根任务状态维**（213 收口后，一行装配，判据 `task_state_188_test.go` 现成）**→ ②拒绝行**（先裁 D43 借名再写）**→ ③时间维**（立票，源＋载体双侧新字段）。门禁：零 Go 改动 ⇒ gofumpt/vet/d22scan **无对象可扫**；零测试运行（65/66/213 三腿在飞，不把别人红绿记本程头上）；`frontend/**`/`design/**` 只读零写面；零 push；AC 框一枚未勾。
- 11-07 19:1x `188-r4`（写腿，只做「根任务状态维」＝取证腿 `9ad02eb50` 裁定顺序①；head 锚 `2fec9f878`）：**一行装配落地**＝`cmd/wisp/run.go` 现量构造 `TaskBackfill`（票面 1855，漂移后现量 1860）补 `State:`，D43 名取 `memory.TaskStateFor(taskLogState(res.Status))`（196 的对表读口，零自造）；三核现跑＝取证腿三条复核成立、**行号漂移 +5**、**票面判据件指认过期**（`internal/panel/task_state_188_test.go` 不存在，现量 `internal/tools/task_state_188_test.go:351`）；诚实边界＝无正文终局（failed/stuck/cancelled）被 Backfill 既有 fail-closed 空正文门整条拒绝、名册里没有行（不归本步不修），运行态与取消兜底归方案乙（要改判 197 根行钉，需编排者具名）；判据新件 `cmd/wisp/task_state_188_test.go`（真组合根三面读数：名册 State=Settling／StateAnswer=(Settling,"")／面板行 statusKnown=true），牙 M1 摘两行装配**必红**（rc=1，红行逐字存 `docs/evidence/s1/188-root-task-state.md` §4），sha256 复原 OK；门禁＝gofumpt 点名件空／vet rc0／build 两包 rc0／`go test ./cmd/wisp -run 'Test188|TestTaskState'` **ok 3.607s**／panel **rc=1 共 4 枚红**（票面预告 2 枚 TS 声明红＋多出 2 枚：`Test187PanelFilesTouchNoPermissionSurface`=187 腿半落缺 `model_switch_187.go`、`TestC21DesignTokensFourWayAgree`=ball 在飞 tokens 面——逐枚归属非本程）／d22scan **rc=1 2 枚红**＝`config_write_187.go:55` phantom-citation（187 在飞未提交件引的证据件盘上不存在，`git log` 空读数为凭）；197 载体钉复跑 **ok 3.265s** 不打红；import 增量现量 0 枚（多返回值接法类型名不出现，票面预估 1-2 行被编译器现量否定）；零触碰 internal/ball、internal/tools、internal/config、frontend/design；只 commit 不 push；AC 框一枚未勾。
- 11-10 18:0x `188-s1`（**停手上报腿**，本程零产码；head 锚 `e90ea3cf0`；证据件 `docs/evidence/s1/188-rejection-row-s1.md`）：接 c2 裁定顺序②「拒绝行（先裁 D43 借名再写）」⇒ **核毕停手（D22③）**，两处卡点：① **D43 借名规格真空未裁**——现读 `internal/statemachine/states.go:11-31` 20 枚名里没有「已拒绝」；spawn 拒绝点现量比 c2 记的四处更多（`internal/tools/subagent_197.go` `:254/:269/:276/:281/:288/:302-304/:309/:316` 全部早退零名册行，唯一 `Record` 在 `:462` 成功路），「已拒绝无源」今天仍成立，而 A1611 之后至本程现读（tail A1680）**无人裁借名**（最近先例形状＝A434 cancel→`StateMuted`）；② **写面不在本腿射程**——落点 `internal/tools/subagent_197.go`（10-14 后执行面在 tools 包），本腿派单写面＝`internal/agent/**`，改它须编排者另派 tools 写腿（同包单写腿）。顺序①根任务态核实在场未动（`cmd/wisp/run.go:2055-2056`）；顺序③时间维按 c2 裁「立票归 197 射程」非本票。门禁现读：`go test ./internal/tools -run TestTaskState188 -count=1` ＝ **ok 0.116s 4 枚 PASS**（AC#3 负向腿 `27 panel non-test sources, 0 write legs`＋五枚正控逐枚 red as required 现跑复现）；`go test ./cmd/wisp -run 'Test188|TestTaskState' -count=1 -timeout 30m` ＝ **ok 3.875s**；`go test ./internal/panel/ -count=1` ＝ **ok 16.638s**（与 A1680 一致今天全绿）；`sh scripts/d22scan.sh` ＝ **clean**（ban#8 internal/=974／cmd/=250，票面基线 433 过期＝仓增）；零 Go 改动 ⇒ gofumpt/vet 无对象可扫（循 c2 先例）。AC#4 分账现口径已写进证据件 §2。纪律：零产码、零删除、只 commit 不 push、AC 框一枚未勾、`frontend/**`／`design/**` 零触碰。
- 11-13 13:0x `188-panel`（写腿，A1827/A1828 派，写面＝`internal/tools/**`＋`internal/panel/**`＋`cmd/wisp/**`；head 锚 `8a8673fe8`）：**产码 1 处＋判据 2 枚，两堆核毕**。证据件 `docs/evidence/s1/188-panel-status-sources.md`。逐格：
  - **行号漂移核毕**（证据件 §1 对照表）：票面 `task.go:110-152`／`Record:123` ⇒ 今天 `Record=task.go:281`、`State` 字段 `task.go:196`、写侧门 `task_backfill.go:130`（票面 113）；c2 记「run.go:1855 不带 State」已过期（r4 落的 `run.go:2185-2189` 带 State）。
  - **堆一（子代理状态）＝有来源全链在场，零改动**：生产者 `subagent_197.go:410-415`（Thinking）／`:450-465`（Settling/Muted/Error），判定门 `task.go:224-234`，载体 `panel/pump.go:202`＋`subagent_roster_197.go:51-143`，读口 `cmd/wisp/panel_pump.go:207-271`，接线 `run.go:1201`。余缺口三枚如实记：拒绝行（八处早退零名册行，D43 借名规格真空未裁）、时间维（归 197）、根行**运行态**空窗（`MarkRoot` 刻意不落，`task.go:417-425` 注释归 A394/196；`cmd/wisp/subagent_carrier_197_test.go:324` 把该空态钉成**正向断言**，补它要改判别人的钉，本程不碰）。
  - **堆二（后台任务状态）＝修「空正文终局整条丢行」**：真缺口在 `task_backfill.go` 旧门 `case text == ""` **整条拒绝** ⇒ failed/stuck/cancelled/truncated 终局（环路 finish 只发 message 不发正文）连状态维一起从名册蒸发，面板读「宿主没有登记这一维」＝谎话（host 明明知道终局）。改门为 `case text == "" && b.State == ""`（`task_backfill.go:123`）：带 D43 名的空正文终局照登（State 真、Text 空、无 artifact），空态＋空正文维持旧拒，非空正文路径逐字节不变。**零自造词汇**：名字仍只从 `memory.TaskStateFor(taskLogState(res.Status))` 对表读数（`run.go:2184`），`statemachine.Valid` 门仍在（:130）。
  - **技术选型（⓪⓪⓪ 自拍）**：主流 harness 面板的终局行常驻可回看 ⇒ 修让面板对终局撒谎的那枚缺陷最有价值；补根行运行态（要改判 A394＋重钉别人正向钉）、拒绝行（规格真空 D22③）、时间维（归 197）三件价值都不如它且碰别人地界。⛔ 未拿"动作小"当理由——不修的代价＝终局栏对失败永远谎报。
  - **成对判据（证据件 §4 逐字）**：正向＝completed→Settling／failed→Error 两源，面板读数随之变（`statusKnown=true status="Error"`）；反向＝门改回旧形状（mutation，sha 备份 `D:\tmp\wisp\20261113\188-panel\task_backfill.go.keep` 改前 `15AC30EB…7013F0` 复原同值）⇒ tools 判据红（`task_state_188_test.go:245` 逐字）＋cmd/wisp 判据红（`roster state=""`／`StateAnswer()=("","宿主没有登记这一维…")`／`statusKnown=false` 逐字）——红在明处不静默；复原后全绿。
  - **门禁读数**：build rc0／vet rc0／`go test -count=1 ./internal/tools/ ./internal/panel/`＝`ok 37.335s`／`ok 20.997s`（panel **今天全绿**）；gofumpt（走 `D:\work\base\gopath\bin\gofumpt.exe`）无输出（首跑点名新判据件一处，按 fumpt 形状重写后清）；`sh scripts/d22scan.sh` **clean**（bans#1-5 internal/=404 cmd/=77；ban#8 internal/=1005 cmd/=276）。⚠ `go test ./cmd/wisp/` 直跑＝`exit status 0xc0000135`（票 98 已定性，不算红），改走 `scripts/wisp-cli-staged-tests.sh`：`-test.run 'Test188'` ＝ 2 枚 PASS。
  - **纪律**：`internal/panel/**` 本程零字节（核毕无缺陷，缺在生产者跳）；`task.output` stub／memory 侧／词表统一一枚未碰；撞「整进程死亡」＝0（0xc0000135＝票 98 加载环境现象，非 314 形状）；零 push；AC 框一枚未勾；临时件只建不删。
---

## ★结案（2026-11-13 13:1x，编排者翻框，非实现者）— **修掉"空正文终局整条丢行"这枚真缺陷**

### 三把尺（我现跑，全过）
- **提交真在**：`e8bb1f2b4`（5 files, +280/−1）
- **件真在**：`docs/evidence/s1/188-panel-status-sources.md`（12417B）
- **自己复跑**：`ok tools 36.065s` ／ `ok panel 22.274s` ／ `ok internal/agent 6.345s`（三包全绿）

### 修的是什么（一句话）
`task_backfill.go` 旧门 `case text == ""` **把整条记录拒掉**，连**状态维**一起扔。而 `failed`/`stuck`/`cancelled`/`truncated` 这四种终局**按构造就没有正文**（环路 finish 只发 message）⇒ **它们在名册里整条蒸发**，面板对着一个**已经明说过的终局**读「宿主没有登记这一维」＝**面板在说谎**。

### 修法（最小且保约定）
```go
-	case text == "":
+	case text == "" && b.State == "":
```
- **有 D43 状态名的空正文终局** ⇒ 现在进名册（行里说清"发生了哪个终局、没有正文"，**不是**"正文是空的"）
- **既无正文也无状态** ⇒ **旧契约逐字保留**（那条记录确实什么都不携带）
⇒ **状态维与正文解耦，零自造词汇、零新契约面。**

### 成对判据（腿做，我核了件）
- **正向**：`failed run: roster state="Error" text=0 bytes, panel statusKnown=true status="Error"`，`--- PASS: Test188FailedRunStillFilesItsD43State`
- **反向（关键）**：把门改回旧形状 ⇒ tools 判据红（`task_state_188_test.go:245: an empty-bodied ending with a D43 state was refused whole`）＋ cmd/wisp 真组合根红（`panel row statusKnown = false`）；sha 备份 `15AC30EB…` 复原同值
- **红在明处，不静默**

### AC #2 — #6 逐格依据（不是印象）
| AC | 依据 |
|---|---|
| **AC#2** | `Record.State`＝`internal/tools/task.go:196`，类型 `statemachine.State`（**D43 对齐靠类型本身**，不自造词表）；本程新判据 `TestTaskState188AC2EmptyBodiedEndingKeepsItsState` 钉"终局状态不再被门吃掉" |
| **AC#3** | `TestTaskState188AC3PanelHasNoWriteLeg` 在 `internal/tools/task_state_188_test.go:413`，**现跑 PASS**（负向腿在位，面板仍无写权） |
| **AC#4** | 证据件 §6 逐条写明欠谁：载体十六字段归 197 已交、出向传输仍零命中归票 35／前端腿、两套词表并表在票 196 名下 |
| **AC#5** | 改动 5 件**逐枚核过**：无一枚是 `PLAN.md`／`specs`／`thresholds.go`／golden／`allowlist.txt` ⇒ **契约轴零字节** |
| **AC#6** | 三包 `-count=1` 全绿（我复跑）；`gofumpt` 无输出；`sh scripts/d22scan.sh` clean；`cmd/wisp` 走 staged 载体 2 枚 PASS |

### ★腿如实上报的三件（**不归本票**，我具名收下）
1. **根行「运行中」空窗**：`MarkRoot` 刻意不落，且 `subagent_carrier_197_test.go:324` 把该空态**钉成正向断言** ⇒ 补它要**改判 `A394`＋重钉别人的钉**。
2. **spawn 被拒的「已拒绝」行**：八处早退零名册行，**要先裁 D43 借哪个名**（规格真空）。
3. **时间维**：归票 197。

**撤销口令：「188 回退」**——`git revert e8bb1f2b4` 即恢复旧门（但面板会对失败行重新说"宿主没登记这一维"）。

