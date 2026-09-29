# nail-1 撞钉预检卡（票 224 / 220 / 213 / 214）

- 只读腿 `nail-1`。锚点 `git rev-parse --short HEAD` = **`98df640a`**（起手现跑）。行号一律按此锚点。
- 起手工作树名册已存 `git-status-start.txt`（153 行，本就含他人脏改动）。逐包名册存 `rosters-by-package.txt`。
- **⛔ 全程零编译／零 `go test`／零 `go build`／零 `go vet`／零 commit／零 push**（本机 self-hosted CI 在跑，另有裁决腿 `235-v1` 占突变窗口）。
- **硬约束遵守**：`frontend/**` 与 `design/**` 零读零引零枚举；不读 `PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`；三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）**一字不读** —— 其射程结论一律引自既有只读普查件（`224-c1`／`213-c1`／`220-c1`）＋整包 `-v` 日志的**用例名/状态**（读日志不算读冻结件），非我亲读源码。
- **四张票面文件名与 `-done` 后缀（起手 `ls`＋`grep` 核对，均无 `-done`＝未结、待派）**：
  - `213-composer-plus-menu-slash-command-catalog.md`
  - `214-composer-plus-menu-attachments-and-context-wired.md`
  - `220-roster-blockedOnApproval-only-sees-L2-and-second-card-silently-drops.md`
  - `224-session-scoped-grant-has-zero-executors-and-no-session-identity.md`
- **今日整包基准（引 `.scratch/wisp/probes/235/orch-rerun/readings.md`＋我从 `full.txt` 复认）**：包级红 4（`cmd/wisp`/`internal/ball`/`internal/panel`/`internal/risk`），名级红 **7**（`grep -cE '^[[:space:]]*--- FAIL'`＝7，无缩进子测试红）。逐包 PASS/FAIL/SKIP（缩进敏感尺，含子测试）：
  - `cmd/wisp` PASS=190 FAIL=1 SKIP=0 ｜ `internal/panel` PASS=148 FAIL=4 SKIP=0 ｜ `internal/risk` PASS=196 FAIL=1 SKIP=1 ｜ `internal/ball` PASS=54 FAIL=1 SKIP=0 ｜ `internal/agent` PASS=94 FAIL=0 ｜ `internal/agent/approval` PASS=56 FAIL=0 SKIP=1 ｜ `internal/tools` **PASS=218 FAIL=0 SKIP=0（全绿＝235-v1 所在包）** ｜ `internal/memory` PASS=67 FAIL=0 SKIP=1 ｜ `internal/perm` PASS=14 FAIL=0。
  - 7 枚红逐名（复认与 readings.md 逐字一致）：`cmd/wisp/TestTicket223HandEditedFsLooseningCostsAnL2Card`、`internal/ball/TestC21TableColourRowsMatchTokensCSS`、`internal/panel/TestApprovalCardViewJSONKeysMatchFrontendTypes`、`internal/panel/TestComposerContractTypesMatchFrontend`、`internal/panel/TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`、`internal/panel/TestC21DesignTokensFourWayAgree`、`internal/risk/TestResolvePerCallBudget`。归因：`internal/panel` 4 枚（含冻结件 `TestC21DesignTokensFourWayAgree`，`Q-52` 已撤＝已知常红）＋`internal/ball` 1 枚＝**在册常红 5**；`cmd/wisp/TestTicket223…` 与 `internal/risk/TestResolvePerCallBudget`＝**计时/争用型 2**（隔离 `-count=3` 各 3/3 绿，见 `iso-223-handedit.txt`／`iso-resolve-budget.txt`）。⇒ **写腿派单里凡列"今天红着的用例"，这 5 枚常红要写明"本就在册、非你弄坏"，2 枚计时红要写明"整包会偶红、安静复量转绿"。**

---

## 卡 1 — 票 224（会话身份候选）

> 票面：`224-session-scoped-grant-has-zero-executors-and-no-session-identity.md`（待派；前置 `A435` 已批：会话＝什么已定、走乙形"宿主随机铸造、不许派生、结束点＝进程退出"；`grants.list`/`grants.revoke` 不在射程，会被 `l2_grant_boundary` 禁名词表直接判死）。只读普查：`.scratch/wisp/probes/224/c1/census.md`（86,894 字节，结论已读）。

### 1. 这轮要新增/改动的具名符号
- **具名符号**：新铸造函数（普查候选名 `mintSessionID`，**今天不存在**）；新只读注入 seam（候选 `GrantSource`/`SessionIDSource`，照 `ModeSource` 那形）；答复词表新增"带会话档的 allow"一枚（现 `Answer` 闭集＝`AnswerAllow`/`AnswerReject`/`AnswerVeto`/`AnswerTimeout`，`internal/tools/gate.go:77-88`）；`agent.ToolRequest` 拟新增会话位（候选 `SessionID`，`internal/agent/tools.go:54` 今天无）；`tool_call.grant_id` 联动（`internal/agent/journal.go:101` 今天硬写 `nil`）。
- **既有（〔建了但没接〕，非本轮新造）**：`InsertGrant`/`ListGrantsBySession`/`RevokeGrant`/`DeleteGrant`（`internal/memory/dao_misc.go`）、`GrantScopeSession="session"`（`models.go:105`）、`ApprovalGrant.SessionID`（`models.go:98`）、`DecisionAllowGrant="allow_session_grant"`（`internal/agent/journal.go:32`，决策词非答复词、零写者）。
- **行为关键词（票 AC 给的行为）**：AC#1 身份唯一铸造点(0→≥1)、AC#2 写/读/失效三件套(读＝"同会话命中→不再弹卡")、AC#3 重启必失效、AC#4 反控"派生式 id 让授权跨会话仍生效"、AC#5 不许把"长期"混进。

