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
⚠ **副作用与选形无关，只与 tag 写法有关（这是本腿写完 §2 后回头推翻自己的一句，详见 §5 枚 2）**：`NewApprovalCardView`（`approval.go:64-67`）与 `CardViewFromDecision`（`:72`）是**同一条构造路**，CLI 诊断腿 `cmd/wisp/panel_assets.go:68` 也走它，而那条路上**根本没有队列**。本腿先前写"走②则 CLI 路线压根不发这枚键"——**不成立**：`encoding/json` 对没有 `omitempty` 的 int 字段**恒发 0** ⇒ 只要 `ApprovalCardView` 上多出 `position`，`wisp panel-assets` 的 stdout 就**必然**带 `"position":0`，而这正是票面 AC#2"不许退化成 0"点名的形状（且 `cardView143` 因 §1.3 的宽松解码**不会红**）。
⇒ 真正决定发不发的不是构造路径而是 **tag 写法**：带 `omitempty` 或指针形才谈得上"没有队列就不发这枚键"，而这一枚恰好适用——`Position` 是 **1-based**（`pending_read.go:44-45` 注释逐字 "the 1-based place in the FIFO"）⇒ 0 从来不是合法序号，被 `omitempty` 省掉的只会是"根本没有队列"那条路。两形差别本腿只登记，不选形。

## 3. 冲突面：哪几把尺会响

> 本节所有"响/不响"都是**射程判断，非内容引用**——判据＝现读 Go 侧尺的代码本体（`internal/panel/*_test.go`、`cmd/wisp/*_test.go`）与它读入的**路径**，⛔ 未读页面侧任何文件内容。页面侧今天到底声明了哪几枚键＝本腿量不到（§6 枚 1）。

### 3.1 全仓只有**两把**双向对账尺（尺：`grep -rn "jsonKeysOf\|tsInterfaceKeys" internal/panel/*_test.go` ⇒ 调用点恰 4 行、分属两枚测试）

