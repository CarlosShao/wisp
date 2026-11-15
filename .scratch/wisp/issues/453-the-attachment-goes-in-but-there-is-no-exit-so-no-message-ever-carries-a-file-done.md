# 453 消息门 handler 本体 ＋ C25 附件源（票 214 AC#5 的移交目标）

> 立票 2026-11-15 15:4x，由编排者立。来源＝票 214 AC#2–#6 收件（`f563a3c0e`）的 AC#5 判语，以及编排者的两条裁决。

**Status**: ready-for-agent
**Blocked by**: 无（可独立开工）
**Blocks**: 票 214 的 AC#5 已按"移交本票"收；本票是那半边的**唯一去处**

## 为什么有这张票

票 214 收了附件的**半路**：面板侧收得进、落盘落到真 artifacts 树、有行为牙。但**出口那一半不存在**：

- `OutgoingMessage{` 生产**构造 0 处**（编排者现跑复核）
- `ForAgent(` 在产码里只有 `internal/panel/composer.go:495` 一处 —— **是定义不是调用者**（"调用者 0"）
- **C25 名册里没有附件源**

⇒ 结果：文件收进来了，**没有任何路径把它送进对话/模型**。渲染形状测试有，**生产链路零**。这是 A1880 的极端形：**连字面钉都只有渲染那一枚**。

## AC

- [x] **AC#1 先证出口不存在**：写一枚测试断言"收进附件后，模型/下一轮 turn 的输入里**看不到**这个文件"——**它必须先绿或先红由你判**，红句原文抄进证据，⛔ 只写"没有出口"这种形容。
- [x] **AC#2 `OutgoingMessage` 有生产构造点 ≥1**：给 `file:line`；并给**会红的行为测试**（A1880：⛔ 只让断言源码字符串的测试红）。
- [x] **AC#3 C25 名册加附件源**：名册改动的 `file:line`；`ForAgent` 有**真调用者**（≥1，`file:line`）。
- [x] **AC#4 界面那条链走通**：`frontend/src/lib/panel.ts:718` 的 `panel.attachment.add`（2026-11-11 构建产物里已有该路由）→ 收件 → 出口。前端归机主外部 agent，你只把**契约**写清。
- [x] **AC#5 门禁**：gofumpt 空／`go vet` rc=0／d22scan 0 违规／path-length GREEN；`go test ./cmd/wisp` 走载具一次一发，**整包读数必须带终点锚是否命中＋最后一条 `=== RUN`**（台账 A1882）。

## 约束

- ⛔ 动 `thresholds.go`／golden／`allowlist.txt`／`docs/PLAN.md`／`docs/specs/**`／D43 表／`frontend/**`／`design/**`。
- ⛔ 改 D36 语义（附件上限沿用编译默认）。
- 本包 `cmd/wisp` 每波只一枚写腿。
---

## 收件（2026-11-15 17:4x，`7b1b90ad6` ＋ 证据 `7bd8b469e`）：**5/5，结案**

**尺1**：两枚提交在；`7b1b90ad6` **10 件**（`cmd/wisp/panel_message_453.go`／`panel_message_453_windows.go`／测试／`panel_preview_192_windows_test.go`／`panel_inbound.go`／`resident_task_source_windows.go`／`resident_windows.go`／`run.go`／`internal/risk/provenance.go`／`provenance_453_test.go`）；★**`docs/specs` 命中 0** —— 我核过，边界守住。
**尺2**：`docs/evidence/s1/453-attachment-exit.md`（8 节：AC#1 红／AC#2 两生产点＋行为牙／AC#3 名册／AC#4 前端契约／AC#5 门禁／文件清单／复跑／残留）。
**尺3 我自己复跑**：`Test453PanelMessageCarriesTheAcceptedAttachmentIntoTheModelsNextTurn` ⇒ **PASS 0.16s**；`./internal/risk -run 453` ⇒ **2 枚全 PASS**（`…IsOnTheSensitiveSourceRoster` 0.00s／`…MarksAreReportedAsThatSource` 0.00s）；门禁我另跑：gofumpt **0 条目**／`go vet ./cmd/wisp ./internal/risk` **rc=0**／`d22scan` **clean rc=0**（scope `internal/=405`／`cmd/=83`／`frontend/=12`）／path-length **GREEN**。

