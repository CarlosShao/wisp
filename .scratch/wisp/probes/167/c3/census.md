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

> 问法拆两问才有效：**（a）快照的出向键里有没有一枚今天空着、名字与语义都还给得起序号的？**（b）如果没有，序号要进来最少要动几跳？

### 2.1 (a) 逐枚过目的结果：**没有可复用的槽**

尺：`grep -rn "json:\"" internal/panel/*.go | grep -v _test.go` 全 **112** 命中逐枚过目并按结构体归并（名册见 §1）。判据＝"这枚键今天恒 0/恒空，且它的名字与语义承载'第几条'不算说谎"。

- **`ApprovalCardView` 里没有任何整数键**（§1.2 那 10 枚＝8 枚字符串/布尔＋2 枚字符串切片）。序号是个整数，**同型可复用的槽＝零枚**。
- 唯一一枚在快照路线上**恒空**的卡面键＝`callChain`（`approval.go:55`，`[]string`）。它恒空是**产码事实**：`liveVerdicts()`（`cmd/wisp/panel_pump.go:65-73`）逐字段抄写里**没有 CallChain**，`NativeVerdict.CallChain`（`pump.go:86`）因此是 nil，`approval.go:95-97` 把 nil 归一成 `[]`。
  ⛔ **但它不能承载序号**，两条独立理由：①型不对（`[]string` vs 序数）；②语义已被占且被**明文保护**——`pump.go:84-86` 逐字："inventing `["session","agent-loop",tool]` here would put a string no code produced onto a security card"；`cmd/wisp/panel_pump.go:50` 逐字 "CallChain stays empty"，`:55-56` 同一意思的第二处措辞（"a sentence no producer wrote onto a security card"），`:53-54` 并把原因钉死："the only caller that ever filled CallChain is the CLI probe"。
  ⇒ **它不是全仓恒空、只是快照路线恒空**：`cmd/wisp/panel_assets.go:72` 真填 `[]string{"cli","panel-assets",*l2}`。动它＝同时动两枚出向消费者。
- 整数键候选逐枚排除（都在别的段，语义已满）：`maxAttachmentBytes`（`composer.go:240`，`pump.go:238-240` 有兜底真值）· `sizeBytes`（`attachments.go:94`/`:370`）· `depth`（`instructions_200.go:52`，**产码有真生产者**：`instructions_200.go:145` `Depth: f.Depth`）· `bytes`/`truncatedBytes`（`:54`/`:56`）· `inFlightSlots`/`poolCap`/`streamElidedRunes`（`subagent_roster_197.go:152`/`:153`/`:131`）。
- ⚠ **`depth` 这一枚键名在快照里已经被占用**（＝指令文件嵌套层数，`instructions_200.go:52`），而且 `instructions_200_test.go:300` 把词面 `"depth":0` 钉进了那一段的字节里。若把"共几张卡"命名成 `depth`：不同父对象（`composer` 段 vs `instructions.files[]`）⇒ **上线上不撞**，但**页面侧读者会撞**——同一份 `panel.ts` 里两枚 `depth` 两种意思。
- 另有一枚**已经存在的"队列总数"读数，但它不在快照 JSON 里**：账本摘要行的 `depth=<len(snap.Pending)>`（`cmd/wisp/panel_pump.go:350` 的 `fmt.Sprintf`，超限时 `:357` 换一条同样带 `depth=` 的窄形），被 `cmd/wisp/panel_pump_test.go:97` 钉住 `rec["depth"] == fmt.Sprint(len(snap.Pending))`。⚠ 它是 **log record 的 `msg` 串里的一个 `k=v` token，不是出向 JSON 的键**——测试先在 `:63` 用 `Msg string \`json:"msg"\`` 解出整条消息，再在 `:73` 用 `strings.Cut(f, "=")` 现拆（tag 常量在 `:59`）。⇒ "共几张"这个数今天**已经外递了，只是递在账本里**，票面 AC#3 要的"可读出口"不认这条。

**(a) 结论：必须新增键**（要 Go 侧出口这条前提立着的话）。唯一不动契约的替代＝页面自己数 `pending[]` 下标（167-a1 §3.1 已具名，本腿不再展开），那是"页面推出来的数"。

### 2.2 (b) 序号要走到卡面，缺的是**三跳**不是一跳

现量链路（每一跳都点名，⛔ 不选形）：

| 跳 | 现状 | file:line |
|---|---|---|
| 真源 | `LiveApproval.Position` 有、产码已填 | `internal/agent/approval/pending_read.go:49`（声明）· `:118`（填，`q.position(it)`） |
| 跳 1 | `liveVerdicts()` 抄了 7 枚字段，`it.Position` 一个字节没抄 | `cmd/wisp/panel_pump.go:65-73` |
| 跳 2 | 载体 `NativeVerdict` 只有 8 枚字段、**无 Position**（它也无 json tag ⇒ 本跳不撞任何契约） | `internal/panel/pump.go:77-91` |
| 跳 3 | 卡面由 `CardView()` 现造：`CardViewFromDecision(ApprovalSubject{…}, risk.Decision{…})`，**两枚入参结构体都没有序号位**，返回值直接 append 进 `cards` | `internal/panel/pump.go:98-112`（`CardView`）· `:215-217`（`Snapshot()` 里 `cards = append(cards, v.CardView())`）· `internal/panel/approval.go:27-36`（`ApprovalSubject`）· `:72`（`CardViewFromDecision`） |
| 落点 | `ApprovalCardView` 加一枚带 tag 的键 ⇒ **撞 Q-51**（射程见 §3） | `internal/panel/approval.go:39-59` |

⚠ 167-c2 §2.1 把这判成"断点两处"（载体无字段＋抄写未抄）——**在"到 NativeVerdict 为止"这个口径下属实**；但 c2 §4① 自己列的改动面是三份文件，本腿把第三跳的**形状**量清：`CardView()` 是纯函数、只吃 `(ApprovalSubject, risk.Decision)`，序号进不了它的任何一枚入参，除非 ①`ApprovalSubject` 也加一枚字段，或 ②在 `pump.go:217` 那行之后对返回值做一次赋值（`cards[i].Position = v.Position`）。
**①有一个 AC#2 型的副作用必须点名**：`NewApprovalCardView`（`approval.go:64-67`）与 `CardViewFromDecision` 是**同一条构造路**，CLI 诊断腿 `cmd/wisp/panel_assets.go:68` 也走它——那条路上**根本没有队列**，序号只能填零 ⇒ `wisp panel-assets` 的 stdout 会打印一枚 `"position":0`，这正是票面 AC#2"不许退化成 0"点名的形状（且 `cardView143` 因 §1.3 的宽松解码**不会红**）。走②则 CLI 路线那枚键压根不发。⇒ **两形的差别在契约上看得见，本腿只登记，不选形。**

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
