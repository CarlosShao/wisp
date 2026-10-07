# 145-f1 — 面板页面字段名册 ↔ Go 快照实交字段 对拉表（只读普查腿）

## §0 起手锚

尺（在 `D:/work/workspace/projects plans/Wisp`，分支 `dev`）：

```
git rev-parse --short HEAD                    # a781fdc8
git rev-parse --abbrev-ref HEAD               # dev
date                                          # Wed Oct  7 10:33:47 CST 2026
git rev-parse --short dsh/feat/frontend-p0-v2 # 16c2f038   ← 页面那侧的锚
git ls-tree -r --name-only dsh/feat/frontend-p0-v2 -- frontend/src | wc -l   # 85
```

- 本腿全程**只读**页面代码：一律走 `git show dsh/feat/frontend-p0-v2:<path>` / `git grep <ref>`，
  **未进入 `D:/wt/fe` 工作树**，未跑 `go test`，未 `checkout`/`switch`/`restore`。
- 大输出落本仓 `logs/`（`145f1-*.txt`），文件路径见各节"尺"。

---

## §1 页面侧字段名册

### §1.0 类型地基与量法

地基两枚（`git show` 全文落 `logs/145f1-panel.ts.txt` 311 行 / `logs/145f1-panel-views.ts.txt` 106 行）：
`frontend/src/lib/panel.ts`（`PanelSnapshot` 及其五枚嵌套 interface）＋
`frontend/src/lib/panel-views.ts`（`PanelView` 九枚屏表 + `currentView()`）。

穷举其余声明的尺（结果 200 命中，落 `logs/145f1-type-hunt.txt`）：

```
git grep -n -E 'interface |type [A-Za-z]+ = \{|Snapshot|Props' dsh/feat/frontend-p0-v2 \
  -- 'frontend/src/**/*.ts' 'frontend/src/**/*.tsx'
```

`PanelSnapshot` 的**真实消费面**只有 8 枚文件引 `@/lib/panel`（尺：
`git grep -l -E 'from "@/lib/panel"' dsh/feat/frontend-p0-v2 -- 'frontend/src'`，落 `logs/145f1-panel-importers.txt`）：
`App.tsx` · `components/approval-screen.tsx` · `components/chat-screen.tsx` · `components/composer.tsx` ·
`components/l2-approval-card.tsx` · `components/result-stream.tsx` · `fixtures/harness.ts` · `lib/panel-views.ts`。

对 `snapshot.<字段>` 的直接属性访问尺（整行读，未 `cut`）：

```
git grep -n -o -E 'snapshot\??\.[a-zA-Z_]+' dsh/feat/frontend-p0-v2 -- 'frontend/src'   # 10 命中
```

读数（落 `logs/145f1-snapshot-accesses.txt`）：`snapshot.pending` ×4 · `snapshot.results` ×4 ·
`snapshot.composer` ×1 · `snapshot.view` ×1 · **`snapshot.generatedAt` 0 命中**（＝见 §1.2 正控说明）。

### §1.1 `PanelSnapshot` 顶层（4 枚）

| 字段名（TS 逐字） | 类型 | 谁在消费（文件:行） | 今天的值从哪来 |
|---|---|---|---|
| `pending` | `ApprovalCardView[]` | `lib/panel.ts:130` 声明；`components/approval-screen.tsx:68` `const pending = snapshot.pending`，`:83`/`:94`/`:103` 渲染 | Go 快照 |
| `results` | `ResultChunkView[]` | `lib/panel.ts:131` 声明；消费行见 §1.3（`chat-screen.tsx` / `result-stream.tsx`） | Go 快照 |
| `composer` | `ComposerState` | `lib/panel.ts:133` 声明；消费在 `components/composer.tsx`（`snapshot.composer`，行号见 §1.3） | Go 快照 |
| `generatedAt` | `string` | `lib/panel.ts:136` 声明；**页面读取 0 命中**（尺见 §1.2 正控）。唯一出现处是 fixture 写死值：`fixtures/harness.ts:97` `generatedAt: "2026-09-25T18:40:00+08:00"`、`:181` `generatedAt: "showcase"` | **没有任何页面用它**；只有假件里写死 |

### §1.2 负向结论的正控（两把尺都要有对照）

