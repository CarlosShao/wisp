# 167-a2 只读普查（有写工具的那一版）：AC#5「草稿」＋ AC#6「崩溃自救」

> 本件由腿 `167-a2` 自写自 commit。前一版同路径件（编排者代落、8,300 字节）**已被本版取代**，本版逐条自跑尺、不抄前版行号；前版的存废与是否合流由编排者裁。

## 0. 起手锚与并发窗口

（本节答：起手 HEAD `f761a017`／时刻 `2026-10-05T09:15:08+0800`／在飞＝写腿 `267-r1`＋只读腿 `ci-phantom-1`／搜索根只 `cmd internal tools docs .scratch scripts`）

- 起手尺：`date`＋`git rev-parse --short HEAD`＝**`f761a017`**，branch `dev`，`git status --porcelain` 起手**非空**（共享工作树里他人的写面：`design/**` 的删除项、`frontend` 一侧、`.scratch` 各件与若干 `commit-msg-*.txt`）。**本件全程不碰它们**，收口尺只数我名下那一枚路径。
- ⚠ **骨架落盘比派单要求的"前 5 次工具调用"晚两次**（自陈，不粉）：第 1–6 次用在起手锚＋读前版同路径件＋读票面 167 全文上，骨架写在第 7 次、骨架 commit＝**`aa016492`**。读数因此不是"先建框后补料"的假骨架，但形状确实迟了两拍，归编排者记账。
- 同机在飞（按派单登记，本程未验证其存在）：写腿 **`267-r1`**（`internal/config`）＋只读腿 **`ci-phantom-1`** ⇒ **零 `go` 命令**（见 §8 名册），所有判据都标〔仅读码〕或〔尺〕。
- **与前一版同路径件的关系**：本版取代它；前版＝编排者代落（`git log --oneline -- .scratch/wisp/probes/167/a2/census.md` 唯一命中 **`8fa4246d`**「ledger(A590) 收 167-a2 代落…」，8,300 字节）。**本版不抄那一份的行号**，§1–§6 每一枚 `file:line` 都是本程自己 `grep -n` 被指的那串字符之后落的名；`167-c1/c2` 的三格读数**未重跑、未覆盖、未冒充本程原创**。
- 词表（自扩，不只派单列的六枚）：`draft`／`unsent`／`not.sent`／`outbox`／`pendingSend`／`retryText`／`composerText`／`stash(ed)Text`／`preserveText`／`retain`／`rehydrate`／`pre.?fill`／`swallow`／`clobber`／`Reset|Clear`／`修法`。

## 1. 草稿：承接体与生产者-消费者名册

（本节答：**无**——Go 侧零枚字段或结构承接"这次没发出去的字"；连"发出去"那条路本身都是空的）

**1.1 入向那一寸（有封套，没有听众）**
- 封套＝`internal/panel/bridge.go:91`（尺：`type ComposerRequest struct`）；装用户原文的槽＝`:98` 逐字 `Text string \`json:"text,omitempty"\``。**它是一次性入向封套，不是存储位**——`ParseComposerRequest`（同文件）解完即交派发。
- 门牌＝`internal/panel/bridge.go:45` 逐字 `MethodMessageSend = "panel.message.send"`；派发表＝`internal/panel/composer_dispatch.go:175`（`func (d *ComposerDispatch) dispatch`），其 case 在 `:192`。
- 生产装配把这一枚处理器**留空**：`cmd/wisp/panel_inbound.go:271` 起那枚 `&panel.ComposerDispatch{`，`:277` 逐字 `Message:    nil,`（同块 `:275 Workspace: nil`、`:276 Attachment: nil`）。⇒ 今天 `panel.message.send` 一到就撞 `internal/panel/composer_dispatch.go:193-194` 的 nil 检查 ⇒ `unattached`（`:230`）。
- 入向有两个真听众，都在别的路由上：`cmd/wisp/panel_inbound.go:163`（`reply, err := disp.Handle(ctx, raw)`）与 `cmd/wisp/panel_host_windows.go:637`（`return m.disp.Handle(ctx, raw)`）。