两把共用同一组 helper（`approval_test.go:174` `jsonKeysOf`、`:193` `tsBlockRe`、`:208 `tsInterfaceKeys`、`:213` `subtract`），输入都是 `os.ReadFile(root/frontend/src/lib/panel.ts)`（`approval_test.go:107`、`composer_test.go:50`）。**两方向各一次 `t.Errorf`**：Go 有 TS 没有 ⇒ "emits [...] that interface ... does not declare"；TS 有 Go 没有 ⇒ "reads [...] that Go ... never sends - those fields render as undefined"。

⚠ **射程的决定性读数是两把尺各自覆盖的结构体名册并不相同**（本腿逐枚抄自 `pairs` 字面量）：

| 尺 | 位置 | 覆盖的 Go 结构体 ↔ TS interface |
|---|---|---|
| A | `composer_test.go:48` `TestComposerContractTypesMatchFrontend` | `Snapshot`↔`PanelSnapshot`(:60) · `ComposerState`↔`ComposerState`(:61) · `ModeView`↔`ComposerMode`(:62) · `WorkspaceView`↔`ComposerWorkspace`(:63) · `AttachmentRef`↔`ComposerAttachment`(:64) · `ResultChunk`↔`ResultChunkView`(:65) |
| B | `approval_test.go:105` `TestApprovalCardViewJSONKeysMatchFrontendTypes` | `ApprovalCardView`↔`ApprovalCardView`(:118) · `ResultChunk`↔`ResultChunkView`(:119) · `Snapshot`↔`PanelSnapshot`(:120) |

⇒ **`GitView`／`GitWorktree`／`InstructionsSection`／`ProjectInstructionFile`／`TaskRosterSection`／`TaskRowView` 不在任何一把尺的名册里**（旁证：`grep -rln 'panel\.ts' internal/panel` 的 9 枚命中里**没有** `git_test.go`、`instructions_200_test.go`、`subagent_roster_197.go`）。这解释了 A395 那条裁度的机制——"嵌在既有段里开子格、不加顶层 key"不是措辞，是**真的有一整族嵌套结构体落在两把尺的射程之外**。

### 3.2 按"序号键放哪儿"逐格给射程（本腿只判响不响，不判红不红）

| 放法 | 响的尺 | 枚数 |
|---|---|---|
| `ApprovalCardView` 加 `position`（§2 的自然落点） | 尺 B 一枚（`:118`） | **1** |
| `ComposerState` 加"共几张"标量 | 尺 A 一枚（`:61`）——先例就是它：`pump.go:31-40` 逐字记录 248 加了 `credentialState`/`credentialKnown` 两枚键、"the page-side interface does not declare them, and the ruler says so out loud" | **1** |
| `Snapshot` 顶层加键 | 尺 A(`:60`)＋尺 B(`:120`)＋**三枚字节级键集钉**（§3.3） | **5** |
| `GitView`／`TaskRosterSection`／`InstructionsSection` 内加子格 | **零枚双向尺**（不在名册）；只可能撞各自包的词面钉（§3.4） | **0** |
| `AttachmentRef` 加键 | 尺 A(`:64`)＋**入向三向对账**（§3.3 末段）——它是唯一一枚两边都在册的结构体 | **2 族** |

⚠ 一处**必须与 167-c2 不同的判读**：c2 §4① 写"撞 Q-51……加键需要页面声明同批"并并列点名 `approval_test.go:105`＋`composer_test.go:49` 两把。现量：`ApprovalCardView.position` 这一枚**只落在尺 B 的射程里**（尺 A 的 `pairs` 压根没有 `ApprovalCardView`）⇒ 需要同批跟上的 TS 声明是**一枚 interface（`ApprovalCardView`）**，不是两枚。尺 A 只有在动 `Snapshot`/`ComposerState` 时才响。

### 3.3 三枚字节级键集钉＋一枚入向三向尺（都是**逐字读尺体**得出的射程）

- `pump_test.go:123`（`TestThePumpBuildsThePacketFromWhatTheHostHolds`，函数起 `:35`）：`strings.Join(sortedCopy(got),",") != "composer,generatedAt,pending,results"`。
- `pump_test.go:291`（`TestPublishHandsTheBytesToTheAttachedExit`，起 `:263`）：同一串字面量，钉的是**出口收到的字节**。
- `subagent_roster_197_test.go:214`：同一串字面量；且 `:207-208` 反向钉 `if _, ok := wire["tasks"]; ok { t.Fatalf(...) }`（无 roster reader 不许发 `tasks`）。
  ⇒ 三枚都只解**顶层** `map[string]json.RawMessage`（`pump_test.go:115`、`:283`、`roster:203`）⇒ **卡面嵌套键不在它们的射程里**；`omitempty` 指针键能绕过它们落地（`tasks`/`instructions` 正是这么进来的，`composer.go:68-73`、`:80-85` 两处注释逐字自述了这个手法）。
- 入向侧另有一把**三向对账**：`l2_grant_boundary_test.go:1854`（registry 读入）／`:1807` `TestJSONKeyDerivationAgreesWithEncodingJSON`（⚠ 本件是冻结件，本腿只读）——它把 AST 读的键册、反射读的键册、`encoding/json` 的 `DisallowUnknownFields` 裁决（`:1779` `decoderKnowsKey`）三方对撞。**其 registry 只有 4 枚入向类型**：`ComposerRequest`/`ModeRequest`/`AttachmentPayload`/`AttachmentRef`（`:1720-1726`）。⇒ 出向加键不触它；但 **`AttachmentRef` 双役**，动它的键会同时被两把尺＋这把三向尺读。
  ⚠ 同处 `:1878` 那枚子测试对 registry 里每一型跑一份 24 枚"判决词"拼写册（词表 `:1882-1887`：`outcome`/`allow`/`grant`/`decision`/`verdict`/`override`/… ）必须**不可绑定**，断言在 `:1901-1906`。⇒ **只有入向**才受这份词表约束；给序号起名 `position`/`queuePosition` 与之无交集（射程判断）。

### 3.4 词面钉与"手抄名册"（会静默过期、不会响的那一类）

- `cmd/wisp/panel_assets_143_test.go:52-61` `cardView143` 手抄 8 枚卡面键、`:112` 裸 `json.Unmarshal` ⇒ 加键**不红**，只是名册少抄一枚（§1.3）。同文件 `:265` 那份 banned 词表含 `rulesHit`/`sessionOverrideBlocked`，但 `:267` 只 `ParseFile("panel_assets.go")` **一枚文件** ⇒ 除非把判决词写进那枚产码文件，否则不响。
- `instructions_200_test.go:300` 钉 `"depth":0`＋`"tier":"project"`＋`"instructions"` 等词面 ⇒ **这一段字节**里 `depth` 已被占；`cmd/wisp/instructions_200r2_test.go:43` 另有一枚本地重declare 的 `Depth int \`json:"depth"\``。
- `cmd/wisp/panel_pump_test.go:97` 钉账本 `depth=` token ＝ `len(snap.Pending)`；`:203` 另钉 `rec["depth"]=="1"`。**它们读的是日志 msg 里的 k=v，不是快照 JSON**（§2.1）⇒ 卡面加 `position` 不触；把"共几张"做成快照键也不触，但会在同一份 `panel.ts` 里造出第二枚 `depth`（§2.1 命名冲突面）。
- `bridge_test.go:131` `TestFrontendComposerRequestsMatchTheEnvelope` 读同一份 `panel.ts`，但**只单向**：要求页面**含**那 4 枚方法名与信封键（方法名循环 `:138-145`、信封键循环 `:146-156`），另钉 `strings.Count(text,"bridge.postMessage") != 2`（`:161`）。⇒ Go 出向加一枚键它**看不见**；它管的是入向方法与"不许开第二条 IPC"。

