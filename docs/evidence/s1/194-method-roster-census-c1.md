# `194-c1` — 面板方法名册：只读代价普查（零产码）

- 程：`194-c1`（只读普查）｜派单：`.scratch/wisp/dispatches/2026-09-28-170x-readonly-194-c1-method-roster-census.md`
- 工单：`.scratch/wisp/issues/194-...-owner-ruled-align-the-code-to-the-spec.md`（**本表不勾它的任何 AC 框**）
- 起手时刻（本机 `date "+%Y-%m-%d %H:%M %z"` 现读）：`2026-09-28 16:12 +0800`
- 起手锚点：HEAD＝`8ca2291f`（`git log -1 --format='%h %s'` 现读，与派单标称的 `8ca2291f` **一致，未漂**）
  - ⚠ 记一笔时钟差：派单面写「派单时刻 `2026-09-28 16:5x`」，本程现读 `date` 是 **16:12 +0800**，比派单标称**早 38 分钟**。本表一律以**现读**为准；差因未定性（可能是派单时刻被前推书写），**不当缺陷、不据此推断任何事**。

## §1 起手脏件名册（写面闸门的起手基准）

尺：`git status --porcelain -- internal/ cmd/`（起手现跑）

```
（空——零枚具名脏件）
```

⇒ 起手名册＝**空集**。⚠ 按派单 §1「写面闸门（`A374`）」：闸门是**终态 == 起手名册（逐枚具名差集为空）**，**不是"必须为空"**。
同树在册飞着的写腿：`181-r2`、`185-r1`。终态若多出任何一枚，**具名登记为"不是我的"，不动、不提交、不还原**。
本程自己的写面只有三处（派单 §1）：本表 · `.scratch/wisp/probes/194/c1/**` · 票 194 的 Progress log 追加。

## §2 Q1 两份名册到底有几枚（逐枚现跑）

### 2.1 规格那份（`docs/specs/SPEC-08-ui-ball-panel.md`，表体逐枚具名）

尺：`Read` 现读 `SPEC-08:155-188`（对应派单给的 `sed -n '158,178p'` 窗口）。表体在 **`SPEC-08:161` 的表头**之下，行 **`163`–`174`**。

**逐枚方法名（把每行里 `/` 并列的两枚拆开各算一枚）**：

| # | 方法名 | 出处行 | 表内标的方向 | 需原生侧授权（逐字） |
|---|---|---|---|---|
| 1 | `panel.resync` | `SPEC-08:163` | **Go→前端推送** | `—` |
| 2 | `tasks.list` | `SPEC-08:164` | invoke | `—` |
| 3 | `task.detail` | `SPEC-08:164` | invoke | `—` |
| 4 | `history.query` | `SPEC-08:165` | invoke | `—` |
| 5 | `transcript.get` | `SPEC-08:165` | invoke | `—` |
| 6 | `approval.current` | `SPEC-08:166` | invoke | `—` |
| 7 | `approval.queue` | `SPEC-08:166` | invoke | `—` |
| 8 | `approval.decide` | `SPEC-08:167` | invoke | 「**「allow」拒绝一切面板来源（F2）；仅 `reject` 可面板发起**」 |
| 9 | `config.get` | `SPEC-08:168` | invoke | 「安全节放宽走 L2 重新确认（SPEC-03 §4.2）」 |
| 10 | `config.set` | `SPEC-08:168` | invoke | 同上一字 |
| 11 | `grants.list` | `SPEC-08:169` | invoke | `—` |
| 12 | `grants.revoke` | `SPEC-08:169` | invoke | `—` |
| 13 | `privacy.purge` | `SPEC-08:170` | invoke | 「purge 需确认」 |
| 14 | `privacy.export` | `SPEC-08:170` | invoke | 「purge 需确认」（同行，只管 purge） |
| 15 | `cost.summary` | `SPEC-08:171` | invoke | `—` |
| 16 | `models.list` | `SPEC-08:172` | invoke | `—` |
| 17 | `models.delete` | `SPEC-08:172` | invoke | `—` |
| 18 | `diagnostics.export` | `SPEC-08:173` | invoke | `—` |
| — | **事件推送那一行**：`task.delta`／`tool.chip`／`approval.request`／`ball.state`／`cost.tick`（5 枚，方向 Go→前端） | `SPEC-08:174` | Go→前端 | `—` |

