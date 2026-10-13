# 票 145 · A 节 —— 快照今天（299 面板宿主落地后）逐件现状（只读腿 `145-a2`）

- 锚 = dev `8e4f6f83`（本腿起手 HEAD，`git rev-parse HEAD` 现跑）。
- 本腿禁跑任何 `go` 命令；结论只来自 Read / `grep -n`（射程内）。跑才知道的格子逐枚登记在末节。
- 每一行 = 结论 + 尺 + 件行号快照。行号随树变动会腐烂，引用前请现跑同一把尺。

---

## A1 `internal/panel/composer.go` 的 `Snapshot` 字段逐枚（现读 `composer.go:57-92`）

**直接字段 = 6 枚**（不是本票标题那句"恰四字段"；`A560`/`145-c2` 已更正为 6，本腿复量一致）。

| # | Go 字段 | json 键 | 行 | 生产者 / 读口（现读） |
|---|---|---|---|---|
| 1 | `Pending []ApprovalCardView` | `pending` | `:58` | 读口 `cmd/wisp/panel_pump.go:58 liveVerdicts` -> 填块 `internal/panel/pump.go:239-245` |
| 2 | `Results []ResultChunk` | `results` | `:59` | 读口 `run.go:711 rt.stream.Chunks`（无中间 reader 函数）-> 填块 `pump.go:247-250` |
| 3 | `Composer ComposerState` | `composer` | `:60` | `pump.go:269 NewComposerState` + `:273-297` 逐段覆盖 |
| 4 | `GeneratedAt string` | `generatedAt` | `:61` | `pump.go:299-303 now()`（`p.src.Now`，nil -> `time.Now`） |
| 5 | `Instructions *InstructionsSection` | `instructions,omitempty` | `:74` | 读口 `panel_pump.go:108 instructionBundle` -> 填块 `pump.go:304-309` |
| 6 | `Tasks *TaskRosterSection` | `tasks,omitempty` | `:91` | 读口 `panel_pump.go:146 taskRosterState` -> 填块 `pump.go:310-328` |

`ComposerState`（`composer.go:235-270`）**11 枚字段**：

| # | 字段 | json 键 | 行 | 填点 |
|---|---|---|---|---|
| 1 | `Mode` | `mode` | `:236` | `pump.go:269` + 不可读时 `:273 ModeUnknownView` |
| 2 | `Workspace` | `workspace` | `:237` | `pump.go:261-264`（读口 `panel_pump.go:83 workspaceView`） |
| 3 | `Attachments` | `attachments` | `:238` | `pump.go:269` 传 nil -> `composer.go:274-276` 归一为 `[]`（**恒空**，见 C） |
| 4 | `AttachmentMIMEs` | `acceptedAttachmentMimes` | `:239` | `composer.go:281 AcceptedMIMETypes()` |
| 5 | `MaxAttachmentB` | `maxAttachmentBytes` | `:240` | `pump.go:265-268`（`AttachmentMax` 未接，零 -> 常量） |
| 6 | `AttachmentReason` | `attachmentError` | `:241` | **泵从不填**（恒空串，见 C） |
| 7 | `Git` | `git` | `:246` | 读口 `panel_pump.go:228 gitView` -> `pump.go:275-279` |
| 8 | `CurrentModel` | `currentModel` | `:255` | 读口 `panel_pump.go:242 currentModel` -> `pump.go:280-286` |
| 9 | `ModelKnown` | `modelKnown` | `:256` | 同上（`= CurrentModel != ""`） |
| 10 | `Credential` | `credentialState` | `:268` | 读口 `panel_config_store.go:285 credentialStatus` -> `pump.go:287-297` |
| 11 | `CredentialKnown` | `credentialKnown` | `:269` | 同上 |

## A2 `cmd/wisp/panel_pump.go` 的填充点逐处（现读）

reader 定义（本文件内）：

| 行 | 函数 | 喂哪一段 |
|---|---|---|
| `:58` | `liveVerdicts()` | `Snapshot.pending`（经 `NativeVerdict.CardView()`） |
| `:83` | `workspaceView()` | `composer.workspace` |
| `:108` | `instructionBundle()` | `Snapshot.instructions` |
| `:146` | `taskRosterState()` | `Snapshot.tasks` |
| `:228` | `gitView()` | `composer.git` |
| `:242` | `currentModel()` | `composer.currentModel` / `modelKnown` |