**1.2 失败路径上"原文"今天去了哪里（三处都不带它）**
- 审计行：`internal/panel/composer_dispatch.go:240` 的 `record` ⇒ `:244` 逐字 `"panel: INBOUND-DISPATCH request=%q method=%q source=%q origin=%q err=%v detail=%q"`，字段里**零枚是 `req.Text`**；`detail` 是 `:245-246` 那句固定散文。
- 回给页面的句子：`internal/panel/bridge.go:157` `RefusedEnvelopeForUser` ⇒ 逐字 `"面板请求被拒绝 [%s %s]：%v"`（method＋requestId＋err），**不含原文**。
- `raw` 本身：`Handle(ctx, raw)` 收进去之后只在 config 那两条路由上继续活着（`composer_dispatch.go:201/:206` 传 `raw`），message 路由 `:196` 只交 `req` ⇒ 函数返回即随栈丢弃。**〔仅读码〕没有一枚 Go 变量把它抄进任何长期容器。**

**1.3 全仓"能装用户原文"的 Go 结构穷举（消费者名册）**
- `internal/panel/composer.go:295`（尺：`type OutgoingMessage struct`）——唯一为"一次发送"而生的文本载体，槽在 `:296` `Text string \`json:"text"\``。**它的产码构造者＝零枚、方法 `ForAgent`（`:308`）产码调用者＝零枚**：尺 `grep -rn "OutgoingMessage" cmd internal tools --include=*.go` 的全部命中＝定义自身三行（`:292/:295/:308`）＋一枚测试 `internal/panel/attachments_test.go:284`（`:285` 调 `ForAgent`）。⇒ 判语：**这条"发送"路今天整枚没接**，谈不上"失败后留住原文"。
- 出向快照里能落自由文本的槽（逐枚看过，都不是草稿位）：`Snapshot`（`internal/panel/composer.go` 内，键名册 `Pending/Results/Composer/GeneratedAt/Instructions/Tasks`）→ `ResultChunk.Text`（`:97`，助手流）／`ApprovalCardView`（`internal/panel/approval.go:39`，`:48 Reason string json:"reason"`、`Args []string`，卡面文本）／`ComposerState`（`internal/panel/composer.go:235`）11 键逐枚＝`mode`/`workspace`/`attachments`/`acceptedAttachmentMimes`/`maxAttachmentBytes`/`attachmentError`/`git`/`currentModel`/`modelKnown`/`credentialState`/`credentialKnown` ⇒ **零枚可装"你没发出去的那句话"的槽**。
- `TaskRowView`／`TaskRosterSection` 在 `internal/panel/subagent_roster_197.go:106`／`:150`（本程只登记位置，其读数属 `167-c2` 那三格，不重跑）。
- 词表扫的净结果：`draft` 在 `cmd`＋`internal`＋`tools` 的 `*.go` 里共 **14 行命中**，逐枚是散文（`internal/observe/logging.go:383`、`internal/ball/sta_release_windows_test.go:62,376`、`internal/config/migrate.go:30,102`、`internal/config/schema.go:22`、`internal/tools/paths_ticket107b_probes_test.go:141`、`cmd/wisp/config_reload_perm_223_windows_test.go:66`、`cmd/wisp/resident_hotkey_v1probe_test.go:72`、`cmd/wisp/slo_exit_os_156_windows_test.go:80,336` 全为"本文第一稿/早前那稿"）或 fixture（`internal/tools/fs_edit_ac34_test.go:75,86` 的 `status: draft` 字符串）。**零枚是字段名、零枚是结构名。**其余同族词（`unsent`／`outbox`／`pendingSend`／`retryText`／`composerText`／`stashed`／`rehydrate`）在 `cmd`＋`internal`＋`tools`＋`scripts` **零命中＝没找到＝这一格今天不存在，不是坏消息**。
- ⚠ 顺带一枚**容易被误当草稿位**的东西：`internal/panel/pump.go:427 streamKept`／`:441 absorb`／`:595 clampLocked` 那族"保留窗"是**助手流的尾部窗口**（本程只登记命名，判它不是草稿位＝读码级：它吸的是 delta，不是 `req.Text`）。

## 2. 草稿：存不存盘、谁清、谁覆盖

（本节答：**盘上零处、会话内零处**；"清"没有责任人；"新写的字不许被覆盖"今天**没有任何一格负责**）

