# 35-a3 · 页面要读的字段 与 Go 侧真有的字段 —— 逐枚对账（只读数、不推测）

腿号 `35-a3`（只读普查腿）· 落笔时刻 **2026-10-07 13:40 +0800**（`date` 现量，非按上一条戳推）
起手锚＝同目录 `anchor.txt`（commit `6a976294`）
两把尺的锚点：
- **会被构建那棵树＝`HEAD`＝`dev`＝`518b89cf`**（提交时刻 2026-10-07 13:30:36 +0800）。依据＝票面 `.scratch/wisp/issues/35-panel-bridge-c17.md:115` 逐字："**接收器／桥的落地判据只能钉在 dev 的树上，⛔ 不许拿分支的行号当凭据**"。
- 页面那棵树＝`dsh/feat/frontend-p0-v2`＝`16c2f038`（2026-09-28 10:18:53 +0800，比 dev 旧 9 天）；`git rev-list --left-right --count`＝dev 领先 **1633** 笔、分支领先 **21** 笔；merge-base `c1b10089`。

本文所有行号形如 `<树或ref>:<file>:<line>`。裸 `file:line` 在本项目不算可核锚点，本文一处都没有。
每条引用前都跑过存在性尺 `git show <ref>:<file> | wc -l`（枚数与行数见 §7）。

---

## 0. 一句话先给结论

**页面对账的两棵树共用同一枚契约件**（`frontend/src/lib/panel.ts` 在 `HEAD` 与 `16c2f038` 上**字节相同**，`diff` 输出 **0 行**、各 **311 行**），
所以"页面读什么"这一侧**两棵树是同一份答案**；分歧全部在 Go 那侧与页面装配层：
dev 的 Go 快照比这份契约件**多 7 枚键**（含子键共 52 枚），页面**一个字都没读**；
页面**读了 2 枚 Go 没有的东西**（`snapshot.view`、入向方法名 `panel.approval.request`），另有 **2 枚页面假设的宿主对象 Go 从不安装**；
而**页面发出的全部 5 枚请求今天一跳都到不了 Go**（传输层断裂，票面已裁形＝甲子形①，见 §4）。

---

## 1. 页面真读了哪些键（逐枚，带锚；静态文案不算）

### 1.0 先剔掉不算对接的起点

派单 §5.1 给的候选起点里，**六枚里有一枚以上根本不从宿主拿数据**，具名报回（尺＝`git grep -o -h "snapshot\.[a-zA-Z_]+|snap\.[a-zA-Z_]+" <ref> -- frontend/src` 全树逐文件计数＋`git grep -n` 命中行逐行读全文）：

| 候选起点 | 真读宿主字段吗 | 现量 |
|---|---|---|
| `frontend/src/App.tsx` | **读** | `snapshot.*` 4 枚读取点在下表 |
| `frontend/src/main.tsx` | **不读**（挂载＋演示横幅） | `dsh/feat/frontend-p0-v2:frontend/src/main.tsx:96-97` `const ROOTS = { product: { … node: <App /> } }`——**无 prop 挂载**；全件读的是 `location.search`／`document`，零宿主键 |
| `components/ai-native/agent-screen.tsx` | **不读** | 名册尺里 **0 命中**（`snapshot.` 与宿主 interface 名都没出现） |
| `components/ai-native/approval-card.tsx` | **不读** | 全件只有 `dsh/feat/frontend-p0-v2:frontend/src/components/ai-native/approval-card.tsx:67` 一枚 `window.setTimeout`＝动效 |
| `components/ai-native/chat-composer.tsx` | **不读** | 0 命中；真正接宿主的是**另一枚** `components/composer.tsx` |
| `components/ai-native/context-cards.tsx` | **不读** | 0 命中 |
| `components/ai-native/records-table.tsx` | **不读** | 11 枚 `window.` 命中全是 `innerWidth`／`addEventListener("pointer*")`＝拖拽定位 |
| `components/ai-native/insight-cards.tsx` | **不读，且自己写了它在等** | `dsh/feat/frontend-p0-v2:frontend/src/components/ai-native/insight-cards.tsx:15` 逐字 "…the cost and context meters are the place it would go, and those are **waiting on the snapshot fields (ticket 167)**" ⇒ 本文只登记它"期待谁喂"，不把它算进读取数 |

**真正从宿主拿数据的文件＝9 枚**：`lib/panel.ts`（契约声明本身，不算读取）、`lib/panel-views.ts`、`App.tsx`、`components/composer.tsx`、`components/approval-screen.tsx`、`components/chat-screen.tsx`、`components/l2-approval-card.tsx`、`components/result-stream.tsx`、`components/git-branch.tsx`（**最后一枚今天没被喂**，见 §3.1）。
`fixtures/harness.ts`·`fixtures/harness-app.ts`·`components/showcase.tsx`·`components/harness/**` 是 `?harness=1/2` 的演示数据，按派单口径**一律不算**（`dsh/feat/frontend-p0-v2:frontend/src/main.tsx:19-21` 逐字："the product panel has no fixture mode"）。

### 1.1 快照顶层（`PanelSnapshot` 的读取点）

