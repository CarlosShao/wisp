# 票 453：消息门 handler 本体 ＋ C25 附件源（交付证据）

- 写腿：453（唯一写面：`cmd/wisp/**`、`internal/panel/**`（只读）、`internal/agent/**`（只读）、`internal/risk/provenance.go` 单文件＋自带测试、`docs/evidence/s1/453-*`）
- 基线：`203fc8898`（`dev`；本腿开工时 HEAD，工作树另有别腿的未跟踪/已改文件；本腿的 commit 走显式 pathspec）
- 代码提交：`7b1b90ad6`（10 文件；`git show --name-only` 的 `docs/` 条数 = 0，`docs/specs/` = 0）
- 前序：票 214 只做了收件半边；本票做出口半边
- 提交分两次：代码（`7b1b90ad6`）与证据（本文件）分开，故代码提交的 `docs/` 为零（§7）

---

## 1. AC#1：出口不存在——先写测试、先跑红灯、红句原文

### 1.1 先跑的测试（实现前，同一枚测试，后由 §2 转绿）

`cmd/wisp/panel_message_453_windows_test.go:185`
`Test453PanelMessageCarriesTheAcceptedAttachmentIntoTheModelsNextTurn`

它走的是**生产链**，不是测试自搭的替身：`newResidentPanelManager(dir, nil)`（常驻 boot 的同一个构造函数）→ 生产收件门 `panel.attachment.add`（真 broker、真 artifacts 树、真账本）→ 生产消息门 `panel.message.send`。

运行形（载具同款：`go test -c` → 三枚 DLL 钉在二进制旁 → `cwd=cmd/wisp` → 五枚语音夹具环境变量；本机无 Git Bash，按 `scripts/wisp-cli-staged-tests.sh` 的同一形在 PowerShell 复现）：

```
go test -c -o .scratch/wisp/probes/453/ac1/wisp.test.exe ./cmd/wisp/
copy third_party/sherpa-onnx/{onnxruntime,sherpa-onnx-c-api,sherpa-onnx-cxx-api}.dll <stage>/
cd cmd/wisp
<stage>/wisp.test.exe -test.count=1 -test.timeout=5m -test.v \
  -test.run=^Test453PanelMessageCarriesTheAcceptedAttachmentIntoTheModelsNextTurn$
```

**红句原文（逐字，未改写）**：

```
wisp: panel: ATTACHMENT-ADD request="r-453-ingest" actor="resident-panel" name="shot-453.png" stored=true deduplicated=false bytes=12 artifact="attachment-1b56b50ac4e976f4.png" err=<nil>
wisp: panel: PREVIEW-QUOTE request="r-453-send" actor="resident-panel" artifact="" textRunes=5 outcome=refused err=panel: 该方法在名册内，但本机未接入处理器: "panel.message.send" 的派发通路属票 35，本机未接入；这条请求里也没有预览引用（path 不是以 "/__artifact/" 开头的预览地址），因此本机既没有送出消息也没有盖戳
    panel_message_453_windows_test.go:113: panel.message.send has no exit: the accepted attachment "shot-453.png" (artifact "attachment-1b56b50ac4e976f4.png") cannot reach the model, so the next turn's input never sees it: panel: 该方法在名册内，但本机未接入处理器: "panel.message.send" 的派发通路属票 35，本机未接入；这条请求里也没有预览引用（path 不是以 "/__artifact/" 开头的预览地址），因此本机既没有送出消息也没有盖戳
--- FAIL: Test453PanelMessageCarriesTheAcceptedAttachmentIntoTheModelsNextTurn (0.17s)
FAIL
```

**这条红灯同时钉住两件事**：

1. **收件是真的**：`stored=true`，`artifact="attachment-1b56b50ac4e976f4.png"`，`bytes=12`——上一半（票 214）确实把字节写进了 artifacts 树。
2. **出口是零**：紧接着的 `panel.message.send` 被 192 预览门按名拒绝，模型/下一轮 turn 的输入里**看不到**这个文件。原句里 `"本机既没有送出消息也没有盖戳"` 就是票面那句"收进附件后，模型/下一轮 turn 的输入里看不到这个文件"的机器读数。

测试断言的是**目的地行为**（消息被接受、模型请求里出现附件），不是"断言没有出口"；后者是空真、A1880 的装饰性证据。实现前它红，实现后它绿——同一枚测试，见 §2.3。