### 2. 正向撞钉尺（逐符号命中；尺＝`grep -rn "<符号>" --include=*_test.go internal/ cmd/ tools/`）
| 符号 | 命中（file:line／那一行在判什么） |
|---|---|
| `InsertGrant` | dao_test.go:383 写行、:391 拒非 `session` scope；**ticket90_persist_test.go:230/239/242（冻结件，仅列名不转述）**；run_mode101_test.go:411 装配内写行、:420 Fatal |
| `ListGrantsBySession` | dao_test.go:394 查；**ticket90:244/258/260（冻结）**；run_mode101:422/424/453/455 同会话查到 1 行 / 新会话查 0 行 |
| `RevokeGrant` | dao_test.go:398/401/408（幂等＋ErrNotFound）——仅 DAO 自测 |
| `DecisionAllowGrant` | **0** |
| `GrantScopeSession` | dao_test.go:384；**ticket90:231（冻结）**；run_mode101:412 |
| `allow_session_grant` | dao_test.go:338 `DecideToolCall(...,"allow_session_grant",&grantID)`、:352 回读 decision 串——DAO 层已通 |
| `DeleteGrant`／`ListAllGrants` | 0／0 |
| `SessionID` | dao_test.go:385/391；**ticket90:234/271/272（冻结）**；run_mode101:415/499 |
| 候选新名 `GrantSource`/`SessionIDSource`/`MatchGrant`/`mintSessionID`/`mintGrant`/`GrantMatch` | **全 0**（`SessionGrant` 只命在测试**函数名** Ticket90/Ticket101，非 Go 标识符）⇒ 用这些新名不会撞既有测试引用 |

### 3. 行为型负向钉（关键词扫：`NeverSilenced`/`免审`/`通行证`/`no.*grant.*source`/`不.*静默`）
尺：`grep -rn --include=*_test.go -E 'NeverSilenced|免审|通行证|grant source' internal/tools internal/risk internal/agent cmd/wisp`（排三枚冻结件）。逐枚读正文后的射程判断：
- ⛔ **`cmd/wisp/run_mode101_test.go:389 TestTicket101SessionGrantDoesNotCrossRestart`（非冻结）**：`:438-441` `if !firstErr { t.Fatalf("boot 1 already let the B-tier .env write through …") }` ＋ `:427-428` 注释逐字 **"the assembly hands the bridge a mode, never a grant source (Options.Confirmations nil)"** ⇒ 这枚钉**把"同会话活授权仍必须被拒"钉成期望**，与 224 AC#2 的"读"（命中→不问）**极性相反**。它今天绿，只因测试自打的字面量 `session-before-restart`（`:400`）≠ 生产会铸造的任何 id。**这是本票最硬的撞面**（详见末节）。
- `internal/tools/ticket90_test.go:357 TestTicket90TaintEscalationNeverSilenced`：R4/污染 verdict 在每一档都必须 raise L2、`SessionOverrideBlocked` 不得丢。⇒ 224 的"命中即不问"这步**绝不许静默 R4/污染/Deny**（红线不可覆盖），实现须保它绿。
- `internal/tools/ticket90_test.go:431 TestTicket90ModeCarriesNoAllowAuthority`：反射核 `tools.Decision` 字段名**不得含 allow/approve/grant、且不得是 `Answer` 型**；`agent.ToolRequest` 字段名**不得含 mode/allow/approve/auto**；`ModeSource` 恰 1 方法。⇒ ① 给 `ToolRequest` 加 `SessionID` 安全；② **不许把"命中"结果存进 `tools.Decision` 的 grant/allow 字段**（要落在 `Silenced.Kept` 或作 route 形参）；③ 新 `GrantSource` 接口不受影响（只核 `ModeSource`）。
- `internal/memory/dao_test.go:391` 拒 `Scope:"global"` ⇒ 224 只可保留 `session` 一档，扩档即红。
- ⚠ **冻结件（不读，引普查）**：`ticket90_persist_test.go:220` 纯 DAO、按字面量查，看不见桥改动；`l2_grant_boundary_test.go` 会判死 `grants.list`/`grants.revoke` 方法名、判死入向新类型的 `Grant`/`Allow` 字段（`Scope` 不在词表＝安全）、并因 registry 腐蚀守卫使"在 `internal/panel` 新增可解码结构体"走不通 ⇒ **落点必须避开 `internal/panel`**（`A435` 已定）。

### 4. 今天绿着的用例名名册（对应包）
- `cmd/wisp` PASS=190／FAIL=1／SKIP=0。相关绿（务必抄进派单让写腿认得出）：`TestTicket101SessionGrantDoesNotCrossRestart`(3.17s)、`TestTicket101ModeSwitchUsesTheRealL2Gate`、`TestTicket101ManualSwitchSurvivesRestart`、`TestTicket101UntouchedConfigRestartsAtDefault`、`TestTicket101UnreadableModeFailsLoudlyAndStrict` 全 PASS。今日红＝`TestTicket223HandEditedFsLooseningCostsAnL2Card`（**计时型，隔离 3/3 绿**）。
- `internal/tools` **PASS=218／FAIL=0（全绿）**。相关绿：`TestTicket90TaintEscalationNeverSilenced`、`TestTicket90ModeCarriesNoAllowAuthority`、`TestTicket90IrreversibleStillAsksInEveryMode`、`TestTicket90TierADenySurvivesEveryMode`、`TestTicket90TaintFlagAndDenyAreNeverSilenced` 全 PASS。**⚠ 这枚包就是 235-v1 在动的包。**
- `internal/risk` PASS=196／FAIL=1／SKIP=1。相关绿：`TestFusionR4BlocksSessionOverride` PASS。今日红＝`TestResolvePerCallBudget`（**计时型，隔离 3/3 绿**）；SKIP＝`TestSyncRegistryProbeLive`。
- `internal/memory` PASS=67／SKIP=1。`TestGrantCostPluginStateDAO` PASS。
- `internal/perm` PASS=14／FAIL=0（含冻结件 5 枚 `TestTicket90…` 全绿）。

