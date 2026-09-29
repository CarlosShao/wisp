# 票 220 只读普查 c1 — 「要修得动它，落点到底在哪」

- 锚点 `git rev-parse --short HEAD` = **6233dedc**（分支 dev，现跑；采集起点 11:28）。
- 复核锚点 = 采集结束时 HEAD 再前移（`git diff --name-status 6233dedc..HEAD` 现跑：动过/新增的 Go 件 = `cmd/wisp/approval_always.go(A) / approval_reply.go(M) / run.go(M) / approval_always_201_test.go(A) / approval_seam_201_test.go(A)`、`internal/agent/approval/replies.go(A)+测试(A)`、`internal/config/allowdirs.go(A)`。**gate.go／queue.go／pending_read.go／pump.go／subagent_roster_197.go／task.go／subagent_197.go／bridge.go／loop.go／l2_grant_boundary_test.go 在该窗口零改动**，本文对它们的行号在新锚点同样成立；`run.go` 相关引用已于 11:46 在新 HEAD 重跑校正，其余 `.Cancel(` 命中行全部原位复认，票 201 新增件里对 `TaskRoster.Cancel` 的调用 = 0 枚，E1 结论不因漂移失效）。
- 取证时间（现跑 `date "+%Y-%m-%d %H:%M %z"`）= **2026-09-29 11:38 +0800**；校正重跑时间 = 同日 **11:46 +0800**（两枚都是当时现跑 date，非推算）。
- 性质：**只读普查**。零编译零测试零写码；`frontend/**`／`design/**` 未读、结论不转述其内容。
- 本文所有行号均为本轮现跑尺所得（命令见下），**没有一枚抄自 `.scratch/wisp/probes/**` 归档**。

## 本轮跑过的尺（可原样重跑，均在仓库根目录）

```
git rev-parse --short HEAD ; git rev-parse --abbrev-ref HEAD
grep -rn "LiveApprovals" --include="*.go" . | grep -v ".scratch"
grep -n "windows" internal/agent/approval/gate.go
grep -n "^func (g \*Gate)" internal/agent/approval/gate.go
grep -rn "g\.windows\|\.windows\b" --include="*.go" . | grep -v ".scratch"
grep -rn "TaskRosterSectionFrom\|ApprovalCardView" --include="*.go" . | grep -v ".scratch"
sed -n '183,196p' internal/agent/approval/ui.go ; sed -n '27,39p' internal/agent/approval/ui.go
sed -n '505,540p' cmd/wisp/run.go ; sed -n '58,76p;130,204p' cmd/wisp/panel_pump.go
sed -n '640,655p' internal/agent/loop.go ; sed -n '1123,1136p' internal/agent/loop.go
grep -rn "liveVerdicts\|Verdicts:" --include="*.go" cmd/wisp internal/panel | grep -v _test
grep -rn "\.Cancel(" --include="*.go" . | grep -v _test | grep -v ".scratch"
grep -rn "AttachCancel\|DetachCancel" --include="*.go" . | grep -v ".scratch"
grep -n "clash\|re-stamp" --include="*_test.go" internal/agent/approval internal/panel cmd/wisp
grep -rn "BlockedOnApproval\|blockedOnApproval" --include="*.go" . | grep -v ".scratch"
grep -n "func Test" internal/agent/approval/pending_read_test.go internal/agent/approval/queue_test.go \
  internal/agent/approval/ticket87_veto_l2_test.go internal/agent/approval/ticket97_alias_direction_test.go \
  internal/panel/subagent_roster_197_test.go cmd/wisp/subagent_blocked_197_test.go
grep -n "grantFieldWords\|grantRouteWords\|grantRoutePrefixes\|grantRouteSuffixes\|inboundTypeRegistry" \
  internal/panel/l2_grant_boundary_test.go
sed -n '183,250p;1034,1062p;1093,1100p;1310,1370p;1590,1670p;1715,1727p;1879,1915p' internal/panel/l2_grant_boundary_test.go
grep -n "task.cancel\|task.list\|task.output\|task.spawn\|C17" docs/PLAN.md | head -45
grep -rn "DEFERRED(" --include="*.go" internal/tools | grep -v _test
grep -n "func BuiltinTaskEntries" -A3 internal/tools/task.go ; sed -n '14,30p;425,475p;493,512p;585,597p' internal/tools/task.go
sed -n '30,45p;183,196p' internal/tools/subagent_197.go
grep -n "WithValue" internal/tools/bridge.go
```

## A. 两条「正在等的东西」各自的容器

### A1 L2 队列容器 `Queue.pending` 〔已证（现读）〕

