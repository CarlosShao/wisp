# 145-f1 — 面板页面字段名册 ↔ Go 快照实交字段 对拉表（只读普查腿）

> 本腿射程＝**页面到底读哪几枚字段**（过去从没量过的那一半）＋**Go 到底交哪几枚**，逐枚对拉。
> 全程只读页面代码（`git show`／`git grep` 走对象层，未进 `D:/wt/fe`，未跑任何 `go test`）。
> 大输出落 `logs/145f1-*.txt`（同名尺可复跑）。凡未量的行一律写"未证"，不填假读数。

## §0 起手锚

```
git rev-parse --short HEAD                     # a781fdc8
git rev-parse --abbrev-ref HEAD                # dev
date                                           # Wed Oct  7 10:33:47 CST 2026
git rev-parse --short dsh/feat/frontend-p0-v2  # 16c2f038   ← 页面那侧的锚（缺了此表不可复现）
git ls-tree -r --name-only dsh/feat/frontend-p0-v2 -- frontend/src | wc -l   # 85（与派单给的枚数一致）
```

终态自查尺：`git status --porcelain -- cmd internal docs frontend` → **空**（起手空、交件后仍空）。

---

## §1 页面侧字段名册

### §1.0 地基与穷举尺

- 地基：`frontend/src/lib/panel.ts`（311 行，全文 `logs/145f1-panel.ts.txt`）＝`PanelSnapshot` 五枚嵌套 interface；
  `frontend/src/lib/panel-views.ts`（106 行，全文 `logs/145f1-panel-views.ts.txt`）＝九枚屏表 + `currentView()`。
- 穷举其余声明（200 命中，落 `logs/145f1-type-hunt.txt`）：
  `git grep -n -E 'interface |type [A-Za-z]+ = \{|Snapshot|Props' dsh/feat/frontend-p0-v2 -- 'frontend/src/**/*.ts' 'frontend/src/**/*.tsx'`
- 引 `@/lib/panel` 的只有 8 枚文件（`logs/145f1-panel-importers.txt`）：
  `App.tsx` · `components/approval-screen.tsx` · `components/chat-screen.tsx` · `components/composer.tsx` ·
  `components/l2-approval-card.tsx` · `components/result-stream.tsx` · `fixtures/harness.ts` · `lib/panel-views.ts`。
  ⇒ **"页面读字段"这件事只可能发生在这 8 枚里**；其余组件全是 props-only（见 §1.3）。

### §1.1 快照契约上、页面逐枚读取的字段（全名册）

尺（**整行读**，未 `cut -c1-150`；每条读数可用 `git show dsh/feat/frontend-p0-v2:<path> | grep -n '<字段名>'` 复算）：

| 归属（TS 声明行） | 字段名（TS 逐字） | 类型 | 谁在消费（文件:行） | 今天的值从哪来 |
|---|---|---|---|---|
| `PanelSnapshot:129` | `pending` | `ApprovalCardView[]` | `approval-screen.tsx:68` `const pending = snapshot.pending`；`:83` `:94` `:103`；`App.tsx:206` `:208` `:321` | Go 快照 |
| 〃 | `results` | `ResultChunkView[]` | `chat-screen.tsx:25` `const results = snapshot.results`（→`:33` 传给 `ResultStream`）；`App.tsx:431` `:457` `:461` | Go 快照 |
| 〃 | `composer` | `ComposerState` | `App.tsx:92` `const composer = snapshot.composer ?? EMPTY_COMPOSER` | Go 快照（**兜底是一枚页面写死的假快照，见 §1.5**） |
| 〃 | `generatedAt` | `string` | **读取 0 命中**（尺见 §1.2 行 1-3）。出现处只有对象字面量：`App.tsx:77`（`EMPTY.generatedAt: ""`）、`fixtures/harness.ts:97` `:181`（写死串） | **无人读**；只有占位假值 |
| `ApprovalCardView:24` | `correlationId` | `string` | `approval-screen.tsx:58` `:104`；`l2-approval-card.tsx:265` `:266`；`result-stream.tsx:34`(chunk 同名) | Go |
| 〃 | `tool` | `string` | `approval-screen.tsx:55`；`l2-approval-card.tsx:268` `:292` | Go |
| 〃 | `args` | `string[]` | `l2-approval-card.tsx:305`（`<ParamsBlock args={view.args}>`） | Go |
| 〃 | `level` | `string` | `approval-screen.tsx:53`；`l2-approval-card.tsx:269` `:293` | Go |
| 〃 | `rulesHit` | `string[]` | `approval-screen.tsx:61`；`l2-approval-card.tsx:312` | Go |
| 〃 | `reason` | `string` | `l2-approval-card.tsx:207` `:216` | Go |
| 〃 | `reasonKnown` | `boolean` | `l2-approval-card.tsx:207` `if (!view.reasonKnown \|\| …)` | Go |
| 〃 | `sessionOverrideBlocked` | `boolean` | `l2-approval-card.tsx:321` | Go |
| 〃 | `callChain` | `string[]` | `l2-approval-card.tsx:316` `:318` | Go |
| 〃 | `decidedBy` | `string` | `l2-approval-card.tsx:338` | Go |
| `ResultChunkView:122` | `correlationId` | `string` | `App.tsx:432`（turn id）`:463`（React key）；`result-stream.tsx:34` | Go |
| 〃 | `text` | `string` | `App.tsx:433` `label: chunk.text.slice(0, 18) \|\| "（空）"`；`App.tsx:464`；`result-stream.tsx:37` | Go（**空串时页面自造"（空）"**，见 §1.5） |
| 〃 | `done` | `boolean` | `result-stream.tsx:37` `{chunk.done ? chunk.text : <RevealText …/>}` | Go |
| `ComposerState:112` | `mode` | `ComposerMode` | `composer.tsx:94` | Go |
| 〃 | `workspace` | `ComposerWorkspace` | `composer.tsx:150` `:152` `:154` `:163`；`App.tsx:242`（监控栏"工作区"行） | Go |
| 〃 | `attachments` | `ComposerAttachment[]` | `composer.tsx:89` `:132` `:134` | Go |
| 〃 | `acceptedAttachmentMimes` | `string[]` | `composer.tsx:199` | Go |
| 〃 | `maxAttachmentBytes` | `number` | `composer.tsx:199` | Go |
| 〃 | `attachmentError` | `string` | `composer.tsx:139` | Go |
| `ComposerMode:70` | `current` | `string` | `composer.tsx:94` `state.mode.current` | Go |
| 〃 | `names` | `string[]` | `composer.tsx:94` | Go |
| 〃 | `l2ConfirmNames` | `string[]` | `composer.tsx:276` | Go |
| `ComposerWorkspace:80` | `set` | `boolean` | `composer.tsx:150` `:152`；`App.tsx:242` | Go |
| 〃 | `spelling` | `string` | `composer.tsx:163` | Go |
| 〃 | `canonical` | `string` | `composer.tsx:150` | Go |
| 〃 | `reparse` | `boolean` | `composer.tsx:154` | Go |
| 〃 | `rewritten` | `boolean` | **读取 0 命中**（全 `frontend/src`；尺与正控见 §1.2 行 4） | Go 交了，页面不读 |
| 〃 | `reason` | `string` | `composer.tsx:150` | Go |
| `ComposerAttachment:95` | `id` | `string` | `composer.tsx:135`（React key） | Go |
| 〃 | `name` | `string` | `composer.tsx:300` `:313` | Go |
| 〃 | `mime` | `string` | `composer.tsx:315` | Go |
| 〃 | `kind` | `string` | `composer.tsx:294` `attachmentIcon(a.kind)` | Go |
| 〃 | `sizeBytes` | `number` | `composer.tsx:315` `bytes(a.sizeBytes)` | Go |
| 〃 | `artifact` | `string` | **读取 0 命中**（全 `frontend/src`；正控同行） | Go 交了，页面不读 |
| 〃 | `stored` | `boolean` | `composer.tsx:89` `:295` | Go |
| 〃 | `deduplicated` | `boolean` | `composer.tsx:318` | Go |
| 〃 | `reason` | `string` | `composer.tsx:302` | Go |
| **`PanelSnapshot` 之外的 `view`**（`panel-views.ts:98` 的 `PanelSnapshot & { view?: string }`） | `view` | `string?` | `panel-views.ts:99` `const named = snapshot.view`；`App.tsx:123` `const hostView = (snapshot as { view?: string }).view` | **Go 从不交**：`panel-views.ts:89-93` 原文"The field does not exist on `PanelSnapshot` yet - ticket 35's pump owns it"⇒ 恒落 `DEFAULT_VIEW = "chat"`（`:80`） |