**枚数（本程现数）**：
- 表体 `SPEC-08:163-173` 展开后＝**18 枚具名** ⇒ **票 194 面上写的「18 枚」是对的，不需更正**（这一支眼睛数赢了）。
- 加上 `SPEC-08:174` 事件推送那一行的 **5 枚** ⇒ 那张表**具名总面＝23 枚**。
- ⚠ **一处具名更正（不是枚数，是行号）**：票 194 AC#2 与本派单 §2 把 allow 侧禁令引作「`SPEC-08:169`」、把 `config.set` 的 L2 重确认引作「`SPEC-08:170`」。**磁盘现读**：allow 侧禁令在 **`SPEC-08:167`**，`config.set` 那一字在 **`SPEC-08:168`**（`169` 是 `grants.*`，`170` 是 `privacy.*` 的「purge 需确认」）。**两处引用各偏 2 行**；引文本身逐字对得上，只是行号错位。⇒ 影响：写腿若照票面行号去定位会改错行（把 `grants.revoke` 那行当 allow 侧禁令）。**本表只更正行号，不改 `docs/specs/**` 一字（AC#6）。**
- ⚠ 另一处口径要摊：`panel.resync`（第 1 枚）**在 invoke 白名单的表里，但方向标的是「Go→前端推送」**——它不属于「前端调后端」这一族。票 33 前例说的「有名无实现」正是它。**它算不算第 19 枚要补的"门"，是裁定面，本程只标不裁。**

### 2.2 代码那份（现量）

尺：`grep -nE '=\s*"panel\.' internal/panel/bridge.go` ⇒ **4 枚**；`grep -cE` 同尺 ⇒ **4**（与票 194 面一致；⚠ 未跑 `grep -c 'panel\.'`，那把尺会把 `bridge.go:35` 的注释 `"panel.*"` 计进来＝5，`33-r1`／`A380` 已具名报过，本程**不复量那把错尺**）。

| # | 常量名 | 字面值 | 出处行 |
|---|---|---|---|
| 1 | `MethodModeRequest` | `panel.mode.request` | `bridge.go:42` |
| 2 | `MethodWorkspaceRequest` | `panel.workspace.request` | `bridge.go:43` |
| 3 | `MethodAttachmentAdd` | `panel.attachment.add` | `bridge.go:44` |
| 4 | `MethodMessageSend` | `panel.message.send` | `bridge.go:45` |

`knownComposerMethod`（`bridge.go:104-110`）的 switch 就是这四枚常量，无第五枚。

### 2.3 交集

- 规格 18 枚（＋事件 5 枚）∩ 代码 4 枚＝**零枚**，与票 194 面「交集为零」一致（本程现读两侧名册，未见任何一枚同名）。
- ⚠ **反向也要摊**：代码那 4 枚 `panel.*` **在规格那张表里一枚都没有**。所以「代码向规格对齐」不等于「只往代码里加 18 枚」——**现有 4 枚是规格没写的既有门**，动它们（改名／换顺序）会撞 `composer_test.go:502` 那枚钉（票 194 AC#3 已具名警告，出处 `181-r1` 的按磁盘源码数枚数的尺）。本程判定：**4 枚保留，名册是"并集"问题不是"替换"问题** ⇒ 这是给编排者的一枚**待裁形状**（见 §6）。

## §3 Q2 逐枚四问主表（18 枚 invoke 侧）

尺（现跑，全量输出存 `.scratch/wisp/probes/194/c1/logs/q2-existence.txt`）：
`for m in …; do grep -rlF --include='*.go' "$m" internal/ cmd/ tools/ | grep -v _test.go; done` ⇒ 逐枚 PROD／TEST 两栏。