| 页面读的键 | 读取锚（页面树） | 谁喂它 |
|---|---|---|
| `snapshot.composer` | `dsh/feat/frontend-p0-v2:frontend/src/App.tsx:92`（`snapshot.composer ?? EMPTY_COMPOSER`）；dev 同义件锚＝`HEAD:frontend/src/App.tsx:77` | Go `Snapshot.Composer` |
| `snapshot.pending` | `…App.tsx:206`／`:208`／`:321`／`:322`（`snapshot.pending.length` 与 `.map`）；dev＝`HEAD:frontend/src/App.tsx:91`／`:96`；另 `…components/approval-screen.tsx:68` | Go `Snapshot.Pending` |
| `snapshot.results` | `…App.tsx:431`／`:457`／`:461`；另 `…components/chat-screen.tsx:25` | Go `Snapshot.Results` |
| `snapshot.view` | `…App.tsx:123`（`(snapshot as { view?: string }).view`）＋ `…lib/panel-views.ts:98-99`（`currentView(snapshot)` 读 `snapshot.view`）；dev＝`HEAD:frontend/src/App.tsx:85-86` | **Go 没有这一枚**（§3.1） |
| `generatedAt` | 声明＝`HEAD:frontend/src/lib/panel.ts:136`；**读取点＝0**（全树 `snapshot.generatedAt` 命中 0；`generatedAt` 4 处命中＝声明＋`EMPTY` 字面量 `HEAD:frontend/src/App.tsx:53`） | Go 有、页面不读（§2.2） |

### 1.2 审批卡（`ApprovalCardView`，页面读 9／声明 10）

页面侧锚全部在 `dsh/feat/frontend-p0-v2:frontend/src/`：
`components/l2-approval-card.tsx:207`（`view.reasonKnown`＋`view.reason`）·`:265`（`requestApprovalResolution(view.correlationId, outcome)`）·`:266`（`onIntent?.(view.correlationId, …)`）·`:312`（`view.rulesHit`）·`:316`/`:318`（`view.callChain`）·`:321`（`view.sessionOverrideBlocked`）·`:338`（`view.decidedBy`）·`components/approval-screen.tsx:58`（`view.correlationId`）·`:61`（`view.rulesHit.length`）·`:104`（`key={view.correlationId}`）·`App.tsx:322`（`key={v.correlationId}`）。

- **读了但页面自己没声明的键＝0**；**声明了但没人读的键＝1**：`tool`（`lib/panel.ts:28` 声明，读取点＝0，命中计数尺见 `.scratch/wisp/probes/35/a3/raw-keyreads-bothtrees.txt`）。
- 结果流：`components/result-stream.tsx:34` 与 `App.tsx:432`/`:462` 读 `chunk.correlationId`；`text`/`done` 由 `chat-screen.tsx` 消费。

### 1.3 输入行（`ComposerState` 及三个子节）

`dsh/feat/frontend-p0-v2:frontend/src/components/composer.tsx`：`:139`/`:140`（`state.attachmentError`）·`:150`/`:152`（`state.workspace.set`/`.canonical`/`.reason`）·`:154`（`.reparse`）·`:163`（`.spelling`）·`:199`（`state.acceptedAttachmentMimes`＋`state.maxAttachmentBytes`）·`:205`（`.join(",")` 进 `accept=`）·`:276`（`state.mode.l2ConfirmNames.includes(m)`）·`:318`（`a.deduplicated`）；`App.tsx:243`/`:303` 也读 `composer.workspace.set`/`.canonical`。
`rewritten`（`lib/panel.ts:89` 声明）**读取点＝0**。

### 1.4 页面发出的入向请求（页→Go；这一跳的"键"同样是字段）

件＝`frontend/src/lib/panel.ts`，**两棵树字节相同 ⇒ 行号两树一致**（下表同时是 `HEAD:` 与 `dsh/feat/frontend-p0-v2:`）：

| 页面发的方法名 | 发送锚 | 封套键 | Go 是否答 |
|---|---|---|---|
| `panel.approval.request` | `HEAD:frontend/src/lib/panel.ts:179-186` | `method`／`correlationId`／`outcome`（`"grant"`/`"refuse"`，`:50`） | **不在名册**（§3.1） |
| `panel.mode.request` | `HEAD:frontend/src/lib/panel.ts:246` | `method`+`requestId`+`source`（`:218-225`）+`to` | 答 |
| `panel.workspace.request` | `HEAD:frontend/src/lib/panel.ts:255` | 同上＋`path` | 答 |
| `panel.attachment.add` | `HEAD:frontend/src/lib/panel.ts:272-277` | 同上＋`name`/`declaredMime`/`sizeBytes`/`dataBase64` | 答 |
| `panel.message.send` | `HEAD:frontend/src/lib/panel.ts:300` | 同上＋`text`/`attachments` | 答 |
| （无）`config.get`／`config.set` | 页面**零发送点**（尺＝`git grep -n -E 'config\.(get|set)' <ref> -- frontend/src`＝**0 命中**） | — | Go 答（§2.2） |

页面假设的宿主对象：`HEAD:frontend/src/lib/panel.ts:152-155`（`window.wispBridge ?? window.chrome?.webview ?? null`）＋ `:140-142`（`postMessage(message: string)`）。

---

## 2. Go 侧真有什么（dev＝会被构建那棵树）

### 2.1 快照本体