- 类型：`pending []*qitem` —— `internal/agent/approval/queue.go:69`（`Queue` 结构体 `:62-79`）；配套索引 `byID map[string]*qitem` `:70`、`alias map[string]map[*qitem]bool`（拒方向别名索引）`:71-76`。`qitem` 结构体 `:31-56`：`Dec tools.Decision / Corr string / Seq uint64 / names []string / bind / grants / state itemState / answer chan answer / replayOf`。
- **谁往里放**：`(*Queue).push` `queue.go:141-169` —— `:163` append 进 `pending`，`:164` 写 `byID`，`:165` 经 `indexLocked`（`:195-204`）登记 `names` 进 `alias`。生产里 push 的唯一调用者 = `Gate.PendingApproval` `internal/agent/approval/gate.go:473`（L2 路径，函数 `:468`）；R7 批量升级支路也走同一函数（`gate.go:236` 调 `PendingApproval`）。
- **谁从里取**：`LiveApprovals`（`pending_read.go:104-122`，只遍历 `q.pending`，`statePending` 过滤 `:112`）、`dropLocked`（`queue.go:279-295`，`:286` 删 byID）、`deliver` `:298-311`、`head` `:494-501`（队头单显）、`Depth` `:133-137`、`view` `:472-480`、`position` `:268-275`。

### A2 L1 窗口容器 `Gate.windows` 〔已证（现读）〕

- 类型：`windows map[string]*window` —— `gate.go:64`（`Gate` 结构体 `:54-69`）。值类型 `window` 结构体 `gate.go:72-77`：**只有 `ctx / corr / tool / vetoes chan Veto` 四枚字段，不携带裁决**（没有 Decision/Reason/RulesHit/Level —— 这对 D 节是硬成本）。
- **谁往里放**：`openWindow` `gate.go:318-326`（`:324` 写 map；`:321` 同 corr 撞名即拒），调用者 = `PendingWindow` `gate.go:249`（L1 路径，函数 `:218`）。键：`corr := orDefaultText(d.CorrelationID, d.TaskID)` `gate.go:247` —— **L1 的 corr 缺省回落到 taskID，就这一处**。
- **谁从里取**：登记删除 = `closeWindow` `:328-332`（`PendingWindow` 的 defer，`:252`）；**全仓唯一一次读出按键查 = `Veto` 里 `w := g.windows[v.CorrelationID]` `:393`**（函数 `:387`）。本轮 `grep -rn 'g\.windows\|\.windows\b'`（非测试、去归档）命中仅 `gate.go:321/324/330/393` 四行自己用。

### A3 `LiveApprovals()` 本体与生产调用者 〔已证（现读）〕

- 实现：`internal/agent/approval/pending_read.go:104-122`。返回 `[]LiveApproval`；`LiveApproval` 结构体 `:46-50` = `{CorrelationID string; Decision tools.Decision; Position int}`，`Decision` 经 `cloneDecision`（`:57-68`，票 146 的深拷钉）逐槽拷贝；`:116` **`CorrelationID: it.Corr`** —— 交出去的是队列重贴后的键，`it.names` 别名不外传。
- **生产里谁在调它（具名）**：`(*agentRuntime).liveVerdicts` —— `cmd/wisp/panel_pump.go:58-76`，`:62` `items := rt.gate.Queue().LiveApprovals()`；这枚 reader 在装配根 `cmd/wisp/run.go:519` 被接成 `PumpSources.Verdicts`（`Verdicts: rt.liveVerdicts`，泵装配 `run.go:518-540`）。消费链：`SnapshotPump.Snapshot` `internal/panel/pump.go:195-202`（`:196-197` 读 Verdicts，`:199-202` 逐枚 `NativeVerdict.CardView()` → `ApprovalCardView`，`pump.go:88`、`internal/panel/approval.go:72-100`），`:263` `snap.Tasks = TaskRosterSectionFrom(p.src.Tasks(), cards)`。**不是「只有测试在调」**：run.go:519 是 `wisp run` 生产装配位。

### A4 L1 窗口的导出面 〔已证（现读）：**没有**〕

