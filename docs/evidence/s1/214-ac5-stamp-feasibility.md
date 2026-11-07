# 214 AC#5「C25 盖戳」可行性取证（只读腿 · 零产码）

- 腿：`214-ac5-stamp-feasibility`；时刻 2026-11-07 18:5x；取证时 HEAD `2b3bc6839`（214 交付件＝`fc76d5d7e` 在库）。
- 起手三核：①票面 AC#5 现量在——`.scratch/wisp/issues/214-composer-plus-menu-attachments-and-context-wired.md:41` 逐字「AC#5：来源戳（C25 现成名）盖上 ⇒ 附件内容进上下文时被当作不可信外部文本处理（票 200 那节"说明文件≠授权"同一条纪律）」；②本件所有引文均现跑 `sed`/`grep` 对过文件与行（行号以取证时 HEAD 为准）；③红绿无涉——本腿零产码、零测试改动。
- ⚠ 又见票面现量腐烂一例：台账 `docs/reports/pending-and-issues.md:14264` 写「`panel_inbound.go:277` 逐字 `Message: nil`」——**行号已过期，现量 :352**（214 那笔在 :318-358 一带插了 attachment 装配段，把 Message 行推下去了）。事实本身仍真，行号引用者须复跑。

## 1. 现状：message 门今天是什么形状（现跑读码）

| 事实 | 读数 | 尺 |
|---|---|---|
| 接口定义 | `MessageRequestHandler` 一方法：`HandleMessageRequest(ctx, req ComposerRequest) error` | `internal/panel/composer_dispatch.go:92-94` |
| 插座字段 | `Message MessageRequestHandler`，注释逐字「Message answers the message request. Untouched by this slice.」 | `internal/panel/composer_dispatch.go:149-150` |
| 派发分支 | `case MethodMessageSend:` → `d.Message == nil` 时 `unattached(req)`＝按名拒绝（`ErrNoHandlerAttached`），不静默 | `internal/panel/composer_dispatch.go:215-219` |
| 装配点 | 逐字 `Message: nil`，同行注释「message = ticket 35, still refused by name when it arrives」 | `cmd/wisp/panel_inbound.go:352`（注释 :347） |
| 谁该填它（代码注释的说法） | 「workspace/message handlers are ticket 186 / ticket 35's land」 | `cmd/wisp/panel_inbound.go:67-69` |
| 谁该填它（**现读票 35 的真身**） | 票 35＝PanelBridge (C17) **传输层**：AC 八枚框（:42-:63）全是白名单 fuzz／correlationId 路由／伪造 allow 拒／resync 无状态／背压／secret 扫描／传输协议——**没有一格覆盖 `panel.message.send` 的 handler 语义** | `.scratch/wisp/issues/35-panel-bridge-c17.md` |
| handler 产码实现数 | **0**（台账 ：14264 同钉；本腿现跑 grep 复认） | `grep -rn "HandleMessageRequest" internal/ cmd/ --include=*.go \| grep -v _test` ⇒ 仅接口声明与派发分支 |
| `OutgoingMessage` 生产构造数 | **0**（台账 ：14264 同钉） | `grep -rn "OutgoingMessage{" internal/ cmd/ \| grep -v _test` ⇒ 0 |
| 附件进 agent 上下文的唯一形状 | `OutgoingMessage.ForAgent(dir)`：文本＋逐附件渲染 `[附件] name=… mime=… kind=… bytes=… path=…`＋拒件逐字重复 | `internal/panel/composer.go:363-410`（渲染体 :379 起） |
| `ForAgent` 生产调用者 | **0 枚**（只有定义、注释提及 `display_scrub_216.go:222`、测试） | `grep -rn "ForAgent" internal/ cmd/ --include=*.go \| grep -v _test` ⇒ composer.go 定义＋display_scrub_216.go 注释 |
| 附件行的读者（已就位） | `Attachments`/`AttachmentReason` 两枚 reader＋`attachmentDraftLedger` | `internal/panel/pump.go:235/:244`；`cmd/wisp/panel_attachment_214.go:66` |