日志保存：`.scratch/wisp/probes/453/ac1/ac1-before-red.txt`（实现前）与 `ac1-after-green.txt`（实现后）。上文 `:113` 是**实现前那一版测试文件**的行号（当时文件只有这一枚用例）；定稿后同一条断言在 `cmd/wisp/panel_message_453_windows_test.go:219`，红句内容不变。

---

## 2. AC#2：`OutgoingMessage` 生产构造点 ＋ `ForAgent` 生产调用点 ＋ 行为牙

### 2.1 生产构造点（file:line）

| 项 | 实现前（票 214 收件 @ `docs/evidence/s1/214-ac2-ac6-receipt.md:145-147`） | 本票后 |
|---|---|---|
| `panel.OutgoingMessage{` 生产构造 | **0** | **1**：`cmd/wisp/panel_message_453.go:131` |
| `ForAgent(` 生产调用 | **0** | **1**：`cmd/wisp/panel_message_453.go:132` |

两行相邻，形式是：

```go
msg := panel.OutgoingMessage{Text: req.Text, Attachments: rows}
rendered := msg.ForAgent(d.drafts.artifactsDir453())
```

`rows` 不是页面递来的行，是**账本按 ID 复核过**的行（`storedSelection453`，`cmd/wisp/panel_message_453.go:168`）：页面能回显账本行，不能凭空造行。

### 2.2 行为牙：删掉 `ForAgent` 就红（变异实测）

变异 M-B：保留 `OutgoingMessage` 构造，把 `rendered := msg.ForAgent(...)` 改成 `rendered := req.Text`（模拟"构造了却不渲染"）。重建载具后：

```
wisp: panel: MESSAGE-SEND request="r-453-send" actor="resident-panel" attachments=1 textRunes=5 outcome=accepted err=<nil>
    panel_message_453_windows_test.go:231: the model's next turn input does not contain "[附件]".
        The input was:
        看看这张图
## scene
...
--- FAIL: Test453PanelMessageCarriesTheAcceptedAttachmentIntoTheModelsNextTurn (0.17s)
```

即：`OutgoingMessage`/`ForAgent` 不是"源码里出现了"就算接线，而是**模型请求里真的会出现 `[附件] ... path=...`**；删掉这一跳，`Test453PanelMessageCarriesTheAcceptedAttachmentIntoTheModelsNextTurn` 立刻红。这满足 A1880 的"行为牙"，不是盘上字面钉。

### 2.3 转绿读数

```
=== RUN   Test453PanelMessageCarriesTheAcceptedAttachmentIntoTheModelsNextTurn
--- PASS: Test453PanelMessageCarriesTheAcceptedAttachmentIntoTheModelsNextTurn (0.13s)
```

同一测试断言模型请求文本同时含：`[附件]`、`name=shot-453.png`、`path=<artifactsDir>/attachment-1b56b50ac4e976f4.png`、用户原话 `看看这张图`。

---

## 3. AC#3：C25 附件源

### 3.1 名册加了什么

- 常量：`risk.SrcPanelAttachment = "panel.attachment"`（`internal/risk/provenance.go:106`）
- 名册：`sensitiveSourceTools` 增加一项（`internal/risk/provenance.go:116`）
- **同一次改动**把名册注释改成真话（`internal/risk/provenance.go:109-112`）：

  > `sensitiveSourceTools is SPEC-06 §5's list PLUS SrcPanelAttachment (ticket 453: the spec sentence for that one is Q-456, pending the owner; ...)`

  **不改 `docs/specs/**`**：`SPEC-06-security-gatekeeping.md:72-74` 的敏感源枚举是契约原文，新增一句属契约文字，已由编排者立 `Q-456`；本腿只在代码侧加名，不改契约。

- 加入的这一项**属于 `SPEC-06-security-gatekeeping.md:72-74` 那一串敏感源清单**（`fs.read` · `search.content` · `clipboard.read` · `sysinfo` 当前窗口标题 · `web.fetch` · `doc.read` · `screen.capture`）：`panel.attachment` 与 `doc.read` 同类——字节来自模型 provenance 之外，进入任务上下文即污染。

### 3.2 行为牙：`IsSensitiveSource` 必须是 true；删名册项就红

`internal/risk/provenance_453_test.go:20`
`Test453PanelAttachmentIsOnTheSensitiveSourceRoster` 断言 `IsSensitiveSource("panel.attachment") == true`。