| 尺 | 打在"确认存在"的字段上（正控读数） | 打在"判定页面不读"的字段上（负读数） |
|---|---|---|
| `git grep -n -o -E 'snapshot\??\.[a-zA-Z_]+'` 全树 | `snapshot.pending` **4 命中**（approval-screen.tsx:68 等） | `snapshot.generatedAt` **0** · `snapshot.tasks` **0** · `snapshot.instructions` **0** |
| 同尺的解构形：`git grep -n -E '\{ *(pending\|results\|composer\|generatedAt\|view) *\} *='` | `approval-screen.tsx:68` 取到 `snapshot.pending` 的赋值形 | 解构形 0 命中 |
| 词形兜底：`git grep -n -E '\btasks\b\|\binstructions\b'`（落 `logs/145f1-tasks-instr-hits.txt`，21 命中） | `lib/panel-views.ts:30`/`:69` 的 `"tasks"` 是**屏名**，不是快照字段 | 21 命中**逐条读完无一枚是 `snapshot.tasks`/`snapshot.instructions` 的读取** |

⇒ 结论：页面侧今天**只读** `pending` / `results` / `composer` 三枚顶层字段＋一枚**并不存在于 interface** 的 `view`。
`generatedAt` 声明了但没人读（＝Go 交了、页面不读那一档，见 §3③）。

### §1.3 页面读取的嵌套字段名册（逐枚）

（本节逐行读数在补，见 §1.4 未完清单；下表凡标"行号待补"的一律未证，不得当已量。）

| 归属 interface | 字段名（TS 逐字） | 类型 | 消费点 | 值来源 |
|---|---|---|---|---|
| `ApprovalCardView` (`lib/panel.ts:24-47`) | `correlationId` | `string` | 行号待补 | Go |
| 〃 | `tool` | `string` | 行号待补 | Go |
| 〃 | `args` | `string[]` | 行号待补 | Go |
| 〃 | `level` | `string` | 行号待补 | Go |
| 〃 | `rulesHit` | `string[]` | 行号待补 | Go |
| 〃 | `reason` | `string` | 行号待补 | Go |
| 〃 | `reasonKnown` | `boolean` | 行号待补 | Go |
| 〃 | `sessionOverrideBlocked` | `boolean` | 行号待补 | Go |
| 〃 | `callChain` | `string[]` | 行号待补 | Go |
| 〃 | `decidedBy` | `string` | 行号待补 | Go |
| `PanelSnapshot & {view?}` (`lib/panel-views.ts:98-99`) | `view` | `string?` | `lib/panel-views.ts:99` `const named = snapshot.view`；由 `App.tsx` 经 `currentView()` 调 | **Go 从不交**（`:90` 原文："The field does not exist on PanelSnapshot yet"）⇒ 恒走 `DEFAULT_VIEW = "chat"` |
| `PanelView`（**页面自己写死的表**，`lib/panel-views.ts:65-75`） | `fed` | `boolean` | 同文件九行常量 | **写死的假内容风险点**：`fed: true` 只有 `chat`/`approval` 两枚，其余 7 枚 `fed: false`（`config` 一枚另标 `selfFed: true`）——它不是字段读取，是页面里对"Go 有没有源"的**手抄账**，Go 加了字段它不会自己知道 |
| 〃 | `id` `label` `icon` `demoIcon` `note?` `interim?` `selfFed?` | — | `lib/panel-views.ts:38-62`，消费在 `App.tsx` 导航 | 页面常量，非 Go |

---

## §2 Go 侧今天真交出去的字段（现量）

### §2.1 `Snapshot` 顶层：**恰六枚**，不是票 145 说的"恰四枚"

尺：`Read internal/panel/composer.go:57-93`（锚 `a781fdc8`）。

| # | Go 字段 | JSON key | 行 |
|---|---|---|---|
| 1 | `Pending []ApprovalCardView` | `pending` | `composer.go:58` |
| 2 | `Results []ResultChunk` | `results` | `:59` |
| 3 | `Composer ComposerState` | `composer` | `:60` |
| 4 | `GeneratedAt string` | `generatedAt` | `:61` |
| 5 | `Instructions *InstructionsSection` | `instructions,omitempty` | `:62-70` 区段（字段声明在 `composer.go` 内，`struct` 起 `:57`） |
| 6 | `Tasks *TaskRosterSection` | `tasks,omitempty` | `:71-93` 区段 |