⇒ **快照契约上页面真读的叶子＝38 枚**（顶层 3 ＋ 卡 10 ＋ chunk 3 ＋ composer 段 6 ＋ mode 3 ＋ workspace 5 ＋ attachment 8）。
⇒ **页面声明了却不读的＝3 枚**（`generatedAt` / `rewritten` / `artifact`）。
⇒ **页面读了但 TS interface 里根本不存在的＝1 枚**（`view`）。

### §1.2 两把尺的正控／负控（0 读数一律带对照才成立）

| # | 尺（完整命令） | 正控读数（确认存在的字段） | 负读数（判定"页面不读"） |
|---|---|---|---|
| 1 | `git grep -n -o -E 'snapshot\??\.[a-zA-Z_]+' dsh/feat/frontend-p0-v2 -- 'frontend/src'`（`logs/145f1-snapshot-accesses.txt`） | `snapshot.pending` **4** ／`snapshot.results` **4** ／`snapshot.composer` **1** | `snapshot.generatedAt` **0** ／`snapshot.tasks` **0** ／`snapshot.instructions` **0** |
| 2 | `grep -n -E 'snapshot[?]?\.|\(snapshot as' logs/145f1-App.tsx.txt`（整文件已落盘，逐行读） | `App.tsx:92` `:206` `:321` `:431` `:457` `:461` | 同文件 `generatedAt` 读取 **0**（只在 `:77` 的字面量里出现） |
| 3 | `git grep -n -E '\btasks\b\|\binstructions\b' dsh/feat/frontend-p0-v2 -- 'frontend/src/**/*.tsx' 'frontend/src/**/*.ts'`（21 命中，`logs/145f1-tasks-instr-hits.txt`，**逐条读完**） | `panel-views.ts:30` `:69` 的 `"tasks"`（＝屏名字符串，非快照字段）；`tasks-screen.tsx:95` 的 `tasks`（＝本地 prop） | 21 命中里**无一枚**是 `snapshot.tasks` / `snapshot.instructions` 的读取 |
| 4 | `git grep -c -E '\.<字段>\b' dsh/feat/frontend-p0-v2 -- 'frontend/src'`（逐字段跑） | `.reparse` **1**（composer.tsx）／`.stored` **2**／`.deduplicated` **1** | `.rewritten` **空**／`.artifact` **空** |
| 5 | 同上尺打在新段上 | `.workspace` **4 枚文件**（`App.tsx` `composer.tsx` `harness/main.tsx` `lib/panel.ts`） | `.git` **空**／`.currentModel` **空**／`.modelKnown` **空**／`.credentialState` **空**／`.credentialKnown` **空** |

### §1.3 页面自己声明的字段名册（快照契约之外＝props-only 组件族）

这族才是"没有输入可画"的真身：**组件已在盘上、类型已声明、没有任何 Go 生产者接到它身上**。
尺：`git show dsh/feat/frontend-p0-v2:frontend/src/<path> | grep -n -A14 -E '^export (interface|type) '`（合并落 `logs/145f1-props-rosters.txt` 153 行，全文逐行读）。