**2.1 唯一带用户文本的持久列，只给"已受理"的任务**
- 建表名册（尺：`grep -n "CREATE TABLE" internal/memory/schema.go`）＝**九枚**：`:23 schema_meta`／`:30 profile`／`:40 memory`／`:53 task_log`／`:69 tool_call`／`:89 approval_grant`／`:102 cost_daily`／`:111 plugin_state`／`:136 provider_health`。
- 其中唯一装用户原话的列＝`internal/memory/schema.go:58` 逐字 `query_text     TEXT NOT NULL,          -- 已脱敏；模式匹配掩码是尽力而为，须标注（§14.4）`。
- 写手链：`internal/agent/journal.go:122`（`func (t *taskJournal) startTask(ctx context.Context, query string)`）⇒ `:129` `j.StartTaskLog(ctx, memory.TaskLog{ID:…, State: "running", QueryText: query})` ⇒ DAO `internal/memory/dao_tasklog.go:24`（`StartTaskLog`，`:33` 那句 INSERT）。`startTask` 的调用点在 loop 内（任务已受理之后）。⇒ **判据：这条列今天只收"已经进入 loop 的那一发"；一次没被受理的发送（`panel.message.send` 走 `unattached` 那条）根本不经过它。**
- 接口面也印证同一件事：`internal/agent/journal.go:18-20` 的 `Journal` 只暴露 `StartTaskLog`/`FinishTaskLog`，**零枚"存一条未发出的草稿"的形状**。
- **第二条持久面不存在**：`ls internal` 的名册＝`agent audio ball buildinfo config llm memory models observe panel perm plugin proc projctx risk secret session speech statemachine streamkey tools watchdog winsec`——**没有 `store` 这个包**（尺 `ls internal/store` 回 `No such file or directory`）⇒ 想落草稿也没有现成的落盘面。
- 数据根下的文件面（具名，〔仅读码〕）：`config.toml`／`logs\`（`cmd/wisp/logsink.go:76` 逐字 `const logDirName = "logs"`、`:87 logSinkDir`）／`models\`／`artifacts\`（附件字节；`internal/panel/attachments.go:61` 逐字 `ArtifactsDir() string`，注释 `:56` 写明"`*memory.Store` satisfies it in production"）／`memory.db`。⇒ **附件字节今天存盘、原文不存盘**，这两件事的方向正好相反。

**2.2 会话内存里也不存**
- `ls internal/session`＝`doc.go grants.go grants_test.go session.go ticket224_pattern_dialect_test.go`；尺 `grep -rn "Role.*user\|\"user\"" internal/session/*.go`（排测试）**零命中＝这个包只管会话与授权，不管话**。
- `ComposerState`（`internal/panel/composer.go:235`）是**出向**快照的一节，11 键里零自由文本位（§1.3 已逐枚列），它**不回写页面输入框**，也不承担草稿。
- 附件族里"清空"这件事在 Go 侧也**没有产码函数**：尺 `grep -rn "func.*Reset\|func.*Clear\|attachments = nil" internal/panel/*.go | grep -v _test` **零命中** ⇒ **谁清它＝零枚**（因为没有需要清的那一格；不是"我没找到"）。

**2.3 "用户后来新写的字不许被它覆盖"今天由谁负责＝没有人**
- Go 侧既没存原文、也不回写文本 ⇒ **这一条今天整个在页面那一寸**（⛔ 本程不读 `frontend/**`，判不了，见 §7-1）。
- 全仓唯一同族"新的不许顶旧的"的**已接产**在配置那一轴：`internal/config/writeguard.go` 的守护写＋两枚尺族（§3 逐枚列名）。它们断的是 `config.toml` 的键，**射程里没有 composer 文本**——⛔ 不许有人把"已有 no-clobber 机制"当成草稿这一格已有人管。

## 3. 草稿：既有尺名册

（本节答：**零枚**。断"草稿不吞／失败原文单独存住／新字不被覆盖"的既有测试＝**0 枚**；下面给邻近者的逐枚名册）

- 直接靶（草稿）：**零枚**。`grep -rn "draft" cmd internal tools --include=*.go` 的 14 枚命中里没有一枚是测试用例名或断言（§1.3 已逐枚具名，其中只有 `internal/tools/fs_edit_ac34_test.go:75,86` 在测试里，而那是 `status: draft` 的**文件内容 fixture**，与 composer 无关）。
- 沾到 `panel.message.send` 的用例**共 15 处**，逐枚看**没有一枚读 `req.Text` 的内容**，全部是"名册／派发表一致性"用途：
  `cmd/wisp/panel_config_248_test.go:467`；`internal/panel/bridge_test.go:62`（封套里塞了 `"看看这个"` 只为验证能解析）、`:140`；`internal/panel/composer_dispatch_test.go:261`、`:292`；`internal/panel/composer_handlers_test.go:303`；`internal/panel/composer_test.go:414`；`internal/panel/git_test.go:382`、`:450`、`:508`；`internal/panel/inbound_roster_253_test.go:73`、`:659`、`:703`。
- 唯一沾"发送文本"的用例＝`internal/panel/attachments_test.go:284` 造 `OutgoingMessage{Text: "看看这个", …}`、`:285` 调 `ForAgent` 后断附件形状（拒收的附件要在文本里逐字重复）。**它断的是"附件别被吃掉"，不是"草稿别被吞"**——同一票面 §「do not eat the user's intent」的**另一维**，具名免得下一位误认。
- no-clobber 那一家（**同族不同靶**，逐枚名＋断什么）：
  - `cmd/wisp/always_write_no_clobber_226_test.go:39` `TestAC1AlwaysBranchDoesNotRevertAHandEditedKey` —— 断"always 那条分支不许把人手改过的键回落"。
  - `internal/config/writeguard_226_test.go` 13 枚：`:63 TestAC1AllowedDirsWriteKeepsAHandEditedKey`、`:102 TestAC1ControlSnapshotWriteIsWhatRevertsTheHandEdit`、`:117 TestAC2GuardedWriteReportsTheKeysItWroteAndTheOnesItKept`、`:146 TestAC2CleanWriteReportsOneKeyAtInfoLevel`、`:164 TestAC2UnreadableFileIsRefusedAndNothingIsWritten`、`:187 TestAC3AdoptionClaimsOnlyWhatThisWriteProduced`、`:234 TestAC3HandAddedEntryIsPreservedButNotAppliedToMemory`、`:269 TestAC4PermissionModeWriteKeepsAHandEditedKey`、`:308 TestAC5SetAllowedDirsPersistsARevocationWithoutAClobber`、`:330 TestAddAllowedDirStillRefusesARelativeHop`、`:360 TestGuardedWriteFailureRollsBackMemory`、`:393 TestGuardedWriteDoesNotRewriteAFileThatAlreadySaysIt`、`:416 TestDiffKeyPathsNamesRealKeyPaths`。⇒ **全是 `config.toml` 的键级 no-clobber**，与 AC#5 那发"失败时整块回填"的变异毫无重叠。
- ⛔ 三枚冻结件按令未读：`internal/panel/tokens_fourway_test.go`、`internal/panel/l2_grant_boundary_test.go`、`internal/perm/ticket90_persist_test.go`。本程在包级 `grep` 里**顺带见到**第二枚的文件名与行号，一律不转写、不据此下结论；若编排者要判"草稿出向字段会不会撞那枚白名单尺"，得由**非本程**的腿去读。

**3.1 C17（`PanelBridge`）里有没有草稿相关方法名——零枚，且本程不新造**
- 代码侧入向名册（逐字，`internal/panel/bridge.go:41-45` 那块常量；枚位尺＝`grep -nE '=\s*"panel\.' internal/panel/bridge.go` 回 42/43/44/45 四行）：`panel.mode.request`／`panel.workspace.request`／`panel.attachment.add`／`panel.message.send`，另有票 248 的两枚 `config.get`/`config.set` 形状（同块注释自述它们**刻意不带 `panel.` 前缀**）。⇒ **无草稿格。**
- 契约文本侧（`docs/PLAN.md`，逐字引 C17 那几枚）：`:1367`「**C17** | **`PanelBridge`** | 前端↔Go 双向通道：`invoke(method, args) → result` + Go→前端事件推送（流式结果、审批请求、任务状态）。**回复必须按 correlationId 路由**；前端**必须无状态**（WebView 销毁后一切从 Go 侧重读）」；`:2434`「**配套**：C17 `PanelBridge` 必须给出**方法白名单**（只有列出的方法可被前端调用）」；`:2968` 起「**C17 `PanelBridge` — 补四项**：① **方法白名单**…② 每个方法标注所需 capability…③ correlationId 路由 + 事件推送的背压与合并（D38(d)）④ **前端无状态的强制手段**：面板每次 `show` 必须发 `panel.resync`…」。⇒ **PLAN.md 的 C17 条文里零枚草稿方法名**，只有 `panel.resync` 这一枚推送名（`docs/specs/SPEC-08-ui-ball-panel.md:150`、`:163` 表首行）。
- `docs/specs/SPEC-08-ui-ball-panel.md:161-174` 那张提案白名单表逐枚（本程全文读完）：`panel.resync`／`tasks.list`／`task.detail`／`history.query`／`transcript.get`／`approval.current`／`approval.queue`／`approval.decide`／`config.get`／`config.set`／`grants.list`／`grants.revoke`／`privacy.purge`／`privacy.export`／`cost.summary`／`models.list`／`models.delete`／`diagnostics.export` ＋ 事件推送 `task.delta` `tool.chip` `approval.request` `ball.state` `cost.tick`。⇒ **零枚草稿/draft 方法**；最接近"草稿出口"的是 `diagnostics.export`，但它是**诊断包**不是输入框（§5 归到崩溃那一格去量）。
- ⛔ 本件不新造任何方法名（派单硬约束）。若 AC#5 要一枚出向草稿字段，那是**§7-5 的裁定面**，不是本件的产出。

## 4. 自救：谁产／谁投／谁落三处齐读

（本节答：产＝`internal/observe/goroutine.go:286` 那枚唯一受许可的协程体；投＝`:312-315` 交给 `r.sink`，默认值在 `NewRegistry` 体内 `:236` 就装好；落＝`cmd/wisp/logsink.go:153` 那台 tee 的 primary → JSONL 文件。**"生产没装 panic sink"那句假话本程独立复核为假**）

**4.1 谁产（panic → 结构化记录）**
- `internal/observe/goroutine.go:286`（`func (r *Registry) run(...)`）——`:285` 注释逐字 `// run is the single sanctioned goroutine body of the codebase.`
- `:294` `defer func() {` ⇒ `:296` `if rec := recover(); rec != nil {` ⇒ `:297-304` 组装 `PanicEvent`，其中 **`:302` 逐字 `Stack:     string(debug.Stack()),`**；`:301` `Recovered: fmt.Sprint(rec)`；`:303` `At: NowWallUTC()`。
- 记录形状：`:153` `type PanicEvent struct`，文本槽 `:159` 逐字 `Stack     string     \`json:"stack"\``（同族键 `goroutine`/`owner`/`root_id`/`error_class`/`recovered`/`at`）。
- 计数：`:308-310` `r.panicCount++`；出口 `:375` `func (r *Registry) PanicCount() uint64`。
- 崩溃之后进程看到的错误：`:317-320` 那句 `&Error{Class: ClassInternal, Detail: fmt.Sprintf("panic in goroutine %q (owner %s): %s", …)}`——**含 recover 值、不含栈**（这枚区别就是当年票 249 误判的地方）。

**4.2 谁投（记录 → sink）**
- `internal/observe/goroutine.go:312` `sink := r.sink`（`sinkMu` 保护）⇒ `:314` `if sink != nil {` ⇒ `:315` `sink(ev)`。
- sink 字段：`:220` `sink PanicSink`；类型：`:165` `type PanicSink func(PanicEvent)`（`:163-164` 注释逐字 `PanicSink receives panic records. The default writes a structured slog error; ticket 08 replaces it with the JSONL pipeline.`）。
- ★ **默认装在装配点**：`:236` 逐字 `sink: defaultPanicSink,`（在 `NewRegistry` 体内）。⇒ `:314` 那句 `if sink != nil` **在生产里不可能走空**。
- 替换口：`:240-241`（`func (r *Registry) SetPanicSink(s PanicSink)`），`:245` `s = defaultPanicSink`（传 nil 会把默认装回来）。
- ⚠⚠ **本程不复演那枚旧错**：尺 `grep -rn "SetPanicSink" cmd internal tools --include=*.go | grep -v _test` 的结果**只命中定义自身那两行（`:240`/`:241`）＝非测试调用者零枚**——这枚读数证的是"没人替换默认"，⛔ 它证不了"没有默认"。当年票 249 正是把这两件事写成了一句（现由 `docs/reports/pending-and-issues.md` 与票 249 票面自纠、票已撤回；本程独立从 `:236` 复核出同一结论）。测试侧调用者＝`internal/observe/goroutine_test.go:117-118` 两行。
- 默认 sink 干什么：`:167` `func defaultPanicSink(ev PanicEvent)` ⇒ `:168` `slog.Error("goroutine panic recovered", …)`，**`:174` 逐字 `"stack", ev.Stack,`**（同块还带 goroutine/owner/root/error_class/recovered/at）。⇒ 它交给**进程默认 logger**，不自己开文件。

**4.3 谁落（默认 logger → 盘）**
- 装默认 logger＝`cmd/wisp/logsink.go:144`（`func installLogSink(dataDir string) (*logSink, error)`）：`:145-147` 空 dataDir 直接拒（逐字 `"wisp: no data dir resolved, so the log sink has nowhere to write"`）⇒ `:148` `dir := logSinkDir(dataDir)`（`:87`，`:76` 逐字 `const logDirName = "logs"`）⇒ `:149` `observe.InitLog(observe.LogConfig{Dir: dir, Level: logSinkLevel})` ⇒ `:153` `slog.SetDefault(slog.New(teeHandler{`。
- **tee 的两条腿**：`:157` 逐字 `primary: p.Handler(),`（`:154-156` 注释逐字说它是 "redactHandler on top of the JSON writer"）；`:160` 逐字 `mirror: slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}),`（`:158-159` 注释逐字 `The mirror deliberately carries no redaction of its own.`）。扇出＝`:208` `func (t teeHandler) Handle(...)` ⇒ `:209` primary 先、`:211` mirror 后（同一枚 record 交两边）。
- 生产装配根（尺：`grep -rn "installLogSink(" cmd --include=*.go | grep -v _test`）＝**四枚**：`cmd/wisp/run.go:222`、`cmd/wisp/resident_windows.go:66`、`cmd/wisp/secret.go:226`、`cmd/wisp/models.go:284`。⇒ **常驻 GUI 那条腿装了 sink**（`resident_windows.go:66`），这枚是"盘上有没有栈"的开关，本程判它**装了**。
- 脱敏闸在落盘之前：`internal/observe/logging.go:88` `red: Redactor{RedactPaths: cfg.RedactPaths}`；`:123` `type redactHandler struct` ⇒ `:132` `func (h *redactHandler) Handle(ctx context.Context, rec slog.Record) error` 逐 attr 过 `h.red.Attr(...)`（文件侧与 stderr 侧共用这一枚流水，mirror 那一侧不过它，见 §4.3 上面那两条腿）⇒ `internal/observe/redact.go:107` `func (r Redactor) Attr(key string, v slog.Value) slog.Attr` 的分支顺序逐枚核过：`:112` 密钥名（`secretWords` ⇒ `secret.RedactSecret`）、`:114` 音频、`:116` 正文、**:120 逐字 `return slog.String(key, r.String(v.String()))`（`slog.KindString` 那一支，`stack` 走的就是这一支）** ⇒ `:133` `func (r Redactor) String(s string) string`：`:138-139` `if r.RedactPaths { s = pathRe.ReplaceAllString(s, "<path>") }` ⇒ **`:141` 逐字 `return truncate(s, MaxLoggedString)`**。
- ★ **栈进文件之前被截到 512 rune**：`:36-37` 逐字 `// MaxLoggedString bounds any single string that reaches the log (rule 4).` ＋ `MaxLoggedString = 512`；`truncate`（`:193`）在 `:198` 逐字 `return string(runes[:max]) + fmt.Sprintf("...(truncated, %d chars total)", len(runes))`——**所以盘上那枚 `stack` 字段＝栈头 512 rune＋一枚可数的总长**，完整栈只在不过脱敏的 stderr 那半。
- ⚠ 顺带一枚**影响"可复制现场"含金量**的现读：`cmd/wisp/logsink.go:149` 那次 `InitLog` **没填 `RedactPaths`** ⇒ 走零值；而配置里 `internal/config/schema.go:509` 逐字 `RedactPaths bool \`toml:"redact_paths" default:"false"\``，且本仓自己的名册早已写明没有任何产码从 `cfg.Privacy` 读它（`cmd/wisp/config_readers_255.go:134` 逐字 `"privacy": hotClaimNoReader + "internal/observe/logging.go:50 [RedactPaths] - the mirror field exists and is filled by callers, never from cfg.Privacy"`）。⇒ **今天盘上那 512 rune 的栈头里的路径是原样路径**（这条对 AC#6 是双刃：可复制性↑、脱敏承诺✗，归编排者与票 264 那条轴一起摆）。
- 产码 `recover()` 全名册（尺：`grep -rn "recover()" cmd internal tools --include=*.go | grep -v _test`＝**9 枚**）：`internal/ball/hotkey_reload.go:146`、`internal/memory/writer.go:144`、`internal/observe/goroutine.go:296`、`internal/risk/assessor.go:278`、`internal/risk/assessor.go:292`、`internal/tools/bridge.go:823`、`internal/tools/grant.go:70`、`internal/tools/mode.go:40` ＋ `cmd/balldebug/main.go:401` 那句**注释里的**"故意不加 recover"（尺命中含它，枚数按产码实体算＝**8 枚真 recover**）。⇒ **`cmd/wisp` 产码零枚 recover**（⚠ 修正前版那句"7 处"）。

## 5. 自救："可复制现场"与"自救指令"两半各自的出口有无

（本节答：现场格式的可复制边界＋自救指令出口＋断"栈真的落盘"的既有尺名册）

未判

## 6. 既有文案有没有写歪（"禁用插件重启"那一形状）

（本节答：现读确认有没有既有文案/注释已写成"能禁用插件重启"，有则逐字引）

未判

## 7. 我判不动／量不到

（本节答：七格具名＋归口；⛔ 没有一枚写成"应该没问题"）

1. **页面那一侧今天自不自存草稿、失败回填会不会吃掉新字**——射程外（硬约束：⛔ 不读 `frontend/**`；票 167 Status 逐字 `skipped=frontend(owner-delegated)`）。⇒ 归 owner 委托的那支前端腿。**这一格决定 AC#5 到底"缺 Go 侧出口"还是"整枚归页面"，本程判不了。**
2. **512 rune 那道闸实际留下几帧、够不够定位一次崩溃**——要一发真 panic 落盘再读 JSONL 才知道；本轮禁跑任何 `go` ⇒〔仅读码〕。⇒ 归落地腿（且落地腿要先解 §7-6 那一寸人工批准，别先跑再问）。
3. **main goroutine 与非 registry 裸协程崩溃有没有 Windows／WER 兜底留痕**——产码 `recover()` 九枚全在 `internal/*`（名册见 §4），`cmd/wisp` 产码**零枚 recover**（尺：`grep -rn "recover()" cmd internal tools --include=*.go | grep -v _test`＝9 行，逐枚已列）；运行时面读不到。⇒ 归人工／落地腿。
4. **这台机器 `%APPDATA%\wisp\logs\` 里现存 jsonl 是否已含 `stack` 字段**——那是用户数据目录，不在授权搜索根内（根只有 `cmd internal tools docs .scratch scripts`），且不能跑进程 ⇒ **判不到**（⛔ 不许有人把它写成"盘上没有"）。
5. **给草稿开一枚出向字段算不算越出票 145 已批那一寸**——票面 AC#7 逐字「`internal/panel/**` 的写面**只在票 145 已批的局部解冻范围内**（composer/pump/panel_pump，且"只到新增字段"为止）」；"新增字段"这五个字能不能盖住"新增一枚承载用户原文的字段"是**裁定面**。⇒ 归编排者（需先读那枚解冻的边界行）。
6. **要不要放宽 `MaxLoggedString`／要不要另开一枚不走脱敏管线的崩溃现场文件**——`internal/observe/redact.go:36` 注释逐字 `bounds any single string that reaches the log (rule 4)` ⇒ 动它＝动既有脱敏规矩＝**人工批准面**（且乙形会新造一条"不掩码的落盘面"，要先回答"栈里带不带密钥"）。⇒ 归编排者摆给机主，⛔ 本程不裁形。
7. **草稿要不要落库**——落库＝改存储形状（`internal/memory/schema.go` 那九枚建表是 SPEC-02 镜像面）；台账已把它列为越界候选（`docs/reports/pending-and-issues.md:10629` 逐字「**草稿落库＝`docs/specs/SPEC-02` 镜像面**（改存储形状要人工批准）」，⛔ 那句不是我量的，是引的）。⇒ 归人工批准。

## 8. 交件判语

（本节答：尺名册证明零 `go`；票面/AC 框/产码/测试零改动；没读 `frontend`、没读 `design`；三枚冻结件未读未引）

**8.1 本程跑过的命令名册（逐枚，证明零 `go`）**
- `date "+%Y-%m-%dT%H:%M:%S%z"`、`git rev-parse --short HEAD`、`git branch --show-current`、`git status --porcelain`（起手与收口）、`git log --oneline -- <路径>`、`git ls-files <路径>`、`git add -- <显式 pathspec>`、`git commit -F <txt> -- <显式 pathspec>`。
- `ls`（`internal/panel`／`internal/agent`／`internal/session`／`internal/watchdog`／`internal/plugin`／`internal`／`internal/store`（不存在，尺回"No such file or directory"）／`.scratch/wisp/probes/167*`／`cmd/wisp` 目录名册）、`wc -l`（票面、前版件、四枚入向文件）。
- `grep -rn`（含 `-i`／`--include=*.go`／`--include=*_test.go`／`--include=*.md`）与 `grep -n` 定位；`sed -n 'A,Bp'` 展开上下文；计数尺一律 `;` 串接，**每一处"零枚"都在正文里翻译成"没找到＝这一格今天不存在"**，不把 `grep -c` 的 rc=1 当失败。
- **`go` 系命令零枚**：没有 `go test`／`go build`／`go vet`／`go run`／`go list`／`gofmt`／`wisp slo`／任何 scripts 门禁。全部读数＝〔仅读码〕／〔尺〕两级，无"跑一遍看看"。

**8.2 零改动自证**
- 产码／测试目录：`git status --porcelain -- cmd internal tools scripts docs` 收口读数见 §8.5；起手那一次同尺＝**零行**。
- 票面：`.scratch/wisp/issues/167-*.md` 一枚未动，AC 框（AC#1–AC#7 的 `[ ]`）一枚未碰；本件也不代它翻勾。
- 本件名下唯一写面＝`.scratch/wisp/probes/167/a2/census.md`＋`.scratch/commit-msg-167a2-*.txt`（临时件按 `issues/README` 规则 8 **只建不删**）。
- Git 纪律：只 commit、**未 push**；每枚 commit 都带显式 pathspec；`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`git add -A`／`git add .` **零枚使用**。
- ⛔ 三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）**未读未引**：包级 `grep` 顺带命中过 `l2_grant_boundary_test.go` 的行号，本件一律不转写（只在 §3 记一句"它存在、本程未读"）。
- `frontend/**`／`design/**`：**零读取**（porcelain 里出现它们的删除项不算读取，本程未 `cat`／未 `sed`／未 `grep` 进那两枚根）。

**8.3 与前轮读数的关系（不重跑、不覆盖、不冒充）**
- `167-c1/c2`（占用／序号／停止）三格：**未重跑**。本件只在 §1.3 登记 `TaskRowView`/`TaskRosterSection` 的定义位置一枚，且明写"其读数属 c2"。
- 前一版同路径件（commit `8fa4246d`，编排者代落）：**本版取代**，其全部 `file:line` 一律重新自跑（`bridge.go:91/:98`、`composer.go:235/:295/:296/:308`、`goroutine.go:153/:167/:174/:236/:240/:241/:286/:294/:296/:302/:312/:314/:315`、`redact.go:36/:37/:107/:120/:133/:141/:193`、`logsink.go:76/:87/:99/:144/:148/:149/:153/:157/:160/:208`、`schema.go` 九枚建表、`dao_tasklog.go:24`、`journal.go:122/:129`、`doctor.go:290/:296`、`diagnostics.go:62` 逐枚核过）；**三处比前版更准或不同**，具名：① 前版只数"入向文本载体几枚"，本件 §1.3 另量了 `OutgoingMessage` 的**产码构造者零枚**这一维；② 前版把 512 那格写成"〔待验，需跑一遍才知道〕"，本件按硬约束不跑，改判为 §7-2〔仅读码〕；③ 本件新量 **`logSink.logDir()` 产码调用者零枚、唯一读者是测试 `cmd/wisp/logsink_test.go:158`**，并新量 **`observe.BuildDiagnosticsBundle` 产码调用者零枚**（前版未提后者，它正是"可复制现场"那一格的既有半成品）。

**8.4 交付形状**
- 九节齐（§0–§8）；本节以下无一处 `未判`／`待填`／`填写中`；每节首行给"本节答"。
- 判语口径：AC#5 草稿＝**无**（三层皆空：字段／库表／入向听众）；AC#6 崩溃自救＝**有但没接**（记录通路在、话说出口那一寸没线，且"完整现场"卡在 512 那道闸＝人工批准面）。

**8.5 收口读数**：本节末尾由最后一枚 Edit 落名册（起手 HEAD `f761a017`；骨架 `aa016492`；§0–§1 `9c290d55`；本节与 §2–§6 的逐节 commit 号由 `git log --oneline -- .scratch/wisp/probes/167/a2/census.md` 现量后填在 §8.5 表内）。