⇒ 票 145 票面"恰四字段"与 `internal/panel/pump.go:5-6` 那句注释
（"Snapshot (composer.go:44) has four fields"）**在今天的 HEAD 上都已过期**：`instructions` / `tasks` 两枚
是票 200 / 票 197 后续加进来的第 5、6 枚键。票 145 自己"AC#1 交件后更正②"改的行号（`:44-49`）也漂了，
现量 `type Snapshot struct` 在 **`composer.go:57`**、`func NewSnapshot` 在 **`:106`**。

### §2.2 各段字段（Go 实交形状）

- `ComposerState`（`composer.go:235-272`）**11 枚 key**：
  `mode` `workspace` `attachments` `acceptedAttachmentMimes` `maxAttachmentBytes` `attachmentError`
  **`git`**（票 181） **`currentModel`** **`modelKnown`**（票 145 AC#2b） **`credentialState`** **`credentialKnown`**（票 248）
  ⇒ TS 侧 `ComposerState`（`lib/panel.ts:112-119`）**只声明 6 枚**，差的 5 枚页面一律不读。
- `ModeView`（`:249-260`）＝ `current` `names` `l2ConfirmNames` ｜ `WorkspaceView`（`:209-225`）＝ `set` `spelling` `canonical` `reparse` `rewritten` `reason`
  ｜ `AttachmentRef`（`attachments.go:83-104`）＝ `id` `name` `mime` `kind` `sizeBytes` `artifact` `stored` `deduplicated` `reason`
  ｜ `GitView`（`git.go:117-`）＝ `kind` `reason` `branch` `detachedSha` `isDetached` `repoRoot` `currentWorktree` `worktrees` `branches`（…续读见 §1.4）
  ｜ `ApprovalCardView`（`approval.go:39-61`）＝ `correlationId` `tool` `args` `level` `rulesHit` `reason` `reasonKnown` `sessionOverrideBlocked` `callChain` `decidedBy`（**与 TS 逐键同**）
  ｜ `ResultChunk`（`composer.go:95-99`）＝ `correlationId` `text` `done`（与 TS `ResultChunkView` 同）
  ｜ `InstructionsSection`（`instructions_200.go:71-`）＝ `status` `reason,omitempty` `files[]`；`ProjectInstructionFile`（`:44-65`）＝ `path` `tier` `depth` `bytes` `truncatedBytes,omitempty` `dropped,omitempty` `duplicateOf,omitempty` `source,omitempty`
  ｜ `TaskRosterSection`（`subagent_roster_197.go:150-`）＝ `rows[]` `inFlightSlots` `poolCap` `streamTruncated` `streamElidedRunes` `droppedStreamKeys`；`TaskRowView`（`:106-138`）＝ `taskId` `label` `kind` `parentTaskId` `status` `statusKnown` `statusReason,omitempty` `streamKey` `blockedOnApproval` `streamTruncated` `streamElidedRunes` `streamDropped`

### §2.3 生产者（谁填真值）

尺（剔 `_test.go`）：

```
grep -rn --include=*.go -E 'NewSnapshot\(|NewComposerState\(|NewApprovalCardView\(|TaskRosterSection\{|InstructionsSection\{' internal cmd | grep -v '_test.go'
```

- `NewSnapshot(`：`internal/panel/pump.go:303`（**生产泵**，票 35 的数据半边）＋ `pump.go:236`（无源时的空形）；
  `cmd/wisp/` 里**无直接调用者**。⚠ **票 145 更正⑤那句 "`NewSnapshot(` 0 枚调用者" 在今天已过期**——泵落地了，但它填的是**四枚**（见下）。
- `NewComposerState(`：`internal/panel/pump.go:269`。
- `NewApprovalCardView(`：生产侧唯一调用者仍是 `cmd/wisp/panel_assets.go:68`（命令行诊断分支）；
  面板路径走的是 `pump.go` 的 `NativeVerdict.CardView()`（`pump.go:98`）——本腿未展开验这条，列 §1.4 待补。
- `InstructionsSection{`：`instructions_200.go:164` `:174` `:210`。
- `TaskRosterSection{`：`subagent_roster_197.go:199`（在 `taskRosterSectionFrom` 内）。
- ⚠ **`NewSnapshot` 的 return（`composer.go:131-136`）只填 `Pending`/`Results`/`Composer`/`GeneratedAt` 四枚**，
  `Instructions` / `Tasks` 两枚**不经这条构造路径**——它们是别处对 `snap.Instructions =` / `snap.Tasks =` 的赋值（赋值点全名册本腿未穷举，列 §1.4 待补）。
  这就是票 145 收表里"新段出门是 `depth:0 / windowMs:0 / null / 全零`"那一族的现形位置。