### AC#1 红句（实现前，逐字）
`panel.message.send` 被 **192 预览门按名拒绝**：
> `wisp: panel: PREVIEW-QUOTE request="r-453-send" … outcome=refused err=panel: 该方法在名册内，但本机未接入处理器: "panel.message.send" 的派发通路属票 35，本机未接入；这条请求里也没有预览引用…因此本机既没有送出消息也没有盖戳`
> `panel_message_453_windows_test.go:113: panel.message.send has no exit: the accepted attachment "shot-453.png" (artifact "attachment-1b56b50ac4e976f4.png") cannot reach the model, so the next turn's input never sees it`
⇒ 断言的是**目的地行为**（模型请求里出现附件），⛔ 不是"断言没有出口"。

### AC#2 两个生产点 ＋ 行为牙
| 项 | 前 | 后 |
|---|---|---|
| `OutgoingMessage{` 生产构造 | 0 | **1**：`cmd/wisp/panel_message_453.go:131` |
| `ForAgent(` 生产调用 | 0 | **1**：`:132` |

**变异 M-B**（保留构造、把 `ForAgent` 换成 `req.Text`）⇒ 红句「the model's next turn input does not contain "[附件]". The input was: 看看这张图」⇒ **删掉这一跳测试立刻失败**，⛔ 不是"源码里出现就算接线"（A1880）。转绿后同一测试断言模型请求含 `[附件]`、`name=shot-453.png`、`path=<artifactsDir>/…png`、用户原话。

### AC#3 C25 附件源（我给的甲裁定，四条边界逐条对上）
- `risk.SrcPanelAttachment = "panel.attachment"`（`provenance.go:106`）入册 `:116` ✔
- **注释同批改成真话**（我现读）：`sensitiveSourceTools is SPEC-06 §5's list **PLUS** SrcPanelAttachment (ticket 453: the spec sentence for that one is Q-456, pending the owner; …)` ✔ —— ⛔ 没留假话
- **`docs/specs/**` 零改动**，并把要加的位置写成 `SPEC-06:72-74`（＝我立的 `Q-456`）✔
- 行为牙：`Test453PanelAttachmentIsOnTheSensitiveSourceRoster` 删名册项即红，红句「a panel attachment is outside content (ticket 214 AC#5) and the roster must say so」✔

### AC#4 前端契约（⛔ 前端零改动，只写清）
现读三条：`panel.ts:718 submitAttachment` 发 `{name, declaredMime, sizeBytes, dataBase64}` → `:736-738 sendMessage(text, attachments)` → `panel.message.send {text, attachments}`；`composer-row.tsx:226-233` 只把 `stored` 的行放进去。
⇒ **结论：契约不需要改前端**。给外部 agent 的一句话：`panel.attachment.add` 收件 → 快照回 `stored=true` 的 `ComposerAttachment` → 原样回显进 `sendMessage`；本票让这些回显行**第一次真的到达模型**。

### AC#5 门禁（含 A1882 三读数）
**整包**：终点锚 `TestZZ315A1` **命中**（最后一枚 PASS）；最后一条 `=== RUN` ＝ `TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask`；`PASS=639 FAIL=2 SKIP=17`，两枚 FAIL 是本机已知红 ⇒ **本发无截断点**，新测试在锚之前。

### 残留（腿明写四条，我认，不并进票面）
1. **C25 戳的是上下文的附件表示，不是字节**（§3.4 —— 这是对的：字节盖戳会把整段二进制搬进污染面）
2. **单槽暂存** `pendingMarks453` 只存一条消息的行；常驻腿一次只接一发（`submitTask` 的 `running` 闸）⇒ 若将来放开并发，**必须改成按 task id 的队列**
3. **CLI 缝（`panel inbound`）无任务管线**：453 门在那里 relay 永远为 nil，按名拒绝 —— 与实现前同形，不是遗漏
4. **`admitTask` 的 mark 只在文本循环的 `AdmitTask` 钩子上**：绕过该钩子、自带 task id 的桥调用不会带附件戳（与 `internal/tools/bridge.go:989-997` 已登同形）