`HEAD:internal/panel/composer.go:57-92`＝`type Snapshot struct`，**6 枚键**：
`:58` `pending` · `:59` `results` · `:60` `composer` · `:61` `generatedAt` · `:74` `instructions,omitempty`（指针） · `:91` `tasks,omitempty`（指针）。
泵只在挂了读者时填后两枚：`HEAD:internal/panel/pump.go:304`（`if p.src.Instructions != nil`）→ `:308`；`:310`（`if p.src.Tasks != nil`）→ `:327`。
`HEAD:internal/panel/composer.go:235-269`＝`ComposerState`，**11 枚键**：`:236` `mode` · `:237` `workspace` · `:238` `attachments` · `:239` `acceptedAttachmentMimes` · `:240` `maxAttachmentBytes` · `:241` `attachmentError` · `:246` `git` · `:255` `currentModel` · `:268` `credentialState` · `:269` `credentialKnown`（＋`:256` `modelKnown`）。
子节：`ModeView` `:141-149`（3）·`WorkspaceView` `:212-224`（6）·`AttachmentRef` `HEAD:internal/panel/attachments.go:84-103`（9）·`ResultChunk` `HEAD:internal/panel/composer.go:96-98`（3）·`ApprovalCardView` `HEAD:internal/panel/approval.go:40-58`（10）·`GitView`＋`GitWorktree` `HEAD:internal/panel/git.go:96-144`（10＋5）·`InstructionsSection`＋文件行 `HEAD:internal/panel/instructions_200.go:46-79`（3＋8）·`TaskRosterSection`＋`TaskRowView` `HEAD:internal/panel/subagent_roster_197.go:107-168`（7＋12）。`CredentialState` 是**字符串型**（`HEAD:internal/panel/config_handlers.go:175` `type CredentialState string`），无子键。
逐枚名册（112 行，含每枚的 `ref:file:line`）＝`.scratch/wisp/probes/35/a3/raw-head-go-jsontags.txt`。

入向名册（Go 答 6 枚）：`HEAD:internal/panel/bridge.go:42`／`:43`／`:44`／`:45` 四枚 `panel.*` ＋ `:66` `config.get`／`:67` `config.set`；守卫 `:146-152`（`knownComposerMethod`，`case` 在 `:148` 一次列全 6 枚）；解析 `:126-144`（`ParseComposerRequest`，三段拒绝：名册 `:132`、来源 `:135`、缺 `requestId` `:139`）。封套字段 `:92-120`（`method`/`requestId`/`source`/`to`/`path`/`text`/`attachments`/`AttachmentPayload`/`configField`/`configProvider`/`configModel`/`configValue`）。
页面页 Go 路由器：`HEAD:internal/panel/composer_dispatch.go:153`（`Handle`）→`:175`（`dispatch`）→ `case` 六枚在 `:177`/`:182`/`:187`/`:192`/`:197`/`:202`。

### 2.2 五种情况计数（**枚数全为我现量的键级读数**）

| 情况 | 枚数 | 说明 |
|---|---|---|
| **同名同义** | **41** | 审批卡 10 ＋ 结果块 3 ＋ `mode` 3 ＋ `workspace` 6 ＋ 附件 9 ＋ `ComposerState` 交集 6 ＋ 快照顶层交集 4。逐枚名册对拉＝`raw-head-go-jsontags.txt` × 页面声明 `HEAD:frontend/src/lib/panel.ts:24-137` |
| **同名不同义** | **0** | 判据：41 枚同名键的**类型形**逐枚对（`[]string`↔`string[]`、`int64`↔`number`、`bool`↔`boolean`）＋语义按两侧原文注释对。**未跑运行时**（⛔ 零 go 命令），所以这一格是"读码判"，不是"实测判"（§6） |
| **异名同义** | **7** | 型名 5 对：`Snapshot`↔`PanelSnapshot` · `ModeView`↔`ComposerMode` · `WorkspaceView`↔`ComposerWorkspace` · `AttachmentRef`↔`ComposerAttachment` · `ResultChunk`↔`ResultChunkView`（这 5 对＋同名的 `ApprovalCardView` 正是 `HEAD:internal/panel/composer_test.go:60-65` 那张对照表的 6 行）。另 2 对跨层：①传输对象——页面 `window.wispBridge`/`window.chrome.webview.postMessage`（`lib/panel.ts:140-155`）↔ Go 绑定名 `wispDispatch`＋`dispatchRaw`（`HEAD:cmd/wisp/panel_host_windows.go:80`/`:630`）；②审批答复——页面 `panel.approval.request`＋`outcome:grant/refuse` ↔ Go 面板路由动词 `panelAllow`/`panelReject`（`HEAD:cmd/wisp/approval_reply.go:185`，经 `Gate.DecideFromPanel`，`HEAD:internal/agent/approval/gate.go:640-642`） |
| **Go 有、页面不读** | **15 项**（展开子键 **45**） | 快照 2：`instructions`（子 11）、`tasks`（子 19）；输入行 5：`git`（子 15）、`currentModel`、`modelKnown`、`credentialState`、`credentialKnown`；`generatedAt`（顶层，声明了但读取点 0，§1.1）＝页面"半读"；入向名册 2：`config.get`、`config.set`（页面零发送点）；信封 4：`configField`/`configProvider`/`configModel`/`configValue` |
| **页面读、Go 没有** | **4** | `snapshot.view`（§3.1-a）· 入向方法 `panel.approval.request`（§3.1-b）· 宿主对象 `window.wispBridge`（§3.1-c）· `window.chrome.webview`（§3.1-c 同族） |

### 2.3 跨树读数（页面分支的 Go 侧比 dev 少一截，别拿它当真身）