- `Gate` 全部导出方法（`grep -n '^func (g \*Gate)'` 现读）：`Queue :135 / Channels :139 / Window :142（返回 time.Duration，是配置长度不是窗口清单）/ AdmitTextTask :160 / PendingWindow :218 / Complete :355 / Veto :387 / LateVeto :448 / PendingApproval :468 / Native :588 / Panel :592 / DecideFromNative :626 / DecideFromPanel :638 / Replay :657`。**没有一枚枚举 `windows`。**
- 两侧 API 面也不含枚举：`NativeAPI` = `Allow/Reject`（`ui.go:143-151`），`PanelAPI` = `Reject/Head/View`（`ui.go:155-160`），全部按 corr 名字单点寻址、且只作用于 L2 队列。
- 宿主侧唯一的 L1 痕迹 = `consoleApprovalUI` 的**计数器**：`cmd/wisp/run.go:1085-1088`（`shown()` 返回 `u.cards`，`:1093-1095` `Prompt` 里 `u.cards++`；`Update` 于 11:46 在新 HEAD 重读为 `:1142`），`windowCount` `run.go:660-665`。**是 int，不是名单** —— 连「哪枚任务在窗口里」都不记得。
- 结论：**L1 短窗口今天在 Go 侧没有任何可枚举出口；票 220 现量表「⛔ L1 短窗口今天不可枚举」这一行复核为真。**

## B. 名册那枚 join 的确切形状

> 路径勘正：任务书写的是 `internal/tools/subagent_roster_197.go` —— 现树**不存在**该文件（`Glob **/subagent_roster_197.go` 现读只命中 `internal\panel\subagent_roster_197.go` 一枚，另有一枚归档快照不算）。tools 侧同名相近的是 `internal/tools/subagent_197.go`（spawn 工具）。以下全部指 panel 那枚。

### B1 `TaskRowView` 逐字段（分母 = **12 枚**）〔已证（现读）〕

`internal/panel/subagent_roster_197.go:101-133`：

| # | 字段 | JSON 键 | 类型 | 行 |
|---|---|---|---|---|
| 1 | TaskID | `taskId` | string | :102 |
| 2 | Label | `label` | string | :103 |
| 3 | Kind | `kind` | string | :104 |
| 4 | ParentTaskID | `parentTaskId` | string | :105 |
| 5 | Status | `status` | string | :109 |
| 6 | StatusKnown | `statusKnown` | bool | :110 |
| 7 | StatusReason | `statusReason,omitempty` | string | :115 |
| 8 | StreamKey | `streamKey` | string | :121 |
| 9 | **BlockedOnApproval** | **`blockedOnApproval`** | **bool** | **:127** |
| 10 | StreamTruncated | `streamTruncated` | bool | :130 |
| 11 | StreamElidedRunes | `streamElidedRunes` | int | :131 |
| 12 | StreamDropped | `streamDropped` | bool | :132 |

（外层 `TaskRosterSection` 另带 6 枚键：`rows/inFlightSlots/poolCap/streamTruncated/streamElidedRunes/droppedStreamKeys`，`:145-164`。）

### B2 `blockedOnApproval` 是哪一行、拿什么当键 〔已证（现读）〕

- 建索引：`:190-195` —— `waiting[card.CorrelationID] = true`（**键 = 卡的 correlation id**，`:192` 空串跳过）。
- 查：`:210` —— `BlockedOnApproval: row.TaskID != "" && waiting[row.TaskID]`（**查 = 行的 task id**）。
- ⇒ 票 220 现量表写的 `:190-210` 区间，本轮精确化为「建 :190-195 / 查 :210」。

### B3 两半键的来源与拼法 〔已证（现读）〕

- **`card.CorrelationID` 正向链**：`internal/agent/loop.go:647`（`req := ToolRequest{TaskID: taskID, CorrelationID: taskID, …}` —— 生产派发时 corr 就是 taskID）→ `internal/tools/bridge.go:246-247`（空 corr 回落 `req.TaskID`）→ Decision 上 corr 两处置法：`bridge.go:330`（`CorrelationID: req.CorrelationID`）与 `bridge.go:961`（`orDefault(req.CorrelationID, req.TaskID)`）→ `Gate.PendingApproval gate.go:468/473` → `Queue.push queue.go:141-169`（**可能在此重贴 `id#seq`，见 C 节**）→ `pending_read.go:116` `CorrelationID: it.Corr` → `cmd/wisp/panel_pump.go:66` → `pump.go:88 CardView()` → `approval.go:40/77` `ApprovalCardView.CorrelationID`。
- **`row.TaskID` 来源**：`cmd/wisp/panel_pump.go:146-204`（`taskRosterState`）—— `:174 TaskID: rec.TaskID`，rec 来自名册 `rt.tasks.Look(id)` `:188` 与 `rt.tasks.Descendants(id)` `:191`；task id 本身由 `internal/agent/loop.go:1123 newTaskID()` 铸（UUID 形 8-4-4-4-12 十六进制，随机 16 字节）。行的登记时机在子代理出生路径（`tools.PublishSubagent`，注释见 `cmd/wisp/subagent_blocked_197_test.go:103-105`）。
- 一卡对得上的原因：生产配对 `CorrelationID == TaskID`（`loop.go:647`；roster 文件自己的注 `:172-175` 就写着这条并声明「不同配法的宿主会读成不阻塞」）。**一旦 push 重贴，这枚前提当场失效。**