**现量总面**：23 枚具名里，**PROD 侧命中字面量的只有 3 枚**（`approval.queue`＝状态机事件名、`approval.decide`＝d22scan 的注释与拒绝文案、`grants.revoke`＝`queue.go:288/374` 真调用），
外加 `approval.request` 命中 `tools/d22scan/selftestsamples.go`（扫描器自测样本，不是实现）。⇒ **其余 19 枚字面量在产码与测试里都是零命中。**
所以本表的"有货"全部来自 **② 最接近的现有地基**（按能力找同物异名），没有一枚是靠同名找到的。

堆别：**A＝有货只差名字**｜**B＝有名字没货（要新建后端）**｜**C＝撞禁区／要先有逐枚批准**

| # | 方法 | 堆 | ① 同名等价实现（现量） | ② 落点／最接近的现有地基（`file:line`） | ③ 撞 AC#2 哪条禁区 | ④ 要不要新门控／审计 |
|---|---|---|---|---|---|---|
| 1 | `panel.resync` | B | 零命中（PROD/TEST 皆无） | 语义地基＝`NewSnapshot`「the one push from native state」`internal/panel/composer.go:71-78` ＋ 出向发布 `publishPanelSnapshot` `cmd/wisp/panel_pump.go:264`；**没有"重发一次"的请求面** | 不撞 | 不需要门控；需要一条"请求→重发"的入向路由（与 18 枚同属零生产听众） |
| 2 | `tasks.list` | A | 零命中 | `internal/memory/dao_tasklog.go:93 func (s *Store) ListTaskLogs`（数据面真在，缺面板路由） | 不撞（纯显示） | 只读，沿用既有审计即可 |
| 3 | `task.detail` | A/B 之间 | 零命中 | 最近地基 `internal/memory/dao_toolcall.go:111 ListToolCallsByTask(taskID)`；⚠ **按 id 取一条任务详情的函数没有**（只有全量列＋按任务列工具调用） | 不撞 | 需要按 id 的入参校验（id 是外部可控串） |
| 4 | `history.query` | A/B 之间 | 零命中 | `ListTaskLogs` `dao_tasklog.go:93`、`ListToolCalls` `dao_toolcall.go:121`、`ListArtifacts` `artifacts.go:96`；⚠ **没有任何查询面**（无过滤/分页/时间窗函数） | 不撞 | 查询面新造 ⇒ 需要边界（行数上限、字段脱敏），对照 `internal/observe/redact.go` |
| 5 | `transcript.get` | B | 零命中 | 面板侧只有内存态流片段 `Snapshot.Results []ResultChunk` `composer.go:59`（不可回查）；存储侧无 transcript 概念（`task_log`／`tool_call` 是最近的形） | 不撞 | ⚠ 逐字引 transcript 属于取证面：`SPEC-08:182` 那句「react-virtual（取证列表上千行）」说明它天然是大结果面 ⇒ 需要上限与脱敏，属新门控 |
| 6 | `approval.current` | A | 零命中 | `Queue.Depth()` `internal/agent/approval/queue.go:133` ＋ 面板卡视图 `Snapshot.Pending []ApprovalCardView` `composer.go:58` ⇒ **队头单显（C18）的形状已在快照里**，缺"单独问一次当前"的入向面 | 不撞（显示侧） | 不需要新门控；⚠ 不得由这一枚回带 decide 结果（AC#4②） |
| 7 | `approval.queue` | A | **PROD 命中＝同名噪声**：`internal/statemachine/events.go:43` 是**状态机事件名**（该行号出自票 194 面，本程**未复量**〔转述，未复量〕） | 深度／队列形状同上 `queue.go:133`；快照已带 `Pending` 全量 | 不撞（显示侧） | 同上；⚠ 与状态机事件**同名不同物**，直接补名册会让两处字面量撞车（判据尺会误认） |
| 8 | `approval.decide` | **C** | **PROD 命中＝d22scan 的注释与拒绝文案**（`tools/d22scan/main.go`、`selftestsamples.go`），非实现 | **两侧都有真实现**：allow 侧 `Queue.allow(corr, nonce)` `queue.go:343`（要 `grantNonce` 原生单用凭据 `queue.go:326`、`grants.spend` `:354`）；reject 侧 `Queue.reject(corr, reason)` `queue.go:393` | **撞①**：`SPEC-08:167` 逐字「**「allow」拒绝一切面板来源（F2）；仅 `reject` 可面板发起**」 | 已有闸门与判据在位（现量 `l2_grant_boundary_test.go` 7 枚：`TestAnsweredPanelRoutesCarryNoApprovalDecision:1229`／`TestRealGuardRefusesEveryAssemblableApprovalRouteName:1530`／`TestNoInboundEnvelopeCanBindAnApprovalVerdict:1547`／`TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes:1595`／`TestGrantWireShapesAreRefusedAtTheDoor:1959`／`TestPlantedGrantWiringGoesRedInASnapshot:2021` 等）⇒ **补这一枚＝只补 `reject` 一支；`allow` 一支不许接** |
| 9 | `config.get` | A/B 之间 | 零命中 | 装载面在（`internal/config/schema.go`、`validate.go`），⚠ **面向外部的 getter 没有**；现量 `internal/config/` 里唯一的 `Set/Get` 形导出是 `Manager.SetPermissionMode` `permmode.go:64` | 显示侧不撞②；⚠ **返回体必须能过滤安全节明文**（`D36` 三档生效级别）——未做 ⇒ 需要新门控 | 需要"读也分段"的口径（安全节里的密引用何面出）——本程未见现成面 |
| 10 | `config.set` | **C** | 零命中 | 同上；**唯一现成可借的写手就是 `Manager.SetPermissionMode` `internal/config/permmode.go:64`** | **撞③**：`SPEC-08:168` 逐字「**安全节放宽走 L2 重新确认（SPEC-03 §4.2）**」；⚠ **并且撞②**：把"档位"接进 `config.set` 就是把 R20／票 92 的权限输入口做成"面板直接改权限状态" | **要新门控**（安全节判定＋L2 重确认＋审计）。⚠ 本程**不裁**这一枚能不能做：今天没有任何一条"只放宽非安全节"的现成分节写面 |
| 11 | `grants.list` | A | 零命中 | `internal/memory/dao_misc.go:90 ListGrants(ctx)`、`:78 ListGrantsBySession(sessionID)` ⇒ **数据面完全在** | 显示侧不撞 | 只读；⚠ 授权记录含 nonce 派生物 ⇒ 出面前必须脱敏（`observe/redact.go`） |
| 12 | `grants.revoke` | **C?** | **PROD 命中＝真实现**：`internal/agent/approval/queue.go:288` `it.grants.revoke()`（dropLocked 内）、`:374`（`revokeGrants`） | 现成函数 `Queue.revokeGrants(corr string)` `queue.go:370`；⚠ 它**只按 correlationId 撤销、且不答复该条**（注释 `queue.go:366` "burns every live nonce … without answering it"） | **可能撞②**（R20／票 92：面板只许显示＋发起请求，不许直接改权限状态）。⚠ 反方证据也在此地：`queue.go:378-380` 逐字「**refusing is the fail-closed direction, which is why the panel may do it (SPEC-06 §9)**」⇒ 撤销属"收紧"一侧。**这一支算不算②＝裁定面，本程只标不裁** | 要：撤销必须落审计（谁按的、烧了哪几枚 nonce），今天 `revokeGrants` 无审计入参 |
| 13 | `privacy.purge` | A＋门控 | 零命中 | `internal/memory/privacy.go:170 PurgePrivacy(ctx, domain)`＋`:47 ListPrivacy`＋按域的 `PurgeArtifacts:300`／`PurgeMemories:155`／`PurgeProfiles:150`／`PurgeTaskLogs:129`／`PurgeToolCalls:146` | 不撞①②③；⚠ 规格自标「**purge 需确认**」（`SPEC-08:170`） | **要确认门＋审计**（不可逆删除）。⚠ **具名更正票面**：票 194「为什么值得做」段写「`privacy.purge`…今天连后端能力都还没有」——**存储层已实现**（上列 7 枚函数），缺的只是面板路由与确认门 |
| 14 | `privacy.export` | A | 零命中 | `internal/memory/privacy.go:198 ExportPrivacy(ctx, domain) ([]byte, error)` | 不撞三条 | 落盘/送出那一跳要新门控；⚠ 另撞 **AGENTS §1.2 的禁止形状**「把宿主内部 artifacts 写入实现成受门控的 Tool」⇒ 交付形态要先想清楚（不是 AC#2 那三条，但同样别顺手动） |
| 15 | `cost.summary` | A | 零命中 | `internal/memory/dao_misc.go:169 ListCostDaily(ctx)` ⇒ **日粒度数据面在** | 不撞 | 只读；⚠ 出向快照**没有 cost 段**（见 §4 `cost.tick`） |
| 16 | `models.list` | B | 零命中 | `internal/models` 现量导出方法名册：`NewManager downloader.go:97`、`Ensure:140`、`VerifyInstalled:228`、`Cancel:246`、`VerifyDir:567` ⇒ **没有面向外部的列举方法**（目录信息在 `ModelEntry`／Options 内，未成面） | 不撞 | 要新面；低风险（只读） |
| 17 | `models.delete` | B＋门控 | 零命中 | **没有删除已装模型的方法**（只有私有的 `removeStaging stagingDir:735/removeStaging:739`，作用域是暂存目录不是安装目录） | 不撞①②③；⚠ 但它是**删用户磁盘上的模型**⇒ 撞 AGENTS §1.2「不许在 `risk.PathResolver` 之外用 `filepath.Clean\|Abs` 做文件系统决策」 | **要新门控＋审计**（不可逆、外部落盘）。票面把它列为"没货"这一支**成立** |
| 18 | `diagnostics.export` | A | 零命中 | `internal/observe/diagnostics.go:62 BuildDiagnosticsBundle(o BundleOptions) (string, error)`＋`writeZipEntry:215` ⇒ **bundle 生成面在**，缺"递给面板/落盘"那一跳 | 不撞三条；⚠ 同 #14 撞 AGENTS §1.2 那条 artifacts 写入禁止形状 | 要：脱敏已在 `observe/redact.go`，交付路径要新门控 |