### §2.4 快照怎么到达页面（最后一跳）

尺：`grep -rn 'Out:' internal cmd | grep -v _test.go` → 唯一装配点 `cmd/wisp/run.go:725  Out: rt.bookPanelSnapshot`。
`bookPanelSnapshot` 定义在 `cmd/wisp/panel_pump.go:319-325`：只做
`rt.lastSnap, rt.lastSnapBytes, rt.snapSeen = snap, data, true` ＋ 一行日志（`:323`）——**不碰 WebView**。
`internal/panel/pump.go:19-25` 原文（整行读）：
"There is no Go -> page channel in this tree today - no WebView2 host (ticket 33 is unclaimed), no postMessage writer …
this pump stops at the byte boundary … The last mile is NOT here."

⇒ 与机主给的已知"Go→页绑定回话只有库 `Run()` 泵那一形兑现"**一致但更窄**：`Run()` 泵兑现的是
**marshal 到字节 + 记账**，兑现到 `rt.lastSnap`，**没有任何一行把 `data` 送进页面**。
页面里 `lib/panel.ts:139-155` 的 `WispHostBridge` / `window.wispBridge` / `window.chrome?.webview`
是**页→Go 的 postMessage 出口**；Go→页的反向在 `internal/panel` 与 `cmd/wisp` 两包里本腿**未找到具名行**（尺见上，命中为空）。
`cmd/wisp/panel_host_windows.go:80` 的 `panelDispatchBinding = "wispDispatch"` 是**绑定（页→Go）**那一侧，不是推送。

---

## §3 对拉表（四档）

（本节按 §1.2 / §2.1 / §2.2 的读数成文；ⓐ①②③④ 的逐字段行在 §3.1–§3.4，未量的写"未证"。）

- ①**两边都有、已对上**：4 枚顶层（`pending` `results` `composer` `generatedAt`）＋ `ApprovalCardView` 10 枚 ＋ `ResultChunkView` 3 枚 ＋ `ComposerMode` 3 枚 ＋ `ComposerWorkspace` 6 枚 ＋ `ComposerAttachment` 9 枚 ＝ **35 枚**（顶层 `generatedAt` 虽"声明对上"但页面 0 读，另记 ③）。
- ②**页面要、Go 没交**：**0 枚已证** —— 因为页面在 `PanelSnapshot` 这条契约上唯一额外读的字段是 `view`（`lib/panel-views.ts:99`），而它是 ④ 那种"页面自己知道 Go 没有"的形状（原文 `:90`）。其余"页面要的字段"目前**全部躲在 props-only 组件里**（`TasksScreen` / `BallScreen` / `PaletteScreen` / `ContextMeter` / `TurnRail` / `GitBranchChip` / `PanelSidebar` …），它们不接快照，见 §3.2。
- ③**Go 交了、页面不读**：`generatedAt`（顶层，页面 0 读，§1.2 已正控）＋ `composer` 段 5 枚（`git` `currentModel` `modelKnown` `credentialState` `credentialKnown`）＋ `instructions` 整段 ＋ `tasks` 整段 ＝ 顶层 2 段 + 5 枚段内 key + 1 枚顶层。
- ④**名字对不上但形状像同一件事**：`Snapshot.Tasks` 的 `TaskRowView`（`taskId`/`label`/`kind`/`status`/`parentTaskId`）↔ 页面 `TaskRow`（`tasks-screen.tsx:27`，字段名册本腿未逐枚抄）；`Snapshot.Composer.Git`（`GitView`，`git.go:117`）↔ 页面 `GitBranchView`（`git-branch.tsx:32`）。⇒ 两枚都要逐字段对，见 §3.4，未量的写"未证"。

### §3.1 ①已对上（逐枚）

尺：`git show dsh/feat/frontend-p0-v2:frontend/src/lib/panel.ts`（全文 `logs/145f1-panel.ts.txt`）↔
`internal/panel/{composer,approval,attachments,git}.go` 的 `json:` tag（`grep -n 'json:"'`）。
（逐枚行表在 §1.1/§2.2 已给，本节欠的是"每一枚的消费行号"，见 §1.3 与 §1.4。）