## C. queue.go 重贴 `id#seq` 那一段

### C1 现读与机理 〔已证（现读）〕

- 确切行：`internal/agent/approval/queue.go:152-154`：
  `if _, clash := q.byID[corr]; clash { corr = fmt.Sprintf("%s#%d", d.CorrelationID, q.seq) }`
  （前置：`:147 seq++`、`:148-151` 空 corr 生成 `approval-<seq>` —— 生成名互异永不撞，所以撞名只会发生在**调用方自带同名 corr** 的场合，生产上即 `corr==taskID` 的同一任务第二张卡。）
- 为什么要重贴：`byID` 是一枚 `map[string]*qitem`（`:70`，`:164` 直接下键写），同 corr 的第二张卡不改名就会**覆盖第一张的索引**。重贴保住了「每张活卡各占一枚键」，同时 `:155-161` 把原始名收进 `qitem.names`（`:36-42` 注释逐字：票 87 —— 在屏的卡要能被显示侧用原名指到）。
- 第二张卡落地时第一张卡的键怎么变：**第一张保持裸名 `taskID` 不变**；第二张变 `taskID#<seq>`。坏形发生在**第一张先被答复**：`deliver :298-311` → `dropLocked :279-295` → `:286 delete(q.byID, it.Corr)` —— 裸名从 byID/LiveApprovals 消失，**只剩 `taskID#seq` 一枚 pending**。此时名册查 `waiting[row.TaskID]` 落空 ⇒ 那一格读 false ⇒「越忙的任务越读起来像没在等」。
- 缓解为什么救不到读数：`names` 别名只服务**拒方向**（`lookupForRefusalLocked :250-265`，`otherNames :174-192` 把 `d.TaskID` 收进别名，`reject` `:403-431` 可达）；**allow 方向**只认精确键（`lookupForAllowLocked :231-236`）。而读方向 `LiveApproval`（`pending_read.go:46-50`）只带 `it.Corr`，`it.names` **没有任何出口**。

### C2 今天有没有尺看得见这一形 〔已证（现读）：**没有，一枚都没有**〕

- 队列侧（`internal/agent/approval`）：`pending_read_test.go` —— `TestLiveApprovalsReportsPendingRowsAndDoesNotConsumeThem :111`（3 枚**互异** corr）、`TestLiveApprovalsSkipsRowsThatAreNotPending :190`、`TestLiveApprovalsOnNilQueueIsNil :224`；`ticket146_liveapprovals_backing_test.go` —— `TestLiveApprovalsSharesNoReferenceSlotWithTheQueue :208`、`TestLiveApprovalsInPlaceWriteCannotReachTheQueueRecord :264`（拷问的是**底层数组共享**，不是键的形状）；`queue_test.go` —— `TestL2NativeAllowExecutesThroughTheBridge :22 / TestPanelSourcedAllowIsRejectedOnEveryForgeableAxis :89 / TestGrantIsSingleUseAndBoundToItsItem :150 / TestL2QueueAutoRejectsAt300sKeepsTaskAliveAndWarnsAt270s :189 / TestReplayRedisplaysUnderAFreshGrant :269 / TestHostUnreachableAndFullQueueFailClosed :318`。**没有任何一条把两枚同 corr 的卡真 push 过。**
- 重贴形状自身的钉只有拒方向：`ticket87_veto_l2_test.go:129 TestAVetoNamedByTheHostsOwnKeyStillRefusesThatCard`（注释 `:127` 明写「when it has re-stamped one」—— 钉的是 veto 还能落到重贴后的卡）；`ticket97_alias_direction_test.go:75 TestAnAliasCanNeverBuyAnAllow / :149 TestEveryRefusalRouteOnAnUnknownEntryStillRefuses`。⇒ 不变式「重贴后仍可按 taskID 落到拒」是**有**的，但它对读方向（名册那格）视而不见。
- 名册侧：`internal/panel/subagent_roster_197_test.go:103 TestTheRosterReaderPutsSubagentsOnTheWire`（join 断言 `:164/:169`，卡是该文件手搓的**单卡**；该文件 `:466` 自陈「join 是对着本文件造的卡断言的」）；`cmd/wisp/subagent_blocked_197_test.go:59 TestRunPacketMarksTheRosterRowACardIsHolding` —— 走真队列真 packet，但判据写死 `len(snap.Pending)==1`（`:155-166`）且只开**一张**卡。**「同任务第二张卡」形状在两层都不存在。**
- L1→名册方向：无尺可能看见（A4：根本没有枚举口）。