**逐枚一句话（三堆归堆结果）**：A 堆 **7 枚**（`tasks.list`／`approval.current`／`grants.list`／`privacy.purge`／`cost.summary`／`privacy.export`／`diagnostics.export`）；
A/B 之间（有数据面、没有查询或读取那一面）**4 枚**（`task.detail`／`history.query`／`config.get`／`approval.queue`）；
B 堆（后端能力真没有）**4 枚**（`panel.resync`／`transcript.get`／`models.list`／`models.delete`）；
C 堆（撞禁区／要先有批准）**3 枚**（`approval.decide`①、`config.set`③＋②、`grants.revoke` 待裁②）。
⚠ 合计 7＋4＋4＋3＝18，与 §2.1 枚数自洽。`privacy.purge` 同时属 A 与"要新门控"，`models.delete` 同时属 B 与"要新门控"——**堆别只回答"货在哪"，不回答"能不能直接接"**。

## §4 Q3 事件推送那一行（`SPEC-08:174` 的 5 枚，Go→前端）

**先钉住"管子"与"内容"各是什么形状（现量）**：
- 管子真在跑：`cmd/wisp/panel_pump.go`（283 行），文件头第 5 行逐字「**THIS IS WHERE THE PACKET STARTS EXISTING. Before these files, panel.Snapshot**」、第 13 行逐字「**no WebView2 host (ticket 33 is unclaimed), no postMessage writer, no local**」⇒ 出向** packet 生成与记账在**，最后一跳到不了前端（宿主未建）。发布函数 `publishPanelSnapshot` `panel_pump.go:264`，记账 `bookPanelSnapshot` `:182`。
- 内容载体只有一个结构：`panel.Snapshot{Pending, Results, Composer, GeneratedAt}` `internal/panel/composer.go:57-62` ⇒ **顶层 4 字段**，展开子结构含 `ApprovalCardView`／`ResultChunk`／`ComposerState`／`ModeView`／`WorkspaceView`／`AttachmentRef`。
- ⚠ **这 5 枚事件名一枚都不存在于产码**（`q2-existence.txt`：`task.delta`／`tool.chip`／`ball.state`／`cost.tick` PROD+TEST 全零；`approval.request` 只命中 `tools/d22scan/selftestsamples.go`＝扫描器自测样本，非实现）。⇒ 前端侧契约以何名字收本程**不判**（`frontend/**` 禁读面）。
- ⚠ 上游"十六字段表"（票 145）：本程**未复量**（不读前端），仅现量到 Go 侧 `go test` 自述的键数对齐结果：`Snapshot <-> PanelSnapshot: 4 JSON keys reconciled`、`ComposerState <-> ComposerState: 7`、`ModeView: 3`、`WorkspaceView: 6`、`AttachmentRef: 9`、`ResultChunk <-> ResultChunkView: 3`（出处 `composer_test.go:74/79` 运行输出）。**"十六"这一枚数以谁为准＝待裁**。

