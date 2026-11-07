# 145 恒空两枚处置（attachments / attachmentError）——裁定：丙＝保留＋具名测试钉住"恒空是有意的"

- 日期：2026-11-07 13:5x–14:1x +08；性质＝写腿（选型权按派单铁律 ⓪⓪⓪ 由本程自拍并写明理由）。
- 起手锚 `df331ec09`；起手 `git status --porcelain internal/panel` ＝**空**（190 腿当时未在写，无冲突）。
- 写面：`internal/panel/composer.go`（两枚字段定义点注释）＋ `internal/panel/composer_test.go`（追加一枚测试）＋本件＋票 145 Progress log 一行。两处显式 pathspec、只 commit 不 push、全程 LF。
- 上游：取证腿 `145-ac6-field-census.md` §2 把 `Attachments`／`AttachmentReason` 列为"答不出清单"，处置归编排者；本程即该处置的执行腿。

## 1. 两枚字段的完整生命周期（本程现跑，尺可复算）

| 环节 | `ComposerState.Attachments`（json: attachments） | `ComposerState.AttachmentReason`（json: attachmentError） |
|---|---|---|
| 定义 | `internal/panel/composer.go:274`（现 `:280`，`[]AttachmentRef`） | `internal/panel/composer.go:277`（现 `:288`，`string`） |
| 生产赋值 | 唯一赋值点＝`NewComposerState`（composer.go:312→现 `:318` 起）；泵 `pump.go:277` 硬传 nil，构造器把 nil 归一为 `[]AttachmentRef{}` | **无**——构造器字面量不填、泵不填；全仓非测试命中＝定义行 1 枚（census 尺复跑同读数） |
| 非测试读取 | 仅 `display_scrub_216.go:240`（洗刷级复制，输入仍来自构造器） | **0 枚** |
| 序列化 | 恒上 wire：`"attachments":[]`（空数组非 null） | 恒上 wire：`"attachmentError":""` |
| 页面消费 | TS `ComposerState` 声明该键（合同尺 12 键 reconciled）；挂载侧 composer.tsx 0 引用（145-f1 ③档）＋Go→页出向跳 0（145-f1 ⓐ） | 同左 |
| 真源 | **不存在**：`panel.attachment.add`（bridge.go:44）的 handler 在 `cmd/wisp/panel_inbound.go:331` 显式 `Attachment: nil`（注释逐字 ticket 92 still refused by name）；`NewAttachmentBroker`／`.Ingest(` 非测试调用者 **0 枚** | 同左（refuse 的 reason 无生产通路可产生） |

## 2. 三选一裁定：**丙**（保留＋具名测试钉住）

**为什么不是甲（接真 reader）**：读口要接的前提是装配根有真源可给。现量**无源**——附件 ingest 的 handler 在生产装配根未接线（`Attachment: nil`），broker 无生产调用者，全仓没有任何代码持有"composer 草案附件列表"。接 reader＝先造一个草案附件仓库＝造源，正是本票"宁缺毋造"禁的形状；且装配根在 `cmd/wisp`（本程禁写面，190 腿在写）。

**为什么不是乙（删掉）**：①这两键不是无人认领的占位——页面契约 `frontend/src/lib/panel.ts` 的 `ComposerState` interface **声明**了它们（合同尺双向），删 Go 键会立刻在 `subtract(tsKeys, goKeys)` 方向新红 `[attachments attachmentError]`（`composer_test.go:73-78` 双向机制现读确认），在册 `files` 红之外加两条，而 `frontend/**` 是本编队零写面 ⇒ 乙**机械上落不了地**。②票 92 的 ingest handler 接线之日这两键就是真载体（`AttachmentRef` 全套 9 键已在盘、`OutgoingMessage.Attachments` 同形），删了再加回＝净负价值。③AC#6 的"答不出＝拿掉"针对的是**落地集新字段**；这两枚是基线字段＋有既定归宿（票 92），不属于装饰品族。

**为什么是丙**：与仓内既有形状同族（`GitViewNotProbed()`／`CredentialUnknown`／mode "unknown"——"没人读"是一个要显式说出的状态）。价值三笔：①AC#6 对这两枚从"答不出"变成"有具名判据"（测试名钉进定义点注释与证据件）；②回归绊线——谁在没接 ingest handler 的情况下往泵/构造器里塞值或把空态改成 null/编造文案，测试四臂即红；③零契约损伤，票 92 接线时照常走 `NewComposerState` 那扇门（测试臂 3 钉了"构造器不拦未来"）。

## 3. 落地清单

| 件 | 改动 |
|---|---|
| `internal/panel/composer.go:274-287` | 两枚字段定义点加注释（裁定、无源现量、钉子测试名）。⚠ gofumpt 要求注释块分段后重排同 struct 字段列宽：`Mode`/`Workspace` 两行与三行附件字段的**纯空白对齐**变动，无语义改动，`files` 键行零字节 |
| `internal/panel/composer_test.go:952-1018`（追加） | `TestComposerAttachmentsStayEmptyUntilTheIngestHandlerExists` 四臂：①裸泵＝空非 nil 数组＋reason 空串；②满配泵（Git/Model/Credential reader 全挂）两枚仍空（空不是"装配太穷"而是"没有 ingest 源"）；③构造器收真行原样透传（钉"这扇门不拦票 92"）；④线上字节＝字面 `[]` 与 `""`（null/编造文案即红） |

## 4. 门禁读数（改前改后各一遍；环境照纪律 unset 四件＋PATH 前缀 build/sherpa-onnx）

- `gofumpt -l internal/panel`（绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`）：改后**空**（rc=0）。
- `go vet ./internal/panel`：**净**。`go build ./internal/panel/ ./cmd/wisp/`：**净**。
- ⚠ `go build ./...` 红——**HEAD 上既有**：`build/mut_rvw.go`（windows-only package main）符号未定义；`git status --porcelain build/` ＝空 ⇒ 该目录树＝HEAD 状态，非本程引入，也非本程写面，照实登记不修。
- `go test ./internal/panel/ -count=1 -timeout 20m`：改前红名册＝`TestComposerContractTypesMatchFrontend`（红在 `files` 键，`composer_test.go:78`）＋`TestC21DesignTokensFourWayAgree`；**改后逐名同册、名数未增**（rc=1，两枚皆在册既有红；145-screen 尺 4 同色）。新用例单跑 `--- PASS`（0.0x s）。⚠ 颜色尺 `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` 今日两轮都不在红册（与 145-r2/r3 记的三枚红差一枚，照实记现量）。
- 新用例钉的两键在合同尺里仍双向 reconciled：`ComposerState <-> ComposerState: 12 JSON keys reconciled`。
- `sh scripts/d22scan.sh`：**clean rc=0**（bans #1-5 `internal/=378`＋`cmd/=58`、ban #6 `frontend/=116`、ban #7 `internal/tools/=48`、ban #10 `internal/risk/=21`、ban #8 `design/=39`／`frontend/=116`／`internal/=939`／`cmd/=224`）。

## 5. 本程没做的（如实）

- 未碰 `frontend/**`／`design/**`／`cmd/wisp/**`／`internal/ball/**`／thresholds／golden；票面 AC 勾与 Status 一字未动（Progress log 只追加）。
- 未给泵加任何读口、未改 `NewComposerState` 签名、未动 `pump.go` 一字节（`pump_test.go:124/:276` 那两枚解冻行也未碰）。
- `TestComposerCurrentModelTravelsOnlyFromItsReader`（AC#6 对象用例）未重跑本体——全量轮无它红，census/筛查腿的绿读数维持。