| 页面类型（文件:行） | 字段名（TS 逐字） | 谁在消费（文件:行） | 今天的值从哪来 |
|---|---|---|---|
| `TaskStatus`＋`TaskRow`（`components/tasks-screen.tsx:25`/`:27`） | `id` `title` `sub` `status` `elapsed?`；status 六态枚举 `running\|approval\|blocked\|queued\|failed\|done` | `tasks-screen.tsx:95` `:98` `:101` `:104`（组件内）；挂载只在 **`showcase.tsx:290` `<TasksScreen tasks={SHOWCASE_TASKS} />`**；**`App.tsx` 不挂它** | **写死的假内容**：`fixtures/harness.ts:185 SHOWCASE_TASKS`（`:350` 等条目逐枚是常量）。`tasks-screen.tsx:12-13` 原文自己写着"feeds nothing back to the host. The tasks themselves are pure props; when Go grows a tasks feed, the snapshot becomes the …" |
| `BallStateRow`（`components/ball-screen.tsx:18`） | `state` `visual` `enter` `exit` | **`showcase.tsx:304`** `<BallScreen rings={SHOWCASE_BALL_RINGS} states={SHOWCASE_BALL_STATES}/>`；`App.tsx` 不挂 | **写死**：`fixtures/harness.ts:225 SHOWCASE_BALL_STATES`、`:274 SHOWCASE_BALL_RINGS` |
| `BallRingDemo`（`ball-screen.tsx:29`） | `label` `progress` `tone` | 同上 | **写死**（同一枚 fixture） |
| `PaletteItem`／`PaletteGroup`（`palette-screen.tsx:31`/`:38`） | `label` `icon?` `hint?`／`name` `items` | `showcase.tsx:279`（fixture）**＋真实面板 `App.tsx:368` `groups={[]} query=""`** | 面板侧＝**空数组占位**（不是假内容，是"没源就不画"）；showcase 侧＝写死 `fixtures/harness.ts:195` |
| `ContextUsage`（`context-meter.tsx:25`） | `tokens?` `cacheHitRate?` `occupancy?` | **真实面板 `App.tsx:473` `<ContextMeter usage={{}} className="mb-1.5" />`** | **空对象占位**。`context-meter.tsx:10` 原文："PanelSnapshot has no field for any of the three yet" |
| `TurnMark`（`turn-rail.tsx:38`） | `id` `label` | `App.tsx:431` 由 `snapshot.results.map` **现场推导**（`:443`-`:450` 交给 `LineSidebar`） | **不是新字段**：从 `results` 的 `correlationId`/`text.slice(0,18)` 客户端算出来的（⇒ 见 §1.5 那句"（空）"） |
| `GitBranchView`（`git-branch.tsx:32`） | `current` `local` `worktrees: {branch,path}[]` `isRepo` | **真实面板 `App.tsx:214-216`**：`<GitBranchChip view={{ current: "", isRepo: false, local: [], worktrees: [] }} onCheckout={() => undefined} onNewWorktree={() => undefined} />` | **页面写死的空假值**＋**注释与代码相反**：`:214` 注释逐字"git branch. **Reads its strings from the snapshot**; a workspace that is not a repo renders nothing at all rather than a chip." ⇒ 实参是四个字面量，**没有一个字节来自快照**。⇒ "宁缺毋造"要防的那一形，且 Go 侧**已经有** `composer.git`（§2.2） |
| `TreeSession`／`PanelSidebarProps`（`panel-sidebar.tsx:51`/`:59`） | `id` `title` `state: "done"\|"failed"\|"streaming"` `meta?`；`activeSessionId` `workspaceName` `workspacePath` `workspaceSessions` `historyGroups: {label, rows}[]` `theme` `settingsOpen` ＋5 枚回调 | **真实面板 `App.tsx:293-301`**：`historyGroups={[]}`、`activeSessionId={local.sessionId}`、其余＝本地 state/回调 | `historyGroups`＝**空数组占位**；`workspaceName`/`workspacePath` 本可来自 `composer.workspace`（页面注释 `:68` 就写着"Path comes from the snapshot's C26-resolved field"）但 `App.tsx:293` 那一段没传 ⇒ 未证是否真传（`--` 见 §5）；`session` 维 Go 侧无具名导出（尺见 §3.2） |
| `MonitorEnvRow`／`MonitorAgentRow`（`harness/monitor-popover.tsx:15`/`:20`） | `label` `value`／`id` `name` `state: "running"\|"done"\|"failed"` | **真实面板 `App.tsx:240-244`** | `env` 第一枚**读快照**（`composer.workspace.set ? composer.workspace.canonical : "尚未收到"`），第二枚**写死** `{ label: "分支", value: "尚未读到" }` ⇒ 而 Go 早就在交 `composer.git.branch`（§2.2/§3④）。`agents`＝待补（本腿未读完该块，未证） |
| `Row`（`config-screen.tsx:70`，内部 interface） | `group` `name` `desc` `blocked` `control?: "alpha"\|"motion"` | `config-screen.tsx:86-` 的 `APPEARANCE_ROWS: readonly Row[]` | **整表写死在组件里**；`blocked` 逐枚是人话（如 `font_size` 那行逐字"快照里没有这个字段"）⇒ 这一族是"页面自己抄了一份 Go 有没有源的账"，Go 变了它不知道 |
| `PanelView`（`lib/panel-views.ts:37`） | `id` `label` `icon` `demoIcon` `note?` `interim?` `fed` `selfFed?` | `panel-views.ts:65-75` 常量表；`App.tsx` 导航经 `viewOf()`/`currentView()` | **页面常量**。⚠ `fed: boolean`（`:54-56`）是页面手抄的"Go 有没有源"账：今天 `fed:true` 只有 `chat`/`approval` 两枚、7 枚 `false`（`config` 一枚另标 `selfFed:true`）⇒ **Go 加了字段这枚 `fed` 不会自己变绿**，是 §3③ 那批新段落地时最容易漏的一格 |