### 5. 与 235-v1 撞面判定
⛔ **必须排在 235-v1 之后／不能并发**。224 的最小落点按普查 B3＝**`internal/tools/bridge.go:288` 与 :289 之间**插"命中降 L0"一步；seam 注入点＝**`cmd/wisp/run.go` 的 `assembleRuntime`**（`rt.store` 与 `rt.bridge` 同处此地）。`internal/tools/**` 与 `cmd/wisp/run.go` **双双命中任务点名的敏感文件**，且 235-v1 正在 `internal/tools` 占突变窗口 ⇒ **同包写腿会互相洗读数，串行**。票 224 禁区亦自陈"与票 201/222/223 同撞 cmd/wisp/internal/agent ⇒ 串行"。

### 6. 门禁射程（d22scan）
- 调用形（⛔ 我未跑）：`bash scripts/d22scan.sh`（历史派单亦写 `sh scripts/d22scan.sh`；⚠ 其日志混扫描器自检 `examined N/1` 行，**别用 `grep -m1 examined` 取分母**）。
- emoji 实际射程（引 AGENTS.md 定案）：扫 `U+1F000–1FAFF`·`U+2200–22FF`·`U+2600–27BF`·`U+2B00–2BFF`·`U+FE0F`·`U+1F1E6–1F1FF`；**注释豁免、字符串不豁免**；箭头 `U+2190–U+21FF`、带圈数字 `U+2460–U+24FF` **不扫**。⇒ 224 若在 Go 侧落"本会话内允许"卡文案：`✓`(U+2713)、`≤`(U+2264) **会被点红**，`→` 不会。
- 路径判定禁令：`filepath.Clean|Abs` 用于 FS 决策（在 `risk.PathResolver` 外）＝即判违规；224 若新增 `approval_grant.pattern` 的路径匹配器**必须在 PathResolver 内**，否则 d22scan 红。

---

## 卡 2 — 票 220（名册 `blockedOnApproval` join 键与 L1 枚举器）

> 票面：`220-…-second-card-silently-drops.md`（待派）。只读普查：`.scratch/wisp/probes/220/c1/census.md`（结论已读；路径勘正：roster 件在 `internal/panel/subagent_roster_197.go`，tools 侧无同名文件）。

### 1. 这轮要新增/改动的具名符号
- **具名符号**：`BlockedOnApproval`(Go 字段)/`blockedOnApproval`(JSON 键)（`internal/panel/subagent_roster_197.go:127`）；join 建索引 `waiting[card.CorrelationID]`（`:190-195`）与查 `waiting[row.TaskID]`（`:210`）；queue 重贴 `fmt.Sprintf("%s#%d", d.CorrelationID, q.seq)`（`internal/agent/approval/queue.go:152-154`）；`LiveApprovals()`（`pending_read.go:104-122`）/`LiveApproval` 结构体（`CorrelationID`/`Decision`/`Position`）；`Gate.windows`（`gate.go:64`）+`window` 结构体（`:72-77`）+`openWindow/closeWindow/Veto`；**AC#2 甲若要新增 L1 只读枚举方法＝新增导出方法名**（须先落台账，且不得写成 `panel.*` 字面量）。AC#2 乙形（候选 D3）＝在 `cmd/wisp/run.go` 的 `consoleApprovalUI` 记活窗口名单，零新增导出名。
- **票给的行为关键词**：join 双查(corr+taskID)、`#序号` 后缀匹配、L1 可见、两形各配正控、不许把布尔改成 D43 态(`Confirming`/`AwaitingApproval` `states.go:21-22`)、界面"还差哪一跳"只写票不代做。