### 3.5 一枚本腿新发现的**盲区**（在册，非本腿发明）

`pump.go:47-49` 逐字："**the reconciliations read json tags, so an untagged exported field passes them while encoding/json still emits it** (registered, not used)"。机制复核属实：`jsonKeysOf`（`approval_test.go:182`）只在 `f.Tag.Get("json") != ""` 时才收键名，而 `encoding/json` 会把无 tag 的导出字段按 **Go 字段名**发出去。
⇒ **给 `ApprovalCardView` 加一枚不带 tag 的 `Position int`，两把尺都不响，而线上会多一枚 `"Position"` 键**（PascalCase，页面按 camelCase 读必然 `undefined`）。这不是建议这么做——恰恰相反，它说明**"尺没响"在本寸地上不等于"没撞契约"**，Q-51 的成本不能只按尺的读数结。

### 3.6 补一枚 §3.1 名册之外的射程：页面侧声明**自己也在 Go 仪器的扫描面里**

落地若要"同批带上页面声明"，那批改动除了 §3.1 两把尺，还落在三族**只读页面树、不读键名**的尺射程里（⚠ 射程判断，非内容引用；本腿未跑它们，见 §6 枚 4）：

- `internal/panel/frontend_hygiene_test.go` 六枚测试：`TestPanelFrontendIsStateless`（`:169`，走 `frontendSrcFiles`→`frontend/src` 的 `.ts/.tsx/.css`，`:79`/`:86`）· `TestFrontendHasNoEmoji`（`:246`，走 `frontend/src`＋`frontend/scripts`）· `TestFrontendNeverNamesAnApprovalDecision`（`:281`，走**整棵 `frontend/`**，但 `:290-292` 明确跳过 `/dist/` 与 `node_modules/`；其尺体 `panelDecisionIdentifierRe` 只有一条正则，`tools/d22scan/../frontend_hygiene_test.go:73`＝`approval\.decide`）。
- CI 仪器 `tools/d22scan/main.go:249` 把 `frontend/` 列为 ban #6 的声明扫描面，ban #8 的 emoji 面覆盖 `design/`＋`frontend/`（同文件 `:23-24`、`:40` 的清单）。
⇒ 三条对"position"这种拼写都**不设障**（判读依据＝上面三把尺的字面量本体），但它们是**真实存在的相邻火源**：页面那份声明里若出现被禁的存储 API／emoji／`approval.decide` 字样，红的是 `internal/panel` 与 CI，不是页面自己的构建。
⚠ 另一枚**布局事实**（据 `.gitignore:22-23` 与 `git ls-files`，非页面内容）：`frontend/` 下受版本控制的条目＝**85 枚**，`frontend/src/lib/panel.ts` **在版本控制里**；`frontend/dist/*` 被 ignore、只留 `!frontend/dist/.gitkeep` 锚（`internal/panel/assets.go:63` 注释自述这枚锚是给 `go:embed` 保命用的）。⇒ **`dist` 是构建产物、不入库**，`cmd/wisp/panel_host_gate_test.go:144-155` 那枚 AC#12 对账钉数的正是"`frontend/dist` 下被 git 跟踪的文件数 vs embed 报的 built 位"。