| 事件 | 有管子？ | **有内容？**（逐枚结论） | 证据（`file:line`） |
|---|---|---|---|
| `task.delta` | 有 | **有内容、没这个名字**：流式增量真在 `Snapshot.Results []ResultChunk{CorrelationID,Text,Done}` `composer.go:59/65-69`，但推送形态是**整份快照重发**，不是 delta 帧 | `composer.go:59`、`composer.go:64` 注释「ResultChunk is one streamed assistant chunk」 |
| `tool.chip` | 有 | **没内容**：`Snapshot` 无工具段，产码无 `tool.chip` 字面量；工具调用数据在存储侧（`ListToolCallsByTask` `dao_toolcall.go:111`）但**从未进过快照** | 反证：`composer.go:57-62` 四字段名册里没有 tool |
| `approval.request` | 有 | **有内容、没这个名字**：审批卡随快照出 `Snapshot.Pending []ApprovalCardView` `composer.go:58`；⚠ 字面量 `panel.approval.request` 现量 **PROD 零命中** | `composer.go:58`、`NewSnapshot(pending …)` `composer.go:76` |
| `ball.state` | 有（但**不在这根管上**） | **这条管没内容**：`Snapshot` 无 ball 段；球的态走**原生窗口另一条路** `Ball.SetState` `internal/ball/ball_windows.go:306`（＋`SetBadge:311`／`SetProgress:319`／`SetBadgeText:327`）⇒ 属"另有一根管"，不是"没货" | `ball_windows.go:306`、`anim.go:40 AnimationPolicy(s statemachine.State)` |
| `cost.tick` | 有 | **没内容**：`Snapshot` 无 cost 段、产码无 tick；数据面在存储侧 `ListCostDaily` `dao_misc.go:169`，**没有任何按 tick 推的函数** | `dao_misc.go:169`＋`composer.go:57-62` 反证 |