### 2. 正向撞钉尺（`grep -rn "<符号>" --include=*_test.go internal/ cmd/ tools/`）
| 符号 | 命中（判什么） |
|---|---|
| `BlockedOnApproval`/`blockedOnApproval` | **正向钉**：subagent_roster_197_test.go:57(镜像字段)/:164 `if !child.BlockedOnApproval`(child 应有)/:169 `if root.BlockedOnApproval`(root 应无＝负向)；subagent_blocked_197_test.go:170/:190/:223/:237（真 packet 上断 child 真/root 假/**答复后 afterRow 必须假**＝负向）；subagent_carrier_197_test.go:82/:311。⚠ 这些单卡正/负向钉的**射程直接盖住 D1 的 join 改动**（见 3）。 |
| `LiveApprovals` | pending_read_test.go:111/121/133/158/170/201/226（ReportsPending…、Skips…、OnNilQueueIsNil）；ticket146_liveapprovals_backing_test.go:208/264/361（反形钉，见 3） |
| `Gate.windows`/`windows[` | 非该字段——命中的是别包同名：ticket90_persist_test.go:80/94、bridge_test.go:84/98/104/107、ticket90_test.go:116/138/230/234（`t90Gate` 的计数窗，非 `Gate.windows`）⇒ **L1 枚举器今天无任何正向钉，也确实无枚举口** |
| `#序号`/`re-stamp`/`clash` | ticket87_veto_l2_test.go:127（重贴后仍可按 host 原 key 落到**拒**方向）；subagent_stream_197_test.go:75 `clash`（流键，另一义）。⇒ 拒方向有钉、读方向无钉（票 220 现量"第二张卡静默掉格今天没有尺看得见"复核为真） |

### 3. 行为型负向钉（关键词扫：`SpellsNoRoute`/`SharesNoReference`/`second switch`/`inbound`/`闭集`/`fifth`）
尺：`grep -rn --include=*_test.go -E 'SpellsNoRoute|SharesNoReference|second switch|inbound|fifth' internal/panel internal/agent/approval`（排冻结件）。判断：
- ⚠ **`internal/panel/composer_dispatch_test.go:373 TestDispatcherSpellsNoRouteLiteralOfItsOwn`（非冻结，PASS 今）**：正则扫 `composer_dispatch.go` 里 `"panel\.[^"]*"` 字面量即判红。⇒ 若 220 的 L1 枚举行在 `internal/panel` 侧写成 `panel.*` 字面量＝红；用 bridge 常量或纯 Go 枚举名（放 approval/cmd 包）则不入射程。
- ⚠ **`internal/agent/approval/ticket146_liveapprovals_backing_test.go:208 TestLiveApprovalsSharesNoReferenceSlotWithTheQueue`**：`typeCensus` 走的是 **`tools.Decision`**（`:218`），核对 `queueSlotPaths` 手维护集，fail-closed（`:135-136` "新增字段类型要连这把尺一起改"）。⇒ **候选 D4（给 `LiveApproval` 加 `Names []string`）不触此钉**（Decision 仍 8 槽），**但**此钉**不保护 `LiveApproval.Names` 本身的底层共用** ⇒ 写腿须自补一枚"Names 不共享队列 backing"的正控，否则是票 146 哲学下的纪律账/未防形。
- ⚠ **join 负向钉（D1 的射程）**：subagent_blocked_197_test.go:190 `if rootRow.BlockedOnApproval`(root 无卡→必须假)、:223 `if afterRow.BlockedOnApproval`(答复后→必须假)；roster_test.go:169 同。**D1 改查询为"corr∪taskID∪`taskID#`前缀"时，若前缀匹配放得过宽，这些"应当假"的负向控制会翻真＝打红。** 写腿要保 root/after 支路仍假。
- ⛔ 冻结件 `l2_grant_boundary`（不读，引普查 220-c1 D0）：AST 只扫 `internal/panel` 目录；`blockedOnApproval` 键含 `approval` 词根但 `TaskRowView` 是**出向**、不在 4 枚 inbound registry ⇒ 今天不被扫。**红线＝新增字段落进那 4 枚入向类型、或新增 `panel.*` 路由字面量，才会红。**

### 4. 今天绿着的用例名名册（对应包）
- `internal/panel` PASS=148／FAIL=4／SKIP=0。相关绿：`TestTheRosterReaderPutsSubagentsOnTheWire`、`TestDispatcherSpellsNoRouteLiteralOfItsOwn`、`TestRosterMismatchBackstopRefusesInsteadOfAccepting`、`TestAC9InboundLegRefusesRosterMethodWithNoHandler` 全 PASS。**今日 4 枚红（在册常红，非写腿弄坏）**＝`TestApprovalCardViewJSONKeysMatchFrontendTypes`/`TestComposerContractTypesMatchFrontend`/`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`/`TestC21DesignTokensFourWayAgree`。
- `internal/agent/approval` PASS=56／FAIL=0／SKIP=1。相关绿：`TestLiveApprovalsReportsPendingRowsAndDoesNotConsumeThem`/`…SkipsRowsThatAreNotPending`/`…OnNilQueueIsNil`/`…SharesNoReferenceSlotWithTheQueue`/`…InPlaceWriteCannotReachTheQueueRecord`/`TestGrantIsSingleUseAndBoundToItsItem`/`TestAnAliasCanNeverBuyAnAllow`/`TestAVetoNamedByTheHostsOwnKeyStillRefusesThatCard` 全 PASS。SKIP＝`TestDefaultDeadlineWallClockMeasurement`。
- `cmd/wisp` PASS=190／FAIL=1。相关绿：`TestRunPacketMarksTheRosterRowACardIsHolding`(1.48s)、`TestRunPacketCarriesTheSubagentItsRosterRowFed`(1.41s)。今日红＝`TestTicket223…`（计时型）。

### 5. 与 235-v1 撞面判定
- **D1（只改 `internal/panel/subagent_roster_197.go`）**：不碰 `internal/tools`／`cmd/wisp/run.go`／`bridge.go` ⇒ **与 235-v1 无直接文件撞面**（但同 `internal/panel` 写面若另有 panel 写腿需互斥）。
- **D2（改 `gate.go` 补枚举口）**：写 `internal/agent/approval/gate.go`——非 235-v1 的 `internal/tools`，但票 201-c2 写面同区 ⇒ 串行 201。
- **D3（cmd/wisp/run.go 记名单）**：⛔ **命中 `cmd/wisp/run.go`＝任务点名敏感文件 ⇒ 与 cmd/wisp 写腿（224/213）串行**。
- 票 220 禁区原文：**"票 201 的写面（`internal/agent/approval`＋`internal/tools`＋`cmd/wisp`）未空出之前不许派本票"** ⇒ 无论走哪支，**须排在 201/235 家族之后**（`internal/tools` 被 235-v1 占＝禁区 blanket 串行）。

### 6. 门禁射程（d22scan）
调用形 `bash scripts/d22scan.sh`（勿 `grep -m1 examined`）。220 Go 侧无界面文案（界面那跳在禁区外、只写票不代做）⇒ emoji 风险低；但若在 Go 侧给名册行加"还在跑/在等批准"人话注记，`✓`/`≤` 在射程内、`→` 不在。**路径判定禁令不适用本票**（不碰 filepath 决策）。

---

## 卡 3 — 票 213（加号菜单·斜杠命令目录，Go 侧前置）

> 票面：`213-composer-plus-menu-slash-command-catalog.md`（待派，owner 09-28 口头接单＝功能请求）。只读普查：`.scratch/wisp/probes/213/c1/inventory.md`（35,987 字节，结论已读）。票面 §09-29 11:3x 已自报一枚挡住本票的钉并定案"不新建登记表、复用 control.go"。

### 1. 这轮要新增/改动的具名符号
- **具名符号（复用为主，票面明令不新造名类）**：控制词解析族 `internal/agent/control.go`（`ControlVerb:17`/`controlWords:36`/`MatchControl:51`/`ControlOutcome:61`/`ControlHandler:73`）；装配落点 `cmd/wisp/run.go` 里 `opts := agent.Options{ … Control: … }`（`Options.Control` 定义 `internal/agent/loop.go:158`，生产今天没人设）；取消命令内部调 `internal/tools/task.go:442 (*TaskRoster).Cancel`（把手已由 `subagent_197.go:330 AttachCancel`/`:391 DetachCancel` 存好）；撞名响亮报照 `internal/tools/registry.go:103 "duplicate tool name (C1 requires global uniqueness)"`。
- **命令 token（票交付 1 初始集，多数能力未落地＝不许先进册）**：`compact`/`goal`/`plan`/`model`/权限预设/`导出`/`反馈`/`取消`。
- **禁止新增的符号**：**新的 C17 入向方法名**（`internal/panel/bridge.go:42-45` 是 4 枚闭集：`MethodModeRequest`/`MethodWorkspaceRequest`/`MethodAttachmentAdd`/`MethodMessageSend`）、`panel.*` 路由字面量、`internal/panel` 里的第二枚"命令登记表/switch 链"、把 `/compact` 等做成模型可自调工具。
- **票给的行为关键词**：AC#1 每条有真执行者、AC#2 没能力进不了册(正控塞假命令必红)、AC#3 档位权限(拒要说为什么)、AC#4 入向没接通前"界面点不了"不许翻勾、AC#5 全走现成审计、命令只有人能执行。

### 2. 正向撞钉尺（`grep -rn "<符号>" --include=*_test.go internal/ cmd/ tools/`）
| 符号 | 命中（判什么） |
|---|---|
| `MatchControl`/`controlWords` | control_test.go:61-62 `if _,ok:=MatchControl(u);ok {…bare-word rule}`（近失词不许误判）；TestControlWordsNeverReachTheProvider(:15)/TestControlWithoutAConsumerIsVisible(:107)/TestNearMissUtterancesGoToTheLoop(:59)/TestDefaultControlCancelsRunningTask(:77) 均在此文件 |
| `TaskRoster.Cancel`/`.Cancel(` | 生产钉全在 internal/tools：subagent_197_test.go:772/797/800（停孩子/停已 join 的/停查不到的）、**task_cancel_221_legs_test.go:181 "BuiltinTaskEntries 没注册 task.cancel（现名册…）⇒ 说明书又在许诺一枚不存在的工具"**、task_cancel_221_test.go:12/50 |
| `BuiltinTaskEntries` | 12+ 命中（pointer_183/185/subagent_197/task_cancel_221/task_output_*/ticket175r2），是工具名册枚数钉 |
| `MethodModeRequest` 等 4 枚 | composer_test.go:411-414、git_test.go:382/450/505、bridge_test.go:56-62/139-140 —— 白名单枚数钉（见 3） |
| `Options.Control` 生产装配 | **无钉把"生产 assembly 的 Options.Control 为空"钉成期望**（harness_test.go:75 是测试自带 setter，非生产钉）⇒ **213 在 run.go 装 `Options.Control` 不撞钉**（这是好消息，写腿可放心接） |