## 4. 停止那一半：4 枚入向名册＋最少一行

> 167-c2 §3.1 已枚过这四枚；本腿**独立复测**（不抄读数），并把"看得见／按得下"两半分开答。`requestTaskStop` 本仓 0 命中这一条 c2 已定案，本腿不复量、不重开。

### 4.1 名册：4 枚"处理体写好、生产调用者零枚"（逐枚 file:line＋方法名＋尺）

| # | 方法名（入向/口） | 声明处 file:line | 本腿复测尺与读数 |
|---|---|---|---|
| 1 | `(*Loop).Steer` | `internal/agent/loop.go:275` `func (l *Loop) Steer(text string) bool {` | 尺 `grep -rn "\.Steer(" cmd internal tools --include=*.go \| grep -v _test.go` ⇒ **零行** |
| 2 | `(*RunningTask).Cancel` | `internal/agent/loop.go:300` `func (t *RunningTask) Cancel() { t.root.Cancel() }` | 尺 `grep -rn "\.Cancel()" cmd internal tools --include=*.go \| grep -v _test.go` ⇒ 12 枚命中**全不是这一枚**：`feedRoot/root/replyRoot/reloadRoot/observe.Root/p` 的 context 取消（`cmd/balldebug/main.go:339,475`、`cmd/wisp/approval_always.go:100`、`cmd/wisp/resident_task_source_windows.go:289,432`、`cmd/wisp/run.go:904,910`、`internal/agent/loop.go:346,531`、`internal/observe/goroutine.go:324`、`internal/observe/logging.go:116`）＋`:300` 自身。⇒ 握着 `*RunningTask` 的人（`cmd/wisp/run.go:1099` 的 `bg`）**从不调它** |
| 3 | `RequestWorkspaceSwitch` | `internal/panel/workspace.go:76` `func RequestWorkspaceSwitch(scope PathScope, input string, audit AuditFunc) (WorkspaceView, error)` | 尺 `grep -rn "RequestWorkspaceSwitch" cmd internal tools --include=*.go \| grep -v _test.go` ⇒ 除声明外**全是注释**（`cmd/wisp/panel_pump.go:80`、`internal/panel/bridge.go:17`、`composer_dispatch.go:74`、`git.go:261`、`instructions_200.go:109`）⇒ 产码调用者 0 |
| 4 | `Options.Control`（`ControlHandler` 字段） | `internal/agent/loop.go:158` `Control ControlHandler`（注释 `:156` "nil = the loop's own default"） | 赋值处 0 枚（c2 尺：`grep -rn "Control:\s" cmd internal tools` ⇒ 只命中 `internal/llm/anthropic/request.go` 的 CacheControl）；loop 自带默认体 `defaultControl`（`loop.go:521` 一带，注释逐字 "covers what the loop itself can do: cancel/stop its running"）⇒ **执行体在、装配处空** |

⚠ 分类要写清：**这四枚里只有枚 3 是面板入向**（它前面是 `ComposerDispatch.Workspace` socket，产码写死 nil，`cmd/wisp/panel_inbound.go:275`）。枚 1/2/4 是 agent 侧的口，票面 AC#4 要的"停止"真正的执行体在枚 2（`RunningTask.Cancel`）与名册句柄（`TaskRoster.AttachCancel`/`Cancel`）上。

### 4.2 "要让停止按钮**看得见**，最少动哪一行"

现量的名册侧事实：`TaskRoster.AttachCancel`（`internal/tools/task.go:420`）存在，全仓**唯一产码调用者是孩子侧** `internal/tools/subagent_197.go:346` `t.d.Roster.AttachCancel(bg.ID, cancelChild)`；根任务从不挂（`cmd/wisp/run.go:1099` `bg := loop.RunAsync(ctx, task)` → `:1105` `rt.tasks.MarkRoot(bg.ID, task)` → `:1106` `res := bg.Wait()`，三行之间没有 `AttachCancel`）。
不挂的后果是**具名的假话**而不是空：对根 id 调 `TaskRoster.Cancel` 走 `internal/tools/task.go:459-460` 答"任务 X 没有在跑（kind=…，state=…），没有可停的句柄"。