> 注：`harness/*`（`main.tsx` 的 `HarnessMainProps`、`fixtures/harness-app.ts` 的 17 枚 `Harness*` interface、`right-rail.tsx` 的 `TreeRow`/`ReviewFile`）是**另一条线**（owner 委托的 harness 骨架，`?harness=1` 才挂，`fixtures/harness.ts:101 HARNESS_BANNER` 自己标"HARNESS 假数据"）。本腿按"面板页面"口径未逐枚展开，列 §5；`App.tsx:311` 真挂了 `SettingsPage`、`:349` 挂 `RightRail reviewFiles={[]} fileTree={[]} terminalLines={[]} browserUrl=""`（**空占位**）。

### §1.4 页面里"写死的假内容"点名（宁缺毋造要防的东西）

1. `App.tsx:54-71 EMPTY_COMPOSER` ＋ `:73-78 EMPTY`：一整套假的 `ComposerState`/`PanelSnapshot`，作为 `snapshot` 缺席时的默认值（`:91` `App({ snapshot = EMPTY })`）。其中 `workspace.reason: "尚未收到原生侧的状态快照"` 是人话兜底，但 `mode.current: "unknown"`、`maxAttachmentBytes: 0` 会被 `composer.tsx:199` 当真值渲染。
2. `App.tsx:214-216`：`GitBranchChip` 的 `view` 四个写死字面量＋两条 `() => undefined` 回调，**注释却声称读自快照**（§1.3 已逐字引）。
3. `App.tsx:243`：监控栏第二行 `value: "尚未读到"`（Go 已有 `git.branch`）。
4. `App.tsx:433`：`chunk.text.slice(0, 18) || "（空）"`——页面替空文本造标签。
5. `fixtures/harness.ts` 的 `SHOWCASE_*`（尺：`git grep -n -E 'export const (SHOWCASE|RB_|HARNESS)[A-Z_]*' … -- 'frontend/src/fixtures'` → 22 枚常量，`fixtures/harness.ts:33/110/177/185/195/225/274/280/287/294/300/307/308/311/331/339/354/362/367/373/386`）＝**全部写死**，且 `SHOWCASE_APPROVAL_SNAPSHOT:177`、`HARNESS_SNAPSHOT:33` 是两枚**冒充 `PanelSnapshot` 的假快照**（`generatedAt: "showcase"`）。它们只进 `showcase.tsx`/`?harness=1`，不进取自真宿主的路径——但**类型上没有任何东西拦住别人把它当真件用**。
6. `config-screen.tsx:86-` `APPEARANCE_ROWS` ＋ `panel-views.ts:65-75` `fed` 表＝页面**手抄的 Go 现状账**。

---

## §2 Go 侧今天真交出去的字段（现量）

### §2.1 `Snapshot` 顶层：**六枚**，不是票 145 说的"恰四枚"

尺：`Read internal/panel/composer.go:57-93`（锚 `a781fdc8`；`type Snapshot struct` 在 **`:57`**，`func NewSnapshot` 在 **`:106`**）。

| # | Go 字段 | JSON key |
|---|---|---|
| 1 | `Pending []ApprovalCardView` | `pending` |
| 2 | `Results []ResultChunk` | `results` |
| 3 | `Composer ComposerState` | `composer` |
| 4 | `GeneratedAt string` | `generatedAt` |
| 5 | `Instructions *InstructionsSection` | `instructions`（`omitempty`，票 200） |
| 6 | `Tasks *TaskRosterSection` | `tasks`（`omitempty`，票 197） |

⇒ 票面"恰四字段"过期；`internal/panel/pump.go:5-6` 的注释"Snapshot (composer.go:44) has four fields"在今天的 HEAD 上**同样过期**（行号也漂了）；票 145 自己更正②给的 `:44-49`／`:79` 也已漂（现量 `:57`／`:106`）。

### §2.2 各段字段（Go 实交形状，逐枚）

尺：`grep -n -A22 'type <名> struct' internal/panel/{composer,approval,attachments,git,instructions_200,subagent_roster_197}.go` ＋ `grep -n 'json:"'`。

- `ApprovalCardView`（`approval.go:39-61`）10 枚：`correlationId` `tool` `args` `level` `rulesHit` `reason` `reasonKnown` `sessionOverrideBlocked` `callChain` `decidedBy` ⇒ **与 TS 逐键同**。
- `ResultChunk`（`composer.go:95-99`）3 枚：`correlationId` `text` `done` ⇒ 同。
- `ComposerState`（`composer.go:235-272`）**11 枚**：`mode` `workspace` `attachments` `acceptedAttachmentMimes` `maxAttachmentBytes` `attachmentError` ＋ **`git`**（票 181） ＋ **`currentModel`** **`modelKnown`**（票 145 AC#2b） ＋ **`credentialState`** **`credentialKnown`**（票 248）
  ⇒ TS `ComposerState` 只声明 6 枚，**差的 5 枚页面一个字都不读**（§1.2 尺 5）。