### 3. 行为型负向钉（关键词扫：`fifth`/`second switch`/`handler registry`/`route.*literal`/`closed set`）
- ⛔⛔ **主互斥（本票最该当场裁的一枚）＝冻结件 `internal/panel/l2_grant_boundary_test.go` 的 `poolJudgedByRealGuard`**（不读，引票 213 §09-29 原文＋普查 213-c1 §2）：它把包里**任何"第二枚 route 形状的 switch/if 链"**报成 —— 注释逐字 **"a second switch/if chain … exactly the shape a handler registry takes"**。**新建命令登记表若落在 `internal/panel`＝被这枚冻结钉直接判死。** ⇒ 票面已定案"复用 control.go、落点出 panel"。**只要写腿把命令名册/派发链写在 `internal/agent` 或 `cmd/wisp`，就不撞；一旦写回 `internal/panel`，就是 09-28 白跑那一形，须"改扫能力＋加正控"当场裁，不能等写腿写完再说。**
- ⛔ **`internal/panel/git_test.go:382-387/:446-450/:503`（非冻结，PASS 今）**：`composerMethodWhitelist` 把 4 枚方法当**一枚闭集**核，注释逐字 **"a fifth method means…"**／**"so a fifth method cannot [be added silently]"**。⇒ 213 新增第 5 枚 `panel.*` 方法＝红（票禁区已写"不新增 C17 方法名"）。composer_test.go:411-414 同族。
- ⚠ **`composer_dispatch_test.go:373 TestDispatcherSpellsNoRouteLiteralOfItsOwn`**：`composer_dispatch.go` 里出现 `"panel.*"` 字面量＝红。⇒ 213 若加执行通道须用 bridge 常量。
- ⚠ **`internal/agent/control_test.go:15 TestControlWordsNeverReachTheProvider`**＋`loop_golden_test.go:186 "not even a fifth call may have entered the provider"`：控制词/命令**绝不许走 LLM provider**。复用 control.go（非 LLM 路）保绿。
- ⚠ **`internal/agent/control_test.go:107 TestControlWithoutAConsumerIsVisible`**（现读）：用 harness 直接测 loop（`大声点`→`StatusControl`、`requests==0`），**不测 cmd/wisp 装配** ⇒ 213 在 run.go 装 Control 不影响它。

