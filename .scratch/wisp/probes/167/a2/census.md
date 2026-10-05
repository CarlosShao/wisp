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

（本节答：存盘与否＋失败路径上谁写谁清＋"新写的字不许被覆盖"今天由谁负责）

未判

## 3. 草稿：既有尺名册

（本节答：断"草稿不吞"的既有测试＝逐枚用例名＋到底断什么；零枚就写零枚）

未判

## 4. 自救：谁产／谁投／谁落三处齐读

（本节答：panic 记录的生产者／投递者／落盘者三处具名 file:line）

未判

## 5. 自救："可复制现场"与"自救指令"两半各自的出口有无

（本节答：现场格式的可复制边界＋自救指令出口＋断"栈真的落盘"的既有尺名册）

未判

## 6. 既有文案有没有写歪（"禁用插件重启"那一形状）

（本节答：现读确认有没有既有文案/注释已写成"能禁用插件重启"，有则逐字引）

未判

## 7. 我判不动／量不到

（本节答：具名＋归口）

未判

## 8. 交件判语

（本节答：跑过的命令名册，须证明零 `go`；票面/AC 框/产码零改动；没读 frontend 与 design）

未判