| 结构 | `16c2f038`（页面分支） | `HEAD`（dev，会被构建） |
|---|---|---|
| `Snapshot` 键 | **4**（`16c2f038:internal/panel/composer.go:58-61`，无 `instructions`/`tasks`） | **6**（`:58-61`＋`:74`＋`:91`） |
| `ComposerState` | **6**（`16c2f038:internal/panel/composer.go:200` 起，无 `git`/`currentModel`/`credentialState`） | **11** |
| 入向名册 | **4 枚**（`16c2f038:internal/panel/bridge.go:42-45`；该件共 **125 行**，`ParseComposerRequest`＝`:84`、`knownComposerMethod`＝`:104`） | **6 枚**（`HEAD:internal/panel/bridge.go:42-45`+`:66-67`；该件 **167 行**） |
| 路由／宿主／CLI 入向腿 | **全树不存在**：`internal/panel/composer_dispatch.go`·`inbound_roster_253_test.go`·`git.go`·`config_handlers.go`·`instructions_200.go`·`subagent_roster_197.go`·`cmd/wisp/panel_inbound.go`·`cmd/wisp/panel_host_windows.go`·`cmd/wisp/panel_resident_windows.go` 逐枚 `git cat-file -e`＝**NO** | 全在 |

⇒ **推论（不是推测，是枚数对拉的直接结果）**：页面那份契约件今天**恰好等于分支的 Go 侧**，而 dev 的 Go 侧已经长出 7 枚键；**把分支合过来不会补上这 7 枚**，因为缺的是页面声明那一侧（`panel.ts` 两树相同）。

---

## 3. 每一枚"页面读、Go 没有"缺的那一跳（四档选一档）

档位名照派单：**①字段没进快照／②进了快照但没泵出／③泵出但名字或类型不对／④根本没有入向通道**。

**a) `snapshot.view`（页面 `…App.tsx:123`＋`…lib/panel-views.ts:98-99`；`HEAD:frontend/src/App.tsx:85-86`）**
**断在①字段没进快照**。判定行＝`HEAD:internal/panel/composer.go:57-92` 全结构读完，**6 枚键里没有 view**；`HEAD:internal/panel/pump.go:234`（`func (p *SnapshotPump) Snapshot()`）的装配段（`:275`/`:284`/`:287`/`:304`/`:310`）也没有任何 view 读者。
旁证（页面自陈）：`dsh/feat/frontend-p0-v2:frontend/src/lib/panel-views.ts:89-92` 逐字 "The field does not exist on PanelSnapshot yet - **ticket 35's pump owns it**"。

**b) 入向方法名 `panel.approval.request`（页面 `HEAD:frontend/src/lib/panel.ts:179-186`）**
**断在④根本没有入向通道**。判定行＝`HEAD:internal/panel/bridge.go:146-152`：`knownComposerMethod` 的 `case` 只列 6 枚，**不答这一枚** ⇒ 走 `:132-134` 拒。
尺（⛔ 不是"按某层符号名扫另一层"）：`git grep -n -E "panel\.approval\.request|approval\.resolve|ApprovalRequest" HEAD -- '*.go'` 的**非 test 命中＝0**（全仓命中只有 `HEAD:internal/panel/composer_test.go:417` 一枚名册字面量与 `HEAD:internal/panel/frontend_hygiene_test.go:32` 一句注释）。
⚠ 这一档还叠了第二层，**必须一起写**：即使把名字加进名册（＝C17 契约变更，人工批准面），页面的 `outcome:"grant"` 那一支**结构上不可能被满足**——`HEAD:internal/agent/approval/gate.go:640-642` 逐字 "DecideFromPanel **refuses every allow** on the route alone"，且 `HEAD:cmd/wisp/approval_reply.go:189-192` 写明 "PanelAPI has no Allow method"。**这是裁决，不是缺口**（AGENTS.md §1.2「面板侧来源的 L2『允许』」）。
⚠ 还有第三层：页面这枚信封**不带 `requestId` 与 `source`**（对比 `lib/panel.ts:218-225` 的 `sendRequest` 才带），而 `HEAD:internal/panel/bridge.go:135`（来源）与 `:139`（缺 requestId）两道拒绝先于路由 ⇒ 属于**③名字对不上之外的"封套形对不上"**，一并记入。

**c) 页面假设的宿主对象 `window.wispBridge` / `window.chrome.webview`（`HEAD:frontend/src/lib/panel.ts:140-155`）**
**断在④根本没有入向通道**（传输层），且**注释描述的机制不存在**。判定行：
- 页面全树对绑定名零引用：`git grep -c wispDispatch <ref> -- frontend` 在**两棵树都＝0 命中**。
- Go 侧只绑不接：`HEAD:cmd/wisp/panel_host_windows.go:80`（`panelDispatchBinding = "wispDispatch"`）·`:405-406`（`w.Bind(panelDispatchBinding, func(raw string) string { reply, _ := m.dispatchRaw(ctx, raw) … })`）·`:630`（`dispatchRaw` 本体）。这条链的**第一跳只能由 Go 测试直接调**，页面从不叫（票面 `…/35-panel-bridge-c17.md:93-97` 的协议链已裁，本腿复跑的码证＝上一条＋"页面对 `wispDispatch` 零引用"）。
- 注释过期：`HEAD:frontend/src/lib/panel.ts:146` 写 "Installed by WebView2's **AddHostObjectToScript** / postMessage pipe"，而 `AddHostObjectToScript` 全仓＝0（票面 `:109` 的现量，本腿未复跑该把尺，具名标为**转述**）。

