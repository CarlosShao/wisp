# 167-c3 — Q-51（卡面 JSON 键冲突）只读普查

> 本文件由只读普查腿 `167-c3` 产出。**零产品码改动**；唯一可写文件即本文件（`.scratch/wisp/probes/167/c3/census.md`）。
> 票面：`.scratch/wisp/issues/167-four-small-outputs-the-panel-needs-occupancy-queue-stop-draft-plus-crash-self-rescue.md`（AC 勾选框一枚未碰）。
> 前程：`.scratch/wisp/probes/167/c2/census.md`（175 行，commit `9a2a3e00`）——本腿**不重复**它的占用条/序号真源/停止名册普查，只答它 §4① 与 §5① 都推给本腿的那一枚 `Q-51`。
> 搜索根一律显式：`cmd internal tools docs scripts .scratch`。`frontend/**`（含 `frontend/dist`）与 `design/**` **未读、未转述**；两层禁令按票面 §禁区执行。⛔ 本腿未跑任何 Go 命令（`go test`/`build`/`vet`/`run`/`gofumpt`/`d22scan`/`staticcheck` 全部为零）。

## 0. 起手锚（逐字读数，同一发命令取）

`date -Iseconds`：

```
2026-10-03T09:12:18+08:00
```

`git log -1 --format='%h %ad' --date=iso`：

```
3736f0dd 2026-10-03 09:11:55 +0800
```

分支（`git rev-parse --abbrev-ref HEAD`）：

```
dev
```

`git status --porcelain | wc -l`（全仓）：**397 行**。
`git status --porcelain cmd internal tools docs scripts`：**空**——产码五根起手全部干净，本腿读到的 `cmd`/`internal`/`tools` 内容即 HEAD `3736f0dd` 的内容。`.scratch` 下的 397 行属他程遗留与在飞写腿的工作面，与本腿写面（仅本文件）不相交。

## 1. 出向 JSON 键名册（逐枚带 json tag＋file:line）

### 1.0 先划边界：谁真的出向

尺（本腿现跑，只读）：`grep -rn "json:\"" internal/agent/approval/*.go` ⇒ **零命中**（连测试文件一起零命中）。
⇒ **`PanelItem`（`internal/agent/approval/ui.go:50`，字段含 `Depth int`，`:56`）与 `LiveApproval`（`internal/agent/approval/pending_read.go:46-54`，字段含 `Position int`，`:52`）身上没有任何 json tag——它们今天不出向**，一次 JSON 都不出。序号真源那一族（167-c2 §2.1）在包内是纯 Go 结构体，进快照必须换一个**有 tag 的载体**。这是 Q-51 的第一步：撞不撞，要先看落在哪一族的哪一枚结构体上。

出向只有一族：`internal/panel` 的快照包（`Snapshot` 及其嵌套）＋同一批 tag 的 CLI 诊断出口（§1.3）。尺：`grep -rn "json:\"" internal/panel/*.go | grep -v _test.go` ⇒ 逐枚过目，全录如下。

### 1.1 快照顶层：`Snapshot`（`internal/panel/composer.go:57`）——4 恒发＋2 可选

| 键 | 行 | 类型 | 发射条件 |
|---|---|---|---|
| `pending` | `composer.go:58` | `[]ApprovalCardView` | 恒发（`NewSnapshot` `:109-111` 把 nil 归一成 `[]`） |
| `results` | `:59` | `[]ResultChunk` | 恒发（`:112-114`） |
| `composer` | `:60` | `ComposerState` | 恒发 |
| `generatedAt` | `:61` | `string` | 恒发 |
| `instructions` | `:74` | `*InstructionsSection` `omitempty` | **只在装了 reader 时发** |
| `tasks` | `:91` | `*TaskRosterSection` `omitempty` | **只在装了 roster reader 时发** |

`pump.go:25` 逐字立着这寸地：**"NO NEW TOP-LEVEL KEYS. Snapshot keeps its four JSON keys and frontend/src/lib/panel.ts is untouched"**，并点名第五枚（"which screen"）＝`Q-51` 领地、票 145 AC#3 已裁"不在这片切片"（`internal/panel/pump.go:25-30`）。

### 1.2 嵌套层（快照里每一枚出向结构体的键名册）

**`ApprovalCardView`（`internal/panel/approval.go:39`）＝序号最近的一枚候选载体**，10 枚键：
`correlationId`:40 · `tool`:41 · `args`:42 · `level`:44 · `rulesHit`:46 · `reason`:48 · `reasonKnown`:52 · `sessionOverrideBlocked`:54 · `callChain`:55 · `decidedBy`:58。
**无 position/depth/index/order 任何一枚**；也无 `omitempty` 位可当"恒 0 待用"（下节 §2 逐枚查）。