## D. 修法候选（≥3 支）与 l2 冻结钉的碰撞判定

### D0 先钉死尺子本身：`internal/panel/l2_grant_boundary_test.go` 现读 〔已证（现读）〕

- 扫描射程：`scanGrantBoundary(t, dir)` 定义 `:1163`，两处调用传入的都是 **`<repo root>/internal/panel` 目录**（`:1230-1231`、`:1573-1574`）。**approval / cmd/wisp / tools 包里的 Go 导出方法名不在这枚钉的扫描面内。**
- 词根表一（JSON 键面）`grantFieldWords` = **20 枚**，`:183-188`：approve / approval / autoapprove / grant / granted / allow / allowed / allowonce / permit / permitted / ratify / authorize / authorised / decision / decide / verdict / outcome / bypass / override。
- 词根表二（路由名面）`grantRouteWords` = **11 枚**，`:192-195`：approve / approval / grant / allow / permit / ratify / authorize / authorised / decide / decision / verdict。
- 拼法收敛：`carriesGrantWord` `:242-249` —— 转小写并剥 `_ - 空格 .` 后做**子串匹配**，所以大小写／下划线／连字符／点分／驼峰全部折叠成同一名。路由名网格 = 8 前缀（`:1316-1319`：panel、panel.review、panel.approval、panel.l2、panel.mode、panel.workspace、panel.attachment、panel.message）× 11 词根 × 3 后缀（`:1324`：`"" / ".request" / ".now"`）；枚数由常量钉死 `:1333-1337`（wantGrantRoutePrefixes=8 / Words=11 / Suffixes=3，注释 `:199-200` 记「每枚词根贡献 48 枚候选名」）。路由形状过滤器 `routeShapedName` `:1049` 起（≥2 段点分）。
- 豁免规则：JSON 键面只对 `inboundTypeRegistry` 里 **4 枚入向类型**成立（`ComposerRequest/ModeRequest/AttachmentPayload/AttachmentRef`，`:1720-1727`）—— 出向类型（Snapshot/ApprovalCardView/TaskRowView）不在键扫内；头注 `:174-181` 逐字把 confirm/accept/intent 这类「命名审批邻居但不携带裁决」的词**刻意留在词表外**。
- `reason` 为什么安全：`reason` **不在两枚词表里的任何一枚**（现读 `:183-195`）；`:1602-1608` 种的诱饵 `plantedVerdictIn` 里 `Reason string json:"reason"` 本身不响，响的是内嵌 `Decision{Verdict}` 那枚 "verdict"。所以 `ApprovalCardView.Reason/ReasonKnown`（approval.go `:45-50`）与 `reject(corr, reason)`（queue.go `:403`）都照常存在。
- ⛔ 该钉一字不许动 —— 以下每支只判定「撞不撞」，不提议改钉。
- 既有事实参照：`blockedOnApproval` 键名含 "approval" 词根，但 TaskRowView 是**出向**类型、不在 4 枚入向 registry 里，今天不被扫（B1 表 #9 就是它）。**新增字段落进那 4 枚类型或新增 `panel.*` 路由字面量才会红。**

### D1 候选一（= 票 220 AC#1 形）：join 改「corr 与 taskID 双查 + `#序号` 后缀匹配」〔已证可达（现读）〕

- 动哪些文件：`internal/panel/subagent_roster_197.go:190-195`（建索引时同时按裸名收）与 `:210`（查改成 `waiting[row.TaskID] || ∃k HasPrefix(k, row.TaskID+"#")`）；该文件现只 import `sort`（`:39`），需加 `strings`。正控测试补两枚：纯 join 层进 `subagent_roster_197_test.go`，生产形状进 `cmd/wisp/subagent_blocked_197_test.go`（同文件现读判据 `:155-166` 的等长 1 卡判据要放宽成 ≥1 —— 那是测试文件不是冻结件）。
- 撞钉判定：**零碰撞** —— 不新增方法名、不动 4 枚入向类型、不写 `panel.*` 字面量。
- 判据在生产里可达：是。链路 run.go:519 → panel_pump.go:62 → pump.go:196-202/:263 全生产件；种形用真 `Gate.PendingApproval`（subagent_blocked_197_test.go:120-126 已有同款手法）连开两卡、答复第一张 ⇒ 改前读 false、改后 true。
- 代价：最小。误匹配面：taskID 为 UUID 形（loop.go:1123-1136 现读），一枚 taskID 恰是另一枚 `taskID#seq` 前缀的概率可忽略。
- 边界：对 L1 那半**零改善**（A4 —— 根本没有东西进 waiting）。