- `ModeView`（`:249`）3 ＝ `current` `names` `l2ConfirmNames`｜`WorkspaceView`（`:209-225`）6 ＝ `set` `spelling` `canonical` `reparse` `rewritten` `reason`｜`AttachmentRef`（`attachments.go:83-104`）9 ＝ `id` `name` `mime` `kind` `sizeBytes` `artifact` `stored` `deduplicated` `reason` ⇒ 三枚均与 TS 逐键同。
- `GitView`（`git.go:117-147`）**10 枚**：`kind` `reason` `branch` `detachedSha` `isDetached` `repoRoot` `currentWorktree` `worktrees` `branches` `switchBlocked`；`GitWorktree`（`git.go:93-`）＝ `path` `branch` `sha` `detached`(…)（尾部未读完，列 §5）。
- `CredentialState`＝**字符串枚举**（`config_handlers.go:175-187`：`unknown` `config_unreadable` `no_ref_declared` `all_recorded`…），不是 struct。
- `InstructionsSection`（`instructions_200.go:71-`）3 ＝ `status` `reason?` `files[]`；`ProjectInstructionFile`（`:44-65`）8 ＝ `path` `tier` `depth` `bytes` `truncatedBytes?` `dropped?` `duplicateOf?` `source?` ⇒ **共 11 叶，页面 0 读**。
- `TaskRosterSection`（`subagent_roster_197.go:150-`）6 ＝ `rows[]` `inFlightSlots` `poolCap` `streamTruncated` `streamElidedRunes` `droppedStreamKeys`；`TaskRowView`（`:106-138`）12 ＝ `taskId` `label` `kind` `parentTaskId` `status` `statusKnown` `statusReason?` `streamKey` `blockedOnApproval` `streamTruncated` `streamElidedRunes` `streamDropped` ⇒ **共 18 叶，页面 0 读**。

### §2.3 生产者（谁填真值）

尺（剔 `_test.go`）：
```
grep -rn --include=*.go -E 'NewSnapshot\(|NewComposerState\(|NewApprovalCardView\(|TaskRosterSection\{|InstructionsSection\{' internal cmd | grep -v '_test.go'
```

| 构造点 | 读数（生产码） |
|---|---|
| `NewSnapshot(` | **`internal/panel/pump.go:303`**（真泵）＋ `pump.go:236`（无源时的空形）。`cmd/wisp/` 内无直接调用者 ⇒ **票 145 更正⑤"`NewSnapshot(` 0 枚调用者"今天已过期** |
| `NewComposerState(` | `pump.go:269` |
| `NewApprovalCardView(` | 生产侧唯一调用者**仍是** `cmd/wisp/panel_assets.go:68`（命令行诊断分支）；面板路径走 `pump.go:98` `NativeVerdict.CardView()`（这条本腿未展开逐字段验，§5） |
| `InstructionsSection{` | `instructions_200.go:164` `:174` `:210`（`:210` 是 `s.Instructions =` 的赋值形） |
| `TaskRosterSection{` | `subagent_roster_197.go:199`（在 `taskRosterSectionFrom` 内） |

⚠ **"写了没接"那一族的现形位置**：`NewSnapshot` 的 return（`composer.go:131-136`）**只填四枚**（`Pending`/`Results`/`Composer`/`GeneratedAt`），第 5、6 枚不经这条构造路径，要靠别处 `snap.Instructions =` / `snap.Tasks =` 的赋值。本腿**没有**穷举出"泵到页面那一站有没有人真的做这两次赋值"（尺与读数见 §5 第 3 条），这正是票 145 收表里"新段出门是 `depth:0 / windowMs:0 / null / 全零`"那一族的落点。

### §2.4 快照怎么到达页面（最后一跳＝没有这一跳）

尺与读数（逐条可复跑）：
```
grep -rn --include=*.go 'Out:' internal cmd | grep -v _test.go     # 唯一装配点 cmd/wisp/run.go:725  Out: rt.bookPanelSnapshot
grep -rn --include=*.go -A8 'func (rt \*agentRuntime) bookPanelSnapshot' cmd/wisp  # cmd/wisp/panel_pump.go:319-325
```
- `bookPanelSnapshot` 只做三件事：`rt.lastSnap, rt.lastSnapBytes, rt.snapSeen = snap, data, true`（`:321`）＋一行日志（`:323`）。**不碰 WebView、不发 postMessage。**
- `internal/panel/pump.go:19-25`（整行读）逐字：*"WHAT THIS FILE IS NOT: the transport. There is no Go -> page channel in this tree today - no WebView2 host (ticket 33 is unclaimed), no postMessage writer, no local HTTP/SSE/websocket server … this pump stops at the byte boundary: Marshal hands out the JSON, Publish hands it to src.Out … The last mile is NOT here."*
- 反向（页→Go）那一跳是**绑定**：`cmd/wisp/panel_host_windows.go:80` `panelDispatchBinding = "wispDispatch"`（页面侧对应 `lib/panel.ts:139-155` 的 `window.wispBridge` / `window.chrome?.webview`）。
- 我在 `internal/panel/**` 与 `cmd/wisp/**` 里**没有读到任何把 `data` 送进页面的具名行**（尺：`grep -rn --include=*.go -E 'PostWebMessage|EvaluateScript|CreateWebMessageAsJson' cmd/wisp internal/panel | grep -v _test` → **0 命中**）。

⇒ 与机主给的"Go→页绑定回话只有库 `Run()` 泵那一形兑现"**方向一致但更窄**：真兑现的是"泵 marshal 到字节＋记账到 `rt.lastSnap`"；**Go→页的推送本身今天零路径**，页面 `snapshot` prop 的实际来源（除 `App.tsx:91` 的 `EMPTY` 默认值与两枚 fixture 外）在本腿射程内**取不到**（§5 第 4 条）。

---

## §3 对拉表（四档）

### §3.1 ①两边都有、已对上＝**38 枚**

顶层 3（`pending` `results` `composer`）＋ `ApprovalCardView` 10 ＋ `ResultChunkView` 3 ＋ `ComposerState` 段 6 ＋ `ComposerMode` 3 ＋ `ComposerWorkspace` 5 ＋ `ComposerAttachment` 8。
尺＝§1.1 全表（每一行都给了 `git show … | grep -n` 可复算的消费点）↔ §2.2 的 `json:` tag。

### §3.2 ②页面要、Go 没交＝**快照契约上 1 枚** ＋ **props-only 族 26 枚叶子**