**d) 出向（Go 快照 → 页面）整条：档位②**
所有 §1 读到的键**都进了快照**（`HEAD:internal/panel/composer.go:57-92`），泵也**真的 marshal**（`HEAD:internal/panel/pump.go:337` `Marshal`／`:348` `Publish`；调用点 `HEAD:cmd/wisp/panel_pump.go:405`），**但字节停在运行时**：`HEAD:cmd/wisp/panel_pump.go:321`（`rt.lastSnap, rt.lastSnapBytes, rt.snapSeen = snap, data, true`），读者 `:383`（`func (rt *agentRuntime) lastPanelSnapshot()`）。
非 test 推送到页面的语句＝**0**：`git grep -n -E '\.Eval\(|\.Init\(|PostWebMessage|CreateWebMessageAsJson' HEAD -- cmd internal ':!*_test.go'` 命中 **0**；同尺同路径只命中 `\.SetHtml(` **3 枚**＝`HEAD:cmd/wisp/panel_host_windows.go:449`/`:468`/`:681`（一次性装载 HTML，不是推送）。
⇒ 页面今天**没有任何一格能被真数据画亮**；`?harness=1/2` 画亮的都是 fixture（`HEAD` 与分支的 `main.tsx` 挂载点＝`HEAD:frontend/src/main.tsx:58`、`dsh/feat/frontend-p0-v2:frontend/src/main.tsx:96-97`，**两处都无 prop**）。

---

## 4. 票 35 落地判据该钉哪一行（35-r1 用）

### 4.1 入向那几枚方法名的当前真身行号

⛔ **只准钉 `HEAD`（＝dev）那一棵**：分支那 6 枚文件里有 5 枚**整文件不存在**（§2.3），钉它＝钉空气。
唯一两枚**两棵树行号相同**的入向锚＝`frontend/src/lib/panel.ts` 的发送点（件字节相同）：

| 要钉的东西 | dev 真身（可用） | 分支对照（⛔ 不作凭据） |
|---|---|---|
| 名册 6 枚 | `HEAD:internal/panel/bridge.go:42`／`:43`／`:44`／`:45`／`:66`／`:67` | `16c2f038:internal/panel/bridge.go:42-45`（**只有 4 枚**） |
| 守卫 | `HEAD:internal/panel/bridge.go:146`（函数）·`:148`（`case` 一列 6 枚） | `16c2f038:internal/panel/bridge.go:104` |
| 解析＋三道拒绝 | `HEAD:internal/panel/bridge.go:126`／`:132`／`:135`／`:139` | `16c2f038:internal/panel/bridge.go:84` |
| 路由器 | `HEAD:internal/panel/composer_dispatch.go:153`／`:175`／`:177`／`:182`／`:187`／`:192`／`:197`／`:202` | **文件不存在** |
| WebView2 绑定 | `HEAD:cmd/wisp/panel_host_windows.go:80`／`:405`／`:406`／`:630` | **文件不存在** |
| CLI 入向腿 | `HEAD:cmd/wisp/panel_inbound.go:103`（`cmdPanelInbound`）·`:163`（`disp.Handle(ctx, raw)`） | **文件不存在** |
| 页面发送点（判据的"原料串"从这里取） | `HEAD:frontend/src/lib/panel.ts:179-186`（approval 信封）·`:218-225`（`sendRequest` 封套）·`:246`·`:255`·`:272-277`·`:300` | 同行号（`diff`＝0 行） |

**判据必须由页面自己报回**（票面 `…/35-panel-bridge-c17.md:63-68` 那枚新框）⇒ 原料串请逐字取
`{"method":"panel.approval.request","correlationId":…,"outcome":…}`（`HEAD:frontend/src/lib/panel.ts:179-186`）
与 `{"method":"panel.mode.request","requestId":"pc-<n>-<uuid>","source":"panel-composer","to":…}`（`:218-225`＋`:246`＋`:205-209`）。
⚠ **本腿给一枚票面没写的提醒**：那枚新框写"取 `panel.approval.request` 的原始串"，可**这一枚今天过不了名册守卫**（§3-b）。若照框直落，`dispatchRaw` 到达性断言会**永远红**——这不是坏消息，这是**唯一能让判据有牙的形状**，但写腿要么把它写成"预期被拒且拒因可指认"，要么由编排者把这枚信封换成名册内的 `panel.mode.request` 做正控＋把 approval 那枚留作负控。**⛔ 本腿不代拍**（选形与名册都属 C17 人工批准面）。

### 4.2 会撞到的既有钉（全部 `grep` 现量，⛔ 未跑测试）