**⇒ 五枚里：2 枚"有内容有形状但换了名字"（`task.delta`／`approval.request`，实为整快照的两段）；1 枚"内容在另一条原生管上"（`ball.state`）；2 枚真没内容（`tool.chip`／`cost.tick`）。**
⚠ 别把这五枚补进 `bridge.go`：`bridge.go` 那把尺只数 `=\s*"panel\.`（入向常量），事件是 Go→前端方向，**放同一张常量表会让"入向白名单"和"出向事件名"混成一枚名册**——这正是票 194 要治的病的镜像形状。

## §5 Q4 三堆分批建议（**不拍板**）

⚠ **贯穿三堆的一条前提（本程现量，逐枚都适用）**：`ComposerDispatch` 的入向入口 `Handle` `internal/panel/composer_dispatch.go:120` 的生产调用者＝**零枚**——尺：`grep -rn "ComposerDispatch" internal/ cmd/ --include='*.go' | grep -v _test.go` 只命中它自己的定义与注释（`:90/:97/:120/:137/:167/:182/:192`）。
⇒ **本表 18 枚里没有任何一枚今天有真听众**；补名册补的是"门牌"，不是"能点的按钮"。**必须等票 33 的真窗口那一环（H2/H3）才有生产调用者＝全部 18 枚＋入向那 4 枚既有门，无一例外。**
⚠ 出向那一侧不对称，要分开说：出向 pump 已在 `cmd/wisp` 里被真调用（`bookPanelSnapshot`/`publishPanelSnapshot` 在 `panel_pump.go`，且 `gate-clauses.sh` 的 G6neg 腿在 `cmd/wisp/run.go:568` 抓到 CloseTask 钩子），所以**出向 5 枚是"有生产作者、缺最后一跳"；入向 18 枚是"连作者都没有"**。