⇒ **最少的一行＝`cmd/wisp/run.go:1105` 之后补一行 `rt.tasks.AttachCancel(bg.ID, bg.Cancel)`**（`bg.Cancel` 即 `loop.go:300` 那枚现成 `func()`，`AttachCancel` 的形参正是 `cancel func()`，`task.go:420`；`task.go:421-423` 自带 nil 防护）。它在 `cmd/wisp` 装配根内、不动 `internal/tools` 一字、不动 `internal/agent` 一字、**零枚新方法名**。

⚠ 但本腿必须把话说全：**这一行买到的是"按得下去的 Go 侧前提"，买不到"看得见"**。快照今天没有任何一格说"这发行可停"：

- `TaskRowView` 的键册（§1.2，12 枚）里**没有 cancellable 位**；`TaskRoster` 的导出面（`task.go:215`–`:447` 共 13 枚方法）里除了**会动手的** `Cancel`（`:447`）之外**没有"能不能停"的谓词**——`RunningSubagentIDs()`（`:402`）/`InFlightSubagents()`（`:390`）都只覆盖子代理，不含根行。
- 根行今天唯一的诚实读数已经是 `statusKnown:false`＋`statusReason:"宿主没有登记这一维"`（reader＝`cmd/wisp/panel_pump.go:171` `rec.Out.StateAnswer()` ⇒ 填 `StatusKnown: status != ""`，`:179`；机制＝`internal/tools/task.go:355-357` `MarkRoot` 逐字 "State is deliberately NOT written here"，归属票 196／裁度 `A394`）。⇒ **这一维被票 196 占着，不许拿来当"可停"**。
- ⇒ 若要把"可停"画出来，落点只剩两枚：①`TaskRowView`/`TaskRosterSection` 加一枚键——按 §3.2 的格表，这一族**在两把双向尺射程之外（0 枚尺响）**，是 Q-51 成本最低的一格；②等票 196 把根行 state 填上后由状态推断。**本腿只登记两枚的存在与射程，不选形**（选形归编排者）。

### 4.3 需不需要新增 C17 方法名：登记为"按得下那一半需要"，不动手

现量入向闭合集：`knownComposerMethod`（`internal/panel/bridge.go:146-153`）的 case 恰 6 枚（4 枚 `panel.*`＋`config.get`/`config.set`）；派发表在 `internal/panel/composer_dispatch.go:175-207`（6 个 case＋`default: d.rosterMismatch(req)`）。今天**无一枚与"停"有关**。
⇒ **"按得下"确实需要一枚新的入向方法名 ⇒ 属契约面、须人工批准**，本腿不自决（票面纪律＋c2 §3.3 引 HANDOVER 的死令）。

两枚锚的**射程**（本腿只读尺体，⚠ 射程判断非内容引用）：
- `internal/panel/git_test.go:381-387`：`wantMethods` 恰 4 枚，与 `whitelistMethodsFromSource(bridge.go)`（`:407` 那把 extractor）做 `equalStrings` 全等对账。
- `internal/panel/git_test.go:517`：对真实 `bridge.go` 再钉一次 `len(...) != 4`（`:505-515` 是同测试的正向对照：往桩文件里种第 5 枚 `panel.*` ⇒ 计数要到 5）。
- ⚠ 两枚都**只数 `panel.` 前缀的字面量**（extractor 走 `panelMethodRe`，`git_test.go:394` = `"(panel\.[a-z0-9_.-]+)"`）。这枚事实不是后门，是**在册的既有手法**：`bridge.go:46-64` 逐字记录了 248 为什么把设置路线拼成不带前缀的 `config.get`/`config.set`、并写明这两枚"carries no `panel.` prefix on purpose … ticket 181's anti-drift nail … counts the `panel.*` literals"。⇒ **仪器射程窄于契约面**：换一种拼法可以不动这两枚计数，但 `knownComposerMethod` 的 case 集（`bridge.go:148`）与冻结件 `l2_grant_boundary_test.go` 的枚举/推导（`:777`、`:1240-1241`、`:1282`、`:2066`）会读它——**这一寸整枚归人工批准，本腿不提议走哪条形**。

## 5. 我可能写错的条目（自我对抗）