| 页面要的字段 | 谁在读/会读 | Go 侧今天有没有 | 真来源（既有包能不能算出来） |
|---|---|---|---|
| `view`（当前屏） | `panel-views.ts:99`、`App.tsx:123` | **无**（TS 里也不存在，靠 `& {view?}` 强转） | 无源风险低但**形状已有名**：九枚屏名在 `panel-views.ts:26-35`；Go 侧不知道这九枚名字（票 145 AC#3 已裁"ⓐ 不落"）⇒ 归 **票 145／`Q-51`／票 35（ⓑ）**，**本腿不新开账** |
| `TaskRow.id/title/sub/status/elapsed` | `tasks-screen.tsx:95-`（只挂 showcase） | **Go 已交同族**（`tasks` 段 18 叶） | **有源且已算出**：`subagent_roster_197.go:106`/`:150`；缺的是"页面改读 `snapshot.tasks`"＋改名（→④）。`sub`/`elapsed` 两枚 **Go 无对应**＝**无源**（尺：`grep -rn --include=*.go -iE 'elapsed' internal | grep -v _test`，本腿未跑＝未证，见 §5） |
| `ContextUsage.tokens/cacheHitRate/occupancy` | `App.tsx:473`（`usage={{}}`） | **无** | **候选有源**：分子 `internal/agent/cost.go:41`（`Usage.InputTokens` 累加）、分母 `internal/config/settings.go:143 SetModelContextWindow`、窗口/预算 `internal/agent/budgets.go`。⇒ **票 167 AC#2 已排这一枚**（本腿不重开） |
| `MonitorEnvRow{label:"分支"}` 的真值 | `App.tsx:243` 写死"尚未读到" | **Go 已交**（`composer.git.branch`） | **有源**，缺页面读 ⇒ 不是"Go 加字段"，见 ③ |
| `MonitorAgentRow.id/name/state` | `App.tsx:240`-块（未读完） | **Go 已交同族**（`TaskRowView.taskId/label/status`） | **有源**（`subagent_roster_197.go:106`）＋**票 188 在册**（它的"没有状态维"读数已被票 197 部分超越，见 §4） |
| `TreeSession.id/title/state/meta` ＋ `historyGroups{label,rows}` | `App.tsx:293`（`historyGroups={[]}`） | **无** | **未证**：尺 `grep -rn --include=*.go -E 'func .*(List|Recent|Sessions)\(' internal/session` → **0 命中**（同族尺在 §3.2 的 `queue.go:148 Depth()` 上出过正读数）⇒ 写"尺取不到具名导出"，**不写成"无源"**；会话维归**票 166**（`panel-sidebar.tsx:71` 注释逐字"Empty today (ticket 166 owns the feed)"） |
| `BallStateRow.state/visual/enter/exit` | `showcase.tsx:304` | **无** | **未证**：`grep -rn --include=*.go -E 'func .*(Transitions|States)\(\)' internal/statemachine` → **0 命中**；D43 转移表在文档侧（`PLAN.md`），非导出尺 ⇒ 归 **票 145 的"ⓐ 当前屏"同族**或另登，**本腿不裁** |
| `PaletteGroup.name/items`＋`PaletteItem.label/icon/hint` | `App.tsx:368`（`groups={[]}`）、`showcase.tsx:279` | **无** | **无源（尺命中不相关）**：`grep -rn --include=*.go -iE 'palette|commandregistry' internal cmd` → 只出 `internal/ball/renderer_windows.go:74-82` 的**配色板 palette**，与命令面板无关 ⇒ 真"无源" |
| `GitBranchView.current/local/worktrees{branch,path}/isRepo` | `App.tsx:214` | **Go 已交**（`composer.git`，同段 10 叶） | **有源**，缺读＋改名（→④） |
| `TurnMark.id/label` | `App.tsx:431` | 页面自己从 `results` 推导 | **不需 Go 加字段**（若要真标签则撞 `results.text` 截断，非契约缺口） |

### §3.3 ③Go 交了、页面不读＝**8 枚具名 key ＋ 2 枚整段（29 叶）**

| Go 交的 | 行 | 页面读数 |
|---|---|---|
| `generatedAt`（顶层） | `composer.go:61` | 声明了、**0 读**（§1.2 尺 1/2）；`App.tsx:77` 只往占位对象里写 |
| `composer.git` | `composer.go:246` | **0 读**（尺 5） |
| `composer.currentModel` / `modelKnown` | `:256` / `:261` | **0 读** |
| `composer.credentialState` / `credentialKnown` | `:271` / `:272` 区 | **0 读** |
| `composer.workspace.rewritten` | `composer.go:221` | **0 读**（TS 声明了，`lib/panel.ts:89`，仍不读） |
| `composer.attachments[].artifact` | `attachments.go:96` | **0 读**（TS 声明 `lib/panel.ts:104`） |
| `instructions` 整段（11 叶） | `composer.go:62-70` | **0 读**（尺 3） |
| `tasks` 整段（18 叶） | `composer.go:71-93` | **0 读**（尺 3） |

⇒ 这一档里 **`git`/`currentModel`/`credentialState`/`tasks`/`instructions` 五枚是"Go 真算了、页面没接"**（不是装饰品：有真生产者），而 `generatedAt`/`rewritten`/`artifact` 三枚是**连页面声明了都不读的形状**。区分这两族很重要——前者只差一次读，后者是契约噪声。

### §3.4 ④名字对不上但形状像同一件事＝**3 族**（逐枚给两边原文）