变异 M-C：从 `sensitiveSourceTools` 删掉 `SrcPanelAttachment`，红句原文：

```
=== RUN   Test453PanelAttachmentIsOnTheSensitiveSourceRoster
    provenance_453_test.go:22: IsSensitiveSource("panel.attachment") = false: a panel attachment is outside content (ticket 214 AC#5) and the roster must say so
--- FAIL: Test453PanelAttachmentIsOnTheSensitiveSourceRoster (0.00s)
```

第二条 `Test453PanelAttachmentMarksAreReportedAsThatSource` 在删名册项后**仍绿**，这是 C25 既有的 fail-closed 形状（非名册 Mark 被接受但记日志，`provenance.go:109-112` 写明），不是本票的漏洞；它钉的是"Mark 用的源名 = 外泄判定报的源名"（改源名即红）。

### 3.3 C25 盖戳点在环回 mint 的 task id 下

链路：

1. 消息门把选中的账本行交给任务入口：`panelMessageDoor453.HandleMessageRequest`（`cmd/wisp/panel_message_453.go:102`）
2. 任务入口 `SubmitPanelMessage453`（`cmd/wisp/resident_task_source_windows.go:525`）先把行**暂存**到 runtime，再 `submitTask(rendered)`（单槽）
3. 文本循环在 `RunAsync` 里 mint task id；`agentRuntime.admitTask` 是这个 id 第一次出现的地方：`rt.markPanelMessage453(taskID)`（`cmd/wisp/run.go:1826`）
4. 盖戳实现：`markPanelMessage453`（`cmd/wisp/panel_message_453.go:255`），源名用 `risk.SrcPanelAttachment`，origin 是 `artifact:<裸文件名>`，内容是"显示名＋artifact 坐标"

**行为牙**：`cmd/wisp/panel_message_453_windows_test.go:272` 用最小 runtime（真 `tools.Bridge`＋真 `approval.Gate`＋真 `risk.Provenance`）走 `admitTask`，然后 `Inspect` 必须报 `SrcTool == panel.attachment`。变异 M-A（删掉 `admitTask` 里那行 mark）红句：

```
    panel_message_453_windows_test.go:289: no taint after a panel attachment reached task "task-453-c25": the attachment's name can leave unmarked
--- FAIL: Test453PanelMessageMarksTheAttachmentsUnderTheTaskScope (0.00s)
```

### 3.4 为什么盖"名字＋artifact 坐标"而不是文件字节

- 附件白名单只有图片/视频（`internal/panel/attachments.go:119-134`），PNG 字节不是文本，C25 的 ≥8 字符连续片段匹配对二进制没有意义；
- 真正会外泄的是**模型上下文里的附件表示**：`[附件] name=shot-453.png ... path=<dir>/attachment-....png`；
- 因此 Mark 的内容是 `display name + "\n" + artifact`，origin 是 `artifact:<artifact>`；这正好是同一份 `ForAgent` 渲染里出现的两个词。文件字节本身另由 `fs.read` 的既有敏感源覆盖（工具读它时）。

残留（明写）：若模型只复述文件内容（不经过 `fs.read`），本戳不命中字节。这是本票的边界：出口证明"附件进了上下文"，不把"字节级污点"伪装成已做。

---

## 4. AC#4：界面那条链的契约（前端零改动）

前端归机主外部 agent，本腿**只读契约、零改动**。现读的两处（`frontend/` 工作树内已有别腿改动，本腿未触碰）：

| 前端 | 行 | 发什么 | 后端接收 |
|---|---|---|---|
| `frontend/src/lib/panel.ts` | `:718` `submitAttachment` | `sendAttachmentBytes(file.name, file.type, byteLength, base64)` | `panel.attachment.add`：`:709-714` 四个字段 |
| `frontend/src/lib/panel.ts` | `:736-738` `sendMessage(text, attachments)` | `panel.message.send` `{text, attachments}` | `ComposerRequest.Attachments []AttachmentRef`（`internal/panel/bridge.go:140`） |
| `frontend/src/components/composer-row.tsx` | `:226-233` | 只把 `stored` 的行放进 `sendMessage` | 门按 ID 回账本复核，非 stored/伪造行拒绝 |