| 钉 | 锚 | 为什么会被撞到 |
|---|---|---|
| 闭名册尺 | `HEAD:internal/panel/inbound_roster_253_test.go:67-75`（名册 6 枚字面量）·`:419`（`TestFullInboundMethodRosterIsClosed`）·`:434`/`:437`/`:441`/`:446`/`:456`/`:465`/`:476`（五个分母） | 新增点分法名 ⇒ 当场红（票面 `:149` 已列禁区） |
| "渲染器只许一扇门" | `HEAD:internal/panel/composer_test.go:502`（`TestTheRendererHoldsExactlyOneDoorToTheHost`）·`:510`（host call sites 要＝2）·`:518`（method 必须是字面量）·`:522`（页面不许命名 Go 不答的路由）·`:533`（正控 `TestPlantedRendererDoorShapesGoRed`） | 落地若在页面加第二通道或新路由名 ⇒ 这枚先红 |
| 双向键对账 ×2 | `HEAD:internal/panel/composer_test.go:48`（`TestComposerContractTypesMatchFrontend`，对照表 `:60-65`，两向 `t.Errorf` 在 `:73-78`，**没有豁免支**）·`HEAD:internal/panel/approval_test.go:105`（`:126-132` 同形） | dev 的 `Snapshot` 比 `PanelSnapshot` 多 `instructions`/`tasks` 两枚 ⇒ 这枚尺**今天要么已红、要么靠我没找到的某处兜住**；本腿⛔不许跑 go test，判不动（§6-b） |
| 快照四键字节钉 | `HEAD:internal/panel/pump_test.go:123` 与 `:291`：逐字 `!= "composer,generatedAt,pending,results"` | 页面读 `view` ⇒ 谁若"顺手"给快照加第 5 枚非 omitempty 键，这两枚立刻红 |
| 名册源码抽取器 | `HEAD:internal/panel/git_test.go:507-510`（把 `bridge.go` 四枚 `panel.*` 常量当源码文本比对）·`:394`（`panelMethodRe`） | 改 `bridge.go` 常量区的**排版／行数**都会打到它 |
| 封套接受／拒绝面 | `HEAD:internal/panel/bridge_test.go:51`（`TestComposerEnvelopeAcceptsItsFourRequests`）·`:77`·`:131`（`TestFrontendComposerRequestsMatchTheEnvelope`） | "Four" 这枚名字与 6 枚现实差 2；改名册必动它 |
| 路由不越权族 | `HEAD:internal/panel/composer_dispatch_test.go:126`/`:201`/`:231`/`:249`/`:374`/`:434`（`TestPanelHostIsAttachedAndNamesTheWindowHops`）/`:708` | 票面 `:139` 说乙形在 `:434` 上"真红"；甲形若碰到 hop 命名同样会红 |
| 出向吞 `nil` 的空桩 | `HEAD:internal/panel/composer_dispatch_test.go:461`（票面 `:107`/`:151` M3 具名：它拼的是 `PostWebMessageAsJson`，库里真名 `PostWebMessageAsJSON`） | "包绿"≠"页面收到"，出向判据不许绕它 |
| "测试扮演页面"族（N6，⛔ 不得抵这一格） | 票面 `:141` 列 `cmd/wisp/panel_resident_windows_test.go:301`/`:801`/`:805`/`:808`/`:864`（**本腿未复跑此把尺**，⛔ 不许当行号凭据引用） | 甲落地后它们照旧绿＝零分辨力 |
| CLI 入向腿 AC#9 | `HEAD:cmd/wisp/panel_inbound_33_test.go:99`/`:125`/`:158`/`:182`/`:393` | 动 `panel_inbound.go` 的签名或拒绝形状会打到 |
| 前端卫生尺 | `HEAD:internal/panel/frontend_hygiene_test.go:32` | ⛔ 不动 `frontend/**` 一字（机主只放读） |

---

## 5. 对账结论一句话 ＋ 照今天这张表落地后仍拿不到输入的格子

**一句话**：**字段名这一层两棵树今天对得上 41 枚、缺的是"跳"不是"名"**——出向缺**泵到页面的那一句**（`HEAD:cmd/wisp/panel_pump.go:321` 之后零推送语句）、入向缺**页面串的到达**（`HEAD:internal/panel/bridge.go:146-152` 不认页面那枚方法名）、页面缺**耳朵与 5 枚键的声明**（`HEAD:frontend/src/lib/panel.ts:112-137` 一个字都没写 `git`/`currentModel`/`credentialState`/`instructions`/`tasks`）。

照今天这张表落甲子形①之后，**仍拿不到输入的格子**（宁可缺，不许造）：

1. **L2 卡的「允许」按钮**——永远拿不到输入：`HEAD:internal/agent/approval/gate.go:640-642`（面板路由对 allow 一律拒）＋`HEAD:cmd/wisp/approval_reply.go:189-192`（PanelAPI 无 Allow）。页面 `outcome:"grant"` 那一支（`HEAD:frontend/src/lib/panel.ts:50`）**无解**，除非人改裁决。
2. **审批答复整跳**（`panel.approval.request`）——不在名册＝C17 人工批准面（AGENTS.md §1「改契约＝人工批准」），甲子形①**不解决它**，只解决"名字在册的四枚能不能到"。
3. **九屏里七屏无输入可画**：`dsh/feat/frontend-p0-v2:frontend/src/lib/panel-views.ts:65-75` 的名册逐枚——`fed: true` 只有 `chat`／`approval` 两枚，其余 **7 枚 `fed:false`**（`palette`/`tasks`/`ball`/`config`/`security`/`privacy`/`cost`），其中 `config` 带 `selfFed:true`。`tasks` 屏虽然 dev 已经有 `snapshot.tasks`（19 枚子键），**页面侧连声明都没有**，仍然画不出来。
4. **git 分支 chip**：`dsh/feat/frontend-p0-v2:frontend/src/App.tsx:217` 逐字 `view={{ current: "", isRepo: false, local: [], worktrees: [] }}`——**手传字面量**，而且这枚字面量的形（`current`/`isRepo`/`local`/`worktrees`）**与 Go 的 `GitView`（`kind`/`reason`/`branch`/`detachedSha`/`isDetached`/`repoRoot`/`currentWorktree`/`worktrees`/`branches`/`switchBlocked`，`HEAD:internal/panel/git.go:119-144`）不同名也不同形**；dev 树更彻底——`HEAD:frontend/src/App.tsx` 共 117 行，**根本没 import GitBranchChip**。
5. **模型格与密钥格**：Go 已填 `currentModel`/`modelKnown`/`credentialState`/`credentialKnown`（`HEAD:internal/panel/pump.go:284-296`），页面**零声明零读取点**。
6. **视图选择（`view`）**：Go 侧零键（§3-a），页面每次回落 `DEFAULT_VIEW`（`dsh/feat/frontend-p0-v2:frontend/src/lib/panel-views.ts:80`）。

---

## 6. 未做／判不动（具名）