1. **"全仓只有两把双向尺"的尺有口径边界**。我的尺是 `grep -rn "jsonKeysOf\|tsInterfaceKeys" internal/panel/*_test.go`（调用点 4 行）。⚠ 它只能证明"`internal/panel` 包里没有第三处"；`cmd/wisp` 的测试**无法调用这两枚未导出 helper**，理论上可以自带一把。补尺：`grep -rln "panel\.ts" cmd internal tools` ⇒ `cmd/wisp` **零命中**；`tools/d22scan/scan_test.go` 的命中（`:616`、`:739`、`:1813`、`:1979`）全是 `t.TempDir()` 里 `seedFile` 种下去的**合成 fixture**（同文件 `:606` `root := t.TempDir()`），不是真树。⇒ 结论站得住，但站的是"两把尺＋一次补尺"，不是"证明全仓唯一"。
2. **§2.2 我原话错了一句，已就地更正**："走②则 CLI 路线那枚键压根不发"——**不成立**。`encoding/json` 对**没有 `omitempty` 的 int 字段恒发 0**，所以只要 `ApprovalCardView` 上多出 `position`，`wisp panel-assets` 的 stdout 无论走①还是②都会带 `"position":0`。正确的说法是：**要么带 `omitempty`／指针形，要么接受诊断输出里多一枚 0**。⚠ 而 `omitempty` 在这一枚上恰好是安全的：`Position` 是 **1-based**（`internal/agent/approval/pending_read.go:44-45` 注释逐字 "the 1-based place in the FIFO"）⇒ 0 从来不是合法序号，被省略的只会是"根本没有队列"的那条路。这句是本腿写完 §2 之后回头推翻自己的，**如果编排者只信 §2 的第一版，就会把一个"选形决定契约"的差别读丢**。
3. **§1.0 那句"approval 包零 json tag"只覆盖了包自身**，而序号真源是经 `LiveApproval.Decision`（`pending_read.go:48`）这枚**跨包类型**（`tools.Decision`）把卡面字段带出来的。补尺：`awk '/^type Decision struct/,/^}/' internal/tools/gate.go`（声明在 `internal/tools/gate.go:15`）里 `json:"` 计数＝**0** ⇒ `tools.Decision` 同样无 tag，§1.0 的结论（这一族今天不出向）不因跨包而松动。`internal/tools` 确有 json tag（`fs.go:114-115`、`fs_edit.go:52-58`、`fs_write.go:225-226` 等），但那些是**工具入参 schema**（模型侧），与面板出向不同一路。
4. **§3.1 的两把尺读的是一张写在测试里的字面量表**（`composer_test.go:55-66`、`approval_test.go:113-121`）。⇒ 我 §3.2 表里"嵌套段 0 枚尺响"这一格**只对当前 HEAD 成立**：谁往 `pairs` 里补一行 `GitView`／`TaskRosterSection`，那一格立刻从 0 变 1，而这一改动**不需要动任何产码、不需要解冻任何冻结件**，所以它是本件里最容易过期的一格。
5. **§3.3 那句"只有入向才受这份词表约束"可能被读成"出向永远不受"**。前一句对今天成立（registry 只有 4 型，`:1720-1726`）；但那把尺的注释 `:142` 明写"一个第二解码目标若不在 registry 里，AST 半就没被对账"⇒ 它是**会被补的**，不是设计成永不出向。我给 `position` 的"无交集"判读只覆盖那 24 枚判决词，不覆盖未来词表。
6. **§4.2 的 `bg.Cancel` 是读码判读，没编译过**。我核了三个形参/返回类型：`RunAsync` 返回 `*RunningTask`（`internal/agent/loop.go:321`）、`Cancel()` 是 `*RunningTask` 的指针接收者方法（`:300`）⇒ `bg.Cancel` 作方法值合法、类型 `func()`；`AttachCancel` 第二形参正是 `cancel func()`（`internal/tools/task.go:420`），且 `:421-423` 自带 `cancel == nil` 防护。⛔ 但"能编过"是 Go 命令才给的读数，本腿禁跑 ⇒ 这一行我只签"形状对"，不签"编译过"。
7. **§1.2 那枚"11 枚"是对 c2"14 枚"的更正，可能是口径差不是事实差**。我的口径＝"带 json tag 且真的发射出去的键"（尺：`grep -c "json:\"" internal/panel/*.go` ⇒ 七个产码文件分别 10/13/11/34/15/11/18，总 **112**，与本件 §1 全册逐枚相加一致）。c2 若把 `PumpSources` 的读口槽位或嵌套段的子键算进"字段"，14 也能自洽。**我只签我的口径下的 11**。
8. **§2 的"必须新增"依赖票面 AC#3 的口径**（"Go 侧可读出口"）。若编排者改判"页面自数 `pending[]` 下标即算 AC#3"（167-a1 §3.1 那条替代），则**一枚键都不用加、Q-51 整个不必动**——那是裁定不是我的读数，本腿不替他选。
9. **我把"看得见"与"按得下"切成两半，可能切错了缝**。c2 也这么切；但 `TaskRoster.Cancel`（`task.go:447`）的返回值 `(bool, string)` 同时带"能不能"与"为什么不能"，理论上"看得见"可以经由它而不必新增出向键。本腿没找到把 `(bool, string)` 递到快照的任何现装读口（`taskRosterState`，`cmd/wisp/panel_pump.go:146-200` 只读 `Look`/`Descendants`/`InFlightSubagents`，不调 `Cancel`——它也不该调，那是会动手的方法），所以维持两半切法，但**登记这一格是我读出来的、不是量出来的**。