结论：**契约不需要改前端**。页面在收件后拿到的是快照里的 `ComposerAttachment`（`panel.ts:95-109`，与 Go 的 `AttachmentRef` 同字段），发送时原样回显；本票让这些回显行第一次真的到达模型。端到端测试用的就是这两个前端的**逐字信封形状**（`attachmentEnvelope453` / `messageEnvelope453`）。

契约一句话（给外部 agent）：`panel.attachment.add {name, declaredMime, sizeBytes, dataBase64}` 收件 → 快照回 `stored=true` 的 `ComposerAttachment` → `panel.message.send {text, attachments:[这些行]}` 出口；后端只认账本里的 `id`。

---

## 5. AC#5：门禁

### 5.1 单包门禁

| 门禁 | 命令 | 读数 |
|---|---|---|
| gofumpt | `gofumpt -l <10 个改动文件>` | **空输出，rc=0** |
| gofmt | `gofmt -l <同 10 个文件>` | **空输出，rc=0** |
| go vet | `go vet ./cmd/wisp/ ./internal/risk/` | **rc=0** |
| 全仓构建 | `go build ./...` | **rc=0** |
| `internal/risk` 全包 | `go test ./internal/risk/ -count=1` | **ok 13.184s** |
| d22scan（自建二进制，`-root` 仓根） | `<stage>/d22scan-453.exe -root <repo>` | **rc=0 `clean - no D22 ban violations`**；`bans #1-5 internal/=405 cmd/=83`、`ban #6 frontend/=123`、`ban #7 internal/tools/=51`、`ban #10 internal/risk/=21`、`ban #8 design/=39 frontend/=123 internal/=1020 cmd/=297` |
| 路径长度（本机无 POSIX sh，按 `scripts/check-path-length-budget.sh` 同式复算） | 名册 68 条 vs 超限 68 条 | **GREEN**：`roster=68 over_budget=68 not_rostered=0 stale_roster=0`；本腿 5 条新路径最长 42 字符、叶名最长 33，均未越两顶帽子（总长 121 / 叶名 100） |

### 5.2 目标测试（载具，`-test.run '^Test453'`）

```
--- PASS: Test453PanelMessageCarriesTheAcceptedAttachmentIntoTheModelsNextTurn (0.28s)
--- PASS: Test453WithoutARelayTheMessageIsRefusedByNameAndRunsNoModel (0.14s)
--- PASS: Test453AForgedRowIsRefusedBeforeAnyRendering (0.15s)
--- PASS: Test453PanelMessageMarksTheAttachmentsUnderTheTaskScope (0.00s)
--- PASS: Test453BindingTheRelayWalksTheResidentChain (0.15s)
PASS
```

回归子集（同载具，`Test453|Test192|Test214|TestTicket114r2|Test323|TestAC1OtherDoors...`）：**31/31 PASS**，含被本票改过断言的 `Test192ResidentAssemblyWiresThePreviewStack`（现在读 `*panelMessageDoor453` 且 `door.preview.src` 与窗口同源，比原来只看 192 门更强）。

### 5.3 整包 `go test ./cmd/wisp`（A1882 三读数齐全）

载具同 §1.1（`final.test.exe`，三枚 DLL 钉在二进制旁，`cwd=cmd/wisp`，五枚语音夹具环境变量），`-test.count=1 -test.timeout=90m -test.v`，全包一次一发。完整日志：`.scratch/wisp/probes/453/ac1/full-package-453.txt`（1,327,613 字节，7,878 行）。

| A1882 读数 | 值 |
|---|---|
| ① 终点锚 `TestZZ315A1` 命中/未命中 | **命中**：三枚全跑（`TestZZ315A1RootLoopFirstBridgeCall` · `TestZZ315A1SiblingChildFirstBridgeCall` · `TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask`）；最后一枚 **PASS**（0.21s） |
| ② 最后一条 `=== RUN` | `TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask`（日志第 7,705 行附近；文件共 7,878 行） |
| ③ 本波新测试在截断点之前/之后 | **之前**：`Test453...` 首条 `=== RUN` 在第 2,887 行，终点锚在第 7,705 行；本次整包**没有被截断**（锚命中即证），故不存在"新测试落在截断点之后"的假绿 |

计数：**PASS=639，FAIL=2，SKIP=17**。两条 FAIL **都是本机已知红、与本票无关**：

```
--- FAIL: TestAC4FocusReturnToPriorWindowGap33r5 (4.87s)
--- FAIL: TestAC246DevLegIgnoresTheTestTaskInjection (37.49s)
```