### D2 候选二（= 票 220 AC#2 甲）：给 `Gate` 补只读 L1 窗口枚举器〔建了但没接的对岸：容器有、面无 —— 现读为「面不存在」〕

- 动哪些文件：`internal/agent/approval/gate.go:72-77`（`window` 结构体得**增储裁决** —— 现读四字段没有 Reason/RulesHit/Level，而 `NativeVerdict→CardView`（pump.go:88 → approval.go:72）要这些才能成卡）、`:318-326`（openWindow 扩参）、`:248-252`（PendingWindow 构造点传入 d）、新增一枚导出方法；再接 `cmd/wisp/panel_pump.go:58-76`（liveVerdicts 并读枚举器，L1 行 Level="L1"）。
- 撞钉判定：方法落在 approval 包 —— **不在钉的扫描面**（D0：只扫 internal/panel 目录）。把 L1 变成 `snap.pending` 的**多出行**不改任何 JSON 键 ⇒ `TestApprovalCardViewJSONKeysMatchFrontendTypes`（approval_test.go:105-132，ApprovalCardView/ResultChunk/Snapshot 三对双向键核）不红。**除非**新枚举器名字被写成 internal/panel 里的点分字面量 —— 不要那样做。
- 命名提示：`Window` 已被占（gate.go:142 返回时长），枚举报需另名；「不新增方法名」是票 197 禁区、票 220 AC#2 甲自带**先落台账批准记录**的前提。
- 判据生产可达：是 —— 种一枚停在 L1 的任务（`PendingWindow` armed-before-prompt，gate.go:254-257；3s 窗口 queue.go:115-122），改前名册那格恒 false（枚举器不存在，见 A4），改后为真。
- 代价：中偏大 —— ①window 扩储触及写腿正活跃的 `gate.go`；②窗口开／关是否已触发泵 publish **本轮未证**（见末节），可能还要补发布触发点。

### D3 候选三（= AC#2 乙的落地形）：不动 approval 包，在装配根把 L1 记进宿主 UI 适配器〔已证（现读）：原料全在，记录不存在〕