| # | 页面侧原文 | Go 侧原文 | 关系／坑 |
|---|---|---|---|
| 1 | `tasks-screen.tsx:27-34` `export interface TaskRow { id: string; title: string; sub: string; status: TaskStatus; elapsed?: string; }`，`:25` `TaskStatus = "running"\|"approval"\|"blocked"\|"queued"\|"failed"\|"done"` | `subagent_roster_197.go:106-138` `type TaskRowView struct { TaskID string \`json:"taskId"\`; Label string \`json:"label"\`; Kind string; ParentTaskID string \`json:"parentTaskId"\`; Status string \`json:"status"\`; StatusKnown bool; StatusReason string; StreamKey string; BlockedOnApproval bool; StreamTruncated bool; StreamElidedRunes int; StreamDropped bool }` | `id↔taskId`、`title↔label`；**status 词汇不同源**：页面六态是人话枚举，Go 的 `Status` 是 **D43 状态机名**（`:114` 注释逐字"a D43 name"，`:120` `StatusReason` 存在的意义就是"没 filing 时不许印一个判决"）⇒ 直接改名对接会把 D43 名当人话画。页面 `sub`/`elapsed` 在 Go 侧**无对应**；Go 的 `kind/parentTaskId/streamKey/blockedOnApproval/…` 页面**无处画**（那才是监控栏要的） |
| 2 | `git-branch.tsx:32-41` `export interface GitBranchView { current: string; local: readonly string[]; worktrees: readonly { branch: string; path: string }[]; isRepo: boolean; }` | `git.go:117-147` `GitView{Kind json:"kind"; Reason; Branch json:"branch"; DetachedSha; IsDetached; RepoRoot; CurrentWorktree; Worktrees []GitWorktree; Branches []string; SwitchBlocked}`；`GitWorktree{Path; Branch; Sha; Detached}` | `current↔branch`、`local↔branches`、`worktrees{branch,path}↔GitWorktree{branch,path,sha,detached}`；**`isRepo` 不是一枚 bool**＝Go 用 `Kind`（四值 `GitKind*`）＋`Reason`＋**`SwitchBlocked`（`:145` 注释逐字："empty would mean 'switchable', which is false in this tree, so it carries the reason"）** ⇒ 页面把 `Kind` 折成 `isRepo` 就**吞掉了"不许切"这条真判决**，且 `onCheckout={() => undefined}` 已经是假可点 |
| 3 | `panel-sidebar.tsx:49` `SessionState = "done"\|"failed"\|"streaming"`；`:51-57 TreeSession{id,title,state,meta?}` | `TaskRowView{taskId,label,status(D43 名),statusKnown,statusReason}` ＋ `ResultChunk{correlationId,text,done}` | **疑似同族、未裁**：会话（sidebar 的历史）≠ 任务（roster）≠ 流（results）。`monitor-popover.tsx:20` 的 `MonitorAgentRow.state: "running"\|"done"\|"failed"` 又开了**第四套**状态词汇。⇒ 四套词汇（D43 / TaskStatus / SessionState / MonitorAgentRow）谁向谁对齐**＝本腿不裁**，只点名"这就是最容易埋 bug 的那格" |

---

## §4 与既有票的分工（本腿不裁归属、不开票、不选形）

### 票 145（`.scratch/wisp/issues/145-panel-snapshot-has-four-fields-…-carrier.md`）

- **AC#1 那张表勾过没有**＝**勾了，且勾到"表本体已由非实现者验收表定档"**的程度：票面 `[x]`，档位**附条件**（六枚 C1–C6 全为记录级、无一枚要求改码），凭据＝`docs/evidence/s1/145-snapshot-field-census-r1.md`（本腿 `ls docs/evidence/s1 | grep '^145'` 现量到 5 枚件：`145-snapshot-field-census-r1.md`／`145-snapshot-fields-landed-r1.md`／`…-r1-accept-r1.md`／`145-snapshot-growth-r2.md`／`145-snapshot-growth-r3.md`）。
  ⚠ 但那份普查是**从 Go 侧倒推**的（台账 `12944` 行原文："AC#1/AC#2 一直只从 Go 侧推……页面到底读哪几枚从没量过"）⇒ **本腿与之不重复**：它交"该加哪几枚"，本腿交"页面真读哪几枚"。§3.2/§3.3 里凡与它撞的字段（`view`／occupancy／queue）我都标了"已在册"。
- **AC#2 那句"TS 侧对齐不写 `frontend/**`，由前端会话自己落"今天还成立吗**＝**不成立**，两条现量出处：
  ① 票 145 自己 `[ ]` **维持未勾**，且收表写明"未勾不是欠账，是这一格的结论本身"（AC#2 落地集＝空）；
  ② 台账 `docs/reports/pending-and-issues.md` 的 `A383`（`:8442`，09-28"禁止任何跨会话联系前端会话"）与 `:12938`（10-07 逐字"⛔ 跨会话转达那条……继续生效"；"写页面文件……我不替他扩大解释"）⇒ **"由前端会话自己落"的收件人今天不存在**，而 `Q-51`（`:1086`，"谁有权同一枚 commit 同时写 `internal/panel/**` 与 `frontend/src/lib/panel.ts`"）**仍挂未答**、`:6339` 已把它降级为"素材没定，答了也白答"。⇒ 归属问题**原样交回机主**，本腿不选形。
- **本腿现量到的票 145 过期点**（票面"更正"节之外的第二层）：①"恰四字段"→ 现量**六枚**；②更正⑤"`NewSnapshot` 0 枚调用者"→ 现量 `pump.go:236`/`:303` 两枚生产调用者；③行号再漂（`:44-49`→`composer.go:57`，`:79`→`:106`）；④`internal/panel/pump.go:5-6` 的注释也过期（"has four fields"）。

### 票 167 / 188 / 242 与本表的重叠