### 4. 今天绿着的用例名名册（对应包）
- `internal/agent` **PASS=94／FAIL=0**：`TestControlWordsNeverReachTheProvider`/`TestControlWithoutAConsumerIsVisible`/`TestNearMissUtterancesGoToTheLoop`/`TestDefaultControlCancelsRunningTask` 全绿。
- `internal/panel` PASS=148／**FAIL=4（在册常红）**／SKIP=0：相关绿 `TestComposerEnvelopeAcceptsItsFourRequests`/`TestDispatcherSpellsNoRouteLiteralOfItsOwn`/`TestAC9ComposerDispatchHasAProductionCaller`/`TestAC1AC2DispatchHopGate133`/`TestFrontendComposerRequestsMatchTheEnvelope` 全 PASS；**今日红 4 枚**同卡 2 列表（含 `TestComposerContractTypesMatchFrontend` 属前端契约常红）。
- `cmd/wisp` PASS=190／FAIL=1（`TestTicket223…` 计时型）。
- `internal/tools` **PASS=218（全绿）**：`task_cancel_221*`/`subagent_197` 等 cancel 相关全绿。**⚠ 235-v1 所在包。**

### 5. 与 235-v1 撞面判定
⛔ **须串行**。213 主落点＝**`cmd/wisp/run.go`（装 `Options.Control`）＝任务点名敏感文件** ⇒ 与 224/201-c2 的 cmd/wisp 写腿互斥。213 **不新增也不改 `internal/tools` 文件**（cancel 只"调用"`TaskRoster.Cancel`，不注册 `task.cancel` 工具——那属票 221 的甲形且已被 `task_cancel_221_test.go:181` 名册枚数钉守着）⇒ 与 235-v1 只在"213 若误改 `internal/tools/**`"时才撞。**派单要把这条界线写死：213 碰 `internal/tools/**`＝越界。**

### 6. 门禁射程（d22scan）
调用形 `bash scripts/d22scan.sh`（勿 `grep -m1 examined`）。**213 的名册是一等界面文案面**：每条"一句话说明"＋"要不要人批"＋"会改系统状态/只递一段话"栏＝Go 侧字符串 ⇒ **`✓`(U+2713)/`≤`(U+2264) 会被点红、`→` 不会、注释豁免但字符串不豁免**；票 216 是同族"scrub 控制字符"的相邻票。**路径判定禁令**不适用（不碰 filepath 决策）。命令若被误做成 C4 工具，另撞 registry 唯一名钉＋`internal/panel` 之外仍可能触发 d22scan 的其它 ban，但主风险是上面那枚 poolJudgedByRealGuard 冻结钉。

---

## 卡 4 — 票 214（加号菜单·附件与上下文接线，Go 侧前置）

> 票面：`214-composer-plus-menu-attachments-and-context-wired.md`（待派）。普查：同批 `.scratch/wisp/probes/213/c1/inventory.md`（含 214 落点，第 5 节）。票面 §09-29 12:1x 定案"只用现成默认、不新增配置键"。

### 1. 这轮要新增/改动的具名符号
- **具名符号**：受理器 `internal/panel/attachments.go`（`NewAttachmentBroker:179`、`(*AttachmentBroker).Ingest:199`、`matchesISOBaseMedia`、`AcceptedMIMETypes:153`、`MaxAttachmentBytes:47`、`NameGuard:165`、`ArtifactSink` 接口）；断线点 `internal/panel/pump.go:226` `NewComposerState(mode, workspace, nil, maxAttachment)`——**附件实参被写死成 `nil`，要换成真读**；装配根 `cmd/wisp/run.go` 给 `ArtifactSink` 一枚真 sink（起手"生产装配根有没有真给它 sink＝未量"）。
- **既有快照键（不许再加新键）**：`composer.go:238-241` `attachments`/`acceptedAttachmentMimes`/`maxAttachmentBytes`/`attachmentError`。
- **禁止新增**：新配置键（属 D36 树＝契约面）、新 C17 方法名、第二条文件读入口、`risk.PathResolver` 外的 `filepath.Clean|Abs`。
- **票给的行为关键词**：AC#1 受理器有生产调用者(0→≥1，在 cmd/wisp 装配路径)、AC#2 超限拒且快照见"被拒＋为什么"、AC#3 伪装按 `matchesISOBaseMedia` 嗅探判不按扩展名、AC#4 越界必红(票 174/175 族)、AC#5 来源戳 C25 现成名盖上、AC#6 没接的三件具名留票。

### 2. 正向撞钉尺（`grep -rn "<符号>" --include=*_test.go internal/ cmd/ tools/`）
| 符号 | 命中（判什么） |
|---|---|
| `NewAttachmentBroker`/`MaxAttachmentBytes`/`NameGuard` | attachments_test.go:91/132/229/232/246/270/333（含 :229 nil sink 拒、:232 nil 校验拒）；composer_test.go:215、pump_test.go:140/142/305/306、panel_pump_test.go:282 —— 正向用例，214 改动要连它们一起保绿 |
| `NewComposerState` | composer_test.go:210、**pump_test.go:301 注释逐字 "NewComposerState owns it and the pump must go through that constructor"**（＝构造器钉，见 3） |
| `AcceptedMIMETypes` | composer_test.go:258 `if len(s.Composer.AttachmentMIMEs)!=len(AcceptedMIMETypes())` |
| `matchesISOBaseMedia` | attachments_test.go 内嗅探用例（TestAttachmentPayloadCarriesTheBytes 族，PASS 今） |
| `ArtifactSink` | 见装配段现读；测试侧以 nil/真 sink 分支覆盖 |
| `Attachment` 方法/键 | subagent_carrier/instructions 无关；无独立生产调用者（票 214 现量"生产零调用者"复认） |