- `credentialStatus` 不在本文件：`cmd/wisp/panel_config_store.go:285`（`configStore` 方法）。
- **装配根唯一** = `cmd/wisp/run.go:699-726`（`panel.NewSnapshotPump(panel.PumpSources{...})`）；
  字面量接 10 根线：`Verdicts:700 / Mode:701 / Workspace:702 / Git:703 / Model:704 / Credential:710 / Results:711 / Instructions:718 / Tasks:724 / Out:725`。
- `PumpSources` 共 13 枚字段（`pump.go:138-215`），**未接 3 枚**：`L1Windows`(:200，nil=诚实无枚举)、`AttachmentMax`(:205)、`Now`(:209)。
- 投递：`pump.go:348 Publish()` -> `json.Marshal` -> `src.Out`；生产实参 `rt.bookPanelSnapshot`（`panel_pump.go:319`）只落**有界摘要行**进 ledger（`panel_pump.go:290-313` 自述 512-rune 界装不下整包）。

## A3 页面侧（dev 树）有没有接收端 / 类型

**结论：接收端 = 0（复认 313-a1 的读数，并在 dev 树这一棵上成立）。类型 = 有，但只声明了旧 4+6 键。**

- 接收端尺：`grep -rn 'addEventListener("message"\|onmessage' frontend/src`（含 `.ts/.tsx`）= **0 命中**。
  正控（证明这把尺有牙）：同树 `addEventListener` 在 8 枚文件命中（`App.tsx:169/:393`、`records-table.tsx:609-611`、`flowchart.tsx:275`、`agent-screen.tsx:226`、`use-dismiss.ts:46-47`、`react-bits/*`）＝全是 keydown/resize/pointer，无 message。
- 唯一宿主引用 = `frontend/src/lib/panel.ts:152-155 hostBridge()`，且只用于**出向** `postMessage`（`:179`、`:218`）。`window.wispBridge` / `chrome.webview` 全树仅 `panel.ts:147/:154`。
- 挂载：`frontend/src/main.tsx:104-110` 渲染 `<App />`（**无 snapshot prop**）=> `App.tsx:91 snapshot = EMPTY`（`EMPTY` 定义 `App.tsx:70-75`）。
- 注释与实现相反：`main.tsx:12` 逐字"Go pushes snapshots through the C17 bridge"，本树无实现（与 313-a1 `d-page.md:9` 同读数）。
- **类型**（`frontend/src/lib/panel.ts`）：
  - `PanelSnapshot`（`:129-137`）= **4 顶层键**（`pending`/`results`/`composer`/`generatedAt`）=> Go 已发的 `instructions`/`tasks` **未声明**。
  - `ComposerState`（`:112-119`）= **6 键**（mode/workspace/attachments/acceptedAttachmentMimes/maxAttachmentBytes/attachmentError）=> Go 已发的 `git`/`currentModel`/`modelKnown`/`credentialState`/`credentialKnown` **未声明**。
- **dev 树实读者**（真面板路径）：`App.tsx` 读 `snapshot.pending`（`:206/:208/:321`）、`snapshot.results`（`:431/:457/:461`）、`snapshot.composer`（`:92`）+ `composer.mode.names/current`（`:181/:335`）+ `composer.workspace.set/canonical`（`:243/:303`）；`approval-screen.tsx:68`（只被 `showcase.tsx:41` 挂=非真面板路径）。
- **读数为 0 的键**（`grep -rn '\.<key>\b' frontend/src` 计数）：`generatedAt`=0、`currentModel`=0、`modelKnown`=0、`credentialState`=0、`credentialKnown`=0、`instructions`=0、`tasks`=0、`git`=0。
- **死组件**（有读取点但 0 引用）：`components/composer.tsx`（`:89/:132/:139/:199/:205` 读 attachments/mimes/maxBytes/attachmentError）= 0 引用；`components/chat-screen.tsx`（`:25` 读 results）= 0 引用（`result-stream.tsx` 只被它引用 => 传递性死）。
  => dev 树上"读 composer 附件四键"的那一族**不在真面板路径**（与 `145-f1` §1.4 的挂载读数同形，但那是另一棵树 `dsh/feat/frontend-p0-v2`；本腿在 dev 树复量）。

## 量不到（本腿禁跑 go；逐枚点名）

1. 在册三红（`TestComposerContractTypesMatchFrontend` / `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` / `TestC21DesignTokensFourWayAgree`）当前颜色——只读静态确认函数在盘。
2. 任何 `Snapshot` 出口字节数 / ledger 512-rune 界的现量。
3. 页面 `tsc` / lint 是否因"类型缺键"报错（页面那侧本腿只读未跑构建）。
