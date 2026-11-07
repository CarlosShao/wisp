# 214 — 附件接线证据件（腿 214-attach-wire，2026-11-07）

本票把附件受理器（`internal/panel/attachments.go`，票 92 建、票 145 钉「零调用者」）接上生产。
写面：`cmd/wisp/**`（本波唯一 cmd/wisp 腿）＋ `internal/panel/pump.go`/`composer.go` 附件两枚字段相关行
＋ `internal/panel/composer_test.go` 那条 145 用例。⛔ 未碰 `internal/ball/**`、`frontend/**`、`design/**`、
thresholds/golden、票面 AC 勾与 Status。

## 0. 起手三核（现跑读数，非转述）

| 核 | 读数 | 结论 |
|---|---|---|
| ①票面点名的串/行还在不在 | `panel_inbound.go:331` `Attachment: nil` **在**；`pump.go` 的 `NewComposerState(..., nil, ...)` 现量在 **`pump.go:277`**（票面写 `:226`，行号已腐烂，事实仍在） | 成立，行号按现量修正 |
| ②引文出处 | `composer.go` 附件两枚字段＝`Attachments`(`json:"attachments"`) 与 `AttachmentReason`(`json:"attachmentError"`)；票面写 `:238-241`，现量 `:281/:288`（行号腐烂，字段仍在） | 成立 |
| ③红绿归属 | 起手时 `grep NewAttachmentBroker\|.Ingest(` 非测试调用者 **0 枚**；145 腿钉子 `TestComposerAttachmentsStayEmptyUntilTheIngestHandlerExists` 在 `composer_test.go:968` | 成立；ⓘ 波内另有 194 只读腿飞 `internal/ball/**`、176 飞 `internal/tools/**`，本腿零触碰 |

## 1. 最小完整形状的判定（铁律 ⓪⓪⓪）

**按主流 harness 的形状选：信封 → 受理器 → 落盘 → 快照可读**，四段全接，缺一段就是 145 那记过的
「字段在、值永远不来」（A408 文案在、控件不在的线侧同形）。因此：

- **门（panel_inbound.go）**：`ComposerDispatch.Attachment` 从 nil 换成真 handler——AC#1 的那枚调用者。
- **机器（panel_attachment_214.go）**：新 handler 走现成 `DecodeAttachmentPayload` → `AttachmentBroker.Ingest`，
  **sink 是真 `memory.Store` 的 artifacts 树**（票 76/79 的名字守卫、O_EXCL、DACL、配额全适用），不造第二落盘家。
  开库是**懒**的（首枚附件才 `memory.Open`）：一条从不收附件的链零新 OS 资源，存不开拒那一发入向而不是炸整链装配。
- **账（attachmentDraftLedger）**：每次尝试（收与拒）记一行，带上限 64 行、人话 reason；
  当前拒句＝最近一次尝试的拒因、被后来成功清空（composer.go 字段注释里写明的裁定）。
- **读（pump.go）**：`PumpSources` 加 `Attachments`/`AttachmentReason` 两枚 reader——与 200/197/220 三票
  同款「carrier + 装配根 reader」形状；装配根 `run.go` 字面量挂两枚 method value。
- **时序（票 186 holder 同款）**：面板链装配早于任务管线 ⇒ ledger 在面板装配处创建、随
  `newResidentPanelManager` 返回值穿出、`startResidentTaskSource` 之后 `bindAttachmentDrafts` 绑上 runtime
  （resident_windows.go，与 `setPanelSnapshotHost` 同一行相邻）。nil 容忍：CLI seam（`panel-inbound`）无快照，
  传 nil ledger——摄取与审计照常，仅无快照行。

**明确不做**（票面 09-29 定案＋AC#6）：不新增配置键（上限用编译默认 `MaxAttachmentBytes`，调键是 D36 契约面）；
不做选文件器、@引用、message 门（前端归外部 agent；message 是票 35 的 socket，仍按名拒）；
不做「限死这件事的界面文案」（票面已具名留给后续票）。

## 2. 145 钉子的更新——为什么这个改动是对的

`TestComposerAttachmentsStayEmptyUntilTheIngestHandlerExists` 钉的是**「没接线」这个状态**。接线后它的前提
（"the only door that could produce real rows ... is refused by name"）失效，红是正确的红。更新为
`TestComposerAttachmentsTravelFromTheIngestDoorToTheWire`，**没有为绿删断言**：

- 保留（在接线后仍然为真）：裸 pump 仍出 `[]`/`""` 不出 null/编造句（臂 1/4）；全 reader 无附件 reader 仍空
  （臂 2）；构造器吃行原样传（臂 5，原臂 3）。
- 新增（更强调言）：有 reader 时行与拒句**逐字**上 wire，含 `"stored":true`、artifact 名、拒因全文（臂 3）；
  reader 返回 nil 仍渲染 `[]` 不渲染 null（臂 4）。
- 注释里写明（composer.go 两枚字段注释同批更新）：恒空从「有意状态」变成「无 reader 装配的诚实空」。

## 3. 读数表

### AC#1 受理器有生产调用者
- 尺（改前）：`grep -rn --include=*.go "NewAttachmentBroker\|\.Ingest(" internal/ cmd/ | grep -v "_test.go"` ⇒ 仅定义 0 调用。
- 尺（改后，现跑）：`cmd/wisp/panel_attachment_214.go:204`（`panel.NewAttachmentBroker(...)`）＋
  `cmd/wisp/panel_attachment_214.go:243`（`broker.Ingest(ctx, src)`）；handler 挂在
  `cmd/wisp/panel_inbound.go:351` `Attachment: attachment`——**在 cmd/wisp 的装配路径上，不在测试里**。0 → 2 处。