**一句话现状**：插座、名册、派发分支、附件草稿账本、渲染函数全在；唯独**没有人实现 handler、没有人构造 OutgoingMessage、没有人调 ForAgent**——附件内容今天在任何一条生产路径上都到不了 agent 上下文，盖戳自然无宿主。

## 2. C25 契约原文摘引（现读）

- **C 表定义**（`docs/PLAN.md:1381`）：
  > **C25** | **`Provenance`** | 工具结果的来源标记（`{tool, origin, sensitive}`）+ 出网前的**粗粒度污染匹配**（规范化后 ≥8 字符连续片段命中即升 L2）。D30① 组合闸门的实现契约 | D33/F4
- **D33/F4 修复模型**（`docs/PLAN.md:2477-2486`）：「**敏感源**（读取即打 `sensitive` 标记）……**判定**：外泄参数中若出现来自敏感源的 **≥8 字符连续片段（规范化后）** → 升 L2 强确认，并在确认卡上**明示『本次操作包含来自 `<源>` 的内容』**」。PLAN.md 原文列的敏感源（fs.read/search.content/clipboard.read/sysinfo 窗口标题/web.fetch/doc.read/screen.capture）**不含面板附件**——附件按 AC#5 的纪律（票 200「说明文件≠授权」）同为不可信外部文本，这是 214 票面自己点的推广方向。
- **产码现状**：引擎在——`internal/risk/provenance.go`（`Mark`/`CheckText`/`Hit`；敏感源名册 ：84-102 九枚：fs.read、search.content、clipboard.read、system.get、web.fetch、doc.read、screen.capture、asr.transcript、task.output）；**装配根已有实例**——`cmd/wisp/run.go:330` `c25 *risk.Provenance`、`:1051` 构造、`:1095`/:1116 下发。⚠ 九枚名册**没有附件源**：给附件盖戳需要在 provenance.go 增一枚源名（代码名册，非冻结 C 号文字，但属敏感源语义面，按纪律须具名登记，⛔ 顺手加）。
- **票 200 先例**（盖「当作不可信文本」的戳长什么样）：`internal/panel/instructions_200.go:63` 逐字「Source is the C25 provenance name stamped on the content ("fs.read")」——戳＝把内容用源名 Mark 进引擎，使后续外泄通道的 ≥8 字符片段匹配能命中并升 L2。

## 3. 三选一裁定：**乙**（必须先落消息门 handler 本体，⌛且今天工单池没有 owner）

**丙排除**：AC#5 逐字在票 214 判据 ：41（本件起手三核①已引）——C25 盖戳**确在 214 契约里**，票面没有记错。

**甲排除（三条独立理由，缺一即可否）**：

1. **「给 message 门接上 socket」不是一枚 method value 的事。** handler 收到 `panel.message.send` 后要：解码信封文本 → 从 draft ledger 取附件行 → 构造 `OutgoingMessage` → 调 `ForAgent` → 把渲染结果喂进 agent 对话（任务管线、LLM 调用、流式回复、上下文预算 D15）。这是整条**面板→agent 对话入口**，超出 214 票面自定射程「接线＋补最小缺口」（票 ：3）不止一个量级；214 剩余的活不存在「只补一枚 nil」的形。
2. **盖戳的时点语义钉死宿主在消息门。** C25 的 `Mark` 是 scope 键控（`provenance.go:394` 一带注释：never-opened scope 也记录 taint，但匹配按 scope 走）；ingest 时附件躺在 draft 里、**没有 task scope 可挂**——在 ingest 处抢盖＝一笔挂错 scope 的死 taint，买不到 AC#5「进上下文时被当作不可信外部文本处理」的保证。有 scope 的时点＝对话发起那一跳＝消息门 handler 内（构造 OutgoingMessage 与喂 agent 之间）。
3. **本波写面互斥。** cmd/wisp 唯一写腿＝188（run.go）、internal/panel＝187——即便裁甲，214 在本波也没有可落的文件。（这条是时运不是法理，法理在前两条。）