- 动哪些文件：`cmd/wisp/run.go:1085-1095`（`consoleApprovalUI.Prompt` 今天只 `u.cards++`；`approval.Prompt` 参数自带 `CorrelationID/TaskID/Tool/Level/Reason/Paths`，ui.go:27-39 现读）—— 把活窗口记成名单，`Update`（`:1142`，11:46 新 HEAD 重读；Dismissed/Started 事件由 gate.go:288/:306/:520 发出）负责销记；再让 `liveVerdicts`（panel_pump.go:58-76）并读。**零新增导出名、零 approval 包改动。**
- 撞钉判定：写面全在 cmd/wisp —— 不在钉射程。internal/panel 侧不改结构 ⇒ 键面/路由面均零碰撞。
- 与票 220 措辞的差：AC#2 乙原文写的是「让 L1 也进 `LiveApprovals()` 的视图（不动方法名）」—— 现读证伪了字面形：**Queue 不引用 Gate（gate.go:54-69 只有 Gate→Queue 单向）**，「视图」只能落在 assembly root（cmd/wisp）或 Gate 自身，所以乙的真实可执行形就是本候选。票 220 复查项「票 146 对 LiveApprovals 的既有判据不被顶红」：天然满足 —— `LiveApprovals` 本体（pending_read.go:104-122）一字不动。
- 代价：中 —— ①cmd/wisp/** 是编排者写腿活跃地界（同文件互斥，票 220 禁区同源）；②「名单 + Event 销记」是新的生命周期正确性，忘销＝假阻塞，正控要含否决支路（gate.go:288 的 Dismissed 与 :306 的 Started 两种闭合）。

### D4 候选四：`LiveApproval` 增带 `Names`（键统一到读方向）〔建了但没接：names 数据在队列里躺着（queue.go:42），读方向没出口〕

- 动哪些文件：`internal/agent/approval/pending_read.go:46-50`（结构体加字段）与 `:115-119`（append 处拷入，**必须过 `cloneBacking :71-76`**）；消费侧 `cmd/wisp/panel_pump.go:64-73` 或 `internal/panel/pump.go:59-90`（NativeVerdict 增 Go 字段供 join，**不进** ApprovalCardView 的 JSON 键）；join 点 `subagent_roster_197.go:190-195`。
- 撞钉判定：Go 方法名/非 registry 类型的 Go 字段 —— 不撞 D0 钉。**但**票 146 的反形钉会动：`TestLiveApprovalsSharesNoReferenceSlotWithTheQueue`（ticket146_..._test.go:208）按**反射走 tools.Decision**，本候选不触 tools.Decision 形状（Decision 仍 8 槽），新槽 `Names []string` 在 LiveApproval 层需自行补克隆断言，否则 pending_read.go:25-38 头注的「八槽」承诺与实际交出的引用槽数目不符 —— 这是纪律账，不是钉。
- 判据生产可达：是；种形同 D1。
- 代价：中；且「给导出结构体加字段」是否算票 220「不新增方法名」的同族要编排者裁定（AC#1 字面只禁「新增导出**方法**」）。对 L1 半仍零改善。

### D5 候选五：`blockedOnApproval` 布尔 → 列表/枚举〔仅形状，别当免费〕

- 动哪些文件：`subagent_roster_197.go:127`（字段类型）＋ `:210`（赋值）＋ 两枚解码镜像 `internal/panel/subagent_roster_197_test.go:37/:57`、`cmd/wisp/subagent_carrier_197_test.go:82/:91` ＋ 消费那格的界面侧（界面地界，本普查不伸）。
- 撞钉判定：键名不变就不入键面扫描；TaskRowView 不在三对双向核（approval_test.go:118 只列 ApprovalCardView/ResultChunk/Snapshot）里 —— **但这是线协议形状变更**，票 220 AC#4 只禁「改成 D43 态名／新造态名」（D43 现名 `Confirming/AwaitingApproval` 在 `internal/statemachine/states.go:21-22`，票面 :16 行号本轮复核为真），改成列表未获票面授权 ⇒ 需编排者先裁。
- 判据生产可达：是（同 D1 种形，格子里能数出几张卡）。
- 代价：跨队协调（Go 测试镜像＋界面消费）最大；收益是把「第二张卡」从补丁变成显式事实。

**组合提示（只述形状）**：D1 单独落地 = AC#1 绿、AC#2 仍红（L1 恒不可见）；AC#2 必须至少配 D2 或 D3 之一；D4/D5 是 D1 的加重版，解决的是同一族两形。

## E. 票 221 顺带格

### E1 `TaskRoster.Cancel` 位置与生产零调用者点数〔已证（现读）〕

- 定义：`internal/tools/task.go:442-463`（`func (r *TaskRoster) Cancel(taskID string) (bool, string)`，函数头注 `:438-441`：只停这一行、不碰父与兄弟）。
- 现跑 `grep -rn '\.Cancel(' --include='*.go' | grep -v _test` 全量点数：`cmd/balldebug/main.go:339/:475`、`cmd/wisp/approval_always.go:100`、`cmd/wisp/run.go:678`、`internal/agent/loop.go:300/:346/:531`、`internal/observe/goroutine.go:324`、`internal/observe/logging.go:116` —— **9 处命中全是 `context.CancelFunc` 的 root.Cancel 或 `RunningTask.Cancel`（loop.go:300 定义体 `t.root.Cancel()`），对 `TaskRoster.Cancel` 的生产调用者 = 0 枚**。票 221 现量表点名唯一非测试命中是注释：`internal/tools/subagent_197.go:38-39`（「the only cancel that exists … read back by TaskRoster.Cancel」）—— 复核为真。
- 句柄今天是活的：`AttachCancel` 生产调用点 `subagent_197.go:330`（定义 task.go:415），`DetachCancel` 子代理 join 时 `subagent_197.go:391`（定义 :429）。⇒ **状态＝〔建了但没接〕：能停、有人给停了留了句柄，就是没人拧。**

### E2 模型读到的谎言与注册表现读 〔已证（现读）〕

- 许诺逐字：`internal/tools/subagent_197.go:189`「它会在任务名册里留下一行有父子关系与状态的记录，可以用 task.cancel 单独停它；」（Description 函数 `:187-192`）。
- `BuiltinTaskEntries` `internal/tools/task.go:595-597`：**今天只返回一枚** —— `{Tool: taskOutput{d: d}, Decl: TaskOutputDecl()}`（task.output，Execute `:498`）。头注 `task.go:18-24` 把 `task.list`／`task.cancel` 标 DEFERRED 并引 PLAN.md §7 `:1531`。

### E3 注册一枚新工具要动的面（逐面，只给形状与代价）

1. **代码 registry 面**：`internal/tools/task.go:595` 追加 Entry + 新 `taskCancel` Tool 实现 + Decl（参照 taskOutput 形状 :493-597）。装配根**无需改**：`cmd/wisp/run.go:424-430` 已 `for … BuiltinTaskEntries … reg.Register(e)`（`Registry.Register` internal/tools/registry.go:95）；模型侧 listTools 自动带出（`internal/agent/loop.go:642 localListTools` 调用位）。
2. **权力面（真正的成本，不在注册）**：「只有父（或用户）能停它的孩子」需要**调用者身份进得 Execute** —— 现读 `TaskDeps` 只有 `{Roster, Paths}`（task.go:30-38），且 `taskOutput.Execute` 无任何调用者概念（:498-512）；bridge 的 execute ctx 只注入过 `hostPathBoxKey`（bridge.go:596，全仓 `WithValue` 现读仅此后花园键）。⇒ 要新增一枚身份载体（ctx key 或既有参数改道），这是甲形的实质工作量。
3. **权威表/登记面（契约级，人工批准）**：D34 内置工具权威表 `docs/PLAN.md:2564`（`task.list / task.cancel` 行，L0/L1、S7）；§7 DEFERRED 登记 `docs/PLAN.md:1531`（完成判据＝「注册表里各有一枚真实现并过契约测试」）；S7 验收行 `docs/PLAN.md:3115` 也点了 task.list/task.cancel。票 221 处置：PLAN.md 冻结文字一字不动、先落 `A##` 台账（先例 A399/A401）。⚠ 现读 `grep DEFERRED(` 全仓 tools 包只命中 bridge.go:687 与 doc.go:47 —— task.go 的 DEFERRED 是**注释散文**（:18-24）不是 `DEFERRED(D-xx)` 标记，AGENTS §1.1 的 1:1 双向对表按哪种口径执行，归编排者裁（见末节）。
4. **C17 面**：task.cancel 属模型侧 C1 工具，**天然不碰 C17**；C17 `PanelBridge` 冻结位在 `docs/PLAN.md:1367`（定义行）、`:2434`（方法白名单要求）、`:2968`（D39 补四项）。只有面板要加「停子任务」按钮才升级为 C17 方法名（票 221 禁区：不新增 C17 方法名 —— 所以甲形今天不需要它）。
5. **终态面**：停掉后名册行落 D43 已有名（`states.go:21-22` 邻域，20 枚名 `:43` 清单），不新造 —— 与票 220 AC#4 同源约束。

## 我没查清／查不动的

1. **L1 窗口开／关是否已触发泵 publish**：`consoleApprovalUI.Prompt`（run.go:1093）之外是否挂了快照发布（panel_pump.go 注释 `:393-397` 只说「不在 delta 上发」）—— 未逐行核。这决定 D2/D3 落地时「那一格何时变真」要不要补发布触发点。
2. **`NativeVerdict.CardView()`（pump.go:88）内部把 `NativeVerdict` 换成 `risk.Decision` 的转换体**没逐行读（只读了输入形状 approval.go:72-100 与调用点 pump.go:199-202）；D2 给 L1 合成行时要先看那枚转换吃什么。
3. **`agent.newTaskID` 与 roster 行的登记时序**（PublishSubagent 精确行号）只到注释级（subagent_blocked_197_test.go:103-105 转述），没读实现行；不影响 B/C 结论。
4. **queue.go:153 用 `d.CorrelationID` 而非生成名**推得「空 corr 永不进撞名支」，靠的是 :149-151 生成名互异——**没有**针对该支路的现存测试可读来背书（推断，非尺）。
5. **`DEFERRED(D-xx)` 代码标记 vs 注释散文**：AGENTS §1.1 的「1:1 双向」对表以哪种为准未查 SPEC-12 §5 原文（本轮射程内只证了 tools 包现读 `grep DEFERRED(` 命中两枚且都与 task.cancel 无关）。
6. **界面那一跳**：修完之后「界面上还差哪一跳」（票 220 AC#5 末句）—— frontend/design 零接触硬约束，未查也不转述。
7. **写腿并发（本轮真实发生，非假设）**：采集期间 HEAD 从 `6233dedc` 前移（票 201-r2 交件）。已于 11:46 现跑 diff＋重读校正：`cmd/wisp/run.go` 的引用（:424/:519/:1085-1095 原位成立，`Update` 由 :1139→**:1142**）；其余被引文件在该窗口零改动。**但派修时若 HEAD 再动，应以下表命令重跑一遍 A/C 节行号，不得引用本件行号当现读。**
8. `subagent_blocked_197_test.go` 种卡手法的精确起点现跑为 `PendingApproval(ctxCard…)` @ `:120`（D1 引 `:120-126` 段）；`PublishSubagent` 那枚注释在 `:103` 一带 —— 两处范围是行级近似，未逐行展开。

---
*（只读普查件。票 220 判据、修法选择与台账登记归编排者。）*