- `ResultChunk`（`composer.go:95`）：`correlationId`:96 · `text`:97 · `done`:98（共 3）。`pump.go:374-377` 注释自述"ResultChunk keeps its three JSON keys … which this leg may not touch"。
- `ComposerState`（`composer.go:235`）：**11 枚**（不是 167-c2 §1.2 写的 14 枚，见本件 §5 枚 3）：
  `mode`:236 · `workspace`:237 · `attachments`:238 · `acceptedAttachmentMimes`:239 · `maxAttachmentBytes`:240 · `attachmentError`:241 · `git`:246 · `currentModel`:255 · `modelKnown`:256 · `credentialState`:268 · `credentialKnown`:269。
  - `ModeView`（`:141`）：`current`:144 · `names`:146 · `l2ConfirmNames`:149。
  - `WorkspaceView`（`:209`）：`set`:212 · `spelling`:214 · `canonical`:216 · `reparse`:218 · `rewritten`:221 · `reason`:224。
  - `GitView`（`git.go:117`）：`kind`:119 · `reason`:122 · `branch`:124 · `detachedSha`:126 · `isDetached`:128 · `repoRoot`:130 · `currentWorktree`:134 · `worktrees`:136 · `branches`:141 · `switchBlocked`:144。
    - `GitWorktree`（`git.go:93`）：`path`:96 · `branch`:98 · `sha`:100 · `detached`:102 · `main`:105。
- `AttachmentRef`（`attachments.go:83`）：`id`:84 · `name`:86 · `mime`:90 · `kind`:92 · `sizeBytes`:94 · `artifact`:96 · `stored`:98 · `deduplicated`:101 · `reason`:103。⚠ **这枚双向在册**（既在快照里、又在入向 registry 里，§3.3）。
- `InstructionsSection`（`instructions_200.go:71`）：`status`:73 · `reason`:77(omitempty) · `files`:79。
  - `ProjectInstructionFile`（`instructions_200.go:44`）：`path`:46 · `tier`:49 · **`depth`:52** · `bytes`:54 · `truncatedBytes`:56 · `dropped`:59 · `duplicateOf`:62 · `source`:64。⚠ **`depth` 这个键名在快照里已经被占用**（语义＝指令文件嵌套层数），§2 的命名冲突面。
- `TaskRosterSection`（`subagent_roster_197.go:145`）：`rows`:148 · `inFlightSlots`:152 · `poolCap`:153 · `streamTruncated`:156 · `streamElidedRunes`:158 · `droppedStreamKeys`:163。
  - `TaskRowView`（`:101`）：`taskId`:102 · `label`:103 · `kind`:104 · `parentTaskId`:105 · `status`:109 · `statusKnown`:110 · `statusReason`:115 · `streamKey`:121 · `blockedOnApproval`:127 · `streamTruncated`:130 · `streamElidedRunes`:131 · `streamDropped`:132。

### 1.3 同一批 tag 的**第二枚出向消费者**（不是快照）

`cmd/wisp/panel_assets.go:68` 造 `panel.NewApprovalCardView(...)`，`:79` `json.NewEncoder(os.Stdout).Encode(...)` ⇒ **`ApprovalCardView` 的键集同时是 `wisp panel-assets` 的 stdout 契约**。它的测试没有豁免：`cmd/wisp/panel_assets_143_test.go:52-61` 的 `cardView143` **逐字重declare 了这 8 枚键**（`correlationId`/`tool`/`level`/`rulesHit`/`reason`/`reasonKnown`/`sessionOverrideBlocked`/`decidedBy`；比产码少 `args`、`callChain` 两枚）。⚠ 但它用裸 `json.Unmarshal`（`:112` 一带，**无 `DisallowUnknownFields`**）⇒ **加一枚键不会让它红**，它只是"名册的第二份手抄"，会在加键后**静默地少抄一枚**。这条差异本件记进 §3.2。

### 1.4 入向键名册（Q-51 只问出向，但同一批尺子两边都读，故列全）

- `ComposerRequest`（`bridge.go:91`）：`method`:92 · `requestId`:93 · `source`:94 · `to`:96 · `path`:97 · `text`:98 · `attachments`:100 · `configField`:117 · `configProvider`:118 · `configModel`:119 · `configValue`:120。
- `ModeRequest`（`composer.go:176`）：`to`:177 · `correlationId`:179。
- `WorkspaceRequest`（`composer.go:228`）：`path`:229。
- `AttachmentPayload`（`attachments.go:367`）：`name`:368 · `declaredMime`:369 · `sizeBytes`:370 · `dataBase64`:371。
- `OutgoingMessage`（`composer.go:295`）：`text`:296 · `attachments`:297（Go→agent 侧复用，不入快照）。

## 2. 序号落点候选：现读有没有"已存在但恒 0/恒空"的字段可承载

（取数中）

## 3. 冲突面：哪几把尺会响

（取数中）

## 4. 停止那一半：4 枚入向名册＋最少一行

（取数中）

## 5. 我可能写错的条目（自我对抗）

（取数中）

## 6. 量不到的地方

（取数中）

## 7. Q-51 的答案一句话

（取数中）

## 8. 交件判语

（取数中）