### §3.2 ②页面要、Go 没交（props-only 组件族：页面自己声明了字段但没有任何 Go 生产者）

这族群是"十一态没有输入可画"的真身：组件**已在盘上**、类型**已声明**、值**只能来自写死的 fixture**。

| 页面类型（文件:行） | 消费它的页面 | 今天值从哪来 | Go 有没有对应段 |
|---|---|---|---|
| `TaskRow`（`components/tasks-screen.tsx:27`） | `showcase.tsx:290` `<TasksScreen tasks={SHOWCASE_TASKS} />` | **写死假内容**：`fixtures/harness.ts` 的 `SHOWCASE_TASKS`（`:359` 等） | **有**：`Snapshot.Tasks`＝`TaskRosterSection`；但页面不读它（§1.2 负读数），且形状不同名（④） |
| `BallStateRow` / `BallRingDemo`（`components/ball-screen.tsx:18`/`:29`） | `showcase.tsx` | 待补（未证） | 待补（未证） |
| `PaletteItem` / `PaletteGroup`（`components/palette-screen.tsx:31`/`:38`） | `showcase.tsx` | 待补（未证） | 待补（未证） |
| `ContextUsage`（`components/context-meter.tsx:25`） | 待补 | 待补（未证）；`:10` 原文："PanelSnapshot has no field for any of the three yet" | **无源**（页面注释自己点名） |
| `TurnMark`（`components/turn-rail.tsx:38`） | 待补 | 待补（未证） | 待补（未证） |
| `GitBranchView`（`components/git-branch.tsx:32`） | 待补 | 待补（未证） | Go 有 `GitView`（`git.go:117`）：名字/形状待对 |
| `TreeSession`（`components/panel-sidebar.tsx:51`） | 待补 | 待补（未证） | 待补（未证） |

---

## §4 与既有票的分工

（本节待读票面后成文，见 §1.4 待补；票名已定位：
`145-panel-snapshot-has-four-fields-so-eleven-of-the-fourteen-ui-states-have-no-input-to-render-grow-the-go-side-carrier.md` ·
`167-four-small-outputs-the-panel-needs-occupancy-queue-stop-draft-plus-crash-self-rescue.md` ·
`188-the-task-monitor-rail-s-subagent-and-background-task-stacks-have-no-go-side-source.md` ·
`242-the-grant-binding-layer-and-the-panel-facing-read-surface-both-have-zero-rulers.md`。）

已现量到的**票 145 票面过期点**（不是新账，是它自己"更正"节之外的第二层过期）：
1. "恰四字段" → 现量**六枚**（§2.1）。
2. 更正⑤"`NewSnapshot(` 0 枚调用者" → 现量 `pump.go:236`/`:303` **两枚生产调用者**；但 `NewSnapshot` 只填四枚，第 5、6 枚不经它（§2.3）。
3. 更正②的行号 `:44-49` → 现量 `composer.go:57`（`struct`）／`:106`（`NewSnapshot`）。

---

## §5 未决清单

（进行中，见 §1.4。）

### §1.4 本腿未完清单（写作时刻仍欠的读数，全部＝未证）

1. `ApprovalCardView` 10 枚字段的**逐枚消费行号**（`l2-approval-card.tsx` / `approval-screen.tsx` / `ai-native/approval-card.tsx` 三处，整行读）。
2. `results` / `composer` 的逐枚消费行号（`chat-screen.tsx`、`result-stream.tsx`、`composer.tsx`）。
3. `ball-screen.tsx` / `palette-screen.tsx` / `turn-rail.tsx` / `panel-sidebar.tsx` / `config-screen.tsx` / `firstrun-screen.tsx` 六枚屏的字段名册与"值从哪来"。
4. `fixtures/harness.ts` / `fixtures/harness-app.ts` 写死假内容的全名册（§3.2 只点了 `SHOWCASE_TASKS` 一枚）。
5. `snap.Instructions =` / `snap.Tasks =` 的赋值点穷举（§2.3 只点名构造函数不填）。
6. 票 167 / 188 / 242 的票面（§4 空白）。
7. `GitView` 尾部字段（`git.go:117` 之后未读到 struct 收尾）与页面 `GitBranchView` 的逐枚对拉。