### 3. 行为型负向钉（关键词扫：`go through that constructor`/`SendsFourKeys`/`carry the four keys`/`inbound`/`fifth`）
尺：`grep -rn --include=*_test.go -E 'go through that constructor|FourKeys|four PanelSnapshot declares|acceptedAttachmentMimes|no_inbound|InboundEnvelope' internal/panel cmd/wisp`（排冻结件）。判断：
- ⚠ **`pump_test.go:301`（非冻结）**：快照的 `acceptedAttachmentMimes`/`maxAttachmentBytes` **必须由 `NewComposerState` 供**，泵不许绕过构造器直接塞字段 ⇒ 214 接线走构造器即绿。
- ⚠ **`pump_test.go:291-292`**：出向字节顶层键**恰**＝`composer,generatedAt,pending,results` ⇒ 214 **不许新增快照顶层键**（票已说"不加键"）。
- ⛔ 冻结件 `l2_grant_boundary`（不读，引普查）：4 枚 inbound 类型（`ComposerRequest`/`ModeRequest`/`AttachmentPayload`/`AttachmentRef`）之外的**新增可解码结构体**会触发 registry 腐蚀 `t.Fatalf`（修 registry 要动冻结件＝走不通）⇒ 214 **须复用现成 `AttachmentPayload`/`AttachmentRef`，不新造入向解码型**。相关今天绿的正控钉：`TestNoInboundEnvelopeCanBindAnApprovalVerdict`、`TestPlantedGrantWiringGoesRedInASnapshot`、`TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes`（均 PASS）。
- ⚠ **d22scan 路径判定禁**（非测试钉、属 CI 仪器）：附件落盘/名字守卫的路径决策**必须在 `risk.PathResolver` 内**；`NameGuard` 现读只"管名字合法、不管显示安全"（213-c1 §Z）⇒ 214 若自拼 `filepath`＝d22scan 判违规（`AGENTS.md:1288`）。

### 4. 今天绿着的用例名名册（对应包）
- `internal/panel` PASS=148／FAIL=4（在册常红：`TestApprovalCardViewJSONKeysMatchFrontendTypes`/`TestComposerContractTypesMatchFrontend`/`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`/`TestC21DesignTokensFourWayAgree`）／SKIP=0。相关绿（务必抄进派单）：`TestAttachmentPayloadCarriesTheBytes`（含 4 子测试全 PASS）、`TestMessageCarriesAttachmentRefsAndRefusals`、`TestSecondIdenticalAttachmentIsDeduplicatedNotOverwritten`、`TestComposerEnvelopeAcceptsItsFourRequests`、`TestComposerStateSurvivesPanelCloseAndReopen`、`TestComposerCurrentModelTravelsOnlyFromItsReader`、`TestAPumpWithoutARosterReaderSendsFourKeys`、`TestThePumpIsDrivenNotJustAssembled`、`TestNoInboundEnvelopeCanBindAnApprovalVerdict` 全 PASS。
- `cmd/wisp` PASS=190／FAIL=1（计时型 `TestTicket223…`）。相关绿：`TestTicket223PanelInboundSaysHotReloadIsDisabled`、`panel_pump_test` 附件断言。

### 5. 与 235-v1 撞面判定
⛔ **须与 cmd/wisp 写腿串行**（214 装配落点＝**`cmd/wisp/run.go`＝任务点名敏感文件**）。214 **不碰 `internal/tools`／`bridge.go`**（受理器在 `internal/panel`，与 235-v1 的 `internal/tools` 无重叠）⇒ 对 235-v1 无直接文件撞面，但 run.go 与 213/224 争同一文件 ⇒ 三者互斥、串行。若 214 的 sink 装配触及 `internal/tools` 类型定义才升级为撞 235-v1。

### 6. 门禁射程（d22scan）
调用形 `bash scripts/d22scan.sh`（勿 `grep -m1 examined`）。**214 有两道 d22scan 射程**：① **被拒附件的"人话原因"是 Go 侧字符串** ⇒ `✓`/`≤` 在射程、`→` 不在（注释豁免、字符串不豁免）；② **路径决策禁用 `filepath.Clean|Abs` 于 `risk.PathResolver` 外**（AGENTS §1.2）⇒ 接 sink/落盘那步若图省事自拼路径＝被点红。落盘根须在授权根内（票 174/175 族）。

---

## 没做完／留给编排者

### A. 正面回答"有没有哪张的票面判据与某枚现存负向钉互斥"（09-28 白跑那一形，须当场裁）

- **⛔ 票 224 —— 是，存在真互斥候选，且必须先裁再派写腿。**
  冲突对：票 224 **AC#2「读」＝"同一次会话内命中授权 → 不再弹卡"** ⟷ **`cmd/wisp/run_mode101_test.go:427-441`（非冻结、今天 PASS）** 现存负向钉："boot 1 已装 live 会话授权，`.env` 写**仍必须被拒**（`if !firstErr → Fatalf`）"，且注释**逐字钉死"assembly hands the bridge a mode, never a grant source"**。这枚钉的**期望值把 224 要推翻的"零执行者＝正确答案"写成了断言**。它今天绿**纯属侥幸**：测试打的会话字面量 `session-before-restart`（`:400`/`:453`）≠ 生产铸造的任何 id。⇒ **裁点**：① 这枚不是冻结件，**可随 224 实现改写期望并补正控**（改写它本身就是 AC#4 反控要做的"派生式 id 让授权跨会话生效必判红"）；② 派单必须写明"run_mode101 的同会话仍拒期望与 224 AC#2 相反，写腿改到它变红是**预期内的改扫**、不是弄坏"，并**先决定 224 的会话身份口径是否保证"生产铸造 id 永不等 `session-before-restart`"**（若否，这枚会静默常绿＝AC#4 假绿，正是 `store.go:19-27` 要挡的形状）。**结论：224 需"改扫能力（改写 run_mode101 期望）＋加正控（调生产铸造函数跨两启动断两 id 不等、用生产 id 写再查 0 行、走真 `Bridge.Execute` 断问与不问）"，三件缺一即恒真。**