**a) 既有表我一份都没读**——派单 §4 点名的 `docs/evidence/s1/145-snapshot-field-census-r1.md`／`145-snapshot-growth-r2.md`／`-r3.md`／`145-snapshot-fields-landed-r1.md`＋`…-accept-r1.md`／`180-182-panel-fields-census-c1.md`／`180-panel-width-ownership-a1.md`／`167-panel-widgets-r1.md`／`.scratch/wisp/probes/253/`（只 `ls`＝`p4`·`r1`·`r3`·`r5` 四子目录，未读）／`.scratch/wisp/probes/242/`（未读）。
⇒ 本文的 Go 侧字段全部**现量自代码**（112 行名册在 `raw-head-go-jsontags.txt`），**没有与那几份表对拉过**；"如果发现某份既有表和盘上代码不一致就具名报回"这一条**本腿没执行**（调用预算换成了直接量码）。代价：若那几份表里有比 `HEAD` 更细的口径，本文不带它。

**b) 判不动（⛔ 零 go 命令，是派单禁的，不是我省的）**
- `HEAD:internal/panel/composer_test.go:48` 与 `approval_test.go:105` 这两枚双向尺**今天到底红不红**：`HEAD:internal/panel/composer.go:86-90` 的原文注释自己写着"这枚尺现在会把 `tasks` 报成 interface 未声明的键"，而我在 `composer_test.go:55-80` 区间**没找到豁免支**（`func jsonKeys` 那个 helper 名我猜错过一次：`git grep -A22 'func jsonKeys'`＝**0 命中**，真实 helper 名未知）。⇒ 需要 `go test ./internal/panel/ -run 'Composer|ApprovalCardView'` 才能定；本腿⛔不许跑。
- `go vet`／`go build` 同理未跑（另一条腿持 `cmd/wisp` 整包读数窗口）。

**c) 读了但不足以定性的**
- `components/ai-native/**` 30 枚里我只逐枚扫了派单点名的 6 枚＋`streaming-text.tsx:25`（"One inline source chip, as the snapshot reports it"——**注释声称等快照，页面里没有任何读取点**，其余 23 枚未开）。
- `HEAD:cmd/wisp/approval_reply.go` 的路由表我只读了 `:160-215` 区段（禁读母本之外，仍按行区间读），**没有逐枚抽出 `panelAllow`/`panelReject` 的确切字符串形状与入参键名**；要把它当"审批答复的 Go 侧真名"用，需下一腿补一把按语义变体（`allow|reject|grant|refuse|outcome|correlationId`）的扫。
- `credentialState` 的词表（`HEAD:internal/panel/config_handlers.go:175-180` 起，我只读到 `CredentialUnknown = "unknown"` 一枚常量，其余取值枚数未点）。

**d) 转述未复跑（凡引票面读数处都已在正文标"票面 `:NN`"）**：`AddHostObjectToScript` 全仓＝0、`panel_resident_windows.go:349-372` 那枚 `post()`（`:360`）非 test 唯一 `Dispatch` 点、`composer_dispatch_test.go:461` 空桩的拼写、`go.mod:19` 的 `// indirect`——四把尺本腿**没有复跑**。

**e) 未做（本腿射程内但其实该做）**：没量 `frontend/dist` 与 `internal/panel/assets.go` 的 fail-closed 形状（票面 `:149` 列为禁区，读它也无助于字段对账）。

---

## 7. 自证与卫生

- **dirty 数三枚读数（同一把尺 `git status --porcelain | wc -l`）**：起手 **752**（13:32）→ 交件前 **781**（13:40，内含本腿 29 枚未跟踪件）→ **终态 752**（13:44，本腿 commit 之后）。
  **差集 100% 落在 `.scratch/wisp/probes/35/a3/**`**：终态尺 `git status --porcelain | grep -c 'probes/35/a3'`＝**0**（本腿 30 枚件已全部入库：commit `6a976294` 3 枚＋commit `0f0fdaf2` 30 枚）。
  **但 781→752 这 29 枚不是全归我**——期间**别的腿动了盘**：起手登记的那 8 枚 ` M .scratch/wisp/probes/161/r6/logs/flip-*.txt` 在终态**已从 status 里消失**（`git status --porcelain | grep -c 'probes/161/r6/logs/flip'`＝**0**），说明那批在飞物在 13:40-13:44 之间被别的腿 commit 掉了。**本腿一枚没提交、没还原、没"顺手清理"**别人的东西（派单原文："数字会动，动的人不是你就不用管，但要逐枚量 mtime 具名登记来历"）；两处在飞物都按此登记：
  - ` M .scratch/wisp/probes/161/r6/logs/flip-*.txt`＝**8 枚**（尺＝`git status --porcelain | grep -c '^ M .scratch/wisp/probes/161/r6/logs/flip'`），mtime 现量样本：`flip-1.txt` 2026-10-03 11:17／`flip-2.txt` 11:19／`flip-3.txt` 11:22（各 117,744 字节）⇒ **10-03 那批腿的旧在飞物，本腿不动、不提交、不还原**。
  - `?? design/`＝**11 枚未跟踪条目**（尺＝`grep -c '^?? design/'`），机主只放**读**，本腿零写入。
  - 其余未跟踪目录（本腿非作者，登记即可）：`.scratch/wisp/probes/35/a1`·`a2`（前一／二条普查腿的件）、`.scratch/ci-logs/run-*.log|err`、`.scratch/.scratch/`。