## 6. 量不到的地方（具名）

1. **页面侧今天到底声明了哪几枚键**——`frontend/src/lib/panel.ts`（尺 A/B 的另一半输入）在两层禁令内，本腿**未读、未转述**。⇒ §3.2 那张格表给的是**尺的射程**（哪些结构体在册、哪把尺的 `pairs` 里有它），**不是"会不会红"的结论**：射程说"卡面加键落在尺 B 一次"，落不下来还得看那枚 interface 现在声明了什么。这一格**要 owner 或他委托的前端 agent 带问题去问他用的那枚前端 agent**——票派单里"⛔ 不许读 `frontend/**` 源码；如需要它的信息，写进量不到"点名的就是这一格。
2. **两把尺今天的颜色**（在册红几枚）——判它要跑 `go test ./internal/panel`，⛔ 本腿一枚 Go 命令都没跑。167-c2 §5 枚 1 已把这条列为量不到并注"在册读数引用要带锚点、本轮不复量⇒当过期"；本腿维持，不复量、不引用旧锚当现数。⇒ 本腿所有"响/不响"只能与"红/绿"分开写，§3 全节按此口径。
3. **`frontend/dist` 里有没有读 JSON 键的静态产物**——两层禁令覆盖 `frontend/**` 整棵（含 `dist`），本腿**没读、没列目录、没 grep**。本腿只从**仓外元数据**（`.gitignore:22-23`、`git ls-files -- frontend`＝85、`internal/panel/assets.go:63` 的注释）拿到三条布局事实并写进 §3.6：dist 被 ignore（留一枚 `.gitkeep` 锚给 `go:embed`）、页面**源码**才是在版本控制里的那 85 枚。**"dist 内容里有没有 JSON 键消费者"这一问的答案在 dist 文件里，不在 git 元数据里 ⇒ 量不到**。旁证一枚（可读范围内的 Go 尺）：`frontend_hygiene_test.go:290-292` 自己**跳过 `/dist/`**，`tools/d22scan` 的扫描面按 `main.go:23-24`/`:40` 的声明走——⇒ Go 侧仪器里没有一把以 dist 的键名为输入，但**这不等于 dist 里没有消费者**（WebView 读的正是 dist）。
4. **`tools/d22scan` 那三族（ban #6／ban #8／emoji 面）今天扫不扫得到新键**——要跑 `sh scripts/d22scan.sh`，禁。§3.6 只写"扫描面声明里有 `frontend/`"这一条读自 `main.go` 的事实，**不写"会不会红"**。
5. **`bg.Cancel` 那一行能否编过**——要跑 `go build ./cmd/wisp`，禁。§5 枚 6 已把签名逐枚核对（三处 file:line），但"编译过"与"能跑"是跑出来的，不是读出来的。
6. **票面 AC#1 要求的"逐枚现量出口状态"里那三枚今天实不实现得了**——本件只答 Q-51 与停止的"看得见/按得下"缝，占用条（AC#2）的分子分母、草稿（AC#5）、崩溃自救（AC#6）三格**一枚未量**（归 167-c2 §1 与后续票），本腿不重复普查、也不替它们说"已量过"。

## 7. Q-51 的答案一句话