**堆 1｜有货只差名字（7 枚＋A/B 4 枚里数据面已在的 `grants.list`/`cost.summary` 类）**
- 清单：`tasks.list`、`approval.current`、`grants.list`、`privacy.export`、`cost.summary`、`diagnostics.export`、`approval.queue`（⚠ 与状态机事件同名不同物，见 §3#7）、`config.get`（只读非安全节那一半）。
- 最小写腿写面范围：`internal/panel/bridge.go`（**只追加常量，不动 42-45 那四枚的顺序与命名**，票 194 AC#3）、`internal/panel/composer_dispatch.go`（dispatch switch `:137` ＋ `rosterMismatch :167` 的名册一致性）、**新建** handler 文件（例 `internal/panel/roster_handlers.go`，⚠ 不碰 `composer_handlers.go` 既有四枚）、**新建**判据文件（⚠ AC#6 点名 `tokens_fourway_test.go`／`composer_test.go`／`l2_grant_boundary_test.go` 三枚不许动 ⇒ 判据只能进新文件）。
- 每枚自带 AC#4 三判据的可行性：**已具备**——白名单外拒答＋审计在 `ParseComposerRequest` `bridge.go:90-92`（`knownComposerMethod` 失败即拒）与 `rosterMismatch` `composer_dispatch.go:167`；归因字段用 `RequestID`（`bridge.go:97-100` 已强制非空）。

**堆 2｜有名字没货（要新建后端能力，4 枚）**
- 清单：`transcript.get`（存储侧无 transcript 概念）、`models.list`（models 包无列举面）、`models.delete`（**无删已装模型的方法**＋要 `risk.PathResolver` 规矩）、`panel.resync`（无"请求重发"入向面）。
- 最小写腿写面范围：`internal/models/*.go`（新增导出面）＋`internal/memory/*.go`（transcript 概念或明说复用 `task_log`／`tool_call`）＋`internal/panel/`（新常量＋新 handler）；`models.delete` 另需 `internal/risk/` 侧的路径决议入口（⚠ 不许在 `PathResolver` 外做文件系统决策，AGENTS §1.2）。
- ⚠ 这一堆的每一枚都是**功能票不是改名票**：派写腿前要先有"这枚功能到底做不做、落在哪个切片"的批准，否则造出来的就是第二枚 `panel.resync`。

**堆 3｜撞禁区要先裁（3 枚）**
- 清单：`approval.decide`（①：allow 侧永拒面板来源，`SPEC-08:167` 逐字；reject 侧可）、`config.set`（③：`SPEC-08:168` 逐字安全节 L2 重确认；**并且**其唯一现成写手 `Manager.SetPermissionMode` `permmode.go:64` 天生撞②）、`grants.revoke`（②待裁：撤销属 fail-closed 收紧一侧，`queue.go:378-380` 有反方逐字证据）。
- 最小写腿写面范围：**本堆今天不该派**——`approval.decide` 若只补 reject 一支要新增"面板来源只能带 reject"的门（可借 `Source` 判定 `bridge.go:93`＋`l2_grant_boundary_test.go` 的 7 枚形状）；`config.set` 要先有"分节写＋安全节判定＋L2 重确认"的**地基**（今天不存在，本程现量）。⇒ **这三枚逐枚都需要 `A##` 人工批准（AC#3），不是写腿能自决的。**

## §6 Q5 债务清点：`C24 GojaHostAPI` 初始集与 `C17` 方法白名单定稿到哪一步

**尺现跑**：`grep -rn "C24" docs/PLAN.md docs/specs/`。现量 10 处，关键的三处：