本波五枚 `Test453*` 在整包里全绿：

```
--- PASS: Test453PanelMessageCarriesTheAcceptedAttachmentIntoTheModelsNextTurn (0.15s)
--- PASS: Test453WithoutARelayTheMessageIsRefusedByNameAndRunsNoModel (0.14s)
--- PASS: Test453AForgedRowIsRefusedBeforeAnyRendering (0.15s)
--- PASS: Test453PanelMessageMarksTheAttachmentsUnderTheTaskScope (0.00s)
--- PASS: Test453BindingTheRelayWalksTheResidentChain (0.15s)
```

已知三红里的第三枚 `Test320FetchLegWritesTheKeyNowhereOnDisk` **本次整包通过**（它红不红取决于 `%TEMP%` 扫描干扰，单跑绿是常态）；本文件按实际读数记，不把它算成本票修好。

A1882 的假绿分支（被 native AV 静默终止、无 SKIP）**本次未发生**：进程跑到终点锚并打印了每一条 `--- PASS/FAIL/SKIP`，不存在"没跑到"的读法。

---

## 6. 本票新增/改动的文件

新增：

- `cmd/wisp/panel_message_453.go`（门＋账本选择＋runtime 暂存＋C25 盖戳）
- `cmd/wisp/panel_message_453_windows.go`（boot 晚绑定）
- `cmd/wisp/panel_message_453_windows_test.go`（5 枚行为测试）
- `internal/risk/provenance_453_test.go`（名册＋引擎两枚）

改动：

- `cmd/wisp/panel_inbound.go`（链上把 192 门包进 453 门）
- `cmd/wisp/run.go`（`pendingMarks453` 字段；`admitTask` 盖戳；`execute` 清理）
- `cmd/wisp/resident_task_source_windows.go`（`SubmitPanelMessage453`）
- `cmd/wisp/resident_windows.go`（boot 绑定一行）
- `cmd/wisp/panel_preview_192_windows_test.go`（同源断言改为复合门，未放宽）
- `internal/risk/provenance.go`（常量＋名册一项＋注释改真话）

**未改**：`docs/specs/**`（0）、`internal/panel/**`（0）、`internal/agent/**`（0）、`frontend/**`（0）、`design/**`（0）、`thresholds.go`/golden/`allowlist.txt`/`docs/PLAN.md`（0）。

---

## 7. 复跑与提交

复跑：见 §1.1 的载具步骤，`-test.run '^Test453'`；风险包 `go test ./internal/risk/ -run '^Test453' -count=1`。

提交（两次，显式 pathspec，不 push）：

1. **代码提交**：`cmd/wisp/panel_message_453.go`、`cmd/wisp/panel_message_453_windows.go`、`cmd/wisp/panel_message_453_windows_test.go`、`cmd/wisp/panel_inbound.go`、`cmd/wisp/run.go`、`cmd/wisp/resident_task_source_windows.go`、`cmd/wisp/resident_windows.go`、`cmd/wisp/panel_preview_192_windows_test.go`、`internal/risk/provenance.go`、`internal/risk/provenance_453_test.go`
2. **证据提交**：`docs/evidence/s1/453-attachment-exit.md`

`git show --name-only 7b1b90ad6` 的 `docs/` 条数：**0**（读数逐字见交付回话；证据文件在第二个提交里，故代码提交的 docs/ 为零，`docs/specs/**` 亦为零）。

---

## 8. 残留（明写，不隐藏）

1. **C25 戳的是上下文的附件表示，不是字节**（§3.4）。
2. **单槽暂存**：`pendingMarks453` 只存一条消息的行。常驻腿一次只接一发（`submitTask` 的 `running` 闸），所以不存在两条消息并发争用；若将来放开并发，这里必须改成按 task id 的队列，不能继续用单槽。
3. **CLI 缝（`panel inbound`）无任务管线**：453 门在那里 relay 永远为 nil，明文消息按名拒绝，与实现前同形；这不是遗漏，是那条腿没有模型可送。
4. **`admitTask` 的 mark 只在文本循环的 AdmitTask 钩子上**：任何绕过该钩子、自带 task id 的桥调用都不会带上附件戳——与 `internal/tools/bridge.go:989-997` 已登记的"无主腿"同一类，本票没有扩这个口。