| 票 | 与本腿表的重叠在哪 | 现量到的状态 |
|---|---|---|
| `167-four-small-outputs…` | **重叠＝§3.2 的 occupancy 那一枚**（`ContextUsage.tokens/cacheHitRate/occupancy`，页面 `App.tsx:473 usage={{}}`）＋ queue/stop/draft 三枚本腿**未展开量**（`snapshot` 上没有、页面也没声明具名字段，故不在 §1 名册里） | Status：`ready-for-agent，但立而不派`，AC#1 `[x]`、AC#2 `[ ]`；票面逐字"**排在票 145 之后——它俩抢同一批快照字段**"。⇒ **occupancy 已有人在排，本腿不另开** |
| `188-the-task-monitor-rail…no-go-side-source` | **重叠＝§3.4①/③ 的 `tasks` 整段与 `MonitorAgentRow`** | Status：`ready-for-agent（先只读设计核，再落地）`，AC#1 `[x]`（＝乙，子代理要进产品）、AC#2-#6 `[ ]`。⚠ **它的"现量"节今天过期**：票面写"全仓没有'子代理'这一维的任何数据模型""后台任务……没有'状态'这一维"，而本腿现量 `internal/panel/subagent_roster_197.go:106`/`:150`（`TaskRowView.Status/StatusKnown`＝D43 名）与 `Snapshot.Tasks` 已上线 ⇒ **源已存在、缺的是页面读**（§3.3）。票面 AC#4 逐字"与票 145 分账：快照载体扩张归票 145，本票只交'源'" |
| `242-the-grant-binding-layer…zero-rulers` | **不重叠**：它量的是 `approval.PanelItem` 的出向读面**不许带 grant**（AC#2 反射扫字段名枚数）与一次性令牌绑定层，**不涉及面板快照的字段名册**。唯一沾边＝它给"面板那面不带令牌"只有三行注释守着 ⇒ 本腿 §3.3 的"Go 交了页面不读"那批**不许被读成"面板能答"**（R20／票 92 口径：面板只显示＋发起请求） | Status：`已立，未派`（09-30），两格 AC 全 `[ ]` |

⇒ **"哪几枚字段其实已有人排了程"清单**（交回机主）：
`view`→票 145 AC#3／`Q-51`／票 35(ⓑ)；occupancy 三枚→**票 167 AC#2**；
tasks/子代理状态维→**票 188 AC#2**（＋票 197 已落源）；session/history→**票 166**（`panel-sidebar.tsx:71` 注释指名）；
`git` 分支显示→**票 181** 已落源（`composer.git`），无人在排"页面读它"；
`currentModel`/`modelKnown`→**票 145 AC#2b** 已落源，无人排"页面读它"；
`credentialState`→**票 248** 已落源，无人排"页面读它"；
`instructions`→**票 200** 已落源，无人排"页面读它"；
palette 命令表／ball 四字段／`TaskRow.sub`/`elapsed`→**本腿取不到源，也无在册票**（§5）。

---

## §5 未决清单（逐条，不留空）

1. **`harness/` 那条线没量**：`fixtures/harness-app.ts` 的 17 枚 `Harness*` interface（`HarnessMainProps`、`HarnessSidebarProps` 等）＋`harness/main.tsx`、`right-rail.tsx`、`sidebar.tsx` 的字段名册。它们是不是"面板页面"口径的一部分，**机主的框里没说**，本腿按"取不到口径"处理、未展开（§1.3 末注）。
2. **`GitWorktree` 尾部与 `GitView` 的四个 `GitKind*` 常量名没读完**（`git.go:93-116`、`git.go:148-`），`§3.4②` 的 worktree 对拉只到 `path/branch/sha/detached` 四枚。
3. **`snap.Instructions =` / `snap.Tasks =` 的赋值点穷举没做**：本腿只证明"`NewSnapshot` 不填这两枚"（`composer.go:131-136`）与 `instructions_200.go:210`/`subagent_roster_197.go:199` 两处构造函数。⇒ **"泵发出的包到底带不带这两段"我没有读数**，票 145 收表那句"新段出门是 null/全零"本腿**既未复现也未推翻**。
4. **页面的 `snapshot` prop 在真宿主里由谁传**：`App.tsx:91` 只写了 `App({ snapshot = EMPTY })`，本腿在 `cmd/wisp` 与 `internal/panel` 里找不到 Go→页那一跳（§2.4 尺 0 命中）。`?harness=1` 那条是 fixture 路径。**"面板今天到底有没有真数据可渲染"这一问题本腿答不出**，且它比票面的"字段少"更前一格。
5. **`internal/session` / `internal/statemachine` 的导出形状取不到**（§3.2 两条尺 0 命中，正控是 `approval/queue.go:148 Depth()` 那条同族尺出了读数）⇒ 只写"尺取不到具名导出"，**不许被下游读成"无源"**。
6. **`TaskRow.sub` / `elapsed`、`PaletteItem.icon` 的 Go 侧有无源**：尺未跑（`elapsed` 那条我在 §3.2 标了"未跑＝未证"）。
7. **与转述冲突的原文两处**：①派单说票 145 写"恰四字段"——现量**六枚**（§2.1）；②派单说"页面十四态里十一态没有输入可画"——`panel-views.ts` 今天只有**九枚屏**（`:26-35`），票 145 自己的更正①也已把"11"改成 **8**（"要凑到 11 必须把行 14 也算成无输入，而行 14 恰恰是 fed 的"）。⇒ "十四"是 `PLAN.md:3473-3488` 那张状态表的行数，**不是屏数**，本腿没有量那张表（不在射程）。
8. **`App.tsx:293-301` 是否真的把 `workspaceName`/`workspacePath` 传给了 `PanelSidebar`**：本腿读到的那段里只有 8 行（`:293-301`），`workspaceName`/`workspacePath`/`workspaceSessions` 三枚**在截到的实参里没出现**（TS 必填 ⇒ 要么在截外、要么编译不过）。⇒ 未证，需逐行读 `App.tsx:285-312` 全文。
9. **`MonitorPopover` 的 `agents` 实参**（`App.tsx:240`-起那块）没读完 ⇒ §1.3 那一行标了未证。
10. 本腿**没有**跑任何 `go test`/`go build`/`go vet`（同仓 `272-r2` 在取整包读数，包级互斥）⇒ §2.3 全部是 grep 级读数，**没有一枚"颜色"可以引**。