- **⛔ 票 213 —— 原始思路与冻结钉互斥，但票面已绕开；派单要把界线钉死。**
  冲突对：票 213 交付 1 的"新建可枚举命令名册（登记表/派发链）"**若落在 `internal/panel`** ⟷ **冻结件 `l2_grant_boundary_test.go` 的 `poolJudgedByRealGuard`**（"a second switch/if chain … exactly the shape a handler registry takes"）。这枚**不能改扫**（冻结件一字不许动）。⇒ **裁点**：票面 §09-29 已定案"复用 `control.go`、落点出 panel"，因此**只要写腿把名册/派发写在 `internal/agent` 或 `cmd/wisp`，就不互斥、AC 可满足**；一旦写回 `internal/panel`，就是白跑那一形（且无解，只能改票）。派单须把"命令名册不得作为 `panel.*` 路由/第二 switch 链出现在 `internal/panel`"写成硬约束，并把"第五枚 C17 方法闭集钉（git_test/composer_test 三处）"列为禁区。**不需改扫能力，需要的是把落点钉死在 panel 之外。**

- **票 220 —— 不互斥（有构造性约束，非硬顶）。** join(D1) 受 subagent 名册的 root/after "应当假"负向控制约束（前缀别放太宽）；D4 给 `LiveApproval` 加 `Names` 不触 `ticket146`（它普查的是 `tools.Decision`），但**该钉不保护 Names 的底层共用**⇒ 写腿须自补正控；L1 枚举器不得写成 `panel.*` 字面量（`TestDispatcherSpellsNoRouteLiteralOfItsOwn`）。均为可满足约束，无须改票面判据。

- **票 214 —— 不互斥（全为可满足的构造性约束）。** 走 `NewComposerState` 构造器、不加快照顶层键、复用现成 `AttachmentPayload`/`AttachmentRef` 不新造入向型（否则撞冻结件 registry 腐蚀、而修 registry 要动冻结件＝走不通）、路径决策只走 `risk.PathResolver`（d22scan）。这些都能在不改尺前提下达成。

### B. 我没做完／做不动的（照实报，不脑补）
1. **三枚冻结件我未亲读**：`l2_grant_boundary_test.go` 的 `poolJudgedByRealGuard`／`inboundTypeRegistry`／`carriesGrantWord` 的**完整断言**只引普查（224-c1 E2/E3、220-c1 D0）；224-c1 §Z 第 7 条自陈也只读了该文五段、未通读 2380 行。⇒ 写腿动 `internal/panel` 前**应自己核 `TestRealGuardRefusesEveryAssemblableApprovalRouteName`／`TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes` 两枚**（普查没逐行）。
2. **`cmd/wisp/run_mode101_test.go` 的行号是漂移态**：起手 `git status` 显示当前 `internal/tools`＋`cmd/wisp` **工作树干净**（224-c1 采集期那两处 M 已被合并/提交），但 235-v1 仍在动 `internal/tools` ⇒ 派 224 时行号须现跑 `grep -n` 复量，勿抄本卡（本卡行号＝我对 `98df640a` 现读）。
3. **`internal/agent/approval/gate.go`/`queue.go`/`pending_read.go` 的行号**引自 220-c1（其锚点 `6233dedc`，与本锚 `98df640a` 不同）⇒ 写腿须复量。
4. **票 214 的"生产装配根有没有真给 `ArtifactSink` 一个 sink"**：票 214 现量表自陈"未量"，我亦未逐行读 run.go 装配段确认，只证受理器生产零调用者成立。
5. **`wisp panel-inbound` 到原生答复那条腿是否会把会话身份/附件带进来**（224-c1 §Z 第 5 条、213-c1 §5 前置）：三票共用的"入向线没接通 ⇒ 界面点不了"这一前置，我**未现读 `panel_inbound.go`/`panel_pump.go` 全貌**，只引普查结论。
6. **d22scan 我未跑**（硬约束）：只给可复跑调用形 `bash scripts/d22scan.sh`；emoji 射程/豁免以 AGENTS.md 定案段为准。

### C. 一句话回编排者（四张各一句）
- **224**：**有钉**——`cmd/wisp/run_mode101_test.go:427-441` 的"同会话活授权仍必须拒"与 AC#2"命中即不问"极性相反、还自陈"never a grant source"；**要改票面/改扫**：派写腿前先裁"改写这枚非冻结钉的期望＋按 AC#4 三件补正控"，并把落点锁在 `internal/tools/bridge.go`+`cmd/wisp/run.go`（⇒ 串行 235-v1）。
- **220**：**有钉但可用**——join 改动受 root/after"应当假"负向控制与 `panel.*` 字面量钉约束、D4 的 `Names` 字段无人护 backing；**不必改票面**，派单把这些约束写清即可，走 D3 时串行 cmd/wisp。
- **213**：**有一枚冻结钉正面盖着**——`l2_grant_boundary` 的 poolJudgedByRealGuard 把"panel 里第二枚 switch/registry 链"判死；**票面已绕开**（复用 control.go、名册出 panel），派单须把"名册不得落 `internal/panel`、不新增第五枚 C17 方法"写成硬约束。
- **214**：**有钉全可满足**——须走 `NewComposerState`、不加快照顶层键、复用现成 `AttachmentPayload/Ref` 不新造入向型、路径只走 `risk.PathResolver`；**不必改票面**，串行 cmd/wisp/run.go。