- **两棵树行数复核（第二遍，同一把尺 `git show <ref>:<path> | wc -l`）**：`HEAD:frontend/src/App.tsx`＝**117**、`HEAD:frontend/src/main.tsx`＝**61**、`16c2f038:frontend/src/App.tsx`＝**531**、`16c2f038:frontend/src/main.tsx`＝**115**、`HEAD:frontend/src/lib/panel.ts`＝`16c2f038:frontend/src/lib/panel.ts`＝**311**（`diff` 二者＝**0 行**）；`HEAD:internal/panel/bridge.go`＝**167**、`16c2f038:internal/panel/bridge.go`＝**125**。⇒ 本文每一枚引用行号都小于所在件行数。
- **本腿落的件**（全在 `.scratch/wisp/probes/35/a3/`，⛔ 无 `*.out`，因根 `.gitignore:8` 是全仓 `*.out`）：
  `anchor.txt`·`raw-branch-frontend-src-files.txt`(85)·`raw-head-frontend-src-files.txt`(64)·`raw-head-go-jsontags.txt`(112)·`raw-keyreads-bothtrees.txt`·`raw-branch-propreads.txt`(27)·`raw-branch-grep-window.txt`(31)·`raw-branch-grep-snapshot.txt`(107)·`raw-branch-snapshotkeyreads.txt`·`head-panel.ts.txt`(311)·`head-App.tsx.txt`(117)·`branch-panel.ts.txt`(311)·`branch-App.tsx.txt`(531)·`branch-main.tsx.txt`(115)·`branch-panel-views.ts.txt`(106)·`branch-composer.tsx.txt`·`branch-git-branch.tsx.txt`·`branch-panel-sidebar.tsx.txt`·`head-bridge.go.txt`(167)·`head-composer.go.txt`(346)·`head-pump.go.txt`(712)·`head-approval.go.txt`(114)·`head-attachments.go.txt`(410)·`head-workspace.go.txt`(182)·`head-composer_dispatch.go.txt`(247)·`head-composer_handlers.go.txt`(165)·`head-panel_inbound.go.txt`(310)·`head-panel_host_windows.go.txt`(785)·`diff-panel.ts.txt`(0)·`field-reconciliation.md`（本文）
  ⚠ 一枚**假阴性标本就地认**：`raw-branch-grep-hostshape.txt`＋`.err` 都是 **0 行**——那次 `git grep` 我 `cd` 进了探针子目录，路径过滤变成 `.scratch/.../frontend/src`（不存在）⇒ **0 命中是尺放错位置，不是"页面无宿主引用"**；正确读数在**重跑件** `raw-branch-grep-window.txt`（31 枚命中）。这两枚空件按"临时件只建不删"**留着并随本件一起入库**，当作"过滤后 grep 空≠没跑"的标本；同样 0 行的 `diff-panel.ts.txt` 是**真读数**（两棵树的 `panel.ts` 字节相同），二者形状一样、含义相反，靠本段区分。
- **纪律自证**：全程只 `git show <ref>:<path>`／`git grep <pat> <ref> -- <path>`／`git ls-tree`／`git cat-file -e` 读分支；**零 `checkout`／`switch`／`worktree`／进那棵工作树**；`frontend/**`·`design/**` 零写入；禁读母本（`docs/PLAN.md`·`pending-and-issues.md`·`HANDOVER.md`·`*.css`·`.github/workflows/*.yml`）**一枚都没整读**（本腿一枚都没读）；只 commit、⛔ 未 push；每笔带显式 pathspec，⛔ `git add -A`／`.`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`。

---

## 8. 与派单／票面原文的冲突——以原文为准，具名报回

1. **票面 `:77` 那句"方法名册现量 6 枚（`internal/panel/bridge.go:42-49`）"的锚不准**（原文枚数对、行号区间错）：`:42-45` 只有 **4 枚**，另两枚在 `:66`/`:67`，`:46-65` 是注释。本文一律按 `:42-45`＋`:66-67` 引。⛔ 我没改票面一字。
2. **派单 §5.1 的候选起点有 5 枚根本不对接**（`main.tsx`·`agent-screen.tsx`·`approval-card.tsx`·`chat-composer.tsx`·`context-cards.tsx`·`records-table.tsx` 六枚里除 `App.tsx` 外全部零宿主读取，`insight-cards.tsx` 只有"在等票 167"的注释）——**照那张表找字段会白扫**。真起点清单见 §1.0。
3. **派单说"页面在另一棵树上、dev 上是另一份较旧的页面源码"——枚数与旧度都对，但有一处比转述更极端**：`frontend/src/lib/panel.ts` **两棵树字节相同**（`diff`＝0 行），所以"页面读什么"这一侧没有新旧两版，只有装配层（`App.tsx` 117↔531）与组件集（64↔85 枚）两版。这条影响 §3-b 的结论写法，我按原文事实写。
4. **票面 `:151` M1 与 `:77` 的枚数（6 枚）我复跑一致**；票面 `:87` 说 `panel_host_windows.go:403-404` 注释写 "four existing methods (bridge.go:42-45…)" 过期——本腿在 `:205` 另读到一句同样过期的话："…listeners (cmd/wisp/panel_inbound.go); the page's postMessage is just a second"，**它把页面的 `postMessage` 当成"第二条腿"来描述，而 §3-c 的量是"页面这条今天一跳都没到 Go"**。⛔ 注释我没动，交给票 33/248 的验收腿。
5. **派单 §1 要求"按字段语义变体扫"这条我撞了一次墙并就地认**：第一次扫入向方法名我用了 `panel\.(approval|mode|…)\.[a-z]+` 的正则形，命中 0（因 `cd` 进子目录，见 §7 的假阴性标本）；重跑改从仓根扫才拿到 31 枚。另外我为 `func jsonKeys` 猜过一次 helper 名（0 命中，真名未知，§6-b）。