**必须新增**——不能复用任何现有键（`ApprovalCardView` 那 10 枚键里**零枚整数槽**，唯一在快照路线上恒空的 `callChain` 是 `[]string` 且被 `pump.go:84-86`／`panel_pump.go:50` 两处明文挡着不许填假话），而新增这一枚的冲突面是**一把尺、一次声明**（尺 B `approval_test.go:105`/`:118`，尺 A 的 `pairs` 里根本没有 `ApprovalCardView`），⛔ **不是 167-c2 写的两把**；凭据行＝`internal/panel/composer_test.go:55-66` 与 `internal/panel/approval_test.go:113-121` 那两张 `pairs` 字面量名册（本腿逐枚抄过）＋`internal/panel/approval.go:39-59` 的键册。
**仍要 owner 答的那半**不是"撞不撞"（这一半已答：不撞键名、只撞尺 B 一次），而是"**谁有权在同一批里写页面那份声明**"——那一半的输入在 `frontend/**` 里，本腿两层禁令不许读（§6 枚 1）。

## 8. 交件判语

- **只读**：**起手**（§0，`3736f0dd`）`git status --porcelain cmd internal tools docs scripts`＝**空**。⚠ **交件时不再为空**，且**没有一行是本腿写的**：此刻五根上是 `M cmd/wisp/config_reload.go`／`M cmd/wisp/config_reload_223_test.go`／`M internal/tools/paths_shortname_252_r1_test.go`＋`?? cmd/wisp/config_readers_255.go`／`?? cmd/wisp/config_receipt_255_test.go`／`?? internal/config/tiers_app_255r2_test.go`——正是派单点名"在飞"的那三枚写腿（`cmd/wisp`／`internal/config`／`internal/tools`＋`internal/risk`）的工作面。
  **本腿自己的面已逐枚核过**：`git show --name-only` 对 `bdbf2a08`/`bcce2b8a`/`e569a507`/`62ca6c34` 各发一次 ⇒ 四发**只含 `.scratch/wisp/probes/167/c3/census.md` 与同目录 `msgN.txt`**。
  ⛔ 本腿未 `add -A`/`.`/`-a`，每发都是 `git add -- <两枚路径> && git commit -F <msg.txt> -- <同两枚路径>` 串成一条命令（索引共享期未带走别人的已 add 面）。
- **零 Go 命令**：未跑 `go test`/`go build`/`go vet`/`go run`/`gofumpt`/`d22scan`/`staticcheck` 中的任何一枚。所有需要跑才能拿到的读数都具名进了 §6，⛔ 未用推测填空。
- **未 push**：本腿共 5 枚 commit——骨架＋§0 `bdbf2a08` · §1 `bcce2b8a` · §2 `e569a507` · §3＋§4 `62ca6c34` · §3.6＋§5＋§6＋§7＋§8＋两处自我更正＝含本节的末枚。`git show --name-only` 逐枚核过：前四发各只含 `census.md`＋同目录一枚 `msgN.txt`。**全部只在本地**，本腿无一次 `git push`。
- **AC 框一枚未碰**：`git status --porcelain .scratch/wisp/issues` 交件读数＝`M 244-…`＋`?? 259-…`（**都是他程**），**167 那枚票面不在其中**；本腿 pathspec 自始至终只含 `.scratch/wisp/probes/167/c3/**`（逐枚核过，见上一条）。
- **禁区**：`frontend/**`（含 `dist`）与 `design/**` **未读、未在正文转述内容**；本件出现的只有仓库元数据（`.gitignore` 行、`git ls-files` 计数）与 Go 侧尺体里指向它们的路径字符串。冻结件（`docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`／`.github/workflows/ci.yml`）一字未动，其中 `l2_grant_boundary_test.go` 本腿**只读尺体**（`:1720-1726`、`:1779`、`:1807`、`:1878-1906`）。
- **临时件只建不删**：`.scratch/wisp/probes/167/c3/` 下 `census.md`＋`msg1..msg5.txt` 五枚消息件全部留存，无删除动作。
- **给下一程的一句话**：本件 §3.2 那张格表是"按落点算尺"的，不是"按尺算落点"的——**分子分母照旧归编排者后裁**（编排者在 167-c2 收档时的裁定有效），本腿只把 Q-51 从"没人量过"改成"量过、且成本是一枚 interface 不是一棵前端树"。