### AC#2 超限拒＋快照可见；限内收
- `Test214OversizeAttachmentIsRefusedBeforeAnyRead`：`MaxAttachmentBytes+1` 声称大小 ⇒
  `ErrAttachmentRejected` +「超过附件上限」，`Open` 恒 nil（未读一字节），artifacts 目录 0 文件，wire 带拒句。**PASS**
- `Test214AttachmentAddEndToEndStoresBytesInArtifacts`：真链真库，11 字节 PNG ⇒ artifacts 恰 1 文件、
  名合 `attachment-[0-9a-f]{16}.png`、字节逐位相同、ledger 行与快照同说。**PASS**

### AC#3 伪装文件按嗅探判
- `Test214DisguisedExecutableIsRefusedByItsBytes`：`MZ` 头改名 `.png` ⇒ 拒因「可执行文件（MS-DOS/PE 头 "MZ"）」，
  0 落盘。**PASS**（`matchesISOBaseMedia` 语义：ftyp＋白名单 brand，语义在 attachments.go:140-149，未动一字）

### AC#4 落盘位置在授权树内
- sink＝`memory.Store.ArtifactsDir()`（`<数据根>/artifacts`，open.go:258），名是内容哈希裸名，
  `memory.ValidArtifactName` 在 broker 内与写时同守卫——**目录不是本票选的，是票 76/79 的树**。
- `Test214PathSpellingDisplayNameIsRefusedAndNothingLands`：`../evil.png` ⇒ 拒因「裸文件名」，0 落盘。**PASS**

### AC#5 来源戳（C25 现成名）
- 现量：C25 的来源戳发生在**消息门**把文本交进上下文时；message socket 今天仍 nil（票 35 的地，
  `panel_inbound.go` 头注与本票禁区都点名）。本票交付的是它的**前半**：`ForAgent` 的渲染形状
  （`Test214IngestRefRenderingCarriesTheArtifactPathAndRefusalsToTheAgent`，PASS）——存字节带 artifact 路径、
  拒件逐字重复、明确「未进入本轮」；消息门接线时 C25 盖戳是其职责（票 200「说明文件≠授权」同纪律）。
- **如实记**：AC#5 的 C25 盖戳半边**未完成**，因为它的宿主（message 门）不存在——具名留给 message 接线票。

### AC#6 三件没接的
- 界面点不了：**是**（frontend 零写面，page 侧无选择器）。没有选择器：**是**（Go 侧也无一枚，载荷走 base64）。
  没有 @引用：**是**（无第二条文件读入口，@引用解析不存在）。逐条留在票面，本票未 pretend。

### 测试与门禁
| 门 | 读数 |
|---|---|
| `gofumpt -l`（12 件点名） | 空（rc0，两件曾要格式化已 `-w` 落） |
| `go vet ./internal/panel ./cmd/wisp` | rc0 |
| `go build ./internal/panel ./cmd/wisp` | rc0 |
| `go test ./internal/panel -count=1` | 本票新钉全绿；红 2 枚＝**既有红**，见下节归属 |
| `go test ./cmd/wisp -run 'Test214|TestComposer' -count=1 -timeout 30m` | **全绿**（8/8：214×8 中 8 PASS） |
| 波及族 `TestTicket114*|Test323*|TestAC9*|TestR9*|TestAC1OtherDoors*` | ok（323 两件未跟踪在飞件盘上现态也过） |
| d22scan | clean（rc0） |
| LF | 12 件全 `w/lf`（`git ls-files --eol` ＋ awk 复核新件） |

### internal/panel 两枚红的归属（诚实上报）
1. `TestComposerContractTypesMatchFrontend`：Go `ComposerState` 多发 `files` 键（票 190 载体，HEAD 已在）、
   page 侧 interface 未声明。**本腿 composer.go 零 json 行改动**（`git diff | grep json:` 空）；
   page 侧声明是 Q-51 名下前端腿的写面。红在改动前后同在。
2. `TestC21DesignTokensFourWayAgree`：`docs/evidence/s1/c21-native-tokens.md` 点名 `Palette.GlassRing/GlassHi/GlassHiSoft`
   在 `internal/ball/tokens.go` 无色面。`git status` 示 `internal/ball/` 5 件在飞改动（wave62 腿写面）。
   与本票写面（panel/cmd-wisp 附件）零交集。

## 4. 提交面（显式 pathspec，只 commit 不 push）

- 修改：`internal/panel/pump.go`、`internal/panel/composer.go`、`internal/panel/composer_test.go`、
  `cmd/wisp/panel_inbound.go`、`cmd/wisp/panel_resident_windows.go`、`cmd/wisp/resident_windows.go`、
  `cmd/wisp/run.go`、`cmd/wisp/panel_mode_confirm_114_windows_test.go`（签名适配，唯一已跟踪调用点）
- 新增：`cmd/wisp/panel_attachment_214.go`、`cmd/wisp/panel_attachment_214_test.go`、本证据件、票面 Progress 一行
- ⛔ 不入：`cmd/wisp/panel_resident_323_windows_test.go`、`cmd/wisp/panel_resident_assembly_323_test.go`
  （**别腿的未跟踪在飞件**；我的单行签名适配只落在盘上工作态，未 add 未 commit——HEAD 内已跟踪调用方
  `panel_mode_confirm_114_windows_test.go` 已适配，committed 树可编译）