- `docs/PLAN.md:1374` 逐字：**「| **C24** | **`GojaHostAPI`** | Tier-2 插件在 JS 侧可见的**唯一**宿主入口全集（方法名 + 对应 capability + RiskLevel）……C16 的 `host_api` semver 就是本契约的 semver | D39（补 C16 的实体） |」** ⇒ 契约**定义了"要有一张全集表"这件事**，⚠ **初始集本身没有枚举**（本程在 PLAN.md 里未见任何一份 C24 方法清单）。
- `docs/specs/SPEC-12-roadmap-governance.md:48` 逐字：**「C24 GojaHostAPI 初始集与 C17 方法白名单定稿（S7/S5 切片卡批准）」** ⇒ **这就是 `AGENTS §2` 那条在册待定项的原文，今天仍原样在册、未销**。
- `docs/specs/SPEC-07-tools-and-plugins.md:119` 与 `:147`：要求 JS 侧唯一宿主入口＝C24，且「C24 契约测试：JS 侧调用白名单外宿主入口 → 必须拒绝；**扫描 goja 绑定表与 C24 文档一致**」⇒ 定稿的验收形状已写、无实现体可对照（`internal/plugin/` 侧本程未查，⚠ 未测项已记 §7）。

**⚠ 本票直接相关的另一处口径（这条比 C24 更要命）**：那张 18 枚表**自己的标题就带着未定稿标记**——`SPEC-08:156` 逐字：**「### 5.2 C17 PanelBridge 方法白名单【SPEC 提案，S5 定稿走契约批准】」**，`SPEC-08:159` 逐字「未列出方法名 → 拒绝并记日志」。
⇒ 「代码向规格对齐」的规格那一侧**是一枚 S5 待批的提案，不是已定稿的契约**；而 `C17` 的权威文字在 `docs/PLAN.md:1367`（**该行号出自票 194 AC#6，本程未复量〔转述，未复量〕**）。**本程只报现状，不替它定稿**（派单 Q5 逐字要求）。

## §7 本程没测／判错的（逐名）

1. **判错并自纠一次**：第 13 枚调用那次「domain 点号字面量普查」（`grep -rhoE '"(task|tool|approval|…)[a-zA-Z._-]*"'`）**没有排除 `_test.go`**，所以它给出的计数（如 `panel.approval.request`＝12、`history_changed`＝5、`panel.review.allow`＝5、`panel.mode.set`＝3）**含测试与扫描器自测样本，不可当作产码事实**。⇒ 本表一律以第 11 枚那次 PROD/TEST 分栏的 `q2-existence.txt` 为准；这次失误的直接后果我已把它转成一条**正面证据**（`panel.approval.request` 实为 PROD 零命中）。
2. **未复量（照抄了他处文字，已就地标注）**：① `approval.queue` 的"状态机事件名"落点 `internal/statemachine/events.go:43`；② `C17` 权威文字行号 `PLAN.md:1367`。两处均标〔转述，未复量〕。
3. **没测**：`internal/plugin/` 里有没有 goja 绑定表（Q5 第三句要求"与 C24 文档一致"，本程只证明文档侧没枚举初始集，未证明实现侧形状）；`internal/perm/`、`internal/session/`、`internal/observe/` 三包的导出名册（我只按概念找函数，未全包普查）；票 145 那张"十六字段表"的真实枚数（前端面，禁读）。
4. **没跑**：`go test ./internal/config/ ./internal/risk/ ./internal/tools/`（票 194 AC#7 要求逐包四连；本程是只读普查、派单 §3 只点名 `./internal/panel/`）⇒ 另三包的今天红态**本程不背书**。全仓面我只经由 `scripts/d22scan.sh` 自带的 `runtests.sh` 见到一行 `packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0`，⚠ **它与本程随后现量的 3 枚 panel 红互相矛盾**（同树在飞件所致，本程不裁谁对、不改尺）。
5. **尺的边界**：§3 的"有货"判定靠**能力近邻**而非同名，属我的解读；若编排者要求"只认同名实现"，则 23 枚里 19 枚应直接归"没货"。两种口径的差我已逐枚写出（PROD 命中列＋地基列分两栏）。

## §8 门禁终态（只读程取数，`A363`：不充当结案凭据）

（待填）

## §9 被拒调用＋零删除自证＋工具调用终值

（待填）

## §10 next＝派写腿之前还缺哪几枚批准

（待填）