**乙的具名（如实，含一个坏消息）**：等「消息门 handler 本体」落地——即 `internal/panel/composer_dispatch.go:150` 的 `Message` 字段从 nil 换成产码实现。⚠ **归属现状**：`panel_inbound.go:67-69`/:347 注释点名「ticket 35's land」，但**现读票 35 的 AC 八枚框全是传输层，没有一格是 handler 本体**；工单池 grep（composer/message/chat/send 族）也没有第二张票认领 ⇒ **今天没有一张现成票的现成格可等**。须编排者立新票（或按规矩给 35 具名增格——改票面判据文字须走审批）。

**给后续票的正面凭据（取证腿顺手量到的）**：盖戳引擎不必新建——`run.go:330` 的 `rt.c25` 已在装配根活着；附件行读者（pump.go 两枚 reader）与 ledger（`panel_attachment_214.go:66`）214 已备好；`ForAgent` 渲染形状（含拒件逐字重复、注明「未进入本轮」）票 214 已交付并有命名测试钉住（`Test214IngestRefRenderingCarriesTheArtifactPathAndRefusalsToTheAgent`）。消息门票的净新增面＝handler 本体＋盖戳一跳＋回执路径。

## 4. 若（未来）消息门落地时，C25 盖戳的精确落点

> 以下是**只读取证腿读码给出**的落点，⛔ 未经实现者起手复算，不构成派单；实现腿起手须按三核纪律逐条复跑。

1. **handler 实现**：新件 `cmd/wisp/panel_message_*.go`（与 `panel_attachment_214.go` 同族形状），实现 `panel.MessageRequestHandler`；payload 解码复用 `bridge.go` 的信封形状。
2. **装配点（唯一一处）**：`cmd/wisp/panel_inbound.go:352` `Message: nil` → 真 handler。`newComposerDispatchChain`（:269）是 CLI 与 resident 共享的构造体（:344-358 一个 struct literal），**一处装配两条缝都活**；盖戳引擎经参数传进 handler（`rt.c25` 同款传法，勿让 handler 自建引擎）。
3. **盖戳点**：handler 内，`OutgoingMessage` 构造 → `ForAgent(dir)` 渲染后、**喂 agent 之前**：对每枚 `Stored` 附件的内容 `c25.Mark(taskScope, <附件源名>, <内容>, <origin=artifact 路径>)`；拒件不盖（它没进上下文，ForAgent 只逐字重复拒因）。
4. **敏感源名**：`internal/risk/provenance.go:84-102` 名册现无附件源，需增一枚（如 `SrcPanelAttachment = "panel.attachment"`）并进 ：100-102 的聚合表；这是新敏感源语义面，证据件与票面都要具名登记。
5. **判据**：AC#5 的尺届时才可写死——「盖了戳的附件文本，经外泄通道参数（notify/clipboard.write/TTS/fs.write 同步盘…）携带 ≥8 字符片段 ⇒ R4 升 L2 且确认卡明示来源」；现量 C25 的通道尺在 `internal/risk/tts_gate_test.go` 一族有现成形状可仿。

## 5. 门禁登记（如实）

- **Go 改动：零** ⇒ **无对象可扫**（gofumpt／d22scan／go test 均无对象；本腿未跑测试、未设环境变量——无测试可跑）。
- 本腿写面：本件＋票面 Progress log 一行；`frontend/**`、`design/**`、`internal/**`、`cmd/**`、`docs/PLAN.md`、`docs/specs/**`、thresholds/golden 零字节。
- 行尾：本件 LF（写入后现跑核过）。
